<template>
  <div class="max-w-7xl mx-auto py-8 px-4">
    <!-- Header -->
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">🛍️ Marketplace</h1>
        <p class="text-gray-500 dark:text-gray-400">Découvrez toutes les boutiques de la communauté</p>
      </div>
      <div class="flex gap-3">
        <NuxtLink to="/shops/my-shops" class="px-4 py-2 bg-white dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 border border-gray-200 dark:border-gray-700 rounded-xl font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors flex items-center gap-2">
          <span>🏪</span> Mes Boutiques
        </NuxtLink>
        <NuxtLink to="/shops/create" class="btn-premium flex items-center gap-2">
          <span>+</span> Créer ma boutique
        </NuxtLink>
      </div>
    </div>

    <!-- Search -->
    <div class="relative mb-8">
      <input 
        v-model="searchQuery" 
        type="text" 
        placeholder="Rechercher une boutique..." 
        class="w-full pl-12 pr-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 transition-shadow"
        @input="debouncedSearch"
      >
      <span class="absolute left-4 top-3.5 text-gray-400 text-xl">🔍</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="i in 6" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 h-72 rounded-2xl"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="shops.length === 0" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
      <div class="text-6xl mb-4">🏪</div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white">Aucune boutique trouvée</h3>
      <p class="text-gray-500 mt-2">Soyez le premier à créer votre boutique !</p>
      <NuxtLink to="/shops/create" class="mt-6 inline-block px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors">
        Créer ma boutique
      </NuxtLink>
    </div>

    <!-- Shops Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <NuxtLink 
        v-for="shop in shops" 
        :key="shop.id" 
        :to="`/shops/${shop.slug}`"
        class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:border-indigo-500/30 transition-all duration-300"
      >
        <!-- Banner -->
        <div class="h-32 bg-gradient-to-br from-indigo-500 to-purple-600 relative overflow-hidden">
          <img v-if="shop.banner_url" :src="shop.banner_url" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" alt="">
          <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
        </div>
        
        <!-- Content -->
        <div class="p-5 relative">
          <!-- Logo -->
          <div class="absolute -top-8 left-5 w-16 h-16 rounded-xl bg-white dark:bg-slate-800 border-4 border-white dark:border-slate-900 shadow-lg overflow-hidden">
            <img v-if="shop.logo_url" :src="shop.logo_url" class="w-full h-full object-cover" alt="">
            <div v-else class="w-full h-full flex items-center justify-center text-2xl bg-gradient-to-br from-indigo-500 to-purple-500 text-white">
              {{ shop.name.charAt(0) }}
            </div>
          </div>
          
          <div class="pt-10">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
              {{ shop.name }}
            </h3>
            <p class="text-gray-500 dark:text-gray-400 text-sm mb-4 line-clamp-2">
              {{ shop.description || 'Découvrez nos produits' }}
            </p>
            
            <!-- Stats -->
            <div class="flex items-center gap-4 text-sm">
              <span class="flex items-center gap-1 text-gray-500">
                <span>📦</span> {{ shop.stats?.total_products || 0 }} produits
              </span>
              <span v-if="shop.stats?.average_rating" class="flex items-center gap-1 text-amber-500">
                <span>⭐</span> {{ shop.stats.average_rating.toFixed(1) }}
              </span>
            </div>
            
            <!-- Tags -->
            <div v-if="shop.tags?.length" class="flex flex-wrap gap-2 mt-3">
              <span v-for="tag in shop.tags.slice(0, 3)" :key="tag" class="px-2 py-1 bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 text-xs rounded-lg">
                {{ tag }}
              </span>
            </div>
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
        :class="[
          'w-10 h-10 rounded-lg font-medium transition-colors',
          p === page 
            ? 'bg-indigo-600 text-white' 
            : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700'
        ]"
      >
        {{ p }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useShopApi, type Shop } from '~/composables/useShopApi'

definePageMeta({
  layout: 'dashboard'
})

const shopApi = useShopApi()

const loading = ref(true)
const shops = ref<Shop[]>([])
const searchQuery = ref('')
const page = ref(1)
const totalPages = ref(1)
const pageSize = 12

let searchTimeout: NodeJS.Timeout | null = null

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    loadShops()
  }, 300)
}

const loadShops = async () => {
  loading.value = true
  try {
    const result = await shopApi.listShops(page.value, pageSize, searchQuery.value)
    shops.value = result.shops || []
    totalPages.value = result.total_pages || 1
  } catch (e) {
    console.error('Failed to load shops', e)
    shops.value = []
  } finally {
    loading.value = false
  }
}

const goToPage = (p: number) => {
  page.value = p
  loadShops()
}

onMounted(() => {
  loadShops()
})
</script>
