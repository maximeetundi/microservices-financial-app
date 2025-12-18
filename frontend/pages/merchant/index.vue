<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-base">💼 Espace Marchand</h1>
          <p class="text-muted mt-1">Recevez des paiements via QR code</p>
        </div>
        <button class="btn-primary flex items-center gap-2" @click="showCreateModal = true">
          <span class="text-xl">+</span>
          Créer un paiement
        </button>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="glass-card p-6 flex items-center gap-4">
          <div class="p-3 rounded-xl bg-primary/10 text-primary text-2xl">💰</div>
          <div>
            <span class="block text-2xl font-bold text-base">{{ formatCurrency(stats.totalAmount) }}</span>
            <span class="text-sm text-muted">Total reçu</span>
          </div>
        </div>
        <div class="glass-card p-6 flex items-center gap-4">
          <div class="p-3 rounded-xl bg-purple-500/10 text-purple-500 text-2xl">📊</div>
          <div>
            <span class="block text-2xl font-bold text-base">{{ stats.totalPayments }}</span>
            <span class="text-sm text-muted">Paiements</span>
          </div>
        </div>
        <div class="glass-card p-6 flex items-center gap-4">
          <div class="p-3 rounded-xl bg-amber-500/10 text-amber-500 text-2xl">⏳</div>
          <div>
            <span class="block text-2xl font-bold text-base">{{ pendingPayments.length }}</span>
            <span class="text-sm text-muted">En attente</span>
          </div>
        </div>
      </div>

      <!-- Payment Requests List -->
      <div class="glass-card">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-6 border-b border-secondary-200 dark:border-secondary-800">
          <h2 class="text-lg font-bold text-base">Demandes de paiement</h2>
          <div class="flex gap-2 bg-surface-hover p-1 rounded-lg">
            <button 
              v-for="tab in ['pending', 'completed']" 
              :key="tab"
              class="px-4 py-2 rounded-md text-sm font-medium transition-all"
              :class="activeTab === tab ? 'bg-bg-surface shadow-sm text-primary' : 'text-muted hover:text-base'"
              @click="activeTab = tab"
            >
              {{ tab === 'pending' ? 'En attente' : 'Terminés' }}
            </button>
          </div>
        </div>

        <div class="p-4">
          <div 
            v-for="payment in filteredPayments" 
            :key="payment.id" 
            class="group flex flex-col sm:flex-row sm:items-center gap-4 p-4 rounded-xl hover:bg-surface-hover transition-colors cursor-pointer border border-transparent hover:border-secondary-200 dark:hover:border-secondary-700"
            @click="showPaymentDetails(payment)"
          >
            <div class="flex-1">
              <div class="font-bold text-base mb-1">{{ payment.title }}</div>
              <div class="flex items-center gap-3 text-xs text-muted">
                <span class="px-2 py-0.5 rounded-full bg-secondary-100 dark:bg-secondary-800 text-secondary-600 dark:text-secondary-400">
                  {{ getTypeLabel(payment.type) }}
                </span>
                <span>{{ formatDate(payment.created_at) }}</span>
              </div>
            </div>
            
            <div class="text-right">
              <div v-if="payment.amount" class="text-lg font-bold text-success">
                {{ formatCurrency(payment.amount, payment.currency) }}
              </div>
              <div v-else class="text-muted italic">Variable</div>
            </div>

            <div class="px-3 py-1 rounded-full text-xs font-bold" :class="getStatusClass(payment.status)">
              {{ getStatusLabel(payment.status) }}
            </div>

            <button class="p-2 rounded-lg hover:bg-secondary-200 dark:hover:bg-secondary-700 text-muted transition-colors" @click.stop="showQRCode(payment)">
              <span class="text-xl">📱</span>
            </button>
          </div>

          <div v-if="filteredPayments.length === 0" class="text-center py-12">
            <span class="text-4xl mb-4 block">📭</span>
            <p class="text-muted mb-4">Aucune demande de paiement</p>
            <button class="btn-primary" @click="showCreateModal = true">
              Créer votre première demande
            </button>
          </div>
        </div>
      </div>

      <!-- Create Payment Modal -->
      <div v-if="showCreateModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4" @click.self="showCreateModal = false">
        <div class="glass-card w-full max-w-lg max-h-[90vh] overflow-y-auto animate-fade-in-up">
          <div class="flex justify-between items-center p-6 border-b border-secondary-200 dark:border-secondary-800">
            <h3 class="text-xl font-bold text-base">Nouvelle demande</h3>
            <button class="text-muted hover:text-base text-2xl" @click="showCreateModal = false">×</button>
          </div>
          
          <form @submit.prevent="createPayment" class="p-6 space-y-6">
            <!-- Payment Type -->
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted">Type de paiement</label>
              <div class="grid grid-cols-2 gap-4">
                <button 
                  type="button"
                  class="p-4 rounded-xl border-2 transition-all flex flex-col items-center gap-2"
                  :class="newPayment.type === 'fixed' 
                    ? 'border-primary bg-primary/5 text-primary' 
                    : 'border-secondary-200 dark:border-secondary-700 text-muted hover:border-primary/50'"
                  @click="newPayment.type = 'fixed'"
                >
                  <span class="text-2xl">🏷️</span>
                  <span class="font-bold">Prix fixe</span>
                </button>
                <button 
                  type="button"
                  class="p-4 rounded-xl border-2 transition-all flex flex-col items-center gap-2"
                  :class="newPayment.type === 'variable' 
                    ? 'border-primary bg-primary/5 text-primary' 
                    : 'border-secondary-200 dark:border-secondary-700 text-muted hover:border-primary/50'"
                  @click="newPayment.type = 'variable'"
                >
                  <span class="text-2xl">🎁</span>
                  <span class="font-bold">Variable</span>
                </button>
              </div>
            </div>

            <!-- Wallet Selection -->
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted">Portefeuille de réception</label>
              <select v-model="newPayment.wallet_id" class="input-field" required>
                <option value="">Sélectionner...</option>
                <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                  {{ wallet.currency }} - {{ formatCurrency(wallet.balance, wallet.currency) }}
                </option>
              </select>
            </div>

            <!-- Title -->
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted">Titre</label>
              <input 
                v-model="newPayment.title" 
                type="text" 
                class="input-field"
                placeholder="Ex: iPhone 15 Pro"
                required
              />
            </div>

            <!-- Amount (for fixed) -->
            <div v-if="newPayment.type === 'fixed'" class="space-y-2">
              <label class="text-sm font-medium text-muted">Montant</label>
              <div class="relative">
                <input 
                  v-model.number="newPayment.amount" 
                  type="number" 
                  step="0.01"
                  min="0"
                  class="input-field pr-16"
                  placeholder="0.00"
                  required
                />
                <span class="absolute right-4 top-1/2 -translate-y-1/2 font-bold text-muted">{{ selectedCurrency }}</span>
              </div>
            </div>

            <!-- Min/Max (for variable) -->
            <div v-if="newPayment.type === 'variable'" class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="text-sm font-medium text-muted">Min (Optionnel)</label>
                <input 
                  v-model.number="newPayment.min_amount" 
                  type="number" 
                  step="0.01"
                  class="input-field"
                  placeholder="0.00"
                />
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium text-muted">Max (Optionnel)</label>
                <input 
                  v-model.number="newPayment.max_amount" 
                  type="number" 
                  step="0.01"
                  class="input-field"
                  placeholder="0.00"
                />
              </div>
            </div>

            <!-- Description -->
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted">Description (Optionnel)</label>
              <textarea 
                v-model="newPayment.description" 
                class="input-field min-h-[100px]"
                placeholder="Détails du produit..."
              ></textarea>
            </div>

             <!-- Expiration -->
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted">Expiration</label>
               <select v-model="newPayment.expires_in_minutes" class="input-field">
                <option :value="60">1 heure</option>
                <option :value="1440">24 heures</option>
                <option :value="10080">7 jours</option>
                <option :value="-1">Jamais</option>
              </select>
            </div>

            <!-- Reusable -->
             <div class="flex items-center gap-3 p-3 rounded-lg bg-surface hover:bg-surface-hover cursor-pointer border border-secondary-200 dark:border-secondary-700">
              <input type="checkbox" id="reusable" v-model="newPayment.reusable" class="w-5 h-5 rounded border-gray-300 text-primary focus:ring-primary" />
              <label for="reusable" class="text-sm font-medium text-base cursor-pointer select-none">Réutilisable (plusieurs paiements possibles)</label>
            </div>

            <div class="flex gap-4 pt-4">
              <button type="button" class="btn-secondary w-full" @click="showCreateModal = false">
                Annuler
              </button>
              <button type="submit" class="btn-primary w-full" :disabled="creating">
                <span v-if="creating" class="loading-spinner w-5 h-5 border-white/30 border-t-white mr-2"></span>
                {{ creating ? 'Création...' : 'Créer le QR code' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- QR Code Modal -->
      <div v-if="showQRModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4" @click.self="showQRModal = false">
        <div class="glass-card w-full max-w-md animate-fade-in-up text-center overflow-hidden">
          <div class="p-6 border-b border-secondary-200 dark:border-secondary-800 relative">
            <h3 class="text-xl font-bold text-base pr-8">{{ selectedPayment?.title }}</h3>
            <button class="absolute right-4 top-4 text-muted hover:text-base text-2xl" @click="showQRModal = false">×</button>
          </div>
          
          <div class="p-8">
            <div class="bg-white p-4 rounded-xl inline-block shadow-lg mb-6">
              <img :src="qrCodeImage" alt="QR Code" class="w-48 h-48 object-contain" />
            </div>
            
            <div class="space-y-4 mb-6 text-sm">
              <div class="flex justify-between items-center py-2 border-b border-secondary-200 dark:border-secondary-800">
                <span class="text-muted">Montant</span>
                <strong class="text-lg text-base">{{ selectedPayment?.amount ? formatCurrency(selectedPayment.amount, selectedPayment.currency) : 'À définir' }}</strong>
              </div>
              <div class="flex justify-between items-center py-2 border-b border-secondary-200 dark:border-secondary-800">
                <span class="text-muted">Statut</span>
                <span class="px-2 py-0.5 rounded-full text-xs font-bold" :class="getStatusClass(selectedPayment?.status)">
                  {{ getStatusLabel(selectedPayment?.status) }}
                </span>
              </div>
               <div v-if="selectedPayment?.expires_at" class="flex justify-between items-center py-2">
                <span class="text-muted">Expire le</span>
                <span class="text-base">{{ formatDate(selectedPayment.expires_at) }}</span>
              </div>
            </div>

            <div class="grid grid-cols-3 gap-2 mb-6">
              <button class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-surface-hover text-muted hover:text-primary transition-colors" @click="copyPaymentLink">
                <span class="text-xl">📋</span>
                <span class="text-xs">Copier</span>
              </button>
              <button class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-surface-hover text-muted hover:text-primary transition-colors" @click="downloadQR">
                <span class="text-xl">⬇️</span>
                <span class="text-xs">DL</span>
              </button>
              <button class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-surface-hover text-muted hover:text-primary transition-colors" @click="sharePayment">
                <span class="text-xl">📤</span>
                <span class="text-xs">Partager</span>
              </button>
            </div>

            <div class="relative">
              <input type="text" :value="selectedPayment?.payment_link" readonly class="input-field text-xs pr-10" />
               <button @click="copyPaymentLink" class="absolute right-2 top-1/2 -translate-y-1/2 text-primary hover:scale-110 transition-transform">
                 📋
               </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

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
  if (!wallets.value || !Array.isArray(wallets.value)) return 'EUR'
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
    // Fix: correctly access data properties from response structure
    payments.value = (paymentsRes.data?.payments || paymentsRes.payments) || []
    wallets.value = (walletsRes.data?.wallets || walletsRes.wallets || walletsRes.data) || []
    
    // Ensure wallets is an array to prevent crashes
    if (!Array.isArray(wallets.value)) {
        console.warn('Wallets response is not an array:', wallets.value)
        wallets.value = []
    }

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
    
    // Handle response structure variations
    const newPaymentReq = response.data?.payment_request || response.payment_request
    const qrCode = response.data?.qr_code_base64 || response.qr_code_base64
    
    if (newPaymentReq) payments.value.unshift(newPaymentReq)
    showCreateModal.value = false
    
    // Show QR code
    selectedPayment.value = newPaymentReq
    qrCodeImage.value = qrCode
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
    alert('Erreur: ' + (error.response?.data?.error || error.message))
  } finally {
    creating.value = false
  }
}

async function showQRCode(payment) {
  selectedPayment.value = payment
  try {
    const response = await merchantApi.getQRCode(payment.id)
    qrCodeImage.value = response.data?.qr_code_base64 || response.qr_code_base64
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
  if (selectedPayment.value?.payment_link) {
      navigator.clipboard.writeText(selectedPayment.value.payment_link)
      alert('Lien copié!')
  }
}

function downloadQR() {
  const link = document.createElement('a')
  link.download = `qr-${selectedPayment.value.id}.png`
  link.href = qrCodeImage.value
  link.click()
}

function sharePayment() {
  if (navigator.share && selectedPayment.value) {
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
  if (!date) return '-'
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

function getStatusClass(status) {
    const classes = {
        pending: 'bg-warning/10 text-warning',
        paid: 'bg-success/10 text-success',
        expired: 'bg-error/10 text-error',
        cancelled: 'bg-secondary-200 text-muted'
    }
    return classes[status] || 'bg-secondary-100 text-muted'
}

definePageMeta({
    middleware: 'auth'
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Custom scrollbar for modal */
.glass-card::-webkit-scrollbar {
  width: 6px;
}
.glass-card::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}
</style>
