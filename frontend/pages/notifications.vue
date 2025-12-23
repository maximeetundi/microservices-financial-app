<template>
  <NuxtLayout name="dashboard">
    <div class="notifications-page">
      <!-- Header Section with Gradient Background -->
      <div class="header-section">
        <div class="header-bg"></div>
        <div class="header-content">
          <div class="header-icon-wrapper">
            <div class="header-icon">
              <svg viewBox="0 0 24 24" fill="none" class="w-8 h-8 text-white">
                <path d="M18 8A6 6 0 106 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M13.73 21a2 2 0 01-3.46 0" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span v-if="unreadCount > 0" class="header-badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
            </div>
          </div>
          <div class="header-text">
            <h1>Centre de Notifications</h1>
            <p v-if="unreadCount > 0">{{ unreadCount }} nouvelle{{ unreadCount > 1 ? 's' : '' }} notification{{ unreadCount > 1 ? 's' : '' }}</p>
            <p v-else>Vous êtes à jour ✨</p>
          </div>
          <div class="header-actions">
            <button 
              v-if="unreadCount > 0"
              @click="markAllAsRead"
              class="action-btn primary"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              Tout lire
            </button>
            <button 
              @click="fetchNotifications"
              :disabled="loading"
              class="action-btn secondary"
              :class="{ 'loading': loading }"
            >
              <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="stats-section">
        <div class="stats-grid">
          <div class="stat-card" v-for="stat in stats" :key="stat.type" @click="activeFilter = stat.type">
            <div class="stat-icon" :class="stat.color">{{ stat.icon }}</div>
            <div class="stat-info">
              <span class="stat-value">{{ stat.count }}</span>
              <span class="stat-label">{{ stat.label }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter Pills -->
      <div class="filter-section">
        <div class="filter-pills">
          <button 
            v-for="filter in filters" 
            :key="filter.id"
            @click="activeFilter = filter.id"
            class="filter-pill"
            :class="{ 'active': activeFilter === filter.id }"
          >
            <span class="filter-icon">{{ filter.icon }}</span>
            <span class="filter-label">{{ filter.label }}</span>
            <span v-if="getFilterCount(filter.id) > 0" class="filter-count">{{ getFilterCount(filter.id) }}</span>
          </button>
        </div>
      </div>

      <!-- Main Content -->
      <div class="content-section">
        <!-- Loading State -->
        <div v-if="loading && notifications.length === 0" class="loading-grid">
          <div v-for="i in 4" :key="i" class="skeleton-card">
            <div class="skeleton-icon"></div>
            <div class="skeleton-content">
              <div class="skeleton-line long"></div>
              <div class="skeleton-line medium"></div>
              <div class="skeleton-line short"></div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredNotifications.length === 0" class="empty-state">
          <div class="empty-illustration">
            <div class="empty-icon">
              <span>{{ activeFilter === 'all' ? '🔔' : getActiveFilterIcon() }}</span>
            </div>
            <div class="empty-rings">
              <div class="ring ring-1"></div>
              <div class="ring ring-2"></div>
              <div class="ring ring-3"></div>
            </div>
          </div>
          <h3>{{ activeFilter === 'all' ? 'Aucune notification' : 'Rien ici' }}</h3>
          <p>{{ activeFilter === 'all' ? 'Vos notifications apparaîtront ici' : `Pas de notifications de type "${getActiveFilterLabel()}"` }}</p>
        </div>

        <!-- Notifications List -->
        <div v-else class="notifications-list">
          <!-- Today section -->
          <div v-if="todayNotifications.length > 0" class="notifications-group">
            <div class="group-header">
              <span class="group-title">Aujourd'hui</span>
              <span class="group-count">{{ todayNotifications.length }}</span>
            </div>
            <TransitionGroup name="notif">
              <div 
                v-for="notification in todayNotifications" 
                :key="notification.id"
                @click="handleNotificationClick(notification)"
                class="notification-item"
                :class="{ 'unread': !notification.is_read }"
              >
                <div class="notif-indicator" :class="{ 'show': !notification.is_read }"></div>
                <div class="notif-icon" :class="getIconColorClass(notification.type)">
                  <span>{{ getTypeIcon(notification.type) }}</span>
                </div>
                <div class="notif-content">
                  <div class="notif-header">
                    <h4>{{ notification.title }}</h4>
                    <span class="notif-time">{{ formatTime(notification.created_at) }}</span>
                  </div>
                  <p class="notif-message">{{ notification.message }}</p>
                  <div v-if="notification.data?.amount" class="notif-badge">
                    <span :class="notification.type === 'transfer' && notification.title.includes('reçu') ? 'positive' : 'neutral'">
                      {{ notification.type === 'transfer' && notification.title.includes('reçu') ? '+' : '' }}{{ formatMoney(notification.data.amount, notification.data.currency) }}
                    </span>
                  </div>
                </div>
                <div class="notif-actions">
                  <button 
                    v-if="!notification.is_read"
                    @click.stop="markAsRead(notification.id)"
                    class="notif-action-btn"
                    title="Marquer comme lu"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                    </svg>
                  </button>
                  <button 
                    @click.stop="deleteNotification(notification.id)"
                    class="notif-action-btn delete"
                    title="Supprimer"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </div>
            </TransitionGroup>
          </div>

          <!-- Earlier section -->
          <div v-if="earlierNotifications.length > 0" class="notifications-group">
            <div class="group-header">
              <span class="group-title">Plus tôt</span>
              <span class="group-count">{{ earlierNotifications.length }}</span>
            </div>
            <TransitionGroup name="notif">
              <div 
                v-for="notification in earlierNotifications" 
                :key="notification.id"
                @click="handleNotificationClick(notification)"
                class="notification-item"
                :class="{ 'unread': !notification.is_read }"
              >
                <div class="notif-indicator" :class="{ 'show': !notification.is_read }"></div>
                <div class="notif-icon" :class="getIconColorClass(notification.type)">
                  <span>{{ getTypeIcon(notification.type) }}</span>
                </div>
                <div class="notif-content">
                  <div class="notif-header">
                    <h4>{{ notification.title }}</h4>
                    <span class="notif-time">{{ formatTime(notification.created_at) }}</span>
                  </div>
                  <p class="notif-message">{{ notification.message }}</p>
                  <div v-if="notification.data?.amount" class="notif-badge">
                    <span :class="notification.type === 'transfer' && notification.title.includes('reçu') ? 'positive' : 'neutral'">
                      {{ notification.type === 'transfer' && notification.title.includes('reçu') ? '+' : '' }}{{ formatMoney(notification.data.amount, notification.data.currency) }}
                    </span>
                  </div>
                </div>
                <div class="notif-actions">
                  <button 
                    v-if="!notification.is_read"
                    @click.stop="markAsRead(notification.id)"
                    class="notif-action-btn"
                    title="Marquer comme lu"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                    </svg>
                  </button>
                  <button 
                    @click.stop="deleteNotification(notification.id)"
                    class="notif-action-btn delete"
                    title="Supprimer"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </div>
            </TransitionGroup>
          </div>
        </div>

        <!-- Load More -->
        <div v-if="hasMore && !loading && notifications.length > 0" class="load-more">
          <button @click="loadMore" class="load-more-btn">
            <span>Charger plus</span>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const { notificationApi } = useApi()

const notifications = ref([])
const unreadCount = ref(0)
const loading = ref(false)
const offset = ref(0)
const limit = 20
const hasMore = ref(true)
const activeFilter = ref('all')

const filters = [
  { id: 'all', label: 'Toutes', icon: '📋' },
  { id: 'transfer', label: 'Transferts', icon: '💸' },
  { id: 'security', label: 'Sécurité', icon: '🔐' },
  { id: 'card', label: 'Cartes', icon: '💳' },
  { id: 'wallet', label: 'Portefeuille', icon: '👛' },
]

const stats = computed(() => [
  { type: 'all', icon: '📋', label: 'Total', count: notifications.value.length, color: 'blue' },
  { type: 'transfer', icon: '💸', label: 'Transferts', count: notifications.value.filter(n => n.type === 'transfer').length, color: 'green' },
  { type: 'security', icon: '🔐', label: 'Sécurité', count: notifications.value.filter(n => n.type === 'security').length, color: 'orange' },
  { type: 'card', icon: '💳', label: 'Cartes', count: notifications.value.filter(n => n.type === 'card').length, color: 'purple' },
])

const filteredNotifications = computed(() => {
  if (activeFilter.value === 'all') return notifications.value
  return notifications.value.filter(n => n.type === activeFilter.value)
})

const todayNotifications = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return filteredNotifications.value.filter(n => new Date(n.created_at) >= today)
})

