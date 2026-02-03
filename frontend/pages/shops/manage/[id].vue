<template>
  <ShopLayout>
    <div class="animate-fade-in-up">
      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mb-4"></div>
        <p class="text-gray-500 dark:text-gray-400">Chargement de votre boutique...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center min-h-[400px] text-center">
        <div class="text-4xl mb-4">⚠️</div>
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Une erreur est survenue</h2>
        <p class="text-gray-500 dark:text-gray-400 mb-6">{{ error }}</p>
        <NuxtLink to="/shops/my-shops" class="btn-primary">Retour aux boutiques</NuxtLink>
      </div>

      <div v-else-if="shop" class="space-y-8">
        <!-- Breadcrumb -->
        <nav class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
          <NuxtLink to="/shops" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Marketplace</NuxtLink>
          <span class="text-gray-300 dark:text-gray-600">/</span>
          <NuxtLink to="/shops/my-shops" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Mes Boutiques</NuxtLink>
          <span class="text-gray-300 dark:text-gray-600">/</span>
          <span class="font-medium text-gray-900 dark:text-white">{{ shop.name }}</span>
        </nav>

        <!-- Shop Header Banner -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl shadow-sm border border-gray-100 dark:border-white/5 overflow-hidden">
          <div class="relative h-48 bg-gradient-to-br from-indigo-500 to-violet-600">
            <div class="absolute inset-0 bg-gradient-to-b from-transparent to-black/30"></div>
            <img v-if="shop.banner_url" :src="shop.banner_url" class="w-full h-full object-cover" alt="Banner">
          </div>
          
          <div class="px-8 pb-8 -mt-10 relative flex flex-wrap gap-6 items-end justify-between">
            <div class="flex items-end gap-6">
              <div class="w-32 h-32 rounded-3xl border-4 border-white dark:border-slate-900 bg-white dark:bg-slate-900 shadow-xl overflow-hidden flex-shrink-0">
                <img v-if="shop.logo_url" :src="shop.logo_url" alt="Logo" class="w-full h-full object-cover">
                <div v-else class="w-full h-full flex items-center justify-center text-4xl font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/20">
                  {{ shop.name.charAt(0) }}
                </div>
              </div>
              
              <div class="pb-2 space-y-2">
                <div class="flex items-center gap-3 flex-wrap">
                  <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white">{{ shop.name }}</h1>
                  <span :class="['px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wide', getStatusClass(shop.status)]">
                    {{ getStatusLabel(shop.status) }}
                  </span>
                </div>
                <p class="text-gray-600 dark:text-gray-400 max-w-xl">{{ shop.description || 'Ajoutez une description pour vos clients' }}</p>
                <div class="flex items-center gap-3 text-sm text-gray-500 dark:text-gray-400 font-medium">
                  <span>📅 Créée le {{ formatDate(shop.created_at) }}</span>
                  <span class="text-gray-300 dark:text-gray-600">•</span>
                  <span>🌍 {{ shop.is_public ? 'Publique' : 'Privée' }}</span>
                </div>
              </div>
            </div>
            
            <div class="flex gap-3 mb-2 w-full sm:w-auto">
              <NuxtLink :to="`/shops/${shop.slug}`" target="_blank" class="flex-1 sm:flex-none flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-white/10 backdrop-blur-md border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-white font-semibold hover:bg-gray-50 dark:hover:bg-slate-800 transition-all">
                <EyeIcon class="w-5 h-5" />
                Voir la boutique
              </NuxtLink>
              <NuxtLink :to="`/shops/manage/${slug}/settings`" class="flex-1 sm:flex-none flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gray-900 dark:bg-white text-white dark:text-gray-900 font-bold hover:opacity-90 transition-all shadow-lg shadow-gray-200 dark:shadow-none">
                <Cog6ToothIcon class="w-5 h-5" />
                Paramètres
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-gradient-to-br from-indigo-600 to-violet-600 rounded-3xl p-6 text-white shadow-lg shadow-indigo-500/20 relative overflow-hidden group">
            <div class="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full blur-3xl group-hover:bg-white/20 transition-all"></div>
            <div class="w-12 h-12 rounded-xl bg-white/20 backdrop-blur-md flex items-center justify-center mb-4">
              <CurrencyEuroIcon class="w-6 h-6 text-white" />
            </div>
            <div class="relative z-10">
              <p class="text-indigo-100 font-medium mb-1">Chiffre d'affaires (Mois)</p>
              <p class="text-3xl font-extrabold mb-2">{{ formatPrice(shop.stats?.total_revenue || 0, shop.currency) }}</p>
              <span class="inline-flex items-center px-2.5 py-1 rounded-lg bg-green-500/20 text-green-100 text-xs font-bold">
                ↗ +12% vs mois dernier
              </span>
            </div>
          </div>

          <div class="bg-white dark:bg-slate-900 rounded-3xl p-6 border border-gray-100 dark:border-white/5 shadow-sm hover:shadow-md transition-shadow group">
            <div class="w-12 h-12 rounded-xl bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <ShoppingBagIcon class="w-6 h-6" />
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 font-medium mb-1">Commandes</p>
              <p class="text-3xl font-extrabold text-gray-900 dark:text-white mb-1">{{ shop.stats?.total_orders || 0 }}</p>
              <p class="text-xs text-gray-400 font-medium">0 en attente</p>
            </div>
          </div>

          <div class="bg-white dark:bg-slate-900 rounded-3xl p-6 border border-gray-100 dark:border-white/5 shadow-sm hover:shadow-md transition-shadow group">
            <div class="w-12 h-12 rounded-xl bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <TagIcon class="w-6 h-6" />
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 font-medium mb-1">Produits actifs</p>
              <p class="text-3xl font-extrabold text-gray-900 dark:text-white mb-1">{{ shop.stats?.total_products || 0 }}</p>
              <p class="text-xs text-gray-400 font-medium">Catalogue en ligne</p>
            </div>
          </div>
        </div>

        <!-- Main Content Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          
          <!-- Quick Actions -->
          <div class="bg-white dark:bg-slate-900 rounded-3xl border border-gray-100 dark:border-white/5 shadow-sm overflow-hidden">
            <div class="p-6 border-b border-gray-50 dark:border-gray-800">
              <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                <span>⚡</span> Actions Rapides
              </h3>
            </div>
            <div class="p-6 space-y-3">
              <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="flex items-center gap-4 p-4 rounded-2xl bg-gray-50 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 border border-transparent hover:border-gray-200 dark:hover:border-gray-700 transition-all group cursor-pointer shadow-sm hover:shadow-md">
                <div class="w-12 h-12 rounded-xl bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-500 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <PlusIcon class="w-6 h-6" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-900 dark:text-white">Ajouter un produit</h4>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Créez une nouvelle fiche produit</p>
                </div>
                <ChevronRightIcon class="w-5 h-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
              </NuxtLink>

              <NuxtLink :to="`/shops/manage/${slug}/categories`" class="flex items-center gap-4 p-4 rounded-2xl bg-gray-50 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 border border-transparent hover:border-gray-200 dark:hover:border-gray-700 transition-all group cursor-pointer shadow-sm hover:shadow-md">
                <div class="w-12 h-12 rounded-xl bg-pink-50 dark:bg-pink-500/10 text-pink-600 dark:text-pink-500 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <TagIcon class="w-6 h-6" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-900 dark:text-white">Gérer les catégories</h4>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Organisez vos rayons</p>
                </div>
                <ChevronRightIcon class="w-5 h-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
              </NuxtLink>

              <NuxtLink :to="`/shops/manage/${slug}/managers`" class="flex items-center gap-4 p-4 rounded-2xl bg-gray-50 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 border border-transparent hover:border-gray-200 dark:hover:border-gray-700 transition-all group cursor-pointer shadow-sm hover:shadow-md">
                <div class="w-12 h-12 rounded-xl bg-green-50 dark:bg-green-500/10 text-green-600 dark:text-green-500 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <UsersIcon class="w-6 h-6" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-900 dark:text-white">Gérer l'équipe</h4>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Ajoutez des vendeurs</p>
                </div>
                <ChevronRightIcon class="w-5 h-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
              </NuxtLink>
              
              <button @click="shareShop" class="w-full flex items-center gap-4 p-4 rounded-2xl bg-gray-50 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 border border-transparent hover:border-gray-200 dark:hover:border-gray-700 transition-all group cursor-pointer shadow-sm hover:shadow-md text-left">
                <div class="w-12 h-12 rounded-xl bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-500 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <ShareIcon class="w-6 h-6" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-900 dark:text-white">Partager la boutique</h4>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Copier le lien public</p>
                </div>
                <ChevronRightIcon class="w-5 h-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
              </button>
            </div>
          </div>

          <!-- Recent Orders -->
          <div class="bg-white dark:bg-slate-900 rounded-3xl border border-gray-100 dark:border-white/5 shadow-sm overflow-hidden flex flex-col h-full">
            <div class="p-6 border-b border-gray-50 dark:border-gray-800 flex justify-between items-center">
              <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                <span>📦</span> Dernières Commandes
              </h3>
              <NuxtLink :to="`/shops/manage/${slug}/orders`" class="text-sm font-bold text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300">Tout voir</NuxtLink>
            </div>
            
            <div class="flex-1 flex flex-col items-center justify-center p-12 text-center">
              <div class="text-6xl mb-4 opacity-50 grayscale">🛍️</div>
              <p class="text-gray-500 dark:text-gray-400 font-medium mb-2">Aucune commande récente</p>
              <span class="text-xs text-indigo-500 dark:text-indigo-400">Partagez votre boutique pour faire vos premières ventes !</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </ShopLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
