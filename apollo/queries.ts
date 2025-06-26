import { gql } from 'graphql-tag'

// XRP Block Queries
export const BlocksXrpGQL = gql`
  query BlocksXrpGQL($network: String!) {
    blocks(network: $network) {
      network
      blockNumber
      minedAt
      txCount
      XRPLedger {
        ledgerHash
        eventsCount
      }
    }
  }
`

// XRP Transaction Queries
export const XRPTransactionGQL = gql`
  query XRPTransactionGQL($hash: String!) {
    xrpTransaction(hash: $hash) {
      account
      amount
      destination
      fee
      flags
      lastLedgerSequence
      offerSequence
      sequence
      signingPubKey
      takerGets
      takerPays
      transactionType
      txnSignature
      date
      hash
      inLedger
      ledgerIndex
      meta
      metadata
      validated
      warnings
      memos
    }
  }
`

// XRP Account Queries
export const XRPAccountTransactionsGQL = gql`
  query XRPAccountTransactionsGQL($address: String!) {
    xrpTransactions(address: $address) {
      amount
      destination
      transactionType
      fee
      hash
      ledgerIndex
    }
  }
`

export const XRPAccountBalancesGQL = gql`
  query XRPAccountBalancesGQL($account: String!) {
    xrpAccountBalances(account: $account) {
      account
      xrpBalance
      xrpPrice
      xrpTokens {
        symbol
        issuer
        name
        balance
        price
        value
      }
    }
  }
`

// XRP Screener Queries
export const XRPScreenerGQL = gql`
  query XRPScreenerGQL {
    xrpScreener {
      currency
      issuerAddress
      icon
      tokenName
      issuerName
      marketcap
      price
      volume24H
    }
  }
`

// XRP AMM Queries
export const XRPAmmPoolsGQL = gql`
  query XRPAmmPoolsGQL {
    xrpAmmPools {
      id
      token1 {
        symbol
        name
        icon
        issuer
      }
      token2 {
        symbol
        name
        icon
        issuer
      }
      liquidity
      volume24h
      fee
      apr
      priceChange24h
    }
  }
`

export const XRPAmmPoolDetailsGQL = gql`
  query XRPAmmPoolDetailsGQL($poolId: String!) {
    xrpAmmPoolDetails(poolId: $poolId) {
      id
      token1 {
        symbol
        name
        icon
        issuer
        price
        priceChange24h
      }
      token2 {
        symbol
        name
        icon
        issuer
        price
        priceChange24h
      }
      liquidity
      volume24h
      fee
      apr
      priceChange24h
      token1Balance
      token2Balance
      totalSupply
      createdAt
      lastUpdated
      transactions {
        hash
        type
        amount
        token
        timestamp
        user
      }
      userPositions {
        user
        poolTokens
        token1Balance
        token2Balance
        share
        value
      }
    }
  }
`

export const XRPAmmUserPositionsGQL = gql`
  query XRPAmmUserPositionsGQL($address: String!) {
    xrpAmmUserPositions(address: $address) {
      poolId
      pool {
        id
        token1 {
          symbol
          name
          icon
        }
        token2 {
          symbol
          name
          icon
        }
        liquidity
        fee
        apr
      }
      poolTokens
      token1Balance
      token2Balance
      share
      value
      lastUpdated
    }
  }
`

export const XRPAmmQuoteGQL = gql`
  query XRPAmmQuoteGQL($poolId: String!, $amount: String!, $fromToken: String!) {
    xrpAmmQuote(poolId: $poolId, amount: $amount, fromToken: $fromToken) {
      inputAmount
      outputAmount
      priceImpact
      fee
      minimumReceived
      price
    }
  }
`

export const XRPAmmTransactionsGQL = gql`
  query XRPAmmTransactionsGQL($poolId: String!, $limit: Int) {
    xrpAmmTransactions(poolId: $poolId, limit: $limit) {
      hash
      type
      amount
      token
      timestamp
      user
      poolId
      fee
      priceImpact
    }
  }
`

// XRP Token Queries
export const XRPTokenPriceGQL = gql`
  query XRPTokenPriceGQL($currency: String!, $issuer: String) {
    xrpTokenPrice(currency: $currency, issuer: $issuer) {
      currency
      issuer
      price
      priceChange24h
      volume24h
      marketcap
      lastUpdated
    }
  }
`

export const XRPTokenBalancesGQL = gql`
  query XRPTokenBalancesGQL($address: String!) {
    xrpTokenBalances(address: $address) {
      currency
      issuer
      balance
      price
      value
      lastUpdated
    }
  }
`

// XRP Heatmap Queries
export const XRPHeatmapGQL = gql`
  query XRPHeatmapGQL($timeFrame: String, $blockSize: String, $limit: Int) {
    xrpHeatmap(timeFrame: $timeFrame, blockSize: $blockSize, limit: $limit) {
      currency
      name
      price_usd
      price_xrp
      price_change_1h
      price_change_4h
      price_change_24h
      price_change_7d
      price_change_30d
      marketcap
      volume_24h
      liquidity
      tvl
      issuer
      hasAmmPool
      hasTrustLine
    }
  }
`

export const XRPAMMHeatmapGQL = gql`
  query XRPAMMHeatmapGQL($timeFrame: String, $blockSize: String, $limit: Int) {
    xrpAmmHeatmap(timeFrame: $timeFrame, blockSize: $blockSize, limit: $limit) {
      poolId
      token1 {
        currency
        symbol
        name
        issuer
      }
      token2 {
        currency
        symbol
        name
        issuer
      }
      liquidity
      volume24h
      volume7d
      fee
      apr
      priceChange24h
      priceChange7d
      tvl
      transactions24h
      uniqueTraders24h
    }
  }
`

export const XRPHeatmapUpdatesGQL = gql`
  query XRPHeatmapUpdatesGQL($currencies: [String!], $timeFrame: String) {
    xrpHeatmapUpdates(currencies: $currencies, timeFrame: $timeFrame) {
      currency
      price_usd
      price_xrp
      price_change_1h
      price_change_4h
      price_change_24h
      price_change_7d
      price_change_30d
      marketcap
      volume_24h
      liquidity
      tvl
    }
  }
`

// XRP AMM Liquidity Queries
export const XRPAMMLiquidityValuesGQL = gql`
  query XRPAMMLiquidityValues {
    xrpAMMLiquidityValues {
      poolId
      asset1ValueUsd
      asset2ValueUsd
      totalLiquidityUsd
      asset1Percentage
      asset2Percentage
      priceImpact
    }
  }
`

export const SimpleAMMLiquidityValuesGQL = gql`
  query SimpleAMMLiquidityValues {
    xrpAMMLiquidityValues {
      poolId
      asset1ValueUsd
      asset2ValueUsd
      totalLiquidityUsd
    }
  }
`

export const GetAllAMMLiquidityValuesGQL = gql`
  query GetAllAMMLiquidityValues {
    xrpAMMLiquidityValues {
      poolId
      asset1 {
        currency
        issuer
      }
      asset2 {
        currency
        issuer
      }
      asset1Balance
      asset2Balance
      asset1ValueUsd
      asset2ValueUsd
      totalLiquidityUsd
      fee
      createdAt
    }
  }
`

// Test Query
export const TestGraphQLGQL = gql`
  query TestGraphQL {
    __schema {
      types {
        name
      }
    }
  }
`
