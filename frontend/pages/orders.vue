<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto animate-fade-in-up">
      <!-- Header Section -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <div class="flex items-center gap-3 mb-2">
            <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white">📋 Mes Ordres</h1>
            <span class="px-3 py-1 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 text-sm rounded-full font-medium">
              {{ allOrders.length }} total
            </span>
          </div>
          <p class="text-gray-500 dark:text-gray-400">Gérez vos ordres de trading et consultez vos factures</p>
        </div>

        <div class="flex gap-3">
          <button 
            @click="exportOrders"
            class="btn-secondary flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
            </svg>
            Exporter
          </button>
          <NuxtLink to="/exchange/trading" class="btn-premium flex items-center gap-2 shadow-lg shadow-indigo-500/20">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Nouvel Ordre
          </NuxtLink>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Total Ordres</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total }}</p>
            </div>
            <div class="w-12 h-12 rounded-xl bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center">
              <span class="text-2xl">📋</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">En attente</p>
              <p class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ stats.pending }}</p>
            </div>
            <div class="w-12 h-12 rounded-xl bg-yellow-100 dark:bg-yellow-900/30 flex items-center justify-center">
              <span class="text-2xl">⏳</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Exécutés</p>
              <p class="text-2xl font-bold text-green-600 dark:text-green-400">{{ stats.filled }}</p>
            </div>
            <div class="w-12 h-12 rounded-xl bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
              <span class="text-2xl">✅</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Volume Total</p>
              <p class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">{{ formatVolume(stats.volume) }}</p>
            </div>
            <div class="w-12 h-12 rounded-xl bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center">
              <span class="text-2xl">📊</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="glass-card p-2 mb-6">
        <div class="flex flex-wrap gap-2">
          <button 
            v-for="filter in filters" 
            :key="filter.key"
            @click="activeFilter = filter.key"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200',
              activeFilter === filter.key 
                ? 'bg-indigo-600 text-white shadow-md' 
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
            ]"
          >
            {{ filter.label }}
            <span 
              v-if="filter.count > 0" 
              :class="[
                'ml-2 px-2 py-0.5 text-xs rounded-full',
                activeFilter === filter.key ? 'bg-white/20' : 'bg-gray-200 dark:bg-gray-700'
              ]"
            >
              {{ filter.count }}
            </span>
          </button>
        </div>
      </div>

      <!-- Search and Filters Bar -->
      <div class="glass-card p-4 mb-6">
        <div class="flex flex-col md:flex-row gap-4">
          <div class="flex-1 relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Rechercher par paire, ID ou statut..."
              class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
          </div>
          <div class="flex gap-3">
            <select 
              v-model="pairFilter" 
              class="px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
            >
              <option value="">Toutes les paires</option>
              <option v-for="pair in availablePairs" :key="pair" :value="pair">{{ pair }}</option>
            </select>
            <select 
              v-model="typeFilter" 
              class="px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
            >
              <option value="">Tous types</option>
              <option value="market">Market</option>
              <option value="limit">Limit</option>
              <option value="stop_loss">Stop Loss</option>
            </select>
            <select 
              v-model="dateFilter" 
              class="px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
            >
              <option value="">Toutes dates</option>
              <option value="today">Aujourd'hui</option>
              <option value="week">Cette semaine</option>
              <option value="month">Ce mois</option>
              <option value="year">Cette année</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="glass-card p-8">
        <div class="space-y-4">
          <div v-for="i in 5" :key="i" class="flex items-center gap-4 animate-pulse">
            <div class="w-12 h-12 rounded-xl bg-gray-200 dark:bg-gray-700"></div>
            <div class="flex-1 space-y-2">
              <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/3"></div>
            </div>
            <div class="h-8 w-24 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredOrders.length === 0" class="glass-card p-12 text-center">
        <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center">
          <span class="text-4xl">📋</span>
        </div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun ordre trouvé</h3>
        <p class="text-gray-500 mb-6">Commencez à trader pour voir vos ordres ici</p>
        <NuxtLink to="/exchange/trading" class="btn-premium inline-flex items-center gap-2">
          <span>Créer un ordre</span>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
          </svg>
        </NuxtLink>
      </div>

      <!-- Orders Grid -->
      <div v-else class="space-y-4">
        <div 
          v-for="order in paginatedOrders" 
          :key="order.id"
          class="glass-card p-5 hover:border-indigo-500/30 transition-all duration-300 group"
        >
          <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
            <!-- Order Info -->
            <div class="flex items-center gap-4">
              <div 
                :class="[
                  'w-14 h-14 rounded-xl flex items-center justify-center text-2xl',
                  order.side === 'buy' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-red-100 dark:bg-red-900/30'
                ]"
              >
                {{ order.side === 'buy' ? '📈' : '📉' }}
              </div>
              <div>
                <div class="flex items-center gap-2 mb-1">
                  <span class="font-bold text-lg text-gray-900 dark:text-white">{{ order.pair }}</span>
                  <span 
                    :class="[
                      'px-2 py-0.5 text-xs rounded-full font-medium',
                      order.side === 'buy' 
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' 
                        : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                    ]"
                  >
                    {{ order.side === 'buy' ? 'ACHAT' : 'VENTE' }}
                  </span>
                </div>
                <div class="flex items-center gap-3 text-sm text-gray-500">
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    {{ formatDate(order.created_at) }}
                  </span>
                  <span>•</span>
                  <span class="capitalize">{{ order.order_type }}</span>
                  <span>•</span>
                  <span class="font-mono text-xs">ID: {{ order.id.slice(-8) }}</span>
                </div>
              </div>
            </div>

            <!-- Amount & Price -->
            <div class="flex items-center gap-8">
              <div class="text-center">
                <p class="text-xs text-gray-500 mb-1">Quantité</p>
                <p class="font-semibold text-gray-900 dark:text-white">{{ order.amount }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500 mb-1">Prix</p>
                <p class="font-semibold text-gray-900 dark:text-white">{{ order.price ? formatPrice(order.price) : 'Market' }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500 mb-1">Exécuté</p>
                <div class="flex items-center gap-2">
                  <div class="w-16 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <div 
                      class="h-full rounded-full transition-all duration-500"
                      :class="order.side === 'buy' ? 'bg-green-500' : 'bg-red-500'"
                      :style="{ width: (order.filled_amount / order.amount * 100) + '%' }"
                    ></div>
                  </div>
                  <span class="text-xs font-medium">{{ Math.round((order.filled_amount || 0) / order.amount * 100) }}%</span>
                </div>
              </div>
            </div>

            <!-- Status & Actions -->
            <div class="flex items-center gap-3">
              <span :class="getStatusClass(order.status)">
                {{ getStatusLabel(order.status) }}
              </span>
              <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button 
                  @click="viewInvoice(order)"
                  class="p-2 rounded-lg bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-200 dark:hover:bg-indigo-900/50 transition-colors"
                  title="Voir la facture"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                </button>
                <button 
                  @click="viewOrderDetails(order)"
                  class="p-2 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
                  title="Détails"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
                <button 
                  v-if="order.status === 'pending'"
                  @click="cancelOrder(order.id)"
                  class="p-2 rounded-lg bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
                  title="Annuler"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="filteredOrders.length > 0" class="mt-8 flex items-center justify-between">
        <div class="text-sm text-gray-500">
          Affichage de {{ (currentPage - 1) * itemsPerPage + 1 }} à {{ Math.min(currentPage * itemsPerPage, filteredOrders.length) }} sur {{ filteredOrders.length }} ordres
        </div>
        <div class="flex items-center gap-2">
          <button 
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Précédent
          </button>
          <div class="flex gap-1">
            <button 
              v-for="page in displayedPages" 
              :key="page"
              @click="currentPage = page"
              :class="[
                'w-10 h-10 rounded-lg font-medium transition-colors',
                currentPage === page 
                  ? 'bg-indigo-600 text-white' 
                  : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
              ]"
            >
              {{ page }}
            </button>
          </div>
          <button 
            @click="currentPage++"
            :disabled="currentPage === totalPages"
            class="px-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Suivant
          </button>
        </div>
      </div>
    </div>

    <!-- Invoice Modal -->
    <InvoiceModal 
      :is-open="showInvoiceModal"
      :order="selectedOrder"
      :invoice="selectedInvoice"
      @close="showInvoiceModal = false"
    />

    <!-- Order Details Modal -->
    <div v-if="showDetailsModal && selectedOrder" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="fixed inset-0 bg-black/60 backdrop-blur-sm" @click="showDetailsModal = false"></div>
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="relative w-full max-w-2xl bg-white dark:bg-slate-900 rounded-2xl shadow-2xl">
          <div class="p-6 border-b border-gray-200 dark:border-gray-800">
            <div class="flex items-center justify-between">
              <h3 class="text-xl font-bold text-gray-900 dark:text-white">Détails de l'ordre</h3>
              <button @click="showDetailsModal = false" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors">
                <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-4">
                <div>
                  <p class="text-sm text-gray-500 mb-1">ID de l'ordre</p>
                  <p class="font-mono text-gray-900 dark:text-white">{{ selectedOrder.id }}</p>
                </div>
                <div>
                  <p class="text-sm text-gray-500 mb-1">Paire</p>
                  <p class="font-semibold text-gray-900 dark:text-white">{{ selectedOrder.pair }}</p>
                </div>
                <div>
                  <p class="text-sm text-gray-500 mb-1">Type</p>
                  <p class="text-gray-900 dark:text-white capitalize">{{ selectedOrder.order_type }}</p>
                </div>
              </div>
              <div class="space-y-4">
                <div>
                  <p class="text-sm text-gray-500 mb-1">Statut</p>
                  <span :class="getStatusClass(selectedOrder.status)">{{ getStatusLabel(selectedOrder.status) }}</span>
                </div>
                <div>
                  <p class="text-sm text-gray-500 mb-1">Quantité</p>
                  <p class="font-semibold text-gray-900 dark:text-white">{{ selectedOrder.amount }}</p>
                </div>
                <div>
                  <p class="text-sm text-gray-500 mb-1">Prix</p>
                  <p class="text-gray-900 dark:text-white">{{ selectedOrder.price ? formatPrice(selectedOrder.price) : 'Market' }}</p>
                </div>
              </div>
            </div>
            <div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-800">
              <div class="flex justify-between items-center">
                <span class="text-gray-600 dark:text-gray-400">Total</span>
                <span class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">{{ calculateOrderTotal(selectedOrder) }}</span>
              </div>
            </div>
            <div class="mt-6 flex gap-3">
              <button 
                @click="viewInvoice(selectedOrder); showDetailsModal = false"
                class="flex-1 btn-premium"
              >
                Voir la facture
              </button>
              <button 
                v-if="selectedOrder.status === 'pending'"
                @click="cancelOrder(selectedOrder.id); showDetailsModal = false"
                class="flex-1 btn-danger"
              >
                Annuler l'ordre
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { exchangeAPI } from '~/composables/useApi'
import InvoiceModal from '~/components/orders/InvoiceModal.vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const authStore = useAuthStore()
const loading = ref(true)
const allOrders = ref([])
const activeFilter = ref('all')
const searchQuery = ref('')
const pairFilter = ref('')
const typeFilter = ref('')
const dateFilter = ref('')
const currentPage = ref(1)
const itemsPerPage = 10
const showInvoiceModal = ref(false)
const showDetailsModal = ref(false)
const selectedOrder = ref(null)
const selectedInvoice = ref(null)

// Stats
const stats = computed(() => ({
  total: allOrders.value.length,
  pending: allOrders.value.filter(o => o.status === 'pending').length,
  filled: allOrders.value.filter(o => o.status === 'filled').length,
  volume: allOrders.value.reduce((acc, o) => acc + (o.amount * (o.price || 0)), 0)
}))

// Available pairs from orders
const availablePairs = computed(() => {
  const pairs = new Set(allOrders.value.map(o => o.pair))
  return Array.from(pairs).sort()
})

// Filter tabs
const filters = computed(() => [
  { key: 'all', label: 'Tous', count: allOrders.value.length },
  { key: 'pending', label: 'En attente', count: allOrders.value.filter(o => o.status === 'pending').length },
  { key: 'filled', label: 'Exécutés', count: allOrders.value.filter(o => o.status === 'filled').length },
  { key: 'cancelled', label: 'Annulés', count: allOrders.value.filter(o => o.status === 'cancelled').length }
])

// Filtered orders
const filteredOrders = computed(() => {
  let filtered = allOrders.value

  // Filter by status tab
  if (activeFilter.value !== 'all') {
    filtered = filtered.filter(order => order.status === activeFilter.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(order => 
      order.pair.toLowerCase().includes(query) ||
      order.id.toLowerCase().includes(query) ||
      order.status.toLowerCase().includes(query)
    )
  }

  // Filter by pair
  if (pairFilter.value) {
    filtered = filtered.filter(order => order.pair === pairFilter.value)
  }

  // Filter by type
  if (typeFilter.value) {
    filtered = filtered.filter(order => order.order_type === typeFilter.value)
  }

  // Filter by date
  if (dateFilter.value) {
    const now = new Date()
    filtered = filtered.filter(order => {
      const orderDate = new Date(order.created_at)
      switch (dateFilter.value) {
        case 'today':
          return orderDate.toDateString() === now.toDateString()
        case 'week':
          const weekAgo = new Date(now - 7 * 24 * 60 * 60 * 1000)
          return orderDate >= weekAgo
        case 'month':
          return orderDate.getMonth() === now.getMonth() && orderDate.getFullYear() === now.getFullYear()
        case 'year':
          return orderDate.getFullYear() === now.getFullYear()
        default:
          return true
      }
    })
  }

  return filtered
})

// Paginated orders
const paginatedOrders = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  return filteredOrders.value.slice(start, start + itemsPerPage)
})

