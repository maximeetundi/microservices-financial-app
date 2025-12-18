<template>
  <div>
    <!-- Header -->
    <div class="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4 animate-fade-in-up">
      <div>
        <h1 class="text-3xl font-bold text-base">Bonjour, {{ userName }} 👋</h1>
        <p class="text-muted mt-1">Voici un aperçu de votre portefeuille aujourd'hui.</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="px-4 py-2 rounded-xl bg-surface border border-secondary-200 dark:border-secondary-800 text-sm font-medium text-muted">
          Dernière MAJ: {{ new Date().toLocaleTimeString() }}
        </div>
        <NuxtLink to="/transfer" class="btn-primary flex items-center gap-2">
          <span>Envoyer</span>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/>
          </svg>
        </NuxtLink>
      </div>
    </div>

    <!-- Portfolio Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8 animate-fade-in-up" style="animation-delay: 0.1s">
      <!-- Total Portfolio Value -->
      <div class="glass-card relative overflow-hidden group">
        <div class="absolute top-0 right-0 w-32 h-32 bg-primary-500/10 rounded-full blur-3xl group-hover:bg-primary-500/20 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-medium text-muted">Valeur Totale</p>
            <div class="p-2 bg-primary-50 dark:bg-primary-900/20 rounded-lg text-primary-600 dark:text-primary-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-bold text-base mb-2">${{ portfolioStats.totalValue?.toLocaleString() || '0' }}</p>
          <div class="flex items-center text-sm">
            <span :class="portfolioStats.dailyChange >= 0 ? 'text-success' : 'text-error'" class="font-medium flex items-center gap-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="portfolioStats.dailyChange >= 0" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"/>
              </svg>
              {{ Math.abs(portfolioStats.dailyChange) }}%
            </span>
            <span class="text-muted ml-2">vs hier</span>
          </div>
        </div>
      </div>

      <!-- Available Cash -->
      <div class="glass-card relative overflow-hidden group">
        <div class="absolute top-0 right-0 w-32 h-32 bg-green-500/10 rounded-full blur-3xl group-hover:bg-green-500/20 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-medium text-muted">Disponible</p>
            <div class="p-2 bg-green-50 dark:bg-green-900/20 rounded-lg text-green-600 dark:text-green-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-bold text-base mb-2">${{ portfolioStats.availableCash?.toLocaleString() || '0' }}</p>
          <p class="text-sm text-muted">Prêt à investir</p>
        </div>
      </div>

      <!-- Total Profit/Loss -->
      <div class="glass-card relative overflow-hidden group">
        <div class="absolute top-0 right-0 w-32 h-32 bg-purple-500/10 rounded-full blur-3xl group-hover:bg-purple-500/20 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-medium text-muted">P&L Total</p>
            <div class="p-2 bg-purple-50 dark:bg-purple-900/20 rounded-lg text-purple-600 dark:text-purple-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 00-2-2z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-bold text-base mb-2" :class="portfolioStats.totalPnL >= 0 ? 'text-success' : 'text-error'">
            {{ portfolioStats.totalPnL >= 0 ? '+' : '' }}${{ Math.abs(portfolioStats.totalPnL || 0).toLocaleString() }}
          </p>
          <p class="text-sm text-muted">Total réalisé</p>
        </div>
      </div>

      <!-- Active Orders -->
      <div class="glass-card relative overflow-hidden group">
        <div class="absolute top-0 right-0 w-32 h-32 bg-orange-500/10 rounded-full blur-3xl group-hover:bg-orange-500/20 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-medium text-muted">Ordres Actifs</p>
            <div class="p-2 bg-orange-50 dark:bg-orange-900/20 rounded-lg text-orange-600 dark:text-orange-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-bold text-base mb-2">{{ portfolioStats.activeOrders || 0 }}</p>
          <p class="text-sm text-muted">En attente</p>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8 animate-fade-in-up" style="animation-delay: 0.2s">
      
      <!-- Portfolio Chart -->
      <div class="glass-card">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-base">Performance</h3>
          <div class="flex space-x-1 bg-secondary-100 dark:bg-secondary-800 p-1 rounded-lg">
            <button v-for="period in chartPeriods" :key="period.value"
                    @click="selectedPeriod = period.value"
                    :class="selectedPeriod === period.value 
                      ? 'bg-white dark:bg-secondary-700 text-primary shadow-sm' 
                      : 'text-muted hover:text-base'"
                    class="px-3 py-1 text-sm font-medium rounded-md transition-all">
              {{ period.label }}
            </button>
          </div>
        </div>
        <div class="h-64 flex items-center justify-center border border-dashed border-secondary-200 dark:border-secondary-700 rounded-xl bg-surface/50">
          <span class="text-muted">Graphique Portefeuille (Chart.js)</span>
        </div>
      </div>

      <!-- Top Holdings -->
      <div class="glass-card">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-base">Top Actifs</h3>
          <NuxtLink to="/wallet" class="text-primary text-sm font-medium hover:text-primary-600 transition-colors">
            Voir tout
          </NuxtLink>
        </div>
        <div class="space-y-3">
          <div v-for="holding in topHoldings" :key="holding.currency" 
               class="flex items-center justify-between p-3 rounded-xl hover:bg-surface-hover transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shadow-lg"
                   :class="getCurrencyColor(holding.currency)">
                {{ getCurrencySymbol(holding.currency) }}
              </div>
              <div>
                <h4 class="font-bold text-base group-hover:text-primary transition-colors">{{ holding.currency }}</h4>
                <p class="text-sm text-muted">{{ holding.amount }} {{ holding.currency }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-base">${{ holding.value?.toLocaleString() }}</p>
              <p class="text-sm font-medium" :class="holding.change >= 0 ? 'text-success' : 'text-error'">
                {{ holding.change >= 0 ? '+' : '' }}{{ holding.change }}%
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Market Overview -->
    <div class="glass-card mb-8 animate-fade-in-up" style="animation-delay: 0.3s">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-bold text-base">Marché</h3>
        <NuxtLink to="/exchange/crypto" class="text-primary text-sm font-medium hover:text-primary-600 transition-colors">
          Voir tout les marchés
        </NuxtLink>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-xs font-semibold text-muted uppercase tracking-wider border-b border-secondary-200 dark:border-secondary-800">
              <th class="pb-4 pl-4">Actif</th>
              <th class="pb-4">Prix</th>
              <th class="pb-4">24h Change</th>
              <th class="pb-4">Volume</th>
              <th class="pb-4 pr-4 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-secondary-100 dark:divide-secondary-800">
            <tr v-for="market in marketData" :key="market.symbol" class="group hover:bg-surface-hover transition-colors">
              <td class="py-4 pl-4">
                <div class="flex items-center">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold text-white mr-3 shadow-md"
                       :class="getCurrencyColor(market.symbol.split('/')[0])">
                    {{ getCurrencySymbol(market.symbol.split('/')[0]) }}
                  </div>
                  <div>
                    <div class="font-bold text-base group-hover:text-primary transition-colors">{{ market.symbol }}</div>
                    <div class="text-xs text-muted">{{ getCurrencyName(market.symbol.split('/')[0]) }}</div>
                  </div>
                </div>
              </td>
              <td class="py-4 font-medium text-base">${{ market.price?.toLocaleString() }}</td>
              <td class="py-4 font-medium" :class="market.change >= 0 ? 'text-success' : 'text-error'">
                {{ market.change >= 0 ? '+' : '' }}{{ market.change }}%
              </td>
              <td class="py-4 text-sm text-muted">${{ market.volume?.toLocaleString() }}</td>
              <td class="py-4 pr-4 text-right">
                <button @click="quickTrade(market)" 
                        class="px-4 py-2 rounded-lg bg-primary-50 dark:bg-primary-900/20 text-primary font-medium text-sm hover:bg-primary-100 dark:hover:bg-primary-900/40 transition-colors">
                  Trader
                </button>
              </td>
            </tr>
          </tbody>
        </table>
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