<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8 transition-colors duration-300">
    <div class="max-w-2xl mx-auto">
      
      <!-- Loading State -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        <p class="mt-4 text-gray-500 font-medium">Chargement de la page...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="!enterprise" class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8 text-center">
        <ExclamationTriangleIcon class="w-16 h-16 text-yellow-500 mx-auto mb-4" />
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Entreprise introuvable</h2>
        <p class="text-gray-500 dark:text-gray-400">Le lien que vous avez suivi semble invalide ou expiré.</p>
      </div>

      <!-- Main Content -->
      <div v-else class="bg-white dark:bg-gray-800 rounded-3xl shadow-xl overflow-hidden ring-1 ring-gray-900/5 transition-all duration-300 transform hover:shadow-2xl">
        
        <!-- Header Branding -->
        <div class="relative bg-gradient-to-br from-primary-600 to-primary-800 p-8 text-center text-white">
            <div class="absolute top-0 left-0 w-full h-full opacity-10 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')]"></div>
            
            <div class="relative z-10">
                <div v-if="enterprise.logo" class="w-24 h-24 mx-auto bg-white rounded-2xl p-1 mb-4 shadow-lg transform -rotate-3 hover:rotate-0 transition-transform duration-300">
                    <img :src="enterprise.logo" alt="Logo" class="w-full h-full object-cover rounded-xl">
                </div>
                <div v-else class="w-24 h-24 mx-auto bg-white/20 backdrop-blur-sm text-white rounded-2xl flex items-center justify-center text-4xl font-bold mb-4 shadow-lg border border-white/30">
                    {{ enterprise.name.charAt(0) }}
                </div>
                
                <h1 class="text-3xl font-extrabold tracking-tight mb-2">{{ enterprise.name }}</h1>
                <div class="inline-flex items-center gap-2 bg-white/20 backdrop-blur-md px-4 py-1.5 rounded-full text-sm font-medium border border-white/20 shadow-sm">
                    <CheckBadgeIcon class="w-4 h-4 text-green-300" />
                    <span>Page Officielle d'Abonnement</span>
                </div>
            </div>
        </div>

        <div class="p-8 space-y-8">
            <!-- Service Selection -->
            <div class="space-y-4">
                <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 flex items-center gap-2">
                    <SparklesIcon class="w-5 h-5 text-primary-500" />
                    Choisir un Service
                </label>
                
                <div class="relative">
                    <select v-model="selectedServiceId" :disabled="isServiceLocked" 
                        class="appearance-none block w-full pl-4 pr-10 py-3 text-base border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-xl transition-shadow shadow-sm hover:shadow-md">
                        <option value="">-- Sélectionner un service --</option>
                        <optgroup v-for="group in displayedGroups" :key="group.id" :label="group.name">
                            <option v-for="svc in group.services" :key="svc.id" :value="svc.id">
                                {{ svc.name }}
                            </option>
                        </optgroup>
                    </select>
                    <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-gray-500">
                        <ChevronUpDownIcon class="h-5 w-5" aria-hidden="true" />
                    </div>
                </div>
            </div>
            
            <div v-if="!displayedGroups.length && !isLoading" class="p-4 bg-yellow-50 text-yellow-700 rounded-lg text-sm">
                Aucun service public disponible pour le moment.
            </div>

            <!-- Selected Service Card -->
            <div v-if="selectedService" class="animate-fade-in-up">
                 <div class="bg-gray-50 dark:bg-gray-700/50 rounded-2xl p-6 border border-gray-100 dark:border-gray-600 relative overflow-hidden group">
                     <!-- ID Badge -->
                     <div class="absolute top-0 right-0 p-4 opacity-50 group-hover:opacity-100 transition-opacity">
                        <TagIcon class="w-12 h-12 text-gray-200 dark:text-gray-600 transform rotate-12" />
                     </div>
                     
                     <div class="relative z-10 grid grid-cols-1 sm:grid-cols-2 gap-6">
                         <div>
                            <p class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-1">Coût</p>
                            <div class="flex items-baseline gap-1">
                                <span class="text-3xl font-bold text-primary-600 dark:text-primary-400">{{ selectedService.base_price }}</span>
                                <span class="text-sm font-medium text-gray-500">{{ getServiceCurrency(selectedService) }}</span>
                            </div>
                         </div>
                         <div>
                            <p class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-1">Fréquence</p>
                            <div class="flex items-center gap-2 text-gray-700 dark:text-gray-200 font-medium">
                                <component :is="getFrequencyIcon(selectedService.billing_frequency)" class="w-5 h-5 text-gray-400" />
                                {{ formatFrequency(selectedService.billing_frequency) }}
                            </div>
                         </div>
                     </div>

                     <!-- Payment Schedule Table -->
                     <div v-if="selectedService.billing_frequency === 'CUSTOM' && selectedService.payment_schedule?.length" class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-600">
                        <div class="flex items-center gap-2 mb-3">
                            <CalendarDaysIcon class="w-5 h-5 text-primary-500" />
                            <h4 class="font-semibold text-sm text-gray-900 dark:text-white">Échéancier de paiement</h4>
                        </div>
                        <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-600 overflow-hidden">
                            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                                <thead class="bg-gray-50 dark:bg-gray-900">
                                    <tr>
                                        <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Libellé</th>
                                        <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Date</th>
                                        <th scope="col" class="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">Montant</th>
                                    </tr>
                                </thead>
                                <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                                    <tr v-for="(item, idx) in selectedService.payment_schedule" :key="idx" class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                        <td class="px-4 py-3 text-sm font-medium text-gray-900 dark:text-white">{{ item.name }}</td>
                                        <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ new Date(item.start_date || item.end_date).toLocaleDateString() }}</td>
                                        <td class="px-4 py-3 text-sm font-bold text-right text-primary-600">{{ item.amount }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                     </div>
                 </div>

                 <!-- Form Section -->
                 <div class="mt-8 space-y-6">
                     <div v-if="selectedService.form_schema?.length">
                         <div class="flex items-center gap-2 mb-4">
                             <DocumentTextIcon class="w-5 h-5 text-primary-500" />
                             <h3 class="text-lg font-bold text-gray-900 dark:text-white">Informations Complémentaires</h3>
                         </div>
                         
                         <div class="grid grid-cols-1 gap-5">
                            <div v-for="field in selectedService.form_schema" :key="field.key" class="group">
                                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 group-hover:text-primary-600 transition-colors">
                                    {{ field.label }} <span v-if="field.required" class="text-red-500">*</span>
                                </label>
                                
                                <div class="relative rounded-md shadow-sm">
                                    <input v-if="field.type !== 'select'" 
                                        v-model="formData[field.key]" 
                                        :type="field.type" 
                                        :required="field.required"
                                        :placeholder="`Entrez ${field.label.toLowerCase()}...`"
                                        class="block w-full rounded-xl border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-primary-500 focus:border-primary-500 py-3 px-4 border shadow-sm transition-all duration-200">
                                    
                                    <select v-else 
                                            v-model="formData[field.key]" 
                                            class="block w-full rounded-xl border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-primary-500 focus:border-primary-500 py-3 px-4 border shadow-sm transition-all duration-200">
                                            <option value="">Sélectionner une option...</option>
                                            <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                                    </select>
                                </div>
                            </div>
                         </div>
                     </div>

                     <!-- External ID -->
                     <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 flex items-center gap-2">
                             <IdentificationIcon class="w-4 h-4 text-gray-400" />
                             Matricule / ID Client
                             <span class="text-xs font-normal text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded-full ml-auto">Optionnel</span>
                        </label>
                        <div class="relative rounded-md shadow-sm">
                             <div class="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center">
                                 <HashtagIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                             </div>
                             <input v-model="externalId" placeholder="ex: MAT-2024-001" 
                                class="block w-full rounded-xl border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-primary-500 focus:border-primary-500 py-3 pl-10 px-4 border shadow-sm" />
                        </div>
                        <p class="text-xs text-gray-500 mt-2 flex items-center gap-1">
                            <InformationCircleIcon class="w-4 h-4 inline" />
                            Laissez ce champ vide si vous n'avez pas encore de matricule.
                        </p>
                     </div>

                     <!-- Submit Button -->
                     <button @click="submitSubscription" :disabled="isSubmitting" 
                        class="w-full relative flex justify-center py-4 px-4 border border-transparent rounded-xl shadow-lg text-base font-bold text-white bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-75 disabled:cursor-not-allowed transform hover:-translate-y-0.5 transition-all duration-200">
                         <span v-if="isSubmitting" class="absolute left-0 inset-y-0 flex items-center pl-3">
                            <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                         </span>
                         {{ isSubmitting ? 'Traitement en cours...' : "Valider l'inscription" }}
                     </button>
                 </div>
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
import { useAuthStore } from '@/stores/auth' // Ensure this store exists or adapt
import { 
    ExclamationTriangleIcon, CheckBadgeIcon, SparklesIcon, ChevronUpDownIcon,
    TagIcon, CalendarDaysIcon, ClockIcon, ArrowPathIcon, BanknotesIcon,
    DocumentTextIcon, IdentificationIcon, HashtagIcon, InformationCircleIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const router = useRouter()
const { authApi } = useApi() 
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

// Compute displayed groups based on visibility and URL params
const displayedGroups = computed(() => {
   if (!enterprise.value?.service_groups) return []
   return enterprise.value.service_groups.map(g => {
       // Filter services within group
       // Deep copy to avoid side effects
       const groupServices = g.services || []
       
       if (g.is_private) {
           // For private groups, only show if specific service matches URL param
           if (!serviceIdQuery) return null // Hide entire group if no direct link
           
           const matchedService = groupServices.find(s => s.id === serviceIdQuery)
           if (matchedService) {
               return { ...g, services: [matchedService] }
           }
           return null // No match in this private group
       }
       
       // Public group: show all
       return { ...g, services: groupServices }
   }).filter(g => g !== null && g.services.length > 0)
})

// Flatten for easier lookup of selected
const selectedService = computed(() => {
    if (!enterprise.value || !selectedServiceId.value) return null
    
    // Search in groups
    for (const group of enterprise.value.service_groups || []) {
        const found = group.services?.find(s => s.id === selectedServiceId.value)
        if (found) return found
    }
    return null
})

const getServiceCurrency = (svc) => {
    if (!enterprise.value || !svc) return 'XOF'
    const group = enterprise.value.service_groups.find(g => g.services.some(s => s.id === svc.id))
    return group ? group.currency : 'XOF'
}

const formatFrequency = (freq) => {
    const map = {
        'DAILY': 'Quotidien',
        'WEEKLY': 'Hebdomadaire',
        'MONTHLY': 'Mensuel',
        'ANNUALLY': 'Annuel',
        'CUSTOM': 'Personnalisé',
        'ONETIME': 'Paiement Unique'
    }
    return map[freq] || freq
}

const getFrequencyIcon = (freq) => {
    switch(freq) {
        case 'ONETIME': return BanknotesIcon
        case 'DAILY': return ClockIcon
        case 'WEEKLY': return ArrowPathIcon
        case 'MONTHLY': return CalendarDaysIcon
        default: return CalendarDaysIcon
    }
}

onMounted(async () => {
    // Check Auth
    if (!authStore.isAuthenticated) {
        // Redirect to login with return URL
        router.push(`/login?redirect=${encodeURIComponent(route.fullPath)}`)
        return
    }

    try {
        isLoading.value = true
        const { data } = await enterpriseAPI.get(enterpriseId)
        enterprise.value = data
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
                alert(`Le champ "${field.label}" est requis pour continuer.`)
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
        
        // Show success visual (e.g. checkmark modal) or redirect
        alert('🎉 Abonnement validé avec succès !')
        router.push('/dashboard') 
    } catch (e) {
        console.error(e)
        const msg = e.response?.data?.error || e.message
        alert("Oups ! Une erreur est survenue : " + msg)
    } finally {
        isSubmitting.value = false
    }
}
</script>

<style scoped>
.animate-fade-in-up {
    animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
