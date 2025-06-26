import { ref, computed, Ref, reactive } from '@nuxtjs/composition-api'
import { useQuery, useSubscription, useMutation } from '@vue/apollo-composable/dist'
import { 
  XRPAccountBalancesQuery,
  XRPAccountTransactionsQuery,
  XRPAccountDataQuery,
  XRPAMMPoolsQuery,
  XRPAMMPoolQuery,
  XRPAMMSwapQuoteQuery,
  XRPAMMLiquidityValueQuery,
  XRPAMMTransactionsQuery,
  AMMHeatmapQuery,
  TopAMMPoolsQuery,
  AMMPoolDetailsQuery,
  XRPTokensQuery,
  XRPTokenQuery,
  AMMPoolsQuery,
  AMMPoolUpdatesSubscription,
  TokenPriceUpdatesSubscription,
  HeatmapUpdatesSubscription,
  XRPAssetInput,
  TimeRange,
  SortField,
  TokenSortField,
  SortOrder
} from '~/types/apollo/main/types'

// Import GraphQL queries
import {
  XRPAccountBalances,
  XRPAccountTransactions,
  XRPAccountData,
  XRPAMMPools,
  XRPAMMPool,
  XRPAMMSwapQuote,
  XRPAMMLiquidityValue,
  XRPAMMTransactions,
  AMMHeatmap,
  TopAMMPools,
  AMMPoolDetails,
  XRPTokens,
  XRPToken,
  AMMPools,
  AMMPoolUpdates,
  TokenPriceUpdates,
  HeatmapUpdates
} from '~/apollo/main/xrp.query.graphql'

export interface XRPAccountState {
  address: string
  xrpBalance: number
  xrpPrice: number
  tokens: Array<{
    symbol: string
    issuer: string
    name: string
    balance: number
    price: number
    value: number
  }>
}

export interface XRPTransactionState {
  account: string
  amount?: any
  destination: string
  fee: string
  flags: number
  lastLedgerSequence: number
  offerSequence: number
  sequence: number
  signingPubKey: string
  takerGets?: any
  takerPays?: any
  transactionType: string
  txnSignature: string
  date: number
  hash: string
  inLedger: number
  ledgerIndex: number
  meta?: any
  metadata?: any
  validated?: boolean
  warnings: Array<any>
  memos: Array<any>
}

export interface XRPAMMPoolState {
  poolId: string
  asset1: { currency: string; issuer?: string }
  asset2: { currency: string; issuer?: string }
  asset1Balance: number
  asset2Balance: number
  lpBalance: number
  fee: number
  tradingVolume24H?: number
  tradingVolume7D?: number
  createdAt: number
  lastUpdated: number
}

export interface XRPSwapQuoteState {
  inputAmount: number
  outputAmount: number
  priceImpact: number
  fee: number
  minimumReceived: number
}

export interface XRPTokenState {
  issuer: string
  currency: string
  name: string
  icon?: string
  description?: string
  marketcap?: number
  price?: number
  volume24h?: number
  volume7d?: number
}

