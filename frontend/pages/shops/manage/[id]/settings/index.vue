<template>
  <div class="container mx-auto px-4 py-8 max-w-4xl">
     <!-- Header -->
    <div class="mb-8 border-b border-gray-200 dark:border-gray-700 pb-5">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
        <Cog6ToothIcon class="h-8 w-8 text-indigo-600" />
        Paramètres de la boutique
      </h1>
      <p class="mt-2 text-gray-500 dark:text-gray-400">Configurez les informations et préférences de votre boutique</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
       <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <form v-else @submit.prevent="saveSettings" class="space-y-8">
      
      <!-- General Info -->
      <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 border-b border-gray-100 dark:border-gray-700 pb-2 flex items-center gap-2">
           <BuildingStorefrontIcon class="h-6 w-6 text-indigo-500" />
           Informations Générales
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nom de la boutique</label>
            <input v-model="form.name" type="text" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description</label>
            <textarea v-model="form.description" rows="3" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500"></textarea>
          </div>
          <!-- Visibility Switch -->
          <div class="flex items-center gap-3">
             <div class="flex items-center h-5">
               <input v-model="form.is_public" id="public" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
             </div>
             <div>
               <label for="public" class="font-medium text-gray-700 dark:text-gray-300">Boutique Publique</label>
               <p class="text-xs text-gray-500">Si décoché, la boutique sera accessible uniquement sur invitation.</p>
             </div>
          </div>
        </div>
      </div>

      <!-- Contact & Address -->
      <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 border-b border-gray-100 dark:border-gray-700 pb-2 flex items-center gap-2">
          <MapPinIcon class="h-6 w-6 text-indigo-500" />
          Contact & Adresse
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email de contact</label>
            <input v-model="form.contact_info.email" type="email" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Téléphone</label>
            <input v-model="form.contact_info.phone" type="tel" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Site Web</label>
            <input v-model="form.contact_info.website" type="url" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500" placeholder="https://">
          </div>
          
          <div class="md:col-span-2 pt-4">
             <h3 class="font-bold text-gray-700 dark:text-gray-300 mb-3">Adresse</h3>
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Rue</label>
            <input v-model="form.address.street" type="text" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Ville</label>
            <input v-model="form.address.city" type="text" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Pays</label>
            <input v-model="form.address.country" type="text" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
          </div>
        </div>
      </div>

      <!-- Shop Settings (Delivery, Orders) -->
      <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 border-b border-gray-100 dark:border-gray-700 pb-2 flex items-center gap-2">
          <TruckIcon class="h-6 w-6 text-indigo-500" />
          Préférences de commande
        </h2>
        <div class="space-y-4">
           
           <div class="flex items-center justify-between">
              <div>
                 <label class="font-medium text-gray-900 dark:text-white">Autoriser le retrait en magasin (Pickup)</label>
                 <p class="text-xs text-gray-500">Les clients peuvent venir récupérer leur commande.</p>
              </div>
              <div class="flex items-center h-5">
                 <input v-model="form.settings.allow_pickup" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
              </div>
           </div>

           <div class="flex items-center justify-between">
              <div>
                 <label class="font-medium text-gray-900 dark:text-white">Autoriser la livraison</label>
                 <p class="text-xs text-gray-500">Vous proposez un service de livraison.</p>
              </div>
              <div class="flex items-center h-5">
                 <input v-model="form.settings.allow_delivery" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
              </div>
           </div>

           <div v-if="form.settings.allow_delivery" class="pl-4 border-l-2 border-indigo-100 dark:border-indigo-900 ml-1">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Frais de livraison standard</label>
              <div class="relative rounded-md shadow-sm max-w-xs">
                 <input v-model.number="form.settings.delivery_fee" type="number" min="0" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700 focus:ring-indigo-500">
              </div>
           </div>

           <div class="flex items-center justify-between pt-4 border-t border-gray-100 dark:border-gray-700">
              <div>
                 <label class="font-medium text-gray-900 dark:text-white">Acceptation automatique</label>
                 <p class="text-xs text-gray-500">Accepter automatiquement les nouvelles commandes payées.</p>
              </div>
              <div class="flex items-center h-5">
                 <input v-model="form.settings.auto_accept_orders" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
              </div>
           </div>

        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-end gap-4 pt-4">
        <NuxtLink :to="`/shops/manage/${slug}`" class="px-6 py-3 border border-gray-300 dark:border-gray-600 rounded-xl text-gray-700 dark:text-gray-300 font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors">
          Annuler
        </NuxtLink>
        <button type="submit" :disabled="saving" class="px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30 flex items-center gap-2">
          <span v-if="saving" class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
          {{ saving ? 'Enregistrement...' : 'Enregistrer les modifications' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
import { useToast } from "vue-toastification"
import { 
  Cog6ToothIcon, 
  BuildingStorefrontIcon, 
  MapPinIcon, 
  TruckIcon 
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const toast = useToast()
const slug = route.params.id as string

const loading = ref(true)
const saving = ref(false)
const shopId = ref('')

const form = ref({
  name: '',
  description: '',
  is_public: true,
  contact_info: {
    email: '',
    phone: '',
    website: ''
  },
  address: {
    street: '',
    city: '',
    country: '',
    postal_code: '',
    state: ''
  },
  settings: {
    allow_pickup: false,
    allow_delivery: false,
    delivery_fee: 0,
    min_order_amount: 0,
    max_order_amount: 0,
    auto_accept_orders: false
  }
})

const fetchShop = async () => {
  try {
    loading.value = true
    const shop = await shopApi.getShop(slug)
    shopId.value = shop.id
    
    // Populate form
    form.value.name = shop.name
    form.value.description = shop.description
    form.value.is_public = shop.is_public
    
    if (shop.contact_info) form.value.contact_info = { ...shop.contact_info }
    if (shop.address) form.value.address = { ...shop.address }
    if (shop.settings) form.value.settings = { ...shop.settings }

  } catch (e) {
    console.error('Failed to load shop', e)
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  if (!shopId.value) return
  
  try {
    saving.value = true
    await shopApi.updateShop(shopId.value, form.value)
    toast.success('Paramètres mis à jour avec succès !')
  } catch (e: any) {
    toast.error('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchShop()
})

definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>
