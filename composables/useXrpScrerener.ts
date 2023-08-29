import * as process from 'process'
import { ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import { useQuery, useSubscription } from '@vue/apollo-composable/dist'
import { Block } from '@/types/apollo/main/types'
import { BlocksSubscriptionGQL, BlocksXrpGQL } from '~/apollo/queries'

type BlockObserver = Block & {
  updateOption?: { status: boolean; color: string | null }
}

export default function () {
  // COMPOSABLES
  const { $emitter } = useContext()

  // STATE
  const loading = ref<boolean>(true)
  const pageNumber = ref<number>(0)
  let updateTimeout: any = null
  const blocks = ref<Block[]>([])
  const currentTime = ref<number>(new Date().getTime() / 1000)

  const { onResult } = useQuery(BlocksXrpGQL, () => ({ network: 'ripple' }), {
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
  })

  const { result: liveBlock } = process.browser
    ? useSubscription(BlocksSubscriptionGQL, () => ({ network: 'ripple' }), {
        fetchPolicy: 'no-cache',
      })
    : { result: ref(null) }

  const nextPage = () => pageNumber.value++

  const currentPage = computed({
    get: () => pageNumber.value + 1,
    set: (v: number) => (pageNumber.value = v - 1),
  })

  // EVENTS
  onResult((queryResult: any) => {
    blocks.value = queryResult.data?.blocks ?? []
    loading.value = queryResult.loading
    currentTime.value = new Date().getTime() / 1000
  })

  watch(liveBlock, (val: any) => {
    const newData: BlockObserver[] | Block[] = val?.block ?? []
    addNewRecords(newData)
    $emitter.emit('onNewBlock', newData)
  })

  function addNewRecords(newRecords: BlockObserver[]) {
    clearTimeout(updateTimeout)
    newRecords = newRecords.map((elem) => ({
      ...elem,
      updateOption: { status: true, color: '#4caf5026' },
    }))
    blocks.value = [...newRecords, ...blocks.value]
    currentTime.value = new Date().getTime() / 1000
    if (blocks.value.length > 25) {
      blocks.value.splice(-newRecords.length)
    }

    updateTimeout = setTimeout(() => {
      blocks.value.forEach((elem: BlockObserver) => {
        if (elem.updateOption) {
          elem.updateOption = { status: false, color: null }
        }
      })
    }, 1000)
  }

  return { blocks, currentPage, loading, currentTime, nextPage, testUpdate: addNewRecords }
}
