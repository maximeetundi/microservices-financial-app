<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold">Approbations en attente</h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">✕</button>
        </div>
      </div>

      <div class="p-6 space-y-4">
        <div v-if="loading" class="text-center py-8">Chargement...</div>

        <div v-else-if="requests.length === 0" class="text-center py-8 text-gray-500">
          Aucune approbation en attente
        </div>

        <div v-else v-for="req in requests" :key="req.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <div class="flex justify-between items-start mb-3">
            <div>
              <div class="font-medium">{{ req.request_type === 'loan' ? 'Prêt' : 'Distribution' }}</div>
              <div class="text-2xl font-bold text-indigo-600">{{ formatAmount(req.amount) }}</div>
              <div class="text-sm text-gray-500">{{ req.description }}</div>
            </div>
            <span :class="['px-2 py-1 rounded text-xs', req.status === 'pending' ? 'bg-yellow-100 text-yellow-800' : 'bg-green-100 text-green-800']">
              {{ req.status }}
            </span>
          </div>

          <div class="flex items-center gap-2 text-sm mb-3">
            <div class="flex-1 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div :style="{width: `${(req.current_approvals / req.required_approvals) * 100}%`}"
                class="bg-green-500 h-2 rounded-full"></div>
            </div>
            <span class="text-gray-600">{{ req.current_approvals }}/{{ req.required_approvals }}</span>
          </div>

          <div v-if="req.status === 'pending'" class="flex gap-2">
            <button @click="vote(req.id, 'approve')" class="flex-1 btn-primary">Approuver</button>
            <button @click="vote(req.id, 'reject')" class="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700">
              Rejeter
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { associationAPI } from '~/composables/useApi'
import { useNotification } from '~/composables/useNotification'

const props = defineProps<{
  show: boolean
  associationId: string
}>()

const emit = emit(['close', 'refresh'])

const { showSuccess, showError } = useNotification()

const loading = ref(false)
const requests = ref<any[]>([])

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'XOF' }).format(amount)
}

const loadRequests = async () => {
  loading.value = true
  try {
    const res = await associationAPI.getPendingApprovals(props.associationId)
    requests.value = res.data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const vote = async (requestId: string, voteType: 'approve' | 'reject') => {
  try {
    await associationAPI.voteOnApproval(requestId, voteType)
    loadRequests()
    emit('refresh')
  } catch (err: any) {
    showError(err.response?.data?.error || 'Erreur lors du vote')
  } finally {
  }
}

watch(() => props.show, (val) => {
  if (val) loadRequests()
})
</script>
