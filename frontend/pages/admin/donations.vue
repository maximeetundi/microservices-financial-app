<template>
  <NuxtLayout name="admin">
    <div class="p-6 lg:p-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">Dons & Solidarité</h1>
          <p class="text-slate-400">Gérez les campagnes et suivez les dons</p>
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
          @click="activeTab = 'campaigns'"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="activeTab === 'campaigns' ? 'text-indigo-400' : 'text-slate-400 hover:text-slate-300'"
        >
          Campagnes
          <div v-if="activeTab === 'campaigns'" class="absolute bottom-0 left-0 w-full h-0.5 bg-indigo-500"></div>
        </button>
        <button 
          @click="activeTab = 'donations'"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="activeTab === 'donations' ? 'text-indigo-400' : 'text-slate-400 hover:text-slate-300'"
        >
          Dons Récents
          <div v-if="activeTab === 'donations'" class="absolute bottom-0 left-0 w-full h-0.5 bg-indigo-500"></div>
        </button>
      </div>

      <!-- Campaigns Table -->
      <div v-if="activeTab === 'campaigns'" class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-800/50 text-slate-400 uppercase text-xs">
              <tr>
                <th class="px-6 py-4">Titre</th>
                <th class="px-6 py-4">Créateur</th>
                <th class="px-6 py-4">Objectif</th>
                <th class="px-6 py-4">Collecté</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50 text-slate-300">
              <tr v-if="loading" class="animate-pulse">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Chargement...</td>
              </tr>
              <tr v-else-if="campaigns.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Aucune campagne trouvée</td>
              </tr>
              <tr v-else v-for="c in campaigns" :key="c.id" class="hover:bg-slate-800/30 transition-colors">
                <td class="px-6 py-4 font-medium text-white">{{ c.title }}</td>
                <td class="px-6 py-4 text-sm">{{ c.creator_id ? c.creator_id.substring(0, 8) + '...' : 'N/A' }}</td>
                <td class="px-6 py-4">{{ formatMoney(c.target_amount, c.currency) }}</td>
                <td class="px-6 py-4">
                  <div class="flex flex-col gap-1">
                    <span>{{ formatMoney(c.collected_amount, c.currency) }}</span>
                    <div class="w-full h-1.5 bg-slate-700 rounded-full overflow-hidden">
                      <div class="h-full bg-emerald-500" :style="{ width: Math.min((c.collected_amount / c.target_amount) * 100, 100) + '%' }"></div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-emerald-500/20 text-emerald-400': c.status === 'active',
                      'bg-amber-500/20 text-amber-400': c.status === 'paused',
                      'bg-slate-500/20 text-slate-400': c.status === 'draft' || c.status === 'completed'
                    }">
                    {{ c.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-400">{{ formatDate(c.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Donations Table -->
      <div v-else class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-800/50 text-slate-400 uppercase text-xs">
              <tr>
                <th class="px-6 py-4">ID Donateur</th>
                <th class="px-6 py-4">Campagne</th>
                <th class="px-6 py-4">Montant</th>
                <th class="px-6 py-4">Message</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50 text-slate-300">
              <tr v-if="loading" class="animate-pulse">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Chargement...</td>
              </tr>
              <tr v-else-if="donations.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-500">Aucun don trouvé</td>
              </tr>
              <tr v-else v-for="d in donations" :key="d.id" class="hover:bg-slate-800/30 transition-colors">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-mono bg-slate-800 px-1.5 py-0.5 rounded">{{ d.donor_id ? d.donor_id.substring(0, 6) : 'Guest' }}</span>
                    <span v-if="d.is_anonymous" class="text-xs text-indigo-400" title="Anonyme">👻</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-sm text-slate-400 font-mono">{{ d.campaign_id }}</td>
                <td class="px-6 py-4 font-bold text-emerald-400">+{{ formatMoney(d.amount, d.currency) }}</td>
                <td class="px-6 py-4 text-sm italic text-slate-500 truncate max-w-xs">{{ d.message || '-' }}</td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-emerald-500/20 text-emerald-400': d.status === 'paid',
                      'bg-amber-500/20 text-amber-400': d.status === 'pending',
                      'bg-red-500/20 text-red-400': d.status === 'failed' || d.status === 'refunded'
                    }">
                    {{ d.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-400">{{ formatDate(d.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'

definePageMeta({
  layout: false,
  middleware: 'admin-auth'
})

const activeTab = ref('campaigns')
const loading = ref(false)
const campaigns = ref([])
const donations = ref([])

const config = useRuntimeConfig()
const apiUrl = config.public.adminApiUrl || 'https://api.admin.maximeetundi.store'

const fetchCampaigns = async () => {
  try {
    const adminToken = localStorage.getItem('adminToken')
    const response = await fetch(`${apiUrl}/api/v1/admin/campaigns?limit=50`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    })
    if (response.ok) {
        campaigns.value = await response.json() || []
    }
  } catch (e) { console.error(e) }
}

const fetchDonations = async () => {
    try {
    const adminToken = localStorage.getItem('adminToken')
    const response = await fetch(`${apiUrl}/api/v1/admin/donations?limit=50`, {
      headers: { 'Authorization': `Bearer ${adminToken}` }
    })
    if (response.ok) {
        donations.value = await response.json() || []
    }
  } catch (e) { console.error(e) }
}

const refreshData = async () => {
  loading.value = true
  await Promise.all([fetchCampaigns(), fetchDonations()])
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
