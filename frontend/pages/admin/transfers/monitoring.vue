<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header with Real-time Stats -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Monitoring des Transferts Fiduciaires
        </h1>
        <p class="text-gray-600 dark:text-gray-400">
          Surveillance en temps réel des mouvements de fonds
        </p>
      </div>

      <!-- Real-Time Metrics -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border-l-4 border-blue-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-1">Dépôts Aujourd'hui</p>
              <p class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ metrics.todayDeposits.count }}
              </p>
              <p class="text-sm text-blue-600 dark:text-blue-400 mt-1">
                {{ formatCurrency(metrics.todayDeposits.volume) }}
              </p>
            </div>
            <icon name="heroicons:arrow-down-circle" class="w-12 h-12 text-blue-500" />
          </div>
        </div>

        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border-l-4 border-green-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-1">Retraits Aujourd'hui</p>
              <p class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ metrics.todayWithdrawals.count }}
              </p>
              <p class="text-sm text-green-600 dark:text-green-400 mt-1">
                {{ formatCurrency(metrics.todayWithdrawals.volume) }}
              </p>
            </div>
            <icon name="heroicons:arrow-up-circle" class="w-12 h-12 text-green-500" />
          </div>
        </div>

        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border-l-4 border-purple-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-1">En Cours</p>
              <p class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ metrics.pending.count }}
              </p>
              <p class="text-sm text-purple-600 dark:text-purple-400 mt-1">
                {{ formatCurrency(metrics.pending.volume) }}
              </p>
            </div>
            <icon name="heroicons:clock" class="w-12 h-12 text-purple-500" />
          </div>
        </div>

        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border-l-4 border-red-500">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-1">Échecs/Annulations</p>
              <p class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ metrics.failed.count }}
              </p>
              <p class="text-sm text-red-600 dark:text-red-400 mt-1">
                Taux: {{ failureRate }}%
              </p>
            </div>
            <icon name="heroicons:exclamation-triangle" class="w-12 h-12 text-red-500" />
          </div>
        </div>
      </div>

      <!-- Platform Wallets Status -->
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 mb-8">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <icon name="heroicons:wallet" class="w-6 h-6" />
          État des Hot Wallets (Operations)
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="wallet in platformWallets"
            :key="wallet.id"
            class="p-4 rounded-lg border border-gray-200 dark:border-gray-700"
          >
            <div class="flex justify-between items-start mb-3">
              <div>
                <p class="text-lg font-bold text-gray-900 dark:text-white">
                  {{ wallet.currency }}
                </p>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ wallet.name }}
                </p>
              </div>
              <span
                :class="[
                  'px-2 py-1 text-xs rounded-full',
                  getWalletHealthClass(wallet)
                ]"
              >
                {{ getWalletHealth(wallet) }}
              </span>
            </div>

            <div class="space-y-2">
              <div>
                <p class="text-xs text-gray-500 dark:text-gray-400">Balance Actuelle</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatCurrency(wallet.balance) }}
                </p>
              </div>

              <!-- Balance Bar -->
              <div class="relative pt-1">
                <div class="overflow-hidden h-2 text-xs flex rounded bg-gray-200 dark:bg-gray-700">
                  <div
                    :style="{ width: `${getBalancePercentage(wallet)}%` }"
                    :class="[
                      'shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center',
                      getBalanceBarColor(wallet)
                    ]"
                  ></div>
                </div>
                <divclass="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
                  <span>Min: {{ formatCurrency(wallet.min_balance || 0) }}</span>
                  <span>Max: {{ formatCurrency(wallet.max_balance || 0) }}</span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                <div>
                  <p class="text-xs text-gray-500">Dépôts 24h</p>
                  <p class="text-sm font-semibold text-blue-600">
                    {{ wallet.deposits_24h || 0 }}
                  </p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Retraits 24h</p>
                  <p class="text-sm font-semibold text-green-600">
                    {{ wallet.withdrawals_24h || 0 }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transactions Stream -->
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <icon name="heroicons:arrows-right-left" class="w-6 h-6" />
            Flux de Transactions en Temps Réel
          </h2>
          <div class="flex gap-3">
            <select
              v-model="filter.type"
              class="px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900"
            >
              <option value="">Tous les types</option>
              <option value="deposit">Dépôts</option>
              <option value="withdrawal">Retraits</option>
            </select>
            <select
              v-model="filter.status"
              class="px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900"
            >
              <option value="">Tous les statuts</option>
              <option value="pending">En cours</option>
              <option value="completed">Complété</option>
              <option value="failed">Échoué</option>
            </select>
          </div>
        </div>

        <!-- Live Transaction Feed -->
        <div class="space-y-3 max-h-96 overflow-y-auto">
          <div
            v-for="tx in filteredTransactions"
            :key="tx.id"
            class="p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <div class="flex justify-between items-start">
              <div class="flex items-start gap-3">
                <div
                  :class="[
                    'p-2 rounded-lg',
                    tx.type === 'deposit' ? 'bg-blue-100 dark:bg-blue-900' : 'bg-green-100 dark:bg-green-900'
                  ]"
                >
                  <icon
                    :name="tx.type === 'deposit' ? 'heroicons:arrow-down' : 'heroicons:arrow-up'"
                    :class="[
                      'w-5 h-5',
                      tx.type === 'deposit' ? 'text-blue-600' : 'text-green-600'
                    ]"
                  />
                </div>

                <div>
                  <p class="font-semibold text-gray-900 dark:text-white">
                    {{ tx.type === 'deposit' ? 'Dépôt' : 'Retrait' }} - {{ tx.provider }}
                  </p>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    User: {{ tx.user_id }} • Ref: {{ tx.reference }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-500 mt-1">
                    {{ formatDateTime(tx.created_at) }}
                  </p>
                </div>
              </div>

              <div class="text-right">
                <p class="text-lg font-bold text-gray-900 dark:text-white">
                  {{ formatCurrency(tx.amount) }} {{ tx.currency }}
                </p>
                <span
                  :class="[
                    'inline-block px-3 py-1 text-xs rounded-full mt-1',
                    getStatusClass(tx.status)
                  ]"
                >
                  {{ getStatusLabel(tx.status) }}
                </span>
              </div>
            </div>

            <!-- Fund Flow Indicator -->
            <div v-if="tx.type === 'deposit'" class="mt-3 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
              <icon name="heroicons:building-library" class="w-4 h-4" />
              <span>Provider</span>
              <icon name="heroicons:arrow-right" class="w-4 h-4 text-blue-500" />
              <icon name="heroicons:wallet" class="w-4 h-4 text-orange-500" />
              <span>Hot Wallet</span>
              <icon name="heroicons:arrow-right" class="w-4 h-4 text-blue-500" />
              <icon name="heroicons:user" class="w-4 h-4 text-purple-500" />
              <span>User</span>
            </div>
            <div v-else class="mt-3 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
              <icon name="heroicons:user" class="w-4 h-4 text-purple-500" />
              <span>User</span>
              <icon name="heroicons:arrow-right" class="w-4 h-4 text-green-500" />
              <icon name="heroicons:wallet" class="w-4 h-4 text-orange-500" />
              <span>Hot Wallet</span>
              <icon name="heroicons:arrow-right" class="w-4 h-4 text-green-500" />
              <icon name="heroicons:building-library" class="w-4 h-4" />
              <span>Provider</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const metrics = ref({
  todayDeposits: { count: 0, volume: 0 },
  todayWithdrawals: { count: 0, volume: 0 },
  pending: { count: 0, volume: 0 },
  failed: { count: 0, volume: 0 }
})

