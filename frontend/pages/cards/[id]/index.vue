<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <NuxtLink to="/cards" class="text-blue-600 hover:text-blue-800 mr-4">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </NuxtLink>
            <h1 class="text-2xl font-bold text-gray-900">Card Management</h1>
          </div>
          <div class="flex items-center space-x-4">
            <button @click="freezeCard" :disabled="card.status === 'frozen'"
                    class="px-4 py-2 rounded-lg text-sm font-medium"
                    :class="card.status === 'frozen' 
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed' 
                      : 'bg-red-100 text-red-700 hover:bg-red-200'">
              {{ card.status === 'frozen' ? 'Frozen' : 'Freeze Card' }}
            </button>
            <button @click="showTopUp = true" 
                    class="bg-blue-600 text-gray-900 px-4 py-2 rounded-lg hover:bg-blue-700">
              Top Up
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8" v-if="card">
      
      <!-- Card Display -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        
        <!-- Card Visual -->
        <div class="bg-white dark:bg-slate-900 rounded-xl shadow-lg p-6 border border-gray-100 dark:border-gray-800">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">Your Card</h3>
          
          <div class="relative rounded-xl p-6 text-white shadow-lg mb-6"
               :class="getCardGradient(card.type)">
            
            <!-- Card Type Badge -->
            <div class="absolute top-4 right-4">
              <span class="px-2 py-1 bg-white/20 backdrop-blur-md rounded text-xs uppercase font-semibold">
                {{ card.type }}
              </span>
            </div>

            <!-- Card Number -->
            <div class="mb-6">
              <div class="text-sm opacity-75 mb-1">Card Number</div>
              <div class="text-lg font-mono tracking-wider">
                {{ showFullNumber ? card.card_number : formatCardNumber(card.card_number) }}
              </div>
              <button @click="showFullNumber = !showFullNumber" 
                      class="text-sm mt-2 bg-white/20 hover:bg-white/30 backdrop-blur-md px-3 py-1 rounded transition-colors">
                {{ showFullNumber ? 'Hide' : 'Show' }} Full Number
              </button>
            </div>

            <!-- Card Details -->
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div>
                <div class="text-xs opacity-75 mb-1">Valid Thru</div>
                <div class="text-sm font-mono">{{ card.expiry_date }}</div>
              </div>
              <div>
                <div class="text-xs opacity-75 mb-1">CVV</div>
                <div class="text-sm font-mono">
                  {{ showCVV ? card.cvv : '***' }}
                  <button @click="showCVV = !showCVV" class="ml-2 text-xs opacity-75 hover:opacity-100 underline">
                    {{ showCVV ? 'Hide' : 'Show' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- Balance -->
            <div>
              <div class="text-sm opacity-75 mb-1">Available Balance</div>
              <div class="text-2xl font-bold">
                ${{ card.balance?.toLocaleString() || '0.00' }}
              </div>
            </div>
          </div>

          <!-- Card Actions -->
          <div class="grid grid-cols-2 gap-4">
            <button @click="showChangePin = true" 
                    class="bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 py-3 px-4 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-700 font-medium transition-colors">
              Change PIN
            </button>
            <button @click="showSettings = true" 
                    class="bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 font-medium transition-colors shadow-lg shadow-blue-500/20">
              Card Settings
            </button>
          </div>
        </div>

        <!-- Card Information -->
        <div class="space-y-6">
          
          <!-- Quick Stats -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Stats</h3>
            <div class="grid grid-cols-2 gap-4">
              <div class="text-center p-4 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-gray-900">${{ (monthlySpending || 0).toLocaleString() }}</div>
                <div class="text-sm text-gray-500">This Month</div>
              </div>
              <div class="text-center p-4 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-gray-900">{{ transactionCount }}</div>
                <div class="text-sm text-gray-500">Transactions</div>
              </div>
              <div class="text-center p-4 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-green-600">${{ (cashbackEarned || 0).toLocaleString() }}</div>
                <div class="text-sm text-gray-500">Cashback Earned</div>
              </div>
              <div class="text-center p-4 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-gray-900">{{ card.status }}</div>
                <div class="text-sm text-gray-500">Status</div>
              </div>
            </div>
          </div>

          <!-- Limits & Controls -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Spending Limits</h3>
            <div class="space-y-4">
              
              <!-- Daily Limit -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-700">Daily Limit</span>
                  <button @click="editLimit('daily')" class="text-blue-600 text-sm hover:text-blue-800">
                    Edit
                  </button>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div class="bg-blue-600 h-2 rounded-full" 
                       :style="{ width: (dailySpent / card.daily_limit * 100) + '%' }"></div>
                </div>
                <div class="flex justify-between text-sm text-gray-500 mt-1">
                  <span>${{ (dailySpent || 0).toLocaleString() }} spent</span>
                  <span>${{ (card.daily_limit || 0).toLocaleString() }} limit</span>
                </div>
              </div>

              <!-- Monthly Limit -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-medium text-gray-700">Monthly Limit</span>
                  <button @click="editLimit('monthly')" class="text-blue-600 text-sm hover:text-blue-800">
                    Edit
                  </button>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div class="bg-green-600 h-2 rounded-full" 
                       :style="{ width: (monthlySpending / card.monthly_limit * 100) + '%' }"></div>
                </div>
                <div class="flex justify-between text-sm text-gray-500 mt-1">
                  <span>${{ (monthlySpending || 0).toLocaleString() }} spent</span>
                  <span>${{ (card.monthly_limit || 0).toLocaleString() }} limit</span>
                </div>
              </div>

              <!-- Transaction Toggles -->
              <div class="pt-4 border-t border-gray-200">
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-gray-700">Online Purchases</span>
                    <button @click="toggleSetting('online_purchases')" 
                            :class="card.settings?.online_purchases 
                              ? 'bg-blue-600' 
                              : 'bg-gray-300'"
                            class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors">
                      <span :class="card.settings?.online_purchases 
                              ? 'translate-x-6' 
                              : 'translate-x-1'"
                            class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"></span>
                    </button>
                  </div>
                  
                  <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-gray-700">ATM Withdrawals</span>
                    <button @click="toggleSetting('atm_withdrawals')" 
                            :class="card.settings?.atm_withdrawals 
                              ? 'bg-blue-600' 
                              : 'bg-gray-300'"
                            class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors">
                      <span :class="card.settings?.atm_withdrawals 
                              ? 'translate-x-6' 
                              : 'translate-x-1'"
                            class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"></span>
                    </button>
                  </div>
                  
                  <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-gray-700">International Transactions</span>
                    <button @click="toggleSetting('international')" 
                            :class="card.settings?.international 
                              ? 'bg-blue-600' 
                              : 'bg-gray-300'"
                            class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors">
                      <span :class="card.settings?.international 
                              ? 'translate-x-6' 
                              : 'translate-x-1'"
                            class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform"></span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-xl shadow-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Recent Transactions</h3>
            <div class="flex space-x-2">
              <select v-model="transactionFilter" class="text-sm border-gray-300 rounded-md">
                <option value="">All Categories</option>
                <option value="shopping">Shopping</option>
                <option value="food">Food & Drink</option>
                <option value="transport">Transport</option>
                <option value="entertainment">Entertainment</option>
              </select>
              <button @click="exportTransactions" 
                      class="text-blue-600 text-sm hover:text-blue-800 font-medium">
                Export
              </button>
            </div>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Merchant</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="transaction in filteredTransactions" :key="transaction.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(transaction.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm font-medium text-gray-900">{{ transaction.merchant_name }}</div>
                  <div class="text-sm text-gray-500">{{ transaction.merchant_location }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-800">
                    {{ transaction.category }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium"
                    :class="transaction.type === 'debit' ? 'text-red-600' : 'text-green-600'">
                  {{ transaction.type === 'debit' ? '-' : '+' }}${{ (transaction.amount || 0).toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full"
                        :class="getTransactionStatusColor(transaction.status)">
                    {{ transaction.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Top Up Modal -->
    <div v-if="showTopUp" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold">Top Up Card</h3>
          <button @click="showTopUp = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitTopUp" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Source Wallet</label>
            <select v-model="topUp.source_wallet" required
                    class="w-full rounded-md border-gray-300 shadow-sm">
              <option value="">Select wallet</option>
              <option v-for="wallet in availableWallets" :key="wallet.id" :value="wallet.id">
                {{ wallet.currency }} - ${{ (wallet.balance || 0).toLocaleString() }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Amount</label>
            <input v-model.number="topUp.amount" type="number" step="0.01" required
                   class="w-full rounded-md border-gray-300 shadow-sm" 
                   placeholder="0.00">
          </div>

          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex justify-between text-sm">
              <span>Top up amount:</span>
              <span class="font-medium">${{ topUp.amount?.toLocaleString() || '0.00' }}</span>
            </div>
            <div class="flex justify-between text-sm mt-1">
              <span>Fee:</span>
              <span class="font-medium text-green-600">Free</span>
            </div>
            <div class="border-t border-gray-200 mt-2 pt-2 flex justify-between">
              <span class="font-medium">Total:</span>
              <span class="font-bold">${{ topUp.amount?.toLocaleString() || '0.00' }}</span>
            </div>
          </div>

          <div class="flex space-x-3">
            <button type="button" @click="showTopUp = false" 
                    class="flex-1 bg-gray-300 text-gray-700 py-3 px-4 rounded-lg hover:bg-gray-400">
              Cancel
            </button>
            <button type="submit" 
                    class="flex-1 bg-blue-600 text-gray-900 py-3 px-4 rounded-lg hover:bg-blue-700">
              Top Up Card
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

// Reactive data
const route = useRoute()
const cardId = route.params.id

const card = ref({})
const transactions = ref([])
const availableWallets = ref([])
const showFullNumber = ref(false)
const showCVV = ref(false)
const showTopUp = ref(false)
const showChangePin = ref(false)
const showSettings = ref(false)
const transactionFilter = ref('')

const topUp = ref({
  source_wallet: '',
  amount: 0
})

// Computed
const monthlySpending = computed(() => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  return transactions.value
    .filter(t => t.type === 'debit' && new Date(t.created_at) >= startOfMonth)
    .reduce((sum, t) => sum + t.amount, 0)
})

const dailySpent = computed(() => {
  const today = new Date().toDateString()
  return transactions.value
    .filter(t => t.type === 'debit' && new Date(t.created_at).toDateString() === today)
    .reduce((sum, t) => sum + t.amount, 0)
})

const transactionCount = computed(() => transactions.value.length)

const cashbackEarned = computed(() => {
  return transactions.value
    .filter(t => t.type === 'debit')
    .reduce((sum, t) => sum + (t.amount * 0.01), 0) // 1% cashback
})

const filteredTransactions = computed(() => {
  if (!transactionFilter.value) return transactions.value
  return transactions.value.filter(t => t.category === transactionFilter.value)
})

// Methods
const getCardGradient = (type) => {
  const gradients = {
    'virtual': 'bg-gradient-to-r from-blue-500 to-blue-600',
    'physical': 'bg-gradient-to-r from-purple-500 to-purple-600',
    'premium': 'bg-gradient-to-r from-gray-800 to-black'
  }
  return gradients[type] || 'bg-gradient-to-r from-gray-400 to-gray-500'
}

const formatCardNumber = (number) => {
  if (!number) return '•••• •••• •••• ••••'
  return '•••• •••• •••• ' + number.slice(-4)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const getTransactionStatusColor = (status) => {
  const colors = {
    'completed': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'declined': 'bg-red-100 text-red-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const freezeCard = async () => {
  try {
    await $fetch(`/api/cards/${cardId}/freeze`, { method: 'POST' })
    card.value.status = 'frozen'
    alert('Card frozen successfully')
  } catch (error) {
    console.error('Error freezing card:', error)
    alert('Failed to freeze card')
  }
}

const toggleSetting = async (setting) => {
  try {
    card.value.settings[setting] = !card.value.settings[setting]
    await $fetch(`/api/cards/${cardId}/settings`, {
      method: 'PUT',
      body: { [setting]: card.value.settings[setting] }
    })
  } catch (error) {
    console.error('Error updating setting:', error)
    card.value.settings[setting] = !card.value.settings[setting] // Revert
  }
}

const editLimit = (type) => {
  const currentLimit = type === 'daily' ? card.value.daily_limit : card.value.monthly_limit
  const newLimit = prompt(`Enter new ${type} limit:`, currentLimit)
  
  if (newLimit && !isNaN(newLimit)) {
    updateLimit(type, parseFloat(newLimit))
  }
}

const updateLimit = async (type, amount) => {
  try {
    await $fetch(`/api/cards/${cardId}/limits`, {
      method: 'PUT',
      body: { [type + '_limit']: amount }
    })
    
    if (type === 'daily') {
      card.value.daily_limit = amount
    } else {
      card.value.monthly_limit = amount
    }
  } catch (error) {
    console.error('Error updating limit:', error)
    alert('Failed to update limit')
  }
}

const submitTopUp = async () => {
  try {
    await $fetch(`/api/cards/${cardId}/topup`, {
      method: 'POST',
      body: topUp.value
    })

    card.value.balance += topUp.value.amount
    topUp.value = { source_wallet: '', amount: 0 }
    showTopUp.value = false

    alert('Card topped up successfully!')
  } catch (error) {
    console.error('Error topping up card:', error)
    alert('Failed to top up card')
  }
}

const exportTransactions = () => {
  // Implementation for exporting transactions
  alert('Export functionality would be implemented here')
}

// Fetch data
const fetchCardDetails = async () => {
  try {
    const data = await $fetch(`/api/cards/${cardId}`)
    card.value = data.card
  } catch (error) {
    console.error('Error fetching card details:', error)
    // Fallback demo data
    card.value = {
      id: cardId,
      type: 'physical',
      card_number: '1234567812345678',
      expiry_date: '12/26',
      cvv: '123',
      status: 'active',
      balance: 2500.00,
      daily_limit: 500.00,
      monthly_limit: 5000.00,
      settings: {
        online_purchases: true,
        atm_withdrawals: true,
        international: false
      }
    }
  }
}

const fetchTransactions = async () => {
  try {
    const data = await $fetch(`/api/cards/${cardId}/transactions`)
    transactions.value = data.transactions
  } catch (error) {
    console.error('Error fetching transactions:', error)
    // Fallback demo data
    transactions.value = [
      {
        id: '1',
        merchant_name: 'Amazon',
        merchant_location: 'Online',
        amount: 89.99,
        type: 'debit',
        category: 'shopping',
        status: 'completed',
        created_at: new Date().toISOString()
      }
    ]
  }
}

const fetchAvailableWallets = async () => {
  try {
    const data = await $fetch('/api/wallet/list')
    availableWallets.value = data.wallets
  } catch (error) {
    console.error('Error fetching wallets:', error)
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchCardDetails(),
    fetchTransactions(),
    fetchAvailableWallets()
  ])
})
</script>