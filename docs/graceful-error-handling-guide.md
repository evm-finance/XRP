# Graceful Error Handling Guide for XRP GraphQL Queries

## Overview
This guide explains the comprehensive graceful error handling system implemented for all XRP GraphQL queries. The system prevents server crashes, provides user-friendly error messages, and ensures a smooth user experience even when backend queries fail.

## System Components

### 1. Error Handler Composable (`useXrpGraphQLErrorHandler.ts`)
**Purpose**: Centralized error handling with retry logic and fallback data support

**Key Features**:
- **Error Categorization**: Automatically categorizes errors as network, GraphQL, validation, or unknown
- **Retry Logic**: Implements exponential backoff retry mechanism
- **User-Friendly Messages**: Converts technical errors to user-friendly messages
- **Fallback Data Support**: Integrates with fallback data system
- **Error State Management**: Tracks error state, retry counts, and retry status

**Configuration Options**:
```typescript
interface ErrorHandlerConfig {
  maxRetries?: number        // Default: 3
  retryDelay?: number        // Default: 1000ms
  enableFallback?: boolean   // Default: true
  showUserErrors?: boolean   // Default: true
  logErrors?: boolean        // Default: true
}
```

### 2. Fallback Data Provider (`useXrpFallbackData.ts`)
**Purpose**: Provides realistic mock data when GraphQL queries fail

**Available Fallback Data**:
- **XRP Screener**: Token list with prices, market caps, volumes
- **XRP Transactions**: Transaction history with realistic data
- **XRP Balances**: Account balances with USD values
- **XRP AMM**: Liquidity pool data for heatmaps
- **XRP Token Mints**: New token information
- **XRP Liquidity Pools**: Pool data for DEX displays

**Usage**:
```typescript
const { getFallbackData } = useXrpFallbackData()
const fallbackData = getFallbackData('XRPScreener')
```

### 3. Enhanced Logging Integration (`useXrpGraphQLWithLogging.ts`)
**Purpose**: Integrates error handling with the existing logging system

**Features**:
- **Automatic Error Handling**: All queries automatically get error handling
- **Fallback Data Integration**: Automatically uses fallback data when available
- **Retry Logic**: Built-in retry with exponential backoff
- **Error State Exposure**: Exposes error state to components

### 4. User-Friendly Error Display (`GraphQLErrorDisplay.vue`)
**Purpose**: Displays graceful error messages to users

**Features**:
- **Error Type Icons**: Different icons for different error types
- **Retry Information**: Shows retry count and progress
- **Retry Button**: Manual retry option for users
- **Technical Details**: Optional technical details for developers
- **Dismissible**: Users can dismiss non-critical errors

## Error Types and Handling

### 1. Network Errors
**Causes**: Connection issues, server down, timeout
**User Message**: "Network connection issue. Please check your internet connection and try again."
**Retry Strategy**: 3 retries with exponential backoff
**Fallback**: Uses cached/mock data if available

### 2. GraphQL Errors
**Causes**: Server-side GraphQL errors, schema issues
**User Message**: "Data loading issue. Please refresh the page or try again later."
**Retry Strategy**: 3 retries with exponential backoff
**Fallback**: Uses cached/mock data if available

