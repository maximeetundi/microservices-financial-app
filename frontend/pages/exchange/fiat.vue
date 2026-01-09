<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="text-center mb-10">
        <h1 class="text-3xl font-bold text-base mb-3">
          💱 Conversion Monnaies Fiduciaires
        </h1>
        <p class="text-muted max-w-2xl mx-auto">
          Convertissez instantanément entre USD, EUR, GBP, JPY, CAD, AUD et plus avec nos meilleurs taux.
        </p>
      </div>

      <!-- Conversion Form -->
      <div class="glass-card mb-8">
        <div class="grid grid-cols-1 md:grid-cols-[1fr_auto_1fr] gap-6 items-start">
          <!-- From Currency -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-muted">Depuis</label>
            <div class="p-4 rounded-2xl bg-surface-hover border border-secondary-200 dark:border-secondary-700 transition-colors focus-within:ring-2 focus-within:ring-primary-500">
              <div class="flex gap-4 mb-4">
                 <select 
                    v-model="fromCurrency"
                    class="bg-transparent text-base font-bold text-lg outline-none cursor-pointer w-full"
                    @change="updateRates"
                  >
                  <option v-for="currency in supportedCurrencies" :key="currency.code" :value="currency.code" class="text-gray-900 dark:text-gray-900">
                    {{ currency.flag }} {{ currency.code }}
                  </option>
                </select>
              </div>
              
              <input
                v-model="fromAmount"
                type="number"
                class="w-full bg-transparent text-3xl font-bold text-base outline-none placeholder-muted/50"
                placeholder="0.00"
                @input="calculateConversion"
                min="1"
                step="0.01"
              />
              
               <div class="mt-2 text-sm text-muted" v-if="sourceWallets.length > 0">
                 Solde: {{ formatMoney(currentSourceBalance) }}
               </div>
            </div>
            
             <select 
                v-if="sourceWallets.length > 0"
                v-model="fromWalletId"
                class="input-field text-sm"
              >
                <option v-for="wallet in sourceWallets" :key="wallet.id" :value="wallet.id">
                  {{ wallet.name }} ({{ wallet.balance }} {{ wallet.currency }})
                </option>
              </select>
          </div>

          <!-- Swap Button -->
          <div class="flex items-center justify-center md:pt-10">
            <button
              @click="swapCurrencies"
              class="p-3 rounded-full bg-surface-hover hover:bg-primary text-muted hover:text-white border border-secondary-200 dark:border-secondary-700 transition-all shadow-lg active:scale-95"
              :disabled="loading"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
            </button>
          </div>

          <!-- To Currency -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-muted">Vers</label>
             <div class="p-4 rounded-2xl bg-surface-hover border border-secondary-200 dark:border-secondary-700 transition-colors">
              <div class="flex gap-4 mb-4">
                 <select 
                    v-model="toCurrency"
                    class="bg-transparent text-base font-bold text-lg outline-none cursor-pointer w-full"
                    @change="updateRates"
                  >
                  <option v-for="currency in supportedCurrencies" :key="currency.code" :value="currency.code" class="text-gray-900 dark:text-white">
                     {{ currency.flag }} {{ currency.code }}
                  </option>
                </select>
              </div>
              
              <div class="relative">
                 <input
                  :value="toAmount"
                  type="text"
                  readonly
                  class="w-full bg-transparent text-3xl font-bold text-base outline-none"
                  placeholder="0.00"
                />
                 <div v-if="loading" class="absolute right-0 top-1/2 transform -translate-y-1/2">
                  <div class="loading-spinner w-6 h-6"></div>
                </div>
              </div>
              
               <div class="mt-2 text-sm text-muted">
                 1 {{ fromCurrency }} = {{ (Number(exchangeRate?.rate) || 0).toFixed(4) }} {{ toCurrency }}
               </div>
            </div>

             <select 
                v-if="destWallets.length > 0"
                v-model="toWalletId"
                class="input-field text-sm"
              >
                <option v-for="wallet in destWallets" :key="wallet.id" :value="wallet.id">
                  {{ wallet.name }} ({{ wallet.balance }} {{ wallet.currency }})
                </option>
                <option value="__new__">➕ Créer un nouveau portefeuille {{ toCurrency }}</option>
              </select>
              
              <!-- No wallet exists for destination currency -->
              <div v-else class="p-4 rounded-xl bg-warning/10 border border-warning/30 text-warning">
                <div class="flex items-center gap-3 mb-3">
                  <span class="text-xl">⚠️</span>
                  <span class="font-medium">Aucun portefeuille {{ toCurrency }}</span>
                </div>
                <p class="text-sm text-muted mb-3">Un nouveau portefeuille {{ toCurrency }} sera créé automatiquement lors de la conversion.</p>
                <button
                  @click="createDestWallet"
                  class="w-full py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary-600 transition-colors text-sm font-medium"
                  :disabled="creatingWallet"
                >
                  <span v-if="creatingWallet">Création en cours...</span>
                  <span v-else>➕ Créer un portefeuille {{ toCurrency }} maintenant</span>
                </button>
              </div>
          </div>
        </div>

        <!-- Exchange Rate Details -->
        <div v-if="exchangeRate" class="mt-8 p-4 rounded-xl bg-surface border border-secondary-100 dark:border-secondary-800">
           <div class="flex justify-between items-center text-sm">
             <span class="text-muted">Frais de service ({{ exchangeRate.fee_percentage }}%)</span>
             <span class="text-base font-medium">{{ calculateFee }} {{ fromCurrency }}</span>
           </div>
             <div class="flex justify-between items-center text-sm mt-2">
             <span class="text-muted">Variation 24h</span>
             <span :class="exchangeRate.change_24h >= 0 ? 'text-success' : 'text-error'" class="font-bold">
                {{ (exchangeRate.change_24h || 0) >= 0 ? '+' : '' }}{{ ((exchangeRate.change_24h || 0) * 100).toFixed(2) }}%
              </span>
           </div>
        </div>

        <!-- Action Button -->
        <div class="mt-8">
          <button
            @click="executeFiatConversion"
            :disabled="!canConvert || loading"
            class="w-full btn-premium py-4 rounded-xl text-lg shadow-lg disabled:opacity-50 disabled:cursor-not-allowed group relative overflow-hidden"
          >
            <span class="relative z-10 flex items-center justify-center gap-2">
               <span v-if="loading" class="loading-spinner w-5 h-5 border-white/30 border-t-white"></span>
               {{ loading ? 'Conversion en cours...' : 'Confirmer la conversion' }}
            </span>
          </button>
        </div>
      </div>

      <!-- Stats / Comparisons -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <!-- Fee Comparison -->
         <div class="glass-card p-6">
           <h3 class="text-lg font-bold text-base mb-4">💰 Économies Estimées</h3>
           <div class="space-y-4">
             <div class="flex items-center justify-between p-3 rounded-lg bg-success/10 border border-success/20">
               <div class="flex items-center gap-3">
                 <div class="w-8 h-8 rounded-full bg-success flex items-center justify-center text-white text-xs">Cb</div>
                 <span class="font-bold text-base">Zekora</span>
               </div>
               <span class="font-bold text-success">{{ ourFee }}%</span>
             </div>
              <div class="flex items-center justify-between p-3 rounded-lg opacity-60">
               <div class="flex items-center gap-3">
                 <div class="w-8 h-8 rounded-full bg-secondary-200 dark:bg-secondary-700 flex items-center justify-center text-xs">🏛️</div>
                 <span class="font-medium text-base">Banque Classique</span>
               </div>
               <span class="font-medium text-base">3.5% + 15€</span>
             </div>
           </div>
           
            <div class="mt-4 pt-4 border-t border-secondary-100 dark:border-secondary-800">
              <p class="text-sm text-muted">
                Vous économisez environ <span class="text-success font-bold">{{ savingsAmount }}€</span> sur cette transaction.
              </p>
            </div>
         </div>

        <!-- Recent Conversions -->
        <div class="glass-card p-6">
          <h3 class="text-lg font-bold text-base mb-4">📊 Récemment</h3>
           <div v-if="recentConversions.length === 0" class="text-center py-8">
             <p class="text-muted text-sm">Aucune conversion récente</p>
           </div>
           <div v-else class="space-y-3 max-h-[200px] overflow-y-auto custom-scrollbar">
             <div 
              v-for="conversion in recentConversions" 
              :key="conversion.id"
              class="flex items-center justify-between p-3 rounded-lg bg-surface-hover/50 hover:bg-surface-hover border border-transparent hover:border-secondary-200 dark:hover:border-secondary-700 transition-colors"
            >
              <div>
                 <p class="text-sm font-bold text-base">
                   {{ conversion.from_amount }} {{ conversion.from_currency }} → {{ conversion.to_amount }} {{ conversion.to_currency }}
                 </p>
                 <p class="text-xs text-muted">{{ formatTime(conversion.created_at) }}</p>
              </div>
              <span class="text-xs px-2 py-1 rounded-full bg-success/10 text-success border border-success/20">
                Succès
              </span>
            </div>
           </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { exchangeAPI, walletAPI } from '~/composables/useApi'

