package xrp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xrpscan/xrpl-go"
	"gorm.io/gorm"
)

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// EnhancedLedgerService provides real XRPL data functionality
type EnhancedLedgerService struct {
	db          *gorm.DB
	xrpConnMgr  *ConnectionManager
	redisClient *redis.Client

	// Runtime state tracking for cache freshness validation
	latestProcessedLedgerIndex int
	cacheFreshnessThreshold    int
	lastCacheCheckTime         time.Time
	cacheCheckInterval         time.Duration

	// Mutex for thread-safe state updates
	stateMutex sync.RWMutex
}

// NewEnhancedLedgerService creates a new enhanced ledger service
func NewEnhancedLedgerService(db *gorm.DB, connMgr *ConnectionManager, redisClient *redis.Client) *EnhancedLedgerService {
	// Get cache freshness threshold from environment or use default
	cacheThreshold := 10 // default threshold
	if thresholdStr := getEnvOrDefault("CACHE_FRESHNESS_THRESHOLD", "10"); thresholdStr != "" {
		if threshold, err := strconv.Atoi(thresholdStr); err == nil && threshold > 0 {
			cacheThreshold = threshold
		}
	}

	return &EnhancedLedgerService{
		db:                      db,
		xrpConnMgr:              connMgr,
		redisClient:             redisClient,
		cacheFreshnessThreshold: cacheThreshold,
		cacheCheckInterval:      30 * time.Second, // Check cache freshness every 30 seconds
		lastCacheCheckTime:      time.Now(),
	}
}

// GetClient returns a pooled XRPL client
func (s *EnhancedLedgerService) GetClient() (*xrpl.Client, error) {
	if s.xrpConnMgr == nil {
		return nil, fmt.Errorf("XRPL connection manager not initialized")
	}
	return s.xrpConnMgr.GetConnection()
}

// GetRealLedgerData fetches actual ledger data from XRPL
func (s *EnhancedLedgerService) GetRealLedgerData(ledgerIndex int) (*XRPLedger, error) {
	client, err := s.GetClient()
	if err != nil {
		return nil, err
	}
	defer s.xrpConnMgr.ReturnConnection(client)

	return s.getRealLedgerDataWithClient(client, ledgerIndex)
}

// getRealLedgerDataWithClient fetches ledger data using a provided client (for batch operations)
func (s *EnhancedLedgerService) getRealLedgerDataWithClient(client *xrpl.Client, ledgerIndex int) (*XRPLedger, error) {
	// Request ledger data with full transaction details
	request := xrpl.BaseRequest{
		"id":           1,
		"command":      "ledger",
		"ledger_index": ledgerIndex,
		"full":         false,
		"accounts":     false,
		"transactions": true,
		"expand":       true,
		"owner_funds":  false,
	}

	response, err := client.Request(request)
	if err != nil {
		s.xrpConnMgr.MarkConnectionUnhealthy(client)
		return nil, fmt.Errorf("failed to request ledger data: %w", err)
	}

	// Parse the response
	jsonStr, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var result XRPLedgerResult
	if err := json.Unmarshal(jsonStr, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ledger result: %w", err)
	}

	if result.Status == "error" {
		return nil, fmt.Errorf("XRPL API error for ledger %d", ledgerIndex)
	}

	// Count events in the ledger
	if result.Result != nil {
		result.Result.CountEvents()
	}

	return result.Result, nil
}

