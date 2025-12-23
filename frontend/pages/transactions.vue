<template>
  <NuxtLayout name="dashboard">
    <div class="tx-page">
      <!-- Header -->
      <div class="page-header">
        <h1>📊 Transactions</h1>
        <p>Historique de vos transactions</p>
      </div>

      <!-- Filters -->
      <div class="filters">
        <select v-model="filterType" class="filter-select">
          <option value="">Tous types</option>
          <option value="deposit">Dépôts</option>
          <option value="withdraw">Retraits</option>
          <option value="transfer">Transferts</option>
          <option value="exchange">Échanges</option>
        </select>
        <select v-model="filterPeriod" class="filter-select">
          <option value="all">Toujours</option>
          <option value="today">Aujourd'hui</option>
          <option value="week">7 jours</option>
          <option value="month">30 jours</option>
        </select>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="state-container">
        <div class="spinner"></div>
        <p>Chargement...</p>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredTransactions.length === 0" class="state-container">
        <span class="empty-icon">📜</span>
        <p>Aucune transaction</p>
      </div>

      <!-- List -->
      <div v-else class="tx-list">
        <div v-for="tx in filteredTransactions" :key="tx.id" class="tx-item">
          <div class="tx-icon" :class="getTypeColorClass(tx.type)">
            {{ getTypeIcon(tx.type) }}
          </div>
          <div class="tx-info">
            <div class="tx-title">{{ tx.title }}</div>
            <div class="tx-desc">{{ tx.description }}</div>
          </div>
          <div class="tx-amount" :class="{ positive: tx.amount >= 0 }">
            <span>{{ tx.amount >= 0 ? '+' : '' }}{{ formatMoney(tx.amount, tx.currency) }}</span>
            <small>{{ formatDate(tx.date) }}</small>
          </div>
        </div>

        <!-- Load More -->
        <button v-if="hasMore" @click="loadMore" class="load-more-btn">
          Charger plus ↓
        </button>
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
    month: 'short'
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

const getTypeColorClass = (type) => {
  const colors = {
    deposit: 'icon-green',
    withdraw: 'icon-red',
    transfer: 'icon-purple',
    exchange: 'icon-blue',
    payment: 'icon-orange'
  }
  return colors[type] || 'icon-gray'
}

const fetchTransactions = async () => {
  loading.value = true
  try {
    const walletsRes = await walletAPI.getAll()
    if (!walletsRes.data?.wallets) {
      transactions.value = []
      return
    }
    
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
              description: tx.description || `${wallet.currency}`,
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
.tx-page {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.page-header {
  margin-bottom: 1rem;
}

.page-header h1 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.page-header p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

.filters {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.filter-select {
  flex: 1;
  min-width: 120px;
  max-width: 100%;
  padding: 0.625rem 0.875rem;
  border-radius: 0.625rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
  cursor: pointer;
}

.filter-select:focus {
  border-color: #6366f1;
}

.state-container {
  text-align: center;
  padding: 3rem 1rem;
  color: #888;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.tx-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  padding: 0.5rem;
}

.tx-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.75rem;
  transition: background 0.2s;
}

.tx-item:active {
  background: rgba(255,255,255,0.05);
}

.tx-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
}

.icon-green { background: rgba(34, 197, 94, 0.15); color: #22c55e; }
.icon-red { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.icon-purple { background: rgba(168, 85, 247, 0.15); color: #a855f7; }
.icon-blue { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.icon-orange { background: rgba(249, 115, 22, 0.15); color: #f97316; }
.icon-gray { background: rgba(107, 114, 128, 0.15); color: #6b7280; }

.tx-info {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.tx-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tx-desc {
  font-size: 0.75rem;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tx-amount {
  text-align: right;
  flex-shrink: 0;
}

.tx-amount span {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: #fff;
}

.tx-amount.positive span {
  color: #22c55e;
}

.tx-amount small {
  font-size: 0.625rem;
  color: #666;
}

.load-more-btn {
  width: 100%;
  padding: 0.75rem;
  border: none;
  background: transparent;
  color: #6366f1;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  border-radius: 0.5rem;
  transition: background 0.2s;
}

.load-more-btn:hover {
  background: rgba(99, 102, 241, 0.1);
}

/* Desktop */
@media (min-width: 640px) {
  .page-header h1 {
    font-size: 1.5rem;
  }
  
  .filters {
    flex-wrap: nowrap;
  }
  
  .filter-select {
    flex: 0 0 auto;
    width: auto;
  }
  
  .tx-list {
    padding: 0.75rem;
  }
  
  .tx-item {
    gap: 1rem;
    padding: 1rem;
  }
  
  .tx-icon {
    width: 3rem;
    height: 3rem;
    font-size: 1.25rem;
  }
  
  .tx-title {
    font-size: 1rem;
  }
  
  .tx-desc {
    font-size: 0.875rem;
  }
  
  .tx-amount span {
    font-size: 1rem;
  }
}
</style>
