<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2 flex items-center gap-3">
          <span class="text-3xl">❤️</span> Mes Favoris
        </h1>
        <p class="text-gray-500">{{ favorites.length }} produits sauvegardés</p>
      </div>
      <button 
        v-if="favorites.length > 0"
        @click="clearAll"
        class="px-4 py-2 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-xl text-sm font-medium transition-colors"
      >
        Tout supprimer
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div v-for="i in 4" :key="i" class="animate-pulse">
        <div class="bg-gray-200 dark:bg-slate-800 aspect-square rounded-2xl mb-4"></div>
        <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded"></div>
      </div>
    </div>

    <!-- Favorites Grid -->
    <div v-else-if="favorites.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <ProductCard 
        v-for="product in favorites" 
        :key="product.id" 
        :product="product" 
        :shop-slug="shopSlug"
        @toggle-favorite="removeFavorite"
      />
    </div>

    <!-- Empty -->
    <div v-else class="text-center py-24 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-gray-800">
      <div class="text-6xl mb-4">💔</div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun favori</h3>
      <p class="text-gray-500 mb-6">Vous n'avez pas encore ajouté de produits à vos favoris.</p>
      <NuxtLink 
        :to="`/shops/${shopSlug}`"
        class="px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors inline-block"
      >
        Découvrir les produits
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useShopApi, type Product } from '~/composables/useShopApi'

definePageMeta({
  layout: 'shop-customer'
})

const route = useRoute()
const shopApi = useShopApi()

const shopSlug = computed(() => route.params.slug as string)
const loading = ref(true)
const favorites = ref<Product[]>([])
const favoriteIds = ref<string[]>([])

const loadFavorites = async () => {
  loading.value = true
  try {
    // Get favorite IDs from localStorage
    const key = `shop_favorites_${shopSlug.value}`
    favoriteIds.value = JSON.parse(localStorage.getItem(key) || '[]')
    
    if (favoriteIds.value.length === 0) {
      favorites.value = []
      return
    }
    
    // Load all products and filter by favorites
    const result = await shopApi.listProducts(shopSlug.value, 1, 100)
    favorites.value = (result.products || []).filter(p => favoriteIds.value.includes(p.id))
    
  } catch (e) {
    console.error('Failed to load favorites', e)
  } finally {
    loading.value = false
  }
}

const removeFavorite = (productId: string) => {
  const key = `shop_favorites_${shopSlug.value}`
  favoriteIds.value = favoriteIds.value.filter(id => id !== productId)
  localStorage.setItem(key, JSON.stringify(favoriteIds.value))
  favorites.value = favorites.value.filter(p => p.id !== productId)
}

const clearAll = () => {
  if (confirm('Supprimer tous les favoris ?')) {
    const key = `shop_favorites_${shopSlug.value}`
    localStorage.setItem(key, '[]')
    favoriteIds.value = []
    favorites.value = []
  }
}

onMounted(loadFavorites)
</script>
