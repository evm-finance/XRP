package xrp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"qc-defi-graphql-server/internal/models"
	"sort"
	"time"

	"gorm.io/gorm"
)

// HybridScreenerService combines external API data with database fallback
type HybridScreenerService struct {
	db         *gorm.DB
	httpClient *http.Client
}

// NewHybridScreenerService creates a new hybrid screener service
func NewHybridScreenerService(db *gorm.DB) *HybridScreenerService {
	return &HybridScreenerService{
		db: db,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetScreenerTokens serves tokens from database (enhanced with external data via background sync)
func (s *HybridScreenerService) GetScreenerTokens() ([]*models.XRPTokenFields, error) {
	log.Printf("🔍 [HYBRID SCREENER] Serving tokens from database (enhanced with external data)")

	// Always serve from database (which is populated by external sync job)
	return s.getTokensFromDatabase()
}

// fetchFromExternalAPI fetches fresh data from XRPL Meta API
func (s *HybridScreenerService) fetchFromExternalAPI() ([]*models.XRPTokenFields, error) {
	log.Printf("🌐 [HYBRID SCREENER] Fetching from https://s1.xrplmeta.org/tokens")

	// Fetch XRP data from database (not external API)
	xrpToken, err := s.getXRPDataFromDatabase()
	if err != nil {
		log.Printf("⚠️ [HYBRID SCREENER] Failed to fetch XRP data from database: %v", err)
	}

	// Fetch other tokens from XRPL Meta
	resp, err := s.httpClient.Get("https://s1.xrplmeta.org/tokens")
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenList models.XRPTokenList
	if err := json.Unmarshal(body, &tokenList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to our format
	var result []*models.XRPTokenFields

	// Add XRP first if available
	if xrpToken != nil {
		result = append(result, xrpToken)
	}

	// Add other tokens with enhanced data from external API
	for _, token := range tokenList.Tokens {
		// Skip tokens without names (poor quality)
		if token.Meta.Token.Name == "" {
			continue
		}

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.Issuer,
			TokenName:     token.Meta.Token.Name,  // ✅ From external API
			IssuerName:    token.Meta.Issuer.Name, // ✅ From external API
			Icon:          token.Meta.Token.Icon,
			Marketcap:     token.Metrics.Marketcap, // ✅ From external API
			Price:         token.Metrics.Price,     // ✅ From external API
			Supply:        0,                       // Not available from external API
			Liquidity:     token.LiquidityUSD,      // ✅ From external API
			Volume24H:     token.Metrics.Volume24H, // ✅ From external API
		})

		// Limit to top 100 tokens
		if len(result) >= 100 {
			break
		}
	}

	// Sort by market cap (largest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Marketcap > result[j].Marketcap
	})

	log.Printf("✅ [HYBRID SCREENER] External API returned %d enhanced tokens", len(result))
	return result, nil
}

// getTokensFromDatabase gets tokens from database (enhanced with external API data)
func (s *HybridScreenerService) getTokensFromDatabase() ([]*models.XRPTokenFields, error) {
	log.Printf("📊 [HYBRID SCREENER] Querying database with external API enhancements")

	var tokens []struct {
		Currency     string  `gorm:"column:currency"`
		Issuer       string  `gorm:"column:issuer"`
		TokenName    string  `gorm:"column:token_name"`
		IssuerName   string  `gorm:"column:issuer_name"`
		Icon         string  `gorm:"column:icon"`
		Price        float64 `gorm:"column:price"`
		Marketcap    float64 `gorm:"column:marketcap"`
		SupplyXrpl   float64 `gorm:"column:supply_xrpl"`
		Volume24H    float64 `gorm:"column:volume_24h"`
	}

	// Query from xrpTokens table which now has external data (issuer_name, external marketcap/volume)
	query := `
		SELECT 
			currency, issuer, 
			COALESCE(token_name, currency) as token_name,
			COALESCE(issuer_name, '') as issuer_name,
			COALESCE(icon, '') as icon,
			COALESCE(price, 0) as price,
			COALESCE(marketcap, 0) as marketcap,
			COALESCE(supply_xrpl, 0) as supply_xrpl,
			COALESCE(volume_24h, 0) as volume_24h
		FROM xrpTokens t
		WHERE EXISTS (
			SELECT 1 FROM xrpAmm_normalized amm 
			WHERE (amm.asset1_currency = t.currency AND amm.asset1_issuer = t.issuer)
			   OR (amm.asset2_currency = t.currency AND amm.asset2_issuer = t.issuer)
			   AND amm.liquidity_usd > 1000
		)
		ORDER BY 
			CASE WHEN marketcap > 0 THEN marketcap ELSE 0 END DESC,
			price DESC, 
			trustlines DESC
		LIMIT 50`

	err := s.db.Raw(query).Scan(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	// Get liquidity for each token from AMM pools and convert to API format
	var result []*models.XRPTokenFields
	for _, token := range tokens {
		// Get liquidity for this token from AMM pools
		liquidity := s.getTokenLiquidity(token.Currency, token.Issuer)

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.Issuer,
			TokenName:     token.TokenName,
			IssuerName:    token.IssuerName,     // ✅ From external API via sync job
			Icon:          token.Icon,
			Marketcap:     token.Marketcap,      // ✅ From external API via sync job
			Price:         token.Price,          // ✅ From our drops conversion fixes
			Supply:        token.SupplyXrpl,
			Liquidity:     liquidity,            // ✅ From our accurate AMM calculations
			Volume24H:     token.Volume24H,      // ✅ From external API via sync job
		})
	}

	log.Printf("✅ [HYBRID SCREENER] Database query returned %d tokens with external enhancements", len(result))
	return result, nil
}

