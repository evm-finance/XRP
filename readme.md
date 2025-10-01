---
purpose: Main XRP backend - Production GraphQL/REST API server with real-time XRPL integration
category: core-api-server
technology: [go, graphql, websocket, mysql, xrpl]
dependencies: [github.com/gorilla/websocket, gorm.io/gorm, xrpl-library]
endpoints: [/graphql, /api/xrp/*, /ws, /test]
database: [xrpTokens, xrpAmm, xrpPools, xrpLedgers]
integration_priority: high
maintainer: Development Team
last_updated: 2024-12-19
load_frequency: core-backend-development
size_mb: 686.39
file_count: 2840
---

# 🚀 Main XRP Backend Service

## **🔌 XRPL Interface Routines**

Here are the main routines that interface directly with XRPL:

**Price service**
- `GetXRPLPrice`: calls amm_info on XRP/RLUSD to get live XRP price
- `GetTokenPriceInXRPDirect(currency, issuer)`: calls amm_info for XRP/token pairs; handles hex/string currency formats; uses ledger_index "validated"
- `getTokenPriceInXRPFromLedger(poolAccount,...)`: calls account_info and account_lines for pool balances
- `GetTokenPriceFromAMM(currency, issuer)`: uses GetXRPLPrice + GetTokenPriceInXRPDirect to compute USD price

**AMM service**
- `GetAMMPoolInfo(account)`: calls account_info and account_objects (type amm); uses ledger_index "validated"
- `CheckSingleAMMPool(...)`: issues amm_info requests (with validated) for pools; handles errors and tracking

**XRP service (account-level API)**
- `GetAccountBalances(address)`: uses xrpl-go client to call account_info and aggregate balances; also obtains XRP price via GetXRPLPrice

**Connection/streaming**
- `XRPLListener.StartListener/runListener`: uses Subscribe to ledger stream and processes StreamLedger
- `ConnectionManager.GetConnection/createNewConnection/isConnectionHealthy`: manages xrpl-go client connections, performs Ping tests

**Screener**
- `ScreenerUpdateJob.getXRPLPriceWithRetry`: invokes priceService.GetXRPLPrice() once per cycle (will hit ledger unless cached/DB/fallback)
- `AMMPriceCollector.CollectPricesFromAMMPools`: fetches XRP price once via priceService.GetXRPLPrice() and, per batch override, reuses it

## **⚡ QUICK CONTEXT**

**Production-grade Go server providing GraphQL + REST + WebSocket APIs for XRP/XRPL ecosystem**

- **Technology**: Go 1.19+ with Gorilla WebSocket, GORM, XRPL libraries
- **Size**: 686.39 MB (1,910-line main.go + extensive internal packages)
- **Integration**: Direct XRPL network connection + MySQL database
- **Status**: ✅ **PRODUCTION ACTIVE** - Core backend for XRP.FINANCE

---

## 🎯 **CORE FUNCTIONALITY ANALYSIS**

### **Primary Services (from main.go:1-1910)**
```go
// Server Architecture (lines 120-180)
type MinimalServer struct {
    xrpHandler        handler.GraphHandler      // GraphQL query handling
    xrpService        *xrp.XRPService          // XRP blockchain integration
    ammService        xrp.AMMServiceInterface   // AMM pool calculations
    enhancedLedgerSvc *xrp.EnhancedLedgerService // Real-time ledger data
    ledgerMonitor     *xrp.XRPLedgerMonitor    // Live network monitoring
    xrpDB             *gorm.DB                 // MySQL database connection
    wsHub             *WSHub                   // WebSocket connection hub
}
```

### **API Endpoint Analysis**

#### **GraphQL API (`/graphql`)**
```yaml
Primary Operations (lines 200-350):
  XRPAccountBalancesGQL:
    Purpose: Get complete account token balances
    Variables: address (string, required)
    Returns: XRPBalance + XRPTokens array
    Real Resolver: s.xrpHandler.GetXRPAccountBalances()
    
  XRPTopAMMPools:
    Purpose: Retrieve top AMM pools by liquidity
    Variables: limit (default: 20)
    Returns: Pool array with liquidity metrics
    Real Resolver: s.xrpHandler.GetXRPTopAMMPools()
    
  XRPAMMLiquidityValue:
    Purpose: Calculate USD liquidity for specific pool
    Variables: poolId (string, required)
    Returns: Total liquidity in USD
    Real Resolver: s.ammService.GetAMMLiquidityValue()
```

#### **REST API (`/api/xrp/*`)**
```yaml
Core Endpoints (lines 400-1350):
  GET /api/xrp/transaction/{hash}: Single transaction details
  GET /api/xrp/ledger/{id}: Complete ledger with transactions
  GET /api/xrp/blocks: Recent ledgers (paginated)
  GET /api/xrp/tokens: Token registry with metadata
  GET /api/xrp/amm/pools: AMM pool listings
  GET /api/xrp/amm/pools/top: Top pools by liquidity
  GET /api/xrp/token/{id}: Individual token details
  GET /api/xrp/pool/{id}: AMM pool complete data
  POST /api/xrp/account/balances: Account balance queries
  GET /api/xrp/account/{address}/transactions: Transaction history
  GET /api/xrp/monitor/status: Service health monitoring
```

#### **WebSocket API (`/ws`)**
```yaml
Real-time Capabilities (lines 1550-1810):
  Protocol: GraphQL-WS + Custom streaming
  Features:
    - Live ledger updates (s.wsHub.BroadcastLedgerUpdate)
    - GraphQL subscriptions
    - XRP transaction notifications
    - AMM pool change alerts
  Message Types:
    - ledger_update: New ledger notifications
    - graphql: GraphQL subscription responses
    - xrp_stream: Real-time XRPL data
```

---

## 🗄️ **DATABASE SCHEMA & INTEGRATION**

### **MySQL Tables (GORM Integration)**
```sql
-- Core XRP Tables
xrpTokens:
  - token_id (Primary Key)
  - symbol, name, issuer
  - total_supply, market_cap, price_usd
  - gecko_id (CoinGecko integration)
  - created_at, updated_at

xrpAmm:
  - pool_id (Primary Key) 
  - token_a, token_b (Foreign Keys)
  - liquidity_usd, volume_24h
  - trading_fee, lp_token_supply
  - pool_state (JSON)

xrpPools:
  - Real-time pool state tracking
  - Liquidity provider positions
  - Historical performance metrics

xrpLedgers:
  - ledger_index (Primary Key)
  - ledger_hash, tx_count
  - validated_at, events_count
  - raw_data (JSON)
```

### **External Integrations**
```yaml
XRPL Network:
  Protocol: WebSocket (wss://xrpl.ws)
  Purpose: Real-time ledger monitoring
  Data Flow: XRPL → EnhancedLedgerService → Database → WebSocket clients
  
CoinGecko API:
  Purpose: Token price feeds and market data
  Integration: Asynchronous price updates
  
Database Connection:
  Driver: GORM (Go ORM)
  Config: configs.NewConfigRepo("./.env")
  Features: Auto-migration, connection pooling
```

---

## 🔧 **DEVELOPMENT CONTEXT**

### **Key Files & Their Purposes**
```bash
main.go (1,910 lines):
  - HTTP server setup and routing
  - GraphQL request processing
  - REST API endpoint handlers
  - WebSocket connection management
  - Real-time data streaming

internal/xrp/:
  - XRPService: Core blockchain operations
  - AMMService: Liquidity calculations  
  - EnhancedLedgerService: Real-time ledger data
  - XRPLedgerMonitor: Network monitoring

internal/graph/handler/:
  - GraphQL resolvers and query handling
  - Data transformation and validation

configs/:
  - Database configuration
  - Environment variable management
  - XRPL network settings
```

### **Critical Dependencies (from go.mod)**
```go
// Core Framework
github.com/gorilla/websocket v1.5.0  // WebSocket implementation
gorm.io/gorm v1.25.0                // Database ORM
gorm.io/driver/mysql v1.5.0         // MySQL driver

// XRPL Integration  
github.com/rubblelabs/ripple v0.0.0  // XRPL client library

// Utilities
github.com/gorilla/mux              // HTTP routing
encoding/json                       // JSON processing (stdlib)
```

### **Performance Characteristics**
```yaml
Concurrent Connections:
  WebSocket: Unlimited (managed by WSHub)
  HTTP: Go default (excellent under load)
  Database: Connection pooling via GORM

Memory Usage:
  Base: ~50MB Go runtime
  Per Connection: ~1-2KB (WebSocket)
  Cache: Variable based on token/pool data

Response Times:
  GraphQL: 50-200ms (database dependent)
  REST: 10-100ms (cached vs real-time)
  WebSocket: <10ms (real-time updates)
```

---

## 🚀 **QUICK START GUIDE**

### **1. Environment Setup**
```bash
cd backends/main/

# Install Go dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit: DATABASE_URL, XRPL_WEBSOCKET_URL, etc.
```

### **2. Database Initialization**
```bash
# Run migrations (if needed)
go run cmd/db_util.go -migrate

# Verify connection
go run cmd/db_util.go -test-connection
```

### **3. Start Development Server**
```bash
# Run main server
go run main.go

# Server starts on localhost:8080
# GraphQL Playground: http://localhost:8080/test
# Health Check: http://localhost:8080/api/xrp/monitor/status
```

### **4. Test Integration**
```bash
# Test GraphQL
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { xrpTopAMMPools { poolId liquidityUsd } }"}'

# Test REST
curl http://localhost:8080/api/xrp/blocks?limit=5

# Test WebSocket (browser console)
const ws = new WebSocket('ws://localhost:8080/ws')
ws.onmessage = (event) => console.log(JSON.parse(event.data))
```

---

## 🔗 **SERVICE DEPENDENCIES & INTEGRATIONS**

### **Internal Dependencies**
```yaml
Frontend Integration:
  - Current: Limited REST API usage
  - Available: Full GraphQL + WebSocket capabilities
  - Optimization: Frontend can leverage real-time features

Database Dependencies:
  - MySQL 8.0+ (required)
  - GORM migrations (auto-handled)
  - Connection pooling (configured)

Configuration Dependencies:
  - .env file (database, XRPL settings)
  - configs/ package (centralized config)
```

### **External Service Dependencies**
```yaml
XRPL Network:
  Status: CRITICAL - Core functionality
  Endpoints: wss://xrpl.ws, wss://s2.ripple.com
  Failover: Multiple rippled servers configured
  
CoinGecko API:
  Status: OPTIONAL - Price data enhancement
  Rate Limits: 50 calls/minute (free tier)
  Fallback: Cached price data available
  
CDN/Static Assets:
  Status: OPTIONAL - Frontend served separately
  Current: No static file serving in this backend
```

---

## 📊 **MONITORING & HEALTH CHECKS**

### **Built-in Monitoring Endpoints**
```yaml
GET /api/xrp/monitor/status:
  Purpose: Service health check
  Returns: Database connectivity, XRPL status, WebSocket connections
  
GET /api/xrp/monitor/stats:
  Purpose: Performance metrics
  Returns: Request counts, response times, error rates
  
WebSocket Health:
  Endpoint: /ws (connection test)
  Metrics: Connected client count, message throughput
```

### **Logging & Debugging**
```go
// Structured logging throughout main.go
log.Printf("📡 Broadcasting ledger %d to %d connected clients", ledger.LedgerIndex, clientCount)
log.Printf("✅ GraphQL success with real data")
log.Printf("❌ Error from real resolver: %v", err)

// Debug endpoints
GET /test: Interactive GraphQL playground
GET /api/xrp/monitor/debug: Detailed system state
```

---

## 🔧 **CONTEXT OPTIMIZATION FEATURES**

### **Meta Tag Implementation Status**
- [ ] **Add meta tags to main.go** (lines 1-50)
- [ ] **Add meta tags to internal/ packages**
- [ ] **Document API endpoint purposes**
- [ ] **Create dependency graph**

### **Documentation Optimization**
- [x] **Complete README** (this file)
- [x] **API endpoint catalog** 
- [x] **Database schema documentation**
- [x] **Development setup guide**

### **Import Graph Dependencies**
```yaml
External Imports:
  - github.com/gorilla/websocket (WebSocket server)
  - gorm.io/gorm (Database ORM)
  - XRPL libraries (Blockchain integration)
  
Internal Imports:
  - configs/ (Configuration management)
  - internal/xrp/ (Core XRP functionality)
  - internal/graph/handler/ (GraphQL handling)
  
Service Relationships:
  Provides: GraphQL API, REST API, WebSocket streaming
  Consumes: MySQL database, XRPL network
  Integrates: Frontend applications, monitoring systems
```

---

## 📚 **RELATED DOCUMENTATION**

### **Backend Infrastructure**
- [Backend Infrastructure Overview](../README.md)
- [Backend Context Optimization Strategy](../../docs/current/BACKEND_CONTEXT_OPTIMIZATION_STRATEGY.md)
- [Service Catalog Summary](../../docs/current/SERVICE_CATALOG_SUMMARY.md)

### **API Documentation**
- [Backend API Complete](../../docs/current/BACKEND_API_COMPLETE.md)
- [API Quick Reference](../../docs/current/API_QUICK_REFERENCE.md)
- [Frontend-Backend Integration Strategy](../../docs/current/FRONTEND_BACKEND_INTEGRATION_STRATEGY.md)

### **Development Resources**
- [Quick Debug Guide](../../docs/current/QUICK_DEBUG_GUIDE.md)
- [Weekly Maintenance Checklist](../../docs/current/WEEKLY_MAINTENANCE_CHECKLIST.md)
- Sample AMM Queries: `sample_amm_queries.md`
- AMM Enhancement Roadmap: `AMM_ENHANCEMENT_ROADMAP.md`

---

**Status**: 🟢 **PRODUCTION ACTIVE** - Core backend serving XRP.FINANCE frontend  
**Optimization Level**: 🔧 **READY FOR CONTEXT OPTIMIZATION** - Meta tags and import graphs pending  
**Integration Readiness**: 🚀 **ADVANCED FEATURES AVAILABLE** - GraphQL + WebSocket ready for frontend enhancement
