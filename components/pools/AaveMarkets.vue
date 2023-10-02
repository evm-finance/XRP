<template>
  <v-card tile outlined>
    <v-skeleton-loader v-if="poolsLoading" type="table-tbody,table-tbody,table-tbody" />
    <client-only>
      <v-data-table
        v-if="!poolsLoading"
        id="curve-pools-grid"
        hide-default-footer
        :headers="cols"
        :items="aavePools"
        :sort-by="['usdBalance']"
        :sort-desc="[true]"
        :items-per-page="100"
        :search="searchString"
        class="elevation-0"
        :expanded.sync="expanded"
        :single-expand="false"
        show-expand
        item-key="name"
        mobile-breakpoint="0"
      >
        <template #expanded-item="{ headers, item }">
          <td :colspan="headers.length"><aave-market-details :pool="item" /></td>
        </template>

        <!--      Cols-->
        <template #item.symbol="{ item }">
          <div class="my-1">
            <v-row no-gutters align="center">
              <v-col cols="3">
                <v-avatar size="24" class="mr-2">
                  <v-img :alt="`${item.symbol} logo`" :src="item.logoUrl" :lazy-src="item.logoUrl" />
                </v-avatar>
              </v-col>
              <v-col>
                <v-row no-gutters>
                  <v-col>
                    <nuxt-link
                      class="text-capitalize font-weight-bold pink--text text-decoration-none d-inline-block text-truncate"
                      style="max-width: 150px"
                      :to="{
                        path: `/token/${item.symbol}`,
                        query: { name: item.name, contract: item.addresses.address, decimals: item.addresses.decimals },
                      }"
                    >
                      {{ item.name }}
                    </nuxt-link>
                  </v-col>
                </v-row>
                <v-row no-gutters>
                  <v-col>
                    <span :class="[textClass, 'text-caption']" v-text="item.symbol" />
                    <v-tooltip v-if="item.symbol !== 'ETH'" bottom color="black">
                      <template #activator="{ on, attrs }">
                        <v-btn icon color="grey" x-small v-bind="attrs" v-on="on" @click="importToMetamask(item.id)">
                          <v-icon size="14">mdi-wallet</v-icon>
                        </v-btn>
                      </template>
                      <span v-text="'Import To Metamask Wallet'" />
                    </v-tooltip>
                  </v-col>
                </v-row>
              </v-col>
            </v-row>
          </div>
        </template>

        <template #item.portfolio.walletBal="{ item }">
          <div v-text="$f(item.portfolio.walletBal, { minDigits: 2, maxDigits: 6 })" />
          <div
            :class="textClass"
            v-text="$f(item.portfolio.walletBal * item.price.priceUsd, { minDigits: 2, pre: '$ ' })"
          />
        </template>

        <template #item.portfolio.totalDeposits="{ item }">
          <div v-text="$f(item.portfolio.totalDeposits, { minDigits: 2, maxDigits: 6 })" />
          <div
            :class="textClass"
            v-text="$f(item.portfolio.totalDeposits * item.price.priceUsd, { minDigits: 2, pre: '$ ' })"
          />
        </template>

        <template #item.portfolio.variableBorrow="{ item }">
          <div v-text="$f(item.portfolio.variableBorrow, { minDigits: 2, maxDigits: 6 })" />
          <div
            :class="textClass"
            v-text="$f(item.portfolio.variableBorrow * item.price.priceUsd, { minDigits: 2, pre: '$ ' })"
          />
        </template>

        <template #item.price.priceUsd="{ item }">
          <div v-text="$f(item.price.priceUsd, { minDigits: 2, pre: '$ ' })" />
        </template>

        <template #item.depositAPY="{ item }">
          <div :class="textClass">
            <span v-if="item.depositAPY === -1">--</span>
            <span v-else v-text="$f(item.depositAPY * 100, { minDigits: 2, after: ' %' })" />
          </div>
        </template>

        <template #item.variableBorrowAPY="{ item }">
          <div :class="textClass">
            <span v-if="item.variableBorrowAPY === -1">--</span>
            <span v-else v-text="$f(item.variableBorrowAPY * 100, { minDigits: 2, after: ' %' })" />
          </div>
        </template>

        <template #item.usdBalance="{ item }">
          <span v-text="$f(item.usdBalance, { minDigits: 2, pre: '$ ', useSymbol: true })" />
        </template>

        <template #item.totalBorrowBalanceUsd="{ item }">
          <span v-text="$f(item.totalBorrowBalanceUsd, { minDigits: 2, pre: '$ ', useSymbol: true })" />
        </template>

        <template #item.link="{ item }">
          <v-tooltip bottom color="black">
            <template #activator="{ on, attrs }">
              <v-btn
                color="primary"
                v-bind="attrs"
                x-small
                class="mb-1 ml-1"
                icon
                v-on="on"
                @click="$copyAddressToClipboard(item.price.id)"
              >
                <v-icon size="18">mdi-content-copy</v-icon>
              </v-btn>
            </template>
            Copy Pool Address
          </v-tooltip>
          <v-tooltip bottom color="black">
            <template #activator="{ on, attrs }">
              <v-btn
                color="primary"
                v-bind="attrs"
                x-small
                class="mb-1 ml-1"
                icon
                v-on="on"
                @click="navigateToExplorer(item.price.id)"
              >
                <v-icon size="18">mdi-open-in-new</v-icon>
              </v-btn>
            </template>
            Open in Block Explorer
          </v-tooltip>
        </template>

        <template #[`item.action`]="{ item }">
          <v-btn
            v-for="action in aaveActions"
            :key="action"
            :disabled="!walletReady"
            text
            outlined
            color="pink"
            class="pa-1 ma-1"
            height="26"
            @click="initAction(item.id, action)"
          >
            <span class="text-caption">{{ action }}</span>
          </v-btn>
        </template>
      </v-data-table>
    </client-only>
  </v-card>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  PropType,
  ref,
  toRefs,
  UnwrapRef,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { AavePoolModel, aaveActions, actionTypes } from '~/composables/useAavePools'
