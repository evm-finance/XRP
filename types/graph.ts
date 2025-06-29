export interface EvmTransaction {
  chainId: number
  timestamp: number
  block: number
  from: string
  gasPrice: number
  transactionFee: number
  gssLimit: number
  hash: string
  isPending: boolean
  status: string
  nonce: number
  to: string
  txDataHex: string
  value: string
  input:
    | {
        methodSigDataStr: string
        inputsSigDataStr: string
        fullFunctionSig: string
        functionName: string
        inputsMap: Record<string, string> | undefined
        argsMap: Record<string, Record<string, string>> | undefined
      }
    | undefined
  logEvents:
    | {
        items: {
          network: string
          contract: string
          name: string
          topic: string
          address: string
          signature: string
          outputDataMapHex: string
          allFunctionParams: Record<string, { indexed: boolean; name: string; type: string; value: any }> | undefined
        }[]
      }
    | undefined
}

export type GraphData = {
  // XRP-specific graph data structure
  // No chainId needed for XRP
}
