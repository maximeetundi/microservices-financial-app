<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">💼 Espace Marchand</h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">Recevez des paiements via QR code</p>
        </div>
        <button class="btn-primary flex items-center gap-2 shadow-lg shadow-primary-500/30 hover:shadow-primary-500/50 transition-all" @click="showCreateModal = true">
          <span class="text-xl">+</span>
          Créer un paiement
        </button>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="glass-card p-6 flex items-center gap-4 hover:scale-105 transition-transform duration-300 border border-gray-100 dark:border-white/10 bg-surface dark:bg-white/5">
          <div class="p-4 rounded-xl bg-blue-100 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 text-2xl shadow-sm">💰</div>
          <div>
            <span class="block text-3xl font-bold text-gray-900 dark:text-white">{{ formatCurrency(stats.totalAmount) }}</span>
            <span class="text-sm text-gray-500 dark:text-gray-400 font-medium">Total reçu</span>
          </div>
        </div>
        <div class="glass-card p-6 flex items-center gap-4 hover:scale-105 transition-transform duration-300 border border-gray-100 dark:border-white/10 bg-surface dark:bg-white/5">
          <div class="p-4 rounded-xl bg-purple-100 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400 text-2xl shadow-sm">📊</div>
          <div>
            <span class="block text-3xl font-bold text-gray-900 dark:text-white">{{ stats.totalPayments }}</span>
            <span class="text-sm text-gray-500 dark:text-gray-400 font-medium">Paiements</span>
          </div>
        </div>
        <div class="glass-card p-6 flex items-center gap-4 hover:scale-105 transition-transform duration-300 border border-gray-100 dark:border-white/10 bg-surface dark:bg-white/5">
          <div class="p-4 rounded-xl bg-amber-100 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400 text-2xl shadow-sm">⏳</div>
          <div>
            <span class="block text-3xl font-bold text-gray-900 dark:text-white">{{ pendingPayments.length }}</span>
            <span class="text-sm text-gray-500 dark:text-gray-400 font-medium">En attente</span>
          </div>
        </div>
      </div>

      <!-- Payment Requests List -->
      <div class="glass-card bg-surface dark:bg-slate-900 border border-gray-100 dark:border-white/10 shadow-xl dark:shadow-none overflow-hidden rounded-2xl">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-6 border-b border-gray-100 dark:border-gray-800">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Demandes de paiement</h2>
          <div class="flex gap-2 bg-gray-100 dark:bg-gray-800 p-1 rounded-xl">
            <button 
              v-for="tab in ['pending', 'completed']" 
              :key="tab"
              class="px-4 py-2 rounded-lg text-sm font-bold transition-all"
              :class="activeTab === tab ? 'bg-white dark:bg-slate-700 shadow-sm text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'"
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
            class="group flex flex-col sm:flex-row sm:items-center gap-4 p-4 rounded-xl hover:bg-gray-50 dark:hover:bg-white/5 transition-colors cursor-pointer border border-transparent hover:border-gray-200 dark:hover:border-gray-700"
            @click="showPaymentDetails(payment)"
          >
            <div class="flex-1">
              <div class="font-bold text-lg text-gray-900 dark:text-white mb-1">{{ payment.title }}</div>
              <div class="flex items-center gap-3 text-xs text-gray-500 dark:text-gray-400 font-medium">
                <span class="px-2.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-700">
                  {{ getTypeLabel(payment.type) }}
                </span>
                <span>{{ formatDate(payment.created_at) }}</span>
              </div>
            </div>
            
            <div class="text-right">
              <div v-if="payment.amount" class="text-xl font-bold text-gray-900 dark:text-white">
                {{ formatCurrency(payment.amount, payment.currency) }}
              </div>
              <div v-else class="text-gray-500 dark:text-gray-400 italic font-medium">Variable</div>
            </div>

            <div class="px-3 py-1 rounded-full text-xs font-bold" :class="getStatusClass(payment.status)">
              {{ getStatusLabel(payment.status) }}
            </div>

            <button class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-400 transition-colors" @click.stop="showQRCode(payment)">
              <span class="text-xl">📱</span>
            </button>
          </div>

          <div v-if="filteredPayments.length === 0" class="text-center py-16">
            <span class="text-6xl mb-6 block opacity-50 grayscale">📭</span>
            <p class="text-gray-500 dark:text-gray-400 text-lg mb-6 font-medium">Aucune demande de paiement</p>
            <button class="btn-primary" @click="showCreateModal = true">
              Créer votre première demande
            </button>
          </div>
        </div>
      </div>

      <!-- Create Payment Modal - REFACTORED FOR BEAUTY -->
      <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 transition-opacity" @click.self="showCreateModal = false">
        <div class="w-full max-w-lg max-h-[85vh] overflow-hidden rounded-2xl bg-white dark:bg-[#1a1b26] shadow-2xl animate-fade-in-up flex flex-col">
          <!-- Modal Header -->
          <div class="flex justify-between items-center p-6 border-b border-gray-100 dark:border-gray-800 bg-white dark:bg-[#1a1b26] sticky top-0 z-10">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">Nouvelle demande</h3>
            <button class="w-8 h-8 flex items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800 text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-lg" @click="showCreateModal = false">×</button>
          </div>
          
          <!-- Modal Body - Scrollable -->
          <div class="p-6 overflow-y-auto custom-scrollbar flex-1 bg-gray-50/50 dark:bg-[#1a1b26]">
            <form @submit.prevent="createPayment" class="space-y-6">
              <!-- Payment Type -->
              <div class="space-y-3">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Type de paiement</label>
                <div class="grid grid-cols-2 gap-4">
                  <button 
                    type="button"
                    class="p-4 rounded-xl border-2 transition-all flex flex-col items-center gap-2 relative overflow-hidden group"
                    :class="newPayment.type === 'fixed' 
                      ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 shadow-md ring-2 ring-indigo-500/20' 
                      : 'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 hover:border-indigo-300 dark:hover:border-indigo-700'"
                    @click="newPayment.type = 'fixed'"
                  >
                    <span class="text-3xl mb-1">🏷️</span>
                    <span class="font-bold">Prix fixe</span>
                    <div v-if="newPayment.type === 'fixed'" class="absolute inset-0 bg-indigo-500/5 dark:bg-indigo-500/10 pointer-events-none"></div>
                  </button>
                  <button 
                    type="button"
                    class="p-4 rounded-xl border-2 transition-all flex flex-col items-center gap-2 relative overflow-hidden group"
                    :class="newPayment.type === 'variable' 
                      ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400 shadow-md ring-2 ring-purple-500/20' 
                      : 'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 hover:border-purple-300 dark:hover:border-purple-700'"
                    @click="newPayment.type = 'variable'"
                  >
                    <span class="text-3xl mb-1">🎁</span>
                    <span class="font-bold">Variable</span>
                     <div v-if="newPayment.type === 'variable'" class="absolute inset-0 bg-purple-500/5 dark:bg-purple-500/10 pointer-events-none"></div>
                  </button>
                </div>
              </div>

              <!-- Wallet Selection -->
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Portefeuille de réception</label>
                <div class="relative">
                   <select v-model="newPayment.wallet_id" class="input-premium w-full p-3 rounded-xl border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white" required>
                    <option value="" disabled selected>Sélectionner...</option>
                    <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                      {{ wallet.currency }} - {{ formatCurrency(wallet.balance, wallet.currency) }}
                    </option>
                  </select>
                   <div class="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-gray-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                  </div>
                </div>
              </div>

              <!-- Title -->
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Titre</label>
                <input 
                  v-model="newPayment.title" 
                  type="text" 
                  class="input-premium w-full p-3 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white placeholder-gray-400"
                  placeholder="Ex: iPhone 15 Pro"
                  required
                />
              </div>

              <!-- Amount (for fixed) -->
              <div v-if="newPayment.type === 'fixed'" class="space-y-2 animate-fade-in-up">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Montant</label>
                <div class="relative">
                  <input 
                    v-model.number="newPayment.amount" 
                    type="number" 
                    step="0.01"
                    min="0"
                    class="input-premium w-full p-3 pr-16 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white font-mono text-lg"
                    placeholder="0.00"
                    required
                  />
                  <span class="absolute right-4 top-1/2 -translate-y-1/2 font-bold text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">{{ selectedCurrency }}</span>
                </div>
              </div>

              <!-- Min/Max (for variable) -->
              <div v-if="newPayment.type === 'variable'" class="grid grid-cols-2 gap-4 animate-fade-in-up">
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Min (Optionnel)</label>
                  <input 
                    v-model.number="newPayment.min_amount" 
                    type="number" 
                    step="0.01"
                    class="input-premium w-full p-3 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white"
                    placeholder="0.00"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Max (Optionnel)</label>
                  <input 
                    v-model.number="newPayment.max_amount" 
                    type="number" 
                    step="0.01"
                    class="input-premium w-full p-3 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white"
                    placeholder="0.00"
                  />
                </div>
              </div>

              <!-- Description -->
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Description (Optionnel)</label>
                <textarea 
                  v-model="newPayment.description" 
                  class="input-premium w-full p-3 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white min-h-[100px] resize-none"
                  placeholder="Détails du produit..."
                ></textarea>
              </div>

               <!-- Expiration -->
              <div class="space-y-2">
                <label class="text-sm font-bold text-gray-700 dark:text-gray-300">Expiration</label>
                <div class="relative">
                 <select v-model="newPayment.expires_in_minutes" class="input-premium w-full p-3 rounded-xl border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 focus:ring-2 focus:ring-indigo-500/50 outline-none transition-all dark:text-white appearance-none">
                  <option :value="60">1 heure</option>
                  <option :value="1440">24 heures</option>
                  <option :value="10080">7 jours</option>
                  <option :value="-1">Jamais</option>
                </select>
                <div class="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-gray-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                  </div>
                </div>
              </div>

              <!-- Reusable -->
               <div class="flex items-center gap-3 p-4 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:border-indigo-300 transition-colors cursor-pointer">
                <input type="checkbox" id="reusable" v-model="newPayment.reusable" class="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 cursor-pointer" />
                <label for="reusable" class="text-sm font-medium text-gray-700 dark:text-gray-300 cursor-pointer select-none flex-1">Réutilisable (plusieurs paiements possibles)</label>
              </div>
            </form>
          </div>
          
           <!-- Modal Footer -->
          <div class="p-6 border-t border-gray-100 dark:border-gray-800 bg-white dark:bg-[#1a1b26] sticky bottom-0 z-10 flex gap-4">
            <button type="button" class="flex-1 py-3 px-4 rounded-xl font-bold text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors" @click="showCreateModal = false">
              Annuler
            </button>
            <button type="submit" class="flex-1 py-3 px-4 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 shadow-lg shadow-indigo-500/30 transition-all flex items-center justify-center" :disabled="creating" @click="createPayment">
              <span v-if="creating" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2"></span>
              {{ creating ? 'Création...' : 'Créer le QR code' }}
            </button>
          </div>
        </div>
      </div>

      <!-- QR Code Modal -->
      <div v-if="showQRModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in-up" @click.self="showQRModal = false">
        <div class="w-full max-w-md bg-white dark:bg-[#1a1b26] rounded-2xl shadow-2xl overflow-hidden relative">
          <div class="p-6 border-b border-gray-100 dark:border-gray-800 relative">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white pr-8">{{ selectedPayment?.title }}</h3>
            <button class="absolute right-4 top-4 text-gray-400 hover:text-gray-600 dark:hover:text-white text-2xl" @click="showQRModal = false">×</button>
          </div>
          
          <div class="p-8 flex flex-col items-center">
            <div class="bg-white p-4 rounded-2xl inline-block shadow-lg mb-8 border border-gray-100">
              <img :src="qrCodeImage" alt="QR Code" class="w-56 h-56 object-contain" />
            </div>
            
            <div class="w-full space-y-4 mb-8 text-sm bg-gray-50 dark:bg-gray-800/50 p-4 rounded-xl">
              <div class="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span class="text-gray-500 dark:text-gray-400">Montant</span>
                <strong class="text-lg text-gray-900 dark:text-white">{{ selectedPayment?.amount ? formatCurrency(selectedPayment.amount, selectedPayment.currency) : 'À définir' }}</strong>
              </div>
              <div class="flex justify-between items-center py-2 border-b border-gray-200 dark:border-gray-700">
                <span class="text-gray-500 dark:text-gray-400">Statut</span>
                <span class="px-2.5 py-1 rounded-full text-xs font-bold capitalize" :class="getStatusClass(selectedPayment?.status)">
                  {{ getStatusLabel(selectedPayment?.status) }}
                </span>
              </div>
               <div v-if="selectedPayment?.expires_at" class="flex justify-between items-center py-2">
                <span class="text-gray-500 dark:text-gray-400">Expire le</span>
                <span class="text-gray-900 dark:text-white font-medium">{{ formatDate(selectedPayment.expires_at) }}</span>
              </div>
            </div>

            <div class="grid grid-cols-3 gap-3 w-full mb-6">
              <button class="flex flex-col items-center gap-1 p-3 rounded-xl bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors" @click="copyPaymentLink">
                <span class="text-xl">📋</span>
                <span class="text-xs font-bold">Copier</span>
              </button>
              <button class="flex flex-col items-center gap-1 p-3 rounded-xl bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors" @click="downloadQR">
                <span class="text-xl">⬇️</span>
                <span class="text-xs font-bold">DL PNG</span>
              </button>
              <button class="flex flex-col items-center gap-1 p-3 rounded-xl bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors" @click="sharePayment">
                <span class="text-xl">📤</span>
                <span class="text-xs font-bold">Partager</span>
              </button>
            </div>

            <div class="relative w-full">
              <input type="text" :value="selectedPayment?.payment_link" readonly class="w-full bg-gray-50 dark:bg-gray-800 border-none rounded-lg text-xs py-3 pl-3 pr-10 text-gray-500 font-mono" />
               <button @click="copyPaymentLink" class="absolute right-2 top-1/2 -translate-y-1/2 text-primary-600 hover:scale-110 transition-transform">
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
        pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
        paid: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
        expired: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
        cancelled: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-400'
    }
    return classes[status] || 'bg-gray-100 text-gray-800'
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

/* Force Visible Scrollbar */
.custom-scrollbar::-webkit-scrollbar {
  width: 8px; /* Force width */
  display: block; /* Ensure display */
}

.custom-scrollbar::-webkit-scrollbar-track {
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}
.dark .custom-scrollbar::-webkit-scrollbar-track {
  background-color: rgba(255, 255, 255, 0.05);
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1; /* Visible slate color by default */
  border-radius: 4px;
  border: 2px solid transparent;
  background-clip: content-box;
}

.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #475569; /* Darker slate for dark mode */
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: #94a3b8;
}

/* Styling for select arrow */
select {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}
</style>
