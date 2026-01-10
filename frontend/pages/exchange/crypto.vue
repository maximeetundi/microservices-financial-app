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
                Acheter
              </button>
              <button
                @click="mode = 'sell'"
                :class="mode === 'sell' ? 'bg-error text-white' : 'bg-white/10 text-muted hover:text-base'"
                class="flex-1 py-3 rounded-xl font-semibold transition-colors"
              >
                Vendre
              </button>
            </div>
            
            <!-- From (Fiat for buy, Crypto for sell) -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-muted mb-2">
                {{ mode === 'buy' ? 'Vous payez' : 'Vous vendez' }}
              </label>
              <div class="flex gap-4">
                <div class="flex-1">
                  <input 
                    v-model.number="fromAmount" 
                    type="number" 
                    placeholder="0.00"
                    class="input-premium w-full text-2xl font-bold"
                    @input="calculateConversion"
                  />
                </div>
                <!-- Currency selector -->
                <select v-model="fromCurrency" @change="onFromCurrencyChange" class="input-premium w-40">
                  <template v-if="mode === 'buy'">
                    <option v-for="wallet in fiatWallets" :key="wallet.id" :value="wallet.currency">
                      {{ getCurrencyEmoji(wallet.currency) }} {{ wallet.currency }}
                    </option>
                  </template>
                  <template v-else>
                    <option v-for="crypto in cryptoCurrencies" :key="crypto.symbol" :value="crypto.symbol">
                      {{ crypto.icon }} {{ crypto.symbol }}
                    </option>
                  </template>
                </select>
              </div>
              <p class="text-xs text-muted mt-2">
                Solde: {{ formatBalance(getWalletBalance(fromCurrency)) }} {{ fromCurrency }}
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

            <!-- To (Crypto for buy, Fiat for sell) -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-muted mb-2">
                {{ mode === 'buy' ? 'Vous recevez' : 'Vous recevez' }}
              </label>
              <div class="flex gap-4">
                <div class="flex-1">
                  <input 
                    v-model="toAmount" 
                    type="text" 
                    placeholder="0.00"
                    class="input-premium w-full text-2xl font-bold"
                    readonly
                  />
                </div>
                <select v-model="toCurrency" @change="calculateConversion" class="input-premium w-40">
                  <template v-if="mode === 'buy'">
                    <option v-for="crypto in cryptoCurrencies" :key="crypto.symbol" :value="crypto.symbol">
                      {{ crypto.icon }} {{ crypto.symbol }}
                    </option>
                  </template>
                  <template v-else>
                    <option v-for="wallet in fiatWallets" :key="wallet.id" :value="wallet.currency">
                      {{ getCurrencyEmoji(wallet.currency) }} {{ wallet.currency }}
                    </option>
                  </template>
                </select>
              </div>
            </div>

            <!-- Rate Info -->
            <div class="p-4 rounded-xl bg-white/5 border border-white/10 mb-6">
              <div class="flex justify-between text-sm mb-2">
                <span class="text-muted">Taux</span>
                <span class="text-base font-medium">
                  1 {{ mode === 'buy' ? toCurrency : fromCurrency }} = 
                  {{ formatMoney(getCurrentRate(), mode === 'buy' ? fromCurrency : toCurrency) }}
                </span>
              </div>
              <div class="flex justify-between text-sm mb-2">
                <span class="text-muted">Frais ({{ feePercentage }}%)</span>
                <span class="text-base">{{ formatMoney(calculatedFee, mode === 'buy' ? fromCurrency : toCurrency) }}</span>
              </div>
              <div class="flex justify-between text-sm font-bold pt-2 border-t border-white/10">
                <span class="text-muted">Vous recevez</span>
                <span class="text-success">{{ toAmount }} {{ toCurrency }}</span>
              </div>
            </div>

            <!-- Execute Button -->
            <button 
              @click="executeExchange"
              :disabled="loading || !fromAmount || fromAmount <= 0 || fromAmount > getWalletBalance(fromCurrency)"
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
                {{ mode === 'buy' ? `Acheter ${toCurrency}` : `Vendre ${fromCurrency}` }}
              </span>
            </button>

            <p v-if="fromAmount > getWalletBalance(fromCurrency)" class="text-error text-sm mt-2 text-center">
              Solde insuffisant
            </p>
          </div>
        </div>

        <!-- Market Prices -->
        <div class="glass-card p-6">
          <h2 class="text-lg font-semibold text-base mb-6">📊 Prix du marché</h2>
          <div class="space-y-3 max-h-[500px] overflow-y-auto custom-scrollbar">
            <div v-for="crypto in markets" :key="crypto.symbol" 
                class="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors cursor-pointer"
                @click="selectCrypto(crypto.symbol)">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold text-sm">{{ crypto.icon || crypto.symbol?.slice(0, 2) }}</span>
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
import { exchangeAPI, walletAPI } from '~/composables/useApi'

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

// Wallets and balances
const wallets = ref([])
const fiatWallets = computed(() => wallets.value.filter(w => w.wallet_type === 'fiat'))
const cryptoWallets = computed(() => wallets.value.filter(w => w.wallet_type === 'crypto'))

// Notification
const notification = ref({
  show: false,
  type: 'success',
  message: ''
})

const showNotification = (type, message) => {
  notification.value = { show: true, type, message }
}

