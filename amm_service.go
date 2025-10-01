package xrp

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"qc-defi-graphql-server/internal/models"

	"github.com/shopspring/decimal"
	xrpl "github.com/xrpscan/xrpl-go"
	"gorm.io/gorm"
)

// --- Lightweight AMM metrics (logs every 60s) ---
type ammErrorCounters struct {
	mu      sync.Mutex
	counts  map[string]int
	started bool
}

// --- Hourly logging of not-found pairs ---
type ammNotFoundTracker struct {
	mu             sync.Mutex
	notFoundPairs  map[string]int // pair -> count
	lastHourlyLog  time.Time
	firstBatchDone bool
}

var ammCounters ammErrorCounters
var notFoundTracker = &ammNotFoundTracker{}

func (c *ammErrorCounters) startTicker() {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return
	}
	c.started = true
	c.mu.Unlock()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			c.mu.Lock()
			if len(c.counts) > 0 {
				snapshot := c.counts
				c.counts = make(map[string]int)
				c.mu.Unlock()

				// Prefer a human-friendly line for actNotFound counts
				if nf, ok := snapshot["error:actNotFound"]; ok && nf > 0 {
					successCount := 0
					if sc, ok := snapshot["success"]; ok {
						successCount = sc
					}
					log.Printf("📊 [AMM METRICS 60s] %d pairs found, %d pairs not found", successCount, nf)
				}

				// Optionally log remaining metrics (excluding actNotFound and success) if any
				remaining := map[string]int{}
				for k, v := range snapshot {
					if k == "error:actNotFound" || k == "success" {
						continue
					}
					remaining[k] = v
				}
				if len(remaining) > 0 {
					log.Printf("📊 [AMM METRICS 60s] %v", remaining)
				}
			} else {
				c.mu.Unlock()
			}
		}
	}()
}

func (c *ammErrorCounters) incr(key string) {
	if !c.started {
		c.startTicker()
	}
	c.mu.Lock()
	if c.counts == nil {
		c.counts = make(map[string]int)
	}
	c.counts[key]++
	c.mu.Unlock()
}

// Methods for ammNotFoundTracker
func (t *ammNotFoundTracker) addNotFoundPair(pair string) {
	t.mu.Lock()
	if t.notFoundPairs == nil {
		t.notFoundPairs = make(map[string]int)
	}
	t.notFoundPairs[pair]++
	t.mu.Unlock()
}

func (t *ammNotFoundTracker) markFirstBatchDone() {
	t.mu.Lock()
	t.firstBatchDone = true
	t.mu.Unlock()
}

func (t *ammNotFoundTracker) shouldLogHourly() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.firstBatchDone {
		return false
	}

	now := time.Now()
	if t.lastHourlyLog.IsZero() || now.Sub(t.lastHourlyLog) >= time.Hour {
		t.lastHourlyLog = now
		return true
	}
	return false
}

func (t *ammNotFoundTracker) logAndClear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.notFoundPairs) == 0 {
		return
	}

	log.Printf("🕐 [AMM HOURLY] Not-found pairs (last hour):")
	for pair, count := range t.notFoundPairs {
		// Decode hex currency codes to ASCII for better readability
		decodedPair := t.decodeHexCurrencyPair(pair)
		log.Printf("   - %s (%s): %d times", decodedPair, pair, count)
	}

	// Clear the map for next hour
	t.notFoundPairs = make(map[string]int)
}

// decodeHexCurrencyPair decodes hex currency codes in a pair to ASCII
func (t *ammNotFoundTracker) decodeHexCurrencyPair(pair string) string {
	parts := strings.Split(pair, "/")
	if len(parts) != 2 {
		return pair // Return as-is if not in expected format
	}

	currency1, currency2 := parts[0], parts[1]

	// Decode currency2 if it's a 40-character hex string
	if len(currency2) == 40 {
		if decoded, err := hex.DecodeString(currency2); err == nil {
			// Find the end of the actual currency (before zero padding)
			var currencyStr string
			for i, b := range decoded {
				if b == 0 {
					currencyStr = string(decoded[:i])
					break
				}
			}
			if len(currencyStr) > 0 {
				return fmt.Sprintf("%s/%s", currency1, currencyStr)
			}
		}
	}

	return pair // Return original if decoding fails
}

// XRP constants
const (
	XRPAccount = "rsoLo2S1kiGeCcn6hCUXVrCpGMWLrRrLZz"
)

// AMMService provides AMM pool discovery and management functionality
type AMMService struct {
	db         *gorm.DB
	xrpConnMgr *ConnectionManager
	requestID  int
	idMutex    sync.Mutex
	// priceService temporarily removed to fix compilation - will add back later
}

// AMMServiceInterface defines the contract for AMM operations
type AMMServiceInterface interface {
	StoreAMMPools(pools []models.AMMInfo) error
	GetHeatmapData(filters HeatmapFilters) ([]AMMHeatmapData, error)
	UpdatePoolMetrics() error
	DiscoverAllAMMPools() ([]models.AMMInfo, error)
	CheckSingleAMMPool(currency1, issuer1, currency2, issuer2 string) (*models.AMMInfo, error)
	GetAllAMMPools() ([]AMMPoolRowForUpdate, error)
	GetAMMPoolInfo(account string) (*models.AMMInfo, error)
	GetAMMSwapQuote(asset1, asset2 map[string]string, amount string) (*models.AMMInfo, error)
	GetAMMLiquidityValue(account string) (float64, error)
	GetAMMTransactions(account string, limit int) ([]interface{}, error)
	CalculatePoolLiquidityUSD(pool models.AMMInfo) float64
	computeLiquidityByXRPAmount() error
	GetTopAMMPoolsByLiquidity(limit int) ([]AMMPoolRow, error)
	// New normalized table methods for improved deduplication and data quality
	StoreAMMPoolsNormalized(pools []models.AMMInfo) error
	GetNormalizedAMMPools(limit int) ([]map[string]interface{}, error)
	DeduplicateExistingData() error
	// Live price fetching methods (no hardcoded fallbacks)
	GetLiveXRPPrice() (decimal.Decimal, error)
}

// HeatmapFilters defines filtering options for AMM heatmap data
type HeatmapFilters struct {
	TimeRange    string  `json:"timeRange"`
	MinLiquidity float64 `json:"minLiquidity"`
	MaxLiquidity float64 `json:"maxLiquidity"`
	SortBy       string  `json:"sortBy"`
	SortOrder    string  `json:"sortOrder"`
	Limit        int     `json:"limit"`
	Offset       int     `json:"offset"`
}

// AMMHeatmapData represents heatmap data for AMM pools
type AMMHeatmapData struct {
	Token          models.XRPTokenData `json:"token"`
	Pools          []models.AMMInfo    `json:"pools"`
	TotalLiquidity float64             `json:"totalLiquidity"`
	Volume24h      float64             `json:"volume24h"`
	Volume7d       float64             `json:"volume7d"`
	MarketCap      float64             `json:"marketCap"`
	PriceChange24h float64             `json:"priceChange24h"`
}

// Move AMMPoolRow to package level so both methods can use it
//
//	type AMMPoolRow struct {
//		Account         string  `json:"account"`
//		Amount          string  `json:"amount"`
//		Amount2Currency string  `json:"amount2currency"`
//		Amount2Issuer   string  `json:"amount2issuer"`
//		Amount2Value    string  `json:"amount2value"`
//		Asset2Frozen    int     `json:"asset2frozen"`
//		LpTokenCurrency string  `json:"lptokencurrency"`
//		LpTokenIssuer   string  `json:"lptokenissuer"`
//		LpTokenValue    string  `json:"lptokenvalue"`
//		TradingFee      int     `json:"tradingfee"`
//		LiquidityUSD    float64 `json:"liquidity_usd"`
//		CreatedAt       int64   `json:"created_at"`
//		LastUpdated     int64   `json:"last_updated"`
//	}
type AMMPoolRow struct {
	Account         string  `json:"account"`
	Amount          string  `json:"amount"`
	Amount2Currency string  `json:"amount2currency"`
	Amount2Issuer   string  `json:"amount2issuer"`
	Amount2Value    string  `json:"amount2value"`
	Asset2Frozen    int     `json:"asset2frozen"`
	LpTokenCurrency string  `json:"lptokencurrency"`
	LpTokenIssuer   string  `json:"lptokenissuer"`
	LpTokenValue    string  `json:"lptokenvalue"`
	TradingFee      int     `json:"tradingfee"`
	LiquidityUSD    float64 `json:"liquidity_usd"`
	CreatedAt       int64   `json:"created_at"`
	LastUpdated     int64   `json:"last_updated"`
}

// NewAMMService creates a new AMM service instance
func NewAMMService(db *gorm.DB, connMgr *ConnectionManager) AMMServiceInterface {
	return &AMMService{
		db:         db,
		xrpConnMgr: connMgr,
		requestID:  1,
		// priceService: NewPriceServiceWithDB(db), // temporarily removed
	}
}

// getNextRequestID returns the next request ID
func (s *AMMService) getNextRequestID() int {
	s.idMutex.Lock()
	defer s.idMutex.Unlock()
	s.requestID++
	return s.requestID
}

// debugTraceMethodCall logs when any AMM service method is called for pipeline debugging
// DISABLED: Excessive logging removed to reduce noise
func (s *AMMService) debugTraceMethodCall(methodName string, args ...interface{}) {
	// Logging disabled to reduce informational noise
	_ = methodName
	_ = args
}

