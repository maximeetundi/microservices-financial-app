<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8">
        <!-- Loading State -->
        <div v-if="isLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p class="text-gray-500">Chargement des détails...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-8">
          <XCircleIcon class="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Erreur</h2>
          <p class="text-gray-500 mb-6">{{ error }}</p>
          <NuxtLink to="/" 
            class="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
            Retour à l'accueil
          </NuxtLink>
        </div>

        <!-- Success State -->
        <div v-else-if="completed" class="text-center py-8">
          <CheckCircleIcon class="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">{{ successTitle }}</h2>
          <p class="text-gray-500 mb-6">{{ successMessage }}</p>
          <NuxtLink :to="successRedirect" 
            class="inline-flex items-center px-6 py-3 bg-green-600 text-white rounded-xl hover:bg-green-700 font-medium">
            {{ successButtonText }}
          </NuxtLink>
        </div>

        <!-- Action Details -->
        <div v-else-if="action">
          <!-- Header with Icon -->
          <div class="text-center mb-6">
            <div :class="['w-20 h-20 rounded-2xl flex items-center justify-center mx-auto mb-4', actionIconBg]">
              <component :is="actionIcon" class="w-10 h-10 text-white" />
            </div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ actionTitle }}</h1>
            <p class="text-gray-500">{{ actionSubtitle }}</p>
          </div>

          <!-- Action Info Card -->
          <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-750 rounded-xl p-5 mb-6">
            <!-- Logo/Icon + Name -->
            <div class="flex items-center gap-4 mb-4">
              <img v-if="action.logo" 
                :src="action.logo" 
                :alt="action.name"
                class="w-14 h-14 rounded-xl object-cover" />
              <div v-else :class="['w-14 h-14 rounded-xl flex items-center justify-center', actionIconBg]">
                <component :is="actionIcon" class="w-7 h-7 text-white" />
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ action.name }}</h3>
                <span v-if="action.type_label" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
                  {{ action.type_label }}
                </span>
              </div>
            </div>
            
            <!-- Details List -->
            <div class="space-y-3 text-sm">
              <div v-for="detail in actionDetails" :key="detail.label" 
                class="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-600 last:border-0">
                <span class="text-gray-500">{{ detail.label }}</span>
                <span class="font-medium text-gray-900 dark:text-white flex items-center gap-1.5">
                  <component v-if="detail.icon" :is="detail.icon" class="w-4 h-4" />
                  {{ detail.value }}
                </span>
              </div>
            </div>
          </div>

          <!-- Description if present -->
          <div v-if="action.description" class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 mb-6">
            <div class="flex gap-3">
              <InformationCircleIcon class="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
              <div class="text-sm text-blue-700 dark:text-blue-300">
                <p class="font-medium mb-1">Détails de l'action</p>
                <p class="text-blue-600 dark:text-blue-400">{{ action.description }}</p>
              </div>
            </div>
          </div>

          <!-- Processing Message -->
          <div v-if="isProcessing" class="text-center py-4 mb-4">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto mb-2"></div>
            <p class="text-gray-500 text-sm">Traitement en cours...</p>
          </div>

          <!-- Actions -->
          <div class="flex gap-3">
            <button @click="handleReject" :disabled="isProcessing"
              class="flex-1 px-4 py-3 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 font-medium transition-colors disabled:opacity-50">
              {{ rejectButtonText }}
            </button>
            <button @click="handleApprove" :disabled="isProcessing"
              class="flex-1 px-4 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-xl hover:from-green-700 hover:to-green-800 font-medium shadow-lg shadow-green-500/25 disabled:opacity-50 flex items-center justify-center gap-2">
              <CheckIcon class="w-5 h-5" />
              {{ approveButtonText }}
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
import { 
  BuildingOffice2Icon, CalendarIcon, CheckCircleIcon, XCircleIcon, 
  InformationCircleIcon, CheckIcon, UserPlusIcon, ShieldCheckIcon,
  UserGroupIcon, CurrencyDollarIcon, DocumentCheckIcon
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const route = useRoute()
const router = useRouter()
const { requirePin, checkPinStatus } = usePin()

const action = ref(null)
const isLoading = ref(true)
const isProcessing = ref(false)
const error = ref('')
const completed = ref(false)

// Computed properties based on action type
const actionType = computed(() => action.value?.action_type || 'default')

const actionIcon = computed(() => {
  const icons = {
    'EMPLOYEE_INVITATION': UserPlusIcon,
    'PROMOTE_ADMIN': ShieldCheckIcon,
    'ADD_CLIENT': UserGroupIcon,
    'PAYMENT_APPROVAL': CurrencyDollarIcon,
    'default': DocumentCheckIcon
  }
  return icons[actionType.value] || icons.default
})

const actionIconBg = computed(() => {
  const backgrounds = {
    'EMPLOYEE_INVITATION': 'bg-gradient-to-br from-green-500 to-emerald-600',
    'PROMOTE_ADMIN': 'bg-gradient-to-br from-purple-500 to-indigo-600',
    'ADD_CLIENT': 'bg-gradient-to-br from-blue-500 to-cyan-600',
    'PAYMENT_APPROVAL': 'bg-gradient-to-br from-amber-500 to-orange-600',
    'default': 'bg-gradient-to-br from-primary-500 to-purple-600'
  }
  return backgrounds[actionType.value] || backgrounds.default
})

const actionTitle = computed(() => {
  const titles = {
    'EMPLOYEE_INVITATION': 'Invitation d\'Emploi',
    'PROMOTE_ADMIN': 'Promotion Administrateur',
    'ADD_CLIENT': 'Ajout de Client',
    'PAYMENT_APPROVAL': 'Approbation de Paiement',
    'default': 'Action à Approuver'
  }
  return titles[actionType.value] || titles.default
})

const actionSubtitle = computed(() => {
  const subtitles = {
    'EMPLOYEE_INVITATION': 'Vous avez été invité(e) à rejoindre une entreprise',
    'PROMOTE_ADMIN': 'Une demande de promotion nécessite votre approbation',
    'ADD_CLIENT': 'Confirmez l\'ajout de ce nouveau client',
    'PAYMENT_APPROVAL': 'Un paiement nécessite votre approbation',
    'default': 'Cette action nécessite votre confirmation'
  }
  return subtitles[actionType.value] || subtitles.default
})

const actionDetails = computed(() => {
  if (!action.value) return []
  
  const details = []
  
  // Common details
  if (action.value.role) {
    details.push({ label: 'Rôle', value: action.value.role })
  }
  if (action.value.profession) {
    details.push({ label: 'Poste', value: action.value.profession })
  }
  if (action.value.initiator_name) {
    details.push({ label: 'Initié par', value: action.value.initiator_name })
  }
  if (action.value.amount) {
    details.push({ label: 'Montant', value: `${action.value.amount} ${action.value.currency || ''}` })
  }
  if (action.value.created_at) {
    details.push({ 
      label: 'Date', 
      value: formatDate(action.value.created_at),
      icon: CalendarIcon
    })
  }
  
  return details
})

const approveButtonText = computed(() => {
  const texts = {
    'EMPLOYEE_INVITATION': 'Accepter l\'invitation',
    'PROMOTE_ADMIN': 'Approuver la promotion',
    'ADD_CLIENT': 'Confirmer l\'ajout',
    'PAYMENT_APPROVAL': 'Approuver le paiement',
    'default': 'Approuver'
  }
  return texts[actionType.value] || texts.default
})

const rejectButtonText = computed(() => {
  const texts = {
    'EMPLOYEE_INVITATION': 'Refuser',
    'default': 'Rejeter'
  }
  return texts[actionType.value] || texts.default
})

const successTitle = computed(() => {
  const titles = {
    'EMPLOYEE_INVITATION': 'Invitation Acceptée !',
    'PROMOTE_ADMIN': 'Promotion Approuvée !',
    'ADD_CLIENT': 'Client Ajouté !',
    'PAYMENT_APPROVAL': 'Paiement Approuvé !',
    'default': 'Action Approuvée !'
  }
  return titles[actionType.value] || titles.default
})

const successMessage = computed(() => {
  if (actionType.value === 'EMPLOYEE_INVITATION') {
    return `Vous faites maintenant partie de ${action.value?.name}`
  }
  return 'L\'action a été traitée avec succès'
})

const successRedirect = computed(() => {
  const redirects = {
    'EMPLOYEE_INVITATION': '/enterprise',
    'PROMOTE_ADMIN': '/enterprise',
    'default': '/'
  }
  return redirects[actionType.value] || redirects.default
})

const successButtonText = computed(() => {
  const texts = {
    'EMPLOYEE_INVITATION': 'Accéder à l\'entreprise',
    'default': 'Continuer'
  }
  return texts[actionType.value] || texts.default
})

// Format date
const formatDate = (date) => date ? new Date(date).toLocaleDateString('fr-FR', {
  day: 'numeric',
  month: 'long',
  year: 'numeric'
}) : '--'

// Load action details
onMounted(async () => {
  const actionId = route.params.id
  const actionTypeParam = route.query.type // e.g., 'invitation', 'approval'
  
  if (!actionId) {
    error.value = 'ID d\'action manquant'
    isLoading.value = false
    return
  }

  try {
    await checkPinStatus()
    
    // Determine which API to call based on type
    let data
    if (actionTypeParam === 'invitation' || route.path.includes('/accept')) {
      // Employee invitation
      const response = await enterpriseAPI.getInvitationDetails(actionId)
      data = {
        ...response.data,
        action_type: 'EMPLOYEE_INVITATION',
        name: response.data.enterprise_name,
        logo: response.data.enterprise_logo,
        type_label: response.data.enterprise_type
      }
    } else {
      // Multi-admin approval or other action types
      // TODO: Add endpoint for getting approval details
      // const response = await enterpriseAPI.getApprovalDetails(actionId)
      // For now, show error
      error.value = 'Type d\'action non supporté'
      isLoading.value = false
      return
    }
    
    action.value = data
  } catch (e) {
    error.value = e.response?.data?.details || e.response?.data?.error || 'Action introuvable ou déjà traitée'
  } finally {
    isLoading.value = false
  }
})

// Handle approve with global PIN modal
const handleApprove = async () => {
  isProcessing.value = true
  
  try {
    const verified = await requirePin(async () => {
      // Call appropriate API based on action type
      if (actionType.value === 'EMPLOYEE_INVITATION') {
        await enterpriseAPI.acceptInvitation({
          employee_id: action.value.employee_id || route.params.id
        })
      } else {
        // Multi-admin approval
        // await enterpriseAPI.approveAction(route.params.id, { pin: encryptedPin })
      }
      completed.value = true
    })
    
    if (!verified) {
      console.log('PIN verification cancelled or failed')
    }
  } catch (e) {
    error.value = e.response?.data?.details || e.response?.data?.error || 'Erreur lors de l\'approbation'
  } finally {
    isProcessing.value = false
  }
}

// Handle reject
const handleReject = () => {
  router.push(successRedirect.value)
}
</script>
