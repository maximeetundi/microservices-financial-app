<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto py-8 px-4">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">📦 Mes Commandes</h1>

      <!-- Loading -->
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 h-32 rounded-xl"></div>
      </div>

      <!-- Empty -->
      <div v-else-if="orders.length === 0" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl">
        <div class="text-6xl mb-4">📦</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucune commande</h3>
        <p class="text-gray-500 mb-6">Vos commandes apparaîtront ici</p>
        <NuxtLink to="/shops" class="btn-premium px-6 py-3">
          Explorer les boutiques
        </NuxtLink>
      </div>

      <!-- Orders List -->
      <div v-else class="space-y-4">
        <NuxtLink 
          v-for="order in orders" 
          :key="order.id"
          :to="`/orders/${order.id}`"
          class="block bg-white dark:bg-slate-900 rounded-xl p-6 border border-gray-100 dark:border-gray-800 hover:border-indigo-500/30 transition-all"
        >
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <div class="flex items-center gap-3 mb-2">
                <span class="text-lg font-bold text-gray-900 dark:text-white">{{ order.order_number }}</span>
                <span :class="getStatusClass(order.order_status)">
                  {{ getStatusLabel(order.order_status) }}
                </span>
                <span :class="getPaymentClass(order.payment_status)">
                  {{ getPaymentLabel(order.payment_status) }}
                </span>
              </div>
              <p class="text-gray-500">{{ order.shop_name }} • {{ formatDate(order.created_at) }}</p>
            </div>
            
            <div class="text-right">
              <div class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">
                {{ formatPrice(order.total_amount, order.currency) }}
              </div>
              <p class="text-sm text-gray-500">{{ order.items.length }} article(s)</p>
            </div>
          </div>
          
          <!-- Items Preview -->
          <div class="flex gap-2 mt-4 overflow-x-auto">
            <div v-for="item in order.items.slice(0, 4)" :key="item.product_id" class="w-12 h-12 rounded-lg bg-gray-100 dark:bg-slate-800 overflow-hidden flex-shrink-0">
              <img v-if="item.product_image" :src="item.product_image" class="w-full h-full object-cover" alt="">
              <div v-else class="w-full h-full flex items-center justify-center text-lg">📦</div>
            </div>
            <div v-if="order.items.length > 4" class="w-12 h-12 rounded-lg bg-gray-100 dark:bg-slate-800 flex items-center justify-center text-sm font-bold text-gray-500">
              +{{ order.items.length - 4 }}
            </div>
          </div>
        </NuxtLink>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center gap-2 mt-8">
        <button 
          v-for="p in totalPages" 
          :key="p"
          @click="goToPage(p)"
          :class="['w-10 h-10 rounded-lg font-medium transition-colors', p === page ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300']"
        >
          {{ p }}
        </button>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useShopApi, type Order } from '~/composables/useShopApi'

definePageMeta({ middleware: 'auth' })

const shopApi = useShopApi()
const loading = ref(true)
const orders = ref<Order[]>([])
const page = ref(1)
const totalPages = ref(1)

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    confirmed: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    processing: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    shipped: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-400',
    delivered: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    cancelled: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
  }
  return `text-xs px-2 py-1 rounded-full ${classes[status] || classes.pending}`
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '⏳ En attente',
    confirmed: '✓ Confirmée',
    processing: '📦 En préparation',
    shipped: '🚚 Expédiée',
    delivered: '✅ Livrée',
    cancelled: '❌ Annulée',
  }
  return labels[status] || status
}

const getPaymentClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
    completed: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    failed: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    refunded: 'bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400',
  }
  return `text-xs px-2 py-1 rounded-full ${classes[status] || classes.pending}`
}

const getPaymentLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '💳 Paiement en attente',
    completed: '💰 Payée',
    failed: '❌ Échec paiement',
    refunded: '↩️ Remboursée',
  }
  return labels[status] || status
}

const loadOrders = async () => {
  loading.value = true
  try {
    const result = await shopApi.listMyOrders(page.value, 20)
    orders.value = result.orders || []
    totalPages.value = result.total_pages || 1
  } catch (e) {
    console.error('Failed to load orders', e)
  } finally {
    loading.value = false
  }
}

const goToPage = (p: number) => {
  page.value = p
  loadOrders()
}

onMounted(loadOrders)
</script>
