<template>
  <v-card height="400" outlined tile>
    <v-container fill-height>
      <v-row justify="center" align="center" align-content="center">
        <v-col cols="12" class="text-center">
          <v-icon style="margin-left: 100px; margin-bottom: -10px" color="red" right>mdi-close-circle</v-icon>
          <div>
            <v-avatar size="80"><img src="/img/metamask.svg" alt="" /></v-avatar>
          </div>
          <div class="title">Wallet Not Connected</div>
          <div class="subtitle-2">To start using EVMX, please connect your Metamask Wallet or search for wallet</div>
          <div class="py-3">
            <v-btn tile depressed @click="dispatch('ui/walletDialogStatus', true)">Connect to Wallet</v-btn>
          </div>
          <div class="py-1">OR</div>
          <div class="px-8">
            <span class="text-center">
              <v-text-field
                solo
                rounded
                dense
                hide-details
                single-line
                placeholder="Search for wallet"
                prepend-inner-icon="mdi-magnify"
                @click="initSearch"
              />
            </span>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </v-card>
</template>

<script lang="ts">
import { defineComponent, useContext, useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'
export default defineComponent({
  setup() {
    // COMPOSABLE
    const { dispatch } = useStore<State>()
    const { $emitter } = useContext()
    const initSearch = () => $emitter.emit('onInitGlobalSearch', '')

    return { dispatch, initSearch }
  },
})
</script>
