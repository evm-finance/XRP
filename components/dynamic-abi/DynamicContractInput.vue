// This component will create a dynamic form that has a ParamInputForm for each parameter
//  
//0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX

<template>
    <div class="dynamic_ui">
      <h3  style="font-size:32px; color:white; font-weight: bold;">Enter Function Parameters</h3>
      <p style="font-size:20px; color:white; font-weight: bold;">Function Selected: {{  }}</p>
      
        <v-row v-for="f in calldataParams" :key="f">
          <v-col> Name:{{ f.name }}</v-col>
          <v-col> Type:{{ f.type }}</v-col>
          <v-text-field
           style="width:300px; background-color: rgba(241, 236, 236, 0.9); font-color:black"></v-text-field>
        </v-row>
      <v-btn color="red" outlined @click="previewTransaction">
        preview transaction
      </v-btn>
    </div>

</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  computed, PropType
} from '@nuxtjs/composition-api'

//import useDynamicAbi from '~/composables/useDynamicAbi'
import TransactionResult from '~/components/common/TransactionResult.vue'
import { inputAbi, CalldataAbi, EventElem } from '~/types/apollo/main/types'



export default defineComponent({
  props:{
    calldataParams:{type: Array as PropType<inputAbi[]>, required: true}
  },
  setup(props,{emit}){


    const previewTransaction = () => {
      console.log(props.calldataParams)
      console.log('previewTransaction')
      emit('previewTransaction')
    }
    
    //const funcs = [{"name":"amount","type":"uint256"},{"name":"to","type":"address"}]
    // const calldata = computed(() => {
    //   return []
    // })
    const cols = computed(() => {
      return [
        {
          text: 'Parameter Name',
          align: 'left',
          value: 'name',
          width: '30px',
        },
        {
          text: 'Type',
          align: 'left',
          value: 'type',
          width: '30px',
        },
        {
          text: '',
          align: 'left',
          value: 'input',
          width: '100px',
        },
      ]
    })

    return {
        cols,
        calldataParams:props.calldataParams,

        previewTransaction,
    }
  }
})

</script>

<style>
.dynamic_ui {
    background-color: rgba(3, 3, 3, 0.9);
    word-wrap: break-word;
    font-size: 14px;
    width: 600px;
    height: 600px;
}

</style>