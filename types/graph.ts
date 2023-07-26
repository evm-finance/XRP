export interface EvmTransaction {
  chainId: number
  timestamp: number
  block: number
  from: string
  gasFees: number
  gssLimit: number
  hash: string
  isPending: boolean
  status: string
  nonce: number
  to: string
  txDataHex: string
  value: string
  Input: {
    methodSigDataStr: string
    inputsSigDataStr: string
    inputsMap: Object | null
    argsMap: Object | null
  }
  logEvents: {
    items: {
      network: string
      contract: string
      name: string
      topic: string
      address: string
      signature: string
      allFunctionParams: Object | null
    }[]
  }
}
