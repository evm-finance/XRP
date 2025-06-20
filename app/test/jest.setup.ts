import Vue from 'vue'
import Vuetify from 'vuetify'
Vue.use(Vuetify)
Vue.config.ignoredElements = [
  'v-icon',
  'nuxt-link',
  'client-only',
  'marquee-text',
  'router-link',
  '$ua',
  'v-progress-circular',
]
jest.setTimeout(20000);
