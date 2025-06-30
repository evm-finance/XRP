import { computed } from '@nuxtjs/composition-api'

export interface XRPScreenerFallback {
  currency: string
  issuerAddress: string
  tokenName: string
  issuerName: string
  price: number
  marketcap: number
  volume24H: number
  icon: string
}

export interface XRPTransactionFallback {
  hash: string
  account: string
  destination?: string
  amount: string
  fee: string
  date: string
  type: string
  status: string
}

export interface XRPBalanceFallback {
  currency: string
  issuer?: string
  balance: string
  valueUsd: number
}

export interface XRPAMMFallback {
  poolId: string
  asset1: {
    currency: string
    issuer: string
  }
  asset2: {
    currency: string
    issuer: string
  }
  asset1ValueUsd: number
  asset2ValueUsd: number
  totalLiquidityUsd: number
  fee: number
  createdAt: string
}

export function useXrpFallbackData() {
  // XRP Screener fallback data
  const screenerFallback = computed<XRPScreenerFallback[]>(() => [
    {
      currency: 'USDC',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'USD Coin',
      issuerName: 'Circle',
      price: 1.0,
      marketcap: 1000000,
      volume24H: 500000,
      icon: '🪙'
    },
    {
      currency: 'USDT',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'Tether',
      issuerName: 'Tether',
      price: 1.0,
      marketcap: 800000,
      volume24H: 400000,
      icon: '💎'
    },
    {
      currency: 'BTC',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'Bitcoin',
      issuerName: 'BitGo',
      price: 45000,
      marketcap: 500000,
      volume24H: 200000,
      icon: '₿'
    },
    {
      currency: 'ETH',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'Ethereum',
      issuerName: 'BitGo',
      price: 3000,
      marketcap: 300000,
      volume24H: 150000,
      icon: 'Ξ'
    },
    {
      currency: 'SOL',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'Solana',
      issuerName: 'Solana',
      price: 100,
      marketcap: 200000,
      volume24H: 100000,
      icon: '◎'
    }
  ])

  // XRP Transactions fallback data
  const transactionsFallback = computed<XRPTransactionFallback[]>(() => [
    {
      hash: 'txn_1234567890abcdef',
      account: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      destination: 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo',
      amount: '1000',
      fee: '0.000012',
      date: new Date(Date.now() - 3600000).toISOString(),
      type: 'Payment',
      status: 'tesSUCCESS'
    },
    {
      hash: 'txn_abcdef1234567890',
              account: 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo',
      destination: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      amount: '500',
      fee: '0.000012',
      date: new Date(Date.now() - 7200000).toISOString(),
      type: 'Payment',
      status: 'tesSUCCESS'
    },
    {
      hash: 'txn_7890abcdef123456',
      account: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      amount: '100',
      fee: '0.000012',
      date: new Date(Date.now() - 10800000).toISOString(),
      type: 'TrustSet',
      status: 'tesSUCCESS'
    }
  ])

  // XRP Balances fallback data
  const balancesFallback = computed<XRPBalanceFallback[]>(() => [
    {
      currency: 'XRP',
      balance: '10000',
      valueUsd: 5000
    },
    {
      currency: 'USDC',
      issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      balance: '1000',
      valueUsd: 1000
    },
    {
      currency: 'USDT',
      issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      balance: '500',
      valueUsd: 500
    },
    {
      currency: 'BTC',
      issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      balance: '0.1',
      valueUsd: 4500
    }
  ])

  // XRP AMM fallback data
  const ammFallback = computed<XRPAMMFallback[]>(() => [
    {
      poolId: '1',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'USD', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      asset1ValueUsd: 5000000,
      asset2ValueUsd: 5000000,
      totalLiquidityUsd: 10000000,
      fee: 0.3,
      createdAt: new Date().toISOString()
    },
    {
      poolId: '2',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'BTC', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      asset1ValueUsd: 3000000,
      asset2ValueUsd: 7000000,
      totalLiquidityUsd: 10000000,
      fee: 0.5,
      createdAt: new Date().toISOString()
    },
    {
      poolId: '3',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'ETH', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      asset1ValueUsd: 2000000,
      asset2ValueUsd: 8000000,
      totalLiquidityUsd: 10000000,
      fee: 0.4,
      createdAt: new Date().toISOString()
    },
    {
      poolId: '4',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'SOL', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      asset1ValueUsd: 1500000,
      asset2ValueUsd: 3500000,
      totalLiquidityUsd: 5000000,
      fee: 0.6,
      createdAt: new Date().toISOString()
    },
    {
      poolId: '5',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'ADA', issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B' },
      asset1ValueUsd: 800000,
      asset2ValueUsd: 1200000,
      totalLiquidityUsd: 2000000,
      fee: 0.7,
      createdAt: new Date().toISOString()
    }
  ])

  // XRP Token Mints fallback data
  const tokenMintsFallback = computed(() => [
    {
      currency: 'NEWTOKEN',
      issuerAddress: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      tokenName: 'New Token',
      issuerName: 'New Issuer',
      totalSupply: '1000000',
      createdAt: new Date(Date.now() - 86400000).toISOString(),
      icon: '🆕'
    },
    {
      currency: 'TESTTOKEN',
              issuerAddress: 'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo',
      tokenName: 'Test Token',
      issuerName: 'Test Issuer',
      totalSupply: '500000',
      createdAt: new Date(Date.now() - 172800000).toISOString(),
      icon: '🧪'
    }
  ])

  // XRP Liquidity Pools fallback data
  const liquidityPoolsFallback = computed(() => [
    {
      poolId: 'pool_1',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'USDC', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      totalLiquidityUsd: 5000000,
      volume24H: 1000000,
      fee: 0.3,
      createdAt: new Date().toISOString()
    },
    {
      poolId: 'pool_2',
      asset1: { currency: 'XRP', issuer: '' },
      asset2: { currency: 'USDT', issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh' },
      totalLiquidityUsd: 3000000,
      volume24H: 800000,
      fee: 0.3,
      createdAt: new Date().toISOString()
    }
  ])

  // Get fallback data by query type
  const getFallbackData = (queryType: string) => {
    switch (queryType) {
      case 'XRPScreener':
        return { xrpScreener: screenerFallback.value }
      case 'XRPTransactions':
      case 'XRPAccountTransactions':
        return { xrpTransactions: transactionsFallback.value }
      case 'XRPBalances':
      case 'XRPAccountBalances':
        return { xrpBalances: balancesFallback.value }
      case 'SimpleAMMLiquidityValues':
      case 'GetAllAMMLiquidityValues':
        return { xrpAMMLiquidityValues: ammFallback.value }
      case 'XRPTokenMints':
        return { xrpTokenMints: tokenMintsFallback.value }
      case 'XRPLiquidityPools':
        return { xrpLiquidityPools: liquidityPoolsFallback.value }
      case 'XRPDefiData':
        return { xrpDefiData: { balances: balancesFallback.value, transactions: transactionsFallback.value } }
      default:
        return null
    }
  }

  return {
    screenerFallback,
    transactionsFallback,
    balancesFallback,
    ammFallback,
    tokenMintsFallback,
    liquidityPoolsFallback,
    getFallbackData
  }
} 