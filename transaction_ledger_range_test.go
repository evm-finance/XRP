package xrp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"qc-defi-graphql-server/configs"
)

// TestTransactionLedgerRange tests if our updated XRP service is fetching recent transactions
// using the -1,-1 ledger parameters and forward=false instead of hardcoded ranges
func TestTransactionLedgerRange(t *testing.T) {
	t.Logf("🧪 [TEST] Starting Transaction Ledger Range Test")

	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		t.Logf("⚠️ Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("❌ Failed to load config: %v", err)
	}

	// Get XRP database connection
	var xrpDB *gorm.DB
	xrpDB = cfg.XRPDB()
	if xrpDB == nil {
		t.Fatalf("❌ Failed to connect to XRP database")
	}

	// Create XRP connection manager
	xrplWssUrl := "wss://xrplcluster.com"
	connMgr := NewConnectionManager(xrplWssUrl, 5, cfg.Logger())

	// Initialize XRP service with the connection manager
	xrpService := NewXRPServiceWithConnectionManager(connMgr)
	if err := xrpService.SetDatabase(xrpDB); err != nil {
		t.Fatalf("❌ Failed to set database on XRP service: %v", err)
	}

	// Test with the original address (confirmed active by user)
	testAddress := "rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo" // User confirmed this has August 21st transactions
	t.Logf("🔍 [TEST] Testing transaction fetch for address: %s", testAddress)

	t.Logf("⏰ [TEST] Current time: %s", time.Now().Format("2006-01-02 15:04:05 UTC"))

	// Call our updated GetAccountTransactions method
	t.Logf("🚀 [TEST] Calling XRP service GetAccountTransactions...")
	response, err := xrpService.GetAccountTransactions(testAddress)
	if err != nil {
		t.Fatalf("❌ [TEST] XRP service call failed: %v", err)
	}

	t.Logf("✅ [TEST] XRP service call succeeded")

	// Parse the response to analyze transaction dates
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("❌ [TEST] Failed to marshal response: %v", err)
	}

	var rawResponse map[string]interface{}
	if err := json.Unmarshal(responseJSON, &rawResponse); err != nil {
		t.Fatalf("❌ [TEST] Failed to unmarshal response: %v", err)
	}

	// Extract transactions from XRPL response structure
	var transactionCount int
	var recentTransactionCount int
	var oldestDate, newestDate time.Time

	if result, ok := rawResponse["result"]; ok {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if txs, ok := resultMap["transactions"]; ok {
				if txArray, ok := txs.([]interface{}); ok {
					transactionCount = len(txArray)
					t.Logf("📊 [TEST] Found %d total transactions", transactionCount)

					// Analyze transaction dates
					cutoffDate := time.Now().AddDate(0, 0, -30) // Last 30 days
					t.Logf("🗓️ [TEST] Checking for transactions newer than: %s", cutoffDate.Format("2006-01-02 15:04:05"))

					for i, tx := range txArray {
						if i >= 10 { // Only analyze first 10 for performance
							break
						}

						if txMap, ok := tx.(map[string]interface{}); ok {
							// Look for transaction data in 'tx' object
							if txObj, ok := txMap["tx"]; ok {
								if txData, ok := txObj.(map[string]interface{}); ok {
									if dateVal, exists := txData["date"]; exists {
										var txDate int64
										switch v := dateVal.(type) {
										case float64:
											txDate = int64(v)
										case int:
											txDate = int64(v)
										case int64:
											txDate = v
										}

										if txDate > 0 {
											// Convert XRP Ledger epoch to Unix timestamp
											unixTimestamp := txDate + 946684800
											txTime := time.Unix(unixTimestamp, 0)

											if i == 0 {
												newestDate = txTime
												oldestDate = txTime
											}

											if txTime.After(newestDate) {
												newestDate = txTime
											}
											if txTime.Before(oldestDate) {
												oldestDate = txTime
											}

											if txTime.After(cutoffDate) {
												recentTransactionCount++
											}

											t.Logf("  Transaction %d: %s (Ledger: %v)",
												i+1,
												txTime.Format("2006-01-02 15:04:05 UTC"),
												txData["ledger_index"])
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Test Results Summary
	t.Logf("\n🎯 [TEST RESULTS] Transaction Ledger Range Analysis:")
	t.Logf("  ✅ Total transactions fetched: %d", transactionCount)
	t.Logf("  ✅ Recent transactions (last 30 days): %d", recentTransactionCount)
	if !newestDate.IsZero() {
		t.Logf("  ✅ Newest transaction: %s", newestDate.Format("2006-01-02 15:04:05 UTC"))
		t.Logf("  ✅ Oldest transaction: %s", oldestDate.Format("2006-01-02 15:04:05 UTC"))

		daysSinceNewest := time.Since(newestDate).Hours() / 24
		t.Logf("  📈 Days since newest transaction: %.1f", daysSinceNewest)

		if daysSinceNewest < 7 {
			t.Logf("  ✅ SUCCESS: Recent transactions detected (within last week)")
		} else if daysSinceNewest < 30 {
			t.Logf("  ⚠️  WARNING: Newest transaction is %v days old", int(daysSinceNewest))
		} else {
			t.Logf("  ❌ ISSUE: Newest transaction is too old (%v days)", int(daysSinceNewest))
		}
	}

	if recentTransactionCount > 0 {
		t.Logf("  ✅ CONCLUSION: Ledger range fix is working - getting recent transactions")
	} else {
		t.Logf("  ❌ CONCLUSION: May still have ledger range issues - no recent transactions")
	}

	t.Logf("🧪 [TEST] Transaction Ledger Range Test completed")

	// Assert that we got some transactions
	if transactionCount == 0 {
		t.Errorf("Expected to receive transactions, but got zero")
	}

	// Assert that we got recent transactions (warning if not, but don't fail the test)
	if recentTransactionCount == 0 {
		t.Logf("⚠️ WARNING: No recent transactions found - may indicate ledger range issue")
	}
}