export default function useXrpGraphQL() {
  // Reactive state
  const currentAddress = ref('')
  const selectedPoolId = ref('')
  const selectedToken = reactive<{ currency: string; issuer: string }>({ currency: '', issuer: '' })

  // Account Queries
  const {
    result: accountBalancesResult,
    loading: accountBalancesLoading,
    error: accountBalancesError,
    refetch: refetchAccountBalances
  } = useQuery<XRPAccountBalancesQuery>(
    XRPAccountBalances,
    () => ({ address: currentAddress.value }),
    () => ({ enabled: !!currentAddress.value, fetchPolicy: 'cache-and-network' })
  )

  const {
    result: accountTransactionsResult,
    loading: accountTransactionsLoading,
    error: accountTransactionsError,
    refetch: refetchAccountTransactions
  } = useQuery<XRPAccountTransactionsQuery>(
    XRPAccountTransactions,
    () => ({ address: currentAddress.value }),
    () => ({ enabled: !!currentAddress.value, fetchPolicy: 'cache-and-network' })
  )

  const {
    result: accountDataResult,
    loading: accountDataLoading,
    error: accountDataError,
    refetch: refetchAccountData
  } = useQuery<XRPAccountDataQuery>(
    XRPAccountData,
    () => ({ address: currentAddress.value }),
    () => ({ enabled: !!currentAddress.value, fetchPolicy: 'cache-and-network' })
  )

  // AMM Queries
  const {
    result: ammPoolsResult,
    loading: ammPoolsLoading,
    error: ammPoolsError,
    refetch: refetchAMMPools
  } = useQuery<XRPAMMPoolsQuery>(
    XRPAMMPools,
    null,
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  const {
    result: ammPoolResult,
    loading: ammPoolLoading,
    error: ammPoolError,
    refetch: refetchAMMPool
  } = useQuery<XRPAMMPoolQuery>(
    XRPAMMPool,
    () => ({
      asset1: { currency: selectedToken.currency, issuer: selectedToken.issuer },
      asset2: { currency: 'XRP' }
    }),
    () => ({ 
      enabled: !!selectedToken.currency && !!selectedToken.issuer,
      fetchPolicy: 'cache-and-network' 
    })
  )

  const {
    result: swapQuoteResult,
    loading: swapQuoteLoading,
    error: swapQuoteError,
    refetch: refetchSwapQuote
  } = useQuery<XRPAMMSwapQuoteQuery>(
    XRPAMMSwapQuote,
    () => ({
      inputAsset: { currency: selectedToken.currency, issuer: selectedToken.issuer },
      outputAsset: { currency: 'XRP' },
      amount: 100 // Default amount, should be configurable
    }),
    () => ({ 
      enabled: !!selectedToken.currency && !!selectedToken.issuer,
      fetchPolicy: 'cache-and-network' 
    })
  )

  const {
    result: liquidityValueResult,
    loading: liquidityValueLoading,
    error: liquidityValueError,
    refetch: refetchLiquidityValue
  } = useQuery<XRPAMMLiquidityValueQuery>(
    XRPAMMLiquidityValue,
    () => ({ poolId: selectedPoolId.value }),
    () => ({ 
      enabled: !!selectedPoolId.value,
      fetchPolicy: 'cache-and-network' 
    })
  )

  const {
    result: ammTransactionsResult,
    loading: ammTransactionsLoading,
    error: ammTransactionsError,
    refetch: refetchAMMTransactions
  } = useQuery<XRPAMMTransactionsQuery>(
    XRPAMMTransactions,
    () => ({ 
      poolId: selectedPoolId.value,
      limit: 50 
    }),
    () => ({ 
      enabled: !!selectedPoolId.value,
      fetchPolicy: 'cache-and-network' 
    })
  )

  // Heatmap Queries
  const {
    result: heatmapResult,
    loading: heatmapLoading,
    error: heatmapError,
    refetch: refetchHeatmap
  } = useQuery<AMMHeatmapQuery>(
    AMMHeatmap,
    () => ({
      timeRange: TimeRange.Day_1,
      sortBy: SortField.Liquidity,
      sortOrder: SortOrder.Desc,
      limit: 100
    }),
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  const {
    result: topPoolsResult,
    loading: topPoolsLoading,
    error: topPoolsError,
    refetch: refetchTopPools
  } = useQuery<TopAMMPoolsQuery>(
    TopAMMPools,
    () => ({ limit: 50 }),
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  const {
    result: poolDetailsResult,
    loading: poolDetailsLoading,
    error: poolDetailsError,
    refetch: refetchPoolDetails
  } = useQuery<AMMPoolDetailsQuery>(
    AMMPoolDetails,
    () => ({ poolAccount: selectedPoolId.value }),
    () => ({ 
      enabled: !!selectedPoolId.value,
      fetchPolicy: 'cache-and-network' 
    })
  )

  // Token Queries
  const {
    result: tokensResult,
    loading: tokensLoading,
    error: tokensError,
    refetch: refetchTokens
  } = useQuery<XRPTokensQuery>(
    XRPTokens,
    () => ({
      limit: 100,
      sortBy: TokenSortField.Marketcap,
      sortOrder: SortOrder.Desc
    }),
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  const {
    result: tokenResult,
    loading: tokenLoading,
    error: tokenError,
    refetch: refetchToken
  } = useQuery<XRPTokenQuery>(
    XRPToken,
    () => ({
      currency: selectedToken.currency,
      issuer: selectedToken.issuer
    }),
    () => ({ 
      enabled: !!selectedToken.currency && !!selectedToken.issuer,
      fetchPolicy: 'cache-and-network' 
    })
  )

  const {
    result: ammPoolsForTokenResult,
    loading: ammPoolsForTokenLoading,
    error: ammPoolsForTokenError,
    refetch: refetchAMMPoolsForToken
  } = useQuery<AMMPoolsQuery>(
    AMMPools,
    () => ({
      currency: selectedToken.currency,
      issuer: selectedToken.issuer,
      limit: 50
    }),
    () => ({ 
      enabled: !!selectedToken.currency && !!selectedToken.issuer,
      fetchPolicy: 'cache-and-network' 
    })
  )

  // Subscriptions
  const {
    result: poolUpdatesResult,
    loading: poolUpdatesLoading,
    error: poolUpdatesError
  } = useSubscription<AMMPoolUpdatesSubscription>(
    AMMPoolUpdates,
    () => ({ poolAccount: selectedPoolId.value }),
    () => ({ enabled: !!selectedPoolId.value })
  )

  const {
    result: tokenPriceUpdatesResult,
    loading: tokenPriceUpdatesLoading,
    error: tokenPriceUpdatesError
  } = useSubscription<TokenPriceUpdatesSubscription>(
    TokenPriceUpdates,
    () => ({
      currency: selectedToken.currency,
      issuer: selectedToken.issuer
    }),
    () => ({ 
      enabled: !!selectedToken.currency && !!selectedToken.issuer 
    })
  )

  const {
    result: heatmapUpdatesResult,
    loading: heatmapUpdatesLoading,
    error: heatmapUpdatesError
  } = useSubscription<HeatmapUpdatesSubscription>(
    HeatmapUpdates,
    () => ({ timeRange: TimeRange.Day_1 }),
    () => ({ enabled: true })
  )

  // Computed properties
  const accountBalances = computed(() => accountBalancesResult.value?.xrpAccountBalances)
  const accountTransactions = computed(() => accountTransactionsResult.value?.xrpAccountTransactions || [])
  const accountData = computed(() => accountDataResult.value?.xrpAccountData)
  const ammPools = computed(() => ammPoolsResult.value?.xrpAMMPools || [])
  const ammPool = computed(() => ammPoolResult.value?.xrpAMMPool)
  const swapQuote = computed(() => swapQuoteResult.value?.xrpAMMSwapQuote)
  const liquidityValue = computed(() => liquidityValueResult.value?.xrpAMMLiquidityValue)
  const ammTransactions = computed(() => ammTransactionsResult.value?.xrpAMMTransactions || [])
  const heatmapData = computed(() => heatmapResult.value?.ammHeatmap || [])
  const topPools = computed(() => topPoolsResult.value?.topAMMPools || [])
  const poolDetails = computed(() => poolDetailsResult.value?.ammPoolDetails)
  const tokens = computed(() => tokensResult.value?.xrpTokens || [])
  const token = computed(() => tokenResult.value?.xrpToken)
  const ammPoolsForToken = computed(() => ammPoolsForTokenResult.value?.ammPools || [])

  // Methods
  const setAddress = (address: string) => {
    currentAddress.value = address
  }

  const setSelectedPool = (poolId: string) => {
    selectedPoolId.value = poolId
  }

  const setSelectedToken = (currency: string, issuer: string) => {
    selectedToken.currency = currency
    selectedToken.issuer = issuer
  }

  const refreshAccountData = async () => {
    await Promise.all([
      refetchAccountBalances(),
      refetchAccountTransactions(),
      refetchAccountData()
    ])
  }

  const refreshAMMData = async () => {
    await Promise.all([
      refetchAMMPools(),
      refetchAMMPool(),
      refetchSwapQuote(),
      refetchLiquidityValue(),
      refetchAMMTransactions()
    ])
  }

  const refreshHeatmapData = async () => {
    await Promise.all([
      refetchHeatmap(),
      refetchTopPools(),
      refetchPoolDetails()
    ])
  }

  const refreshTokenData = async () => {
    await Promise.all([
      refetchTokens(),
      refetchToken(),
      refetchAMMPoolsForToken()
    ])
  }

  // Loading states
  const isLoading = computed(() => 
    accountBalancesLoading.value ||
    accountTransactionsLoading.value ||
    accountDataLoading.value ||
    ammPoolsLoading.value ||
    ammPoolLoading.value ||
    swapQuoteLoading.value ||
    liquidityValueLoading.value ||
    ammTransactionsLoading.value ||
    heatmapLoading.value ||
    topPoolsLoading.value ||
    poolDetailsLoading.value ||
    tokensLoading.value ||
    tokenLoading.value ||
    ammPoolsForTokenLoading.value
  )

  const hasError = computed(() => 
    accountBalancesError.value ||
    accountTransactionsError.value ||
    accountDataError.value ||
    ammPoolsError.value ||
    ammPoolError.value ||
    swapQuoteError.value ||
    liquidityValueError.value ||
    ammTransactionsError.value ||
    heatmapError.value ||
    topPoolsError.value ||
    poolDetailsError.value ||
    tokensError.value ||
    tokenError.value ||
    ammPoolsForTokenError.value
  )

  return {
    // State
    currentAddress,
    selectedPoolId,
    selectedToken,
    
    // Computed
    accountBalances,
    accountTransactions,
    accountData,
    ammPools,
    ammPool,
    swapQuote,
    liquidityValue,
    ammTransactions,
    heatmapData,
    topPools,
    poolDetails,
    tokens,
    token,
    ammPoolsForToken,
    isLoading,
    hasError,
    
    // Loading states
    accountBalancesLoading,
    accountTransactionsLoading,
    accountDataLoading,
    ammPoolsLoading,
    ammPoolLoading,
    swapQuoteLoading,
    liquidityValueLoading,
    ammTransactionsLoading,
    heatmapLoading,
    topPoolsLoading,
    poolDetailsLoading,
    tokensLoading,
    tokenLoading,
    ammPoolsForTokenLoading,
    
    // Error states
    accountBalancesError,
    accountTransactionsError,
    accountDataError,
    ammPoolsError,
    ammPoolError,
    swapQuoteError,
    liquidityValueError,
    ammTransactionsError,
    heatmapError,
    topPoolsError,
    poolDetailsError,
    tokensError,
    tokenError,
    ammPoolsForTokenError,
    
    // Methods
    setAddress,
    setSelectedPool,
    setSelectedToken,
    refreshAccountData,
    refreshAMMData,
    refreshHeatmapData,
    refreshTokenData,
    
    // Subscriptions
    poolUpdatesResult,
    tokenPriceUpdatesResult,
    heatmapUpdatesResult
  }
} 