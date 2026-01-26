<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-base mb-2">₿ Acheter / Vendre Crypto</h1>
          <p class="text-muted">Achetez et vendez des cryptomonnaies instantanément</p>
        </div>
        <NuxtLink to="/exchange" class="text-primary hover:text-primary/80 flex items-center gap-2 font-medium">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
          Retour
        </NuxtLink>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Exchange Form -->
        <div class="lg:col-span-2">
          <div class="glass-card p-6">
            <!-- Buy/Sell Toggle -->
            <div class="flex gap-2 mb-6">
              <button
                @click="mode = 'buy'"
                :class="mode === 'buy' ? 'bg-success text-white' : 'bg-white/10 text-muted hover:text-base'"
                class="flex-1 py-3 rounded-xl font-semibold transition-colors"
              >
                🛒 Acheter
              </button>
              <button
                @click="mode = 'sell'"
                :class="mode === 'sell' ? 'bg-error text-white' : 'bg-white/10 text-muted hover:text-base'"
                class="flex-1 py-3 rounded-xl font-semibold transition-colors"
              >
                💰 Vendre
              </button>
            </div>
            
            <!-- FROM Section -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-muted mb-2">
                {{ mode === 'buy' ? '💳 Vous payez avec' : '📤 Vous vendez' }}
              </label>
              
              <!-- Wallet Selector -->
              <div class="mb-3">
                <select 
                  v-model="selectedFromWalletId" 
                  @change="onWalletChange"
                  class="input-premium w-full bg-slate-800 text-white border-slate-600"
                >
                  <option value="" disabled>Sélectionner un portefeuille</option>
                  <option v-for="wallet in fromWalletOptions" :key="wallet.id" :value="wallet.id">
                    {{ getCurrencyEmoji(wallet.currency) }} {{ wallet.currency }} - 
                    Solde: {{ formatBalance(wallet.balance) }}
                  </option>
                </select>
              </div>
              
              <!-- Amount Input -->
              <div class="relative">
                <input 
                  v-model.number="fromAmount" 
                  type="number" 
                  placeholder="Entrez le montant"
                  class="input-premium w-full text-2xl font-bold bg-slate-800 text-white border-slate-600 pr-20 h-16"
                  @input="calculateConversion"
                />
                <span class="absolute right-4 top-1/2 -translate-y-1/2 text-lg font-semibold text-muted">
                  {{ fromCurrency }}
                </span>
              </div>
              
              <p class="text-xs text-muted mt-2">
                Solde disponible: <span class="text-success font-medium">{{ formatBalance(getWalletBalance(fromCurrency)) }} {{ fromCurrency }}</span>
              </p>
            </div>

            <!-- Swap Button -->
            <div class="flex justify-center my-4">
              <button @click="swapMode" class="p-3 rounded-xl bg-primary/20 hover:bg-primary/30 transition-colors">
                <svg class="w-6 h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
                </svg>
              </button>
            </div>

            <!-- TO Section -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-muted mb-2">
                {{ mode === 'buy' ? '📥 Vous recevez' : '💵 Vous recevez' }}
              </label>
              
              <!-- Crypto/Fiat Selector -->
              <div class="mb-3">
                <select v-model="toCurrency" @change="calculateConversion" class="input-premium w-full bg-slate-800 text-white border-slate-600">
                  <template v-if="mode === 'buy'">
                    <option v-for="crypto in cryptoCurrencies" :key="crypto.symbol" :value="crypto.symbol">
                      {{ crypto.icon }} {{ crypto.symbol }} - {{ crypto.name }}
                    </option>
                  </template>
                  <template v-else>
                    <option v-for="wallet in fiatWallets" :key="wallet.id" :value="wallet.currency">
                      {{ getCurrencyEmoji(wallet.currency) }} {{ wallet.currency }}
                    </option>
                  </template>
                </select>
              </div>
              
              <!-- Destination Wallet Selector -->
              <div class="mb-3">
                <select 
                  v-model="selectedToWalletId" 
                  class="input-premium w-full bg-slate-800 text-white border-slate-600"
                >
                  <option v-for="wallet in toWalletOptions" :key="wallet.id" :value="wallet.id">
                    {{ wallet.id === 'create_new' ? '➕ ' : '' }}
                    {{ wallet.name || wallet.currency }} 
                    {{ wallet.balance !== undefined && wallet.id !== 'create_new' ? `(${formatBalance(wallet.balance)} ${wallet.currency})` : '' }}
                  </option>
                </select>
              </div>

              <!-- Amount Display -->
              <div class="relative">
                <input 
                  v-model="toAmount" 
                  type="text" 
                  placeholder="0.00"
                  class="input-premium w-full text-2xl font-bold bg-slate-700 text-success border-slate-600 pr-20 h-16"
                  readonly
                />
                <span class="absolute right-4 top-1/2 -translate-y-1/2 text-lg font-semibold text-muted">
                  {{ toCurrency }}
                </span>
              </div>
            </div>

            <!-- Rate Info -->
            <div class="p-4 rounded-xl bg-slate-800/50 border border-slate-700 mb-6">
              <div class="flex justify-between text-sm mb-2">
                <span class="text-muted">📊 Taux</span>
                <span class="text-base font-medium">
                  1 {{ mode === 'buy' ? toCurrency : fromCurrency }} = 
                  {{ formatMoney(getCurrentRate(), mode === 'buy' ? fromCurrency : toCurrency) }}
                </span>
              </div>
              <div class="flex justify-between text-sm mb-2">
                <span class="text-muted">💸 Frais ({{ feePercentage }}%)</span>
                <span class="text-base">{{ formatMoney(calculatedFee, mode === 'buy' ? fromCurrency : toCurrency) }}</span>
              </div>
              <div class="flex justify-between text-sm font-bold pt-2 border-t border-slate-700">
                <span class="text-muted">✨ Vous recevez</span>
                <span class="text-success text-lg">{{ toAmount }} {{ toCurrency }}</span>
              </div>
            </div>

            <!-- Execute Button -->
            <button 
              @click="executeExchange"
              :disabled="loading || !fromAmount || fromAmount <= 0 || fromAmount > getWalletBalance(fromCurrency) || !selectedFromWalletId"
              class="btn-premium w-full py-4 text-lg font-bold disabled:opacity-50"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                Traitement...
              </span>
              <span v-else>
                {{ mode === 'buy' ? `🛒 Acheter ${toCurrency}` : `💰 Vendre ${fromCurrency}` }}
              </span>
            </button>

            <p v-if="!selectedFromWalletId" class="text-warning text-sm mt-2 text-center">
              ⚠️ Veuillez sélectionner un portefeuille
            </p>
            <p v-else-if="fromAmount > getWalletBalance(fromCurrency)" class="text-error text-sm mt-2 text-center">
              ❌ Solde insuffisant
            </p>
          </div>

          <!-- Crypto Address Section (for receiving crypto) -->
          <div v-if="mode === 'buy' && selectedCryptoWallet" class="glass-card p-6 mt-6">
            <h3 class="text-lg font-semibold text-base mb-4">📬 Adresse de réception {{ toCurrency }}</h3>
            
            <div v-if="getCryptoWalletAddress()" class="space-y-4">
              <div class="p-4 bg-slate-800 rounded-xl">
                <p class="text-xs text-muted mb-2">Votre adresse {{ toCurrency }}:</p>
                <div class="flex items-center gap-2">
                  <code class="text-sm text-success break-all flex-1">{{ getCryptoWalletAddress() }}</code>
                  <button @click="copyAddress" class="p-2 bg-primary/20 rounded-lg hover:bg-primary/30">
                    <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                    </svg>
                  </button>
                </div>
              </div>
              <p class="text-xs text-muted">
                💡 Utilisez cette adresse pour recevoir des {{ toCurrency }} depuis un autre portefeuille ou exchange.
              </p>
            </div>
            
            <div v-else class="text-center py-4">
              <p class="text-muted mb-4">Aucune adresse {{ toCurrency }} générée</p>
              <button 
                @click="generateAddress" 
                :disabled="generatingAddress"
                class="btn-secondary px-6 py-2"
              >
                <span v-if="generatingAddress">Génération...</span>
                <span v-else>🔑 Générer une adresse</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Market Prices -->
        <div class="glass-card p-6 h-fit">
          <h2 class="text-lg font-semibold text-base mb-6">📊 Prix du marché</h2>
          <div class="space-y-3 max-h-[500px] overflow-y-auto custom-scrollbar">
            <div v-for="crypto in markets" :key="crypto.symbol" 
                class="flex items-center justify-between p-3 rounded-xl bg-slate-800/50 border border-slate-700 hover:bg-slate-700/50 transition-colors cursor-pointer"
                @click="selectCrypto(crypto.symbol)">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold text-sm">{{ crypto.icon }}</span>
                </div>
                <div>
                  <p class="font-medium text-base">{{ crypto.name }}</p>
                  <p class="text-xs text-muted">{{ crypto.symbol }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-base">${{ formatPrice(crypto.price) }}</p>
                <p class="text-xs" :class="(crypto.change || 0) >= 0 ? 'text-success' : 'text-error'">
                  {{ (crypto.change || 0) >= 0 ? '+' : '' }}{{ Number(crypto.change || 0).toFixed(2) }}%
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Notification Modal -->
      <Teleport to="body">
        <Transition name="modal-fade">
          <div v-if="notification.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" @click="notification.show = false"></div>
            <div class="relative w-full max-w-md transform transition-all">
              <div class="glass-card p-6 text-center">
                <div class="mb-4">
                  <div v-if="notification.type === 'success'" class="w-16 h-16 mx-auto rounded-full bg-success/20 flex items-center justify-center">
                    <svg class="w-8 h-8 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                    </svg>
                  </div>
                  <div v-else class="w-16 h-16 mx-auto rounded-full bg-error/20 flex items-center justify-center">
                    <svg class="w-8 h-8 text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                    </svg>
                  </div>
                </div>
                <h3 class="text-xl font-bold text-base mb-2">
                  {{ notification.type === 'success' ? 'Succès!' : 'Erreur' }}
                </h3>
                <p class="text-muted mb-6">{{ notification.message }}</p>
                <button
                  @click="notification.show = false"
                  :class="notification.type === 'success' ? 'bg-success hover:bg-success/80' : 'bg-error hover:bg-error/80'"
                  class="w-full py-3 px-6 rounded-xl text-white font-medium transition-colors"
                >
                  Compris
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { walletAPI, exchangeAPI } from '~/composables/useApi'
import { usePin } from '~/composables/usePin'
import { useWalletStore } from '~/stores/wallet'
import { useExchangeStore } from '~/stores/exchange'
import { storeToRefs } from 'pinia'

// Page meta
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { requirePin } = usePin()
const walletStore = useWalletStore()
const exchangeStore = useExchangeStore()

// Store refs
const { wallets } = storeToRefs(walletStore)
const { cryptoRates } = storeToRefs(exchangeStore)

// Mode: buy or sell
const mode = ref('buy')

// Reactive data
const loading = ref(false)
const fromCurrency = ref('USD')
const toCurrency = ref('BTC')
const fromAmount = ref(100)
const toAmount = ref('0')
const feePercentage = ref(0.5)
const calculatedFee = ref(0)
const generatingAddress = ref(false)

// Wallet selection
const selectedFromWalletId = ref('')
const selectedToWalletId = ref('create_new') // Default to create new or select existing

// Computed wallet options
const fromWalletOptions = computed(() => {
  if (mode.value === 'buy') {
    // When buying crypto, pay with fiat
    return wallets.value.filter(w => w.wallet_type === 'fiat')
  } else {
    // When selling crypto, pay with crypto
    return wallets.value.filter(w => w.wallet_type === 'crypto')
  }
})

// Destination wallet options based on selected currency
const toWalletOptions = computed(() => {
    // Current target currency
    const target = mode.value === 'buy' ? toCurrency.value : toCurrency.value

  if (mode.value === 'buy') {
    const existing = wallets.value.filter(w => w.currency === target && w.wallet_type === 'crypto')
    return [
      ...existing,
      { id: 'create_new', name: '➕ Créer un nouveau portefeuille', currency: target, balance: 0 }
    ]
  } else {
    // Selling crypto -> receive fiat
    const existing = wallets.value.filter(w => w.currency === target && w.wallet_type === 'fiat')
    // Usually we don't creating fiat wallets on the fly as easily, but let's allow selecting existing
    return existing
  }
})

const fiatWallets = computed(() => wallets.value.filter(w => w.wallet_type === 'fiat'))

// Selected crypto wallet (for address display)
const selectedCryptoWallet = computed(() => {
  if (mode.value === 'buy') {
    // If specific wallet selected, show it. If create_new, show nothing or indication
    if (selectedToWalletId.value && selectedToWalletId.value !== 'create_new') {
        return wallets.value.find(w => w.id === selectedToWalletId.value)
    }
    // Fallback to finding one if not explicitly selected logic wasn't there
    return wallets.value.find(w => w.currency === toCurrency.value && w.wallet_type === 'crypto')
  }
  return null
})

// Notification
const notification = ref({
  show: false,
  type: 'success',
  message: ''
})

const showNotification = (type, message) => {
  notification.value = { show: true, type, message }
}

// Market Data from Store
// Map store rates to list usable by UI
const markets = computed(() => {
    if (cryptoRates.value.length === 0) return []
    return cryptoRates.value.map(r => {
        // Parse symbol, e.g., "BTC/USD" -> "BTC"
        let symbol = r.symbol || r.pair || ''
        if (!symbol && r.from_currency && r.to_currency) {
            symbol = `${r.from_currency}/${r.to_currency}`
        }
        if (!symbol) return null // Skip invalid entries
        
        const base = symbol.includes('/') ? symbol.split('/')[0] : symbol
        
        return {
            symbol: base,
            name: getCurrencyName(base),
            icon: getCurrencySymbol(base), 
            price: r.price || r.rate,
            change: r.change_24h || 0,
            bgColor: getBgColor(base)
        }
    }).filter(m => m !== null) // Filter out nulls
})

const cryptoCurrencies = computed(() => markets.value)

const getBgColor = (symbol) => {
  const colors = {
    BTC: 'bg-orange-500',
    ETH: 'bg-blue-500',
    SOL: 'bg-purple-500',
    XRP: 'bg-slate-500',
    BNB: 'bg-yellow-500',
    ADA: 'bg-blue-600',
    DOGE: 'bg-amber-400',
    DOT: 'bg-pink-500',
    USDT: 'bg-green-500',
    USDC: 'bg-blue-400'
  }
  return colors[symbol] || 'bg-indigo-500'
}

const getCurrencyEmoji = (currency) => {
  const emojis = {
    USD: '🇺🇸', EUR: '🇪🇺', GBP: '🇬🇧', XOF: '💰', XAF: '💰',
    NGN: '🇳🇬', KES: '🇰🇪', ZAR: '🇿🇦', MAD: '🇲🇦', CAD: '🇨🇦',
    AUD: '🇦🇺', JPY: '🇯🇵', CHF: '🇨🇭', CNY: '🇨🇳', INR: '🇮🇳',
    BTC: '₿', ETH: 'Ξ', SOL: '◎', USDT: '₮', USDC: '$'
  }
  return emojis[currency] || '💵'
}

const getCurrencySymbol = (currency) => {
  const symbols = {
    'BTC': '₿', 'ETH': 'Ξ', 'LTC': 'Ł', 'ADA': '₳', 
    'USD': '$', 'EUR': '€', 'GBP': '£'
  }
  return symbols[currency] || currency.substring(0, 2)
}

const getCurrencyName = (currency) => {
  const names = {
    'BTC': 'Bitcoin', 'ETH': 'Ethereum', 'LTC': 'Litecoin', 
    'USD': 'US Dollar', 'EUR': 'Euro', 'SOL': 'Solana', 'XOF': 'Franc CFA'
  }
  return names[currency] || currency
}

const getWalletBalance = (currency) => {
  const wallet = wallets.value.find(w => w.currency === currency)
  return wallet?.balance || 0
}

const getCurrentRate = () => {
    // Rely on Exchange Store getter
    const from = mode.value === 'buy' ? fromCurrency.value : fromCurrency.value // If buy: paying with fromCurrency (USD). rate is 1? No.
    // Rate display: 1 Target = X Source
    // If Buying BTC with USD: Display 1 BTC = 50000 USD
    
    if (mode.value === 'buy') {
        const target = toCurrency.value
        const source = fromCurrency.value
        return exchangeStore.getRate(target, source)
    } else {
        // Selling BTC for USD: Display 1 BTC = 50000 USD
        const target = fromCurrency.value // BTC
        const source = toCurrency.value // USD
        return exchangeStore.getRate(target, source)
    }
}

const getRateForConversion = () => {
    // If buying: Source (USD) -> Target (BTC)
    // AmountFrom (USD) / Rate (USD/BTC) ? No.
    // We typically calculate: Amount * Price.
    // Price of 1 Target in Source Terms.
    if (mode.value === 'buy') {
         // Buying BTC with USD. Price is e.g. 50000.
         // 100 USD / 50000 = 0.002 BTC.
         return exchangeStore.getRate(toCurrency.value, fromCurrency.value)
    } else {
        // Selling BTC for USD.
        // 1 BTC * 50000 = 50000 USD.
        return exchangeStore.getRate(fromCurrency.value, toCurrency.value)
    }
}

const formatBalance = (amount) => {
  return Number(amount || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 8 })
}

