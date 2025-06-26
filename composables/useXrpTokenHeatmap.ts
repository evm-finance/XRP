import { computed, onBeforeUnmount, onMounted, ref, useContext, watch } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable'
import { XRPScreenerGQL } from '~/apollo/queries'
import type { HeatmapData, HeatmapRowData, HeatmapUpdateData } from '~/types/heatmap'
import type { HeatmapIntervals } from '~/types/state'
import useHeatmapConfigs from '~/composables/useHeatmapConfigs'

// Color function for XRP token heatmap tiles
const setXrpTokenColor = (x: number, blueTile: boolean) => {
  if (x * 100 > 0 && x * 100 <= 1) {
    return blueTile ? '#1A2A52' : '#71c175'
  } else if (x * 100 > 1 && x * 100 <= 2.5) {
    return blueTile ? '#234e91' : '#4eb153'
  } else if (x * 100 > 2.5 && x * 100 <= 5) {
    return blueTile ? '#336cc2' : '#3e8e42'
  } else if (x * 100 > 5) {
    return blueTile ? '#5898f8' : '#2f6a32'
  } else if (x * 100 <= 0 && x * 100 >= -1) {
    return blueTile ? '#EF9A9A' : '#ff8080'
  } else if (x * 100 < -1 && x * 100 >= -2.5) {
    return blueTile ? '#EF5350' : '#ff4d4d'
  } else if (x * 100 < -2.5 && x * 100 >= -5) {
    return blueTile ? '#D32F2F' : '#ff1a1a'
  } else if (x * 100 < -5) {
    return blueTile ? '#B71C1C' : '#e60000'
  }
  return ''
}

export default function useXrpTokenHeatmap(userCanAccessTrend: any) {
  const { $f } = useContext()
  const { $axios } = useContext()
  
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

  // Apollo query for XRP screener data
  const { onResult } = useQuery(XRPScreenerGQL, { 
    fetchPolicy: 'no-cache', 
    pollInterval: 60000 
  })

  // Transform XRP token data to heatmap format
  const transformedData = computed<HeatmapRowData[]>(() => {
    if (!rowData.value.length) return []

    return rowData.value
      .slice(0, numOfCoins.value)
      .map((item) => {
        // Calculate price change percentage (mock data for now)
        const priceChange = (Math.random() - 0.5) * 10 // -5% to +5%
        
        return {
          qc_key: `${item.currency}_${item.issuerAddress}`,
          name: item.tokenName || item.currency,
          symbol: item.currency,
          price: item.price || 0,
          price24h: priceChange,
          price1h: priceChange * 0.1,
          price7d: priceChange * 0.7,
          marketcap: item.marketcap || 0,
          volume24h: item.volume24H || 0,
          issuer: item.issuerAddress,
          issuerName: item.issuerName,
          icon: item.icon,
          color: setXrpTokenColor(priceChange / 100, blueTile.value),
          size: blockSize.value === 'marketcap_index' ? item.marketcap || 0 : item.volume24H || 0,
        }
      })
      .sort((a, b) => b.size - a.size)
  })

  // Tile text template
  const tileText = computed(() => {
    const timeFrameMap = {
      '1h': 'price1h',
      '24h': 'price24h',
      '7d': 'price7d',
    }
    const fieldName = timeFrameMap[timeFrame.value as keyof typeof timeFrameMap] || 'price24h'
    
    return `{name}
{${fieldName}:numberFormat.2}%`
  })

  // Tile tooltip template
  const tileTooltip = computed(() => {
    const timeFrameMap = {
      '1h': 'price1h',
      '24h': 'price24h',
      '7d': 'price7d',
    }
    const fieldName = timeFrameMap[timeFrame.value as keyof typeof timeFrameMap] || 'price24h'
    
    return `{name}
Price: $${$f('{price}', { minDigits: 6, after: '' })}
Market Cap: $${$f('{marketcap}', { minDigits: 0, after: '' })}
24h Volume: $${$f('{volume24h}', { minDigits: 0, after: '' })}
${timeFrame.value} Change: {${fieldName}:numberFormat.2}%
Issuer: {issuerName}`
  })

  // Group data by gainers and losers
  const gainersAndLosers = computed<HeatmapData[]>(() => {
    if (displayGainersAndLosers.value) {
      const resp: HeatmapData[] = [
        { name: 'Gainers', children: [] },
        { name: 'Losers', children: [] },
      ]
      transformedData.value.forEach((elem) => {
        const timeFrameMap = {
          '1h': 'price1h',
          '24h': 'price24h',
          '7d': 'price7d',
        }
        const fieldName = timeFrameMap[timeFrame.value as keyof typeof timeFrameMap] || 'price24h'
        const change = elem[fieldName as keyof typeof elem] as number
        
        change > 0 ? resp[0].children.push(elem) : resp[1].children.push(elem)
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

  // Handle query results
  onResult((queryResult: any) => {
    if (queryResult.data?.xrpScreener) {
      rowData.value = queryResult.data.xrpScreener
      heatmapData.value = gainersAndLosers.value
    }
    loading.value = false
  })

  // Watch for configuration changes
  watch([tileText, tileTooltip, displayGainersAndLosers], () => {
    heatmapData.value = gainersAndLosers.value
  })

  watch([numOfCoins, timeFrame, blockSize, blueTile], () => {
    heatmapData.value = gainersAndLosers.value
  })

  // Initialize with mock data for development
  const initializeMockData = () => {
    const mockTokens = [
      { currency: 'USDC', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'USD Coin', issuerName: 'Circle', marketcap: 1000000, price: 1.0, volume24H: 500000, icon: '🪙' },
      { currency: 'USDT', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Tether', issuerName: 'Tether', marketcap: 800000, price: 1.0, volume24H: 400000, icon: '💎' },
      { currency: 'BTC', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Bitcoin', issuerName: 'BitGo', marketcap: 500000, price: 45000, volume24H: 200000, icon: '₿' },
      { currency: 'ETH', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Ethereum', issuerName: 'BitGo', marketcap: 300000, price: 3000, volume24H: 150000, icon: 'Ξ' },
      { currency: 'SOL', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Solana', issuerName: 'Solana', marketcap: 200000, price: 100, volume24H: 100000, icon: '◎' },
      { currency: 'ADA', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Cardano', issuerName: 'IOG', marketcap: 150000, price: 0.5, volume24H: 75000, icon: '₳' },
      { currency: 'DOT', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Polkadot', issuerName: 'Web3 Foundation', marketcap: 120000, price: 8, volume24H: 60000, icon: '●' },
      { currency: 'LINK', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Chainlink', issuerName: 'Chainlink', marketcap: 100000, price: 15, volume24H: 50000, icon: '🔗' },
      { currency: 'UNI', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Uniswap', issuerName: 'Uniswap', marketcap: 80000, price: 8, volume24H: 40000, icon: '🦄' },
      { currency: 'AAVE', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Aave', issuerName: 'Aave', marketcap: 60000, price: 300, volume24H: 30000, icon: '👻' },
    ]
    
    rowData.value = mockTokens
    heatmapData.value = gainersAndLosers.value
    loading.value = false
  }

  onMounted(() => {
    // Initialize with mock data for now
    initializeMockData()
  })

  return {
    loading,
    heatmapData,
    tileText,
    tileTooltip,
    updateData,
    transformedData,
    gainersAndLosers,
  }
} 