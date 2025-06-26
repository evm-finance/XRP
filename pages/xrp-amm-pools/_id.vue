<template>
  <div>
    <v-row no-gutters justify="center">
      <v-col cols="12" lg="10">
        <!-- Pool Header -->
        <v-row justify="center" class="mb-4">
          <v-col cols="12">
            <v-card tile outlined>
              <v-card-title class="d-flex align-center">
                <div class="d-flex align-center">
                  <v-avatar size="40" class="mr-3">
                    <v-img
                      :src="$imageUrlBySymbol(poolData.token1.symbol.toLowerCase())"
                      :lazy-src="$imageUrlBySymbol(poolData.token1.symbol.toLowerCase())"
                    />
                  </v-avatar>
                  <span class="text-h5 mr-2">{{ poolData.token1.symbol }}</span>
                  <v-icon class="mx-2">mdi-swap-horizontal</v-icon>
                  <v-avatar size="40" class="mr-3">
                    <v-img
                      :src="$imageUrlBySymbol(poolData.token2.symbol.toLowerCase())"
                      :lazy-src="$imageUrlBySymbol(poolData.token2.symbol.toLowerCase())"
                    />
                  </v-avatar>
                  <span class="text-h5">{{ poolData.token2.symbol }}</span>
                </div>
                <v-spacer />
                <v-chip color="primary" large>
                  Pool ID: {{ $route.params.id }}
                </v-chip>
              </v-card-title>
              
              <v-card-text>
                <v-row>
                  <v-col cols="12" md="3">
                    <div class="text-center">
                      <div class="text-h6 font-weight-bold">${{ formatNumber(poolData.liquidity) }}</div>
                      <div class="text-caption grey--text">Total Liquidity</div>
                    </div>
                  </v-col>
                  <v-col cols="12" md="3">
                    <div class="text-center">
                      <div class="text-h6 font-weight-bold text-success">${{ formatNumber(poolData.volume24h) }}</div>
                      <div class="text-caption grey--text">24h Volume</div>
                    </div>
                  </v-col>
                  <v-col cols="12" md="3">
                    <div class="text-center">
                      <div class="text-h6 font-weight-bold text-info">{{ formatPercentage(poolData.apr) }}%</div>
                      <div class="text-caption grey--text">APR</div>
                    </div>
                  </v-col>
                  <v-col cols="12" md="3">
                    <div class="text-center">
                      <div class="text-h6 font-weight-bold text-warning">{{ formatFee(poolData.fee) }}%</div>
                      <div class="text-caption grey--text">Trading Fee</div>
                    </div>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Main Content -->
        <v-row>
          <!-- Left Column - Pool Actions -->
          <v-col cols="12" lg="4">
            <v-row>
              <!-- Pool Balances -->
              <v-col cols="12">
                <v-card tile outlined class="mb-4">
                  <v-card-title>Pool Balances</v-card-title>
                  <v-card-text>
                    <div class="mb-3">
                      <div class="d-flex justify-space-between align-center mb-2">
                        <div class="d-flex align-center">
                          <v-avatar size="24" class="mr-2">
                            <v-img :src="$imageUrlBySymbol(poolData.token1.symbol.toLowerCase())" />
                          </v-avatar>
                          <span>{{ poolData.token1.symbol }}</span>
                        </div>
                        <span class="font-weight-medium">{{ formatBalance(poolData.token1Balance) }}</span>
                      </div>
                      <div class="d-flex justify-space-between align-center">
                        <div class="d-flex align-center">
                          <v-avatar size="24" class="mr-2">
                            <v-img :src="$imageUrlBySymbol(poolData.token2.symbol.toLowerCase())" />
                          </v-avatar>
                          <span>{{ poolData.token2.symbol }}</span>
                        </div>
                        <span class="font-weight-medium">{{ formatBalance(poolData.token2Balance) }}</span>
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- User Pool Token Balance -->
              <v-col cols="12">
                <v-card tile outlined class="mb-4">
                  <v-card-title>Your Pool Tokens</v-card-title>
                  <v-card-text>
                    <div v-if="userPoolTokens > 0">
                      <div class="d-flex justify-space-between align-center mb-2">
                        <span>Pool Tokens</span>
                        <span class="font-weight-medium">{{ formatBalance(userPoolTokens) }}</span>
                      </div>
                      <div class="d-flex justify-space-between align-center">
                        <span>Value</span>
                        <span class="font-weight-medium text-success">${{ formatNumber(userPoolTokensValue) }}</span>
                      </div>
                      <div class="d-flex justify-space-between align-center">
                        <span>Share</span>
                        <span class="font-weight-medium">{{ formatPercentage(userPoolShare) }}%</span>
                      </div>
                    </div>
                    <div v-else class="text-center grey--text">
                      <v-icon size="48" class="mb-2">mdi-wallet-outline</v-icon>
                      <div>No pool tokens</div>
                      <div class="text-caption">Deposit to earn fees</div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Quick Actions -->
              <v-col cols="12">
                <v-card tile outlined class="mb-4">
                  <v-card-title>Quick Actions</v-card-title>
                  <v-card-text>
                    <v-btn
                      block
                      color="primary"
                      class="mb-2"
                      @click="actionDialog.openDialog()"
                    >
                      <v-icon left>mdi-plus</v-icon>
                      Deposit/Withdraw
                    </v-btn>
                    
                    <v-btn
                      block
                      outlined
                      @click="swapDialog.openDialog()"
                    >
                      <v-icon left>mdi-swap-horizontal</v-icon>
                      Swap
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-col>

          <!-- Right Column - Charts and History -->
          <v-col cols="12" lg="8">
            <v-row>
              <!-- Price Chart -->
              <v-col cols="12">
                <v-card tile outlined class="mb-4">
                  <v-card-title class="d-flex align-center justify-space-between">
                    <span>Price Chart</span>
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
                    <div v-if="chartData.length > 0" class="text-center">
                      <!-- Placeholder for price chart - would integrate with charting library -->
                      <div class="text-h4 grey--text">Price Chart</div>
                      <div class="text-body-2 grey--text">Data points: {{ chartData.length }}</div>
                      <div class="text-caption grey--text">Timeframe: {{ selectedTimeframe }}</div>
                    </div>
                    <div v-else class="text-center grey--text">
                      <v-icon size="48" class="mb-2">mdi-chart-line</v-icon>
                      <div>No chart data available</div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Transaction History Tabs -->
              <v-col cols="12">
                <v-card tile outlined>
                  <v-card-title>
                    <v-tabs v-model="activeTab" background-color="transparent">
                      <v-tab>Pool Transactions</v-tab>
                      <v-tab>Your Transactions</v-tab>
                    </v-tabs>
                  </v-card-title>
                  
                  <v-card-text>
                    <v-tabs-items v-model="activeTab">
                      <!-- Pool Transactions -->
                      <v-tab-item>
                        <v-data-table
                          :headers="poolTxHeaders"
                          :items="poolTransactions"
                          :loading="loading"
                          :items-per-page="10"
                          class="elevation-0"
                          hide-default-footer
                        >
                          <template #item.hash="{ item }">
                            <div class="d-flex align-center">
                              <span class="font-family-mono font-weight-medium">{{ formatHash(item.hash) }}</span>
                              <v-btn
                                icon
                                x-small
                                class="ml-1"
                                @click="copyToClipboard(item.hash)"
                              >
                                <v-icon size="16">mdi-content-copy</v-icon>
                              </v-btn>
                              <v-btn
                                icon
                                x-small
                                class="ml-1"
                                @click="openTransaction(item.hash)"
                              >
                                <v-icon size="16">mdi-open-in-new</v-icon>
                              </v-btn>
                            </div>
                          </template>
                          
                          <template #item.type="{ item }">
                            <v-chip
                              :color="getTransactionTypeColor(item.type)"
                              small
                              text-color="white"
                            >
                              {{ item.type }}
                            </v-chip>
                          </template>
                          
                          <template #item.amount="{ item }">
                            <span class="font-weight-medium">{{ formatAmount(item.amount) }}</span>
                          </template>
                          
                          <template #item.fee="{ item }">
                            <span class="grey--text">{{ formatFee(item.fee) }}</span>
                          </template>
                        </v-data-table>
                      </v-tab-item>

                      <!-- User Transactions -->
                      <v-tab-item>
                        <v-data-table
                          :headers="userTxHeaders"
                          :items="userTransactions"
                          :loading="loading"
                          :items-per-page="10"
                          class="elevation-0"
                          hide-default-footer
                        >
                          <template #item.hash="{ item }">
                            <div class="d-flex align-center">
                              <span class="font-family-mono font-weight-medium">{{ formatHash(item.hash) }}</span>
                              <v-btn
                                icon
                                x-small
                                class="ml-1"
                                @click="copyToClipboard(item.hash)"
                              >
                                <v-icon size="16">mdi-content-copy</v-icon>
                              </v-btn>
                              <v-btn
                                icon
                                x-small
                                class="ml-1"
                                @click="openTransaction(item.hash)"
                              >
                                <v-icon size="16">mdi-open-in-new</v-icon>
                              </v-btn>
                            </div>
                          </template>
                          
                          <template #item.token="{ item }">
                            <div class="d-flex align-center">
                              <v-avatar size="20" class="mr-2">
                                <v-img :src="$imageUrlBySymbol(item.token.toLowerCase())" />
                              </v-avatar>
                              <span>{{ item.token }}</span>
                            </div>
                          </template>
                          
                          <template #item.type="{ item }">
                            <v-chip
                              :color="getTransactionTypeColor(item.type)"
                              small
                              text-color="white"
                            >
                              {{ item.type }}
                            </v-chip>
                          </template>
                          
                          <template #item.amount="{ item }">
                            <span class="font-weight-medium">{{ formatAmount(item.amount) }}</span>
                          </template>
                        </v-data-table>
                      </v-tab-item>
                    </v-tabs-items>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <!-- Action Dialogs -->
    <xrp-amm-action-dialog
      ref="actionDialog"
      :pool="poolData"
    />
    
    <xrp-amm-swap-dialog
      ref="swapDialog"
      :pool="poolData"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, useRoute } from '@nuxtjs/composition-api'
