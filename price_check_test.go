package xrp

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestPriceCheckMinimal(t *testing.T) {
	// Connect directly to the database without using the config system
	dsn := "coindatabase:Smokey1!@tcp(coindatabase2.cbejwbliijcq.us-west-2.rds.amazonaws.com:3306)/XRP?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Log("🔍 Checking database prices (minimal test)...")

	// Test 1: Check XRP price from XRP/RLUSD pool
	t.Log("\n=== Test 1: XRP Price from XRP/RLUSD Pool ===")

	var xrpPrice float64
	row := db.QueryRow(`
		SELECT 
			CASE 
				WHEN amount2currency = ? AND amount2issuer = ? THEN 
					CAST(amount2value AS FLOAT) / CAST(amount AS FLOAT)
				ELSE 
					CAST(amount AS FLOAT) / CAST(amount2value AS FLOAT)
			END as xrp_price
		FROM xrpAmm
		WHERE (amount2currency = ? AND amount2issuer = ?) 
		   OR (amount2currency = 'XRP' AND amount2issuer = '')
		LIMIT 1
	`, RLUSDCurrency, RLUSDIssuer, RLUSDCurrency, RLUSDIssuer)

	if err := row.Scan(&xrpPrice); err != nil {
		t.Logf("❌ No XRP/RLUSD pool found: %v", err)
	} else {
		t.Logf("✅ XRP Price from AMM pool: $%.6f", xrpPrice)

		// Check if price is reasonable (between $0.10 and $10)
		if xrpPrice < 0.10 || xrpPrice > 10.0 {
			t.Logf("⚠️  Warning: XRP price seems unusual: $%.6f", xrpPrice)
		}
	}

	// Test 2: Check stored token prices in xrp_token_prices table
	t.Log("\n=== Test 2: Stored Token Prices ===")

	rows, err := db.Query(`
		SELECT currency, issuer, price_usd, price_xrp, source, last_updated
		FROM xrp_token_prices
		ORDER BY last_updated DESC
		LIMIT 10
	`)

	if err != nil {
		t.Logf("❌ Failed to query token prices: %v", err)
	} else {
		defer rows.Close()

		count := 0
		for rows.Next() {
			var currency, issuer, source string
			var priceUSD, priceXRP float64
			var lastUpdated time.Time

			if err := rows.Scan(&currency, &issuer, &priceUSD, &priceXRP, &source, &lastUpdated); err != nil {
				continue
			}

			age := time.Since(lastUpdated)
			status := "✅"
			if age > 5*time.Minute {
				status = "⚠️ "
			}
			if age > 30*time.Minute {
				status = "❌"
			}

			issuerShort := issuer
			if len(issuerShort) > 8 {
				issuerShort = issuerShort[:8] + "..."
			}

			t.Logf("%s Token: %s/%s | USD: $%.6f | XRP: %.6f | Source: %s | Age: %v",
				status, currency, issuerShort, priceUSD, priceXRP, source, age)
			count++
		}
		t.Logf("Found %d token prices in database", count)
	}

	// Test 3: Check AMM pool liquidity values
	t.Log("\n=== Test 3: AMM Pool Liquidity ===")

	rows2, err := db.Query(`
		SELECT account, liquidity_usd, last_updated
		FROM xrpAmm
		WHERE liquidity_usd > 0
		ORDER BY liquidity_usd DESC
		LIMIT 10
	`)

	if err != nil {
		t.Logf("❌ Failed to query pool liquidity: %v", err)
	} else {
		defer rows2.Close()

		t.Log("Top 10 pools by liquidity:")
		rank := 1
		for rows2.Next() {
			var account string
			var liquidityUSD float64
			var lastUpdated time.Time

			if err := rows2.Scan(&account, &liquidityUSD, &lastUpdated); err != nil {
				continue
			}

			age := time.Since(lastUpdated)
			status := "✅"
			if age > 5*time.Minute {
				status = "⚠️ "
			}
			if age > 30*time.Minute {
				status = "❌"
			}

			accountShort := account
			if len(accountShort) > 10 {
				accountShort = accountShort[:10] + "..."
			}

			t.Logf("%s #%d Pool: %s | Liquidity: $%.2f | Age: %v",
				status, rank, accountShort, liquidityUSD, age)
			rank++
		}
	}

	// Test 4: Check for stale prices
	t.Log("\n=== Test 4: Stale Price Check ===")

	var staleCount int
	row = db.QueryRow(`
		SELECT COUNT(*) 
		FROM xrp_token_prices 
		WHERE last_updated < ?
	`, time.Now().Add(-30*time.Minute))

	if err := row.Scan(&staleCount); err != nil {
		t.Logf("❌ Failed to count stale prices: %v", err)
	} else {
		if staleCount > 0 {
			t.Logf("⚠️  Found %d tokens with prices older than 30 minutes", staleCount)
		} else {
			t.Logf("✅ All token prices are fresh (updated within 30 minutes)")
		}
	}

	// Test 5: Price consistency check
	t.Log("\n=== Test 5: Price Consistency Check ===")

	// Check if XRP price from pools matches stored XRP price
	var storedXRPPrice float64
	row = db.QueryRow(`
		SELECT price_usd 
		FROM xrp_token_prices 
		WHERE currency = 'XRP' AND issuer = ''
		ORDER BY last_updated DESC
		LIMIT 1
	`)

	if err := row.Scan(&storedXRPPrice); err == nil {
		if xrpPrice > 0 {
			priceDiff := abs(xrpPrice - storedXRPPrice)
			percentDiff := (priceDiff / xrpPrice) * 100

			if percentDiff > 5 {
				t.Logf("⚠️  Price mismatch: AMM XRP price ($%.6f) differs from stored price ($%.6f) by %.2f%%",
					xrpPrice, storedXRPPrice, percentDiff)
			} else {
				t.Logf("✅ Price consistency: AMM and stored XRP prices match within 5%%")
			}
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
