package xrp

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// LiquidityCalculator handles liquidity calculation from existing AMM pool data
// This implements Phase 3 of the XRP Screener Enhancement Spec
type LiquidityCalculator struct {
	db *gorm.DB
}

// NewLiquidityCalculator creates a new liquidity calculator
func NewLiquidityCalculator(db *gorm.DB) *LiquidityCalculator {
	return &LiquidityCalculator{
		db: db,
	}
}

// CalculateTokenLiquidity calculates total liquidity for a specific token from AMM pools
func (calc *LiquidityCalculator) CalculateTokenLiquidity(currency, issuer string) (decimal.Decimal, error) {
	// Input validation
	if currency == "" {
		return decimal.Zero, fmt.Errorf("currency cannot be empty")
	}
	if issuer == "" {
		return decimal.Zero, fmt.Errorf("issuer cannot be empty")
	}

	fmt.Printf("💧 [LIQUIDITY] Calculating liquidity for %s issued by %s\n", currency, issuer)

	var pools []struct {
		Account      string  `gorm:"column:account"`
		LiquidityUSD float64 `gorm:"column:liquidity_usd"`
	}

	// Query all AMM pools containing this token (either as asset1 or asset2)
	query := `
		SELECT account, liquidity_usd
		FROM xrpAmm_normalized 
		WHERE (
			(asset1_currency = ? AND asset1_issuer = ?) 
			OR (asset2_currency = ? AND asset2_issuer = ?)
		)
		AND liquidity_usd IS NOT NULL 
		AND liquidity_usd > 0
		ORDER BY liquidity_usd DESC`

	err := calc.db.Raw(query, currency, issuer, currency, issuer).Scan(&pools).Error
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to query AMM pools: %w", err)
	}

	if len(pools) == 0 {
		fmt.Printf("💧 [LIQUIDITY] No AMM pools found for %s\n", currency)
		return decimal.Zero, nil
	}

	// Sum liquidity from all pools containing this token
	totalLiquidity := decimal.Zero
	for _, pool := range pools {
		poolLiquidity := decimal.NewFromFloat(pool.LiquidityUSD)
		totalLiquidity = totalLiquidity.Add(poolLiquidity)

		fmt.Printf("💧 [LIQUIDITY] Pool %s: $%.2f USD\n", pool.Account[:10]+"...", pool.LiquidityUSD)
	}

	fmt.Printf("✅ [LIQUIDITY] Total liquidity for %s: $%.2f USD (%d pools)\n",
		currency, totalLiquidity.InexactFloat64(), len(pools))

	return totalLiquidity, nil
}

// PopulateAllTokenLiquidity calculates and updates liquidity for all tokens
func (calc *LiquidityCalculator) PopulateAllTokenLiquidity() error {
	fmt.Println("🚀 Starting Token Liquidity Population")
	fmt.Println("======================================")
	fmt.Println("📋 Phase 3: Calculating liquidity from existing AMM pool data")
	fmt.Println()

	// Get all tokens that need liquidity data
	tokens, err := calc.getTokensForLiquidityCalculation()
	if err != nil {
		return fmt.Errorf("failed to get tokens for liquidity calculation: %w", err)
	}

	if len(tokens) == 0 {
		fmt.Println("⚠️ No tokens found needing liquidity calculation")
		return nil
	}

	fmt.Printf("📊 Found %d tokens for liquidity calculation\n", len(tokens))

	// Process tokens in batches
	batchSize := 10
	totalProcessed := 0
	totalUpdated := 0
	totalWithLiquidity := 0

	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := tokens[i:end]
		fmt.Printf("\n🔄 Processing batch %d-%d (%d tokens)...\n", i+1, end, len(batch))

		for j, token := range batch {
			fmt.Printf("  [%d/%d] %s (%s)... ", j+1, len(batch), token.Currency, token.Name)

			// Calculate liquidity for this token
			liquidity, err := calc.CalculateTokenLiquidity(token.Currency, token.Issuer)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			// Update database with calculated liquidity
			liquidityFloat, _ := liquidity.Float64()
			success := calc.updateTokenLiquidity(token.Currency, token.Issuer, liquidityFloat)

			if success {
				if liquidityFloat > 0 {
					fmt.Printf("✅ $%.2f\n", liquidityFloat)
					totalWithLiquidity++
				} else {
					fmt.Printf("✅ No liquidity\n")
				}
				totalUpdated++
			} else {
				fmt.Printf("❌ Update failed\n")
			}

			totalProcessed++
		}

		// Brief pause between batches
		if end < len(tokens) {
			fmt.Printf("⏸️  Pausing 1 second between batches...\n")
		}
	}

	// Summary
	fmt.Printf("\n🎯 LIQUIDITY POPULATION COMPLETED\n")
	fmt.Printf("=================================\n")
	fmt.Printf("📊 Tokens Processed: %d\n", totalProcessed)
	fmt.Printf("✅ Successfully Updated: %d\n", totalUpdated)
	fmt.Printf("💧 With Liquidity Data: %d\n", totalWithLiquidity)
	fmt.Printf("📈 Liquidity Success Rate: %.1f%%\n", float64(totalWithLiquidity)/float64(totalProcessed)*100)

	return nil
}

// TokenLiquidityData represents token data for liquidity calculation
type TokenLiquidityData struct {
	Currency string `gorm:"column:currency"`
	Issuer   string `gorm:"column:issuer"`
	Name     string `gorm:"column:name"`
}

