<template>
  <div>
    <aave-actions
      v-if="pool"
      ref="aaveComponent"
      :type="`card`"
      :health-factor="healthFactor"
      :total-borrowed-usd="totalBorrowedUsd"
      :total-collateral-usd="totalCollateralUsd"
      :max-ltv="maxLtv"
      :pool-data="pool"
      :pool-action="aaveActionType"
      version="v3"
      @transaction-success="updatePortfolio()"
      @toggle-reserve-details="reserveDetailDialog = !reserveDetailDialog"
    />

    <v-col v-else class="pa-0">
      <v-card disabled tile outlined height="500">
        <v-card-title class="font-weight-bold subtitle-1 py-3 text-capitalize">
          <v-avatar size="26" class="mr-2"><v-img :src="$imageUrlBySymbol(`aave`)"></v-img></v-avatar>AAVE V2
        </v-card-title>
        <v-divider />
        <v-card-text>This token is not supported by AAVE</v-card-text>
      </v-card>
    </v-col>

    <v-dialog v-model="reserveDetailDialog" max-width="800">
      <v-card tile outlined class="pa-4">
        <v-card-title class="pa-0">
          Reserve Details<v-spacer />
          <v-icon @click="reserveDetailDialog = !reserveDetailDialog"> mdi-close </v-icon>
        </v-card-title>
        <aave-market-details :pool="pool" :show-balance-chart="false"></aave-market-details>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref, useRoute, watch } from '@nuxtjs/composition-api'
import AaveMarketDetails from '~/components/pools/AaveMarketDetails.vue'
import AaveActions from '~/components/pools/AaveActions.vue'
import usePortfolio, { PortfolioMap } from '~/composables/usePortfolio'
import useAavePools, { AavePoolModel, actionTypes } from '~/composables/useAavePools'
import UseAavePoolsStats from '~/composables/useAavePoolsStats'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'

export default defineComponent({
  components: {
    AaveMarketDetails,
    AaveActions,
  },

  setup() {
    const route = useRoute()
    const { walletReady, account, chainId } = inject(WEB3_PLUGIN_KEY) as Web3
    const portfolio = ref<PortfolioMap>({})
    const aaveActionType = ref<actionTypes>('deposit')
    const reserveDetailDialog = ref(false)

    const { aavePoolsData, loading } = useAavePools(chainId)

    const addresses = computed(() =>
      aavePoolsData.value.reduce((elem, item) => ({ ...elem, [item.id]: item.addresses }), {})
    )
    const { fetchPortfolioNew } = usePortfolio(addresses)
    const updatePortfolio = async () => (portfolio.value = await fetchPortfolioNew())
    const pools = computed(() => {
      const pools: AavePoolModel[] = []
      aavePoolsData.value.forEach((elem) => {
        elem.portfolio = portfolio.value[elem.id] || elem.portfolio
        pools.push(elem)
      })
      return pools
    })

    const { totalCollateralUsd, totalBorrowedUsd, healthFactor, maxLtv } = UseAavePoolsStats(pools)
    const pool = computed(() => pools.value.find((elem) => elem.symbol === route.value.params?.id ?? ''))

    // WATCHERS
    watch([aavePoolsData, walletReady, account, chainId], async () => {
      if (!loading.value) await updatePortfolio()
    })

    return {
      pool,
      pools,
      totalCollateralUsd,
      totalBorrowedUsd,
      healthFactor,
      maxLtv,
      aaveActionType,
      reserveDetailDialog,

      updatePortfolio,
    }
  },
})
</script>
