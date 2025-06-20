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
  
  // Mock data generation for development
  const generateMockData = () => {
    // Mock token data
    tokenData.value = {
      ...tokenData.value,
      tokenName: `${tokenCurrency.value} Token`,
      issuerName: 'Sample Issuer',
      price: 0.001234,
      marketcap: 1234567,
      volume24H: 98765,
      websiteUrl: 'https://example.com',
      sourceCodeUrl: 'https://github.com/example',
      telegramUrl: 'https://t.me/example',
      twitterUrl: 'https://twitter.com/example',
      discordUrl: 'https://discord.gg/example',
    }
    
    // Mock balances
    tokenBalances.value = [
      {
        account: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        balance: 1000000,
        value: 1234.56,
      },
      {
        account: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
        balance: 500000,
        value: 617.28,
      },
    ]
    
    // Mock transactions
    walletTransactions.value = [
      {
        hash: 'E0C1D4B24D76B4180D2C96450438A0BE14304E69EDFBE91DF6211C923B344401',
        account: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        destination: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
        transactionType: 'Payment',
        amount: 1000000,
        currency: tokenCurrency.value,
        fee: 12,
        date: new Date().toISOString(),
      },
      {
        hash: 'F1D2E5C35E87C5291E3D07561549B1CF25415F7AFEGC02EG7322D034C455502',
        account: 'rYYYYYYYYYYYYYYYYYYYYYYYYYYYY',
        destination: 'rXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        transactionType: 'OfferCreate',
        amount: 500000,
        currency: tokenCurrency.value,
        fee: 12,
        date: new Date(Date.now() - 86400000).toISOString(),
      },
    ]
    
    // Mock AMM chart data
    ammChartData.value = Array.from({ length: 24 }, (_, i) => ({
      time: new Date(Date.now() - (23 - i) * 3600000).toISOString(),
      price: 0.001 + Math.random() * 0.0005,
      volume: Math.random() * 1000,
    }))
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
        tokenName: tokenInfo.tokenName,
        issuerName: tokenInfo.issuerName,
        price: tokenInfo.price,
        marketcap: tokenInfo.marketcap,
        volume24H: tokenInfo.volume24H,
      }
    }
    
    loading.value = false
  })
  
  onTransactionsResult((queryResult: any) => {
    const transactions = queryResult.data?.xrpTransactions ?? []
    walletTransactions.value = transactions.map((tx: any) => ({
      ...tx,
      currency: tokenCurrency.value,
    }))
  })
  
  onBalancesResult((queryResult: any) => {
    const balances = queryResult.data?.xrpAccountBalances?.xrpTokens ?? []
    tokenBalances.value = balances
      .filter((token: any) => token.symbol === tokenCurrency.value)
      .map((token: any) => ({
        account: token.issuer,
        balance: token.balance,
        value: token.value,
      }))
  })
  
  // Initialize with mock data for development
  generateMockData()
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