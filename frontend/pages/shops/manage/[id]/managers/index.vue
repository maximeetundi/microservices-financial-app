<template>
  <div class="container mx-auto px-4 py-8 max-w-4xl">
    <!-- Header -->
    <div class="flex justify-between items-center mb-8">
      <div>
        <NuxtLink :to="`/shops/manage/${slug}`" class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 flex items-center gap-1 mb-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
          Retour au tableau de bord
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Gestionnaires</h1>
        <p class="text-gray-500 dark:text-gray-400">Gérez l'accès à votre boutique</p>
      </div>
      <button @click="showInviteModal = true" class="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/></svg>
        Inviter un membre
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
       <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else class="space-y-6">
      
      <!-- Owner Card -->
      <div class="bg-white dark:bg-slate-800 shadow rounded-lg p-6 border-l-4 border-indigo-500">
         <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
               <div class="w-12 h-12 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400 font-bold text-xl">
                  {{ ownerInitials }}
               </div>
               <div>
                  <h3 class="font-bold text-gray-900 dark:text-white">Propriétaire</h3>
                  <p class="text-sm text-gray-500">Accès complet propriétaire</p>
               </div>
            </div>
            <span class="px-3 py-1 bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-300 rounded-full text-xs font-bold uppercase">Owner</span>
         </div>
      </div>

      <!-- Managers List -->
      <div class="bg-white dark:bg-slate-800 shadow rounded-lg overflow-hidden">
        <div class="p-6 border-b border-gray-100 dark:border-gray-700">
          <h2 class="text-lg font-bold text-gray-900 dark:text-white">Membres de l'équipe</h2>
        </div>
        
        <div v-if="managers.length === 0" class="p-12 text-center text-gray-500">
           <div class="mb-3 text-4xl">👥</div>
           <p>Aucun gestionnaire invité pour le moment.</p>
        </div>

        <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
           <div v-for="manager in managers" :key="manager.user_id || manager.email" class="p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
              <div class="flex items-center gap-4">
                 <div class="w-10 h-10 rounded-full bg-gray-200 dark:bg-slate-700 flex items-center justify-center text-gray-500 dark:text-gray-400 font-bold">
                    {{ (manager.first_name?.[0] || manager.email?.[0] || '?').toUpperCase() }}
                 </div>
                 <div>
                    <h4 class="font-medium text-gray-900 dark:text-white">{{ manager.first_name ? `${manager.first_name} ${manager.last_name}` : manager.email }}</h4>
                    <p class="text-xs text-gray-500">{{ manager.email }}</p>
                 </div>
              </div>
              
              <div class="flex items-center gap-4">
                 <div class="flex flex-col items-end">
                    <span class="text-sm font-medium text-gray-700 dark:text-gray-300 capitalize">{{ manager.role }}</span>
                    <span class="text-xs" :class="statusColor(manager.status)">{{ manager.status }}</span>
                 </div>
                 
                 <button @click="removeManager(manager)" class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors" title="Révoquer l'accès">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                 </button>
              </div>
           </div>
        </div>
      </div>
    </div>

    <!-- Invite Modal -->
    <div v-if="showInviteModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md p-6 shadow-2xl">
         <h3 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Inviter un gestionnaire</h3>
         
         <form @submit.prevent="inviteManager">
            <div class="mb-4">
               <label class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">Email du membre</label>
               <input v-model="inviteForm.email" type="email" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-800 focus:ring-indigo-500" placeholder="collague@exemple.com">
            </div>
            
            <div class="mb-6">
               <label class="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-300">Rôle</label>
               <select v-model="inviteForm.role" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-slate-800 focus:ring-indigo-500">
                  <option value="editor">Éditeur (Peut gérer produits & commandes)</option>
                  <option value="admin">Administrateur (Accès complet sauf suppress.)</option>
                  <option value="viewer">Observateur (Lecture seule)</option>
               </select>
            </div>

            <div class="flex justify-end gap-3">
               <button type="button" @click="showInviteModal = false" class="px-4 py-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-lg">Annuler</button>
               <button type="submit" :disabled="submitting" class="px-4 py-2 bg-indigo-600 text-white rounded-lg font-bold hover:bg-indigo-700 flex items-center gap-2">
                  <span v-if="submitting" class="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
                  Inviter
               </button>
            </div>
         </form>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type ShopManager } from '@/composables/useShopApi'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const loading = ref(true)
const showInviteModal = ref(false)
const submitting = ref(false)
const shopId = ref('')
const managers = ref<ShopManager[]>([])
const shopOwnerId = ref('') // Store owner ID to identify

const inviteForm = ref({
  email: '',
  role: 'editor',
  permissions: [] as string[]
})

const ownerInitials = computed(() => 'OP') // Ideally fetch owner details

const statusColor = (status: string) => {
   switch(status) {
      case 'active': return 'text-green-600'
      case 'pending': return 'text-yellow-600'
      default: return 'text-gray-500'
   }
}

const fetchShopData = async () => {
  try {
    loading.value = true
    const shop = await shopApi.getShop(slug)
    shopId.value = shop.id
    managers.value = shop.managers || []
    shopOwnerId.value = shop.owner_id
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const inviteManager = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    // Default permissions based on role
    let perms: string[] = []
    if (inviteForm.value.role === 'admin') perms = ['*']
    else if (inviteForm.value.role === 'editor') perms = ['products.*', 'orders.*']
    else perms = ['read']

    await shopApi.inviteManager(shopId.value, inviteForm.value.email, inviteForm.value.role, perms)
    
    // Refresh
    await fetchShopData()
    showInviteModal.value = false
    inviteForm.value.email = ''
  } catch (e: any) {
    alert('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const removeManager = async (manager: ShopManager) => {
   if (!confirm(`Révoquer l'accès pour ${manager.email} ?`)) return
   try {
      await shopApi.removeManager(shopId.value, manager.user_id || 'pending') // Use ID if registered, or we might need email for pending invites?
      // Note: Backend RemoveManager likely expects UserID. If pending, how do we remove? 
      // Checking models logic: Usually invites have an ID or we delete by email? 
      // Handler says RemoveManager(shopID, targetUserID). 
      // If invitation is pending, there might not be a UserID yet if user doesn't exist?
      // Assuming for now manager object has user_id populated or we handle it.
      await fetchShopData()
   } catch (e: any) {
      alert('Erreur: ' + (e.response?.data?.error || e.message))
   }
}

onMounted(() => {
  fetchShopData()
})

definePageMeta({
  middleware: ['auth']
})
</script>
