<template>
  <div class="h-screen overflow-hidden bg-gray-50 dark:bg-gray-900 flex flex-col">
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileMenuOpen" 
      class="fixed inset-0 bg-black/50 z-40 backdrop-blur-sm"
      @click="isMobileMenuOpen = false"
    ></div>

    <!-- Mobile Sidebar (Hidden by default) -->
    <aside 
      v-if="isMobileMenuOpen"
      class="fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out flex flex-col"
    >
      <!-- Close Button -->
      <div class="h-16 flex items-center justify-between px-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <div class="h-8 w-8 rounded bg-indigo-600 flex items-center justify-center text-white font-bold">
            {{ shopName ? shopName.charAt(0).toUpperCase() : 'S' }}
          </div>
          <span class="font-bold text-gray-900 dark:text-white text-sm">{{ shopName || 'Boutique' }}</span>
        </div>
        <button
          @click="isMobileMenuOpen = false"
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700 transition-colors"
        >
          <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Mobile Navigation -->
      <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
        <NuxtLink
          to="/shops/my-shops"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700"
          @click="isMobileMenuOpen = false"
        >
          <span class="mr-3">←</span>
          Mes Boutiques
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}`) && !isActive(`/shops/manage/${slug}/`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <HomeIcon class="h-5 w-5 mr-3" />
          Tableau de bord
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/products`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/products`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <ShoppingBagIcon class="h-5 w-5 mr-3" />
          Produits
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/orders`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/orders`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <ClipboardDocumentListIcon class="h-5 w-5 mr-3" />
          Commandes
        </NuxtLink>
        
        <NuxtLink 
          :to="`/shops/manage/${slug}/categories`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/categories`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <TagIcon class="h-5 w-5 mr-3" />
          Catégories
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/settings`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/settings`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <Cog6ToothIcon class="h-5 w-5 mr-3" />
          Paramètres
        </NuxtLink>
      </nav>
    </aside>

    <div class="flex flex-1 min-h-0">
      <ShopSidebar class="hidden lg:flex" />

      <div class="flex-1 flex flex-col min-h-0">
        <!-- Main Header with Navigation -->
        <header class="bg-white dark:bg-slate-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-30">
          <div class="px-4 sm:px-6 lg:px-8">
            <!-- Top Row -->
            <div class="flex items-center justify-between h-16">
              <!-- Left: Hamburger + Shop Info -->
              <div class="flex items-center gap-4">
                <!-- Hamburger Menu (Mobile) -->
                <button
                  @click="isMobileMenuOpen = true"
                  class="lg:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700 transition-colors"
                >
                  <Bars3Icon class="h-6 w-6 text-gray-600 dark:text-gray-300" />
                </button>

                <!-- Shop Logo and Name -->
                <div class="flex items-center gap-3">
                  <div class="h-10 w-10 rounded bg-indigo-600 flex items-center justify-center text-white font-bold text-lg">
                    {{ shopName ? shopName.charAt(0).toUpperCase() : 'S' }}
                  </div>
                  <div>
                    <h1 class="text-lg font-bold text-gray-900 dark:text-white">{{ shopName || 'Boutique' }}</h1>
                    <p class="text-sm text-gray-600 dark:text-gray-400">Gestion de la boutique</p>
                  </div>
                </div>
              </div>

              <!-- Right: Actions -->
              <div class="flex items-center gap-3">
                <NuxtLink
                  to="/shops/my-shops"
                  class="hidden sm:inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors text-sm font-medium"
                >
                  ← Mes boutiques
                </NuxtLink>
                <NuxtLink 
                  :to="`/shops/${slug}`"
                  target="_blank"
                  class="inline-flex items-center px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  <EyeIcon class="h-4 w-4 mr-2" />
                  Voir la boutique
                </NuxtLink>
              </div>
            </div>
          </div>
        </header>

        <!-- Page Content -->
        <main class="flex-1 overflow-y-auto min-h-0">
          <slot />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
import ShopSidebar from '~/components/shops/ShopSidebar.vue'
import { 
  HomeIcon, 
  ShoppingBagIcon, 
  ClipboardDocumentListIcon, 
  TagIcon, 
  Cog6ToothIcon, 
  EyeIcon,
  Bars3Icon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string
const isMobileMenuOpen = ref(false)
const shopName = ref('')

const isActive = (path: string) => {
  // Exact match or subpath match? 
  // For dashboard, exact match. For others, startsWith
  if (path.endsWith(`/${slug}`)) {
     return route.path === path
  }
  return route.path.startsWith(path)
}

// Fetch shop name for sidebar
watchEffect(async () => {
   if (slug) {
     try {
       const shop = await shopApi.getShop(slug)
       shopName.value = shop.name
     } catch (e) {
       console.error('Layout failed to load shop', e)
     }
   }
})
</script>
