<template>
  <NuxtLayout name="dashboard">
    <div class="chat-page">
      <!-- Chat Header -->
      <div class="chat-header glass-card">
        <div class="header-left">
          <NuxtLink to="/support" class="back-btn">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </NuxtLink>
          <div class="agent-info">
            <div class="agent-avatar-wrapper">
              <div class="agent-avatar" :class="conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
                {{ conversation.agent_type === 'ai' ? '🤖' : '👤' }}
              </div>
              <div class="online-indicator"></div>
            </div>
            <div>
              <h1 class="agent-name">{{ conversation.agent_type === 'ai' ? 'Assistant IA' : agentName }}</h1>
              <p class="agent-status">
                <span class="status-dot"></span>
                En ligne
              </p>
            </div>
          </div>
        </div>
        <div class="header-actions">
          <button 
            v-if="conversation.agent_type === 'ai' && conversation.status !== 'escalated'"
            @click="escalateToHuman"
            class="btn-escalate"
          >
            👤 Humain
          </button>
          <button @click="closeConversation" class="btn-close">
            Fermer
          </button>
        </div>
      </div>

      <!-- Messages Container -->
      <div ref="messagesContainer" class="messages-container">
        <!-- Date Separator -->
        <div class="date-separator">
          <span>{{ formatDateHeader(new Date()) }}</span>
        </div>

        <!-- Messages -->
        <div 
          v-for="message in messages" 
          :key="message.id"
          class="message-wrapper"
          :class="message.sender_type === 'user' ? 'message-user' : 'message-agent'"
        >
          <!-- Agent/AI Message -->
          <div v-if="message.sender_type !== 'user'" class="message-bubble-wrapper">
            <div class="avatar-small" :class="message.sender_type === 'system' ? 'bg-gray-500/20' : (conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20')">
              {{ message.sender_type === 'system' ? '🔔' : (conversation.agent_type === 'ai' ? '🤖' : '👤') }}
            </div>
            <div>
              <div class="message-bubble agent-bubble">
                <p>{{ message.content }}</p>
              </div>
              <p class="message-time">{{ formatTime(message.created_at) }}</p>
            </div>
          </div>

          <!-- User Message -->
          <div v-else class="message-bubble-wrapper user-wrapper">
            <div>
              <div class="message-bubble user-bubble">
                <p>{{ message.content }}</p>
              </div>
              <p class="message-time user-time">
                {{ formatTime(message.created_at) }}
                <span v-if="message.is_read" class="read-indicator">✓✓</span>
              </p>
            </div>
          </div>
        </div>

        <!-- Typing Indicator -->
        <div v-if="isTyping" class="message-wrapper message-agent">
          <div class="message-bubble-wrapper">
            <div class="avatar-small" :class="conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
              {{ conversation.agent_type === 'ai' ? '🤖' : '👤' }}
            </div>
            <div class="message-bubble agent-bubble typing-bubble">
              <div class="typing-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Replies -->
      <div v-if="quickReplies.length > 0" class="quick-replies">
        <button 
          v-for="reply in quickReplies" 
          :key="reply"
          @click="sendQuickReply(reply)"
          class="quick-reply-btn"
        >
          {{ reply }}
        </button>
      </div>

      <!-- Input Area -->
      <div class="input-area glass-card">
        <form @submit.prevent="sendMessage" class="input-form">
          <textarea
            ref="messageInput"
            v-model="newMessage"
            @keydown.enter.exact.prevent="sendMessage"
            placeholder="Écrivez votre message..."
            rows="1"
            class="message-input input-premium"
            :disabled="sending"
          ></textarea>
          <button 
            type="submit"
            :disabled="!newMessage.trim() || sending"
            class="send-btn btn-premium"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
            </svg>
          </button>
        </form>
      </div>

      <!-- Rating Modal -->
      <Teleport to="body">
        <div v-if="showRatingModal" class="modal-overlay" @click.self="showRatingModal = false">
          <div class="modal-content glass-card">
            <h3 class="modal-title">Comment évaluez-vous cette conversation ?</h3>
            <div class="rating-stars">
              <button 
                v-for="star in 5" 
                :key="star"
                @click="rating = star"
                class="star-btn"
              >
                {{ star <= rating ? '⭐' : '☆' }}
              </button>
            </div>
            <textarea
              v-model="feedback"
              placeholder="Commentaire (optionnel)"
              class="feedback-input input-premium"
              rows="3"
            ></textarea>
            <div class="modal-actions">
              <button @click="submitRating" class="btn-premium">
                Envoyer
              </button>
              <button @click="showRatingModal = false; navigateToSupport()" class="btn-secondary-premium">
                Passer
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { supportAPI } from '~/composables/useApi'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()

