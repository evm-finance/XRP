<!-- // This component will create a dynamic form that has a ParamInputForm for each parameter
//  

<template>
    <div>

    </div>

</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  PropType,
  ref,
  toRefs,
  UnwrapRef,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import { DefiEvents, EmitEvents } from '~/types/events'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Chain } from '~/types/apollo/main/types'
import useDynamicAbi from '~/composables/useDynamicAbi'
import TransactionResult from '~/components/common/TransactionResult.vue'



</script>

<style>

</style> -->