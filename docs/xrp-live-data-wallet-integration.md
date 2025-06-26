# XRP Live Data and Wallet Integration

## Overview

This document describes the implementation of live data integration and wallet connectivity for the XRP AMM functionality. The system provides real-time data from GraphQL APIs and supports multiple wallet types for transaction signing and submission.

## Architecture

### 1. GraphQL API Integration

#### AMM-Specific Queries

The system includes comprehensive GraphQL queries for AMM operations:

```typescript
// AMM Pools
export const XRPAmmPoolsGQL = gql`
  query XRPAmmPoolsGQL {
    xrpAmmPools {
      id
      token1 { symbol, name, icon, issuer }
      token2 { symbol, name, icon, issuer }
      liquidity
      volume24h
      fee
      apr
      priceChange24h
      token1Balance
      token2Balance
      totalSupply
      createdAt
      lastUpdated
    }
  }
`

// Pool Details
export const XRPAmmPoolDetailsGQL = gql`
  query XRPAmmPoolDetailsGQL($poolId: String!) {
    xrpAmmPoolDetails(poolId: $poolId) {
      // ... pool details with transactions and user positions
    }
  }
`

// User Positions
export const XRPAmmUserPositionsGQL = gql`
  query XRPAmmUserPositionsGQL($address: String!) {
    xrpAmmUserPositions(address: $address) {
      // ... user positions across all pools
    }
  }
`

// Swap Quotes
export const XRPAmmQuoteGQL = gql`
  query XRPAmmQuoteGQL($poolId: String!, $amount: String!, $fromToken: String!) {
    xrpAmmQuote(poolId: $poolId, amount: $amount, fromToken: $fromToken) {
      inputAmount
      outputAmount
      priceImpact
      fee
      minimumReceived
      price
    }
  }
`
```

#### Live Data Composable

The `useXrpAmmLiveData` composable provides:

- **Real-time polling** with configurable intervals
- **Automatic data refresh** when wallet connects
- **Error handling** and retry logic
- **Query functions** for specific data needs
- **User-specific data** when wallet is connected

```typescript
const {
  ammPools,           // All AMM pools
  userPositions,      // User's positions
  totalUserValue,     // Total value of user positions
  loading,           // Loading state
  error,             // Error state
  refreshAll,        // Manual refresh
  getPoolDetails,    // Get specific pool details
  getQuote,          // Get swap quote
  getTokenPrice,     // Get token price
  getUserTokenBalances // Get user token balances
} = useXrpAmmLiveData()
```

### 2. Wallet Integration

#### Enhanced XRP Wallet

The `useEnhancedXrpWallet` composable supports multiple wallet types:

- **GEM Wallet** - Native XRP wallet
- **Xaman (XUMM)** - Popular XRP wallet
- **MetaMask XRP Snap** - MetaMask extension

```typescript
const {
  connectWallet,           // Connect specific wallet type
  disconnectWallet,        // Disconnect wallet
  isWalletReady,          // Wallet connection status
  address,                // Wallet address
  currentWalletType,      // Current wallet type
  canSignTransactions,    // Can sign transactions
  signAMMDepositTransaction,   // Sign AMM deposit
  signAMMWithdrawTransaction,  // Sign AMM withdraw
  signAMMTradeTransaction,     // Sign AMM trade
} = useEnhancedXrpWallet()
```

#### Transaction Signing

Each wallet type provides transaction signing capabilities:

```typescript
// AMM Deposit Transaction
const depositParams = {
  amount: "1000",
  amount2Currency: "USDC",
  amount2Issuer: "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
  amount2Value: "950"
}

const signedTx = await signAMMDepositTransaction(depositParams)

// AMM Trade Transaction
const tradeParams = {
  amount: "100",
  amount2Currency: "USDC",
  amount2Issuer: "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
  amount2Value: "95.5"
}

const signedTx = await signAMMTradeTransaction(tradeParams)
```

### 3. XRPL Transaction Submission

#### Transaction Composable

The `useXrplTransaction` composable handles:

- **Transaction submission** to XRPL
- **Confirmation tracking** with timeout
- **Error handling** and retry logic
- **Transaction status** monitoring

```typescript
const {
  submitTransaction,      // Submit transaction
  waitForConfirmation,    // Wait for confirmation
  submitAndConfirm,       // Submit and confirm
  getTransactionStatus,   // Get transaction status
  submitting,            // Submission state
  confirming,            // Confirmation state
  error,                 // Error state
} = useXrplTransaction()
```

