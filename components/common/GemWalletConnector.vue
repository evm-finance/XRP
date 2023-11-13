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
              <v-avatar tile size="24" class="rounded"><v-img :src="`/img/gem-wallet.svg`" /></v-avatar>
              <v-icon small>mdi-chevron-down</v-icon>
            </div>
          </v-btn>
        </div>
      </template>
      <v-card outlined tile>
        <v-list dense class="pa-0 ma-0">
          <v-list-item-group color="primary">
            <v-list-item v-for="(item, i) in walletActions" :key="i" @click="item.action">
              <v-list-item-icon><v-icon size="18" v-text="item.icon" /></v-list-item-icon>
              <v-list-item-content><v-list-item-title />{{ item.text }}</v-list-item-content>
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
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const { dispatch } = useStore<State>()
    const { connectWallet, isWalletReady, address, disconnectWallet } = inject(XRP_PLUGIN_KEY) as XrpClient
    const { $copyAddressToClipboard } = useContext()

    const copyAccountAddress = async () => await $copyAddressToClipboard(address.value)

    // METHODS
    const walletActions = computed(() => {
      return [
        { text: 'Copy Address', icon: 'mdi-content-copy', action: copyAccountAddress },
        { text: 'View in Explorer', icon: 'mdi-open-in-new', action: navigateToExplorer },
        { text: 'Disconnect Wallet', icon: 'mdi-power', action: disconnectWallet },
      ]
    })

    function navigateToExplorer() {
      try {
        window.open(`https://xrpscan.com/account/${address.value}`, '_blank')
      } catch {}
    }

    const connectToXRPWallet = () => {
      connectWallet()
    }

    return {
      dispatch,
      isWalletReady,
      walletActions,
      connectToXRPWallet,
    }
  },
})
</script>
