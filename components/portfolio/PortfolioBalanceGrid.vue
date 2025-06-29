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
import { useHelpers } from '~/composables/useHelpers'

type Props = {
  data: Balance
  height?: number
}

export default defineComponent({
  props: {
    data: { type: Object as () => Balance, required: true },
    height: { type: Number, default: 400 },
  },
  setup(props) {
    // COMPOSABLES
    const { getCurrentNetwork } = useHelpers()
    // DATA
    const columns = [
      { text: 'Token', value: 'symbol', width: '30%' },
      { text: 'Balance', value: 'balance', width: '25%' },
      { text: 'Price', value: 'price', width: '20%' },
      { text: 'Value', value: 'value', width: '25%' },
    ]

    // COMPUTED
    const totalBalance = computed(() => props.data.items.reduce((n, { quote }) => n + quote, 0))
    const balanceItems = computed<BalanceItem[]>(() => props.data.items.filter((elem) => elem.quote > 0))
    const chainData = computed(() => getCurrentNetwork())
    return { columns, balanceItems, height: props.height, chainData, totalBalance }
  },
})
</script>
