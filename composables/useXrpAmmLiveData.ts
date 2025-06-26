import { ref, computed, watch, onMounted, onBeforeUnmount } from '@nuxtjs/composition-api'
import { useQuery, useSubscription } from '@vue/apollo-composable'
import { 
  XRPAmmPoolsGQL, 
  XRPAmmPoolDetailsGQL, 
  XRPAmmUserPositionsGQL,
  XRPAmmQuoteGQL,
  XRPAmmTransactionsGQL,
  XRPTokenPriceGQL,
  XRPTokenBalancesGQL
} from '~/apollo/queries'
import { useEnhancedXrpWallet } from '~/composables/useEnhancedXrpWallet'

interface XrpToken {
  symbol: string
  name: string
  icon: string
  issuer?: string
  price?: number
  priceChange24h?: number
}

interface XrpAmmPool {
  id: string
  token1: XrpToken
  token2: XrpToken
  liquidity: number
  volume24h: number
  fee: number
  apr: number
  priceChange24h: number
  token1Balance: number
  token2Balance: number
  totalSupply: number
  createdAt: string
  lastUpdated: string
}

interface XrpAmmPoolDetails extends XrpAmmPool {
  transactions: {
    hash: string
    type: string
    amount: number
    token: string
    timestamp: string
    user: string
  }[]
  userPositions: {
    user: string
    poolTokens: number
    token1Balance: number
    token2Balance: number
    share: number
    value: number
  }[]
}

interface XrpAmmQuote {
  inputAmount: string
  outputAmount: string
  priceImpact: number
  fee: number
  minimumReceived: string
  price: number
}

interface XrpUserPosition {
  poolId: string
  pool: XrpAmmPool
  poolTokens: number
  token1Balance: number
  token2Balance: number
  share: number
  value: number
  lastUpdated: string
}

