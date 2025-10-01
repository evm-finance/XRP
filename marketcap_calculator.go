package xrp

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MarketCapCalculator handles market cap calculation using the formula: price * supply = marketcap
// This implements the final step of Phase 2: XRP Screener Enhancement Spec
type MarketCapCalculator struct {
	db *gorm.DB
}

// NewMarketCapCalculator creates a new market cap calculator
func NewMarketCapCalculator(db *gorm.DB) *MarketCapCalculator {
	return &MarketCapCalculator{
		db: db,
	}
}

// RecalculateMarketCaps updates market cap for all tokens using price * supply formula
func (calc *MarketCapCalculator) RecalculateMarketCaps() error {
	fmt.Println("🚀 Recalculating Market Caps using Price * Supply Formula")
	fmt.Println("=========================================================")
	fmt.Println("📋 Implementing XRP_SCREENER_ENHANCEMENT_SPEC.md Phase 2 Final Step")
	fmt.Println()

	// Get tokens that have both price and supply data
	tokens, err := calc.getTokensWithPriceAndSupply()
	if err != nil {
		return fmt.Errorf("failed to get tokens with price and supply data: %w", err)
	}

	if len(tokens) == 0 {
		fmt.Println("⚠️ No tokens found with both price and supply data")
		return nil
	}

	fmt.Printf("📊 Found %d tokens with price and supply data for market cap calculation\n", len(tokens))

	// Process each token
	updated := 0
	errors := 0

	for i, token := range tokens {
		fmt.Printf("[%d/%d] %s (%s)...", i+1, len(tokens), token.Currency, token.Name)

		success := calc.updateTokenMarketCap(token)
		if success {
			fmt.Printf(" ✅ Updated\n")
			updated++
		} else {
			fmt.Printf(" ❌ Failed\n")
			errors++
		}
	}

	// Summary
	fmt.Printf("\n🎯 MARKET CAP RECALCULATION COMPLETED\n")
	fmt.Printf("====================================\n")
	fmt.Printf("📊 Tokens Processed: %d\n", len(tokens))
	fmt.Printf("✅ Successfully Updated: %d\n", updated)
	fmt.Printf("❌ Errors: %d\n", errors)
	fmt.Printf("📈 Success Rate: %.1f%%\n", float64(updated)/float64(len(tokens))*100)

	return nil
}

// TokenMarketCapData represents token data needed for market cap calculation
type TokenMarketCapData struct {
	Currency string          `gorm:"column:currency"`
	Issuer   string          `gorm:"column:issuer"`
	Name     string          `gorm:"column:name"`
	Price    decimal.Decimal `gorm:"column:price"`
	Supply   decimal.Decimal `gorm:"column:supply_xrpl"`
}

// getTokensWithPriceAndSupply fetches tokens that have both price and supply data
func (calc *MarketCapCalculator) getTokensWithPriceAndSupply() ([]TokenMarketCapData, error) {
	var tokens []TokenMarketCapData

	query := `
		SELECT currency, issuer, token_name as name, price, supply_xrpl
		FROM xrpTokens
		WHERE currency != '' AND issuer != ''
		AND price > 0
		AND supply_xrpl IS NOT NULL
		AND supply_xrpl > 0
		ORDER BY supply_xrpl DESC
		LIMIT 100`

	err := calc.db.Raw(query).Scan(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query tokens with price and supply: %w", err)
	}

	return tokens, nil
}

// updateTokenMarketCap calculates and updates market cap for a single token
func (calc *MarketCapCalculator) updateTokenMarketCap(token TokenMarketCapData) bool {
	// Calculate market cap: price * supply
	marketCap := token.Price.Mul(token.Supply)

	// Update database with calculated market cap
	updateSQL := `
		UPDATE xrpTokens 
		SET marketcap = ?
		WHERE currency = ? AND issuer = ?`

	marketCapFloat, _ := marketCap.Float64()

	err := calc.db.Exec(updateSQL, marketCapFloat, token.Currency, token.Issuer).Error
	if err != nil {
		log.Printf("Failed to update market cap for %s/%s: %v", token.Currency, token.Issuer, err)
		return false
	}

	// Log the calculation details
	priceFloat, _ := token.Price.Float64()
	supplyFloat, _ := token.Supply.Float64()

	if len(token.Name) > 0 {
		log.Printf("📊 %s (%s): $%.6f × %.2f = $%.2f",
			token.Currency, token.Name, priceFloat, supplyFloat, marketCapFloat)
	}

	return true
}

// RecalculateMarketCapForToken updates market cap for a specific token
func (calc *MarketCapCalculator) RecalculateMarketCapForToken(currency, issuer string) error {
	// Get token data
	var token TokenMarketCapData
	err := calc.db.Raw(`
		SELECT currency, issuer, name, price, supply_xrpl
		FROM xrpTokens 
		WHERE currency = ? AND issuer = ?
		AND price > 0 
		AND supply_xrpl IS NOT NULL 
		AND supply_xrpl > 0
	`, currency, issuer).Scan(&token).Error

	if err != nil {
		return fmt.Errorf("failed to get token data: %w", err)
	}

	// Calculate and update market cap
	success := calc.updateTokenMarketCap(token)
	if !success {
		return fmt.Errorf("failed to update market cap for %s/%s", currency, issuer)
	}

	fmt.Printf("✅ Market cap updated for %s using price * supply formula\n", currency)
	return nil
}

// GetMarketCapCalculationPreview shows what the market cap calculation would produce
func (calc *MarketCapCalculator) GetMarketCapCalculationPreview() error {
	fmt.Println("🔍 Market Cap Calculation Preview")
	fmt.Println("=================================")

	tokens, err := calc.getTokensWithPriceAndSupply()
	if err != nil {
		return fmt.Errorf("failed to get tokens: %w", err)
	}

	if len(tokens) == 0 {
		fmt.Println("⚠️ No tokens available for market cap calculation")
		return nil
	}

	// Show preview for first 5 tokens
	previewCount := 5
	if len(tokens) < previewCount {
		previewCount = len(tokens)
	}

	fmt.Printf("📊 Preview of market cap calculations (top %d tokens):\n\n", previewCount)

	for i := 0; i < previewCount; i++ {
		token := tokens[i]

		// Calculate market cap
		marketCap := token.Price.Mul(token.Supply)

		priceFloat, _ := token.Price.Float64()
		supplyFloat, _ := token.Supply.Float64()
		marketCapFloat, _ := marketCap.Float64()

		name := token.Name
		if name == "" {
			name = "Unknown"
		}

		fmt.Printf("  [%d] %s (%s)\n", i+1, token.Currency, name)
		fmt.Printf("      Price: $%.6f\n", priceFloat)
		fmt.Printf("      Supply: %.2f tokens\n", supplyFloat)
		fmt.Printf("      Market Cap: $%.2f (price × supply)\n", marketCapFloat)
		fmt.Printf("\n")
	}

	fmt.Printf("💡 Formula: Market Cap = Price × Supply\n")
	fmt.Printf("📋 Ready to update %d tokens with calculated market caps\n", len(tokens))

	return nil
}
















