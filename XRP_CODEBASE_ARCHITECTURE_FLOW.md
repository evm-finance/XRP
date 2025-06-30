# XRP Codebase Architecture & Data Flow Guide

This document provides a comprehensive mapping of data flow from **GraphQL queries → composables → components** for all XRP features in the codebase. Perfect for understanding how everything connects for validation and development.

## 🏗️ **OVERALL ARCHITECTURE PATTERN**

```
📊 GraphQL Query Files (.graphql)
    ↓
🔧 Apollo Queries (queries.ts) 
    ↓
🎯 Composables (useXrp*.ts)
    ↓
🖼️ Vue Components (.vue)
    ↓
📄 Pages (pages/*.vue)
```

---

## 🎯 **KEY VALIDATION AREAS MAPPING**

### **1. XRP BALANCES**

#### **Data Flow Chain:**
```
apollo/main/xrp.query.graphql (XRPAccountBalances)
    ↓
apollo/queries.ts (XRPAccountBalancesGQL)
    ↓
composables/useXrpGraphQLWithLogging.ts (logging wrapper)
    ↓
components/xrp/xrpBalances.vue (main component)
    ↓
components/portfolio/XrpBalanceWidget.vue (widget)
    ↓
pages/xrp-balances.vue (standalone page)
    ↓
pages/xrp-screener.vue (integrated widget)
    ↓
pages/portfolio/balances.vue (portfolio integration)
```

#### **Key Files:**
- **GraphQL Query:** `apollo/main/xrp.query.graphql` → `XRPAccountBalances`
- **Apollo Export:** `apollo/queries.ts` → `XRPAccountBalancesGQL`
- **Main Component:** `components/xrp/xrpBalances.vue`
- **Widget Component:** `components/portfolio/XrpBalanceWidget.vue`
- **Standalone Page:** `pages/xrp-balances.vue`
- **Logging:** `composables/useXrpGraphQLWithLogging.ts`

#### **Query Structure:**
```graphql
query XRPAccountBalances($account: String!) {
  xrpAccountBalances(account: $account) {
    account
    xrpBalance
    xrpPrice
    xrpTokens {
      symbol, issuer, name, balance, price, value
    }
  }
}
```

#### **Data Transformation:**
```typescript
// In components/xrp/xrpBalances.vue
balancesRawData.value = transformedData: XRPBalanceElem[]
// ↓ 
screenerDataFormatted (computed with formatting)
// ↓
v-data-table display with copy functionality
```

---

### **2. XRP TRANSACTIONS**

#### **Data Flow Chain:**
```
apollo/main/xrp.query.graphql (XRPAccountTransactions)
    ↓
apollo/queries.ts (XRPAccountTransactionsGQL)  
    ↓
composables/useXrpGraphQLWithLogging.ts (logging wrapper)
    ↓
components/xrp/xrpAccountHistory.vue (main component)
    ↓
composables/useXrpToken.ts (token page integration)
    ↓
pages/xrp-transactions.vue (standalone page)
    ↓
pages/xrp-token/_id.vue (token page integration)
```

#### **Key Files:**
- **GraphQL Query:** `apollo/main/xrp.query.graphql` → `XRPAccountTransactions`
- **Apollo Export:** `apollo/queries.ts` → `XRPAccountTransactionsGQL`
- **Main Component:** `components/xrp/xrpAccountHistory.vue`
- **Token Integration:** `composables/useXrpToken.ts`
- **Standalone Page:** `pages/xrp-transactions.vue`
- **Explorer Integration:** `pages/xrp-explorer/tx/_id.vue`

#### **Query Structure:**
```graphql
query XRPAccountTransactions($address: String!) {
  xrpAccountTransactions(address: $address) {
    account, amount, destination, transactionType, fee, hash, 
    ledgerIndex, date, validated, flags, sequence, etc.
  }
}
```

#### **Data Transformation:**
```typescript
// In components/xrp/xrpAccountHistory.vue
transactionRawData.value: XRPTransactionElem[]
// ↓
transactionDataFormatted (computed with formatting)
// ↓
v-data-table with action color coding & explorer links
```

---

### **3. XRP AMM HEATMAP**

