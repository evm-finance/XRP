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
            <v-dialog v-model="dialog" v-if="dialog" width="600" @click:outside="handleClickOutside">
                <dynamic-contract-input :calldataParams=selectedAbi :payable=payableFunction :functionName=selectedFunction :address=address :contract=contractInstance @previewTransaction="sendTransaction"></dynamic-contract-input>
            </v-dialog>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, inject, reactive, Ref, del, nextTick } from '@nuxtjs/composition-api';
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
import { Chain, AbiElem, EventElem, AbiEvent, CalldataAbi, inputAbi, calldataTemplate } from '~/types/apollo/main/types'
import { ConstructorFragment } from 'ethers/lib/utils';

export default defineComponent({
    components:{
    AbiInputForm,
    ContractAddressForm,
    FunctionButtons,
    DynamicContractInput
},
setup() {

    const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
    const { allowedSpending, approveMaxSpending } = useERC20()
    const renderButtons = ref(false)
    const dialog = ref(false)
    const address = ref('')
    const callValue = ref(0)
    const contractInstance = ref<ethers.Contract | null>(null); 
    //const functions : Ref<abiTemplate[]> = ref([])
    type rawAbi = AbiElem | AbiEvent
    //type inputAbi = CalldataAbi | EventElem
    const functions : Ref<rawAbi[]> = ref([])
    const calldata : Ref<calldataTemplate[]> = ref([])
    const selectedType : Ref<string> = ref('')
    const selectedFunction : Ref<string> = ref('')
    const payableFunction = ref(false)
    const selectedAbi : Ref<CalldataAbi[] | EventElem[]> = ref([])

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve,ms));

    const handleClickOutside = () => 
    {
        console.log('clicked outside, closing')
        dialog.value=false
        selectedFunction.value = ''
        selectedAbi.value=[]
        console.log(selectedAbi.value,selectedFunction.value)

    }
    
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
                    
                    var temp = functions.value[i] as AbiElem
                    console.log('found',temp)
                    if('payable' in temp)
                    {
                        
                        payableFunction.value = temp.payable

                        // if(payableFunction.value)
                        //     console.log('payable')
                        // else
                        //     console.log('was not payable function')
                    }
                        

                    
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

        contractInstance.value = new ethers.Contract(
                address.value,
                data,
                signer.value
            )
        
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
            // const depositCall = await contract.functions[selectedFunction.value]({value:callValue.value,gasLimit,})
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
    
    const openWindow = async (name: string,type: string) => {
        //
        console.log(dialog.value)
        console.log("window opened", 'name:',name,'type:',type)
        selectedFunction.value = name
        selectedType.value = type
        console.log(selectedType.value,selectedFunction.value,dialog.value,selectedAbi.value)
        setAbi()
        
        await delay(100)
        dialog.value=true

    }

    return {
        // functions
        generateUi,
        openWindow,
        loadAgain,
        setAddress,
        sendTransaction,
        handleClickOutside,
        
        //data
        renderButtons,
        functions,
        dialog,
        selectedFunction,
        selectedAbi,
        payableFunction,
        callValue,
        address,
        contractInstance

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

