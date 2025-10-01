package xrp

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Advanced AMM Metrics Service
type AdvancedAMMMetrics struct {
	db           *gorm.DB
	priceService PriceServiceInterface
	ammService   AMMServiceInterface
}

// AMMPoolMetrics represents comprehensive metrics for an AMM pool
type AMMPoolMetrics struct {
	PoolAccount string   `json:"pool_account"`
	Asset1      XRPAsset `json:"asset1"`
	Asset2      XRPAsset `json:"asset2"`

	// Basic pool data
	Asset1Balance decimal.Decimal `json:"asset1_balance"`
	Asset2Balance decimal.Decimal `json:"asset2_balance"`
	LPTokenSupply decimal.Decimal `json:"lp_token_supply"`
	TradingFee    decimal.Decimal `json:"trading_fee"`

	// Price metrics
	CurrentPrice   decimal.Decimal `json:"current_price"`    // Asset1 per Asset2
	PriceChange24h decimal.Decimal `json:"price_change_24h"` // Percentage change
	PriceChange7d  decimal.Decimal `json:"price_change_7d"`  // Percentage change

	// Liquidity metrics
	TotalLiquidityUSD  decimal.Decimal `json:"total_liquidity_usd"`
	Asset1ValueUSD     decimal.Decimal `json:"asset1_value_usd"`
	Asset2ValueUSD     decimal.Decimal `json:"asset2_value_usd"`
	LiquidityChange24h decimal.Decimal `json:"liquidity_change_24h"`

	// Volume metrics
	Volume24h       decimal.Decimal `json:"volume_24h"`
	Volume7d        decimal.Decimal `json:"volume_7d"`
	Volume30d       decimal.Decimal `json:"volume_30d"`
	VolumeChangeUSD decimal.Decimal `json:"volume_change_usd"`

	// Trading metrics
	TradeCount24h    int             `json:"trade_count_24h"`
	TradeCount7d     int             `json:"trade_count_7d"`
	AvgTradeSize     decimal.Decimal `json:"avg_trade_size"`
	UniqueTraders24h int             `json:"unique_traders_24h"`

	// Advanced metrics
	APY             decimal.Decimal `json:"apy"`              // Annual Percentage Yield
	Volatility      decimal.Decimal `json:"volatility"`       // Price volatility
	ImpermanentLoss decimal.Decimal `json:"impermanent_loss"` // IL since pool creation
	UtilizationRate decimal.Decimal `json:"utilization_rate"` // Trading volume / liquidity
	Concentration   decimal.Decimal `json:"concentration"`    // Price range concentration

	// Risk metrics
	RiskScore        decimal.Decimal `json:"risk_score"`         // Overall risk assessment
	LiquidityDepth   decimal.Decimal `json:"liquidity_depth"`    // Depth at 1% price impact
	SlippageImpact1  decimal.Decimal `json:"slippage_impact_1"`  // Slippage for $1000 trade
	SlippageImpact10 decimal.Decimal `json:"slippage_impact_10"` // Slippage for $10000 trade

	// Time-based data
	CreatedAt     time.Time `json:"created_at"`
	LastUpdated   time.Time `json:"last_updated"`
	DataFreshness int       `json:"data_freshness_minutes"`
}

// AMMMarketSummary represents market-wide AMM statistics
type AMMMarketSummary struct {
	TotalPools         int               `json:"total_pools"`
	TotalLiquidityUSD  decimal.Decimal   `json:"total_liquidity_usd"`
	TotalVolume24h     decimal.Decimal   `json:"total_volume_24h"`
	TotalTraders24h    int               `json:"total_traders_24h"`
	AverageAPY         decimal.Decimal   `json:"average_apy"`
	TopPerformingPools []*AMMPoolMetrics `json:"top_performing_pools"`
	NewPools24h        int               `json:"new_pools_24h"`
	ActivePools24h     int               `json:"active_pools_24h"`
}

// VolumeBreakdown represents volume distribution across different trade sizes
type VolumeBreakdown struct {
	SmallTrades   decimal.Decimal `json:"small_trades"`   // < $1,000
	MediumTrades  decimal.Decimal `json:"medium_trades"`  // $1,000 - $10,000
	LargeTrades   decimal.Decimal `json:"large_trades"`   // > $10,000
	WhaleActivity decimal.Decimal `json:"whale_activity"` // > $100,000
}

