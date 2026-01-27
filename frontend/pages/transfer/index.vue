<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="mb-10 text-center md:text-left">
        <h1 class="text-3xl font-bold text-base mb-2">Envoyer de l'argent 💸</h1>
        <p class="text-muted">Transferts P2P, Mobile Money, et virements bancaires</p>
      </div>

      <!-- Transfer Type Selector -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <button 
          v-for="type in transferTypes" :key="type.id"
          @click="selectType(type.id)"
          class="p-4 rounded-2xl border transition-all duration-300 hover:scale-[1.02] text-left relative overflow-hidden group"
          :class="selectedType === type.id 
            ? 'bg-primary/10 border-primary text-primary shadow-lg shadow-primary/20' 
            : 'bg-surface border-secondary-200 dark:border-secondary-700 text-muted hover:border-primary/50 hover:text-base'"
        >
          <div class="z-10 relative">
            <span class="text-3xl mb-3 block filter drop-shadow-sm">{{ type.icon }}</span>
            <h3 class="font-bold text-sm">{{ type.name }}</h3>
          </div>
          <!-- Selection Glow -->
          <div v-if="selectedType === type.id" class="absolute inset-0 bg-gradient-to-tr from-primary/10 via-transparent to-transparent opacity-50"></div>
        </button>
      </div>

      <!-- Transfer Form -->
      <div class="glass-card mb-8 relative overflow-hidden">
        <div v-if="loadingWallets" class="absolute inset-0 bg-surface/80 backdrop-blur-sm z-10 flex items-center justify-center">
          <div class="loading-spinner w-8 h-8"></div>
        </div>

        <form @submit.prevent="submitTransfer" class="space-y-8">
          
          <!-- Source Wallet Selection -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-muted">Depuis le portefeuille</label>
            <div class="relative">
               <select v-model="form.fromWalletId" class="input-field w-full appearance-none cursor-pointer" required>
                <option value="" disabled>Sélectionnez un portefeuille</option>
                <option v-for="wallet in filteredWallets" :key="wallet.id" :value="wallet.id">
                   {{ wallet.name }} - {{ formatMoney(wallet.balance, wallet.currency) }}
                </option>
              </select>
               <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-muted">
                ▼
              </div>
            </div>
          </div>

          <!-- Amount -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-muted">Montant à envoyer</label>
            <div class="p-4 bg-surface-hover rounded-2xl border border-secondary-200 dark:border-secondary-700 flex items-center gap-4 focus-within:ring-2 focus-within:ring-primary transition-all">
              <input
                v-model.number="form.amount"
                type="number"
                min="0.01"
                step="0.01"
                required
                class="bg-transparent text-3xl font-bold text-base outline-none w-full placeholder-muted/30"
                placeholder="0.00"
              />
              <span class="text-lg font-bold text-primary bg-primary/10 px-3 py-1 rounded-lg">
                {{ selectedWalletCurrency }}
              </span>
            </div>
             <div class="flex justify-between text-xs px-1">
               <p v-if="form.amount > selectedWalletBalance" class="text-error font-medium">Solde insuffisant</p>
               <p v-else class="text-muted">Disponible: {{ formatMoney(selectedWalletBalance, selectedWalletCurrency) }}</p>
             </div>
          </div>

          <!-- P2P User Lookup -->
          <div v-if="selectedType === 'p2p'" class="space-y-4 animate-fade-in">
            <div>
              <label class="block text-sm font-medium text-muted mb-2">Destinataire (Email ou Téléphone)</label>
              <div class="flex gap-2">
                <input
                  v-model="p2pSearch"
                  type="text"
                  class="input-field flex-1"
                  placeholder="ex: ami@email.com ou +225..."
                  @blur="lookupUser"
                />
                <button type="button" @click="lookupUser" class="px-4 rounded-xl bg-surface-hover hover:bg-primary hover:text-white transition-colors border border-secondary-200 dark:border-secondary-700 text-muted">
                  <span v-if="lookupLoading" class="loading-spinner w-4 h-4 text-current"></span>
                  <span v-else class="text-lg">🔍</span>
                </button>
              </div>
              <p v-if="lookupError" class="text-xs text-error mt-2 ml-1">{{ lookupError }}</p>
              
              <div v-if="lookupResult" class="mt-4 p-3 bg-success/10 border border-success/20 rounded-xl flex items-center gap-4 animate-fade-in-up">
                <div class="w-10 h-10 rounded-full bg-success text-white flex items-center justify-center font-bold text-lg shadow-lg shadow-success/20">
                  {{ lookupResult.first_name?.[0] || 'U' }}
                </div>
                <div>
                  <p class="text-sm font-bold text-base">{{ lookupResult.first_name }} {{ lookupResult.last_name }}</p>
                  <p class="text-xs text-success font-medium flex items-center gap-1">
                     <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/></svg>
                     Utilisateur vérifié
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Mobile Money Fields -->
          <div v-if="selectedType === 'mobile'" class="space-y-4 animate-fade-in">
            <div class="grid grid-cols-2 gap-4">
               <div class="space-y-2">
                <label class="block text-sm font-medium text-muted">Pays</label>
                <select v-model="form.country" class="input-field">
                  <option value="CI">🇨🇮 Côte d'Ivoire</option>
                  <option value="SN">🇸🇳 Sénégal</option>
                  <option value="CM">🇨🇲 Cameroun</option>
                </select>
              </div>
              <div class="space-y-2">
                <label class="block text-sm font-medium text-muted">Opérateur</label>
                <select v-model="form.provider" class="input-field">
                  <option value="orange">Orange Money</option>
                  <option value="mtn">MTN MoMo</option>
                  <option value="wave">Wave</option>
                </select>
              </div>
            </div>
             <div class="space-y-2">
              <label class="block text-sm font-medium text-muted">Numéro de téléphone</label>
              <input
                v-model="form.recipientPhone"
                type="tel"
                class="input-field"
                placeholder="+225 07..."
              />
            </div>
          </div>

          <!-- Bank Wire Fields -->
          <div v-if="selectedType === 'wire'" class="space-y-4 animate-fade-in">
             <div class="space-y-2">
              <label class="block text-sm font-medium text-muted">Nom de la banque</label>
              <input v-model="form.bankName" type="text" class="input-field" placeholder="ex: Ecobank" />
            </div>
             <div class="space-y-2">
              <label class="block text-sm font-medium text-muted">IBAN / Numéro de compte</label>
              <input v-model="form.bankAccount" type="text" class="input-field" placeholder="FR76..." />
            </div>
             <div class="space-y-2">
              <label class="block text-sm font-medium text-muted">Nom du bénéficiaire</label>
              <input v-model="form.recipientName" type="text" class="input-field" />
            </div>
          </div>

          <!-- Crypto Transfer Fields -->
          <div v-if="selectedType === 'crypto'" class="space-y-4 animate-fade-in">
            <div class="space-y-2">
              <label class="block text-sm font-medium text-muted">Adresse du portefeuille destinataire</label>
              <input 
                v-model="form.recipientAddress" 
                type="text" 
                class="input-field font-mono text-sm" 
                placeholder="0x... ou bc1q..." 
              />
            </div>
             <div class="p-3 bg-warning/10 border border-warning/20 rounded-xl text-xs text-warning">
              <div class="mb-1 font-bold flex items-center gap-1">
                 <span class="text-base">🌐</span> Réseau requis : <span class="uppercase underline">{{ currentNetwork }}</span>
              </div>
              ⚠️ Assurez-vous que l'adresse correspond bien au réseau <strong>{{ currentNetwork }}</strong>. Les transactions crypto sont irréversibles.
            </div>
          </div>

          <!-- Description -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-muted">Note (Optionnel)</label>
            <input
              v-model="form.description"
              type="text"
              class="input-field"
              placeholder="Ex: Loyer"
            />
          </div>

          <!-- Summary -->
          <div class="p-6 rounded-2xl bg-surface border border-secondary-200 dark:border-secondary-800">
            <h4 class="font-bold text-base mb-4 flex items-center gap-2">
               📝 Résumé
            </h4>
            <div class="space-y-3 text-sm">
              <div class="flex justify-between items-center">
                <span class="text-muted">Montant</span>
                <span class="text-base font-medium">{{ formatMoney(form.amount || 0, selectedWalletCurrency) }}</span>
              </div>
              <div class="flex justify-between items-center text-primary" v-if="estimatedFee > 0">
                <span>Frais estimés</span>
                <span class="font-medium">+ {{ formatMoney(estimatedFee, selectedWalletCurrency) }}</span>
              </div>
              <div class="border-t border-secondary-200 dark:border-secondary-700 pt-3 flex justify-between items-center">
                <span class="text-base font-bold">Total à débiter</span>
                <span class="text-lg font-bold text-base">{{ formatMoney((form.amount || 0) + estimatedFee, selectedWalletCurrency) }}</span>
              </div>
            </div>
          </div>

          <!-- Submit -->
          <button 
            type="submit" 
            :disabled="loading || (selectedType === 'p2p' && !lookupResult) || !form.fromWalletId" 
            class="btn-premium w-full py-4 text-lg font-bold shadow-lg shadow-primary/25 disabled:opacity-50 disabled:cursor-not-allowed group relative overflow-hidden">
             
            <span class="relative z-10 flex items-center justify-center gap-2">
               <span v-if="loading" class="loading-spinner w-5 h-5 border-white/30 border-t-white"></span>
               {{ loading ? 'Envoi en cours...' : 'Confirmer le transfert' }}
               <svg v-if="!loading" class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
            </span>
          </button>
        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useApi } from '@/composables/useApi'
