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
        <div class="flex items-center gap-3">
          <button v-if="isOwner && event?.status !== 'cancelled'" 
                  @click="confirmCancelEvent" 
                  class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2">
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
             Annuler l'événement
          </button>
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
          <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Total vendu</p>
            <p class="text-3xl font-bold text-gray-900 dark:text-white mt-1">{{ tickets.length }}</p>
          </div>
          <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Revenus</p>
            <p class="text-3xl font-bold text-emerald-600 mt-1">{{ formatAmount(totalRevenue) }} <span class="text-sm font-normal text-gray-500">{{ event?.currency }}</span></p>
          </div>
          <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Utilisés</p>
            <p class="text-3xl font-bold text-indigo-600 mt-1">{{ usedCount }}</p>
          </div>
          <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Remboursés</p>
            <p class="text-3xl font-bold text-red-600 mt-1">{{ refundedCount }}</p>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="tickets.length === 0" class="glass-card p-16 text-center bg-white dark:bg-slate-800/50 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="text-6xl mb-4">🎫</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun ticket vendu</h3>
          <p class="text-gray-500 dark:text-gray-400">Les tickets vendus apparaîtront ici.</p>
        </div>

        <!-- Tickets Table -->
        <div v-else class="glass-card overflow-hidden bg-white dark:bg-slate-800/50 border border-gray-100 dark:border-gray-700 shadow-sm rounded-xl">
          <table class="table-premium w-full">
            <thead>
              <tr class="bg-gray-50/80 dark:bg-slate-800/80 border-b border-gray-100 dark:border-gray-700">
                <th class="text-left px-6 py-4">Acheteur</th>
                <th class="text-left px-6 py-4">Type</th>
                <th class="text-left px-6 py-4">Code</th>
                <th class="text-left px-6 py-4">Prix</th>
                <th class="text-left px-6 py-4">Statut</th>
                <th class="text-left px-6 py-4">Date</th>
                <th class="text-right px-6 py-4">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ticket in tickets" :key="ticket.id" @click.stop="viewDetails(ticket)" class="border-t border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-slate-800/30 cursor-pointer transition-colors">
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
                 <td class="px-6 py-4 text-right" @click.stop>
                   <button v-if="ticket.status === 'paid' && isOwner" 
                           @click="confirmRefund(ticket)"
                           class="text-sm text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/30 px-3 py-1 rounded-md transition-colors border border-red-200 dark:border-red-800">
                     Rembourser
                   </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Details Modal -->
      <Teleport to="body">
        <div v-if="showDetailsModal && selectedTicket" class="fixed inset-0 z-40 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="showDetailsModal = false">
          <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200">
            <!-- Modal Header -->
            <div class="p-6 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center bg-gray-50/50 dark:bg-slate-800/50">
              <div>
                <h3 class="text-xl font-bold text-gray-900 dark:text-white">Détails du Ticket</h3>
                <p class="text-sm text-gray-500 dark:text-gray-400">Code: {{ selectedTicket.ticket_code }}</p>
              </div>
              <button @click="showDetailsModal = false" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>

            <!-- Modal Content -->
            <div class="p-6 space-y-6 max-h-[70vh] overflow-y-auto custom-scrollbar">
              
              <!-- Status & Amount -->
              <div class="flex justify-between items-center p-4 rounded-xl bg-gray-50 dark:bg-slate-800">
                <div>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Statut</p>
                  <span :class="getStatusClass(selectedTicket.status)" class="px-2 py-1 rounded-full text-xs font-medium inline-block mt-1">
                    {{ getStatusLabel(selectedTicket.status) }}
                  </span>
                </div>
                <div class="text-right">
                   <p class="text-sm text-gray-500 dark:text-gray-400">Prix payé</p>
                   <p class="text-lg font-bold text-gray-900 dark:text-white">{{ formatAmount(selectedTicket.price) }} {{ selectedTicket.currency }}</p>
                </div>
              </div>

               <!-- Action Buttons in Modal -->
              <div v-if="selectedTicket.status === 'paid' && isOwner" class="flex justify-end pt-2">
                 <button @click="confirmRefund(selectedTicket); showDetailsModal = false;" 
                         class="px-4 py-2 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 border border-red-200 dark:border-red-800 rounded-lg hover:bg-red-100 transition-colors w-full sm:w-auto">
                    Rembourser ce ticket
                 </button>
              </div>

              <!-- Ticket Info -->
              <div>
                <h4 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Information du Ticket</h4>
                <div class="space-y-3">
                  <div class="flex justify-between border-b border-gray-100 dark:border-gray-800 pb-2">
                    <span class="text-gray-600 dark:text-gray-300">Type de ticket</span>
                    <span class="font-medium text-gray-900 dark:text-white" :style="{ color: selectedTicket.tier_color }">{{ selectedTicket.tier_icon }} {{ selectedTicket.tier_name }}</span>
                  </div>
                  <div class="flex justify-between border-b border-gray-100 dark:border-gray-800 pb-2">
                    <span class="text-gray-600 dark:text-gray-300">Date d'achat</span>
                    <span class="font-medium text-gray-900 dark:text-white">{{ formatDate(selectedTicket.created_at) }}</span>
                  </div>
                </div>
              </div>

              <!-- Participant Data (Form Fields) -->
              <div v-if="selectedTicket.form_data && Object.keys(selectedTicket.form_data).length > 0">
                <h4 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Informations du Participant</h4>
                <div class="space-y-3 bg-gray-50 dark:bg-slate-800/50 p-4 rounded-xl">
                  <div v-for="(value, key) in selectedTicket.form_data" :key="key" class="flex flex-col sm:flex-row sm:justify-between border-b border-gray-200 dark:border-gray-700 last:border-0 pb-2 last:pb-0">
                    <span class="text-sm text-gray-500 dark:text-gray-400 mb-1 sm:mb-0">{{ getLabelForField(key) }}</span>
                    <span class="font-medium text-gray-900 dark:text-white text-right break-words">{{ value }}</span>
                  </div>
                </div>
              </div>

              <!-- Transaction ID -->
              <div v-if="selectedTicket.transaction_id" class="text-center pt-2">
                <p class="text-xs text-gray-400 font-mono">Ref: {{ selectedTicket.transaction_id }}</p>
              </div>

            </div>

            <!-- Footer -->
            <div class="p-4 border-t border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-slate-800/50 flex justify-end">
              <button @click="showDetailsModal = false" class="px-4 py-2 bg-white dark:bg-slate-700 border border-gray-200 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors">
                Fermer
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'auth'
})

import { ticketAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const eventId = computed(() => route.params.id as string)

const loading = ref(true)
const error = ref('')
const event = ref<any>(null)
const tickets = ref<any[]>([])

const totalRevenue = computed(() => tickets.value.reduce((sum, t) => sum + (t.status !== 'refunded' ? (t.price || 0) : 0), 0))
const usedCount = computed(() => tickets.value.filter(t => t.status === 'used').length)
const pendingCount = computed(() => tickets.value.filter(t => t.status === 'pending').length)
const refundedCount = computed(() => tickets.value.filter(t => t.status === 'refunded').length)

const isOwner = computed(() => {
  return event.value?.creator_id === authStore.user?.id
})

const loadData = async () => {
  loading.value = true
  error.value = ''
  try {
    const [eventRes, ticketsRes] = await Promise.all([
      ticketAPI.getEvent(eventId.value),
      ticketAPI.getEventTickets(eventId.value, 100, 0)
    ])
    event.value = eventRes.data?.event || eventRes.data
    tickets.value = ticketsRes.data?.tickets || ticketsRes.data || []
  } catch (err: any) {
    console.error('Failed to load data:', err)
    error.value = err.response?.data?.error || 'Erreur lors du chargement des tickets'
  } finally {
    loading.value = false
  }
}

const confirmRefund = async (ticket: any) => {
    if(!confirm(`Êtes-vous sûr de vouloir rembourser le ticket ${ticket.ticket_code} ? Cette action est irréversible et les fonds seront renvoyés à l'acheteur.`)) {
        return;
    }
    
    try {
        await ticketAPI.refundTicket(ticket.id)
        // Optimistic update
        const t = tickets.value.find(t => t.id === ticket.id)
        if(t) t.status = 'refunded'
        alert('Ticket remboursé avec succès')
    } catch(err: any) {
        alert(err.response?.data?.error || 'Erreur lors du remboursement')
    }
}

const confirmCancelEvent = async () => {
    if(!confirm("ATTENTION : Êtes-vous sûr de vouloir annuler cet événement ? TOUS les tickets vendus seront REMBOURSÉS automatiquement. Cette action est IRRÉVERSIBLE.")) {
        return;
    }
     if(!confirm("Dernière vérification : Voulez-vous vraiment procéder à l'annulation et au remboursement général ?")) {
        return;
    }

    try {
        await ticketAPI.cancelEvent(eventId.value)
        await loadData() // Reload to see updates
        alert('Événement annulé et remboursements initiés.')
    } catch(err: any) {
        alert(err.response?.data?.error || "Erreur lors de l'annulation")
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
    case 'paid': return 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400'
    case 'cancelled': return 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
    case 'refunded': return 'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 line-through'
    default: return 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'used': return '✓ Utilisé'
    case 'paid': return 'Confirmé'
    case 'cancelled': return 'Annulé'
    case 'pending': return 'En attente'
    case 'refunded': return 'Remboursé'
    default: return status
  }
}

const showDetailsModal = ref(false)
const selectedTicket = ref<any>(null)

const viewDetails = (ticket: any) => {
  selectedTicket.value = ticket
  showDetailsModal.value = true
}

const getLabelForField = (fieldName: string) => {
  if (!event.value?.form_fields) return fieldName
  const field = event.value.form_fields.find((f: any) => f.name === fieldName)
  return field ? field.label : fieldName
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

