import { useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'

export const useHelpers = () => {
  const { state } = useStore<State>()

  const isNativeToken = (symbol: string) => {
    return state.configs.networks[0]?.symbol === symbol
  }

  const getCurrentNetwork = () => {
    return state.configs.networks[0]
  }

  function debounceAsync<T, Callback extends (...args: any[]) => Promise<T>>(
    callback: Callback,
    wait: number
  ): (...args: Parameters<Callback>) => Promise<T> {
    let timeoutId: number | null = null

    return (...args: any[]) => {
      if (timeoutId) {
        clearTimeout(timeoutId)
      }

      return new Promise<T>((resolve) => {
        const timeoutPromise = new Promise<void>((resolve) => {
          // @ts-ignore
          timeoutId = setTimeout(resolve, wait)
        })
        timeoutPromise.then(async () => {
          // eslint-disable-next-line n/no-callback-literal
          resolve(await callback(...args))
        })
      })
    }
  }

  return {
    isNativeToken,
    getCurrentNetwork,
    debounceAsync,
  }
}
