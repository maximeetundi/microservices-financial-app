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
          @click="openNotificationDetail(notif)"
          :class="['notif-item', { unread: !notif.is_read }]"
        >
          <div class="notif-icon">{{ getTypeIcon(notif.type) }}</div>
          <div class="notif-body">
            <div class="notif-title">{{ notif.title }}</div>
            <div class="notif-msg">{{ notif.message }}</div>
            <div class="notif-time">{{ formatTime(notif.created_at) }}</div>
            
            <!-- Prominent action button for ANY notification with action_url -->
            <button v-if="notif.action_url" 
              @click.stop="navigateTo(notif.action_url)" 
              class="action-btn-primary">
              {{ getActionLabel(notif) }} →
            </button>
          </div>
          <div class="notif-actions">
            <button v-if="!notif.is_read" @click.stop="markAsRead(notif.id)" class="action-mark" title="Marquer comme lu">✓</button>
            <button @click.stop="deleteNotification(notif.id)" class="action-delete" title="Supprimer">✕</button>
          </div>
        </div>
      </div>

      <!-- Notification Detail Modal -->
      <Teleport to="body">
        <div v-if="selectedNotification" class="modal-overlay" @click.self="closeDetailModal">
          <div class="modal-content">
            <button @click="closeDetailModal" class="modal-close">✕</button>
            <div class="modal-icon-large">{{ getTypeIcon(selectedNotification.type) }}</div>
            <h3 class="modal-title">{{ selectedNotification.title }}</h3>
            <p class="modal-message">{{ selectedNotification.message }}</p>
            <div class="modal-meta">
              <span class="modal-type">{{ getTypeLabel(selectedNotification.type) }}</span>
              <span class="modal-date">{{ formatFullDate(selectedNotification.created_at) }}</span>
            </div>
            <div v-if="selectedNotification.action_url" class="modal-actions mb-4">
              <button @click="navigateTo(selectedNotification.action_url)" class="modal-btn-primary w-full">
                Voir les détails →
              </button>
            </div>
            <button @click="closeDetailModal" class="modal-btn">Fermer</button>
          </div>
        </div>
      </Teleport>

      <!-- Load More -->
      <div v-if="hasMore && !loading && notifications.length > 0" class="load-more">
        <button @click="loadMore">Charger plus ↓</button>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useApi } from '~/composables/useApi'

const { notificationApi } = useApi()

const notifications = ref([])
const unreadCount = ref(0)
const loading = ref(false)
const offset = ref(0)
const limit = 20
const hasMore = ref(true)
const activeFilter = ref('all')
const selectedNotification = ref(null)
const router = useRouter()

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

