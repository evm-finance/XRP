<template>
  <div>
    <v-row justify="center">
      <v-col lg="9" md="12">
        <v-row no-gutters>
          <v-col><h1 class="text-h4">Transaction Details</h1></v-col>
        </v-row>
        <v-row>
          <v-btn-toggle v-model="defaultView" color="primary" class="ml-3">
            <v-btn v-for="item in viewOptions" :key="item.value" tile outlined small :value="item.value">
              {{ item.text }}
            </v-btn>
          </v-btn-toggle>
        </v-row>

        <v-row v-if="tx" justify="center" class="my-1">
          <v-col v-if="defaultView === 'overview'">
            <v-card tile outlined>
              <v-card-title>Overview</v-card-title>
              <v-simple-table class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="width: 300px">Transaction Hash:</td>
                      <td>{{ tx.hash }} <v-icon>mid-copy</v-icon></td>
                    </tr>

                    <tr>
                      <td class="grey--text">Status:</td>
                      <td>
                        <v-chip color="green" small label outlined>{{ tx.status }}</v-chip>
                      </td>
                    </tr>
                    <tr>
                      <td class="grey--text">Block:</td>
                      <td>
                        <v-icon color="green" size="18" class="mr-1">mdi-check-circle-outline</v-icon> {{ tx.block }}
                      </td>
                    </tr>

                    <tr>
                      <td class="grey--text">Timestamp:</td>
                      <td><v-icon size="18" class="mr-2"> mdi-clock-time-four-outline</v-icon>{{ tx.dateTime }}</td>
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
                      <td>{{ tx.from }}</td>
                    </tr>

                    <tr>
                      <td class="grey--text">To:</td>
                      <td>{{ tx.to }}</td>
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
                      <td><v-icon size="16" color="grey">mdi-ethereum</v-icon> {{ tx.value }}</td>
                    </tr>

                    <tr>
                      <td class="grey--text">Transaction Fee:</td>
                      <td>{{ tx.txFeeInETH }}</td>
                    </tr>

                    <tr>
                      <td class="grey--text">Gas Price:</td>
                      <td>
                        {{ tx.gasPriceInGwei }}<span class="grey--text mx-2">({{ tx.gasPriceInETH }})</span>
                      </td>
                    </tr>

                    <tr>
                      <td class="grey--text">Gas Limit:</td>
                      <td>{{ tx.gssLimit }}</td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>

            <v-card v-if="tx.input" tile outlined class="my-4">
              <v-simple-table class="">
                <template #default>
                  <tbody>
                    <tr>
                      <td class="grey--text" style="min-width: 300px">
                        <div>Input:</div>
                      </td>
                      <td class="pa-2">
                        <div style="background: #202121; height: 300px" class="pa-4 overflow-y-auto">
                          <div v-if="currentInputOption === 'original'" style="width: 1000px">
                            {{ tx.input.methodSigDataStr + tx.input.inputsSigDataStr }}
                          </div>
                          <div v-if="currentInputOption === 'default'" style="width: 1000px">
                            <div>{{ tx.input.fullFunctionSig }}</div>
                            <div class="my-2">Method ID: {{ '0x' + tx.input.methodSigDataStr }}</div>
                            <div v-for="(item, i) in tx.dataToArrayOfStrings" :key="i">
                              <span class="mr-1">[{{ i }}]:</span> {{ item }}
                            </div>
                          </div>
                        </div>

                        <div>
                          <v-menu offset-y>
                            <template #activator="{ attrs, on }">
                              <v-btn class="white--text mt-3" v-bind="attrs" depressed tile small v-on="on">
                                View Input As <v-icon small class="ml-1">mdi-chevron-down</v-icon>
                              </v-btn>
                            </template>
                            <v-list dense>
                              <v-list-item-group v-model="currentInputOption" color="primary">
                                <v-list-item
                                  v-for="item in inputOptions"
                                  :key="item.value"
                                  :value="item.value"
                                  link
                                  dense
                                >
                                  <v-list-item-title>{{ item.text }}</v-list-item-title>
                                </v-list-item>
                              </v-list-item-group>
                            </v-list>
                          </v-menu>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>
          </v-col>

          <v-col v-if="defaultView === 'input'">
            <v-card v-if="tx.input" tile outlined>
              <v-card-title>Input Details</v-card-title>
              <v-simple-table>
                <template #default>
                  <thead>
                    <tr>
                      <th class="text-left">Name</th>
                      <th class="text-left">Type</th>
                      <th class="text-left">Value</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in inputDetails" :key="item.name">
                      <td>{{ item.name }}</td>
                      <td>{{ item.type }}</td>
                      <td>{{ item.value }}</td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </v-card>
          </v-col>

          <v-col v-if="defaultView === 'log'">
            <v-card v-if="tx.logEvents" tile outlined>
              <v-card-title>Logs Events</v-card-title>
              <v-card v-for="(item, i) in logEvents" :key="i" tile outlined class="mb-6">
                <v-simple-table class="">
                  <template #default>
                    <tbody>
                      <tr>
                        <td class="grey--text" style="width: 300px">Address</td>
                        <td>
                          <span class="primary--text">{{ item.address }}</span>
                        </td>
                      </tr>

                      <tr>
                        <td class="grey--text">Name</td>
                        <td>
                          <span class="pink--text">{{ item.name }}</span>
                        </td>
                      </tr>

                      <tr>
                        <td class="grey--text">Signature</td>
                        <td>{{ item.signature }}</td>
                      </tr>

                      <tr>
                        <td class="grey--text">Topics</td>
                        <td>
                          <div class="py-1">
                            <div>
                              <v-btn class="white--text px-0 mr-2" depressed tile small width="10">0</v-btn>
                              {{ item.topic }}
                            </div>
                            <div v-for="(t, idx) in item.indexedVals" :key="idx" class="py-1">
                              <v-btn class="white--text px-0 mr-2" depressed tile small width="10">{{ idx + 1 }}</v-btn>
                              {{ t }}
                            </div>
                          </div>
                        </td>
                      </tr>

                      <tr>
                        <td class="grey--text">Data</td>
                        <td>
                          <div class="py-2">
                            <v-simple-table v-if="currentInputOption === 'default'" dense>
                              <template #default>
                                <thead>
                                  <tr>
                                    <th class="text-left">Name</th>
                                    <th class="text-left">Type</th>
                                    <th class="text-left">Value</th>
                                  </tr>
                                </thead>
                                <tbody>
                                  <tr v-for="f in item.functionParams" :key="f.name">
                                    <td>{{ f.name }}</td>
                                    <td class="grey--text">{{ f.type }}</td>
                                    <td>{{ f.value }}</td>
                                  </tr>
                                </tbody>
                              </template>
                            </v-simple-table>

                            <div
                              v-if="currentInputOption === 'original'"
                              style="background: #202121; height: 150px; max-width: 900px"
                              class="pa-4 overflow-y-auto"
                            >
                              {{ item.outputDataMapHex }}
                            </div>
                            <div>
                              <v-menu offset-y>
                                <template #activator="{ attrs, on }">
                                  <v-btn class="white--text mt-3" v-bind="attrs" depressed tile small v-on="on">
                                    View Input As <v-icon small class="ml-1">mdi-chevron-down</v-icon>
                                  </v-btn>
                                </template>
                                <v-list dense>
                                  <v-list-item-group v-model="currentInputOption" color="primary">
                                    <v-list-item
                                      v-for="item in inputOptions"
                                      :key="item.value"
                                      :value="item.value"
                                      link
                                      dense
                                    >
                                      <v-list-item-title>{{ item.text }}</v-list-item-title>
                                    </v-list-item>
                                  </v-list-item-group>
                                </v-list>
                              </v-menu>
                            </div>
                          </div>
                        </td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-card>
            </v-card>
          </v-col>
        </v-row>
        <v-row v-else>
          <div class="text-h4 ml-3 pt-6">{{ error }}</div>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, Ref, ref, useRoute } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { EvmTransactionGQL } from '~/apollo/queries'
