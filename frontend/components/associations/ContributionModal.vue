<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">Payer une cotisation</h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">✕</button>
        </div>
      </div>

      <div class="p-6 space-y-4">
        <!-- Amount -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Montant</label>
          <div class="relative">
            <input v-model.number="form.amount" type="number" step="100"
                   class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                   placeholder="0" />
            <span class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500">{{ currency }}</span>
          </div>
        </div>

        <!-- Period -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Période</label>
          <select v-model="form.period" 
                  class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500">
            <option value="janvier_2026">Janvier 2026</option>
            <option value="fevrier_2026">Février 2026</option>
            <option value="mars_2026">Mars 2026</option>
            <option value="avril_2026">Avril 2026</option>
          </select>
        </div>

        <!-- Wallet Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Portefeuille source</label>
          <div v-if="loadingWallets" class="text-center py-4 text-gray-500">Chargement...</div>
          <div v-else class="space-y-2">
            <div v-for="wallet in wallets" :key="wallet.id"
                 @click="form.wallet_id = wallet.id"
                 class="p-3 border rounded-lg cursor-pointer transition-all"
                 :class="form.wallet_id === wallet.id ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-900/20' : 'border-gray-200 dark:border-gray-700 hover:border-indigo-300'">
              <div class="flex justify-between items-center">
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">{{ wallet.currency }}</div>
                  <div class="text-xs text-gray-500">{{ wallet.type }}</div>
                </div>
                <div class="text-right">
                  <div class="font-bold text-gray-900 dark:text-white">{{ formatBalance(wallet.balance, wallet.currency) }}</div>
                  <div v-if="form.wallet_id === wallet.id" class="text-xs text-indigo-600">Sélectionné</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- PIN -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Code PIN</label>
          <input v-model="form.pin" type="password" maxlength="5" inputmode="numeric"
                 class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-center text-2xl tracking-widest focus:ring-2 focus:ring-indigo-500"
                 placeholder="•••••" />
        </div>

        <!-- Error -->
        <div v-if="error" class="p-3 bg-red-50 dark:bg-red-900/20 text-red-600 rounded-lg text-sm">
          {{ error }}
        </div>
      </div>

      <div class="p-6 border-t border-gray-100 dark:border-gray-700 flex space-x-3">
        <button @click="$emit('close')" 
                class="flex-1 px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-700">
          Annuler
        </button>
        <button @click="submit" :disabled="loading || !isValid"
                class="flex-1 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white rounded-lg font-medium flex items-center justify-center">
          <span v-if="loading" class="animate-spin mr-2 h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
          Payer {{ formatBalance(form.amount, currency) }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

const props = defineProps<{
  show: boolean
  associationId: string
  currency: string
}>()

const emit = defineEmits(['close', 'success'])

const { walletApi, associationApi } = useApi()

const loading = ref(false)
const loadingWallets = ref(true)
const error = ref('')
const wallets = ref<any[]>([])

const form = ref({
  amount: 0,
  period: 'janvier_2026',
  wallet_id: '',
  pin: '',
  description: ''
})

const isValid = computed(() => {
  return form.value.amount > 0 && form.value.wallet_id && form.value.pin.length === 5
})

const formatBalance = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount || 0)
}

const loadWallets = async () => {
  loadingWallets.value = true
  try {
    const response = await walletApi.getWallets()
    wallets.value = response.data?.filter((w: any) => w.currency === props.currency) || response.data || []
    if (wallets.value.length > 0 && !form.value.wallet_id) {
      form.value.wallet_id = wallets.value[0].id
    }
  } catch (err) {
    console.error('Failed to load wallets', err)
  } finally {
    loadingWallets.value = false
  }
}

const submit = async () => {
  if (!isValid.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    await associationApi.recordContribution(props.associationId, {
      wallet_id: form.value.wallet_id,
      pin: form.value.pin,
      amount: form.value.amount,
      period: form.value.period,
      description: `Cotisation ${form.value.period}`
    })
    emit('success')
    emit('close')
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Échec du paiement. Vérifiez votre solde et votre PIN.'
  } finally {
    loading.value = false
  }
}

watch(() => props.show, (val) => {
  if (val) {
    loadWallets()
    form.value.pin = ''
    error.value = ''
  }
})

onMounted(() => {
  if (props.show) {
    loadWallets()
  }
})
</script>