// CheckSingleAMMPool checks if a specific AMM pool exists and returns its information
func (s *AMMService) CheckSingleAMMPool(currency1, issuer1, currency2, issuer2 string) (*models.AMMInfo, error) {
	s.debugTraceMethodCall("CheckSingleAMMPool", currency1, issuer1, currency2, issuer2)
	// Try to get a connection with retries
	var client *xrpl.Client
	var err error

	// Retry up to 3 times if connection pool is full
	for attempt := 0; attempt < 3; attempt++ {
		client, err = s.xrpConnMgr.GetConnection()
		if err == nil {
			break
		}

		// If pool is full, wait a bit and retry
		if attempt < 2 {
			time.Sleep(time.Millisecond * 100 * time.Duration(attempt+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL connection after retries: %v", err)
	}

	// Always return the connection to the pool when done
	defer s.xrpConnMgr.ReturnConnection(client)

	return s.checkSingleAMMPoolWithClient(client, currency1, issuer1, currency2, issuer2)
}

// isHexString checks if a string is a valid hex string (40 characters)
func (s *AMMService) isHexString(str string) bool {
	if len(str) != 40 {
		return false
	}
	// Check if all characters are valid hex digits
	for _, char := range str {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}

// convertCurrencyToXRPLFormat converts currency codes to the format required by XRPL API
func (s *AMMService) convertCurrencyToXRPLFormat(currency string) string {
	// Standard 3-letter ISO codes work as-is
	if len(currency) <= 3 {
		return currency
	}

	// If it's already a 40-character hex string, use as-is
	if s.isHexString(currency) {
		return currency
	}

	// Non-standard currencies like RLUSD need hex encoding
	switch currency {
	case "RLUSD":
		return "524C555344000000000000000000000000000000" // RLUSD in hex
	default:
		// For other non-standard currencies, convert to hex if needed
		if len(currency) > 3 {
			// Convert string to hex (pad to 40 characters)
			hex := fmt.Sprintf("%X", currency)
			if len(hex) < 40 {
				hex = hex + strings.Repeat("0", 40-len(hex))
			}
			return hex
		}
		return currency
	}
}

// convertCurrencyToXRPLHex converts currency codes to hex format for database storage
// Returns empty string for standard 3-letter currencies that don't need hex encoding
func (s *AMMService) convertCurrencyToXRPLHex(currency string) string {
	// Standard 3-letter ISO codes don't need hex storage - store empty string
	if len(currency) <= 3 {
		return ""
	}

	// If it's already a 40-character hex string, store as-is
	if s.isHexString(currency) {
		return currency
	}

	// Non-standard currencies like RLUSD need hex encoding
	switch currency {
	case "RLUSD":
		return "524C555344000000000000000000000000000000" // RLUSD in hex
	default:
		// For other non-standard currencies, convert to hex
		if len(currency) > 3 {
			// Convert string to hex (pad to 40 characters)
			hex := fmt.Sprintf("%X", currency)
			if len(hex) < 40 {
				hex = hex + strings.Repeat("0", 40-len(hex))
			}
			return hex
		}
		return ""
	}
}

// getHexCurrencyFromDB retrieves hex currency format from database if available
func (s *AMMService) getHexCurrencyFromDB(currency, issuer string) string {
	// For XRP or standard 3-letter currencies, return as-is (no hex needed)
	if len(currency) <= 3 {
		return currency
	}

	// Check if we already have the hex format stored in database
	type CurrencyHex struct {
		Asset1CurrencyHex string `gorm:"column:asset1_currency_hex"`
		Asset2CurrencyHex string `gorm:"column:asset2_currency_hex"`
	}

	var result CurrencyHex
	query := `
		SELECT asset1_currency_hex, asset2_currency_hex 
		FROM xrpAmm_normalized 
		WHERE (asset1_currency = ? AND asset1_issuer = ?) 
		   OR (asset2_currency = ? AND asset2_issuer = ?)
		LIMIT 1
	`

	err := s.db.Raw(query, currency, issuer, currency, issuer).Scan(&result).Error
	if err != nil {
		// If database lookup fails, fall back to on-the-fly conversion
		return s.convertCurrencyToXRPLFormat(currency)
	}

	// Return the hex format from database if available
	if result.Asset1CurrencyHex != "" {
		return result.Asset1CurrencyHex
	}
	if result.Asset2CurrencyHex != "" {
		return result.Asset2CurrencyHex
	}

	// Fall back to on-the-fly conversion if no hex format in database
	return s.convertCurrencyToXRPLFormat(currency)
}

// checkSingleAMMPoolWithClient checks a specific AMM pool using a provided client (for batch operations)
func (s *AMMService) checkSingleAMMPoolWithClient(client *xrpl.Client, currency1, issuer1, currency2, issuer2 string) (*models.AMMInfo, error) {
	// Use hex currency formats from database for XRPL API compatibility
	xrplCurrency1 := s.getHexCurrencyFromDB(currency1, issuer1)
	xrplCurrency2 := s.getHexCurrencyFromDB(currency2, issuer2)

	asset1 := map[string]string{"currency": xrplCurrency1}
	if issuer1 != "" {
		asset1["issuer"] = issuer1
	}
	asset2 := map[string]string{"currency": xrplCurrency2}
	if issuer2 != "" {
		asset2["issuer"] = issuer2
	}
	request := map[string]interface{}{
		"command":      "amm_info",
		"asset":        asset1,
		"asset2":       asset2,
		"ledger_index": "validated",
	}
	response, err := client.Request(request)
	if err != nil {
		s.xrpConnMgr.MarkConnectionUnhealthy(client)
		ammCounters.incr("request_error")
		return nil, fmt.Errorf("XRPL API error: %v", err)
	}
	// Only log error responses below

	// Check for XRPL API errors
	if status, ok := response["status"].(string); ok && status == "error" {
		if errorCode, ok := response["error"].(string); ok {
			ammCounters.incr("error:" + errorCode)

			// Skip retries for actNotFound and issueMalformed errors
			if errorCode == "actNotFound" || errorCode == "issueMalformed" {
				// Track not-found pairs for hourly logging
				if errorCode == "actNotFound" {
					pair := fmt.Sprintf("%s/%s", currency1, currency2)
					notFoundTracker.addNotFoundPair(pair)
				}
				return nil, fmt.Errorf("XRPL API error: %s", errorCode)
			}

			// Fallbacks for other errors: try string currency and freshly derived ASCII→hex
			// (Removed actNotFound retry logic since we now skip retries for this error)

			return nil, fmt.Errorf("XRPL API error: %s", errorCode)
		}
		ammCounters.incr("error:unknown")
		return nil, fmt.Errorf("XRPL API error: unknown error")
	}

	// Check if we got a successful response
	if status, exists := response["status"]; !exists || status != "success" {
		ammCounters.incr("unexpected_status")
		return nil, fmt.Errorf("XRPL API error: unexpected status")
	}

	// Prefer validated responses; if result.validated == false, ignore as ambiguous
	if result, ok := response["result"].(map[string]interface{}); ok {
		if validated, okv := result["validated"].(bool); okv && !validated {
			ammCounters.incr("non_validated")
			return nil, fmt.Errorf("non-validated amm_info result")
		}
	}

	// Debug: Log the actual response structure (commented out to reduce log noise)
	// log.Printf("🔍 [DEBUG] AMM API response type: %T", response)
	// log.Printf("🔍 [DEBUG] AMM API response content: %+v", response)

	// Additional debug for parseAMMInfoResponse (commented out to reduce log noise)
	// log.Printf("🔍 [DEBUG] About to call parseAMMInfoResponse with type: %T", response)

	// Parse the response into AMMInfo
	ammInfo := parseAMMInfoResponse(response)
	if ammInfo == nil {
		return nil, fmt.Errorf("failed to parse AMM response")
	}

	// Count successful responses
	ammCounters.incr("success")

	return ammInfo, nil
}

// parseAMMInfoResponse parses and validates the XRPL response into AMMInfo
func parseAMMInfoResponse(response interface{}) *models.AMMInfo {
	// Basic validation of response structure
	if response == nil {
		log.Printf("❌ AMM response is nil")
		return nil
	}

	// Handle different response types more robustly
	var responseMap map[string]interface{}

	switch v := response.(type) {
	case map[string]interface{}:
		// log.Printf("✅ [DEBUG] Response is map[string]interface{}")
		responseMap = v
	case string:
		// log.Printf("🔍 [DEBUG] Response is string, attempting JSON parse")
		// Try to parse JSON string
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(v), &jsonResponse); err != nil {
			log.Printf("❌ AMM response is a string but not valid JSON: %s", v)
			return nil
		}
		responseMap = jsonResponse
	default:
		// Handle xrpl.BaseResponse and other types by converting to map
		// log.Printf("🔧 [DEBUG] Response is type %T, attempting to convert to map", response)

		// Use reflection to convert to map[string]interface{}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("❌ Failed to marshal response to JSON: %v", err)
			return nil
		}

		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(responseBytes, &jsonResponse); err != nil {
			log.Printf("❌ Failed to unmarshal response from JSON: %v", err)
			return nil
		}

		// log.Printf("✅ [DEBUG] Successfully converted %T to map[string]interface{}", response)
		responseMap = jsonResponse
	}

	// Validate status field
	status, ok := responseMap["status"].(string)
	if !ok {
		log.Printf("❌ AMM response missing status field")
		return nil
	}

	if status != "success" {
		errorMsg := "unknown error"
		if errorCode, exists := responseMap["error"].(string); exists {
			errorMsg = errorCode
		}
		log.Printf("❌ AMM API returned error status: %s (%s)", status, errorMsg)
		return nil
	}

	// Validate result structure exists
	resultData, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		log.Printf("❌ AMM response missing result field")
		return nil
	}

	// Validate AMM data exists
	ammData, ok := resultData["amm"].(map[string]interface{})
	if !ok {
		log.Printf("❌ AMM response missing amm data field")
		return nil
	}

	// Validate critical AMM fields
	account, ok := ammData["account"].(string)
	if !ok || account == "" {
		log.Printf("❌ AMM response missing or invalid account field")
		return nil
	}

	// Validate account format (basic XRPL address check)
	if len(account) < 25 || len(account) > 35 || account[0] != 'r' {
		log.Printf("❌ AMM account has invalid XRPL address format: %s", account)
		return nil
	}

	// log.Printf("✅ AMM response validation passed for pool: %s", account)
	// Marshal and unmarshal to handle the response structure
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ JSON marshal error: %v", err)
		return nil
	}

	var parsedResponse struct {
		Result struct {
			Amm struct {
				Account      string          `json:"account"`
				Amount       json.RawMessage `json:"amount"`  // Handle both string and object formats
				Amount2      json.RawMessage `json:"amount2"` // Handle both string and object formats
				Asset2Frozen bool            `json:"asset2_frozen"`
				AuctionSlot  struct {
					Account      string `json:"account"`
					AuthAccounts []struct {
						Account string `json:"account"`
					} `json:"auth_accounts"`
					DiscountedFee int    `json:"discounted_fee"`
					Expiration    string `json:"expiration"`
					Price         struct {
						Currency string `json:"currency"`
						Issuer   string `json:"issuer"`
						Value    string `json:"value"`
					} `json:"price"`
					TimeInterval int `json:"time_interval"`
				} `json:"auction_slot"`
				LpToken struct {
					Currency string `json:"currency"`
					Issuer   string `json:"issuer"`
					Value    string `json:"value"`
				} `json:"lp_token"`
				TradingFee int `json:"trading_fee"`
				VoteSlots  []struct {
					Account    string `json:"account"`
					TradingFee int    `json:"trading_fee"`
					VoteWeight int    `json:"vote_weight"`
				} `json:"vote_slots"`
			} `json:"amm"`
			LedgerCurrentIndex int  `json:"ledger_current_index"`
			Validated          bool `json:"validated"`
		} `json:"result"`
		Status string `json:"status"`
		Type   string `json:"type"`
	}

	if err := json.Unmarshal(responseJSON, &parsedResponse); err != nil {
		log.Printf("❌ JSON unmarshal error: %v", err)
		return nil
	}

	// Convert to AMMInfo - we need to handle the Amount2 field conversion
	ammInfo := &models.AMMInfo{
		Status: parsedResponse.Status,
		Type:   parsedResponse.Type,
	}

	// Copy the Result struct, handling Amount conversion
	ammInfo.Result.Amm.Account = parsedResponse.Result.Amm.Account

	// Convert Amount to RawMessage format - handle both string and object formats
	if len(parsedResponse.Result.Amm.Amount) > 0 {
		// Try to unmarshal as a string first (XRP amounts come as strings)
		var strValue string
		if err := json.Unmarshal(parsedResponse.Result.Amm.Amount, &strValue); err == nil {
			// It's a string (XRP amount in drops), create normalized structure
			amountData := map[string]string{
				"currency": "XRP",
				"issuer":   "",
				"value":    strValue,
			}
			amountJSON, _ := json.Marshal(amountData)
			ammInfo.Result.Amm.Amount = amountJSON
		} else {
			// Try to unmarshal as an object (token amounts come as objects)
			var objValue struct {
				Currency string `json:"currency"`
				Issuer   string `json:"issuer"`
				Value    string `json:"value"`
			}
			if err := json.Unmarshal(parsedResponse.Result.Amm.Amount, &objValue); err == nil {
				// It's an object, normalize the structure
				amountData := map[string]string{
					"currency": objValue.Currency,
					"issuer":   objValue.Issuer,
					"value":    objValue.Value,
				}
				amountJSON, _ := json.Marshal(amountData)
				ammInfo.Result.Amm.Amount = amountJSON
			} else {
				// Log parsing failure for debugging
				log.Printf("⚠️ [AMM PARSE] Failed to parse Amount field as string or object: %s", string(parsedResponse.Result.Amm.Amount))
				// Create a safe fallback structure
				amountData := map[string]string{
					"currency": "XRP",
					"issuer":   "",
					"value":    "0",
				}
				amountJSON, _ := json.Marshal(amountData)
				ammInfo.Result.Amm.Amount = amountJSON
			}
		}
	}

	ammInfo.Result.Amm.Asset2Frozen = parsedResponse.Result.Amm.Asset2Frozen
	ammInfo.Result.Amm.AuctionSlot = parsedResponse.Result.Amm.AuctionSlot

	// Convert Amount2 to RawMessage format - handle both string and object formats
	if len(parsedResponse.Result.Amm.Amount2) > 0 {
		// Try to unmarshal as a string first (XRP amounts come as strings)
		var strValue string
		if err := json.Unmarshal(parsedResponse.Result.Amm.Amount2, &strValue); err == nil {
			// It's a string (XRP amount in drops), create normalized structure
			amount2Data := map[string]string{
				"currency": "XRP",
				"issuer":   "",
				"value":    strValue,
			}
			amount2JSON, _ := json.Marshal(amount2Data)
			ammInfo.Result.Amm.Amount2 = amount2JSON
		} else {
			// Try to unmarshal as an object (token amounts come as objects)
			var objValue struct {
				Currency string `json:"currency"`
				Issuer   string `json:"issuer"`
				Value    string `json:"value"`
			}
			if err := json.Unmarshal(parsedResponse.Result.Amm.Amount2, &objValue); err == nil {
				// It's an object, normalize the structure
				amount2Data := map[string]string{
					"currency": objValue.Currency,
					"issuer":   objValue.Issuer,
					"value":    objValue.Value,
				}
				amount2JSON, _ := json.Marshal(amount2Data)
				ammInfo.Result.Amm.Amount2 = amount2JSON
			} else {
				// Log parsing failure for debugging
				log.Printf("⚠️ [AMM PARSE] Failed to parse Amount2 field as string or object: %s", string(parsedResponse.Result.Amm.Amount2))
				// Create a safe fallback structure
				amount2Data := map[string]string{
					"currency": "XRP",
					"issuer":   "",
					"value":    "0",
				}
				amount2JSON, _ := json.Marshal(amount2Data)
				ammInfo.Result.Amm.Amount2 = amount2JSON
			}
		}
	}

	// log.Printf("✅ AMM info parsed successfully for pool: %s", ammInfo.Result.Amm.Account)
	return ammInfo
}