const messagesContainer = ref(null)
const messageInput = ref(null)

const conversation = ref({
  id: '',
  subject: '',
  agent_type: 'ai',
  status: 'active'
})

const agentName = ref('Support Zekora')
const messages = ref([])
const newMessage = ref('')
const isTyping = ref(false)
const sending = ref(false)
const showRatingModal = ref(false)
const rating = ref(5)
const feedback = ref('')

const quickReplies = ref([
  'Solde du compte',
  'Frais de transfert',
  'Commander une carte',
  'Parler à un humain'
])

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || sending.value) return

  const content = newMessage.value.trim()
  newMessage.value = ''
  sending.value = true

  // Add user message immediately
  const userMessage = {
    id: 'msg-' + Date.now(),
    sender_type: 'user',
    sender_name: 'Vous',
    content: content,
    created_at: new Date().toISOString(),
    is_read: true
  }
  messages.value.push(userMessage)
  scrollToBottom()

  // Show typing indicator
  isTyping.value = true

  try {
    // Send to backend
    if (conversation.value.id) {
      await supportAPI.sendMessage(conversation.value.id, content)
    }

    // Simulate delay for AI response
    await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 1500))

    // If AI, generate response
    if (conversation.value.agent_type === 'ai') {
      const aiResponse = generateAIResponse(content)
      messages.value.push({
        id: 'msg-' + Date.now(),
        sender_type: 'agent',
        sender_name: 'Assistant IA',
        content: aiResponse,
        created_at: new Date().toISOString()
      })
    }
  } catch (error) {
    console.error('Error sending message:', error)
  } finally {
    isTyping.value = false
    sending.value = false
    scrollToBottom()
    messageInput.value?.focus()
  }
}

const sendQuickReply = (reply) => {
  newMessage.value = reply
  sendMessage()
}

const generateAIResponse = (message) => {
  const lower = message.toLowerCase()
  
  if (lower.includes('solde') || lower.includes('balance')) {
    return "Pour consulter votre solde, rendez-vous sur la page d'accueil de l'application. Votre solde total s'affiche en haut de l'écran.\n\n💡 Astuce : Vous pouvez également voir le détail par devise en cliquant sur 'Mes Wallets'."
  }
  if (lower.includes('frais') || lower.includes('commission')) {
    return "📊 Voici nos frais :\n\n• Transferts SEPA : Gratuit\n• Crypto-Crypto : 0.5%\n• Fiat-Crypto : 0.75%\n• Fiat-Fiat : 0.15-0.25%\n\nNous sommes jusqu'à 8x moins chers que les banques traditionnelles !"
  }
  if (lower.includes('carte')) {
    return "Pour commander votre carte Zekora :\n\n1. Allez dans le menu 'Cartes'\n2. Cliquez sur 'Commander une carte'\n3. Choisissez entre virtuelle (gratuite) ou physique (9.99€)\n4. Suivez les étapes de personnalisation\n\nVotre carte virtuelle est disponible instantanément !"
  }
  if (lower.includes('humain') || lower.includes('agent')) {
    return "Je comprends que vous souhaitez parler à un conseiller humain. Utilisez le bouton '👤 Humain' en haut de l'écran pour être mis en relation avec un de nos conseillers.\n\n⏱️ Temps d'attente estimé : 2-5 minutes."
  }
  if (lower.includes('merci') || lower.includes('super') || lower.includes('parfait')) {
    return "Je vous en prie ! 😊 Ravi d'avoir pu vous aider. N'hésitez pas si vous avez d'autres questions.\n\nBonne journée et à bientôt sur Zekora !"
  }
  
  return "Je comprends votre demande. Pourriez-vous me donner plus de détails ?\n\nJe peux vous aider avec :\n• 💳 Compte et cartes\n• 💸 Transferts\n• ₿ Cryptomonnaies\n• 📊 Frais\n• 🔐 Sécurité\n\nOu utilisez le bouton '👤 Humain' pour une assistance personnalisée."
}

