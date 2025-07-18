<template>
  <v-app>
    <v-app-bar app elevation="0" outlined height="60" class="hidden-xs-and-down">
      <v-app-bar-nav-icon class="d-flex d-sm-none" @click="drawer = true"></v-app-bar-nav-icon>

      <v-avatar size="32" class="rounded-0">
        <v-img :src="imageUrl" :lazy-src="imageUrl" />
      </v-avatar>

      <nuxt-link to="/" class="text-decoration-none mr-10" style="color: inherit">
        <v-toolbar-title class="ml-2 d-flex">
          <div>
            <div class="subtitle-1 mb-n2">XRP</div>
            <div class="caption grey--text">finance</div>
          </div>
        </v-toolbar-title>
      </nuxt-link>

      <client-only>
        <div v-if="$vuetify.breakpoint.smAndUp" class="text-no-wrap">
          <nuxt-link
            v-for="(link, i) in links"
            :key="i"
            active-class="underline-glow-active"
            class="underline-glow-hover white--text mr-4 pb-2 px-1"
            tile
            text
            :to="link.to"
            >{{ link.name }}
          </nuxt-link>
        </div>
      </client-only>

      <v-spacer />

      <client-only>
        <!-- <global-search /> -->
        <wallet-connector class="hidden-sm-and-down" />
        <!-- <enhanced-xrp-wallet-connector class="hidden-sm-and-down" /> -->
      </client-only>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <nuxt />
      </v-container>
    </v-main>

    <v-navigation-drawer v-model="drawer" fixed temporary>
      <v-list elevation="0">
        <v-list-item-group color="primary">
          <v-list-item v-for="(item, index) in links" :key="index" :to="item.to">
            <v-list-item-title>{{ item.name }}</v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <client-only>
      <!-- <enhanced-wallet-select-dialog /> -->
      <main-footer />
    </client-only>
  </v-app>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useStore } from '@nuxtjs/composition-api'
import WalletConnector from '~/components/common/WalletConnector.vue'
import { useInitTheme } from '~/composables/useInitTheme'
// import EnhancedWalletSelectDialog from '~/components/common/EnhancedWalletSelectDialog.vue'
import { State } from '~/types/state'
import MainFooter from '~/components/common/ui/footers/MainFooter.vue'
// import GlobalSearch from '~/components/common/GlobalSearch.vue'
// import EnhancedXrpWalletConnector from '~/components/common/EnhancedXrpWalletConnector.vue'

export default defineComponent({
  components: {
    // EnhancedXrpWalletConnector,
    // GlobalSearch,
    MainFooter,
    // EnhancedWalletSelectDialog,
    WalletConnector,
  },
  setup() {
    // STATE
    const searchOverlay = ref(true)
    const notificationComponent = ref<Notification>()
    // const globalSearchComponentRef = ref<any>(null) // Removed since GlobalSearch is disabled
    const drawer = ref(false)
    const links = ref([
      { name: 'XRP-Explorer', to: '/xrp-explorer' },
      { name: 'XRP-Screener', to: '/xrp-screener' },
      { name: 'XRP-Balances', to: '/xrp-balances' },
      { name: 'XRP-Transactions', to: '/xrp-transactions' },
      { name: 'XRP-Terminal', to: '/xrp-terminal' },
      { name: 'XRP-AMM-Heatmap', to: '/xrp-amm-heatmap' },
      { name: 'XRP-Portfolio', to: '/xrp-portfolio' },
      { name: 'XRP-Trust-Lines', to: '/xrp-trust-lines' },
      { name: 'XRP-Token-Mints', to: '/xrp-token-mints' },
    ])

    // COMPOSABLE
    const { state } = useStore<State>()
    useInitTheme()

    // const onInitSearch = () => {
    //   try {
    //     globalSearchComponentRef.value.openDialog()
    //   } catch (e) {}
    // }

    // COMPUTED
    const imageUrl = computed(() => `/img/logo/evmfinance-logo.svg`)

    return {
      title: state.configs.title,
      links,
      notificationComponent,
      imageUrl,
      drawer,
      searchOverlay,
      // globalSearchComponentRef, // Removed since GlobalSearch is disabled
      // onInitSearch, // Removed since GlobalSearch is disabled
    }
  },
})
</script>
