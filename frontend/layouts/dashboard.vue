<template>
  <div class="min-h-screen flex bg-base transition-colors duration-300">
    <!-- Mobile Overlay -->
    <div 
      v-if="sidebarOpen" 
      class="fixed inset-0 bg-black/50 backdrop-blur-sm z-40 lg:hidden"
      @click="sidebarOpen = false"
    ></div>

    <!-- Sidebar -->
    <aside 
      :class="[
        'fixed h-full bg-slate-100 dark:bg-surface border-r border-secondary-200 dark:border-secondary-800 flex flex-col transition-all duration-300 z-50',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
        'w-72'
      ]"
    >
      <!-- Logo & Close Button -->
      <div class="p-6 border-b border-secondary-100 dark:border-secondary-800 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img src="/logo.png" alt="Zekora Logo" class="w-10 h-10 rounded-xl shadow-lg shadow-primary-500/20 object-cover" />
          <div>
            <h1 class="text-lg font-bold bg-gradient-to-r from-secondary-900 to-secondary-600 dark:from-white dark:to-secondary-400 bg-clip-text text-transparent">Zekora</h1>
            <p class="text-xs text-muted font-medium">Premium Banking</p>
          </div>
        </div>
        <!-- Close button (mobile only) -->
        <button 
          @click="sidebarOpen = false" 
          class="lg:hidden p-2 rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-800 transition-colors"
        >
          <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 py-6 px-4 space-y-1 overflow-y-auto">
        <NuxtLink to="/dashboard" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">📊</span>
          <span>Tableau de bord</span>
        </NuxtLink>

        <NuxtLink to="/wallet" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">👛</span>
          <span>Portefeuilles</span>
        </NuxtLink>

        <NuxtLink to="/cards" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">💳</span>
          <span>Mes Cartes</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Professionnel</p>
        </div>

        <NuxtLink to="/enterprise" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🏢</span>
          <span>Entreprise</span>
        </NuxtLink>

        <NuxtLink to="/shops" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🛍️</span>
          <span>Boutiques</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Échange</p>
        </div>

        <NuxtLink to="/exchange/crypto" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">₿</span>
          <span>Crypto</span>
        </NuxtLink>

        <NuxtLink to="/exchange/fiat" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">💱</span>
          <span>Fiat</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Communauté</p>
        </div>

        <NuxtLink to="/events" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🎫</span>
          <span>Événements</span>
        </NuxtLink>

        <NuxtLink to="/donations" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🤲</span>
          <span>Dons & Solidarité</span>
        </NuxtLink>

        <div class="pt-4 pb-2">
          <p class="text-xs font-semibold text-muted px-4 uppercase tracking-wider">Opérations</p>
        </div>

        <NuxtLink to="/transactions" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">📝</span>
          <span>Historique</span>
        </NuxtLink>

        <NuxtLink to="/transfer" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">💸</span>
          <span>Virements</span>
        </NuxtLink>
        
        <NuxtLink to="/merchant" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🏪</span>
          <span>Espace Marchand</span>
        </NuxtLink>

        <NuxtLink to="/scan" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">📷</span>
          <span>Scanner / Payer</span>
        </NuxtLink>

        <NuxtLink to="/support" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🎧</span>
          <span>Assistance</span>
        </NuxtLink>

        <NuxtLink to="/notifications" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">🔔</span>
          <span>Notifications</span>
        </NuxtLink>

        <NuxtLink to="/messages" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">💬</span>
          <span>Messages</span>
        </NuxtLink>

        <NuxtLink to="/settings" class="nav-item" active-class="active" @click="closeSidebarOnMobile">
          <span class="icon">⚙️</span>
          <span>Paramètres</span>
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
    <main class="flex-1 lg:ml-72 transition-all duration-300">
      <!-- Mobile Header with Hamburger -->
      <header class="lg:hidden sticky top-0 z-30 bg-slate-100/95 dark:bg-surface/95 backdrop-blur-md border-b border-secondary-200 dark:border-secondary-800 px-4 py-3 flex items-center justify-between">
        <button 
          @click="sidebarOpen = true" 
          class="p-2 rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-800 transition-colors"
        >
          <svg class="w-6 h-6 text-base" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        
        <div class="flex items-center gap-2">
          <img src="/logo.png" alt="Zekora Logo" class="w-8 h-8 rounded-lg object-cover" />
          <span class="font-bold text-base">Zekora</span>
        </div>

        <div class="flex items-center gap-1">
          <NotificationCenter />
          <ThemeToggle />
        </div>
      </header>

      <!-- Page Content -->
      <div class="p-4 lg:p-8 w-full max-w-full overflow-x-hidden">
        <div class="w-full max-w-7xl mx-auto animate-fade-in-up">
          <slot />
        </div>
      </div>
    </main>

    <!-- Floating Action Button (FAB) -->
    <!-- Backdrop to close menu when clicking outside -->
    <div 
      v-if="showFabMenu" 
      class="fixed inset-0 z-40" 
      @click="showFabMenu = false"
    ></div>
    
    <div class="fixed bottom-6 right-6 z-50 flex flex-col gap-3 items-end">
      <!-- Menu Items -->
      <transition 
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="transform translate-y-4 opacity-0 scale-90"
        enter-to-class="transform translate-y-0 opacity-100 scale-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="transform translate-y-0 opacity-100 scale-100"
        leave-to-class="transform translate-y-4 opacity-0 scale-90"
      >
        <div v-if="showFabMenu" class="flex flex-col gap-3 items-end mb-2">
          <NuxtLink 
            to="/scan" 
            @click="showFabMenu = false"
            class="flex items-center gap-3 px-4 py-2 bg-white dark:bg-surface text-base font-semibold rounded-full shadow-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-secondary-200 dark:border-secondary-700"
          >
            <span class="whitespace-nowrap">Scanner QR</span>
            <span class="w-10 h-10 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-xl">📷</span>
          </NuxtLink>

          <NuxtLink 
            to="/transfer" 
            @click="showFabMenu = false"
            class="flex items-center gap-3 px-4 py-2 bg-white dark:bg-surface text-base font-semibold rounded-full shadow-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-secondary-200 dark:border-secondary-700"
          >
            <span class="whitespace-nowrap">Virement</span>
            <span class="w-10 h-10 rounded-full bg-green-100 text-green-600 flex items-center justify-center text-xl">💸</span>
          </NuxtLink>
        </div>
      </transition>

      <!-- Main Trigger Button -->
      <button 
        @click.stop="showFabMenu = !showFabMenu"
        class="w-14 h-14 rounded-full bg-gradient-to-r from-primary-600 to-primary-500 text-white shadow-xl shadow-primary-500/30 flex items-center justify-center hover:scale-105 active:scale-95 transition-all duration-300"
      >
        <transition mode="out-in">
          <svg v-if="!showFabMenu" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          <svg v-else class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </transition>
      </button>
    </div>

    <!-- Global Modals -->
    <GlobalModal />
    <PinSetupModal />
    <PinModal 
      v-model:isOpen="pinState.showVerifyModal"
      @success="onPinSuccess"
      @close="onPinClose"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { computed } from 'vue'
