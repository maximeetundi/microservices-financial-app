<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileMenuOpen" 
      class="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm"
      @click="isMobileMenuOpen = false"
    ></div>

    

    <aside 
      class="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out lg:transform-none flex flex-col"
      :class="isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Shop Logo & Name -->
      <div class="h-16 flex items-center px-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center gap-3 flex-1 min-w-0">
          <div v-if="shop?.logo_url" class="h-10 w-10 rounded-xl overflow-hidden flex-shrink-0">
            <img :src="shop.logo_url" class="h-full w-full object-cover" />
          </div>
          <div v-else class="h-10 w-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg flex-shrink-0">
            {{ shop?.name?.charAt(0)?.toUpperCase() || 'S' }}
          </div>
          <div class="min-w-0 flex-1">
            <span class="font-bold text-gray-900 dark:text-white truncate block">
              {{ shop?.name || 'Chargement...' }}
            </span>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              Boutique
            </span>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="p-4 space-y-1 overflow-y-auto">
        <!-- NAVIGATION Section -->
        <div class="pb-2">
          <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4 mb-2">Navigation</p>
        </div>

        <NuxtLink 
          :to="`/shops/${shopSlug}`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isExactActive ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">🏠</span>
          Accueil
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/${shopSlug}/categories`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive('categories') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">📂</span>
          Toutes les Catégories
        </NuxtLink>

        <!-- MON ESPACE Section -->
        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4 mb-2">Mon Espace</p>
        </div>

        <NuxtLink 
          :to="`/shops/${shopSlug}/favorites`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive('favorites') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">❤️</span>
          Mes Favoris
          <span v-if="favoritesCount > 0" class="ml-auto text-xs bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 px-2 py-0.5 rounded-full font-semibold">{{ favoritesCount }}</span>
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/${shopSlug}/my-orders`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive('my-orders') ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">📦</span>
          Mes Commandes
        </NuxtLink>

        <NuxtLink 
          to="/cart"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isCartActive ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3 text-lg">🛒</span>
          Mon Panier
          <span v-if="cartStore.itemCount > 0" class="ml-auto text-xs bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 px-2 py-0.5 rounded-full font-semibold">{{ cartStore.itemCount }}</span>
        </NuxtLink>

        <!-- EXPLORER Section -->
        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4 mb-2">Explorer</p>
        </div>

        <div v-if="categories.length > 0" class="space-y-1">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 px-4 py-1">🏷️ Catégories</p>
          
          <NuxtLink 
            v-for="cat in categories.slice(0, 8)" 
            :key="cat.id"
            :to="`/shops/${shopSlug}?category=${cat.slug}`"
            class="flex items-center px-8 py-2 text-sm font-medium rounded-lg transition-colors group"
            :class="route.query.category === cat.slug ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <span class="mr-2 text-base">{{ cat.icon || '📁' }}</span>
            <span class="truncate flex-1">{{ cat.name }}</span>
            <span class="ml-auto text-xs text-gray-400">{{ cat.product_count || 0 }}</span>
          </NuxtLink>

          <NuxtLink 
            v-if="categories.length > 8"
            :to="`/shops/${shopSlug}/categories`"
            class="flex items-center px-8 py-1.5 text-xs font-medium text-indigo-600 dark:text-indigo-400 hover:underline"
            @click="isMobileMenuOpen = false"
          >
            → Voir toutes les catégories
          </NuxtLink>
        </div>

        <!-- Tags Section -->
        <div v-if="popularTags.length > 0" class="pt-3">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 px-4 py-1">#️⃣ Tags Populaires</p>
          <div class="px-4 flex flex-wrap gap-1.5 mt-2">
            <button 
              v-for="tag in popularTags.slice(0, 8)" 
              :key="tag"
              @click="searchByTag(tag)"
              class="px-2 py-1 text-xs font-medium bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-gray-400 rounded hover:bg-indigo-100 dark:hover:bg-indigo-900/30 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
            >
              #{{ tag }}
            </button>
          </div>
        </div>
      </nav>

      
      <!-- Bottom Actions -->
      <div class="p-4 border-t border-gray-200 dark:border-gray-700">
        <NuxtLink 
          to="/shops"
          class="flex items-center justify-center w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors"
        >
          ← Retour au Marketplace
        </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Mobile Header -->
      <header class="bg-white dark:bg-slate-800 shadow-sm lg:hidden h-16 flex items-center justify-between px-4 z-30 border-b border-gray-200 dark:border-gray-700">
        <button @click="isMobileMenuOpen = true" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
          <span class="text-2xl">☰</span>
        </button>
        <div class="flex items-center gap-2">
          <div v-if="shop?.logo_url" class="h-8 w-8 rounded-lg overflow-hidden">
            <img :src="shop.logo_url" class="h-full w-full object-cover" />
          </div>
          <span class="font-bold text-gray-900 dark:text-white">{{ shop?.name }}</span>
        </div>
        <NuxtLink to="/cart" class="relative">
          <span class="text-2xl">🛒</span>
          <span v-if="cartStore.itemCount > 0" class="absolute -top-2 -right-2 w-5 h-5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center">
            {{ cartStore.itemCount }}
          </span>
        </NuxtLink>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-8">
        <div class="max-w-7xl mx-auto">
          <slot />
        </div>
      </main>
    </div>

    <!-- Floating Cart Button -->
    <FloatingCartButton />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Shop, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import FloatingCartButton from '~/components/common/FloatingCartButton.vue'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()

const shopSlug = computed(() => route.params.slug as string)
const isMobileMenuOpen = ref(false)
const shop = ref<Shop | null>(null)
const categories = ref<Category[]>([])
const popularTags = ref<string[]>([])

const favoritesCount = computed(() => {
  if (typeof window === 'undefined') return 0
  const favs = JSON.parse(localStorage.getItem(`shop_favorites_${shopSlug.value}`) || '[]')
  return favs.length
})

const formatPrice = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: shop.value?.currency || 'XOF' }).format(amount)
}

const isExactActive = computed(() => {
  return route.path === `/shops/${shopSlug.value}` && !route.query.category
})

const isCartActive = computed(() => {
  return route.path === '/cart'
})

const isActive = (section: string) => {
  return route.path.includes(`/shops/${shopSlug.value}/${section}`)
}

const searchByTag = (tag: string) => {
  router.push({ path: `/shops/${shopSlug.value}`, query: { tag } })
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
    try {
      const prodResult = await shopApi.listProducts(shopSlug.value, 1, 50)
      const allTags = prodResult.products?.flatMap(p => p.tags || []) || []
      const tagCounts = allTags.reduce((acc: Record<string, number>, tag: string) => {
        acc[tag] = (acc[tag] || 0) + 1
        return acc
      }, {})
      popularTags.value = Object.entries(tagCounts)
        .sort((a, b) => (b[1] as number) - (a[1] as number))
        .slice(0, 10)
        .map(([tag]) => tag)
    } catch (e) {
      console.warn('Failed to load tags', e)
    }
    
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
