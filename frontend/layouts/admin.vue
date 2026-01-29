<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
    <!-- Sidebar -->
    <div class="fixed left-0 top-0 h-full w-64 bg-slate-900/80 backdrop-blur-xl border-r border-slate-700/50 z-40">
      <!-- Logo -->
      <div class="p-6 border-b border-slate-700/50">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600 flex items-center justify-center">
            <span class="text-xl">🛡️</span>
          </div>
          <div>
            <h1 class="font-bold text-white">Admin Panel</h1>
            <p class="text-xs text-slate-400">Zekora</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="p-4 space-y-1">
        <NuxtLink to="/admin/dashboard" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/dashboard') }">
          <span class="text-lg">📊</span>
          <span>Dashboard</span>
        </NuxtLink>
        <NuxtLink to="/admin/payments" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/payments') }">
          <span class="text-lg">💳</span>
          <span>Agrégateurs Paiement</span>
        </NuxtLink>
        <NuxtLink to="/admin/fees"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/fees') }">
          <span class="text-lg">💰</span>
          <span>Frais & Commissions</span>
        </NuxtLink>
        <NuxtLink to="/admin/users" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/users') }">
          <span class="text-lg">👥</span>
          <span>Utilisateurs</span>
        </NuxtLink>
        <NuxtLink to="/admin/donations" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/donations') }">
          <span class="text-lg">🤲</span>
          <span>Dons & Solidarité</span>
        </NuxtLink>
        <NuxtLink to="/admin/events" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/events') }">
          <span class="text-lg">🎫</span>
          <span>Événements</span>
        </NuxtLink>
        <NuxtLink to="/admin/kyc" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/kyc') }">
          <span class="text-lg">📋</span>
          <span>KYC</span>
        </NuxtLink>
        <NuxtLink to="/admin/transactions" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/transactions') }">
          <span class="text-lg">💸</span>
          <span>Transactions</span>
        </NuxtLink>
        <NuxtLink to="/admin/analytics" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/analytics') }">
          <span class="text-lg">📈</span>
          <span>Analytics</span>
        </NuxtLink>
        <NuxtLink to="/admin/platform" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/platform') }">
          <span class="text-lg">🏦</span>
          <span>Comptes Plateforme</span>
        </NuxtLink>
        <NuxtLink to="/admin/settings" 
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-300 hover:text-white hover:bg-slate-800/50 transition-all"
          :class="{ 'bg-indigo-500/20 text-indigo-400': isActive('/admin/settings') }">
          <span class="text-lg">⚙️</span>
          <span>Global Settings</span>
        </NuxtLink>
      </nav>

      <!-- User Info -->
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-slate-700/50">
        <div class="flex items-center gap-3 px-4 py-3">
          <div class="w-10 h-10 rounded-full bg-indigo-500/20 flex items-center justify-center">
            <span class="text-lg">👤</span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-white truncate">{{ adminUser?.first_name || 'Admin' }}</p>
            <p class="text-xs text-slate-400 truncate">{{ adminUser?.email || '' }}</p>
          </div>
          <button @click="logout" class="p-2 hover:bg-red-500/20 rounded-lg text-slate-400 hover:text-red-400 transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="ml-64">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const adminUser = ref(null)

const isActive = (path) => route.path === path || route.path.startsWith(path + '/')

const logout = () => {
  localStorage.removeItem('adminToken')
  localStorage.removeItem('adminUser')
  router.push('/admin/login')
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('adminUser')
    if (stored) {
      try {
        adminUser.value = JSON.parse(stored)
      } catch {}
    }
  }
})
</script>
