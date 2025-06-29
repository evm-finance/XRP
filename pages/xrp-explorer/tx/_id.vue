<template>
  <div>
    <v-row justify="center">
      <v-col lg="9" md="12">
        <v-row no-gutters>
          <v-col><h1 class="text-h4">Transaction summary</h1></v-col>
        </v-row>

        <client-only>
          <v-row v-if="txData" justify="center" class="my-1">
            <v-col lg="8">
              <v-card tile outlined height="100%">
                <v-simple-table class="">
                  <template #default>
                    <tbody>
                      <tr>
                        <td><span class="font-weight-bold">Specification</span></td>
                        <td>
                          <v-chip class="ma-2" label small outlined color="green">
                            {{ txData.transactionType }}
                          </v-chip>
                        </td>
                      </tr>

                      <tr>
                        <td>Tx hash:</td>
                        <td class="grey--text">{{ txData.hash }}</td>
                      </tr>

                      <tr>
                        <td>Date:</td>
                        <td class="grey--text">{{ txData.date }}</td>
                      </tr>
                      <tr>
                        <td>Source:</td>
                        <td class="grey--text">{{ txData.account }}</td>
                      </tr>

                      <tr>
                        <td>Offer:</td>
                        <td v-if="takerGets" class="grey--text">
                          {{ takerGets.value }}<span class="mx-2">{{ takerGets.currency }}</span>
                        </td>
                      </tr>

                      <tr>
                        <td>Offer:</td>
                        <td v-if="txData.takerPays" class="grey--text">{{ parseFloat(txData.takerPays) / 10 ** 6 }}</td>
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
                        <td><span class="font-weight-bold">Ledger</span></td>
                        <td>
                          <v-chip class="ma-2" label small outlined color="green">
                            {{ txData.meta?.TransactionResult ?? '-' }}
                          </v-chip>
                        </td>
                      </tr>

                      <tr>
                        <td><span>Ledger:</span></td>
                        <td>
                          <nuxt-link :to="`/xrp-explorer/ledger/${txData.ledgerIndex}`">
                            {{ txData.ledgerIndex }}
                          </nuxt-link>
                        </td>
                      </tr>

                      <tr>
                        <td>Index:</td>
                        <td class="grey--text">{{ txData.flags }}</td>
                      </tr>

                      <tr>
                        <td>Tx seq::</td>
                        <td class="grey--text">{{ txData.sequence }}</td>
                      </tr>
                      <tr>
                        <td>Fee:</td>
                        <td class="grey--text">{{ parseFloat(txData.fee) / 10 ** 6 }}</td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-card>
            </v-col>

            <v-col v-for="(nodeHeader, i) in txData.meta.AffectedNodes" :key="i" cols="12">
              <v-card v-for="(elems, header) in nodeHeader" :key="header" tile outlined>
                <v-card-title>{{ header }}</v-card-title>
                <v-simple-table class="">
                  <template #default>
                    <tbody>
                      <tr v-for="(k, v) in elems" :key="v">
                        <td>
                          <span>{{ v }}</span>
                        </td>
                        <td>
                          <div v-if="v === 'LedgerEntryType'">
                            <v-chip class="ma-2" label small outlined color="primary">
                              {{ k }}
                            </v-chip>
                          </div>
                          <div v-else class="grey--text">{{ k }}</div>
                        </td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-card>
            </v-col>
          </v-row>
        </client-only>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useRoute } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'
import { XRPTransactionGQL } from '~/apollo/queries'
import { XrpTransaction } from '~/types/apollo/main/types'

export default defineComponent({
  components: {},
  setup() {
    const route = useRoute()
    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    
    const hash = computed(() => route.value.params?.id ?? '')
    
    // Enhanced GraphQL query with logging
    const { result, error } = process.browser
      ? useLoggedQuery(XRPTransactionGQL, () => ({
          hash: hash.value,
        }))
      : { result: ref(null), error: ref(null) }
      
    const txData = computed<XrpTransaction>(() => result.value?.xrpTransaction ?? null)
    const takerGets = computed(() => txData.value?.takerGets ?? null)

    // Enhanced error handling
    if (error.value) {
      console.error('XRP Transaction query failed:', error.value)
    }

    return { txData, takerGets }
  },
  head: {},
})
</script>

<style scoped></style>
