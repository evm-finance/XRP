import * as process from 'process'
import { plainToClass } from 'class-transformer'
import { useQuery } from '@vue/apollo-composable/dist'
import { computed, inject, reactive, Ref, ref, useRoute, watch } from '@nuxtjs/composition-api'
import { TransactionsGQL } from '~/apollo/main/portfolio.query.graphql'

import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { LogEvent, TransactionItem, TxDetail } from '~/types/apollo/main/types'

export class TransactionModel implements TransactionItem {
  readonly blockSignedAt!: string
  readonly fromAddress!: string
  readonly fromAddressIsContract!: boolean
  readonly fromAddressName!: string
  readonly fromAddressSymbol!: string
  readonly gasPrice!: number
  readonly gasQuote!: number
  readonly gasSpent!: number
  readonly txDetails!: TxDetail[]
  readonly logEvents!: LogEvent[]
  readonly successful!: boolean
  readonly toAddress!: string
  readonly toAddressIsContract!: boolean
  readonly toAddressName!: string
  readonly toAddressSymbol!: string
  readonly txHash!: string
  readonly value!: string
  readonly valueQuote!: number
  readonly callFunction!: string
  readonly decodedFunctionID!: string

  get txDate(): string {
    return new Date(this.blockSignedAt).toLocaleString()
  }

  get txDetail(): TxDetail | undefined {
    if (this.methodTextRenderer === 'Transfer' && this.txDetails) {
      return this.txDetails[0]
    }
  }

  get methodTextRenderer(): string {
    const split: string[] = this.callFunction.split('(')
    if (split.length > 1) {
      const callName = split[0][0].toUpperCase() + split[0].slice(1)
      if (callName === 'Multicall') {
        return 'Swap'
      }
      return callName
    }

    if (this.logEvents.length === 0) {
      return 'Transfer'
    }

    if (this.logEvents.length > 2) {
      return 'Swap'
    }

    return 'Unknown'
  }

  get txnValue(): { amount: number; priceUsd: number } {
    if (this.txDetail) {
      const formattedValue = +this.txDetail.value / 10 ** this.txDetail.tokenContractDecimals
      return {
        priceUsd: this.txDetail.priceUSD * formattedValue,
        amount: formattedValue,
      }
    }

    return {
      priceUsd: this.valueQuote,
      amount: +this.value / 10 ** 18,
    }
  }

  get txnFee(): number {
    return (this.gasPrice / 10 ** 18) * this.gasSpent
  }

  get isSuccess(): { text: string; color: string } {
    return this.successful ? { text: 'Success', color: 'green' } : { text: 'Failed', color: 'pink' }
  }
}
type AnyFunction = (...args: any[]) => void
export default function (chainId: Ref<number>) {
  // STATE
  const route = useRoute()
  const loading = ref(true)
  const currentPage = ref(1)
  const hasMore = ref(false)
  const pagination = reactive({ page: 0, total: 1, perPage: 15, visible: 30 })

  // COMPOSABLES
  const { account: metamaskWalletAddress, walletReady } = inject(WEB3_PLUGIN_KEY) as Web3

  const customWalletAddress = computed(() => route.value.query?.wallet ?? '') as Ref<string>

  const isWalletReady = computed(() => walletReady.value || customWalletAddress.value?.length)

  const account = computed(() =>
    customWalletAddress.value?.trim().length ? customWalletAddress.value?.trim() : metamaskWalletAddress.value
  )

  const myFunction: AnyFunction = (_: any) => {}
  const { result, error, onResult } = process.browser
    ? useQuery(
        TransactionsGQL,
        () => ({
          chainId: chainId.value ?? 1,
          address: account.value ?? '',
          pageNumber: pagination.page,
          pageSize: pagination.perPage,
        }),
        { fetchPolicy: 'no-cache' }
      )
    : { result: ref(null), error: ref(null), onResult: myFunction }

  // COMPUTED
  const transactionsData = computed(
    () => plainToClass(TransactionModel, result.value?.transactions?.items as TransactionModel[]) ?? []
  ) as Ref<TransactionModel[]>

  // EVENTS
  onResult((queryResult) => {
    hasMore.value = queryResult.data.transactions.pagination.hasMore
    loading.value = queryResult.loading

    if (hasMore.value && pagination.total <= currentPage.value) {
      pagination.total++
    }
  })

  // On chain network or account change, reset currentPage and pagination.total
  watch([chainId, account], () => {
    loading.value = true
    currentPage.value = 1
    pagination.total = 1
  })

  // Offset page index
  watch(currentPage, (newVal) => {
    loading.value = true
    pagination.page = newVal - 1
  })

  // METHODS
  function isInboundRenderer(item: TransactionModel): { color: string; text: string } {
    if (item.logEvents && item.logEvents.length > 0) {
      if (item.methodTextRenderer === 'Transfer') {
        const param = Object.values(item.logEvents[0].decoded.params).filter((el: any) => el.name === 'to')[0]
        if (param && param.value === account.value?.toLowerCase()) {
          return { color: 'green', text: 'IN' }
        }
      }
    }

    if (item.toAddress.toLowerCase() === item.fromAddress.toLowerCase()) {
      return { color: 'grey', text: 'SELF' }
    }

    if (item.toAddress.toLowerCase() === account.value?.toLowerCase()) {
      return { color: 'green', text: 'IN' }
    }

    return { color: 'pink', text: 'OUT' }
  }

  return {
    loading,
    error,
    walletReady,
    account,
    transactionsData,
    hasMore,
    currentPage,
    pagination,
    isWalletReady,

    isInboundRenderer,
  }
}
