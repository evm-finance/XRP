import { gql } from 'graphql-tag'

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

export const BlocksSubscriptionGQL = gql`
  subscription BlocksStreamGQL($network: String!) {
    block(network: $network) {
      network
      blockNumber
      minedAt
      txCount
      swapCount
      pairCreatedCount
      mintCount
      metrics {
        items {
          totalLiquidity
          change1H
          token0Symbol
          token1Symbol
        }
      }
      XRPLedger {
        ledgerHash
        eventsCount
      }
    }
  }
`

export const EvmTransactionGQL = gql`
  query EvmTransactionGQL($chainId: Int!, $hash: String!) {
    evmTransaction(chainId: $chainId, hash: $hash) {
      chainId
      timestamp
      block
      status
      from
      gasPrice
      transactionFee
      gssLimit
      hash
      isPending
      nonce
      to
      txDataHex
      value
      input {
        methodSigDataStr
        inputsSigDataStr
        inputsMap
        argsMap
        fullFunctionSig
        functionName
      }
      logEvents {
        items {
          network
          contract
          name
          topic
          address
          signature
          allFunctionParams
          outputDataMapHex
        }
      }
    }
  }
`
export const XRPTransactionGQL = gql`
  query EvmTransactionGQL($hash: String!) {
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

export const XRPDefiDataGQL = gql`
query XRPDefiDataGQL ($address: String!) {
  account
  xrpBalance
  xrpPrice
  xrpTransactions {
      amount
      destination
      transactionType
      fee
      hash
      ledgerIndex
  }
  xrpTokenBalances {
    tokenSymbol
    tokenIssuer
    tokenName
    balance
    priceXrp
  }
}`

export const XRPAccountTransactionsGQL = gql`
query XRPAccountTransactionsGQL ($address: String!) {
  xrpTransactions(address: $address) {
      amount
      destination
      transactionType
      fee
      hash
      ledgerIndex
  }
}`

export const XRPAccountBalancesGQL = gql`
query XRPAccountBalancesGQL ($account: String!) {
  
  xrpAccountBalances (account: $account) {
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
}`

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
      token1Balance
      token2Balance
      totalSupply
      createdAt
      lastUpdated
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

export const AaveMarketsQGL = gql`
  query AaveMarketsQGL($chainId: Int!, $version: String!) {
    aavePools(chainId: $chainId, version: $version) {
      id
      underlyingAsset
      name
      symbol
      decimals
      totalLiquidity
      price {
        id
        priceInEth
        priceUsd
      }
      liquidityRate
      stableBorrowRate
      variableBorrowRate
      aEmissionPerSecond
      vEmissionPerSecond
      sEmissionPerSecond
      availableLiquidity
      utilizationRate
      totalATokenSupply
      totalCurrentVariableDebt
      totalPrincipalStableDebt
      addresses {
        aTokenAddress
        aTokenSymbol
        stableDebtTokenAddress
        variableDebtTokenAddress
        decimals
        address
        symbol
      }
      portfolioVal {
        totalDeposits
        walletBal
        stableBorrow
        variableBorrow
      }
      totalLiquidityAsCollateral
      baseLTVasCollateral
      reserveLiquidationThreshold
      reserveLiquidationBonus
      usageAsCollateralEnabled

      borrowingEnabled
      stableBorrowRateEnabled
    }
  }
`

export const SupportedChainsGQL = gql`
  query SupportedChainsGQL {
    networks {
      id
      chainIdentifier
      name
      symbol
      nativeTokenSymbol
      rpcUrl
      blockExplorerUrl
      dex {
        name
        value
        symbol
      }
      weth {
        chainId
        address
        symbol
      }
    }
  }
`
