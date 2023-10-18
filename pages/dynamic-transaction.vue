<template>
    <div>
        <v-row>
            <abi-input-form></abi-input-form>
        </v-row>
    </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api';
import AbiInputForm from '~/components/dynamic-abi/AbiInputForm.vue'

export default defineComponent({
    components:{
    AbiInputForm
},
setup() {

}

})

</script>

