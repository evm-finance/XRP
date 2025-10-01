package xrp

import (
	"fmt"
	"log"
	"qc-defi-graphql-server/internal/models"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ScreenerService handles optimized screener queries using XRPL metadata instead of database tables
// This service provides high-performance queries for tokens with meaningful liquidity (>$1000)
type ScreenerService struct {
	db *gorm.DB
}

// NewScreenerService creates a new screener service instance
func NewScreenerService(db *gorm.DB) *ScreenerService {
	return &ScreenerService{
		db: db,
	}
}

// ScreenerToken represents the optimized screener token structure
type ScreenerToken struct {
	ID            int             `gorm:"column:id"`
	Currency      string          `gorm:"column:currency"`
	IssuerAddress string          `gorm:"column:issuer_address"`
	IssuerName    string          `gorm:"column:issuer_name"`
	TokenName     string          `gorm:"column:token_name"`
	Icon          string          `gorm:"column:icon"`
	Price         decimal.Decimal `gorm:"column:price"`
	Marketcap     decimal.Decimal `gorm:"column:marketcap"`
	SupplyXrpl    decimal.Decimal `gorm:"column:supply_xrpl"`
	LiquidityUSD  decimal.Decimal `gorm:"column:liquidity_total"`
	Trustlines    int             `gorm:"column:trustlines"`
	Volume24H     decimal.Decimal `gorm:"column:volume_24h"`
	UpdatedAt     string          `gorm:"column:updated_at"`
}

// GetScreenerTokens retrieves tokens directly from xrpScreener table with all fields
// This matches the reference card specification and includes marketcap and volume_24h
func (s *ScreenerService) GetScreenerTokens() ([]*models.XRPTokenFields, error) {
	startTime := time.Now()
	log.Printf("🔍 [SCREENER SERVICE] === SCREENER SERVICE START ===")

	if s.db == nil {
		log.Printf("❌ [SCREENER SERVICE] CRITICAL: Database connection is nil")
		return nil, fmt.Errorf("database connection not available")
	}

	log.Printf("🔍 [SCREENER SERVICE] ✅ Database connection available, querying xrpScreener table directly")

	var screenerTokens []struct {
		Currency      string  `gorm:"column:currency"`
		IssuerAddress string  `gorm:"column:issuer_address"`
		IssuerName    string  `gorm:"column:issuer_name"`
		TokenName     string  `gorm:"column:token_name"`
		Icon          string  `gorm:"column:icon"`
		Price         float64 `gorm:"column:price"`
		Marketcap     float64 `gorm:"column:marketcap"`
		Supply        float64 `gorm:"column:supply_xrpl"`
		LiquidityUSD  float64 `gorm:"column:liquidity_total"`
		Volume24H     float64 `gorm:"column:volume_24h"`
		Trustlines    int     `gorm:"column:trustlines"`
	}

	// Query xrpScreener table directly with all fields
	query := `
		SELECT 
			currency,
			issuer_address,
			COALESCE(issuer_name, '') as issuer_name,
			COALESCE(token_name, currency) as token_name,
			COALESCE(icon, '') as icon,
			COALESCE(price, 0) as price,
			COALESCE(marketcap, 0) as marketcap,
			COALESCE(supply_xrpl, 0) as supply_xrpl,
			COALESCE(liquidity_total, 0) as liquidity_total,
			COALESCE(volume_24h, 0) as volume_24h,
			COALESCE(trustlines, 0) as trustlines
		FROM xrpScreener 
		WHERE currency != '' AND issuer_address != ''
		ORDER BY liquidity_total DESC, marketcap DESC, price DESC
		LIMIT 100`

	log.Printf("🔍 [SCREENER SERVICE] Executing SQL query: %s", query)

	// First check if table exists and has data
	var tableCount int64
	countErr := s.db.Raw("SELECT COUNT(*) FROM xrpScreener").Scan(&tableCount).Error
	if countErr != nil {
		log.Printf("❌ [SCREENER SERVICE] CRITICAL: Cannot count xrpScreener table: %v", countErr)
		return nil, fmt.Errorf("cannot access xrpScreener table: %w", countErr)
	}
	log.Printf("🔍 [SCREENER SERVICE] xrpScreener table has %d total records", tableCount)

	// Check filtered count
	var filteredCount int64
	filteredCountErr := s.db.Raw("SELECT COUNT(*) FROM xrpScreener WHERE currency != '' AND issuer_address != ''").Scan(&filteredCount).Error
	if filteredCountErr != nil {
		log.Printf("❌ [SCREENER SERVICE] Cannot count filtered xrpScreener records: %v", filteredCountErr)
	} else {
		log.Printf("🔍 [SCREENER SERVICE] xrpScreener has %d records matching WHERE conditions", filteredCount)
	}

	queryStart := time.Now()
	queryResult := s.db.Raw(query).Scan(&screenerTokens)
	queryDuration := time.Since(queryStart)

	log.Printf("🔍 [SCREENER SERVICE] Raw GORM result - Error: %v, RowsAffected: %d", queryResult.Error, queryResult.RowsAffected)

	if queryResult.Error != nil {
		log.Printf("❌ [SCREENER SERVICE] Database query failed after %v: %v", queryDuration, queryResult.Error)
		log.Printf("❌ [SCREENER SERVICE] Error type: %T", queryResult.Error)
		log.Printf("❌ [SCREENER SERVICE] Query was: %s", query)
		return nil, fmt.Errorf("failed to query xrpScreener table: %w", queryResult.Error)
	}

	log.Printf("✅ [SCREENER SERVICE] Database query succeeded in %v", queryDuration)
	log.Printf("📊 [SCREENER SERVICE] GORM returned %d rows, scanned into %d tokens", queryResult.RowsAffected, len(screenerTokens))

	if len(screenerTokens) == 0 {
		log.Printf("❌ [SCREENER SERVICE] CRITICAL: Query returned 0 tokens despite table having %d records", tableCount)
		log.Printf("❌ [SCREENER SERVICE] This indicates a problem with the WHERE conditions or data structure")
		log.Printf("❌ [SCREENER SERVICE] Table count: %d, Filtered count: %d, Query result: %d", tableCount, filteredCount, len(screenerTokens))

		// Let's see what the actual data looks like
		var sampleRows []map[string]interface{}
		sampleResult := s.db.Raw("SELECT currency, issuer_address, issuer_name, token_name, price, marketcap, volume_24h FROM xrpScreener LIMIT 3").Find(&sampleRows)
		if sampleResult.Error == nil {
			log.Printf("🔍 [SCREENER SERVICE] Sample data from xrpScreener table:")
			for i, row := range sampleRows {
				log.Printf("  Row %d: %v", i+1, row)
			}
		} else {
			log.Printf("❌ [SCREENER SERVICE] Cannot query sample data: %v", sampleResult.Error)
		}

		return nil, fmt.Errorf("xrpScreener query returned no results despite table having %d records - data structure issue", tableCount)
	}

	// Convert to expected format with detailed logging
	log.Printf("🔄 [SCREENER SERVICE] Converting %d raw tokens to XRPTokenFields format", len(screenerTokens))

	var result []*models.XRPTokenFields
	tokensWithMarketCap := 0
	tokensWithVolume := 0
	tokensWithPrices := 0

	for i, token := range screenerTokens {
		// Log first 3 tokens for debugging
		if i < 3 {
			log.Printf("🔍 [SCREENER SERVICE] Token %d: %s (%s) by %s - Price: $%.6f, MarketCap: $%.2f, Volume24h: $%.2f, Liquidity: $%.2f",
				i+1, token.TokenName, token.Currency, token.IssuerName, token.Price, token.Marketcap, token.Volume24H, token.LiquidityUSD)
		}

		// Calculate marketcap if missing using simple formula: price * supply
		marketcap := token.Marketcap
		if marketcap == 0 && token.Price > 0 && token.Supply > 0 {
			marketcap = token.Price * token.Supply
		}

		// Track data completeness
		if marketcap > 0 {
			tokensWithMarketCap++
		}
		if token.Volume24H > 0 {
			tokensWithVolume++
		}
		if token.Price > 0 {
			tokensWithPrices++
		}

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.IssuerAddress,
			TokenName:     token.TokenName,
			IssuerName:    token.IssuerName, // ✅ Now stored in screener table
			Icon:          token.Icon,
			Marketcap:     marketcap, // ✅ Calculated from price * supply if needed
			Price:         token.Price,
			Supply:        token.Supply,
			Liquidity:     token.LiquidityUSD,
			Volume24H:     token.Volume24H, // ✅ Now included
			Trustlines:    token.Trustlines,
		})
	}

	totalDuration := time.Since(startTime)

	log.Printf("📊 [SCREENER SERVICE] Data completeness: %d with prices | %d with marketcap | %d with volume24h",
		tokensWithPrices, tokensWithMarketCap, tokensWithVolume)
	log.Printf("✅ [SCREENER SERVICE] === SCREENER SERVICE COMPLETE === (%v total)", totalDuration)
	log.Printf("📊 [SCREENER SERVICE] Returning %d tokens with complete data including marketcap and volume_24h", len(result))

	return result, nil
}

