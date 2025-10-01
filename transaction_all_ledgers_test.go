package xrp

import (
	"fmt"
	"qc-defi-graphql-server/configs"
	"strings"
	"testing"
	"time"

	xrpl "github.com/xrpscan/xrpl-go"
)

// TestTransactionFetchAllLedgers tests fetching ALL transactions using -1,-1 parameters
func TestTransactionFetchAllLedgers(t *testing.T) {
	fmt.Println("🧪 TESTING TRANSACTION FETCH FROM ALL LEDGERS (-1, -1)")
	fmt.Println(strings.Repeat("=", 58))

	// Connect using config pattern
	conf, err := configs.NewConfigRepo("../../.env")
	if err != nil {
		t.Fatalf("❌ Failed to create config: %v", err)
	}

	db := conf.XRPDB()
	if db == nil {
		t.Fatal("❌ Failed to get XRPDB connection")
	}

	fmt.Println("✅ Database connection established")

	// Test with a well-known XRP address
	testAddress := "rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo"

	fmt.Printf("\n🔍 Testing ALL LEDGERS (-1,-1) vs Current (10M) for address: %s\n", testAddress)

	// Create XRP service
	xrpService := NewXRPService("wss://xrplcluster.com", 1, nil)
	if err := xrpService.SetDatabase(db); err != nil {
		t.Fatalf("❌ Failed to set database: %v", err)
	}

	// Get connection for direct XRPL testing
	client, err := xrpService.GetClient()
	if err != nil {
		t.Fatalf("❌ Failed to get XRPL client: %v", err)
	}
	defer xrpService.GetConnectionManager().ReturnConnection(client)

	// Test 1: Current approach (10M ledgers)
	fmt.Println("\n📊 Test 1: Current 10M Ledger Approach")
	currentResult := testSpecificLedgerRange(client, testAddress, -10000000, -1, "Current 10M", t)

	// Test 2: ALL ledgers approach (-1, -1)
	fmt.Println("\n📊 Test 2: ALL Ledgers Approach (-1, -1)")
	allResult := testSpecificLedgerRange(client, testAddress, -1, -1, "ALL Ledgers", t)

	// Test 3: Other ranges for comparison
	fmt.Println("\n📊 Test 3: Other Range Comparisons")
	yearResult := testSpecificLedgerRange(client, testAddress, -28800000, -1, "~1 Year", t)
	monthResult := testSpecificLedgerRange(client, testAddress, -864000, -1, "~1 Month", t)

	// Results Summary
	fmt.Println("\n🎯 LEDGER RANGE COMPARISON RESULTS:")
	fmt.Println("====================================")
	printTestResult("Current (10M ledgers)", currentResult)
	printTestResult("ALL Ledgers (-1,-1)", allResult)
	printTestResult("1 Year (28.8M ledgers)", yearResult)
	printTestResult("1 Month (864k ledgers)", monthResult)

	// Analysis and Recommendations
	fmt.Println("\n💡 ANALYSIS & RECOMMENDATIONS:")
	fmt.Println("===============================")

	if allResult != nil && currentResult != nil {
		if allResult.TxCount > currentResult.TxCount {
			diff := allResult.TxCount - currentResult.TxCount
			fmt.Printf("✅ ALL LEDGERS finds %d MORE transactions (%d vs %d)\n",
				diff, allResult.TxCount, currentResult.TxCount)
			fmt.Println("🚀 RECOMMEND: Switch to (-1,-1) for complete transaction history")
		} else if allResult.TxCount == currentResult.TxCount {
			fmt.Printf("⚖️  ALL LEDGERS returns SAME count as current (%d transactions)\n", allResult.TxCount)
			fmt.Println("💡 SUGGEST: Use (-1,-1) for future-proof complete coverage")
		} else {
			fmt.Printf("⚠️  ALL LEDGERS returned FEWER transactions (%d vs %d) - investigate\n",
				allResult.TxCount, currentResult.TxCount)
		}

		// Performance comparison
		if allResult.Duration > currentResult.Duration*2 {
			fmt.Printf("⏱️  ALL LEDGERS takes longer (%v vs %v)\n", allResult.Duration, currentResult.Duration)
			fmt.Println("💡 Consider pagination for large result sets in production")
		} else {
			fmt.Printf("⚡ Performance acceptable: ALL LEDGERS (%v) vs Current (%v)\n",
				allResult.Duration, currentResult.Duration)
		}
	}

	t.Logf("All ledgers test completed - compared multiple approaches")
}

