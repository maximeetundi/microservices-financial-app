// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: false },
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxtjs/color-mode'
  ],
  colorMode: {
    classSuffix: '',
    preference: 'dark',
    fallback: 'dark'
  },
  css: ['~/assets/css/main.css'],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: 'https://api.app.tech-afm.com',
      appName: 'Zekora',
      appVersion: '1.0.0',
      cryptoNetwork: process.env.CRYPTO_NETWORK || 'mainnet'
    }
  },
  ssr: true,
  app: {
    head: {
      title: 'Zekora - Secure Digital Banking',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Secure digital banking with cryptocurrency and fiat currency support' },
        { name: 'format-detection', content: 'telephone=no' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
        { rel: 'manifest', href: '/site.webmanifest' }
      ]
    }
  },
  build: {
    transpile: ['chart.js']
  },
  nitro: {
    compressPublicAssets: true,
    routeRules: {
      '/**': {
        headers: {
          'Content-Security-Policy': "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://www.paypal.com https://www.sandbox.paypal.com https://www.paypalobjects.com; connect-src 'self' https://www.paypal.com https://www.sandbox.paypal.com https://api-m.paypal.com https://api-m.sandbox.paypal.com; img-src 'self' data: https://www.paypalobjects.com https://www.paypal.com https://www.sandbox.paypal.com; style-src 'self' 'unsafe-inline'; frame-src https://www.paypal.com https://www.sandbox.paypal.com;"
        }
      },
      '/_nuxt/**': {
        headers: { 'Cache-Control': 'public, max-age=31536000, immutable' }
      }
    }
  },
  experimental: {
    payloadExtraction: false
  }
})
