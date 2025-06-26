<template>
  <div>
    <v-container fluid>
      <!-- Header -->
      <v-row>
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2">mdi-history</v-icon>
              <span class="text-h6">XRP Transaction History</span>
              <v-spacer />
              <v-btn
                icon
                @click="refreshData"
                :loading="loading"
              >
                <v-icon>mdi-refresh</v-icon>
              </v-btn>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>

      <!-- Address Input -->
      <v-row v-if="!address">
        <v-col cols="12" md="8" lg="6">
          <v-card tile outlined>
            <v-card-text>
              <v-text-field
                v-model="addressInput"
                label="Enter XRP Address"
                placeholder="rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
                outlined
                dense
                :rules="[addressRule]"
                @keyup.enter="loadAddress"
              >
                <template #append>
                  <v-btn
                    @click="loadAddress"
                    :disabled="!isValidAddress"
                    :loading="loading"
                    color="primary"
                  >
                    Load
                  </v-btn>
                </template>
              </v-text-field>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Account Info -->
      <v-row v-if="address">
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title>
              <v-icon class="mr-2">mdi-account</v-icon>
              Account: {{ formatAddress(address) }}
              <v-btn
                icon
                x-small
                @click="copyToClipboard(address)"
                class="ml-2"
              >
                <v-icon>mdi-content-copy</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-row>
                <v-col cols="12" md="4">
                  <v-list dense>
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Total Transactions</v-list-item-subtitle>
                        <v-list-item-title class="text-h6">
                          {{ transactions.length }}
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </v-col>
                
                <v-col cols="12" md="4">
                  <v-list dense>
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Filter</v-list-item-subtitle>
                        <v-list-item-title>
                          <v-select
                            v-model="transactionType"
                            :items="transactionTypes"
                            dense
                            outlined
                            hide-details
                            style="max-width: 200px"
                          />
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </v-col>
                
                <v-col cols="12" md="4">
                  <v-list dense>
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Search</v-list-item-subtitle>
                        <v-list-item-title>
                          <v-text-field
                            v-model="searchQuery"
                            placeholder="Search transactions..."
                            dense
                            outlined
                            hide-details
                            prepend-inner-icon="mdi-magnify"
                            style="max-width: 200px"
                          />
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Transactions Table -->
      <v-row v-if="address">
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title>
              <v-icon class="mr-2">mdi-format-list-bulleted</v-icon>
              Transactions
            </v-card-title>
            <v-card-text>
              <v-data-table
                :headers="transactionHeaders"
                :items="filteredTransactions"
                :loading="loading"
                :search="searchQuery"
                dense
                class="elevation-0"
                :items-per-page="20"
                :footer-props="{
                  'items-per-page-options': [10, 20, 50, 100]
                }"
              >
                <template #item.transactionType="{ item }">
                  <v-chip
                    :color="getTransactionTypeColor(item.transactionType)"
                    x-small
                    text-color="white"
                  >
                    {{ item.transactionType }}
                  </v-chip>
                </template>

                <template #item.amount="{ item }">
                  <span :class="getAmountClass(item)">
                    {{ formatAmount(item.amount) }}
                  </span>
                </template>

                <template #item.fee="{ item }">
                  {{ formatFee(item.fee) }}
                </template>

                <template #item.hash="{ item }">
                  <div class="d-flex align-center">
                    <span class="font-family-mono caption">{{ formatHash(item.hash) }}</span>
                    <v-btn
                      icon
                      x-small
                      @click="copyToClipboard(item.hash)"
                      class="ml-1"
                    >
                      <v-icon>mdi-content-copy</v-icon>
                    </v-btn>
                  </div>
                </template>

                <template #item.ledgerIndex="{ item }">
                  <v-btn
                    small
                    text
                    color="primary"
                    @click="viewLedger(item.ledgerIndex)"
                  >
                    {{ item.ledgerIndex }}
                  </v-btn>
                </template>

                <template #item.actions="{ item }">
                  <v-btn
                    small
                    text
                    color="primary"
                    @click="viewTransaction(item.hash)"
                  >
                    View
                  </v-btn>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Error Message -->
      <v-row v-if="error">
        <v-col cols="12">
          <v-alert
            type="error"
            dismissible
            @click:close="error = null"
          >
            {{ error }}
          </v-alert>
        </v-col>
      </v-row>

      <!-- Loading -->
      <v-row v-if="loading && !address">
        <v-col cols="12">
          <v-skeleton-loader
            type="card"
            class="mx-auto"
            max-width="400"
          />
        </v-col>
      </v-row>

      <!-- No Data -->
      <v-row v-if="address && !loading && transactions.length === 0">
        <v-col cols="12">
          <v-alert
            type="info"
            text
          >
            No transactions found for this address.
          </v-alert>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, useContext } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable'
import { XRPAccountTransactionsGQL } from '~/apollo/queries'