// GetScreenerTokensDynamic retrieves tokens using XRPL metadata instead of database table (FALLBACK)
// Returns only tokens with liquidity > $1000 for better user experience
func (s *ScreenerService) GetScreenerTokensDynamic() ([]*models.XRPTokenFields, error) {
	if s.db == nil {
		log.Printf("❌ [SCREENER] Database connection is nil")
		return nil, fmt.Errorf("database connection not available")
	}

	log.Printf("🔍 [SCREENER] Starting query using XRPL metadata and AMM data")

	// Get tokens with liquidity from AMM pools
	var tokensWithLiquidity []struct {
		Currency       string  `gorm:"column:currency"`
		Issuer         string  `gorm:"column:issuer"`
		TotalLiquidity float64 `gorm:"column:total_liquidity"`
		PoolCount      int     `gorm:"column:pool_count"`
	}

	// Query AMM pools to get tokens with meaningful liquidity
	query := `
		WITH TopPools AS (
			SELECT * FROM xrpAmm_normalized 
			WHERE liquidity_usd > 1000 
			ORDER BY liquidity_usd DESC 
			LIMIT 500
		),
		TokensFromPools AS (
			SELECT 
				asset1_currency as currency,
				asset1_issuer as issuer,
				liquidity_usd
			FROM TopPools
			WHERE asset1_currency != '' 
				AND asset1_currency != 'XRP' 
				AND asset1_issuer != ''
			UNION ALL
			SELECT 
				asset2_currency as currency,
				asset2_issuer as issuer,
				liquidity_usd
			FROM TopPools
			WHERE asset2_currency != '' 
				AND asset2_currency != 'XRP' 
				AND asset2_issuer != ''
		)
		SELECT 
			currency,
			issuer,
			SUM(liquidity_usd) as total_liquidity,
			COUNT(*) as pool_count
		FROM TokensFromPools
		GROUP BY currency, issuer
		HAVING SUM(liquidity_usd) > 1000.00
		ORDER BY total_liquidity DESC
		LIMIT 50`

	err := s.db.Raw(query).Scan(&tokensWithLiquidity).Error
	if err != nil {
		log.Printf("❌ [SCREENER] AMM query failed: %v", err)
		return nil, fmt.Errorf("failed to query AMM tokens: %w", err)
	}

	log.Printf("✅ [SCREENER] Found %d tokens with liquidity > $1000", len(tokensWithLiquidity))

	if len(tokensWithLiquidity) == 0 {
		log.Printf("⚠️ [SCREENER] No tokens found with liquidity > $1000")
		return []*models.XRPTokenFields{}, nil
	}

	// Get basic token info from xrpTokens table for names and icons
	var result []*models.XRPTokenFields
	for _, token := range tokensWithLiquidity {
		// Get token metadata from xrpTokens table
		var tokenInfo struct {
			Name       string  `gorm:"column:name"`
			Icon       string  `gorm:"column:icon"`
			Price      float64 `gorm:"column:price"`
			Supply     float64 `gorm:"column:supply_xrpl"`
			Holders    *int    `gorm:"column:holders"`
			Marketcap  float64 `gorm:"column:marketcap"`
			Trustlines *int    `gorm:"column:trustlines"`
		}

		// Query token metadata with enhanced fields
		metadataQuery := `
			SELECT 
				COALESCE(name, ?) as name,
				COALESCE(icon, '') as icon,
				COALESCE(price, 0) as price,
				COALESCE(supply_xrpl, 0) as supply_xrpl,
				holders,
				COALESCE(marketcap, 0) as marketcap,
				trustlines
			FROM xrpTokens 
			WHERE currency = ? AND issuer = ?
			LIMIT 1`

		err := s.db.Raw(metadataQuery, token.Currency, token.Currency, token.Issuer).Scan(&tokenInfo).Error
		if err != nil {
			log.Printf("⚠️ [SCREENER] Failed to get metadata for %s/%s: %v", token.Currency, token.Issuer, err)
			// Continue with default values
		} else {
			// ENHANCED LOGGING: Log token metadata field existence and values
			log.Printf("📊 [SCREENER] Token metadata for %s/%s:", token.Currency, token.Issuer)
			log.Printf("📊 [SCREENER] - Name: %s", tokenInfo.Name)
			log.Printf("📊 [SCREENER] - Price: %.6f", tokenInfo.Price)
			log.Printf("📊 [SCREENER] - Supply: %.6f", tokenInfo.Supply)
			log.Printf("📊 [SCREENER] - Holders: %v (exists: %t)", tokenInfo.Holders, tokenInfo.Holders != nil)
			log.Printf("📊 [SCREENER] - Marketcap: %.6f", tokenInfo.Marketcap)
			log.Printf("📊 [SCREENER] - Trustlines: %v (exists: %t)", tokenInfo.Trustlines, tokenInfo.Trustlines != nil)
			log.Printf("📊 [SCREENER] - Icon: %s", tokenInfo.Icon)
		}

		// Calculate market cap if we have price and supply
		marketcap := 0.0
		if tokenInfo.Price > 0 && tokenInfo.Supply > 0 {
			marketcap = tokenInfo.Price * tokenInfo.Supply
		}

		// Use token name or fallback to currency
		tokenName := tokenInfo.Name
		if tokenName == "" {
			tokenName = token.Currency
		}

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.Issuer,
			TokenName:     tokenName,
			IssuerName:    "", // Deprecated field, keeping for backwards compatibility
			Icon:          tokenInfo.Icon,
			Marketcap:     marketcap,
			Price:         tokenInfo.Price,
			Supply:        tokenInfo.Supply,
			Liquidity:     token.TotalLiquidity,
		})
	}

	// Sort by liquidity (highest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Liquidity > result[j].Liquidity
	})

	log.Printf("✅ [SCREENER] Converted %d tokens to API format", len(result))

	// Log first few tokens for debugging
	for i, token := range result {
		if i < 3 {
			log.Printf("📋 [SCREENER] Sample Token %d: %s (%s), Issuer: %.16s..., Liquidity: %.2f",
				i+1, token.TokenName, token.Currency, token.IssuerAddress, token.Liquidity)
		}
	}

	return result, nil
}

