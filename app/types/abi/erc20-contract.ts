import {
  ContractTransaction,
  ContractInterface,
  BytesLike as Arrayish,
  BigNumber,
  BigNumberish,
} from 'ethers'
import { EthersContractContextV5 } from 'ethereum-abi-types-generator'

export type ContractContext = EthersContractContextV5<
  Erc20Contract,
  Erc20ContractMethodNames,
  Erc20ContractEventsContext,
  Erc20ContractEvents
>

export declare type EventFilter = {
  address?: string
  topics?: Array<string>
  fromBlock?: string | number
  toBlock?: string | number
}

export interface ContractTransactionOverrides {
  /**
   * The maximum units of gas for the transaction to use
   */
  gasLimit?: number
  /**
   * The price (in wei) per unit of gas
   */
  gasPrice?: BigNumber | string | number | Promise<any>
  /**
   * The nonce to use in the transaction
   */
  nonce?: number
  /**
   * The amount to send with the transaction (i.e. msg.value)
   */
  value?: BigNumber | string | number | Promise<any>
  /**
   * The chain ID (or network ID) to use
   */
  chainId?: number
}

export interface ContractCallOverrides {
  /**
   * The address to execute the call as
   */
  from?: string
  /**
   * The maximum units of gas for the transaction to use
   */
  gasLimit?: number
}
export type Erc20ContractEvents = 'Transfer' | 'Approval'
export interface Erc20ContractEventsContext {
  Transfer(...parameters: any): EventFilter
  Approval(...parameters: any): EventFilter
}
export type Erc20ContractMethodNames =
  | 'name'
  | 'approve'
  | 'totalSupply'
  | 'transferFrom'
  | 'decimals'
  | 'transferAndCall'
  | 'decreaseApproval'
  | 'balanceOf'
  | 'symbol'
  | 'transfer'
  | 'increaseApproval'
  | 'allowance'
  | 'new'
export interface TransferEventEmittedResponse {
  from: string
  to: string
  value: BigNumberish
  data: Arrayish
}
export interface ApprovalEventEmittedResponse {
  owner: string
  spender: string
  value: BigNumberish
}
export interface Erc20Contract {
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  name(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _spender Type: address, Indexed: false
   * @param _value Type: uint256, Indexed: false
   */
  approve(
    _spender: string,
    _value: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  totalSupply(overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _from Type: address, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _value Type: uint256, Indexed: false
   */
  transferFrom(
    _from: string,
    _to: string,
    _value: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  decimals(overrides?: ContractCallOverrides): Promise<number>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _to Type: address, Indexed: false
   * @param _value Type: uint256, Indexed: false
   * @param _data Type: bytes, Indexed: false
   */
  transferAndCall(
    _to: string,
    _value: BigNumberish,
    _data: Arrayish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _spender Type: address, Indexed: false
   * @param _subtractedValue Type: uint256, Indexed: false
   */
  decreaseApproval(
    _spender: string,
    _subtractedValue: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param _owner Type: address, Indexed: false
   */
  balanceOf(
    _owner: string,
    overrides?: ContractCallOverrides
  ): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  symbol(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _to Type: address, Indexed: false
   * @param _value Type: uint256, Indexed: false
   */
  transfer(
    _to: string,
    _value: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _spender Type: address, Indexed: false
   * @param _addedValue Type: uint256, Indexed: false
   */
  increaseApproval(
    _spender: string,
    _addedValue: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param _owner Type: address, Indexed: false
   * @param _spender Type: address, Indexed: false
   */
  allowance(
    _owner: string,
    _spender: string,
    overrides?: ContractCallOverrides
  ): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: constructor
   */
  'new'(overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
}