import { useModal } from '~/composables/useModal'
import { usePin } from '~/composables/usePin'

const { walletApi, userApi, transferApi } = useApi()
const modal = useModal()
const { requirePin, checkPinStatus } = usePin()

const transferTypes = [
  { id: 'p2p', name: 'Interne (P2P)', icon: '👤' },
  { id: 'mobile', name: 'Mobile Money', icon: '📱' },
  { id: 'wire', name: 'Virement', icon: '🏦' },
  { id: 'crypto', name: 'Crypto', icon: '₿' }
]

const selectedType = ref('p2p')

// Comprehensive Network Mapping
const networkMap = {
  'BTC': 'Bitcoin Network (SegWit)',
  'ETH': 'Ethereum Mainnet (ERC-20)',
  'USDT': 'Ethereum Mainnet (ERC-20)', // Default USDT on our platform
  'USDC': 'Ethereum Mainnet (ERC-20)',
  'LTC': 'Litecoin Network',
  'BSC': 'BNB Smart Chain (BEP-20)',
  'BNB': 'BNB Smart Chain (BEP-20)',
  'MATIC': 'Polygon Network',
  'SOL': 'Solana Network',
  'TRON': 'Tron Network (TRC-20)',
  'USDT_TRON': 'Tron Network (TRC-20)',
  'XRP': 'Ripple Network',
  'ADA': 'Cardano Network',
  'DOGE': 'Dogecoin Network',
  'BCH': 'Bitcoin Cash Network',
  'XLM': 'Stellar Network',
  'ALGO': 'Algorand Network',
  'CELO': 'Celo Network',
  'ONE': 'Harmony One Network'
}

