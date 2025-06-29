import { ref, reactive } from '@nuxtjs/composition-api'

export interface GraphQLLogEntry {
  id: string
  queryName: string
  operation: 'query' | 'mutation' | 'subscription'
  queryString?: string
  variables?: any
  startTime: number
  endTime?: number
  duration?: number
  success: boolean
  error?: any
  userContext?: {
    page: string
    walletAddress?: string
    timestamp: string
  }
  networkInfo?: {
    endpoint: string
    userAgent: string
  }
}

export interface GraphQLErrorStats {
  totalErrors: number
  errorsByQuery: Record<string, number>
  errorsByType: Record<string, number>
  lastError?: GraphQLLogEntry
}

export interface GraphQLPerformanceStats {
  totalQueries: number
  averageResponseTime: number
  slowestQueries: Array<{ queryName: string; duration: number }>
  queriesByEndpoint: Record<string, number>
}

class XRPGraphQLLogger {
  private logs = reactive<GraphQLLogEntry[]>([])
  private errorStats = reactive<GraphQLErrorStats>({
    totalErrors: 0,
    errorsByQuery: {},
    errorsByType: {}
  })
  private performanceStats = reactive<GraphQLPerformanceStats>({
    totalQueries: 0,
    averageResponseTime: 0,
    slowestQueries: [],
    queriesByEndpoint: {}
  })

  private formatQueryForLogging(queryString: string): string {
    // Remove extra whitespace and format for logging
    return queryString
      .replace(/\s+/g, ' ')
      .replace(/\s*{\s*/g, ' { ')
      .replace(/\s*}\s*/g, ' } ')
      .replace(/\s*\(\s*/g, ' ( ')
      .replace(/\s*\)\s*/g, ' ) ')
      .trim()
  }

  private getCurrentPage(): string {
    if (process.client) {
      return window.location.pathname
    }
    return 'server-side'
  }

  private getWalletAddress(): string | undefined {
    // Integrate with wallet store to get current wallet address
    if (process.client) {
      // Try to get from localStorage or sessionStorage
      const storedAddress = localStorage.getItem('xrp-wallet-address') || 
                           sessionStorage.getItem('xrp-wallet-address')
      if (storedAddress) {
        return storedAddress
      }
      
      // Try to get from window object if available
      if ((window as any).xrpWalletAddress) {
        return (window as any).xrpWalletAddress
      }
    }
    return undefined
  }

  private getNetworkInfo() {
    return {
      endpoint: process.env.BASE_GRAPHQL_SERVER_URL || 'http://127.0.0.1:8080/query',
      userAgent: process.client ? window.navigator.userAgent : 'server-side'
    }
  }

  private categorizeError(error: any): string {
    if (error.networkError) {
      if (error.networkError.statusCode === 401) return 'authentication'
      if (error.networkError.statusCode === 403) return 'authorization'
      if (error.networkError.statusCode >= 500) return 'server_error'
      return 'network'
    }
    if (error.graphQLErrors && error.graphQLErrors.length > 0) {
      const graphQLError = error.graphQLErrors[0]
      if (graphQLError.extensions?.code === 'UNAUTHENTICATED') return 'authentication'
      if (graphQLError.extensions?.code === 'FORBIDDEN') return 'authorization'
      if (graphQLError.extensions?.code === 'VALIDATION') return 'validation'
      return 'graphql'
    }
    return 'unknown'
  }

  private updateErrorStats(entry: GraphQLLogEntry) {
    this.errorStats.totalErrors++
    
    // Update errors by query
    if (!this.errorStats.errorsByQuery[entry.queryName]) {
      this.errorStats.errorsByQuery[entry.queryName] = 0
    }
    this.errorStats.errorsByQuery[entry.queryName]++

    // Update errors by type
    if (entry.error) {
      const errorType = this.categorizeError(entry.error)
      if (!this.errorStats.errorsByType[errorType]) {
        this.errorStats.errorsByType[errorType] = 0
      }
      this.errorStats.errorsByType[errorType]++
    }

    this.errorStats.lastError = entry
  }