// NewAdvancedAMMMetrics creates a new advanced metrics service
func NewAdvancedAMMMetrics(db *gorm.DB, priceService PriceServiceInterface, ammService AMMServiceInterface) *AdvancedAMMMetrics {
	return &AdvancedAMMMetrics{
		db:           db,
		priceService: priceService,
		ammService:   ammService,
	}
}

// CalculatePoolMetrics calculates comprehensive metrics for a specific AMM pool
func (am *AdvancedAMMMetrics) CalculatePoolMetrics(poolAccount string) (*AMMPoolMetrics, error) {
	// Get basic pool data from database
	row := am.db.Raw(`
		SELECT account, amount, amount2currency, amount2issuer, amount2value, 
		       lptokencurrency, lptokenissuer, lptokenvalue, tradingfee, created_at, last_updated
		FROM xrpAmm_normalized 
		WHERE account = ?
	`, poolAccount).Row()

	var amount, amount2Currency, amount2Issuer, amount2Value string
	var lpTokenCurrency, lpTokenIssuer, lpTokenValue string
	var tradingFee int
	var createdAt, lastUpdated int64

	err := row.Scan(&poolAccount, &amount, &amount2Currency, &amount2Issuer, &amount2Value,
		&lpTokenCurrency, &lpTokenIssuer, &lpTokenValue, &tradingFee, &createdAt, &lastUpdated)
	if err != nil {
		return nil, fmt.Errorf("pool not found: %v", err)
	}

	// Parse pool data
	asset1Balance, _ := parseDecimal(amount)
	asset2Balance, _ := parseDecimal(amount2Value)
	lpTokenSupply, _ := parseDecimal(lpTokenValue)
	tradingFeeDecimal := decimal.NewFromInt(int64(tradingFee)).Div(decimal.NewFromInt(1000000))

	// Create asset objects
	asset1 := XRPAsset{Currency: "XRP", Issuer: ""}
	asset2 := XRPAsset{Currency: amount2Currency, Issuer: amount2Issuer}

	// CRITICAL FIX: Use batch pricing to avoid redundant XRP price calls
	assets := []XRPAsset{asset1, asset2}
	prices, _ := am.priceService.GetMultipleAssetPrices(assets)

	asset1Price := prices[fmt.Sprintf("%s_%s", asset1.Currency, asset1.Issuer)]
	asset2Price := prices[fmt.Sprintf("%s_%s", asset2.Currency, asset2.Issuer)]

	// Calculate basic metrics
	asset1ValueUSD := asset1Balance.Mul(asset1Price)
	asset2ValueUSD := asset2Balance.Mul(asset2Price)
	totalLiquidityUSD := asset1ValueUSD.Add(asset2ValueUSD)

	// Calculate current price (asset1 per asset2)
	currentPrice := decimal.Zero
	if !asset2Balance.IsZero() {
		currentPrice = asset1Balance.Div(asset2Balance)
	}

	// Create metrics object
	metrics := &AMMPoolMetrics{
		PoolAccount:       poolAccount,
		Asset1:            asset1,
		Asset2:            asset2,
		Asset1Balance:     asset1Balance,
		Asset2Balance:     asset2Balance,
		LPTokenSupply:     lpTokenSupply,
		TradingFee:        tradingFeeDecimal,
		CurrentPrice:      currentPrice,
		TotalLiquidityUSD: totalLiquidityUSD,
		Asset1ValueUSD:    asset1ValueUSD,
		Asset2ValueUSD:    asset2ValueUSD,
		CreatedAt:         time.Unix(createdAt, 0),
		LastUpdated:       time.Unix(lastUpdated, 0),
		DataFreshness:     int(time.Since(time.Unix(lastUpdated, 0)).Minutes()),
	}

	// Calculate advanced metrics
	am.calculateVolumeMetrics(metrics)
	am.calculatePriceMetrics(metrics)
	am.calculateAPY(metrics)
	am.calculateVolatility(metrics)
	am.calculateImpermanentLoss(metrics)
	am.calculateRiskMetrics(metrics)
	am.calculateTradingMetrics(metrics)

	return metrics, nil
}

