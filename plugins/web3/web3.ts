import { reactive, Ref, computed, toRefs, onGlobalSetup, provide, ref } from '@nuxtjs/composition-api'
import { Context } from '@nuxt/types'
import { ethers } from 'ethers'
import { MetamaskConnector } from '~/plugins/web3/metamask.connector'
import { ConnectorInterface, Web3ErrorInterface } from '~/plugins/web3/connector'
import { Cookies } from '~/types/cookies'
import { Network } from '~/types/global'

export type Wallet = 'metamask'

type WalletState = 'connecting' | 'connected' | 'disconnected'
export const WEB3_PLUGIN_KEY = '$web3'

export type Web3 = {
  connectWallet: (wallet: Wallet) => Promise<void>
  disconnectWallet: () => void
  resetErrors: () => void
  changeNetwork: (chain: Network) => void
  importTokenToMetamask: (params: { address: string; symbol: string; decimals: number; image: string }) => Promise<void>
  provider: Ref<ethers.providers.Web3Provider | null>
  account: Ref<string>
  chainId: Ref<number | null>
  connector: Ref<ConnectorInterface | null>
  walletState: Ref<WalletState>
  walletReady: Ref<boolean>
  error: Ref<Web3ErrorInterface | null>
  signer: Ref<any>
  allNetworks: Ref<Network[]>
  currentNetwork: Ref<Network | null>
  getCustomProviderByNetworkId: (networkId: string) => ethers.providers.JsonRpcProvider | null
  getNetworkById: (networkId: string) => Network | null
  getNetworkByChainNumber: (chainIdentifier: number) => Network | null
  isWrapped: (address: string, network: Network) => boolean
}

type PluginState = {
  connector: ConnectorInterface | null
  walletState: WalletState
}

const WalletConnectorDictionary: Record<Wallet, ConnectorInterface> = {
  metamask: new MetamaskConnector(),
}

export default (context: Context): void => {
  onGlobalSetup(async () => {
    const pluginState = reactive<PluginState>({
      connector: null,
      walletState: 'disconnected',
    })

    const errorStatus: Ref<Web3ErrorInterface> = ref({ status: false, message: null })

    // COMPUTED
    const account = computed(() => pluginState.connector?.account ?? '')
    const chainId = computed(() => pluginState.connector?.chainId ?? null)
    const provider = computed(() => pluginState.connector?.provider ?? null)
    const signer = computed(() => pluginState.connector?.provider?.getSigner())
    const walletReady = computed(() => {
      return !!(pluginState.connector && pluginState.connector.provider && pluginState.walletState === 'connected')
    })
    const allNetworks = computed<Network[]>(() => context.store.state.configs.networks)
    const currentNetwork = computed<Network | null>(
      () => allNetworks.value.find((elem) => elem.chainIdentifier === chainId.value) ?? null
    )

    // METHODS
    const connectWallet = async (wallet: Wallet) => {
      pluginState.walletState = 'connecting'

      try {
        if (!wallet) {
          throw new Error('Please provide a wallet to facilitate a web3 connection.')
        }
        const connector = WalletConnectorDictionary[wallet]
        pluginState.connector = connector
        pluginState.connector.resetErrors()

        if (!connector) {
          throw new Error(`Wallet [${wallet}] is not supported yet. Please contact the dev team to add this connector.`)
        }
        const { account, error } = await pluginState.connector.connect()

        // Toggling Error if account is not detected or Metamask is not installed
        if (!account) {
          errorStatus.value = { status: error.status, message: error.message }
          pluginState.walletState = 'disconnected'
          return
        }
        if (account) {
          pluginState.walletState = 'connected'
          context.$cookies.set(Cookies.walletConnected, wallet)
        }
      } catch (err) {
        pluginState.walletState = 'disconnected'
      }
    }
    const getCustomProviderByNetworkId = (networkId: string): ethers.providers.JsonRpcProvider | null => {
      const network = allNetworks.value.find((elem) => elem.id === networkId)
      if (network) {
        return new ethers.providers.JsonRpcProvider(network.rpcUrl)
      }
      return null
    }
    const getNetworkById = (networkId: string): Network | null => {
      const network = allNetworks.value.find((elem) => elem.id === networkId)
      if (network) {
        return network
      }
      return null
    }

    const getNetworkByChainNumber = (chainIdentifier: number): Network | null => {
      const network = allNetworks.value.find((elem) => elem.chainIdentifier === chainIdentifier)
      if (network) {
        return network
      }
      return null
    }

    const isWrapped = (address: string, network: Network): boolean =>
      address.toLowerCase() === network.weth.address.toLowerCase()

    const disconnectWallet = () => {
      if (!pluginState.connector) {
        throw new Error('Cannot disconnect a wallet. No wallet currently connected.')
      }
      const connector = pluginState.connector as ConnectorInterface
      connector.handleDisconnect()
      pluginState.connector = null
      pluginState.walletState = 'disconnected'
      context.$cookies.remove(Cookies.walletConnected)
    }

    const resetErrors = () => {
      errorStatus.value = { status: false, message: null }
      if (pluginState.connector) {
        pluginState.connector.resetErrors()
      }
    }

    const changeNetwork = (chain: Network) => {
      if (pluginState.connector) {
        pluginState.connector.handleChanChange(chain)
      }
    }

    const importTokenToMetamask = async (params: {
      address: string
      symbol: string
      decimals: number
      image: string
    }) => {
      if (pluginState.connector) {
        await pluginState.connector.importTokenToMetamask(params)
      }
    }

    const plugin: Web3 = {
      connectWallet,
      disconnectWallet,
      changeNetwork,
      importTokenToMetamask,
      resetErrors,
      getCustomProviderByNetworkId,
      getNetworkById,
      getNetworkByChainNumber,
      ...toRefs(pluginState),
      isWrapped,
      account,
      chainId,
      provider,
      signer,
      walletReady,
      error: errorStatus,
      allNetworks,
      currentNetwork,
    }

    provide(WEB3_PLUGIN_KEY, plugin)

    const alreadyConnected: Wallet | undefined = context.$cookies.get(Cookies.walletConnected)
    if (alreadyConnected) {
      await connectWallet(alreadyConnected)
    }
  })
}