#### Transaction Flow

1. **Sign Transaction** - Wallet signs the transaction
2. **Submit to XRPL** - Submit via WebSocket connection
3. **Wait for Confirmation** - Poll for ledger inclusion
4. **Update UI** - Refresh data and show results

```typescript
// Complete transaction flow
const executeTransaction = async (transactionParams) => {
  // 1. Sign with wallet
  const signedTx = await signTransaction(transactionParams)
  
  // 2. Submit and confirm
  const receipt = await submitAndConfirm(signedTx)
  
  // 3. Refresh data
  await refreshAll()
  
  return receipt
}
```

## Implementation Details

### 1. AMM Transactions

#### Deposit/Withdraw

```typescript
// useXrpAmmTransactions composable
const { deposit, withdraw, txLoading, receipt, error } = useXrpAmmTransactions(pool, amount)

// Deposit to pool
await deposit()

// Withdraw from pool
await withdraw()
```

#### Swap Operations

```typescript
// useXrpAmmSwap composable
const { swap, loading, quote, errorMessage } = useXrpAmmSwap(fromToken, toToken, amount, pool)

// Execute swap
await swap()
```

### 2. Real-time Data Updates

#### Polling Configuration

```typescript
// Configurable polling intervals
const pollInterval = ref(30000) // 30 seconds for pools
const poolDetailsInterval = 15000 // 15 seconds for pool details
const transactionsInterval = 10000 // 10 seconds for transactions
const tokenPricesInterval = 60000 // 1 minute for token prices
```

#### Data Refresh Triggers

- **Wallet connection** - Load user-specific data
- **Transaction completion** - Refresh all data
- **Manual refresh** - User-initiated refresh
- **Periodic polling** - Automatic updates

### 3. Error Handling

#### GraphQL Errors

```typescript
// Error handling in live data composable
const handleError = (err: any, context: string) => {
  console.error(`Error in ${context}:`, err)
  error.value = `Failed to load ${context}: ${err.message}`
}

// Watch for errors
watch(ammPoolsError, (err) => {
  if (err) handleError(err, 'AMM pools')
})
```

#### Transaction Errors

```typescript
// Transaction error handling
try {
  const receipt = await submitAndConfirm(signedTx)
  // Handle success
} catch (err) {
  error.value = err.message || 'Transaction failed'
  // Handle error
}
```

## Usage Examples

### 1. Display Live AMM Data

```vue
<template>
  <div>
    <div v-if="loading">Loading...</div>
    <div v-else>
      <div v-for="pool in ammPools" :key="pool.id">
        <h3>{{ pool.token1.symbol }}/{{ pool.token2.symbol }}</h3>
        <p>Liquidity: {{ formatNumber(pool.liquidity) }}</p>
        <p>Volume 24h: {{ formatNumber(pool.volume24h) }}</p>
        <p>APR: {{ formatPercentage(pool.apr) }}%</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  setup() {
    const { ammPools, loading } = useXrpAmmLiveData()
    const { formatNumber, formatPercentage } = useXrpFormatters()
    
    return {
      ammPools,
      loading,
      formatNumber,
      formatPercentage
    }
  }
}
</script>
```

### 2. Execute AMM Swap

```vue
<template>
  <div>
    <input v-model="amount" placeholder="Amount" />
    <button 
      :disabled="!actionButton.status" 
      @click="swap"
    >
      {{ actionButton.message }}
    </button>
    
    <div v-if="quote">
      <p>Quote: {{ quote }}</p>
    </div>
    
    <div v-if="receipt">
      <p>Transaction Hash: {{ receipt.hash }}</p>
      <p>Status: {{ receipt.status === 1 ? 'Success' : 'Failed' }}</p>
    </div>
  </div>
</template>

<script>
export default {
  setup() {
    const { swap, actionButton, quote, receipt, loading } = useXrpAmmSwap(
      fromToken, 
      toToken, 
      amount, 
      pool
    )
    
    return {
      swap,
      actionButton,
      quote,
      receipt,
      loading
    }
  }
}
</script>
```

### 3. Connect Wallet

