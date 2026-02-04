<template>
  <div class="space-y-12">
    <!-- Hero Slider -->
    <section class="relative rounded-3xl overflow-hidden bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-500 h-[400px] lg:h-[500px]">
      <div class="absolute inset-0 bg-black/20"></div>
      <div class="absolute inset-0 flex items-center">
        <div class="max-w-2xl mx-auto text-center px-6">
          <span class="inline-block px-4 py-1.5 bg-white/20 backdrop-blur rounded-full text-white text-sm font-medium mb-6">
            🎉 Bienvenue dans notre boutique
          </span>
          <h1 class="text-4xl lg:text-6xl font-bold text-white mb-6 leading-tight">
            {{ shop?.name || 'Notre Boutique' }}
          </h1>
          <p class="text-xl text-white/90 mb-8">
            {{ shop?.description || 'Découvrez nos produits de qualité à des prix imbattables' }}
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button 
              @click="scrollToProducts"
              class="px-8 py-4 bg-white text-indigo-600 font-bold rounded-xl hover:bg-gray-100 transition-colors shadow-xl"
            >
              Voir les produits
            </button>
            <NuxtLink 
              :to="`/shops/${shopSlug}/categories`"
              class="px-8 py-4 bg-white/20 backdrop-blur text-white font-bold rounded-xl hover:bg-white/30 transition-colors"
            >
              Nos catégories
            </NuxtLink>
          </div>
        </div>
      </div>
      <!-- Decorative elements -->
      <div class="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-gray-100 dark:from-slate-950 to-transparent"></div>
    </section>

    <!-- Features Bar -->
    <section class="grid grid-cols-2 lg:grid-cols-4 gap-4 -mt-16 relative z-10">
      <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-xl flex items-center gap-4">
        <div class="w-14 h-14 bg-indigo-100 dark:bg-indigo-900/30 rounded-xl flex items-center justify-center text-2xl">🚚</div>
        <div>
          <h3 class="font-bold text-gray-900 dark:text-white">Livraison rapide</h3>
          <p class="text-sm text-gray-500">Partout au Sénégal</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-xl flex items-center gap-4">
        <div class="w-14 h-14 bg-green-100 dark:bg-green-900/30 rounded-xl flex items-center justify-center text-2xl">🔒</div>
        <div>
          <h3 class="font-bold text-gray-900 dark:text-white">Paiement sécurisé</h3>
          <p class="text-sm text-gray-500">100% sécurisé</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-xl flex items-center gap-4">
        <div class="w-14 h-14 bg-amber-100 dark:bg-amber-900/30 rounded-xl flex items-center justify-center text-2xl">⭐</div>
        <div>
          <h3 class="font-bold text-gray-900 dark:text-white">Qualité garantie</h3>
          <p class="text-sm text-gray-500">Produits vérifiés</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-xl flex items-center gap-4">
        <div class="w-14 h-14 bg-purple-100 dark:bg-purple-900/30 rounded-xl flex items-center justify-center text-2xl">💬</div>
        <div>
          <h3 class="font-bold text-gray-900 dark:text-white">Support 24/7</h3>
          <p class="text-sm text-gray-500">À votre écoute</p>
        </div>
      </div>
    </section>

    <!-- Categories Section -->
    <section v-if="categories.length > 0">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white">Nos Catégories</h2>
          <p class="text-gray-500 mt-1">Explorez notre sélection</p>
        </div>
        <NuxtLink 
          :to="`/shops/${shopSlug}/categories`"
          class="hidden sm:flex items-center gap-2 px-4 py-2 text-indigo-600 dark:text-indigo-400 font-medium hover:bg-indigo-50 dark:hover:bg-indigo-900/20 rounded-xl transition-colors"
        >
          Voir tout <span>→</span>
        </NuxtLink>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        <NuxtLink 
          v-for="cat in categories.slice(0, 6)" 
          :key="cat.id"
          :to="`/shops/${shopSlug}?category=${cat.slug}`"
          class="group bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800 hover:shadow-2xl hover:border-indigo-500/50 hover:-translate-y-2 transition-all duration-300 text-center"
        >
          <div class="w-20 h-20 mx-auto mb-4 bg-gradient-to-br from-indigo-100 to-purple-100 dark:from-slate-800 dark:to-slate-700 rounded-2xl flex items-center justify-center text-4xl group-hover:scale-110 transition-transform">
            {{ cat.icon || '📦' }}
          </div>
          <h3 class="font-bold text-gray-900 dark:text-white mb-1 group-hover:text-indigo-600 transition-colors">{{ cat.name }}</h3>
          <p class="text-sm text-gray-500">{{ cat.product_count || 0 }} produits</p>
        </NuxtLink>
      </div>
    </section>

    <!-- Featured Products -->
    <section v-if="featuredProducts.length > 0" class="bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-slate-900 dark:to-slate-800 rounded-3xl p-8">
      <div class="flex items-center justify-between mb-8">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 bg-amber-500 rounded-xl flex items-center justify-center text-2xl animate-pulse">⭐</div>
          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white">Produits Vedettes</h2>
            <p class="text-gray-500">Nos meilleures ventes</p>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <ProductCard 
          v-for="product in featuredProducts.slice(0, 4)" 
          :key="product.id" 
          :product="product" 
          :shop-slug="shopSlug"
          @toggle-favorite="toggleFavorite"
        />
      </div>
    </section>

    <!-- All Products -->
    <section id="products-section">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 mb-8">
        <div>
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
            {{ activeFilters ? 'Résultats de recherche' : 'Tous les produits' }}
          </h2>
          <p class="text-gray-500 mt-1">{{ totalProducts }} produits disponibles</p>
        </div>
        
        <div class="flex flex-wrap items-center gap-3">
          <!-- Active Filters -->
          <template v-if="activeFilters">
            <span 
              v-if="route.query.search"
              class="px-3 py-1.5 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-full text-sm flex items-center gap-2"
            >
              "{{ route.query.search }}"
              <button @click="clearFilter('search')" class="hover:text-red-500 font-bold">×</button>
            </span>
            <span 
              v-if="route.query.category"
              class="px-3 py-1.5 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 rounded-full text-sm flex items-center gap-2"
            >
              {{ findCategoryName(route.query.category as string) }}
              <button @click="clearFilter('category')" class="hover:text-red-500 font-bold">×</button>
            </span>
            <button @click="clearAllFilters" class="text-sm text-gray-500 hover:text-red-500 underline">
              Effacer
            </button>
          </template>

          <!-- Sort -->
          <select 
            v-model="sortBy"
            @change="loadProducts"
            class="px-4 py-2.5 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="newest">Plus récents</option>
            <option value="price_asc">Prix croissant</option>
            <option value="price_desc">Prix décroissant</option>
            <option value="popular">Populaires</option>
          </select>

          <!-- View Mode -->
          <div class="hidden lg:flex bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
            <button 
              @click="viewMode = 'grid'"
              :class="viewMode === 'grid' ? 'bg-indigo-600 text-white' : 'text-gray-500'"
              class="p-2.5 transition-colors"
            >⊞</button>
            <button 
              @click="viewMode = 'list'"
              :class="viewMode === 'list' ? 'bg-indigo-600 text-white' : 'text-gray-500'"
              class="p-2.5 transition-colors"
            >☰</button>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div v-for="i in 8" :key="i" class="bg-white dark:bg-slate-900 rounded-2xl overflow-hidden">
          <div class="animate-pulse">
            <div class="bg-gray-200 dark:bg-slate-800 aspect-square"></div>
            <div class="p-4 space-y-3">
              <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Products Grid -->
      <div v-else-if="products.length > 0" :class="viewMode === 'grid' ? 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6' : 'space-y-4'">
        <ProductCard 
          v-for="product in products" 
          :key="product.id" 
          :product="product" 
          :shop-slug="shopSlug"
          :list-mode="viewMode === 'list'"
          @toggle-favorite="toggleFavorite"
        />
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20 bg-white dark:bg-slate-900 rounded-3xl">
        <div class="w-24 h-24 mx-auto mb-6 bg-gray-100 dark:bg-slate-800 rounded-full flex items-center justify-center text-5xl">🔍</div>
        <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Aucun produit trouvé</h3>
        <p class="text-gray-500 mb-8 max-w-md mx-auto">Nous n'avons trouvé aucun produit correspondant à vos critères. Essayez d'autres filtres.</p>
        <button 
          @click="clearAllFilters"
          class="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-colors shadow-lg shadow-indigo-500/30"
        >
          Voir tous les produits
        </button>
      </div>
    </section>

    <!-- Newsletter -->
    <section class="bg-gradient-to-r from-indigo-600 to-purple-600 rounded-3xl p-8 lg:p-12 text-center">
      <h2 class="text-3xl font-bold text-white mb-4">Restez informé</h2>
      <p class="text-white/80 mb-8 max-w-lg mx-auto">Inscrivez-vous à notre newsletter pour recevoir nos offres exclusives et nouveautés.</p>
      <form @submit.prevent class="flex flex-col sm:flex-row gap-4 max-w-md mx-auto">
        <input 
          type="email" 
          placeholder="Votre email" 
          class="flex-1 px-6 py-4 rounded-xl focus:outline-none focus:ring-2 focus:ring-white/50"
        >
        <button type="submit" class="px-8 py-4 bg-white text-indigo-600 font-bold rounded-xl hover:bg-gray-100 transition-colors shadow-lg">
          S'inscrire
        </button>
      </form>
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
const totalProducts = ref(0)
const sortBy = ref('newest')
const viewMode = ref('grid')

const activeFilters = computed(() => {
  return route.query.search || route.query.category || route.query.tag
})

const findCategoryName = (slug: string) => {
  return categories.value.find(c => c.slug === slug)?.name || slug
}

const scrollToProducts = () => {
  document.getElementById('products-section')?.scrollIntoView({ behavior: 'smooth' })
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
    options.sort = sortBy.value
    
    const result = await shopApi.listProducts(shopSlug.value, 1, 100, options)
    products.value = result.products || []
    totalProducts.value = result.total || products.value.length
    
    // Featured = products with is_featured flag
    if (!activeFilters.value) {
      featuredProducts.value = products.value.filter(p => p.is_featured).slice(0, 4)
    }
    
  } catch (e) {
    console.error('Failed to load products', e)
  } finally {
    loading.value = false
  }
}

// Watch for URL and categories changes
watch(() => route.query, loadProducts, { immediate: true })
watch(categories, () => {
  if (categories.value.length > 0 && products.value.length === 0) {
    loadProducts()
  }
})
</script>
