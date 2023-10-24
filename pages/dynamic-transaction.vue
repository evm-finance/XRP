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
            <v-dialog v-model="dialog" width="400">
                <dynamic-contract-ui @previewTransaction="sendTransaction"></dynamic-contract-ui>
            </v-dialog>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, inject, reactive, Ref } from '@nuxtjs/composition-api';
import AbiInputForm from '~/components/dynamic-abi/AbiInputForm.vue'
import ContractAddressForm from '~/components/dynamic-abi/ContractAddressForm.vue'
import FunctionButtons from '~/components/dynamic-abi/FunctionButtons.vue'
import DynamicContractUi from '~/components/dynamic-abi/DynamicContractUi.vue';
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'
import { BigNumber, ethers } from 'ethers'
// import { State } from '~/types/state'
// import { DefiEvents, EmitEvents } from '~/types/events'
import type { ContractTransaction } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Chain } from '~/types/apollo/main/types'
import { ConstructorFragment } from 'ethers/lib/utils';

export default defineComponent({
    components:{
    AbiInputForm,
    ContractAddressForm,
    FunctionButtons,
    DynamicContractUi
},
setup() {

    const ESTIMATED_GAS_FEE = async (tx: ContractTransaction): Promise<BigNumber> =>
        (await provider.value?.estimateGas(<any>tx)) ?? BigNumber.from('0')
    const ERC20_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`125`).div('100').toNumber()
    const NATIVE_ETH_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`175`).div('100').toNumber()
    const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
    const { allowedSpending, approveMaxSpending } = useERC20()
    const renderButtons = ref(false)
    const dialog = ref(false)
    const address = ref('')
    const functions = ref([])
    const calldata = ref([])
    const selectedType = ref('')
    const selectedFunction = ref('')

    const delay = (ms) => new Promise(resolve => setTimeout(resolve,ms));
    
    const loadAgain = async () => {
        console.log("reloading")
        if(renderButtons.value == true)
        {
            renderButtons.value = false
            await delay(100)
        }
        renderButtons.value = true
    }

    const setAddress = (addressInput) => {
        //console.log(addressInput)
        address.value = addressInput
        //console.log(addressInput)
    }

    const generateUi = (data) => {

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
        const num = 1
        calldata.value.push(num)
        console.log(num)
        console.log(selectedFunction.value,signer.value,account.value,selectedType.value,chainId.value)
        try{
            const contract = new ethers.Contract(
                address.value,
                functions.value,
                signer.value
            )
            console.log(contract)
            var testTx = await contract.populateTransaction[selectedFunction.value]({value:calldata.value});
            console.log(testTx)
            const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(testTx)
            const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
            console.log(gasLimit)
            console.log(contract.functions[selectedFunction.value])
            const depositCall = await contract.functions[selectedFunction.value]({value:calldata.value,gasLimit})
            console.log(depositCall)
            const resp = await depositCall.wait()


            //contract.
        }
        catch (error)
        {
            console.log('failed to instantiate contract, try harder')
        }
        //contractInstance.{{function}}
    }
    //0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc
    
    const openWindow = (name,type) => {
        //
        if(dialog.value)
        {
            dialog.value = false
        }
        console.log("starting")
        //console.log(name,type)
        selectedFunction.value = name
        selectedType.value = type.type
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
        dialog


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

