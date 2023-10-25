<template>
  <div ref="chartDiv" class="mt-4" :style="{ width: '100%', height: `${chartHeight}px` }"></div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  inject,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  useContext,
  useStore,
  watch,
} from '@nuxtjs/composition-api'
import { Balance } from '~/types/apollo/main/types'
import { State } from '~/types/state'
import { Web3, WEB3_PLUGIN_KEY } from '~/plugins/web3/web3'
let am4core: any = null
let am4charts: any = null
let am4themesAnimated: any = null
let am4themesDark: any = null
if (process.browser) {
  am4core = require('@amcharts/amcharts4/core')
  am4charts = require('@amcharts/amcharts4/charts')
  am4themesAnimated = require('@amcharts/amcharts4/themes/animated')
  am4themesDark = require('@amcharts/amcharts4/themes/dark')
}
type Props = {
  balances: Balance[]
}

export default defineComponent<Props>({
  props: {
    balances: { type: Array as PropType<Balance[]>, default: () => [] },
  },

  setup(props) {
    const { env } = useContext()
    const { state } = useStore<State>()
    const { getNetworkByChainNumber } = inject(WEB3_PLUGIN_KEY) as Web3
    const chartHeight = ref(300)
    const chartDiv = ref(null)

    let chart: any = null
    let pieChart: any = null
    let pieSeries: any = null

    const data = computed(() => {
      const chartData: { category: string; value: number; breakdown: Record<string, any> }[] = []

      props.balances.forEach((elem) => {
        const breakDown: { category: string; value: number }[] = []
        elem.items.forEach((bal) => {
          const val = bal.quote
          if (val > 1 / 10 ** 6) {
            breakDown.push({ category: `${bal.contractTickerSymbol.slice(0, 16)}`, value: val })
          }
        })
        const info = getNetworkByChainNumber(elem.chainId)
        if (info) {
          chartData.push({
            category: info.name,
            value: elem.items.reduce((n, { quote }) => n + quote, 0),
            breakdown: breakDown.sort((a, b) => (a.value < b.value ? 1 : -1)),
          })
        }
      })
      return chartData
    })

    onMounted(() => renderChart())
    onBeforeUnmount(() => (chart ? chart.dispose() : null))
    watch(
      () => data.value,
      () => {
        pieChart.data = data.value
      }
    )

    function renderChart() {
      if (state.ui.theme === 'dark') {
        am4core.useTheme(am4themesDark.default)
      } else {
        am4core.unuseTheme(am4themesDark.default)
        am4core.useTheme(am4themesAnimated.default)
      }
      am4core.addLicense(env.amChartLicense)

      // Create chart instance
      chart = am4core.create(chartDiv.value, am4core.Container)
      chart.width = am4core.percent(100)
      chart.height = am4core.percent(100)
      chart.layout = 'horizontal'

      /**
       * Pie chart
       */

      // Create chart instance
      pieChart = chart.createChild(am4charts.PieChart)
      pieChart.data = data.value
      pieChart.innerRadius = am4core.percent(65)
      pieChart.width = 600
      // Add and configure Series
      pieSeries = pieChart.series.push(new am4charts.PieSeries())
      pieSeries.dataFields.value = 'value'
      pieSeries.dataFields.category = 'category'
      pieSeries.slices.template.propertyFields.fill = 'color'
      pieSeries.labels.template.disabled = true

      // Set up labels
      const label1 = pieChart.seriesContainer.createChild(am4core.Label)
      label1.text = ''
      label1.horizontalCenter = 'middle'
      label1.fontSize = 35
      label1.fontWeight = 400
      label1.dy = -30

      const label2 = pieChart.seriesContainer.createChild(am4core.Label)
      label2.text = ''
      label2.horizontalCenter = 'middle'
      label2.fontSize = 12
      label2.dy = 20

      pieChart.legend = new am4charts.Legend()
      pieChart.legend.position = 'right'
      pieChart.legend.maxWidth = 300

      // Auto-select first slice on load
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      pieChart.events.on('ready', (ev: any) => {
        try {
          pieSeries.slices.getIndex(0).isActive = true
        } catch (e) {}
      })

      // Set up toggling events
      pieSeries.slices.template.events.on('toggled', (ev: any) => {
        if (ev.target.isActive) {
          // Untoggle other slices
          pieSeries.slices.each((slice: any) => {
            if (slice !== ev.target) {
              slice.isActive = false
            }
          })

          // Update column chart
          columnSeries.appeared = false
          columnChart.data = ev.target.dataItem.dataContext.breakdown
          columnSeries.fill = ev.target.fill
          columnSeries.reinit()

          // Update labels
          label1.text = pieChart.numberFormatter.format(ev.target.dataItem.values.value.percent, "#.'%'")
          label1.fill = ev.target.fill

          label2.text = ev.target.dataItem.category
        }
      })

      /**
       * Column chart
       */

      // Create chart instance
      const columnChart = chart.createChild(am4charts.XYChart)

      // Create axes
      const categoryAxis = columnChart.yAxes.push(new am4charts.CategoryAxis())
      categoryAxis.dataFields.category = 'category'
      categoryAxis.renderer.grid.template.location = 0
      categoryAxis.renderer.inversed = true
      categoryAxis.renderer.minGridDistance = 20
      categoryAxis.renderer.labels.template.fontSize = 13

      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const valueAxis = columnChart.xAxes.push(new am4charts.ValueAxis())

      // Create series
      const columnSeries = columnChart.series.push(new am4charts.ColumnSeries())
      columnSeries.dataFields.valueX = 'value'
      columnSeries.dataFields.categoryY = 'category'
      columnSeries.columns.template.strokeWidth = 0
      columnSeries.columns.template.tooltipText = '{categoryY}: [bold]' + '$' + '{valueX}[/]'

      // Add simple vertical scrollbar
      columnChart.scrollbarY = new am4core.Scrollbar()

      // Pre-zoom 8 first elements
      chart.events.on('ready', () => categoryAxis.zoomToIndexes(0, 8))
    }

    return { chartDiv, chartHeight }
  },
})
</script>
