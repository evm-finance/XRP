package xrp

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	dbpackage "qc-defi-graphql-server/internal/database"
	"qc-defi-graphql-server/internal/models"
)

func TestNewXRPTokenService(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	if tokenService == nil {
		t.Fatal("Expected token service to be created, got nil")
	}
}

func TestFetchTokensFromAPI(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	tokens, err := tokenService.FetchTokensFromAPI()
	if err != nil {
		t.Fatalf("Failed to fetch tokens from API: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("Expected tokens to be fetched, got empty slice")
	}

	// Verify token structure
	for i, token := range tokens {
		if token.Currency == "" {
			t.Errorf("Token %d has empty currency", i)
		}
		if token.Issuer == "" {
			t.Errorf("Token %d has empty issuer", i)
		}
	}
}

func TestStoreTokens(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	// First, ensure the database and tables exist
	err := dbService.CreateXRPDatabase()
	if err != nil {
		t.Fatalf("Failed to create XRP database: %v", err)
	}

	if concrete, ok := dbService.(*dbpackage.DatabaseCreateImpl); ok {
		concrete.CreateXRPTokensTable()
	} else {
		t.Log("Could not call CreateXRPTokensTable: not a *DatabaseCreateImpl")
	}

	// Fetch tokens from API
	tokens, err := tokenService.FetchTokensFromAPI()
	if err != nil {
		t.Fatalf("Failed to fetch tokens: %v", err)
	}

	// Store only first 5 tokens for testing
	if len(tokens) > 5 {
		tokens = tokens[:5]
	}

	err = tokenService.StoreTokens(tokens)
	if err != nil {
		t.Fatalf("Failed to store tokens: %v", err)
	}
}

func TestGetTokensWithPools(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	// Ensure database and tables exist
	err := dbService.CreateXRPDatabase()
	if err != nil {
		t.Fatalf("Failed to create XRP database: %v", err)
	}

	if concrete, ok := dbService.(*dbpackage.DatabaseCreateImpl); ok {
		concrete.CreateXRPTokensTable()
		concrete.CreateXRPAMMTable()
	} else {
		t.Log("Could not call CreateXRPTokensTable/CreateXRPAMMTable: not a *DatabaseCreateImpl")
	}

	// This test will likely return empty results since we haven't populated AMM data yet
	tokens, err := tokenService.GetTokensWithPools()
	if err != nil {
		t.Fatalf("Failed to get tokens with pools: %v", err)
	}

	// For now, we just verify the function doesn't crash
	_ = tokens
}

func TestUpdateTokenMetrics(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	// Ensure database and tables exist
	err := dbService.CreateXRPDatabase()
	if err != nil {
		t.Fatalf("Failed to create XRP database: %v", err)
	}

	if concrete, ok := dbService.(*dbpackage.DatabaseCreateImpl); ok {
		concrete.CreateXRPTokensTable()
	} else {
		t.Log("Could not call CreateXRPTokensTable: not a *DatabaseCreateImpl")
	}

	err = tokenService.UpdateTokenMetrics()
	if err != nil {
		t.Fatalf("Failed to update token metrics: %v", err)
	}
}

func TestTokenPaginationDiscovery(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	t.Log("Testing token pagination discovery...")

	// Test different page sizes
	pageSizes := []int{50, 100, 200}

	for _, limit := range pageSizes {
		t.Logf("\n--- Testing with limit: %d ---", limit)

		// Discover total pages
		totalPages, err := tokenService.DiscoverTotalPages(limit)
		if err != nil {
			t.Logf("❌ Failed to discover pages with limit %d: %v", limit, err)
			continue
		}

		t.Logf("✅ Discovered %d pages with limit %d", totalPages, limit)

		// Fetch first page to see how many tokens we get
		firstPageTokens, err := tokenService.FetchTokensWithPagination(1, limit)
		if err != nil {
			t.Logf("❌ Failed to fetch first page with limit %d: %v", limit, err)
			continue
		}

		t.Logf("✅ First page: %d tokens", len(firstPageTokens))

		// Try to fetch the last page
		if totalPages > 1 {
			lastPageTokens, err := tokenService.FetchTokensWithPagination(totalPages, limit)
			if err != nil {
				t.Logf("❌ Failed to fetch last page (%d) with limit %d: %v", totalPages, limit, err)
			} else {
				t.Logf("✅ Last page (%d): %d tokens", totalPages, len(lastPageTokens))
			}
		}

		// Calculate estimated total tokens
		estimatedTotal := totalPages * limit
		t.Logf("📊 Estimated total tokens: %d (pages: %d × limit: %d)", estimatedTotal, totalPages, limit)
	}
}

func TestFetchAllTokensWithPagination(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	t.Log("Testing fetch all tokens with pagination and incremental storage...")

	// Fetch all tokens using pagination (tokens are stored incrementally)
	allTokens, err := tokenService.FetchAllTokensWithPagination()
	if err != nil {
		t.Fatalf("Failed to fetch all tokens: %v", err)
	}

	t.Logf("✅ Successfully fetched and stored %d total tokens using pagination", len(allTokens))

	// Show some sample tokens
	if len(allTokens) > 0 {
		t.Logf("\n📋 Sample tokens:")
		for i := 0; i < 5 && i < len(allTokens); i++ {
			token := allTokens[i]
			t.Logf("  %d. %s (%s) - Issuer: %s",
				i+1, token.Meta.Token.Name, token.Currency, token.Issuer)
		}
	}

	// Verify tokens are actually stored in database
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	var storedTokenCount int
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&storedTokenCount).Error
	if err != nil {
		t.Fatalf("Failed to count stored tokens: %v", err)
	}

	t.Logf("💾 Tokens stored in database: %d", storedTokenCount)

	// Verify we have a reasonable number of tokens
	if storedTokenCount < 100 {
		t.Logf("Warning: Only %d tokens stored, expected many more", storedTokenCount)
	} else {
		t.Logf("✅ Successfully stored %d tokens in database", storedTokenCount)
	}
}

