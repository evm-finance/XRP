# Enhanced XRP Wallet Integration

This document describes the enhanced wallet integration that supports multiple XRP wallet types including GEM Wallet, Xaman (XUMM), and MetaMask XRP Snap.

## Supported Wallets

### 1. GEM Wallet
- **Type**: Browser Extension
- **Installation**: https://gemwallet.app/
- **Features**: Basic connection, address retrieval
- **Status**: ✅ Implemented

### 2. Xaman (XUMM)
- **Type**: Mobile Wallet with Web Integration
- **Installation**: https://xaman.app/
- **Features**: Full connection, transaction signing, account info
- **Status**: ✅ Implemented

### 3. MetaMask XRP Snap
- **Type**: MetaMask Extension + XRP Snap
- **Installation**: https://snaps.metamask.io/snap/npm/@metamask/xrp-snap/
- **Features**: Full connection, transaction signing, account info, balance
- **Status**: ✅ Implemented

## Architecture

### Core Components

#### 1. Wallet Connectors
- `plugins/web3/xaman.connector.ts` - Xaman wallet connector
- `plugins/web3/metamask-xrp-snap.connector.ts` - MetaMask XRP snap connector
- `plugins/web3/enhanced-xrp.client.ts` - Unified XRP client

#### 2. UI Components
- `components/common/EnhancedWalletSelectDialog.vue` - Enhanced wallet selection dialog
- `components/common/EnhancedXrpWalletConnector.vue` - Enhanced wallet connector component

#### 3. Composables
- `composables/useEnhancedXrpWallet.ts` - Easy-to-use wallet composable

## Usage

### Basic Wallet Connection

```typescript
import useEnhancedXrpWallet from '~/composables/useEnhancedXrpWallet'

export default {
  setup() {
    const {
      connectWallet,
      disconnectWallet,
      isWalletReady,
      address,
      currentWalletType,
      error
    } = useEnhancedXrpWallet()

    // Connect to specific wallet
    const connectToXaman = async () => {
      await connectWallet('xaman')
    }

    const connectToMetamaskXrp = async () => {
      await connectWallet('metamask-xrp-snap')
    }

    const connectToGem = async () => {
      await connectWallet('gem')
    }

    return {
      connectToXaman,
      connectToMetamaskXrp,
      connectToGem,
      disconnectWallet,
      isWalletReady,
      address,
      currentWalletType,
      error
    }
  }
}
```

### Transaction Signing

```typescript
import useEnhancedXrpWallet from '~/composables/useEnhancedXrpWallet'

export default {
  setup() {
    const {
      signPaymentTransaction,
      signTrustSetTransaction,
      signAMMTradeTransaction,
      canSignTransactions
    } = useEnhancedXrpWallet()

    // Sign a payment transaction
    const sendPayment = async () => {
      if (!canSignTransactions.value) {
        console.error('No compatible wallet connected')
        return
      }

      const result = await signPaymentTransaction({
        destination: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        amount: '1000000', // 1 XRP in drops
      })

      if (result.success) {
        console.log('Transaction signed:', result.signedTx)
      } else {
        console.error('Transaction failed:', result.error)
      }
    }

    // Sign a trust line transaction
    const setTrustLine = async () => {
      const result = await signTrustSetTransaction({
        currency: 'USDC',
        issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        limit: '1000000'
      })

      if (result.success) {
        console.log('Trust line set:', result.signedTx)
      }
    }

    // Sign an AMM trade
    const tradeAMM = async () => {
      const result = await signAMMTradeTransaction({
        amount: '1000000', // XRP amount
        amount2Currency: 'USDC',
        amount2Issuer: 'rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh',
        amount2Value: '1.0'
      })

      if (result.success) {
        console.log('AMM trade signed:', result.signedTx)
      }
    }

    return {
      sendPayment,
      setTrustLine,
      tradeAMM
    }
  }
}
```

### Account Information

```typescript
import useEnhancedXrpWallet from '~/composables/useEnhancedXrpWallet'

export default {
  setup() {
    const {
      getAccountInfo,
      getBalance,
      canGetAccountInfo,
      canGetBalance
    } = useEnhancedXrpWallet()

    // Get account information
    const fetchAccountInfo = async () => {
      if (!canGetAccountInfo.value) {
        console.error('Account info not available for this wallet')
        return
      }

      const result = await getAccountInfo()
      if (result.success) {
        console.log('Account info:', result.accountInfo)
      }
    }

    // Get balance (MetaMask XRP Snap only)
    const fetchBalance = async () => {
      if (!canGetBalance.value) {
        console.error('Balance not available for this wallet')
        return
      }

      const result = await getBalance()
      if (result.success) {
        console.log('Balance:', result.balance)
      }
    }

    return {
      fetchAccountInfo,
      fetchBalance
    }
  }
}
```

## Wallet Capabilities

| Feature | GEM Wallet | Xaman | MetaMask XRP Snap |
|---------|------------|-------|-------------------|
| Connection | ✅ | ✅ | ✅ |
| Address Retrieval | ✅ | ✅ | ✅ |
| Transaction Signing | ❌ | ✅ | ✅ |
| Account Info | ❌ | ✅ | ✅ |
| Balance | ❌ | ❌ | ✅ |
| AMM Operations | ❌ | ✅ | ✅ |

## Implementation Details

### Wallet Detection

Each wallet connector implements detection methods:

