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
      <!-- Left Panel - Heatmap -->
      <v-flex xs12 lg="8">
        <v-card tile outlined height="600" class="mb-4">
          <v-card-title class="py-2">
            <span class="subtitle-2">Token Heatmap</span>
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

          <xrp-heatmap-chart
            :data="heatmapData"
            :update-data="updateData"
            :tile-body="tileText"
            :tile-tooltip="tileTooltip"
            :block-size="blockSize"
            :chart-height="540"
          />
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
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, onBeforeUnmount } from '@nuxtjs/composition-api'
import { useXrpHeatmap } from '~/composables/useXrpHeatmap'
import { useXrpConfigs } from '~/composables/useXrpConfigs'
import { useXrpFormatters } from '~/composables/useXrpFormatters'
import { useXrpGridRenderers } from '~/composables/useXrpGridRenderers'
import XrpHeatmapChart from '~/components/xrp/XrpHeatmapChart.vue'
import XrpHeatmapConfigMenu from '~/components/xrp/XrpHeatmapConfigMenu.vue'

export default defineComponent({
  name: 'XrpTerminal',
  components: {
    XrpHeatmapChart,
    XrpHeatmapConfigMenu,
  },
  setup() {
    // Composables
    const { heatmapData, updateData, loading, tileText, tileTooltip, fetchHeatmapData } = useXrpHeatmap()
    const {
      displayMode,
      autoRefresh,
      refreshInterval,
      compactMode,
      displayFavorites,
      blockSize,
      screenerSortBy,
      screenerSortOrder,
      screenerPageSize,
    } = useXrpConfigs()
    const { formatXrpPrice, formatMarketCap, formatVolume, formatIssuerAddress, formatPercentageChange, formatTimestamp } = useXrpFormatters()

    // Local state
    const showSettings = ref(false)
    const showScreenerFilters = ref(false)
    const isFullscreen = ref(false)
    const searchQuery = ref('')
    const filters = ref({
      minMarketCap: 0,
      maxMarketCap: Infinity,
      hasAmmPool: false,
      hasTrustLine: false,
    })

    // Market stats
    const totalMarketCap = computed(() => 
      heatmapData.value.reduce((sum, token) => sum + token.marketcap, 0)
    )
    const totalVolume = computed(() => 
      heatmapData.value.reduce((sum, token) => sum + token.volume_24h, 0)
    )
    const activeTokens = computed(() => heatmapData.value.length)
    const ammPools = computed(() => 
      heatmapData.value.filter(token => token.hasAmmPool).length
    )

    // Display options
    const displayModeOptions = [
      { text: 'USD', value: 'usd' },
      { text: 'XRP', value: 'xrp' },
    ]

    const refreshIntervalOptions = [
      { text: '5 seconds', value: 5000 },
      { text: '10 seconds', value: 10000 },
      { text: '30 seconds', value: 30000 },
      { text: '1 minute', value: 60000 },
    ]

    // Screener headers
    const screenerHeaders = [
      { text: 'Token', value: 'name', sortable: false },
      { text: 'Price', value: 'price_usd', sortable: true },
      { text: 'Market Cap', value: 'marketcap', sortable: true },
      { text: 'Volume 24h', value: 'volume_24h', sortable: true },
      { text: 'Issuer', value: 'issuer', sortable: false },
      { text: 'AMM Pool', value: 'ammPool', sortable: false },
      { text: 'Actions', value: 'actions', sortable: false },
    ]

    // Filtered tokens for screener
    const filteredTokens = computed(() => {
      let tokens = heatmapData.value

      // Apply search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        tokens = tokens.filter(token => 
          token.name.toLowerCase().includes(query) ||
          token.currency.toLowerCase().includes(query)
        )
      }

      // Apply market cap filters
      if (filters.value.minMarketCap > 0) {
        tokens = tokens.filter(token => token.marketcap >= filters.value.minMarketCap)
      }
      if (filters.value.maxMarketCap < Infinity) {
        tokens = tokens.filter(token => token.marketcap <= filters.value.maxMarketCap)
      }

      // Apply AMM pool filter
      if (filters.value.hasAmmPool) {
        tokens = tokens.filter(token => token.hasAmmPool)
      }

      // Apply trust line filter
      if (filters.value.hasTrustLine) {
        tokens = tokens.filter(token => token.hasTrustLine)
      }

      return tokens
    })

    // Recent activity (mock data)
    const recentActivity = ref([
      {
        id: 1,
        type: 'Token Listed',
        description: 'SOLO token added to AMM pool',
        timestamp: new Date(Date.now() - 5 * 60 * 1000),
      },
      {
        id: 2,
        type: 'Price Update',
        description: 'XRP price increased by 2.5%',
        timestamp: new Date(Date.now() - 10 * 60 * 1000),
      },
      {
        id: 3,
        type: 'New AMM Pool',
        description: 'USDC/XRP pool created',
        timestamp: new Date(Date.now() - 15 * 60 * 1000),
      },
    ])

    // Methods
    const refreshAll = async () => {
      await fetchHeatmapData()
    }

    const toggleFullscreen = () => {
      isFullscreen.value = !isFullscreen.value
      if (isFullscreen.value) {
        document.documentElement.requestFullscreen()
      } else {
        document.exitFullscreen()
      }
    }

    const openTokenScreener = () => {
      window.$nuxt.$router.push('/xrp-token-screener')
    }

    const openAmmExplorer = () => {
      window.$nuxt.$router.push('/xrp-amm-explorer')
    }

    const openWalletManager = () => {
      window.$nuxt.$router.push('/xrp-wallet-manager')
    }

    const openTransactionExplorer = () => {
      window.$nuxt.$router.push('/xrp-explorer')
    }

    const getTokenIcon = (currency: string) => {
      return `/img/tokens/${currency.toLowerCase()}.svg`
    }

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        window.$nuxt.$emit('show-toast', { message: 'Copied to clipboard', type: 'success' })
      } catch (err) {
        console.error('Failed to copy to clipboard:', err)
      }
    }

    const viewTokenDetails = (token: any) => {
      window.$nuxt.$router.push(`/token/${token.currency}`)
    }

    // Lifecycle
    onMounted(() => {
      fetchHeatmapData()
    })

    return {
      // Data
      heatmapData,
      updateData,
      loading,
      tileText,
      tileTooltip,
      showSettings,
      showScreenerFilters,
      isFullscreen,
      searchQuery,
      filters,
      recentActivity,

      // Computed
      totalMarketCap,
      totalVolume,
      activeTokens,
      ammPools,
      displayModeOptions,
      refreshIntervalOptions,
      screenerHeaders,
      filteredTokens,

      // Configs
      displayMode,
      autoRefresh,
      refreshInterval,
      compactMode,
      displayFavorites,
      blockSize,
      screenerSortBy,
      screenerSortOrder,
      screenerPageSize,

      // Formatters
      formatXrpPrice,
      formatMarketCap,
      formatVolume,
      formatIssuerAddress,
      formatPercentageChange,
      formatTimestamp,

      // Methods
      refreshAll,
      toggleFullscreen,
      openTokenScreener,
      openAmmExplorer,
      openWalletManager,
      openTransactionExplorer,
      getTokenIcon,
      copyToClipboard,
      viewTokenDetails,
    }
  },
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