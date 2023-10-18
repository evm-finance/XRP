import type { ActionTree, MutationTree, GetterTree } from 'vuex'
import { Context } from '@nuxt/types'
import { Chain } from '~/types/apollo/main/types'
import { SupportedChainsGQL } from '~/apollo/main/config.query.graphql'
import { ConfigState } from '~/types/state'
import { SearchResult } from '~/types/global'

const defaultChain: Chain = {
  weth: {
    chainId: 1,
    symbol: 'WETH',
    decimals: 18,
    address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
    name: 'Wrapped Ether',
  },
  dex: [
    { value: 'uniswap_v2', name: 'Uniswap V2', symbol: 'UNI' },
    { value: 'uniswap_v3', name: 'Uniswap V3', symbol: 'UNI' },
  ],
  id: 'ethereum',
  blockExplorerUrl: 'https://etherscan.io/',
  chainIdentifier: 1,
  name: 'Ethereum Main Net',
  rpcUrl: 'https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161/',
  symbol: 'ETH',
}

export const state = () =>
  ({
    title: 'EVM',
    globalStats: null,
    gasStats: null,
    chains: [],
    protocols: [],
    // -----------> ETH OP BSC MATIC FTM  BASE   BCH    ARB    AVAX  
    balancesChains: [1, 10, 56, 137, 250, 8453, 10000, 42161, 43114],
    globalSearchResult: [],
  } as ConfigState)

export const mutations: MutationTree<ConfigState> = {
  SET_CONFIG: (state, { chains }) => (state.chains = chains),
  SET_SEARCH_RESULT: (state, { searchResult }) => (state.globalSearchResult = searchResult),
}

export const actions: ActionTree<ConfigState, ConfigState> = {
  async initConfigs({ commit }, context: Context): Promise<void> {
    try {
      const client = context.app.apolloProvider?.defaultClient
      const query = await client?.query({ query: SupportedChainsGQL, fetchPolicy: 'no-cache' })
      if (query && query.data) {
        commit('SET_CONFIG', query.data)
      }
    } catch {}
  },

  async searchResult({ commit }, searchResult: SearchResult) {
    await commit('SET_SEARCH_RESULT', { searchResult })
  },
}
export const getters: GetterTree<ConfigState, ConfigState> = {
  chainInfo: (state: ConfigState) => (chainId: number) =>
    state.chains.find((elem: Chain) => elem.chainIdentifier === chainId) ?? defaultChain,
} as any
