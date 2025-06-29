import { ref, computed, watch } from '@nuxtjs/composition-api'
import { ApolloError } from '@apollo/client/errors'

export interface ErrorHandlerConfig {
  maxRetries?: number
  retryDelay?: number
  enableFallback?: boolean
  showUserErrors?: boolean
  logErrors?: boolean
}

export interface FallbackData {
  [key: string]: any
}

export interface ErrorState {
  hasError: boolean
  errorMessage: string
  errorType: 'network' | 'graphql' | 'validation' | 'unknown'
  retryCount: number
  lastError?: ApolloError
  isRetrying: boolean
}

export function useXrpGraphQLErrorHandler(config: ErrorHandlerConfig = {}) {
  const {
    maxRetries = 3,
    retryDelay = 1000,
    enableFallback = true,
    showUserErrors = true,
    logErrors = true
  } = config

  // Error state
  const errorState = ref<ErrorState>({
    hasError: false,
    errorMessage: '',
    errorType: 'unknown',
    retryCount: 0,
    isRetrying: false
  })

  // Fallback data storage
  const fallbackData = ref<FallbackData>({})

  // Retry timeout reference
  let retryTimeout: NodeJS.Timeout | null = null

  // Error categorization
  const categorizeError = (error: ApolloError): 'network' | 'graphql' | 'validation' | 'unknown' => {
    if (error.networkError) {
      return 'network'
    }
    if (error.graphQLErrors && error.graphQLErrors.length > 0) {
      return 'graphql'
    }
    if (error.message.includes('validation') || error.message.includes('invalid')) {
      return 'validation'
    }
    return 'unknown'
  }

  // User-friendly error messages
  const getUserFriendlyMessage = (errorType: string, originalMessage: string): string => {
    switch (errorType) {
      case 'network':
        return 'Network connection issue. Please check your internet connection and try again.'
      case 'graphql':
        return 'Data loading issue. Please refresh the page or try again later.'
      case 'validation':
        return 'Invalid request. Please check your input and try again.'
      case 'unknown':
      default:
        return 'An unexpected error occurred. Please try again or contact support if the problem persists.'
    }
  }

  // Handle error with retry logic
  const handleError = async (
    error: ApolloError,
    retryFunction?: () => Promise<any>,
    fallbackDataKey?: string
  ): Promise<boolean> => {
    const errorType = categorizeError(error)
    const userMessage = getUserFriendlyMessage(errorType, error.message)

    // Log error if enabled
    if (logErrors) {
      console.error('GraphQL Error:', {
        type: errorType,
        message: error.message,
        originalError: error,
        retryCount: errorState.value.retryCount,
        timestamp: new Date().toISOString()
      })
    }

    // Update error state
    errorState.value = {
      hasError: true,
      errorMessage: showUserErrors ? userMessage : error.message,
      errorType,
      retryCount: errorState.value.retryCount,
      lastError: error,
      isRetrying: false
    }

    // Check if we should retry
    if (errorState.value.retryCount < maxRetries && retryFunction) {
      return await retryWithBackoff(retryFunction)
    }

    // Use fallback data if available and enabled
    if (enableFallback && fallbackDataKey && fallbackData.value[fallbackDataKey]) {
      console.log(`Using fallback data for: ${fallbackDataKey}`)
      return true
    }

    return false
  }

  // Retry with exponential backoff
  const retryWithBackoff = async (retryFunction: () => Promise<any>): Promise<boolean> => {
    errorState.value.isRetrying = true
    errorState.value.retryCount++

    const delay = retryDelay * Math.pow(2, errorState.value.retryCount - 1)

    try {
      await new Promise(resolve => {
        retryTimeout = setTimeout(resolve, delay)
      })

      const result = await retryFunction()
      
      // Success - clear error state
      errorState.value = {
        hasError: false,
        errorMessage: '',
        errorType: 'unknown',
        retryCount: 0,
        isRetrying: false
      }

      return true
    } catch (error) {
      // Retry failed
      errorState.value.isRetrying = false
      return false
    }
  }

  // Set fallback data
  const setFallbackData = (key: string, data: any) => {
    fallbackData.value[key] = data
  }

  // Get fallback data
  const getFallbackData = (key: string) => {
    return fallbackData.value[key] || null
  }

  // Clear error state
  const clearError = () => {
    errorState.value = {
      hasError: false,
      errorMessage: '',
      errorType: 'unknown',
      retryCount: 0,
      isRetrying: false
    }
  }

  // Reset retry count
  const resetRetryCount = () => {
    errorState.value.retryCount = 0
  }

  // Cleanup on unmount
  const cleanup = () => {
    if (retryTimeout) {
      clearTimeout(retryTimeout)
      retryTimeout = null
    }
  }

  // Computed properties
  const canRetry = computed(() => errorState.value.retryCount < maxRetries)
  const isNetworkError = computed(() => errorState.value.errorType === 'network')
  const isGraphQLError = computed(() => errorState.value.errorType === 'graphql')
  const isValidationError = computed(() => errorState.value.errorType === 'validation')

  return {
    // State
    errorState: computed(() => errorState.value),
    canRetry,
    isNetworkError,
    isGraphQLError,
    isValidationError,

    // Methods
    handleError,
    setFallbackData,
    getFallbackData,
    clearError,
    resetRetryCount,
    cleanup
  }
} 