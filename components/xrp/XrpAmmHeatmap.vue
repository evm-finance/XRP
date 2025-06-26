<template>
  <v-card tile outlined :height="height">
    <v-card-title class="notranslate py-1 mt-1 text-no-wrap">
      <span class="subtitle-1">XRP AMM Heatmap</span>

      <v-spacer />
      <v-btn icon elevation="0" class="mx-2" height="20" width="20" @click="refreshData">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>

      <v-btn
        text
        color="primary"
        to="/xrp-amm-heatmap"
      >
        View Full
        <v-icon right>mdi-arrow-right</v-icon>
      </v-btn>
    </v-card-title>

    <v-card-text v-if="isLoading" class="d-flex justify-center align-center" style="height: 400px;">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </v-card-text>

    <v-card-text v-else-if="!isGraphQLEndpointAvailable" class="d-flex justify-center align-center" style="height: 400px;">
      <v-alert type="warning" text>
        <div class="text-center">
          <v-icon large color="warning" class="mb-2">mdi-alert-circle</v-icon>
          <div class="text-h6 mb-2">GraphQL Endpoint Not Configured</div>
          <div class="text-body-2">Please configure BASE_GRAPHQL_SERVER_URL environment variable</div>
        </div>
      </v-alert>
    </v-card-text>

    <v-card-text v-else-if="hasError" class="d-flex justify-center align-center" style="height: 400px;">
      <v-alert type="warning" text>
        <div class="text-center">
          <v-icon large color="warning" class="mb-2">mdi-alert-circle</v-icon>
          <div class="text-h6 mb-2">Using Demo Data</div>
          <div class="text-body-2">GraphQL server returned an error. Showing sample AMM data for demonstration.</div>
          <v-btn color="primary" class="mt-2" @click="refreshData">Retry Connection</v-btn>
        </div>
      </v-alert>
    </v-card-text>

    <xrp-heatmap-chart
      v-else
      data-nosnippet
      :block-size="blockSize"
      :chart-height="chartHeight"
      :data="heatmapData"
      :tile-body="'Pool: {poolId}<br/>Liquidity: ${totalLiquidityUsd}'"
      :tile-tooltip="'Pool: {poolId}<br/>Total Liquidity: ${totalLiquidityUsd}<br/>Asset 1: ${asset1ValueUsd}<br/>Asset 2: ${asset2ValueUsd}'"
    />
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'
import XrpHeatmapChart from '~/components/xrp/XrpHeatmapChart.vue'
import { useXrpHeatmap } from '~/composables/useXrpHeatmap'

export default defineComponent({
  components: {
    XrpHeatmapChart,
  },
  props: {
    height: { type: [String, Number], default: 548 },
    userCanAccessTrend: { type: Boolean, default: false },
  },
  setup(props) {
    // Use the AMM heatmap composable
    const {
      heatmapData,
      processedData,
      isLoading,
      hasError,
      refreshData,
      blockSize,
      isGraphQLEndpointAvailable,
    } = useXrpHeatmap()

    // COMPUTED
    const chartHeight = computed(() => Number(props.height) - 80) // Account for header

    return {
      isLoading,
      heatmapData,
      processedData,
      hasError,
      refreshData,
      blockSize,
      chartHeight,
      isGraphQLEndpointAvailable,
    }
  }
})
</script> 