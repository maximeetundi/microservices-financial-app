<template>
  <div class="min-h-screen bg-gray-50 dark:bg-slate-950 flex flex-col">
    <!-- Top Bar -->
    <div class="bg-indigo-600 text-white text-xs py-1.5">
      <div class="max-w-7xl mx-auto px-4 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <span>🚚 Livraison gratuite à partir de 25 000 FCFA</span>
          <span class="hidden sm:inline">•</span>
          <span class="hidden sm:inline">⭐ Satisfaction garantie</span>
        </div>
        <NuxtLink to="/shops" class="hover:underline flex items-center gap-1">
          ← Retour Marketplace
        </NuxtLink>
      </div>
    </div>

    <!-- Main Header -->
    <header class="bg-white dark:bg-slate-900 border-b border-gray-200 dark:border-gray-800 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16 gap-4">
          <!-- Logo & Shop Name -->
          <NuxtLink :to="`/shops/${shopSlug}`" class="flex items-center gap-3 flex-shrink-0">
            <div v-if="shop?.logo_url" class="w-10 h-10 rounded-xl overflow-hidden shadow-sm">
              <img :src="shop.logo_url" class="w-full h-full object-cover" alt="">
            </div>
            <div v-else class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg shadow-sm">
              {{ shop?.name?.charAt(0)?.toUpperCase() || 'S' }}
            </div>
            <div class="hidden sm:block">
              <h1 class="font-bold text-gray-900 dark:text-white text-lg leading-tight">{{ shop?.name || 'Boutique' }}</h1>
              <div class="flex items-center gap-1 text-xs text-gray-500">
                <span class="text-amber-500">★</span>
                <span>{{ shop?.stats?.average_rating?.toFixed(1) || '4.8' }}</span>
                <span class="mx-1">•</span>
                <span>{{ shop?.stats?.total_products || 0 }} produits</span>
              </div>
            </div>
          </NuxtLink>

          <!-- Search Bar -->
          <div class="flex-1 max-w-xl hidden md:block">
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Rechercher dans cette boutique..."
                class="w-full pl-10 pr-4 py-2.5 bg-gray-100 dark:bg-slate-800 border-0 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:bg-white dark:focus:bg-slate-700 transition-all"
                @keyup.enter="performSearch"
              >
              <svg class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
              </svg>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2">
            <!-- Search Mobile -->
            <button @click="showMobileSearch = !showMobileSearch" class="md:hidden p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
              </svg>
            </button>

            <!-- Favorites -->
            <NuxtLink
              :to="`/shops/${shopSlug}/favorites`"
              class="relative p-2 text-gray-500 hover:text-red-500 dark:text-gray-400 dark:hover:text-red-400 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
              </svg>
              <span v-if="favoritesCount > 0" class="absolute -top-0.5 -right-0.5 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
                {{ favoritesCount > 9 ? '9+' : favoritesCount }}
              </span>
            </NuxtLink>

            <!-- Cart -->
            <NuxtLink
              :to="`/shops/${shopSlug}/cart`"
              class="relative flex items-center gap-2 px-3 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
              </svg>
              <span class="hidden sm:inline font-medium text-sm">Panier</span>
              <span v-if="cartStore.itemCount > 0" class="min-w-[20px] h-5 bg-white text-indigo-600 text-xs font-bold rounded-full flex items-center justify-center px-1">
                {{ cartStore.itemCount }}
              </span>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Navigation Bar -->
      <nav class="border-t border-gray-100 dark:border-gray-800 bg-white dark:bg-slate-900">
        <div class="max-w-7xl mx-auto px-4">
          <div class="flex items-center gap-1 overflow-x-auto scrollbar-hide py-2 -mx-4 px-4">
            <NuxtLink
              :to="`/shops/${shopSlug}`"
              class="flex-shrink-0 px-4 py-2 text-sm font-medium rounded-lg transition-colors whitespace-nowrap"
              :class="isExactHome ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800'"
            >
              🏠 Accueil
            </NuxtLink>

            <NuxtLink
              :to="`/shops/${shopSlug}/categories`"
              class="flex-shrink-0 px-4 py-2 text-sm font-medium rounded-lg transition-colors whitespace-nowrap"
              :class="isActive('categories') ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800'"
            >
              📂 Catégories
            </NuxtLink>

            <!-- Divider -->
            <div class="w-px h-6 bg-gray-200 dark:bg-gray-700 mx-1 flex-shrink-0"></div>

            <!-- Category Links -->
            <NuxtLink
              v-for="cat in categories.slice(0, 6)"
              :key="cat.id"
              :to="`/shops/${shopSlug}?category=${cat.slug}`"
              class="flex-shrink-0 px-3 py-2 text-sm font-medium rounded-lg transition-colors whitespace-nowrap flex items-center gap-1.5"
              :class="route.query.category === cat.slug ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800'"
            >
              <span>{{ cat.icon || '📁' }}</span>
              <span>{{ cat.name }}</span>
            </NuxtLink>

            <!-- More Categories -->
            <div v-if="categories.length > 6" class="relative flex-shrink-0">
              <button
                @click="showMoreCategories = !showMoreCategories"
                class="px-3 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-lg transition-colors whitespace-nowrap flex items-center gap-1"
              >
                Plus
                <svg class="w-4 h-4" :class="showMoreCategories ? 'rotate-180' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </button>
              <!-- Dropdown -->
              <div
                v-if="showMoreCategories"
                class="absolute top-full left-0 mt-1 w-56 bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 py-2 z-50"
              >
                <NuxtLink
                  v-for="cat in categories.slice(6)"
                  :key="cat.id"
                  :to="`/shops/${shopSlug}?category=${cat.slug}`"
                  class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700"
                  @click="showMoreCategories = false"
                >
                  <span>{{ cat.icon || '📁' }}</span>
                  <span>{{ cat.name }}</span>
                  <span class="ml-auto text-xs text-gray-400">{{ cat.product_count || 0 }}</span>
                </NuxtLink>
              </div>
            </div>

            <!-- Spacer -->
            <div class="flex-1"></div>

            <!-- User Menu -->
            <NuxtLink
              :to="`/shops/${shopSlug}/my-orders`"
              class="flex-shrink-0 px-3 py-2 text-sm font-medium rounded-lg transition-colors whitespace-nowrap"
              :class="isActive('my-orders') ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-slate-800'"
            >
              📦 Mes Commandes
            </NuxtLink>
          </div>
        </div>
      </nav>

      <!-- Mobile Search (Expandable) -->
      <div v-if="showMobileSearch" class="md:hidden border-t border-gray-100 dark:border-gray-800 p-3 bg-gray-50 dark:bg-slate-800">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher..."
            class="w-full pl-10 pr-4 py-2.5 bg-white dark:bg-slate-700 border border-gray-200 dark:border-gray-600 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500"
            @keyup.enter="performSearch"
          >
          <svg class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
        </div>
      </div>
    </header>

    <!-- Click outside to close dropdowns -->
    <div
      v-if="showMoreCategories"
      class="fixed inset-0 z-40"
      @click="showMoreCategories = false"
    ></div>

    <!-- Main Content -->
    <main class="flex-1">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <slot />
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white dark:bg-slate-900 border-t border-gray-200 dark:border-gray-800 mt-auto">
      <div class="max-w-7xl mx-auto px-4 py-8">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mb-8">
          <!-- About -->
          <div>
            <div class="flex items-center gap-2 mb-4">
              <div v-if="shop?.logo_url" class="w-8 h-8 rounded-lg overflow-hidden">
                <img :src="shop.logo_url" class="w-full h-full object-cover" alt="">
              </div>
              <span class="font-bold text-gray-900 dark:text-white">{{ shop?.name }}</span>
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 line-clamp-3">
              {{ shop?.description || 'Bienvenue dans notre boutique en ligne.' }}
            </p>
          </div>

          <!-- Quick Links -->
          <div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-4">Navigation</h4>
            <ul class="space-y-2 text-sm">
              <li><NuxtLink :to="`/shops/${shopSlug}`" class="text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400">Accueil</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/categories`" class="text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400">Catégories</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/my-orders`" class="text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400">Mes Commandes</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/favorites`" class="text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400">Favoris</NuxtLink></li>
            </ul>
          </div>

          <!-- Categories -->
          <div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-4">Catégories</h4>
            <ul class="space-y-2 text-sm">
              <li v-for="cat in categories.slice(0, 5)" :key="cat.id">
                <NuxtLink :to="`/shops/${shopSlug}?category=${cat.slug}`" class="text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400">
                  {{ cat.name }}
                </NuxtLink>
              </li>
            </ul>
          </div>

          <!-- Contact -->
          <div>
            <h4 class="font-bold text-gray-900 dark:text-white mb-4">Contact</h4>
            <ul class="space-y-2 text-sm text-gray-500 dark:text-gray-400">
              <li v-if="shop?.phone" class="flex items-center gap-2">
                <span>📞</span> {{ shop.phone }}
              </li>
              <li v-if="shop?.email" class="flex items-center gap-2">
                <span>✉️</span> {{ shop.email }}
              </li>
              <li v-if="shop?.address" class="flex items-center gap-2">
                <span>📍</span> {{ shop.address }}
              </li>
            </ul>
          </div>
        </div>

        <!-- Bottom Bar -->
        <div class="pt-6 border-t border-gray-200 dark:border-gray-800 flex flex-col md:flex-row items-center justify-between gap-4">
          <p class="text-sm text-gray-500 dark:text-gray-400">
            © {{ new Date().getFullYear() }} {{ shop?.name }}. Tous droits réservés.
          </p>
          <div class="flex items-center gap-4">
            <span class="text-xs text-gray-400">Propulsé par</span>
            <NuxtLink to="/" class="text-sm font-bold text-indigo-600 dark:text-indigo-400 hover:underline">
              Zekora
            </NuxtLink>
          </div>
        </div>
      </div>
    </footer>

    <!-- Floating Cart Button (Mobile) -->
    <div v-if="cartStore.itemCount > 0" class="fixed bottom-4 right-4 z-40 md:hidden">
      <NuxtLink
        :to="`/shops/${shopSlug}/cart`"
        class="flex items-center gap-2 px-4 py-3 bg-indigo-600 text-white rounded-full shadow-lg shadow-indigo-500/30"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
        </svg>
        <span class="font-bold">{{ cartStore.itemCount }}</span>
        <span class="text-sm">•</span>
        <span class="text-sm font-medium">{{ formatPrice(cartStore.total) }}</span>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()

