import {
  ContractTransaction,
  ContractInterface,
  BytesLike as Arrayish,
  BigNumber,
  BigNumberish,
} from 'ethers'
import { EthersContractContextV5 } from 'ethereum-abi-types-generator'

export type ContractContext = EthersContractContextV5<
  UniswapRouterV2Contract,
  UniswapRouterV2ContractMethodNames,
  UniswapRouterV2ContractEventsContext,
  UniswapRouterV2ContractEvents
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
export type UniswapRouterV2ContractEvents = undefined
export interface UniswapRouterV2ContractEventsContext {}
export type UniswapRouterV2ContractMethodNames =
  | 'new'
  | 'FACTORY'
  | 'PAIR'
  | 'WETH'
  | 'addLiquidity'
  | 'addLiquidityETH'
  | 'getAmountIn'
  | 'getAmountOut'
  | 'getAmountsIn'
  | 'getAmountsOut'
  | 'pairFor'
  | 'quote'
  | 'removeLiquidity'
  | 'removeLiquidityETH'
  | 'removeLiquidityETHSupportingFeeOnTransferTokens'
  | 'removeLiquidityETHWithPermit'
  | 'removeLiquidityETHWithPermitSupportingFeeOnTransferTokens'
  | 'removeLiquidityWithPermit'
  | 'swapETHForExactTokens'
  | 'swapExactETHForTokens'
  | 'swapExactETHForTokensSupportingFeeOnTransferTokens'
  | 'swapExactTokensForETH'
  | 'swapExactTokensForETHSupportingFeeOnTransferTokens'
  | 'swapExactTokensForTokens'
  | 'swapExactTokensForTokensSupportingFeeOnTransferTokens'
  | 'swapTokensForExactETH'
  | 'swapTokensForExactTokens'
