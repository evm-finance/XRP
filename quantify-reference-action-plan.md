# Quantify Reference - Action Plan

## Overview
This document provides a detailed, actionable plan for completing the implementation of features inspired by the `quantify/qc-front-end-server` reference project.

## Current Status Assessment

### ✅ Completed (Ready for Production)
- **Core Infrastructure**: Configuration management, data formatting, grid renderers
- **Heatmap System**: Real-time data, multiple timeframes, configurable options
- **Terminal Interface**: Multi-panel layout, market stats, token screener
- **Data Integration**: GraphQL API, WebSocket updates, error handling

### 🔄 In Progress (Needs Completion)
- **Performance Optimizations**: Virtual scrolling, memoization, lazy loading
- **Advanced Analytics**: Technical indicators, portfolio tracking
- **Export Features**: CSV/JSON export, chart export, social sharing

### 📋 Planned (Not Started)
- **UI/UX Enhancements**: Drag & drop, keyboard shortcuts, context menus
- **Advanced Features**: Customizable dashboards, alert system

## Phase 1: Performance Optimizations (Priority: High)

### Task 1.1: Virtual Scrolling Implementation
**Estimated Time**: 3-4 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Virtual Scrolling Composable**
```typescript
// composables/useVirtualScrolling.ts
export function useVirtualScrolling<T>(
  items: Ref<T[]>,
  itemHeight: number = 50,
  overscan: number = 5
) {
  const containerRef = ref<HTMLElement>()
  const scrollTop = ref(0)
  const containerHeight = ref(0)
  
  const visibleItems = computed(() => {
    const startIndex = Math.max(0, Math.floor(scrollTop.value / itemHeight) - overscan)
    const endIndex = Math.min(
      items.value.length,
      Math.ceil((scrollTop.value + containerHeight.value) / itemHeight) + overscan
    )
    
    return items.value.slice(startIndex, endIndex).map((item, index) => ({
      ...item,
      virtualIndex: startIndex + index,
      style: {
        transform: `translateY(${(startIndex + index) * itemHeight}px)`,
        height: `${itemHeight}px`,
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
      },
    }))
  })
  
  return {
    containerRef,
    visibleItems,
    totalHeight: computed(() => items.value.length * itemHeight),
    onScroll: (event: Event) => {
      scrollTop.value = (event.target as HTMLElement).scrollTop
    },
  }
}
```

2. **Update Token Screener Component**
```vue
<!-- components/xrp/XrpTokenScreener.vue -->
<template>
  <div class="virtual-scroll-container" ref="containerRef" @scroll="onScroll">
    <div :style="{ height: `${totalHeight}px`, position: 'relative' }">
      <div
        v-for="item in visibleItems"
        :key="item.virtualIndex"
        :style="item.style"
        class="token-row"
      >
        <token-row :token="item" />
      </div>
    </div>
  </div>
</template>

<script>
import { useVirtualScrolling } from '~/composables/useVirtualScrolling'

export default {
  setup() {
    const { filteredTokens } = useXrpScreener()
    const { containerRef, visibleItems, totalHeight, onScroll } = useVirtualScrolling(
      filteredTokens,
      60 // token row height
    )
    
    return {
      containerRef,
      visibleItems,
      totalHeight,
      onScroll,
    }
  },
}
</script>
```

**Acceptance Criteria:**
- [ ] Handles 10,000+ tokens without performance issues
- [ ] Smooth scrolling with no lag
- [ ] Proper item positioning and sizing
- [ ] Memory usage stays constant regardless of dataset size

