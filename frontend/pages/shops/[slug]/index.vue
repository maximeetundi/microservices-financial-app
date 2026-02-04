<template>
  <div class="space-y-10">
    <!-- Hero Banner with Advanced Search -->
    <div class="relative bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 rounded-3xl p-8 lg:p-12 overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_50%_50%,rgba(255,255,255,0.3)_1px,transparent_1px)] bg-[length:24px_24px]"></div>
      
      <div class="relative z-10 max-w-3xl mx-auto text-center">
        <h1 class="text-3xl lg:text-5xl font-bold text-white mb-4">
          Bienvenue chez {{ shop?.name || 'la boutique' }}
        </h1>
        <p class="text-white/80 text-lg mb-8">
          {{ shop?.description || 'Découvrez nos produits de qualité' }}
        </p>
        
        <!-- Advanced Search -->
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-4 shadow-2xl">
          <div class="flex flex-col lg:flex-row gap-3">
            <div class="flex-1 relative">
              <input 
                v-model="searchQuery"
                type="text" 
                placeholder="Que recherchez-vous ?" 
                class="w-full pl-12 pr-4 py-4 bg-gray-50 dark:bg-slate-800 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 dark:text-white text-base"
              >
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-xl">🔍</span>
            </div>
            <select 
              v-model="selectedCategory"
              class="px-4 py-4 bg-gray-50 dark:bg-slate-800 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 dark:text-white min-w-[180px]"
            >
              <option value="">Toutes catégories</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.slug">{{ cat.name }}</option>
            </select>
            <div class="flex gap-2">
              <input 
                v-model.number="priceMin"
                type="number" 
                placeholder="Prix min" 
                class="w-28 px-4 py-4 bg-gray-50 dark:bg-slate-800 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 dark:text-white"
              >
              <input 
                v-model.number="priceMax"
                type="number" 
                placeholder="Prix max" 
                class="w-28 px-4 py-4 bg-gray-50 dark:bg-slate-800 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 dark:text-white"
              >
            </div>
            <button 
              @click="applyFilters"
              class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-colors shadow-lg shadow-indigo-500/30"
            >
              Rechercher
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Categories Grid -->
    <section v-if="categories.length > 0">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
          <span class="w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-xl flex items-center justify-center">📂</span>
          Catégories
        </h2>
        <NuxtLink :to="`/shops/${shopSlug}/categories`" class="text-indigo-600 dark:text-indigo-400 font-medium hover:underline">
          Voir tout →
        </NuxtLink>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        <NuxtLink 
          v-for="cat in categories.slice(0, 6)" 
          :key="cat.id"
          :to="`/shops/${shopSlug}?category=${cat.slug}`"
          class="group bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:border-indigo-500/30 hover:-translate-y-1 transition-all text-center"
        >
          <div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-indigo-100 to-purple-100 dark:from-slate-800 dark:to-slate-700 rounded-2xl flex items-center justify-center text-3xl group-hover:scale-110 transition-transform">
            {{ cat.icon || '📦' }}
          </div>
          <h3 class="font-bold text-gray-900 dark:text-white mb-1 truncate">{{ cat.name }}</h3>
          <p class="text-sm text-gray-500">{{ cat.product_count || 0 }} produits</p>
        </NuxtLink>
      </div>
    </section>

    <!-- Featured Products -->
    <section v-if="featuredProducts.length > 0">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
          <span class="w-10 h-10 bg-amber-100 dark:bg-amber-900/30 rounded-xl flex items-center justify-center">⭐</span>
          Produits Vedettes
        </h2>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <ProductCard 
          v-for="product in featuredProducts" 
          :key="product.id" 
          :product="product" 
          :shop-slug="shopSlug"
          @toggle-favorite="toggleFavorite"
        />
      </div>
    </section>

    <!-- All Products -->
    <section>
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
          <span class="w-10 h-10 bg-emerald-100 dark:bg-emerald-900/30 rounded-xl flex items-center justify-center">🛍️</span>
          {{ activeFilters ? 'Résultats' : 'Tous les produits' }}
          <span v-if="products.length > 0" class="text-base font-normal text-gray-500">({{ products.length }} produits)</span>
        </h2>
        
        <!-- Sort -->
        <select 
          v-model="sortBy"
          @change="loadProducts"
          class="px-4 py-2 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500"
        >
          <option value="newest">Plus récents</option>
          <option value="price_asc">Prix croissant</option>
          <option value="price_desc">Prix décroissant</option>
          <option value="popular">Populaires</option>
        </select>
      </div>

      <!-- Active Filters -->
      <div v-if="activeFilters" class="flex flex-wrap gap-2 mb-6">
        <span 
          v-if="route.query.search"
          class="px-3 py-1.5 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-full text-sm flex items-center gap-2"
        >
          Recherche: {{ route.query.search }}
          <button @click="clearFilter('search')" class="hover:text-red-500">×</button>
        </span>
        <span 
          v-if="route.query.category"
          class="px-3 py-1.5 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 rounded-full text-sm flex items-center gap-2"
        >
          Catégorie: {{ route.query.category }}
          <button @click="clearFilter('category')" class="hover:text-red-500">×</button>
        </span>
        <span 
          v-if="route.query.tag"
          class="px-3 py-1.5 bg-pink-100 dark:bg-pink-900/30 text-pink-600 dark:text-pink-400 rounded-full text-sm flex items-center gap-2"
        >
          Tag: #{{ route.query.tag }}
          <button @click="clearFilter('tag')" class="hover:text-red-500">×</button>
        </span>
        <button 
          @click="clearAllFilters"
          class="px-3 py-1.5 text-gray-500 hover:text-red-500 text-sm underline"
        >
          Effacer tout
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div v-for="i in 8" :key="i" class="animate-pulse">
          <div class="bg-gray-200 dark:bg-slate-800 aspect-square rounded-2xl mb-4"></div>
          <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded mb-2"></div>
          <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded w-1/2"></div>
        </div>
      </div>

      <!-- Products Grid -->
      <div v-else-if="products.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <ProductCard 
          v-for="product in products" 
          :key="product.id" 
          :product="product" 
          :shop-slug="shopSlug"
          @toggle-favorite="toggleFavorite"
        />
      </div>

      <!-- Empty -->
      <div v-else class="text-center py-16 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-gray-800">
        <div class="text-6xl mb-4">🔍</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun produit trouvé</h3>
        <p class="text-gray-500 mb-6">Essayez d'autres critères de recherche</p>
        <button @click="clearAllFilters" class="px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors">
          Voir tous les produits
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import { useShopApi, type Shop, type Product, type Category } from '~/composables/useShopApi'

