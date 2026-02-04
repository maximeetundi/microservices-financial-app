<template>
  <div class="min-h-screen bg-gray-100 dark:bg-slate-950">
    <!-- Top Bar -->
    <div class="bg-slate-900 text-white text-xs py-2">
      <div class="max-w-7xl mx-auto px-4 flex justify-between items-center">
        <div class="flex items-center gap-4">
          <span>📞 +221 XX XXX XX XX</span>
          <span class="hidden sm:inline">|</span>
          <span class="hidden sm:inline">📧 contact@{{ shop?.name || 'boutique' }}</span>
        </div>
        <div class="flex items-center gap-4">
          <NuxtLink to="/shops" class="hover:text-indigo-400 transition-colors">🏪 Marketplace</NuxtLink>
          <span>|</span>
          <button @click="toggleDarkMode" class="hover:text-indigo-400">
            {{ isDark ? '☀️' : '🌙' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Main Header -->
    <header class="bg-white dark:bg-slate-900 shadow-lg sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-20">
          <!-- Logo & Shop Name -->
          <NuxtLink :to="`/shops/${shopSlug}`" class="flex items-center gap-4">
            <div v-if="shop?.logo_url" class="w-14 h-14 rounded-xl overflow-hidden shadow-lg">
              <img :src="shop.logo_url" class="w-full h-full object-cover" alt="">
            </div>
            <div v-else class="w-14 h-14 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold shadow-lg">
              {{ shop?.name?.charAt(0) || 'S' }}
            </div>
            <div class="hidden md:block">
              <h1 class="text-xl font-bold text-gray-900 dark:text-white">{{ shop?.name || 'Boutique' }}</h1>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ shop?.description?.slice(0, 30) }}...</p>
            </div>
          </NuxtLink>

          <!-- Search Bar (Desktop) -->
          <div class="hidden lg:flex flex-1 max-w-xl mx-8">
            <div class="flex w-full">
              <select 
                v-model="searchCategory"
                class="px-4 py-3 bg-gray-100 dark:bg-slate-800 border-r border-gray-200 dark:border-gray-700 rounded-l-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">Toutes catégories</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.slug">{{ cat.name }}</option>
              </select>
              <input 
                v-model="searchQuery"
                @keyup.enter="doSearch"
                type="text" 
                placeholder="Rechercher des produits..." 
                class="flex-1 px-4 py-3 bg-gray-100 dark:bg-slate-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
              <button 
                @click="doSearch"
                class="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-r-xl transition-colors"
              >
                🔍
              </button>
            </div>
          </div>

          <!-- Right Actions -->
          <div class="flex items-center gap-2 sm:gap-4">
            <!-- Mobile Search Toggle -->
            <button @click="showMobileSearch = !showMobileSearch" class="lg:hidden p-3 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl">
              🔍
            </button>

            <!-- Favorites -->
            <NuxtLink 
              :to="`/shops/${shopSlug}/favorites`"
              class="relative p-3 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl transition-colors"
            >
              <span class="text-xl">❤️</span>
              <span v-if="favoritesCount > 0" class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center font-bold">
                {{ favoritesCount }}
              </span>
            </NuxtLink>

            <!-- Cart -->
            <NuxtLink 
              :to="`/shops/${shopSlug}/cart`"
              class="relative flex items-center gap-2 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl transition-colors"
            >
              <span class="text-xl">🛒</span>
              <span class="hidden sm:inline font-medium">{{ formatPrice(cartStore.total) }}</span>
              <span v-if="cartStore.itemCount > 0" class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white text-xs rounded-full flex items-center justify-center font-bold shadow-lg">
                {{ cartStore.itemCount }}
              </span>
            </NuxtLink>

            <!-- User Menu -->
            <div class="relative hidden sm:block">
              <button 
                @click="showUserMenu = !showUserMenu"
                class="p-3 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                <span class="text-xl">👤</span>
              </button>
              <div 
                v-if="showUserMenu"
                class="absolute right-0 top-full mt-2 w-48 bg-white dark:bg-slate-800 rounded-xl shadow-2xl border border-gray-100 dark:border-gray-700 py-2 z-50"
              >
                <NuxtLink :to="`/shops/${shopSlug}/my-orders`" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700">📦 Mes commandes</NuxtLink>
                <NuxtLink :to="`/shops/${shopSlug}/favorites`" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700">❤️ Mes favoris</NuxtLink>
                <hr class="my-2 border-gray-100 dark:border-gray-700">
                <NuxtLink to="/profile" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700">⚙️ Mon compte</NuxtLink>
              </div>
            </div>

            <!-- Mobile Menu -->
            <button @click="showMobileMenu = true" class="lg:hidden p-3 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl">
              ☰
            </button>
          </div>
        </div>

        <!-- Mobile Search -->
        <div v-if="showMobileSearch" class="lg:hidden pb-4">
          <div class="flex">
            <input 
              v-model="searchQuery"
              @keyup.enter="doSearch"
              type="text" 
              placeholder="Rechercher..." 
              class="flex-1 px-4 py-3 bg-gray-100 dark:bg-slate-800 text-gray-900 dark:text-white rounded-l-xl focus:outline-none"
            >
            <button @click="doSearch" class="px-4 py-3 bg-indigo-600 text-white rounded-r-xl">🔍</button>
          </div>
        </div>
      </div>

      <!-- Category Navigation -->
      <div class="border-t border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-slate-800/50">
        <div class="max-w-7xl mx-auto px-4">
          <div class="flex items-center gap-1 overflow-x-auto py-3 scrollbar-hide">
            <NuxtLink 
              :to="`/shops/${shopSlug}`"
              class="flex-shrink-0 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              :class="!route.query.category ? 'bg-indigo-600 text-white' : 'text-gray-600 dark:text-gray-300 hover:bg-white dark:hover:bg-slate-700'"
            >
              🏠 Tous
            </NuxtLink>
            <NuxtLink 
              v-for="cat in categories" 
              :key="cat.id"
              :to="`/shops/${shopSlug}?category=${cat.slug}`"
              class="flex-shrink-0 px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap"
              :class="route.query.category === cat.slug ? 'bg-indigo-600 text-white' : 'text-gray-600 dark:text-gray-300 hover:bg-white dark:hover:bg-slate-700'"
            >
              {{ cat.icon || '📁' }} {{ cat.name }}
            </NuxtLink>
            <NuxtLink 
              :to="`/shops/${shopSlug}/categories`"
              class="flex-shrink-0 px-4 py-2 text-indigo-600 dark:text-indigo-400 text-sm font-medium hover:underline"
            >
              Voir tout →
            </NuxtLink>
          </div>
        </div>
      </div>
    </header>

    <!-- Mobile Sidebar -->
    <div v-if="showMobileMenu" class="fixed inset-0 z-50 lg:hidden">
      <div class="absolute inset-0 bg-black/50" @click="showMobileMenu = false"></div>
      <aside class="absolute inset-y-0 left-0 w-80 bg-white dark:bg-slate-900 shadow-2xl overflow-y-auto">
        <!-- Shop Info -->
        <div class="p-6 bg-gradient-to-r from-indigo-600 to-purple-600 text-white">
          <div class="flex items-center gap-4">
            <div v-if="shop?.logo_url" class="w-16 h-16 rounded-xl overflow-hidden">
              <img :src="shop.logo_url" class="w-full h-full object-cover" alt="">
            </div>
            <div v-else class="w-16 h-16 rounded-xl bg-white/20 flex items-center justify-center text-2xl font-bold">
              {{ shop?.name?.charAt(0) || 'S' }}
            </div>
            <div>
              <h2 class="font-bold text-lg">{{ shop?.name }}</h2>
              <p class="text-sm text-white/80">{{ categories.length }} catégories</p>
            </div>
          </div>
        </div>

        <!-- Menu Links -->
        <nav class="p-4 space-y-1">
          <NuxtLink 
            :to="`/shops/${shopSlug}`"
            @click="showMobileMenu = false"
            class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800"
          >
            <span class="text-xl">🏠</span>
            Accueil
          </NuxtLink>
          <NuxtLink 
            :to="`/shops/${shopSlug}/categories`"
            @click="showMobileMenu = false"
            class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800"
          >
            <span class="text-xl">📂</span>
            Catégories
          </NuxtLink>
          <NuxtLink 
            :to="`/shops/${shopSlug}/favorites`"
            @click="showMobileMenu = false"
            class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800"
          >
            <span class="text-xl">❤️</span>
            Favoris
          </NuxtLink>
          <NuxtLink 
            :to="`/shops/${shopSlug}/my-orders`"
            @click="showMobileMenu = false"
            class="flex items-center gap-3 px-4 py-3 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800"
          >
            <span class="text-xl">📦</span>
            Mes commandes
          </NuxtLink>

          <hr class="my-4 border-gray-100 dark:border-gray-800">

          <p class="px-4 text-xs font-bold text-gray-400 uppercase">Catégories</p>
          <NuxtLink 
            v-for="cat in categories" 
            :key="cat.id"
            :to="`/shops/${shopSlug}?category=${cat.slug}`"
            @click="showMobileMenu = false"
            class="flex items-center justify-between px-4 py-2 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-slate-800"
          >
            <span>{{ cat.icon || '📁' }} {{ cat.name }}</span>
            <span class="text-xs text-gray-400">{{ cat.product_count || 0 }}</span>
          </NuxtLink>
        </nav>

        <!-- Back to Marketplace -->
        <div class="p-4 border-t border-gray-100 dark:border-gray-800">
          <NuxtLink 
            to="/shops"
            class="flex items-center justify-center gap-2 w-full px-4 py-3 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-xl font-medium"
          >
            ← Retour au Marketplace
          </NuxtLink>
        </div>
      </aside>
    </div>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 py-8">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-slate-900 text-white mt-16">
      <div class="max-w-7xl mx-auto px-4 py-12">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <!-- Shop Info -->
          <div class="md:col-span-2">
            <div class="flex items-center gap-3 mb-4">
              <div v-if="shop?.logo_url" class="w-12 h-12 rounded-xl overflow-hidden">
                <img :src="shop.logo_url" class="w-full h-full object-cover" alt="">
              </div>
              <h3 class="text-xl font-bold">{{ shop?.name }}</h3>
            </div>
            <p class="text-gray-400 text-sm mb-4">{{ shop?.description }}</p>
            <div class="flex gap-3">
              <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-indigo-600 transition-colors">📘</a>
              <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-indigo-600 transition-colors">📸</a>
              <a href="#" class="w-10 h-10 bg-white/10 rounded-full flex items-center justify-center hover:bg-indigo-600 transition-colors">🐦</a>
            </div>
          </div>

          <!-- Quick Links -->
          <div>
            <h4 class="font-bold mb-4">Liens rapides</h4>
            <ul class="space-y-2 text-sm text-gray-400">
              <li><NuxtLink :to="`/shops/${shopSlug}`" class="hover:text-white">Accueil</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/categories`" class="hover:text-white">Catégories</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/favorites`" class="hover:text-white">Favoris</NuxtLink></li>
              <li><NuxtLink :to="`/shops/${shopSlug}/my-orders`" class="hover:text-white">Mes commandes</NuxtLink></li>
            </ul>
          </div>

          <!-- Contact -->
          <div>
            <h4 class="font-bold mb-4">Contact</h4>
            <ul class="space-y-2 text-sm text-gray-400">
              <li>📧 contact@{{ shop?.name?.toLowerCase().replace(/\s/g, '') }}.com</li>
              <li>📞 +221 XX XXX XX XX</li>
              <li>📍 Dakar, Sénégal</li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-800 mt-8 pt-8 flex flex-col md:flex-row justify-between items-center text-sm text-gray-500">
          <p>© 2024 {{ shop?.name }}. Tous droits réservés.</p>
          <div class="flex gap-4 mt-4 md:mt-0">
            <span>🔒 Paiement sécurisé</span>
            <span>🚚 Livraison rapide</span>
            <span>💬 Support 24/7</span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import { useColorMode } from '#imports'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()
const colorMode = useColorMode()

const shopSlug = computed(() => route.params.slug as string)
const shop = ref<Shop | null>(null)
const categories = ref<Category[]>([])

const showMobileMenu = ref(false)
const showMobileSearch = ref(false)
const showUserMenu = ref(false)
const searchQuery = ref('')
const searchCategory = ref('')
const isDark = computed(() => colorMode.value === 'dark')

const favoritesCount = computed(() => {
  if (typeof window === 'undefined') return 0
  const favs = JSON.parse(localStorage.getItem(`shop_favorites_${shopSlug.value}`) || '[]')
  return favs.length
})

const formatPrice = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: shop.value?.currency || 'XOF' }).format(amount)
}

const toggleDarkMode = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}

const doSearch = () => {
  const query: Record<string, string> = {}
  if (searchQuery.value) query.search = searchQuery.value
  if (searchCategory.value) query.category = searchCategory.value
  router.push({ path: `/shops/${shopSlug.value}`, query })
  showMobileSearch.value = false
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

// Close menus on click outside
onMounted(() => {
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('[data-user-menu]')) {
      showUserMenu.value = false
    }
  })
})

// Provide to child pages
provide('shop', shop)
provide('shopSlug', shopSlug)
provide('categories', categories)
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
