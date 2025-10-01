package xrp

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ScreenerUpdateJob manages periodic screener table updates with AMM data
type ScreenerUpdateJob struct {
	db                *gorm.DB
	priceService      PriceServiceInterface
	updateInterval    time.Duration
	running           bool
	updateCount       int64
	errorCount        int64
	lastUpdateTime    time.Time
	lastUpdateSuccess bool
	mu                sync.RWMutex
}

// ScreenerUpdateStats provides statistics about the screener update job
type ScreenerUpdateStats struct {
	Running           bool          `json:"running"`
	LastUpdateTime    time.Time     `json:"last_update_time"`
	LastUpdateSuccess bool          `json:"last_update_success"`
	UpdateCount       int64         `json:"update_count"`
	ErrorCount        int64         `json:"error_count"`
	UpdateInterval    time.Duration `json:"update_interval"`
	NextUpdateIn      time.Duration `json:"next_update_in"`
}

// NewScreenerUpdateJob creates a new screener update job
func NewScreenerUpdateJob(db *gorm.DB, priceService PriceServiceInterface, updateInterval time.Duration) *ScreenerUpdateJob {
	return &ScreenerUpdateJob{
		db:             db,
		priceService:   priceService,
		updateInterval: updateInterval,
		running:        false,
	}
}

// Start begins the continuous update process
func (job *ScreenerUpdateJob) Start() error {
	job.mu.Lock()
	defer job.mu.Unlock()

	if job.running {
		return fmt.Errorf("screener update job is already running")
	}

	job.running = true
	go job.runUpdateLoop()

	// Check if screener table is empty and run immediate update if needed
	go func() {
		var count int64
		job.db.Table("xrpScreener").Count(&count)
		if count == 0 {
			log.Printf("⚠️ [SCREENER UPDATE JOB] Screener table is empty, running immediate population...")
			job.performUpdate()
		}
	}()

	log.Printf("✅ [SCREENER UPDATE JOB] Started with %v update interval", job.updateInterval)
	return nil
}

// Stop halts the update process
func (job *ScreenerUpdateJob) Stop() error {
	job.mu.Lock()
	defer job.mu.Unlock()

	if !job.running {
		return fmt.Errorf("screener update job is not running")
	}

	job.running = false

	log.Printf("🛑 [SCREENER UPDATE JOB] Stopped after %d updates (%d errors)", job.updateCount, job.errorCount)
	return nil
}

// GetStats returns current job statistics
func (job *ScreenerUpdateJob) GetStats() ScreenerUpdateStats {
	job.mu.RLock()
	defer job.mu.RUnlock()

	stats := ScreenerUpdateStats{
		Running:           job.running,
		LastUpdateTime:    job.lastUpdateTime,
		LastUpdateSuccess: job.lastUpdateSuccess,
		UpdateCount:       job.updateCount,
		ErrorCount:        job.errorCount,
		UpdateInterval:    job.updateInterval,
	}

	if job.running && !job.lastUpdateTime.IsZero() {
		nextUpdate := job.lastUpdateTime.Add(job.updateInterval)
		stats.NextUpdateIn = time.Until(nextUpdate)
	}

	return stats
}

// runUpdateLoop is the main update loop
func (job *ScreenerUpdateJob) runUpdateLoop() {
	log.Printf("🔄 [SCREENER UPDATE JOB] Background loop started - running initial update immediately")

	// Run initial update immediately
	job.performUpdate()

	ticker := time.NewTicker(job.updateInterval)
	defer ticker.Stop()

	// Heartbeat ticker for status logging (every 5 minutes)
	heartbeatTicker := time.NewTicker(5 * time.Minute)
	defer heartbeatTicker.Stop()

	log.Printf("⏰ [SCREENER UPDATE JOB] Timer set for %v intervals - waiting for next update cycle", job.updateInterval)

	for {
		select {
		case <-ticker.C:
			if !job.running {
				return
			}
			job.performUpdate()

		case <-heartbeatTicker.C:
			job.mu.RLock()
			if job.running {
				nextUpdate := job.lastUpdateTime.Add(job.updateInterval)
				timeUntilNext := time.Until(nextUpdate)
				log.Printf("💓 [SCREENER UPDATE JOB] Heartbeat - Status: Running | Next update in: %v | Total cycles: %d", timeUntilNext, job.updateCount)
			}
			job.mu.RUnlock()
		}
	}
}

