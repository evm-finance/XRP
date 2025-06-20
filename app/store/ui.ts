/* eslint-disable @typescript-eslint/no-unused-vars */
import type { ActionTree, MutationTree } from 'vuex'
import { ThemeOptions } from '~/types/state'

export const state = () => ({
  theme: 'dark' as ThemeOptions,
  dark: {
    overlayColor: 'grey darken-4',
    // header
    headerTextClass: 'white--text',
    innerCardLighten: 'grey--text text--lighten-1',
    // Json logs for transactions dialog
    jsonLogs: 'grey darken-4',
    // texts
    subTextColor: 'grey--text',
  },
  light: {
    overlayColor: 'grey lighten-4',
    headerTextClass: 'black--text',
    innerCardLighten: 'black--text',
    jsonLogs: 'grey lighten-2',
    subTextColor: 'grey--text text--darken-1',
  },
  walletSelectionDialog: false as boolean,
})
export type UiState = ReturnType<typeof state>

export const mutations: MutationTree<UiState> = {
  SET_THEME: (state: UiState, theme: ThemeOptions) => (state.theme = theme),
  SET_WALLET_DIALOG_STATUS: (state: UiState, status: boolean) => (state.walletSelectionDialog = status),
}

export const actions: ActionTree<UiState, UiState> = {
  changeTheme({ commit }, theme: string): void {
    commit('SET_THEME', theme)
  },

  walletDialogStatus({ commit, state }, status: boolean): void {
    commit('SET_WALLET_DIALOG_STATUS', status)
  },
}