import { usePin } from '~/composables/usePin'
import PinModal from '~/components/common/PinModal.vue'

const authStore = useAuthStore()
const sidebarOpen = ref(false)
const showFabMenu = ref(false)
const { checkPinStatus, hasPin, showPinSetup, state: pinState, executePendingAction, closeModals } = usePin()

const onPinSuccess = (pin) => {
    executePendingAction(pin)
}

const onPinClose = () => {
    closeModals()
}

// Check PIN status on mount (only once per session)
// Check PIN status on mount
onMounted(async () => {
  try {
    // Always check PIN status to ensure state is synced with backend
    const hasPinSet = await checkPinStatus()
    
    // If user has NO PIN, prompt them to set it up
    if (!hasPinSet) {
      showPinSetup()
    }
  } catch (error) {
    console.error('Failed to check PIN status:', error)
  }
})

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

const closeSidebarOnMobile = () => {
  // Only close on mobile (smaller than lg breakpoint)
  if (window.innerWidth < 1024) {
    sidebarOpen.value = false
  }
}
</script>

<style scoped>
.nav-item {
  @apply flex items-center gap-3 px-4 py-3 text-sm font-medium text-muted rounded-xl transition-all duration-200 hover:bg-surface-hover hover:text-base;
}

.nav-item.active {
  @apply bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400;
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