#### **Data Flow Chain:**
```
apollo/main/heatmap.query.graphql (XRPAMMHeatmap, GetAllAMMLiquidityValues)
    ↓
apollo/queries.ts (XRPAMMHeatmapGQL, XRPAMMLiquidityValuesGQL)
    ↓
composables/useXrpHeatmap.ts (main heatmap logic)
    ↓
composables/useXrpAmmHeatmap.ts (AMM-specific logic)
    ↓
components/xrp/XrpHeatmapChart.vue (chart component)
    ↓
components/xrp/XrpAmmHeatmap.vue (AMM heatmap wrapper)
    ↓
pages/xrp-amm-heatmap.vue (dedicated page)
    ↓
pages/xrp-heatmap.vue (general heatmap page)
    ↓
components/xrp/XrpTerminal.vue (terminal integration)
```

#### **Key Files:**
- **GraphQL Queries:** `apollo/main/heatmap.query.graphql` → `XRPAMMHeatmap`, `GetAllAMMLiquidityValues`
- **Apollo Exports:** `apollo/queries.ts` → `XRPAMMHeatmapGQL`, `XRPAMMLiquidityValuesGQL`
- **Main Composables:** `composables/useXrpHeatmap.ts`, `composables/useXrpAmmHeatmap.ts`
- **Chart Component:** `components/xrp/XrpHeatmapChart.vue`
- **AMM Component:** `components/xrp/XrpAmmHeatmap.vue`
- **Dedicated Pages:** `pages/xrp-amm-heatmap.vue`, `pages/xrp-heatmap.vue`

#### **Query Structure:**
```graphql
query XRPAMMHeatmap($timeFrame: String, $blockSize: String, $limit: Int) {
  xrpAmmHeatmap(timeFrame: $timeFrame, blockSize: $blockSize, limit: $limit) {
    poolId, token1 { currency, symbol, name, issuer }, 
    token2 { currency, symbol, name, issuer },
    liquidity, volume24h, fee, apr, priceChange24h, tvl
  }
}
```

#### **Data Transformation:**
```typescript
// In composables/useXrpHeatmap.ts
heatmapData.value: XrpAmmHeatmapData[]
// ↓
processedData (computed grid layout)
// ↓
XrpHeatmapChart (AmCharts visualization)
// ↓
Data table with pool information
```

---

### **4. XRP TOKEN PAGES**

#### **Data Flow Chain:**
```
apollo/main/token.query.graphql (XRPScreener, XRPAccountTransactions, XRPAccountBalances)
    ↓
apollo/queries.ts (XRPScreenerGQL, XRPAccountTransactionsGQL, XRPAccountBalancesGQL)
    ↓
composables/useXrpToken.ts (main token logic)
    ↓
pages/xrp-token/_id.vue (token detail page)
    ↓
pages/xrp-screener.vue (token listing with links)
    ↓
components/trading/XrpAmmSwap.vue (swap integration)
```

#### **Key Files:**
- **GraphQL Queries:** `apollo/main/token.query.graphql` → `XRPScreener`, transactions, balances
- **Apollo Exports:** `apollo/queries.ts` → `XRPScreenerGQL`, etc.
- **Main Composable:** `composables/useXrpToken.ts`
- **Token Page:** `pages/xrp-token/_id.vue`
- **Screener Page:** `pages/xrp-screener.vue`
- **Swap Integration:** `components/trading/XrpAmmSwap.vue`

#### **Query Structure:**
```graphql
query XRPScreener {
  xrpScreener {
    currency, issuerAddress, icon, tokenName, issuerName,
    marketcap, price, volume24H
  }
}
```

#### **Data Transformation:**
```typescript
// In composables/useXrpToken.ts
tokenData.value: XRPTokenData
tokenBalances.value: array
walletTransactions.value: array
// ↓
pages/xrp-token/_id.vue (formatted display)
// ↓
Copy functionality, USD/XRP toggle, AMM integration
```

---

### **5. XRP AMM LIQUIDITY POOL PAGES**

#### **Data Flow Chain:**
```
apollo/main/xrp.query.graphql (XRPAMMPools, XRPAMMPool, XRPAMMTransactions)
    ↓
apollo/queries.ts (XRPAmmPoolsGQL, XRPAmmPoolDetailsGQL, XRPAmmTransactionsGQL)
    ↓
composables/useXrpAmmLiveData.ts (live pool data)
    ↓
composables/useXrpAmmOperations.ts (pool operations)
    ↓
composables/useXrpAmmTransactions.ts (pool transactions)
    ↓
pages/xrp-amm-pools.vue (pool listing)
    ↓
pages/xrp-amm-pools/_id.vue (individual pool)
    ↓
components/xrp/XrpAmmActionForm.vue (deposit/withdraw forms)
    ↓
components/xrp/XrpAmmSwapDialog.vue (swap dialog)
```

