<template>
  <div>
    <v-container fluid>
      <!-- Header -->
      <v-row>
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2">mdi-wallet</v-icon>
              <span class="text-h6">XRP Account Balances</span>
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
                placeholder="rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo"
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
              
              <!-- Default Address Examples -->
              <v-row class="mt-3">
                <v-col cols="12">
                  <v-subheader class="px-0">Example Addresses:</v-subheader>
                  <v-chip-group>
                    <v-chip
                      v-for="exampleAddress in exampleAddresses"
                      :key="exampleAddress"
                      small
                      outlined
                      clickable
                      @click="loadExampleAddress(exampleAddress)"
                      class="font-family-mono"
                    >
                      {{ formatExampleAddress(exampleAddress) }}
                    </v-chip>
                  </v-chip-group>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Account Info -->
      <v-row v-if="address && accountData">
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title>
              <v-icon class="mr-2">mdi-account</v-icon>
              Account Information
            </v-card-title>
            <v-card-text>
              <v-row>
                <v-col cols="12" md="6">
                  <v-list dense>
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Address</v-list-item-subtitle>
                        <v-list-item-title class="font-family-mono">
                          {{ address }}
                          <v-btn
                            icon
                            x-small
                            @click="copyToClipboard(address)"
                            class="ml-2"
                          >
                            <v-icon>mdi-content-copy</v-icon>
                          </v-btn>
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                    
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>XRP Balance</v-list-item-subtitle>
                        <v-list-item-title>
                          {{ formatXRPBalance(accountData.xrpBalance) }} XRP
                          <span class="caption grey--text ml-2">
                            ≈ ${{ formatUSD(accountData.xrpBalance * (accountData.xrpPrice || 0)) }}
                          </span>
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </v-col>
                
                <v-col cols="12" md="6">
                  <v-list dense>
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Total Value</v-list-item-subtitle>
                        <v-list-item-title class="text-h6">
                          ${{ formatUSD(totalValue) }}
                        </v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                    
                    <v-list-item>
                      <v-list-item-content>
                        <v-list-item-subtitle>Token Count</v-list-item-subtitle>
                        <v-list-item-title>
                          {{ accountData.xrpTokens?.length || 0 }} tokens
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

      <!-- Token Balances -->
      <v-row v-if="address && accountData?.xrpTokens?.length">
        <v-col cols="12">
          <v-card tile outlined>
            <v-card-title>
              <v-icon class="mr-2">mdi-currency-usd</v-icon>
              Token Balances
            </v-card-title>
            <v-card-text>
              <v-data-table
                :headers="tokenHeaders"
                :items="accountData.xrpTokens"
                :loading="loading"
                :search="tokenSearch"
                dense
                class="elevation-0"
              >
                <template #top>
                  <v-text-field
                    v-model="tokenSearch"
                    label="Search tokens"
                    prepend-inner-icon="mdi-magnify"
                    outlined
                    dense
                    hide-details
                    class="mx-4 mt-4"
                    style="max-width: 300px"
                  />
                </template>

                <template #item.symbol="{ item }">
                  <div class="d-flex align-center">
                    <v-avatar size="24" class="mr-2">
                      <v-img :src="$imageUrlBySymbol(item.symbol.toLowerCase())" />
                    </v-avatar>
                    <span>{{ item.symbol }}</span>
                  </div>
                </template>

                <template #item.balance="{ item }">
                  {{ formatTokenBalance(item.balance) }}
                </template>

                <template #item.price="{ item }">
                  ${{ formatUSD(item.price) }}
                </template>

                <template #item.value="{ item }">
                  ${{ formatUSD(item.value) }}
                </template>

                <template #item.issuer="{ item }">
                  <div class="d-flex align-center">
                    <span class="font-family-mono caption">{{ formatIssuerAddress(item.issuer) }}</span>
                    <v-btn
                      icon
                      x-small
                      @click="copyToClipboard(item.issuer)"
                      class="ml-1"
                    >
                      <v-icon>mdi-content-copy</v-icon>
                    </v-btn>
                  </div>
                </template>

                <template #item.actions="{ item }">
                  <v-btn
                    small
                    text
                    color="primary"
                    @click="viewTokenDetails(item)"
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
    </v-container>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, useContext } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable'
