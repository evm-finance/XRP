<template>
  <v-card tile outlined>
    <v-card-title class="pa-0 ma-0">
      <v-col class="d-flex">
        <v-avatar size="24px">
          <v-img :src="$imageUrlBySymbol('xrp')" />
        </v-avatar>
        <h4 class="text-subtitle-1 pl-3 text-truncate">XRP Ledger</h4>
      </v-col>
      <v-col cols="4" class="text-right">
        <h4
          class="text-subtitle-1 text-truncate pink--text font-weight-medium"
          v-text="$f(totalBalance, { pre: '$ ', minDigits: 2 })"
        />
      </v-col>
    </v-card-title>
    <v-divider />
    
    <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />
    
    <v-data-table
      v-else
      id="xrp-balances-grid"
      :headers="columns"
      :items="balanceItems"
      :sort-desc="[true]"
      :height="height"
      :items-per-page="10 * 10 ** 12"
      class="elevation-0"
      :mobile-breakpoint="0"
      hide-default-footer
    >
      <template #[`item.currency`]="{ item }">
        <div class="text-no-wrap overflow-x-hidden">
          <v-avatar size="20" class="mr-2">
            <img :src="$imageUrlBySymbol(item.currency.toLowerCase())" alt="" />
          </v-avatar>
          <span class="text-capitalize white--text">
            {{ item.currency }}
          </span>
        </div>
      </template>

      <template #[`item.balance`]="{ item }">
        <span
          class="grey--text"
          v-text="formatBalance(item.balance)"
        />
      </template>

      <template #[`item.price`]="{ item }">
        <span
          class="grey--text"
          v-text="$f(item.price, { pre: '$ ', minDigits: 2, maxDigits: 6 })"
        />
      </template>

      <template #[`item.value`]="{ item }">
        <span v-text="$f(item.value, { pre: '$ ', minDigits: 2, maxDigits: 4 })" />
      </template>
    </v-data-table>
    
    <div v-if="!loading && balanceItems.length === 0" class="text-center pa-4">
      <v-icon size="48" color="grey lighten-2">mdi-wallet-outline</v-icon>
      <div class="text-subtitle-2 grey--text mt-2">No XRP balances</div>
      <div class="text-caption grey--text">Connect your XRP wallet to view balances</div>
    </div>
    
    <v-divider />
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, inject, onMounted } from '@nuxtjs/composition-api'
import { XRPAccountBalancesGQL } from '~/apollo/queries'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import useXrpGraphQLWithLogging from '~/composables/useXrpGraphQLWithLogging'

interface XRPBalanceItem {
  currency: string
  name: string
  balance: number
  price: number
  value: number
  issuer?: string
}

export default defineComponent({
  setup() {
    const { $f } = useContext()
    const { address, isWalletReady } = inject(XRP_PLUGIN_KEY) as XrpClient
    
    const loading = ref(true)
    const balanceData = ref<XRPBalanceItem[]>([])
    
    // Use connected wallet address or fallback to default
    const accountAddress = computed(() => address.value || 'rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj')
    
    // Log the query content BEFORE making the call
    console.log('🚀 [BEFORE QUERY] XrpBalanceWidget - XRPAccountBalancesGQL:', {
        query: XRPAccountBalancesGQL.loc?.source.body,
        variables: { account: accountAddress.value },
        timestamp: new Date().toISOString()
    })

    const { useLoggedQuery } = useXrpGraphQLWithLogging()
    const { onResult } = useLoggedQuery(
      XRPAccountBalancesGQL, 
      () => ({ account: accountAddress.value }), 
      { 
          fetchPolicy: 'no-cache',
          context: {
              queryName: 'XRPAccountBalances',
              component: 'XrpBalanceWidget',
              purpose: 'XRP balance widget display'
          }
      }
    )

    const columns = [
      {
        text: 'Token',
        align: 'start',
        value: 'currency',
        class: 'px-2',
        width: 60,
      },
      {
        text: 'Balance',
        align: 'start',
        value: 'balance',
        class: 'px-2',
        width: 60,
      },
      {
        text: 'Price',
        align: 'start',
        value: 'price',
        width: 80,
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'text-truncate'],
      },
      {
        text: 'Value',
        align: 'start',
        value: 'value',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'text-truncate'],
        width: 90,
      },
    ]
    
    const height = 450

    const totalBalance = computed(() => balanceData.value.reduce((sum, item) => sum + item.value, 0))
    const balanceItems = computed(() => balanceData.value.filter((item) => item.value > 0))

    const formatBalance = (balance: number): string => {
      if (balance === 0) return '0'
      if (balance < 0.000001) return '< 0.000001'
      return balance.toFixed(6)
    }

    // Handle query results
    onResult((queryResult: any) => {
      if (queryResult.data?.xrpAccountBalances) {
        const data = queryResult.data.xrpAccountBalances
        balanceData.value = data.xrpBalances || []
        
        // Add XRP balance if available
        if (data.xrpBalance > 0) {
          balanceData.value.unshift({
            currency: 'XRP',
            name: 'XRP',
            balance: data.xrpBalance,
            price: data.xrpPrice || 0,
            value: data.xrpBalance * (data.xrpPrice || 0)
          })
        }
      }
      loading.value = false
    })

    onMounted(() => {
      // Initial load
    })

    return {
      columns,
      balanceItems,
      height,
      totalBalance,
      loading,
      formatBalance,
      isWalletReady
    }
  }
})
</script> 