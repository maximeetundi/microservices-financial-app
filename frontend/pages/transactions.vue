<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-base mb-2">Historique des Transactions 📊</h1>
          <p class="text-muted">Consultez toutes vos transactions</p>
        </div>
        <div class="flex gap-3">
          <select v-model="filterType" class="select-field">
            <option value="">Tous les types</option>
            <option value="deposit">Dépôts</option>
            <option value="withdraw">Retraits</option>
            <option value="transfer">Transferts</option>
            <option value="exchange">Échanges</option>
          </select>
          <select v-model="filterPeriod" class="select-field">
            <option value="all">Toute période</option>
            <option value="today">Aujourd'hui</option>
            <option value="week">Cette semaine</option>
            <option value="month">Ce mois</option>
          </select>
        </div>
      </div>

      <!-- Transactions List -->
      <div class="glass-card p-6">
        <div v-if="loading" class="py-16 text-center">
          <div class="loading-spinner w-8 h-8 mx-auto mb-4"></div>
          <p class="text-muted">Chargement des transactions...</p>
        </div>

        <div v-else-if="filteredTransactions.length > 0" class="space-y-3">
          <div v-for="tx in filteredTransactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl hover:bg-surface-hover transition-colors border border-transparent hover:border-secondary-100 dark:hover:border-secondary-800">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center font-bold text-lg shadow-inner" 
                  :class="getTypeColor(tx.type)">
                {{ getTypeIcon(tx.type) }}
              </div>
              <div>
                <p class="font-bold text-base">{{ tx.title }}</p>
                <p class="text-sm text-muted">{{ tx.description }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-lg" :class="tx.amount >= 0 ? 'text-emerald-500' : 'text-base'">
                {{ tx.amount >= 0 ? '+' : '' }}{{ formatMoney(tx.amount, tx.currency) }}
              </p>
              <p class="text-xs text-muted">{{ formatDate(tx.date) }}</p>
            </div>
          </div>

          <!-- Load More -->
          <button v-if="hasMore" @click="loadMore" 
                  class="w-full py-3 text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded-xl transition-colors font-medium">
            Charger plus
          </button>
        </div>

        <div v-else class="text-center py-16">
          <div class="w-16 h-16 rounded-full bg-surface-hover flex items-center justify-center mx-auto mb-4">
             <span class="text-3xl grayscale opacity-50">📜</span>
          </div>
          <p class="text-muted font-medium mb-2">Aucune transaction</p>
          <p class="text-sm text-muted">Vos transactions apparaîtront ici</p>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { walletAPI } from '~/composables/useApi'

const transactions = ref([])
const loading = ref(false)
const filterType = ref('')
const filterPeriod = ref('all')
const offset = ref(0)
const limit = 50
const hasMore = ref(false)

const filteredTransactions = computed(() => {
  let result = [...transactions.value]
  
  if (filterType.value) {
    result = result.filter(tx => tx.type === filterType.value)
  }
  
  if (filterPeriod.value !== 'all') {
    const now = new Date()
    const startDate = new Date()
    
    switch (filterPeriod.value) {
      case 'today':
        startDate.setHours(0, 0, 0, 0)
        break
      case 'week':
        startDate.setDate(now.getDate() - 7)
        break
      case 'month':
        startDate.setMonth(now.getMonth() - 1)
        break
    }
    
    result = result.filter(tx => new Date(tx.date) >= startDate)
  }
  
  return result
})

const formatMoney = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('fr-FR', { 
    style: 'currency', 
    currency 
  }).format(Math.abs(amount))
}

const formatDate = (date) => {
  return new Intl.DateTimeFormat('fr-FR', { 
    day: '2-digit', 
    month: 'short',
    year: 'numeric',
    hour: '2-digit', 
    minute: '2-digit' 
  }).format(new Date(date))
}

const getTypeIcon = (type) => {
  const icons = {
    deposit: '↓',
    withdraw: '↑',
    transfer: '💸',
    exchange: '💱',
    payment: '💳'
  }
  return icons[type] || '💰'
}

const getTypeColor = (type) => {
  const colors = {
    deposit: 'bg-emerald-500/10 text-emerald-500',
    withdraw: 'bg-rose-500/10 text-rose-500',
    transfer: 'bg-purple-500/10 text-purple-500',
    exchange: 'bg-blue-500/10 text-blue-500',
    payment: 'bg-orange-500/10 text-orange-500'
  }
  return colors[type] || 'bg-gray-500/10 text-gray-500'
}

const fetchTransactions = async () => {
  loading.value = true
  try {
    // Get all wallets first
    const walletsRes = await walletAPI.getAll()
    if (!walletsRes.data?.wallets) {
      transactions.value = []
      return
    }
    
    // Fetch transactions from all wallets
    const allTransactions = []
    for (const wallet of walletsRes.data.wallets) {
      try {
        const txRes = await walletAPI.getTransactions(wallet.id, limit, offset.value)
        if (txRes.data?.transactions) {
          txRes.data.transactions.forEach(tx => {
            allTransactions.push({
              id: tx.id,
              type: tx.transaction_type || 'transfer',
              title: getTransactionTitle(tx.transaction_type),
              description: tx.description || `${wallet.currency} - ${tx.reference_id || 'N/A'}`,
              amount: tx.from_wallet_id === wallet.id ? -tx.amount : tx.amount,
              currency: tx.currency || wallet.currency,
              date: tx.created_at
            })
          })
        }
      } catch (e) {
        console.error(`Error fetching transactions for wallet ${wallet.id}:`, e)
      }
    }
    
    // Sort by date descending
    allTransactions.sort((a, b) => new Date(b.date) - new Date(a.date))
    transactions.value = allTransactions
    hasMore.value = allTransactions.length >= limit
  } catch (e) {
    console.error('Error fetching transactions:', e)
    transactions.value = []
  } finally {
    loading.value = false
  }
}

const getTransactionTitle = (type) => {
  const titles = {
    deposit: 'Dépôt',
    withdraw: 'Retrait',
    transfer: 'Transfert',
    exchange: 'Échange',
    payment: 'Paiement'
  }
  return titles[type] || 'Transaction'
}

const loadMore = async () => {
  offset.value += limit
  await fetchTransactions()
}

onMounted(() => {
  fetchTransactions()
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.select-field {
  @apply px-4 py-2 rounded-xl bg-surface border border-secondary-200 dark:border-secondary-700 text-base focus:outline-none focus:ring-2 focus:ring-primary-500;
}
</style>
