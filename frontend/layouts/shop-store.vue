<template>
  <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-slate-950 transition-colors duration-300">
    <!-- Navbar -->
    <header class="sticky top-0 z-40 w-full bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between gap-4">
          <!-- Logo -->
          <NuxtLink to="/shops" class="flex items-center gap-3 flex-shrink-0 group">
             <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-xl shadow-lg group-hover:scale-105 transition-transform">Z</div>
             <div>
                <h1 class="font-bold text-xl bg-gradient-to-r from-indigo-600 to-purple-600 dark:from-indigo-400 dark:to-purple-400 bg-clip-text text-transparent">Zekora</h1>
                <p class="text-[10px] text-gray-500 dark:text-gray-400 font-medium uppercase tracking-wider">Marketplace</p>
             </div>
          </NuxtLink>

          <!-- Search Bar (Desktop) -->
          <div class="hidden md:flex flex-1 max-w-lg mx-auto">
             <div class="relative w-full group">
                <input 
                  v-model="searchQuery"
                  @keyup.enter="handleSearch"
                  type="text" 
                  placeholder="Rechercher un produit..." 
                  class="w-full pl-10 pr-4 py-2.5 bg-gray-100 dark:bg-slate-800 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 dark:text-white transition-all group-hover:bg-white dark:group-hover:bg-slate-700 shadow-sm"
                >
                <button @click="handleSearch" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 group-hover:text-indigo-500 transition-colors">🔍</button>
             </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-3 sm:gap-4">
             <button class="md:hidden p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-slate-800 rounded-xl transition-colors">
               <span class="text-xl">🔍</span>
             </button>

             <ThemeToggle />
             
             <!-- Cart -->
             <NuxtLink to="/cart" class="relative p-2.5 text-gray-500 dark:text-gray-400 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600 dark:hover:text-indigo-400 rounded-xl transition-all">
               <span class="text-xl">🛒</span>
               <span v-if="cartStore.itemCount > 0" class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center border-2 border-white dark:border-slate-900 shadow-sm animate-bounce">
                 {{ cartStore.itemCount }}
               </span>
             </NuxtLink>

             <!-- User Menu -->
             <div class="relative group">
                <button class="flex items-center gap-2 p-1.5 rounded-full border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 hover:shadow-md hover:border-indigo-200 dark:hover:border-indigo-800 transition-all">
                   <div class="w-8 h-8 rounded-full bg-indigo-100 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-bold text-sm">
                     {{ userInitials }}
                   </div>
                   <span class="hidden sm:block text-sm font-medium text-gray-700 dark:text-gray-200 pr-2 max-w-[100px] truncate">
                     {{ userName }}
                   </span>
                </button>
                
                <!-- Dropdown -->
                <div class="absolute right-0 top-full mt-2 w-48 bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 opacity-0 translate-y-2 group-hover:opacity-100 group-hover:translate-y-0 invisible group-hover:visible transition-all transform origin-top-right z-50">
                   <div class="p-2 space-y-1">
                      <NuxtLink to="/dashboard" class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700 rounded-lg">
                        <span>📊</span> Tableau de bord
                      </NuxtLink>
                      <NuxtLink to="/orders" class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700 rounded-lg">
                        <span>📦</span> Mes commandes
                      </NuxtLink>
                       <div class="h-px bg-gray-100 dark:bg-gray-700 my-1"></div>
                      <button @click="authStore.logout()" class="w-full text-left flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg">
                        <span>🚪</span> Déconnexion
                      </button>
                   </div>
                </div>
             </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
       <div class="animate-fade-in">
         <slot />
       </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white dark:bg-slate-900 border-t border-gray-200 dark:border-gray-800 py-12 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div class="col-span-1 md:col-span-2">
            <h2 class="font-bold text-2xl text-gray-900 dark:text-white mb-4 flex items-center gap-2">
               <span class="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white text-sm">Z</span>
               Zekora Marketplace
            </h2>
            <p class="text-gray-500 dark:text-gray-400 max-w-md text-sm leading-relaxed">
              Découvrez des produits uniques de vendeurs de confiance, locaux et internationaux. 
              La meilleure expérience d'achat en ligne avec paiements sécurisés et crypto.
            </p>
          </div>
          <div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-4">Liens Utiles</h3>
            <ul class="space-y-2 text-sm">
              <li><NuxtLink to="/shops" class="text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Toutes les boutiques</NuxtLink></li>
              <li><NuxtLink to="/about" class="text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">À propos</NuxtLink></li>
              <li><NuxtLink to="/support" class="text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Support & Aide</NuxtLink></li>
            </ul>
          </div>
          <div>
             <h3 class="font-bold text-gray-900 dark:text-white mb-4">Légal</h3>
            <ul class="space-y-2 text-sm">
              <li><NuxtLink to="/terms" class="text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Conditions d'utilisation</NuxtLink></li>
              <li><NuxtLink to="/privacy" class="text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Confidentialité</NuxtLink></li>
            </ul>
          </div>
        </div>
        <div class="pt-8 border-t border-gray-100 dark:border-gray-800 text-center text-gray-400 text-xs">
           &copy; {{ new Date().getFullYear() }} Zekora. Tous droits réservés. Fièrement propulsé par Zekora Pay.
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useCartStore } from '~/stores/cart'

const authStore = useAuthStore()
const cartStore = useCartStore()
const route = useRoute()
const router = useRouter()
const searchQuery = ref('')

// Initialize search from URL
searchQuery.value = (route.query.search as string) || ''

const userInitials = computed(() => {
  if (authStore.user) {
    return `${authStore.user.first_name?.[0] || ''}${authStore.user.last_name?.[0] || ''}`
  }
  return 'G'
})

const userName = computed(() => {
    if (authStore.user) return `${authStore.user.first_name}`
    return 'Guest'
})

const handleSearch = () => {
    const slug = route.params.slug as string
    if (slug) {
        router.push({ 
            path: `/shops/${slug}`, 
            query: { ...route.query, search: searchQuery.value || undefined } 
        })
    }
}

onMounted(() => {
    cartStore.loadFromStorage()
})

</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
