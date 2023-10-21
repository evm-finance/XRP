// This Component will be a generic input form for smart contracts
// There will be rules enforced for address, uint, string, and bytes types
// Type will be passed in as props

// This component will accept a string from an input form and interpret it as a JSON object

<template>
    <div >
        <v-row>
            <h3 style="margin-bottom: 25px; margin-right:10px; color: white"> Enter Contract Address: </h3>
            <form @submit.prevent="submitForm">
                <input v-model="inputString" type="text" class="contractAddressForm" style="margin-bottom: 25px; margin-right:10px; height:38px; color: rgba(255, 255, 255, 0.89);"></input>

            </form>
            <v-btn type="submit" style="color: red">submit</v-btn>
            <h3 style="margin-left: 10px; font-size: 18px;"> Network: </h3>
        </v-row>
    </div>
</template>

<script>

import { ref } from 'vue';
import { defineComponent } from '@nuxtjs/composition-api';

export default defineComponent({
    setup({ emit }) {
        const inputString = ref('');

        const submitForm = () => {
            try {
                // console.log(jsonData)
                emit('addressEntered', inputString.value)
            }
            catch (error) {
                console.error('Invalid JSON string');
            }
        }

        return {
            inputString,
            submitForm,
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
    height: 30px;
}

</style>
