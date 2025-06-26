<template>
  <v-menu v-model="searchMenuToggle" offset-y>
    <template #activator="{ attrs, on }">
      <v-text-field
        v-model="search"
        :loading="loading"
        clearable
        solo
        dense
        hide-details
        single-line
        placeholder="Search by Token address / Wallet / Txn Hash / XRP Address"
        v-bind="attrs"
        prepend-inner-icon="mdi-magnify"
        v-on="on"
      />
    </template>

    <v-list>
      <v-list-item v-for="(item, i) in searchResult" :key="i" exact :to="item.to">
        <v-list-item-avatar size="24">
          <v-img
            :src="getNetworkIcon(item)"
            :lazy-src="getNetworkIcon(item)"
          />
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>{{ item.network?.name || 'XRP Ledger' }}</v-list-item-title>
          <v-list-item-subtitle>{{ item.desc }}</v-list-item-subtitle>
        </v-list-item-content>

        <v-list-item-action>
          <span>
            <v-chip v-if="item.isWallet" x-small color="green darken-2">Wallet</v-chip>
            <v-chip v-if="item.isContract" x-small color="deep-purple darken-2">Contract</v-chip>
            <v-chip v-if="item.isTransaction" x-small color="blue darken-2">Transaction</v-chip>
            <v-chip v-if="item.isXRPLedger" x-small color="orange darken-2">XRP Ledger</v-chip>
            <v-chip v-if="item.isXRPAccount" x-small color="green darken-2">XRP Account</v-chip>
            <v-chip v-if="item.isXRPTransaction" x-small color="blue darken-2">XRP Transaction</v-chip>
            <v-icon
              v-if="!item.isWallet && !item.isContract && !item.isTransaction && !item.isXRPLedger && !item.isXRPAccount && !item.isXRPTransaction"
              size="18"
              color="red"
            >
              mdi-close
            </v-icon>
          </span>
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
<script lang="ts">
import { defineComponent, inject, ref, watch } from '@nuxtjs/composition-api'
import { ethers } from 'ethers'
import useERC20 from '~/composables/useERC20'
import { useHelpers } from '~/composables/useHelpers'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
import { SearchResult, Network } from '~/types/global'

