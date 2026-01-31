<template>
  <div class="shop-dashboard">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement de votre boutique...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h2>Une erreur est survenue</h2>
      <p>{{ error }}</p>
      <NuxtLink to="/shops/my-shops" class="btn-primary">Retour aux boutiques</NuxtLink>
    </div>

    <div v-else-if="shop" class="dashboard-content">
      <!-- Breadcrumb -->
      <nav class="breadcrumb">
        <NuxtLink to="/shops" class="breadcrumb-item">Marketplace</NuxtLink>
        <span class="separator">/</span>
        <NuxtLink to="/shops/my-shops" class="breadcrumb-item">Mes Boutiques</NuxtLink>
        <span class="separator">/</span>
        <span class="current">{{ shop.name }}</span>
      </nav>

      <!-- Shop Header Banner -->
      <div class="shop-header">
        <div class="header-bg">
          <div class="gradient-overlay"></div>
          <img v-if="shop.banner_url" :src="shop.banner_url" class="banner-img" alt="Banner">
        </div>
        
        <div class="header-content">
          <div class="shop-identity">
            <div class="shop-logo">
              <img v-if="shop.logo_url" :src="shop.logo_url" alt="Logo">
              <div v-else class="logo-placeholder">{{ shop.name.charAt(0) }}</div>
            </div>
            <div class="shop-text">
              <div class="flex items-center gap-3">
                <h1>{{ shop.name }}</h1>
                <span :class="['status-badge', shop.status]">{{ getStatusLabel(shop.status) }}</span>
              </div>
              <p class="description">{{ shop.description || 'Ajoutez une description pour vos clients' }}</p>
              <div class="meta-info">
                <span>📅 Créée le {{ formatDate(shop.created_at) }}</span>
                <span class="divider">•</span>
                <span>🌍 {{ shop.is_public ? 'Publique' : 'Privée' }}</span>
              </div>
            </div>
          </div>
          
          <div class="header-actions">
            <NuxtLink :to="`/shops/${shop.slug}`" target="_blank" class="btn-glass">
              <EyeIcon class="w-5 h-5" />
              Voir la boutique
            </NuxtLink>
            <NuxtLink :to="`/shops/manage/${slug}/settings`" class="btn-primary">
              <Cog6ToothIcon class="w-5 h-5" />
              Paramètres
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="stats-grid">
        <div class="stat-card primary">
          <div class="stat-icon-wrapper">
            <CurrencyEuroIcon class="w-6 h-6" />
          </div>
          <div class="stat-details">
            <span class="stat-label">Chiffre d'affaires (Mois)</span>
            <span class="stat-value">{{ formatPrice(shop.stats?.total_revenue || 0, shop.currency) }}</span>
            <span class="stat-trend positive">↗ +12% vs mois dernier</span>
          </div>
          <div class="stat-chart-placeholder"></div>
        </div>

        <div class="stat-card">
          <div class="stat-icon-wrapper blue">
            <shopping-bag-icon class="w-6 h-6" />
          </div>
          <div class="stat-details">
            <span class="stat-label">Commandes</span>
            <span class="stat-value">{{ shop.stats?.total_orders || 0 }}</span>
            <span class="stat-sub">0 en attente</span>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon-wrapper purple">
            <TagIcon class="w-6 h-6" />
          </div>
          <div class="stat-details">
            <span class="stat-label">Produits actifs</span>
            <span class="stat-value">{{ shop.stats?.total_products || 0 }}</span>
            <span class="stat-sub">Catalogue en ligne</span>
          </div>
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="content-grid">
        
        <!-- Quick Actions -->
        <div class="dashboard-card actions-card">
          <div class="card-header">
            <h3>⚡ Actions Rapides</h3>
          </div>
          <div class="actions-list">
            <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="action-item">
              <div class="action-icon orange"><PlusIcon class="w-6 h-6" /></div>
              <div class="action-text">
                <h4>Ajouter un produit</h4>
                <p>Créez une nouvelle fiche produit</p>
              </div>
              <ChevronRightIcon class="arrow-icon" />
            </NuxtLink>

            <NuxtLink :to="`/shops/manage/${slug}/categories`" class="action-item">
              <div class="action-icon pink"><TagIcon class="w-6 h-6" /></div>
              <div class="action-text">
                <h4>Gérer les catégories</h4>
                <p>Organisez vos rayons</p>
              </div>
              <ChevronRightIcon class="arrow-icon" />
            </NuxtLink>

            <NuxtLink :to="`/shops/manage/${slug}/managers`" class="action-item">
              <div class="action-icon green"><UsersIcon class="w-6 h-6" /></div>
              <div class="action-text">
                <h4>Gérer l'équipe</h4>
                <p>Ajoutez des vendeurs</p>
              </div>
              <ChevronRightIcon class="arrow-icon" />
            </NuxtLink>
            
            <button class="action-item" @click="shareShop">
              <div class="action-icon blue"><ShareIcon class="w-6 h-6" /></div>
              <div class="action-text">
                <h4>Partager la boutique</h4>
                <p>Copier le lien public</p>
              </div>
              <ChevronRightIcon class="arrow-icon" />
            </button>
          </div>
        </div>

        <!-- Recent Orders -->
        <div class="dashboard-card orders-card">
          <div class="card-header">
            <h3>📦 Dernières Commandes</h3>
            <NuxtLink :to="`/shops/manage/${slug}/orders`" class="link-more">Tout voir</NuxtLink>
          </div>
          
          <div class="empty-orders">
            <div class="empty-illustration">🛍️</div>
            <p>Aucune commande récente</p>
            <span class="empty-hint">Partagez votre boutique pour faire vos premières ventes !</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi } from '@/composables/useShopApi'
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
  layout: 'shop-admin'
})
</script>