### Task 1.2: Memoization Implementation
**Estimated Time**: 2-3 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Memoization Composable**
```typescript
// composables/useMemoization.ts
export function useMemoization<T>(
  fn: () => T,
  deps: Ref<any>[],
  options: { maxSize?: number; ttl?: number } = {}
) {
  const { maxSize = 100, ttl = 5 * 60 * 1000 } = options // 5 minutes default
  const cache = new Map<string, { value: T; timestamp: number }>()
  
  const memoized = computed(() => {
    const key = deps.map(dep => JSON.stringify(dep.value)).join('|')
    const now = Date.now()
    
    // Check cache
    if (cache.has(key)) {
      const cached = cache.get(key)!
      if (now - cached.timestamp < ttl) {
        return cached.value
      } else {
        cache.delete(key)
      }
    }
    
    // Calculate new value
    const result = fn()
    cache.set(key, { value: result, timestamp: now })
    
    // Cleanup if cache is too large
    if (cache.size > maxSize) {
      const oldestKey = cache.keys().next().value
      cache.delete(oldestKey)
    }
    
    return result
  })
  
  return memoized
}
```

2. **Apply to Expensive Computations**
```typescript
// composables/useXrpHeatmap.ts
const processedHeatmapData = useMemoization(
  () => processTokensForHeatmap(heatmapData.value),
  [heatmapData, timeFrame, blockSize, blueTile],
  { maxSize: 50, ttl: 30000 } // 30 seconds cache
)
```

**Acceptance Criteria:**
- [ ] Expensive computations cached appropriately
- [ ] Cache size limited to prevent memory leaks
- [ ] Cache invalidation works correctly
- [ ] Performance improvement measurable

### Task 1.3: Lazy Loading Implementation
**Estimated Time**: 2-3 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Lazy Loading Components**
```vue
<!-- components/xrp/XrpLazyHeatmap.vue -->
<template>
  <div>
    <div v-if="loading" class="loading-placeholder">
      <v-skeleton-loader type="card" height="600" />
    </div>
    <xrp-heatmap-chart v-else v-bind="$props" />
  </div>
</template>

<script>
export default {
  name: 'XrpLazyHeatmap',
  components: {
    XrpHeatmapChart: () => import('./XrpHeatmapChart.vue'),
  },
  props: {
    data: Array,
    updateData: Object,
    tileBody: String,
    tileTooltip: String,
    blockSize: String,
    chartHeight: [String, Number],
  },
  data() {
    return {
      loading: true,
    }
  },
  mounted() {
    // Simulate loading time for heavy component
    setTimeout(() => {
      this.loading = false
    }, 200)
  },
}
</script>
```

2. **Update Bundle Configuration**
```javascript
// nuxt.config.js
export default {
  build: {
    optimization: {
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
          },
          xrp: {
            test: /[\\/]components[\\/]xrp[\\/]/,
            name: 'xrp-components',
            chunks: 'all',
          },
          charts: {
            test: /[\\/]components[\\/]xrp[\\/].*[Cc]hart/,
            name: 'chart-components',
            chunks: 'all',
          },
        },
      },
    },
  },
}
```

**Acceptance Criteria:**
- [ ] Heavy components load only when needed
- [ ] Loading states provide good user feedback
- [ ] Bundle size reduced by at least 30%
- [ ] Initial page load time improved

## Phase 2: Advanced Analytics (Priority: High)