func TestCheckStoredTokenCount(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	var storedTokenCount int
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&storedTokenCount).Error
	if err != nil {
		t.Fatalf("Failed to count stored tokens: %v", err)
	}

	t.Logf("💾 Tokens currently stored in database: %d", storedTokenCount)

	// Show some sample tokens
	var sampleTokens []models.XRPTokenData
	err = db.Limit(5).Find(&sampleTokens).Error
	if err != nil {
		t.Fatalf("Failed to fetch sample tokens: %v", err)
	}

	if len(sampleTokens) > 0 {
		t.Logf("\n📋 Sample stored tokens:")
		for i, token := range sampleTokens {
			t.Logf("  %d. %s (%s) - Issuer: %s",
				i+1, token.Meta.Token.Name, token.Currency, token.Issuer)
		}
	}
}

func TestResumeTokenCollectionFromPage39(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	t.Log("Testing resume token collection from page 39...")

	// Resume from page 39 (where we left off)
	allTokens, err := tokenService.FetchAllTokensWithPaginationFromPage(39)
	if err != nil {
		t.Fatalf("Failed to resume token collection: %v", err)
	}

	t.Logf("✅ Successfully resumed and collected %d additional tokens", len(allTokens))

	// Verify total tokens in database
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	var storedTokenCount int
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&storedTokenCount).Error
	if err != nil {
		t.Fatalf("Failed to count stored tokens: %v", err)
	}

	t.Logf("💾 Total tokens now stored in database: %d", storedTokenCount)
}

func TestResumeTokenCollectionFromPage77(t *testing.T) {
	dbService := dbpackage.NewDatabaseCreateService()
	tokenService := NewXRPTokenService(dbService)

	t.Log("Testing resume token collection from page 77...")

	// Resume from page 77 (where we left off)
	allTokens, err := tokenService.FetchAllTokensWithPaginationFromPage(77)
	if err != nil {
		t.Fatalf("Failed to resume token collection: %v", err)
	}

	t.Logf("✅ Successfully resumed and collected %d additional tokens", len(allTokens))

	// Verify total tokens in database
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	var storedTokenCount int
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&storedTokenCount).Error
	if err != nil {
		t.Fatalf("Failed to count stored tokens: %v", err)
	}

	t.Logf("💾 Total tokens now stored in database: %d", storedTokenCount)
}

