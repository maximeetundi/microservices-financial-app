<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="showVerifyModal" 
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        @click.self="handleClose"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-surface rounded-2xl shadow-2xl max-w-md w-full p-6">
          <!-- Close Button -->
          <button 
            @click="handleClose"
            class="absolute top-4 right-4 p-2 rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-800 transition-colors"
          >
            <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>

          <!-- Header -->
          <div class="text-center mb-6">
            <div class="w-16 h-16 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mx-auto mb-4">
              <span class="text-3xl">🔒</span>
            </div>
            <h3 class="text-xl font-bold text-base">Vérification PIN</h3>
            <p class="text-muted mt-2">
              Entrez votre code PIN pour continuer
            </p>
          </div>

          <!-- Error/Lock Message -->
          <div v-if="errorMessage" class="mb-4 p-3 rounded-lg bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm text-center">
            {{ errorMessage }}
            <span v-if="attemptsLeft !== null && attemptsLeft > 0" class="block mt-1">
              {{ attemptsLeft }} tentative(s) restante(s)
            </span>
          </div>

          <div v-if="isLocked" class="mb-4 p-3 rounded-lg bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400 text-sm text-center">
            ⏳ PIN temporairement bloqué. Réessayez dans quelques minutes.
          </div>
          
          <!-- PIN Input -->
          <div class="mb-6">
            <div class="flex justify-center gap-2">
              <input 
                v-for="(_, index) in 5" 
                :key="index"
                ref="pinInputs"
                type="password"
                inputmode="numeric"
                pattern="[0-9]*"
                maxlength="1"
                :value="pin[index] || ''"
                :disabled="isLocked"
                @input="(e) => handlePinInput(e, index)"
                @keydown="(e) => handlePinKeydown(e, index)"
                @paste="handlePaste"
                class="w-12 h-14 text-center text-2xl font-bold rounded-xl border-2 border-secondary-200 dark:border-secondary-700 bg-surface focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              />
            </div>
          </div>

          <!-- Submit Button -->
          <button 
            @click="handleSubmit"
            :disabled="pin.length !== 5 || isLoading || isLocked"
            class="w-full py-3 px-4 rounded-xl font-medium bg-primary-600 hover:bg-primary-700 text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="isLoading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isLoading ? 'Vérification...' : 'Vérifier' }}</span>
          </button>

          <!-- Forgot PIN -->
          <p class="text-xs text-muted text-center mt-4">
            <a href="#" class="text-primary-600 hover:text-primary-700 dark:text-primary-400">
              PIN oublié ?
            </a>
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { usePin } from '~/composables/usePin'

const { showVerifyModal, isLoading, verifyPin, closeModals, executePendingAction } = usePin()

const pin = ref('')
const errorMessage = ref('')
const attemptsLeft = ref<number | null>(null)
const isLocked = ref(false)
const pinInputs = ref<HTMLInputElement[]>([])

// Reset when modal opens
watch(() => showVerifyModal.value, (isOpen) => {
  if (isOpen) {
    pin.value = ''
    errorMessage.value = ''
    attemptsLeft.value = null
    isLocked.value = false
    nextTick(() => {
      pinInputs.value[0]?.focus()
    })
  }
})

const handlePinInput = (event: Event, index: number) => {
  const target = event.target as HTMLInputElement
  const value = target.value.replace(/\D/g, '')
  
  if (value) {
    pin.value = pin.value.slice(0, index) + value + pin.value.slice(index + 1)
    if (index < 4) {
      nextTick(() => pinInputs.value[index + 1]?.focus())
    }
  }
}

const handlePinKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'Backspace' && !pin.value[index] && index > 0) {
    pin.value = pin.value.slice(0, index - 1) + pin.value.slice(index)
    nextTick(() => pinInputs.value[index - 1]?.focus())
  }
  
  // Submit on Enter when PIN is complete
  if (event.key === 'Enter' && pin.value.length === 5) {
    handleSubmit()
  }
}

const handlePaste = async (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text').replace(/\D/g, '').slice(0, 5)
  if (pastedData) {
    pin.value = pastedData
  }
}

const handleClose = () => {
  closeModals()
}

const handleSubmit = async () => {
  errorMessage.value = ''
  
  const result = await verifyPin(pin.value)
  
  if (result.valid) {
    executePendingAction()
  } else {
    errorMessage.value = result.message || 'PIN incorrect'
    attemptsLeft.value = result.attemptsLeft ?? null
    
    // Check if locked
    if (attemptsLeft.value === 0) {
      isLocked.value = true
    }
    
    // Clear PIN and refocus
    pin.value = ''
    nextTick(() => {
      pinInputs.value[0]?.focus()
    })
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.9) translateY(20px);
}
</style>
