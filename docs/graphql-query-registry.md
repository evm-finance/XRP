# GraphQL Query Registry

## Overview
This registry provides a comprehensive mapping of all GraphQL queries in the XRP project, including their parameters, return types, usage locations, and dependencies.

## XRP Account Queries

### XRPAccountBalancesGQL
**File**: `apollo/queries.ts:58-75`
**Parameters**: `$account: String!`
**Return Type**: `XRPAccountBalances`
**Usage Locations**:
- `pages/xrp-balances.vue` - Main balances page
- `pages/xrp-portfolio.vue` - Portfolio page
- `components/xrp/xrpBalances.vue` - Balance display component
- `components/portfolio/XrpBalanceWidget.vue` - Portfolio widget
- `composables/useXrpGraphQL.ts` - Central query management
- `composables/useXrpToken.ts` - Token management

**Error Handling**: Basic error states in most locations
**Loading States**: Available in some locations
**Refetch**: Available in some locations

### XRPAccountTransactionsGQL
**File**: `apollo/queries.ts:45-57`
**Parameters**: `$address: String!`
**Return Type**: `XRPTransactions`
**Usage Locations**:
- `pages/xrp-transactions.vue` - Main transactions page
- `components/xrp/xrpAccountHistory.vue` - Account history component
- `composables/useXrpGraphQL.ts` - Central query management
- `composables/useXrpToken.ts` - Token management
- `composables/useXrpScrerener.ts` - Screener functionality

**Error Handling**: Basic error states in most locations
**Loading States**: Available in some locations
**Refetch**: Available in some locations

### XRPTransactionGQL
**File**: `apollo/queries.ts:18-44`
**Parameters**: `$hash: String!`
**Return Type**: `XRPTransaction`
**Usage Locations**:
- `pages/xrp-explorer/tx/_id.vue` - Transaction detail page

**Error Handling**: Basic error handling
**Loading States**: No explicit loading state
**Refetch**: Not available

## XRP AMM Queries

### XRPAmmPoolsGQL
**File**: `apollo/queries.ts:77-100`
**Parameters**: None
**Return Type**: `XRPAmmPools`
**Usage Locations**:
- `pages/xrp-amm-pools.vue` - AMM pools page
- `composables/useXrpGraphQL.ts` - Central query management

**Error Handling**: onError callback in pools page
**Loading States**: No explicit loading state
**Refetch**: Not available

### XRPAmmPoolDetailsGQL
**File**: `apollo/queries.ts:102-140`
**Parameters**: `$poolId: String!`
**Return Type**: `XRPAmmPoolDetails`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### XRPAmmUserPositionsGQL
**File**: `apollo/queries.ts:142-170`
**Parameters**: `$address: String!`
**Return Type**: `XRPAmmUserPositions`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### XRPAmmQuoteGQL
**File**: `apollo/queries.ts:172-185`
**Parameters**: `$poolId: String!, $amount: String!, $fromToken: String!`
**Return Type**: `XRPAmmQuote`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### XRPAmmTransactionsGQL
**File**: `apollo/queries.ts:187-202`
**Parameters**: `$poolId: String!, $limit: Int`
**Return Type**: `XRPAmmTransactions`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

## XRP Token Queries

### XRPTokenPriceGQL
**File**: `apollo/queries.ts:205-220`
**Parameters**: `$currency: String!, $issuer: String`
**Return Type**: `XRPTokenPrice`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### XRPTokenBalancesGQL
**File**: `apollo/queries.ts:222-237`
**Parameters**: `$address: String!`
**Return Type**: `XRPTokenBalances`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

## XRP Screener Queries

### XRPScreenerGQL
**File**: `apollo/queries.ts:65-76`
**Parameters**: None
**Return Type**: `XRPScreener`
**Usage Locations**:
- `pages/xrp-screener.vue` - Main screener page
- `composables/useXrpGraphQL.ts` - Central query management
- `composables/useXrpToken.ts` - Token management
- `composables/useXrpTokenHeatmap.ts` - Token heatmap

**Error Handling**: Minimal error handling
**Loading States**: No explicit loading state
**Refetch**: Not available
**Polling**: 60-second interval in screener page

## XRP Heatmap Queries

### XRPHeatmapGQL
**File**: `apollo/queries.ts:240-265`
**Parameters**: `$timeFrame: String, $blockSize: String, $limit: Int`
**Return Type**: `XRPHeatmap`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### XRPAMMHeatmapGQL
**File**: `apollo/queries.ts:267-295`
**Parameters**: `$timeFrame: String, $blockSize: String, $limit: Int`
**Return Type**: `XRPAMMHeatmap`
**Usage Locations**:
- `composables/useXrpHeatmap.ts` - Heatmap composable

**Error Handling**: Basic error handling
**Loading States**: Available
**Refetch**: Available

### XRPHeatmapUpdatesGQL
**File**: `apollo/queries.ts:297-315`
**Parameters**: `$currencies: [String!], $timeFrame: String`
**Return Type**: `XRPHeatmapUpdates`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

## XRP Block Queries

### BlocksXrpGQL
**File**: `apollo/queries.ts:4-17`
**Parameters**: `$network: String!`
**Return Type**: `Blocks`
**Usage Locations**:
- `composables/useXrpScrerener.ts` - Screener composable

