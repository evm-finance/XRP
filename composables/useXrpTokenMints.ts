import { ref, computed } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'
import { XRPTokenMintsGQL, XRPLiquidityPoolsGQL } from '~/apollo/queries'

// Interfaces
export interface XRPTokenMint {
  currency: string
  issuerAddress: string
  tokenName: string
  issuerName: string
  mintDate: string
  initialSupply: number
  currentSupply: number
  price: number
  marketcap: number
  volume24H: number
  holders: number
  liquidityPools: XRPLiquidityPool[]
}

export interface XRPLiquidityPool {
  poolId: string
  asset1: XRPTokenAsset
  asset2: XRPTokenAsset
  liquidity: number
  volume24H: number
  volume7D: number
  fees24H: number
  fees7D: number
  apr: number
  tvl: number
  priceChange24H: number
  priceChange7D: number
  transactions24H: number
  uniqueTraders24H: number
}

export interface XRPTokenAsset {
  currency: string
  issuer?: string
  symbol: string
  name: string
}

export interface TokenMintFilters {
  timeRange: '24h' | '7d' | '30d' | 'all'
  minMarketcap: number
  minVolume: number
  hasLiquidity: boolean
}

export interface LiquidityPoolFilters {
  minTvl: number
  minVolume: number
  minApr: number
  sortBy: 'tvl' | 'volume' | 'apr' | 'fees'
  sortOrder: 'asc' | 'desc'
}

interface TokenMint {
  id: string
  name: string
  symbol: string
  issuer: string
  currency: string
  marketcap: number
  volume24H: number
  price: number
  change24H: number
  createdAt: string
}

interface LiquidityPool {
  id: string
  name: string
  token1: string
  token2: string
  tvl: number
  volume24H: number
  apr: number
  fees: number
  priceChange24H: number
  trades24H: number
  uniqueTraders: number
}

