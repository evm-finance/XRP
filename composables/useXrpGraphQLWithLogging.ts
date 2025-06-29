import { useQuery, useMutation, useSubscription } from '@vue/apollo-composable/dist'
import { computed, watch, onBeforeUnmount } from '@nuxtjs/composition-api'
import useXrpGraphQLLogging from './useXrpGraphQLLogging'
import { useXrpGraphQLErrorHandler } from './useXrpGraphQLErrorHandler'
import { useXrpFallbackData } from './useXrpFallbackData'

export default function useXrpGraphQLWithLogging() {
  const { startQuery, completeQuery } = useXrpGraphQLLogging()
  const { getFallbackData } = useXrpFallbackData()

  // Wrapper for useQuery with logging and graceful error handling
  const useLoggedQuery = (query: any, variables?: any, options?: any) => {
    const queryName = query.definitions?.[0]?.name?.value || 'UnknownQuery'
    
    // Extract the actual GraphQL query string
    const queryString = query.loc?.source?.body || query.toString()
    
    // Start logging the query
    const logId = startQuery(queryName, 'query', variables, queryString)
    
    // Create error handler for this query
    const errorHandler = useXrpGraphQLErrorHandler({
      maxRetries: 3,
      retryDelay: 1000,
      enableFallback: true,
      showUserErrors: true,
      logErrors: true
    })

    // Set fallback data for this query
    const fallbackData = getFallbackData(queryName)
    if (fallbackData) {
      errorHandler.setFallbackData(queryName, fallbackData)
    }
    
    // Use the original useQuery
    const result = useQuery(query, variables, options)
    
    // Watch for results and errors with graceful handling
    watch(result.result, (newResult) => {
      if (newResult) {
        completeQuery(logId, true)
        errorHandler.clearError() // Clear any previous errors
      }
    }, { immediate: true })
    
    watch(result.error, async (newError) => {
      if (newError) {
        completeQuery(logId, false, newError)
        
        // Handle error gracefully
        const retryFunction = () => {
          const refetchResult = result.refetch()
          return refetchResult ? refetchResult : Promise.resolve()
        }
        const success = await errorHandler.handleError(newError, retryFunction, queryName)
        
        if (success && fallbackData) {
          // Use fallback data if retry failed but we have fallback
          console.log(`Using fallback data for ${queryName} after error handling`)
        }
      }
    }, { immediate: true })
    
    // Cleanup on unmount
    onBeforeUnmount(() => {
      errorHandler.cleanup()
    })
    
    return {
      ...result,
      errorState: errorHandler.errorState,
      canRetry: errorHandler.canRetry,
      isNetworkError: errorHandler.isNetworkError,
      isGraphQLError: errorHandler.isGraphQLError,
      isValidationError: errorHandler.isValidationError,
      clearError: errorHandler.clearError,
      resetRetryCount: errorHandler.resetRetryCount
    }
  }

  // Wrapper for useMutation with logging and graceful error handling
  const useLoggedMutation = (mutation: any, options?: any) => {
    const mutationName = mutation.definitions?.[0]?.name?.value || 'UnknownMutation'
    
    // Extract the actual GraphQL mutation string
    const mutationString = mutation.loc?.source?.body || mutation.toString()
    
    // Create error handler for this mutation
    const errorHandler = useXrpGraphQLErrorHandler({
      maxRetries: 2, // Fewer retries for mutations
      retryDelay: 2000,
      enableFallback: false, // No fallback for mutations
      showUserErrors: true,
      logErrors: true
    })
    
    const result = useMutation(mutation, options)
    
    // For mutations, we need to wrap the mutate function
    const originalMutate = result.mutate
    result.mutate = async (...args: any[]) => {
      const logId = startQuery(mutationName, 'mutation', args[0], mutationString)
      
      try {
        const response = await originalMutate(...args)
        completeQuery(logId, true)
        errorHandler.clearError()
        return response
      } catch (error) {
        completeQuery(logId, false, error)
        
        // Handle error gracefully
        const retryFunction = () => originalMutate(...args)
        await errorHandler.handleError(error, retryFunction)
        
        throw error // Re-throw for component handling
      }
    }
    
    // Cleanup on unmount
    onBeforeUnmount(() => {
      errorHandler.cleanup()
    })
    
    return {
      ...result,
      errorState: errorHandler.errorState,
      canRetry: errorHandler.canRetry,
      isNetworkError: errorHandler.isNetworkError,
      isGraphQLError: errorHandler.isGraphQLError,
      isValidationError: errorHandler.isValidationError,
      clearError: errorHandler.clearError,
      resetRetryCount: errorHandler.resetRetryCount
    }
  }

  // Wrapper for useSubscription with logging and graceful error handling
  const useLoggedSubscription = (subscription: any, variables?: any, options?: any) => {
    const subscriptionName = subscription.definitions?.[0]?.name?.value || 'UnknownSubscription'
    
    // Extract the actual GraphQL subscription string
    const subscriptionString = subscription.loc?.source?.body || subscription.toString()
    
    // Start logging the subscription
    const logId = startQuery(subscriptionName, 'subscription', variables, subscriptionString)
    
    // Create error handler for this subscription
    const errorHandler = useXrpGraphQLErrorHandler({
      maxRetries: 5, // More retries for subscriptions
      retryDelay: 3000,
      enableFallback: true,
      showUserErrors: false, // Don't show user errors for subscriptions
      logErrors: true
    })
    
    // Use the original useSubscription
    const result = useSubscription(subscription, variables, options)
    
    // Watch for results and errors with graceful handling
    watch(result.result, (newResult) => {
      if (newResult) {
        completeQuery(logId, true)
        errorHandler.clearError()
      }
    }, { immediate: true })
    
    watch(result.error, async (newError) => {
      if (newError) {
        completeQuery(logId, false, newError)
        
        // Handle error gracefully - subscriptions auto-retry
        await errorHandler.handleError(newError, undefined, subscriptionName)
      }
    }, { immediate: true })
    
    // Cleanup on unmount
    onBeforeUnmount(() => {
      errorHandler.cleanup()
    })
    
    return {
      ...result,
      errorState: errorHandler.errorState,
      canRetry: errorHandler.canRetry,
      isNetworkError: errorHandler.isNetworkError,
      isGraphQLError: errorHandler.isGraphQLError,
      isValidationError: errorHandler.isValidationError,
      clearError: errorHandler.clearError,
      resetRetryCount: errorHandler.resetRetryCount
    }
  }

  return {
    useLoggedQuery,
    useLoggedMutation,
    useLoggedSubscription
  }
} 