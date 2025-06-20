# MASTER XRP DOCUMENTATION

This document provides comprehensive documentation for all XRP-related files in the project.

## Table of Contents
1. [Composables](#composables)
2. [Components](#components)
3. [Pages](#pages)
4. [Plugins](#plugins)
5. [Types](#types)
6. [Configuration](#configuration)

## Composables

### useXrpAmm.ts
**Location**: `composables/useXrpAmm.ts`
**Purpose**: Provides XRP AMM (Automated Market Maker) trading functionality
**Key Features**:
- XRPL client initialization and connection
- Account balance retrieval
- AMM swap quote calculation
- Transaction preparation and submission
- Error handling and state management

**Key Functions**:
- `initializeClient()`: Connects to XRPL network
- `getAccountInfo()`: Retrieves account balances
- `calculateSwapQuote()`: Calculates expected output for swaps
- `performSwap()`: Executes AMM trades
- `resetTransaction()`: Resets transaction state

**Interfaces**:
- `XRPToken`: Represents XRP tokens with currency, issuer, symbol, name, decimals
- `XRPSwapParams`: Parameters for swap operations
- `XRPSwapResult`: Result of swap operations

### useXrpTrade.ts
**Location**: `composables/useXrpTrade.ts`
**Purpose**: Provides XRP trading functionality using GEM wallet
**Key Features**:
- GEM wallet integration
- Offer creation for trading
- Buy/sell operations
- Dialog management

**Key Functions**:
- `connectWallet()`: Connects to GEM wallet
- `buy()`: Creates buy offers
- `sell()`: Creates sell offers
- `openDialog()` / `closeDialog()`: Manages UI dialogs

### useXrpAccounts.ts
**Location**: `composables/useXrpAccounts.ts`
**Purpose**: Manages XRP account data and transactions
**Key Features**:
- Account balance queries
- Transaction history
- Apollo GraphQL integration

### useXrpScrerener.ts
**Location**: `composables/useXrpScrerener.ts`
**Purpose**: Provides XRP token screener functionality
**Key Features**:
- Token listing and filtering
- Price and volume data
- Market statistics

## Components

### Trading Components

#### XrpAmmSwap.vue
**Location**: `components/trading/XrpAmmSwap.vue`
**Purpose**: Main XRP AMM swap interface component
**Key Features**:
- Token input/output fields
- Swap quote display
- Transaction execution
- Error handling
- Loading states

**Props**:
- `height`: Component height
- `width`: Component width
- `inToken`: Input token configuration
- `outToken`: Output token configuration

#### XrpTokenInputField.vue
**Location**: `components/trading/XrpTokenInputField.vue`
**Purpose**: Token input field for XRP trading
**Key Features**:
- Amount input with validation
- Token selection button
- Balance display
- Price display

**Props**:
- `formTradeDirection`: Input or output direction
- `token`: Token configuration
- `balance`: User's token balance
- `fiatPrice`: Token's fiat price
- `expectedConvertQuote`: Expected output amount
- `loading`: Loading state

#### XrpTokenMenuDialog.vue
**Location**: `components/trading/XrpTokenMenuDialog.vue`
**Purpose**: Token selection dialog for XRP trading
**Key Features**:
- Token search functionality
- Common XRP tokens list
- Token selection with icons

#### XrpOfferInput.vue
**Location**: `components/trading/XrpOfferInput.vue`
**Purpose**: XRP offer input component (legacy)
**Status**: Basic implementation, needs enhancement

### XRP-Specific Components

#### xrpBalances.vue
**Location**: `components/xrp/xrpBalances.vue`
**Purpose**: Displays XRP account balances
**Key Features**:
- Apollo GraphQL integration
- Balance formatting
- Price display
- Issuer address display

#### xrpAccountHistory.vue
**Location**: `components/xrp/xrpAccountHistory.vue`
**Purpose**: Displays XRP account transaction history
**Key Features**:
- Transaction grid display
- Column configuration
- Data formatting

### Common Components

#### GemWalletConnector.vue
**Location**: `components/common/GemWalletConnector.vue`
**Purpose**: GEM wallet connection interface
**Key Features**:
- Wallet connection status
- Address display
- Wallet actions menu
- Explorer navigation

#### GlobalSearch.vue
**Location**: `components/common/GlobalSearch.vue`
**Purpose**: Global search functionality including XRP
**Key Features**:
- XRP ledger search
- Transaction hash validation
- Address type detection
- Search result categorization

## Pages

### XRP Portfolio Page
**Location**: `pages/xrp-portfolio.vue`
**Purpose**: XRP portfolio overview
**Components Used**:
- `xrpAccountHistory`
- `xrpBalances`

### XRP Screener Page
**Location**: `pages/xrp-screener.vue`
**Purpose**: XRP token screener interface

### XRP Explorer Pages
**Location**: `pages/xrp-explorer/`
**Purpose**: XRP ledger exploration
- `index.vue`: Main explorer interface
- `ledger/_id.vue`: Individual ledger view
- `tx/_id.vue`: Transaction details

### Swap Page
**Location**: `pages/swap.vue`
**Purpose**: Trading interface
**Components Used**:
- `EvmSwap`: EVM trading interface
- `XrpAmmSwap`: XRP AMM trading interface

## Plugins

### XRP Client Plugin
**Location**: `plugins/web3/xrp.client.ts`
**Purpose**: XRP wallet and client management
**Key Features**:
- GEM wallet integration
- Address management
- Connection state management
- Event handling

**Key Functions**:
- `connectWallet()`: Connects to GEM wallet
- `disconnectWallet()`: Disconnects wallet
- `getAddress()`: Retrieves wallet address
- `isInstalled()`: Checks wallet installation

## Types

### XRP Types
**Location**: `types/apollo/main/types.ts`
**Purpose**: TypeScript type definitions for XRP data

**Key Types**:
- `XRPDefiData`: XRP DeFi data structure
- `XRPAccountBalances`: Account balance data
- `XrpTransaction`: Transaction data structure

## Configuration

### Package Configuration
**Location**: `package.json`
**XRP Dependencies**:
- `xrpl`: XRPL JavaScript library
- `@gemwallet/api`: GEM wallet API

### NPM Configuration
**Location**: `.npmrc`
**Purpose**: Resolves dependency conflicts
**Settings**:
- `engine-strict=false`: Suppresses engine warnings
- `legacy-peer-deps=true`: Uses legacy peer dependency resolution

## Recent Updates

### Task 6: AMM Interface and Swap Functionality - COMPLETED
**Date**: Current session
**Changes Made**:
1. **Installed xrpl.js**: Added XRPL JavaScript library for blockchain interactions
2. **Created useXrpAmm composable**: Full AMM trading functionality with quote calculation
3. **Created XrpAmmSwap component**: Complete swap interface with token selection
4. **Created XrpTokenInputField component**: Token input with balance display
5. **Created XrpTokenMenuDialog component**: Token selection dialog
6. **Updated swap page**: Integrated both EVM and XRP AMM interfaces
7. **Added .npmrc**: Resolved Apollo Federation dependency conflicts

**Files Created/Modified**:
- `composables/useXrpAmm.ts` (NEW)
- `components/trading/XrpAmmSwap.vue` (NEW)
- `components/trading/XrpTokenInputField.vue` (NEW)
- `components/trading/XrpTokenMenuDialog.vue` (NEW)
- `pages/swap.vue` (MODIFIED)
- `.npmrc` (NEW)
- `package.json` (MODIFIED - added xrpl dependency)

## Next Steps

### Task 7: Update and improve XRP token pages
- Remove Uniswap/Aave interfaces from XRP pages
- Add XRP screener fields (issuer, price, volume)
- Add copy icon for issuer addresses
- Add AMM chart data with time tabs
- Add USD/XRP price toggle

### Task 8: Add block summary records and analytics
- Implement block reader for summary records
- Display top traded tokens and volumes
- Add heatmap for XRPL liquidity pairings
- Create analytics dashboard

## Notes

- All XRP components follow the existing EVM component patterns for consistency
- GEM wallet integration is partially implemented and needs enhancement
- XRPL client connection is established but transaction signing needs GEM wallet integration
- Apollo Federation dependency conflicts have been resolved with .npmrc configuration 