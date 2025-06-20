<template>
  <v-card outlined tile :disabled="!isWalletReady" :height="height" :width="width">
    <div
      v-if="!receipt"
      :class="height === '100%' ? '' : 'overflow-y-auto'"
      :style="height === '100%' ? '' : { height: `${height}px` }"
    >
      <v-card-title class="subtitle-1 font-weight-medium py-3">
        <v-avatar size="26" class="mr-2"><v-img :src="$imageUrlBySymbol(`xrp`)"></v-img></v-avatar>XRP AMM
        <v-spacer />
      </v-card-title>
      <v-divider class="mb-4" />
      <v-row no-gutters justify="center" class="pa-2">
        <v-col cols="12">
          <v-row no-gutters>
            <v-col class="px-3">
              <xrp-token-input-field
                form-trade-direction="EXACT_INPUT"
                :token="fromToken"
                :balance="fromTokenBalance"
                :fiat-price="fromTokenFiatPrice"
                :expected-convert-quote="expectedConvertQuote"
                :loading="loading"
                @on-value-changed="onAmountChange"
                @on-token-menu-open="onToggleTokenMenu('EXACT_INPUT')"
              />
            </v-col>
          </v-row>
        </v-col>

        <v-btn class="mt-n1 mb-3" text outlined fab small @click="onTokenChange">
          <v-icon size="28"> mdi-swap-horizontal-variant</v-icon>
        </v-btn>

        <v-col cols="12">
          <v-row no-gutters>
            <v-col class="px-3">
              <xrp-token-input-field
                form-trade-direction="EXACT_OUTPUT"
                :token="toToken"
                :balance="toTokenBalance"
                :fiat-price="toTokenFiatPrice"
                :expected-convert-quote="expectedConvertQuote"
                :loading="loading"
                @on-value-changed="onAmountChange"
                @on-token-menu-open="onToggleTokenMenu('EXACT_OUTPUT')"
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
              <v-card v-if="expectedConvertQuote > 0" tile outlined class="mt-2">
                <div v-if="loading" class="d-flex pa-2" style="height: 36px">
                  <span class="grey--text">
                    <v-progress-circular size="20" indeterminate color="pink" class="mr-1" /> Calculating Swap
                  </span>
                </div>

                <div v-else class="d-flex pa-2" style="height: 36px">
                  <span>Expected Output: {{ expectedConvertQuote.toFixed(6) }} {{ toToken.symbol }}</span>
                  <v-spacer />
                  <v-btn color="grey lighten-1" height="22" width="22" icon @click="expand = !expand">
                    <v-icon size="22">mdi-chevron-{{ expand ? 'up' : 'down' }}</v-icon>
                  </v-btn>
                </div>

                <v-simple-table v-if="expand && !loading" dense>
                  <template #default>
                    <tbody>
                      <tr>
                        <td class="caption grey--text text--lighten-1">Expected Output</td>
                        <td class="caption text-right">{{ expectedConvertQuote.toFixed(6) }} {{ toToken.symbol }}</td>
                      </tr>
                      <tr>
                        <td class="caption grey--text text--lighten-1">Slippage</td>
                        <td class="caption text-right">{{ slippage }}%</td>
                      </tr>
                      <tr>
                        <td class="caption grey--text text--lighten-1">Network Fee</td>
                        <td class="caption text-right">~0.00001 XRP</td>
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

        <xrp-token-menu-dialog
          ref="tokenMenuDialogComponent"
          @on-token-select="onTokenSelect"
        />
      </v-row>
    </div>
    <v-row v-if="receipt" class="pa-2">
      <v-col cols="12">
        <transaction-result
          :receipt="receipt"
          :is-tx-mined="isTxMined"
          success-message="Swap completed successfully!"
          :scroll-height="logScrollHeight"
          @on-result-closed="resetTransaction"
        />
      </v-col>
    </v-row>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, onMounted, watch } from '@nuxtjs/composition-api'
import { inject } from '@nuxtjs/composition-api'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import { XRPToken } from '~/composables/useXrpAmm'
import useXrpAmm from '~/composables/useXrpAmm'
import TransactionResult from '~/components/common/TransactionResult.vue'

const defaultXRP: XRPToken = {
  currency: 'XRP',
  symbol: 'XRP',
  name: 'XRP',
  decimals: 6
}

const defaultUSDC: XRPToken = {
  currency: 'USDC',
  issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
  symbol: 'USDC',
  name: 'USD Coin',
  decimals: 6
}

type Props = {
  height: string
  width: string
  inToken: XRPToken
  outToken: XRPToken
}

export default defineComponent<Props>({
  components: { TransactionResult },

  props: {
    height: { type: String, default: '100%' },
    width: { type: String, default: '450' },
    logScrollHeight: { type: String, default: '140' },
    inToken: { type: Object as PropType<XRPToken>, default: () => defaultXRP },
    outToken: { type: Object as PropType<XRPToken>, default: () => defaultUSDC },
  },

  setup(props) {
    const { isWalletReady } = inject(XRP_PLUGIN_KEY) as XrpClient
    
    const expand = ref(false)
    const tokenMenuDialogComponent = ref<any>(null)
    const fromToken = ref(props.inToken)
    const toToken = ref(props.outToken)
    const amount = ref(0)
    const slippage = ref(0.5)

    const {
      loading,
      errorMessage,
      txLoading,
      receipt,
      isTxMined,
      expectedConvertQuote,
      fromTokenBalance,
      toTokenBalance,
      fromTokenFiatPrice,
      toTokenFiatPrice,
      actionButton,
      initialize,
      performSwap,
      resetTransaction,
      watchAndCalculate
    } = useXrpAmm(fromToken, toToken, amount, slippage)

    const onToggleTokenMenu = (type: string): void => {
      tokenMenuDialogComponent.value.dialog = true
    }

    const onTokenSelect = (token: XRPToken) => {
      // For now, just swap the tokens
      const tempFromToken = toToken.value
      const tempToToken = fromToken.value
      fromToken.value = tempFromToken
      toToken.value = tempToToken
    }

    const onAmountChange = ({ value }: { value: number }) => {
      amount.value = value
    }

    const onTokenChange = () => {
      const tempFromToken = toToken.value
      const tempToToken = fromToken.value
      fromToken.value = tempFromToken
      toToken.value = tempToToken
    }

    const swap = async () => {
      const result = await performSwap()
      if (result.success) {
        console.log('Swap successful:', result)
      } else {
        console.error('Swap failed:', result.error)
      }
    }

    // Watch for amount changes to recalculate quote
    watch(amount, () => {
      watchAndCalculate()
    })

    // Initialize on mount
    onMounted(() => {
      initialize()
    })

    return {
      expand,
      tokenMenuDialogComponent,
      fromToken,
      toToken,
      slippage,
      isWalletReady,

      // From composable
      loading,
      errorMessage,
      txLoading,
      receipt,
      isTxMined,
      expectedConvertQuote,
      fromTokenBalance,
      toTokenBalance,
      fromTokenFiatPrice,
      toTokenFiatPrice,
      actionButton,

      onToggleTokenMenu,
      onTokenSelect,
      onAmountChange,
      onTokenChange,
      swap,
      resetTransaction,
    }
  },
})
</script>

<style scoped></style> 