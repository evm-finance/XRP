import { computed, onGlobalSetup, provide, reactive, ref, Ref } from '@nuxtjs/composition-api'
import { getAddress, isInstalled, on } from '@gemwallet/api'
import { Context } from '@nuxt/types'
import { Cookies } from '~/types/cookies'

export const XRP_PLUGIN_KEY = '$xrp'
const wallet: string = 'xrp'

export type XrpClient = {
  connectWallet: () => void
  disconnectWallet: () => void
  isWalletReady: Ref<boolean>
  address: Ref<string>
  error: Ref<string | null>
}

type XrpState = {
  installed: boolean
  address: string
}

export default (context: Context): void => {
  onGlobalSetup(async () => {
    const state = reactive<XrpState>({
      installed: false,
      address: '',
    })

    const address = computed(() => state.address)
    const errorStatus = ref<string | null>(null)
    const isWalletReady = computed(() => state.address.length > 0)

    const connectWallet = async () => {
      try {
        errorStatus.value = null
        const installed = await isInstalled()
        if (installed.result.isInstalled) {
          const address = await getAddress()
          state.address = address.result?.address ?? ''
          const inFifteenMinutes = new Date(new Date().getTime() + 15 * 60 * 1000)
          context.$cookies.set(Cookies.expWalletConnected, wallet, { expires: inFifteenMinutes })
        } else {
          errorStatus.value = 'GEM wallet not installed. Go to https://gemwallet.app/ to install.'
        }
      } catch (error) {
        errorStatus.value = "Can't connect to XRP wallet"
      }
    }
    const disconnectWallet = () => {
      state.address = ''
      errorStatus.value = null
      context.$cookies.remove(Cookies.expWalletConnected)
    }

    on('login', async () => {
      await connectWallet()
    })
    on('logout', () => {
      disconnectWallet()
    })

    const plugin: XrpClient = {
      connectWallet,
      disconnectWallet,
      address,
      error: errorStatus,
      isWalletReady,
    }

    provide(XRP_PLUGIN_KEY, plugin)

    const alreadyConnected: string | undefined = context.$cookies.get(Cookies.expWalletConnected)
    if (alreadyConnected) {
      await connectWallet()
    }
  })
}
