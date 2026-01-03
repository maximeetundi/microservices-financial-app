<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto py-6 px-4">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div class="flex items-center gap-4">
          <button @click="navigateTo(`/events/${eventId}`)" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors">
            <svg class="w-6 h-6 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
          </button>
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">🎟️ Tickets Vendus</h1>
            <p class="text-gray-500 dark:text-gray-400">{{ event?.title || 'Chargement...' }}</p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="px-3 py-1 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-sm font-medium">
            {{ tickets.length }} tickets
          </span>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="loading-spinner"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="glass-card p-8 text-center">
        <p class="text-red-500 dark:text-red-400 mb-4">{{ error }}</p>
        <button @click="loadData" class="btn-premium">Réessayer</button>
      </div>

      <!-- Tickets List -->
      <div v-else>
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div class="glass-card p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Total vendu</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ tickets.length }}</p>
          </div>
          <div class="glass-card p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Revenus</p>
            <p class="text-2xl font-bold text-emerald-600">{{ formatAmount(totalRevenue) }} {{ event?.currency }}</p>
          </div>
          <div class="glass-card p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Utilisés</p>
            <p class="text-2xl font-bold text-indigo-600">{{ usedCount }}</p>
          </div>
          <div class="glass-card p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">En attente</p>
            <p class="text-2xl font-bold text-amber-600">{{ pendingCount }}</p>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="tickets.length === 0" class="glass-card p-16 text-center">
          <div class="text-6xl mb-4">🎫</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun ticket vendu</h3>
          <p class="text-gray-500 dark:text-gray-400">Les tickets vendus apparaîtront ici.</p>
        </div>

        <!-- Tickets Table -->
        <div v-else class="glass-card overflow-hidden">
          <table class="table-premium w-full">
            <thead>
              <tr class="bg-gray-50 dark:bg-slate-800/50">
                <th class="text-left px-6 py-4">Acheteur</th>
                <th class="text-left px-6 py-4">Type</th>
                <th class="text-left px-6 py-4">Code</th>
                <th class="text-left px-6 py-4">Prix</th>
                <th class="text-left px-6 py-4">Statut</th>
                <th class="text-left px-6 py-4">Date</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ticket in tickets" :key="ticket.id" class="border-t border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-slate-800/30">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400 font-bold">
                      {{ getInitials(ticket.form_data) }}
                    </div>
                    <div>
                      <p class="font-medium text-gray-900 dark:text-white">{{ getBuyerName(ticket.form_data) }}</p>
                      <p class="text-sm text-gray-500 dark:text-gray-400">{{ getBuyerEmail(ticket.form_data) }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span class="inline-flex items-center gap-1 px-2 py-1 rounded-full text-sm" 
                        :style="{ backgroundColor: ticket.tier_color + '20', color: ticket.tier_color }">
                    {{ ticket.tier_icon }} {{ ticket.tier_name }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <code class="px-2 py-1 bg-gray-100 dark:bg-slate-700 rounded text-sm font-mono text-gray-700 dark:text-gray-300">
                    {{ ticket.ticket_code }}
                  </code>
                </td>
                <td class="px-6 py-4 font-semibold text-gray-900 dark:text-white">
                  {{ formatAmount(ticket.price) }} {{ ticket.currency }}
                </td>
                <td class="px-6 py-4">
                  <span :class="getStatusClass(ticket.status)" class="px-2 py-1 rounded-full text-xs font-medium">
                    {{ getStatusLabel(ticket.status) }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-gray-500 dark:text-gray-400">
                  {{ formatDate(ticket.created_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'auth'
})

const route = useRoute()
const { ticketApi } = useApi()

const eventId = computed(() => route.params.id as string)

const loading = ref(true)
const error = ref('')
const event = ref<any>(null)
const tickets = ref<any[]>([])

const totalRevenue = computed(() => tickets.value.reduce((sum, t) => sum + (t.price || 0), 0))
const usedCount = computed(() => tickets.value.filter(t => t.status === 'used').length)
const pendingCount = computed(() => tickets.value.filter(t => t.status === 'pending' || t.status === 'confirmed').length)

const loadData = async () => {
  loading.value = true
  error.value = ''
  try {
    const [eventRes, ticketsRes] = await Promise.all([
      ticketApi.getEvent(eventId.value),
      ticketApi.getEventTickets(eventId.value, 100, 0)
    ])
    event.value = eventRes.data
    tickets.value = ticketsRes.data?.tickets || ticketsRes.data || []
  } catch (err: any) {
    console.error('Failed to load data:', err)
    error.value = err.response?.data?.error || 'Erreur lors du chargement des tickets'
  } finally {
    loading.value = false
  }
}

const getInitials = (formData: any) => {
  if (!formData) return '?'
  const name = formData.name || formData.nom || formData.full_name || ''
  return name.split(' ').map((n: string) => n[0]).join('').toUpperCase().slice(0, 2) || '?'
}

const getBuyerName = (formData: any) => {
  if (!formData) return 'Anonyme'
  return formData.name || formData.nom || formData.full_name || 'Anonyme'
}

const getBuyerEmail = (formData: any) => {
  if (!formData) return ''
  return formData.email || formData.Email || ''
}

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('fr-FR').format(amount || 0)
}

const formatDate = (date: string) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('fr-FR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'used': return 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400'
    case 'confirmed': return 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400'
    case 'cancelled': return 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
    default: return 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'used': return '✓ Utilisé'
    case 'confirmed': return 'Confirmé'
    case 'cancelled': return 'Annulé'
    case 'pending': return 'En attente'
    default: return status
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
