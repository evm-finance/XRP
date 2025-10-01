package xrp

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"qc-defi-graphql-server/internal/models"

	"github.com/shopspring/decimal"
)

// AMMDiscoveryV2Stats tracks statistics for the enhanced discovery process
type AMMDiscoveryV2Stats struct {
	TotalTokensProcessed   int       `json:"total_tokens_processed"`
	TotalPoolsDiscovered   int       `json:"total_pools_discovered"`
	TotalPoolsStored       int       `json:"total_pools_stored"`
	TotalPoolsFiltered     int       `json:"total_pools_filtered"`
	TotalDuplicatesSkipped int       `json:"total_duplicates_skipped"`
	TotalErrors            int       `json:"total_errors"`
	ProcessingTimeMs       int64     `json:"processing_time_ms"`
	StartTime              time.Time `json:"start_time"`
	EndTime                time.Time `json:"end_time"`
	BatchSize              int       `json:"batch_size"`
	MinLiquidityUSD        float64   `json:"min_liquidity_usd"`
	DiscoveredPools        []string  `json:"discovered_pools"`
	ErrorMessages          []string  `json:"error_messages"`
}

// DiscoverAMMPoolsV2 discovers AMM pools using enhanced logic and stores them in xrpAmm_v2
func (s *AMMService) DiscoverAMMPoolsV2(batchSize int, minLiquidityUSD float64, maxTokens int) (*AMMDiscoveryV2Stats, error) {
	s.debugTraceMethodCall("DiscoverAMMPoolsV2", batchSize, minLiquidityUSD, maxTokens)

	startTime := time.Now()
	stats := &AMMDiscoveryV2Stats{
		StartTime:       startTime,
		BatchSize:       batchSize,
		MinLiquidityUSD: minLiquidityUSD,
		DiscoveredPools: []string{},
		ErrorMessages:   []string{},
	}

	// Default values
	if batchSize <= 0 {
		batchSize = 100
		stats.BatchSize = batchSize
	}
	if minLiquidityUSD <= 0 {
		minLiquidityUSD = 100.0
		stats.MinLiquidityUSD = minLiquidityUSD
	}
	if maxTokens <= 0 {
		maxTokens = 10000 // Reasonable default limit
	}

	log.Printf("🚀 Starting Enhanced AMM Pool Discovery V2...")
	log.Printf("   Batch size: %d", batchSize)
	log.Printf("   Min liquidity: $%.2f USD", minLiquidityUSD)
	log.Printf("   Max tokens to process: %d", maxTokens)

	// Process tokens in batches to discover pools
	offset := 0
	for {
		log.Printf("🔄 Processing token batch %d (offset: %d)...", (offset/batchSize)+1, offset)

		// Get batch of tokens from database
		tokens, err := s.getTokenBatch(batchSize, offset)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to fetch token batch at offset %d: %v", offset, err)
			stats.ErrorMessages = append(stats.ErrorMessages, errorMsg)
			stats.TotalErrors++
			log.Printf("❌ %s", errorMsg)
			break
		}

		if len(tokens) == 0 {
			log.Printf("✅ No more tokens to process")
			break
		}

		stats.TotalTokensProcessed += len(tokens)
		log.Printf("📊 Processing %d tokens in this batch...", len(tokens))

		// Check if we've reached the token limit
		if stats.TotalTokensProcessed >= maxTokens {
			log.Printf("⚠️ Reached maximum token limit (%d), stopping discovery", maxTokens)
			break
		}

		// Discover pools for this batch of tokens
		batchStats, err := s.discoverPoolsForTokenBatch(tokens, minLiquidityUSD)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to process token batch: %v", err)
			stats.ErrorMessages = append(stats.ErrorMessages, errorMsg)
			stats.TotalErrors++
			log.Printf("❌ %s", errorMsg)
		} else {
			// Aggregate batch statistics
			stats.TotalPoolsDiscovered += batchStats.PoolsDiscovered
			stats.TotalPoolsStored += batchStats.PoolsStored
			stats.TotalPoolsFiltered += batchStats.PoolsFiltered
			stats.TotalDuplicatesSkipped += batchStats.DuplicatesSkipped
			stats.TotalErrors += batchStats.Errors
			stats.DiscoveredPools = append(stats.DiscoveredPools, batchStats.DiscoveredPools...)
			stats.ErrorMessages = append(stats.ErrorMessages, batchStats.ErrorMessages...)
		}

		offset += batchSize

		// Log progress every few batches
		if (offset/batchSize)%5 == 0 {
			log.Printf("📈 Progress Update:")
			log.Printf("   Tokens processed: %d", stats.TotalTokensProcessed)
			log.Printf("   Pools discovered: %d", stats.TotalPoolsDiscovered)
			log.Printf("   Pools stored: %d", stats.TotalPoolsStored)
			log.Printf("   Duplicates skipped: %d", stats.TotalDuplicatesSkipped)
		}
	}

	// Finalize statistics
	stats.EndTime = time.Now()
	stats.ProcessingTimeMs = stats.EndTime.Sub(startTime).Milliseconds()

	log.Printf("🎉 Enhanced AMM Pool Discovery V2 Complete!")
	log.Printf("📊 Final Statistics:")
	log.Printf("   Tokens processed: %d", stats.TotalTokensProcessed)
	log.Printf("   Pools discovered: %d", stats.TotalPoolsDiscovered)
	log.Printf("   Pools stored: %d", stats.TotalPoolsStored)
	log.Printf("   Pools filtered out: %d", stats.TotalPoolsFiltered)
	log.Printf("   Duplicates skipped: %d", stats.TotalDuplicatesSkipped)
	log.Printf("   Errors encountered: %d", stats.TotalErrors)
	log.Printf("   Processing time: %dms", stats.ProcessingTimeMs)

	return stats, nil
}

