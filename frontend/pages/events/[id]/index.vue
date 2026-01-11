<template>
  <NuxtLayout name="dashboard">
    <div class="event-detail-page">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement de l'événement...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <div class="error-icon">❌</div>
        <h2>Événement non trouvé</h2>
        <NuxtLink to="/events" class="btn-back">Retour aux événements</NuxtLink>
      </div>

      <!-- Event Content -->
      <template v-else-if="event">
        <!-- Hero Section -->
        <div class="hero-section" :style="event.cover_image ? { backgroundImage: `url(${event.cover_image})` } : {}">
          <div class="hero-overlay">
            <NuxtLink to="/events" class="back-btn">← Retour</NuxtLink>
            <div class="hero-content">
              <span :class="['status-badge', event.status]">{{ getStatusLabel(event.status) }}</span>
              <h1>{{ event.title }}</h1>
              <div class="meta-info">
                <span>📍 {{ event.location || 'Non défini' }}</span>
                <span>📅 {{ formatDate(event.start_date) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="content-grid">
          <!-- Left: Event Info -->
          <div class="event-info">
            <section class="info-section">
              <h3>📝 Description</h3>
              <p>{{ event.description || 'Aucune description disponible.' }}</p>
            </section>

            <section class="info-section">
              <h3>📅 Dates importantes</h3>
              <div class="dates-grid">
                <div class="date-item">
                  <span class="date-label">Début</span>
                  <span class="date-value">{{ formatDate(event.start_date) }}</span>
                </div>
                <div class="date-item">
                  <span class="date-label">Fin</span>
                  <span class="date-value">{{ formatDate(event.end_date) }}</span>
                </div>
                <div class="date-item">
                  <span class="date-label">Ventes à partir du</span>
                  <span class="date-value">{{ formatDate(event.sale_start_date) }}</span>
                </div>
                <div class="date-item">
                  <span class="date-label">Ventes jusqu'au</span>
                  <span class="date-value">{{ formatDate(event.sale_end_date) }}</span>
                </div>
              </div>
            </section>

            <!-- QR Code for organizer -->
            <section v-if="isOwner" class="info-section qr-section">
              <h3>🔲 Code de l'événement</h3>
              <div class="qr-container">
                <div class="event-code-box" @click="copyEventCode">
                  <span class="event-code">{{ event.event_code }}</span>
                  <span class="copy-icon">📋</span>
                </div>
                <button class="btn-qr-modal" @click="showQRModal = true">
                  🔍 Voir QR Code & Télécharger
                </button>
                <p class="qr-hint">Partagez ce code pour permettre aux participants de trouver votre événement</p>
              </div>
            </section>
          </div>

          <!-- Right: Tickets -->
          <div class="tickets-section">
            <h3>🎫 Choisir un ticket</h3>
            
            <div class="tiers-list">
              <div 
                v-for="tier in event.ticket_tiers" 
                :key="tier.id" 
                :class="['tier-card', { selected: selectedTier?.id === tier.id, 'sold-out': isSoldOut(tier) }]"
                @click="!isSoldOut(tier) && selectTier(tier)"
              >
                <div class="tier-header" :style="{ borderColor: tier.color }">
                  <span class="tier-icon">{{ tier.icon }}</span>
                  <div class="tier-info">
                    <h4>{{ tier.name }}</h4>
                    <p class="tier-desc">{{ tier.description }}</p>
                  </div>
                  <div class="tier-price">
                    <span class="price">{{ formatAmount(tier.price) }}</span>
                    <span class="currency">{{ event.currency || 'XOF' }}</span>
                  </div>
                </div>
                <div class="tier-footer">
                  <span v-if="isSoldOut(tier)" class="sold-out-badge">Épuisé</span>
                  <span v-else-if="tier.quantity > 0" class="remaining">{{ tier.quantity - tier.sold }} restants</span>
                  <span v-else class="remaining">Disponible</span>
                </div>
              </div>
            </div>

            <!-- Purchase Button -->
            <button 
              v-if="selectedTier && !isOwner" 
              class="btn-purchase"
              @click="showPurchaseModal = true"
            >
              Acheter - {{ formatAmount(selectedTier.price) }} {{ event.currency || 'XOF' }}
            </button>

            <!-- Organizer Actions -->
            <div v-if="isOwner" class="organizer-actions">
              <button v-if="event.status === 'draft'" @click="publishEvent" class="btn-publish">
                🚀 Publier l'événement
              </button>
              
              <div class="actions-row">
                <NuxtLink :to="`/events/${event.id}/edit`" class="btn-edit">
                  ✏️ Modifier
                </NuxtLink>
                <button @click="confirmDelete" class="btn-delete">
                  🗑️ Supprimer
                </button>
              </div>

              <button @click="showScannerModal = true" class="btn-scan">
                📷 Scanner les tickets
              </button>
              <NuxtLink :to="`/events/${event.id}/tickets`" class="btn-secondary">
                📊 Voir les tickets vendus
              </NuxtLink>
            </div>
          </div>
        </div>
      </template>

      <!-- Purchase Modal - Step 1: Info + Wallet -->
      <Teleport to="body">
        <div v-if="showPurchaseModal" class="modal-overlay" @click.self="showPurchaseModal = false">
          <div class="purchase-modal">
            <button class="close-btn" @click="showPurchaseModal = false">✕</button>
            <h2>Acheter un ticket</h2>
            <p class="modal-subtitle">{{ selectedTier?.icon }} {{ selectedTier?.name }} - {{ formatAmount(selectedTier?.price || 0) }} {{ event?.currency }}</p>
            
            <!-- Form Fields -->
            <form @submit.prevent="openPinModal">
              <div v-for="field in event?.form_fields" :key="field.name" class="form-group">
                <label>{{ field.label }} <span v-if="field.required">*</span></label>
                
                <!-- Select Field -->
                <select 
                  v-if="field.type === 'select'"
                  v-model="formData[field.name]"
                  :required="field.required"
                >
                  <option value="">Sélectionner...</option>
                  <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                </select>

                <!-- Text/Email/Phone Field -->
                <input 
                  v-else
                  v-model="formData[field.name]" 
                  :type="field.type === 'email' ? 'email' : field.type === 'phone' ? 'tel' : 'text'"
                  :required="field.required"
                  :placeholder="field.label"
                />
              </div>

              <div class="form-group">
                <label>Portefeuille *</label>
                <select v-model="selectedWalletId" required>
                  <option value="">Sélectionner...</option>
                  <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                    {{ wallet.currency }} - {{ formatAmount(wallet.balance) }}
                  </option>
                </select>
              </div>

              <div class="purchase-summary">
                <div class="summary-row">
                  <span>🎫 {{ selectedTier?.name }}</span>
                  <span class="summary-price">{{ formatAmount(selectedTier?.price || 0) }} {{ event?.currency }}</span>
                </div>
              </div>

              <button type="submit" class="btn-confirm">
                Continuer vers le paiement →
              </button>
            </form>
          </div>
        </div>
      </Teleport>

      <!-- PIN Confirmation Modal - Step 2 -->
      <Teleport to="body">
        <div v-if="showPinModal" class="modal-overlay">
          <div class="pin-modal">
            <div class="pin-header">
              <div class="pin-icon">🔐</div>
              <h3>Confirmation du paiement</h3>
              <p class="pin-subtitle">Entrez votre code PIN pour confirmer</p>
            </div>

            <div class="pin-amount">
              <span class="pin-amount-value">{{ formatAmount(selectedTier?.price || 0) }}</span>
              <span class="pin-amount-currency">{{ event?.currency }}</span>
            </div>

            <div class="pin-input-container">
              <input 
                v-model="pin"
                type="password"
                maxlength="5"
                class="pin-input"
                placeholder="•••••"
                autofocus
                @keyup.enter="purchaseTicket"
              />
              <p class="pin-hint">Code à 5 chiffres</p>
            </div>

            <p v-if="purchaseError" class="pin-error">{{ purchaseError }}</p>

            <div class="pin-actions">
              <button @click="closePinModal" class="btn-cancel" :disabled="purchasing">
                Annuler
              </button>
              <button @click="purchaseTicket" :disabled="purchasing || pin.length !== 5" class="btn-pay">
                <span v-if="purchasing" class="spinner-small"></span>
                {{ purchasing ? 'Paiement...' : 'Payer' }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Success Modal -->
      <Teleport to="body">
        <div v-if="purchaseSuccess" class="modal-overlay">
          <div class="success-modal">
            <div class="success-icon">✅</div>
            <h2>Achat réussi !</h2>
            <p>Votre ticket a été acheté avec succès.</p>
            <img v-if="purchasedTicket?.qr_code" :src="purchasedTicket.qr_code" alt="QR" class="ticket-qr" />
            <p class="ticket-code">{{ purchasedTicket?.ticket_code }}</p>
            <NuxtLink to="/events" class="btn-primary">Voir mes tickets</NuxtLink>
          </div>
        </div>
      </Teleport>

      <!-- QR Code Modal for Organizer -->
      <Teleport to="body">
        <div v-if="showQRModal" class="modal-overlay" @click.self="showQRModal = false">
          <div class="qr-modal">
            <button class="close-btn" @click="showQRModal = false">✕</button>
            
            <!-- Large QR Code -->
            <div class="qr-large-container">
              <img v-if="event?.qr_code" :src="event.qr_code" alt="QR Code" class="qr-large" />
              <div v-else class="qr-placeholder-large">
                <span>📱</span>
                <p>QR Code non disponible</p>
              </div>
            </div>

            <!-- Event Info -->
            <div class="qr-modal-info">
              <div class="info-row">
                <span class="label">Événement</span>
                <span class="value">{{ event?.title }}</span>
              </div>
              <div class="info-row">
                <span class="label">Statut</span>
                <span :class="['status-pill', event?.status]">{{ getStatusLabel(event?.status) }}</span>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="qr-actions">
              <button class="qr-action-btn" @click="copyEventCode">
                <span class="action-icon">📋</span>
                <span>Code</span>
              </button>
              <button class="qr-action-btn" @click="downloadQRCode">
                <span class="action-icon">⬇️</span>
                <span>DL PNG</span>
              </button>
              <button class="qr-action-btn" @click="shareEvent">
                <span class="action-icon">📤</span>
                <span>Partager</span>
              </button>
            </div>

            <!-- Event Code -->
            <div class="qr-code-display">
              <span class="qr-code-label">CODE DE L'ÉVÉNEMENT</span>
              <div class="qr-code-value" @click="copyEventCode">
                <span>{{ event?.event_code }}</span>
                <span class="copy-btn">📋</span>
              </div>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Ticket Scanner Modal for Organizers -->
      <Teleport to="body">
        <div v-if="showScannerModal" class="modal-overlay" @click.self="closeScannerModal">
          <div class="scanner-modal">
            <button class="close-btn" @click="closeScannerModal">✕</button>
            <h2>📷 Scanner un ticket</h2>
            <p class="modal-subtitle">Scannez le QR code du ticket pour vérifier sa validité</p>
            
            <!-- Camera View -->
            <div class="scanner-camera-container">
              <video ref="ticketScannerVideo" class="scanner-video" playsinline></video>
              <div class="scanner-overlay">
                <div class="scanner-frame">
                  <div class="scanner-corner top-left"></div>
                  <div class="scanner-corner top-right"></div>
                  <div class="scanner-corner bottom-left"></div>
                  <div class="scanner-corner bottom-right"></div>
                  <div v-if="scannerActive" class="scanner-line"></div>
                </div>
              </div>
            </div>

            <!-- Manual Code Entry -->
            <div class="manual-entry">
              <input 
                v-model="manualTicketCode" 
                placeholder="Ou entrez le code TKT-XXXXX"
                class="ticket-code-input"
              />
              <button @click="verifyManualCode" :disabled="!manualTicketCode" class="btn-verify">
                Vérifier
              </button>
            </div>

            <!-- Scan Result -->
            <div v-if="scanResult" :class="['scan-result', scanResult.valid ? 'valid' : 'invalid']">
              <div class="result-icon">{{ scanResult.valid ? '✅' : '❌' }}</div>
              <h3>{{ scanResult.valid ? 'Ticket Valide!' : 'Ticket Invalide' }}</h3>
              <p>{{ scanResult.message }}</p>
              <div v-if="scanResult.valid && scanResult.ticket" class="ticket-info">
                <p><strong>Ticket:</strong> {{ scanResult.ticket.ticket_code }}</p>
                <p><strong>Type:</strong> {{ scanResult.ticket.tier_name }}</p>
                <p><strong>Statut:</strong> {{ scanResult.ticket.used ? 'Déjà utilisé' : 'Non utilisé' }}</p>
              </div>
              <button v-if="scanResult.valid && !scanResult.ticket?.used" @click="markTicketAsUsed" class="btn-mark-used">
                ✓ Marquer comme utilisé
              </button>
              <button @click="resetScanner" class="btn-scan-another">
                Scanner un autre ticket
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ticketAPI, walletAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'
import QrScanner from 'qr-scanner'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const eventId = route.params.id

const loading = ref(true)
const error = ref(false)
const event = ref(null)
const selectedTier = ref(null)
const showPurchaseModal = ref(false)
const purchasing = ref(false)
const purchaseSuccess = ref(false)
const purchasedTicket = ref(null)
const wallets = ref([])
const selectedWalletId = ref('')
const pin = ref('')
const formData = reactive({})
const showQRModal = ref(false)
const showPinModal = ref(false)
const purchaseError = ref('')

// Ticket Scanner state
const showScannerModal = ref(false)
const ticketScannerVideo = ref(null)
const scannerActive = ref(false)
const manualTicketCode = ref('')
const scanResult = ref(null)
const ticketScanner = ref(null)

// Auth store for user info
const authStore = useAuthStore()

// Check if current user is the event owner
const isOwner = computed(() => {
  if (!event.value) return false
  const userId = authStore.user?.id
  console.log('DEBUG isOwner:', { creatorId: event.value.creator_id, userId, isOwner: event.value.creator_id === userId })
  return event.value.creator_id === userId
})

const loadEvent = async () => {
  loading.value = true
  error.value = false
  try {
    const res = await ticketAPI.getEvent(eventId)
    event.value = res.data?.event || res.data
    
    // Initialize form data with field names
    if (event.value?.form_fields) {
      event.value.form_fields.forEach(f => {
        formData[f.name] = ''
      })
    }
  } catch (e) {
    console.error('Failed to load event:', e)
    error.value = true
  } finally {
    loading.value = false
  }
}

const loadWallets = async () => {
  try {
    const res = await walletAPI.getWallets()
    wallets.value = res.data?.wallets || []
  } catch (e) {
    console.error('Failed to load wallets:', e)
  }
}

const selectTier = (tier) => {
  selectedTier.value = tier
}

const isSoldOut = (tier) => {
  if (tier.quantity === -1) return false
  return tier.sold >= tier.quantity
}

// Open PIN confirmation modal after form is filled
const openPinModal = () => {
  if (!selectedWalletId.value) {
    alert('Veuillez sélectionner un portefeuille')
    return
  }
  showPurchaseModal.value = false
  pin.value = ''
  purchaseError.value = ''
  showPinModal.value = true
}

// Close PIN modal and go back to purchase modal
const closePinModal = () => {
  showPinModal.value = false
  pin.value = ''
  purchaseError.value = ''
  showPurchaseModal.value = true
}

const purchaseTicket = async () => {
  if (!selectedTier.value || !selectedWalletId.value || pin.value.length !== 5) return
  
  purchasing.value = true
  purchaseError.value = ''
  try {
    const res = await ticketAPI.purchaseTicket({
      event_id: eventId,
      tier_id: selectedTier.value.id,
      form_data: formData,
      wallet_id: selectedWalletId.value,
      pin: pin.value
    })
    
    purchasedTicket.value = res.data?.ticket
    showPinModal.value = false
    purchaseSuccess.value = true
  } catch (e) {
    purchaseError.value = e.response?.data?.error || 'Erreur lors du paiement. Vérifiez votre PIN et réessayez.'
  } finally {
    purchasing.value = false
  }
}

const publishEvent = async () => {
  try {
    await ticketAPI.publishEvent(eventId)
    loadEvent()
  } catch (e) {
    alert(e.response?.data?.error || 'Erreur lors de la publication')
  }
}

const confirmDelete = async () => {
  if (confirm('Êtes-vous sûr de vouloir supprimer cet événement ? Cette action est irréversible.')) {
    try {
      await ticketAPI.deleteEvent(eventId)
      router.push('/events')
    } catch (e) {
      alert(e.response?.data?.error || 'Erreur lors de la suppression')
    }
  }
}

const formatDate = (date) => {
  if (!date) return 'Non défini'
  return new Date(date).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('fr-FR').format(amount || 0)
}

const getStatusLabel = (status) => {
  const labels = { draft: 'Brouillon', active: 'Actif', ended: 'Terminé', cancelled: 'Annulé' }
  return labels[status] || status
}

const copyEventCode = async () => {
  if (!event.value?.event_code) return
  try {
    await navigator.clipboard.writeText(event.value.event_code)
    alert('Code copié: ' + event.value.event_code)
  } catch (e) {
    console.error('Failed to copy', e)
  }
}

const downloadQRCode = () => {
  if (!event.value?.qr_code) {
    alert('QR Code non disponible')
    return
  }
  
  const link = document.createElement('a')
  link.href = event.value.qr_code
  link.download = `event-${event.value.event_code}-qr.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const shareEvent = async () => {
  const shareData = {
    title: event.value?.title || 'Événement',
    text: `Rejoignez l'événement "${event.value?.title}" - Code: ${event.value?.event_code}`,
    url: window.location.href
  }
  
  try {
    if (navigator.share) {
      await navigator.share(shareData)
    } else {
      await navigator.clipboard.writeText(`${shareData.text}\n${shareData.url}`)
      alert('Lien copié dans le presse-papiers!')
    }
  } catch (e) {
    console.error('Share failed:', e)
  }
}

// ========== TICKET SCANNER FUNCTIONS ==========

// Start ticket scanner when modal opens
watch(showScannerModal, async (isOpen) => {
  if (isOpen) {
    await nextTick()
    startTicketScanner()
  } else {
    stopTicketScanner()
  }
})

const startTicketScanner = async () => {
  scanResult.value = null
  manualTicketCode.value = ''
  
  if (!ticketScannerVideo.value) {
    console.error('Video element not found')
    return
  }

  try {
    ticketScanner.value = new QrScanner(
      ticketScannerVideo.value,
      (result) => handleTicketScan(result.data),
      {
        returnDetailedScanResult: true,
        highlightScanRegion: false,
        highlightCodeOutline: true,
        preferredCamera: 'environment'
      }
    )
    await ticketScanner.value.start()
    scannerActive.value = true
  } catch (e) {
    console.error('Scanner error:', e)
  }
}

const stopTicketScanner = () => {
  if (ticketScanner.value) {
    ticketScanner.value.stop()
    ticketScanner.value.destroy()
    ticketScanner.value = null
  }
  scannerActive.value = false
}

const closeScannerModal = () => {
  stopTicketScanner()
  showScannerModal.value = false
  scanResult.value = null
}

const handleTicketScan = async (data) => {
  if (!data) return
  
  // Stop scanning while processing
  stopTicketScanner()
  
  // Extract ticket code from scanned data
  let ticketCode = data.trim()
  
  // Handle ZEKORA_TICKET: prefix
  if (ticketCode.startsWith('ZEKORA_TICKET:')) {
    ticketCode = ticketCode.replace('ZEKORA_TICKET:', '')
  }
  
  await verifyTicket(ticketCode)
}

const verifyManualCode = async () => {
  if (!manualTicketCode.value) return
  stopTicketScanner()
  await verifyTicket(manualTicketCode.value.trim())
}

const verifyTicket = async (ticketCode) => {
  try {
    const res = await ticketAPI.verifyTicket(ticketCode)
    const ticket = res.data?.ticket
    
    if (ticket) {
      scanResult.value = {
        valid: true,
        message: ticket.status === 'used' ? 'Ce ticket a déjà été utilisé' : 'Ticket valide et prêt à être utilisé',
        ticket: ticket
      }
    } else {
      scanResult.value = {
        valid: false,
        message: 'Ticket non trouvé pour cet événement'
      }
    }
  } catch (e) {
    scanResult.value = {
      valid: false,
      message: e.response?.data?.error || 'Erreur lors de la vérification'
    }
  }
}

const markTicketAsUsed = async () => {
  if (!scanResult.value?.ticket?.id) return
  
  try {
    await ticketAPI.useTicket(scanResult.value.ticket.id)
    scanResult.value.ticket.status = 'used'
    scanResult.value.message = 'Ticket marqué comme utilisé ✓'
    alert('Ticket validé avec succès!')
  } catch (e) {
    alert(e.response?.data?.error || 'Erreur lors de la validation')
  }
}

const resetScanner = () => {
  scanResult.value = null
  manualTicketCode.value = ''
  startTicketScanner()
}

onMounted(async () => {
  // Ensure auth store is initialized (user data loaded)
  await authStore.initializeAuth()
  loadEvent()
  loadWallets()
})

onUnmounted(() => {
  stopTicketScanner()
})
</script>

<style scoped>
.event-detail-page {
  min-height: 100vh;
}

.loading-state, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  color: var(--text-muted);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.btn-back {
  display: inline-block;
  padding: 10px 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  text-decoration: none;
  color: var(--text-primary);
  margin-top: 16px;
}

/* Hero Section */
.hero-section {
  height: 300px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  background-size: cover;
  background-position: center;
  position: relative;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, transparent 0%, rgba(0,0,0,0.8) 100%);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 24px;
}

.back-btn {
  color: white;
  text-decoration: none;
  font-size: 14px;
}

.hero-content {
  color: white;
}

.hero-content h1 {
  font-size: 32px;
  font-weight: 700;
  margin: 8px 0;
}

.meta-info {
  display: flex;
  gap: 24px;
  opacity: 0.9;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.draft { background: #6b7280; }
.status-badge.active { background: #10b981; }
.status-badge.ended { background: #f59e0b; }
.status-badge.cancelled { background: #ef4444; }

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 32px;
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
}

.info-section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
}

.info-section h3 {
  font-size: 18px;
  margin-bottom: 16px;
  color: var(--text-primary);
}

.info-section p {
  color: var(--text-muted);
  line-height: 1.6;
}

.dates-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.date-item {
  display: flex;
  flex-direction: column;
}

.date-label {
  font-size: 12px;
  color: var(--text-muted);
}

.date-value {
  font-weight: 600;
  color: var(--text-primary);
}

/* QR Section */
.qr-section .qr-container {
  text-align: center;
}

.qr-image {
  width: 180px;
  height: 180px;
  border-radius: 12px;
}

.qr-placeholder {
  width: 180px;
  height: 180px;
  background: var(--surface-hover);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
}

.qr-placeholder span {
  font-size: 48px;
  margin-bottom: 8px;
}

.qr-placeholder p {
  font-size: 12px;
  color: var(--text-muted);
}

.event-code-box {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--surface-hover);
  padding: 12px 20px;
  border-radius: 10px;
  cursor: pointer;
  margin: 12px 0;
  transition: all 0.2s;
}

.event-code-box:hover {
  background: rgba(99, 102, 241, 0.2);
}

.event-code {
  font-family: monospace;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.copy-icon {
  font-size: 16px;
  opacity: 0.7;
}

.qr-hint {
  font-size: 13px;
  color: var(--text-muted);
}

.btn-qr-modal {
  display: block;
  width: 100%;
  padding: 12px;
  margin-top: 12px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-qr-modal:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4);
}

/* QR Modal */
.qr-modal {
  background: var(--surface);
  border-radius: 24px;
  padding: 32px;
  max-width: 400px;
  width: 90%;
  position: relative;
  text-align: center;
}

.qr-large-container {
  margin-bottom: 24px;
}

.qr-large {
  width: 220px;
  height: 220px;
  border-radius: 16px;
  border: 4px solid var(--border);
}

.qr-placeholder-large {
  width: 220px;
  height: 220px;
  background: var(--surface-hover);
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.qr-placeholder-large span {
  font-size: 64px;
  margin-bottom: 8px;
}

.qr-placeholder-large p {
  color: var(--text-muted);
  font-size: 14px;
}

.qr-modal-info {
  background: var(--surface-hover);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.info-row:not(:last-child) {
  border-bottom: 1px solid var(--border);
}

.info-row .label {
  color: var(--text-muted);
  font-size: 14px;
}

.info-row .value {
  color: var(--text-primary);
  font-weight: 600;
}

.status-pill {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-pill.draft { background: #6b7280; color: white; }
.status-pill.active { background: #10b981; color: white; }
.status-pill.ended { background: #f59e0b; color: white; }
.status-pill.cancelled { background: #ef4444; color: white; }

.qr-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 20px;
}

.qr-action-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 16px 12px;
  background: var(--surface-hover);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.qr-action-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: #6366f1;
}

.action-icon {
  font-size: 24px;
}

.qr-action-btn span:last-child {
  font-size: 12px;
  font-weight: 500;
}

.qr-code-display {
  background: var(--surface-hover);
  border-radius: 12px;
  padding: 16px;
}

.qr-code-label {
  display: block;
  font-size: 11px;
  color: var(--text-muted);
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.qr-code-value {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: var(--surface);
  border: 1px solid #6366f1;
  border-radius: 10px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.qr-code-value:hover {
  background: rgba(99, 102, 241, 0.1);
}

.qr-code-value span:first-child {
  font-family: monospace;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.copy-btn {
  font-size: 16px;
}

/* Tickets Section */
.tickets-section {
  position: sticky;
  top: 24px;
}

.tickets-section h3 {
  font-size: 20px;
  margin-bottom: 16px;
  color: var(--text-primary);
}

.tiers-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.tier-card {
  background: var(--surface);
  border: 2px solid var(--border);
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.tier-card:hover:not(.sold-out) {
  border-color: #6366f1;
}

.tier-card.selected {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.tier-card.sold-out {
  opacity: 0.5;
  cursor: not-allowed;
}

.tier-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.tier-icon {
  font-size: 32px;
}

.tier-info {
  flex: 1;
}

.tier-info h4 {
  margin: 0;
  color: var(--text-primary);
}

.tier-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 4px 0 0;
}

.tier-price {
  text-align: right;
}

.tier-price .price {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
}

.tier-price .currency {
  font-size: 12px;
  color: var(--text-muted);
}

.tier-footer {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--border);
}

.remaining, .sold-out-badge {
  font-size: 12px;
}

.remaining { color: var(--text-muted); }
.sold-out-badge { color: #ef4444; font-weight: 600; }

.btn-purchase {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-purchase:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.organizer-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.btn-publish {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}

.btn-secondary {
  display: block;
  text-align: center;
  padding: 14px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text-primary);
  text-decoration: none;
  font-weight: 600;
}

.actions-row {
  display: flex;
  gap: 12px;
}

.btn-edit, .btn-delete {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  border-radius: 12px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  text-decoration: none;
  font-size: 14px;
}

.btn-edit {
  background: var(--surface-hover);
  color: var(--text-primary);
  border: 1px solid var(--border);
}

.btn-edit:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: #6366f1;
}

.btn-delete {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.2);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
  overflow-y: auto;
}

.purchase-modal, .success-modal {
  background: var(--surface);
  border-radius: 20px;
  padding: 28px;
  max-width: 440px;
  width: 100%;
  position: relative;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid var(--border);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: var(--surface-hover);
  border: 1px solid var(--border);
  font-size: 18px;
  color: var(--text-muted);
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--error);
  color: white;
  border-color: var(--error);
}

.purchase-modal h2 {
  margin: 0 0 8px;
  color: var(--text-primary);
  font-size: 22px;
}

.modal-subtitle {
  color: var(--text-muted);
  margin-bottom: 24px;
  font-size: 15px;
}

.form-group {
  margin-bottom: 18px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.form-group label span {
  color: var(--error);
}

.form-group input, .form-group select {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid var(--border);
  border-radius: 12px;
  background: var(--background);
  color: var(--text-primary);
  font-size: 15px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-group input:focus, .form-group select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.form-group input::placeholder {
  color: var(--text-muted);
}

.form-group select option {
  background: var(--surface);
  color: var(--text-primary);
  padding: 10px;
}

.btn-confirm {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 12px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-confirm:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);
}

.btn-confirm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.success-modal {
  text-align: center;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.ticket-qr {
  width: 180px;
  height: 180px;
  border-radius: 12px;
  margin: 16px 0;
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
  margin-bottom: 20px;
}

.btn-primary {
  display: inline-block;
  padding: 14px 32px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 600;
}

@media (max-width: 900px) {
  .content-grid {
    grid-template-columns: 1fr;
    padding: 20px;
  }
  
  .tickets-section {
    position: static;
  }
}

/* Scan Button */
.btn-scan {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-scan:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(245, 158, 11, 0.4);
}

/* Scanner Modal */
.scanner-modal {
  background: var(--surface);
  border-radius: 24px;
  padding: 24px;
  max-width: 420px;
  width: 90%;
  position: relative;
  text-align: center;
}

.scanner-modal h2 {
  margin: 0 0 8px;
  color: var(--text-primary);
}

.scanner-modal .modal-subtitle {
  color: var(--text-muted);
  margin-bottom: 20px;
}

.scanner-camera-container {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  border-radius: 16px;
  overflow: hidden;
  background: #000;
  margin-bottom: 16px;
}

.scanner-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scanner-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.scanner-frame {
  width: 70%;
  height: 70%;
  position: relative;
}

.scanner-corner {
  position: absolute;
  width: 30px;
  height: 30px;
  border-color: #f59e0b;
  border-style: solid;
  border-width: 0;
}

.scanner-corner.top-left {
  top: 0; left: 0;
  border-top-width: 4px;
  border-left-width: 4px;
  border-top-left-radius: 12px;
}

.scanner-corner.top-right {
  top: 0; right: 0;
  border-top-width: 4px;
  border-right-width: 4px;
  border-top-right-radius: 12px;
}

.scanner-corner.bottom-left {
  bottom: 0; left: 0;
  border-bottom-width: 4px;
  border-left-width: 4px;
  border-bottom-left-radius: 12px;
}

.scanner-corner.bottom-right {
  bottom: 0; right: 0;
  border-bottom-width: 4px;
  border-right-width: 4px;
  border-bottom-right-radius: 12px;
}

.scanner-line {
  position: absolute;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, transparent, #f59e0b, transparent);
  animation: scan-anim 2s linear infinite;
}

@keyframes scan-anim {
  0% { top: 0; }
  50% { top: 100%; }
  100% { top: 0; }
}

.manual-entry {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.ticket-code-input {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--surface);
  color: var(--text-primary);
  font-family: monospace;
  text-align: center;
}

.btn-verify {
  padding: 12px 20px;
  background: #6366f1;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}

.btn-verify:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.scan-result {
  padding: 20px;
  border-radius: 16px;
  margin-top: 16px;
}

.scan-result.valid {
  background: rgba(16, 185, 129, 0.1);
  border: 2px solid #10b981;
}

.scan-result.invalid {
  background: rgba(239, 68, 68, 0.1);
  border: 2px solid #ef4444;
}

.result-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.scan-result h3 {
  margin: 0 0 8px;
  color: var(--text-primary);
}

.scan-result p {
  color: var(--text-muted);
  margin: 0 0 12px;
}

.ticket-info {
  background: var(--surface-hover);
  padding: 12px;
  border-radius: 10px;
  margin: 12px 0;
  text-align: left;
}

.ticket-info p {
  margin: 4px 0;
  font-size: 14px;
}

.btn-mark-used {
  width: 100%;
  padding: 14px;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 8px;
}

.btn-scan-another {
  width: 100%;
  padding: 12px;
  background: var(--surface-hover);
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-radius: 12px;
  font-weight: 500;
  cursor: pointer;
}

/* Purchase Summary */
.purchase-summary {
  background: var(--surface-hover);
  border-radius: 12px;
  padding: 16px;
  margin-top: 16px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-primary);
  font-weight: 500;
}

.summary-price {
  font-size: 18px;
  font-weight: 700;
  color: var(--primary);
}

/* PIN Confirmation Modal */
.pin-modal {
  background: var(--surface);
  border-radius: 24px;
  padding: 32px;
  max-width: 380px;
  width: 95%;
  text-align: center;
  border: 1px solid var(--border);
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
}

.pin-header {
  margin-bottom: 24px;
}

.pin-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.pin-header h3 {
  color: var(--text-primary);
  font-size: 22px;
  margin: 0 0 8px;
}

.pin-subtitle {
  color: var(--text-muted);
  font-size: 14px;
  margin: 0;
}

.pin-amount {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 24px;
}

.pin-amount-value {
  display: block;
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
}

.pin-amount-currency {
  color: var(--text-muted);
  font-size: 14px;
}

.pin-input-container {
  margin-bottom: 20px;
}

.pin-input {
  width: 100%;
  padding: 18px;
  text-align: center;
  font-size: 28px;
  font-weight: 600;
  letter-spacing: 12px;
  border: 2px solid var(--border);
  border-radius: 16px;
  background: var(--background);
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.pin-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.15);
}

.pin-input::placeholder {
  color: var(--text-muted);
  letter-spacing: 8px;
}

.pin-hint {
  color: var(--text-muted);
  font-size: 12px;
  margin-top: 8px;
}

.pin-error {
  color: var(--error);
  font-size: 14px;
  margin-bottom: 16px;
  padding: 12px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 10px;
}

.pin-actions {
  display: flex;
  gap: 12px;
}

.btn-cancel {
  flex: 1;
  padding: 14px;
  background: var(--surface-hover);
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background: var(--border);
}

.btn-pay {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-pay:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);
}

.btn-pay:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner-small {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
