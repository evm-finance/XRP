<template>
    <div justify="center">
        <div class="padding">
            <v-row>
                <contract-address-form @addressEntered="setAddress"></contract-address-form>
            </v-row>
            <v-row>
                <abi-input-form @abiEntered="generateUi" @reload="loadAgain"></abi-input-form>
            </v-row>
            <v-row v-if="renderButtons">
                <function-buttons :functionNames="functions"
                @functionSelected="openWindow"></function-buttons>
            </v-row>
            <v-dialog v-model="dialog" width="600">
                <dynamic-contract-input :calldataParams=selectedAbi :functionName=selectedFunction @previewTransaction="sendTransaction"></dynamic-contract-input>
            </v-dialog>
        </div>
    </div>
</template>

<script lang="ts">
import { onClickOutside } from '@vueuse/core'
import { defineComponent, ref, inject, reactive, Ref, del } from '@nuxtjs/composition-api';
import AbiInputForm from '~/components/dynamic-abi/AbiInputForm.vue'
import ContractAddressForm from '~/components/dynamic-abi/ContractAddressForm.vue'
import FunctionButtons from '~/components/dynamic-abi/FunctionButtons.vue'
import DynamicContractInput from '~/components/dynamic-abi/DynamicContractInput.vue';
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'
import { BigNumber, ethers } from 'ethers'
// import { State } from '~/types/state'
// import { DefiEvents, EmitEvents } from '~/types/events'
import type { ContractTransaction } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Chain, AbiElem, EventElem, AbiEvent, CalldataAbi, inputAbi } from '~/types/apollo/main/types'
import { ConstructorFragment } from 'ethers/lib/utils';

export default defineComponent({
    components:{
    AbiInputForm,
    ContractAddressForm,
    FunctionButtons,
    DynamicContractInput
},
setup() {

    const ESTIMATED_GAS_FEE = async (tx: any): Promise<BigNumber> =>
        (await provider.value?.estimateGas(<any>tx)) ?? BigNumber.from('0')
    const ERC20_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`125`).div('100').toNumber()
    const NATIVE_ETH_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`175`).div('100').toNumber()
    const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
    type calldataElement = string | BigNumber
    type abiTemplate = {
        name: string,
        type: string,
        value: calldataElement | null
    }
    const { allowedSpending, approveMaxSpending } = useERC20()
    const renderButtons = ref(false)
    const dialog = ref(false)
    const address = ref('')
    //const functions : Ref<abiTemplate[]> = ref([])
    type rawAbi = AbiElem | AbiEvent
    type inputAbi = CalldataAbi | EventElem
    const functions : Ref<rawAbi[]> = ref([])
    const calldata : Ref<calldataElement[]> = ref([])
    const selectedType = ref('')
    const selectedFunction = ref('')
    const selectedAbi : Ref<inputAbi[]> = ref([])

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve,ms));
    
    const loadAgain = async () => {
        console.log("reloading")
        if(renderButtons.value == true)
        {
            renderButtons.value = false
            await delay(100)
        }
        renderButtons.value = true
    }

    const setAbi = () => {
        for(let i = 0; i < functions.value.length; i++)
            {
                // console.log('lost')
                // console.log(functions.value[i].type,functions.value[i].name)
                // console.log(selectedType.value,selectedFunction.value)
                if(functions.value[i].type == selectedType.value && functions.value[i].name == selectedFunction.value)
                {
                    console.log('found')
                    //selectedAbi.value = functions.value[i].inputs
                    for(let j = 0; j < functions.value[i].inputs.length; j++)
                    {
                        selectedAbi.value[j] = (functions.value[i].inputs[j])
                    }
                }
                //console.log(data[i])
                //functions.value.push(jsonData[i])
            }
    }

    const setAddress = (addressInput: string) => {
        //console.log(addressInput)
        address.value = addressInput
        //console.log(addressInput)
    }

    const generateUi = (data: Array<any>) => {

        console.log('starting generateUI')
        console.log(data)
        console.log(address.value)

        
        for(let i = 0; i < data.length; i++)
            {
                //console.log(data[i])
                //functions.value.push(jsonData[i])
            }
            functions.value = data
        // console.log('render buttons',renderButtons.value)
        // console.log('finised generateUI')
    }

    const sendTransaction = async () => {
        console.log(selectedFunction.value,signer.value,account.value,selectedType.value,chainId.value)
        try{
            const contract = new ethers.Contract(
                address.value,
                functions.value,
                signer.value
            )
            //setAbi()
            console.log(selectedAbi.value)
            console.log(contract)
            // var testTx = await contract.populateTransaction[selectedFunction.value]();
            // console.log(testTx)
            // const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(testTx)
            // const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
            // console.log(gasLimit)
            // console.log(contract.functions[selectedFunction.value])
            // const depositCall = await contract.functions[selectedFunction.value]({value:calldata.value,gasLimit})
            // console.log(depositCall)
            // const resp = await depositCall.wait()


            //contract.
        }
        catch (error)
        {
            console.log('failed to instantiate contract, try harder')
        }
        //contractInstance.{{function}}
    }
    //0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc
    
    const openWindow = (name: string,type: string) => {
        //
        if(dialog.value)
        {
            dialog.value = false
        }
        console.log("window opened", 'name:',name,'type:',type)
        selectedFunction.value = name
        selectedType.value = type

        setAbi()
        delay(100)

        console.log("finishing")
        
        console.log()

        dialog.value=true
    }

    return {
        // functions
        generateUi,
        openWindow,
        loadAgain,
        setAddress,
        sendTransaction,
        
        //data
        renderButtons,
        functions,
        dialog,
        selectedFunction,
        selectedAbi

    }
}
})

</script>

<style scoped>
.padding {
    margin-top: 30px;
    padding-left: 100px;
}
</style>

