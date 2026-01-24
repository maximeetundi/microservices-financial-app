<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileMenuOpen" 
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      @click="isMobileMenuOpen = false"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out lg:transform-none"
      :class="isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo / Shop Name -->
      <div class="h-16 flex items-center px-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center gap-2">
           <div class="h-8 w-8 rounded bg-indigo-600 flex items-center justify-center text-white font-bold text-lg">
             {{ shopName ? shopName.charAt(0).toUpperCase() : 'S' }}
           </div>
           <span class="font-bold text-gray-900 dark:text-white truncate max-w-[150px]">
             {{ shopName || 'Chargement...' }}
           </span>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="p-4 space-y-1">
        <NuxtLink 
          :to="`/shops/manage/${slug}`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}`) && !isActive(`/shops/manage/${slug}/`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <HomeIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`/shops/manage/${slug}`) && !isActive(`/shops/manage/${slug}/`) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300'" />
          Tableau de bord
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/products`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/products`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <ShoppingBagIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`/shops/manage/${slug}/products`) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300'" />
          Produits
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/orders`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/orders`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <ClipboardDocumentListIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`/shops/manage/${slug}/orders`) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300'" />
          Commandes
        </NuxtLink>
        
        <NuxtLink 
          :to="`/shops/manage/${slug}/categories`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/categories`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <TagIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`/shops/manage/${slug}/categories`) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300'" />
          Catégories
        </NuxtLink>

        <NuxtLink 
          :to="`/shops/manage/${slug}/settings`"
          class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors group"
          :class="isActive(`/shops/manage/${slug}/settings`) ? 'bg-indigo-50 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <Cog6ToothIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`/shops/manage/${slug}/settings`) ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300'" />
          Paramètres
        </NuxtLink>
      </nav>

      <!-- Bottom Actions -->
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 dark:border-gray-700">
         <NuxtLink 
           :to="`/shops/${slug}`"
           target="_blank"
           class="flex items-center justify-center w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600"
         >
           <EyeIcon class="h-4 w-4 mr-2" />
           Voir la boutique
         </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Mobile Header -->
      <header class="bg-white dark:bg-slate-800 shadow-sm lg:hidden h-16 flex items-center justify-between px-4 z-30">
        <button @click="isMobileMenuOpen = true" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
          <Bars3Icon class="h-6 w-6" />
        </button>
        <span class="font-bold text-gray-900 dark:text-white">{{ shopName }}</span>
        <div class="w-6"></div> <!-- Spacer for centering -->
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-8">
        <div class="max-w-7xl mx-auto">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
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
