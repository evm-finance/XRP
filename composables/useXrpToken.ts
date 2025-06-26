import { computed, ref, useRoute } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPScreenerGQL, XRPAccountTransactionsGQL, XRPAccountBalancesGQL } from '~/apollo/queries'

interface XRPTokenData {
  currency: string
  issuerAddress: string
  tokenName: string
  issuerName: string
  price: number
  marketcap: number
  volume24H: number
  websiteUrl?: string
  sourceCodeUrl?: string
  telegramUrl?: string
  twitterUrl?: string
  discordUrl?: string
}

interface XRPBalance {
  account: string
  balance: number
  value: number
}

interface XRPTransaction {
  hash: string
  account: string
  destination: string
  transactionType: string
  amount: number
  currency?: string
  fee: number
  date: string
}

interface XRPAMMData {
  time: string
  price: number
  volume: number
}

export default function useXrpToken() {
  const route = useRoute()
  const loading = ref(true)
  
  // Get token currency and issuer from route
  const tokenCurrency = computed(() => route.value.params?.id ?? '')
  const issuerAddress = computed(() => route.value.query?.issuer as string ?? '')
  
  // State
  const selectedTimeframe = ref('1d')
  const priceDisplayMode = ref('xrp')
  const xrpPrice = ref(0.5) // Mock XRP price, would come from API
  
  // Token data
  const tokenData = ref<XRPTokenData>({
    currency: tokenCurrency.value,
    issuerAddress: issuerAddress.value,
    tokenName: '',
    issuerName: '',
    price: 0,
    marketcap: 0,
    volume24H: 0,
  })
  
  // Balances and transactions
  const tokenBalances = ref<XRPBalance[]>([])
  const walletTransactions = ref<XRPTransaction[]>([])
  const ammChartData = ref<XRPAMMData[]>([])
  
  // Queries
  const { onResult: onScreenerResult } = useQuery(XRPScreenerGQL, { 
    fetchPolicy: 'no-cache', 
    pollInterval: 60000 
  })
  
  const { onResult: onTransactionsResult } = useQuery(
    XRPAccountTransactionsGQL,
    () => ({ address: issuerAddress.value }),
    { fetchPolicy: 'no-cache', pollInterval: 30000 }
  )
  
  const { onResult: onBalancesResult } = useQuery(
    XRPAccountBalancesGQL,
    () => ({ account: issuerAddress.value }),
    { fetchPolicy: 'no-cache', pollInterval: 30000 }
  )
  
  // Computed
  const transactionHeaders = computed(() => [
    { text: 'Hash', value: 'hash', width: '200' },
    { text: 'Type', value: 'transactionType', width: '120' },
    { text: 'Amount', value: 'amount', width: '120' },
    { text: 'Fee', value: 'fee', width: '80' },
    { text: 'Date', value: 'date', width: '120' },
  ])
  
  // Methods
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      // Could add toast notification here
    } catch (err) {
      console.error('Failed to copy text: ', err)
    }
  }
  
  const formatAddress = (address: string) => {
    return `${address.slice(0, 10)}...${address.slice(-10)}`
  }
  
  const formatHash = (hash: string) => {
    return `${hash.slice(0, 10)}...${hash.slice(-10)}`
  }
  
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 6,
      maximumFractionDigits: 6,
    }).format(price)
  }
  
  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount)
  }
  
  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString()
  }
  
  const getTransactionTypeColor = (type: string) => {
    switch (type) {
      case 'Payment': return 'green'
      case 'OfferCreate': return 'blue'
      case 'OfferCancel': return 'orange'
      case 'TrustSet': return 'purple'
      default: return 'grey'
    }
  }
  
  // Event handlers
  onScreenerResult((queryResult: any) => {
    const screenerData = queryResult.data?.xrpScreener ?? []
    const tokenInfo = screenerData.find((t: any) => 
      t.currency === tokenCurrency.value && t.issuerAddress === issuerAddress.value
    )
    
    if (tokenInfo) {
      tokenData.value = {
        ...tokenData.value,
        tokenName: tokenInfo.tokenName || `${tokenCurrency.value} Token`,
        issuerName: tokenInfo.issuerName || 'Unknown Issuer',
        price: tokenInfo.price || 0,
        marketcap: tokenInfo.marketcap || 0,
        volume24H: tokenInfo.volume24H || 0,
        websiteUrl: tokenInfo.websiteUrl,
        sourceCodeUrl: tokenInfo.sourceCodeUrl,
        telegramUrl: tokenInfo.telegramUrl,
        twitterUrl: tokenInfo.twitterUrl,
        discordUrl: tokenInfo.discordUrl,
      }
    } else {
      // Set default values if token not found in screener
      tokenData.value = {
        ...tokenData.value,
        tokenName: `${tokenCurrency.value} Token`,
        issuerName: 'Unknown Issuer',
        price: 0,
        marketcap: 0,
        volume24H: 0,
      }
    }
    
    loading.value = false
  })
  
  onTransactionsResult((queryResult: any) => {
    const transactions = queryResult.data?.xrpTransactions ?? []
    walletTransactions.value = transactions.map((tx: any) => ({
      hash: tx.hash || '',
      account: tx.account || '',
      destination: tx.destination || '',
      transactionType: tx.transactionType || 'Unknown',
      amount: typeof tx.amount === 'string' ? parseFloat(tx.amount) : (tx.amount || 0),
      currency: tokenCurrency.value,
      fee: typeof tx.fee === 'string' ? parseFloat(tx.fee) : (tx.fee || 0),
      date: tx.date || new Date().toISOString(),
    }))
  })
  
  onBalancesResult((queryResult: any) => {
    const balances = queryResult.data?.xrpAccountBalances?.xrpTokens ?? []
    tokenBalances.value = balances
      .filter((token: any) => token.symbol === tokenCurrency.value)
      .map((token: any) => ({
        account: token.issuer || '',
        balance: token.balance || 0,
        value: token.value || 0,
      }))
  })

  // Add error handling for all queries
  const { onError: onScreenerError } = useQuery(XRPScreenerGQL, { 
    fetchPolicy: 'no-cache', 
    pollInterval: 60000 
  })

  const { onError: onTransactionsError } = useQuery(
    XRPAccountTransactionsGQL,
    () => ({ address: issuerAddress.value }),
    { fetchPolicy: 'no-cache', pollInterval: 30000 }
  )

  const { onError: onBalancesError } = useQuery(
    XRPAccountBalancesGQL,
    () => ({ account: issuerAddress.value }),
    { fetchPolicy: 'no-cache', pollInterval: 30000 }
  )

  onScreenerError((error: any) => {
    console.error('GraphQL Error in useXrpToken screener:', error)
    loading.value = false
  })

  onTransactionsError((error: any) => {
    console.error('GraphQL Error in useXrpToken transactions:', error)
  })

  onBalancesError((error: any) => {
    console.error('GraphQL Error in useXrpToken balances:', error)
  })

  // Remove mock data generation - rely entirely on live data
  // generateMockData() - REMOVED
  loading.value = false
  
  return {
    loading,
    tokenData,
    tokenBalances,
    walletTransactions,
    ammChartData,
    selectedTimeframe,
    priceDisplayMode,
    xrpPrice,
    transactionHeaders,
    copyToClipboard,
    formatAddress,
    formatHash,
    formatPrice,
    formatAmount,
    formatDate,
    getTransactionTypeColor,
  }
} 