const formatMoney = (amount, currency = 'USD') => {
  try {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)
  } catch {
    return `${Number(amount).toFixed(2)} ${currency}`
  }
}

const formatPrice = (price) => {
  if (price == null) return '0'
  return Number(price) >= 1 ? Number(price).toLocaleString() : Number(price).toFixed(6)
}

const onWalletChange = () => {
  const wallet = wallets.value.find(w => w.id === selectedFromWalletId.value)
  if (wallet) {
    fromCurrency.value = wallet.currency
    calculateConversion()
  }
}

const calculateConversion = () => {
    // Using store rate
    const rate = getRateForConversion() // Price of 1 item in counter currency
    
    const fee = fromAmount.value * (feePercentage.value / 100)
    calculatedFee.value = fee
    const netAmount = fromAmount.value - fee

    if (mode.value === 'buy') {
        // Buying. fromAmount is in Fiat (e.g. USD). Rate is BTC in USD (e.g. 50000).
        // Result is BTC.
        // USD / (USD/BTC) = BTC
        if (rate > 0) {
            toAmount.value = (netAmount / rate).toFixed(8)
        } else {
            toAmount.value = '0'
        }
        
        // Auto-select destination wallet if exists
        const matching = wallets.value.find(w => w.currency === toCurrency.value && w.wallet_type === 'crypto')
        selectedToWalletId.value = matching ? matching.id : 'create_new'
    } else {
        // Selling. fromAmount is in Crypto (e.g. BTC). Rate is BTC in USD (e.g. 50000).
        // Result is USD.
        // BTC * (USD/BTC) = USD
         toAmount.value = (netAmount * rate).toFixed(2)
         
        // Auto-select destination fiat wallet
        const matching = wallets.value.find(w => w.currency === toCurrency.value && w.wallet_type === 'fiat')
        selectedToWalletId.value = matching ? matching.id : ''
    }
}

