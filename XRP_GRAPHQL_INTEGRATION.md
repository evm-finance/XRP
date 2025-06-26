# XRP GraphQL Integration Guide

This document outlines the integration of new XRP-related GraphQL models from the backend server into the frontend application.

## Backend GraphQL Models

The backend server at `D:/DefiBackend/newest/qc-defi-graphql-server` provides the following XRP-related GraphQL models:

### 1. XRP Account Data
- **XRPAccountBalances**: Account balances including XRP and token balances
- **XRPAccountTransactions**: Account transaction history
- **XRPAccountData**: Combined account data (balances + transactions)

### 2. XRP AMM (Automated Market Maker)
- **XRPAMMPool**: AMM pool information with asset balances and trading volume
- **XRPAMMSwapQuote**: Swap quote calculations with price impact and fees
- **XRPAMMLiquidityCalculation**: Liquidity value calculations in USD
- **AMMTransaction**: AMM transaction history

### 3. XRP Tokens
- **XRPToken**: Token information including metadata, market cap, and volume
- **AMMHeatmapData**: Token heatmap data for visualization
- **AMMPool**: AMM pool details for specific tokens

### 4. XRP Ledger Data
- **XRPLedger**: Ledger information and metadata
- **XRPLedgerData**: Detailed ledger data including transactions

## Frontend Integration

### 1. GraphQL Queries

Created `apollo/main/xrp.query.graphql` with comprehensive queries for:

#### Account Queries
```graphql
query XRPAccountBalances($address: String!) {
  xrpAccountBalances(address: $address) {
    account
    xrpBalance
    xrpPrice
    xrpTokens {
      symbol
      issuer
      name
      balance
      price
      value
    }
  }
}
```

#### AMM Queries
```graphql
query XRPAMMPools {
  xrpAMMPools {
    poolId
    asset1 { currency issuer }
    asset2 { currency issuer }
    asset1Balance
    asset2Balance
    lpBalance
    fee
    tradingVolume24H
    tradingVolume7D
    createdAt
    lastUpdated
  }
}
```

#### Token Queries
```graphql
query XRPTokens($limit: Int, $offset: Int, $sortBy: TokenSortField, $sortOrder: SortOrder) {
  xrpTokens(limit: $limit, offset: $offset, sortBy: $sortBy, sortOrder: $sortOrder) {
    issuer
    currency
    name
    icon
    description
    marketcap
    price
    volume24h
    volume7d
  }
}
```

### 2. TypeScript Types

Updated `types/apollo/main/types.ts` with comprehensive type definitions:

#### Core Types
```typescript
export type XRPAsset = {
  __typename?: 'XRPAsset';
  currency: Scalars['String'];
  issuer?: Maybe<Scalars['String']>;
};

export type XRPAMMPool = {
  __typename?: 'XRPAMMPool';
  poolId: Scalars['String'];
  asset1: XRPAsset;
  asset2: XRPAsset;
  asset1Balance: Scalars['Float'];
  asset2Balance: Scalars['Float'];
  lpBalance: Scalars['Float'];
  fee: Scalars['Float'];
  tradingVolume24H?: Maybe<Scalars['Float']>;
  tradingVolume7D?: Maybe<Scalars['Float']>;
  createdAt: Scalars['Int'];
  lastUpdated: Scalars['Int'];
};
```

#### Enums
```typescript
export enum TimeRange {
  Hour_1 = 'HOUR_1',
  Day_1 = 'DAY_1',
  Day_7 = 'DAY_7',
  Day_30 = 'DAY_30'
}

export enum SortField {
  Liquidity = 'LIQUIDITY',
  Volume_24H = 'VOLUME_24H',
  Volume_7D = 'VOLUME_7D',
  Market_Cap = 'MARKET_CAP',
  Price_Change_24H = 'PRICE_CHANGE_24H'
}
```

### 3. Composables

Created comprehensive composables for XRP functionality:

#### useXrpGraphQL
Main composable for GraphQL-based XRP operations:
- Account data queries
- AMM pool queries
- Token queries
- Real-time subscriptions
- Loading and error states

#### useXrpAmmOperations
Composable for AMM operations:
- Swap functionality
- Liquidity provision/removal
- Pool creation
- Gas estimation

#### useXrpTokenManagement
Composable for token management:
- Trust line operations
- Token minting/burning
- Token transfers
- Freeze/unfreeze operations

## Usage Examples

### 1. Account Data
```typescript
import useXrpGraphQL from '~/composables/useXrpGraphQL'

export default {
  setup() {
    const {
      setAddress,
      accountBalances,
      accountTransactions,
      isLoading,
      hasError
    } = useXrpGraphQL()

    // Set account address
    setAddress('rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX')

    return {
      accountBalances,
      accountTransactions,
      isLoading,
      hasError
    }
  }
}
```

