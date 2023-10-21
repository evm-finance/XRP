import { computed, inject, reactive, ref, Ref } from '@nuxtjs/composition-api'
import { BigNumber, ethers } from 'ethers'
import type { ContractTransaction } from 'ethers'
import type { AavePoolModel, aaveVersion } from '~/composables/useAavePools'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import wethGatewayAbi from '~/constracts/abi/aave/wethGatewayAbi.json'
import lendingPoolAbi from '~/constracts/abi/aave/lendingPoolAbi.json'
import { AaveLendingPoolContract } from '~/types/abi/aave-lending-pool-contract'
import { AaveWethPoolContract } from '~/types/abi/aave-weth-pool-contract'
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'

interface MarketConfig {
  name: string
  AAVE_LENDING_POOL_ADDRESS: string
  AAVE_WETH_GATEWAY_ADDRESS: string
  AWETH_ADDRESS: string
  WETH_VARIABLE_BORROW: string
  WRAPPED_ETH: string
}

const configsV2: { [key: number]: MarketConfig } = {
  1: {
    name: 'ETH mainnet',
    AAVE_LENDING_POOL_ADDRESS: '0x7d2768dE32b0b80b7a3454c06BdAc94A69DDc7A9',
    AAVE_WETH_GATEWAY_ADDRESS: '0xcc9a0B7c43DC2a5F023Bb9b738E45B0Ef6B06E04',
    AWETH_ADDRESS: '0x030bA81f1c18d280636F32af80b9AAd02Cf0854e',
    WETH_VARIABLE_BORROW: '0xF63B34710400CAd3e044cFfDcAb00a0f32E33eCf',
    WRAPPED_ETH: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
  },

  137: {
    name: 'Polygon (Matic) ',
    AAVE_LENDING_POOL_ADDRESS: '0x8dFf5E27EA6b7AC08EbFdf9eB090F32ee9a30fcf',
    AAVE_WETH_GATEWAY_ADDRESS: '0xbEadf48d62aCC944a06EEaE0A9054A90E5A7dc97',
    AWETH_ADDRESS: '0x8dF3aad3a84da6b69A4DA8aeC3eA40d9091B2Ac4',
    WETH_VARIABLE_BORROW: '0x59e8E9100cbfCBCBAdf86b9279fa61526bBB8765',
    WRAPPED_ETH: '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270',
  },

  // 1337: {
  //   name: 'Polygon (Matic) ',
  //   AAVE_LENDING_POOL_ADDRESS: '0x8dFf5E27EA6b7AC08EbFdf9eB090F32ee9a30fcf',
  //   AAVE_WETH_GATEWAY_ADDRESS: '0xbEadf48d62aCC944a06EEaE0A9054A90E5A7dc97',
  //   AWETH_ADDRESS: '0x8dF3aad3a84da6b69A4DA8aeC3eA40d9091B2Ac4',
  //   WETH_VARIABLE_BORROW: '0x59e8E9100cbfCBCBAdf86b9279fa61526bBB8765',
  //   WRAPPED_ETH: '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270',
  // },

  31337: {
    name: 'Polygon (Matic) ',
    AAVE_LENDING_POOL_ADDRESS: '0x8dFf5E27EA6b7AC08EbFdf9eB090F32ee9a30fcf',
    AAVE_WETH_GATEWAY_ADDRESS: '0xbEadf48d62aCC944a06EEaE0A9054A90E5A7dc97',
    AWETH_ADDRESS: '0x8dF3aad3a84da6b69A4DA8aeC3eA40d9091B2Ac4',
    WETH_VARIABLE_BORROW: '0x59e8E9100cbfCBCBAdf86b9279fa61526bBB8765',
    WRAPPED_ETH: '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270',
  },
}

