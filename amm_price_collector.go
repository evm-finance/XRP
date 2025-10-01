package xrp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"qc-defi-graphql-server/internal/models"

	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

// AMMPriceCollector collects and manages token prices from AMM pools
type AMMPriceCollector struct {
	db           *gorm.DB
	ammService   AMMServiceInterface
	priceService PriceServiceInterface
	cache        map[string]TokenPrice
	cacheMutex   sync.RWMutex
	lastUpdate   time.Time
	updateMutex  sync.Mutex
}

// TokenPrice represents a token's price data
type TokenPrice struct {
	Currency    string          `json:"currency"`
	Issuer      string          `json:"issuer"`
	PriceUSD    decimal.Decimal `json:"price_usd"`
	PriceXRP    decimal.Decimal `json:"price_xrp"`
	Source      string          `json:"source"` // "amm_xrp", "amm_rlusd", "calculated"
	PoolID      string          `json:"pool_id"`
	Liquidity   decimal.Decimal `json:"liquidity"`
	Volume24h   decimal.Decimal `json:"volume_24h"`
	LastUpdated time.Time       `json:"last_updated"`
}

// PriceCollectionResult represents the result of a price collection run
type PriceCollectionResult struct {
	TotalTokens    int                    `json:"total_tokens"`
	PricedTokens   int                    `json:"priced_tokens"`
	FailedTokens   int                    `json:"failed_tokens"`
	Prices         map[string]TokenPrice  `json:"prices"`
	Errors         map[string]string      `json:"errors"`
	CollectionTime time.Duration          `json:"collection_time"`
	Summary        map[string]interface{} `json:"summary"`
}

// NewAMMPriceCollector creates a new AMM price collector
func NewAMMPriceCollector(db *gorm.DB, ammService AMMServiceInterface, priceService PriceServiceInterface) *AMMPriceCollector {
	return &AMMPriceCollector{
		db:           db,
		ammService:   ammService,
		priceService: priceService,
		cache:        make(map[string]TokenPrice),
		lastUpdate:   time.Time{},
	}
}

