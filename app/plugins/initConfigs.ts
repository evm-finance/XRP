import { Context, Plugin } from '@nuxt/types'
const globalHelper: Plugin = async (context: Context) => await context.store.dispatch('configs/initConfigs', context)
export default globalHelper
