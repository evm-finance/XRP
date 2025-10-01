package xrp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	xrpl "github.com/xrpscan/xrpl-go"
)

type qcResult struct {
	price decimal.Decimal
	err   error
}

// getXRPLPriceFromQuantifyCryptoThrottled ensures at most one QuantifyCrypto request every 10 seconds
func (p *PriceService) getXRPLPriceFromQuantifyCryptoThrottled() (decimal.Decimal, error) {
	now := time.Now()

	p.qcMutex.Lock()
	// If last successful request is within 10s, return cached last result immediately
	if !p.qcLastRequestAt.IsZero() && now.Sub(p.qcLastRequestAt) < 10*time.Second {
		price, err := p.qcLastPrice, p.qcLastErr
		p.qcMutex.Unlock()
		return price, err
	}

	// If a request is already in-flight, wait on a channel
	if p.qcInFlight {
		ch := make(chan qcResult, 1)
		p.qcWaiters = append(p.qcWaiters, ch)
		p.qcMutex.Unlock()
		res := <-ch
		return res.price, res.err
	}

	// Mark in-flight and proceed to make a single request
	p.qcInFlight = true
	p.qcMutex.Unlock()

	// Perform the actual request outside the lock
	// Log removed to reduce log spam
	price, err := p.getXRPLPriceFromQuantifyCrypto()

	// Update last result and notify any waiters
	p.qcMutex.Lock()
	p.qcLastRequestAt = time.Now()
	p.qcLastPrice = price
	p.qcLastErr = err
	waiters := p.qcWaiters
	p.qcWaiters = nil
	p.qcInFlight = false
	p.qcMutex.Unlock()

	// Notify waiters
	if len(waiters) > 0 {
		res := qcResult{price: price, err: err}
		for _, ch := range waiters {
			ch <- res
			close(ch)
		}
	}

	return price, err
}

// PriceService manages token pricing
type PriceService struct {
	cache             map[string]PriceCache
	db                *gorm.DB
	connectionManager *ConnectionManager
	cacheMutex        sync.RWMutex // EMERGENCY: Add mutex for thread safety
	// Per-batch XRP price override: when active, GetXRPLPrice returns this value
	batchMutex    sync.RWMutex
	batchActive   bool
	batchXRPPrice decimal.Decimal
	// QuantifyCrypto throttling: ensure <=1 request per 10s and dedupe concurrent callers
	qcMutex         sync.Mutex
	qcLastRequestAt time.Time
	qcInFlight      bool
	qcWaiters       []chan qcResult
	qcLastPrice     decimal.Decimal
	qcLastErr       error
}

type PriceCache struct {
	Price     decimal.Decimal
	Timestamp time.Time
	Source    string
}

type PriceServiceInterface interface {
	GetAssetPrice(asset XRPAsset) (decimal.Decimal, error)
	GetXRPLPrice() (decimal.Decimal, error)
	GetRLUSDPrice() (decimal.Decimal, error)
	GetTokenPriceViaTwoStepConversion(currency, issuer string) (decimal.Decimal, error)
	GetTokenPriceInXRP(currency, issuer string) (decimal.Decimal, error)
	RefreshPrices()

	// New methods to eliminate redundant XRP price calls
	GetTokenPriceWithXRPPrice(currency, issuer string, xrpPriceUSD decimal.Decimal) (decimal.Decimal, error)
	GetAssetPriceWithXRP(asset XRPAsset, xrpPriceUSD decimal.Decimal) (decimal.Decimal, error)
	GetMultipleAssetPrices(assets []XRPAsset) (map[string]decimal.Decimal, error)

	// Batch helpers
	SetBatchXRPPrice(price decimal.Decimal)
	ClearBatchXRPPrice()
}

// SetBatchXRPPrice activates batch override and stores the provided XRP price
func (p *PriceService) SetBatchXRPPrice(price decimal.Decimal) {
	p.batchMutex.Lock()
	p.batchActive = true
	p.batchXRPPrice = price
	p.batchMutex.Unlock()
}

// ClearBatchXRPPrice deactivates batch override
func (p *PriceService) ClearBatchXRPPrice() {
	p.batchMutex.Lock()
	p.batchActive = false
	p.batchXRPPrice = decimal.Zero
	p.batchMutex.Unlock()
}

// CoinGecko API response structures (only for XRP)
type CoinGeckoXRPLResponse struct {
	Xrp struct {
		Usd decimal.Decimal `json:"usd"`
	} `json:"xrp"`
}

func NewPriceService() PriceServiceInterface {
	return &PriceService{
		cache:             make(map[string]PriceCache),
		db:                nil, // Will be set by the caller
		connectionManager: nil, // Will be set by the caller
	}
}

// NewPriceServiceWithDB creates a price service with a shared database connection
func NewPriceServiceWithDB(db *gorm.DB) PriceServiceInterface {
	return &PriceService{
		cache:             make(map[string]PriceCache),
		db:                db,
		connectionManager: nil, // Will be set by the caller
	}
}

// NewPriceServiceWithDBAndConn creates a price service with both database and connection manager
func NewPriceServiceWithDBAndConn(db *gorm.DB, connMgr *ConnectionManager) PriceServiceInterface {
	return &PriceService{
		cache:             make(map[string]PriceCache),
		db:                db,
		connectionManager: connMgr,
	}
}

