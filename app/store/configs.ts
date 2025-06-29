import type { ActionTree, MutationTree, GetterTree } from 'vuex'
import { Context } from '@nuxt/types'
import { ConfigState } from '~/types/state'
import { SearchResult, Network } from '~/types/global'

const defaultChain: Network = {
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
    globalStats: null,
    gasStats: null,
    networks: [],
    protocols: [],
    balancesChains: ['xrp'],
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
  chainInfo: (state: ConfigState) => (chainId: number) =>
    state.networks.find((elem: Network) => elem.id === chainId.toString()) ?? defaultChain,
} as any