definePageMeta({
  layout: 'shop-customer'
})

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()

const shopSlug = computed(() => route.params.slug as string)
const shop = inject<Ref<Shop | null>>('shop', ref(null))
const categories = inject<Ref<Category[]>>('categories', ref([]))

const loading = ref(true)
const products = ref<Product[]>([])
const featuredProducts = ref<Product[]>([])

// Search & Filter State
const searchQuery = ref('')
const selectedCategory = ref('')
const priceMin = ref<number | null>(null)
const priceMax = ref<number | null>(null)
const sortBy = ref('newest')

const activeFilters = computed(() => {
  return route.query.search || route.query.category || route.query.tag || route.query.price_min || route.query.price_max
})

const applyFilters = () => {
  const query: Record<string, any> = {}
  if (searchQuery.value) query.search = searchQuery.value
  if (selectedCategory.value) query.category = selectedCategory.value
  if (priceMin.value) query.price_min = priceMin.value
  if (priceMax.value) query.price_max = priceMax.value
  
  router.push({ path: route.path, query })
}

const clearFilter = (key: string) => {
  const query = { ...route.query }
  delete query[key]
  router.push({ path: route.path, query })
}

const clearAllFilters = () => {
  router.push({ path: route.path })
}

const toggleFavorite = (productId: string) => {
  const key = `shop_favorites_${shopSlug.value}`
  const favs = JSON.parse(localStorage.getItem(key) || '[]')
  const index = favs.indexOf(productId)
  if (index > -1) {
    favs.splice(index, 1)
  } else {
    favs.push(productId)
  }
  localStorage.setItem(key, JSON.stringify(favs))
}

const loadProducts = async () => {
  loading.value = true
  try {
    const options: Record<string, any> = {}
    if (route.query.search) options.search = route.query.search
    if (route.query.category) options.category = route.query.category
    if (route.query.tag) options.tag = route.query.tag
    if (route.query.price_min) options.price_min = route.query.price_min
    if (route.query.price_max) options.price_max = route.query.price_max
    options.sort = sortBy.value
    
    const result = await shopApi.listProducts(shopSlug.value, 1, 100, options)
    products.value = result.products || []
    
    // Featured = products with is_featured flag
    featuredProducts.value = products.value.filter(p => p.is_featured).slice(0, 4)
    
  } catch (e) {
    console.error('Failed to load products', e)
  } finally {
    loading.value = false
  }
}

// Initialize from URL
onMounted(() => {
  searchQuery.value = (route.query.search as string) || ''
  selectedCategory.value = (route.query.category as string) || ''
  priceMin.value = route.query.price_min ? Number(route.query.price_min) : null
  priceMax.value = route.query.price_max ? Number(route.query.price_max) : null
})

// Watch for URL changes
watch(() => route.query, loadProducts, { immediate: true })
</script>
