package xrp

import "github.com/shopspring/decimal"

// AMMPriceCalculator provides methods for calculating prices within AMM pools.
type AMMPriceCalculator struct{}

// NewAMMPriceCalculator creates a new AMMPriceCalculator.
func NewAMMPriceCalculator() *AMMPriceCalculator {
	return &AMMPriceCalculator{}
}

// CalculatePrice calculates the price of one asset in terms of another.
func (c *AMMPriceCalculator) CalculatePrice(balanceA, balanceB decimal.Decimal) decimal.Decimal {
	if balanceB.IsZero() {
		return decimal.Zero
	}
	return balanceA.Div(balanceB)
}
