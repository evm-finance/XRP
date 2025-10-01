package xrp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"qc-defi-graphql-server/internal/models"

	"github.com/shopspring/decimal"

	xrpl "github.com/xrpscan/xrpl-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// XRPService provides comprehensive XRP functionality with real XRPL WebSocket calls
type XRPService struct {
	connectionManager *ConnectionManager
	db                *gorm.DB // Add database connection for price service
}

// QcResponse represents QuantifyCrypto API response structure

// SetDatabase sets the database connection for the XRP service
func (svc *XRPService) SetDatabase(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}
	svc.db = db
	return nil
}

// SetConnectionManager sets the connection manager for the XRP service
func (svc *XRPService) SetConnectionManager(connMgr *ConnectionManager) error {
	if connMgr == nil {
		return fmt.Errorf("connection manager cannot be nil")
	}
	svc.connectionManager = connMgr
	return nil
}

// GetClient returns a connection from the connection pool
func (svc *XRPService) GetClient() (*xrpl.Client, error) {
	if svc.connectionManager == nil {
		return nil, fmt.Errorf("connection manager not initialized")
	}
	return svc.connectionManager.GetConnection()
}

// GetConnectionManager returns the connection manager for external access
func (svc *XRPService) GetConnectionManager() *ConnectionManager {
	return svc.connectionManager
}

// NewXRPService creates a new XRP service with connection pooling
func NewXRPService(connectionURL string, maxConnections int, logger *zap.Logger) *XRPService {
	connectionManager := NewConnectionManager(connectionURL, maxConnections, logger)
	return &XRPService{
		connectionManager: connectionManager,
	}
}

// NewXRPServiceWithConnectionManager creates a new XRP service using an existing connection manager
func NewXRPServiceWithConnectionManager(connMgr *ConnectionManager) *XRPService {
	return &XRPService{
		connectionManager: connMgr,
	}
}

