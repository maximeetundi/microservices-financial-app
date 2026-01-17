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

    <!-- Modal Overlay -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="isOpen" class="notification-overlay" @click.self="isOpen = false">
          <div class="notification-modal">
            <!-- Header -->
            <div class="notification-header">
              <div class="header-left">
                <div class="header-icon">🔔</div>
                <h3 class="header-title">Notifications</h3>
                <span v-if="unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
              </div>
              <div class="header-actions">
                <button 
                  v-if="unreadCount > 0" 
                  @click="markAllRead" 
                  class="mark-all-btn"
                >
                  ✓ Tout marquer lu
                </button>
                <button @click="isOpen = false" class="close-btn">✕</button>
              </div>
            </div>

            <!-- Loading -->
            <div v-if="loading" class="loading-state">
              <div class="loading-spinner"></div>
              <p>Chargement...</p>
            </div>

            <!-- Empty State -->
            <div v-else-if="notifications.length === 0" class="empty-state">
              <div class="empty-icon">🔔</div>
              <p class="empty-title">Aucune notification</p>
              <p class="empty-subtitle">Vous serez notifié des nouvelles activités ici</p>
            </div>

            <!-- Notification List -->
            <div v-else class="notification-list">
              <div 
                v-for="notif in notifications" 
                :key="notif.id"
                class="notification-item"
                :class="{ 'unread': !notif.is_read }"
              >
                <div class="notification-icon" :class="getIconClass(notif.type)">
                  {{ getIcon(notif.type) }}
                </div>
                <div class="notification-content" @click="handleClick(notif)">
                  <p class="notification-title">{{ notif.title }}</p>
                  <p class="notification-message">{{ notif.message }}</p>
                  <p class="notification-time">{{ formatTime(notif.created_at) }}</p>
                  
                  <!-- Action Button for notifications with action_url -->
                  <button 
                    v-if="notif.action_url" 
                    @click.stop="navigateToAction(notif)"
                    class="action-button"
                  >
                    {{ getActionLabel(notif) }} →
                  </button>
                </div>
                <div v-if="!notif.is_read" class="unread-dot"></div>
              </div>
            </div>

            <!-- Footer -->
            <div class="notification-footer">
              <NuxtLink 
                to="/notifications" 
                class="view-all-link"
                @click="isOpen = false"
              >
                Voir toutes les notifications →
              </NuxtLink>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { notificationAPI } from '~/composables/useApi'
import { useWalletStore } from '~/stores/wallet'

const router = useRouter()
const walletStore = useWalletStore()
const isOpen = ref(false)
const loading = ref(false)
const notifications = ref([])
const unreadCount = ref(0)
const bellButton = ref(null)

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
    notifications.value = processNotifications(res.data?.notifications || [])
  } catch (e) {
    console.error('Failed to fetch notifications:', e)
  } finally {
    loading.value = false
  }
}

// Process notifications to extract action_url from data
const processNotifications = (list) => {
  return list.map(n => {
    let meta = {}
    try {
      if (n.data) {
        meta = typeof n.data === 'string' ? JSON.parse(n.data) : n.data
      }
    } catch (e) {
      console.warn('Failed to parse notification data', e)
    }
    
    // Start with existing action_url or link from meta
    let url = n.action_url || meta.link || null
    
    // Generate action_url based on notification type if not already set
    if (!url) {
      const type = n.type?.toLowerCase() || ''
      const refId = meta.transfer_id || meta.reference || meta.reference_id || meta.id
      
      switch (type) {
        case 'transfer':
        case 'transaction':
        case 'payment':
          url = refId ? `/transactions?ref=${refId}` : '/transactions'
          break
        case 'wallet':
          url = '/wallet'
          break
        case 'card':
          url = '/cards'
          break
        case 'security':
        case 'kyc':
          url = '/settings'
          break
        case 'enterprise':
          if (meta.action === 'accept_invitation' && meta.employee_id) {
            url = `/enterprise/accept?id=${meta.employee_id}`
          } else {
            url = '/enterprise'
          }
          break
        default:
          // Check for special actions
          if (meta.action === 'accept_invitation' && meta.employee_id) {
            url = `/enterprise/accept?id=${meta.employee_id}`
          }
          break
      }
    }
    
    return { ...n, meta, action_url: url }
  })
}