// performUpdate executes a single update cycle
func (job *ScreenerUpdateJob) performUpdate() {
	startTime := time.Now()
	nextCycleNum := job.updateCount + 1
	log.Printf("🔄 [SCREENER UPDATE JOB] ============ STARTING UPDATE CYCLE #%d ============", nextCycleNum)
	log.Printf("🕐 [SCREENER UPDATE JOB] Cycle start time: %s", startTime.Format("2006-01-02 15:04:05 UTC"))

	// Perform the actual screener table update
	updateError := job.updateScreenerTable()

	// Update job statistics
	job.mu.Lock()
	job.lastUpdateTime = time.Now()
	job.lastUpdateSuccess = updateError == nil

	if updateError == nil {
		job.updateCount++
		duration := time.Since(startTime)
		log.Printf("✅ [SCREENER UPDATE JOB] ============ CYCLE #%d COMPLETED ============", job.updateCount)
		log.Printf("⏱️  [SCREENER UPDATE JOB] Cycle duration: %v", duration)
		log.Printf("📊 [SCREENER UPDATE JOB] Total successful cycles: %d", job.updateCount)
		log.Printf("⏰ [SCREENER UPDATE JOB] Next update in: %v", job.updateInterval)

		// Log periodic summary every 10 cycles
		if job.updateCount%10 == 0 {
			log.Printf("📈 [SCREENER UPDATE JOB] ===== 10-CYCLE SUMMARY =====")
			log.Printf("📈 [SCREENER UPDATE JOB] Total successful updates: %d", job.updateCount)
			log.Printf("📈 [SCREENER UPDATE JOB] Total errors: %d", job.errorCount)
			if job.updateCount+job.errorCount > 0 {
				successRate := float64(job.updateCount) / float64(job.updateCount+job.errorCount) * 100
				log.Printf("📈 [SCREENER UPDATE JOB] Success rate: %.1f%%", successRate)
			}
			log.Printf("📈 [SCREENER UPDATE JOB] Runtime: %v", time.Since(startTime))
		}
	} else {
		job.errorCount++
		log.Printf("❌ [SCREENER UPDATE JOB] ============ CYCLE #%d FAILED ============", nextCycleNum)
		log.Printf("💥 [SCREENER UPDATE JOB] Error: %v", updateError)
		log.Printf("📊 [SCREENER UPDATE JOB] Total errors: %d", job.errorCount)
		log.Printf("⏰ [SCREENER UPDATE JOB] Will retry in: %v", job.updateInterval)
	}
	job.mu.Unlock()
}

// updateScreenerTable performs the actual screener table update with AMM data
func (job *ScreenerUpdateJob) updateScreenerTable() error {
	log.Printf("🔍 [SCREENER UPDATE] Starting screener table refresh...")

	// FIXED: Removed destructive DELETE operation
	// Now using UPSERT pattern to preserve data if update fails

	// CRITICAL FIX: Get XRP price ONCE for the entire update cycle with retry logic
	xrpPriceUSD, err := job.getXRPLPriceWithRetry()
	if err != nil {
		log.Printf("⚠️ [SCREENER UPDATE] Failed to get XRP price after retries, skipping this cycle: %v", err)
		return fmt.Errorf("failed to get XRP price: %w", err)
	}
	log.Printf("💵 [SCREENER UPDATE] XRP price for this cycle: $%s", xrpPriceUSD.StringFixed(6))

	// Step 1: Get token data with AMM liquidity and calculate live prices
	log.Printf("🔄 [SCREENER UPDATE] Step 1: Fetching token data with AMM liquidity for live price calculation...")

	// First, get the basic token data with liquidity
	tokensWithLiquidity, err := job.getTokensWithLiquidity()
	if err != nil {
		return fmt.Errorf("failed to get tokens with liquidity: %w", err)
	}

	log.Printf("💰 [SCREENER UPDATE] Found %d tokens with liquidity >$1000", len(tokensWithLiquidity))

	// Step 2: Calculate live prices and upsert into screener table
	log.Printf("📊 [SCREENER UPDATE] Step 2: Calculating live prices and upserting into screener table...")

	tokensInserted := 0
	for i, token := range tokensWithLiquidity {
		if i%10 == 0 {
			log.Printf("Progress: %d/%d tokens processed", i, len(tokensWithLiquidity))
		}

		if err := job.insertTokenWithLivePrice(token, xrpPriceUSD); err != nil {
			log.Printf("⚠️ [SCREENER UPDATE] Failed to insert token %s/%s: %v", token.Currency, token.Issuer, err)
			continue
		}
		tokensInserted++
	}

	log.Printf("✅ [SCREENER UPDATE] Upserted %d tokens with live prices", tokensInserted)

	// Step 3: Clean up obsolete records (tokens that no longer meet criteria)
	// This happens ONLY after successful updates to avoid data loss
	if tokensInserted > 0 {
		log.Printf("🧹 [SCREENER UPDATE] Step 3: Cleaning up obsolete records...")
		cleanupQuery := `
			DELETE FROM xrpScreener 
			WHERE updated_at < DATE_SUB(NOW(), INTERVAL 1 HOUR)
			AND liquidity_total < 1000`
		cleanupResult := job.db.Exec(cleanupQuery)
		if cleanupResult.Error == nil {
			log.Printf("🗑️ [SCREENER UPDATE] Removed %d obsolete records", cleanupResult.RowsAffected)
		}
	}

	// Step 4: Validation - verify table has reasonable data
	var recordCount int64
	countResult := job.db.Table("xrpScreener").Count(&recordCount)
	if countResult.Error != nil {
		return fmt.Errorf("failed to validate screener table: %w", countResult.Error)
	}

	if recordCount == 0 {
		return fmt.Errorf("screener table has no records - check data source and liquidity criteria")
	}

	// Step 5: Performance metrics
	originalTableSize := 9938 // Known full xrpTokens table size
	performanceImprovement := float64(originalTableSize) / float64(recordCount)

	log.Printf("📊 [SCREENER UPDATE] Performance Metrics:")
	log.Printf("📊   Original table size: %d records", originalTableSize)
	log.Printf("📊   Optimized table size: %d records", recordCount)
	log.Printf("📊   Data reduction: %.1f%% (%.1fx improvement)", (1.0-float64(recordCount)/float64(originalTableSize))*100, performanceImprovement)
	log.Printf("✅ [SCREENER UPDATE] Successfully updated screener table")

	return nil
}

