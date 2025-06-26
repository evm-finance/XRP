import { ref, computed } from '@nuxtjs/composition-api'

export const useXrpMockData = () => {
  // Mock XRP account data
  const mockAccountData = ref({
    account: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
    xrpBalance: 1000000000, // 1000 XRP in drops
    xrpPrice: 0.5,
    xrpTokens: [
      {
        symbol: 'USDC',
        issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        name: 'USD Coin',
        balance: 1000,
        price: 1.0,
        value: 1000
      },
      {
        symbol: 'BTC',
        issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        name: 'Bitcoin',
        balance: 0.5,
        price: 50000,
        value: 25000
      }
    ]
  })

  // Mock XRP transaction data
  const mockTransactions = ref([
    {
      amount: '1000000', // 1 XRP in drops
      destination: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      transactionType: 'Payment',
      fee: '12',
      hash: 'A1B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456',
      ledgerIndex: 12345678
    },
    {
      amount: '500000', // 0.5 XRP in drops
      destination: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
      transactionType: 'AMMTrade',
      fee: '15',
      hash: 'B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456A1',
      ledgerIndex: 12345679
    }
  ])

  // Mock XRP addresses for testing
  const mockAddresses = [
    'rMV5cxLAKs8SuoZ8Ly8geDSnXgf9gui6Fo',
    'rDodqfAoF8pVh2SoUwhQRfvkqrs4wwxUrz',
    'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
    'rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe',
    'rUCzEr6jrEyMpjhs4wSdQdz4g8Y382txf7'
  ]

  // Mock transaction hashes for testing
  const mockTransactionHashes = [
    'A1B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456',
    'B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456A1',
    'C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456A1B2'
  ]

  // Mock ledger indices for testing
  const mockLedgerIndices = [
    12345678,
    12345679,
    12345680,
    12345681
  ]

  // Get mock data for testing
  const getMockAccountData = (address: string) => {
    if (mockAddresses.includes(address)) {
      return { ...mockAccountData.value, account: address }
    }
    return null
  }

  const getMockTransactions = (address: string) => {
    if (mockAddresses.includes(address)) {
      return mockTransactions.value
    }
    return []
  }

  const getMockTransaction = (hash: string) => {
    if (mockTransactionHashes.includes(hash)) {
      return {
        account: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        amount: '1000000',
        destination: 'rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe',
        fee: '12',
        flags: 0,
        lastLedgerSequence: 12345678,
        offerSequence: 0,
        sequence: 1,
        signingPubKey: 'mock-signing-key',
        takerGets: '1000000',
        takerPays: '1000000',
        transactionType: 'Payment',
        txnSignature: 'mock-signature',
        date: new Date().toISOString(),
        hash: hash,
        inLedger: 12345678,
        ledgerIndex: 12345678,
        meta: {},
        metadata: {},
        validated: true,
        warnings: [],
        memos: []
      }
    }
    return null
  }

  const getMockLedger = (index: number) => {
    if (mockLedgerIndices.includes(index)) {
      return {
        ledgerIndex: index,
        ledgerHash: `mock-ledger-hash-${index}`,
        parentHash: `mock-parent-hash-${index}`,
        totalCoins: '100000000000000000',
        parentCloseTime: new Date().toISOString(),
        closeTime: new Date().toISOString(),
        closeTimeRes: 30,
        closeFlags: 0,
        accountHash: `mock-account-hash-${index}`,
        transactionHash: `mock-transaction-hash-${index}`,
        validated: true
      }
    }
    return null
  }

  return {
    mockAccountData,
    mockTransactions,
    mockAddresses,
    mockTransactionHashes,
    mockLedgerIndices,
    getMockAccountData,
    getMockTransactions,
    getMockTransaction,
    getMockLedger
  }
} 