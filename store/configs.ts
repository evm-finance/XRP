import type { ActionTree, MutationTree, GetterTree } from 'vuex'
import { Context } from '@nuxt/types'
import { ConfigState } from '~/types/state'
import { SearchResult, Network } from '~/types/global'

const defaultChain: Network = {
  id: '1',
  name: 'Ethereum',
  symbol: 'ETH',
  nativeTokenSymbol: 'ETH',
  rpcUrl: 'https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161',
  blockExplorerUrl: 'https://etherscan.io',
  dex: [],
  weth: {
    chainId: 1,
    address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
    symbol: 'WETH',
    name: 'Wrapped Ether',
    decimals: 18
  }
}

export const state = () =>
  ({
    title: 'XRP',
    globalStats: null,
    gasStats: null,
    networks: [],
    protocols: [],
    balancesChains: [1, 10, 56, 10000, 137, 250, 42161, 43114],
    globalSearchResult: [],
  } as ConfigState)

export const mutations: MutationTree<ConfigState> = {
  SET_CONFIG: (state, { networks }) => (state.networks = networks),
  SET_SEARCH_RESULT: (state, { searchResult }) => (state.globalSearchResult = searchResult),
}

export const actions: ActionTree<ConfigState, ConfigState> = {
  async initConfigs({ commit }, context: Context): Promise<void> {
    // XRP-only configuration - no need for network queries
    try {
      // Initialize with default XRP network
      const xrpNetwork: Network = {
        id: 'xrp',
        name: 'XRP Ledger',
        symbol: 'XRP',
        nativeTokenSymbol: 'XRP',
        rpcUrl: 'wss://xrplcluster.com',
        blockExplorerUrl: 'https://livenet.xrpl.org',
        dex: [],
        weth: {
          chainId: 0,
          address: '',
          symbol: 'XRP',
          name: 'XRP',
          decimals: 6
        }
      }
      commit('SET_CONFIG', { networks: [xrpNetwork] })
    } catch {}
  },

  async searchResult({ commit }, searchResult: SearchResult) {
    await commit('SET_SEARCH_RESULT', { searchResult })
  },
}

export const getters: GetterTree<ConfigState, ConfigState> = {
  chainInfo: (state: ConfigState) => (chainId: number) =>
    state.networks.find((elem: Network) => elem.id === chainId.toString()) ?? defaultChain,
} as any

export const getNetworkByChainId = (chainId: number): Network => {
  // This function should be used within a component with access to the store
  // For now, return defaultChain as fallback
  return defaultChain
}
