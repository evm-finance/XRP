/*
import { provide, onGlobalSetup, defineNuxtPlugin } from '@nuxtjs/composition-api'
import { DefaultApolloClient } from '@vue/apollo-composable/dist'
// import { onError } from '@apollo/client/link/error'
export default defineNuxtPlugin(({ app }) =>
  onGlobalSetup(() => {
    // console.log(DefaultApolloClient, app.apolloProvider.defaultClient)

    provide(DefaultApolloClient, app.apolloProvider?.defaultClient)

    /!*   const link = onError(({ graphQLErrors, networkError }) => {
      if (graphQLErrors)
        graphQLErrors.map(({ message, locations, path }) =>
          console.log(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`)
        )

      if (networkError) console.log(`[Network error]: ${networkError}`)
    })

    console.log(link) *!/
  })
)
*/

import { Context } from '@nuxt/types'
import { provide, onGlobalSetup, defineNuxtPlugin } from '@nuxtjs/composition-api'
import { DefaultApolloClient } from '@vue/apollo-composable'

/**
 * This plugin will connect @nuxt/apollojs with @vue/apollo-composable
 */

export default defineNuxtPlugin(({ app }: Context): void => {
  onGlobalSetup(() => {
    provide(DefaultApolloClient, app.apolloProvider?.defaultClient)
  })
})
