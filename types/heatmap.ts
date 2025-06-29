import type { XRPToken } from './global'

export interface HeatmapData {
  name: string
  value: number
  color?: string
  children?: HeatmapRowData[]
}

export interface HeatmapRowData {
  name: string
  value: number
  color?: string
  children?: HeatmapData[]
  qc_key?: string
  liquidity?: number
  volume24h?: number
  apr?: number
  fee?: number
  priceChange24h?: number
  marketcap?: number
  volume24H?: number
  price?: number
  price24h?: number
  price1h?: number
  price7d?: number
  symbol?: string
  issuer?: string
  issuerName?: string
  icon?: string
  size?: number
  token1?: XRPToken & { issuer?: string; currency?: string; icon?: string }
  token2?: XRPToken & { issuer?: string; currency?: string; icon?: string }
}

export interface HeatmapUpdateData {
  name: string
  value: number
  color?: string
  [key: string]: any
}

export interface heatmapConfigInterface {
  timeFrame: string
  blockSize: string
  blueTile: boolean
  displayFavorites: boolean
  numOfCoins: number
  displayGainersAndLosers: boolean
  grouped: boolean
} 