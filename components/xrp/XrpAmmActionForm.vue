<template>
  <v-row no-gutters>
    <v-col cols="12">
      <v-text-field
        v-model="amountInput"
        :label="`${action === 'deposit' ? 'Deposit' : 'Withdraw'} Amount`"
        type="number"
        outlined
        dense
        :rules="rules"
        @input="onAmountChange"
      />
    </v-col>
    
    <v-col cols="12" class="mt-2">
      <div class="d-flex justify-space-between align-center">
        <span class="text-caption grey--text">
          {{ action === 'deposit' ? 'Available' : 'Pool Tokens' }}: 
          {{ formatBalance(availableBalance) }}
        </span>
        <v-btn
          text
          small
          color="primary"
          @click="setMaxAmount"
        >
          MAX
        </v-btn>
      </div>
    </v-col>
    
    <v-col v-if="action === 'deposit'" cols="12" class="mt-2">
      <div class="text-caption grey--text">
        You'll receive approximately {{ formatBalance(estimatedPoolTokens) }} pool tokens
      </div>
    </v-col>
    
    <v-col v-if="action === 'withdraw'" cols="12" class="mt-2">
      <div class="text-caption grey--text">
        You'll receive approximately {{ formatBalance(estimatedReturn) }} {{ pool.token1.symbol }}
      </div>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, ref, watch } from '@nuxtjs/composition-api'

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
  props: {
    action: { type: String, required: true },
    pool: { type: Object as () => XrpAmmPool, required: true },
    amount: { type: Number, default: 0 },
    rules: { type: Array, default: () => [] },
  },
  setup(props, { emit }) {
    const amountInput = ref('')
    
    // Mock user balances - in real implementation this would come from wallet
    const userToken1Balance = ref(1000) // Mock XRP balance
    const userToken2Balance = ref(500)  // Mock token2 balance
    const userPoolTokens = ref(100)     // Mock pool tokens
    
    const availableBalance = computed(() => {
      if (props.action === 'deposit') {
        return userToken1Balance.value
      } else {
        return userPoolTokens.value
      }
    })
    
    const estimatedPoolTokens = computed(() => {
      if (props.action === 'deposit' && props.amount > 0) {
        // Simple estimation - in real implementation this would use AMM math
        return props.amount * 0.95
      }
      return 0
    })
    
    const estimatedReturn = computed(() => {
      if (props.action === 'withdraw' && props.amount > 0) {
        // Simple estimation - in real implementation this would use AMM math
        return props.amount * 0.95
      }
      return 0
    })
    
    const formatBalance = (value: number): string => {
      return value.toFixed(6)
    }
    
    const onAmountChange = () => {
      const value = parseFloat(amountInput.value) || 0
      emit('on-value-changed', value)
    }
    
    const setMaxAmount = () => {
      amountInput.value = availableBalance.value.toString()
      onAmountChange()
    }
    
    // Watch for external amount changes
    watch(() => props.amount, (newAmount) => {
      if (newAmount !== parseFloat(amountInput.value) || 0) {
        amountInput.value = newAmount.toString()
      }
    })
    
    return {
      amountInput,
      availableBalance,
      estimatedPoolTokens,
      estimatedReturn,
      formatBalance,
      onAmountChange,
      setMaxAmount,
    }
  },
})
</script> 