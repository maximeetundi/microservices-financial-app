<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <UserGroupIcon class="w-6 h-6 text-purple-500" />
          Gestion des Abonnés
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Clients et élèves inscrits à vos services
        </p>
      </div>
      <div class="flex gap-3">
        <button @click="downloadExport" class="px-4 py-2 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors flex items-center gap-2">
          <ArrowDownTrayIcon class="w-4 h-4" />
          Exporter
        </button>
        <button @click="showAddModal = true" 
          class="px-5 py-2.5 bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-700 hover:to-purple-800 text-white rounded-xl text-sm font-semibold shadow-lg shadow-purple-500/25 transition-all flex items-center gap-2">
          <PlusIcon class="w-5 h-5" />
          Nouveau Client
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-3">
      <select v-model="selectedService" @change="fetchClients"
        class="px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500">
        <option value="">Tous les services</option>
        <option v-for="svc in allServices" :key="svc.id" :value="svc.id">{{ svc.name }}</option>
      </select>
    </div>

    <!-- Clients Table -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div v-if="isLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
      </div>

      <div v-else-if="clients.length === 0" class="text-center py-16">
        <UsersIcon class="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <h4 class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">Aucun abonné</h4>
        <p class="text-gray-500 text-sm">Les clients apparaîtront ici après inscription</p>
      </div>

      <table v-else class="w-full text-left text-sm">
        <thead class="bg-gray-50 dark:bg-gray-900/50 text-gray-600 dark:text-gray-400 text-xs uppercase">
          <tr>
            <th class="px-6 py-4 font-semibold">Client</th>
            <th class="px-6 py-4 font-semibold">Matricule</th>
            <th class="px-6 py-4 font-semibold">Service</th>
            <th class="px-6 py-4 font-semibold text-right">Montant</th>
            <th class="px-6 py-4 font-semibold">Prochaine Facture</th>
            <th class="px-6 py-4 font-semibold text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
          <tr v-for="sub in clients" :key="sub.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
            <td class="px-6 py-4 font-medium text-gray-900 dark:text-white">{{ sub.client_name }}</td>
            <td class="px-6 py-4">
              <span v-if="sub.external_id" class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-xs font-mono text-gray-600 dark:text-gray-300">
                {{ sub.external_id }}
              </span>
              <span v-else class="text-gray-400">--</span>
            </td>
            <td class="px-6 py-4 text-gray-600 dark:text-gray-300">{{ getServiceName(sub.service_id) }}</td>
            <td class="px-6 py-4 text-right font-medium text-gray-900 dark:text-white">
              {{ formatCurrency(sub.amount) }}
              <span class="text-xs text-gray-400 block">{{ formatFrequency(sub.billing_frequency) }}</span>
            </td>
            <td class="px-6 py-4 text-gray-500">{{ formatDate(sub.next_billing_at) }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="confirmCancel(sub)" class="text-red-500 hover:text-red-700 text-xs font-medium hover:underline">
                Résilier
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Client Modal -->
    <AddClientModal 
      v-if="showAddModal"
      :enterprise="enterprise"
      @close="showAddModal = false"
      @added="handleClientAdded" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { UserGroupIcon, PlusIcon, ArrowDownTrayIcon, UsersIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const props = defineProps({
  enterprise: { type: Object, required: true }
})

const { requirePin, checkPinStatus } = usePin()

const clients = ref([])
const isLoading = ref(true)
const selectedService = ref('')
const showAddModal = ref(false)

// Check PIN status on mount
onMounted(async () => {
  await checkPinStatus()
  fetchClients()
})

const allServices = computed(() => {
  return (props.enterprise.service_groups || []).flatMap(g => g.services || [])
})

const getServiceName = (id) => {
  const svc = allServices.value.find(s => s.id === id)
  return svc?.name || id
}

const fetchClients = async () => {
  isLoading.value = true
  try {
    const { data } = await enterpriseAPI.getSubscriptions(props.enterprise.id, selectedService.value)
    clients.value = (data || []).sort((a, b) => (a.client_name || '').localeCompare(b.client_name || ''))
  } catch (e) {
    console.error('Failed to fetch clients', e)
  } finally {
    isLoading.value = false
  }
}

const formatCurrency = (amount) => new Intl.NumberFormat('fr-FR').format(amount || 0) + ' XOF'
const formatFrequency = (freq) => {
  const map = { MONTHLY: 'Mensuel', WEEKLY: 'Hebdo', ANNUALLY: 'Annuel', ONETIME: 'Unique' }
  return map[freq] || freq
}
const formatDate = (date) => date ? new Date(date).toLocaleDateString() : '--'

const downloadExport = () => {
  if (!clients.value.length) return
  let csv = 'Nom,Service,Matricule,Montant,Fréquence,Prochaine Facture\n'
  clients.value.forEach(sub => {
    csv += `"${sub.client_name}","${getServiceName(sub.service_id)}","${sub.external_id || ''}",${sub.amount},"${sub.billing_frequency}","${formatDate(sub.next_billing_at)}"\n`
  })
  const blob = new Blob([csv], { type: 'text/csv' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `clients_${props.enterprise.name}_${new Date().toISOString().split('T')[0]}.csv`
  a.click()
}

// 🔒 Cancel subscription with PIN verification
const confirmCancel = async (sub) => {
  if (!confirm(`Résilier l'abonnement de ${sub.client_name} ?`)) return
  
  await requirePin(async () => {
    try {
      await enterpriseAPI.cancelSubscription(props.enterprise.id, sub.id)
      fetchClients()
      alert('Abonnement résilié avec succès')
    } catch (e) {
      alert('Erreur: ' + e.message)
    }
  })
}

const handleClientAdded = () => {
  showAddModal.value = false
  fetchClients()
}
</script>
