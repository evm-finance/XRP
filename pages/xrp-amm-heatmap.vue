<template>
  <div>
    <v-container fluid class="pa-0">
      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-title class="py-3">
              <div class="d-flex align-center">
                <v-icon large class="mr-3" color="primary">💧</v-icon>
                <div>
                  <h1 class="text-h4 font-weight-bold">XRP AMM Heatmap</h1>
                  <p class="text-subtitle-1 grey--text mb-0">
                    Visualize liquidity across XRP Automated Market Makers
                  </p>
                </div>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>

      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-text>
              <v-row>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold">{{ totalPools }}</div>
                    <div class="text-caption grey--text">Total Pools</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-success">{{ highLiquidityCount }}</div>
                    <div class="text-caption grey--text">High Liquidity</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold text-warning">{{ mediumLiquidityCount }}</div>
                    <div class="text-caption grey--text">Medium Liquidity</div>
                  </div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-center">
                    <div class="text-h6 font-weight-bold">${{ totalLiquidity }}</div>
                    <div class="text-caption grey--text">Total Liquidity</div>
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Controls -->
      <v-row no-gutters class="mb-4">
        <v-col>
          <v-card tile outlined>
            <v-card-text>
              <v-row align="center">
                <v-col cols="12" md="3">
                  <v-select
                    v-model="blockSize"
                    :items="blockSizeOptions"
                    label="Tile Size"
                    outlined
                    dense
                  />
                </v-col>
                <v-col cols="12" md="3">
                  <v-btn
                    color="primary"
                    @click="refreshData"
                    :loading="isLoading"
                    block
                  >
                    <v-icon left>mdi-refresh</v-icon>
                    Refresh
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Heatmap Chart -->
      <v-row no-gutters>
        <v-col>
          <v-card tile outlined class="mb-4">
            <v-card-text class="pa-0">
              <div v-if="isLoading" class="text-center py-8">
                <v-progress-circular indeterminate color="primary" size="64" />
                <div class="mt-4 text-h6">Loading AMM data...</div>
              </div>
              <div v-else-if="!isGraphQLEndpointAvailable" class="text-center py-8">
                <v-icon color="warning" size="64">mdi-alert-circle</v-icon>
                <div class="mt-4 text-h6">GraphQL Endpoint Not Configured</div>
                <div class="text-body-1 mt-2">Please configure BASE_GRAPHQL_SERVER_URL environment variable</div>
                <div class="text-caption mt-2">This is required to fetch real-time AMM data</div>
              </div>
              <div v-else-if="hasError" class="text-center py-8">
                <v-icon color="error" size="64">mdi-alert-circle</v-icon>
                <div class="mt-4 text-h6">Error loading AMM data</div>
                <div class="text-body-1 mt-2">Please check your network connection and try again</div>
                <v-btn color="primary" @click="refreshData" class="mt-4">
                  Try Again
                </v-btn>
              </div>
              <div v-else>
                <xrp-heatmap-chart
                  :data="heatmapData"
                  :tile-body="'Pool: {poolId}<br/>Liquidity: ${totalLiquidityUsd}'"
                  :tile-tooltip="'Pool: {poolId}<br/>Total Liquidity: ${totalLiquidityUsd}<br/>Asset 1: ${asset1ValueUsd}<br/>Asset 2: ${asset2ValueUsd}'"
                  :chart-height="heatmapHeight"
                  :block-size="blockSize"
                />
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Pools Table -->
      <v-row no-gutters class="mt-4">
        <v-col>
          <v-card tile outlined>
            <v-card-title class="py-3">
              <span class="text-h6">Top XRP AMM Pools by Liquidity</span>
              <v-spacer />
              <v-btn
                text
                color="primary"
                @click="refreshData"
                :loading="isLoading"
              >
                <v-icon left>mdi-refresh</v-icon>
                Refresh
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-data-table
                :headers="headers"
                :items="topPools"
                :loading="isLoading"
                :items-per-page="10"
                class="elevation-0"
                hide-default-footer
              >
                <template #item.poolId="{ item }">
                  <div class="d-flex align-center">
                    <div>
                      <div class="font-weight-medium">{{ item.poolId }}</div>
                      <div class="text-caption grey--text">Pool ID</div>
                    </div>
                  </div>
                </template>
                <template #item.asset1="{ item }">
                  <div>
                    <div class="font-weight-medium">{{ item.asset1.currency }}</div>
                    <div class="text-caption text-muted">{{ formatIssuerAddress(item.asset1.issuer) }}</div>
                  </div>
                </template>
                
                <template #item.asset2="{ item }">
                  <div>
                    <div class="font-weight-medium">{{ item.asset2.currency }}</div>
                    <div class="text-caption text-muted">{{ formatIssuerAddress(item.asset2.issuer) }}</div>
                  </div>
                </template>
                
                <template #item.totalLiquidityUsd="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.totalLiquidityUsd) }}</span>
                </template>
                
                <template #item.asset1ValueUsd="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.asset1ValueUsd) }}</span>
                </template>
                
                <template #item.asset2ValueUsd="{ item }">
                  <span class="font-weight-medium">${{ formatNumber(item.asset2ValueUsd) }}</span>
                </template>
                
                <template #item.fee="{ item }">
                  <span class="font-weight-medium">{{ formatPercentage(item.fee) }}%</span>
                </template>
                <template #item.actions="{ item }">
                  <v-btn
                    small
                    text
                    color="primary"
                    @click="viewPool(item)"
                  >
                    View
                  </v-btn>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref, useContext } from '@nuxtjs/composition-api'
