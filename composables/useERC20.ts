import { inject } from '@nuxtjs/composition-api'
import { ethers } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import erc20Abi from '~/constracts/abi/erc20Abi.json'
import { Erc20Contract } from '~/types/abi/erc20-contract'
import { useHelpers } from '~/composables/useHelpers'
import { Chain, UniswapToken } from '~/types/apollo/main/types'
import { ERC20Balance } from '~/types/global'

const MAX_UINT256 = '115792089237316195423570985008687907853269984665640564039457584007913129639935'

export default function useERC20() {
  const { account, signer } = inject(WEB3_PLUGIN_KEY) as Web3
  const { isNativeToken } = useHelpers()

  const allowedSpending = async (tokenAddress: string, routerToApproveAddress: string): Promise<number> => {
    try {
      const contract = new ethers.Contract(tokenAddress, erc20Abi, signer.value) as unknown as Erc20Contract
      const checkApprove = await contract.allowance(account.value, routerToApproveAddress)
      return parseInt(checkApprove.toBigInt().toString())
    } catch (err) {
      return 0
    }
  }
  const getUniswapTokenByAddress = async (
    address: string,
    network: Chain,
    provider: any
  ): Promise<UniswapToken | null> => {
    try {
      const contract = new ethers.Contract(address, erc20Abi, provider) as unknown as Erc20Contract
      const symbol = await contract.symbol()
      const name = await contract.name()
      const decimals = await contract.decimals()
      return { address, chainId: network.chainIdentifier, name, symbol, decimals }
    } catch (e) {
      console.log(e)
      return null
    }
  }

  const approveMaxSpending = async (
    tokenAddress: string,
    routerToApproveAddress: string
  ): Promise<{ isCompleted: boolean; receipt: any }> => {
    try {
      const contract = new ethers.Contract(tokenAddress, erc20Abi, signer.value) as unknown as Erc20Contract
      const data = await contract.approve(routerToApproveAddress, MAX_UINT256)
      return { isCompleted: true, receipt: data }
    } catch (error) {
      return { isCompleted: false, receipt: error }
    }
  }

  async function ERC20Balance(
    tokenAddress: string,
    decimals: number,
    uniqueKeyOrSymbol: string,
    account: string,
    provider: any
  ): Promise<ERC20Balance> {
    try {
      const contract = new ethers.Contract(tokenAddress, erc20Abi, provider) as unknown as Erc20Contract
      const bal = await contract.balanceOf(account)
      const balance = parseFloat(bal.toString()) / 10 ** decimals
      return { address: tokenAddress, balance, isNative: false, isWrapped: false, uniqueKeyOrSymbol }
    } catch (error) {
      return { address: tokenAddress, balance: 0, isNative: false, isWrapped: false, uniqueKeyOrSymbol }
    }
  }

  /*
  async function ERC20Balance(
    tokenAddress: string,
    decimals: number,
    uniqueKeyOrSymbol?: string
  ): Promise<ERC20Balance> {
    try {
      const contract = new ethers.Contract(tokenAddress, erc20Abi, signer.value) as unknown as Erc20Contract
      const bal = await contract.balanceOf(account.value)
      const balance = parseFloat(bal.toString()) / 10 ** decimals
      return { address: tokenAddress, balance, isNative: false, uniqueKeyOrSymbol }
    } catch (error) {
      return { address: tokenAddress, balance: 0, isNative: false, uniqueKeyOrSymbol }
    }
  }
  */

  async function nativeNetworkBalance(
    address: string,
    uniqueKeyOrSymbol: string,
    provider: any,
    account: string
  ): Promise<ERC20Balance> {
    try {
      const balanceInWei = await provider.getBalance(account)
      const balanceInEth = ethers.utils.formatUnits(balanceInWei)
      return {
        address,
        balance: parseFloat(balanceInEth),
        isNative: true,
        uniqueKeyOrSymbol,
      }
    } catch (error) {
      return { address: account, balance: 0, isNative: true, uniqueKeyOrSymbol }
    }
  }

  /*  async function nativeNetworkBalance(address: string, uniqueKeyOrSymbol: string): Promise<ERC20Balance> {
    try {
      const balanceInWei = (await provider.value?.getBalance(account.value)) ?? BigNumber.from('0')
      const balanceInEth = ethers.utils.formatUnits(balanceInWei)
      return {
        address,
        balance: parseFloat(balanceInEth),
        isNative: true,
        uniqueKeyOrSymbol,
      }
    } catch (error) {
      return { address: account.value, balance: 0, isNative: true, uniqueKeyOrSymbol }
    }
  } */

  async function balanceMulticall(
    tokens: UniswapToken[],
    account: string,
    provider: any,
    chainId: number = 1
  ): Promise<ERC20Balance[]> {
    const multCalls: Promise<ERC20Balance>[] = []
    tokens.forEach((token) => {
      const isNative = isNativeToken(chainId, token.symbol)
      const request = isNative
        ? nativeNetworkBalance(token.address, `NATIVE_${token.symbol}`, provider, account)
        : ERC20Balance(token.address, token.decimals, token.symbol, account, provider)
      multCalls.push(request)
    })
    return await Promise.all(multCalls)
  }

  const isWrappedToken = (tokenAddress: string, network: Chain): boolean => {
    switch (network.id) {
      case 'ethereum':
        return tokenAddress.toLowerCase() === network.weth.address.toLowerCase()
      case 'polygon-pos':
        return tokenAddress.toLowerCase() === network.weth.address.toLowerCase()
      case 'binance-smart-chain':
        return tokenAddress.toLowerCase() === network.weth.address.toLowerCase()
      default:
        return false
    }
  }

  const getBalance = async (
    tokenAddress: string,
    symbol: string,
    decimals: number,
    network: Chain,
    account: string,
    provider: any
  ) => {
    let resp: ERC20Balance
    const isWrapped = isWrappedToken(tokenAddress, network)
    if (isWrapped) {
      const native = await nativeNetworkBalance(tokenAddress, `NATIVE_${symbol}`, provider, account)
      const erc20Bal = await ERC20Balance(tokenAddress, decimals, symbol, account, provider)
      resp = {
        isNative: false,
        isWrapped: true,
        address: tokenAddress,
        balance: erc20Bal.balance,
        nativeBalance: native.balance,
        uniqueKeyOrSymbol: symbol,
        nativeSymbol: network.symbol,
      }
    } else {
      resp = await ERC20Balance(tokenAddress, decimals, symbol, account, provider)
    }
    return resp
  }

  async function poolBalanceMulticall(
    tokens: UniswapToken[],
    network: Chain,
    account: string,
    provider: any
  ): Promise<ERC20Balance[]> {
    const multCalls: Promise<any>[] = []
    tokens.forEach((token) => {
      const request = getBalance(token.address, token.symbol, token.decimals, network, account, provider)
      multCalls.push(request)
    })
    return await Promise.all(multCalls)
  }

  return {
    allowedSpending,
    approveMaxSpending,
    ERC20Balance,
    nativeNetworkBalance,
    balanceMulticall,
    poolBalanceMulticall,
    getUniswapTokenByAddress,
  }
}