// GetCurrentLedgerIndex gets the current ledger index from XRPL network
func (s *EnhancedLedgerService) GetCurrentLedgerIndex() (int, error) {
	// Get client with timeout handling
	client, err := s.GetClient()
	if err != nil {
		return 0, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer s.xrpConnMgr.ReturnConnection(client)

	// Use the correct XRPL library API
	request := xrpl.BaseRequest{
		"id":      1,
		"command": "ledger_current",
	}

	// Add timeout handling - create a channel for the request
	type result struct {
		response map[string]interface{}
		err      error
	}

	resultChan := make(chan result, 1)

	// Run the XRPL request in a goroutine with timeout
	go func() {
		response, err := client.Request(request)
		resultChan <- result{response: response, err: err}
	}()

	// Wait for result or timeout
	select {
	case res := <-resultChan:
		if res.err != nil {
			log.Printf("⚠️ XRPL network call failed: %v", res.err)
			return 0, fmt.Errorf("failed to get current ledger from XRPL: %w", res.err)
		}

		// Parse the response using the original logic
		if result, ok := res.response["result"].(map[string]interface{}); ok {
			if ledgerIndex, ok := result["ledger_current_index"].(float64); ok {
				return int(ledgerIndex), nil
			}
		}

		return 0, fmt.Errorf("failed to parse current ledger index from XRPL response")

	case <-time.After(5 * time.Second):
		log.Printf("⚠️ XRPL network call timed out after 5 seconds")
		return 0, fmt.Errorf("XRPL network call timed out")
	}
}

// GetRecentLedgers fetches recent ledger data using Redis caching with freshness validation
func (s *EnhancedLedgerService) GetRecentLedgers(limit int) ([]map[string]interface{}, error) {
	// First check if cache is fresh before using it
	isFresh, gap, err := s.IsCacheFresh()
	if err != nil {
		log.Printf("⚠️ Cache freshness check failed: %v, falling back to live XRPL data", err)
		isFresh = false
	}

	// If cache is fresh, try to get cached ledgers
	if isFresh {
		if cachedLedgers, err := s.getCachedRecentLedgers(limit); err == nil && len(cachedLedgers) > 0 {
			// Validate cached data before returning
			validLedgers := make([]map[string]interface{}, 0, len(cachedLedgers))
			for _, ledger := range cachedLedgers {
				if s.ValidateLedgerData(ledger) {
					validLedgers = append(validLedgers, ledger)
				} else {
					log.Printf("⚠️ Skipping invalid cached ledger: %v", ledger)
				}
			}

			if len(validLedgers) > 0 {
				log.Printf("✅ Retrieved %d valid recent ledgers from fresh Redis cache (gap=%d)", len(validLedgers), gap)
				return validLedgers, nil
			} else {
				log.Printf("⚠️ No valid ledgers in cache, falling back to live XRPL data")
			}
		}
	} else {
		log.Printf("🔄 Cache is stale (gap=%d), fetching live XRPL data", gap)
	}

	// Cache miss, stale, or invalid - fetch from live XRPL
	log.Printf("🔄 Fetching recent ledgers from live XRPL network...")

	// Get current ledger index from XRPL with timeout and error handling
	currentIndex, err := s.GetCurrentLedgerIndex()
	if err != nil {
		log.Printf("⚠️ XRPL connection issue: %v", err)

		// Try to get ANY cached data as fallback, even if stale
		if fallbackLedgers, fallbackErr := s.getCachedRecentLedgers(limit * 2); fallbackErr == nil && len(fallbackLedgers) > 0 {
			log.Printf("✅ Using fallback cached ledgers (%d available) despite XRPL issues", len(fallbackLedgers))
			if len(fallbackLedgers) > limit {
				fallbackLedgers = fallbackLedgers[:limit]
			}
			return fallbackLedgers, nil
		}

		// If no cache available, return a more specific error
		return nil, fmt.Errorf("XRPL temporarily unavailable and no cached data available: %w", err)
	}

	var ledgers []map[string]interface{}

	// Fetch the last 'limit' ledgers using a SINGLE connection to prevent pool exhaustion
	client, err := s.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get XRPL client: %w", err)
	}
	defer s.xrpConnMgr.ReturnConnection(client)

	for i := 0; i < limit; i++ {
		ledgerIndex := currentIndex - i

		// Use the SAME client for all requests to avoid connection pool exhaustion
		ledger, err := s.getRealLedgerDataWithClient(client, ledgerIndex)
		if err != nil {
			log.Printf("⚠️ Skipping ledger %d - could not get real XRPL data: %v", ledgerIndex, err)
			continue // Skip this ledger - no fake data
		}

		if ledger == nil {
			log.Printf("⚠️ Skipping ledger %d - empty ledger data from XRPL", ledgerIndex)
			continue // Skip this ledger - no fake data
		}

		// Convert the full ledger object to map[string]interface{} for JSON response
		ledgerBytes, err := json.Marshal(ledger)
		if err != nil {
			log.Printf("⚠️ Skipping ledger %d - failed to marshal: %v", ledgerIndex, err)
			continue
		}

		var ledgerMap map[string]interface{}
		if err := json.Unmarshal(ledgerBytes, &ledgerMap); err != nil {
			log.Printf("⚠️ Skipping ledger %d - failed to unmarshal to map: %v", ledgerIndex, err)
			continue
		}

		// Validate the ledger data before adding
		if !s.ValidateLedgerData(ledgerMap) {
			log.Printf("⚠️ Skipping invalid ledger %d from XRPL", ledgerIndex)
			continue
		}

		ledgers = append(ledgers, ledgerMap)

		// Cache each ledger individually for future requests
		s.cacheLedgerData(ledgerIndex, ledgerMap)

		// Update the latest processed ledger index
		s.UpdateLatestProcessedLedger(ledgerIndex)

		// log.Printf("✅ Retrieved and cached real ledger %d with %d transactions", ledgerIndex, ledger.TxCount)
	}

	// If we got at least some ledgers, return them
	if len(ledgers) > 0 {
		// Cache the full result set for faster subsequent requests
		s.cacheRecentLedgers(ledgers)
		log.Printf("✅ Retrieved %d real recent XRPL ledgers and cached for future requests", len(ledgers))
		return ledgers, nil
	}

	// If no new ledgers but we have cache, try cache again with more lenient parameters
	if fallbackLedgers, fallbackErr := s.getCachedRecentLedgers(limit * 3); fallbackErr == nil && len(fallbackLedgers) > 0 {
		log.Printf("✅ Using extended fallback cached ledgers (%d available)", len(fallbackLedgers))
		if len(fallbackLedgers) > limit {
			fallbackLedgers = fallbackLedgers[:limit]
		}
		return fallbackLedgers, nil
	}

	// Absolutely no data available
	return nil, fmt.Errorf("no XRPL ledger data available and no cache")
}

