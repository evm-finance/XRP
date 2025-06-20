<template>
  <v-dialog v-model="dialog" max-width="400">
    <v-card>
      <v-card-title>
        Select Token
        <v-spacer />
        <v-btn icon @click="dialog = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      
      <v-card-text>
        <v-text-field
          v-model="search"
          placeholder="Search tokens..."
          prepend-inner-icon="mdi-magnify"
          outlined
          dense
          hide-details
          class="mb-4"
        />
        
        <v-list>
          <v-list-item
            v-for="token in filteredTokens"
            :key="`${token.currency}-${token.issuer || 'XRP'}`"
            @click="selectToken(token)"
          >
            <v-list-item-avatar size="32">
              <v-img :src="$imageUrlBySymbol(token.symbol.toLowerCase())" />
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ token.symbol }}</v-list-item-title>
              <v-list-item-subtitle>{{ token.name }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from '@nuxtjs/composition-api'
import { XRPToken } from '~/composables/useXrpAmm'

export default defineComponent({
  setup(props, { emit }) {
    const dialog = ref(false)
    const search = ref('')

    // Common XRP tokens
    const tokens: XRPToken[] = [
      {
        currency: 'XRP',
        symbol: 'XRP',
        name: 'XRP',
        decimals: 6
      },
      {
        currency: 'USDC',
        issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
        symbol: 'USDC',
        name: 'USD Coin',
        decimals: 6
      },
      {
        currency: 'USDT',
        issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
        symbol: 'USDT',
        name: 'Tether USD',
        decimals: 6
      },
      {
        currency: 'ETH',
        issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
        symbol: 'ETH',
        name: 'Ethereum',
        decimals: 6
      },
      {
        currency: 'BTC',
        issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
        symbol: 'BTC',
        name: 'Bitcoin',
        decimals: 6
      }
    ]

    const filteredTokens = computed(() => {
      if (!search.value) return tokens
      const searchLower = search.value.toLowerCase()
      return tokens.filter(token => 
        token.symbol.toLowerCase().includes(searchLower) ||
        token.name.toLowerCase().includes(searchLower)
      )
    })

    const selectToken = (token: XRPToken) => {
      emit('on-token-select', token)
      dialog.value = false
      search.value = ''
    }

    return {
      dialog,
      search,
      filteredTokens,
      selectToken
    }
  }
})
</script>

<style scoped></style> 