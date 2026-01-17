<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-md p-6 shadow-2xl">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <ShieldCheckIcon class="w-6 h-6 text-purple-500" />
          Promouvoir en Administrateur
        </h3>

        <!-- Employee Info -->
        <div class="bg-gray-50 dark:bg-gray-700 rounded-xl p-4 mb-6">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-purple-600 text-white flex items-center justify-center font-bold">
              {{ (employee.first_name?.charAt(0) || '') + (employee.last_name?.charAt(0) || '') }}
            </div>
            <div>
              <p class="font-medium text-gray-900 dark:text-white">{{ employee.first_name }} {{ employee.last_name }}</p>
              <p class="text-sm text-gray-500">{{ employee.profession || 'Employé' }}</p>
            </div>
          </div>
        </div>

        <!-- Permissions Selection -->
        <div class="mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Permissions à accorder</h4>
          <div class="grid grid-cols-2 gap-3">
            <label v-for="perm in permissionList" :key="perm.key" 
              class="flex items-center gap-2 text-sm cursor-pointer p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
              <input type="checkbox" v-model="selectedPermissions[perm.key]" 
                class="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500" />
              <span class="text-gray-700 dark:text-gray-300">{{ perm.label }}</span>
            </label>
          </div>
        </div>

        <!-- Multi-Approval Notice -->
        <div class="mb-6 p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
          <div class="flex gap-3">
            <ExclamationTriangleIcon class="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" />
            <div class="text-sm">
              <p class="font-medium text-amber-700 dark:text-amber-300 mb-1">Approbation multi-admin requise</p>
              <p class="text-amber-600 dark:text-amber-400">
                Cette action nécessite l'approbation de tous les administrateurs existants.
              </p>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <p class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
        </div>

        <!-- Actions -->
        <div class="flex justify-end gap-3">
          <button @click="$emit('close')" 
            class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
            Annuler
          </button>
          <button @click="initiatePromotion" :disabled="isSubmitting"
            class="px-5 py-2.5 bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-700 hover:to-purple-800 text-white rounded-xl font-medium shadow-lg shadow-purple-500/25 disabled:opacity-50 flex items-center gap-2">
            <div v-if="isSubmitting" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            {{ isSubmitting ? '' : 'Initier la Promotion' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { XMarkIcon, ShieldCheckIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  employee: { type: Object, required: true },
  enterpriseId: { type: String, required: true }
})

const emit = defineEmits(['close', 'promoted'])

const isSubmitting = ref(false)
const errorMessage = ref('')

// Permission definitions
const permissionList = [
  { key: 'can_invite_employees', label: 'Inviter employés' },
  { key: 'can_terminate_employees', label: 'Licencier' },
  { key: 'can_manage_payroll', label: 'Gérer paie' },
  { key: 'can_manage_services', label: 'Gérer services' },
  { key: 'can_manage_settings', label: 'Paramètres' },
  { key: 'can_manage_wallets', label: 'Wallets' },
  { key: 'can_approve_actions', label: 'Approuver actions' }
]

// All permissions selected by default
const selectedPermissions = reactive(
  Object.fromEntries(permissionList.map(p => [p.key, true]))
)

const initiatePromotion = async () => {
  isSubmitting.value = true
  errorMessage.value = ''

  try {
    // Use the multi-approval system to initiate promotion
    await enterpriseAPI.initiateAction(props.enterpriseId, {
      action_type: 'PROMOTE_ADMIN',
      action_name: `Promouvoir ${props.employee.first_name} ${props.employee.last_name} en admin`,
      description: `Promotion de ${props.employee.first_name} ${props.employee.last_name} au rôle d'administrateur`,
      payload: {
        employee_id: props.employee.id || props.employee._id,
        new_role: 'ADMIN',
        permissions: selectedPermissions
      }
    })
    
    emit('promoted')
  } catch (e) {
    console.error('Promotion initiation failed:', e)
    errorMessage.value = e.response?.data?.details || e.response?.data?.error || 'Erreur lors de l\'initiation de la promotion'
  } finally {
    isSubmitting.value = false
  }
}
</script>
