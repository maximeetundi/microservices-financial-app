<template>
  <NuxtLayout name="dashboard">
    <div class="h-[calc(100vh-120px)] flex bg-white dark:bg-gray-900 rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-800 shadow-lg">
      <!-- Sidebar Conversations -->
      <div class="w-full md:w-96 border-r border-gray-200 dark:border-gray-800 flex flex-col">
        <!-- Header -->
        <div class="bg-gray-50 dark:bg-gray-800 px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Messages</h2>
          <div class="flex items-center gap-2">
            <button @click="showContactsModal = true" class="w-10 h-10 rounded-full bg-blue-500 hover:bg-blue-600 text-white flex items-center justify-center transition-colors" title="Mes Contacts">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
            <button @click="showNewConversationModal = true" class="w-10 h-10 rounded-full bg-green-500 hover:bg-green-600 text-white flex items-center justify-center transition-colors" title="Nouvelle conversation">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Search -->
        <div class="p-3 bg-gray-50 dark:bg-gray-800">
          <div class="relative">
            <input v-model="searchQuery" type="text" placeholder="Rechercher une conversation..."
              class="w-full pl-10 pr-4 py-2 rounded-lg bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 text-sm focus:ring-2 focus:ring-green-500" />
            <svg class="w-5 h-5 absolute left-3 top-2.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>

        <!-- Tabs -->
        <div class="flex border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
          <button @click="activeConvType = 'users'" :class="['flex-1 py-3 text-sm font-medium transition-colors', activeConvType === 'users' ? 'text-green-600 border-b-2 border-green-600' : 'text-gray-500 hover:text-gray-700']">
            Utilisateurs
          </button>
          <button @click="activeConvType = 'associations'" :class="['flex-1 py-3 text-sm font-medium transition-colors', activeConvType === 'associations' ? 'text-green-600 border-b-2 border-green-600' : 'text-gray-500 hover:text-gray-700']">
            Associations
          </button>
        </div>

        <!-- Conversations List -->
        <div class="flex-1 overflow-y-auto">
          <div v-if="activeConvType === 'users'">
            <div v-for="conv in userConversations" :key="conv.id" @click="selectConversation(conv)"
              :class="['group flex items-center gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors border-b border-gray-100 dark:border-gray-800', selectedConv?.id === conv.id && 'bg-gray-100 dark:bg-gray-800']">
              <div class="w-12 h-12 rounded-full bg-green-500 text-white flex items-center justify-center font-bold">
                {{ getOtherParticipantName(conv)?.[0]?.toUpperCase() || 'U' }}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-baseline">
                  <h3 class="font-medium text-gray-900 dark:text-white truncate">{{ getOtherParticipantName(conv) }}</h3>
                <span class="text-xs text-gray-500"><ClientOnly>{{ formatTime(conv.lastMessageAt) }}</ClientOnly></span>
                </div>
                <p class="text-sm text-gray-500 truncate">{{ conv.lastMessage }}</p>
              </div>
              <div v-if="conv.unreadCount" class="w-6 h-6 rounded-full bg-green-500 text-white text-xs flex items-center justify-center font-bold">
                {{ conv.unreadCount }}
              </div>
              <button @click.stop="deleteConversation(conv)" class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-full opacity-0 group-hover:opacity-100 hover:opacity-100 transition-all" title="Supprimer">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
            
            <div v-if="userConversations.length === 0" class="p-8 text-center text-gray-500">
              <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <p class="font-medium mb-2">Aucune conversation</p>
              <p class="text-sm text-gray-400">Cliquez sur + pour démarrer une nouvelle conversation</p>
            </div>
          </div>

          <div v-else>
            <div v-for="assoc in associationChats" :key="assoc.id" @click="selectAssociationChat(assoc)"
              :class="['flex items-center gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors border-b border-gray-100 dark:border-gray-800', selectedAssoc?.id === assoc.id && 'bg-gray-100 dark:bg-gray-800']">
              <div class="w-12 h-12 rounded-full bg-indigo-500 text-white flex items-center justify-center font-bold">
                👥
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-medium text-gray-900 dark:text-white truncate">{{ assoc.name }}</h3>
                <p class="text-sm text-gray-500 truncate">{{ assoc.total_members }} membres</p>
              </div>
            </div>
          </div>

          <div v-if="(activeConvType === 'users' && userConversations.length === 0) || (activeConvType === 'associations' && associationChats.length === 0)" class="p-8 text-center text-gray-500">
            <p>Aucune conversation</p>
          </div>
        </div>
      </div>

      <!-- Chat Area -->
      <div v-if="selectedConv || selectedAssoc" class="flex-1 flex flex-col">
        <!-- Chat Header -->
        <div class="bg-gray-50 dark:bg-gray-800 px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3">
          <div class="relative">
            <div class="w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center font-bold">
              {{ selectedConv ? getOtherParticipantName(selectedConv)?.[0]?.toUpperCase() || 'U' : '👥' }}
            </div>
            <!-- Online indicator dot -->
            <div v-if="selectedUserStatus === 'En ligne'" class="absolute bottom-0 right-0 w-3 h-3 bg-green-400 border-2 border-white dark:border-gray-800 rounded-full"></div>
          </div>
          <div class="flex-1">
            <h3 class="font-medium text-gray-900 dark:text-white">{{ selectedConv ? getOtherParticipantName(selectedConv) : selectedAssoc?.name }}</h3>
            <p :class="['text-xs', selectedUserStatus === 'En ligne' ? 'text-green-500' : 'text-gray-500']">
              {{ selectedAssoc ? `${selectedAssoc.total_members || 0} membres` : selectedUserStatus || 'Chargement...' }}
            </p>
          </div>
        </div>

        <!-- Messages -->
        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50 dark:bg-gray-900" style="background-image: url('data:image/svg+xml,%3Csvg width=&quot;100&quot; height=&quot;100&quot; xmlns=&quot;http://www.w3.org/2000/svg&quot;%3E%3Cg opacity=&quot;0.05&quot;%3E%3Cpath d=&quot;M10 10h80v80H10z&quot; fill=&quot;none&quot; stroke=&quot;%23000&quot;/%3E%3C/g%3E%3C/svg%3E');">
          <MessageBubble 
            v-for="msg in messages" 
            :key="msg.id" 
            :message="msg" 
            :show-sender-name="!!selectedAssoc"
            @openImage="openImageModal" 
          />
        </div>

        <!-- Input Component -->
        <MessageInput 
          :conversation-id="selectedConv?.id" 
          :association-id="selectedAssoc?.id"
          @messageSent="handleMessageSent" 
        />
      </div>

      <!-- Empty State -->
      <div v-else class="flex-1 flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div class="text-center">
          <svg class="w-32 h-32 mx-auto text-gray-300 dark:text-gray-700 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <h3 class="text-xl font-medium text-gray-600 dark:text-gray-400 mb-2">Sélectionnez une conversation</h3>
          <p class="text-gray-500">Vos messages apparaîtront ici</p>
        </div>
      </div>
    </div>

    <!-- New Conversation Modal -->
    <NewConversationModal :show="showNewConversationModal" @close="showNewConversationModal = false" @userSelected="handleUserSelected" />
    <ContactsModal :show="showContactsModal" @close="showContactsModal = false" @startConversation="handleContactConversation" />

    <!-- Image Modal -->
    <Teleport to="body">
      <div v-if="imageModalUrl" class="fixed inset-0 z-[9999] bg-black/90 flex items-center justify-center p-4" @click="closeImageModal">
        <button @click="closeImageModal" class="absolute top-4 right-4 w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition-colors">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <img :src="imageModalUrl" class="max-w-full max-h-full object-contain" @click.stop />
      </div>
    </Teleport>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { associationAPI, contactsAPI } from '~/composables/useApi'
