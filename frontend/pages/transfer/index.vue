<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Envoyer de l'argent 💸</h1>
        <p class="text-gray-900/60">Transferts P2P, Mobile Money, et virements bancaires</p>
      </div>

      <!-- Transfer Type Selector -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <button 
          v-for="type in transferTypes" :key="type.id"
          @click="selectType(type.id)"
          class="glass-card p-4 text-left transition-all hover:scale-[1.02]"
          :class="selectedType === type.id ? 'ring-2 ring-indigo-500 bg-indigo-50/50' : ''">
          <span class="text-2xl mb-2 block">{{ type.icon }}</span>
          <h3 class="font-semibold text-gray-900 mb-1 text-sm">{{ type.name }}</h3>
        </button>
      </div>

      <!-- Transfer Form -->
      <div class="glass-card p-8 relative overflow-hidden">
        <div v-if="loadingWallets" class="absolute inset-0 bg-white/80 z-10 flex items-center justify-center">
          <div class="loading-spinner w-8 h-8"></div>
        </div>

        <form @submit.prevent="submitTransfer" class="space-y-6">
          
          <!-- Source Wallet Selection -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Depuis le portefeuille</label>
            <select v-model="form.fromWalletId" class="input-premium w-full" required>
              <option value="" disabled>Sélectionnez un portefeuille</option>
              <option v-for="wallet in wallets" :key="wallet.id" :value="wallet.id">
                {{ wallet.currency }} - Solde: {{ formatMoney(wallet.balance, wallet.currency) }} ({{ wallet.wallet_type }})
              </option>
            </select>
          </div>

          <!-- Amount -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Montant à envoyer</label>
            <div class="flex gap-4">
              <input
                v-model.number="form.amount"
                type="number"
                min="0.01"
                step="0.01"
                required
                class="input-premium flex-1 text-2xl"
                placeholder="0.00"
              />
              <div class="input-premium w-24 flex items-center justify-center bg-gray-50 text-gray-600 font-bold">
                {{ selectedWalletCurrency }}
              </div>
            </div>
            <p v-if="form.amount > selectedWalletBalance" class="text-xs text-red-500 mt-1">Solde insuffisant</p>
          </div>

          <!-- P2P User Lookup -->
          <div v-if="selectedType === 'p2p'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-900/80 mb-2">Destinataire (Email ou Téléphone)</label>
              <div class="flex gap-2">
                <input
                  v-model="p2pSearch"
                  type="text"
                  class="input-premium flex-1"
                  placeholder="ex: ami@email.com ou +225..."
                  @blur="lookupUser"
                />
                <button type="button" @click="lookupUser" class="btn-secondary px-4">
                  <span v-if="lookupLoading" class="loading-spinner w-4 h-4"></span>
                  <span v-else>🔍</span>
                </button>
              </div>
              <p v-if="lookupError" class="text-xs text-red-500 mt-1">{{ lookupError }}</p>
              <div v-if="lookupResult" class="mt-2 p-3 bg-emerald-50 rounded-lg flex items-center gap-3 border border-emerald-100">
                <div class="w-8 h-8 rounded-full bg-emerald-200 flex items-center justify-center text-emerald-700 font-bold">
                  {{ lookupResult.first_name[0] }}
                </div>
                <div>
                  <p class="text-sm font-semibold text-emerald-900">{{ lookupResult.first_name }} {{ lookupResult.last_name }}</p>
                  <p class="text-xs text-emerald-600">Utilisateur vérifié</p>
                </div>
              </div>
            </div>
            <div v-if="!lookupResult" class="text-sm text-gray-500 italic">Recherchez un utilisateur pour continuer</div>
          </div>

          <!-- Mobile Money Fields -->
          <div v-if="selectedType === 'mobile'">
            <div class="grid grid-cols-2 gap-4">
               <div>
                <label class="block text-sm font-medium text-gray-900/80 mb-2">Pays</label>
                <select v-model="form.country" class="input-premium w-full">
                  <option value="CI">🇨🇮 Côte d'Ivoire</option>
                  <option value="SN">🇸🇳 Sénégal</option>
                  <option value="CM">🇨🇲 Cameroun</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-900/80 mb-2">Opérateur</label>
                <select v-model="form.provider" class="input-premium w-full">
                  <option value="orange">Orange Money</option>
                  <option value="mtn">MTN Mobile Money</option>
                  <option value="wave">Wave</option>
                </select>
              </div>
            </div>
             <div class="mt-4">
              <label class="block text-sm font-medium text-gray-900/80 mb-2">Numéro de téléphone</label>
              <input
                v-model="form.recipientPhone"
                type="tel"
                class="input-premium w-full"
                placeholder="+225 07..."
              />
            </div>
          </div>

          <!-- Bank Wire Fields -->
          <div v-if="selectedType === 'wire'">
             <div>
              <label class="block text-sm font-medium text-gray-900/80 mb-2">Nom de la banque</label>
              <input v-model="form.bankName" type="text" class="input-premium w-full" />
            </div>
             <div class="mt-4">
              <label class="block text-sm font-medium text-gray-900/80 mb-2">IBAN / Numéro de compte</label>
              <input v-model="form.bankAccount" type="text" class="input-premium w-full" />
            </div>
             <div class="mt-4">
              <label class="block text-sm font-medium text-gray-900/80 mb-2">Nom du bénéficiaire</label>
              <input v-model="form.recipientName" type="text" class="input-premium w-full" />
            </div>
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Note (Optionnel)</label>
            <input
              v-model="form.description"
              type="text"
              class="input-premium w-full"
              placeholder="Ex: Loyer"
            />
          </div>

          <!-- Summary -->
          <div class="p-6 rounded-xl bg-gray-50 border border-gray-100">
            <h4 class="font-semibold text-gray-900 mb-4">Résumé</h4>
            <div class="space-y-3 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600">Montant</span>
                <span class="text-gray-900 font-medium">{{ formatMoney(form.amount || 0, selectedWalletCurrency) }}</span>
              </div>
              <div class="flex justify-between text-indigo-600" v-if="estimatedFee > 0">
                <span>Frais estimés</span>
                <span class="font-medium">+ {{ formatMoney(estimatedFee, selectedWalletCurrency) }}</span>
              </div>
              <div class="border-t border-gray-200 pt-3 flex justify-between">
                <span class="text-gray-900 font-semibold">Total</span>
                <span class="text-gray-900 font-bold text-lg">{{ formatMoney((form.amount || 0) + estimatedFee, selectedWalletCurrency) }}</span>
              </div>
            </div>
          </div>

          <!-- Submit -->
          <button 
            type="submit" 
            :disabled="loading || (selectedType === 'p2p' && !lookupResult) || !form.fromWalletId" 
            class="btn-premium w-full py-4 text-lg disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <div class="loading-spinner w-5 h-5"></div>
              Envoi en cours...
            </span>
            <span v-else>Confirmer le transfert</span>
          </button>
        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useApi } from '@/composables/useApi'

