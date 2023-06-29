<template>
  <div>
    <v-card tile outlined height="1300">
      <v-skeleton-loader v-if="loadingRef" type="table-tbody,table-tbody,table-tbody" />
      <client-only>
        <v-data-table
          v-if="!loadingRef"
          hide-default-footer
          :headers="cols"
          :items="screenerDataFormatted"
          :items-per-page="25"
          class="elevation-0 row-height-50"
          mobile-breakpoint="0"
          :item-class="rowClass"
        >
          <template #item.blockNumber="{ item }">
            <div>
              <nuxt-link class="pink--text" :to="`/xrp-explorer/ledger/${item.blockNumber}`">
                {{ item.blockNumber }}
              </nuxt-link>
              <div class="grey--text text-caption">
                {{ item.hashShort }}
              </div>
            </div>
          </template>

          <template #item.events="{ item }">
            <span v-for="(i, d) in item.XRPLedger.eventsCount" :key="d">
              <v-chip class="mx-2" label :color="eventColor(d)" small outlined>{{ `${d}  ${i}` }}</v-chip>
            </span>
          </template>
        </v-data-table>
      </client-only>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, toRefs, useContext } from '@nuxtjs/composition-api'
import { Block } from '~/types/apollo/main/types'

type Props = {
  blocks: Block[]
  loading: boolean
  currentTime: number
}
export default defineComponent<Props>({
  props: {
    blocks: { type: Array as PropType<Block[]>, default: () => [] },
    loading: { type: Boolean, default: true },
    currentTime: { type: Number, default: 0 },
  },
  setup(props) {
    // COMPOSABLE
    const { $copyAddressToClipboard } = useContext()
    const loadingRef = toRefs(props).loading
    const currentTimeRef = toRefs(props).currentTime
    const screenerDataFormatted = computed(() =>
      props.blocks.map((elem) => ({
        ...elem,
        hashShort: `${elem.XRPLedger.ledgerHash.slice(0, 10)}........${elem.XRPLedger.ledgerHash.slice(
          elem.XRPLedger.ledgerHash.length - 10,
          elem.XRPLedger.ledgerHash.length
        )}`,
        secondsAgo: `${(currentTimeRef.value - elem.minedAt).toFixed(0)} Seconds ago`,
      }))
    )

    const eventColor = (event: any): string => {
      const colors: { [key: string]: string } = {
        OfferCreate: 'green',
        TicketCreate: 'green',
        Payment: 'green',
        NFTokenCreateOffer: 'green',
        NFTokenCancelOffer: 'red',
        NFTokenMint: 'primary',
        OfferCancel: 'red',
        TrustSet: 'orange',
      }

      // @ts-ignore
      return Object.hasOwn(colors, event) ? colors[event] : 'grey'
    }

    const rowClass = (item: any) => {
      if (item.updateOption) {
        return item.updateOption.status ? 'on-update-flash' : ''
      }
    }

    const cols = computed(() => {
      return [
        {
          text: 'Ledger',
          align: 'left',
          value: 'blockNumber',
          width: '300',
          class: ['px-4', 'text-truncate'],
          cellClass: ['px-4', 'text-truncate'],
          sortable: false,
        },

        {
          text: 'Σ(Tx)',
          align: 'left',
          value: 'txCount',
          width: '100',
          sortable: false,
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
        },

        {
          text: 'Time',
          align: 'left',
          value: 'secondsAgo',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate', 'grey--text'],
          width: '200',
          sortable: false,
        },

        {
          text: 'Events',
          align: 'left',
          value: 'events',
          class: ['px-2', 'text-truncate'],
          cellClass: ['px-2', 'text-truncate'],
          sortable: false,
        },
      ]
    })

    return {
      screenerDataFormatted,
      cols,
      loadingRef,
      copyAddressToClipboard: $copyAddressToClipboard,
      eventColor,
      rowClass,
    }
  },
})
</script>
<style lang="css">
.on-update-flash {
  background-color: #4caf5026;
}
</style>
