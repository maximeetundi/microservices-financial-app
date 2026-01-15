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
      <div v-for="wallet in wallets" :key="wallet.id"
        class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 relative overflow-hidden group hover:shadow-lg transition-all">
        
        <!-- Default Badge -->
        <div v-if="wallet.id === enterprise.default_wallet_id" 
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
            Transactions
          </button>
          <button v-if="wallet.id !== enterprise.default_wallet_id" 
            @click="setAsDefault(wallet)"
            class="flex-1 px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 dark:text-gray-400 dark:hover:bg-gray-700 rounded-lg transition-colors">
            Définir par défaut
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!wallets.length" class="col-span-full text-center py-12">
        <WalletIcon class="w-16 h-16 mx-auto text-gray-300 dark:text-gray-600 mb-4" />
        <p class="text-gray-500 dark:text-gray-400">Aucun portefeuille trouvé</p>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { WalletIcon, PlusIcon } from '@heroicons/vue/24/outline'
import { walletAPI, enterpriseAPI } from '@/composables/useApi'

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

const newWallet = ref({
  currency: 'XOF',
  name: ''
})

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

const loadWallets = async () => {
  isLoading.value = true
  try {
    // Load wallets by IDs stored in enterprise
    const walletIds = [props.enterprise.default_wallet_id, ...(props.enterprise.wallet_ids || [])].filter(Boolean)
    const walletPromises = walletIds.map(id => walletAPI.get(id).catch(() => null))
    const results = await Promise.all(walletPromises)
    wallets.value = results.filter(r => r?.data).map(r => r.data)
  } catch (error) {
    console.error('Failed to load wallets:', error)
  } finally {
    isLoading.value = false
  }
}

const createWallet = async () => {
  isCreating.value = true
  try {
    const { data } = await walletAPI.create({
      currency: newWallet.value.currency,
      wallet_type: 'business',
      name: newWallet.value.name || `Portefeuille ${newWallet.value.currency}`
    })
    
    // Add wallet ID to enterprise
    const updatedWalletIds = [...(props.enterprise.wallet_ids || []), data.id]
    await enterpriseAPI.update(props.enterprise.id, {
      ...props.enterprise,
      wallet_ids: updatedWalletIds
    })
    
    emit('update', { ...props.enterprise, wallet_ids: updatedWalletIds })
    showCreateWallet.value = false
    newWallet.value = { currency: 'XOF', name: '' }
    await loadWallets()
  } catch (error) {
    console.error('Failed to create wallet:', error)
    alert('Erreur lors de la création du portefeuille')
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

const viewTransactions = (wallet) => {
  // Navigate to transactions tab with wallet filter
  navigateTo(`/wallets/${wallet.id}`)
}

onMounted(loadWallets)
</script>
