# XRP Tokens Table Migration Analysis

## 🚨 CRITICAL ISSUE IDENTIFIED
The current `xrpTokens` table is **obsolete** with poor data quality:
- Only 5% of tokens have names
- Inconsistent price data ($0 price but $5M market cap)
- Missing proper currency/issuer unique key constraints

The `xrpScreener` table has **excellent data quality**:
- 100% price coverage with live calculations
- Fresh data (all updated today)
- Proper currency/issuer unique key structure
- Active price service integration

## 📋 MIGRATION PLAN

### Step 1: Drop Deprecated Table
**CRITICAL**: Must drop `xrpTokens` table FIRST before creating new one
```sql
DROP TABLE IF EXISTS xrpTokens;
```

### Step 2: Create New xrpTokens Table
Use `xrpScreener` structure as template:
```sql
CREATE TABLE xrpTokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    currency VARCHAR(40) NOT NULL COMMENT 'Token currency code',
    issuer VARCHAR(34) NOT NULL COMMENT 'Token issuer address', 
    token_name VARCHAR(255) DEFAULT NULL COMMENT 'Human readable token name',
    icon VARCHAR(500) DEFAULT NULL COMMENT 'Token icon URL',
    price DECIMAL(20,8) DEFAULT 0 COMMENT 'Current token price in USD',
    marketcap DECIMAL(20,2) DEFAULT 0 COMMENT 'Market capitalization',
    supply_xrpl DECIMAL(20,8) DEFAULT NULL COMMENT 'Total token supply',
    liquidity_total DECIMAL(20,2) DEFAULT 0 COMMENT 'Total liquidity in USD',
    trustlines INT DEFAULT 0 COMMENT 'Number of trustlines',
    volume_24h DECIMAL(20,2) DEFAULT 0 COMMENT 'Trading volume 24h',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_token (currency, issuer),
    INDEX idx_liquidity (liquidity_total),
    INDEX idx_marketcap (marketcap),
    INDEX idx_price (price),
    INDEX idx_updated (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Step 3: Migrate Data
Copy data from `xrpScreener` to new `xrpTokens`:
```sql
INSERT INTO xrpTokens (
    currency, issuer, token_name, icon, price, marketcap, 
    supply_xrpl, liquidity_total, trustlines, volume_24h, updated_at
)
SELECT 
    currency, issuer, token_name, icon, price, marketcap,
    supply_xrpl, liquidity_total, trustlines, volume_24h, updated_at
FROM xrpScreener;
```

## 📁 FILES REQUIRING UPDATES

### Core Service Files
1. **`internal/xrp/xrp_service.go`** (Lines 123-129)
   - `GetTokenList()` method - UPDATE query to use new schema
   - Query currently uses: `SELECT currency, issuer, CONVERT(name USING utf8) as name, icon, price, marketcap, supply_xrpl, liquidity_total FROM xrpTokens`
   - **Action**: Update to use new column names and structure

2. **`internal/xrp/token_service.go`** (Multiple locations)
   - `StoreTokens()` method - UPDATE INSERT/UPDATE statements
   - `GetTokensFromDatabase()` method - UPDATE query structure
   - **Action**: Update all SQL statements to use new schema

3. **`internal/xrp/amm_price_collector.go`** (Lines 439-443)
   - `StorePricesToDatabase()` method - UPDATE to new table structure
   - Currently: `UPDATE xrpTokens SET price = ?, marketcap = ?, volume24h = ?, updated_at = NOW() WHERE currency = ? AND issuer = ?`
   - **Action**: Update to use new column names

### API Handler Files
4. **`main.go`** (Lines 796-804, 628-632)
   - `getToken()` handler - UPDATE XRP token query
   - `getTokenMints()` handler - UPDATE query structure
   - **Action**: Update all xrpTokens table references

5. **`internal/graph/handler/xrp.go`** 
   - GraphQL resolvers - UPDATE any xrpTokens references
   - **Action**: Check and update GraphQL token queries

### Database Schema Files
6. **`internal/database/database_create.go`** (Lines 122-148)
   - `CreateXRPTokensTable()` method - REPLACE with new schema
   - **Action**: Replace entire method with new table structure

### Test Files
7. **`internal/xrp/testing/token_handlers_test.go`**
   - **Action**: Update tests to validate new schema

8. **All other test files** referencing xrpTokens
   - **Action**: Update queries to use new structure

### Migration Scripts
9. **`scripts/` directory** (Multiple files)
   - Various scripts populating or querying xrpTokens
   - **Action**: Update to use new schema

## 🔧 SCHEMA CHANGES REQUIRED

### Column Mapping: Old → New
```
Old xrpTokens          → New xrpTokens
=================      ================
issuer                 → issuer (same)
currency               → currency (same)  
name                   → token_name
icon                   → icon (same)
price                  → price (same)
marketcap              → marketcap (same)
supply                 → supply_xrpl
volume24h              → volume_24h
trustlines             → trustlines (same)
(missing)              → liquidity_total (NEW)
(missing)              → updated_at (NEW)
(missing)              → UNIQUE KEY (currency, issuer)
```

### Critical Differences
1. **Unique Key**: New table has `UNIQUE KEY (currency, issuer)` constraint
2. **Column Names**: `name` → `token_name`, `supply` → `supply_xrpl`  
3. **New Columns**: `liquidity_total`, proper `updated_at` tracking
4. **Data Types**: Better precision for decimal fields
5. **Indexes**: Proper indexing for performance

## ⚠️ MIGRATION RISKS

1. **Data Loss Risk**: Must drop old table BEFORE creating new one
2. **Downtime Risk**: API endpoints will fail during migration
3. **Reference Risk**: All code references must be updated simultaneously
4. **Constraint Risk**: New unique key constraint may reject duplicate data

## ✅ VALIDATION CHECKLIST

Before migration:
- [ ] All file references documented
- [ ] New schema validated against xrpScreener
- [ ] Migration script tested on staging
- [ ] All code updates prepared
- [ ] Rollback plan prepared

After migration:
- [ ] New table created successfully
- [ ] Data migrated from xrpScreener
- [ ] All API endpoints working
- [ ] Tests passing with new schema
- [ ] Performance validated

## 🎯 EXPECTED BENEFITS

1. **Data Quality**: 100% price coverage (vs 40% current)
2. **Consistency**: Single source of truth for token data
3. **Performance**: Better indexing and structure
4. **Maintainability**: Currency/issuer unique keys
5. **Integration**: Seamless with existing price service