// Removed relative import for ShopLayout as usually components are auto-imported or use #components alias. 
// However, assuming user structure:
import ShopLayout from '@/components/shops/ShopLayout.vue'
import { 
  BuildingStorefrontIcon, 
  CurrencyEuroIcon,
  ClipboardDocumentListIcon,
  ShoppingBagIcon,
  PlusIcon,
  TagIcon,
  UsersIcon,
  Cog6ToothIcon,
  EyeIcon,
  ShareIcon,
  ChevronRightIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const shop = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

// Animations on mount
onMounted(async () => {
  try {
    loading.value = true
    const response = await shopApi.getShop(slug)
    shop.value = response
  } catch (e: any) {
    console.error('Error fetching shop:', e)
    error.value = e.message || 'Impossible de charger la boutique'
  } finally {
    loading.value = false
  }
})

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: 'Actif',
    pending: 'En attente',
    suspended: 'Suspendu',
    draft: 'Brouillon'
  }
  return labels[status] || status
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    active: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    suspended: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    draft: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400'
  }
  return classes[status] || classes.draft
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('fr-FR', {
    day: 'numeric', 
    month: 'long', 
    year: 'numeric'
  })
}

const formatPrice = (price: number, currency = 'XOF') => {
  return new Intl.NumberFormat('fr-FR', { 
    style: 'currency', 
    currency: currency 
  }).format(price)
}

const shareShop = () => {
  if (navigator.share) {
    navigator.share({
      title: shop.value.name,
      text: shop.value.description,
      url: `${window.location.origin}/shops/${shop.value.slug}`
    })
  } else {
    // Fallback copy to clipboard
    navigator.clipboard.writeText(`${window.location.origin}/shops/${shop.value.slug}`)
    alert('Lien copié dans le presse-papier !')
  }
}

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>
