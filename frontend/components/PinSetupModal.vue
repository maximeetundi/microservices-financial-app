<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="showSetupModal" 
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
      >
        <!-- Backdrop (no close on click - mandatory) -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-surface rounded-2xl shadow-2xl max-w-md w-full p-6">
          <!-- Header -->
          <div class="text-center mb-6">
            <div class="w-16 h-16 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mx-auto mb-4">
              <span class="text-3xl">🔐</span>
            </div>
            <h3 class="text-xl font-bold text-base">Définir votre PIN</h3>
            <p class="text-muted mt-2">
              Créez un code PIN à 5 chiffres pour sécuriser vos transactions
            </p>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="mb-4 p-3 rounded-lg bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm text-center">
            {{ errorMessage }}
          </div>
          
          <!-- PIN Input -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-muted mb-2">PIN (5 chiffres)</label>
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
                @input="(e) => handlePinInput(e, index)"
                @keydown="(e) => handlePinKeydown(e, index)"
                @paste="handlePaste"
                class="w-12 h-14 text-center text-2xl font-bold rounded-xl border-2 border-secondary-200 dark:border-secondary-700 bg-surface focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 outline-none transition-all"
              />
            </div>
          </div>

          <!-- Confirm PIN Input -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-muted mb-2">Confirmer le PIN</label>
            <div class="flex justify-center gap-2">
              <input 
                v-for="(_, index) in 5" 
                :key="'confirm-' + index"
                ref="confirmPinInputs"
                type="password"
                inputmode="numeric"
                pattern="[0-9]*"
                maxlength="1"
                :value="confirmPin[index] || ''"
                @input="(e) => handleConfirmPinInput(e, index)"
                @keydown="(e) => handleConfirmPinKeydown(e, index)"
                class="w-12 h-14 text-center text-2xl font-bold rounded-xl border-2 border-secondary-200 dark:border-secondary-700 bg-surface focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 outline-none transition-all"
              />
            </div>
          </div>

          <!-- Submit Button -->
          <button 
            @click="handleSubmit"
            :disabled="!canSubmit || isLoading"
            class="w-full py-3 px-4 rounded-xl font-medium bg-primary-600 hover:bg-primary-700 text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="isLoading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isLoading ? 'Définition en cours...' : 'Définir le PIN' }}</span>
          </button>
          
          <!-- Info -->
          <p class="text-xs text-muted text-center mt-4">
            Ce PIN sera requis pour toutes les transactions sensibles
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import { usePin } from '~/composables/usePin'

const { showSetupModal, isLoading, setupPin, executePendingAction } = usePin()

const pin = ref('')
const confirmPin = ref('')
const errorMessage = ref('')
const pinInputs = ref<HTMLInputElement[]>([])
const confirmPinInputs = ref<HTMLInputElement[]>([])

const canSubmit = computed(() => {
  return pin.value.length === 5 && confirmPin.value.length === 5
})

// Reset when modal opens
watch(() => showSetupModal.value, (isOpen) => {
  if (isOpen) {
    pin.value = ''
    confirmPin.value = ''
    errorMessage.value = ''
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
    } else {
      nextTick(() => confirmPinInputs.value[0]?.focus())
    }
  }
}

const handlePinKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'Backspace' && !pin.value[index] && index > 0) {
    pin.value = pin.value.slice(0, index - 1) + pin.value.slice(index)
    nextTick(() => pinInputs.value[index - 1]?.focus())
  }
}

const handleConfirmPinInput = (event: Event, index: number) => {
  const target = event.target as HTMLInputElement
  const value = target.value.replace(/\D/g, '')
  
  if (value) {
    confirmPin.value = confirmPin.value.slice(0, index) + value + confirmPin.value.slice(index + 1)
    if (index < 4) {
      nextTick(() => confirmPinInputs.value[index + 1]?.focus())
    }
  }
}

const handleConfirmPinKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'Backspace' && !confirmPin.value[index] && index > 0) {
    confirmPin.value = confirmPin.value.slice(0, index - 1) + confirmPin.value.slice(index)
    nextTick(() => confirmPinInputs.value[index - 1]?.focus())
  }
}

const handlePaste = async (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text').replace(/\D/g, '').slice(0, 5)
  if (pastedData) {
    pin.value = pastedData
    if (pastedData.length === 5) {
      nextTick(() => confirmPinInputs.value[0]?.focus())
    }
  }
}

const handleSubmit = async () => {
  errorMessage.value = ''
  
  if (pin.value !== confirmPin.value) {
    errorMessage.value = 'Les PINs ne correspondent pas'
    return
  }
  
  if (!/^\d{5}$/.test(pin.value)) {
    errorMessage.value = 'Le PIN doit contenir exactement 5 chiffres'
    return
  }
  
  const result = await setupPin(pin.value, confirmPin.value)
  
  if (result.success) {
    executePendingAction()
  } else {
    errorMessage.value = result.message
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
