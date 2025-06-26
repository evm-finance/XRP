# Mock Data Replacement Plan - High Priority

## Overview
This document outlines the systematic replacement of mock data with live GraphQL data across the XRP frontend application. This is a **HIGH PRIORITY** task to transition from development to production-ready functionality.

## Current Mock Data Usage Analysis

### 1. **CRITICAL - Account & Transaction Data**
**Files with Mock Data:**
- `composables/useXrpMockData.ts` - Main mock data source
- `components/xrp/xrpAccountHistory.vue` - Transaction history
- `components/xrp/xrpBalances.vue` - Account balances
- `pages/xrp-balances.vue` - Balance display page
- `pages/xrp-transactions.vue` - Transaction display page

**GraphQL Queries Available:**
- ✅ `XRPAccountBalancesGQL` - Account balances
- ✅ `XRPAccountTransactionsGQL` - Transaction history
- ✅ `XRPTransactionGQL` - Individual transaction details

**Status:** READY TO IMPLEMENT

### 2. **CRITICAL - Token & Screener Data**
**Files with Mock Data:**
- `composables/useXrpToken.ts` - Token details
- `pages/xrp-token/_id.vue` - Token page
- `pages/xrp-screener.vue` - Token screener

**GraphQL Queries Available:**
- ✅ `XRPScreenerGQL` - Token screener data
- ✅ `XRPAccountBalancesGQL` - Token balances

**Status:** READY TO IMPLEMENT

### 3. **HIGH PRIORITY - AMM Pool Data**
**Files with Mock Data:**
- `pages/xrp-amm-pools.vue` - AMM pools list
- `pages/xrp-amm-pools/_id.vue` - Individual pool details
- `pages/xrp-amm-heatmap.vue` - AMM heatmap
- `composables/useXrpTokenMints.ts` - Token mints and pools

**GraphQL Queries Available:**
- ✅ `XRPAmmPoolsGQL` - All AMM pools
- ✅ `XRPAmmPoolDetailsGQL` - Individual pool details
- ✅ `XRPAmmUserPositionsGQL` - User positions

**Status:** READY TO IMPLEMENT

### 4. **MEDIUM PRIORITY - Heatmap Data**
**Files with Mock Data:**
- `composables/useXrpHeatmap.ts` - Heatmap data
- `pages/xrp-heatmap.vue` - Token heatmap

**GraphQL Queries Available:**
- ❌ Missing heatmap-specific queries
- ⚠️ Need to create `XRPHeatmapGQL` query

**Status:** NEEDS QUERY CREATION

## Implementation Plan

### Phase 1: Account & Transaction Data (Week 1)
**Priority: CRITICAL**

#### Day 1-2: Replace Account Balance Mock Data
1. **Update `components/xrp/xrpBalances.vue`**
   - Remove `generateMockData()` calls
   - Connect to `XRPAccountBalancesGQL`
   - Add proper error handling and loading states
   - Test with real addresses

2. **Update `components/xrp/xrpAccountHistory.vue`**
   - Remove `generateMockTransactions()` calls
   - Connect to `XRPAccountTransactionsGQL`
   - Add proper error handling and loading states
   - Test with real addresses

#### Day 3-4: Replace Transaction Display Mock Data
1. **Update `pages/xrp-balances.vue`**
   - Remove mock data generation
   - Connect to live GraphQL queries
   - Add proper error handling

2. **Update `pages/xrp-transactions.vue`**
   - Remove mock data generation
   - Connect to live GraphQL queries
   - Add proper error handling

#### Day 5: Clean Up Mock Data Composable
1. **Update `composables/useXrpMockData.ts`**
   - Remove account and transaction mock data
   - Keep only essential testing data
   - Add deprecation warnings

### Phase 2: Token & Screener Data (Week 2)
**Priority: CRITICAL**

#### Day 1-2: Replace Token Page Mock Data
1. **Update `composables/useXrpToken.ts`**
   - Remove `generateMockData()` calls
   - Connect to `XRPScreenerGQL` and `XRPAccountBalancesGQL`
   - Add proper error handling and loading states

2. **Update `pages/xrp-token/_id.vue`**
   - Remove mock data generation
   - Connect to live GraphQL queries
   - Add proper error handling

#### Day 3-4: Replace Screener Mock Data
1. **Update `pages/xrp-screener.vue`**
   - Remove mock data generation
   - Connect to `XRPScreenerGQL`
   - Add proper error handling and loading states

#### Day 5: Test Token Functionality
1. **Test all token-related pages**
   - Verify data loading from GraphQL
   - Test error scenarios
   - Verify navigation between pages

### Phase 3: AMM Pool Data (Week 3)
**Priority: HIGH**

#### Day 1-2: Replace AMM Pools Mock Data
1. **Update `pages/xrp-amm-pools.vue`**
   - Remove mock AMM data
   - Connect to `XRPAmmPoolsGQL`
   - Add proper error handling and loading states

2. **Update `pages/xrp-amm-pools/_id.vue`**
   - Remove mock pool details
   - Connect to `XRPAmmPoolDetailsGQL`
   - Add proper error handling

#### Day 3-4: Replace Token Mints Mock Data
1. **Update `composables/useXrpTokenMints.ts`**
   - Remove `generateMockData()` calls
   - Connect to AMM pool queries
   - Add proper error handling

