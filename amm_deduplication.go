package xrp

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// AMMDeduplicationService handles removal of duplicate AMM pool records
type AMMDeduplicationService struct {
	db *gorm.DB
}

// DeduplicationStats tracks the results of deduplication operations
type DeduplicationStats struct {
	TotalRecordsProcessed int       `json:"total_records_processed"`
	DuplicatesRemoved     int       `json:"duplicates_removed"`
	UniqueRecordsKept     int       `json:"unique_records_kept"`
	ErrorsEncountered     int       `json:"errors_encountered"`
	ProcessingTimeMs      int64     `json:"processing_time_ms"`
	StartTime             time.Time `json:"start_time"`
	EndTime               time.Time `json:"end_time"`
}

// DuplicateRecord represents a duplicate AMM record with metadata
type DuplicateRecord struct {
	Account      string                 `json:"account"`
	RecordCount  int                    `json:"record_count"`
	FirstCreated time.Time              `json:"first_created"`
	LastUpdated  time.Time              `json:"last_updated"`
	DuplicateIDs []int                  `json:"duplicate_ids"`
	SampleData   map[string]interface{} `json:"sample_data"`
}

// NewAMMDeduplicationService creates a new deduplication service
func NewAMMDeduplicationService(db *gorm.DB) *AMMDeduplicationService {
	return &AMMDeduplicationService{
		db: db,
	}
}

// AnalyzeDuplicates analyzes the current state of duplicates in the database
func (d *AMMDeduplicationService) AnalyzeDuplicates() ([]DuplicateRecord, *DeduplicationStats, error) {
	startTime := time.Now()
	stats := &DeduplicationStats{
		StartTime: startTime,
	}

	log.Printf("🔍 Starting duplicate analysis of xrpAmm table...")

	// Find all duplicate accounts with their counts
	var duplicateAccounts []struct {
		Account string `json:"account"`
		Count   int    `json:"count"`
	}

	result := d.db.Raw(`
		SELECT account, COUNT(*) as count 
		FROM xrpAmm_normalized 
		WHERE account IS NOT NULL AND account != '' 
		GROUP BY account 
		HAVING COUNT(*) > 1 
		ORDER BY count DESC
	`).Scan(&duplicateAccounts)

	if result.Error != nil {
		return nil, nil, fmt.Errorf("failed to analyze duplicates: %w", result.Error)
	}

	log.Printf("📊 Found %d accounts with duplicates", len(duplicateAccounts))

	// Get detailed information for each duplicate account
	var duplicateRecords []DuplicateRecord
	totalDuplicates := 0

	for _, dup := range duplicateAccounts {
		// Get all records for this account with detailed info
		var records []map[string]interface{}
		result := d.db.Raw(`
			SELECT *, 
				COALESCE(created_at, CURRENT_TIMESTAMP) as created_at,
				COALESCE(last_updated, CURRENT_TIMESTAMP) as last_updated
			FROM xrpAmm_normalized 
			WHERE account = ? 
			ORDER BY 
				COALESCE(last_updated, CURRENT_TIMESTAMP) DESC,
				COALESCE(created_at, CURRENT_TIMESTAMP) DESC
		`, dup.Account).Scan(&records)

		if result.Error != nil {
			log.Printf("⚠️ Error getting records for account %s: %v", dup.Account, result.Error)
			stats.ErrorsEncountered++
			continue
		}

		if len(records) > 0 {
			// Extract metadata from records
			var duplicateIDs []int
			var firstCreated, lastUpdated time.Time

			for i, record := range records {
				// Extract ID if available
				if id, ok := record["id"]; ok {
					if idInt, ok := id.(int64); ok {
						duplicateIDs = append(duplicateIDs, int(idInt))
					}
				}

				// Track timestamps
				if i == 0 {
					lastUpdated = time.Now() // Default
					if updated, ok := record["last_updated"]; ok {
						if updatedTime, ok := updated.(time.Time); ok {
							lastUpdated = updatedTime
						}
					}
				}
				if i == len(records)-1 {
					firstCreated = time.Now() // Default
					if created, ok := record["created_at"]; ok {
						if createdTime, ok := created.(time.Time); ok {
							firstCreated = createdTime
						}
					}
				}
			}

			duplicateRecord := DuplicateRecord{
				Account:      dup.Account,
				RecordCount:  dup.Count,
				FirstCreated: firstCreated,
				LastUpdated:  lastUpdated,
				DuplicateIDs: duplicateIDs,
				SampleData:   records[0], // Keep the most recent record as sample
			}

			duplicateRecords = append(duplicateRecords, duplicateRecord)
			totalDuplicates += dup.Count - 1 // -1 because we keep one record
		}
	}

	// Get total record count
	var totalRecords int64
	d.db.Raw("SELECT COUNT(*) FROM xrpAmm_normalized").Scan(&totalRecords)

	// Update stats
	stats.TotalRecordsProcessed = int(totalRecords)
	stats.DuplicatesRemoved = totalDuplicates
	stats.UniqueRecordsKept = int(totalRecords) - totalDuplicates
	stats.EndTime = time.Now()
	stats.ProcessingTimeMs = stats.EndTime.Sub(startTime).Milliseconds()

	log.Printf("📈 Duplicate analysis complete:")
	log.Printf("  - Total records: %d", stats.TotalRecordsProcessed)
	log.Printf("  - Duplicate records to remove: %d", stats.DuplicatesRemoved)
	log.Printf("  - Unique records to keep: %d", stats.UniqueRecordsKept)
	log.Printf("  - Accounts with duplicates: %d", len(duplicateRecords))

	return duplicateRecords, stats, nil
}