// getCachedRecentLedgers attempts to retrieve recent ledgers from Redis cache
func (s *EnhancedLedgerService) getCachedRecentLedgers(limit int) ([]map[string]interface{}, error) {
	if s.redisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	ctx := context.Background()

	// First try the simple TTL cache for recent ledgers (fastest)
	cachedData, err := s.redisClient.Get(ctx, "xrp-recent-ledgers").Result()
	if err == nil && cachedData != "" {
		var ledgers []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &ledgers); err == nil {
			// Apply limit to cached data
			if len(ledgers) > limit {
				ledgers = ledgers[:limit]
			}
			return ledgers, nil
		}
	}

	// Return empty result to force fresh data fetch - do NOT automatically clear stream
	// Stream clearing should only be done manually to resolve ID conflicts
	return nil, fmt.Errorf("no cached ledgers found - TTL cache expired")
}

// cacheLedgerData caches individual ledger data in Redis
func (s *EnhancedLedgerService) cacheLedgerData(ledgerIndex int, ledgerData map[string]interface{}) {
	if s.redisClient == nil {
		log.Printf("⚠️ Redis client not initialized, skipping cache for ledger %d", ledgerIndex)
		return
	}

	ctx := context.Background()

	// Marshal ledger data to JSON for storage
	jsonData, err := json.Marshal(ledgerData)
	if err != nil {
		log.Printf("❌ Failed to marshal ledger %d for caching: %v", ledgerIndex, err)
		return
	}

	// Use Redis streams with ledger numbers as IDs - ledger numbers are always increasing and unique
	// Store with ledger number as ID and keep last 17280 entries (~24 hours of ledgers at 5sec intervals)
	result, err := s.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "xrp-ledgers",
		ID:     fmt.Sprintf("%d", ledgerIndex), // Use ledger number as ID - always increasing
		MaxLen: 17280,                          // ~24 hours of ledgers
		Values: map[string]interface{}{
			"data":         string(jsonData),
			"ledger_index": ledgerIndex, // Store ledger index as field for easy querying
		},
	}).Result()

	if err != nil {
		log.Printf("❌ Redis XAdd error for ledger %d: %v, result: %s", ledgerIndex, err, result)

		// Debug: Log current Redis stream state when we get ID conflicts
		s.debugLogRedisStreamState(ledgerIndex)
		return
	}

	// Success - ledger cached to Redis stream
}

