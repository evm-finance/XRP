import { computed, onGlobalSetup, provide, reactive, ref, Ref } from '@nuxtjs/composition-api'
import { getAddress, isInstalled, on } from '@gemwallet/api'
import { Context } from '@nuxt/types'
import { Cookies } from '~/types/cookies'
import { XamanConnector } from './xaman.connector'
import { MetamaskXrpSnapConnector } from './metamask-xrp-snap.connector'

export const ENHANCED_XRP_PLUGIN_KEY = '$enhancedXrp'

export type XrpWalletType = 'gem' | 'xaman' | 'metamask-xrp-snap'

export type EnhancedXrpClient = {
  connectWallet: (walletType: XrpWalletType) => Promise<void>
  disconnectWallet: () => void
  isWalletReady: Ref<boolean>
  address: Ref<string>
  error: Ref<string | null>
  currentWalletType: Ref<XrpWalletType | null>
  getWalletConnector: (walletType: XrpWalletType) => any
  signTransaction: (transaction: any) => Promise<{ success: boolean; signedTx?: any; error?: string }>
  getAccountInfo: () => Promise<{ success: boolean; accountInfo?: any; error?: string }>
  getBalance: () => Promise<{ success: boolean; balance?: string; error?: string }>
}

type XrpState = {
  installed: boolean
  address: string
  walletType: XrpWalletType | null
}

export default (context: Context): void => {
  onGlobalSetup(async () => {
    const state = reactive<XrpState>({
      installed: false,
      address: '',
      walletType: null,
    })

    const address = computed(() => state.address)
    const errorStatus = ref<string | null>(null)
    const isWalletReady = computed(() => state.address.length > 0)
    const currentWalletType = computed(() => state.walletType)

    // Wallet connectors
    const xamanConnector = new XamanConnector()
    const metamaskXrpSnapConnector = new MetamaskXrpSnapConnector()

    const getWalletConnector = (walletType: XrpWalletType) => {
      switch (walletType) {
        case 'xaman':
          return xamanConnector
        case 'metamask-xrp-snap':
          return metamaskXrpSnapConnector
        default:
          return null
      }
    }

    const connectGemWallet = async (): Promise<{ success: boolean; address?: string; error?: string }> => {
      try {
        const installed = await isInstalled()
        if (installed.result.isInstalled) {
          const addressResult = await getAddress()
          const address = addressResult.result?.address ?? ''
          return { success: true, address }
        } else {
          return { 
            success: false, 
            error: 'GEM wallet not installed. Go to https://gemwallet.app/ to install.' 
          }
        }
      } catch (error) {
        return { success: false, error: "Can't connect to GEM wallet" }
      }
    }

    const connectWallet = async (walletType: XrpWalletType) => {
      try {
        errorStatus.value = null
        
        let result: { success: boolean; address?: string; error?: string } = { success: false }

        switch (walletType) {
          case 'gem':
            result = await connectGemWallet()
            break
          case 'xaman':
            const xamanResult = await xamanConnector.connect()
            result = {
              success: !!xamanResult.account,
              address: xamanResult.account || undefined,
              error: xamanResult.error.message || undefined
            }
            break
          case 'metamask-xrp-snap':
            const metamaskResult = await metamaskXrpSnapConnector.connect()
            result = {
              success: !!metamaskResult.account,
              address: metamaskResult.account || undefined,
              error: metamaskResult.error.message || undefined
            }
            break
          default:
            result = { success: false, error: 'Unsupported wallet type' }
        }

        if (result.success && result.address) {
          state.address = result.address
          state.walletType = walletType
          const inFifteenMinutes = new Date(new Date().getTime() + 15 * 60 * 1000)
          context.$cookies.set(Cookies.expWalletConnected, walletType, { expires: inFifteenMinutes })
        } else {
          errorStatus.value = result.error || 'Failed to connect to wallet'
        }
      } catch (error) {
        errorStatus.value = "Can't connect to XRP wallet"
      }
    }

    const disconnectWallet = () => {
      state.address = ''
      state.walletType = null
      errorStatus.value = null
      context.$cookies.remove(Cookies.expWalletConnected)
      
      // Disconnect from specific wallet connectors
      if (state.walletType === 'xaman') {
        xamanConnector.disconnect()
      } else if (state.walletType === 'metamask-xrp-snap') {
        metamaskXrpSnapConnector.disconnect()
      }
    }

    const signTransaction = async (transaction: any): Promise<{ success: boolean; signedTx?: any; error?: string }> => {
      if (!state.walletType) {
        return { success: false, error: 'No wallet connected' }
      }

      switch (state.walletType) {
        case 'gem':
          // GEM wallet transaction signing would be implemented here
          return { success: false, error: 'GEM wallet transaction signing not implemented yet' }
        case 'xaman':
          return await xamanConnector.signTransaction(transaction)
        case 'metamask-xrp-snap':
          return await metamaskXrpSnapConnector.signTransaction(transaction)
        default:
          return { success: false, error: 'Unsupported wallet type for transaction signing' }
      }
    }

    const getAccountInfo = async (): Promise<{ success: boolean; accountInfo?: any; error?: string }> => {
      if (!state.walletType) {
        return { success: false, error: 'No wallet connected' }
      }

      switch (state.walletType) {
        case 'gem':
          // GEM wallet account info would be implemented here
          return { success: false, error: 'GEM wallet account info not implemented yet' }
        case 'xaman':
          return await xamanConnector.getAccountInfo()
        case 'metamask-xrp-snap':
          return await metamaskXrpSnapConnector.getAccountInfo()
        default:
          return { success: false, error: 'Unsupported wallet type for account info' }
      }
    }

    const getBalance = async (): Promise<{ success: boolean; balance?: string; error?: string }> => {
      if (!state.walletType) {
        return { success: false, error: 'No wallet connected' }
      }

      switch (state.walletType) {
        case 'gem':
          // GEM wallet balance would be implemented here
          return { success: false, error: 'GEM wallet balance not implemented yet' }
        case 'xaman':
          // Xaman balance would be implemented here
          return { success: false, error: 'Xaman balance not implemented yet' }
        case 'metamask-xrp-snap':
          return await metamaskXrpSnapConnector.getBalance()
        default:
          return { success: false, error: 'Unsupported wallet type for balance' }
      }
    }

    // GEM wallet event listeners
    on('login', async () => {
      if (state.walletType === 'gem') {
        await connectWallet('gem')
      }
    })
    
    on('logout', () => {
      if (state.walletType === 'gem') {
        disconnectWallet()
      }
    })

    const plugin: EnhancedXrpClient = {
      connectWallet,
      disconnectWallet,
      address,
      error: errorStatus,
      isWalletReady,
      currentWalletType,
      getWalletConnector,
      signTransaction,
      getAccountInfo,
      getBalance,
    }

    provide(ENHANCED_XRP_PLUGIN_KEY, plugin)

    // Auto-connect if previously connected
    const alreadyConnected: string | undefined = context.$cookies.get(Cookies.expWalletConnected)
    if (alreadyConnected && ['gem', 'xaman', 'metamask-xrp-snap'].includes(alreadyConnected)) {
      await connectWallet(alreadyConnected as XrpWalletType)
    }
  })
} 