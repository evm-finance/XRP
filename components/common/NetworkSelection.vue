<template>
  <v-menu v-if="currentNetwork" :close-on-content-click="false" nudge-bottom="13" offset-y>
    <template #activator="{ on, attrs }">
      <v-btn elevation="0" class="px-2" color="transparent" tile v-bind="attrs" :disabled="!walletReady" v-on="on">
        <v-avatar size="26" tile>
          <v-img :src="$imageUrlBySymbol(currentNetwork.symbol)" :lazy-src="$imageUrlBySymbol(currentNetwork.symbol)" />
        </v-avatar>
        <v-icon right>mdi-chevron-down</v-icon>
      </v-btn>
    </template>

    <v-card outlined tile>
      <v-row no-gutters>
        <v-col>
          <v-list dense>
            <v-list-item-group>
              <v-list-item
                v-for="item in allNetworks"
                :key="item.chainIdentifier"
                no-action
                :value="item"
                @change="onSelect(item)"
              >
                <v-list-item-avatar size="24">
                  <v-img :src="$imageUrlBySymbol(item.symbol)" :lazy-src="$imageUrlBySymbol(item.symbol)" />
                </v-list-item-avatar>
                <v-list-item-content><v-list-item-title />{{ item.name }}</v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-col>
      </v-row>
    </v-card>
  </v-menu>
</template>

<script lang="ts">
import { defineComponent, inject } from '@nuxtjs/composition-api'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'

export default defineComponent({
  setup() {
    const {
      allNetworks,
      changeNetwork: onSelect,
      currentNetwork,
      chainId,
      walletReady,
    } = inject(WEB3_PLUGIN_KEY) as Web3

    return {
      allNetworks,
      currentNetwork,
      walletReady,
      chainId,
      onSelect,
    }
  },
})
</script>