### Task 2.1: Technical Indicators
**Estimated Time**: 4-5 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Technical Indicators Composable**
```typescript
// composables/useTechnicalIndicators.ts
export function useTechnicalIndicators() {
  const calculateSMA = (prices: number[], period: number) => {
    if (prices.length < period) return null
    const sum = prices.slice(-period).reduce((acc, price) => acc + price, 0)
    return sum / period
  }
  
  const calculateEMA = (prices: number[], period: number) => {
    if (prices.length < period) return null
    const multiplier = 2 / (period + 1)
    let ema = prices[0]
    
    for (let i = 1; i < prices.length; i++) {
      ema = (prices[i] * multiplier) + (ema * (1 - multiplier))
    }
    
    return ema
  }
  
  const calculateRSI = (prices: number[], period: number = 14) => {
    if (prices.length < period + 1) return null
    
    const gains = []
    const losses = []
    
    for (let i = 1; i < prices.length; i++) {
      const change = prices[i] - prices[i - 1]
      gains.push(change > 0 ? change : 0)
      losses.push(change < 0 ? Math.abs(change) : 0)
    }
    
    const avgGain = gains.slice(-period).reduce((sum, gain) => sum + gain, 0) / period
    const avgLoss = losses.slice(-period).reduce((sum, loss) => sum + loss, 0) / period
    
    if (avgLoss === 0) return 100
    
    const rs = avgGain / avgLoss
    return 100 - (100 / (1 + rs))
  }
  
  const calculateMACD = (prices: number[], fastPeriod: number = 12, slowPeriod: number = 26) => {
    const ema12 = calculateEMA(prices, fastPeriod)
    const ema26 = calculateEMA(prices, slowPeriod)
    
    if (!ema12 || !ema26) return null
    
    return {
      macd: ema12 - ema26,
      signal: calculateEMA([ema12 - ema26], 9),
      histogram: (ema12 - ema26) - (calculateEMA([ema12 - ema26], 9) || 0),
    }
  }
  
  const calculateBollingerBands = (prices: number[], period: number = 20, stdDev: number = 2) => {
    if (prices.length < period) return null
    
    const sma = calculateSMA(prices, period)
    if (!sma) return null
    
    const squaredDiffs = prices.slice(-period).map(price => Math.pow(price - sma, 2))
    const variance = squaredDiffs.reduce((sum, diff) => sum + diff, 0) / period
    const standardDeviation = Math.sqrt(variance)
    
    return {
      upper: sma + (stdDev * standardDeviation),
      middle: sma,
      lower: sma - (stdDev * standardDeviation),
    }
  }
  
  return {
    calculateSMA,
    calculateEMA,
    calculateRSI,
    calculateMACD,
    calculateBollingerBands,
  }
}
```

2. **Create Technical Analysis Component**
```vue
<!-- components/xrp/XrpTechnicalAnalysis.vue -->
<template>
  <v-card tile outlined>
    <v-card-title class="py-2">
      <span class="subtitle-2">Technical Analysis</span>
      <v-spacer />
      <v-select
        v-model="selectedIndicator"
        :items="indicatorOptions"
        label="Indicator"
        dense
        outlined
        hide-details
        style="max-width: 150px;"
      />
    </v-card-title>
    
    <v-card-text class="py-2">
      <div v-if="selectedIndicator === 'rsi'">
        <div class="d-flex justify-space-between align-center">
          <span>RSI (14): {{ rsiValue?.toFixed(2) || 'N/A' }}</span>
          <v-chip
            :color="getRSIColor(rsiValue)"
            small
            outlined
          >
            {{ getRSISignal(rsiValue) }}
          </v-chip>
        </div>
        <v-progress-linear
          :value="rsiValue || 0"
          :color="getRSIColor(rsiValue)"
          height="8"
          class="mt-2"
        />
      </div>
      
      <div v-else-if="selectedIndicator === 'macd'">
        <div class="d-flex justify-space-between">
          <div>
            <div class="text-caption">MACD: {{ macdValue?.macd?.toFixed(4) || 'N/A' }}</div>
            <div class="text-caption">Signal: {{ macdValue?.signal?.toFixed(4) || 'N/A' }}</div>
            <div class="text-caption">Histogram: {{ macdValue?.histogram?.toFixed(4) || 'N/A' }}</div>
          </div>
          <v-chip
            :color="getMACDColor(macdValue)"
            small
            outlined
          >
            {{ getMACDSignal(macdValue) }}
          </v-chip>
        </div>
      </div>
      
      <div v-else-if="selectedIndicator === 'bollinger'">
        <div class="d-flex justify-space-between">
          <div>
            <div class="text-caption">Upper: {{ bollingerValue?.upper?.toFixed(4) || 'N/A' }}</div>
            <div class="text-caption">Middle: {{ bollingerValue?.middle?.toFixed(4) || 'N/A' }}</div>
            <div class="text-caption">Lower: {{ bollingerValue?.lower?.toFixed(4) || 'N/A' }}</div>
          </div>
          <v-chip
            :color="getBollingerColor(bollingerValue)"
            small
            outlined
          >
            {{ getBollingerSignal(bollingerValue) }}
          </v-chip>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
import { useTechnicalIndicators } from '~/composables/useTechnicalIndicators'

export default {
  props: {
    prices: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const selectedIndicator = ref('rsi')
    const { calculateRSI, calculateMACD, calculateBollingerBands } = useTechnicalIndicators()
    
    const indicatorOptions = [
      { text: 'RSI', value: 'rsi' },
      { text: 'MACD', value: 'macd' },
      { text: 'Bollinger Bands', value: 'bollinger' },
    ]
    
    const rsiValue = computed(() => calculateRSI(props.prices))
    const macdValue = computed(() => calculateMACD(props.prices))
    const bollingerValue = computed(() => calculateBollingerBands(props.prices))
    
    const getRSIColor = (value) => {
      if (!value) return 'grey'
      if (value > 70) return 'error'
      if (value < 30) return 'success'
      return 'warning'
    }
    
    const getRSISignal = (value) => {
      if (!value) return 'N/A'
      if (value > 70) return 'Overbought'
      if (value < 30) return 'Oversold'
      return 'Neutral'
    }
    
    // Similar functions for MACD and Bollinger Bands...
    
    return {
      selectedIndicator,
      indicatorOptions,
      rsiValue,
      macdValue,
      bollingerValue,
      getRSIColor,
      getRSISignal,
      // ... other functions
    }
  },
}
</script>
```

