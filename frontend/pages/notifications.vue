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
  @apply min-h-screen pb-12;
}

/* Header Section */
.header-section {
  @apply relative overflow-hidden mb-8;
}

.header-bg {
  @apply absolute inset-0 bg-gradient-to-br from-primary via-primary-600 to-indigo-700;
}

.header-bg::before {
  content: '';
  @apply absolute inset-0;
  background: radial-gradient(circle at 30% 20%, rgba(255,255,255,0.15) 0%, transparent 50%);
}

.header-content {
  @apply relative z-10 flex flex-wrap items-center gap-6 p-8;
}

.header-icon-wrapper {
  @apply flex-shrink-0;
}

.header-icon {
  @apply relative w-16 h-16 rounded-2xl bg-white/20 backdrop-blur flex items-center justify-center shadow-lg;
}

.header-badge {
  @apply absolute -top-2 -right-2 min-w-[24px] h-6 px-1.5 rounded-full bg-red-500 text-white text-xs font-bold flex items-center justify-center shadow-lg;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.header-text {
  @apply flex-1;
}

.header-text h1 {
  @apply text-2xl md:text-3xl font-bold text-white mb-1;
}

.header-text p {
  @apply text-white/80;
}

.header-actions {
  @apply flex gap-3;
}

.action-btn {
  @apply flex items-center gap-2 px-4 py-2.5 rounded-xl font-medium transition-all duration-200;
}

.action-btn.primary {
  @apply bg-white text-primary hover:bg-white/90 shadow-lg;
}

.action-btn.secondary {
  @apply bg-white/20 text-white hover:bg-white/30 backdrop-blur;
}

.action-btn.loading {
  @apply opacity-75 cursor-not-allowed;
}

/* Stats Section */
.stats-section {
  @apply px-4 md:px-8 -mt-6 relative z-20 mb-8;
}

.stats-grid {
  @apply grid grid-cols-2 md:grid-cols-4 gap-4;
}

.stat-card {
  @apply flex items-center gap-4 p-4 bg-surface rounded-2xl shadow-lg border border-secondary-100 dark:border-secondary-800 cursor-pointer hover:shadow-xl transition-all duration-200 hover:-translate-y-0.5;
}

.stat-icon {
  @apply w-12 h-12 rounded-xl flex items-center justify-center text-2xl;
}

.stat-icon.blue { @apply bg-blue-100 dark:bg-blue-900/30; }
.stat-icon.green { @apply bg-emerald-100 dark:bg-emerald-900/30; }
.stat-icon.orange { @apply bg-orange-100 dark:bg-orange-900/30; }
.stat-icon.purple { @apply bg-purple-100 dark:bg-purple-900/30; }

.stat-info {
  @apply flex flex-col;
}

.stat-value {
  @apply text-2xl font-bold text-base;
}

.stat-label {
  @apply text-sm text-muted;
}

/* Filter Section */
.filter-section {
  @apply px-4 md:px-8 mb-6;
}

.filter-pills {
  @apply flex gap-2 overflow-x-auto pb-2 scrollbar-hide;
}

.filter-pill {
  @apply flex items-center gap-2 px-4 py-2.5 rounded-full font-medium text-sm whitespace-nowrap transition-all duration-200;
  @apply bg-surface border border-secondary-200 dark:border-secondary-700 text-muted hover:text-base hover:border-primary/50;
}

.filter-pill.active {
  @apply bg-primary text-white border-primary shadow-lg shadow-primary/25;
}

.filter-icon {
  @apply text-base;
}

.filter-count {
  @apply ml-1 px-1.5 py-0.5 rounded-full text-xs font-bold;
  @apply bg-secondary-200 dark:bg-secondary-700 text-muted;
}

.filter-pill.active .filter-count {
  @apply bg-white/20 text-white;
}

/* Content Section */
.content-section {
  @apply px-4 md:px-8;
}

/* Loading Grid */
.loading-grid {
  @apply space-y-4;
}

.skeleton-card {
  @apply flex items-start gap-4 p-5 bg-surface rounded-2xl border border-secondary-100 dark:border-secondary-800 animate-pulse;
}

.skeleton-icon {
  @apply w-12 h-12 rounded-xl bg-secondary-200 dark:bg-secondary-700 flex-shrink-0;
}

.skeleton-content {
  @apply flex-1 space-y-3;
}

.skeleton-line {
  @apply h-4 rounded-full bg-secondary-200 dark:bg-secondary-700;
}

.skeleton-line.long { @apply w-3/4; }
.skeleton-line.medium { @apply w-1/2; }
.skeleton-line.short { @apply w-1/4; }

/* Empty State */
.empty-state {
  @apply text-center py-16;
}

.empty-illustration {
  @apply relative w-32 h-32 mx-auto mb-6;
}

.empty-icon {
  @apply absolute inset-0 flex items-center justify-center text-5xl z-10;
}

.empty-rings {
  @apply absolute inset-0;
}

.ring {
  @apply absolute rounded-full border-2 border-primary/20;
  animation: ripple 3s ease-out infinite;
}

.ring-1 { @apply inset-0; animation-delay: 0s; }
.ring-2 { @apply inset-2; animation-delay: 1s; }
.ring-3 { @apply inset-4; animation-delay: 2s; }

@keyframes ripple {
  0% { transform: scale(1); opacity: 0.5; }
  100% { transform: scale(1.5); opacity: 0; }
}

.empty-state h3 {
  @apply text-xl font-bold text-base mb-2;
}

.empty-state p {
  @apply text-muted;
}

/* Notifications List */
.notifications-list {
  @apply space-y-6;
}

.notifications-group {
  @apply space-y-3;
}

.group-header {
  @apply flex items-center gap-3 px-2;
}

.group-title {
  @apply text-sm font-semibold text-muted uppercase tracking-wide;
}

.group-count {
  @apply px-2 py-0.5 rounded-full text-xs font-bold bg-secondary-200 dark:bg-secondary-700 text-muted;
}

/* Notification Item */
.notification-item {
  @apply relative flex items-start gap-4 p-5 bg-surface rounded-2xl border border-secondary-100 dark:border-secondary-800 cursor-pointer transition-all duration-200;
  @apply hover:shadow-lg hover:-translate-y-0.5 hover:border-primary/30;
}

.notification-item.unread {
  @apply bg-primary/5 border-primary/20;
}

.notif-indicator {
  @apply absolute left-0 top-1/2 -translate-y-1/2 w-1 h-0 bg-primary rounded-r-full transition-all duration-200;
}

.notif-indicator.show {
  @apply h-12;
}

.notif-icon {
  @apply w-12 h-12 rounded-xl flex items-center justify-center text-xl flex-shrink-0;
}

.icon-green { @apply bg-emerald-100 dark:bg-emerald-900/30; }
.icon-blue { @apply bg-blue-100 dark:bg-blue-900/30; }
.icon-orange { @apply bg-orange-100 dark:bg-orange-900/30; }
.icon-purple { @apply bg-purple-100 dark:bg-purple-900/30; }
.icon-indigo { @apply bg-indigo-100 dark:bg-indigo-900/30; }
.icon-teal { @apply bg-teal-100 dark:bg-teal-900/30; }
.icon-gray { @apply bg-secondary-100 dark:bg-secondary-800; }

.notif-content {
  @apply flex-1 min-w-0;
}

.notif-header {
  @apply flex items-start justify-between gap-4 mb-1;
}

.notif-header h4 {
  @apply font-semibold text-base truncate;
}

.notif-time {
  @apply text-xs text-muted whitespace-nowrap;
}

.notif-message {
  @apply text-sm text-muted line-clamp-2;
}

.notif-badge {
  @apply mt-2;
}

.notif-badge span {
  @apply inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-bold;
}

.notif-badge .positive {
  @apply bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400;
}

.notif-badge .neutral {
  @apply bg-secondary-100 dark:bg-secondary-800 text-muted;
}

.notif-actions {
  @apply flex gap-1 opacity-0 transition-opacity duration-200;
}

.notification-item:hover .notif-actions {
  @apply opacity-100;
}

.notif-action-btn {
  @apply p-2 rounded-lg transition-colors duration-200;
  @apply bg-secondary-100 dark:bg-secondary-800 text-muted hover:text-base hover:bg-secondary-200 dark:hover:bg-secondary-700;
}

.notif-action-btn.delete:hover {
  @apply bg-red-100 dark:bg-red-900/30 text-red-600;
}

/* Load More */
.load-more {
  @apply text-center mt-8;
}

.load-more-btn {
  @apply inline-flex items-center gap-2 px-6 py-3 rounded-xl font-medium transition-all duration-200;
  @apply bg-surface border border-secondary-200 dark:border-secondary-700 text-muted hover:text-base hover:border-primary/50 hover:shadow-lg;
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

/* Utility */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