export default defineComponent({
  setup() {
    const { allNetworks, getCustomProviderByNetworkId } = inject(WEB3_PLUGIN_KEY) as Web3
    const { debounceAsync } = useHelpers()
    const { getUniswapTokenByAddress } = useERC20()

    const searchResult = ref<SearchResult[]>([])
    const searchMenuToggle = ref(false)
    const dialog = ref(false)
    const search = ref<string | null>('')
    const loading = ref(false)

    // XRP address validation
    const isValidXRPAddress = (address: string): boolean => {
      // XRP addresses are typically 25-35 characters long and start with 'r'
      return /^r[a-km-zA-HJ-NP-Z1-9]{25,34}$/.test(address)
    }

    // XRP transaction hash validation
    const isValidXRPTransactionHash = (hash: string): boolean => {
      // XRP transaction hashes are 64 characters long hex strings
      return /^[A-Fa-f0-9]{64}$/.test(hash)
    }

    // XRP ledger index validation
    const isValidXRPLedgerIndex = (ledger: string): boolean => {
      // XRP ledger indices are numeric
      return /^\d+$/.test(ledger)
    }

    // Get network icon
    const getNetworkIcon = (item: SearchResult): string => {
      if (item.isXRPLedger || item.isXRPAccount || item.isXRPTransaction) {
        return '/img/xrp-icon.svg'
      }
      return this.$imageUrlBySymbol(item.network?.symbol ?? '')
    }

    const checkXRPAddress = async (address: string): Promise<SearchResult[]> => {
      const resp: SearchResult[] = []

      try {
        // Check if it's a valid XRP address
        if (isValidXRPAddress(address)) {
          resp.push({
            desc: 'XRP Account - View Balances & Transactions',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: false,
            isXRPAccount: true,
            searchString: address,
            to: `/xrp-balances?address=${address}`,
          })

          resp.push({
            desc: 'XRP Account - View Transaction History',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: false,
            isXRPAccount: true,
            isXRPTransaction: true,
            searchString: address,
            to: `/xrp-transactions?address=${address}`,
          })

          return resp
        }

        // Check if it's a valid XRP transaction hash
        if (isValidXRPTransactionHash(address)) {
          resp.push({
            desc: 'XRP Transaction - View Details',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: false,
            isXRPTransaction: true,
            searchString: address,
            to: `/xrp-explorer/tx/${address}`,
          })

          return resp
        }

        // Check if it's a valid XRP ledger index
        if (isValidXRPLedgerIndex(address)) {
          resp.push({
            desc: 'XRP Ledger - View Details',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: true,
            searchString: address,
            to: `/xrp-explorer/ledger/${address}`,
          })

          return resp
        }

        // If none of the above, return invalid
        resp.push({
          desc: 'Invalid XRP address, transaction hash, or ledger index',
          network: null,
          isWallet: false,
          isContract: false,
          isXRPLedger: false,
          searchString: address,
          to: '#',
        })

        return resp
      } catch (error) {
        return [
          {
            desc: 'Error checking XRP address',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: false,
            searchString: address,
            to: '#',
          },
        ]
      }
    }

    const checkAddressType = async (address: string, network: Network, provider: any): Promise<SearchResult[]> => {
      try {
        const resp: SearchResult[] = []

        // Legacy XRP ledger check (8 characters)
        if (address.length === 8) {
          resp.push({
            desc: 'Valid XRP ledger',
            network: null,
            isWallet: false,
            isContract: false,
            isXRPLedger: true,
            searchString: address,
            to: `/xrp-explorer/ledger/${address}`,
          })
          return resp
        }

        const isValidTx = await isValidTransactionHash(address, provider)
        if (isValidTx) {
          resp.push({
            desc: 'Valid transaction hash',
            network,
            isWallet: false,
            isContract: false,
            isTransaction: true,
            searchString: address,
            to: `/tx/${address}?chainId=${network.chainIdentifier}`,
          })
          return resp
        }

        // Get the code at the address
        const code = await provider.getCode(address)

        // If the code is "0x", it is not a contract address (wallet address)
        if (code === '0x') {
          const transactionCount = await provider.getTransactionCount(address)

          // Check ETH balance
          const balance = await provider.getBalance(address)

          if (transactionCount > 0) {
            resp.push({
              desc: 'Address is a wallet (has transaction history)',
              network,
              isWallet: true,
              isContract: false,
              searchString: address,
              to: `/portfolio/balances?wallet=${address}&chainId=${network.chainIdentifier}`,
            })
          } else if (balance.gt(ethers.constants.Zero)) {
            resp.push({
              desc: 'Address is a wallet (holds Ether)',
              network,
              isWallet: true,
              isContract: false,
              searchString: address,
              to: `/portfolio/balances?wallet=${address}&chainId=${network.chainIdentifier}`,
            })
          } else {
            resp.push({
              desc: 'Address is a wallet',
              network,
              isWallet: true,
              isContract: false,
              searchString: address,
              to: `/portfolio/balances?wallet=${address}&chainId=${network.chainIdentifier}`,
            })
          }
        } else {
          const contract = await getUniswapTokenByAddress(address, network, provider)
          resp.push({
            desc: `${contract?.name ?? 'UNKNOWN'}`,
            network,
            isWallet: false,
            isContract: true,
            searchString: address,
            to: `/token/${contract?.symbol}?contact=${address}&decimals=${contract?.decimals}`,
          })

          resp.push({
            desc: `Balances of ${contract?.name} contract`,
            network,
            isWallet: true,
            isContract: false,
            searchString: address,
            to: `/portfolio/balances?wallet=${address}&chainId=${network.chainIdentifier}`,
          })
        }
        return resp
      } catch (error) {
        return [
          {
            desc: 'Invalid address or tx hash',
            network,
            isWallet: false,
            isContract: false,
            searchString: address,
            to: '#',
          },
        ]
      }
    }

    async function isValidTransactionHash(transactionHash: string, provider: any): Promise<boolean> {
      try {
        const transaction = await provider.getTransaction(transactionHash)
        return !!transaction
      } catch (error) {
        return false
      }
    }

    const onSearch = async (search: string | null): Promise<SearchResult[]> => {
      let result: SearchResult[][] = []
      if (search === null || search.length < 1) {
        searchResult.value = []
        return []
      }

      // First check if it's an XRP address/transaction/ledger
      const xrpResults = await checkXRPAddress(search)
      if (xrpResults.length > 0 && !xrpResults[0].desc.includes('Invalid')) {
        return xrpResults
      }

      // If not XRP, check EVM networks
      const multCalls: Promise<SearchResult[]>[] = []
      allNetworks.value.forEach((elem) => {
        const provider = getCustomProviderByNetworkId(elem.id)
        const request = checkAddressType(search ?? '', elem, provider)
        multCalls.push(request)
      })
      result = await Promise.all(multCalls)
      const res = result.reduce((r, c) => r.concat(c), [])

      // If we found XRP results but they were invalid, still return them
      if (xrpResults.length > 0 && xrpResults[0].desc.includes('Invalid')) {
        return xrpResults
      }

      return res
    }

    const onSearchChange = debounceAsync(async (searchVal: string) => {
      searchResult.value = []
      searchResult.value = await onSearch(searchVal)
      searchMenuToggle.value = true
    }, 1000)

    watch(search, async (searchVal) => {
      if (searchVal == null) {
        return
      }
      loading.value = true
      await onSearchChange(searchVal)
      loading.value = false
    })

    function openDialog() {
      dialog.value = true
    }

    return {
      loading,
      searchResult,
      dialog,
      search,
      searchMenuToggle,
      onSearch,
      openDialog,
      getNetworkIcon,
    }
  },
})
</script>
<style scoped lang="css">
>>> .my-custom-dialog {
  align-self: flex-start;
}
</style>