// Reactive data
const loading = ref(false)
const fromCurrency = ref('USD')
const toCurrency = ref('EUR')
const fromAmount = ref(1000)
const toAmount = ref(0)
const exchangeRate = ref(null)
const recentConversions = ref([])

const wallets = ref([])
const loadingWallets = ref(false)
const creatingWallet = ref(false)
const fromWalletId = ref('')
const toWalletId = ref('')

// Supported currencies - ALL WORLD CURRENCIES
const supportedCurrencies = ref([
  // Major currencies
  { code: 'USD', name: 'US Dollar', flag: '🇺🇸' },
  { code: 'EUR', name: 'Euro', flag: '🇪🇺' },
  { code: 'GBP', name: 'British Pound', flag: '🇬🇧' },
  { code: 'JPY', name: 'Japanese Yen', flag: '🇯🇵' },
  { code: 'CHF', name: 'Swiss Franc', flag: '🇨🇭' },
  // Americas
  { code: 'CAD', name: 'Canadian Dollar', flag: '🇨🇦' },
  { code: 'MXN', name: 'Mexican Peso', flag: '🇲🇽' },
  { code: 'BRL', name: 'Brazilian Real', flag: '🇧🇷' },
  { code: 'ARS', name: 'Argentine Peso', flag: '🇦🇷' },
  { code: 'CLP', name: 'Chilean Peso', flag: '🇨🇱' },
  { code: 'COP', name: 'Colombian Peso', flag: '🇨🇴' },
  { code: 'PEN', name: 'Peruvian Sol', flag: '🇵🇪' },
  // Europe
  { code: 'NOK', name: 'Norwegian Krone', flag: '🇳🇴' },
  { code: 'SEK', name: 'Swedish Krona', flag: '🇸🇪' },
  { code: 'DKK', name: 'Danish Krone', flag: '🇩🇰' },
  { code: 'PLN', name: 'Polish Złoty', flag: '🇵🇱' },
  { code: 'CZK', name: 'Czech Koruna', flag: '🇨🇿' },
  { code: 'HUF', name: 'Hungarian Forint', flag: '🇭🇺' },
  { code: 'RON', name: 'Romanian Leu', flag: '🇷🇴' },
  { code: 'TRY', name: 'Turkish Lira', flag: '🇹🇷' },
  { code: 'RUB', name: 'Russian Ruble', flag: '🇷🇺' },
  // Asia
  { code: 'CNY', name: 'Chinese Yuan', flag: '🇨🇳' },
  { code: 'HKD', name: 'Hong Kong Dollar', flag: '🇭🇰' },
  { code: 'SGD', name: 'Singapore Dollar', flag: '🇸🇬' },
  { code: 'KRW', name: 'South Korean Won', flag: '🇰🇷' },
  { code: 'INR', name: 'Indian Rupee', flag: '🇮🇳' },
  { code: 'IDR', name: 'Indonesian Rupiah', flag: '🇮🇩' },
  { code: 'MYR', name: 'Malaysian Ringgit', flag: '🇲🇾' },
  { code: 'THB', name: 'Thai Baht', flag: '🇹🇭' },
  { code: 'PHP', name: 'Philippine Peso', flag: '🇵🇭' },
  { code: 'VND', name: 'Vietnamese Đồng', flag: '🇻🇳' },
  { code: 'PKR', name: 'Pakistani Rupee', flag: '🇵🇰' },
  { code: 'BDT', name: 'Bangladeshi Taka', flag: '🇧🇩' },
  // Middle East
  { code: 'AED', name: 'UAE Dirham', flag: '🇦🇪' },
  { code: 'SAR', name: 'Saudi Riyal', flag: '🇸🇦' },
  { code: 'QAR', name: 'Qatari Riyal', flag: '🇶🇦' },
  { code: 'KWD', name: 'Kuwaiti Dinar', flag: '🇰🇼' },
  { code: 'BHD', name: 'Bahraini Dinar', flag: '🇧🇭' },
  { code: 'OMR', name: 'Omani Rial', flag: '🇴🇲' },
  { code: 'ILS', name: 'Israeli Shekel', flag: '🇮🇱' },
  { code: 'EGP', name: 'Egyptian Pound', flag: '🇪🇬' },
  // Africa
  { code: 'XAF', name: 'Franc CFA (CEMAC)', flag: '🇨🇲' },
  { code: 'XOF', name: 'Franc CFA (UEMOA)', flag: '🇸🇳' },
  { code: 'NGN', name: 'Nigerian Naira', flag: '🇳🇬' },
  { code: 'ZAR', name: 'South African Rand', flag: '🇿🇦' },
  { code: 'KES', name: 'Kenyan Shilling', flag: '🇰🇪' },
  { code: 'GHS', name: 'Ghanaian Cedi', flag: '🇬🇭' },
  { code: 'MAD', name: 'Moroccan Dirham', flag: '🇲🇦' },
  { code: 'TND', name: 'Tunisian Dinar', flag: '🇹🇳' },
  { code: 'DZD', name: 'Algerian Dinar', flag: '🇩🇿' },
  { code: 'UGX', name: 'Ugandan Shilling', flag: '🇺🇬' },
  { code: 'TZS', name: 'Tanzanian Shilling', flag: '🇹🇿' },
  { code: 'RWF', name: 'Rwandan Franc', flag: '🇷🇼' },
  { code: 'ETB', name: 'Ethiopian Birr', flag: '🇪🇹' },
  // Oceania
  { code: 'AUD', name: 'Australian Dollar', flag: '🇦🇺' },
  { code: 'NZD', name: 'New Zealand Dollar', flag: '🇳🇿' },
  { code: 'FJD', name: 'Fijian Dollar', flag: '🇫🇯' },
])