### 2. AMM Operations
```typescript
import useXrpAmmOperations from '~/composables/useXrpAmmOperations'

export default {
  setup() {
    const {
      setSwapParams,
      getSwapQuote,
      executeSwap,
      swapQuote,
      swapLoading,
      swapError
    } = useXrpAmmOperations()

    const performSwap = async () => {
      setSwapParams({
        inputAsset: { currency: 'USD', issuer: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX' },
        outputAsset: { currency: 'XRP' },
        amount: 100,
        slippage: 0.5
      })

      const quote = await getSwapQuote()
      if (quote) {
        const result = await executeSwap()
        console.log('Swap result:', result)
      }
    }

    return {
      performSwap,
      swapQuote,
      swapLoading,
      swapError
    }
  }
}
```

### 3. Token Management
```typescript
import useXrpTokenManagement from '~/composables/useXrpTokenManagement'

export default {
  setup() {
    const {
      setTrustLineParams,
      setTrustLine,
      fetchUserTokens,
      userTokens,
      trustLineLoading
    } = useXrpTokenManagement()

    const setupTrustLine = async () => {
      setTrustLineParams({
        currency: 'USD',
        issuer: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        limit: 10000
      })

      const result = await setTrustLine()
      console.log('Trust line result:', result)
    }

    return {
      setupTrustLine,
      userTokens,
      trustLineLoading
    }
  }
}
```

## Real-time Subscriptions

The integration includes real-time subscriptions for:

### AMM Pool Updates
```typescript
const {
  result: poolUpdatesResult,
  loading: poolUpdatesLoading
} = useSubscription<AMMPoolUpdatesSubscription>(
  AMMPoolUpdates,
  () => ({ poolAccount: selectedPoolId.value }),
  () => ({ enabled: !!selectedPoolId.value })
)
```

### Token Price Updates
```typescript
const {
  result: tokenPriceUpdatesResult,
  loading: tokenPriceUpdatesLoading
} = useSubscription<TokenPriceUpdatesSubscription>(
  TokenPriceUpdates,
  () => ({
    currency: selectedToken.currency,
    issuer: selectedToken.issuer
  }),
  () => ({ enabled: !!selectedToken.currency && !!selectedToken.issuer })
)
```

### Heatmap Updates
```typescript
const {
  result: heatmapUpdatesResult,
  loading: heatmapUpdatesLoading
} = useSubscription<HeatmapUpdatesSubscription>(
  HeatmapUpdates,
  () => ({ timeRange: TimeRange.Day_1 }),
  () => ({ enabled: true })
)
```

## Configuration

### 1. Apollo Client Configuration
The Apollo client is configured in `nuxt.config.js`:
```javascript
apollo: {
  clientConfigs: {
    default: {
      httpEndpoint: process.env.BASE_GRAPHQL_SERVER_URL,
      wsEndpoint: process.env.BASE_GRAPHQL_WEBSOCKET_URL,
      websocketsOnly: false,
    },
  },
},
```

### 2. Environment Variables
Ensure the following environment variables are set:
```bash
BASE_GRAPHQL_SERVER_URL=http://localhost:8080/query
BASE_GRAPHQL_WEBSOCKET_URL=ws://localhost:8080/query
```

## Error Handling

All composables include comprehensive error handling:

```typescript
const hasError = computed(() => 
  accountBalancesError.value ||
  accountTransactionsError.value ||
  ammPoolsError.value ||
  // ... other error states
)
```

## Loading States

Loading states are provided for all operations:

```typescript
const isLoading = computed(() => 
  accountBalancesLoading.value ||
  accountTransactionsLoading.value ||
  ammPoolsLoading.value ||
  // ... other loading states
)
```

## Next Steps

1. **Update Components**: Update existing XRP components to use the new GraphQL composables
2. **Add New Components**: Create new components for AMM operations, token management, etc.
3. **Testing**: Add comprehensive tests for the new GraphQL functionality
4. **Documentation**: Update component documentation to reflect the new GraphQL integration
5. **Performance**: Monitor and optimize GraphQL query performance

## Migration from Existing XRP Code

The existing XRP functionality in the codebase can be gradually migrated to use the new GraphQL-based composables:

1. **Replace Direct XRPL API Calls**: Use GraphQL queries instead of direct XRPL API calls
2. **Update State Management**: Use the reactive state from composables instead of manual state management
3. **Implement Real-time Updates**: Use subscriptions for real-time data updates
4. **Standardize Error Handling**: Use the standardized error handling from composables

This integration provides a robust, type-safe, and scalable foundation for XRP functionality in the frontend application. 