**Acceptance Criteria:**
- [ ] RSI calculation accurate and responsive
- [ ] MACD calculation with signal line and histogram
- [ ] Bollinger Bands with upper, middle, lower bands
- [ ] Visual indicators with color coding
- [ ] Real-time updates when price data changes

### Task 2.2: Portfolio Analytics
**Estimated Time**: 3-4 days
**Dependencies**: Technical indicators

**Implementation Steps:**
1. **Create Portfolio Analytics Composable**
```typescript
// composables/usePortfolioAnalytics.ts
export function usePortfolioAnalytics() {
  const calculatePortfolioValue = (holdings: any[], prices: any[]) => {
    return holdings.reduce((total, holding) => {
      const price = prices.find(p => p.currency === holding.currency)?.price_usd || 0
      return total + (holding.balance * price)
    }, 0)
  }
  
  const calculatePortfolioReturn = (initialValue: number, currentValue: number) => {
    if (initialValue === 0) return 0
    return ((currentValue - initialValue) / initialValue) * 100
  }
  
  const calculateSharpeRatio = (returns: number[], riskFreeRate: number = 0.02) => {
    if (returns.length === 0) return 0
    
    const avgReturn = returns.reduce((sum, ret) => sum + ret, 0) / returns.length
    const variance = returns.reduce((sum, ret) => sum + Math.pow(ret - avgReturn, 2), 0) / returns.length
    const stdDev = Math.sqrt(variance)
    
    if (stdDev === 0) return 0
    
    return (avgReturn - riskFreeRate) / stdDev
  }
  
  const calculateMaxDrawdown = (values: number[]) => {
    let maxDrawdown = 0
    let peak = values[0]
    
    for (const value of values) {
      if (value > peak) {
        peak = value
      }
      const drawdown = (peak - value) / peak
      if (drawdown > maxDrawdown) {
        maxDrawdown = drawdown
      }
    }
    
    return maxDrawdown * 100
  }
  
  const calculateVolatility = (returns: number[]) => {
    if (returns.length === 0) return 0
    
    const avgReturn = returns.reduce((sum, ret) => sum + ret, 0) / returns.length
    const variance = returns.reduce((sum, ret) => sum + Math.pow(ret - avgReturn, 2), 0) / returns.length
    
    return Math.sqrt(variance) * 100
  }
  
  return {
    calculatePortfolioValue,
    calculatePortfolioReturn,
    calculateSharpeRatio,
    calculateMaxDrawdown,
    calculateVolatility,
  }
}
```

