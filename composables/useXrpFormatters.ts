import { useContext } from '@nuxtjs/composition-api'

export interface XrpFormatterOptions {
  currency?: string
  locale?: string
  maxDigits?: number
  minDigits?: number
  autoDigits?: boolean
}

export function useXrpFormatters() {
  const { $f } = useContext()

  /**
   * Format XRP price with appropriate decimal places
   */
  const formatXrpPrice = (value: number, options: XrpFormatterOptions = {}) => {
    const { currency = 'USD', locale = 'en-US', autoDigits = true } = options
    
    if (!value || isNaN(value)) {
      return '-.--'
    }

    // XRP-specific formatting logic
    if (value >= 10) {
      return new Intl.NumberFormat(locale, {
        currency,
        style: 'currency',
        maximumFractionDigits: 2,
        minimumFractionDigits: 2,
      }).format(value)
    } else if (value <= 0.0001) {
      return new Intl.NumberFormat(locale, {
        currency,
        style: 'currency',
        maximumFractionDigits: 6,
        minimumFractionDigits: 6,
      }).format(value)
    } else {
      return new Intl.NumberFormat(locale, {
        currency,
        style: 'currency',
        maximumFractionDigits: 4,
        minimumFractionDigits: 4,
      }).format(value)
    }
  }

  /**
   * Format issuer address with copy functionality
   */
  const formatIssuerAddress = (address: string, maxLength: number = 12) => {
    if (!address) return '-'
    
    if (address.length <= maxLength) {
      return address
    }
    
    const start = address.substring(0, maxLength / 2)
    const end = address.substring(address.length - maxLength / 2)
    return `${start}...${end}`
  }

  /**
   * Format trust line information
   */
  const formatTrustLine = (limit: number, balance: number, currency: string = 'USD') => {
    if (limit === 0) return 'No Trust Line'
    
    const limitFormatted = formatXrpPrice(limit, { currency })
    const balanceFormatted = formatXrpPrice(balance, { currency })
    const percentage = limit > 0 ? (balance / limit) * 100 : 0
    
    return {
      limit: limitFormatted,
      balance: balanceFormatted,
      percentage: `${percentage.toFixed(2)}%`,
      status: balance > 0 ? 'Active' : 'Inactive'
    }
  }

  /**
   * Format AMM pool information
   */
  const formatAmmPool = (pool: any) => {
    if (!pool) return '-'
    
    return {
      asset1: pool.asset1?.currency || 'XRP',
      asset2: pool.asset2?.currency || 'USD',
      liquidity: formatXrpPrice(pool.liquidity || 0),
      volume24h: formatXrpPrice(pool.volume24h || 0),
      fee: `${(pool.fee || 0) * 100}%`
    }
  }

  /**
   * Format percentage change with color coding
   */
  const formatPercentageChange = (value: number, options: XrpFormatterOptions = {}) => {
    const { maxDigits = 2, minDigits = 2 } = options
    
    if (!value || isNaN(value)) {
      return '-.--%'
    }

    const sign = value >= 0 ? '+' : ''
    const formatted = $f(value * 100, { 
      maxDigits, 
      minDigits, 
      after: '%' 
    })
    
    return `${sign}${formatted}`
  }

  /**
   * Format market cap with appropriate units
   */
  const formatMarketCap = (value: number, currency: string = 'USD') => {
    if (!value || isNaN(value)) {
      return '-'
    }

    if (value >= 1e12) {
      return `${(value / 1e12).toFixed(2)}T ${currency}`
    } else if (value >= 1e9) {
      return `${(value / 1e9).toFixed(2)}B ${currency}`
    } else if (value >= 1e6) {
      return `${(value / 1e6).toFixed(2)}M ${currency}`
    } else if (value >= 1e3) {
      return `${(value / 1e3).toFixed(2)}K ${currency}`
    } else {
      return formatXrpPrice(value, { currency })
    }
  }

  /**
   * Format volume with appropriate units
   */
  const formatVolume = (value: number, currency: string = 'USD') => {
    return formatMarketCap(value, currency)
  }

  /**
   * Format transaction hash
   */
  const formatTransactionHash = (hash: string, maxLength: number = 16) => {
    if (!hash) return '-'
    
    if (hash.length <= maxLength) {
      return hash
    }
    
    const start = hash.substring(0, maxLength / 2)
    const end = hash.substring(hash.length - maxLength / 2)
    return `${start}...${end}`
  }

  /**
   * Format ledger index
   */
  const formatLedgerIndex = (index: number) => {
    if (!index || isNaN(index)) return '-'
    return index.toLocaleString()
  }

  /**
   * Format timestamp
   */
  const formatTimestamp = (timestamp: string | number, format: 'relative' | 'absolute' = 'relative') => {
    if (!timestamp) return '-'
    
    const date = new Date(timestamp)
    
    if (format === 'relative') {
      const now = new Date()
      const diff = now.getTime() - date.getTime()
      const seconds = Math.floor(diff / 1000)
      const minutes = Math.floor(seconds / 60)
      const hours = Math.floor(minutes / 60)
      const days = Math.floor(hours / 24)
      
      if (days > 0) return `${days}d ago`
      if (hours > 0) return `${hours}h ago`
      if (minutes > 0) return `${minutes}m ago`
      return `${seconds}s ago`
    } else {
      return date.toLocaleString()
    }
  }

  /**
   * Format XRP amount (drops to XRP)
   */
  const formatXrpAmount = (drops: number, options: XrpFormatterOptions = {}) => {
    const { maxDigits = 6, minDigits = 0 } = options
    
    if (!drops || isNaN(drops)) {
      return '0 XRP'
    }
    
    const xrp = drops / 1000000 // 1 XRP = 1,000,000 drops
    return `${$f(xrp, { maxDigits, minDigits })} XRP`
  }

  return {
    formatXrpPrice,
    formatIssuerAddress,
    formatTrustLine,
    formatAmmPool,
    formatPercentageChange,
    formatMarketCap,
    formatVolume,
    formatTransactionHash,
    formatLedgerIndex,
    formatTimestamp,
    formatXrpAmount
  }
} 