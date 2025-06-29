import { reactive, Ref, computed, toRefs, onGlobalSetup, provide, ref } from '@nuxtjs/composition-api'
import { Context } from '@nuxt/types'
import { MetamaskConnector } from './metamask.connector'
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
  provider: Ref<any | null>
  account: Ref<string>
  chainId: Ref<number | null>
  connector: Ref<ConnectorInterface | null>
  walletState: Ref<WalletState>
  walletReady: Ref<boolean>
  error: Ref<Web3ErrorInterface | null>
  signer: Ref<any>
  allNetworks: Ref<Network[]>
  currentNetwork: Ref<Network | null>
  getCustomProviderByNetworkId: (networkId: string) => any | null
  getNetworkById: (networkId: string) => Network | null
  getNetworkByChainNumber: (chainId: number) => Network | null
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

    const errorStatus: Ref<Web3ErrorInterface> = ref({ code: 0, message: '' })

    // COMPUTED
    const account = computed(() => pluginState.connector?.getAccount() ?? '')
    const chainId = computed(() => null) // XRP doesn't use chainId
    const provider = computed(() => null) // XRP doesn't use providers
    const signer = computed(() => null) // XRP doesn't use signers
    const walletReady = computed(() => {
      return !!(pluginState.connector && pluginState.connector.isConnected() && pluginState.walletState === 'connected')
    })
    const allNetworks = ref<Network[]>([])
    const currentNetwork = computed<Network | null>(
      () => allNetworks.value.find((elem) => elem.id === 'xrp') ?? null
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

        if (!connector) {
          throw new Error(`Wallet [${wallet}] is not supported yet. Please contact the dev team to add this connector.`)
        }
        
        await pluginState.connector.connect()
        const account = pluginState.connector.getAccount()

        // Toggling Error if account is not detected
        if (!account) {
          errorStatus.value = { code: 1, message: 'No account detected' }
          pluginState.walletState = 'disconnected'
          return
        }
        
        if (account) {
          pluginState.walletState = 'connected'
          const inFifteenMinutes = new Date(new Date().getTime() + 15 * 60 * 1000)
          context.$cookies.set(Cookies.walletConnected, wallet, { expires: inFifteenMinutes })
        }
      } catch (err) {
        pluginState.walletState = 'disconnected'
        errorStatus.value = { code: 1, message: err instanceof Error ? err.message : 'Connection failed' }
      }
    }
    
    const getCustomProviderByNetworkId = (networkId: string): any | null => {
      // XRP doesn't use providers
      return null
    }
    
    const getNetworkById = (networkId: string): Network | null => {
      const network = allNetworks.value.find((elem) => elem.id === networkId)
      if (network) {
        return network
      }
      return null
    }

    const getNetworkByChainNumber = (chainId: number): Network | null => {
      const network = allNetworks.value.find((elem) => elem.id === chainId.toString())
      if (network) {
        return network
      }
      return null
    }

    const isWrapped = (address: string, network: Network): boolean => {
      // XRP doesn't have wrapped tokens like ETH
      return false
    }

    const disconnectWallet = () => {
      if (!pluginState.connector) {
        throw new Error('Cannot disconnect a wallet. No wallet currently connected.')
      }
      pluginState.connector.disconnect()
      pluginState.connector = null
      pluginState.walletState = 'disconnected'
      context.$cookies.remove(Cookies.walletConnected)
    }

    const resetErrors = () => {
      errorStatus.value = { code: 0, message: '' }
    }

    const changeNetwork = (chain: Network) => {
      // XRP doesn't support network switching
    }

    const importTokenToMetamask = async (params: {
      address: string
      symbol: string
      decimals: number
      image: string
    }) => {
      // XRP doesn't use Metamask token import
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

