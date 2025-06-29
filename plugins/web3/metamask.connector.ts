// Minimal Metamask connector for XRP-only functionality
import { ConnectorInterface } from './connector'

export class MetamaskConnector implements ConnectorInterface {
  async connect(): Promise<void> {
    // XRP-only implementation - no Metamask needed
    throw new Error('Metamask not supported in XRP-only mode')
  }

  async disconnect(): Promise<void> {
    // XRP-only implementation
  }

  isConnected(): boolean {
    return false
  }

  getAccount(): string | null {
    return null
  }

  getNetwork(): any {
    return null
  }
} 