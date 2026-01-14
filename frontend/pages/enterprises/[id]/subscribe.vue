<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <!-- Loading State -->
      <div v-if="isLoading" class="p-8 text-center text-gray-500">
        Chargement...
      </div>

      <div v-else-if="!enterprise" class="p-8 text-center text-red-500">
        Entreprise introuvable.
      </div>

      <!-- Content -->
      <div v-else>
        <!-- Header / Branding -->
        <div class="bg-primary-600 p-6 text-center">
            <div v-if="enterprise.logo" class="w-20 h-20 mx-auto bg-white rounded-full p-1 mb-3">
                <img :src="enterprise.logo" alt="Logo" class="w-full h-full object-cover rounded-full">
            </div>
            <div v-else class="w-20 h-20 mx-auto bg-white/20 text-white rounded-full flex items-center justify-center text-3xl font-bold mb-3">
                {{ enterprise.name.charAt(0) }}
            </div>
            <h1 class="text-2xl font-bold text-white">{{ enterprise.name }}</h1>
            <p class="text-primary-100 text-sm mt-1">S'abonner à un service</p>
        </div>

        <div class="p-6 space-y-6">
            <!-- Service Selection -->
            <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Service souhaité</label>
                <select v-model="selectedServiceId" :disabled="isServiceLocked" class="block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500 py-2 px-3 border">
                    <option value="">Sélectionner un service...</option>
                    <option v-for="svc in enterprise.custom_services" :key="svc.id" :value="svc.id">
                        {{ svc.name }} - {{ svc.base_price }} {{ enterprise.settings?.currency }}
                    </option>
                </select>
            </div>

            <div v-if="selectedService" class="space-y-4">
                 <div class="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg text-sm text-blue-800 dark:text-blue-200">
                     <p><strong>Fréquence :</strong> {{ formatFrequency(selectedService.billing_frequency) }}</p>
                     <p><strong>Montant :</strong> {{ selectedService.base_price }} {{ enterprise.settings?.currency }}</p>
                     
                     <!-- Payment Schedule Display -->
                     <div v-if="selectedService.billing_frequency === 'CUSTOM' && selectedService.payment_schedule?.length" class="mt-3 bg-white dark:bg-gray-800 rounded p-2 border border-blue-100 dark:border-blue-800">
                        <p class="font-semibold text-xs uppercase mb-2 text-gray-500">Calendrier de paiement</p>
                        <ul class="space-y-1">
                            <li v-for="(item, idx) in selectedService.payment_schedule" :key="idx" class="flex justify-between text-xs">
                                <span>{{ item.name }} ({{ new Date(item.start_date || item.end_date).toLocaleDateString() }})</span>
                                <span class="font-medium">{{ item.amount }} {{ enterprise.settings?.currency }}</span>
                            </li>
                        </ul>
                     </div>
                 </div>

                 <!-- Dynamic Form Fields -->
                 <div v-if="selectedService.form_schema?.length">
                     <h3 class="font-medium text-gray-900 dark:text-white border-b pb-2">Informations requises</h3>
                     <div class="space-y-3 pt-2">
                        <div v-for="field in selectedService.form_schema" :key="field.key">
                            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                {{ field.label }} <span v-if="field.required" class="text-red-500">*</span>
                            </label>
                            
                            <input v-if="field.type !== 'select'" 
                                   v-model="formData[field.key]" 
                                   :type="field.type" 
                                   :required="field.required"
                                   class="block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white shadow-sm py-2 px-3 border">
                                   
                            <select v-else 
                                    v-model="formData[field.key]" 
                                    class="block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white shadow-sm py-2 px-3 border">
                                    <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                            </select>
                        </div>
                     </div>
                 </div>

                 <!-- Matricule / External ID -->
                 <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Matricule / ID Client (Optionnel)
                    </label>
                    <input v-model="externalId" placeholder="ex: MAT-001" class="block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white shadow-sm py-2 px-3 border">
                    <p class="text-xs text-gray-400 mt-1">Laissez vide si on vous l'attribue automatiquement.</p>
                 </div>

                 <!-- Submit Button -->
                 <button @click="submitSubscription" :disabled="isSubmitting" class="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50">
                     {{ isSubmitting ? 'Traitement...' : "S'abonner" }}
                 </button>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { enterpriseAPI, useApi } from '@/composables/useApi'

const route = useRoute()
const router = useRouter()
const { authApi } = useApi() // Access auth store effectively via API or directly?
// Assuming we have an auth store or way to get user info.
// In this codebase, useApi mainly returns APIs. We might need useAuthStore if available.
import { useAuthStore } from '@/stores/auth' // Guessing store path based on context

const authStore = useAuthStore()

const enterpriseId = route.params.id
const serviceIdQuery = route.query.service_id

const enterprise = ref(null)
const isLoading = ref(true)
const isSubmitting = ref(false)

const selectedServiceId = ref(serviceIdQuery || '')
const isServiceLocked = ref(!!serviceIdQuery) // Lock if provided in URL

const formData = ref({})
const externalId = ref('')

const selectedService = computed(() => {
    if (!enterprise.value || !selectedServiceId.value) return null
    return enterprise.value.custom_services.find(s => s.id === selectedServiceId.value)
})

const formatFrequency = (freq) => {
    const map = {
        'DAILY': 'Quotidien',
        'WEEKLY': 'Hebdomadaire',
        'MONTHLY': 'Mensuel',
        'ANNUALLY': 'Annuel',
        'CUSTOM': 'Personnalisé',
        'ONETIME': 'Ponctuel'
    }
    return map[freq] || freq
}

onMounted(async () => {
    // Check Auth
    if (!authStore.isAuthenticated) {
        // Redirect to login with return URL
        // Simple redirect for now
        router.push(`/login?redirect=${encodeURIComponent(route.fullPath)}`)
        return
    }

    try {
        isLoading.value = true
        const { data } = await enterpriseAPI.get(enterpriseId)
        enterprise.value = data
        
        // Initialize formData defaults if needed
    } catch (e) {
        console.error('Failed to load enterprise', e)
    } finally {
        isLoading.value = false
    }
})

const submitSubscription = async () => {
    if (!selectedService.value) return
    
    // Validate required fields
    if (selectedService.value.form_schema) {
        for (const field of selectedService.value.form_schema) {
            if (field.required && !formData.value[field.key]) {
                alert(`Le champ "${field.label}" est requis.`)
                return
            }
        }
    }

    isSubmitting.value = true
    try {
        const payload = {
            client_id: authStore.user?.id,
            client_name: authStore.user?.name,
            service_id: selectedService.value.id,
            billing_frequency: selectedService.value.billing_frequency,
            amount: selectedService.value.base_price,
            form_data: formData.value,
            external_id: externalId.value
        }

        await enterpriseAPI.createSubscription(enterpriseId, payload)
        
        alert('Abonnement réussi !')
        router.push('/dashboard') // Or User's Subscription List
    } catch (e) {
        console.error(e)
        alert("Erreur lors de l'abonnement: " + (e.response?.data?.error || e.message))
    } finally {
        isSubmitting.value = false
    }
}
</script>
