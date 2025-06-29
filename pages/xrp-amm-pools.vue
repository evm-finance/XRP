<template>
  <div>
    <v-row no-gutters justify="center">
      <v-col cols="12" md="10">
        <v-row justify="center">
          <v-col cols="12">
            <h1 class="text-h4">XRP AMM Pools</h1>
          </v-col>
        </v-row>
        
        <!-- XRP Balance Widget -->
        <v-row justify="center" class="mb-4">
          <v-col cols="12" lg="6">
            <xrp-balance-widget />
          </v-col>
        </v-row>
        
        <v-row justify="center">
          <v-col md="12">
            <v-card tile outlined height="1330">
              <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />
              <client-only>
                <v-data-table
                  v-if="!loading"
                  hide-default-footer
                  :headers="cols"
                  :items="poolsDataFormatted"
                  :items-per-page="25"
                  class="elevation-0 row-height-50"
                  mobile-breakpoint="0"
                >
                  <template #[`item.pair`]="{ item }">
                    <div class="text-no-wrap overflow-x-hidden">
                      <div class="d-flex align-center">
                        <v-avatar size="20" class="mr-2">
                          <v-img
                            :src="$imageUrlBySymbol(item.token1.symbol.toLowerCase())"
                            :lazy-src="$imageUrlBySymbol(item.token1.symbol.toLowerCase())"
                          />
                        </v-avatar>
                        <span class="mr-1">{{ item.token1.symbol }}</span>
                        <span class="mx-1">/</span>
                        <v-avatar size="20" class="mr-2">
                          <v-img
                            :src="$imageUrlBySymbol(item.token2.symbol.toLowerCase())"
                            :lazy-src="$imageUrlBySymbol(item.token2.symbol.toLowerCase())"
                          />
                        </v-avatar>
                        <nuxt-link
                          class="text-decoration-none white--text"
                          :to="`/xrp-amm-pools/${item.id}`"
                        >
                          {{ item.token2.symbol }}
                        </nuxt-link>
                      </div>
                      <div class="text-caption grey--text">
                        {{ item.token2.issuer ? `${item.token2.issuer.slice(0, 8)}...` : 'Native' }}
                      </div>
                    </div>
                  </template>

                  <template #[`item.token2IssuerShort`]="{ item }">
                    <div class="d-flex align-center">
                      <span class="font-family-mono">{{ item.token2IssuerShort }}</span>
                      <v-btn
                        v-if="item.token2.issuer"
                        icon
                        x-small
                        class="ml-1"
                        @click="copyToClipboard(item.token2.issuer)"
                      >
                        <v-icon size="16">mdi-content-copy</v-icon>
                      </v-btn>
                    </div>
                  </template>

                  <template #[`item.liquidity`]="{ item }">
                    <span class="font-weight-medium">${{ formatNumber(item.liquidity) }}</span>
                  </template>

                  <template #[`item.volume24h`]="{ item }">
                    <span class="font-weight-medium">${{ formatNumber(item.volume24h) }}</span>
                  </template>

                  <template #[`item.apr`]="{ item }">
                    <span class="font-weight-medium text-success">{{ formatPercentage(item.apr) }}%</span>
                  </template>

                  <template #[`item.fee`]="{ item }">
                    <span class="font-weight-medium">{{ formatFee(item.fee) }}%</span>
                  </template>

                  <template #[`item.priceChange24h`]="{ item }">
                    <span :class="item.priceChange24h >= 0 ? 'text-success' : 'text-error'">
                      {{ item.priceChange24h >= 0 ? '+' : '' }}{{ formatPercentage(item.priceChange24h) }}%
                    </span>
                  </template>

                  <template #[`item.actions`]="{ item }">
                    <v-btn
                      small
                      text
                      color="primary"
                      @click="viewPool(item)"
                    >
                      View
                    </v-btn>
                  </template>
                </v-data-table>
              </client-only>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, inject } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'
import { XRPAmmPoolsGQL } from '~/apollo/queries'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import XrpBalanceWidget from '~/components/portfolio/XrpBalanceWidget.vue'

interface XRPAmmPool {
  id: string
  token1: { symbol: string; name: string; icon: string }
  token2: { symbol: string; name: string; icon: string; issuer?: string }
  liquidity: number
  volume24h: number
  fee: number
  apr: number
  priceChange24h: number
}

