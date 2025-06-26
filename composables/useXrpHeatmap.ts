import { computed, onBeforeUnmount, onMounted, Ref, ref, useContext, useStore, watch } from '@nuxtjs/composition-api'
import { useXrpConfigs, XrpTimeFrame, XrpBlockSize } from './useXrpConfigs'
import { useQuery } from '@vue/apollo-composable/dist'
import { GetAllAMMLiquidityValuesGQL, TestGraphQLGQL, SimpleAMMLiquidityValuesGQL } from '~/apollo/queries'
import { State } from '~/types/state'

export interface XrpAmmHeatmapData {
  poolId: string
  asset1: {
    currency: string
    issuer: string
  }
  asset2: {
    currency: string
    issuer: string
  }
  asset1Balance: number
  asset2Balance: number
  asset1ValueUsd: number
  asset2ValueUsd: number
  totalLiquidityUsd: number
  fee: number
  createdAt: string
  // Additional fields for display
  color: string
}

const setColor = (liquidity: number, maxLiquidity: number) => {
  const percentage = (liquidity / maxLiquidity) * 100
  
  if (percentage >= 80) {
    return '#2f6a32' // Dark green for highest liquidity
  } else if (percentage >= 60) {
    return '#3e8e42' // Medium green
  } else if (percentage >= 40) {
    return '#4eb153' // Light green
  } else if (percentage >= 20) {
    return '#71c175' // Very light green
  } else {
    return '#a5d6a7' // Lightest green for lowest liquidity
  }
}

