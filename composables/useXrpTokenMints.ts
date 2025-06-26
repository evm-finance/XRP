import { ref, computed } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable'
import { gql } from 'graphql-tag'

// GraphQL queries
const XRPTokenMintsGQL = gql`
  query XRPTokenMintsGQL($limit: Int = 50) {
    xrpTokenMints(limit: $limit) {
      currency
      issuerAddress
      tokenName
      issuerName
      mintDate
      initialSupply
      currentSupply
      price
      marketcap
      volume24H
      holders
      liquidityPools {
        poolId
        asset1
        asset2
        liquidity
        volume24H
        fees24H
        apr
      }
    }
  }
`

const XRPLiquidityPoolsGQL = gql`
  query XRPLiquidityPoolsGQL($limit: Int = 50) {
    xrpLiquidityPools(limit: $limit) {
      poolId
      asset1 {
        currency
        issuer
        symbol
        name
      }
      asset2 {
        currency
        issuer
        symbol
        name
      }
      liquidity
      volume24H
      volume7D
      fees24H
      fees7D
      apr
      tvl
      priceChange24H
      priceChange7D
      transactions24H
      uniqueTraders24H
    }
  }
`

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

  // Queries
  const { onResult: onTokenMintsResult } = useQuery(XRPTokenMintsGQL, { 
    fetchPolicy: 'no-cache', 
    pollInterval: 60000 
  })
  
  const { onResult: onLiquidityPoolsResult } = useQuery(XRPLiquidityPoolsGQL, { 
    fetchPolicy: 'no-cache', 
    pollInterval: 30000 
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

  // Mock data generation for development
  const generateMockData = () => {
    const mockTokenMints: TokenMint[] = [
      {
        id: '1',
        name: 'Mock Token 1',
        symbol: 'MT1',
        issuer: 'rMockIssuer1Address123456789',
        currency: 'MT1',
        marketcap: 1000000,
        volume24H: 50000,
        price: 0.1,
        change24H: 5.2,
        createdAt: new Date().toISOString()
      },
      {
        id: '2',
        name: 'Mock Token 2',
        symbol: 'MT2',
        issuer: 'rMockIssuer2Address123456789',
        currency: 'MT2',
        marketcap: 2500000,
        volume24H: 75000,
        price: 0.25,
        change24H: -2.1,
        createdAt: new Date().toISOString()
      }
    ]

    const mockLiquidityPools: LiquidityPool[] = [
      {
        id: '1',
        name: 'MT1/XRP Pool',
        token1: 'MT1',
        token2: 'XRP',
        tvl: 500000,
        volume24H: 25000,
        apr: 12.5,
        fees: 0.3,
        priceChange24H: 1.2,
        trades24H: 150,
        uniqueTraders: 45
      },
      {
        id: '2',
        name: 'MT2/XRP Pool',
        token1: 'MT2',
        token2: 'XRP',
        tvl: 750000,
        volume24H: 35000,
        apr: 15.2,
        fees: 0.3,
        priceChange24H: -0.8,
        trades24H: 200,
        uniqueTraders: 60
      }
    ]

    return { mockTokenMints, mockLiquidityPools }
  }

  const fetchTokenMints = async () => {
    loading.value = true
    error.value = null
    
    try {
      // TODO: Replace with actual GraphQL query
      // const { data } = await $apollo.query({
      //   query: require('~/apollo/main/token-mints.query.graphql')
      // })
      
      // For now, use mock data
      const { mockTokenMints } = generateMockData()
      tokenMints.value = mockTokenMints
    } catch (err) {
      error.value = 'Failed to fetch token mints'
      console.error('Error fetching token mints:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchLiquidityPools = async () => {
    loading.value = true
    error.value = null
    
    try {
      // TODO: Replace with actual GraphQL query
      // const { data } = await $apollo.query({
      //   query: require('~/apollo/main/liquidity-pools.query.graphql')
      // })
      
      // For now, use mock data
      const { mockLiquidityPools } = generateMockData()
      liquidityPools.value = mockLiquidityPools
    } catch (err) {
      error.value = 'Failed to fetch liquidity pools'
      console.error('Error fetching liquidity pools:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchData = async () => {
    await Promise.all([fetchTokenMints(), fetchLiquidityPools()])
  }

  // Event handlers
  onTokenMintsResult((queryResult: any) => {
    const mints = queryResult.data?.xrpTokenMints ?? []
    tokenMints.value = mints
    loading.value = false
  })

  onLiquidityPoolsResult((queryResult: any) => {
    const pools = queryResult.data?.xrpLiquidityPools ?? []
    liquidityPools.value = pools
  })

  // Initialize with mock data for development
  generateMockData()
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
    fetchLiquidityPools
  }
} 