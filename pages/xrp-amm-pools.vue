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
import { computed, defineComponent, ref, useContext } from '@nuxtjs/composition-api'
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
    const loading = ref(true)
    const poolsRawData = ref<XRPAmmPool[]>([])

    // Mock AMM data - in real implementation this would come from XRPL AMM API
    const mockAmmData = [
      {
        id: 'amm_1',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'USDC', name: 'USD Coin', icon: '🪙', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 2500000,
        volume24h: 1800000,
        fee: 0.003,
        apr: 12.5,
        priceChange24h: 2.3,
      },
      {
        id: 'amm_2',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'USDT', name: 'Tether', icon: '💎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 1800000,
        volume24h: 1200000,
        fee: 0.003,
        apr: 15.2,
        priceChange24h: -1.8,
      },
      {
        id: 'amm_3',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'BTC', name: 'Bitcoin', icon: '₿', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 1200000,
        volume24h: 800000,
        fee: 0.003,
        apr: 18.7,
        priceChange24h: 4.1,
      },
      {
        id: 'amm_4',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'ETH', name: 'Ethereum', icon: 'Ξ', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 900000,
        volume24h: 600000,
        fee: 0.003,
        apr: 22.1,
        priceChange24h: -0.5,
      },
      {
        id: 'amm_5',
        token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
        token2: { symbol: 'USDC', name: 'USD Coin', icon: '🪙', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 800000,
        volume24h: 400000,
        fee: 0.001,
        apr: 8.3,
        priceChange24h: 0.1,
      },
      {
        id: 'amm_6',
        token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
        token2: { symbol: 'USDT', name: 'Tether', icon: '💎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 600000,
        volume24h: 300000,
        fee: 0.001,
        apr: 9.8,
        priceChange24h: -0.2,
      },
      {
        id: 'amm_7',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'SOL', name: 'Solana', icon: '◎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 500000,
        volume24h: 250000,
        fee: 0.003,
        apr: 25.4,
        priceChange24h: 6.7,
      },
      {
        id: 'amm_8',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'ADA', name: 'Cardano', icon: '₳', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 400000,
        volume24h: 200000,
        fee: 0.003,
        apr: 28.9,
        priceChange24h: -2.1,
      },
    ]

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

    // Initialize with mock data
    const initializeData = () => {
      poolsRawData.value = mockAmmData
      loading.value = false
    }

    // Initialize data on mount
    initializeData()

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