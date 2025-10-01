package xrp

// XRPOverview represents the QuantifyCrypto XRP overview used by the Terminal page
type XRPOverview struct {
	// Pricing
	PriceUSD     float64 `json:"priceUSD"`
	PriceBTC     float64 `json:"priceBTC"`
	Change24hPct float64 `json:"change24hPct"`
	High24hUSD   float64 `json:"high24hUSD"`
	Low24hUSD    float64 `json:"low24hUSD"`

	// Market stats
	MarketCapUSD      float64 `json:"marketCapUSD"`
	Volume24hUSD      float64 `json:"volume24hUSD"`
	CirculatingSupply float64 `json:"circulatingSupply"`

	// Technicals
	SupportUSD     float64 `json:"supportUSD"`
	ResistanceUSD  float64 `json:"resistanceUSD"`
	SafetyScore    float64 `json:"safetyScore"`
	TechnicalScore float64 `json:"technicalScore"`

	// QC Trend scores (0-100)
	TrendScores map[string]float64 `json:"trendScores"`
}