// GetAssetPrice gets the USD price for any XRPL asset using the correct two-step strategy
func (p *PriceService) GetAssetPrice(asset XRPAsset) (decimal.Decimal, error) {
	log.Printf("🔍 [PRICE SERVICE] GetAssetPrice called for %s/%s", asset.Currency, asset.Issuer)
	log.Printf("🔍 [PRICE SERVICE] Connection manager available: %t", p.connectionManager != nil)
	log.Printf("🔍 [PRICE SERVICE] Database available: %t", p.db != nil)

	// Check cache first
	cacheKey := fmt.Sprintf("%s_%s", asset.Currency, asset.Issuer)
	if cached, exists := p.cache[cacheKey]; exists && time.Since(cached.Timestamp) < 5*time.Minute {
		log.Printf("✅ [PRICE SERVICE] Returning cached price for %s/%s: $%s (age: %v)",
			asset.Currency, asset.Issuer, cached.Price.String(), time.Since(cached.Timestamp))
		return cached.Price, nil
	}

	log.Printf("🔍 [PRICE SERVICE] No cached price found, calculating fresh price")

	var priceUSD decimal.Decimal
	var err error

	// Handle different asset types with the correct strategy
	switch {
	case asset.Currency == "XRP":
		// For XRP, get price directly from XRP/RLUSD pool
		log.Printf("🔍 [PRICE SERVICE] Asset type: XRP - getting price from XRP/RLUSD pool")
		priceUSD, err = p.GetXRPLPrice()
	case asset.Currency == "RLUSD" && asset.Issuer == "rMxCKbEDwqr76QuheSUMdEGf4B9xJ8m5De":
		// RLUSD is pegged to USD
		log.Printf("🔍 [PRICE SERVICE] Asset type: RLUSD - using $1.00 peg")
		priceUSD, err = p.GetRLUSDPrice()
	default:
		// For all other tokens:
		// Step 1: Get token price in XRP from XRP/Token pool
		// Step 2: Convert XRP amount to USD using XRP/RLUSD rate
		log.Printf("🔍 [PRICE SERVICE] Asset type: Token - using two-step conversion (Token→XRP→USD)")
		log.Printf("🔍 [PRICE SERVICE] Starting two-step conversion for %s/%s", asset.Currency, asset.Issuer)
		priceUSD, err = p.GetTokenPriceViaTwoStepConversion(asset.Currency, asset.Issuer)

	}

	if err != nil {
		log.Printf("❌ [PRICE SERVICE] Price calculation failed for %s/%s: %v", asset.Currency, asset.Issuer, err)
		return decimal.Zero, fmt.Errorf("failed to get price for %s/%s: %v", asset.Currency, asset.Issuer, err)
	}

	// Cache the result
	p.cache[cacheKey] = PriceCache{
		Price:     priceUSD,
		Timestamp: time.Now(),
		Source:    "two_step_conversion",
	}

	log.Printf("✅ [PRICE SERVICE] Successfully calculated price for %s/%s: $%s",
		asset.Currency, asset.Issuer, priceUSD.String())
	log.Printf("🔍 [PRICE SERVICE] Caching price for 5 minutes")
	return priceUSD, nil
}

