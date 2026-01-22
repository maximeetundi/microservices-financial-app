<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto py-8 px-4">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">🏪 Mes Boutiques</h1>
          <p class="text-gray-500">Gérez vos boutiques et vos ventes</p>
        </div>
        <NuxtLink to="/shops/create" class="btn-premium flex items-center gap-2">
          <span>+</span> Nouvelle boutique
        </NuxtLink>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div v-for="i in 2" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 h-48 rounded-2xl"></div>
      </div>

      <!-- Empty -->
      <div v-else-if="shops.length === 0" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl">
        <div class="text-6xl mb-4">🏪</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucune boutique</h3>
        <p class="text-gray-500 mb-6">Créez votre première boutique pour commencer à vendre</p>
        <NuxtLink to="/shops/create" class="btn-premium px-6 py-3">
          Créer ma boutique
        </NuxtLink>
      </div>

      <!-- Shops List -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div 
          v-for="shop in shops" 
          :key="shop.id"
          class="bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800"
        >
          <!-- Banner -->
          <div class="h-24 bg-gradient-to-br from-indigo-500 to-purple-600 relative">
            <img v-if="shop.banner_url" :src="shop.banner_url" class="w-full h-full object-cover" alt="">
          </div>
          
          <div class="p-5 relative">
            <!-- Logo -->
            <div class="absolute -top-8 left-5 w-14 h-14 rounded-xl bg-white dark:bg-slate-800 border-4 border-white dark:border-slate-900 shadow-lg overflow-hidden">
              <img v-if="shop.logo_url" :src="shop.logo_url" class="w-full h-full object-cover" alt="">
              <div v-else class="w-full h-full flex items-center justify-center text-xl bg-gradient-to-br from-indigo-500 to-purple-500 text-white">
                {{ shop.name.charAt(0) }}
              </div>
            </div>
            
            <div class="pt-8">
              <div class="flex items-start justify-between mb-3">
                <div>
                  <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ shop.name }}</h3>
                  <span :class="['text-xs px-2 py-1 rounded-full', shop.status === 'active' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-600']">
                    {{ shop.status === 'active' ? '🟢 Active' : '⏸️ Inactive' }}
                  </span>
                </div>
                <span :class="['text-xs px-2 py-1 rounded-full', shop.is_public ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400' : 'bg-gray-100 text-gray-600']">
                  {{ shop.is_public ? '🌍 Publique' : '🔒 Privée' }}
                </span>
              </div>
              
              <!-- Stats -->
              <div class="grid grid-cols-3 gap-4 py-4 border-t border-gray-100 dark:border-gray-800">
                <div class="text-center">
                  <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ shop.stats?.total_products || 0 }}</div>
                  <div class="text-xs text-gray-500">Produits</div>
                </div>
                <div class="text-center">
                  <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ shop.stats?.total_orders || 0 }}</div>
                  <div class="text-xs text-gray-500">Commandes</div>
                </div>
                <div class="text-center">
                  <div class="text-lg font-bold text-green-600">{{ formatPrice(shop.stats?.total_revenue || 0, shop.currency) }}</div>
                  <div class="text-xs text-gray-500">Revenus</div>
                </div>
              </div>
              
              <!-- Actions -->
              <div class="flex gap-2 mt-4">
                <NuxtLink :to="`/shops/${shop.slug}`" class="flex-1 py-2 text-center bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors">
                  👁️ Voir
                </NuxtLink>
                <NuxtLink :to="`/shops/manage/${shop.id}`" class="flex-1 py-2 text-center bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors">
                  ⚙️ Gérer
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useShopApi, type Shop } from '~/composables/useShopApi'

definePageMeta({ middleware: 'auth' })

const shopApi = useShopApi()
const loading = ref(true)
const shops = ref<Shop[]>([])

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

onMounted(async () => {
  try {
    const result = await shopApi.getMyShops()
    shops.value = result.shops || []
  } catch (e) {
    console.error('Failed to load shops', e)
  } finally {
    loading.value = false
  }
})
</script>