export default defineComponent({
  components: {
    XrpBalanceWidget,
  },
  setup() {
    const { $f } = useContext()
    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    
    const loading = ref(true)
    const poolsRawData = ref<XRPAmmPool[]>([])

    // Enhanced GraphQL query with logging
    const { onResult, onError } = useLoggedQuery(
      XRPAmmPoolsGQL,
      null,
      { fetchPolicy: 'no-cache', pollInterval: 30000 }
    )

    // Handle query results
    onResult((queryResult: any) => {
      if (queryResult.data?.xrpAmmPools) {
        const pools = queryResult.data.xrpAmmPools
        
        // Transform GraphQL data to component format
        const transformedPools: XRPAmmPool[] = pools.map((pool: any) => ({
          id: pool.id || '',
          token1: {
            symbol: pool.token1?.symbol || 'XRP',
            name: pool.token1?.name || 'XRP',
            icon: pool.token1?.icon || '⚡',
            issuer: pool.token1?.issuer || ''
          },
          token2: {
            symbol: pool.token2?.symbol || '',
            name: pool.token2?.name || '',
            icon: pool.token2?.icon || '🪙',
            issuer: pool.token2?.issuer || ''
          },
          liquidity: pool.liquidity || 0,
          volume24h: pool.volume24h || 0,
          fee: pool.fee || 0,
          apr: pool.apr || 0,
          priceChange24h: pool.priceChange24h || 0,
          token1Balance: pool.token1Balance || 0,
          token2Balance: pool.token2Balance || 0,
          totalSupply: pool.totalSupply || 0,
          createdAt: pool.createdAt || '',
          lastUpdated: pool.lastUpdated || ''
        }))
        
        poolsRawData.value = transformedPools
      } else {
        // No pools found
        poolsRawData.value = []
      }
      loading.value = false
    })

    // Add error handling
    onError((error: any) => {
      console.error('GraphQL Error in xrp-amm-pools:', error)
      loading.value = false
      // You could add error state handling here
    })

    const poolsDataFormatted = computed(() =>
      poolsRawData.value.map((elem) => ({
        ...elem,
        token2IssuerShort: elem.token2.issuer 
          ? `${elem.token2.issuer.slice(0, 10)}........${elem.token2.issuer.slice(
              elem.token2.issuer.length - 10,
              elem.token2.issuer.length
            )}`
          : 'Native',
        liquidityFormatted: $f(elem.liquidity, { minDigits: 0, after: '' }),
        volume24hFormatted: $f(elem.volume24h, { minDigits: 0, after: '' }),
      }))
    )

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // You could add a toast notification here
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }

    const formatNumber = (value: number): string => {
      return $f(value, { minDigits: 0, after: '' })
    }

    const formatPercentage = (value: number): string => {
      return value.toFixed(1)
    }

    const formatFee = (value: number): string => {
      return (value * 100).toFixed(3)
    }

    const viewPool = (pool: XRPAmmPool) => {
      // Navigate to pool detail page
      window.location.href = `/xrp-amm-pools/${pool.id}`
    }

    const cols = computed(() => {
      return [
        {
          text: 'Pool',
          align: 'left',
          value: 'pair',
          width: '200',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate'],
          sortable: true,
        },
        {
          text: 'Token 2 Issuer',
          align: 'left',
          value: 'token2IssuerShort',
          width: '150',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
        },
        {
          text: 'Liquidity',
          align: 'right',
          value: 'liquidity',
          width: '120',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: '24h Volume',
          align: 'right',
          value: 'volume24h',
          width: '120',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'APR',
          align: 'right',
          value: 'apr',
          width: '80',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Fee',
          align: 'right',
          value: 'fee',
          width: '80',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: '24h Change',
          align: 'right',
          value: 'priceChange24h',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Actions',
          align: 'center',
          value: 'actions',
          width: '80',
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
      ]
    })

    return {
      loading,
      cols,
      poolsDataFormatted,
      copyToClipboard,
      formatNumber,
      formatPercentage,
      formatFee,
      viewPool,
    }
  },
  head: {},
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style> 