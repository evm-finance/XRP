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
          <v-btn class="mt-1 px-2 subtitle-2 text-capitalize font-weight-regular" text tile v-bind="attrs" v-on="on">
            <div>
              <v-avatar tile size="24" class="rounded"><v-img :src="avatar" /></v-avatar>
              <span v-if="$vuetify.breakpoint.mdAndUp" class="ml-2" v-text="accountShortCut" />
              <v-icon small>mdi-chevron-down</v-icon>
            </div>
          </v-btn>
        </div>
      </template>
      <v-card outlined tile>
        <v-row no-gutters align="center">
          <v-col cols="3" class="pa-2">
            <v-avatar tile size="36" class="rounded"><v-img :src="avatar" alt="" /></v-avatar>
          </v-col>
          <v-col> <span class="text-subtitle-1 font-weight-bold ml-2" v-text="accountShortCut" /> </v-col>
        </v-row>
        <v-divider class="my-1" />
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
    <v-btn v-else tile depressed @click="connectWallet">Connect XRP Wallet</v-btn>
  </client-only>
</template>

<script lang="ts">
import { computed, defineComponent, inject, useContext } from '@nuxtjs/composition-api'
import { create } from 'blockies-ts'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const { connectWallet, disconnectWallet, address, isWalletReady } = inject(XRP_PLUGIN_KEY) as XrpClient
    const { $copyAddressToClipboard } = useContext()

    // COMPUTED
    const avatar = computed(() => create({ seed: address.value }).toDataURL() || '')
    const accountShortCut = computed(
      () => `${address.value.slice(0, 5)}.....${address.value.slice(address.value.length - 5, address.value.length)}`
    )

    // METHODS
    const walletActions = computed(() => {
      return [
        { text: 'Copy Address', icon: 'mdi-content-copy', action: copyAccountAddress },
        { text: 'View on XRP Explorer', icon: 'mdi-open-in-new', action: navigateToExplorer },
        { text: 'Disconnect Wallet', icon: 'mdi-power', action: disconnectWallet },
      ]
    })

    const copyAccountAddress = async () => await $copyAddressToClipboard(address.value)

    function navigateToExplorer() {
      try {
        window.open(`https://xrpscan.com/account/${address.value}`, '_blank')
      } catch {}
    }

    return {
      connectWallet,
      disconnectWallet,
      isWalletReady,
      address,
      avatar,
      accountShortCut,
      walletActions,
    }
  },
})
</script>