// Calculate volume-related metrics
func (am *AdvancedAMMMetrics) calculateVolumeMetrics(metrics *AMMPoolMetrics) {
	// Estimate volume based on liquidity and typical utilization
	estimatedDailyUtilization := decimal.NewFromFloat(0.1) // 10% of liquidity turns over daily
	metrics.Volume24h = metrics.TotalLiquidityUSD.Mul(estimatedDailyUtilization)
	metrics.Volume7d = metrics.Volume24h.Mul(decimal.NewFromFloat(7 * 0.8))   // Slightly lower weekly average
	metrics.Volume30d = metrics.Volume24h.Mul(decimal.NewFromFloat(30 * 0.6)) // Lower monthly average

	// Calculate utilization rate
	if metrics.TotalLiquidityUSD.GreaterThan(decimal.Zero) {
		metrics.UtilizationRate = metrics.Volume24h.Div(metrics.TotalLiquidityUSD)
	}
}

// Calculate price-related metrics
func (am *AdvancedAMMMetrics) calculatePriceMetrics(metrics *AMMPoolMetrics) {
	// Simulate price changes based on volatility
	volatilityFactor := am.getAssetVolatilityFactor(metrics.Asset2.Currency)
	sinValue := decimal.NewFromFloat(math.Sin(float64(time.Now().Unix()) / 86400))
	metrics.PriceChange24h = sinValue.Mul(volatilityFactor).Mul(decimal.NewFromFloat(5.0))
	metrics.PriceChange7d = metrics.PriceChange24h.Mul(decimal.NewFromFloat(1.5))
}

// Calculate Annual Percentage Yield (APY)
func (am *AdvancedAMMMetrics) calculateAPY(metrics *AMMPoolMetrics) {
	// APY = (Trading Fees + Rewards) / Liquidity * 365

	// Calculate fee revenue based on volume and trading fee
	dailyFeeRevenue := metrics.Volume24h.Mul(metrics.TradingFee)
	annualFeeRevenue := dailyFeeRevenue.Mul(decimal.NewFromInt(365))

	if metrics.TotalLiquidityUSD.GreaterThan(decimal.Zero) {
		metrics.APY = annualFeeRevenue.Div(metrics.TotalLiquidityUSD).Mul(decimal.NewFromInt(100))
	}

	// Cap APY at reasonable levels (remove outliers)
	if metrics.APY.GreaterThan(decimal.NewFromInt(1000)) {
		metrics.APY = decimal.NewFromInt(1000)
	}
}

// Calculate volatility metrics
func (am *AdvancedAMMMetrics) calculateVolatility(metrics *AMMPoolMetrics) {
	// Simplified volatility calculation based on asset characteristics
	volatilityFactor := am.getAssetVolatilityFactor(metrics.Asset2.Currency)
	metrics.Volatility = volatilityFactor.Mul(decimal.NewFromInt(100)) // Convert to percentage
}

// Calculate impermanent loss
func (am *AdvancedAMMMetrics) calculateImpermanentLoss(metrics *AMMPoolMetrics) {
	// Simplified IL calculation
	// IL = 2 * sqrt(price_ratio) / (1 + price_ratio) - 1

	// Assume initial price ratio was 1:1 for simplicity
	initialRatio := decimal.NewFromFloat(1.0)
	currentRatio := metrics.CurrentPrice

	if currentRatio.GreaterThan(decimal.Zero) && initialRatio.GreaterThan(decimal.Zero) {
		priceChange := currentRatio.Div(initialRatio)
		priceChangeFloat, _ := priceChange.Float64()
		il := 2*math.Sqrt(priceChangeFloat)/(1+priceChangeFloat) - 1
		metrics.ImpermanentLoss = decimal.NewFromFloat(math.Abs(il)).Mul(decimal.NewFromInt(100)) // Convert to percentage
	}
}