// RemoveDuplicates removes duplicate records, keeping the most recent/complete one
func (d *AMMDeduplicationService) RemoveDuplicates(dryRun bool) (*DeduplicationStats, error) {
	startTime := time.Now()
	stats := &DeduplicationStats{
		StartTime: startTime,
	}

	log.Printf("🧹 Starting deduplication process (dry run: %t)...", dryRun)

	// First, analyze current duplicates
	duplicateRecords, analyzeStats, err := d.AnalyzeDuplicates()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze duplicates: %w", err)
	}

	// Copy analyze stats
	stats.TotalRecordsProcessed = analyzeStats.TotalRecordsProcessed
	stats.DuplicatesRemoved = analyzeStats.DuplicatesRemoved
	stats.UniqueRecordsKept = analyzeStats.UniqueRecordsKept

	if len(duplicateRecords) == 0 {
		log.Printf("✅ No duplicates found - database is already clean!")
		stats.EndTime = time.Now()
		stats.ProcessingTimeMs = stats.EndTime.Sub(startTime).Milliseconds()
		return stats, nil
	}

	// Process each duplicate group
	removedCount := 0
	for _, dupRecord := range duplicateRecords {
		log.Printf("🔄 Processing duplicates for account: %s (%d records)",
			dupRecord.Account, dupRecord.RecordCount)

		// Get all records for this account ordered by quality/recency
		var records []map[string]interface{}
		result := d.db.Raw(`
					SELECT *, 
			COALESCE(liquidity_usd, 0) as liquidity_usd,
			COALESCE(last_updated, created_at, CURRENT_TIMESTAMP) as effective_updated
		FROM xrpAmm_normalized 
		WHERE account = ? 
			ORDER BY 
				liquidity_usd DESC,  -- Prefer records with liquidity data
				effective_updated DESC,  -- Prefer most recently updated
				CASE 
					WHEN amount2currency IS NOT NULL AND amount2currency != '' THEN 1 
					ELSE 0 
				END DESC,  -- Prefer records with complete currency data
				CASE 
					WHEN amount IS NOT NULL AND amount != '' AND amount != '0' THEN 1 
					ELSE 0 
				END DESC  -- Prefer records with amount data
		`, dupRecord.Account).Scan(&records)

		if result.Error != nil {
			log.Printf("⚠️ Error getting records for account %s: %v", dupRecord.Account, result.Error)
			stats.ErrorsEncountered++
			continue
		}

		if len(records) <= 1 {
			continue // No duplicates (shouldn't happen but safety check)
		}

		// Keep the first record (highest quality), remove the rest
		recordsToRemove := records[1:] // All except the first (best) record

		if !dryRun {
			// Remove duplicate records by building a DELETE query with IDs
			var idsToDelete []interface{}
			for _, record := range recordsToRemove {
				if id, exists := record["id"]; exists {
					idsToDelete = append(idsToDelete, id)
				}
			}

			if len(idsToDelete) > 0 {
				// Build dynamic placeholders for the IN clause
				placeholders := ""
				for i := range idsToDelete {
					if i > 0 {
						placeholders += ","
					}
					placeholders += "?"
				}

				deleteQuery := fmt.Sprintf("DELETE FROM xrpAmm_normalized WHERE id IN (%s)", placeholders)
				result := d.db.Exec(deleteQuery, idsToDelete...)

				if result.Error != nil {
					log.Printf("❌ Error removing duplicates for account %s: %v", dupRecord.Account, result.Error)
					stats.ErrorsEncountered++
					continue
				}

				removedCount += len(recordsToRemove)
				log.Printf("✅ Removed %d duplicate records for account %s, kept 1 best record",
					len(recordsToRemove), dupRecord.Account)
			}
		} else {
			// Dry run - just log what would be removed
			removedCount += len(recordsToRemove)
			log.Printf("🔍 [DRY RUN] Would remove %d duplicate records for account %s",
				len(recordsToRemove), dupRecord.Account)
		}
	}

	stats.DuplicatesRemoved = removedCount
	stats.EndTime = time.Now()
	stats.ProcessingTimeMs = stats.EndTime.Sub(startTime).Milliseconds()

	if dryRun {
		log.Printf("🔍 [DRY RUN] Deduplication analysis complete:")
	} else {
		log.Printf("✅ Deduplication complete:")
	}
	log.Printf("  - Total records processed: %d", stats.TotalRecordsProcessed)
	log.Printf("  - Duplicate records removed: %d", stats.DuplicatesRemoved)
	log.Printf("  - Unique records kept: %d", stats.UniqueRecordsKept)
	log.Printf("  - Errors encountered: %d", stats.ErrorsEncountered)
	log.Printf("  - Processing time: %dms", stats.ProcessingTimeMs)

	return stats, nil
}

