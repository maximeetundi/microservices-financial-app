<template>
  <div class="container mx-auto px-4 py-8 max-w-4xl">
    <!-- Header -->
    <div class="flex justify-between items-center mb-8">
      <div>
        <NuxtLink :to="`/shops/manage/${slug}`" class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 flex items-center gap-1 mb-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
          Retour au tableau de bord
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Catégories</h1>
        <p class="text-gray-500 dark:text-gray-400">Organisez vos produits par rayons</p>
      </div>
      <button @click="showCreateModal = true" class="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
        Nouvelle catégorie
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
       <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Content -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-8">
       <!-- List -->
       <div class="space-y-4">
          <div v-if="categories.length === 0" class="text-center py-12 bg-gray-50 dark:bg-slate-800 rounded-lg border border-gray-100 dark:border-gray-700">
             <p class="text-gray-500">Aucune catégorie pour le moment.</p>
          </div>
          
          <div v-else class="bg-white dark:bg-slate-800 shadow rounded-lg overflow-hidden border border-gray-100 dark:border-gray-700">
             <div v-for="category in categories" :key="category.id" class="p-4 border-b border-gray-100 dark:border-gray-700 last:border-0 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                <div class="flex items-center gap-3">
                   <div class="w-10 h-10 rounded-lg bg-gray-100 dark:bg-slate-700 flex items-center justify-center text-xl">
                      ⚡
                   </div>
                   <div>
                      <h3 class="font-bold text-gray-900 dark:text-white">{{ category.name }}</h3>
                      <p class="text-xs text-gray-500">{{ category.product_count || 0 }} produits</p>
                   </div>
                </div>
                <div class="flex items-center gap-2">
                   <button @click="deleteCategory(category)" class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                   </button>
                </div>
             </div>
          </div>
       </div>

       <!-- Quick Helper -->
       <div class="bg-blue-50 dark:bg-blue-900/20 p-6 rounded-lg text-blue-800 dark:text-blue-200 self-start">
         <h3 class="font-bold mb-2">Pourquoi des catégories ?</h3>
         <p class="text-sm mb-4">Les catégories aident vos clients à trouver plus facilement ce qu'ils cherchent. Une boutique bien organisée vend mieux !</p>
         <ul class="text-sm space-y-2 list-disc list-inside opacity-80">
            <li>Essayez de rester simple (5-8 catégories max)</li>
            <li>Utilisez des noms clairs</li>
            <li>Regroupez les produits similaires</li>
         </ul>
       </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md p-6 shadow-2xl">
         <h3 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Nouvelle Catégorie</h3>
         
         <form @submit.prevent="createCategory">
            <div class="mb-4">
               <label class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">Nom</label>
               <input v-model="newCategory.name" type="text" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-800 focus:ring-indigo-500" placeholder="ex: Vêtements Femme">
            </div>
            
            <div class="mb-6">
               <label class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">Description</label>
               <textarea v-model="newCategory.description" rows="3" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-800 focus:ring-indigo-500"></textarea>
            </div>

            <div class="flex justify-end gap-3">
               <button type="button" @click="showCreateModal = false" class="px-4 py-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-lg">Annuler</button>
               <button type="submit" :disabled="submitting" class="px-4 py-2 bg-indigo-600 text-white rounded-lg font-bold hover:bg-indigo-700 flex items-center gap-2">
                  <span v-if="submitting" class="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
                  Créer
               </button>
            </div>
         </form>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type Category } from '@/composables/useShopApi'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const categories = ref<Category[]>([])
const loading = ref(true)
const showCreateModal = ref(false)
const submitting = ref(false)

const newCategory = ref({
  name: '',
  description: ''
})

const shopId = ref('')

const fetchCategories = async () => {
  try {
    loading.value = true
    const res = await shopApi.listCategories(slug)
    categories.value = res.categories || []
    
    // Also fetch shop ID if not already done
    if (!shopId.value) {
       const shop = await shopApi.getShop(slug)
       shopId.value = shop.id
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const createCategory = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    await shopApi.createCategory({
       shop_id: shopId.value,
       ...newCategory.value
    })
    
    // Reset and refresh
    newCategory.value = { name: '', description: '' }
    showCreateModal.value = false
    await fetchCategories()
  } catch (e) {
    alert('Erreur: ' + (e.message || 'Impossible de créer la catégorie'))
  } finally {
    submitting.value = false
  }
}

const deleteCategory = async (cat: Category) => {
   if (!confirm(`Supprimer la catégorie "${cat.name}" ?`)) return
   try {
      await shopApi.deleteCategory(cat.id)
      await fetchCategories()
   } catch (e) {
      alert('Impossible de supprimer')
   }
}

onMounted(() => {
  fetchCategories()
})

definePageMeta({
  middleware: ['auth']
})
</script>