// FindAMMPools finds AMM pools for a given token
func (s *AMMService) FindAMMPools(token models.XRPTokenData) ([]models.AMMInfo, error) {
	_, err := s.xrpConnMgr.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL connection: %v", err)
	}
	var pools []models.AMMInfo

	// Check for XRP pairs (native XRP) - both orders
	if token.Currency != "XRP" {
		// Order 1: Token/XRP
		pool, err := s.CheckSingleAMMPool(token.Currency, token.Issuer, "XRP", "")
		if err == nil && pool != nil {
			pools = append(pools, *pool)
		}

		// Order 2: XRP/Token
		pool, err = s.CheckSingleAMMPool("XRP", "", token.Currency, token.Issuer)
		if err == nil && pool != nil {
			pools = append(pools, *pool)
		}
	}

	// Check for RLUSD pairs - both orders
	if token.Currency != RLUSDCurrency {
		// Order 1: Token/RLUSD
		pool, err := s.CheckSingleAMMPool(token.Currency, token.Issuer, RLUSDCurrency, RLUSDIssuer)
		if err == nil && pool != nil {
			pools = append(pools, *pool)
		}

		// Order 2: RLUSD/Token
		pool, err = s.CheckSingleAMMPool(RLUSDCurrency, RLUSDIssuer, token.Currency, token.Issuer)
		if err == nil && pool != nil {
			pools = append(pools, *pool)
		}
	}

	return pools, nil
}

// DiscoverAllAMMPools discovers all AMM pools by querying tokens from the database
func (s *AMMService) DiscoverAllAMMPools() ([]models.AMMInfo, error) {
	var allPools []models.AMMInfo
	totalFoundPools := 0
	totalCheckedPools := 0
	totalStoredPools := 0
	totalFilteredPools := 0
	totalProcessedTokens := 0
	minLiquidityUSD := 100.0
	batchSize := 100
	offset := 0

	// AMM pool discovery started (logging reduced)

	for {
		// Get batch of tokens from database - fetch ALL rows and handle conversion in Go
		rows, err := s.db.Raw(`
			SELECT currency, issuer 
			FROM xrpTokens 
			LIMIT ? OFFSET ?
		`, batchSize, offset).Rows()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tokens batch (offset %d): %v", offset, err)
		}

		var tokens []struct {
			Currency string
			Issuer   string
		}

		// Use a map to track unique currency/issuer pairs
		uniqueTokens := make(map[string]bool)

		for rows.Next() {
			var currencyInterface interface{}
			var issuer string

			// Scan into interface{} first to handle byte arrays
			if err := rows.Scan(&currencyInterface, &issuer); err != nil {
				continue
			}

			// Convert currency to string, handling byte arrays
			var currency string
			switch v := currencyInterface.(type) {
			case string:
				currency = v
			case []byte:
				// Convert byte array to ASCII string for currency codes
				currency = string(v)
			default:
				// Try to convert to string
				currency = fmt.Sprintf("%v", v)
			}

			// Skip if currency is empty or invalid
			if currency == "" || currency == "NULL" {
				continue
			}

			// Create unique key for deduplication
			uniqueKey := currency + "|" + issuer
			if uniqueTokens[uniqueKey] {
				continue // Skip if we've already seen this currency/issuer pair
			}
			uniqueTokens[uniqueKey] = true

			tokens = append(tokens, struct {
				Currency string
				Issuer   string
			}{
				Currency: currency,
				Issuer:   issuer,
			})
		}
		rows.Close()

		// If no more tokens, break the loop
		if len(tokens) == 0 {
			break
		}

		// fmt.Printf("📦 Processing batch %d: %d tokens (offset: %d)\n", (offset/batchSize)+1, len(tokens), offset)

		// Log some sample tokens being processed (commented out to reduce log noise)
		// if len(tokens) > 0 {
		// 	fmt.Printf("📋 Sample tokens in this batch:\n")
		// 	for i := 0; i < min(5, len(tokens)); i++ {
		// 		fmt.Printf("   - Currency: %s, Issuer: %s\n", tokens[i].Currency, tokens[i].Issuer)
		// 	}
		// 	if len(tokens) > 5 {
		// 		fmt.Printf("   ... and %d more tokens\n", len(tokens)-5)
		// 	}
		// }

		// Process this batch of tokens using connection pooling optimization
		batchFoundPools := 0
		batchCheckedPools := 0
		batchStoredPools := 0
		batchFilteredPools := 0

		// Get XRP price ONCE for the entire batch to avoid redundant calls
		var batchXRPPrice decimal.Decimal
		batchXRPPrice, err = s.GetLiveXRPPrice()
		if err != nil {
			fmt.Printf("❌ Failed to get XRP price for batch: %v\n", err)
			continue // Skip this batch if we can't get XRP price
		}

		// Mark first batch as done for hourly logging
		notFoundTracker.markFirstBatchDone()

		// Get a single connection for the entire batch to avoid pool exhaustion
		client, err := s.xrpConnMgr.GetConnection()
		if err != nil {
			fmt.Printf("❌ Failed to get XRPL connection for batch: %v\n", err)
			continue // Skip this batch
		}

		// Process all tokens in this batch with the same connection
		for _, token := range tokens {
			// Token checking (debug logging removed)

			// Check for XRP pairs - both orders
			if token.Currency != "XRP" {
				// Order 1: Token/XRP
				pool, err := s.checkSingleAMMPoolWithClient(client, token.Currency, token.Issuer, "XRP", "")
				batchCheckedPools++
				if err == nil && pool != nil {
					// Check liquidity before storing
					if s.checkAndStorePoolWithLiquidityAndXRPPrice(*pool, minLiquidityUSD, batchXRPPrice) {
						allPools = append(allPools, *pool)
						batchFoundPools++
						batchStoredPools++
						// fmt.Printf("✅ Stored XRP pool: %s/%s - Liquidity: $%.2f USD\n", token.Currency, "XRP", liquidityUSD)
					} else {
						batchFilteredPools++
						// fmt.Printf("❌ Filtered XRP pool: %s/%s - Liquidity: $%.2f USD (below $%.2f threshold)\n", token.Currency, "XRP", liquidityUSD, minLiquidityUSD)
					}
				}

				// Order 2: XRP/Token
				pool, err = s.checkSingleAMMPoolWithClient(client, "XRP", "", token.Currency, token.Issuer)
				batchCheckedPools++
				if err == nil && pool != nil {
					// Check liquidity before storing
					if s.checkAndStorePoolWithLiquidityAndXRPPrice(*pool, minLiquidityUSD, batchXRPPrice) {
						allPools = append(allPools, *pool)
						batchFoundPools++
						batchStoredPools++
						// fmt.Printf("✅ Stored XRP pool: %s/%s - Liquidity: $%.2f USD\n", "XRP", token.Currency, liquidityUSD)
					} else {
						batchFilteredPools++
						// fmt.Printf("❌ Filtered XRP pool: %s/%s - Liquidity: $%.2f USD (below $%.2f threshold)\n", "XRP", token.Currency, liquidityUSD, minLiquidityUSD)
					}
				}
			}

			// Only check XRP pairs - RLUSD pairs are not needed for this system
		}

		// Return the connection to the pool after processing the entire batch
		s.xrpConnMgr.ReturnConnection(client)

		// Update totals
		totalFoundPools += batchFoundPools
		totalCheckedPools += batchCheckedPools
		totalStoredPools += batchStoredPools
		totalFilteredPools += batchFilteredPools
		totalProcessedTokens += len(tokens)

		// Check for hourly logging of not-found pairs
		if notFoundTracker.shouldLogHourly() {
			notFoundTracker.logAndClear()
		}

		// fmt.Printf("✅ Batch %d complete: checked %d pools, found %d pools, stored %d pools, filtered %d pools\n",
		// 	(offset/batchSize)+1, batchCheckedPools, batchFoundPools, batchStoredPools, batchFilteredPools)

		// Move to next batch
		offset += batchSize
	}

	// fmt.Printf("🎉 Discovery complete! Processed %d tokens in %d batches.\n", totalProcessedTokens, (offset / batchSize))
	// fmt.Printf("📊 Totals: checked %d pools, found %d AMM pools, stored %d pools (min liquidity: $%.2f USD), filtered out %d pools\n",
	// 	totalCheckedPools, totalFoundPools, totalStoredPools, minLiquidityUSD, totalFilteredPools)

	return allPools, nil
}

// checkAndStorePoolWithLiquidityAndXRPPrice checks if a pool meets the minimum liquidity threshold and stores it if it does
// Uses pre-fetched XRP price to avoid redundant price calls
func (s *AMMService) checkAndStorePoolWithLiquidityAndXRPPrice(pool models.AMMInfo, minLiquidityUSD float64, xrpPriceUSD decimal.Decimal) bool {
	// Convert string amounts to float64 and handle drops conversion
	amount1Drops, _ := strconv.ParseFloat(pool.GetAmountValue(), 64)
	amount2, _ := strconv.ParseFloat(pool.GetAmount2Value(), 64)

	// Convert XRP drops to XRP (GetAmountValue returns drops)
	amount1 := amount1Drops / 1000000.0

	// Determine which is XRP/RLUSD and which is the token
	var baseAmount float64
	var tokenCurrency, tokenIssuer string
	var isXRPBase bool

	if pool.GetAmount2Currency() == "XRP" {
		// amount is token, amount2 is XRP
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = true
	} else if pool.GetAmount2Currency() == RLUSDCurrency && pool.GetAmount2Issuer() == RLUSDIssuer {
		// amount is token, amount2 is RLUSD
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = false
	} else {
		// amount is XRP/RLUSD, amount2 is token
		baseAmount = amount1
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = (pool.GetAmountValue() != "" && pool.GetAmount2Currency() != RLUSDCurrency)
	}

	// XRP amounts from GetAmountValue() are in drops format, need conversion
	if isXRPBase && baseAmount > 0 {
		// CRITICAL: Convert drops to XRP for liquidity calculations
		baseAmount = baseAmount / 1000000.0
	}

	// Calculate liquidity in USD using provided XRP price
	var liquidityUSD decimal.Decimal
	if isXRPBase {
		// For XRP base, use provided XRP price
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2)) // Double the base amount for total liquidity
	} else {
		// For RLUSD base, RLUSD is approximately $1
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(decimal.NewFromInt(2)) // Double the base amount for total liquidity
	}

	// Check if liquidity meets minimum threshold
	if liquidityUSD.LessThan(decimal.NewFromFloat(minLiquidityUSD)) {
		// Pool filtered below threshold (logging reduced)
		return false
	}

	// Check if pool already exists in database to prevent duplicates
	poolExists, err := s.CheckPoolExistsInDB(pool.Result.Amm.Account)
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to check pool existence for %s: %v\n", pool.Result.Amm.Account, err)
		// Continue with storage in case of check error
	} else if poolExists {
		// fmt.Printf("⏭️  Pool %s already exists in database, skipping duplicate\n", pool.Result.Amm.Account)
		return false
	}

	// Also check by currency pair to catch different orderings
	pairExists, _, err := s.CheckPoolExistsByCurrencyPair("XRP", "", tokenCurrency, tokenIssuer)
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to check currency pair existence: %v\n", err)
	} else if pairExists {
		// fmt.Printf("⏭️  Pool pair XRP/%s already exists, skipping duplicate\n", tokenCurrency)
		return false
	}

	// Store the pool in the database using StoreAMMPoolsNormalized
	err = s.StoreAMMPoolsNormalized([]models.AMMInfo{pool})
	if err != nil {
		fmt.Printf("❌ Failed to store pool %s: %v\n", pool.Result.Amm.Account, err)
		return false
	}

	// fmt.Printf("✅ Successfully stored pool %s with liquidity $%.2f USD\n", pool.Result.Amm.Account, liquidityUSD.InexactFloat64())
	return true
}

