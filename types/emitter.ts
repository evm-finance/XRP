import mitt from 'mitt'
import { Block } from '~/types/apollo/main/types'

type Events = {
  priceStream: string
  onInitGlobalSearch: string
  onNewBlock: Block[]
}
const emitter = mitt<Events>()
export default emitter
