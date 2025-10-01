# REST API Endpoints - Updated for Normalized Data

## Overview
All XRP AMM endpoints have been updated to use the **normalized database** with **deduplicated, high-quality data**. The migration removed 255 duplicate records and improved data integrity significantly.

## Endpoint Summary

| Endpoint | Method | Description | Status |
|----------|--------|-------------|--------|
| `/api/xrp/amm-pools` | GET | Top AMM pools (updated) | ✅ Updated |
| `/api/xrp/amm-normalized` | GET | Advanced AMM pools with filtering | 🆕 New |
| `/api/xrp/amm-stats` | GET | Migration and deduplication statistics | 🆕 New |
| `/api/admin/data-quality` | GET | Data quality validation metrics | 🆕 New |
| `/api/xrp/account-balances` | POST | XRP account balances | ✅ Existing |
| `/api/xrp/screener` | GET | Token screener | ✅ Existing |
| `/health` | GET | Server health check | ✅ Existing |

---

## Updated Endpoints

### GET `/api/xrp/amm-pools`
**Updated to use normalized table with improved performance**

**Parameters:**
- `limit` (optional): Number of pools to return (max 100, default 20)

**Response:**
```json
{
  "pools": [
    {
      "account": "rPool123...",
      "amount": "1000000",
      "amount2currency": "USD",
      "amount2issuer": "rIssuer...",
      "amount2value": "2500.50",
      "liquidity_usd": 5000.50,
      "asset1_value_usd": 2500.25,
      "asset2_value_usd": 2500.25,
      "last_liquidity_update": "2025-07-29T14:20:00Z",
      "created_at": "2025-07-20T10:00:00Z",
      "last_updated": "2025-07-29T14:20:00Z"
    }
  ],
  "count": 20,
  "source": "normalized",
  "deduplication_applied": true
}
```

**Changes:**
- ✅ **No more duplicates** - uses `xrpAmm_normalized` table
- ✅ **Better performance** - simplified query without deduplication logic
- ✅ **Additional fields** - `asset1_value_usd`, `asset2_value_usd`, `last_liquidity_update`
- ✅ **Quality indicators** - response includes source and deduplication status

---

## New Endpoints

### GET `/api/xrp/amm-normalized`
**Advanced AMM pools endpoint with filtering capabilities**

**Parameters:**
- `limit` (optional): Number of pools (max 200, default 50)
- `min_liquidity` (optional): Minimum liquidity in USD (default 0)

**Example Request:**
```
GET /api/xrp/amm-normalized?limit=10&min_liquidity=1000
```

**Response:**
```json
{
  "success": true,
  "data": {
    "pools": [
      {
        "account": "rPool123...",
        "amount": "1000000",
        "amount2currency": "USD",
        "amount2issuer": "rIssuer...",
        "amount2value": "2500.50",
        "asset2frozen": 0,
        "lptokencurrency": "LP",
        "lptokenissuer": "rLP...",
        "lptokenvalue": "1000.0",
        "tradingfee": 30,
        "liquidity_usd": 5000.50,
        "asset1_value_usd": 2500.25,
        "asset2_value_usd": 2500.25,
        "last_liquidity_update": "2025-07-29T14:20:00Z",
        "created_at": "2025-07-20T10:00:00Z",
        "last_updated": "2025-07-29T14:20:00Z"
      }
    ],
    "count": 10,
    "filters": {
      "limit": 10,
      "min_liquidity": 1000
    },
    "source": "xrpAmm_normalized",
    "note": "Deduplicated data with improved quality"
  }
}
```

**Use Cases:**
- High-liquidity pool discovery
- Advanced filtering for trading interfaces
- Data analysis and reporting
- Quality-focused pool selection

### GET `/api/xrp/amm-stats`
**Migration statistics and deduplication metrics**

**Response:**
```json
{
  "success": true,
  "data": {
    "migration_status": "completed",
    "original_records": 1423,
    "normalized_records": 1168,
    "duplicates_removed": 255,
    "deduplication_ratio": "1.2:1",
    "storage_space_saved": "17.9%",
    "data_quality": {
      "null_accounts": 0,
      "duplicate_accounts": 0,
      "data_integrity": "excellent"
    }
  }
}
```

**Use Cases:**
- Migration validation
- Data quality reporting
- Performance metrics
- Administrative monitoring

### GET `/api/admin/data-quality`
**Data validation and quality metrics**

