package xrp

import (
	"fmt"
	"os"
	"qc-defi-graphql-server/configs"
	"strings"
	"testing"
	"time"
)

var xrpService *XRPService

func init() {
	xrpService = &XRPService{}
}

func TestGetTokenList(t *testing.T) {
	tokens, err := xrpService.GetTokenList()
	fmt.Println(tokens, err)

	for index, value := range tokens {
		fmt.Println(index, value)
	}
}

func TestGetAccountTransactions_Newest(t *testing.T) {
	service := XRPService{}
	response, err := service.GetAccountTransactions("rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj")
	if err != nil {
		t.Errorf("Error getting account transactions: %v", err)
	}
	fmt.Printf("Account transactions response: %+v\n", response)
}

func TestGetAccountBalances_Newest(t *testing.T) {
	service := XRPService{}
	balances, err := service.GetAccountBalances("rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj")
	if err != nil {
		t.Errorf("Error getting account balances: %v", err)
	}
	fmt.Printf("Account balances: %+v\n", balances)
}

func TestConfigDebug(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	// Print the raw connection strings loaded from env
	t.Logf("Loaded DeFiDB DSN: %s", conf.EnvVars().DeFiDB)
	t.Logf("Loaded XRPDB DSN: %s", conf.EnvVars().XRPDB)
}

func TestTokenDatabaseCount(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Count total tokens
	var tokenCount int
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&tokenCount).Error
	if err != nil {
		t.Fatalf("Failed to count tokens: %v", err)
	}

	t.Logf("Total tokens in database: %d", tokenCount)

	// Get sample tokens
	var sampleTokens []struct {
		Issuer    string  `json:"issuer"`
		Currency  string  `json:"currency"`
		Name      string  `json:"name"`
		Marketcap float64 `json:"marketCap"`
		Price     float64 `json:"price"`
		Volume24H float64 `json:"volume24H"`
	}

	err = dbConn.Raw("SELECT issuer, currency, name, marketcap, price, volume24h FROM xrpTokens LIMIT 10").Scan(&sampleTokens).Error
	if err != nil {
		t.Fatalf("Failed to get sample tokens: %v", err)
	}

	t.Logf("Sample tokens:")
	for i, token := range sampleTokens {
		t.Logf("  %d. %s (%s) - Marketcap: $%.2f, Price: $%.6f, Volume24H: $%.2f",
			i+1, token.Name, token.Currency, token.Marketcap, token.Price, token.Volume24H)
	}
}

func TestAMMDatabaseCount(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Check if xrpAmm table exists
	var tableExists int
	err = dbConn.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'XRP' AND table_name = 'xrpAmm'").Scan(&tableExists).Error
	if err != nil {
		t.Fatalf("Failed to check if xrpAmm table exists: %v", err)
	}

	if tableExists == 0 {
		t.Log("xrpAmm table does not exist yet")
		return
	}

	// Count total AMM pools
	var ammCount int
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpAmm_normalized").Scan(&ammCount).Error
	if err != nil {
		t.Fatalf("Failed to count AMM pools: %v", err)
	}

	t.Logf("Total AMM pools in database: %d", ammCount)

	// Print actual AMM pool contents using the current schema
	t.Logf("Sample AMM pools (current schema):")
	type PoolRow struct {
		Account        string
		Asset1Amount   string
		Asset1Currency string
		Asset1Issuer   string
		Asset2Amount   string
		TradingFee     int
		CreatedAt      int64
		LastUpdated    int64
	}
	var samplePools []PoolRow
	err = dbConn.Raw(`SELECT account, asset1_amount, asset1_currency, asset1_issuer, asset2_amount, trading_fee, created_at, last_updated FROM xrpAmm_normalized LIMIT 10`).Scan(&samplePools).Error
	if err != nil {
		t.Fatalf("Failed to get sample AMM pools: %v", err)
	}
	for i, pool := range samplePools {
		t.Logf("  %d. Account: %s | Asset1: %s | Currency: %s | Issuer: %s | Asset2: %s | Fee: %d | Created: %d | Updated: %d",
			i+1, pool.Account, pool.Asset1Amount, pool.Asset1Currency, pool.Asset1Issuer, pool.Asset2Amount, pool.TradingFee, pool.CreatedAt, pool.LastUpdated)
	}
}

