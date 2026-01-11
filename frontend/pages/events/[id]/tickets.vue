<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto py-6 px-4">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
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
        <div class="flex flex-wrap items-center gap-3">
          <button v-if="isOwner && event?.status !== 'cancelled'" 
                  @click="openCancelModal" 
                  class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2 shadow-lg shadow-red-500/20">
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
             Annuler l'événement
          </button>
        </div>
      </div>

      <!-- Stats Cards -->
      <div v-if="!loading && event" class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
          <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Total vendu</p>
          <p class="text-3xl font-bold text-gray-900 dark:text-white mt-1">{{ stats?.total_sold || 0 }}</p>
        </div>
        <div class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700">
          <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">Revenus</p>
          <p class="text-3xl font-bold text-emerald-600 mt-1">{{ formatAmount(stats?.total_revenue || 0) }} <span class="text-sm font-normal text-gray-500">{{ event?.currency }}</span></p>
        </div>
        
        <!-- Tier Stats Loop -->
        <div v-for="tier in tierStats" :key="tier.id" class="glass-card p-4 bg-white dark:bg-slate-800/50 shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden">
             <div class="absolute right-0 top-0 p-2 opacity-10">
                 <span class="text-4xl">{{ tier.icon }}</span>
             </div>
             <p class="text-sm text-gray-500 dark:text-gray-400 font-medium truncate">{{ tier.name }}</p>
             <div class="flex items-baseline gap-1 mt-1">
                 <p class="text-2xl font-bold text-indigo-600">{{ tier.sold }}</p>
                 <span class="text-xs text-gray-400">/ {{ tier.capacity }} vendus</span>
             </div>
             <div class="w-full bg-gray-200 dark:bg-gray-700 h-1 mt-2 rounded-full overflow-hidden">
                 <div class="h-full bg-indigo-500 rounded-full" :style="{ width: `${(tier.sold / tier.capacity) * 100}%` }"></div>
             </div>
        </div>
      </div>

      <!-- Search Bar -->
      <div class="mb-6">
        <div class="relative max-w-md">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Rechercher (Nom, Tel, Code)..." 
            class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-shadow"
          >
          <svg class="w-5 h-5 text-gray-400 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="loading-spinner"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="glass-card p-8 text-center bg-white dark:bg-slate-800">
        <p class="text-red-500 dark:text-red-400 mb-4">{{ error }}</p>
        <button @click="loadData" class="btn-premium">Réessayer</button>
      </div>

      <!-- Tickets List -->
      <div v-else>
        <!-- Empty State -->
        <div v-if="tickets.length === 0" class="glass-card p-16 text-center bg-white dark:bg-slate-800 shadow-sm border border-gray-100 dark:border-gray-700">
          <div class="text-6xl mb-4">🎫</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun ticket vendu</h3>
          <p class="text-gray-500 dark:text-gray-400">Les tickets vendus apparaîtront ici.</p>
        </div>

        <!-- Tickets Table -->
        <div v-else class="glass-card bg-white dark:bg-slate-800 shadow-sm border border-gray-100 dark:border-gray-700 rounded-xl flex flex-col overflow-hidden">
          <div class="overflow-x-auto">
              <table class="table-premium w-full text-left whitespace-nowrap">
                <thead>
                  <tr class="bg-gray-50 dark:bg-slate-700/50 border-b border-gray-200 dark:border-gray-700">
                    <th class="px-6 py-4 text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Acheteur / Contact</th>
                    <th class="px-6 py-4 text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Type</th>
                    <th class="px-6 py-4 text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">Statut</th>
                    <th class="px-6 py-4 text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider text-right">Prix</th>
                    <th class="px-6 py-4 text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
                  <tr v-for="ticket in tickets" :key="ticket.id" 
                      @click="viewDetails(ticket)" 
                      class="hover:bg-gray-50 dark:hover:bg-slate-700/30 cursor-pointer transition-colors group bg-white dark:bg-transparent">
                    
                    <!-- Buyer Name -->
                    <td class="px-6 py-4">
                      <div class="flex items-center gap-3">
                        <div class="w-8 h-8 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400 font-bold text-xs">
                          {{ getInitials(ticket.form_data) }}
                        </div>
                        <div class="flex flex-col">
                             <span class="font-medium text-gray-900 dark:text-white truncate max-w-[150px]" :title="getBuyerName(ticket.form_data)">
                                {{ getBuyerName(ticket.form_data) }}
                             </span>
                             <span class="text-xs text-gray-500 dark:text-gray-400">{{ getBuyerContact(ticket.form_data) }}</span>
                        </div>
                      </div>
                    </td>

                    <!-- Ticket Type -->
                    <td class="px-6 py-4">
                      <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium" 
                            :style="{ backgroundColor: ticket.tier_color + '20', color: ticket.tier_color }">
                        {{ ticket.tier_icon }} {{ ticket.tier_name }}
                      </span>
                    </td>

                    <!-- Status -->
                    <td class="px-6 py-4">
                      <span :class="getStatusClass(ticket.status)" class="px-2.5 py-1 rounded-full text-xs font-medium">
                        {{ getStatusLabel(ticket.status) }}
                      </span>
                    </td>

                    <!-- Price -->
                    <td class="px-6 py-4 text-right font-medium text-gray-900 dark:text-white">
                      {{ formatAmount(ticket.price) }} <span class="text-xs text-gray-500">{{ ticket.currency }}</span>
                    </td>

                    <!-- Actions -->
                    <td class="px-6 py-4 text-right" @click.stop>
                        <!-- ALWAYS VISIBLE REFUND BUTTON -->
                        <button v-if="ticket.status === 'paid' && isOwner" 
                                @click="openRefundModal(ticket)"
                                class="text-xs font-medium text-red-600 hover:text-red-700 bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:hover:bg-red-900/40 px-3 py-1.5 rounded-lg border border-red-200 dark:border-red-900/30 transition-colors">
                            Rembourser
                        </button>
                    </td>
                  </tr>
                </tbody>
              </table>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="p-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between bg-gray-50/50 dark:bg-slate-800/50">
              <span class="text-sm text-gray-600 dark:text-gray-400">
                  Page {{ currentPage }} sur {{ totalPages }}
              </span>
              <div class="flex gap-2">
                  <button @click="changePage(currentPage - 1)" :disabled="currentPage === 1" 
                          class="px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-white dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300">
                      Précédent
                  </button>
                  <button @click="changePage(currentPage + 1)" :disabled="currentPage === totalPages" 
                          class="px-3 py-1.5 text-sm rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-white dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300">
                      Suivant
                  </button>
              </div>
          </div>
        </div>
      </div>

      <!-- Modals -->
      
      <!-- Details Modal -->
      <Teleport to="body">
        <div v-if="showDetailsModal && selectedTicket" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-md" @click.self="showDetailsModal = false">
          <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200 border border-gray-100 dark:border-gray-800">
            <div class="p-6 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center bg-gray-50/50 dark:bg-slate-800/50">
              <h3 class="text-lg font-bold text-gray-900 dark:text-white">Détails du Ticket</h3>
              <button @click="showDetailsModal = false" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg></button>
            </div>
            <div class="p-6 max-h-[70vh] overflow-y-auto custom-scrollbar space-y-6">
                 <!-- Buyer Full Info -->
                 <div class="text-center">
                     <div class="w-16 h-16 mx-auto rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400 font-bold text-2xl mb-3">
                         {{ getInitials(selectedTicket.form_data) }}
                     </div>
                     <h4 class="text-xl font-bold text-gray-900 dark:text-white">{{ getBuyerName(selectedTicket.form_data) }}</h4>
                     <p class="text-sm text-gray-500 dark:text-gray-400">{{ getBuyerEmail(selectedTicket.form_data) }}</p>
                 </div>

                 <div class="grid grid-cols-2 gap-4">
                     <div class="p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 border border-gray-100 dark:border-gray-700">
                         <p class="text-xs text-gray-500 uppercase">Status</p>
                         <p class="font-medium mt-1" :class="getStatusColor(selectedTicket.status)">{{ getStatusLabel(selectedTicket.status) }}</p>
                     </div>
                     <div class="p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 border border-gray-100 dark:border-gray-700">
                         <p class="text-xs text-gray-500 uppercase">Prix Payé</p>
                         <p class="font-medium mt-1 text-gray-900 dark:text-white">{{ formatAmount(selectedTicket.price) }} {{ selectedTicket.currency }}</p>
                     </div>
                 </div>
 
                 <!-- Full Form Data -->
                 <div v-if="selectedTicket.form_data" class="space-y-3">
                     <h5 class="text-sm font-semibold text-gray-900 dark:text-white">Données du participant</h5>
                     <div class="bg-gray-50 dark:bg-slate-800/50 rounded-xl p-4 space-y-2">
                         <div v-for="(val, key) in selectedTicket.form_data" :key="key" class="flex justify-between text-sm">
                             <span class="text-gray-500">{{ getLabelForField(key) }}</span>
                             <span class="font-medium text-gray-900 dark:text-white">{{ val }}</span>
                         </div>
                     </div>
                 </div>

                 <div class="text-center">
                     <p class="text-xs text-mono text-gray-400">ID: {{ selectedTicket.ticket_code }}</p>
                 </div>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Refund Confirmation Modal -->
      <Teleport to="body">
          <div v-if="showRefundModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showRefundModal = false">
              <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md shadow-2xl p-6 border border-gray-100 dark:border-gray-800 animate-in fade-in zoom-in duration-200">
                  <div class="w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/30 text-red-600 flex items-center justify-center mb-4 mx-auto">
                      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                  </div>
                  <h3 class="text-xl font-bold text-gray-900 dark:text-white text-center mb-2">Confirmer le remboursement</h3>
                  <p class="text-gray-500 dark:text-gray-400 text-center text-sm mb-6">
                      Voulez-vous vraiment rembourser le ticket <strong>{{ ticketToRefund?.ticket_code }}</strong> de <strong>{{ getBuyerName(ticketToRefund?.form_data) }}</strong> ? <br>
                      montant: <span class="font-bold text-gray-900 dark:text-white">{{ formatAmount(ticketToRefund?.price) }} {{ ticketToRefund?.currency }}</span><br><br>
                      Cette action est irréversible.
                  </p>
                  <div class="flex gap-3">
                      <button @click="showRefundModal = false" class="flex-1 px-4 py-2 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors font-medium">Annuler</button>
                      <button @click="openPinForRefund" class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium shadow-lg shadow-red-600/20">Continuer</button>
                  </div>
              </div>
          </div>
      </Teleport>

      <!-- PIN Verification Modal -->
      <Teleport to="body">
          <div v-if="showPinModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showPinModal = false">
              <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-xs shadow-2xl p-6 border border-gray-100 dark:border-gray-800 animate-in fade-in zoom-in duration-200">
                  <h3 class="text-lg font-bold text-gray-900 dark:text-white text-center mb-4">Sécurité</h3>
                  <p class="text-sm text-gray-500 dark:text-gray-400 text-center mb-4">Entrez votre code PIN pour confirmer le remboursement.</p>
                  
                  <div class="mb-4">
                      <input v-model="pinCode" type="password" maxlength="5" placeholder="• • • • •" 
                             class="w-full text-center text-2xl tracking-widest py-2 rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                             :disabled="verifyingPin"
                             @keyup.enter="confirmRefund"
                      >
                  </div>

                  <button @click="confirmRefund" :disabled="verifyingPin || pinCode.length < 5" 
                          class="w-full py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition-colors flex justify-center items-center disabled:opacity-50 disabled:cursor-not-allowed">
                      <svg v-if="verifyingPin" class="animate-spin h-5 w-5 mr-2 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                      {{ verifyingPin ? 'Vérification...' : 'Confirmer' }}
                  </button>
                  <button @click="showPinModal = false" class="w-full mt-2 py-2 text-sm text-gray-500 hover:text-gray-700 dark:hover:text-gray-300">Annuler</button>
              </div>
          </div>
      </Teleport>

      <!-- Cancel Event Confirmation Modal -->
      <Teleport to="body">
          <div v-if="showCancelModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md">
              <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md shadow-2xl p-6 border border-red-500/20 animate-in fade-in zoom-in duration-200">
                  <div class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-900/30 text-red-600 flex items-center justify-center mb-4 mx-auto animate-pulse">
                      <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
                  </div>
                  <h3 class="text-xl font-bold text-red-600 dark:text-red-500 text-center mb-2">ANNULATION D'ÉVÉNEMENT</h3>
                  <p class="text-gray-600 dark:text-gray-300 text-center text-sm mb-6 leading-relaxed">
                      Attention ! Vous êtes sur le point d'annuler l'événement <strong>{{ event?.title }}</strong>.<br><br>
                      🔴 <strong>Tous les tickets vendus ({{ tickets.length }}) seront REMBOURSÉS automatiquement.</strong><br>
                      🔴 Cette action est <strong>IRRÉVERSIBLE</strong>.
                  </p>
                  
                  <div class="bg-red-50 dark:bg-red-900/10 p-4 rounded-lg border border-red-100 dark:border-red-900/20 mb-6">
                      <label class="flex items-start gap-3 cursor-pointer">
                          <input type="checkbox" v-model="confirmCancelCheck" class="mt-1 w-4 h-4 text-red-600 rounded border-gray-300 focus:ring-red-500">
                          <span class="text-sm text-red-800 dark:text-red-300">Je comprends que cette action va déclencher le remboursement de tous les participants et annuler l'événement définitivement.</span>
                      </label>
                  </div>

                  <div class="flex gap-3">
                      <button @click="showCancelModal = false; confirmCancelCheck = false" class="flex-1 px-4 py-2 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors font-medium">Ne pas annuler</button>
                      <button @click="processCancelEvent" :disabled="!confirmCancelCheck" class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium shadow-lg shadow-red-600/20 disabled:opacity-50 disabled:cursor-not-allowed">
                          CONFIRMER L'ANNULATION
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

