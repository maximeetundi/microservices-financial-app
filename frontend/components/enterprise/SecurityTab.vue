<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <ShieldCheckIcon class="w-6 h-6 text-primary-500" />
            Sécurité & Approbations
          </h2>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            Configurez les règles de validation multi-signatures pour les actions sensibles
          </p>
        </div>
        
        <button @click="addPolicy" 
          class="px-5 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl font-medium transition-all shadow-lg shadow-primary-500/25 flex items-center gap-2">
          <PlusIcon class="w-5 h-5" />
          Nouvelle Règle
        </button>
      </div>
    </div>

    <!-- Info Banner -->
    <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 flex gap-3">
      <InformationCircleIcon class="w-6 h-6 text-blue-500 flex-shrink-0" />
      <div class="text-sm text-blue-800 dark:text-blue-200">
        <p class="font-medium">Comment ça fonctionne ?</p>
        <p class="mt-1">Quand une action sensible est initiée, tous les administrateurs désignés reçoivent une notification. 
        L'action sera exécutée uniquement après le nombre requis d'approbations (avec PIN/mot de passe).</p>
      </div>
    </div>

    <!-- Pending Approvals -->
    <div v-if="pendingApprovals.length" class="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-2xl p-6">
      <h3 class="text-lg font-bold text-amber-800 dark:text-amber-200 mb-4 flex items-center gap-2">
        <ClockIcon class="w-5 h-5" />
        Approbations en attente ({{ pendingApprovals.length }})
      </h3>
      
      <div class="space-y-3">
        <div v-for="approval in pendingApprovals" :key="approval.id"
          class="bg-white dark:bg-gray-800 rounded-xl p-4 border border-amber-100 dark:border-amber-800">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="font-medium text-gray-900 dark:text-white">{{ approval.action_name }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ approval.description }}</p>
              <div class="flex items-center gap-4 mt-2 text-xs text-gray-500">
                <span>Initié par: {{ approval.initiator_name }}</span>
                <span>{{ formatDate(approval.created_at) }}</span>
                <span class="font-medium text-primary-600">{{ approval.approvals?.length || 0 }}/{{ approval.required_approvals }} approbations</span>
              </div>
            </div>
            <div class="flex gap-2">
              <button @click="openApprovalModal(approval, 'approve')"
                class="px-3 py-1.5 bg-green-100 text-green-700 hover:bg-green-200 rounded-lg text-sm font-medium transition-colors">
                Approuver
              </button>
              <button @click="openApprovalModal(approval, 'reject')"
                class="px-3 py-1.5 bg-red-100 text-red-700 hover:bg-red-200 rounded-lg text-sm font-medium transition-colors">
                Rejeter
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Security Policies -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white">Règles de sécurité</h3>
      </div>

      <div v-if="!policies.length" class="p-12 text-center">
        <ShieldExclamationIcon class="w-16 h-16 mx-auto text-gray-300 dark:text-gray-600 mb-4" />
        <p class="text-gray-500 dark:text-gray-400">Aucune règle de sécurité configurée</p>
        <p class="text-sm text-gray-400 mt-1">Ajoutez des règles pour exiger des validations multi-signatures</p>
      </div>

      <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
        <div v-for="(policy, idx) in policies" :key="policy.id"
          class="p-6 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
          <div class="flex items-start justify-between gap-4">
            <!-- Policy Info -->
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <input v-model="policy.name" 
                  class="text-lg font-semibold text-gray-900 dark:text-white bg-transparent border-0 focus:ring-0 p-0 flex-1" 
                  placeholder="Nom de la règle" />
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="policy.enabled" class="sr-only peer" />
                  <div class="w-11 h-6 bg-gray-200 peer-focus:ring-4 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary-600"></div>
                </label>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
                <!-- Action Type -->
                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Type d'action</label>
                  <select v-model="policy.action_type" 
                    class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm">
                    <option value="TRANSACTION">Transactions</option>
                    <option value="PAYROLL">Paiement salaires</option>
                    <option value="EMPLOYEE_TERMINATE">Licenciement</option>
                    <option value="EMPLOYEE_PROMOTE">Promotion</option>
                    <option value="SETTINGS_CHANGE">Paramètres</option>
                    <option value="WALLET_CREATE">Création wallet</option>
                    <option value="ADMIN_CHANGE">Changement admin</option>
                    <option value="SERVICE_CREATE">Création service</option>
                    <option value="INVOICE_BATCH">Factures en lot</option>
                    <option value="ENTERPRISE_DELETE">⚠️ Suppression entreprise</option>
                  </select>
                </div>

                <!-- Approval Mode -->
                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Mode d'approbation</label>
                  <select v-model="policy.approval_mode" @change="updateApprovalMode(policy)"
                    class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm">
                    <option value="all">Tous les admins</option>
                    <option value="majority">Majorité</option>
                    <option value="count">Nombre précis</option>
                  </select>
                </div>

                <!-- Min Approvals (if count mode) -->
                <div v-if="policy.approval_mode === 'count'">
                  <label class="block text-xs font-medium text-gray-500 mb-1">Nombre minimum</label>
                  <input v-model.number="policy.min_approvals" type="number" min="1"
                    class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm" />
                </div>

                <!-- Threshold Amount -->
                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Seuil montant (optionnel)</label>
                  <input v-model.number="policy.threshold_amount" type="number" placeholder="0 = toujours"
                    class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm" />
                </div>

                <!-- Expiration -->
                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Délai d'expiration (heures)</label>
                  <input v-model.number="policy.expiration_hours" type="number" min="1"
                    class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm" />
                </div>
              </div>
            </div>

            <!-- Delete -->
            <button @click="removePolicy(idx)" class="p-2 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg">
              <TrashIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Save Button -->
    <div class="flex justify-end">
      <button @click="savePolicies" :disabled="isSaving"
        class="px-6 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl font-medium hover:from-primary-700 hover:to-primary-800 disabled:opacity-50 flex items-center gap-2 shadow-lg shadow-primary-500/25 transition-all">
        <span v-if="isSaving" class="animate-spin">⟳</span>
        {{ isSaving ? 'Enregistrement...' : 'Enregistrer les règles' }}
      </button>
    </div>

    <!-- Rejection Reason Modal (shown before PIN) -->
    <Teleport to="body">
      <div v-if="showRejectReasonModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2">Rejeter l'action</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-4">{{ selectedApproval?.action_name }}</p>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Raison du rejet (optionnel)</label>
            <textarea v-model="rejectReason" rows="3" placeholder="Expliquez pourquoi vous rejetez cette action..."
              class="w-full px-4 py-2 rounded-xl border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white"></textarea>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button @click="showRejectReasonModal = false" 
              class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
              Annuler
            </button>
            <button @click="proceedToPin"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700">
              Continuer
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- PIN Modal (reusable component) -->
    <PinModal 
      v-model:is-open="showPinModal"
      :title="approvalAction === 'approve' ? 'Approuver l\'action' : 'Rejeter l\'action'"
      :description="selectedApproval?.action_name || 'Confirmez votre identité'"
      @success="handlePinSuccess"
      @close="handlePinClose"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { 
  ShieldCheckIcon, PlusIcon, TrashIcon, ClockIcon, 
  InformationCircleIcon, ShieldExclamationIcon 
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import PinModal from '@/components/common/PinModal.vue'

const props = defineProps({
  enterprise: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

const policies = ref([])
const pendingApprovals = ref([])
const isSaving = ref(false)
const isLoading = ref(false)

// Approval modal state
const showPinModal = ref(false)
const showRejectReasonModal = ref(false)
const selectedApproval = ref(null)
const approvalAction = ref('approve')
const rejectReason = ref('')
const isApproving = ref(false)

// Initialize policies from enterprise
watch(() => props.enterprise, (ent) => {
  if (ent?.security_policies) {
    policies.value = JSON.parse(JSON.stringify(ent.security_policies)).map(p => ({
      ...p,
      approval_mode: p.require_all_admins ? 'all' : (p.require_majority ? 'majority' : 'count')
    }))
  }
}, { immediate: true })

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' 
  })
}

