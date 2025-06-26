<template>
  <xrpAccountHistory />
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import xrpAccountHistory from '~/components/xrp/xrpAccountHistory.vue'

export default defineComponent({
  components: { xrpAccountHistory },
})
</script> 