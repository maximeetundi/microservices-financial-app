<template>
  <div class="orders-page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">📦 Commandes</h1>
          <p class="page-subtitle">Gérez les commandes de votre boutique</p>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card pending">
          <div class="stat-icon">⏳</div>
          <div class="stat-info">
            <span class="stat-value">{{ pendingCount }}</span>
            <span class="stat-label">En attente</span>
          </div>
        </div>
        <div class="stat-card processing">
          <div class="stat-icon">📦</div>
          <div class="stat-info">
            <span class="stat-value">{{ processingCount }}</span>
            <span class="stat-label">En préparation</span>
          </div>
        </div>
        <div class="stat-card shipped">
          <div class="stat-icon">🚚</div>
          <div class="stat-info">
            <span class="stat-value">{{ shippedCount }}</span>
            <span class="stat-label">Expédiées</span>
          </div>
        </div>
        <div class="stat-card completed">
          <div class="stat-icon">✅</div>
          <div class="stat-info">
            <span class="stat-value">{{ deliveredCount }}</span>
            <span class="stat-label">Livrées</span>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-tabs">
          <button 
            v-for="tab in statusTabs" 
            :key="tab.key"
            @click="statusFilter = tab.key"
            :class="['filter-tab', statusFilter === tab.key && 'active']"
          >
            {{ tab.label }}
          </button>
        </div>
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Rechercher par numéro..."
            class="search-input"
          >
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des commandes...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredOrders.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <h3>Aucune commande</h3>
        <p>{{ statusFilter ? 'Aucune commande avec ce statut' : 'Partagez votre boutique pour recevoir des commandes' }}</p>
      </div>

      <!-- Orders List -->
      <div v-else class="orders-list">
        <div 
          v-for="order in filteredOrders" 
          :key="order.id" 
          class="order-card"
        >
          <div class="order-header">
            <div class="order-info">
              <span class="order-number">{{ order.order_number }}</span>
              <span :class="['status-badge', order.order_status]">
                {{ getStatusLabel(order.order_status) }}
              </span>
            </div>
            <div class="order-date">{{ formatDate(order.created_at) }}</div>
          </div>

          <div class="order-customer">
            <span class="customer-icon">👤</span>
            <span class="customer-name">{{ order.buyer_name || 'Client' }}</span>
          </div>

          <div class="order-items">
            <div v-for="item in order.items.slice(0, 3)" :key="item.product_id" class="item-preview">
              <div class="item-image">
                <img v-if="item.product_image" :src="item.product_image" alt="">
                <span v-else>📦</span>
              </div>
              <span class="item-qty">x{{ item.quantity }}</span>
            </div>
            <div v-if="order.items.length > 3" class="item-more">
              +{{ order.items.length - 3 }}
            </div>
          </div>

          <div class="order-footer">
            <div class="order-total">
              <span class="total-label">Total</span>
              <span class="total-amount">{{ formatPrice(order.total_amount, order.currency) }}</span>
            </div>

            <div class="order-actions">
              <!-- Quick Actions based on status -->
              <button 
                v-if="order.order_status === 'pending'"
                @click="updateStatus(order, 'confirmed')"
                class="btn-action confirm"
                :disabled="updating === order.id"
              >
                ✓ Confirmer
              </button>
              <button 
                v-if="order.order_status === 'confirmed'"
                @click="updateStatus(order, 'processing')"
                class="btn-action process"
                :disabled="updating === order.id"
              >
                📦 Préparer
              </button>
              <button 
                v-if="order.order_status === 'processing'"
                @click="updateStatus(order, 'shipped')"
                class="btn-action ship"
                :disabled="updating === order.id"
              >
                🚚 Expédier
              </button>
              <button 
                v-if="order.order_status === 'shipped'"
                @click="updateStatus(order, 'delivered')"
                class="btn-action deliver"
                :disabled="updating === order.id"
              >
                ✅ Livrée
              </button>
              <button @click="viewDetails(order)" class="btn-action details">
                Détails
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button 
          v-for="p in totalPages" 
          :key="p"
          @click="goToPage(p)"
          :class="['page-btn', p === page && 'active']"
        >
          {{ p }}
        </button>
      </div>

      <!-- Order Detail Modal -->
      <Teleport to="body">
        <div v-if="selectedOrder" class="modal-overlay" @click.self="selectedOrder = null">
          <div class="modal-content">
            <button class="close-btn" @click="selectedOrder = null">✕</button>
            
            <h2 class="modal-title">{{ selectedOrder.order_number }}</h2>
            
            <div class="detail-section">
              <h4>Client</h4>
              <p>{{ selectedOrder.buyer_name || 'Non spécifié' }}</p>
            </div>

            <div class="detail-section">
              <h4>Articles</h4>
              <div v-for="item in selectedOrder.items" :key="item.product_id" class="detail-item">
                <span>{{ item.product_name }} x{{ item.quantity }}</span>
                <span>{{ formatPrice(item.total_price, selectedOrder.currency) }}</span>
              </div>
            </div>

            <div class="detail-section">
              <h4>Livraison</h4>
              <p>{{ selectedOrder.delivery_type === 'pickup' ? '🏃 Retrait' : '📦 Livraison' }}</p>
              <p v-if="selectedOrder.shipping_address">
                {{ selectedOrder.shipping_address.street }}, {{ selectedOrder.shipping_address.city }}
              </p>
            </div>

            <div class="detail-section">
              <h4>Total</h4>
              <p class="total-large">{{ formatPrice(selectedOrder.total_amount, selectedOrder.currency) }}</p>
            </div>

            <div v-if="selectedOrder.notes" class="detail-section">
              <h4>Notes</h4>
              <p class="notes">{{ selectedOrder.notes }}</p>
            </div>

            <div class="modal-actions">
              <select v-model="newStatus" class="status-select">
                <option value="">Changer le statut...</option>
                <option value="confirmed">✓ Confirmée</option>
                <option value="processing">📦 En préparation</option>
                <option value="shipped">🚚 Expédiée</option>
                <option value="delivered">✅ Livrée</option>
                <option value="cancelled">❌ Annulée</option>
              </select>
              <button 
                v-if="newStatus" 
                @click="applyStatusChange" 
                class="btn-apply"
                :disabled="updating === selectedOrder.id"
              >
                Appliquer
              </button>
            </div>
          </div>
        </div>
      </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type Order } from '@/composables/useShopApi'
