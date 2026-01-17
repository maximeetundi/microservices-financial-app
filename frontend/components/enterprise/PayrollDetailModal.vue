<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <BanknotesIcon class="w-6 h-6 text-green-500" />
          Détails de la Paie
        </h3>

        <!-- Run Info -->
        <div class="space-y-4 mb-6">
          <div class="flex justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-500">Période</span>
            <span class="font-medium text-gray-900 dark:text-white">{{ formatPeriod }}</span>
          </div>
          <div class="flex justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-500">Statut</span>
            <span :class="getStatusClass(run.status)" class="px-2 py-1 rounded text-xs font-medium">
              {{ formatStatus(run.status) }}
            </span>
          </div>
          <div class="flex justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-500">Total versé</span>
            <span class="font-bold text-green-600">{{ formatCurrency(run.total_amount) }}</span>
          </div>
          <div class="flex justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-500">Employés payés</span>
            <span class="font-medium text-gray-900 dark:text-white">{{ run.employee_count || 0 }}</span>
          </div>
          <div class="flex justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-500">Exécuté le</span>
            <span class="text-gray-900 dark:text-white">{{ formatDate(run.executed_at) }}</span>
          </div>
        </div>

        <!-- Close Button -->
        <div class="flex justify-end">
          <button @click="$emit('close')" 
            class="px-5 py-2.5 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-xl font-medium transition-colors">
            Fermer
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { XMarkIcon, BanknotesIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  run: { type: Object, required: true }
})

defineEmits(['close'])

const formatPeriod = computed(() => {
  const months = ['Jan', 'Fév', 'Mar', 'Avr', 'Mai', 'Juin', 'Juil', 'Août', 'Sep', 'Oct', 'Nov', 'Déc']
  return `${months[props.run.period_month - 1]} ${props.run.period_year}`
})

const formatCurrency = (amount) => new Intl.NumberFormat('fr-FR').format(amount || 0) + ' XOF'
const formatDate = (date) => date ? new Date(date).toLocaleDateString('fr-FR', { dateStyle: 'long' }) : '--'

const formatStatus = (status) => {
  const map = { COMPLETED: 'Terminé', PROCESSING: 'En cours', DRAFT: 'Brouillon', FAILED: 'Échoué' }
  return map[status] || status
}

const getStatusClass = (status) => {
  const map = {
    COMPLETED: 'bg-green-100 text-green-700',
    PROCESSING: 'bg-blue-100 text-blue-700',
    DRAFT: 'bg-gray-100 text-gray-700',
    FAILED: 'bg-red-100 text-red-700'
  }
  return map[status] || 'bg-gray-100 text-gray-700'
}
</script>
