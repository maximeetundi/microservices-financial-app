<template>
  <NuxtLayout name="dashboard">
    <div class="events-page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">🎫 Mes Événements</h1>
          <p class="page-subtitle">Créez et gérez vos événements avec billetterie</p>
        </div>
        <NuxtLink to="/events/create" class="btn-create">
          <span class="icon">+</span>
          Créer un événement
        </NuxtLink>
      </div>

      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">🎪</div>
          <div class="stat-info">
            <span class="stat-value">{{ events.length }}</span>
            <span class="stat-label">Événements</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">🎟️</div>
          <div class="stat-info">
            <span class="stat-value">{{ totalTicketsSold }}</span>
            <span class="stat-label">Tickets vendus</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">💰</div>
          <div class="stat-info">
            <span class="stat-value">{{ formatAmount(totalRevenue) }}</span>
            <span class="stat-label">Revenus</span>
          </div>
        </div>
      </div>

      <!-- Events List -->
      <div class="events-section">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Chargement des événements...</p>
        </div>

        <div v-else-if="events.length === 0" class="empty-state">
          <div class="empty-icon">🎭</div>
          <h3>Aucun événement</h3>
          <p>Créez votre premier événement pour commencer à vendre des tickets</p>
          <NuxtLink to="/events/create" class="btn-primary">Créer un événement</NuxtLink>
        </div>

        <div v-else class="events-grid">
          <div 
            v-for="event in events" 
            :key="event.id" 
            class="event-card"
            @click="goToEvent(event.id)"
          >
            <div class="event-image">
              <img v-if="event.cover_image" :src="event.cover_image" :alt="event.title" />
              <div v-else class="placeholder-image">🎪</div>
              <span :class="['status-badge', event.status]">{{ getStatusLabel(event.status) }}</span>
            </div>
            <div class="event-content">
              <h3 class="event-title">{{ event.title }}</h3>
              <div class="event-meta">
                <span class="meta-item">📍 {{ event.location || 'Non défini' }}</span>
                <span class="meta-item">📅 {{ formatDate(event.start_date) }}</span>
              </div>
              <div class="event-stats">
                <div class="stat">
                  <span class="stat-num">{{ event.total_sold || 0 }}</span>
                  <span class="stat-text">vendus</span>
                </div>
                <div class="stat">
                  <span class="stat-num">{{ formatAmount(event.total_revenue || 0) }}</span>
                  <span class="stat-text">revenus</span>
                </div>
              </div>
              <div class="event-tiers">
                <span 
                  v-for="tier in event.ticket_tiers?.slice(0, 3)" 
                  :key="tier.id"
                  class="tier-badge"
                  :style="{ background: tier.color || '#6366f1' }"
                >
                  {{ tier.icon }} {{ tier.name }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- My Purchased Tickets Section -->
      <div class="tickets-section">
        <h2 class="section-title">🎟️ Mes Tickets Achetés</h2>
        
        <div v-if="myTickets.length === 0" class="empty-tickets">
          <p>Vous n'avez pas encore de tickets</p>
          <NuxtLink to="/events/browse" class="btn-secondary">Découvrir des événements</NuxtLink>
        </div>

        <div v-else class="tickets-grid">
          <div v-for="ticket in myTickets" :key="ticket.id" class="ticket-card" @click="showTicket(ticket)">
            <div class="ticket-left">
              <span class="ticket-icon">{{ ticket.tier_icon }}</span>
              <div class="ticket-info">
                <h4>{{ ticket.event_title }}</h4>
                <p class="ticket-tier">{{ ticket.tier_name }}</p>
                <p class="ticket-date">{{ formatDate(ticket.event_date) }}</p>
              </div>
            </div>
            <div class="ticket-right">
              <span :class="['ticket-status', ticket.status]">
                {{ ticket.status === 'paid' ? 'Valide' : ticket.status === 'used' ? 'Utilisé' : ticket.status }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Ticket QR Modal -->
      <Teleport to="body">
        <div v-if="showQRModal" class="modal-overlay" @click.self="showQRModal = false">
          <div class="qr-modal">
            <button class="close-btn" @click="showQRModal = false">✕</button>
            <h3>{{ selectedTicket?.event_title }}</h3>
            <p class="tier-info">{{ selectedTicket?.tier_icon }} {{ selectedTicket?.tier_name }}</p>
            <img :src="selectedTicket?.qr_code" alt="QR Code" class="qr-image" />
            <p class="ticket-code">{{ selectedTicket?.ticket_code }}</p>
            <p class="status-info" :class="selectedTicket?.status">
              {{ selectedTicket?.status === 'paid' ? '✅ Ticket valide' : '⚠️ Ticket déjà utilisé' }}
            </p>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ticketAPI } from '~/composables/useApi'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const router = useRouter()
const loading = ref(true)
const events = ref([])
const myTickets = ref([])
const showQRModal = ref(false)
const selectedTicket = ref(null)

const totalTicketsSold = computed(() => events.value.reduce((sum, e) => sum + (e.total_sold || 0), 0))
const totalRevenue = computed(() => events.value.reduce((sum, e) => sum + (e.total_revenue || 0), 0))

const loadData = async () => {
  loading.value = true
  try {
    const [eventsRes, ticketsRes] = await Promise.all([
      ticketAPI.getMyEvents(),
      ticketAPI.getMyTickets()
    ])
    events.value = eventsRes.data?.events || []
    myTickets.value = ticketsRes.data?.tickets || []
  } catch (e) {
    console.error('Failed to load data:', e)
  } finally {
    loading.value = false
  }
}

const goToEvent = (id) => {
  router.push(`/events/${id}`)
}

const showTicket = (ticket) => {
  selectedTicket.value = ticket
  showQRModal.value = true
}

const formatDate = (date) => {
  if (!date) return 'Non défini'
  return new Date(date).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('fr-FR').format(amount) + ' XOF'
}

const getStatusLabel = (status) => {
  const labels = {
    draft: 'Brouillon',
    active: 'Actif',
    ended: 'Terminé',
    cancelled: 'Annulé'
  }
  return labels[status] || status
}

onMounted(loadData)
</script>

<style scoped>
.events-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.page-subtitle {
  color: var(--text-muted);
  margin-top: 4px;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  padding: 12px 24px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.btn-create .icon {
  font-size: 20px;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 32px;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 12px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  color: var(--text-muted);
  font-size: 14px;
}

/* Events Grid */
.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.event-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.event-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.event-image {
  height: 160px;
  position: relative;
  background: linear-gradient(135deg, #1e1e2e 0%, #2d2d3a 100%);
}

.event-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder-image {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
}

.status-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.draft { background: #6b7280; color: white; }
.status-badge.active { background: #10b981; color: white; }
.status-badge.ended { background: #f59e0b; color: white; }
.status-badge.cancelled { background: #ef4444; color: white; }

.event-content {
  padding: 16px;
}

.event-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 12px 0;
}

.event-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 12px;
}

.meta-item {
  font-size: 13px;
  color: var(--text-muted);
}

.event-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 12px;
}

.event-stats .stat {
  display: flex;
  flex-direction: column;
}

.stat-num {
  font-weight: 700;
  color: var(--text-primary);
}

.stat-text {
  font-size: 12px;
  color: var(--text-muted);
}

.event-tiers {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tier-badge {
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 12px;
  color: white;
}

/* Tickets Section */
.tickets-section {
  margin-top: 48px;
}

.section-title {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.tickets-grid {
  display: grid;
  gap: 12px;
}

.ticket-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.ticket-card:hover {
  background: var(--surface-hover);
}

.ticket-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ticket-icon {
  font-size: 32px;
}

.ticket-info h4 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.ticket-tier {
  font-size: 14px;
  color: var(--text-muted);
  margin: 2px 0;
}

.ticket-date {
  font-size: 12px;
  color: var(--text-muted);
}

.ticket-status {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.ticket-status.paid { background: #10b981; color: white; }
.ticket-status.used { background: #6b7280; color: white; }

/* Empty & Loading States */
.empty-state, .loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.btn-primary, .btn-secondary {
  display: inline-block;
  padding: 12px 24px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  margin-top: 16px;
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
}

.btn-secondary {
  background: var(--surface);
  border: 1px solid var(--border);
  color: var(--text-primary);
}

/* QR Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.qr-modal {
  background: var(--surface);
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  max-width: 360px;
  position: relative;
}

.close-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 20px;
  cursor: pointer;
}

.qr-modal h3 {
  margin: 0 0 8px;
  color: var(--text-primary);
}

.tier-info {
  color: var(--text-muted);
  margin-bottom: 20px;
}

.qr-image {
  width: 200px;
  height: 200px;
  border-radius: 12px;
  margin-bottom: 16px;
}

.ticket-code {
  font-family: monospace;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  background: var(--surface-hover);
  padding: 8px 16px;
  border-radius: 8px;
  display: inline-block;
  margin-bottom: 12px;
}

.status-info {
  font-weight: 600;
}

.status-info.paid { color: #10b981; }
.status-info.used { color: #6b7280; }

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .events-grid {
    grid-template-columns: 1fr;
  }
}
</style>