const escalateToHuman = async () => {
  isTyping.value = true
  
  try {
    // Call backend API to escalate
    if (conversation.value.id && !conversation.value.id.startsWith('demo-')) {
      await supportAPI.escalate(conversation.value.id, 'User requested human agent')
    }
  } catch (error) {
    console.error('Error escalating:', error)
  }
  
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  messages.value.push({
    id: 'msg-' + Date.now(),
    sender_type: 'system',
    sender_name: 'Système',
    content: "🔔 Votre demande a été transférée à un conseiller humain. Un agent va prendre en charge votre conversation sous peu.\n\n⏱️ Temps d'attente estimé : 2-5 minutes.",
    created_at: new Date().toISOString()
  })
  
  conversation.value.status = 'escalated'
  conversation.value.agent_type = 'human'
  agentName.value = 'Conseiller Support'
  quickReplies.value = []
  isTyping.value = false
  scrollToBottom()
}


const closeConversation = () => {
  showRatingModal.value = true
}

const submitRating = async () => {
  try {
    // Call backend API to close with rating
    if (conversation.value.id && !conversation.value.id.startsWith('demo-')) {
      await supportAPI.closeTicket(conversation.value.id, rating.value, feedback.value)
    }
  } catch (error) {
    console.error('Error submitting rating:', error)
  }
  showRatingModal.value = false
  navigateToSupport()
}

const navigateToSupport = () => {
  router.push('/support')
}