// cacheRecentLedgers caches the full recent ledgers result set
func (s *EnhancedLedgerService) cacheRecentLedgers(ledgers []map[string]interface{}) {
	if s.redisClient == nil {
		log.Printf("⚠️ Redis client not initialized, skipping recent ledgers cache")
		return
	}

	ctx := context.Background()

	// Marshal the full ledgers array to JSON
	jsonData, err := json.Marshal(ledgers)
	if err != nil {
		log.Printf("❌ Failed to marshal recent ledgers for caching: %v", err)
		return
	}

	// Cache the full result set with a short TTL (30 seconds)
	// This provides fast access for immediate subsequent requests
	err = s.redisClient.Set(ctx, "xrp-recent-ledgers", string(jsonData), 30*time.Second).Err()
	if err != nil {
		log.Printf("❌ Failed to cache recent ledgers in Redis: %v", err)
		return
	}

	log.Printf("✅ Cached %d recent ledgers in Redis with 30s TTL", len(ledgers))
}

// formatTime converts Unix timestamp to human readable format
func (s *EnhancedLedgerService) formatTime(timestamp int64) string {
	// Convert Ripple epoch to Unix epoch and format
	// Ripple epoch starts at January 1, 2000 00:00:00 UTC
	unixTime := timestamp + 946684800
	return fmt.Sprintf("2025-Jan-01 00:00:%02d.000000000 UTC", unixTime%60)
}

// getActivityLevel returns activity level description based on transaction count
func (s *EnhancedLedgerService) getActivityLevel(txCount int) string {
	switch {
	case txCount == 0:
		return "idle"
	case txCount < 30:
		return "low"
	case txCount < 100:
		return "moderate"
	case txCount < 200:
		return "high"
	default:
		return "very_high"
	}
}

// getCurrentTime returns current Unix timestamp
func (s *EnhancedLedgerService) getCurrentTime() int64 {
	return time.Now().Unix()
}

// GetLedgerTransactionCount gets the number of transactions in a ledger
func (s *EnhancedLedgerService) GetLedgerTransactionCount(ledgerIndex int) (int, error) {
	ledger, err := s.GetRealLedgerData(ledgerIndex)
	if err != nil {
		return 0, err
	}

	if ledger != nil && ledger.Ledger != nil {
		return len(ledger.Ledger.Transactions), nil
	}

	return 0, nil
}

// GetAllLedgerTransactions returns all transactions from a specific ledger
// This method returns ALL transactions (not just AMM), following the reference implementation
func (s *EnhancedLedgerService) GetAllLedgerTransactions(ledgerIndex int) ([]*XRPTransaction, error) {
	log.Printf("🔍 Fetching ledger %d transactions", ledgerIndex)

	ledger, err := s.GetRealLedgerData(ledgerIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger data: %w", err)
	}

	if ledger == nil || ledger.Ledger == nil {
		return []*XRPTransaction{}, nil
	}

	transactions := ledger.Ledger.Transactions
	log.Printf("✅ Retrieved %d transactions from ledger %d", len(transactions), ledgerIndex)

	return transactions, nil
}

// GetLedgerWithAllTransactions returns the complete ledger data including all transactions
// This provides the full ledger structure for detailed analysis
func (s *EnhancedLedgerService) GetLedgerWithAllTransactions(ledgerIndex int) (*XRPLedger, error) {
	log.Printf("🔍 Fetching complete ledger %d", ledgerIndex)

	ledger, err := s.GetRealLedgerData(ledgerIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get complete ledger data: %w", err)
	}

	if ledger != nil && ledger.Ledger != nil {
		transactionCount := len(ledger.Ledger.Transactions)
		log.Printf("✅ Retrieved complete ledger %d with %d transactions", ledgerIndex, transactionCount)

		// Log transaction count summary instead of detailed breakdown
		log.Printf("📊 Ledger %d contains %d transactions", ledgerIndex, len(ledger.Ledger.Transactions))
	}

	return ledger, nil
}

// ProcessLedgerForAMMEvents processes a ledger and extracts AMM-related transactions
func (s *EnhancedLedgerService) ProcessLedgerForAMMEvents(ledgerIndex int) ([]*AMMTransaction, error) {
	ledger, err := s.GetRealLedgerData(ledgerIndex)
	if err != nil {
		return nil, err
	}

	var ammTransactions []*AMMTransaction

	if ledger != nil && ledger.Ledger != nil {
		for _, tx := range ledger.Ledger.Transactions {
			if s.isAMMTransaction(tx) {
				ammTx := s.parseAMMTransaction(tx)
				if ammTx != nil {
					ammTransactions = append(ammTransactions, ammTx)
				}
			}
		}
	}

	// AMM transactions processed
	return ammTransactions, nil
}

