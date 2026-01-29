<template>
  <NuxtLayout name="admin">
    <div class="p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">Analytics</h1>
        <p class="text-slate-400">Vue d'ensemble de l'activité de la plateforme.</p>
      </div>

      <!-- KPI Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 rounded-xl bg-indigo-500/20 flex items-center justify-center text-indigo-400 text-2xl">
              👥
            </div>
            <div>
              <div class="text-slate-400 text-sm">Utilisateurs</div>
              <div class="text-2xl font-bold text-white">{{ stats.totalUsers }}</div>
            </div>
          </div>
          <div class="text-xs text-emerald-400 flex items-center gap-1">
            <span>↑ 12%</span>
            <span class="text-slate-500">vs mois dernier</span>
          </div>
        </div>

        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500/20 flex items-center justify-center text-emerald-400 text-2xl">
              📋
            </div>
            <div>
              <div class="text-slate-400 text-sm">KYC en attente</div>
              <div class="text-2xl font-bold text-white">{{ stats.pendingKYC }}</div>
            </div>
          </div>
          <div class="text-xs text-slate-500">
            Documents à vérifier
          </div>
        </div>

        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 rounded-xl bg-purple-500/20 flex items-center justify-center text-purple-400 text-2xl">
              💳
            </div>
            <div>
              <div class="text-slate-400 text-sm">Volume Trans. (24h)</div>
              <div class="text-2xl font-bold text-white">45.2K €</div>
            </div>
          </div>
          <div class="text-xs text-emerald-400 flex items-center gap-1">
            <span>↑ 5%</span>
            <span class="text-slate-500">vs hier</span>
          </div>
        </div>

        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 rounded-xl bg-amber-500/20 flex items-center justify-center text-amber-400 text-2xl">
              🛡️
            </div>
            <div>
              <div class="text-slate-400 text-sm">Comptes bloqués</div>
              <div class="text-2xl font-bold text-white">3</div>
            </div>
          </div>
          <div class="text-xs text-slate-500">
            Nécessitent une action
          </div>
        </div>
      </div>

      <!-- Charts Section (Placeholder) -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50 h-80 flex flex-col items-center justify-center text-slate-500">
          <span class="text-4xl mb-4">📊</span>
          <p>Graphique: Nouveaux utilisateurs (Coming Soon)</p>
        </div>
        <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/50 h-80 flex flex-col items-center justify-center text-slate-500">
          <span class="text-4xl mb-4">🥧</span>
          <p>Graphique: Répartition des volumes (Coming Soon)</p>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { adminUserAPI, adminKYCAPI } = useApi()

const stats = ref({
  totalUsers: 0,
  pendingKYC: 0
})

onMounted(async () => {
  try {
    // Fetch stats in parallel
    const [usersRes, kycRes] = await Promise.all([
      adminUserAPI.getUsers(1, 0),
      adminKYCAPI.getRequests('pending', 1, 0)
    ])
    
    stats.value.totalUsers = usersRes.data?.total || 0
    stats.value.pendingKYC = kycRes.data?.total || 0
    
  } catch (error) {
    console.error("Failed to fetch analytics stats", error)
  }
})
</script>