// getTokenLiquidity calculates accurate liquidity from AMM pools
func (s *HybridScreenerService) getTokenLiquidity(currency, issuer string) float64 {
	var liquidityData struct {
		TotalLiquidity float64 `gorm:"column:total_liquidity"`
	}

	query := `
		SELECT SUM(liquidity_usd) as total_liquidity
		FROM xrpAmm_normalized 
		WHERE (asset1_currency = ? AND asset1_issuer = ?) 
		   OR (asset2_currency = ? AND asset2_issuer = ?)
		   AND liquidity_usd > 0`

	err := s.db.Raw(query, currency, issuer, currency, issuer).Scan(&liquidityData).Error
	if err != nil {
		return 0
	}

	return liquidityData.TotalLiquidity
}

// enhanceWithDatabaseLiquidity adds liquidity data from our AMM database
func (s *HybridScreenerService) enhanceWithDatabaseLiquidity(tokens []*models.XRPTokenFields) []*models.XRPTokenFields {
	log.Printf("🔧 [HYBRID SCREENER] Enhancing %d tokens with database liquidity data", len(tokens))

	for _, token := range tokens {
		if token.Currency == "XRP" {
			continue // Skip XRP
		}

		// Get liquidity for this token from our AMM pools
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

		err := s.db.Raw(liquidityQuery, token.Currency, token.IssuerAddress,
			token.Currency, token.IssuerAddress).Scan(&liquidityData).Error

		if err == nil && liquidityData.TotalLiquidity > 0 {
			// Use our calculated liquidity if it's higher (more accurate)
			if liquidityData.TotalLiquidity > token.Liquidity {
				log.Printf("🔧 [HYBRID SCREENER] Enhanced %s liquidity: %.2f → %.2f",
					token.Currency, token.Liquidity, liquidityData.TotalLiquidity)
				token.Liquidity = liquidityData.TotalLiquidity
			}
		}
	}

	return tokens
}

// getXRPDataFromDatabase gets XRP data from our database (using XRP/RLUSD price calculation)
func (s *HybridScreenerService) getXRPDataFromDatabase() (*models.XRPTokenFields, error) {
	var xrpData struct {
		Price     float64 `gorm:"column:price"`
		MarketCap float64 `gorm:"column:marketcap"`
		Supply    float64 `gorm:"column:supply_xrpl"`
		Volume24H float64 `gorm:"column:volume_24h"`
	}

	// Get XRP data from our database (populated by our price service)
	err := s.db.Raw(`
		SELECT 
			COALESCE(price, 0) as price,
			COALESCE(marketcap, 0) as marketcap,
			COALESCE(supply_xrpl, 0) as supply_xrpl,
			COALESCE(volume_24h, 0) as volume_24h
		FROM xrpTokens 
		WHERE currency = 'XRP' OR currency = ''
		LIMIT 1
	`).Scan(&xrpData).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get XRP data from database: %w", err)
	}

	return &models.XRPTokenFields{
		Currency:      "XRP",
		IssuerAddress: "",
		TokenName:     "XRP",
		IssuerName:    "XRP Ledger",
		Icon:          "",
		Marketcap:     xrpData.MarketCap,  // From database
		Price:         xrpData.Price,      // From our XRP/RLUSD calculation
		Supply:        xrpData.Supply,     // From database
		Liquidity:     0,                  // XRP doesn't have AMM pools
		Volume24H:     xrpData.Volume24H,  // From external sync if available
	}, nil
}
