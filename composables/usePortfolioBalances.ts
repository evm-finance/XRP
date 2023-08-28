import * as process from 'process'
import { useQuery } from '@vue/apollo-composable/dist'
import { computed, inject, Ref, ref, useRoute, useStore, watch } from '@nuxtjs/composition-api'
import { BalancesGQL } from '~/apollo/main/portfolio.query.graphql'
import { BalanceItem, Balance } from '~/types/apollo/main/types'
import { State } from '~/types/state'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'

type AnyFunction = (...args: any[]) => void

export default function () {
  // STATE
  const loading = ref(true)

  // COMPOSABLES
  const { state } = useStore<State>()
  const { account: metamaskWalletAddress } = inject(WEB3_PLUGIN_KEY) as Web3
  const route = useRoute()

  const customWalletAddress = computed(() => route.value.query?.wallet ?? '') as Ref<string>

  const account = computed(() =>
    customWalletAddress.value?.trim().length ? customWalletAddress.value?.trim() : metamaskWalletAddress.value
  )
  const isWalletReady = computed(() => account.value.length)

  const myFunction: AnyFunction = (_: any) => {}
  const { result, error, onResult } = process.browser
    ? useQuery(BalancesGQL, () => ({ chainIds: state.configs.balancesChains, address: account.value }), {
        fetchPolicy: 'no-cache',
      })
    : { result: ref(null), error: ref(null), onResult: myFunction }

  // COMPUtED
  const balanceData = computed(() => result.value?.balances ?? []) as Ref<Balance[]>
  const totalBalance = computed(() =>
    balanceData.value
      .map((elem: Balance) => elem.items.reduce((n: number, el: BalanceItem) => n + el.quote, 0))
      .reduce((n: number, curr) => n + curr, 0)
  )

  // EVENTS
  onResult((queryResult) => (loading.value = queryResult.loading))

  // WATCHERS
  watch(account, () => (loading.value = true))

  return {
    balanceData,
    totalBalance,
    isWalletReady,
    error,
    loading,
    account,
  }
}
