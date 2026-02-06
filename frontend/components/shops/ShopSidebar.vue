<template>
  <aside class="w-64 bg-surface border-r border-secondary-200 dark:border-secondary-800 flex-shrink-0 flex flex-col lg:rounded-l-2xl sticky top-0 h-screen overflow-hidden">
    <!-- Back to App -->
    <div class="p-4 pb-0">
      <NuxtLink to="/shops/my-shops" class="flex items-center text-sm font-medium text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400 transition-colors">
        <span class="mr-2">←</span> Mes Boutiques
      </NuxtLink>
    </div>

    <!-- Shop Context Header -->
    <div class="h-16 flex items-center px-4 border-b border-gray-200 dark:border-gray-700 mx-4 mt-2 mb-2">
      <div class="flex items-center gap-3 w-full">
        <div class="h-8 w-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white font-bold text-lg flex-shrink-0 shadow-sm">
          {{ shopName ? shopName.charAt(0).toUpperCase() : 'S' }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500 font-medium uppercase tracking-wider mb-0.5 dark:text-gray-400">Boutique</p>
          <h2 class="font-bold text-gray-900 dark:text-white truncate text-sm leading-tight">
            {{ shopName || 'Chargement...' }}
          </h2>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 min-h-0 p-3 space-y-1 overflow-y-auto">
      <NuxtLink 
        :to="`/shops/manage/${slug}`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}`) && !isActive(`/shops/manage/${slug}/`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <HomeIcon class="h-5 w-5 mr-3" />
        Tableau de bord
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Catalogue</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/products`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}/products`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <ShoppingBagIcon class="h-5 w-5 mr-3" />
        Produits
      </NuxtLink>
      
      <NuxtLink 
        :to="`/shops/manage/${slug}/categories`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}/categories`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <TagIcon class="h-5 w-5 mr-3" />
        Catégories
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Ventes</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/orders`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}/orders`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <ClipboardDocumentListIcon class="h-5 w-5 mr-3" />
        Commandes
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Administration</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/managers`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}/managers`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <UsersIcon class="h-5 w-5 mr-3" />
        Équipe
      </NuxtLink>

      <NuxtLink 
        :to="`/shops/manage/${slug}/settings`"
        class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-all mb-0.5"
        :class="isActive(`/shops/manage/${slug}/settings`) 
          ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400' 
          : 'text-gray-600 dark:text-slate-400 hover:bg-gray-50 hover:text-gray-900 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
      >
        <Cog6ToothIcon class="h-5 w-5 mr-3" />
        Paramètres
      </NuxtLink>
    </nav>

    <!-- Bottom Actions -->
    <!-- Bottom Actions -->
    <div class="p-4 border-t border-secondary-200 dark:border-secondary-800 bg-surface-hover/50 flex items-center gap-2 flex-shrink-0">
      <ThemeToggle />
      <NuxtLink 
        :to="`/shops/${slug}`"
        target="_blank"
        class="flex-1 flex items-center justify-center px-4 py-2 border border-secondary-200 dark:border-secondary-700 rounded-lg shadow-sm text-sm font-medium text-base bg-surface hover:bg-surface-hover transition-colors"
      >
        <EyeIcon class="h-4 w-4 mr-2" />
        Voir le site
      </NuxtLink>
    </div>
  </aside>
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
  UsersIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string
const shopName = ref('')

const isActive = (path: string) => {
  if (path.endsWith(`/${slug}`)) {
     return route.path === path
  }
  return route.path.startsWith(path)
}

watchEffect(async () => {
   if (slug) {
     try {
       const shop = await shopApi.getShop(slug)
       shopName.value = shop.name
     } catch (e) {
       console.error('Sidebar failed to load shop', e)
     }
   }
})
</script>