import api, { ticketAPI, userAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const eventId = computed(() => route.params.id as string)

// State
const loading = ref(true)
const error = ref('')
const event = ref<any>(null)
const tickets = ref<any[]>([])
const stats = ref<any>(null)
const searchQuery = ref('')

// Pagination
const currentPage = ref(1)
const itemsPerPage = 50
const totalPages = computed(() => Math.ceil((stats.value?.total_sold || 0) / itemsPerPage))

// Computed Stats
const tierStats = computed(() => {
    if (!event.value?.ticket_tiers) return []
    return event.value.ticket_tiers.map((t: any) => ({
        id: t.id,
        name: t.name,
        icon: t.icon,
        sold: t.sold || 0,
        capacity: t.quantity,
        color: t.color
    }))
})

const isOwner = computed(() => event.value?.creator_id === authStore.user?.id)

// Data Loading
const loadData = async () => {
  loading.value = true
  error.value = ''
  try {
    const offset = (currentPage.value - 1) * itemsPerPage
    
    // Fetch Event (for tiers/stats) & Tickets (paginated)
    const [eventRes, ticketsRes] = await Promise.all([
      ticketAPI.getEvent(eventId.value),
      ticketAPI.getEventTickets(eventId.value, itemsPerPage, offset)
    ])
    
    event.value = eventRes.data?.event || eventRes.data
    tickets.value = ticketsRes.data?.tickets || ticketsRes.data || []
    
    // Set Stats generic object for total counts
    stats.value = {
        total_sold: event.value.total_sold || tickets.value.length, 
        total_revenue: event.value.total_revenue || 0
    }

  } catch (err: any) {
    console.error('Failed to load data:', err)
    error.value = err.response?.data?.error || 'Erreur lors du chargement des tickets'
  } finally {
    loading.value = false
  }
}

// Search with Debounce
let searchTimeout: any
const handleSearch = () => {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        currentPage.value = 1
        // Trigger loadData passing search query
        // We'll modify loadData to read searchQuery
        loadDataWithSearch()
    }, 500)
}

