<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <BanknotesIcon class="w-6 h-6 text-green-500" />
          Gestion de la Paie
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Exécutez et suivez les versements de salaires
        </p>
      </div>
      <button @click="initiateRunPayroll" 
        class="px-5 py-2.5 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 text-white rounded-xl text-sm font-semibold shadow-lg shadow-green-500/25 transition-all flex items-center gap-2">
        <PlayIcon class="w-5 h-5" />
        Exécuter la Paie
      </button>
    </div>

    <!-- Period Selector -->
    <div class="flex items-center gap-4 p-4 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700">
      <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Période:</label>
      <select v-model="selectedMonth" class="px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700">
        <option v-for="m in 12" :key="m" :value="m">{{ monthNames[m-1] }}</option>
      </select>
      <select v-model="selectedYear" class="px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700">
        <option v-for="y in years" :key="y" :value="y">{{ y }}</option>
      </select>
      <button @click="fetchPayrollRuns" class="px-3 py-2 bg-gray-100 dark:bg-gray-700 rounded-lg text-sm hover:bg-gray-200 dark:hover:bg-gray-600">
        Actualiser
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Total Employés</p>
        <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalEmployees }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Masse Salariale</p>
        <p class="text-2xl font-bold text-green-600">{{ formatCurrency(stats.totalAmount) }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">Versements ce mois</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.paidCount }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-100 dark:border-gray-700">
        <p class="text-sm text-gray-500">En attente</p>
        <p class="text-2xl font-bold text-amber-600">{{ stats.pendingCount }}</p>
      </div>
    </div>

    <!-- Payroll Runs History -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-4 border-b border-gray-100 dark:border-gray-700">
        <h4 class="font-semibold text-gray-900 dark:text-white">Historique des Exécutions</h4>
      </div>

      <div v-if="isLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
      </div>

      <div v-else-if="payrollRuns.length === 0" class="text-center py-12 text-gray-500">
        <CalendarIcon class="w-12 h-12 mx-auto mb-3 opacity-50" />
        <p>Aucune exécution de paie pour cette période</p>
      </div>

      <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
        <div v-for="run in payrollRuns" :key="run.id" 
          class="p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer"
          @click="selectedRun = run">
          <div class="flex items-center gap-4">
            <div :class="getStatusBadgeClass(run.status)" class="w-10 h-10 rounded-xl flex items-center justify-center">
              <CheckIcon v-if="run.status === 'COMPLETED'" class="w-5 h-5" />
              <ClockIcon v-else-if="run.status === 'PROCESSING'" class="w-5 h-5 animate-spin" />
              <XMarkIcon v-else class="w-5 h-5" />
            </div>
            <div>
              <p class="font-medium text-gray-900 dark:text-white">
                {{ monthNames[run.period_month - 1] }} {{ run.period_year }}
              </p>
              <p class="text-sm text-gray-500">
                {{ run.total_employees }} employés • {{ formatCurrency(run.total_amount) }}
              </p>
            </div>
          </div>
          <div class="text-right">
            <span :class="getStatusClass(run.status)" class="px-2.5 py-1 rounded-full text-xs font-medium">
              {{ formatStatus(run.status) }}
            </span>
            <p class="text-xs text-gray-400 mt-1">{{ formatDate(run.executed_at) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Run Detail Modal -->
    <PayrollDetailModal 
      v-if="selectedRun" 
      :run="selectedRun"
      @close="selectedRun = null" />

    <!-- Run Payroll Modal -->
    <RunPayrollModal 
      v-if="showRunPayroll"
      :enterprise-id="enterpriseId"
      :month="selectedMonth"
      :year="selectedYear"
      @close="showRunPayroll = false"
      @executed="handlePayrollExecuted" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  BanknotesIcon, PlayIcon, CalendarIcon, CheckIcon, ClockIcon, XMarkIcon
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const props = defineProps({
  enterpriseId: { type: String, required: true }
})

const { requirePin, checkPinStatus } = usePin()

const payrollRuns = ref([])
const isLoading = ref(true)
const selectedRun = ref(null)
const showRunPayroll = ref(false)

const currentDate = new Date()
const selectedMonth = ref(currentDate.getMonth() + 1)
const selectedYear = ref(currentDate.getFullYear())

const monthNames = ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Septembre', 'Octobre', 'Novembre', 'Décembre']
const years = computed(() => {
  const curr = currentDate.getFullYear()
  return [curr - 1, curr, curr + 1]
})

// Check PIN status on mount
onMounted(async () => {
  await checkPinStatus()
  fetchPayrollRuns()
})

const stats = computed(() => {
  const runs = payrollRuns.value.filter(r => r.period_month === selectedMonth.value && r.period_year === selectedYear.value)
  const completed = runs.filter(r => r.status === 'COMPLETED')
  return {
    totalEmployees: completed.reduce((s, r) => s + (r.total_employees || 0), 0),
    totalAmount: completed.reduce((s, r) => s + (r.total_amount || 0), 0),
    paidCount: completed.length,
    pendingCount: runs.filter(r => r.status === 'PROCESSING' || r.status === 'DRAFT').length
  }
})

const fetchPayrollRuns = async () => {
  isLoading.value = true
  try {
    const { data } = await enterpriseAPI.getPayrollRuns(props.enterpriseId, selectedYear.value)
    payrollRuns.value = data || []
  } catch (e) {
    console.error('Failed to fetch payroll runs', e)
    payrollRuns.value = []
  } finally {
    isLoading.value = false
  }
}

const formatCurrency = (amount) => new Intl.NumberFormat('fr-FR').format(amount || 0) + ' XOF'
const formatDate = (date) => date ? new Date(date).toLocaleString() : '--'
const formatStatus = (status) => {
  const map = { DRAFT: 'Brouillon', PROCESSING: 'En cours', COMPLETED: 'Complété', FAILED: 'Échoué' }
  return map[status] || status
}

const getStatusClass = (status) => {
  switch (status) {
    case 'COMPLETED': return 'bg-green-100 text-green-700'
    case 'PROCESSING': return 'bg-blue-100 text-blue-700'
    case 'FAILED': return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'COMPLETED': return 'bg-green-100 text-green-600'
    case 'PROCESSING': return 'bg-blue-100 text-blue-600'
    case 'FAILED': return 'bg-red-100 text-red-600'
    default: return 'bg-gray-100 text-gray-600'
  }
}

// 🔒 Run Payroll with PIN verification
const initiateRunPayroll = async () => {
  await requirePin(async () => {
    showRunPayroll.value = true
  })
}

const handlePayrollExecuted = () => {
  showRunPayroll.value = false
  fetchPayrollRuns()
}
</script>
