import { onMounted } from '@nuxtjs/composition-api'

export const useInitTheme = () => {
  onMounted(() => {
    // Initialize theme on mount
    const theme = localStorage.getItem('theme') || 'dark'
    document.documentElement.setAttribute('data-theme', theme)
  })

  return {
    // Theme utilities can be added here
  }
} 