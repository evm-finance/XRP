<template>
  <v-card tile outlined :height="height">
    <v-card-title class="notranslate py-1 mt-1 text-no-wrap">
      <span class="subtitle-1">XRP AMM Heatmap</span>

      <v-spacer />
      <v-btn icon elevation="0" class="mx-2" height="20" width="20" @click="displayFavorites = !displayFavorites">
        <v-icon :color="displayFavorites ? 'orange' : ''">{{ icons.mdiStar }} </v-icon>
      </v-btn>

      <heatmap-config-menu icon></heatmap-config-menu>
    </v-card-title>

    <v-card-text v-if="loading" class="d-flex justify-center align-center" style="height: 400px;">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </v-card-text>

    <heatmap-chart
      v-else
      data-nosnippet
      :block-size="blockSize"
      :chart-height="chartHeight"
      :data="heatmapData"
      :tile-body="tileText"
      :tile-tooltip="tileTooltip"
      :update-data="updateData"
    />
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'
import { mdiStar, mdiCog } from '@mdi/js'
import HeatmapChart from '~/components/HeatmapChart.vue'
import useXrpAmmHeatmap from '~/composables/useXrpAmmHeatmap'
import HeatmapConfigMenu from '~/components/HeatmapConfigMenu.vue'
import useHeatmapConfigs from '~/composables/useHeatmapConfigs'

export default defineComponent({
  components: {
    HeatmapConfigMenu,
    HeatmapChart,
  },
  props: {
    height: { type: [String, Number], default: 548 },
    userCanAccessTrend: { type: Boolean, default: false },
  },
  setup(props) {
    // COMPOSABLES
    const icons = { mdiStar, mdiCog }
    const USER_CAN_READ_TREND_DATA = computed(() => props.userCanAccessTrend)
    const { heatmapData, tileText, tileTooltip, updateData, loading } = useXrpAmmHeatmap(USER_CAN_READ_TREND_DATA)
    const { displayFavorites, blockSize } = useHeatmapConfigs()

    // Calculate chart height (subtract header height)
    const chartHeight = computed(() => {
      const heightNum = typeof props.height === 'string' ? parseInt(props.height) : props.height
      return heightNum - 60 // Subtract header height
    })

    return {
      icons,
      heatmapData,
      tileText,
      tileTooltip,
      updateData,
      loading,
      displayFavorites,
      blockSize,
      chartHeight,
    }
  },
})
</script> 