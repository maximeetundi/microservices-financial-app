<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Orders</h1>
          <div class="flex items-center space-x-4">
            <NuxtLink to="/exchange/trading" 
                      class="bg-blue-600 text-gray-900 px-4 py-2 rounded-lg hover:bg-blue-700">
              New Order
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filter Tabs -->
      <div class="bg-white rounded-xl shadow-lg p-6 mb-8">
        <div class="flex space-x-6 border-b">
          <button v-for="filter in filters" :key="filter.key"
                  @click="activeFilter = filter.key"
                  :class="activeFilter === filter.key 
                    ? 'border-blue-500 text-blue-600 border-b-2 pb-2' 
                    : 'text-gray-500 hover:text-gray-700 pb-2'">
            {{ filter.label }}
            <span v-if="filter.count" class="ml-2 px-2 py-1 text-xs rounded-full bg-gray-200">
              {{ filter.count }}
            </span>
          </button>
        </div>

        <!-- Search and Filters -->
        <div class="mt-4 flex items-center space-x-4">
          <div class="flex-1">
            <input v-model="searchQuery" type="text" placeholder="Search orders..."
                   class="w-full rounded-md border-gray-300 shadow-sm">
          </div>
          <select v-model="pairFilter" class="rounded-md border-gray-300 shadow-sm">
            <option value="">All Pairs</option>
            <option value="BTC/USD">BTC/USD</option>
            <option value="ETH/USD">ETH/USD</option>
            <option value="EUR/USD">EUR/USD</option>
          </select>
          <select v-model="typeFilter" class="rounded-md border-gray-300 shadow-sm">
            <option value="">All Types</option>
            <option value="market">Market</option>
            <option value="limit">Limit</option>
            <option value="stop_loss">Stop Loss</option>
          </select>
        </div>
      </div>

      <!-- Orders Table -->
      <div class="bg-white rounded-xl shadow-lg overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Pair</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Side</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Filled</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="order in filteredOrders" :key="order.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(order.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ order.pair }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="[
                    'px-2 py-1 text-xs rounded-full',
                    order.side === 'buy' ? 'text-green-600 bg-green-100' : 'text-red-600 bg-red-100'
                  ]">
                    {{ order.side.toUpperCase() }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ order.order_type }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ order.amount }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ order.price ? '$' + order.price.toLocaleString() : 'Market' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ order.filled_amount || 0 }} / {{ order.amount }}
                  <div class="w-full bg-gray-200 rounded-full h-1 mt-1">
                    <div class="bg-blue-600 h-1 rounded-full" 
                         :style="{ width: (order.filled_amount / order.amount * 100) + '%' }"></div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getStatusColor(order.status)"
                        class="px-2 py-1 text-xs rounded-full">
                    {{ order.status }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button v-if="order.status === 'pending'" 
                          @click="cancelOrder(order.id)"
                          class="text-red-600 hover:text-red-900">
                    Cancel
                  </button>
                  <button @click="viewOrderDetails(order)"
                          class="text-blue-600 hover:text-blue-900 ml-2">
                    Details
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="mt-6 flex items-center justify-between">
        <div class="text-sm text-gray-500">
          Showing {{ filteredOrders.length }} of {{ allOrders.length }} orders
        </div>
        <div class="flex space-x-2">
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50">Previous</button>
          <button class="px-3 py-1 text-sm bg-blue-600 text-gray-900 rounded">1</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50">2</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50">Next</button>
        </div>
      </div>
    </div>

    <!-- Order Details Modal -->
    <div v-if="selectedOrder" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-lg w-full mx-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold">Order Details</h3>
          <button @click="selectedOrder = null" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-500">Order ID</label>
              <div class="text-sm text-gray-900">{{ selectedOrder.id }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Status</label>
              <span :class="getStatusColor(selectedOrder.status)" 
                    class="px-2 py-1 text-xs rounded-full">
                {{ selectedOrder.status }}
              </span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-500">Pair</label>
              <div class="text-sm text-gray-900">{{ selectedOrder.pair }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Side</label>
              <div class="text-sm text-gray-900">{{ selectedOrder.side }}</div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-500">Amount</label>
              <div class="text-sm text-gray-900">{{ selectedOrder.amount }}</div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Price</label>
              <div class="text-sm text-gray-900">
                {{ selectedOrder.price ? '$' + selectedOrder.price.toLocaleString() : 'Market Price' }}
              </div>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-500">Created At</label>
            <div class="text-sm text-gray-900">{{ formatDate(selectedOrder.created_at) }}</div>
          </div>
        </div>

        <div class="mt-6 flex space-x-3">
          <button v-if="selectedOrder.status === 'pending'" 
                  @click="cancelOrder(selectedOrder.id); selectedOrder = null"
                  class="flex-1 bg-red-600 text-gray-900 py-2 px-4 rounded-lg hover:bg-red-700">
            Cancel Order
          </button>
          <button @click="selectedOrder = null" 
                  class="flex-1 bg-gray-300 text-gray-700 py-2 px-4 rounded-lg hover:bg-gray-400">
            Close
          </button>
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
const allOrders = ref([])
const activeFilter = ref('all')
const searchQuery = ref('')
const pairFilter = ref('')
const typeFilter = ref('')
const selectedOrder = ref(null)

const filters = computed(() => [
  { key: 'all', label: 'All Orders', count: allOrders.value.length },
  { key: 'pending', label: 'Pending', count: allOrders.value.filter(o => o.status === 'pending').length },
  { key: 'filled', label: 'Filled', count: allOrders.value.filter(o => o.status === 'filled').length },
  { key: 'cancelled', label: 'Cancelled', count: allOrders.value.filter(o => o.status === 'cancelled').length }
])

const filteredOrders = computed(() => {
  let filtered = allOrders.value

  // Filter by status
  if (activeFilter.value !== 'all') {
    filtered = filtered.filter(order => order.status === activeFilter.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    filtered = filtered.filter(order => 
      order.pair.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      order.id.toLowerCase().includes(searchQuery.value.toLowerCase())
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

  return filtered
})

// Methods
const getStatusColor = (status) => {
  const colors = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'filled': 'bg-green-100 text-green-800',
    'cancelled': 'bg-red-100 text-red-800',
    'partially_filled': 'bg-blue-100 text-blue-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const cancelOrder = async (orderId) => {
  try {
    await $fetch(`/api/trading/orders/${orderId}/cancel`, { method: 'POST' })
    
    // Update local state
    const order = allOrders.value.find(o => o.id === orderId)
    if (order) {
      order.status = 'cancelled'
    }

    alert('Order cancelled successfully')
  } catch (error) {
    console.error('Error cancelling order:', error)
    alert('Error cancelling order')
  }
}

const viewOrderDetails = (order) => {
  selectedOrder.value = order
}

// Fetch orders
const fetchOrders = async () => {
  try {
    const { data } = await $fetch('/api/trading/orders')
    allOrders.value = data.orders || []
  } catch (error) {
    console.error('Error fetching orders:', error)
    // Fallback demo data
    allOrders.value = [
      {
        id: '1',
        pair: 'BTC/USD',
        side: 'buy',
        order_type: 'market',
        amount: 0.1,
        price: null,
        filled_amount: 0.1,
        status: 'filled',
        created_at: '2024-01-01T10:00:00Z'
      },
      {
        id: '2',
        pair: 'ETH/USD',
        side: 'sell',
        order_type: 'limit',
        amount: 2.5,
        price: 2500,
        filled_amount: 0,
        status: 'pending',
        created_at: '2024-01-01T11:00:00Z'
      }
    ]
  }
}

// Lifecycle
onMounted(async () => {
  await fetchOrders()
})
</script>