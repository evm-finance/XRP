<template>
  <v-row no-gutters justify="center"> <evm-verse-swap /></v-row>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import EvmVerseSwap from '~/components/trading/EvmVerseSwap.vue'

export default defineComponent({
  components: { EvmVerseSwap },
})
</script>
