<template>
    <div>
        <h1 class="text-h4">XRP Account Balances</h1>

    </div>
</template>
<script lang="ts">
import { computed, defineComponent, ref, useContext } from '@nuxtjs/composition-api'
import { plainToClass } from 'class-transformer'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPAccountBalancesGQL } from '~/apollo/queries'

//import { State } from '~/types/state'
//import { result } from '~/composables/useXrpAccounts'

interface XRPBalanceElem {
    icon: string
    issuerAddress: string
    currency: string
    name: string
    balance: number
    price: number
    value: number
}

export default defineComponent({
        
    setup() {
        const { $f } = useContext()
    const loading = ref(true)
    const { onResult } = useQuery(XRPAccountBalancesGQL, () => ({address:'rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj'}), { fetchPolicy: 'no-cache'})

    const balancesRawData = ref<XRPBalanceElem[]>([])
        const screenerDataFormatted = computed(() =>
        balancesRawData.value.map((elem) => ({
        ...elem,
        currencyShort: elem.currency.length > 20 ? elem.currency.substring(0, 20) + '...' : elem.currency,
        issuerAddressShort: `${elem.issuerAddress.slice(0, 10)}........${elem.issuerAddress.slice(
          elem.issuerAddress.length - 10,
          elem.issuerAddress.length
        )}`,
        priceFormatted: $f(elem.price, { minDigits: 6, after: '' }),
      }))
    )

        const cols = computed(() => {
            return [
            {
          text: 'Currency',
          align: 'left',
          value: 'currencySymbol',
          width: '150',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate'],
          sortable: true,
        },
        {
          text: 'Issuer Name',
          align: 'left',
          value: 'issuerName',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Issuer Address',
          align: 'left',
          value: 'issuerAddressShort',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
        },

        {
          text: 'Price (XRP)',
          align: 'left',
          value: 'priceFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Balance',
          align: 'left',
          value: 'priceFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
            ]
        })
    
        return {cols}
    
    }
       
})

</script>