// ForceUpdate triggers an immediate update outside the regular schedule
func (job *ScreenerUpdateJob) ForceUpdate() error {
	job.mu.RLock()
	if !job.running {
		job.mu.RUnlock()
		return fmt.Errorf("screener update job is not running")
	}
	job.mu.RUnlock()

	log.Printf("🚀 [SCREENER UPDATE JOB] Force update requested")
	job.performUpdate()
	return nil
}

// TokenWithLiquidity represents a token with its AMM liquidity data
type TokenWithLiquidity struct {
	Currency     string
	Issuer       string
	TokenName    string
	Icon         string
	SupplyXrpl   decimal.Decimal
	LiquidityUsd decimal.Decimal
	Trustlines   int
	Volume24h    decimal.Decimal
}

// getTokensWithLiquidity fetches tokens that have meaningful AMM liquidity (>$1000)
func (job *ScreenerUpdateJob) getTokensWithLiquidity() ([]TokenWithLiquidity, error) {
	var tokens []TokenWithLiquidity

	query := `
		SELECT 
			t.currency,
			t.issuer,
			COALESCE(t.token_name, t.currency) as token_name,
			COALESCE(t.icon, '') as icon,
			COALESCE(t.supply_xrpl, 0) as supply_xrpl,
			SUM(COALESCE(amm.liquidity_usd, 0)) as liquidity_total,
			COALESCE(t.trustlines, 0) as trustlines,
			COALESCE(t.volume_24h, 0) as volume_24h
		FROM xrpTokens t
		LEFT JOIN xrpAmm_normalized amm ON (
			(amm.asset1_currency COLLATE utf8mb4_unicode_ci = t.currency COLLATE utf8mb4_unicode_ci 
			 AND amm.asset1_issuer COLLATE utf8mb4_unicode_ci = t.issuer COLLATE utf8mb4_unicode_ci)
			OR (amm.asset2_currency COLLATE utf8mb4_unicode_ci = t.currency COLLATE utf8mb4_unicode_ci 
			    AND amm.asset2_issuer COLLATE utf8mb4_unicode_ci = t.issuer COLLATE utf8mb4_unicode_ci)
		)
		WHERE t.currency != '' AND t.issuer != ''
		GROUP BY t.currency, t.issuer, t.token_name, t.icon, t.supply_xrpl, t.trustlines, t.volume_24h
		HAVING SUM(COALESCE(amm.liquidity_usd, 0)) > 1000.00
		ORDER BY liquidity_total DESC
		LIMIT 250`

	err := job.db.Raw(query).Scan(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens with liquidity: %w", err)
	}

	return tokens, nil
}

// getTokenPriceInXRPFromDatabase calculates token price in XRP using database AMM pool data
func (job *ScreenerUpdateJob) getTokenPriceInXRPFromDatabase(currency, issuer string) (decimal.Decimal, error) {
	// Query AMM pool data from database where token is paired with XRP
	var xrpAmount, tokenAmount decimal.Decimal

	row := job.db.Raw(`
		SELECT 
			CASE 
				WHEN asset1_currency = 'XRP' THEN asset1_amount
				WHEN asset2_currency = 'XRP' THEN asset2_amount
			END as xrp_amount,
			CASE 
				WHEN asset1_currency = ? AND asset1_issuer = ? THEN asset1_amount
				WHEN asset2_currency = ? AND asset2_issuer = ? THEN asset2_amount
			END as token_amount
		FROM xrpAmm_normalized
		WHERE (
			(asset1_currency = 'XRP' AND asset2_currency = ? AND asset2_issuer = ?) OR
			(asset2_currency = 'XRP' AND asset1_currency = ? AND asset1_issuer = ?)
		)
		LIMIT 1
	`, currency, issuer, currency, issuer, currency, issuer, currency, issuer).Row()

	err := row.Scan(&xrpAmount, &tokenAmount)
	if err != nil {
		return decimal.Zero, fmt.Errorf("no XRP/%s pool found: %w", currency, err)
	}

	// Price = XRP amount / Token amount (how much XRP per token)
	if tokenAmount.IsZero() {
		return decimal.Zero, fmt.Errorf("token amount is zero")
	}

	return xrpAmount.Div(tokenAmount), nil
}

