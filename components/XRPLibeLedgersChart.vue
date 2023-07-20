<template>
  <div ref="chartDiv" :style="{ width: '100%', height: `${height}px` }"></div>
</template>

<script lang="ts">
import { defineComponent, onBeforeUnmount, onMounted, PropType, ref, useContext } from '@nuxtjs/composition-api'
import { Block } from '~/types/apollo/main/types'
import emitter from '~/types/emitter'

let am4core: any = null
let am4charts: any = null
let am4themesDark: any = null
if (process.browser) {
  am4core = require('@amcharts/amcharts4/core')
  am4charts = require('@amcharts/amcharts4/charts')
  am4themesDark = require('@amcharts/amcharts4/themes/dark')
}

type Props = {
  data: Block[]
  chartHeight: number
}

export default defineComponent<Props>({
  props: {
    data: { type: Array as PropType<Block[]>, default: () => [] },
    chartHeight: { type: Number, default: 220 },
  },

  setup(props) {
    // STATE
    const chartDiv = ref(null)
    const ledgers = ref<Block[]>(
      props.data
        .map((e) => ({
          ...e,
          date: e.minedAt * 1000,
        }))
        .sort((a, b) => (a.date < b.date ? -1 : 1))
    )
    let chart: any = null

    // COMPUTED
    const { env } = useContext()

    function renderChart() {
      am4core.useTheme(am4themesDark.default)

      am4core.addLicense(env.amChartLicense)

      chart = am4core.create(chartDiv.value, am4charts.XYChart)
      chart.hiddenState.properties.opacity = 0

      chart.padding(0, 0, 0, 0)

      chart.zoomOutButton.disabled = true

      const data = []
      let visits = 10
      let i = 0

      for (i = 0; i <= 25; i++) {
        visits -= Math.round((Math.random() < 0.5 ? 1 : -1) * Math.random() * 10)
        data.push({ date: new Date().setSeconds(i - 30), value: visits })
      }

      chart.data = ledgers.value

      const dateAxis = chart.xAxes.push(new am4charts.DateAxis())
      dateAxis.renderer.grid.template.location = 0
      dateAxis.renderer.minGridDistance = 30
      dateAxis.dateFormats.setKey('second', 'ss')
      dateAxis.periodChangeDateFormats.setKey('second', '[bold]h:mm a')
      dateAxis.periodChangeDateFormats.setKey('minute', '[bold]h:mm a')
      dateAxis.periodChangeDateFormats.setKey('hour', '[bold]h:mm a')
      dateAxis.renderer.inside = true
      dateAxis.renderer.axisFills.template.disabled = true
      dateAxis.renderer.ticks.template.disabled = true

      const valueAxis = chart.yAxes.push(new am4charts.ValueAxis())
      valueAxis.tooltip.disabled = true
      valueAxis.interpolationDuration = 500
      valueAxis.rangeChangeDuration = 500
      valueAxis.renderer.inside = true
      valueAxis.renderer.minLabelPosition = 0.05
      valueAxis.renderer.maxLabelPosition = 0.95
      valueAxis.renderer.axisFills.template.disabled = true
      valueAxis.renderer.ticks.template.disabled = true

      const series = chart.series.push(new am4charts.LineSeries())
      series.dataFields.dateX = 'date'
      series.dataFields.valueY = 'txCount'
      series.interpolationDuration = 500
      series.defaultState.transitionDuration = 0
      series.smoothing = 'monotoneX'
      series.tooltipText = '{dateX}: [b]{txCount}[/]'

      chart.events.on('datavalidated', () => {
        dateAxis.zoom({ start: 1 / 15, end: 1.2 }, false, true)
      })

      dateAxis.interpolationDuration = 500
      dateAxis.rangeChangeDuration = 500
      series.fillOpacity = 1
      const gradient = new am4core.LinearGradient()
      gradient.addColor(chart.colors.getIndex(0), 0.2)
      gradient.addColor(chart.colors.getIndex(0), 0)
      series.fill = gradient

      // this makes date axis labels to fade out
      dateAxis.renderer.labels.template.adapter.add('fillOpacity', (_: any, target: any) => {
        const dataItem = target.dataItem
        return dataItem.position
      })

      // need to set this, otherwise fillOpacity is not changed and not set
      dateAxis.events.on('validated', () => {
        am4core.iter.each(dateAxis.renderer.labels.iterator(), function (label: any) {
          // eslint-disable-next-line no-self-assign
          label.fillOpacity = label.fillOpacity
        })
      })

      // this makes date axis labels which are at equal minutes to be rotated
      dateAxis.renderer.labels.template.adapter.add('rotation', function (_: any, target: any) {
        const dataItem = target.dataItem
        if (
          dataItem.date &&
          dataItem.date.getTime() === am4core.time.round(new Date(dataItem.date.getTime()), 'minute').getTime()
        ) {
          target.verticalCenter = 'middle'
          target.horizontalCenter = 'left'
          return -90
        } else {
          target.verticalCenter = 'bottom'
          target.horizontalCenter = 'middle'
          return 0
        }
      })

      // bullet at the front of the line
      const bullet = series.createChild(am4charts.CircleBullet)
      bullet.circle.radius = 5
      bullet.fillOpacity = 1
      bullet.fill = chart.colors.getIndex(0)
      bullet.isMeasured = false

      series.events.on('validated', () => {
        bullet.moveTo(series.dataItems.last.point)
        bullet.validatePosition()
      })
    }

    onMounted(() => {
      renderChart()
    })

    emitter.on('onNewBlock', (data: Block[]) => {
      if (data.length) {
        const block: any = data[0]
        block.date = block.minedAt * 1000
        chart.addData(block, 1)
      }
    })

    onBeforeUnmount(() => {
      chart.dispose()
    })

    return {
      height: props.chartHeight,
      chartDiv,
    }
  },
})
</script>