import api from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'
import MessageBubble from '~/components/messages/MessageBubble.vue'
import MessageInput from '~/components/messages/MessageInput.vue'
import NewConversationModal from '~/components/messages/NewConversationModal.vue'
import ContactsModal from '~/components/messages/ContactsModal.vue'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

// WebSocket connection
let ws: WebSocket | null = null
let presenceInterval: ReturnType<typeof setInterval> | null = null

const searchQuery = ref('')
const activeConvType = ref<'users' | 'associations'>('users')
const selectedConv = ref<any>(null)
const selectedAssoc = ref<any>(null)
const messages = ref<any[]>([])
const messagesContainer = ref<HTMLElement>()
const imageModalUrl = ref<string | null>(null)
const showNewConversationModal = ref(false)
const showContactsModal = ref(false)
const userConversations = ref<any[]>([])
const associationChats = ref<any[]>([])
const onlineStatus = ref<Record<string, string>>({})
const currentUserId = ref<string>('')
const syncedContacts = ref<Array<{id: string, phone: string, email: string, name: string}>>([])

// Initialize auth store
const authStore = useAuthStore()

// Get current user ID from auth store
const getCurrentUserId = () => {
  // Try from Pinia store first
  if (authStore.user?.id) {
    return authStore.user.id
  }
  // Fallback to localStorage (in case store not yet initialized)
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return user.id || ''
  } catch {
    return ''
  }
}