import { State } from '~/types/state'
import AaveMarketDetails from '~/components/pools/AaveMarketDetails.vue'
import { EmitEvents } from '~/types/events'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Chain } from '~/types/apollo/main/types'
type Props = {
  pools: AavePoolModel[]
  loading: boolean
  search: string | null
  hideActionCols: boolean
  expandFirstRow: boolean
}
export default defineComponent<Props>({
  components: { AaveMarketDetails },
  props: {
    pools: { type: Array as PropType<UnwrapRef<AavePoolModel[]>>, default: () => [], required: true },
    loading: { type: Boolean, default: true },
    search: { type: String as PropType<string | null>, default: null, required: false },
    hideActionCols: { type: Boolean, default: false },
    expandFirstRow: { type: Boolean, default: false },
  },
  setup(props, { emit }) {
    // COMPOSABLE
    const { state, getters } = useStore<State>()
    const { walletReady, importTokenToMetamask, chainId } = inject(WEB3_PLUGIN_KEY) as Web3

    // COMPUTED
    const textClass = computed(() => state.ui[state.ui.theme].innerCardLighten)

    // STATE
    const expanded: any = ref<AavePoolModel[]>([])
    const datatable: any = ref(null)
    const cols = computed(() => {
      const cols = [
        {
          text: 'Assets',
          align: 'left',
          value: 'symbol',
          width: '230',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Your Balance',
          align: 'left',
          value: 'portfolio.walletBal',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Your Deposits',
          align: 'left',
          value: 'portfolio.totalDeposits',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Your  Borrows',
          align: 'left',
          value: 'portfolio.variableBorrow',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Asset Price',
          align: 'left',
          value: 'price.priceUsd',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: '',
          value: 'action',
          width: 300,
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Deposit APY',
          align: 'center',
          value: 'depositAPY',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Borrow APY, variable',
          align: 'center',
          value: 'variableBorrowAPY',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Total Deposited',
          align: 'center',
          value: 'usdBalance',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Total Borrowed',
          align: 'center',
          value: 'totalBorrowBalanceUsd',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: '',
          align: 'right',
          value: 'link',
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
      ]
      return props.hideActionCols ? cols.filter((item) => item.value !== 'action') : cols
    })

    const aavePools = toRefs(props).pools
    const poolsLoading = toRefs(props).loading
    const searchString = toRefs(props).search

    // METHODS
    function initAction(poolAddress: string, action: actionTypes) {
      const pool: AavePoolModel | undefined = aavePools.value.find((elem) => elem.id === poolAddress)
      if (pool) {
        emit(EmitEvents.initAction, { action, pool })
      }
    }

    async function importToMetamask(poolId: string) {
      const pool: AavePoolModel | undefined = aavePools.value.find((elem) => elem.id === poolId)
      if (pool) {
        await importTokenToMetamask({
          address: pool.underlyingAsset,
          symbol: pool.symbol,
          decimals: pool.decimals,
          image: pool.logoUrl,
        })
      }
    }
    // METHODS
    function navigateToExplorer(address: string) {
      const currentChain: Chain = getters['configs/chainInfo'](chainId.value ?? 1)
      const url = `${currentChain.blockExplorerUrl}address/${address}`
      window.open(url)
    }

    function expendSingleRow() {
      if (props.expandFirstRow && aavePools.value.length > 0) {
        const firstRow: AavePoolModel = aavePools.value[0]
        expanded.value.push(firstRow)
      }
    }

    // WATCHERS
    watch(aavePools, () => expendSingleRow())

    return {
      // DATA
      cols,
      aavePools,
      textClass,
      expanded,
      poolsLoading,
      aaveActions,
      walletReady,
      searchString,
      datatable,

      //  METHODS
      navigateToExplorer,
      initAction,
      importToMetamask,
    }
  },
})
</script>
