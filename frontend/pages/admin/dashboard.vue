<template>
  <NuxtLayout name="admin">
    <div class="p-6 lg:p-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">👋 Bienvenue, {{ adminName }}</h1>
        <p class="text-slate-400">Vue d'ensemble de votre plateforme</p>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">Utilisateurs</p>
              <p class="text-2xl font-bold text-white">{{ stats.users.toLocaleString() }}</p>
              <p class="text-xs text-emerald-400">+12% ce mois</p>
            </div>
            <div class="p-3 bg-indigo-500/20 rounded-xl">
              <span class="text-2xl">👥</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">Transactions</p>
              <p class="text-2xl font-bold text-white">{{ stats.transactions.toLocaleString() }}</p>
              <p class="text-xs text-emerald-400">+8% ce mois</p>
            </div>
            <div class="p-3 bg-emerald-500/20 rounded-xl">
              <span class="text-2xl">💸</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">Volume total</p>
              <p class="text-2xl font-bold text-white">{{ formatMoney(stats.volume) }}</p>
              <p class="text-xs text-emerald-400">+25% ce mois</p>
            </div>
            <div class="p-3 bg-purple-500/20 rounded-xl">
              <span class="text-2xl">📊</span>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">KYC en attente</p>
              <p class="text-2xl font-bold text-white">{{ stats.pendingKyc }}</p>
              <p class="text-xs text-amber-400">Nécessite attention</p>
            </div>
            <div class="p-3 bg-amber-500/20 rounded-xl">
              <span class="text-2xl">📋</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <NuxtLink to="/admin/payments" class="glass-card p-6 hover:border-indigo-500/50 transition-all group">
          <div class="flex items-center gap-4">
            <div class="p-4 bg-indigo-500/20 rounded-xl group-hover:scale-110 transition-transform">
              <span class="text-3xl">💳</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-white">Agrégateurs de Paiement</h3>
              <p class="text-sm text-slate-400">Configurer Flutterwave, CinetPay...</p>
            </div>
          </div>
        </NuxtLink>
        <div class="glass-card p-6 hover:border-emerald-500/50 transition-all group cursor-pointer">
          <div class="flex items-center gap-4">
            <div class="p-4 bg-emerald-500/20 rounded-xl group-hover:scale-110 transition-transform">
              <span class="text-3xl">👥</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-white">Gestion Utilisateurs</h3>
              <p class="text-sm text-slate-400">Voir et gérer les comptes</p>
            </div>
          </div>
        </div>
        <div class="glass-card p-6 hover:border-purple-500/50 transition-all group cursor-pointer">
          <div class="flex items-center gap-4">
            <div class="p-4 bg-purple-500/20 rounded-xl group-hover:scale-110 transition-transform">
              <span class="text-3xl">📊</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-white">Analytics</h3>
              <p class="text-sm text-slate-400">Rapports et statistiques</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="glass-card">
        <div class="p-4 border-b border-slate-700/50">
          <h2 class="text-lg font-semibold text-white">Activité récente</h2>
        </div>
        <div class="divide-y divide-slate-700/50">
          <div v-for="activity in recentActivities" :key="activity.id" class="p-4 flex items-center gap-4">
            <div class="p-2 rounded-lg" :class="activity.bgColor">
              <span class="text-lg">{{ activity.icon }}</span>
            </div>
            <div class="flex-1">
              <p class="text-white">{{ activity.title }}</p>
              <p class="text-sm text-slate-400">{{ activity.description }}</p>
            </div>
            <span class="text-xs text-slate-500">{{ activity.time }}</span>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const stats = ref({
  users: 1247,
  transactions: 8453,
  volume: 125000000,
  pendingKyc: 23
})

const recentActivities = ref([
  { id: 1, icon: '👤', title: 'Nouvel utilisateur inscrit', description: 'john.doe@email.com', time: 'Il y a 5 min', bgColor: 'bg-indigo-500/20' },
  { id: 2, icon: '✅', title: 'KYC approuvé', description: 'Marie Dupont', time: 'Il y a 12 min', bgColor: 'bg-emerald-500/20' },
  { id: 3, icon: '💰', title: 'Dépôt effectué', description: '500,000 XOF via Orange Money', time: 'Il y a 25 min', bgColor: 'bg-purple-500/20' },
  { id: 4, icon: '🔒', title: 'Compte bloqué', description: 'suspicious.user@email.com', time: 'Il y a 1h', bgColor: 'bg-red-500/20' },
])

const adminName = computed(() => {
  if (typeof window !== 'undefined') {
    const admin = localStorage.getItem('adminUser')
    if (admin) {
      try {
        const parsed = JSON.parse(admin)
        return parsed.first_name || 'Admin'
      } catch { return 'Admin' }
    }
  }
  return 'Admin'
})

const formatMoney = (amount) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'XOF', maximumFractionDigits: 0 }).format(amount)
}

definePageMeta({
  layout: false,
  middleware: 'admin-auth'
})
</script>

<style scoped>
.glass-card {
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(71, 85, 105, 0.5);
  border-radius: 1rem;
}
</style>
