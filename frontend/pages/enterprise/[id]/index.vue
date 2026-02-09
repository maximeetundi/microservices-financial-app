<template>
  <NuxtLayout name="enterprise">
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Tableau de Bord</h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">Vue d'ensemble de {{ enterprise?.name }}</p>
        </div>
        <button 
          @click="showQRModal = true" 
          class="inline-flex items-center gap-2 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-colors"
        >
          <QrCodeIcon class="w-5 h-5" />
          Codes QR
        </button>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-xl bg-blue-100 dark:bg-blue-900/30">
              <UsersIcon class="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.employees }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Employés</p>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-xl bg-purple-100 dark:bg-purple-900/30">
              <UserGroupIcon class="w-6 h-6 text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.clients }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Abonnés</p>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-xl bg-amber-100 dark:bg-amber-900/30">
              <FolderIcon class="w-6 h-6 text-amber-600 dark:text-amber-400" />
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.services }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Services</p>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-xl bg-green-100 dark:bg-green-900/30">
              <BanknotesIcon class="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ formatCurrency(stats.revenue) }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Revenus</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Actions Rapides</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/employees`"
            class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group"
          >
            <div class="p-3 rounded-xl bg-blue-100 dark:bg-blue-900/30 group-hover:scale-110 transition-transform">
              <UsersIcon class="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Gérer Employés</span>
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/payroll`"
            class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group"
          >
            <div class="p-3 rounded-xl bg-green-100 dark:bg-green-900/30 group-hover:scale-110 transition-transform">
              <BanknotesIcon class="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Lancer Paie</span>
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/services`"
            class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group"
          >
            <div class="p-3 rounded-xl bg-amber-100 dark:bg-amber-900/30 group-hover:scale-110 transition-transform">
              <FolderIcon class="w-6 h-6 text-amber-600 dark:text-amber-400" />
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Configurer Services</span>
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/settings`"
            class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors group"
          >
            <div class="p-3 rounded-xl bg-gray-100 dark:bg-gray-700 group-hover:scale-110 transition-transform">
              <Cog6ToothIcon class="w-6 h-6 text-gray-600 dark:text-gray-400" />
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Paramètres</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Activité Récente</h2>
        <div v-if="recentActivity.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
          Aucune activité récente
        </div>
        <div v-else class="space-y-3">
          <div v-for="activity in recentActivity" :key="activity.id" 
            class="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
          >
            <div class="p-2 rounded-lg" :class="getActivityColor(activity.type)">
              <component :is="getActivityIcon(activity.type)" class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ activity.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ activity.description }}</p>
            </div>
            <span class="text-xs text-gray-400">{{ formatTime(activity.timestamp) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- QR Code Modal -->
    <QRCodeModal 
      :is-open="showQRModal" 
      :enterprise="enterprise"
      :enterprise-id="String(enterpriseId || '')"
      @close="showQRModal = false" 
    />
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, inject, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { enterpriseAPI } from '@/composables/useApi'
import { 
  UsersIcon, 
  UserGroupIcon, 
  FolderIcon, 
  BanknotesIcon, 
  Cog6ToothIcon,
  QrCodeIcon,
  ArrowPathIcon,
  CheckCircleIcon,
  ExclamationCircleIcon
} from '@heroicons/vue/24/outline'
import QRCodeModal from '@/components/enterprise/QRCodeModal.vue'

const route = useRoute()
const enterpriseId = computed(() => route.params.id)
const enterprise = inject('enterprise', ref(null))
const showQRModal = ref(false)

const stats = ref({
  employees: 0,
  clients: 0,
  services: 0,
  revenue: 0
})

const recentActivity = ref([])

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'XAF',
    minimumFractionDigits: 0
  }).format(amount || 0)
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

const getActivityIcon = (type) => {
  const icons = {
    employee: UsersIcon,
    payment: BanknotesIcon,
    service: FolderIcon,
    default: CheckCircleIcon
  }
  return icons[type] || icons.default
}

const getActivityColor = (type) => {
  const colors = {
    employee: 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400',
    payment: 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400',
    service: 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400',
    default: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400'
  }
  return colors[type] || colors.default
}

onMounted(async () => {
  if (enterpriseId.value) {
    try {
      // Fetch employees count
      const { data: employees } = await enterpriseAPI.listEmployees(enterpriseId.value)
      stats.value.employees = employees?.length || 0

      // Fetch clients count
      const { data: clients } = await enterpriseAPI.getSubscriptions(enterpriseId.value)
      stats.value.clients = clients?.length || 0

      // Calculate services count from enterprise
      if (enterprise.value?.service_groups) {
        stats.value.services = enterprise.value.service_groups.reduce(
          (sum, g) => sum + (g.services?.length || 0), 0
        )
      }
    } catch (e) {
      console.error('Failed to fetch stats:', e)
    }
  }
})
</script>
