<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto py-8 px-4">
      <!-- Back -->
      <NuxtLink to="/orders" class="inline-flex items-center text-gray-500 hover:text-indigo-600 mb-6 transition-colors">
        <span class="mr-2">←</span> Retour aux commandes
      </NuxtLink>

      <!-- Loading -->
      <div v-if="loading" class="space-y-6">
        <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-40 rounded-2xl"></div>
        <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-64 rounded-2xl"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl">
        <div class="text-6xl mb-4">😞</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white">Commande non trouvée</h3>
        <NuxtLink to="/orders" class="mt-4 inline-block px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold">
          Retour aux commandes
        </NuxtLink>
      </div>

      <div v-else-if="order">
        <!-- Header Card -->
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800 mb-6">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-1">{{ order.order_number }}</h1>
              <p class="text-gray-500">Passée le {{ formatDate(order.created_at) }}</p>
            </div>
            <div class="flex gap-2">
              <span :class="getStatusClass(order.order_status)">{{ getStatusLabel(order.order_status) }}</span>
              <span :class="getPaymentClass(order.payment_status)">{{ getPaymentLabel(order.payment_status) }}</span>
            </div>
          </div>

          <!-- Status Timeline -->
          <div class="flex items-center justify-between relative mb-2">
            <div class="absolute top-4 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 rounded"></div>
            <div class="absolute top-4 left-0 h-1 bg-indigo-500 rounded transition-all" :style="{ width: getProgressWidth() }"></div>
            
            <div v-for="(step, idx) in statusSteps" :key="step.key" class="relative z-10 flex flex-col items-center">
              <div :class="[
                'w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm transition-all',
                isStepCompleted(step.key) 
                  ? 'bg-indigo-600 text-white' 
                  : isStepCurrent(step.key)
                    ? 'bg-indigo-100 dark:bg-indigo-900 text-indigo-600 dark:text-indigo-400 ring-4 ring-indigo-500/30'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-400'
              ]">
                {{ isStepCompleted(step.key) ? '✓' : idx + 1 }}
              </div>
              <span class="text-xs mt-2 text-gray-500 dark:text-gray-400 text-center">{{ step.label }}</span>
            </div>
          </div>
        </div>

        <!-- Shop Info -->
        <NuxtLink :to="`/shops/${order.shop_slug || order.shop_id}`" class="block bg-white dark:bg-slate-900 rounded-xl p-4 border border-gray-100 dark:border-gray-800 mb-6 hover:border-indigo-500/30 transition-all">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white text-xl font-bold">
              🏪
            </div>
            <div>
              <h3 class="font-bold text-gray-900 dark:text-white">{{ order.shop_name }}</h3>
              <p class="text-sm text-gray-500">Voir la boutique →</p>
            </div>
          </div>
        </NuxtLink>

        <!-- Items -->
        <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-gray-800 mb-6 overflow-hidden">
          <div class="p-4 border-b border-gray-100 dark:border-gray-800">
            <h3 class="font-bold text-gray-900 dark:text-white">📦 Articles ({{ order.items.length }})</h3>
          </div>
          
          <div class="divide-y divide-gray-100 dark:divide-gray-800">
            <div v-for="item in order.items" :key="item.product_id" class="p-4 flex gap-4">
              <div class="w-16 h-16 rounded-lg bg-gray-100 dark:bg-slate-800 overflow-hidden flex-shrink-0">
                <img v-if="item.product_image" :src="item.product_image" class="w-full h-full object-cover" alt="">
                <div v-else class="w-full h-full flex items-center justify-center text-2xl">📦</div>
              </div>
              <div class="flex-1 min-w-0">
                <h4 class="font-semibold text-gray-900 dark:text-white truncate">{{ item.product_name }}</h4>
                <p class="text-sm text-gray-500">Quantité: {{ item.quantity }}</p>
              </div>
              <div class="text-right">
                <p class="font-bold text-gray-900 dark:text-white">{{ formatPrice(item.total_price, order.currency) }}</p>
                <p class="text-xs text-gray-500">{{ formatPrice(item.unit_price, order.currency) }} / unité</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Summary -->
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800 mb-6">
          <h3 class="font-bold text-gray-900 dark:text-white mb-4">💰 Récapitulatif</h3>
          <div class="space-y-3">
            <div class="flex justify-between text-gray-600 dark:text-gray-400">
              <span>Sous-total</span>
              <span>{{ formatPrice(order.sub_total, order.currency) }}</span>
            </div>
            <div v-if="order.delivery_fee" class="flex justify-between text-gray-600 dark:text-gray-400">
              <span>Frais de livraison</span>
              <span>{{ formatPrice(order.delivery_fee, order.currency) }}</span>
            </div>
            <div class="border-t border-gray-100 dark:border-gray-800 pt-3 flex justify-between text-lg font-bold text-gray-900 dark:text-white">
              <span>Total</span>
              <span class="text-indigo-600 dark:text-indigo-400">{{ formatPrice(order.total_amount, order.currency) }}</span>
            </div>
          </div>
        </div>

        <!-- Delivery Info -->
        <div v-if="order.shipping_address || order.delivery_type" class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
          <h3 class="font-bold text-gray-900 dark:text-white mb-4">🚚 Livraison</h3>
          <div class="flex items-center gap-3 mb-4">
            <span class="text-2xl">{{ order.delivery_type === 'pickup' ? '🏃' : '📦' }}</span>
            <span class="font-medium text-gray-900 dark:text-white">
              {{ order.delivery_type === 'pickup' ? 'Retrait en magasin' : 'Livraison à domicile' }}
            </span>
          </div>
          <div v-if="order.shipping_address" class="text-gray-600 dark:text-gray-400">
            <p>{{ order.shipping_address.street }}</p>
            <p>{{ order.shipping_address.city }}, {{ order.shipping_address.state }} {{ order.shipping_address.postal_code }}</p>
            <p>{{ order.shipping_address.country }}</p>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useShopApi, type Order } from '~/composables/useShopApi'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const shopApi = useShopApi()