<style scoped>
.shop-dashboard {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  animation: fadeIn 0.5s ease-out;
}

/* Header & Banner */
.shop-header {
  position: relative;
  border-radius: 24px;
  overflow: hidden;
  background: white;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  margin-bottom: 32px;
}

.header-bg {
  height: 200px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  position: relative;
}

.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gradient-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0) 0%, rgba(0,0,0,0.4) 100%);
}

.header-content {
  padding: 0 32px 32px;
  margin-top: -60px;
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 24px;
}

.shop-identity {
  display: flex;
  gap: 24px;
  align-items: flex-end;
}

.shop-logo {
  width: 120px;
  height: 120px;
  border-radius: 24px;
  border: 4px solid white;
  background: white;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  overflow: hidden;
  flex-shrink: 0;
}

.shop-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.logo-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 800;
  color: #6366f1;
  background: #eff6ff;
}

.shop-text h1 {
  font-size: 32px;
  font-weight: 800;
  color: #111827;
  letter-spacing: -0.5px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.active { background: #dcfce7; color: #166534; }
.status-badge.pending { background: #fef9c3; color: #854d0e; }
.status-badge.suspended { background: #fee2e2; color: #991b1b; }

.description {
  color: #6b7280;
  margin-top: 8px;
  max-width: 600px;
  line-height: 1.5;
}

.meta-info {
  display: flex;
  gap: 8px;
  color: #9ca3af;
  font-size: 14px;
  margin-top: 8px;
}

/* Stats */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border: 1px solid #f3f4f6;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.stat-card.primary {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.stat-card.primary .stat-label, 
.stat-card.primary .stat-sub { color: rgba(255,255,255,0.8); }
.stat-card.primary .stat-value { color: white; }
.stat-card.primary .stat-icon-wrapper { background: rgba(255,255,255,0.2); color: white; }
.stat-card.primary .stat-trend { background: rgba(255,255,255,0.2); color: white; }

.stat-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  background: #eff6ff;
  color: #6366f1;
}

.stat-icon-wrapper.blue { background: #eff6ff; color: #3b82f6; }
.stat-icon-wrapper.purple { background: #f3f0ff; color: #a855f7; }

.stat-details {
  display: flex;
  flex-direction: column;
}

.stat-label { font-size: 14px; color: #6b7280; font-weight: 500; }
.stat-value { font-size: 32px; font-weight: 800; margin: 4px 0; color: #111827; }
.stat-sub { font-size: 13px; color: #9ca3af; }

.stat-trend {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  margin-top: 8px;
  background: #dcfce7;
  color: #166534;
  width: fit-content;
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.dashboard-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  border: 1px solid #f3f4f6;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.action-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border: 1px solid transparent;
  border-radius: 16px;
  margin-bottom: 12px;
  transition: all 0.2s;
  cursor: pointer;
  width: 100%;
  text-align: left;
}

.action-item:hover {
  background: white;
  border-color: #e5e7eb;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  transform: translateX(4px);
}

.action-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 20px;
}

.action-icon.orange { background: #fff7ed; color: #f97316; }
.action-icon.pink { background: #fdf2f8; color: #db2777; }
.action-icon.green { background: #f0fdf4; color: #16a34a; }
.action-icon.blue { background: #eff6ff; color: #2563eb; }

.action-text h4 { font-weight: 600; color: #1f2937; margin: 0; }
.action-text p { font-size: 13px; color: #6b7280; margin: 2px 0 0; }
.arrow-icon { width: 16px; color: #d1d5db; margin-left: auto; }

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 24px;
  font-size: 14px;
  color: #6b7280;
}

.breadcrumb-item { color: #6b7280; text-decoration: none; transition: color 0.2s; }
.breadcrumb-item:hover { color: #6366f1; }
.current { color: #111827; font-weight: 500; }
.separator { color: #d1d5db; }

/* Buttons */
.btn-primary, .btn-glass {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.2s;
  cursor: pointer;
}

.btn-primary {
  background: linear-gradient(135deg, #111827 0%, #374151 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.3);
}

.btn-glass {
  background: rgba(255,255,255,0.8);
  backdrop-filter: blur(10px);
  color: #1f2937;
  border: 1px solid rgba(0,0,0,0.05);
}

.btn-glass:hover {
  background: white;
}

/* Empty State */
.empty-orders {
  text-align: center;
  padding: 40px;
  color: #9ca3af;
}

.empty-illustration { font-size: 48px; margin-bottom: 16px; opacity: 0.5; }
.empty-hint { display: block; margin-top: 8px; font-size: 13px; color: #6366f1; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 1024px) {
  .content-grid { grid-template-columns: 1fr; }
  .header-content { flex-direction: column; align-items: flex-start; }
  .header-actions { width: 100%; display: flex; gap: 12px; }
  .btn-primary, .btn-glass { flex: 1; justify-content: center; }
}
</style>
