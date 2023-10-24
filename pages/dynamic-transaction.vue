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

export default defineComponent({
    components:{
    AbiInputForm,
    ContractAddressForm,
    FunctionButtons,
    DynamicContractUi
},
setup() {

    const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
    const { allowedSpending, approveMaxSpending } = useERC20()
    const renderButtons = ref(false)
    const dialog = ref(false)
    const address = ref('')
    const functions = ref([])
    const contractInstance = ref()
    const selectedFunction = ref('')
    const calldataParams = ref([])

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
                console.log(data[i])
                //functions.value.push(jsonData[i])
            }
        functions.value = data
        // console.log('render buttons',renderButtons.value)
        // console.log('finised generateUI')
    }

    const sendTransaction = () => {
        console.log(selectedFunction.value,calldataParams.value,signer.value,account.value)
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
        console.log(functions)
        console.log(name,type)
        try{
            const contract = new ethers.Contract(
                address.value,
                functions.value,
                signer.value
            )
            console.log(contract)
            contractInstance.value=contract
            dialog.value=true
        }
        catch (error)
        {
            console.log('failed to instantiate contract, try harder')
        }
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

