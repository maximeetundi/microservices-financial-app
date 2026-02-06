<template>
  <div class="space-y-10">
    <!-- Hero Banner -->
    <section
      class="relative rounded-2xl overflow-hidden h-[280px] md:h-[360px]"
      :style="shop?.banner_url ? `background-image: url('${shop.banner_url}'); background-size: cover; background-position: center;` : ''"
      :class="!shop?.banner_url ? 'bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500' : ''"
    >
      <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent"></div>
      <div class="absolute inset-0 flex items-end">
        <div class="w-full p-6 md:p-10">
          <div class="max-w-2xl">
            <h1 class="text-2xl md:text-4xl font-bold text-white mb-3 drop-shadow-lg">
              {{ shop?.name || 'Bienvenue' }}
            </h1>
            <p class="text-white/90 text-sm md:text-base mb-5 line-clamp-2 drop-shadow">
              {{ shop?.description || 'Découvrez nos produits de qualité' }}
            </p>
            <div class="flex flex-wrap gap-3">
              <button
                @click="scrollToProducts"
                class="px-5 py-2.5 bg-white text-indigo-600 font-semibold rounded-xl hover:bg-gray-100 transition-all shadow-lg text-sm"
              >
                Voir les produits
              </button>
              <NuxtLink
                :to="`/shops/${shopSlug}/categories`"
                class="px-5 py-2.5 bg-white/20 backdrop-blur text-white font-semibold rounded-xl hover:bg-white/30 transition-all border border-white/30 text-sm"
              >
                Explorer les catégories
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Trust Badges -->
    <section class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <div class="bg-white dark:bg-slate-900 rounded-xl p-4 flex items-center gap-3 border border-gray-100 dark:border-gray-800">
        <div class="w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center text-xl">🚚</div>
        <div>
          <p class="font-semibold text-gray-900 dark:text-white text-sm">Livraison rapide</p>
          <p class="text-xs text-gray-500">Partout au Sénégal</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-xl p-4 flex items-center gap-3 border border-gray-100 dark:border-gray-800">
        <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center text-xl">🔒</div>
        <div>
          <p class="font-semibold text-gray-900 dark:text-white text-sm">Paiement sécurisé</p>
          <p class="text-xs text-gray-500">100% sécurisé</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-xl p-4 flex items-center gap-3 border border-gray-100 dark:border-gray-800">
        <div class="w-10 h-10 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center text-xl">⭐</div>
        <div>
          <p class="font-semibold text-gray-900 dark:text-white text-sm">Qualité garantie</p>
          <p class="text-xs text-gray-500">Produits vérifiés</p>
        </div>
      </div>
      <div class="bg-white dark:bg-slate-900 rounded-xl p-4 flex items-center gap-3 border border-gray-100 dark:border-gray-800">
        <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center text-xl">💬</div>
        <div>
          <p class="font-semibold text-gray-900 dark:text-white text-sm">Support 24/7</p>
          <p class="text-xs text-gray-500">À votre écoute</p>
        </div>
      </div>
    </section>

    <!-- Categories Preview -->
    <section v-if="categories.length > 0">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white">Nos Catégories</h2>
        <NuxtLink
          :to="`/shops/${shopSlug}/categories`"
          class="text-sm text-indigo-600 dark:text-indigo-400 font-medium hover:underline flex items-center gap-1"
        >
          Voir tout <span>→</span>
        </NuxtLink>
      </div>
      <div class="flex gap-3 overflow-x-auto pb-2 scrollbar-hide -mx-4 px-4">
        <NuxtLink
          v-for="cat in categories.slice(0, 8)"
          :key="cat.id"
          :to="`/shops/${shopSlug}?category=${cat.slug}`"
          class="flex-shrink-0 group"
        >
          <div class="w-20 h-20 md:w-24 md:h-24 bg-gradient-to-br from-indigo-100 to-purple-100 dark:from-slate-800 dark:to-slate-700 rounded-2xl flex items-center justify-center text-3xl md:text-4xl mb-2 group-hover:scale-105 transition-transform border-2 border-transparent group-hover:border-indigo-500">
            {{ cat.icon || '📦' }}
          </div>
          <p class="text-xs font-medium text-gray-700 dark:text-gray-300 text-center truncate w-20 md:w-24">{{ cat.name }}</p>
        </NuxtLink>
      </div>
    </section>

    <!-- Featured Products -->
    <section v-if="featuredProducts.length > 0" class="bg-gradient-to-r from-amber-50 to-orange-50 dark:from-slate-900 dark:to-slate-800 rounded-2xl p-5 md:p-6">
      <div class="flex items-center gap-3 mb-5">
        <div class="w-10 h-10 bg-amber-500 rounded-xl flex items-center justify-center text-xl">⭐</div>
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Produits Vedettes</h2>
          <p class="text-sm text-gray-500">Nos meilleures ventes</p>
        </div>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <ProductCard
          v-for="product in featuredProducts.slice(0, 4)"
          :key="product.id"
          :product="product"
          :shop-slug="shopSlug"
          @add-to-cart="addToCart"
          @toggle-favorite="toggleFavorite"
        />
      </div>
    </section>

    <!-- All Products Section -->
    <section id="products-section">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-5">
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">
            {{ hasFilters ? 'Résultats' : 'Tous les produits' }}
          </h2>
          <p class="text-sm text-gray-500">{{ totalProducts }} produit{{ totalProducts > 1 ? 's' : '' }}</p>
        </div>

        <div class="flex items-center gap-3">
          <!-- Active Filters -->
          <div v-if="hasFilters" class="flex items-center gap-2 flex-wrap">
            <span
              v-if="route.query.search"
              class="px-3 py-1 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-full text-xs font-medium flex items-center gap-1"
            >
              "{{ route.query.search }}"
              <button @click="clearFilter('search')" class="hover:text-red-500 ml-1">×</button>
            </span>
            <span
              v-if="route.query.category"
              class="px-3 py-1 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 rounded-full text-xs font-medium flex items-center gap-1"
            >
              {{ getCategoryName(route.query.category as string) }}
              <button @click="clearFilter('category')" class="hover:text-red-500 ml-1">×</button>
            </span>
            <button @click="clearAllFilters" class="text-xs text-gray-500 hover:text-red-500 underline">
              Effacer tout
            </button>
          </div>

          <!-- Sort -->
          <select
            v-model="sortBy"
            @change="loadProducts"
            class="px-3 py-2 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500"
          >
            <option value="newest">Plus récents</option>
            <option value="price_asc">Prix ↑</option>
            <option value="price_desc">Prix ↓</option>
            <option value="popular">Populaires</option>
          </select>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div v-for="i in 8" :key="i" class="bg-white dark:bg-slate-900 rounded-xl overflow-hidden border border-gray-100 dark:border-gray-800">
          <div class="animate-pulse">
            <div class="bg-gray-200 dark:bg-slate-800 aspect-square"></div>
            <div class="p-3 space-y-2">
              <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Products Grid -->
      <div v-else-if="products.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <ProductCard
          v-for="product in products"
          :key="product.id"
          :product="product"
          :shop-slug="shopSlug"
          @add-to-cart="addToCart"
          @toggle-favorite="toggleFavorite"
        />
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white dark:bg-slate-900 rounded-2xl border border-dashed border-gray-200 dark:border-gray-800">
        <div class="w-20 h-20 mx-auto mb-4 bg-gray-100 dark:bg-slate-800 rounded-full flex items-center justify-center text-4xl">🔍</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun produit trouvé</h3>
        <p class="text-gray-500 mb-6 text-sm max-w-sm mx-auto">Nous n'avons trouvé aucun produit correspondant à vos critères.</p>
        <button
          @click="clearAllFilters"
          class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl transition-colors text-sm"
        >
          Voir tous les produits
        </button>
      </div>

      <!-- Load More -->
      <div v-if="products.length > 0 && products.length < totalProducts" class="text-center mt-8">
        <button
          @click="loadMore"
          :disabled="loadingMore"
          class="px-6 py-3 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 rounded-xl font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors disabled:opacity-50"
        >
          {{ loadingMore ? 'Chargement...' : 'Charger plus de produits' }}
        </button>
      </div>
    </section>

    <!-- Newsletter -->
    <section class="bg-gradient-to-r from-indigo-600 to-purple-600 rounded-2xl p-6 md:p-8 text-center">
      <h2 class="text-2xl font-bold text-white mb-2">Restez informé</h2>
      <p class="text-white/80 mb-6 text-sm max-w-md mx-auto">Inscrivez-vous pour recevoir nos offres exclusives et nouveautés.</p>
      <form @submit.prevent="subscribeNewsletter" class="flex flex-col sm:flex-row gap-3 max-w-md mx-auto">
        <input
          v-model="newsletterEmail"
          type="email"
          placeholder="Votre email"
          class="flex-1 px-4 py-3 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-white/50"
          required
        >
        <button type="submit" class="px-6 py-3 bg-white text-indigo-600 font-semibold rounded-xl hover:bg-gray-100 transition-colors text-sm">
          S'inscrire
        </button>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
