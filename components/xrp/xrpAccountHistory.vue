<template>
    <div>
        <v-card tile outlined>
            <v-card-title class="text-h5">
                <v-avatar size="32" class="mr-2">
                    <v-img :src="$imageUrlBySymbol('xrp')" />
                </v-avatar>
                XRP Account History
                <v-spacer />
                <v-chip v-if="totalTransactions > 0" color="primary" small>
                    {{ totalTransactions }} transactions
                </v-chip>
            </v-card-title>
            
            <v-card-text>
                <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />
                
                <client-only>
                    <v-data-table
                        v-if="!loading"
                        hide-default-footer
                        :headers="cols"
                        :items="transactionDataFormatted"
                        :items-per-page="25"
                        class="elevation-0 row-height-50"
                        mobile-breakpoint="0"
                        dense
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
                        
                        <template #item.from="{ item }">
                            <div class="d-flex align-center">
                                <span class="font-family-mono">{{ formatAddress(item.from) }}</span>
                                <v-btn
                                    icon
                                    x-small
                                    class="ml-1"
                                    @click="copyToClipboard(item.from)"
                                >
                                    <v-icon size="16">mdi-content-copy</v-icon>
                                </v-btn>
                            </div>
                        </template>
                        
                        <template #item.to="{ item }">
                            <div class="d-flex align-center">
                                <span class="font-family-mono">{{ formatAddress(item.to) }}</span>
                                <v-btn
                                    icon
                                    x-small
                                    class="ml-1"
                                    @click="copyToClipboard(item.to)"
                                >
                                    <v-icon size="16">mdi-content-copy</v-icon>
                                </v-btn>
                            </div>
                        </template>
                        
                        <template #item.action="{ item }">
                            <v-chip
                                :color="getActionColor(item.action)"
                                small
                                text-color="white"
                            >
                                {{ item.action }}
                            </v-chip>
                        </template>
                        
                        <template #item.amount="{ item }">
                            <span class="font-weight-medium">{{ formatAmount(item.amount) }}</span>
                        </template>
                        
                        <template #item.fee="{ item }">
                            <span class="grey--text">{{ formatFee(item.fee) }}</span>
                        </template>
                    </v-data-table>
                </client-only>
                
                <div v-if="!loading && transactionDataFormatted.length === 0" class="text-center pa-4">
                    <v-icon size="64" color="grey lighten-2">mdi-history</v-icon>
                    <div class="text-h6 grey--text mt-2">No transactions found</div>
                    <div class="text-caption grey--text">Connect your XRP wallet to view transaction history</div>
                </div>
            </v-card-text>
        </v-card>
    </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, inject, onMounted } from '@nuxtjs/composition-api'
import { plainToClass } from 'class-transformer'
import { XRPAccountTransactionsGQL } from '~/apollo/queries'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'

interface XRPTransactionElem {
    hash: string
    from: string
    action: string
    to: string
    amount: number
    fee: number
    timestamp?: string
    ledgerIndex?: number
}