// Total pages
const totalPages = computed(() => Math.ceil(filteredOrders.value.length / itemsPerPage))

// Displayed page numbers
const displayedPages = computed(() => {
  const pages = []
  const maxDisplayed = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxDisplayed / 2))
  let end = Math.min(totalPages.value, start + maxDisplayed - 1)
  
  if (end - start < maxDisplayed - 1) {
    start = Math.max(1, end - maxDisplayed + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Formatters
const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatPrice = (price) => {
  if (!price) return '-'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'USD' }).format(price)
}

const formatVolume = (volume) => {
  if (volume >= 1000000) {
    return (volume / 1000000).toFixed(2) + 'M'
  } else if (volume >= 1000) {
    return (volume / 1000).toFixed(2) + 'K'
  }
  return volume.toFixed(2)
}

const calculateOrderTotal = (order) => {
  if (!order || !order.price) return '-'
  return formatPrice(order.price * order.amount)
}

// Status helpers
const getStatusClass = (status) => {
  const classes = {
    'pending': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
    'filled': 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
    'cancelled': 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
    'partially_filled': 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
  }
  return `px-3 py-1 text-xs font-medium rounded-full ${classes[status] || classes.pending}`
}

const getStatusLabel = (status) => {
  const labels = {
    'pending': '⏳ En attente',
    'filled': '✅ Exécuté',
    'cancelled': '❌ Annulé',
    'partially_filled': '🔄 Partiel'
  }
  return labels[status] || status
}

