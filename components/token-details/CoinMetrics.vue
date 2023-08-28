<template>
  <v-card tile outlined class="pa-4" min-height="100">
    <v-row>
      <v-col cols="12" md="6">
        <v-row no-gutters>
          <v-col>
            <h2 :class="['text-h4 font-weight-medium', 'ml-1', 'text-no-wrap', priceTextClass]">
              {{ $f(tokenPriceUsd, { pre: '$ ', minDigits: 2, maxDigits: 4 }) }}
              <span class="grey--text text-caption">USD</span>
              <span class="ml-3">
                <v-btn
                  dark
                  depressed
                  tile
                  height="22"
                  :color="tokenPrice24h > 0 ? 'green' : 'red'"
                  class="px-2 subtitle-2"
                >
                  {{ $f(tokenPrice24h * 100, { minDigits: 2, after: ' %' }) }}
                </v-btn>
              </span>
            </h2>
          </v-col>

          <v-col cols="12">
            <div class="font-weight-medium subtitle-1 ml-2 text-no-wrap">
              <!--              <span v-text="$f(tokenPriceEth, { minDigits: 2, maxDigits: 8 })" />-->
              <span v-text="'-'" />
              <span class="font-weight-regular ml-1 grey--text text-caption">ETH</span>
            </div>
          </v-col>
        </v-row>
      </v-col>

      <v-col cols="12" md="6">
        <v-row>
          <v-col v-for="(item, i) in metrics" :key="i" cols="6" class="pa-4">
            <v-row no-gutters class="text-no-wrap">
              <v-col class="grey--text subtitle-2">
                <v-row no-gutters>
                  <v-avatar size="24" class="mr-3">
                    <v-icon color="primary" size="24">{{ item.icon }}</v-icon>
                  </v-avatar>
                  <div>
                    {{ item.title }}
                  </div>
                </v-row>
                <div class="pl-9 text-h6 text-truncate white--text">
                  {{ item.value }}
                </div>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-card>
</template>

<script lang="ts">
import { computed, defineComponent, ref, toRefs, useContext, watch } from '@nuxtjs/composition-api'

type Props = {
  priceUsd: number
  priceEth: number
  price24h: number
  marketCap: number
  circSupply: number
}
export default defineComponent<Props>({
  name: 'CoinMetrics',

  props: {
    priceUsd: { type: Number, required: true },
    priceEth: { type: Number, required: true },
    price24h: { type: Number, required: true },
    marketCap: { type: Number, default: 0 },
    circSupply: { type: Number, default: 0 },
  },

  setup(props) {
    // COMPOSABELES
    const { $f } = useContext()

    // STATE
    const tokenPriceUsd = toRefs(props).priceUsd
    const tokenPriceEth = toRefs(props).priceEth
    const tokenPrice24h = toRefs(props).price24h
    const priceTextClass = ref('')

    const metrics = computed(() => [
      {
        title: 'Market Cap',
        value: $f(props.marketCap, { useSymbol: true, pre: '$' }),
        icon: 'mdi-chart-line-variant',
      },
      {
        title: 'Circ. Supply',
        value: $f(props.circSupply, { useSymbol: true, pre: '$' }),
        icon: 'mdi-web',
      },
    ])

    // WATCHERS
    watch(tokenPriceUsd, (newPrice, prevPrice) =>
      newPrice > prevPrice ? (priceTextClass.value = 'green--text') : (priceTextClass.value = 'red--text')
    )

    return {
      tokenPriceUsd,
      tokenPriceEth,
      tokenPrice24h,
      priceTextClass,

      metrics,
    }
  },
})
</script>
