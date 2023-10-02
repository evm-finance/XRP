export interface Addresses {
  aTokenAddress: string
  aTokenSymbol: string
  stableDebtTokenAddress: string
  variableDebtTokenAddress: string
  decimals: number
  address: string
  symbol: string
}

export interface PortfolioVal {
  totalDeposits: number
  walletBal: number
  stableBorrow: number
  variableBorrow: number
}

export interface Price {
  id: string
  priceInEth: number
  priceUsd: number
}

export interface AaveMarket {
  id: string
  underlyingAsset: string
  name: string
  symbol: string
  decimals: number
  totalLiquidity: number
  price: Price
  liquidityRate: number
  stableBorrowRate: number
  variableBorrowRate: number
  aEmissionPerSecond: number
  vEmissionPerSecond: number
  sEmissionPerSecond: number
  availableLiquidity: number
  utilizationRate: number
  totalATokenSupply: number
  totalCurrentVariableDebt: number
  totalPrincipalStableDebt: number
  addresses: Addresses
  portfolioVal: PortfolioVal
  totalLiquidityAsCollateral: number
  baseLTVasCollateral: number
  reserveLiquidationThreshold: number
  reserveLiquidationBonus: number
  usageAsCollateralEnabled: boolean
  borrowingEnabled: boolean
  stableBorrowRateEnabled: boolean
}
