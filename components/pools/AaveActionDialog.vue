<template>
  <v-dialog v-model="dialog" width="400" max-width="400">
    <v-card v-if="pool && dialog" tile outlined max-width="400" class="pa-4" height="100%">
      <div>  
        <v-row no-gutters align="center">
          <v-col cols="11">
            <h6 class="text-h6"><span class="text-capitalize" v-text="action" /> {{ pool.symbol }}</h6>
          </v-col>
          <v-spacer />
          <v-col>
            <div class="text-right"><v-icon @click="dialog = false">mdi-close</v-icon></div>
          </v-col>
        </v-row>
        <transaction-result
          v-if="receipt"
          :receipt="receipt"
          :is-tx-mined="isTxMined"
          :success-message="txOptions.successMessage"
          @on-result-closed="dialog = false"
        />
        <div v-else>
          <v-row no-gutters class="mb-8">
            <v-col cols="12" class="pt-2">
              <div class="text-center">
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
              </div>
            </v-col>
          </v-row>

          <!--          <v-row>-->
          <!--            <v-col class="py-0">-->
          <!--              <a v-if="isChainAndMarketMismatched" href="#" class="text-decoration-none" @click="changeToRequiredChain">-->
          <!--                <small class="grey&#45;&#45;text">-->
          <!--                  Please Switch To-->
          <!--                  <span-->
          <!--                    class="red&#45;&#45;text text&#45;&#45;lighten-1 mr-4 font-weight-bold"-->
          <!--                    v-text="isChainAndMarketMismatched.label"-->
          <!--                  />-->
          <!--                </small>-->
          <!--              </a>-->
          <!--            </v-col>-->
          <!--          </v-row>-->
          <aave-action-form
            :amount="amount"
            :price-usd="pool.price.priceUsd"
            :symbol="pool.symbol"
            :logo="pool.logoUrl"
            :allowance="txOptions.allowance"
            :rules="txOptions.rules"
            @on-value-changed="onFormValueChanged"
          />

          <v-row no-gutters>
            <v-col>
              <small>Transaction overview</small>
              <v-card tile outlined>
                <v-simple-table>
                  <template #default>
                    <tbody>
                      <tr v-for="(item, i) in txOptions.overview" :key="i">
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
            class="text-capitalize my-4"
            block
            color="primary"
            :loading="txLoading"
            @click="txOptions.txMethod"
          >
            {{ action }}
          </v-btn>
        </div>
      </div>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  reactive,
  Ref,
  ref,
  toRefs,
  useContext,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { AavePoolModel, actionTypes, aaveActions } from '~/composables/useAavePools'
import { State } from '~/types/state'
import useAaveTransactions from '~/composables/useAaveTransactions'
import TransactionResult from '~/components/common/TransactionResult.vue'
import { EmitEvents } from '~/types/events'
import AaveActionForm from '~/components/pools/AaveActionForm.vue'

type Props = {
  totalCollateralUsd: number
  totalBorrowedUsd: number
  healthFactor: number
  maxLtv: number
  someData: number
}

interface TxOption {
  overview: { text: string; value: string | boolean; class?: string | string[] }[]
  rules: ((v: string) => boolean | string)[]
  allowance: { text: string; value: string | number }
  successMessage?: string
  txMethod: (interestRateMode?: number) => Promise<void>
}
type ActionOptions = {
  [actionType in actionTypes]: TxOption
}

export default defineComponent<Props>({
  components: {
    AaveActionForm,
    TransactionResult,
  },
  props: {
    healthFactor: { type: Number, default: 0, required: true },
    totalCollateralUsd: { type: Number, default: 0, required: true },
    totalBorrowedUsd: { type: Number, default: 0, required: true },
    maxLtv: { type: Number, default: 0, required: true },
  },

  setup(props, { emit }) {
    // STATE
    const dialog = ref(false)
    const action = ref<actionTypes>('deposit')
    const pool = ref() as Ref<AavePoolModel>
    const amount = ref<number>(0)
    const isActionButtonDisabled = ref(true)
    const isVariableBorrow = ref(true)

    const marketStats = reactive({
      healthFactor: toRefs(props).healthFactor,
      totalCollateralUsd: toRefs(props).totalCollateralUsd,
      totalBorrowedUsd: toRefs(props).totalBorrowedUsd,
      maxLTV: toRefs(props).maxLtv,
    })

    // COMPOSABLES
    const { $f } = useContext()
    const { state } = useStore<State>()
    const { txLoading, receipt, isTxMined, deposit, borrow, repay, withdraw, resetToDefault } = useAaveTransactions(
      pool,
      amount
    )
    // const { isChainAndMarketMismatched, changeToRequiredChain } = useMarketSelector()

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
            },
            {
              text: 'Deposit APR',
              value:
                pool.value.depositAPR === -1 ? '--' : $f(pool.value.depositAPR * 100, { minDigits: 2, after: ' %' }),
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
            },
            {
              text: 'Used as collateral',
              class: pool.value.usageAsCollateralEnabled ? ['green--text'] : ['red--text'],
              value: pool.value.usageAsCollateralEnabled ? 'Yes' : 'No',
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
                }
              : {
                  text: 'Stable Borrow APY',
                  value:
                    pool.value.stableBorrowAPY === -1
                      ? '--'
                      : $f(pool.value.stableBorrowAPY * 100, { minDigits: 2, after: ' %' }),
                },

            isVariableBorrow.value
              ? {
                  text: 'Variable Borrow APR',
                  value:
                    pool.value.variableBorrowAPR === -1
                      ? '--'
                      : $f(pool.value.variableBorrowAPR * 100, { minDigits: 2, after: ' %' }),
                }
              : {
                  text: 'Stable Borrow APR',
                  value:
                    pool.value.stableBorrowAPR === -1
                      ? '--'
                      : $f(pool.value.stableBorrowAPR * 100, { minDigits: 2, after: ' %' }),
                },

            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
            },
            {
              text: 'Used as collateral',
              class: pool.value.usageAsCollateralEnabled ? ['green--text'] : ['red--text'],
              value: pool.value.usageAsCollateralEnabled ? 'Yes' : 'No',
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
            },
            {
              text: 'Your Wallet Balance',
              value: `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })}  ${pool.value.symbol}`,
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
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
            },
            {
              text: 'Your Wallet Balance',
              value: `${$f(pool.value.portfolio.walletBal, { minDigits: 6 })}  ${pool.value.symbol}`,
            },
            {
              text: 'Health Factor',
              value: isNaN(marketStats.healthFactor) ? '--' : $f(marketStats.healthFactor, { minDigits: 2 }),
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

    function init(currentAction: actionTypes, currentPool: AavePoolModel) {
      pool.value = currentPool
      action.value = currentAction
      dialog.value = true
    }

    function onFormValueChanged({ isFormValid, amountVal }: { isFormValid: boolean; amountVal: number }) {
      isActionButtonDisabled.value = !isFormValid
      amount.value = amountVal
    }

    // Watch if dialog closed reset tx details
    watch(dialog, (newVal: boolean) => (!newVal ? resetToDefault() : null))
    watch(isTxMined, (newVal: boolean) => (newVal ? emit(EmitEvents.transactionSuccess, newVal) : null))

    return {
      dialog,
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

      // METHODS
      init,
      onFormValueChanged,
      deposit,
      borrow,
    }
  },
})
</script>
