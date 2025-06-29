<template>
  <div v-if="customWalletAddress">
    <span class="mr-2 grey--text">Wallet:</span>
    <span>
      <v-avatar v-if="customWalletAddress" tile size="24" class="rounded mr-1">
        <img :src="avatarUrl" alt="" />
      </v-avatar>
      <v-chip label color="deep-purple darken-2" close close-icon="mdi-close" small @click:close="onCloseCustomAddress">
        {{ customWalletAddress }}
      </v-chip>
    </span>
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, toRefs, useStore } from '@nuxtjs/composition-api'
import { create } from 'blockies-ts'
import { State } from '~/types/state'
type Props = {
  selectedChainId: string
}
export default defineComponent<Props>({
  props: {
    selectedChainId: { type: String, required: true },
  },
  setup(props) {
    const { state, dispatch } = useStore<State>()
    const networkId = toRefs(props).selectedChainId

    const customWalletAddress = computed<any>({
      get: () => {
        const address =
          state.configs.globalSearchResult.find((elem) => elem.network === networkId.value)?.address ?? null
        return address === null
          ? null
          : `${address.slice(0, 5)}.....${address.slice(address.length - 5, address.length)}`
      },
      set: (value: any) => dispatch('configs/searchResult', value),
    })
    const avatarUrl = computed(() => create({ seed: customWalletAddress.value ?? '' }).toDataURL() ?? '')

    function onCloseCustomAddress() {
      customWalletAddress.value = []
    }

    return {
      avatarUrl,
      customWalletAddress,
      onCloseCustomAddress,
    }
  },
})
</script>