func TestTokenDiscoveryAndDatabaseDump(t *testing.T) {
	// Create a fresh database service and connection for this test only
	dbService := dbpackage.NewDatabaseCreateService()
	db, err := dbService.GetConnectionXRPDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	// Ensure the xrpTokens table exists and is empty for this test
	dbService.DropXRPTokensTable()
	dbService.CreateXRPTokensTable()

	tokenService := NewXRPTokenService(dbService)

	t.Log("🧪 Testing token discovery for 100 tokens and database dump...")

	// Step 1: Fetch 100 tokens from API
	t.Log("📡 Fetching 100 tokens from API...")
	tokens, err := tokenService.FetchTokensWithPagination(1, 100)
	if err != nil {
		t.Fatalf("Failed to fetch tokens: %v", err)
	}

	t.Logf("✅ Fetched %d tokens from API", len(tokens))

	// Show sample tokens before storage
	t.Log("📋 Sample tokens from API:")
	for i := 0; i < 5 && i < len(tokens); i++ {
		token := tokens[i]
		t.Logf("  %d. Currency: '%s' (type: %T), Issuer: '%s' (type: %T), Name: '%s'",
			i+1, token.Currency, token.Currency, token.Issuer, token.Issuer, token.Meta.Token.Name)
	}

	// Step 2: Store tokens in database
	t.Log("💾 Storing tokens in database...")
	err = tokenService.StoreTokens(tokens)
	if err != nil {
		t.Fatalf("Failed to store tokens: %v", err)
	}

	t.Logf("✅ Successfully stored %d tokens in database", len(tokens))

	// Step 3: Dump database to verify storage
	t.Log("🔍 Dumping database to verify storage...")

	// Get total count
	var totalCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&totalCount).Error
	if err != nil {
		t.Fatalf("Failed to count tokens: %v", err)
	}

	t.Logf("📊 Total tokens in database: %d", totalCount)

	// Get sample tokens from database
	var dbTokens []struct {
		ID          uint
		Currency    interface{}
		Issuer      string
		Name        string
		Icon        string
		Description string
		Trustlines  int
		Holders     int
		Supply      string
		Marketcap   float64
		Price       float64
		Volume24H   float64
		Volume7D    string
	}

	err = db.Raw(`
		SELECT id, currency, issuer, name, icon, description, 
		       trustlines, holders, supply, marketcap, price, volume24h, volume7d
		FROM xrpTokens 
		ORDER BY id DESC 
		LIMIT 20
	`).Scan(&dbTokens).Error
	if err != nil {
		t.Fatalf("Failed to query tokens from database: %v", err)
	}

	t.Logf("📋 Database dump (last 20 tokens):")
	for i, token := range dbTokens {
		// Check currency type and value
		var currencyStr string
		var currencyType string
		switch v := token.Currency.(type) {
		case string:
			currencyStr = v
			currencyType = "string"
		case []byte:
			currencyStr = string(v)
			currencyType = "[]byte"
		default:
			currencyStr = fmt.Sprintf("%v", v)
			currencyType = fmt.Sprintf("%T", v)
		}

		t.Logf("  %d. ID: %d, Currency: '%s' (type: %s), Issuer: '%s', Name: '%s'",
			i+1, token.ID, currencyStr, currencyType, token.Issuer, token.Name)
	}

	// Step 4: Verify data types
	t.Log("🔍 Verifying data types...")
	var stringCount, byteCount, otherCount int64

	// Count by currency type
	rows, err := db.Raw("SELECT currency FROM xrpTokens LIMIT 100").Rows()
	if err != nil {
		t.Fatalf("Failed to query currency types: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var currencyInterface interface{}
		if err := rows.Scan(&currencyInterface); err != nil {
			continue
		}

		switch currencyInterface.(type) {
		case string:
			stringCount++
		case []byte:
			byteCount++
		default:
			otherCount++
		}
	}

	t.Logf("📊 Currency data types in database:")
	t.Logf("  String: %d", stringCount)
	t.Logf("  Byte array: %d", byteCount)
	t.Logf("  Other: %d", otherCount)

	// Step 5: Test currency filtering
	t.Log("🔍 Testing currency filtering...")
	var filteredCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency IS NOT NULL AND currency != ''").Scan(&filteredCount).Error
	if err != nil {
		t.Fatalf("Failed to count filtered tokens: %v", err)
	}

	t.Logf("📊 Tokens with non-null, non-empty currency: %d", filteredCount)

	// Step 6: Test specific currency queries
	t.Log("🔍 Testing specific currency queries...")
	var xrpCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens WHERE currency = 'XRP'").Scan(&xrpCount).Error
	if err != nil {
		t.Logf("⚠️  Could not count XRP tokens: %v", err)
	} else {
		t.Logf("📊 XRP tokens: %d", xrpCount)
	}

	// Test for 3-character currency codes (common pattern)
	var threeCharCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens WHERE LENGTH(currency) = 3").Scan(&threeCharCount).Error
	if err != nil {
		t.Logf("⚠️  Could not count 3-character currencies: %v", err)
	} else {
		t.Logf("📊 3-character currency codes: %d", threeCharCount)
	}

	t.Log("✅ Token discovery and database dump test completed!")
}

