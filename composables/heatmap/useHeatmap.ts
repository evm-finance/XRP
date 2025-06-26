import { ref, computed } from '@nuxtjs/composition-api'

export const useHeatmap = () => {
  const data = ref([])
  const loading = ref(false)
  const error = ref(null)

  const loadData = async () => {
    loading.value = true
    try {
      // Mock data loading
      await new Promise(resolve => setTimeout(resolve, 1000))
      data.value = []
    } catch (err: any) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  return {
    data,
    loading,
    error,
    loadData
  }
} 