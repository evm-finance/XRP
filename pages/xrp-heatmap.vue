<template>
  <div>
    <v-container fluid class="pa-0">
      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-title class="py-3">
              <div class="d-flex align-center">
                <v-icon large class="mr-3" color="primary">🔥</v-icon>
                <div>
                  <h1 class="text-h4 font-weight-bold">XRP Token Heatmap</h1>
                  <p class="text-subtitle-1 grey--text mb-0">
                    Visualize XRP token performance across the XRP Ledger
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
                    <div class="text-h6 font-weight-bold">{{ totalTokens }}</div>
                    <div class="text-caption grey--text">Total Tokens</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-success">{{ gainersCount }}</div>
                    <div class="text-caption grey--text">Gainers (24h)</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-error">{{ losersCount }}</div>
                    <div class="text-caption grey--text">Losers (24h)</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold">${{ totalMarketCap }}</div>
                    <div class="text-caption grey--text">Total Market Cap</div>
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row no-gutters>
        <v-col>
          <xrp-token-heatmap :height="heatmapHeight" :user-can-access-trend="true" />
        </v-col>
      </v-row>

      <v-row no-gutters class="mt-4">
        <v-col>
          <v-card tile outlined>
            <v-card-title class="py-3">
              <span class="text-h6">Top XRP Tokens by Market Cap</span>
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
                :items="topTokens"
                :loading="loading"
                :items-per-page="10"
                class="elevation-0"
                hide-default-footer
              >
                <template #item.tokenName="{ item }">
                  <div class="d-flex align-center">
                    <span class="mr-2">{{ item.icon }}</span>
                    <div>
                      <div class="font-weight-medium">{{ item.tokenName }}</div>
                      <div class="text-caption grey--text">{{ item.currency }}</div>
                    </div>
                  </div>
                </template>
                <template #item.price="{ item }">
                  <span class="font-weight-medium">${{ formatPrice(item.price) }}</span>
                </template>
                <template #item.price24h="{ item }">
                  <span :class="item.price24h >= 0 ? 'text-success' : 'text-error'">
                    {{ item.price24h >= 0 ? '+' : '' }}{{ formatPercentage(item.price24h) }}%
                  </span>
                </template>
                <template #item.marketcap="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.marketcap) }}</span>
                </template>
                <template #item.volume24h="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.volume24h) }}</span>
                </template>
                <template #item.issuer="{ item }">
                  <div class="d-flex align-center">
                    <span class="text-caption font-family-mono">{{ formatAddress(item.issuer) }}</span>
                    <v-btn
                      icon
                      x-small
                      class="ml-1"
                      @click="copyToClipboard(item.issuer)"
                    >
                      <v-icon x-small>mdi-content-copy</v-icon>
                    </v-btn>
                  </div>
                </template>
                <template #item.actions="{ item }">
                  <v-btn
                    small
                    text
                    color="primary"
                    @click="viewToken(item)"
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
import { useQuery } from '@vue/apollo-composable'
import { XRPScreenerGQL } from '~/apollo/queries'
import XrpTokenHeatmap from '~/components/xrp/XrpTokenHeatmap.vue'

