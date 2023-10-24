<!-- // This Component will be a generic input form for smart contracts
// There will be rules enforced for address, uint, string, and bytes types
// Type will be passed in as props

// This component will accept a string from an input form and interpret it as a JSON object

<template>
    <div>
        <v-row>
            <p> Param Name: {{ paramName }} </p>
            <form @submit.prevent="submitForm">
                <input v-model="inputString" type="text" class="param_form_input" style="margin-bottom: 10px; color: rgba(255, 255, 255, 0.89)"></textarea>

            </form>
        </v-row>
    </div>
</template>

<script>

import { defineComponent, ref, props } from '@nuxtjs/composition-api';
import { stringify } from 'postcss';

export default defineComponent({
    props: {
        paramName: String,
        paramType: String
    },
    setup(props, { emit }) {
        const inputString = ref('');

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
            paramName: props.paramName, 
            paramType: props.paramType
        }
    }
})

</script>

<style>
.param_form_input {
    background-color: rgba(255, 255, 255, 0.048);
    word-wrap: break-word;
    font-size: 18px;
    width: 500px;
    height: 30px;
}

</style> -->
