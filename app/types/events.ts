import { Block } from '~/types/apollo/main/types'

export enum EmitEvents {
  toggleNavigationMenu = 'toggle-navigation-menu',
  initAction = 'init-action',
  onValueChanged = 'on-value-changed',
  onTokenSelectChange = 'on-token-select-change',
  onResultClosed = 'on-result-closed',
  transactionSuccess = 'transaction-success',
  onNetworkSelectChange = 'on-network-select-change',
  onIntervalChange = 'on-interval-change',
  navigateToExplorer = 'navigate-to-explorer',
  onXrpInputOpen = 'on-xrp-input-open',
}
 
export enum DefiEvents {
  toggleActionDialog = 'toggle-action-dialog',
  toggleDepositModal = 'toggle-deposit-modal',
  toggleLendModal = 'toggle-lend-modal',
  toggleWithdrawModal = 'toggle-withdraw-modal',
}

export enum Events {
  onTokenSelect = 'on-token-select',
  onTokenMenuOpen = 'on-token-menu-open',
  onPoolSelect = 'on-pool-select',
  onPoolMenuOpen = 'on-pool-menu-open',
  onAccountSelect = 'on-account-select',
  onAccountMenuOpen = 'on-account-menu-open',
}

export type EmitterEvents = {
  priceStream: string
  onInitGlobalSearch: string
  onNewBlock: Block[]
}
