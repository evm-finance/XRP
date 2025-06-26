<template>
  <v-card tile outlined height="548">
    <v-card-title class="notranslate py-1 mt-1 text-no-wrap">
      <span class="subtitle-1">Heatmap</span>

      <v-spacer />
      <v-btn icon elevation="0" class="mx-2" height="20" width="20" @click="displayFavorites = !displayFavorites">
        <v-icon :color="displayFavorites ? 'orange' : ''">{{ icons.mdiStar }} </v-icon>
      </v-btn>

      <heatmap-config-menu icon></heatmap-config-menu>
    </v-card-title>

    <heatmap-chart
      data-nosnippet
      :block-size="blockSize"
      :chart-height="506"
      :data="heatmapData"
      :tile-body="tileText"
      :tile-tooltip="tileTooltip"
      :update-data="updateData"
    />
  </v-card>
</template>

<script>
import { computed, defineComponent } from '@nuxtjs/composition-api'
import { mdiCheckboxBlankOutline, mdiCheckboxMarked, mdiChevronDown, mdiStar, mdiCog } from '@mdi/js'
import HeatmapChart from '~/components/heatmaps/HeatmapChart'
import useHeatmap from '~/composables/heatmap/useHeatmap'
import HeatmapConfigMenu from '~/components/heatmaps/HeatmapConfigMenu'
import useHeatmapConfigs from '~/composables/heatmap/useHeatmapConfigs'
export default defineComponent({
  components: {
    HeatmapConfigMenu,
    HeatmapChart,
  },
  props: {
    userCanAccessTrend: { type: Boolean, default: false },
  },
  setup(props) {
    // COMPOSABLES
    const icons = { mdiChevronDown, mdiCheckboxBlankOutline, mdiCheckboxMarked, mdiStar, mdiCog }
    const USER_CAN_READ_TREND_DATA = computed(() => props.userCanAccessTrend)
    const { heatmapData, tileText, tileTooltip, updateData } = useHeatmap(USER_CAN_READ_TREND_DATA)
    const { displayFavorites, blockSize } = useHeatmapConfigs()

    return {
      icons,
      heatmapData,
      tileText,
      tileTooltip,
      updateData,
      displayFavorites,
      blockSize,
    }
  },
})
</script>
