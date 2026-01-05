<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full mx-4">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Demander un prêt</h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">✕</button>
        </div>
      </div>

      <div class="p-6 space-y-4">
        <!-- Amount -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Montant demandé</label>
          <div class="relative">
            <input v-model.number="form.amount" type="number" step="1000"
                   class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                   placeholder="0" />
            <span class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500">{{ currency }}</span>
          </div>
        </div>

        <!-- Duration -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Durée (mois)</label>
          <select v-model.number="form.duration"
                  class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500">
            <option :value="1">1 mois</option>
            <option :value="2">2 mois</option>
            <option :value="3">3 mois</option>
            <option :value="6">6 mois</option>
            <option :value="12">12 mois</option>
          </select>
        </div>

        <!-- Interest Rate Info -->
        <div class="bg-indigo-50 dark:bg-indigo-900/20 p-4 rounded-lg">
          <div class="flex justify-between items-center">
            <span class="text-sm text-indigo-700 dark:text-indigo-300">Taux d'intérêt</span>
            <span class="font-bold text-indigo-700 dark:text-indigo-300">{{ interestRate }}%</span>
          </div>
          <div class="flex justify-between items-center mt-2">
            <span class="text-sm text-indigo-700 dark:text-indigo-300">Montant total à rembourser</span>
            <span class="font-bold text-indigo-700 dark:text-indigo-300">{{ formatBalance(totalRepayment) }}</span>
          </div>
        </div>

        <!-- Reason -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Motif</label>
          <textarea v-model="form.reason" rows="2"
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                    placeholder="Raison de votre demande..."></textarea>
        </div>

        <!-- Error -->
        <div v-if="error" class="p-3 bg-red-50 dark:bg-red-900/20 text-red-600 rounded-lg text-sm">{{ error }}</div>
      </div>

      <div class="p-6 border-t border-gray-100 dark:border-gray-700 flex space-x-3">
        <button @click="$emit('close')" 
                class="flex-1 px-4 py-3 border border-gray-300 rounded-lg text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50">
          Annuler
        </button>
        <button @click="submit" :disabled="loading || !isValid"
                class="flex-1 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white rounded-lg font-medium flex items-center justify-center">
          <span v-if="loading" class="animate-spin mr-2 h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
          Demander
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { associationAPI } from '~/composables/useApi'

const props = defineProps<{
  show: boolean
  associationId: string
  currency: string
  interestRate?: number
}>()

const emit = defineEmits(['close', 'success'])


const loading = ref(false)
const error = ref('')

const form = ref({
  amount: 0,
  duration: 3,
  reason: ''
})

const isValid = computed(() => form.value.amount > 0 && form.value.duration > 0)

const totalRepayment = computed(() => {
  const rate = props.interestRate || 5
  return form.value.amount * (1 + rate / 100)
})

const formatBalance = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: props.currency || 'XOF' }).format(amount || 0)
}

const submit = async () => {
  if (!isValid.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    await associationAPI.requestLoan(props.associationId, {
      amount: form.value.amount,
      interest_rate: props.interestRate || 5,
      duration: form.value.duration,
      reason: form.value.reason
    })
    emit('success')
    emit('close')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Échec de la demande'
  } finally {
    loading.value = false
  }
}
</script>
