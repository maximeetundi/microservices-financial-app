<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">Acheter Crypto ₿</h1>
          <p class="text-slate-400">Achetez et vendez des cryptomonnaies</p>
        </div>
        <NuxtLink to="/exchange" class="text-indigo-400 hover:text-indigo-300 flex items-center gap-2">
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
            <h2 class="text-lg font-semibold text-white mb-6">Échanger</h2>
            
            <!-- From Currency -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-slate-400 mb-2">Vous payez</label>
              <div class="flex gap-4">
                <div class="flex-1">
                  <input 
                    v-model.number="fromAmount" 
                    type="number" 
                    placeholder="0.00"
                    class="input-premium w-full text-2xl font-bold"
                    @input="calculateToAmount"
                  />
                </div>
                <select v-model="fromCurrency" @change="calculateToAmount" class="input-premium w-32">
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                  <option value="GBP">GBP</option>
                </select>
              </div>
              <p class="text-xs text-slate-500 mt-2">Solde disponible: {{ formatMoney(balance[fromCurrency] || 0, fromCurrency) }}</p>
            </div>

            <!-- Swap Button -->
            <div class="flex justify-center my-4">
              <button @click="swapCurrencies" class="p-3 rounded-xl bg-indigo-500/20 hover:bg-indigo-500/30 transition-colors">
                <svg class="w-6 h-6 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
                </svg>
              </button>
            </div>

            <!-- To Currency -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-slate-400 mb-2">Vous recevez</label>
              <div class="flex gap-4">
                <div class="flex-1">
                  <input 
                    v-model.number="toAmount" 
                    type="number" 
                    placeholder="0.00"
                    class="input-premium w-full text-2xl font-bold"
                    readonly
                  />
                </div>
                <select v-model="toCurrency" @change="calculateToAmount" class="input-premium w-32">
                  <option value="BTC">BTC</option>
                  <option value="ETH">ETH</option>
                  <option value="USDT">USDT</option>
                  <option value="SOL">SOL</option>
                  <option value="XRP">XRP</option>
                </select>
              </div>
            </div>

            <!-- Rate Info -->
            <div class="p-4 rounded-xl bg-slate-800/50 border border-slate-700/50 mb-6">
              <div class="flex justify-between text-sm mb-2">
                <span class="text-slate-400">Taux</span>
                <span class="text-white">1 {{ toCurrency }} = {{ formatMoney(rates[toCurrency] || 0, fromCurrency) }}</span>
              </div>
              <div class="flex justify-between text-sm mb-2">
                <span class="text-slate-400">Frais (0.5%)</span>
                <span class="text-white">{{ formatMoney(fee, fromCurrency) }}</span>
              </div>
              <div class="flex justify-between text-sm font-bold pt-2 border-t border-slate-700">
                <span class="text-slate-300">Total</span>
                <span class="text-white">{{ formatMoney(fromAmount + fee, fromCurrency) }}</span>
              </div>
            </div>

            <!-- Execute Button -->
            <button 
              @click="executeExchange"
              :disabled="loading || !fromAmount || fromAmount <= 0"
              class="btn-premium w-full py-4 text-lg font-bold disabled:opacity-50"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                Traitement...
              </span>
              <span v-else>Acheter {{ toCurrency }}</span>
            </button>
          </div>
        </div>

        <!-- Market Prices -->
        <div class="glass-card p-6">
          <h2 class="text-lg font-semibold text-white mb-6">Prix du marché</h2>
          <div class="space-y-4">
            <div v-for="crypto in markets" :key="crypto.symbol" 
                class="flex items-center justify-between p-3 rounded-xl bg-slate-800/50 border border-slate-700/50 hover:bg-slate-700/50 transition-colors cursor-pointer"
                @click="selectCrypto(crypto.symbol)">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center" :class="crypto.bgColor">
                  <span class="text-white font-bold text-sm">{{ crypto.symbol?.slice(0, 2) || '??' }}</span>
                </div>
                <div>
                  <p class="font-medium text-white">{{ crypto.name }}</p>
                  <p class="text-xs text-slate-400">{{ crypto.symbol }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-white">${{ formatPrice(crypto.price) }}</p>
                <p class="text-xs" :class="(crypto.change || 0) >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ (crypto.change || 0) >= 0 ? '+' : '' }}{{ (crypto.change || 0).toFixed(2) }}%
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Success/Error Modal -->
      <div v-if="showResult" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
        <div class="glass-card p-8 max-w-md w-full mx-4 text-center">
          <div v-if="success" class="w-16 h-16 rounded-full bg-emerald-500/20 flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          <div v-else class="w-16 h-16 rounded-full bg-red-500/20 flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
          <h3 class="text-xl font-bold text-white mb-2">{{ success ? 'Échange réussi!' : 'Erreur' }}</h3>
          <p class="text-slate-400 mb-6">{{ resultMessage }}</p>
          <button @click="showResult = false" class="btn-premium px-8 py-3">Fermer</button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { exchangeAPI, walletAPI } from '~/composables/useApi'

const fromCurrency = ref('USD')
const toCurrency = ref('BTC')
const fromAmount = ref(100)
const toAmount = ref(0)
const loading = ref(false)
const showResult = ref(false)
const success = ref(false)
const resultMessage = ref('')
const balance = ref({ USD: 0, EUR: 0, GBP: 0 })

const rates = ref({
  BTC: 45000,
  ETH: 3000,
  USDT: 1,
  SOL: 100,
  XRP: 0.60
})

const markets = ref([
  { symbol: 'BTC', name: 'Bitcoin', price: 45000, change: 2.5, bgColor: 'bg-orange-500' },
  { symbol: 'ETH', name: 'Ethereum', price: 3000, change: 1.8, bgColor: 'bg-blue-500' },
  { symbol: 'SOL', name: 'Solana', price: 100, change: -0.5, bgColor: 'bg-purple-500' },
  { symbol: 'XRP', name: 'Ripple', price: 0.60, change: 3.2, bgColor: 'bg-slate-600' }
])

const fee = computed(() => fromAmount.value * 0.005)

const formatMoney = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount)
}