import { useXrpHeatmap } from '~/composables/useXrpHeatmap'
import { useXrpConfigs } from '~/composables/useXrpConfigs'
import XrpHeatmapChart from '~/components/xrp/XrpHeatmapChart.vue'

export default defineComponent({
  components: {
    XrpHeatmapChart,
  },
  setup() {
    const { $f } = useContext()
    
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

    // Use XRP configs for block size options
    const { blockSizeOptions } = useXrpConfigs()

    // Heatmap height
    const heatmapHeight = ref(600)

    // Table headers
    const headers = [
      { text: 'Pool ID', value: 'poolId', align: 'left' },
      { text: 'Asset 1', value: 'asset1', align: 'left' },
      { text: 'Asset 2', value: 'asset2', align: 'left' },
      { text: 'Total Liquidity USD', value: 'totalLiquidityUsd', align: 'right' },
      { text: 'Asset 1 Value USD', value: 'asset1ValueUsd', align: 'right' },
      { text: 'Asset 2 Value USD', value: 'asset2ValueUsd', align: 'right' },
      { text: 'Fee', value: 'fee', align: 'right' },
    ]

    // Computed properties
    const topPools = computed(() => {
      return heatmapData.value
        .sort((a, b) => b.totalLiquidityUsd - a.totalLiquidityUsd)
        .slice(0, 20)
    })

    const totalPools = computed(() => heatmapData.value.length)
    
    const highLiquidityCount = computed(() => 
      heatmapData.value.filter(pool => pool.totalLiquidityUsd > 1000000).length
    )
    
    const mediumLiquidityCount = computed(() => 
      heatmapData.value.filter(pool => 
        pool.totalLiquidityUsd > 100000 && pool.totalLiquidityUsd <= 1000000
      ).length
    )
    
    const totalLiquidity = computed(() => 
      $f(heatmapData.value.reduce((sum, pool) => sum + pool.totalLiquidityUsd, 0), { 
        minDigits: 0, 
        after: '' 
      })
    )

    // Methods
    const formatNumber = (num: number) => {
      return $f(num, { minDigits: 0, after: '' })
    }

    const formatPercentage = (percentage: number) => {
      return $f(percentage, { minDigits: 2, after: '' })
    }

    const viewPool = (pool: any) => {
      // Navigate to pool page
      window.location.href = `/xrp-amm-pools/${pool.poolId}`
    }

    onMounted(() => {
      // Set responsive height
      if (process.client) {
        heatmapHeight.value = window.innerHeight * 0.6
      }
    })

    return {
      heatmapData,
      processedData,
      isLoading,
      hasError,
      refreshData,
      blockSize,
      blockSizeOptions,
      topPools,
      totalPools,
      highLiquidityCount,
      mediumLiquidityCount,
      totalLiquidity,
      headers,
      formatNumber,
      formatPercentage,
      viewPool,
      isGraphQLEndpointAvailable,
    }
  },
})
</script> 