2. **Create Portfolio Dashboard Component**
```vue
<!-- components/xrp/XrpPortfolioDashboard.vue -->
<template>
  <v-card tile outlined>
    <v-card-title class="py-2">
      <span class="subtitle-2">Portfolio Analytics</span>
    </v-card-title>
    
    <v-card-text class="py-2">
      <v-row dense>
        <v-col cols="12" sm="6" md="3">
          <div class="text-center">
            <div class="text-h6">{{ formatCurrency(portfolioValue) }}</div>
            <div class="text-caption">Total Value</div>
          </div>
        </v-col>
        
        <v-col cols="12" sm="6" md="3">
          <div class="text-center">
            <div class="text-h6" :class="portfolioReturn >= 0 ? 'success--text' : 'error--text'">
              {{ formatPercentage(portfolioReturn) }}
            </div>
            <div class="text-caption">Total Return</div>
          </div>
        </v-col>
        
        <v-col cols="12" sm="6" md="3">
          <div class="text-center">
            <div class="text-h6">{{ sharpeRatio.toFixed(2) }}</div>
            <div class="text-caption">Sharpe Ratio</div>
          </div>
        </v-col>
        
        <v-col cols="12" sm="6" md="3">
          <div class="text-center">
            <div class="text-h6 error--text">{{ maxDrawdown.toFixed(2) }}%</div>
            <div class="text-caption">Max Drawdown</div>
          </div>
        </v-col>
      </v-row>
      
      <v-divider class="my-3" />
      
      <div class="text-caption mb-2">Asset Allocation</div>
      <div v-for="holding in holdings" :key="holding.currency" class="d-flex justify-space-between align-center mb-1">
        <div class="d-flex align-center">
          <v-avatar size="24" class="mr-2">
            <v-img :src="getTokenIcon(holding.currency)" />
          </v-avatar>
          <span>{{ holding.currency }}</span>
        </div>
        <span>{{ formatPercentage(holding.allocation) }}</span>
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
import { usePortfolioAnalytics } from '~/composables/usePortfolioAnalytics'

export default {
  props: {
    holdings: {
      type: Array,
      required: true,
    },
    prices: {
      type: Array,
      required: true,
    },
    historicalData: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const {
      calculatePortfolioValue,
      calculatePortfolioReturn,
      calculateSharpeRatio,
      calculateMaxDrawdown,
      calculateVolatility,
    } = usePortfolioAnalytics()
    
    const portfolioValue = computed(() => calculatePortfolioValue(props.holdings, props.prices))
    const portfolioReturn = computed(() => {
      // Calculate return based on historical data
      if (props.historicalData.length < 2) return 0
      const initialValue = props.historicalData[0].value
      const currentValue = props.historicalData[props.historicalData.length - 1].value
      return calculatePortfolioReturn(initialValue, currentValue)
    })
    
    const returns = computed(() => {
      // Calculate daily returns from historical data
      const returns = []
      for (let i = 1; i < props.historicalData.length; i++) {
        const prevValue = props.historicalData[i - 1].value
        const currentValue = props.historicalData[i].value
        returns.push((currentValue - prevValue) / prevValue)
      }
      return returns
    })
    
    const sharpeRatio = computed(() => calculateSharpeRatio(returns.value))
    const maxDrawdown = computed(() => calculateMaxDrawdown(props.historicalData.map(d => d.value)))
    const volatility = computed(() => calculateVolatility(returns.value))
    
    return {
      portfolioValue,
      portfolioReturn,
      sharpeRatio,
      maxDrawdown,
      volatility,
    }
  },
}
</script>
```

