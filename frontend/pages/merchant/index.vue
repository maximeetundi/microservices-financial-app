<template>
  <NuxtLayout name="dashboard">
    <div class="merchant-page">
      <!-- Header -->
      <div class="page-header">
        <div class="header-content">
          <h1>💼 Espace Marchand</h1>
          <p class="subtitle">Recevez des paiements via QR code</p>
        </div>
        <button class="btn-primary" @click="showCreateModal = true">
          <span class="icon">➕</span>
          Créer un paiement
        </button>
      </div>

      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <span class="stat-value">{{ formatCurrency(stats.totalAmount) }}</span>
            <span class="stat-label">Total reçu</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-content">
            <span class="stat-value">{{ stats.totalPayments }}</span>
            <span class="stat-label">Paiements</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <span class="stat-value">{{ pendingPayments.length }}</span>
            <span class="stat-label">En attente</span>
          </div>
        </div>
      </div>

      <!-- Payment Requests List -->
      <div class="section">
        <div class="section-header">
          <h2>Demandes de paiement</h2>
          <div class="tabs">
            <button 
              :class="{ active: activeTab === 'pending' }" 
              @click="activeTab = 'pending'"
            >
              En attente
            </button>
            <button 
              :class="{ active: activeTab === 'completed' }" 
              @click="activeTab = 'completed'"
            >
              Terminés
            </button>
          </div>
        </div>

        <div class="payments-list">
          <div 
            v-for="payment in filteredPayments" 
            :key="payment.id" 
            class="payment-card"
            @click="showPaymentDetails(payment)"
          >
            <div class="payment-info">
              <div class="payment-title">{{ payment.title }}</div>
              <div class="payment-meta">
                <span class="payment-type">{{ getTypeLabel(payment.type) }}</span>
                <span class="payment-date">{{ formatDate(payment.created_at) }}</span>
              </div>
            </div>
            <div class="payment-amount">
              <span v-if="payment.amount" class="amount">
                {{ formatCurrency(payment.amount, payment.currency) }}
              </span>
              <span v-else class="variable">Variable</span>
            </div>
            <div class="payment-status" :class="payment.status">
              {{ getStatusLabel(payment.status) }}
            </div>
            <button class="btn-icon" @click.stop="showQRCode(payment)">
              <span>📱</span>
            </button>
          </div>

          <div v-if="filteredPayments.length === 0" class="empty-state">
            <span class="empty-icon">📭</span>
            <p>Aucune demande de paiement</p>
            <button class="btn-secondary" @click="showCreateModal = true">
              Créer votre première demande
            </button>
          </div>
        </div>
      </div>

      <!-- Create Payment Modal -->
      <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
        <div class="modal">
          <div class="modal-header">
            <h3>Nouvelle demande de paiement</h3>
            <button class="btn-close" @click="showCreateModal = false">✕</button>
          </div>
          
          <form @submit.prevent="createPayment" class="modal-body">
            <!-- Payment Type -->
            <div class="form-group">
              <label>Type de paiement</label>
              <div class="type-selector">
                <button 
                  type="button"
                  :class="{ active: newPayment.type === 'fixed' }"
                  @click="newPayment.type = 'fixed'"
                >
                  <span class="type-icon">🏷️</span>
                  <span class="type-label">Prix fixe</span>
                </button>
                <button 
                  type="button"
                  :class="{ active: newPayment.type === 'variable' }"
                  @click="newPayment.type = 'variable'"
                >
                  <span class="type-icon">🎁</span>
                  <span class="type-label">Variable</span>
                </button>
              </div>
            </div>

            <!-- Wallet Selection -->
            <div class="form-group">
              <label>Portefeuille de réception</label>
              <select v-model="newPayment.wallet_id" required>
                <option value="">Sélectionner...</option>
                <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                  {{ wallet.currency }} - {{ formatCurrency(wallet.balance, wallet.currency) }}
                </option>
              </select>
            </div>

            <!-- Title -->
            <div class="form-group">
              <label>Titre</label>
              <input 
                v-model="newPayment.title" 
                type="text" 
                placeholder="Ex: iPhone 15 Pro"
                required
              />
            </div>

            <!-- Amount (for fixed) -->
            <div v-if="newPayment.type === 'fixed'" class="form-group">
              <label>Montant</label>
              <div class="amount-input">
                <input 
                  v-model.number="newPayment.amount" 
                  type="number" 
                  step="0.01"
                  min="0"
                  placeholder="0.00"
                  required
                />
                <span class="currency">{{ selectedCurrency }}</span>
              </div>
            </div>

            <!-- Min/Max (for variable) -->
            <div v-if="newPayment.type === 'variable'" class="form-row">
              <div class="form-group">
                <label>Montant min (optionnel)</label>
                <input 
                  v-model.number="newPayment.min_amount" 
                  type="number" 
                  step="0.01"
                  placeholder="0.00"
                />
              </div>
              <div class="form-group">
                <label>Montant max (optionnel)</label>
                <input 
                  v-model.number="newPayment.max_amount" 
                  type="number" 
                  step="0.01"
                  placeholder="0.00"
                />
              </div>
            </div>

            <!-- Description -->
            <div class="form-group">
              <label>Description (optionnel)</label>
              <textarea 
                v-model="newPayment.description" 
                placeholder="Description du produit ou service..."
                rows="3"
              ></textarea>
            </div>

            <!-- Expiration -->
            <div class="form-group">
              <label>Expiration</label>
              <select v-model="newPayment.expires_in_minutes">
                <option :value="60">1 heure</option>
                <option :value="1440">24 heures</option>
                <option :value="10080">7 jours</option>
                <option :value="-1">Jamais</option>
              </select>
            </div>

            <!-- Reusable -->
            <div class="form-group checkbox">
              <input type="checkbox" id="reusable" v-model="newPayment.reusable" />
              <label for="reusable">Réutilisable (plusieurs paiements)</label>
            </div>

            <div class="modal-footer">
              <button type="button" class="btn-secondary" @click="showCreateModal = false">
                Annuler
              </button>
              <button type="submit" class="btn-primary" :disabled="creating">
                {{ creating ? 'Création...' : 'Créer le QR code' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- QR Code Modal -->
      <div v-if="showQRModal" class="modal-overlay" @click.self="showQRModal = false">
        <div class="modal qr-modal">
          <div class="modal-header">
            <h3>{{ selectedPayment?.title }}</h3>
            <button class="btn-close" @click="showQRModal = false">✕</button>
          </div>
          
          <div class="qr-content">
            <div class="qr-code">
              <img :src="qrCodeImage" alt="QR Code" />
            </div>
            
            <div class="payment-details">
              <div v-if="selectedPayment?.amount" class="detail-row">
                <span>Montant:</span>
                <strong>{{ formatCurrency(selectedPayment.amount, selectedPayment.currency) }}</strong>
              </div>
              <div v-else class="detail-row">
                <span>Montant:</span>
                <strong>À définir par le client</strong>
              </div>
              <div class="detail-row">
                <span>Statut:</span>
                <span :class="'status ' + selectedPayment?.status">
                  {{ getStatusLabel(selectedPayment?.status) }}
                </span>
              </div>
              <div v-if="selectedPayment?.expires_at" class="detail-row">
                <span>Expire:</span>
                <span>{{ formatDate(selectedPayment.expires_at) }}</span>
              </div>
            </div>

            <div class="qr-actions">
              <button class="btn-secondary" @click="copyPaymentLink">
                📋 Copier le lien
              </button>
              <button class="btn-secondary" @click="downloadQR">
                ⬇️ Télécharger
              </button>
              <button class="btn-secondary" @click="sharePayment">
                📤 Partager
              </button>
            </div>

            <div class="payment-link">
              <input type="text" :value="selectedPayment?.payment_link" readonly />
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const { merchantApi, walletApi } = useApi()

const showCreateModal = ref(false)
const showQRModal = ref(false)
const creating = ref(false)
const activeTab = ref('pending')
const selectedPayment = ref(null)
const qrCodeImage = ref('')

const payments = ref([])
const wallets = ref([])
const stats = ref({
  totalAmount: 0,
  totalPayments: 0,
  totalFees: 0
})

const newPayment = ref({
  type: 'fixed',
  wallet_id: '',
  title: '',
  amount: null,
  min_amount: null,
  max_amount: null,
  description: '',
  expires_in_minutes: 60,
  reusable: false
})

const selectedCurrency = computed(() => {
  const wallet = wallets.value.find(w => w.id === newPayment.value.wallet_id)
  return wallet?.currency || 'EUR'
})

const pendingPayments = computed(() => 
  payments.value.filter(p => p.status === 'pending')
)

const filteredPayments = computed(() => {
  if (activeTab.value === 'pending') {
    return payments.value.filter(p => p.status === 'pending')
  }
  return payments.value.filter(p => p.status !== 'pending')
})

onMounted(async () => {
  await loadData()
})

async function loadData() {
  try {
    const [paymentsRes, walletsRes] = await Promise.all([
      merchantApi.getPayments(),
      walletApi.getWallets()
    ])
    payments.value = paymentsRes.payments || []
    wallets.value = walletsRes || []
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

async function createPayment() {
  creating.value = true
  try {
    const response = await merchantApi.createPayment({
      ...newPayment.value,
      currency: selectedCurrency.value
    })
    
    payments.value.unshift(response.payment_request)
    showCreateModal.value = false
    
    // Show QR code
    selectedPayment.value = response.payment_request
    qrCodeImage.value = response.qr_code_base64
    showQRModal.value = true
    
    // Reset form
    newPayment.value = {
      type: 'fixed',
      wallet_id: '',
      title: '',
      amount: null,
      min_amount: null,
      max_amount: null,
      description: '',
      expires_in_minutes: 60,
      reusable: false
    }
  } catch (error) {
    alert('Erreur: ' + error.message)
  } finally {
    creating.value = false
  }
}

async function showQRCode(payment) {
  selectedPayment.value = payment
  try {
    const response = await merchantApi.getQRCode(payment.id)
    qrCodeImage.value = response.qr_code_base64
    showQRModal.value = true
  } catch (error) {
    console.error('Failed to get QR code:', error)
  }
}

function showPaymentDetails(payment) {
  selectedPayment.value = payment
  showQRCode(payment)
}

function copyPaymentLink() {
  navigator.clipboard.writeText(selectedPayment.value.payment_link)
  alert('Lien copié!')
}

function downloadQR() {
  const link = document.createElement('a')
  link.download = `qr-${selectedPayment.value.id}.png`
  link.href = qrCodeImage.value
  link.click()
}

function sharePayment() {
  if (navigator.share) {
    navigator.share({
      title: selectedPayment.value.title,
      text: `Paiement: ${selectedPayment.value.title}`,
      url: selectedPayment.value.payment_link
    })
  } else {
    copyPaymentLink()
  }
}

function formatCurrency(amount, currency = 'EUR') {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency
  }).format(amount || 0)
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getTypeLabel(type) {
  const labels = {
    fixed: 'Prix fixe',
    variable: 'Variable',
    invoice: 'Facture'
  }
  return labels[type] || type
}

function getStatusLabel(status) {
  const labels = {
    pending: 'En attente',
    paid: 'Payé',
    expired: 'Expiré',
    cancelled: 'Annulé'
  }
  return labels[status] || status
}
</script>

<style scoped>
.merchant-page {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 1.75rem;
  color: var(--text-primary);
}

.subtitle {
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 1rem;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  border: 1px solid var(--border-color);
}

.stat-icon {
  font-size: 2rem;
}

.stat-value {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.section {
  background: var(--card-bg);
  border-radius: 1rem;
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.section-header h2 {
  font-size: 1.125rem;
  color: var(--text-primary);
}

.tabs {
  display: flex;
  gap: 0.5rem;
}

.tabs button {
  padding: 0.5rem 1rem;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.tabs button.active {
  background: var(--primary-color);
  color: white;
}

.payments-list {
  padding: 1rem;
}

.payment-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border-radius: 0.75rem;
  background: rgba(255, 255, 255, 0.02);
  margin-bottom: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.payment-card:hover {
  background: rgba(255, 255, 255, 0.05);
}

.payment-info {
  flex: 1;
}

.payment-title {
  font-weight: 600;
  color: var(--text-primary);
}

.payment-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.payment-amount .amount {
  font-weight: 700;
  font-size: 1.125rem;
  color: var(--success-color);
}

.payment-amount .variable {
  color: var(--text-secondary);
  font-style: italic;
}

.payment-status {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.payment-status.pending {
  background: rgba(251, 191, 36, 0.1);
  color: #fbbf24;
}

.payment-status.paid {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.payment-status.expired {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 1rem;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal {
  background: var(--card-bg);
  border-radius: 1rem;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid var(--border-color);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  font-size: 1.25rem;
  color: var(--text-primary);
}

.btn-close {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: var(--text-secondary);
  cursor: pointer;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  color: var(--text-primary);
  font-size: 1rem;
}

.form-group.checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-group.checkbox input {
  width: auto;
}

.form-group.checkbox label {
  margin: 0;
}

.type-selector {
  display: flex;
  gap: 1rem;
}

.type-selector button {
  flex: 1;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.type-selector button.active {
  border-color: var(--primary-color);
  background: rgba(99, 102, 241, 0.1);
}

.type-icon {
  font-size: 1.5rem;
}

.type-label {
  color: var(--text-primary);
  font-weight: 500;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.amount-input {
  position: relative;
}

.amount-input .currency {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
}

.modal-footer {
  display: flex;
  gap: 1rem;
  padding-top: 1rem;
}

.modal-footer button {
  flex: 1;
}

/* QR Modal */
.qr-modal {
  max-width: 400px;
}

.qr-content {
  padding: 1.5rem;
  text-align: center;
}

.qr-code {
  background: white;
  padding: 1rem;
  border-radius: 1rem;
  display: inline-block;
  margin-bottom: 1.5rem;
}

.qr-code img {
  width: 200px;
  height: 200px;
}

.payment-details {
  text-align: left;
  margin-bottom: 1.5rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--border-color);
}

.qr-actions {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.qr-actions button {
  flex: 1;
  padding: 0.5rem;
  font-size: 0.75rem;
}

.payment-link input {
  width: 100%;
  padding: 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.75rem;
  text-align: center;
}

/* Buttons */
.btn-primary {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon {
  background: transparent;
  border: none;
  font-size: 1.25rem;
  cursor: pointer;
  padding: 0.5rem;
}

.status.paid {
  color: #22c55e;
}

.status.pending {
  color: #fbbf24;
}

.status.expired {
  color: #ef4444;
}
</style>
