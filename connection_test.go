package xrp

import (
	"testing"

	"qc-defi-graphql-server/configs"
)

func TestDatabaseConnection(t *testing.T) {
	// Use the proper config system
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	db := conf.XRPDB()
	if db == nil {
		t.Fatalf("Failed to get XRP database connection")
	}

	t.Log("✅ Successfully connected to XRP database via config system")

	// Check what tables actually exist
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}
	defer rows.Close()

	t.Log("\n=== All Tables in XRP Database ===")
	count := 0
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		t.Logf("📋 Table: %s", tableName)
		count++
	}
	t.Logf("Total tables found: %d", count)

	// Specifically check for xrpAmm_normalized table
	t.Log("\n=== Checking xrpAmm_normalized Table ===")
	var ammCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized").Scan(&ammCount).Error
	if err != nil {
		t.Logf("❌ Error accessing xrpAmm_normalized table: %v", err)
	} else {
		t.Logf("✅ xrpAmm_normalized table exists with %d rows", ammCount)

		// Get a sample row
		type AmmRow struct {
			Account        string
			Asset1Amount   string
			Asset1Currency string
			Asset1Issuer   string
			Asset2Amount   string
		}

		var sampleRow AmmRow
		err = db.Raw("SELECT account, asset1_amount, asset1_currency, asset1_issuer, asset2_amount FROM xrpAmm_normalized LIMIT 1").Scan(&sampleRow).Error
		if err != nil {
			t.Logf("❌ Error reading from xrpAmm_normalized: %v", err)
		} else {
			accountShort := sampleRow.Account
			if len(accountShort) > 10 {
				accountShort = accountShort[:10] + "..."
			}
			issuerShort := sampleRow.Asset1Issuer
			if len(issuerShort) > 10 {
				issuerShort = issuerShort[:10] + "..."
			}
			t.Logf("📊 Sample xrpAmm_normalized row: Account=%s, Asset1=%s, Currency=%s, Issuer=%s, Asset2=%s",
				accountShort, sampleRow.Asset1Amount, sampleRow.Asset1Currency, issuerShort, sampleRow.Asset2Amount)
		}
	}

	// Specifically check for xrpTokens table
	t.Log("\n=== Checking xrpTokens Table ===")
	var tokenCount int64
	err = db.Raw("SELECT COUNT(*) FROM xrpTokens").Scan(&tokenCount).Error
	if err != nil {
		t.Logf("❌ Error accessing xrpTokens table: %v", err)
	} else {
		t.Logf("✅ xrpTokens table exists with %d rows", tokenCount)

		// Get a sample row
		type TokenRow struct {
			Currency string
			Issuer   string
			Name     string
			Price    float64
		}

		var sampleRow TokenRow
		err = db.Raw("SELECT currency, issuer, name, price FROM xrpTokens LIMIT 1").Scan(&sampleRow).Error
		if err != nil {
			t.Logf("❌ Error reading from xrpTokens: %v", err)
		} else {
			issuerShort := sampleRow.Issuer
			if len(issuerShort) > 10 {
				issuerShort = issuerShort[:10] + "..."
			}
			t.Logf("📊 Sample xrpTokens row: Currency=%s, Issuer=%s, Name=%s, Price=%.6f",
				sampleRow.Currency, issuerShort, sampleRow.Name, sampleRow.Price)
		}
	}
}