export function useXrpHeatmap() {
  const { state } = useStore<State>()
  
  const {
    blockSize,
  } = useXrpConfigs()

  // Test GraphQL endpoint first
  const { result: testResult, error: testError } = useQuery(
    TestGraphQLGQL,
    {},
    {
      errorPolicy: 'all',
      fetchPolicy: 'network-only',
    }
  )

  // Log test query results
  watch(testError, (newError) => {
    if (newError) {
      console.error('GraphQL Test Query Error:', {
        message: newError.message,
        graphQLErrors: newError.graphQLErrors,
        networkError: newError.networkError,
      })
    } else if (testResult.value) {
      console.log('GraphQL Test Query Success:', testResult.value)
    }
  })

  // GraphQL query for AMM liquidity data - try simple version first
  const { result: ammData, loading, error, refetch } = useQuery(
    SimpleAMMLiquidityValuesGQL,
    {},
    {
      errorPolicy: 'all',
      fetchPolicy: 'network-only', // Force network request
    }
  )

  // Log detailed error information
  watch(error, (newError) => {
    if (newError) {
      console.error('GraphQL Error Details:', {
        message: newError.message,
        graphQLErrors: newError.graphQLErrors,
        networkError: newError.networkError,
        extraInfo: newError.extraInfo
      })
    }
  })

  // Check if GraphQL endpoint is available
  const isGraphQLEndpointAvailable = computed(() => {
    return process.env.BASE_GRAPHQL_SERVER_URL && process.env.BASE_GRAPHQL_SERVER_URL !== ''
  })

  // Mock data for when GraphQL endpoint fails
  const mockAmmData = computed(() => ({
    xrpAMMLiquidityValues: [
      {
        poolId: '1',
        asset1ValueUsd: 5000000,
        asset2ValueUsd: 5000000,
        totalLiquidityUsd: 10000000,
        asset1: { currency: 'XRP', issuer: null },
        asset2: { currency: 'USD', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
        fee: 0.3,
        createdAt: new Date().toISOString()
      },
      {
        poolId: '2',
        asset1ValueUsd: 3000000,
        asset2ValueUsd: 7000000,
        totalLiquidityUsd: 10000000,
        asset1: { currency: 'XRP', issuer: null },
        asset2: { currency: 'BTC', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
        fee: 0.5,
        createdAt: new Date().toISOString()
      },
      {
        poolId: '3',
        asset1ValueUsd: 2000000,
        asset2ValueUsd: 8000000,
        totalLiquidityUsd: 10000000,
        asset1: { currency: 'XRP', issuer: null },
        asset2: { currency: 'ETH', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
        fee: 0.4,
        createdAt: new Date().toISOString()
      },
      {
        poolId: '4',
        asset1ValueUsd: 1500000,
        asset2ValueUsd: 3500000,
        totalLiquidityUsd: 5000000,
        asset1: { currency: 'XRP', issuer: null },
        asset2: { currency: 'SOL', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
        fee: 0.6,
        createdAt: new Date().toISOString()
      },
      {
        poolId: '5',
        asset1ValueUsd: 800000,
        asset2ValueUsd: 1200000,
        totalLiquidityUsd: 2000000,
        asset1: { currency: 'XRP', issuer: null },
        asset2: { currency: 'ADA', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
        fee: 0.7,
        createdAt: new Date().toISOString()
      }
    ]
  }))

  const heatmapData = computed(() => {
    // Use mock data if GraphQL fails or returns no data
    const data = hasError.value || !ammData.value?.xrpAMMLiquidityValues ? mockAmmData.value : ammData.value
    
    if (!data?.xrpAMMLiquidityValues) return []
    
    const pools = data.xrpAMMLiquidityValues as XrpAmmHeatmapData[]
    
    // Sort pools based on the selected block size
    const sortedPools = [...pools].sort((a, b) => {
      let aValue: number
      let bValue: number
      
      switch (blockSize.value) {
        case 'totalLiquidityUsd':
          aValue = a.totalLiquidityUsd
          bValue = b.totalLiquidityUsd
          break
        case 'asset1ValueUsd':
          aValue = a.asset1ValueUsd
          bValue = b.asset1ValueUsd
          break
        case 'asset2ValueUsd':
          aValue = a.asset2ValueUsd
          bValue = b.asset2ValueUsd
          break
        default:
          aValue = a.totalLiquidityUsd
          bValue = b.totalLiquidityUsd
      }
      
      return bValue - aValue // Always sort descending
    })

    // Calculate max value for color scaling based on selected block size
    let maxValue = 1
    switch (blockSize.value) {
      case 'totalLiquidityUsd':
        maxValue = Math.max(...sortedPools.map(pool => pool.totalLiquidityUsd), 1)
        break
      case 'asset1ValueUsd':
        maxValue = Math.max(...sortedPools.map(pool => pool.asset1ValueUsd), 1)
        break
      case 'asset2ValueUsd':
        maxValue = Math.max(...sortedPools.map(pool => pool.asset2ValueUsd), 1)
        break
      default:
        maxValue = Math.max(...sortedPools.map(pool => pool.totalLiquidityUsd), 1)
    }

    // Add colors based on the selected value
    return sortedPools.map(pool => {
      let valueForColor: number
      switch (blockSize.value) {
        case 'totalLiquidityUsd':
          valueForColor = pool.totalLiquidityUsd
          break
        case 'asset1ValueUsd':
          valueForColor = pool.asset1ValueUsd
          break
        case 'asset2ValueUsd':
          valueForColor = pool.asset2ValueUsd
          break
        default:
          valueForColor = pool.totalLiquidityUsd
      }
      
      return {
        ...pool,
        color: setColor(valueForColor, maxValue)
      }
    })
  })

  const processedData = computed(() => {
    const data = heatmapData.value
    if (!data.length) return []

    // Group data by block size for grid layout
    const groupedData: XrpAmmHeatmapData[][] = []
    const itemsPerRow = 6 // Fixed number of items per row

    for (let i = 0; i < data.length; i += itemsPerRow) {
      groupedData.push(data.slice(i, i + itemsPerRow))
    }

    return groupedData
  })

  const isLoading = computed(() => loading.value)
  const hasError = computed(() => !!error.value)

  const refreshData = () => {
    refetch()
  }

  return {
    heatmapData,
    processedData,
    isLoading,
    hasError,
    refreshData,
    blockSize,
    isGraphQLEndpointAvailable,
  }
} 