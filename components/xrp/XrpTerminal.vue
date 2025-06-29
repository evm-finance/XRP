// DISABLED - XrpTerminal component temporarily disabled
// This component has complex dependencies that need to be resolved separately

<template>
  <div class="xrp-terminal">
    <!-- Terminal Header -->
    <v-card tile outlined class="mb-4">
      <v-card-title class="py-2">
        <span class="subtitle-1">XRP Terminal</span>
        <v-spacer />
        
        <!-- Quick Actions -->
        <v-btn
          icon
          small
          class="mx-1"
          @click="refreshAll"
          :loading="loading"
        >
          <v-icon>mdi-refresh</v-icon>
        </v-btn>
        
        <v-btn
          icon
          small
          class="mx-1"
          @click="toggleFullscreen"
        >
          <v-icon>{{ isFullscreen ? 'mdi-fullscreen-exit' : 'mdi-fullscreen' }}</v-icon>
        </v-btn>
        
        <v-btn
          icon
          small
          class="mx-1"
          @click="showSettings = !showSettings"
        >
          <v-icon>mdi-cog</v-icon>
        </v-btn>
      </v-card-title>
    </v-card>

    <!-- Settings Panel -->
    <v-expand-transition>
      <v-card v-if="showSettings" tile outlined class="mb-4">
        <v-card-text class="py-2">
          <v-row dense>
            <v-col cols="12" sm="6" md="3">
              <v-select
                v-model="displayMode"
                :items="displayModeOptions"
                label="Display Mode"
                dense
                outlined
              />
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-select
                v-model="refreshInterval"
                :items="refreshIntervalOptions"
                label="Refresh Interval"
                dense
                outlined
              />
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-switch
                v-model="autoRefresh"
                label="Auto Refresh"
                dense
              />
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-switch
                v-model="compactMode"
                label="Compact Mode"
                dense
              />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-expand-transition>

    <!-- Main Terminal Layout -->
    <v-layout wrap>
      <!-- Left Panel - Heatmaps -->
      <v-flex xs12 lg="8">
        <v-card tile outlined height="600" class="mb-4">
          <v-card-title class="py-2">
            <v-tabs v-model="activeHeatmapTab" background-color="transparent" color="primary">
              <v-tab value="token">
                <v-icon left>mdi-fire</v-icon>
                Token Heatmap
              </v-tab>
              <v-tab value="amm">
                <v-icon left>mdi-water</v-icon>
                AMM Heatmap
              </v-tab>
            </v-tabs>
            <v-spacer />
            
            <!-- Heatmap Controls -->
            <v-btn
              icon
              small
              class="mx-1"
              @click="displayFavorites = !displayFavorites"
            >
              <v-icon :color="displayFavorites ? 'orange' : ''">mdi-star</v-icon>
            </v-btn>
            
            <xrp-heatmap-config-menu icon />
          </v-card-title>

          <v-window v-model="activeHeatmapTab">
            <v-window-item value="token">
              <xrp-heatmap-chart
                :data="heatmapData"
                :update-data="refreshData"
                :tile-body="tileText"
                :tile-tooltip="tileTooltip"
                :block-size="blockSizeValue"
                :chart-height="chartHeight"
              />
            </v-window-item>
            <v-window-item value="amm">
              <xrp-amm-heatmap
                :height="chartHeight"
                :user-can-access-trend="true"
              />
            </v-window-item>
          </v-window>
        </v-card>
      </v-flex>

      <!-- Right Panel - Market Stats & Quick Actions -->
      <v-flex xs12 lg="4">
        <v-layout column>
          <!-- Market Stats -->
          <v-flex>
            <v-card tile outlined class="mb-4">
              <v-card-title class="py-2">
                <span class="subtitle-2">Market Stats</span>
              </v-card-title>
              <v-card-text class="py-2">
                <v-list dense>
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title class="text-caption">Total Market Cap</v-list-item-title>
                      <v-list-item-subtitle class="text-h6">
                        {{ formatMarketCap(totalMarketCap) }}
                      </v-list-item-subtitle>
                    </v-list-item-content>
                  </v-list-item>
                  
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title class="text-caption">24h Volume</v-list-item-title>
                      <v-list-item-subtitle class="text-h6">
                        {{ formatVolume(totalVolume) }}
                      </v-list-item-subtitle>
                    </v-list-item-content>
                  </v-list-item>
                  
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title class="text-caption">Active Tokens</v-list-item-title>
                      <v-list-item-subtitle class="text-h6">
                        {{ activeTokens }}
                      </v-list-item-subtitle>
                    </v-list-item-content>
                  </v-list-item>
                  
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title class="text-caption">AMM Pools</v-list-item-title>
                      <v-list-item-subtitle class="text-h6">
                        {{ ammPools }}
                      </v-list-item-subtitle>
                    </v-list-item-content>
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-flex>

          <!-- Quick Actions -->
          <v-flex>
            <v-card tile outlined class="mb-4">
              <v-card-title class="py-2">
                <span class="subtitle-2">Quick Actions</span>
              </v-card-title>
              <v-card-text class="py-2">
                <v-btn
                  block
                  color="primary"
                  class="mb-2"
                  @click="openTokenScreener"
                >
                  <v-icon left>mdi-view-list</v-icon>
                  Token Screener
                </v-btn>
                
                <v-btn
                  block
                  outlined
                  class="mb-2"
                  @click="openAmmExplorer"
                >
                  <v-icon left>mdi-swap-horizontal</v-icon>
                  AMM Explorer
                </v-btn>
                
                <v-btn
                  block
                  outlined
                  class="mb-2"
                  @click="openWalletManager"
                >
                  <v-icon left>mdi-wallet</v-icon>
                  Wallet Manager
                </v-btn>
                
                <v-btn
                  block
                  outlined
                  @click="openTransactionExplorer"
                >
                  <v-icon left>mdi-history</v-icon>
                  Transaction Explorer
                </v-btn>
              </v-card-text>
            </v-card>
          </v-flex>

          <!-- Recent Activity -->
          <v-flex>
            <v-card tile outlined>
              <v-card-title class="py-2">
                <span class="subtitle-2">Recent Activity</span>
              </v-card-title>
              <v-card-text class="py-2">
                <v-list dense>
                  <v-list-item
                    v-for="activity in recentActivity"
                    :key="activity.id"
                    class="px-0"
                  >
                    <v-list-item-content>
                      <v-list-item-title class="text-caption">
                        {{ activity.type }}
                      </v-list-item-title>
                      <v-list-item-subtitle class="text-body-2">
                        {{ activity.description }}
                      </v-list-item-subtitle>
                      <v-list-item-subtitle class="text-caption">
                        {{ formatTimestamp(activity.timestamp) }}
                      </v-list-item-subtitle>
                    </v-list-item-content>
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-flex>
    </v-layout>

    <!-- Bottom Panel - Token Screener -->
    <v-layout>
      <v-flex xs12>
        <v-card tile outlined class="mt-4">
          <v-card-title class="py-2">
            <span class="subtitle-2">Token Screener</span>
            <v-spacer />
            
            <!-- Screener Controls -->
            <v-text-field
              v-model="searchQuery"
              placeholder="Search tokens..."
              dense
              outlined
              hide-details
              class="mx-2"
              style="max-width: 200px;"
            >
              <template #prepend-inner>
                <v-icon>mdi-magnify</v-icon>
              </template>
            </v-text-field>
            
            <v-btn
              icon
              small
              class="mx-1"
              @click="showScreenerFilters = !showScreenerFilters"
            >
              <v-icon>mdi-filter</v-icon>
            </v-btn>
          </v-card-title>

          <!-- Screener Filters -->
          <v-expand-transition>
            <v-card-text v-if="showScreenerFilters" class="py-2">
              <v-row dense>
                <v-col cols="12" sm="6" md="3">
                  <v-text-field
                    v-model.number="filters.minMarketCap"
                    label="Min Market Cap"
                    type="number"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-text-field
                    v-model.number="filters.maxMarketCap"
                    label="Max Market Cap"
                    type="number"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-switch
                    v-model="filters.hasAmmPool"
                    label="Has AMM Pool"
                    dense
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-switch
                    v-model="filters.hasTrustLine"
                    label="Has Trust Line"
                    dense
                  />
                </v-col>
              </v-row>
            </v-card-text>
          </v-expand-transition>

          <!-- Token Grid -->
          <v-data-table
            :headers="screenerHeaders"
            :items="filteredTokens"
            :loading="loading"
            :items-per-page="screenerPageSize"
            :sort-by="screenerSortBy"
            :sort-desc="screenerSortOrder === 'desc'"
            class="elevation-0"
            dense
          >
            <template #item.name="{ item }">
              <div class="d-flex align-center">
                <v-avatar size="24" class="mr-2">
                  <v-img :src="getTokenIcon(item.currency)" />
                </v-avatar>
                <div>
                  <div class="font-weight-medium">{{ item.name }}</div>
                  <div class="text-caption">{{ item.currency }}</div>
                </div>
              </div>
            </template>

            <template #item.price="{ item }">
              <div>
                <div class="font-weight-medium">{{ formatXrpPrice(item.price_usd) }}</div>
                <div class="text-caption" :class="item.price_change_24h >= 0 ? 'success--text' : 'error--text'">
                  {{ formatPercentageChange(item.price_change_24h) }}
                </div>
              </div>
            </template>

            <template #item.marketcap="{ item }">
              {{ formatMarketCap(item.marketcap) }}
            </template>

            <template #item.volume="{ item }">
              {{ formatVolume(item.volume_24h) }}
            </template>

            <template #item.issuer="{ item }">
              <div class="d-flex align-center">
                <span class="text-caption font-family-monospace">{{ formatIssuerAddress(item.issuer) }}</span>
                <v-btn
                  icon
                  x-small
                  class="ml-1"
                  @click="copyToClipboard(item.issuer)"
                >
                  <v-icon x-small>mdi-content-copy</v-icon>
                </v-btn>
              </div>
            </template>

            <template #item.ammPool="{ item }">
              <v-chip
                :color="item.hasAmmPool ? 'success' : 'grey'"
                small
                outlined
              >
                {{ item.hasAmmPool ? 'Yes' : 'No' }}
              </v-chip>
            </template>

            <template #item.actions="{ item }">
              <v-btn
                icon
                x-small
                @click="viewTokenDetails(item)"
              >
                <v-icon x-small>mdi-open-in-new</v-icon>
              </v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-flex>
    </v-layout>

    <!-- AMM Pools Section -->
    <v-layout>
      <v-flex xs12>
        <v-card tile outlined class="mt-4">
          <v-card-title class="py-2">
            <span class="subtitle-2">AMM Pools</span>
            <v-spacer />
            
            <!-- AMM Controls -->
            <v-text-field
              v-model="ammSearchQuery"
              placeholder="Search pools..."
              dense
              outlined
              hide-details
              class="mx-2"
              style="max-width: 200px;"
            >
              <template #prepend-inner>
                <v-icon>mdi-magnify</v-icon>
              </template>
            </v-text-field>
            
            <v-btn
              icon
              small
              class="mx-1"
              @click="showAmmFilters = !showAmmFilters"
            >
              <v-icon>mdi-filter</v-icon>
            </v-btn>
          </v-card-title>

          <!-- AMM Filters -->
          <v-expand-transition>
            <v-card-text v-if="showAmmFilters" class="py-2">
              <v-row dense>
                <v-col cols="12" sm="6" md="3">
                  <v-text-field
                    v-model.number="ammFilters.minLiquidity"
                    label="Min Liquidity"
                    type="number"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-text-field
                    v-model.number="ammFilters.minVolume"
                    label="Min Volume"
                    type="number"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-text-field
                    v-model.number="ammFilters.minApr"
                    label="Min APR"
                    type="number"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <v-select
                    v-model="ammFilters.feeTier"
                    :items="feeTierOptions"
                    label="Fee Tier"
                    dense
                    outlined
                    hide-details
                  />
                </v-col>
              </v-row>
            </v-card-text>
          </v-expand-transition>

          <!-- AMM Pools Grid -->
          <v-data-table
            :headers="ammHeaders"
            :items="filteredAmmPools"
            :loading="loading"
            :items-per-page="ammPageSize"
            :sort-by="ammSortBy"
            :sort-desc="ammSortOrder === 'desc'"
            class="elevation-0"
            dense
          >
            <template #item.pair="{ item }">
              <div class="d-flex align-center">
                <div class="d-flex align-center mr-2">
                  <v-avatar size="24" class="mr-1">
                    <v-img :src="getTokenIcon(item.token1.symbol)" />
                  </v-avatar>
                  <v-avatar size="24" class="mr-1">
                    <v-img :src="getTokenIcon(item.token2.symbol)" />
                  </v-avatar>
                </div>
                <div>
                  <div class="font-weight-medium">{{ item.token1.symbol }}/{{ item.token2.symbol }}</div>
                  <div class="text-caption">{{ item.token1.name }}/{{ item.token2.name }}</div>
                </div>
              </div>
            </template>

            <template #item.liquidity="{ item }">
              {{ formatMarketCap(item.liquidity) }}
            </template>

            <template #item.volume="{ item }">
              {{ formatVolume(item.volume24h) }}
            </template>

            <template #item.fee="{ item }">
              {{ (item.fee * 100).toFixed(3) }}%
            </template>

            <template #item.apr="{ item }">
              <div>
                <div class="font-weight-medium">{{ item.apr.toFixed(2) }}%</div>
                <div class="text-caption" :class="item.priceChange24h >= 0 ? 'success--text' : 'error--text'">
                  {{ formatPercentageChange(item.priceChange24h) }}
                </div>
              </div>
            </template>

            <template #item.actions="{ item }">
              <v-btn
                icon
                x-small
                class="mr-1"
                @click="viewPoolDetails(item)"
              >
                <v-icon x-small>mdi-open-in-new</v-icon>
              </v-btn>
              <v-btn
                icon
                x-small
                @click="openPoolActions(item)"
              >
                <v-icon x-small>mdi-swap-horizontal</v-icon>
              </v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-flex>
    </v-layout>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted } from '@nuxtjs/composition-api'
