import { useStore } from '@nuxtjs/composition-api'
import { State } from '~/types/state'

export function useHelpers() {
  const { state } = useStore<State>()
  const isNativeToken = (chainId: number, symbol: string) => {
    const chainIdAdjusted: number = chainId === 1337 ? 1 : chainId
    const chain = state.configs.chains.find(
      (elem) => elem.chainIdentifier === chainIdAdjusted && elem.symbol === symbol
    )
    return !!chain
  }

  // function debounce(callback: any, delay = 300) {
  //   let timeout: any = null
  //   return (...args: any) => {
  //     clearTimeout(timeout)
  //     timeout = setTimeout(() => callback.apply(args), delay)
  //   }
  // }

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

  return { isNativeToken, debounceAsync }
}
