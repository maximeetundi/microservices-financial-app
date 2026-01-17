<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-md p-6 shadow-2xl max-h-[90vh] overflow-y-auto">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <UserPlusIcon class="w-6 h-6 text-green-500" />
          Inviter un Employé
        </h3>

        <!-- Consent Notice -->
        <div class="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200 dark:border-blue-800">
          <div class="flex gap-3">
            <InformationCircleIcon class="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
            <div class="text-sm text-blue-700 dark:text-blue-300">
              <p class="font-medium mb-1">Processus avec consentement</p>
              <p class="text-blue-600 dark:text-blue-400">L'employé recevra une notification pour rejoindre l'entreprise.</p>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="mb-4 p-4 bg-red-50 dark:bg-red-900/20 rounded-xl border border-red-200 dark:border-red-800">
          <div class="flex gap-3">
            <XMarkIcon class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
            <div class="text-sm text-red-700 dark:text-red-300">
              <p class="font-medium mb-1">Erreur</p>
              <p class="text-red-600 dark:text-red-400">{{ errorMessage }}</p>
            </div>
            <button @click="errorMessage = ''" class="ml-auto text-red-500 hover:text-red-700">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <form @submit.prevent="sendInvite" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Prénom *</label>
              <input v-model="form.first_name" type="text" required
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nom *</label>
              <input v-model="form.last_name" type="text" required
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
            </div>
          </div>

          <!-- Position Selection -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Poste / Métier</label>
            <div class="relative">
              <select v-model="selectedPositionId" @change="handlePositionChange"
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 appearance-none">
                <option value="">Sélectionner ou saisir manuellement...</option>
                <option v-for="pos in jobPositions" :key="pos.id" :value="pos.id">
                  {{ pos.name }} ({{ pos.department }})
                </option>
              </select>
              <BriefcaseIcon class="w-5 h-5 text-gray-400 absolute right-3 top-2.5 pointer-events-none" />
            </div>
            
            <!-- Manual input if no position selected -->
            <input v-if="!selectedPositionId" v-model="form.profession" type="text" placeholder="Saisir l'intitulé du poste..."
              class="mt-2 w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
          </div>

          <!-- Contact -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Email ou Téléphone *</label>
            <input v-model="form.contact" type="text" required placeholder="email@example.com ou +221..."
              class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
          </div>

          <!-- Hire Date -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Date d'embauche</label>
            <input v-model="form.hire_date" type="date" required
              class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
            <p class="text-xs text-gray-500 mt-1">Utilisé pour calculer l'ancienneté</p>
          </div>

          <!-- Salary Setup -->
          <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
            <div class="flex items-center justify-between mb-3">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Configuration salaire</span>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" v-model="showSalary" class="rounded border-gray-300 text-primary-600" />
                <span class="text-xs text-gray-500">Configurer maintenant</span>
              </label>
            </div>
            
            <div v-if="showSalary" class="space-y-3">
              <div>
                <label class="block text-xs text-gray-500 mb-1">Salaire de base</label>
                <div class="flex items-center gap-2">
                  <input v-model.number="form.salary_config.base_amount" type="number"
                    class="flex-1 px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-600 text-sm" />
                  <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ form.salary_config.currency || 'XOF' }}</span>
                </div>
              </div>
              <div>
                <label class="block text-xs text-gray-500 mb-1">Fréquence</label>
                <select v-model="form.salary_config.frequency"
                  class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-600 text-sm">
                  <option value="MONTHLY">Mensuel</option>
                  <option value="WEEKLY">Hebdomadaire</option>
                </select>
              </div>
            </div>
          </div>

          <div class="flex gap-3 pt-4">
            <button type="button" @click="$emit('close')"
              class="flex-1 px-4 py-2.5 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              Annuler
            </button>
            <button type="submit" :disabled="isLoading"
              class="flex-1 px-4 py-2.5 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-xl font-medium hover:from-green-700 hover:to-green-800 disabled:opacity-50 transition-all">
              {{ isLoading ? 'Envoi...' : 'Envoyer l\'invitation' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { UserPlusIcon, XMarkIcon, InformationCircleIcon, BriefcaseIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  enterpriseId: { type: String, required: true },
  jobPositions: { type: Array, default: () => [] }
})

const emit = defineEmits(['close', 'invited'])

const isLoading = ref(false)
const showSalary = ref(false)
const selectedPositionId = ref('')
const errorMessage = ref('')

const form = reactive({
  first_name: '',
  last_name: '',
  profession: '',
  contact: '',
  hire_date: new Date().toISOString().split('T')[0],
  salary_config: {
    base_amount: 0,
    frequency: 'MONTHLY',
    currency: 'XOF', // Default
    deductions: [],
    bonuses: []
  }
})

// Auto-fill when position changes
const handlePositionChange = () => {
  if (!selectedPositionId.value) {
    form.profession = ''
    return
  }
  
  const position = props.jobPositions.find(p => p.id === selectedPositionId.value)
  if (position) {
    form.profession = position.name
    form.salary_config.base_amount = position.default_salary
    form.salary_config.currency = position.currency || 'XOF'
    showSalary.value = true
  }
}

const sendInvite = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    // Validate required fields
    if (!form.first_name || !form.last_name) {
      errorMessage.value = 'Le prénom et le nom sont requis.'
      isLoading.value = false
      return
    }
    if (!form.contact) {
      errorMessage.value = 'L\'email ou le téléphone est requis.'
      isLoading.value = false
      return
    }
    
    const payload = {
      first_name: form.first_name,
      last_name: form.last_name,
      profession: form.profession,
      position_id: selectedPositionId.value || undefined,
      hire_date: new Date(form.hire_date).toISOString(),
      email: form.contact.includes('@') ? form.contact : '',
      phone_number: !form.contact.includes('@') ? form.contact : '',
      salary_config: showSalary.value ? form.salary_config : undefined
    }
    
    await enterpriseAPI.inviteEmployee(props.enterpriseId, payload)
    emit('invited')
  } catch (e) {
    console.error('Invite failed:', e)
    // Parse different error formats - prioritize 'details' which contains the real error message
    let msg = 'Erreur inconnue'
    if (e.response?.data?.details) {
      msg = e.response.data.details
    } else if (e.response?.data?.error && e.response.data.error !== 'Failed to invite employee') {
      msg = e.response.data.error
    } else if (e.response?.data?.message) {
      msg = e.response.data.message
    } else if (e.message === 'Network Error') {
      msg = 'Erreur réseau. Vérifiez votre connexion.'
    } else if (e.message) {
      msg = e.message
    }
    errorMessage.value = msg
  } finally {
    isLoading.value = false
  }
}
</script>
