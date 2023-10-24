// This component will create a dynamic form that has a ParamInputForm for each parameter
//  

<template>
    <div>
      <h3  style="font-size:32px; color:green; font-weight: bold;">Enter Function Parameters</h3>
      <p style="font-size:20px; color:green; font-weight: bold;">Function Selected: </p>



      <v-btn color="red" outlined @click="previewTransaction">
        preview transaction
      </v-btn>
    </div>

</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  computed
} from '@nuxtjs/composition-api'

//import useDynamicAbi from '~/composables/useDynamicAbi'
import TransactionResult from '~/components/common/TransactionResult.vue'

export default defineComponent({
  setup(_,{emit}){

    const previewTransaction = () => {
      console.log('previewTransaction')
      emit('previewTransaction')
    }

    const cols = computed(() => {
      return [
        {
          text: 'Parameter Name',
          align: 'left',
          value: 'name',
          width: '50',
        },
        {
          text: 'Type',
          align: 'left',
          value: 'type',
          width: '50',
        },
        {
          text: '',
          align: 'left',
          value: 'input',
          width: '100',
        },
      ]
    })

    return {
        cols,
        previewTransaction
    }
  }
})

</script>

<style>

</style>