export default function useXrpAmmLiveData() {
  const { address, isWalletReady } = useEnhancedXrpWallet()
  
  // State
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastUpdate = ref<Date | null>(null)
  
  // Polling configuration
  const pollInterval = ref(30000) // 30 seconds
  const isPolling = ref(false)
  let pollTimer: NodeJS.Timeout | null = null

  // AMM Pools Query
  const { 
    result: ammPoolsResult, 
    loading: ammPoolsLoading, 
    error: ammPoolsError,
    refetch: refetchAmmPools 
  } = useQuery(XRPAmmPoolsGQL, null, {
    fetchPolicy: 'cache-and-network',
    pollInterval: pollInterval.value,
    errorPolicy: 'all'
  })

  // User Positions Query (only when wallet is connected)
  const { 
    result: userPositionsResult, 
    loading: userPositionsLoading, 
    error: userPositionsError,
    refetch: refetchUserPositions 
  } = useQuery(XRPAmmUserPositionsGQL, 
    computed(() => ({ address: address.value })),
    () => ({
      enabled: isWalletReady.value && !!address.value,
      fetchPolicy: 'cache-and-network',
      pollInterval: pollInterval.value,
      errorPolicy: 'all'
    })
  )

  // Computed data
  const ammPools = computed<XrpAmmPool[]>(() => {
    return ammPoolsResult.value?.xrpAmmPools || []
  })

  const userPositions = computed<XrpUserPosition[]>(() => {
    return userPositionsResult.value?.xrpAmmUserPositions || []
  })

  const totalUserValue = computed(() => {
    return userPositions.value.reduce((sum, position) => sum + position.value, 0)
  })

  // Pool details query function
  const getPoolDetails = (poolId: string) => {
    const { result, loading, error, refetch } = useQuery(
      XRPAmmPoolDetailsGQL,
      { poolId },
      {
        fetchPolicy: 'cache-and-network',
        pollInterval: 15000, // 15 seconds for pool details
        errorPolicy: 'all'
      }
    )

    return {
      poolDetails: computed<XrpAmmPoolDetails | null>(() => result.value?.xrpAmmPoolDetails || null),
      loading,
      error,
      refetch
    }
  }

  // Quote query function
  const getQuote = (poolId: string, amount: string, fromToken: string) => {
    const { result, loading, error, refetch } = useQuery(
      XRPAmmQuoteGQL,
      { poolId, amount, fromToken },
      {
        fetchPolicy: 'cache-first',
        errorPolicy: 'all'
      }
    )

    return {
      quote: computed<XrpAmmQuote | null>(() => result.value?.xrpAmmQuote || null),
      loading,
      error,
      refetch
    }
  }

  // Pool transactions query function
  const getPoolTransactions = (poolId: string, limit: number = 50) => {
    const { result, loading, error, refetch } = useQuery(
      XRPAmmTransactionsGQL,
      { poolId, limit },
      {
        fetchPolicy: 'cache-and-network',
        pollInterval: 10000, // 10 seconds for transactions
        errorPolicy: 'all'
      }
    )

    return {
      transactions: computed(() => result.value?.xrpAmmTransactions || []),
      loading,
      error,
      refetch
    }
  }

  // Token price query function
  const getTokenPrice = (currency: string, issuer?: string) => {
    const { result, loading, error, refetch } = useQuery(
      XRPTokenPriceGQL,
      { currency, issuer },
      {
        fetchPolicy: 'cache-and-network',
        pollInterval: 60000, // 1 minute for token prices
        errorPolicy: 'all'
      }
    )

    return {
      tokenPrice: computed(() => result.value?.xrpTokenPrice || null),
      loading,
      error,
      refetch
    }
  }

  // User token balances query function
  const getUserTokenBalances = () => {
    const { result, loading, error, refetch } = useQuery(
      XRPTokenBalancesGQL,
      computed(() => ({ address: address.value })),
      () => ({
        enabled: isWalletReady.value && !!address.value,
        fetchPolicy: 'cache-and-network',
        pollInterval: 30000, // 30 seconds for balances
        errorPolicy: 'all'
      })
    )

    return {
      tokenBalances: computed(() => result.value?.xrpTokenBalances || []),
      loading,
      error,
      refetch
    }
  }

  // Error handling
  const handleError = (err: any, context: string) => {
    console.error(`Error in ${context}:`, err)
    error.value = `Failed to load ${context}: ${err.message}`
  }

  // Watch for errors
  watch(ammPoolsError, (err) => {
    if (err) handleError(err, 'AMM pools')
  })

  watch(userPositionsError, (err) => {
    if (err) handleError(err, 'user positions')
  })

  // Polling management
  const startPolling = () => {
    if (isPolling.value) return
    
    isPolling.value = true
    pollTimer = setInterval(() => {
      refetchAmmPools()
      if (isWalletReady.value && address.value) {
        refetchUserPositions()
      }
      lastUpdate.value = new Date()
    }, pollInterval.value)
  }

  const stopPolling = () => {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    isPolling.value = false
  }

  // Manual refresh
  const refreshAll = async () => {
    loading.value = true
    error.value = null
    
    try {
      await Promise.all([
        refetchAmmPools(),
        isWalletReady.value && address.value ? refetchUserPositions() : Promise.resolve()
      ])
      lastUpdate.value = new Date()
    } catch (err: any) {
      handleError(err, 'data refresh')
    } finally {
      loading.value = false
    }
  }

  // Lifecycle
  onMounted(() => {
    startPolling()
  })

  onBeforeUnmount(() => {
    stopPolling()
  })

  // Watch wallet connection for user data
  watch([isWalletReady, address], ([ready, addr]) => {
    if (ready && addr) {
      refetchUserPositions()
    }
  })

  return {
    // State
    loading: computed(() => ammPoolsLoading.value || userPositionsLoading.value || loading.value),
    error,
    lastUpdate,
    isPolling,

    // Data
    ammPools,
    userPositions,
    totalUserValue,

    // Query functions
    getPoolDetails,
    getQuote,
    getPoolTransactions,
    getTokenPrice,
    getUserTokenBalances,

    // Actions
    refreshAll,
    startPolling,
    stopPolling,
    refetchAmmPools,
    refetchUserPositions,
  }
} 