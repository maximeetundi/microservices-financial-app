<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white mb-2">My Wallets 👛</h1>
          <p class="text-gray-500 dark:text-gray-400">Manage your fiat and crypto currencies</p>
        </div>
        <div class="flex gap-3">
           <NuxtLink to="/cards" class="btn-secondary flex items-center gap-2">
            <span class="text-xl">💳</span>
            <span>Commander une carte</span>
          </NuxtLink>
          <button @click="showCreateWallet = true" class="btn-primary flex items-center gap-2 shadow-lg shadow-indigo-500/30">
            <span class="text-xl">+</span>
            <span>Nouveau Portefeuille</span>
          </button>
        </div>
      </div>

      <!-- Total Balance Card -->
      <div class="glass-card mb-8 p-8 bg-white dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 relative overflow-hidden group">
        <!-- Background Effects -->
        <div class="absolute top-0 right-0 w-64 h-64 bg-indigo-500/10 rounded-full blur-3xl group-hover:bg-indigo-500/20 transition-all duration-500"></div>
        
        <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 relative z-10">
          <div>
            <p class="text-gray-500 dark:text-gray-400 font-medium mb-1 uppercase tracking-wider text-sm">Valeur Totale</p>
            <div class="flex items-baseline gap-3">
              <h2 class="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300">
                {{ formatMoney(totalBalance) }}
              </h2>
            </div>
            <p class="text-sm mt-3 flex items-center gap-2 text-emerald-600 dark:text-emerald-400 font-medium bg-emerald-50 dark:bg-emerald-500/10 px-3 py-1 rounded-full w-fit">
              <span class="relative flex h-2 w-2">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
              </span>
              Calculé en temps réel
            </p>
          </div>
          <div class="flex flex-wrap gap-3">
             <!-- Recharger / Deposit -->
            <button @click="openTopUpModal" class="flex items-center gap-2 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-700 shadow-lg shadow-indigo-500/30 transition-all transform hover:-translate-y-0.5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              Recharger
            </button>
            <!-- Envoyer / Send -->
            <NuxtLink to="/transfer" class="flex items-center gap-2 px-6 py-3 rounded-xl bg-white dark:bg-slate-800 text-gray-700 dark:text-white border border-gray-200 dark:border-gray-700 font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-all transform hover:-translate-y-0.5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
              Envoyer
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Wallets Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="wallet in wallets" :key="wallet.id" 
            class="glass-card p-6 cursor-pointer hover:border-indigo-500 dark:hover:border-indigo-500 transition-all duration-300 relative overflow-hidden group bg-white dark:bg-slate-900 border border-gray-200 dark:border-white/10"
            :class="{'ring-2 ring-indigo-500 dark:ring-indigo-400': selectedWallet?.id === wallet.id}"
            @click="selectWallet(wallet)">
            
           <!-- Background Glow -->
           <div class="absolute -right-10 -top-10 w-32 h-32 rounded-full blur-3xl opacity-0 group-hover:opacity-10 transition-opacity duration-500"
                :class="getCurrencyBg(wallet.currency)"></div>

          <div class="flex items-center justify-between mb-6 relative z-10">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shadow-md text-white font-bold" 
                   :class="getCurrencyBg(wallet.currency)">
                {{ getCurrencyIcon(wallet.currency) }}
              </div>
              <div>
                <p class="font-bold text-gray-900 dark:text-white text-lg">{{ wallet.currency }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 font-medium uppercase tracking-wider">{{ wallet.name }}</p>
              </div>
            </div>
            <div class="flex flex-col items-end gap-1">
               <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide border" 
                  :class="wallet.status === 'active' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-500/20' : 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-200 dark:border-amber-500/20'">
              {{ wallet.status }}
            </span>
             <span class="text-[10px] font-mono text-gray-400">{{ wallet.type }}</span>
            </div>
          </div>

          <div class="mb-4 relative z-10">
            <p class="text-2xl font-extrabold text-gray-900 dark:text-white tracking-tight">
              {{ formatCrypto(wallet.balance, wallet.currency) }}
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">
              ≈ {{ formatMoney(wallet.balanceUSD) }}
            </p>
          </div>

          <div class="flex justify-between items-center text-sm relative z-10 pt-4 border-t border-gray-100 dark:border-gray-800">
            <span class="text-gray-400 font-medium">Adresse</span>
            <div class="flex items-center gap-2">
                 <span class="text-gray-600 dark:text-gray-300 font-mono text-xs bg-gray-100 dark:bg-slate-800 px-2 py-1 rounded">{{ truncateAddress(wallet.address) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Selected Wallet Actions (If any layout specific needed, handled by modal) -->

      <!-- Create Wallet Modal -->
      <div v-if="showCreateWallet" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-fade-in-up">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-8 max-w-md w-full shadow-2xl border border-gray-100 dark:border-gray-800">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">Nouveau Portefeuille</h3>
            <button @click="showCreateWallet = false" class="text-gray-400 hover:text-gray-600 dark:hover:text-white transition-colors">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <form @submit.prevent="createWallet" class="space-y-5">
            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Type de devise</label>
              <select v-model="newWallet.type" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
                <option value="fiat">Devise Fiat (USD, EUR...)</option>
                <option value="crypto">Crypto-monnaie (BTC, ETH...)</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Devise</label>
              <select v-model="newWallet.currency" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
                <template v-if="newWallet.type === 'fiat'">
                  <option value="USD">🇺🇸 Dollar US (USD)</option>
                  <option value="EUR">🇪🇺 Euro (EUR)</option>
                  <option value="GBP">🇬🇧 Livre Sterling (GBP)</option>
                  <option value="XOF">🇨🇮 Franc CFA (XOF)</option>
                  <option value="XAF">🇨🇲 Franc CFA (XAF)</option>
                </template>
                <template v-else>
                  <option value="BTC">₿ Bitcoin (BTC)</option>
                  <option value="ETH">Ξ Ethereum (ETH)</option>
                  <option value="USDT">💵 Tether (USDT)</option>
                  <option value="SOL">◎ Solana (SOL)</option>
                </template>
              </select>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Nom du portefeuille</label>
              <input 
                v-model="newWallet.name" 
                type="text" 
                placeholder="ex: Mon portefeuille principal"
                class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
              />
            </div>

            <div class="flex gap-3 mt-8">
              <button 
                type="button" 
                @click="showCreateWallet = false"
                class="flex-1 py-3 px-4 rounded-xl font-bold text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors"
              >
                Annuler
              </button>
              <button 
                type="submit" 
                :disabled="creatingWallet"
                class="flex-1 py-3 px-4 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/30"
              >
                {{ creatingWallet ? 'Création...' : 'Créer le portefeuille' }}
              </button>
            </div>
          </form>
        </div>
      </div>

       <!-- Recharge / Top Up Modal -->
      <div v-if="showTopUpModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-fade-in-up">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-0 max-w-lg w-full shadow-2xl border border-gray-100 dark:border-gray-800 overflow-hidden flex flex-col max-h-[90vh]">
           <div class="p-6 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center">
             <h3 class="text-xl font-bold text-gray-900 dark:text-white">Recharger</h3>
             <button @click="showTopUpModal = false" class="text-gray-400 hover:text-gray-900 dark:hover:text-white">✕</button>
           </div>
           
           <div class="p-6 overflow-y-auto custom-scrollbar">
             <!-- Wallet Selection if likely needed, or display current -->
             <div class="mb-6 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl flex items-center gap-4">
               <div class="w-10 h-10 rounded-full flex items-center justify-center text-xl font-bold text-white shadow-sm" :class="getCurrencyBg(selectedWallet?.currency)">
                 {{ getCurrencyIcon(selectedWallet?.currency) }}
               </div>
               <div>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Portefeuille cible</p>
                  <p class="font-bold text-gray-900 dark:text-white">{{ selectedWallet?.name }} ({{ selectedWallet?.currency }})</p>
               </div>
             </div>

             <!-- Conditional Content based on Wallet Type -->
             <div v-if="selectedWallet?.type === 'fiat'">
                <h4 class="font-bold text-gray-900 dark:text-white mb-4">Moyens de paiement</h4>
                <div class="grid grid-cols-1 gap-3">
                  <button class="flex items-center gap-4 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-indigo-500 dark:hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-500/10 transition-all group">
                    <span class="text-2xl">📱</span>
                    <div class="text-left">
                      <p class="font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400">Mobile Money</p>
                      <p class="text-xs text-gray-500 dark:text-gray-400">Orange, MTN, Wave</p>
                    </div>
                  </button>
                   <button class="flex items-center gap-4 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-indigo-500 dark:hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-500/10 transition-all group">
                    <span class="text-2xl">🏦</span>
                    <div class="text-left">
                      <p class="font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400">Virement Bancaire</p>
                      <p class="text-xs text-gray-500 dark:text-gray-400">IBAN / RIB</p>
                    </div>
                  </button>
                    <button class="flex items-center gap-4 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-indigo-500 dark:hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-500/10 transition-all group">
                    <span class="text-2xl">💳</span>
                    <div class="text-left">
                      <p class="font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400">Carte Bancaire</p>
                      <p class="text-xs text-gray-500 dark:text-gray-400">Visa, Mastercard</p>
                    </div>
                  </button>
                </div>
             </div>

             <div v-else>
               <!-- Crypto Deposit -->
                <div class="text-center">
                  <div class="bg-white p-4 rounded-xl inline-block shadow-lg border border-gray-100 mb-6">
                    <!-- Placeholder QR -->
                    <div class="w-48 h-48 bg-gray-100 flex items-center justify-center text-gray-400 text-xs">QR Code {{ selectedWallet?.currency }}</div>
                  </div>
                  <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">Adresse de dépôt {{ selectedWallet?.currency }}</p>
                  <div class="flex items-center gap-2 bg-gray-100 dark:bg-slate-800 p-3 rounded-xl border border-gray-200 dark:border-gray-700">
                    <code class="text-sm font-mono flex-1 break-all text-gray-800 dark:text-gray-200">{{ selectedWallet?.address }}</code>
                    <button @click="copyAddress" class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">📋</button>
                  </div>
                  <p class="text-xs text-amber-600 mt-4 bg-amber-50 dark:bg-amber-900/20 p-3 rounded-lg">
                    ⚠️ Envoyez uniquement du <strong>{{ selectedWallet?.currency }}</strong> sur cette adresse. Tout autre jeton sera perdu.
                  </p>
                </div>
             </div>
           </div>
        </div>
      </div>

    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { walletAPI } from '~/composables/useApi'
import { useRouter } from 'vue-router'

const router = useRouter()

// Wallets - will be loaded from API
const wallets = ref([])

const selectedWallet = ref(null)
const showCreateWallet = ref(false)
const showTopUpModal = ref(false)
const creatingWallet = ref(false)

// New wallet form
const newWallet = ref({
  type: 'fiat',
  currency: 'USD',
  name: ''
})

const totalBalance = computed(() => {
  return wallets.value.reduce((sum, w) => {
    const val = Number(w.balanceUSD)
    return sum + (isNaN(val) ? 0 : val)
  }, 0)
})

const formatMoney = (amount) => {
  const val = Number(amount)
  if (amount === undefined || amount === null || isNaN(val)) return '$0.00'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val)
}

const formatCrypto = (amount, currency) => {
  const val = Number(amount)
  if (amount === undefined || amount === null || isNaN(val)) return `0.00 ${currency}`
  if (['BTC', 'ETH', 'SOL'].includes(currency)) {
    return `${val.toFixed(6)} ${currency}`
  }
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(val)
}

const getCurrencyIcon = (currency) => {
  const icons = { BTC: '₿', ETH: 'Ξ', USD: '$', EUR: '€', GBP: '£', SOL: '◎' }
  return icons[currency] || '💰'
}

const getCurrencyBg = (currency) => {
  const bgs = { 
    BTC: 'bg-amber-500', 
    ETH: 'bg-indigo-500', 
    USD: 'bg-emerald-500', 
    EUR: 'bg-blue-600', 
    SOL: 'bg-purple-500',
    GBP: 'bg-rose-500',
    XOF: 'bg-cyan-600',
    XAF: 'bg-teal-600'
  }
  return bgs[currency] || 'bg-slate-500'
}

const truncateAddress = (address) => {
  if (!address || address === 'N/A') return 'N/A'
  return `${address.slice(0, 10)}...${address.slice(-4)}`
}

const selectWallet = (wallet) => {
  selectedWallet.value = wallet
}

const openTopUpModal = () => {
    if (!selectedWallet.value && wallets.value.length > 0) {
        selectedWallet.value = wallets.value[0]
    }
    showTopUpModal.value = true
}

const copyAddress = () => {
  if (selectedWallet.value?.address) {
    navigator.clipboard.writeText(selectedWallet.value.address)
    alert('Adresse copiée !')
  }
}

const fetchWallets = async () => {
  try {
    const response = await walletAPI.getAll()
    if (response.data?.wallets) {
      wallets.value = response.data.wallets
    }
  } catch (e) {
    console.log('Using mock data or API error')
     wallets.value = [
      { id: 1, type: 'fiat', currency: 'USD', name: 'Main USD', balance: 1500.50, balanceUSD: 1500.50, status: 'active', address: 'USD-1234-5678' },
      { id: 2, type: 'crypto', currency: 'BTC', name: 'Bitcoin Vault', balance: 0.045, balanceUSD: 1950.00, status: 'active', address: 'bc1q...3k4j' },
    ]
  }
}

const createWallet = async () => {
  creatingWallet.value = true
  try {
    const response = await walletAPI.create({
      currency: newWallet.value.currency,
      name: newWallet.value.name || `Mon Portefeuille ${newWallet.value.currency}`,
      wallet_type: newWallet.value.type
    })
    
    const wallet = response.data?.wallet || response.data || response
    
    if (wallet) {
      wallets.value.push(wallet)
      showCreateWallet.value = false
      newWallet.value = { type: 'fiat', currency: 'USD', name: '' }
      fetchWallets() 
    }
  } catch (e) {
    console.error('Error creating wallet:', e)
  } finally {
    creatingWallet.value = false
  }
}

onMounted(() => {
  fetchWallets()
  if (wallets.value.length > 0) {
    selectedWallet.value = wallets.value[0]
  }
})

definePageMeta({
  layout: false, // Explicitly set layout to false since we use NuxtLayout inside template
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

/* Custom Scrollbar for Modal */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0,0,0,0.05);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(100,116,139,0.5);
  border-radius: 3px;
}
</style>