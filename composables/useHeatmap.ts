import { computed, onBeforeUnmount, onMounted, Ref, ref, useContext, useStore, watch } from '@nuxtjs/composition-api'
import { APIEndpoints, LongPooling } from '~/types/enums'
import useCurrency from '~/composables/useCurrency'
import type { HeatmapData, HeatmapRowData, HeatmapUpdateData } from '~/types/heatmap'
import type { HeatmapIntervals } from '~/types/state'
import { State } from '~/types/state'
import useExchange from '~/composables/useExchange'
import useHeatmapConfigs from '~/composables/heatmap/useHeatmapConfigs'
import useBrowser from '~/composables/useBrowser'
const setColor = (x: number, blueTile: boolean) => {
  if (x * 100 > 0 && x * 100 <= 1) {
    return blueTile ? '#1A2A52' : '#71c175'
  } else if (x * 100 > 1 && x * 100 <= 2.5) {
    return blueTile ? '#234e91' : '#4eb153'
  } else if (x * 100 > 2.5 && x * 100 <= 5) {
    return blueTile ? '#336cc2' : '#3e8e42'
  } else if (x * 100 > 5) {
    return blueTile ? '#5898f8' : '#2f6a32'
  } else if (x * 100 <= 0 && x * 100 >= -1) {
    return blueTile ? '#EF9A9A' : '#ff8080'
  } else if (x * 100 < -1 && x * 100 >= -2.5) {
    return blueTile ? '#EF5350' : '#ff4d4d'
  } else if (x * 100 < -2.5 && x * 100 >= -5) {
    return blueTile ? '#D32F2F' : '#ff1a1a'
  } else if (x * 100 < -5) {
    return blueTile ? '#B71C1C' : '#e60000'
  }
  return ''
}

const tileConfigs: { [key in HeatmapIntervals]: { text: string; toolTip: string } } = {
  '1h': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price1h} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              4H: {price4h}%
              7 DAY: {price1week}%
              30 DAY: {price30day}%
              1 YEAR: {price1year}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },
  '4h': {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price4h} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              7 DAY: {price1week}%
              30 DAY: {price30day}%
              1 YEAR: {price1year}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },

  daily: {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price24h} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              4H: {price4h}%
              7 DAY: {price1week}%
              30 DAY: {price30day}%
              1 YEAR: {price1year}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },

  weekly: {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price1week} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              4H: {price4h}%
              1 DAY: {price24h}%
              30 DAY: {price30day}%
              1 YEAR: {price1year}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },

  monthly: {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price30day} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              4H: {price4h}%
              1 DAY: {price24h}%
              7 DAY: {price1week}%
              1 YEAR: {price1year}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },

  yearly: {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price1year} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              4H: {price4h}%
              1 DAY: {price24h}%
              7 DAY: {price1week}%
              30 DAY: {price30day}%
              1 YEAR TO DATE: {price_year_to_date}% 
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },

  ytd: {
    text: ` [font-size: {fontSize}px font-weight: 400;]{qc_key}[/]
            [font-size: {fontSizeLev2}px; font-weight: 400;] {currency} {price_usd}
            {price_year_to_date} %[/]
          `,

    toolTip: `[bold]{symbol_name}[/]
              ---------------------
              1H: {price1h}%
              4H: {price4h}%
              1 DAY: {price24h}%
              7 DAY: {price1week}%
              30 DAY: {price30day}%
              1 YEAR: {price1year}%
              TREND MEAN: {qma_score}
              RSI 2H: {rsi2h}
              MC: {marketcapVal}`,
  },
}

