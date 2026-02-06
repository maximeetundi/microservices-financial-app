<template>
  <div class="min-h-screen bg-gray-50 dark:bg-surface">
    <!-- Search Header -->
    <div class="bg-white dark:bg-surface border-b border-gray-200 dark:border-secondary-800 sticky top-0 z-40">
      <div class="max-w-4xl mx-auto px-4 py-4">
        <div class="flex items-center gap-4">
          <!-- Back Button -->
          <button
            @click="goBack"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-secondary-800 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
          </button>
          
          <!-- Search Input -->
          <div class="flex-1 relative">
            <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input
              v-model="searchQuery"
              @input="performSearch"
              type="text"
              placeholder="Rechercher des produits, boutiques..."
              class="w-full pl-10 pr-4 py-2 bg-gray-100 dark:bg-secondary-800 border border-gray-200 dark:border-secondary-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-transparent"
              autofocus
            />
            <button
              v-if="searchQuery"
              @click="clearSearch"
              class="absolute right-3 top-1/2 transform -translate-y-1/2 p-1 rounded-full hover:bg-gray-200 dark:hover:bg-secondary-700 transition-colors"
            >
              <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
        
        <!-- Search Filters -->
        <div class="flex gap-2 mt-3 overflow-x-auto pb-2">
          <button
            v-for="filter in searchFilters"
            :key="filter.key"
            @click="toggleFilter(filter.key)"
            :class="[
              'px-3 py-1 rounded-full text-sm font-medium whitespace-nowrap transition-colors',
              activeFilters.includes(filter.key)
                ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400'
                : 'bg-gray-100 dark:bg-secondary-800 text-gray-600 dark:text-gray-400'
            ]"
          >
            {{ filter.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Search Results -->
    <div class="max-w-4xl mx-auto px-4 py-6">
      <!-- Loading State -->
      <div v-if="searching" class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
        <span class="ml-3 text-gray-600 dark:text-gray-400">Recherche en cours...</span>
      </div>

      <!-- Results -->
      <div v-else-if="searchQuery && searchResults.length > 0">
        <!-- Results Count -->
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
          {{ searchResults.length }} résultat{{ searchResults.length > 1 ? 's' : '' }} pour "{{ searchQuery }}"
        </p>

        <!-- Products Grid -->
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-8">
          <div
            v-for="item in searchResults.filter(r => r.type === 'product')"
            :key="item.id"
            class="bg-white dark:bg-surface rounded-xl overflow-hidden border border-gray-200 dark:border-secondary-800 hover:shadow-lg transition-shadow"
          >
            <img
              :src="item.image || '/placeholder-product.png'"
              :alt="item.name"
              class="w-full h-40 object-cover bg-gray-200 dark:bg-secondary-700"
            />
            <div class="p-3">
              <h3 class="font-semibold text-gray-900 dark:text-white text-sm truncate mb-1">
                {{ item.name }}
              </h3>
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                {{ item.shop_name }}
              </p>
              <p class="text-lg font-bold text-orange-600">
                {{ formatMoney(item.price) }}
              </p>
              <button
                @click="addToCart(item)"
                class="w-full mt-2 px-3 py-2 bg-orange-600 text-white rounded-lg text-sm font-medium hover:bg-orange-700 transition-colors"
              >
                Ajouter au panier
              </button>
            </div>
          </div>
        </div>

        <!-- Shops List -->
        <div v-if="searchResults.filter(r => r.type === 'shop').length > 0" class="mb-8">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Boutiques</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="shop in searchResults.filter(r => r.type === 'shop')"
              :key="shop.id"
              class="bg-white dark:bg-surface rounded-xl p-4 border border-gray-200 dark:border-secondary-800 hover:shadow-lg transition-shadow"
            >
              <div class="flex gap-4">
                <img
                  :src="shop.logo || '/placeholder-shop.png'"
                  :alt="shop.name"
                  class="w-16 h-16 rounded-lg object-cover bg-gray-200 dark:bg-secondary-700"
                />
                <div class="flex-1">
                  <h3 class="font-semibold text-gray-900 dark:text-white mb-1">
                    {{ shop.name }}
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400 mb-2 line-clamp-2">
                    {{ shop.description }}
                  </p>
                  <NuxtLink
                    :to="`/shops/${shop.slug}`"
                    class="inline-flex items-center px-3 py-1 bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400 rounded-lg text-sm font-medium hover:bg-orange-200 dark:hover:bg-orange-900/50 transition-colors"
                  >
                    Visiter la boutique
                  </NuxtLink>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No Results -->
      <div v-else-if="searchQuery && !searching" class="text-center py-16">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Aucun résultat trouvé
        </h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          Essayez d'autres mots-clés ou explorez nos boutiques
        </p>
        <NuxtLink
          to="/shops"
          class="inline-flex items-center px-6 py-3 bg-orange-600 text-white rounded-xl font-medium hover:bg-orange-700 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/>
          </svg>
          Parcourir les boutiques
        </NuxtLink>
      </div>

      <!-- Initial State -->
      <div v-else class="text-center py-16">
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Que recherchez-vous ?
        </h3>
        <p class="text-gray-600 dark:text-gray-400">
          Tapez des mots-clés pour trouver des produits ou des boutiques
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '~/stores/cart'

const router = useRouter()
const cartStore = useCartStore()

// State
const searchQuery = ref('')
const searchResults = ref([])
const searching = ref(false)
const activeFilters = ref([])

// Search filters
const searchFilters = [
  { key: 'products', label: 'Produits' },
  { key: 'shops', label: 'Boutiques' },
  { key: 'featured', label: 'Vedettes' }
]

// Page meta
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

// Methods
const formatMoney = (amount) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'XOF'
  }).format(amount)
}

const performSearch = async () => {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }

  searching.value = true
  
  try {
    // Simuler une recherche (remplacer par un vrai appel API)
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Mock results - à remplacer par un vrai appel API
    const mockProducts = [
      {
        id: '1',
        name: 'iPhone 14 Pro',
        price: 850000,
        image: '/placeholder-product.png',
        shop_name: 'TechStore',
        shop_id: 'shop1',
        type: 'product'
      },
      {
        id: '2',
        name: 'MacBook Air M2',
        price: 650000,
        image: '/placeholder-product.png',
        shop_name: 'TechStore',
        shop_id: 'shop1',
        type: 'product'
      }
    ]
    
    const mockShops = [
      {
        id: 'shop1',
        name: 'TechStore',
        slug: 'techstore',
        description: 'Votre boutique de produits technologiques',
        logo: '/placeholder-shop.png',
        type: 'shop'
      }
    ]
    
    let results = []
    
    if (activeFilters.value.includes('products') || activeFilters.value.length === 0) {
      results = [...results, ...mockProducts]
    }
    
    if (activeFilters.value.includes('shops') || activeFilters.value.length === 0) {
      results = [...results, ...mockShops]
    }
    
    searchResults.value = results
  } catch (error) {
    console.error('Search error:', error)
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

const clearSearch = () => {
  searchQuery.value = ''
  searchResults.value = []
}

const toggleFilter = (filterKey) => {
  const index = activeFilters.value.indexOf(filterKey)
  if (index > -1) {
    activeFilters.value.splice(index, 1)
  } else {
    activeFilters.value.push(filterKey)
  }
  performSearch()
}

const addToCart = (product) => {
  cartStore.addToCart(
    product,
    1,
    product.shop_id,
    product.shop_name
  )
}

const goBack = () => {
  router.back()
}

// Initialize
onMounted(() => {
  cartStore.loadCart()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