// WebSocket URL
const getWsUrl = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//api.app.maximeetundi.store/messaging-service/ws/chat`
}

// Connect to WebSocket
const connectWebSocket = () => {
  if (ws?.readyState === WebSocket.OPEN) return
  
  const userId = getCurrentUserId()
  if (!userId) return
  
  currentUserId.value = userId
  
  const url = `${getWsUrl()}?user_id=${userId}`
  console.log('WebSocket connecting:', url)
  
  ws = new WebSocket(url)
  
  ws.onopen = () => {
    console.log('WebSocket connected')
  }
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      handleWebSocketMessage(data)
    } catch (e) {
      console.error('WebSocket parse error:', e)
    }
  }
  
  ws.onclose = (event) => {
    console.log('WebSocket disconnected:', event.code)
    // Auto-reconnect after 3 seconds
    if (event.code !== 1000) {
      setTimeout(connectWebSocket, 3000)
    }
  }
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }
}

// Handle incoming WebSocket messages
const handleWebSocketMessage = (data: any) => {
  switch (data.type) {
    case 'new_message':
      // Add new message if we're in the same conversation
      if (selectedConv.value?.id === data.conversation_id) {
        messages.value.push({
          ...data.content,
          isMine: data.sender_id === currentUserId.value
        })
        nextTick(() => scrollToBottom())
      }
      // Update conversation list
      updateConversationLastMessage(data.conversation_id, data.content)
      break
      
    case 'typing':
      // Handle typing indicator
      console.log('User typing in', data.conversation_id)
      break
      
    case 'read':
      // Handle read receipt
      break
      
    case 'presence':
      // Update online status
      if (data.user_id && data.status) {
        onlineStatus.value[data.user_id] = data.status
      }
      break
  }
}

// Update conversation last message
const updateConversationLastMessage = (convId: string, message: any) => {
  const index = userConversations.value.findIndex(c => c.id === convId)
  if (index > -1) {
    userConversations.value[index].last_message = message
    userConversations.value[index].updated_at = new Date().toISOString()
    // Move to top
    const conv = userConversations.value.splice(index, 1)[0]
    userConversations.value.unshift(conv)
  }
}

// Update presence (heartbeat)
// Enrich conversation participants with user phone/email data
const enrichConversationParticipants = async () => {
  // Collect user IDs of other participants who don't have phone/email
  const userIdsToFetch = new Set<string>()
  
  for (const conv of userConversations.value) {
    for (const p of conv.participants || []) {
      if (p.user_id !== currentUserId.value && !p.phone && !p.email) {
        userIdsToFetch.add(p.user_id)
      }
    }
  }
  
  if (userIdsToFetch.size === 0) return
  
  // Fetch user data for each missing participant
  const userDataMap: Record<string, any> = {}
  
  for (const userId of userIdsToFetch) {
    try {
      const res = await api.get(`/auth-service/api/v1/users/${userId}`)
      if (res.data) {
        userDataMap[userId] = res.data
      }
    } catch (e) {
      // Ignore errors for individual user fetches
    }
  }
  
  // Update conversations with fetched user data
  for (const conv of userConversations.value) {
    for (const p of conv.participants || []) {
      if (p.user_id !== currentUserId.value && userDataMap[p.user_id]) {
        const userData = userDataMap[p.user_id]
        // Only update if data is missing
        if (!p.phone && userData.phone) p.phone = userData.phone
        if (!p.email && userData.email) p.email = userData.email
      }
    }
  }
}