**Error Handling**: No explicit error handling
**Loading States**: No explicit loading state
**Refetch**: Not available

## XRP AMM Liquidity Queries

### XRPAMMLiquidityValuesGQL
**File**: `apollo/queries.ts:318-330`
**Parameters**: None
**Return Type**: `XRPAMMLiquidityValues`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### SimpleAMMLiquidityValuesGQL
**File**: `apollo/queries.ts:332-340`
**Parameters**: None
**Return Type**: `SimpleAMMLiquidityValues`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

### GetAllAMMLiquidityValuesGQL
**File**: `apollo/queries.ts:342-365`
**Parameters**: None
**Return Type**: `GetAllAMMLiquidityValues`
**Usage Locations**: None identified yet

**Error Handling**: Not implemented
**Loading States**: Not implemented
**Refetch**: Not available

## Test Queries

### TestGraphQLGQL
**File**: `apollo/queries.ts:367-375`
**Parameters**: None
**Return Type**: `__Schema`
**Usage Locations**:
- `composables/useXrpHeatmap.ts` - Heatmap composable

**Error Handling**: Basic error handling
**Loading States**: Available
**Refetch**: Not available

## Query Dependencies and Relationships

### Account Data Dependencies
```
XRPAccountBalancesGQL
├── Used by: XRP balance displays
├── Dependencies: None
└── Related: XRPAccountTransactionsGQL

XRPAccountTransactionsGQL
├── Used by: Transaction history displays
├── Dependencies: None
└── Related: XRPTransactionGQL (for details)
```

### AMM Data Dependencies
```
XRPAmmPoolsGQL
├── Used by: Pool listings
├── Dependencies: None
└── Related: XRPAmmPoolDetailsGQL, XRPAmmUserPositionsGQL

XRPAmmPoolDetailsGQL
├── Used by: Pool detail pages
├── Dependencies: XRPAmmPoolsGQL
└── Related: XRPAmmTransactionsGQL, XRPAmmQuoteGQL
```

### Token Data Dependencies
```
XRPScreenerGQL
├── Used by: Token listings and heatmaps
├── Dependencies: None
└── Related: XRPTokenPriceGQL, XRPTokenBalancesGQL

XRPHeatmapGQL
├── Used by: Token heatmap visualizations
├── Dependencies: XRPScreenerGQL
└── Related: XRPHeatmapUpdatesGQL
```

## Critical Query Analysis

### High Priority Queries (Most Used)
1. **XRPAccountBalancesGQL** - 6 usage locations
2. **XRPAccountTransactionsGQL** - 5 usage locations
3. **XRPScreenerGQL** - 4 usage locations
4. **XRPAmmPoolsGQL** - 2 usage locations

### Medium Priority Queries (Some Usage)
1. **XRPTransactionGQL** - 1 usage location
2. **XRPAMMHeatmapGQL** - 1 usage location
3. **BlocksXrpGQL** - 1 usage location
4. **TestGraphQLGQL** - 1 usage location

### Low Priority Queries (No Usage Yet)
1. **XRPAmmPoolDetailsGQL**
2. **XRPAmmUserPositionsGQL**
3. **XRPAmmQuoteGQL**
4. **XRPAmmTransactionsGQL**
5. **XRPTokenPriceGQL**
6. **XRPTokenBalancesGQL**
7. **XRPHeatmapGQL**
8. **XRPHeatmapUpdatesGQL**
9. **XRPAMMLiquidityValuesGQL**
10. **SimpleAMMLiquidityValuesGQL**
11. **GetAllAMMLiquidityValuesGQL**

## Error Handling Status

### Good Error Handling
- `pages/xrp-transactions.vue` - Basic error states
- `pages/xrp-balances.vue` - Basic error states
- `pages/xrp-amm-pools.vue` - onError callbacks
- `composables/useXrpGraphQL.ts` - Basic error states

### Poor Error Handling
- `pages/xrp-screener.vue` - Minimal error handling
- `pages/xrp-explorer/tx/_id.vue` - Basic error handling
- `pages/xrp-explorer/ledger/_id.vue` - No error handling
- Most composables - No explicit error handling

### No Error Handling
- All unused queries
- Some components with basic error handling

## Recommendations

### Immediate Actions (Step 2 Priority)
1. **Add error handling to high-priority queries first**
   - XRPAccountBalancesGQL (6 locations)
   - XRPAccountTransactionsGQL (5 locations)
   - XRPScreenerGQL (4 locations)

2. **Implement logging for critical pages**
   - xrp-screener.vue (high traffic)
   - xrp-transactions.vue (critical functionality)
   - xrp-balances.vue (critical functionality)

3. **Standardize error handling patterns**
   - Use consistent error handling across all query usages
   - Implement retry logic for failed queries
   - Add fallback data for offline scenarios

### Future Actions
1. **Remove or implement unused queries**
2. **Consolidate duplicate query definitions**
3. **Add comprehensive loading states**
4. **Implement query caching strategies**

## Next Steps
- Complete Step 1.4: Document current error handling patterns
- Begin Step 2: Implement logging for XRP-specific queries
- Focus on high-priority queries and critical pages 