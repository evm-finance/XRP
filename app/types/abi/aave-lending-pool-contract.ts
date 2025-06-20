import { ContractTransaction, ContractInterface, BytesLike as Arrayish, BigNumber, BigNumberish } from 'ethers'
import { EthersContractContextV5 } from 'ethereum-abi-types-generator'

export type ContractContext = EthersContractContextV5<
  AaveLendingPoolContract,
  AaveLendingPoolContractMethodNames,
  AaveLendingPoolContractEventsContext,
  AaveLendingPoolContractEvents
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
export type AaveLendingPoolContractEvents =
  | 'Borrow'
  | 'Deposit'
  | 'FlashLoan'
  | 'LiquidationCall'
  | 'Paused'
  | 'RebalanceStableBorrowRate'
  | 'Repay'
  | 'ReserveDataUpdated'
  | 'ReserveUsedAsCollateralDisabled'
  | 'ReserveUsedAsCollateralEnabled'
  | 'Swap'
  | 'Unpaused'
  | 'Withdraw'
export interface AaveLendingPoolContractEventsContext {
  Borrow(...parameters: any): EventFilter
  Deposit(...parameters: any): EventFilter
  FlashLoan(...parameters: any): EventFilter
  LiquidationCall(...parameters: any): EventFilter
  Paused(...parameters: any): EventFilter
  RebalanceStableBorrowRate(...parameters: any): EventFilter
  Repay(...parameters: any): EventFilter
  ReserveDataUpdated(...parameters: any): EventFilter
  ReserveUsedAsCollateralDisabled(...parameters: any): EventFilter
  ReserveUsedAsCollateralEnabled(...parameters: any): EventFilter
  Swap(...parameters: any): EventFilter
  Unpaused(...parameters: any): EventFilter
  Withdraw(...parameters: any): EventFilter
}
export type AaveLendingPoolContractMethodNames =
  | 'FLASHLOAN_PREMIUM_TOTAL'
  | 'LENDINGPOOL_REVISION'
  | 'MAX_NUMBER_RESERVES'
  | 'MAX_STABLE_RATE_BORROW_SIZE_PERCENT'
  | 'borrow'
  | 'deposit'
  | 'finalizeTransfer'
  | 'flashLoan'
  | 'getAddressesProvider'
  | 'getConfiguration'
  | 'getReserveData'
  | 'getReserveNormalizedIncome'
  | 'getReserveNormalizedVariableDebt'
  | 'getReservesList'
  | 'getUserAccountData'
  | 'getUserConfiguration'
  | 'initReserve'
  | 'initialize'
  | 'liquidationCall'
  | 'paused'
  | 'rebalanceStableBorrowRate'
  | 'repay'
  | 'setConfiguration'
  | 'setPause'
  | 'setReserveInterestRateStrategyAddress'
  | 'setUserUseReserveAsCollateral'
  | 'swapBorrowRateMode'
  | 'withdraw'
