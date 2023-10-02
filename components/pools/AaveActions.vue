<template>
  <v-row no-gutters>
    <v-col v-if="pool">
      <v-card tile outlined :height="type === 'dialog' ? '100%' : '500'">
        <div class="overflow-y-auto" :style="type === 'dialog' ? { height: '100%' } : { height: '500px' }">
          <v-row v-if="type === 'card'" no-gutters>
            <v-col>
              <v-card-title class="font-weight-bold subtitle-1 py-3 text-capitalize">
                <v-avatar size="26" class="mr-2"><v-img :src="$imageUrlBySymbol(`aave`)"></v-img></v-avatar>AAVE V2

                <v-spacer />
                <v-btn x-small icon @click="$emit('toggle-reserve-details')">
                  <v-icon>mdi-information-outline</v-icon>
                </v-btn>
              </v-card-title>

              <v-divider />
            </v-col>
          </v-row>

          <v-card-title v-if="type === 'dialog'" class="py-3">
            <h6 class="text-h6"><span class="text-capitalize" v-text="action" /> {{ pool.symbol }}</h6>
            <v-spacer />
            <v-icon @click="toggleDialog">mdi-close</v-icon>
          </v-card-title>

          <div class="px-4 pb-4 pt-2">
            <v-row v-if="receipt">
              <v-col cols="12">
                <transaction-result
                  :receipt="receipt"
                  :is-tx-mined="isTxMined"
                  :success-message="txOptions.successMessage"
                  @on-result-closed="resetToDefault"
                />
              </v-col>
            </v-row>

            <div v-else>
              <v-row no-gutters>
                <v-col cols="12" class="text-center pb-2 mt-0">
                  <v-btn-toggle v-model="action" v color="primary" mandatory group>
                    <v-btn
                      v-for="elem in aaveActions"
                      :key="elem"
                      small
                      class="ma-0"
                      height="32"
                      :value="elem"
                      depressed
                      outlined
                    >
                      {{ elem }}
                    </v-btn>
                  </v-btn-toggle>
                </v-col>

                <v-col v-if="type === 'card'" cols="12">
                  <v-row no-gutters>
                    <v-col cols="12">
                      <v-divider />
                    </v-col>

                    <v-col v-for="(item, i) in walletSummary" :key="i" cols="4" class="pa-3 grey--text text-center">
                      <v-row>
                        <v-col cols="12" style="min-height: 118px">
                          <div class="caption font-weight-medium text-no-wrap" v-text="item.text" />
                          <div class="subtitle-1 white--text text-no-wrap" v-text="item.value" />
                          <div
                            class="subtitle-2 font-weight-regular text-no-wrap text-truncate"
                            v-text="item.priceUsd"
                          />

                          <v-chip v-if="item.apy" pill outlined class="pa-2" text-color="white" label x-small>
                            <span v-if="item.apy === -1">--.--</span>
                            <span v-else v-text="`${$f(item.apy * 100, { maxDigits: 2, after: ' %' })} APY`" />
                          </v-chip>
                        </v-col>

                        <v-divider v-if="i < walletSummary.length - 1" style="" vertical />
                      </v-row>
                    </v-col>
                    <v-col cols="12">
                      <v-divider />
                    </v-col>
                  </v-row>
                </v-col>
              </v-row>

              <aave-action-form
                :amount="amount"
                :price-usd="pool.price.priceUsd"
                :symbol="pool.symbol"
                :logo="pool.logoUrl"
                :allowance="txOptions.allowance"
                :rules="txOptions.rules"
                :disabled="!walletReady"
                @on-value-changed="onFormValueChanged"
              />

              <v-row no-gutters>
                <v-col>
                  <small v-if="type === 'dialog'">Transaction overview</small>
                  <v-card tile outlined>
                    <v-simple-table>
                      <template #default>
                        <tbody>
                          <tr
                            v-for="(item, i) in type === 'card'
                              ? txOptions.overview.filter((option) => option.isVisible)
                              : txOptions.overview"
                            :key="i"
                          >
                            <td :class="[textClass]" v-text="item.text" />
                            <td :class="[item.class]" v-text="item.value" />
                          </tr>
                        </tbody>
                      </template>
                    </v-simple-table>
                  </v-card>
                </v-col>
              </v-row>
              <v-btn
                tile
                :disabled="isActionButtonDisabled"
                class="text-capitalize mt-4"
                block
                color="primary"
                :loading="txLoading"
                @click="txOptions.txMethod"
              >
                {{ action }}
              </v-btn>
            </div>
          </div>
        </div>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  inject,
  onMounted,
  PropType,
  reactive,
  Ref,
  ref,
  toRefs,
  useContext,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { AavePoolModel, actionTypes, aaveActions, aaveVersion } from '~/composables/useAavePools'