import ShopLayout from '@/components/shops/ShopLayout.vue'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const orders = ref<Order[]>([])
const loading = ref(true)
const updating = ref<string | null>(null)
const statusFilter = ref('')
const searchQuery = ref('')
const page = ref(1)
const totalPages = ref(1)
const selectedOrder = ref<Order | null>(null)
const newStatus = ref('')
const shopId = ref('')

const statusTabs = [
  { key: '', label: 'Toutes' },
  { key: 'pending', label: '⏳ En attente' },
  { key: 'confirmed', label: '✓ Confirmées' },
  { key: 'processing', label: '📦 En préparation' },
  { key: 'shipped', label: '🚚 Expédiées' },
  { key: 'delivered', label: '✅ Livrées' },
]

const pendingCount = computed(() => orders.value.filter(o => o.order_status === 'pending').length)
const processingCount = computed(() => orders.value.filter(o => o.order_status === 'processing' || o.order_status === 'confirmed').length)
const shippedCount = computed(() => orders.value.filter(o => o.order_status === 'shipped').length)
const deliveredCount = computed(() => orders.value.filter(o => o.order_status === 'delivered').length)

const filteredOrders = computed(() => {
  let result = orders.value
  if (statusFilter.value) {
    result = result.filter(o => o.order_status === statusFilter.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(o => 
      o.order_number.toLowerCase().includes(q) ||
      o.buyer_name?.toLowerCase().includes(q)
    )
  }
  return result
})

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '⏳ En attente',
    confirmed: '✓ Confirmée',
    processing: '📦 Préparation',
    shipped: '🚚 Expédiée',
    delivered: '✅ Livrée',
    cancelled: '❌ Annulée',
  }
  return labels[status] || status
}

