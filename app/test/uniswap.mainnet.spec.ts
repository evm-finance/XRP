import { ethers } from 'ethers'
import { TradePath, UniswapV2, ROUTER } from '~/lib/uniswap-v2/uniswapV2'
import erc20Abi from '~/constracts/abi/erc20Abi.json'
import { Erc20Contract } from '~/types/abi/erc20-contract'

describe('Verse uniswap-v2 MainNet Test', () => {
  const usdcToken = {
    chainId: 1,
    address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
    symbol: 'USDC',
    name: 'USDCoin',
    decimals: 6,
  }

  const daiToken = {
    chainId: 1,
    address: '0x6b175474e89094c44da98b954eedeac495271d0f',
    symbol: 'DAI',
    name: 'Dai',
    decimals: 18,
  }

  const wethToken = {
    address: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
    chainId: 1,
    decimals: 18,
    name: 'Wrapped Ether',
    symbol: 'WETH',
  }
  const verseToken = {
    chainId: 1,
    address: '0x249ca82617ec3dfb2589c4c17ab7ec9765350a18',
    name: 'Verse',
    symbol: 'VERSE',
    decimals: 18,
  }

  const wbtc = {
    chainId: 1,
    address: '0x2260fac5e5542a773aa44fbcfedf7c193bc2c599',
    decimals: 8,
    symbol: 'WBTC',
    name: 'WBTC',
  }

  const provider = new ethers.providers.JsonRpcProvider(
    process.env.MAIN_NET_FORK_URL || ''
  )
  const signer = provider.getSigner()

  it('ETH > ERC20 EXACT_INPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(wethToken, usdcToken, TradePath.ethToErc20, userWallet, provider)
    const txData = await uniswap.trade('1', 'EXACT_INPUT')
    expect(txData).not.toBeUndefined()
    console.log(txData)
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
    await provider.getGasPrice()
  })

  it('ETH > ERC20 EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(wethToken, daiToken, TradePath.ethToErc20, userWallet, provider)
    const txData = await uniswap.trade('500', 'EXACT_OUTPUT')
    expect(txData).not.toBeUndefined()
    console.log(txData)
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ERC20 EXACT_INPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(wbtc, wethToken, TradePath.erc20ToErc20, userWallet, provider)
    const txData = await uniswap.trade('1', 'EXACT_INPUT')
    console.log(txData)
    // expect(txData).not.toBeUndefined()
    // console.log(txData)
    // txData.transaction.gasPrice = await provider.getGasPrice()
    // txData.transaction.gasLimit = 250000
    // const contract = new ethers.Contract(wethToken.address, erc20Abi, signer) as unknown as Erc20Contract
    // await contract.approve(ROUTER(1), '115792089237316195423570985008687907853269984665640564039457584007913129639935')
    // const tx = await signer.sendTransaction(txData.transaction)
    // const resp = await tx.wait()
    // console.log(resp)
    // expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ERC20 EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(daiToken, verseToken, TradePath.erc20ToErc20, userWallet, provider)
    const txData = await uniswap.trade('3000', 'EXACT_OUTPUT')
    console.log(txData)
    expect(txData).not.toBeUndefined()
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const contract = new ethers.Contract(daiToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(ROUTER(1), '115792089237316195423570985008687907853269984665640564039457584007913129639935')
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ETH  EXACT_INPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(usdcToken, wethToken, TradePath.erc20ToEth, userWallet, provider)
    const txData = await uniswap.trade('50', 'EXACT_INPUT')
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    console.log(txData)
    const contract = new ethers.Contract(usdcToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(ROUTER(1), '115792089237316195423570985008687907853269984665640564039457584007913129639935')
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ETH  EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(usdcToken, wethToken, TradePath.erc20ToEth, userWallet, provider)
    const txData = await uniswap.trade('0.1', 'EXACT_OUTPUT')
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    console.log(txData)
    const contract = new ethers.Contract(usdcToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(ROUTER(1), '115792089237316195423570985008687907853269984665640564039457584007913129639935')
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })
})
