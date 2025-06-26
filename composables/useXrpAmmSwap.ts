import { computed, ref, useContext, watch } from '@nuxtjs/composition-api'
import { useEnhancedXrpWallet } from '~/composables/useEnhancedXrpWallet'
import { useXrpAmmLiveData } from '~/composables/useXrpAmmLiveData'
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

export default function useXrpAmmSwap(
  fromToken: any,
  toToken: any,
  fromAmount: any,
  pool: XrpAmmPool
) {
  const { $f } = useContext()
  const { 
    wallet, 
    isWalletReady, 
    signAMMTradeTransaction,
    canSignTransactions,
    address
  } = useEnhancedXrpWallet()
  
  const { getQuote, getUserTokenBalances, refreshAll } = useXrpAmmLiveData()
  const { submitAndConfirm, submitting, confirming, error: txError, clearTransactionState } = useXrplTransaction()
  
  // State
  const loading = ref(false)
  const quote = ref('')
  const gasFeeUSD = ref(0.01)
  const errorMessage = ref('')
  const txLoading = ref(false)
  const receipt = ref<any>(null)
  const isTxMined = ref(false)
  const txHash = ref<string | null>(null)
  
  // Transaction status tracking
  const txStatus = ref<'idle' | 'pending' | 'success' | 'error'>('idle')
  
  // Get user token balances
  const { tokenBalances, loading: balancesLoading } = getUserTokenBalances()
  
  // Get live quote when amount changes
  const { quote: liveQuote, loading: quoteLoading, error: quoteError } = getQuote(
    pool.id,
    computed(() => fromAmount.value?.toString() || '0'),
    computed(() => fromToken.value?.symbol || '')
  )
  
  // Computed balances
  const fromTokenBalance = computed(() => {
    if (!fromToken.value) return 0
    const balance = tokenBalances.value.find(b => 
      b.currency === fromToken.value.symbol && 
      b.issuer === fromToken.value.issuer
    )
    return balance?.balance || 0
  })
  
  const toTokenBalance = computed(() => {
    if (!toToken.value) return 0
    const balance = tokenBalances.value.find(b => 
      b.currency === toToken.value.symbol && 
      b.issuer === toToken.value.issuer
    )
    return balance?.balance || 0
  })
  
  // Update quote from live data
  watch(liveQuote, (newQuote) => {
    if (newQuote && fromAmount.value > 0) {
      quote.value = `1 ${fromToken.value?.symbol} = ${newQuote.price.toFixed(6)} ${toToken.value?.symbol}`
    } else {
      quote.value = ''
    }
  })
  
  // Watch for quote errors
  watch(quoteError, (err) => {
    if (err) {
      errorMessage.value = 'Failed to get quote'
    } else {
      errorMessage.value = ''
    }
  })
  
  // Validate transaction parameters
  const validateTransaction = (amount: number) => {
    if (!isWalletReady.value) {
      throw new Error('Wallet not connected')
    }
    
    if (!canSignTransactions.value) {
      throw new Error('Wallet does not support transaction signing')
    }
    
    if (!amount || amount <= 0) {
      throw new Error('Invalid amount')
    }
    
    if (amount > fromTokenBalance.value) {
      throw new Error('Insufficient balance')
    }
    
    if (!address.value) {
      throw new Error('No wallet address available')
    }
    
    if (!liveQuote.value) {
      throw new Error('No quote available')
    }
  }
  
  // Action button state
  const actionButton = computed(() => {
    if (!isWalletReady.value) {
      return { status: false, message: 'Connect Wallet' }
    }
    
    if (!canSignTransactions.value) {
      return { status: false, message: 'Wallet not supported' }
    }
    
    if (!fromAmount.value || fromAmount.value <= 0) {
      return { status: false, message: 'Enter an amount' }
    }
    
    if (fromAmount.value > fromTokenBalance.value) {
      return { status: false, message: 'Insufficient balance' }
    }
    
    if (quoteLoading.value) {
      return { status: false, message: 'Getting quote...' }
    }
    
    if (errorMessage.value) {
      return { status: false, message: 'Quote error' }
    }
    
    if (txLoading.value || submitting.value || confirming.value) {
      return { status: false, message: 'Swapping...' }
    }
    
    if (!liveQuote.value) {
      return { status: false, message: 'No quote available' }
    }
    
    return { status: true, message: 'Swap' }
  })
  
  // Execute swap
  const swap = async () => {
    validateTransaction(fromAmount.value)
    
    txLoading.value = true
    txStatus.value = 'pending'
    errorMessage.value = ''
    receipt.value = null
    txHash.value = null
    
    try {
      // Create AMM trade transaction
      const tradeParams = {
        amount: fromAmount.value.toString(),
        amount2Currency: toToken.value.symbol,
        amount2Issuer: toToken.value.issuer,
        amount2Value: liveQuote.value.outputAmount
      }
      
      console.log('Creating AMM trade transaction:', tradeParams)
      
      // Sign transaction with wallet
      const signedTx = await signAMMTradeTransaction(tradeParams)
      
      // Submit and confirm transaction
      const txReceipt = await submitAndConfirm(signedTx)
      
      // Update state
      txHash.value = txReceipt.hash
      receipt.value = {
        ...txReceipt,
        to: pool.id,
        poolId: pool.id,
        inputAmount: fromAmount.value,
        outputAmount: liveQuote.value.outputAmount,
        priceImpact: liveQuote.value.priceImpact
      }
      
      isTxMined.value = true
      txStatus.value = 'success'
      
      // Refresh data
      await refreshAll()
      
      console.log(`Successfully swapped ${fromAmount.value} ${fromToken.value.symbol} for ${liveQuote.value.outputAmount} ${toToken.value.symbol}`)
      
    } catch (err: any) {
      errorMessage.value = err.message || 'Swap failed'
      txStatus.value = 'error'
      console.error('Swap error:', err)
    } finally {
      txLoading.value = false
    }
  }
  
  // Clear trade state
  const clearTrade = () => {
    loading.value = false
    quote.value = ''
    errorMessage.value = ''
    txLoading.value = false
    receipt.value = null
    isTxMined.value = false
    txHash.value = null
    txStatus.value = 'idle'
    clearTransactionState()
  }
  
  // Watch for amount changes to recalculate quote
  watch(fromAmount, () => {
    if (fromAmount.value > 0) {
      // Quote will be updated automatically via the GraphQL query
    } else {
      quote.value = ''
    }
  })
  
  // Watch for token changes to clear state
  watch([fromToken, toToken], () => {
    clearTrade()
  })
  
  // Combined loading state
  const isLoading = computed(() => 
    quoteLoading.value || 
    balancesLoading.value || 
    loading.value || 
    txLoading.value || 
    submitting.value || 
    confirming.value
  )
  
  // Combined error state
  const hasError = computed(() => 
    !!errorMessage.value || 
    !!txError.value
  )
  
  return {
    // State
    loading: isLoading,
    quote,
    gasFeeUSD,
    errorMessage: computed(() => errorMessage.value || txError.value),
    txLoading,
    receipt,
    isTxMined,
    txHash,
    txStatus: computed(() => txStatus.value),
    fromTokenBalance,
    toTokenBalance,
    
    // Computed
    actionButton,
    
    // Live data
    liveQuote,
    
    // Methods
    swap,
    clearTrade,
  }
} 