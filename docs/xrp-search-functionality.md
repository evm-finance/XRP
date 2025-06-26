# XRP Search Functionality

## Overview

The search bar has been enhanced to support XRP address lookup for balance and transaction lookup. Users can now search for XRP addresses, transaction hashes, and ledger indices directly from the global search bar.

## Features

### 1. XRP Address Search
- **Format**: XRP addresses start with 'r' and are 25-34 characters long
- **Example**: `rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh`
- **Results**: Shows options to view balances and transaction history

### 2. XRP Transaction Hash Search
- **Format**: 64-character hexadecimal string
- **Example**: `A1B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456`
- **Results**: Direct link to transaction details

### 3. XRP Ledger Index Search
- **Format**: Numeric ledger index
- **Example**: `12345678`
- **Results**: Direct link to ledger details

## Search Results

When searching for XRP-related data, the search results will show:

- **XRP Account**: Links to balance and transaction pages
- **XRP Transaction**: Links to transaction details
- **XRP Ledger**: Links to ledger details
- **Color-coded chips**: Different colors for different types of results

## Pages Created

### 1. XRP Balances Page (`/xrp-balances`)
- **URL**: `/xrp-balances?address=<xrp_address>`
- **Features**:
  - Display XRP balance and USD value
  - Show all token balances with prices
  - Copy address functionality
  - Search and filter tokens
  - Link to token details

### 2. XRP Transactions Page (`/xrp-transactions`)
- **URL**: `/xrp-transactions?address=<xrp_address>`
- **Features**:
  - Display transaction history
  - Filter by transaction type
  - Search transactions
  - Copy transaction hashes
  - Link to transaction and ledger details

## Implementation Details

### Search Validation
- **XRP Address**: `/^r[a-km-zA-HJ-NP-Z1-9]{25,34}$/`
- **Transaction Hash**: `/^[A-Fa-f0-9]{64}$/`
- **Ledger Index**: `/^\d+$/`

### GraphQL Integration
The search uses existing GraphQL queries:
- `XRPAccountBalancesGQL` for balance lookup
- `XRPAccountTransactionsGQL` for transaction history
- `XRPTransactionGQL` for transaction details

### Mock Data
For testing purposes, mock data is available in `composables/useXrpMockData.ts`:
- Sample XRP addresses
- Sample transaction hashes
- Sample ledger indices
- Mock account and transaction data

## Usage Examples

### Search for XRP Address
1. Type an XRP address in the search bar
2. Select "XRP Account - View Balances & Transactions"
3. View account balances and token holdings
4. Or select "XRP Account - View Transaction History"
5. View complete transaction history

### Search for Transaction
1. Type a transaction hash in the search bar
2. Select "XRP Transaction - View Details"
3. View detailed transaction information

### Search for Ledger
1. Type a ledger index in the search bar
2. Select "XRP Ledger - View Details"
3. View ledger information

## Navigation

The new pages are accessible via:
- Direct URL with address parameter
- Global search bar
- Navigation menu (XRP-Balances, XRP-Transactions)

## Testing

To test the functionality:
1. Use mock addresses: `rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh`
2. Use mock transaction hashes: `A1B2C3D4E5F6789012345678901234567890ABCDEF1234567890ABCDEF123456`
3. Use mock ledger indices: `12345678`

## Future Enhancements

- Real-time balance updates
- Transaction notifications
- Advanced filtering options
- Export functionality
- Mobile optimization
- Integration with XRP wallets 