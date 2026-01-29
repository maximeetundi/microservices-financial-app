<template>
  <NuxtLayout name="admin">
    <div class="p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">Vérification KYC</h1>
        <p class="text-slate-400">Queue de validation des documents d'identité.</p>
      </div>

       <!-- Tabs -->
      <div class="flex space-x-4 mb-6 border-b border-slate-700">
        <button 
          v-for="status in ['pending', 'approved', 'rejected', 'all']"
          :key="status"
          @click="activeStatus = status; offset = 0; fetchRequests()"
          class="pb-2 px-4 text-sm font-medium transition-colors relative capitalize"
          :class="activeStatus === status ? 'text-indigo-400 border-b-2 border-indigo-400' : 'text-slate-400 hover:text-white'"
        >
          {{ status === 'pending' ? 'En attente' : status === 'approved' ? 'Validés' : status === 'rejected' ? 'Rejetés' : 'Tous' }}
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
                <th class="px-6 py-4">Utilisateur</th>
                <th class="px-6 py-4">Type Document</th>
                <th class="px-6 py-4">Fichier</th>
                <th class="px-6 py-4">Date de soumission</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr v-for="req in requests" :key="req.document.id" class="hover:bg-slate-700/20 transition-colors">
                <td class="px-6 py-4">
                  <div class="font-medium text-white">{{ req.user_name }}</div>
                  <div class="text-xs text-slate-400">{{ req.user_email }}</div>
                </td>
                <td class="px-6 py-4">
                  <div class="capitalize text-slate-300">{{ req.document.type }}</div>
                  <div v-if="req.document.identity_sub_type" class="text-xs text-slate-500 capitalize">
                    {{ req.document.identity_sub_type }}
                  </div>
                </td>
                <td class="px-6 py-4 text-sm text-indigo-400">
                  <a :href="req.document.file_path" target="_blank" class="hover:underline flex items-center gap-1">
                    📄 {{ req.document.file_name }}
                  </a>
                </td>
                <td class="px-6 py-4 text-xs text-slate-400">
                  {{ formatDate(req.document.uploaded_at) }}
                </td>
                <td class="px-6 py-4">
                   <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="{
                      'bg-emerald-500/10 text-emerald-400': req.document.status === 'approved',
                      'bg-red-500/10 text-red-400': req.document.status === 'rejected',
                      'bg-amber-500/10 text-amber-400': req.document.status === 'pending'
                    }">
                    {{ req.document.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right">
                  <button 
                    v-if="req.document.status === 'pending'"
                    @click="openReviewModal(req)" 
                    class="bg-indigo-500 hover:bg-indigo-600 text-white px-3 py-1 rounded-lg text-sm font-medium transition-colors"
                  >
                    Examiner
                  </button>
                  <span v-else class="text-xs text-slate-500">
                     {{ req.document.status === 'approved' ? 'Validé' : 'Rejeté' }}
                     <span v-if="req.document.reviewed_at">le {{ formatDate(req.document.reviewed_at) }}</span>
                  </span>
                </td>
              </tr>
              <tr v-if="requests.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-400">
                  Aucun document à vérifier.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="px-6 py-4 border-t border-slate-700/50 flex justify-between items-center">
          <button @click="prevPage" :disabled="offset === 0" class="text-slate-400 hover:text-white disabled:opacity-50">
            ← Précédent
          </button>
           <span class="text-slate-500 text-sm">Total: {{ total }}</span>
          <button @click="nextPage" :disabled="offset + limit >= total" class="text-slate-400 hover:text-white disabled:opacity-50">
            Suivant →
          </button>
        </div>
      </div>
    </div>

    <!-- Review Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/90 backdrop-blur-sm">
      <div class="bg-slate-800 rounded-2xl border border-slate-700 w-full max-w-4xl shadow-xl max-h-[90vh] overflow-y-auto" @click.stop>
        <div class="p-6 border-b border-slate-700 flex justify-between items-center sticky top-0 bg-slate-800 z-10">
          <div>
            <h2 class="text-xl font-bold text-white">Examen du document</h2>
            <p class="text-sm text-slate-400">{{ selectedReq?.user_name }} - {{ selectedReq?.document.type }}</p>
          </div>
          <button @click="closeModal" class="text-slate-400 hover:text-white">✕</button>
        </div>
        
        <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-8">
           <!-- Document Preview -->
           <div class="bg-slate-900 rounded-xl p-2 border border-slate-700 flex items-center justify-center min-h-[400px]">
             <img 
               v-if="isImage(selectedReq?.document.mime_type)" 
               :src="selectedReq?.document.file_path" 
               class="max-w-full max-h-[500px] object-contain"
             />
             <iframe 
               v-else-if="selectedReq?.document.mime_type === 'application/pdf'" 
               :src="selectedReq?.document.file_path" 
               class="w-full h-[500px]"
             ></iframe>
             <div v-else class="text-center p-8">
               <p class="text-slate-400 mb-4">Aperçu non disponible pour ce format</p>
               <a :href="selectedReq?.document.file_path" target="_blank" class="text-indigo-400 underline">Télécharger le fichier</a>
             </div>
           </div>

           <!-- Actions -->
           <div class="space-y-6">
             <div class="bg-slate-700/30 rounded-xl p-4 border border-slate-700">
               <h3 class="font-medium text-white mb-4">Détails détectés</h3>
               <dl class="space-y-2 text-sm">
                 <div class="flex justify-between">
                   <dt class="text-slate-400">Numéro document (OCR)</dt>
                   <dd class="text-white font-mono">{{ selectedReq?.document.document_number || 'Non détecté' }}</dd>
                 </div>
                 <div class="flex justify-between">
                   <dt class="text-slate-400">Date expiration</dt>
                   <dd class="text-white font-mono">{{ selectedReq?.document.expiry_date || 'Non détectée' }}</dd>
                 </div>
               </dl>
             </div>

             <div v-if="reviewAction === 'reject'" class="space-y-2">
               <label class="text-sm font-medium text-slate-300">Raison du rejet</label>
               <textarea 
                 v-model="rejectionReason"
                 class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-red-500 h-24"
                 placeholder="Expliquez pourquoi le document est rejeté..."
               ></textarea>
             </div>

             <div class="flex flex-col gap-3">
               <div v-if="reviewAction !== 'reject'" class="grid grid-cols-2 gap-3">
                 <button 
                   @click="reviewAction = 'reject'" 
                   class="px-4 py-3 rounded-xl border border-red-500/30 text-red-400 hover:bg-red-500/10 transition-colors"
                 >
                   Refuser
                 </button>
                 <button 
                   @click="submitReview('approved')" 
                   :disabled="submitting"
                   class="px-4 py-3 rounded-xl bg-emerald-500 hover:bg-emerald-600 text-white font-medium transition-colors"
                 >
                   <span v-if="submitting">Traitement...</span>
                   <span v-else>✅ Valider le document</span>
                 </button>
               </div>

               <div v-else class="grid grid-cols-2 gap-3">
                 <button 
                   @click="reviewAction = null" 
                   class="px-4 py-3 rounded-xl border border-slate-600 text-slate-300 hover:bg-slate-700 transition-colors"
                 >
                   Annuler
                 </button>
                 <button 
                   @click="submitReview('rejected')" 
                   :disabled="submitting || !rejectionReason"
                   class="px-4 py-3 rounded-xl bg-red-500 hover:bg-red-600 text-white font-medium transition-colors"
                 >
                   <span v-if="submitting">Traitement...</span>
                   <span v-else>🚫 Confirmer le rejet</span>
                 </button>
               </div>
             </div>
           </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { adminKYCAPI } = useApi()

const loading = ref(true)
const requests = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)
const activeStatus = ref('pending')

const showModal = ref(false)
const selectedReq = ref(null)
const reviewAction = ref(null)
const rejectionReason = ref('')
const submitting = ref(false)

const fetchRequests = async () => {
  loading.value = true
  try {
    const response = await adminKYCAPI.getRequests(activeStatus.value, limit.value, offset.value)
    requests.value = response.data?.requests || []
    total.value = response.data?.total || 0
  } catch (error) {
    console.error("Failed to fetch KYC request", error)
  } finally {
    loading.value = false
  }
}

const nextPage = () => {
  if (offset.value + limit.value < total.value) {
    offset.value += limit.value
    fetchRequests()
  }
}

const prevPage = () => {
  if (offset.value >= limit.value) {
    offset.value -= limit.value
    fetchRequests()
  }
}

const openReviewModal = (req) => {
  selectedReq.value = req
  reviewAction.value = null
  rejectionReason.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedReq.value = null
}

const submitReview = async (status) => {
  if (!selectedReq.value) return
  
  submitting.value = true
  try {
    await adminKYCAPI.reviewRequest(selectedReq.value.document.id, status, rejectionReason.value)
    
    // Remove from list if viewing pending, or update status
    if (activeStatus.value === 'pending') {
      requests.value = requests.value.filter(r => r.document.id !== selectedReq.value.document.id)
      total.value--
    } else {
      fetchRequests()
    }
    
    closeModal()
  } catch (error) {
    console.error("Failed to review document", error)
    alert("Erreur: " + (error.response?.data?.error || error.message))
  } finally {
    submitting.value = false
  }
}

const isImage = (mimeType) => {
  return mimeType && mimeType.startsWith('image/')
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('fr-FR')
}

onMounted(() => {
  fetchRequests()
})
</script>
