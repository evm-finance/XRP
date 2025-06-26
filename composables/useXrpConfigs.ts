import { computed, useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'

export type XrpTimeFrame = '1h' | '4h' | '1d' | '1w' | '1m' | '3m' | '1y'
export type XrpBlockSize = 'marketcap' | 'volume' | 'price' | 'liquidity' | 'tvl' | 'totalLiquidityUsd' | 'asset1ValueUsd' | 'asset2ValueUsd'
export type XrpDisplayMode = 'usd' | 'xrp'
export type XrpTheme = 'light' | 'dark' | 'auto'

export interface XrpHeatmapConfig {
  timeFrame: XrpTimeFrame
  blockSize: XrpBlockSize
  displayFavorites: boolean
  numOfTokens: number
  displayGainersAndLosers: boolean
  blueTile: boolean
}

export interface XrpDisplayConfig {
  mode: XrpDisplayMode
  theme: XrpTheme
  compactMode: boolean
  showTooltips: boolean
  autoRefresh: boolean
  refreshInterval: number
}

export interface XrpWalletConfig {
  autoConnect: boolean
  preferredWallet: string
  showBalances: boolean
  showTransactions: boolean
  notifications: boolean
}

export interface XrpScreenerConfig {
  sortBy: string
  sortOrder: 'asc' | 'desc'
  filters: {
    minMarketCap: number
    maxMarketCap: number
    minVolume: number
    maxVolume: number
    hasAmmPool: boolean
    hasTrustLine: boolean
  }
  columns: string[]
  pageSize: number
}

// Configuration options
const timeFrameOptions: { text: string; value: XrpTimeFrame }[] = [
  { text: '1 Hour', value: '1h' },
  { text: '4 Hours', value: '4h' },
  { text: '1 Day', value: '1d' },
  { text: '1 Week', value: '1w' },
  { text: '1 Month', value: '1m' },
  { text: '3 Months', value: '3m' },
  { text: '1 Year', value: '1y' },
]

const blockSizeOptions: { text: string; value: XrpBlockSize }[] = [
  { text: 'Market Cap', value: 'marketcap' },
  { text: 'Volume', value: 'volume' },
  { text: 'Price', value: 'price' },
  { text: 'Liquidity', value: 'liquidity' },
  { text: 'TVL', value: 'tvl' },
  { text: 'Total Liquidity USD', value: 'totalLiquidityUsd' },
  { text: 'Asset 1 Value USD', value: 'asset1ValueUsd' },
  { text: 'Asset 2 Value USD', value: 'asset2ValueUsd' },
]

const numOfTokensOptions: number[] = [10, 25, 50, 100, 200, 500, 1000]

const refreshIntervalOptions: { text: string; value: number }[] = [
  { text: '5 seconds', value: 5000 },
  { text: '10 seconds', value: 10000 },
  { text: '30 seconds', value: 30000 },
  { text: '1 minute', value: 60000 },
  { text: '5 minutes', value: 300000 },
]

export function useXrpConfigs() {
  const { state, dispatch } = useStore<State>()

  // Heatmap Configuration
  const timeFrame = computed({
    get: () => (state.xrp?.heatmap?.timeFrame ?? '1d') as XrpTimeFrame,
    set: (value: XrpTimeFrame) => dispatch('xrp/setHeatmapTimeFrame', value),
  })

  const blockSize = computed({
    get: () => (state.xrp?.heatmap?.blockSize ?? 'marketcap') as XrpBlockSize,
    set: (value: XrpBlockSize) => dispatch('xrp/setHeatmapBlockSize', value),
  })

  const displayFavorites = computed({
    get: () => state.xrp?.heatmap?.displayFavorites ?? false,
    set: (value: boolean) => dispatch('xrp/setHeatmapDisplayFavorites', value),
  })

  const numOfTokens = computed({
    get: () => state.xrp?.heatmap?.numOfTokens ?? 100,
    set: (value: number) => dispatch('xrp/setHeatmapNumOfTokens', value),
  })

  const displayGainersAndLosers = computed({
    get: () => state.xrp?.heatmap?.displayGainersAndLosers ?? false,
    set: (value: boolean) => dispatch('xrp/setHeatmapDisplayGainersAndLosers', value),
  })

  const blueTile = computed({
    get: () => state.xrp?.heatmap?.blueTile ?? false,
    set: (value: boolean) => dispatch('xrp/setHeatmapBlueTile', value),
  })

  // Display Configuration
  const displayMode = computed({
    get: () => (state.xrp?.display?.mode ?? 'usd') as XrpDisplayMode,
    set: (value: XrpDisplayMode) => dispatch('xrp/setDisplayMode', value),
  })

  const theme = computed({
    get: () => (state.xrp?.display?.theme ?? 'auto') as XrpTheme,
    set: (value: XrpTheme) => dispatch('xrp/setTheme', value),
  })

  const compactMode = computed({
    get: () => state.xrp?.display?.compactMode ?? false,
    set: (value: boolean) => dispatch('xrp/setCompactMode', value),
  })

  const showTooltips = computed({
    get: () => state.xrp?.display?.showTooltips ?? true,
    set: (value: boolean) => dispatch('xrp/setShowTooltips', value),
  })

  const autoRefresh = computed({
    get: () => state.xrp?.display?.autoRefresh ?? false,
    set: (value: boolean) => dispatch('xrp/setAutoRefresh', value),
  })

  const refreshInterval = computed({
    get: () => state.xrp?.display?.refreshInterval ?? 30000,
    set: (value: number) => dispatch('xrp/setRefreshInterval', value),
  })

  // Wallet Configuration
  const autoConnect = computed({
    get: () => state.xrp?.wallet?.autoConnect ?? false,
    set: (value: boolean) => dispatch('xrp/setWalletAutoConnect', value),
  })

  const preferredWallet = computed({
    get: () => state.xrp?.wallet?.preferredWallet ?? '',
    set: (value: string) => dispatch('xrp/setPreferredWallet', value),
  })

  const showBalances = computed({
    get: () => state.xrp?.wallet?.showBalances ?? true,
    set: (value: boolean) => dispatch('xrp/setShowBalances', value),
  })

  const showTransactions = computed({
    get: () => state.xrp?.wallet?.showTransactions ?? true,
    set: (value: boolean) => dispatch('xrp/setShowTransactions', value),
  })

  const notifications = computed({
    get: () => state.xrp?.wallet?.notifications ?? true,
    set: (value: boolean) => dispatch('xrp/setNotifications', value),
  })

  // Screener Configuration
  const screenerSortBy = computed({
    get: () => state.xrp?.screener?.sortBy ?? 'marketcap',
    set: (value: string) => dispatch('xrp/setScreenerSortBy', value),
  })

  const screenerSortOrder = computed({
    get: () => (state.xrp?.screener?.sortOrder ?? 'desc') as 'asc' | 'desc',
    set: (value: 'asc' | 'desc') => dispatch('xrp/setScreenerSortOrder', value),
  })

  const screenerFilters = computed({
    get: () => state.xrp?.screener?.filters ?? {
      minMarketCap: 0,
      maxMarketCap: Infinity,
      minVolume: 0,
      maxVolume: Infinity,
      hasAmmPool: false,
      hasTrustLine: false,
    },
    set: (value: any) => dispatch('xrp/setScreenerFilters', value),
  })

  const screenerColumns = computed({
    get: () => state.xrp?.screener?.columns ?? ['name', 'price', 'marketcap', 'volume'],
    set: (value: string[]) => dispatch('xrp/setScreenerColumns', value),
  })

  const screenerPageSize = computed({
    get: () => state.xrp?.screener?.pageSize ?? 25,
    set: (value: number) => dispatch('xrp/setScreenerPageSize', value),
  })

  // Computed properties for UI
  const blockSizeName = computed(() => 
    blockSizeOptions.find((elem) => elem.value === blockSize.value) || blockSizeOptions[0]
  )

  const timeFrameName = computed(() => 
    timeFrameOptions.find((elem) => elem.value === timeFrame.value) || timeFrameOptions[2]
  )

  const refreshIntervalName = computed(() => 
    refreshIntervalOptions.find((elem) => elem.value === refreshInterval.value) || refreshIntervalOptions[2]
  )

  // Configuration presets
  const presets = {
    default: {
      heatmap: {
        timeFrame: '1d' as XrpTimeFrame,
        blockSize: 'marketcap' as XrpBlockSize,
        displayFavorites: false,
        numOfTokens: 100,
        displayGainersAndLosers: false,
        blueTile: false,
      },
      display: {
        mode: 'usd' as XrpDisplayMode,
        theme: 'auto' as XrpTheme,
        compactMode: false,
        showTooltips: true,
        autoRefresh: false,
        refreshInterval: 30000,
      },
      wallet: {
        autoConnect: false,
        preferredWallet: '',
        showBalances: true,
        showTransactions: true,
        notifications: true,
      },
      screener: {
        sortBy: 'marketcap',
        sortOrder: 'desc' as 'desc',
        filters: {
          minMarketCap: 0,
          maxMarketCap: Infinity,
          minVolume: 0,
          maxVolume: Infinity,
          hasAmmPool: false,
          hasTrustLine: false,
        },
        columns: ['name', 'price', 'marketcap', 'volume'],
        pageSize: 25,
      },
    },
    trader: {
      heatmap: {
        timeFrame: '1h' as XrpTimeFrame,
        blockSize: 'volume' as XrpBlockSize,
        displayFavorites: true,
        numOfTokens: 200,
        displayGainersAndLosers: true,
        blueTile: true,
      },
      display: {
        mode: 'usd' as XrpDisplayMode,
        theme: 'dark' as XrpTheme,
        compactMode: true,
        showTooltips: true,
        autoRefresh: true,
        refreshInterval: 10000,
      },
      wallet: {
        autoConnect: true,
        preferredWallet: 'gem',
        showBalances: true,
        showTransactions: true,
        notifications: true,
      },
      screener: {
        sortBy: 'volume',
        sortOrder: 'desc' as 'desc',
        filters: {
          minMarketCap: 1000000,
          maxMarketCap: Infinity,
          minVolume: 100000,
          maxVolume: Infinity,
          hasAmmPool: true,
          hasTrustLine: false,
        },
        columns: ['name', 'price', 'volume', 'ammPool', 'trustLine'],
        pageSize: 50,
      },
    },
    analyst: {
      heatmap: {
        timeFrame: '1w' as XrpTimeFrame,
        blockSize: 'tvl' as XrpBlockSize,
        displayFavorites: false,
        numOfTokens: 500,
        displayGainersAndLosers: false,
        blueTile: false,
      },
      display: {
        mode: 'usd' as XrpDisplayMode,
        theme: 'light' as XrpTheme,
        compactMode: false,
        showTooltips: true,
        autoRefresh: false,
        refreshInterval: 60000,
      },
      wallet: {
        autoConnect: false,
        preferredWallet: '',
        showBalances: false,
        showTransactions: false,
        notifications: false,
      },
      screener: {
        sortBy: 'marketcap',
        sortOrder: 'desc' as 'desc',
        filters: {
          minMarketCap: 0,
          maxMarketCap: Infinity,
          minVolume: 0,
          maxVolume: Infinity,
          hasAmmPool: false,
          hasTrustLine: false,
        },
        columns: ['name', 'price', 'marketcap', 'volume', 'issuer', 'ammPool'],
        pageSize: 100,
      },
    },
  }

  // Apply preset configuration
  const applyPreset = (presetName: keyof typeof presets) => {
    const preset = presets[presetName]
    if (preset) {
      dispatch('xrp/applyPreset', preset)
    }
  }

  // Reset to default configuration
  const resetToDefault = () => {
    applyPreset('default')
  }

  // Save configuration to localStorage
  const saveConfig = () => {
    const config = {
      heatmap: {
        timeFrame: timeFrame.value,
        blockSize: blockSize.value,
        displayFavorites: displayFavorites.value,
        numOfTokens: numOfTokens.value,
        displayGainersAndLosers: displayGainersAndLosers.value,
        blueTile: blueTile.value,
      },
      display: {
        mode: displayMode.value,
        theme: theme.value,
        compactMode: compactMode.value,
        showTooltips: showTooltips.value,
        autoRefresh: autoRefresh.value,
        refreshInterval: refreshInterval.value,
      },
      wallet: {
        autoConnect: autoConnect.value,
        preferredWallet: preferredWallet.value,
        showBalances: showBalances.value,
        showTransactions: showTransactions.value,
        notifications: notifications.value,
      },
      screener: {
        sortBy: screenerSortBy.value,
        sortOrder: screenerSortOrder.value,
        filters: screenerFilters.value,
        columns: screenerColumns.value,
        pageSize: screenerPageSize.value,
      },
    }
    
    localStorage.setItem('xrp-config', JSON.stringify(config))
  }

  // Load configuration from localStorage
  const loadConfig = () => {
    const saved = localStorage.getItem('xrp-config')
    if (saved) {
      try {
        const config = JSON.parse(saved)
        dispatch('xrp/loadConfig', config)
      } catch (error) {
        console.error('Failed to load XRP configuration:', error)
      }
    }
  }

  return {
    // Heatmap configuration
    timeFrame,
    timeFrameOptions,
    blockSize,
    blockSizeOptions,
    blockSizeName,
    displayFavorites,
    numOfTokens,
    numOfTokensOptions,
    displayGainersAndLosers,
    blueTile,

    // Display configuration
    displayMode,
    theme,
    compactMode,
    showTooltips,
    autoRefresh,
    refreshInterval,
    refreshIntervalOptions,
    refreshIntervalName,

    // Wallet configuration
    autoConnect,
    preferredWallet,
    showBalances,
    showTransactions,
    notifications,

    // Screener configuration
    screenerSortBy,
    screenerSortOrder,
    screenerFilters,
    screenerColumns,
    screenerPageSize,

    // Computed properties
    timeFrameName,

    // Presets and utilities
    presets,
    applyPreset,
    resetToDefault,
    saveConfig,
    loadConfig,
  }
} 