const updatePresence = async () => {
  try {
    await api.post('/auth-service/api/v1/users/presence')
  } catch (e: any) {
    // Ignore errors - presence endpoint may not exist
    if (e?.response?.status !== 404) {
      console.debug('Presence update skipped')
    }
  }
}

// Get the name of the other participant in a conversation (not "Moi")
// Priority: synced contact name > phone > email > fallback
const getOtherParticipantName = (conv: any) => {
  if (!conv) return 'Inconnu'
  
  // If conv has participants array, find the one that's not the current user
  const participants = conv.participants || []
  for (const p of participants) {
    const uid = p.user_id || p
    if (uid !== currentUserId.value) {
      const phone = p.phone
      const email = p.email
      
      // Check if this phone/email is in synced contacts
      // For now, we show phone/email first (like WhatsApp when contact is not saved)
      // If user has synced contacts from mobile, the name will be in the contact
      
      // Show phone if available (privacy first - like WhatsApp)
      if (phone && phone.length > 0) {
        // Check syncedContacts for this phone to show contact name
        const contact = syncedContacts.value.find(c => c.phone === phone)
        if (contact) {
          return contact.name
        }
        return formatPhoneNumber(phone)
      }
      // Show email if available
      if (email && email.length > 0) {
        // Check syncedContacts for this email
        const contact = syncedContacts.value.find(c => c.email === email)
        if (contact) {
          return contact.name
        }
        return email
      }
      // Fallback to generic
      return 'Utilisateur'
    }
  }
  
  // If conv has other_user_phone, use it
  if (conv.other_user_phone) return formatPhoneNumber(conv.other_user_phone)
  if (conv.other_user_email) return conv.other_user_email
  
  return conv.phone || 'Conversation'
}

// Get current user name from localStorage
const getCurrentUserName = () => {
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return user.first_name ? `${user.first_name} ${user.last_name || ''}`.trim() : ''
  } catch {
    return ''
  }
}

