// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt'
  ],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080/api/v1',
      appName: 'CryptoBank',
      appVersion: '1.0.0'
    }
  },
  ssr: true,
  nitro: {
    experimental: {
      wasm: true
    }
  },
  app: {
    head: {
      title: 'CryptoBank - Secure Digital Banking',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Secure digital banking with cryptocurrency and fiat currency support' },
        { name: 'format-detection', content: 'telephone=no' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },
  tailwindcss: {
    cssPath: '~/assets/css/main.css',
  },
  build: {
    transpile: ['chart.js']
  }
})