func TestFetchSecondPageOfTokens_APIOnly(t *testing.T) {
	tokenService := &XRPTokenService{
		dbService: nil, // Not used for API fetch
		client:    &http.Client{Timeout: 30 * time.Second},
		apiURL:    "https://s1.xrplmeta.org",
	}

	pageSize := 100

	t.Log("Fetching first page of tokens...")
	firstPage, err := tokenService.FetchTokensWithPagination(1, pageSize)
	if err != nil {
		t.Fatalf("Failed to fetch first page: %v", err)
	}
	t.Logf("First page: %d tokens", len(firstPage))
	for i := 0; i < 5 && i < len(firstPage); i++ {
		t.Logf("  %d. %s (%s) - Issuer: %s", i+1, firstPage[i].Meta.Token.Name, firstPage[i].Currency, firstPage[i].Issuer)
	}

	t.Log("Fetching second page of tokens...")
	secondPage, err := tokenService.FetchTokensWithPagination(2, pageSize)
	if err != nil {
		t.Fatalf("Failed to fetch second page: %v", err)
	}
	t.Logf("Second page: %d tokens", len(secondPage))
	for i := 0; i < 5 && i < len(secondPage); i++ {
		t.Logf("  %d. %s (%s) - Issuer: %s", i+1, secondPage[i].Meta.Token.Name, secondPage[i].Currency, secondPage[i].Issuer)
	}

	if len(firstPage) > 0 && len(secondPage) > 0 && firstPage[0].Currency == secondPage[0].Currency && firstPage[0].Issuer == secondPage[0].Issuer {
		t.Errorf("First token on page 1 and 2 are the same, pagination may not be working!")
	}
}

func TestFetchLargeTokenLimit(t *testing.T) {
	tokenService := &XRPTokenService{
		dbService: nil, // Not used for API fetch
		client:    &http.Client{Timeout: 30 * time.Second},
		apiURL:    "https://s1.xrplmeta.org",
	}

	largeLimit := 10000

	t.Logf("Fetching tokens with large limit: %d", largeLimit)
	tokens, err := tokenService.FetchTokensWithPagination(1, largeLimit)
	if err != nil {
		t.Fatalf("Failed to fetch tokens with large limit: %v", err)
	}

	t.Logf("Total tokens returned: %d", len(tokens))

	// Show first 10 tokens
	t.Log("First 10 tokens:")
	for i := 0; i < 10 && i < len(tokens); i++ {
		t.Logf("  %d. %s (%s) - Issuer: %s", i+1, tokens[i].Meta.Token.Name, tokens[i].Currency, tokens[i].Issuer)
	}

	// Show last 10 tokens if we have more than 10
	if len(tokens) > 10 {
		t.Log("Last 10 tokens:")
		start := len(tokens) - 10
		for i := start; i < len(tokens); i++ {
			t.Logf("  %d. %s (%s) - Issuer: %s", i+1, tokens[i].Meta.Token.Name, tokens[i].Currency, tokens[i].Issuer)
		}
	}

	// Check if we got the full limit or fewer tokens
	if len(tokens) == largeLimit {
		t.Logf("⚠️  Got exactly %d tokens - there might be more available", largeLimit)
	} else if len(tokens) < largeLimit {
		t.Logf("✅ Got %d tokens (less than limit %d) - this might be all available tokens", len(tokens), largeLimit)
	}
}
