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
import { XrpAmmHeatmapData } from '~/composables/useXrpHeatmap'
import type { XrpBlockSize } from '~/composables/useXrpConfigs'

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
  data: XrpAmmHeatmapData[]
  blockSize: XrpBlockSize
}

export default defineComponent<Props>({
  props: {
    tileBody: { type: String, default: '' },
    tileTooltip: { type: String, default: '' },
    chartHeight: { type: [String, Number] as PropType<string | number>, default: '600' },
    data: { type: Array as PropType<XrpAmmHeatmapData[]>, default: () => [] },
    blockSize: { type: String as PropType<XrpBlockSize>, default: () => 'totalLiquidityUsd' },
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
        const poolId = target.dataItem.dataContext.poolId
        
        // Set URL for AMM pool page
        target.url = `/xrp-amm-pools/${poolId}`

        let fontSize: any =
          (target.availableWidth / (poolId.length * 0.9)) * 0.75
        let fontSizeLev2: any =
          (target.availableWidth / (target.dataItem.dataContext.totalLiquidityUsd?.toString().length * 0.9 || 10)) * 0.5

        if (target.availableHeight < fontSize * 2) {
          fontSize = target.availableHeight / 2.5
          fontSizeLev2 = fontSize / 2.5
        }

        if (fontSizeLev2 > 18) {
          fontSizeLev2 = 18
        }

        // Get the value based on block size
        let displayValue = target.dataItem.dataContext.totalLiquidityUsd || 0
        if (tileSize.value === 'asset1ValueUsd') {
          displayValue = target.dataItem.dataContext.asset1ValueUsd || 0
        } else if (tileSize.value === 'asset2ValueUsd') {
          displayValue = target.dataItem.dataContext.asset2ValueUsd || 0
        }

        return tileTemplate
          .replace('{fontSize}', fontSize)
          .replace('{fontSizeLev2}', fontSizeLev2)
          .replace('{poolId}', poolId)
          .replace('{totalLiquidityUsd}', displayValue.toFixed(2))
          .replace('{asset1ValueUsd}', target.dataItem.dataContext.asset1ValueUsd?.toFixed(2) || '0')
          .replace('{asset2ValueUsd}', target.dataItem.dataContext.asset2ValueUsd?.toFixed(2) || '0')
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
      chart.value.dataFields.name = 'poolId'
      chart.value.dataFields.color = 'color'

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

      // Add click handler for navigation
      level0Column.value.events.on('hit', (ev: any) => {
        const url = ev.target.url
        if (url) {
          window.location.href = url
        }
      })

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

    watch(toolTipTemplate, (toolTipNewData) => (level0Column.value.tooltipText = toolTipNewData))

    onBeforeUnmount(() => chart.value.dispose())

    return { chartDiv }
  },
})
</script> 