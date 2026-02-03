<template>
  <NuxtLayout name="dashboard">
    <div class="h-[calc(100vh-60px)] md:h-[calc(100vh-120px)] flex bg-white dark:bg-gray-900 rounded-none md:rounded-2xl overflow-hidden border-0 md:border border-gray-200 dark:border-gray-800 shadow-none md:shadow-xl max-w-full">
      <!-- Sidebar Conversations - Hide on mobile when chat is selected -->
      <div :class="['w-full md:w-96 border-r border-gray-200 dark:border-gray-700 flex flex-col bg-white dark:bg-gray-900', selectedConv ? 'hidden md:flex' : 'flex']">
        <!-- Header -->
        <div class="bg-gradient-to-r from-green-500 to-green-600 px-4 py-4 flex items-center justify-between">
          <h2 class="text-xl font-bold text-white">Messages</h2>
          <div class="flex items-center gap-2">
            <button @click="showContactsModal = true" class="w-10 h-10 rounded-full bg-white/20 hover:bg-white/30 text-white flex items-center justify-center transition-all hover:scale-105" title="Mes Contacts">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
            <button @click="showNewConversationModal = true" class="w-10 h-10 rounded-full bg-white/20 hover:bg-white/30 text-white flex items-center justify-center transition-all hover:scale-105" title="Nouvelle conversation">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Search -->
        <div class="p-3 bg-gray-50 dark:bg-gray-800">
          <div class="relative">
            <input v-model="searchQuery" type="text" placeholder="Rechercher..."
              class="w-full pl-10 pr-4 py-2.5 rounded-full bg-white dark:bg-gray-700 border-0 text-sm focus:ring-2 focus:ring-green-500 shadow-sm" />
            <svg class="w-5 h-5 absolute left-3.5 top-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>

        <!-- Conversations List -->
        <div class="flex-1 overflow-y-auto">
          <div v-for="conv in filteredConversations" :key="conv.id" @click="selectConversation(conv)"
            :class="['group flex items-center gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 cursor-pointer transition-all border-l-4', 
              selectedConv?.id === conv.id ? 'bg-green-50 dark:bg-green-900/20 border-green-500' : 'border-transparent hover:border-gray-200']">
            
            <!-- Avatar with online indicator -->
            <div class="relative flex-shrink-0">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-green-400 to-green-600 text-white flex items-center justify-center font-bold text-lg shadow-md">
                {{ getOtherParticipantName(conv)?.[0]?.toUpperCase() || 'U' }}
              </div>
              <!-- Online indicator -->
              <div v-if="getParticipantStatus(conv) === 'online'" 
                class="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-500 border-2 border-white dark:border-gray-900 rounded-full"></div>
            </div>
            
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-baseline">
                <h3 class="font-semibold text-gray-900 dark:text-white truncate">{{ getOtherParticipantName(conv) }}</h3>
                <span class="text-xs text-gray-500 flex-shrink-0 ml-2"><ClientOnly>{{ formatTime(conv.lastMessageAt || conv.updated_at) }}</ClientOnly></span>
              </div>
              <div class="flex items-center gap-1">
                <!-- Typing indicator in list -->
                <span v-if="typingUsers[getOtherParticipantId(conv)]" class="text-sm text-green-600 dark:text-green-400 italic">
                  écrit...
                </span>
                <!-- Last message with read status -->
                <template v-else>
                  <svg v-if="conv.lastMessageMine && conv.lastMessageRead" class="w-4 h-4 text-blue-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <p class="text-sm text-gray-500 dark:text-gray-400 truncate">{{ conv.lastMessage || conv.last_message?.content || 'Nouvelle conversation' }}</p>
                </template>
              </div>
            </div>
            
            <!-- Unread badge -->
            <div v-if="(conv.unread_count || conv.unreadCount) && (conv.unread_count || conv.unreadCount) > 0" 
              class="w-6 h-6 rounded-full bg-green-500 text-white text-xs flex items-center justify-center font-bold shadow-sm">
              {{ (conv.unread_count || conv.unreadCount) > 9 ? '9+' : (conv.unread_count || conv.unreadCount) }}
            </div>

            
            <!-- Delete button -->
            <button @click.stop="deleteConversation(conv)" 
              class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-full opacity-0 group-hover:opacity-100 transition-all" title="Supprimer">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
          
          <div v-if="filteredConversations.length === 0" class="p-8 text-center">
            <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
              <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <p class="font-medium text-gray-900 dark:text-white mb-1">Aucune conversation</p>
            <p class="text-sm text-gray-500">Cliquez sur + pour commencer</p>
          </div>
        </div>
      </div>

      <!-- Chat Area -->
      <div v-if="selectedConv" class="flex-1 flex flex-col w-full max-w-full overflow-hidden bg-gray-50 dark:bg-gray-900">
        <!-- Chat Header -->
        <div class="bg-white dark:bg-gray-800 px-3 md:px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3 shadow-sm">
          <!-- Back button for mobile -->
          <button @click="goBackToList" class="md:hidden p-2 -ml-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors">
            <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          
          <!-- Avatar with status -->
          <div class="relative">
            <div class="w-11 h-11 rounded-full bg-gradient-to-br from-green-400 to-green-600 text-white flex items-center justify-center font-bold shadow-md">
              {{ getOtherParticipantName(selectedConv)?.[0]?.toUpperCase() || 'U' }}
            </div>
            <!-- Online indicator -->
            <div v-if="selectedUserStatus === 'En ligne'" 
              class="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white dark:border-gray-800 rounded-full"></div>
          </div>
          
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-gray-900 dark:text-white truncate">{{ getOtherParticipantName(selectedConv) }}</h3>
            <p :class="['text-xs transition-colors', 
              isOtherTyping ? 'text-green-600 dark:text-green-400 font-medium' :
              selectedUserStatus === 'En ligne' ? 'text-green-500' : 'text-gray-500']">
              {{ isOtherTyping ? 'écrit...' : selectedUserStatus }}
            </p>
          </div>
          
          <!-- Actions -->
          <div class="flex items-center gap-1">
            <button class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors">
              <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
            </button>
            <button class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors">
              <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
            <button class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors">
              <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Messages Area -->
        <div ref="messagesContainer" class="flex-1 overflow-y-auto overflow-x-hidden p-3 md:p-4 space-y-1" 
          :style="{ backgroundImage: `url('data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' viewBox=\'0 0 60 60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%239C92AC\' fill-opacity=\'0.03\'%3E%3Cpath d=\'M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z\'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E')` }">
          
          <!-- Date separator -->
          <div class="flex items-center justify-center my-4">
            <span class="px-3 py-1 bg-white dark:bg-gray-800 rounded-full text-xs text-gray-500 shadow-sm">
              Aujourd'hui
            </span>
          </div>
          
          <MessageBubble 
            v-for="msg in messages" 
            :key="msg.id" 
            :message="msg" 
            :show-sender-name="false"
            @openImage="openImageModal" 
          />
          
          <!-- Typing indicator -->
          <TypingIndicator v-if="isOtherTyping" :name="getOtherParticipantName(selectedConv)" />
        </div>

        <!-- Input Component -->
        <MessageInput 
          :conversation-id="selectedConv?.id" 
          @messageSent="handleMessageSent"
          @typing="handleMyTyping"
        />
      </div>

      <!-- Empty State -->
      <div v-else class="hidden md:flex flex-1 items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
        <div class="text-center">
          <div class="w-32 h-32 mx-auto mb-6 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center shadow-xl">
            <svg class="w-16 h-16 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Vos Messages</h3>
          <p class="text-gray-500 max-w-xs mx-auto">Sélectionnez une conversation ou démarrez une nouvelle discussion</p>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <NewConversationModal :show="showNewConversationModal" @close="showNewConversationModal = false" @userSelected="handleUserSelected" />
    <ContactsModal :show="showContactsModal" @close="showContactsModal = false" @startConversation="handleContactConversation" />

    <!-- Image Modal -->
    <Teleport to="body">
      <div v-if="imageModalUrl" class="fixed inset-0 z-[9999] bg-black/95 flex items-center justify-center p-4" @click="closeImageModal">
        <button @click="closeImageModal" class="absolute top-4 right-4 w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition-colors">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <img :src="imageModalUrl" class="max-w-full max-h-full object-contain rounded-lg shadow-2xl" @click.stop />
      </div>
    </Teleport>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import api, { contactsAPI, messagingAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'
import MessageBubble from '~/components/messages/MessageBubble.vue'
import MessageInput from '~/components/messages/MessageInput.vue'
import TypingIndicator from '~/components/messages/TypingIndicator.vue'
import NewConversationModal from '~/components/messages/NewConversationModal.vue'
import ContactsModal from '~/components/messages/ContactsModal.vue'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

// WebSocket connection
let ws: WebSocket | null = null
let presenceInterval: ReturnType<typeof setInterval> | null = null
let chatActivityInterval: ReturnType<typeof setInterval> | null = null

const searchQuery = ref('')
const selectedConv = ref<any>(null)
const messages = ref<any[]>([])
const messagesContainer = ref<HTMLElement>()
const imageModalUrl = ref<string | null>(null)
const showNewConversationModal = ref(false)
const showContactsModal = ref(false)
const userConversations = ref<any[]>([])
const onlineStatus = ref<Record<string, string>>({})
const typingUsers = ref<Record<string, boolean>>({})
const currentUserId = ref<string>('')
const syncedContacts = ref<Array<{id: string, phone: string, email: string, name: string}>>([])
const isOtherTyping = ref(false)

// Initialize auth store
const authStore = useAuthStore()

// Filtered conversations
const filteredConversations = computed(() => {
  if (!searchQuery.value.trim()) return userConversations.value
  const query = searchQuery.value.toLowerCase()
  return userConversations.value.filter(conv => {
    const name = getOtherParticipantName(conv)?.toLowerCase() || ''
    return name.includes(query)
  })
})

// Get current user ID from auth store
const getCurrentUserId = () => {
  if (authStore.user?.id) return authStore.user.id
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
  return `${protocol}//api.app.tech-afm.com/messaging-service/ws/chat`
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
      if (selectedConv.value?.id === data.conversation_id) {
        messages.value.push({
          ...data.content,
          isMine: data.sender_id === currentUserId.value
        })
        nextTick(() => scrollToBottom())
        
        // Mark as read
        if (data.sender_id !== currentUserId.value) {
          markMessageAsRead(data.content.id)
        }
      }
      updateConversationLastMessage(data.conversation_id, data.content)
      break
      
    case 'typing':
      if (data.conversation_id === selectedConv.value?.id && data.user_id !== currentUserId.value) {
        isOtherTyping.value = data.is_typing
        typingUsers.value[data.user_id] = data.is_typing
        
        // Auto-clear after 3 seconds
        if (data.is_typing) {
          setTimeout(() => {
            isOtherTyping.value = false
            typingUsers.value[data.user_id] = false
          }, 3000)
        }
      }
      break
      
    case 'read':
      // Update message read status
      if (data.conversation_id === selectedConv.value?.id) {
        messages.value.forEach(msg => {
          if (msg.id === data.message_id) {
            msg.read_at = data.read_at
          }
        })
      }
      break
      
    case 'presence':
      if (data.user_id && data.status) {
        onlineStatus.value[data.user_id] = data.status
      }
      break
  }
}

