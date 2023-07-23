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
              <v-simple-table v-if="tx" class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="width: 300px">Transaction Hash:</td>
                      <td>{{ tx.hash }} <v-icon>mid-copy</v-icon></td>
                    </tr>

                    <tr>
                      <td class="grey--text">Status:</td>
                      <td>
                        <v-chip color="green" small label outlined>{{ tx.isPending }}</v-chip>
                      </td>
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
        <!--        <v-row>-->
        <!--          <v-col cols="12">-->
        <!--            <v-card tile outlined>-->
        <!--              <v-card-title> <span class="text-h6">Transaction</span></v-card-title>-->
        <!--              <v-divider />-->
        <!--              <v-data-table-->
        <!--                v-if="ledger"-->
        <!--                hide-default-footer-->
        <!--                :headers="cols"-->
        <!--                :items-per-page="200"-->
        <!--                :items="ledger.XRPTransactions.items"-->
        <!--                class="elevation-0 row-height-50"-->
        <!--                mobile-breakpoint="0"-->
        <!--              >-->
        <!--                <template #item.hash="{ item }">-->
        <!--                  <span>{{ item.hash }}</span>-->
        <!--                </template>-->
        <!--              </v-data-table>-->
        <!--            </v-card>-->
        <!--          </v-col>-->
        <!--        </v-row>-->
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { EvmTransactionGQL } from '~/apollo/queries'
import { EvmTransaction } from '~/types/graph'
export default defineComponent({
  components: {},
  setup() {
    const { result, error, onError } = useQuery(
      EvmTransactionGQL,
      () => ({ chainId: 1, hash: '0x100364c3733e65f044bac70534e8b17187ed64fa89d8b45883ff081f67379178' }),
      { fetchPolicy: 'no-cache' }
    )
    onError((e) => {
      console.log('EEEEEE', e)
    })
    console.log(result.value, error.value)
    const txData = computed<EvmTransaction>(() => result.value?.evmTransaction ?? null)
    const txDataFormatted = computed(() => ({ ...txData.value, some: 'value' }))
    return { tx: txDataFormatted }
  },
  head: {},
})
</script>

<style scoped></style>
