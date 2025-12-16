<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">Mes Portefeuilles 👛</h1>
          <p class="text-gray-500">Gérez vos devises fiat et crypto</p>
        </div>
        <button @click="showCreateWallet = true" class="btn-premium flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Nouveau Portefeuille
        </button>
      </div>

      <!-- Total Balance -->
      <div class="bg-white rounded-xl shadow-lg p-8 mb-8">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 mb-2">Valeur totale</p>
            <h2 class="text-4xl font-bold text-gray-900 mb-2">{{ formatMoney(totalBalance) }}</h2>
            <p class="text-sm" v-if="totalBalance > 0">
              <span class="text-emerald-600">Calculé en temps réel</span>
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
                <p class="font-semibold text-gray-900">{{ wallet.currency }}</p>
                <p class="text-sm text-gray-900/50">{{ wallet.name }}</p>
              </div>
            </div>
            <span class="badge" :class="wallet.status === 'active' ? 'badge-success' : 'badge-warning'">
              {{ wallet.status }}
            </span>
          </div>

          <div class="mb-4">
            <p class="text-2xl font-bold text-gray-900">
              {{ formatCrypto(wallet.balance, wallet.currency) }}
            </p>
            <p class="text-sm text-gray-900/50">
              ≈ {{ formatMoney(wallet.balanceUSD) }}
            </p>
          </div>

          <div class="flex justify-between text-sm">
            <span class="text-gray-900/50">Adresse</span>
            <span class="text-gray-900 font-mono">{{ truncateAddress(wallet.address) }}</span>
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
              <h3 class="text-xl font-bold text-gray-900">{{ selectedWallet.currency }}</h3>
              <p class="text-gray-900/50">{{ selectedWallet.name }}</p>
            </div>
          </div>
          <div class="flex gap-3">
            <button @click="showDeposit = true" class="btn-success">Recevoir</button>
            <button @click="showSend = true" class="btn-premium">Envoyer</button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Solde</p>
            <p class="text-gray-900 font-bold">{{ formatCrypto(selectedWallet.balance, selectedWallet.currency) }}</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Valeur USD</p>
            <p class="text-gray-900 font-bold">{{ formatMoney(selectedWallet.balanceUSD) }}</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Variation 24h</p>
            <p class="text-emerald-400 font-bold">+5.2%</p>
          </div>
          <div>
            <p class="text-gray-900/50 text-sm mb-1">Réseau</p>
            <p class="text-gray-900">{{ selectedWallet.network || 'Mainnet' }}</p>
          </div>
        </div>

        <!-- Address -->
        <div class="mt-6 p-4 rounded-xl bg-white/5">
          <p class="text-gray-900/50 text-sm mb-2">Adresse de réception</p>
          <div class="flex items-center gap-3">
            <p class="text-gray-900 font-mono text-sm flex-1 truncate">{{ selectedWallet.address }}</p>
            <button @click="copyAddress" class="btn-secondary-premium text-sm py-2 px-3">
              📋 Copier
            </button>
          </div>
        </div>
      </div>

      <!-- Transactions -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">Historique des transactions</h3>
        
        <div class="space-y-3">
          <div v-for="tx in transactions" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl bg-white/5">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center"
                  :class="tx.type === 'receive' ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'">
                {{ tx.type === 'receive' ? '↓' : '↑' }}
              </div>
              <div>
                <p class="font-medium text-gray-900">{{ tx.type === 'receive' ? 'Reçu' : 'Envoyé' }}</p>
                <p class="text-sm text-gray-900/50">{{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold" :class="tx.type === 'receive' ? 'text-green-400' : 'text-gray-900'">
                {{ tx.type === 'receive' ? '+' : '-' }}{{ tx.amount }} {{ tx.currency }}
              </p>
              <p class="text-xs text-gray-900/40">{{ truncateAddress(tx.hash) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Wallet Modal -->
    <div v-if="showCreateWallet" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-2xl p-8 max-w-md w-full mx-4 shadow-xl">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-gray-900">Nouveau Portefeuille</h3>
          <button @click="showCreateWallet = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <form @submit.prevent="createWallet" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Type de devise</label>
            <select v-model="newWallet.type" class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500">
              <option value="fiat">Devise Fiat (USD, EUR...)</option>
              <option value="crypto">Crypto-monnaie (BTC, ETH...)</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Devise</label>
            <select v-model="newWallet.currency" class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500">
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
            <label class="block text-sm font-medium text-gray-700 mb-2">Nom du portefeuille</label>
            <input 
              v-model="newWallet.name" 
              type="text" 
              placeholder="ex: Mon portefeuille principal"
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div class="flex gap-3 mt-6">
            <button 
              type="button" 
              @click="showCreateWallet = false"
              class="flex-1 py-3 px-4 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50"
            >
              Annuler
            </button>
            <button 
              type="submit" 
              :disabled="creatingWallet"
              class="flex-1 py-3 px-4 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 disabled:opacity-50"
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

const createWallet = async () => {
  creatingWallet.value = true
  try {
    const response = await walletAPI.create({
      currency: newWallet.value.currency,
      name: newWallet.value.name || `Mon Portefeuille ${newWallet.value.currency}`,
      wallet_type: newWallet.value.type
    })
    if (response.data) {
      wallets.value.push(response.data.wallet || response.data)
      showCreateWallet.value = false
      newWallet.value = { type: 'fiat', currency: 'USD', name: '' }
      alert('Portefeuille créé avec succès!')
      fetchWallets() // Refresh the list
    }
  } catch (e) {
    console.error('Error creating wallet:', e)
    const error = e as any
    alert(error.response?.data?.error || 'Erreur lors de la création du portefeuille')
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
  layout: false,
  middleware: 'auth'
})
</script>