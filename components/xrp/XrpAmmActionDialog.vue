<template>
  <v-dialog v-model="dialog" width="500" max-width="500">
    <v-card v-if="pool && dialog" tile outlined max-width="500" class="pa-4" height="100%">
      <div>  
        <v-row no-gutters align="center">
          <v-col cols="11">
            <h6 class="text-h6">
              <span class="text-capitalize" v-text="action" /> 
              {{ pool.token1.symbol }}/{{ pool.token2.symbol }}
            </h6>
          </v-col>
          <v-spacer />
          <v-col>
            <div class="text-right">
              <v-icon @click="dialog = false">mdi-close</v-icon>
            </div>
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
                <v-btn-toggle v-model="action" color="primary" mandatory group>
                  <v-btn
                    v-for="elem in ammActions"
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

          <xrp-amm-action-form
            :action="action"
            :pool="pool"
            :amount="amount"
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
  defineComponent,
  reactive,
  ref,
  toRefs,
  useContext,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import useXrpAmmTransactions from '~/composables/useXrpAmmTransactions'
import XrpTransactionResult from '~/components/common/XrpTransactionResult.vue'
import XrpAmmActionForm from '~/components/xrp/XrpAmmActionForm.vue'

interface XrpAmmPool {
  id: string
  token1: { symbol: string; name: string; icon: string }
  token2: { symbol: string; name: string; icon: string; issuer?: string }
  liquidity: number
  volume24h: number
  fee: number
  apr: number
  priceChange24h: number
  token1Balance: number
  token2Balance: number
}

type ActionTypes = 'deposit' | 'withdraw'
type ActionOptions = {
  [actionType in ActionTypes]: {
    overview: { text: string; value: string | boolean; class?: string | string[] }[]
    rules: ((v: string) => boolean | string)[]
    successMessage?: string
    txMethod: () => Promise<void>
  }
}

export default defineComponent({
  components: {
    XrpAmmActionForm,
    XrpTransactionResult,
  },
  props: {
    pool: { type: Object as () => XrpAmmPool, required: true },
  },
  setup(props, { emit }) {
    // STATE
    const dialog = ref(false)
    const action = ref<ActionTypes>('deposit')
    const amount = ref<number>(0)
    const isActionButtonDisabled = ref(true)

    const ammActions: ActionTypes[] = ['deposit', 'withdraw']

    // COMPOSABLES
    const { $f } = useContext()
    const { state } = useStore<State>()
    const { txLoading, receipt, isTxMined, deposit, withdraw, resetToDefault } = useXrpAmmTransactions(
      props.pool,
      amount
    )

    // COMPUTED
    const textClass = computed(() => state.ui[state.ui.theme].innerCardLighten)
    
    const txOptions = computed<ActionOptions>(() => ({
      deposit: {
        overview: [
          { text: 'Pool', value: `${props.pool.token1.symbol}/${props.pool.token2.symbol}` },
          { text: 'Amount', value: `${$f(amount.value, { minDigits: 6, after: '' })} ${props.pool.token1.symbol}` },
          { text: 'Estimated Pool Tokens', value: `${$f(amount.value * 0.95, { minDigits: 6, after: '' })} LP` },
          { text: 'Fee', value: `${$f(props.pool.fee * 100, { minDigits: 3, after: '' })}%` },
        ],
        rules: [
          (v: string) => !!v || 'Amount is required',
          (v: string) => !isNaN(parseFloat(v)) || 'Must be a number',
          (v: string) => parseFloat(v) > 0 || 'Must be greater than 0',
        ],
        successMessage: 'Successfully deposited to pool',
        txMethod: deposit,
      },
      withdraw: {
        overview: [
          { text: 'Pool', value: `${props.pool.token1.symbol}/${props.pool.token2.symbol}` },
          { text: 'Pool Tokens', value: `${$f(amount.value, { minDigits: 6, after: '' })} LP` },
          { text: 'Estimated Return', value: `${$f(amount.value * 0.95, { minDigits: 6, after: '' })} ${props.pool.token1.symbol}` },
          { text: 'Fee', value: `${$f(props.pool.fee * 100, { minDigits: 3, after: '' })}%` },
        ],
        rules: [
          (v: string) => !!v || 'Amount is required',
          (v: string) => !isNaN(parseFloat(v)) || 'Must be a number',
          (v: string) => parseFloat(v) > 0 || 'Must be greater than 0',
        ],
        successMessage: 'Successfully withdrawn from pool',
        txMethod: withdraw,
      },
    }))

    // METHODS
    const onFormValueChanged = (value: number) => {
      amount.value = value
      isActionButtonDisabled.value = value <= 0
    }

    const openDialog = () => {
      dialog.value = true
      resetToDefault()
    }

    const closeDialog = () => {
      dialog.value = false
      resetToDefault()
    }

    // Watch for action changes
    watch(action, () => {
      amount.value = 0
      isActionButtonDisabled.value = true
    })

    return {
      // State
      dialog,
      action,
      amount,
      isActionButtonDisabled,
      ammActions,
      
      // Computed
      textClass,
      txOptions,
      
      // Methods
      onFormValueChanged,
      openDialog,
      closeDialog,
      
      // Transaction state
      txLoading,
      receipt,
      isTxMined,
    }
  },
})
</script> 