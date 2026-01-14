<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <UsersIcon class="w-6 h-6 text-primary-500" />
          Annuaire des Employés
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Gérez votre équipe et leurs informations de paie
        </p>
      </div>
      <button @click="showInviteModal = true" 
        class="px-5 py-2.5 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 text-white rounded-xl text-sm font-semibold shadow-lg shadow-green-500/25 transition-all flex items-center gap-2 transform hover:-translate-y-0.5">
        <UserPlusIcon class="w-5 h-5" />
        Inviter un Employé
      </button>
    </div>

    <!-- Stats Strip -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Total</p>
        <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ employees.length }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Actifs</p>
        <p class="text-2xl font-bold text-green-600">{{ activeCount }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">En attente</p>
        <p class="text-2xl font-bold text-amber-600">{{ pendingCount }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Masse salariale</p>
        <p class="text-2xl font-bold text-blue-600">{{ formatCurrency(totalMasseSalariale) }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-3">
      <div class="relative">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input v-model="searchQuery" 
          type="text" 
          placeholder="Rechercher par nom..." 
          class="pl-10 pr-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 w-64" />
      </div>
      <select v-model="statusFilter" 
        class="px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500">
        <option value="">Tous les statuts</option>
        <option value="ACTIVE">Actifs</option>
        <option value="PENDING_INVITE">En attente</option>
        <option value="TERMINATED">Licenciés</option>
      </select>
    </div>

    <!-- Table -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div v-if="isLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
        <p class="mt-4 text-gray-500">Chargement...</p>
      </div>

      <div v-else-if="filteredEmployees.length === 0" class="text-center py-16">
        <UserGroupIcon class="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <h4 class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">Aucun employé trouvé</h4>
        <p class="text-gray-500 text-sm mb-6">Invitez votre premier membre d'équipe</p>
        <button @click="showInviteModal = true" 
          class="px-5 py-2.5 bg-primary-600 hover:bg-primary-700 text-white rounded-xl text-sm font-medium transition-colors">
          Inviter un employé
        </button>
      </div>

      <table v-else class="w-full text-left text-sm">
        <thead class="bg-gray-50 dark:bg-gray-900/50 text-gray-600 dark:text-gray-400 text-xs uppercase tracking-wider">
          <tr>
            <th class="px-6 py-4 font-semibold">Employé</th>
            <th class="px-6 py-4 font-semibold">Poste</th>
            <th class="px-6 py-4 font-semibold">Statut</th>
            <th class="px-6 py-4 font-semibold text-right">Salaire Net</th>
            <th class="px-6 py-4 font-semibold text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
          <tr v-for="emp in filteredEmployees" :key="emp.id" 
            class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors group">
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-primary-600 text-white flex items-center justify-center font-bold text-sm">
                  {{ (emp.first_name?.charAt(0) || '') + (emp.last_name?.charAt(0) || '') }}
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ emp.first_name }} {{ emp.last_name }}</p>
                  <p class="text-xs text-gray-500">{{ emp.email || emp.phone_number }}</p>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-gray-600 dark:text-gray-300">{{ emp.profession || '--' }}</td>
            <td class="px-6 py-4">
              <span :class="getStatusClass(emp.status)" 
                class="px-2.5 py-1 rounded-full text-xs font-medium border">
                {{ formatStatus(emp.status) }}
              </span>
            </td>
            <td class="px-6 py-4 text-right font-medium text-gray-900 dark:text-white">
              {{ formatCurrency(emp.salary_config?.net_payable || 0) }}
            </td>
            <td class="px-6 py-4 text-right">
              <div class="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click="viewEmployee(emp)" 
                  class="p-2 text-gray-400 hover:text-primary-600 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded-lg"
                  title="Voir détails">
                  <EyeIcon class="w-4 h-4" />
                </button>
                <button @click="editSalary(emp)" 
                  class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg"
                  title="Configurer salaire">
                  <CurrencyDollarIcon class="w-4 h-4" />
                </button>
                <button v-if="emp.status === 'ACTIVE'" @click="confirmTerminate(emp)" 
                  class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg"
                  title="Licencier">
                  <XCircleIcon class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Invite Modal -->
    <InviteEmployeeModal 
      v-if="showInviteModal" 
      :enterprise-id="enterpriseId"
      @close="showInviteModal = false"
      @invited="handleInvited" />

    <!-- Employee Detail Modal -->
    <EmployeeDetailModal 
      v-if="selectedEmployee" 
      :employee="selectedEmployee"
      @close="selectedEmployee = null"
      @updated="handleEmployeeUpdated" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { 
  UsersIcon, UserPlusIcon, UserGroupIcon, MagnifyingGlassIcon,
  EyeIcon, CurrencyDollarIcon, XCircleIcon
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  enterpriseId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['refresh'])

const employees = ref([])
const isLoading = ref(true)
const searchQuery = ref('')
const statusFilter = ref('')
const showInviteModal = ref(false)
const selectedEmployee = ref(null)

// Computed
const filteredEmployees = computed(() => {
  return employees.value.filter(emp => {
    const matchesSearch = !searchQuery.value || 
      `${emp.first_name} ${emp.last_name}`.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !statusFilter.value || emp.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const activeCount = computed(() => employees.value.filter(e => e.status === 'ACTIVE').length)
const pendingCount = computed(() => employees.value.filter(e => e.status === 'PENDING_INVITE').length)
const totalMasseSalariale = computed(() => 
  employees.value
    .filter(e => e.status === 'ACTIVE')
    .reduce((sum, e) => sum + (e.salary_config?.net_payable || 0), 0)
)

// Methods
const fetchEmployees = async () => {
  isLoading.value = true
  try {
    const { data } = await enterpriseAPI.listEmployees(props.enterpriseId)
    employees.value = data || []
  } catch (e) {
    console.error('Failed to fetch employees', e)
    employees.value = []
  } finally {
    isLoading.value = false
  }
}

const getStatusClass = (status) => {
  switch(status) {
    case 'ACTIVE': return 'bg-green-50 text-green-700 border-green-200 dark:bg-green-900/20 dark:text-green-400 dark:border-green-800'
    case 'PENDING_INVITE': return 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-900/20 dark:text-amber-400 dark:border-amber-800'
    case 'TERMINATED': return 'bg-red-50 text-red-700 border-red-200 dark:bg-red-900/20 dark:text-red-400 dark:border-red-800'
    default: return 'bg-gray-50 text-gray-700 border-gray-200'
  }
}

const formatStatus = (status) => {
  const map = {
    'ACTIVE': 'Actif',
    'PENDING_INVITE': 'En attente',
    'TERMINATED': 'Licencié'
  }
  return map[status] || status
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('fr-FR').format(amount) + ' XOF'
}

const viewEmployee = (emp) => {
  selectedEmployee.value = emp
}

const editSalary = (emp) => {
  selectedEmployee.value = { ...emp, _editingSalary: true }
}

const confirmTerminate = async (emp) => {
  if (confirm(`Êtes-vous sûr de vouloir licencier ${emp.first_name} ${emp.last_name} ?`)) {
    try {
      // TODO: Call API to terminate
      // await enterpriseAPI.terminateEmployee(props.enterpriseId, emp.id)
      await fetchEmployees()
    } catch (e) {
      alert('Erreur lors du licenciement')
    }
  }
}

const handleInvited = () => {
  showInviteModal.value = false
  fetchEmployees()
}

const handleEmployeeUpdated = () => {
  selectedEmployee.value = null
  fetchEmployees()
}

// Lifecycle
onMounted(fetchEmployees)

watch(() => props.enterpriseId, fetchEmployees)
</script>
