<template>
  <v-card tile outlined>
    <v-card-title v-if="chainData" class="pa-0 ma-0">
      <v-col class="d-flex">
        <v-avatar size="24px">
          <v-img :src="$imageUrlBySymbol(chainData.symbol)" @error="$setAltImageUrl"></v-img>
        </v-avatar>
        <h4 class="text-subtitle-1 pl-3 text-truncate" v-text="chainData.name" />
      </v-col>
      <v-col cols="4" class="text-right">
        <h4
          class="text-subtitle-1 text-truncate pink--text font-weight-medium"
          v-text="$f(totalBalance, { pre: '$ ', minDigits: 2 })"
        />
      </v-col>
    </v-card-title>
    <v-divider />
    <v-data-table
      id="balances-grid"
      :headers="columns"
      :items="balanceItems"
      :sort-desc="[true]"
      :height="height"
      :items-per-page="10 * 10 ** 12"
      class="elevation-0"
      :mobile-breakpoint="0"
      hide-default-footer
    >
      <template #[`item.contractTickerSymbol`]="{ item }">
        <div class="text-no-wrap overflow-x-hidden">
          <v-avatar size="20" class="mr-2">
            <img :src="$imageUrlBySymbol(item.contractTickerSymbol)" alt="" @error="$setAltImageUrl" />
          </v-avatar>
          <nuxt-link
            class="text-capitalize text-decoration-none white--text"
            :to="{
              path: `/token/${item.contractTickerSymbol}`,
              query: { contract: item.contractAddress, decimals: item.contractDecimals },
            }"
          >
            {{ item.contractTickerSymbol }}
          </nuxt-link>
        </div>
      </template>

      <template #[`item.balance`]="{ item }">
        <span
          class="grey--text"
          v-text="$f(item.balance / Math.pow(10, item.contractDecimals), { minDigits: 2, maxDigits: 4 })"
        />
      </template>

      <template #[`item.quoteRate`]="{ item }">
        <span
          class="grey--text"
          v-text="item.disableQuoteRate ? '-' : $f(item.quoteRate, { pre: '$ ', minDigits: 2, maxDigits: 6 })"
        />
      </template>

      <template #[`item.quote`]="{ item }">
        <span v-text="item.disableQuoteRate ? '-' : $f(item.quote, { pre: '$ ', minDigits: 2, maxDigits: 4 })" />
      </template>
    </v-data-table>
    <v-divider />
  </v-card>
</template>
<script lang="ts">
import { computed, defineComponent, inject, PropType } from '@nuxtjs/composition-api'
import { BalanceItem, Balance } from '~/types/apollo/main/types'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'

type Props = {
  data: Balance
}

export default defineComponent<Props>({
  props: {
    data: { type: Object as PropType<Balance>, default: () => ({}) },
  },

  setup(props) {
    // COMPOSABLES
    const { getNetworkByChainNumber } = inject(WEB3_PLUGIN_KEY) as Web3
    // DATA
    const columns = [
      {
        text: 'Token',
        align: 'start',
        value: 'contractTickerSymbol',
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
        value: 'quoteRate',
        width: 80,
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'text-truncate'],
      },
      {
        text: 'Value',
        align: 'start',
        value: 'quote',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'text-truncate'],
        width: 90,
      },
    ]
    const height = 450

    // COMPUTED
    const totalBalance = computed(() => props.data.items.reduce((n, { quote }) => n + quote, 0))
    const balanceItems = computed<BalanceItem[]>(() => props.data.items.filter((elem) => elem.quote > 0))
    const chainData = computed(() => getNetworkByChainNumber(props.data.chainId))
    return { columns, balanceItems, height, chainData, totalBalance }
  },
})
</script>