**Response:**
```json
{
  "success": true,
  "data": {
    "table": "xrpAmm_normalized",
    "validation_results": {
      "invalid_account_formats": 0,
      "unreasonable_liquidity_values": 0,
      "valid_pools": 1168
    },
    "quality_score": "excellent",
    "recommendations": [
      "Data quality is excellent after normalization",
      "All pools have valid XRP account formats",
      "Liquidity values are within reasonable ranges"
    ]
  }
}
```

**Use Cases:**
- Data quality monitoring
- Validation reporting
- Administrative oversight
- Health checks

---

## Migration Benefits

### Performance Improvements
- ✅ **17.9% storage reduction** (255 duplicates removed)
- ✅ **Faster queries** - no deduplication logic needed
- ✅ **Better indexing** - primary keys and unique constraints
- ✅ **Simplified code** - cleaner, more maintainable queries

### Data Quality Improvements
- ✅ **Zero null accounts** - perfect data integrity
- ✅ **Zero duplicates** - primary key enforcement
- ✅ **Valid XRP addresses** - all accounts properly formatted
- ✅ **Reasonable liquidity values** - data validation applied
- ✅ **Consistent timestamps** - proper audit trail

### Frontend Benefits
- ✅ **Reliable data** - no duplicate pool entries
- ✅ **Better UX** - consistent pool information
- ✅ **Faster loading** - improved query performance
- ✅ **Advanced filtering** - new filtering capabilities
- ✅ **Quality metrics** - transparency about data quality

---

## Error Handling

### Standard Error Response
```json
{
  "error": "Error description"
}
```

### Common Error Codes
- `400` - Bad Request (invalid parameters)
- `500` - Internal Server Error (database/server issues)
- `200` - Success

---

## Testing

### Test Script
Run the endpoint testing script:
```bash
go run scripts/test_rest_endpoints.go
```

### Manual Testing
```bash
# Test updated AMM pools endpoint
curl http://localhost:8080/api/xrp/amm-pools?limit=5

# Test new normalized endpoint with filtering
curl "http://localhost:8080/api/xrp/amm-normalized?limit=10&min_liquidity=1000"

# Test migration statistics
curl http://localhost:8080/api/xrp/amm-stats

# Test data quality
curl http://localhost:8080/api/admin/data-quality

# Health check
curl http://localhost:8080/health
```

---

## Frontend Integration

### Key Changes for Frontend
1. **Same endpoints work better** - existing `/api/xrp/amm-pools` now returns cleaner data
2. **New filtering options** - use `/api/xrp/amm-normalized` for advanced filtering
3. **Quality transparency** - responses include data source and quality indicators
4. **Additional fields** - access to `asset1_value_usd`, `asset2_value_usd`, `last_liquidity_update`

### Recommended Usage
- **General pool listing**: Use `/api/xrp/amm-pools` (existing functionality, improved data)
- **Advanced filtering**: Use `/api/xrp/amm-normalized` for min liquidity thresholds
- **Analytics dashboards**: Use `/api/xrp/amm-stats` for migration metrics
- **Admin panels**: Use `/api/admin/data-quality` for monitoring

---

## Database Schema

### Normalized Table: `xrpAmm_normalized`
```sql
CREATE TABLE xrpAmm_normalized (
    account VARCHAR(128) PRIMARY KEY,
    amount VARCHAR(128) NOT NULL,
    amount2currency VARCHAR(128),
    amount2issuer VARCHAR(128),
    amount2value VARCHAR(128),
    asset2frozen INT DEFAULT 0,
    lptokencurrency VARCHAR(128),
    lptokenissuer VARCHAR(128),
    lptokenvalue VARCHAR(128),
    tradingfee INT DEFAULT 0,
    liquidity_usd DECIMAL(20,8) DEFAULT 0,
    asset1_value_usd DECIMAL(20,8) DEFAULT 0,
    asset2_value_usd DECIMAL(20,8) DEFAULT 0,
    last_liquidity_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_asset_pair (amount2currency, amount2issuer),
    INDEX idx_liquidity (liquidity_usd),
    INDEX idx_currency (amount2currency),
    INDEX idx_last_updated (last_updated)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**Key Improvements:**
- ✅ **Primary key** on `account` (prevents duplicates)
- ✅ **Unique constraint** on asset pairs
- ✅ **Optimized indexes** for common queries
- ✅ **Proper data types** for USD values
- ✅ **Timestamp tracking** for audit trail 