export interface UniswapRouterV2Contract {
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: constructor
   * @param _factory Type: address, Indexed: false
   * @param _WETH Type: address, Indexed: false
   */
  'new'(
    _factory: string,
    _WETH: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  FACTORY(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  PAIR(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  WETH(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _tokenA Type: address, Indexed: false
   * @param _tokenB Type: address, Indexed: false
   * @param _amountADesired Type: uint256, Indexed: false
   * @param _amountBDesired Type: uint256, Indexed: false
   * @param _amountAMin Type: uint256, Indexed: false
   * @param _amountBMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  addLiquidity(
    _tokenA: string,
    _tokenB: string,
    _amountADesired: BigNumberish,
    _amountBDesired: BigNumberish,
    _amountAMin: BigNumberish,
    _amountBMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param _token Type: address, Indexed: false
   * @param _amountTokenDesired Type: uint256, Indexed: false
   * @param _amountTokenMin Type: uint256, Indexed: false
   * @param _amountETHMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  addLiquidityETH(
    _token: string,
    _amountTokenDesired: BigNumberish,
    _amountTokenMin: BigNumberish,
    _amountETHMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: pure
   * Type: function
   * @param _amountOut Type: uint256, Indexed: false
   * @param _reserveIn Type: uint256, Indexed: false
   * @param _reserveOut Type: uint256, Indexed: false
   */
  getAmountIn(
    _amountOut: BigNumberish,
    _reserveIn: BigNumberish,
    _reserveOut: BigNumberish,
    overrides?: ContractCallOverrides
  ): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: pure
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _reserveIn Type: uint256, Indexed: false
   * @param _reserveOut Type: uint256, Indexed: false
   */
  getAmountOut(
    _amountIn: BigNumberish,
    _reserveIn: BigNumberish,
    _reserveOut: BigNumberish,
    overrides?: ContractCallOverrides
  ): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param _amountOut Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   */
  getAmountsIn(
    _amountOut: BigNumberish,
    _path: string[],
    overrides?: ContractCallOverrides
  ): Promise<BigNumber[]>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   */
  getAmountsOut(
    _amountIn: BigNumberish,
    _path: string[],
    overrides?: ContractCallOverrides
  ): Promise<BigNumber[]>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param _factory Type: address, Indexed: false
   * @param _tokenA Type: address, Indexed: false
   * @param _tokenB Type: address, Indexed: false
   */
  pairFor(
    _factory: string,
    _tokenA: string,
    _tokenB: string,
    overrides?: ContractCallOverrides
  ): Promise<string>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: pure
   * Type: function
   * @param _amountA Type: uint256, Indexed: false
   * @param _reserveA Type: uint256, Indexed: false
   * @param _reserveB Type: uint256, Indexed: false
   */
  quote(
    _amountA: BigNumberish,
    _reserveA: BigNumberish,
    _reserveB: BigNumberish,
    overrides?: ContractCallOverrides
  ): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _tokenA Type: address, Indexed: false
   * @param _tokenB Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountAMin Type: uint256, Indexed: false
   * @param _amountBMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  removeLiquidity(
    _tokenA: string,
    _tokenB: string,
    _liquidity: BigNumberish,
    _amountAMin: BigNumberish,
    _amountBMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _token Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountTokenMin Type: uint256, Indexed: false
   * @param _amountETHMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  removeLiquidityETH(
    _token: string,
    _liquidity: BigNumberish,
    _amountTokenMin: BigNumberish,
    _amountETHMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _token Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountTokenMin Type: uint256, Indexed: false
   * @param _amountETHMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  removeLiquidityETHSupportingFeeOnTransferTokens(
    _token: string,
    _liquidity: BigNumberish,
    _amountTokenMin: BigNumberish,
    _amountETHMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _token Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountTokenMin Type: uint256, Indexed: false
   * @param _amountETHMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   * @param _approveMax Type: bool, Indexed: false
   * @param _v Type: uint8, Indexed: false
   * @param _r Type: bytes32, Indexed: false
   * @param _s Type: bytes32, Indexed: false
   */
  removeLiquidityETHWithPermit(
    _token: string,
    _liquidity: BigNumberish,
    _amountTokenMin: BigNumberish,
    _amountETHMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    _approveMax: boolean,
    _v: BigNumberish,
    _r: Arrayish,
    _s: Arrayish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _token Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountTokenMin Type: uint256, Indexed: false
   * @param _amountETHMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   * @param _approveMax Type: bool, Indexed: false
   * @param _v Type: uint8, Indexed: false
   * @param _r Type: bytes32, Indexed: false
   * @param _s Type: bytes32, Indexed: false
   */
  removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(
    _token: string,
    _liquidity: BigNumberish,
    _amountTokenMin: BigNumberish,
    _amountETHMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    _approveMax: boolean,
    _v: BigNumberish,
    _r: Arrayish,
    _s: Arrayish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _tokenA Type: address, Indexed: false
   * @param _tokenB Type: address, Indexed: false
   * @param _liquidity Type: uint256, Indexed: false
   * @param _amountAMin Type: uint256, Indexed: false
   * @param _amountBMin Type: uint256, Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   * @param _approveMax Type: bool, Indexed: false
   * @param _v Type: uint8, Indexed: false
   * @param _r Type: bytes32, Indexed: false
   * @param _s Type: bytes32, Indexed: false
   */
  removeLiquidityWithPermit(
    _tokenA: string,
    _tokenB: string,
    _liquidity: BigNumberish,
    _amountAMin: BigNumberish,
    _amountBMin: BigNumberish,
    _to: string,
    _deadline: BigNumberish,
    _approveMax: boolean,
    _v: BigNumberish,
    _r: Arrayish,
    _s: Arrayish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param _amountOut Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapETHForExactTokens(
    _amountOut: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactETHForTokens(
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: true
   * Constant: false
   * StateMutability: payable
   * Type: function
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactETHForTokensSupportingFeeOnTransferTokens(
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactTokensForETH(
    _amountIn: BigNumberish,
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactTokensForETHSupportingFeeOnTransferTokens(
    _amountIn: BigNumberish,
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactTokensForTokens(
    _amountIn: BigNumberish,
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountIn Type: uint256, Indexed: false
   * @param _amountOutMin Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapExactTokensForTokensSupportingFeeOnTransferTokens(
    _amountIn: BigNumberish,
    _amountOutMin: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountOut Type: uint256, Indexed: false
   * @param _amountInMax Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapTokensForExactETH(
    _amountOut: BigNumberish,
    _amountInMax: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param _amountOut Type: uint256, Indexed: false
   * @param _amountInMax Type: uint256, Indexed: false
   * @param _path Type: address[], Indexed: false
   * @param _to Type: address, Indexed: false
   * @param _deadline Type: uint256, Indexed: false
   */
  swapTokensForExactTokens(
    _amountOut: BigNumberish,
    _amountInMax: BigNumberish,
    _path: string[],
    _to: string,
    _deadline: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
}
