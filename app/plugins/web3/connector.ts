import { ethers } from 'ethers'
import { Maybe } from '~/types/apollo/main/types'
import { Network } from '~/types/global'

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
  handleChanChange: (chain: Network) => void
  importTokenToMetamask: (params: { address: string; symbol: string; decimals: number; image: string }) => Promise<void>
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
