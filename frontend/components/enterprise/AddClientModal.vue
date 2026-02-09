<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <UserPlusIcon class="w-6 h-6 text-indigo-500" />
          Ajouter un Client / Élève
        </h3>

        <!-- Step 1: Search User -->
        <div class="mb-5">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Rechercher un utilisateur</label>
          <div class="flex gap-2">
            <input v-model="searchQuery" placeholder="Téléphone ou Email"
              class="flex-1 px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
            <button @click="searchUser" :disabled="isSearching" type="button"
              class="px-4 py-2.5 bg-primary-100 text-primary-700 rounded-xl font-medium hover:bg-primary-200 disabled:opacity-50">
              {{ isSearching ? '...' : '🔍' }}
            </button>
          </div>
          <div v-if="foundUser" class="mt-2 p-3 bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 rounded-xl flex justify-between items-center">
            <span class="flex items-center gap-2">
              <CheckCircleIcon class="w-5 h-5" />
              {{ foundUser.name }}
            </span>
            <button @click="foundUser = null; searchQuery = ''" class="text-gray-500 hover:text-gray-700">×</button>
          </div>
        </div>

        <!-- Step 2: Select Service -->
        <div v-if="foundUser" class="mb-5">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Service ou Classe</label>
          <select v-model="form.service_id"
            class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500">
            <option value="">-- Sélectionner --</option>
            <optgroup v-for="group in (enterprise?.service_groups || [])" :key="group.id" :label="group.name">
              <option v-for="svc in group.services" :key="svc.id" :value="svc.id">{{ svc.name }}</option>
            </optgroup>
          </select>
        </div>

        <!-- External ID -->
        <div v-if="foundUser" class="mb-5">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Matricule / ID Externe (Optionnel)</label>
          <input v-model="form.external_id" placeholder="ex: MAT-2024-001"
            class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500" />
        </div>

        <!-- Dynamic Form Fields -->
        <div v-if="foundUser && selectedServiceSchema?.length" class="mb-5 p-4 bg-gray-50 dark:bg-gray-700 rounded-xl">
          <h4 class="text-xs font-bold text-gray-500 uppercase mb-3">Informations Complémentaires</h4>
          <div class="space-y-3">
            <div v-for="field in selectedServiceSchema" :key="field.key">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                {{ field.label }} <span v-if="field.required" class="text-red-500">*</span>
              </label>
              <input v-if="field.type !== 'select'" v-model="form.form_data[field.key]" :type="field.type" :required="field.required"
                class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-600 text-sm" />
              <select v-else v-model="form.form_data[field.key]"
                class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-600 text-sm">
                <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button type="button" @click="$emit('close')"
            class="flex-1 px-4 py-2.5 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700">
            Annuler
          </button>
          <button @click="addClient" :disabled="!foundUser || !form.service_id || isAdding"
            class="flex-1 px-4 py-2.5 bg-gradient-to-r from-purple-600 to-purple-700 text-white rounded-xl font-medium hover:from-purple-700 hover:to-purple-800 disabled:opacity-50">
            {{ isAdding ? 'Ajout...' : 'Ajouter' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { UserPlusIcon, XMarkIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI, authAPI } from '@/composables/useApi'

const props = defineProps({
  enterprise: { type: Object, required: false, default: null },
  enterpriseId: { type: String, required: true }
})

const emit = defineEmits(['close', 'added'])

const searchQuery = ref('')
const foundUser = ref(null)
const isSearching = ref(false)
const isAdding = ref(false)

const form = reactive({
  service_id: '',
  external_id: '',
  form_data: {}
})

const allServices = computed(() => {
  return (props.enterprise?.service_groups || []).flatMap(g => g.services || [])
})

const selectedServiceSchema = computed(() => {
  if (!form.service_id) return []
  const svc = allServices.value.find(s => s.id === form.service_id)
  return svc?.form_schema || []
})

const searchUser = async () => {
  if (!searchQuery.value) return
  isSearching.value = true
  foundUser.value = null
  try {
    // Detect if input is email or phone
    const isEmail = searchQuery.value.includes('@')
    const params = isEmail 
      ? { email: searchQuery.value } 
      : { phone: searchQuery.value }
    
    const { data } = await authAPI.lookup(params)
    if (data) {
      foundUser.value = {
        id: data.id,
        name: `${data.first_name || ''} ${data.last_name || ''}`.trim() || data.email || data.phone_number
      }
    } else {
      alert('Utilisateur non trouvé')
    }
  } catch (e) {
    console.error('User lookup failed:', e)
    alert('Utilisateur non trouvé')
  } finally {
    isSearching.value = false
  }
}

const addClient = async () => {
  if (!foundUser.value || !form.service_id) return
  isAdding.value = true
  try {
    const svc = allServices.value.find(s => s.id === form.service_id)
    const payload = {
      client_id: foundUser.value.id,
      client_name: foundUser.value.name,
      service_id: form.service_id,
      external_id: form.external_id,
      form_data: form.form_data,
      amount: svc?.base_price || 0,
      billing_frequency: svc?.billing_frequency || 'MONTHLY'
    }
    if (!props.enterpriseId) throw new Error('enterpriseId manquant')
    await enterpriseAPI.createSubscription(props.enterpriseId, payload)
    emit('added')
  } catch (e) {
    alert('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    isAdding.value = false
  }
}
</script>
