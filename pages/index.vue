<template>
  <v-row no-gutters justify="center">
    <!--    <vue-particles />-->

    <v-col cols="12" style="z-index: 1">
      <v-row justify="center" class="text-center">
        <v-col cols="12" class="mt-6">
          <h1 class="text-h2 font-weight-regular">
            EVM Finance <br />
            <span class="text-h4 font-weight-light mt-3 grey--text"> Unified Interface for Decentralized Finance </span>
          </h1>
        </v-col>

        <v-col cols="12">
          <p class="text-h4 font-weight-light">
            Token Pages • Aave Professional Interface • Multi Chain Balances • Advanced Analytics
          </p>
        </v-col>
      </v-row>
    </v-col>

    <v-col cols="12" style="z-index: 1">
      <!--      <defi-node-tree />-->
    </v-col>

    <v-col cols="11" xl="7" lg="10" class="text-center mt-16 pt-16" style="z-index: 1">
      <v-row>
        <v-col cols="12">
          <h2 class="text-h4 text-left">Roadmap</h2>
        </v-col>
        <v-col cols="12" class="justify-center">
          <v-row justify="center">
            <v-col v-for="item in items" :key="item.header" cols="12" sm="6" md="4">
              <v-card class="pa-4">
                <v-avatar size="40" class="mt-5">
                  <v-icon size="30" color="primary">{{ item.icon }}</v-icon>
                </v-avatar>

                <p class="text-h5 pt-5 px-5" v-text="item.header" />
                <p class="pt-5 px-5 grey--text" v-text="item.desc" />
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useStore } from '@nuxtjs/composition-api'

import { State } from '~/types/state'
import { useMetaTags } from '~/composables/useMetaTags'
// import OutlineGlow from '~/components/common/ui/custom/OutlineGlow.vue'
// import DefiNodeTree from '~/components/common/DefiNodeTree.vue'

export default defineComponent({
  // components: { DefiNodeTree, OutlineGlow },
  setup() {
    // STATE
    const items = ref([
      {
        color: 'primary lighten-2',
        icon: 'mdi-tablet-dashboard',
        header: 'Aave Professional Trader',
        desc: 'Easier Interface using Aave smart contracts',
      },
      {
        color: 'red lighten-2',
        icon: 'mdi-wallet-outline',
        header: 'Asset Management',
        desc: 'Full Wallet Balance, Chain Balances, Token Balance',
      },

      {
        color: 'orange lighten-2',
        icon: 'mdi-desktop-mac-dashboard',
        header: 'Token Pages',
        desc: 'Balances, investment options, and metrics on a single screen',
      },
      {
        color: 'red lighten-2',
        icon: 'mdi-history',
        header: 'Transaction History',
        desc: 'Best in class transaction displays',
      },
      {
        color: 'primary lighten-2',
        icon: 'mdi-transit-connection-variant',
        header: 'Self Custody',
        desc: 'Support for Ethereum Mainnet, Binance Smart Chain, Matic, Fantom Opera, Avalanche networks and more',
      },
      {
        color: 'red lighten-2',
        icon: 'mdi-rocket-launch-outline',
        header: 'NFT Launch',
        desc: 'Check this site for future updates',
      },
    ])

    // COMPOSABLE
    const store = useStore<State>()

    // COMPUTED
    const ui = computed(() => store.state.ui)
    const theme = computed(() => store.state.ui.theme)

    return { items, ui, theme }
  },
  head() {
    return useMetaTags('homepage')
  },
})
</script>

<style lang="scss">
.vue-typer .custom.char.typed {
  color: #e91e63ff;
  align-items: center;
}

.vue-typer .custom.caret {
  animation: rocking 1s ease-in-out 0s infinite;
}

.vue-typer .custom.caret {
  width: 5px;
  background-color: #e91e63ff;
  margin-right: 50px;
}
</style>