// Crypto currencies with prices
const cryptoCurrencies = ref([
  { symbol: 'BTC', name: 'Bitcoin', icon: '₿', price: 90000, change: 2.5 },
  { symbol: 'ETH', name: 'Ethereum', icon: 'Ξ', price: 3100, change: 1.8 },
  { symbol: 'SOL', name: 'Solana', icon: '◎', price: 180, change: 4.2 },
  { symbol: 'XRP', name: 'Ripple', icon: '✕', price: 2.2, change: 3.2 },
  { symbol: 'BNB', name: 'BNB', icon: '◆', price: 520, change: 1.1 },
  { symbol: 'ADA', name: 'Cardano', icon: '₳', price: 0.95, change: -0.8 },
  { symbol: 'DOGE', name: 'Dogecoin', icon: 'Ð', price: 0.35, change: 5.5 },
  { symbol: 'DOT', name: 'Polkadot', icon: '●', price: 7.5, change: -1.2 },
  { symbol: 'USDT', name: 'Tether', icon: '₮', price: 1, change: 0 },
  { symbol: 'USDC', name: 'USD Coin', icon: '$', price: 1, change: 0 },
])

const markets = computed(() => cryptoCurrencies.value.map(c => ({
  ...c,
  bgColor: getBgColor(c.symbol)
})))

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
    AUD: '🇦🇺', JPY: '🇯🇵', CHF: '🇨🇭', CNY: '🇨🇳', INR: '🇮🇳'
  }
  return emojis[currency] || '💵'
}

const getWalletBalance = (currency) => {
  const wallet = wallets.value.find(w => w.currency === currency)
  return wallet?.balance || 0
}

const getCurrentRate = () => {
  const crypto = mode.value === 'buy' ? toCurrency.value : fromCurrency.value
  const cryptoData = cryptoCurrencies.value.find(c => c.symbol === crypto)
  return cryptoData?.price || 1
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

const calculateConversion = () => {
  const rate = getCurrentRate()
  const fee = fromAmount.value * (feePercentage.value / 100)
  calculatedFee.value = fee
  
  if (mode.value === 'buy') {
    // Buying crypto with fiat: amount / rate
    const netAmount = fromAmount.value - fee
    toAmount.value = (netAmount / rate).toFixed(8)
  } else {
    // Selling crypto for fiat: amount * rate
    const netAmount = fromAmount.value - fee
    toAmount.value = (netAmount * rate).toFixed(2)
  }
}

const swapMode = () => {
  mode.value = mode.value === 'buy' ? 'sell' : 'buy'
  // Swap currencies
  const tempFrom = fromCurrency.value
  fromCurrency.value = toCurrency.value
  toCurrency.value = tempFrom
  calculateConversion()
}

const onFromCurrencyChange = () => {
  calculateConversion()
}

const selectCrypto = (symbol) => {
  if (mode.value === 'buy') {
    toCurrency.value = symbol
  } else {
    fromCurrency.value = symbol
  }
  calculateConversion()
}

const executeExchange = async () => {
  loading.value = true
  try {
    // Get the source and destination wallets
    const sourceWallet = wallets.value.find(w => w.currency === fromCurrency.value)
    let destWallet = wallets.value.find(w => w.currency === toCurrency.value)
    
    if (!sourceWallet) {
      throw new Error(`Portefeuille ${fromCurrency.value} non trouvé`)
    }
    
    // Create destination wallet if it doesn't exist
    if (!destWallet) {
      const crypto = cryptoCurrencies.value.find(c => c.symbol === toCurrency.value)
      const walletType = crypto ? 'crypto' : 'fiat'
      const walletName = crypto ? `${crypto.name} Wallet` : `${toCurrency.value} Wallet`
      
      const { data: newWalletRes } = await walletAPI.createWallet({
        currency: toCurrency.value,
        name: walletName,
        wallet_type: walletType
      })
      destWallet = newWalletRes.wallet || newWalletRes
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
    const { data: exchangeRes } = await exchangeAPI.executeExchange(
      quote.id,
      sourceWallet.id,
      destWallet.id
    )
    
    const exchange = exchangeRes.exchange || exchangeRes
    
    // Success!
    showNotification('success', 
      mode.value === 'buy' 
        ? `Vous avez acheté ${toAmount.value} ${toCurrency.value} avec succès!`
        : `Vous avez vendu ${fromAmount.value} ${fromCurrency.value} avec succès!`
    )
    
    // Refresh wallets
    await fetchWallets()
    
    // Reset form
    fromAmount.value = 0
    toAmount.value = '0'
    
  } catch (error) {
    console.error('Exchange error:', error)
    showNotification('error', error.response?.data?.error || error.message || 'Erreur lors de l\'échange')
  } finally {
    loading.value = false
  }
}

const fetchWallets = async () => {
  try {
    const { data } = await walletAPI.getAll()
    wallets.value = data.wallets || []
    
    // Set default from currency based on first fiat wallet
    if (fiatWallets.value.length > 0 && mode.value === 'buy') {
      fromCurrency.value = fiatWallets.value[0].currency
    }
  } catch (error) {
    console.error('Error fetching wallets:', error)
  }
}

const fetchRates = async () => {
  try {
    const { data } = await exchangeAPI.getRates()
    if (data.rates) {
      // Update crypto prices from API
      Object.entries(data.rates).forEach(([pair, rate]) => {
        const crypto = cryptoCurrencies.value.find(c => pair.includes(c.symbol))
        if (crypto && rate.Rate) {
          crypto.price = rate.Rate
        }
      })
    }
  } catch (error) {
    console.error('Error fetching rates:', error)
  }
}

// Watch mode changes
watch(mode, () => {
  if (mode.value === 'buy') {
    fromCurrency.value = fiatWallets.value[0]?.currency || 'USD'
    toCurrency.value = 'BTC'
  } else {
    fromCurrency.value = 'BTC'
    toCurrency.value = fiatWallets.value[0]?.currency || 'USD'
  }
  calculateConversion()
})

onMounted(async () => {
  await fetchWallets()
  await fetchRates()
  calculateConversion()
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
