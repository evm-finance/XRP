<template>
  <div>
    <v-row justify="center">
      <v-col md="12" lg="9">
        <v-row>
          <v-col cols="5">
            <v-card tile outlined height="320" class="pa-2">
              <v-skeleton-loader v-if="loading" type="card-heading, image" height="400" />
              <div v-else>
                <span class="text-subtitle-2">Events Composition last 5 min</span>
                <client-only>
                  <aave-composition-chart
                    :data="events"
                    :chart-height="280"
                    :labels-disabled="false"
                    :ticks-disabled="false"
                  />
                </client-only>
              </div>
            </v-card>
          </v-col>
          <v-col cols="7">
            <v-card tile outlined height="320">
              <client-only>
                <x-r-p-libe-ledgers-chart v-if="blocks.length" :chart-height="320" :data="blocks" />
              </client-only>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col md="12" lg="9"><x-r-p-grid :blocks="blocks" :loading="loading" :current-time="currentTime" /> </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'
import XRPLibeLedgersChart from '~/components/XRPLibeLedgersChart.vue'
import useXrpScrerener from '~/composables/useXrpScrerener'

export default defineComponent({
  components: { XRPLibeLedgersChart },
  setup() {
    const { blocks, loading, currentTime } = useXrpScrerener()

    const events = computed(() => {
      const sumObject: Record<string, number> = {}

      for (const obj of blocks.value) {
        for (const key in obj.XRPLedger.eventsCount) {
          // eslint-disable-next-line no-prototype-builtins
          if (sumObject.hasOwnProperty(key)) {
            sumObject[key] += obj.XRPLedger?.eventsCount[key]
          } else {
            sumObject[key] = obj.XRPLedger?.eventsCount[key]
          }
        }
      }

      return Object.entries(sumObject)
        .map(([key, value]) => ({ name: key, value }))
        .sort((a, b) => (a.name < b.name ? -1 : 1))
    })

    return { blocks, loading, currentTime, events }
  },
  head: {},
})
</script>

<style scoped></style>