// Actions
const viewOrderDetails = (order) => {
  selectedOrder.value = order
  showDetailsModal.value = true
}

const viewInvoice = (order) => {
  selectedOrder.value = order
  selectedInvoice.value = {
    invoice_number: `INV-${order.id.slice(-8).toUpperCase()}`,
    created_at: order.created_at,
    payment_method: 'Carte bancaire',
    transaction_id: `TXN-${Date.now()}`,
    status: 'paid'
  }
  showInvoiceModal.value = true
}

const cancelOrder = async (orderId) => {
  if (!confirm('Êtes-vous sûr de vouloir annuler cet ordre ?')) return
  
  try {
    await exchangeAPI.cancelOrder(orderId)
    const order = allOrders.value.find(o => o.id === orderId)
    if (order) {
      order.status = 'cancelled'
    }
    alert('Ordre annulé avec succès')
  } catch (error) {
    console.error('Error cancelling order:', error)
    alert('Erreur lors de l\'annulation de l\'ordre')
  }
}

const exportOrders = () => {
  const csvContent = [
    ['ID', 'Paire', 'Type', 'Côté', 'Quantité', 'Prix', 'Statut', 'Date'].join(','),
    ...filteredOrders.value.map(o => [
      o.id,
      o.pair,
      o.order_type,
      o.side,
      o.amount,
      o.price || 'Market',
      o.status,
      new Date(o.created_at).toISOString()
    ].join(','))
  ].join('\n')
  
  const blob = new Blob([csvContent], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `orders_${new Date().toISOString().split('T')[0]}.csv`
  a.click()
  window.URL.revokeObjectURL(url)
}

// Fetch orders
const fetchOrders = async () => {
  loading.value = true
  try {
    const response = await exchangeAPI.getUserOrders()
    allOrders.value = response.data?.orders || []
  } catch (error) {
    console.error('Error fetching orders:', error)
    // Demo data fallback
    allOrders.value = [
      {
        id: 'ord-001',
        pair: 'BTC/USD',
        side: 'buy',
        order_type: 'market',
        amount: 0.5,
        price: 45000,
        filled_amount: 0.5,
        status: 'filled',
        created_at: new Date(Date.now() - 86400000).toISOString()
      },
      {
        id: 'ord-002',
        pair: 'ETH/USD',
        side: 'sell',
        order_type: 'limit',
        amount: 5,
        price: 3200,
        filled_amount: 2,
        status: 'partially_filled',
        created_at: new Date(Date.now() - 172800000).toISOString()
      },
      {
        id: 'ord-003',
        pair: 'EUR/USD',
        side: 'buy',
        order_type: 'limit',
        amount: 10000,
        price: 1.08,
        filled_amount: 0,
        status: 'pending',
        created_at: new Date().toISOString()
      }
    ]
  } finally {
    loading.value = false
  }
}

// Reset page when filters change
watch([activeFilter, searchQuery, pairFilter, typeFilter, dateFilter], () => {
  currentPage.value = 1
})

onMounted(fetchOrders)
</script>
