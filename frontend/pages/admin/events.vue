<template>
  <NuxtLayout name="admin">
    <div class="p-6 lg:p-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">Événements</h1>
          <p class="text-slate-400">Gérez les événements et la billetterie</p>
        </div>
        <div class="flex gap-2">
           <button @click="refreshData" class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg transition-colors flex items-center gap-2">
            <span>🔄</span> Refresh
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex gap-4 border-b border-slate-700 mb-6">
        <button 
          @click="activeTab = 'events'"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="activeTab === 'events' ? 'text-indigo-400' : 'text-slate-400 hover:text-slate-300'"
        >
          Événements
          <div v-if="activeTab === 'events'" class="absolute bottom-0 left-0 w-full h-0.5 bg-indigo-500"></div>
        </button>
        <button 
          @click="activeTab = 'tickets'"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="activeTab === 'tickets' ? 'text-indigo-400' : 'text-slate-400 hover:text-slate-300'"
        >
          Billets Vendus
          <div v-if="activeTab === 'tickets'" class="absolute bottom-0 left-0 w-full h-0.5 bg-indigo-500"></div>
        </button>
      </div>

      <!-- Events Table -->
      <div v-if="activeTab === 'events'" class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-800/50 text-slate-400 uppercase text-xs">
              <tr>
                <th class="px-6 py-4">Titre</th>
                <th class="px-6 py-4">Organisateur</th>
                <th class="px-6 py-4">Lieu</th>
                <th class="px-6 py-4">Date Début</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4">Créé le</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50 text-slate-300">
              <tr v-if="loading" class="animate-pulse">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Chargement...</td>
              </tr>
              <tr v-else-if="events.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Aucun événement trouvé</td>
              </tr>
              <tr v-else v-for="e in events" :key="e.id" class="hover:bg-slate-800/30 transition-colors">
                <td class="px-6 py-4 font-medium text-white">{{ e.title }}</td>
                <td class="px-6 py-4 text-sm">{{ e.creator_id ? e.creator_id.substring(0, 8) + '...' : 'N/A' }}</td>
                <td class="px-6 py-4">{{ e.location || 'En ligne' }}</td>
                <td class="px-6 py-4">{{ formatDate(e.start_date) }}</td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-emerald-500/20 text-emerald-400': e.status === 'active',
                      'bg-slate-500/20 text-slate-400': e.status === 'draft',
                      'bg-red-500/20 text-red-400': e.status === 'cancelled',
                      'bg-blue-500/20 text-blue-400': e.status === 'ended'
                    }">
                    {{ e.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-400">{{ formatDate(e.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Tickets Table -->
      <div v-else class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-800/50 text-slate-400 uppercase text-xs">
              <tr>
                <th class="px-6 py-4">Code</th>
                <th class="px-6 py-4">Événement ID</th>
                <th class="px-6 py-4">Acheteur</th>
                <th class="px-6 py-4">Prix</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4">Date Achat</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50 text-slate-300">
              <tr v-if="loading" class="animate-pulse">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Chargement...</td>
              </tr>
              <tr v-else-if="tickets.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Aucun billet vendu trouvé</td>
              </tr>
              <tr v-else v-for="t in tickets" :key="t.id" class="hover:bg-slate-800/30 transition-colors">
                <td class="px-6 py-4 font-mono text-indigo-400 font-bold">{{ t.ticket_code }}</td>
                <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ t.event_id }}</td>
                <td class="px-6 py-4 text-sm">{{ t.buyer_id ? t.buyer_id.substring(0, 8) + '...' : 'N/A' }}</td>
                <td class="px-6 py-4">{{ formatMoney(t.price, t.currency || 'XOF') }}</td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-emerald-500/20 text-emerald-400': t.status === 'valid',
                      'bg-blue-500/20 text-blue-400': t.status === 'used',
                      'bg-red-500/20 text-red-400': t.status === 'refunded' || t.status === 'cancelled'
                    }">
                    {{ t.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-400">{{ formatDate(t.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'

definePageMeta({
  layout: false,
  middleware: 'admin-auth'
})

const activeTab = ref('events')
const loading = ref(false)
const events = ref([])
const tickets = ref([])

const config = useRuntimeConfig()
const apiUrl = config.public.adminApiUrl || 'https://api.admin.maximeetundi.store'

const fetchEvents = async () => {
  try {
    const adminToken = localStorage.getItem('adminToken')
    const response = await fetch(`${apiUrl}/api/v1/admin/events?limit=50`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    })
    if (response.ok) {
        events.value = await response.json() || []
    }
  } catch (e) { console.error(e) }
}

const fetchTickets = async () => {
    try {
    const adminToken = localStorage.getItem('adminToken')
    const response = await fetch(`${apiUrl}/api/v1/admin/sold-tickets?limit=50`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    })
    if (response.ok) {
        tickets.value = await response.json() || []
    }
  } catch (e) { console.error(e) }
}

const refreshData = async () => {
  loading.value = true
  await Promise.all([fetchEvents(), fetchTickets()])
  loading.value = false
}

onMounted(() => {
  refreshData()
})

const formatMoney = (amount, currency = 'XOF') => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency }).format(amount)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.glass-card {
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(71, 85, 105, 0.5);
  border-radius: 1rem;
}
</style>