const swapMode = () => {
  mode.value = mode.value === 'buy' ? 'sell' : 'buy'
}

const selectCrypto = (symbol) => {
  if (mode.value === 'buy') {
    toCurrency.value = symbol
  } else {
    // For sell, need to select wallet with this crypto
    const cryptoWallet = wallets.value.find(w => w.currency === symbol)
    if (cryptoWallet) {
      selectedFromWalletId.value = cryptoWallet.id
      fromCurrency.value = symbol
    }
  }
  calculateConversion()
}

const copyAddress = () => {
  const address = getCryptoWalletAddress()
  if (address) {
    navigator.clipboard.writeText(address)
    showNotification('success', 'Adresse copiée!')
  }
}

const getCryptoWalletAddress = () => {
  if (!selectedCryptoWallet.value) return null
  // Handle different possible field names
  return selectedCryptoWallet.value.wallet_address || 
         selectedCryptoWallet.value.address || 
         selectedCryptoWallet.value.WalletAddress || 
         null
}

const generateAddress = async () => {
  generatingAddress.value = true
  try {
    // API: POST /wallet-service/api/v1/wallets/{id}/address
    // This creates a new deposit address for the crypto wallet
    let wallet = selectedCryptoWallet.value
    
    if (!wallet) {
      // Create the crypto wallet first
      const { data: newWalletRes } = await walletAPI.createWallet({
        currency: toCurrency.value,
        name: `${toCurrency.value} Wallet`,
        wallet_type: 'crypto'
      })
      wallet = newWalletRes.wallet || newWalletRes
    }
    
    // Generate address (this endpoint needs to exist in wallet-service)
    // For now, simulate with a placeholder
    showNotification('success', `Portefeuille ${toCurrency.value} créé! L'adresse sera générée automatiquement.`)
    
    // Refresh wallets
    await walletStore.fetchWallets()
  } catch (error) {
    showNotification('error', 'Erreur lors de la génération de l\'adresse')
  } finally {
    generatingAddress.value = false
  }
}