// checkAndStorePoolWithLiquidity checks if a pool meets the minimum liquidity threshold and stores it if it does
func (s *AMMService) checkAndStorePoolWithLiquidity(pool models.AMMInfo, minLiquidityUSD float64) bool {
	// Convert string amounts to float64 and handle drops conversion
	amount1Drops, _ := strconv.ParseFloat(pool.GetAmountValue(), 64)
	amount2, _ := strconv.ParseFloat(pool.GetAmount2Value(), 64)

	// Convert XRP drops to XRP (GetAmountValue returns drops)
	amount1 := amount1Drops / 1000000.0

	// Determine which is XRP/RLUSD and which is the token
	var baseAmount, tokenAmount float64
	var tokenCurrency, tokenIssuer string
	var isXRPBase bool

	if pool.GetAmount2Currency() == "XRP" {
		// amount is token, amount2 is XRP
		tokenAmount = amount1
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = true
	} else if pool.GetAmount2Currency() == RLUSDCurrency && pool.GetAmount2Issuer() == RLUSDIssuer {
		// amount is token, amount2 is RLUSD
		tokenAmount = amount1
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = false
	} else {
		// amount is XRP/RLUSD, amount2 is token
		baseAmount = amount1
		tokenAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = (pool.GetAmountValue() != "" && pool.GetAmount2Currency() != RLUSDCurrency)
	}

	// XRP amounts from GetAmountValue() are in drops format, need conversion
	if isXRPBase && baseAmount > 0 {
		// CRITICAL: Convert drops to XRP for liquidity calculations
		baseAmount = baseAmount / 1000000.0
	}

	// Calculate liquidity in USD using live price
	var liquidityUSD decimal.Decimal
	if isXRPBase {
		// For XRP base, get live XRP price and calculate liquidity
		xrpPriceUSD, err := s.GetLiveXRPPrice()
		if err != nil {
			fmt.Printf("❌ Cannot calculate liquidity - live XRP price unavailable: %v\n", err)
			return false // Skip this pool if we can't get live price
		}
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2)) // Double the base amount for total liquidity
		// fmt.Printf("💰 XRP pool liquidity calc: %.2f XRP × $%.2f × 2 = $%.2f USD\n",
		// 	baseAmount, xrpPriceUSD.InexactFloat64(), liquidityUSD.InexactFloat64())
	} else {
		// For RLUSD base, RLUSD is approximately $1
		liquidityUSD = decimal.NewFromFloat(baseAmount).Mul(decimal.NewFromInt(2)) // Double the base amount for total liquidity
		// fmt.Printf("💰 RLUSD pool liquidity calc: %.2f RLUSD × $1.00 × 2 = $%.2f USD\n",
		// 	baseAmount, liquidityUSD.InexactFloat64())
	}

	// Check if liquidity meets minimum threshold
	if liquidityUSD.LessThan(decimal.NewFromFloat(minLiquidityUSD)) {
		// Pool filtered below threshold (logging reduced)
		return false
	}

	// Check if pool already exists in database to prevent duplicates
	poolExists, err := s.CheckPoolExistsInDB(pool.Result.Amm.Account)
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to check pool existence for %s: %v\n", pool.Result.Amm.Account, err)
		// Continue with storage in case of check error
	} else if poolExists {
		// fmt.Printf("⏭️  Pool %s already exists in database, skipping duplicate\n", pool.Result.Amm.Account)
		return false
	}

	// Also check by currency pair to catch different orderings
	pairExists, existingAccount, err := s.CheckPoolExistsByCurrencyPair("XRP", "", tokenCurrency, tokenIssuer)
	if err != nil {
		fmt.Printf("⚠️  Warning: failed to check currency pair existence: %v\n", err)
	} else if pairExists && existingAccount != pool.Result.Amm.Account {
		fmt.Printf("⏭️  Pool for currency pair %s/%s already exists as account %s, skipping duplicate account %s\n",
			"XRP", tokenCurrency, existingAccount, pool.Result.Amm.Account)
		return false
	}

	// Store in optimized xrpAmm_v2 database with enhanced constraints and validation
	err = s.storeAMMPoolOptimized(pool, baseAmount, tokenAmount, tokenCurrency, tokenIssuer, liquidityUSD, isXRPBase)

	if err != nil {
		log.Printf("❌ [AMM ERROR] Failed to store pool %s: %v", pool.Result.Amm.Account, err)
		return false
	}

	fmt.Printf("✅ Stored new pool %s: %s/%s - Liquidity: $%.2f USD\n",
		pool.Result.Amm.Account, "XRP", tokenCurrency, liquidityUSD.InexactFloat64())

	return true
}

// storeAMMPoolOptimized stores AMM pool data with enhanced validation and deduplication
func (s *AMMService) storeAMMPoolOptimized(pool models.AMMInfo, baseAmount, tokenAmount float64, tokenCurrency, tokenIssuer string, liquidityUSD decimal.Decimal, isXRPBase bool) error {
	// Additional validation before storage
	if pool.Result.Amm.Account == "" {
		return fmt.Errorf("cannot store pool with empty account")
	}

	if tokenCurrency == "" || tokenIssuer == "" {
		return fmt.Errorf("cannot store pool with incomplete token information")
	}

	if baseAmount <= 0 || tokenAmount <= 0 {
		return fmt.Errorf("cannot store pool with invalid amounts")
	}

	// Determine asset currencies and values for proper normalization
	var asset1Currency, asset1Issuer string
	var asset1Value, asset2Value float64

	if isXRPBase {
		asset1Currency = "XRP"
		asset1Issuer = ""
		_, _ = asset1Currency, asset1Issuer // Mark as used
		asset1Value = baseAmount
		asset2Value = tokenAmount
	} else {
		asset1Currency = tokenCurrency
		asset1Issuer = tokenIssuer
		asset1Value = tokenAmount
		asset2Value = baseAmount
	}

	// Use INSERT ... ON DUPLICATE KEY UPDATE for upsert behavior with proper constraints
	query := `
		INSERT INTO xrpAmm_normalized (
			account,
			asset1_amount,
			asset2_currency,
			asset2_issuer,
			asset2_amount,
			asset2_frozen,
			lp_token_currency,
			lp_token_issuer,
			lp_token_amount,
			trading_fee,
			liquidity_usd,
			created_at,
			last_updated
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			asset1_amount = VALUES(asset1_amount),
			asset2_currency = VALUES(asset2_currency),
			asset2_issuer = VALUES(asset2_issuer),
			asset2_amount = VALUES(asset2_amount),
			asset2_frozen = VALUES(asset2_frozen),
			lp_token_currency = VALUES(lp_token_currency),
			lp_token_issuer = VALUES(lp_token_issuer),
			lp_token_amount = VALUES(lp_token_amount),
			trading_fee = VALUES(trading_fee),
			liquidity_usd = VALUES(liquidity_usd),
			last_updated = NOW()
	`

	// Parse LP token data
	var lpTokenValue float64
	var lpTokenCurrency, lpTokenIssuer string

	if pool.Result.Amm.LpToken.Value != "" {
		lpTokenValue, _ = strconv.ParseFloat(pool.Result.Amm.LpToken.Value, 64)
	}
	lpTokenCurrency = pool.Result.Amm.LpToken.Currency
	lpTokenIssuer = pool.Result.Amm.LpToken.Issuer

	// Execute the optimized insert/update
	err := s.db.Exec(query,
		pool.Result.Amm.Account,      // account
		asset1Value,                  // asset1_amount
		tokenCurrency,                // asset2_currency
		tokenIssuer,                  // asset2_issuer
		asset2Value,                  // asset2_amount
		pool.Result.Amm.Asset2Frozen, // asset2_frozen
		lpTokenCurrency,              // lp_token_currency
		lpTokenIssuer,                // lp_token_issuer
		lpTokenValue,                 // lp_token_amount
		pool.Result.Amm.TradingFee,   // trading_fee
		liquidityUSD,                 // liquidity_usd
	).Error

	if err != nil {
		log.Printf("❌ Failed to store AMM pool %s: %v", pool.Result.Amm.Account, err)
		return err
	}

	log.Printf("✅ [AMM SUCCESS] Successfully stored/updated AMM pool %s with $%.2f USD liquidity", pool.Result.Amm.Account, liquidityUSD.InexactFloat64())
	return nil
}

// GetHeatmapData returns heatmap data for AMM pools
func (s *AMMService) GetHeatmapData(filters HeatmapFilters) ([]AMMHeatmapData, error) {
	// Build the base query using the correct table
	query := `
		SELECT 
			xa.account,
			xa.asset1_amount,
			xa.asset2_currency,
			xa.asset2_issuer,
			xa.asset2_amount,
			xa.trading_fee,
			xa.liquidity_usd,
			xa.created_at,
			xa.last_updated,
			xt1.name as token1_name,
			xt1.icon as token1_icon,
			xt1.marketcap as token1_marketcap,
			xt1.price as token1_price,
			xt1.volume24h as token1_volume24h,
			xt2.name as token2_name,
			xt2.icon as token2_icon,
			xt2.marketcap as token2_marketcap,
			xt2.price as token2_price,
			xt2.volume24h as token2_volume24h
		FROM xrpAmm_normalized xa
		LEFT JOIN xrpTokens xt1 ON xa.asset2_currency = xt1.currency AND xa.asset2_issuer = xt1.issuer
		LEFT JOIN xrpTokens xt2 ON xa.asset2_currency = xt2.currency AND xa.asset2_issuer = xt2.issuer
		WHERE xa.liquidity_usd > 0
	`

	var args []interface{}

	// Apply liquidity filters
	if filters.MinLiquidity > 0 {
		query += " AND xa.liquidity_usd >= ?"
		args = append(args, filters.MinLiquidity)
	}
	if filters.MaxLiquidity > 0 {
		query += " AND xa.liquidity_usd <= ?"
		args = append(args, filters.MaxLiquidity)
	}

	// Apply sorting
	switch filters.SortBy {
	case "LIQUIDITY":
		query += " ORDER BY xa.liquidity_usd"
	case "VOLUME":
		query += " ORDER BY COALESCE(xt1.volume24h, 0) + COALESCE(xt2.volume24h, 0)"
	case "MARKETCAP":
		query += " ORDER BY COALESCE(xt1.marketcap, 0) + COALESCE(xt2.marketcap, 0)"
	case "CREATED":
		query += " ORDER BY xa.created_at"
	default:
		query += " ORDER BY xa.liquidity_usd"
	}

	if filters.SortOrder == "ASC" {
		query += " ASC"
	} else {
		query += " DESC"
	}

	// Apply limit and offset
	if filters.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filters.Limit)
	}
	if filters.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filters.Offset)
	}

	// Execute query
	rows, err := s.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query AMM heatmap data: %w", err)
	}
	defer rows.Close()

	var heatmapData []AMMHeatmapData
	poolMap := make(map[string]*AMMHeatmapData)

	for rows.Next() {
		var (
			account, amount, amount2Currency, amount2Issuer, amount2Value         string
			tradingFee                                                            int
			liquidityUSD                                                          float64
			createdAt, lastUpdated                                                int64
			token1Name, token1Icon, token1Marketcap, token1Price, token1Volume24h string
			token2Name, token2Icon, token2Marketcap, token2Price, token2Volume24h string
		)

		err := rows.Scan(
			&account, &amount, &amount2Currency, &amount2Issuer, &amount2Value,
			&tradingFee, &liquidityUSD, &createdAt, &lastUpdated,
			&token1Name, &token1Icon, &token1Marketcap, &token1Price, &token1Volume24h,
			&token2Name, &token2Icon, &token2Marketcap, &token2Price, &token2Volume24h,
		)
		if err != nil {
			continue // Skip invalid rows
		}

		// Parse string values to float64 for calculations
		var token1MarketcapFloat, token1PriceFloat, token1Volume24hFloat float64
		var token2MarketcapFloat, token2PriceFloat, token2Volume24hFloat float64

		if token1Marketcap != "" {
			token1MarketcapFloat, _ = strconv.ParseFloat(token1Marketcap, 64)
		}
		if token1Price != "" {
			token1PriceFloat, _ = strconv.ParseFloat(token1Price, 64)
			_ = token1PriceFloat // Mark as used
		}
		if token1Volume24h != "" {
			token1Volume24hFloat, _ = strconv.ParseFloat(token1Volume24h, 64)
		}
		if token2Marketcap != "" {
			token2MarketcapFloat, _ = strconv.ParseFloat(token2Marketcap, 64)
		}
		if token2Price != "" {
			token2PriceFloat, _ = strconv.ParseFloat(token2Price, 64)
		}
		if token2Volume24h != "" {
			token2Volume24hFloat, _ = strconv.ParseFloat(token2Volume24h, 64)
		}

		// Create or update heatmap data entry
		key := fmt.Sprintf("%s_%s_%s_%s", "XRP", "", amount2Currency, amount2Issuer)

		if heatmapEntry, exists := poolMap[key]; exists {
			// Update existing entry with additional pool data
			ammInfo := models.AMMInfo{}

			// Set basic fields
			ammInfo.Result.Amm.Account = account

			// Set Amount using RawMessage
			amountData := map[string]string{
				"currency": "XRP",
				"issuer":   "",
				"value":    amount,
			}
			amountJSON, _ := json.Marshal(amountData)
			ammInfo.Result.Amm.Amount = amountJSON

			// Set Amount2 using RawMessage
			amount2Data := map[string]string{
				"currency": amount2Currency,
				"issuer":   amount2Issuer,
				"value":    amount2Value,
			}
			amount2JSON, _ := json.Marshal(amount2Data)
			ammInfo.Result.Amm.Amount2 = amount2JSON

			ammInfo.Result.Amm.TradingFee = tradingFee
			ammInfo.Result.LedgerCurrentIndex = 0
			ammInfo.Result.Validated = false
			heatmapEntry.Pools = append(heatmapEntry.Pools, ammInfo)
			heatmapEntry.TotalLiquidity += liquidityUSD
			poolMap[key] = heatmapEntry
		} else {
			// Create new heatmap entry
			token := models.XRPTokenData{
				Currency: amount2Currency,
				Issuer:   amount2Issuer,
				Meta: struct {
					Token struct {
						Description    string `json:"description"`
						SelfAssessment bool   `json:"self_assessment"`
						TrustLevel     int    `json:"trust_level"`
						Name           string `json:"name"`
						Icon           string `json:"icon"`
						Weblinks       []struct {
							Url   string `json:"url"`
							Type  string `json:"type"`
							Title string `json:"title"`
						} `json:"weblinks"`
					} `json:"token"`
					Issuer struct {
						Description string `json:"description"`
						Domain      string `json:"domain"`
						Followers   int    `json:"followers"`
						Icon        string `json:"icon"`
						Kyc         bool   `json:"kyc"`
						Name        string `json:"name"`
						TrustLevel  int    `json:"trust_level"`
						Weblinks    []struct {
							Url  string `json:"url"`
							Type string `json:"type"`
						} `json:"weblinks"`
					} `json:"issuer"`
				}{
					Token: struct {
						Description    string `json:"description"`
						SelfAssessment bool   `json:"self_assessment"`
						TrustLevel     int    `json:"trust_level"`
						Name           string `json:"name"`
						Icon           string `json:"icon"`
						Weblinks       []struct {
							Url   string `json:"url"`
							Type  string `json:"type"`
							Title string `json:"title"`
						} `json:"weblinks"`
					}{
						Name: token2Name,
						Icon: token2Icon,
					},
				},
				Metrics: struct {
					Trustlines   int     `json:"trustlines"`
					Holders      int     `json:"holders"`
					Supply       string  `json:"supply"`
					Marketcap    float64 `json:"marketCap,string"`
					Price        float64 `json:"price,string"`
					Volume24H    float64 `json:"volume24H,string"`
					Volume7D     string  `json:"volume_7d"`
					Exchanges24H string  `json:"exchanges_24h"`
					Exchanges7D  string  `json:"exchanges_7d"`
					Takers24H    string  `json:"takers_24h"`
					Takers7D     string  `json:"takers_7d"`
				}{
					Marketcap: token2MarketcapFloat,
					Price:     token2PriceFloat,
					Volume24H: token2Volume24hFloat,
				},
			}

			ammInfo := models.AMMInfo{}

			// Set basic fields
			ammInfo.Result.Amm.Account = account

			// Set Amount using RawMessage
			amountData := map[string]string{
				"currency": "XRP",
				"issuer":   "",
				"value":    amount,
			}
			amountJSON, _ := json.Marshal(amountData)
			ammInfo.Result.Amm.Amount = amountJSON

			// Set Amount2 using RawMessage
			amount2Data := map[string]string{
				"currency": amount2Currency,
				"issuer":   amount2Issuer,
				"value":    amount2Value,
			}
			amount2JSON, _ := json.Marshal(amount2Data)
			ammInfo.Result.Amm.Amount2 = amount2JSON

			ammInfo.Result.Amm.TradingFee = tradingFee
			ammInfo.Result.LedgerCurrentIndex = 0
			ammInfo.Result.Validated = false

			heatmapEntry := &AMMHeatmapData{
				Token:          token,
				Pools:          []models.AMMInfo{ammInfo},
				TotalLiquidity: liquidityUSD,
				Volume24h:      token1Volume24hFloat + token2Volume24hFloat,
				Volume7d:       0, // TODO: Calculate 7-day volume
				MarketCap:      token1MarketcapFloat + token2MarketcapFloat,
				PriceChange24h: 0, // TODO: Calculate price change
			}

			poolMap[key] = heatmapEntry
		}
	}

	// Convert map to slice
	for _, entry := range poolMap {
		heatmapData = append(heatmapData, *entry)
	}

	return heatmapData, nil
}

