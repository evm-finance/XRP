<template>
  <v-row>
    <v-col>
      <v-row>
        <v-col class="mb-1 caption" v-text="actionAllowance.text + ':'" />
        <v-col class="text-right caption">
          <span class="font-weight-medium" v-text="actionAllowance.value" />
        </v-col>
      </v-row>
      <v-form ref="form">
        <v-text-field
          v-model="amount"
          :rules="fromRules"
          outlined
          color="pink"
          clearable
          dense
          class="rounded-0 text-right"
          placeholder="Number of Tokens"
          :hint="amountInUsd"
          :disabled="disabled"
        >
          <template #prepend-inner>
            <v-row class="ml-0 my-0 mr-1">
              <div class="text-no-wrap">
                <v-avatar size="19"><v-img :alt="logo" :src="logo" /></v-avatar>
                <span class="mx-2 subtitle-2" :class="[textClass]" v-text="symbol" />
              </div>
              <v-divider class="mx-3" vertical></v-divider>
            </v-row>
          </template>
        </v-text-field>
      </v-form>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, ref, PropType, computed, useStore, watch, useContext } from '@nuxtjs/composition-api'
import { State } from '~/types/state'
import { EmitEvents } from '~/types/events'

type arrayOfOfFunction = ((v: string) => boolean | string)[]

type Props = {
  symbol: string
  priceUsd: number
  logo: string
  rules: arrayOfOfFunction
  allowance?: { text: string; value: string | number }
  disabled: boolean
}

export default defineComponent<Props>({
  props: {
    disabled: { type: Boolean, default: false },
    symbol: { type: String, default: '', required: true },
    priceUsd: { type: Number, default: 0, required: true },
    logo: { type: String, default: '', required: true },
    rules: { type: Array as PropType<arrayOfOfFunction>, default: () => [], required: true },
    allowance: {
      type: Object as PropType<{ text: string; value: string | number }>,
      required: false,
      default: () => ({ text: '', value: '' }),
    },
  },
  setup(props, { emit }) {
    // COMPUTED
    const { state } = useStore<State>()
    const { $f } = useContext()

    // STATE
    const amountVal = ref<number>(0)
    const form = ref<any>(null)
    const amount = computed({
      get() {
        return new Intl.NumberFormat('en').format(amountVal.value)
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
    // COMPUTED
    const textClass = computed(() => state.ui[state.ui.theme].innerCardLighten)
    const amountInUsd = computed(() => $f(amountVal.value * props.priceUsd, { minDigits: 2, pre: '$ ' }))
    const actionAllowance = computed(() => props.allowance)
    const fromRules = computed(() => props.rules)

    // WATCH
    watch(amountVal, () =>
      emit(EmitEvents.onValueChanged, { isFormValid: form.value?.validate(), amountVal: amountVal.value })
    )

    return {
      amount,
      amountInUsd,
      textClass,
      form,
      actionAllowance,
      fromRules,
    }
  },
})
</script>