#### **Key Files:**
- **GraphQL Queries:** `apollo/main/xrp.query.graphql` → `XRPAMMPools`, `XRPAMMPool`, `XRPAMMTransactions`
- **Apollo Exports:** `apollo/queries.ts` → `XRPAmmPoolsGQL`, etc.
- **Main Composables:** `composables/useXrpAmmLiveData.ts`, `composables/useXrpAmmOperations.ts`
- **Pool Listing:** `pages/xrp-amm-pools.vue`
- **Pool Details:** `pages/xrp-amm-pools/_id.vue`
- **Action Components:** `components/xrp/XrpAmmActionForm.vue`, `components/xrp/XrpAmmSwapDialog.vue`

#### **Query Structure:**
```graphql
query XRPAMMPools {
  xrpAMMPools {
    poolId, asset1 { currency, issuer }, asset2 { currency, issuer },
    asset1Balance, asset2Balance, lpBalance, fee, tradingVolume24H, tradingVolume7D
  }
}
```

#### **Data Transformation:**
```typescript
// In composables/useXrpAmmLiveData.ts
poolDetails.value: PoolData
// ↓
pages/xrp-amm-pools/_id.vue (formatted display)
// ↓
Deposit/Withdraw/Swap actions with transaction handling
```

---

## 🔐 **WALLET INTEGRATION ARCHITECTURE**

### **Current Implementation (GEM Wallet):**

```
plugins/web3/xrp.client.ts (main XRP client)
    ↓
@gemwallet/api (GEM wallet SDK)
    ↓
components/common/GemWalletConnector.vue (UI component)
    ↓
components/common/WalletSelectDialog.vue (selection dialog)
    ↓
All XRP components (inject wallet state)
```

### **Enhanced Implementation (Multiple Wallets - Disabled):**

```
plugins/web3/enhanced-xrp.client.ts (enhanced client)
    ↓
plugins/web3/xaman.connector.ts (Xaman connector)
    ↓
plugins/web3/metamask-xrp-snap.connector.ts (MetaMask snap)
    ↓
composables/useEnhancedXrpWallet.ts (composable)
    ↓
components/common/EnhancedWalletSelectDialog.vue (UI)
    ↓
components/common/EnhancedXrpWalletConnector.vue (connector)
```

#### **Key Wallet Files:**
- **Main Client:** `plugins/web3/xrp.client.ts`
- **Enhanced Client:** `plugins/web3/enhanced-xrp.client.ts` (disabled)
- **Connectors:** `plugins/web3/xaman.connector.ts`, `plugins/web3/metamask-xrp-snap.connector.ts`
- **Composables:** `composables/useEnhancedXrpWallet.ts`
- **UI Components:** `components/common/GemWalletConnector.vue`, `components/common/WalletSelectDialog.vue`

---

## 🛠️ **CRITICAL COMPOSABLES BREAKDOWN**

### **Core XRP Composables:**

| Composable | Purpose | Key Methods | GraphQL Queries |
|------------|---------|-------------|-----------------|
| `useXrpGraphQLWithLogging.ts` | GraphQL logging wrapper | `useLoggedQuery()` | All XRP queries |
| `useXrpAmm.ts` | AMM swap functionality | `performSwap()`, `calculateSwapQuote()` | AMM-related |
| `useXrpToken.ts` | Token page data | Token info, transactions, balances | `XRPScreener`, `XRPAccountTransactions` |
| `useXrpHeatmap.ts` | Heatmap visualization | `refreshData()`, color calculations | `GetAllAMMLiquidityValues` |
| `useXrpAmmLiveData.ts` | Live AMM data | `getPoolDetails()`, `getUserTokenBalances()` | `XRPAmmPools`, `XRPTokenBalances` |
| `useXrpAmmOperations.ts` | AMM operations | `swap()`, `addLiquidity()`, `removeLiquidity()` | Quote and pool queries |
| `useXrpFormatters.ts` | Data formatting | `formatXrpPrice()`, `formatIssuerAddress()` | None (utility) |

### **Wallet Composables:**

| Composable | Purpose | Key Methods | Dependencies |
|------------|---------|-------------|--------------|
| `useXrpTrade.ts` | Basic trading with GEM | `buy()`, `sell()`, `connectWallet()` | `@gemwallet/api` |
| `useEnhancedXrpWallet.ts` | Multi-wallet support (disabled) | `connectWallet()`, `signTransaction()` | Enhanced client |
| `useXrplTransaction.ts` | Transaction handling | `submitTransaction()`, `waitForConfirmation()` | xrpl.js |