// GetScreenerTokenBySymbol retrieves a specific token using XRPL metadata
func (s *ScreenerService) GetScreenerTokenBySymbol(currency, issuer string) (*models.XRPTokenFields, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	// Get liquidity for this specific token from AMM pools
	var liquidityData struct {
		TotalLiquidity float64 `gorm:"column:total_liquidity"`
		PoolCount      int     `gorm:"column:pool_count"`
	}

	liquidityQuery := `
		SELECT 
			SUM(liquidity_usd) as total_liquidity,
			COUNT(*) as pool_count
		FROM xrpAmm_normalized 
		WHERE (asset1_currency = ? AND asset1_issuer = ?) 
		   OR (asset2_currency = ? AND asset2_issuer = ?)
		   AND liquidity_usd > 0`

	err := s.db.Raw(liquidityQuery, currency, issuer, currency, issuer).Scan(&liquidityData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query liquidity for %s/%s: %w", currency, issuer, err)
	}

	// Get token metadata from xrpTokens table
	var tokenInfo struct {
		Name   string  `gorm:"column:name"`
		Icon   string  `gorm:"column:icon"`
		Price  float64 `gorm:"column:price"`
		Supply float64 `gorm:"column:supply_xrpl"`
	}

	metadataQuery := `
		SELECT 
			COALESCE(name, ?) as name,
			COALESCE(icon, '') as icon,
			COALESCE(price, 0) as price,
			COALESCE(supply_xrpl, 0) as supply_xrpl
		FROM xrpTokens 
		WHERE currency = ? AND issuer = ?
		LIMIT 1`

	err = s.db.Raw(metadataQuery, currency, currency, issuer).Scan(&tokenInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query token metadata for %s/%s: %w", currency, issuer, err)
	}

	// Calculate market cap if we have price and supply
	marketcap := 0.0
	if tokenInfo.Price > 0 && tokenInfo.Supply > 0 {
		marketcap = tokenInfo.Price * tokenInfo.Supply
	}

	// Use token name or fallback to currency
	tokenName := tokenInfo.Name
	if tokenName == "" {
		tokenName = currency
	}

	return &models.XRPTokenFields{
		Currency:      currency,
		IssuerAddress: issuer,
		TokenName:     tokenName,
		IssuerName:    "",
		Icon:          tokenInfo.Icon,
		Marketcap:     marketcap,
		Price:         tokenInfo.Price,
		Supply:        tokenInfo.Supply,
		Liquidity:     liquidityData.TotalLiquidity,
	}, nil
}

