<template>
  <v-row justify="center">
    <v-col lg="10" md="12">
      <v-row v-if="!loading">
        <v-col cols="12" lg="4">
          <v-row>
            <v-col sm="12" md="6" lg="12">
              <coin-profile
                height="364"
                :symbol="tokenData.symbolName"
                :qc-key="tokenData.qcKey"
                :aave-symbol="tokenData.aaveSymbol"
                :is-aave-token="tokenData.isAaveToken"
                :is-q-c-token="tokenData.isQCToken"
                :contract-address="contractAddress"
                :website-links="tokenData.websiteUrl"
                :bit-bucket-repos="tokenData.bitbucketRepos"
                :git-hub-repos="tokenData.githubRepos"
                :explorer-urls="tokenData.explorerUrls"
                :telegram-channel-id="tokenData.telegramChannelId"
                :twitter-url="tokenData.twitterUrl"
                :discord-channel-id="tokenData.discordChannelId"
                :subreddit-url="tokenData.subredditUrl"
                :facebook-url="tokenData.facebookUrl"
                :coin-description="tokenData.coinDescription"
              />
            </v-col>

            <v-col sm="12" md="6" lg="12">
              <token-balances
                height="100%"
                :balances="tokenData.balances"
                :price-usd="tokenData.price.priceUsd"
                :contract="contractAddress"
              />
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12" lg="8">
          <v-row>
            <v-col cols="12">
              <coin-metrics
                :price-usd="tokenData.price.priceUsd"
                :price-eth="tokenData.price.priceEth"
                :price24h="tokenData.price24h"
                :market-cap="tokenData.marketcap"
                :circ-supply="tokenData.circulatingSupply"
              />
            </v-col>

            <v-col cols="12" md="6">
              <token-aave-assets />
            </v-col>

            <v-col cols="12" md="6">
              <evm-verse-swap
                v-if="isVerse"
                height="500"
                :in-token="{
                  symbol: tokenData.qcKey,
                  name: tokenData.symbolName,
                  chainId: tokenData.chainId,
                  address: tokenData.contractAddress,
                  decimals: tokenData.decimals,
                }"
                width="100%"
              />
              <evm-swap
                v-else
                :in-token="{
                  symbol: tokenData.qcKey,
                  name: tokenData.symbolName,
                  chainId: tokenData.chainId,
                  address: tokenData.contractAddress,
                  decimals: tokenData.decimals,
                }"
                height="500"
                width="100%"
              />
            </v-col>
          </v-row>
        </v-col>

        <v-col v-if="!tokenTransactions.loading.value && tokenData" cols="12">
          <token-transactions-grid
            :account="tokenTransactions.account.value"
            :contract-name="tokenData.symbolName"
            :transactions-data="tokenTransactions.transfersData.value"
          />

          <v-pagination
            v-model="tokenTransactions.currentPage.value"
            class="mt-4"
            :total-visible="tokenTransactions.pagination.visible"
            :length="tokenTransactions.pagination.total"
          ></v-pagination>
        </v-col>

        <v-col cols="12">
          <v-card tile outlined height="600">
            <token-prices-chart :coin-gecko-id="tokenData.coinGeckoID" />
          </v-card>
        </v-col>
      </v-row>

      <v-row v-if="loading">
        <v-col cols="12" lg="4">
          <v-row>
            <v-col>
              <v-card height="364" tile outlined>
                <v-skeleton-loader type="article,list-item-three-line@4" height="364" />
              </v-card>
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-card height="230" tile outlined>
                <v-skeleton-loader type="table-heading,divider,table-tbody" height="230" />
              </v-card>
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12" lg="8">
          <v-row>
            <v-col>
              <v-card height="102" tile outlined class="px-3 pt-2">
                <v-skeleton-loader type="table-heading,text" height="102" />
              </v-card>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-card tile outlined height="492">
                <v-skeleton-loader type="table-heading,divider,image@3" height="492" />
              </v-card>
            </v-col>
            <v-col cols="6">
              <v-card tile outlined height="492">
                <v-skeleton-loader type="table-heading,divider,image@3" height="492" />
              </v-card>
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12">
          <v-card tile outlined height="447.59">
            <v-skeleton-loader type="table-heading,divider,table-tbody" tile height="447.59" />
          </v-card>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, useRoute } from '@nuxtjs/composition-api'
import useToken from '~/composables/useToken'
import CoinProfile from '~/components/token-details/CoinProfile.vue'
import CoinMetrics from '~/components/token-details/CoinMetrics.vue'
import TokenBalances from '~/components/token-details/TokenBalances.vue'
import TokenPricesChart from '~/components/token-details/TokenPricesChart.vue'
import { useMetaTags } from '~/composables/useMetaTags'
import TokenTransactionsGrid from '~/components/transactions/TokenTransactionsGrid.vue'
import useTokenTransactions from '~/composables/useTokenTransactions'
import EvmSwap from '~/components/trading/EvmSwap.vue'
import TokenAaveAssets from '~/components/token-details/TokenAaveAssets.vue'
import EvmVerseSwap from '~/components/trading/EvmVerseSwap.vue'

export default defineComponent({
  components: {
    EvmVerseSwap,
    TokenAaveAssets,
    EvmSwap,
    CoinMetrics,
    TokenTransactionsGrid,
    TokenPricesChart,
    TokenBalances,

    CoinProfile,
  },
  setup() {
    const { params, query } = useRoute().value
    const isVerse = computed(() => params.id === 'VERSE')

    // Token Data
    const { tokenData, loading, contractAddress } = useToken()

    // Token Transactions
    const tokenTransactions = useTokenTransactions()

    useMetaTags('tokenPage', `/token/${params.id}`, query.name as string, params.id)

    return {
      isVerse,
      tokenData,
      loading,
      contractAddress,
      tokenTransactions,
    }
  },
  head: {},
})
</script>