const loadDataWithSearch = async () => {
  loading.value = true
  try {
    const offset = (currentPage.value - 1) * itemsPerPage
    const searchRes = await api.get(`/ticket-service/api/v1/events/${eventId.value}/tickets?limit=${itemsPerPage}&offset=${offset}&search=${encodeURIComponent(searchQuery.value)}`)
    tickets.value = searchRes.data?.tickets || []
    
    // Refresh stats if needed? Or keep event stats global?
    // Let's keep existing event stats but update tickets list.
  } catch (err) {
      console.error(err)
  } finally {
      loading.value = false
  }
}

watch(searchQuery, handleSearch)

// Helpers
const getInitials = (formData: any) => {
  if (!formData) return '?'
  const name = formData.name || formData.nom || formData.full_name || ''
  return name.split(' ').map((n: string) => n[0]).join('').toUpperCase().slice(0, 2) || '?'
}

const getBuyerName = (formData: any) => {
  if (!formData) return 'Anonyme'
  return formData.name || formData.nom || formData.full_name || 'Anonyme'
}
const getBuyerEmail = (formData: any) => formData?.email || formData?.Email || ''

const getBuyerContact = (formData: any) => {
    if (!formData) return ''
    // Prioritize Phone, then Email
    const phone = formData.phone || formData.telephone || formData.mobile || formData.tel
    if (phone) return phone
    return formData.email || formData.Email || ''
}

