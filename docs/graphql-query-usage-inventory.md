# GraphQL Query Usage Inventory

## Overview
This document maps all GraphQL query usage across the codebase, identifying which components, pages, and composables use specific queries.

## Pages Using GraphQL Queries

### 1. `pages/xrp-transactions.vue`
**Status**: Active XRP transaction page
**Queries Used**:
- `XRPAccountTransactionsGQL` - via `useQuery`
- **Error Handling**: Basic error state with `queryError`
- **Loading State**: `queryLoading` state
- **Refetch**: `refetch` function available

### 2. `pages/xrp-screener.vue`
**Status**: Active XRP token screener page
**Queries Used**:
- `XRPScreenerGQL` - via `useQuery`
- **Error Handling**: Minimal error handling
- **Loading State**: No explicit loading state
- **Polling**: 60-second poll interval

### 3. `pages/xrp-explorer/tx/_id.vue`
**Status**: XRP transaction detail page
**Queries Used**:
- `XRPTransactionGQL` - via `useQuery`
- **Error Handling**: Basic error handling
- **Loading State**: No explicit loading state

### 4. `pages/xrp-explorer/ledger/_id.vue`
**Status**: XRP ledger detail page
**Queries Used**:
- `BlockGQL` - via `useQuery`
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

### 5. `pages/xrp-balances.vue`
**Status**: XRP balances page
**Queries Used**:
- `XRPAccountBalancesGQL` - via `useQuery`
- **Error Handling**: Basic error state with `queryError`
- **Loading State**: `queryLoading` state
- **Refetch**: `refetch` function available

### 6. `pages/xrp-amm-pools.vue`
**Status**: XRP AMM pools page
**Queries Used**:
- `XRPAmmPoolsGQL` - via `useQuery`
- **Error Handling**: `onError` callback
- **Loading State**: No explicit loading state

## Components Using GraphQL Queries

### 1. `components/xrp/xrpBalances.vue`
**Status**: XRP balance display component
**Queries Used**:
- `XRPAccountBalancesGQL` - via `useQuery`
- **Error Handling**: `onError` callback
- **Loading State**: No explicit loading state

### 2. `components/xrp/xrpAccountHistory.vue`
**Status**: XRP account history component
**Queries Used**:
- `XRPAccountTransactionsGQL` - via `useQuery`
- **Error Handling**: `onError` callback
- **Loading State**: No explicit loading state

### 3. `components/portfolio/XrpBalanceWidget.vue`
**Status**: Portfolio balance widget
**Queries Used**:
- `XRPAccountBalancesGQL` - via `useQuery`
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

### 4. `components/common/GasInfo.vue`
**Status**: Gas information component
**Queries Used**:
- `GasGQL` - via `useQuery` (non-XRP query)
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

## Composables Using GraphQL Queries

### 1. `composables/useXrpAccounts.ts`
**Status**: XRP account management composable
**Queries Used**:
- `XRPDefiDataGQL` - via `useQuery`
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

### 2. `composables/useXrpAmmOperations.ts`
**Status**: XRP AMM operations composable
**Queries Used**:
- Uses `useMutation` (not queries)
- **Error Handling**: No explicit error handling

### 3. `composables/useXrpGraphQL.ts`
**Status**: Main XRP GraphQL composable
**Queries Used**:
- `XRPAccountBalancesGQL` - via `useQuery`
- `XRPAccountTransactionsGQL` - via `useQuery`
- `XRPScreenerGQL` - via `useQuery`
- `XRPAmmPoolsGQL` - via `useQuery`
- **Error Handling**: Basic error states
- **Loading State**: Loading states available
- **Refetch**: Refetch functions available

### 4. `composables/useXrpScrerener.ts`
**Status**: XRP screener composable
**Queries Used**:
- `BlocksXrpGQL` - via `useQuery` and `useSubscription`
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

### 5. `composables/useXrpTokenHeatmap.ts`
**Status**: XRP token heatmap composable
**Queries Used**:
- `XRPScreenerGQL` - via `useQuery`
- **Error Handling**: No explicit error handling
- **Loading State**: No explicit loading state

### 6. `composables/useXrpToken.ts`
**Status**: XRP token management composable
**Queries Used**:
- `XRPScreenerGQL` - via `useQuery`
- `XRPAccountTransactionsGQL` - via `useQuery`
- `XRPAccountBalancesGQL` - via `useQuery`
- **Error Handling**: `onError` callbacks for all queries
- **Loading State**: No explicit loading state

