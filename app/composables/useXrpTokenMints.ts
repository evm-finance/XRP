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

export default function useXrpTokenMints() {
  // State
  const loading = ref(false)
  const tokenMints = ref<XRPTokenMint[]>([])
  const liquidityPools = ref<XRPLiquidityPool[]>([])
  const selectedTimeRange = ref<'24h' | '7d' | '30d' | 'all'>('7d')
  const selectedSortBy = ref<'tvl' | 'volume' | 'apr' | 'fees'>('tvl')
  const selectedSortOrder = ref<'asc' | 'desc'>('desc')
  
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
      
      filtered = filtered.filter(mint => new Date(mint.mintDate).getTime() > cutoff)
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
      filtered = filtered.filter(mint => mint.liquidityPools.length > 0)
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
          comparison = a.fees24H - b.fees24H
          break
      }
      return sortOrder === 'asc' ? comparison : -comparison
    })

    return filtered
  })

  // Table headers
  const tokenMintHeaders = computed(() => [
    { text: 'Token', value: 'token', width: '200' },
    { text: 'Issuer', value: 'issuer', width: '200' },
    { text: 'Mint Date', value: 'mintDate', width: '120' },
    { text: 'Initial Supply', value: 'initialSupply', width: '120' },
    { text: 'Current Supply', value: 'currentSupply', width: '120' },
    { text: 'Price', value: 'price', width: '100' },
    { text: 'Market Cap', value: 'marketcap', width: '120' },
    { text: '24H Volume', value: 'volume24H', width: '120' },
    { text: 'Holders', value: 'holders', width: '100' },
    { text: 'Liquidity Pools', value: 'liquidityPools', width: '120' },
  ])

  const liquidityPoolHeaders = computed(() => [
    { text: 'Pool', value: 'pool', width: '250' },
    { text: 'TVL', value: 'tvl', width: '120' },
    { text: '24H Volume', value: 'volume24H', width: '120' },
    { text: '24H Fees', value: 'fees24H', width: '120' },
    { text: 'APR', value: 'apr', width: '100' },
    { text: '24H Change', value: 'priceChange24H', width: '120' },
    { text: '7D Change', value: 'priceChange7D', width: '120' },
    { text: '24H Trades', value: 'transactions24H', width: '120' },
    { text: 'Unique Traders', value: 'uniqueTraders24H', width: '120' },
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
    // Mock token mints
    tokenMints.value = Array.from({ length: 20 }, (_, i) => ({
      currency: `TOKEN${i + 1}`,
      issuerAddress: `r${Math.random().toString(36).substring(2, 34).toUpperCase()}`,
      tokenName: `Token ${i + 1}`,
      issuerName: `Issuer ${i + 1}`,
      mintDate: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      initialSupply: Math.random() * 1000000,
      currentSupply: Math.random() * 1000000,
      price: Math.random() * 0.01,
      marketcap: Math.random() * 1000000,
      volume24H: Math.random() * 50000,
      holders: Math.floor(Math.random() * 1000),
      liquidityPools: Array.from({ length: Math.floor(Math.random() * 3) }, (_, j) => ({
        poolId: `pool_${i}_${j}`,
        asset1: { currency: 'XRP', symbol: 'XRP', name: 'XRP' },
        asset2: { currency: `TOKEN${i + 1}`, issuer: `r${Math.random().toString(36).substring(2, 34).toUpperCase()}`, symbol: `TOKEN${i + 1}`, name: `Token ${i + 1}` },
        liquidity: Math.random() * 100000,
        volume24H: Math.random() * 10000,
        volume7D: Math.random() * 50000,
        fees24H: Math.random() * 100,
        fees7D: Math.random() * 500,
        apr: Math.random() * 50,
        tvl: Math.random() * 500000,
        priceChange24H: (Math.random() - 0.5) * 20,
        priceChange7D: (Math.random() - 0.5) * 50,
        transactions24H: Math.floor(Math.random() * 1000),
        uniqueTraders24H: Math.floor(Math.random() * 100)
      }))
    }))

    // Mock liquidity pools
    liquidityPools.value = Array.from({ length: 30 }, (_, i) => ({
      poolId: `pool_${i}`,
      asset1: { currency: 'XRP', symbol: 'XRP', name: 'XRP' },
      asset2: { 
        currency: `TOKEN${i + 1}`, 
        issuer: `r${Math.random().toString(36).substring(2, 34).toUpperCase()}`, 
        symbol: `TOKEN${i + 1}`, 
        name: `Token ${i + 1}` 
      },
      liquidity: Math.random() * 1000000,
      volume24H: Math.random() * 100000,
      volume7D: Math.random() * 500000,
      fees24H: Math.random() * 1000,
      fees7D: Math.random() * 5000,
      apr: Math.random() * 100,
      tvl: Math.random() * 5000000,
      priceChange24H: (Math.random() - 0.5) * 20,
      priceChange7D: (Math.random() - 0.5) * 50,
      transactions24H: Math.floor(Math.random() * 1000),
      uniqueTraders24H: Math.floor(Math.random() * 100)
    }))
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
    copyToClipboard
  }
} 