// UpdatePoolMetrics updates pool metrics
func (s *AMMService) UpdatePoolMetrics() error {
	// UpdatePoolMetrics called (debug logging reduced)

	// Rate limit check: Test XRPL connection before processing pools
	if s.xrpConnMgr != nil {
		testClient, err := s.xrpConnMgr.GetConnection()
		if err != nil {
			log.Printf("🛑 [AMM RATE LIMIT] Cannot get XRPL connection, skipping pool updates: %v", err)
			return fmt.Errorf("XRPL connection unavailable, skipping pool updates: %w", err)
		}

		// Test the connection with a simple ping
		testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Fix: use defer to ensure cancel is called
		done := make(chan error, 1)
		go func() {
			done <- testClient.Ping([]byte("PING"))
		}()

		select {
		case pingErr := <-done:
			if pingErr != nil {
				s.xrpConnMgr.ReturnConnection(testClient)
				log.Printf("🛑 [AMM RATE LIMIT] XRPL connection test failed, skipping pool updates: %v", pingErr)
				return fmt.Errorf("XRPL connection test failed, skipping pool updates: %w", pingErr)
			}
		case <-testCtx.Done():
			s.xrpConnMgr.ReturnConnection(testClient)
			log.Printf("🛑 [AMM RATE LIMIT] XRPL connection timeout, skipping pool updates")
			return fmt.Errorf("XRPL connection timeout, skipping pool updates")
		}

		s.xrpConnMgr.ReturnConnection(testClient)
	}

	// Get all existing pools from database
	poolRows, err := s.GetAllAMMPools()
	if err != nil {
		log.Printf("❌ [AMM ERROR] Failed to get existing pools: %v", err)
		return fmt.Errorf("failed to get existing pools: %w", err)
	}

	// Found pools in database (debug logging reduced)

	var updatedPools []models.AMMInfo
	successCount := 0
	errorCount := 0
	skippedCount := 0

	// Process each pool
	for _, poolRow := range poolRows {
		// Convert byte arrays to strings
		account := string(poolRow.Account)
		amount2Currency := string(poolRow.Asset2Currency)
		amount2Issuer := string(poolRow.Asset2Issuer)

		// Processing pool (debug logging reduced)

		// FIXED: Process ALL pools including XRP pairs - they are essential for pricing!
		// XRP pairs (amount2Currency="") are the PRIMARY source of token prices
		if amount2Currency == "" {
			// Processing XRP pair (debug logging removed)
		}

		// Get fresh pool data from XRPL with circuit breaker for rate limiting
		// Fetching fresh data (debug logging removed)
		poolInfo, err := s.CheckSingleAMMPool("XRP", "", amount2Currency, amount2Issuer)
		if err != nil {
			// Circuit breaker: If we're getting rate limited, trigger global backoff and stop
			if strings.Contains(err.Error(), "policy violation") ||
				strings.Contains(err.Error(), "IP limit") ||
				strings.Contains(err.Error(), "close 1008") ||
				strings.Contains(err.Error(), "close sent") ||
				strings.Contains(err.Error(), "WS read error") ||
				strings.Contains(err.Error(), "subscription error") {
				log.Printf("🛑 [AMM CIRCUIT BREAKER] XRPL rate limit detected, pausing requests for 15s")
				s.xrpConnMgr.TriggerRateLimitBackoff(15 * time.Second)
				break
			}

			// Skip logging for all errors to reduce log noise (error counters still track them)
			errorCount++
			continue
		}

		if poolInfo != nil {
			// Successfully got fresh data (debug logging removed)
			updatedPools = append(updatedPools, *poolInfo)
			successCount++
		} else {
			log.Printf("⚠️ [AMM WARNING] No pool info returned for: %s", account)
			errorCount++
		}
	}

	// Update summary (debug logging reduced)

	// Store updated pool data in database using COMPLETE StoreAMMPools method
	if len(updatedPools) > 0 {
		log.Printf("💾 [AMM INFO] Storing %d updated pools to database", len(updatedPools))
		if err := s.StoreAMMPoolsNormalized(updatedPools); err != nil {
			log.Printf("❌ [AMM ERROR] Failed to store updated pools: %v", err)
			return fmt.Errorf("failed to store updated pools: %w", err)
		}
		log.Printf("✅ [AMM SUCCESS] Successfully stored %d pools", len(updatedPools))
	} else {
		log.Printf("⚠️ [AMM WARNING] No pools to store - all updates failed or were skipped")
	}

	log.Printf("✅ [AMM SUCCESS] Pool metrics update complete: %d updated, %d errors, %d skipped", successCount, errorCount, skippedCount)
	return nil
}

// calculatePoolLiquidityUSD calculates the liquidity in USD for a given pool
func (s *AMMService) calculatePoolLiquidityUSD(pool models.AMMInfo) float64 {
	// SIMPLIFIED: Database schema is clear - no need to "determine" which is which
	// amount = XRP amount (always)
	// amount2 = Token amount (always)

	xrpAmountDrops, _ := strconv.ParseFloat(pool.GetAmountValue(), 64)
	// Token amount calculation (available for future use in more complex liquidity calculations)
	tokenAmount, _ := strconv.ParseFloat(pool.GetAmount2Value(), 64)
	_ = tokenAmount // Prevent unused variable warning

	// CRITICAL: Convert drops to XRP - GetAmountValue() returns drops format
	xrpAmount := xrpAmountDrops / 1000000.0

	// Calculate liquidity in USD using live XRP price
	xrpPriceUSD, err := s.GetLiveXRPPrice()
	if err != nil {
		fmt.Printf("❌ Cannot calculate liquidity - live XRP price unavailable: %v\n", err)
		return 0 // Return 0 liquidity if we can't get live price
	}

	// Liquidity = XRP amount * XRP price * 2 (both sides of pool)
	liquidityUSD, _ := decimal.NewFromFloat(xrpAmount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2)).Float64()

	// Liquidity calculation (debug logging reduced)

	return liquidityUSD
}

