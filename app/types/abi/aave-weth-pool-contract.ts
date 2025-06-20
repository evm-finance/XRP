import { ContractTransaction, ContractInterface, BytesLike as Arrayish, BigNumber, BigNumberish } from 'ethers'
import { EthersContractContextV5 } from 'ethereum-abi-types-generator'

export type ContractContext = EthersContractContextV5<
  AaveWethPoolContract,
  AaveWethPoolContractMethodNames,
  AaveWethPoolContractEventsContext,
  AaveWethPoolContractEvents
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
export type AaveWethPoolContractEvents = 'OwnershipTransferred'
export interface AaveWethPoolContractEventsContext {
  OwnershipTransferred(...parameters: any): EventFilter
}
export type AaveWethPoolContractMethodNames =
  | 'new'
  | 'authorizeLendingPool'
  | 'borrowETH'
  | 'depositETH'
  | 'emergencyEtherTransfer'
  | 'emergencyTokenTransfer'
  | 'getWETHAddress'
  | 'owner'
  | 'renounceOwnership'
  | 'repayETH'
  | 'transferOwnership'
  | 'withdrawETH'
export interface OwnershipTransferredEventEmittedResponse {
  previousOwner: string
  newOwner: string
}
export interface AaveWethPoolContract {
  populateTransaction: any
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: constructor
   * @param weth Type: address, Indexed: false
   */
  'new'(weth: string, overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param lendingPool Type: address, Indexed: false
   */
  authorizeLendingPool(lendingPool: string, overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param lendingPool Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param interesRateMode Type: uint256, Indexed: false
   * @param referralCode Type: uint16, Indexed: false
   */
  borrowETH(
    lendingPool: string,
    amount: BigNumberish,
    interesRateMode: BigNumberish,
    referralCode: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param lendingPool Type: address, Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   * @param referralCode Type: uint16, Indexed: false
   */
  depositETH(
    lendingPool: string,
    onBehalfOf: string,
    referralCode: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param to Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   */
  emergencyEtherTransfer(
    to: string,
    amount: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param token Type: address, Indexed: false
   * @param to Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   */
  emergencyTokenTransfer(
    token: string,
    to: string,
    amount: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  getWETHAddress(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  owner(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   */
  renounceOwnership(overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param lendingPool Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param rateMode Type: uint256, Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   */
  repayETH(
    lendingPool: string,
    amount: BigNumberish,
    rateMode: BigNumberish,
    onBehalfOf: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param newOwner Type: address, Indexed: false
   */
  transferOwnership(newOwner: string, overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param lendingPool Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param to Type: address, Indexed: false
   */
  withdrawETH(
    lendingPool: string,
    amount: BigNumberish,
    to: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
}
