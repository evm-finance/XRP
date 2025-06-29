# XRP GraphQL Logging Usage Example

## Overview
This document shows how to use the enhanced logging system that captures the exact GraphQL queries being executed.

## Basic Usage

### 1. Using the Logged Query Wrapper

```typescript
// Instead of this:
import { useQuery } from '@vue/apollo-composable'
import { XRPAccountBalancesGQL } from '~/apollo/queries'

const { result, loading, error } = useQuery(
  XRPAccountBalancesGQL,
  () => ({ account: currentAddress.value })
)

// Use this:
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'
import { XRPAccountBalancesGQL } from '~/apollo/queries'

const { useLoggedQuery } = useXrpGraphQLWithLogging()

const { result, loading, error } = useLoggedQuery(
  XRPAccountBalancesGQL,
  () => ({ account: currentAddress.value })
)
```

### 2. What Gets Logged

When you use the logged wrapper, you'll see output like this in the console:

```javascript
// Query Start
🚀 XRP GraphQL Query Started: XRPAccountBalancesGQL {
  id: "XRPAccountBalancesGQL_1703123456789_abc123def",
  operation: "query",
  queryString: "query XRPAccountBalancesGQL($account: String!) { xrpAccountBalances(account: $account) { account xrpBalance xrpPrice xrpTokens { symbol issuer name balance price value } } }",
  variables: { account: "rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj" },
  page: "/xrp-balances",
  timestamp: "2023-12-21T10:30:45.123Z"
}

// Query Success
✅ XRP GraphQL Query Completed: XRPAccountBalancesGQL {
  duration: "245ms",
  page: "/xrp-balances"
}

// Query Failure
❌ XRP GraphQL Query Failed: XRPAccountBalancesGQL {
  error: { networkError: { statusCode: 500, message: "Internal Server Error" } },
  duration: "5000ms",
  page: "/xrp-balances",
  variables: { account: "rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj" },
  queryString: "query XRPAccountBalancesGQL($account: String!) { xrpAccountBalances(account: $account) { account xrpBalance xrpPrice xrpTokens { symbol issuer name balance price value } } }"
}
```

## Integration Examples

### 1. Update useXrpGraphQL.ts

```typescript
// composables/useXrpGraphQL.ts
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'
import { XRPAccountBalancesGQL, XRPAccountTransactionsGQL } from '~/apollo/queries'

export default function useXrpGraphQL() {
  const { useLoggedQuery } = useXrpGraphQLWithLogging()
  
  // Replace useQuery with useLoggedQuery
  const {
    result: accountBalancesResult,
    loading: accountBalancesLoading,
    error: accountBalancesError,
    refetch: refetchAccountBalances
  } = useLoggedQuery(
    XRPAccountBalancesGQL,
    () => ({ account: currentAddress.value }),
    () => ({ enabled: !!currentAddress.value, fetchPolicy: 'cache-and-network' })
  )

  const {
    result: accountTransactionsResult,
    loading: accountTransactionsLoading,
    error: accountTransactionsError,
    refetch: refetchAccountTransactions
  } = useLoggedQuery(
    XRPAccountTransactionsGQL,
    () => ({ address: currentAddress.value }),
    () => ({ enabled: !!currentAddress.value, fetchPolicy: 'cache-and-network' })
  )

  // ... rest of the composable
}
```

### 2. Update Critical Pages

```vue
<!-- pages/xrp-screener.vue -->
<template>
  <div>
    <!-- Your existing template -->
  </div>
</template>

<script>
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'
import { XRPScreenerGQL } from '~/apollo/queries'

export default {
  setup() {
    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    
    // Replace the existing useQuery with logged version
    const { onResult } = useLoggedQuery(
      XRPScreenerGQL, 
      {}, 
      { 
        fetchPolicy: 'no-cache', 
        pollInterval: 60000 
      }
    )

    // ... rest of your setup
  }
}
</script>
```

## Logging Features

### 1. Query String Capture
- **Exact GraphQL query**: The actual query string being sent to the server
- **Formatted for readability**: Whitespace normalized for easier debugging
- **Variables included**: All query parameters are logged

### 2. Error Context
- **Full error details**: Network errors, GraphQL errors, validation errors
- **Query context**: The exact query that failed
- **User context**: Page, wallet address, timestamp
- **Network info**: Endpoint, user agent

### 3. Performance Tracking
- **Query duration**: How long each query takes
- **Slow query identification**: Queries taking longer than expected
- **Success/failure rates**: Track query reliability

### 4. Development vs Production
- **Development**: Detailed console logging for debugging
- **Production**: Structured logging for external services
- **Error categorization**: Network, authentication, validation, etc.

## Benefits

### 1. Debugging
- **Exact query reproduction**: Copy the logged query to test in GraphQL Playground
- **Variable inspection**: See exactly what parameters were sent
- **Error correlation**: Match errors to specific queries

### 2. Performance Monitoring
- **Identify slow queries**: Track which queries are taking too long
- **Success rate tracking**: Monitor query reliability
- **Usage patterns**: Understand which queries are used most

### 3. Error Prevention
- **Early detection**: Catch issues before they affect users
- **Pattern recognition**: Identify recurring error patterns
- **Proactive fixes**: Address issues before they become critical

## Next Steps

1. **Integrate with critical files**:
   - `composables/useXrpGraphQL.ts`
   - `pages/xrp-screener.vue`
   - `pages/xrp-transactions.vue`
   - `pages/xrp-balances.vue`

2. **Add error correlation IDs** for better tracking

3. **Implement external logging service** for production

4. **Create monitoring dashboard** for query performance 