// isAMMTransaction checks if a transaction is AMM-related
func (s *EnhancedLedgerService) isAMMTransaction(tx *XRPTransaction) bool {
	if tx == nil {
		return false
	}

	switch tx.TransactionType {
	case "AMMCreate", "AMMDeposit", "AMMWithdraw", "AMMVote", "AMMBid":
		return true
	default:
		return false
	}
}

// parseAMMTransaction converts an XRP transaction to an AMM transaction structure
func (s *EnhancedLedgerService) parseAMMTransaction(tx *XRPTransaction) *AMMTransaction {
	if tx == nil {
		return nil
	}

	ammTx := &AMMTransaction{
		TransactionHash: tx.Hash,
		TransactionType: tx.TransactionType,
		Account:         tx.Account,
		LedgerIndex:     tx.LedgerIndex,
		Timestamp:       int64(tx.Date),
		Metadata:        make(map[string]interface{}),
	}

	// Add transaction-specific metadata
	ammTx.Metadata["fee"] = tx.Fee
	ammTx.Metadata["flags"] = tx.Flags
	ammTx.Metadata["sequence"] = tx.Sequence

	if tx.Amount != nil {
		ammTx.Metadata["amount"] = tx.Amount
	}

	return ammTx
}

// Helper method to count AMM transactions in a list
func (s *EnhancedLedgerService) countAMMTransactions(transactions []*XRPTransaction) int {
	count := 0
	for _, tx := range transactions {
		if s.isAMMTransaction(tx) {
			count++
		}
	}
	return count
}

// ClearRedisStream clears the xrp-ledgers Redis stream to resolve ID conflicts
func (s *EnhancedLedgerService) ClearRedisStream() error {
	if s.redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	ctx := context.Background()

	// Check current stream info before clearing
	log.Printf("🔍 CLEAR: Checking current stream info before clearing...")
	streamInfo, err := s.redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		log.Printf("🔍 CLEAR: Error getting stream info (stream may not exist): %v", err)
	} else {
		log.Printf("🔍 CLEAR: Stream Length: %d", streamInfo.Length)
		log.Printf("🔍 CLEAR: First Entry ID: %s", streamInfo.FirstEntry.ID)
		log.Printf("🔍 CLEAR: Last Entry ID: %s", streamInfo.LastEntry.ID)
	}

	// Clear the stream by deleting it entirely
	log.Printf("🔍 CLEAR: Clearing Redis stream 'xrp-ledgers'...")
	result, err := s.redisClient.Del(ctx, "xrp-ledgers").Result()
	if err != nil {
		log.Printf("❌ CLEAR: Failed to delete stream: %v", err)
		return fmt.Errorf("failed to delete Redis stream: %w", err)
	}

	if result == 1 {
		log.Printf("✅ CLEAR: Successfully deleted 'xrp-ledgers' stream")
	} else {
		log.Printf("⚠️ CLEAR: Stream 'xrp-ledgers' did not exist or was already empty")
	}

	// Verify the stream is cleared
	log.Printf("🔍 CLEAR: Verifying stream is cleared...")
	streamInfo, err = s.redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		log.Printf("✅ CLEAR: Stream 'xrp-ledgers' no longer exists: %v", err)
	} else {
		log.Printf("⚠️ CLEAR: Stream still exists with length: %d", streamInfo.Length)
	}

	log.Printf("✅ CLEAR: Redis stream cleared successfully. New ledger data will use auto-generated timestamp IDs.")
	return nil
}

