<template>
  <div>
    <h1>Test Page</h1>
    <p>If you can see this, Nuxt is working!</p>
    <p>GraphQL Server URL: {{ graphqlUrl }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    return {
      graphqlUrl: process.env.BASE_GRAPHQL_SERVER_URL || 'Not set'
    }
  }
})
</script> 