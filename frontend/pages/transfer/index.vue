<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Transfers</h1>
          <div class="flex items-center space-x-4">
            <button @click="showNewTransfer = true" 
                    class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
              New Transfer
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Transfer Options -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
        
        <!-- Crypto Transfer -->
        <div @click="selectTransferType('crypto')" 
             class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow cursor-pointer group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-orange-100 rounded-lg">
              <svg class="w-8 h-8 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Crypto Transfer</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Send Bitcoin, Ethereum and other cryptocurrencies to any wallet address worldwide.
          </p>
          <div class="text-sm text-gray-500">
            <div class="flex justify-between mb-1">
              <span>Network Fee:</span>
              <span class="text-green-600">Dynamic</span>
            </div>
            <div class="flex justify-between">
              <span>Speed:</span>
              <span>1-60 minutes</span>
            </div>
          </div>
        </div>

        <!-- Fiat Transfer -->
        <div @click="selectTransferType('fiat')" 
             class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow cursor-pointer group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-green-100 rounded-lg">
              <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Bank Transfer</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Send traditional currencies via SEPA, SWIFT, or domestic bank transfers.
          </p>
          <div class="text-sm text-gray-500">
            <div class="flex justify-between mb-1">
              <span>Fee:</span>
              <span class="text-green-600">From $2.50</span>
            </div>
            <div class="flex justify-between">
              <span>Speed:</span>
              <span>1-3 business days</span>
            </div>
          </div>
        </div>

        <!-- Instant Transfer -->
        <div @click="selectTransferType('instant')" 
             class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow cursor-pointer group">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-purple-100 rounded-lg">
              <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M13 10V3L4 14h7v7l9-11h-7z"></path>
              </svg>
            </div>
            <h3 class="ml-4 text-xl font-semibold text-gray-900">Instant Transfer</h3>
          </div>
          <p class="text-gray-600 mb-4">
            Instant transfers between Crypto Bank users with zero fees.
          </p>
          <div class="text-sm text-gray-500">
            <div class="flex justify-between mb-1">
              <span>Fee:</span>
              <span class="text-green-600">Free</span>
            </div>
            <div class="flex justify-between">
              <span>Speed:</span>
              <span>Instant</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transfers -->
      <div class="bg-white rounded-xl shadow-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900">Recent Transfers</h2>
            <div class="flex space-x-2">
              <select v-model="filterStatus" class="text-sm border-gray-300 rounded-md">
                <option value="">All Status</option>
                <option value="completed">Completed</option>
                <option value="pending">Pending</option>
                <option value="failed">Failed</option>
              </select>
            </div>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">To</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="transfer in filteredTransfers" :key="transfer.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(transfer.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full"
                        :class="getTypeColor(transfer.type)">
                    {{ transfer.type }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ transfer.amount }} {{ transfer.currency }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ transfer.recipient_name || transfer.recipient_address?.substring(0, 20) + '...' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-full"
                        :class="getStatusColor(transfer.status)">
                    {{ transfer.status }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button @click="viewTransfer(transfer)" 
                          class="text-blue-600 hover:text-blue-900">
                    View
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- New Transfer Modal -->
    <div v-if="showNewTransfer" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-2xl w-full mx-4 max-h-screen overflow-y-auto">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold">New Transfer</h3>
          <button @click="showNewTransfer = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Transfer Type Selection -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-3">Transfer Type</label>
          <div class="grid grid-cols-3 gap-3">
            <button v-for="type in transferTypes" :key="type.value"
                    @click="newTransfer.type = type.value"
                    :class="newTransfer.type === type.value 
                      ? 'border-blue-500 bg-blue-50 text-blue-700' 
                      : 'border-gray-300 text-gray-700'"
                    class="p-3 border-2 rounded-lg text-sm font-medium hover:border-blue-300">
              {{ type.label }}
            </button>
          </div>
        </div>

        <!-- Transfer Form -->
        <form @submit.prevent="submitTransfer" class="space-y-6">
          <!-- From Wallet -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">From Wallet</label>
            <select v-model="newTransfer.from_wallet" required
                    class="w-full rounded-md border-gray-300 shadow-sm">
              <option value="">Select wallet</option>
              <option v-for="wallet in userWallets" :key="wallet.id" :value="wallet.id">
                {{ wallet.currency }} - {{ wallet.balance }} {{ wallet.currency }}
              </option>
            </select>
          </div>

          <!-- Amount -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Amount</label>
            <input v-model.number="newTransfer.amount" type="number" step="0.00000001" required
                   class="w-full rounded-md border-gray-300 shadow-sm" 
                   placeholder="0.00000000">
          </div>

          <!-- Recipient (varies by type) -->
          <div v-if="newTransfer.type === 'crypto'">
            <label class="block text-sm font-medium text-gray-700 mb-2">Recipient Address</label>
            <input v-model="newTransfer.recipient_address" type="text" required
                   class="w-full rounded-md border-gray-300 shadow-sm" 
                   placeholder="Enter wallet address">
          </div>

          <div v-if="newTransfer.type === 'fiat'">
            <label class="block text-sm font-medium text-gray-700 mb-2">Bank Details</label>
            <div class="grid grid-cols-2 gap-4">
              <input v-model="newTransfer.bank_account" type="text" required
                     class="rounded-md border-gray-300 shadow-sm" 
                     placeholder="Account Number/IBAN">
              <input v-model="newTransfer.bank_code" type="text" required
                     class="rounded-md border-gray-300 shadow-sm" 
                     placeholder="Bank Code/SWIFT">
            </div>
            <input v-model="newTransfer.recipient_name" type="text" required
                   class="w-full mt-3 rounded-md border-gray-300 shadow-sm" 
                   placeholder="Recipient Name">
          </div>

          <div v-if="newTransfer.type === 'instant'">
            <label class="block text-sm font-medium text-gray-700 mb-2">Recipient</label>
            <input v-model="newTransfer.recipient_email" type="email" required
                   class="w-full rounded-md border-gray-300 shadow-sm" 
                   placeholder="user@email.com">
          </div>

          <!-- Notes -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Notes (Optional)</label>
            <textarea v-model="newTransfer.notes" rows="3"
                      class="w-full rounded-md border-gray-300 shadow-sm" 
                      placeholder="Payment reference or notes"></textarea>
          </div>

          <!-- Fee Information -->
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex justify-between items-center">
              <span class="text-sm font-medium text-gray-700">Estimated Fee:</span>
              <span class="text-sm font-semibold">{{ estimatedFee }} {{ selectedCurrency }}</span>
            </div>
            <div class="flex justify-between items-center mt-2">
              <span class="text-sm font-medium text-gray-700">You will send:</span>
              <span class="text-sm font-semibold">{{ totalAmount }} {{ selectedCurrency }}</span>
            </div>
          </div>

          <!-- Submit Buttons -->
          <div class="flex space-x-3">
            <button type="button" @click="showNewTransfer = false" 
                    class="flex-1 bg-gray-300 text-gray-700 py-3 px-4 rounded-lg hover:bg-gray-400">
              Cancel
            </button>
            <button type="submit" 
                    class="flex-1 bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700">
              Send Transfer
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
const showNewTransfer = ref(false)
const filterStatus = ref('')
const transfers = ref([])
const userWallets = ref([])

const newTransfer = ref({
  type: 'crypto',
  from_wallet: '',
  amount: 0,
  recipient_address: '',
  recipient_email: '',
  recipient_name: '',
  bank_account: '',
  bank_code: '',
  notes: ''
})

const transferTypes = [
  { value: 'crypto', label: 'Crypto' },
  { value: 'fiat', label: 'Bank' },
  { value: 'instant', label: 'Instant' }
]

// Computed
const filteredTransfers = computed(() => {
  if (!filterStatus.value) return transfers.value
  return transfers.value.filter(t => t.status === filterStatus.value)
})

const selectedCurrency = computed(() => {
  const wallet = userWallets.value.find(w => w.id === newTransfer.value.from_wallet)
  return wallet?.currency || 'USD'
})

const estimatedFee = computed(() => {
  const amount = newTransfer.value.amount || 0
  switch (newTransfer.value.type) {
    case 'crypto': return (amount * 0.0025).toFixed(8) // 0.25%
    case 'fiat': return Math.max(2.5, amount * 0.001).toFixed(2) // Min $2.50 or 0.1%
    case 'instant': return '0.00'
    default: return '0.00'
  }
})

const totalAmount = computed(() => {
  return (parseFloat(newTransfer.value.amount || 0) + parseFloat(estimatedFee.value)).toFixed(8)
})

// Methods
const selectTransferType = (type) => {
  newTransfer.value.type = type
  showNewTransfer.value = true
}

const getTypeColor = (type) => {
  const colors = {
    'crypto': 'bg-orange-100 text-orange-800',
    'fiat': 'bg-green-100 text-green-800',
    'instant': 'bg-purple-100 text-purple-800'
  }
  return colors[type] || 'bg-gray-100 text-gray-800'
}

const getStatusColor = (status) => {
  const colors = {
    'completed': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'failed': 'bg-red-100 text-red-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

const submitTransfer = async () => {
  try {
    const response = await $fetch('/api/transfer/create', {
      method: 'POST',
      body: newTransfer.value
    })

    // Reset form and close modal
    newTransfer.value = {
      type: 'crypto',
      from_wallet: '',
      amount: 0,
      recipient_address: '',
      recipient_email: '',
      recipient_name: '',
      bank_account: '',
      bank_code: '',
      notes: ''
    }
    showNewTransfer.value = false

    // Refresh transfers list
    await fetchTransfers()

    alert('Transfer submitted successfully!')
  } catch (error) {
    console.error('Transfer error:', error)
    alert('Transfer failed. Please try again.')
  }
}

const viewTransfer = (transfer) => {
  navigateTo(`/transfer/${transfer.id}`)
}

// Fetch data
const fetchTransfers = async () => {
  try {
    const data = await $fetch('/api/transfer/history')
    transfers.value = data.transfers || []
  } catch (error) {
    console.error('Error fetching transfers:', error)
    // Fallback demo data
    transfers.value = [
      {
        id: '1',
        type: 'crypto',
        amount: 0.1,
        currency: 'BTC',
        recipient_address: '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
        status: 'completed',
        created_at: new Date().toISOString()
      },
      {
        id: '2',
        type: 'fiat',
        amount: 1000,
        currency: 'USD',
        recipient_name: 'John Doe',
        status: 'pending',
        created_at: new Date().toISOString()
      }
    ]
  }
}

const fetchUserWallets = async () => {
  try {
    const data = await $fetch('/api/wallet/list')
    userWallets.value = data.wallets || []
  } catch (error) {
    console.error('Error fetching wallets:', error)
    // Fallback demo data
    userWallets.value = [
      { id: '1', currency: 'BTC', balance: 0.5 },
      { id: '2', currency: 'ETH', balance: 10.0 },
      { id: '3', currency: 'USD', balance: 5000.0 }
    ]
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchTransfers(),
    fetchUserWallets()
  ])
})
</script>