import { computed, ref, useContext } from '@nuxtjs/composition-api'
import useEnhancedXrpWallet from '~/composables/useEnhancedXrpWallet'
import useXrpAmmLiveData from '~/composables/useXrpAmmLiveData'
import { useXrplTransaction } from '~/composables/useXrplTransaction'

interface XrpAmmPool {
  id: string
  token1: XrpToken
  token2: XrpToken
  liquidity: number
  volume24h: number
  fee: number
  apr: number
  priceChange24h: number
  token1Balance: number
  token2Balance: number
}

interface XrpToken {
  symbol: string
  name: string
  icon: string
  issuer?: string
}

export default function useXrpAmmTransactions(pool: XrpAmmPool, amount: any) {
  const { $f } = useContext()
  const { 
    wallet, 
    isWalletReady, 
    signAMMDepositTransaction,
    signAMMWithdrawTransaction,
    signAMMTradeTransaction,
    canSignTransactions,
    address
  } = useEnhancedXrpWallet()
  
  const { refreshAll } = useXrpAmmLiveData()
  const { submitAndConfirm, submitting, confirming, error: txError, clearTransactionState } = useXrplTransaction()
  
  // Transaction state
  const txLoading = ref(false)
  const receipt = ref<any>(null)
  const isTxMined = ref(false)
  const error = ref<string | null>(null)
  const txHash = ref<string | null>(null)
  
  // Transaction status tracking
  const txStatus = ref<'idle' | 'pending' | 'success' | 'error'>('idle')
  
  // Validate transaction parameters
  const validateTransaction = (action: string, amount: number) => {
    if (!isWalletReady.value) {
      throw new Error('Wallet not connected')
    }
    
    if (!canSignTransactions.value) {
      throw new Error('Wallet does not support transaction signing')
    }
    
    if (!amount || amount <= 0) {
      throw new Error('Invalid amount')
    }
    
    if (!address.value) {
      throw new Error('No wallet address available')
    }
  }
  
  // Deposit to AMM pool
  const deposit = async () => {
    validateTransaction('deposit', amount.value)
    
    txLoading.value = true
    txStatus.value = 'pending'
    error.value = null
    receipt.value = null
    txHash.value = null
    
    try {
      // Create AMM deposit transaction
      const depositParams = {
        amount: amount.value.toString(),
        amount2Currency: pool.token2.symbol,
        amount2Issuer: pool.token2.issuer,
        amount2Value: (amount.value * 0.95).toString() // Estimate based on pool ratio
      }
      
      console.log('Creating AMM deposit transaction:', depositParams)
      
      // Sign transaction with wallet
      const signedTx = await signAMMDepositTransaction()
      
      // Submit and confirm transaction
      const txReceipt = await submitAndConfirm(signedTx)
      
      // Update state
      txHash.value = txReceipt.hash
      receipt.value = {
        ...txReceipt,
        to: pool.id,
        poolId: pool.id
      }
      
      isTxMined.value = true
      txStatus.value = 'success'
      
      // Refresh data
      await refreshAll()
      
      console.log(`Successfully deposited ${amount.value} ${pool.token1.symbol} to pool ${pool.id}`)
      
    } catch (err: any) {
      error.value = err.message || 'Deposit failed'
      txStatus.value = 'error'
      console.error('Deposit error:', err)
    } finally {
      txLoading.value = false
    }
  }
  
  // Withdraw from AMM pool
  const withdraw = async () => {
    validateTransaction('withdraw', amount.value)
    
    txLoading.value = true
    txStatus.value = 'pending'
    error.value = null
    receipt.value = null
    txHash.value = null
    
    try {
      // Create AMM withdraw transaction
      const withdrawParams = {
        amount: amount.value.toString(),
        amount2Currency: pool.token1.symbol,
        amount2Issuer: pool.token1.issuer,
        amount2Value: (amount.value * 0.95).toString() // Estimate based on pool ratio
      }
      
      console.log('Creating AMM withdraw transaction:', withdrawParams)
      
      // Sign transaction with wallet
      const signedTx = await signAMMWithdrawTransaction()
      
      // Submit and confirm transaction
      const txReceipt = await submitAndConfirm(signedTx)
      
      // Update state
      txHash.value = txReceipt.hash
      receipt.value = {
        ...txReceipt,
        to: pool.id,
        poolId: pool.id
      }
      
      isTxMined.value = true
      txStatus.value = 'success'
      
      // Refresh data
      await refreshAll()
      
      console.log(`Successfully withdrew ${amount.value} LP tokens from pool ${pool.id}`)
      
    } catch (err: any) {
      error.value = err.message || 'Withdraw failed'
      txStatus.value = 'error'
      console.error('Withdraw error:', err)
    } finally {
      txLoading.value = false
    }
  }
  
  // Reset transaction state
  const resetToDefault = () => {
    txLoading.value = false
    txStatus.value = 'idle'
    receipt.value = null
    isTxMined.value = false
    error.value = null
    txHash.value = null
    clearTransactionState()
  }
  
  // Transaction status
  const txStatusComputed = computed(() => txStatus.value)
  
  // Can perform transaction
  const canPerformTransaction = computed(() => {
    return isWalletReady.value && 
           canSignTransactions.value && 
           amount.value > 0 && 
           txStatus.value === 'idle' &&
           !submitting.value &&
           !confirming.value
  })
  
  // Combined loading state
  const isLoading = computed(() => 
    txLoading.value || 
    submitting.value || 
    confirming.value
  )
  
  // Combined error state
  const hasError = computed(() => 
    !!error.value || 
    !!txError.value
  )
  
  return {
    // State
    txLoading: isLoading,
    receipt,
    isTxMined,
    error: computed(() => error.value || txError.value),
    txHash,
    txStatus: txStatusComputed,
    
    // Computed
    canPerformTransaction,
    
    // Methods
    deposit,
    withdraw,
    resetToDefault,
  }
} 