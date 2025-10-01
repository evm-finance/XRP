package xrp

import (
	"qc-defi-graphql-server/configs"
	"testing"
)

// TestTokenEndpointDebug reproduces token endpoint issues across multiple tokens
func TestTokenEndpointDebug(t *testing.T) {
	t.Log("🔍 Debugging Token Endpoint Issues (Multiple Tokens)")
	t.Log("=================================================")

	// Load configuration
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("❌ Failed to load config: %v", err)
	}

	// Get database connection
	xrpDB := conf.XRPDB()
	if xrpDB == nil {
		t.Fatalf("❌ Failed to get XRP database connection")
	}

	// Test multiple tokens that are causing issues
	testTokens := []struct {
		name     string
		currency string
		issuer   string
	}{
		{"ARMY Token", "ARMY", "rGG3wQ4kUzd7Jnmk1n5NWPZjjut62kCBfC"},
		{"MAG Token", "MAG", "rXmagwMmnFtVet3uL26Q2iwk287SRvVMJ"},
		{"USD Token", "USD", "rhub8VRN55s94qWKDv6jmDy1pUykJzF3wq"},
	}

	for _, token := range testTokens {
		t.Run(token.name, func(t *testing.T) {
			testCurrency := token.currency
			testIssuer := token.issuer

			t.Logf("🔍 Testing token: %s/%s", testCurrency, testIssuer)

			// First, check if this token exists in our database
			var tokenExists bool
			err := xrpDB.Raw(`
				SELECT EXISTS(
					SELECT 1 FROM xrpTokens 
					WHERE currency = ? AND issuer = ?
				) as token_exists
			`, testCurrency, testIssuer).Scan(&tokenExists).Error

			if err != nil {
				t.Errorf("❌ Database query to check token existence failed: %v", err)
				return
			}

			t.Logf("📊 Token %s/%s exists in database: %v", testCurrency, testIssuer, tokenExists)

			// Test the exact query from the endpoint handler that's failing
			var token struct {
				Currency    string  `json:"currency"`
				Issuer      string  `json:"issuer"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Supply      string  `json:"supply"`
				Trustlines  int64   `json:"trustlines"`
				Holders     int64   `json:"holders"`
				Price       float64 `json:"price"`
				MarketCap   float64 `json:"marketCap"`
				Volume24h   float64 `json:"volume24H"`
				Volume7d    float64 `json:"volume7d"`
				Icon        string  `json:"icon"`
				Weblinks    string  `json:"weblinks"`
			}

			// This is the FIXED query from the endpoint handler (now using only existing columns)
			query := `
				SELECT currency, issuer, token_name as name, '' as description, supply_xrpl as supply, trustlines, 0 as holders, 
					   price, marketcap, volume_24h as volume24h, volume_24h as volume7d, icon, '' as weblinks
				FROM xrpTokens 
				WHERE currency = ? AND issuer = ?
			`

			t.Log("🔍 Running fixed endpoint query...")
			err = xrpDB.Raw(query, testCurrency, testIssuer).Scan(&token).Error
			if err != nil {
				t.Logf("❌ Database query still failing: %v", err)
			} else if token.Currency == "" {
				t.Logf("📊 Token not found in database (will return 404)")
			} else {
				t.Logf("✅ Token found: %s/%s (Price: $%.6f)", token.Currency, token.Issuer, token.Price)

				// Test price calculation for this token
				if token.Price <= 0 {
					t.Logf("⚠️  Token has no price data - testing price calculation...")

					// Test live price calculation - will implement after fixing core query issues
					t.Logf("🔍 Token %s needs price calculation (current: $%.6f)", testCurrency, token.Price)
				}
			}
		})
	}

	// Test a broader sample of tokens to see scope of issues
	t.Log("\n🔍 Testing broader token sample for systematic issues...")

	var sampleTokens []struct {
		Currency string  `json:"currency"`
		Issuer   string  `json:"issuer"`
		Name     string  `json:"token_name"`
		Price    float64 `json:"price"`
	}

	err := xrpDB.Raw(`
		SELECT currency, issuer, token_name, price
		FROM xrpTokens 
		WHERE currency != '' AND issuer != ''
		ORDER BY trustlines DESC
		LIMIT 10
	`).Scan(&sampleTokens).Error

	if err != nil {
		t.Errorf("❌ Failed to get token sample: %v", err)
	} else {
		t.Logf("📊 Testing %d top tokens for endpoint compatibility:", len(sampleTokens))

		for i, sampleToken := range sampleTokens {
			// Test the fixed query on each token
			var testResult struct {
				Currency string  `json:"currency"`
				Price    float64 `json:"price"`
			}

			testQuery := `
				SELECT currency, price
				FROM xrpTokens 
				WHERE currency = ? AND issuer = ?
			`

			err := xrpDB.Raw(testQuery, sampleToken.Currency, sampleToken.Issuer).Scan(&testResult).Error
			if err != nil {
				t.Errorf("❌ Token %d (%s) query failed: %v", i+1, sampleToken.Currency, err)
			} else {
				priceStatus := "✅"
				if testResult.Price <= 0 {
					priceStatus = "⚠️ NO PRICE"
				}
				t.Logf("   %d. %s/%s %s ($%.6f)", i+1, testResult.Currency, sampleToken.Issuer[:8]+"...", priceStatus, testResult.Price)
			}
		}
	}
}