```vue
<template>
  <div>
    <div v-if="!isWalletReady">
      <button @click="connectWallet('xaman')">Connect Xaman</button>
      <button @click="connectWallet('metamask-xrp-snap')">Connect MetaMask</button>
    </div>
    
    <div v-else>
      <p>Connected: {{ address }}</p>
      <p>Wallet: {{ currentWalletType }}</p>
      <button @click="disconnectWallet">Disconnect</button>
    </div>
  </div>
</template>

<script>
export default {
  setup() {
    const { 
      connectWallet, 
      disconnectWallet, 
      isWalletReady, 
      address, 
      currentWalletType 
    } = useEnhancedXrpWallet()
    
    return {
      connectWallet,
      disconnectWallet,
      isWalletReady,
      address,
      currentWalletType
    }
  }
}
</script>
```

## Configuration

### 1. GraphQL Endpoint

Configure the GraphQL endpoint in your environment:

```env
GRAPHQL_ENDPOINT=https://your-backend.com/graphql
```

### 2. XRPL Network

Configure XRPL network settings:

```typescript
// XRPL network configuration
const XRPL_CONFIG = {
  mainnet: {
    wsUrl: 'wss://xrplcluster.com',
    explorer: 'https://livenet.xrpl.org'
  },
  testnet: {
    wsUrl: 'wss://s.altnet.rippletest.net:51233',
    explorer: 'https://testnet.xrpl.org'
  }
}
```

### 3. Polling Intervals

Configure polling intervals for different data types:

```typescript
// Polling configuration
const POLLING_CONFIG = {
  ammPools: 30000,        // 30 seconds
  poolDetails: 15000,     // 15 seconds
  transactions: 10000,    // 10 seconds
  tokenPrices: 60000,     // 1 minute
  userBalances: 30000     // 30 seconds
}
```

## Testing

### 1. Mock Data

The system includes mock data for development:

```typescript
// Mock AMM pools data
const mockAmmPools = [
  {
    id: 'XRP_USDC_pool',
    token1: { symbol: 'XRP', name: 'Ripple', icon: '🪙' },
    token2: { symbol: 'USDC', name: 'USD Coin', icon: '💎' },
    liquidity: 2500000,
    volume24h: 125000,
    fee: 0.003,
    apr: 12.5
  }
]
```

### 2. Transaction Simulation

Transactions are simulated for development:

```typescript
// Simulate transaction submission
const submitTransaction = async (signedTx: any) => {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  // Generate mock transaction hash
  const mockHash = '0x' + Math.random().toString(16).substr(2, 64)
  
  return {
    hash: mockHash,
    status: 'success',
    ledgerIndex: Math.floor(Math.random() * 1000000),
    timestamp: Date.now()
  }
}
```

## Performance Considerations

### 1. Polling Optimization

- **Configurable intervals** based on data volatility
- **Conditional polling** only when needed
- **Error backoff** to prevent excessive requests

### 2. Caching

- **GraphQL caching** with Apollo Client
- **Local state management** for frequently accessed data
- **Optimistic updates** for better UX

### 3. Memory Management

- **Cleanup on unmount** to prevent memory leaks
- **Limited data retention** for large datasets
- **Efficient re-renders** with computed properties

## Security Considerations

### 1. Wallet Security

- **Secure transaction signing** with wallet providers
- **No private key storage** in the application
- **Transaction validation** before submission

### 2. Data Validation

- **Input sanitization** for all user inputs
- **Transaction parameter validation** before signing
- **Error handling** for malformed data

### 3. Network Security

- **HTTPS/WSS connections** for all API calls
- **Request signing** for authenticated endpoints
- **Rate limiting** to prevent abuse

## Future Enhancements

### 1. WebSocket Integration

- **Real-time updates** via WebSocket subscriptions
- **Live transaction tracking** with immediate updates
- **Market data streaming** for price updates

### 2. Advanced Features

- **Batch transactions** for multiple operations
- **Transaction queuing** for better UX
- **Advanced error recovery** with retry mechanisms

### 3. Performance Improvements

- **Data compression** for large datasets
- **Lazy loading** for non-critical data
- **Background sync** for offline support

## Conclusion

The live data and wallet integration provides a robust foundation for XRP AMM operations. The system supports multiple wallet types, real-time data updates, and comprehensive error handling. The modular architecture allows for easy extension and maintenance.

Key benefits:
- **Real-time data** with configurable polling
- **Multiple wallet support** for user flexibility
- **Comprehensive error handling** for reliability
- **Modular architecture** for maintainability
- **Performance optimized** for smooth UX 