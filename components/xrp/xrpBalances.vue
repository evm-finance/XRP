<template>
  <v-card tile outlined>
    <v-card-title class="pa-0 ma-0">
      <v-col class="d-flex">
        <v-avatar size="24px">
          <v-img src="/img/xrp-logo.png" @error="$setAltImageUrl"></v-img>
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
    
    <!-- Wallet Connection Status -->
    <v-card-text v-if="!isWalletReady" class="text-center pa-4">
      <v-btn color="primary" @click="connectWallet" :loading="connecting">
        Connect GEM Wallet
      </v-btn>
      <p class="text-caption mt-2 grey--text">Connect your GEM wallet to view XRP balances</p>
    </v-card-text>

    <!-- Loading State -->
    <v-card-text v-else-if="loading" class="text-center pa-4">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
      <p class="text-caption mt-2 grey--text">Loading XRP balances...</p>
    </v-card-text>

    <!-- Balances Table -->
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
            <img :src="$imageUrlBySymbol(item.currency)" alt="" @error="$setAltImageUrl" />
          </v-avatar>
          <span class="text-capitalize">{{ item.currency }}</span>
        </div>
      </template>

      <template #[`item.balance`]="{ item }">
        <span class="grey--text" v-text="$f(item.balance, { minDigits: 2, maxDigits: 6 })" />
      </template>

      <template #[`item.price`]="{ item }">
        <span class="grey--text" v-text="$f(item.price, { pre: '$ ', minDigits: 2, maxDigits: 6 })" />
      </template>

      <template #[`item.value`]="{ item }">
        <span v-text="$f(item.value, { pre: '$ ', minDigits: 2, maxDigits: 4 })" />
      </template>
    </v-data-table>
    <v-divider />
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, inject } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPAccountBalancesGQL } from '~/apollo/queries'
import { XrpClient, XRP_PLUGIN_KEY } from '~/plugins/web3/xrp.client'

interface XRPBalanceElem {
  issuer: string
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
    const connecting = ref(false)
    
    // Inject XRP wallet plugin
    const xrpClient = inject(XRP_PLUGIN_KEY) as XrpClient
    const isWalletReady = computed(() => xrpClient?.isWalletReady?.value || false)
    const walletAddress = computed(() => xrpClient?.address?.value || '')
    
    // GraphQL query for XRP balances
    const { onResult } = useQuery(
      XRPAccountBalancesGQL, 
      () => ({ account: walletAddress.value || 'rMjRc6Xyz5KHHDizJeVU63ducoaqWb1NSj' }), 
      { fetchPolicy: 'no-cache', enabled: isWalletReady }
    )

    const balancesRawData = ref<XRPBalanceElem[]>([])
    
    // Format balance data
    const balanceItems = computed(() =>
      balancesRawData.value.map((elem) => ({
        ...elem,
        currencyShort: elem.currency.length > 20 ? elem.currency.substring(0, 20) + '...' : elem.currency,
        issuerAddressShort: `${elem.issuer.slice(0, 10)}........${elem.issuer.slice(
          elem.issuer.length - 10,
          elem.issuer.length
        )}`,
        priceFormatted: $f(elem.price, { minDigits: 6, after: '' }),
      }))
    )

    // Calculate total balance
    const totalBalance = computed(() => 
      balanceItems.value.reduce((n, { value }) => n + value, 0)
    )

    // Table columns
    const columns = [
      {
        text: 'Currency',
        align: 'start',
        value: 'currency',
        class: 'px-2',
        width: 60,
      },
      {
        text: 'Token Name',
        align: 'start',
        value: 'name',
        class: 'px-2',
        width: 100,
      },
      {
        text: 'Balance',
        align: 'start',
        value: 'balance',
        class: 'px-2',
        width: 80,
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

    // Connect wallet function
    const connectWallet = async () => {
      if (xrpClient?.connectWallet) {
        connecting.value = true
        try {
          await xrpClient.connectWallet()
        } catch (error) {
          console.error('Failed to connect wallet:', error)
        } finally {
          connecting.value = false
        }
      }
    }

    // Handle GraphQL results
    onResult((queryResult: any) => {
      balancesRawData.value = queryResult.data?.xrpAccountBalances?.xrpBalances || []
      loading.value = false
    })

    return {
      columns,
      balanceItems,
      height,
      loading,
      connecting,
      isWalletReady,
      walletAddress,
      totalBalance,
      connectWallet
    }
  }
})
</script>