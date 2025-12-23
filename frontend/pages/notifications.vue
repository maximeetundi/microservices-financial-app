<template>
  <NuxtLayout name="dashboard">
    <div class="notif-page">
      <!-- Simple Header -->
      <div class="page-header">
        <h1>🔔 Notifications</h1>
        <div class="header-actions">
          <span v-if="unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
          <button v-if="unreadCount > 0" @click="markAllAsRead" class="btn-text">
            Tout lire
          </button>
          <button @click="fetchNotifications" :disabled="loading" class="btn-icon">
            <svg :class="{ 'spin': loading }" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Filter Pills -->
      <div class="filters">
        <button 
          v-for="filter in filters" 
          :key="filter.id"
          @click="activeFilter = filter.id"
          :class="['filter-btn', { active: activeFilter === filter.id }]"
        >
          {{ filter.icon }} {{ filter.label }}
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading && notifications.length === 0" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement...</p>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredNotifications.length === 0" class="empty-state">
        <span class="empty-icon">🔔</span>
        <p>Aucune notification</p>
      </div>

      <!-- List -->
      <div v-else class="notif-list">
        <div 
          v-for="notif in filteredNotifications" 
          :key="notif.id"
          @click="handleNotificationClick(notif)"
          :class="['notif-item', { unread: !notif.is_read }]"
        >
          <div class="notif-icon">{{ getTypeIcon(notif.type) }}</div>
          <div class="notif-body">
            <div class="notif-title">{{ notif.title }}</div>
            <div class="notif-msg">{{ notif.message }}</div>
            <div class="notif-time">{{ formatTime(notif.created_at) }}</div>
          </div>
          <div class="notif-actions">
            <button v-if="!notif.is_read" @click.stop="markAsRead(notif.id)" class="action-mark">✓</button>
            <button @click.stop="deleteNotification(notif.id)" class="action-delete">✕</button>
          </div>
        </div>
      </div>

      <!-- Load More -->
      <div v-if="hasMore && !loading && notifications.length > 0" class="load-more">
        <button @click="loadMore">Charger plus ↓</button>
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
    notifications.value = []
    unreadCount.value = 0
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
}

const getTypeIcon = (type) => {
  const icons = { transfer: '💸', card: '💳', security: '🔐', wallet: '👛', kyc: '✅' }
  return icons[type] || '🔔'
}

const formatTime = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return 'À l\'instant'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h`
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })
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
.notif-page {
  width: 100%;
  max-width: 100%;
  padding: 0;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  margin: -1rem -1rem 1rem -1rem;
  width: calc(100% + 2rem);
}

.page-header h1 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.unread-badge {
  background: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.125rem 0.5rem;
  border-radius: 999px;
}

.btn-text {
  background: rgba(255,255,255,0.2);
  color: white;
  border: none;
  padding: 0.375rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  cursor: pointer;
}

.btn-icon {
  background: rgba(255,255,255,0.2);
  color: white;
  border: none;
  width: 2rem;
  height: 2rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.filters {
  display: flex;
  gap: 0.5rem;
  padding: 0 0 1rem 0;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.filters::-webkit-scrollbar {
  display: none;
}

.filter-btn {
  flex-shrink: 0;
  padding: 0.5rem 0.875rem;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.05);
  color: #888;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn.active {
  background: #6366f1;
  color: white;
  border-color: #6366f1;
}

.loading-state, .empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #888;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.notif-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.notif-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.875rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 0.75rem;
  cursor: pointer;
  transition: background 0.2s;
}

.notif-item:active {
  background: rgba(255,255,255,0.08);
}

.notif-item.unread {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
}

.notif-icon {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.625rem;
  background: rgba(99, 102, 241, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
}

.notif-body {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.notif-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notif-msg {
  font-size: 0.75rem;
  color: #888;
  margin-top: 0.125rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notif-time {
  font-size: 0.625rem;
  color: #666;
  margin-top: 0.25rem;
}

.notif-actions {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex-shrink: 0;
}

.action-mark, .action-delete {
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 0.375rem;
  border: none;
  font-size: 0.75rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-mark {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.action-delete {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.load-more {
  text-align: center;
  padding: 1rem;
}

.load-more button {
  padding: 0.75rem 1.5rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: transparent;
  color: #888;
  font-size: 0.875rem;
  cursor: pointer;
}

/* Desktop enhancements */
@media (min-width: 640px) {
  .page-header {
    padding: 1.5rem;
    border-radius: 1rem;
    margin: 0 0 1.5rem 0;
    width: 100%;
  }
  
  .page-header h1 {
    font-size: 1.5rem;
  }
  
  .notif-item {
    padding: 1rem;
  }
  
  .notif-icon {
    width: 2.75rem;
    height: 2.75rem;
    font-size: 1.25rem;
  }
  
  .notif-title {
    font-size: 1rem;
  }
  
  .notif-msg {
    font-size: 0.875rem;
  }
}
</style>
