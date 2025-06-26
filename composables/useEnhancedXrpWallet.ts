import { inject, computed } from '@nuxtjs/composition-api'
import { EnhancedXrpClient, ENHANCED_XRP_PLUGIN_KEY, XrpWalletType } from '~/plugins/web3/enhanced-xrp.client'

export default function useEnhancedXrpWallet() {
  const enhancedXrpClient = inject(ENHANCED_XRP_PLUGIN_KEY) as EnhancedXrpClient

  if (!enhancedXrpClient) {
    throw new Error('Enhanced XRP client not found. Make sure the plugin is properly registered.')
  }

  // Wallet connection
  const connectWallet = async (walletType: XrpWalletType) => {
    return await enhancedXrpClient.connectWallet(walletType)
  }

  const disconnectWallet = () => {
    enhancedXrpClient.disconnectWallet()
  }

  // Wallet state
  const isWalletReady = computed(() => enhancedXrpClient.isWalletReady.value)
  const address = computed(() => enhancedXrpClient.address.value)
  const currentWalletType = computed(() => enhancedXrpClient.currentWalletType.value)
  const error = computed(() => enhancedXrpClient.error.value)

  // Wallet operations
  const signTransaction = async (transaction: any) => {
    return await enhancedXrpClient.signTransaction(transaction)
  }

  const getAccountInfo = async () => {
    return await enhancedXrpClient.getAccountInfo()
  }

  const getBalance = async () => {
    return await enhancedXrpClient.getBalance()
  }

  // Wallet type helpers
  const isGemWallet = computed(() => currentWalletType.value === 'gem')
  const isXamanWallet = computed(() => currentWalletType.value === 'xaman')
  const isMetamaskXrpWallet = computed(() => currentWalletType.value === 'metamask-xrp-snap')

  // Wallet capabilities
  const canSignTransactions = computed(() => {
    return isWalletReady.value && (isXamanWallet.value || isMetamaskXrpWallet.value)
  })

  const canGetAccountInfo = computed(() => {
    return isWalletReady.value && (isXamanWallet.value || isMetamaskXrpWallet.value)
  })

  const canGetBalance = computed(() => {
    return isWalletReady.value && isMetamaskXrpWallet.value
  })

  // Transaction helpers
  const signPaymentTransaction = async (params: {
    destination: string
    amount: string
    currency?: string
    issuer?: string
  }) => {
    const transaction = {
      TransactionType: 'Payment',
      Account: address.value,
      Destination: params.destination,
      Amount: params.currency ? {
        currency: params.currency,
        issuer: params.issuer,
        value: params.amount
      } : params.amount
    }

    return await signTransaction(transaction)
  }

  const signTrustSetTransaction = async (params: {
    currency: string
    issuer: string
    limit: string
  }) => {
    const transaction = {
      TransactionType: 'TrustSet',
      Account: address.value,
      LimitAmount: {
        currency: params.currency,
        issuer: params.issuer,
        value: params.limit
      }
    }

    return await signTransaction(transaction)
  }

  const signAMMDepositTransaction = async (params: {
    amount: string
    amount2Currency?: string
    amount2Issuer?: string
    amount2Value?: string
  }) => {
    const transaction = {
      TransactionType: 'AMMDeposit',
      Account: address.value,
      Amount: params.amount,
      ...(params.amount2Currency && {
        Amount2: {
          currency: params.amount2Currency,
          issuer: params.amount2Issuer,
          value: params.amount2Value
        }
      })
    }

    return await signTransaction(transaction)
  }

  const signAMMWithdrawTransaction = async (params: {
    amount: string
    amount2Currency?: string
    amount2Issuer?: string
    amount2Value?: string
  }) => {
    const transaction = {
      TransactionType: 'AMMWithdraw',
      Account: address.value,
      Amount: params.amount,
      ...(params.amount2Currency && {
        Amount2: {
          currency: params.amount2Currency,
          issuer: params.amount2Issuer,
          value: params.amount2Value
        }
      })
    }

    return await signTransaction(transaction)
  }

  const signAMMTradeTransaction = async (params: {
    amount: string
    amount2Currency: string
    amount2Issuer: string
    amount2Value: string
  }) => {
    const transaction = {
      TransactionType: 'AMMTrade',
      Account: address.value,
      Amount: params.amount,
      Amount2: {
        currency: params.amount2Currency,
        issuer: params.amount2Issuer,
        value: params.amount2Value
      }
    }

    return await signTransaction(transaction)
  }

  return {
    // Connection
    connectWallet,
    disconnectWallet,
    
    // State
    isWalletReady,
    address,
    currentWalletType,
    error,
    
    // Operations
    signTransaction,
    getAccountInfo,
    getBalance,
    
    // Wallet type helpers
    isGemWallet,
    isXamanWallet,
    isMetamaskXrpWallet,
    
    // Capabilities
    canSignTransactions,
    canGetAccountInfo,
    canGetBalance,
    
    // Transaction helpers
    signPaymentTransaction,
    signTrustSetTransaction,
    signAMMDepositTransaction,
    signAMMWithdrawTransaction,
    signAMMTradeTransaction,
  }
} 