```typescript
// Xaman detection
private isXamanInstalled(): boolean {
  return typeof window !== 'undefined' && 
         (window as any).xumm !== undefined && 
         (window as any).xumm.xapp !== undefined
}

// MetaMask detection
private isMetamaskInstalled(): boolean {
  return typeof window !== 'undefined' && 
         (window as any).ethereum !== undefined &&
         (window as any).ethereum.isMetaMask === true
}
```

### Auto-Installation

The MetaMask XRP Snap connector automatically installs the snap if not present:

```typescript
private async installXrpSnap(): Promise<{ success: boolean; error?: string }> {
  try {
    const ethereum = (window as any).ethereum
    
    await ethereum.request({
      method: 'wallet_installSnaps',
      params: {
        [this.snapId]: {
          version: 'latest'
        }
      }
    })
    
    return { success: true }
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Failed to install XRP snap'
    }
  }
}
```

### Error Handling

All connectors implement comprehensive error handling:

```typescript
async connect(): Promise<{ account: string | null; error: Web3ErrorInterface }> {
  try {
    this.resetErrors()
    
    // Check if wallet is installed
    if (!this.isWalletInstalled()) {
      this.error = {
        status: true,
        message: 'Wallet not installed. Please install the wallet first.'
      }
      return { account: null, error: this.error }
    }

    // Attempt connection
    const result = await this.requestConnection()
    
    if (result.success && result.address) {
      this.account = result.address
      return { account: this.account, error: this.error }
    } else {
      this.error = {
        status: true,
        message: result.error || 'Failed to connect to wallet'
      }
      return { account: null, error: this.error }
    }
  } catch (err) {
    this.error = {
      status: true,
      message: 'Failed to connect to wallet'
    }
    return { account: null, error: this.error }
  }
}
```

## UI Integration

### Wallet Selection Dialog

The enhanced wallet selection dialog provides:

- Clear separation between Ethereum and XRP wallets
- Wallet-specific descriptions and icons
- Help links for wallet installation
- Real-time connection status
- Error handling and display

### Wallet Connector Component

The enhanced wallet connector shows:

- Current wallet type with appropriate icon
- Truncated address display
- Wallet-specific actions (copy address, view in explorer, etc.)
- Account info and balance (where supported)
- Disconnect functionality

## Configuration

### Plugin Registration

The enhanced XRP client is registered in `nuxt.config.js`:

```javascript
plugins: [
  '~/plugins/web3/web3.ts',
  '~/plugins/web3/xrp.client.ts',
  '~/plugins/web3/enhanced-xrp.client.ts', // Enhanced client
  '~/plugins/typer.client.ts',
],
```

### Cookie Management

Wallet connections are persisted using cookies:

```typescript
// Save connection
const inFifteenMinutes = new Date(new Date().getTime() + 15 * 60 * 1000)
context.$cookies.set(Cookies.expWalletConnected, walletType, { expires: inFifteenMinutes })

// Restore connection
const alreadyConnected: string | undefined = context.$cookies.get(Cookies.expWalletConnected)
if (alreadyConnected && ['gem', 'xaman', 'metamask-xrp-snap'].includes(alreadyConnected)) {
  await connectWallet(alreadyConnected as XrpWalletType)
}
```

## Testing

### Development Testing

1. **GEM Wallet**: Install the browser extension and test basic connection
2. **Xaman**: Use the mobile app with web integration
3. **MetaMask XRP Snap**: Install MetaMask and the XRP snap

### Test Scenarios

- [ ] Wallet installation detection
- [ ] Connection flow for each wallet type
- [ ] Transaction signing with different transaction types
- [ ] Error handling for various failure scenarios
- [ ] Auto-reconnection on page refresh
- [ ] Wallet switching
- [ ] Disconnection and cleanup

## Future Enhancements

### Planned Features

1. **GEM Wallet Enhancement**: Add transaction signing support
2. **Ledger Integration**: Add support for Ledger hardware wallets
3. **WalletConnect**: Add WalletConnect v2 support for mobile wallets
4. **Multi-wallet Support**: Allow multiple wallets to be connected simultaneously
5. **Transaction History**: Add transaction history retrieval
6. **Network Switching**: Support for testnet/mainnet switching

### Performance Optimizations

1. **Lazy Loading**: Load wallet connectors only when needed
2. **Connection Pooling**: Reuse connections where possible
3. **Caching**: Cache account info and balances
4. **Background Sync**: Sync wallet state in background

## Troubleshooting

### Common Issues

1. **Wallet not detected**: Ensure wallet is properly installed and accessible
2. **Connection fails**: Check browser permissions and wallet state
3. **Transaction signing fails**: Verify wallet has sufficient permissions
4. **Auto-reconnection fails**: Clear cookies and reconnect manually

### Debug Mode

Enable debug logging by setting:

```typescript
// In development
console.log('Wallet state:', {
  isWalletReady: isWalletReady.value,
  currentWalletType: currentWalletType.value,
  address: address.value,
  error: error.value
})
```

## Security Considerations

1. **Private Key Security**: Never expose private keys in the application
2. **Transaction Validation**: Always validate transactions before signing
3. **Error Handling**: Don't expose sensitive information in error messages
4. **Connection Persistence**: Use secure cookie settings
5. **HTTPS Only**: Ensure all wallet connections use HTTPS in production

## Support

For issues or questions about the enhanced wallet integration:

1. Check the wallet-specific documentation
2. Review the error messages in the browser console
3. Test with different wallet types to isolate issues
4. Verify wallet installation and permissions 