**Acceptance Criteria:**
- [ ] Portfolio value calculation accurate
- [ ] Return calculation with proper time periods
- [ ] Risk metrics (Sharpe ratio, max drawdown, volatility)
- [ ] Asset allocation visualization
- [ ] Real-time updates when holdings change

## Phase 3: Export Features (Priority: Medium)

### Task 3.1: Data Export Implementation
**Estimated Time**: 2-3 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Export Utilities**
```typescript
// composables/useExport.ts
export function useExport() {
  const exportToCSV = (data: any[], filename: string) => {
    if (data.length === 0) return
    
    const headers = Object.keys(data[0])
    const csvContent = [
      headers.join(','),
      ...data.map(row => headers.map(header => `"${row[header]}"`).join(','))
    ].join('\n')
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `${filename}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }
  
  const exportToJSON = (data: any, filename: string) => {
    const jsonContent = JSON.stringify(data, null, 2)
    const blob = new Blob([jsonContent], { type: 'application/json' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `${filename}.json`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }
  
  const exportChartToPNG = (chartRef: HTMLElement, filename: string) => {
    html2canvas(chartRef).then(canvas => {
      const link = document.createElement('a')
      link.download = `${filename}.png`
      link.href = canvas.toDataURL()
      link.click()
    })
  }
  
  return {
    exportToCSV,
    exportToJSON,
    exportChartToPNG,
  }
}
```

2. **Add Export Buttons to Components**
```vue
<!-- Add to XrpTerminal.vue -->
<template>
  <v-btn
    icon
    small
    class="mx-1"
    @click="showExportMenu = !showExportMenu"
  >
    <v-icon>mdi-download</v-icon>
  </v-btn>
  
  <v-menu
    v-model="showExportMenu"
    :close-on-content-click="false"
    offset-y
  >
    <v-list dense>
      <v-list-item @click="exportHeatmapData('csv')">
        <v-list-item-icon>
          <v-icon>mdi-file-csv</v-icon>
        </v-list-item-icon>
        <v-list-item-title>Export to CSV</v-list-item-title>
      </v-list-item>
      
      <v-list-item @click="exportHeatmapData('json')">
        <v-list-item-icon>
          <v-icon>mdi-code-json</v-icon>
        </v-list-item-icon>
        <v-list-item-title>Export to JSON</v-list-item-title>
      </v-list-item>
      
      <v-list-item @click="exportHeatmapData('png')">
        <v-list-item-icon>
          <v-icon>mdi-image</v-icon>
        </v-list-item-icon>
        <v-list-item-title>Export as PNG</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { useExport } from '~/composables/useExport'

export default {
  setup() {
    const { exportToCSV, exportToJSON, exportChartToPNG } = useExport()
    
    const exportHeatmapData = (format) => {
      switch (format) {
        case 'csv':
          exportToCSV(heatmapData.value, 'xrp-heatmap')
          break
        case 'json':
          exportToJSON(heatmapData.value, 'xrp-heatmap')
          break
        case 'png':
          exportChartToPNG(chartRef.value, 'xrp-heatmap')
          break
      }
      showExportMenu.value = false
    }
    
    return {
      exportHeatmapData,
    }
  },
}
</script>
```

**Acceptance Criteria:**
- [ ] CSV export with proper formatting
- [ ] JSON export with all data fields
- [ ] PNG export of charts and heatmaps
- [ ] Proper file naming and download
- [ ] Export menu accessible from main components

## Phase 4: UI/UX Enhancements (Priority: Medium)

### Task 4.1: Keyboard Shortcuts
**Estimated Time**: 2-3 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Keyboard Shortcuts Composable**
```typescript
// composables/useKeyboardShortcuts.ts
export function useKeyboardShortcuts() {
  const shortcuts = new Map()
  
  const registerShortcut = (key: string, callback: Function, description: string) => {
    shortcuts.set(key, { callback, description })
  }
  
  const unregisterShortcut = (key: string) => {
    shortcuts.delete(key)
  }
  
  const handleKeydown = (event: KeyboardEvent) => {
    const key = buildKeyString(event)
    const shortcut = shortcuts.get(key)
    
    if (shortcut) {
      event.preventDefault()
      shortcut.callback()
    }
  }
  
  const buildKeyString = (event: KeyboardEvent) => {
    const parts = []
    if (event.ctrlKey) parts.push('Ctrl')
    if (event.altKey) parts.push('Alt')
    if (event.shiftKey) parts.push('Shift')
    if (event.metaKey) parts.push('Meta')
    parts.push(event.key.toUpperCase())
    return parts.join('+')
  }
  
  const getShortcutsList = () => {
    return Array.from(shortcuts.entries()).map(([key, { description }]) => ({
      key,
      description,
    }))
  }
  
  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })
  
  onBeforeUnmount(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
  
  return {
    registerShortcut,
    unregisterShortcut,
    getShortcutsList,
  }
}
```

2. **Implement Shortcuts in Terminal**
```typescript
// In XrpTerminal.vue setup()
const { registerShortcut } = useKeyboardShortcuts()