export default function useXrpTokenMints() {
  const { useLoggedQuery } = useXrpGraphQLWithLogging()
  
  // State
  const loading = ref(false)
  const tokenMints = ref<TokenMint[]>([])
  const liquidityPools = ref<LiquidityPool[]>([])
  const selectedTimeRange = ref<'24h' | '7d' | '30d' | 'all'>('7d')
  const selectedSortBy = ref<'tvl' | 'volume' | 'apr' | 'fees'>('tvl')
  const selectedSortOrder = ref<'asc' | 'desc'>('desc')
  const error = ref<string | null>(null)
  
  // Filters
  const tokenMintFilters = ref<TokenMintFilters>({
    timeRange: '7d',
    minMarketcap: 0,
    minVolume: 0,
    hasLiquidity: false
  })
  
  const liquidityPoolFilters = ref<LiquidityPoolFilters>({
    minTvl: 0,
    minVolume: 0,
    minApr: 0,
    sortBy: 'tvl',
    sortOrder: 'desc'
  })

  // Log the query content BEFORE making the calls
  console.log('🚀 [BEFORE QUERY] useXrpTokenMints - XRPTokenMintsGQL:', {
    query: XRPTokenMintsGQL.loc?.source.body,
    variables: {},
    timestamp: new Date().toISOString()
  })

  // Queries
  const { onResult: onTokenMintsResult, refetch: refetchTokenMints } = useLoggedQuery(XRPTokenMintsGQL, {
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
    context: {
      queryName: 'XRPTokenMints',
      component: 'useXrpTokenMints',
      purpose: 'XRP token mints data'
    }
  })

  console.log('🚀 [BEFORE QUERY] useXrpTokenMints - XRPLiquidityPoolsGQL:', {
    query: XRPLiquidityPoolsGQL.loc?.source.body,
    variables: {},
    timestamp: new Date().toISOString()
  })

  const { onResult: onLiquidityPoolsResult, refetch: refetchLiquidityPools } = useLoggedQuery(XRPLiquidityPoolsGQL, {
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
    context: {
      queryName: 'XRPLiquidityPools',
      component: 'useXrpTokenMints',
      purpose: 'XRP liquidity pools data'
    }
  })

  // Computed
  const filteredTokenMints = computed(() => {
    let filtered = tokenMints.value

    // Filter by time range
    if (tokenMintFilters.value.timeRange !== 'all') {
      const now = new Date()
      const timeRanges = {
        '24h': 24 * 60 * 60 * 1000,
        '7d': 7 * 24 * 60 * 60 * 1000,
        '30d': 30 * 24 * 60 * 60 * 1000
      }
      const cutoff = now.getTime() - timeRanges[tokenMintFilters.value.timeRange]
      
      filtered = filtered.filter(mint => new Date(mint.createdAt).getTime() > cutoff)
    }

    // Filter by minimum market cap
    if (tokenMintFilters.value.minMarketcap > 0) {
      filtered = filtered.filter(mint => mint.marketcap >= tokenMintFilters.value.minMarketcap)
    }

    // Filter by minimum volume
    if (tokenMintFilters.value.minVolume > 0) {
      filtered = filtered.filter(mint => mint.volume24H >= tokenMintFilters.value.minVolume)
    }

    // Filter by liquidity
    if (tokenMintFilters.value.hasLiquidity) {
      filtered = filtered.filter(mint => liquidityPools.value.some(pool => pool.token1 === mint.currency || pool.token2 === mint.currency))
    }

    return filtered
  })

  const filteredLiquidityPools = computed(() => {
    let filtered = liquidityPools.value

    // Filter by minimum TVL
    if (liquidityPoolFilters.value.minTvl > 0) {
      filtered = filtered.filter(pool => pool.tvl >= liquidityPoolFilters.value.minTvl)
    }

    // Filter by minimum volume
    if (liquidityPoolFilters.value.minVolume > 0) {
      filtered = filtered.filter(pool => pool.volume24H >= liquidityPoolFilters.value.minVolume)
    }

    // Filter by minimum APR
    if (liquidityPoolFilters.value.minApr > 0) {
      filtered = filtered.filter(pool => pool.apr >= liquidityPoolFilters.value.minApr)
    }

    // Sort pools
    const sortBy = liquidityPoolFilters.value.sortBy
    const sortOrder = liquidityPoolFilters.value.sortOrder
    
    filtered.sort((a, b) => {
      let comparison = 0
      switch (sortBy) {
        case 'tvl':
          comparison = a.tvl - b.tvl
          break
        case 'volume':
          comparison = a.volume24H - b.volume24H
          break
        case 'apr':
          comparison = a.apr - b.apr
          break
        case 'fees':
          comparison = a.fees - b.fees
          break
      }
      return sortOrder === 'asc' ? comparison : -comparison
    })

    return filtered
  })

  // Table headers
  const tokenMintHeaders = computed(() => [
    { text: 'Token', value: 'name', width: '200' },
    { text: 'Issuer', value: 'issuer', width: '200' },
    { text: 'Mint Date', value: 'createdAt', width: '120' },
    { text: 'Initial Supply', value: 'initialSupply', width: '120' },
    { text: 'Current Supply', value: 'currentSupply', width: '120' },
    { text: 'Price', value: 'price', width: '100' },
    { text: 'Market Cap', value: 'marketcap', width: '120' },
    { text: '24H Volume', value: 'volume24H', width: '120' },
    { text: 'Liquidity Pools', value: 'liquidityPools', width: '120' },
  ])

  const liquidityPoolHeaders = computed(() => [
    { text: 'Pool', value: 'name', width: '250' },
    { text: 'TVL', value: 'tvl', width: '120' },
    { text: '24H Volume', value: 'volume24H', width: '120' },
    { text: 'APR', value: 'apr', width: '100' },
    { text: '24H Change', value: 'priceChange24H', width: '120' },
    { text: '7D Change', value: 'priceChange7D', width: '120' },
    { text: '24H Trades', value: 'trades24H', width: '120' },
    { text: 'Unique Traders', value: 'uniqueTraders', width: '120' },
  ])

  // Methods
  const formatAddress = (address: string) => {
    return `${address.slice(0, 10)}...${address.slice(-10)}`
  }

  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 6,
    }).format(amount)
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount)
  }

  const formatPercentage = (value: number) => {
    return `${value > 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      // Could add toast notification here
    } catch (err) {
      console.error('Failed to copy text: ', err)
    }
  }

  // Data fetching functions
  const fetchData = async () => {
    loading.value = true
    error.value = null
    try {
      await refetchTokenMints()
      await refetchLiquidityPools()
    } catch (err) {
      error.value = 'Failed to fetch data'
      console.error('Error fetching data:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchTokenMints = async () => {
    try {
      await refetchTokenMints()
    } catch (err) {
      error.value = 'Failed to fetch token mints'
      console.error('Error fetching token mints:', err)
    }
  }

  const fetchLiquidityPools = async () => {
    try {
      await refetchLiquidityPools()
    } catch (err) {
      error.value = 'Failed to fetch liquidity pools'
      console.error('Error fetching liquidity pools:', err)
    }
  }

  // Event handlers
  onTokenMintsResult((queryResult: any) => {
    const mints = queryResult.data?.xrpTokenMints ?? []
    
    // Transform GraphQL data to component format
    const transformedMints: TokenMint[] = mints.map((mint: any) => ({
      id: mint.currency || '',
      name: mint.tokenName || mint.currency || '',
      symbol: mint.currency || '',
      issuer: mint.issuerAddress || '',
      currency: mint.currency || '',
      marketcap: mint.marketcap || 0,
      volume24H: mint.volume24H || 0,
      price: mint.price || 0,
      change24H: mint.change24H || 0,
      createdAt: mint.mintDate || new Date().toISOString()
    }))
    
    tokenMints.value = transformedMints
    loading.value = false
  })

  onLiquidityPoolsResult((queryResult: any) => {
    const pools = queryResult.data?.xrpLiquidityPools ?? []
    
    // Transform GraphQL data to component format
    const transformedPools: LiquidityPool[] = pools.map((pool: any) => ({
      id: pool.poolId || '',
      name: `${pool.asset1?.symbol || ''}/${pool.asset2?.symbol || ''} Pool`,
      token1: pool.asset1?.symbol || '',
      token2: pool.asset2?.symbol || '',
      tvl: pool.tvl || 0,
      volume24H: pool.volume24H || 0,
      apr: pool.apr || 0,
      fees: pool.fees24H || 0,
      priceChange24H: pool.priceChange24H || 0,
      trades24H: pool.transactions24H || 0,
      uniqueTraders: pool.uniqueTraders24H || 0
    }))
    
    liquidityPools.value = transformedPools
  })

  // Error handling queries with enhanced logging
  const { onError: onTokenMintsError } = useLoggedQuery(
    XRPTokenMintsGQL,
    {
      fetchPolicy: 'no-cache',
      pollInterval: 60000,
      context: {
        queryName: 'XRPTokenMints',
        component: 'useXrpTokenMints',
        purpose: 'XRP token mints data (error handling)'
      }
    }
  )

  const { onError: onLiquidityPoolsError } = useLoggedQuery(
    XRPLiquidityPoolsGQL,
    {
      fetchPolicy: 'no-cache',
      pollInterval: 60000,
      context: {
        queryName: 'XRPLiquidityPools',
        component: 'useXrpTokenMints',
        purpose: 'XRP liquidity pools data (error handling)'
      }
    }
  )

  onTokenMintsError((error: any) => {
    console.error('GraphQL Error in useXrpTokenMints token mints:', error)
    loading.value = false
    error.value = 'Failed to fetch token mints'
  })

  onLiquidityPoolsError((error: any) => {
    console.error('GraphQL Error in useXrpTokenMints liquidity pools:', error)
    error.value = 'Failed to fetch liquidity pools'
  })

  // Remove mock data generation - rely entirely on live data
  // generateMockData() - REMOVED
  loading.value = false

  return {
    // State
    loading,
    tokenMints,
    liquidityPools,
    selectedTimeRange,
    selectedSortBy,
    selectedSortOrder,
    tokenMintFilters,
    liquidityPoolFilters,
    error,
    
    // Computed
    filteredTokenMints,
    filteredLiquidityPools,
    tokenMintHeaders,
    liquidityPoolHeaders,
    
    // Methods
    formatAddress,
    formatAmount,
    formatCurrency,
    formatPercentage,
    formatDate,
    copyToClipboard,
    fetchData,
    fetchTokenMints,
    fetchLiquidityPools,
  }
} 