const earlierNotifications = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return filteredNotifications.value.filter(n => new Date(n.created_at) < today)
})

const getFilterCount = (filterId) => {
  if (filterId === 'all') return notifications.value.length
  return notifications.value.filter(n => n.type === filterId).length
}

const getActiveFilterIcon = () => {
  const filter = filters.find(f => f.id === activeFilter.value)
  return filter?.icon || '📋'
}

const getActiveFilterLabel = () => {
  const filter = filters.find(f => f.id === activeFilter.value)
  return filter?.label || ''
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const [notifRes, countRes] = await Promise.all([
      notificationApi.getAll({ limit, offset: 0 }),
      notificationApi.getUnreadCount()
    ])
    notifications.value = notifRes.data.notifications || []
    unreadCount.value = countRes.data.unread_count || 0
    offset.value = limit
    hasMore.value = notifications.value.length >= limit
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    // Demo data
    notifications.value = [
      { id: '1', type: 'transfer', title: '💸 Argent envoyé', message: 'Vous avez envoyé 5,000 XOF avec succès.', is_read: false, created_at: new Date(), data: { amount: 5000, currency: 'XOF' } },
      { id: '2', type: 'transfer', title: '💰 Argent reçu', message: 'Vous avez reçu 10,000 XOF de John Doe.', is_read: false, created_at: new Date(Date.now() - 3600000), data: { amount: 10000, currency: 'XOF' } },
      { id: '3', type: 'security', title: '🔐 Nouvelle connexion', message: 'Connexion depuis Chrome sur Windows', is_read: true, created_at: new Date(Date.now() - 86400000) },
      { id: '4', type: 'card', title: '💳 Paiement carte', message: 'Paiement de 2,500 XOF chez Carrefour', is_read: true, created_at: new Date(Date.now() - 172800000), data: { amount: 2500, currency: 'XOF' } },
    ]
    unreadCount.value = 2
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  loading.value = true
  try {
    const res = await notificationApi.getAll({ limit, offset: offset.value })
    const newNotifs = res.data.notifications || []
    notifications.value.push(...newNotifs)
    offset.value += limit
    hasMore.value = newNotifs.length >= limit
  } catch (error) {
    console.error('Failed to load more:', error)
  } finally {
    loading.value = false
  }
}