import { State } from '~/types/state'
import AaveActionForm from '~/components/pools/AaveActionForm.vue'
import useAaveTransactions from '~/composables/useAaveTransactions'
import TransactionResult from '~/components/common/TransactionResult.vue'
import { DefiEvents, EmitEvents } from '~/types/events'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'

type Props = {
  // eslint-disable-next-line no-use-before-define
  type: aaveComponentType
  totalCollateralUsd: number
  totalBorrowedUsd: number
  healthFactor: number
  maxLtv: number
  someData: number
  poolData: AavePoolModel
  poolAction: actionTypes
  version: aaveVersion
}

interface TxOption {
  overview: { text: string; value: string | boolean; class?: string | string[]; isVisible: boolean }[]
  rules: ((v: string) => boolean | string)[]
  allowance: { text: string; value: string | number }
  successMessage?: string
  txMethod: (interestRateMode?: number) => Promise<void>
}
type ActionOptions = {
  [actionType in actionTypes]: TxOption
}

type aaveComponentType = 'dialog' | 'card'

export default defineComponent<Props>({
  components: {
    TransactionResult,
    AaveActionForm,
  },
  props: {
    type: { type: String as PropType<aaveComponentType>, required: true },
    poolData: { type: Object as PropType<AavePoolModel>, required: true },
    poolAction: { type: String as PropType<actionTypes>, required: true },
    healthFactor: { type: Number, default: 0, required: true },
    totalCollateralUsd: { type: Number, default: 0, required: true },
    totalBorrowedUsd: { type: Number, default: 0, required: true },
    maxLtv: { type: Number, default: 0, required: true },
    version: { type: String as PropType<aaveVersion>, required: true },
  },

  setup(props, { emit }) {
    // STATE
    const marketDetailDialog = ref(false)
    const dialog = ref(false)
    const amount = ref<number>(0)
    const isActionButtonDisabled = ref(true)
    const isVariableBorrow = ref(true)

    const action = ref<actionTypes>('deposit')
    const pool = ref() as Ref<AavePoolModel>
    const version = ref<aaveVersion>('v2')

    const marketStats = reactive({
      healthFactor: toRefs(props).healthFactor,
      totalCollateralUsd: toRefs(props).totalCollateralUsd,
      totalBorrowedUsd: toRefs(props).totalBorrowedUsd,
      maxLTV: toRefs(props).maxLtv,
    })

    // COMPOSABLES
    const { walletReady } = inject(WEB3_PLUGIN_KEY) as Web3
    const { $f } = useContext()
    const { state } = useStore<State>()
    const { txLoading, receipt, isTxMined, deposit, borrow, repay, withdraw, resetToDefault } = useAaveTransactions(
      pool,
      amount,
      version
    )

    // COMPUTED
    const textClass = computed(() => state.ui[state.ui.theme].innerCardLighten)
    const allowedToBorrow = computed(() => {
      const allowed =
        ((marketStats.totalCollateralUsd * marketStats.maxLTV) / 100 - marketStats.totalBorrowedUsd) /
        pool.value.price.priceUsd
      if (allowed > 0) {
        return allowed
      } else return 0
    })
    const allowedToWithdraw = computed(() => {
      if (pool.value.portfolio.totalDeposits === 0) {
        return 0
      }
      const totalAllowedToWithdraw = (allowedToBorrow.value * 100) / marketStats.maxLTV
      if (totalAllowedToWithdraw > pool.value.portfolio.totalDeposits) {
        return pool.value.portfolio.totalDeposits
      } else return totalAllowedToWithdraw
    })
    const rules = computed(() => {
      return {
        required: (v: string) => !!v || 'Value is required',
        mustBeNumber: (v: string) => !isNaN(parseFloat(v)) || 'Must be Number',
        musBeMoreThen0: (v: string) => parseFloat(v) > 0 || 'Must be more then 0',
        walletBalCheck: (v: string) =>
          (!!v && pool.value && pool.value.portfolio.walletBal >= parseFloat(v.replace(/,/g, ''))) ||
          'Insufficient Funds',
      }
    })
    const txOptions: ComputedRef<TxOption> = computed(() => {
      const stats: ActionOptions = {
        deposit: {
          txMethod: deposit,
          overview: [
            {
              text: 'Deposit APY',
              value:
                pool.value.depositAPY === -1 ? '--' : $f(pool.value.depositAPY * 100, { minDigits: 2, after: ' %' }),
              isVisible: false,
            },
            {
              text: 'Deposit APR',
              value:
                pool.value.depositAPR === -1 ? '--' : $f(pool.value.depositAPR * 100, { minDigits: 2, after: ' %' }),
              isVisible: false,
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
              isVisible: true,
            },
            {
              text: 'Used as collateral',
              class: pool.value.usageAsCollateralEnabled ? ['green--text'] : ['red--text'],
              value: pool.value.usageAsCollateralEnabled ? 'Yes' : 'No',
              isVisible: true,
            },
          ],
          rules: [
            rules.value.required,
            rules.value.mustBeNumber,
            rules.value.musBeMoreThen0,
            rules.value.walletBalCheck,
          ],
          allowance: {
            text: 'Allowed To Deposit',
            value: `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })} ${pool.value.symbol}`,
          },
          successMessage: `You Deposited ${amount.value} ${pool.value.symbol}`,
        },
        borrow: {
          txMethod: borrow,
          overview: [
            isVariableBorrow.value
              ? {
                  text: 'Variable Borrow APY',
                  value:
                    pool.value.variableBorrowAPY === -1
                      ? '--'
                      : $f(pool.value.variableBorrowAPY * 100, { minDigits: 2, after: ' %' }),
                  isVisible: false,
                }
              : {
                  text: 'Stable Borrow APY',
                  value:
                    pool.value.stableBorrowAPY === -1
                      ? '--'
                      : $f(pool.value.stableBorrowAPY * 100, { minDigits: 2, after: ' %' }),
                  isVisible: false,
                },

            isVariableBorrow.value
              ? {
                  text: 'Variable Borrow APR',
                  value:
                    pool.value.variableBorrowAPR === -1
                      ? '--'
                      : $f(pool.value.variableBorrowAPR * 100, { minDigits: 2, after: ' %' }),
                  isVisible: false,
                }
              : {
                  text: 'Stable Borrow APR',
                  value:
                    pool.value.stableBorrowAPR === -1
                      ? '--'
                      : $f(pool.value.stableBorrowAPR * 100, { minDigits: 2, after: ' %' }),
                  isVisible: false,
                },

            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
              isVisible: true,
            },
            {
              text: 'Used as collateral',
              class: pool.value.usageAsCollateralEnabled ? ['green--text'] : ['red--text'],
              value: pool.value.usageAsCollateralEnabled ? 'Yes' : 'No',
              isVisible: true,
            },
          ],
          rules: [rules.value.required, rules.value.mustBeNumber, rules.value.musBeMoreThen0],
          allowance: {
            text: 'Allowed To Borrow',
            value: `${$f(allowedToBorrow.value, { minDigits: 6 })}  ${pool.value.symbol}`,
          },
          successMessage: `You Borrowed ${amount.value} ${pool.value.symbol}`,
        },
        repay: {
          txMethod: repay,
          overview: [
            {
              text: 'Your Borrows',
              value: `${$f(pool.value.portfolio.variableBorrow, { minDigits: 6 })}  ${pool.value.symbol}`,
              isVisible: false,
            },
            {
              text: 'Your Wallet Balance',
              value: `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })}  ${pool.value.symbol}`,
              isVisible: false,
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
              isVisible: true,
            },
          ],
          rules: [rules.value.required, rules.value.mustBeNumber, rules.value.musBeMoreThen0],
          allowance: {
            text: 'Allowed to Repay',
            value: `${
              pool.value.portfolio.walletBal >= pool.value.portfolio.variableBorrow
                ? `${$f(pool.value.portfolio.variableBorrow, { minDigits: 6 })}  ${pool.value.symbol}`
                : `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })} ${pool.value.symbol}`
            }`,
          },
          successMessage: `You Repayed  ${amount.value} ${pool.value.symbol}`,
        },
        withdraw: {
          txMethod: withdraw,
          overview: [
            {
              text: 'Your Deposits',
              value: `${$f(pool.value.portfolio.totalDeposits, { minDigits: 6 })}  ${pool.value.symbol}`,
              isVisible: false,
            },
            {
              text: 'Your Wallet Balance',
              value: `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })}  ${pool.value.symbol}`,
              isVisible: false,
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
              isVisible: true,
            },
          ],
          rules: [rules.value.required, rules.value.mustBeNumber, rules.value.musBeMoreThen0],
          allowance: {
            text: 'Allowed To Withdraw',
            value: `${$f(allowedToWithdraw.value, { minDigits: 6 })}  ${pool.value.symbol}`,
          },
        },
      }
      return stats[action.value]
    })

    const walletSummary = computed<{ text: string; value: string; priceUsd: string; apy?: number }[]>(() => {
      return [
        {
          text: 'Your Deposits',
          value: $f(pool.value.portfolioVal.totalDeposits, { minDigits: 2, maxDigits: 4 }),
          priceUsd: $f(pool.value.portfolioVal.totalDeposits * pool.value.price.priceUsd, {
            pre: '$ ',
            minDigits: 2,
            maxDigits: 4,
          }),
          apy: pool.value.depositAPY,
        },
        {
          text: 'Your Balance',
          value: $f(pool.value.portfolioVal.walletBal, { minDigits: 2, maxDigits: 4 }),
          priceUsd: $f(pool.value.portfolioVal.walletBal * pool.value.price.priceUsd, {
            pre: '$ ',
            minDigits: 2,
            maxDigits: 4,
          }),
        },
        {
          text: 'Your Borrows',
          value: $f(pool.value.portfolioVal.variableBorrow, { minDigits: 2, maxDigits: 4 }),
          priceUsd: $f(pool.value.portfolioVal.variableBorrow * pool.value.price.priceUsd, {
            pre: '$ ',
            minDigits: 2,
            maxDigits: 4,
          }),
          apy: pool.value.variableBorrowAPY,
        },
      ]
    })

    function onFormValueChanged({ isFormValid, amountVal }: { isFormValid: boolean; amountVal: number }) {
      isActionButtonDisabled.value = !isFormValid
      amount.value = amountVal
    }

    function toggleDialog() {
      dialog.value = false
      emit(DefiEvents.toggleActionDialog)
    }

    watch(action, (newVal: actionTypes) => (newVal ? emit(EmitEvents.onValueChanged, newVal) : null))
    watch(isTxMined, (newVal: boolean) => (newVal ? emit(EmitEvents.transactionSuccess, newVal) : null))
    watch(
      () => props.poolData,
      (newVal: AavePoolModel) => (pool.value = newVal)
    )

    onMounted(() => {
      action.value = props.poolAction
      pool.value = props.poolData
      version.value = props.version
    })

    return {
      dialog,
      marketDetailDialog,
      action,
      isVariableBorrow,
      pool,
      txOptions,
      textClass,
      marketStats,
      amount,
      rules,
      isActionButtonDisabled,
      txLoading,
      receipt,
      isTxMined,
      aaveActions,
      walletSummary,
      walletReady,

      // METHODS
      onFormValueChanged,
      deposit,
      borrow,
      resetToDefault,
      toggleDialog,
    }
  },
})
</script>