// CalculatePoolLiquidityUSD is a public wrapper for calculatePoolLiquidityUSD
func (s *AMMService) CalculatePoolLiquidityUSD(pool models.AMMInfo) float64 {
	return s.calculatePoolLiquidityUSD(pool)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetAllAMMPools retrieves all AMM pools from the database
// AMMPoolRowForUpdate is a subset of AMMPoolRow with only the fields needed for updates
type AMMPoolRowForUpdate struct {
	Account        []byte    `gorm:"column:account" json:"account"`
	Asset1Amount   []byte    `gorm:"column:asset1_amount" json:"asset1_amount"`
	Asset2Currency []byte    `gorm:"column:asset2_currency" json:"asset2_currency"`
	Asset2Issuer   []byte    `gorm:"column:asset2_issuer" json:"asset2_issuer"`
	Asset2Amount   []byte    `gorm:"column:asset2_amount" json:"asset2_amount"`
	Asset2Frozen   int       `gorm:"column:asset2_frozen" json:"asset2_frozen"`
	TradingFee     int       `gorm:"column:trading_fee" json:"trading_fee"`
	LiquidityUSD   float64   `gorm:"column:liquidity_usd" json:"liquidity_usd"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	LastUpdated    time.Time `gorm:"column:last_updated" json:"last_updated"`
}

func (s *AMMService) GetAllAMMPools() ([]AMMPoolRowForUpdate, error) {
	var poolRows []AMMPoolRowForUpdate
	err := s.db.Raw(`SELECT account, asset1_amount, asset2_currency, asset2_issuer, asset2_amount, asset2_frozen, trading_fee, liquidity_usd, created_at, last_updated FROM xrpAmm_normalized`).Scan(&poolRows).Error
	if err != nil {
		return nil, err
	}
	return poolRows, nil
}

// GetAMMPoolInfo retrieves information for a specific AMM pool by account
func (s *AMMService) GetAMMPoolInfo(account string) (*models.AMMInfo, error) {
	if account == "" {
		return nil, fmt.Errorf("account cannot be empty")
	}

	// Get connection from pool
	client, err := s.xrpConnMgr.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL connection: %v", err)
	}
	defer s.xrpConnMgr.ReturnConnection(client)

	// Use account_info to get AMM account details first
	accountRequest := map[string]interface{}{
		"command":      "account_info",
		"account":      account,
		"ledger_index": "validated",
	}

	if reqBytes, _ := json.Marshal(accountRequest); reqBytes != nil {
		log.Printf("🔍 [XRPL DEBUG] account_info request -> %s", string(reqBytes))
	}
	accountResponse, err := client.Request(accountRequest)
	if err != nil {
		s.xrpConnMgr.MarkConnectionUnhealthy(client)
		return nil, fmt.Errorf("XRPL account_info error: %v", err)
	}
	if respBytes, _ := json.Marshal(accountResponse); respBytes != nil {
		log.Printf("🔍 [XRPL DEBUG] account_info response <- %s", string(respBytes))
	}

	// Check for errors
	if status, ok := accountResponse["status"].(string); ok && status == "error" {
		if errorCode, ok := accountResponse["error"].(string); ok {
			if errorCode == "actNotFound" {
				return nil, fmt.Errorf("AMM pool account not found: %s", account)
			}
			return nil, fmt.Errorf("XRPL API error: %s", errorCode)
		}
		return nil, fmt.Errorf("XRPL API error: unknown error")
	}

	// Try to get AMM info using account objects
	objectsRequest := map[string]interface{}{
		"command":      "account_objects",
		"account":      account,
		"type":         "amm",
		"ledger_index": "validated",
	}

	if reqBytes, _ := json.Marshal(objectsRequest); reqBytes != nil {
		log.Printf("🔍 [XRPL DEBUG] account_objects request -> %s", string(reqBytes))
	}
	objectsResponse, err := client.Request(objectsRequest)
	if err != nil {
		s.xrpConnMgr.MarkConnectionUnhealthy(client)
		return nil, fmt.Errorf("XRPL account_objects error: %v", err)
	}
	if respBytes, _ := json.Marshal(objectsResponse); respBytes != nil {
		log.Printf("🔍 [XRPL DEBUG] account_objects response <- %s", string(respBytes))
	}

	// Parse the response to extract AMM information
	ammInfo := s.parseAMMAccountObjectsResponse(objectsResponse, account)
	if ammInfo == nil {
		return nil, fmt.Errorf("failed to parse AMM account objects for %s", account)
	}

	return ammInfo, nil
}

// parseAMMAccountObjectsResponse parses account objects response to extract AMM info
func (s *AMMService) parseAMMAccountObjectsResponse(response interface{}, account string) *models.AMMInfo {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil
	}

	result, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		return nil
	}

	accountObjects, ok := result["account_objects"].([]interface{})
	if !ok || len(accountObjects) == 0 {
		return nil
	}

	// Look for AMM object
	for _, obj := range accountObjects {
		objectMap, ok := obj.(map[string]interface{})
		if !ok {
			continue
		}

		ledgerEntryType, ok := objectMap["LedgerEntryType"].(string)
		if !ok || ledgerEntryType != "AMM" {
			continue
		}

		// Found AMM object, parse it
		ammInfo := &models.AMMInfo{}
		ammInfo.Result.Amm.Account = account

		// Parse asset1 (Amount)
		if amount, ok := objectMap["Amount"].(string); ok {
			// Set Amount using RawMessage
			amountData := map[string]string{
				"currency": "XRP",
				"issuer":   "",
				"value":    amount,
			}
			amountJSON, _ := json.Marshal(amountData)
			ammInfo.Result.Amm.Amount = amountJSON
		}

		// Parse asset2 (Amount2)
		if amount2, ok := objectMap["Amount2"].(map[string]interface{}); ok {
			// Convert Amount2 to RawMessage format
			amount2Data := make(map[string]string)
			if currency, ok := amount2["currency"].(string); ok {
				amount2Data["currency"] = currency
			}
			if issuer, ok := amount2["issuer"].(string); ok {
				amount2Data["issuer"] = issuer
			}
			if value, ok := amount2["value"].(string); ok {
				amount2Data["value"] = value
			}
			amount2JSON, _ := json.Marshal(amount2Data)
			ammInfo.Result.Amm.Amount2 = amount2JSON
		}

		// Parse LP Token
		if lpToken, ok := objectMap["LPTokenBalance"].(map[string]interface{}); ok {
			if currency, ok := lpToken["currency"].(string); ok {
				ammInfo.Result.Amm.LpToken.Currency = currency
			}
			if issuer, ok := lpToken["issuer"].(string); ok {
				ammInfo.Result.Amm.LpToken.Issuer = issuer
			}
			if value, ok := lpToken["value"].(string); ok {
				ammInfo.Result.Amm.LpToken.Value = value
			}
		}

		// Parse trading fee
		if tradingFee, ok := objectMap["TradingFee"].(float64); ok {
			ammInfo.Result.Amm.TradingFee = int(tradingFee)
		}

		return ammInfo
	}

	return nil
}

// GetAMMSwapQuote gets a real swap quote for AMM trading using XRPL path_find
func (s *AMMService) GetAMMSwapQuote(asset1, asset2 map[string]string, amount string) (*models.AMMInfo, error) {
	if asset1 == nil || asset2 == nil || amount == "" {
		return nil, fmt.Errorf("invalid input parameters")
	}

	// Validate amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}

	// Get connection from pool
	client, err := s.xrpConnMgr.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL connection: %v", err)
	}
	defer s.xrpConnMgr.ReturnConnection(client)

	// Build source and destination currency specs
	var sourceCurrency, destCurrency interface{}

	// Handle source currency (asset1)
	if asset1["currency"] == "XRP" {
		// XRP is represented as a string amount
		sourceCurrency = fmt.Sprintf("%.0f", amountFloat*1000000) // Convert to drops
	} else {
		sourceCurrency = map[string]string{
			"currency": asset1["currency"],
			"issuer":   asset1["issuer"],
			"value":    amount,
		}
	}

	// Handle destination currency (asset2)
	if asset2["currency"] == "XRP" {
		destCurrency = "XRP"
	} else {
		destCurrency = map[string]string{
			"currency": asset2["currency"],
			"issuer":   asset2["issuer"],
		}
	}

	// Use ripple_path_find to get swap quote
	pathRequest := map[string]interface{}{
		"command":             "ripple_path_find",
		"source_account":      "rSourceAccount123", // Placeholder - in real implementation this would be user's account
		"destination_account": "rDestAccount123",   // Placeholder - could be same as source for swaps
		"destination_amount":  destCurrency,
		"send_max":            sourceCurrency,
		"ledger_index":        "validated",
	}

	pathResponse, err := client.Request(pathRequest)
	if err != nil {
		s.xrpConnMgr.MarkConnectionUnhealthy(client)
		return nil, fmt.Errorf("XRPL path_find error: %v", err)
	}

	// Check for errors
	if status, ok := pathResponse["status"].(string); ok && status == "error" {
		if errorCode, ok := pathResponse["error"].(string); ok {
			return nil, fmt.Errorf("XRPL API error: %s", errorCode)
		}
		return nil, fmt.Errorf("XRPL API error: unknown error")
	}

	// Parse the path finding response into a swap quote
	swapQuote := s.parseSwapQuoteResponse(pathResponse, asset1, asset2, amount)
	if swapQuote == nil {
		return nil, fmt.Errorf("failed to parse swap quote response")
	}

	return swapQuote, nil
}

// parseSwapQuoteResponse converts path_find response to AMMInfo structure for compatibility
func (s *AMMService) parseSwapQuoteResponse(response interface{}, asset1, asset2 map[string]string, inputAmount string) *models.AMMInfo {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil
	}

	result, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		return nil
	}

	alternatives, ok := result["alternatives"].([]interface{})
	if !ok || len(alternatives) == 0 {
		return nil
	}

	// Take the first (best) alternative
	bestPath, ok := alternatives[0].(map[string]interface{})
	if !ok {
		return nil
	}

	// Extract destination amount (what we would receive)
	destAmount, ok := bestPath["destination_amount"]
	if !ok {
		return nil
	}

	// Calculate output amount
	var outputAmount float64
	if destAmountStr, ok := destAmount.(string); ok {
		// XRP amount already in XRP format
		if val, err := strconv.ParseFloat(destAmountStr, 64); err == nil {
			outputAmount = val // No conversion needed - already in XRP format
		}
	} else if destAmountMap, ok := destAmount.(map[string]interface{}); ok {
		// Token amount
		if value, ok := destAmountMap["value"].(string); ok {
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				outputAmount = val
			}
		}
	}

	// Create a mock AMMInfo response for compatibility with existing code
	inputAmountFloat, _ := strconv.ParseFloat(inputAmount, 64)

	// Calculate basic price impact (simplified)
	priceImpact := 0.01 // 1% default
	if outputAmount > 0 && inputAmountFloat > 0 {
		rate := outputAmount / inputAmountFloat
		if rate < 0.99 {
			priceImpact = 1.0 - rate
		}
	}

	ammInfo := &models.AMMInfo{}
	ammInfo.Result.Amm.Account = "swap_quote_simulation"

	// Set Amount using RawMessage
	amountData := map[string]string{
		"currency": "XRP",
		"issuer":   "",
		"value":    inputAmount,
	}
	amountJSON, _ := json.Marshal(amountData)
	ammInfo.Result.Amm.Amount = amountJSON

	// Set Amount2 using RawMessage
	amount2Data := map[string]string{
		"currency": asset2["currency"],
		"issuer":   asset2["issuer"],
		"value":    fmt.Sprintf("%.6f", outputAmount),
	}
	amount2JSON, _ := json.Marshal(amount2Data)
	ammInfo.Result.Amm.Amount2 = amount2JSON

	ammInfo.Result.Amm.TradingFee = int(priceImpact * 10000) // Store price impact as trading fee for now

	return ammInfo
}

// GetAMMLiquidityValue gets the liquidity value for an AMM pool with real calculation
func (s *AMMService) GetAMMLiquidityValue(account string) (float64, error) {
	// Get the actual pool info
	poolInfo, err := s.GetAMMPoolInfo(account)
	if err != nil {
		return 0, fmt.Errorf("failed to get pool info: %v", err)
	}

	if poolInfo == nil {
		return 0, fmt.Errorf("pool not found")
	}

	// Calculate liquidity using the existing method
	liquidityUSD := s.CalculatePoolLiquidityUSD(*poolInfo)

	return liquidityUSD, nil
}

// GetAMMTransactions gets transactions for an AMM pool
func (s *AMMService) GetAMMTransactions(account string, limit int) ([]interface{}, error) {
	// For now, return an empty slice since we don't have transaction data in the database yet
	// In the future, this would query a transactions table
	return []interface{}{}, nil
}

// computeLiquidityByXRPAmount calculates USD liquidity for AMM pools
// by taking the XRP amount, converting to USD, and doubling it
func (s *AMMService) computeLiquidityByXRPAmount() error {
	// Get current XRP price in USD from live price service
	xrpPriceUSD, err := s.GetLiveXRPPrice()
	if err != nil {
		return fmt.Errorf("CRITICAL: Cannot compute liquidity without live XRP price: %v", err)
	}

	fmt.Printf("Current XRP price: $%.6f USD (source: %s)\n", xrpPriceUSD.InexactFloat64(),
		func() string {
			if err != nil {
				return "fallback"
			} else {
				return "live AMM"
			}
		}())

	// Get all AMM pools from database
	var pools []struct {
		Account         string  `json:"account"`
		Asset1AmountStr string  `json:"asset1_amount_str"`
		LiquidityUSD    float64 `json:"liquidity_usd"`
	}

	err = s.db.Raw("SELECT account, asset1_amount as asset1_amount_str, liquidity_usd FROM xrpAmm_normalized").Scan(&pools).Error
	if err != nil {
		return fmt.Errorf("failed to query AMM pools: %w", err)
	}

	fmt.Printf("Found %d AMM pools to process\n", len(pools))

	// Process pools in batches for better performance
	batchSize := 100
	totalBatches := (len(pools) + batchSize - 1) / batchSize
	updatedCount := 0

	for batch := 0; batch < totalBatches; batch++ {
		start := batch * batchSize
		end := start + batchSize
		if end > len(pools) {
			end = len(pools)
		}

		batchPools := pools[start:end]

		// Build batch update query
		var updateQueries []string
		var args []interface{}

		for _, pool := range batchPools {
			// CRITICAL: Parse XRP amount from database (stored in drops format)
			xrpAmountDrops, err := decimal.NewFromString(pool.Asset1AmountStr)
			if err != nil {
				continue // Skip invalid amounts
			}

			// Convert drops to XRP (CRITICAL FIX)
			xrpAmount := xrpAmountDrops.Div(decimal.NewFromFloat(1000000.0))

			// Calculate USD liquidity: (XRP amount * XRP price * 2) with proper conversion
			liquidityUSD := xrpAmount.Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2))

			liquidityFloat, _ := liquidityUSD.Float64()
			updateQueries = append(updateQueries, "WHEN ? THEN ?")
			args = append(args, pool.Account, liquidityFloat)
		}

		if len(updateQueries) > 0 {
			// Build the CASE statement for batch update
			caseStatement := strings.Join(updateQueries, " ")
			query := fmt.Sprintf(`
				UPDATE xrpAmm_normalized 
				SET liquidity_usd = CASE account %s END,
				    last_updated = ?
				WHERE account IN (?)
			`, caseStatement)

			// Add the timestamp and account list to args
			args = append(args, time.Now().Unix())
			accountList := make([]string, len(batchPools))
			for i, pool := range batchPools {
				accountList[i] = pool.Account
			}
			args = append(args, accountList)

			// Execute batch update
			err = s.db.Exec(query, args...).Error
			if err != nil {
				fmt.Printf("Error updating batch %d: %v\n", batch+1, err)
				continue
			}

			updatedCount += len(batchPools)
			fmt.Printf("Updated batch %d/%d (%d pools)...\n", batch+1, totalBatches, len(batchPools))
		}
	}

	log.Printf("✅ [AMM SUCCESS] Updated liquidity for %d pools", updatedCount)
	return nil
}

// GetTopAMMPoolsByLiquidity returns the top X AMM pools by liquidity_usd, filtering out duplicates
func (s *AMMService) GetTopAMMPoolsByLiquidity(limit int) ([]AMMPoolRow, error) {
	var poolRows []AMMPoolRow
	err := s.db.Raw(`
		SELECT DISTINCT * FROM xrpAmm
		WHERE liquidity_usd > 0
		ORDER BY liquidity_usd DESC
		LIMIT ?
	`, limit).Scan(&poolRows).Error
	if err != nil {
		return nil, err
	}
	return poolRows, nil
}

// StoreAMMPools implements the missing method to store multiple AMM pools using the same pattern as checkAndStorePoolWithLiquidity
func (s *AMMService) StoreAMMPools(pools []models.AMMInfo) error {
	// StoreAMMPools called (debug logging reduced)
	if len(pools) == 0 {
		return nil
	}

	// Debug: Log details about the first pool to understand structure (commented out to reduce noise)
	// if len(pools) > 0 {
	// 	firstPool := pools[0]
	// 	log.Printf("First pool analysis: Account=%s, Amount=%s", firstPool.Result.Amm.Account, firstPool.GetAmountValue())
	// }

	log.Printf("⚠️ [DEPRECATED] StoreAMMPools() writes to old xrpAmm table - use StoreAMMPoolsNormalized() instead!")
	log.Printf("📦 [STORE AMM POOLS] Storing %d AMM pools to DEPRECATED xrpAmm table...", len(pools))

	// Process pools in batches for better performance
	batchSize := 50
	totalBatches := (len(pools) + batchSize - 1) / batchSize
	storedCount := 0
	errorCount := 0

	for batch := 0; batch < totalBatches; batch++ {
		start := batch * batchSize
		end := start + batchSize
		if end > len(pools) {
			end = len(pools)
		}

		batchPools := pools[start:end]
		log.Printf("📋 [STORE AMM POOLS] Processing batch %d/%d (%d pools)...", batch+1, totalBatches, len(batchPools))

		// Process each pool in the batch
		for _, pool := range batchPools {
			// Processing pool (debug logging removed)

			if err := s.storeIndividualAMMPool(pool); err != nil {
				log.Printf("❌ [AMM ERROR] Failed to store pool %s: %v", pool.Result.Amm.Account, err)
				log.Printf("⚠️ [STORE AMM POOLS] Failed to store pool %s: %v", pool.Result.Amm.Account, err)
				errorCount++
			} else {
				// Successfully stored pool (debug logging removed)
				storedCount++
			}
		}
	}

	log.Printf("✅ [STORE AMM POOLS] Completed: %d stored, %d errors", storedCount, errorCount)

	if errorCount > 0 && storedCount == 0 {
		return fmt.Errorf("failed to store any AMM pools: %d errors", errorCount)
	}

	return nil
}

// storeIndividualAMMPool stores a single AMM pool using the same logic as checkAndStorePoolWithLiquidity
func (s *AMMService) storeIndividualAMMPool(pool models.AMMInfo) error {
	// Processing pool (debug logging reduced)

	// Amount parsing (errors handled)
	amount1Drops, _ := strconv.ParseFloat(pool.GetAmountValue(), 64)
	amount2, _ := strconv.ParseFloat(pool.GetAmount2Value(), 64)

	// Convert XRP drops to XRP for amount1 (from GetAmountValue)
	amount1 := amount1Drops / 1000000.0

	// Raw amounts processing (debug logging removed)

	// Determine which is XRP/RLUSD and which is the token (same logic as working code)
	var baseAmount, tokenAmount float64
	var tokenCurrency, tokenIssuer string
	var isXRPBase bool

	if pool.GetAmount2Currency() == "XRP" {
		// amount is token, amount2 is XRP
		tokenAmount = amount1
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = true
	} else if pool.GetAmount2Currency() == RLUSDCurrency && pool.GetAmount2Issuer() == RLUSDIssuer {
		// amount is token, amount2 is RLUSD
		tokenAmount = amount1
		baseAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = false
	} else {
		// amount is XRP/RLUSD, amount2 is token
		baseAmount = amount1
		tokenAmount = amount2
		tokenCurrency = pool.GetAmount2Currency()
		tokenIssuer = pool.GetAmount2Issuer()
		isXRPBase = (pool.GetAmountValue() != "" && pool.GetAmount2Currency() != RLUSDCurrency)
	}

	// XRP amounts from GetAmountValue() are in drops format - conversion needed
	if isXRPBase && baseAmount > 0 {
		// XRP amount converted from drops (debug logging removed)
		// Note: Drops conversion applied above (amount1Drops / 1000000.0)
	}

	// Final values for storage (debug logging removed)
	fmt.Printf("   - Is XRP Base: %t\n", isXRPBase)

	// Store in database using EXACT same SQL as working code
	// Executing database INSERT (debug logging removed)

	err := s.db.Exec(`
		INSERT INTO xrpAmm_normalized (
			account, asset1_amount, asset2_currency, asset2_issuer, asset2_amount, 
			trading_fee, liquidity_usd, asset1_value_usd, asset2_value_usd,
			created_at, last_updated
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
		asset1_amount = VALUES(asset1_amount),
		asset2_currency = VALUES(asset2_currency),
		asset2_issuer = VALUES(asset2_issuer),
		asset2_amount = VALUES(asset2_amount),
		trading_fee = VALUES(trading_fee),
		liquidity_usd = VALUES(liquidity_usd),
		asset1_value_usd = VALUES(asset1_value_usd),
		asset2_value_usd = VALUES(asset2_value_usd),
		last_updated = NOW(),
		last_liquidity_update = NOW()
	`, pool.Result.Amm.Account, baseAmount, tokenCurrency, tokenIssuer, tokenAmount,
		pool.Result.Amm.TradingFee, 0.0, 0.0, 0.0).Error

	if err != nil {
		log.Printf("❌ [AMM ERROR] Database error for pool %s: %v", pool.Result.Amm.Account, err)
		log.Printf("🔍 [AMM ERROR] Error type: %T", err)
		return fmt.Errorf("failed to store AMM pool %s: %w", pool.Result.Amm.Account, err)
	}

	// Successfully stored pool (debug logging removed)
	return nil
}

// StoreAMMPoolsNormalized stores AMM pools in the normalized database schema with proper deduplication
func (s *AMMService) StoreAMMPoolsNormalized(pools []models.AMMInfo) error {
	fmt.Printf("🚀 [AMM NORMALIZED] Starting normalized storage for %d pools\n", len(pools))

	if len(pools) == 0 {
		fmt.Printf("⚠️ [AMM NORMALIZED] No pools provided - returning early\n")
		return nil
	}

	// Track statistics
	storedCount := 0
	errorCount := 0
	duplicateCount := 0

	// Process pools individually with proper deduplication
	for _, pool := range pools {
		// Processing pool (debug logging reduced)

		// Check if pool already exists to avoid unnecessary processing
		exists, err := s.poolExistsInNormalizedTable(pool.Result.Amm.Account)
		if err != nil {
			log.Printf("❌ [AMM ERROR] Error checking pool existence: %v", err)
			errorCount++
			continue
		}

		if exists {
			// Pool already exists, updating (debug logging removed)
		}

		// Store or update the pool
		if err := s.storeNormalizedAMMPool(pool); err != nil {
			if isDuplicateAssetPairError(err) {
				fmt.Printf("⚠️ [AMM NORMALIZED] Duplicate asset pair for pool %s (this is expected)\n", pool.Result.Amm.Account)
				duplicateCount++
			} else {
				log.Printf("❌ [AMM ERROR] Failed to store pool %s: %v", pool.Result.Amm.Account, err)
				errorCount++
			}
		} else {
			// Pool stored/updated successfully (debug logging removed)
			storedCount++
		}
	}

	fmt.Printf("📊 [AMM NORMALIZED] Storage completed: %d stored/updated, %d duplicates skipped, %d errors\n",
		storedCount, duplicateCount, errorCount)

	// Consider it successful if we stored at least some pools
	if storedCount > 0 || duplicateCount > 0 {
		return nil
	}

	if errorCount > 0 {
		return fmt.Errorf("failed to store any AMM pools: %d errors", errorCount)
	}

	return nil
}

// storeNormalizedAMMPool stores a single AMM pool in the normalized table structure
func (s *AMMService) storeNormalizedAMMPool(pool models.AMMInfo) error {
	// Parse amounts and handle drops conversion
	amount1Drops, err1 := strconv.ParseFloat(pool.GetAmountValue(), 64)
	amount2, err2 := strconv.ParseFloat(pool.GetAmount2Value(), 64)

	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to parse amounts: amount1=%v, amount2=%v", err1, err2)
	}

	// Convert XRP drops to XRP (GetAmountValue returns drops)
	amount1 := amount1Drops / 1000000.0

	// Normalize the asset pair information
	asset1Currency, asset1Issuer, asset1Amount := s.normalizeAsset("XRP", "", amount1)
	asset2Currency, asset2Issuer, asset2Amount := s.normalizeAsset(
		pool.GetAmount2Currency(),
		pool.GetAmount2Issuer(),
		amount2,
	)

	// XRP amounts converted from drops to XRP above
	// GetAmountValue() returns drops format, conversion applied before normalization

	// Parse LP token information
	lpTokenAmount := 0.0
	if pool.Result.Amm.LpToken.Value != "" {
		if val, err := strconv.ParseFloat(pool.Result.Amm.LpToken.Value, 64); err == nil {
			lpTokenAmount = val
		}
	}

	// Calculate basic liquidity USD (will be improved with real price feeds later)
	liquidityUSD := s.calculateBasicLiquidityUSD(asset1Currency, asset1Amount, asset2Currency, asset2Amount)

	// Generate hex formats for XRPL API compatibility
	asset1CurrencyHex := s.convertCurrencyToXRPLHex(asset1Currency)
	asset2CurrencyHex := s.convertCurrencyToXRPLHex(asset2Currency)

	// Insert or update using the normalized schema with dual currency storage
	err := s.db.Exec(`
		INSERT INTO xrpAmm_normalized (
			account, 
			asset1_currency, asset1_currency_hex, asset1_issuer, asset1_amount,
			asset2_currency, asset2_currency_hex, asset2_issuer, asset2_amount,
			lp_token_currency, lp_token_issuer, lp_token_amount,
			trading_fee, asset2_frozen, liquidity_usd,
			last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			asset1_amount = VALUES(asset1_amount),
			asset1_currency_hex = VALUES(asset1_currency_hex),
			asset2_amount = VALUES(asset2_amount),
			asset2_currency_hex = VALUES(asset2_currency_hex),
			lp_token_amount = VALUES(lp_token_amount),
			trading_fee = VALUES(trading_fee),
			asset2_frozen = VALUES(asset2_frozen),
			liquidity_usd = VALUES(liquidity_usd),
			last_updated = NOW()
	`,
		pool.Result.Amm.Account,
		asset1Currency, asset1CurrencyHex, asset1Issuer, asset1Amount,
		asset2Currency, asset2CurrencyHex, asset2Issuer, asset2Amount,
		pool.Result.Amm.LpToken.Currency, pool.Result.Amm.LpToken.Issuer, lpTokenAmount,
		pool.Result.Amm.TradingFee, pool.Result.Amm.Asset2Frozen, liquidityUSD,
	).Error

	return err
}

// normalizeAsset normalizes asset information for consistent storage
func (s *AMMService) normalizeAsset(currency, issuer string, amount float64) (string, string, float64) {
	// Normalize XRP
	if currency == "" || currency == "XRP" {
		return "XRP", "", amount
	}

	// Normalize currency codes (handle hex currency codes)
	normalizedCurrency := s.normalizeCurrencyCode(currency)

	// Normalize issuer
	normalizedIssuer := strings.TrimSpace(issuer)

	return normalizedCurrency, normalizedIssuer, amount
}

// normalizeCurrencyCode converts hex currency codes to readable format where possible
func (s *AMMService) normalizeCurrencyCode(currency string) string {
	// Trim whitespace
	currency = strings.TrimSpace(currency)

	// If it's a 40-character hex string, try to decode it
	if len(currency) == 40 {
		// Remove trailing zeros
		currency = strings.TrimRight(currency, "0")

		// Try to decode as hex
		if decoded, err := hex.DecodeString(currency); err == nil {
			// Convert to string and remove null bytes
			readable := string(bytes.TrimRight(decoded, "\x00"))
			if len(readable) >= 3 && len(readable) <= 20 && isValidCurrencyCode(readable) {
				return strings.ToUpper(readable)
			}
		}
	}

	// Return as-is if not hex or couldn't decode
	return strings.ToUpper(currency)
}

// isValidCurrencyCode checks if a string looks like a valid currency code
func isValidCurrencyCode(code string) bool {
	// Basic validation: alphanumeric characters, reasonable length
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	for _, r := range code {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// calculateBasicLiquidityUSD provides a basic USD liquidity calculation
func (s *AMMService) calculateBasicLiquidityUSD(asset1Currency string, asset1Amount float64, asset2Currency string, asset2Amount float64) float64 {
	// Get live XRP price directly from XRPL ledger
	xrpPriceUSD, err := s.GetLiveXRPPrice()
	if err != nil {
		fmt.Printf("❌ Cannot get live XRP price for liquidity calculation: %v\n", err)
		return 0 // Return 0 if we can't get live price
	}
	rlusdPriceUSD := 1.0 // RLUSD is pegged to USD

	var totalUSD decimal.Decimal

	if asset1Currency == "XRP" {
		totalUSD = totalUSD.Add(decimal.NewFromFloat(asset1Amount).Mul(xrpPriceUSD))
	} else if asset1Currency == "RLUSD" {
		totalUSD = totalUSD.Add(decimal.NewFromFloat(asset1Amount).Mul(decimal.NewFromFloat(rlusdPriceUSD)))
	}

	if asset2Currency == "XRP" {
		totalUSD = totalUSD.Add(decimal.NewFromFloat(asset2Amount).Mul(xrpPriceUSD))
	} else if asset2Currency == "RLUSD" {
		totalUSD = totalUSD.Add(decimal.NewFromFloat(asset2Amount).Mul(decimal.NewFromFloat(rlusdPriceUSD)))
	}

	// If neither asset is XRP or RLUSD, estimate based on XRP amount
	if totalUSD.IsZero() && (asset1Currency == "XRP" || asset2Currency == "XRP") {
		if asset1Currency == "XRP" {
			totalUSD = decimal.NewFromFloat(asset1Amount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2)) // Assume both sides roughly equal
		} else if asset2Currency == "XRP" {
			totalUSD = decimal.NewFromFloat(asset2Amount).Mul(xrpPriceUSD).Mul(decimal.NewFromInt(2))
		}
	}

	return totalUSD.InexactFloat64()
}

// getRealXRPPrice is DEPRECATED - tries to read from dropped xrpAmm table
// Use GetLiveXRPPrice() instead for live XRPL ledger data
func (s *AMMService) getRealXRPPrice() float64 {
	// First try to get from the normalized table
	var xrpAmount, rlusdAmount float64
	err := s.db.Raw(`
		SELECT asset1_amount, asset2_amount 
		FROM xrpAmm_normalized 
		WHERE (asset1_currency = 'XRP' AND asset2_currency = 'RLUSD')
		   OR (asset1_currency = 'RLUSD' AND asset2_currency = 'XRP')
		ORDER BY last_updated DESC
		LIMIT 1
	`).Row().Scan(&xrpAmount, &rlusdAmount)

	if err == nil && xrpAmount > 0 && rlusdAmount > 0 {
		// Calculate price based on the pool ratio
		price := rlusdAmount / xrpAmount // RLUSD per XRP (since RLUSD ≈ $1)
		if price > 0 {                   // Only check for positive price
			return price
		}
	}

	// DEPRECATED: Old fallback to dropped xrpAmm table - this will fail
	// The xrpAmm table has been dropped - should only use live XRPL data
	fmt.Printf("❌ CRITICAL: Attempted to read from dropped xrpAmm table - using live XRPL price instead\n")
	// No hardcoded fallbacks - return 0 if no live data available
	fmt.Printf("⚠️  Warning: No live XRP price data available from any AMM pool\n")
	return 0
}

// poolExistsInNormalizedTable checks if a pool already exists in the normalized table
func (s *AMMService) poolExistsInNormalizedTable(account string) (bool, error) {
	var count int64
	err := s.db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized WHERE account = ?", account).Scan(&count).Error
	return count > 0, err
}

// isDuplicateAssetPairError checks if an error is due to duplicate asset pair constraint
func isDuplicateAssetPairError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "unique_asset_pair") ||
		strings.Contains(errStr, "duplicate entry") && strings.Contains(errStr, "unique_asset_pair")
}

// GetNormalizedAMMPools retrieves AMM pools from the normalized table
func (s *AMMService) GetNormalizedAMMPools(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 100
	}

	var pools []map[string]interface{}
	err := s.db.Raw(`
		SELECT 
			account,
			asset1_currency,
			asset1_issuer,
			asset1_amount,
			asset2_currency, 
			asset2_issuer,
			asset2_amount,
			lp_token_currency,
			lp_token_issuer,
			lp_token_amount,
			trading_fee,
			asset2_frozen,
			liquidity_usd,
			volume_24h,
			volume_7d,
			last_updated,
			created_at
		FROM xrpAmm_normalized 
		ORDER BY liquidity_usd DESC, last_updated DESC
		LIMIT ?
	`, limit).Scan(&pools).Error

	return pools, err
}

// DeduplicateExistingData removes duplicates from the old table and migrates to normalized structure
func (s *AMMService) DeduplicateExistingData() error {
	fmt.Printf("🧹 [DEDUPLICATION] Starting deduplication process...\n")

	// First, get all pools from the old table
	var oldPools []map[string]interface{}
	err := s.db.Raw(`
		SELECT DISTINCT 
			account, asset1_amount, asset2_currency, asset2_issuer, asset2_amount, 
			trading_fee, asset2_frozen, lp_token_currency, lp_token_issuer, lp_token_amount,
			COALESCE(liquidity_usd, 0) as liquidity_usd,
			COALESCE(created_at, UNIX_TIMESTAMP()) as created_at,
			COALESCE(last_updated, UNIX_TIMESTAMP()) as last_updated
		FROM xrpAmm_normalized 
		WHERE account IS NOT NULL AND account != ''
		ORDER BY account
	`).Scan(&oldPools).Error

	if err != nil {
		return fmt.Errorf("failed to read old pools: %v", err)
	}

	fmt.Printf("📊 [DEDUPLICATION] Found %d pools in old table\n", len(oldPools))

	// Group by account to find duplicates
	poolGroups := make(map[string][]map[string]interface{})
	for _, pool := range oldPools {
		account := pool["account"].(string)
		poolGroups[account] = append(poolGroups[account], pool)
	}

	duplicateCount := 0
	processedCount := 0

	// Process each group
	for account, pools := range poolGroups {
		if len(pools) > 1 {
			// Found duplicates (debug logging reduced)
			duplicateCount += len(pools) - 1
		}

		// Take the most recent pool (or first if timestamps are the same)
		selectedPool := pools[0]
		for _, pool := range pools[1:] {
			if lastUpdated, ok := pool["last_updated"].(int64); ok {
				if selectedLastUpdated, ok := selectedPool["last_updated"].(int64); ok {
					if lastUpdated > selectedLastUpdated {
						selectedPool = pool
					}
				}
			}
		}

		// Convert to normalized format and store
		if err := s.convertAndStoreNormalizedPool(selectedPool); err != nil {
			fmt.Printf("❌ [DEDUPLICATION] Failed to convert pool %s: %v\n", account, err)
		} else {
			processedCount++
		}
	}

	fmt.Printf("✅ [DEDUPLICATION] Completed: %d pools processed, %d duplicates removed\n", processedCount, duplicateCount)
	return nil
}

// convertAndStoreNormalizedPool converts an old pool record to normalized format
func (s *AMMService) convertAndStoreNormalizedPool(oldPool map[string]interface{}) error {
	account := oldPool["account"].(string)

	// Parse amounts
	amount1 := s.parseFloatFromInterface(oldPool["asset1_amount"])
	amount2 := s.parseFloatFromInterface(oldPool["asset2_amount"])

	// Get currency information
	amount2Currency := s.parseStringFromInterface(oldPool["asset2_currency"])
	amount2Issuer := s.parseStringFromInterface(oldPool["asset2_issuer"])

	// Normalize assets (old format: amount is usually XRP, amount2 is token)
	asset1Currency, asset1Issuer, asset1Amount := s.normalizeAsset("XRP", "", amount1)
	asset2Currency, asset2Issuer, asset2Amount := s.normalizeAsset(amount2Currency, amount2Issuer, amount2)

	// XRP amounts from GetAmountValue() converted from drops to XRP above

	// Get other fields
	tradingFee := s.parseIntFromInterface(oldPool["trading_fee"])
	asset2Frozen := s.parseIntFromInterface(oldPool["asset2_frozen"]) == 1
	liquidityUSD := s.parseFloatFromInterface(oldPool["liquidity_usd"])

	lpTokenCurrency := s.parseStringFromInterface(oldPool["lp_token_currency"])
	lpTokenIssuer := s.parseStringFromInterface(oldPool["lp_token_issuer"])
	lpTokenAmount := s.parseFloatFromInterface(oldPool["lp_token_amount"])

	// Store in normalized table
	err := s.db.Exec(`
		INSERT INTO xrpAmm_normalized (
			account, 
			asset1_currency, asset1_issuer, asset1_amount,
			asset2_currency, asset2_issuer, asset2_amount,
			lp_token_currency, lp_token_issuer, lp_token_amount,
			trading_fee, asset2_frozen, liquidity_usd,
			last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			asset1_amount = VALUES(asset1_amount),
			asset2_amount = VALUES(asset2_amount),
			lp_token_amount = VALUES(lp_token_amount),
			trading_fee = VALUES(trading_fee),
			asset2_frozen = VALUES(asset2_frozen),
			liquidity_usd = VALUES(liquidity_usd),
			last_updated = NOW()
	`,
		account,
		asset1Currency, asset1Issuer, asset1Amount,
		asset2Currency, asset2Issuer, asset2Amount,
		lpTokenCurrency, lpTokenIssuer, lpTokenAmount,
		tradingFee, asset2Frozen, liquidityUSD,
	).Error

	return err
}

// Helper functions for type conversion
func (s *AMMService) parseFloatFromInterface(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	case int64:
		return float64(v)
	case int:
		return float64(v)
	}
	return 0
}

func (s *AMMService) parseStringFromInterface(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

func (s *AMMService) parseIntFromInterface(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetLiveXRPPrice gets the current XRP price from LIVE XRPL ledger data
func (s *AMMService) GetLiveXRPPrice() (decimal.Decimal, error) {
	// Use the price service to get live data directly from XRPL ledger
	priceService := NewPriceServiceWithDBAndConn(s.db, s.xrpConnMgr)

	// Get live XRP price from XRPL ledger (no database dependency)
	return priceService.GetXRPLPrice()
}

// GetLiveXRPPriceWithFallback is DEPRECATED - use GetLiveXRPPrice() with proper error handling
// This method has been removed to eliminate hardcoded price fallbacks

// CheckPoolExistsInDB checks if an AMM pool already exists in the database
func (s *AMMService) CheckPoolExistsInDB(account string) (bool, error) {
	if account == "" {
		return false, nil
	}

	var count int64
	err := s.db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized WHERE account = ?", account).Scan(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check pool existence for account %s: %w", account, err)
	}

	return count > 0, nil
}

// CheckPoolExistsInDBNormalized checks if an AMM pool already exists in the normalized database
func (s *AMMService) CheckPoolExistsInDBNormalized(account string) (bool, error) {
	if account == "" {
		return false, nil
	}

	var count int64
	err := s.db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized WHERE account = ?", account).Scan(&count).Error
	if err != nil {
		// If normalized table doesn't exist, return false (not an error)
		if strings.Contains(err.Error(), "doesn't exist") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check pool existence in normalized table for account %s: %w", account, err)
	}

	return count > 0, nil
}

// CheckPoolExistsByCurrencyPair checks if a pool exists for a specific currency pair
func (s *AMMService) CheckPoolExistsByCurrencyPair(currency1, issuer1, currency2, issuer2 string) (bool, string, error) {
	// Check both possible orderings of the currency pair
	var account string
	var count int64

	// Query 1: currency1/currency2 order
	query1 := `
		SELECT account, COUNT(*) as cnt
		FROM xrpAmm_normalized 
		WHERE (
			(asset2_currency = ? AND asset2_issuer = ? AND 
			 ((? = 'XRP' AND asset1_amount REGEXP '^[0-9]+$') OR (asset2_currency = ? AND asset2_issuer = ?)))
			OR
			(? = 'XRP' AND asset1_amount REGEXP '^[0-9]+$' AND asset2_currency = ? AND asset2_issuer = ?)
		)
		GROUP BY account
		ORDER BY cnt DESC
		LIMIT 1
	`

	row := s.db.Raw(query1, currency1, issuer1, currency2, currency2, issuer2, currency2, currency1, issuer1).Row()

	err := row.Scan(&account, &count)
	if err == nil && count > 0 {
		return true, account, nil
	}

	// Query 2: currency2/currency1 order (reversed)
	query2 := `
		SELECT account, COUNT(*) as cnt
		FROM xrpAmm_normalized 
		WHERE (
			(asset2_currency = ? AND asset2_issuer = ? AND 
			 ((? = 'XRP' AND asset1_amount REGEXP '^[0-9]+$') OR (asset2_currency = ? AND asset2_issuer = ?)))
			OR
			(? = 'XRP' AND asset1_amount REGEXP '^[0-9]+$' AND asset2_currency = ? AND asset2_issuer = ?)
		)
		GROUP BY account
		ORDER BY cnt DESC
		LIMIT 1
	`

	row2 := s.db.Raw(query2, currency2, issuer2, currency1, currency1, issuer1, currency1, currency2, issuer2).Row()

	err2 := row2.Scan(&account, &count)
	if err2 == nil && count > 0 {
		return true, account, nil
	}

	return false, "", nil
}
