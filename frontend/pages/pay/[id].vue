<template>
  <div class="pay-page">
    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du paiement...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error-state">
      <span class="error-icon">❌</span>
      <h2>Paiement introuvable</h2>
      <p>{{ error }}</p>
      <NuxtLink to="/" class="btn-secondary">Retour à l'accueil</NuxtLink>
    </div>

    <!-- Expired -->
    <div v-else-if="payment?.status === 'expired'" class="expired-state">
      <span class="status-icon">⏰</span>
      <h2>Paiement expiré</h2>
      <p>Cette demande de paiement a expiré.</p>
      <NuxtLink to="/" class="btn-secondary">Retour à l'accueil</NuxtLink>
    </div>

    <!-- Already Paid -->
    <div v-else-if="payment?.is_paid" class="paid-state">
      <span class="status-icon">✅</span>
      <h2>Déjà payé</h2>
      <p>Cette demande de paiement a déjà été réglée.</p>
      <NuxtLink to="/" class="btn-secondary">Retour à l'accueil</NuxtLink>
    </div>

    <!-- Payment Form -->
    <div v-else-if="payment" class="payment-container">
      <div class="payment-header">
        <div class="merchant-info">
          <div class="merchant-avatar">💼</div>
          <span class="merchant-name">Marchand</span>
        </div>
      </div>

      <div class="payment-card">
        <h1 class="payment-title">{{ payment.title }}</h1>
        <p v-if="payment.description" class="payment-description">
          {{ payment.description }}
        </p>

        <!-- Fixed Amount -->
        <div v-if="payment.amount" class="amount-display">
          <span class="amount">{{ formatCurrency(payment.amount, payment.currency) }}</span>
        </div>

        <!-- Variable Amount -->
        <div v-else class="amount-input-section">
          <label>Montant à payer</label>
          <div class="amount-input-wrapper">
            <input 
              v-model.number="amountToPay" 
              type="number" 
              step="0.01"
              :min="payment.min_amount || 0"
              :max="payment.max_amount"
              :placeholder="getAmountPlaceholder()"
              required
            />
            <span class="currency">{{ payment.currency }}</span>
          </div>
          <div v-if="payment.min_amount || payment.max_amount" class="amount-limits">
            <span v-if="payment.min_amount">Min: {{ formatCurrency(payment.min_amount, payment.currency) }}</span>
            <span v-if="payment.max_amount">Max: {{ formatCurrency(payment.max_amount, payment.currency) }}</span>
          </div>
        </div>

        <!-- Expiration Warning -->
        <div v-if="payment.expires_at && !payment.never_expires" class="expiration-warning">
          <span class="icon">⏳</span>
          Expire {{ formatRelativeTime(payment.expires_at) }}
        </div>

        <!-- Login Required -->
        <div v-if="!isAuthenticated" class="login-required">
          <p>Connectez-vous pour payer</p>
          <NuxtLink :to="`/auth/login?redirect=/pay/${payment.payment_id}`" class="btn-primary">
            Se connecter
          </NuxtLink>
          <NuxtLink to="/auth/register" class="btn-secondary">
            Créer un compte
          </NuxtLink>
        </div>

        <!-- Payment Form -->
        <div v-else class="pay-form">
          <div class="wallet-select">
            <label>Payer depuis</label>
            <select v-model="selectedWallet" required>
              <option value="">Sélectionner un portefeuille</option>
              <option 
                v-for="wallet in compatibleWallets" 
                :key="wallet.id" 
                :value="wallet.id"
                :disabled="wallet.balance < (payment.amount || amountToPay)"
              >
                {{ wallet.currency }} - {{ formatCurrency(wallet.balance, wallet.currency) }}
                <span v-if="wallet.balance < (payment.amount || amountToPay)"> (Solde insuffisant)</span>
              </option>
            </select>
          </div>

          <button 
            class="btn-pay" 
            @click="processPayment"
            :disabled="!canPay || processing"
          >
            <span v-if="processing" class="spinner-small"></span>
            <span v-else>
              Payer {{ formatCurrency(payment.amount || amountToPay, payment.currency) }}
            </span>
          </button>

          <p class="secure-note">
            🔒 Paiement sécurisé
          </p>
        </div>
      </div>

      <!-- Success Modal -->
      <div v-if="paymentSuccess" class="success-overlay">
        <div class="success-modal">
          <div class="success-icon">✅</div>
          <h2>Paiement réussi!</h2>
          <p>Votre paiement de {{ formatCurrency(paidAmount, payment.currency) }} a été effectué.</p>
          <NuxtLink to="/dashboard" class="btn-primary">
            Voir mes transactions
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useApi } from '~/composables/useApi'

const route = useRoute()
const { paymentApi, walletApi } = useApi()

const loading = ref(true)
const error = ref(null)
const processing = ref(false)
const paymentSuccess = ref(false)
const paidAmount = ref(0)

const payment = ref(null)
const wallets = ref([])
const selectedWallet = ref('')
const amountToPay = ref(null)

const isAuthenticated = computed(() => {
  if (process.client) {
    return !!localStorage.getItem('token')
  }
  return false
})

const compatibleWallets = computed(() => {
  if (!payment.value) return []
  return wallets.value.filter(w => w.currency === payment.value.currency)
})

