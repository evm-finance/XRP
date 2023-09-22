<template>
    <div>
        <v-row>
            <xrpAccountHistory></xrpAccountHistory>
        </v-row>
        <v-row>
            <xrpBalances></xrpBalances>
        </v-row>
    </div>
</template>
<script lang="ts">
import { defineComponent, ref, computed, watch, useContext, Ref, reactive, inject } from '@nuxtjs/composition-api'
import xrpAccountHistory from '~/components/xrp/xrpAccountHistory.vue'
import xrpBalances from '~/components/xrp/xrpBalances.vue'


export default defineComponent({
  components: { xrpBalances,  xrpAccountHistory},
})

</script>
