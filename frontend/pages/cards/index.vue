<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">Mes Cartes 💳</h1>
          <p class="text-gray-900/60">Gérez vos cartes virtuelles et physiques</p>
        </div>
        <NuxtLink to="/cards/new" class="btn-premium flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Nouvelle Carte
        </NuxtLink>
      </div>

      <!-- Cards Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="card in cards" :key="card.id" 
            class="credit-card cursor-pointer group hover:scale-105 transition-transform"
            :class="getCardClass(card.type)"
            @click="selectCard(card)">
          <div class="relative z-10 h-full flex flex-col">
            <div class="flex justify-between items-start mb-6">
              <div>
                <span class="badge" :class="card.status === 'active' ? 'badge-success' : 'badge-warning'">
                  {{ card.status === 'active' ? 'Active' : card.status }}
                </span>
              </div>
              <span class="text-2xl">{{ card.type === 'virtual' ? '🌐' : '💳' }}</span>
            </div>
            
            <div class="flex-1">
              <p class="text-xl font-mono text-gray-900 tracking-wider mb-2">
                •••• •••• •••• {{ card.last4 || '0000' }}
              </p>
              <p class="text-gray-900/60 text-sm uppercase">{{ card.name || 'Ma Carte' }}</p>
            </div>

            <div class="flex justify-between items-end mt-4">
              <div>
                <p class="text-xs text-gray-900/50">Solde disponible</p>
                <p class="text-lg font-bold text-gray-900">{{ formatMoney(card.balance, card.currency) }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-gray-900/50">Expire</p>
                <p class="text-gray-900 font-medium">{{ card.expiry || '12/28' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Add Card Button -->
        <NuxtLink to="/cards/new" class="credit-card flex items-center justify-center border-2 border-dashed border-white/20 hover:border-indigo-500 transition-colors" style="background: transparent;">
          <div class="text-center">
            <div class="w-16 h-16 rounded-full bg-white/10 flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-gray-900/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
            </div>
            <p class="text-gray-900/60 font-medium">Commander une carte</p>
          </div>
        </NuxtLink>
      </div>

      <!-- Selected Card Details -->
      <div v-if="selectedCard" class="glass-card p-6 mb-8">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">Détails de la carte</h3>
          <div class="flex gap-2">
            <button @click="freezeCard" class="btn-secondary-premium text-sm">
              {{ selectedCard.status === 'frozen' ? '🔓 Dégeler' : '🔒 Geler' }}
            </button>
            <button @click="showTopUp = true" class="btn-success text-sm">
              💳 Recharger
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Numéro</p>
            <p class="text-gray-900 font-mono">•••• •••• •••• {{ selectedCard.last4 }}</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Date d'expiration</p>
            <p class="text-gray-900">{{ selectedCard.expiry }}</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">CVV</p>
            <p class="text-gray-900 font-mono">•••</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Limite quotidienne</p>
            <p class="text-gray-900">{{ formatMoney(selectedCard.dailyLimit || 5000, selectedCard.currency) }}</p>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">Transactions récentes</h3>
        
        <div v-if="transactions.length > 0" class="space-y-3">
          <div v-for="tx in transactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl bg-white/5 hover:bg-white/10 transition-colors">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center" 
                  :class="tx.type === 'credit' ? 'bg-green-500/20' : 'bg-red-500/20'">
                <span class="text-lg">{{ tx.type === 'credit' ? '↓' : '↑' }}</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">{{ tx.description }}</p>
                <p class="text-sm text-gray-900/50">{{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold" :class="tx.type === 'credit' ? 'text-green-400' : 'text-gray-900'">
                {{ tx.type === 'credit' ? '+' : '-' }}{{ formatMoney(tx.amount, tx.currency) }}
              </p>
              <p class="text-xs text-gray-900/40">{{ tx.category }}</p>
            </div>
          </div>
        </div>
        
        <div v-else class="text-center py-12 text-gray-900/50">
          Aucune transaction récente
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { cardAPI } from '~/composables/useApi'

// Cards - will be loaded from API
const cards = ref([])

const selectedCard = ref(null)
const showTopUp = ref(false)

// Transactions - will be loaded from API
const transactions = ref([])

const formatMoney = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

const getCardClass = (type) => {
  return type === 'virtual' ? 'credit-card-virtual' : 'credit-card-physical'
}

const selectCard = (card) => {
  selectedCard.value = card
}

const freezeCard = async () => {
  if (!selectedCard.value) return
  // API call
  selectedCard.value.status = selectedCard.value.status === 'frozen' ? 'active' : 'frozen'
}

const fetchCards = async () => {
  try {
    const response = await cardAPI.getAll()
    if (response.data?.cards) {
      cards.value = response.data.cards
    }
  } catch (e) {
    console.log('Using mock data')
  }
}

onMounted(() => {
  fetchCards()
  if (cards.value.length > 0) {
    selectedCard.value = cards.value[0]
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>