// GetXRPLPrice fetches XRP price with EMERGENCY CACHING to prevent rate limiting
func (p *PriceService) GetXRPLPrice() (decimal.Decimal, error) {
	// If a batch override is active, return it immediately
	p.batchMutex.RLock()
	if p.batchActive {
		price := p.batchXRPPrice
		p.batchMutex.RUnlock()
		if price.GreaterThan(decimal.Zero) {
			return price, nil
		}
		// fall through if zero/invalid
	} else {
		p.batchMutex.RUnlock()
	}
	// EMERGENCY: Check cache first to prevent excessive XRPL API calls
	p.cacheMutex.RLock()
	if cached, exists := p.cache["XRP_PRICE"]; exists {
		if time.Since(cached.Timestamp) < 30*time.Second { // 30s cache for QC and ledger-derived XRP price
			p.cacheMutex.RUnlock()
			// Using cached XRP price (debug logging removed)
			return cached.Price, nil
		}
	}
	p.cacheMutex.RUnlock()

	// Try to get XRP price from database first as fallback
	if dbPrice, err := p.getXRPLPriceFromDatabase(); err == nil && dbPrice.GreaterThan(decimal.Zero) {
		// Cache the database price for future use
		p.cacheMutex.Lock()
		p.cache["XRP_PRICE"] = PriceCache{
			Price:     dbPrice,
			Timestamp: time.Now(),
			Source:    "database_fallback",
		}
		p.cacheMutex.Unlock()
		return dbPrice, nil
	}
	// Database price unavailable, will try QuantifyCrypto fallback
	// Log removed to reduce log spam

	// Try QuantifyCrypto API as final fallback when database fails (throttled)
	if qcPrice, err := p.getXRPLPriceFromQuantifyCryptoThrottled(); err == nil && qcPrice.GreaterThan(decimal.Zero) {
		// Cache the QuantifyCrypto price for future use
		p.cacheMutex.Lock()
		p.cache["XRP_PRICE"] = PriceCache{
			Price:     qcPrice,
			Timestamp: time.Now(),
			Source:    "quantifycrypto_fallback",
		}
		p.cacheMutex.Unlock()
		return qcPrice, nil
	}

	// Fetching LIVE XRP price from XRPL ledger (debug logging removed)

	if p.connectionManager == nil {
		return decimal.Zero, fmt.Errorf("connection manager not available for live price fetching")
	}

	// Get live connection to XRPL with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to get connection with timeout
	done := make(chan *xrpl.Client, 1)
	errChan := make(chan error, 1)

	go func() {
		client, err := p.connectionManager.GetConnection()
		if err != nil {
			errChan <- err
		} else {
			done <- client
		}
	}()

	var xrplClient *xrpl.Client
	select {
	case xrplClient = <-done:
		// Connection obtained successfully
	case err := <-errChan:
		return decimal.Zero, fmt.Errorf("failed to get XRPL connection: %w", err)
	case <-ctx.Done():
		return decimal.Zero, fmt.Errorf("XRPL connection timeout: %w", ctx.Err())
	}
	defer p.connectionManager.ReturnConnection(xrplClient)

	// Call AMM info directly for XRP/RLUSD pool to get LIVE data
	// Use hex-encoded RLUSD currency format for XRPL API calls
	rlusdCurrencyHex := "524C555344000000000000000000000000000000"
	ammInfoRequest := map[string]interface{}{
		"id":      1,
		"command": "amm_info",
		"asset": map[string]interface{}{
			"currency": "XRP",
		},
		"asset2": map[string]interface{}{
			"currency": rlusdCurrencyHex,
			"issuer":   "rMxCKbEDwqr76QuheSUMdEGf4B9xJ8m5De",
		},
		"ledger_index": "validated",
	}

	// // DEBUG logging removed Calling XRPL ledger for LIVE XRP/RLUSD AMM info\n")
	response, err := xrplClient.Request(ammInfoRequest)
	if err != nil {
		// Check for rate limiting and don't cache the failure
		if strings.Contains(err.Error(), "policy violation") ||
			strings.Contains(err.Error(), "IP limit") ||
			strings.Contains(err.Error(), "close 1008") ||
			strings.Contains(err.Error(), "close sent") ||
			strings.Contains(err.Error(), "WS read error") ||
			strings.Contains(err.Error(), "subscription error") {
			if p.connectionManager != nil {
				p.connectionManager.TriggerRateLimitBackoff(15 * time.Second)
			}
			return decimal.Zero, fmt.Errorf("XRPL rate limit reached: %w", err)
		}
		return decimal.Zero, fmt.Errorf("failed to get live AMM info from XRPL: %w", err)
	}

	// Check for XRPL API errors first
	if errorStr, hasError := response["error"].(string); hasError {
		log.Printf("❌ [ERROR] XRPL API error for XRP/RLUSD: %s", errorStr)
		if errorStr == "actNotFound" {
			return decimal.Zero, fmt.Errorf("XRP/RLUSD AMM pool not found on XRPL (critical infrastructure issue)")
		}
		return decimal.Zero, fmt.Errorf("XRPL API error for XRP/RLUSD: %s", errorStr)
	}

	// Parse the live AMM response - XRP can be in amount or amount2
	var xrpBalance, rlusdBalance decimal.Decimal

	if result, ok := response["result"].(map[string]interface{}); ok {
		if ammData, ok := result["amm"].(map[string]interface{}); ok {
			// Check amount field
			if amountStr, ok := ammData["amount"].(string); ok {
				// amount is XRP (in drops format from XRPL amm_info command)
				// XRP value conversion (debug logging removed)
				if xrpValue, err := decimal.NewFromString(amountStr); err == nil {
					xrpBalance = xrpValue.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
				}
				// amount2 should be RLUSD (object)
				if amount2Obj, ok := ammData["amount2"].(map[string]interface{}); ok {
					if value, ok := amount2Obj["value"].(string); ok {
						// DEBUG logging removed
						if balance, err := decimal.NewFromString(value); err == nil {
							rlusdBalance = balance
							// DEBUG logging removed
						} else {
							fmt.Printf("ERROR: Failed to parse RLUSD amount2 value '%s': %v\n", value, err)
						}
					}
				}
			} else if amountObj, ok := ammData["amount"].(map[string]interface{}); ok {
				// amount is RLUSD (object)
				if value, ok := amountObj["value"].(string); ok {
					// DEBUG logging removed RLUSD amount object value: %s\n", value)
					if balance, err := decimal.NewFromString(value); err == nil {
						rlusdBalance = balance
						// DEBUG logging removed
					} else {
						fmt.Printf("ERROR: Failed to parse RLUSD amount value '%s': %v\n", value, err)
					}
				}
				// amount2 should be XRP (string in drops format from XRPL amm_info)
				if amount2Str, ok := ammData["amount2"].(string); ok {
					// DEBUG logging removed [XRP/RLUSD amm amount2]: Raw XRP value in drops: %s\n", amount2Str)
					if xrpValue, err := decimal.NewFromString(amount2Str); err == nil {
						// DEBUG logging removed [XRP/RLUSD amm amount2]: Converting drops to XRP: %s drops\n", xrpValue.String())
						xrpBalance = xrpValue.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
						// DEBUG logging removed [XRP/RLUSD amm amount2]: Final XRP balance: %s XRP\n", xrpBalance.String())
					}
				}
			}
		}
	}

	// Validate balances before calculation
	if xrpBalance.LessThanOrEqual(decimal.Zero) || rlusdBalance.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("invalid pool balances: XRP=%s, RLUSD=%s - AMM response parsing failed", xrpBalance.String(), rlusdBalance.String())
	}

	// Calculate XRP price: RLUSD per XRP (since RLUSD ≈ $1 USD)
	xrpPrice := rlusdBalance.Div(xrpBalance)

	// LIVE XRP price from XRPL ledger (debug logging removed)

	// EMERGENCY: Cache the price for 60 seconds to prevent excessive API calls
	p.cacheMutex.Lock()
	if p.cache == nil {
		p.cache = make(map[string]PriceCache)
	}
	p.cache["XRP_PRICE"] = PriceCache{
		Price:     xrpPrice,
		Timestamp: time.Now(),
		Source:    "XRPL_LIVE_CACHED",
	}
	p.cacheMutex.Unlock()
	// Cached XRP price (debug logging removed)

	return xrpPrice, nil
}

// GetRLUSDPrice gets RLUSD price from market data
func (p *PriceService) GetRLUSDPrice() (decimal.Decimal, error) {
	// RLUSD is a stablecoin pegged to USD, so it should be close to $1
	// In production, this could fetch from multiple AMM pools or oracles
	// For now, return 1.0 as RLUSD ≈ $1 USD
	return decimal.NewFromFloat(1.0), nil
}

// GetTokenPrice fetches the USD price of a token using AMM quotes only
func (p *PriceService) GetTokenPrice(asset XRPAsset) (decimal.Decimal, error) {
	// For XRP itself, get price from XRP/RLUSD pool
	if asset.Currency == "XRP" {
		return p.GetXRPLPrice()
	}

	// For RLUSD, get real market price
	if asset.Currency == "RLUSD" && asset.Issuer == "rMxCKbEDwqr76QuheSUMdEGf4B9xJ8m5De" {
		return p.GetRLUSDPrice()
	}

	// Try to get price from token/RLUSD pool first
	price, err := p.GetTokenPriceFromAMM(asset.Currency, asset.Issuer)
	if err == nil && price.GreaterThan(decimal.Zero) {
		return price, nil
	}

	return decimal.Zero, fmt.Errorf("no AMM pool available for token %s/%s", asset.Currency, asset.Issuer)
}