export default function (USER_CAN_ACCESS_TREND_DATA: Ref<boolean> = ref(false)) {
  const { $f } = useContext()
  const { state } = useStore<State>()
  const { $axios } = useContext()
  const { currency } = useCurrency()
  const { exchanges } = useExchange()
  const { isActiveTab } = useBrowser()
  const {
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
  } = useHeatmapConfigs()

  const loading = ref(true)
  const rowData = ref<HeatmapRowData[]>([])
  const heatmapData = ref<HeatmapData[]>([])

  let apiPoolingTimeout: any = null

  const favorite = computed<string | null>(() =>
    displayFavorites.value ? state.global.screener.favoriteCoins.join(',') : null
  )
  const exchange = computed<string | null>(() => (exchanges.value.length ? exchanges.value.join(',') : null))
  const tileText = computed(() => tileConfigs[timeFrame.value].text)
  const tileTooltip = computed(() => tileConfigs[timeFrame.value].toolTip)

  const fieldNameByPerformancePeriod = computed<keyof HeatmapRowData>(() => {
    const frames: { [key in HeatmapIntervals]: keyof HeatmapRowData } = {
      '1h': 'price1h',
      '4h': 'price4h',
      daily: 'price24h',
      weekly: 'price1week',
      monthly: 'price30day',
      yearly: 'price1year',
      ytd: 'price_year_to_date',
    }
    const fieldByTime: keyof HeatmapRowData = frames[timeFrame.value] ?? 'price24h'
    return fieldByTime
  })

  const transformedData = computed<HeatmapRowData[]>(() => {
    return rowData.value.map((elem) => ({
      ...elem,
      qma_score: USER_CAN_ACCESS_TREND_DATA.value ? elem.qma_score : '--',
      color: setColor(<number>elem[fieldNameByPerformancePeriod.value], blueTile.value),
      price_usd: $f(<number>elem.price_usd, {
        autoDigits: true,
        currencySymbol: true,
        code: currency.value.code,
        locale: currency.value.locale,
      }),
      price1h: parseFloat($f(elem.price1h * 100, {})),
      price4h: parseFloat($f(elem.price4h * 100, {})),
      price24h: parseFloat($f(elem.price24h * 100, {})),
      price24hAbs: Math.abs(elem.price24h * 100),
      price1week: parseFloat($f(elem.price1week * 100, {})),
      price30day: parseFloat($f(elem.price30day * 100, {})),
      price1year: parseFloat($f(elem.price1year * 100, {})),
      price_year_to_date: parseFloat($f(elem.price_year_to_date * 100, {})),
      marketcap: Math.round(elem.marketcap / 10 ** 6),
      marketcapVal: $f(elem.marketcap, {
        useSymbol: true,
        currencySymbol: true,
        code: currency.value.code,
        locale: currency.value.locale,
      }),
    }))
  })

  const gainersAndLosers = computed<HeatmapData[]>(() => {
    if (displayGainersAndLosers.value) {
      const resp: HeatmapData[] = [
        { name: 'Gainers', children: [] },
        { name: 'Losers', children: [] },
      ]
      transformedData.value.forEach((elem) => {
        ;(elem[fieldNameByPerformancePeriod.value] ?? elem.price24h) > 0
          ? resp[0].children.push(elem)
          : resp[1].children.push(elem)
      })
      return resp
    } else {
      return [{ name: '', children: transformedData.value }]
    }
  })

  const updateData = computed<HeatmapUpdateData>(() =>
    transformedData.value.reduce((elem, item) => ({ ...elem, [item.qc_key]: item }), {})
  )

  async function fetchData() {
    try {
      const { data } = await $axios.get(APIEndpoints.heatmap, {
        params: {
          currency: currency.value.code,
          num_of_coins: numOfCoins.value,
          favorites: favorite.value,
          exchange: exchange.value,
        },
      })
      rowData.value = data.data ?? []
    } catch (error: any) {
      console.log(error.response ? error.responce : 'something went wrong')
      clearTimeout(apiPoolingTimeout)
    }
  }

  async function initLoad() {
    loading.value = true
    await fetchData()
    heatmapData.value = gainersAndLosers.value
    loading.value = false
  }

  onMounted(async () => {
    await initLoad()
    loading.value = false
    apiPoolingTimeout = setInterval(() => fetchData(), LongPooling.default)
  })

  // Watch tile Text and tooltip and update heatmap Data
  watch([tileText, tileTooltip, displayGainersAndLosers], () => {
    heatmapData.value = gainersAndLosers.value
  })

  watch([favorite, numOfCoins, currency, exchange], () => initLoad())

  watch(isActiveTab, (newValue: boolean) => {
    if (!newValue) {
      clearTimeout(apiPoolingTimeout)
      apiPoolingTimeout = null
    } else {
      apiPoolingTimeout = setInterval(() => fetchData(), LongPooling.default)
    }
  })

  onBeforeUnmount(() => {
    clearTimeout(apiPoolingTimeout)
    apiPoolingTimeout = null
  })

  return {
    fetchData,
    heatmapData,
    tileText,
    tileTooltip,
    timeFrameOptions,
    timeFrame,
    updateData,
    blockSize,
    blockSizeOptions,
    blockSizeName,
    blueTile,
    displayFavorites,
    numOfCoins,
    numOfCoinsOptions,
    displayGainersAndLosers,
    loading,
  }
}
