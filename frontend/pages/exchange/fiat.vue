<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          💱 Conversion Monnaies Fiduciaires
        </h1>
        <p class="text-gray-600">
          Convertissez instantanément entre USD, EUR, GBP, JPY, CAD, AUD et plus
        </p>
      </div>

      <!-- Conversion Form -->
      <div class="bg-white rounded-lg shadow-lg p-6 mb-8">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- From Currency -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Depuis
            </label>
            <div class="relative">
              <select 
                v-model="fromCurrency"
                class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                @change="updateRates"
              >
                <option v-for="currency in supportedCurrencies" :key="currency.code" :value="currency.code">
                  {{ currency.flag }} {{ currency.code }} - {{ currency.name }}
                </option>
              </select>
            </div>
            <div class="mt-2">
              <input
                v-model="fromAmount"
                type="number"
                placeholder="Montant"
                class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                @input="calculateConversion"
                min="1"
                step="0.01"
              />
            </div>
          </div>

          <!-- Swap Button -->
          <div class="flex items-center justify-center">
            <button
              @click="swapCurrencies"
              class="bg-blue-500 hover:bg-blue-600 text-gray-900 p-3 rounded-full transition-colors duration-200"
              :disabled="loading"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </button>
          </div>

          <!-- To Currency -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Vers
            </label>
            <div class="relative">
              <select 
                v-model="toCurrency"
                class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                @change="updateRates"
              >
                <option v-for="currency in supportedCurrencies" :key="currency.code" :value="currency.code">
                  {{ currency.flag }} {{ currency.code }} - {{ currency.name }}
                </option>
              </select>
            </div>
            <div class="mt-2 relative">
              <input
                :value="toAmount"
                type="text"
                readonly
                placeholder="Montant converti"
                class="w-full p-3 border border-gray-300 rounded-lg bg-gray-50"
              />
              <div v-if="loading" class="absolute right-3 top-1/2 transform -translate-y-1/2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-500"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Exchange Rate Info -->
        <div v-if="exchangeRate" class="mt-6 p-4 bg-gray-50 rounded-lg">
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span class="text-gray-600">Taux de change:</span>
              <div class="font-semibold">{{ exchangeRate.rate?.toFixed(4) }}</div>
            </div>
            <div>
              <span class="text-gray-600">Frais CryptoBank:</span>
              <div class="font-semibold text-green-600">{{ exchangeRate.fee_percentage }}%</div>
            </div>
            <div>
              <span class="text-gray-600">Variation 24h:</span>
              <div :class="exchangeRate.change_24h >= 0 ? 'text-green-600' : 'text-red-600'" class="font-semibold">
                {{ exchangeRate.change_24h >= 0 ? '+' : '' }}{{ (exchangeRate.change_24h * 100).toFixed(2) }}%
              </div>
            </div>
            <div>
              <span class="text-gray-600">Dernière MAJ:</span>
              <div class="font-semibold">{{ formatTime(exchangeRate.last_updated) }}</div>
            </div>
          </div>
        </div>

        <!-- Conversion Button -->
        <div class="mt-6">
          <button
            @click="executeFiatConversion"
            :disabled="!canConvert || loading"
            class="w-full bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 disabled:from-gray-400 disabled:to-gray-400 text-gray-900 font-semibold py-4 px-6 rounded-lg transition-all duration-200 transform hover:scale-105 disabled:scale-100"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
              Conversion en cours...
            </span>
            <span v-else>
              💱 Convertir {{ fromAmount }} {{ fromCurrency }} → {{ toCurrency }}
            </span>
          </button>
        </div>
      </div>

      <!-- Fee Comparison -->
      <div class="bg-white rounded-lg shadow-lg p-6 mb-8">
        <h3 class="text-xl font-bold mb-4">💰 Comparaison des Frais</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="p-4 border-2 border-green-200 rounded-lg bg-green-50">
            <h4 class="font-bold text-green-800">🏦 CryptoBank</h4>
            <div class="text-2xl font-bold text-green-600">{{ ourFee }}%</div>
            <p class="text-sm text-green-700">Frais transparents</p>
          </div>
          <div class="p-4 border border-gray-200 rounded-lg">
            <h4 class="font-bold text-gray-800">🏪 Banques Traditionnelles</h4>
            <div class="text-2xl font-bold text-red-600">3.5%</div>
            <p class="text-sm text-gray-600">+ frais fixes 15€</p>
          </div>
          <div class="p-4 border border-gray-200 rounded-lg">
            <h4 class="font-bold text-gray-800">💸 Services de Transfert</h4>
            <div class="text-2xl font-bold text-orange-600">2.0%</div>
            <p class="text-sm text-gray-600">+ frais fixes 10€</p>
          </div>
        </div>
        <div class="mt-4 p-3 bg-blue-50 rounded-lg">
          <p class="text-sm text-blue-800">
            💡 <strong>Économisez {{ savingsAmount }}€</strong> par rapport aux banques traditionnelles sur cette conversion !
          </p>
        </div>
      </div>

      <!-- Recent Conversions -->
      <div class="bg-white rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold mb-4">📊 Conversions Récentes</h3>
        <div v-if="recentConversions.length === 0" class="text-center py-8 text-gray-500">
          <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <p>Aucune conversion récente</p>
        </div>
        <div v-else class="space-y-3">
          <div 
            v-for="conversion in recentConversions" 
            :key="conversion.id"
            class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50"
          >
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                <span class="text-blue-600 font-semibold text-xs">💱</span>
              </div>
              <div>
                <div class="font-semibold">
                  {{ conversion.from_amount }} {{ conversion.from_currency }} → {{ conversion.to_amount }} {{ conversion.to_currency }}
                </div>
                <div class="text-sm text-gray-500">
                  Taux: {{ conversion.exchange_rate }} • Frais: {{ conversion.fee }}€
                </div>
              </div>
            </div>
            <div class="text-right">
              <div class="text-sm font-semibold" :class="conversion.status === 'completed' ? 'text-green-600' : 'text-yellow-600'">
                {{ conversion.status === 'completed' ? '✅ Terminée' : '⏳ En cours' }}
              </div>
              <div class="text-xs text-gray-500">{{ formatTime(conversion.created_at) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Reactive data
const loading = ref(false)
const fromCurrency = ref('USD')
const toCurrency = ref('EUR')
const fromAmount = ref(1000)
const toAmount = ref(0)
const exchangeRate = ref(null)
const recentConversions = ref([])

// Supported currencies
const supportedCurrencies = ref([
  { code: 'USD', name: 'US Dollar', flag: '🇺🇸' },
  { code: 'EUR', name: 'Euro', flag: '🇪🇺' },
  { code: 'GBP', name: 'British Pound', flag: '🇬🇧' },
  { code: 'JPY', name: 'Japanese Yen', flag: '🇯🇵' },
  { code: 'CAD', name: 'Canadian Dollar', flag: '🇨🇦' },
  { code: 'AUD', name: 'Australian Dollar', flag: '🇦🇺' },
  { code: 'CHF', name: 'Swiss Franc', flag: '🇨🇭' },
  { code: 'SEK', name: 'Swedish Krona', flag: '🇸🇪' },
  { code: 'NOK', name: 'Norwegian Krone', flag: '🇳🇴' },
])

// Computed properties
const canConvert = computed(() => {
  return fromAmount.value > 0 && fromCurrency.value !== toCurrency.value
})

const ourFee = computed(() => {
  if (!exchangeRate.value) return '0.25'
  return (exchangeRate.value.fee_percentage || 0.25).toFixed(2)
})

const savingsAmount = computed(() => {
  const traditionalFee = (fromAmount.value * 0.035) + 15
  const ourFeeAmount = fromAmount.value * (parseFloat(ourFee.value) / 100)
  return Math.max(0, traditionalFee - ourFeeAmount).toFixed(2)
})

// Methods
const updateRates = async () => {
  if (fromCurrency.value === toCurrency.value) return
  
  loading.value = true
  try {
    const response = await $fetch(`/gateway/exchange-service/fiat/rates/${fromCurrency.value}/${toCurrency.value}`)
    exchangeRate.value = response
    calculateConversion()
  } catch (error) {
    console.error('Error fetching rates:', error)
  } finally {
    loading.value = false
  }
}

const calculateConversion = () => {
  if (exchangeRate.value && fromAmount.value > 0) {
    const converted = fromAmount.value * exchangeRate.value.rate
    const fee = fromAmount.value * (exchangeRate.value.fee_percentage / 100)
    toAmount.value = (converted - (fee * exchangeRate.value.rate)).toFixed(2)
  } else {
    toAmount.value = 0
  }
}

const swapCurrencies = () => {
  const temp = fromCurrency.value
  fromCurrency.value = toCurrency.value
  toCurrency.value = temp
  updateRates()
}

const executeFiatConversion = async () => {
  loading.value = true
  try {
    // First get quote
    const quoteResponse = await $fetch('/gateway/exchange-service/fiat/quote', {
      method: 'POST',
      body: {
        from_currency: fromCurrency.value,
        to_currency: toCurrency.value,
        amount: fromAmount.value
      }
    })

    // Execute conversion
    const conversionResponse = await $fetch('/gateway/exchange-service/fiat/execute', {
      method: 'POST', 
      body: {
        from_wallet_id: 'user-wallet-' + fromCurrency.value.toLowerCase(),
        to_wallet_id: 'user-wallet-' + toCurrency.value.toLowerCase(),
        from_currency: fromCurrency.value,
        to_currency: toCurrency.value,
        amount: fromAmount.value
      }
    })

    // Add to recent conversions
    recentConversions.value.unshift({
      id: Date.now(),
      from_amount: fromAmount.value,
      from_currency: fromCurrency.value,
      to_amount: toAmount.value,
      to_currency: toCurrency.value,
      exchange_rate: exchangeRate.value.rate.toFixed(4),
      fee: (fromAmount.value * (exchangeRate.value.fee_percentage / 100)).toFixed(2),
      status: 'completed',
      created_at: new Date().toISOString()
    })

    // Show success notification
    showNotification('Conversion réalisée avec succès!', 'success')
    
  } catch (error) {
    console.error('Conversion error:', error)
    showNotification('Erreur lors de la conversion', 'error')
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

const showNotification = (message, type) => {
  // Simple notification - in real app would use a proper notification system
  console.log(`${type.toUpperCase()}: ${message}`)
}

// Lifecycle
onMounted(() => {
  updateRates()
})

// SEO
definePageMeta({
  title: 'Conversion Fiat - CryptoBank',
  description: 'Convertissez instantanément entre différentes monnaies fiduciaires avec des frais réduits'
})
</script>