import { ComputedRef, inject } from '@nuxtjs/composition-api'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import useERC20 from '~/composables/useERC20'
import { AaveAddress, AavePortfolio } from '~/types/aave'

export type PortfolioMap = { [poolId: string]: AavePortfolio }

export default function (addresses: ComputedRef<{ [id: string]: AaveAddress }>) {
  // COMPOSABLES
  const { account, provider, isWrapped, currentNetwork } = inject(WEB3_PLUGIN_KEY) as Web3
  const { ERC20Balance, nativeNetworkBalance } = useERC20()

  const getPoolPortfolio = async (
    poolId: string,
    tokenAddress: string,
    aTokenAddress: string,
    variableBorrowTokenAddress: string,
    decimals: number,
    account: string,
    provider: any,
    isWrapped: boolean = false,
    symbol: string = '',
    networkSymbol: string = ''
  ): Promise<{ key: string; value: AavePortfolio }> => {
    let nativeBalance = 0
    const tokenBalance = await ERC20Balance(tokenAddress, decimals, '', account, provider)
    const depositBalance = await ERC20Balance(aTokenAddress, decimals, '', account, provider)
    const variableBorrowBalance = await ERC20Balance(variableBorrowTokenAddress, decimals, '', account, provider)

    if (isWrapped) {
      const nativeBal = await nativeNetworkBalance('', '', provider, account)
      nativeBalance = nativeBal.balance
    }

    const portfolio: AavePortfolio = {
      symbol,
      networkSymbol,
      nativeBalance,
      isWrapped,
      walletBal: tokenBalance.balance,
      totalDeposits: depositBalance.balance,
      stableBorrow: 0,
      variableBorrow: variableBorrowBalance.balance,
    }
    return { key: poolId, value: portfolio }
  }

  const fetchPortfolioNew = async (): Promise<PortfolioMap> => {
    const multiCalls: Promise<{ key: string; value: AavePortfolio }>[] = []

    Object.keys(addresses.value).forEach((key) => {
      const tokenAddress = addresses.value[key].address
      const symbol = addresses.value[key].symbol
      const aTokenAddress = addresses.value[key].aTokenAddress
      const variableBorrowTokenAddress = addresses.value[key].variableDebtTokenAddress
      const decimals = addresses.value[key].decimals

      if (currentNetwork.value) {
        const isWrappedToken = isWrapped(tokenAddress, currentNetwork.value)
        multiCalls.push(
          getPoolPortfolio(
            key,
            tokenAddress,
            aTokenAddress,
            variableBorrowTokenAddress,
            decimals,
            account.value,
            provider.value,
            isWrappedToken,
            symbol,
            currentNetwork.value?.nativeTokenSymbol
          )
        )
      }
    })

    const p = await Promise.all(multiCalls)
    return p.reduce((elem, item) => ({ ...elem, [item.key]: item.value }), {})
  }

  return { fetchPortfolioNew }
}