// Format phone number for display
const formatPhoneNumber = (phone: string | undefined) => {
  if (!phone) return ''
  // Format as +XXX XX XX XX XX (if long enough)
  if (phone.length >= 10) {
    return phone.replace(/(\d{3})(\d{2})(\d{2})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
  }
  return phone
}

// Get display status for selected conversation
const selectedUserStatus = computed(() => {
  if (!selectedConv.value) return ''
  const participants = selectedConv.value.participants || []
  for (const p of participants) {
    const uid = p.user_id || p
    if (uid !== currentUserId.value) {
      const status = onlineStatus.value[uid]
      if (status === 'online') return 'En ligne'
      if (status === 'away') return 'Absent'
      return 'Hors ligne'
    }
  }
  return ''
})

const formatTime = (date: any) => {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const hours = diff / (1000 * 60 * 60)
  
  if (hours < 24) return d.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
  if (hours < 48) return 'Hier'
  return d.toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit' })
}

const selectConversation = (conv: any) => {
  selectedConv.value = conv
  selectedAssoc.value = null
  loadMessages()
  // Load presence for other participant
  loadParticipantPresence(conv)
}

const loadParticipantPresence = async (conv: any) => {
  const participants = conv.participants || []
  const otherIds = participants
    .map((p: any) => p.user_id || p)
    .filter((id: string) => id !== currentUserId.value)
  
  if (otherIds.length > 0) {
    try {
      const res = await api.post('/auth-service/api/v1/users/presence/batch', {
        user_ids: otherIds
      })
      for (const p of (res.data?.presences || [])) {
        onlineStatus.value[p.user_id] = p.status
      }
    } catch (e) {
      // Ignore
    }
  }
}

const selectAssociationChat = (assoc: any) => {
  selectedAssoc.value = assoc
  selectedConv.value = null
  loadAssociationMessages()
}

const loadMessages = async () => {
  if (!selectedConv.value) return
  try {
    const res = await api.get(`/messaging-service/api/v1/conversations/${selectedConv.value.id}/messages`)
    const userId = currentUserId.value
    
    // Get other participant IDs from the conversation
    const otherParticipantIds = (selectedConv.value.participants || [])
      .filter((p: any) => (p.user_id || p) !== userId)
      .map((p: any) => String(p.user_id || p).trim())
    
    console.log('Current user ID:', userId)
    console.log('Other participant IDs:', otherParticipantIds)
    
    messages.value = (res.data?.messages || []).map((m: any) => {
      const senderId = String(m.sender_id || '').trim()
      // Message is mine if sender is NOT one of the other participants
      const isMine = otherParticipantIds.length > 0 
        ? !otherParticipantIds.includes(senderId)
        : senderId === String(userId).trim()
      
      console.log('Message sender:', senderId, '| isMine:', isMine)
      return {
        ...m,
        isMine
      }
    })
    await nextTick()
    scrollToBottom()
  } catch (err) {
    console.error(err)
  }
}

const loadAssociationMessages = async () => {
  if (!selectedAssoc.value) return
  try {
    const res = await api.get(`/messaging-service/api/v1/associations/${selectedAssoc.value.id}/chat`)
    messages.value = (res.data?.messages || []).map((m: any) => ({
      ...m,
      isMine: m.sender_id === currentUserId.value
    }))
    await nextTick()
    scrollToBottom()
  } catch (err) {
    console.error(err)
  }
}

const loadConversations = async () => {
  try {
    const res = await api.get('/messaging-service/api/v1/conversations')
    userConversations.value = res.data?.conversations || []
  } catch (err) {
    console.error('Failed to load conversations:', err)
  }
}

const handleUserSelected = (conversation: any) => {
  userConversations.value.unshift(conversation)
  selectConversation(conversation)
}

const handleContactConversation = async (user: any) => {
  showContactsModal.value = false
  // Create or get existing conversation with this user
  try {
    const res = await api.post('/messaging-service/api/v1/conversations', {
      participant_id: user.id,
      participant_name: user.contactName || user.name || 'Utilisateur',
      participant_email: user.email || '',
      participant_phone: user.phone || '',
      my_name: authStore.user?.firstName || 'Moi'
    })
    const conversation = res.data
    userConversations.value.unshift(conversation)
    selectConversation(conversation)
  } catch (e) {
    console.error('Failed to create conversation:', e)
  }
}

const deleteConversation = async (conv: any) => {
  if (!confirm('Supprimer cette conversation ?')) return
  
  try {
    await api.delete(`/messaging-service/api/v1/conversations/${conv.id}`)
    userConversations.value = userConversations.value.filter(c => c.id !== conv.id)
    
    // Clear selection if this was the selected conversation
    if (selectedConv.value?.id === conv.id) {
      selectedConv.value = null
      messages.value = []
    }
  } catch (e) {
    console.error('Failed to delete conversation:', e)
    alert('Erreur lors de la suppression')
  }
}

const handleMessageSent = () => {
  if (selectedAssoc.value) {
    loadAssociationMessages()
  } else if (selectedConv.value) {
    loadMessages()
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const openImageModal = (url: string) => {
  imageModalUrl.value = url
}

const closeImageModal = () => {
  imageModalUrl.value = null
}

onMounted(async () => {
  // Wait for auth store to be ready
  if (!authStore.user) {
    await authStore.initializeAuth()
  }
  
  // Get user ID from store (now should be available)
  currentUserId.value = authStore.user?.id || getCurrentUserId()
  console.log('Current user ID after init:', currentUserId.value)
  
  try {
    const [convRes, assocRes, contactsRes] = await Promise.all([
      api.get('/messaging-service/api/v1/conversations'),
      associationAPI.getAll(),
      contactsAPI.getAll()
    ])
    userConversations.value = convRes.data?.conversations || []
    associationChats.value = assocRes.data || []
    syncedContacts.value = contactsRes.data?.contacts || []
    
    // Enrich participants with user data (phone/email) for conversations missing this info
    await enrichConversationParticipants()
  } catch (err) {
    console.error(err)
  }
  
  // Connect WebSocket
  connectWebSocket()
  
  // Start presence heartbeat (every 60 seconds)
  updatePresence()
  presenceInterval = setInterval(updatePresence, 60000)
})

onUnmounted(() => {
  // Cleanup
  if (ws) {
    ws.close(1000)
    ws = null
  }
  if (presenceInterval) {
    clearInterval(presenceInterval)
  }
})
</script>
