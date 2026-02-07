<template>
  <div class="px-4 py-8 w-full">
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

      <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 border-b border-gray-100 dark:border-gray-700 pb-2 flex items-center gap-2">
          <PhotoIcon class="h-6 w-6 text-indigo-500" />
          Images
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Logo</label>
            <div class="aspect-square border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl flex items-center justify-center cursor-pointer hover:border-indigo-500 transition-colors">
              <input type="file" @change="uploadLogo" accept="image/*" class="hidden" ref="logoInput">
              <div v-if="form.logo_url" class="w-full h-full">
                <img :src="form.logo_url" class="w-full h-full object-cover rounded-xl">
              </div>
              <div v-else @click="($refs.logoInput as HTMLInputElement).click()" class="text-center p-4">
                <div class="text-4xl mb-2">📷</div>
                <p class="text-sm text-gray-500">Cliquer pour ajouter</p>
              </div>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Bannière</label>
            <div class="aspect-video border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl flex items-center justify-center cursor-pointer hover:border-indigo-500 transition-colors">
              <input type="file" @change="uploadBanner" accept="image/*" class="hidden" ref="bannerInput">
              <div v-if="form.banner_url" class="w-full h-full">
                <img :src="form.banner_url" class="w-full h-full object-cover rounded-xl">
              </div>
              <div v-else @click="($refs.bannerInput as HTMLInputElement).click()" class="text-center p-4">
                <div class="text-4xl mb-2">🖼️</div>
                <p class="text-sm text-gray-500">Cliquer pour ajouter</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-slate-800 shadow-sm ring-1 ring-gray-900/5 rounded-xl p-6">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 border-b border-gray-100 dark:border-gray-700 pb-2 flex items-center gap-2">
          <span class="text-indigo-500">🧩</span>
          Cartes de confiance
        </h2>

        <div class="space-y-4">
          <div
            v-for="(badge, idx) in trustBadges"
            :key="badge.key + '-' + idx"
            class="p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-900/30"
          >
            <div class="flex flex-col lg:flex-row lg:items-center gap-4">
              <div class="flex items-center gap-3">
                <input v-model="badge.enabled" type="checkbox" class="h-4 w-4 text-indigo-600 border-gray-300 rounded" />
                <input v-model="badge.icon" type="text" class="w-16 rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700" placeholder="🚚" />
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-3 flex-1">
                <input v-model="badge.title" type="text" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700" placeholder="Titre" />
                <input v-model="badge.subtitle" type="text" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-700" placeholder="Sous-titre" />
              </div>

              <div class="flex items-center gap-2 justify-end">
                <button type="button" @click="moveBadge(idx, -1)" :disabled="idx === 0" class="px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 text-sm font-bold disabled:opacity-50">
                  ↑
                </button>
                <button type="button" @click="moveBadge(idx, 1)" :disabled="idx === trustBadges.length - 1" class="px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 text-sm font-bold disabled:opacity-50">
                  ↓
                </button>
                <button type="button" @click="removeBadge(idx)" class="px-3 py-2 rounded-lg border border-red-300 text-red-600 text-sm font-bold hover:bg-red-50 dark:hover:bg-red-900/20">
                  Supprimer
                </button>
              </div>
            </div>
          </div>

          <div class="flex justify-end">
            <button type="button" @click="addBadge" class="px-5 py-2.5 bg-indigo-50 text-indigo-600 rounded-xl font-bold hover:bg-indigo-100 transition-colors">
              + Ajouter une carte
            </button>
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
  TruckIcon,
  PhotoIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const toast = useToast()
const slug = route.params.id as string

const loading = ref(true)
const saving = ref(false)
const shopId = ref('')

const trustBadges = ref<any[]>([])

const form = ref({
  name: '',
  description: '',
  is_public: true,
  logo_url: '',
  banner_url: '',
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
    form.value.logo_url = shop.logo_url || ''
    form.value.banner_url = shop.banner_url || ''
    
    if (shop.contact_info) form.value.contact_info = { ...shop.contact_info }
    if (shop.address) form.value.address = { ...shop.address }
    if (shop.settings) form.value.settings = { ...shop.settings }

    const defaults = [
      { key: 'fast_delivery', icon: '🚚', title: 'Livraison rapide', subtitle: 'Partout au Sénégal', enabled: true, order: 1 },
      { key: 'secure_payment', icon: '🔒', title: 'Paiement sécurisé', subtitle: '100% sécurisé', enabled: true, order: 2 },
      { key: 'quality_guarantee', icon: '⭐', title: 'Qualité garantie', subtitle: 'Produits vérifiés', enabled: true, order: 3 },
      { key: 'support_24_7', icon: '💬', title: 'Support 24/7', subtitle: 'À votre écoute', enabled: true, order: 4 },
    ]
    trustBadges.value = (shop.trust_badges && shop.trust_badges.length) ? [...shop.trust_badges] : defaults
    trustBadges.value = trustBadges.value.sort((a: any, b: any) => (a.order || 0) - (b.order || 0))

  } catch (e) {
    console.error('Failed to load shop', e)
  } finally {
    loading.value = false
  }
}

const uploadLogo = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const result = await shopApi.uploadMedia(file)
    form.value.logo_url = result.url
  } catch (e: any) {
    toast.error('Échec de l\'upload du logo. Veuillez réessayer.')
  }
}

const uploadBanner = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const result = await shopApi.uploadMedia(file)
    form.value.banner_url = result.url
  } catch (e: any) {
    toast.error('Échec de l\'upload de la bannière. Veuillez réessayer.')
  }
}

const reindexOrders = () => {
  trustBadges.value = trustBadges.value.map((b: any, i: number) => ({ ...b, order: i + 1 }))
}

const addBadge = () => {
  trustBadges.value.push({
    key: `custom_${Date.now()}`,
    icon: '✨',
    title: 'Nouveau',
    subtitle: 'Description',
    enabled: true,
    order: trustBadges.value.length + 1
  })
}

const removeBadge = (idx: number) => {
  trustBadges.value.splice(idx, 1)
  reindexOrders()
}

const moveBadge = (idx: number, dir: number) => {
  const newIdx = idx + dir
  if (newIdx < 0 || newIdx >= trustBadges.value.length) return
  const copy = [...trustBadges.value]
  const tmp = copy[idx]
  copy[idx] = copy[newIdx]
  copy[newIdx] = tmp
  trustBadges.value = copy
  reindexOrders()
}

const saveSettings = async () => {
  if (!shopId.value) return
  
  try {
    saving.value = true
    await shopApi.updateShop(shopId.value, form.value)
    reindexOrders()
    await shopApi.updateTrustBadges(shopId.value, trustBadges.value)
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
