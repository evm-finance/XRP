import { Maybe } from '~/types/apollo/main/types'

export interface Web3ErrorInterface {
  status: boolean
  message: string | null
}
export interface ConnectorInterface {
  id: string
  provider: any
  account: string | null
  chainId: number | null
  active: boolean
  error: Web3ErrorInterface
  connect: () => Promise<{ account: string | null; error: Web3ErrorInterface }>
  handleDisconnect: () => void
  registerListeners: () => void
  resetErrors: () => void
}

export interface ChainChangeParamInterface {
  chainId: string
  chainName: string
  nativeCurrency: {
    name: string
    symbol: string
    decimals: number
  }
  blockExplorerUrls?: (Maybe<string> | undefined)[]
  rpcUrls: string[]
}
