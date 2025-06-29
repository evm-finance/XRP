<template>
  <div class="text-center">
    <v-icon class="mt-12" size="48" :color="txData.color">{{ txData.icon }}</v-icon>
    <h5 class="text-h5 mt-4" v-text="txData.message" />
    <p v-if="isTxMined && message" class="grey--text" v-text="message" />

    <v-row style="margin-top: 90px">
      <v-col class="text-left mb-1">
        <small><a class="grey--text" href="#" @click="showLogs = !showLogs">Transaction Details</a></small>
      </v-col>
      <v-col class="text-right grey--text mb-1">
        <small v-if="isTxMined && txHash">
          <a class="grey--text" target="_blank" :href="txHash">View on Explorer</a>
          <v-icon color="grey" size="14">mdi-open-in-new</v-icon>
        </small>
      </v-col>
    </v-row>

    <v-card v-if="showLogs" elevation="0" class="overflow-auto text-left" :height="scrollHeight">
      <pre
        class="overflow-y-auto subtitle-2 font-weight-regular grey--text"
        v-text="JSON.stringify(txReceipt, null, 2)"
      />
    </v-card>

    <v-btn tile class="text-capitalize" block color="primary" @click="onClose"> OK, Close it </v-btn>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from '@nuxtjs/composition-api'

interface XrpTransactionReceipt {
  hash?: string
  account?: string
  destination?: string
  transactionType?: string
  amount?: number
  currency?: string
  fee?: number
  date?: string
  meta?: any
  [key: string]: any
}

type Props = {
  receipt: XrpTransactionReceipt | Error
  isTxMined: boolean
  successMessage?: string | null
}

export default defineComponent<Props>({
  props: {
    receipt: { type: [Object, Error] as any, required: true },
    scrollHeight: { type: String, default: '140' },
    isTxMined: { type: Boolean, default: false, required: true },
    successMessage: { type: String as PropType<string | null>, default: null, required: false },
  },
  setup(props, { emit }) {
    const showLogs = ref(false)

    const txHash = computed(() => {
      if (props.receipt && typeof props.receipt === 'object' && 'hash' in props.receipt) {
        return `https://livenet.xrpl.org/transactions/${props.receipt.hash}`
      }
      return null
    })

    const txData = computed(() => {
      let data: { message: string; icon: string; color: string }
      props.isTxMined
        ? (data = { message: 'Transaction Successful', icon: 'mdi-check-circle-outline', color: 'green' })
        : (data = { message: 'Something Went Wrong', icon: 'mdi-alert-circle-outline', color: 'red' })
      return data
    })

    const onClose = () => emit('close')

    return {
      // COMPUTED
      txReceipt: props.receipt,
      message: props.successMessage,
      txMined: props.isTxMined,
      showLogs,
      txData,
      txHash,

      // METHODS
      onClose,
    }
  },
})
</script> 