<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            💳 Mes Cartes
          </h1>
          <p class="text-gray-600">
            Gérez vos cartes prépayées crypto et fiat
          </p>
        </div>
        <div class="flex space-x-3 mt-4 md:mt-0">
          <button
            @click="showCreateCardModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors duration-200"
          >
            + Nouvelle Carte
          </button>
          <button
            @click="showGiftCardModal = true"
            class="bg-purple-600 hover:bg-purple-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors duration-200"
          >
            🎁 Carte Cadeau
          </button>
        </div>
      </div>

      <!-- Quick Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-blue-100 text-blue-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Cartes</p>
              <p class="text-2xl font-semibold text-gray-900">{{ totalCards }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-green-100 text-green-600">
              <span class="text-xl">💰</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Solde Total</p>
              <p class="text-2xl font-semibold text-gray-900">{{ totalBalance }}€</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-purple-100 text-purple-600">
              <span class="text-xl">📱</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Cartes Virtuelles</p>
              <p class="text-2xl font-semibold text-gray-900">{{ virtualCards }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-yellow-100 text-yellow-600">
              <span class="text-xl">🏪</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Cartes Physiques</p>
              <p class="text-2xl font-semibold text-gray-900">{{ physicalCards }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Cards List -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div
          v-for="card in userCards"
          :key="card.id"
          class="relative group"
        >
          <!-- Card Component -->
          <div 
            class="relative h-48 rounded-xl shadow-lg overflow-hidden cursor-pointer transform transition-all duration-300 group-hover:scale-105"
            :class="getCardStyle(card)"
            @click="selectCard(card)"
          >
            <!-- Card Background -->
            <div class="absolute inset-0 bg-gradient-to-br opacity-90" :class="getCardGradient(card.currency)"></div>
            
            <!-- Card Content -->
            <div class="relative h-full p-6 flex flex-col justify-between text-white">
              <div class="flex justify-between items-start">
                <div>
                  <p class="text-sm opacity-80">{{ card.card_type.toUpperCase() }} CARD</p>
                  <p class="font-bold text-lg">{{ card.currency }}</p>
                </div>
                <div class="text-right">
                  <div class="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">
                    <span v-if="card.is_virtual">📱</span>
                    <span v-else>💳</span>
                  </div>
                </div>
              </div>

              <div>
                <p class="text-lg font-mono tracking-wider">{{ card.card_number }}</p>
                <div class="flex justify-between items-end mt-3">
                  <div>
                    <p class="text-xs opacity-80">SOLDE</p>
                    <p class="font-bold text-lg">{{ formatMoney(card.balance, card.currency) }}</p>
                  </div>
                  <div class="text-right">
                    <p class="text-xs opacity-80">STATUS</p>
                    <span 
                      class="inline-block px-2 py-1 rounded-full text-xs font-semibold"
                      :class="getStatusClass(card.status)"
                    >
                      {{ card.status.toUpperCase() }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Card Actions Overlay -->
            <div class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
              <div class="flex space-x-3">
                <button
                  @click.stop="loadCard(card)"
                  class="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg font-semibold transition-all duration-200"
                >
                  💰 Charger
                </button>
                <button
                  @click.stop="viewCard(card)"
                  class="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg font-semibold transition-all duration-200"
                >
                  👁️ Détails
                </button>
                <button
                  @click.stop="manageCard(card)"
                  class="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg font-semibold transition-all duration-200"
                >
                  ⚙️ Gérer
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Add New Card -->
        <div 
          @click="showCreateCardModal = true"
          class="h-48 rounded-xl border-2 border-dashed border-gray-300 hover:border-blue-400 cursor-pointer transition-colors duration-200 flex items-center justify-center group"
        >
          <div class="text-center">
            <svg class="w-12 h-12 text-gray-400 group-hover:text-blue-500 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <p class="text-gray-500 group-hover:text-blue-500 font-semibold">Ajouter une carte</p>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-lg shadow-lg p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-gray-900">📊 Transactions Récentes</h3>
          <button class="text-blue-600 hover:text-blue-800 font-semibold">
            Voir tout
          </button>
        </div>
        
        <div v-if="recentTransactions.length === 0" class="text-center py-8 text-gray-500">
          <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <p>Aucune transaction récente</p>
        </div>

        <div v-else class="space-y-3">
          <div 
            v-for="transaction in recentTransactions"
            :key="transaction.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
          >
            <div class="flex items-center space-x-4">
              <div class="w-10 h-10 rounded-full flex items-center justify-center" :class="getTransactionIcon(transaction.transaction_type)">
                <span class="text-white font-semibold">
                  {{ getTransactionEmoji(transaction.transaction_type) }}
                </span>
              </div>
              <div>
                <p class="font-semibold text-gray-900">
                  {{ getTransactionDescription(transaction) }}
                </p>
                <p class="text-sm text-gray-500">
                  {{ transaction.merchant_name || 'CryptoBank' }} • {{ formatTime(transaction.created_at) }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold" :class="transaction.transaction_type === 'load' ? 'text-green-600' : 'text-gray-900'">
                {{ transaction.transaction_type === 'load' ? '+' : '-' }}{{ formatMoney(transaction.amount, transaction.currency) }}
              </p>
              <span 
                class="inline-block px-2 py-1 rounded-full text-xs font-semibold"
                :class="getStatusClass(transaction.status)"
              >
                {{ transaction.status }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Create Card Modal -->
      <div v-if="showCreateCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-8 max-w-md w-full mx-4">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold">Créer une Nouvelle Carte</h3>
            <button @click="showCreateCardModal = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <button
                @click="createCardType = 'virtual'"
                class="p-4 border-2 rounded-lg text-center transition-all"
                :class="createCardType === 'virtual' ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300'"
              >
                <div class="text-2xl mb-2">📱</div>
                <div class="font-semibold">Virtuelle</div>
                <div class="text-sm text-gray-500">Instantanée</div>
              </button>
              
              <button
                @click="createCardType = 'prepaid'"
                class="p-4 border-2 rounded-lg text-center transition-all"
                :class="createCardType === 'prepaid' ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300'"
              >
                <div class="text-2xl mb-2">💳</div>
                <div class="font-semibold">Physique</div>
                <div class="text-sm text-gray-500">Livraison</div>
              </button>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Devise</label>
              <select v-model="newCard.currency" class="w-full p-3 border border-gray-300 rounded-lg">
                <option value="USD">🇺🇸 USD - Dollar US</option>
                <option value="EUR">🇪🇺 EUR - Euro</option>
                <option value="GBP">🇬🇧 GBP - Livre Sterling</option>
                <option value="BTC">₿ BTC - Bitcoin</option>
                <option value="ETH">Ξ ETH - Ethereum</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nom sur la carte</label>
              <input 
                v-model="newCard.cardholder_name"
                type="text"
                placeholder="John Doe"
                class="w-full p-3 border border-gray-300 rounded-lg"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Montant initial</label>
              <input 
                v-model="newCard.initial_amount"
                type="number"
                placeholder="100"
                min="10"
                step="0.01"
                class="w-full p-3 border border-gray-300 rounded-lg"
              />
            </div>

            <div class="flex space-x-3 mt-6">
              <button
                @click="createCard"
                :disabled="createCardLoading"
                class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-3 rounded-lg transition-colors"
              >
                <span v-if="createCardLoading">Création...</span>
                <span v-else>Créer la Carte</span>
              </button>
              <button
                @click="showCreateCardModal = false"
                class="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-700 font-semibold py-3 rounded-lg transition-colors"
              >
                Annuler
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Gift Card Modal -->
      <div v-if="showGiftCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-8 max-w-md w-full mx-4">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold">🎁 Créer une Carte Cadeau</h3>
            <button @click="showGiftCardModal = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Montant</label>
              <input 
                v-model="giftCard.amount"
                type="number"
                placeholder="50"
                min="10"
                step="0.01"
                class="w-full p-3 border border-gray-300 rounded-lg"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Devise</label>
              <select v-model="giftCard.currency" class="w-full p-3 border border-gray-300 rounded-lg">
                <option value="USD">🇺🇸 USD</option>
                <option value="EUR">🇪🇺 EUR</option>
                <option value="BTC">₿ BTC</option>
                <option value="ETH">Ξ ETH</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Email du destinataire</label>
              <input 
                v-model="giftCard.recipient_email"
                type="email"
                placeholder="destinataire@email.com"
                class="w-full p-3 border border-gray-300 rounded-lg"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Message</label>
              <textarea 
                v-model="giftCard.message"
                placeholder="Joyeux anniversaire !"
                rows="3"
                class="w-full p-3 border border-gray-300 rounded-lg"
              ></textarea>
            </div>

            <div class="flex space-x-3 mt-6">
              <button
                @click="createGiftCard"
                :disabled="giftCardLoading"
                class="flex-1 bg-purple-600 hover:bg-purple-700 disabled:bg-gray-400 text-white font-semibold py-3 rounded-lg transition-colors"
              >
                <span v-if="giftCardLoading">Création...</span>
                <span v-else>🎁 Créer & Envoyer</span>
              </button>
              <button
                @click="showGiftCardModal = false"
                class="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-700 font-semibold py-3 rounded-lg transition-colors"
              >
                Annuler
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Reactive data
const userCards = ref([])
const recentTransactions = ref([])
const showCreateCardModal = ref(false)
const showGiftCardModal = ref(false)
const createCardLoading = ref(false)
const giftCardLoading = ref(false)
const createCardType = ref('virtual')

const newCard = ref({
  currency: 'USD',
  cardholder_name: '',
  initial_amount: 100,
})

const giftCard = ref({
  amount: 50,
  currency: 'USD',
  recipient_email: '',
  message: ''
})

// Computed properties
const totalCards = computed(() => userCards.value.length)
const totalBalance = computed(() => {
  return userCards.value.reduce((sum, card) => sum + card.balance, 0).toFixed(2)
})
const virtualCards = computed(() => {
  return userCards.value.filter(card => card.is_virtual).length
})
const physicalCards = computed(() => {
  return userCards.value.filter(card => !card.is_virtual).length
})

// Methods
const loadUserCards = async () => {
  try {
    // Mock data for demonstration
    userCards.value = [
      {
        id: '1',
        card_number: '•••• •••• •••• 1234',
        card_type: 'virtual',
        currency: 'USD',
        balance: 1250.50,
        status: 'active',
        is_virtual: true,
        created_at: new Date().toISOString()
      },
      {
        id: '2', 
        card_number: '•••• •••• •••• 5678',
        card_type: 'prepaid',
        currency: 'EUR',
        balance: 850.25,
        status: 'active',
        is_virtual: false,
        created_at: new Date().toISOString()
      },
      {
        id: '3',
        card_number: '•••• •••• •••• 9012',
        card_type: 'virtual',
        currency: 'BTC',
        balance: 0.05,
        status: 'active',
        is_virtual: true,
        created_at: new Date().toISOString()
      }
    ]

    recentTransactions.value = [
      {
        id: '1',
        card_id: '1',
        transaction_type: 'purchase',
        amount: 25.99,
        currency: 'USD',
        merchant_name: 'Amazon',
        status: 'completed',
        created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString()
      },
      {
        id: '2',
        card_id: '1', 
        transaction_type: 'load',
        amount: 100.00,
        currency: 'USD',
        status: 'completed',
        created_at: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString()
      }
    ]
  } catch (error) {
    console.error('Error loading cards:', error)
  }
}

const createCard = async () => {
  createCardLoading.value = true
  try {
    const cardData = {
      card_type: createCardType.value,
      currency: newCard.value.currency,
      cardholder_name: newCard.value.cardholder_name,
      initial_amount: newCard.value.initial_amount,
      card_category: 'personal'
    }

    // Mock card creation
    const newCardData = {
      id: Date.now().toString(),
      card_number: '•••• •••• •••• ' + Math.floor(Math.random() * 9000 + 1000),
      card_type: createCardType.value,
      currency: newCard.value.currency,
      balance: newCard.value.initial_amount,
      status: 'active',
      is_virtual: createCardType.value === 'virtual',
      cardholder_name: newCard.value.cardholder_name,
      created_at: new Date().toISOString()
    }

    userCards.value.unshift(newCardData)
    showCreateCardModal.value = false
    
    // Reset form
    newCard.value = {
      currency: 'USD',
      cardholder_name: '',
      initial_amount: 100
    }

    showNotification('Carte créée avec succès !', 'success')
  } catch (error) {
    console.error('Error creating card:', error)
    showNotification('Erreur lors de la création de la carte', 'error')
  } finally {
    createCardLoading.value = false
  }
}

const createGiftCard = async () => {
  giftCardLoading.value = true
  try {
    // Mock gift card creation
    const giftCardData = {
      amount: giftCard.value.amount,
      currency: giftCard.value.currency,
      recipient_email: giftCard.value.recipient_email,
      message: giftCard.value.message,
      design: 'birthday',
      source_wallet_id: 'mock-wallet-id'
    }

    showGiftCardModal.value = false
    
    // Reset form
    giftCard.value = {
      amount: 50,
      currency: 'USD', 
      recipient_email: '',
      message: ''
    }

    showNotification('Carte cadeau créée et envoyée !', 'success')
  } catch (error) {
    console.error('Error creating gift card:', error)
    showNotification('Erreur lors de la création de la carte cadeau', 'error')
  } finally {
    giftCardLoading.value = false
  }
}

const loadCard = (card) => {
  console.log('Loading card:', card.id)
  // TODO: Implement card loading modal
}

const viewCard = (card) => {
  console.log('Viewing card:', card.id) 
  // TODO: Navigate to card details page
}

const manageCard = (card) => {
  console.log('Managing card:', card.id)
  // TODO: Implement card management modal
}

const selectCard = (card) => {
  console.log('Selected card:', card.id)
  // TODO: Navigate to card details page
}

// Utility methods
const formatMoney = (amount, currency) => {
  if (currency === 'BTC') {
    return `${amount.toFixed(8)} ₿`
  } else if (currency === 'ETH') {
    return `${amount.toFixed(6)} Ξ`
  } else {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(amount)
  }
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit'
  })
}

const getCardStyle = (card) => {
  return 'transform-gpu'
}

const getCardGradient = (currency) => {
  const gradients = {
    USD: 'from-blue-500 to-blue-700',
    EUR: 'from-purple-500 to-purple-700',
    GBP: 'from-green-500 to-green-700',
    BTC: 'from-orange-500 to-yellow-600',
    ETH: 'from-indigo-500 to-purple-600'
  }
  return gradients[currency] || 'from-gray-500 to-gray-700'
}

const getStatusClass = (status) => {
  const classes = {
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-gray-100 text-gray-800',
    blocked: 'bg-red-100 text-red-800',
    pending: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getTransactionIcon = (type) => {
  const icons = {
    purchase: 'bg-blue-500',
    load: 'bg-green-500', 
    withdrawal: 'bg-red-500',
    refund: 'bg-yellow-500'
  }
  return icons[type] || 'bg-gray-500'
}

const getTransactionEmoji = (type) => {
  const emojis = {
    purchase: '🛒',
    load: '💰',
    withdrawal: '💸',
    refund: '↩️'
  }
  return emojis[type] || '💳'
}

const getTransactionDescription = (transaction) => {
  const descriptions = {
    purchase: `Achat chez ${transaction.merchant_name || 'Commerçant'}`,
    load: 'Rechargement de carte',
    withdrawal: 'Retrait ATM',
    refund: 'Remboursement'
  }
  return descriptions[transaction.transaction_type] || 'Transaction'
}

const showNotification = (message, type) => {
  console.log(`${type.toUpperCase()}: ${message}`)
}

// Lifecycle
onMounted(() => {
  loadUserCards()
})

// SEO
definePageMeta({
  title: 'Mes Cartes - CryptoBank',
  description: 'Gérez vos cartes prépayées crypto et fiat, cartes virtuelles et physiques'
})
</script>

<style scoped>
/* Additional custom styles if needed */
</style>