package xrp

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabaseSchema(t *testing.T) {
	// Connect directly to the database
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

	t.Log("🔍 Checking database schema...")

	// List all tables in the XRP database
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		t.Fatalf("Failed to list tables: %v", err)
	}
	defer rows.Close()

	t.Log("\n=== Available Tables ===")
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

	// Check for AMM-related tables
	t.Log("\n=== AMM-Related Tables ===")
	ammTables := []string{
		"xrpAmm",
		"xrp_amm",
		"amm_pools",
		"xrp_token_prices",
		"token_prices",
		"prices",
	}

	for _, tableName := range ammTables {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'XRP' AND table_name = ?", tableName).Scan(&exists)
		if err != nil {
			t.Logf("❌ Error checking table %s: %v", tableName, err)
		} else if exists > 0 {
			t.Logf("✅ Found table: %s", tableName)
			
			// Get table structure
			rows, err := db.Query("DESCRIBE " + tableName)
			if err != nil {
				t.Logf("   ❌ Failed to describe table: %v", err)
			} else {
				defer rows.Close()
				t.Logf("   📊 Columns:")
				for rows.Next() {
					var field, typ, null, key, defaultVal, extra sql.NullString
					if err := rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra); err != nil {
						continue
					}
					t.Logf("      - %s (%s)", field.String, typ.String)
				}
			}
		} else {
			t.Logf("❌ Table not found: %s", tableName)
		}
	}

	// Check for any table with "price" in the name
	t.Log("\n=== Tables with 'price' in name ===")
	rows2, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'XRP' AND table_name LIKE '%price%'")
	if err != nil {
		t.Logf("❌ Failed to search for price tables: %v", err)
	} else {
		defer rows2.Close()
		found := false
		for rows2.Next() {
			var tableName string
			if err := rows2.Scan(&tableName); err != nil {
				continue
			}
			t.Logf("💰 Found price table: %s", tableName)
			found = true
		}
		if !found {
			t.Log("❌ No tables with 'price' in name found")
		}
	}

	// Check for any table with "amm" in the name
	t.Log("\n=== Tables with 'amm' in name ===")
	rows3, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'XRP' AND table_name LIKE '%amm%'")
	if err != nil {
		t.Logf("❌ Failed to search for AMM tables: %v", err)
	} else {
		defer rows3.Close()
		found := false
		for rows3.Next() {
			var tableName string
			if err := rows3.Scan(&tableName); err != nil {
				continue
			}
			t.Logf("🏊 Found AMM table: %s", tableName)
			found = true
		}
		if !found {
			t.Log("❌ No tables with 'amm' in name found")
		}
	}
} 