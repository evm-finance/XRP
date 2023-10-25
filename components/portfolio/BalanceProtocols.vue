<template>
  <div>
    <v-row class="mt-6" no-gutters>
      <v-avatar size="30" class="mt-1">
        <v-img alt="aave logo" :src="$imageUrlBySymbol('aave')" :lazy-src="$imageUrlBySymbol('aave')" />
      </v-avatar>

      <div class="text-h5 font-weight-medium mt-1 ml-3">AAVE V2</div>
    </v-row>

    <v-row>
      <v-col v-for="(poolElem, i) in pools" :key="i">
        <v-card tile outlined height="100%">
          <v-card-title class="pa-0 ma-0">
            <v-row class="pa-2">
              <v-col class="d-flex" lg="3">
                <v-avatar size="24px">
                  <v-img
                    :src="$imageUrlBySymbol(poolElem.chain.symbol)"
                    :lazy-src="$imageUrlBySymbol(poolElem.chain.symbol)"
                  />
                </v-avatar>
                <h4 class="text-subtitle-1 pl-3 text-truncate">{{ poolElem.chain.name }}</h4>
              </v-col>
              <v-col class="text-right hidden-sm-and-down" lg="9">
                <span v-for="(balance, index) in poolElem.balances" :key="index" class="text-caption grey--text mr-2">
                  {{ balance.text }}:
                  <span
                    class="text-subtitle-1 text-truncate pink--text font-weight-medium"
                    v-text="$f(balance.value, { pre: '$', minDigits: 2 })"
                  />
                </span>
              </v-col>
            </v-row>
          </v-card-title>
          <v-divider />
          <v-data-table
            id="balances-grid"
            :headers="columns"
            :items="poolElem.pool"
            :sort-desc="[true]"
            height="100%"
            :items-per-page="10 * 10 ** 12"
            class="elevation-0"
            :mobile-breakpoint="0"
            hide-default-footer
          >
            <template #[`item.symbol`]="{ item }">
              <div class="text-no-wrap overflow-x-hidden">
                <v-avatar size="20" class="mr-2">
                  <v-img
                    :src="$imageUrlBySymbol(item.symbol)"
                    :lazy-src="$imageUrlBySymbol(item.symbol)"
                    @error="$setAltImageUrl"
                  />
                </v-avatar>
                <nuxt-link
                  class="text-capitalize text-decoration-none white--text"
                  :to="{
                    path: `/token/${item.symbol}`,
                    query: { contract: item.addresses.address, decimals: item.addresses.decimals },
                  }"
                  v-text="item.symbol"
                />
              </div>
            </template>

            <template #item.portfolio.walletBal="{ item }">
              <div v-text="$f(item.portfolio.walletBal, { minDigits: 2, maxDigits: 6 })" />
              <div class="grey--text" v-text="$f(item.portfolio.walletBal * item.price.priceUsd, { pre: '$ ' })" />
            </template>

            <template #item.portfolio.totalDeposits="{ item }">
              <div v-text="$f(item.portfolio.totalDeposits, { minDigits: 2, maxDigits: 6 })" />
              <div class="grey--text" v-text="$f(item.portfolio.totalDeposits * item.price.priceUsd, { pre: '$ ' })" />
            </template>

            <template #item.portfolio.variableBorrow="{ item }">
              <div v-text="$f(item.portfolio.variableBorrow, { minDigits: 2, maxDigits: 6 })" />
              <div class="grey--text" v-text="$f(item.portfolio.variableBorrow * item.price.priceUsd, { pre: '$ ' })" />
            </template>

            <template #item.depositAPY="{ item }">
              <div class="grey--text">
                <span v-if="item.depositAPY === -1">--</span>
                <span v-else v-text="$f(item.depositAPY * 100, { minDigits: 2, after: ' %' })" />
              </div>
            </template>

            <template #item.variableBorrowAPY="{ item }">
              <div class="grey--text">
                <span v-if="item.variableBorrowAPY === -1">--</span>
                <span v-else v-text="$f(item.variableBorrowAPY * 100, { minDigits: 2, after: ' %' })" />
              </div>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, PropType, useStore } from '@nuxtjs/composition-api'
import { plainToClass } from 'class-transformer'
import { Balance } from '~/types/apollo/main/types'
import { State } from '~/types/state'
import { AavePoolModel } from '~/composables/useAavePools'
import { Network } from '~/types/global'

type Props = {
  balances: Balance[]
}

export default defineComponent<Props>({
  props: {
    balances: {
      type: Array as PropType<Balance[]>,
      default: () => [],
    },
  },

  setup(props) {
    const { getters } = useStore<State>()
    const columns = [
      {
        text: 'Token',
        align: 'start',
        value: 'symbol',
        class: 'px-2',
        width: 60,
      },

      {
        text: 'Balance',
        align: 'left',
        value: 'portfolio.walletBal',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'py-1', 'text-truncate'],
      },

      {
        text: 'Deposits',
        align: 'left',
        value: 'portfolio.totalDeposits',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'py-1', 'text-truncate'],
      },

      {
        text: ' Borrows',
        align: 'left',
        value: 'portfolio.variableBorrow',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'py-1', 'text-truncate'],
      },
      {
        text: 'Deposit APY',
        align: 'center',
        value: 'depositAPY',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'py-1', 'text-truncate'],
      },
      {
        text: 'Borrow APY',
        align: 'center',
        value: 'variableBorrowAPY',
        class: ['px-2', 'text-truncate'],
        cellClass: ['px-2', 'py-1', 'text-truncate'],
      },
    ]

    const pools: ComputedRef<{ chain: Network; pool: AavePoolModel[]; balances: { text: string; value: number }[] }[]> =
      computed(() =>
        props.balances
          .map((elem) => ({
            chain: getters['configs/chainInfo'](elem.chainId),
            pool: plainToClass(AavePoolModel, elem.aavePools as AavePoolModel[]),
            balances: [
              {
                text: 'Net Worth',
                value:
                  elem.aavePools.reduce((a, b) => a + b.price.priceUsd * b.portfolioVal.totalDeposits, 0) -
                  elem.aavePools.reduce((a, b) => a + b.price.priceUsd * b.portfolioVal.variableBorrow, 0),
              },
              {
                text: 'Deposits',
                value: elem.aavePools.reduce((a, b) => a + b.price.priceUsd * b.portfolioVal.totalDeposits, 0),
              },
              {
                text: 'Borrows',
                value: elem.aavePools.reduce((a, b) => a + b.price.priceUsd * b.portfolioVal.variableBorrow, 0),
              },
            ],
          }))
          .filter((elem) => elem.pool.length > 0)
      )

    return { columns, pools }
  },
})
</script>
