import { computed, onBeforeUnmount, onMounted, ref, useContext, watch } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'
import { XRPAMMHeatmapGQL } from '~/apollo/queries'
import type { HeatmapData, HeatmapRowData, HeatmapUpdateData } from '~/types/heatmap'
import type { HeatmapIntervals } from '~/types/state'
import useHeatmapConfigs from '~/composables/useHeatmapConfigs'

// Color function for XRP AMM heatmap tiles based on liquidity
const setXrpAmmColor = (liquidity: number, maxLiquidity: number, blueTile: boolean) => {
  const percentage = liquidity / maxLiquidity
  
  if (percentage > 0.8) {
    return blueTile ? '#5898f8' : '#2f6a32'
  } else if (percentage > 0.6) {
    return blueTile ? '#336cc2' : '#3e8e42'
  } else if (percentage > 0.4) {
    return blueTile ? '#234e91' : '#4eb153'
  } else if (percentage > 0.2) {
    return blueTile ? '#1A2A52' : '#71c175'
  } else {
    return blueTile ? '#0f1a2a' : '#4a7c59'
  }
}

export default function useXrpAmmHeatmap(userCanAccessTrend: any) {
  const { $f } = useContext()
  const { useLoggedQuery } = useXrpGraphQLWithLogging()
  
  // State
  const loading = ref(true)
  const rowData = ref<any[]>([])
  const heatmapData = ref<HeatmapData[]>([])
  
  // Configuration
  const {
    timeFrame,
    blockSize,
    blueTile,
    displayFavorites,
    numOfCoins,
    displayGainersAndLosers,
  } = useHeatmapConfigs()

  // Log the query content BEFORE making the call
  console.log('🚀 [BEFORE QUERY] useXrpAmmHeatmap - XRPAMMHeatmapGQL:', {
    query: XRPAMMHeatmapGQL.loc?.source.body,
    variables: {
      timeFrame: timeFrame.value,
      blockSize: blockSize.value,
      limit: numOfCoins.value
    },
    timestamp: new Date().toISOString()
  })

  // Live AMM heatmap data
  const { result: ammResult, loading: ammLoading, error, refetch } = useLoggedQuery(
    XRPAMMHeatmapGQL,
    computed(() => ({
      timeFrame: timeFrame.value,
      blockSize: blockSize.value,
      limit: numOfCoins.value
    }))
  )

  // Internal error handling - prevent Apollo errors from bubbling up
  const internalError = ref<string | null>(null)
  
  // Watch for Apollo errors and handle them internally
  watch(error, (newError) => {
    if (newError) {
      console.warn('GraphQL error in useXrpAmmHeatmap, using fallback data:', newError)
      internalError.value = 'Network error: using fallback data'
    } else {
      internalError.value = null
    }
  })

  // Error state for graceful fallback
  const hasError = computed(() => !!internalError.value || !!error.value)

  // Transform AMM data to heatmap format
  const transformedData = computed<HeatmapRowData[]>(() => {
    const data = ammResult.value?.xrpAmmHeatmap || []
    if (!data.length) return []

    const maxLiquidity = Math.max(...data.map((item: any) => item.liquidity))

    return data.map((item: any): HeatmapRowData => {
      const pairName = `${item.token1.symbol}/${item.token2.symbol}`
      return {
        qc_key: item.poolId,
        name: pairName,
        value: item.liquidity,
        symbol: pairName,
        token1: item.token1,
        token2: item.token2,
        liquidity: item.liquidity,
        volume24h: item.volume24h,
        fee: item.fee,
        apr: item.apr,
        priceChange24h: item.priceChange24h,
        color: setXrpAmmColor(item.liquidity, maxLiquidity, blueTile.value),
        size: blockSize.value === 'marketcap' ? item.liquidity : item.volume24h
      }
    }).sort((a: HeatmapRowData, b: HeatmapRowData) => (b?.size ?? 0) - (a?.size ?? 0))
  })

  // Tile text template
  const tileText = computed(() => {
    return `{name}\n{liquidity}`
  })

  // Tile tooltip template
  const tileTooltip = computed(() => {
    return (item: any) => {
      const name = item.name || 'Unknown'
      const liquidity = item.liquidity || 0
      const volume24h = item.volume24h || 0
      const apr = item.apr || 0
      const fee = item.fee || 0
      const priceChange24h = item.priceChange24h || 0
      
      return `{name}
Liquidity: $${$f(liquidity, { minDigits: 0, after: '' })}
24h Volume: $${$f(volume24h, { minDigits: 0, after: '' })}
APR: {apr:numberFormat.1}%
Fee: {fee:numberFormat.3}%
24h Change: {priceChange24h:numberFormat.2}%`
    }
  })

  // Group data by liquidity tiers
  const liquidityGroups = computed<HeatmapData[]>(() => {
    if (displayGainersAndLosers.value) {
      const resp: HeatmapData[] = [
        { name: 'High Liquidity', value: 0, children: [] },
        { name: 'Low Liquidity', value: 0, children: [] },
      ]
      const data = transformedData.value || []
      const avgLiquidity = data.length ? data.reduce((sum, item) => sum + (item.liquidity || 0), 0) / data.length : 0
      data.forEach((elem) => {
        if (!elem || typeof elem.liquidity !== 'number') return
        elem.liquidity > avgLiquidity ? resp[0].children!.push(elem) : resp[1].children!.push(elem)
      })
      return resp
    } else {
      return [{ name: '', value: 0, children: transformedData.value || [] }]
    }
  })

  // Update data for real-time updates
  const updateData = computed<HeatmapUpdateData>(() => {
    const result: HeatmapUpdateData = { name: '', value: 0 }
    transformedData.value.forEach((item) => {
      if (item.qc_key) {
        result[item.qc_key] = item
      }
    })
    return result
  })

  // Watch for configuration changes
  watch([tileText, tileTooltip, displayGainersAndLosers], () => {
    heatmapData.value = liquidityGroups.value
  })

  watch([numOfCoins, timeFrame, blockSize, blueTile], () => {
    heatmapData.value = liquidityGroups.value
  })

  // Initialize data
  const initializeData = () => {
    heatmapData.value = liquidityGroups.value
    loading.value = false
  }

  onMounted(() => {
    // Initialize with mock AMM data
    initializeData()
  })

  return {
    loading,
    heatmapData,
    tileText,
    tileTooltip,
    updateData,
    transformedData,
    liquidityGroups,
    ammResult,
    ammLoading,
    hasError,
    refetch,
  }
} 