// TokenBatchStats tracks statistics for a single batch of token processing
type TokenBatchStats struct {
	PoolsDiscovered   int      `json:"pools_discovered"`
	PoolsStored       int      `json:"pools_stored"`
	PoolsFiltered     int      `json:"pools_filtered"`
	DuplicatesSkipped int      `json:"duplicates_skipped"`
	Errors            int      `json:"errors"`
	DiscoveredPools   []string `json:"discovered_pools"`
	ErrorMessages     []string `json:"error_messages"`
}

// getTokenBatch retrieves a batch of tokens from the database
func (s *AMMService) getTokenBatch(batchSize, offset int) ([]TokenRecord, error) {
	// Get tokens from database with proper currency/issuer handling
	rows, err := s.db.Raw(`
		SELECT currency, issuer 
		FROM xrpTokens 
		WHERE currency IS NOT NULL 
			AND currency != '' 
			AND currency != 'NULL'
			AND issuer IS NOT NULL 
			AND issuer != ''
		ORDER BY currency, issuer
		LIMIT ? OFFSET ?
	`, batchSize, offset).Rows()

	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var tokens []TokenRecord
	uniqueTokens := make(map[string]bool)

	for rows.Next() {
		var currencyInterface interface{}
		var issuer string

		// Scan with proper type handling
		if err := rows.Scan(&currencyInterface, &issuer); err != nil {
			continue
		}

		// Convert currency to string handling byte arrays
		var currency string
		switch v := currencyInterface.(type) {
		case string:
			currency = v
		case []byte:
			currency = string(v)
		default:
			currency = fmt.Sprintf("%v", v)
		}

		// Skip invalid currencies
		if currency == "" || currency == "NULL" || issuer == "" {
			continue
		}

		// Deduplicate tokens
		uniqueKey := currency + "|" + issuer
		if uniqueTokens[uniqueKey] {
			continue
		}
		uniqueTokens[uniqueKey] = true

		tokens = append(tokens, TokenRecord{
			Currency: currency,
			Issuer:   issuer,
		})
	}

	return tokens, nil
}

// TokenRecord represents a token from the database
type TokenRecord struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
}