const configsV3: { [key: number]: MarketConfig } = {
  1: {
    name: 'ETH mainnet',
    AAVE_LENDING_POOL_ADDRESS: '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2',
    AAVE_WETH_GATEWAY_ADDRESS: '',
    AWETH_ADDRESS: '',
    WETH_VARIABLE_BORROW: '',
    WRAPPED_ETH: '',
  },

  137: {
    name: 'Polygon (Matic) ',
    AAVE_LENDING_POOL_ADDRESS: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_WETH_GATEWAY_ADDRESS: '',
    AWETH_ADDRESS: '',
    WETH_VARIABLE_BORROW: '',
    WRAPPED_ETH: '',
  },

  43114: {
    name: 'Avalanche',
    AAVE_LENDING_POOL_ADDRESS: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_WETH_GATEWAY_ADDRESS: '',
    AWETH_ADDRESS: '',
    WETH_VARIABLE_BORROW: '',
    WRAPPED_ETH: '',
  },

  10: {
    name: 'Optimism',
    AAVE_LENDING_POOL_ADDRESS: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_WETH_GATEWAY_ADDRESS: '',
    AWETH_ADDRESS: '',
    WETH_VARIABLE_BORROW: '',
    WRAPPED_ETH: '',
  },

  // 1337: {
  //   name: 'Polygon (Matic) ',
  //   AAVE_LENDING_POOL_ADDRESS: '0x8dFf5E27EA6b7AC08EbFdf9eB090F32ee9a30fcf',
  //   AAVE_WETH_GATEWAY_ADDRESS: '0xbEadf48d62aCC944a06EEaE0A9054A90E5A7dc97',
  //   AWETH_ADDRESS: '0x8dF3aad3a84da6b69A4DA8aeC3eA40d9091B2Ac4',
  //   WETH_VARIABLE_BORROW: '0x59e8E9100cbfCBCBAdf86b9279fa61526bBB8765',
  //   WRAPPED_ETH: '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270',
  // },

  31337: {
    name: 'Optimism',
    AAVE_LENDING_POOL_ADDRESS: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_WETH_GATEWAY_ADDRESS: '',
    AWETH_ADDRESS: '',
    WETH_VARIABLE_BORROW: '',
    WRAPPED_ETH: '',
  },
}

interface Transaction {
  loading: boolean
  isCompleted: boolean
  receipt: any
}

