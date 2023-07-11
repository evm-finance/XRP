import { computed, ref, useContext, inject, Ref, watch } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { UniswapToken, UniswapTokenSource } from '@/types/apollo/main/types'
import { UniswapTokensGQL } from '~/apollo/main/config.query.graphql'
import { WEB3_PLUGIN_KEY, Web3 } from '~/plugins/web3/web3'
import { SUPPORTED_CHAINS as VerseSupportedChains } from '~/lib/uniswap-v2/uniswapV2'

interface SwapTokenInterface {
  name: string
  source: UniswapTokenSource
  image: string
  active: boolean
}

export default function (tokenListType: Ref<'uniswap' | 'verse'> = ref('uniswap')) {
  const { $imageUrlBySymbol } = useContext()
  const { account, chainId } = inject(WEB3_PLUGIN_KEY) as Web3

  // STATE
  const loading = ref<boolean>(true)
  const searchString = ref<string>('')
  const pageSize = ref<number>(20)
  const pageNumber = ref<number>(0)
  const sources = ref<SwapTokenInterface[]>([
    // {
    //   name: 'Your Wallet',
    //   source: UniswapTokenSource.UserBalances,
    //   image: $imageUrlBySymbol('wallet'),
    //   active: true,
    // },
    {
      name: 'CoinGecko',
      source: UniswapTokenSource.CoinGecko,
      image: $imageUrlBySymbol('coingecko'),
      active: false,
    },
    {
      name: 'Uniswap',
      source: UniswapTokenSource.Uniswap,
      image: $imageUrlBySymbol('uni'),
      active: true,
    },
    {
      name: 'Aave',
      source: UniswapTokenSource.Aave,
      image: $imageUrlBySymbol('aave'),
      active: false,
    },
  ])

  // COMPOSABLES
  const source = computed(() => sources.value.filter((elem) => elem.active === true).map(({ source }) => source))
  const { result, onResult } = useQuery(
    UniswapTokensGQL,
    () => ({
      chainId: 1,
      userWallet: account.value ?? '',
      searchString: searchString.value ?? '',
      source: source.value,
      pageNumber: pageNumber.value,
      pageSize: pageSize.value,
    }),
    { fetchPolicy: 'no-cache' }
  )

  const verseTokenList = (chain: number) => {
    const token: { [key in VerseSupportedChains]: UniswapToken[] } = {
      1: [
        {
          chainId: 1,
          address: '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee',
          symbol: 'ETH',
          name: 'Ether',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x4d224452801aced8b2f0aebe155379bb5d594381',
          name: 'ApeCoin',
          symbol: 'APE',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce',
          name: 'SHIBA INU',
          symbol: 'SHIB',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x6b175474e89094c44da98b954eedeac495271d0f',
          name: 'Dai Stablecoin',
          symbol: 'DAI',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
          name: 'USD Coin',
          symbol: 'USDC',
          decimals: 6,
        },
        {
          chainId: 1,
          address: '0xdac17f958d2ee523a2206206994597c13d831ec7',
          name: 'Tether USD',
          symbol: 'USDT',
          decimals: 6,
        },
        {
          chainId: 1,
          address: '0x249ca82617ec3dfb2589c4c17ab7ec9765350a18',
          name: 'Verse',
          symbol: 'VERSE',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0xd9c2d319cd7e6177336b0a9c93c21cb48d84fb54',
          name: 'HAPI',
          symbol: 'HAPI',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x0f5d2fb29fb7d3cfee444a200298f468908cc942',
          name: 'Decentraland MANA',
          symbol: 'MANA',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x2260fac5e5542a773aa44fbcfedf7c193bc2c599',
          name: 'Wrapped BTC',
          symbol: 'WBTC',
          decimals: 8,
        },
        {
          chainId: 1,
          address: '0xf629cbd94d3791c9250152bd8dfbdf380e2a3b9c',
          name: 'Enjin Coin',
          symbol: 'ENJ',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e',
          name: 'yearn.finance',
          symbol: 'YFI',
          decimals: 18,
        },
        {
          chainId: 1,
          address: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
          name: 'Wrapped Ether',
          symbol: 'WETH',
          decimals: 18,
        },
      ],
      10000: [
        {
          address: '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee',
          name: 'Bitcoin Cash',
          symbol: 'BCH',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x3743eC0673453E5009310C727Ba4eaF7b3a1cc04',
          name: 'Wrapped Bitcoin Cash',
          symbol: 'WBCH',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x73BE9c8Edf5e951c9a0762EA2b1DE8c8F38B5e91',
          name: 'Tango Token',
          symbol: 'TANGO',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x0b00366fbf7037e9d75e4a569ab27dab84759302',
          name: 'Law Token',
          symbol: 'LAW',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x77CB87b57F54667978Eb1B199b28a0db8C8E1c0B',
          name: 'Green Ben',
          symbol: 'EBEN',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x4b85a666dec7c959e88b97814e46113601b07e57',
          name: 'GoCrypto',
          symbol: 'GOC',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x7b2b3c5308ab5b2a1d9a94d20d35ccdf61e05b72',
          name: 'flexUSD',
          symbol: 'flexUSD',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x98dd7ec28fb43b3c4c770ae532417015fa939dd3',
          name: 'FLEX Coin',
          symbol: 'FLEX',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x6732e55ac3eca734f54c26bd8df4eed52fb79a6e',
          name: 'Joystick',
          symbol: 'JOY',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0xb5b1939ef0a3743d0ae9282dba62312b614a5ac0',
          name: 'Potato',
          symbol: 'POTA',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x49f9ecf126b6ddf51c731614755508a4674ba7ed',
          name: 'Xolos',
          symbol: 'RMZ',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x5fa664f69c2a4a3ec94fac3cbf7049bd9ca73129',
          name: 'MistToken',
          symbol: 'MIST',
          chainId: 10000,
          decimals: 18,
        },
        {
          address: '0x35b3ee79e1a7775ce0c11bd8cd416630e07b0d6f',
          name: 'Bitcoin Cash Name Service',
          symbol: 'LNS',
          chainId: 10000,
          decimals: 18,
        },
      ],
    }
    return chain in token ? token[<VerseSupportedChains>chain] : token[10000]
  }

  const isUniswapType = computed(() => tokenListType.value === 'uniswap')
  const hasMoreElements = computed(() =>
    isUniswapType.value ? result.value?.uniswapTokens.pagination.hasMore ?? false : false
  ) as Ref<boolean>
  const rawTokenList = ref<UniswapToken[]>([])
  const tokenList = computed(() => (isUniswapType.value ? rawTokenList.value : verseTokenList(chainId.value ?? 1)))

  // EVENTS
  onResult((queryResult) => {
    queryResult.data.uniswapTokens.items.forEach((elem: UniswapToken) => rawTokenList.value.push(elem))
    loading.value = queryResult.loading
  })

  // METHODS
  function onSearchChange(val: string) {
    loading.value = true
    const searchTimeout: any = null
    clearTimeout(searchTimeout)
    setTimeout(() => {
      searchString.value = val
      rawTokenList.value = []
    }, 200)
  }

  function nextPage() {
    pageNumber.value++
  }

  watch(source, () => {
    searchString.value = ''
    rawTokenList.value = []
  })

  return {
    loading,
    searchString,
    tokenList,
    sources,
    source,
    hasMoreElements,
    onSearchChange,
    nextPage,
  }
}