const markAsRead = async (id) => {
  try {
    await notificationApi.markAsRead(id)
    const notif = notifications.value.find(n => n.id === id)
    if (notif) {
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    unreadCount.value = 0
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

const deleteNotification = async (id) => {
  try {
    await notificationApi.delete(id)
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      if (!notifications.value[index].is_read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      notifications.value.splice(index, 1)
    }
  } catch (error) {
    console.error('Failed to delete notification:', error)
  }
}

const handleNotificationClick = (notification) => {
  if (!notification.is_read) {
    markAsRead(notification.id)
  }
  if (notification.type === 'transfer') {
    navigateTo('/transactions')
  } else if (notification.type === 'card') {
    navigateTo('/cards')
  } else if (notification.type === 'security') {
    navigateTo('/settings/security')
  }
}

const getTypeIcon = (type) => {
  const icons = {
    transfer: '💸',
    transaction: '💳',
    security: '🔐',
    card: '💳',
    wallet: '👛',
    kyc: '✅',
  }
  return icons[type] || '🔔'
}

const getIconColorClass = (type) => {
  const classes = {
    transfer: 'icon-green',
    transaction: 'icon-blue',
    security: 'icon-orange',
    card: 'icon-purple',
    wallet: 'icon-indigo',
    kyc: 'icon-teal',
  }
  return classes[type] || 'icon-gray'
}

const formatTime = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  
  if (diff < 60000) return 'À l\'instant'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}j`
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })
}

const formatMoney = (amount, currency = 'XOF') => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

onMounted(() => {
  fetchNotifications()
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.notifications-page {
  min-height: 100vh;
  padding-bottom: 3rem;
  max-width: 100%;
  overflow-x: hidden;
}

/* Header Section */
.header-section {
  position: relative;
  overflow: hidden;
  margin-bottom: 2rem;
}

.header-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 50%, #4338ca 100%);
}

.header-bg::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 30% 20%, rgba(255,255,255,0.15) 0%, transparent 50%);
}

.header-content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem 1rem;
}

.header-icon-wrapper {
  flex-shrink: 0;
}

.header-icon {
  position: relative;
  width: 3.5rem;
  height: 3.5rem;
  border-radius: 1rem;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.header-icon svg {
  width: 1.5rem;
  height: 1.5rem;
  color: white;
}

.header-badge {
  position: absolute;
  top: -0.5rem;
  right: -0.5rem;
  min-width: 20px;
  height: 20px;
  padding: 0 0.375rem;
  border-radius: 9999px;
  background: #ef4444;
  color: white;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(239,68,68,0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.header-text {
  flex: 1;
  min-width: 0;
}

.header-text h1 {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.header-text p {
  font-size: 0.875rem;
  color: rgba(255,255,255,0.8);
}

.header-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  border-radius: 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 0.2s;
  border: none;
  cursor: pointer;
}

.action-btn svg {
  width: 1rem;
  height: 1rem;
}

.action-btn.primary {
  background: white;
  color: #6366f1;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}

.action-btn.primary:hover {
  background: rgba(255,255,255,0.9);
}

.action-btn.secondary {
  background: rgba(255,255,255,0.2);
  color: white;
  backdrop-filter: blur(8px);
}

.action-btn.secondary:hover {
  background: rgba(255,255,255,0.3);
}

.action-btn.loading {
  opacity: 0.75;
  cursor: not-allowed;
}

/* Stats Section */
.stats-section {
  padding: 0 1rem;
  margin-top: -1rem;
  position: relative;
  z-index: 20;
  margin-bottom: 1.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem;
  background: var(--color-surface, #1a1a2e);
  border-radius: 1rem;
  border: 1px solid rgba(255,255,255,0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}

.stat-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.stat-icon.blue { background: rgba(59, 130, 246, 0.2); }
.stat-icon.green { background: rgba(34, 197, 94, 0.2); }
.stat-icon.orange { background: rgba(249, 115, 22, 0.2); }
.stat-icon.purple { background: rgba(168, 85, 247, 0.2); }

.stat-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text, #fff);
}

.stat-label {
  font-size: 0.75rem;
  color: var(--color-muted, #888);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Filter Section */
.filter-section {
  padding: 0 1rem;
  margin-bottom: 1rem;
}

.filter-pills {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
  padding-bottom: 0.5rem;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.filter-pills::-webkit-scrollbar {
  display: none;
}

.filter-pill {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
  background: var(--color-surface, #1a1a2e);
  border: 1px solid rgba(255,255,255,0.1);
  color: var(--color-muted, #888);
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.filter-pill:hover {
  border-color: rgba(99, 102, 241, 0.5);
  color: var(--color-text, #fff);
}

.filter-pill.active {
  background: #6366f1;
  color: white;
  border-color: #6366f1;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.filter-icon {
  font-size: 0.875rem;
}

.filter-count {
  margin-left: 0.25rem;
  padding: 0.125rem 0.375rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 700;
  background: rgba(255,255,255,0.1);
  color: inherit;
}

.filter-pill.active .filter-count {
  background: rgba(255,255,255,0.2);
}

/* Content Section */
.content-section {
  padding: 0 1rem;
}

/* Loading Grid */
.loading-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.skeleton-card {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--color-surface, #1a1a2e);
  border-radius: 1rem;
  border: 1px solid rgba(255,255,255,0.1);
}

.skeleton-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  background: rgba(255,255,255,0.1);
  flex-shrink: 0;
  animation: shimmer 1.5s infinite;
}

.skeleton-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.skeleton-line {
  height: 0.75rem;
  border-radius: 9999px;
  background: rgba(255,255,255,0.1);
  animation: shimmer 1.5s infinite;
}

.skeleton-line.long { width: 75%; }
.skeleton-line.medium { width: 50%; }
.skeleton-line.short { width: 25%; }

@keyframes shimmer {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-illustration {
  position: relative;
  width: 6rem;
  height: 6rem;
  margin: 0 auto 1.5rem;
}

.empty-icon {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  z-index: 10;
}

.empty-rings {
  position: absolute;
  inset: 0;
}

.ring {
  position: absolute;
  border-radius: 50%;
  border: 2px solid rgba(99, 102, 241, 0.2);
  animation: ripple 3s ease-out infinite;
}

.ring-1 { inset: 0; animation-delay: 0s; }
.ring-2 { inset: 0.5rem; animation-delay: 1s; }
.ring-3 { inset: 1rem; animation-delay: 2s; }

@keyframes ripple {
  0% { transform: scale(1); opacity: 0.5; }
  100% { transform: scale(1.5); opacity: 0; }
}

.empty-state h3 {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--color-text, #fff);
  margin-bottom: 0.5rem;
}

.empty-state p {
  font-size: 0.875rem;
  color: var(--color-muted, #888);
}

/* Notifications List */
.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.notifications-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0 0.5rem;
}

.group-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-muted, #888);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.group-count {
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 700;
  background: rgba(255,255,255,0.1);
  color: var(--color-muted, #888);
}

/* Notification Item */
.notification-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--color-surface, #1a1a2e);
  border-radius: 1rem;
  border: 1px solid rgba(255,255,255,0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.notification-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
  border-color: rgba(99, 102, 241, 0.3);
}

.notification-item.unread {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
}

.notif-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background: #6366f1;
  border-radius: 0 2px 2px 0;
  transition: height 0.2s;
}

.notif-indicator.show {
  height: 2.5rem;
}

.notif-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
  flex-shrink: 0;
}

.icon-green { background: rgba(34, 197, 94, 0.2); }
.icon-blue { background: rgba(59, 130, 246, 0.2); }
.icon-orange { background: rgba(249, 115, 22, 0.2); }
.icon-purple { background: rgba(168, 85, 247, 0.2); }
.icon-indigo { background: rgba(99, 102, 241, 0.2); }
.icon-teal { background: rgba(20, 184, 166, 0.2); }
.icon-gray { background: rgba(107, 114, 128, 0.2); }

.notif-content {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.notif-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.notif-header h4 {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text, #fff);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notif-time {
  font-size: 0.625rem;
  color: var(--color-muted, #666);
  white-space: nowrap;
  flex-shrink: 0;
}

.notif-message {
  font-size: 0.75rem;
  color: var(--color-muted, #888);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.notif-badge {
  margin-top: 0.5rem;
}

.notif-badge span {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  font-size: 0.625rem;
  font-weight: 700;
}

.notif-badge .positive {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.notif-badge .neutral {
  background: rgba(107, 114, 128, 0.2);
  color: #888;
}

.notif-actions {
  display: flex;
  gap: 0.25rem;
  opacity: 0;
  transition: opacity 0.2s;
  flex-shrink: 0;
}

.notification-item:hover .notif-actions {
  opacity: 1;
}

/* Always show actions on touch devices */
@media (hover: none) {
  .notif-actions {
    opacity: 1;
  }
}

.notif-action-btn {
  padding: 0.375rem;
  border-radius: 0.5rem;
  background: rgba(255,255,255,0.1);
  color: var(--color-muted, #888);
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.notif-action-btn svg {
  width: 0.875rem;
  height: 0.875rem;
}

.notif-action-btn:hover {
  background: rgba(255,255,255,0.2);
  color: var(--color-text, #fff);
}

.notif-action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

/* Load More */
.load-more {
  text-align: center;
  margin-top: 2rem;
}

.load-more-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  background: var(--color-surface, #1a1a2e);
  border: 1px solid rgba(255,255,255,0.1);
  color: var(--color-muted, #888);
  cursor: pointer;
  transition: all 0.2s;
}

.load-more-btn svg {
  width: 1rem;
  height: 1rem;
}

.load-more-btn:hover {
  border-color: rgba(99, 102, 241, 0.5);
  color: var(--color-text, #fff);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

/* Transitions */
.notif-enter-active,
.notif-leave-active {
  transition: all 0.3s ease;
}

.notif-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.notif-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* Desktop styles */
@media (min-width: 768px) {
  .header-content {
    padding: 2rem;
    gap: 1.5rem;
  }
  
  .header-icon {
    width: 4rem;
    height: 4rem;
  }
  
  .header-icon svg {
    width: 2rem;
    height: 2rem;
  }
  
  .header-text h1 {
    font-size: 1.75rem;
  }
  
  .action-btn {
    padding: 0.625rem 1rem;
    font-size: 0.875rem;
  }
  
  .stats-section {
    padding: 0 2rem;
    margin-top: -1.5rem;
  }
  
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
  }
  
  .stat-card {
    padding: 1rem;
  }
  
  .stat-icon {
    width: 3rem;
    height: 3rem;
    font-size: 1.5rem;
  }
  
  .stat-value {
    font-size: 1.5rem;
  }
  
  .filter-section {
    padding: 0 2rem;
  }
  
  .content-section {
    padding: 0 2rem;
  }
  
  .notification-item {
    padding: 1.25rem;
    gap: 1rem;
  }
  
  .notif-icon {
    width: 3rem;
    height: 3rem;
    font-size: 1.25rem;
  }
  
  .notif-header h4 {
    font-size: 1rem;
  }
  
  .notif-message {
    font-size: 0.875rem;
  }
}
</style>