const canPay = computed(() => {
  if (!selectedWallet.value) return false
  
  const amount = payment.value?.amount || amountToPay.value
  if (!amount || amount <= 0) return false
  
  const wallet = wallets.value.find(w => w.id === selectedWallet.value)
  if (!wallet || wallet.balance < amount) return false
  
  return true
})

onMounted(async () => {
  await loadPayment()
  if (isAuthenticated.value) {
    await loadWallets()
  }
})

async function loadPayment() {
  loading.value = true
  try {
    const paymentId = route.params.id
    payment.value = await paymentApi.getPaymentDetails(paymentId)
  } catch (err) {
    error.value = err.message || 'Paiement introuvable'
  } finally {
    loading.value = false
  }
}

async function loadWallets() {
  try {
    wallets.value = await walletApi.getWallets()
  } catch (err) {
    console.error('Failed to load wallets:', err)
  }
}

async function processPayment() {
  if (!canPay.value) return
  
  processing.value = true
  try {
    const amount = payment.value.amount || amountToPay.value
    
    await paymentApi.payPayment(payment.value.payment_id, {
      from_wallet_id: selectedWallet.value,
      amount: amount
    })
    
    paidAmount.value = amount
    paymentSuccess.value = true
  } catch (err) {
    alert('Erreur: ' + (err.message || 'Paiement échoué'))
  } finally {
    processing.value = false
  }
}

function formatCurrency(amount, currency = 'EUR') {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency
  }).format(amount || 0)
}

function formatRelativeTime(date) {
  const now = new Date()
  const target = new Date(date)
  const diff = target - now
  
  if (diff < 0) return 'expiré'
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `dans ${days} jour(s)`
  if (hours > 0) return `dans ${hours} heure(s)`
  return `dans ${minutes} minute(s)`
}

function getAmountPlaceholder() {
  if (payment.value?.min_amount) {
    return `Minimum ${payment.value.min_amount}`
  }
  return 'Entrez le montant'
}
</script>

<style scoped>
.pay-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a3e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.loading-state,
.error-state,
.expired-state,
.paid-state {
  text-align: center;
  color: var(--text-primary);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon,
.status-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 1rem;
}

.payment-container {
  width: 100%;
  max-width: 420px;
}

.payment-header {
  text-align: center;
  margin-bottom: 1rem;
}

.merchant-info {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.5rem 1rem;
  border-radius: 2rem;
}

.merchant-avatar {
  font-size: 1.25rem;
}

.merchant-name {
  color: var(--text-primary);
  font-weight: 500;
}

.payment-card {
  background: var(--card-bg);
  border-radius: 1.5rem;
  padding: 2rem;
  border: 1px solid var(--border-color);
}

.payment-title {
  font-size: 1.5rem;
  color: var(--text-primary);
  text-align: center;
  margin-bottom: 0.5rem;
}

.payment-description {
  color: var(--text-secondary);
  text-align: center;
  margin-bottom: 1.5rem;
}

.amount-display {
  text-align: center;
  margin: 2rem 0;
}

.amount-display .amount {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.amount-input-section {
  margin: 1.5rem 0;
}

.amount-input-section label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.amount-input-wrapper {
  position: relative;
}

.amount-input-wrapper input {
  width: 100%;
  padding: 1rem;
  font-size: 1.5rem;
  text-align: center;
  background: rgba(255, 255, 255, 0.05);
  border: 2px solid var(--border-color);
  border-radius: 0.75rem;
  color: var(--text-primary);
}

.amount-input-wrapper input:focus {
  border-color: var(--primary-color);
  outline: none;
}

.amount-input-wrapper .currency {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
}

.amount-limits {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin-top: 0.5rem;
}

.expiration-warning {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: rgba(251, 191, 36, 0.1);
  border-radius: 0.5rem;
  color: #fbbf24;
  font-size: 0.875rem;
  margin: 1rem 0;
}

.login-required {
  text-align: center;
  padding: 1.5rem 0;
}

.login-required p {
  color: var(--text-secondary);
  margin-bottom: 1rem;
}

.login-required .btn-primary,
.login-required .btn-secondary {
  display: block;
  width: 100%;
  margin-bottom: 0.5rem;
  text-decoration: none;
  text-align: center;
}

.pay-form {
  margin-top: 1.5rem;
}

.wallet-select {
  margin-bottom: 1rem;
}

.wallet-select label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.wallet-select select {
  width: 100%;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  color: var(--text-primary);
  font-size: 1rem;
}

.btn-pay {
  width: 100%;
  padding: 1rem;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  border: none;
  border-radius: 0.75rem;
  color: white;
  font-size: 1.125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.btn-pay:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3);
}

.btn-pay:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.secure-note {
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.75rem;
  margin-top: 1rem;
}

/* Success Modal */
.success-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.success-modal {
  background: var(--card-bg);
  border-radius: 1.5rem;
  padding: 3rem 2rem;
  text-align: center;
  max-width: 360px;
  border: 1px solid var(--border-color);
}

.success-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.success-modal h2 {
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.success-modal p {
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
}

.success-modal .btn-primary {
  display: inline-block;
  padding: 0.75rem 2rem;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  border-radius: 0.5rem;
  text-decoration: none;
  font-weight: 600;
}

/* Buttons */
.btn-primary {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  font-weight: 600;
  cursor: pointer;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
  text-decoration: none;
  display: inline-block;
}
</style>
