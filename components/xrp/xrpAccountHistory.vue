<template>
    <div>
        <h2 class="text-h4">XRP Account History</h2>

    </div>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, PropType, useStore } from '@nuxtjs/composition-api'
import { plainToClass } from 'class-transformer'
import { State } from '~/types/state'
//import { result } from '~/composables/useXrpAccounts'



interface XRPAccountHistoryElem {
    hash: string
    from: string
    action: string
    to: string
    amount: number
    fee: number
}

export default defineComponent({
        
    setup() {

        const cols = computed(() => {
            return [
                {
                    text:'Hash',
                    align:'left',
                    value:'hash',
                    width:'150',
                    class: ['px-4', 'text-truncate'],
                    cellClass: ['px-4', 'text-truncate'],
                    sortable: false,
                },
                {
                    text:'From',
                    align:'left',
                    value:'from',
                    width:'150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text:'Action',
                    align:'left',
                    value:'action',
                    width:'150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text:'To',
                    align:'left',
                    value:'to',
                    width:'150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text:'Amount',
                    align:'left',
                    value:'amount',
                    width:'150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
                {
                    text:'Fee',
                    align:'left',
                    value:'fee',
                    width:'150',
                    class: ['px-2', 'text-truncate'],
                    cellClass: ['px-2', 'text-truncate'],
                    sortable: false,
                },
            ]
        })
    
        return { cols }
    
    }
       
})

</script>