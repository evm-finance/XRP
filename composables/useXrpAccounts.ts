import { ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import { useQuery, useSubscription } from '@vue/apollo-composable/dist'
import { Block } from '@/types/apollo/main/types'
import { XRPDefiDataGQL } from '~/apollo/queries'

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

  const { defiData } = useQuery(XRPDefiDataGQL, () => ({ network: 'ripple' }), {
    fetchPolicy: 'no-cache',
  })

  // EVENTS



  return { defiData }
}