  private updatePerformanceStats(entry: GraphQLLogEntry) {
    this.performanceStats.totalQueries++
    
    if (entry.duration) {
      // Update average response time
      const totalTime = this.performanceStats.averageResponseTime * (this.performanceStats.totalQueries - 1) + entry.duration
      this.performanceStats.averageResponseTime = totalTime / this.performanceStats.totalQueries

      // Update slowest queries
      this.performanceStats.slowestQueries.push({
        queryName: entry.queryName,
        duration: entry.duration
      })
      this.performanceStats.slowestQueries.sort((a, b) => b.duration - a.duration)
      this.performanceStats.slowestQueries = this.performanceStats.slowestQueries.slice(0, 10) // Keep top 10
    }

    // Update queries by endpoint
    const endpoint = entry.networkInfo?.endpoint || 'unknown'
    if (!this.performanceStats.queriesByEndpoint[endpoint]) {
      this.performanceStats.queriesByEndpoint[endpoint] = 0
    }
    this.performanceStats.queriesByEndpoint[endpoint]++
  }

  startQuery(queryName: string, operation: 'query' | 'mutation' | 'subscription' = 'query', variables?: any, queryString?: string): string {
    const id = `${queryName}_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    
    const entry: GraphQLLogEntry = {
      id,
      queryName,
      operation,
      queryString,
      variables,
      startTime: Date.now(),
      success: false,
      userContext: {
        page: this.getCurrentPage(),
        walletAddress: this.getWalletAddress(),
        timestamp: new Date().toISOString()
      },
      networkInfo: this.getNetworkInfo()
    }

    this.logs.push(entry)
    
    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
      console.log(`🚀 XRP GraphQL Query Started: ${queryName}`, {
        id,
        operation,
        queryString: queryString ? this.formatQueryForLogging(queryString) : 'Not provided',
        variables,
        page: entry.userContext?.page || 'unknown',
        timestamp: entry.userContext?.timestamp || new Date().toISOString()
      })
    }

    return id
  }

  completeQuery(id: string, success: boolean, error?: any) {
    const entry = this.logs.find(log => log.id === id)
    if (!entry) {
      console.warn(`XRP GraphQL Logger: Could not find log entry for id ${id}`)
      return
    }

    entry.endTime = Date.now()
    entry.duration = entry.endTime - entry.startTime
    entry.success = success
    entry.error = error

    if (success) {
      this.updatePerformanceStats(entry)
      
      // Log to console in development
      if (process.env.NODE_ENV === 'development') {
        console.log(`✅ XRP GraphQL Query Completed: ${entry.queryName}`, {
          duration: `${entry.duration}ms`,
          page: entry.userContext?.page
        })
      }
    } else {
      this.updateErrorStats(entry)
      this.updatePerformanceStats(entry)
      
      // Log error to console
      console.error(`❌ XRP GraphQL Query Failed: ${entry.queryName}`, {
        error: entry.error,
        duration: entry.duration ? `${entry.duration}ms` : 'unknown',
        page: entry.userContext?.page,
        variables: entry.variables,
        queryString: entry.queryString ? this.formatQueryForLogging(entry.queryString) : 'Not provided'
      })

      // Log to external service in production
      if (process.env.NODE_ENV === 'production') {
        this.logToExternalService(entry)
      }
    }
  }

  private logToExternalService(entry: GraphQLLogEntry) {
    // Implement external logging service (e.g., Sentry, LogRocket, etc.)
    // For now, we'll implement a basic external logging mechanism
    
    const logData = {
      timestamp: entry.userContext?.timestamp,
      queryName: entry.queryName,
      operation: entry.operation,
      error: entry.error,
      userContext: entry.userContext,
      networkInfo: entry.networkInfo,
      duration: entry.duration,
      variables: entry.variables,
      queryString: entry.queryString ? this.formatQueryForLogging(entry.queryString) : undefined
    }

    // Log to console with structured format for external services
    console.error('Production Error Log:', JSON.stringify(logData, null, 2))

    // Send to external service if configured
    if (process.env.EXTERNAL_LOGGING_ENDPOINT) {
      // In a real implementation, you would send this to your logging service
      // fetch(process.env.EXTERNAL_LOGGING_ENDPOINT, {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(logData)
      // }).catch(err => console.error('Failed to send to external logging service:', err))
    }

    // Store in localStorage for debugging (limited to last 50 errors)
    if (process.client) {
      try {
        const storedErrors = JSON.parse(localStorage.getItem('xrp-graphql-errors') || '[]')
        storedErrors.unshift(logData)
        storedErrors.splice(50) // Keep only last 50 errors
        localStorage.setItem('xrp-graphql-errors', JSON.stringify(storedErrors))
      } catch (err) {
        console.error('Failed to store error in localStorage:', err)
      }
    }
  }

  // Public getters for reactive data
  getLogs() {
    return this.logs
  }

  getErrorStats() {
    return this.errorStats
  }

  getPerformanceStats() {
    return this.performanceStats
  }

  // Utility methods
  getRecentErrors(limit: number = 10): GraphQLLogEntry[] {
    return this.logs
      .filter(log => !log.success)
      .sort((a, b) => (b.endTime || 0) - (a.endTime || 0))
      .slice(0, limit)
  }

  getSlowQueries(threshold: number = 1000): GraphQLLogEntry[] {
    return this.logs
      .filter(log => log.duration && log.duration > threshold)
      .sort((a, b) => (b.duration || 0) - (a.duration || 0))
  }

  clearLogs() {
    this.logs.splice(0, this.logs.length)
  }

  // Get stored errors from localStorage for debugging
  getStoredErrors(): any[] {
    if (process.client) {
      try {
        return JSON.parse(localStorage.getItem('xrp-graphql-errors') || '[]')
      } catch (err) {
        console.error('Failed to parse stored errors:', err)
        return []
      }
    }
    return []
  }

  // Clear stored errors from localStorage
  clearStoredErrors() {
    if (process.client) {
      localStorage.removeItem('xrp-graphql-errors')
    }
  }

  exportLogs(): string {
    return JSON.stringify({
      logs: this.logs,
      errorStats: this.errorStats,
      performanceStats: this.performanceStats,
      exportTime: new Date().toISOString()
    }, null, 2)
  }
}

// Create singleton instance
const xrpGraphQLLogger = new XRPGraphQLLogger()

export default function useXrpGraphQLLogging() {
  return {
    // Logger methods
    startQuery: xrpGraphQLLogger.startQuery.bind(xrpGraphQLLogger),
    completeQuery: xrpGraphQLLogger.completeQuery.bind(xrpGraphQLLogger),
    
    // Data getters
    logs: xrpGraphQLLogger.getLogs(),
    errorStats: xrpGraphQLLogger.getErrorStats(),
    performanceStats: xrpGraphQLLogger.getPerformanceStats(),
    
    // Utility methods
    getRecentErrors: xrpGraphQLLogger.getRecentErrors.bind(xrpGraphQLLogger),
    getSlowQueries: xrpGraphQLLogger.getSlowQueries.bind(xrpGraphQLLogger),
    clearLogs: xrpGraphQLLogger.clearLogs.bind(xrpGraphQLLogger),
    getStoredErrors: xrpGraphQLLogger.getStoredErrors.bind(xrpGraphQLLogger),
    clearStoredErrors: xrpGraphQLLogger.clearStoredErrors.bind(xrpGraphQLLogger),
    exportLogs: xrpGraphQLLogger.exportLogs.bind(xrpGraphQLLogger)
  }
} 