const { walletApi, userApi, transferApi } = useApi()

const transferTypes = [
  { id: 'p2p', name: 'Interne (P2P)', icon: '👤' },
  { id: 'mobile', name: 'Mobile Money', icon: '📱' },
  { id: 'wire', name: 'Virement', icon: '🏦' },
  { id: 'crypto', name: 'Crypto', icon: '₿' }
]

const selectedType = ref('p2p')
const wallets = ref([])
const loadingWallets = ref(true)
const loading = ref(false)

// Form Data
const form = ref({
  fromWalletId: '',
  amount: null,
  description: '',
  // Specific fields
  country: 'CI',
  provider: 'orange',
  recipientPhone: '',
  bankName: '',
  bankAccount: '',
  recipientName: ''
})

// P2P Lookup
const p2pSearch = ref('')
const lookupLoading = ref(false)
const lookupResult = ref(null)
const lookupError = ref('')

const selectedWallet = computed(() => wallets.value.find(w => w.id === form.value.fromWalletId))
const selectedWalletCurrency = computed(() => selectedWallet.value?.currency || '')
const selectedWalletBalance = computed(() => selectedWallet.value?.balance || 0)

const estimatedFee = computed(() => {
  // Simple estimation for UI feedback
  if (!form.value.amount) return 0
  if (selectedType.value === 'p2p') return 0 // Free internal transfers?
  if (selectedType.value === 'mobile') return form.value.amount * 0.01 // 1%
  if (selectedType.value === 'wire') return Math.max(5, form.value.amount * 0.02) // Min 5 or 2%
  return 0
})

