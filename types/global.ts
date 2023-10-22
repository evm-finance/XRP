import { UniswapToken } from '~/types/apollo/main/types'

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
  blockExplorerUrl: string
  chainIdentifier: number
  dex: Dex[]
  id: string
  name: string
  rpcUrl: string
  symbol: string
  weth: UniswapToken
}

export type SearchResult = {
  isWallet: boolean
  isContract: boolean
  isTransaction?: boolean
  isXRPLedger?: boolean
  searchString: string
  network: Network | null
  desc: string
  to: string
}