const formatAmount = (amount: number) => new Intl.NumberFormat('fr-FR').format(amount || 0)

const getStatusClass = (status: string) => {
  switch (status) {
    case 'paid': return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
    case 'used': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
    case 'cancelled': 
    case 'refunded': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400 line-through'
    default: return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  }
}
const getStatusColor = (status: string) => {
    switch(status) {
        case 'paid': return 'text-emerald-600'
        case 'used': return 'text-blue-600'
        case 'refunded': return 'text-red-600'
        default: return 'text-amber-600'
    }
}
const getStatusLabel = (status: string) => {
  switch (status) {
    case 'paid': return 'Confirmé'
    case 'used': return 'Utilisé'
    case 'cancelled': return 'Annulé'
    case 'refunded': return 'Remboursé'
    case 'pending': return 'En attente'
    default: return status
  }
}

const getLabelForField = (fieldName: string) => {
  if (!event.value?.form_fields) return fieldName
  const field = event.value.form_fields.find((f: any) => f.name === fieldName)
  return field ? field.label : fieldName
}

// Modals Logic
const showDetailsModal = ref(false)
const selectedTicket = ref<any>(null)
const viewDetails = (ticket: any) => {
  selectedTicket.value = ticket
  showDetailsModal.value = true
}

// PIN & Refund Logic
const showRefundModal = ref(false)
const showPinModal = ref(false)
const pinCode = ref('')
const verifyingPin = ref(false)
const ticketToRefund = ref<any>(null)

