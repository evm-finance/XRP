<template>
    <v-row>
      <v-col>
        <v-form ref="form">
          <v-text-field
            v-model="amount"
            hide-details
            outlined
            color="pink"
            dense
            class="rounded-0 px-0"
            :disabled="loading"
            @input="onInput"
          >
            <template #append>
              <v-divider class="mx-3 mb-2" vertical />
              <img
                width="24"
                height="24"
                style="margin-top: 2px"
                :src="$imageUrlBySymbol(token.symbol)"
                alt=""
                @error="$setAltImageUrl"
              />
              <v-btn text class="mb-2 py-0" tile small :disabled="loading" @click="$emit(menuOpenEvent)">
                <span class="text-subtitle-2" v-text="token.symbol" />
                <v-icon right dark> mdi-chevron-down</v-icon>
              </v-btn>
            </template>
          </v-text-field>
  
          <v-row class="mt-2 grey--text text--lighten-1 caption" no-gutters>
            <v-col cols="7">
              <v-skeleton-loader v-if="loading" type="heading" tile height="20" width="250" />
              <span v-else class="mr-0" v-text="`${fiatPrice ? $f(fiatPrice, { pre: '$ ', maxDigits: 2 }) : ''}`" />
            </v-col>
            <v-col class="text-right text-truncate" style="max-width: 100%">
              <v-skeleton-loader v-if="loading" type="heading" tile height="20" class="ml-8" width="300" />
              <span v-else>
                <span class="mr-0" v-text="`Balance: ${$f(balance, { maxDigits: 6 })}`" />
                <v-btn
                  v-if="isMaximumButtonVisible"
                  tile
                  color="primary"
                  class="pa-0 mb-1"
                  x-small
                  text
                  depressed
                  @click="setMax"
                  >MAX</v-btn
                >
              </span>
            </v-col>
          </v-row>
        </v-form>
      </v-col>
    </v-row>
  </template>
  
  <script lang="ts">
  import { computed, defineComponent, PropType, ref, toRefs, watch } from '@nuxtjs/composition-api'
  import { TradeType } from '@uniswap/sdk'
  import { UniswapToken } from '~/types/apollo/main/types'
  import { EmitEvents } from '~/types/events'
  
  type Props = {
    token: UniswapToken
    formTradeDirection: keyof typeof TradeType
    tradeDirection: keyof typeof TradeType
    balance: string
    expectedConvertQuote: number
    loading: boolean
    fiatPrice: number | null
  }
  
  export default defineComponent<Props>({
    props: {
      token: { type: Object as PropType<UniswapToken | null>, required: true },
      formTradeDirection: { type: String as PropType<keyof typeof TradeType>, required: true },
      tradeDirection: { type: String as PropType<keyof typeof TradeType>, required: true },
      balance: { type: Number, default: 0 },
      expectedConvertQuote: { type: Number, required: true },
      fiatPrice: { type: Number as PropType<null | number>, default: null },
      loading: { type: Boolean, default: false },
    },
  
    setup(props, { emit }) {
      // STATE
      const amountVal = ref<number>(0)
      const menuOpenEvent = EmitEvents.onUniswapTokenMenuOpen
      const isActive = computed(() => props.formTradeDirection === props.tradeDirection)
  
      const amount = computed({
        get() {
          return new Intl.NumberFormat('en', { minimumFractionDigits: 0, maximumFractionDigits: 20 }).format(
            amountVal.value
          )
        },
        set(value: string) {
          if (value && value.length > 0) {
            const num = parseFloat(value.replace(/,/g, ''))
            if (!isNaN(num)) {
              amountVal.value = parseFloat(value.replace(/,/g, ''))
            } else amountVal.value = 0
          } else amountVal.value = 0
        },
      })
  
      const isMaximumButtonVisible = computed(
        () => Number(props.balance) > 0 && Number(props.balance) !== amountVal.value
      )
      function onMenuOpen() {
        emit(menuOpenEvent)
      }
  
      function onInput() {
        // const searchTimeout: any = null
        // clearTimeout(searchTimeout)
        // setTimeout(() => {
        emit(EmitEvents.onValueChanged, { direction: props.formTradeDirection, value: amountVal.value })
        // }, 1500)
      }
  
      function setMax() {
        amountVal.value = Number(props.balance)
        emit(EmitEvents.onValueChanged, { direction: props.formTradeDirection, value: amountVal.value })
      }
  
      // Watch  of expectedConvertQuote for both directions, if field isn't active set incoming val in the form
      watch(toRefs(props).expectedConvertQuote, (val: number) => (!isActive.value ? (amount.value = `${val}`) : null))
  
      return { amount, currency, issuer }
    },
  })
  </script>
  