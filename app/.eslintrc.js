module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
  },
  extends: ['@nuxtjs/eslint-config-typescript', 'plugin:prettier/recommended', 'plugin:nuxt/recommended'],
  ignorePatterns: ['types/apollo/main/*'],
  plugins: [],
  // add your custom rules here
  rules: {
    'vue/valid-v-slot': ['error', { allowModifiers: true }],
    'prettier/prettier': ['error', { printWidth: 120 }],
  },
}
