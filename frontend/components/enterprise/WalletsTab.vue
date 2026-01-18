<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <WalletIcon class="w-6 h-6 text-primary-500" />
            Portefeuilles
          </h2>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            Gérez les portefeuilles de votre entreprise
          </p>
        </div>
        
        <button 
          @click="showCreateWallet = true" 
          class="px-5 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl font-medium transition-all shadow-lg shadow-primary-500/25 flex items-center gap-2">
          <PlusIcon class="w-5 h-5" />
          Nouveau Portefeuille
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600"></div>
    </div>

    <!-- Wallets Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="wallet in wallets" :key="wallet?.id || wallet?._id"
        class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 relative overflow-hidden group hover:shadow-lg transition-all">
        
        <!-- Default Badge -->
        <div v-if="enterprise?.default_wallet_id && wallet?.id === enterprise.default_wallet_id" 
          class="absolute top-3 right-3 px-2 py-1 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 text-xs font-medium rounded-full">
          Par défaut
        </div>

        <!-- Currency Icon -->
        <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center text-white text-xl font-bold mb-4 shadow-lg">
          {{ getCurrencySymbol(wallet.currency) }}
        </div>

        <!-- Balance -->
        <div class="mb-4">
          <p class="text-sm text-gray-500 dark:text-gray-400">Solde disponible</p>
          <p class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ formatAmount(wallet.balance, wallet.currency) }}
          </p>
        </div>

        <!-- Info -->
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-500">Devise</span>
            <span class="font-medium text-gray-900 dark:text-white">{{ wallet.currency }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500">Type</span>
            <span class="font-medium text-gray-900 dark:text-white">{{ wallet.wallet_type || 'Business' }}</span>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-2 mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
          <button @click="viewTransactions(wallet)" 
            class="flex-1 px-3 py-2 text-sm font-medium text-primary-600 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded-lg transition-colors">
            📋 Transactions
          </button>
          <button @click="initiateSend(wallet)" 
            class="flex-1 px-3 py-2 text-sm font-medium text-green-600 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-lg transition-colors">
            ↗️ Envoyer
          </button>
        </div>
        <div class="flex gap-2 mt-2">
          <button v-if="!enterprise?.default_wallet_id || wallet?.id !== enterprise.default_wallet_id" 
            @click="setAsDefault(wallet)"
            class="flex-1 px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 dark:text-gray-400 dark:hover:bg-gray-700 rounded-lg transition-colors">
            ⭐ Définir par défaut
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!wallets.length" class="col-span-full text-center py-12">
        <WalletIcon class="w-16 h-16 mx-auto text-gray-300 dark:text-gray-600 mb-4" />
        <p class="text-gray-500 dark:text-gray-400">Aucun portefeuille trouvé</p>
        <p class="text-sm text-gray-400 mt-2">Créez un nouveau portefeuille pour commencer</p>
      </div>
    </div>

    <!-- Create Wallet Modal -->
    <Teleport to="body">
      <div v-if="showCreateWallet" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Nouveau Portefeuille</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Devise</label>
              <select v-model="newWallet.currency" 
                class="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white">
                <option value="XOF">XOF - Franc CFA</option>
                <option value="EUR">EUR - Euro</option>
                <option value="USD">USD - Dollar</option>
                <option value="XAF">XAF - Franc CFA (BEAC)</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nom (optionnel)</label>
              <input v-model="newWallet.name" type="text" placeholder="Ex: Caisse principale"
                class="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button @click="showCreateWallet = false" 
              class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
              Annuler
            </button>
            <button @click="createWallet" :disabled="isCreating"
              class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50">
              {{ isCreating ? 'Création...' : 'Créer' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Transactions Modal -->
    <Teleport to="body">
      <div v-if="showTransactions" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-2xl max-h-[80vh] overflow-hidden">
          <div class="p-6 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">
              Transactions - {{ selectedWallet?.currency }}
            </h3>
            <button @click="showTransactions = false" class="text-gray-400 hover:text-gray-600 text-2xl">&times;</button>
          </div>
          
          <div v-if="loadingTransactions" class="p-12 text-center">
            <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600 mx-auto"></div>
          </div>
          
          <div v-else class="divide-y divide-gray-100 dark:divide-gray-700 max-h-96 overflow-y-auto">
            <div v-for="tx in transactions" :key="tx.id" class="p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div :class="['w-10 h-10 rounded-full flex items-center justify-center text-white text-lg',
                    tx.type === 'credit' || tx.type === 'deposit' ? 'bg-green-500' : 'bg-red-500']">
                    {{ tx.type === 'credit' || tx.type === 'deposit' ? '↓' : '↑' }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900 dark:text-white">{{ tx.description || tx.type }}</p>
                    <p class="text-xs text-gray-500">{{ formatDate(tx.created_at) }}</p>
                  </div>
                </div>
                <div :class="['font-bold', tx.type === 'credit' || tx.type === 'deposit' ? 'text-green-600' : 'text-red-600']">
                  {{ tx.type === 'credit' || tx.type === 'deposit' ? '+' : '-' }}{{ formatAmount(tx.amount, selectedWallet?.currency) }}
                </div>
              </div>
            </div>
            
            <div v-if="!transactions.length" class="p-12 text-center text-gray-500">
              Aucune transaction
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Send Money Modal (Transfer-Style) -->
    <Teleport to="body">
      <div v-if="showSendModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 overflow-y-auto">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-lg my-8">
          <!-- Header -->
          <div class="p-6 border-b border-gray-100 dark:border-gray-700 flex justify-between items-center">
            <div>
              <h3 class="text-lg font-bold text-gray-900 dark:text-white">
                💸 Envoyer de l'argent
              </h3>
              <p class="text-sm text-gray-500 mt-1">Depuis le portefeuille {{ selectedWallet?.currency }}</p>
            </div>
            <button @click="showSendModal = false" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
              ✕
            </button>
          </div>

          <div class="p-6 space-y-6">
            <!-- Balance Display -->
            <div class="p-4 bg-primary-50 dark:bg-primary-900/20 rounded-xl">
              <p class="text-sm text-gray-500">Solde disponible</p>
              <p class="text-2xl font-bold text-primary-600">{{ formatAmount(selectedWallet?.balance, selectedWallet?.currency) }}</p>
            </div>

            <!-- Amount Input -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Montant à envoyer</label>
              <div class="flex items-center gap-3 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-xl border border-gray-200 dark:border-gray-600 focus-within:ring-2 focus-within:ring-primary-500">
                <input v-model.number="sendForm.amount" type="number" step="0.01" min="0" 
                  placeholder="0.00"
                  class="flex-1 text-2xl font-bold bg-transparent outline-none text-gray-900 dark:text-white"
                />
                <span class="text-lg font-bold text-primary-600 bg-primary-100 dark:bg-primary-900/30 px-3 py-1 rounded-lg">
                  {{ selectedWallet?.currency }}
                </span>
              </div>
              <p v-if="sendForm.amount > (selectedWallet?.balance || 0)" class="text-xs text-red-500 mt-1">Solde insuffisant</p>
              <p v-else class="text-xs text-gray-400 mt-1">Disponible: {{ formatAmount(selectedWallet?.balance, selectedWallet?.currency) }}</p>
            </div>

            <!-- Recipient Input with Lookup -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Destinataire (Email ou Téléphone)</label>
              <div class="flex gap-2">
                <input v-model="sendForm.recipient" type="text" 
                  placeholder="ex: ami@email.com ou +225..."
                  class="flex-1 px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 dark:text-white"
                  @blur="lookupRecipient"
                />
                <button type="button" @click="lookupRecipient" 
                  class="px-4 rounded-xl bg-gray-100 dark:bg-gray-700 hover:bg-primary-100 dark:hover:bg-primary-900/30 transition-colors border border-gray-200 dark:border-gray-600">
                  <span v-if="lookupLoading" class="animate-spin">⟳</span>
                  <span v-else>🔍</span>
                </button>
              </div>
              <!-- Lookup Result -->
              <div v-if="lookupResult" class="mt-3 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center font-bold shadow-lg">
                  {{ lookupResult.first_name?.[0] || 'U' }}
                </div>
                <div>
                  <p class="font-bold text-gray-900 dark:text-white">{{ lookupResult.first_name }} {{ lookupResult.last_name }}</p>
                  <p class="text-xs text-green-600 flex items-center gap-1">✓ Utilisateur vérifié</p>
                </div>
              </div>
              <p v-if="lookupError" class="text-xs text-red-500 mt-2">{{ lookupError }}</p>
            </div>

            <!-- Description -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Note (Optionnel)</label>
              <input v-model="sendForm.description" type="text" placeholder="Ex: Paiement fournisseur"
                class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 dark:text-white"
              />
            </div>

            <!-- Summary -->
            <div class="p-4 rounded-xl bg-gray-50 dark:bg-gray-700/50 border border-gray-200 dark:border-gray-600">
              <h4 class="font-bold text-gray-900 dark:text-white mb-3 flex items-center gap-2">📝 Résumé</h4>
              <div class="space-y-2 text-sm">
                <div class="flex justify-between">
                  <span class="text-gray-500">Montant</span>
                  <span class="font-medium text-gray-900 dark:text-white">{{ formatAmount(sendForm.amount || 0, selectedWallet?.currency) }}</span>
                </div>
                <div class="flex justify-between text-primary-600" v-if="estimatedFee > 0">
                  <span>Frais estimés</span>
                  <span class="font-medium">+ {{ formatAmount(estimatedFee, selectedWallet?.currency) }}</span>
                </div>
                <div class="border-t border-gray-200 dark:border-gray-600 pt-2 flex justify-between">
                  <span class="font-bold text-gray-900 dark:text-white">Total à débiter</span>
                  <span class="font-bold text-lg text-gray-900 dark:text-white">{{ formatAmount((sendForm.amount || 0) + estimatedFee, selectedWallet?.currency) }}</span>
                </div>
              </div>
            </div>

            <!-- Multi-Admin Notice -->
            <div class="p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl text-sm text-blue-700 dark:text-blue-300">
              ℹ️ Cette transaction nécessitera l'approbation d'un autre administrateur avant exécution.
            </div>

            <!-- Error -->
            <div v-if="sendError" class="p-3 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg text-sm">
              {{ sendError }}
            </div>
          </div>

          <!-- Footer -->
          <div class="p-6 border-t border-gray-100 dark:border-gray-700 flex gap-3">
            <button @click="showSendModal = false" 
              class="flex-1 px-4 py-3 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl font-medium">
              Annuler
            </button>
            <button @click="confirmSend" 
              :disabled="isSending || !sendForm.recipient || !sendForm.amount || sendForm.amount > (selectedWallet?.balance || 0)"
              class="flex-1 px-4 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-xl font-bold hover:from-green-700 hover:to-green-800 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2">
              <span v-if="isSending" class="animate-spin">⟳</span>
              {{ isSending ? 'Création...' : 'Confirmer →' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- PIN Modal for transaction confirmation -->
    <PinModal 
      :isOpen="showPinModal"
      title="Confirmer la transaction"
      description="Entrez votre code PIN pour valider et initier le processus d'approbation."
      @update:isOpen="showPinModal = $event"
      @success="onPinSuccess"
      @close="showPinModal = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { WalletIcon, PlusIcon } from '@heroicons/vue/24/outline'
import { walletAPI, enterpriseAPI, transferAPI, userAPI } from '@/composables/useApi'
import PinModal from '@/components/common/PinModal.vue'

const props = defineProps({
  enterprise: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

const wallets = ref([])
const isLoading = ref(true)
const showCreateWallet = ref(false)
const isCreating = ref(false)

// Transactions state
const showTransactions = ref(false)
const selectedWallet = ref(null)
const transactions = ref([])
const loadingTransactions = ref(false)

// Send money state
const showSendModal = ref(false)
const isSending = ref(false)
const sendError = ref('')
const sendForm = ref({
  recipient: '',
  amount: 0,
  description: ''
})

const newWallet = ref({
  currency: 'XOF',
  name: ''
})

// Recipient lookup state
const lookupLoading = ref(false)
const lookupResult = ref(null)
const lookupError = ref('')

// PIN modal state
const showPinModal = ref(false)

// Estimated fee (0 for internal P2P transfers)
const estimatedFee = computed(() => 0)

const getCurrencySymbol = (currency) => {
  const symbols = { XOF: 'F', XAF: 'F', EUR: '€', USD: '$', GBP: '£' }
  return symbols[currency] || currency?.charAt(0) || '?'
}

const formatAmount = (amount, currency) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency || 'XOF',
    minimumFractionDigits: 0
  }).format(amount || 0)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadWallets = async () => {
  isLoading.value = true
  try {
    const { data } = await walletAPI.getAll()
    console.log('Wallets response:', data)
    
    let allWallets = []
    if (Array.isArray(data)) {
      allWallets = data
    } else if (data?.wallets && Array.isArray(data.wallets)) {
      allWallets = data.wallets
    }
    
    const walletIds = [props.enterprise.default_wallet_id, ...(props.enterprise.wallet_ids || [])].filter(Boolean)
    
    if (walletIds.length > 0) {
      wallets.value = allWallets.filter(w => walletIds.includes(w.id))
    } else {
      wallets.value = allWallets.filter(w => w.wallet_type === 'fiat')
    }
    
    console.log('Filtered wallets:', wallets.value)
  } catch (error) {
    console.error('Failed to load wallets:', error)
    wallets.value = []
  } finally {
    isLoading.value = false
  }
}

const createWallet = async () => {
  isCreating.value = true
  try {
    const { data } = await walletAPI.create({
      currency: newWallet.value.currency,
      wallet_type: 'fiat',
      name: newWallet.value.name || `${props.enterprise.name} - ${newWallet.value.currency}`
    })
    
    const updatedWalletIds = [...(props.enterprise.wallet_ids || []), data.id]
    const updates = {
      ...props.enterprise,
      wallet_ids: updatedWalletIds
    }
    
    if (!props.enterprise.default_wallet_id) {
      updates.default_wallet_id = data.id
    }
    
    await enterpriseAPI.update(props.enterprise.id, updates)
    
    emit('update', updates)
    showCreateWallet.value = false
    newWallet.value = { currency: 'XOF', name: '' }
    await loadWallets()
  } catch (error) {
    console.error('Failed to create wallet:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || error.message || 'Erreur inconnue'
    alert(`Erreur lors de la création du portefeuille:\n${errorMsg}`)
  } finally {
    isCreating.value = false
  }
}

const setAsDefault = async (wallet) => {
  try {
    await enterpriseAPI.update(props.enterprise.id, {
      ...props.enterprise,
      default_wallet_id: wallet.id
    })
    emit('update', { ...props.enterprise, default_wallet_id: wallet.id })
  } catch (error) {
    console.error('Failed to set default wallet:', error)
  }
}

const viewTransactions = async (wallet) => {
  selectedWallet.value = wallet
  showTransactions.value = true
  loadingTransactions.value = true
  
  try {
    // Fetch transactions for this wallet
    const { data } = await walletAPI.getTransactions(wallet.id)
    transactions.value = data?.transactions || data || []
  } catch (error) {
    console.error('Failed to load transactions:', error)
    transactions.value = []
  } finally {
    loadingTransactions.value = false
  }
}

const initiateSend = (wallet) => {
  selectedWallet.value = wallet
  sendForm.value = { recipient: '', amount: 0, description: '' }
  sendError.value = ''
  lookupResult.value = null
  lookupError.value = ''
  showSendModal.value = true
}

// Lookup recipient user
const lookupRecipient = async () => {
  if (!sendForm.value.recipient || sendForm.value.recipient.length < 3) return
  
  lookupLoading.value = true
  lookupError.value = ''
  lookupResult.value = null
  
  try {
    const isEmail = sendForm.value.recipient.includes('@')
    const query = isEmail 
      ? { email: sendForm.value.recipient } 
      : { phone: sendForm.value.recipient }
    const res = await userAPI.lookup(query)
    lookupResult.value = res.data
  } catch (e) {
    lookupError.value = 'Utilisateur introuvable'
  } finally {
    lookupLoading.value = false
  }
}

// Step 1: Show PIN modal to confirm
const confirmSend = () => {
  if (!selectedWallet.value || !sendForm.value.recipient || !sendForm.value.amount) return
  showPinModal.value = true
}

// Step 2: After PIN validated, initiate the approval request
const onPinSuccess = async (encryptedPin) => {
  showPinModal.value = false
  
  if (!selectedWallet.value || !sendForm.value.recipient || !sendForm.value.amount) return
  
  isSending.value = true
  sendError.value = ''
  
  try {
    // Use the enterprise approval system instead of direct transfer
    // This creates a pending action that requires admin approval
    const enterpriseId = props.enterprise?.id || props.enterprise?._id
    if (!enterpriseId) {
      throw new Error('Enterprise ID not found')
    }
    
    const { data } = await enterpriseAPI.initiateAction(enterpriseId, {
      action_type: 'TRANSACTION',
      action_name: `Transfert ${formatAmount(sendForm.value.amount, selectedWallet.value.currency)}`,
      description: sendForm.value.description || `Transfert vers ${sendForm.value.recipient}`,
      amount: sendForm.value.amount,
      currency: selectedWallet.value.currency,
      payload: {
        recipient_identifier: sendForm.value.recipient,
        amount: sendForm.value.amount,
        currency: selectedWallet.value.currency,
        description: sendForm.value.description || `Transfert ${props.enterprise?.name || 'Entreprise'}`,
        source_wallet_id: selectedWallet.value.id,
        initiator_pin: encryptedPin // PIN of the initiator
      }
    })
    
    showSendModal.value = false
    
    if (data?.requires_approval || data?.approval) {
      // Redirect to approval page - admin can approve from there
      alert('Transaction créée! Elle nécessite l\'approbation d\'un autre administrateur.')
      // Navigate to the approval confirmation page
      const approvalId = data.approval?.id || data.approval?._id
      if (approvalId) {
        navigateTo(`/enterprise/approvals/${approvalId}`)
      }
    } else {
      // No approval needed - execute directly
      alert('Transfert initié avec succès!')
      await loadWallets()
    }
  } catch (error) {
    console.error('Failed to initiate transaction:', error)
    sendError.value = error.response?.data?.error || error.response?.data?.message || error.message || 'Erreur lors de l\'initiation'
  } finally {
    isSending.value = false
  }
}

onMounted(loadWallets)
</script>
