<template>
  <v-row justify="center">
    <v-col lg="10" md="12">
      <v-row v-if="!loading">
        <v-col cols="12" lg="4">
          <v-row>
            <v-col sm="12" md="6" lg="12">
              <!-- XRP Token Profile -->
              <v-card tile outlined height="364">
                <v-card-title class="d-flex align-center">
                  <v-avatar size="40" class="mr-3">
                    <v-img
                      :src="$imageUrlBySymbol(tokenData.currency?.toLowerCase())"
                      :lazy-src="$imageUrlBySymbol(tokenData.currency?.toLowerCase())"
                    />
                  </v-avatar>
                  <div>
                    <div class="text-h6">{{ tokenData.currency }}</div>
                    <div class="text-subtitle-2 grey--text">{{ tokenData.tokenName }}</div>
                  </div>
                </v-card-title>
                
                <v-card-text>
                  <!-- Issuer Information -->
                  <div class="mb-4">
                    <div class="text-subtitle-2 mb-2">Issuer Information</div>
                    <div class="d-flex align-center mb-2">
                      <span class="text-body-2 mr-2">Name:</span>
                      <span class="text-body-1">{{ tokenData.issuerName || 'Unknown' }}</span>
                    </div>
                    <div class="d-flex align-center">
                      <span class="text-body-2 mr-2">Address:</span>
                      <span class="text-body-1 font-family-mono">{{ tokenData.issuerAddress }}</span>
                      <v-btn
                        icon
                        x-small
                        class="ml-1"
                        @click="copyToClipboard(tokenData.issuerAddress)"
                      >
                        <v-icon size="16">mdi-content-copy</v-icon>
                      </v-btn>
                    </div>
                  </div>

                  <!-- Links -->
                  <div v-if="tokenData.websiteUrl" class="mb-2">
                    <v-btn
                      text
                      small
                      color="primary"
                      :href="tokenData.websiteUrl"
                      target="_blank"
                    >
                      <v-icon left small>mdi-web</v-icon>
                      Website
                    </v-btn>
                  </div>
                  
                  <div v-if="tokenData.sourceCodeUrl" class="mb-2">
                    <v-btn
                      text
                      small
                      color="primary"
                      :href="tokenData.sourceCodeUrl"
                      target="_blank"
                    >
                      <v-icon left small>mdi-code-braces</v-icon>
                      Source Code
                    </v-btn>
                  </div>

                  <!-- Community Links -->
                  <div v-if="tokenData.telegramUrl || tokenData.twitterUrl || tokenData.discordUrl" class="mb-2">
                    <div class="text-subtitle-2 mb-1">Community</div>
                    <div class="d-flex flex-wrap">
                      <v-btn
                        v-if="tokenData.telegramUrl"
                        icon
                        x-small
                        class="mr-1"
                        :href="tokenData.telegramUrl"
                        target="_blank"
                      >
                        <v-icon size="16">mdi-telegram</v-icon>
                      </v-btn>
                      <v-btn
                        v-if="tokenData.twitterUrl"
                        icon
                        x-small
                        class="mr-1"
                        :href="tokenData.twitterUrl"
                        target="_blank"
                      >
                        <v-icon size="16">mdi-twitter</v-icon>
                      </v-btn>
                      <v-btn
                        v-if="tokenData.discordUrl"
                        icon
                        x-small
                        class="mr-1"
                        :href="tokenData.discordUrl"
                        target="_blank"
                      >
                        <v-icon size="16">mdi-discord</v-icon>
                      </v-btn>
                    </div>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>

            <v-col sm="12" md="6" lg="12">
              <!-- XRP Token Balances -->
              <v-card tile outlined height="230">
                <v-card-title>Token Balances</v-card-title>
                <v-card-text>
                  <div v-if="tokenBalances.length > 0">
                    <div
                      v-for="balance in tokenBalances"
                      :key="balance.account"
                      class="d-flex justify-space-between align-center mb-2"
                    >
                      <div>
                        <div class="text-body-2">{{ formatAddress(balance.account) }}</div>
                        <div class="text-caption grey--text">{{ balance.balance }} {{ tokenData.currency }}</div>
                      </div>
                      <div class="text-right">
                        <div class="text-body-2">{{ formatPrice(balance.value) }} XRP</div>
                        <div class="text-caption grey--text">{{ formatPrice(balance.value * xrpPrice) }} USD</div>
                      </div>
                    </div>
                  </div>
                  <div v-else class="text-center grey--text">
                    No balances found
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12" lg="8">
          <v-row>
            <v-col cols="12">
              <!-- XRP Token Metrics -->
              <v-card tile outlined height="102" class="px-3 pt-2">
                <v-row>
                  <v-col cols="3">
                    <div class="text-caption grey--text">Price (XRP)</div>
                    <div class="text-h6">{{ formatPrice(tokenData.price) }}</div>
                  </v-col>
                  <v-col cols="3">
                    <div class="text-caption grey--text">Price (USD)</div>
                    <div class="text-h6">{{ formatPrice(tokenData.price * xrpPrice) }}</div>
                  </v-col>
                  <v-col cols="3">
                    <div class="text-caption grey--text">Market Cap</div>
                    <div class="text-h6">{{ formatPrice(tokenData.marketcap) }} XRP</div>
                  </v-col>
                  <v-col cols="3">
                    <div class="text-caption grey--text">Volume 24H</div>
                    <div class="text-h6">{{ formatPrice(tokenData.volume24H) }} XRP</div>
                  </v-col>
                </v-row>
              </v-card>
            </v-col>

            <v-col cols="12" md="6">
              <!-- XRP AMM Chart -->
              <v-card tile outlined height="500">
                <v-card-title class="d-flex align-center justify-space-between">
                  <span>AMM Chart</span>
                  <div class="d-flex">
                    <v-btn-toggle
                      v-model="selectedTimeframe"
                      mandatory
                      dense
                    >
                      <v-btn value="1h" small>1H</v-btn>
                      <v-btn value="1d" small>1D</v-btn>
                      <v-btn value="1w" small>1W</v-btn>
                    </v-btn-toggle>
                  </div>
                </v-card-title>
                <v-card-text>
                  <div v-if="ammChartData.length > 0" class="text-center">
                    <!-- Placeholder for AMM chart - would integrate with charting library -->
                    <div class="text-h4 grey--text">AMM Chart</div>
                    <div class="text-body-2 grey--text">Data points: {{ ammChartData.length }}</div>
                    <div class="text-caption grey--text">Timeframe: {{ selectedTimeframe }}</div>
                  </div>
                  <div v-else class="text-center grey--text">
                    <v-icon size="48" class="mb-2">mdi-chart-line</v-icon>
                    <div>No AMM data available</div>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>

            <v-col cols="12" md="6">
              <!-- XRP AMM Swap -->
              <xrp-amm-swap
                height="500"
                :in-token="{
                  symbol: tokenData.currency,
                  name: tokenData.tokenName,
                  issuer: tokenData.issuerAddress,
                  decimals: 6,
                }"
                width="100%"
              />
            </v-col>
          </v-row>
        </v-col>

        <!-- Wallet Transactions -->
        <v-col v-if="walletTransactions.length > 0" cols="12">
          <v-card tile outlined>
            <v-card-title>Wallet Transactions</v-card-title>
            <v-card-text>
              <v-data-table
                :headers="transactionHeaders"
                :items="walletTransactions"
                :items-per-page="10"
                class="elevation-0"
              >
                <template #[`item.hash`]="{ item }">
                  <nuxt-link
                    :to="`/xrp-explorer/tx/${item.hash}`"
                    class="text-decoration-none"
                  >
                    {{ formatHash(item.hash) }}
                  </nuxt-link>
                </template>
                
                <template #[`item.amount`]="{ item }">
                  <span :class="item.transactionType === 'Payment' ? 'green--text' : 'blue--text'">
                    {{ formatAmount(item.amount) }} {{ item.currency || 'XRP' }}
                  </span>
                </template>
                
                <template #[`item.transactionType`]="{ item }">
                  <v-chip
                    :color="getTransactionTypeColor(item.transactionType)"
                    small
                    text-color="white"
                  >
                    {{ item.transactionType }}
                  </v-chip>
                </template>
                
                <template #[`item.date`]="{ item }">
                  {{ formatDate(item.date) }}
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Price Chart -->
        <v-col cols="12">
          <v-card tile outlined height="600">
            <v-card-title class="d-flex align-center justify-space-between">
              <span>Price Chart</span>
              <v-btn-toggle
                v-model="priceDisplayMode"
                mandatory
                dense
              >
                <v-btn value="xrp" small>XRP</v-btn>
                <v-btn value="usd" small>USD</v-btn>
              </v-btn-toggle>
            </v-card-title>
            <v-card-text>
              <div class="text-center">
                <!-- Placeholder for price chart - would integrate with charting library -->
                <div class="text-h4 grey--text">Price Chart</div>
                <div class="text-body-2 grey--text">Display mode: {{ priceDisplayMode.toUpperCase() }}</div>
                <div class="text-caption grey--text">Token: {{ tokenData.currency }}</div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Loading State -->
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
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useRoute, useContext } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPScreenerGQL, XRPAccountTransactionsGQL } from '~/apollo/queries'
import XrpAmmSwap from '~/components/trading/XrpAmmSwap.vue'