// testSpecificLedgerRange tests a specific ledger range configuration
func testSpecificLedgerRange(client *xrpl.Client, address string, minLedger, maxLedger int, name string, t *testing.T) *LedgerTestResult {
	fmt.Printf("   🔄 Testing %s (min:%d, max:%d)...\n", name, minLedger, maxLedger)

	startTime := time.Now()

	// Create XRPL request
	transactionRequest := xrpl.BaseRequest{
		"id":               2,
		"command":          "account_tx",
		"account":          address,
		"ledger_index_min": minLedger,
		"ledger_index_max": maxLedger,
		"binary":           false,
		"forward":          true,
		"limit":            500,
	}

	response, err := client.Request(transactionRequest)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("   ❌ %s failed after %v: %v\n", name, duration, err)
		return nil
	}

	// Analyze response
	result := &LedgerTestResult{
		Name:     name,
		Duration: duration,
		MinReq:   minLedger,
		MaxReq:   maxLedger,
	}

	// Parse XRPL response
	responseMap := map[string]interface{}(response)
	if resultData, ok := responseMap["result"].(map[string]interface{}); ok {
		// Get actual ledger range returned by XRPL
		if ledgerMin, ok := resultData["ledger_index_min"].(float64); ok {
			result.MinActual = int(ledgerMin)
		}
		if ledgerMax, ok := resultData["ledger_index_max"].(float64); ok {
			result.MaxActual = int(ledgerMax)
		}
		if limit, ok := resultData["limit"].(float64); ok {
			result.Limit = int(limit)
		}
		if validated, ok := resultData["validated"].(bool); ok {
			result.Validated = validated
		}

		// Count transactions
		if txs, ok := resultData["transactions"].([]interface{}); ok {
			result.TxCount = len(txs)

			// Get date range if available
			if len(txs) > 0 {
				result.DateInfo = getDateInfo(txs)
			}
		}

		// Check for pagination
		if marker, ok := resultData["marker"]; ok {
			result.HasMore = true
			result.Marker = fmt.Sprintf("%v", marker)
		}
	}

	fmt.Printf("   ✅ %s: %d txs, %v, ledgers %d-%d\n",
		name, result.TxCount, duration, result.MinActual, result.MaxActual)

	return result
}

// LedgerTestResult holds results from a ledger range test
type LedgerTestResult struct {
	Name      string
	Duration  time.Duration
	MinReq    int    // Requested min ledger
	MaxReq    int    // Requested max ledger
	MinActual int    // Actual min ledger from XRPL
	MaxActual int    // Actual max ledger from XRPL
	Limit     int    // Limit used by XRPL
	TxCount   int    // Number of transactions returned
	Validated bool   // Whether results are validated
	HasMore   bool   // Whether more data is available
	Marker    string // Pagination marker if present
	DateInfo  string // Date range info
}

// getDateInfo extracts date information from transactions
func getDateInfo(txs []interface{}) string {
	if len(txs) == 0 {
		return "no dates"
	}

	// Check first and last transactions for dates
	var newest, oldest time.Time
	validDates := 0

	for _, idx := range []int{0, len(txs) - 1} {
		if txMap, ok := txs[idx].(map[string]interface{}); ok {
			var txDate time.Time

			// Try to find date field
			if date, ok := txMap["date"].(float64); ok {
				rippleEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
				txDate = rippleEpoch.Add(time.Duration(date) * time.Second)
			} else if tx, ok := txMap["tx"].(map[string]interface{}); ok {
				if date, ok := tx["date"].(float64); ok {
					rippleEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
					txDate = rippleEpoch.Add(time.Duration(date) * time.Second)
				}
			}

			if !txDate.IsZero() {
				validDates++
				if validDates == 1 || txDate.After(newest) {
					newest = txDate
				}
				if validDates == 1 || txDate.Before(oldest) {
					oldest = txDate
				}
			}
		}
	}

	if validDates == 0 {
		return "no valid dates"
	}

	if validDates == 1 {
		return newest.Format("2006-01-02")
	}

	span := newest.Sub(oldest)
	daysSpan := int(span.Hours() / 24)

	recentInfo := ""
	if time.Since(newest).Hours()/24 < 7 {
		recentInfo = " (recent)"
	}

	return fmt.Sprintf("%s to %s (%dd)%s", oldest.Format("2006-01-02"), newest.Format("2006-01-02"), daysSpan, recentInfo)
}

// printTestResult prints formatted test result
func printTestResult(name string, result *LedgerTestResult) {
	if result == nil {
		fmt.Printf("❌ %-20s: FAILED\n", name)
		return
	}

	status := "✅"
	if !result.Validated {
		status = "⚠️ "
	}

	moreInfo := ""
	if result.HasMore {
		moreInfo = " (+more)"
	}

	fmt.Printf("%s %-20s: %3d txs | %8v | %s%s\n",
		status, name, result.TxCount, result.Duration, result.DateInfo, moreInfo)
}
