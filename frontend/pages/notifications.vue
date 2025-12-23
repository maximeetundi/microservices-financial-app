<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-3xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-2xl font-bold text-base">Notifications 🔔</h1>
          <p class="text-muted mt-1">
            {{ unreadCount > 0 ? `${unreadCount} non lue(s)` : 'Toutes lues' }}
          </p>
        </div>
        <div class="flex gap-3">
          <button 
            v-if="unreadCount > 0"
            @click="markAllAsRead"
            class="px-4 py-2 rounded-xl text-sm font-medium bg-primary/10 text-primary hover:bg-primary/20 transition-colors"
          >
            Tout marquer comme lu
          </button>
          <button 
            @click="fetchNotifications"
            :disabled="loading"
            class="p-2 rounded-xl bg-surface-hover hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-colors"
          >
            <svg class="w-5 h-5 text-muted" :class="{'animate-spin': loading}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="flex gap-2 mb-6 overflow-x-auto pb-2">
        <button 
          v-for="filter in filters" 
          :key="filter.id"
          @click="activeFilter = filter.id"
          class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-all"
          :class="activeFilter === filter.id 
            ? 'bg-primary text-white shadow-lg shadow-primary/25' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          {{ filter.icon }} {{ filter.label }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading && notifications.length === 0" class="space-y-4">
        <div v-for="i in 5" :key="i" class="notification-skeleton">
          <div class="skeleton-icon"></div>
          <div class="skeleton-content">
            <div class="skeleton-line w-2/3"></div>
            <div class="skeleton-line w-1/2"></div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredNotifications.length === 0" class="text-center py-16">
        <div class="w-20 h-20 mx-auto mb-4 rounded-2xl bg-secondary-100 dark:bg-secondary-800 flex items-center justify-center">
          <span class="text-4xl">📭</span>
        </div>
        <h3 class="text-lg font-medium text-base mb-2">Aucune notification</h3>
        <p class="text-muted text-sm">
          {{ activeFilter === 'all' ? 'Vous n\'avez pas encore de notifications' : 'Aucune notification dans cette catégorie' }}
        </p>
      </div>

      <!-- Notifications List -->
      <div v-else class="space-y-3">
        <TransitionGroup name="notification-list">
          <div 
            v-for="notification in filteredNotifications" 
            :key="notification.id"
            @click="handleNotificationClick(notification)"
            class="notification-card group"
            :class="{ 'unread': !notification.is_read }"
          >
            <!-- Type Icon -->
            <div class="notification-icon" :class="getIconClass(notification.type)">
              <span>{{ getTypeIcon(notification.type) }}</span>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1 min-w-0">
                  <h4 class="font-semibold text-base truncate">
                    {{ notification.title }}
                  </h4>
                  <p class="text-sm text-muted mt-1 line-clamp-2">
                    {{ notification.message }}
                  </p>
                </div>
                <span class="text-xs text-muted whitespace-nowrap">
                  {{ formatTime(notification.created_at) }}
                </span>
              </div>
              
              <!-- Data Badge -->
              <div v-if="notification.data?.amount" class="mt-2">
                <span class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-medium bg-success/10 text-success">
                  {{ formatMoney(notification.data.amount, notification.data.currency) }}
                </span>
              </div>
            </div>

            <!-- Unread Indicator -->
            <div v-if="!notification.is_read" class="unread-dot"></div>

            <!-- Actions -->
            <div class="notification-actions">
              <button 
                v-if="!notification.is_read"
                @click.stop="markAsRead(notification.id)"
                class="action-btn"
                title="Marquer comme lu"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
              </button>
              <button 
                @click.stop="deleteNotification(notification.id)"
                class="action-btn text-error"
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

      <!-- Load More -->
      <div v-if="hasMore && !loading" class="text-center mt-8">
        <button 
          @click="loadMore"
          class="px-6 py-3 rounded-xl font-medium bg-surface-hover hover:bg-secondary-200 dark:hover:bg-secondary-700 text-muted hover:text-base transition-all"
        >
          Charger plus
        </button>
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
  { id: 'transaction', label: 'Transactions', icon: '💳' },
  { id: 'security', label: 'Sécurité', icon: '🔐' },
  { id: 'card', label: 'Cartes', icon: '💳' },
]

const filteredNotifications = computed(() => {
  if (activeFilter.value === 'all') return notifications.value
  return notifications.value.filter(n => n.type === activeFilter.value)
})

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
    // Mock data for demo
    notifications.value = [
      { id: '1', type: 'transfer', title: '💸 Argent envoyé', message: 'Vous avez envoyé 5,000 XOF avec succès.', is_read: false, created_at: new Date(), data: { amount: 5000, currency: 'XOF' } },
      { id: '2', type: 'transfer', title: '💰 Argent reçu', message: 'Vous avez reçu 10,000 XOF de John Doe.', is_read: false, created_at: new Date(Date.now() - 3600000), data: { amount: 10000, currency: 'XOF' } },
      { id: '3', type: 'security', title: '🔐 Nouvelle connexion', message: 'Connexion depuis Chrome sur Windows', is_read: true, created_at: new Date(Date.now() - 86400000) },
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
  // Navigate based on type
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

const getIconClass = (type) => {
  const classes = {
    transfer: 'bg-emerald-100 dark:bg-emerald-900/30',
    transaction: 'bg-blue-100 dark:bg-blue-900/30',
    security: 'bg-amber-100 dark:bg-amber-900/30',
    card: 'bg-purple-100 dark:bg-purple-900/30',
    wallet: 'bg-indigo-100 dark:bg-indigo-900/30',
    kyc: 'bg-green-100 dark:bg-green-900/30',
  }
  return classes[type] || 'bg-gray-100 dark:bg-gray-900/30'
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
.notification-card {
  @apply relative flex items-start gap-4 p-4 rounded-2xl bg-surface border border-secondary-200 dark:border-secondary-700 cursor-pointer transition-all;
}

.notification-card:hover {
  @apply bg-surface-hover shadow-lg;
}

.notification-card.unread {
  @apply bg-primary/5 border-primary/20;
}

.notification-icon {
  @apply w-12 h-12 rounded-xl flex items-center justify-center text-xl flex-shrink-0;
}

.unread-dot {
  @apply absolute top-4 right-4 w-2.5 h-2.5 rounded-full bg-primary shadow-lg shadow-primary/50;
}

.notification-actions {
  @apply absolute top-1/2 right-4 -translate-y-1/2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity;
}

.action-btn {
  @apply p-2 rounded-lg bg-surface hover:bg-secondary-200 dark:hover:bg-secondary-700 text-muted hover:text-base transition-colors;
}

.action-btn.text-error:hover {
  color: rgb(239 68 68);
  background-color: rgb(239 68 68 / 0.1);
}

/* Skeleton */
.notification-skeleton {
  @apply flex items-start gap-4 p-4 rounded-2xl bg-surface animate-pulse;
}

.skeleton-icon {
  @apply w-12 h-12 rounded-xl bg-secondary-200 dark:bg-secondary-700;
}

.skeleton-content {
  @apply flex-1 space-y-2;
}

.skeleton-line {
  @apply h-4 rounded bg-secondary-200 dark:bg-secondary-700;
}

/* List animation */
.notification-list-enter-active,
.notification-list-leave-active {
  transition: all 0.3s ease;
}

.notification-list-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.notification-list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