// Mark message as read
const markMessageAsRead = async (messageId: string) => {
  try {
    await messagingAPI.markAsRead(messageId)
  } catch (e) {
    // Ignore
  }
}

// Send typing indicator
const handleMyTyping = (isTyping: boolean) => {
  if (ws?.readyState === WebSocket.OPEN && selectedConv.value) {
    ws.send(JSON.stringify({
      type: 'typing',
      conversation_id: selectedConv.value.id,
      is_typing: isTyping
    }))
  }
}

// Update conversation last message
const updateConversationLastMessage = (convId: string, message: any) => {
  const index = userConversations.value.findIndex(c => c.id === convId)
  if (index > -1) {
    userConversations.value[index].lastMessage = message.content
    userConversations.value[index].lastMessageAt = new Date().toISOString()
    userConversations.value[index].lastMessageMine = message.sender_id === currentUserId.value
    const conv = userConversations.value.splice(index, 1)[0]
    userConversations.value.unshift(conv)
  }
}

// Enrich conversation participants
const enrichConversationParticipants = async () => {
  const userIdsToFetch = new Set<string>()
  
  for (const conv of userConversations.value) {
    for (const p of conv.participants || []) {
      if (p.user_id !== currentUserId.value && !p.phone && !p.email) {
        userIdsToFetch.add(p.user_id)
      }
    }
  }
  
  if (userIdsToFetch.size === 0) return
  
  const userDataMap: Record<string, any> = {}
  
  for (const userId of userIdsToFetch) {
    try {
      const res = await api.get(`/auth-service/api/v1/users/${userId}`)
      if (res.data) {
        userDataMap[userId] = res.data
      }
    } catch {}
  }
  
  for (const conv of userConversations.value) {
    for (const p of conv.participants || []) {
      if (p.user_id !== currentUserId.value && userDataMap[p.user_id]) {
        const userData = userDataMap[p.user_id]
        if (!p.phone && userData.phone) p.phone = userData.phone
        if (!p.email && userData.email) p.email = userData.email
      }
    }
  }
}