import { useXrpAmmLiveData } from '~/composables/useXrpAmmLiveData'
import XrpAmmActionDialog from '~/components/xrp/XrpAmmActionDialog.vue'
import XrpAmmSwapDialog from '~/components/xrp/XrpAmmSwapDialog.vue'

interface PoolData {
  id: string
  token1: { symbol: string; name: string; icon: string }
  token2: { symbol: string; name: string; icon: string; issuer?: string }
  liquidity: number
  volume24h: number
  fee: number
  apr: number
  priceChange24h: number
  token1Balance: number
  token2Balance: number
  totalSupply: number
  createdAt: string
  lastUpdated: string
  transactions: {
    hash: string
    type: string
    amount: number
    token: string
    timestamp: string
    user: string
  }[]
  userPositions: {
    user: string
    poolTokens: number
    token1Balance: number
    token2Balance: number
    share: number
    value: number
  }[]
}

interface Transaction {
  hash: string
  type: string
  amount: number
  fee: number
  timestamp: string
  token?: string
}

export default defineComponent({
  components: {
    XrpAmmActionDialog,
    XrpAmmSwapDialog,
  },
  setup() {
    const { $f } = useContext()
    const route = useRoute()
    const { getPoolDetails, getUserTokenBalances } = useXrpAmmLiveData()
    
    // State
    const loading = ref(true)
    const activeTab = ref(0)
    const selectedTimeframe = ref('1d')
    
    // Dialog states
    const showDepositDialog = ref(false)
    const showWithdrawDialog = ref(false)
    const showSwapDialog = ref(false)
    
    // Form data
    const depositAmount1 = ref('')
    const depositAmount2 = ref('')
    const withdrawAmount = ref('')
    const swapFrom = ref('')
    const swapTo = ref('')
    const swapAmount = ref('')
    
    // Loading states
    const depositing = ref(false)
    const withdrawing = ref(false)
    const swapping = ref(false)
    
    // Get pool details from GraphQL
    const { poolDetails, loading: poolLoading, error: poolError, refetch: refetchPool } = getPoolDetails(route.value.params.id)
    
    // Get user token balances
    const { tokenBalances, loading: balancesLoading } = getUserTokenBalances()
    
    // Computed pool data
    const poolData = computed<PoolData>(() => {
      if (!poolDetails.value) {
        return {
          id: route.value.params.id,
          token1: { symbol: 'XRP', name: 'Ripple', icon: '🪙' },
          token2: { symbol: 'USDC', name: 'USD Coin', icon: '💎', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
          liquidity: 0,
          volume24h: 0,
          fee: 0.003,
          apr: 0,
          priceChange24h: 0,
          token1Balance: 0,
          token2Balance: 0,
          totalSupply: 0,
          createdAt: '',
          lastUpdated: '',
          transactions: [],
          userPositions: []
        }
      }
      
      return poolDetails.value
    })
    
    // User position in this pool
    const userPosition = computed(() => {
      if (!poolDetails.value?.userPositions) return null
      return poolDetails.value.userPositions.find(pos => pos.user === 'current-user') || null
    })
    
    // User pool tokens
    const userPoolTokens = computed(() => userPosition.value?.poolTokens || 0)
    
    // User pool tokens value
    const userPoolTokensValue = computed(() => userPosition.value?.value || 0)
    
    // User pool share percentage
    const userPoolShare = computed(() => {
      if (!poolData.value.totalSupply || !userPoolTokens.value) return 0
      return (userPoolTokens.value / poolData.value.totalSupply) * 100
    })
    
    // Chart data (mock for now)
    const chartData = ref([])
    
    // Transaction data
    const poolTransactions = ref<Transaction[]>([
      {
        hash: '0x1234567890abcdef',
        type: 'Swap',
        amount: 1000,
        fee: 0.003,
        timestamp: '2024-01-01T12:00:00Z',
      },
      {
        hash: '0xabcdef1234567890',
        type: 'Deposit',
        amount: 500,
        fee: 0.001,
        timestamp: '2024-01-01T11:00:00Z',
      },
    ])
    
    const userTransactions = ref<Transaction[]>([
      {
        hash: '0xuser1234567890',
        type: 'Swap',
        amount: 100,
        fee: 0.003,
        timestamp: '2024-01-01T12:00:00Z',
        token: 'XRP',
      },
      {
        hash: '0xuserabcdef123456',
        type: 'Deposit',
        amount: 50,
        fee: 0.001,
        timestamp: '2024-01-01T11:00:00Z',
        token: 'USDC',
      },
    ])
    
    // Table headers
    const poolTxHeaders = [
      { text: 'Hash', value: 'hash', sortable: false },
      { text: 'Type', value: 'type', align: 'center' },
      { text: 'Amount', value: 'amount', align: 'right' },
      { text: 'Fee', value: 'fee', align: 'right' },
      { text: 'Time', value: 'timestamp', align: 'right' },
    ]
    
    const userTxHeaders = [
      { text: 'Hash', value: 'hash', sortable: false },
      { text: 'Token', value: 'token', align: 'left' },
      { text: 'Type', value: 'type', align: 'center' },
      { text: 'Amount', value: 'amount', align: 'right' },
      { text: 'Time', value: 'timestamp', align: 'right' },
    ]
    
    // Methods
    const formatNumber = (value: number): string => {
      return $f(value, { minDigits: 0, after: '' })
    }
    
    const formatPercentage = (value: number): string => {
      return value.toFixed(2)
    }
    
    const formatFee = (value: number): string => {
      return (value * 100).toFixed(3)
    }
    
    const formatBalance = (value: number): string => {
      return $f(value, { minDigits: 6, after: '' })
    }
    
    const formatHash = (hash: string): string => {
      return `${hash.slice(0, 8)}...${hash.slice(-8)}`
    }
    
    const formatAmount = (amount: number): string => {
      return $f(amount, { minDigits: 2, after: '' })
    }
    
    const getTransactionTypeColor = (type: string): string => {
      switch (type.toLowerCase()) {
        case 'swap': return 'primary'
        case 'deposit': return 'success'
        case 'withdraw': return 'warning'
        default: return 'grey'
      }
    }
    
    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }
    
    const openTransaction = (hash: string) => {
      window.open(`/xrp-explorer/tx/${hash}`, '_blank')
    }
    
    const deposit = async () => {
      depositing.value = true
      try {
        // Deposit logic would be handled by the action dialog
        await new Promise(resolve => setTimeout(resolve, 2000))
        await refetchPool()
      } catch (error) {
        console.error('Deposit failed:', error)
      } finally {
        depositing.value = false
      }
    }
    
    const withdraw = async () => {
      withdrawing.value = true
      try {
        // Withdraw logic would be handled by the action dialog
        await new Promise(resolve => setTimeout(resolve, 2000))
        await refetchPool()
      } catch (error) {
        console.error('Withdraw failed:', error)
      } finally {
        withdrawing.value = false
      }
    }
    
    const swap = async () => {
      swapping.value = true
      try {
        // Swap logic would be handled by the swap dialog
        await new Promise(resolve => setTimeout(resolve, 2000))
        await refetchPool()
      } catch (error) {
        console.error('Swap failed:', error)
      } finally {
        swapping.value = false
      }
    }
    
    // Initialize data
    const initializeData = () => {
      // Pool data will be loaded automatically via GraphQL
      // Chart data would be loaded here in real implementation
    }
    
    initializeData()

    const actionDialog = ref()
    const swapDialog = ref()

    return {
      // Data
      loading,
      activeTab,
      selectedTimeframe,
      poolData,
      userPoolTokens,
      userPoolTokensValue,
      userPoolShare,
      chartData,
      poolTransactions,
      userTransactions,
      poolTxHeaders,
      userTxHeaders,
      
      // Dialog states
      showDepositDialog,
      showWithdrawDialog,
      showSwapDialog,
      
      // Form data
      depositAmount1,
      depositAmount2,
      withdrawAmount,
      swapFrom,
      swapTo,
      swapAmount,
      
      // Loading states
      depositing,
      withdrawing,
      swapping,
      
      // Methods
      formatNumber,
      formatPercentage,
      formatFee,
      formatBalance,
      formatHash,
      formatAmount,
      getTransactionTypeColor,
      copyToClipboard,
      openTransaction,
      deposit,
      withdraw,
      swap,
      actionDialog,
      swapDialog,
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