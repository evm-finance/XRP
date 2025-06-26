import { computed, onBeforeUnmount, onMounted, Ref, ref, useContext, useStore, watch } from '@nuxtjs/composition-api'
import { useXrpConfigs, XrpTimeFrame, XrpBlockSize } from './useXrpConfigs'
import { useXrpFormatters } from './useXrpFormatters'
import { useXrpGraphQL } from './useXrpGraphQL'

export interface XrpHeatmapData {
  currency: string
  name: string
  price_usd: number
  price_xrp: number
  price_change_1h: number
  price_change_4h: number
  price_change_24h: number
  price_change_7d: number
  price_change_30d: number
  marketcap: number
  volume_24h: number
  liquidity: number
  tvl: number
  issuer: string
  hasAmmPool: boolean
  hasTrustLine: boolean
  color: string
  children?: XrpHeatmapData[]
}

export interface XrpHeatmapUpdateData {
  [currency: string]: {
    price_usd: number
    price_xrp: number
    price_change_1h: number
    price_change_4h: number
    price_change_24h: number
    price_change_7d: number
    price_change_30d: number
    marketcap: number
    volume_24h: number
    liquidity: number
    tvl: number
    color: string
  }
}

const setColor = (change: number, blueTile: boolean = false) => {
  const absChange = Math.abs(change)
  
  if (change > 0) {
    if (absChange <= 1) {
      return blueTile ? '#1A2A52' : '#71c175'
    } else if (absChange <= 2.5) {
      return blueTile ? '#234e91' : '#4eb153'
    } else if (absChange <= 5) {
      return blueTile ? '#336cc2' : '#3e8e42'
    } else {
      return blueTile ? '#5898f8' : '#2f6a32'
    }
  } else {
    if (absChange <= 1) {
      return blueTile ? '#EF9A9A' : '#ff8080'
    } else if (absChange <= 2.5) {
      return blueTile ? '#EF5350' : '#ff4d4d'
    } else if (absChange <= 5) {
      return blueTile ? '#D32F2F' : '#ff1a1a'
    } else {
      return blueTile ? '#B71C1C' : '#e60000'
    }
  }
}

