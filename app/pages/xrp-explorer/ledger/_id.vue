<template>
  <div>
    <v-row justify="center">
      <v-col lg="9" md="12">
        <v-row no-gutters>
          <v-col><h1 class="text-h4">Ledger Summary</h1></v-col>
        </v-row>
        <client-only>
          <v-row v-if="ledger" justify="center" class="my-1">
            <v-col lg="8">
              <v-card tile outlined height="100%">
                <v-simple-table class="">
                  <template #default>
                    <tbody>
                      <tr>
                        <td><span class="text-h6">Ledger Index</span></td>
                        <td class="grey--text">
                          <v-btn icon small class="mb-1" @click="navigateToLedger(ledger.blockNumber - 1)">
                            <v-icon size="16">mdi-arrow-left</v-icon>
                          </v-btn>
                          <span class="text-h6 px-2 green--text">{{ ledger.blockNumber }}</span>
                          <v-btn icon small class="mb-1" @click="navigateToLedger(ledger.blockNumber + 1)">
                            <v-icon size="16">mdi-arrow-right</v-icon>
                          </v-btn>
                        </td>
                      </tr>

                      <tr>
                        <td>Ledger hash:</td>
                        <td class="grey--text">{{ ledger.XRPLedger.ledgerHash }}</td>
                      </tr>

                      <tr>
                        <td>Parent hash:</td>
                        <td class="grey--text">{{ ledger.XRPLedger.ledger.parentHash }}</td>
                      </tr>
                      <tr>
                        <td>Txs hash:</td>
                        <td class="grey--text">{{ ledger.XRPLedger.ledger.transactionHash }}</td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-card>
            </v-col>

            <v-col lg="4">
              <v-card tile outlined height="100%">
                <v-simple-table class="">
                  <template #default>
                    <tbody>
                      <tr>
                        <td><span class="text-h6">Properties</span></td>
                        <td></td>
                      </tr>

                      <tr>
                        <td>Closed on:</td>
                        <td class="grey--text">{{ ledger.XRPLedger.ledger.closeTimeHuman }}</td>
                      </tr>

                      <tr>
                        <td>Total coins:</td>
                        <td class="grey--text">
                          {{ $f(parseFloat(ledger.XRPLedger.ledger.totalCoins) / 100, { maxDigits: 0 }) }} XRP
                        </td>
                      </tr>
                      <tr>
                        <td>Σ Fee::</td>
                        <td class="grey--text">-</td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-card>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-card tile outlined>
                <v-card-title> <span class="text-h6">Transaction</span></v-card-title>
                <v-divider />
                <v-data-table
                  v-if="ledger"
                  hide-default-footer
                  :headers="cols"
                  :items-per-page="200"
                  :items="transactions"
                  class="elevation-0 row-height-50"
                  mobile-breakpoint="0"
                >
                  <template #item.hash="{ item }">
                    <nuxt-link :to="`/xrp-explorer/tx/${item.hash}`" class="pink--text">{{ item.hashShort }}</nuxt-link>
                  </template>

                  <template #item.account="{ item }">
                    <div class="grey--text">{{ item.accountShort }}</div>
                  </template>

                  <template #item.transactionType="{ item }">
                    <v-chip class="ma-2" label small outlined :color="eventColor(item.transactionType)">
                      {{ item.transactionType }}
                    </v-chip>
                  </template>

                  <template #item.destination="{ item }">
                    <div class="grey--text">{{ item.toFormatted }}</div>
                  </template>
                  <template #item.amount="{ item }">
                    <div>
                      {{ item.amountFormatted.value }}
                      <span class="grey--text ml-2">{{ item.amountFormatted.currency }}</span>
                    </div>
                  </template>
                </v-data-table>
              </v-card>
            </v-col>
          </v-row>
        </client-only>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, useRoute, useRouter, ref, watch } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'