// Computed properties
const canConvert = computed(() => {
  // Allow conversion if we have source wallet and either dest wallet exists OR we'll create one
  const hasDestWallet = toWalletId.value && toWalletId.value !== '__new__'
  const willCreateWallet = destWallets.value.length === 0 || toWalletId.value === '__new__'
  return fromAmount.value > 0 && fromCurrency.value !== toCurrency.value && fromWalletId.value && (hasDestWallet || willCreateWallet)
})

const sourceWallets = computed(() => wallets.value.filter(w => w.currency === fromCurrency.value))
const destWallets = computed(() => wallets.value.filter(w => w.currency === toCurrency.value))

const currentSourceBalance = computed(() => {
    const w = wallets.value.find(w => w.id === fromWalletId.value)
    return w ? w.balance : 0
})

const ourFee = computed(() => {
  if (!exchangeRate.value) return '0.25'
  return Number(exchangeRate.value.fee_percentage || 0.25).toFixed(2)
})

const calculateFee = computed(() => {
    return (fromAmount.value * (parseFloat(ourFee.value) / 100)).toFixed(2)
})

const savingsAmount = computed(() => {
  const traditionalFee = (fromAmount.value * 0.035) + 15
  const ourFeeAmount = fromAmount.value * (parseFloat(ourFee.value) / 100)
  return Math.max(0, traditionalFee - ourFeeAmount).toFixed(2)
})