const formatPrice = (price) => {
  if (price == null) return '0'
  return price >= 1 ? (price || 0).toLocaleString() : (price || 0).toFixed(4)
}

const calculateToAmount = () => {
  const rate = rates.value[toCurrency.value] || 1
  toAmount.value = (fromAmount.value - fee.value) / rate
}

const swapCurrencies = () => {
  const temp = fromCurrency.value
  fromCurrency.value = toCurrency.value
  toCurrency.value = temp
  calculateToAmount()
}

const selectCrypto = (symbol) => {
  toCurrency.value = symbol
  calculateToAmount()
}

const executeExchange = async () => {
  loading.value = true
  try {
    await exchangeAPI.execute(fromCurrency.value, toCurrency.value, fromAmount.value)
    success.value = true
    resultMessage.value = `Vous avez acheté ${toAmount.value.toFixed(8)} ${toCurrency.value}`
    showResult.value = true
    fromAmount.value = 0
    toAmount.value = 0
  } catch (e) {
    success.value = false
    resultMessage.value = e.response?.data?.error || 'Erreur lors de l\'échange'
    showResult.value = true
  } finally {
    loading.value = false
  }
}

const fetchData = async () => {
  try {
    const [ratesRes, walletsRes, marketsRes] = await Promise.all([
      exchangeAPI.getRates().catch(() => ({ data: {} })),
      walletAPI.getAll().catch(() => ({ data: { wallets: [] } })),
      exchangeAPI.getMarkets().catch(() => ({ data: { markets: [] } }))
    ])
    
    if (ratesRes.data?.rates) {
      Object.entries(ratesRes.data.rates).forEach(([key, val]) => {
        if (key.includes('BTC')) rates.value.BTC = val.Rate
        if (key.includes('ETH')) rates.value.ETH = val.Rate
      })
    }
    
    if (walletsRes.data?.wallets) {
      walletsRes.data.wallets.forEach(w => {
        balance.value[w.currency] = w.balance
      })
    }

    if (marketsRes.data?.markets) {
      markets.value = marketsRes.data.markets.slice(0, 4).map(m => ({
        symbol: m.Symbol,
        name: m.BaseAsset,
        price: m.Price,
        change: m.Change24h,
        bgColor: 'bg-indigo-500'
      }))
    }

    calculateToAmount()
  } catch (e) {
    console.error('Error fetching data:', e)
  }
}

onMounted(fetchData)

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>