const shopSlug = computed(() => route.params.slug as string)
const shop = ref<Shop | null>(null)
const categories = ref<Category[]>([])
const searchQuery = ref('')
const showMobileSearch = ref(false)
const showMoreCategories = ref(false)

const favoritesCount = computed(() => {
  if (typeof window === 'undefined') return 0
  try {
    const favs = JSON.parse(localStorage.getItem(`shop_favorites_${shopSlug.value}`) || '[]')
    return favs.length
  } catch {
    return 0
  }
})

const isExactHome = computed(() => {
  return route.path === `/shops/${shopSlug.value}` && !route.query.category && !route.query.search
})

const isActive = (section: string) => {
  return route.path.includes(`/shops/${shopSlug.value}/${section}`)
}

const formatPrice = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: shop.value?.currency || 'XOF',
    minimumFractionDigits: 0
  }).format(amount)
}

const performSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      path: `/shops/${shopSlug.value}`,
      query: { search: searchQuery.value.trim() }
    })
    showMobileSearch.value = false
  }
}

const loadShopData = async () => {
  if (!shopSlug.value) return

  try {
    shop.value = await shopApi.getShop(shopSlug.value)

    const catResult = await shopApi.listCategories(shopSlug.value)
    categories.value = catResult.categories || []
  } catch (e) {
    console.error('Failed to load shop data', e)
  }
}

watch(shopSlug, loadShopData, { immediate: true })

onMounted(() => {
  cartStore.loadFromStorage()
})

// Provide to child pages
provide('shop', shop)
provide('shopSlug', shopSlug)
provide('categories', categories)
provide('formatPrice', formatPrice)
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