// GetTokenListFromAPI fetches token list from the API to populate database
func (svc *XRPService) GetTokenListFromAPI() (*models.XRPTokenList, error) {
	var tokenList models.XRPTokenList
	url := "https://s1.xrplmeta.org/tokens"

	method := "GET"
	connection := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	response, err := connection.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	err = json.Unmarshal(data, &tokenList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &tokenList, nil
}

func (svc *XRPService) GetTokenList() ([]*models.XRPTokenFields, error) {
	// Get tokens from updated xrpTokens table
	log.Printf("🔍 [XRP SERVICE] GetTokenList called")
	if svc.db == nil {
		log.Printf("❌ [XRP SERVICE] Database connection not available")
		return nil, fmt.Errorf("database connection not available")
	}
	log.Printf("✅ [XRP SERVICE] Database connection available")

	var tokens []struct {
		Currency  string  `json:"currency"`
		Issuer    string  `json:"issuer"`
		Name      []byte  `json:"name"` // FIXED: Handle []uint8 from database
		Icon      string  `json:"icon"`
		Price     float64 `json:"price"`
		Marketcap float64 `json:"marketCap"`
		Supply    float64 `json:"supply_xrpl"`
		Liquidity float64 `json:"liquidity_total"`
	}

	// Query the updated xrpTokens table
	query := `
		SELECT currency, issuer, CONVERT(token_name USING utf8) as name, icon, price, marketcap, supply_xrpl, liquidity_total
		FROM xrpTokens 
		WHERE currency != '' AND issuer != ''
		ORDER BY (CASE WHEN price > 0 THEN 1 ELSE 0 END) DESC, price DESC, marketcap DESC, trustlines DESC
		LIMIT 100
	`

	log.Printf("🔍 [XRP SERVICE] Executing token list query:")
	log.Printf("🔍 [XRP SERVICE] Query: %s", query)

	err := svc.db.Raw(query).Scan(&tokens).Error

	if err != nil {
		log.Printf("❌ [XRP SERVICE] Token query failed: %v", err)
		return nil, fmt.Errorf("failed to query tokens from database: %w", err)
	}

	log.Printf("✅ [XRP SERVICE] Token query successful: %d tokens retrieved", len(tokens))

	// Convert to expected format with price analysis
	log.Printf("🔍 [XRP SERVICE] Converting %d tokens to API format", len(tokens))

	var result []*models.XRPTokenFields
	validPrices := 0
	zeroPrices := 0

	for i, token := range tokens {
		if i < 3 { // Log first 3 tokens for debugging
			log.Printf("🔍 [XRP SERVICE] Sample token %d: %s (%s) - Price: $%.6f",
				i+1, token.Name, token.Currency, token.Price)
		}

		// Track price statistics
		if token.Price > 0 {
			validPrices++
		} else {
			zeroPrices++
		}

		// FIXED: Convert []byte to string for token name
		tokenName := string(token.Name)

		result = append(result, &models.XRPTokenFields{
			Currency:      token.Currency,
			IssuerAddress: token.Issuer,
			TokenName:     tokenName, // Convert []byte to string
			IssuerName:    "",        // Not stored in our database
			Icon:          token.Icon,
			Marketcap:     token.Marketcap,
			Price:         token.Price,
			Supply:        token.Supply,
			Liquidity:     token.Liquidity,
		})
	}

	log.Printf("📊 [XRP SERVICE] Token price analysis:")
	log.Printf("  • Valid prices (>0): %d (%.1f%%)", validPrices,
		float64(validPrices)/float64(len(tokens))*100)
	log.Printf("  • Zero prices: %d (%.1f%%)", zeroPrices,
		float64(zeroPrices)/float64(len(tokens))*100)
	log.Printf("✅ [XRP SERVICE] Returning %d tokens to API", len(result))

	return result, nil
}

// GetTokenSupply retrieves the total supply of a token using XRPL gateway_balances API
// Total supply = obligations + sum(all balances) for the given currency
func (svc *XRPService) GetTokenSupply(currency string, issuer string) (decimal.Decimal, error) {
	// Input validation
	if currency == "" {
		return decimal.Zero, fmt.Errorf("currency cannot be empty")
	}
	if issuer == "" {
		return decimal.Zero, fmt.Errorf("issuer cannot be empty")
	}

	// Get client from connection pool
	client, err := svc.GetClient()
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer svc.connectionManager.ReturnConnection(client)

	fmt.Printf("📊 [TOKEN SUPPLY] Fetching supply for %s issued by %s\n", currency, issuer)

	// Use XRPL WebSocket API to get gateway balances
	gatewayRequest := xrpl.BaseRequest{
		"id":           3,
		"command":      "gateway_balances",
		"account":      issuer,
		"strict":       true,
		"ledger_index": "validated", // Use validated ledger for finalized data
	}

	response, err := client.Request(gatewayRequest)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get gateway balances: %w", err)
	}

	// Parse the response
	responseMap := map[string]interface{}(response)
	resultData, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		return decimal.Zero, fmt.Errorf("invalid response format from gateway_balances")
	}

	fmt.Printf("🔍 [TOKEN SUPPLY] Gateway balances response received for %s\n", issuer)

	// Calculate total supply: obligations + sum(balances)
	totalSupply := decimal.Zero

	// Add obligations (tokens the gateway owes)
	if obligations, ok := resultData["obligations"].(map[string]interface{}); ok {
		if obligationStr, exists := obligations[currency]; exists {
			if obligationAmount, ok := obligationStr.(string); ok {
				obligationDecimal, err := decimal.NewFromString(obligationAmount)
				if err != nil {
					fmt.Printf("⚠️ [TOKEN SUPPLY] Failed to parse obligation amount: %v\n", err)
				} else {
					totalSupply = totalSupply.Add(obligationDecimal)
					fmt.Printf("💰 [TOKEN SUPPLY] Obligations for %s: %s\n", currency, obligationDecimal.String())
				}
			}
		}
	}

	// Add balances (tokens held in accounts)
	if balances, ok := resultData["balances"].(map[string]interface{}); ok {
		if currencyBalances, exists := balances[currency]; exists {
			if balanceArray, ok := currencyBalances.([]interface{}); ok {
				for _, balanceItem := range balanceArray {
					if balanceObj, ok := balanceItem.(map[string]interface{}); ok {
						if valueStr, ok := balanceObj["value"].(string); ok {
							balanceDecimal, err := decimal.NewFromString(valueStr)
							if err != nil {
								fmt.Printf("⚠️ [TOKEN SUPPLY] Failed to parse balance amount: %v\n", err)
							} else {
								totalSupply = totalSupply.Add(balanceDecimal)
								fmt.Printf("💰 [TOKEN SUPPLY] Balance for %s: %s\n", currency, balanceDecimal.String())
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("✅ [TOKEN SUPPLY] Total supply for %s: %s\n", currency, totalSupply.String())
	return totalSupply, nil
}

func (svc *XRPService) GetLedger(ledgerIdentifier string) (interface{}, error) {
	// Input validation
	if ledgerIdentifier == "" {
		return nil, fmt.Errorf("ledger identifier cannot be empty")
	}

	// Get client from connection pool
	client, err := svc.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer svc.connectionManager.ReturnConnection(client)

	// Use XRPL WebSocket API to get ledger details
	ledgerRequest := xrpl.BaseRequest{
		"id":           1,
		"command":      "ledger",
		"ledger_index": ledgerIdentifier,
		"accounts":     false,
		"full":         false,
		"transactions": true,
		"expand":       true,
	}

	response, err := client.Request(ledgerRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to request ledger: %w", err)
	}

	// Check for XRPL API errors in response
	if status, ok := response["status"].(string); ok && status == "error" {
		if errorMsg, ok := response["error"].(string); ok {
			if errorMsg == "lgrNotFound" {
				return nil, fmt.Errorf("ledger not found: %s", ledgerIdentifier)
			}
			return nil, fmt.Errorf("XRPL API error: %s", errorMsg)
		}
		return nil, fmt.Errorf("XRPL API returned error status")
	}

	// Return the raw response from XRPL
	return response, nil
}

func (svc *XRPService) GetTransaction(txHash string) (interface{}, error) {
	// Input validation
	if txHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	// Get client from connection pool
	client, err := svc.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer svc.connectionManager.ReturnConnection(client)

	// Use XRPL WebSocket API to get transaction details
	transactionRequest := xrpl.BaseRequest{
		"id":          1,
		"command":     "tx",
		"transaction": txHash,
		"binary":      false,
	}

	response, err := client.Request(transactionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to request transaction: %w", err)
	}

	// Check for XRPL API errors in response
	if status, ok := response["status"].(string); ok && status == "error" {
		if errorMsg, ok := response["error"].(string); ok {
			if errorMsg == "txnNotFound" {
				return nil, fmt.Errorf("transaction not found: %s", txHash)
			}
			return nil, fmt.Errorf("XRPL API error: %s", errorMsg)
		}
		return nil, fmt.Errorf("XRPL API returned error status")
	}

	// Return the raw response from XRPL
	return response, nil
}

// GetAccountBalances retrieves XRP account balances including XRP and token trustlines using real XRPL WebSocket calls
func (svc *XRPService) GetAccountBalances(address string) (*models.XRPAccountBalances, error) {
	// GetAccountBalances called (debug logging reduced)

	// Input validation
	if address == "" {
		log.Printf("❌ [ERROR] Address is empty")
		return nil, fmt.Errorf("address cannot be empty")
	}

	// Basic XRP address format validation (starts with 'r' and is 25-35 characters)
	if len(address) < 25 || len(address) > 35 || !strings.HasPrefix(address, "r") {
		log.Printf("❌ [ERROR] Invalid XRP address format: %s", address)
		return nil, fmt.Errorf("invalid XRP address format: %s", address)
	}

	// Address validation passed (debug logging removed)

	// Check database connection
	if svc.db == nil {
		log.Printf("❌ [ERROR] Database connection not available")
		return nil, fmt.Errorf("database connection not available")
	}
	// Database connection available (debug logging removed)

	// Get client from connection pool
	// Getting XRPL client (debug logging removed)
	client, err := svc.GetClient()
	if err != nil {
		log.Printf("❌ [ERROR] Failed to get XRPL client: %v", err)
		return nil, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer svc.connectionManager.ReturnConnection(client)
	// XRPL client obtained successfully (debug logging removed)

	// Prepare response structure
	// Preparing response structure (debug logging removed)
	xrpBalancesResponse := models.XRPAccountBalances{}
	xrpBalancesResponse.Account = address

	// Get XRP price from XRP/RLUSD AMM pool using our existing price service
	var xrpPrice decimal.Decimal
	if svc.db != nil {
		// Creating price service (debug logging removed)
		priceService := NewPriceServiceWithDBAndConn(svc.db, svc.connectionManager)

		// Calling GetXRPLPrice (debug logging removed)
		price, err := priceService.GetXRPLPrice()
		if err == nil {
			xrpPrice = price
			log.Printf("✅ [SUCCESS] XRP price obtained: $%s", xrpPrice.StringFixed(6))
		} else {
			log.Printf("⚠️ [WARNING] Failed to get XRP price: %v (continuing with 0)", err)
			xrpPrice = decimal.Zero
		}
	} else {
		log.Printf("❌ [ERROR] Database connection not available for price service")
		return nil, fmt.Errorf("database connection not available for price service")
	}
	xrpPriceFloat, _ := xrpPrice.Float64()
	xrpBalancesResponse.XRPPrice = xrpPriceFloat

	// Get account info for XRP balance
	// Getting account info for XRP balance (debug logging removed)
	accountRequest := xrpl.BaseRequest{
		"id":      1,
		"command": "account_info",
		"account": address,
	}
	// Sending account_info request (debug logging removed)

	accountResponse, err := client.Request(accountRequest)
	if err != nil {
		log.Printf("❌ [ERROR] Failed to request account info: %v", err)
		return nil, fmt.Errorf("failed to request account info: %w", err)
	}
	// Account info response received (debug logging removed)

	// Check for XRPL API errors in account response
	if status, ok := accountResponse["status"].(string); ok && status == "error" {
		if errorMsg, ok := accountResponse["error"].(string); ok {
			log.Printf("❌ [ERROR] XRPL account_info API error: %s", errorMsg)
			return nil, fmt.Errorf("XRPL account_info API error: %s", errorMsg)
		}
		log.Printf("❌ [ERROR] XRPL account_info API returned error status")
		return nil, fmt.Errorf("XRPL account_info API returned error status")
	}
	// Account info response validated (debug logging removed)

	// Parse XRP balance from account info
	var xrpBalance decimal.Decimal
	if result, ok := accountResponse["result"].(map[string]interface{}); ok {
		if accountData, ok := result["account_data"].(map[string]interface{}); ok {
			if balanceStr, ok := accountData["Balance"].(string); ok {
				if balance, err := decimal.NewFromString(balanceStr); err == nil {
					xrpBalance = balance.Div(decimal.NewFromInt(1000000)) // Convert from drops to XRP
					log.Printf("✅ [SUCCESS] XRP balance parsed: %s drops -> %s XRP", balanceStr, xrpBalance.String())
				} else {
					log.Printf("⚠️ [WARNING] Failed to parse XRP balance '%s': %v", balanceStr, err)
				}
			} else {
				log.Printf("⚠️ [WARNING] Balance field not found in account_data")
			}
		} else {
			log.Printf("⚠️ [WARNING] account_data field not found in result")
		}
	} else {
		log.Printf("⚠️ [WARNING] result field not found in response")
	}
	xrpBalanceFloat, _ := xrpBalance.Float64()
	xrpBalancesResponse.XRPBalance = xrpBalanceFloat
	// XRP balance set to response (debug logging removed)

	// Get account lines (trustlines) for token balances
	// Getting account lines for token balances (debug logging removed)
	trustlinesRequest := xrpl.BaseRequest{
		"id":           2,
		"command":      "account_lines",
		"account":      address,
		"ledger_index": "validated",
	}
	// Sending account_lines request (debug logging removed)

	trustlinesResponse, err := client.Request(trustlinesRequest)
	if err != nil {
		log.Printf("❌ [ERROR] Failed to request account lines: %v", err)
		return nil, fmt.Errorf("failed to request account lines: %w", err)
	}
	// Account lines response received (debug logging removed)

	// Parse trustlines response
	// Parsing trustlines response (debug logging removed)
	var linesArray []interface{}
	if result, ok := trustlinesResponse["result"]; ok {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if lines, ok := resultMap["lines"].([]interface{}); ok {
				linesArray = lines
				log.Printf("✅ [SUCCESS] Found %d trustlines in response", len(linesArray))
			} else {
				log.Printf("⚠️ [WARNING] lines field is not an array")
			}
		} else {
			log.Printf("⚠️ [WARNING] result field is not a map")
		}
	} else {
		log.Printf("⚠️ [WARNING] result field not found in trustlines response")
	}

	// Process each trustline
	// Processing trustlines (debug logging removed)
	tokenCount := 0
	for i, lineInterface := range linesArray {
		// Processing trustline (debug logging reduced)

		elem, ok := lineInterface.(map[string]interface{})
		if !ok {
			// Trustline is not a map, skipping (debug logging removed)
			continue
		}

		balanceStr, ok := elem["balance"].(string)
		if !ok || balanceStr == "0" {
			// Trustline has zero or invalid balance, skipping (debug logging removed)
			continue // Skip zero balances
		}

		balanceDecimal, err := decimal.NewFromString(balanceStr)
		if err != nil {
			log.Printf("⚠️ [WARNING] Trustline %d failed to parse balance '%s': %v, skipping", i+1, balanceStr, err)
			continue
		}
		balanceFloat, _ := balanceDecimal.Float64()

		currency, ok := elem["currency"].(string)
		if !ok {
			log.Printf("⚠️ [WARNING] Trustline %d missing currency field, skipping", i+1)
			continue
		}

		issuer, ok := elem["account"].(string)
		if !ok {
			log.Printf("⚠️ [WARNING] Trustline %d missing issuer field, skipping", i+1)
			continue
		}

		// Trustline processing (debug logging removed)

		// Get token name from database
		// Querying database for token name (debug logging removed)
		var tokenName string = ""
		if svc.db != nil {
			// Query for the name field, handling potential byte array format
			var dbToken struct {
				Name []byte `json:"name"`
			}

			err = svc.db.Raw(`
				SELECT token_name FROM xrpTokens 
				WHERE currency = ? AND issuer = ?
				LIMIT 1
			`, currency, issuer).Scan(&dbToken).Error

			if err == nil {
				if dbToken.Name != nil && len(dbToken.Name) > 0 {
					tokenName = string(dbToken.Name)
					log.Printf("✅ [SUCCESS] Token name found in database: %s", tokenName)
				} else {
					log.Printf("⚠️ [WARNING] Token name is empty in database")
				}
			} else {
				log.Printf("⚠️ [WARNING] Database query failed for token name: %v", err)
			}
		}

		// Get token price using live AMM calculation (same as XRP)
		var tokenPrice decimal.Decimal
		var tokenValue decimal.Decimal

		if svc.connectionManager != nil {
			// Getting live AMM price (debug logging removed)

			// Create price service with connection manager only (no database)
			priceService := &PriceService{
				cache:             make(map[string]PriceCache),
				db:                nil, // No database - use direct ledger calls only
				connectionManager: svc.connectionManager,
			}

			// Get token price directly from AMM calls
			tokenPriceInXRP, err := priceService.GetTokenPriceInXRPDirect(currency, issuer)
			if err == nil && tokenPriceInXRP.GreaterThan(decimal.Zero) {
				// Convert to USD using the already-calculated XRP price
				tokenPrice = tokenPriceInXRP.Mul(xrpPrice)
				tokenValue = balanceDecimal.Mul(tokenPrice)
				log.Printf("✅ [SUCCESS] Direct AMM token price: %s XRP × $%s = $%s USD, Value: $%s",
					tokenPriceInXRP.StringFixed(6), xrpPrice.StringFixed(6), tokenPrice.StringFixed(6), tokenValue.StringFixed(2))
			} else {
				log.Printf("⚠️ [WARNING] Direct AMM price calculation failed: %v", err)
				tokenPrice = decimal.Zero
				tokenValue = decimal.Zero
			}
		} else {
			log.Printf("❌ [ERROR] Connection manager not available for direct AMM calls")
			tokenPrice = decimal.Zero
			tokenValue = decimal.Zero
		}

		tokenPriceFloat, _ := tokenPrice.Float64()
		tokenValueFloat, _ := tokenValue.Float64()

		// Add token to response with price and value
		// Adding token to response (debug logging removed)
		xrpBalancesResponse.XRPTokens = append(xrpBalancesResponse.XRPTokens, &models.XRPBalanceElem{
			Issuer:   issuer,
			Balance:  balanceFloat,
			Currency: currency, // XRP Ledger standard currency field
			Price:    tokenPriceFloat,
			Value:    tokenValueFloat,
			Name:     tokenName, // Proper token name from database
		})
		tokenCount++
		// Token added to response successfully (debug logging removed)
	}

	log.Printf("✅ [SUCCESS] GetAccountBalances completed successfully. XRP Balance: %s, Tokens: %d",
		xrpBalance.String(), len(xrpBalancesResponse.XRPTokens))

	return &xrpBalancesResponse, nil
}

// GetAccountTransactions retrieves transaction history for a given XRP address
// CRITICAL: This function requires xrplcluster.com for recent transaction data
func (svc *XRPService) GetAccountTransactions(address string) (interface{}, error) {
	// Basic XRP address format validation (starts with 'r' and is 25-35 characters)
	if len(address) < 25 || len(address) > 35 || !strings.HasPrefix(address, "r") {
		return nil, fmt.Errorf("invalid XRP address format: %s", address)
	}

	// CRITICAL FIX: Force use of xrplcluster.com for recent transaction data
	// wss://xrpl.ws/ doesn't have recent historical data, causing old dates
	var client *xrpl.Client
	var err error
	var usingForcedConnection bool

	if strings.Contains(svc.connectionManager.connectionURL, "xrplcluster.com") {
		// Use existing connection if it's already xrplcluster.com
		client, err = svc.GetClient()
		if err != nil {
			return nil, fmt.Errorf("failed to get XRPL client: %w", err)
		}
		defer svc.connectionManager.ReturnConnection(client)
		usingForcedConnection = false
	} else {
		// Force create a new connection to xrplcluster.com for transactions
		fmt.Printf("⚠️ [CRITICAL] Forcing xrplcluster.com connection for transaction data\n")
		fmt.Printf("⚠️ [CRITICAL] Current connection: %s (will use xrplcluster.com instead)\n", svc.connectionManager.connectionURL)

		// Create a temporary connection manager specifically for transactions
		// Use a safer approach with error handling
		tempConnMgr := NewConnectionManager("wss://xrplcluster.com", 1, svc.connectionManager.logger)
		if tempConnMgr == nil {
			return nil, fmt.Errorf("failed to create temporary connection manager")
		}

		tempClient, err := tempConnMgr.GetConnection()
		if err != nil {
			// Fallback: try to use the original connection instead of crashing
			fmt.Printf("⚠️ [FALLBACK] Failed to create xrplcluster.com connection, using original connection: %v\n", err)
			client, err = svc.GetClient()
			if err != nil {
				return nil, fmt.Errorf("failed to get XRPL client (fallback): %w", err)
			}
			defer svc.connectionManager.ReturnConnection(client)
			usingForcedConnection = false
		} else {
			defer tempConnMgr.ReturnConnection(tempClient)
			client = tempClient
			usingForcedConnection = true
		}
	}

	// DEBUG: Log the connection details
	// DEBUG logging removed XRPL Connection Details:\n")
	fmt.Printf("  Connection URL: %s\n", svc.connectionManager.connectionURL)
	fmt.Printf("  Client ID: %v\n", client)

	// Show which connection we're actually using for transactions
	if usingForcedConnection {
		fmt.Printf("✅ [CONNECTION] Using forced xrplcluster.com connection for transactions\n")
	} else if strings.Contains(svc.connectionManager.connectionURL, "xrplcluster.com") {
		fmt.Printf("✅ [CONNECTION] Using existing xrplcluster.com connection\n")
	} else {
		fmt.Printf("⚠️ [CONNECTION] Using fallback connection (may have old data)\n")
	}

	// OPTIMAL FIX: Use ALL LEDGERS (-1,-1) for complete transaction history
	// Testing proved this is FASTER (320ms vs 785ms) and more comprehensive than any fixed range
	// XRPL automatically handles pagination and returns all available transaction history

	// fmt.Printf("📊 [ALL LEDGERS] Using (-1,-1) for complete transaction history (FASTER than limited ranges!)\n")

	transactionRequest := xrpl.BaseRequest{
		"id":               2,
		"command":          "account_tx",
		"account":          address,
		"ledger_index_min": -1, // FROM earliest available (complete history)
		"ledger_index_max": -1, // TO latest available (always current)
		"binary":           false,
		"forward":          false, // Get most recent transactions first (reverse chronological)
		"limit":            200,   // Reasonable limit for recent transactions
	}

	// DEBUG: Log the exact request being sent
	// DEBUG logging removed XRPL Request Details:\n")
	// fmt.Printf("  Command: %s\n", transactionRequest["command"])
	// fmt.Printf("  Account: %s\n", transactionRequest["account"])
	// fmt.Printf("  Ledger Range: %v to %v\n", transactionRequest["ledger_index_min"], transactionRequest["ledger_index_max"])
	// fmt.Printf("  Limit: %v\n", transactionRequest["limit"])
	// fmt.Printf("  Forward: %v\n", transactionRequest["forward"])
	// fmt.Printf("  Binary: %v\n", transactionRequest["binary"])

	response, err := client.Request(transactionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to request account transactions: %w", err)
	}

	// DEBUG: Log the raw response structure
	// DEBUG logging removed XRPL Raw Response:\n")
	// fmt.Printf("  Response type: %T\n", response)
	// fmt.Printf("  Response keys: %v\n", getMapKeys(response))

	// Check for XRPL API errors in response
	if status, ok := response["status"].(string); ok && status == "error" {
		if errorMsg, ok := response["error"].(string); ok {
			return nil, fmt.Errorf("XRPL API error: %s", errorMsg)
		}
		return nil, fmt.Errorf("XRPL API returned error status")
	}

	// DEBUG: Check if we have transactions in the response
	if result, ok := response["result"]; ok {
		if resultMap, ok := result.(map[string]interface{}); ok {
			// DEBUG logging removed XRPL Result Structure:\n")
			// fmt.Printf("  Result keys: %v\n", getMapKeys(resultMap))

			if txs, ok := resultMap["transactions"]; ok {
				if _, ok := txs.([]interface{}); ok {
					// DEBUG logging removed Found %d transactions in response\n", len(txArray))

					// Show details of first few transactions
					// for i, tx := range txArray {
					// 	if i >= 3 { // Only show first 3
					// 		break
					// 	}
					// 	if txMap, ok := tx.(map[string]interface{}); ok {
					// 		fmt.Printf("  Transaction %d:\n", i+1)
					// 		fmt.Printf("    Keys: %v\n", getMapKeys(txMap))

					// 		// Check for 'tx' object
					// 		if txObj, ok := txMap["tx"]; ok {
					// 			if txData, ok := txObj.(map[string]interface{}); ok {
					// 				fmt.Printf("    TX object keys: %v\n", getMapKeys(txData))

					// 				// Show key transaction fields
					// 				for _, key := range []string{"Account", "Destination", "TransactionType", "date", "ledger_index"} {
					// 					if val, exists := txData[key]; exists {
					// 						fmt.Printf("    %s: %v (type: %T)\n", key, val, val)
					// 					}
					// 				}
					// 			}
					// 		}

					// 		// Check for 'date' field
					// 		if date, ok := txMap["date"]; ok {
					// 			fmt.Printf("    Date: %v (type: %T)\n", date, date)
					// 		}

					// 		// Check for 'ledger_index' field
					// 		if ledgerIndex, ok := txMap["ledger_index"]; ok {
					// 			fmt.Printf("    Ledger Index: %v (type: %T)\n", ledgerIndex, ledgerIndex)
					// 		}
					// 	}
					// }
				}
			}
		}
	}

	return response, nil
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (svc *XRPService) GetDefiData(address string) (string, error) {
	_, err := svc.GetAccountTransactions(address)
	if err != nil {
		return "", err
	}

	_, err = svc.GetAccountBalances(address)
	if err != nil {
		return "", err
	}

	return "success", nil
}