import { inject, ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Product, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'

definePageMeta({
  layout: 'shop-layout'
})

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()

const shopSlug = computed(() => route.params.slug as string)
const shop = inject<Ref<Shop | null>>('shop', ref(null))
const categories = inject<Ref<Category[]>>('categories', ref([]))

const loading = ref(true)
const loadingMore = ref(false)
const products = ref<Product[]>([])
const featuredProducts = ref<Product[]>([])
const totalProducts = ref(0)
const currentPage = ref(1)
const sortBy = ref('newest')
const newsletterEmail = ref('')

const hasFilters = computed(() => {
  return !!route.query.search || !!route.query.category || !!route.query.tag
})

const getCategoryName = (slug: string) => {
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

const addToCart = (product: Product) => {
  cartStore.addItem({
    id: product.id,
    name: product.name,
    price: product.price,
    image: product.images?.[0] || '',
    quantity: 1,
    shopId: shop.value?.id || 'unknown',
    shopName: shop.value?.name || 'Boutique'
  })
}

const loadProducts = async (reset = true) => {
  if (reset) {
    loading.value = true
    currentPage.value = 1
  }

  try {
    const options: Record<string, any> = { sort: sortBy.value }
    if (route.query.search) options.search = route.query.search
    if (route.query.category) options.category = route.query.category
    if (route.query.tag) options.tag = route.query.tag

    const result = await shopApi.listProducts(shopSlug.value, currentPage.value, 20, options)

    if (reset) {
      products.value = result.products || []
    } else {
      products.value = [...products.value, ...(result.products || [])]
    }

    totalProducts.value = result.total || products.value.length

    // Load featured products on first load without filters
    if (reset && !hasFilters.value) {
      featuredProducts.value = (result.products || []).filter((p: Product) => p.is_featured).slice(0, 4)
    }
  } catch (e) {
    console.error('Failed to load products', e)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = async () => {
  loadingMore.value = true
  currentPage.value++
  await loadProducts(false)
}

const subscribeNewsletter = () => {
  // TODO: Implement newsletter subscription
  alert(`Merci ! Vous serez notifié à ${newsletterEmail.value}`)
  newsletterEmail.value = ''
}

// Watch for URL changes
watch(() => route.query, () => loadProducts(true), { immediate: true })

// Watch for categories to be loaded
watch(categories, () => {
  if (categories.value.length > 0 && products.value.length === 0 && !loading.value) {
    loadProducts(true)
  }
})

// Initialize on mount
onMounted(() => {
  cartStore.loadCart()
  // Load products if not already loading
  if (!loading.value && products.value.length === 0) {
    loadProducts(true)
  }
})
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
