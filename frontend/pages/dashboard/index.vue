<template>
  <div>
    <!-- Header -->
    <div class="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4 animate-fade-in-up">
      <div>
        <h1 class="text-3xl font-extrabold text-gray-900 dark:text-gray-100">Bonjour, {{ userName }} 👋</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Voici un aperçu de votre portefeuille aujourd'hui.</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="px-4 py-2 rounded-xl bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 text-sm font-medium text-gray-600 dark:text-gray-300 shadow-sm">
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
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none transition-all duration-300 hover:transform hover:-translate-y-1 backdrop-blur-sm">
        <div class="absolute top-0 right-0 w-32 h-32 bg-indigo-500/5 dark:bg-indigo-500/20 rounded-full blur-3xl group-hover:bg-indigo-500/10 dark:group-hover:bg-indigo-500/30 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Valeur Totale</p>
            <div class="p-2.5 bg-indigo-50 dark:bg-indigo-500/20 rounded-xl text-indigo-600 dark:text-indigo-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-2 tracking-tight">${{ portfolioStats.totalValue?.toLocaleString() || '0' }}</p>
          <div class="flex items-center text-sm">
            <span :class="portfolioStats.dailyChange >= 0 ? 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-500/20' : 'text-rose-600 bg-rose-50 dark:text-rose-400 dark:bg-rose-500/20'" class="px-2 py-0.5 rounded-full font-bold flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="portfolioStats.dailyChange >= 0" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"/>
              </svg>
              {{ Math.abs(portfolioStats.dailyChange) }}%
            </span>
            <span class="text-gray-400 ml-2">vs hier</span>
          </div>
        </div>
      </div>

      <!-- Available Cash -->
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none transition-all duration-300 hover:transform hover:-translate-y-1 backdrop-blur-sm">
        <div class="absolute top-0 right-0 w-32 h-32 bg-emerald-500/5 dark:bg-emerald-500/20 rounded-full blur-3xl group-hover:bg-emerald-500/10 dark:group-hover:bg-emerald-500/30 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Disponible</p>
            <div class="p-2.5 bg-emerald-50 dark:bg-emerald-500/20 rounded-xl text-emerald-600 dark:text-emerald-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-2 tracking-tight">${{ portfolioStats.availableCash?.toLocaleString() || '0' }}</p>
          <p class="text-sm text-gray-400">Prêt à investir</p>
        </div>
      </div>

      <!-- Total Profit/Loss -->
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none transition-all duration-300 hover:transform hover:-translate-y-1 backdrop-blur-sm">
        <div class="absolute top-0 right-0 w-32 h-32 bg-indigo-500/5 dark:bg-indigo-500/20 rounded-full blur-3xl group-hover:bg-indigo-500/10 dark:group-hover:bg-indigo-500/30 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">P&L Total</p>
            <div class="p-2.5 bg-indigo-50 dark:bg-indigo-500/20 rounded-xl text-indigo-600 dark:text-indigo-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 00-2-2z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-2 tracking-tight" :class="portfolioStats.totalPnL >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'">
            {{ portfolioStats.totalPnL >= 0 ? '+' : '' }}${{ Math.abs(portfolioStats.totalPnL || 0).toLocaleString() }}
          </p>
          <p class="text-sm text-gray-400">Total réalisé</p>
        </div>
      </div>

      <!-- Active Orders -->
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none transition-all duration-300 hover:transform hover:-translate-y-1 backdrop-blur-sm">
        <div class="absolute top-0 right-0 w-32 h-32 bg-amber-500/5 dark:bg-amber-500/20 rounded-full blur-3xl group-hover:bg-amber-500/10 dark:group-hover:bg-amber-500/30 transition-all"></div>
        <div class="relative">
          <div class="flex items-center justify-between mb-4">
            <p class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Ordres Actifs</p>
            <div class="p-2.5 bg-amber-50 dark:bg-amber-500/20 rounded-xl text-amber-600 dark:text-amber-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-3xl font-extrabold text-gray-900 dark:text-gray-100 mb-2 tracking-tight">{{ portfolioStats.activeOrders || 0 }}</p>
          <p class="text-sm text-gray-400">En attente</p>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="mb-8 animate-fade-in-up" style="animation-delay: 0.15s">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
        <span>⚡</span> Actions Rapides
      </h2>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <!-- Transfer -->
        <NuxtLink to="/transfer" class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-md hover:shadow-xl hover:bg-indigo-50/50 dark:hover:bg-slate-800 transition-all duration-300 flex flex-col items-center gap-3 backdrop-blur-sm">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-indigo-600 flex items-center justify-center text-white shadow-lg shadow-indigo-500/30 group-hover:scale-110 transition-transform duration-300">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
            </svg>
          </div>
          <span class="font-bold text-gray-700 dark:text-gray-100 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">Virement</span>
        </NuxtLink>

        <!-- Exchange -->
        <NuxtLink to="/exchange/crypto" class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-md hover:shadow-xl hover:bg-violet-50/50 dark:hover:bg-slate-800 transition-all duration-300 flex flex-col items-center gap-3 backdrop-blur-sm">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-violet-500 to-violet-600 flex items-center justify-center text-white shadow-lg shadow-violet-500/30 group-hover:scale-110 transition-transform duration-300">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
            </svg>
          </div>
          <span class="font-bold text-gray-700 dark:text-gray-100 group-hover:text-violet-600 dark:group-hover:text-violet-400 transition-colors">Échange</span>
        </NuxtLink>

        <!-- Wallet -->
        <NuxtLink to="/wallet" class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-md hover:shadow-xl hover:bg-emerald-50/50 dark:hover:bg-slate-800 transition-all duration-300 flex flex-col items-center gap-3 backdrop-blur-sm">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center text-white shadow-lg shadow-emerald-500/30 group-hover:scale-110 transition-transform duration-300">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
            </svg>
          </div>
          <span class="font-bold text-gray-700 dark:text-gray-100 group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">Portefeuille</span>
        </NuxtLink>

        <!-- Cards -->
        <NuxtLink to="/cards" class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 shadow-md hover:shadow-xl hover:bg-amber-50/50 dark:hover:bg-slate-800 transition-all duration-300 flex flex-col items-center gap-3 backdrop-blur-sm">
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-amber-500 to-amber-600 flex items-center justify-center text-white shadow-lg shadow-amber-500/30 group-hover:scale-110 transition-transform duration-300">
            <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
            </svg>
          </div>
          <span class="font-bold text-gray-700 dark:text-gray-100 group-hover:text-amber-600 dark:group-hover:text-amber-400 transition-colors">Mes Cartes</span>
        </NuxtLink>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8 animate-fade-in-up" style="animation-delay: 0.2s">
      
      <!-- Portfolio Chart -->
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/70 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none backdrop-blur-sm">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white">Performance</h3>
          <div class="flex space-x-1 bg-gray-100 dark:bg-slate-800 p-1 rounded-lg">
            <button v-for="period in chartPeriods" :key="period.value"
                    @click="selectedPeriod = period.value"
                    :class="selectedPeriod === period.value 
                      ? 'bg-white dark:bg-slate-700 text-indigo-600 dark:text-indigo-400 shadow-sm' 
                      : 'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'"
                    class="px-3 py-1 text-sm font-medium rounded-md transition-all">
              {{ period.label }}
            </button>
          </div>
        </div>
        <div class="h-64 flex items-center justify-center border border-dashed border-gray-200 dark:border-white/10 rounded-xl bg-gray-50 dark:bg-slate-800/60">
          <span class="text-gray-400 dark:text-gray-500">Graphique Portefeuille (Chart.js)</span>
        </div>
      </div>

      <!-- Top Holdings -->
      <div class="relative overflow-hidden group rounded-2xl p-6 bg-white dark:bg-slate-900/70 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none backdrop-blur-sm">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white">Top Actifs</h3>
          <NuxtLink to="/wallet" class="text-indigo-600 dark:text-indigo-400 text-sm font-bold hover:text-indigo-700 dark:hover:text-indigo-300 transition-colors">
            Voir tout
          </NuxtLink>
        </div>
        <div class="space-y-3">
          <div v-for="holding in topHoldings" :key="holding.currency" 
               class="flex items-center justify-between p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors cursor-pointer group">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shadow-md transform group-hover:scale-110 transition-transform"
                   :class="getCurrencyColor(holding.currency)">
                {{ getCurrencySymbol(holding.currency) }}
              </div>
              <div>
                <h4 class="font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">{{ holding.currency }}</h4>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ holding.amount }} {{ holding.currency }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900 dark:text-white">${{ holding.value?.toLocaleString() }}</p>
              <p class="text-sm font-bold" :class="holding.change >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'">
                {{ holding.change >= 0 ? '+' : '' }}{{ holding.change }}%
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Market Overview -->
    <div class="relative overflow-hidden mb-8 animate-fade-in-up p-6 bg-white dark:bg-slate-900/70 border border-gray-200 dark:border-white/10 shadow-lg dark:shadow-none backdrop-blur-sm" style="animation-delay: 0.3s">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white">Marché</h3>
        <NuxtLink to="/exchange/crypto" class="text-indigo-600 dark:text-indigo-400 text-sm font-bold hover:text-indigo-700 dark:hover:text-indigo-300 transition-colors">
          Voir tout les marchés
        </NuxtLink>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider border-b border-gray-200 dark:border-white/10">
              <th class="pb-4 pl-4">Actif</th>
              <th class="pb-4">Prix</th>
              <th class="pb-4">24h Change</th>
              <th class="pb-4">Volume</th>
              <th class="pb-4 pr-4 text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <tr v-for="market in marketData" :key="market.symbol" class="group hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
              <td class="py-4 pl-4">
                <div class="flex items-center">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold text-white mr-3 shadow-md"
                       :class="getCurrencyColor(market.symbol.split('/')[0])">
                    {{ getCurrencySymbol(market.symbol.split('/')[0]) }}
                  </div>
                  <div>
                    <div class="font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">{{ market.symbol }}</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400">{{ getCurrencyName(market.symbol.split('/')[0]) }}</div>
                  </div>
                </div>
              </td>
              <td class="py-4 font-bold text-gray-900 dark:text-white">${{ market.price?.toLocaleString() }}</td>
              <td class="py-4 font-bold" :class="market.change >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'">
                {{ market.change >= 0 ? '+' : '' }}{{ market.change }}%
              </td>
              <td class="py-4 text-sm text-gray-500 dark:text-gray-400">${{ market.volume?.toLocaleString() }}</td>
              <td class="py-4 pr-4 text-right">
                <button @click="quickTrade(market)" 
                        class="px-4 py-2 rounded-lg bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 font-bold text-sm hover:bg-indigo-100 dark:hover:bg-indigo-500/20 transition-colors">
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
import { useAuthStore } from '~/stores/auth'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const authStore = useAuthStore()
const userName = computed(() => authStore.user?.first_name || 'Utilisateur')

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
    'BTC': 'bg-amber-500',
    'ETH': 'bg-indigo-500',
    'LTC': 'bg-gray-400',
    'ADA': 'bg-blue-600',
    'USD': 'bg-emerald-500',
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

const quickTrade = (market) => {
  navigateTo(`/exchange/trading?pair=${market.symbol}`)
}

// Fetch data functions
const fetchDashboardData = async () => {
    // Demo data for now
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
}

const fetchMarketData = async () => {
    // Demo market data
    marketData.value = [
      { symbol: 'BTC/USD', price: 43500, change: 2.3, volume: 1234567 },
      { symbol: 'ETH/USD', price: 2450, change: -1.2, volume: 987654 },
      { symbol: 'LTC/USD', price: 72, change: 0.8, volume: 543210 },
      { symbol: 'ADA/USD', price: 0.52, change: 3.1, volume: 876543 }
    ]
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchDashboardData(),
    fetchMarketData()
  ])
})
</script>