export default defineComponent({
    setup() {
        const { $f } = useContext()
        const { address, isWalletReady } = inject(XRP_PLUGIN_KEY) as XrpClient
        
        const loading = ref(true)
        const transactionRawData = ref<XRPTransactionElem[]>([])
        
        // Use connected wallet address or fallback to default
        const accountAddress = computed(() => address.value || 'rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj')
        
        // Log the query content BEFORE making the call
        console.log('🚀 [BEFORE QUERY] xrpAccountHistory - XRPAccountTransactionsGQL:', {
            query: XRPAccountTransactionsGQL.loc?.source.body,
            variables: { address: accountAddress.value },
            timestamp: new Date().toISOString()
        })

        const { useLoggedQuery } = useXrpGraphQLWithLogging()
        const { onResult } = useLoggedQuery(
            XRPAccountTransactionsGQL, 
            () => ({ address: accountAddress.value }), 
            { 
                fetchPolicy: 'no-cache',
                context: {
                    queryName: 'XRPAccountTransactions',
                    component: 'xrpAccountHistory',
                    purpose: 'XRP account transaction history'
                }
            }
        )

        const transactionDataFormatted = computed(() =>
            transactionRawData.value.map((elem) => ({
                ...elem,
                amountFormatted: formatAmount(elem.amount),
                feeFormatted: formatFee(elem.fee),
            }))
        )

        const totalTransactions = computed(() => transactionRawData.value.length)

        const cols = computed(() => {
            return [
                {
                    text: 'Hash',
                    align: 'left',
                    value: 'hash',
                    width: '200',
                    class: ['px-4', 'text-truncate'],
                    cellClass: ['px-4', 'text-truncate'],
                    sortable: false,
                },
                {
                    text: 'From',
                    align: 'left',
                    value: 'from',
                    width: '150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text: 'Action',
                    align: 'center',
                    value: 'action',
                    width: '120',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: true,
                },
                {
                    text: 'To',
                    align: 'left',
                    value: 'to',
                    width: '150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text: 'Amount',
                    align: 'right',
                    value: 'amount',
                    width: '120',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: true,
                },
                {
                    text: 'Fee',
                    align: 'right',
                    value: 'fee',
                    width: '100',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: true,
                },
            ]
        })

        const formatHash = (hash: string): string => {
            if (hash.length <= 16) return hash
            return `${hash.slice(0, 8)}...${hash.slice(-8)}`
        }

        const formatAddress = (address: string): string => {
            if (address.length <= 16) return address
            return `${address.slice(0, 8)}...${address.slice(-8)}`
        }

        const formatAmount = (amount: number): string => {
            if (amount === 0) return '0'
            if (amount < 0.000001) return '< 0.000001'
            return amount.toFixed(6)
        }

        const formatFee = (fee: number): string => {
            return `${fee} drops`
        }

        const getActionColor = (action: string): string => {
            switch (action.toLowerCase()) {
                case 'payment':
                    return 'green'
                case 'offercreate':
                    return 'blue'
                case 'offercancel':
                    return 'orange'
                case 'trustset':
                    return 'purple'
                case 'accountset':
                    return 'grey'
                default:
                    return 'primary'
            }
        }

        const copyToClipboard = async (text: string) => {
            try {
                await navigator.clipboard.writeText(text)
                // You could add a toast notification here
            } catch (err) {
                console.error('Failed to copy text: ', err)
            }
        }

        const openTransaction = (hash: string) => {
            window.open(`/xrp-explorer/tx/${hash}`, '_blank')
        }

        // Handle query results
        onResult((queryResult: any) => {
            if (queryResult.data?.xrpTransactions) {
                const transactions = queryResult.data.xrpTransactions
                
                // Transform GraphQL data to component format
                const transformedTransactions: XRPTransactionElem[] = transactions.map((tx: any) => ({
                    hash: tx.hash || '',
                    from: tx.account || '',
                    to: tx.destination || '',
                    action: tx.transactionType || 'Unknown',
                    amount: typeof tx.amount === 'string' ? parseFloat(tx.amount) : (tx.amount || 0),
                    fee: typeof tx.fee === 'string' ? parseFloat(tx.fee) : (tx.fee || 0),
                    timestamp: tx.date ? new Date(tx.date).toISOString() : new Date().toISOString(),
                    ledgerIndex: tx.ledgerIndex || 0
                }))
                
                transactionRawData.value = transformedTransactions
            } else {
                // No transactions found
                transactionRawData.value = []
            }
            loading.value = false
        })

        // Add error handling
        const { onError } = useLoggedQuery(
            XRPAccountTransactionsGQL, 
            () => ({ address: accountAddress.value }), 
            { 
                fetchPolicy: 'no-cache',
                context: {
                    queryName: 'XRPAccountTransactions',
                    component: 'xrpAccountHistory',
                    purpose: 'XRP account transaction history'
                }
            }
        )

        onError((error: any) => {
            console.error('GraphQL Error in xrpAccountHistory:', error)
            loading.value = false
            // You could add error state handling here
        })

        // Watch for wallet address changes
        onMounted(() => {
            // Initial load
        })

        return {
            cols,
            transactionDataFormatted,
            loading,
            totalTransactions,
            formatHash,
            formatAddress,
            formatAmount,
            formatFee,
            getActionColor,
            copyToClipboard,
            openTransaction,
            isWalletReady
        }
    }
})
</script>

<style scoped>
.font-family-mono {
    font-family: 'Courier New', monospace;
}
</style>