import { XRPAccountBalancesGQL } from '~/apollo/queries'

export default defineComponent({
  setup() {
    const { $f } = useContext()
    const { query } = useRoute()
    
    // State
    const address = ref<string>('')
    const addressInput = ref<string>('')
    const loading = ref(false)
    const error = ref<string | null>(null)
    const tokenSearch = ref('')

    // Example addresses
    const exampleAddresses = [
      'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo',
      'rDodqfAoF8pVh2SoUwhQRfvkqrs4wwxUrz'
    ]

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

    // GraphQL query
    const { result, loading: queryLoading, error: queryError, refetch } = useQuery(
      XRPAccountBalancesGQL,
      computed(() => ({ account: address.value })),
      () => ({
        enabled: !!address.value,
        fetchPolicy: 'cache-and-network',
        errorPolicy: 'all'
      })
    )

    // Computed data
    const accountData = computed(() => result.value?.xrpAccountBalances || null)
    
    const totalValue = computed(() => {
      if (!accountData.value) return 0
      const xrpValue = (accountData.value.xrpBalance || 0) * (accountData.value.xrpPrice || 0)
      const tokenValue = (accountData.value.xrpTokens || []).reduce((sum, token) => sum + (token.value || 0), 0)
      return xrpValue + tokenValue
    })

    // Table headers
    const tokenHeaders = [
      { text: 'Token', value: 'symbol', sortable: true },
      { text: 'Balance', value: 'balance', sortable: true, align: 'right' },
      { text: 'Price', value: 'price', sortable: true, align: 'right' },
      { text: 'Value', value: 'value', sortable: true, align: 'right' },
      { text: 'Issuer', value: 'issuer', sortable: false },
      { text: 'Actions', value: 'actions', sortable: false, align: 'center' },
    ]

    // Methods
    const loadAddress = () => {
      if (isValidAddress.value) {
        address.value = addressInput.value
        loadAddressData()
      }
    }

    const loadExampleAddress = (exampleAddress: string) => {
      addressInput.value = exampleAddress
      address.value = exampleAddress
      loadAddressData()
    }

    const loadAddressData = async () => {
      if (!address.value) return
      
      loading.value = true
      error.value = null
      
      try {
        await refetch()
        if (queryError.value) {
          error.value = 'Failed to load account data'
        }
      } catch (err: any) {
        error.value = err.message || 'Failed to load account data'
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

    const viewTokenDetails = (token: any) => {
      // Navigate to token details page
      window.$nuxt.$router.push(`/token/${token.symbol}?issuer=${token.issuer}`)
    }

    // Formatting functions
    const formatXRPBalance = (balance: number) => {
      return $f(balance || 0, { minDigits: 6, after: '' })
    }

    const formatTokenBalance = (balance: number) => {
      return $f(balance || 0, { minDigits: 2, after: '' })
    }

    const formatUSD = (value: number) => {
      return $f(value || 0, { minDigits: 2, after: '' })
    }

    const formatIssuerAddress = (address: string) => {
      if (!address) return ''
      return address.length > 20 ? `${address.substring(0, 10)}...${address.substring(address.length - 10)}` : address
    }

    const formatExampleAddress = (address: string) => {
      if (!address) return ''
      return address.length > 20 ? `${address.substring(0, 10)}...${address.substring(address.length - 10)}` : address
    }

    return {
      // State
      address,
      addressInput,
      loading: computed(() => loading.value || queryLoading.value),
      error,
      tokenSearch,
      
      // Computed
      isValidAddress,
      accountData,
      totalValue,
      tokenHeaders,
      
      // Methods
      addressRule,
      loadAddress,
      refreshData,
      copyToClipboard,
      viewTokenDetails,
      
      // Formatting
      formatXRPBalance,
      formatTokenBalance,
      formatUSD,
      formatIssuerAddress,
      
      // Example addresses
      exampleAddresses,
      loadExampleAddress,
      formatExampleAddress,
    }
  },
})
</script> 