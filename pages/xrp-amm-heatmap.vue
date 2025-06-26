<template>
  <div>
    <v-container fluid class="pa-0">
      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-title class="py-3">
              <div class="d-flex align-center">
                <v-icon large class="mr-3" color="primary">💧</v-icon>
                <div>
                  <h1 class="text-h4 font-weight-bold">XRP AMM Heatmap</h1>
                  <p class="text-subtitle-1 grey--text mb-0">
                    Visualize liquidity across XRP Automated Market Makers
                  </p>
                </div>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>

      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-text>
              <v-row>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold">{{ totalPools }}</div>
                    <div class="text-caption grey--text">Total AMM Pools</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-success">${{ totalLiquidity }}</div>
                    <div class="text-caption grey--text">Total Liquidity</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-info">${{ totalVolume }}</div>
                    <div class="text-caption grey--text">24h Volume</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-warning">{{ avgApr }}%</div>
                    <div class="text-caption grey--text">Average APR</div>
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row no-gutters>
        <v-col>
          <xrp-amm-heatmap :height="heatmapHeight" :user-can-access-trend="true" />
        </v-col>
      </v-row>

      <v-row no-gutters class="mt-4">
        <v-col>
          <v-card tile outlined>
            <v-card-title class="py-3">
              <span class="text-h6">Top AMM Pools by Liquidity</span>
              <v-spacer />
              <v-btn
                text
                color="primary"
                @click="refreshData"
                :loading="loading"
              >
                <v-icon left>mdi-refresh</v-icon>
                Refresh
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-data-table
                :headers="tableHeaders"
                :items="topPools"
                :loading="loading"
                :items-per-page="10"
                class="elevation-0"
                hide-default-footer
              >
                <template #item.pair="{ item }">
                  <div class="d-flex align-center">
                    <div class="d-flex align-center mr-2">
                      <span class="mr-1">{{ item.token1.icon }}</span>
                      <span class="mr-1">/</span>
                      <span>{{ item.token2.icon }}</span>
                    </div>
                    <div>
                      <div class="font-weight-medium">{{ item.pair }}</div>
                      <div class="text-caption grey--text">{{ item.token2.issuer ? item.token2.issuer.slice(0, 8) + '...' : 'Native' }}</div>
                    </div>
                  </div>
                </template>
                <template #item.liquidity="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.liquidity) }}</span>
                </template>
                <template #item.volume24h="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.volume24h) }}</span>
                </template>
                <template #item.apr="{ item }">
                  <span class="font-weight-medium text-success">{{ formatPercentage(item.apr) }}%</span>
                </template>
                <template #item.fee="{ item }">
                  <span class="font-weight-medium">{{ formatFee(item.fee) }}%</span>
                </template>
                <template #item.priceChange24h="{ item }">
                  <span :class="item.priceChange24h >= 0 ? 'text-success' : 'text-error'">
                    {{ item.priceChange24h >= 0 ? '+' : '' }}{{ formatPercentage(item.priceChange24h) }}%
                  </span>
                </template>
                <template #item.actions="{ item }">
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
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref } from '@nuxtjs/composition-api'
import XrpAmmHeatmap from '~/components/xrp/XrpAmmHeatmap.vue'

export default defineComponent({
  components: {
    XrpAmmHeatmap,
  },
  setup() {
    const { $f } = useContext()
    const loading = ref(false)
    const poolsData = ref<any[]>([])
    
    // Heatmap height
    const heatmapHeight = ref(600)
    
    // Table headers
    const tableHeaders = [
      { text: 'Pool', value: 'pair', sortable: false },
      { text: 'Liquidity', value: 'liquidity', align: 'right' },
      { text: '24h Volume', value: 'volume24h', align: 'right' },
      { text: 'APR', value: 'apr', align: 'right' },
      { text: 'Fee', value: 'fee', align: 'right' },
      { text: '24h Change', value: 'priceChange24h', align: 'right' },
      { text: 'Actions', value: 'actions', sortable: false, align: 'center' },
    ]

    // Mock AMM data
    const ammData = ref([
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
      {
        id: 'amm_9',
        token1: { symbol: 'rLUSD', name: 'rLUSD', icon: '💵' },
        token2: { symbol: 'BTC', name: 'Bitcoin', icon: '₿', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 350000,
        volume24h: 150000,
        fee: 0.001,
        apr: 12.6,
        priceChange24h: 1.8,
      },
      {
        id: 'amm_10',
        token1: { symbol: 'XRP', name: 'XRP', icon: '⚡' },
        token2: { symbol: 'DOT', name: 'Polkadot', icon: '●', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
        liquidity: 300000,
        volume24h: 120000,
        fee: 0.003,
        apr: 32.1,
        priceChange24h: 3.4,
      },
    ])

    // Computed properties
    const topPools = computed(() => {
      return ammData.value
        .map(item => ({
          ...item,
          pair: `${item.token1.symbol}/${item.token2.symbol}`,
        }))
        .sort((a, b) => b.liquidity - a.liquidity)
        .slice(0, 20)
    })

    const totalPools = computed(() => ammData.value.length)
    const totalLiquidity = computed(() => 
      $f(ammData.value.reduce((sum, pool) => sum + pool.liquidity, 0), { 
        minDigits: 0, 
        after: '' 
      })
    )
    const totalVolume = computed(() => 
      $f(ammData.value.reduce((sum, pool) => sum + pool.volume24h, 0), { 
        minDigits: 0, 
        after: '' 
      })
    )
    const avgApr = computed(() => {
      const avg = ammData.value.reduce((sum, pool) => sum + pool.apr, 0) / ammData.value.length
      return $f(avg, { minDigits: 1, after: '' })
    })

    // Methods
    const formatNumber = (num: number) => {
      return $f(num, { minDigits: 0, after: '' })
    }

    const formatPercentage = (percentage: number) => {
      return $f(percentage, { minDigits: 2, after: '' })
    }

    const formatFee = (fee: number) => {
      return $f(fee * 100, { minDigits: 3, after: '' })
    }

    const viewPool = (pool: any) => {
      // Navigate to pool details page (to be implemented)
      console.log('View pool:', pool.id)
    }

    const refreshData = async () => {
      loading.value = true
      // Simulate refresh
      setTimeout(() => {
        loading.value = false
      }, 1000)
    }

    onMounted(() => {
      // Set responsive height
      if (process.client) {
        heatmapHeight.value = window.innerHeight * 0.6
      }
    })

    return {
      loading,
      poolsData,
      topPools,
      totalPools,
      totalLiquidity,
      totalVolume,
      avgApr,
      heatmapHeight,
      tableHeaders,
      formatNumber,
      formatPercentage,
      formatFee,
      viewPool,
      refreshData,
    }
  },
  head() {
    return {
      title: 'XRP AMM Heatmap - EVM Finance',
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: 'Visualize liquidity across XRP Automated Market Makers with our interactive heatmap. Track AMM pools, liquidity, and yield opportunities.',
        },
      ],
    }
  },
})
</script> 