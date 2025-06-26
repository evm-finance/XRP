<template>
  <v-dialog v-model="dialog" width="500" max-width="500">
    <v-card tile outlined max-width="500" class="pa-4" height="100%">
      <div>
        <v-row no-gutters align="center">
          <v-col cols="11">
            <h6 class="text-h6">Swap {{ pool.token1.symbol }}/{{ pool.token2.symbol }}</h6>
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
          success-message="Swap completed successfully"
          @on-result-closed="dialog = false"
        />
        
        <div v-else>
          <v-row no-gutters class="pa-2">
            <v-col cols="12">
              <v-row no-gutters>
                <v-col class="px-3">
                  <token-input-field
                    :token="fromToken"
                    :balance="fromTokenBalance"
                    :amount="fromAmount"
                    :loading="loading"
                    @on-value-changed="onFromAmountChange"
                  />
                </v-col>
              </v-row>
            </v-col>

            <v-btn class="mt-n1 mb-3" text outlined fab small @click="swapTokens">
              <v-icon size="28">mdi-swap-horizontal-variant</v-icon>
            </v-btn>

            <v-col cols="12">
              <v-row no-gutters>
                <v-col class="px-3">
                  <token-input-field
                    :token="toToken"
                    :balance="toTokenBalance"
                    :amount="toAmount"
                    :loading="loading"
                    :readonly="true"
                  />
                </v-col>
              </v-row>
            </v-col>

            <v-col v-if="errorMessage" cols="12">
              <v-alert text type="error" dismissible class="ma-2">{{ errorMessage }}</v-alert>
            </v-col>

            <v-col v-else cols="12">
              <v-row class="pa-3 caption font-weight-medium">
                <v-col cols="12">
                  <v-card tile outlined class="mt-2">
                    <div v-if="loading" class="d-flex pa-2" style="height: 36px">
                      <span class="grey--text">
                        <v-progress-circular size="20" indeterminate color="pink" class="mr-1" /> 
                        Fetching Latest Prices
                      </span>
                    </div>

                    <div v-else class="d-flex pa-2" style="height: 36px">
                      <span v-text="quote">
                        <span class="ml-2 grey--text">($ 1.0000)</span>
                      </span>
                      <v-spacer />
                      <div v-if="!expand" class="pr-2">
                        <v-icon color="grey lighten-1" size="17">mdi-gas-station</v-icon>
                        <span class="grey--text text--lighten-1" v-text="$f(gasFeeUSD, { maxDigits: 2, pre: '$' })" />
                      </div>
                      <v-btn color="grey lighten-1" height="22" width="22" icon @click="expand = !expand">
                        <v-icon size="22">mdi-chevron-{{ expand ? 'up' : 'down' }}</v-icon>
                      </v-btn>
                    </div>

                    <v-simple-table v-if="expand && !loading" dense>
                      <template #default>
                        <tbody>
                          <tr v-for="(item, i) in details" :key="i">
                            <td class="caption grey--text text--lighten-1" v-text="item.text" />
                            <td class="caption text-right" v-text="item.value" />
                          </tr>
                        </tbody>
                      </template>
                    </v-simple-table>
                  </v-card>
                </v-col>
              </v-row>
            </v-col>

            <v-col cols="12" class="pa-3 pt-2">
              <v-btn
                tile
                block
                color="primary"
                :disabled="!actionButton.status || txLoading"
                :loading="txLoading"
                @click="swap"
              >
                {{ actionButton.message }}
              </v-btn>
            </v-col>
          </v-row>
        </div>
      </div>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext, watch } from '@nuxtjs/composition-api'
import useXrpAmmSwap from '~/composables/useXrpAmmSwap'
import TransactionResult from '~/components/common/TransactionResult.vue'
import TokenInputField from '~/components/trading/TokenInputField.vue'

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

export default defineComponent({
  components: {
    TransactionResult,
    TokenInputField,
  },
  props: {
    pool: { type: Object as () => XrpAmmPool, required: true },
  },
  setup(props) {
    const { $f } = useContext()
    
    // State
    const dialog = ref(false)
    const expand = ref(false)
    const fromAmount = ref(0)
    const toAmount = ref(0)
    const fromToken = ref(props.pool.token1)
    const toToken = ref(props.pool.token2)
    
    // Swap composable
    const {
      fromTokenBalance,
      toTokenBalance,
      loading,
      quote,
      gasFeeUSD,
      actionButton,
      errorMessage,
      txLoading,
      receipt,
      isTxMined,
      swap: executeSwap,
      clearTrade,
    } = useXrpAmmSwap(fromToken, toToken, fromAmount, props.pool)
    
    // Computed
    const details = computed(() => [
      { text: 'Expected Output', value: `${toAmount.value.toFixed(6)} ${toToken.value.symbol}` },
      { text: 'Price Impact', value: '< 0.01%' },
      { text: 'Pool Fee', value: `${(props.pool.fee * 100).toFixed(3)}%` },
      { text: 'Minimum Received', value: `${(toAmount.value * 0.99).toFixed(6)} ${toToken.value.symbol}` },
    ])
    
    // Methods
    const onFromAmountChange = (value: number) => {
      fromAmount.value = value
      // Calculate expected output (simplified)
      toAmount.value = value * 0.997 // Account for 0.3% fee
    }
    
    const swapTokens = () => {
      const temp = fromToken.value
      fromToken.value = toToken.value
      toToken.value = temp
      
      const tempAmount = fromAmount.value
      fromAmount.value = toAmount.value
      toAmount.value = tempAmount
    }
    
    const swap = async () => {
      await executeSwap()
    }
    
    const openDialog = () => {
      dialog.value = true
      clearTrade()
    }
    
    const closeDialog = () => {
      dialog.value = false
      clearTrade()
    }
    
    // Watch for pool changes
    watch(() => props.pool, (newPool) => {
      fromToken.value = newPool.token1
      toToken.value = newPool.token2
      fromAmount.value = 0
      toAmount.value = 0
    })
    
    return {
      // State
      dialog,
      expand,
      fromAmount,
      toAmount,
      fromToken,
      toToken,
      
      // Computed
      details,
      
      // Methods
      onFromAmountChange,
      swapTokens,
      swap,
      openDialog,
      closeDialog,
      
      // Swap state
      fromTokenBalance,
      toTokenBalance,
      loading,
      quote,
      gasFeeUSD,
      actionButton,
      errorMessage,
      txLoading,
      receipt,
      isTxMined,
    }
  },
})
</script> 