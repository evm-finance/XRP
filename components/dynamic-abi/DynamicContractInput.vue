// This component will create a dynamic form that has a ParamInputForm for each parameter
//  
//0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX

<template>
    <div class="dynamic_ui">
      <h3  style="font-size:32px; color:white; font-weight: bold;margin-bottom: 10px;">Enter Function Parameters</h3>
      <p style="font-size:20px; color:white; font-weight: bold;margin-bottom: 10px;">Function Selected: {{ functionName }}</p>
      <v-col>
          <v-row v-if="payable"> 
          <v-col style="font-size:20px; color:white; font-weight: bold;margin-bottom: 20px;"> Value
          </v-col>
          <v-text-field
            v-model="callValueInput" style="width:300px; background-color: rgba(241, 236, 236, 0.9); font-color:black">
          </v-text-field>
          </v-row>
          <v-row v-for="(item,index) in calldataParams" :key="index">
            <v-col style="font-size:20px; color:white; font-weight: bold; margin-bottom: 20px;"> Name:{{ item.name }} </v-col>
            <v-col style="font-size:20px; color:white; font-weight: bold; margin-bottom: 20px;"> Type:{{ item.type }} </v-col>
            <v-text-field
              v-model="calldataObject[index]" style="width:300px; background-color: rgba(241, 236, 236, 0.9); font-color:black" type="text" > 
            </v-text-field>
          </v-row>
          <v-btn style="margin-top:50px;" color="red" outlined @click="previewTransaction">
            preview transaction
          </v-btn>
        </v-col> 
      




    </div>

</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  computed, PropType, inject, Ref, reactive
} from '@nuxtjs/composition-api'

//import useDynamicAbi from '~/composables/useDynamicAbi'
import TransactionResult from '~/components/common/TransactionResult.vue'
import { calldataTemplate, CalldataAbi, EventElem } from '~/types/apollo/main/types'
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'
import { BigNumber, ethers } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'




export default defineComponent({
  props:{
    calldataParams:{type: Array as PropType<CalldataAbi[] | EventElem[]>, required: true},
    functionName:{type: String as PropType <string>, required: true},
    payable:{type:Boolean, required:true},
    address:{type: String, required:true},
    contract:{type:ethers.Contract, required:true}
  },
  setup(props,{emit}){


    const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
    const { allowedSpending, approveMaxSpending } = useERC20()
    const callValueInput = ref(0)
    const calldataObject : Ref<(string | number)[]> = ref([])

    const ERC20_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`125`).div('100').toNumber()
  const NATIVE_ETH_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`175`).div('100').toNumber()
  const ESTIMATED_GAS_FEE = async (tx: any): Promise<BigNumber> =>
    (await provider.value?.estimateGas(<any>tx)) ?? BigNumber.from('0')



      //0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX

    const previewTransaction = async () => {
      console.log(props.calldataParams)
      console.log('previewTransaction')
      console.log(props.contract)
      emit('previewTransaction')
      var finalCalldata: (string | number)[] = []
      for (let i = 0; i < calldataObject.value.length; i++) 
      {
        console.log('lost')
        if (props.calldataParams[i].type[0] == 'u' )
        {
          console.log('found', calldataObject.value[i])
          finalCalldata[i] = Number(calldataObject.value[i])
        }
        else
        {
          console.log('not found')
          finalCalldata[i] = calldataObject.value[i]
        }
      }

      try{
        console.log('testing tx')
        console.log(props.functionName)
        console.log(...finalCalldata)
              //0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX

        var testTx = await props.contract.populateTransaction[props.functionName](...finalCalldata);
        console.log(testTx)
        const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(testTx)
        const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
        console.log('gas limit',gasLimit)
        console.log(props.contract.functions[props.functionName])
        const depositCall = await props.contract.functions[props.functionName]({value: callValueInput.value, gasLimit, data:testTx.data})
        console.log(depositCall)
        const resp = await depositCall.wait()
      }

      catch(error)
      {
        console.log(error)
        console.log('failed to generate transaction, hurry up, its getting late')
      }


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
        functionName: props.functionName,
        payabe: props.functionName,
        callValueInput,
        calldataObject,


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