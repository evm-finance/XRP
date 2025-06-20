<template>
  <v-card outlined tile class="pa-3">
    <v-row no-gutters align="center">
      <v-col cols="8">
        <v-text-field
          v-model="amount"
          :loading="loading"
          :disabled="loading"
          type="number"
          placeholder="0.0"
          hide-details
          dense
          outlined
          @input="onInput"
        />
      </v-col>
      <v-col cols="4" class="pl-2">
        <v-btn
          text
          block
          outlined
          tile
          @click="$emit('on-token-menu-open')"
        >
          <v-avatar size="20" class="mr-1">
            <v-img :src="$imageUrlBySymbol(token.symbol.toLowerCase())" />
          </v-avatar>
          {{ token.symbol }}
          <v-icon small class="ml-1">mdi-chevron-down</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    
    <v-row no-gutters class="mt-2">
      <v-col cols="6">
        <span class="caption grey--text">
          Balance: {{ formatBalance(balance) }}
        </span>
      </v-col>
      <v-col cols="6" class="text-right">
        <span class="caption grey--text">
          {{ fiatPrice > 0 ? `$${fiatPrice.toFixed(2)}` : '' }}
        </span>
      </v-col>
    </v-row>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, watch } from '@nuxtjs/composition-api'
import { XRPToken } from '~/composables/useXrpAmm'

export default defineComponent({
  props: {
    formTradeDirection: { type: String, default: 'EXACT_INPUT' },
    token: { type: Object as PropType<XRPToken>, required: true },
    balance: { type: Number, default: 0 },
    fiatPrice: { type: Number, default: 0 },
    expectedConvertQuote: { type: Number, default: 0 },
    loading: { type: Boolean, default: false }
  },

  setup(props, { emit }) {
    const amount = ref('')

    const formatBalance = (balance: number): string => {
      if (balance === 0) return '0'
      if (balance < 0.000001) return '< 0.000001'
      return balance.toFixed(6)
    }

    const onInput = () => {
      const numValue = parseFloat(amount.value) || 0
      emit('on-value-changed', { value: numValue })
    }

    // Watch for expected convert quote changes (for output field)
    watch(() => props.expectedConvertQuote, (newValue) => {
      if (props.formTradeDirection === 'EXACT_OUTPUT' && newValue > 0) {
        amount.value = newValue.toFixed(6)
      }
    })

    return {
      amount,
      formatBalance,
      onInput
    }
  }
})
</script>

<style scoped></style> 