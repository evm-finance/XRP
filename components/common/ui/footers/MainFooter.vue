<template>
  <v-footer padless class="mt-10 overflow-hidden" style="z-index: 1">
    <v-row no-gutters>
      <v-col class="pa-0">
        <v-card flat tile outlined class="text-center">
          <v-card-text class="px-2 pb-2">
            <v-row>
              <v-col cols="7" class="text-left">
                <v-row>
                  <v-col class="d-sm-flex">
                    <div v-for="link in footerLinks" :key="link.title">
                      <nuxt-link
                        :to="link.url"
                        active-class="underline-glow-active"
                        class="white--text text-uppercase caption underline-glow-hover ml-3 mr-3 font-weight-medium pb-1"
                        v-text="link.title"
                      />
                    </div>
                  </v-col>
                </v-row>

                <div class="d-sm-flex font-weight-medium grey--text py-2 pl-3">
                  &copy;
                  {{ new Date().getFullYear() }} - EVM Finance
                  <div class="caption text-no-wrap pl-sm-2" style="margin-top: 2px">
                    contact us at
                    <a
                      href="mailto:admin@quantifycrypto.net"
                      class="text-decoration-none grey--text token-link font-weight-medium"
                    >
                      admin@quantifycrypto.net
                    </a>
                  </div>
                </div>
              </v-col>
              <v-col class="text-right" align-self="end">
                <div class="mt-sm-n13">
                  <v-btn
                    v-for="icon in communityLinks"
                    :key="icon.url"
                    :href="icon.url"
                    target="_blank"
                    class="mx-1 text-hover-primary grey--text"
                    icon
                  >
                    <v-icon v-if="icon.icon.startsWith('mdi')" size="20">
                      {{ icon.icon }}
                    </v-icon>

                    <v-avatar v-else size="20">
                      <v-img class="color-icon-primary" :src="icon.icon" :lazy-src="icon.icon" />
                    </v-avatar>
                  </v-btn>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-footer>
</template>

<script lang="ts">
import { defineComponent, useStore, computed } from '@nuxtjs/composition-api'
import { State } from '~/types/state'

export default defineComponent({
  name: 'MainFooter',
  setup() {
    const { state } = useStore<State>()
    const headerTextClass = computed(() => state.ui[state.ui.theme].headerTextClass)

    const footerLinks = [
      { title: 'Terms and Conditions', url: '/terms-and-conditions' },
      { title: 'Meet Our Team', url: '/team' },
      { title: 'FAQ', url: '/faq' },
      { title: 'About', url: '/about' },
      { title: 'Verse', url: '/verse' },
    ]

    const communityLinks = [
      { icon: 'mdi-twitter', url: 'https://twitter.com/EVMX_IO' },
      { icon: 'mdi-linkedin', url: 'https://www.linkedin.com/company/quantifycrypto' },
      { icon: '/img/discord.svg', url: 'https://discord.gg/NZCtdZPCag' },
      {
        icon: '/img/telegram.svg',
        url: 'https://t.me/Quantify_Crypto',
        style: { filter: 'invert(56%) sepia(95%) saturate(0%) hue-rotate(185deg) brightness(93%) contrast(98%)' },
      },
      { icon: '/img/lenster.svg', url: 'https://lenster.xyz/u/evmfi.lens' },
    ]

    return { footerLinks, communityLinks, headerTextClass }
  },
})
</script>
