package xrp

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SupplyPopulationJob handles background population of token supply data
// This implements Phase 2 Step 3 of the XRP Screener Enhancement Spec
type SupplyPopulationJob struct {
	xrpService *XRPService
	db         *gorm.DB
	batchSize  int
	maxRetries int
}

// NewSupplyPopulationJob creates a new supply population job
func NewSupplyPopulationJob(xrpService *XRPService, batchSize int) *SupplyPopulationJob {
	return &SupplyPopulationJob{
		xrpService: xrpService,
		db:         xrpService.db,
		batchSize:  batchSize,
		maxRetries: 3,
	}
}

// PopulateTokenSupplies processes tokens in batches to populate supply data
func (job *SupplyPopulationJob) PopulateTokenSupplies() error {
	fmt.Println("🚀 Starting Token Supply Population Job")
	fmt.Println("======================================")
	fmt.Printf("📋 Batch size: %d tokens\n", job.batchSize)
	fmt.Printf("🔄 Max retries: %d per token\n", job.maxRetries)
	fmt.Println()

	// Get tokens that need supply data (either NULL or very old)
	tokens, err := job.getTokensNeedingSupplyData()
	if err != nil {
		return fmt.Errorf("failed to get tokens needing supply data: %w", err)
	}

	if len(tokens) == 0 {
		fmt.Println("✅ No tokens need supply data updates")
		return nil
	}

	fmt.Printf("📊 Found %d tokens needing supply data updates\n", len(tokens))

	// Process tokens in batches
	totalProcessed := 0
	totalUpdated := 0
	totalErrors := 0

	for i := 0; i < len(tokens); i += job.batchSize {
		end := i + job.batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := tokens[i:end]
		fmt.Printf("\n🔄 Processing batch %d-%d (%d tokens)...\n", i+1, end, len(batch))

		batchUpdated, batchErrors := job.processBatch(batch)
		totalProcessed += len(batch)
		totalUpdated += batchUpdated
		totalErrors += batchErrors

		// Brief pause between batches to avoid overwhelming XRPL
		if i+job.batchSize < len(tokens) {
			fmt.Printf("⏸️  Pausing 2 seconds between batches...\n")
			time.Sleep(2 * time.Second)
		}
	}

	// Final summary
	fmt.Printf("\n🎯 SUPPLY POPULATION JOB COMPLETED\n")
	fmt.Printf("==================================\n")
	fmt.Printf("📊 Tokens Processed: %d\n", totalProcessed)
	fmt.Printf("✅ Successfully Updated: %d\n", totalUpdated)
	fmt.Printf("❌ Errors: %d\n", totalErrors)
	fmt.Printf("📈 Success Rate: %.1f%%\n", float64(totalUpdated)/float64(totalProcessed)*100)

	return nil
}

// TokenSupplyRecord represents a token that needs supply data
type TokenSupplyRecord struct {
	Currency string `gorm:"column:currency"`
	Issuer   string `gorm:"column:issuer"`
	Name     string `gorm:"column:name"`
}

// getTokensNeedingSupplyData fetches tokens that need supply data updates
func (job *SupplyPopulationJob) getTokensNeedingSupplyData() ([]TokenSupplyRecord, error) {
	var tokens []TokenSupplyRecord

	// Get tokens where supply_xrpl is NULL or data is older than 24 hours
	query := `
		SELECT currency, issuer, token_name as name
		FROM xrpTokens 
		WHERE currency != '' AND issuer != '' 
		AND (
			supply_xrpl IS NULL 
			OR supply_last_updated IS NULL
			OR supply_last_updated < DATE_SUB(NOW(), INTERVAL 24 HOUR)
		)
		ORDER BY marketcap DESC, trustlines DESC
		LIMIT 50`

	err := job.db.Raw(query).Scan(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens needing supply data: %w", err)
	}

	return tokens, nil
}

