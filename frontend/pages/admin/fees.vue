<template>
  <NuxtLayout name="admin">
    <div class="p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">Gestion des Frais</h1>
        <p class="text-slate-400">Configurez les frais dynamiques pour les services Wallet et Exchange.</p>
      </div>

      <!-- Tabs -->
      <div class="flex space-x-4 mb-6 border-b border-slate-700">
        <button 
          @click="activeTab = 'wallet'"
          class="pb-2 px-4 text-sm font-medium transition-colors relative"
          :class="activeTab === 'wallet' ? 'text-indigo-400 border-b-2 border-indigo-400' : 'text-slate-400 hover:text-white'"
        >
          Wallet Service
        </button>
        <button 
          @click="activeTab = 'exchange'"
          class="pb-2 px-4 text-sm font-medium transition-colors relative"
          :class="activeTab === 'exchange' ? 'text-indigo-400 border-b-2 border-indigo-400' : 'text-slate-400 hover:text-white'"
        >
          Exchange Service
        </button>
      </div>

      <!-- Content -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
      </div>

      <div v-else class="bg-slate-800/50 backdrop-blur-xl rounded-2xl border border-slate-700/50 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-900/50 text-slate-400 text-xs uppercase font-medium">
              <tr>
                <th class="px-6 py-4">Clé / Description</th>
                <th class="px-6 py-4">Type</th>
                <th class="px-6 py-4 text-right">Valeur</th>
                <th class="px-6 py-4 text-right">Min / Max</th>
                <th class="px-6 py-4 text-center">Statut</th>
                <th class="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr v-for="fee in currentFess" :key="fee.id" class="hover:bg-slate-700/20 transition-colors">
                <td class="px-6 py-4">
                  <div class="text-white font-medium">{{ formatKey(fee.key) }}</div>
                  <div class="text-xs text-slate-400">{{ fee.description }}</div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="{
                      'bg-blue-500/10 text-blue-400': fee.type === 'percentage',
                      'bg-green-500/10 text-green-400': fee.type === 'fixed',
                      'bg-purple-500/10 text-purple-400': fee.type === 'hybrid'
                    }">
                    {{ fee.type }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right font-mono text-slate-300">
                  <div v-if="fee.type === 'percentage' || fee.type === 'hybrid'">
                    {{ fee.percentage }}%
                  </div>
                  <div v-if="fee.type === 'fixed' || fee.type === 'hybrid'" class="text-xs text-slate-400">
                    + {{ fee.fixed_amount }} {{ fee.currency }}
                  </div>
                </td>
                <td class="px-6 py-4 text-right font-mono text-xs text-slate-400">
                  <div>Min: {{ fee.min_fee }}</div>
                  <div>Max: {{ fee.max_fee > 0 ? fee.max_fee : '∞' }}</div>
                </td>
                <td class="px-6 py-4 text-center">
                   <div class="w-2 h-2 rounded-full mx-auto" :class="fee.is_active ? 'bg-emerald-400' : 'bg-red-400'"></div>
                </td>
                <td class="px-6 py-4 text-right">
                  <button @click="openEditModal(fee)" class="text-indigo-400 hover:text-indigo-300 font-medium text-sm">
                    Modifier
                  </button>
                </td>
              </tr>
              <tr v-if="currentFess.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-400">
                  Aucune configuration de frais trouvée.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/80 backdrop-blur-sm">
      <div class="bg-slate-800 rounded-2xl border border-slate-700 w-full max-w-lg shadow-xl" @click.stop>
        <div class="p-6 border-b border-slate-700 flex justify-between items-center">
          <h2 class="text-xl font-bold text-white">Modifier les frais</h2>
          <button @click="closeModal" class="text-slate-400 hover:text-white">✕</button>
        </div>
        
        <div class="p-6 space-y-4">
           <div>
            <label class="block text-sm font-medium text-slate-400 mb-1">Description</label>
            <input v-model="editingFee.description" type="text" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
               <label class="block text-sm font-medium text-slate-400 mb-1">Type</label>
               <select v-model="editingFee.type" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
                 <option value="percentage">Pourcentage</option>
                 <option value="fixed">Fixe</option>
                 <option value="hybrid">Hybride</option>
               </select>
            </div>
             <div>
               <label class="block text-sm font-medium text-slate-400 mb-1">Statut</label>
               <select v-model="editingFee.is_active" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
                 <option :value="true">Actif</option>
                 <option :value="false">Inactif</option>
               </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
             <div v-if="editingFee.type !== 'fixed'">
               <label class="block text-sm font-medium text-slate-400 mb-1">Pourcentage (%)</label>
               <input v-model.number="editingFee.percentage" type="number" step="0.01" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
             </div>
             <div v-if="editingFee.type !== 'percentage'">
               <label class="block text-sm font-medium text-slate-400 mb-1">Montant Fixe</label>
               <input v-model.number="editingFee.fixed_amount" type="number" step="0.01" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
             </div>
          </div>

           <div class="grid grid-cols-2 gap-4">
             <div>
               <label class="block text-sm font-medium text-slate-400 mb-1">Min Fee</label>
               <input v-model.number="editingFee.min_fee" type="number" step="0.01" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
             </div>
             <div>
               <label class="block text-sm font-medium text-slate-400 mb-1">Max Fee (0 = infini)</label>
               <input v-model.number="editingFee.max_fee" type="number" step="0.01" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500 transition-colors">
             </div>
          </div>

        </div>

        <div class="p-6 border-t border-slate-700 flex justify-end gap-3">
          <button @click="closeModal" class="px-4 py-2 rounded-xl text-slate-300 hover:text-white hover:bg-slate-700 transition-colors">Annuler</button>
          <button @click="saveFee" :disabled="saving" class="px-4 py-2 rounded-xl bg-indigo-500 hover:bg-indigo-600 text-white font-medium transition-colors flex items-center gap-2">
            <span v-if="saving" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></span>
            <span>Enregistrer</span>
          </button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { adminFeeApi } = useApi()

const activeTab = ref('wallet')
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const walletFees = ref([])
const exchangeFees = ref([])
const editingFee = ref({})

const currentFess = computed(() => {
  return activeTab.value === 'wallet' ? walletFees.value : exchangeFees.value
})

const fetchFees = async () => {
  loading.value = true
  try {
    const [walletRes, exchangeRes] = await Promise.all([
      adminFeeApi.getWalletFees(),
      adminFeeApi.getExchangeFees()
    ])
    walletFees.value = walletRes.data || []
    exchangeFees.value = exchangeRes.data || []
  } catch (error) {
    console.error("Failed to fetch fees", error)
    // Could add toast notification here
  } finally {
    loading.value = false
  }
}

const formatKey = (key) => {
  return key.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}

const openEditModal = (fee) => {
  editingFee.value = JSON.parse(JSON.stringify(fee)) // Deep copy
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingFee.value = {}
}

const saveFee = async () => {
  saving.value = true
  try {
    if (activeTab.value === 'wallet') {
      await adminFeeApi.updateWalletFee(editingFee.value)
      // Update local state vs re-fetch
       const index = walletFees.value.findIndex(f => f.id === editingFee.value.id)
       if (index !== -1) walletFees.value[index] = editingFee.value
    } else {
      await adminFeeApi.updateExchangeFee(editingFee.value)
       const index = exchangeFees.value.findIndex(f => f.id === editingFee.value.id)
       if (index !== -1) exchangeFees.value[index] = editingFee.value
    }
    closeModal()
  } catch (error) {
    console.error("Failed to update fee", error)
    alert("Erreur lors de la mise à jour : " + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchFees()
})
</script>
