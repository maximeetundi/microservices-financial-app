<template>
  <div class="container mx-auto px-4 py-8">
    <!-- Header -->
    <div class="flex justify-between items-center mb-8">
      <div>
        <NuxtLink :to="`/shops/manage/${slug}`" class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 flex items-center gap-1 mb-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
          Retour au tableau de bord
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Produits</h1>
        <p class="text-gray-500 dark:text-gray-400">Gérez le catalogue de votre boutique</p>
      </div>
      <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
        Nouveau produit
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 p-4 rounded-lg mb-8">
      {{ error }}
    </div>

    <!-- Empty State -->
    <div v-else-if="products.length === 0" class="text-center py-16 bg-white dark:bg-slate-800 rounded-lg border border-gray-100 dark:border-gray-700 shadow-sm">
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 mb-4">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">Aucun produit</h3>
      <p class="text-gray-500 dark:text-gray-400 mb-6">Commencez par ajouter votre premier produit à la boutique.</p>
      <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="text-indigo-600 dark:text-indigo-400 font-medium hover:underline">
        Créer un produit →
      </NuxtLink>
    </div>

    <!-- Products List -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="product in products" :key="product.id" class="bg-white dark:bg-slate-800 rounded-lg border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition-shadow overflow-hidden group">
        <div class="relative h-48 bg-gray-100 dark:bg-slate-700">
          <img v-if="product.images && product.images.length > 0" :src="product.images[0]" :alt="product.name" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full flex items-center justify-center text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-slate-800">
            <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
          </div>
          
          <div class="absolute top-2 right-2">
             <span class="px-2 py-1 text-xs font-bold rounded-full bg-white/90 dark:bg-slate-900/90 shadow-sm"
               :class="product.status === 'active' ? 'text-green-600' : 'text-gray-500'">
               {{ product.status }}
             </span>
          </div>
        </div>
        
        <div class="p-4">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-1 truncate">{{ product.name }}</h3>
          <div class="flex justify-between items-center mb-4">
            <span class="text-indigo-600 dark:text-indigo-400 font-bold">{{ product.price }} {{ product.currency }}</span>
            <span class="text-sm text-gray-500">Stock: {{ product.stock }}</span>
          </div>
          
          <div class="flex gap-2 pt-4 border-t border-gray-100 dark:border-gray-700">
             <NuxtLink :to="`/shops/manage/${slug}/products/${product.id}`" class="flex-1 py-2 text-center text-sm font-medium text-gray-700 dark:text-gray-200 bg-gray-50 dark:bg-slate-700 rounded hover:bg-gray-100 dark:hover:bg-slate-600 transition-colors">
               Modifier
             </NuxtLink>
             <button @click="deleteProduct(product)" class="p-2 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors">
               <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
             </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type Product } from '@/composables/useShopApi'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const products = ref<Product[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const fetchProducts = async () => {
  try {
    loading.value = true
    const response = await shopApi.listProducts(slug)
    products.value = response.products || []
  } catch (e: any) {
    console.error('Error fetching products:', e)
    error.value = e.message || 'Impossible de charger les produits'
  } finally {
    loading.value = false
  }
}

const deleteProduct = async (product: Product) => {
  if (!confirm(`Voulez-vous vraiment supprimer le produit "${product.name}" ?`)) return

  try {
    await shopApi.deleteProduct(product.id)
    // Remove from local list
    products.value = products.value.filter(p => p.id !== product.id)
  } catch (e: any) {
    alert('Erreur: ' + e.message)
  }
}

onMounted(() => {
  fetchProducts()
})

definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>
