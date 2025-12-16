<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex flex-col">
    <!-- Header -->
    <div class="bg-white/5 backdrop-blur-lg border-b border-white/10 flex-shrink-0">
      <div class="max-w-5xl mx-auto px-4 py-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <NuxtLink to="/support" class="text-gray-900/60 hover:text-gray-900 transition">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </NuxtLink>
            <div class="flex items-center gap-3">
              <div class="relative">
                <div class="w-10 h-10 rounded-full flex items-center justify-center text-xl" :class="conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
                  {{ conversation.agent_type === 'ai' ? '🤖' : '👤' }}
                </div>
                <div class="absolute -bottom-1 -right-1 w-3 h-3 rounded-full bg-green-500 border-2 border-slate-900"></div>
              </div>
              <div>
                <h1 class="text-gray-900 font-semibold text-sm">{{ conversation.agent_type === 'ai' ? 'Assistant IA' : agentName }}</h1>
                <p class="text-green-400 text-xs flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-green-400 animate-pulse"></span>
                  En ligne
                </p>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button 
              v-if="conversation.agent_type === 'ai' && conversation.status !== 'escalated'"
              @click="escalateToHuman"
              class="px-3 py-1.5 bg-orange-500/20 text-orange-400 text-sm rounded-lg hover:bg-orange-500/30 transition"
            >
              👤 Parler à un humain
            </button>
            <button 
              @click="closeConversation"
              class="px-3 py-1.5 bg-white/10 text-gray-900/60 text-sm rounded-lg hover:bg-white/20 transition"
            >
              Fermer
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Messages Area -->
    <div ref="messagesContainer" class="flex-1 overflow-y-auto px-4 py-6">
      <div class="max-w-3xl mx-auto space-y-4">
        <!-- Date Separator -->
        <div class="flex items-center justify-center">
          <div class="bg-white/10 px-3 py-1 rounded-full">
            <span class="text-gray-900/50 text-xs">{{ formatDateHeader(new Date()) }}</span>
          </div>
        </div>

        <!-- Messages -->
        <div 
          v-for="(message, index) in messages" 
          :key="message.id"
          class="flex" 
          :class="message.sender_type === 'user' ? 'justify-end' : 'justify-start'"
        >
          <!-- Agent/AI Message -->
          <div v-if="message.sender_type !== 'user'" class="flex items-start gap-3 max-w-[80%]">
            <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm flex-shrink-0" :class="message.sender_type === 'system' ? 'bg-gray-500/20' : (conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20')">
              {{ message.sender_type === 'system' ? '🔔' : (conversation.agent_type === 'ai' ? '🤖' : '👤') }}
            </div>
            <div>
              <div class="bg-white/10 backdrop-blur-lg rounded-2xl rounded-tl-sm px-4 py-3 border border-white/5">
                <p class="text-gray-900 whitespace-pre-wrap">{{ message.content }}</p>
              </div>
              <p class="text-gray-900/40 text-xs mt-1 ml-1">{{ formatTime(message.created_at) }}</p>
            </div>
          </div>

          <!-- User Message -->
          <div v-else class="flex items-start gap-3 max-w-[80%]">
            <div class="order-2">
              <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl rounded-tr-sm px-4 py-3">
                <p class="text-gray-900 whitespace-pre-wrap">{{ message.content }}</p>
              </div>
              <p class="text-gray-900/40 text-xs mt-1 mr-1 text-right">
                {{ formatTime(message.created_at) }}
                <span v-if="message.is_read" class="text-blue-400 ml-1">✓✓</span>
              </p>
            </div>
          </div>
        </div>

        <!-- Typing Indicator -->
        <div v-if="isTyping" class="flex items-start gap-3">
          <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm" :class="conversation.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
            {{ conversation.agent_type === 'ai' ? '🤖' : '👤' }}
          </div>
          <div class="bg-white/10 backdrop-blur-lg rounded-2xl rounded-tl-sm px-4 py-3 border border-white/5">
            <div class="flex gap-1">
              <span class="w-2 h-2 bg-white/50 rounded-full animate-bounce" style="animation-delay: 0ms"></span>
              <span class="w-2 h-2 bg-white/50 rounded-full animate-bounce" style="animation-delay: 150ms"></span>
              <span class="w-2 h-2 bg-white/50 rounded-full animate-bounce" style="animation-delay: 300ms"></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Replies -->
    <div v-if="quickReplies.length > 0" class="px-4 pb-2 flex-shrink-0">
      <div class="max-w-3xl mx-auto">
        <div class="flex flex-wrap gap-2">
          <button 
            v-for="reply in quickReplies" 
            :key="reply"
            @click="sendQuickReply(reply)"
            class="px-3 py-1.5 bg-white/10 text-gray-900/80 text-sm rounded-full hover:bg-white/20 transition border border-white/10"
          >
            {{ reply }}
          </button>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="bg-white/5 backdrop-blur-lg border-t border-white/10 p-4 flex-shrink-0">
      <div class="max-w-3xl mx-auto">
        <form @submit.prevent="sendMessage" class="flex items-end gap-3">
          <div class="flex-1 relative">
            <textarea
              ref="messageInput"
              v-model="newMessage"
              @keydown.enter.exact.prevent="sendMessage"
              placeholder="Écrivez votre message..."
              rows="1"
              class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-900 placeholder-white/40 focus:outline-none focus:border-blue-500 transition resize-none max-h-32"
              :disabled="sending"
            ></textarea>
          </div>
          <button 
            type="submit"
            :disabled="!newMessage.trim() || sending"
            class="p-3 bg-gradient-to-r from-blue-500 to-purple-600 text-gray-900 rounded-xl hover:from-blue-600 hover:to-purple-700 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
            </svg>
          </button>
        </form>
      </div>
    </div>

    <!-- Satisfaction Modal -->
    <div v-if="showRatingModal" class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
      <div class="bg-slate-800 rounded-2xl p-6 max-w-md w-full mx-4 border border-white/10">
        <h3 class="text-xl font-bold text-gray-900 mb-4 text-center">Comment évaluez-vous cette conversation ?</h3>
        <div class="flex justify-center gap-2 mb-6">
          <button 
            v-for="star in 5" 
            :key="star"
            @click="rating = star"
            class="text-3xl transition-transform hover:scale-110"
          >
            {{ star <= rating ? '⭐' : '☆' }}
          </button>
        </div>
        <textarea
          v-model="feedback"
          placeholder="Commentaire (optionnel)"
          class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-900 placeholder-white/40 focus:outline-none focus:border-blue-500 transition resize-none mb-4"
          rows="3"
        ></textarea>
        <div class="flex gap-3">
          <button 
            @click="submitRating"
            class="flex-1 py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-gray-900 rounded-xl"
          >
            Envoyer
          </button>
          <button 
            @click="showRatingModal = false; navigateToSupport()"
            class="flex-1 py-3 bg-white/10 text-gray-900 rounded-xl"
          >
            Passer
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

definePageMeta({
  layout: 'default',
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

const agentName = ref('Support CryptoBank')
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
    // Simulate API call
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
    return "Pour commander votre carte CryptoBank :\n\n1. Allez dans le menu 'Cartes'\n2. Cliquez sur 'Commander une carte'\n3. Choisissez entre virtuelle (gratuite) ou physique (9.99€)\n4. Suivez les étapes de personnalisation\n\nVotre carte virtuelle est disponible instantanément !"
  }
  if (lower.includes('humain') || lower.includes('agent')) {
    return "Je comprends que vous souhaitez parler à un conseiller humain. Utilisez le bouton '👤 Parler à un humain' en haut de l'écran pour être mis en relation avec un de nos conseillers.\n\n⏱️ Temps d'attente estimé : 2-5 minutes."
  }
  if (lower.includes('merci') || lower.includes('super') || lower.includes('parfait')) {
    return "Je vous en prie ! 😊 Ravi d'avoir pu vous aider. N'hésitez pas si vous avez d'autres questions.\n\nBonne journée et à bientôt sur CryptoBank !"
  }
  
  return "Je comprends votre demande. Pourriez-vous me donner plus de détails ?\n\nJe peux vous aider avec :\n• 💳 Compte et cartes\n• 💸 Transferts\n• ₿ Cryptomonnaies\n• 📊 Frais\n• 🔐 Sécurité\n\nOu utilisez le bouton '👤 Parler à un humain' pour une assistance personnalisée."
}

const escalateToHuman = async () => {
  isTyping.value = true
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
  if (!id) {
    router.push('/support')
    return
  }

  conversation.value.id = id

  // Load initial welcome message
  messages.value = [
    {
      id: 'welcome',
      sender_type: 'agent',
      sender_name: 'Assistant IA',
      content: "Bonjour ! 👋 Je suis l'assistant virtuel CryptoBank. Je suis là pour vous aider 24/7.\n\nVoici ce que je peux faire pour vous :\n• 💳 Assistance cartes bancaires\n• 💸 Aide aux transferts\n• ₿ Questions sur les cryptomonnaies\n• 📊 Informations sur les frais\n• 🔐 Sécurité du compte\n\nComment puis-je vous aider ?",
      created_at: new Date().toISOString()
    }
  ]
  scrollToBottom()
}

watch(() => messages.value.length, scrollToBottom)

onMounted(() => {
  loadConversation()
})
</script>

<style scoped>
.animate-bounce {
  animation: bounce 1.4s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
</style>
