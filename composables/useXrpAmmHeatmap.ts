import { computed, onBeforeUnmount, onMounted, ref, useContext, watch } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable'
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

  // Mock AMM data - in real implementation this would come from XRPL AMM API
  const ammData = ref([
    {
      id: 'amm_1',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'USDC', name: 'USD Coin', icon: '🪙', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 2500000,
      volume24h: 1800000,
      fee: 0.003,
      apr: 12.5,
      priceChange24h: 2.3,
    },
    {
      id: 'amm_2',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'USDT', name: 'Tether', icon: '💎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 1800000,
      volume24h: 1200000,
      fee: 0.003,
      apr: 15.2,
      priceChange24h: -1.8,
    },
    {
      id: 'amm_3',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'BTC', name: 'Bitcoin', icon: '₿', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 1200000,
      volume24h: 800000,
      fee: 0.003,
      apr: 18.7,
      priceChange24h: 4.1,
    },
    {
      id: 'amm_4',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'ETH', name: 'Ethereum', icon: 'Ξ', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 900000,
      volume24h: 600000,
      fee: 0.003,
      apr: 22.1,
      priceChange24h: -0.5,
    },
    {
      id: 'amm_5',
      token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
      token2: { symbol: 'USDC', name: 'USD Coin', icon: '🪙', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 800000,
      volume24h: 400000,
      fee: 0.001,
      apr: 8.3,
      priceChange24h: 0.1,
    },
    {
      id: 'amm_6',
      token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
      token2: { symbol: 'USDT', name: 'Tether', icon: '💎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 600000,
      volume24h: 300000,
      fee: 0.001,
      apr: 9.8,
      priceChange24h: -0.2,
    },
    {
      id: 'amm_7',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'SOL', name: 'Solana', icon: '◎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 500000,
      volume24h: 250000,
      fee: 0.003,
      apr: 25.4,
      priceChange24h: 6.7,
    },
    {
      id: 'amm_8',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'ADA', name: 'Cardano', icon: '₳', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 400000,
      volume24h: 200000,
      fee: 0.003,
      apr: 28.9,
      priceChange24h: -2.1,
    },
    {
      id: 'amm_9',
      token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
      token2: { symbol: 'BTC', name: 'Bitcoin', icon: '₿', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 350000,
      volume24h: 150000,
      fee: 0.001,
      apr: 12.6,
      priceChange24h: 1.8,
    },
    {
      id: 'amm_10',
      token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
      token2: { symbol: 'DOT', name: 'Polkadot', icon: '●', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      liquidity: 300000,
      volume24h: 120000,
      fee: 0.003,
      apr: 32.1,
      priceChange24h: 3.4,
    },
  ])

  // Transform AMM data to heatmap format
  const transformedData = computed<HeatmapRowData[]>(() => {
    if (!ammData.value.length) return []

    const maxLiquidity = Math.max(...ammData.value.map(item => item.liquidity))

    return ammData.value
      .slice(0, numOfCoins.value)
      .map((item) => {
        const pairName = `${item.token1.symbol}/${item.token2.symbol}`
        const token2Info = item.token2.issuer ? ` (${item.token2.issuer.slice(0, 8)}...)` : ''
        
        return {
          qc_key: item.id,
          name: pairName,
          symbol: pairName,
          token1: item.token1,
          token2: item.token2,
          liquidity: item.liquidity,
          volume24h: item.volume24h,
          fee: item.fee,
          apr: item.apr,
          priceChange24h: item.priceChange24h,
          color: setXrpAmmColor(item.liquidity, maxLiquidity, blueTile.value),
          size: blockSize.value === 'marketcap_index' ? item.liquidity : item.volume24h,
          // Additional fields for tooltip
          token2Issuer: item.token2.issuer || '',
          token2Name: item.token2.name,
        }
      })
      .sort((a, b) => b.size - a.size)
  })

  // Tile text template
  const tileText = computed(() => {
    return `{name}
$${$f('{liquidity}', { minDigits: 0, after: '' })}`
  })

  // Tile tooltip template
  const tileTooltip = computed(() => {
    return `{name}
Liquidity: $${$f('{liquidity}', { minDigits: 0, after: '' })}
24h Volume: $${$f('{volume24h}', { minDigits: 0, after: '' })}
APR: {apr:numberFormat.1}%
Fee: {fee:numberFormat.3}%
24h Change: {priceChange24h:numberFormat.2}%
Token 2: {token2Name}{token2Issuer ? ' (' + token2Issuer + ')' : ''}`
  })

  // Group data by liquidity tiers
  const liquidityGroups = computed<HeatmapData[]>(() => {
    if (displayGainersAndLosers.value) {
      const resp: HeatmapData[] = [
        { name: 'High Liquidity', children: [] },
        { name: 'Low Liquidity', children: [] },
      ]
      const avgLiquidity = transformedData.value.reduce((sum, item) => sum + item.liquidity, 0) / transformedData.value.length
      
      transformedData.value.forEach((elem) => {
        elem.liquidity > avgLiquidity ? resp[0].children.push(elem) : resp[1].children.push(elem)
      })
      return resp
    } else {
      return [{ name: '', children: transformedData.value }]
    }
  })

  // Update data for real-time updates
  const updateData = computed<HeatmapUpdateData>(() =>
    transformedData.value.reduce((elem, item) => ({ ...elem, [item.qc_key]: item }), {})
  )

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
    ammData,
  }
} 