const executeExchange = async () => {
  if (!selectedFromWalletId.value) {
    showNotification('error', 'Veuillez sélectionner un portefeuille source')
    return
  }
  
  if (!selectedToWalletId.value && mode.value === 'sell') {
      showNotification('error', 'Veuillez sélectionner un portefeuille de destination')
      return
  }

  await requirePin(async () => {
    loading.value = true
    try {
      let destWalletId = selectedToWalletId.value
      
      // Determine Destination Wallet
      if (mode.value === 'buy') {
          if (destWalletId === 'create_new') {
             // Create new wallet
            const cryptoData = markets.value.find(c => c.symbol === toCurrency.value)
            const walletName = cryptoData ? `${cryptoData.name} Wallet` : `${toCurrency.value} Wallet`
            
            const { data: newWalletRes } = await walletAPI.createWallet({
              currency: toCurrency.value,
              name: walletName,
              wallet_type: 'crypto'
            })
            const newWallet = newWalletRes.wallet || newWalletRes
            destWalletId = newWallet.id
          }
      } 
      // If sell mode, destWalletId is already selected (fiat wallet) or user must select one
      
      if (!destWalletId) {
          throw new Error("Portefeuille de destination manquant")
      }

      // 1. Get Quote
      const { data: quoteRes } = await exchangeAPI.getQuote(
        fromCurrency.value, 
        toCurrency.value, 
        fromAmount.value
      )
      const quote = quoteRes.quote
      
      if (!quote?.id) {
        throw new Error('Échec de la création du devis')
      }
      
      // 2. Execute Exchange
      await exchangeAPI.executeExchange(
        quote.id,
        selectedFromWalletId.value,
        destWalletId
      )
      
      showNotification('success', 
        mode.value === 'buy' 
          ? `🎉 Vous avez acheté ${toAmount.value} ${toCurrency.value}!`
          : `💰 Vous avez vendu ${fromAmount.value} ${fromCurrency.value}!`
      )
      
      // Refresh wallets
      await walletStore.fetchWallets()
      
      // Reset form (keep small amount to avoid empty input)
      fromAmount.value = 100
      toAmount.value = '0'
      calculateConversion()
      
    } catch (error) {
      console.error('Exchange error:', error)
      showNotification('error', error.response?.data?.error || error.message || 'Erreur lors de l\'échange')
    } finally {
      loading.value = false
    }
  })
}

// Watch mode changes
watch(mode, () => {
  selectedFromWalletId.value = ''
  if (fromWalletOptions.value.length > 0) {
    selectedFromWalletId.value = fromWalletOptions.value[0].id
    fromCurrency.value = fromWalletOptions.value[0].currency
  }
  
  if (mode.value === 'buy') {
    toCurrency.value = 'BTC'
  } else {
    toCurrency.value = fiatWallets.value[0]?.currency || 'USD'
  }
  calculateConversion()
})

// Watch rates to recalculate if they change
watch(cryptoRates, () => {
    calculateConversion()
})

onMounted(async () => {
  // Initialize stores
  await Promise.all([
      walletStore.fetchWallets(),
      exchangeStore.fetchRates()
  ])
  
  // Set defaults
  if (fromWalletOptions.value.length > 0 && !selectedFromWalletId.value) {
      selectedFromWalletId.value = fromWalletOptions.value[0].id
      fromCurrency.value = fromWalletOptions.value[0].currency
  }
  
  calculateConversion()
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

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.5);
  border-radius: 2px;
}

/* Modal animations */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: all 0.3s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
.modal-fade-enter-from .glass-card,
.modal-fade-leave-to .glass-card {
  transform: scale(0.9);
}
</style>