const formatTime = (dateString) => {
  return new Date(dateString).toLocaleTimeString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDateHeader = (date) => {
  const today = new Date()
  if (date.toDateString() === today.toDateString()) {
    return "Aujourd'hui"
  }
  return date.toLocaleDateString('fr-FR', { weekday: 'long', day: 'numeric', month: 'long' })
}

const loadConversation = async () => {
  const id = route.query.id
  const agentType = route.query.agent_type || 'ai'
  
  if (!id) {
    router.push('/support')
    return
  }

  conversation.value.id = id
  conversation.value.agent_type = agentType
  
  // Set agent name based on type
  if (agentType === 'human') {
    agentName.value = 'Conseiller Support'
    quickReplies.value = [] // No quick replies for human agent
  }

  // Try to load existing messages from backend
  if (!id.startsWith('demo-')) {
    try {
      const response = await supportAPI.getMessages(id)
      if (response.data?.messages && response.data.messages.length > 0) {
        messages.value = response.data.messages.map(msg => ({
          id: msg.id,
          sender_type: msg.sender_type || (msg.is_agent ? 'agent' : 'user'),
          sender_name: msg.sender_name || (msg.is_agent ? 'Support' : 'Vous'),
          content: msg.content || msg.message,
          created_at: msg.created_at,
          is_read: msg.is_read
        }))
        scrollToBottom()
        return
      }
    } catch (error) {
      console.error('Error loading messages:', error)
    }
  }

  // Load initial welcome message based on agent type
  if (agentType === 'human') {
    messages.value = [
      {
        id: 'welcome',
        sender_type: 'system',
        sender_name: 'Système',
        content: "👤 Vous avez demandé à parler à un conseiller humain.\n\n⏱️ Un de nos conseillers va prendre en charge votre conversation sous peu.\n\nTemps d'attente estimé : 2-5 minutes.\n\nMerci de patienter.",
        created_at: new Date().toISOString()
      }
    ]
  } else {
    messages.value = [
      {
        id: 'welcome',
        sender_type: 'agent',
        sender_name: 'Assistant IA',
        content: "Bonjour ! 👋 Je suis l'assistant virtuel Zekora. Je suis là pour vous aider 24/7.\n\nVoici ce que je peux faire pour vous :\n• 💳 Assistance cartes bancaires\n• 💸 Aide aux transferts\n• ₿ Questions sur les cryptomonnaies\n• 📊 Informations sur les frais\n• 🔐 Sécurité du compte\n\nComment puis-je vous aider ?",
        created_at: new Date().toISOString()
      }
    ]
  }
  scrollToBottom()
}

watch(() => messages.value.length, scrollToBottom)

// Polling for new messages every 5 seconds (for human agent responses)
let pollingInterval = null

const startPolling = () => {
  if (pollingInterval) return
  
  pollingInterval = setInterval(async () => {
    // Only poll if we have a valid conversation with a human agent
    if (!conversation.value.id || conversation.value.id.startsWith('demo-')) return
    if (conversation.value.agent_type !== 'human') return
    if (conversation.value.status === 'closed' || conversation.value.status === 'resolved') return
    
    try {
      const response = await supportAPI.getMessages(conversation.value.id)
      if (response.data?.messages) {
        const newMessages = response.data.messages
        // Only update if we have more messages
        if (newMessages.length > messages.value.length) {
          const currentIds = new Set(messages.value.map((m: any) => m.id))
          const newOnes = newMessages.filter((m: any) => !currentIds.has(m.id))
          
          if (newOnes.length > 0) {
            // Add only truly new messages
            for (const msg of newOnes) {
              messages.value.push({
                id: msg.id,
                sender_type: msg.sender_type || (msg.is_agent ? 'agent' : 'user'),
                sender_name: msg.sender_name || (msg.is_agent ? 'Support' : 'Vous'),
                content: msg.content || msg.message,
                created_at: msg.created_at,
                is_read: msg.is_read
              })
            }
            scrollToBottom()
          }
        }
      }
    } catch (error) {
      console.error('Error polling messages:', error)
    }
  }, 5000) // Poll every 5 seconds
}

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

onMounted(() => {
  loadConversation()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.chat-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 8rem);
  max-height: calc(100vh - 8rem);
  gap: 1rem;
}

@media (max-width: 1024px) {
  .chat-page {
    height: calc(100vh - 5rem);
    max-height: calc(100vh - 5rem);
  }
}

/* Chat Header */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1rem;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.back-btn {
  padding: 0.5rem;
  border-radius: 0.5rem;
  color: #64748b;
  transition: all 0.2s;
}

.back-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1e293b;
}

.dark .back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.agent-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.agent-avatar-wrapper {
  position: relative;
}

.agent-avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.online-indicator {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 0.75rem;
  height: 0.75rem;
  border-radius: 9999px;
  background: #22c55e;
  border: 2px solid white;
}

.dark .online-indicator {
  border-color: #0f172a;
}

.agent-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.dark .agent-name {
  color: white;
}

.agent-status {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.75rem;
  color: #22c55e;
  margin: 0;
}

