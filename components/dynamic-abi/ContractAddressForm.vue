// This Component will be a generic input form for smart contracts
// There will be rules enforced for address, uint, string, and bytes types
// Type will be passed in as props

// This component will accept a string from an input form and interpret it as a JSON object

<template>
    <div >
        <v-row>
            <h3 style="margin-right:10px; color: white"> Enter Contract Address: </h3>
            <form @submit.prevent="enterContract">
                <textarea v-model="inputString" type="text" class="contractAddressForm" style="margin-right:10px; height:30px; color: rgba(255, 255, 255, 0.89);"></textarea>
                <v-btn outlined type="enterContract" style="color: red; margin-bottom:25px;">submit</v-btn>

            </form>

        </v-row>
    </div>
</template>

<script>

import { ref } from 'vue';
import { defineComponent } from '@nuxtjs/composition-api';

export default defineComponent({
    setup(_,{ emit }) {
        const inputString = ref('');

        const enterContract = () => {
            try {
                // console.log(jsonData)
                //console.log('hello world')
                console.log(inputString.value)
                emit('addressEntered', inputString.value)
            }
            catch (error) {
                console.error('Invalid Address');
            }
        }

        return {
            inputString,
            enterContract,
        }
    }
})

</script>

<style>
.contractAddressForm {
    background-color: rgba(255, 255, 255, 0.048);
    word-wrap: break-word;
    font-size: 18px;
    width: 500px;
}

</style>
