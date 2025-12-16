// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: false },
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt'
  ],
  css: ['~/assets/css/main.css'],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: 'https://api.app.maximeetundi.store',
      appName: 'CryptoBank',
      appVersion: '1.0.0'
    }
  },
  ssr: true,
  nitro: {
    prerender: {
      failOnError: false,
      crawlLinks: false
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
  build: {
    transpile: ['chart.js']
  }
})