export default defineComponent({
  components: {
    XrpTokenHeatmap,
  },
  setup() {
    const { $f } = useContext()
    const loading = ref(true)
    const tokensData = ref<any[]>([])
    
    // Heatmap height
    const heatmapHeight = ref(600)
    
    // Apollo query for XRP screener data
    const { onResult, refetch } = useQuery(XRPScreenerGQL, { 
      fetchPolicy: 'no-cache', 
      pollInterval: 60000 
    })

    // Table headers
    const tableHeaders = [
      { text: 'Token', value: 'tokenName', sortable: false },
      { text: 'Price', value: 'price', align: 'right' },
      { text: '24h Change', value: 'price24h', align: 'right' },
      { text: 'Market Cap', value: 'marketcap', align: 'right' },
      { text: '24h Volume', value: 'volume24h', align: 'right' },
      { text: 'Issuer', value: 'issuer', sortable: false },
      { text: 'Actions', value: 'actions', sortable: false, align: 'center' },
    ]

    // Computed properties
    const topTokens = computed(() => {
      return tokensData.value
        .map(item => ({
          ...item,
          price24h: (Math.random() - 0.5) * 10, // Mock price change
        }))
        .sort((a, b) => (b.marketcap || 0) - (a.marketcap || 0))
        .slice(0, 20)
    })

    const totalTokens = computed(() => tokensData.value.length)
    const gainersCount = computed(() => 
      topTokens.value.filter(token => token.price24h > 0).length
    )
    const losersCount = computed(() => 
      topTokens.value.filter(token => token.price24h < 0).length
    )
    const totalMarketCap = computed(() => 
      $f(tokensData.value.reduce((sum, token) => sum + (token.marketcap || 0), 0), { 
        minDigits: 0, 
        after: '' 
      })
    )

    // Methods
    const formatPrice = (price: number) => {
      return $f(price, { minDigits: 6, after: '' })
    }

    const formatPercentage = (percentage: number) => {
      return $f(percentage, { minDigits: 2, after: '' })
    }

    const formatNumber = (num: number) => {
      return $f(num, { minDigits: 0, after: '' })
    }

    const formatAddress = (address: string) => {
      return `${address.slice(0, 10)}...${address.slice(-10)}`
    }

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // You could add a toast notification here
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }

    const viewToken = (token: any) => {
      const route = useRoute()
      const router = useRouter()
      router.push({
        path: `/token/${token.currency}/xrp`,
        query: { issuer: token.issuerAddress }
      })
    }

    const refreshData = async () => {
      loading.value = true
      await refetch()
    }

    // Handle query results
    onResult((queryResult: any) => {
      if (queryResult.data?.xrpScreener) {
        tokensData.value = queryResult.data.xrpScreener
      } else {
        // Initialize with mock data for development
        tokensData.value = [
          { currency: 'USDC', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'USD Coin', issuerName: 'Circle', marketcap: 1000000, price: 1.0, volume24H: 500000, icon: '🪙' },
          { currency: 'USDT', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Tether', issuerName: 'Tether', marketcap: 800000, price: 1.0, volume24H: 400000, icon: '💎' },
          { currency: 'BTC', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Bitcoin', issuerName: 'BitGo', marketcap: 500000, price: 45000, volume24H: 200000, icon: '₿' },
          { currency: 'ETH', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Ethereum', issuerName: 'BitGo', marketcap: 300000, price: 3000, volume24H: 150000, icon: 'Ξ' },
          { currency: 'SOL', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Solana', issuerName: 'Solana', marketcap: 200000, price: 100, volume24H: 100000, icon: '◎' },
          { currency: 'ADA', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Cardano', issuerName: 'IOG', marketcap: 150000, price: 0.5, volume24H: 75000, icon: '₳' },
          { currency: 'DOT', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Polkadot', issuerName: 'Web3 Foundation', marketcap: 120000, price: 8, volume24H: 60000, icon: '●' },
          { currency: 'LINK', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Chainlink', issuerName: 'Chainlink', marketcap: 100000, price: 15, volume24H: 50000, icon: '🔗' },
          { currency: 'UNI', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Uniswap', issuerName: 'Uniswap', marketcap: 80000, price: 8, volume24H: 40000, icon: '🦄' },
          { currency: 'AAVE', issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh', tokenName: 'Aave', issuerName: 'Aave', marketcap: 60000, price: 300, volume24H: 30000, icon: '👻' },
        ]
      }
      loading.value = false
    })

    onMounted(() => {
      // Set responsive height
      if (process.client) {
        heatmapHeight.value = window.innerHeight * 0.6
      }
    })

    return {
      loading,
      tokensData,
      topTokens,
      totalTokens,
      gainersCount,
      losersCount,
      totalMarketCap,
      heatmapHeight,
      tableHeaders,
      formatPrice,
      formatPercentage,
      formatNumber,
      formatAddress,
      copyToClipboard,
      viewToken,
      refreshData,
    }
  },
  head() {
    return {
      title: 'XRP Token Heatmap - EVM Finance',
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: 'Visualize XRP token performance across the XRP Ledger with our interactive heatmap. Track price changes, market cap, and volume for all XRP tokens.',
        },
      ],
    }
  },
})
</script> 