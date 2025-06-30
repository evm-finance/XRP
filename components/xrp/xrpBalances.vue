<template>
  <div>
    <v-card tile outlined>
      <v-card-title class="text-h5">
        <v-avatar size="32" class="mr-2">
          <v-img :src="$imageUrlBySymbol('xrp')" />
        </v-avatar>
        XRP Account Balances
        <v-spacer />
        <v-chip v-if="totalValue > 0" color="primary" small>
          Total: ${{ formatCurrency(totalValue) }}
        </v-chip>
      </v-card-title>
      
      <v-card-text>
        <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />
        
        <client-only>
          <v-data-table
            v-if="!loading"
            hide-default-footer
            :headers="cols"
            :items="screenerDataFormatted"
            :items-per-page="25"
            class="elevation-0 row-height-50"
            mobile-breakpoint="0"
            dense
          >
            <template #item.currency="{ item }">
              <div class="d-flex align-center">
                <v-avatar size="24" class="mr-2">
                  <v-img :src="$imageUrlBySymbol(item.currency.toLowerCase())" />
                </v-avatar>
                <span class="font-weight-medium">{{ item.currency }}</span>
              </div>
            </template>
            
            <template #item.issuerAddressShort="{ item }">
              <div class="d-flex align-center">
                <span class="font-family-mono">{{ item.issuerAddressShort }}</span>
                <v-btn
                  icon
                  x-small
                  class="ml-1"
                  @click="copyToClipboard(item.issuer)"
                >
                  <v-icon size="16">mdi-content-copy</v-icon>
                </v-btn>
              </div>
            </template>
            
            <template #item.balance="{ item }">
              <span class="font-weight-medium">{{ formatBalance(item.balance) }}</span>
            </template>
            
            <template #item.value="{ item }">
              <span class="font-weight-medium">${{ formatCurrency(item.value) }}</span>
            </template>
          </v-data-table>
        </client-only>
        
        <div v-if="!loading && screenerDataFormatted.length === 0" class="text-center pa-4">
          <v-icon size="64" color="grey lighten-2">mdi-wallet-outline</v-icon>
          <div class="text-h6 grey--text mt-2">No balances found</div>
          <div class="text-caption grey--text">Connect your XRP wallet to view balances</div>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, inject, onMounted } from '@nuxtjs/composition-api'
import { plainToClass } from 'class-transformer'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPAccountBalancesGQL } from '~/apollo/queries'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'

//import { State } from '~/types/state'
//import { result } from '~/composables/useXrpAccounts'

interface XRPBalanceElem {
    issuer: string
    currency: string
    name: string
    balance: number
    price: number
    value: number
}

