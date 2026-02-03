<template>
  <div class="products-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">📦 Produits</h1>
        <p class="page-subtitle">Gérez le catalogue de votre boutique</p>
      </div>
      <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="btn-create">
        <span class="icon">+</span>
        Nouveau produit
      </NuxtLink>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <span class="stat-value">{{ products.length }}</span>
          <span class="stat-label">Produits</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-info">
          <span class="stat-value">{{ activeProducts }}</span>
          <span class="stat-label">Actifs</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-info">
          <span class="stat-value">{{ totalStock }}</span>
          <span class="stat-label">En stock</span>
        </div>
      </div>
    </div>

    <!-- Search & Filters -->
    <div class="filters-bar">
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Rechercher un produit..."
          class="search-input"
          @input="filterProducts"
        >
      </div>
      <select v-model="statusFilter" class="filter-select" @change="filterProducts">
        <option value="">Tous les statuts</option>
        <option value="active">Actifs</option>
        <option value="draft">Brouillons</option>
        <option value="archived">Archivés</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des produits...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredProducts.length === 0" class="empty-state">
      <div class="empty-icon">📦</div>
      <h3>Aucun produit</h3>
      <p>Commencez par ajouter votre premier produit à la boutique</p>
      <NuxtLink :to="`/shops/manage/${slug}/products/create`" class="btn-primary">
        Créer un produit
      </NuxtLink>
    </div>

    <!-- Products Grid -->
    <div v-else class="products-grid">
      <div 
        v-for="product in filteredProducts" 
        :key="product.id" 
        class="product-card"
      >
        <!-- Product Image -->
        <div class="product-image">
          <img 
            v-if="product.images && product.images.length > 0" 
            :src="product.images[0]" 
            :alt="product.name"
          >
          <div v-else class="placeholder-image">📦</div>
          
          <!-- Status Badge -->
          <span 
            class="status-badge"
            :class="product.status"
          >
            {{ getStatusLabel(product.status) }}
          </span>

          <!-- Featured Badge -->
          <span v-if="product.is_featured" class="featured-badge">⭐ En vedette</span>
        </div>

        <!-- Product Content -->
        <div class="product-content">
          <h3 class="product-name">{{ product.name }}</h3>
          
          <div class="product-price">
            <span class="price-current">{{ formatPrice(product.price, product.currency) }}</span>
            <span v-if="product.compare_at_price && product.compare_at_price > product.price" class="price-compare">
              {{ formatPrice(product.compare_at_price, product.currency) }}
            </span>
          </div>

          <div class="product-meta">
            <span class="meta-item">📦 Stock: {{ product.stock }}</span>
            <span class="meta-item">👁️ {{ product.view_count || 0 }} vues</span>
            <span class="meta-item">🛒 {{ product.sold_count || 0 }} vendus</span>
          </div>
        </div>

        <!-- Product Actions -->
        <div class="product-actions">
          <NuxtLink 
            :to="`/shops/manage/${slug}/products/${product.id}`" 
            class="btn-action edit"
            title="Modifier"
          >
            ✏️ Modifier
          </NuxtLink>
          <button @click="deleteProduct(product)" class="btn-action delete" title="Supprimer">
            🗑️
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type Product } from '@/composables/useShopApi'
import ShopLayout from '@/components/shops/ShopLayout.vue'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const products = ref<Product[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const searchQuery = ref('')
const statusFilter = ref('')

const activeProducts = computed(() => 
  products.value.filter(p => p.status === 'active').length
)

const totalStock = computed(() => 
  products.value.reduce((sum, p) => sum + (p.stock || 0), 0)
)

const filteredProducts = computed(() => {
  let result = products.value
  
  if (statusFilter.value) {
    result = result.filter(p => p.status === statusFilter.value)
  }
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(p => 
      p.name.toLowerCase().includes(query) ||
      p.description?.toLowerCase().includes(query)
    )
  }
  
  return result
})

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: 'Actif',
    draft: 'Brouillon',
    archived: 'Archivé'
  }
  return labels[status] || status
}

const formatPrice = (price: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency || 'XOF'
  }).format(price)
}

const fetchProducts = async () => {
  try {
    loading.value = true
    const response = await shopApi.listProducts(slug)
    products.value = response.products || []
  } catch (e: any) {
    console.error('Error fetching products:', e)
    error.value = e.message || 'Impossible de charger les produits'
  } finally {
    loading.value = false
  }
}

const deleteProduct = async (product: Product) => {
  if (!confirm(`Voulez-vous vraiment supprimer le produit "${product.name}" ?`)) return

  try {
    await shopApi.deleteProduct(product.id)
    products.value = products.value.filter(p => p.id !== product.id)
  } catch (e: any) {
    alert('Erreur: ' + e.message)
  }
}

const filterProducts = () => {
  // Filtering is handled by computed property
}

onMounted(fetchProducts)

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>

<style scoped>
.products-page {
  padding: 24px;
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.page-subtitle {
  color: var(--text-muted, #6b7280);
  margin-top: 4px;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  padding: 12px 24px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 28px;
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 12px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
}

.stat-label {
  color: var(--text-muted, #6b7280);
  font-size: 14px;
}

/* Filters */
.filters-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 250px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 44px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 15px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
}

.search-input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.filter-select {
  padding: 12px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 15px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
  cursor: pointer;
}

/* Products Grid */
.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.product-card {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.product-image {
  position: relative;
  height: 180px;
  background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%);
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder-image {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
}

.status-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.active {
  background: #10b981;
  color: white;
}

.status-badge.draft {
  background: #6b7280;
  color: white;
}

.status-badge.archived {
  background: #f59e0b;
  color: white;
}

.featured-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
  border-radius: 16px;
  font-size: 11px;
  font-weight: 600;
}

.product-content {
  padding: 16px;
}

.product-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin: 0 0 8px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-price {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.price-current {
  font-size: 18px;
  font-weight: 700;
  color: #6366f1;
}

.price-compare {
  font-size: 14px;
  color: var(--text-muted, #6b7280);
  text-decoration: line-through;
}

.product-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.meta-item {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
}

.product-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border, #e5e7eb);
  background: var(--surface-hover, #f9fafb);
}

.btn-action {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
  text-decoration: none;
  text-align: center;
}

.btn-action.edit {
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;
}

.btn-action.edit:hover {
  background: rgba(99, 102, 241, 0.2);
}

.btn-action.delete {
  flex: 0;
  padding: 10px 14px;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.btn-action.delete:hover {
  background: rgba(239, 68, 68, 0.2);
}

/* Empty & Loading */
.empty-state, .loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted, #6b7280);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.btn-primary {
  display: inline-block;
  margin-top: 16px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border-radius: 12px;
  font-weight: 600;
  text-decoration: none;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border, #e5e7eb);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .products-grid {
    grid-template-columns: 1fr;
  }
  
  .filters-bar {
    flex-direction: column;
  }
}
</style>
