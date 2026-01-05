<template>
  <div class="space-y-6">
    <!-- Emergency Fund Header -->
    <div v-if="!fund" class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
      <div class="text-center py-8">
        <h3 class="text-lg font-bold mb-2">Caisse de Secours</h3>
        <p class="text-gray-500 mb-4">Fonds d'urgence pour événements familiaux (deuil, mariage, naissance, maladie)</p>
        <button @click="showCreateModal = true" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg font-medium">
          Créer la Caisse de Secours
        </button>
      </div>
    </div>

    <!-- Emergency Fund Stats -->
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="text-sm text-gray-500 mb-1">Solde disponible</div>
          <div class="text-3xl font-bold text-green-600">{{ formatCurrency(fund.balance) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="text-sm text-gray-500 mb-1">Cotisation mensuelle</div>
          <div class="text-3xl font-bold text-indigo-600">{{ formatCurrency(fund.monthly_contribution) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 flex items-center justify-center">
          <button @click="showContributeModal = true" class="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-medium">
            💰 Payer ma cotisation
          </button>
        </div>
      </div>

      <!-- Request Withdrawal Button -->
      <div class="mb-6">
        <button @click="showWithdrawalModal = true" class="bg-orange-600 hover:bg-orange-700 text-white px-6 py-3 rounded-lg font-medium">
          📤 Demander une aide
        </button>
        <span class="text-sm text-gray-500 ml-3">Nécessite l'approbation de 4 membres sur 5</span>
      </div>

      <!-- Withdrawals List -->
      <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
        <div class="p-6 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-bold">Demandes d'aide</h3>
        </div>
        <div class="divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="withdrawal in withdrawals" :key="withdrawal.id" class="p-6">
            <div class="flex justify-between items-start">
              <div>
                <div class="flex items-center gap-2 mb-2">
                  <span class="text-2xl">{{ getEventIcon(withdrawal.event_type) }}</span>
                  <span class="font-medium">{{ getEventLabel(withdrawal.event_type) }}</span>
                  <span :class="['px-2 py-1 rounded-full text-xs font-medium', getStatusColor(withdrawal.status)]">
                    {{ withdrawal.status }}
                  </span>
                </div>
                <p class="text-gray-600 dark:text-gray-400 text-sm">{{ withdrawal.reason }}</p>
                <p class="text-xs text-gray-500 mt-1">Bénéficiaire: {{ withdrawal.beneficiary_name }}</p>
              </div>
              <div class="text-right">
                <div class="text-2xl font-bold text-orange-600">{{ formatCurrency(withdrawal.amount) }}</div>
                <div class="text-xs text-gray-500">{{ new Date(withdrawal.created_at).toLocaleDateString() }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Fund Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Créer Caisse de Secours</h3>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Cotisation mensuelle par membre</label>
          <input v-model.number="createForm.monthlyContribution" type="number" step="100" class="w-full px-4 py-2 rounded-lg border" placeholder="2000">
        </div>
        <div class="flex gap-3">
          <button @click="showCreateModal = false" class="flex-1 px-4 py-2 border rounded-lg">Annuler</button>
          <button @click="createFund" class="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-lg">Créer</button>
        </div>
      </div>
    </div>

    <!-- Contribute Modal -->
    <div v-if="showContributeModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Payer cotisation mensuelle</h3>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Période</label>
          <select v-model="contributeForm.period" class="w-full px-4 py-2 rounded-lg border">
            <option value="janvier_2026">Janvier 2026</option>
            <option value="fevrier_2026">Février 2026</option>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Montant</label>
          <input :value="fund?.monthly_contribution" readonly class="w-full px-4 py-2 rounded-lg border bg-gray-100">
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">PIN</label>
          <input v-model="contributeForm.pin" type="password" maxlength="5" class="w-full px-4 py-2 rounded-lg border text-center">
        </div>
        <div class="flex gap-3">
          <button @click="showContributeModal = false" class="flex-1 px-4 py-2 border rounded-lg">Annuler</button>
          <button @click="contribute" class="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg">Payer</button>
        </div>
      </div>
    </div>

    <!-- Withdrawal Request Modal -->
    <div v-if="showWithdrawalModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Demander une aide</h3>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Type d'événement</label>
          <select v-model="withdrawalForm.event_type" class="w-full px-4 py-2 rounded-lg border">
            <option value="deuil">🕊️ Deuil</option>
            <option value="mariage">💍 Mariage</option>
            <option value="naissance">👶 Naissance</option>
            <option value="maladie">🏥 Maladie</option>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Montant demandé</label>
          <input v-model.number="withdrawalForm.amount" type="number" class="w-full px-4 py-2 rounded-lg border">
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium mb-2">Motif</label>
          <textarea v-model="withdrawalForm.reason" rows="3" class="w-full px-4 py-2 rounded-lg border"></textarea>
        </div>
        <div class="flex gap-3">
          <button @click="showWithdrawalModal = false" class="flex-1 px-4 py-2 border rounded-lg">Annuler</button>
          <button @click="requestWithdrawal" class="flex-1 px-4 py-2 bg-orange-600 text-white rounded-lg">Demander</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { associationAPI } from '~/composables/useApi'
import { useNotification } from '~/composables/useNotification'

const props = defineProps<{ associationId: string }>()

const { showSuccess, showError, showWarning } = useNotification()

const fund = ref<any>(null)
const withdrawals = ref<any[]>([])
const showCreateModal = ref(false)
const showContributeModal = ref(false)
const showWithdrawalModal = ref(false)

const createForm = ref({ monthlyContribution: 2000 })
const contributeForm = ref({ period: 'janvier_2026', pin: '', wallet_id: '' })
const withdrawalForm = ref({ event_type: 'deuil', amount: 0, reason: '', beneficiary_id: '' })

const loadData = async () => {
  try {
    const res = await associationAPI.getEmergencyFund(props.associationId)
    fund.value = res.data
    loadWithdrawals()
  } catch (err) {
    fund.value = null
  }
}

const loadWithdrawals = async () => {
  try {
    const res = await associationAPI.getEmergencyWithdrawals(props.associationId)
    withdrawals.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

const createFund = async () => {
  try {
    await associationAPI.createEmergencyFund(props.associationId, createForm.value.monthlyContribution)
    showCreateModal.value = false
    showSuccess('Caisse de secours créée avec succès !')
    loadData()
  } catch (err: any) {
    showError(err.response?.data?.error || 'Erreur lors de la création')
  }
}

const contribute = async () => {
  try {
    await associationAPI.contributeToEmergencyFund(props.associationId, {
      period: contributeForm.value.period,
      wallet_id: contributeForm.value.wallet_id || 'default',
      pin: contributeForm.value.pin,
      amount: fund.value.monthly_contribution
    })
    showContributeModal.value = false
    showSuccess('Cotisation mensuelle payée avec succès !')
    contributeForm.value.pin = ''
    loadData()
  } catch (err: any) {
    showError(err.response?.data?.error || 'Erreur lors du paiement')
  }
}

const requestWithdrawal = async () => {
  try {
    await associationAPI.requestEmergencyWithdrawal(props.associationId, {
      ...withdrawalForm.value,
      beneficiary_id: withdrawalForm.value.beneficiary_id || 'self'
    })
    showWithdrawalModal.value = false
    showSuccess('Demande créée ! En attente de 4 approbations sur 5.', 'Demande enregistrée')
    withdrawalForm.value = { event_type: 'deuil', amount: 0, reason: '', beneficiary_id: '' }
    loadWithdrawals()
  } catch (err: any) {
    showError(err.response?.data?.error || 'Erreur lors de la demande')
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'XOF' }).format(amount || 0)
}

const getEventIcon = (type: string) => {
  const icons: any = { deuil: '🕊️', mariage: '💍', naissance: '👶', maladie: '🏥' }
  return icons[type] || '📌'
}

const getEventLabel = (type: string) => {
  const labels: any = { deuil: 'Deuil', mariage: 'Mariage', naissance: 'Naissance', maladie: 'Maladie' }
  return labels[type] || type
}

const getStatusColor = (status: string) => {
  if (status === 'approved') return 'bg-green-100 text-green-700'
  if (status === 'rejected') return 'bg-red-100 text-red-700'
  return 'bg-yellow-100 text-yellow-700'
}

onMounted(() => loadData())
</script>
