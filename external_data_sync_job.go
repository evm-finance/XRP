package xrp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"qc-defi-graphql-server/internal/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ExternalDataSyncJob syncs external API data to database for enhanced screener functionality
type ExternalDataSyncJob struct {
	db             *gorm.DB
	httpClient     *http.Client
	interval       time.Duration
	running        bool
	stopChan       chan struct{}
	mutex          sync.RWMutex
	lastUpdateTime time.Time
	stats          ExternalDataSyncStats
}

// ExternalDataSyncStats tracks job performance
type ExternalDataSyncStats struct {
	LastRun            time.Time `json:"last_run"`
	TotalTokensUpdated int       `json:"total_tokens_updated"`
	IssuerNamesAdded   int       `json:"issuer_names_added"`
	MarketCapsUpdated  int       `json:"market_caps_updated"`
	VolumeDataUpdated  int       `json:"volume_data_updated"`
	APICallsSuccessful int       `json:"api_calls_successful"`
	APICallsFailed     int       `json:"api_calls_failed"`
	LastError          string    `json:"last_error"`
}

// NewExternalDataSyncJob creates a new external data sync job
func NewExternalDataSyncJob(db *gorm.DB, interval time.Duration) *ExternalDataSyncJob {
	return &ExternalDataSyncJob{
		db:       db,
		interval: interval,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		stopChan: make(chan struct{}),
	}
}

// Start begins the external data sync job
func (job *ExternalDataSyncJob) Start() error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	if job.running {
		return fmt.Errorf("external data sync job is already running")
	}

	job.running = true
	log.Printf("🔄 [EXTERNAL SYNC] Starting external data sync job (interval: %v)", job.interval)

	go job.runSyncLoop()
	return nil
}

// Stop stops the external data sync job
func (job *ExternalDataSyncJob) Stop() error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	if !job.running {
		return fmt.Errorf("external data sync job is not running")
	}

	job.running = false
	close(job.stopChan)
	log.Printf("🛑 [EXTERNAL SYNC] External data sync job stopped")
	return nil
}

// runSyncLoop runs the main sync loop
func (job *ExternalDataSyncJob) runSyncLoop() {
	// Run immediately on start
	if err := job.syncExternalData(); err != nil {
		log.Printf("❌ [EXTERNAL SYNC] Initial sync failed: %v", err)
	}

	ticker := time.NewTicker(job.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := job.syncExternalData(); err != nil {
				log.Printf("❌ [EXTERNAL SYNC] Periodic sync failed: %v", err)
			}
		case <-job.stopChan:
			log.Printf("🛑 [EXTERNAL SYNC] Sync loop stopped")
			return
		}
	}
}

// syncExternalData fetches data from external APIs and updates database
func (job *ExternalDataSyncJob) syncExternalData() error {
	startTime := time.Now()
	log.Printf("🔄 [EXTERNAL SYNC] Starting external data synchronization")

	// Fetch token data from XRPL Meta API
	tokens, err := job.fetchFromXRPLMetaAPI()
	if err != nil {
		job.updateStats("", 0, 0, 0, 0, 0, 1, err.Error())
		return fmt.Errorf("failed to fetch from XRPL Meta API: %w", err)
	}

	log.Printf("✅ [EXTERNAL SYNC] Fetched %d tokens from external API", len(tokens))

	// Update database with external data
	updated, issuerNames, marketCaps, volumes := job.updateTokensInDatabase(tokens)

	// Update statistics
	job.updateStats(startTime.Format(time.RFC3339), updated, issuerNames, marketCaps, volumes, 1, 0, "")

	// Update screener table with new data
	screenerUpdated := job.updateScreenerTable()

	duration := time.Since(startTime)
	log.Printf("✅ [EXTERNAL SYNC] Sync completed in %v - Updated: %d tokens, %d screener records",
		duration, updated, screenerUpdated)

	return nil
}

