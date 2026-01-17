<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8">
        <!-- Loading State -->
        <div v-if="isLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p class="text-gray-500">Chargement de l'invitation...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-8">
          <XCircleIcon class="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Erreur</h2>
          <p class="text-gray-500 mb-6">{{ error }}</p>
          <NuxtLink to="/enterprise" 
            class="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
            Retour aux entreprises
          </NuxtLink>
        </div>

        <!-- Success State -->
        <div v-else-if="accepted" class="text-center py-8">
          <CheckCircleIcon class="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Invitation Acceptée !</h2>
          <p class="text-gray-500 mb-6">Vous faites maintenant partie de {{ invitation?.enterprise_name }}</p>
          <NuxtLink to="/enterprise" 
            class="inline-flex items-center px-6 py-3 bg-green-600 text-white rounded-xl hover:bg-green-700 font-medium">
            Accéder à l'entreprise
          </NuxtLink>
        </div>

        <!-- Invitation Details -->
        <div v-else-if="invitation">
          <div class="text-center mb-6">
            <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-600 to-purple-600 flex items-center justify-center mx-auto mb-4">
              <BuildingOffice2Icon class="w-10 h-10 text-white" />
            </div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Invitation d'Emploi</h1>
            <p class="text-gray-500">Vous avez été invité(e) à rejoindre une entreprise</p>
          </div>

          <!-- Enterprise Info Card -->
          <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-750 rounded-xl p-5 mb-6">
            <div class="flex items-center gap-4 mb-4">
              <div class="w-14 h-14 rounded-xl bg-gradient-to-br from-primary-500 to-purple-600 flex items-center justify-center">
                <BuildingOffice2Icon class="w-7 h-7 text-white" />
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ invitation.enterprise_name }}</h3>
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
                  {{ invitation.enterprise_type || 'Entreprise' }}
                </span>
              </div>
            </div>
            
            <div class="space-y-3 text-sm">
              <div class="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-600">
                <span class="text-gray-500">Poste proposé</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ invitation.profession || 'Employé' }}</span>
              </div>
              <div v-if="invitation.inviter_name" class="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-600">
                <span class="text-gray-500">Invité par</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ invitation.inviter_name }}</span>
              </div>
              <div class="flex items-center justify-between py-2">
                <span class="text-gray-500">Date d'invitation</span>
                <span class="font-medium text-gray-900 dark:text-white flex items-center gap-1.5">
                  <CalendarIcon class="w-4 h-4" />
                  {{ formatDate(invitation.invited_at) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Info Notice -->
          <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 mb-6">
            <div class="flex gap-3">
              <InformationCircleIcon class="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
              <div class="text-sm text-blue-700 dark:text-blue-300">
                <p class="font-medium mb-1">En acceptant cette invitation :</p>
                <ul class="list-disc list-inside text-blue-600 dark:text-blue-400 space-y-1">
                  <li>Vous rejoindrez l'entreprise en tant qu'employé</li>
                  <li>Vous pourrez accéder aux services de l'entreprise</li>
                  <li>Une confirmation par code PIN sera requise</li>
                </ul>
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
            <button @click="rejectInvitation" :disabled="isProcessing"
              class="flex-1 px-4 py-3 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 font-medium transition-colors disabled:opacity-50">
              Refuser
            </button>
            <button @click="handleAccept" :disabled="isProcessing"
              class="flex-1 px-4 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-xl hover:from-green-700 hover:to-green-800 font-medium shadow-lg shadow-green-500/25 disabled:opacity-50 flex items-center justify-center gap-2">
              <CheckIcon class="w-5 h-5" />
              Accepter l'invitation
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  BuildingOffice2Icon, CalendarIcon, CheckCircleIcon, XCircleIcon, 
  InformationCircleIcon, CheckIcon
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const route = useRoute()
const router = useRouter()
const { requirePin, checkPinStatus } = usePin()

const invitation = ref(null)
const isLoading = ref(true)
const isProcessing = ref(false)
const error = ref('')
const accepted = ref(false)

// Format date
const formatDate = (date) => date ? new Date(date).toLocaleDateString('fr-FR', {
  day: 'numeric',
  month: 'long',
  year: 'numeric'
}) : '--'

// Load invitation details
onMounted(async () => {
  const employeeId = route.query.id
  if (!employeeId) {
    error.value = 'ID d\'invitation manquant'
    isLoading.value = false
    return
  }

  try {
    // Check PIN status first
    await checkPinStatus()
    
    // TODO: Ideally fetch real invitation details from backend
    // For now, use query params or defaults
    invitation.value = {
      id: employeeId,
      enterprise_name: route.query.enterprise || 'Entreprise',
      enterprise_type: route.query.type || '',
      profession: route.query.profession || 'Employé',
      inviter_name: route.query.inviter || '',
      invited_at: route.query.date || new Date().toISOString()
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Invitation introuvable'
  } finally {
    isLoading.value = false
  }
})

// Handle accept with global PIN modal
const handleAccept = async () => {
  isProcessing.value = true
  
  try {
    // Use global secure PIN verification modal
    const verified = await requirePin(async () => {
      // This callback runs after PIN is verified
      await enterpriseAPI.acceptInvitation({
        employee_id: invitation.value.id
      })
      accepted.value = true
    })
    
    if (!verified) {
      // User cancelled PIN or PIN failed
      console.log('PIN verification cancelled or failed')
    }
  } catch (e) {
    const errMsg = e.response?.data?.details || e.response?.data?.error || 'Erreur lors de l\'acceptation'
    error.value = errMsg
  } finally {
    isProcessing.value = false
  }
}

// Reject invitation
const rejectInvitation = () => {
  router.push('/enterprise')
}
</script>
