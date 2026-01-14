<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden shadow-2xl flex flex-col">
        <!-- Header -->
        <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-primary-50 to-white dark:from-gray-800 dark:to-gray-800">
          <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
            <XMarkIcon class="w-5 h-5" />
          </button>

          <div class="flex items-center gap-4">
            <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 text-white flex items-center justify-center font-bold text-xl shadow-lg">
              {{ initials }}
            </div>
            <div>
              <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                {{ employee.first_name }} {{ employee.last_name }}
              </h3>
              <p class="text-gray-500">{{ employee.profession || 'Poste non défini' }}</p>
              <span :class="getStatusClass(employee.status)" 
                class="inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium mt-2">
                {{ formatStatus(employee.status) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Content -->
        <div class="p-6 overflow-y-auto flex-1 space-y-6">
          <!-- Contact Info -->
          <div class="grid grid-cols-2 gap-4">
            <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
              <p class="text-xs text-gray-500 uppercase mb-1">Email</p>
              <p class="font-medium text-gray-900 dark:text-white">{{ employee.email || '--' }}</p>
            </div>
            <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
              <p class="text-xs text-gray-500 uppercase mb-1">Téléphone</p>
              <p class="font-medium text-gray-900 dark:text-white">{{ employee.phone_number || '--' }}</p>
            </div>
          </div>

          <!-- Salary Configuration -->
          <div class="p-4 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200 dark:border-green-800">
            <div class="flex items-center justify-between mb-4">
              <h4 class="font-semibold text-green-800 dark:text-green-200 flex items-center gap-2">
                <BanknotesIcon class="w-5 h-5" />
                Configuration Salaire
              </h4>
              <button @click="editingSalary = !editingSalary" 
                class="text-sm text-green-600 hover:text-green-700 font-medium">
                {{ editingSalary ? 'Annuler' : 'Modifier' }}
              </button>
            </div>

            <div v-if="!editingSalary" class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div>
                <p class="text-xs text-gray-500">Salaire Brut</p>
                <p class="font-bold text-lg text-gray-900 dark:text-white">{{ formatCurrency(salaryConfig.base_amount) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Déductions</p>
                <p class="font-bold text-lg text-red-600">- {{ formatCurrency(totalDeductions) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Primes</p>
                <p class="font-bold text-lg text-green-600">+ {{ formatCurrency(totalBonuses) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Salaire Net</p>
                <p class="font-bold text-lg text-primary-600">{{ formatCurrency(netPayable) }}</p>
              </div>
            </div>

            <!-- Salary Editor -->
            <div v-else class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Salaire de base</label>
                  <input v-model.number="editForm.base_amount" type="number"
                    class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Fréquence</label>
                  <select v-model="editForm.frequency"
                    class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700">
                    <option value="MONTHLY">Mensuel</option>
                    <option value="WEEKLY">Hebdomadaire</option>
                  </select>
                </div>
              </div>

              <!-- Deductions -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Déductions</label>
                  <button @click="addDeduction" class="text-xs text-primary-600">+ Ajouter</button>
                </div>
                <div v-for="(ded, idx) in editForm.deductions" :key="idx" class="flex gap-2 mb-2">
                  <input v-model="ded.name" placeholder="Nom" class="flex-1 px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                  <select v-model="ded.type" class="px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm">
                    <option value="FIXED">Fixe</option>
                    <option value="PERCENTAGE">%</option>
                  </select>
                  <input v-model.number="ded.value" type="number" class="w-20 px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                  <button @click="editForm.deductions.splice(idx, 1)" class="text-red-500"><XMarkIcon class="w-4 h-4" /></button>
                </div>
              </div>

              <!-- Bonuses -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Primes & Bonus</label>
                  <button @click="addBonus" class="text-xs text-primary-600">+ Ajouter</button>
                </div>
                <div v-for="(bon, idx) in editForm.bonuses" :key="idx" class="flex gap-2 mb-2">
                  <input v-model="bon.name" placeholder="Nom" class="flex-1 px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                  <select v-model="bon.type" class="px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm">
                    <option value="FIXED">Fixe</option>
                    <option value="PERCENTAGE">%</option>
                  </select>
                  <input v-model.number="bon.value" type="number" class="w-20 px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                  <button @click="editForm.bonuses.splice(idx, 1)" class="text-red-500"><XMarkIcon class="w-4 h-4" /></button>
                </div>
              </div>

              <button @click="saveSalary" :disabled="isSaving"
                class="w-full py-2.5 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium disabled:opacity-50">
                {{ isSaving ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </div>
          </div>

          <!-- Career History -->
          <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
            <h4 class="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <ClockIcon class="w-5 h-5" />
              Historique de Carrière
            </h4>
            <div v-if="!employee.history?.length" class="text-sm text-gray-500 text-center py-4">
              Aucun événement enregistré
            </div>
            <div v-else class="space-y-3">
              <div v-for="(event, idx) in employee.history" :key="idx" class="flex gap-3">
                <div class="w-8 h-8 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center flex-shrink-0">
                  <ArrowTrendingUpIcon v-if="event.type === 'PROMOTION'" class="w-4 h-4" />
                  <BanknotesIcon v-else-if="event.type === 'SALARY_CHANGE'" class="w-4 h-4" />
                  <XCircleIcon v-else class="w-4 h-4" />
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white text-sm">{{ event.description }}</p>
                  <p class="text-xs text-gray-500">{{ new Date(event.date).toLocaleDateString() }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="p-4 border-t border-gray-100 dark:border-gray-700 flex justify-end gap-3">
          <button @click="$emit('close')" class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg">
            Fermer
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { 
  XMarkIcon, BanknotesIcon, ClockIcon, ArrowTrendingUpIcon, XCircleIcon
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  employee: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'updated'])

const editingSalary = ref(props.employee._editingSalary || false)
const isSaving = ref(false)

const salaryConfig = computed(() => props.employee.salary_config || {})
const editForm = reactive({
  base_amount: salaryConfig.value.base_amount || 0,
  frequency: salaryConfig.value.frequency || 'MONTHLY',
  deductions: [...(salaryConfig.value.deductions || [])],
  bonuses: [...(salaryConfig.value.bonuses || [])]
})

const initials = computed(() => {
  return (props.employee.first_name?.charAt(0) || '') + (props.employee.last_name?.charAt(0) || '')
})

const totalDeductions = computed(() => {
  return (salaryConfig.value.deductions || []).reduce((sum, d) => {
    if (d.type === 'PERCENTAGE') {
      return sum + (salaryConfig.value.base_amount * d.value / 100)
    }
    return sum + (d.value || 0)
  }, 0)
})

const totalBonuses = computed(() => {
  return (salaryConfig.value.bonuses || []).reduce((sum, b) => {
    if (b.type === 'PERCENTAGE') {
      return sum + (salaryConfig.value.base_amount * b.value / 100)
    }
    return sum + (b.value || 0)
  }, 0)
})

const netPayable = computed(() => {
  return (salaryConfig.value.base_amount || 0) - totalDeductions.value + totalBonuses.value
})

const getStatusClass = (status) => {
  switch(status) {
    case 'ACTIVE': return 'bg-green-100 text-green-700'
    case 'PENDING_INVITE': return 'bg-amber-100 text-amber-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}

const formatStatus = (status) => {
  const map = { 'ACTIVE': 'Actif', 'PENDING_INVITE': 'En attente', 'TERMINATED': 'Licencié' }
  return map[status] || status
}

const formatCurrency = (amount) => new Intl.NumberFormat('fr-FR').format(amount || 0) + ' XOF'

const addDeduction = () => {
  editForm.deductions.push({ name: '', type: 'FIXED', value: 0 })
}

const addBonus = () => {
  editForm.bonuses.push({ name: '', type: 'FIXED', value: 0 })
}

const saveSalary = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update salary
    // await enterpriseAPI.updateEmployeeSalary(props.employee.enterprise_id, props.employee.id, editForm)
    emit('updated')
  } catch (e) {
    alert('Erreur: ' + e.message)
  } finally {
    isSaving.value = false
  }
}
</script>
