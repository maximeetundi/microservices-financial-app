<template>
  <div class="space-y-8">
    <!-- Job Positions Section -->
    <section class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-blue-50 to-white dark:from-blue-900/20 dark:to-gray-800">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
              <BriefcaseIcon class="w-5 h-5 text-blue-500" />
              Postes & Métiers
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Définissez les postes avec leurs salaires par défaut
            </p>
          </div>
          <button @click="addPosition" 
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-medium transition-colors flex items-center gap-2">
            <PlusIcon class="w-4 h-4" />
            Ajouter un poste
          </button>
        </div>
      </div>
      
      <div class="p-6">
        <!-- Empty State -->
        <div v-if="!positions.length" class="text-center py-8 text-gray-400">
          <BriefcaseIcon class="w-12 h-12 mx-auto mb-2 opacity-50" />
          <p>Aucun poste défini</p>
          <p class="text-sm">Ajoutez des postes pour faciliter l'invitation des employés</p>
        </div>

        <!-- Positions List -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="(pos, idx) in positions" :key="pos.id" 
            class="p-4 bg-gray-50 dark:bg-gray-700/50 rounded-xl border border-gray-100 dark:border-gray-600">
            <div class="flex justify-between items-start mb-3">
              <input v-model="pos.name" type="text" placeholder="Nom du poste"
                class="bg-transparent text-lg font-semibold text-gray-900 dark:text-white border-none focus:outline-none w-full" />
              <button @click="removePosition(idx)" class="p-1 text-red-400 hover:text-red-600">
                <TrashIcon class="w-4 h-4" />
              </button>
            </div>
            
            <div class="space-y-2">
              <input v-model="pos.department" type="text" placeholder="Département (ex: Finance)"
                class="w-full px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700" />
              
              <div class="flex gap-2">
                <input v-model.number="pos.default_salary" type="number" placeholder="Salaire"
                  class="flex-1 px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700" />
                <select v-model="pos.currency" 
                  class="px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-600 dark:bg-gray-700">
                  <option value="XOF">XOF</option>
                  <option value="EUR">EUR</option>
                  <option value="USD">USD</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Admin Permissions Section -->
    <section class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-purple-50 to-white dark:from-purple-900/20 dark:to-gray-800">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <ShieldCheckIcon class="w-5 h-5 text-purple-500" />
          Gestion des Administrateurs
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Définissez les permissions de chaque administrateur
        </p>
      </div>
      
      <div class="p-6">
        <!-- Admin List -->
        <div v-if="!admins.length" class="text-center py-8 text-gray-400">
          <ShieldCheckIcon class="w-12 h-12 mx-auto mb-2 opacity-50" />
          <p>Aucun administrateur trouvé</p>
        </div>

        <div v-else class="space-y-4">
          <div v-for="admin in admins" :key="admin.id" 
            class="p-4 bg-gray-50 dark:bg-gray-700/50 rounded-xl border border-gray-100 dark:border-gray-600">
            <div class="flex items-center justify-between mb-4">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-purple-600 text-white flex items-center justify-center font-bold">
                  {{ (admin.first_name?.charAt(0) || '') + (admin.last_name?.charAt(0) || '') }}
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ admin.first_name }} {{ admin.last_name }}</p>
                  <p class="text-xs text-gray-500">{{ admin.role === 'OWNER' ? 'Propriétaire (toutes permissions)' : 'Administrateur' }}</p>
                </div>
              </div>
              <span v-if="admin.role === 'OWNER'" class="px-2 py-1 bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 text-xs font-medium rounded-full">
                Super Admin
              </span>
            </div>

            <!-- Permissions Grid -->
            <div v-if="admin.role !== 'OWNER'" class="grid grid-cols-2 md:grid-cols-4 gap-3">
              <label v-for="perm in permissionList" :key="perm.key" 
                class="flex items-center gap-2 text-sm cursor-pointer">
                <input type="checkbox" v-model="admin.permissions[perm.key]" 
                  class="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500" />
                <span class="text-gray-700 dark:text-gray-300">{{ perm.label }}</span>
              </label>
            </div>
            
            <!-- Owner has all permissions by default -->
            <div v-else class="text-sm text-gray-500 dark:text-gray-400 italic">
              Le propriétaire a automatiquement toutes les permissions et ne peut pas être modifié.
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Save Button -->
    <div class="flex justify-end">
      <button @click="save" :disabled="isSaving"
        class="px-6 py-3 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl font-medium transition-all shadow-lg disabled:opacity-50">
        {{ isSaving ? 'Enregistrement...' : 'Enregistrer les modifications' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { BriefcaseIcon, ShieldCheckIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  enterprise: { type: Object, required: true },
  employees: { type: Array, default: () => [] }
})

const emit = defineEmits(['update', 'save'])

const isSaving = ref(false)

// Job Positions
const positions = computed({
  get: () => props.enterprise.job_positions || [],
  set: (val) => emit('update', { ...props.enterprise, job_positions: val })
})

const addPosition = () => {
  const newPositions = [...positions.value, {
    id: 'pos_' + Date.now(),
    name: '',
    department: '',
    default_salary: 0,
    currency: 'XOF'
  }]
  emit('update', { ...props.enterprise, job_positions: newPositions })
}

const removePosition = (idx) => {
  const newPositions = [...positions.value]
  newPositions.splice(idx, 1)
  emit('update', { ...props.enterprise, job_positions: newPositions })
}

// Admins (filter employees with ADMIN or OWNER role)
const admins = computed(() => {
  return props.employees.filter(e => e.role === 'ADMIN' || e.role === 'OWNER')
})

// Permission definitions
const permissionList = [
  { key: 'can_invite_employees', label: 'Inviter employés' },
  { key: 'can_terminate_employees', label: 'Licencier' },
  { key: 'can_manage_payroll', label: 'Gérer paie' },
  { key: 'can_manage_services', label: 'Gérer services' },
  { key: 'can_manage_settings', label: 'Paramètres' },
  { key: 'can_manage_wallets', label: 'Wallets' },
  { key: 'can_approve_actions', label: 'Approuver actions' },
  { key: 'can_manage_admins', label: 'Gérer admins' }
]

const save = async () => {
  isSaving.value = true
  try {
    await enterpriseAPI.update(props.enterprise.id, props.enterprise)
    
    // Update admin permissions
    for (const admin of admins.value) {
      if (admin.role !== 'OWNER') {
        await enterpriseAPI.promoteEmployee(admin.id, { permissions: admin.permissions })
      }
    }
    
    alert('✅ Modifications enregistrées')
    emit('save')
  } catch (error) {
    console.error('Save failed:', error)
    alert('Erreur: ' + (error.response?.data?.error || error.message))
  } finally {
    isSaving.value = false
  }
}
</script>
