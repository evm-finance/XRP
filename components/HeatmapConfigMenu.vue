<template>
  <v-menu
    v-model="settingsMenu"
    :close-on-content-click="false"
    :nudge-width="100"
    offset-x
    max-width="300"
    min-width="300"
  >
    <template #activator="{ on, attrs }">
      <v-btn :fab="!icon" :icon="icon" height="32" width="32" v-bind="attrs" v-on="on">
        <v-icon class="pa-0 ma-0">{{ icons.mdiCog }}</v-icon>
      </v-btn>
    </template>

    <v-card tile outlined class="pa-2">
      <v-row no-gutters class="px-3 py-1">
        <v-col cols="12">
          <span class="text-subtitle-2">Settings</span>
        </v-col>
        <v-col cols="12">
          <span class="text-caption grey--text lighten-2"> Customize your Heatmap experience. </span></v-col
        >
      </v-row>
      <v-divider class="my-1" />
      <v-row no-gutters class="justify-center align-center">
        <v-col cols="10" class="px-3 text-no-wrap"><span class="text-body-2">Blue Tiles</span> </v-col>
        <v-col>
          <div class="text-right">
            <v-checkbox v-model="blueTile" hide-details class="ma-0 py-1" s />
          </div>
        </v-col>
      </v-row>
      <v-row no-gutters class="justify-center align-center">
        <v-col cols="10" class="px-3 text-no-wrap"><span class="text-body-2">Gainers And Losers</span> </v-col>
        <v-col>
          <div class="text-right">
            <v-checkbox v-model="displayGainersAndLosers" hide-details class="ma-0 py-1" />
          </div>
        </v-col>
      </v-row>

      <v-row no-gutters class="justify-center align-center">
        <v-col cols="6" class="px-3 text-no-wrap"><span class="text-body-2">Num Of Coins</span> </v-col>
        <v-col cols="6">
          <div class="text-right">
            <v-menu offset-y right close-on-content-click>
              <template #activator="{ on }">
                <v-btn tile depressed color="transparent" class="pl-0" v-on="on">
                  <span class="text-subtitle-2 text-capitalize ml-2">
                    {{ numOfCoins }}
                  </span>
                  <v-icon right dark>{{ icons.mdiChevronDown }}</v-icon>
                </v-btn>
              </template>

              <v-list dense col="3" class="notranslate" color="#121212">
                <v-list-item-group v-model="numOfCoins" active-class="" color="primary" mandatory>
                  <v-list-item v-for="item in numOfCoinsOptions" :key="item" :value="item" class="my-0 py-0">
                    <v-list-item-content class="my-0 py-0">
                      <v-list-item-title v-text="item" />
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-menu>
          </div>
        </v-col>
      </v-row>
      <v-row no-gutters class="justify-center align-center">
        <v-col cols="4" class="px-3 text-no-wrap"><span class="text-body-2">Block Size</span> </v-col>
        <v-col cols="8">
          <div class="text-right">
            <v-menu offset-y right close-on-content-click>
              <template #activator="{ on }">
                <v-btn tile depressed color="transparent" class="pl-0" v-on="on">
                  <span class="text-subtitle-2 text-capitalize ml-2">
                    {{ blockSizeName.text }}
                  </span>
                  <v-icon right dark>{{ icons.mdiChevronDown }}</v-icon>
                </v-btn>
              </template>

              <v-list dense col="3" class="notranslate" color="#121212">
                <v-list-item-group v-model="blockSize" active-class="" color="primary" mandatory>
                  <v-list-item v-for="item in blockSizeOptions" :key="item.value" class="my-0 py-0" :value="item.value">
                    <v-list-item-content class="my-0 py-0">
                      <v-list-item-title v-text="item.text" />
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-menu>
          </div>
        </v-col>
      </v-row>

      <v-row no-gutters class="justify-center align-center">
        <v-col cols="4" class="px-3 text-no-wrap"><span class="text-body-2">Performance</span> </v-col>
        <v-col cols="8">
          <div class="text-right">
            <v-menu offset-y right close-on-content-click>
              <template #activator="{ on }">
                <v-btn tile depressed color="transparent" class="pl-0" v-on="on">
                  <span class="text-subtitle-2 text-capitalize ml-2">
                    {{ timeFrame }}
                  </span>
                  <v-icon right dark>{{ icons.mdiChevronDown }}</v-icon>
                </v-btn>
              </template>

              <v-list dense col="3" class="notranslate" color="#121212">
                <v-list-item-group v-model="timeFrame" active-class="" color="primary" mandatory>
                  <v-list-item v-for="item in timeFrameOptions" :key="item.value" class="my-0 py-0" :value="item.value">
                    <v-list-item-content class="my-0 py-0">
                      <v-list-item-title v-text="item.text" />
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-menu>
          </div>
        </v-col>
      </v-row>
    </v-card>
  </v-menu>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'
import { mdiCog, mdiChevronDown } from '@mdi/js'
import useHeatmapConfigs from '~/composables/heatmap/useHeatmapConfigs'

export default defineComponent({
  props: { icon: { type: Boolean, default: false } },
  setup() {
    const icons = { mdiCog, mdiChevronDown }
    const settingsMenu = ref(false)

    const {
      timeFrame,
      timeFrameOptions,
      blockSize,
      blockSizeOptions,
      blockSizeName,
      blueTile,
      displayFavorites,
      numOfCoins,
      numOfCoinsOptions,
      displayGainersAndLosers,
    } = useHeatmapConfigs()

    return {
      icons,
      settingsMenu,
      timeFrame,
      timeFrameOptions,
      blockSize,
      blockSizeOptions,
      blockSizeName,
      blueTile,
      displayFavorites,
      numOfCoins,
      numOfCoinsOptions,
      displayGainersAndLosers,
    }
  },
})
</script>
