import colors from 'vuetify/es5/util/colors'

const graphqlServerUrl = process.env.BASE_GRAPHQL_SERVER_URL || 'http://127.0.0.1:8080/query'
console.log('GraphQL Server URL:', graphqlServerUrl)

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
    '~/plugins/apolloClient.ts',
    '~/plugins/emitter.client.ts',
    '~/plugins/web3/xrp.client.ts',
    '~/plugins/typer.client.ts',
  ],

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

    '@nuxtjs/composition-api/module',
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
        httpEndpoint: graphqlServerUrl,
        wsEndpoint: process.env.BASE_GRAPHQL_WEBSOCKET_URL || 'ws://127.0.0.1:8080/query',
        websocketsOnly: false,
        // Improved error handling for missing server
        onError: (error) => {
          console.error('GraphQL server error:', {
            message: error.message,
            networkError: error.networkError,
            graphQLErrors: error.graphQLErrors,
            timestamp: new Date().toISOString(),
            endpoint: graphqlServerUrl
          })
          
          // Log to external service if available
          if (process.env.NODE_ENV === 'production') {
            // TODO: Add proper error reporting service
            console.warn('Production error - should be reported to monitoring service')
          }
        }
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
    extend(config, ctx) {
      // Exclude quantify subproject from build
      config.watchOptions = {
        ...config.watchOptions,
        ignored: [
          ...(config.watchOptions?.ignored || []),
          '**/quantify/**',
          '**/qc-front-end-server/**'
        ]
      }
      
      // Exclude quantify directory from webpack processing
      if (config.module && config.module.rules) {
        config.module.rules.push({
          test: /\.(js|ts|vue)$/,
          exclude: [
            /node_modules/,
            /quantify/,
            /qc-front-end-server/
          ]
        })
      }
    },
    transpile: [
      // 'tslib',
      // '@apollo/client',
      '@apollo/client/core',
      '@vue/apollo-composable',
      'ts-invariant',
      // '@vue/apollo-option',
      '@noble/curves',
      '@noble/hashes',
      'xrpl',
    ],
  },

  hooks: {
    render: {
      // errorMiddleware(app) {
      //   // eslint-disable-next-line @typescript-eslint/no-unused-vars,n/handle-callback-err
      //   app.use((error, req, res, next) => {
      //     console.error(error)
      //     res.writeHead(307, { Location: '/about' })
      //     res.end()
      //   })
      // },
    },
  },
  env: {
    amChartLicense: process.env.AMCHARTS_LICENSE,
    baseURL: process.env.BASE_URL,
  },

  server: { port: 3000, host: process.env.SERVER_HOST },
}
