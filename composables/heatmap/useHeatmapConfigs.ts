import { ref, computed } from '@nuxtjs/composition-api'
import { useXrpConfigs } from '~/composables/useXrpConfigs'

export const useHeatmapConfigs = () => {
  const configs = ref([])
  const activeConfig = ref(null)
  const loading = ref(false)

  // Use XRP configs for AMM heatmap
  const {
    blockSize,
    blockSizeOptions,
    blueTile,
    displayGainersAndLosers,
  } = useXrpConfigs()

  const loadConfigs = async () => {
    loading.value = true
    try {
      // Mock configs loading
      await new Promise(resolve => setTimeout(resolve, 500))
      configs.value = []
    } catch (err) {
      console.error('Failed to load configs:', err)
    } finally {
      loading.value = false
    }
  }

  const setActiveConfig = (config: any) => {
    activeConfig.value = config
  }

  return {
    configs,
    activeConfig,
    loading,
    loadConfigs,
    setActiveConfig,
    // AMM heatmap specific configs
    blockSize,
    blockSizeOptions,
    blueTile,
    displayGainersAndLosers,
  }
} 