const addPolicy = () => {
  policies.value.push({
    id: `policy_${Date.now()}`,
    name: 'Nouvelle règle',
    action_type: 'TRANSACTION',
    enabled: true,
    approval_mode: 'majority',
    min_approvals: 2,
    require_majority: true,
    require_all_admins: false,
    threshold_amount: 0,
    expiration_hours: 24
  })
}

const removePolicy = (idx) => {
  policies.value.splice(idx, 1)
}

const updateApprovalMode = (policy) => {
  policy.require_all_admins = policy.approval_mode === 'all'
  policy.require_majority = policy.approval_mode === 'majority'
  if (policy.approval_mode === 'count' && !policy.min_approvals) {
    policy.min_approvals = 2
  }
}

const savePolicies = async () => {
  isSaving.value = true
  try {
    // Convert approval_mode back to backend format
    const formattedPolicies = policies.value.map(p => ({
      id: p.id,
      name: p.name,
      action_type: p.action_type,
      enabled: p.enabled,
      min_approvals: p.approval_mode === 'count' ? p.min_approvals : 0,
      require_majority: p.approval_mode === 'majority',
      require_all_admins: p.approval_mode === 'all',
      threshold_amount: p.threshold_amount || 0,
      expiration_hours: p.expiration_hours || 24
    }))

    const updated = {
      ...props.enterprise,
      security_policies: formattedPolicies
    }
    await enterpriseAPI.update(props.enterprise.id, updated)
    emit('update', updated)
    alert('✅ Règles de sécurité enregistrées')
  } catch (error) {
    console.error('Failed to save policies:', error)
    alert('Erreur lors de l\'enregistrement')
  } finally {
    isSaving.value = false
  }
}

