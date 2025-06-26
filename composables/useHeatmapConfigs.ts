import { computed, useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'

export type HeatmapIntervals = '1h' | '4h' | '1d' | '1w' | '1m' | '3m' | '1y'
export type HeatmapTileSize = 'marketcap' | 'volume' | 'price' | 'liquidity' | 'tvl' | 'price24hAbs' | 'volume24h' | 'price_usd'

const numOfCoinsOptions: number[] = [10, 20, 36, 50, 100, 200, 300, 500]

const timeFrameOptions: { text: string; value: HeatmapIntervals }[] = [
  { text: '1 Hour', value: '1h' },
  { text: '4 Hours', value: '4h' },
  { text: '1 Day', value: '1d' },
  { text: '1 Week', value: '1w' },
  { text: '1 Month', value: '1m' },
  { text: '3 Months', value: '3m' },
  { text: '1 Year', value: '1y' },
]

const blockSizeOptions: { text: string; value: HeatmapTileSize }[] = [
  { text: 'Market Cap', value: 'marketcap' },
  { text: 'Volume', value: 'volume' },
  { text: 'Price', value: 'price' },
  { text: 'Liquidity', value: 'liquidity' },
  { text: 'TVL', value: 'tvl' },
  { text: 'Price 24h ', value: 'price24hAbs' },
  { text: 'Volume 24h ', value: 'volume24h' },
  { text: 'Price USD ', value: 'price_usd' },
]

export default function () {
  const { state, dispatch } = useStore<State>()

  const blockSizeName = computed(() => blockSizeOptions.find((elem) => elem.value === blockSize.value) ?? '')

  const timeFrame = computed({
    get: () => (state.global?.heatmap?.timeFrame ?? '1d') as HeatmapIntervals,
    set: (value: HeatmapIntervals) => dispatch('global/timeFrame', value),
  })
  const blockSize = computed({
    get: () => (state.global?.heatmap?.blockSize ?? 'marketcap') as HeatmapTileSize,
    set: (value: HeatmapTileSize) => dispatch('global/blockSize', value),
  })
  const blueTile = computed({
    get: () => state.global?.heatmap?.blueTile ?? false,
    set: (value: boolean) => dispatch('global/blueTile', value),
  })
  const displayFavorites = computed({
    get: () => state.global?.heatmap?.displayFavorites ?? false,
    set: (value: boolean) => dispatch('global/displayFavorites', value),
  })
  const numOfCoins = computed({
    get: () => state.global?.heatmap?.numOfCoins ?? 100,
    set: (value: number) => dispatch('global/numOfCoins', value),
  })
  const displayGainersAndLosers = computed({
    get: () => state.global?.heatmap?.displayGainersAndLosers ?? false,
    set: (value: boolean) => dispatch('global/displayGainersAndLosers', value),
  })

  return {
    timeFrame,
    timeFrameOptions,
    blockSize,
    blockSizeOptions,
    blockSizeName,
    blueTile,
    displayFavorites,
    numOfCoins,
    numOfCoinsOptions,
    displayGainersAndLosers,
  }
}