const formatPrice = (price: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(price)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

const fetchOrders = async () => {
  try {
    loading.value = true
    
    // Get shop ID first
    if (!shopId.value) {
      const shop = await shopApi.getShop(slug)
      shopId.value = shop.id
    }

    const response = await shopApi.listShopOrders(shopId.value, page.value, 50)
    orders.value = response.orders || []
    totalPages.value = response.total_pages || 1
  } catch (e: any) {
    console.error('Error fetching orders:', e)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (order: Order, status: string) => {
  try {
    updating.value = order.id
    await shopApi.updateOrderStatus(order.id, status)
    order.order_status = status
  } catch (e: any) {
    alert('Erreur: ' + (e.message || 'Impossible de mettre à jour'))
  } finally {
    updating.value = null
  }
}

const viewDetails = (order: Order) => {
  selectedOrder.value = order
  newStatus.value = ''
}

const applyStatusChange = async () => {
  if (!selectedOrder.value || !newStatus.value) return
  await updateStatus(selectedOrder.value, newStatus.value)
  newStatus.value = ''
}

const goToPage = (p: number) => {
  page.value = p
  fetchOrders()
}

onMounted(fetchOrders)

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>

<style scoped>
.orders-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
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

/* Stats */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
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
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-card.pending { border-left: 4px solid #f59e0b; }
.stat-card.processing { border-left: 4px solid #8b5cf6; }
.stat-card.shipped { border-left: 4px solid #06b6d4; }
.stat-card.completed { border-left: 4px solid #10b981; }

.stat-icon {
  font-size: 28px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
}

.stat-label {
  color: var(--text-muted, #6b7280);
  font-size: 13px;
}

/* Filters */
.filters-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
}

.filter-tab {
  padding: 8px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 20px;
  background: var(--surface, white);
  color: var(--text-muted, #6b7280);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.filter-tab:hover {
  border-color: #6366f1;
  color: #6366f1;
}

.filter-tab.active {
  background: #6366f1;
  border-color: #6366f1;
  color: white;
}

.search-box {
  position: relative;
  min-width: 200px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
}

.search-input {
  width: 100%;
  padding: 10px 16px 10px 40px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 14px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
}

.search-input:focus {
  outline: none;
  border-color: #6366f1;
}

/* Orders List */
.orders-list {
  display: grid;
  gap: 16px;
}

.order-card {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  padding: 20px;
  transition: all 0.2s;
}

.order-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.order-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-number {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
}

.status-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.pending { background: #fef3c7; color: #92400e; }
.status-badge.confirmed { background: #dbeafe; color: #1d4ed8; }
.status-badge.processing { background: #ede9fe; color: #7c3aed; }
.status-badge.shipped { background: #cffafe; color: #0891b2; }
.status-badge.delivered { background: #d1fae5; color: #059669; }
.status-badge.cancelled { background: #fee2e2; color: #dc2626; }

:global(.dark) .status-badge.pending { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
:global(.dark) .status-badge.confirmed { background: rgba(59, 130, 246, 0.2); color: #60a5fa; }
:global(.dark) .status-badge.processing { background: rgba(139, 92, 246, 0.2); color: #a78bfa; }
:global(.dark) .status-badge.shipped { background: rgba(6, 182, 212, 0.2); color: #22d3ee; }
:global(.dark) .status-badge.delivered { background: rgba(16, 185, 129, 0.2); color: #34d399; }

.order-date {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
}

.order-customer {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-primary, #1f2937);
  margin-bottom: 12px;
}

.order-items {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.item-preview {
  position: relative;
  width: 48px;
  height: 48px;
}

.item-image {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  background: var(--surface-hover, #f3f4f6);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-qty {
  position: absolute;
  bottom: -4px;
  right: -4px;
  background: #6366f1;
  color: white;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 10px;
}

.item-more {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: var(--surface-hover, #f3f4f6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted, #6b7280);
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid var(--border, #e5e7eb);
}

.total-label {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
  display: block;
}

.total-amount {
  font-size: 20px;
  font-weight: 700;
  color: #6366f1;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.btn-action {
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-action.confirm { background: #dbeafe; color: #1d4ed8; }
.btn-action.confirm:hover { background: #bfdbfe; }

.btn-action.process { background: #ede9fe; color: #7c3aed; }
.btn-action.process:hover { background: #ddd6fe; }

.btn-action.ship { background: #cffafe; color: #0891b2; }
.btn-action.ship:hover { background: #a5f3fc; }

.btn-action.deliver { background: #d1fae5; color: #059669; }
.btn-action.deliver:hover { background: #a7f3d0; }

.btn-action.details { 
  background: var(--surface-hover, #f3f4f6); 
  color: var(--text-primary, #1f2937); 
}
.btn-action.details:hover { background: #e5e7eb; }

.btn-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 32px;
}

.page-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  border: 1px solid var(--border, #e5e7eb);
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn.active {
  background: #6366f1;
  border-color: #6366f1;
  color: white;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.modal-content {
  background: var(--surface, white);
  border-radius: 20px;
  width: 100%;
  max-width: 500px;
  padding: 32px;
  position: relative;
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text-muted, #6b7280);
}

.modal-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: var(--text-primary, #1f2937);
}

.detail-section {
  margin-bottom: 20px;
}

.detail-section h4 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted, #6b7280);
  margin: 0 0 8px 0;
  text-transform: uppercase;
}

.detail-section p {
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--border, #e5e7eb);
  font-size: 14px;
}

.total-large {
  font-size: 24px;
  font-weight: 700;
  color: #6366f1;
}

.notes {
  background: var(--surface-hover, #f9fafb);
  padding: 12px;
  border-radius: 8px;
  font-style: italic;
}

.modal-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.status-select {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 14px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
}

.btn-apply {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
}

.btn-apply:disabled {
  opacity: 0.5;
}

/* Dark Mode */
:global(.dark) .order-card,
:global(.dark) .stat-card {
  background: #1e293b;
  border-color: #334155;
}

:global(.dark) .modal-content {
  background: #1e293b;
}

:global(.dark) .filter-tab {
  background: #1e293b;
  border-color: #334155;
}

@media (max-width: 768px) {
  .filters-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-tabs {
    width: 100%;
  }
  
  .search-box {
    width: 100%;
  }
  
  .order-footer {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .order-actions {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>
