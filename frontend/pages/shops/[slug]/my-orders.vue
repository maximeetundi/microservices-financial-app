<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2 flex items-center gap-3">
          <span class="text-3xl">📦</span> Mes Commandes
        </h1>
        <p class="text-gray-500">Historique de vos achats dans cette boutique</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="animate-pulse bg-gray-200 dark:bg-slate-800 h-24 rounded-2xl"></div>
    </div>

    <!-- Orders List -->
    <div v-else-if="orders.length > 0" class="space-y-4">
      <div 
        v-for="order in orders" 
        :key="order.id"
        class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-gray-800 p-6 hover:shadow-lg transition-shadow"
      >
        <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
          <!-- Order Info -->
          <div class="flex items-start gap-4">
            <div class="w-16 h-16 bg-indigo-100 dark:bg-indigo-900/30 rounded-xl flex items-center justify-center text-2xl flex-shrink-0">
              📦
            </div>
            <div>
              <h3 class="font-bold text-gray-900 dark:text-white mb-1">
                Commande #{{ order.id.slice(-8).toUpperCase() }}
              </h3>
              <p class="text-sm text-gray-500 mb-2">
                {{ new Date(order.created_at).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' }) }}
              </p>
              <div class="flex flex-wrap gap-2">
                <span 
                  :class="statusClass(order.status)"
                  class="px-3 py-1 rounded-full text-xs font-bold"
                >
                  {{ statusLabel(order.status) }}
                </span>
                <span class="px-3 py-1 bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-gray-400 rounded-full text-xs">
                  {{ order.items?.length || 0 }} article(s)
                </span>
              </div>
            </div>
          </div>
          
          <!-- Amount & Actions -->
          <div class="flex items-center gap-6">
            <div class="text-right">
              <p class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">
                {{ formatPrice(order.total_amount, order.currency) }}
              </p>
            </div>
            <NuxtLink 
              :to="`/orders/${order.id}`"
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-medium transition-colors"
            >
              Détails
            </NuxtLink>
          </div>
        </div>
        
        <!-- Items Preview -->
        <div v-if="order.items?.length" class="mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
          <div class="flex gap-3 overflow-x-auto pb-2">
            <div 
              v-for="item in order.items.slice(0, 4)" 
              :key="item.product_id"
              class="w-16 h-16 bg-gray-100 dark:bg-slate-800 rounded-lg flex-shrink-0 overflow-hidden"
            >
              <img 
                v-if="item.image_url" 
                :src="item.image_url" 
                class="w-full h-full object-cover"
              >
              <div v-else class="w-full h-full flex items-center justify-center text-xl">📦</div>
            </div>
            <div 
              v-if="order.items.length > 4"
              class="w-16 h-16 bg-gray-100 dark:bg-slate-800 rounded-lg flex-shrink-0 flex items-center justify-center text-sm font-bold text-gray-500"
            >
              +{{ order.items.length - 4 }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="text-center py-24 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-gray-800">
      <div class="text-6xl mb-4">🛒</div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucune commande</h3>
      <p class="text-gray-500 mb-6">Vous n'avez pas encore passé de commande dans cette boutique.</p>
      <NuxtLink 
        :to="`/shops/${shopSlug}`"
        class="px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors inline-block"
      >
        Commencer mes achats
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useShopApi } from '~/composables/useShopApi'

definePageMeta({
  layout: 'shop-layout'
})

const route = useRoute()
const shopApi = useShopApi()

const shopSlug = computed(() => route.params.slug as string)
const loading = ref(true)
const orders = ref<any[]>([])

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const statusClass = (status: string) => {
  const classes: Record<string, string> = {
    'pending': 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    'confirmed': 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    'processing': 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
    'shipped': 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    'delivered': 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    'cancelled': 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
  }
  return classes[status?.toLowerCase()] || classes['pending']
}

const statusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'pending': 'En attente',
    'confirmed': 'Confirmée',
    'processing': 'En préparation',
    'shipped': 'Expédiée',
    'delivered': 'Livrée',
    'cancelled': 'Annulée',
  }
  return labels[status?.toLowerCase()] || status
}

const loadOrders = async () => {
  loading.value = true
  try {
    // TODO: Filter by shop when API supports it
    const result = await shopApi.listMyOrders()
    // Filter orders by this shop
    orders.value = (result.orders || []).filter((o: any) => {
      // Check if order belongs to this shop
      return true // For now, show all orders - filter can be added later
    })
  } catch (e) {
    console.error('Failed to load orders', e)
  } finally {
    loading.value = false
  }
}

onMounted(loadOrders)
</script>