const NOTIFICATION_TRANSLATIONS = {
  'wallet.created': {
    title: '🎉 Portefeuille créé',
    message: 'Votre nouveau portefeuille a été activé avec succès.'
  },
  'wallet.updated': {
    title: '👛 Portefeuille mis à jour',
    message: 'Les informations de votre portefeuille ont été modifiées.'
  },
  'user.balance_updated': {
    title: '💰 Solde mis à jour',
    message: 'Le solde de votre compte a changé suite à une opération.'
  },
  'transfer.created': {
    title: '💸 Nouveau transfert',
    message: 'Un transfert a été initié.'
  },
  'transfer.completed': {
    title: '✅ Transfert réussi',
    message: 'Le transfert a été traité avec succès.'
  },
  'transfer.received': {
    title: '📥 Fonds reçus',
    message: 'Vous avez reçu de l\'argent sur votre compte.'
  },
  'transfer.failed': {
    title: '❌ Échec du transfert',
    message: 'Le transfert n\'a pas pu être finalisé.'
  },
  'security.login': {
    title: '🔐 Nouvelle connexion',
    message: 'Une nouvelle connexion a été détectée.'
  },
  'security.password_changed': {
    title: '🔐 Mot de passe modifié',
    message: 'Votre mot de passe a été mis à jour avec succès.'
  },
  'kyc.submitted': {
    title: '📋 KYC Soumis',
    message: 'Vos documents sont en cours d\'examen.'
  },
  'kyc.approved': {
    title: '✅ Identité vérifiée',
    message: 'Félicitations ! Votre identité a été vérifiée.'
  },
  'kyc.rejected': {
    title: '⚠️ Vérification refusée',
    message: 'Veuillez vérifier les documents requis et réessayer.'
  },
  'card.created': {
    title: '💳 Carte créée',
    message: 'Votre nouvelle carte virtuelle est prête.'
  },
  'card.frozen': {
    title: '❄️ Carte gelée',
    message: 'Votre carte a été temporairement bloquée.'
  },
  'card.unfrozen': {
    title: '🔥 Carte débloquée',
    message: 'Votre carte est à nouveau active.'
  }
}

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
        
        // Translation Logic
        let title = n.title
        let message = n.message

        // Helper to find translation
        const getTrans = (key) => NOTIFICATION_TRANSLATIONS[key]
        
        // Translate Title
        const titleTrans = getTrans(n.title)
        if (titleTrans) {
            title = titleTrans.title
        } else {
            // Fallback: prettier formatting if it looks like a key (has dots or underscores)
            if (title && (title.includes('.') || title.includes('_'))) {
                title = title.replace(/[._]/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
            }
        }

        // Translate Message
        const msgTrans = getTrans(n.message)
        if (msgTrans) {
            message = msgTrans.message
        } else if (n.message === n.title && titleTrans) {
            // If message is same as title (raw key repeated), use the translated message from title
            message = titleTrans.message
        } else if (message && (message.includes('.') && !message.includes(' '))) {
             // Fallback for message if it looks like a key and has no translation
             message = message.replace(/[._]/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
        }

        // Action URL Logic
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
                url = `/approve/${meta.employee_id}?type=invitation`
              } else if (meta.action === 'approve_action' && meta.approval_id) {
                url = `/approve/${meta.approval_id}?type=approval`
              } else {
                url = '/enterprise'
              }
              break
            default:
              // Check for special actions
              if (meta.action === 'accept_invitation' && meta.employee_id) {
                url = `/approve/${meta.employee_id}?type=invitation`
              } else if (meta.action === 'approve_action' && meta.approval_id) {
                url = `/approve/${meta.approval_id}?type=approval`
              }
              break
          }
        }
        
        return {
            ...n,
            title,
            message,
            meta,
            action_url: url
        }
    })
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const [notifRes, countRes] = await Promise.all([
      notificationApi.getAll({ limit, offset: 0 }),
      notificationApi.getUnreadCount()
    ])
    const rawList = notifRes.data.notifications || []
    notifications.value = processNotifications(rawList)
    unreadCount.value = countRes.data.unread_count || 0
    offset.value = limit
    hasMore.value = rawList.length >= limit
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
    const rawList = res.data.notifications || []
    const processed = processNotifications(rawList)
    notifications.value.push(...processed)
    offset.value += limit
    hasMore.value = rawList.length >= limit
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

const navigateTo = (url) => {
    if(!url) return
    router.push(url)
    closeDetailModal()
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

const formatFullDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('fr-FR', {
    weekday: 'long',
    year: 'numeric', 
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getTypeLabel = (type) => {
  const labels = { transfer: 'Transfert', card: 'Carte', security: 'Sécurité', wallet: 'Portefeuille', kyc: 'Vérification' }
  return labels[type] || 'Notification'
}

// Action label based on notification type
const getActionLabel = (notif) => {
  if (notif.meta?.action === 'accept_invitation') return '✅ Accepter l\'invitation'
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

const openNotificationDetail = (notification) => {
  selectedNotification.value = notification
  if (!notification.is_read) {
    markAsRead(notification.id)
  }
}

const closeDetailModal = () => {
  selectedNotification.value = null
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

.action-btn-primary {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

.action-btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.notif-actions {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex-shrink: 0;
}

.action-mark, .action-delete, .action-view {
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

.action-view {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
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

/* Notification Detail Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 1rem;
}

.modal-content {
  background: linear-gradient(180deg, #1e1e2e, #15151f);
  border-radius: 1.5rem;
  padding: 2rem;
  max-width: 400px;
  width: 100%;
  text-align: center;
  position: relative;
  border: 1px solid rgba(255,255,255,0.1);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.modal-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: rgba(255,255,255,0.1);
  border: none;
  color: #888;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  cursor: pointer;
  font-size: 0.875rem;
}

.modal-close:hover {
  background: rgba(255,255,255,0.2);
  color: #fff;
}

.modal-icon-large {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
  margin-bottom: 0.75rem;
}

.modal-message {
  font-size: 0.9375rem;
  color: #ccc;
  line-height: 1.6;
  margin-bottom: 1.25rem;
  white-space: pre-wrap;
  word-break: break-word;
}

.modal-meta {
  display: flex;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  font-size: 0.75rem;
}

.modal-type {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
}

.modal-date {
  color: #666;
}

.modal-btn {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.modal-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.modal-btn-primary {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
  margin-bottom: 0.5rem;
}

.modal-btn-primary:hover {
  transform: scale(1.02);
}
</style>