import { EvmTransaction } from '~/types/graph'

interface FunctionParams {
  name: string
  type: string
  indexed: boolean
  value: any
}

export default defineComponent({
  components: {},
  setup() {
    const route = useRoute()
    const hash = computed(() => route.value.params?.id ?? '')
    const chainId = computed(() => route.value.query.chainId ?? '1') as Ref<string>

    const { result, error } = process.browser
      ? useQuery(EvmTransactionGQL, () => ({
          chainId: parseFloat(chainId.value),
          hash: hash.value,
        }))
      : { result: ref(null), error: ref('') }

    const txData = computed<EvmTransaction>(() => result.value?.evmTransaction ?? null)
    const txDataFormatted = computed(() =>
      txData.value
        ? {
            ...txData.value,
            dateTime: new Date(txData.value.timestamp * 1000),
            txFeeInETH: `${txData.value.transactionFee * Math.pow(10, -18)} ETH`,
            gasPriceInGwei: `${txData.value.gasPrice * Math.pow(10, -9)} Gwei`,
            gasPriceInETH: `${(txData.value.gasPrice * Math.pow(10, -18)).toFixed(18)} ETH`,
            dataToArrayOfStrings: toArrayOfString(txData.value.input?.inputsSigDataStr ?? '', 64),
          }
        : null
    )
    const inputDetails = computed<FunctionParams[]>(() => {
      const d = sortObjectByKey(txDataFormatted.value?.input?.argsMap ?? {})
      d.forEach((elem) => {
        elem.value = txDataFormatted.value?.input?.inputsMap
          ? Object.prototype.hasOwnProperty.call(txDataFormatted.value?.input?.inputsMap ?? {}, elem.name)
            ? txDataFormatted.value?.input?.inputsMap[elem.name]
            : ''
          : ''
      })
      return d
    })
    const logEvents = computed(() =>
      txDataFormatted.value?.logEvents?.items.map((elem) => ({
        ...elem,
        functionParams: getAllLogsData(elem?.allFunctionParams ?? {}),
        indexedVals: getAllLogsData(elem?.allFunctionParams ?? {})
          .filter((elem) => elem.indexed)
          .map((elem) => elem.value),
      }))
    )

    const currentInputOption = ref('default')
    const inputOptions = ref([
      { text: 'Default', value: 'default' },
      { text: 'Original', value: 'original' },
    ])
    const defaultView = ref('overview')
    const viewOptions = ref([
      { text: 'Overview', value: 'overview' },
      { text: 'Input Details', value: 'input' },
      { text: 'Log Events', value: 'log' },
    ])

    const toArrayOfString = (inputString: string, partLength: number): string[] => {
      const arrayOfStrings = []
      let currentIndex = 0

      while (currentIndex < inputString.length) {
        arrayOfStrings.push(inputString.substr(currentIndex, partLength))
        currentIndex += partLength
      }

      return arrayOfStrings
    }

    function sortObjectByKey(obj: { [key: string]: { [key: string]: string } }): FunctionParams[] {
      const res: FunctionParams[] = []
      const sortedKeys = Object.keys(obj).sort()
      for (const key of sortedKeys) {
        const val = obj[key]
        res.push({ name: Object.keys(val)[0], type: Object.values(val)[0], value: '', indexed: false })
      }
      return res
    }

    function getAllLogsData(
      obj: Record<string, { indexed: boolean; name: string; type: string; value: any }>
    ): FunctionParams[] {
      const res: FunctionParams[] = []
      const sortedKeys = Object.keys(obj).sort()

      for (const key of sortedKeys) {
        const val = obj[key]
        res.push({ name: val.name, type: val.type, value: val.value, indexed: val.indexed })
      }
      return res
    }

    return {
      tx: txDataFormatted,
      currentInputOption,
      inputOptions,
      inputDetails,
      logEvents,
      defaultView,
      viewOptions,
      error,
    }
  },
  head: {},
})
</script>

<style scoped></style>
