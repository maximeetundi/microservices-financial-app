<template>
  <aside class="w-64 bg-white dark:bg-slate-900 border-r border-gray-200 dark:border-gray-700 flex-shrink-0 flex flex-col h-full rounded-l-2xl">
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
    <nav class="flex-1 p-3 space-y-1 overflow-y-auto">
      <NuxtLink 
        :to="`/shops/manage/${slug}`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}`) && !isActive(`/shops/manage/${slug}/`) }"
      >
        <HomeIcon class="h-5 w-5 mr-3" />
        Tableau de bord
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Catalogue</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/products`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}/products`) }"
      >
        <ShoppingBagIcon class="h-5 w-5 mr-3" />
        Produits
      </NuxtLink>
      
      <NuxtLink 
        :to="`/shops/manage/${slug}/categories`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}/categories`) }"
      >
        <TagIcon class="h-5 w-5 mr-3" />
        Catégories
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Ventes</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/orders`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}/orders`) }"
      >
        <ClipboardDocumentListIcon class="h-5 w-5 mr-3" />
        Commandes
      </NuxtLink>

      <div class="pt-4 pb-2 px-3">
        <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Administration</p>
      </div>

      <NuxtLink 
        :to="`/shops/manage/${slug}/managers`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}/managers`) }"
      >
        <UsersIcon class="h-5 w-5 mr-3" />
        Équipe
      </NuxtLink>

      <NuxtLink 
        :to="`/shops/manage/${slug}/settings`"
        class="nav-item"
        :class="{ 'active': isActive(`/shops/manage/${slug}/settings`) }"
      >
        <Cog6ToothIcon class="h-5 w-5 mr-3" />
        Paramètres
      </NuxtLink>
    </nav>

    <!-- Bottom Actions -->
    <div class="p-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800/50">
      <NuxtLink 
        :to="`/shops/${slug}`"
        target="_blank"
        class="flex items-center justify-center w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors"
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

<style scoped>
.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
  border-radius: 8px;
  transition: all 0.2s;
  margin-bottom: 2px;
}

.dark .nav-item {
  color: #94a3b8;
}

.nav-item:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.dark .nav-item:hover {
  background-color: #1e293b;
  color: #f1f5f9;
}

.nav-item.active {
  background-color: #eff6ff;
  color: #4f46e5;
}

.dark .nav-item.active {
  background-color: rgba(99, 102, 241, 0.15);
  color: #818cf8;
}
</style>