2. **Update `pages/xrp-amm-heatmap.vue`**
   - Remove mock AMM data
   - Connect to live AMM pool data
   - Add proper error handling

#### Day 5: Test AMM Functionality
1. **Test all AMM-related pages**
   - Verify data loading from GraphQL
   - Test error scenarios
   - Verify navigation between pages

### Phase 4: Heatmap Data (Week 4)
**Priority: MEDIUM**

#### Day 1-2: Create Heatmap GraphQL Queries
1. **Create `apollo/main/heatmap.query.graphql`**
   - Add `XRPHeatmapGQL` query
   - Add `XRPAMMHeatmapGQL` query
   - Add proper TypeScript types

#### Day 3-4: Replace Heatmap Mock Data
1. **Update `composables/useXrpHeatmap.ts`**
   - Remove `generateMockHeatmapData()` calls
   - Connect to new heatmap queries
   - Add proper error handling

2. **Update `pages/xrp-heatmap.vue`**
   - Remove mock data generation
   - Connect to live heatmap data
   - Add proper error handling

#### Day 5: Test Heatmap Functionality
1. **Test heatmap pages**
   - Verify data loading from GraphQL
   - Test error scenarios
   - Verify navigation to token/pool pages

## Technical Implementation Details

### 1. Error Handling Strategy
```typescript
// Standard error handling pattern
const { result, loading, error } = useQuery(XRPAccountBalancesGQL, 
  { account: address },
  () => ({
    enabled: !!address,
    fetchPolicy: 'cache-and-network',
    errorPolicy: 'all',
    notifyOnNetworkStatusChange: true
  })
)

// Handle loading states
if (loading.value) {
  return { loading: true, data: null, error: null }
}

// Handle errors
if (error.value) {
  console.error('GraphQL Error:', error.value)
  return { loading: false, data: null, error: error.value.message }
}

// Return data
return { loading: false, data: result.value?.xrpAccountBalances, error: null }
```

### 2. Loading States
```vue
<template>
  <div>
    <v-skeleton-loader v-if="loading" type="table-tbody" />
    <v-data-table v-else-if="!error" :items="data" />
    <v-alert v-else type="error">{{ error }}</v-alert>
  </div>
</template>
```

### 3. Data Transformation
```typescript
// Transform GraphQL data to component format
const transformAccountData = (graphqlData: any) => {
  return {
    account: graphqlData.account,
    xrpBalance: graphqlData.xrpBalance,
    xrpPrice: graphqlData.xrpPrice,
    tokens: graphqlData.xrpTokens.map((token: any) => ({
      symbol: token.symbol,
      issuer: token.issuer,
      name: token.name,
      balance: token.balance,
      price: token.price,
      value: token.value
    }))
  }
}
```

### 4. Fallback Strategy
```typescript
// Fallback to mock data only in development
const useLiveData = () => {
  const isDevelopment = process.env.NODE_ENV === 'development'
  
  if (isDevelopment && !graphqlData) {
    console.warn('Falling back to mock data in development')
    return mockData
  }
  
  return graphqlData
}
```

## Testing Strategy

### 1. Unit Tests
- Test data transformation functions
- Test error handling scenarios
- Test loading states

### 2. Integration Tests
- Test GraphQL query execution
- Test component rendering with live data
- Test error scenarios

### 3. Manual Testing
- Test with real XRP addresses
- Test with various data scenarios
- Test error handling and loading states

## Success Criteria

### Phase 1 Success Criteria
- [ ] All account balance pages load live data
- [ ] All transaction history pages load live data
- [ ] No mock data used for account/transaction display
- [ ] Proper error handling implemented
- [ ] Loading states implemented

### Phase 2 Success Criteria
- [ ] All token pages load live data
- [ ] Screener loads live data
- [ ] No mock data used for token display
- [ ] Proper error handling implemented
- [ ] Loading states implemented

### Phase 3 Success Criteria
- [ ] All AMM pool pages load live data
- [ ] Token mints page loads live data
- [ ] No mock data used for AMM display
- [ ] Proper error handling implemented
- [ ] Loading states implemented

### Phase 4 Success Criteria
- [ ] All heatmap pages load live data
- [ ] No mock data used for heatmap display
- [ ] Proper error handling implemented
- [ ] Loading states implemented

## Risk Mitigation

### 1. GraphQL API Availability
- Implement retry logic for failed requests
- Add fallback to cached data when API is unavailable
- Monitor API health and performance

### 2. Data Consistency
- Implement data validation
- Add data transformation error handling
- Monitor data quality and consistency

### 3. Performance
- Implement proper caching strategies
- Optimize query frequency
- Monitor query performance

## Timeline

- **Week 1**: Account & Transaction Data (CRITICAL)
- **Week 2**: Token & Screener Data (CRITICAL)
- **Week 3**: AMM Pool Data (HIGH)
- **Week 4**: Heatmap Data (MEDIUM)

**Total Estimated Time: 4 weeks**

## Next Steps

1. **Immediate Action**: Start Phase 1 with account balance replacement
2. **Daily Progress**: Update this document with completed tasks
3. **Weekly Review**: Assess progress and adjust timeline if needed
4. **Testing**: Implement comprehensive testing for each phase

This plan ensures systematic removal of mock data while maintaining application stability and user experience. 