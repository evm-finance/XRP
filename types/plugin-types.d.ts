import Vue from 'vue'
import { Emitter } from 'mitt'
import { EmitterEvents } from '~/types/events'

interface Params {
  minDigits?: number
  maxDigits?: number
  pre?: string
  after?: string
  useSymbol?: boolean
}
declare module '@nuxt/types' {
  interface Context {
    $f(val: number, params: Params): string
    $copyAddressToClipboard(value: string): Promise<void>
    $truncateAddress(address: string, zeroIndexTo: number, endIndexMinus: number): string
    $imageUrlBySymbol(symbol: string | null): string
    $applyPtcChange(val: number): { value: string; color: string; icon: string | null }
    $emitter: Emitter<EmitterEvents>
  }
}
declare module 'vue/types/vue' {
  interface Vue {
    $f(val: number, params: Params): string
    $copyAddressToClipboard(value: string): Promise<void>
    $setAltImageUrl(event: any): string
    $imageUrlBySymbol(symbol: string | null): string
    $truncateAddress(address: string, zeroIndexTo: number, endIndexMinus: number): string
    $navigateToExplorer(address: string, type?: 'tx' | 'address', blockExplorerUrl?: string): string
    $applyPtcChange(val: number): { value: string; color: string; icon: string | null }
    $emitter: Emitter<EmitterEvents>
  }
}