// insertTokenWithLivePrice calculates live AMM price and inserts token into screener table
func (job *ScreenerUpdateJob) insertTokenWithLivePrice(token TokenWithLiquidity, xrpPriceUSD decimal.Decimal) error {
	// Get token price in XRP from DATABASE to avoid XRPL API calls
	tokenPriceInXRP, err := job.getTokenPriceInXRPFromDatabase(token.Currency, token.Issuer)
	if err != nil {
		log.Printf("⚠️ [SCREENER PRICE] Failed to get token price in XRP for %s/%s, using 0: %v",
			token.Currency, token.Issuer, err)
		tokenPriceInXRP = decimal.Zero
	}

	// Calculate USD price using the cached XRP price
	livePrice := tokenPriceInXRP.Mul(xrpPriceUSD)

	// Calculate market cap: price * supply
	var marketCap decimal.Decimal
	if token.SupplyXrpl.GreaterThan(decimal.Zero) && livePrice.GreaterThan(decimal.Zero) {
		marketCap = livePrice.Mul(token.SupplyXrpl)
	} else {
		marketCap = decimal.Zero
	}

	// Convert to float64 for database storage
	priceFloat, _ := livePrice.Float64()
	marketCapFloat, _ := marketCap.Float64()
	supplyFloat, _ := token.SupplyXrpl.Float64()
	liquidityFloat, _ := token.LiquidityUsd.Float64()
	volumeFloat, _ := token.Volume24h.Float64()

	// UPSERT into screener table with live price data
	// This preserves existing data if update fails partway through
	upsertQuery := `
		INSERT INTO xrpScreener (
			currency, issuer_address, issuer_name, token_name, icon, price, marketcap, 
			supply_xrpl, liquidity_total, trustlines, volume_24h, created_at, updated_at
		) VALUES (?, ?, '', ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			token_name = VALUES(token_name),
			issuer_name = VALUES(issuer_name),
			icon = VALUES(icon),
			price = VALUES(price),
			marketcap = VALUES(marketcap),
			supply_xrpl = VALUES(supply_xrpl),
			liquidity_total = VALUES(liquidity_total),
			trustlines = VALUES(trustlines),
			volume_24h = VALUES(volume_24h),
			updated_at = NOW()`

	result := job.db.Exec(upsertQuery,
		token.Currency,
		token.Issuer,
		token.TokenName,
		token.Icon,
		priceFloat,
		marketCapFloat,
		supplyFloat,
		liquidityFloat,
		token.Trustlines,
		volumeFloat,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to insert token into screener table: %w", result.Error)
	}

	// Log removed to reduce log noise
	// if priceFloat > 0 {
	// 	log.Printf("💰 [SCREENER PRICE] %s/%s: $%.6f (Market Cap: $%.2f)",
	// 		token.Currency, token.Issuer, priceFloat, marketCapFloat)
	// }

	return nil
}

// getXRPLPriceWithRetry attempts to get XRP price with exponential backoff
func (job *ScreenerUpdateJob) getXRPLPriceWithRetry() (decimal.Decimal, error) {
	maxRetries := 2               // Reduced from 3 to minimize XRPL load
	baseDelay := time.Second * 30 // Increased from 2s to 30s for rate limit protection

	for attempt := 0; attempt < maxRetries; attempt++ {
		price, err := job.priceService.GetXRPLPrice()
		if err == nil {
			return price, nil
		}

		// Check if it's a rate limit error - if so, don't retry
		if strings.Contains(err.Error(), "policy violation") ||
			strings.Contains(err.Error(), "IP limit") ||
			strings.Contains(err.Error(), "close 1008") {
			log.Printf("🛑 [SCREENER UPDATE] XRPL rate limit detected, using database fallback")
			return decimal.Zero, err
		}

		log.Printf("⚠️ [SCREENER UPDATE] XRP price attempt %d/%d failed: %v",
			attempt+1, maxRetries, err)

		// Don't wait after the last attempt
		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			log.Printf("⏳ [SCREENER UPDATE] Waiting %v before retry...", delay)
			time.Sleep(delay)
		}
	}

	return decimal.Zero, fmt.Errorf("failed to get XRP price after %d attempts", maxRetries)
}