const updatePresence = async () => {
  try {
    await api.post('/auth-service/api/v1/users/presence')
  } catch {}
}

// Get the other participant's ID
const getOtherParticipantId = (conv: any) => {
  const participants = conv.participants || []
  for (const p of participants) {
    const uid = p.user_id || p
    if (uid !== currentUserId.value) return uid
  }
  return ''
}

// Get participant status - first check backend data, then local cache
const getParticipantStatus = (conv: any) => {
  const participants = conv.participants || []
  for (const p of participants) {
    const uid = p.user_id || p
    if (uid !== currentUserId.value) {
      // Check if backend already provided online status
      if (p.online === true) return 'online'
      // Fall back to local cache
      return onlineStatus.value[uid] || 'offline'
    }
  }
  return 'offline'
}


// Get the name of the other participant
const getOtherParticipantName = (conv: any) => {
  if (!conv) return 'Inconnu'
  
  const participants = conv.participants || []
  for (const p of participants) {
    const uid = p.user_id || p
    if (uid !== currentUserId.value) {
      const phone = p.phone
      const email = p.email
      
      if (phone && phone.length > 0) {
        const contact = syncedContacts.value.find(c => c.phone === phone)
        if (contact) return contact.name
        return formatPhoneNumber(phone)
      }
      if (email && email.length > 0) {
        const contact = syncedContacts.value.find(c => c.email === email)
        if (contact) return contact.name
        return email
      }
      return 'Utilisateur'
    }
  }
  
  if (conv.other_user_phone) return formatPhoneNumber(conv.other_user_phone)
  if (conv.other_user_email) return conv.other_user_email
  
  return conv.phone || 'Conversation'
}

