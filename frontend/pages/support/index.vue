<template>
  <NuxtLayout name="dashboard">
    <div class="support-page">
      <!-- Header -->
      <div class="page-header">
        <h1>🎧 Centre d'Assistance</h1>
        <p class="header-subtitle">Nous sommes là pour vous aider 24/7</p>
      </div>

      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card stat-card-green">
          <div class="stat-icon bg-emerald-500/20">
            <svg class="w-6 h-6 text-emerald-500 dark:text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <div class="stat-content">
            <p class="stat-label">Temps de réponse</p>
            <p class="stat-value">~2 min</p>
          </div>
        </div>
        <div class="stat-card stat-card-blue">
          <div class="stat-icon bg-blue-500/20">
            <svg class="w-6 h-6 text-blue-500 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="stat-content">
            <p class="stat-label">Disponibilité</p>
            <p class="stat-value">24h/24 7j/7</p>
          </div>
        </div>
        <div class="stat-card stat-card-purple">
          <div class="stat-icon bg-purple-500/20">
            <svg class="w-6 h-6 text-purple-500 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </div>
          <div class="stat-content">
            <p class="stat-label">Satisfaction</p>
            <p class="stat-value">98%</p>
          </div>
        </div>
      </div>

      <!-- Agent Selection -->
      <div v-if="!selectedAgent && !activeConversation" class="section">
        <h2 class="section-title">Comment souhaitez-vous être aidé ?</h2>
        <div class="agents-grid">
          <!-- AI Agent Card -->
          <div 
            @click="selectAgent('ai')"
            class="agent-card agent-card-ai"
          >
            <div class="agent-icon">🤖</div>
            <h3 class="agent-title">Assistant IA</h3>
            <p class="agent-description">Réponses instantanées 24/7 pour les questions courantes.</p>
            <ul class="agent-features">
              <li>
                <svg class="feature-check" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Réponse instantanée
              </li>
              <li>
                <svg class="feature-check" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Disponible 24h/24
              </li>
              <li>
                <svg class="feature-check" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Escalade vers humain
              </li>
            </ul>
          </div>

          <!-- Human Agent Card -->
          <div 
            @click="selectAgent('human')"
            class="agent-card agent-card-human"
          >
            <div class="agent-icon">👤</div>
            <h3 class="agent-title">Conseiller Humain</h3>
            <p class="agent-description">Pour les demandes complexes nécessitant une intervention.</p>
            <ul class="agent-features">
              <li>
                <svg class="feature-check" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Analyse personnalisée
              </li>
              <li>
                <svg class="feature-check" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Peut prendre des actions
              </li>
              <li>
                <svg class="feature-check-warning" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                </svg>
                Attente : 2-5 min
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Start Conversation Form -->
      <div v-if="selectedAgent && !activeConversation" class="section">
        <div class="glass-card form-card">
          <div class="form-header">
            <button @click="selectedAgent = null" class="back-btn">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <div class="form-header-info">
              <div class="agent-avatar" :class="selectedAgent === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
                {{ selectedAgent === 'ai' ? '🤖' : '👤' }}
              </div>
              <div>
                <h3 class="form-header-title">{{ selectedAgent === 'ai' ? 'Assistant IA' : 'Conseiller Humain' }}</h3>
                <p class="form-header-subtitle">{{ selectedAgent === 'ai' ? 'Réponse instantanée' : 'Temps d\'attente: ~3 min' }}</p>
              </div>
            </div>
          </div>

          <form @submit.prevent="startConversation" class="conversation-form">
            <div class="form-group">
              <label class="form-label">Sujet <span class="required-hint">(min. 5 caractères)</span></label>
              <input 
                v-model="newConversation.subject" 
                type="text" 
                placeholder="Ex: Problème avec mon transfert"
                class="input-premium"
                :class="{ 'input-error': newConversation.subject.length > 0 && newConversation.subject.length < 5 }"
                minlength="5"
                required
              />
              <p v-if="newConversation.subject.length > 0 && newConversation.subject.length < 5" class="error-hint">
                Le sujet doit contenir au moins 5 caractères ({{ newConversation.subject.length }}/5)
              </p>
            </div>

            <div class="form-group">
              <label class="form-label">Catégorie</label>
              <select 
                v-model="newConversation.category"
                class="input-premium"
                required
              >
                <option value="" disabled>Sélectionner une catégorie</option>
                <option value="account">💳 Compte & Cartes</option>
                <option value="transfer">💸 Transferts & Virements</option>
                <option value="crypto">₿ Cryptomonnaies</option>
                <option value="security">🔐 Sécurité</option>
                <option value="fees">📊 Frais & Facturation</option>
                <option value="other">❓ Autre</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">Décrivez votre demande</label>
              <textarea 
                v-model="newConversation.message"
                rows="4"
                placeholder="Expliquez votre problème ou question..."
                class="input-premium"
                required
              ></textarea>
            </div>

            <button 
              type="submit"
              :disabled="loading"
              class="btn-premium submit-btn"
            >
              <span v-if="loading" class="loading-content">
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
      <div v-if="!selectedAgent && !activeConversation && conversations.length > 0" class="section">
        <h2 class="section-title">Vos conversations</h2>
        <div class="conversations-list">
          <div 
            v-for="conv in conversations" 
            :key="conv.id"
            @click="openConversation(conv)"
            class="conversation-item"
          >
            <div class="conversation-avatar" :class="conv.agent_type === 'ai' ? 'bg-blue-500/20' : 'bg-emerald-500/20'">
              {{ conv.agent_type === 'ai' ? '🤖' : '👤' }}
            </div>
            <div class="conversation-info">
              <h3 class="conversation-subject">{{ conv.subject }}</h3>
              <p class="conversation-preview">{{ conv.last_message }}</p>
            </div>
            <div class="conversation-meta">
              <span :class="conv.agent_type === 'ai' ? 'agent-badge-ai' : 'agent-badge-human'" class="agent-type-badge">
                {{ conv.agent_type === 'ai' ? '🤖 IA' : '👤 Humain' }}
              </span>
              <span :class="getStatusClass(conv.status)" class="status-badge">
                {{ getStatusLabel(conv.status) }}
              </span>
              <p class="conversation-time">{{ formatDate(conv.updated_at) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!selectedAgent && !activeConversation && conversations.length === 0 && !loading" class="empty-conversations">
        <div class="glass-card empty-card">
          <span class="empty-icon">💬</span>
          <h3 class="empty-title">Aucune conversation</h3>
          <p class="empty-text">Sélectionnez un type d'assistance ci-dessus pour commencer.</p>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { supportAPI } from '~/composables/useApi'

definePageMeta({
  layout: false,
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
  // Validate subject length
  if (newConversation.value.subject.length < 5) {
    alert('Le sujet doit contenir au moins 5 caractères.')
    return
  }
  
  loading.value = true
  
  try {
    const response = await supportAPI.createTicket({
      subject: newConversation.value.subject,
      category: newConversation.value.category,
      description: newConversation.value.message,
      priority: 'normal',
      agent_type: selectedAgent.value
    })
    
    // Navigate to chat with ticket ID and agent_type
    const ticketId = response.data?.ticket?.id || response.data?.conversation?.id || response.data?.id
    if (ticketId) {
      router.push(`/support/chat?id=${ticketId}&agent_type=${selectedAgent.value}`)
    } else {
      console.error('No ticket ID returned from API')
      alert('Erreur: Impossible de créer la conversation. Veuillez réessayer.')
    }
  } catch (error) {
    console.error('Error starting conversation:', error)
    alert('Erreur de connexion au service support. Veuillez réessayer.')
  } finally {
    loading.value = false
  }
}

const openConversation = (conv) => {
  // Navigate to chat with conversation ID and agent_type
  router.push(`/support/chat?id=${conv.id}&agent_type=${conv.agent_type || 'ai'}`)
}

const fetchConversations = async () => {
  loading.value = true
  try {
    const response = await supportAPI.getTickets()
    // Map conversations from backend response
    conversations.value = (response.data?.conversations || []).map(conv => ({
      id: conv.id,
      subject: conv.subject,
      agent_type: conv.agent_type || 'ai',
      status: conv.status,
      last_message: conv.last_message || conv.subject,
      updated_at: conv.updated_at || conv.created_at
    }))
  } catch (error) {
    console.error('Error fetching conversations:', error)
    conversations.value = []
  } finally {
    loading.value = false
  }
}


const getStatusClass = (status) => {
  const classes = {
    'open': 'badge-info',
    'pending': 'badge-warning',
    'active': 'badge-success',
    'resolved': 'badge-success',
    'closed': 'badge-muted',
    'escalated': 'badge-warning'
  }
  return classes[status] || 'badge-muted'
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

// Polling for conversations list every 15 seconds
let conversationsPollingInterval = null

const startConversationsPolling = () => {
  if (conversationsPollingInterval) return
  
  conversationsPollingInterval = setInterval(async () => {
    // Only poll if not in active conversation and not selecting agent
    if (selectedAgent.value || activeConversation.value) return
    
    try {
      const response = await supportAPI.getTickets()
      if (response.data?.conversations) {
        conversations.value = response.data.conversations.map(conv => ({
          id: conv.id,
          subject: conv.subject,
          agent_type: conv.agent_type || 'ai',
          status: conv.status,
          last_message: conv.last_message || conv.subject,
          updated_at: conv.updated_at || conv.created_at
        }))
      }
    } catch (error) {
      console.error('Error polling conversations:', error)
    }
  }, 15000) // Poll every 15 seconds
}

const stopConversationsPolling = () => {
  if (conversationsPollingInterval) {
    clearInterval(conversationsPollingInterval)
    conversationsPollingInterval = null
  }
}

onMounted(() => {
  fetchConversations()
  startConversationsPolling()
})

onUnmounted(() => {
  stopConversationsPolling()
})
</script>

<style scoped>
.support-page {
  width: 100%;
  max-width: 100%;
}

/* Header */
.page-header {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  padding: 1.5rem;
  border-radius: 1rem;
  margin-bottom: 1.5rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.25rem 0;
}

.header-subtitle {
  font-size: 0.875rem;
  opacity: 0.9;
  margin: 0;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 2rem;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: #64748b;
  margin: 0 0 0.25rem 0;
}

.dark .stat-label {
  color: rgba(255, 255, 255, 0.6);
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.dark .stat-value {
  color: white;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
}

/* Section */
.section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 1rem 0;
}

.dark .section-title {
  color: white;
}

/* Agents Grid */
.agents-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

@media (max-width: 768px) {
  .agents-grid {
    grid-template-columns: 1fr;
  }
}

.agent-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 1.5rem;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.dark .agent-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.agent-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(99, 102, 241, 0.15);
}

.agent-card-ai:hover {
  border-color: rgba(99, 102, 241, 0.5);
}

.agent-card-human:hover {
  border-color: rgba(16, 185, 129, 0.5);
}

.agent-icon {
  width: 4rem;
  height: 4rem;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin-bottom: 1rem;
}

.agent-card-ai .agent-icon {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(139, 92, 246, 0.2));
}

.agent-card-human .agent-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(20, 184, 166, 0.2));
}

