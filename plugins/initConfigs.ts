import { Context, Plugin } from '@nuxt/types'
const globalHelper: Plugin = async (context: Context) => await context.store.dispatch('configs/initConfigs', context)

// Validate GraphQL server configuration
if (!process.env.BASE_GRAPHQL_SERVER_URL) {
  console.warn('⚠️  WARNING: BASE_GRAPHQL_SERVER_URL environment variable is not set!')
  console.warn('   Using default local development endpoint: http://127.0.0.1:8080/query')
  console.warn('   For production, please set BASE_GRAPHQL_SERVER_URL in your environment variables.')
} else {
  console.log('✅ GraphQL server URL configured:', process.env.BASE_GRAPHQL_SERVER_URL)
}

if (!process.env.BASE_GRAPHQL_WEBSOCKET_URL) {
  console.warn('⚠️  WARNING: BASE_GRAPHQL_WEBSOCKET_URL environment variable is not set!')
  console.warn('   Using default local development endpoint: ws://127.0.0.1:8080/query')
  console.warn('   For production, please set BASE_GRAPHQL_WEBSOCKET_URL in your environment variables.')
} else {
  console.log('✅ GraphQL WebSocket URL configured:', process.env.BASE_GRAPHQL_WEBSOCKET_URL)
}

export default globalHelper
