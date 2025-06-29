import { XRPToken } from '~/types/apollo/main/types'

export interface ERC20Balance {
  address: string
  balance: number
  isNative: boolean
  isWrapped?: boolean
  nativeBalance?: number
  nativeSymbol?: string
  uniqueKeyOrSymbol?: string
}

export type Dex = {
  name: string
  symbol: string
  value: string
}

export type Network = {
  id: string
  name: string
  symbol: string
  nativeTokenSymbol: string
  rpcUrl: string
  blockExplorerUrl: string
  xrp: XRPToken
  dex: Array<{
    value: string
    name: string
    symbol: string
  }>
}

export type SearchResult = {
  id: string
  name: string
  symbol: string
  type: 'token' | 'pool' | 'account'
  address?: string
  issuer?: string
  network?: string
}
