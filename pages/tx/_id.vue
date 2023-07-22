<template>
  <div>
    <v-row justify="center">
      <v-col lg="9" md="12">
        <v-row no-gutters>
          <v-col><h1 class="text-h4">Transaction Details</h1></v-col>
        </v-row>

        <v-row justify="center" class="my-1">
          <v-col>
            <v-card tile outlined>
              <v-simple-table class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="width: 300px">Transaction Hash:</td>
                      <td>
                        0x8444890f84d1b18d099e57c8ef257183e7f8a9e1e1f4bbb52f1026d2265c1af9 <v-icon>mid-copy</v-icon>
                      </td>
                    </tr>

                    <tr>
                      <td class="grey--text">Status:</td>
                      <td><v-chip color="green" small label outlined>Success</v-chip></td>
                    </tr>
                    <tr>
                      <td class="grey--text">Block:</td>
                      <td><v-icon color="green" size="18" class="mr-1">mdi-check-circle-outline</v-icon> 17687982</td>
                    </tr>

                    <tr>
                      <td class="grey--text">Timestamp:</td>
                      <td>
                        <v-icon size="18" class="mr-1"> mdi-clock-time-four-outline</v-icon> 14 mins ago (Jul-13-2023
                        11:55:59 PM +UTC)
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>

            <v-card tile outlined class="my-4">
              <v-simple-table class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="width: 300px">From:</td>
                      <td>0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5</td>
                    </tr>

                    <tr>
                      <td class="grey--text">To:</td>
                      <td>0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5</td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>

            <v-card tile outlined>
              <v-simple-table class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="width: 300px">Value:</td>
                      <td>
                        <v-icon size="16" color="grey">mdi-ethereum</v-icon> 0.157085003937486727 ETH
                        <v-chip color="grey" class="ml-2" label outlined small> $315</v-chip>
                      </td>
                    </tr>

                    <tr>
                      <td class="grey--text">Transaction Fee:</td>
                      <td>
                        0.157085003937486727 ETH <v-chip color="grey" class="ml-2" label outlined small> $315</v-chip>
                      </td>
                    </tr>

                    <tr>
                      <td class="grey--text">Gas Price:</td>
                      <td>27.097807961 Gwei <span class="grey--text mx-2">(0.000000027097807961 ETH)</span></td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>
          </v-col>

          <!--            <v-col lg="4">-->
          <!--              <v-card tile outlined height="100%">-->
          <!--                <v-simple-table class="">-->
          <!--                  <template #default>-->
          <!--                    <tbody>-->
          <!--                      <tr>-->
          <!--                        <td><span class="text-h6">Properties</span></td>-->
          <!--                        <td></td>-->
          <!--                      </tr>-->

          <!--                      <tr>-->
          <!--                        <td>Closed on:</td>-->
          <!--                        <td class="grey&#45;&#45;text">{{ ledger.XRPLedger.ledger.closeTimeHuman }}</td>-->
          <!--                      </tr>-->

          <!--                      <tr>-->
          <!--                        <td>Total coins:</td>-->
          <!--                        <td class="grey&#45;&#45;text">-->
          <!--                          {{ $f(parseFloat(ledger.XRPLedger.ledger.totalCoins) / 100, { maxDigits: 0 }) }} XRP-->
          <!--                        </td>-->
          <!--                      </tr>-->
          <!--                      <tr>-->
          <!--                        <td>Σ Fee::</td>-->
          <!--                        <td class="grey&#45;&#45;text">-</td>-->
          <!--                      </tr>-->
          <!--                    </tbody>-->
          <!--                  </template>-->
          <!--                </v-simple-table>-->
          <!--              </v-card>-->
          <!--            </v-col>-->
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
                :items="ledger.XRPTransactions.items"
                class="elevation-0 row-height-50"
                mobile-breakpoint="0"
              >
                <template #item.hash="{ item }">
                  <span>{{ item.hash }}</span>
                </template>
              </v-data-table>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, useRoute, useRouter } from '@nuxtjs/composition-api'
import { Block } from '~/types/apollo/main/types'
export default defineComponent({
  components: {},
  setup() {
    const route = useRoute()
    const router = useRouter()
    // const ledgerIndex = computed(() => route.value.params?.id ?? 0)
    // const { result } = useQuery(BlockGQL, () => ({ network: 'ripple', blockNumber: ledgerIndex.value }), {
    //   fetchPolicy: 'no-cache',
    // })
    const ledger = computed<Block | null>(() => null)
    const navigateToLedger = (ledger: number) => router.push(`/xrp-explorer/ledger/${ledger}`)
    const cols = computed(() => {
      return [
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
          sortable: false,
          class: ['px-2', 'text-truncate', 'grey--text'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
        },

        {
          text: 'To',
          align: 'left',
          value: 'destination',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
          width: '200',
          sortable: false,
        },

        {
          text: '',
          align: 'left',
          value: 'transactionType',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          sortable: false,
        },
      ]
    })

    onMounted(() => {
      console.log(router, route)
    })

    return { ledger, navigateToLedger, cols }
  },
  head: {},
})
</script>

<style scoped></style>
