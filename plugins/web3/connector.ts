// Minimal connector interface for XRP-only functionality
export interface ConnectorInterface {
  connect(): Promise<void>
  disconnect(): Promise<void>
  isConnected(): boolean
  getAccount(): string | null
  getNetwork(): any
}

export interface Web3ErrorInterface {
  code: number
  message: string
  data?: any
} 