// getTokensForLiquidityCalculation fetches tokens that need liquidity calculation
func (calc *LiquidityCalculator) getTokensForLiquidityCalculation() ([]TokenLiquidityData, error) {
	var tokens []TokenLiquidityData

	// Get tokens that exist in xrpTokens but don't have recent liquidity data
	query := `
		SELECT currency, issuer, name
		FROM xrpTokens 
		WHERE currency != '' AND issuer != ''
		AND (
			liquidity_total IS NULL 
			OR supply_last_updated IS NULL
			OR supply_last_updated < DATE_SUB(NOW(), INTERVAL 6 HOUR)
		)
		ORDER BY marketcap DESC, trustlines DESC
		LIMIT 100`

	err := calc.db.Raw(query).Scan(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens for liquidity calculation: %w", err)
	}

	return tokens, nil
}

// updateTokenLiquidity updates the liquidity_total field for a token
func (calc *LiquidityCalculator) updateTokenLiquidity(currency, issuer string, liquidity float64) bool {
	updateSQL := `
		UPDATE xrpTokens 
		SET liquidity_total = ?, supply_last_updated = NOW()
		WHERE currency = ? AND issuer = ?`

	err := calc.db.Exec(updateSQL, liquidity, currency, issuer).Error
	if err != nil {
		log.Printf("Failed to update liquidity for %s/%s: %v", currency, issuer, err)
		return false
	}

	return true
}

// GetLiquidityStatistics provides overview of liquidity data in the system
func (calc *LiquidityCalculator) GetLiquidityStatistics() error {
	fmt.Println("📊 Liquidity Statistics Overview")
	fmt.Println("================================")

	// Overall stats
	var overallStats struct {
		TotalTokens    int     `gorm:"column:total_tokens"`
		WithLiquidity  int     `gorm:"column:with_liquidity"`
		AvgLiquidity   float64 `gorm:"column:avg_liquidity"`
		TotalLiquidity float64 `gorm:"column:total_liquidity"`
		MaxLiquidity   float64 `gorm:"column:max_liquidity"`
	}

	err := calc.db.Raw(`
		SELECT 
			COUNT(*) as total_tokens,
			SUM(CASE WHEN liquidity_total IS NOT NULL AND liquidity_total > 0 THEN 1 ELSE 0 END) as with_liquidity,
			AVG(CASE WHEN liquidity_total > 0 THEN liquidity_total ELSE NULL END) as avg_liquidity,
			SUM(CASE WHEN liquidity_total > 0 THEN liquidity_total ELSE 0 END) as total_liquidity,
			MAX(liquidity_total) as max_liquidity
		FROM xrpTokens 
		WHERE currency != '' AND issuer != ''
	`).Scan(&overallStats).Error

	if err != nil {
		return fmt.Errorf("failed to get liquidity statistics: %v", err)
	}

	fmt.Printf("📈 Overall Liquidity Statistics:\n")
	fmt.Printf("  Total tokens: %d\n", overallStats.TotalTokens)
	fmt.Printf("  With liquidity data: %d\n", overallStats.WithLiquidity)
	fmt.Printf("  Average liquidity: $%.2f\n", overallStats.AvgLiquidity)
	fmt.Printf("  Total liquidity: $%.2f\n", overallStats.TotalLiquidity)
	fmt.Printf("  Max liquidity: $%.2f\n", overallStats.MaxLiquidity)

	// Top liquidity tokens
	fmt.Printf("\n🏆 Top 5 Tokens by Liquidity:\n")
	var topTokens []struct {
		Currency     string  `gorm:"column:currency"`
		Issuer       string  `gorm:"column:issuer"`
		Name         string  `gorm:"column:name"`
		LiquidityUSD float64 `gorm:"column:liquidity_total"`
	}

	err = calc.db.Raw(`
		SELECT currency, issuer, name, liquidity_total
		FROM xrpTokens 
		WHERE liquidity_total IS NOT NULL AND liquidity_total > 0
		ORDER BY liquidity_total DESC
		LIMIT 5
	`).Scan(&topTokens).Error

	if err != nil {
		return fmt.Errorf("failed to get top liquidity tokens: %v", err)
	}

	for i, token := range topTokens {
		name := token.Name
		if name == "" {
			name = "Unknown"
		}
		fmt.Printf("  [%d] %s (%s): $%.2f USD\n", i+1, token.Currency, name, token.LiquidityUSD)
	}

	// AMM pool count verification
	var poolStats struct {
		TotalPools   int     `gorm:"column:total_pools"`
		AvgLiquidity float64 `gorm:"column:avg_liquidity"`
		TotalPoolLiq float64 `gorm:"column:total_pool_liq"`
	}

	err = calc.db.Raw(`
		SELECT 
			COUNT(*) as total_pools,
			AVG(liquidity_usd) as avg_liquidity,
			SUM(liquidity_usd) as total_pool_liq
		FROM xrpAmm_normalized 
		WHERE liquidity_usd IS NOT NULL AND liquidity_usd > 0
	`).Scan(&poolStats).Error

	if err != nil {
		return fmt.Errorf("failed to get pool statistics: %v", err)
	}

	fmt.Printf("\n🏊 AMM Pool Statistics:\n")
	fmt.Printf("  Total AMM pools: %d\n", poolStats.TotalPools)
	fmt.Printf("  Average pool liquidity: $%.2f\n", poolStats.AvgLiquidity)
	fmt.Printf("  Total pool liquidity: $%.2f\n", poolStats.TotalPoolLiq)

	return nil
}

// RunLiquidityCalculation is a convenience function to run liquidity calculation
func RunLiquidityCalculation(db *gorm.DB) error {
	calculator := NewLiquidityCalculator(db)
	return calculator.PopulateAllTokenLiquidity()
}

