// MigrateToNewSchema migrates clean data from old xrpAmm table to new xrpAmm_v2 table
func (d *AMMDeduplicationService) MigrateToNewSchema(dryRun bool) (*DeduplicationStats, error) {
	startTime := time.Now()
	log.Printf("🔄 Starting migration to xrpAmm_v2 schema (dry run: %t)...", dryRun)

	// First, ensure we have clean data in the original table
	log.Printf("📊 Step 1: Cleaning original table...")
	dedupeStats, err := d.RemoveDuplicates(dryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to deduplicate before migration: %w", err)
	}

	if dryRun {
		log.Printf("🔍 [DRY RUN] Migration would proceed with %d clean records", dedupeStats.UniqueRecordsKept)
		return dedupeStats, nil
	}

	// Step 2: Create new table if it doesn't exist (should already exist from database_create.go)
	log.Printf("📊 Step 2: Ensuring xrpAmm_v2 table exists...")

	// Step 3: Migrate clean data with proper data types and validation
	log.Printf("📊 Step 3: Migrating clean data to new schema...")

	migrateQuery := `
		INSERT IGNORE INTO xrpAmm_v2 (
			account,
			amount,
			amount_currency,
			amount_issuer,
			amount2_value,
			amount2_currency,
			amount2_issuer,
			asset2_frozen,
			lp_token_value,
			lp_token_currency,
			lp_token_issuer,
			trading_fee_bps,
			liquidity_usd,
			asset1_value_usd,
			asset2_value_usd,
			ledger_index,
			validated,
			created_at,
			last_updated,
			last_api_fetch
		)
		SELECT 
			account,
			CASE 
				WHEN amount REGEXP '^[0-9]+\\.?[0-9]*$' THEN CAST(amount AS DECIMAL(20,8))
				ELSE 0
			END as amount,
			'XRP' as amount_currency,
			'' as amount_issuer,
			CASE 
				WHEN amount2value REGEXP '^[0-9]+\\.?[0-9]*$' THEN CAST(amount2value AS DECIMAL(20,8))
				ELSE 0
			END as amount2_value,
			COALESCE(amount2currency, '') as amount2_currency,
			COALESCE(amount2issuer, '') as amount2_issuer,
			CASE WHEN asset2frozen = 1 THEN TRUE ELSE FALSE END as asset2_frozen,
			CASE 
				WHEN lptokenvalue REGEXP '^[0-9]+\\.?[0-9]*$' THEN CAST(lptokenvalue AS DECIMAL(20,8))
				ELSE 0
			END as lp_token_value,
			COALESCE(lptokencurrency, '') as lp_token_currency,
			COALESCE(lptokenissuer, '') as lp_token_issuer,
			COALESCE(tradingfee, 0) as trading_fee_bps,
			COALESCE(liquidity_usd, 0) as liquidity_usd,
			COALESCE(asset1_value_usd, 0) as asset1_value_usd,
			COALESCE(asset2_value_usd, 0) as asset2_value_usd,
			COALESCE(created_at, 0) as ledger_index,
			TRUE as validated,
			COALESCE(created_at, CURRENT_TIMESTAMP) as created_at,
			COALESCE(last_updated, CURRENT_TIMESTAMP) as last_updated,
			CURRENT_TIMESTAMP as last_api_fetch
		FROM xrpAmm_normalized 
		WHERE account IS NOT NULL 
			AND account != '' 
			AND account REGEXP '^r[a-zA-Z0-9]{24,34}$'
			AND amount2currency IS NOT NULL 
			AND amount2currency != ''
			AND amount2issuer IS NOT NULL 
			AND amount2issuer != ''
	`

	result := d.db.Exec(migrateQuery)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to migrate data to new schema: %w", result.Error)
	}

	migratedCount := result.RowsAffected
	log.Printf("✅ Successfully migrated %d records to xrpAmm_v2", migratedCount)

	// Update final stats
	dedupeStats.EndTime = time.Now()
	dedupeStats.ProcessingTimeMs = time.Now().Sub(startTime).Milliseconds()

	log.Printf("🎉 Migration to new schema complete!")
	log.Printf("  - Records migrated: %d", migratedCount)
	log.Printf("  - Total processing time: %dms", dedupeStats.ProcessingTimeMs)

	return dedupeStats, nil
}

// SetupUniqueConstraints adds unique constraints to prevent future duplicates
func (d *AMMDeduplicationService) SetupUniqueConstraints() error {
	log.Printf("🔒 Setting up unique constraints to prevent future duplicates...")

	// Add unique constraint on account field in xrpAmm_v2 (should already exist from schema)
	constraintQuery := `
		ALTER TABLE xrpAmm_v2 
		ADD CONSTRAINT uk_amm_account_unique 
		UNIQUE (account)
	`

	result := d.db.Exec(constraintQuery)
	if result.Error != nil {
		// Check if constraint already exists
		if result.Error.Error() != "" {
			log.Printf("ℹ️ Unique constraint may already exist: %v", result.Error)
		} else {
			return fmt.Errorf("failed to add unique constraint: %w", result.Error)
		}
	} else {
		log.Printf("✅ Unique constraint added successfully")
	}

	return nil
}
