<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-7xl mx-auto">
      <!-- Welcome Header -->
      <div class="mb-8">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900 mb-2">
              Bienvenue sur CryptoBank 🏦
            </h1>
            <p class="text-gray-600">
              Tableau de bord - {{ new Date().toLocaleDateString('fr-FR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
            </p>
          </div>
          <div class="flex space-x-3 mt-4 md:mt-0">
            <NuxtLink
              to="/exchange/fiat"
              class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors duration-200"
            >
              💱 Convertir Fiat
            </NuxtLink>
            <NuxtLink
              to="/exchange/crypto"
              class="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors duration-200"
            >
              ₿ Acheter Crypto
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Portfolio Overview -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-blue-100 text-sm font-medium">Solde Total</p>
              <p class="text-3xl font-bold">{{ formatMoney(totalBalance, 'USD') }}</p>
              <p class="text-blue-100 text-sm mt-1">
                <span class="text-green-300">↗ +5.2%</span> aujourd'hui
              </p>
            </div>
            <div class="bg-white/20 p-3 rounded-full">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
              </svg>
            </div>
          </div>
        </div>

        <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-green-100 text-sm font-medium">Portefeuille Crypto</p>
              <p class="text-3xl font-bold">{{ formatMoney(cryptoBalance, 'USD') }}</p>
              <p class="text-green-100 text-sm mt-1">
                <span class="text-green-300">↗ +12.8%</span> cette semaine
              </p>
            </div>
            <div class="bg-white/20 p-3 rounded-full">
              <span class="text-2xl">₿</span>
            </div>
          </div>
        </div>

        <div class="bg-gradient-to-r from-purple-500 to-purple-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-purple-100 text-sm font-medium">Cartes Actives</p>
              <p class="text-3xl font-bold">{{ activeCards }}</p>
              <p class="text-purple-100 text-sm mt-1">
                Solde: {{ formatMoney(cardsBalance, 'USD') }}
              </p>
            </div>
            <div class="bg-white/20 p-3 rounded-full">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
              </svg>
            </div>
          </div>
        </div>

        <div class="bg-gradient-to-r from-orange-500 to-orange-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-orange-100 text-sm font-medium">Transferts ce mois</p>
              <p class="text-3xl font-bold">{{ monthlyTransfers }}</p>
              <p class="text-orange-100 text-sm mt-1">
                Volume: {{ formatMoney(monthlyVolume, 'USD') }}
              </p>
            </div>
            <div class="bg-white/20 p-3 rounded-full">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white rounded-lg shadow-lg p-6 mb-8">
        <h3 class="text-xl font-bold text-gray-900 mb-6">🚀 Actions Rapides</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <NuxtLink
            to="/exchange/crypto"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-blue-400 hover:bg-blue-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">₿</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-blue-600">Acheter Crypto</span>
          </NuxtLink>

          <NuxtLink
            to="/exchange/fiat"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-green-400 hover:bg-green-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">💱</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-green-600">Convertir Devises</span>
          </NuxtLink>

          <NuxtLink
            to="/cards"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-purple-400 hover:bg-purple-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">💳</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-purple-600">Nouvelle Carte</span>
          </NuxtLink>

          <NuxtLink
            to="/transfers"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-orange-400 hover:bg-orange-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">💸</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-orange-600">Envoyer Argent</span>
          </NuxtLink>

          <NuxtLink
            to="/wallets"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-indigo-400 hover:bg-indigo-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">👛</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-indigo-600">Mes Portefeuilles</span>
          </NuxtLink>

          <button
            @click="showQRScanner = true"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-gray-200 hover:border-pink-400 hover:bg-pink-50 transition-all duration-200 group"
          >
            <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">📱</div>
            <span class="text-sm font-semibold text-gray-700 group-hover:text-pink-600">Scanner QR</span>
          </button>
        </div>
      </div>

      <!-- Market Overview and Recent Activity -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Crypto Market -->
        <div class="bg-white rounded-lg shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold text-gray-900">📊 Marchés Crypto</h3>
            <button class="text-blue-600 hover:text-blue-800 font-semibold">
              Voir tout
            </button>
          </div>
          
          <div class="space-y-4">
            <div v-for="crypto in cryptoMarkets" :key="crypto.symbol" class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-full flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold">{{ crypto.symbol.slice(0, 2) }}</span>
                </div>
                <div>
                  <p class="font-semibold text-gray-900">{{ crypto.name }}</p>
                  <p class="text-sm text-gray-500">{{ crypto.symbol }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">${{ crypto.price.toLocaleString() }}</p>
                <p class="text-sm" :class="crypto.change >= 0 ? 'text-green-600' : 'text-red-600'">
                  {{ crypto.change >= 0 ? '+' : '' }}{{ crypto.change.toFixed(2) }}%
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white rounded-lg shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold text-gray-900">🕒 Activité Récente</h3>
            <button class="text-blue-600 hover:text-blue-800 font-semibold">
              Voir tout
            </button>
          </div>
          
          <div class="space-y-4">
            <div v-for="activity in recentActivities" :key="activity.id" class="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50">
              <div class="w-10 h-10 rounded-full flex items-center justify-center" :class="activity.bgColor">
                <span class="text-white">{{ activity.icon }}</span>
              </div>
              <div class="flex-1">
                <p class="font-semibold text-gray-900">{{ activity.title }}</p>
                <p class="text-sm text-gray-500">{{ activity.description }}</p>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">{{ activity.amount }}</p>
                <p class="text-xs text-gray-500">{{ formatTime(activity.time) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Fiat Exchange Rates -->
      <div class="mt-8 bg-white rounded-lg shadow-lg p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-gray-900">💱 Taux de Change Fiduciaires</h3>
          <button 
            @click="refreshRates"
            class="text-blue-600 hover:text-blue-800 font-semibold flex items-center"
          >
            <svg class="w-4 h-4 mr-1" :class="refreshing ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualiser
          </button>
        </div>
        
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div v-for="rate in fiatRates" :key="rate.pair" class="p-4 border border-gray-200 rounded-lg text-center hover:border-blue-300 transition-colors">
            <p class="text-sm font-medium text-gray-600">{{ rate.pair }}</p>
            <p class="text-lg font-bold text-gray-900">{{ rate.rate.toFixed(4) }}</p>
            <p class="text-xs" :class="rate.change >= 0 ? 'text-green-600' : 'text-red-600'">
              {{ rate.change >= 0 ? '+' : '' }}{{ (rate.change * 100).toFixed(2) }}%
            </p>
          </div>
        </div>
      </div>

      <!-- Quick Converters -->
      <div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Crypto Converter -->
        <div class="bg-white rounded-lg shadow-lg p-6">
          <h3 class="text-lg font-bold text-gray-900 mb-4">⚡ Convertisseur Crypto Express</h3>
          <div class="space-y-3">
            <div class="grid grid-cols-2 gap-3">
              <input 
                v-model="cryptoAmount"
                type="number"
                placeholder="1"
                class="p-2 border border-gray-300 rounded text-center"
                @input="calculateCrypto"
              />
              <select v-model="cryptoFrom" class="p-2 border border-gray-300 rounded" @change="calculateCrypto">
                <option value="BTC">₿ BTC</option>
                <option value="ETH">Ξ ETH</option>
                <option value="USD">$ USD</option>
              </select>
            </div>
            <div class="text-center">
              <svg class="w-6 h-6 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <input 
                :value="cryptoResult"
                readonly
                class="p-2 border border-gray-300 rounded text-center bg-gray-50"
              />
              <select v-model="cryptoTo" class="p-2 border border-gray-300 rounded" @change="calculateCrypto">
                <option value="USD">$ USD</option>
                <option value="BTC">₿ BTC</option>
                <option value="ETH">Ξ ETH</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Fiat Converter -->
        <div class="bg-white rounded-lg shadow-lg p-6">
          <h3 class="text-lg font-bold text-gray-900 mb-4">💰 Convertisseur Fiat Express</h3>
          <div class="space-y-3">
            <div class="grid grid-cols-2 gap-3">
              <input 
                v-model="fiatAmount"
                type="number"
                placeholder="1000"
                class="p-2 border border-gray-300 rounded text-center"
                @input="calculateFiat"
              />
              <select v-model="fiatFrom" class="p-2 border border-gray-300 rounded" @change="calculateFiat">
                <option value="USD">🇺🇸 USD</option>
                <option value="EUR">🇪🇺 EUR</option>
                <option value="GBP">🇬🇧 GBP</option>
                <option value="JPY">🇯🇵 JPY</option>
              </select>
            </div>
            <div class="text-center">
              <svg class="w-6 h-6 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <input 
                :value="fiatResult"
                readonly
                class="p-2 border border-gray-300 rounded text-center bg-gray-50"
              />
              <select v-model="fiatTo" class="p-2 border border-gray-300 rounded" @change="calculateFiat">
                <option value="EUR">🇪🇺 EUR</option>
                <option value="USD">🇺🇸 USD</option>
                <option value="GBP">🇬🇧 GBP</option>
                <option value="JPY">🇯🇵 JPY</option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Reactive data
const totalBalance = ref(25750.80)
const cryptoBalance = ref(12500.00)
const cardsBalance = ref(3250.50)
const activeCards = ref(4)
const monthlyTransfers = ref(23)
const monthlyVolume = ref(15750.00)
const refreshing = ref(false)
const showQRScanner = ref(false)

// Converters
const cryptoAmount = ref(1)
const cryptoFrom = ref('BTC')
const cryptoTo = ref('USD')
const cryptoResult = ref('45000.00')

const fiatAmount = ref(1000)
const fiatFrom = ref('USD')
const fiatTo = ref('EUR')
const fiatResult = ref('845.60')

// Market data
const cryptoMarkets = ref([
  { symbol: 'BTC', name: 'Bitcoin', price: 45125, change: 2.35, bgColor: 'bg-orange-500' },
  { symbol: 'ETH', name: 'Ethereum', price: 3024, change: 1.87, bgColor: 'bg-blue-500' },
  { symbol: 'BNB', name: 'Binance Coin', price: 245, change: -0.52, bgColor: 'bg-yellow-500' },
  { symbol: 'ADA', name: 'Cardano', price: 0.48, change: 3.21, bgColor: 'bg-indigo-500' },
])

const fiatRates = ref([
  { pair: 'EUR/USD', rate: 1.0856, change: 0.0023 },
  { pair: 'GBP/USD', rate: 1.2687, change: -0.0012 },
  { pair: 'USD/JPY', rate: 149.23, change: 0.0156 },
  { pair: 'USD/CAD', rate: 1.3567, change: 0.0089 },
  { pair: 'AUD/USD', rate: 0.6564, change: -0.0034 },
  { pair: 'USD/CHF', rate: 0.8912, change: 0.0067 }
])

const recentActivities = ref([
  {
    id: 1,
    icon: '₿',
    title: 'Achat Bitcoin',
    description: '0.02 BTC acheté',
    amount: '$900.50',
    time: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    bgColor: 'bg-orange-500'
  },
  {
    id: 2,
    icon: '💳',
    title: 'Rechargement carte',
    description: 'Carte virtuelle USD',
    amount: '$200.00',
    time: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(),
    bgColor: 'bg-blue-500'
  },
  {
    id: 3,
    icon: '💱',
    title: 'Conversion USD → EUR',
    description: '1000 USD convertis',
    amount: '€845.60',
    time: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(),
    bgColor: 'bg-green-500'
  },
  {
    id: 4,
    icon: '💸',
    title: 'Transfert Mobile Money',
    description: 'Envoyé à +233241234567',
    amount: '$50.00',
    time: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
    bgColor: 'bg-purple-500'
  }
])

// Methods
const calculateCrypto = () => {
  // Mock calculation
  if (cryptoFrom.value === 'BTC' && cryptoTo.value === 'USD') {
    cryptoResult.value = (cryptoAmount.value * 45125).toFixed(2)
  } else if (cryptoFrom.value === 'USD' && cryptoTo.value === 'BTC') {
    cryptoResult.value = (cryptoAmount.value / 45125).toFixed(8)
  } else if (cryptoFrom.value === 'ETH' && cryptoTo.value === 'USD') {
    cryptoResult.value = (cryptoAmount.value * 3024).toFixed(2)
  } else {
    cryptoResult.value = '0.00'
  }
}

const calculateFiat = () => {
  // Mock calculation using real-ish rates
  const rates = {
    'USD-EUR': 0.8456,
    'EUR-USD': 1.1826,
    'USD-GBP': 0.7821,
    'GBP-USD': 1.2787,
    'USD-JPY': 149.23,
    'JPY-USD': 0.0067
  }
  
  const key = `${fiatFrom.value}-${fiatTo.value}`
  const rate = rates[key] || 1.0
  fiatResult.value = (fiatAmount.value * rate).toFixed(2)
}

const refreshRates = async () => {
  refreshing.value = true
  setTimeout(() => {
    // Simulate rate refresh
    fiatRates.value.forEach(rate => {
      rate.change = (Math.random() - 0.5) * 0.02 // Random change ±1%
      rate.rate = rate.rate * (1 + rate.change)
    })
    refreshing.value = false
  }, 1000)
}

const formatMoney = (amount, currency) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency
  }).format(amount)
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit'
  })
}

// Lifecycle
onMounted(() => {
  calculateCrypto()
  calculateFiat()
})

// SEO
definePageMeta({
  title: 'Tableau de Bord - CryptoBank',
  description: 'Gérez vos finances crypto et fiat depuis un tableau de bord unifié'
})
</script>