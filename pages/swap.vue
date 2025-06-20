<template>
  <v-row no-gutters justify="center">
    <v-col cols="12" md="6" class="pa-2">
      <evm-swap />
    </v-col>
    <v-col cols="12" md="6" class="pa-2">
      <xrp-amm-swap />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import EvmSwap from '~/components/trading/EvmSwap.vue'
import XrpAmmSwap from '~/components/trading/XrpAmmSwap.vue'

export default defineComponent({
  components: { EvmSwap, XrpAmmSwap },
})
</script>
