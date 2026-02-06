<template>
  <div class="tx-page">
    <!-- Header -->
    <div class="page-header">
      <h1>📊 Transactions</h1>
      <p>Historique de vos transactions</p>
    </div>

    <!-- Filters -->
    <div class="filters">
      <select v-model="filterType" class="filter-select">
        <option value="">Tous types</option>
        <option value="deposit">Dépôts</option>
        <option value="withdraw">Retraits</option>
        <option value="transfer">Transferts</option>
        <option value="exchange">Échanges</option>
      </select>
      <select v-model="filterPeriod" class="filter-select">
        <option value="all">Toujours</option>
        <option value="today">Aujourd'hui</option>
        <option value="week">7 jours</option>
        <option value="month">30 jours</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading && transactions.length === 0" class="state-container">
      <div class="spinner"></div>
      <p>Chargement...</p>
    </div>

    <!-- Empty -->
    <div v-else-if="filteredTransactions.length === 0" class="state-container">
      <span class="empty-icon">📜</span>
      <p>Aucune transaction</p>
    </div>

    <!-- List -->
    <div v-else class="tx-list">
      <div v-for="tx in filteredTransactions" :key="tx.id" class="tx-item" @click="openTxDetail(tx)">
        <div class="tx-icon" :class="getTypeColorClass(tx.type)">
          {{ getTypeIcon(tx.type) }}
        </div>
        <div class="tx-info">
          <div class="tx-title">{{ tx.title }}</div>
          <div class="tx-desc">{{ tx.description }}</div>
        </div>
        <div class="tx-amount" :class="{ positive: tx.amount >= 0 }">
          <span>{{ tx.amount >= 0 ? '+' : '' }}{{ formatMoney(tx.amount, tx.currency) }}</span>
          <small>{{ formatDate(tx.date) }}</small>
        </div>
      </div>

      <!-- Load More -->
      <button v-if="hasMore" @click="loadMore" class="load-more-btn">
        Charger plus ↓
      </button>
    </div>

    <!-- Transaction Detail Modal -->
    <Teleport to="body">
      <div v-if="selectedTx" class="modal-overlay" @click.self="closeTxDetail">
        <div class="modal-content">
          <button @click="closeTxDetail" class="modal-close">✕</button>
          <div class="modal-icon-lg" :class="getTypeColorClass(selectedTx.type)">
            {{ getTypeIcon(selectedTx.type) }}
          </div>
          <h3 class="modal-title">{{ selectedTx.title }}</h3>
          <p class="modal-amount" :class="{ positive: selectedTx.amount >= 0 }">
            {{ selectedTx.amount >= 0 ? '+' : '' }}{{ formatMoney(selectedTx.amount, selectedTx.currency) }}
          </p>
          
          <div class="tx-details">
            <div class="detail-row">
              <span class="detail-label">ID Transaction</span>
              <span class="detail-value id-value">{{ selectedTx.id }}</span>
            </div>
            <div class="detail-row" v-if="selectedTx.reference">
              <span class="detail-label">Référence</span>
              <span class="detail-value">{{ selectedTx.reference }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Type</span>
              <span class="detail-value">{{ getTransactionTitle(selectedTx.type) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Devise</span>
              <span class="detail-value">{{ selectedTx.currency }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Statut</span>
              <span class="detail-value status-badge" :class="selectedTx.status">{{ getStatusLabel(selectedTx.status) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Date</span>
              <span class="detail-value">{{ formatFullDate(selectedTx.date) }}</span>
            </div>
            <div class="detail-row" v-if="selectedTx.description">
              <span class="detail-label">Description</span>
              <span class="detail-value">{{ selectedTx.description }}</span>
            </div>

             <!-- Extended Details -->
             <div v-if="loadingDetails" class="detail-row justify-center py-4">
                <div class="spinner-sm"></div>
             </div>
             
             <template v-else>
                 <!-- Sender -->
                 <div class="detail-row" v-if="senderInfo">
                     <span class="detail-label">De (Remetteur)</span>
                     <div class="text-right flex flex-col items-end">
                         <span class="detail-value font-bold">{{ senderInfo.name }}</span>
                         <span v-if="senderInfo.phone" class="text-xs text-gray-400">{{ senderInfo.phone }}</span>
                         <span v-if="senderInfo.email" class="text-xs text-gray-400">{{ senderInfo.email }}</span>
                     </div>
                 </div>

                 <!-- Receiver -->
                 <div class="detail-row" v-if="receiverInfo">
                     <span class="detail-label">À (Récepteur)</span>
                     <div class="text-right flex flex-col items-end">
                          <span class="detail-value font-bold">{{ receiverInfo.name }}</span>
                          <span v-if="receiverInfo.phone" class="text-xs text-gray-400">{{ receiverInfo.phone }}</span>
                          <span v-if="receiverInfo.email" class="text-xs text-gray-400">{{ receiverInfo.email }}</span>
                     </div>
                 </div>
             </template>

             <!-- Transfer Actions (Cancel/Reverse) -->
             <div v-if="selectedTx.type === 'transfer'" class="mt-6 pt-6 border-t border-gray-700">
                 <!-- Sender Cancel (< 5 mins) -->
                 <button v-if="selectedTx.amount < 0 && canCancel(selectedTx)" @click="initiateAction('cancel')" class="w-full py-2 bg-red-500/10 hover:bg-red-500/20 text-red-500 rounded-lg text-sm font-bold border border-red-500/20 transition-colors mb-2">
                     Annuler le transfert ({{ getRemainingCancelTime(selectedTx) }})
                 </button>
                 
                 <!-- Recipient Reverse (< 1 week, completed) -->
                 <button v-if="selectedTx.amount > 0 && selectedTx.status === 'completed' && canReverse(selectedTx)" @click="initiateAction('reverse')" class="w-full py-2 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-lg text-sm font-bold border border-indigo-500/20 transition-colors mb-2">
                     Renvoyer les fonds
                 </button>
             </div>

          </div>
          
          <button @click="closeTxDetail" class="modal-btn">Fermer</button>
        </div>
      </div>
    </Teleport>

    <!-- Action Reason Modal -->
    <Teleport to="body">
        <div v-if="showActionModal" class="modal-overlay" @click.self="closeActionModal">
            <div class="modal-content">
                <h3 class="modal-title mb-2">
                    {{ actionType === 'cancel' ? 'Annuler le transfert' : 'Renvoyer les fonds' }}
                </h3>
                <p class="text-gray-400 text-sm mb-6">
                    {{ actionType === 'cancel' 
                        ? 'Les transferts peuvent être annulés uniquement dans les 5 minutes suivant l\'envoi, si le destinataire n\'a pas encore utilisé les fonds.' 
                        : 'Vous pouvez renvoyer les fonds dans un délai de 7 jours si vous ne les avez pas utilisés.' 
                    }}
                </p>

                <div class="text-left mb-6">
                    <label class="block text-sm font-medium text-gray-400 mb-2">Motif (Obligatoire)</label>
                    <textarea 
                        v-model="actionReason"
                        rows="3"
                        class="w-full bg-black/20 border border-gray-700 rounded-xl p-3 text-white focus:border-indigo-500 outline-none"
                        :placeholder="actionType === 'cancel' ? 'Erreur de montant, destinataire incorrect...' : 'Erreur de réception...'"
                    ></textarea>
                </div>

                <div class="flex gap-3">
                    <button @click="closeActionModal" class="flex-1 py-3 bg-gray-700/50 hover:bg-gray-700 text-white rounded-xl font-medium transition-colors">Retour</button>
                    <button 
                        @click="submitAction" 
                        :disabled="!actionReason.trim() || actionLoading"
                        class="flex-1 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2"
                    >
                        <span v-if="actionLoading" class="spinner-sm border-white w-4 h-4"></span>
                        {{ actionType === 'cancel' ? 'Confirmer Annulation' : 'Confirmer Renvoi' }}
                    </button>
                </div>
            </div>
        </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { walletAPI, transferAPI, userAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const transactions = ref([])
const loading = ref(false)
const filterType = ref('')
const filterPeriod = ref('all')
const offset = ref(0)
const limit = 50
const hasMore = ref(false)
const selectedTx = ref(null)

// Extended Details
const loadingDetails = ref(false)
const senderInfo = ref(null)
const receiverInfo = ref(null)

// Time Helpers
const canCancel = (tx) => {
    if (tx.status !== 'pending' && tx.status !== 'processing' && tx.status !== 'completed') return false
    const now = new Date()
    const txDate = new Date(tx.date)
    const diffMins = (now - txDate) / 1000 / 60
    return diffMins < 5
}

const getRemainingCancelTime = (tx) => {
    const now = new Date()
    const txDate = new Date(tx.date)
    const diffSecs = (now - txDate) / 1000
    const remaining = 300 - diffSecs // 5 mins * 60
    if (remaining <= 0) return '0s'
    const mins = Math.floor(remaining / 60)
    const secs = Math.floor(remaining % 60)
    return `${mins}m ${secs}s`
}

const canReverse = (tx) => {
    const now = new Date()
    const txDate = new Date(tx.date)
    const diffHours = (now - txDate) / 1000 / 60 / 60
    return diffHours < 168 // 7 days * 24h
}

const filteredTransactions = computed(() => {
  let result = [...transactions.value]
  
  if (filterType.value) {
    result = result.filter(tx => tx.type === filterType.value)
  }
  
  if (filterPeriod.value !== 'all') {
    const now = new Date()
    const startDate = new Date()
    
    switch (filterPeriod.value) {
      case 'today':
        startDate.setHours(0, 0, 0, 0)
        break
      case 'week':
        startDate.setDate(now.getDate() - 7)
        break
      case 'month':
        startDate.setMonth(now.getMonth() - 1)
        break
    }
    
    result = result.filter(tx => new Date(tx.date) >= startDate)
  }
  
  return result
})

const formatMoney = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('fr-FR', { 
    style: 'currency', 
    currency 
  }).format(Math.abs(amount))
}

const formatDate = (date) => {
  return new Intl.DateTimeFormat('fr-FR', { 
    day: '2-digit', 
    month: 'short'
  }).format(new Date(date))
}

const getTypeIcon = (type) => {
  const icons = {
    deposit: '↓',
    withdraw: '↑',
    transfer: '💸',
    exchange: '💱',
    payment: '💳',
    refund: '↩️',
    donation: '💖'
  }
  return icons[type] || '💰'
}

const getTypeColorClass = (type) => {
  const colors = {
    deposit: 'icon-green',
    withdraw: 'icon-red',
    transfer: 'icon-purple',
    exchange: 'icon-blue',
    payment: 'icon-orange',
    refund: 'icon-purple',
    donation: 'icon-purple'
  }
  return colors[type] || 'icon-gray'
}

const fetchTransactions = async () => {
  loading.value = true
  try {
    const walletsRes = await walletAPI.getAll()
    if (!walletsRes.data?.wallets) {
      transactions.value = []
      return
    }
    
    const allTransactions = []
    for (const wallet of walletsRes.data.wallets) {
      try {
        const txRes = await walletAPI.getTransactions(wallet.id, limit, offset.value)
        if (txRes.data?.transactions) {
          txRes.data.transactions.forEach(tx => {
            let txType = tx.transaction_type || 'transfer'
            if (tx.reference && tx.reference.startsWith('DON-')) txType = 'donation'
            if (tx.reference && tx.reference.startsWith('REF-')) txType = 'refund'

            allTransactions.push({
              id: tx.id,
              type: txType,
              title: getTransactionTitle(txType),
              description: tx.description || `${wallet.currency}`,
              amount: tx.from_wallet_id === wallet.id ? -tx.amount : tx.amount,
              currency: tx.currency || wallet.currency,
              date: tx.created_at,
              reference: tx.reference,
              status: tx.status
            })
          })
        }
      } catch (e) {
        console.error(`Error fetching transactions for wallet ${wallet.id}:`, e)
      }
    }
    
    allTransactions.sort((a, b) => new Date(b.date) - new Date(a.date))
    transactions.value = allTransactions
    hasMore.value = allTransactions.length >= limit
  } catch (e) {
    console.error('Error fetching transactions:', e)
  } finally {
    loading.value = false
  }
}

const getTransactionTitle = (type) => {
  const titles = {
    deposit: 'Dépôt',
    withdraw: 'Retrait',
    transfer: 'Transfert',
    exchange: 'Échange',
    payment: 'Paiement',
    refund: 'Remboursement',
    donation: 'Donation'
  }
  return titles[type] || 'Transaction'
}

const getStatusLabel = (status) => {
  const labels = {
    completed: 'Complété',
    pending: 'En attente',
    processing: 'En cours',
    failed: 'Échoué',
    cancelled: 'Annulé'
  }
  return labels[status] || status || 'Complété'
}

const formatFullDate = (date) => {
  return new Intl.DateTimeFormat('fr-FR', {
    weekday: 'long',
    year: 'numeric', 
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(date))
}

const loadTransferDetails = async (tx) => {
    loadingDetails.value = true
    senderInfo.value = null
    receiverInfo.value = null
    
    try {
        // Fetch details for transfer, donation, or refund types
        const isTransferType = ['transfer', 'donation', 'refund'].includes(tx.type) || tx.reference?.startsWith('DON-') || tx.reference?.startsWith('REF-')
        
        if (isTransferType && tx.reference) {
            try {
                let transferId = tx.id // Use TX ID for fetch if possible, or try reference logic
                
                // Fetch full transfer details (enriched by backend)
                // Assuming API can take UUID
                const res = await transferAPI.get(transferId)
                const transfer = res.data

                // Use enriched details if available
                if (transfer.sender_details) {
                    senderInfo.value = {
                        name: transfer.sender_details.name,
                        email: transfer.sender_details.email,
                        phone: transfer.sender_details.phone
                    }
                } else if (transfer.user_id) {
                     // Fallback to fetching user if not enriched (backwards compat)
                    try {
                        const uRes = await userAPI.getById(transfer.user_id)
                        const u = uRes.data
                        senderInfo.value = {
                            name: `${u.first_name} ${u.last_name}`.trim() || u.username || 'Utilisateur',
                            phone: u.phone,
                            email: u.email
                        }
                    } catch (e) { console.warn('Fetch sender failed', e) }
                }

                if (transfer.recipient_details) {
                    receiverInfo.value = {
                        name: transfer.recipient_details.name,
                        email: transfer.recipient_details.email,
                        phone: transfer.recipient_details.phone
                    }
                } else if(transfer.recipient_email || transfer.recipient_phone) {
                     receiverInfo.value = {
                         name: transfer.recipient_name || 'Contact Externe',
                         email: transfer.recipient_email,
                         phone: transfer.recipient_phone
                     }
                }
            } catch (e) {
                console.warn("Failed to fetch transfer details", e)
            }
        }
    } finally {
        loadingDetails.value = false
    }
}

const openTxDetail = (tx) => {
  selectedTx.value = tx
  loadTransferDetails(tx)
}

const closeTxDetail = () => {
  selectedTx.value = null
  senderInfo.value = null
  receiverInfo.value = null
  // Clean query param
  const url = new URL(window.location.href);
  url.searchParams.delete("id");
  url.searchParams.delete("ref");
  window.history.replaceState({}, "", url);
}

const loadMore = async () => {
  offset.value += limit
  await fetchTransactions()
}

// === Transfer Actions Logic ===
const showActionModal = ref(false)
const actionType = ref('cancel') // 'cancel' | 'reverse'
const actionReason = ref('')
const actionLoading = ref(false)

const initiateAction = (type) => {
    actionType.value = type
    actionReason.value = ''
    showActionModal.value = true
}

const closeActionModal = () => {
    showActionModal.value = false
    actionReason.value = ''
}

const submitAction = async () => {
    if (!selectedTx.value || !selectedTx.value.id) return
    
    actionLoading.value = true
    try {
        const transferId = selectedTx.value.id 
        
        if (actionType.value === 'cancel') {
            await transferAPI.cancelTransfer(transferId, actionReason.value)
            alert("Transfert annulé avec succès.")
        } else {
            await transferAPI.reverseTransfer(transferId, actionReason.value)
            alert("Fonds renvoyés avec succès.")
        }
        
        closeActionModal()
        closeTxDetail()
        
        // Reload transactions to reflect status change
        // Reset list and reload
        offset.value = 0
        transactions.value = []
        await fetchTransactions()
        
    } catch (e) {
        console.error(e)
        alert(e.response?.data?.error || "Une erreur est survenue lors de l'opération.")
    } finally {
        actionLoading.value = false
    }
}

onMounted(async () => {
  await fetchTransactions()
  
  // Check for deep link
  const refId = route.query.ref || route.query.id
  if (refId) {
      // Find in current list
      const tx = transactions.value.find(t => t.id === refId || t.reference === refId)
      if (tx) {
          openTxDetail(tx)
      } else {
         // If not found in list (maybe older?), try to construct one from Transfer API
         try {
             // If ID format is UUID, it might be Transfer ID.
             const res = await transferAPI.get(refId)
             const t = res.data
             const mockTx = {
                 id: t.id, // Using Transfer ID as ID if TX not found
                 type: 'transfer',
                 title: 'Transfert (Archive)',
                 description: t.description || 'Transfert',
                 amount: t.amount, // Direction unknown without wallet context, assume +? Or check user_id
                 currency: t.currency,
                 date: t.created_at,
                 reference: t.id,
                 status: t.status
             }
             openTxDetail(mockTx)
         } catch(e) { console.error("Deep link failed", e) }
      }
  }
})

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})
</script>

<style scoped>
.tx-page {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.page-header {
  margin-bottom: 1rem;
}

.page-header h1 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.page-header p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

.filters {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.filter-select {
  flex: 1;
  min-width: 120px;
  max-width: 100%;
  padding: 0.625rem 0.875rem;
  border-radius: 0.625rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
  cursor: pointer;
}

.filter-select:focus {
  border-color: #6366f1;
}

.state-container {
  text-align: center;
  padding: 3rem 1rem;
  color: #888;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

.spinner-sm {
  width: 1.5rem;
  height: 1.5rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  margin: 0 auto;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.tx-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  padding: 0.5rem;
}

.tx-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.75rem;
  transition: background 0.2s;
  cursor: pointer;
}

.tx-item:hover {
  background: rgba(255,255,255,0.05);
}

.tx-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
}

.icon-green { background: rgba(34, 197, 94, 0.15); color: #22c55e; }
.icon-red { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.icon-purple { background: rgba(168, 85, 247, 0.15); color: #a855f7; }
.icon-blue { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.icon-orange { background: rgba(249, 115, 22, 0.15); color: #f97316; }
.icon-gray { background: rgba(107, 114, 128, 0.15); color: #6b7280; }

.tx-info {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.tx-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tx-desc {
  font-size: 0.75rem;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tx-amount {
  text-align: right;
  flex-shrink: 0;
}

.tx-amount span {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: #fff;
}

.tx-amount.positive span {
  color: #22c55e;
}

.tx-amount small {
  font-size: 0.625rem;
  color: #666;
}

.load-more-btn {
  width: 100%;
  padding: 0.75rem;
  border: none;
  background: transparent;
  color: #6366f1;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  border-radius: 0.5rem;
  transition: background 0.2s;
}

.load-more-btn:hover {
  background: rgba(99, 102, 241, 0.1);
}

/* Desktop */
@media (min-width: 640px) {
  .page-header h1 {
    font-size: 1.5rem;
  }
  
  .filters {
    flex-wrap: nowrap;
  }
  
  .filter-select {
    flex: 0 0 auto;
    width: auto;
  }
  
  .tx-list {
    padding: 0.75rem;
  }
  
  .tx-item {
    gap: 1rem;
    padding: 1rem;
  }
  
  .tx-icon {
    width: 3rem;
    height: 3rem;
    font-size: 1.25rem;
  }
  
  .tx-title {
    font-size: 1rem;
  }
  
  .tx-desc {
    font-size: 0.875rem;
  }
  
  .tx-amount span {
    font-size: 1rem;
  }
}

/* Transaction Detail Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 1rem;
}

.modal-content {
  background: linear-gradient(180deg, #1e1e2e, #15151f);
  border-radius: 1.5rem;
  padding: 1.5rem;
  max-width: 420px;
  width: 100%;
  text-align: center;
  position: relative;
  border: 1px solid rgba(255,255,255,0.1);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  max-height: 90vh;
  overflow-y: auto;
}

.modal-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: rgba(255,255,255,0.1);
  border: none;
  color: #888;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  cursor: pointer;
  font-size: 0.875rem;
}

.modal-close:hover {
  background: rgba(255,255,255,0.2);
  color: #fff;
}

.modal-icon-lg {
  width: 4rem;
  height: 4rem;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  margin: 0 auto 1rem;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
  margin-bottom: 0.5rem;
}

.modal-amount {
  font-size: 1.75rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 1.5rem;
}

.modal-amount.positive {
  color: #22c55e;
}

.tx-details {
  background: rgba(0,0,0,0.3);
  border-radius: 1rem;
  padding: 1rem;
  margin-bottom: 1.5rem;
  text-align: left;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 0.5rem 0;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 0.75rem;
  color: #888;
  flex-shrink: 0;
}

.detail-value {
  font-size: 0.875rem;
  color: #fff;
  text-align: right;
  word-break: break-all;
  max-width: 60%;
}

.id-value {
  font-family: monospace;
  font-size: 0.75rem;
  color: #6366f1;
}

.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-badge.completed { background: rgba(34,197,94,0.2); color: #22c55e; }
.status-badge.pending { background: rgba(245,158,11,0.2); color: #f59e0b; }
.status-badge.processing { background: rgba(59,130,246,0.2); color: #3b82f6; }
.status-badge.failed { background: rgba(239,68,68,0.2); color: #ef4444; }
.status-badge.cancelled { background: rgba(107,114,128,0.2); color: #6b7280; }

.modal-btn {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
  width: 100%;
}

.modal-btn:hover {
  transform: scale(1.02);
}
</style>