// GetScreenerStats provides statistics about the screener data
func (s *ScreenerService) GetScreenerStats() (map[string]interface{}, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	var stats struct {
		TotalTokens     int     `gorm:"column:total_tokens"`
		AvgLiquidity    float64 `gorm:"column:avg_liquidity"`
		MaxLiquidity    float64 `gorm:"column:max_liquidity"`
		TotalLiquidity  float64 `gorm:"column:total_liquidity"`
		RecentlyUpdated int     `gorm:"column:recently_updated"`
		AvgMarketCap    float64 `gorm:"column:avg_marketcap"`
	}

	// Query AMM-based statistics
	query := `
		WITH TokenLiquidity AS (
			SELECT 
				asset1_currency as currency,
				asset1_issuer as issuer,
				SUM(liquidity_usd) as total_liquidity
			FROM xrpAmm_normalized
			WHERE asset1_currency != '' AND asset1_currency != 'XRP' AND asset1_issuer != ''
			GROUP BY asset1_currency, asset1_issuer
			UNION ALL
			SELECT 
				asset2_currency as currency,
				asset2_issuer as issuer,
				SUM(liquidity_usd) as total_liquidity
			FROM xrpAmm_normalized
			WHERE asset2_currency != '' AND asset2_currency != 'XRP' AND asset2_issuer != ''
			GROUP BY asset2_currency, asset2_issuer
		)
		SELECT 
			COUNT(DISTINCT CONCAT(currency, issuer)) as total_tokens,
			AVG(total_liquidity) as avg_liquidity,
			MAX(total_liquidity) as max_liquidity,
			SUM(total_liquidity) as total_liquidity,
			COUNT(CASE WHEN total_liquidity > 1000.00 THEN 1 END) as recently_updated,
			0 as avg_marketcap
		FROM TokenLiquidity
		WHERE total_liquidity > 1000.00`

	err := s.db.Raw(query).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get screener stats: %w", err)
	}

	return map[string]interface{}{
		"total_tokens":     stats.TotalTokens,
		"avg_liquidity":    stats.AvgLiquidity,
		"max_liquidity":    stats.MaxLiquidity,
		"total_liquidity":  stats.TotalLiquidity,
		"recently_updated": stats.RecentlyUpdated,
		"avg_marketcap":    stats.AvgMarketCap,
		"performance_gain": "XRPL metadata + AMM data",
		"criteria":         "liquidity > $1000 USD from AMM pools",
		"data_source":      "XRPL metadata + xrpAmm_normalized",
	}, nil
}

