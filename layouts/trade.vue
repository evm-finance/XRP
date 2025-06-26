<template>
  <v-app>
    <v-app-bar app elevation="0" height="60" style="border-bottom: none">
      <v-avatar size="32" class="rounded-0">
        <v-img :src="imageUrl" :lazy-src="imageUrl" />
      </v-avatar>

      <a href="https://evm.finance" class="text-decoration-none mr-10" style="color: inherit">
        <v-toolbar-title class="ml-2 d-flex">
          <div>
            <div class="subtitle-1 mb-n2">EVM</div>
            <div class="caption grey--text">finance</div>
          </div>
        </v-toolbar-title>
      </a>

      <v-spacer />

      <div class="pr-3">
        <v-card tile color="grey darken-4"> </v-card>
      </div>

      <wallet-connector />
      <wallet-select-dialog />
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'
import WalletConnector from '~/components/common/WalletConnector.vue'
import { useInitTheme } from '~/composables/useInitTheme'
import WalletSelectDialog from '~/components/common/WalletSelectDialog.vue'

export default defineComponent({
  components: { WalletSelectDialog, WalletConnector },
  setup() {
    // COMPOSABLE
    useInitTheme()

    // COMPUTED
    const imageUrl = computed(() => `/img/logo/evmfinance-logo.svg`)
    return { imageUrl }
  },
})
</script>