// Methods
const fetchWallets = async () => {
  loadingWallets.value = true
  try {
    const res = await walletAPI.getAll() // corrected to match wallet page usage
    wallets.value = res.data.wallets || []
    updateWalletSelection()
  } catch (e) {
    console.error('Failed to fetch wallets', e)
     // Mock for dev
    wallets.value = [
         { id: 1, currency: 'USD', balance: 5000, name: 'Main USD' },
         { id: 2, currency: 'EUR', balance: 1200, name: 'Euro Savings' }
    ]
    updateWalletSelection()
  } finally {
    loadingWallets.value = false
  }
}

const updateWalletSelection = () => {
  const source = sourceWallets.value[0]
  if (source) fromWalletId.value = source.id
  else fromWalletId.value = ''

  const dest = destWallets.value[0]
  if (dest) toWalletId.value = dest.id
  else toWalletId.value = ''
}

const updateRates = async () => {
  if (fromCurrency.value === toCurrency.value) return
  
  // Update wallet selection if currency changed
  updateWalletSelection()

  loading.value = true
  try {
    const { data } = await exchangeAPI.getRate(fromCurrency.value, toCurrency.value)
    exchangeRate.value = data
    calculateConversion()
  } catch (error) {
    console.error('Error fetching rates:', error)
    // mock rate
    exchangeRate.value = { rate: 0.92, fee_percentage: 0.25, change_24h: 0.05 }
    calculateConversion()
  } finally {
    loading.value = false
  }
}