import { BlockGQL } from '~/apollo/main/token.query.graphql'
import { Block, XrpTransaction } from '~/types/apollo/main/types'
export default defineComponent({
  components: {},
  setup() {
    const route = useRoute()
    const router = useRouter()
    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    
    const ledgerIndex = computed(() => route.value.params?.id ?? 0)
    
    // Internal error handling - prevent Apollo errors from bubbling up
    const internalError = ref<string | null>(null)
    
    const { result, error } = useLoggedQuery(
      BlockGQL, 
      () => ({ network: 'ripple', blockNumber: ledgerIndex.value }), 
      {
        fetchPolicy: 'no-cache',
        context: {
          queryName: 'Block',
          component: 'xrp-explorer-ledger',
          purpose: 'XRP ledger block data'
        }
      }
    )
    
    // Watch for Apollo errors and handle them internally
    watch(error, (newError) => {
      if (newError) {
        console.warn('GraphQL error in xrp-explorer-ledger, using fallback data:', newError)
        internalError.value = 'Network error: using fallback data'
      } else {
        internalError.value = null
      }
    })
    
    // Error state for graceful fallback
    const hasError = computed(() => !!internalError.value || !!error.value)
    
    const ledger = computed<Block | null>(() => result.value?.block ?? null)

    const transactions = computed(() =>
      ledger.value?.XRPTransactions.items
        .map((elem) => ({
          ...elem,
          hashShort: `${elem.hash.slice(0, 10)}........${elem.hash.slice(elem.hash.length - 10, elem.hash.length)}`,
          accountShort: `${elem.account.slice(0, 20)}....`,
          toFormatted: elem?.destination ? `${elem.destination.slice(0, 20)}....` : 'XRPL',
          amountFormatted: handleAmount(elem),
          fee: `${parseFloat(elem.fee) / 10 ** 6}  XRP`,
        }))
        .sort((a, b) => {
          if (a.metadata.TransactionIndex < b.metadata.TransactionIndex) return -1
          if (a.metadata.TransactionIndex > b.metadata.TransactionIndex) return 1
          return 0
        })
    )

    function handleAmount(data: XrpTransaction) {
      if (typeof data.amount === 'string') {
        return { value: parseFloat(data.amount) / 10 ** 6, currency: 'XRP' }
      }
      // if (data.amount === null) {
      //   const affectedNodes = data.metadata.AffectedNodes
      //   for (const node of affectedNodes) {
      //     if (node.CreatedNode && node.CreatedNode.NewFields) {
      //       const takerGets = node.CreatedNode.NewFields.TakerGets
      //       console.log('TakerGets:', takerGets)
      //       if (typeof takerGets === 'object') {
      //         return { value: parseFloat(takerGets?.value ?? 0), currency: takerGets?.currency ?? '' }
      //       } else {
      //         continue
      //       }
      //     } else if (node.ModifiedNode && node.ModifiedNode.NewFields) {
      //       const takerGets = node.ModifiedNode.NewFields.TakerGets
      //       console.log('TakerGets:', takerGets)
      //     }
      //   }
      // }
      return { value: '-', currency: '' }
    }

    const navigateToLedger = (ledger: number) => router.push(`/xrp-explorer/ledger/${ledger}`)
    const cols = computed(() => {
      return [
        {
          text: 'Index',
          align: 'left',
          value: 'metadata.TransactionIndex',
          width: '20',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate', 'grey--text'],
          sortable: false,
        },

        {
          text: 'Hash',
          align: 'left',
          value: 'hash',
          width: '100',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate', 'grey--text'],
          sortable: false,
        },

        {
          text: 'From',
          align: 'left',
          value: 'account',
          width: '100',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate', 'grey--text'],
          sortable: false,
        },

        {
          text: '',
          align: 'left',
          value: 'transactionType',
          class: ['px-2', 'text-truncate'],
          width: '100',
          cellClass: ['px-2', 'text-truncate'],
          sortable: false,
        },

        {
          text: 'To',
          align: 'left',
          value: 'destination',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
          width: '100',
          sortable: false,
        },

        {
          text: 'Amount',
          align: 'left',
          value: 'amount',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          width: '100',
          sortable: false,
        },
        {
          text: 'Amount',
          align: 'left',
          value: 'fee',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
          width: '100',
          sortable: false,
        },
      ]
    })
    const eventColor = (event: any): string => {
      const colors: { [key: string]: string } = {
        OfferCreate: 'green',
        TicketCreate: 'green',
        Payment: 'green',
        NFTokenCreateOffer: 'green',
        NFTokenCancelOffer: 'red ',
        NFTokenMint: 'primary',
        OfferCancel: 'red',
        TrustSet: 'orange',
      }

      // @ts-ignore
      return Object.hasOwn(colors, event) ? colors[event] : 'grey'
    }

    return { ledger, navigateToLedger, cols, transactions, eventColor }
  },
  head: {},
})
</script>

<style scoped></style>
