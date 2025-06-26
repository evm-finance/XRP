<template>
  <div>
    <v-row no-gutters justify="center">
      <v-col cols="12" md="11">
        <v-row justify="center">
          <v-col cols="12">
            <h1 class="text-h4">XRP Token Mints & Liquidity Pools</h1>
            <p class="text-body-2 grey--text mt-2">
              Discover recently minted XRP tokens and explore liquidity pools with detailed analytics
            </p>
          </v-col>
        </v-row>

        <!-- Stats Cards -->
        <v-row justify="center" class="mb-6">
          <v-col cols="12" sm="6" md="3">
            <v-card tile outlined class="text-center pa-4">
              <div class="text-h6 font-weight-bold primary--text">{{ formatAmount(filteredTokenMints.length) }}</div>
              <div class="text-caption grey--text">New Tokens</div>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card tile outlined class="text-center pa-4">
              <div class="text-h6 font-weight-bold green--text">{{ formatCurrency(totalMarketCap) }}</div>
              <div class="text-caption grey--text">Total Market Cap</div>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card tile outlined class="text-center pa-4">
              <div class="text-h6 font-weight-bold blue--text">{{ formatCurrency(totalVolume24H) }}</div>
              <div class="text-caption grey--text">24H Volume</div>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card tile outlined class="text-center pa-4">
              <div class="text-h6 font-weight-bold orange--text">{{ formatCurrency(totalTvl) }}</div>
              <div class="text-caption grey--text">Total TVL</div>
            </v-card>
          </v-col>
        </v-row>

        <!-- Token Mints Section -->
        <v-row justify="center" class="mb-8">
          <v-col cols="12">
            <v-card tile outlined>
              <v-card-title class="d-flex align-center justify-space-between">
                <span>Recent Token Mints</span>
                <div class="d-flex align-center">
                  <!-- Time Range Filter -->
                  <v-btn-toggle
                    v-model="selectedTimeRange"
                    mandatory
                    dense
                    class="mr-4"
                  >
                    <v-btn value="24h" small>24H</v-btn>
                    <v-btn value="7d" small>7D</v-btn>
                    <v-btn value="30d" small>30D</v-btn>
                    <v-btn value="all" small>ALL</v-btn>
                  </v-btn-toggle>
                  
                  <!-- Refresh Button -->
                  <v-btn
                    icon
                    @click="refreshData"
                    :loading="loading"
                  >
                    <v-icon>mdi-refresh</v-icon>
                  </v-btn>
                </div>
              </v-card-title>
              
              <!-- Filters -->
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="6" md="3">
                    <v-text-field
                      v-model.number="tokenMintFilters.minMarketcap"
                      label="Min Market Cap ($)"
                      type="number"
                      outlined
                      dense
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-text-field
                      v-model.number="tokenMintFilters.minVolume"
                      label="Min 24H Volume ($)"
                      type="number"
                      outlined
                      dense
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-checkbox
                      v-model="tokenMintFilters.hasLiquidity"
                      label="Has Liquidity"
                      dense
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-btn
                      color="primary"
                      outlined
                      @click="clearTokenMintFilters"
                    >
                      Clear Filters
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-text>
                <v-data-table
                  v-if="!loading"
                  :headers="tokenMintHeaders"
                  :items="filteredTokenMints"
                  :items-per-page="10"
                  class="elevation-0"
                  mobile-breakpoint="0"
                >
                  <template #[`item.token`]="{ item }">
                    <div class="d-flex align-center">
                      <v-avatar size="32" class="mr-2">
                        <v-img
                          :src="getImageUrl(item.currency)"
                          :lazy-src="getImageUrl(item.currency)"
                        />
                      </v-avatar>
                      <div>
                        <div class="font-weight-medium">{{ item.tokenName }}</div>
                        <div class="text-caption grey--text">{{ item.currency }}</div>
                      </div>
                    </div>
                  </template>

                  <template #[`item.issuer`]="{ item }">
                    <div class="d-flex align-center">
                      <span class="font-family-mono">{{ formatAddress(item.issuerAddress) }}</span>
                      <v-btn
                        icon
                        x-small
                        class="ml-1"
                        @click="copyToClipboard(item.issuerAddress)"
                      >
                        <v-icon size="16">mdi-content-copy</v-icon>
                      </v-btn>
                    </div>
                  </template>

                  <template #[`item.mintDate`]="{ item }">
                    <span class="text-no-wrap">{{ formatDate(item.mintDate) }}</span>
                  </template>

                  <template #[`item.initialSupply`]="{ item }">
                    <span>{{ formatAmount(item.initialSupply) }}</span>
                  </template>

                  <template #[`item.currentSupply`]="{ item }">
                    <span>{{ formatAmount(item.currentSupply) }}</span>
                  </template>

                  <template #[`item.price`]="{ item }">
                    <span class="font-weight-medium">{{ formatCurrency(item.price) }}</span>
                  </template>

                  <template #[`item.marketcap`]="{ item }">
                    <span class="font-weight-medium">{{ formatCurrency(item.marketcap) }}</span>
                  </template>

                  <template #[`item.volume24H`]="{ item }">
                    <span>{{ formatCurrency(item.volume24H) }}</span>
                  </template>

                  <template #[`item.holders`]="{ item }">
                    <span>{{ formatAmount(item.holders) }}</span>
                  </template>

                  <template #[`item.liquidityPools`]="{ item }">
                    <div class="d-flex align-center">
                      <v-chip
                        v-if="item.liquidityPools.length > 0"
                        small
                        color="success"
                        outlined
                      >
                        {{ item.liquidityPools.length }} pools
                      </v-chip>
                      <span v-else class="grey--text text-caption">No pools</span>
                    </div>
                  </template>
                </v-data-table>

                <div v-else class="text-center pa-4">
                  <v-progress-circular indeterminate color="primary" />
                  <div class="mt-2">Loading token mints...</div>
                </div>

                <div v-if="!loading && filteredTokenMints.length === 0" class="text-center pa-4">
                  <v-icon size="48" class="mb-2 grey--text">mdi-coin-off</v-icon>
                  <div class="text-h6 grey--text">No token mints found</div>
                  <div class="text-body-2 grey--text">Try adjusting your filters</div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Liquidity Pools Section -->
        <v-row justify="center">
          <v-col cols="12">
            <v-card tile outlined>
              <v-card-title class="d-flex align-center justify-space-between">
                <span>Liquidity Pools</span>
                <div class="d-flex align-center">
                  <!-- Sort Options -->
                  <v-select
                    v-model="selectedSortBy"
                    :items="sortOptions"
                    item-text="label"
                    item-value="value"
                    outlined
                    dense
                    class="mr-2"
                    style="min-width: 120px"
                  />
                  <v-btn-toggle
                    v-model="selectedSortOrder"
                    mandatory
                    dense
                    class="mr-4"
                  >
                    <v-btn value="desc" small>
                      <v-icon small>mdi-sort-descending</v-icon>
                    </v-btn>
                    <v-btn value="asc" small>
                      <v-icon small>mdi-sort-ascending</v-icon>
                    </v-btn>
                  </v-btn-toggle>
                </div>
              </v-card-title>

              <!-- Pool Filters -->
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="6" md="3">
                    <v-text-field
                      v-model.number="liquidityPoolFilters.minTvl"
                      label="Min TVL ($)"
                      type="number"
                      outlined
                      dense
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-text-field
                      v-model.number="liquidityPoolFilters.minVolume"
                      label="Min 24H Volume ($)"
                      type="number"
                      outlined
                      dense
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-text-field
                      v-model.number="liquidityPoolFilters.minApr"
                      label="Min APR (%)"
                      type="number"
                      outlined
                      dense
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-btn
                      color="primary"
                      outlined
                      @click="clearLiquidityPoolFilters"
                    >
                      Clear Filters
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-text>
                <v-data-table
                  v-if="!loading"
                  :headers="liquidityPoolHeaders"
                  :items="filteredLiquidityPools"
                  :items-per-page="10"
                  class="elevation-0"
                  mobile-breakpoint="0"
                >
                  <template #[`item.pool`]="{ item }">
                    <div class="d-flex align-center">
                      <div class="d-flex align-center mr-2">
                        <v-avatar size="24" class="mr-1">
                          <v-img
                            :src="getImageUrl(item.asset1.symbol)"
                            :lazy-src="getImageUrl(item.asset1.symbol)"
                          />
                        </v-avatar>
                        <v-avatar size="24">
                          <v-img
                            :src="getImageUrl(item.asset2.symbol)"
                            :lazy-src="getImageUrl(item.asset2.symbol)"
                          />
                        </v-avatar>
                      </div>
                      <div>
                        <div class="font-weight-medium">{{ item.asset1.symbol }}/{{ item.asset2.symbol }}</div>
                        <div class="text-caption grey--text">{{ item.asset1.name }}/{{ item.asset2.name }}</div>
                      </div>
                    </div>
                  </template>

                  <template #[`item.tvl`]="{ item }">
                    <span class="font-weight-medium">{{ formatCurrency(item.tvl) }}</span>
                  </template>

                  <template #[`item.volume24H`]="{ item }">
                    <span>{{ formatCurrency(item.volume24H) }}</span>
                  </template>

                  <template #[`item.fees24H`]="{ item }">
                    <span>{{ formatCurrency(item.fees24H) }}</span>
                  </template>

                  <template #[`item.apr`]="{ item }">
                    <span class="font-weight-medium green--text">{{ item.apr.toFixed(2) }}%</span>
                  </template>

                  <template #[`item.priceChange24H`]="{ item }">
                    <span :class="item.priceChange24H >= 0 ? 'green--text' : 'red--text'">
                      {{ formatPercentage(item.priceChange24H) }}
                    </span>
                  </template>

                  <template #[`item.priceChange7D`]="{ item }">
                    <span :class="item.priceChange7D >= 0 ? 'green--text' : 'red--text'">
                      {{ formatPercentage(item.priceChange7D) }}
                    </span>
                  </template>

                  <template #[`item.transactions24H`]="{ item }">
                    <span>{{ formatAmount(item.transactions24H) }}</span>
                  </template>

                  <template #[`item.uniqueTraders24H`]="{ item }">
                    <span>{{ formatAmount(item.uniqueTraders24H) }}</span>
                  </template>
                </v-data-table>

                <div v-else class="text-center pa-4">
                  <v-progress-circular indeterminate color="primary" />
                  <div class="mt-2">Loading liquidity pools...</div>
                </div>

                <div v-if="!loading && filteredLiquidityPools.length === 0" class="text-center pa-4">
                  <v-icon size="48" class="mb-2 grey--text">mdi-water-off</v-icon>
                  <div class="text-h6 grey--text">No liquidity pools found</div>
                  <div class="text-body-2 grey--text">Try adjusting your filters</div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed, watch } from '@nuxtjs/composition-api'