// GetTokenPriceFromAMM gets token price by calling XRPL ledger directly for real-time pool balances
func (p *PriceService) GetTokenPriceFromAMM(currency, issuer string) (decimal.Decimal, error) {
	if p.db == nil {
		return decimal.Zero, fmt.Errorf("database connection not available")
	}

	// Looking for AMM pool for token

	// First, find the AMM account from normalized database (better performance, no duplicates)
	query := `
		SELECT account, asset2_currency, asset2_issuer
		FROM xrpAmm_normalized 
		WHERE (asset2_currency = ? AND asset2_issuer = ?) 
		   OR (CAST(asset2_currency AS CHAR) = CAST(? AS CHAR) AND asset2_issuer = ?)
		LIMIT 1
	`

	row := p.db.Raw(query, currency, issuer, currency, issuer).Row()

	var account, asset2currency, asset2issuer string
	if err := row.Scan(&account, &asset2currency, &asset2issuer); err != nil {
		return decimal.Zero, fmt.Errorf("no AMM pool available for token %s/%s", currency, issuer)
	}

	// Found AMM pool account for token

	// Now call XRPL ledger directly to get real-time pool balances
	return p.getPoolBalancesFromLedger(account, currency, issuer)
}

// getPoolBalancesFromLedger calls the XRPL ledger directly to get current AMM pool balances
func (p *PriceService) getPoolBalancesFromLedger(poolAccount, tokenCurrency, tokenIssuer string) (decimal.Decimal, error) {
	// Use existing connection manager if available
	if p.connectionManager == nil {
		return decimal.Zero, fmt.Errorf("connection manager not available for price service")
	}

	xrplClient, err := p.connectionManager.GetConnection()
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRPL connection: %w", err)
	}
	defer p.connectionManager.ReturnConnection(xrplClient)

	// Get account info for the AMM pool account to see current balances
	accountInfoRequest := map[string]interface{}{
		"id":           1,
		"command":      "account_info",
		"account":      poolAccount,
		"ledger_index": "validated",
	}

	fmt.Printf("Calling XRPL ledger for pool account: %s\n", poolAccount)
	response, err := xrplClient.Request(accountInfoRequest)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get account info from XRPL: %w", err)
	}

	// Parse XRP balance from account info
	var xrpBalanceDecimal decimal.Decimal
	if result, ok := response["result"].(map[string]interface{}); ok {
		if accountData, ok := result["account_data"].(map[string]interface{}); ok {
			if balanceStr, ok := accountData["Balance"].(string); ok {
				// Balance from AMM pool account
				if balance, err := decimal.NewFromString(balanceStr); err == nil {
					xrpBalanceDecimal = balance
					fmt.Printf("Pool XRP balance from ledger: %s XRP\n", xrpBalanceDecimal.String())
				}
			}
		}
	}

	// Get account lines to see token holdings
	accountLinesRequest := map[string]interface{}{
		"id":           2,
		"command":      "account_lines",
		"account":      poolAccount,
		"ledger_index": "validated",
	}

	linesResponse, err := xrplClient.Request(accountLinesRequest)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get account lines from XRPL: %w", err)
	}

	// Find the specific token balance
	var tokenBalanceDecimal decimal.Decimal
	foundToken := false
	if result, ok := linesResponse["result"].(map[string]interface{}); ok {
		if lines, ok := result["lines"].([]interface{}); ok {
			for _, line := range lines {
				if lineMap, ok := line.(map[string]interface{}); ok {
					lineCurrency, _ := lineMap["currency"].(string)

					// Match the token we're looking for by currency and issuer
					lineAccount, _ := lineMap["account"].(string)
					if lineCurrency == tokenCurrency && lineAccount == tokenIssuer {
						if balance, err := decimal.NewFromString(lineMap["balance"].(string)); err == nil {
							tokenBalanceDecimal = balance.Abs() // Use absolute value for safety
							foundToken = true
							fmt.Printf("Found token %s balance from ledger: %s\n", tokenCurrency, tokenBalanceDecimal.String())
							break
						}
					}
				}
			}
		}
	}

	if !foundToken {
		return decimal.Zero, fmt.Errorf("token %s not found in pool %s", tokenCurrency, poolAccount)
	}

	// CRITICAL: Validate that we have actual balances to work with
	if xrpBalanceDecimal.IsZero() || tokenBalanceDecimal.IsZero() {
		return decimal.Zero, fmt.Errorf("invalid pool balances for %s/%s: XRP=%s, Token=%s",
			tokenCurrency, tokenIssuer, xrpBalanceDecimal.String(), tokenBalanceDecimal.String())
	}

	// Price is Token per XRP
	price := tokenBalanceDecimal.Div(xrpBalanceDecimal)
	fmt.Printf("Calculated price for %s: %s tokens per XRP\n", tokenCurrency, price.String())

	return price, nil
}

// RefreshPrices clears the cache to force fresh price fetching
func (p *PriceService) RefreshPrices() {
	p.cache = make(map[string]PriceCache)
}

