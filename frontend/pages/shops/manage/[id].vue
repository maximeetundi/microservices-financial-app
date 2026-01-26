<template>
  <div class="container mx-auto px-4 py-8">
    <div class="mb-6">
      <NuxtLink to="/shops" class="inline-flex items-center text-gray-500 hover:text-indigo-600 transition-colors">
        <span class="mr-2">←</span> Retour au Marketplace
      </NuxtLink>
    </div>
    <div v-if="loading" class="flex justify-center items-center min-h-[50vh]">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="error" class="text-center py-12">
      <div class="text-red-500 mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Erreur de chargement</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error }}</p>
      <NuxtLink to="/shops/my-shops" class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
        Retour à mes boutiques
      </NuxtLink>
    </div>

    <div v-else-if="shop" class="space-y-6">
      <!-- Header -->
      <div class="bg-white dark:bg-slate-800 shadow rounded-lg p-6">
        <div class="flex justify-between items-start">
          <div>
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">{{ shop.name }}</h1>
            <p class="text-gray-500 dark:text-gray-400 mb-4">{{ shop.description }}</p>
            <div class="flex items-center space-x-4">
              <span 
                class="px-3 py-1 rounded-full text-sm font-medium"
                :class="{
                  'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': shop.status === 'active',
                  'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': shop.status === 'pending',
                  'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': shop.status === 'suspended'
                }"
              >
                {{ shop.status }}
              </span>
              <span class="text-sm text-gray-500 dark:text-gray-400">
                Créée le {{ new Date(shop.created_at).toLocaleDateString() }}
              </span>
            </div>
          </div>
          <div class="flex space-x-3">
            <NuxtLink :to="`/shops/${shop.slug}`" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600">
              Voir la boutique
            </NuxtLink>
            <button class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700">
              Modifier
            </button>
          </div>
        </div>
      </div>

      <!-- Dashboard Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Stats Card -->
        <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
          <div class="flex items-center gap-2 mb-4">
             <CurrencyEuroIcon class="h-6 w-6 text-indigo-500" />
             <h3 class="text-lg font-medium text-gray-900 dark:text-white">Ventes du mois</h3>
          </div>
          <div class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">0.00 €</div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Vs mois dernier: 0%</p>
        </div>

        <!-- Orders Card -->
        <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
          <div class="flex items-center gap-2 mb-4">
             <ClipboardDocumentListIcon class="h-6 w-6 text-indigo-500" />
             <h3 class="text-lg font-medium text-gray-900 dark:text-white">Commandes</h3>
          </div>
          <div class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">0</div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">En attente: 0</p>
        </div>

        <!-- Products Card -->
        <NuxtLink :to="`/shops/manage/${slug}/products`" class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6 hover:shadow-md transition-shadow cursor-pointer border border-transparent hover:border-indigo-500">
          <div class="flex items-center gap-2 mb-4">
             <ShoppingBagIcon class="h-6 w-6 text-indigo-500" />
             <h3 class="text-lg font-medium text-gray-900 dark:text-white">Produits</h3>
          </div>
          <div class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">0</div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">En stock</p>
        </NuxtLink>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Gestion rapide</h3>
          <div class="space-y-3">
            <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="w-full flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group">
              <span class="font-medium text-gray-700 dark:text-gray-200 flex items-center gap-3">
                <PlusIcon class="h-5 w-5 text-gray-400 group-hover:text-indigo-500" />
                Ajouter un produit
              </span>
              <span class="text-gray-400">→</span>
            </NuxtLink>
            <NuxtLink :to="`/shops/manage/${slug}/categories`" class="w-full flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group">
              <span class="font-medium text-gray-700 dark:text-gray-200 flex items-center gap-3">
                <TagIcon class="h-5 w-5 text-gray-400 group-hover:text-indigo-500" />
                Gérer les catégories
              </span>
              <span class="text-gray-400">→</span>
            </NuxtLink>
            <NuxtLink :to="`/shops/manage/${slug}/managers`" class="w-full flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group">
              <span class="font-medium text-gray-700 dark:text-gray-200 flex items-center gap-3">
                <UsersIcon class="h-5 w-5 text-gray-400 group-hover:text-indigo-500" />
                Gérer l'équipe
              </span>
              <span class="text-gray-400">→</span>
            </NuxtLink>
            <NuxtLink :to="`/shops/manage/${slug}/settings`" class="w-full flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group">
              <span class="font-medium text-gray-700 dark:text-gray-200 flex items-center gap-3">
                <Cog6ToothIcon class="h-5 w-5 text-gray-400 group-hover:text-indigo-500" />
                Paramètres & Paiement
              </span>
              <span class="text-gray-400">→</span>
            </NuxtLink>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Dernières commandes</h3>
          <div class="text-center py-8 text-gray-500 dark:text-gray-400">
            Aucune commande récente
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
import { 
  BuildingStorefrontIcon, 
  PencilSquareIcon,
  CurrencyEuroIcon,
  ClipboardDocumentListIcon,
  ShoppingBagIcon,
  PlusIcon,
  TagIcon,
  UsersIcon,
  Cog6ToothIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const shop = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    loading.value = true
    const response = await shopApi.getShop(slug)
    shop.value = response
  } catch (e: any) {
    console.error('Error fetching shop:', e)
    error.value = e.message || 'Impossible de charger la boutique'
  } finally {
    loading.value = false
  }
})

definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>
