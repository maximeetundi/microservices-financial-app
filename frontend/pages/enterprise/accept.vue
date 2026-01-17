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
          <div class="text-center mb-8">
            <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-600 to-purple-600 flex items-center justify-center mx-auto mb-4">
              <BuildingOffice2Icon class="w-10 h-10 text-white" />
            </div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Invitation d'Emploi</h1>
            <p class="text-gray-500">Vous êtes invité(e) à rejoindre</p>
          </div>

          <!-- Enterprise Info -->
          <div class="bg-gray-50 dark:bg-gray-700 rounded-xl p-4 mb-6">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-1">{{ invitation.enterprise_name }}</h3>
            <p class="text-sm text-gray-500">{{ invitation.profession || 'Employé' }}</p>
            <div class="flex items-center gap-2 mt-3 text-sm text-gray-500">
              <CalendarIcon class="w-4 h-4" />
              <span>Invité le {{ formatDate(invitation.invited_at) }}</span>
            </div>
          </div>

          <!-- PIN Input -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Code PIN de sécurité
            </label>
            <div class="flex justify-center gap-3">
              <input
                v-for="(_, index) in 5"
                :key="index"
                ref="pinInputs"
                type="password"
                maxlength="1"
                inputmode="numeric"
                class="w-12 h-14 text-center text-2xl font-bold rounded-xl border-2 border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all"
                @input="handlePinInput($event, index)"
                @keydown.backspace="handleBackspace($event, index)"
              />
            </div>
            <p v-if="pinError" class="text-red-500 text-sm text-center mt-2">{{ pinError }}</p>
          </div>

          <!-- Actions -->
          <div class="flex gap-3">
            <button @click="rejectInvitation" :disabled="isProcessing"
              class="flex-1 px-4 py-3 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 font-medium transition-colors">
              Refuser
            </button>
            <button @click="acceptInvitation" :disabled="isProcessing || pin.length !== 5"
              class="flex-1 px-4 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-xl hover:from-green-700 hover:to-green-800 font-medium shadow-lg shadow-green-500/25 disabled:opacity-50 flex items-center justify-center gap-2">
              <div v-if="isProcessing" class="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
              {{ isProcessing ? '' : 'Accepter' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  BuildingOffice2Icon, CalendarIcon, CheckCircleIcon, XCircleIcon 
} from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'
import { usePin } from '@/composables/usePin'

const route = useRoute()
const router = useRouter()
const { verifyPin, checkPinStatus, hasPin } = usePin()

const invitation = ref(null)
const isLoading = ref(true)
const isProcessing = ref(false)
const error = ref('')
const pinError = ref('')
const accepted = ref(false)
const pinInputs = ref([])

// PIN handling
const pin = computed(() => {
  return pinInputs.value.map(input => input?.value || '').join('')
})

const handlePinInput = (event, index) => {
  const value = event.target.value
  if (value && index < 4) {
    pinInputs.value[index + 1]?.focus()
  }
  pinError.value = ''
}

const handleBackspace = (event, index) => {
  if (!event.target.value && index > 0) {
    pinInputs.value[index - 1]?.focus()
  }
}

// Format date
const formatDate = (date) => date ? new Date(date).toLocaleDateString('fr-FR') : '--'

// Load invitation
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
    
    // Fetch employee/invitation details
    invitation.value = {
      id: employeeId,
      enterprise_name: route.query.enterprise || 'Entreprise',
      profession: route.query.profession || '',
      invited_at: new Date().toISOString()
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Invitation introuvable'
  } finally {
    isLoading.value = false
  }
})

// Accept invitation with global PIN verification
const acceptInvitation = async () => {
  if (pin.value.length !== 5) {
    pinError.value = 'Veuillez entrer votre code PIN à 5 chiffres'
    return
  }

  isProcessing.value = true
  pinError.value = ''

  try {
    // First verify PIN using global verification system
    const pinResult = await verifyPin(pin.value)
    if (!pinResult.valid) {
      pinError.value = pinResult.message || 'PIN incorrect'
      if (pinResult.attemptsLeft !== undefined) {
        pinError.value += ` (${pinResult.attemptsLeft} essais restants)`
      }
      isProcessing.value = false
      return
    }

    // PIN verified, now accept invitation
    await enterpriseAPI.acceptInvitation({
      employee_id: invitation.value.id,
      pin: pin.value
    })
    accepted.value = true
  } catch (e) {
    const errMsg = e.response?.data?.details || e.response?.data?.error || 'Erreur lors de l\'acceptation'
    if (errMsg.toLowerCase().includes('pin')) {
      pinError.value = errMsg
    } else {
      error.value = errMsg
    }
  } finally {
    isProcessing.value = false
  }
}

// Reject invitation (optional - could navigate away)
const rejectInvitation = () => {
  router.push('/enterprise')
}
</script>