.agent-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}

.dark .agent-title {
  color: white;
}

.agent-description {
  font-size: 0.875rem;
  color: #64748b;
  margin: 0 0 1rem 0;
}

.dark .agent-description {
  color: rgba(255, 255, 255, 0.7);
}

.agent-features {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.agent-features li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #475569;
}

.dark .agent-features li {
  color: rgba(255, 255, 255, 0.8);
}

.feature-check {
  width: 1.25rem;
  height: 1.25rem;
  color: #22c55e;
  flex-shrink: 0;
}

.feature-check-warning {
  width: 1.25rem;
  height: 1.25rem;
  color: #f59e0b;
  flex-shrink: 0;
}

/* Form Card */
.form-card {
  max-width: 36rem;
  margin: 0 auto;
  padding: 1.5rem;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.back-btn {
  padding: 0.5rem;
  border-radius: 0.5rem;
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
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

.form-header-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.agent-avatar {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.form-header-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.dark .form-header-title {
  color: white;
}

.form-header-subtitle {
  font-size: 0.75rem;
  color: #64748b;
  margin: 0;
}

.dark .form-header-subtitle {
  color: rgba(255, 255, 255, 0.6);
}

.conversation-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #475569;
}

.dark .form-label {
  color: rgba(255, 255, 255, 0.8);
}

.submit-btn {
  width: 100%;
  padding: 1rem;
  margin-top: 0.5rem;
}

.loading-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

/* Conversations List */
.conversations-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.dark .conversation-item {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.conversation-item:hover {
  background: rgba(99, 102, 241, 0.05);
  border-color: rgba(99, 102, 241, 0.2);
}

.dark .conversation-item:hover {
  background: rgba(99, 102, 241, 0.1);
}

.conversation-avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.conversation-info {
  flex: 1;
  min-width: 0;
}

.conversation-subject {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 0.25rem 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .conversation-subject {
  color: white;
}

.conversation-preview {
  font-size: 0.8125rem;
  color: #64748b;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .conversation-preview {
  color: rgba(255, 255, 255, 0.6);
}

.conversation-meta {
  text-align: right;
  flex-shrink: 0;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.625rem;
  border-radius: 9999px;
  font-size: 0.6875rem;
  font-weight: 600;
}

.badge-info {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.badge-success {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.badge-warning {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.badge-muted {
  background: rgba(100, 116, 139, 0.15);
  color: #64748b;
}

.conversation-time {
  font-size: 0.6875rem;
  color: #94a3b8;
  margin: 0.25rem 0 0 0;
}

/* Empty State */
.empty-conversations {
  margin-top: 2rem;
}

.empty-card {
  text-align: center;
  padding: 3rem 1.5rem;
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #475569;
  margin: 0 0 0.5rem 0;
}

.dark .empty-title {
  color: rgba(255, 255, 255, 0.8);
}

.empty-text {
  font-size: 0.875rem;
  color: #94a3b8;
  margin: 0;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .page-header {
    padding: 1rem;
    margin: -1rem -1rem 1rem -1rem;
    border-radius: 0;
    width: calc(100% + 2rem);
  }
  
  .page-header h1 {
    font-size: 1.25rem;
  }
  
  .stat-card {
    padding: 0.875rem;
  }
  
  .stat-icon {
    width: 2.5rem;
    height: 2.5rem;
  }
  
  .stat-value {
    font-size: 1.125rem;
  }
  
  .agent-card {
    padding: 1.25rem;
  }
  
  .agent-icon {
    width: 3rem;
    height: 3rem;
    font-size: 1.5rem;
  }
  
  .agent-title {
    font-size: 1.125rem;
  }
  
  .form-card {
    padding: 1rem;
  }
  
  .conversation-item {
    flex-wrap: wrap;
  }
  
  .conversation-meta {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 0.5rem;
    padding-left: 3.5rem;
  }
}

/* Validation styles */
.required-hint {
  font-size: 0.75rem;
  font-weight: 400;
  color: #94a3b8;
}

.input-error {
  border-color: #ef4444 !important;
  box-shadow: 0 0 0 1px rgba(239, 68, 68, 0.2) !important;
}

.error-hint {
  font-size: 0.75rem;
  color: #ef4444;
  margin: 0.25rem 0 0 0;
}

/* Agent type badges */
.agent-type-badge {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.agent-badge-ai {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.agent-badge-human {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}
</style>
