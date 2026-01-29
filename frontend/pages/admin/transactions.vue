<template>
  <NuxtLayout name="admin">
    <div class="p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">Transactions</h1>
        <p class="text-slate-400">Historique global des transactions de la plateforme.</p>
      </div>

      <!-- Filters (Placeholder) -->
      <div class="mb-6 flex gap-4">
        <div class="relative flex-1 max-w-md">
          <input 
            type="text" 
            placeholder="Rechercher par ID, référence..." 
            class="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 text-white focus:outline-none focus:border-indigo-500"
          >
        </div>
        <button class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-xl text-slate-300 hover:text-white">
          Filtres
        </button>
      </div>

      <!-- Content -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
      </div>

      <div v-else class="bg-slate-800/50 backdrop-blur-xl rounded-2xl border border-slate-700/50 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-900/50 text-slate-400 text-xs uppercase font-medium">
              <tr>
                <th class="px-6 py-4">Date</th>
                <th class="px-6 py-4">Type</th>
                <th class="px-6 py-4">Compte / Wallet</th>
                <th class="px-6 py-4 text-right">Montant</th>
                <th class="px-6 py-4 text-right">Référence</th>
                <th class="px-6 py-4">Description</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr v-for="tx in transactions" :key="tx.id" class="hover:bg-slate-700/20 transition-colors">
                <td class="px-6 py-4 whitespace-nowrap text-slate-300">
                  {{ formatDate(tx.created_at) }}
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="{
                      'bg-emerald-500/10 text-emerald-400': tx.type === 'credit' || tx.type === 'deposit',
                      'bg-red-500/10 text-red-400': tx.type === 'debit' || tx.type === 'withdrawal',
                      'bg-blue-500/10 text-blue-400': tx.type === 'exchange'
                    }">
                    {{ tx.type }}
                  </span>
                </td>
                <td class="px-6 py-4 text-slate-300">
                  <div class="flex flex-col">
                    <span class="font-medium text-white">{{ tx.currency }}</span>
                    <span class="text-xs text-slate-500">{{ tx.platform_account_id }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-right font-mono" 
                    :class="tx.amount >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ formatAmount(tx.amount, tx.currency) }}
                </td>
                <td class="px-6 py-4 text-right text-xs font-mono text-slate-500">
                  {{ tx.reference_id || '-' }}
                </td>
                <td class="px-6 py-4 text-sm text-slate-400 max-w-xs truncate">
                  {{ tx.description }}
                </td>
              </tr>
              <tr v-if="transactions.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-slate-400">
                  Aucune transaction trouvée.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="px-6 py-4 border-t border-slate-700/50 flex justify-between items-center">
          <button @click="prevPage" :disabled="offset === 0" class="text-slate-400 hover:text-white disabled:opacity-50">
            ← Précédent
          </button>
          <span class="text-slate-500 text-sm">Page {{ currentPage }}</span>
          <button @click="nextPage" :disabled="transactions.length < limit" class="text-slate-400 hover:text-white disabled:opacity-50">
            Suivant →
          </button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useApi } from '@/composables/useApi'

const { adminPlatformAPI } = useApi()

const loading = ref(true)
const transactions = ref([])
const limit = ref(50)
const offset = ref(0)

const currentPage = computed(() => Math.floor(offset.value / limit.value) + 1)

const fetchTransactions = async () => {
  loading.value = true
  try {
    const response = await adminPlatformAPI.getTransactions(limit.value, offset.value)
    transactions.value = response.data?.transactions || []
  } catch (error) {
    console.error("Failed to fetch transactions", error)
  } finally {
    loading.value = false
  }
}

const nextPage = () => {
  offset.value += limit.value
  fetchTransactions()
}

const prevPage = () => {
  if (offset.value >= limit.value) {
    offset.value -= limit.value
    fetchTransactions()
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('fr-FR')
}

const formatAmount = (amount, currency) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'EUR' }).format(amount)
}

onMounted(() => {
  fetchTransactions()
})
</script>
