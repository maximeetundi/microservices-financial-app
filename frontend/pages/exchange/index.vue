<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Exchange Center</h1>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-500">24h Volume:</span>
            <span class="text-lg font-semibold text-blue-600">{{ formattedVolume }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Exchange Options -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
        
        <!-- Crypto Exchange -->
        <NuxtLink to="/exchange/crypto" 
                  class="bg-white rounded-xl shadow-lg p-8 hover:shadow-xl transition-shadow group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-orange-100 rounded-lg">
              <svg class="w-8 h-8 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Cryptocurrency</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Trade Bitcoin, Ethereum, and 100+ cryptocurrencies with instant execution and competitive rates.
          </p>
          <div class="flex items-center text-blue-600 group-hover:text-blue-800">
            <span class="text-sm font-medium">Start Trading</span>
            <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          </div>
          <div class="mt-4 text-sm text-gray-500">
            <div class="flex justify-between">
              <span>BTC/USD:</span>
              <span :class="cryptoPrices.BTC.change >= 0 ? 'text-green-600' : 'text-red-600'">
                ${{ formatPrice(cryptoPrices.BTC.price) }} ({{ formatChange(cryptoPrices.BTC.change) }})
              </span>
            </div>
            <div class="flex justify-between">
              <span>ETH/USD:</span>
              <span :class="cryptoPrices.ETH.change >= 0 ? 'text-green-600' : 'text-red-600'">
                ${{ formatPrice(cryptoPrices.ETH.price) }} ({{ formatChange(cryptoPrices.ETH.change) }})
              </span>
            </div>
          </div>
        </NuxtLink>

        <!-- Fiat Exchange -->
        <NuxtLink to="/exchange/fiat" 
                  class="bg-white rounded-xl shadow-lg p-8 hover:shadow-xl transition-shadow group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-green-100 rounded-lg">
              <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Fiat Currencies</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Exchange between major world currencies with bank-level security and transparent fees.
          </p>
          <div class="flex items-center text-blue-600 group-hover:text-blue-800">
            <span class="text-sm font-medium">Exchange Now</span>
            <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          </div>
          <div class="mt-4 text-sm text-gray-500">
            <div class="flex justify-between">
              <span>EUR/USD:</span>
              <span :class="(fiatRates.EUR_USD?.change || 0) >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ (fiatRates.EUR_USD?.rate || 0).toFixed(4) }} ({{ formatChange(fiatRates.EUR_USD?.change || 0) }})
              </span>
            </div>
            <div class="flex justify-between">
              <span>GBP/USD:</span>
              <span :class="(fiatRates.GBP_USD?.change || 0) >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ (fiatRates.GBP_USD?.rate || 0).toFixed(4) }} ({{ formatChange(fiatRates.GBP_USD?.change || 0) }})
              </span>
            </div>
          </div>
        </NuxtLink>

        <!-- Advanced Trading -->
        <NuxtLink to="/exchange/trading" 
                  class="bg-white rounded-xl shadow-lg p-8 hover:shadow-xl transition-shadow group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-purple-100 rounded-lg">
              <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 00-2-2z"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Advanced Trading</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Professional trading interface with limit orders, stop-loss, and portfolio management tools.
          </p>
          <div class="flex items-center text-blue-600 group-hover:text-blue-800">
            <span class="text-sm font-medium">Start Trading</span>
            <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          </div>
          <div class="mt-4 text-sm text-gray-500">
            <div class="flex justify-between">
              <span>24h Volume:</span>
              <span>{{ formattedVolume }}</span>
            </div>
            <div class="flex justify-between">
              <span>Active Pairs:</span>
              <span>{{ markets.length }}+</span>
            </div>
          </div>
        </NuxtLink>
      </div>

      <!-- Market Overview -->
      <div class="bg-white rounded-xl shadow-lg p-6 mb-8">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold text-gray-900">Market Overview</h2>
          <div class="flex space-x-4">
            <button @click="marketTab = 'crypto'" 
                    :class="marketTab === 'crypto' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500'"
                    class="pb-2">Crypto</button>
            <button @click="marketTab = 'fiat'" 
                    :class="marketTab === 'fiat' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500'"
                    class="pb-2">Fiat</button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="market in filteredMarkets" :key="market.symbol" 
               class="p-4 border rounded-lg hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between mb-2">
              <h3 class="font-semibold">{{ market.symbol }}</h3>
              <span :class="market.change >= 0 ? 'text-green-600' : 'text-red-600'" 
                    class="text-sm">
                {{ market.change >= 0 ? '+' : '' }}{{ market.change }}%
              </span>
            </div>
            <div class="text-lg font-bold">${{ (market.price || 0).toLocaleString() }}</div>
            <div class="text-sm text-gray-500">Vol: ${{ (market.volume || 0).toLocaleString() }}</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <!-- Quick Convert -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Convert</h3>
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">From</label>
                <select v-model="quickConvert.from" class="w-full rounded-md border-gray-300 shadow-sm">
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                  <option value="BTC">BTC</option>
                  <option value="ETH">ETH</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">To</label>
                <select v-model="quickConvert.to" class="w-full rounded-md border-gray-300 shadow-sm">
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                  <option value="BTC">BTC</option>
                  <option value="ETH">ETH</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Amount</label>
              <input v-model.number="quickConvert.amount" type="number" step="0.01" 
                     class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
            </div>
            <div v-if="convertedAmount" class="p-3 bg-gray-50 rounded-lg">
              <div class="text-sm text-gray-600">You will receive:</div>
              <div class="text-lg font-semibold">{{ convertedAmount }} {{ quickConvert.to }}</div>
            </div>
            <button @click="performQuickConvert" 
                    class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700">
              Convert
            </button>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
          <div class="space-y-4">
            <div v-for="activity in recentActivity" :key="activity.id" 
                 class="flex items-center justify-between py-2 border-b last:border-b-0">
              <div>
                <div class="font-medium text-sm">{{ activity.type }}</div>
                <div class="text-xs text-gray-500">{{ activity.pair }} • {{ formatDate(activity.date) }}</div>
              </div>
              <div class="text-right">
                <div class="text-sm font-medium">{{ activity.amount }}</div>
                <div class="text-xs text-gray-500">{{ activity.status }}</div>
              </div>
            </div>
          </div>
          <NuxtLink to="/orders" class="block mt-4 text-blue-600 text-sm hover:text-blue-800">
            View All Activity →
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { exchangeAPI } from '~/composables/useApi'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

// Reactive data
const marketTab = ref('crypto')
const quickConvert = ref({
  from: 'USD',
  to: 'BTC',
  amount: 1000
})
const convertedAmount = ref(null)
const recentActivity = ref([])
const markets = ref([])
const loading = ref(true)

// Dynamic price data for header cards
const totalVolume = ref(0)
const cryptoPrices = ref({
  BTC: { price: 0, change: 0 },
  ETH: { price: 0, change: 0 }
})
const fiatRates = ref({
  EUR_USD: { rate: 0, change: 0 },
  GBP_USD: { rate: 0, change: 0 }
})

// Computed
const filteredMarkets = computed(() => {
  return markets.value.filter(market => market.type === marketTab.value)
})

const formattedVolume = computed(() => {
  if (totalVolume.value >= 1e9) return `$${(totalVolume.value / 1e9).toFixed(1)}B`
  if (totalVolume.value >= 1e6) return `$${(totalVolume.value / 1e6).toFixed(1)}M`
  return `$${(totalVolume.value || 0).toLocaleString()}`
})

// Methods
const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const formatPrice = (price) => {
  if (price == null) return '0'
  if (price >= 1000) return (price || 0).toLocaleString()
  if (price >= 1) return (price || 0).toFixed(2)
  return (price || 0).toFixed(6)
}

const formatChange = (change) => {
  const val = change || 0
  const sign = val >= 0 ? '+' : ''
  return `${sign}${val.toFixed(1)}%`
}

const performQuickConvert = async () => {
  try {
    if (!quickConvert.value.amount || quickConvert.value.amount <= 0) return

    const { data } = await exchangeAPI.convert(
      quickConvert.value.from,
      quickConvert.value.to,
      quickConvert.value.amount
    )
    if (data) {
      convertedAmount.value = data.to_amount
    }
  } catch (error) {
    console.error('Error converting:', error)
  }
}

// Fetch data
const fetchRecentActivity = async () => {
  try {
    const { data } = await exchangeAPI.getHistory(5)
    if (data && data.exchanges) {
      recentActivity.value = data.exchanges.map(ex => ({
        id: ex.id,
        type: 'Exchange',
        pair: `${ex.from_currency}/${ex.to_currency}`,
        amount: `${ex.from_amount} ${ex.from_currency}`,
        status: ex.status,
        date: ex.created_at || new Date().toISOString()
      }))
    }
  } catch (error) {
    console.error('Error fetching recent activity:', error)
  }
}

const fetchMarkets = async () => {
  try {
     const { data } = await exchangeAPI.getMarkets()
     if (data && data.markets) {
       // Transform API response to UI format
       markets.value = data.markets.map(m => ({
         symbol: m.Symbol || m.symbol,
         price: m.Price || m.price || 0,
         change: m.Change24h || m.change_24h || 0,
         volume: m.Volume24h || m.volume_24h || 0,
         type: 'crypto'
       }))
       
       // Calculate total volume
       totalVolume.value = markets.value.reduce((sum, m) => sum + (m.volume || 0), 0)
       
       // Extract BTC and ETH prices for the cards
       const btc = markets.value.find(m => m.symbol?.includes('BTC'))
       const eth = markets.value.find(m => m.symbol?.includes('ETH'))
       
       if (btc) {
         cryptoPrices.value.BTC = { price: btc.price, change: btc.change }
       }
       if (eth) {
         cryptoPrices.value.ETH = { price: eth.price, change: eth.change }
       }
     }
  } catch (error) {
     console.error('Error fetching markets:', error)
  }
}

const fetchRates = async () => {
  try {
    const { data } = await exchangeAPI.getRates()
    if (data && data.rates) {
      // Extract EUR/USD and GBP/USD rates
      const eurUsd = data.rates['EUR/USD'] || data.rates['EUR_USD']
      const gbpUsd = data.rates['GBP/USD'] || data.rates['GBP_USD']
      
      if (eurUsd) {
        fiatRates.value.EUR_USD = { 
          rate: eurUsd.Rate || eurUsd.rate || eurUsd, 
          change: eurUsd.Change24h || eurUsd.change_24h || 0 
        }
      }
      if (gbpUsd) {
        fiatRates.value.GBP_USD = { 
          rate: gbpUsd.Rate || gbpUsd.rate || gbpUsd, 
          change: gbpUsd.Change24h || gbpUsd.change_24h || 0 
        }
      }
      
      // Add fiat markets to the markets array
      if (!markets.value.some(m => m.type === 'fiat')) {
        Object.entries(data.rates).forEach(([pair, rateData]) => {
          if (!pair.includes('BTC') && !pair.includes('ETH')) {
            markets.value.push({
              symbol: pair.replace('_', '/'),
              price: rateData.Rate || rateData.rate || rateData,
              change: rateData.Change24h || rateData.change_24h || 0,
              volume: 0,
              type: 'fiat'
            })
          }
        })
      }
    }
  } catch (error) {
    console.error('Error fetching rates:', error)
  }
}

// Auto-update converted amount when inputs change
watch([() => quickConvert.value.from, () => quickConvert.value.to, () => quickConvert.value.amount], 
  async () => {
    if (quickConvert.value.amount > 0) {
      await performQuickConvert()
    }
  }
)

// Lifecycle
onMounted(async () => {
  loading.value = true
  try {
    // Parallel fetch
    await Promise.all([
      fetchMarkets(),
      fetchRates(),
      fetchRecentActivity(),
      performQuickConvert()
    ])
  } finally {
    loading.value = false
  }
})
</script>