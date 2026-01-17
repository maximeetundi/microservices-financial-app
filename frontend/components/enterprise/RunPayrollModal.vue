<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <PlayIcon class="w-6 h-6 text-green-500" />
          Exécuter la Paie
        </h3>

        <!-- Preview Loading -->
        <div v-if="isLoadingPreview" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto mb-3"></div>
          <p class="text-gray-500">Calcul en cours...</p>
        </div>

        <!-- Preview Results -->
        <div v-else-if="preview" class="space-y-4 mb-6">
          <div class="p-4 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200 dark:border-green-800">
            <p class="text-sm text-green-700 dark:text-green-300 mb-1">Total à verser</p>
            <p class="text-2xl font-bold text-green-600">{{ formatCurrency(preview.total_amount) }}</p>
          </div>
          
          <div class="grid grid-cols-2 gap-4">
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded-lg text-center">
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ preview.employee_count }}</p>
              <p class="text-xs text-gray-500">Employés</p>
            </div>
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded-lg text-center">
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ formatPeriod }}</p>
              <p class="text-xs text-gray-500">Période</p>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="p-3 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
            <p class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
          </div>

          <!-- PIN Verification -->
          <div class="p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
            <p class="text-sm text-amber-700 dark:text-amber-300 mb-2">Vérification PIN requise</p>
            <p class="text-xs text-amber-600 dark:text-amber-400">Cette action nécessite votre code PIN de sécurité.</p>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex justify-end gap-3 mt-6">
          <button @click="$emit('close')" 
            class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
            Annuler
          </button>
          <button @click="executePayroll" :disabled="isExecuting || !preview"
            class="px-5 py-2.5 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 text-white rounded-xl font-medium shadow-lg shadow-green-500/25 disabled:opacity-50 flex items-center gap-2">
            <div v-if="isExecuting" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            {{ isExecuting ? 'Exécution...' : 'Confirmer et Exécuter' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { XMarkIcon, PlayIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const props = defineProps({
  enterpriseId: { type: String, required: true },
  month: { type: Number, required: true },
  year: { type: Number, required: true }
})

const emit = defineEmits(['close', 'executed'])

const { requirePin } = usePin()

const preview = ref(null)
const isLoadingPreview = ref(true)
const isExecuting = ref(false)
const errorMessage = ref('')

const formatPeriod = computed(() => {
  const months = ['Jan', 'Fév', 'Mar', 'Avr', 'Mai', 'Juin', 'Juil', 'Août', 'Sep', 'Oct', 'Nov', 'Déc']
  return `${months[props.month - 1]} ${props.year}`
})

const formatCurrency = (amount) => new Intl.NumberFormat('fr-FR').format(amount || 0) + ' XOF'

onMounted(async () => {
  try {
    const { data } = await enterpriseAPI.previewPayroll(props.enterpriseId)
    preview.value = data
  } catch (e) {
    console.error('Failed to preview payroll:', e)
    errorMessage.value = e.response?.data?.error || 'Erreur lors du calcul'
  } finally {
    isLoadingPreview.value = false
  }
})

const executePayroll = async () => {
  const pinVerified = await requirePin('Exécution de la paie')
  if (!pinVerified) return

  isExecuting.value = true
  errorMessage.value = ''
  
  try {
    await enterpriseAPI.executePayroll(props.enterpriseId)
    emit('executed')
    emit('close')
  } catch (e) {
    console.error('Payroll execution failed:', e)
    errorMessage.value = e.response?.data?.details || e.response?.data?.error || 'Erreur lors de l\'exécution'
  } finally {
    isExecuting.value = false
  }
}
</script>