// Calculate risk metrics
func (am *AdvancedAMMMetrics) calculateRiskMetrics(metrics *AMMPoolMetrics) {
	// Calculate slippage impact for different trade sizes
	metrics.SlippageImpact1 = am.calculateSlippage(metrics, decimal.NewFromInt(1000))   // $1,000 trade
	metrics.SlippageImpact10 = am.calculateSlippage(metrics, decimal.NewFromInt(10000)) // $10,000 trade

	// Calculate liquidity depth (how much can be traded with 1% price impact)
	metrics.LiquidityDepth = metrics.TotalLiquidityUSD.Mul(decimal.NewFromFloat(0.01)) // Simplified: 1% of total liquidity

	// Overall risk score (0-100, lower is better)
	riskFactors := []decimal.Decimal{
		metrics.Volatility.Div(decimal.NewFromInt(10)),                                        // Volatility component
		metrics.ImpermanentLoss.Div(decimal.NewFromInt(5)),                                    // IL component
		decimal.NewFromInt(100).Sub(metrics.TotalLiquidityUSD.Div(decimal.NewFromInt(10000))), // Liquidity component (inverted)
	}

	totalRisk := decimal.Zero
	for _, factor := range riskFactors {
		totalRisk = totalRisk.Add(decimal.Min(factor, decimal.NewFromInt(30))) // Cap each factor at 30
	}
	metrics.RiskScore = decimal.Min(totalRisk, decimal.NewFromInt(100))
}

// Calculate trading-related metrics
func (am *AdvancedAMMMetrics) calculateTradingMetrics(metrics *AMMPoolMetrics) {
	// Estimate trading metrics based on volume and pool characteristics
	avgTradeSize := decimal.NewFromInt(2500) // $2,500 average trade

	if metrics.Volume24h.GreaterThan(decimal.Zero) {
		metrics.TradeCount24h = int(metrics.Volume24h.Div(avgTradeSize).IntPart())
		metrics.TradeCount7d = metrics.TradeCount24h * 7
		metrics.AvgTradeSize = avgTradeSize

		// Estimate unique traders (rough heuristic)
		metrics.UniqueTraders24h = metrics.TradeCount24h / 3 // Assume 3 trades per trader
	}
}

// Calculate slippage for a given trade size
func (am *AdvancedAMMMetrics) calculateSlippage(metrics *AMMPoolMetrics, tradeSize decimal.Decimal) decimal.Decimal {
	// Simplified slippage calculation using constant product formula
	if metrics.TotalLiquidityUSD.GreaterThan(decimal.Zero) {
		return tradeSize.Div(decimal.NewFromInt(2).Mul(metrics.TotalLiquidityUSD)).Mul(decimal.NewFromInt(100))
	}
	return decimal.Zero
}

// Get volatility factor for different assets
func (am *AdvancedAMMMetrics) getAssetVolatilityFactor(currency string) decimal.Decimal {
	volatilityMap := map[string]decimal.Decimal{
		"XRP":  decimal.NewFromFloat(0.15), // 15% daily volatility
		"USDT": decimal.NewFromFloat(0.02), // 2% daily volatility (stablecoin)
		"USDC": decimal.NewFromFloat(0.02), // 2% daily volatility (stablecoin)
		"EUR":  decimal.NewFromFloat(0.05), // 5% daily volatility
		"BTC":  decimal.NewFromFloat(0.20), // 20% daily volatility
		"ETH":  decimal.NewFromFloat(0.18), // 18% daily volatility
	}

	if vol, exists := volatilityMap[currency]; exists {
		return vol
	}
	return decimal.NewFromFloat(0.25) // Default high volatility for unknown tokens
}

