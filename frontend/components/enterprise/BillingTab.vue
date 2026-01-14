<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <DocumentTextIcon class="w-6 h-6 text-blue-500" />
          Facturation en Masse
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Importez ou saisissez les données de consommation
        </p>
      </div>
    </div>

    <!-- Service Selection -->
    <div class="p-4 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Service Concerné</label>
      <select v-model="selectedService" @change="loadSubscribers"
        class="w-full md:w-1/2 px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500">
        <option value="">-- Sélectionner un service --</option>
        <optgroup v-for="group in enterprise.service_groups" :key="group.id" :label="group.name">
          <option v-for="svc in group.services" :key="svc.id" :value="svc.id">
            {{ svc.name }} ({{ svc.billing_type === 'USAGE' ? 'Compteur' : 'Fixe' }})
          </option>
        </optgroup>
      </select>
    </div>

    <!-- Mode Toggle -->
    <div v-if="selectedService" class="flex gap-3">
      <button @click="billingMode = 'IMPORT'" 
        :class="['px-5 py-2.5 rounded-xl text-sm font-medium border flex items-center gap-2 transition-all', 
          billingMode === 'IMPORT' ? 'bg-blue-100 text-blue-700 border-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-700' : 'bg-gray-100 text-gray-600 border-gray-200 dark:bg-gray-700 dark:text-gray-400']">
        <DocumentArrowUpIcon class="w-5 h-5" />
        Importer CSV
      </button>
      <button @click="billingMode = 'MANUAL'; loadSubscribers()" 
        :class="['px-5 py-2.5 rounded-xl text-sm font-medium border flex items-center gap-2 transition-all', 
          billingMode === 'MANUAL' ? 'bg-purple-100 text-purple-700 border-purple-200 dark:bg-purple-900/30 dark:text-purple-300 dark:border-purple-700' : 'bg-gray-100 text-gray-600 border-gray-200 dark:bg-gray-700 dark:text-gray-400']">
        <PencilSquareIcon class="w-5 h-5" />
        Saisie Manuelle
      </button>
    </div>

    <!-- CSV Import Mode -->
    <div v-if="billingMode === 'IMPORT' && selectedService" 
      class="p-8 bg-white dark:bg-gray-800 rounded-2xl border-2 border-dashed border-gray-200 dark:border-gray-600">
      <div class="text-center">
        <input type="file" @change="handleFileUpload" accept=".csv,.xlsx" class="hidden" id="file-upload" />
        <label for="file-upload" class="cursor-pointer">
          <div class="w-16 h-16 mx-auto mb-4 bg-blue-100 dark:bg-blue-900/30 rounded-2xl flex items-center justify-center">
            <DocumentArrowUpIcon class="w-8 h-8 text-blue-500" />
          </div>
          <p class="text-gray-700 dark:text-gray-300 font-medium mb-1">Cliquez pour sélectionner un fichier</p>
          <p class="text-sm text-gray-500">CSV ou Excel - Colonnes: ID Client, Consommation/Montant</p>
          <div v-if="importFile" class="mt-4 inline-flex items-center gap-2 px-4 py-2 bg-green-100 text-green-700 rounded-lg">
            <CheckIcon class="w-5 h-5" />
            {{ importFile.name }}
          </div>
        </label>
      </div>
      <div v-if="importFile" class="flex justify-center mt-6">
        <button @click="uploadFile" :disabled="isUploading"
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium disabled:opacity-50 flex items-center gap-2">
          {{ isUploading ? 'Import en cours...' : 'Importer et Générer' }}
        </button>
      </div>
    </div>

    <!-- Manual Entry Mode -->
    <div v-if="billingMode === 'MANUAL' && selectedService" class="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div v-if="isLoadingSubscribers" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
      </div>

      <div v-else-if="subscribers.length === 0" class="text-center py-12 text-gray-500">
        <UsersIcon class="w-12 h-12 mx-auto mb-3 opacity-50" />
        <p>Aucun abonné éligible pour ce service</p>
      </div>

      <div v-else>
        <div class="overflow-x-auto">
          <table class="w-full text-sm text-left">
            <thead class="bg-gray-50 dark:bg-gray-900/50 text-xs text-gray-600 dark:text-gray-400 uppercase">
              <tr>
                <th class="px-4 py-3">Client</th>
                <th class="px-4 py-3">ID Abonnement</th>
                <th class="px-4 py-3 text-center">Consommation</th>
                <th class="px-4 py-3 text-center">Montant</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
              <tr v-for="sub in subscribers" :key="sub.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                <td class="px-4 py-3 font-medium text-gray-900 dark:text-white">{{ sub.client_name || 'Inconnu' }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 font-mono">{{ sub.id.substring(0, 8) }}...</td>
                <td class="px-4 py-3 text-center">
                  <input v-model.number="entries[sub.id].consumption" type="number" step="0.1"
                    class="w-24 px-2 py-1.5 text-right border rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700" />
                </td>
                <td class="px-4 py-3 text-center">
                  <input v-model.number="entries[sub.id].amount" type="number" step="50"
                    class="w-28 px-2 py-1.5 text-right border rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="p-4 border-t border-gray-100 dark:border-gray-700 flex justify-between items-center">
          <span class="text-sm text-gray-500">{{ subscribers.length }} abonnés</span>
          <div class="flex gap-3">
            <button @click="loadSubscribers" class="px-4 py-2 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">
              Actualiser
            </button>
            <button @click="submitManualBatch" :disabled="isGenerating"
              class="px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium disabled:opacity-50">
              {{ isGenerating ? 'Génération...' : 'Valider et Générer' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { DocumentTextIcon, DocumentArrowUpIcon, PencilSquareIcon, CheckIcon, UsersIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const props = defineProps({
  enterprise: { type: Object, required: true }
})

const selectedService = ref('')
const billingMode = ref('IMPORT')
const importFile = ref(null)
const isUploading = ref(false)

const subscribers = ref([])
const entries = reactive({})
const isLoadingSubscribers = ref(false)
const isGenerating = ref(false)

const handleFileUpload = (event) => {
  importFile.value = event.target.files[0]
}

const uploadFile = async () => {
  if (!importFile.value) return
  isUploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    formData.append('service_id', selectedService.value)
    formData.append('col_client_idx', '0')
    formData.append('col_amount_idx', '1')
    formData.append('col_consumption_idx', '2')

    await enterpriseAPI.importInvoices(props.enterprise.id, formData)
    alert('Import réussi ! Factures générées en brouillon.')
    importFile.value = null
  } catch (e) {
    alert('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    isUploading.value = false
  }
}

const loadSubscribers = async () => {
  if (!selectedService.value) return
  isLoadingSubscribers.value = true
  try {
    const { data } = await enterpriseAPI.getSubscriptions(props.enterprise.id, selectedService.value)
    subscribers.value = data || []
    
    // Initialize entries
    Object.keys(entries).forEach(k => delete entries[k])
    data.forEach(sub => {
      entries[sub.id] = { amount: 0, consumption: 0 }
    })
  } catch (e) {
    console.error('Failed to load subscribers', e)
  } finally {
    isLoadingSubscribers.value = false
  }
}

const submitManualBatch = async () => {
  if (!confirm('Générer les factures ?')) return
  isGenerating.value = true
  try {
    const items = Object.entries(entries)
      .map(([subId, val]) => ({
        subscription_id: subId,
        amount: parseFloat(val.amount) || 0,
        consumption: parseFloat(val.consumption) || 0
      }))
      .filter(i => i.amount > 0 || i.consumption > 0)

    if (items.length === 0) {
      alert('Aucune donnée saisie')
      return
    }

    await enterpriseAPI.createBatchInvoices(props.enterprise.id, items)
    alert('Factures générées avec succès')
  } catch (e) {
    alert('Erreur: ' + e.message)
  } finally {
    isGenerating.value = false
  }
}
</script>
