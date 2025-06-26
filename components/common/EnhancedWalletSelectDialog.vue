<template>
  <v-dialog v-model="dialog" max-width="600">
    <v-card tile outlined class="pa-4">
      <v-row>
        <v-col class="pa-4">
          <h5 class="text-h6 mb-2">Connect to a Wallet</h5>
          <p class="subtitle-2 font-weight-regular grey--text text--lighten-1">
            By connecting a wallet, I agree to EVM Finance
            <nuxt-link to="/terms-and-conditions" class="text-decoration-none">Terms and Conditions</nuxt-link>
          </p>
          
          <!-- Error Alerts -->
          <v-alert v-model="errorAlert" color="error" dense dismissible>{{ error.message }}</v-alert>
          <v-alert v-model="xrpErrAlert" color="error" dense dismissible>{{ xrpError }}</v-alert>
          
          <!-- Wallet Options -->
          <div class="mt-4">
            <h6 class="text-subtitle-1 font-weight-bold mb-3">Ethereum Wallets</h6>
            
            <v-list-item link class="grey--text mb-2" @click="connectWallet('metamask')">
              <v-list-item-avatar tile width="50" height="50">
                <img width="60" height="60" :src="`/img/metamask.svg`" alt="metamask" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="text-h6">MetaMask</v-list-item-title>
                <v-list-item-subtitle>Connect to Ethereum network</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action v-if="walletReady">
                <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
              </v-list-item-action>
            </v-list-item>
          </div>

          <v-divider class="my-4"></v-divider>

          <div>
            <h6 class="text-subtitle-1 font-weight-bold mb-3">XRP Wallets</h6>
            
            <!-- GEM Wallet -->
            <v-list-item link class="grey--text mb-2" @click="connectXrpWallet('gem')">
              <v-list-item-avatar tile width="50" height="50">
                <img width="60" height="60" :src="`/img/gem-wallet.svg`" alt="gem-wallet" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="text-h6">GEM Wallet</v-list-item-title>
                <v-list-item-subtitle>Browser extension for XRP Ledger</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action v-if="isXrpWalletReady && currentWalletType === 'gem'">
                <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
              </v-list-item-action>
            </v-list-item>

            <!-- Xaman (XUMM) Wallet -->
            <v-list-item link class="grey--text mb-2" @click="connectXrpWallet('xaman')">
              <v-list-item-avatar tile width="50" height="50">
                <img width="60" height="60" :src="`/img/xaman-wallet.svg`" alt="xaman-wallet" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="text-h6">Xaman (XUMM)</v-list-item-title>
                <v-list-item-subtitle>Mobile wallet for XRP Ledger</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action v-if="isXrpWalletReady && currentWalletType === 'xaman'">
                <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
              </v-list-item-action>
            </v-list-item>

            <!-- MetaMask XRP Snap -->
            <v-list-item link class="grey--text mb-2" @click="connectXrpWallet('metamask-xrp-snap')">
              <v-list-item-avatar tile width="50" height="50">
                <img width="60" height="60" :src="`/img/metamask.svg`" alt="metamask-xrp" />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="text-h6">MetaMask XRP Snap</v-list-item-title>
                <v-list-item-subtitle>XRP support via MetaMask snap</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action v-if="isXrpWalletReady && currentWalletType === 'metamask-xrp-snap'">
                <v-icon color="green" size="20">mdi-checkbox-marked-circle</v-icon>
              </v-list-item-action>
            </v-list-item>
          </div>

          <!-- Help Section -->
          <div class="grey--text mt-4">
            <h6 class="text-subtitle-1 font-weight-bold">Need Help?</h6>
            <small>
              <div class="mb-2">
                <strong>GEM Wallet:</strong> 
                <a href="https://gemwallet.app/" target="_blank" class="text-decoration-none">
                  Install GEM Wallet <v-icon size="14" color="primary">mdi-open-in-new</v-icon>
                </a>
              </div>
              <div class="mb-2">
                <strong>Xaman (XUMM):</strong> 
                <a href="https://xaman.app/" target="_blank" class="text-decoration-none">
                  Download Xaman <v-icon size="14" color="primary">mdi-open-in-new</v-icon>
                </a>
              </div>
              <div class="mb-2">
                <strong>MetaMask XRP Snap:</strong> 
                <a href="https://snaps.metamask.io/snap/npm/@metamask/xrp-snap/" target="_blank" class="text-decoration-none">
                  Install XRP Snap <v-icon size="14" color="primary">mdi-open-in-new</v-icon>
                </a>
              </div>
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
import { EnhancedXrpClient, ENHANCED_XRP_PLUGIN_KEY, XrpWalletType } from '~/plugins/web3/enhanced-xrp.client'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const store = useStore<State>()
    const { connectWallet, resetErrors, walletReady, error } = inject(WEB3_PLUGIN_KEY) as Web3
    const {
      connectWallet: connectXrpWallet,
      isWalletReady: isXrpWalletReady,
      error: xrpError,
      currentWalletType,
    } = inject(ENHANCED_XRP_PLUGIN_KEY) as EnhancedXrpClient

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
        xrpErrAlert.value = false
        resetErrors()
      },
    })

    // WATCHES
    /** Display Error alert if web3 error field is true */
    watch(error, (newVal: Web3ErrorInterface | null) => {
      if (newVal?.status) {
        errorAlert.value = true
      }
    })

    watch(xrpError, (newVal) => {
      if (newVal) {
        xrpErrAlert.value = true
      }
    })

    /** Close dialog when wallet is connected */
    watch(walletReady, (val: boolean) => {
      if (val) {
        setTimeout(() => (dialog.value = false), 2000)
      }
    })

    watch(isXrpWalletReady, (val: boolean) => {
      if (val) {
        setTimeout(() => (dialog.value = false), 2000)
      }
    })

    return {
      connectWallet,
      connectXrpWallet,
      resetErrors,
      xrpErrAlert,
      isXrpWalletReady,
      currentWalletType,
      dialog,
      walletReady,
      error,
      xrpError,
      errorAlert,
    }
  },
})
</script> 