// debugLogRedisStreamState logs the current Redis stream state for debugging ID conflicts
func (s *EnhancedLedgerService) debugLogRedisStreamState(currentLedgerIndex int) {
	// DEBUG logging removed Starting Redis stream debug for ledger %d", currentLedgerIndex)

	if s.redisClient == nil {
		// DEBUG logging removed Redis client not available for stream debugging")
		return
	}

	// DEBUG logging removed Redis client available, attempting stream info...")

	ctx := context.Background()

	// Get stream info (debug logging removed)
	_, err := s.redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		// DEBUG logging removed Error getting stream info: %v", err)
		return
	}

	// DEBUG logging removed Redis Stream Info for ledger %d conflict:", currentLedgerIndex)
	// DEBUG logging removed Stream Length: %d", streamInfo.Length)
	// DEBUG logging removed Last Entry ID in stream: %s", streamInfo.LastEntry.ID)
	// DEBUG logging removed First Entry ID in stream: %s", streamInfo.FirstEntry.ID)

	// Get the last 5 entries to see what's at the end of the stream
	result, err := s.redisClient.XRevRangeN(ctx, "xrp-ledgers", "+", "-", 5).Result()
	if err != nil {
		// DEBUG logging removed Error reading last entries: %v", err)
		return
	}

	// Process stream entries (debug logging removed)
	for _, entry := range result {
		// Extract ledger index from data (debug logging removed)
		ledgerIndexFromStream := "unknown"
		if data, ok := entry.Values["data"]; ok {
			dataStr := data.(string)
			// Simple check for ledger_index in JSON data
			if len(dataStr) > 50 {
				// Look for ledger_index pattern in JSON
				if strings.Contains(dataStr, `"ledger_index":`) {
					// This is a very basic extraction - in real implementation we'd parse JSON
					ledgerIndexFromStream = "found_in_data"
				}
			}
		}
		if ledgerIndex, ok := entry.Values["ledger_index"]; ok {
			ledgerIndexFromStream = fmt.Sprintf("%v", ledgerIndex)
		}

		// Stream entry logging (debug logging removed)
		_ = ledgerIndexFromStream // Prevent unused variable warning
	}

	// Attempted to add ledger (debug logging removed)
}

// GetLatestCachedLedgerIndex gets the most recent ledger index from Redis cache
func (s *EnhancedLedgerService) GetLatestCachedLedgerIndex() (int, error) {
	s.stateMutex.RLock()
	defer s.stateMutex.RUnlock()

	if s.redisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	// Create context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get the most recent entry from Redis stream with timeout
	result, err := s.redisClient.XRevRangeN(ctx, "xrp-ledgers", "+", "-", 1).Result()
	if err != nil {
		log.Printf("⚠️ Redis stream read failed (stream may be empty): %v", err)
		return 0, fmt.Errorf("failed to get latest cached ledger: %w", err)
	}

	if len(result) == 0 {
		log.Printf("⚠️ Redis stream is empty - no cached ledgers found")
		return 0, fmt.Errorf("no cached ledgers found")
	}

	// Parse the ledger index from the stream ID (ledger number)
	entry := result[0]
	ledgerIndexStr := strings.Split(entry.ID, "-")[0]

	ledgerIndex, err := strconv.Atoi(ledgerIndexStr)
	if err != nil {
		log.Printf("⚠️ Failed to parse ledger index from stream ID '%s': %v", entry.ID, err)
		return 0, fmt.Errorf("failed to parse ledger index from cache: %w", err)
	}

	return ledgerIndex, nil
}

// IsCacheFresh checks if the Redis cache is fresh compared to current XRPL state
func (s *EnhancedLedgerService) IsCacheFresh() (bool, int, error) {
	// Use read lock for most operations, only lock for writing when updating lastCacheCheckTime
	s.stateMutex.RLock()
	lastCheck := s.lastCacheCheckTime
	interval := s.cacheCheckInterval
	s.stateMutex.RUnlock()

	// Check if we need to perform a cache freshness check
	if time.Since(lastCheck) < interval {
		// Return cached result if check was recent
		return true, 0, nil
	}

	// Update the check time with write lock
	s.stateMutex.Lock()
	s.lastCacheCheckTime = time.Now()
	s.stateMutex.Unlock()

	// Get current XRPL ledger index with timeout
	currentXRPLIndex, err := s.GetCurrentLedgerIndex()
	if err != nil {
		log.Printf("⚠️ Failed to get current XRPL ledger index: %v", err)
		// Don't block the request - assume cache is stale and fetch fresh data
		return false, 0, nil
	}

	// Get latest cached ledger index with graceful fallback
	latestCachedIndex, err := s.GetLatestCachedLedgerIndex()
	if err != nil {
		log.Printf("⚠️ Failed to get latest cached ledger index: %v", err)
		// Don't block the request - assume cache is stale and fetch fresh data
		return false, 0, nil
	}

	// Calculate the gap between XRPL and cache
	gap := currentXRPLIndex - latestCachedIndex

	log.Printf("🔍 Cache freshness check: XRPL=%d, Cache=%d, Gap=%d, Threshold=%d",
		currentXRPLIndex, latestCachedIndex, gap, s.cacheFreshnessThreshold)

	// Cache is fresh if gap is within threshold
	isFresh := gap <= s.cacheFreshnessThreshold

	if !isFresh {
		log.Printf("⚠️ Cache is stale: gap=%d > threshold=%d", gap, s.cacheFreshnessThreshold)
	} else {
		log.Printf("✅ Cache is fresh: gap=%d <= threshold=%d", gap, s.cacheFreshnessThreshold)
	}

	return isFresh, gap, nil
}

