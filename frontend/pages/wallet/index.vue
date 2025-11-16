<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Digital Wallets</h1>
          <div class="flex items-center space-x-4">
            <div class="text-right">
              <div class="text-sm text-gray-500">Total Portfolio Value</div>
              <div class="text-2xl font-bold text-gray-900">${{ totalValue.toLocaleString() }}</div>
            </div>
            <button @click="showCreateWallet = true" 
                    class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
              Create Wallet
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      
      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <button @click="navigateTo('/transfer')" 
                class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow text-left group">
          <div class="flex items-center">
            <div class="p-3 bg-blue-100 rounded-lg group-hover:bg-blue-200">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="font-semibold text-gray-900">Send</h3>
              <p class="text-sm text-gray-500">Transfer funds</p>
            </div>
          </div>
        </button>

        <button @click="showReceive = true" 
                class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow text-left group">
          <div class="flex items-center">
            <div class="p-3 bg-green-100 rounded-lg group-hover:bg-green-200">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M7 16l-4-4m0 0l4-4m-4 4h18"></path>
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="font-semibold text-gray-900">Receive</h3>
              <p class="text-sm text-gray-500">Get funds</p>
            </div>
          </div>
        </button>

        <button @click="navigateTo('/exchange')" 
                class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow text-left group">
          <div class="flex items-center">
            <div class="p-3 bg-purple-100 rounded-lg group-hover:bg-purple-200">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"></path>
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="font-semibold text-gray-900">Exchange</h3>
              <p class="text-sm text-gray-500">Swap currencies</p>
            </div>
          </div>
        </button>

        <button @click="showBuyModal = true" 
                class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow text-left group">
          <div class="flex items-center">
            <div class="p-3 bg-yellow-100 rounded-lg group-hover:bg-yellow-200">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="font-semibold text-gray-900">Buy</h3>
              <p class="text-sm text-gray-500">Add funds</p>
            </div>
          </div>
        </button>
      </div>

      <!-- Wallet Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-8">
        <div v-for="wallet in wallets" :key="wallet.id" 
             class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow">
          
          <!-- Wallet Header -->
          <div class="p-6 border-b border-gray-100">
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <div class="w-12 h-12 rounded-full flex items-center justify-center text-white font-bold text-lg"
                     :class="getCurrencyColor(wallet.currency)">
                  {{ getCurrencySymbol(wallet.currency) }}
                </div>
                <div class="ml-4">
                  <h3 class="text-lg font-semibold text-gray-900">{{ wallet.currency }}</h3>
                  <p class="text-sm text-gray-500">{{ getCurrencyName(wallet.currency) }}</p>
                </div>
              </div>
              <div class="text-right">
                <div class="text-2xl font-bold text-gray-900">
                  {{ formatBalance(wallet.balance, wallet.currency) }}
                </div>
                <div class="text-sm text-gray-500">
                  ≈ ${{ (wallet.balance * (wallet.usd_rate || 1)).toLocaleString() }}
                </div>
              </div>
            </div>
          </div>

          <!-- Wallet Actions -->
          <div class="p-6">
            <div class="grid grid-cols-2 gap-3 mb-4">
              <button @click="sendFunds(wallet)" 
                      class="bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 text-sm">
                Send
              </button>
              <button @click="receiveFunds(wallet)" 
                      class="bg-gray-100 text-gray-700 py-2 px-4 rounded-lg hover:bg-gray-200 text-sm">
                Receive
              </button>
            </div>

            <!-- Recent Transactions -->
            <div class="space-y-2">
              <h4 class="text-sm font-medium text-gray-700 mb-2">Recent Activity</h4>
              <div v-if="getWalletTransactions(wallet.id).length === 0" 
                   class="text-center py-4 text-sm text-gray-500">
                No recent transactions
              </div>
              <div v-else v-for="tx in getWalletTransactions(wallet.id).slice(0, 3)" 
                   :key="tx.id" class="flex items-center justify-between py-2">
                <div class="flex items-center">
                  <div class="w-2 h-2 rounded-full mr-2"
                       :class="tx.type === 'credit' ? 'bg-green-500' : 'bg-red-500'"></div>
                  <div class="text-sm">
                    <div class="font-medium">{{ tx.type === 'credit' ? 'Received' : 'Sent' }}</div>
                    <div class="text-xs text-gray-500">{{ formatDate(tx.created_at) }}</div>
                  </div>
                </div>
                <div class="text-sm font-medium"
                     :class="tx.type === 'credit' ? 'text-green-600' : 'text-red-600'">
                  {{ tx.type === 'credit' ? '+' : '-' }}{{ formatBalance(tx.amount, wallet.currency) }}
                </div>
              </div>
              <button @click="viewWalletHistory(wallet)" 
                      class="w-full text-blue-600 text-sm hover:text-blue-800 mt-2">
                View All Transactions →
              </button>
            </div>
          </div>
        </div>

        <!-- Add Wallet Card -->
        <div @click="showCreateWallet = true" 
             class="border-2 border-dashed border-gray-300 rounded-xl p-6 flex flex-col items-center justify-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer group">
          <div class="w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center group-hover:bg-blue-200 mb-4">
            <svg class="w-6 h-6 text-gray-400 group-hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-900 group-hover:text-blue-900 mb-2">Add Wallet</h3>
          <p class="text-sm text-gray-500 text-center group-hover:text-blue-700">
            Create a new wallet for cryptocurrencies or fiat currencies
          </p>
        </div>
      </div>

      <!-- Portfolio Overview -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        
        <!-- Portfolio Allocation -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Portfolio Allocation</h3>
          <div class="space-y-4">
            <div v-for="wallet in wallets" :key="wallet.id" 
                 class="flex items-center justify-between">
              <div class="flex items-center">
                <div class="w-3 h-3 rounded-full mr-3"
                     :class="getCurrencyColor(wallet.currency)"></div>
                <span class="font-medium">{{ wallet.currency }}</span>
              </div>
              <div class="text-right">
                <div class="font-medium">${{ (wallet.balance * (wallet.usd_rate || 1)).toLocaleString() }}</div>
                <div class="text-sm text-gray-500">{{ getPercentage(wallet) }}%</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Market Prices -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Market Prices</h3>
          <div class="space-y-3">
            <div v-for="price in marketPrices" :key="price.currency" 
                 class="flex items-center justify-between py-2">
              <div class="flex items-center">
                <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold mr-3"
                     :class="getCurrencyColor(price.currency)">
                  {{ getCurrencySymbol(price.currency) }}
                </div>
                <div>
                  <div class="font-medium">{{ price.currency }}</div>
                  <div class="text-sm text-gray-500">{{ getCurrencyName(price.currency) }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-medium">${{ price.price?.toLocaleString() }}</div>
                <div class="text-sm" :class="price.change >= 0 ? 'text-green-600' : 'text-red-600'">
                  {{ price.change >= 0 ? '+' : '' }}{{ price.change }}%
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Wallet Modal -->
    <div v-if="showCreateWallet" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold">Create New Wallet</h3>
          <button @click="showCreateWallet = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form @submit.prevent="createWallet" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Currency</label>
            <select v-model="newWallet.currency" required
                    class="w-full rounded-md border-gray-300 shadow-sm">
              <option value="">Select currency</option>
              <optgroup label="Cryptocurrencies">
                <option value="BTC">Bitcoin (BTC)</option>
                <option value="ETH">Ethereum (ETH)</option>
                <option value="LTC">Litecoin (LTC)</option>
                <option value="ADA">Cardano (ADA)</option>
                <option value="DOT">Polkadot (DOT)</option>
                <option value="XRP">Ripple (XRP)</option>
              </optgroup>
              <optgroup label="Fiat Currencies">
                <option value="USD">US Dollar (USD)</option>
                <option value="EUR">Euro (EUR)</option>
                <option value="GBP">British Pound (GBP)</option>
                <option value="JPY">Japanese Yen (JPY)</option>
                <option value="CAD">Canadian Dollar (CAD)</option>
              </optgroup>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Wallet Name (Optional)</label>
            <input v-model="newWallet.name" type="text"
                   class="w-full rounded-md border-gray-300 shadow-sm" 
                   placeholder="My Bitcoin Wallet">
          </div>

          <div class="bg-blue-50 rounded-lg p-4">
            <h4 class="text-sm font-medium text-blue-900 mb-2">Security Notice</h4>
            <p class="text-sm text-blue-800">
              Your wallet will be secured with industry-standard encryption. 
              Make sure to backup your recovery phrase safely.
            </p>
          </div>

          <div class="flex space-x-3">
            <button type="button" @click="showCreateWallet = false" 
                    class="flex-1 bg-gray-300 text-gray-700 py-3 px-4 rounded-lg hover:bg-gray-400">
              Cancel
            </button>
            <button type="submit" 
                    class="flex-1 bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700">
              Create Wallet
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Receive Modal -->
    <div v-if="showReceive" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold">Receive Funds</h3>
          <button @click="showReceive = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <div class="text-center">
          <div class="mb-4">
            <select v-model="selectedReceiveWallet" class="w-full rounded-md border-gray-300 shadow-sm">
              <option value="">Select wallet</option>
              <option v-for="wallet in wallets" :key="wallet.id" :value="wallet">
                {{ wallet.currency }} Wallet
              </option>
            </select>
          </div>

          <div v-if="selectedReceiveWallet" class="space-y-4">
            <div class="bg-white p-4 rounded-lg border-2 border-gray-200">
              <!-- QR Code placeholder -->
              <div class="w-48 h-48 mx-auto bg-gray-100 rounded-lg flex items-center justify-center mb-4">
                <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M12 12h-6.01M12 12v2m-6-8h2m12 0h2"></path>
                </svg>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Wallet Address</label>
              <div class="flex">
                <input :value="selectedReceiveWallet.address" readonly
                       class="flex-1 rounded-l-md border-gray-300 bg-gray-50">
                <button @click="copyToClipboard(selectedReceiveWallet.address)"
                        class="bg-blue-600 text-white px-4 rounded-r-md hover:bg-blue-700">
                  Copy
                </button>
              </div>
            </div>

            <p class="text-sm text-gray-500">
              Send {{ selectedReceiveWallet.currency }} to this address to receive funds in your wallet.
            </p>
          </div>
        </div>
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
const showCreateWallet = ref(false)
const showReceive = ref(false)
const showBuyModal = ref(false)
const selectedReceiveWallet = ref(null)
const wallets = ref([])
const transactions = ref([])
const marketPrices = ref([])

const newWallet = ref({
  currency: '',
  name: ''
})

// Computed
const totalValue = computed(() => {
  return wallets.value.reduce((total, wallet) => {
    return total + (wallet.balance * (wallet.usd_rate || 1))
  }, 0)
})

// Methods
const getCurrencyColor = (currency) => {
  const colors = {
    'BTC': 'bg-orange-500',
    'ETH': 'bg-blue-500',
    'LTC': 'bg-gray-400',
    'ADA': 'bg-blue-600',
    'DOT': 'bg-pink-500',
    'XRP': 'bg-blue-400',
    'USD': 'bg-green-500',
    'EUR': 'bg-blue-800',
    'GBP': 'bg-red-600',
    'JPY': 'bg-red-500',
    'CAD': 'bg-red-700'
  }
  return colors[currency] || 'bg-gray-500'
}

const getCurrencySymbol = (currency) => {
  const symbols = {
    'BTC': '₿',
    'ETH': 'Ξ',
    'LTC': 'Ł',
    'ADA': '₳',
    'DOT': '●',
    'XRP': '✕',
    'USD': '$',
    'EUR': '€',
    'GBP': '£',
    'JPY': '¥',
    'CAD': 'C$'
  }
  return symbols[currency] || currency.substring(0, 2)
}

const getCurrencyName = (currency) => {
  const names = {
    'BTC': 'Bitcoin',
    'ETH': 'Ethereum',
    'LTC': 'Litecoin',
    'ADA': 'Cardano',
    'DOT': 'Polkadot',
    'XRP': 'Ripple',
    'USD': 'US Dollar',
    'EUR': 'Euro',
    'GBP': 'British Pound',
    'JPY': 'Japanese Yen',
    'CAD': 'Canadian Dollar'
  }
  return names[currency] || currency
}

const formatBalance = (balance, currency) => {
  if (['USD', 'EUR', 'GBP', 'JPY', 'CAD'].includes(currency)) {
    return balance.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }
  return balance.toFixed(8)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const getPercentage = (wallet) => {
  const walletValue = wallet.balance * (wallet.usd_rate || 1)
  return totalValue.value > 0 ? ((walletValue / totalValue.value) * 100).toFixed(1) : 0
}

const getWalletTransactions = (walletId) => {
  return transactions.value.filter(tx => tx.wallet_id === walletId)
}

const createWallet = async () => {
  try {
    const response = await $fetch('/api/wallet/create', {
      method: 'POST',
      body: newWallet.value
    })

    // Reset form and close modal
    newWallet.value = { currency: '', name: '' }
    showCreateWallet.value = false

    // Refresh wallets
    await fetchWallets()

    alert('Wallet created successfully!')
  } catch (error) {
    console.error('Wallet creation error:', error)
    alert('Failed to create wallet. Please try again.')
  }
}

const sendFunds = (wallet) => {
  navigateTo(`/transfer?from=${wallet.id}`)
}

const receiveFunds = (wallet) => {
  selectedReceiveWallet.value = wallet
  showReceive.value = true
}

const viewWalletHistory = (wallet) => {
  navigateTo(`/wallet/${wallet.id}/history`)
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('Address copied to clipboard!')
  } catch (error) {
    console.error('Copy failed:', error)
  }
}

// Fetch data
const fetchWallets = async () => {
  try {
    const data = await $fetch('/api/wallet/list')
    wallets.value = data.wallets || []
  } catch (error) {
    console.error('Error fetching wallets:', error)
    // Fallback demo data
    wallets.value = [
      {
        id: '1',
        currency: 'BTC',
        balance: 0.25,
        usd_rate: 43500,
        address: '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa'
      },
      {
        id: '2',
        currency: 'ETH',
        balance: 5.0,
        usd_rate: 2450,
        address: '0x742d35Cc6634C0532925a3b8d6C4a3d8a7b7d8db'
      },
      {
        id: '3',
        currency: 'USD',
        balance: 2500.00,
        usd_rate: 1,
        address: 'USD-WALLET-001'
      }
    ]
  }
}

const fetchTransactions = async () => {
  try {
    const data = await $fetch('/api/wallet/transactions')
    transactions.value = data.transactions || []
  } catch (error) {
    console.error('Error fetching transactions:', error)
  }
}

const fetchMarketPrices = async () => {
  try {
    const data = await $fetch('/api/market/prices')
    marketPrices.value = data.prices || []
  } catch (error) {
    console.error('Error fetching market prices:', error)
    // Fallback demo data
    marketPrices.value = [
      { currency: 'BTC', price: 43500, change: 2.3 },
      { currency: 'ETH', price: 2450, change: -1.2 },
      { currency: 'LTC', price: 72, change: 0.8 },
      { currency: 'ADA', price: 0.52, change: 3.1 }
    ]
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchWallets(),
    fetchTransactions(),
    fetchMarketPrices()
  ])
})
</script>