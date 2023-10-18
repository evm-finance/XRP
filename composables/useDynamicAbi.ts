// This composable library will perfom the following functions
// 1. Hold the rules for ParamInputForm components
// 2. Compute function selector based on ABI
// 3. Prepare function call using ethersjs, taking inputs from Param Forms
import type { ContractTransaction } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'


export type ContractABI = {
    __typename?: 'ContractABI',

  }
  