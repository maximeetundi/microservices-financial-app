<template>
  <div class="notification-center">
    <!-- Bell Button -->
    <button 
      @click="toggleDropdown" 
      class="relative p-2 rounded-xl hover:bg-surface-hover transition-colors"
      ref="bellButton"
    >
      <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
      </svg>
      <!-- Unread Badge -->
      <span 
        v-if="unreadCount > 0" 
        class="absolute -top-1 -right-1 min-w-[20px] h-5 flex items-center justify-center bg-red-500 text-white text-xs font-bold rounded-full px-1"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>

    <!-- Backdrop (click outside to close) -->
    <Teleport to="body">
      <div v-if="isOpen" class="notification-backdrop" @click="isOpen = false"></div>
    </Teleport>

    <!-- Dropdown (teleported to body to avoid overflow issues) -->
    <Teleport to="body">
      <Transition name="dropdown">
        <div 
          v-if="isOpen" 
          class="notification-dropdown"
          :style="dropdownStyle"
        >
          <!-- Header -->
          <div class="notification-header">
            <h3 class="font-semibold text-base">Notifications</h3>
            <button 
              v-if="unreadCount > 0" 
              @click="markAllRead" 
              class="text-xs text-primary hover:underline whitespace-nowrap"
            >
              Tout marquer comme lu
            </button>
          </div>

          <!-- Loading -->
          <div v-if="loading" class="p-8 flex justify-center">
            <div class="loading-spinner w-6 h-6"></div>
          </div>

          <!-- Empty State -->
          <div v-else-if="notifications.length === 0" class="p-8 text-center">
            <div class="text-4xl mb-2">🔔</div>
            <p class="text-muted text-sm">Aucune notification</p>
          </div>

          <!-- Notification List -->
          <div v-else class="notification-list">
            <div 
              v-for="notif in notifications" 
              :key="notif.id"
              @click="handleClick(notif)"
              class="notification-item"
              :class="{ 'unread': !notif.is_read }"
            >
              <div class="flex items-start gap-3">
                <!-- Icon -->
                <div 
                  class="notification-icon"
                  :class="getIconClass(notif.type)"
                >
                  {{ getIcon(notif.type) }}
                </div>
                <!-- Content -->
                <div class="notification-content">
                  <p class="notification-title">{{ notif.title }}</p>
                  <p class="notification-message">{{ notif.message }}</p>
                  <p class="notification-time">{{ formatTime(notif.created_at) }}</p>
                </div>
                <!-- Unread dot -->
                <div v-if="!notif.is_read" class="unread-dot"></div>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="notification-footer">
            <NuxtLink 
              to="/notifications" 
              class="block text-center text-sm text-primary hover:underline py-2"
              @click="isOpen = false"
            >
              Voir toutes les notifications
            </NuxtLink>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { notificationAPI } from '~/composables/useApi'

const isOpen = ref(false)
const loading = ref(false)
const notifications = ref([])
const unreadCount = ref(0)
const bellButton = ref(null)
const dropdownPosition = ref({ top: 0, right: 0 })

let pollInterval = null

const dropdownStyle = computed(() => ({
  position: 'fixed',
  top: `${dropdownPosition.value.top}px`,
  right: `${dropdownPosition.value.right}px`,
  zIndex: 9999
}))

const updateDropdownPosition = () => {
  if (bellButton.value) {
    const rect = bellButton.value.getBoundingClientRect()
    dropdownPosition.value = {
      top: rect.bottom + 8,
      right: Math.max(16, window.innerWidth - rect.right)
    }
  }
}

const toggleDropdown = async () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    await nextTick()
    updateDropdownPosition()
    if (notifications.value.length === 0) {
      fetchNotifications()
    }
  }
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await notificationAPI.getAll(10, 0)
    notifications.value = res.data?.notifications || []
  } catch (e) {
    console.error('Failed to fetch notifications:', e)
  } finally {
    loading.value = false
  }
}

const fetchUnreadCount = async () => {
  try {
    const res = await notificationAPI.getUnreadCount()
    unreadCount.value = res.data?.unread_count || 0
  } catch (e) {
    // Silently fail
  }
}

const handleClick = async (notif) => {
  if (!notif.is_read) {
    try {
      await notificationAPI.markAsRead(notif.id)
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (e) {
      console.error('Failed to mark as read:', e)
    }
  }
}

const markAllRead = async () => {
  try {
    await notificationAPI.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    unreadCount.value = 0
  } catch (e) {
    console.error('Failed to mark all as read:', e)
  }
}

const getIcon = (type) => {
  const icons = {
    transfer: '💸',
    card: '💳',
    security: '🔒',
    promotion: '🎁',
    wallet: '💰',
    payment: '💳',
    default: '🔔'
  }
  return icons[type] || icons.default
}

const getIconClass = (type) => {
  const classes = {
    transfer: 'bg-green-500/20 text-green-500',
    card: 'bg-purple-500/20 text-purple-500',
    security: 'bg-red-500/20 text-red-500',
    promotion: 'bg-yellow-500/20 text-yellow-500',
    wallet: 'bg-blue-500/20 text-blue-500',
    payment: 'bg-indigo-500/20 text-indigo-500',
    default: 'bg-gray-500/20 text-gray-500'
  }
  return classes[type] || classes.default
}

const formatTime = (date) => {
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  
  if (diff < 60000) return 'À l\'instant'
  if (diff < 3600000) return `Il y a ${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `Il y a ${Math.floor(diff / 3600000)} h`
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })
}

// Handle window resize
const handleResize = () => {
  if (isOpen.value) {
    updateDropdownPosition()
  }
}

onMounted(() => {
  fetchUnreadCount()
  pollInterval = setInterval(fetchUnreadCount, 30000)
  window.addEventListener('resize', handleResize)
  window.addEventListener('scroll', handleResize, true)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('scroll', handleResize, true)
})
</script>

<style scoped>
.notification-center {
  position: relative;
  display: inline-block;
}

.notification-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9998;
  background: transparent;
}

.notification-dropdown {
  width: min(320px, calc(100vw - 32px));
  max-height: min(400px, calc(100vh - 100px));
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: var(--color-surface, #1a1a2e);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.notification-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.notification-list {
  flex: 1;
  overflow-y: auto;
  overscroll-behavior: contain;
}

.notification-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.notification-item.unread {
  background: rgba(99, 102, 241, 0.1);
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.notification-title {
  font-weight: 500;
  font-size: 14px;
  color: var(--color-text, #fff);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notification-message {
  font-size: 12px;
  color: var(--color-muted, #888);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
  margin-top: 2px;
}

.notification-time {
  font-size: 11px;
  color: var(--color-muted, #666);
  margin-top: 4px;
}

.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-primary, #6366f1);
  flex-shrink: 0;
  margin-top: 6px;
}

.notification-footer {
  padding: 8px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

/* Transitions */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}

/* Mobile adjustments */
@media (max-width: 480px) {
  .notification-dropdown {
    width: calc(100vw - 24px);
    right: 12px !important;
  }
}
</style>

