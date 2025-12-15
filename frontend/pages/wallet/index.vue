<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">Mes Portefeuilles 👛</h1>
          <p class="text-white/60">Gérez vos devises fiat et crypto</p>
        </div>
        <button @click="showCreateWallet = true" class="btn-premium flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Nouveau Portefeuille
        </button>
      </div>

      <!-- Total Balance -->
      <div class="glass-card p-8 mb-8">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-white/60 mb-2">Valeur totale</p>
            <h2 class="text-4xl font-bold text-white mb-2">{{ formatMoney(totalBalance) }}</h2>
            <p class="text-sm">
              <span class="text-emerald-400">+$1,234.56 (5.2%)</span>
              <span class="text-white/40"> aujourd'hui</span>
            </p>
          </div>
          <div class="flex gap-4">
            <button class="btn-success">Déposer</button>
            <button class="btn-secondary-premium">Retirer</button>
          </div>
        </div>
      </div>

      <!-- Wallets Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="wallet in wallets" :key="wallet.id" 
            class="glass-card p-6 cursor-pointer hover:scale-[1.02] transition-transform"
            :class="wallet.type === 'crypto' ? 'stat-card-orange' : 'stat-card-blue'"
            @click="selectWallet(wallet)">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" 
                  :class="getCurrencyBg(wallet.currency)">
                {{ getCurrencyIcon(wallet.currency) }}
              </div>
              <div>
                <p class="font-semibold text-white">{{ wallet.currency }}</p>
                <p class="text-sm text-white/50">{{ wallet.name }}</p>
              </div>
            </div>
            <span class="badge" :class="wallet.status === 'active' ? 'badge-success' : 'badge-warning'">
              {{ wallet.status }}
            </span>
          </div>

          <div class="mb-4">
            <p class="text-2xl font-bold text-white">
              {{ formatCrypto(wallet.balance, wallet.currency) }}
            </p>
            <p class="text-sm text-white/50">
              ≈ {{ formatMoney(wallet.balanceUSD) }}
            </p>
          </div>

          <div class="flex justify-between text-sm">
            <span class="text-white/50">Adresse</span>
            <span class="text-white font-mono">{{ truncateAddress(wallet.address) }}</span>
          </div>
        </div>
      </div>

      <!-- Selected Wallet Details -->
      <div v-if="selectedWallet" class="glass-card p-6 mb-8">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 rounded-2xl flex items-center justify-center text-3xl"
                :class="getCurrencyBg(selectedWallet.currency)">
              {{ getCurrencyIcon(selectedWallet.currency) }}
            </div>
            <div>
              <h3 class="text-xl font-bold text-white">{{ selectedWallet.currency }}</h3>
              <p class="text-white/50">{{ selectedWallet.name }}</p>
            </div>
          </div>
          <div class="flex gap-3">
            <button @click="showDeposit = true" class="btn-success">Recevoir</button>
            <button @click="showSend = true" class="btn-premium">Envoyer</button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <p class="text-white/50 text-sm mb-1">Solde</p>
            <p class="text-white font-bold">{{ formatCrypto(selectedWallet.balance, selectedWallet.currency) }}</p>
          </div>
          <div>
            <p class="text-white/50 text-sm mb-1">Valeur USD</p>
            <p class="text-white font-bold">{{ formatMoney(selectedWallet.balanceUSD) }}</p>
          </div>
          <div>
            <p class="text-white/50 text-sm mb-1">Variation 24h</p>
            <p class="text-emerald-400 font-bold">+5.2%</p>
          </div>
          <div>
            <p class="text-white/50 text-sm mb-1">Réseau</p>
            <p class="text-white">{{ selectedWallet.network || 'Mainnet' }}</p>
          </div>
        </div>

        <!-- Address -->
        <div class="mt-6 p-4 rounded-xl bg-white/5">
          <p class="text-white/50 text-sm mb-2">Adresse de réception</p>
          <div class="flex items-center gap-3">
            <p class="text-white font-mono text-sm flex-1 truncate">{{ selectedWallet.address }}</p>
            <button @click="copyAddress" class="btn-secondary-premium text-sm py-2 px-3">
              📋 Copier
            </button>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-white mb-6">Historique des transactions</h3>
        
        <div class="space-y-3">
          <div v-for="tx in transactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl bg-white/5">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center"
                  :class="tx.type === 'receive' ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'">
                {{ tx.type === 'receive' ? '↓' : '↑' }}
              </div>
              <div>
                <p class="font-medium text-white">{{ tx.type === 'receive' ? 'Reçu' : 'Envoyé' }}</p>
                <p class="text-sm text-white/50">{{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold" :class="tx.type === 'receive' ? 'text-green-400' : 'text-white'">
                {{ tx.type === 'receive' ? '+' : '-' }}{{ tx.amount }} {{ tx.currency }}
              </p>
              <p class="text-xs text-white/40">{{ truncateAddress(tx.hash) }}</p>
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

const wallets = ref([
  { id: '1', currency: 'USD', name: 'Dollar US', balance: 5420.50, balanceUSD: 5420.50, type: 'fiat', status: 'active', address: 'N/A' },
  { id: '2', currency: 'EUR', name: 'Euro', balance: 2340.80, balanceUSD: 2540.12, type: 'fiat', status: 'active', address: 'N/A' },
  { id: '3', currency: 'BTC', name: 'Bitcoin', balance: 0.245, balanceUSD: 11025.50, type: 'crypto', status: 'active', address: '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa', network: 'Bitcoin Mainnet' },
  { id: '4', currency: 'ETH', name: 'Ethereum', balance: 2.45, balanceUSD: 7399.80, type: 'crypto', status: 'active', address: '0x742d35Cc6634C0532925a3b844Bc9e7595f2bD', network: 'Ethereum Mainnet' },
])

const selectedWallet = ref(null)
const showCreateWallet = ref(false)
const showDeposit = ref(false)
const showSend = ref(false)

const transactions = ref([
  { id: '1', type: 'receive', amount: 0.05, currency: 'BTC', date: new Date(), hash: '0x1234...abcd' },
  { id: '2', type: 'send', amount: 500, currency: 'USD', date: new Date(Date.now() - 86400000), hash: 'N/A' },
])

const totalBalance = computed(() => {
  return wallets.value.reduce((sum, w) => sum + w.balanceUSD, 0)
})

const formatMoney = (amount) => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount)
}

const formatCrypto = (amount, currency) => {
  if (['BTC', 'ETH', 'SOL'].includes(currency)) {
    return `${amount.toFixed(6)} ${currency}`
  }
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

const getCurrencyIcon = (currency) => {
  const icons = { BTC: '₿', ETH: 'Ξ', USD: '$', EUR: '€', GBP: '£', SOL: '◎' }
  return icons[currency] || '💰'
}

const getCurrencyBg = (currency) => {
  const bgs = { BTC: 'bg-orange-500', ETH: 'bg-blue-500', USD: 'bg-green-500', EUR: 'bg-indigo-500', SOL: 'bg-purple-500' }
  return bgs[currency] || 'bg-gray-500'
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
  }
}

const fetchWallets = async () => {
  try {
    const response = await walletAPI.getAll()
    if (response.data?.wallets) {
      wallets.value = response.data.wallets
    }
  } catch (e) {
    console.log('Using mock data')
  }
}

onMounted(() => {
  fetchWallets()
  if (wallets.value.length > 0) {
    selectedWallet.value = wallets.value[0]
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>