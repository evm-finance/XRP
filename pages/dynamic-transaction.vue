<template>
    <div justify="center">
        <div class="padding">
            <v-row>
                <contract-address-form></contract-address-form>
            </v-row>
            <v-row>
                <abi-input-form @abiEntered="generateUi"></abi-input-form>
            </v-row>
        </div>

    </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api';
import AbiInputForm from '~/components/dynamic-abi/AbiInputForm.vue'
import ContractAddressForm from '~/components/dynamic-abi/ContractAddressForm.vue'


export default defineComponent({
    components:{
    AbiInputForm,
    ContractAddressForm
},
setup() {

    var renderButtons = false
    
    const generateUi = (jsonData) => {
        console.log(jsonData)
        renderButtons = true
        //
    }
    

    return {
        generateUi,
        renderButtons
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

