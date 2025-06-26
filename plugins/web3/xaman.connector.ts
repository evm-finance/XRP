import { ConnectorInterface, Web3ErrorInterface } from '~/plugins/web3/connector'

export class XamanConnector implements ConnectorInterface {
  private account: string = ''
  private error: Web3ErrorInterface = { status: false, message: null }

  async connect(): Promise<{ account: string | null; error: Web3ErrorInterface }> {
    try {
      this.resetErrors()
      
      // Check if Xaman is installed
      if (!this.isXamanInstalled()) {
        this.error = {
          status: true,
          message: 'Xaman wallet not installed. Please install Xaman from https://xaman.app/'
        }
        return { account: null, error: this.error }
      }

      // Request connection to Xaman
      const result = await this.requestConnection()
      
      if (result.success && result.address) {
        this.account = result.address
        return { account: this.account, error: this.error }
      } else {
        this.error = {
          status: true,
          message: result.error || 'Failed to connect to Xaman wallet'
        }
        return { account: null, error: this.error }
      }
    } catch (err) {
      this.error = {
        status: true,
        message: 'Failed to connect to Xaman wallet'
      }
      return { account: null, error: this.error }
    }
  }

  private isXamanInstalled(): boolean {
    return typeof window !== 'undefined' && 
           (window as any).xumm !== undefined && 
           (window as any).xumm.xapp !== undefined
  }

  private async requestConnection(): Promise<{ success: boolean; address?: string; error?: string }> {
    try {
      const xumm = (window as any).xumm
      
      // Request user to sign in
      const signIn = await xumm.xapp.signin()
      
      if (signIn && signIn.me && signIn.me.account) {
        return {
          success: true,
          address: signIn.me.account
        }
      } else {
        return {
          success: false,
          error: 'User rejected connection'
        }
      }
    } catch (error) {
      return {
        success: false,
        error: 'Failed to connect to Xaman'
      }
    }
  }

  async signTransaction(transaction: any): Promise<{ success: boolean; signedTx?: any; error?: string }> {
    try {
      const xumm = (window as any).xumm
      
      // Submit transaction for signing
      const result = await xumm.xapp.signAndSubmit(transaction)
      
      if (result && result.result && result.result.txid) {
        return {
          success: true,
          signedTx: result.result
        }
      } else {
        return {
          success: false,
          error: 'Transaction signing failed'
        }
      }
    } catch (error) {
      return {
        success: false,
        error: 'Failed to sign transaction'
      }
    }
  }

  async getAccountInfo(): Promise<{ success: boolean; accountInfo?: any; error?: string }> {
    try {
      const xumm = (window as any).xumm
      
      const accountInfo = await xumm.xapp.accountInfo()
      
      if (accountInfo) {
        return {
          success: true,
          accountInfo
        }
      } else {
        return {
          success: false,
          error: 'Failed to get account info'
        }
      }
    } catch (error) {
      return {
        success: false,
        error: 'Failed to get account info'
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