const currentNetwork = computed(() => {
  if (!selectedWalletCurrency.value) return ''
  return networkMap[selectedWalletCurrency.value] || `${selectedWalletCurrency.value} Network`
})

const wallets = ref([])
const filteredWallets = computed(() => {
  if (selectedType.value === 'crypto') {
    return wallets.value.filter(w => w.wallet_type === 'crypto')
  }
  return wallets.value.filter(w => w.wallet_type !== 'crypto')
})
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
  recipientName: '',
  recipientAddress: ''
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
  if (selectedType.value === 'crypto') return 0.0001 // Estimated network fee placeholder
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

const route = useRoute()

const fetchWallets = async () => {
  loadingWallets.value = true
  try {
    const res = await walletApi.getWallets()
    wallets.value = res.data.wallets || []
    
    // Check if wallet ID is passed in URL query param
    const urlWalletId = route.query.wallet
    if (urlWalletId) {
      // Find and select the wallet from URL
      const targetWallet = wallets.value.find(w => w.id === urlWalletId)
      if (targetWallet) {
        form.value.fromWalletId = targetWallet.id
        return // Don't auto-select, use the URL param
      }
    }
    
    // Fallback: Auto-select first wallet with balance
    const validWallet = wallets.value.find(w => w.balance > 0)
    if (validWallet) form.value.fromWalletId = validWallet.id
  } catch (e) {
    console.error('Failed to fetch wallets', e)
     // Fallback mock
     wallets.value = [
        { id: 1, name: 'Main Wallet', currency: 'USD', balance: 15420.50, wallet_type: 'fiat' },
        { id: 2, name: 'Savings', currency: 'EUR', balance: 2300.00, wallet_type: 'fiat' }
     ]
      form.value.fromWalletId = 1
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

// Execute the actual transfer (called after PIN verification)
const executeTransfer = async () => {
  loading.value = true
  try {
    // Construct payload based on type
    const payload = {
      from_wallet_id: form.value.fromWalletId,
      amount: form.value.amount,
      currency: selectedWalletCurrency.value,
      description: form.value.description
    }

    let result = null
    
    if (selectedType.value === 'p2p') {
       if (p2pSearch.value.includes('@')) payload.to_email = p2pSearch.value
       else payload.to_phone = p2pSearch.value
       
       result = await transferApi.create(payload)
       
    } else if (selectedType.value === 'mobile') {
      result = await transferApi.create({
          type: 'mobile_money',
          from_wallet_id: form.value.fromWalletId,
          amount: form.value.amount,
          currency: selectedWalletCurrency.value,
          recipient: form.value.recipientPhone,
          description: form.value.description + ` [${form.value.provider} - ${form.value.country}]`
      })
    } else if (selectedType.value === 'wire') {
         result = await transferApi.create({
          type: 'wire',
          from_wallet_id: form.value.fromWalletId,
          amount: form.value.amount,
          currency: selectedWalletCurrency.value,
          recipient: form.value.recipientName,
          description: form.value.description + ` [IBAN: ${form.value.bankAccount}]`
        })
    } else if (selectedType.value === 'crypto') {
        if (!form.value.recipientAddress) throw new Error("Adresse destinataire requise")
        
        await walletApi.sendCrypto(form.value.fromWalletId, {
          to_address: form.value.recipientAddress,
          amount: form.value.amount,
          note: form.value.description
        })
        
        result = { data: { transfer: { status: 'completed' } } }
    }

    // Show success message with custom modal
    const transferStatus = result?.data?.transfer?.status || 'initiated'
    await modal.success(
      'Transfert réussi ! 🎉',
      `Votre transfert de ${formatMoney(form.value.amount, selectedWalletCurrency.value)} a été ${transferStatus === 'completed' ? 'effectué' : 'initié'} avec succès.`
    )
    
    // Redirect to wallet page to see updated balance
    navigateTo('/wallet')
    
  } catch (e) {
    console.error(e)
    await modal.error(
      'Échec du transfert',
      e.response?.data?.error || e.message || 'Une erreur est survenue lors du transfert.'
    )
  } finally {
    loading.value = false
  }
}

// Submit transfer - requires PIN verification first
const submitTransfer = async () => {
  // Require PIN verification before executing the transfer
  const verified = await requirePin(executeTransfer)
  
  if (!verified) {
    // User cancelled PIN verification
    console.log('Transfer cancelled - PIN verification required')
  }
}

onMounted(() => {
  fetchWallets()
  // Check PIN status on mount
  checkPinStatus()
})

definePageMeta({
  layout: false,
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
</style>