<template>
  <v-dialog
    v-model="dialog"
    max-width="580"
    width="580"
    persistent
    transition="dialog-top-transition"
    content-class="my-custom-dialog"
  >
    <v-card tile outlined height="auto" min-height="200" width="600">
      <div class="text-right mb-3">
        <v-icon class="my-2 ma-3" @click="dialog = !dialog">mdi-close</v-icon>
      </div>
      <v-row class="px-4">
        <v-col>
          <v-text-field
            v-model="search"
            solo
            dense
            clearable
            hide-details
            placeholder="Search for wallet"
            @keydown.enter.prevent="onSearch"
          >
            <template #append-outer>
              <v-btn depressed style="top: -6px" @click="onSearch">
                <v-icon size="20"> mdi-magnify</v-icon>
              </v-btn>
            </template>
          </v-text-field>
        </v-col>
      </v-row>

      <div v-if="loading" class="text-center mt-4">
        <v-progress-circular indeterminate color="primary" />
      </div>

      <v-row v-else no-gutters justify="center">
        <v-col cols="12">
          <v-list dense>
            <v-list-item-group>
              <v-list-item v-for="item in searchResult" :key="item.network.chainIdentifier" no-action :value="item">
                <v-list-item-avatar size="24">
                  <v-img
                    :src="$imageUrlBySymbol(item.network.symbol)"
                    :lazy-src="$imageUrlBySymbol(item.network.symbol)"
                  />
                </v-list-item-avatar>
                <v-list-item-content>
                  <v-list-item-title v-text="item.network.name" />
                  <v-list-item-subtitle v-text="item.desc" />
                </v-list-item-content>

                <v-list-item-action>
                  <span>
                    <v-chip v-if="item.isWallet" x-small color="green  darken-2">Wallet</v-chip>
                    <v-chip v-if="item.isContract" x-small color="deep-purple darken-2">Contract</v-chip>
                    <v-icon v-if="!item.isWallet && !item.isContract" size="18" color="red">mdi-close</v-icon>
                  </span>
                </v-list-item-action>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-col>

        <v-col v-if="searchResult.length" cols="12">
          <div class="text-center my-4">
            <v-btn text color="primary" @click="useSearch">OK use this address</v-btn>
          </div>
        </v-col>
      </v-row>
    </v-card>
  </v-dialog>
</template>
<script lang="ts">
import * as process from 'process'
import { computed, defineComponent, inject, ref, useContext, useStore } from '@nuxtjs/composition-api'
import { ethers } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { Chain } from '~/types/apollo/main/types'
import { SearchResult } from '~/types/global'
import { State } from '~/types/state'
export default defineComponent({
  setup() {
    const { $emitter } = useContext()
    const { state, dispatch } = useStore<State>()
    const { allNetworks, getCustomProviderByNetworkId } = inject(WEB3_PLUGIN_KEY) as Web3
    const searchResult = computed({
      get: () => state.configs.globalSearchResult,
      set: (value: any) => dispatch('configs/searchResult', value),
    })

    const dialog = ref(false)
    const search = ref<string | null>('')
    const loading = ref(false)

    const checkAddressType = async (address: string, network: Chain, provider: any): Promise<SearchResult> => {
      try {
        // Get the code at the address
        const code = await provider.getCode(address)

        // If the code is "0x", it is not a contract address (wallet address)
        if (code === '0x') {
          const transactionCount = await provider.getTransactionCount(address)
          if (transactionCount > 0) {
            return {
              desc: 'Address is a wallet (has transaction history)',
              network,
              isWallet: true,
              isContract: false,
              searchString: address,
            }
          }

          // Check ETH balance
          const balance = await provider.getBalance(address)
          if (balance.gt(ethers.constants.Zero)) {
            return {
              desc: 'Address is a wallet (holds Ether)',
              network,
              isWallet: true,
              isContract: false,
              searchString: address,
            }
          }
          return { desc: 'Valid wallet', network, isWallet: true, isContract: false, searchString: address }
        } else {
          return { desc: 'invalid wallet', network, isWallet: false, isContract: true, searchString: address }
        }
      } catch (error) {
        return { desc: 'Invalid address', network, isWallet: false, isContract: false, searchString: address }
      }
    }

    const onSearch = async () => {
      if (search.value === null || search.value.length < 1) {
        searchResult.value = []
        return
      }
      loading.value = true
      const multCalls: Promise<SearchResult>[] = []
      allNetworks.value.forEach((elem) => {
        const provider = getCustomProviderByNetworkId(elem.id)
        const request = checkAddressType(search.value ?? '', elem, provider)
        multCalls.push(request)
      })
      searchResult.value = await Promise.all(multCalls)
      loading.value = false
    }

    function openDialog() {
      dialog.value = true
    }
    function useSearch() {
      dialog.value = false
    }
    // TODO: Change to VUEX handler instead of event emitter
    // eslint-disable-next-line no-unused-expressions
    $emitter?.on('onInitGlobalSearch', () => openDialog())

    return {
      loading,
      searchResult,
      dialog,
      search,
      onSearch,
      openDialog,
      useSearch,
    }
  },
})
</script>
<style scoped lang="css">
>>> .my-custom-dialog {
  align-self: flex-start;
}
</style>
