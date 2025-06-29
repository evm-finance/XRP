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

export interface XRPToken {
  symbol: string
  name: string
  decimals: number
}

export interface Network {
  id: string
  name: string
  symbol: string
  nativeTokenSymbol: string
  rpcUrl: string
  blockExplorerUrl: string
  dex: Dex[]
  xrp: XRPToken
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