export interface BorrowEventEmittedResponse {
  reserve: string
  user: string
  onBehalfOf: string
  amount: BigNumberish
  borrowRateMode: BigNumberish
  borrowRate: BigNumberish
  referral: BigNumberish
}
export interface DepositEventEmittedResponse {
  reserve: string
  user: string
  onBehalfOf: string
  amount: BigNumberish
  referral: BigNumberish
}
export interface FlashLoanEventEmittedResponse {
  target: string
  initiator: string
  asset: string
  amount: BigNumberish
  premium: BigNumberish
  referralCode: BigNumberish
}
export interface LiquidationCallEventEmittedResponse {
  collateralAsset: string
  debtAsset: string
  user: string
  debtToCover: BigNumberish
  liquidatedCollateralAmount: BigNumberish
  liquidator: string
  receiveAToken: boolean
}
export interface RebalanceStableBorrowRateEventEmittedResponse {
  reserve: string
  user: string
}
export interface RepayEventEmittedResponse {
  reserve: string
  user: string
  repayer: string
  amount: BigNumberish
}
export interface ReserveDataUpdatedEventEmittedResponse {
  reserve: string
  liquidityRate: BigNumberish
  stableBorrowRate: BigNumberish
  variableBorrowRate: BigNumberish
  liquidityIndex: BigNumberish
  variableBorrowIndex: BigNumberish
}
export interface ReserveUsedAsCollateralDisabledEventEmittedResponse {
  reserve: string
  user: string
}
export interface ReserveUsedAsCollateralEnabledEventEmittedResponse {
  reserve: string
  user: string
}
export interface SwapEventEmittedResponse {
  reserve: string
  user: string
  rateMode: BigNumberish
}
export interface WithdrawEventEmittedResponse {
  reserve: string
  user: string
  to: string
  amount: BigNumberish
}
export interface ReserveconfigurationmapResponse {
  data: BigNumber
  0: BigNumber
}
export interface ConfigurationResponse {
  data: BigNumber
  0: ConfigurationResponse
}
export interface ReservedataResponse {
  configuration: ConfigurationResponse
  0: ConfigurationResponse
  liquidityIndex: BigNumber
  1: BigNumber
  variableBorrowIndex: BigNumber
  2: BigNumber
  currentLiquidityRate: BigNumber
  3: BigNumber
  currentVariableBorrowRate: BigNumber
  4: BigNumber
  currentStableBorrowRate: BigNumber
  5: BigNumber
  lastUpdateTimestamp: number
  6: number
  aTokenAddress: string
  7: string
  stableDebtTokenAddress: string
  8: string
  variableDebtTokenAddress: string
  9: string
  interestRateStrategyAddress: string
  10: string
  id: number
  11: number
}
export interface GetUserAccountDataResponse {
  totalCollateralETH: BigNumber
  0: BigNumber
  totalDebtETH: BigNumber
  1: BigNumber
  availableBorrowsETH: BigNumber
  2: BigNumber
  currentLiquidationThreshold: BigNumber
  3: BigNumber
  ltv: BigNumber
  4: BigNumber
  healthFactor: BigNumber
  5: BigNumber
  length: 6
}
export interface UserconfigurationmapResponse {
  data: BigNumber
  0: BigNumber
}
export interface AaveLendingPoolContract {
  populateTransaction: any
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  FLASHLOAN_PREMIUM_TOTAL(overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  LENDINGPOOL_REVISION(overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  MAX_NUMBER_RESERVES(overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  MAX_STABLE_RATE_BORROW_SIZE_PERCENT(overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param interestRateMode Type: uint256, Indexed: false
   * @param referralCode Type: uint16, Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   */
  borrow(
    asset: string,
    amount: BigNumberish,
    interestRateMode: BigNumberish,
    referralCode: BigNumberish,
    onBehalfOf: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   * @param referralCode Type: uint16, Indexed: false
   */
  deposit(
    asset: string,
    amount: BigNumberish,
    onBehalfOf: string,
    referralCode: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param from Type: address, Indexed: false
   * @param to Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param balanceFromBefore Type: uint256, Indexed: false
   * @param balanceToBefore Type: uint256, Indexed: false
   */
  finalizeTransfer(
    asset: string,
    from: string,
    to: string,
    amount: BigNumberish,
    balanceFromBefore: BigNumberish,
    balanceToBefore: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param receiverAddress Type: address, Indexed: false
   * @param assets Type: address[], Indexed: false
   * @param amounts Type: uint256[], Indexed: false
   * @param modes Type: uint256[], Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   * @param params Type: bytes, Indexed: false
   * @param referralCode Type: uint16, Indexed: false
   */
  flashLoan(
    receiverAddress: string,
    assets: string[],
    amounts: BigNumberish[],
    modes: BigNumberish[],
    onBehalfOf: string,
    params: Arrayish,
    referralCode: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  getAddressesProvider(overrides?: ContractCallOverrides): Promise<string>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param asset Type: address, Indexed: false
   */
  getConfiguration(asset: string, overrides?: ContractCallOverrides): Promise<ReserveconfigurationmapResponse>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param asset Type: address, Indexed: false
   */
  getReserveData(asset: string, overrides?: ContractCallOverrides): Promise<ReservedataResponse>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param asset Type: address, Indexed: false
   */
  getReserveNormalizedIncome(asset: string, overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param asset Type: address, Indexed: false
   */
  getReserveNormalizedVariableDebt(asset: string, overrides?: ContractCallOverrides): Promise<BigNumber>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  getReservesList(overrides?: ContractCallOverrides): Promise<string[]>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param user Type: address, Indexed: false
   */
  getUserAccountData(user: string, overrides?: ContractCallOverrides): Promise<GetUserAccountDataResponse>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   * @param user Type: address, Indexed: false
   */
  getUserConfiguration(user: string, overrides?: ContractCallOverrides): Promise<UserconfigurationmapResponse>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param aTokenAddress Type: address, Indexed: false
   * @param stableDebtAddress Type: address, Indexed: false
   * @param variableDebtAddress Type: address, Indexed: false
   * @param interestRateStrategyAddress Type: address, Indexed: false
   */
  initReserve(
    asset: string,
    aTokenAddress: string,
    stableDebtAddress: string,
    variableDebtAddress: string,
    interestRateStrategyAddress: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param provider Type: address, Indexed: false
   */
  initialize(provider: string, overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param collateralAsset Type: address, Indexed: false
   * @param debtAsset Type: address, Indexed: false
   * @param user Type: address, Indexed: false
   * @param debtToCover Type: uint256, Indexed: false
   * @param receiveAToken Type: bool, Indexed: false
   */
  liquidationCall(
    collateralAsset: string,
    debtAsset: string,
    user: string,
    debtToCover: BigNumberish,
    receiveAToken: boolean,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: true
   * StateMutability: view
   * Type: function
   */
  paused(overrides?: ContractCallOverrides): Promise<boolean>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param user Type: address, Indexed: false
   */
  rebalanceStableBorrowRate(
    asset: string,
    user: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param rateMode Type: uint256, Indexed: false
   * @param onBehalfOf Type: address, Indexed: false
   */
  repay(
    asset: string,
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
   * @param asset Type: address, Indexed: false
   * @param configuration Type: uint256, Indexed: false
   */
  setConfiguration(
    asset: string,
    configuration: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param val Type: bool, Indexed: false
   */
  setPause(val: boolean, overrides?: ContractTransactionOverrides): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param rateStrategyAddress Type: address, Indexed: false
   */
  setReserveInterestRateStrategyAddress(
    asset: string,
    rateStrategyAddress: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param useAsCollateral Type: bool, Indexed: false
   */
  setUserUseReserveAsCollateral(
    asset: string,
    useAsCollateral: boolean,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param rateMode Type: uint256, Indexed: false
   */
  swapBorrowRateMode(
    asset: string,
    rateMode: BigNumberish,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
  /**
   * Payable: false
   * Constant: false
   * StateMutability: nonpayable
   * Type: function
   * @param asset Type: address, Indexed: false
   * @param amount Type: uint256, Indexed: false
   * @param to Type: address, Indexed: false
   */
  withdraw(
    asset: string,
    amount: BigNumberish,
    to: string,
    overrides?: ContractTransactionOverrides
  ): Promise<ContractTransaction>
}
