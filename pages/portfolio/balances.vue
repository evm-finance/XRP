<template>
  <v-row no-gutters justify="center">
    <v-col cols="12" md="10">
      <v-row no-gutters>
        <v-col>
          <h1 class="text-h4">
            DeFi Balances
            <info-tooltip :text="messages.tooltips.balanceHeaderDesc" />
            <nuxt-link :to="`/portfolio/transactions?wallet=${account}`" class="text-caption">
              Navigate to Transaction History
            </nuxt-link>
          </h1>
        </v-col>

        <v-col cols="3" class="pt-2">
          <div class="text-right">
            <custom-wallet-indicator selected-chain-id="ethereum" />
          </div>
        </v-col>
      </v-row>

      <v-row v-if="!isWalletReady" justify="center" no-gutters class="pt-3">
        <v-col cols="12" sm="10" lg="6">
          <connect-wallet-memo />
        </v-col>
      </v-row>

      <v-row v-else justify="center">
        <v-col>
          <v-row>
            <v-col lg="4" cols="12">
              <v-card class="py-2 px-4" height="336" tile outlined>
                <v-dialog v-model="loading" width="450">
                  <v-card tile outlined class="pa-6">
                    <v-row no-gutters>
                      <v-col cols="3">
                        <v-progress-circular color="pink" indeterminate size="64" />
                      </v-col>
                      <v-col class="text-left">
                        <div>Hang Tight!</div>
                        <div class="grey--text">Scouring the Blockchain for your Riches.</div>
                        <div class="grey--text">It may take a minute.</div>
                      </v-col>
                    </v-row>
                  </v-card>
                </v-dialog>
                <v-skeleton-loader v-if="loading" type="heading,table-tbody,table-tbody, table-tbody" height="280" />
                <div v-else>
                  <h4 :class="['text-subtitle-1 text-truncate pink--text font-weight-medium']">Wallet Balance</h4>
                  <div
                    class="d-inline-block text-truncate text-h3"
                    v-text="$f(totalBalance, { pre: '$ ', minDigits: 2 })"
                  />
                  <v-simple-table class="mt-3">
                    <template #default>
                      <tbody class="text-subtitle-1 text-no-wrap">
                        <tr v-for="(elem, i) in stats" :key="i">
                          <td>
                            <div class="text-no-wrap overflow-x-hidden">
                              <v-avatar size="20" class="mr-2">
                                <v-img
                                  :src="$imageUrlBySymbol(elem.symbol)"
                                  :lazy-src="$imageUrlBySymbol(elem.symbol)"
                                />
                              </v-avatar>
                              {{ elem.name }}
                            </div>
                          </td>
                          <td class="grey--text" v-text="$f(elem.total, { pre: '$ ', minDigits: 2 })" />
                        </tr>
                        <tr v-if="xrpTotalBalance > 0">
                          <td>
                            <div class="text-no-wrap overflow-x-hidden">
                              <v-avatar size="20" class="mr-2">
                                <v-img
                                  :src="$imageUrlBySymbol('xrp')"
                                  :lazy-src="$imageUrlBySymbol('xrp')"
                                />
                              </v-avatar>
                              XRP Ledger
                            </div>
                          </td>
                          <td class="grey--text" v-text="$f(xrpTotalBalance, { pre: '$ ', minDigits: 2 })" />
                        </tr>
                      </tbody>
                    </template>
                  </v-simple-table>
                </div>
              </v-card>
            </v-col>
            <v-col class="hidden-sm-and-down">
              <v-card tile outlined height="100%">
                <v-skeleton-loader v-if="loading" type="image, image" height="336" />
                <client-only>
                  <balances-chart v-if="!loading" :balances="balanceData" />
                </client-only>
              </v-card>
            </v-col>
          </v-row>
          <v-row v-if="loading">
            <v-col v-for="i in 4" :key="i" cols="6">
              <v-card height="450" tile outlined class="pa-2">
                <v-skeleton-loader :loading="loading" type="heading,table-tbody,table-tbody" height="420" />
              </v-card>
            </v-col>
          </v-row>

          <v-row v-show="!loading">
            <v-col v-for="balance in balanceData" :key="balance.chainId" lg="6" md="12">
              <portfolio-balance-grid :data="balance" />
            </v-col>
            
            <v-col lg="6" md="12">
              <xrp-balance-widget />
            </v-col>
          </v-row>

          <!--          <balance-protocols v-show="!loading" :balances="balanceData" />-->
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, inject, ref, useRoute } from '@nuxtjs/composition-api'
import usePortfolioBalances from '~/composables/usePortfolioBalances'
import PortfolioBalanceGrid from '~/components/portfolio/PortfolioBalanceGrid.vue'
import XrpBalanceWidget from '~/components/portfolio/XrpBalanceWidget.vue'
import BalancesChart from '~/components/portfolio/BalancesChart.vue'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
// import BalanceProtocols from '~/components/portfolio/BalanceProtocols.vue'
import ConnectWalletMemo from '~/components/common/ConnectWalletMemo.vue'
import { useMetaTags } from '~/composables/useMetaTags'
import InfoTooltip from '~/components/common/ui/InfoTooltip.vue'
import { messages } from '~/constants/messages'
import CustomWalletIndicator from '~/components/common/CustomWalletIndicator.vue'

export default defineComponent({
  components: {
    CustomWalletIndicator,
    InfoTooltip,
    ConnectWalletMemo,
    // BalanceProtocols,
    BalancesChart,
    PortfolioBalanceGrid,
    XrpBalanceWidget,
  },
  setup() {
    // COMPOSABLES
    const { getNetworkByChainNumber } = inject(WEB3_PLUGIN_KEY) as Web3
    const { loading, balanceData, error, totalBalance, isWalletReady, account } = usePortfolioBalances()

    // COMPUTED
    const stats = computed(() =>
      balanceData.value.map((item) => ({
        name: getNetworkByChainNumber(item.chainId)?.name ?? '',
        symbol: getNetworkByChainNumber(item.chainId)?.symbol ?? '',
        total: item.items.reduce((n, { quote }) => n + quote, 0),
      }))
    )

    // Mock XRP total balance for now - in real implementation this would come from XRP widget
    const xrpTotalBalance = ref(0)

    // META TAGS
    useMetaTags('balances', useRoute().value.path)
    const dialog = ref(true)
    return { 
      //methods
      //data
      loading, balanceData, error, stats, totalBalance, isWalletReady, messages, dialog, account, xrpTotalBalance 
    }
  },
  head: {},
})
</script>
<style scoped>
>>> .v-data-table > .v-data-table__wrapper > table > tbody > tr > td {
  height: 32px;
}
</style>
