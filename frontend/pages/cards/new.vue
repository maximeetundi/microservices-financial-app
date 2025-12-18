<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-2xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="mb-8">
        <NuxtLink to="/cards" class="text-sm text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white mb-4 inline-flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
          Retour aux cartes
        </NuxtLink>
        <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white mb-2">Commander une carte 💳</h1>
        <p class="text-gray-500 dark:text-gray-400">Choisissez le type de carte qui vous convient</p>
      </div>

      <div class="bg-white dark:bg-slate-900 rounded-2xl p-8 shadow-xl border border-gray-100 dark:border-gray-800">
        <form @submit.prevent="createCard" class="space-y-6">
          
          <!-- Card Type Selection -->
          <div>
            <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-4">Type de carte</label>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div 
                class="border-2 rounded-xl p-4 cursor-pointer transition-all relative overflow-hidden group"
                :class="form.card_type === 'virtual' ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-900/10' : 'border-gray-200 dark:border-gray-700 hover:border-indigo-300'"
                @click="form.card_type = 'virtual'"
              >
                <div class="flex items-center gap-3 mb-2">
                  <span class="text-2xl">🌐</span>
                  <span class="font-bold text-gray-900 dark:text-white">Virtuelle</span>
                </div>
                <p class="text-xs text-gray-500 dark:text-gray-400">Pour les achats en ligne. Disponible immédiatement.</p>
                <div v-if="form.card_type === 'virtual'" class="absolute top-2 right-2 text-indigo-600">
                  <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                </div>
              </div>

              <div 
                class="border-2 rounded-xl p-4 cursor-pointer transition-all relative overflow-hidden group opacity-50"
                title="Bientôt disponible"
              >
                <div class="flex items-center gap-3 mb-2">
                  <span class="text-2xl">💳</span>
                  <span class="font-bold text-gray-900 dark:text-white">Physique</span>
                </div>
                <p class="text-xs text-gray-500 dark:text-gray-400">Pour les retraits guichet et paiements en magasin.</p>
                <span class="absolute bottom-2 right-2 text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded-full font-bold">Bientôt</span>
              </div>
            </div>
          </div>

          <!-- Currency -->
          <div>
            <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Devise</label>
            <select v-model="form.currency" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
              <option value="USD">🇺🇸 Dollar US (USD)</option>
              <option value="EUR">🇪🇺 Euro (EUR)</option>
              <option value="GBP">🇬🇧 Livre Sterling (GBP)</option>
            </select>
          </div>

          <!-- Cardholder Name -->
          <div>
            <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Nom sur la carte</label>
            <input 
              v-model="form.cardholder_name" 
              type="text" 
              placeholder="votre nom complet"
              class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
            />
          </div>

          <!-- Source Wallet (if needed for funding or context) -->
          <!-- Note: Backend might require source_wallet_id if there is an initial load fee or load amount, 
               but CreateCardRequest says InitialAmount is optional. Let's assume we can create empty for now or select a wallet. -->
          
          <div v-if="wallets.length > 0">
             <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Portefeuille de financement (source des fonds)</label>
             <select v-model="form.source_wallet_id" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all mb-4">
                <option value="">-- Sélectionner un portefeuille --</option>
                <option v-for="w in wallets" :key="w.id" :value="w.id">{{ w.name }} ({{ w.balance }} {{ w.currency }})</option>
             </select>

             <div v-if="form.source_wallet_id">
                 <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Montant initial à créditer</label>
                 <div class="relative">
                    <input 
                      v-model="form.initial_amount" 
                      type="number" 
                      min="0"
                      step="0.01"
                      placeholder="0.00"
                      class="input-premium w-full p-3 pl-10 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
                    />
                    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <span class="text-gray-500 sm:text-sm">$</span>
                    </div>
                 </div>
                 <p class="text-xs text-gray-500 mt-1">Ce montant sera débité de votre portefeuille sélectionné.</p>
             </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="p-4 rounded-xl bg-rose-50 text-rose-600 border border-rose-100 text-sm">
            {{ error }}
          </div>

          <!-- Submit -->
          <button 
            type="submit" 
            :disabled="loading"
            class="w-full py-4 px-6 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/30 flex justify-center items-center gap-2"
          >
            <span v-if="loading" class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
            {{ loading ? 'Création en cours...' : 'Créer ma carte maintenant' }}
          </button>

        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { cardAPI, walletAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')
const wallets = ref([])

const form = ref({
  card_type: 'virtual',
  card_category: 'personal', // Default
  currency: 'USD',
  cardholder_name: '', 
  source_wallet_id: '',
  initial_amount: 0
})

const fetchWallets = async () => {
    try {
        const res = await walletAPI.getAll()
        if (res.data?.wallets) {
             wallets.value = res.data.wallets.filter(w => w.wallet_type !== 'crypto')
        }
    } catch (e) {
        console.error("Error fetching wallets", e)
    }
}

const createCard = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const payload = {
        card_type: form.value.card_type,
        card_category: form.value.card_category,
        currency: form.value.currency,
        cardholder_name: form.value.cardholder_name,
        initial_amount: Number(form.value.initial_amount),
        source_wallet_id: form.value.source_wallet_id || undefined
    }
    
    await cardAPI.create(payload)

    router.push('/cards')
  } catch (e) {
    console.error(e)
    error.value = e.response?.data?.error || e.message || "Une erreur est survenue lors de la création de la carte."
  } finally {
    loading.value = false
  }
}

onMounted(() => {
    fetchWallets()
    if (authStore.user?.first_name && authStore.user?.last_name) {
        form.value.cardholder_name = `${authStore.user.first_name} ${authStore.user.last_name}`
    } else if (authStore.user?.name) {
        form.value.cardholder_name = authStore.user.name
    }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>
