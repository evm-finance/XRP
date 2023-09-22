import { ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import { useQuery, useSubscription } from '@vue/apollo-composable/dist'
import { XRPDefiDataGQL } from '~/apollo/queries'


export default function useXrpAccounts() {
  // COMPOSABLES
  const { $emitter } = useContext()

  const { result } = useQuery(XRPDefiDataGQL, () => ({ account: 'rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj' }), {
    fetchPolicy: 'no-cache',
  })

  
  // EVENTS



  return { result }
}
