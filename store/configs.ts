import type { ActionTree, MutationTree, GetterTree } from 'vuex'
import { Context } from '@nuxt/types'
import { ConfigState } from '~/types/state'
import { SearchResult, Network } from '~/types/global'

const defaultNetwork: Network = {
  xrp: {
    symbol: 'XRP',
    decimals: 6,
    name: 'XRP',
  },
  dex: [
    { value: 'xrp_amm', name: 'XRP AMM', symbol: 'XRP' },
  ],
  id: 'xrp',
  blockExplorerUrl: 'https://livenet.xrpl.org/',
  name: 'XRP Ledger',
  rpcUrl: 'wss://xrplcluster.com/',
  symbol: 'XRP',
  nativeTokenSymbol: 'XRP',
}

export const state = () =>
  ({
    title: 'XRP Finance',
    networks: [defaultNetwork],
    balancesChains: ['xrp'],
    protocols: [],
    globalSearchResult: [],
  } as ConfigState)

export const mutations: MutationTree<ConfigState> = {
  SET_CONFIG: (state, { networks }) => (state.networks = networks),
  SET_SEARCH_RESULT: (state, { searchResult }) => (state.globalSearchResult = searchResult),
}

export const actions: ActionTree<ConfigState, ConfigState> = {
  async searchResult({ commit }, searchResult: SearchResult) {
    await commit('SET_SEARCH_RESULT', { searchResult })
  },
}

export const getters: GetterTree<ConfigState, ConfigState> = {
  currentNetwork: (state: ConfigState) => state.networks[0] ?? defaultNetwork,
} as any
