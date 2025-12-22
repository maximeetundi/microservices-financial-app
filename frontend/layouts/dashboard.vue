<template>
  <div class="min-h-screen flex bg-base transition-colors duration-300">
    <!-- Sidebar -->
    <aside class="w-72 fixed h-full bg-surface border-r border-secondary-200 dark:border-secondary-800 flex flex-col transition-all duration-300 z-50">
      <!-- Logo -->
      <div class="p-6 border-b border-secondary-100 dark:border-secondary-800">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shadow-lg shadow-primary-500/20">
            <span class="text-2xl">🏦</span>
          </div>
          <div>
            <h1 class="text-lg font-bold bg-gradient-to-r from-secondary-900 to-secondary-600 dark:from-white dark:to-secondary-400 bg-clip-text text-transparent">CryptoBank</h1>
            <p class="text-xs text-muted font-medium">Premium Banking</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 py-6 px-4 space-y-1 overflow-y-auto">
        <NuxtLink to="/dashboard" class="nav-item" active-class="active">
          <span class="icon">📊</span>
          <span>Tableau de bord</span>
        </NuxtLink>

        <NuxtLink to="/wallet" class="nav-item" active-class="active">
          <span class="icon">👛</span>
          <span>Portefeuilles</span>
        </NuxtLink>

        <NuxtLink to="/cards" class="nav-item" active-class="active">
          <span class="icon">💳</span>
          <span>Mes Cartes</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Échange</p>
        </div>

        <NuxtLink to="/exchange/crypto" class="nav-item" active-class="active">
          <span class="icon">₿</span>
          <span>Crypto</span>
        </NuxtLink>

        <NuxtLink to="/exchange/fiat" class="nav-item" active-class="active">
          <span class="icon">💱</span>
          <span>Fiat</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Opérations</p>
        </div>

        <NuxtLink to="/transfer" class="nav-item" active-class="active">
          <span class="icon">💸</span>
          <span>Virements</span>
        </NuxtLink>
        
        <!-- Only show merchant link if user is merchant - TODO: logic -->
        <NuxtLink to="/merchant" class="nav-item" active-class="active">
          <span class="icon">🏪</span>
          <span>Espace Marchand</span>
        </NuxtLink>
      </nav>

      <!-- User Section -->
      <div class="p-4 border-t border-secondary-100 dark:border-secondary-800 bg-surface-hover/30">
        <div class="flex items-center gap-3 p-3 rounded-xl hover:bg-white/50 dark:hover:bg-black/20 transition-colors">
          <div class="w-10 h-10 rounded-full bg-gradient-to-tr from-primary-400 to-secondary-400 flex items-center justify-center text-white font-bold shadow-md">
            {{ userInitials }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-base truncate">{{ userName }}</p>
            <p class="text-xs text-muted truncate">{{ userEmail }}</p>
          </div>
          
          <NotificationCenter />
          <ThemeToggle />
          
          <button @click="handleLogout" class="p-2 text-muted hover:text-error dark:hover:text-red-400 transition-colors" title="Déconnexion">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 ml-72 p-8 transition-all duration-300">
      <div class="max-w-7xl mx-auto animate-fade-in-up">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { computed } from 'vue'

const authStore = useAuthStore()

const userName = computed(() => {
  if (authStore.user) {
    return `${authStore.user.first_name || ''} ${authStore.user.last_name || ''}`
  }
  return 'Utilisateur'
})

const userEmail = computed(() => authStore.user?.email || '')

const userInitials = computed(() => {
  if (authStore.user) {
    return `${authStore.user.first_name?.[0] || ''}${authStore.user.last_name?.[0] || ''}`
  }
  return 'U'
})

const handleLogout = () => {
  authStore.logout()
}
</script>

<style scoped>
.nav-item {
  @apply flex items-center gap-3 px-4 py-3 text-sm font-medium text-muted rounded-xl transition-all duration-200 hover:bg-surface-hover hover:text-base;
}

.nav-item.active {
  @apply bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 shadow-sm;
}

.nav-item .icon {
  @apply text-lg opacity-70 transition-opacity;
}

.nav-item:hover .icon,
.nav-item.active .icon {
  @apply opacity-100;
}

.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