// CollectPricesFromAMMPools collects prices for all tokens using discovered AMM pools
func (c *AMMPriceCollector) CollectPricesFromAMMPools() (*PriceCollectionResult, error) {
	c.updateMutex.Lock()
	defer c.updateMutex.Unlock()

	startTime := time.Now()
	result := &PriceCollectionResult{
		Prices: make(map[string]TokenPrice),
		Errors: make(map[string]string),
	}

	log.Println("🔍 Starting AMM price collection...")

	// Get all AMM pools from normalized database with clean data (FIXED: correct column names)
	rows, err := c.db.Raw(`
		SELECT 
			account,
			asset1_amount,
			asset2_currency,
			asset2_issuer,
			asset2_amount,
			lp_token_currency,
			lp_token_issuer,
			lp_token_amount,
			trading_fee
		FROM xrpAmm_normalized 
		WHERE account IS NOT NULL
		ORDER BY liquidity_usd DESC
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AMM pools: %w", err)
	}
	defer rows.Close()

	var ammPools []struct {
		Account         string
		Asset1Amount    string
		Asset2Currency  string
		Asset2Issuer    string
		Asset2Amount    string
		LpTokenCurrency string
		LpTokenIssuer   string
		LpTokenAmount   string
		TradingFee      int
	}

	for rows.Next() {
		var pool struct {
			Account         string
			Asset1Amount    string
			Asset2Currency  string
			Asset2Issuer    string
			Asset2Amount    string
			LpTokenCurrency string
			LpTokenIssuer   string
			LpTokenAmount   string
			TradingFee      int
		}
		err := rows.Scan(
			&pool.Account,
			&pool.Asset1Amount,
			&pool.Asset2Currency,
			&pool.Asset2Issuer,
			&pool.Asset2Amount,
			&pool.LpTokenCurrency,
			&pool.LpTokenIssuer,
			&pool.LpTokenAmount,
			&pool.TradingFee,
		)
		if err != nil {
			log.Printf("❌ Failed to scan AMM pool row: %v", err)
			continue
		}
		ammPools = append(ammPools, pool)
	}

	result.TotalTokens = len(ammPools)
	log.Printf("📊 Found %d AMM pools to process", len(ammPools))

	// CRITICAL FIX: Get XRP price ONCE for all pools
	xrpPriceUSD, err := c.priceService.GetXRPLPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRP price: %w", err)
	}

	// Activate per-batch override so nested calls reuse this price without XRPL/QC
	c.priceService.SetBatchXRPPrice(xrpPriceUSD)
	defer c.priceService.ClearBatchXRPPrice()
	log.Printf("💵 [AMM COLLECTOR] Got XRP price once: $%s (saving %d redundant calls)",
		xrpPriceUSD.StringFixed(6), len(ammPools)-1)

	// Process each pool to extract price information
	for _, pool := range ammPools {
		tokenPrice, err := c.extractPriceFromPoolData(pool, xrpPriceUSD)
		if err != nil {
			key := fmt.Sprintf("%s_%s", pool.Asset2Currency, pool.Asset2Issuer)
			result.Errors[key] = err.Error()
			result.FailedTokens++
			// Log removed to reduce log noise
			// log.Printf("❌ Failed to extract price from pool %s: %v", pool.Account, err)
			continue
		}

		if tokenPrice != nil {
			key := fmt.Sprintf("%s_%s", tokenPrice.Currency, tokenPrice.Issuer)
			result.Prices[key] = *tokenPrice
			c.cacheMutex.Lock()
			c.cache[key] = *tokenPrice
			c.cacheMutex.Unlock()
			result.PricedTokens++
		}
	}

	// Update collection timestamp
	c.lastUpdate = time.Now()
	result.CollectionTime = time.Since(startTime)

	// Generate summary
	result.Summary = c.generateSummary(result)

	log.Printf("✅ Price collection complete: %d priced, %d failed, took %v",
		result.PricedTokens, result.FailedTokens, result.CollectionTime)

	return result, nil
}

// extractPriceFromPoolData extracts price information from raw AMM pool data (FIXED: correct column names)
func (c *AMMPriceCollector) extractPriceFromPoolData(pool struct {
	Account         string
	Asset1Amount    string
	Asset2Currency  string
	Asset2Issuer    string
	Asset2Amount    string
	LpTokenCurrency string
	LpTokenIssuer   string
	LpTokenAmount   string
	TradingFee      int
}, xrpPriceUSD decimal.Decimal) (*TokenPrice, error) {
	// FIXED: Database schema uses asset1_amount/asset2_amount naming
	// Asset1Amount = XRP amount (always)
	// Asset2Currency/Issuer/Amount = Token details (always)

	tokenCurrency := pool.Asset2Currency
	tokenIssuer := pool.Asset2Issuer

	// CRITICAL FIX: Parse XRP amount (database stores in drops format, needs conversion to XRP)
	xrpAmountDrops, err := decimal.NewFromString(pool.Asset1Amount)
	if err != nil {
		log.Printf("❌ [PRICE COLLECTOR] Failed to parse XRP amount '%s': %v", pool.Asset1Amount, err)
		return nil, fmt.Errorf("failed to parse XRP amount: %w", err)
	}
	// CRITICAL: Convert drops to XRP - database stores XRP amounts in drops format
	xrpAmount := xrpAmountDrops.Div(decimal.NewFromFloat(1000000.0))

	// Parse token amount
	tokenAmount, err := decimal.NewFromString(pool.Asset2Amount)
	if err != nil {
		log.Printf("❌ [PRICE COLLECTOR] Failed to parse token amount '%s': %v", pool.Asset2Amount, err)
		return nil, fmt.Errorf("failed to parse token amount: %w", err)
	}

	// Calculate token price: XRP amount / Token amount = XRP per token
	if xrpAmount.LessThanOrEqual(decimal.Zero) || tokenAmount.LessThanOrEqual(decimal.Zero) {
		log.Printf("❌ [PRICE COLLECTOR] Invalid amounts: XRP=%s, Token=%s", xrpAmount.String(), tokenAmount.String())
		return nil, fmt.Errorf("invalid pool amounts: XRP=%s, token=%s", xrpAmount.String(), tokenAmount.String())
	}

	// Simple AMM price using constant product formula
	// Price = amount of base currency / amount of quote currency
	// For XRP/Token pool: Token price in XRP = XRP amount / Token amount
	tokenPriceInXRP := xrpAmount.Div(tokenAmount)

	// FIXED: Remove incorrect "unreasonable" price validation
	// Many XRPL tokens legitimately have microscopic prices due to high supply or micro-cap status
	// Previous threshold (0.000001 XRP) incorrectly rejected 98%+ of valid token prices
	// Only reject truly impossible values (negative, zero, or calculation errors)
	if tokenPriceInXRP.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("invalid price calculation: %s XRP per token", tokenPriceInXRP.String())
	}

	// Accept ALL positive prices, even microscopic ones (e.g., 0.0000000000000001 XRP)
	// XRPL tokens with trillion+ supply legitimately have tiny per-unit values

	// Convert to USD price using provided XRP price (already fetched once)
	priceUSD := tokenPriceInXRP.Mul(xrpPriceUSD)

	// Calculate liquidity: XRP amount * XRP price
	liquidity := xrpAmount.Mul(xrpPriceUSD)

	return &TokenPrice{
		Currency:    tokenCurrency,
		Issuer:      tokenIssuer,
		PriceUSD:    priceUSD,
		PriceXRP:    tokenPriceInXRP,
		Source:      "amm_xrp",
		PoolID:      pool.Account,
		Liquidity:   liquidity,
		Volume24h:   decimal.Zero, // Would need to calculate from transactions
		LastUpdated: time.Now(),
	}, nil
}

// extractPriceFromPool extracts price information from a single AMM pool (original method for compatibility)
func (c *AMMPriceCollector) extractPriceFromPool(pool models.AMMInfo) (*TokenPrice, error) {
	// Determine which asset is the token and which is the base (XRP or RLUSD)
	var tokenCurrency, tokenIssuer, baseCurrency, baseIssuer string
	var tokenAmount, baseAmount decimal.Decimal
	var err error

	// Check if asset1 is XRP (base)
	if pool.GetAmountValue() != "" {
		// Asset1 is XRP
		baseCurrency = "XRP"
		baseIssuer = ""
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()

		// Parse amounts
		baseAmount, err = decimal.NewFromString(pool.GetAmountValue())
		if err != nil {
			return nil, fmt.Errorf("failed to parse XRP amount: %w", err)
		}
		baseAmount = baseAmount.Div(decimal.NewFromFloat(1000000.0)) // Convert drops to XRP from GetAmountValue()

		// Handle both string and object formats for Amount2
		var strValue string
		if err := json.Unmarshal(pool.Result.Amm.Amount2, &strValue); err == nil {
			// It's a string (XRP amount in drops), convert to XRP
			if val, err := decimal.NewFromString(strValue); err == nil {
				tokenAmount = val.Div(decimal.NewFromFloat(1000000.0)) // Convert drops to XRP when Amount2 is XRP
			}
		} else {
			// Try to unmarshal as an object (token amounts)
			var amount2 struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(pool.Result.Amm.Amount2, &amount2); err != nil {
				return nil, fmt.Errorf("failed to parse token amount: %w", err)
			}
			if val, err := decimal.NewFromString(amount2.Value); err == nil {
				tokenAmount = val
			}
		}
	} else {
		// Asset2 is XRP, Asset1 is the token
		baseCurrency = pool.GetAmount2Currency()
		baseIssuer = pool.GetAmount2Issuer()
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()

		// This case needs more complex parsing based on the pool structure
		return nil, fmt.Errorf("complex pool structure not yet supported")
	}

	// Calculate token price in base currency
	if baseAmount.LessThanOrEqual(decimal.Zero) || tokenAmount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("invalid pool amounts: base=%s, token=%s", baseAmount.String(), tokenAmount.String())
	}

	tokenPriceInBase := baseAmount.Div(tokenAmount)

	// Convert to USD price
	var priceUSD decimal.Decimal
	var source string

	if baseCurrency == "XRP" {
		// Get XRP price in USD
		xrpPrice, err := c.priceService.GetXRPLPrice()
		if err != nil {
			return nil, fmt.Errorf("failed to get XRP price: %w", err)
		}
		priceUSD = tokenPriceInBase.Mul(xrpPrice)
		source = "amm_xrp"
	} else if baseCurrency == RLUSDCurrency && baseIssuer == RLUSDIssuer {
		// RLUSD is approximately $1
		priceUSD = tokenPriceInBase
		source = "amm_rlusd"
	} else {
		// Try to get base currency price
		baseAsset := XRPAsset{Currency: baseCurrency, Issuer: baseIssuer}
		basePrice, err := c.priceService.GetAssetPrice(baseAsset)
		if err != nil {
			return nil, fmt.Errorf("failed to get base currency price: %w", err)
		}
		priceUSD = tokenPriceInBase.Mul(basePrice)
		source = "calculated"
	}

	// Calculate liquidity (simplified)
	liquidity := baseAmount.Mul(tokenPriceInBase)

	return &TokenPrice{
		Currency:    tokenCurrency,
		Issuer:      tokenIssuer,
		PriceUSD:    priceUSD,
		PriceXRP:    tokenPriceInBase,
		Source:      source,
		PoolID:      pool.Result.Amm.Account,
		Liquidity:   liquidity,
		Volume24h:   decimal.Zero, // Would need to calculate from transactions
		LastUpdated: time.Now(),
	}, nil
}

// GetTokenPrice retrieves the current price for a specific token
func (c *AMMPriceCollector) GetTokenPrice(currency, issuer string) (*TokenPrice, error) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	key := fmt.Sprintf("%s_%s", currency, issuer)
	if price, exists := c.cache[key]; exists {
		return &price, nil
	}

	return nil, fmt.Errorf("price not found for token %s/%s", currency, issuer)
}

// GetAllPrices returns all cached token prices
func (c *AMMPriceCollector) GetAllPrices() map[string]TokenPrice {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	result := make(map[string]TokenPrice)
	for k, v := range c.cache {
		result[k] = v
	}
	return result
}

// RefreshPrices forces a refresh of all prices
func (c *AMMPriceCollector) RefreshPrices() error {
	_, err := c.CollectPricesFromAMMPools()
	return err
}

// GetLastUpdateTime returns when prices were last updated
func (c *AMMPriceCollector) GetLastUpdateTime() time.Time {
	return c.lastUpdate
}

// generateSummary generates a summary of the price collection results
func (c *AMMPriceCollector) generateSummary(result *PriceCollectionResult) map[string]interface{} {
	summary := make(map[string]interface{})

	// Price statistics
	var totalLiquidity, avgPrice, minPrice, maxPrice decimal.Decimal
	var priceCount int

	for _, price := range result.Prices {
		if price.PriceUSD.GreaterThan(decimal.Zero) {
			totalLiquidity = totalLiquidity.Add(price.Liquidity)
			avgPrice = avgPrice.Add(price.PriceUSD)
			if minPrice.IsZero() || price.PriceUSD.LessThan(minPrice) {
				minPrice = price.PriceUSD
			}
			if price.PriceUSD.GreaterThan(maxPrice) {
				maxPrice = price.PriceUSD
			}
			priceCount++
		}
	}

	if priceCount > 0 {
		avgPrice = avgPrice.Div(decimal.NewFromInt(int64(priceCount)))
	}

	// Source breakdown
	sourceBreakdown := make(map[string]int)
	for _, price := range result.Prices {
		sourceBreakdown[price.Source]++
	}

	summary["total_liquidity"] = totalLiquidity
	summary["average_price"] = avgPrice
	summary["min_price"] = minPrice
	summary["max_price"] = maxPrice
	summary["source_breakdown"] = sourceBreakdown
	summary["success_rate"] = float64(result.PricedTokens) / float64(result.TotalTokens) * 100

	return summary
}

// StorePricesToDatabase stores the collected prices to the database
func (c *AMMPriceCollector) StorePricesToDatabase(prices map[string]TokenPrice) error {
	log.Printf("💾 [PRICE STORAGE] Storing %d token prices to xrpToken table", len(prices))

	for _, price := range prices {
		// Calculate market cap properly: price * supply
		// First get the current supply for this token
		var supply decimal.Decimal
		err := c.db.Raw(`
			SELECT COALESCE(supply_xrpl, 0) as supply 
			FROM xrpTokens 
			WHERE currency = ? AND issuer = ?
		`, price.Currency, price.Issuer).Scan(&supply).Error

		if err != nil {
			log.Printf("⚠️ [PRICE STORAGE] Failed to get supply for %s/%s, using 0: %v", price.Currency, price.Issuer, err)
			supply = decimal.Zero
		}

		// Calculate market cap: price * supply (proper formula)
		var marketCap decimal.Decimal
		if supply.GreaterThan(decimal.Zero) && price.PriceUSD.GreaterThan(decimal.Zero) {
			marketCap = price.PriceUSD.Mul(supply)
		} else {
			marketCap = decimal.Zero // Set to 0 if no supply data available
		}

		// Update existing xrpToken records with price data and calculated market cap
		err = c.db.Exec(`
			UPDATE xrpTokens 
			SET price = ?, marketcap = ?, volume_24h = ?, updated_at = NOW()
			WHERE currency = ? AND issuer = ?
		`, price.PriceUSD, marketCap, price.Volume24h, price.Currency, price.Issuer).Error

		if err != nil {
			log.Printf("❌ [PRICE STORAGE] Failed to update price for %s/%s: %v", price.Currency, price.Issuer, err)
			continue
		}

		// Log the calculation details
		// Log removed to reduce log noise
		// supplyFloat, _ := supply.Float64()
		// marketCapFloat, _ := marketCap.Float64()
		// log.Printf("✅ [PRICE STORAGE] Updated %s/%s: Price=$%s, Supply=%.2f, MarketCap=$%.2f",
		// 	price.Currency, price.Issuer, price.PriceUSD.StringFixed(6), supplyFloat, marketCapFloat)
	}

	log.Printf("✅ [PRICE STORAGE] Price storage complete")
	return nil
}

// GetPriceHistory retrieves price history for a token
func (c *AMMPriceCollector) GetPriceHistory(currency, issuer string, limit int) ([]TokenPrice, error) {
	var result []TokenPrice

	// Query xrpTokens table for price history
	query := `SELECT currency, issuer, price, marketcap, volume_24h as volume24h, updated_at 
			  FROM xrpTokens 
			  WHERE currency = ? AND issuer = ?
			  ORDER BY updated_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := c.db.Raw(query, currency, issuer).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token TokenPrice
		var updatedAt time.Time

		err := rows.Scan(&token.Currency, &token.Issuer, &token.PriceUSD, &token.Liquidity, &token.Volume24h, &updatedAt)
		if err != nil {
			continue
		}

		token.LastUpdated = updatedAt
		result = append(result, token)
	}

	return result, nil
}
