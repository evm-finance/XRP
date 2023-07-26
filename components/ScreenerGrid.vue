<template>
  <div>
    <v-card tile outlined height="1400">
      <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />

      <client-only>
        <v-data-table
          v-if="!loading"
          id="curve-pools-grid"
          hide-default-footer
          :headers="cols"
          :items="screenerDataFormatted"
          :items-per-page="25"
          class="elevation-0"
          :single-expand="false"
          item-key="address"
          mobile-breakpoint="0"
          @update:options="(options) => loadPage(options)"
        >
          <!--      Cols-->
          <template #item.rank="{ item }">
            <div class="my-1">
              <v-row no-gutters align="center">
                <v-col cols="2">
                  <v-avatar size="24" class="ml-2">
                    <img :src="$imageUrlBySymbol(item.token0Symbol)" alt="" @error="$setAltImageUrl" />
                  </v-avatar>
                </v-col>
                <v-col>
                  <v-row no-gutters>
                    <v-col>
                      <nuxt-link
                        class="text-capitalize font-weight-bold pink--text text-decoration-none"
                        :to="{
                          path: `/token/${item.token0Symbol}`,
                          query: {
                            name: item.token0Symbol,
                            contract: item.token0Address,
                            decimals: item.token0Decimals,
                          },
                        }"
                        v-text="item.token0Name"
                      />
                    </v-col>
                  </v-row>
                  <v-row no-gutters>
                    <v-col>
                      <span class="grey--text text-caption" v-text="`${item.token0Symbol}/${item.token1Symbol}`" />
                      <v-tooltip bottom color="black">
                        <template #activator="{ on, attrs }">
                          <v-btn
                            icon
                            color="grey"
                            x-small
                            v-bind="attrs"
                            v-on="on"
                            @click="copyAddressToClipboard(item.address)"
                          >
                            <v-icon size="12">mdi-content-copy</v-icon>
                          </v-btn>
                        </template>
                        <span v-text="'Copy Address'" />
                      </v-tooltip>
                    </v-col>
                  </v-row>
                </v-col>
              </v-row>
            </div>
          </template>

          <template #item.token_0_price_usd="{ item }">
            <!--          {{ item.price }}-->

            <span :class="`${item.priceUpdateOptions.color}--text`">
              {{ item.price }}
              <v-icon v-if="item.priceUpdateOptions.status" :color="item.priceUpdateOptions.color" size="14"
                >{{ item.priceUpdateOptions.icon }}
              </v-icon>
            </span>
          </template>

          <!--          <template #item.balance="{ item }">
            <div v-text="item.balance" />
            <div class="grey&#45;&#45;text" v-text="item.balanceValue" />
            <div class="grey&#45;&#45;text" v-text="item.balanceValue." />
          </template>
          -->
          <template #item.balanceNew="{ item }">
            <div v-if="item.balanceNew">
              <div v-if="!item.balanceNew.isWrapped">
                <div v-text="item.balanceNew.balance" />
                <div
                  class="grey--text"
                  v-text="$f(item.balanceNew.balance * item.token0PriceUSD, { minDigits: 2, pre: '$ ' })"
                />
              </div>

              <div v-if="item.balanceNew.isWrapped">
                <v-row no-gutters>
                  <v-col>
                    <span class="grey--text" style="display: inline-block; width: 50px">
                      {{ item.balanceNew.nativeSymbol }}
                    </span>
                    <span class="mx-3">{{ $f(item.balanceNew.nativeBalance, { minDigits: 6 }) }}</span>
                    <span class="grey--text">
                      {{ $f(item.balanceNew.nativeBalance * item.token0PriceUSD, { minDigits: 4, pre: '$ ' }) }}
                    </span>
                  </v-col>
                </v-row>
                <v-row no-gutters>
                  <v-col>
                    <span class="grey--text" style="display: inline-block; width: 50px">
                      {{ item.balanceNew.uniqueKeyOrSymbol }}
                    </span>
                    <span class="mx-3">{{ $f(item.balanceNew.balance, { minDigits: 6 }) }}</span>
                    <span class="grey--text">
                      {{ $f(item.balanceNew.balance * item.token0PriceUSD, { minDigits: 4, pre: '$ ' }) }}
                    </span>
                  </v-col>
                </v-row>
              </div>
            </div>
            <div v-else>--</div>
          </template>

          <template #item.change_5min="{ item }">
            <span :class="`${item.change5min.color}--text`">
              {{ item.change5min.value }}
              <v-icon v-if="item.change5min.icon" :color="item.change5min.color" size="14"
                >{{ item.change5min.icon }}
              </v-icon>
            </span>
          </template>

          <template #item.change_1h="{ item }">
            <span :class="`${item.change1h.color}--text`">
              {{ item.change1h.value }}
              <v-icon v-if="item.change1h.icon" :color="item.change1h.color" size="14"
                >{{ item.change1h.icon }}
              </v-icon>
            </span>
          </template>

          <template #item.change_24h="{ item }">
            <span :class="`${item.change24h.color}--text`">
              {{ item.change24h.value }}
              <v-icon v-if="item.change24h.icon" :color="item.change24h.color" size="14"
                >{{ item.change24h.icon }}
              </v-icon>
            </span>
          </template>

          <template #item.aavePools="{ item }">
            <v-menu v-if="item.aavePool" left nudge-left="50">
              <template #activator="{ on: menu, attrs }">
                <v-tooltip top color="black">
                  <template #activator="{ on: tooltip }">
                    <v-btn color="primary" icon dark v-bind="attrs" v-on="{ ...tooltip, ...menu }">
                      <v-avatar size="20">
                        <img :src="$imageUrlBySymbol('aave')" alt="" @error="$setAltImageUrl" />
                      </v-avatar>
                    </v-btn>
                  </template>
                  <div class="text-center">
                    <div>{{ item.token0Symbol }} can be used in Aave V2.</div>
                    <div>Click for pool details</div>
                  </div>
                </v-tooltip>
              </template>
              <v-card width="1000" outlined>
                <aave-market-details :pool="item.aavePool" :show-balance-chart="false" />
                <div class="text-center mb-2">
                  <nuxt-link to="/markets/aave" class="text-center text-caption grey--text"
                    >Navigate to AAVE dashboard
                    <v-icon size="12" class="ml-1" color="grey">mdi-open-in-new</v-icon></nuxt-link
                  >
                </div>
              </v-card>
            </v-menu>
          </template>
        </v-data-table>
        <v-pagination v-model="currentPage" :length="currentPage + 1" class="mt-4" />
      </client-only>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs, useContext, watchEffect } from '@nuxtjs/composition-api'