// discoverPoolsForTokenBatch discovers AMM pools for a batch of tokens
func (s *AMMService) discoverPoolsForTokenBatch(tokens []TokenRecord, minLiquidityUSD float64) (*TokenBatchStats, error) {
	stats := &TokenBatchStats{
		DiscoveredPools: []string{},
		ErrorMessages:   []string{},
	}

	// Use concurrent processing for better performance
	var wg sync.WaitGroup
	var mutex sync.Mutex
	semaphore := make(chan struct{}, 10) // Limit concurrent connections

	for _, token := range tokens {
		wg.Add(1)
		go func(token TokenRecord) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			// Discover pools for this token
			poolStats := s.discoverPoolsForSingleToken(token, minLiquidityUSD)

			// Aggregate statistics (thread-safe)
			mutex.Lock()
			stats.PoolsDiscovered += poolStats.PoolsDiscovered
			stats.PoolsStored += poolStats.PoolsStored
			stats.PoolsFiltered += poolStats.PoolsFiltered
			stats.DuplicatesSkipped += poolStats.DuplicatesSkipped
			stats.Errors += poolStats.Errors
			stats.DiscoveredPools = append(stats.DiscoveredPools, poolStats.DiscoveredPools...)
			stats.ErrorMessages = append(stats.ErrorMessages, poolStats.ErrorMessages...)
			mutex.Unlock()
		}(token)
	}

	wg.Wait() // Wait for all goroutines to complete
	return stats, nil
}

// discoverPoolsForSingleToken discovers AMM pools for a single token
func (s *AMMService) discoverPoolsForSingleToken(token TokenRecord, minLiquidityUSD float64) *TokenBatchStats {
	stats := &TokenBatchStats{
		DiscoveredPools: []string{},
		ErrorMessages:   []string{},
	}

	// Check pools in both directions, similar to original logic
	poolChecks := []struct {
		asset1Currency string
		asset1Issuer   string
		asset2Currency string
		asset2Issuer   string
		description    string
	}{
		// Token/XRP pools
		{token.Currency, token.Issuer, "XRP", "", "Token/XRP"},
		{"XRP", "", token.Currency, token.Issuer, "XRP/Token"},

		// Token/RLUSD pools
		{token.Currency, token.Issuer, RLUSDCurrency, RLUSDIssuer, "Token/RLUSD"},
		{RLUSDCurrency, RLUSDIssuer, token.Currency, token.Issuer, "RLUSD/Token"},
	}

	for _, check := range poolChecks {
		pool, err := s.CheckSingleAMMPool(check.asset1Currency, check.asset1Issuer, check.asset2Currency, check.asset2Issuer)
		if err != nil {
			// Log error but continue processing
			errorMsg := fmt.Sprintf("Error checking %s pool for %s/%s: %v", check.description, token.Currency, token.Issuer, err)
			stats.ErrorMessages = append(stats.ErrorMessages, errorMsg)
			stats.Errors++
			continue
		}

		if pool != nil {
			stats.PoolsDiscovered++

			// Use enhanced storage logic to store the pool
			stored, filtered, duplicate := s.storePoolWithEnhancedLogic(*pool, minLiquidityUSD)

			if stored {
				stats.PoolsStored++
				stats.DiscoveredPools = append(stats.DiscoveredPools, pool.Result.Amm.Account)
				log.Printf("✅ Stored %s pool: %s", check.description, pool.Result.Amm.Account)
			} else if filtered {
				stats.PoolsFiltered++
				log.Printf("⏭️ Filtered %s pool: %s (below liquidity threshold)", check.description, pool.Result.Amm.Account)
			} else if duplicate {
				stats.DuplicatesSkipped++
				log.Printf("⏭️ Skipped duplicate %s pool: %s", check.description, pool.Result.Amm.Account)
			}
		}
	}

	return stats
}

