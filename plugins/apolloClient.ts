import { provide, onGlobalSetup, defineNuxtPlugin } from '@nuxtjs/composition-api'
import { DefaultApolloClient } from '@vue/apollo-composable/dist'

export default defineNuxtPlugin(({ app }) =>
  onGlobalSetup(() => {
    // Get the Apollo client
    const apolloClient = app.apolloProvider?.defaultClient
    
    provide(DefaultApolloClient, apolloClient)
  })
)
