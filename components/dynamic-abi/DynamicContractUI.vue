// This component will create a dynamic form that has a ParamInputForm for each parameter
//  
//0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX

<template>
    <div class="dynamic_ui">
      <h3  style="font-size:32px; color:white; font-weight: bold;">Enter Function Parameters</h3>
      <p style="font-size:20px; color:white; font-weight: bold;">Function Selected: {{  }}</p>
      <v-col v-for >

        
      </v-col>


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
        previewTransaction,
    }
  }
})

</script>

<style>
.dynamic_ui {
    background-color: rgb(10, 10, 10);
    word-wrap: break-word;
    font-size: 14px;
    width: 600px;
    height: 800px;
}

</style>