import { useXrpHeatmap } from '~/composables/useXrpHeatmap'
import { useXrpConfigs } from '~/composables/useXrpConfigs'
import { useXrpFormatters } from '~/composables/useXrpFormatters'
import XrpHeatmapChart from '~/components/xrp/XrpHeatmapChart.vue'
import XrpHeatmapConfigMenu from '~/components/xrp/XrpHeatmapConfigMenu.vue'
import XrpAmmHeatmap from '~/components/xrp/XrpAmmHeatmap.vue'

export default defineComponent({
  name: 'XrpTerminal',
  components: {
    XrpHeatmapChart,
    XrpHeatmapConfigMenu,
    XrpAmmHeatmap,
  },
  setup() {
    // COMPOSABLES
    const { heatmapData, isLoading, hasError, refreshData } = useXrpHeatmap()
    const { timeFrame, blockSize, displayFavorites, numOfTokens, displayGainersAndLosers, blueTile } = useXrpConfigs()
    const { formatXrpPrice, formatMarketCap, formatVolume, formatIssuerAddress, formatPercentageChange, formatTimestamp } = useXrpFormatters()

    // REACTIVE DATA
    const chartHeight = ref(540)
    const blockSizeValue = computed(() => blockSize.value)
    const activeHeatmapTab = ref('token')

    // METHODS
    const refreshAll = () => {
      // Implement refreshAll logic
    }

    const toggleFullscreen = () => {
      // Implement toggleFullscreen logic
    }

    // LIFECYCLE
    onMounted(() => {
      // Implement onMounted logic
    })

    return {
      // Data
      heatmapData,
      isLoading,
      hasError,
      refreshData,
      
      // Configs
      timeFrame,
      blockSize,
      displayFavorites,
      numOfTokens,
      displayGainersAndLosers,
      blueTile,
      
      // Methods
      refreshAll,
      toggleFullscreen,
      
      // Reactive Data
      activeHeatmapTab,
      chartHeight,
      blockSizeValue,
    }
  }
})
</script>

<style scoped>
.xrp-terminal {
  padding: 16px;
}

.font-family-monospace {
  font-family: 'Courier New', monospace;
}
</style> 