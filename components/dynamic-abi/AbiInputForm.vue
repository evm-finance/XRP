// This component will accept a string from an input form and interpret it as a JSON object

<template>
    <div>
        <h3 style="margin-bottom: 15px; color: white"> Enter ABI Below </h3>
        <v-row>
            <form @submit.prevent="submitForm">
                <textarea v-model="inputString" type="text" class="abiInput" style="margin-bottom: 10px; color: rgba(255, 255, 255, 0.89)"></textarea>
                <v-row justify="center" >
                    <v-btn type="submit" style="color: red">Generate Function Buttons</v-btn>
                </v-row> 
            </form>
        </v-row>
    </div>
</template>

<script>

import { ref } from 'vue';
import { defineComponent } from '@nuxtjs/composition-api';

export default defineComponent({
    setup(_, { emit }) {
        const inputString = ref('');
        const jsonData = ref('');

        const submitForm = () => {
            try {
                const jsonData = JSON.parse(inputString.value)
                // console.log(jsonData)
                emit('abiEntered', jsonData)
            }
            catch (error) {
                console.error('Invalid JSON string');
            }
        }

        return {
            inputString,
            submitForm,
            jsonData
        }
    }
})

</script>

<style>
.abiInput {
    background-color: rgba(255, 255, 255, 0.048);
    word-wrap: break-word;
    font-size: 14px;
    width: 1000px;
    height: 250px;
}

</style>