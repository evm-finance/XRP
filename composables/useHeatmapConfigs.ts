import { computed, useStore } from '@nuxtjs/composition-api'
import { HeatmapIntervals, HeatmapTileSize, State } from '~/types/state'

const numOfCoinsOptions: number[] = [10, 20, 36, 50, 100, 200, 300, 500]

const timeFrameOptions: { text: string; value: HeatmapIntervals }[] = [
  { text: '1H', value: '1h' },
  { text: '4H', value: '4h' },
  { text: 'Daily', value: 'daily' },
  { text: 'Weekly', value: 'weekly' },
  { text: 'Monthly', value: 'monthly' },
  { text: 'Yearly', value: 'yearly' },
  { text: 'Year To Date', value: 'ytd' },
]

const blockSizeOptions: { text: string; value: HeatmapTileSize }[] = [
  { text: 'QC Index', value: 'marketcap_index' },
  { text: 'MarketCap', value: 'marketcap' },
  { text: 'Price 24h ', value: 'price24hAbs' },
  { text: 'Volume 24h ', value: 'volume24h' },
  { text: 'Price USD ', value: 'price_usd' },
]

export default function () {
  const { state, dispatch } = useStore<State>()

  const blockSizeName = computed(() => blockSizeOptions.find((elem) => elem.value === blockSize.value) ?? '')

  const timeFrame = computed<HeatmapIntervals>({
    get: () => state.global.heatmap.timeFrame,
    set: (value: HeatmapIntervals) => dispatch('global/timeFrame', value),
  })
  const blockSize = computed<HeatmapTileSize>({
    get: () => state.global.heatmap.blockSize,
    set: (value: HeatmapTileSize) => dispatch('global/blockSize', value),
  })
  const blueTile = computed<boolean>({
    get: () => state.global.heatmap.blueTile,
    set: (value: boolean) => dispatch('global/blueTile', value),
  })
  const displayFavorites = computed<boolean>({
    get: () => state.global.heatmap.displayFavorites,
    set: (value: boolean) => dispatch('global/displayFavorites', value),
  })
  const numOfCoins = computed<number>({
    get: () => state.global.heatmap.numOfCoins,
    set: (value: number) => dispatch('global/numOfCoins', value),
  })
  const displayGainersAndLosers = computed<boolean>({
    get: () => state.global.heatmap.displayGainersAndLosers,
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
