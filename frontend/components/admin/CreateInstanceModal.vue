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
                {{ instance ? 'Modifier l\'Instance' : 'Créer une Nouvelle Instance' }}
              </DialogTitle>

              <form @submit.prevent="submit" class="space-y-6">
                <!-- Aggregator Selection -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Agrégateur *
                  </label>
                  <select
                    v-model="form.aggregator_id"
                    :disabled="!!instance"
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
                    required
                  >
                    <option value="">Sélectionner un agrégateur</option>
                    <option v-for="agg in aggregators" :key="agg.id" :value="agg.id">
                      {{ agg.provider_name }}
                    </option>
                  </select>
                </div>

                <!-- Instance Name -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Nom de l'Instance *
                  </label>
                  <input
                    v-model="form.instance_name"
                    type="text"
                    placeholder="Ex: Flutterwave Nigeria Main"
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    required
                  />
                </div>

                <!-- API Credentials -->
                <div class="space-y-4 border-t border-gray-200 dark:border-gray-700 pt-4">
                  <h4 class="text-sm font-medium text-gray-900 dark:text-white">
                    Identifiants API
                  </h4>
                  
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <!-- Common Keys -->
                    <div v-if="shouldShowField('api_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">API Key</label>
                       <input v-model="form.credentials.api_key" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    <div v-if="shouldShowField('public_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Public Key</label>
                       <input v-model="form.credentials.public_key" type="text" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    <div v-if="shouldShowField('secret_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Secret Key</label>
                       <input v-model="form.credentials.secret_key" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    
                    <!-- Lygos / Shop Specific -->
                    <div v-if="shouldShowField('shop_name')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Shop Name</label>
                       <input v-model="form.credentials.shop_name" type="text" placeholder="Ex: Ma Boutique" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>

                    <!-- Client ID / Secret -->
                    <div v-if="shouldShowField('client_id')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Client ID</label>
                       <input v-model="form.credentials.client_id" type="text" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    <div v-if="shouldShowField('client_secret')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Client Secret</label>
                       <input v-model="form.credentials.client_secret" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>

                    <!-- Other Provider Specifics -->
                     <div v-if="shouldShowField('site_id')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Site ID</label>
                       <input v-model="form.credentials.site_id" type="text" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    <div v-if="shouldShowField('merchant_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Merchant Key</label>
                       <input v-model="form.credentials.merchant_key" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                     <div v-if="shouldShowField('subscription_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Subscription Key</label>
                       <input v-model="form.credentials.subscription_key" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                     <div v-if="shouldShowField('api_user')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">API User</label>
                       <input v-model="form.credentials.api_user" type="text" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                    
                    <!-- Advanced -->
                     <div v-if="shouldShowField('encryption_key')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Encryption Key</label>
                       <input v-model="form.credentials.encryption_key" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                     <div v-if="shouldShowField('webhook_secret')">
                       <label class="block text-xs font-medium text-gray-500 mb-1">Webhook Secret</label>
                       <input v-model="form.credentials.webhook_secret" type="password" class="w-full px-3 py-2 text-sm rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900" />
                    </div>
                  </div>
                </div>

                <!-- Priority & Mode -->
                <div class="grid grid-cols-2 gap-4">
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
                  <div class="flex items-center pt-8">
                    <input
                      v-model="form.is_test_mode"
                      type="checkbox"
                      class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                    />
                    <label class="ml-3 text-sm font-medium text-gray-700 dark:text-gray-300">
                      Mode Test
                    </label>
                  </div>
                </div>

                <!-- Limits -->
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Limite Quotidienne
                    </label>
                    <input
                      v-model.number="form.daily_limit"
                      type="number"
                      placeholder="Illimité"
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Limite Mensuelle
                    </label>
                    <input
                      v-model.number="form.monthly_limit"
                      type="number"
                      placeholder="Illimité"
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                <!-- Notes -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Notes
                  </label>
                  <textarea
                    v-model="form.notes"
                    rows="3"
                    placeholder="Notes internes..."
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500"
                  ></textarea>
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
                    {{ instance ? 'Mettre à jour' : 'Créer l\'Instance' }}
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
import { ref, reactive, watch } from 'vue'
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'

const props = defineProps<{
  isOpen: boolean
  instance?: any
  aggregators: any[]
}>()

const emit = defineEmits(['close', 'submit'])

const loading = ref(false)

const form = reactive({
  aggregator_id: '',
  instance_name: '',
  credentials: {
    public_key: '',
    secret_key: '',
    api_key: '',
    shop_name: '',
    client_id: '',
    client_secret: '',
    site_id: '',
    merchant_key: '',
    subscription_key: '',
    api_user: '',
    encryption_key: '',
    webhook_secret: ''
  },
  priority: 50,
  is_test_mode: true,
  daily_limit: null,
  monthly_limit: null,
  notes: ''
})

const shouldShowField = (field: string) => {
  // Show all fields for now to be safe, or logic based on selected aggregator
  // For Lygos we need api_key and shop_name
  return true 
}

watch(() => props.instance, (newVal) => {
  if (newVal) {
    Object.assign(form, {
      aggregator_id: newVal.aggregator_id,
      instance_name: newVal.instance_name,
      // Note: Credentials are not returned in the list for security
      // We start with empty credentials when editing
      credentials: { 
        public_key: '',
        secret_key: '', 
        api_key: '',
        shop_name: '',
        client_id: '',
        client_secret: '',
        site_id: '',
        merchant_key: '',
        subscription_key: '',
        api_user: '',
        encryption_key: '',
        webhook_secret: ''
      },
      priority: newVal.priority,
      is_test_mode: newVal.is_test_mode,
      daily_limit: newVal.daily_limit,
      monthly_limit: newVal.monthly_limit,
      notes: newVal.notes
    })
  }
}, { immediate: true })

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