// UpdateLatestProcessedLedger updates the latest processed ledger index
func (s *EnhancedLedgerService) UpdateLatestProcessedLedger(ledgerIndex int) {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if ledgerIndex > s.latestProcessedLedgerIndex {
		s.latestProcessedLedgerIndex = ledgerIndex
		log.Printf("📊 Updated latest processed ledger index: %d", ledgerIndex)
	}
}

// GetLatestProcessedLedger returns the latest processed ledger index
func (s *EnhancedLedgerService) GetLatestProcessedLedger() int {
	s.stateMutex.RLock()
	defer s.stateMutex.RUnlock()
	return s.latestProcessedLedgerIndex
}

// InvalidateCache clears the Redis cache and forces fresh data fetch
func (s *EnhancedLedgerService) InvalidateCache() error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	ctx := context.Background()

	// Clear the Redis stream
	_, err := s.redisClient.Del(ctx, "xrp-ledgers").Result()
	if err != nil {
		log.Printf("❌ Failed to clear Redis cache: %v", err)
		return fmt.Errorf("failed to clear Redis cache: %w", err)
	}

	// Clear the recent ledgers cache
	_, err = s.redisClient.Del(ctx, "xrp-recent-ledgers").Result()
	if err != nil {
		log.Printf("❌ Failed to clear recent ledgers cache: %v", err)
		return fmt.Errorf("failed to clear recent ledgers cache: %w", err)
	}

	log.Printf("🗑️ Cache invalidated successfully")
	return nil
}

// WarmCacheAfterInvalidation warms the cache with fresh data after invalidation
func (s *EnhancedLedgerService) WarmCacheAfterInvalidation(limit int) error {
	log.Printf("🔥 Warming cache with fresh XRPL data after invalidation...")

	// Fetch fresh ledgers to warm the cache
	ledgers, err := s.GetRecentLedgers(limit)
	if err != nil {
		return fmt.Errorf("failed to fetch ledgers for cache warming: %w", err)
	}

	log.Printf("✅ Cache warmed with %d fresh ledgers", len(ledgers))
	return nil
}

// InvalidateAndWarmCache invalidates cache and immediately warms it with fresh data
func (s *EnhancedLedgerService) InvalidateAndWarmCache(limit int) error {
	// Invalidate the cache
	if err := s.InvalidateCache(); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	// Warm the cache with fresh data
	if err := s.WarmCacheAfterInvalidation(limit); err != nil {
		return fmt.Errorf("failed to warm cache after invalidation: %w", err)
	}

	return nil
}

// ValidateLedgerData checks if ledger data is valid (no zero ledger numbers, etc.)
func (s *EnhancedLedgerService) ValidateLedgerData(ledger map[string]interface{}) bool {
	// Check for zero ledger index
	if ledgerIndex, ok := ledger["ledger_index"]; ok {
		if index, isInt := ledgerIndex.(int); isInt && index <= 0 {
			log.Printf("⚠️ Invalid ledger data: zero or negative ledger_index=%d", index)
			return false
		}
		if index, isFloat := ledgerIndex.(float64); isFloat && index <= 0 {
			log.Printf("⚠️ Invalid ledger data: zero or negative ledger_index=%f", index)
			return false
		}
	}

	// Check for empty ledger hash
	if ledgerHash, ok := ledger["ledger_hash"]; ok {
		if hash, isString := ledgerHash.(string); isString && (hash == "" || hash == "0") {
			log.Printf("⚠️ Invalid ledger data: empty or zero ledger_hash=%s", hash)
			return false
		}
	}

	return true
}