// GetTokenPriceViaTwoStepConversion implements the correct pricing strategy using direct AMM calls:
// Step 1: Get XRP price in USD from direct XRP/RLUSD AMM call
// Step 2: Get token price in XRP from direct XRP/Token AMM call
// Step 3: Convert to USD (Token price in XRP × XRP price in USD)
func (p *PriceService) GetTokenPriceViaTwoStepConversion(currency, issuer string) (decimal.Decimal, error) {
	// Starting two-step conversion (debug logging removed)

	// Step 1: Get XRP price in USD from XRP/RLUSD pool (already working)
	xrpPriceUSD, err := p.GetXRPLPrice()
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRP price in USD: %w", err)
	}
	// Step 1 - XRP price in USD (debug logging removed)

	// Step 2: Get token price in XRP from direct XRP/Token AMM call
	tokenPriceInXRP, err := p.GetTokenPriceInXRPDirect(currency, issuer)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get token price in XRP: %w", err)
	}
	// Step 2 - Token price in XRP (debug logging removed)

	// Step 3: Convert to USD (Token price in XRP × XRP price in USD)
	tokenPriceUSD := tokenPriceInXRP.Mul(xrpPriceUSD)
	// Step 3 - Final conversion (debug logging removed)

	return tokenPriceUSD, nil
}

// GetTokenPriceInXRPDirect gets a token's price in XRP using direct AMM info call (bypasses database)
func (p *PriceService) GetTokenPriceInXRPDirect(currency, issuer string) (decimal.Decimal, error) {
	if p.connectionManager == nil {
		return decimal.Zero, fmt.Errorf("connection manager not available for direct AMM calls")
	}

	// Getting token price directly from XRPL AMM (debug logging removed)

	// Get live connection to XRPL
	xrplClient, err := p.connectionManager.GetConnection()
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRPL connection: %w", err)
	}
	defer p.connectionManager.ReturnConnection(xrplClient)

	// Try with original currency format first, then fallback to opposite format
	var response map[string]interface{}

	if len(currency) != 40 {
		// Method 1: Try with string currency directly first (3-6 characters or non-hex)
		// Attempting with string currency (debug logging removed)

		ammInfoRequest := map[string]interface{}{
			"id":      1,
			"command": "amm_info",
			"asset": map[string]interface{}{
				"currency": "XRP",
			},
			"asset2": map[string]interface{}{
				"currency": currency,
				"issuer":   issuer,
			},
			"ledger_index": "validated",
		}

		response, err = xrplClient.Request(ammInfoRequest)
		if err != nil {
			log.Printf("📉 [PRICE METRICS] amm_info request_error")
			return decimal.Zero, fmt.Errorf("failed to get AMM info: %w", err)
		}

		// Check if we got actNotFound error - if so, try hex conversion
		if errorStr, hasError := response["error"].(string); hasError && errorStr == "actNotFound" {
			// String currency failed with actNotFound, trying hex conversion (debug logging removed)

			// Method 2: Convert string to hex and try again (pad to 20 bytes)
			if len(currency) <= 20 { // Only convert if reasonable length
				currencyBytes := []byte(currency)
				for len(currencyBytes) < 20 {
					currencyBytes = append(currencyBytes, 0)
				}
				currencyHex := strings.ToUpper(hex.EncodeToString(currencyBytes))
				// DEBUG logging removed Converted %s to hex: %s\n", currency, currencyHex)

				ammInfoRequestHex := map[string]interface{}{
					"id":      2,
					"command": "amm_info",
					"asset": map[string]interface{}{
						"currency": "XRP",
					},
					"asset2": map[string]interface{}{
						"currency": currencyHex,
						"issuer":   issuer,
					},
					"ledger_index": "validated",
				}

				response, err = xrplClient.Request(ammInfoRequestHex)
				if err != nil {
					log.Printf("📉 [PRICE METRICS] amm_info request_error")
					return decimal.Zero, fmt.Errorf("failed to get AMM info with hex conversion: %w", err)
				}
			} else {
				// DEBUG logging removed Currency too long for hex conversion: %d characters\n", len(currency))
			}
		}
	} else if len(currency) == 40 {
		// Method 1: Try with hex currency first
		// DEBUG logging removed Attempting with hex currency: %s\n", currency)

		ammInfoRequest := map[string]interface{}{
			"id":      1,
			"command": "amm_info",
			"asset": map[string]interface{}{
				"currency": "XRP",
			},
			"asset2": map[string]interface{}{
				"currency": currency,
				"issuer":   issuer,
			},
			"ledger_index": "validated",
		}

		response, err = xrplClient.Request(ammInfoRequest)
		if err != nil {
			log.Printf("📉 [PRICE METRICS] amm_info request_error")
			return decimal.Zero, fmt.Errorf("failed to get AMM info: %w", err)
		}

		// Check if we got actNotFound error - if so, try converting hex back to string
		if errorStr, hasError := response["error"].(string); hasError && errorStr == "actNotFound" {
			// DEBUG logging removed Hex currency failed with actNotFound, trying to convert back to string\n")

			// Method 2: Try to decode hex back to 3-letter string
			currencyBytes, hexErr := hex.DecodeString(currency)
			if hexErr == nil {
				// Find the end of the actual currency (before zero padding)
				var currencyStr string
				for i, b := range currencyBytes {
					if b == 0 {
						currencyStr = string(currencyBytes[:i])
						break
					}
				}

				if len(currencyStr) >= 3 && len(currencyStr) <= 6 {
					// DEBUG logging removed Converted hex %s back to string: %s\n", currency, currencyStr)

					ammInfoRequestStr := map[string]interface{}{
						"id":      2,
						"command": "amm_info",
						"asset": map[string]interface{}{
							"currency": "XRP",
						},
						"asset2": map[string]interface{}{
							"currency": currencyStr,
							"issuer":   issuer,
						},
						"ledger_index": "validated",
					}

					response, err = xrplClient.Request(ammInfoRequestStr)
					if err != nil {
						log.Printf("📉 [PRICE METRICS] amm_info request_error")
						return decimal.Zero, fmt.Errorf("failed to get AMM info with string conversion: %w", err)
					}
				} else {
					// DEBUG logging removed Hex decode didn't produce a valid currency (3-6 chars): %s (length: %d)\n", currencyStr, len(currencyStr))
				}
			} else {
				// DEBUG logging removed Failed to decode hex currency: %v\n", hexErr)
			}
		}
	} else {
		// Unknown format - try as-is
		// DEBUG logging removed Unknown currency format (length %d): %s\n", len(currency), currency)

		ammInfoRequest := map[string]interface{}{
			"id":      1,
			"command": "amm_info",
			"asset": map[string]interface{}{
				"currency": "XRP",
			},
			"asset2": map[string]interface{}{
				"currency": currency,
				"issuer":   issuer,
			},
			"ledger_index": "validated",
		}

		response, err = xrplClient.Request(ammInfoRequest)
		if err != nil {
			log.Printf("📉 [PRICE METRICS] amm_info request_error")
			return decimal.Zero, fmt.Errorf("failed to get AMM info: %w", err)
		}
	}

	// DEBUG logging removed Final AMM response for XRP/%s: %+v\n", currency, response)

	// Check for errors in the final response
	if errorStr, hasError := response["error"].(string); hasError {
		log.Printf("📉 [PRICE METRICS] error:%s", errorStr)
		// DEBUG logging removed Final XRPL API error for XRP/%s: %s\n", currency, errorStr)
		if errorStr == "actNotFound" {
			return decimal.Zero, fmt.Errorf("AMM pool not found on XRPL for XRP/%s (tried both original and fallback formats)", currency)
		}
		return decimal.Zero, fmt.Errorf("XRPL API error for XRP/%s: %s", currency, errorStr)
	}

	// Reject non-validated results
	if result, ok := response["result"].(map[string]interface{}); ok {
		if validated, okv := result["validated"].(bool); okv && !validated {
			return decimal.Zero, fmt.Errorf("non-validated amm_info result")
		}
	}

	// Parse the AMM response - XRP can be in amount or amount2
	var xrpBalance, tokenBalance decimal.Decimal

	if result, ok := response["result"].(map[string]interface{}); ok {
		if ammData, ok := result["amm"].(map[string]interface{}); ok {
			// Check amount field
			if amountStr, ok := ammData["amount"].(string); ok {
				// amount is XRP (in drops format from XRPL amm_info command)
				// DEBUG logging removed [Token/XRP amm amount]: Raw XRP value in drops: %s\n", amountStr)
				if xrpValue, err := decimal.NewFromString(amountStr); err == nil {
					// DEBUG logging removed [Token/XRP amm amount]: Converting drops to XRP: %s drops\n", xrpValue.String())
					xrpBalance = xrpValue.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
					// DEBUG logging removed [Token/XRP amm amount]: Final XRP balance: %s XRP\n", xrpBalance.String())
				}
				// amount2 should be token (object)
				if amount2Obj, ok := ammData["amount2"].(map[string]interface{}); ok {
					if value, ok := amount2Obj["value"].(string); ok {
						// DEBUG logging removed Token amount2 object value: %s\n", value)
						if balance, err := decimal.NewFromString(value); err == nil {
							tokenBalance = balance
							// DEBUG logging removed Successfully parsed token balance: %s\n", tokenBalance.String())
						} else {
							fmt.Printf("ERROR: Failed to parse token amount2 value '%s': %v\n", value, err)
						}
					}
				}
			} else if amountObj, ok := ammData["amount"].(map[string]interface{}); ok {
				// amount is token (object)
				if value, ok := amountObj["value"].(string); ok {
					// DEBUG logging removed Token amount object value: %s\n", value)
					if balance, err := decimal.NewFromString(value); err == nil {
						tokenBalance = balance
						// DEBUG logging removed Successfully parsed token balance: %s\n", tokenBalance.String())
					} else {
						fmt.Printf("ERROR: Failed to parse token amount value '%s': %v\n", value, err)
					}
				}
				// amount2 should be XRP (string in drops format from XRPL amm_info)
				if amount2Str, ok := ammData["amount2"].(string); ok {
					// DEBUG logging removed [Token/XRP amm amount2]: Raw XRP value in drops: %s\n", amount2Str)
					if xrpValue, err := decimal.NewFromString(amount2Str); err == nil {
						// DEBUG logging removed [Token/XRP amm amount2]: Converting drops to XRP: %s drops\n", xrpValue.String())
						xrpBalance = xrpValue.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
						// DEBUG logging removed [Token/XRP amm amount2]: Final XRP balance: %s XRP\n", xrpBalance.String())
					}
				}
			}
		}
	}

	// Calculate token price in XRP: XRP_balance / Token_balance
	tokenPriceInXRP := xrpBalance.Div(tokenBalance)

	// Direct AMM price calculation (debug logging removed)

	// Returning token price in XRP (debug logging removed)
	return tokenPriceInXRP, nil
}