const platformWallets = ref([])
const transactions = ref([])
const filter = ref({ type: '', status: '' })

const failureRate = computed(() => {
  const total = metrics.value.todayDeposits.count + metrics.value.todayWithdrawals.count
  return total > 0 ? ((metrics.value.failed.count / total) * 100).toFixed(2) : 0
})

const filteredTransactions = computed(() => {
  return transactions.value.filter(tx => {
    if (filter.value.type && tx.type !== filter.value.type) return false
    if (filter.value.status && tx.status !== filter.value.status) return false
    return true
  })
})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('fr-FR').format(amount)
}

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString('fr-FR')
}

const getStatusClass = (status: string) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    completed: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    failed: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusLabel = (status: string) => {
  const labels = {
    pending: 'En cours',
    completed: 'Complété',
    failed: 'Échoué'
  }
  return labels[status] || status
}

const getWalletHealth = (wallet: any) => {
  if (wallet.balance < (wallet.min_balance || 0)) return 'Critique'
  if (wallet.balance < (wallet.min_balance || 0) * 2) return 'Bas'
  return 'Normal'
}

const getWalletHealthClass = (wallet: any) => {
  const health = getWalletHealth(wallet)
  if (health === 'Critique') return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
  if (health === 'Bas') return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
  return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
}

const getBalancePercentage = (wallet: any) => {
  const max = wallet.max_balance || wallet.balance * 2
  return Math.min(100, (wallet.balance / max) * 100)
}

const getBalanceBarColor = (wallet: any) => {
  const health = getWalletHealth(wallet)
  if (health === 'Critique') return 'bg-red-500'
  if (health === 'Bas') return 'bg-yellow-500'
  return 'bg-green-500'
}

let refreshInterval: any

onMounted(() => {
  // Initial load
  fetchMetrics()
  fetchWallets()
  fetchTransactions()

  // Refresh every 5 seconds
  refreshInterval = setInterval(() => {
    fetchMetrics()
    fetchWallets()
    fetchTransactions()
  }, 5000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})

const fetchMetrics = async () => {
  // TODO: Fetch from API
}

const fetchWallets = async () => {
  // TODO: Fetch platform wallets
}

const fetchTransactions = async () => {
  // TODO: Fetch recent transactions
}
</script>
