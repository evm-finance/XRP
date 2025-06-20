import { ethers } from 'ethers'
import { TradePath, UniswapV2, ROUTER } from '~/lib/uniswap-v2/uniswapV2'
import { Erc20Contract } from '~/types/abi/erc20-contract'
import erc20Abi from '~/constracts/abi/erc20Abi.json'


describe('Verse uniswap-v2 BCH Chain Test', () => {
  const wbchToken = {
    address: '0x3743eC0673453E5009310C727Ba4eaF7b3a1cc04',
    chainId: 10000,
    decimals: 18,
    name: 'Wrapped Bitcoin Cas',
    symbol: 'WBCH',
  }

  const lawToken = {
    address: '0x0b00366fbf7037e9d75e4a569ab27dab84759302',
    chainId: 10000,
    decimals: 18,
    name: 'Law Token',
    symbol: 'LAW',
  }
  const joyToken = {
    address: '0x6732e55ac3eca734f54c26bd8df4eed52fb79a6e',
    chainId: 10000,
    decimals: 18,
    name: 'Joystick',
    symbol: 'JOY',
  }

  const provider = new ethers.providers.JsonRpcProvider(process.env.MAIN_NET_FORK_URL || '')
  const signer = provider.getSigner()

  it('ETH > ERC20 EXACT_INPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(wbchToken, lawToken, TradePath.ethToErc20, userWallet, provider, 10000)
    const txData = await uniswap.trade('10', 'EXACT_INPUT')
    expect(txData).not.toBeUndefined()
    console.log(txData)
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ETH > ERC20 EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(wbchToken, joyToken, TradePath.ethToErc20, userWallet, provider, 10000)
    const txData = await uniswap.trade('1', 'EXACT_OUTPUT')
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
    const uniswap = new UniswapV2(lawToken, joyToken, TradePath.erc20ToErc20, userWallet, provider, 10000)
    const txData = await uniswap.trade('0.1', 'EXACT_INPUT')
    console.log(txData)
    expect(txData).not.toBeUndefined()
    console.log(txData)
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const contract = new ethers.Contract(lawToken.address, erc20Abi, signer) as unknown as Erc20Contract
    const balance = await contract.balanceOf(userWallet)
    console.log('BBBBB', balance)
    const approve = await contract.approve(
      ROUTER(10000),
      '115792089237316195423570985008687907853269984665640564039457584007913129639935'
    )
    // console.log(res)
    await approve.wait()
    // console.log(approvalRes)
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ERC20 EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(lawToken, joyToken, TradePath.erc20ToErc20, userWallet, provider, 10000)
    const txData = await uniswap.trade('1', 'EXACT_OUTPUT')
    console.log(txData)
    expect(txData).not.toBeUndefined()
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    const contract = new ethers.Contract(lawToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(
      ROUTER(10000),
      '115792089237316195423570985008687907853269984665640564039457584007913129639935'
    )
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ETH  EXACT_INPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(lawToken, wbchToken, TradePath.erc20ToEth, userWallet, provider, 10000)
    const txData = await uniswap.trade('0.1', 'EXACT_INPUT')
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    console.log(txData)
    const contract = new ethers.Contract(lawToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(
      ROUTER(10000),
      '115792089237316195423570985008687907853269984665640564039457584007913129639935'
    )
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })

  it('ERC20 > ETH  EXACT_OUTPUT', async () => {
    const userWallet = await signer.getAddress()
    const uniswap = new UniswapV2(lawToken, wbchToken, TradePath.erc20ToEth, userWallet, provider, 10000)
    const txData = await uniswap.trade('00.1', 'EXACT_OUTPUT')
    txData.transaction.gasPrice = await provider.getGasPrice()
    txData.transaction.gasLimit = 250000
    console.log(txData)
    const contract = new ethers.Contract(lawToken.address, erc20Abi, signer) as unknown as Erc20Contract
    await contract.approve(
      ROUTER(10000),
      '115792089237316195423570985008687907853269984665640564039457584007913129639935'
    )
    const tx = await signer.sendTransaction(txData.transaction)
    const resp = await tx.wait()
    console.log(resp)
    expect(resp).not.toBeUndefined()
  })
})
