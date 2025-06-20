<template>
  <div>
    <v-row no-gutters justify="center">
      <v-col cols="12" md="10">
        <v-row justify="center">
          <v-col cols="12">
            <h1 class="text-h4">XRP Screener</h1>
          </v-col>
          <!--          <v-btn @click="connectWallet()"> </v-btn>-->
        </v-row>
        
        <!-- XRP Balance Widget -->
        <v-row justify="center" class="mb-4">
          <v-col cols="12" lg="6">
            <xrp-balance-widget />
          </v-col>
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
                  <template #[`item.currency`]="{ item }">
                    <div class="text-no-wrap overflow-x-hidden">
                      <v-avatar size="20" class="mr-2">
                        <v-img
                          :src="$imageUrlBySymbol(item.currency.toLowerCase())"
                          :lazy-src="$imageUrlBySymbol(item.currency.toLowerCase())"
                        />
                      </v-avatar>
                      <nuxt-link
                        class="text-capitalize text-decoration-none white--text"
                        :to="`/token/${item.currency}?issuer=${item.issuerAddress}`"
                      >
                        {{ item.currency }}
                      </nuxt-link>
                    </div>
                  </template>

                  <template #[`item.issuerAddressShort`]="{ item }">
                    <div class="d-flex align-center">
                      <span class="font-family-mono">{{ item.issuerAddressShort }}</span>
                      <v-btn
                        icon
                        x-small
                        class="ml-1"
                        @click="copyToClipboard(item.issuerAddress)"
                      >
                        <v-icon size="16">mdi-content-copy</v-icon>
                      </v-btn>
                    </div>
                  </template>
                </v-data-table>
              </client-only>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, useContext } from '@nuxtjs/composition-api'
import { useQuery } from '@vue/apollo-composable/dist'
import { XRPScreenerGQL } from '~/apollo/queries'
import XrpBalanceWidget from '~/components/portfolio/XrpBalanceWidget.vue'
// import  useXrpTrade from '~/composables/useXrpTrade'

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
  components: {
    XrpBalanceWidget,
  },
  setup() {
    const { $f } = useContext()
    const loading = ref(true)
    const screenerRawData = ref<XRPScreenerElem[]>([])
    // const offers = ref<offerTypes[]>([])
    const { onResult } = useQuery(XRPScreenerGQL, { fetchPolicy: 'no-cache', pollInterval: 60000 })
    // const { buy, sell, connectWallet, isOpen, openDialog, closeDialog } = useXrpTrade()
    // type offerTypes = 'buy' | 'sell'

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

    const copyToClipboard = async (text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        // You could add a toast notification here
      } catch (err) {
        console.error('Failed to copy text: ', err)
      }
    }

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
        // {
        //   text: '',
        //   value: 'buy',
        //   width: 50,
        //   sortable: false,
        //   class: ['px-2', 'text-truncate'],
        //   cellClass: ['px-2', 'text-truncate'],
        // },
        // {
        //   text: '',
        //   value: 'sell',
        //   width: 50,
        //   sortable: false,
        //   class: ['px-2', 'text-truncate'],
        //   cellClass: ['px-2', 'text-truncate'],
        // },
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
      copyToClipboard,
      // buy,
      // sell,
      // connectWallet,
      // isOpen,
      // openDialog,
      // closeDialog,
    }
  },
  head: {},
})
</script>

<style scoped>
.font-family-mono {
  font-family: 'Courier New', monospace;
}
</style>