const loadPendingApprovals = async () => {
  try {
    const { data } = await enterpriseAPI.getPendingApprovals(props.enterprise.id)
    pendingApprovals.value = data || []
  } catch (error) {
    console.error('Failed to load pending approvals:', error)
    pendingApprovals.value = []
  }
}

// Open approval flow
const openApprovalModal = (approval, action) => {
  selectedApproval.value = approval
  approvalAction.value = action
  rejectReason.value = ''
  
  if (action === 'reject') {
    // Show rejection reason modal first
    showRejectReasonModal.value = true
  } else {
    // Go straight to PIN
    showPinModal.value = true
  }
}

// After entering rejection reason, proceed to PIN
const proceedToPin = () => {
  showRejectReasonModal.value = false
  showPinModal.value = true
}

// Handle PIN verification success
const handlePinSuccess = async (encryptedPin) => {
  showPinModal.value = false
  isApproving.value = true
  
  try {
    if (approvalAction.value === 'approve') {
      await enterpriseAPI.approveAction(selectedApproval.value.id, { 
        pin: encryptedPin 
      })
    } else {
      await enterpriseAPI.rejectAction(selectedApproval.value.id, { 
        pin: encryptedPin,
        reason: rejectReason.value
      })
    }
    
    alert(approvalAction.value === 'approve' ? '✅ Action approuvée' : '❌ Action rejetée')
    await loadPendingApprovals()
  } catch (error) {
    console.error('Failed to process approval:', error)
    alert('Erreur: ' + (error.response?.data?.error || error.message))
  } finally {
    isApproving.value = false
    selectedApproval.value = null
  }
}

// Handle PIN modal close
const handlePinClose = () => {
  showPinModal.value = false
  selectedApproval.value = null
}

onMounted(loadPendingApprovals)
</script>