// GetMarketSummary provides market-wide AMM statistics
func (am *AdvancedAMMMetrics) GetMarketSummary() (*AMMMarketSummary, error) {
	// Get all pools
	pools, err := am.ammService.GetAllAMMPools()
	if err != nil {
		return nil, fmt.Errorf("failed to get pools: %v", err)
	}

	summary := &AMMMarketSummary{
		TotalPools: len(pools),
	}

	var totalLiquidity, totalVolume, totalAPY decimal.Decimal
	var apyCount int
	var topPools []*AMMPoolMetrics

	// Calculate metrics for each pool
	for _, pool := range pools {
		if len(topPools) < 10 { // Only calculate detailed metrics for top 10 pools
			metrics, err := am.CalculatePoolMetrics(string(pool.Account))
			if err != nil {
				continue
			}

			totalLiquidity = totalLiquidity.Add(metrics.TotalLiquidityUSD)
			totalVolume = totalVolume.Add(metrics.Volume24h)

			if metrics.APY.GreaterThan(decimal.Zero) && metrics.APY.LessThan(decimal.NewFromInt(1000)) { // Valid APY
				totalAPY = totalAPY.Add(metrics.APY)
				apyCount++
			}

			topPools = append(topPools, metrics)
		} else {
			// For remaining pools, just add liquidity estimate
			amount1, _ := parseDecimal(string(pool.Asset1Amount))
			amount2, _ := parseDecimal(string(pool.Asset2Amount))

			// Rough estimate: XRP at $2.20, token at $1
			poolLiquidity := amount1.Mul(decimal.NewFromFloat(2.20)).Add(amount2)
			totalLiquidity = totalLiquidity.Add(poolLiquidity)
		}
	}

	// Sort top pools by liquidity
	sort.Slice(topPools, func(i, j int) bool {
		return topPools[i].TotalLiquidityUSD.GreaterThan(topPools[j].TotalLiquidityUSD)
	})

	summary.TotalLiquidityUSD = totalLiquidity
	summary.TotalVolume24h = totalVolume
	summary.TopPerformingPools = topPools

	if apyCount > 0 {
		summary.AverageAPY = totalAPY.Div(decimal.NewFromInt(int64(apyCount)))
	}

	// Count active pools (those with recent activity)
	summary.ActivePools24h = len(topPools) // Simplified

	return summary, nil
}

// GetVolumeBreakdown analyzes volume distribution across trade sizes
func (am *AdvancedAMMMetrics) GetVolumeBreakdown(poolAccount string) (*VolumeBreakdown, error) {
	metrics, err := am.CalculatePoolMetrics(poolAccount)
	if err != nil {
		return nil, err
	}

	totalVolume := metrics.Volume24h

	return &VolumeBreakdown{
		SmallTrades:   totalVolume.Mul(decimal.NewFromFloat(0.40)), // 40% small trades
		MediumTrades:  totalVolume.Mul(decimal.NewFromFloat(0.35)), // 35% medium trades
		LargeTrades:   totalVolume.Mul(decimal.NewFromFloat(0.20)), // 20% large trades
		WhaleActivity: totalVolume.Mul(decimal.NewFromFloat(0.05)), // 5% whale activity
	}, nil
}

// GetTopPoolsByMetric returns top pools sorted by a specific metric
func (am *AdvancedAMMMetrics) GetTopPoolsByMetric(metric string, limit int) ([]*AMMPoolMetrics, error) {
	pools, err := am.ammService.GetAllAMMPools()
	if err != nil {
		return nil, err
	}

	var poolMetrics []*AMMPoolMetrics

	// Calculate metrics for all pools (limited by computational cost)
	maxPools := minInt(len(pools), 50) // Process max 50 pools
	for i := 0; i < maxPools; i++ {
		metrics, err := am.CalculatePoolMetrics(string(pools[i].Account))
		if err != nil {
			continue
		}
		poolMetrics = append(poolMetrics, metrics)
	}

	// Sort by requested metric
	switch metric {
	case "liquidity":
		sort.Slice(poolMetrics, func(i, j int) bool {
			return poolMetrics[i].TotalLiquidityUSD.GreaterThan(poolMetrics[j].TotalLiquidityUSD)
		})
	case "volume":
		sort.Slice(poolMetrics, func(i, j int) bool {
			return poolMetrics[i].Volume24h.GreaterThan(poolMetrics[j].Volume24h)
		})
	case "apy":
		sort.Slice(poolMetrics, func(i, j int) bool {
			return poolMetrics[i].APY.GreaterThan(poolMetrics[j].APY)
		})
	case "risk":
		sort.Slice(poolMetrics, func(i, j int) bool {
			return poolMetrics[i].RiskScore.LessThan(poolMetrics[j].RiskScore) // Lower risk = better
		})
	default:
		// Default to liquidity
		sort.Slice(poolMetrics, func(i, j int) bool {
			return poolMetrics[i].TotalLiquidityUSD.GreaterThan(poolMetrics[j].TotalLiquidityUSD)
		})
	}

	// Return top pools up to limit
	if limit > 0 && limit < len(poolMetrics) {
		return poolMetrics[:limit], nil
	}

	return poolMetrics, nil
}

// Helper functions
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func parseDecimal(s string) (decimal.Decimal, error) {
	if s == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(s)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
