import colors from 'vuetify/es5/util/colors'

export default {
  // Global page headers (https://go.nuxtjs.dev/config-head)
  head: {
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: '' },
      { name: 'format-detection', content: 'telephone=no' },
    ],
  },

  // Global CSS (https://go.nuxtjs.dev/config-css)
  css: [],

  plugins: [
    '~/plugins/initConfigs.ts',
    '~/plugins/helper.ts',
    // '~/plugins/errorhandler.apollo.js',
    '~/plugins/apolloClient.ts',

    '~/plugins/web3/web3.ts',
    '~/plugins/typer.client.ts',
  ],

  ssr: true,
  target: 'server',
  components: true,

  buildModules: ['@nuxt/typescript-build', '@nuxtjs/composition-api/module', '@nuxtjs/vuetify'],

  modules: [
    'cookie-universal-nuxt',
    [
      '@nuxtjs/google-analytics',
      {
        ua: process.env.GA_ID,
        debug: { sendHitTask: true },
      },
    ],
    '@nuxtjs/apollo',
    [
      'nuxt-compress',
      {
        gzip: {
          threshold: 8192,
        },
        brotli: {
          threshold: 8192,
        },
      },
    ],
  ],

  // Apollo client setup
  apollo: {
    clientConfigs: {
      default: {
        httpEndpoint: process.env.BASE_GRAPHQL_SERVER_URL,
        wsEndpoint: process.env.BASE_GRAPHQL_WEBSOCKET_URL,
        websocketsOnly: false,
      },
    },
  },

  // Axios module configuration (https://go.nuxtjs.dev/config-axios)
  axios: {
    baseURL: process.env.BASE_URL,
    withCredentials: true,
    debug: false,
  },

  sitemap: {
    hostname: process.env.BASE_URL ?? '',
    gzip: true,
  },

  webfontloader: {
    google: {
      families: ['Roboto:100,300,400,500,700,900&display=swap'],
    },
  },
  // Vuetify module configuration (https://go.nuxtjs.dev/config-vuetify)
  loading: false,
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    treeShake: true,
    theme: {
      // dark: true,
      themes: {
        dark: {
          primary: '#536af6',
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },

        light: {
          primary: '#536af6',
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
    },
  },

  // Build Configuration (https://go.nuxtjs.dev/config-build)
  build: {
    analyze: false,
    build: {},
    filenames: {
      chunk: ({ isDev }) => (isDev ? '[name].js' : '[id].[contenthash].js'),
    },

    extractCSS: false,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    extend(config, ctx) {},
    transpile: ['@apollo/client/core', '@vue/apollo-composable', 'ts-invariant'],
  },

  hooks: {
    render: {
      errorMiddleware(app) {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        app.use((error, req, res, next) => {
          res.writeHead(307, { Location: '/about' })
          res.end()
        })
      },
    },
  },
  env: {
    amChartLicense: process.env.AMCHARTS_LICENSE,
    baseURL: process.env.BASE_URL,
  },

  server: { port: 3000, host: process.env.SERVER_HOST },
}
