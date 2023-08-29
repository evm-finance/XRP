<template>
  <v-row no-gutters justify="center">
    <v-col cols="12" md="10">
      <v-row justify="center">
        <v-col cols="12">
          <h1 class="text-h4">XRP Screener</h1>
        </v-col>
        <v-btn
        @click="connectWallet()">
        </v-btn>
      </v-row>

      <v-row justify="center">
        <v-col md="12">
          <v-card tile outlined height="1330">
            <v-skeleton-loader v-if="loading" type="table-tbody,table-tbody,table-tbody" />
            <client-only>
              <v-data-table
                v-if="!loading"
                hide-default-footer
                :headers="cols"
                :items="screenerDataFormatted"
                :items-per-page="25"
                class="elevation-0 row-height-50"
                mobile-breakpoint="0"
              >
                <template #item.currency="{ item }">
                  <div class="my-1">
                    <v-row no-gutters align="center">
                      <v-col cols="2" class="mr-3">
                        <v-avatar size="24" class="ml-2">
                          <img :src="item.icon" alt="" @error="$setAltImageUrl" />
                        </v-avatar>
                      </v-col>
                      <v-col>
                        <v-row no-gutters>
                          <v-col>
                            <nuxt-link class="text-capitalize font-weight-bold pink--text text-decoration-none" to="#">
                              {{ item.tokenName }}</nuxt-link
                            >
                          </v-col>
                        </v-row>
                        <v-row no-gutters>
                          <v-col>
                            <span class="grey--text text-caption">{{ item.currencyShort }}</span>
                          </v-col>
                        </v-row>
                      </v-col>
                    </v-row>
                  </div>
                </template>

                <!--                <template #item.blockNumber="{ item }">-->
                <!--                  <div>-->
                <!--                    <nuxt-link class="pink&#45;&#45;text" :to="`/xrp-explorer/ledger/${item.blockNumber}`">-->
                <!--                      {{ item.blockNumber }}-->
                <!--                    </nuxt-link>-->
                <!--                    <div class="grey&#45;&#45;text text-caption">-->
                <!--                      {{ item.hashShort }}-->
                <!--                    </div>-->
                <!--                  </div>-->
                <!--                </template>-->

                <!--                <template #item.events="{ item }">-->
                <!--                  <span v-for="(i, d) in item.XRPLedger.eventsCount" :key="d">-->
                <!--                    <v-chip class="mx-2" label :color="eventColor(d)" small outlined>{{ `${d}  ${i}` }}</v-chip>-->
                <!--                  </span>-->
                <!--                </template>-->

                <template #item.buy="{ item }">
          <v-btn
            text:
            outlined:
            color="green"
            class="pa-1 ma-1"
            height="26"
            @click="buy()"
          >
            <span class="text-caption">{{ item.value }}</span>
          </v-btn>
        </template>

        <template #item.sell="{ item }">
          <v-btn
            text:
            outlined:
            color="pink"
            class="pa-1 ma-1"
            height="26"
            @click="sell()"
          >
            <span class="text-caption">{{ item.value }}</span>
          </v-btn>
        </template>
              </v-data-table>
            </client-only>
          </v-card>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPScreenerGQL } from '~/apollo/queries'
import  useXrpTrade from '~/composables/useXrpTrade'
import { isInstalled, getAddress } from "@gemwallet/api";
interface XRPScreenerElem {
  currency: string
  issuerAddress: string
  icon: string
  tokenName: string
  issuerName: string
  marketcap: number
  price: number
  volume24H: number
}

export default defineComponent({
  setup() {
    const { $f } = useContext()
    const loading = ref(true)
    const screenerRawData = ref<XRPScreenerElem[]>([])
    const offers = ref<offerTypes[]>([])
    const { onResult } = useQuery(XRPScreenerGQL, { fetchPolicy: 'no-cache', pollInterval: 60000 })
    const {buy , sell, connectWallet } = useXrpTrade()
    type offerTypes = 'buy' | 'sell'

    const screenerDataFormatted = computed(() =>
      screenerRawData.value.map((elem) => ({
        ...elem,
        currencyShort: elem.currency.length > 20 ? elem.currency.substring(0, 20) + '...' : elem.currency,
        issuerAddressShort: `${elem.issuerAddress.slice(0, 10)}........${elem.issuerAddress.slice(
          elem.issuerAddress.length - 10,
          elem.issuerAddress.length
        )}`,
        priceFormatted: $f(elem.price, { minDigits: 6, after: '' }),
        marketCapFormatted: $f(elem.marketcap, { minDigits: 2, after: '' }),
        volume24HFormatted: $f(elem.volume24H, { minDigits: 2, after: '' }),
      }))
    )

    // EVENTS
    onResult((queryResult: any) => {
      screenerRawData.value = queryResult.data?.xrpScreener ?? []
      loading.value = false
    })

    // const screenerData = computed(() => result.value?.xrpScreener ?? [])

    const cols = computed(() => {
      return [
        {
          text: 'Currency',
          align: 'left',
          value: 'currency',
          width: '150',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate'],
          sortable: true,
        },
        {
          text: '',
          value: 'buy',
          width: 50,
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: '',
          value: 'sell',
          width: 50,
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
        {
          text: 'Issuer Name',
          align: 'left',
          value: 'issuerName',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Issuer Address',
          align: 'left',
          value: 'issuerAddressShort',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
        },

        {
          text: 'Price (XRP)',
          align: 'left',
          value: 'priceFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'MarketCap (XRP)',
          align: 'left',
          value: 'marketCapFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Volume 24H (XRP)',
          align: 'left',
          value: 'volume24HFormatted',
          width: '100',
          sortable: true,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },
      ]
    })

    return {
      loading,
      cols,
      screenerDataFormatted,
      buy, 
      sell,
      connectWallet
    }
  },
  head: {},
})
</script>
