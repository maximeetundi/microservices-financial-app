<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
    <!-- Header -->
    <div class="bg-white/5 backdrop-blur-lg border-b border-white/10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <NuxtLink to="/dashboard" class="text-gray-900/60 hover:text-gray-900 transition">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </NuxtLink>
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Centre d'Assistance</h1>
              <p class="text-gray-900/60 text-sm">Nous sommes là pour vous aider 24/7</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Quick Stats -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-green-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <p class="text-gray-900/60 text-sm">Temps de réponse</p>
              <p class="text-gray-900 text-xl font-bold">~2 min</p>
            </div>
          </div>
        </div>
        <div class="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-blue-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-gray-900/60 text-sm">Disponibilité</p>
              <p class="text-gray-900 text-xl font-bold">24h/24 7j/7</p>
            </div>
          </div>
        </div>
        <div class="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-full bg-purple-500/20 flex items-center justify-center">
              <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
            </div>
            <div>
              <p class="text-gray-900/60 text-sm">Satisfaction</p>
              <p class="text-gray-900 text-xl font-bold">98%</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Agent Selection -->
      <div class="mb-8" v-if="!selectedAgent && !activeConversation">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Comment souhaitez-vous être aidé ?</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- AI Agent Card -->
          <div 
            @click="selectAgent('ai')"
            class="group relative bg-gradient-to-br from-blue-500/20 to-purple-500/20 backdrop-blur-lg rounded-2xl p-8 border border-white/10 cursor-pointer hover:border-blue-400/50 transition-all duration-300 hover:scale-[1.02]"
          >
            <div class="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-purple-500/10 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <div class="relative">
              <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center mb-6 text-4xl">
                🤖
              </div>
              <h3 class="text-2xl font-bold text-gray-900 mb-2">Assistant IA</h3>
              <p class="text-gray-900/70 mb-4">Réponses instantanées 24/7 pour les questions courantes. Idéal pour les demandes simples.</p>
              <ul class="space-y-2">
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Réponse instantanée
                </li>
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Disponible 24h/24
                </li>
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Escalade vers humain possible
                </li>
              </ul>
            </div>
          </div>

          <!-- Human Agent Card -->
          <div 
            @click="selectAgent('human')"
            class="group relative bg-gradient-to-br from-emerald-500/20 to-teal-500/20 backdrop-blur-lg rounded-2xl p-8 border border-white/10 cursor-pointer hover:border-emerald-400/50 transition-all duration-300 hover:scale-[1.02]"
          >
            <div class="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-teal-500/10 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <div class="relative">
              <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center mb-6 text-4xl">
                👤
              </div>
              <h3 class="text-2xl font-bold text-gray-900 mb-2">Conseiller Humain</h3>
              <p class="text-gray-900/70 mb-4">Pour les demandes complexes nécessitant une intervention ou une décision spécifique.</p>
              <ul class="space-y-2">
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Analyse personnalisée
                </li>
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Peut prendre des actions
                </li>
                <li class="flex items-center gap-2 text-gray-900/80">
                  <svg class="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                  </svg>
                  Attente : 2-5 min
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Start Conversation Form -->
      <div v-if="selectedAgent && !activeConversation" class="max-w-2xl mx-auto">
        <div class="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10">
          <div class="flex items-center gap-4 mb-6">
            <button @click="selectedAgent = null" class="text-gray-900/60 hover:text-gray-900">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" :class="selectedAgent === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
                {{ selectedAgent === 'ai' ? '🤖' : '👤' }}
              </div>
              <div>
                <h3 class="text-gray-900 font-semibold">{{ selectedAgent === 'ai' ? 'Assistant IA' : 'Conseiller Humain' }}</h3>
                <p class="text-gray-900/60 text-sm">{{ selectedAgent === 'ai' ? 'Réponse instantanée' : 'Temps d\'attente: ~3 min' }}</p>
              </div>
            </div>
          </div>

          <form @submit.prevent="startConversation" class="space-y-6">
            <div>
              <label class="block text-gray-900/80 mb-2 text-sm font-medium">Sujet</label>
              <input 
                v-model="newConversation.subject" 
                type="text" 
                placeholder="Ex: Problème avec mon transfert"
                class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-900 placeholder-white/40 focus:outline-none focus:border-blue-500 transition"
                required
              />
            </div>

            <div>
              <label class="block text-gray-900/80 mb-2 text-sm font-medium">Catégorie</label>
              <select 
                v-model="newConversation.category"
                class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-900 focus:outline-none focus:border-blue-500 transition"
                required
              >
                <option value="" disabled class="text-gray-900">Sélectionner une catégorie</option>
                <option value="account" class="text-gray-900">💳 Compte & Cartes</option>
                <option value="transfer" class="text-gray-900">💸 Transferts & Virements</option>
                <option value="crypto" class="text-gray-900">₿ Cryptomonnaies</option>
                <option value="security" class="text-gray-900">🔐 Sécurité</option>
                <option value="fees" class="text-gray-900">📊 Frais & Facturation</option>
                <option value="other" class="text-gray-900">❓ Autre</option>
              </select>
            </div>

            <div>
              <label class="block text-gray-900/80 mb-2 text-sm font-medium">Décrivez votre demande</label>
              <textarea 
                v-model="newConversation.message"
                rows="4"
                placeholder="Expliquez votre problème ou question..."
                class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-900 placeholder-white/40 focus:outline-none focus:border-blue-500 transition resize-none"
                required
              ></textarea>
            </div>

            <button 
              type="submit"
              :disabled="loading"
              class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-gray-900 font-semibold rounded-xl hover:from-blue-600 hover:to-purple-700 transition-all duration-300 disabled:opacity-50"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Connexion en cours...
              </span>
              <span v-else>Démarrer la conversation</span>
            </button>
          </form>
        </div>
      </div>

      <!-- Existing Conversations -->
      <div v-if="!selectedAgent && !activeConversation && conversations.length > 0" class="mt-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Vos conversations</h2>
        <div class="space-y-4">
          <div 
            v-for="conv in conversations" 
            :key="conv.id"
            @click="openConversation(conv)"
            class="bg-white/5 backdrop-blur-lg rounded-xl p-4 border border-white/10 cursor-pointer hover:bg-white/10 transition"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-full flex items-center justify-center text-xl" :class="conv.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
                  {{ conv.agent_type === 'ai' ? '🤖' : '👤' }}
                </div>
                <div>
                  <h3 class="text-gray-900 font-medium">{{ conv.subject }}</h3>
                  <p class="text-gray-900/60 text-sm truncate max-w-md">{{ conv.last_message }}</p>
                </div>
              </div>
              <div class="text-right">
                <span :class="getStatusClass(conv.status)" class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getStatusLabel(conv.status) }}
                </span>
                <p class="text-gray-900/40 text-xs mt-1">{{ formatDate(conv.updated_at) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { supportAPI } from '~/composables/useApi'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const router = useRouter()

const selectedAgent = ref(null)
const activeConversation = ref(null)
const conversations = ref([])
const loading = ref(false)

const newConversation = ref({
  subject: '',
  category: '',
  message: ''
})

const selectAgent = (type) => {
  selectedAgent.value = type
}

const startConversation = async () => {
  loading.value = true
  
  try {
    const response = await supportAPI.createTicket({
      subject: newConversation.value.subject,
      category: newConversation.value.category,
      description: newConversation.value.message,
      priority: 'normal'
    })
    
    // Navigate to chat with ticket ID
    const ticketId = response.data?.ticket?.id || response.data?.id || 'demo-' + Date.now()
    router.push(`/support/chat?id=${ticketId}`)
  } catch (error) {
    console.error('Error starting conversation:', error)
    // For demo, navigate anyway
    router.push('/support/chat?id=demo-' + Date.now())
  } finally {
    loading.value = false
  }
}

const openConversation = (conv) => {
  router.push(`/support/chat?id=${conv.id}`)
}

const fetchConversations = async () => {
  try {
    const response = await supportAPI.getTickets()
    // Map tickets to conversation format for UI
    conversations.value = (response.data?.tickets || []).map(ticket => ({
      id: ticket.id,
      subject: ticket.subject,
      agent_type: 'ai', // Default to AI
      status: ticket.status,
      last_message: ticket.description,
      updated_at: ticket.updated_at || ticket.created_at
    }))
  } catch (error) {
    // Demo data fallback
    conversations.value = [
      {
        id: '1',
        subject: 'Question sur mon transfert',
        agent_type: 'ai',
        status: 'resolved',
        last_message: 'Votre transfert a été traité avec succès.',
        updated_at: new Date().toISOString()
      }
    ]
  }
}

const getStatusClass = (status) => {
  const classes = {
    'open': 'bg-blue-500/20 text-blue-400',
    'pending': 'bg-yellow-500/20 text-yellow-400',
    'active': 'bg-green-500/20 text-green-400',
    'resolved': 'bg-emerald-500/20 text-emerald-400',
    'closed': 'bg-gray-500/20 text-gray-400',
    'escalated': 'bg-orange-500/20 text-orange-400'
  }
  return classes[status] || 'bg-gray-500/20 text-gray-400'
}

const getStatusLabel = (status) => {
  const labels = {
    'open': 'En attente',
    'pending': 'En cours',
    'active': 'Actif',
    'resolved': 'Résolu',
    'closed': 'Fermé',
    'escalated': 'Escaladé'
  }
  return labels[status] || status
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchConversations()
})
</script>
