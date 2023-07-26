import mitt from 'mitt'
import { defineNuxtPlugin } from '@nuxtjs/composition-api'
import Vue from 'vue'
import { EmitterEvents } from '~/types/events'

export default defineNuxtPlugin((context) => {
  context.$emitter = mitt<EmitterEvents>()
  Vue.prototype.$emmiter = mitt<EmitterEvents>()
})