const tileConfigs: { [key in XrpTimeFrame]: { text: string; toolTip: string } } = {
  '1h': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_1h} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '4h': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_4h} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '1d': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_24h} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '1w': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_7d} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '1m': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_30d} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '3m': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_30d} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
  '1y': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{currency}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {price_usd}
            {price_change_30d} %[/]
          `,
    toolTip: `[bold]{name}[/]
              ---------------------
              1H: {price_change_1h}%
              4H: {price_change_4h}%
              24H: {price_change_24h}%
              7D: {price_change_7d}%
              30D: {price_change_30d}%
              Market Cap: {marketcap}
              Volume 24H: {volume_24h}
              Liquidity: {liquidity}
              TVL: {tvl}`,
  },
}

export function useXrpHeatmap(USER_CAN_ACCESS_TREND_DATA: Ref<boolean> = ref(false)) {
  const { $axios } = useContext()
  const { state } = useStore()
  const { formatXrpPrice, formatMarketCap, formatVolume } = useXrpFormatters()
  const { getTokens } = useXrpGraphQL()
  
  const {
    timeFrame,
    blockSize,
    displayFavorites,
    numOfTokens,
    displayGainersAndLosers,
    blueTile,
  } = useXrpConfigs()

  const loading = ref(true)
  const error = ref<string | null>(null)
  const heatmapData = ref<XrpHeatmapData[]>([])
  const updateData = ref<XrpHeatmapUpdateData>({})

  let apiPollingTimeout: any = null
  let websocketConnection: any = null

  // Computed properties
  const favorite = computed<string | null>(() =>
    displayFavorites.value ? state.xrp?.favorites?.join(',') || null : null
  )

  const tileText = computed(() => tileConfigs[timeFrame.value].text)
  const tileTooltip = computed(() => tileConfigs[timeFrame.value].toolTip)

  const fieldNameByPerformancePeriod = computed<keyof XrpHeatmapData>(() => {
    const frames: { [key in XrpTimeFrame]: keyof XrpHeatmapData } = {
      '1h': 'price_change_1h',
      '4h': 'price_change_4h',
      '1d': 'price_change_24h',
      '1w': 'price_change_7d',
      '1m': 'price_change_30d',
      '3m': 'price_change_30d',
      '1y': 'price_change_30d',
    }
    return frames[timeFrame.value]
  })

  const fieldNameByBlockSize = computed<keyof XrpHeatmapData>(() => {
    const sizes: { [key in XrpBlockSize]: keyof XrpHeatmapData } = {
      marketcap: 'marketcap',
      volume: 'volume_24h',
      price: 'price_usd',
      liquidity: 'liquidity',
      tvl: 'tvl',
    }
    return sizes[blockSize.value]
  })

  // Data fetching functions
  const fetchHeatmapData = async () => {
    try {
      loading.value = true
      error.value = null

      // Fetch tokens data from GraphQL
      const response = await getTokens({
        limit: numOfTokens.value,
        favorites: favorite.value,
        sortBy: fieldNameByBlockSize.value,
        sortOrder: 'desc',
      })

      if (response.data?.tokens) {
        const tokens = response.data.tokens
        const processedData = processTokensForHeatmap(tokens)
        heatmapData.value = processedData
      }
    } catch (err) {
      console.error('Error fetching heatmap data:', err)
      error.value = 'Failed to load heatmap data'
      
      // Fallback to mock data for development
      if (process.env.NODE_ENV === 'development') {
        heatmapData.value = generateMockHeatmapData()
      }
    } finally {
      loading.value = false
    }
  }

  const processTokensForHeatmap = (tokens: any[]): XrpHeatmapData[] => {
    return tokens.map(token => {
      const changeField = fieldNameByPerformancePeriod.value
      const change = token[changeField] || 0
      
      return {
        currency: token.currency,
        name: token.name || token.currency,
        price_usd: token.price_usd || 0,
        price_xrp: token.price_xrp || 0,
        price_change_1h: token.price_change_1h || 0,
        price_change_4h: token.price_change_4h || 0,
        price_change_24h: token.price_change_24h || 0,
        price_change_7d: token.price_change_7d || 0,
        price_change_30d: token.price_change_30d || 0,
        marketcap: token.marketcap || 0,
        volume_24h: token.volume_24h || 0,
        liquidity: token.liquidity || 0,
        tvl: token.tvl || 0,
        issuer: token.issuer || '',
        hasAmmPool: token.hasAmmPool || false,
        hasTrustLine: token.hasTrustLine || false,
        color: setColor(change, blueTile.value),
        [fieldNameByBlockSize.value]: token[fieldNameByBlockSize.value] || 0,
      }
    })
  }

  const generateMockHeatmapData = (): XrpHeatmapData[] => {
    const mockTokens = [
      { currency: 'XRP', name: 'Ripple', price_usd: 0.52, price_change_24h: 2.5 },
      { currency: 'USDT', name: 'Tether USD', price_usd: 1.00, price_change_24h: 0.1 },
      { currency: 'USDC', name: 'USD Coin', price_usd: 1.00, price_change_24h: -0.1 },
      { currency: 'SOLO', name: 'Sologenic', price_usd: 0.15, price_change_24h: 5.2 },
      { currency: 'CSC', name: 'CasinoCoin', price_usd: 0.008, price_change_24h: -1.8 },
      { currency: 'FLR', name: 'Flare', price_usd: 0.025, price_change_24h: 3.1 },
      { currency: 'SGB', name: 'Songbird', price_usd: 0.012, price_change_24h: -2.3 },
      { currency: 'XDC', name: 'XinFin', price_usd: 0.045, price_change_24h: 1.7 },
    ]

    return mockTokens.map(token => {
      const change = token.price_change_24h || 0
      return {
        currency: token.currency,
        name: token.name,
        price_usd: token.price_usd,
        price_xrp: token.price_usd / 0.52, // Approximate XRP price
        price_change_1h: (Math.random() - 0.5) * 2,
        price_change_4h: (Math.random() - 0.5) * 3,
        price_change_24h: change,
        price_change_7d: (Math.random() - 0.5) * 10,
        price_change_30d: (Math.random() - 0.5) * 20,
        marketcap: token.price_usd * (Math.random() * 1000000 + 100000),
        volume_24h: token.price_usd * (Math.random() * 100000 + 10000),
        liquidity: token.price_usd * (Math.random() * 50000 + 5000),
        tvl: token.price_usd * (Math.random() * 200000 + 20000),
        issuer: `r${Math.random().toString(36).substring(2, 15)}`,
        hasAmmPool: Math.random() > 0.3,
        hasTrustLine: Math.random() > 0.2,
        color: setColor(change, blueTile.value),
        [fieldNameByBlockSize.value]: token.price_usd * (Math.random() * 1000000 + 100000),
      }
    })
  }

  // Real-time updates
  const setupWebSocketConnection = () => {
    if (websocketConnection) {
      websocketConnection.close()
    }

    // WebSocket connection for real-time updates
    // This would connect to your backend WebSocket endpoint
    try {
      websocketConnection = new WebSocket(process.env.XRP_WEBSOCKET_URL || 'ws://localhost:3001/ws')
      
      websocketConnection.onmessage = (event: MessageEvent) => {
        const data = JSON.parse(event.data)
        if (data.type === 'token_update') {
          updateTokenData(data.tokens)
        }
      }

      websocketConnection.onerror = (error: Event) => {
        console.error('WebSocket error:', error)
      }

      websocketConnection.onclose = () => {
        console.log('WebSocket connection closed')
        // Attempt to reconnect after 5 seconds
        setTimeout(setupWebSocketConnection, 5000)
      }
    } catch (err) {
      console.error('Failed to setup WebSocket connection:', err)
    }
  }

  const updateTokenData = (tokens: any[]) => {
    const newUpdateData: XrpHeatmapUpdateData = {}
    
    tokens.forEach(token => {
      const changeField = fieldNameByPerformancePeriod.value
      const change = token[changeField] || 0
      
      newUpdateData[token.currency] = {
        price_usd: token.price_usd || 0,
        price_xrp: token.price_xrp || 0,
        price_change_1h: token.price_change_1h || 0,
        price_change_4h: token.price_change_4h || 0,
        price_change_24h: token.price_change_24h || 0,
        price_change_7d: token.price_change_7d || 0,
        price_change_30d: token.price_change_30d || 0,
        marketcap: token.marketcap || 0,
        volume_24h: token.volume_24h || 0,
        liquidity: token.liquidity || 0,
        tvl: token.tvl || 0,
        color: setColor(change, blueTile.value),
      }
    })

    updateData.value = { ...updateData.value, ...newUpdateData }
  }

  // Polling for updates (fallback when WebSocket is not available)
  const startPolling = () => {
    if (apiPollingTimeout) {
      clearTimeout(apiPollingTimeout)
    }

    apiPollingTimeout = setTimeout(async () => {
      try {
        // Fetch updated data
        const response = await $axios.get('/api/xrp/tokens/updates', {
          params: {
            currencies: heatmapData.value.map(token => token.currency).join(','),
            timeFrame: timeFrame.value,
          }
        })

        if (response.data?.updates) {
          updateTokenData(response.data.updates)
        }
      } catch (err) {
        console.error('Error polling for updates:', err)
      }

      // Continue polling
      startPolling()
    }, 30000) // Poll every 30 seconds
  }

  const stopPolling = () => {
    if (apiPollingTimeout) {
      clearTimeout(apiPollingTimeout)
      apiPollingTimeout = null
    }
  }

  // Watchers
  watch([timeFrame, blockSize, displayFavorites, numOfTokens], () => {
    fetchHeatmapData()
  })

  watch(blueTile, () => {
    // Update colors when blue tile setting changes
    heatmapData.value = heatmapData.value.map(token => {
      const changeField = fieldNameByPerformancePeriod.value
      const change = token[changeField] || 0
      return {
        ...token,
        color: setColor(change, blueTile.value),
      }
    })
  })

  // Lifecycle
  onMounted(() => {
    fetchHeatmapData()
    setupWebSocketConnection()
    startPolling()
  })

  onBeforeUnmount(() => {
    stopPolling()
    if (websocketConnection) {
      websocketConnection.close()
    }
  })

  return {
    heatmapData,
    updateData,
    loading,
    error,
    tileText,
    tileTooltip,
    fetchHeatmapData,
    setupWebSocketConnection,
    startPolling,
    stopPolling,
  }
} 