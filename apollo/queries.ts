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
