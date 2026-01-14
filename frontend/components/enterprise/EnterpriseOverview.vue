<template>
  <div class="space-y-8">
    <!-- Welcome Banner -->
    <div class="bg-gradient-to-r from-indigo-600 via-purple-600 to-primary-600 rounded-2xl p-8 text-white shadow-xl relative overflow-hidden">
      <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4xIj48cGF0aCBkPSJNMzYgMzRjMC0yIDItNCAyLTRzMiAyIDIgNGMwIDItMiA0LTIgNHMtMi0yLTItNHoiLz48L2c+PC9nPjwvc3ZnPg==')] opacity-20"></div>
      
      <div class="relative z-10 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h2 class="text-2xl md:text-3xl font-bold mb-2">Bienvenue sur {{ enterprise.name }}</h2>
          <p class="text-white/80 max-w-xl">
            Accédez rapidement à vos outils de gestion. Configurez vos services, gérez votre personnel et suivez vos encaissements en temps réel.
          </p>
        </div>
        <div class="flex gap-3">
          <button @click="$emit('navigate', 'Settings')" 
            class="px-4 py-2.5 bg-white/20 hover:bg-white/30 backdrop-blur-sm rounded-xl text-sm font-medium transition-all flex items-center gap-2 border border-white/20">
            <Cog6ToothIcon class="w-5 h-5" />
            Paramètres
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-sm border border-gray-100 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-blue-100 dark:bg-blue-900/30 text-blue-600 flex items-center justify-center">
            <UsersIcon class="w-5 h-5" />
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.employees }}</p>
            <p class="text-xs text-gray-500">Employés</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-sm border border-gray-100 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-purple-100 dark:bg-purple-900/30 text-purple-600 flex items-center justify-center">
            <UserGroupIcon class="w-5 h-5" />
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.clients }}</p>
            <p class="text-xs text-gray-500">Abonnés</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-sm border border-gray-100 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-green-100 dark:bg-green-900/30 text-green-600 flex items-center justify-center">
            <BoltIcon class="w-5 h-5" />
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.services }}</p>
            <p class="text-xs text-gray-500">Services</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-sm border border-gray-100 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-amber-100 dark:bg-amber-900/30 text-amber-600 flex items-center justify-center">
            <BanknotesIcon class="w-5 h-5" />
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ formatCurrency(stats.revenue) }}</p>
            <p class="text-xs text-gray-500">Ce mois</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Add Service Action -->
      <div @click="$emit('navigate', 'Settings')" 
        class="group bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-lg transition-all cursor-pointer border border-gray-100 dark:border-gray-700 relative overflow-hidden">
        <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-blue-100 to-transparent dark:from-blue-900/20 rounded-bl-full"></div>
        <div class="relative z-10">
          <div class="w-14 h-14 bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-blue-500/30 group-hover:scale-110 transition-transform">
            <BoltIcon class="w-7 h-7" />
          </div>
          <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Nouveau Service</h3>
          <p class="text-sm text-gray-500">Ajoutez une prestation, une classe ou un abonnement.</p>
        </div>
      </div>

      <!-- Add Member Action -->
      <div @click="$emit('navigate', 'Employees')" 
        class="group bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-lg transition-all cursor-pointer border border-gray-100 dark:border-gray-700 relative overflow-hidden">
        <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-green-100 to-transparent dark:from-green-900/20 rounded-bl-full"></div>
        <div class="relative z-10">
          <div class="w-14 h-14 bg-gradient-to-br from-green-500 to-green-600 text-white rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-green-500/30 group-hover:scale-110 transition-transform">
            <UserPlusIcon class="w-7 h-7" />
          </div>
          <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Inviter Membre</h3>
          <p class="text-sm text-gray-500">Ajoutez des employés à votre équipe.</p>
        </div>
      </div>

      <!-- Billing Action -->
      <div @click="$emit('navigate', 'Billing')" 
        class="group bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-lg transition-all cursor-pointer border border-gray-100 dark:border-gray-700 relative overflow-hidden">
        <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-purple-100 to-transparent dark:from-purple-900/20 rounded-bl-full"></div>
        <div class="relative z-10">
          <div class="w-14 h-14 bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-purple-500/30 group-hover:scale-110 transition-transform">
            <DocumentTextIcon class="w-7 h-7" />
          </div>
          <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2">Facturation</h3>
          <p class="text-sm text-gray-500">Générez des factures ou saisissez des consommations.</p>
        </div>
      </div>
    </div>

    <!-- Bottom Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Quick Stats -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
        <h3 class="font-bold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
          <ChartBarIcon class="w-5 h-5 text-primary-500" />
          Aperçu Rapide
        </h3>
        <div class="grid grid-cols-2 gap-4">
          <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
            <span class="text-sm text-gray-500">Groupes de services</span>
            <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">{{ enterprise.service_groups?.length || 0 }}</p>
          </div>
          <div class="p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
            <span class="text-sm text-gray-500">Services actifs</span>
            <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">{{ totalServices }}</p>
          </div>
        </div>
      </div>

      <!-- QR Code Preview -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700 flex items-center justify-between">
        <div>
          <h3 class="font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <QrCodeIcon class="w-5 h-5 text-primary-500" />
            QR Code Public
          </h3>
          <p class="text-sm text-gray-500 mt-1 max-w-xs">Permettez à vos clients de s'abonner en scannant votre code unique.</p>
          <button @click="$emit('show-qr')" 
            class="mt-4 text-primary-600 font-medium hover:text-primary-700 text-sm flex items-center gap-1 group">
            Voir les codes 
            <ArrowRightIcon class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
          </button>
        </div>
        <div class="w-24 h-24 bg-white p-2 rounded-xl border border-gray-200 shadow-sm flex-shrink-0">
          <img v-if="qrCodeUrl" :src="qrCodeUrl" class="w-full h-full object-contain" />
          <div v-else class="w-full h-full flex items-center justify-center">
            <QrCodeIcon class="w-12 h-12 text-gray-300" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { 
  UsersIcon, UserGroupIcon, BoltIcon, BanknotesIcon, UserPlusIcon,
  DocumentTextIcon, ChartBarIcon, QrCodeIcon, ArrowRightIcon, Cog6ToothIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  enterprise: {
    type: Object,
    required: true
  },
  stats: {
    type: Object,
    default: () => ({ employees: 0, clients: 0, services: 0, revenue: 0 })
  }
})

const emit = defineEmits(['navigate', 'show-qr'])

const totalServices = computed(() => {
  return (props.enterprise.service_groups || []).reduce((sum, g) => sum + (g.services?.length || 0), 0)
})

const qrCodeUrl = computed(() => {
  if (!props.enterprise?.id) return ''
  return `/enterprise-service/api/v1/enterprises/${props.enterprise.id}/qrcode`
})

const formatCurrency = (amount) => {
  if (!amount) return '--'
  return new Intl.NumberFormat('fr-FR', { notation: 'compact' }).format(amount)
}
</script>
