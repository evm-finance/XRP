<template>
  <v-row no-gutters justify="center"> <evm-swap /> </v-row>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import EvmSwap from '~/components/trading/EvmSwap.vue'

export default defineComponent({
  components: { EvmSwap },
})
</script>
