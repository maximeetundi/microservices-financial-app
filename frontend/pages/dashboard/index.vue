<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Welcome back!</h1>
            <p class="text-gray-600 mt-1">Here's what's happening with your crypto portfolio today.</p>
          </div>
          <div class="text-right">
            <div class="text-sm text-gray-500">Last updated</div>
            <div class="text-sm font-medium">{{ new Date().toLocaleString() }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      
      <!-- Portfolio Overview Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- Total Portfolio Value -->
        <div class="bg-white rounded-xl shadow-lg p-6 border-l-4 border-blue-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total Portfolio Value</p>
              <p class="text-3xl font-bold text-gray-900">${{ portfolioStats.totalValue?.toLocaleString() || '0' }}</p>
              <p class="text-sm mt-2" :class="portfolioStats.dailyChange >= 0 ? 'text-green-600' : 'text-red-600'">
                <span>{{ portfolioStats.dailyChange >= 0 ? '+' : '' }}{{ portfolioStats.dailyChange }}%</span>
                <span class="text-gray-500 ml-1">today</span>
              </p>
            </div>
            <div class="p-3 bg-blue-100 rounded-lg">
              <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 00-2-2z"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Available Cash -->
        <div class="bg-white rounded-xl shadow-lg p-6 border-l-4 border-green-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Available Cash</p>
              <p class="text-3xl font-bold text-gray-900">${{ portfolioStats.availableCash?.toLocaleString() || '0' }}</p>
              <p class="text-sm text-gray-500 mt-2">Ready to invest</p>
            </div>
            <div class="p-3 bg-green-100 rounded-lg">
              <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Total Profit/Loss -->
        <div class="bg-white rounded-xl shadow-lg p-6 border-l-4 border-purple-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total P&L</p>
              <p class="text-3xl font-bold" :class="portfolioStats.totalPnL >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ portfolioStats.totalPnL >= 0 ? '+' : '' }}${{ Math.abs(portfolioStats.totalPnL || 0).toLocaleString() }}
              </p>
              <p class="text-sm text-gray-500 mt-2">
                {{ portfolioStats.totalPnLPercent >= 0 ? '+' : '' }}{{ portfolioStats.totalPnLPercent }}% all time
              </p>
            </div>
            <div class="p-3 bg-purple-100 rounded-lg">
              <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Active Orders -->
        <div class="bg-white rounded-xl shadow-lg p-6 border-l-4 border-yellow-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Active Orders</p>
              <p class="text-3xl font-bold text-gray-900">{{ portfolioStats.activeOrders || 0 }}</p>
              <p class="text-sm text-gray-500 mt-2">Pending execution</p>
            </div>
            <div class="p-3 bg-yellow-100 rounded-lg">
              <svg class="w-8 h-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        
        <!-- Portfolio Chart -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Portfolio Performance</h3>
            <div class="flex space-x-2">
              <button v-for="period in chartPeriods" :key="period.value"
                      @click="selectedPeriod = period.value"
                      :class="selectedPeriod === period.value 
                        ? 'bg-blue-600 text-gray-900' 
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                      class="px-3 py-1 text-sm rounded">
                {{ period.label }}
              </button>
            </div>
          </div>
          <div class="h-64 bg-gray-50 rounded-lg flex items-center justify-center">
            <span class="text-gray-500">Portfolio Chart (Chart.js integration)</span>
          </div>
        </div>

        <!-- Top Holdings -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-6">Top Holdings</h3>
          <div class="space-y-4">
            <div v-for="holding in topHoldings" :key="holding.currency" 
                 class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
              <div class="flex items-center">
                <div class="w-10 h-10 rounded-full flex items-center justify-center text-gray-900 font-bold"
                     :class="getCurrencyColor(holding.currency)">
                  {{ getCurrencySymbol(holding.currency) }}
                </div>
                <div class="ml-4">
                  <h4 class="font-medium text-gray-900">{{ holding.currency }}</h4>
                  <p class="text-sm text-gray-500">{{ holding.amount }} {{ holding.currency }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-medium text-gray-900">${{ holding.value?.toLocaleString() }}</p>
                <p class="text-sm" :class="holding.change >= 0 ? 'text-green-600' : 'text-red-600'">
                  {{ holding.change >= 0 ? '+' : '' }}{{ holding.change }}%
                </p>
              </div>
            </div>
          </div>
          <NuxtLink to="/portfolio" class="block mt-4 text-blue-600 text-sm hover:text-blue-800">
            View Full Portfolio →
          </NuxtLink>
        </div>
      </div>

      <!-- Quick Actions & Market Overview -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        
        <!-- Quick Actions -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-6">Quick Actions</h3>
          <div class="grid grid-cols-2 gap-4">
            <NuxtLink to="/wallet" 
                      class="flex flex-col items-center p-4 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors group">
              <div class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center mb-3 group-hover:bg-blue-700">
                <svg class="w-6 h-6 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
                </svg>
              </div>
              <span class="text-sm font-medium text-gray-700">Wallets</span>
            </NuxtLink>

            <NuxtLink to="/exchange" 
                      class="flex flex-col items-center p-4 bg-green-50 rounded-lg hover:bg-green-100 transition-colors group">
              <div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mb-3 group-hover:bg-green-700">
                <svg class="w-6 h-6 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"></path>
                </svg>
              </div>
              <span class="text-sm font-medium text-gray-700">Exchange</span>
            </NuxtLink>

            <NuxtLink to="/transfer" 
                      class="flex flex-col items-center p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors group">
              <div class="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center mb-3 group-hover:bg-purple-700">
                <svg class="w-6 h-6 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
                </svg>
              </div>
              <span class="text-sm font-medium text-gray-700">Transfer</span>
            </NuxtLink>

            <NuxtLink to="/cards" 
                      class="flex flex-col items-center p-4 bg-yellow-50 rounded-lg hover:bg-yellow-100 transition-colors group">
              <div class="w-12 h-12 bg-yellow-600 rounded-lg flex items-center justify-center mb-3 group-hover:bg-yellow-700">
                <svg class="w-6 h-6 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
                </svg>
              </div>
              <span class="text-sm font-medium text-gray-700">Cards</span>
            </NuxtLink>
          </div>
        </div>

        <!-- Market Overview -->
        <div class="lg:col-span-2 bg-white rounded-xl shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Market Overview</h3>
            <NuxtLink to="/exchange/trading" class="text-blue-600 text-sm hover:text-blue-800">
              View All Markets →
            </NuxtLink>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full">
              <thead>
                <tr class="text-left text-sm text-gray-500 border-b">
                  <th class="pb-3">Asset</th>
                  <th class="pb-3">Price</th>
                  <th class="pb-3">24h Change</th>
                  <th class="pb-3">Volume</th>
                  <th class="pb-3">Action</th>
                </tr>
              </thead>
              <tbody class="space-y-2">
                <tr v-for="market in marketData" :key="market.symbol" class="border-b border-gray-50">
                  <td class="py-3">
                    <div class="flex items-center">
                      <div class="w-8 h-8 rounded-full flex items-center justify-center text-gray-900 text-sm font-bold mr-3"
                           :class="getCurrencyColor(market.symbol.split('/')[0])">
                        {{ getCurrencySymbol(market.symbol.split('/')[0]) }}
                      </div>
                      <div>
                        <div class="font-medium text-gray-900">{{ market.symbol }}</div>
                        <div class="text-sm text-gray-500">{{ getCurrencyName(market.symbol.split('/')[0]) }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="py-3 font-medium">${{ market.price?.toLocaleString() }}</td>
                  <td class="py-3" :class="market.change >= 0 ? 'text-green-600' : 'text-red-600'">
                    {{ market.change >= 0 ? '+' : '' }}{{ market.change }}%
                  </td>
                  <td class="py-3 text-sm text-gray-500">${{ market.volume?.toLocaleString() }}</td>
                  <td class="py-3">
                    <button @click="quickTrade(market)" 
                            class="text-blue-600 text-sm hover:text-blue-800 font-medium">
                      Trade
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white rounded-xl shadow-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Recent Activity</h3>
            <NuxtLink to="/orders" class="text-blue-600 text-sm hover:text-blue-800">
              View All →
            </NuxtLink>
          </div>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Asset</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="activity in recentActivity" :key="activity.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full"
                        :class="getActivityTypeColor(activity.type)">
                    {{ activity.type }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ activity.asset }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ activity.amount }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full"
                        :class="getStatusColor(activity.status)">
                    {{ activity.status }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(activity.date) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

// Reactive data
const selectedPeriod = ref('7d')
const portfolioStats = ref({})
const topHoldings = ref([])
const marketData = ref([])
const recentActivity = ref([])

const chartPeriods = [
  { value: '1d', label: '1D' },
  { value: '7d', label: '7D' },
  { value: '1m', label: '1M' },
  { value: '3m', label: '3M' },
  { value: '1y', label: '1Y' }
]

// Methods
const getCurrencyColor = (currency) => {
  const colors = {
    'BTC': 'bg-orange-500',
    'ETH': 'bg-blue-500',
    'LTC': 'bg-gray-400',
    'ADA': 'bg-blue-600',
    'USD': 'bg-green-500',
    'EUR': 'bg-blue-800'
  }
  return colors[currency] || 'bg-gray-500'
}

const getCurrencySymbol = (currency) => {
  const symbols = {
    'BTC': '₿',
    'ETH': 'Ξ',
    'LTC': 'Ł',
    'ADA': '₳',
    'USD': '$',
    'EUR': '€'
  }
  return symbols[currency] || currency.substring(0, 2)
}

const getCurrencyName = (currency) => {
  const names = {
    'BTC': 'Bitcoin',
    'ETH': 'Ethereum',
    'LTC': 'Litecoin',
    'ADA': 'Cardano',
    'USD': 'US Dollar',
    'EUR': 'Euro'
  }
  return names[currency] || currency
}

const getActivityTypeColor = (type) => {
  const colors = {
    'buy': 'bg-green-100 text-green-800',
    'sell': 'bg-red-100 text-red-800',
    'transfer': 'bg-blue-100 text-blue-800',
    'exchange': 'bg-purple-100 text-purple-800'
  }
  return colors[type] || 'bg-gray-100 text-gray-800'
}

const getStatusColor = (status) => {
  const colors = {
    'completed': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'failed': 'bg-red-100 text-red-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const quickTrade = (market) => {
  navigateTo(`/exchange/trading?pair=${market.symbol}`)
}

// Fetch data functions
const fetchDashboardData = async () => {
  try {
    const data = await $fetch('/api/dashboard/overview')
    portfolioStats.value = data.stats || {}
    topHoldings.value = data.holdings || []
    recentActivity.value = data.activity || []
  } catch (error) {
    console.error('Error fetching dashboard data:', error)
    // Fallback demo data
    portfolioStats.value = {
      totalValue: 15750.50,
      availableCash: 2500.00,
      totalPnL: 2150.75,
      totalPnLPercent: 15.8,
      dailyChange: 2.3,
      activeOrders: 3
    }

    topHoldings.value = [
      {
        currency: 'BTC',
        amount: 0.25,
        value: 10875,
        change: 2.3
      },
      {
        currency: 'ETH',
        amount: 2.0,
        value: 4900,
        change: -1.2
      },
      {
        currency: 'USD',
        amount: 2500,
        value: 2500,
        change: 0
      }
    ]

    recentActivity.value = [
      {
        id: '1',
        type: 'buy',
        asset: 'BTC',
        amount: '0.1 BTC',
        status: 'completed',
        date: new Date().toISOString()
      },
      {
        id: '2',
        type: 'transfer',
        asset: 'ETH',
        amount: '1.5 ETH',
        status: 'pending',
        date: new Date().toISOString()
      }
    ]
  }
}

const fetchMarketData = async () => {
  try {
    const data = await $fetch('/api/market/overview')
    marketData.value = data.markets || []
  } catch (error) {
    console.error('Error fetching market data:', error)
    // Fallback demo data
    marketData.value = [
      { symbol: 'BTC/USD', price: 43500, change: 2.3, volume: 1234567 },
      { symbol: 'ETH/USD', price: 2450, change: -1.2, volume: 987654 },
      { symbol: 'LTC/USD', price: 72, change: 0.8, volume: 543210 },
      { symbol: 'ADA/USD', price: 0.52, change: 3.1, volume: 876543 }
    ]
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchDashboardData(),
    fetchMarketData()
  ])
})
</script>