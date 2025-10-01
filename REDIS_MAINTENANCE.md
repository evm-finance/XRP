# Redis Stream Maintenance

This document describes how to maintain the Redis stream used for XRP ledger data caching.

## Problem

When Redis stream entries use sequential integer IDs (like ledger indexes), you may encounter the error:
```
WRONGTYPE Operation against a key holding the wrong kind of value
```
or
```
ERR The ID specified in XADD is equal or smaller than the target stream top item
```

This happens when trying to add entries with IDs that are not greater than the last entry's ID.

## Solution

Use the test-based maintenance approach to safely clear and reset the Redis stream.

## Available Tests

### 1. Check Stream State (Diagnostic)
```bash
cd backends/main
go test -v ./internal/xrp -run TestRedisStreamInfo
```

This test shows the current state of the Redis stream without making any changes:
- Stream length
- First and last entry IDs
- Recent entries

### 2. Clear Redis Stream (Maintenance)
```bash
cd backends/main
go test -v ./internal/xrp -run TestClearRedisStream
```

This test safely clears the Redis stream:
- Connects to Redis using the same credentials as the server
- Shows current stream state
- Deletes the entire stream
- Verifies the stream is cleared

### 3. Test Stream Recovery
```bash
cd backends/main
go test -v ./internal/xrp -run TestRedisStreamRecovery
```

This test verifies that auto-generated IDs work correctly after clearing:
- Adds a test entry with auto-generated ID
- Verifies the entry was added successfully
- Cleans up the test entry

## Usage Workflow

1. **Check current state**:
   ```bash
   go test -v ./internal/xrp -run TestRedisStreamInfo
   ```

2. **Clear the stream** (if needed):
   ```bash
   go test -v ./internal/xrp -run TestClearRedisStream
   ```

3. **Verify recovery**:
   ```bash
   go test -v ./internal/xrp -run TestRedisStreamRecovery
   ```

4. **Restart your server** to begin using the cleared stream with auto-generated IDs.

## Why This Approach?

✅ **Secure**: No external API endpoints that could be misused
✅ **Safe**: Uses the same Redis connection as production
✅ **Testable**: Verifies operations and can be run repeatedly
✅ **Documented**: Clear logs and verification steps
✅ **Following patterns**: Similar to database setup/teardown tests

## After Clearing

Once the stream is cleared:
- New ledger data will use ledger numbers as stream IDs
- Ledger numbers are always increasing and unique, making them perfect for Redis streams
- The server will automatically create new entries with ledger number IDs
- ID conflicts should be resolved
- Stream will continue working normally with proper ordering

## Notes

- Tests are skipped in short mode (`go test -short`) to prevent accidental execution
- All tests include connectivity verification before performing operations
- Redis credentials are the same as used by the production server
- Tests include proper cleanup and verification steps 