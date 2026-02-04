<template>
  <div class="min-h-screen bg-gray-50 dark:bg-slate-950 flex">
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileMenuOpen" 
      class="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm"
      @click="isMobileMenuOpen = false"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed lg:static inset-y-0 left-0 z-50 w-72 bg-white dark:bg-slate-900 border-r border-gray-200 dark:border-gray-800 transform transition-transform duration-300 ease-in-out lg:transform-none flex flex-col shadow-xl lg:shadow-none"
      :class="isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Shop Header -->
      <div class="p-5 border-b border-gray-100 dark:border-gray-800">
        <div class="flex items-center gap-4">
          <div v-if="shop?.logo_url" class="h-14 w-14 rounded-2xl overflow-hidden flex-shrink-0 shadow-lg">
            <img :src="shop.logo_url" class="h-full w-full object-cover" />
          </div>
          <div v-else class="h-14 w-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-2xl flex-shrink-0 shadow-lg">
            {{ shop?.name?.charAt(0)?.toUpperCase() || 'S' }}
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="font-bold text-lg text-gray-900 dark:text-white truncate">
              {{ shop?.name || 'Chargement...' }}
            </h2>
            <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ shop?.description || 'Boutique en ligne' }}</p>
          </div>
        </div>
      </div>

      <!-- Search -->
      <div class="p-4">
        <div class="relative">
          <input 
            v-model="searchQuery"
            @keyup.enter="handleSearch"
            type="text" 
            placeholder="Rechercher..." 
            class="w-full pl-10 pr-4 py-3 bg-gray-100 dark:bg-slate-800 border-none rounded-xl focus:ring-2 focus:ring-indigo-500 dark:text-white text-sm transition-all"
          >
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">🔍</span>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-3 space-y-1 overflow-y-auto">
        <!-- Main Navigation -->
        <NuxtLink 
          :to="`/shops/${shopSlug}`"
          class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all group"
          :class="isExactActive ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">🏠</span>
          Accueil
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/${shopSlug}/categories`"
          class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all group"
          :class="isActive('categories') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">📂</span>
          Catégories
          <span class="ml-auto text-xs bg-gray-100 dark:bg-slate-700 px-2 py-0.5 rounded-full">{{ categories.length }}</span>
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/${shopSlug}/favorites`"
          class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all group"
          :class="isActive('favorites') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">❤️</span>
          Mes Favoris
          <span v-if="favoritesCount > 0" class="ml-auto text-xs bg-red-100 dark:bg-red-900/30 text-red-600 px-2 py-0.5 rounded-full">{{ favoritesCount }}</span>
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/${shopSlug}/my-orders`"
          class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all group"
          :class="isActive('my-orders') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-800'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">📦</span>
          Mes Commandes
        </NuxtLink>

        <!-- Categories Section -->
        <div class="pt-6 pb-2">
          <p class="text-xs font-bold text-gray-400 uppercase tracking-wider px-4">Catégories</p>
        </div>

        <NuxtLink 
          v-for="cat in categories.slice(0, 8)" 
          :key="cat.id"
          :to="`/shops/${shopSlug}?category=${cat.slug}`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-xl transition-all group text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-slate-800 hover:text-gray-900 dark:hover:text-white"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-base">{{ cat.icon || '📁' }}</span>
          <span class="truncate">{{ cat.name }}</span>
          <span class="ml-auto text-xs text-gray-400">{{ cat.product_count || 0 }}</span>
        </NuxtLink>

        <NuxtLink 
          v-if="categories.length > 8"
          :to="`/shops/${shopSlug}/categories`"
          class="flex items-center px-4 py-2 text-xs font-medium text-indigo-600 dark:text-indigo-400 hover:underline"
          @click="isMobileMenuOpen = false"
        >
          Voir toutes les catégories →
        </NuxtLink>

        <!-- Tags Section -->
        <div class="pt-6 pb-2">
          <p class="text-xs font-bold text-gray-400 uppercase tracking-wider px-4">Tags Populaires</p>
        </div>

        <div class="px-4 flex flex-wrap gap-2">
          <button 
            v-for="tag in popularTags.slice(0, 10)" 
            :key="tag"
            @click="searchByTag(tag)"
            class="px-3 py-1.5 text-xs font-medium bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-gray-400 rounded-lg hover:bg-indigo-100 dark:hover:bg-indigo-900/30 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >
            #{{ tag }}
          </button>
        </div>

        <!-- User Stats -->
        <div class="pt-6 pb-2">
          <p class="text-xs font-bold text-gray-400 uppercase tracking-wider px-4">Mes Statistiques</p>
        </div>

        <div class="mx-4 p-4 bg-gradient-to-br from-indigo-50 to-purple-50 dark:from-slate-800 dark:to-slate-800 rounded-2xl">
          <div class="grid grid-cols-2 gap-4 text-center">
            <div>
              <p class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">{{ userStats.orders }}</p>
              <p class="text-xs text-gray-500">Commandes</p>
            </div>
            <div>
              <p class="text-2xl font-bold text-pink-600 dark:text-pink-400">{{ userStats.favorites }}</p>
              <p class="text-xs text-gray-500">Favoris</p>
            </div>
          </div>
        </div>
      </nav>

      <!-- Cart Summary -->
      <div class="p-4 border-t border-gray-100 dark:border-gray-800">
        <NuxtLink 
          :to="`/shops/${shopSlug}/cart`"
          class="flex items-center justify-between w-full px-4 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-medium transition-all shadow-lg shadow-indigo-500/30"
        >
          <div class="flex items-center gap-3">
            <span class="text-xl">🛒</span>
            <span>Panier</span>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="cartStore.itemCount > 0" class="w-6 h-6 bg-white/20 rounded-full flex items-center justify-center text-sm font-bold">
              {{ cartStore.itemCount }}
            </span>
            <span class="font-bold">{{ formatPrice(cartStore.total) }}</span>
          </div>
        </NuxtLink>
      </div>

      <!-- Back to Marketplace -->
      <div class="p-4 border-t border-gray-100 dark:border-gray-800">
        <NuxtLink 
          to="/shops"
          class="flex items-center justify-center w-full px-4 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl text-sm font-medium text-gray-600 dark:text-gray-400 bg-white dark:bg-slate-800 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
        >
          ← Retour au Marketplace
        </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Mobile Header -->
      <header class="bg-white dark:bg-slate-900 shadow-sm lg:hidden h-16 flex items-center justify-between px-4 z-30 border-b border-gray-200 dark:border-gray-800">
        <button @click="isMobileMenuOpen = true" class="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl transition-colors">
          <span class="text-2xl">☰</span>
        </button>
        <div class="flex items-center gap-3">
          <div v-if="shop?.logo_url" class="h-8 w-8 rounded-lg overflow-hidden">
            <img :src="shop.logo_url" class="h-full w-full object-cover" />
          </div>
          <span class="font-bold text-gray-900 dark:text-white">{{ shop?.name }}</span>
        </div>
        <NuxtLink :to="`/shops/${shopSlug}/cart`" class="p-2 relative">
          <span class="text-2xl">🛒</span>
          <span v-if="cartStore.itemCount > 0" class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center">
            {{ cartStore.itemCount }}
          </span>
        </NuxtLink>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto">
        <div class="p-4 lg:p-8">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()
const authStore = useAuthStore()

const shopSlug = computed(() => route.params.slug as string)
const isMobileMenuOpen = ref(false)
const shop = ref<Shop | null>(null)
const categories = ref<Category[]>([])
const popularTags = ref<string[]>([])
const searchQuery = ref('')

// User stats for this shop
const userStats = ref({
  orders: 0,
  favorites: 0
})

// Favorites count from localStorage
const favoritesCount = computed(() => {
  if (typeof window === 'undefined') return 0
  const favs = JSON.parse(localStorage.getItem(`shop_favorites_${shopSlug.value}`) || '[]')
  return favs.length
})

const formatPrice = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: shop.value?.currency || 'XOF' }).format(amount)
}

const isExactActive = computed(() => {
  return route.path === `/shops/${shopSlug.value}`
})

const isActive = (section: string) => {
  return route.path.includes(`/shops/${shopSlug.value}/${section}`)
}

const handleSearch = () => {
  router.push({ 
    path: `/shops/${shopSlug.value}`, 
    query: { search: searchQuery.value || undefined } 
  })
  isMobileMenuOpen.value = false
}

const searchByTag = (tag: string) => {
  router.push({ 
    path: `/shops/${shopSlug.value}`, 
    query: { tag } 
  })
  isMobileMenuOpen.value = false
}

const loadShopData = async () => {
  if (!shopSlug.value) return
  
  try {
    shop.value = await shopApi.getShop(shopSlug.value)
    
    // Load categories
    const catResult = await shopApi.listCategories(shopSlug.value)
    categories.value = catResult.categories || []
    
    // Load products to extract tags
    const prodResult = await shopApi.listProducts(shopSlug.value, 1, 100)
    const allTags = prodResult.products?.flatMap(p => p.tags || []) || []
    const tagCounts = allTags.reduce((acc: Record<string, number>, tag: string) => {
      acc[tag] = (acc[tag] || 0) + 1
      return acc
    }, {})
    popularTags.value = Object.entries(tagCounts)
      .sort((a, b) => (b[1] as number) - (a[1] as number))
      .slice(0, 10)
      .map(([tag]) => tag)
    
    // Load user stats (simplified)
    userStats.value.favorites = favoritesCount.value
    
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
provide('popularTags', popularTags)
</script>

<style scoped>
/* Custom scrollbar */
nav::-webkit-scrollbar {
  width: 4px;
}
nav::-webkit-scrollbar-thumb {
  @apply bg-gray-200 dark:bg-gray-700 rounded-full;
}
</style>
