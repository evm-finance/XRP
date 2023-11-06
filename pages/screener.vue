<template>
  <div>
    <v-row no-gutters>
      <v-col>
        <v-btn-toggle v-model="selectedChainId" dark mandatory tile color="primary">
          <v-btn
            v-for="(item, index) in supportedNetworks"
            :key="index"
            :value="item.id"
            tile
            depressed
            height="36"
            color="transparent"
          >
            <v-avatar size="24" class="mr-2">
              <img :src="$imageUrlBySymbol(item.symbol)" alt="" @error="$setAltImageUrl" />
            </v-avatar>
            <span class="text-body-2 text-capitalize">{{ item.name }}</span>
          </v-btn>
        </v-btn-toggle>

        <v-divider vertical class="mx-4" />

        <v-btn-toggle v-if="selectedNetwork" v-model="selectedDexId" dark mandatory rounded color="primary">
          <v-btn
            v-for="(item, index) in selectedNetwork.dex"
            :key="index"
            :value="item.value"
            rounded
            depressed
            height="26"
            color="transparent"
          >
            <v-avatar size="16" class="mr-2">
              <img :src="$imageUrlBySymbol(item.symbol)" alt="" @error="$setAltImageUrl" />
            </v-avatar>
            <span class="text-body-2 text-capitalize">{{ item.name }}</span>
          </v-btn>
        </v-btn-toggle>
      </v-col>

      <v-col cols="3" class="pt-2">
        <div class="text-right">
          <custom-wallet-indicator selected-chain-id="ethereum" />
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-col lg="9" md="6">
        <screener-grid
          :search="searchString"
          :dex-id="selectedDexId"
          :network-id="selectedChainId"
          :aave-pools="aavePoolsMap"
        />
      </v-col>
      <v-col v-if="$vuetify.breakpoint.mdAndUp" lg="3" md="6">
        <alert-stream :network-id="selectedChainId" />
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref } from '@nuxtjs/composition-api'
import ScreenerGrid from '~/components/ScreenerGrid.vue'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import useAavePools, { AavePoolModel } from '~/composables/useAavePools'
import CustomWalletIndicator from '~/components/common/CustomWalletIndicator.vue'
import { Network } from '~/types/global'

export default defineComponent({
  components: {
    CustomWalletIndicator,
    ScreenerGrid,
  },
  setup() {
    const { allNetworks } = inject(WEB3_PLUGIN_KEY) as Web3
    const supportedNetworks = computed(() =>
      allNetworks.value.filter((elem) =>
        ['ethereum', 'polygon-pos', 'binance-smart-chain', 'avalanche'].includes(elem.id)
      )
    )
    const searchString = ref('')
    const selectedChainId = ref('ethereum')
    const selectedDexId = ref('uniswap_v3')
    const selectedNetwork = computed<Network | null>(
      () => allNetworks.value.find((elem) => elem.id === selectedChainId.value) ?? null
    )
    const chainId = computed(() => selectedNetwork.value?.chainIdentifier ?? null)
    const { aavePoolsData } = useAavePools(chainId)
    const aavePoolsMap = computed<{ [a: string]: AavePoolModel }>(() =>
      aavePoolsData.value.reduce((obj, item) => ({ ...obj, [item.addresses.address.toLowerCase()]: item }), {})
    )

    return {
      allNetworks,
      selectedDexId,
      searchString,
      selectedChainId,
      selectedNetwork,
      aavePoolsMap,
      supportedNetworks,
    }
  },
  head: {},
})
</script>

<style scoped></style>
