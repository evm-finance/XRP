// This component will accept a string from an input form and interpret it as a JSON object

<template>
        <v-row justify="center">Enter Contract ABI
            <form @submit.prevent="submitForm">
                <input v-model="inputString" type="text">
                <button type="submit">Submit</button>
            </form>
        </v-row> 
</template>

<script>

import { ref } from 'vue';
import { defineComponent } from '@nuxtjs/composition-api';


export default defineComponent({
    setup() {
        const inputString = ref('');
        const jsonData = ref('')

        const submitForm = () => {
            try {
                const jsonData = JSON.parse(inputString.value)
                console.log(jsonData)
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