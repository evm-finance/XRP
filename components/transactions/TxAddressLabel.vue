<template>
  <div class="text-no-wrap d-flex grey--text text--lighten-1">
    {{ label }}
    <nuxt-link v-if="isContract" class="ml-1 cursor-copy white--text" :to="`/portfolio/balances?wallet=${address}`">
      <v-avatar v-if="symbol" size="16" style="margin-top: -2px; margin-right: 2px; margin-left: 2px">
        <img alt="" :src="$imageUrlBySymbol(symbol.toLowerCase())" @error="$setAltImageUrl" />
      </v-avatar>
      <v-icon v-else class="ml-1 mt-n1" small>mdi-file-sign</v-icon>
      {{ !name ? $truncateAddress(address, 4, 10) : name }}
    </nuxt-link>

    <nuxt-link
      v-else
      :to="`/portfolio/balances?wallet=${address}`"
      :class="['cursor-copy', 'ml-1', address === walletAddress ? 'pink--text' : 'white--text']"
    >
      {{ $truncateAddress(address, 4, 10) }}
    </nuxt-link>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from '@nuxtjs/composition-api'
type navigateToExplorerType = (address: string, type: 'tx' | 'address', blockExplorerUrl: string) => void
export default defineComponent({
  props: {
    label: { type: String, required: true },
    walletAddress: { type: String, required: true },
    isContract: { type: Boolean, default: false },
    address: { type: String, required: true },
    name: { type: String, default: '' },
    symbol: { type: String, default: '' },
    navigateToExplorer: { type: Function as PropType<navigateToExplorerType>, required: true },
    blockExplorerUrl: { type: String, default: '' },
  },
})
</script>
