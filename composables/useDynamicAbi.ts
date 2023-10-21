// This composable library will perfom the following functions
// 1. Hold the rules for ParamInputForm components
// 2. Compute function selector based on ABI
// 3. Prepare function call using ethersjs, taking inputs from Param Forms
import type { ContractTransaction } from 'ethers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { AbiElem, CalldataAbi, ContractAbi } from '~/types/apollo/main/types'
import { computed, inject, reactive, Ref } from '@nuxtjs/composition-api'
import { BigNumber, ethers } from 'ethers'
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'

//

interface Transaction {
  loading: boolean
  isCompleted: boolean
  receipt: any
}


export default function useDynamicAbi()
{

  const { signer, account, chainId, provider } = inject(WEB3_PLUGIN_KEY) as Web3
  const { allowedSpending, approveMaxSpending } = useERC20()

  const txData = reactive<Transaction>({
    loading: false,
    isCompleted: false,
    receipt: null,
  })

  const txLoading = computed(() => txData.loading)
  const receipt = computed(() => txData.receipt)
  const isTxMined = computed(() => !!(txData.receipt && txData.receipt.transactionHash && txData.receipt.blockNumber))

  const ESTIMATED_GAS_FEE = async (tx: ContractTransaction): Promise<BigNumber> =>
    (await provider.value?.estimateGas(<any>tx)) ?? BigNumber.from('0')

    // METHODS
    function resetToDefault() {
      txData.loading = false
      txData.isCompleted = false
      txData.receipt = null
    }



}
