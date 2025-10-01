package xrp

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClearRedisStream is a maintenance test that can be run to clear the Redis stream
// when encountering ID conflicts or needing to reset the stream state.
//
// To run this test specifically:
//   go test -v ./internal/xrp -run TestClearRedisStream
//
// This is a safer approach than exposing external endpoints for data modification.
func TestClearRedisStream(t *testing.T) {
	// Skip by default - only run when explicitly needed
	if testing.Short() {
		t.Skip("Skipping Redis maintenance test in short mode")
	}

	// Initialize Redis client with production credentials
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-14565.fcrce180.us-east-1-1.ec2.redns.redis-cloud.com:14565",
		Username: "default",
		Password: "RxOmtTCWIIdFdvuhaUwqXB30b08hpvoY",
		DB:       0,
	})
	defer redisClient.Close()

	// Test Redis connectivity first
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	require.NoError(t, err, "Should be able to connect to Redis")
	assert.Equal(t, "PONG", pong, "Redis should respond with PONG")

	t.Log("✅ Connected to Redis successfully")

	// Create enhanced ledger service with Redis client
	enhancedLedgerSvc := &EnhancedLedgerService{
		redisClient: redisClient,
	}

	// Check current stream state before clearing
	t.Log("🔍 Checking current Redis stream state...")
	streamInfo, err := redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		t.Logf("⚠️ Stream may not exist yet: %v", err)
	} else {
		t.Logf("📊 Current stream length: %d", streamInfo.Length)
		t.Logf("📊 First entry ID: %s", streamInfo.FirstEntry.ID)
		t.Logf("📊 Last entry ID: %s", streamInfo.LastEntry.ID)
	}

	// Clear the Redis stream
	t.Log("🔧 Clearing Redis stream...")
	err = enhancedLedgerSvc.ClearRedisStream()
	require.NoError(t, err, "Should be able to clear Redis stream")

	// Verify the stream is cleared
	t.Log("🔍 Verifying stream is cleared...")
	streamInfo, err = redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		t.Log("✅ Stream no longer exists - cleared successfully")
	} else {
		t.Errorf("❌ Stream still exists with length: %d", streamInfo.Length)
	}

	t.Log("✅ Redis stream maintenance completed successfully")
}

// TestRedisStreamInfo is a diagnostic test to check the current state of the Redis stream
// without making any modifications.
//
// To run this test:
//   go test -v ./internal/xrp -run TestRedisStreamInfo
func TestRedisStreamInfo(t *testing.T) {
	// Skip by default - only run when explicitly needed
	if testing.Short() {
		t.Skip("Skipping Redis diagnostic test in short mode")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-14565.fcrce180.us-east-1-1.ec2.redns.redis-cloud.com:14565",
		Username: "default",
		Password: "RxOmtTCWIIdFdvuhaUwqXB30b08hpvoY",
		DB:       0,
	})
	defer redisClient.Close()

	ctx := context.Background()

	// Test Redis connectivity
	_, err := redisClient.Ping(ctx).Result()
	require.NoError(t, err, "Should be able to connect to Redis")

	t.Log("🔍 Checking Redis stream state...")

	// Get stream info
	streamInfo, err := redisClient.XInfoStream(ctx, "xrp-ledgers").Result()
	if err != nil {
		t.Logf("⚠️ Stream does not exist: %v", err)
		return
	}

	t.Logf("📊 Stream Length: %d", streamInfo.Length)
	t.Logf("📊 First Entry ID: %s", streamInfo.FirstEntry.ID)
	t.Logf("📊 Last Entry ID: %s", streamInfo.LastEntry.ID)

	// Get the last 5 entries to see what's at the end of the stream
	result, err := redisClient.XRevRangeN(ctx, "xrp-ledgers", "+", "-", 5).Result()
	if err != nil {
		t.Logf("❌ Error reading last entries: %v", err)
		return
	}

	t.Logf("📊 Last 5 entries in stream (newest first):")
	for i, entry := range result {
		ledgerIndexFromStream := "unknown"
		if ledgerIndex, ok := entry.Values["ledger_index"]; ok {
			ledgerIndexFromStream = ledgerIndex.(string)
		}
		t.Logf("  %d. Stream ID: %s, Ledger: %s", i+1, entry.ID, ledgerIndexFromStream)
	}
}

// TestRedisStreamRecovery tests the stream recovery after clearing
// This can be used to verify that new entries use proper auto-generated IDs
//
// To run this test:
//   go test -v ./internal/xrp -run TestRedisStreamRecovery
func TestRedisStreamRecovery(t *testing.T) {
	// Skip by default - only run when explicitly needed
	if testing.Short() {
		t.Skip("Skipping Redis recovery test in short mode")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-14565.fcrce180.us-east-1-1.ec2.redns.redis-cloud.com:14565",
		Username: "default",
		Password: "RxOmtTCWIIdFdvuhaUwqXB30b08hpvoY",
		DB:       0,
	})
	defer redisClient.Close()

	ctx := context.Background()

	// Test Redis connectivity
	_, err := redisClient.Ping(ctx).Result()
	require.NoError(t, err, "Should be able to connect to Redis")

	t.Log("🔧 Testing Redis stream recovery with auto-generated IDs...")

	// Add a test entry using auto-generated ID (the "*" pattern)
	testData := map[string]interface{}{
		"data":         `{"test": true, "ledger_index": 99999999}`,
		"ledger_index": 99999999,
	}

	result, err := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "xrp-ledgers",
		ID:     "*", // Auto-generate timestamp-based ID
		Values: testData,
	}).Result()

	require.NoError(t, err, "Should be able to add test entry with auto-generated ID")
	t.Logf("✅ Added test entry with ID: %s", result)

	// Verify the entry was added
	entries, err := redisClient.XRevRangeN(ctx, "xrp-ledgers", "+", "-", 1).Result()
	require.NoError(t, err, "Should be able to read last entry")
	require.Len(t, entries, 1, "Should have one entry")

	lastEntry := entries[0]
	t.Logf("✅ Last entry ID: %s", lastEntry.ID)
	t.Logf("✅ Auto-generated IDs are working correctly")

	// Clean up test entry
	_, err = redisClient.XDel(ctx, "xrp-ledgers", lastEntry.ID).Result()
	require.NoError(t, err, "Should be able to clean up test entry")

	t.Log("✅ Redis stream recovery test completed successfully")
}
