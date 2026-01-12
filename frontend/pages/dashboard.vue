<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Bonjour, {{ userName }} 👋
        </h1>
        <p class="text-gray-500 dark:text-slate-400">
          {{ new Date().toLocaleDateString('fr-FR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
        </p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="stat-card stat-card-blue">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-blue-500/30 flex items-center justify-center">
              <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
              </svg>
            </div>
            <span class="badge badge-success">+5.2%</span>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm mb-1">Solde Total</p>
          <div v-if="loading && walletStore.totalBalance === 0" class="h-8 w-32 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse"></div>
          <p v-else class="text-2xl font-bold text-gray-900 dark:text-white">{{ formatMoney(walletStore.totalBalance) }}</p>
        </div>

        <div class="stat-card stat-card-green">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500/30 flex items-center justify-center">
              <span class="text-xl text-emerald-400">₿</span>
            </div>
            <span class="badge badge-success">+12.8%</span>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm mb-1">Crypto Portfolio</p>
          <div v-if="loading && walletStore.cryptoBalance === 0" class="h-8 w-32 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse"></div>
          <p v-else class="text-2xl font-bold text-gray-900 dark:text-white">{{ formatMoney(walletStore.cryptoBalance) }}</p>
        </div>

        <div class="stat-card stat-card-purple">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-purple-500/30 flex items-center justify-center">
              <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
              </svg>
            </div>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm mb-1">Cartes Actives</p>
          <div v-if="loading" class="h-8 w-16 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse mb-1"></div>
          <p v-else class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.activeCards }}</p>
          <div v-if="loading" class="h-4 w-24 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse"></div>
          <p v-else class="text-sm text-gray-400 dark:text-slate-500">Solde: {{ formatMoney(stats.cardsBalance) }}</p>
        </div>

        <div class="stat-card stat-card-orange">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-orange-500/30 flex items-center justify-center">
              <svg class="w-6 h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
              </svg>
            </div>
          </div>
          <p class="text-gray-500 dark:text-slate-400 text-sm mb-1">Transferts ce mois</p>
          <div v-if="loading" class="h-8 w-16 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse mb-1"></div>
          <p v-else class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.monthlyTransfers }}</p>
          <div v-if="loading" class="h-4 w-24 bg-gray-200 dark:bg-slate-700 rounded-md animate-pulse"></div>
          <p v-else class="text-sm text-gray-400 dark:text-slate-500">Volume: {{ formatMoney(stats.monthlyVolume) }}</p>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="glass-card p-6 mb-8">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">🚀 Actions Rapides</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <NuxtLink to="/scan" class="quick-action-btn">
            <span class="text-3xl mb-2">📷</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Scanner</span>
          </NuxtLink>

          <NuxtLink to="/transfer" class="quick-action-btn">
            <span class="text-3xl mb-2">💸</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Envoyer</span>
          </NuxtLink>

          <NuxtLink to="/cards" class="quick-action-btn">
            <span class="text-3xl mb-2">💳</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Mes Cartes</span>
          </NuxtLink>

          <NuxtLink to="/wallet" class="quick-action-btn">
            <span class="text-3xl mb-2">👛</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Portefeuilles</span>
          </NuxtLink>

          <NuxtLink to="/exchange/crypto" class="quick-action-btn">
            <span class="text-3xl mb-2">₿</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Acheter Crypto</span>
          </NuxtLink>

          <NuxtLink to="/exchange/fiat" class="quick-action-btn">
            <span class="text-3xl mb-2">💱</span>
            <span class="text-sm font-medium text-gray-600 dark:text-slate-300">Convertir</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        <!-- Crypto Markets -->
        <div class="glass-card p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">📊 Marchés Crypto</h3>
            <NuxtLink to="/exchange/crypto" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
              Voir tout →
            </NuxtLink>
          </div>
          
          <div class="space-y-4">
            <div v-for="crypto in cryptoMarkets" :key="crypto.symbol" 
                class="flex items-center justify-between p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 hover:bg-gray-100 dark:hover:bg-slate-700/50 transition-colors cursor-pointer border border-gray-100 dark:border-slate-700/50">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-xl flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold">{{ crypto.symbol?.slice(0, 2) || '??' }}</span>
                </div>
                <div>
                  <p class="font-semibold text-gray-900 dark:text-white">{{ crypto.name }}</p>
                  <p class="text-sm text-gray-500 dark:text-slate-400">{{ crypto.symbol }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900 dark:text-white">${{ (crypto.price || 0).toLocaleString() }}</p>
                <p class="text-sm" :class="(crypto.change || 0) >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ (crypto.change || 0) >= 0 ? '+' : '' }}{{ (crypto.change || 0).toFixed(2) }}%
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="glass-card p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">🕒 Activité Récente</h3>
            <NuxtLink to="/transactions" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
              Voir tout →
            </NuxtLink>
          </div>
          
          <div class="space-y-4">
            <div v-for="activity in recentActivities" :key="activity.id" 
                class="flex items-center gap-4 p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 border border-gray-100 dark:border-slate-700/50">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center" :class="activity.bgColor">
                <span class="text-lg">{{ activity.icon }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <p class="font-medium text-gray-900 dark:text-white truncate">{{ activity.title }}</p>
                <p class="text-sm text-gray-500 dark:text-slate-400 truncate">{{ activity.description }}</p>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900 dark:text-white">{{ activity.amount }}</p>
                <p class="text-xs text-gray-400 dark:text-slate-500">{{ formatTime(activity.time) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- My Cards Section -->
      <div class="glass-card p-6 mb-8">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">💳 Mes Cartes</h3>
          <NuxtLink to="/cards" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium">
            Gérer →
          </NuxtLink>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Dynamic Cards from API -->
          <div v-for="card in userCards" :key="card.id" class="credit-card" :class="card.is_virtual ? 'credit-card-virtual' : 'credit-card-physical'">
            <div class="flex justify-between items-start mb-8">
              <span class="text-white/70 font-medium">{{ card.is_virtual ? 'Carte Virtuelle' : 'Carte Physique' }}</span>
              <span class="text-2xl">{{ card.is_virtual ? '💳' : '🔒' }}</span>
            </div>
            <p class="text-xl font-mono text-white tracking-wider mb-4">{{ card.card_number || '•••• •••• •••• ••••' }}</p>
            <div class="flex justify-between items-end">
              <div>
                <p class="text-xs text-white/60">Solde</p>
                <p class="text-lg font-bold text-white">{{ formatMoney(card.balance || 0) }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-white/60">Expire</p>
                <p class="text-white font-medium">{{ card.expiry_month }}/{{ card.expiry_year }}</p>
              </div>
            </div>
          </div>

          <!-- No cards message -->
          <div v-if="userCards.length === 0" class="col-span-full text-center py-8 text-slate-400">
            <p class="mb-4">Vous n'avez pas encore de carte.</p>
            <NuxtLink to="/cards/new" class="text-indigo-400 hover:text-indigo-300 font-medium">
              Créer une carte →
            </NuxtLink>
          </div>

          <!-- Add Card -->
          <NuxtLink v-if="userCards.length > 0" to="/cards/new" class="credit-card flex items-center justify-center border-2 border-dashed border-slate-300 dark:border-slate-600 hover:border-indigo-500 transition-colors bg-gray-50 dark:bg-slate-800/30">
            <div class="text-center">
              <div class="w-16 h-16 rounded-full bg-gray-200 dark:bg-slate-700/50 flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-gray-400 dark:text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </div>
              <p class="text-gray-500 dark:text-slate-400 font-medium">Nouvelle Carte</p>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- Exchange Rates -->
      <div class="glass-card p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">💱 Taux de Change</h3>
          <button @click="refreshRates" class="flex items-center gap-2 text-indigo-400 hover:text-indigo-300 text-sm font-medium">
            <svg class="w-4 h-4" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Actualiser
          </button>
        </div>
        
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div v-for="rate in fiatRates" :key="rate.pair" 
              class="p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 text-center hover:bg-gray-100 dark:hover:bg-slate-700/50 transition-colors border border-gray-100 dark:border-slate-700/50">
            <p class="text-sm font-medium text-gray-500 dark:text-slate-400 mb-1">{{ rate.pair }}</p>
            <p class="text-lg font-bold text-gray-900 dark:text-white">{{ rate.rate.toFixed(4) }}</p>
            <p class="text-xs" :class="rate.change >= 0 ? 'text-emerald-500 dark:text-emerald-400' : 'text-red-500 dark:text-red-400'">
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
import { useWalletStore } from '~/stores/wallet'
import { dashboardAPI, exchangeAPI, cardAPI, transferAPI } from '~/composables/useApi'

const authStore = useAuthStore()
const walletStore = useWalletStore()

const userName = computed(() => authStore.user?.first_name || 'Utilisateur')

// Loading state is now a mix of local data fetching and store loading
// If store has data, we don't show full page loading
const loading = ref(true)
const refreshing = ref(false)
const showQRModal = ref(false)

// Stats - Cards/Transfers still loaded locally for now
const stats = ref({
  activeCards: 0,
  cardsBalance: 0,
  monthlyTransfers: 0,
  monthlyVolume: 0
})

// Market data
const cryptoMarkets = ref([])

// User cards
const userCards = ref([])

// Fiat rates
const fiatRates = ref([])

// Recent activities
const recentActivities = ref([])

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
    const response = await exchangeAPI.getRates()
    if (response.data) {
       fiatRates.value = Object.entries(response.data).map(([pair, rate]) => ({
         pair,
         rate: rate.Rate,
         change: rate.Change24h || 0
       })).filter(r => r.pair.includes('USD') || r.pair.includes('EUR')).slice(0, 6)
    }
  } catch (e) {
    console.error('Failed to refresh rates:', e)
  } finally {
    refreshing.value = false
  }
}

const fetchDashboardData = async () => {
  // Initialize wallet store (loads from cache)
  walletStore.initialize()
  
  // If we have cached wallet data, don't show main skeleton
  if (walletStore.wallets.length > 0) {
    loading.value = false
  } else {
    loading.value = true
  }

  // Ensure user profile is loaded
  if (!authStore.user?.first_name) {
    try {
        await authStore.fetchUserProfile()
    } catch (e) {
        console.error('Failed to fetch user profile:', e)
    }
  }

  try {
    // Trigger wallet fetch in background (store handles logic)
    walletStore.fetchWallets()

    // Fetch other data in parallel
    const [summaryRes, activityRes, marketsRes, cardsRes, ratesRes, transfersRes] = await Promise.all([
      dashboardAPI.getSummary().catch(() => null),
      dashboardAPI.getRecentActivity().catch(() => ({ data: [] })),
      exchangeAPI.getMarkets().catch(() => ({ data: [] })),
      cardAPI.getAll().catch(() => ({ data: { cards: [] } })),
      exchangeAPI.getRates().catch(() => ({ data: {} })),
      transferAPI.getAll(100, 0).catch(() => ({ data: { transfers: [] } }))
    ])
    
    // Fallback to summary stats if wallet store empty/fails (rare)
    if (summaryRes?.data && walletStore.totalBalance === 0) {
      // We prefer store calculation, but this is a backup
    }
    
    // Load cards stats
    if (cardsRes?.data?.cards) {
      userCards.value = cardsRes.data.cards.slice(0, 3)
      stats.value.activeCards = cardsRes.data.cards.filter(c => c.status === 'active').length
      stats.value.cardsBalance = cardsRes.data.cards.reduce((sum, c) => sum + (Number(c.balance) || 0), 0)
    }
    
    if (activityRes?.data?.activities) recentActivities.value = activityRes.data.activities.map(a => ({
      ...a,
      bgColor: a.type === 'credit' ? 'bg-green-500/20' : 'bg-red-500/20',
      icon: a.type === 'credit' ? '↓' : '↑'
    }))
    
    if (marketsRes?.data?.markets) {
       cryptoMarkets.value = marketsRes.data.markets.slice(0, 4).map(m => ({
         name: m.BaseAsset,
         symbol: m.Symbol,
         price: m.Price,
         change: m.Change24h,
         bgColor: 'bg-orange-500/20' 
       }))
    }
    
    if (ratesRes?.data?.rates) {
       fiatRates.value = Object.entries(ratesRes.data.rates).map(([pair, rate]) => ({
         pair,
         rate: rate.Rate || rate,
         change: rate.Change24h || 0
       })).filter(r => r.pair.includes('USD') || r.pair.includes('EUR')).slice(0, 6)
    }

    // Calculate monthly transfers stats
    if (transfersRes?.data?.transfers) {
      // Simple calculation
      stats.value.monthlyTransfers = transfersRes.data.transfers.length
      // Volume calculation could be complex, omitting for now or assuming pre-calculated
      // Use summaryRes if available for perfect volume
      if (summaryRes?.data) {
          stats.value.monthlyVolume = summaryRes.data.monthlyVolume
          stats.value.monthlyTransfers = summaryRes.data.monthlyTransfers
      }
    }

  } catch (error) {
    console.error('Error fetching dashboard data:', error)
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
  /* Light Mode Default - Softer */
  background: rgba(241, 245, 249, 0.5); /* Slate-100/50% */
  border: 1px solid rgba(226, 232, 240, 0.8); /* Slate-200/80% */
  color: #475569;
}

.dark .quick-action-btn {
  /* Dark Mode Overrides */
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: white;
}

.quick-action-btn:hover {
  background: #ffffff;
  border-color: #c7d2fe; /* indigo-200 */
  color: #4f46e5; /* indigo-600 */
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.1);
}

.dark .quick-action-btn:hover {
  background: rgba(99, 102, 241, 0.15);
  border-color: rgba(99, 102, 241, 0.3);
  color: white;
  box-shadow: none;
}
</style>