import { Pool, PriceStream } from '~/types/apollo/main/types'
import useTokenScreener from '~/composables/useTokenScreener'
import { AavePoolModel } from '~/composables/useAavePools'
import AaveMarketDetails from '~/components/pools/AaveMarketDetails.vue'

type ScreenerItemObserver = Pool & {
  priceUpdateOptions: { status: boolean; color: string | null; icon: string | null }
}
type Props = {
  search: string | null
  networkId: string
  dexId: string
  aavePools: { [a: string]: AavePoolModel }
}
export default defineComponent<Props>({
  components: { AaveMarketDetails },
  props: {
    search: { type: String as PropType<string | null>, default: null, required: false },
    networkId: { type: String as PropType<string>, default: 'ethereum' },
    dexId: { type: String as PropType<string>, default: 'uniswap_v3' },
    aavePools: { type: Object as PropType<{ [a: string]: AavePoolModel }>, default: () => ({}) },
  },

  setup(props) {
    // COMPOSABLE
    const { $copyAddressToClipboard, $f, $applyPtcChange, $emitter } = useContext()
    let priceUpdateTimeout: any = null

    // STATE
    const screenerDataObserver = ref<ScreenerItemObserver[]>([])
    const sortBy = ref<string>('rank')
    const order = ref<string>('asc')
    const dex = toRefs(props).dexId
    const network = toRefs(props).networkId
    const { screenerData, currentPage, balanceMap, loading } = useTokenScreener(network, dex, sortBy, order)

    watchEffect(
      () =>
        (screenerDataObserver.value = screenerData.value.map((elem) => ({
          ...elem,
          priceUpdateOptions: { status: false, color: null, icon: null },
        })))
    )

    const screenerDataFormatted = computed(() =>
      screenerDataObserver.value.map((elem) => ({
        ...elem,
        price: $f(elem.token0PriceUSD, { pre: '$ ', minDigits: 2, maxDigits: 6 }),
        balanceNew: Object.prototype.hasOwnProperty.call(balanceMap.value, elem.token0Address)
          ? balanceMap.value[elem.token0Address]
          : null,
        aavePool: Object.prototype.hasOwnProperty.call(props.aavePools, elem.token0Address.toLowerCase())
          ? props.aavePools[elem.token0Address]
          : null,
        change5min: $applyPtcChange(elem.change5Min),
        change1h: $applyPtcChange(elem.change1h),
        change24h: $applyPtcChange(elem.change24h),
      }))
    )

    const updateData = (input: PriceStream[]) => {
      clearTimeout(priceUpdateTimeout)
      input.forEach((inputElem) => {
        const currentGridRaw = screenerDataObserver.value.find((elem) => elem.address === inputElem.pairAddress)
        if (currentGridRaw) {
          if (inputElem.token0PriceUSD > currentGridRaw.token0PriceUSD) {
            currentGridRaw.priceUpdateOptions = { status: true, color: 'green', icon: 'mdi-arrow-up' }
          } else {
            currentGridRaw.priceUpdateOptions = { status: true, color: 'red', icon: 'mdi-arrow-down' }
          }
          currentGridRaw.token0PriceUSD = inputElem.token0PriceUSD
        }
      })

      priceUpdateTimeout = setTimeout(() => {
        input.forEach((inputElem) => {
          const currentGridRaw = screenerDataObserver.value.find((elem) => elem.address === inputElem.pairAddress)
          if (currentGridRaw) {
            currentGridRaw.priceUpdateOptions = { status: false, color: null, icon: null }
          }
        })
      }, 2000)
    }

    const loadPage = (options: any) => {
      const opt = { sortBy: 'rank', order: 'asc' }

      if (options.sortBy.length > 0) {
        opt.sortBy = options.sortBy[0]
      }
      if (options.sortDesc.length > 0) {
        opt.order = options.sortDesc[0] ? 'desc' : 'asc'
      }

      sortBy.value = opt.sortBy
      order.value = opt.order
    }

    $emitter?.on('priceStream', (val: any) => {
      const prices: PriceStream[] = val.priceStream ?? []
      updateData(prices)
    })

    const cols = computed(() => {
      return [
        {
          text: 'Assets',
          align: 'left',
          value: 'rank',
          width: '400',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          sortable: true,
        },

        {
          text: 'Price',
          align: 'left',
          value: 'token_0_price_usd',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          width: '200',
          sortable: true,
        },

        {
          text: 'Your Balance',
          align: 'left',
          value: 'balanceNew',
          width: '200',
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: '',
          align: 'left',
          value: 'priceUpdated.status',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          sortable: false,
        },

        {
          text: '% 5m',
          align: 'left',
          value: 'change_5min',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: '% 1h',
          align: 'left',
          value: 'change_1h',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: '% 24h',
          align: 'left',
          value: 'change_24h',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: '',
          width: '60px',
          align: 'center',
          value: 'aavePools',
          class: ['px-6', 'text-truncate'],
          cellClass: ['px-6', 'text-truncate'],
          sortable: false,
        },
      ]
    })

    // const dataLoading = toRefs(props).loading
    // const searchString = toRefs(props).search

    // METHODS
    /*    function initAction(poolAddress: string, action: actionTypes) {
      const pool: AavePoolModel | undefined = aavePools.value.find((elem) => elem.id === poolAddress)
      if (pool) {
        emit(EmitEvents.initAction, { action, pool })
      }
    } */

    /*    async function importToMetamask(poolId: string) {
      const pool: AavePoolModel | undefined = aavePools.value.find((elem) => elem.id === poolId)
      if (pool) {
        await importTokenToMetamask({
          address: pool.underlyingAsset,
          symbol: pool.symbol,
          decimals: pool.decimals,
          image: pool.logoUrl,
        })
      }
    } */

    // METHODS
    /* function navigateToExplorer(address: string) {
   /!*   const currentChain: Chain = getters['configs/chainInfo'](chainId.value ?? 1)
      const url = `${currentChain.blockExplorerUrl}address/${address}`
      window.open(url)*!/
    } */

    /*   function expendSingleRow() {
      if (props.expandFirstRow && aavePools.value.length > 0) {
        const firstRow: AavePoolModel = aavePools.value[0]
        expanded.value.push(firstRow)
      }
    } */

    // WATCHERS
    // watch(aavePools, () => expendSingleRow())

    return {
      screenerData,
      screenerDataFormatted,
      cols,
      currentPage,
      loading,
      dex,
      loadPage,
      updateData,
      copyAddressToClipboard: $copyAddressToClipboard,
    }
  },
})
</script>
