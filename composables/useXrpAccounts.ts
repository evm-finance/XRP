import { ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import { useSubscription } from '@vue/apollo-composable/dist'
import useXrpGraphQLWithLogging from './useXrpGraphQLWithLogging'
import { XRPDefiDataGQL } from '~/apollo/queries'

export function useXrpAccounts() {
  // COMPOSABLES
  const { $emitter } = useContext()
  const { useLoggedQuery } = useXrpGraphQLWithLogging()

  // Log the query content BEFORE making the call
  console.log('🚀 [BEFORE QUERY] useXrpAccounts - XRPDefiDataGQL:', {
    query: XRPDefiDataGQL.loc?.source.body,
    variables: { account: 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo' },
    timestamp: new Date().toISOString()
  })

  // GraphQL query for XRP DeFi data with enhanced logging
  const { result } = useLoggedQuery(XRPDefiDataGQL, () => ({ account: 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo' }), {
    fetchPolicy: 'no-cache',
    pollInterval: 30000,
    context: {
      queryName: 'XRPDefiData',
      component: 'useXrpAccounts',
      purpose: 'XRP account DeFi data'
    }
  })

  
  // EVENTS



  return { result }
}