// GetTokenPriceInXRP gets a token's price in XRP using the XRP/Token AMM pool
func (p *PriceService) GetTokenPriceInXRP(currency, issuer string) (decimal.Decimal, error) {
	if p.db == nil {
		return decimal.Zero, fmt.Errorf("database connection not available")
	}

	// DEBUG logging removed Looking for XRP/%s pool in database...\n", currency)

	// Find AMM pool where one asset is XRP and the other is our token
	// XRP is represented as native currency (no issuer)
	query := `
		SELECT account, asset1_amount, asset2_currency, asset2_issuer, asset2_amount
		FROM xrpAmm_normalized 
		WHERE (
			-- Case 1: XRP is asset1 (native), token is asset2
			(asset1_currency = 'XRP' AND asset2_currency = ? AND asset2_issuer = ?) OR
			-- Case 2: Token is asset1, XRP is asset2
			(asset1_currency = ? AND asset1_issuer = ? AND asset2_currency = 'XRP')
		)
		LIMIT 1
	`

	row := p.db.Raw(query, currency, issuer, currency, issuer).Row()

	var account, asset1_amount, asset2_currency, asset2_issuer, asset2_amount string
	if err := row.Scan(&account, &asset1_amount, &asset2_currency, &asset2_issuer, &asset2_amount); err != nil {
		return decimal.Zero, fmt.Errorf("no XRP/%s AMM pool found in database", currency)
	}

	// DEBUG logging removed Found XRP/%s pool: account=%s\n", currency, account)

	// Now get current balances from XRPL ledger
	return p.getTokenPriceInXRPFromLedger(account, currency, issuer, asset1_amount, asset2_currency, asset2_amount)
}

