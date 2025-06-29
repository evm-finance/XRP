import { Network, SearchResult } from '~/types/global'

export type ThemeOptions = 'dark' | 'light'

export type HeatmapIntervals = '1h' | '4h' | '1d' | '1w' | '1m' | '3m' | '1y'
export type HeatmapTileSize = 'marketcap' | 'volume' | 'price' | 'change'

export interface ConfigState {
  title: string
  networks: Network[]
  balancesChains: string[]
  protocols: { name: string; symbol: string; id: string }[]
  globalSearchResult: SearchResult[]
}

export interface UiState {
  theme: ThemeOptions
  gasStats: XRPGasStats[] | null
  searchResults: SearchResult[]
  dark: { [key: string]: string }
  light: { [key: string]: string }
  walletSelectionDialog: boolean
}

export interface XRPGasStats {
  network: string
  fee: number
  symbol: string
}

export interface WalletState {
  address: string | null
  isWalletConnected: boolean
  totalBalance: string
}

export interface State {
  configs: ConfigState
  ui: UiState
  wallet: WalletState
  xrp?: {
    favorites?: string[]
    heatmap?: {
      timeFrame?: string
      blockSize?: string
      displayFavorites?: boolean
      numOfTokens?: number
      displayGainersAndLosers?: boolean
      blueTile?: boolean
    }
    display?: {
      mode?: string
      theme?: string
      compactMode?: boolean
      showTooltips?: boolean
      autoRefresh?: boolean
      refreshInterval?: number
    }
    wallet?: {
      autoConnect?: boolean
      preferredWallet?: string
      showBalances?: boolean
      showTransactions?: boolean
      notifications?: boolean
    }
    screener?: {
      sortBy?: string
      sortOrder?: string
      filters?: any
      columns?: string[]
      pageSize?: number
    }
  }
  global?: {
    heatmap?: {
      timeFrame?: string
      blockSize?: string
      displayFavorites?: boolean
      numOfCoins?: number
      displayGainersAndLosers?: boolean
      blueTile?: boolean
    }
  }
}