### 3. Validation Errors
**Causes**: Invalid query parameters, malformed requests
**User Message**: "Invalid request. Please check your input and try again."
**Retry Strategy**: No retry (validation errors won't succeed)
**Fallback**: Uses cached/mock data if available

### 4. Unknown Errors
**Causes**: Unexpected errors, unhandled exceptions
**User Message**: "An unexpected error occurred. Please try again or contact support if the problem persists."
**Retry Strategy**: 3 retries with exponential backoff
**Fallback**: Uses cached/mock data if available

## Implementation Examples

### Basic Usage in Components
```typescript
// In a Vue component
export default defineComponent({
  setup() {
    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    
    const { 
      onResult, 
      errorState, 
      canRetry, 
      clearError,
      refetch 
    } = useLoggedQuery(XRPScreenerGQL, {
      context: {
        queryName: 'XRPScreener',
        component: 'xrp-screener',
        purpose: 'XRP token screener data'
      }
    })

    const handleRetry = () => refetch()
    const handleErrorDismiss = () => clearError()

    return {
      errorState,
      canRetry,
      handleRetry,
      handleErrorDismiss
    }
  }
})
```

### Template Integration
```vue
<template>
  <!-- Error Display -->
  <graph-q-l-error-display
    :error-state="errorState"
    :can-retry="canRetry"
    :max-retries="3"
    @retry="handleRetry"
    @dismiss="handleErrorDismiss"
  />
  
  <!-- Main Content -->
  <div v-if="!errorState.hasError">
    <!-- Your content here -->
  </div>
</template>
```

### Custom Error Handler Configuration
```typescript
// For mutations (fewer retries, no fallback)
const errorHandler = useXrpGraphQLErrorHandler({
  maxRetries: 2,
  retryDelay: 2000,
  enableFallback: false,
  showUserErrors: true
})

// For subscriptions (more retries, silent errors)
const errorHandler = useXrpGraphQLErrorHandler({
  maxRetries: 5,
  retryDelay: 3000,
  enableFallback: true,
  showUserErrors: false
})
```

## Retry Logic

### Exponential Backoff
The system uses exponential backoff to prevent overwhelming the server:

```typescript
const delay = retryDelay * Math.pow(2, retryCount - 1)
```

**Example Retry Delays**:
- 1st retry: 1000ms
- 2nd retry: 2000ms  
- 3rd retry: 4000ms
- 4th retry: 8000ms

### Retry Conditions
- **Queries**: Retry on network and GraphQL errors
- **Mutations**: Retry on network errors only
- **Subscriptions**: Retry on all errors with more attempts

## Fallback Data Strategy

### When Fallback Data is Used
1. **Query fails after all retries**
2. **Fallback data is available for the query type**
3. **Fallback is enabled in configuration**

### Fallback Data Types
```typescript
// Screener fallback
{
  xrpScreener: [
    {
      currency: 'USDC',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'USD Coin',
      price: 1.0,
      marketcap: 1000000,
      volume24H: 500000
    }
    // ... more tokens
  ]
}

// AMM fallback
{
  xrpAMMLiquidityValues: [
    {
      poolId: '1',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'USD', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      totalLiquidityUsd: 10000000,
      fee: 0.3
    }
    // ... more pools
  ]
}
```

## Error State Management

### Error State Interface
```typescript
interface ErrorState {
  hasError: boolean
  errorMessage: string
  errorType: 'network' | 'graphql' | 'validation' | 'unknown'
  retryCount: number
  lastError?: ApolloError
  isRetrying: boolean
}
```

### Error State Lifecycle
1. **Initial**: `hasError: false`
2. **Error Occurs**: `hasError: true`, error details set
3. **Retry Starts**: `isRetrying: true`
4. **Retry Success**: `hasError: false`, state cleared
5. **Retry Fails**: `isRetrying: false`, retry count incremented
6. **Max Retries**: Fallback data used or final error state

## Best Practices

### 1. Always Use Error Display Component
```vue
<graph-q-l-error-display
  :error-state="errorState"
  :can-retry="canRetry"
  @retry="handleRetry"
  @dismiss="handleErrorDismiss"
/>
```

### 2. Handle Different Query Types Appropriately
- **Queries**: Enable fallback, show user errors
- **Mutations**: Disable fallback, show user errors
- **Subscriptions**: Enable fallback, hide user errors

### 3. Provide Meaningful Context
```typescript
context: {
  queryName: 'XRPScreener',
  component: 'xrp-screener',
  purpose: 'XRP token screener data for main page'
}
```

### 4. Test Error Scenarios
- Test network disconnection
- Test server errors
- Test validation errors
- Verify fallback data works
- Check retry logic

## Monitoring and Debugging

### Error Logging
All errors are logged with:
- Error type and message
- Component and query context
- Retry count and timing
- Fallback data usage

### Error Metrics
Track:
- Error frequency by type
- Retry success rates
- Fallback data usage
- User error dismissal rates

### Debug Mode
Enable technical details in error display:
```vue
<graph-q-l-error-display
  :show-technical-details-by-default="true"
/>
```

## Conclusion

The graceful error handling system ensures that:
- **Users never see crashes**: All errors are caught and handled
- **User experience is smooth**: Fallback data provides continuity
- **Errors are recoverable**: Retry logic handles transient issues
- **Errors are informative**: User-friendly messages explain issues
- **Development is easier**: Comprehensive logging and debugging tools

This system makes the XRP application resilient to backend failures and provides a professional user experience even when things go wrong. 