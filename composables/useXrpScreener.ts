import * as process from 'process'
import { ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'
import { Block } from '@/types/apollo/main/types'
import { BlocksXrpGQL } from '~/apollo/queries'

type BlockObserver = Block & {
  updateOption?: { status: boolean; color: string | null }
}

export function useXrpScreener() {
  console.log('🔍 [DEBUG] useXrpScreener() called')
  
  const { $f } = useContext()
  const { useLoggedQuery, useLoggedSubscription } = useXrpGraphQLWithLogging()

  // COMPOSABLES
  const { $emitter } = useContext()

  // STATE
  const loading = ref<boolean>(true)
  const pageNumber = ref<number>(0)
  let updateTimeout: any = null
  const blocks = ref<Block[]>([])
  const currentTime = ref<number>(new Date().getTime() / 1000)

  console.log('🔍 [DEBUG] useXrpScreener: About to execute GraphQL query')

  // Log the query content BEFORE making the call
  console.log('🚀 [BEFORE QUERY] useXrpScreener - BlocksXrpGQL:', {
    query: BlocksXrpGQL.loc?.source.body,
    variables: { network: 'ripple' },
    timestamp: new Date().toISOString()
  })

  // GraphQL query for XRP blocks with enhanced logging
  const { onResult } = useLoggedQuery(BlocksXrpGQL, () => ({ network: 'ripple' }), {
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
    context: {
      queryName: 'BlocksXrp',
      component: 'useXrpScreener',
      purpose: 'XRP blockchain data for screener'
    }
  })

  console.log('🔍 [DEBUG] useXrpScreener: GraphQL query executed, setting up result handler')

  const { result: liveBlock } = process.browser
    ? useLoggedSubscription(BlocksXrpGQL, () => ({ network: 'ripple' }), {
        fetchPolicy: 'no-cache',
        context: {
          queryName: 'BlocksXrp',
          component: 'useXrpScreener',
          purpose: 'XRP blockchain live data subscription'
        }
      })
    : { result: ref(null) }

  const nextPage = () => pageNumber.value++

  const currentPage = computed({
    get: () => pageNumber.value + 1,
    set: (v: number) => (pageNumber.value = v - 1),
  })

  // EVENTS
  onResult((queryResult: any) => {
    console.log('🔍 [DEBUG] useXrpScreener onResult called:', {
      hasData: !!queryResult.data,
      blocksCount: queryResult.data?.blocks?.length || 0,
      loading: queryResult.loading
    })
    
    blocks.value = queryResult.data?.blocks ?? []
    loading.value = queryResult.loading
    currentTime.value = new Date().getTime() / 1000
  })

  watch(liveBlock, (val: any) => {
    console.log('🔍 [DEBUG] useXrpScreener liveBlock watch triggered:', {
      hasData: !!val,
      blockCount: val?.block?.length || 0
    })
    
    const newData: BlockObserver[] | Block[] = val?.block ?? []
    addNewRecords(newData)
    $emitter.emit('onNewBlock', newData)
  })

  function addNewRecords(newRecords: BlockObserver[]) {
    console.log('🔍 [DEBUG] useXrpScreener addNewRecords called:', newRecords.length)
    
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

  console.log('🔍 [DEBUG] useXrpScreener: Returning composable result')
  return { blocks, currentPage, loading, currentTime, nextPage, testUpdate: addNewRecords }
}

// Add default export for compatibility
export default useXrpScreener