// storePoolWithEnhancedLogic stores a pool using the enhanced validation and storage logic
func (s *AMMService) storePoolWithEnhancedLogic(pool models.AMMInfo, minLiquidityUSD float64) (stored, filtered, duplicate bool) {
	// Parse amounts and handle drops conversion
	amount1Drops, err := strconv.ParseFloat(pool.GetAmountValue(), 64)
	if err != nil {
		log.Printf("❌ Invalid amount in pool %s: %v", pool.Result.Amm.Account, err)
		return false, false, false
	}

	// Convert XRP drops to XRP (GetAmountValue returns drops)
	amount1 := amount1Drops / 1000000.0

	amount2, err := strconv.ParseFloat(pool.GetAmount2Value(), 64)
	if err != nil {
		log.Printf("❌ Invalid amount2 in pool %s: %v", pool.Result.Amm.Account, err)
		return false, false, false
	}

	// Determine asset configuration using same logic as original
	var baseAmount, tokenAmount float64
	var tokenCurrency, tokenIssuer string
	var isXRPBase bool

	if pool.GetAmount2Currency() == "XRP" {
		tokenAmount = amount1
		baseAmount = amount2 / 1000000.0 // Convert drops to XRP from GetAmountValue()
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = true
	} else if pool.GetAmount2Currency() == RLUSDCurrency && pool.GetAmount2Issuer() == RLUSDIssuer {
		tokenAmount = amount1
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = false
	} else {
		baseAmount = amount1 / 1000000.0 // Convert drops to XRP from GetAmountValue()
		tokenAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = true
	}

	// Calculate liquidity like the original method
	var liquidityUSD decimal.Decimal
	if isXRPBase {
		xrpPriceUSD, err := s.GetLiveXRPPrice()
		if err != nil {
			fmt.Printf("❌ Cannot calculate liquidity for XRP pool - live price unavailable: %v\n", err)
			return false, false, true // error - skip this pool
		}
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2))
	} else {
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(decimal.NewFromInt(2)) // RLUSD ≈ $1
	}

	// Check liquidity threshold
	if liquidityUSD.LessThan(decimal.NewFromFloat(minLiquidityUSD)) {
		return false, true, false // filtered
	}

	// Check for duplicates (enhanced duplicate checking)
	exists, err := s.CheckPoolExistsInV2Table(pool.Result.Amm.Account)
	if err != nil {
		log.Printf("⚠️ Warning: failed to check pool existence: %v", err)
	} else if exists {
		return false, false, true // duplicate
	}

	// Store using optimized method
	err = s.storeAMMPoolOptimized(pool, baseAmount, tokenAmount, tokenCurrency, tokenIssuer, liquidityUSD, isXRPBase)
	if err != nil {
		log.Printf("❌ Failed to store pool %s: %v", pool.Result.Amm.Account, err)
		return false, false, false
	}

	return true, false, false // stored successfully
}

// CheckPoolExistsInV2Table checks if a pool already exists in the xrpAmm_v2 table
func (s *AMMService) CheckPoolExistsInV2Table(account string) (bool, error) {
	var count int64
	err := s.db.Raw("SELECT COUNT(*) FROM xrpAmm_v2 WHERE account = ?", account).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DiscoverAndStoreSpecificTokenPools discovers pools for specific tokens (utility method)
func (s *AMMService) DiscoverAndStoreSpecificTokenPools(currencies []string, issuers []string, minLiquidityUSD float64) (*AMMDiscoveryV2Stats, error) {
	s.debugTraceMethodCall("DiscoverAndStoreSpecificTokenPools", len(currencies), minLiquidityUSD)

	if len(currencies) != len(issuers) {
		return nil, fmt.Errorf("currencies and issuers arrays must have the same length")
	}

	startTime := time.Now()
	stats := &AMMDiscoveryV2Stats{
		StartTime:       startTime,
		MinLiquidityUSD: minLiquidityUSD,
		DiscoveredPools: []string{},
		ErrorMessages:   []string{},
	}

	// Create token records from provided data
	var tokens []TokenRecord
	for i, currency := range currencies {
		tokens = append(tokens, TokenRecord{
			Currency: currency,
			Issuer:   issuers[i],
		})
	}

	stats.TotalTokensProcessed = len(tokens)
	log.Printf("🔍 Discovering pools for %d specific tokens...", len(tokens))

	// Process the specific tokens
	batchStats, err := s.discoverPoolsForTokenBatch(tokens, minLiquidityUSD)
	if err != nil {
		return nil, fmt.Errorf("failed to process specific tokens: %w", err)
	}

	// Aggregate statistics
	stats.TotalPoolsDiscovered = batchStats.PoolsDiscovered
	stats.TotalPoolsStored = batchStats.PoolsStored
	stats.TotalPoolsFiltered = batchStats.PoolsFiltered
	stats.TotalDuplicatesSkipped = batchStats.DuplicatesSkipped
	stats.TotalErrors = batchStats.Errors
	stats.DiscoveredPools = batchStats.DiscoveredPools
	stats.ErrorMessages = batchStats.ErrorMessages

	// Finalize statistics
	stats.EndTime = time.Now()
	stats.ProcessingTimeMs = stats.EndTime.Sub(startTime).Milliseconds()

	log.Printf("✅ Specific token pool discovery complete: %d pools stored", stats.TotalPoolsStored)
	return stats, nil
}