registerShortcut('Ctrl+R', refreshAll, 'Refresh all data')
registerShortcut('Ctrl+F', () => searchQuery.value = '', 'Clear search')
registerShortcut('Ctrl+S', () => showSettings.value = !showSettings.value, 'Toggle settings')
registerShortcut('F11', toggleFullscreen, 'Toggle fullscreen')
registerShortcut('Escape', () => showSettings.value = false, 'Close settings')
registerShortcut('Ctrl+E', exportHeatmapData, 'Export data')
```

**Acceptance Criteria:**
- [ ] All shortcuts work correctly
- [ ] No conflicts with browser shortcuts
- [ ] Shortcuts list available in help menu
- [ ] Visual feedback when shortcuts are used

### Task 4.2: Context Menus
**Estimated Time**: 2-3 days
**Dependencies**: None

**Implementation Steps:**
1. **Create Context Menu Component**
```vue
<!-- components/xrp/XrpContextMenu.vue -->
<template>
  <v-menu
    v-model="show"
    :position-x="x"
    :position-y="y"
    :close-on-content-click="false"
  >
    <v-list dense>
      <v-list-item @click="viewDetails">
        <v-list-item-icon>
          <v-icon>mdi-open-in-new</v-icon>
        </v-list-item-icon>
        <v-list-item-title>View Details</v-list-item-title>
      </v-list-item>
      
      <v-list-item @click="addToFavorites">
        <v-list-item-icon>
          <v-icon>mdi-star</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ isFavorite ? 'Remove from' : 'Add to' }} Favorites</v-list-item-title>
      </v-list-item>
      
      <v-list-item @click="copyInfo">
        <v-list-item-icon>
          <v-icon>mdi-content-copy</v-icon>
        </v-list-item-icon>
        <v-list-item-title>Copy Info</v-list-item-title>
      </v-list-item>
      
      <v-divider />
      
      <v-list-item @click="exportData">
        <v-list-item-icon>
          <v-icon>mdi-download</v-icon>
        </v-list-item-icon>
        <v-list-item-title>Export Data</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
