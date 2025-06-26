import { ref, computed, inject } from '@nuxtjs/composition-api'
import { XRP_PLUGIN_KEY, XrpClient } from '~/plugins/web3/xrp.client'
import { setTrustline } from '@gemwallet/api'

export interface TrustLine {
  currency: string
  issuer: string
  limit: number
  balance: number
  flags: string[]
}

export interface PopularToken {
  currency: string
  name: string
  issuer: string
}

export interface TrustLineForm {
  currency: string
  issuer: string
  limit: string
  flags: string[]
}

export default function useXrpTrustLines() {
  const xrpClient = inject(XRP_PLUGIN_KEY) as XrpClient
  
  // State
  const loading = ref(false)
  const trustLines = ref<TrustLine[]>([])
  const editingTrustLine = ref<TrustLine | null>(null)
  const deletingTrustLine = ref<TrustLine | null>(null)
  
  // Form data
  const newTrustLine = ref<TrustLineForm>({
    currency: '',
    issuer: '',
    limit: '',
    flags: []
  })

  // Computed
  const isWalletReady = computed(() => xrpClient.isWalletReady.value)
  const address = computed(() => xrpClient.address.value)

  // Trust line flags
  const trustLineFlags = [
    { label: 'Set Authorization', value: 'tfSetfAuth' },
    { label: 'Set No Ripple', value: 'tfSetNoRipple' },
    { label: 'Clear No Ripple', value: 'tfClearNoRipple' },
    { label: 'Set Freeze', value: 'tfSetFreeze' },
    { label: 'Clear Freeze', value: 'tfClearFreeze' }
  ]

  // Popular tokens
  const popularTokens: PopularToken[] = [
    {
      currency: 'USDC',
      name: 'USD Coin',
      issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q'
    },
    {
      currency: 'USDT',
      name: 'Tether USD',
      issuer: 'rchGBxcD1A1C2tdxF6papQYZ8kjRKMYcL'
    },
    {
      currency: 'BTC',
      name: 'Bitcoin',
      issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B'
    },
    {
      currency: 'ETH',
      name: 'Ethereum',
      issuer: 'rcA8X3TVMST1n3CJeAdGk1RdRCHii7N2h'
    }
  ]

  // Table headers
  const trustLineHeaders = [
    { text: 'Currency', value: 'currency', width: '150' },
    { text: 'Issuer', value: 'issuer', width: '200' },
    { text: 'Limit', value: 'limit', width: '120' },
    { text: 'Balance', value: 'balance', width: '120' },
    { text: 'Flags', value: 'flags', width: '150' },
    { text: 'Actions', value: 'actions', width: '100', sortable: false }
  ]

  // Form validation rules
  const rules = {
    required: (v: any) => !!v || 'This field is required',
    currency: (v: string) => /^[A-Z]{3,4}$/.test(v) || 'Currency must be 3-4 uppercase letters',
    address: (v: string) => /^r[a-zA-Z0-9]{25,34}$/.test(v) || 'Invalid XRP address format',
    positive: (v: number) => v > 0 || 'Value must be positive'
  }

  // Methods
  const connectWallet = () => {
    xrpClient.connectWallet()
  }

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

  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 6,
    }).format(amount)
  }

  const loadTrustLines = async () => {
    if (!isWalletReady.value) return
    
    loading.value = true
    try {
      // Mock data for development - would be replaced with actual API call
      trustLines.value = [
        {
          currency: 'USDC',
          issuer: 'rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q',
          limit: 10000,
          balance: 500,
          flags: ['tfSetfAuth']
        },
        {
          currency: 'BTC',
          issuer: 'rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B',
          limit: 100,
          balance: 0,
          flags: []
        }
      ]
    } catch (error) {
      console.error('Failed to load trust lines:', error)
    } finally {
      loading.value = false
    }
  }

  const createTrustLine = async () => {
    if (!isWalletReady.value) return
    
    loading.value = true
    try {
      const payload = {
        limitAmount: {
          currency: newTrustLine.value.currency,
          issuer: newTrustLine.value.issuer,
          value: newTrustLine.value.limit.toString()
        },
        flags: newTrustLine.value.flags.reduce((acc, flag) => {
          acc[flag] = true
          return acc
        }, {} as any)
      }

      const response = await setTrustline(payload)
      console.log('Trust line created:', response)
      
      // Reset form
      newTrustLine.value = {
        currency: '',
        issuer: '',
        limit: '',
        flags: []
      }
      
      // Reload trust lines
      await loadTrustLines()
    } catch (error) {
      console.error('Failed to create trust line:', error)
    } finally {
      loading.value = false
    }
  }

  const quickCreateTrustLine = (token: PopularToken) => {
    newTrustLine.value.currency = token.currency
    newTrustLine.value.issuer = token.issuer
    newTrustLine.value.limit = '1000'
    newTrustLine.value.flags = []
  }

  const editTrustLine = (trustLine: TrustLine) => {
    editingTrustLine.value = { ...trustLine }
  }

  const updateTrustLine = async () => {
    if (!editingTrustLine.value) return
    
    loading.value = true
    try {
      // Mock update - would be replaced with actual API call
      const index = trustLines.value.findIndex(tl => 
        tl.currency === editingTrustLine.value?.currency && 
        tl.issuer === editingTrustLine.value?.issuer
      )
      
      if (index !== -1) {
        trustLines.value[index] = { ...editingTrustLine.value }
      }
      
      editingTrustLine.value = null
    } catch (error) {
      console.error('Failed to update trust line:', error)
    } finally {
      loading.value = false
    }
  }

  const deleteTrustLine = (trustLine: TrustLine) => {
    deletingTrustLine.value = trustLine
  }

  const confirmDeleteTrustLine = async () => {
    if (!deletingTrustLine.value) return
    
    loading.value = true
    try {
      // Mock delete - would be replaced with actual API call
      trustLines.value = trustLines.value.filter(tl => 
        !(tl.currency === deletingTrustLine.value?.currency && 
          tl.issuer === deletingTrustLine.value?.issuer)
      )
      
      deletingTrustLine.value = null
    } catch (error) {
      console.error('Failed to delete trust line:', error)
    } finally {
      loading.value = false
    }
  }

  const refreshTrustLines = () => {
    loadTrustLines()
  }

  return {
    // State
    loading,
    trustLines,
    editingTrustLine,
    deletingTrustLine,
    newTrustLine,
    
    // Computed
    isWalletReady,
    address,
    
    // Data
    trustLineFlags,
    popularTokens,
    trustLineHeaders,
    rules,
    
    // Methods
    connectWallet,
    copyToClipboard,
    formatAddress,
    formatAmount,
    createTrustLine,
    quickCreateTrustLine,
    editTrustLine,
    updateTrustLine,
    deleteTrustLine,
    confirmDeleteTrustLine,
    refreshTrustLines
  }
} 