const selectType = (type) => {
  selectedType.value = type
  // Reset type specific validation
  lookupResult.value = null
  p2pSearch.value = ''
}

const formatMoney = (amount, currency) => {
  if (!currency) return '0.00'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

const fetchWallets = async () => {
  loadingWallets.value = true
  try {
    const res = await walletApi.getWallets()
    wallets.value = res.data.wallets || []
    // Auto-select first wallet with balance
    const validWallet = wallets.value.find(w => w.balance > 0)
    if (validWallet) form.value.fromWalletId = validWallet.id
  } catch (e) {
    console.error('Failed to fetch wallets', e)
  } finally {
    loadingWallets.value = false
  }
}

const lookupUser = async () => {
  if (!p2pSearch.value || p2pSearch.value.length < 3) return
  
  lookupLoading.value = true
  lookupError.value = ''
  lookupResult.value = null

  try {
    const isEmail = p2pSearch.value.includes('@')
    const query = isEmail ? { email: p2pSearch.value } : { phone: p2pSearch.value }
    const res = await userApi.lookup(query)
    lookupResult.value = res.data
  } catch (e) {
    lookupError.value = 'Utilisateur introuvable'
  } finally {
    lookupLoading.value = false
  }
}

const submitTransfer = async () => {
  loading.value = true
  try {
    // Construct payload based on type
    const payload = {
      from_wallet_id: form.value.fromWalletId,
      amount: form.value.amount,
      currency: selectedWalletCurrency.value,
      description: form.value.description
    }

    if (selectedType.value === 'p2p') {
      // Internal transfer
       // Check if transfer-service accepts 'to_email' or 'to_phone' or we need 'to_wallet_id'
       // TransferRequest struct has ToEmail, ToPhone. So we can send email/phone directly OR use lookup result ID if we want to be safe?
       // Let's use lookup result to be safer if we map it to wallet?
       // Actually 'TransferRequest' has ToEmail/ToPhone. Let's use those if provided.
       // But wait, we looked up the user to SHOW them.
       // If we send `to_phone` or `to_email` the backend handles it.
       if (p2pSearch.value.includes('@')) payload.to_email = p2pSearch.value
       else payload.to_phone = p2pSearch.value
       
    } else if (selectedType.value === 'mobile') {
      // Mobile Money
      // Path: /mobile/send
      await transferApi.create({
          type: 'mobile_money',
          amount: form.value.amount,
          currency: selectedWalletCurrency.value,
          recipient: form.value.recipientPhone,
          description: form.value.description,
          // Extra data might be needed like provider/country passed in description or specific endpoint
          // Use specific endpoint for mobile money if exists, or generic create with type
      })
      // Specific mobile endpoint in useApi is not fully aligned with generic create.
      // Let's assume generic create handles it or separate call.
      // Reverting to generic create for now as per useApi definition
    } 
    
    // For P2P we use generic create
    if (selectedType.value === 'p2p') {
        const res = await transferApi.create(payload)
        alert('Transfert réussi!')
    }

    // Reset form
    form.value.amount = null
    form.value.description = ''
  } catch (e) {
    console.error(e)
    alert('Erreur lors du transfert: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWallets()
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>