const orderId = computed(() => route.params.id as string)
const loading = ref(true)
const error = ref(false)
const order = ref<Order | null>(null)

const statusSteps = [
  { key: 'pending', label: 'En attente' },
  { key: 'confirmed', label: 'Confirmée' },
  { key: 'processing', label: 'Préparation' },
  { key: 'shipped', label: 'Expédiée' },
  { key: 'delivered', label: 'Livrée' }
]

const statusOrder = ['pending', 'confirmed', 'processing', 'shipped', 'delivered']

const isStepCompleted = (stepKey: string) => {
  const currentIdx = statusOrder.indexOf(order.value?.order_status || '')
  const stepIdx = statusOrder.indexOf(stepKey)
  return stepIdx < currentIdx
}

const isStepCurrent = (stepKey: string) => {
  return order.value?.order_status === stepKey
}

const getProgressWidth = () => {
  const currentIdx = statusOrder.indexOf(order.value?.order_status || '')
  if (currentIdx <= 0) return '0%'
  return `${(currentIdx / (statusOrder.length - 1)) * 100}%`
}

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
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
  return `text-xs px-3 py-1 rounded-full font-medium ${classes[status] || classes.pending}`
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '⏳ En attente',
    confirmed: '✓ Confirmée',
    processing: '📦 Préparation',
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
  return `text-xs px-3 py-1 rounded-full font-medium ${classes[status] || classes.pending}`
}

const getPaymentLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '💳 En attente',
    completed: '💰 Payée',
    failed: '❌ Échec',
    refunded: '↩️ Remboursée',
  }
  return labels[status] || status
}

const loadOrder = async () => {
  loading.value = true
  error.value = false
  try {
    order.value = await shopApi.getOrder(orderId.value)
  } catch (e) {
    console.error('Failed to load order', e)
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(loadOrder)
</script>