const calculateConversion = () => {
  if (exchangeRate.value && fromAmount.value > 0) {
    const rateVal = exchangeRate.value.Rate || exchangeRate.value.rate || 1
    const pFee = exchangeRate.value.FeePercentage || exchangeRate.value.fee_percentage || 0.5
    
    // Recalculate based on normalized values
    const rawConverted = fromAmount.value * rateVal
    toAmount.value = (rawConverted * (1 - pFee/100)).toFixed(2)
  } else {
    toAmount.value = 0
  }
}

const swapCurrencies = () => {
  const temp = fromCurrency.value
  fromCurrency.value = toCurrency.value
  toCurrency.value = temp
  updateWalletSelection()
  updateRates()
}

const createDestWallet = async () => {
  creatingWallet.value = true
  try {
    const currencyInfo = supportedCurrencies.value.find(c => c.code === toCurrency.value)
    const walletName = currencyInfo ? `${currencyInfo.name} Wallet` : `${toCurrency.value} Wallet`
    
    const { data } = await walletAPI.createWallet({
      currency: toCurrency.value,
      name: walletName,
      wallet_type: 'fiat'
    })
    
    // Refresh wallets and select the new one
    await fetchWallets()
    
    if (data?.id) {
      toWalletId.value = data.id
    }
    
    alert(`Portefeuille ${toCurrency.value} créé avec succès !`)
  } catch (error) {
    console.error('Failed to create wallet:', error)
    alert('Erreur lors de la création du portefeuille')
  } finally {
    creatingWallet.value = false
  }
}

const executeFiatConversion = async () => {
  loading.value = true
  try {
    let destWalletIdToUse = toWalletId.value
    
    // Create destination wallet if needed
    if (!destWalletIdToUse || destWalletIdToUse === '__new__' || destWallets.value.length === 0) {
      const currencyInfo = supportedCurrencies.value.find(c => c.code === toCurrency.value)
      const walletName = currencyInfo ? `${currencyInfo.name} Wallet` : `${toCurrency.value} Wallet`
      
      const { data: newWallet } = await walletAPI.createWallet({
        currency: toCurrency.value,
        name: walletName,
        wallet_type: 'fiat'
      })
      
      if (!newWallet?.id) throw new Error('Failed to create destination wallet')
      destWalletIdToUse = newWallet.id
    }
    
    // 1. Get Quote
    const { data: quote } = await exchangeAPI.getQuote(fromCurrency.value, toCurrency.value, fromAmount.value)
    
    if (!quote || !quote.ID) throw new Error('Failed to get quote')

    // 2. Execute Exchange
    const { data: exchange } = await exchangeAPI.executeExchange(
        quote.ID,
        fromWalletId.value,
        destWalletIdToUse
    )

    // Refresh wallets
    fetchWallets()

    // Add to recent conversions
    recentConversions.value.unshift({
      id: exchange.ID || Date.now(),
      from_amount: exchange.FromAmount || fromAmount.value,
      from_currency: exchange.FromCurrency || fromCurrency.value,
      to_amount: exchange.ToAmount || toAmount.value,
      to_currency: exchange.ToCurrency || toCurrency.value,
      created_at: new Date().toISOString()
    })

    alert('Conversion réussie !')
    
  } catch (error) {
    console.error('Conversion error:', error)
    alert('Erreur lors de la conversion: ' + (error.message || 'Unknown error'))
  } finally {
    loading.value = false
  }
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit'
  })
}

const formatMoney = (amount) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: fromCurrency.value }).format(amount)
}

// Lifecycle
onMounted(() => {
  fetchWallets()
  updateRates()
})

// SEO
definePageMeta({
  title: 'Conversion Fiat - Zekora',
  description: 'Convertissez instantanément entre différentes monnaies fiduciaires avec des frais réduits'
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
</style>