import { ref } from '@nuxtjs/composition-api'
import { createOffer, getAddress, isInstalled } from '@gemwallet/api'

export interface BaseTransactionRequest {
  // Integer amount of XRP, in drops, to be destroyed as a cost for distributing this transaction to the network.
  // Some transaction types have different minimum requirements.
  fee?: string
  // The sequence number of the account sending the transaction. A transaction is only valid if the Sequence number is
  // exactly 1 greater than the previous transaction from the same account. The special case 0 means the transaction is
  // using a Ticket instead.
  sequence?: number
  // Hash value identifying another transaction. If provided, this transaction is only valid if the sending account's
  // previously-sent transaction matches the provided hash.
  accountTxnID?: string
  // Highest ledger index this transaction can appear in. Specifying this field places a strict upper limit on how long
  // the transaction can wait to be validated or rejected.
  lastLedgerSequence?: number
  // Additional arbitrary information used to identify this transaction.
  // Each attribute of each memo must be hex encoded.
  memos?: Memo[]
  // Array of objects that represent a multi-signature which authorizes this transaction.
  signers?: Signer[]
  // Arbitrary integer used to identify the reason for this payment, or a sender on whose behalf this transaction is
  // made. Conventionally, a refund should specify the initial payment's SourceTag as the refund payment's
  // DestinationTag.
  sourceTag?: number
  // Hex representation of the public key that corresponds to the private key used to sign this transaction. If an empty
  // string, indicates a multi-signature is present in the Signers field instead.
  signingPubKey?: string
  // The sequence number of the ticket to use in place of a Sequence number. If this is provided, Sequence must be 0.
  // Cannot be used with AccountTxnID.
  ticketSequence?: number
  // The signature that verifies this transaction as originating from the account it says it is from.
  txnSignature?: string
}

export interface SetTrustlineRequest extends BaseTransactionRequest {
  // The maximum amount of currency that can be exchanged to the trustline
  limitAmount: IssuedCurrencyAmount
  // Flags to set on the transaction
  flags?: TrustSetFlags
}

export type TrustSetFlags =
  | {
      tfSetfAuth?: boolean
      tfSetNoRipple?: boolean
      tfClearNoRipple?: boolean
      tfSetFreeze?: boolean
      tfClearFreeze?: boolean
    }
  | number

export interface IssuedCurrencyAmount {
  currency: string
  issuer: string
  value: string
}

export interface Memo {
  memo: {
    memoType?: string
    memoData?: string
    memoFormat?: string
  }
}

export interface Signer {
  signer: {
    account: string
    txnSignature: string
    signingPubKey: string
  }
}

export type offerTypes = 'buy' | 'sell'
export const xrpOffers = ref<Array<offerTypes>>(['buy', 'sell'])

export interface CreateOfferRequest extends BaseTransactionRequest {
  flags?: OfferCreateFlagsInterface
  // Time after which the Offer is no longer active, in seconds since the Ripple Epoch.
  expiration?: number
  // An Offer to delete first, specified in the same way as OfferCancel.
  offerSequence?: number
  // The amount and type of currency being sold.
  takerGets: Amount
  // The amount and type of currency being bought.
  takerPays: Amount
}

export type OfferCreateFlagsInterface =
  | {
      tfPassive?: boolean
      tfImmediateOrCancel?: boolean
      tfFillOrKill?: boolean
      tfSell?: boolean
    }
  | number

export type Amount =
  | {
      currency: string
      issuer: string
      value: string
    }
  | string

const payload = {
  takerPays: {
    currency: 'ETH',
    issuer: 'rcA8X3TVMST1n3CJeAdGk1RdRCHii7N2h',
    value: '0.0003',
  },
  takerGets: '1000000', // 1 XRP
  flags: {},
  fee: '199',
}

const isOpen = ref(false)

export default function useXrpTrade() {
  async function connectWallet() {
    isInstalled().then((response: any) => {
      if (response.result.isInstalled) {
        getAddress().then((response: any) => {
          console.log(`Your address: ${response.result?.address}`)
        })
      }
    })
  }

  function buy() {
    console.log('buy')
    console.log(payload)
    isInstalled().then((response: any) => {
      if (response.result.isInstalled) {
        createOffer(payload).then((response: any) => {
          console.log('Transaction Hash: ', response.result?.hash)
        })
      }
    })
  }

  function openDialog() {
    isOpen.value = true
    console.log(isOpen)
  }

  function closeDialog() {
    isOpen.value = false
    console.log(isOpen)
  }

  function sell() {
    console.log('sell')
  }

  return {
    buy,
    sell,
    connectWallet,
    isOpen,
    openDialog,
    closeDialog,
  }
}