.status-dot {
  width: 0.375rem;
  height: 0.375rem;
  border-radius: 9999px;
  background: #22c55e;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-escalate {
  padding: 0.5rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 0.5rem;
  border: none;
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-escalate:hover {
  background: rgba(245, 158, 11, 0.25);
}

.btn-close {
  padding: 0.5rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 0.5rem;
  border: none;
  background: rgba(100, 116, 139, 0.1);
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-close:hover {
  background: rgba(100, 116, 139, 0.2);
}

/* Messages Container */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.date-separator {
  display: flex;
  justify-content: center;
  margin-bottom: 0.5rem;
}

.date-separator span {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.6875rem;
  background: rgba(100, 116, 139, 0.1);
  color: #64748b;
}

.dark .date-separator span {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
}

.message-wrapper {
  display: flex;
}

.message-user {
  justify-content: flex-end;
}

.message-agent {
  justify-content: flex-start;
}

.message-bubble-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  max-width: 80%;
}

.user-wrapper {
  flex-direction: row-reverse;
}

.avatar-small {
  width: 2rem;
  height: 2rem;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.message-bubble {
  padding: 0.75rem 1rem;
  border-radius: 1rem;
  max-width: 100%;
}

.message-bubble p {
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.agent-bubble {
  background: rgba(100, 116, 139, 0.1);
  color: #1e293b;
  border-top-left-radius: 0.25rem;
}

.dark .agent-bubble {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.user-bubble {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border-top-right-radius: 0.25rem;
}

.message-time {
  font-size: 0.625rem;
  color: #94a3b8;
  margin: 0.25rem 0 0 0.5rem;
}

.user-time {
  text-align: right;
  margin-right: 0.5rem;
  margin-left: 0;
}

.read-indicator {
  color: #6366f1;
  margin-left: 0.25rem;
}

/* Typing Indicator */
.typing-bubble {
  padding: 1rem;
}

.typing-dots {
  display: flex;
  gap: 0.25rem;
}

.typing-dots span {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 9999px;
  background: #64748b;
  animation: bounce 1.4s infinite;
}

.dark .typing-dots span {
  background: rgba(255, 255, 255, 0.5);
}

.typing-dots span:nth-child(2) { animation-delay: 0.15s; }
.typing-dots span:nth-child(3) { animation-delay: 0.3s; }

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

/* Quick Replies */
.quick-replies {
  display: flex;
  gap: 0.5rem;
  padding: 0 0.5rem;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.quick-reply-btn {
  padding: 0.5rem 0.875rem;
  font-size: 0.75rem;
  border-radius: 9999px;
  border: 1px solid rgba(99, 102, 241, 0.3);
  background: rgba(99, 102, 241, 0.05);
  color: #6366f1;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.dark .quick-reply-btn {
  border-color: rgba(99, 102, 241, 0.3);
  background: rgba(99, 102, 241, 0.1);
  color: #818cf8;
}

.quick-reply-btn:hover {
  background: rgba(99, 102, 241, 0.15);
  border-color: rgba(99, 102, 241, 0.5);
}

/* Input Area */
.input-area {
  padding: 1rem;
  flex-shrink: 0;
}

.input-form {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
}

.message-input {
  flex: 1;
  resize: none;
  max-height: 8rem;
  min-height: 2.75rem;
}

.send-btn {
  padding: 0.75rem;
  flex-shrink: 0;
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.modal-content {
  width: 100%;
  max-width: 24rem;
  padding: 1.5rem;
}

.modal-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1e293b;
  text-align: center;
  margin: 0 0 1.25rem 0;
}

.dark .modal-title {
  color: white;
}

.rating-stars {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 1.25rem;
}

.star-btn {
  font-size: 2rem;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: transform 0.2s;
}

.star-btn:hover {
  transform: scale(1.1);
}

.feedback-input {
  margin-bottom: 1.25rem;
  resize: none;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
}

.modal-actions button {
  flex: 1;
  padding: 0.875rem;
}

/* Responsive */
@media (max-width: 640px) {
  .chat-header {
    padding: 0.75rem;
  }
  
  .agent-name {
    font-size: 0.875rem;
  }
  
  .btn-escalate, .btn-close {
    padding: 0.375rem 0.5rem;
    font-size: 0.6875rem;
  }
  
  .messages-container {
    padding: 0.75rem;
  }
  
  .message-bubble-wrapper {
    max-width: 90%;
  }
  
  .quick-replies {
    overflow-x: auto;
    flex-wrap: nowrap;
    padding-bottom: 0.5rem;
  }
  
  .input-area {
    padding: 0.75rem;
  }
}
</style>
