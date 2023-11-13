<template>
  <v-dialog v-model="dialog" max-width="525">
    <v-card tile outlined class="pa-4">
      <v-row>
        <v-col class="pa-4">
          <h5 class="text-h6 mb-2">Connect to a Wallet</h5>
          <p class="subtitle-2 font-weight-regular grey--text text--lighten-1">
            By connecting a wallet, I agree to EVM Finance
            <nuxt-link to="/terms-and-conditions" class="text-decoration-none">Terms and Conditions</nuxt-link>
          </p>
          <v-alert v-model="errorAlert" color="error" dense dismissible>{{ error.message }}</v-alert>
          <v-alert v-model="xrpErrAlert" color="error" dense dismissible>{{ gemWalletError }}</v-alert>
          <v-list-item link class="grey--text mb-2" @click="connectWallet('metamask')">
            <v-list-item-avatar tile width="50" height="50">
              <img width="60" height="60" :src="`/img/metamask.svg`" alt="metamask" />
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title class="text-h6"> Metamask</v-list-item-title>
            </v-list-item-content>
            <v-list-item-action v-if="walletReady">
              <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
            </v-list-item-action>
          </v-list-item>

          <v-list-item link class="grey--text mb-2" @click="connectGemWallet">
            <v-list-item-avatar tile width="50" height="50">
              <img width="60" height="60" :src="`/img/gem-wallet.svg`" alt="metamask" />
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title class="text-h6"> Gem Wallet</v-list-item-title>
            </v-list-item-content>
            <v-list-item-action v-if="isXRPWalletReady">
              <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
            </v-list-item-action>
          </v-list-item>

          <div class="grey--text">
            <h6 class="text-subtitle-1 font-weight-bold">New to Ethereum?</h6>
            <small>
              EVM Finance is a DeFi app on Ethereum. Set up an Ethereum Wallet to Invest and Trade here.
              <span>
                <a href="https://ethereum.org/en/wallets/" target="_blank" class="text-decoration-none">
                  Learn More <v-icon size="14" color="primary">mdi-open-in-new</v-icon>
                </a>
              </span>
            </small>
          </div>
        </v-col>
      </v-row>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useStore, inject, watch } from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Web3ErrorInterface } from '~/plugins/web3/connector'
import { XrpClient, XRP_PLUGIN_KEY } from '~/plugins/web3/xrp.client'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const store = useStore<State>()
    const { connectWallet, resetErrors, walletReady, error } = inject(WEB3_PLUGIN_KEY) as Web3
    const {
      connectWallet: connectGemWallet,
      isWalletReady: isXRPWalletReady,
      error: gemWalletError,
    } = inject(XRP_PLUGIN_KEY) as XrpClient

    // STATE
    const errorAlert = ref(false)
    const xrpErrAlert = ref(false)

    // COMPUTED
    const dialog = computed({
      get() {
        return store.state.ui.walletSelectionDialog
      },
      set(value) {
        store.dispatch('ui/walletDialogStatus', value)

        /**  Reset Error id when dialog opens or closes */
        errorAlert.value = false
        resetErrors()
      },
    })

    // WATCHES
    /** Display Error alert if wev3 error field is true */
    watch(error, (newVal: Web3ErrorInterface | null) => {
      if (newVal?.status) {
        errorAlert.value = true
      }
    })

    watch(gemWalletError, (newVal) => {
      if (newVal) {
        xrpErrAlert.value = true
      }
    })

    /** Display Error alert if wev3 error field is true */
    watch(walletReady, (val: boolean) => {
      if (val) {
        setTimeout(() => (dialog.value = false), 2000)
      }
    })

    watch(isXRPWalletReady, (val: boolean) => {
      if (val) {
        setTimeout(() => (dialog.value = false), 2000)
      }
    })

    return {
      connectWallet,
      resetErrors,
      connectGemWallet,
      xrpErrAlert,
      isXRPWalletReady,
      dialog,
      walletReady,
      error,
      errorAlert,
      gemWalletError,
    }
  },
})
</script>