---

## 📊 **GRAPHQL QUERY FILES STRUCTURE**

### **Main Query Files:**
```
apollo/main/
├── xrp.query.graphql          # Core XRP queries (accounts, transactions, AMM)
├── heatmap.query.graphql      # Heatmap-specific queries  
├── token.query.graphql        # Token and screener queries
├── portfolio.query.graphql    # Portfolio-specific queries
└── config.query.graphql       # Configuration queries
```

### **Query Exports (`apollo/queries.ts`):**
```typescript
// Account queries
export const XRPAccountBalancesGQL
export const XRPAccountTransactionsGQL
export const XRPAccountDataGQL

// AMM queries  
export const XRPAmmPoolsGQL
export const XRPAMMHeatmapGQL
export const XRPAMMLiquidityValuesGQL

// Token queries
export const XRPScreenerGQL
export const XRPTokenPriceGQL
export const XRPTokenBalancesGQL

// Heatmap queries
export const XRPHeatmapGQL
export const XRPAMMHeatmapGQL
```

---

## 🚨 **VALIDATION CHECKPOINTS**

### **For Each Component, Validate:**

1. **GraphQL Query Structure** - Does the query match the expected schema?
2. **Data Transformation** - Is GraphQL data properly transformed to component format?
3. **Error Handling** - Are Apollo errors handled gracefully with fallbacks?
4. **Loading States** - Are loading indicators shown during data fetching?
5. **Wallet Integration** - Does the component work with/without wallet connection?
6. **Copy Functionality** - Do copy-to-clipboard features work for addresses/hashes?
7. **Navigation Links** - Do links to explorer/token pages work correctly?
8. **Real-time Updates** - Are polling intervals and refresh mechanisms working?

### **For Wallet Actions, Validate:**

1. **GEM Wallet Detection** - Is wallet installed and accessible?
2. **Connection Flow** - Does the connection process work end-to-end?
3. **Address Retrieval** - Is the wallet address properly fetched and displayed?
4. **Transaction Signing** - Can transactions be signed and submitted?
5. **Error Handling** - Are wallet errors caught and displayed to users?
6. **State Persistence** - Does wallet connection persist across page refreshes?

---

## 🔍 **DEBUGGING GUIDELINES**

### **Common Issues & Solutions:**

1. **GraphQL Errors:**
   - Check `useXrpGraphQLWithLogging.ts` console logs
   - Verify query variables match schema
   - Ensure fallback data is available

2. **Wallet Connection Issues:**
   - Verify GEM wallet is installed
   - Check browser console for SDK errors
   - Ensure cookies are being set correctly

3. **Data Display Problems:**
   - Check data transformation in computed properties
   - Verify formatting functions are working
   - Ensure v-data-table configurations are correct

4. **Missing Features:**
   - Check if enhanced wallet features are disabled
   - Verify composable imports and dependencies
   - Ensure proper plugin registration in `nuxt.config.js`

---

## 📁 **QUICK FILE REFERENCE**

### **Key Directories:**
```
├── apollo/main/              # GraphQL query definitions
├── apollo/queries.ts         # Query exports
├── composables/              # XRP business logic
├── components/xrp/           # XRP-specific components
├── components/trading/       # Trading components (AMM swap)
├── components/portfolio/     # Portfolio components (balance widgets)
├── pages/xrp-*              # XRP pages
├── plugins/web3/            # Wallet integration
└── types/apollo/main/       # TypeScript types
```

### **Critical Files for Validation:**
- `composables/useXrpGraphQLWithLogging.ts` - GraphQL error handling
- `apollo/queries.ts` - All query definitions
- `plugins/web3/xrp.client.ts` - Main wallet client
- `components/xrp/xrpBalances.vue` - Balance display
- `components/xrp/xrpAccountHistory.vue` - Transaction display
- `pages/xrp-amm-pools/_id.vue` - AMM pool details
- `composables/useXrpToken.ts` - Token page logic

### **🎯 Test Wallet Addresses:**

**Primary Test Address (Default for all components):**
```
rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo
```

**Secondary Test Address:**
```
rDodqfAoF8pVh2SoUwhQRfvkqrs4wwxUrz
```

**Usage Pattern:**
```typescript
// All components use consistent default
const accountAddress = computed(() => 
  address.value || 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo'
)
```

---

This guide provides the complete architecture understanding needed for validation and development. Each component follows the same pattern: **GraphQL Query → Composable Logic → Vue Component → Page Integration**. 