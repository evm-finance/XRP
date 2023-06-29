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
