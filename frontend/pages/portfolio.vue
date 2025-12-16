<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">Portfolio</h1>
          <div class="flex items-center space-x-4">
            <NuxtLink to="/exchange/trading" 
                      class="bg-blue-600 text-gray-900 px-4 py-2 rounded-lg hover:bg-blue-700">
              Trade
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Portfolio Summary -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-xl shadow-lg p-6">
          <div class="text-sm font-medium text-gray-500 mb-1">Total Value</div>
          <div class="text-2xl font-bold text-gray-900">${{ portfolio.totalValue?.toLocaleString() || '0' }}</div>
          <div class="text-sm mt-2" :class="portfolio.performance?.dayReturn >= 0 ? 'text-green-600' : 'text-red-600'">
            {{ portfolio.performance?.dayReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.dayReturn }}% today
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6">
          <div class="text-sm font-medium text-gray-500 mb-1">Day P&L</div>
          <div class="text-2xl font-bold" :class="portfolio.performance?.dayReturnValue >= 0 ? 'text-green-600' : 'text-red-600'">
            ${{ portfolio.performance?.dayReturnValue?.toLocaleString() || '0' }}
          </div>
          <div class="text-sm text-gray-500 mt-2">24h Change</div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6">
          <div class="text-sm font-medium text-gray-500 mb-1">Total Return</div>
          <div class="text-2xl font-bold" :class="portfolio.performance?.totalReturn >= 0 ? 'text-green-600' : 'text-red-600'">
            {{ portfolio.performance?.totalReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.totalReturn }}%
          </div>
          <div class="text-sm text-gray-500 mt-2">${{ portfolio.performance?.totalReturnValue?.toLocaleString() || '0' }}</div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6">
          <div class="text-sm font-medium text-gray-500 mb-1">Assets</div>
          <div class="text-2xl font-bold text-gray-900">{{ portfolio.holdings?.length || 0 }}</div>
          <div class="text-sm text-gray-500 mt-2">Currencies</div>
        </div>
      </div>

      <!-- Performance Chart -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        <div class="lg:col-span-2 bg-white rounded-xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Performance</h2>
          <div class="grid grid-cols-4 gap-4 mb-6">
            <div class="text-center">
              <div class="text-sm text-gray-500">1 Day</div>
              <div class="font-semibold" :class="portfolio.performance?.dayReturn >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ portfolio.performance?.dayReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.dayReturn }}%
              </div>
            </div>
            <div class="text-center">
              <div class="text-sm text-gray-500">1 Week</div>
              <div class="font-semibold" :class="portfolio.performance?.weekReturn >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ portfolio.performance?.weekReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.weekReturn }}%
              </div>
            </div>
            <div class="text-center">
              <div class="text-sm text-gray-500">1 Month</div>
              <div class="font-semibold" :class="portfolio.performance?.monthReturn >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ portfolio.performance?.monthReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.monthReturn }}%
              </div>
            </div>
            <div class="text-center">
              <div class="text-sm text-gray-500">1 Year</div>
              <div class="font-semibold" :class="portfolio.performance?.yearReturn >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ portfolio.performance?.yearReturn >= 0 ? '+' : '' }}{{ portfolio.performance?.yearReturn }}%
              </div>
            </div>
          </div>
          <div class="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
            <span class="text-gray-500">Performance Chart (Chart.js integration)</span>
          </div>
        </div>

        <!-- Asset Allocation -->
        <div class="bg-white rounded-xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Asset Allocation</h2>
          <div class="space-y-4">
            <div v-for="holding in portfolio.holdings" :key="holding.currency" class="space-y-2">
              <div class="flex justify-between items-center">
                <span class="font-medium">{{ holding.currency }}</span>
                <span class="text-sm text-gray-500">{{ holding.percentage }}%</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div class="h-2 rounded-full transition-all duration-300"
                     :class="getCurrencyColor(holding.currency)"
                     :style="{ width: holding.percentage + '%' }">
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Holdings Table -->
      <div class="bg-white rounded-xl shadow-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-semibold text-gray-900">Holdings</h2>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Asset</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Value</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Allocation</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">24h Change</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="holding in portfolio.holdings" :key="holding.currency" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-8 w-8 rounded-full flex items-center justify-center text-gray-900 text-sm font-bold"
                         :class="getCurrencyColor(holding.currency)">
                      {{ holding.currency.substring(0, 2) }}
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ holding.currency }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatAmount(holding.amount, holding.currency) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${{ holding.value?.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ holding.percentage }}%
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm" 
                    :class="holding.change_24h >= 0 ? 'text-green-600' : 'text-red-600'">
                  <div>{{ holding.change_24h >= 0 ? '+' : '' }}{{ holding.change_24h }}%</div>
                  <div class="text-xs">
                    ${{ holding.changeValue24h >= 0 ? '+' : '' }}{{ holding.changeValue24h?.toLocaleString() }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  <button @click="buyCurrency(holding.currency)" 
                          class="text-green-600 hover:text-green-900">Buy</button>
                  <button @click="sellCurrency(holding.currency)" 
                          class="text-red-600 hover:text-red-900">Sell</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Page meta
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

// Reactive data
const portfolio = ref({
  totalValue: 0,
  holdings: [],
  performance: {}
})

// Methods
const getCurrencyColor = (currency) => {
  const colors = {
    'BTC': 'bg-orange-500',
    'ETH': 'bg-blue-500',
    'USD': 'bg-green-500',
    'EUR': 'bg-purple-500',
    'GBP': 'bg-red-500'
  }
  return colors[currency] || 'bg-gray-500'
}

const formatAmount = (amount, currency) => {
  if (currency === 'USD' || currency === 'EUR' || currency === 'GBP') {
    return amount?.toLocaleString()
  }
  return amount?.toFixed(6)
}

const buyCurrency = (currency) => {
  navigateTo(`/exchange/trading?buy=${currency}`)
}

const sellCurrency = (currency) => {
  navigateTo(`/exchange/trading?sell=${currency}`)
}

// Fetch data
const fetchPortfolio = async () => {
  try {
    const data = await $fetch('/api/trading/portfolio')
    portfolio.value = data
  } catch (error) {
    console.error('Error fetching portfolio:', error)
    // Fallback data for demo
    portfolio.value = {
      totalValue: 50000,
      holdings: [
        {
          currency: 'BTC',
          amount: 0.5,
          value: 21750,
          percentage: 43.5,
          change_24h: 2.3,
          changeValue24h: 500.25
        },
        {
          currency: 'ETH',
          amount: 10,
          value: 24500,
          percentage: 49.0,
          change_24h: -1.2,
          changeValue24h: -294.0
        },
        {
          currency: 'USD',
          amount: 3750,
          value: 3750,
          percentage: 7.5,
          change_24h: 0.0,
          changeValue24h: 0.0
        }
      ],
      performance: {
        totalReturn: 5.2,
        totalReturnValue: 520.0,
        dayReturn: 0.8,
        dayReturnValue: 80.0,
        weekReturn: 2.1,
        monthReturn: 8.5,
        yearReturn: 45.2
      }
    }
  }
}

// Lifecycle
onMounted(async () => {
  await fetchPortfolio()
})
</script>