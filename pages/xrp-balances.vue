<template>
  <xrpBalances />
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import xrpBalances from '~/components/xrp/xrpBalances.vue'

export default defineComponent({
  components: { xrpBalances },
})
</script> 