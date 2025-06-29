# GraphQL Query Inventory Documentation

## Overview
This document provides a comprehensive inventory of all GraphQL queries in the XRP project, organized by file and functionality.

## File Structure

### 1. `apollo/queries.ts` - Main Query Definitions
**Status**: Primary query file with XRP-specific queries
**Lines**: 392 total

#### XRP Block Queries
- `BlocksXrpGQL` - Query for XRP ledger blocks with transaction counts
- Parameters: `$network: String!`
- Returns: Block data with XRPLedger information

#### XRP Transaction Queries
- `XRPTransactionGQL` - Query for individual XRP transaction details
- Parameters: `$hash: String!`
- Returns: Complete transaction data including metadata, memos, validation status

#### XRP Account Queries
- `XRPAccountTransactionsGQL` - Query for account transaction history
- Parameters: `$address: String!`
- Returns: Transaction list for specific address

- `XRPAccountBalancesGQL` - Query for account balances and token holdings
- Parameters: `$account: String!`
- Returns: XRP balance, price, and token balances

#### XRP Screener Queries
- `XRPScreenerGQL` - Query for token screener data
- Parameters: None
- Returns: Token list with market data (currency, issuer, price, volume)

#### XRP AMM Queries
- `XRPAmmPoolsGQL` - Query for AMM pool list
- Parameters: None
- Returns: Pool list with token pairs, liquidity, volume, fees

- `XRPAmmPoolDetailsGQL` - Query for detailed pool information
- Parameters: `$poolId: String!`
- Returns: Detailed pool data including transactions and user positions

- `XRPAmmUserPositionsGQL` - Query for user's AMM positions
- Parameters: `$address: String!`
- Returns: User's pool positions and shares

- `XRPAmmQuoteGQL` - Query for swap quotes
- Parameters: `$poolId: String!, $amount: String!, $fromToken: String!`
- Returns: Swap quote with price impact and fees

- `XRPAmmTransactionsGQL` - Query for pool transaction history
- Parameters: `$poolId: String!, $limit: Int`
- Returns: Transaction list for specific pool

#### XRP Token Queries
- `XRPTokenPriceGQL` - Query for token price data
- Parameters: `$currency: String!, $issuer: String`
- Returns: Price, market cap, volume data

- `XRPTokenBalancesGQL` - Query for token balances
- Parameters: `$address: String!`
- Returns: Token balances for specific address

#### XRP Heatmap Queries
- `XRPHeatmapGQL` - Query for token heatmap data
- Parameters: `$timeFrame: String, $blockSize: String, $limit: Int`
- Returns: Comprehensive token data for heatmap visualization

- `XRPAMMHeatmapGQL` - Query for AMM heatmap data
- Parameters: `$timeFrame: String, $blockSize: String, $limit: Int`
- Returns: AMM pool data for heatmap visualization

- `XRPHeatmapUpdatesGQL` - Query for heatmap updates
- Parameters: `$currencies: [String!], $timeFrame: String`
- Returns: Updated price and market data

#### XRP AMM Liquidity Queries
- `XRPAMMLiquidityValuesGQL` - Query for AMM liquidity values
- Parameters: None
- Returns: Liquidity data with percentages and price impact

- `SimpleAMMLiquidityValuesGQL` - Simplified liquidity query
- Parameters: None
- Returns: Basic liquidity values

- `GetAllAMMLiquidityValuesGQL` - Comprehensive liquidity query
- Parameters: None
- Returns: Complete liquidity data with asset details

#### Test Query
- `TestGraphQLGQL` - Schema introspection query
- Parameters: None
- Returns: GraphQL schema types

### 2. `apollo/main/xrp.query.graphql` - XRP-Specific Queries
**Status**: Dedicated XRP query file
**Lines**: 399 total

#### Account Queries
- `XRPAccountBalances` - Account balance query (duplicate of queries.ts)
- `XRPAccountTransactions` - Account transaction query (duplicate of queries.ts)
- `XRPAccountData` - Combined account data query

#### AMM Queries
- `XRPAMMPools` - AMM pools query
- `XRPAMMPool` - Specific pool query with asset parameters
- `XRPAMMSwapQuote` - Swap quote with asset inputs
- `XRPAMMLiquidityValue` - Pool liquidity value query
- `GetAllAMMLiquidityValues` - All liquidity values (duplicate)
- `XRPAMMTransactions` - AMM transactions with filtering

### 3. `apollo/main/portfolio.query.graphql` - Portfolio Queries
**Status**: Portfolio-specific queries
**Lines**: 55 total

#### Portfolio Queries
- `XRPAccountBalances` - Account balances (duplicate)
- `XRPAccountTransactions` - Account transactions (duplicate)
- `XRPDefiData` - Combined DeFi data for address

### 4. `apollo/main/token.query.graphql` - Token Queries
**Status**: Token and block queries
**Lines**: 235 total

#### Token Queries
- `XRPTokenPrice` - Token price query (duplicate)
- `XRPTokenBalances` - Token balances query (duplicate)
- `XRPScreener` - Screener query (duplicate)

#### Block Queries
- `BlocksXrpGQL` - XRP blocks query (duplicate)
- `BlockGQL` - Specific block query with detailed ledger data

#### Legacy Queries (Non-XRP)
- `TokenQueryGQL` - Legacy token query for EVM chains
- `DailyChartGQL` - Daily chart data
- `ScreenerGQL` - Legacy screener for EVM pools
- `TimeGQL` - Time subscription
- `PriceStreamGQL` - Price stream subscription

### 5. `apollo/main/config.query.graphql` - Configuration Queries
**Status**: Configuration and test queries
**Lines**: 11 total

#### Config Queries
- `TestGraphQL` - Schema test query (duplicate)

### 6. `apollo/main/heatmap.query.graphql` - Heatmap Queries
**Status**: Heatmap-specific queries
**Lines**: Not yet examined

## Query Categories

### XRP-Specific Queries (Critical)
1. **Account Management**
   - Account balances
   - Account transactions
   - Account data

2. **AMM Operations**
   - Pool listings
   - Pool details
   - Swap quotes
   - User positions
   - Transaction history

3. **Token Operations**
   - Token prices
   - Token balances
   - Screener data

4. **Block/Transaction Data**
   - Block information
   - Transaction details
   - Ledger data

5. **Heatmap Data**
   - Token heatmaps
   - AMM heatmaps
   - Real-time updates

### Legacy/Non-XRP Queries (Lower Priority)
1. **EVM Chain Queries**
   - Token queries for EVM chains
   - Pool screener for EVM DEXes
   - Price streams

2. **Configuration Queries**
   - Schema introspection
   - Time subscriptions

## Duplicate Queries Identified
Several queries appear in multiple files with identical or similar functionality:
- `XRPAccountBalances` - appears in 3 files
- `XRPAccountTransactions` - appears in 3 files
- `XRPTokenPrice` - appears in 2 files
- `XRPTokenBalances` - appears in 2 files
- `XRPScreener` - appears in 2 files
- `TestGraphQL` - appears in 3 files

## Recommendations
1. **Consolidate duplicate queries** into single locations
2. **Prioritize XRP-specific queries** for error handling improvements
3. **Remove or deprecate legacy EVM queries** that are no longer needed
4. **Standardize query naming conventions** across files
5. **Add proper error handling** to all XRP queries first

## Next Steps
- Complete Step 1.2: Inventory all GraphQL query usage in components and pages
- Create query registry with usage mapping
- Document current error handling patterns 