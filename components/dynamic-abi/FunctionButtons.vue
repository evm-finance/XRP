<template>
    <div>
        <div class = "function_buttons">
            <p style="font-size: 18px; font-weight: bold;"> Action Methods: Send a transaction to the contract </p>
            <v-row>
                <v-btn 
                    v-for="func in functionMethods"
                    :key="func.name"
                    outlined
                    text
                    style="color: red"
                    @click="initFunction(func.name,func)"> {{ func.name }}
                </v-btn>
            </v-row>
        </div>
        <div class = "view_buttons">
            <p style="font-size: 18px; font-weight: bold;"> View Methods: Read data from the contract </p>
            <v-row>
                <v-btn 
                    outlined
                    v-for="item in viewMethods"
                    :key="item.name"
                    text
                    style="color: red"
                    @click="initFunction(item.name,item)"> {{ item.name }}
                </v-btn>
            </v-row>
        </div>
        <div class = "event_buttons">
            <p style="font-size: 18px; font-weight: bold;"> Events: Look up an event emitted by the contract </p>
            <v-row>
                <v-btn 
                    outlined
                    v-for="event in events"
                    :key="event.name"
                    text
                    style="color: red"
                    @click="initFunction(event.name,event)"> {{ event.name }}
                </v-btn>
            </v-row>
        </div>
    </div>
</template>

<script>
import { defineComponent, props, ref } from '@nuxtjs/composition-api';
import { AbiElem, CalldataAbi, ContractAbi } from '~/types/apollo/main/types'

export default defineComponent({

    props: {
        functionNames: {type: Array, required: true}
    },
    setup(props, {emit})
    {

        const functionMethods = ref([])
        const viewMethods = ref([])
        const events = ref([])
        const initFunction = (name, type) => {
            console.log(props.functionNames)
            emit('functionSelected',name,type.type)
            //console.log(name,type.type)
        }
        //0xdAC17F958D2ee523a2206206994597C13D831ec7 USDT on ETH
        //0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7 WAVAX
        const parseProps = () => {
            // console.log("started parsing props")
            // console.log('props',props)
            // console.log('props.name',props.name)
            // console.log('props.functionNames',props.functionNames)
            for(let i = 0; i < props.functionNames.length; i++)
            {
                //console.log(props.functionNames)
                if(props.functionNames[i].constant == true) {
                    viewMethods.value.push(props.functionNames[i])
                }
                else {
                    if(props.functionNames[i].type == "event") {
                        events.value.push(props.functionNames[i])
                    }
                    else {
                        functionMethods.value.push(props.functionNames[i])
                    }
                }
            }
            console.log(functionMethods.value,viewMethods.value)
        }

        parseProps()

        return {
            initFunction,
            functionMethods,
            viewMethods,
            events
        }
    }
})
</script>
<style>
.function_buttons {
    margin-top: 30px;
}
.view_buttons {
    margin-top: 30px;
}
.event_buttons {
    margin-top: 30px;
}
</style>