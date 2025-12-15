<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">
          Bonjour, {{ userName }} 👋
        </h1>
        <p class="text-white/60">
          {{ new Date().toLocaleDateString('fr-FR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
        </p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="stat-card stat-card-blue">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-blue-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
              </svg>
            </div>
            <span class="badge badge-success">+5.2%</span>
          </div>
          <p class="text-white/60 text-sm mb-1">Solde Total</p>
          <p class="text-2xl font-bold text-white">{{ formatMoney(stats.totalBalance) }}</p>
        </div>

        <div class="stat-card stat-card-green">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500/20 flex items-center justify-center">
              <span class="text-xl">₿</span>
            </div>
            <span class="badge badge-success">+12.8%</span>
          </div>
          <p class="text-white/60 text-sm mb-1">Crypto Portfolio</p>
          <p class="text-2xl font-bold text-white">{{ formatMoney(stats.cryptoBalance) }}</p>
        </div>

        <div class="stat-card stat-card-purple">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-purple-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
              </svg>
            </div>
          </div>
          <p class="text-white/60 text-sm mb-1">Cartes Actives</p>
          <p class="text-2xl font-bold text-white">{{ stats.activeCards }}</p>
          <p class="text-sm text-white/40">Solde: {{ formatMoney(stats.cardsBalance) }}</p>
        </div>

        <div class="stat-card stat-card-orange">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-orange-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
              </svg>
            </div>
          </div>
          <p class="text-white/60 text-sm mb-1">Transferts ce mois</p>
          <p class="text-2xl font-bold text-white">{{ stats.monthlyTransfers }}</p>
          <p class="text-sm text-white/40">Volume: {{ formatMoney(stats.monthlyVolume) }}</p>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="glass-card p-6 mb-8">
        <h3 class="text-lg font-semibold text-white mb-6">🚀 Actions Rapides</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <NuxtLink to="/exchange/crypto" class="quick-action-btn">
            <span class="text-3xl mb-2">₿</span>
            <span class="text-sm font-medium">Acheter Crypto</span>
          </NuxtLink>

          <NuxtLink to="/exchange/fiat" class="quick-action-btn">
            <span class="text-3xl mb-2">💱</span>
            <span class="text-sm font-medium">Convertir</span>
          </NuxtLink>

          <NuxtLink to="/cards" class="quick-action-btn">
            <span class="text-3xl mb-2">💳</span>
            <span class="text-sm font-medium">Mes Cartes</span>
          </NuxtLink>

          <NuxtLink to="/transfer" class="quick-action-btn">
            <span class="text-3xl mb-2">💸</span>
            <span class="text-sm font-medium">Envoyer</span>
          </NuxtLink>

          <NuxtLink to="/wallet" class="quick-action-btn">
            <span class="text-3xl mb-2">👛</span>
            <span class="text-sm font-medium">Portefeuilles</span>
          </NuxtLink>

          <button @click="showQRModal = true" class="quick-action-btn">
            <span class="text-3xl mb-2">📱</span>
            <span class="text-sm font-medium">Recevoir</span>
          </button>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        <!-- Crypto Markets -->
        <div class="glass-card p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-white">📊 Marchés Crypto</h3>
            <NuxtLink to="/exchange/crypto" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
              Voir tout →
            </NuxtLink>
          </div>
          
          <div class="space-y-4">
            <div v-for="crypto in cryptoMarkets" :key="crypto.symbol" 
                class="flex items-center justify-between p-4 rounded-xl bg-white/5 hover:bg-white/10 transition-colors cursor-pointer">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-xl flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold">{{ crypto.symbol.slice(0, 2) }}</span>
                </div>
                <div>
                  <p class="font-semibold text-white">{{ crypto.name }}</p>
                  <p class="text-sm text-white/50">{{ crypto.symbol }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-white">${{ crypto.price.toLocaleString() }}</p>
                <p class="text-sm" :class="crypto.change >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ crypto.change >= 0 ? '+' : '' }}{{ crypto.change.toFixed(2) }}%
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="glass-card p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-white">🕒 Activité Récente</h3>
            <NuxtLink to="/transactions" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
              Voir tout →
            </NuxtLink>
          </div>
          
          <div class="space-y-4">
            <div v-for="activity in recentActivities" :key="activity.id" 
                class="flex items-center gap-4 p-4 rounded-xl bg-white/5">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center" :class="activity.bgColor">
                <span class="text-lg">{{ activity.icon }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <p class="font-medium text-white truncate">{{ activity.title }}</p>
                <p class="text-sm text-white/50 truncate">{{ activity.description }}</p>
              </div>
              <div class="text-right">
                <p class="font-semibold text-white">{{ activity.amount }}</p>
                <p class="text-xs text-white/40">{{ formatTime(activity.time) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- My Cards Section -->
      <div class="glass-card p-6 mb-8">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-white">💳 Mes Cartes</h3>
          <NuxtLink to="/cards" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
            Gérer →
          </NuxtLink>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Virtual Card -->
          <div class="credit-card credit-card-virtual">
            <div class="flex justify-between items-start mb-8">
              <span class="text-white/80 font-medium">Carte Virtuelle</span>
              <span class="text-2xl">💳</span>
            </div>
            <p class="text-xl font-mono text-white tracking-wider mb-4">•••• •••• •••• 4521</p>
            <div class="flex justify-between items-end">
              <div>
                <p class="text-xs text-white/60">Solde</p>
                <p class="text-lg font-bold text-white">$2,450.00</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-white/60">Expire</p>
                <p class="text-white font-medium">12/26</p>
              </div>
            </div>
          </div>

          <!-- Physical Card -->
          <div class="credit-card credit-card-physical">
            <div class="flex justify-between items-start mb-8">
              <span class="text-white/80 font-medium">Carte Physique</span>
              <span class="text-2xl">🔒</span>
            </div>
            <p class="text-xl font-mono text-white tracking-wider mb-4">•••• •••• •••• 8732</p>
            <div class="flex justify-between items-end">
              <div>
                <p class="text-xs text-white/60">Solde</p>
                <p class="text-lg font-bold text-white">$800.50</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-white/60">Expire</p>
                <p class="text-white font-medium">08/27</p>
              </div>
            </div>
          </div>

          <!-- Add Card -->
          <NuxtLink to="/cards/new" class="credit-card flex items-center justify-center border-2 border-dashed border-white/20 hover:border-indigo-500 transition-colors" style="background: transparent;">
            <div class="text-center">
              <div class="w-16 h-16 rounded-full bg-white/10 flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-white/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </div>
              <p class="text-white/60 font-medium">Nouvelle Carte</p>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- Exchange Rates -->
      <div class="glass-card p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-white">💱 Taux de Change</h3>
          <button @click="refreshRates" class="flex items-center gap-2 text-indigo-400 hover:text-indigo-300 text-sm font-medium">
            <svg class="w-4 h-4" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Actualiser
          </button>
        </div>
        
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div v-for="rate in fiatRates" :key="rate.pair" 
              class="p-4 rounded-xl bg-white/5 text-center hover:bg-white/10 transition-colors">
            <p class="text-sm font-medium text-white/60 mb-1">{{ rate.pair }}</p>
            <p class="text-lg font-bold text-white">{{ rate.rate.toFixed(4) }}</p>
            <p class="text-xs" :class="rate.change >= 0 ? 'text-emerald-400' : 'text-red-400'">
              {{ rate.change >= 0 ? '+' : '' }}{{ (rate.change * 100).toFixed(2) }}%
            </p>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { dashboardAPI, exchangeAPI } from '~/composables/useApi'

const authStore = useAuthStore()

const userName = computed(() => authStore.user?.firstName || 'Utilisateur')

const loading = ref(true)
const refreshing = ref(false)
const showQRModal = ref(false)

// Stats
const stats = ref({
  totalBalance: 25750.80,
  cryptoBalance: 12500.00,
  cardsBalance: 3250.50,
  activeCards: 4,
  monthlyTransfers: 23,
  monthlyVolume: 15750.00
})

// Market data
const cryptoMarkets = ref([
  { symbol: 'BTC', name: 'Bitcoin', price: 45125, change: 2.35, bgColor: 'bg-orange-500' },
  { symbol: 'ETH', name: 'Ethereum', price: 3024, change: 1.87, bgColor: 'bg-blue-500' },
  { symbol: 'BNB', name: 'BNB', price: 245, change: -0.52, bgColor: 'bg-yellow-500' },
  { symbol: 'SOL', name: 'Solana', price: 98, change: 5.21, bgColor: 'bg-purple-500' },
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
  { id: 1, icon: '₿', title: 'Achat Bitcoin', description: '0.02 BTC', amount: '$900.50', time: new Date(Date.now() - 2 * 60 * 60 * 1000), bgColor: 'bg-orange-500' },
  { id: 2, icon: '💳', title: 'Rechargement carte', description: 'Carte virtuelle', amount: '$200.00', time: new Date(Date.now() - 4 * 60 * 60 * 1000), bgColor: 'bg-blue-500' },
  { id: 3, icon: '💱', title: 'Conversion USD → EUR', description: '1000 USD', amount: '€845.60', time: new Date(Date.now() - 6 * 60 * 60 * 1000), bgColor: 'bg-green-500' },
  { id: 4, icon: '💸', title: 'Transfert', description: 'Mobile Money', amount: '$50.00', time: new Date(Date.now() - 12 * 60 * 60 * 1000), bgColor: 'bg-purple-500' },
])

// Methods
const formatMoney = (amount) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount)
}

const formatTime = (date) => {
  return new Date(date).toLocaleString('fr-FR', { hour: '2-digit', minute: '2-digit', day: '2-digit', month: '2-digit' })
}

const refreshRates = async () => {
  refreshing.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    fiatRates.value.forEach(rate => {
      rate.change = (Math.random() - 0.5) * 0.02
      rate.rate = rate.rate * (1 + rate.change)
    })
  } finally {
    refreshing.value = false
  }
}

const fetchDashboardData = async () => {
  try {
    // Try to fetch from API
    const [summaryRes, activityRes] = await Promise.all([
      dashboardAPI.getSummary().catch(() => null),
      dashboardAPI.getRecentActivity().catch(() => null)
    ])
    
    if (summaryRes?.data) {
      stats.value = summaryRes.data
    }
    if (activityRes?.data) {
      recentActivities.value = activityRes.data
    }
  } catch (error) {
    console.log('Using mock data')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.quick-action-btn {
  @apply flex flex-col items-center justify-center p-4 rounded-xl transition-all duration-300;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: white;
}

.quick-action-btn:hover {
  background: rgba(99, 102, 241, 0.15);
  border-color: rgba(99, 102, 241, 0.3);
  transform: translateY(-2px);
}
</style>