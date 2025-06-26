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

export interface Network {
  id: string
  name: string
  symbol: string
  nativeTokenSymbol: string
  rpcUrl: string
  blockExplorerUrl: string
  dex: Dex[]
  weth: UniswapToken
}

export type SearchResult = {
  isWallet: boolean
  isContract: boolean
  isTransaction?: boolean
  isXRPLedger?: boolean
  isXRPAccount?: boolean
  isXRPTransaction?: boolean
  searchString: string
  network: Network | null
  desc: string
  to: string
}