interface XRPTokenData {
  currency: string
  issuerAddress: string
  tokenName: string
  issuerName: string
  price: number
  marketcap: number
  volume24H: number
  websiteUrl?: string
  sourceCodeUrl?: string
  telegramUrl?: string
  twitterUrl?: string
  discordUrl?: string
}

interface XRPBalance {
  account: string
  balance: number
  value: number
}

interface XRPTransaction {
  hash: string
  account: string
  destination: string
  transactionType: string
  amount: number
  currency?: string
  fee: number
  date: string
}

export default defineComponent({
  components: {
    XrpAmmSwap,
  },
  setup() {
    const { $f } = useContext()
    const route = useRoute()
    const loading = ref(true)
    
    // Get token currency and issuer from route
    const tokenCurrency = computed(() => route.value.params?.id ?? '')
    const issuerAddress = computed(() => route.value.query?.issuer as string ?? '')
    
    // State
    const selectedTimeframe = ref('1d')
    const priceDisplayMode = ref('xrp')
    const xrpPrice = ref(0.5) // Mock XRP price, would come from API
    
    // Token data
    const tokenData = ref<XRPTokenData>({
      currency: tokenCurrency.value,
      issuerAddress: issuerAddress.value,
      tokenName: '',
      issuerName: '',
      price: 0,
      marketcap: 0,
      volume24H: 0,
    })
    
    // Balances and transactions
    const tokenBalances = ref<XRPBalance[]>([])
    const walletTransactions = ref<XRPTransaction[]>([])
    const ammChartData = ref<any[]>([])
    
    // Queries
    const { onResult: onScreenerResult } = useQuery(XRPScreenerGQL, { 
      fetchPolicy: 'no-cache', 
      pollInterval: 60000 
    })
    
    const { onResult: onTransactionsResult } = useQuery(
      XRPAccountTransactionsGQL,
      () => ({ address: issuerAddress.value }),
      { fetchPolicy: 'no-cache', pollInterval: 30000 }
    )
    
    // Computed
    const transactionHeaders = computed(() => [
      { text: 'Hash', value: 'hash', width: '200' },
      { text: 'Type', value: 'transactionType', width: '120' },
      { text: 'Amount', value: 'amount', width: '120' },
      { text: 'Fee', value: 'fee', width: '80' },
      { text: 'Date', value: 'date', width: '120' },
    ])
    
    // Methods
    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // Could add toast notification here
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }
    
    const formatAddress = (address: string) => {
      return `${address.slice(0, 10)}...${address.slice(-10)}`
    }
    
    const formatHash = (hash: string) => {
      return `${hash.slice(0, 10)}...${hash.slice(-10)}`
    }
    
    const formatPrice = (price: number) => {
      return $f(price, { minDigits: 6, after: '' })
    }
    
    const formatAmount = (amount: number) => {
      return $f(amount, { minDigits: 2, after: '' })
    }
    
    const formatDate = (date: string) => {
      return new Date(date).toLocaleDateString()
    }
    
    const getTransactionTypeColor = (type: string) => {
      switch (type) {
        case 'Payment': return 'green'
        case 'OfferCreate': return 'blue'
        case 'OfferCancel': return 'orange'
        case 'TrustSet': return 'purple'
        default: return 'grey'
      }
    }
    
    // Mock data generation for development
    const generateMockData = () => {
      // Mock token data
      tokenData.value = {
        ...tokenData.value,
        tokenName: `${tokenCurrency.value} Token`,
        issuerName: 'Sample Issuer',
        price: 0.001234,
        marketcap: 1234567,
        volume24H: 98765,
        websiteUrl: 'https://example.com',
        sourceCodeUrl: 'https://github.com/example',
        telegramUrl: 'https://t.me/example',
        twitterUrl: 'https://twitter.com/example',
        discordUrl: 'https://discord.gg/example',
      }
      
      // Mock balances
      tokenBalances.value = [
        {
          account: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
          balance: 1000000,
          value: 1234.56,
        },
        {
          account: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
          balance: 500000,
          value: 617.28,
        },
      ]
      
      // Mock transactions
      walletTransactions.value = [
        {
          hash: 'E0C1D4B24D76B4180D2C96450438A0BE14304E69EDFBE91DF6211C923B344401',
          account: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
          destination: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
          transactionType: 'Payment',
          amount: 1000000,
          currency: tokenCurrency.value,
          fee: 12,
          date: new Date().toISOString(),
        },
        {
          hash: 'F1D2E5C35E87C5291E3D07561549B1CF25415F7AFEGC02EG7322D034C455502',
          account: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
          destination: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
          transactionType: 'OfferCreate',
          amount: 500000,
          currency: tokenCurrency.value,
          fee: 12,
          date: new Date(Date.now() - 86400000).toISOString(),
        },
      ]
      
      // Mock AMM chart data
      ammChartData.value = Array.from({ length: 24 }, (_, i) => ({
        time: new Date(Date.now() - (23 - i) * 3600000).toISOString(),
        price: 0.001 + Math.random() * 0.0005,
        volume: Math.random() * 1000,
      }))
    }
    
    // Event handlers
    onScreenerResult((queryResult: any) => {
      const screenerData = queryResult.data?.xrpScreener ?? []
      const tokenInfo = screenerData.find((t: any) => 
        t.currency === tokenCurrency.value && t.issuerAddress === issuerAddress.value
      )
      
      if (tokenInfo) {
        tokenData.value = {
          ...tokenData.value,
          tokenName: tokenInfo.tokenName,
          issuerName: tokenInfo.issuerName,
          price: tokenInfo.price,
          marketcap: tokenInfo.marketcap,
          volume24H: tokenInfo.volume24H,
        }
      }
      
      loading.value = false
    })
    
    onTransactionsResult((queryResult: any) => {
      const transactions = queryResult.data?.xrpTransactions ?? []
      walletTransactions.value = transactions.map((tx: any) => ({
        ...tx,
        currency: tokenCurrency.value,
      }))
    })
    
    // Initialize with mock data for development
    generateMockData()
    loading.value = false
    
    return {
      loading,
      tokenData,
      tokenBalances,
      walletTransactions,
      ammChartData,
      selectedTimeframe,
      priceDisplayMode,
      xrpPrice,
      transactionHeaders,
      copyToClipboard,
      formatAddress,
      formatHash,
      formatPrice,
      formatAmount,
      formatDate,
      getTransactionTypeColor,
    }
  },
  head: {
    title: 'XRP Token Details',
  },
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style> 