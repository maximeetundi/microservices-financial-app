<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          Gestion des Instances d'Agrégateurs
        </h1>
        <p class="mt-2 text-gray-600 dark:text-gray-400">
          Gérez les instances de paiement et leurs hot wallets associés
        </p>
      </div>

      <!-- Actions -->
      <div class="mb-6 flex justify-between items-center">
        <div class="flex gap-3">
          <button
            @click="openCreateModal"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
          >
            <icon name="heroicons:plus" class="w-5 h-5" />
            Nouvelle Instance
          </button>
        </div>

        <div class="flex gap-3">
          <select
            v-model="filterAggregator"
            class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
          >
            <option value="">Tous les agrégateurs</option>
            <option v-for="agg in aggregators" :key="agg.id" :value="agg.provider_code">
              {{ agg.provider_name }}
            </option>
          </select>
        </div>
      </div>

      <!-- Instances List -->
      <div class="grid gap-6">
        <div
          v-for="instance in filteredInstances"
          :key="instance.id"
          class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-200 dark:border-gray-700"
        >
          <!-- Instance Header -->
          <div class="flex justify-between items-start mb-4">
            <div class="flex items-center gap-4">
              <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                <icon name="heroicons:server" class="w-8 h-8 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ instance.instance_name }}
                </h3>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ instance.provider_name }} - Priority: {{ instance.priority }}
                </p>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <!-- Status Badge -->
              <span
                :class="[
                  'px-3 py-1 rounded-full text-sm font-medium',
                  instance.enabled
                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                    : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                ]"
              >
                {{ instance.enabled ? 'Actif' : 'Inactif' }}
              </span>

              <!-- Test Mode -->
              <span
                v-if="instance.is_test_mode"
                class="px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
              >
                Test
              </span>

              <!-- Actions -->
              <button
                @click="openInstanceEdit(instance)"
                class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
              >
                <icon name="heroicons:pencil" class="w-5 h-5 text-gray-600 dark:text-gray-400" />
              </button>
            </div>
          </div>

          <!-- Statistics -->
          <div class="grid grid-cols-4 gap-4 mb-6 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Transactions</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">
                {{ instance.total_transactions || 0 }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Volume Total</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">
                {{ formatCurrency(instance.total_volume || 0) }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Usage Quotidien</p>
              <p class="text-lg font-bold text-gray-900 dark:text-white">
                {{ instance.daily_limit ? `${Math.round((instance.daily_usage / instance.daily_limit) * 100)}%` : 'N/A' }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Dernière Utilisation</p>
              <p class="text-sm text-gray-900 dark:text-white">
                {{ instance.last_used_at ? formatDate(instance.last_used_at) : 'Jamais' }}
              </p>
            </div>
          </div>

          <!-- Wallets Section -->
          <div>
            <div class="flex justify-between items-center mb-3">
              <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-300">
                Hot Wallets Liés ({{ instance.wallets?.length || 0 }})
              </h4>
              <button
                @click="openLinkWallet(instance)"
                class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 flex items-center gap-1"
              >
                <icon name="heroicons:plus-circle" class="w-4 h-4" />
                Lier un Wallet
              </button>
            </div>

            <div v-if="instance.wallets && instance.wallets.length > 0" class="space-y-2">
              <div
                v-for="wallet in instance.wallets"
                :key="wallet.id"
                class="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700"
              >
                <div class="flex items-center gap-3">
                  <icon
                    :name="wallet.is_primary ? 'heroicons:star-solid' : 'heroicons:star'"
                    :class="[
                      'w-5 h-5',
                      wallet.is_primary ? 'text-yellow-500' : 'text-gray-400'
                    ]"
                  />
                  <div>
                    <p class="font-medium text-gray-900 dark:text-white">
                      {{ wallet.wallet_currency }} Wallet
                      <span v-if="wallet.is_primary" class="text-xs text-yellow-600">(Primary)</span>
                    </p>
                    <p class="text-sm text-gray-600 dark:text-gray-400">
                      Balance: {{ formatCurrency(wallet.wallet_balance) }} {{ wallet.wallet_currency }}
                    </p>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <!-- Auto-Recharge Badge -->
                  <span
                    v-if="wallet.auto_recharge_enabled"
                    class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                  >
                    Auto-Recharge
                  </span>

                  <!-- Status -->
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      getWalletStatusClass(wallet)
                    ]"
                  >
                    {{ getWalletStatus(wallet) }}
                  </span>

                  <button
                    @click="openWalletConfig(instance, wallet)"
                    class="p-1 hover:bg-gray-200 dark:hover:bg-gray-800 rounded"
                  >
                    <icon name="heroicons:cog-6-tooth" class="w-4 h-4 text-gray-600 dark:text-gray-400" />
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="text-center py-6 text-gray-500 dark:text-gray-400">
              Aucun wallet lié à cette instance
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-if="filteredInstances.length === 0"
        class="text-center py-12 bg-white dark:bg-gray-800 rounded-xl shadow-lg"
      >
        <icon name="heroicons:server" class="w-16 h-16 mx-auto text-gray-400 mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
          Aucune instance trouvée
        </h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          Créez votre première instance pour commencer
        </p>
        <button
          @click="showCreateModal = true"
          class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Créer une instance
        </button>
      </div>
    </div>

    <!-- Modals -->
    <CreateInstanceModal
      :is-open="showCreateModal"
      :instance="selectedInstance"
      :aggregators="aggregators"
      @close="showCreateModal = false"
      @submit="handleInstanceSubmit"
    />

    <LinkWalletModal
      :is-open="showLinkWalletModal"
      :instance="selectedInstance"
      @close="showLinkWalletModal = false"
      @submit="handleLinkWallet"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAggregatorInstances } from '@/composables/useAggregatorInstances'