### 7. `composables/useXrpTokenMints.ts`
**Status**: XRP token mints composable
**Queries Used**:
- `XRPTokenMintsGQL` - via `useQuery`
- `XRPLiquidityPoolsGQL` - via `useQuery`
- **Error Handling**: `onError` callbacks for all queries
- **Loading State**: No explicit loading state
- **Refetch**: Refetch functions available

### 8. `composables/useXrpHeatmap.ts`
**Status**: XRP heatmap composable
**Queries Used**:
- `TestGraphQLGQL` - via `useQuery`
- `XRPAMMHeatmapGQL` - via `useQuery`
- **Error Handling**: Basic error handling
- **Loading State**: Loading state available
- **Refetch**: Refetch function available

## Query Usage Patterns

### Most Frequently Used Queries
1. **`XRPAccountBalancesGQL`** - Used in 6 files
   - Pages: xrp-balances.vue, xrp-portfolio.vue
   - Components: xrpBalances.vue, XrpBalanceWidget.vue
   - Composables: useXrpGraphQL.ts, useXrpToken.ts

2. **`XRPAccountTransactionsGQL`** - Used in 5 files
   - Pages: xrp-transactions.vue
   - Components: xrpAccountHistory.vue
   - Composables: useXrpGraphQL.ts, useXrpToken.ts, useXrpScrerener.ts

3. **`XRPScreenerGQL`** - Used in 4 files
   - Pages: xrp-screener.vue
   - Composables: useXrpGraphQL.ts, useXrpToken.ts, useXrpTokenHeatmap.ts

4. **`XRPAmmPoolsGQL`** - Used in 2 files
   - Pages: xrp-amm-pools.vue
   - Composables: useXrpGraphQL.ts

### Error Handling Patterns
1. **Basic Error Handling** (Most Common)
   - Uses `error` state from `useQuery`
   - Displays error messages in UI
   - Examples: xrp-transactions.vue, xrp-balances.vue

2. **onError Callbacks**
   - Uses `onError` callback for custom error handling
   - Examples: xrp-amm-pools.vue, xrpBalances.vue, useXrpToken.ts

3. **No Error Handling** (Problematic)
   - No explicit error handling
   - Examples: xrp-screener.vue, xrp-explorer pages, most composables

### Loading State Patterns
1. **Explicit Loading States**
   - Uses `loading` state from `useQuery`
   - Examples: xrp-transactions.vue, xrp-balances.vue, useXrpGraphQL.ts

2. **No Loading States** (Problematic)
   - No explicit loading state management
   - Examples: Most components and composables

### Refetch Patterns
1. **Refetch Functions Available**
   - Uses `refetch` function from `useQuery`
   - Examples: xrp-transactions.vue, xrp-balances.vue, useXrpGraphQL.ts

2. **No Refetch Capability**
   - No refetch function exposed
   - Examples: Most components and composables

## Critical Issues Identified

### 1. Inconsistent Error Handling
- **Problem**: Many components have no error handling
- **Impact**: Server crashes when queries fail
- **Priority**: HIGH

### 2. Missing Loading States
- **Problem**: Most components don't show loading states
- **Impact**: Poor user experience during data fetching
- **Priority**: MEDIUM

### 3. No Retry Logic
- **Problem**: No automatic retry on query failures
- **Impact**: Queries fail permanently on temporary issues
- **Priority**: HIGH

### 4. No Offline Handling
- **Problem**: No fallback when GraphQL server is unavailable
- **Impact**: Complete app failure when backend is down
- **Priority**: HIGH

### 5. Duplicate Query Usage
- **Problem**: Same queries used in multiple places
- **Impact**: Inconsistent error handling and loading states
- **Priority**: MEDIUM

## Recommendations for Step 2

### High Priority Files for Error Handling
1. `pages/xrp-screener.vue` - No error handling
2. `pages/xrp-explorer/tx/_id.vue` - No error handling
3. `pages/xrp-explorer/ledger/_id.vue` - No error handling
4. `components/xrp/xrpBalances.vue` - Basic error handling
5. `components/xrp/xrpAccountHistory.vue` - Basic error handling
6. All composables - Inconsistent error handling

### Files to Prioritize for Logging
1. `composables/useXrpGraphQL.ts` - Central query management
2. `pages/xrp-screener.vue` - High-traffic page
3. `pages/xrp-transactions.vue` - Critical functionality
4. `pages/xrp-balances.vue` - Critical functionality

## Next Steps
- Complete Step 1.3: Create comprehensive query registry
- Begin Step 2: Implement logging for XRP-specific queries
- Focus on critical files identified above 