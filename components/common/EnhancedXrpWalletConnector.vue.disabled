<template>
  <client-only>
    <v-menu
      v-if="isWalletReady"
      :close-on-content-click="false"
      :nudge-width="250"
      nudge-bottom="8"
      offset-y
      max-width="230"
    >
      <template #activator="{ on, attrs }">
        <div class="d-flex">
          <v-btn class="mt-0 px-2 subtitle-2 text-capitalize font-weight-regular" text tile v-bind="attrs" v-on="on">
            <div>
              <v-avatar tile size="24" class="rounded">
                <v-img :src="walletIcon" />
              </v-avatar>
              <v-icon small>mdi-chevron-down</v-icon>
            </div>
          </v-btn>
        </div>
      </template>
      <v-card outlined tile>
        <v-list dense class="pa-0 ma-0">
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption grey--text">
                {{ walletName }}
              </v-list-item-title>
              <v-list-item-subtitle class="text-caption font-family-mono">
                {{ formatAddress(address) }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-divider></v-divider>
          <v-list-item-group color="primary">
            <v-list-item v-for="(item, i) in walletActions" :key="i" @click="item.action">
              <v-list-item-icon><v-icon size="18" v-text="item.icon" /></v-list-item-icon>
              <v-list-item-content><v-list-item-title>{{ item.text }}</v-list-item-title></v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-card>
    </v-menu>

    <v-btn v-else class="ma-0 pa-0" text @click="dispatch('ui/walletDialogStatus', true)">
      <v-icon class="ma-0 pa-0">mdi-plus</v-icon>
    </v-btn>
  </client-only>
</template>

<script lang="ts">
import { defineComponent, useStore, inject, computed, useContext } from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import { EnhancedXrpClient, ENHANCED_XRP_PLUGIN_KEY, XrpWalletType } from '~/plugins/web3/enhanced-xrp.client'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const { dispatch } = useStore<State>()
    const { 
      isWalletReady, 
      address, 
      disconnectWallet, 
      currentWalletType,
      getAccountInfo,
      getBalance
    } = inject(ENHANCED_XRP_PLUGIN_KEY) as EnhancedXrpClient
    const { $copyAddressToClipboard } = useContext()

    const copyAccountAddress = async () => await $copyAddressToClipboard(address.value)

    const navigateToExplorer = () => {
      try {
        window.open(`https://xrpscan.com/account/${address.value}`, '_blank')
      } catch {}
    }

    const viewAccountInfo = async () => {
      try {
        const result = await getAccountInfo()
        if (result.success && result.accountInfo) {
          console.log('Account Info:', result.accountInfo)
          // You could show this in a dialog or notification
        }
      } catch (error) {
        console.error('Failed to get account info:', error)
      }
    }

    const viewBalance = async () => {
      try {
        const result = await getBalance()
        if (result.success && result.balance) {
          console.log('Balance:', result.balance)
          // You could show this in a dialog or notification
        }
      } catch (error) {
        console.error('Failed to get balance:', error)
      }
    }

    // Computed properties
    const walletIcon = computed(() => {
      switch (currentWalletType.value) {
        case 'gem':
          return '/img/gem-wallet.svg'
        case 'xaman':
          return '/img/xaman-wallet.svg'
        case 'metamask-xrp-snap':
          return '/img/metamask.svg'
        default:
          return '/img/gem-wallet.svg'
      }
    })

    const walletName = computed(() => {
      switch (currentWalletType.value) {
        case 'gem':
          return 'GEM Wallet'
        case 'xaman':
          return 'Xaman (XUMM)'
        case 'metamask-xrp-snap':
          return 'MetaMask XRP'
        default:
          return 'XRP Wallet'
      }
    })

    const walletActions = computed(() => {
      const actions = [
        { text: 'Copy Address', icon: 'mdi-content-copy', action: copyAccountAddress },
        { text: 'View in Explorer', icon: 'mdi-open-in-new', action: navigateToExplorer },
      ]

      // Add wallet-specific actions
      if (currentWalletType.value === 'xaman' || currentWalletType.value === 'metamask-xrp-snap') {
        actions.push(
          { text: 'Account Info', icon: 'mdi-information', action: viewAccountInfo },
          { text: 'View Balance', icon: 'mdi-wallet', action: viewBalance }
        )
      }

      actions.push({ text: 'Disconnect Wallet', icon: 'mdi-power', action: disconnectWallet })
      
      return actions
    })

    const formatAddress = (addr: string) => {
      if (!addr) return ''
      return `${addr.slice(0, 10)}...${addr.slice(-10)}`
    }

    return {
      dispatch,
      isWalletReady,
      address,
      currentWalletType,
      walletIcon,
      walletName,
      walletActions,
      formatAddress,
    }
  },
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style> 