import { usePaymentProviders } from '@/composables/usePaymentProviders'
import CreateInstanceModal from '@/components/admin/CreateInstanceModal.vue'
import LinkWalletModal from '@/components/admin/LinkWalletModal.vue'

// Composables
const { 
  instances, 
  loading: instancesLoading, 
  fetchInstances, 
  createInstance, 
  updateInstance,
  linkWallet,
  unlinkWallet,
  updateWallet
} = useAggregatorInstances()

const { 
  providers: aggregators, 
  loading: providersLoading, 
  fetchPaymentProviders 
} = usePaymentProviders()

// State
const filterAggregator = ref('')
const showCreateModal = ref(false)
const showLinkWalletModal = ref(false)
const selectedInstance = ref<any>(null)

// Computed
const filteredInstances = computed(() => {
  if (!filterAggregator.value) return instances.value
  return instances.value.filter(i => i.provider_code === filterAggregator.value)
})

const loading = computed(() => instancesLoading.value || providersLoading.value)

// Methods
const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('fr-FR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getWalletStatus = (wallet: any) => {
  if (!wallet.enabled) return 'Désactivé'
  // Fix: wallet_balance might be undefined if not populated instantly
  const balance = wallet.wallet_balance || 0
  if (wallet.min_balance && balance < wallet.min_balance) return 'Insuffisant'
  if (wallet.max_balance && balance > wallet.max_balance) return 'Trop élevé'
  return 'Disponible'
}

const getWalletStatusClass = (wallet: any) => {
  const status = getWalletStatus(wallet)
  if (status === 'Disponible') return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
  if (status === 'Désactivé') return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
  return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
}

// Actions
const openCreateModal = () => {
  selectedInstance.value = null
  showCreateModal.value = true
}

const openInstanceEdit = (instance: any) => {
  selectedInstance.value = instance
  showCreateModal.value = true
}

const openLinkWallet = (instance: any) => {
  selectedInstance.value = instance
  showLinkWalletModal.value = true
}

const openWalletConfig = (instance: any, wallet: any) => {
  // TODO: Implement wallet config modal
  console.log('Configure wallet', wallet.id)
}

const handleInstanceSubmit = async (data: any) => {
  try {
    if (selectedInstance.value) {
      await updateInstance(selectedInstance.value.id, {
        ...data,
        aggregator_id: selectedInstance.value.provider_id
      })
    } else {
      await createInstance(data)
    }
    showCreateModal.value = false
    // Refresh to get latest state
    await fetchInstances()
  } catch (err) {
    console.error('Failed to save instance:', err)
    // You might want to show a toast here
  }
}

const handleLinkWallet = async (data: any) => {
  if (!selectedInstance.value) return
  
  try {
    await linkWallet(selectedInstance.value.id, data)
    showLinkWalletModal.value = false
    await fetchInstances()
  } catch (err) {
    console.error('Failed to link wallet:', err)
  }
}

// Load data on mount
onMounted(async () => {
  await Promise.all([
    fetchInstances(),
    fetchPaymentProviders()
  ])
})
</script>
