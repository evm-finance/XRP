<template>
  <v-card tile outlined :height="height">
    <v-card-title class="notranslate py-1 mt-1 text-no-wrap">
      <span class="subtitle-1">XRP Token Heatmap</span>

      <v-spacer />
      <v-btn icon elevation="0" class="mx-2" height="20" width="20" @click="displayFavorites = !displayFavorites">
        <v-icon :color="displayFavorites ? 'orange' : ''">{{ icons.mdiStar }} </v-icon>
      </v-btn>

      <heatmap-config-menu icon></heatmap-config-menu>
    </v-card-title>

    <v-card-text v-if="loading" class="d-flex justify-center align-center" style="height: 400px;">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </v-card-text>

    <xrp-heatmap-chart
      v-else
      data-nosnippet
      :block-size="blockSizeValue"
      :chart-height="chartHeight"
      :data="heatmapData"
      :tile-body="tileText"
      :tile-tooltip="tileTooltip"
      :update-data="updateData"
      heatmap-type="token"
    />
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from '@nuxtjs/composition-api'
import { mdiStar, mdiCog } from '@mdi/js'
import XrpHeatmapChart from '~/components/xrp/XrpHeatmapChart.vue'
import useXrpTokenHeatmap from '~/composables/useXrpTokenHeatmap'
import HeatmapConfigMenu from '~/components/heatmaps/HeatmapConfigMenu.vue'
import useHeatmapConfigs from '~/composables/useHeatmapConfigs'

export default defineComponent({
  components: {
    HeatmapConfigMenu,
    XrpHeatmapChart,
  },
  props: {
    height: { type: [String, Number], default: 548 },
    userCanAccessTrend: { type: Boolean, default: false },
  },
  setup(props) {
    // COMPOSABLES
    const icons = { mdiStar, mdiCog }
    const { loading, heatmapData, tileText, tileTooltip, updateData } = useXrpTokenHeatmap()
    const { displayFavorites, blueTile, blockSize, timeFrame } = useHeatmapConfigs()

    // COMPUTED
    const chartHeight = computed(() => Number(props.height))
    const blockSizeValue = computed(() => blockSize.value)

    return {
      icons,
      loading,
      heatmapData,
      tileText,
      tileTooltip,
      updateData,
      chartHeight,
      blockSizeValue,
      displayFavorites,
      blueTile,
      timeFrame,
    }
  }
})
</script> 