func TestAMMCollection(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Get unique token pairs from xrpAmm_normalized table
	var uniquePairs []struct {
		Asset1Currency string `json:"asset1_currency"`
		Asset1Issuer   string `json:"asset1_issuer"`
		Asset2Currency string `json:"asset2_currency"`
		Asset2Issuer   string `json:"asset2_issuer"`
		Count          int    `json:"count"`
	}

	err = dbConn.Raw(`
		SELECT 
			asset1_currency, asset1_issuer, asset2_currency, asset2_issuer, COUNT(*) as count
		FROM xrpAmm_normalized 
		GROUP BY asset1_currency, asset1_issuer, asset2_currency, asset2_issuer
		ORDER BY count DESC
		LIMIT 20
	`).Scan(&uniquePairs).Error

	if err != nil {
		t.Fatalf("Failed to get unique AMM pairs: %v", err)
	}

	t.Logf("Unique AMM pairs (top 20):")
	for i, pair := range uniquePairs {
		t.Logf("  %d. %s/%s - %s/%s (Count: %d)",
			i+1, pair.Asset1Currency, pair.Asset1Issuer, pair.Asset2Currency, pair.Asset2Issuer, pair.Count)
	}
}

func TestDumpAllXRPTokens(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Create log file
	logFile, err := os.Create("xrp_tokens_dump.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Write header to log
	fmt.Fprintf(logFile, "XRP Tokens Database Dump - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(logFile, "%s\n\n", strings.Repeat("=", 80))

	// Get all tokens
	rows, err := dbConn.Raw("SELECT * FROM xrpTokens").Rows()
	if err != nil {
		t.Fatalf("Failed to query tokens: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("Failed to get columns: %v", err)
	}

	// Write column headers
	fmt.Fprintf(logFile, "Columns: %v\n\n", columns)

	// Prepare scan variables
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Write each row
	rowCount := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			fmt.Fprintf(logFile, "Error scanning row %d: %v\n", rowCount+1, err)
			continue
		}

		rowCount++
		fmt.Fprintf(logFile, "Row %d:\n", rowCount)
		for i, col := range columns {
			val := values[i]
			if val == nil {
				fmt.Fprintf(logFile, "  %s: NULL\n", col)
			} else {
				fmt.Fprintf(logFile, "  %s: %v\n", col, val)
			}
		}
		fmt.Fprintf(logFile, "\n")
	}

	fmt.Fprintf(logFile, "Total rows dumped: %d\n", rowCount)
	t.Logf("✅ Dumped %d tokens to xrp_tokens_dump.log", rowCount)
}

func TestDumpAllAMMPools(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Create log file
	logFile, err := os.Create("xrp_amm_pools_dump.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Write header to log
	fmt.Fprintf(logFile, "XRP AMM Pools Database Dump - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(logFile, "%s\n\n", strings.Repeat("=", 80))

	// Get all AMM pools
	rows, err := dbConn.Raw("SELECT * FROM xrpAmm_normalized").Rows()
	if err != nil {
		t.Fatalf("Failed to query AMM pools: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("Failed to get columns: %v", err)
	}

	// Write column headers
	fmt.Fprintf(logFile, "Columns: %v\n\n", columns)

	// Prepare scan variables
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Write each row
	rowCount := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			fmt.Fprintf(logFile, "Error scanning row %d: %v\n", rowCount+1, err)
			continue
		}

		rowCount++
		fmt.Fprintf(logFile, "Row %d:\n", rowCount)
		for i, col := range columns {
			val := values[i]
			if val == nil {
				fmt.Fprintf(logFile, "  %s: NULL\n", col)
			} else {
				fmt.Fprintf(logFile, "  %s: %v\n", col, val)
			}
		}
		fmt.Fprintf(logFile, "\n")
	}

	fmt.Fprintf(logFile, "Total rows dumped: %d\n", rowCount)
	t.Logf("✅ Dumped %d AMM pools to xrp_amm_pools_dump.log", rowCount)
}

func TestFixCurrencyDataType(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// First, let's check what data types we have in the currency field
	var sampleTokens []struct {
		ID       uint
		Currency interface{} // Use interface{} to catch any type
		Issuer   string
	}

	err = dbConn.Raw("SELECT id, currency, issuer FROM xrpTokens LIMIT 10").Scan(&sampleTokens).Error
	if err != nil {
		t.Fatalf("Failed to query sample tokens: %v", err)
	}

	fmt.Printf("Sample tokens currency types:\n")
	for _, token := range sampleTokens {
		fmt.Printf("ID: %d, Currency: %v (type: %T), Issuer: %s\n",
			token.ID, token.Currency, token.Currency, token.Issuer)
	}

	// Check for byte array currencies
	var byteArrayCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND LENGTH(currency) = 3").Scan(&byteArrayCount).Error
	if err != nil {
		t.Fatalf("Failed to count byte array currencies: %v", err)
	}

	fmt.Printf("Tokens with 3-character currency codes: %d\n", byteArrayCount)

	// Check for empty/null currencies
	var nullCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NULL OR currency = ''").Scan(&nullCount).Error
	if err != nil {
		t.Fatalf("Failed to count null currencies: %v", err)
	}

	fmt.Printf("Tokens with null/empty currency: %d\n", nullCount)

	// Check for valid string currencies
	var validCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND currency != '' AND LENGTH(currency) > 3").Scan(&validCount).Error
	if err != nil {
		t.Fatalf("Failed to count valid currencies: %v", err)
	}

	fmt.Printf("Tokens with valid string currency: %d\n", validCount)

	// If we have byte array currencies, convert them to hex strings
	if byteArrayCount > 0 {
		fmt.Printf("Converting byte array currencies to hex strings...\n")

		// Get all tokens with 3-character currency codes
		var byteArrayTokens []struct {
			ID       uint
			Currency []byte
			Issuer   string
		}

		err = dbConn.Raw("SELECT id, currency, issuer FROM xrpTokens WHERE currency IS NOT NULL AND LENGTH(currency) = 3").Scan(&byteArrayTokens).Error
		if err != nil {
			t.Fatalf("Failed to query byte array tokens: %v", err)
		}

		// Convert and update each token
		for _, token := range byteArrayTokens {
			hexCurrency := fmt.Sprintf("%X", token.Currency)
			fmt.Printf("Converting ID %d: %v -> %s\n", token.ID, token.Currency, hexCurrency)

			err = dbConn.Exec("UPDATE xrpTokens SET currency = ? WHERE id = ?", hexCurrency, token.ID).Error
			if err != nil {
				t.Fatalf("Failed to update token %d: %v", token.ID, err)
			}
		}

		fmt.Printf("Successfully converted %d byte array currencies to hex strings\n", len(byteArrayTokens))
	}

	// Verify the fix worked
	var finalValidCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND currency != ''").Scan(&finalValidCount).Error
	if err != nil {
		t.Fatalf("Failed to count final valid currencies: %v", err)
	}

	fmt.Printf("Final count of tokens with valid currency: %d\n", finalValidCount)
}

func TestCheckCurrencyDataTypes(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	dbConn := conf.XRPDB()

	// Check total count
	var totalCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&totalCount).Error
	if err != nil {
		t.Fatalf("Failed to count total tokens: %v", err)
	}
	fmt.Printf("Total tokens in database: %d\n", totalCount)

	// Check count with the same query as AMM discovery
	var filteredCount int64
	err = dbConn.Raw("SELECT COUNT(DISTINCT currency, issuer) FROM xrpTokens WHERE currency IS NOT NULL AND currency != ''").Scan(&filteredCount).Error
	if err != nil {
		t.Fatalf("Failed to count filtered tokens: %v", err)
	}
	fmt.Printf("Tokens with non-null, non-empty currency: %d\n", filteredCount)

	// Check for different currency data types
	var sampleTokens []struct {
		ID       uint
		Currency interface{}
		Issuer   string
	}

	err = dbConn.Raw("SELECT id, currency, issuer FROM xrpTokens LIMIT 20").Scan(&sampleTokens).Error
	if err != nil {
		t.Fatalf("Failed to query sample tokens: %v", err)
	}

	fmt.Printf("\nSample tokens:\n")
	for _, token := range sampleTokens {
		fmt.Printf("ID: %d, Currency: %v (type: %T), Issuer: %s\n",
			token.ID, token.Currency, token.Currency, token.Issuer)
	}

	// Check for specific patterns
	var nullCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NULL").Scan(&nullCount).Error
	if err != nil {
		t.Fatalf("Failed to count null currencies: %v", err)
	}
	fmt.Printf("\nTokens with NULL currency: %d\n", nullCount)

	var emptyCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency = ''").Scan(&emptyCount).Error
	if err != nil {
		t.Fatalf("Failed to count empty currencies: %v", err)
	}
	fmt.Printf("Tokens with empty currency: %d\n", emptyCount)

	// Check for byte arrays (3-character length)
	var byteArrayCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND LENGTH(currency) = 3").Scan(&byteArrayCount).Error
	if err != nil {
		t.Fatalf("Failed to count byte array currencies: %v", err)
	}
	fmt.Printf("Tokens with 3-character currency (likely byte arrays): %d\n", byteArrayCount)

	// Check for hex strings (6-character length)
	var hexCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND LENGTH(currency) = 6").Scan(&hexCount).Error
	if err != nil {
		t.Fatalf("Failed to count hex currencies: %v", err)
	}
	fmt.Printf("Tokens with 6-character currency (likely hex): %d\n", hexCount)

	// Check for longer strings
	var longStringCount int64
	err = dbConn.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND LENGTH(currency) > 6").Scan(&longStringCount).Error
	if err != nil {
		t.Fatalf("Failed to count long string currencies: %v", err)
	}
	fmt.Printf("Tokens with long currency strings: %d\n", longStringCount)
}
