package xrp_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"qc-defi-graphql-server/configs"
	. "qc-defi-graphql-server/internal/xrp"
)

func TestDebugXRPPriceCalculation(t *testing.T) {
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	db := conf.XRPDB()
	if db == nil {
		t.Fatalf("Failed to get XRP database connection")
	}

	t.Log("🔍 Debugging XRP price calculation...")

	// First, let's see what XRP/RLUSD pools exist
	t.Log("\n=== Checking for XRP/RLUSD pools ===")

	// Check for pools with RLUSD
	rows, err := db.Raw(`
		SELECT account, amount, amount2currency, amount2issuer, amount2value
		FROM xrpAmm
		WHERE amount2currency = ? AND amount2issuer = ?
		LIMIT 5
	`, RLUSDCurrency, RLUSDIssuer).Rows()
	if err != nil {
		t.Fatalf("Failed to query RLUSD pools: %v", err)
	}
	defer rows.Close()

	t.Log("Pools with RLUSD as amount2:")
	for rows.Next() {
		var account, amount, amount2currency, amount2issuer, amount2value string
		if err := rows.Scan(&account, &amount, &amount2currency, &amount2issuer, &amount2value); err != nil {
			continue
		}
		t.Logf("Account: %s | Amount: %s | Amount2: %s/%s = %s",
			account, amount, amount2currency, amount2issuer, amount2value)
	}

	// Check for pools with XRP
	rows2, err := db.Raw(`
		SELECT account, amount, amount2currency, amount2issuer, amount2value
		FROM xrpAmm
		WHERE amount2currency = 'XRP' AND amount2issuer = ''
		LIMIT 5
	`).Rows()
	if err != nil {
		t.Fatalf("Failed to query XRP pools: %v", err)
	}
	defer rows2.Close()

	t.Log("\nPools with XRP as amount2:")
	for rows2.Next() {
		var account, amount, amount2currency, amount2issuer, amount2value string
		if err := rows2.Scan(&account, &amount, &amount2currency, &amount2issuer, &amount2value); err != nil {
			continue
		}
		t.Logf("Account: %s | Amount: %s | Amount2: XRP = %s",
			account, amount, amount2value)
	}

	// Now let's test the actual price service
	t.Log("\n=== Testing Price Service ===")
	priceService := NewPriceServiceWithDB(db)

	xrpPrice, err := priceService.GetXRPLPrice()
	if err != nil {
		t.Logf("❌ Error getting XRP price: %v", err)
	} else {
		t.Logf("✅ XRP price: $%.6f USD", xrpPrice.InexactFloat64())
	}

	// Let's also check what the current query in GetXRPLPrice returns
	t.Log("\n=== Testing the exact query from GetXRPLPrice ===")
	row := db.Raw(`
		SELECT amount, amount2value, amount2currency, amount2issuer
		FROM xrpAmm
		WHERE (amount2currency = ? AND amount2issuer = ?) OR (amount2currency = ? AND amount2issuer = '')
		LIMIT 1
	`, RLUSDCurrency, RLUSDIssuer, "XRP").Row()

	var amount, amount2value, amount2currency, amount2issuer string
	if err := row.Scan(&amount, &amount2value, &amount2currency, &amount2issuer); err != nil {
		t.Logf("❌ Error scanning row: %v", err)
	} else {
		t.Logf("✅ Found pool: Amount=%s, Amount2Value=%s, Amount2Currency=%s, Amount2Issuer=%s",
			amount, amount2value, amount2currency, amount2issuer)

		// Calculate price manually
		var xrpBalance, rlusdBalance decimal.Decimal
		if amount2currency == RLUSDCurrency && amount2issuer == RLUSDIssuer {
			// amount is XRP (in drops), amount2value is RLUSD
			xrpBalance, _ = decimal.NewFromString(amount)
			xrpBalance = xrpBalance.Div(decimal.NewFromFloat(1000000.0)) // drops to XRP
			rlusdBalance, _ = decimal.NewFromString(amount2value)
			t.Logf("Case 1: XRP=%s, RLUSD=%s", xrpBalance.String(), rlusdBalance.String())
		} else if amount2currency == "XRP" {
			// amount2value is XRP (in drops), amount is RLUSD
			xrpBalance, _ = decimal.NewFromString(amount2value)
			xrpBalance = xrpBalance.Div(decimal.NewFromFloat(1000000.0))
			rlusdBalance, _ = decimal.NewFromString(amount)
			t.Logf("Case 2: XRP=%s, RLUSD=%s", xrpBalance.String(), rlusdBalance.String())
		} else {
			t.Logf("❌ Unexpected pool structure")
		}

		if xrpBalance.GreaterThan(decimal.Zero) && rlusdBalance.GreaterThan(decimal.Zero) {
			price := rlusdBalance.Div(xrpBalance)
			t.Logf("💰 Calculated price: RLUSD/XRP = %s", price.String())
			t.Logf("💰 This means 1 XRP = %s RLUSD ≈ $%s USD", price.String(), price.String())
		}
	}
}
