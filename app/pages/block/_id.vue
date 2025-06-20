<template>
  <h1>Block</h1>
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