// Get action label based on notification type
const getActionLabel = (notif) => {
  if (notif.meta?.action === 'accept_invitation') return '✅ Accepter'
  const type = notif.type?.toLowerCase() || ''
  switch (type) {
    case 'transfer':
    case 'transaction':
    case 'payment':
      return '📋 Voir transaction'
    case 'wallet':
      return '👛 Voir portefeuille'
    case 'card':
      return '💳 Voir carte'
    case 'security':
      return '🔐 Paramètres'
    case 'kyc':
      return '✅ Vérification'
    case 'enterprise':
      return '🏢 Entreprise'
    default:
      return '👁️ Voir détails'
  }
}

// Navigate to action URL
const navigateToAction = async (notif) => {
  if (!notif.is_read) {
    try {
      await notificationAPI.markAsRead(notif.id)
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (e) {
      console.error('Failed to mark as read:', e)
    }
  }
  isOpen.value = false
  if (notif.action_url) {
    router.push(notif.action_url)
  }
}

const fetchUnreadCount = async () => {
  try {
    const res = await notificationAPI.getUnreadCount()
    let count = res.data?.unread_count || 0
    
    // Check for new notifications to trigger balance updates
    if (count > unreadCount.value) {
       // Fetch latest notifications to identify type
       try {
         const latestRes = await notificationAPI.getAll(5, 0)
         const latestNotifs = latestRes.data?.notifications || []
         
         const hasBalanceUpdate = latestNotifs.some(n => 
           !n.is_read && ['payment', 'wallet', 'transfer', 'transaction', 'deposit'].includes(n.type)
         )
         
         if (hasBalanceUpdate) {
           console.log('New transaction notification detected, refreshing balance...')
           walletStore.fetchWallets() 
           // Also refresh transactions list if we are on that page? 
           // For now just wallets (balances)
         }
       } catch (err) {
         console.error('Failed to check notification type', err)
       }
    }

    // If user is on support chat page, auto-mark support notifications as read
    // to avoid notification sounds/badges during active chat
    const currentPath = router.currentRoute.value?.path || ''
    if (currentPath.startsWith('/support/chat') && count > 0) {
      // Fetch notifications to mark support ones as read
      const notifRes = await notificationAPI.getAll(20, 0)
      const notifications = notifRes.data?.notifications || []
      const supportNotifs = notifications.filter(n => 
        !n.is_read && (n.type === 'support' || n.type === 'conversation' || n.type === 'ticket')
      )
      
      // Mark support notifications as read silently
      for (const notif of supportNotifs) {
        try {
          await notificationAPI.markAsRead(notif.id)
          count = Math.max(0, count - 1)
        } catch {
          // Ignore errors
        }
      }
    }
    
    unreadCount.value = count
  } catch (e) {
    // Silently fail
  }
}

const handleClick = async (notif) => {
  // Mark as read
  if (!notif.is_read) {
    try {
      await notificationAPI.markAsRead(notif.id)
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (e) {
      console.error('Failed to mark as read:', e)
    }
  }
  
  // Navigate based on notification type and reference
  isOpen.value = false
  
  // Prioritize action_url if present
  if (notif.action_url) {
      router.push(notif.action_url)
      return
  }
  
  const type = notif.type?.toLowerCase() || ''
  const refId = notif.reference_id || notif.data?.reference_id || notif.data?.id
  
  if (type === 'support' || type === 'conversation' || type === 'ticket') {
    // Navigate to support conversation
    if (refId) {
      router.push(`/support/chat?id=${refId}&agent_type=human`)
    } else {
      router.push('/support')
    }
  } else if (type === 'transfer' || type === 'transaction') {
    // Navigate to transaction details
    if (refId) {
        router.push(`/transactions?ref=${refId}`)
    } else {
        router.push('/transactions') // or /wallet if transactions page doesn't exist broadly, but transactions is better
    }
  } else if (type === 'card') {
    router.push('/cards')
  } else if (type === 'wallet' || type === 'payment') {
    if (refId && type === 'payment') {
         router.push(`/transactions?ref=${refId}`)
    } else {
         router.push('/wallet')
    }
  } else if (type === 'security') {
    router.push('/settings')
  } else if (type === 'kyc') {
    router.push('/settings')
  } else {
    // Default: go to notifications page
    router.push('/notifications')
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
    transfer: 'icon-green',
    card: 'icon-purple',
    security: 'icon-red',
    promotion: 'icon-yellow',
    wallet: 'icon-blue',
    payment: 'icon-indigo',
    default: 'icon-gray'
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

// Close on escape key
const handleKeydown = (e) => {
  if (e.key === 'Escape' && isOpen.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  fetchUnreadCount()
  pollInterval = setInterval(fetchUnreadCount, 30000)
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.notification-center {
  position: relative;
  display: inline-block;
}

/* Overlay - centered modal */
.notification-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  padding: 20px;
}

/* Modal container */
.notification-modal {
  width: 100%;
  max-width: 450px;
  max-height: calc(100vh - 40px);
  background: linear-gradient(145deg, #1e1e2e 0%, #181825 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Header */
.notification-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(139, 92, 246, 0.1) 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  font-size: 24px;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.unread-badge {
  background: #ef4444;
  color: white;
  font-size: 12px;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 10px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mark-all-btn {
  font-size: 12px;
  color: #a5b4fc;
  background: rgba(99, 102, 241, 0.2);
  border: none;
  padding: 6px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.mark-all-btn:hover {
  background: rgba(99, 102, 241, 0.4);
  color: #fff;
}

.close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: #888;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

/* Loading state */
.loading-state {
  padding: 60px 20px;
  text-align: center;
  color: #888;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  margin: 0 auto 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Empty state */
.empty-state {
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: #fff;
  margin-bottom: 4px;
}

.empty-subtitle {
  font-size: 13px;
  color: #666;
}

/* Notification list */
.notification-list {
  flex: 1;
  overflow-y: auto;
  max-height: 400px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 20px;
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
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.icon-green { background: rgba(34, 197, 94, 0.2); }
.icon-purple { background: rgba(168, 85, 247, 0.2); }
.icon-red { background: rgba(239, 68, 68, 0.2); }
.icon-yellow { background: rgba(234, 179, 8, 0.2); }
.icon-blue { background: rgba(59, 130, 246, 0.2); }
.icon-indigo { background: rgba(99, 102, 241, 0.2); }
.icon-gray { background: rgba(107, 114, 128, 0.2); }

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  margin-bottom: 3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notification-message {
  font-size: 13px;
  color: #888;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: #666;
  margin-top: 4px;
}

.action-button {
  margin-top: 8px;
  padding: 6px 12px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  font-size: 12px;
  font-weight: 600;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 6px rgba(99, 102, 241, 0.3);
}

.action-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(99, 102, 241, 0.4);
}

.unread-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #6366f1;
  flex-shrink: 0;
  margin-top: 4px;
}

/* Footer */
.notification-footer {
  padding: 14px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  text-align: center;
}

.view-all-link {
  font-size: 14px;
  color: #a5b4fc;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.view-all-link:hover {
  color: #818cf8;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .notification-modal,
.modal-leave-to .notification-modal {
  transform: scale(0.9) translateY(-20px);
  opacity: 0;
}

.modal-enter-active .notification-modal,
.modal-leave-active .notification-modal {
  transition: all 0.3s ease;
}

/* Mobile adjustments */
@media (max-width: 480px) {
  .notification-overlay {
    padding: 10px;
    align-items: flex-end;
  }
  
  .notification-modal {
    max-width: 100%;
    max-height: 70vh;
    border-radius: 20px 20px 0 0;
  }
  
  .header-actions {
    gap: 4px;
  }
  
  .mark-all-btn {
    padding: 4px 8px;
    font-size: 11px;
  }
}
</style>