// ValidateScreenerData checks if the required data sources are available
func (s *ScreenerService) ValidateScreenerData() error {
	if s.db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Check if AMM table exists and has data
	var ammCount int
	err := s.db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized WHERE liquidity_usd > 1000").Scan(&ammCount).Error
	if err != nil {
		return fmt.Errorf("failed to check AMM data: %w", err)
	}

	log.Printf("✅ [SCREENER] AMM validation: %d pools with liquidity > $1000", ammCount)

	if ammCount == 0 {
		log.Printf("⚠️ [SCREENER] No AMM pools with liquidity > $1000 found")
	}

	// Check if xrpTokens table exists
	var tokenCount int
	err = s.db.Raw("SELECT COUNT(*) FROM xrpTokens LIMIT 1").Scan(&tokenCount).Error
	if err != nil {
		return fmt.Errorf("failed to check token metadata: %w", err)
	}

	log.Printf("✅ [SCREENER] Token metadata validation: %d tokens available", tokenCount)

	return nil
}

// GetTopTokensByLiquidity returns the top tokens by liquidity for showcase/testing
func (s *ScreenerService) GetTopTokensByLiquidity(limit int) ([]*models.XRPTokenFields, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	if limit <= 0 || limit > 50 {
		limit = 10 // Default to top 10
	}

	// Get top tokens by liquidity from AMM data
	var topTokens []struct {
		Currency       string  `gorm:"column:currency"`
		Issuer         string  `gorm:"column:issuer"`
		TotalLiquidity float64 `gorm:"column:total_liquidity"`
	}

	query := `
		WITH TokenLiquidity AS (
			SELECT 
				asset1_currency as currency,
				asset1_issuer as issuer,
				SUM(liquidity_usd) as total_liquidity
			FROM xrpAmm_normalized
			WHERE asset1_currency != '' AND asset1_currency != 'XRP' AND asset1_issuer != ''
			GROUP BY asset1_currency, asset1_issuer
			UNION ALL
			SELECT 
				asset2_currency as currency,
				asset2_issuer as issuer,
				SUM(liquidity_usd) as total_liquidity
			FROM xrpAmm_normalized
			WHERE asset2_currency != '' AND asset2_currency != 'XRP' AND asset2_issuer != ''
			GROUP BY asset2_currency, asset2_issuer
		)
		SELECT 
			currency,
			issuer,
			SUM(total_liquidity) as total_liquidity
		FROM TokenLiquidity
		WHERE total_liquidity > 1000.00
		GROUP BY currency, issuer
		ORDER BY total_liquidity DESC
		LIMIT ?`

	err := s.db.Raw(query, limit).Scan(&topTokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query top tokens: %w", err)
	}

	// Convert to API format with metadata
	var result []*models.XRPTokenFields
	for _, token := range topTokens {
		// Get token metadata
		var tokenInfo struct {
			Name   string  `gorm:"column:name"`
			Icon   string  `gorm:"column:icon"`
			Price  float64 `gorm:"column:price"`
			Supply float64 `gorm:"column:supply_xrpl"`
		}

		metadataQuery := `
			SELECT 
				COALESCE(name, ?) as name,
				COALESCE(icon, '') as icon,
				COALESCE(price, 0) as price,
				COALESCE(supply_xrpl, 0) as supply_xrpl
			FROM xrpTokens 
			WHERE currency = ? AND issuer = ?
			LIMIT 1`

		err := s.db.Raw(metadataQuery, token.Currency, token.Currency, token.Issuer).Scan(&tokenInfo).Error
		if err != nil {
			log.Printf("⚠️ [SCREENER] Failed to get metadata for %s/%s: %v", token.Currency, token.Issuer, err)
		}

		// Calculate market cap
		marketcap := 0.0
		if tokenInfo.Price > 0 && tokenInfo.Supply > 0 {
			marketcap = tokenInfo.Price * tokenInfo.Supply
		}

		// Use token name or fallback to currency
		tokenName := tokenInfo.Name
		if tokenName == "" {
			tokenName = token.Currency
		}

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.Issuer,
			TokenName:     tokenName,
			IssuerName:    "",
			Icon:          tokenInfo.Icon,
			Marketcap:     marketcap,
			Price:         tokenInfo.Price,
			Supply:        tokenInfo.Supply,
			Liquidity:     token.TotalLiquidity,
		})
	}

	return result, nil
}
