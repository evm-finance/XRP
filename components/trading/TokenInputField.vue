<template>
  <v-card tile outlined class="pa-3">
    <div class="d-flex justify-space-between align-center mb-2">
      <span class="text-caption grey--text">{{ label || 'Amount' }}</span>
      <span class="text-caption grey--text">
        Balance: {{ formatBalance(balance) }}
      </span>
    </div>
    
    <div class="d-flex align-center">
      <v-text-field
        v-model="inputValue"
        :placeholder="placeholder"
        type="number"
        outlined
        dense
        hide-details
        :readonly="readonly"
        @input="onInputChange"
      />
      
      <v-btn
        text
        small
        color="primary"
        class="ml-2"
        @click="setMaxAmount"
      >
        MAX
      </v-btn>
    </div>
    
    <div class="d-flex justify-space-between align-center mt-2">
      <div class="d-flex align-center">
        <v-avatar size="24" class="mr-2">
          <v-img :src="tokenIcon" />
        </v-avatar>
        <span class="text-body-2 font-weight-medium">{{ token.symbol }}</span>
      </div>
      
      <div v-if="amountUSD" class="text-caption grey--text">
        ${{ formatUSD(amountUSD) }}
      </div>
    </div>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, ref, watch } from '@nuxtjs/composition-api'

interface Token {
  symbol: string
  name: string
  icon: string
  issuer?: string
}

export default defineComponent({
  props: {
    token: { type: Object as () => Token, required: true },
    balance: { type: Number, default: 0 },
    amount: { type: Number, default: 0 },
    loading: { type: Boolean, default: false },
    readonly: { type: Boolean, default: false },
    label: { type: String, default: '' },
    placeholder: { type: String, default: '0.0' },
  },
  setup(props, { emit }) {
    const inputValue = ref('')
    
    const tokenIcon = computed(() => {
      // In real implementation, this would use proper token icons
      return props.token.icon || '🪙'
    })
    
    const amountUSD = computed(() => {
      // Mock USD value - in real implementation this would use price feeds
      if (!props.amount) return 0
      const mockPrice = props.token.symbol === 'XRP' ? 0.5 : 1.0
      return props.amount * mockPrice
    })
    
    const formatBalance = (value: number): string => {
      return value.toFixed(6)
    }
    
    const formatUSD = (value: number): string => {
      return value.toFixed(2)
    }
    
    const onInputChange = () => {
      const value = parseFloat(inputValue.value) || 0
      emit('on-value-changed', value)
    }
    
    const setMaxAmount = () => {
      inputValue.value = props.balance.toString()
      onInputChange()
    }
    
    // Watch for external amount changes
    watch(() => props.amount, (newAmount) => {
      if (newAmount !== parseFloat(inputValue.value) || 0) {
        inputValue.value = newAmount.toString()
      }
    })
    
    return {
      inputValue,
      tokenIcon,
      amountUSD,
      formatBalance,
      formatUSD,
      onInputChange,
      setMaxAmount,
    }
  },
})
</script>
