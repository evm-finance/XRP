<template>
  <div></div>
  <!--  <v-row v-if="isChainAndMarketMismatched" justify="center" no-gutters>-->
  <!--    <v-dialog v-model="dialog" class="rounded-0" persistent max-width="380">-->
  <!--      <v-card tile outlined elevation="0">-->
  <!--        <v-card-title class="text-no-wrap">-->
  <!--          Current Network - {{ currentNetwork }} <v-spacer /><v-icon @click="dialog = false">mdi-close</v-icon>-->
  <!--        </v-card-title>-->
  <!--        <v-row justify="center" class="pt-2 pb-0" no-gutters>-->
  <!--          <v-col class="text-center" cols="12">-->
  <!--            <div class="pb-3">-->
  <!--              <v-avatar>-->
  <!--                <img :alt="`${isChainAndMarketMismatched.label} logo`" :src="isChainAndMarketMismatched.logoUrl" />-->
  <!--              </v-avatar>-->
  <!--            </div>-->
  <!--          </v-col>-->
  <!--        </v-row>-->

  <!--        <v-card-text class="text-center">-->
  <!--          Please switch to <span class="subtitle-2 pink&#45;&#45;text">{{ isChainAndMarketMismatched.label }}</span> <br />-->
  <!--          to perform any actions.-->
  <!--        </v-card-text>-->

  <!--        <v-card-actions class="justify-center pb-4">-->
  <!--          <v-btn depressed block tile class="text-transform-none px-3 ma-0" height="32" @click="changeToRequiredChain">-->
  <!--            Switch network-->
  <!--          </v-btn>-->
  <!--        </v-card-actions>-->
  <!--      </v-card>-->
  <!--    </v-dialog>-->
  <!--  </v-row>-->
</template>

<script lang="ts">
import { ref, defineComponent, inject, useStore, computed } from '@nuxtjs/composition-api'
// import useMarketSelector from '~/composables/useMarketSelector'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { State } from '~/types/state'

export default defineComponent({
  name: 'SwitchNetworkDialog',
  setup() {
    const { state } = useStore<State>()
    const { chainId } = inject(WEB3_PLUGIN_KEY) as Web3
    // const { changeToRequiredChain, isChainAndMarketMismatched } = useMarketSelector()

    const dialog = ref(false)

    // get currentNetwork name, if not supported, show chain id
    const currentNetwork = computed<string>(() => {
      if (chainId.value) {
        const found = state.configs.networks.find((chain) => chain.chainIdentifier === chainId.value)

        if (found) {
          return found.name
        }

        return `Chain ID: ${chainId.value}`
      }
      return ''
    })

    function toggleDialog(value: boolean) {
      dialog.value = value
    }

    return {
      dialog,
      chainId,
      currentNetwork,
      toggleDialog,
    }
  },
})
</script>