const openRefundModal = (ticket: any) => {
    ticketToRefund.value = ticket
    showRefundModal.value = true
}

const openPinForRefund = () => {
    showRefundModal.value = false // Close warning
    showPinModal.value = true     // Open PIN
    pinCode.value = ''
}

const confirmRefund = async () => {
    if (!pinCode.value || pinCode.value.length < 5) {
        alert("Veuillez entrer un code PIN valide (5 chiffres).")
        return
    }
    
    verifyingPin.value = true
    try {
        // Verify PIN
        await userAPI.verifyPin({ pin: pinCode.value })
        
        // If success, proceed to refund
        await executeRefund()
        showPinModal.value = false
    } catch (err: any) {
        alert(err.response?.data?.error || "Code PIN incorrect.")
    } finally {
        verifyingPin.value = false
    }
}

const executeRefund = async () => {
    if (!ticketToRefund.value) return
    try {
        await ticketAPI.refundTicket(ticketToRefund.value.id)
        // Optimistic update
        const t = tickets.value.find(t => t.id === ticketToRefund.value.id)
        if(t) t.status = 'refunded'
        alert("Remboursement effectué avec succès.")
    } catch(err: any) {
        alert(err.response?.data?.error || 'Erreur lors du remboursement')
    }
}

const showCancelModal = ref(false)
const confirmCancelCheck = ref(false)
const openCancelModal = () => {
    showCancelModal.value = true
    confirmCancelCheck.value = false
}
const processCancelEvent = async () => {
    try {
        await ticketAPI.cancelEvent(eventId.value)
        showCancelModal.value = false
        alert('Événement annulé et remboursements initiés.')
        loadData()
    } catch(err: any) {
        alert(err.response?.data?.error || "Erreur lors de l'annulation")
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
@keyframes spin { to { transform: rotate(360deg); } }
</style>