export default {
  props: {
    show: Boolean,
    x: Number,
    y: Number,
    item: Object,
  },
  emits: ['update:show', 'view-details', 'toggle-favorite', 'copy-info', 'export-data'],
  setup(props, { emit }) {
    const isFavorite = computed(() => {
      // Check if item is in favorites
      return false // Implement based on your favorites system
    })
    
    const viewDetails = () => {
      emit('view-details', props.item)
      emit('update:show', false)
    }
    
    const addToFavorites = () => {
      emit('toggle-favorite', props.item)
      emit('update:show', false)
    }
    
    const copyInfo = () => {
      emit('copy-info', props.item)
      emit('update:show', false)
    }
    
    const exportData = () => {
      emit('export-data', props.item)
      emit('update:show', false)
    }
    
    return {
      isFavorite,
      viewDetails,
      addToFavorites,
      copyInfo,
      exportData,
    }
  },
}
</script>
```

2. **Add Context Menu to Grid Components**
```vue
<!-- In grid components -->
<template>
  <div
    @contextmenu="showContextMenu"
    @click="handleClick"
  >
    <!-- Grid content -->
  </div>
  
  <xrp-context-menu
    v-model:show="contextMenu.show"
    :x="contextMenu.x"
    :y="contextMenu.y"
    :item="contextMenu.item"
    @view-details="viewDetails"
    @toggle-favorite="toggleFavorite"
    @copy-info="copyInfo"
    @export-data="exportData"
  />
</template>

<script>
const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  item: null,
})

const showContextMenu = (event) => {
  event.preventDefault()
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    item: event.target.dataset.item ? JSON.parse(event.target.dataset.item) : null,
  }
}
</script>
```

**Acceptance Criteria:**
- [ ] Context menu appears on right-click
- [ ] Menu positioned correctly relative to cursor
- [ ] All menu actions work properly
- [ ] Menu closes when clicking outside
- [ ] Menu items are contextually relevant

## Implementation Timeline

### Week 1: Performance Optimizations
- **Days 1-2**: Virtual scrolling implementation
- **Days 3-4**: Memoization and caching
- **Day 5**: Lazy loading and bundle optimization

### Week 2: Advanced Analytics
- **Days 1-3**: Technical indicators (RSI, MACD, Bollinger Bands)
- **Days 4-5**: Portfolio analytics and risk metrics

### Week 3: Export Features & UI Enhancements
- **Days 1-2**: Data export functionality
- **Days 3-4**: Keyboard shortcuts
- **Day 5**: Context menus

### Week 4: Testing & Polish
- **Days 1-3**: Comprehensive testing
- **Days 4-5**: Bug fixes and performance tuning

## Success Metrics

### Performance Metrics
- [ ] Page load time < 2 seconds
- [ ] Virtual scrolling handles 10,000+ items smoothly
- [ ] Memory usage stays constant with large datasets
- [ ] Bundle size reduced by 30%

### Feature Metrics
- [ ] All technical indicators calculate correctly
- [ ] Portfolio analytics provide accurate metrics
- [ ] Export functionality works for all formats
- [ ] Keyboard shortcuts improve user efficiency

### User Experience Metrics
- [ ] Context menus provide intuitive interactions
- [ ] Loading states provide clear feedback
- [ ] Error handling is graceful and informative
- [ ] Mobile responsiveness maintained

## Risk Mitigation

### Technical Risks
- **Performance Issues**: Implement virtual scrolling and memoization first
- **Browser Compatibility**: Test across major browsers
- **Memory Leaks**: Proper cleanup in lifecycle hooks

### Feature Risks
- **Complex Calculations**: Start with simple indicators, add complexity gradually
- **Data Accuracy**: Validate all calculations with known datasets
- **User Adoption**: Provide clear documentation and help

## Conclusion

This action plan provides a structured approach to implementing the remaining features inspired by the quantify reference project. The plan prioritizes performance optimizations and advanced analytics, which will provide the most value to users.

**Key Success Factors:**
1. **Performance First**: Implement optimizations before adding features
2. **Incremental Development**: Build and test features incrementally
3. **User Feedback**: Gather feedback early and often
4. **Quality Assurance**: Comprehensive testing at each phase

**Expected Outcomes:**
- Professional-grade cryptocurrency analytics platform
- Excellent performance with large datasets
- Advanced technical analysis capabilities
- Intuitive and efficient user interface

The foundation established through the reference project analysis provides a solid base for these enhancements, ensuring a high-quality, maintainable, and scalable application. 