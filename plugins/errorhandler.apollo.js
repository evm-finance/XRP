// export default (graphqlError, { store, error, redirect, route }) => {
//   console.log("EEEEEEEEEEEEEEE", graphqlError)
//
//   const { networkError, message, gqlError, graphqlErrors } = graphqlError
//
//   // handle error
//
//   return error({ statusCode: 503, message })
// }
import { onError } from '@apollo/client/link/error'

export default ({ graphQLErrors, networkError, operation, forward }, nuxtContext) => {
  console.log('Global error handler')
  console.log(graphQLErrors, networkError, operation, forward)
  console.log(nuxtContext)
}
