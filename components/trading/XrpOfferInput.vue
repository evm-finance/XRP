<template>
    <div>
        <v-btn @click="openDialog">Open Dialog</v-btn>
    </div>
</template>
  
  <script lang="ts">
  import { computed, defineComponent, PropType, ref, Ref, toRefs, watch } from '@nuxtjs/composition-api'
  import { EmitEvents } from '~/types/events'
  
  type Props = {
    value: number
  }
  
  export default defineComponent<Props>({
    props: {
        value: { type: Number, default: 0 },
    },
  
    setup(props, { emit }) {
      // STATE
      const amountVal = ref<number>(0)
      const menuOpenEvent = EmitEvents.onXrpInputOpen
        const isOpen = ref(('false'))

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

      function openDialog()
      {
        console.log("dialog opened")
      }
  
      function onMenuOpen() {
        emit(menuOpenEvent)
      }
  
      function onInput() {
        // const searchTimeout: any = null
        // clearTimeout(searchTimeout)
        // setTimeout(() => {
        emit(EmitEvents.onValueChanged, { value: amountVal.value })
        // }, 1500)
      }
  
  

      return { openDialog, isOpen }
    },
  })
  </script>
  