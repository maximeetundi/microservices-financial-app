<template>
  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="close" class="relative z-50">
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black/30 backdrop-blur-sm" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel class="w-full max-w-2xl transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 shadow-xl transition-all">
              <DialogTitle as="h3" class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                Lier un Wallet à {{ instance?.instance_name }}
              </DialogTitle>

              <form @submit.prevent="submit" class="space-y-6">
                <!-- Wallet Selection -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Hot Wallet (Operations) *
                  </label>
                  <select
                    v-model="form.hot_wallet_id"
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    required
                  >
                    <option value="">Sélectionner un wallet operations</option>
                    <option v-for="wallet in availableWallets" :key="wallet.id" :value="wallet.id">
                      {{ wallet.currency }} - {{ wallet.name }} (Balance: {{ formatCurrency(wallet.balance) }})
                    </option>
                  </select>
                </div>

                <!-- Primary & Priority -->
                <div class="grid grid-cols-2 gap-4">
                  <div class="flex items-center">
                    <input
                      v-model="form.is_primary"
                      type="checkbox"
                      class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                    />
                    <label class="ml-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                      <icon name="heroicons:star-solid" class="w-4 h-4 inline text-yellow-500" />
                      Wallet Principal
                    </label>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Priorité
                    </label>
                    <input
                      v-model.number="form.priority"
                      type="number"
                      min="1"
                      max="100"
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                <!-- Balance Constraints -->
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Balance Minimum
                    </label>
                    <input
                      v-model.number="form.min_balance"
                      type="number"
                      placeholder="0"
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Balance Maximum
                    </label>
                    <input
                      v-model.number="form.max_balance"
                      type="number"
                      placeholder="Illimité"
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                <!-- Auto-Recharge Section -->
                <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4">
                  <div class="flex items-center justify-between">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Auto-Recharge depuis Réserve
                    </label>
                    <input
                      v-model="form.auto_recharge_enabled"
                      type="checkbox"
                      class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  <div v-if="form.auto_recharge_enabled" class="space-y-4 pt-2 border-t border-gray-200 dark:border-gray-700">
                    <!-- Source Wallet -->
                    <div>
                      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Cold Wallet Source (Réserve)
                      </label>
                      <select
                        v-model="form.recharge_source_wallet_id"
                        class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                      >
                        <option value="">Sélectionner wallet réserve</option>
                        <option v-for="wallet in reserveWallets" :key="wallet.id" :value="wallet.id">
                          {{ wallet.currency }} - {{ wallet.name }} ({{ formatCurrency(wallet.balance) }})
                        </option>
                      </select>
                    </div>

                    <!-- Thresholds -->
                    <div class="grid grid-cols-2 gap-4">
                      <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                          Seuil de Déclenchement
                        </label>
                        <input
                          v-model.number="form.recharge_threshold"
                          type="number"
                          placeholder="Ex: 50000"
                          class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                        />
                        <p class="mt-1 text-xs text-gray-500">Recharger si balance &lt;</p>
                      </div>
                      <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                          Montant Cible
                        </label>
                        <input
                          v-model.number="form.recharge_target"
                          type="number"
                          placeholder="Ex: 200000"
                          class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                        />
                        <p class="mt-1 text-xs text-gray-500">Recharger jusqu'à</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Actions -->
                <div class="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                  <button
                    type="button"
                    @click="close"
                    class="px-6 py-3 rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    type="submit"
                    :disabled="loading"
                    class="px-6 py-3 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors disabled:opacity-50 flex items-center gap-2"
                  >
                    <icon v-if="loading" name="svg-spinners:ring-resize" class="w-5 h-5" />
                    Lier le Wallet
                  </button>
                </div>
              </form>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'

const props = defineProps<{
  isOpen: boolean
  instance: any
  platformWallets: any[]
}>()

const emit = defineEmits(['close', 'submit'])

const loading = ref(false)

const form = reactive({
  hot_wallet_id: '',
  is_primary: false,
  priority: 50,
  min_balance: null,
  max_balance: null,
  auto_recharge_enabled: false,
  recharge_threshold: null,
  recharge_target: null,
  recharge_source_wallet_id: ''
})

// Filter wallets by type
const availableWallets = computed(() => {
  return props.platformWallets.filter(w => w.account_type === 'operations')
})

const reserveWallets = computed(() => {
  return props.platformWallets.filter(w => w.account_type === 'reserve')
})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const close = () => {
  emit('close')
}

const submit = async () => {
  loading.value = true
  try {
    await emit('submit', form)
    close()
  } finally {
    loading.value = false
  }
}
</script>
