<template>
  <div data-nosnippet>
    <div ref="chartDiv" :style="{ width: '100%', height: `${chartHeight}px` }"></div>
  </div>
</template>

<script lang="ts">
/* eslint-disable camelcase */

import {
  computed,
  defineComponent,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  toRefs,
  watch,
} from '@nuxtjs/composition-api'
import { HeatmapData, HeatmapUpdateData } from '~/types/heatmap'
import type { HeatmapTileSize } from '~/types/state'

let am4core: any = null
let am4charts: any = null
if (process.browser) {
  am4core = require('@amcharts/amcharts4/core')
  am4charts = require('@amcharts/amcharts4/charts')
}

type Props = {
  tileBody: string
  tileTooltip: string
  chartHeight: string
  data: HeatmapData[]
  updateData: HeatmapUpdateData
  blockSize: HeatmapTileSize
}

export default defineComponent<Props>({
  props: {
    tileBody: { type: String, default: '' },
    tileTooltip: { type: String, default: '' },
    chartHeight: { type: [String, Number] as PropType<string | number>, default: '600' },
    data: { type: Array as PropType<HeatmapData[]>, default: () => [] },
    updateData: { type: Object as PropType<HeatmapUpdateData>, default: () => ({}) },
    blockSize: { type: String as PropType<HeatmapTileSize>, default: () => 'marketcap_index' },
  },
  setup(props) {
    const chartDiv = ref(null)
    const chart: any = ref(null)
    const level0: any = ref(null)
    const level0Column: any = ref(null)
    const level0Bullet: any = ref(null)

    const initData = computed(() => props.data)
    const tileTemplate = toRefs(props).tileBody
    const toolTipTemplate = computed(() => props.tileTooltip)
    const tileSize = computed(() => props.blockSize)

    const tileRenderer = (_: any, target: any, tileTemplate: string) => {
      try {
        const key = target.dataItem.dataContext.dataContext.qc_key
        target.url = `/coins/${key}`

        let fontSize: any =
          (target.availableWidth / (target.dataItem.dataContext.dataContext.qc_key.length * 0.9)) * 0.75
        let fontSizeLev2: any =
          (target.availableWidth / (target.dataItem.dataContext.dataContext.price_usd.toString().length * 0.9)) * 0.5

        if (target.availableHeight < fontSize * 2) {
          fontSize = target.availableHeight / 2.5
          fontSizeLev2 = fontSize / 2.5
        }

        if (fontSizeLev2 > 18) {
          fontSizeLev2 = 18
        }

        return tileTemplate.replace('{fontSize}', fontSize).replace('{fontSizeLev2}', fontSizeLev2)
      } catch {}
    }

    function renderChart() {
      am4core.addLicense('CH187387301')
      chart.value = am4core.create(chartDiv.value, am4charts.TreeMap)
      chart.value.hiddenState.properties.opacity = 0
      chart.value.padding(0, 0, 0, 0)

      chart.value.data = initData.value
      chart.value.colors.step = 2

      /* Define data fields */
      chart.value.dataFields.value = tileSize.value
      chart.value.dataFields.name = 'symbol_name'
      chart.value.dataFields.color = 'color'
      chart.value.dataFields.children = 'children'

      // level 0 series template
      const level0SeriesTemplate = chart.value.seriesTemplates.create('0')
      const level0ColumnTemplate = level0SeriesTemplate.columns.template
      level0ColumnTemplate.column.cornerRadius(0, 0, 0, 0)
      level0ColumnTemplate.fillOpacity = 0
      level0ColumnTemplate.strokeWidth = 4
      level0ColumnTemplate.strokeOpacity = 0

      /* Configure top-level series */
      level0.value = chart.value.seriesTemplates.create('1')
      level0Column.value = level0.value.columns.template
      level0Column.value.stroke = am4core.color('#252525')
      level0Column.value.strokeWidth = 1
      level0Column.value.tooltipText = toolTipTemplate.value

      /* Add bullet labels */
      level0Bullet.value = level0.value.bullets.push(new am4charts.LabelBullet())
      level0Bullet.value.locationY = 0.5
      level0Bullet.value.locationX = 0.5
      level0Bullet.value.label.fill = am4core.color('#fff')
      level0Bullet.value.label.adapter.add('text', (_: any, target: any) => tileRenderer(_, target, tileTemplate.value))
    }

    onMounted(() => renderChart())

    watch(initData, (newData) => (chart.value.data = newData))

    watch(tileSize, (newTileSize) => {
      chart.value.dataFields.value = newTileSize
      chart.value.invalidateData()
    })

    const update = () => {}

    watch(toolTipTemplate, (toolTipNewData) => (level0Column.value.tooltipText = toolTipNewData))

    watch(
      () => props.updateData,
      (newData: HeatmapUpdateData) => {
        const chartData: HeatmapData[] = chart.value.data
        for (const child of chartData) {
          for (const element of child.children) {
            if (Object.hasOwnProperty.call(newData, element.qc_key)) {
              element.price_usd = newData[element.qc_key].price_usd

              element.price1h = newData[element.qc_key].price1h
              element.price4h = newData[element.qc_key].price4h
              element.price24h = newData[element.qc_key].price24h
              element.price24hAbs = newData[element.qc_key].price24hAbs
              element.price1week = newData[element.qc_key].price1week
              element.price30day = newData[element.qc_key].price30day
              element.price1year = newData[element.qc_key].price1year
              element.price_year_to_date = newData[element.qc_key].price_year_to_date

              element.qma_score = newData[element.qc_key].qma_score
              element.rsi2h = newData[element.qc_key].rsi2h
              element.marketcap = newData[element.qc_key].marketcap
              element.marketcapVal = newData[element.qc_key].marketcapVal
              element.marketcap_index = newData[element.qc_key].marketcap_index
              element.price24hAbs = newData[element.qc_key].price24hAbs
              element.volume24h = newData[element.qc_key].volume24h

              element.color = newData[element.qc_key].color
            }
          }
        }
        chart.value.invalidateRawData()
        try {
          for (let i = 0; i < chart.value.dataItems.length; i++) {
            const dataItem = chart.value.dataItems.getIndex(i)
            for (let c = 0; c < dataItem.children.length; c++) {
              const child = dataItem.children.getIndex(c)
              const qcKey = child.seriesDataItem.dataContext.dataContext.qc_key
              if (Object.hasOwnProperty.call(newData, qcKey)) {
                const color = newData[qcKey].color
                child.seriesDataItem.column.fill = am4core.color(color)
              }
            }
          }
        } catch {}
      }
    )
    onBeforeUnmount(() => chart.value.dispose())

    return { chartDiv, update }
  },
})
</script>
