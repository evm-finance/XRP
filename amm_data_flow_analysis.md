# XRP AMM Data Flow & Troubleshooting Guide

## 1. Token Fetcher Step
- **Purpose:** Fetches all unique tokens from the XRPL network.
- **Process:**
  - Scans the XRPL for all issued tokens (IOUs).
  - Stores each unique token (by currency code + issuer) in the `xrpToken` table.
- **Result:**
  - ~7,200 unique tokens stored in the token database.
- **Potential Issues:**
  - Duplicate tokens if not deduplicated by (currency, issuer).
  - Incomplete or malformed token data if fetcher logic is incorrect.

### Key Files for Token Fetching:
- `cmd/collect_xrp_data/collect_xrp_data.go` - Main token collection script
- `internal/xrp/token_service.go` - Token service for database operations
- `internal/xrp/token_model.go` - Token data models and structures

## 2. AMM Pool Discovery Step
- **Purpose:** Discover all AMM pools on XRPL using the list of known tokens.
- **Process:**
  - For each token (or token pair), queries the XRPL for AMM pools.
  - Each discovered pool is stored in the `xrpAmm` table.
  - Each pool record includes: `account`, `amount`, `amount2currency`, `amount2issuer`, `amount2value`, etc.
- **Result:**
  - Over 2,000 pool records inserted into the `xrpAmm` table.
  - **But:** Only 68 unique pool accounts (distinct `account` values).
- **Potential Issues:**
  - **Duplication:** Multiple records per pool account (e.g., for different tokens, timestamps, or duplicate discovery runs).
  - **Empty/Null Accounts:** Some records have empty or null `account` fields (data quality issue).
  - **Schema:** Table is not normalized; historical or duplicate entries are not deduplicated.
  - **Discovery Logic:** If the discovery script does not check for existing pools before inserting, duplicates will accumulate.

### Key Files for AMM Pool Discovery:
- `cmd/amm_discovery.go` - Main AMM discovery command
- `internal/xrp/amm_service.go` - AMM service for pool operations
- `internal/xrp/amm_model.go` - AMM pool data models
- `cmd_utils/run_amm_pool_discovery.go` - AMM pool discovery utility (deprecated)

## 3. Liquidity Update Step
- **Purpose:** Calculate and update the USD liquidity for each pool.
- **Process:**
  - For each unique pool account, fetches all records from `xrpAmm`.
  - Calculates liquidity using the latest or best available data.
  - Updates all rows for that account with the new liquidity values.
- **Result:**
  - All rows for each pool account are updated (not just one per pool).
  - Test output: 68 unique pools, 1,171 total rows updated.
- **Potential Issues:**
  - **Bulk Update:** All records for a pool are updated, even if only one is current/valid.
  - **Historical Data:** Old or duplicate records remain in the table, inflating row count.
  - **Scan Errors:** If schema changes or NULLs are not handled, scan errors can occur (now fixed).

### Key Files for Liquidity Updates:
- `cmd/liquidity_calculation.go` - Main liquidity calculation command
- `internal/xrp/amm_liquidity_service.go` - Liquidity calculation service
- `internal/xrp/amm_liquidity_test.go` - Test suite for liquidity operations
- `internal/xrp/price_service.go` - Price service for USD conversions

## 4. Database & Configuration Files:
- `internal/database/database_create.go` - Database schema and table definitions
- `configs/mysql.go` - Database connection configuration
- `configs/eth.go` - XRPL connection configuration

## 5. GraphQL & API Files:
- `internal/graph/handler/amm_handler.go` - GraphQL handlers for AMM data
- `internal/graph/resolver/amm_resolver.go` - GraphQL resolvers
- `internal/graph/schema/amm.graphqls` - GraphQL schema definitions

## 6. Relationship Between Steps
- **Token Fetcher → Pool Discovery:**
  - The list of tokens is used to search for all possible AMM pools.
  - If the token list has duplicates, pool discovery may be redundant or inefficient.
- **Pool Discovery → Liquidity Update:**
  - All discovered pools are stored in `xrpAmm` (with possible duplicates).
  - Liquidity update operates on all records for each unique pool account.

## 7. Troubleshooting & Recommendations
- **Deduplication:**
  - Ensure pool discovery checks for existing pools before inserting new records.
  - Consider normalizing the `xrpAmm` table: one row per pool account, with historical data in a separate table if needed.
- **Data Quality:**
  - Clean up records with empty/null `account` fields.
  - Add constraints or validation to prevent insertion of incomplete records.
- **Efficiency:**
  - Optimize discovery and update scripts to avoid redundant work.
- **Schema:**
  - Review if all columns are needed in the main table, or if some should be moved to a history/audit table.

---

**Summary:**
- The current process results in many duplicate/historical records per pool account due to lack of deduplication and normalization.
- The liquidity update step updates all records for each pool, not just the latest.
- Improving deduplication and table normalization will reduce row count and improve data quality. 