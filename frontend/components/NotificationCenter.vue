<template>
  <div class="relative">
    <!-- Bell Button -->
    <button 
      @click="toggleDropdown" 
      class="relative p-2 rounded-xl hover:bg-surface-hover transition-colors"
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

    <!-- Dropdown -->
    <Transition name="dropdown">
      <div 
        v-if="isOpen" 
        class="absolute right-0 top-12 w-80 max-h-96 overflow-y-auto glass-card border border-secondary-200 dark:border-secondary-700 shadow-xl z-50 rounded-xl"
      >
        <!-- Header -->
        <div class="flex items-center justify-between p-3 border-b border-secondary-200 dark:border-secondary-700">
          <h3 class="font-semibold text-base">Notifications</h3>
          <button 
            v-if="unreadCount > 0" 
            @click="markAllRead" 
            class="text-xs text-primary hover:underline"
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
        <div v-else class="divide-y divide-secondary-200 dark:divide-secondary-700">
          <div 
            v-for="notif in notifications" 
            :key="notif.id"
            @click="handleClick(notif)"
            class="p-3 hover:bg-surface-hover cursor-pointer transition-colors"
            :class="{ 'bg-primary/5': !notif.is_read }"
          >
            <div class="flex items-start gap-3">
              <!-- Icon -->
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center text-lg"
                :class="getIconClass(notif.type)"
              >
                {{ getIcon(notif.type) }}
              </div>
              <!-- Content -->
              <div class="flex-1 min-w-0">
                <p class="font-medium text-sm text-base truncate">{{ notif.title }}</p>
                <p class="text-xs text-muted line-clamp-2">{{ notif.message }}</p>
                <p class="text-xs text-muted mt-1">{{ formatTime(notif.created_at) }}</p>
              </div>
              <!-- Unread dot -->
              <div v-if="!notif.is_read" class="w-2 h-2 rounded-full bg-primary mt-2"></div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="p-2 border-t border-secondary-200 dark:border-secondary-700">
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

    <!-- Backdrop -->
    <div v-if="isOpen" @click="isOpen = false" class="fixed inset-0 z-40"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { notificationAPI } from '~/composables/useApi'

const isOpen = ref(false)
const loading = ref(false)
const notifications = ref([])
const unreadCount = ref(0)

let pollInterval = null

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value && notifications.value.length === 0) {
    fetchNotifications()
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
  // Could navigate to related page based on notif.type and notif.data
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

onMounted(() => {
  fetchUnreadCount()
  // Poll every 30 seconds
  pollInterval = setInterval(fetchUnreadCount, 30000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
