export type AaveAddress = {
  aTokenAddress: string
  aTokenSymbol: string
  address: string
  decimals: number
  stableDebtTokenAddress: string
  symbol: string
  variableDebtTokenAddress: string
}

export type AavePoolPrice = {
  id: string
  priceInEth: number
  priceUsd: number
}

export type AavePortfolio = {
  stableBorrow: number
  symbol: string
  totalDeposits: number
  variableBorrow: number
  walletBal: number
  nativeBalance: number
  isWrapped: boolean
  networkSymbol: string
}

export type AavePool = {
  aEmissionPerSecond: number
  addresses: AaveAddress
  availableLiquidity: number
  baseLTVasCollateral: number
  borrowingEnabled: boolean
  decimals: number
  id: string
  liquidityRate: number
  name: string
  portfolioVal: AavePortfolio
  price: AavePoolPrice
  reserveLiquidationBonus: number
  reserveLiquidationThreshold: number
  sEmissionPerSecond: number
  stableBorrowRate: number
  stableBorrowRateEnabled: boolean
  symbol: string
  totalATokenSupply: number
  totalCurrentVariableDebt: number
  totalLiquidity: number
  totalLiquidityAsCollateral: number
  totalPrincipalStableDebt: number
  underlyingAsset: string
  usageAsCollateralEnabled: boolean
  utilizationRate: number
  vEmissionPerSecond: number
  variableBorrowRate: number
}