// processBatch handles a batch of tokens for supply population
func (job *SupplyPopulationJob) processBatch(tokens []TokenSupplyRecord) (updated int, errors int) {
	for i, token := range tokens {
		fmt.Printf("  [%d/%d] %s (%s)... ", i+1, len(tokens), token.Currency, token.Name)

		success := job.processTokenWithRetries(token)
		if success {
			fmt.Printf("✅ Updated\n")
			updated++
		} else {
			fmt.Printf("❌ Failed\n")
			errors++
		}
	}

	return updated, errors
}

// processTokenWithRetries processes a single token with retry logic
func (job *SupplyPopulationJob) processTokenWithRetries(token TokenSupplyRecord) bool {
	for attempt := 1; attempt <= job.maxRetries; attempt++ {
		if job.processToken(token) {
			return true
		}

		if attempt < job.maxRetries {
			// Exponential backoff: 1s, 2s, 4s...
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			time.Sleep(backoff)
		}
	}
	return false
}

// processToken fetches and updates supply data for a single token
func (job *SupplyPopulationJob) processToken(token TokenSupplyRecord) bool {
	// Fetch supply from XRPL
	supply, err := job.xrpService.GetTokenSupply(token.Currency, token.Issuer)
	if err != nil {
		log.Printf("Failed to get supply for %s/%s: %v", token.Currency, token.Issuer, err)
		return false
	}

	// Update database with supply data
	updateSQL := `
		UPDATE xrpTokens 
		SET supply_xrpl = ?, supply_last_updated = NOW()
		WHERE currency = ? AND issuer = ?`

	// Convert decimal to float64 for database storage
	supplyFloat, _ := supply.Float64()

	err = job.db.Exec(updateSQL, supplyFloat, token.Currency, token.Issuer).Error
	if err != nil {
		log.Printf("Failed to update supply for %s/%s: %v", token.Currency, token.Issuer, err)
		return false
	}

	return true
}

// RunSupplyPopulationJob is a convenience function to run the job with default settings
func RunSupplyPopulationJob(xrpService *XRPService) error {
	job := NewSupplyPopulationJob(xrpService, 5) // Process 5 tokens at a time
	return job.PopulateTokenSupplies()
}

// UpdateTokenSupplyAndMarketCap updates both supply and recalculates market cap
// This implements the price * supply = marketcap formula from the spec
func (job *SupplyPopulationJob) UpdateTokenSupplyAndMarketCap(currency, issuer string) error {
	// Fetch current token data
	var tokenData struct {
		Price decimal.Decimal `gorm:"column:price"`
	}

	err := job.db.Raw("SELECT price FROM xrpTokens WHERE currency = ? AND issuer = ?", currency, issuer).Scan(&tokenData).Error
	if err != nil {
		return fmt.Errorf("failed to get token price: %w", err)
	}

	// Fetch supply from XRPL using validated ledger
	supply, err := job.xrpService.GetTokenSupply(currency, issuer)
	if err != nil {
		return fmt.Errorf("failed to get token supply: %w", err)
	}

	// Calculate market cap: price * supply
	marketCap := tokenData.Price.Mul(supply)

	// Update both supply and market cap
	updateSQL := `
		UPDATE xrpTokens 
		SET supply_xrpl = ?, 
		    marketcap = ?, 
		    supply_last_updated = NOW()
		WHERE currency = ? AND issuer = ?`

	supplyFloat, _ := supply.Float64()
	marketCapFloat, _ := marketCap.Float64()

	err = job.db.Exec(updateSQL, supplyFloat, marketCapFloat, currency, issuer).Error
	if err != nil {
		return fmt.Errorf("failed to update token data: %w", err)
	}

	priceFloat, _ := tokenData.Price.Float64()
	fmt.Printf("✅ Updated %s: Supply=%.2f, Price=%.6f, MarketCap=%.2f\n",
		currency, supplyFloat, priceFloat, marketCapFloat)

	return nil
}