// fetchFromXRPLMetaAPI fetches token data from the XRPL Meta API
func (job *ExternalDataSyncJob) fetchFromXRPLMetaAPI() ([]models.XRPTokenData, error) {
	log.Printf("🔍 [EXTERNAL SYNC] Fetching from https://s1.xrplmeta.org/tokens")

	resp, err := job.httpClient.Get("https://s1.xrplmeta.org/tokens")
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenList models.XRPTokenList
	if err := json.Unmarshal(body, &tokenList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	log.Printf("✅ [EXTERNAL SYNC] Successfully parsed %d tokens", len(tokenList.Tokens))
	return tokenList.Tokens, nil
}

// updateTokensInDatabase updates the xrpTokens table with external API data
func (job *ExternalDataSyncJob) updateTokensInDatabase(tokens []models.XRPTokenData) (int, int, int, int) {
	updated := 0
	issuerNamesAdded := 0
	marketCapsUpdated := 0
	volumesUpdated := 0

	for _, token := range tokens {
		// Skip tokens without proper identifiers
		if token.Currency == "" || token.Issuer == "" {
			continue
		}

		// Prepare data for update
		issuerName := ""
		if token.Meta.Issuer.Name != "" {
			issuerName = token.Meta.Issuer.Name
			issuerNamesAdded++
		}

		tokenName := ""
		if token.Meta.Token.Name != "" {
			tokenName = token.Meta.Token.Name
		}

		marketCap := token.Metrics.Marketcap
		volume24h := token.Metrics.Volume24H
		price := token.Metrics.Price

		// Update xrpTokens table with external data
		updateSQL := `
			UPDATE xrpTokens 
			SET 
				issuer_name = COALESCE(NULLIF(?, ''), issuer_name),
				token_name = COALESCE(NULLIF(?, ''), token_name),
				marketcap = CASE WHEN ? > 0 THEN ? ELSE marketcap END,
				volume_24h = CASE WHEN ? > 0 THEN ? ELSE volume_24h END,
				price = CASE WHEN ? > 0 THEN ? ELSE price END,
				external_data_updated = NOW()
			WHERE currency = ? AND issuer = ?`

		result := job.db.Exec(updateSQL,
			issuerName, tokenName,
			marketCap, marketCap,
			volume24h, volume24h,
			price, price,
			token.Currency, token.Issuer)

		if result.Error != nil {
			log.Printf("⚠️ [EXTERNAL SYNC] Failed to update %s/%s: %v", token.Currency, token.Issuer, result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			updated++
			if marketCap > 0 {
				marketCapsUpdated++
			}
			if volume24h > 0 {
				volumesUpdated++
			}
		}
	}

	log.Printf("✅ [EXTERNAL SYNC] Database update: %d tokens, %d issuer names, %d market caps, %d volumes",
		updated, issuerNamesAdded, marketCapsUpdated, volumesUpdated)

	return updated, issuerNamesAdded, marketCapsUpdated, volumesUpdated
}

// updateScreenerTable syncs updated data to the screener table
func (job *ExternalDataSyncJob) updateScreenerTable() int {
	log.Printf("🔄 [EXTERNAL SYNC] Updating screener table with external data")

	updateSQL := `
		UPDATE xrpScreener s
		INNER JOIN xrpTokens t ON s.currency = t.currency AND s.issuer_address = t.issuer
		SET 
			s.marketcap = COALESCE(t.marketcap, s.marketcap),
			s.volume_24h = COALESCE(t.volume_24h, s.volume_24h),
			s.updated_at = NOW()
		WHERE t.external_data_updated > DATE_SUB(NOW(), INTERVAL 1 HOUR)`

	result := job.db.Exec(updateSQL)
	if result.Error != nil {
		log.Printf("❌ [EXTERNAL SYNC] Failed to update screener table: %v", result.Error)
		return 0
	}

	log.Printf("✅ [EXTERNAL SYNC] Updated %d screener records", result.RowsAffected)
	return int(result.RowsAffected)
}

// GetStats returns current job statistics
func (job *ExternalDataSyncJob) GetStats() ExternalDataSyncStats {
	job.mutex.RLock()
	defer job.mutex.RUnlock()
	return job.stats
}

// updateStats updates job statistics
func (job *ExternalDataSyncJob) updateStats(lastRun string, updated, issuerNames, marketCaps, volumes, successfulCalls, failedCalls int, lastError string) {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	if lastRun != "" {
		if parsedTime, err := time.Parse(time.RFC3339, lastRun); err == nil {
			job.stats.LastRun = parsedTime
		}
	}

	job.stats.TotalTokensUpdated += updated
	job.stats.IssuerNamesAdded += issuerNames
	job.stats.MarketCapsUpdated += marketCaps
	job.stats.VolumeDataUpdated += volumes
	job.stats.APICallsSuccessful += successfulCalls
	job.stats.APICallsFailed += failedCalls
	job.stats.LastError = lastError
}

// ForceSync forces an immediate synchronization
func (job *ExternalDataSyncJob) ForceSync() error {
	log.Printf("🚀 [EXTERNAL SYNC] Forcing immediate external data sync")
	return job.syncExternalData()
}
