<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-base mb-2">Mes Portefeuilles 👛</h1>
          <p class="text-muted">Gérez vos devises fiat et crypto</p>
        </div>
        <button @click="showCreateWallet = true" class="btn-primary flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Nouveau Portefeuille
        </button>
      </div>

      <!-- Total Balance -->
      <div class="glass-card mb-8 group">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-muted mb-2 font-medium">Valeur totale</p>
            <h2 class="text-4xl font-bold text-base mb-2 bg-gradient-to-r from-primary-400 to-secondary-400 bg-clip-text text-transparent inline-block">{{ formatMoney(totalBalance) }}</h2>
            <p class="text-sm" v-if="totalBalance > 0">
              <span class="text-success font-medium flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-success animate-pulse"></span>
                Calculé en temps réel
              </span>
            </p>
          </div>
          <div class="flex gap-4">
            <button class="px-6 py-2 rounded-xl bg-success/10 text-success border border-success/20 font-medium hover:bg-success/20 transition-all">Déposer</button>
            <button class="px-6 py-2 rounded-xl bg-secondary-100 dark:bg-secondary-800 text-base border border-secondary-200 dark:border-secondary-700 font-medium hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-all">Retirer</button>
          </div>
        </div>
      </div>

      <!-- Wallets Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="wallet in wallets" :key="wallet.id" 
            class="glass-card p-6 cursor-pointer hover:scale-[1.02] transition-transform relative overflow-hidden group border-l-4"
            :class="wallet.type === 'crypto' ? 'border-l-orange-500' : 'border-l-blue-500'"
            @click="selectWallet(wallet)">
            
           <!-- Background Glow -->
           <div class="absolute -right-4 -top-4 w-24 h-24 rounded-full blur-2xl opacity-10 transition-opacity group-hover:opacity-20"
                :class="wallet.type === 'crypto' ? 'bg-orange-500' : 'bg-blue-500'"></div>

          <div class="flex items-center justify-between mb-4 relative z-10">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shadow-lg text-white" 
                  :class="getCurrencyBg(wallet.currency)">
                {{ getCurrencyIcon(wallet.currency) }}
              </div>
              <div>
                <p class="font-bold text-base">{{ wallet.currency }}</p>
                <p class="text-xs text-muted font-medium uppercase tracking-wider">{{ wallet.name }}</p>
              </div>
            </div>
            <span class="px-2 py-1 rounded-full text-xs font-semibold border" 
                  :class="wallet.status === 'active' ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20' : 'bg-amber-500/10 text-amber-500 border-amber-500/20'">
              {{ wallet.status }}
            </span>
          </div>

          <div class="mb-4 relative z-10">
            <p class="text-2xl font-bold text-base tracking-tight">
              {{ formatCrypto(wallet.balance, wallet.currency) }}
            </p>
            <p class="text-sm text-muted font-medium">
              ≈ {{ formatMoney(wallet.balanceUSD) }}
            </p>
          </div>

          <div class="flex justify-between text-sm relative z-10 pt-4 border-t border-secondary-100 dark:border-secondary-800">
            <span class="text-muted">Adresse</span>
            <span class="text-base font-mono">{{ truncateAddress(wallet.address) }}</span>
          </div>
        </div>
      </div>

      <!-- Selected Wallet Details -->
      <div v-if="selectedWallet" class="glass-card p-8 mb-8">
        <div class="flex items-center justify-between mb-8">
          <div class="flex items-center gap-5">
            <div class="w-16 h-16 rounded-2xl flex items-center justify-center text-3xl shadow-lg text-white"
                :class="getCurrencyBg(selectedWallet.currency)">
              {{ getCurrencyIcon(selectedWallet.currency) }}
            </div>
            <div>
              <h3 class="text-2xl font-bold text-base">{{ selectedWallet.currency }}</h3>
              <p class="text-muted">{{ selectedWallet.name }}</p>
            </div>
          </div>
          <div class="flex gap-3">
            <button @click="showDeposit = true" class="px-6 py-2 rounded-xl bg-success text-white font-medium shadow-lg shadow-success/20 hover:bg-success-600 transition-all">Recevoir</button>
            <button @click="showSend = true" class="btn-primary">Envoyer</button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 p-6 rounded-xl bg-surface-hover/50 border border-secondary-100 dark:border-secondary-800/50">
          <div>
            <p class="text-muted text-sm mb-1 font-medium">Solde</p>
            <p class="text-base font-bold text-lg">{{ formatCrypto(selectedWallet.balance, selectedWallet.currency) }}</p>
          </div>
          <div>
            <p class="text-muted text-sm mb-1 font-medium">Valeur USD</p>
            <p class="text-base font-bold text-lg">{{ formatMoney(selectedWallet.balanceUSD) }}</p>
          </div>
          <div>
            <p class="text-muted text-sm mb-1 font-medium">Variation 24h</p>
            <p class="text-success font-bold text-lg">+5.2%</p>
          </div>
          <div>
            <p class="text-muted text-sm mb-1 font-medium">Réseau</p>
            <p class="text-base font-medium">{{ selectedWallet.network || 'Mainnet' }}</p>
          </div>
        </div>

        <!-- Address -->
        <div class="mt-8">
          <p class="text-muted text-sm mb-2 font-medium">Adresse de réception</p>
          <div class="flex items-center gap-3 p-4 rounded-xl bg-surface border border-secondary-200 dark:border-secondary-800">
            <p class="text-base font-mono text-sm flex-1 truncate">{{ selectedWallet.address }}</p>
            <button @click="copyAddress" class="text-primary hover:text-primary-600 font-medium text-sm transition-colors">
              📋 Copier
            </button>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-bold text-base mb-6">Historique des transactions</h3>
        
        <div class="space-y-3">
          <div v-for="tx in transactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl hover:bg-surface-hover transition-colors border border-transparent hover:border-secondary-100 dark:hover:border-secondary-800">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full flex items-center justify-center font-bold"
                  :class="tx.type === 'receive' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-rose-500/10 text-rose-500'">
                {{ tx.type === 'receive' ? '↓' : '↑' }}
              </div>
              <div>
                <p class="font-bold text-base">{{ tx.type === 'receive' ? 'Reçu' : 'Envoyé' }}</p>
                <p class="text-sm text-muted">{{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold" :class="tx.type === 'receive' ? 'text-emerald-500' : 'text-base'">
                {{ tx.type === 'receive' ? '+' : '-' }}{{ tx.amount }} {{ tx.currency }}
              </p>
              <p class="text-xs text-muted font-mono">{{ truncateAddress(tx.hash) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Wallet Modal -->
    <div v-if="showCreateWallet" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div class="bg-surface rounded-2xl p-8 max-w-md w-full mx-4 shadow-2xl border border-secondary-200 dark:border-secondary-700">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-base">Nouveau Portefeuille</h3>
          <button @click="showCreateWallet = false" class="text-muted hover:text-base transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <form @submit.prevent="createWallet" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-muted mb-2">Type de devise</label>
            <select v-model="newWallet.type" class="input-field">
              <option value="fiat">Devise Fiat (USD, EUR...)</option>
              <option value="crypto">Crypto-monnaie (BTC, ETH...)</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-muted mb-2">Devise</label>
            <select v-model="newWallet.currency" class="input-field">
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
            <label class="block text-sm font-medium text-muted mb-2">Nom du portefeuille</label>
            <input 
              v-model="newWallet.name" 
              type="text" 
              placeholder="ex: Mon portefeuille principal"
              class="input-field"
            />
          </div>

          <div class="flex gap-3 mt-8">
            <button 
              type="button" 
              @click="showCreateWallet = false"
              class="flex-1 py-3 px-4 border border-secondary-200 dark:border-secondary-700 text-muted hover:text-base rounded-xl hover:bg-surface-hover transition-colors font-medium"
            >
              Annuler
            </button>
            <button 
              type="submit" 
              :disabled="creatingWallet"
              class="flex-1 btn-primary"
            >
              {{ creatingWallet ? 'Création...' : 'Créer le portefeuille' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { walletAPI } from '~/composables/useApi'

// Wallets - will be loaded from API
const wallets = ref([])

const selectedWallet = ref(null)
const showCreateWallet = ref(false)
const showDeposit = ref(false)
const showSend = ref(false)
const creatingWallet = ref(false)

// New wallet form
const newWallet = ref({
  type: 'fiat',
  currency: 'USD',
  name: ''
})

// Transactions - will be loaded from API
const transactions = ref([])

const totalBalance = computed(() => {
  return wallets.value.reduce((sum, w) => {
    const val = Number(w.balanceUSD)
    return sum + (isNaN(val) ? 0 : val)
  }, 0)
})

const formatMoney = (amount) => {
  const val = Number(amount)
  // Handle NaN or undefined/null gracefully
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

const formatDate = (date) => {
  if (!date) return ''
  return new Intl.DateTimeFormat('fr-FR', { 
    day: '2-digit', 
    month: 'short', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(date))
}

const getCurrencyIcon = (currency) => {
  const icons = { BTC: '₿', ETH: 'Ξ', USD: '$', EUR: '€', GBP: '£', SOL: '◎' }
  return icons[currency] || '💰'
}

const getCurrencyBg = (currency) => {
  // Use Tailwind colors directly
  const bgs = { 
    BTC: 'bg-orange-500', 
    ETH: 'bg-blue-500', 
    USD: 'bg-emerald-500', 
    EUR: 'bg-indigo-500', 
    SOL: 'bg-purple-500',
    GBP: 'bg-rose-500',
    XOF: 'bg-cyan-600',
    XAF: 'bg-teal-600'
  }
  return bgs[currency] || 'bg-slate-500'
}

const truncateAddress = (address) => {
  if (!address || address === 'N/A') return 'N/A'
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

const selectWallet = (wallet) => {
  selectedWallet.value = wallet
}

const copyAddress = () => {
  if (selectedWallet.value?.address) {
    navigator.clipboard.writeText(selectedWallet.value.address)
    // Optional: show toast
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
    // Mock data for dev
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
    
    // Handle different response structures gracefully
    const wallet = response.data?.wallet || response.data || response
    
    if (wallet) {
      wallets.value.push(wallet)
      showCreateWallet.value = false
      newWallet.value = { type: 'fiat', currency: 'USD', name: '' }
      // alert('Portefeuille créé avec succès!') // Removed alert for smoother UX
      fetchWallets() 
    }
  } catch (e) {
    console.error('Error creating wallet:', e)
    // alert(e?.response?.data?.error || 'Erreur lors de la création du portefeuille')
  } finally {
    creatingWallet.value = false
  }
}

onMounted(() => {
  fetchWallets()
  if (wallets.value.length > 0) {
    selectedWallet.value = wallets.value[0]
  }
  
  // Mock transactions
  transactions.value = [
    { id: 1, type: 'receive', amount: 0.01, currency: 'BTC', date: new Date(), hash: '0x123...abc' },
    { id: 2, type: 'send', amount: 50, currency: 'USD', date: new Date(Date.now() - 86400000), hash: 'TX-999-888' }
  ]
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