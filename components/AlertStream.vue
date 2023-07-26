<template>
  <v-card tile outlined height="1400" class="overflow-y-auto">
    <client-only>
      <div
        v-for="(item, i) in blocksFormatted"
        :key="i"
        :style="{ 'background-color': item.updateOption ? item.updateOption.color : '' }"
      >
        <v-list dense two-line class="py-0">
          <v-list-item-group color="primary">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>
                  <span class="grey--text text-caption">Block #</span>
                  <span class="primary--text">{{ item.blockNumber }}</span>
                </v-list-item-title>
                <v-list-item-subtitle>
                  <span class="grey--text text-caption">Mined:</span>{{ item.minedAt }}
                </v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-icon>
                <v-icon size="18" color="primary">mdi-open-in-new</v-icon>
              </v-list-item-icon>
            </v-list-item>
          </v-list-item-group>
        </v-list>
        <v-row no-gutters>
          <v-col v-for="(stat, statIndex) in item.stats" :key="statIndex" cols="3" class="pa-3 text-center">
            <v-row>
              <v-col cols="12">
                <div>
                  <span class="text-h4">{{ stat.value }}</span>
                </div>
                <div class="caption font-weight-medium grey--text">{{ stat.text }}</div>
              </v-col>
              <v-divider v-if="i < 3" style="" vertical />
            </v-row>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-carousel height="200" hide-delimiter-background :show-arrows="false" delimiter-icon="mdi-dot">
              <v-carousel-item v-for="(slide, slideIndex) in item.sliced" :key="slideIndex">
                <v-simple-table dense>
                  <template #default>
                    <tbody>
                      <tr v-for="(token, tokenIndex) in slide" :key="tokenIndex">
                        <td class="py-1">
                          <v-avatar size="20">
                            <img :src="$imageUrlBySymbol(token.token0Symbol)" alt="" @error="$setAltImageUrl" />
                          </v-avatar>
                          <v-avatar size="20">
                            <img :src="$imageUrlBySymbol(token.token1Symbol)" alt="" @error="$setAltImageUrl" />
                          </v-avatar>
                          <span class="ml-3">
                            {{ token.token0Symbol }}<span class="grey--text">/{{ token.token1Symbol }}</span></span
                          >
                        </td>
                        <td>
                          <span :class="`${token.changeApplied.color}--text`">
                            {{ token.changeApplied.value }}
                            <v-icon v-if="token.changeApplied.icon" :color="token.changeApplied.color" size="14"
                              >{{ token.changeApplied.icon }}
                            </v-icon>
                          </span>
                        </td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </v-carousel-item>
            </v-carousel>
          </v-col>
        </v-row>
        <v-divider />
      </div>
    </client-only>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRefs, useContext, watch } from '@nuxtjs/composition-api'
import { useQuery, useSubscription } from '@vue/apollo-composable/dist'
import { BlocksGQL, BlocksStreamGQL } from '~/apollo/main/token.query.graphql'
import { Block, BlockMetric } from '~/types/apollo/main/types'

type BlockMetricObserver = BlockMetric & {
  changeApplied: { value: string; color: string; icon: string | null }
}

type BlockObserver = Block & {
  stats: { text: string; value: number }[]
  sliced: BlockMetricObserver[][]
  updateOption?: { status: boolean; color: string | null }
}
type Props = {
  networkId: string
}
export default defineComponent<Props>({
  props: {
    networkId: { type: String as PropType<string>, default: 'ethereum' },
  },
  setup(props) {
    const { $applyPtcChange } = useContext()
    const blocks = ref<BlockObserver[]>([])
    let updateTimeout: any = null
    const network = toRefs(props).networkId

    const { onResult } = useQuery(BlocksGQL, () => ({ network: network.value }), { fetchPolicy: 'no-cache' })
    const { result: liveBlock } = useSubscription(BlocksStreamGQL, () => ({ network: network.value }), {
      fetchPolicy: 'no-cache',
    })
    onResult(({ data }) => {
      blocks.value = data?.blocks ?? []
    })

    const blocksFormatted = computed<BlockObserver[]>(() => {
      if (blocks.value.length === 0) {
        return []
      }

      blocks.value.forEach((elem) => elem.metrics.items.sort((a, b) => b.totalLiquidity - a.totalLiquidity))

      return blocks.value.map((elem) => ({
        ...elem,
        stats: [
          { text: 'All trans.', value: elem.txCount },
          { text: 'Swaps', value: elem.swapCount },
          { text: 'New Pairs', value: elem.pairCreatedCount },
          { text: 'Minted', value: elem.mintCount },
        ],
        sliced: elem.metrics.items.reduce((result: any, object: any, index) => {
          const arrayIndex = Math.floor(index / 4)
          if (!result[arrayIndex]) {
            result[arrayIndex] = []
          }
          object.changeApplied = $applyPtcChange(object.change1H)
          result[arrayIndex].push(object)
          return result
        }, []),
      }))
    })

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    function addNewRecords(newRecords: BlockObserver[]) {
      clearTimeout(updateTimeout)
      newRecords = newRecords.map((elem) => ({
        ...elem,
        updateOption: { status: true, color: '#4caf5026' },
      }))
      blocks.value = [...newRecords, ...blocks.value]
      if (blocks.value.length > 4) {
        blocks.value.splice(-newRecords.length)
      }

      updateTimeout = setTimeout(() => {
        blocks.value.forEach((elem) => {
          if (elem.updateOption) {
            elem.updateOption = { status: false, color: null }
          }
        })
      }, 1000)
    }

    watch(liveBlock, (val: any) => {
      const newData: BlockObserver[] = val?.block ?? []
      addNewRecords(newData)
    })

    return {
      blocks,
      blocksFormatted,
    }
  },
})
</script>
<style lang="css">
.v-carousel__controls__item.v-btn.v-btn--icon {
  background-color: #ebece9; /* Background color of non-active ones */
  height: 8px; /* Height you want */
  width: 8px; /* Width you want */
  border-radius: 100%; /* Remove default border radius */
}

.v-carousel__controls__item.v-btn.v-btn--icon.v-btn--active {
  background-color: #e91e63; /* Colour for active one */
}

.v-carousel__controls__item.v-btn.v-btn--icon:hover {
  background-color: #e91e63; /* You might also want to customise the hover effect */
}
/*.v-btn__content .v-icon {
  display: none; !* Removes the default icon *!
}*/
</style>