// Format phone number
const formatPhoneNumber = (phone: string | undefined) => {
  if (!phone) return ''
  if (phone.length >= 10) {
    return phone.replace(/(\d{3})(\d{2})(\d{2})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
  }
  return phone
}

// Get display status
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
  isOtherTyping.value = false
  loadMessages()
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
    } catch {}
  }
}

const loadMessages = async () => {
  if (!selectedConv.value) return
  try {
    const res = await messagingAPI.getMessages(selectedConv.value.id)
    const userId = currentUserId.value
    
    const otherParticipantIds = (selectedConv.value.participants || [])
      .filter((p: any) => (p.user_id || p) !== userId)
      .map((p: any) => String(p.user_id || p).trim())
    
    messages.value = (res.data?.messages || []).map((m: any) => {
      const senderId = String(m.sender_id || '').trim()
      const isMine = otherParticipantIds.length > 0 
        ? !otherParticipantIds.includes(senderId)
        : senderId === String(userId).trim()
      
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

const loadConversations = async () => {
  try {
    const res = await messagingAPI.getConversations()
    userConversations.value = res.data?.conversations || []
  } catch (err) {
    console.error('Failed to load conversations:', err)
  }
}

const goBackToList = () => {
  selectedConv.value = null
  messages.value = []
}

const handleUserSelected = (conversation: any) => {
  userConversations.value.unshift(conversation)
  selectConversation(conversation)
}

const handleContactConversation = async (user: any) => {
  showContactsModal.value = false
  try {
    const res = await messagingAPI.createConversation({
      participant_id: user.id,
      participant_name: user.contactName || user.name || 'Utilisateur'
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
    await messagingAPI.deleteConversation(conv.id)
    userConversations.value = userConversations.value.filter(c => c.id !== conv.id)
    
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
  if (selectedConv.value) {
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

const setChatActivity = async (active: boolean) => {
  try {
    await api.post('/auth-service/api/v1/users/chat-activity', { active })
  } catch {}
}

onMounted(async () => {
  if (!authStore.user) {
    await authStore.initializeAuth()
  }
  
  currentUserId.value = authStore.user?.id || getCurrentUserId()
  
  try {
    const [convRes, contactsRes] = await Promise.all([
      messagingAPI.getConversations(),
      contactsAPI.getAll()
    ])
    userConversations.value = convRes.data?.conversations || []
    syncedContacts.value = contactsRes.data?.contacts || []
    
    await enrichConversationParticipants()
  } catch (err) {
    console.error(err)
  }
  
  connectWebSocket()
  
  updatePresence()
  presenceInterval = setInterval(updatePresence, 60000)
  
  setChatActivity(true)
  chatActivityInterval = setInterval(() => setChatActivity(true), 30000)
})

onUnmounted(() => {
  if (ws) {
    ws.close(1000)
    ws = null
  }
  if (presenceInterval) clearInterval(presenceInterval)
  if (chatActivityInterval) clearInterval(chatActivityInterval)
  
  setChatActivity(false)
})
</script>