export default function useAaveTransactions(
  pool: Ref<AavePoolModel>,
  amountInTokens: Ref<number>,
  version: Ref<aaveVersion> = ref('v3')
) {
  // COMPUTED
  const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
  const { allowedSpending, approveMaxSpending } = useERC20()
  const { isNativeToken } = useHelpers()

  // STATE
  const txData = reactive<Transaction>({
    loading: false,
    isCompleted: false,
    receipt: null,
  })

  const configs = computed(() => (version.value === 'v2' ? configsV2 : configsV3))

  // COMPUTED
  const AAVE_LENDING_POOL_ADDRESS = computed(() =>
    chainId.value && chainId.value in configs.value ? configs.value[chainId.value].AAVE_LENDING_POOL_ADDRESS : ''
  )

  const AAVE_WETH_GATEWAY_ADDRESS = computed(() =>
    chainId.value && chainId.value in configs.value ? configs.value[chainId.value].AAVE_WETH_GATEWAY_ADDRESS : ''
  )

  const AWETH_ADDRESS = computed(() =>
    chainId.value && chainId.value in configs.value ? configs.value[chainId.value].AWETH_ADDRESS : ''
  )

  const WRAPPED_ETH = computed(() =>
    chainId.value && chainId.value in configs.value ? configs.value[chainId.value].WRAPPED_ETH : ''
  )

  const txLoading = computed(() => txData.loading)
  const receipt = computed(() => txData.receipt)
  const isTxMined = computed(() => !!(txData.receipt && txData.receipt.transactionHash && txData.receipt.blockNumber))

  const ERC20_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`125`).div('100').toNumber()
  const NATIVE_ETH_GAS_LIMIT = (estimatedGas: BigNumber): number => estimatedGas.mul(`175`).div('100').toNumber()
  const ESTIMATED_GAS_FEE = async (tx: ContractTransaction): Promise<BigNumber> =>
    (await provider.value?.estimateGas(<any>tx)) ?? BigNumber.from('0')

  // METHODS
  function resetToDefault() {
    txData.loading = false
    txData.isCompleted = false
    txData.receipt = null
  }

  async function depositETHWeb3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_WETH_GATEWAY_ADDRESS.value,
        wethGatewayAbi,
        signer.value
      ) as unknown as AaveWethPoolContract

      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.depositETH(
        AAVE_LENDING_POOL_ADDRESS.value,
        account.value,
        0,
        {
          value: amount,
        }
      )
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = NATIVE_ETH_GAS_LIMIT(estimatedGas)
      const depositCall = await contract.depositETH(AAVE_LENDING_POOL_ADDRESS.value, account.value, 0, {
        value: amount,
        gasLimit,
      })
      const resp = await depositCall.wait()
      return { isCompleted: true, receipt: resp }
    } catch (error) {
      return { isCompleted: false, receipt: error }
    }
  }

  async function withdrawETHWeb3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_WETH_GATEWAY_ADDRESS.value,
        wethGatewayAbi,
        signer.value
      ) as unknown as AaveWethPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.withdrawETH(
        AAVE_LENDING_POOL_ADDRESS.value,
        amount,
        account.value
      )
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = NATIVE_ETH_GAS_LIMIT(estimatedGas)
      const call = await contract.withdrawETH(AAVE_LENDING_POOL_ADDRESS.value, amount, account.value, { gasLimit })
      const resp = await call.wait()
      return { isCompleted: true, receipt: resp }
    } catch (error) {
      return { isCompleted: false, receipt: error }
    }
  }

  async function repayETHWeb3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_WETH_GATEWAY_ADDRESS.value,
        wethGatewayAbi,
        signer.value
      ) as unknown as AaveWethPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.repayETH(
        AAVE_LENDING_POOL_ADDRESS.value,
        amount,
        2,
        account.value,
        { value: amount }
      )
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = NATIVE_ETH_GAS_LIMIT(estimatedGas)
      const call = await contract.repayETH(AAVE_LENDING_POOL_ADDRESS.value, amount, 2, account.value, {
        value: amount,
        gasLimit,
      })
      const resp = await call.wait()
      return { isCompleted: true, receipt: resp }
    } catch (error) {
      return { isCompleted: false, receipt: error }
    }
  }

  async function depositERC20Web3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    console.log(AAVE_LENDING_POOL_ADDRESS.value, 'AAAAAAAAAAA')

    try {
      const contract = new ethers.Contract(
        AAVE_LENDING_POOL_ADDRESS.value,
        lendingPoolAbi,
        signer.value
      ) as unknown as AaveLendingPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.deposit(
        pool.value.underlyingAsset,
        amount,
        account.value,
        0
      )
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
      const depositCall = await contract.deposit(pool.value.underlyingAsset, amount, account.value, 0, { gasLimit })
      const resp = await depositCall.wait()
      return { isCompleted: true, receipt: resp }
    } catch (error) {
      return { isCompleted: false, receipt: error }
    }
  }

  // async function depositERC20Web3CallV3(): Promise<{ isCompleted: boolean; receipt: any }> {
  //   try {
  //     const contract = new ethers.Contract(
  //       AAVE_LENDING_POOL_ADDRESS.value,
  //       aaveLandingPoolV3,
  //       signer.value
  //     ) as unknown as AaveLendingPoolContractV3
  //     const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
  //     const populatedTx = await contract.populateTransaction.deposit(
  //       pool.value.underlyingAsset,
  //       amount,
  //       account.value,
  //       0
  //     )
  //     const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
  //     const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
  //     const depositCall = await contract.deposit(pool.value.underlyingAsset, amount, account.value, 0, { gasLimit })
  //     const resp = await depositCall.wait()
  //     return { isCompleted: true, receipt: resp }
  //   } catch (error) {
  //     return { isCompleted: false, receipt: error }
  //   }
  // }

  async function borrowERC20Web3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_LENDING_POOL_ADDRESS.value,
        lendingPoolAbi,
        signer.value
      ) as unknown as AaveLendingPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.borrow(
        pool.value.underlyingAsset,
        amount,
        2,
        0,
        account.value
      )
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
      const borrowCall = await contract.borrow(pool.value.underlyingAsset, amount, 2, 0, account.value, { gasLimit })
      const resp = await borrowCall.wait()
      return { isCompleted: true, receipt: resp }
    } catch (err) {
      return { isCompleted: false, receipt: err }
    }
  }

  async function withdrawERC20Web3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_LENDING_POOL_ADDRESS.value,
        lendingPoolAbi,
        signer.value
      ) as unknown as AaveLendingPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.withdraw(pool.value.underlyingAsset, amount, account.value)
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
      const withdrawCall = await contract.withdraw(pool.value.underlyingAsset, amount, account.value, { gasLimit })
      const resp = await withdrawCall.wait()
      return { isCompleted: true, receipt: resp }
    } catch (err) {
      return { isCompleted: false, receipt: err }
    }
  }

  async function repayERC20Web3Call(): Promise<{ isCompleted: boolean; receipt: any }> {
    try {
      const contract = new ethers.Contract(
        AAVE_LENDING_POOL_ADDRESS.value,
        lendingPoolAbi,
        signer.value
      ) as unknown as AaveLendingPoolContract
      const amount = ethers.utils.parseUnits(`${amountInTokens.value}`, pool.value.decimals)
      const populatedTx = await contract.populateTransaction.repay(pool.value.underlyingAsset, amount, 2, account.value)
      const estimatedGas: BigNumber = await ESTIMATED_GAS_FEE(populatedTx)
      const gasLimit = ERC20_GAS_LIMIT(estimatedGas)
      const repayCall = await contract.repay(pool.value.underlyingAsset, amount, 2, account.value, { gasLimit })
      const resp = await repayCall.wait()
      return { isCompleted: true, receipt: resp }
    } catch (err) {
      return { isCompleted: false, receipt: err }
    }
  }

  async function deposit() {
    txData.loading = true
    const isNative = isNativeToken(chainId.value ?? 1, pool.value.symbol)
    if (!isNative) {
      const allowed = await allowedSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      }
    }
    let result: { isCompleted: boolean; receipt: any }
    isNative ? (result = await depositETHWeb3Call()) : (result = await depositERC20Web3Call())
    txData.isCompleted = result.isCompleted
    txData.receipt = result.receipt
    txData.loading = false
  }

  async function borrow() {
    txData.loading = true
    const isNative = isNativeToken(chainId.value ?? 1, pool.value.symbol)

    if (!isNative) {
      const allowed = await allowedSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      }
    }
    const result = await borrowERC20Web3Call()
    txData.isCompleted = result.isCompleted
    txData.receipt = result.receipt
    txData.loading = false
  }

  async function repay() {
    txData.loading = true
    const isNative = isNativeToken(chainId.value ?? 1, pool.value.symbol)
    if (!isNative) {
      const allowed = await allowedSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      }
    } else {
      const allowed = await allowedSpending(WRAPPED_ETH.value, AAVE_WETH_GATEWAY_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(WRAPPED_ETH.value, AAVE_WETH_GATEWAY_ADDRESS.value)
      }
    }
    let result: { isCompleted: boolean; receipt: any }
    isNative ? (result = await repayETHWeb3Call()) : (result = await repayERC20Web3Call())
    txData.isCompleted = result.isCompleted
    txData.receipt = result.receipt
    txData.loading = false
  }

  async function withdraw() {
    txData.loading = true
    const isNative = isNativeToken(chainId.value ?? 1, pool.value.symbol)
    if (!isNative) {
      const allowed = await allowedSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(pool.value.underlyingAsset, AAVE_LENDING_POOL_ADDRESS.value)
      }
    } else {
      const allowed = await allowedSpending(AWETH_ADDRESS.value, AAVE_WETH_GATEWAY_ADDRESS.value)
      if (allowed < amountInTokens.value) {
        await approveMaxSpending(AWETH_ADDRESS.value, AAVE_WETH_GATEWAY_ADDRESS.value)
      }
    }
    let result: { isCompleted: boolean; receipt: any }
    isNative ? (result = await withdrawETHWeb3Call()) : (result = await withdrawERC20Web3Call())
    txData.isCompleted = result.isCompleted
    txData.receipt = result.receipt
    txData.loading = false
  }
  return { txLoading, receipt, isTxMined, deposit, borrow, repay, withdraw, resetToDefault }
}