export default defineComponent({
    setup() {
        const { $f } = useContext()
        const { address, isWalletReady } = inject(XRP_PLUGIN_KEY) as XrpClient
        
        const loading = ref(true)
        const balancesRawData = ref<XRPBalanceElem[]>([])
        
        // Use connected wallet address or fallback to primary test address
        const accountAddress = computed(() => address.value || 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo')
        
        // Log the query content BEFORE making the call
        console.log('🚀 [BEFORE QUERY] xrpBalances - XRPAccountBalancesGQL:', {
            query: XRPAccountBalancesGQL.loc?.source.body,
            variables: { address: accountAddress.value },
            timestamp: new Date().toISOString(),
            component: 'xrpBalances.vue',
            graphqlEndpoint: process.env.BASE_GRAPHQL_SERVER_URL || 'http://127.0.0.1:8080/query'
        })

        const { useLoggedQuery } = useXrpGraphQLWithLogging()
        const { onResult } = useLoggedQuery(
            XRPAccountBalancesGQL, 
            () => ({ address: accountAddress.value }), 
            { 
                fetchPolicy: 'no-cache',
                context: {
                    queryName: 'XRPAccountBalances',
                    component: 'xrpBalances',
                    purpose: 'XRP account balances display'
                }
            }
        )

        const screenerDataFormatted = computed(() =>
            balancesRawData.value.map((elem) => ({
            ...elem,
            currencyShort: elem.currency.length > 20 ? elem.currency.substring(0, 20) + '...' : elem.currency,
            issuerAddressShort: `${elem.issuer.slice(0, 10)}........${elem.issuer.slice(
              elem.issuer.length - 10,
              elem.issuer.length
            )}`,
            priceFormatted: $f(elem.price, { minDigits: 6, after: '' }),
          }))
        )

        const totalValue = computed(() => 
            balancesRawData.value.reduce((sum, item) => sum + item.value, 0)
        )

        const cols = computed(() => {
            return [
            {
          text: 'Currency',
          align: 'left',
          value: 'currency',
          width: '150',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate'],
          sortable: true,
        },
        {
          text: 'Token Name',
          align: 'left',
          value: 'name',
          width: '120',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Issuer Address',
          align: 'left',
          value: 'issuerAddressShort',
          width: '150',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
        },
        {
          text: 'Price',
          align: 'right',
          value: 'priceFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Balance',
          align: 'right',
          value: 'balance',
          width: '120',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Value',
          align: 'right',
          value: 'value',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
            ]
        })
    
        const formatBalance = (balance: number): string => {
            if (balance === 0) return '0'
            if (balance < 0.000001) return '< 0.000001'
            return balance.toFixed(6)
        }

        const formatCurrency = (value: number): string => {
            return $f(value, { minDigits: 2, maxDigits: 2 })
        }

        const copyToClipboard = async (text: string) => {
            try {
                await navigator.clipboard.writeText(text)
                // You could add a toast notification here
            } catch (err) {
                console.error('Failed to copy text: ', err)
            }
        }

        // Handle query results
        onResult((queryResult: any) => {
            console.log('🎯 [XRP BALANCES QUERY RESULT]', {
                timestamp: new Date().toISOString(),
                component: 'xrpBalances.vue',
                accountAddress: accountAddress.value,
                hasData: !!queryResult.data,
                hasXrpAccountBalances: !!queryResult.data?.xrpAccountBalances,
                fullResult: queryResult,
                dataStructure: queryResult.data ? Object.keys(queryResult.data) : 'no data',
                loadingState: loading.value
            })

            if (queryResult.data?.xrpAccountBalances) {
                const data = queryResult.data.xrpAccountBalances
                
                console.log('📊 [XRP BALANCES DATA RECEIVED]', {
                    timestamp: new Date().toISOString(),
                    account: data.account,
                    xrpBalance: data.xrpBalance,
                    xrpPrice: data.xrpPrice,
                    tokensCount: data.xrpTokens ? data.xrpTokens.length : 0,
                    tokensData: data.xrpTokens,
                    fullData: data
                })
                
                // Transform GraphQL data to component format
                const transformedData: XRPBalanceElem[] = []
                
                // Add XRP balance if available
                if (data.xrpBalance && data.xrpBalance > 0) {
                    const xrpBalance = {
                        issuer: '',
                        currency: 'XRP',
                        name: 'XRP',
                        balance: data.xrpBalance,
                        price: data.xrpPrice || 0,
                        value: data.xrpBalance * (data.xrpPrice || 0)
                    }
                    transformedData.push(xrpBalance)
                    console.log('💰 [XRP BALANCE ADDED]', xrpBalance)
                }
                
                // Add token balances
                if (data.xrpTokens && Array.isArray(data.xrpTokens)) {
                    data.xrpTokens.forEach((token: any, index: number) => {
                        const tokenBalance = {
                            issuer: token.issuer || '',
                            currency: token.symbol || '',
                            name: token.name || token.symbol || '',
                            balance: token.balance || 0,
                            price: token.price || 0,
                            value: token.value || 0
                        }
                        transformedData.push(tokenBalance)
                        console.log(`🪙 [TOKEN BALANCE ${index + 1} ADDED]`, tokenBalance)
                    })
                }
                
                console.log('✅ [TRANSFORMED DATA COMPLETE]', {
                    totalItems: transformedData.length,
                    totalValue: transformedData.reduce((sum, item) => sum + item.value, 0),
                    transformedData
                })
                
                balancesRawData.value = transformedData
            } else {
                console.log('❌ [NO XRP BALANCE DATA]', {
                    timestamp: new Date().toISOString(),
                    queryResult,
                    message: 'No xrpAccountBalances found in response'
                })
                balancesRawData.value = []
            }
            loading.value = false
        })

        // Add error handling with detailed logging
        const { onError } = useLoggedQuery(
            XRPAccountBalancesGQL, 
            () => ({ address: accountAddress.value }), 
            { 
                fetchPolicy: 'no-cache',
                context: {
                    queryName: 'XRPAccountBalances',
                    component: 'xrpBalances',
                    purpose: 'XRP account balances error handling'
                }
            }
        )

        onError((error) => {
            console.error('🚨 [XRP BALANCES QUERY ERROR]', {
                timestamp: new Date().toISOString(),
                component: 'xrpBalances.vue',
                accountAddress: accountAddress.value,
                errorMessage: error.message,
                networkError: error.networkError,
                graphQLErrors: error.graphQLErrors,
                fullError: error,
                graphqlEndpoint: process.env.BASE_GRAPHQL_SERVER_URL || 'http://127.0.0.1:8080/query'
            })
            loading.value = false
        })

        // Watch for wallet address changes
        onMounted(() => {
            // Initial load
        })

        const error = ref(null)

        return {
            cols, 
            screenerDataFormatted,
            loading,
            totalValue,
            formatBalance,
            formatCurrency,
            copyToClipboard,
            isWalletReady,
            error
        }
    }
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style>