export default defineComponent({
  setup() {
    const { $f } = useContext()
    const { query } = useRoute()
    
    // State
    const address = ref<string>('')
    const addressInput = ref<string>('')
    const loading = ref(false)
    const error = ref<string | null>(null)
    const searchQuery = ref('')
    const transactionType = ref('all')

    // Get address from query params
    onMounted(() => {
      if (query.address) {
        address.value = query.address as string
        loadAddressData()
      }
    })

    // XRP address validation
    const isValidAddress = computed(() => {
      return /^r[a-km-zA-HJ-NP-Z1-9]{25,34}$/.test(addressInput.value)
    })

    const addressRule = (value: string) => {
      if (!value) return 'Address is required'
      if (!isValidAddress.value) return 'Invalid XRP address format'
      return true
    }

    // Transaction types
    const transactionTypes = [
      { text: 'All Types', value: 'all' },
      { text: 'Payment', value: 'Payment' },
      { text: 'OfferCreate', value: 'OfferCreate' },
      { text: 'OfferCancel', value: 'OfferCancel' },
      { text: 'TrustSet', value: 'TrustSet' },
      { text: 'AccountSet', value: 'AccountSet' },
      { text: 'AMMDeposit', value: 'AMMDeposit' },
      { text: 'AMMWithdraw', value: 'AMMWithdraw' },
      { text: 'AMMTrade', value: 'AMMTrade' },
    ]

    // GraphQL query
    const { result, loading: queryLoading, error: queryError, refetch } = useQuery(
      XRPAccountTransactionsGQL,
      computed(() => ({ address: address.value })),
      () => ({
        enabled: !!address.value,
        fetchPolicy: 'cache-and-network',
        errorPolicy: 'all'
      })
    )

    // Computed data
    const transactions = computed(() => result.value?.xrpTransactions || [])
    
    const filteredTransactions = computed(() => {
      let filtered = transactions.value
      
      if (transactionType.value !== 'all') {
        filtered = filtered.filter(tx => tx.transactionType === transactionType.value)
      }
      
      return filtered
    })

    // Table headers
    const transactionHeaders = [
      { text: 'Type', value: 'transactionType', sortable: true },
      { text: 'Amount', value: 'amount', sortable: true, align: 'right' },
      { text: 'Destination', value: 'destination', sortable: false },
      { text: 'Fee', value: 'fee', sortable: true, align: 'right' },
      { text: 'Hash', value: 'hash', sortable: false },
      { text: 'Ledger', value: 'ledgerIndex', sortable: true, align: 'center' },
      { text: 'Actions', value: 'actions', sortable: false, align: 'center' },
    ]

    // Methods
    const loadAddress = () => {
      if (isValidAddress.value) {
        address.value = addressInput.value
        loadAddressData()
      }
    }

    const loadAddressData = async () => {
      if (!address.value) return
      
      loading.value = true
      error.value = null
      
      try {
        await refetch()
        if (queryError.value) {
          error.value = 'Failed to load transaction data'
        }
      } catch (err: any) {
        error.value = err.message || 'Failed to load transaction data'
      } finally {
        loading.value = false
      }
    }

    const refreshData = () => {
      loadAddressData()
    }

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // You can add a toast notification here
      } catch (err) {
        console.error('Failed to copy to clipboard:', err)
      }
    }

    const viewTransaction = (hash: string) => {
      window.$nuxt.$router.push(`/xrp-explorer/tx/${hash}`)
    }

    const viewLedger = (ledgerIndex: number) => {
      window.$nuxt.$router.push(`/xrp-explorer/ledger/${ledgerIndex}`)
    }

    // Formatting functions
    const formatAddress = (address: string) => {
      if (!address) return ''
      return address.length > 20 ? `${address.substring(0, 10)}...${address.substring(address.length - 10)}` : address
    }

    const formatAmount = (amount: any) => {
      if (!amount) return '0'
      
      // Handle different amount formats (XRP drops, IOU amounts)
      if (typeof amount === 'string') {
        // If it's a string, it might be drops (1 XRP = 1,000,000 drops)
        const drops = parseFloat(amount)
        if (drops >= 1000000) {
          return `${$f(drops / 1000000, { minDigits: 6, after: '' })} XRP`
        } else {
          return `${$f(drops, { minDigits: 0, after: '' })} drops`
        }
      } else if (typeof amount === 'object' && amount.currency) {
        // IOU amount
        return `${$f(amount.value, { minDigits: 2, after: '' })} ${amount.currency}`
      }
      
      return $f(amount || 0, { minDigits: 2, after: '' })
    }

    const formatFee = (fee: any) => {
      if (!fee) return '0'
      const drops = parseFloat(fee)
      return `${$f(drops / 1000000, { minDigits: 6, after: '' })} XRP`
    }

    const formatHash = (hash: string) => {
      if (!hash) return ''
      return hash.length > 20 ? `${hash.substring(0, 10)}...${hash.substring(hash.length - 10)}` : hash
    }

    const getTransactionTypeColor = (type: string) => {
      const colors: { [key: string]: string } = {
        Payment: 'green',
        OfferCreate: 'blue',
        OfferCancel: 'orange',
        TrustSet: 'purple',
        AccountSet: 'indigo',
        AMMDeposit: 'teal',
        AMMWithdraw: 'cyan',
        AMMTrade: 'deep-purple',
      }
      return colors[type] || 'grey'
    }

    const getAmountClass = (item: any) => {
      // Determine if this is an incoming or outgoing transaction
      // This is a simplified logic - you might need to adjust based on your data structure
      if (item.transactionType === 'Payment' && item.destination) {
        return 'green--text'
      }
      return 'red--text'
    }

    return {
      // State
      address,
      addressInput,
      loading: computed(() => loading.value || queryLoading.value),
      error,
      searchQuery,
      transactionType,
      
      // Computed
      isValidAddress,
      transactions,
      filteredTransactions,
      transactionHeaders,
      transactionTypes,
      
      // Methods
      addressRule,
      loadAddress,
      refreshData,
      copyToClipboard,
      viewTransaction,
      viewLedger,
      
      // Formatting
      formatAddress,
      formatAmount,
      formatFee,
      formatHash,
      getTransactionTypeColor,
      getAmountClass,
    }
  },
})
</script> 