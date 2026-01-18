<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-2xl mx-auto py-8 px-4">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-20">
        <div class="text-6xl mb-4">❌</div>
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Erreur</h2>
        <p class="text-gray-500">{{ error }}</p>
        <button @click="$router.back()" class="mt-4 px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg">
          Retour
        </button>
      </div>

      <!-- Approval Content -->
      <div v-else-if="approval" class="space-y-6">
        <!-- Header -->
        <div class="text-center mb-8">
          <div class="w-20 h-20 rounded-full bg-yellow-100 dark:bg-yellow-900/30 flex items-center justify-center mx-auto mb-4">
            <span class="text-4xl">⏳</span>
          </div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Confirmation Requise</h1>
          <p class="text-gray-500 mt-2">Cette action nécessite l'approbation d'un administrateur</p>
        </div>

        <!-- Status Badge -->
        <div class="flex justify-center mb-6">
          <span :class="statusClass" class="px-4 py-2 rounded-full text-sm font-bold">
            {{ statusLabel }}
          </span>
        </div>

        <!-- Action Details Card -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
          <h3 class="font-bold text-gray-900 dark:text-white mb-4">{{ approval.action_name }}</h3>
          
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Type</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ approval.action_type }}</span>
            </div>
            
            <div v-if="approval.amount" class="flex justify-between text-sm">
              <span class="text-gray-500">Montant</span>
              <span class="font-bold text-xl text-primary-600">{{ formatAmount(approval.amount, approval.currency) }}</span>
            </div>

            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Description</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ approval.description }}</span>
            </div>

            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Initié par</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ approval.initiator_name }}</span>
            </div>

            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Créé le</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatDate(approval.created_at) }}</span>
            </div>

            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Expire le</span>
              <span class="font-medium text-orange-500">{{ formatDate(approval.expires_at) }}</span>
            </div>
          </div>
        </div>

        <!-- Approvals Progress -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
          <h4 class="font-bold text-gray-900 dark:text-white mb-4">Approbations</h4>
          
          <div class="flex justify-between items-center mb-4">
            <span class="text-gray-500">Progression</span>
            <span class="font-bold text-primary-600">
              {{ approval.approvals?.filter(a => a.decision === 'APPROVED').length || 0 }} / {{ approval.required_approvals }}
            </span>
          </div>

          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 mb-4">
            <div 
              class="bg-primary-600 h-2 rounded-full transition-all" 
              :style="{ width: approvalProgress + '%' }"
            ></div>
          </div>

          <!-- Individual Votes -->
          <div v-if="approval.approvals?.length > 0" class="space-y-2 mt-4">
            <div v-for="vote in approval.approvals" :key="vote.admin_user_id" 
              class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
              <span class="font-medium text-gray-900 dark:text-white">{{ vote.admin_name }}</span>
              <span :class="vote.decision === 'APPROVED' ? 'text-green-600' : 'text-red-600'" class="font-bold">
                {{ vote.decision === 'APPROVED' ? '✓ Approuvé' : '✕ Rejeté' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Action Buttons (only show if can vote) -->
        <div v-if="canVote && approval.status === 'PENDING'" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <button 
              @click="showRejectModal = true"
              class="px-6 py-4 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded-xl font-bold hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors">
              ✕ Rejeter
            </button>
            <button 
              @click="showApproveModal = true"
              class="px-6 py-4 bg-green-600 text-white rounded-xl font-bold hover:bg-green-700 transition-colors">
              ✓ Approuver
            </button>
          </div>
        </div>

        <!-- Already Voted Message -->
        <div v-else-if="hasVoted" class="text-center py-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl">
          <span class="text-blue-600 dark:text-blue-400 font-medium">
            ℹ️ Vous avez déjà voté sur cette action
          </span>
        </div>

        <!-- Back Button -->
        <button @click="$router.back()" class="w-full py-3 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-xl transition-colors">
          ← Retour
        </button>
      </div>
    </div>

    <!-- PIN Modal for Approve -->
    <PinModal 
      :isOpen="showApproveModal"
      title="Confirmer l'approbation"
      description="Entrez votre code PIN pour approuver cette action."
      @update:isOpen="showApproveModal = $event"
      @success="handleApprove"
      @close="showApproveModal = false"
    />

    <!-- Reject Modal -->
    <Teleport to="body">
      <div v-if="showRejectModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Rejeter cette action?</h3>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Motif (optionnel)</label>
            <textarea v-model="rejectReason" rows="3" placeholder="Indiquez la raison du rejet..."
              class="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white"></textarea>
          </div>

          <div class="flex gap-3">
            <button @click="showRejectModal = false" class="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-lg">
              Annuler
            </button>
            <button @click="showPinForReject = true; showRejectModal = false" class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg">
              Confirmer
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- PIN Modal for Reject -->
    <PinModal 
      :isOpen="showPinForReject"
      title="Confirmer le rejet"
      description="Entrez votre code PIN pour rejeter cette action."
      @update:isOpen="showPinForReject = $event"
      @success="handleReject"
      @close="showPinForReject = false"
    />
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { enterpriseAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'
import PinModal from '~/components/common/PinModal.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const user = computed(() => authStore.user)

const approvalId = computed(() => route.params.id as string)
const approval = ref<any>(null)
const loading = ref(true)
const error = ref('')

// Modals
const showApproveModal = ref(false)
const showRejectModal = ref(false)
const showPinForReject = ref(false)
const rejectReason = ref('')

// Computed
const canVote = computed(() => {
  if (!approval.value || !user.value) return false
  // Can vote if admin and hasn't voted yet
  return !approval.value.approvals?.some((a: any) => a.admin_user_id === user.value.id)
})

const hasVoted = computed(() => {
  if (!approval.value || !user.value) return false
  return approval.value.approvals?.some((a: any) => a.admin_user_id === user.value.id)
})

const approvalProgress = computed(() => {
  if (!approval.value) return 0
  const approved = approval.value.approvals?.filter((a: any) => a.decision === 'APPROVED').length || 0
  return Math.min(100, (approved / approval.value.required_approvals) * 100)
})

const statusClass = computed(() => {
  if (!approval.value) return 'bg-gray-100 text-gray-600'
  switch (approval.value.status) {
    case 'PENDING': return 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400'
    case 'APPROVED': return 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400'
    case 'REJECTED': return 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
    case 'EXECUTED': return 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400'
    case 'EXPIRED': return 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
    default: return 'bg-gray-100 text-gray-600'
  }
})

const statusLabel = computed(() => {
  if (!approval.value) return ''
  switch (approval.value.status) {
    case 'PENDING': return '⏳ En attente d\'approbation'
    case 'APPROVED': return '✓ Approuvé'
    case 'REJECTED': return '✕ Rejeté'
    case 'EXECUTED': return '✅ Exécuté'
    case 'EXPIRED': return '⏰ Expiré'
    default: return approval.value.status
  }
})

const formatAmount = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency || 'XOF',
    minimumFractionDigits: 0
  }).format(amount || 0)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadApproval = async () => {
  loading.value = true
  error.value = ''
  try {
    // We need to get the approval by ID - this might need an API endpoint
    // For now, we'll try to get pending approvals from all enterprises
    // A proper implementation would have a direct GET /approvals/:id endpoint
    const { data } = await enterpriseAPI.getApprovalById(approvalId.value)
    approval.value = data
  } catch (e: any) {
    console.error('Failed to load approval:', e)
    error.value = e.response?.data?.error || 'Impossible de charger cette approbation'
  } finally {
    loading.value = false
  }
}

const handleApprove = async (encryptedPin: string) => {
  showApproveModal.value = false
  loading.value = true
  try {
    await enterpriseAPI.approveAction(approvalId.value, { pin: encryptedPin })
    alert('Action approuvée avec succès!')
    await loadApproval()
  } catch (e: any) {
    alert(e.response?.data?.error || 'Erreur lors de l\'approbation')
  } finally {
    loading.value = false
  }
}

const handleReject = async (encryptedPin: string) => {
  showPinForReject.value = false
  loading.value = true
  try {
    await enterpriseAPI.rejectAction(approvalId.value, { 
      pin: encryptedPin, 
      reason: rejectReason.value 
    })
    alert('Action rejetée.')
    await loadApproval()
  } catch (e: any) {
    alert(e.response?.data?.error || 'Erreur lors du rejet')
  } finally {
    loading.value = false
    rejectReason.value = ''
  }
}

onMounted(loadApproval)

definePageMeta({
  middleware: 'auth'
})
</script>
