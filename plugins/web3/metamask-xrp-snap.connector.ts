import { ConnectorInterface, Web3ErrorInterface } from '~/plugins/web3/connector'

export class MetamaskXrpSnapConnector implements ConnectorInterface {
  private account: string = ''
  private error: Web3ErrorInterface = { status: false, message: null }
  private snapId: string = 'npm:@metamask/xrp-snap'

  async connect(): Promise<{ account: string | null; error: Web3ErrorInterface }> {
    try {
      this.resetErrors()
      
      // Check if MetaMask is installed
      if (!this.isMetamaskInstalled()) {
        this.error = {
          status: true,
          message: 'MetaMask not installed. Please install MetaMask first.'
        }
        return { account: null, error: this.error }
      }

      // Check if XRP snap is installed
      const isSnapInstalled = await this.isXrpSnapInstalled()
      if (!isSnapInstalled) {
        // Install the XRP snap
        const installResult = await this.installXrpSnap()
        if (!installResult.success) {
          return { account: null, error: { status: true, message: installResult.error } }
        }
      }

      // Get XRP account
      const accountResult = await this.getXrpAccount()
      
      if (accountResult.success && accountResult.address) {
        this.account = accountResult.address
        return { account: this.account, error: this.error }
      } else {
        this.error = {
          status: true,
          message: accountResult.error || 'Failed to get XRP account'
        }
        return { account: null, error: this.error }
      }
    } catch (err) {
      this.error = {
        status: true,
        message: 'Failed to connect to MetaMask XRP snap'
      }
      return { account: null, error: this.error }
    }
  }

  private isMetamaskInstalled(): boolean {
    return typeof window !== 'undefined' && 
           (window as any).ethereum !== undefined &&
           (window as any).ethereum.isMetaMask === true
  }

  private async isXrpSnapInstalled(): Promise<boolean> {
    try {
      const ethereum = (window as any).ethereum
      const snaps = await ethereum.request({
        method: 'wallet_getSnaps'
      })
      
      return Object.keys(snaps).includes(this.snapId)
    } catch (error) {
      return false
    }
  }

  private async installXrpSnap(): Promise<{ success: boolean; error?: string }> {
    try {
      const ethereum = (window as any).ethereum
      
      await ethereum.request({
        method: 'wallet_installSnaps',
        params: {
          [this.snapId]: {
            version: 'latest'
          }
        }
      })
      
      return { success: true }
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to install XRP snap'
      }
    }
  }

  private async getXrpAccount(): Promise<{ success: boolean; address?: string; error?: string }> {
    try {
      const ethereum = (window as any).ethereum
      
      const result = await ethereum.request({
        method: 'wallet_invokeSnap',
        params: {
          snapId: this.snapId,
          request: {
            method: 'getAccount'
          }
        }
      })
      
      if (result && result.address) {
        return {
          success: true,
          address: result.address
        }
      } else {
        return {
          success: false,
          error: 'Failed to get XRP account'
        }
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to get XRP account'
      }
    }
  }

  async signTransaction(transaction: any): Promise<{ success: boolean; signedTx?: any; error?: string }> {
    try {
      const ethereum = (window as any).ethereum
      
      const result = await ethereum.request({
        method: 'wallet_invokeSnap',
        params: {
          snapId: this.snapId,
          request: {
            method: 'signTransaction',
            params: {
              transaction
            }
          }
        }
      })
      
      if (result && result.signedTransaction) {
        return {
          success: true,
          signedTx: result.signedTransaction
        }
      } else {
        return {
          success: false,
          error: 'Transaction signing failed'
        }
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to sign transaction'
      }
    }
  }

  async getAccountInfo(): Promise<{ success: boolean; accountInfo?: any; error?: string }> {
    try {
      const ethereum = (window as any).ethereum
      
      const result = await ethereum.request({
        method: 'wallet_invokeSnap',
        params: {
          snapId: this.snapId,
          request: {
            method: 'getAccountInfo'
          }
        }
      })
      
      if (result) {
        return {
          success: true,
          accountInfo: result
        }
      } else {
        return {
          success: false,
          error: 'Failed to get account info'
        }
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to get account info'
      }
    }
  }

  async getBalance(): Promise<{ success: boolean; balance?: string; error?: string }> {
    try {
      const ethereum = (window as any).ethereum
      
      const result = await ethereum.request({
        method: 'wallet_invokeSnap',
        params: {
          snapId: this.snapId,
          request: {
            method: 'getBalance'
          }
        }
      })
      
      if (result && result.balance) {
        return {
          success: true,
          balance: result.balance
        }
      } else {
        return {
          success: false,
          error: 'Failed to get balance'
        }
      }
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to get balance'
      }
    }
  }

  resetErrors(): void {
    this.error = { status: false, message: null }
  }

  disconnect(): void {
    this.account = ''
    this.resetErrors()
  }

  getAccount(): string {
    return this.account
  }

  getError(): Web3ErrorInterface {
    return this.error
  }
} 