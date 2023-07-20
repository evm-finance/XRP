<template>
  <v-menu v-if="ethMainNetInfo" :close-on-content-click="false" nudge-width="600" nudge-top="-46" left max-width="600">
    <template #activator="{ on, attrs }">
      <div class="d-flex">
        <v-btn class="subtitle-2 px-2 mr-2 text-capitalize font-weight-regular" text tile v-bind="attrs" v-on="on">
          <div>
            <v-icon>mdi-gas-station</v-icon>
            <span class="ml-1"> {{ ethMainNetInfo.gas.proposeGasPrice }}</span>
            <v-icon small>mdi-chevron-down</v-icon>
          </div>
        </v-btn>
      </div>
    </template>

    <v-card outlined tile class="pa-2">
      <v-simple-table>
        <template #default>
          <thead>
            <tr>
              <th></th>
              <th></th>
              <th>Low</th>
              <th>Average</th>
              <th>Fast</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="(elem, i) in gasData" :key="i">
              <td style="max-width: 22px; padding-right: 0">
                <v-img :src="$imageUrlBySymbol(elem.symbol)" :lazy-src="$imageUrlBySymbol(elem.symbol)" width="22" />
              </td>
              <td class="pt-1 text-no-wrap">
                <a
                  class="text-subtitle-2 text-decoration-none"
                  :href="elem.blockExplorer"
                  target="_blank"
                  v-text="elem.name"
                />
                <div class="mt-n1">
                  <small>
                    <a
                      target="_blank"
                      class="text-decoration-none grey--text"
                      :href="`${elem.blockExplorer}block/${elem.gas.lastBlock}/`"
                      v-text="elem.gas.lastBlock"
                    />
                  </small>
                </div>
              </td>
              <td class="green--text">
                {{ elem.gas.safeGasPrice === 0 ? '--' : `${$f(elem.gas.safeGasPrice, { minDigits: 0 })} gwei` }}
              </td>
              <td class="primary--text text--lighten-1">
                {{ elem.gas.proposeGasPrice === 0 ? '--' : `${$f(elem.gas.proposeGasPrice, { minDigits: 0 })} gwei` }}
              </td>
              <td class="red--text">
                {{ elem.gas.fastGasPrice === 0 ? '--' : `${$f(elem.gas.fastGasPrice, { minDigits: 0 })} gwei` }}
              </td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
  </v-menu>
</template>
<script lang="ts">
import { computed, defineComponent, Ref } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { GasGQL } from '~/apollo/main/config.query.graphql'
import { GasStats } from '~/types/apollo/main/types'

export default defineComponent({
  setup() {
    // COMPOSABLE
    const { result } = useQuery(GasGQL, null, { fetchPolicy: 'no-cache' })
    const gasData = computed(() => result.value?.gas ?? []) as Ref<GasStats[]>
    const ethMainNetInfo = computed(() => gasData.value.find((elem) => elem.symbol === 'ETH'))

    return {
      // COMPUTED
      gasData,
      ethMainNetInfo,
    }
  },
})
</script>
