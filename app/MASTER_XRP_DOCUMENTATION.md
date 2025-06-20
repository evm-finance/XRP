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
- Copy to clipboard functionality
- Total value calculation
- Wallet-aware data loading

#### xrpAccountHistory.vue
**Location**: `components/xrp/xrpAccountHistory.vue`
**Purpose**: Displays XRP account transaction history
**Key Features**:
- Transaction grid display
- Column configuration
- Data formatting
- Copy to clipboard functionality
- Transaction hash linking
- Action type color coding
- Mock transaction data generation

### Portfolio Components

#### XrpBalanceWidget.vue
**Location**: `components/portfolio/XrpBalanceWidget.vue`
**Purpose**: XRP balance widget for portfolio pages
**Key Features**:
- Compact balance display
- Integration with portfolio layout
- Wallet-aware data loading
- Consistent styling with EVM components
- Total balance calculation

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
**Components Used**:
- `XrpBalanceWidget`
- Enhanced token table with copy functionality

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

### Portfolio Balances Page
**Location**: `pages/portfolio/balances.vue`
**Purpose**: Multi-chain balance display
**Components Used**:
- `PortfolioBalanceGrid`: EVM balance grids
- `XrpBalanceWidget`: XRP balance widget
- `BalancesChart`: Balance visualization

### XRP Trust Lines Page
**Location**: `pages/xrp-trust-lines.vue`
**Purpose**: Comprehensive Trust line management interface
**Components Used**:
- `useXrpTrustLines`: Reusable composable for Trust line operations and data management

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

### Task 9: Add Trust line interface - COMPLETED
**Date**: Current session
**Changes Made**:
1. **Created XRP Trust Lines page**: New `/xrp-trust-lines.vue` page for comprehensive Trust line management
2. **Created useXrpTrustLines composable**: Reusable composable for Trust line operations and data management
3. **Integrated GEM wallet**: Full wallet integration for Trust line operations
4. **Added Trust line creation**: Form for creating new Trust lines with currency, issuer, limit, and flags
5. **Added Trust line management**: View, edit, and delete existing Trust lines
6. **Added popular tokens section**: Quick Trust line creation for common tokens (USDC, USDT, BTC, ETH)
7. **Added Trust line flags support**: Support for all XRP Trust line flags (tfSetfAuth, tfSetNoRipple, etc.)
8. **Added navigation integration**: Trust lines page accessible from main navigation menu
9. **Added copy functionality**: Copy addresses and transaction hashes to clipboard
10. **Added form validation**: Comprehensive validation for currency codes, addresses, and amounts
11. **Added responsive design**: Mobile-friendly interface with proper data tables
12. **Added mock data**: Development data for testing Trust line functionality

**Files Created/Modified**:
- `app/pages/xrp-trust-lines.vue` (NEW)
- `app/composables/useXrpTrustLines.ts` (NEW)
- `components/common/ui/menu/MainNavigationMenu.vue` (MODIFIED - added navigation link)

**Key Features Added**:
- Complete Trust line management interface
- GEM wallet integration for Trust line operations
- Create, read, update, delete Trust line operations
- Popular tokens quick setup
- Trust line flags configuration
- Address and transaction copying
- Form validation and error handling
- Responsive data tables
- Navigation integration
- Mock data for development
- TypeScript interfaces for type safety

**Technical Implementation**:
- Uses `@gemwallet/api` for Trust line operations
- Integrates with existing XRP client plugin
- Follows existing component patterns and styling
- Implements proper error handling and loading states
- Uses Vue 3 Composition API for reactive state management
- Includes comprehensive TypeScript type definitions

### Task 7: Update and improve XRP token pages - COMPLETED
**Date**: Current session
**Changes Made**:
1. **Created XRP-specific token page**: New `/xrp-token/_id.vue` page for XRP tokens with issuer addresses
2. **Removed EVM interfaces**: No Uniswap/Aave components for XRP tokens
3. **Added XRP screener fields**: Issuer name, issuer address, price, market cap, volume
4. **Added copy functionality**: Copy icon for issuer addresses with clipboard integration
5. **Added AMM chart data**: Placeholder for AMM charts with time tabs (1H, 1D, 1W)
6. **Added wallet transactions**: Transaction table with proper linking to XRP explorer
7. **Added USD/XRP price toggle**: Price display mode switching between XRP and USD
8. **Created useXrpToken composable**: Reusable composable for XRP token data management
9. **Updated XRP screener links**: Links now point to XRP token page instead of EVM token page

**Files Created/Modified**:
- `pages/xrp-token/_id.vue` (NEW)
- `app/composables/useXrpToken.ts` (NEW)
- `pages/xrp-screener.vue` (MODIFIED - updated links)

**Key Features Added**:
- XRP-specific token profile with issuer information
- Copy-to-clipboard functionality for addresses
- Community links (Telegram, Twitter, Discord)
- Token balances display with XRP and USD values
- AMM chart placeholder with timeframe selection
- XRP AMM swap integration
- Wallet transaction history with transaction type color coding
- Transaction hash linking to XRP explorer
- Price display toggle between XRP and USD
- Responsive design matching existing EVM components
- Apollo GraphQL integration for real data
- Mock data generation for development

### Task 6: AMM Interface and Swap Functionality - COMPLETED
**Date**: Previous session
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

### Task 8: Enhanced XRP Balances and Transactions Display - COMPLETED
**Date**: Current session
**Changes Made**:
1. **Enhanced xrpBalances component**: Added proper data table display, copy functionality, total value calculation
2. **Enhanced xrpAccountHistory component**: Added transaction table with copy functionality, action color coding, transaction linking
3. **Created XrpBalanceWidget component**: Portfolio-compatible balance widget
4. **Updated portfolio balances page**: Integrated XRP balance widget alongside EVM grids
5. **Updated XRP screener page**: Added balance widget and enhanced token table with copy functionality
6. **Added wallet-aware data loading**: Components now use connected wallet address when available

**Files Created/Modified**:
- `components/xrp/xrpBalances.vue` (ENHANCED)
- `components/xrp/xrpAccountHistory.vue` (ENHANCED)
- `components/portfolio/XrpBalanceWidget.vue` (NEW)
- `pages/portfolio/balances.vue` (MODIFIED)
- `pages/xrp-screener.vue` (MODIFIED)

**Key Features Added**:
- Copy to clipboard functionality for addresses and transaction hashes
- Transaction linking to XRP explorer
- Action type color coding (Payment=green, OfferCreate=blue, etc.)
- Wallet-aware data loading with fallback to default address
- Consistent styling with existing EVM components
- Total balance calculations and display
- Enhanced error states and loading indicators

## Next Steps

### Task 9: Add block summary records and analytics
- Implement block reader for summary records
- Display top traded tokens and volumes
- Add heatmap for XRPL liquidity pairings
- Create analytics dashboard

## Notes

- All XRP components follow the existing EVM component patterns for consistency
- GEM wallet integration is partially implemented and needs enhancement
- XRPL client connection is established but transaction signing needs GEM wallet integration
- Apollo Federation dependency conflicts have been resolved with .npmrc configuration
- XRP balance and transaction components now provide full functionality with proper data display
- Copy functionality and transaction linking enhance user experience
- Wallet-aware components automatically use connected wallet when available 