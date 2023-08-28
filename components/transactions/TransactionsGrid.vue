<template>
  <v-card tile outlined>
    <v-data-table
      :items-per-page="itemsPerPage"
      :headers="cols"
      :items="transactions"
      item-key="txHash"
      class="elevation-0"
      show-expand
      single-expand
      hide-default-footer
      :mobile-breakpoint="0"
    >
      <!--  expanded buttons   -->
      <template #item.data-table-expand="{ item, isExpanded, expand }">
        <div style="width: 16px">
          <v-btn
            v-if="item.logEvents.length > 0 && item.methodTextRenderer !== 'Transfer'"
            small
            icon
            @click="expand(!isExpanded)"
          >
            <v-icon>mdi-chevron-{{ isExpanded ? 'up' : 'down' }}</v-icon>
          </v-btn>
        </div>
      </template>

      <!-- expanded column items -->
      <template #expanded-item="{ headers, item }">
        <td :colspan="headers.length">
          <log-events
            :key="item.txHash"
            :wallet-address="account"
            :log-events="item.logEvents"
            @navigate-to-explorer="$navigateToExplorer"
          />
        </td>
      </template>

      <!--   Txn Date   -->
      <template #[`item.blockSignedAt`]="{ item }">
        <span class="text-no-wrap grey--text" v-text="item.txDate"></span>
      </template>

      <!--   Method   -->
      <template #[`item.method`]="{ item }">
        <div class="text-truncate">
          <v-tooltip right color="black">
            <template #activator="{ on, attrs }">
              <v-btn
                outlined
                small
                width="120"
                depressed
                rounded
                class="cursor-text"
                color="primary"
                v-bind="attrs"
                v-on="on"
              >
                <span
                  class="text-capitalize px-2 font-weight-regular white--text text-truncate"
                  style="max-width: 120px"
                  v-text="item.methodTextRenderer"
                />
              </v-btn>
            </template>
            <span class="text-capitalize" v-text="item.methodTextRenderer" />
          </v-tooltip>
        </div>
      </template>

      <!--   Txn Hash   -->
      <template #[`item.txHash`]="{ item }">
        <nuxt-link class="text-decoration-none grey--text" :to="`/tx/${item.txHash}?chainId=${chainId}`">{{
          $truncateAddress(item.txHash, 6, 6)
        }}</nuxt-link>
      </template>

      <!--   isInbound   -->
      <template #[`item.isInbound`]="{ item }">
        <div class="text-no-wrap">
          <v-btn
            style="pointer-events: none"
            width="50"
            :color="isInboundRenderer(item).color"
            small
            class="rounded-sm"
            text
          >
            {{ isInboundRenderer(item).text }}
          </v-btn>
        </div>
      </template>

      <!--   from and to   -->
      <template #[`item.fromTo`]="{ item }">
        <div class="text-no-wrap">
          <template v-if="item.txDetail">
            <tx-address-label
              label="From:"
              :wallet-address="account.toLowerCase()"
              :address="item.txDetail.fromAddress.toLowerCase()"
              :navigate-to-explorer="$navigateToExplorer"
            />

            <tx-address-label
              label="To:"
              :wallet-address="account.toLowerCase()"
              :address="item.txDetail.toAddress.toLowerCase()"
              :navigate-to-explorer="$navigateToExplorer"
            />
          </template>

          <template v-else>
            <tx-address-label
              label="From:"
              :wallet-address="account.toLowerCase()"
              :address="item.fromAddress.toString().toLowerCase()"
              :symbol="item.fromAddressSymbol.toLowerCase()"
              :name="item.fromAddressName"
              :is-contract="item.fromAddressIsContract"
              :navigate-to-explorer="$navigateToExplorer"
            />

            <tx-address-label
              label="To:"
              :wallet-address="account.toLowerCase()"
              :address="item.toAddress.toString().toLowerCase()"
              :symbol="item.toAddressSymbol.toLowerCase()"
              :name="item.toAddressName"
              :is-contract="item.toAddressIsContract"
              :navigate-to-explorer="$navigateToExplorer"
            />
          </template>
        </div>
      </template>

      <!--   value   -->
      <template #[`item.value`]="{ item }">
        <div
          class="text-no-wrap grey--text"
          v-text="
            `${$f(item.txnValue.amount, { maxDigits: 2 })}
              ${item.txDetail ? item.txDetail.tokenContractSymbol : chainSymbol}`
          "
        />
        <div class="text-no-wrap" v-text="$f(item.txnValue.priceUsd, { minDigits: 2, pre: '$ ' })" />
      </template>

      <!--   txn fee   -->
      <template #[`item.txnFee`]="{ item }">
        <div class="text-no-wrap grey--text" v-text="`${$f(item.txnFee, { minDigits: 2 })} ${chainSymbol}`" />
        <div class="text-no-wrap" v-text="$f(item.gasQuote, { minDigits: 2, pre: '$ ' })" />
      </template>

      <!--   status   -->
      <template #[`item.isSuccess`]="{ item }">
        <v-btn
          width="90"
          outlined
          small
          depressed
          rounded
          class="text-capitalize cursor-text caption"
          :color="item.isSuccess.color"
        >
          <span class="white--text" v-text="item.isSuccess.text" />
        </v-btn>
      </template>

      <template #[`item.action`]="{ item }">
        <v-tooltip bottom color="black">
          <template #activator="{ on, attrs }">
            <v-btn
              color="primary"
              v-bind="attrs"
              x-small
              class="mb-1 ml-1"
              icon
              v-on="on"
              @click="$copyAddressToClipboard(item.txHash)"
            >
              <v-icon size="18">mdi-content-copy</v-icon>
            </v-btn>
          </template>
          Copy Transaction Hash
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
              @click="$navigateToExplorer(item.txHash, 'tx', currentChain.blockExplorerUrl)"
            >
              <v-icon size="18">mdi-open-in-new</v-icon>
            </v-btn>
          </template>
          Open in Block Explorer
        </v-tooltip>
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, inject, PropType, ref, useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import LogEvents from '~/components/transactions/LogEvents.vue'
import TxAddressLabel from '~/components/transactions/TxAddressLabel.vue'
import { TransactionModel } from '~/composables/useTransactions'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
type inboundFunctionType = (item: TransactionModel) => { color: string; text: string }
export default defineComponent({
  name: 'TransactionsGrid',
  components: { TxAddressLabel, LogEvents },
  props: {
    isInboundRenderer: { type: Function as PropType<inboundFunctionType>, required: true },
    transactions: { type: Array as PropType<TransactionModel[]>, default: () => [] },
    itemsPerPage: { type: Number, default: 15 },
    chainId: { type: Number, default: 1 },
    account: { type: String, default: '' },
  },
  setup() {
    // STATE
    const cols = ref([
      {
        text: '',
        value: 'data-table-expand',
      },
      {
        text: 'Txn Date',
        align: 'start',
        value: 'blockSignedAt',
        class: 'py-2',
        cellClass: 'py-4',
      },
      {
        text: 'Method',
        align: 'start',
        value: 'method',
        class: 'py-2',
        sortable: false,
      },
      {
        text: 'Txn Hash',
        align: 'start',
        value: 'txHash',
        class: 'py-2',
        sortable: false,
      },
      {
        text: '',
        align: 'start',
        value: 'isInbound',
        class: 'py-2',
        sortable: false,
      },
      {
        text: '',
        align: 'start',
        value: 'fromTo',
        class: 'py-2',
        sortable: false,
      },
      {
        text: 'Value',
        align: 'start',
        value: 'value',
        class: 'py-2',
      },
      {
        text: 'Txn Fee',
        align: 'start',
        value: 'txnFee',
        class: 'py-2',
      },
      {
        text: 'Status',
        align: 'start',
        value: 'isSuccess',
        class: 'py-2',
      },

      {
        text: '',
        sortable: false,
        align: 'start',
        value: 'action',
        class: 'py-2',
        width: '100px',
      },
    ])

    const { getters } = useStore<State>()
    const { chainId } = inject(WEB3_PLUGIN_KEY) as Web3

    // COMPUTED
    const currentChain = computed(() => getters['configs/chainInfo'](chainId.value ?? 1))
    const chainSymbol = computed(() => currentChain.value.symbol)

    return {
      cols,
      chainSymbol,
      currentChain,
    }
  },
})
</script>
