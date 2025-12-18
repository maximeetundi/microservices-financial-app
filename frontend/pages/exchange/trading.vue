<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Trading</h1>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-500">Total Balance:</span>
            <span class="text-lg font-semibold text-green-600">${{ portfolio.totalValue?.toLocaleString() || '0' }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <!-- Trading Interface -->
        <div class="lg:col-span-2 space-y-6">
          
          <!-- Market Tickers -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Market Tickers</h2>
            <div class="overflow-x-auto">
              <table class="min-w-full">
                <thead>
                  <tr class="text-left text-sm font-medium text-gray-500 border-b">
                    <th class="pb-2">Pair</th>
                    <th class="pb-2">Price</th>
                    <th class="pb-2">24h Change</th>
                    <th class="pb-2">24h Volume</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="ticker in tickers" :key="ticker.symbol" 
                      class="border-b hover:bg-gray-50 cursor-pointer"
                      @click="selectPair(ticker.symbol)">
                    <td class="py-3 font-medium">{{ ticker.symbol }}</td>
                    <td class="py-3">${{ ticker.price?.toLocaleString() }}</td>
                    <td class="py-3" :class="ticker.change_24h >= 0 ? 'text-green-600' : 'text-red-600'">
                      {{ ticker.change_24h >= 0 ? '+' : '' }}{{ ticker.change_24h }}%
                    </td>
                    <td class="py-3">${{ ticker.volume_24h?.toLocaleString() }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Trading Form -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-gray-900">Place Order</h2>
              <div class="text-sm text-gray-500">Selected: {{ selectedPair }}</div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Buy Section -->
              <div class="border border-green-200 rounded-lg p-4 bg-green-50">
                <h3 class="text-lg font-medium text-green-800 mb-4">Buy {{ selectedPair }}</h3>
                <form @submit.prevent="placeBuyOrder">
                  <div class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Order Type</label>
                      <select v-model="buyOrder.type" class="w-full rounded-md border-gray-300 shadow-sm">
                        <option value="market">Market</option>
                        <option value="limit">Limit</option>
                        <option value="stop_loss">Stop Loss</option>
                      </select>
                    </div>
                    
                    <div v-if="buyOrder.type === 'limit'">
                      <label class="block text-sm font-medium text-gray-700 mb-2">Price</label>
                      <input v-model.number="buyOrder.price" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <div v-if="buyOrder.type === 'stop_loss'">
                      <label class="block text-sm font-medium text-gray-700 mb-2">Stop Price</label>
                      <input v-model.number="buyOrder.stopPrice" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Amount</label>
                      <input v-model.number="buyOrder.amount" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <button type="submit" :disabled="!isValidBuyOrder" 
                            class="w-full bg-green-600 text-gray-900 py-2 px-4 rounded-md hover:bg-green-700 disabled:opacity-50">
                      {{ buyOrder.type === 'market' ? 'Buy Now' : 'Place Buy Order' }}
                    </button>
                  </div>
                </form>
              </div>

              <!-- Sell Section -->
              <div class="border border-red-200 rounded-lg p-4 bg-red-50">
                <h3 class="text-lg font-medium text-red-800 mb-4">Sell {{ selectedPair }}</h3>
                <form @submit.prevent="placeSellOrder">
                  <div class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Order Type</label>
                      <select v-model="sellOrder.type" class="w-full rounded-md border-gray-300 shadow-sm">
                        <option value="market">Market</option>
                        <option value="limit">Limit</option>
                        <option value="stop_loss">Stop Loss</option>
                      </select>
                    </div>

                    <div v-if="sellOrder.type === 'limit'">
                      <label class="block text-sm font-medium text-gray-700 mb-2">Price</label>
                      <input v-model.number="sellOrder.price" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <div v-if="sellOrder.type === 'stop_loss'">
                      <label class="block text-sm font-medium text-gray-700 mb-2">Stop Price</label>
                      <input v-model.number="sellOrder.stopPrice" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Amount</label>
                      <input v-model.number="sellOrder.amount" type="number" step="0.01" 
                             class="w-full rounded-md border-gray-300 shadow-sm" placeholder="0.00">
                    </div>

                    <button type="submit" :disabled="!isValidSellOrder" 
                            class="w-full bg-red-600 text-gray-900 py-2 px-4 rounded-md hover:bg-red-700 disabled:opacity-50">
                      {{ sellOrder.type === 'market' ? 'Sell Now' : 'Place Sell Order' }}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Portfolio Summary -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Portfolio</h2>
            <div class="space-y-3">
              <div v-for="holding in portfolio.holdings" :key="holding.currency" 
                   class="flex items-center justify-between">
                <div>
                  <div class="font-medium">{{ holding.currency }}</div>
                  <div class="text-sm text-gray-500">{{ holding.amount }}</div>
                </div>
                <div class="text-right">
                  <div class="font-medium">${{ holding.value?.toLocaleString() }}</div>
                  <div class="text-sm" :class="holding.change_24h >= 0 ? 'text-green-600' : 'text-red-600'">
                    {{ holding.change_24h >= 0 ? '+' : '' }}{{ holding.change_24h }}%
                  </div>
                </div>
              </div>
            </div>
            <NuxtLink to="/portfolio" class="block mt-4 text-blue-600 text-sm hover:text-blue-800">
              View Full Portfolio →
            </NuxtLink>
          </div>

          <!-- Recent Orders -->
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Recent Orders</h2>
            <div class="space-y-3">
              <div v-for="order in recentOrders" :key="order.id" 
                   class="flex items-center justify-between py-2 border-b last:border-b-0">
                <div>
                  <div class="font-medium text-sm">{{ order.pair }}</div>
                  <div class="text-xs text-gray-500">{{ order.side }} • {{ order.order_type }}</div>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium">{{ order.amount }}</div>
                  <div class="text-xs" :class="getStatusColor(order.status)">{{ order.status }}</div>
                </div>
              </div>
            </div>
            <button @click="viewAllOrders" class="block mt-4 text-blue-600 text-sm hover:text-blue-800">
              View All Orders →
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

import { ref, computed, onMounted } from 'vue'
import { exchangeAPI } from '~/composables/useApi'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

// Reactive data
const selectedPair = ref('BTC/USD')
const tickers = ref([])
const portfolio = ref({ totalValue: 0, holdings: [] })
const recentOrders = ref([])

const buyOrder = ref({
  type: 'market',
  amount: 0,
  price: 0,
  stopPrice: 0
})

const sellOrder = ref({
  type: 'market',
  amount: 0,
  price: 0,
  stopPrice: 0
})

// Computed properties
const isValidBuyOrder = computed(() => {
  return buyOrder.value.amount > 0 && 
         (buyOrder.value.type === 'market' || buyOrder.value.price > 0 || buyOrder.value.stopPrice > 0)
})

const isValidSellOrder = computed(() => {
  return sellOrder.value.amount > 0 && 
         (sellOrder.value.type === 'market' || sellOrder.value.price > 0 || sellOrder.value.stopPrice > 0)
})

// Methods
const selectPair = (pair) => {
  selectedPair.value = pair
}

const placeBuyOrder = async () => {
  try {
     const [base, quote] = selectedPair.value.split('/')
     const { data } = await exchangeAPI.buyCrypto(
         base, 
         buyOrder.value.amount, 
         'wallet', // Default payment method
         buyOrder.value.type,
         buyOrder.value.price
     )

    // Reset form
    buyOrder.value = { type: 'market', amount: 0, price: 0, stopPrice: 0 }
    
    // Refresh orders
    await fetchRecentOrders()
    
    // Show success message
    alert('Buy order placed successfully!')
  } catch (error) {
    console.error('Error placing buy order:', error)
    alert('Error placing buy order')
  }
}

const placeSellOrder = async () => {
  try {
    const [base, quote] = selectedPair.value.split('/')
    // NOTE: destinationWalletId needed for sellCrypto implies where fiat goes. 
    // Assuming backend handles it or we need a wallet selection in UI. 
    // Passing 'default' or similar if API allows, or a real ID if we fetched wallets.
    const { data } = await exchangeAPI.sellCrypto(
        base,
        sellOrder.value.amount,
        'wallet-id-placeholder', // TODO: User should select fiat wallet
        sellOrder.value.type,
        sellOrder.value.price
    )

    // Reset form
    sellOrder.value = { type: 'market', amount: 0, price: 0, stopPrice: 0 }
    
    // Refresh orders
    await fetchRecentOrders()
    
    // Show success message
    alert('Sell order placed successfully!')
  } catch (error) {
    console.error('Error placing sell order:', error)
    alert('Error placing sell order')
  }
}

const getStatusColor = (status) => {
  switch (status) {
    case 'filled': return 'text-green-600'
    case 'pending': return 'text-yellow-600'
    case 'cancelled': return 'text-red-600'
    default: return 'text-gray-600'
  }
}

const viewAllOrders = () => {
  navigateTo('/orders')
}

// Fetch data functions
const fetchTickers = async () => {
  try {
    const { data } = await exchangeAPI.getMarkets()
    if (data) {
        tickers.value = data.map(m => ({
            symbol: m.Symbol,
            price: m.Price,
            change_24h: m.Change24h,
            volume_24h: m.Volume24h
        }))
    }
  } catch (error) {
    console.error('Error fetching tickers:', error)
  }
}

const fetchPortfolio = async () => {
  try {
    const { data } = await exchangeAPI.getTradingPortfolio()
    if (data) {
        portfolio.value = {
            totalValue: data.TotalValue,
            holdings: data.Holdings.map(h => ({
                currency: h.Asset,
                amount: h.Amount,
                value: h.Value,
                change_24h: 0 // Backend might need to provide this
            }))
        }
    }
  } catch (error) {
    console.error('Error fetching portfolio:', error)
  }
}

const fetchRecentOrders = async () => {
  try {
    const { data } = await exchangeAPI.getOrders()
    if (data) {
        recentOrders.value = data.map(o => ({
            id: o.ID,
            pair: `${o.ToCurrency}/${o.FromCurrency}`, // Adjust based on Side
            side: o.Side,
            order_type: o.OrderType,
            amount: o.Amount,
            status: o.Status
        })).slice(0, 5)
    }
  } catch (error) {
    console.error('Error fetching recent orders:', error)
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchTickers(),
    fetchPortfolio(),
    fetchRecentOrders()
  ])
})