// getTokenPriceInXRPFromLedger gets real-time token price in XRP from XRPL ledger
func (p *PriceService) getTokenPriceInXRPFromLedger(poolAccount, tokenCurrency, tokenIssuer, dbAmount, dbAmount2Currency, dbAmount2Value string) (decimal.Decimal, error) {
	// Use existing connection manager if available
	if p.connectionManager == nil {
		return decimal.Zero, fmt.Errorf("connection manager not available for price service")
	}

	xrplClient, err := p.connectionManager.GetConnection()
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRPL connection: %w", err)
	}
	defer p.connectionManager.ReturnConnection(xrplClient)

	// DEBUG logging removed Getting real-time balances for pool %s\n", poolAccount)

	// Get XRP balance from account info
	accountInfoRequest := map[string]interface{}{
		"id":           1,
		"command":      "account_info",
		"account":      poolAccount,
		"ledger_index": "validated",
	}

	response, err := xrplClient.Request(accountInfoRequest)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get account info from XRPL: %w", err)
	}

	// Parse XRP balance
	var xrpBalance decimal.Decimal
	if result, ok := response["result"].(map[string]interface{}); ok {
		if accountData, ok := result["account_data"].(map[string]interface{}); ok {
			if balanceStr, ok := accountData["Balance"].(string); ok {
				// Balance from AMM pool account (in drops format from XRPL API)
				if balance, err := decimal.NewFromString(balanceStr); err == nil {
					xrpBalance = balance.Div(decimal.NewFromInt(1000000)) // Convert drops to XRP
					// DEBUG logging removed Pool XRP balance: %s drops → %s XRP\n", balanceStr, xrpBalance.String())
				}
			}
		}
	}

	// Get token balance from account lines
	accountLinesRequest := map[string]interface{}{
		"id":           2,
		"command":      "account_lines",
		"account":      poolAccount,
		"ledger_index": "validated",
	}

	linesResponse, err := xrplClient.Request(accountLinesRequest)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get account lines from XRPL: %w", err)
	}

	// Find the specific token balance
	var tokenBalance decimal.Decimal
	if result, ok := linesResponse["result"].(map[string]interface{}); ok {
		if lines, ok := result["lines"].([]interface{}); ok {
			for _, line := range lines {
				if lineMap, ok := line.(map[string]interface{}); ok {
					currency, _ := lineMap["currency"].(string)
					account, _ := lineMap["account"].(string)
					balanceStr, _ := lineMap["balance"].(string)

					// Match our token
					if currency == tokenCurrency && account == tokenIssuer {
						if balance, err := decimal.NewFromString(balanceStr); err == nil {
							tokenBalance = balance
							// DEBUG logging removed Pool token balance: %s %s\n", tokenBalance.String(), tokenCurrency)
							break
						}
					}
				}
			}
		}
	}

	// Validate balances
	if xrpBalance.LessThanOrEqual(decimal.Zero) || tokenBalance.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("invalid pool balances: XRP=%s, %s=%s", xrpBalance.String(), tokenCurrency, tokenBalance.String())
	}

	// Calculate token price in XRP using constant product formula
	// Price = XRP_balance / Token_balance (how much XRP per 1 token)
	tokenPriceInXRP := xrpBalance.Div(tokenBalance)

	// Price calculation (debug logging removed)

	return tokenPriceInXRP, nil
}

// GetTokenPriceWithXRPPrice calculates token price using pre-fetched XRP price to avoid redundant XRPL calls
func (p *PriceService) GetTokenPriceWithXRPPrice(currency, issuer string, xrpPriceUSD decimal.Decimal) (decimal.Decimal, error) {
	log.Printf("🔍 [PRICE SERVICE] GetTokenPriceWithXRPPrice called for %s/%s with XRP price $%s",
		currency, issuer, xrpPriceUSD.StringFixed(6))

	// Check cache first
	cacheKey := fmt.Sprintf("%s_%s", currency, issuer)
	p.cacheMutex.RLock()
	if cached, exists := p.cache[cacheKey]; exists && time.Since(cached.Timestamp) < 5*time.Minute {
		p.cacheMutex.RUnlock()
		log.Printf("✅ [PRICE SERVICE] Returning cached price for %s/%s: $%s",
			currency, issuer, cached.Price.String())
		return cached.Price, nil
	}
	p.cacheMutex.RUnlock()

	// Skip Step 1 - we already have XRP price
	log.Printf("🔍 [PRICE SERVICE] Using provided XRP price: $%s (skipping GetXRPLPrice call)",
		xrpPriceUSD.StringFixed(6))

	// Step 2: Get token price in XRP from direct XRP/Token AMM call
	tokenPriceInXRP, err := p.GetTokenPriceInXRPDirect(currency, issuer)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get token price in XRP: %w", err)
	}
	log.Printf("🔍 [PRICE SERVICE] Token price in XRP: %s XRP", tokenPriceInXRP.StringFixed(6))

	// Step 3: Convert to USD using provided XRP price
	tokenPriceUSD := tokenPriceInXRP.Mul(xrpPriceUSD)
	log.Printf("✅ [PRICE SERVICE] Final price: %s XRP × $%s = $%s USD",
		tokenPriceInXRP.StringFixed(6), xrpPriceUSD.StringFixed(6), tokenPriceUSD.StringFixed(6))

	// Cache the result
	p.cacheMutex.Lock()
	p.cache[cacheKey] = PriceCache{
		Price:     tokenPriceUSD,
		Timestamp: time.Now(),
		Source:    "two_step_conversion_shared_xrp",
	}
	p.cacheMutex.Unlock()

	return tokenPriceUSD, nil
}

