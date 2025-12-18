<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-base mb-2">Mes Cartes 💳</h1>
          <p class="text-muted">Gérez vos cartes virtuelles et physiques</p>
        </div>
        <NuxtLink to="/cards/new" class="btn-primary flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Nouvelle Carte
        </NuxtLink>
      </div>

      <!-- Cards Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="card in cards" :key="card.id" 
            class="credit-card cursor-pointer group hover:scale-[1.02] transition-transform relative overflow-hidden shadow-xl"
            :class="getCardClass(card.type)"
            @click="selectCard(card)">
            
           <!-- Shine Effect -->
           <div class="absolute inset-0 bg-gradient-to-tr from-white/0 via-white/10 to-white/0 opacity-0 group-hover:opacity-100 transition-opacity duration-700 pointer-events-none transform translate-x-[-100%] group-hover:translate-x-[100%]"></div>

          <div class="relative z-10 h-full flex flex-col text-white">
            <div class="flex justify-between items-start mb-6">
              <div>
                <span class="px-2 py-1 rounded-full text-xs font-semibold backdrop-blur-md bg-white/20 border border-white/10">
                  {{ card.status === 'active' ? 'Active' : card.status }}
                </span>
              </div>
              <span class="text-2xl opacity-80">{{ card.type === 'virtual' ? '🌐' : '💳' }}</span>
            </div>
            
            <div class="flex-1 flex flex-col justify-center">
              <p class="text-xl font-mono tracking-[0.15em] mb-2 drop-shadow-md">
                •••• •••• •••• {{ card.last4 || '0000' }}
              </p>
              <p class="text-white/60 text-xs uppercase tracking-widest">{{ card.name || 'Ma Carte' }}</p>
            </div>

            <div class="flex justify-between items-end mt-4">
              <div>
                <p class="text-[10px] uppercase tracking-wider text-white/50 mb-1">Solde disponible</p>
                <p class="text-lg font-bold tracking-tight">{{ formatMoney(card.balance, card.currency) }}</p>
              </div>
              <div class="text-right">
                <p class="text-[10px] uppercase tracking-wider text-white/50 mb-1">Expire</p>
                <p class="font-mono">{{ card.expiry || '12/28' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Add Card Button -->
        <NuxtLink to="/cards/new" class="group flex flex-col items-center justify-center h-full min-h-[220px] rounded-3xl border-2 border-dashed border-secondary-200 dark:border-secondary-700 hover:border-primary-500 dark:hover:border-primary-500 transition-all bg-surface/50 hover:bg-surface-hover/50">
          <div class="text-center group-hover:scale-105 transition-transform duration-300">
            <div class="w-16 h-16 rounded-full bg-surface shadow-inner flex items-center justify-center mx-auto mb-4 group-hover:shadow-primary-500/20 group-hover:shadow-lg transition-all">
              <svg class="w-8 h-8 text-muted group-hover:text-primary-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
            </div>
            <p class="text-muted font-medium group-hover:text-primary-500 transition-colors">Commander une carte</p>
          </div>
        </NuxtLink>
      </div>

      <!-- Selected Card Details -->
      <div v-if="selectedCard" class="glass-card p-8 mb-8">
        <div class="flex items-center justify-between mb-8">
          <div>
            <h3 class="text-xl font-bold text-base">Détails de la carte</h3>
            <p class="text-muted text-sm">Gérez les paramètres de votre carte</p>
          </div>
          <div class="flex gap-3">
            <button @click="freezeCard" class="px-4 py-2 rounded-xl text-sm font-medium transition-colors border"
                :class="selectedCard.status === 'frozen' 
                  ? 'bg-amber-500/10 text-amber-500 border-amber-500/20 hover:bg-amber-500/20' 
                  : 'bg-secondary-100 dark:bg-secondary-800 text-muted border-secondary-200 dark:border-secondary-700 hover:text-base hover:border-secondary-300'">
              {{ selectedCard.status === 'frozen' ? '🔓 Dégeler la carte' : '🔒 Geler la carte' }}
            </button>
            <button @click="showTopUp = true" class="px-5 py-2 rounded-xl bg-success text-white font-medium text-sm shadow-lg shadow-success/20 hover:bg-success-600 transition-all hover:-translate-y-0.5">
               Recharger
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 p-6 rounded-xl bg-surface-hover/50 border border-secondary-100 dark:border-secondary-800/50">
          <div>
            <p class="text-muted text-xs uppercase tracking-wider mb-2 font-semibold">Numéro</p>
            <p class="text-base font-mono text-lg">•••• •••• •••• {{ selectedCard.last4 }}</p>
          </div>
          <div>
            <p class="text-muted text-xs uppercase tracking-wider mb-2 font-semibold">Expiration</p>
            <p class="text-base text-lg">{{ selectedCard.expiry }}</p>
          </div>
          <div>
            <p class="text-muted text-xs uppercase tracking-wider mb-2 font-semibold">CVV</p>
            <p class="text-base font-mono text-lg">•••</p>
          </div>
          <div>
            <p class="text-muted text-xs uppercase tracking-wider mb-2 font-semibold">Limite quotidienne</p>
            <p class="text-base font-medium text-lg">{{ formatMoney(selectedCard.dailyLimit || 5000, selectedCard.currency) }}</p>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-bold text-base mb-6">Transactions récentes</h3>
        
        <div v-if="loading" class="py-10 text-center text-muted">A few seconds please...</div>
        <div v-else-if="transactions.length > 0" class="space-y-3">
          <div v-for="tx in transactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl hover:bg-surface-hover transition-colors border border-transparent hover:border-secondary-100 dark:hover:border-secondary-800">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center font-bold text-lg shadow-inner" 
                  :class="tx.type === 'credit' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-rose-500/10 text-rose-500'">
                {{ tx.type === 'credit' ? '↓' : '↑' }}
              </div>
              <div>
                <p class="font-bold text-base">{{ tx.description }}</p>
                <p class="text-sm text-muted">{{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-lg" :class="tx.type === 'credit' ? 'text-emerald-500' : 'text-base'">
                {{ tx.type === 'credit' ? '+' : '-' }}{{ formatMoney(tx.amount, tx.currency) }}
              </p>
              <p class="text-xs text-muted font-medium uppercase tracking-wide">{{ tx.category }}</p>
            </div>
          </div>
        </div>
        
        <div v-else class="text-center py-16">
          <div class="w-16 h-16 rounded-full bg-surface-hover flex items-center justify-center mx-auto mb-4">
             <span class="text-3xl grayscale opacity-50">🧾</span>
          </div>
          <p class="text-muted font-medium">Aucune transaction récente</p>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { cardAPI } from '~/composables/useApi'

// Cards - will be loaded from API
const cards = ref([])

const selectedCard = ref(null)
const showTopUp = ref(false)

// Transactions - will be loaded from API
const transactions = ref([])
const loading = ref(false)

const formatMoney = (amount, currency = 'USD') => {
  const val = Number(amount)
  if (isNaN(val)) return '0.00'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(val)
}

const formatDate = (date) => {
  if (!date) return ''
  return new Intl.DateTimeFormat('fr-FR', { 
    day: '2-digit', 
    month: 'short', 
    hour: '2-digit', 
    minute: '2-digit' 
  }).format(new Date(date))
}

const getCardClass = (type) => {
  return type === 'virtual' ? 'credit-card-virtual' : 'credit-card-physical'
}

const selectCard = (card) => {
  selectedCard.value = card
  fetchTransactions(card.id)
}

const freezeCard = async () => {
  if (!selectedCard.value) return
  
  try {
    if (selectedCard.value.status === 'frozen') {
      await cardAPI.unfreeze(selectedCard.value.id)
      selectedCard.value.status = 'active'
    } else {
      await cardAPI.freeze(selectedCard.value.id)
      selectedCard.value.status = 'frozen'
    }
  } catch (e) {
    console.error('Failed to toggle card status', e)
    alert("Impossible de changer le statut de la carte : " + (e.response?.data?.error || e.message))
  }
}

const fetchCards = async () => {
  loading.value = true
  try {
    const response = await cardAPI.getAll()
    if (response.data?.cards) {
      cards.value = response.data.cards.map(c => ({
        id: c.id,
        type: c.card_type || 'virtual',
        name: c.cardholder_name || (c.card_category === 'personal' ? 'Personnelle' : 'Business'),
        last4: c.card_number ? c.card_number.slice(-4) : '****',
        expiry: `${String(c.expiry_month).padStart(2, '0')}/${String(c.expiry_year).slice(-2)}`,
        balance: c.balance,
        currency: c.currency,
        status: c.status, // active, frozen
        dailyLimit: c.daily_limit,
        isVirtual: c.is_virtual
      }))
    } else {
       cards.value = []
    }
  } catch (e) {
    console.error('Error fetching cards:', e)
    cards.value = [] // No mock data fallback, real implementation
  } finally {
    loading.value = false
  }
}

const fetchTransactions = async (cardId) => {
  if (!cardId) return
  try {
    const response = await cardAPI.getTransactions(cardId)
    // Map backend transactions to frontend format if needed
    // Backend: transaction_type, amount, currency, description, created_at...
    if (response.data?.transactions) {
      transactions.value = response.data.transactions.map(tx => ({
        id: tx.id,
        type: tx.amount < 0 || tx.transaction_type === 'purchase' ? 'debit' : 'credit',
        amount: Math.abs(tx.amount),
        currency: tx.currency,
        description: tx.merchant_name || tx.description || 'Transaction',
        category: tx.merchant_category || 'Général',
        date: tx.created_at
      }))
    } else {
      transactions.value = []
    }
  } catch (e) {
    console.error('Error fetching transactions:', e)
    transactions.value = []
  }
}

onMounted(async () => {
  await fetchCards()
  if (cards.value.length > 0) {
    selectCard(cards.value[0])
  }
})

// Watch for selection changes to refresh transactions if needed
watch(selectedCard, (newCard) => {
  if (newCard) {
    fetchTransactions(newCard.id)
  }
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
</style>