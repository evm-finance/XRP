import { ref, computed, Ref, reactive } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'

import {
  XRPAccountBalancesGQL, 
  XRPAccountTransactionsGQL, 
  XRPTransactionGQL,
  XRPScreenerGQL,
  XRPAmmPoolsGQL
} from '~/apollo/queries'

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
  const { useLoggedQuery } = useXrpGraphQLWithLogging()
  
  // Reactive state
  const currentAddress = ref('')
  const selectedPoolId = ref('')
  const selectedToken = reactive<{ currency: string; issuer: string }>({ currency: '', issuer: '' })

  // Log the query content BEFORE making the calls
  console.log('🚀 [BEFORE QUERY] useXrpGraphQL - XRPAccountBalancesGQL:', {
    query: XRPAccountBalancesGQL.loc?.source.body,
    variables: { account: currentAddress.value },
    timestamp: new Date().toISOString()
  })

  // Account Queries with Enhanced Logging
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

  console.log('🚀 [BEFORE QUERY] useXrpGraphQL - XRPAccountTransactionsGQL:', {
    query: XRPAccountTransactionsGQL.loc?.source.body,
    variables: { address: currentAddress.value },
    timestamp: new Date().toISOString()
  })

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

  console.log('🚀 [BEFORE QUERY] useXrpGraphQL - XRPScreenerGQL:', {
    query: XRPScreenerGQL.loc?.source.body,
    variables: {},
    timestamp: new Date().toISOString()
  })

  // Screener Query with Enhanced Logging
  const {
    result: screenerResult,
    loading: screenerLoading,
    error: screenerError,
    refetch: refetchScreener
  } = useLoggedQuery(
    XRPScreenerGQL,
    {},
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  console.log('🚀 [BEFORE QUERY] useXrpGraphQL - XRPAmmPoolsGQL:', {
    query: XRPAmmPoolsGQL.loc?.source.body,
    variables: {},
    timestamp: new Date().toISOString()
  })

  // AMM Pools Query with Enhanced Logging
  const {
    result: ammPoolsResult,
    loading: ammPoolsLoading,
    error: ammPoolsError,
    refetch: refetchAmmPools
  } = useLoggedQuery(
    XRPAmmPoolsGQL,
    {},
    () => ({ fetchPolicy: 'cache-and-network' })
  )

  // Computed properties
  const accountData = computed(() => accountBalancesResult.value?.xrpAccountBalances || null)
  const accountTransactions = computed(() => accountTransactionsResult.value?.xrpTransactions || [])
  const screenerData = computed(() => screenerResult.value?.xrpScreener || [])
  const ammPools = computed(() => ammPoolsResult.value?.xrpAmmPools || [])

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
    if (currentAddress.value) {
      await refetchAccountBalances()
      await refetchAccountTransactions()
    }
  }

  const refreshAMMData = async () => {
    await refetchAmmPools()
  }

  const refreshHeatmapData = async () => {
    await refetchScreener()
  }

  const refreshTokenData = async () => {
    // Token data refresh logic
  }

  return {
    // State
    currentAddress,
    selectedPoolId,
    selectedToken,
    
    // Data
    accountData,
    accountTransactions,
    screenerData,
    ammPools,
    
    // Loading states
    accountBalancesLoading,
    accountTransactionsLoading,
    screenerLoading,
    ammPoolsLoading,
    
    // Error states
    accountBalancesError,
    accountTransactionsError,
    screenerError,
    ammPoolsError,
    
    // Methods
    setAddress,
    setSelectedPool,
    setSelectedToken,
    refreshAccountData,
    refreshAMMData,
    refreshHeatmapData,
    refreshTokenData,
    
    // Query methods for external use
    getAccountBalances: (address: string) => useLoggedQuery(
      XRPAccountBalancesGQL,
      { account: address },
      () => ({ enabled: !!address, fetchPolicy: 'cache-and-network' })
    ),
    getAccountTransactions: (address: string) => useLoggedQuery(
      XRPAccountTransactionsGQL,
      { address },
      () => ({ enabled: !!address, fetchPolicy: 'cache-and-network' })
    ),
    getTransaction: (hash: string) => useLoggedQuery(
      XRPTransactionGQL,
      { hash },
      () => ({ enabled: !!hash, fetchPolicy: 'cache-and-network' })
    ),
    getScreenerData: () => useLoggedQuery(
      XRPScreenerGQL,
      {},
      () => ({ fetchPolicy: 'cache-and-network' })
    ),
    getAmmPools: () => useLoggedQuery(
      XRPAmmPoolsGQL,
      {},
      () => ({ fetchPolicy: 'cache-and-network' })
    )
  }
} 