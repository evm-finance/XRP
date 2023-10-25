<template>
  <v-row>
    <v-col lg="4" md="12" cols="12">
      <v-card tile outlined elevation="0" class="pa-4" height="100%">
        <v-row no-gutters>
          <v-avatar size="30" class="mt-1">
            <v-img alt="Avatar" :src="protocolImage" :lazy-src="protocolImage" />
          </v-avatar>

          <h2 class="text-h4 font-weight-medium ml-3">
            {{ protocolName }}
          </h2>
          <v-btn-toggle
            v-model="currenVersion"
            mandatory
            rounded
            color="primary"
            background-color="transparent"
            class="pt-3 ml-4"
          >
            <v-btn v-for="(item, index) in versions" :key="index" :value="item" rounded height="24" color="transparent">
              <span class="text-body-2 text-capitalize">{{ item }}</span>
            </v-btn>
          </v-btn-toggle>
        </v-row>

        <v-row no-gutters class="mb-2">
          <client-only>
            <v-col cols="12" class="mt-2">
              <v-chip color="grey darken-4" label small>Protocol</v-chip>
              <v-chip color="grey darken-4" label small>DeFi</v-chip>
            </v-col>
          </client-only>
        </v-row>
      </v-card>
    </v-col>
    <v-col lg="8" md="12">
      <v-card tile outlined elevation="0" class="pa-4" height="100%">
        <v-row>
          <v-col>
            <v-row no-gutters>
              <v-col cols="10">
                <h1 class="headline" v-text="protocolPageTitle" />
              </v-col>
              <v-col class="text-right">
                <v-btn
                  width="20"
                  height="20"
                  class="mx-1 pa-0"
                  color="primary"
                  icon
                  target="_blank"
                  :href="`https://twitter.com/${protocolTwitter}`"
                >
                  <v-icon size="20">mdi-twitter</v-icon>
                </v-btn>

                <v-btn width="20" height="20" class="pa-0" color="primary" icon target="_blank" :href="protocolUrl">
                  <v-icon size="20">mdi-web</v-icon>
                </v-btn>
              </v-col>
            </v-row>
            <v-row no-gutters>
              <div
                class="text-subtitle-2 text-sm-subtitle-1 font-weight-regular grey--text"
                v-text="protocolDescription"
              />
            </v-row>
          </v-col>
        </v-row>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, computed, ref, watch } from '@nuxtjs/composition-api'
import { aaveVersion } from '~/composables/useAavePools'

type Props = {
  name: string
  address: string
  symbol: string
  url: string
  twitter: string
  title: string
  description: string
}
export default defineComponent<Props>({
  props: {
    name: { type: String, required: true },
    symbol: { type: String, required: true },
    url: { type: String, required: true },
    title: { type: String, required: true },
    description: { type: String, required: true },
    twitter: { type: String, required: true },
  },

  setup(props, { emit }) {
    const protocolImage = computed(
      () =>
        `https://quantifycrypto.s3-us-west-2.amazonaws.com/pictures/crypto-img/32/icon/${props.symbol.toLowerCase()}.png`
    )
    const currenVersion = ref<aaveVersion>('v3')
    const versions = ref<aaveVersion[]>(['v2', 'v3'])

    watch(currenVersion, (v) => emit('on-version-changed', v))

    return {
      protocolName: props.name,
      protocolSymbol: props.symbol,
      protocolUrl: props.url,
      protocolPageTitle: props.title,
      protocolDescription: props.description,
      protocolTwitter: props.twitter,
      protocolImage,
      currenVersion,
      versions,
    }
  },
})
</script>

<style scoped></style>