// GetAssetPriceWithXRP gets asset price using pre-fetched XRP price to avoid redundant calls
func (p *PriceService) GetAssetPriceWithXRP(asset XRPAsset, xrpPriceUSD decimal.Decimal) (decimal.Decimal, error) {
	log.Printf("🔍 [PRICE SERVICE] GetAssetPriceWithXRP called for %s/%s with XRP price $%s",
		asset.Currency, asset.Issuer, xrpPriceUSD.StringFixed(6))

	// Check cache first
	cacheKey := fmt.Sprintf("%s_%s", asset.Currency, asset.Issuer)
	p.cacheMutex.RLock()
	if cached, exists := p.cache[cacheKey]; exists && time.Since(cached.Timestamp) < 5*time.Minute {
		p.cacheMutex.RUnlock()
		log.Printf("✅ [PRICE SERVICE] Returning cached price for %s/%s: $%s",
			asset.Currency, asset.Issuer, cached.Price.String())
		return cached.Price, nil
	}
	p.cacheMutex.RUnlock()

	var priceUSD decimal.Decimal
	var err error

	// Handle different asset types
	switch {
	case asset.Currency == "XRP":
		// For XRP, use the provided price directly
		log.Printf("🔍 [PRICE SERVICE] Asset is XRP - using provided price $%s", xrpPriceUSD.StringFixed(6))
		priceUSD = xrpPriceUSD
	case asset.Currency == "RLUSD" && asset.Issuer == "rMxCKbEDwqr76QuheSUMdEGf4B9xJ8m5De":
		// RLUSD is pegged to USD
		log.Printf("🔍 [PRICE SERVICE] Asset is RLUSD - using $1.00 peg")
		priceUSD = decimal.NewFromFloat(1.0)
	default:
		// For all other tokens, use the new method with shared XRP price
		log.Printf("🔍 [PRICE SERVICE] Token asset - using two-step conversion with shared XRP price")
		priceUSD, err = p.GetTokenPriceWithXRPPrice(asset.Currency, asset.Issuer, xrpPriceUSD)
		if err != nil {
			return decimal.Zero, fmt.Errorf("failed to get price for %s/%s: %v", asset.Currency, asset.Issuer, err)
		}
	}

	return priceUSD, nil
}

// GetMultipleAssetPrices gets prices for multiple assets with a single XRP price fetch
func (p *PriceService) GetMultipleAssetPrices(assets []XRPAsset) (map[string]decimal.Decimal, error) {
	log.Printf("🔍 [PRICE SERVICE] GetMultipleAssetPrices called for %d assets", len(assets))

	if len(assets) == 0 {
		return make(map[string]decimal.Decimal), nil
	}

	// Get XRP price ONCE for all calculations
	xrpPriceUSD, err := p.GetXRPLPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRP price: %w", err)
	}
	log.Printf("✅ [PRICE SERVICE] Got XRP price once: $%s (saving %d redundant calls)",
		xrpPriceUSD.StringFixed(6), len(assets)-1)

	// Calculate all asset prices using the shared XRP price
	prices := make(map[string]decimal.Decimal)
	successCount := 0
	failureCount := 0

	for _, asset := range assets {
		assetKey := fmt.Sprintf("%s_%s", asset.Currency, asset.Issuer)

		price, err := p.GetAssetPriceWithXRP(asset, xrpPriceUSD)
		if err != nil {
			log.Printf("⚠️ [PRICE SERVICE] Failed to get price for %s: %v", assetKey, err)
			prices[assetKey] = decimal.Zero
			failureCount++
			continue
		}

		prices[assetKey] = price
		successCount++
	}

	log.Printf("✅ [PRICE SERVICE] Batch pricing complete: %d successful, %d failed out of %d assets",
		successCount, failureCount, len(assets))

	return prices, nil
}

// getXRPLPriceFromDatabase attempts to get XRP price from cached database data
func (p *PriceService) getXRPLPriceFromDatabase() (decimal.Decimal, error) {
	if p.db == nil {
		return decimal.Zero, fmt.Errorf("database not available")
	}

	// Try to get XRP price from xrpTokens table (XRP has empty issuer)
	var price float64
	err := p.db.Raw(`
		SELECT price FROM xrpTokens 
		WHERE currency = 'XRP' AND (issuer = '' OR issuer IS NULL)
		ORDER BY updated_at DESC 
		LIMIT 1
	`).Scan(&price).Error

	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRP price from database: %w", err)
	}

	if price <= 0 {
		return decimal.Zero, fmt.Errorf("no valid XRP price found in database")
	}

	return decimal.NewFromFloat(price), nil
}

// getXRPLPriceFromQuantifyCrypto fetches XRP price from QuantifyCrypto API as fallback
func (p *PriceService) getXRPLPriceFromQuantifyCrypto() (decimal.Decimal, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request to QuantifyCrypto API
	req, err := http.NewRequest("GET", "https://quantifycrypto.com/api/v1/coins/XRP?currency=USD&include_signals=true&signal_type=trend", nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to create QuantifyCrypto request: %w", err)
	}

	// Add required headers (keys should be provided via env; these literals are placeholders in legacy flow)
	req.Header.Set("accept", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to fetch XRP price from QuantifyCrypto: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("QuantifyCrypto API returned status %d", resp.StatusCode)
	}

	// Parse JSON response
	var qcResponse struct {
		Data struct {
			CoinPrice float64 `json:"coin_price"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&qcResponse); err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse QuantifyCrypto response: %w", err)
	}

	if qcResponse.Data.CoinPrice <= 0 {
		return decimal.Zero, fmt.Errorf("invalid XRP price from QuantifyCrypto: %f", qcResponse.Data.CoinPrice)
	}

	price := decimal.NewFromFloat(qcResponse.Data.CoinPrice)
	log.Printf("💰 [PRICE SERVICE] Got XRP price from QuantifyCrypto fallback: $%s", price.StringFixed(6))

	return price, nil
}
