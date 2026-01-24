<template>
  <div class="container mx-auto px-4 py-8 max-w-2xl">
    <!-- Header -->
    <div class="mb-8">
      <NuxtLink :to="`/shops/manage/${slug}/products`" class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 flex items-center gap-1 mb-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
        Retour aux produits
      </NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Nouveau Produit</h1>
    </div>

    <!-- Form -->
    <form @submit.prevent="submitProduct" class="bg-white dark:bg-slate-800 shadow rounded-lg p-6 space-y-6">
      
      <!-- Basic Info -->
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nom du produit *</label>
          <input v-model="form.name" type="text" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500 focus:border-indigo-500" placeholder="ex: T-shirt Vintage">
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Prix *</label>
          <div class="relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-500 sm:text-sm">$</span>
            </div>
            <input v-model.number="form.price" type="number" min="0" step="0.01" required class="pl-7 w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500 focus:border-indigo-500" placeholder="0.00">
             <div class="absolute inset-y-0 right-0 max-w-[100px]">
               <select v-model="form.currency" class="h-full py-0 pl-2 pr-7 border-transparent bg-transparent text-gray-500 sm:text-sm rounded-md focus:ring-indigo-500 focus:border-indigo-500">
                 <option>USD</option>
                 <option>EUR</option>
                 <option>XOF</option>
               </select>
             </div>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description</label>
          <textarea v-model="form.description" rows="4" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500 focus:border-indigo-500"></textarea>
        </div>
      </div>

      <!-- Inventory -->
      <div class="grid grid-cols-2 gap-4">
        <div>
           <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Stock</label>
           <input v-model.number="form.stock" type="number" min="0" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500 focus:border-indigo-500">
        </div>
        <div>
           <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Poids (kg)</label>
           <input v-model.number="form.weight" type="number" min="0" step="0.1" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500 focus:border-indigo-500">
        </div>
      </div>

      <!-- Visibility -->
      <div class="flex items-center gap-3 py-2">
         <div class="flex items-center h-5">
           <input v-model="form.is_featured" id="featured" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
         </div>
         <label for="featured" class="font-medium text-gray-700 dark:text-gray-300 select-none">Mettre en avant ce produit</label>
      </div>

      <div class="pt-4 border-t border-gray-100 dark:border-gray-700 flex justify-end gap-3">
        <NuxtLink :to="`/shops/manage/${slug}/products`" class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium">
          Annuler
        </NuxtLink>
        <button type="submit" :disabled="submitting" class="px-6 py-2 bg-indigo-600 text-white rounded-lg font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30 flex items-center">
          <span v-if="submitting" class="animate-spin mr-2 h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
          {{ submitting ? 'Création...' : 'Créer le produit' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const slug = route.params.id as string
const shopId = ref('') // We don't have shopID from url slug directly in create, need to fetch shop or rely on backend to resolve slug? 
// The backend create product expects shop_id (ID not slug). We need to fetch shop first to get ID.

const form = ref({
  name: '',
  description: '',
  price: 0,
  currency: 'USD',
  stock: 10,
  weight: 0,
  is_featured: false,
  status: 'active'
})

const submitting = ref(false)

const submitProduct = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    await shopApi.createProduct({
      shop_id: shopId.value,
      ...form.value
    })
    router.push(`/shops/manage/${slug}/products`)
  } catch (e: any) {
    alert('Erreur: ' + e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
     // Fetch shop to get the currency and ID
     const shop = await shopApi.getShop(slug)
     shopId.value = shop.id
     form.value.currency = shop.currency
  } catch (e) {
    console.error('Failed to load shop info', e)
  }
})

definePageMeta({
  middleware: ['auth']
})
</script>