import useXrpTokenMints from '~/composables/useXrpTokenMints'

export default defineComponent({
  setup() {
    const {
      loading,
      tokenMints,
      liquidityPools,
      selectedTimeRange,
      selectedSortBy,
      selectedSortOrder,
      tokenMintFilters,
      liquidityPoolFilters,
      filteredTokenMints,
      filteredLiquidityPools,
      tokenMintHeaders,
      liquidityPoolHeaders,
      formatAddress,
      formatAmount,
      formatCurrency,
      formatPercentage,
      formatDate,
      copyToClipboard
    } = useXrpTokenMints()

    // Sort options for liquidity pools
    const sortOptions = [
      { label: 'TVL', value: 'tvl' },
      { label: 'Volume', value: 'volume' },
      { label: 'APR', value: 'apr' },
      { label: 'Fees', value: 'fees' }
    ]

    // Computed stats
    const totalMarketCap = computed(() => {
      return filteredTokenMints.value.reduce((sum, mint) => sum + mint.marketcap, 0)
    })

    const totalVolume24H = computed(() => {
      return filteredTokenMints.value.reduce((sum, mint) => sum + mint.volume24H, 0)
    })

    const totalTvl = computed(() => {
      return filteredLiquidityPools.value.reduce((sum, pool) => sum + pool.tvl, 0)
    })

    // Methods
    const refreshData = () => {
      // This would trigger a refetch of the data
      console.log('Refreshing data...')
    }

    const clearTokenMintFilters = () => {
      tokenMintFilters.value = {
        timeRange: '7d',
        minMarketcap: 0,
        minVolume: 0,
        hasLiquidity: false
      }
    }

    const clearLiquidityPoolFilters = () => {
      liquidityPoolFilters.value = {
        minTvl: 0,
        minVolume: 0,
        minApr: 0,
        sortBy: 'tvl',
        sortOrder: 'desc'
      }
    }

    const getImageUrl = (symbol: string) => {
      try {
        return (this as any).$imageUrlBySymbol(symbol.toLowerCase())
      } catch {
        return '/img/default-token.png'
      }
    }

    // Watch for filter changes
    watch(selectedTimeRange, (newValue) => {
      tokenMintFilters.value.timeRange = newValue
    })

    watch(selectedSortBy, (newValue) => {
      liquidityPoolFilters.value.sortBy = newValue
    })

    watch(selectedSortOrder, (newValue) => {
      liquidityPoolFilters.value.sortOrder = newValue
    })

    return {
      // State
      loading,
      selectedTimeRange,
      selectedSortBy,
      selectedSortOrder,
      tokenMintFilters,
      liquidityPoolFilters,
      
      // Computed
      filteredTokenMints,
      filteredLiquidityPools,
      tokenMintHeaders,
      liquidityPoolHeaders,
      totalMarketCap,
      totalVolume24H,
      totalTvl,
      sortOptions,
      
      // Methods
      formatAddress,
      formatAmount,
      formatCurrency,
      formatPercentage,
      formatDate,
      copyToClipboard,
      refreshData,
      clearTokenMintFilters,
      clearLiquidityPoolFilters,
      getImageUrl
    }
  },
  head: {
    title: 'XRP Token Mints & Liquidity Pools',
  },
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style> 