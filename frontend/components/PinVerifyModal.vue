<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="showVerifyModal" 
        class="pin-overlay focus-trap"
        @click.self="handleClose"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/90 backdrop-blur-sm"></div>
        
        <!-- Modal Content -->
        <div class="pin-modal animate-in fade-in zoom-in duration-200" :class="{ 'shake': errorMessage }">
          <!-- Close Button -->
          <button 
            @click="handleClose"
            class="absolute top-4 right-4 p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>

          <!-- Header -->
          <div class="pin-icon mb-4 text-5xl">🔒</div>
          <h2 class="text-2xl font-bold mb-2 text-white">Vérification PIN</h2>
          <p class="text-sm mb-2 text-gray-400">
               Entrez votre code PIN pour continuer
               <br><span class="text-xs text-indigo-400 mt-1 block">🛡️ Clavier sécurisé anti-keylogger activé</span>
          </p>

          <!-- Error/Lock Message -->
          <div v-if="errorMessage" class="mb-4 p-3 rounded-lg bg-red-500/10 text-red-400 text-sm text-center border border-red-500/20">
            {{ errorMessage }}
            <span v-if="attemptsLeft !== null && attemptsLeft > 0" class="block mt-1">
              {{ attemptsLeft }} tentative(s) restante(s)
            </span>
          </div>

          <div v-if="isLocked" class="mb-4 p-3 rounded-lg bg-yellow-500/10 text-yellow-400 text-sm text-center border border-yellow-500/20">
            ⏳ PIN temporairement bloqué. Réessayez dans quelques minutes.
          </div>
          
          <!-- PIN Display -->
          <div class="pin-input-container flex gap-2 sm:gap-3 justify-center mb-6">
            <div
              v-for="(digit, i) in 5"
              :key="i"
              class="w-10 sm:w-12 h-12 sm:h-14 border-2 rounded-xl text-2xl flex items-center justify-center transition-all duration-200"
              :class="[
                  digit ? 'bg-indigo-500/10 border-indigo-500 text-white' : 'bg-white/5 border-white/10 text-gray-500',
                  currentLength === i ? 'ring-2 ring-indigo-500 ring-offset-2 ring-offset-[#1a1a2e] border-indigo-500' : ''
              ]"
            >
              <span class="transition-opacity duration-200" :class="digit ? 'opacity-100' : 'opacity-0'">•</span>
            </div>
          </div>

          <!-- Randomized Keypad -->
          <div class="grid grid-cols-3 gap-3 mb-6 max-w-[280px] mx-auto select-none">
              <button v-for="key in shuffledKeys" :key="key"
                      @click="handleKeyPress(key)"
                      :disabled="isLocked"
                      class="h-14 rounded-xl bg-white/5 hover:bg-white/10 active:bg-indigo-600 active:scale-95 transition-all text-xl font-bold text-white border border-white/10 flex items-center justify-center shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
              >
                  {{ key }}
              </button>
              <button @click="handleClear" :disabled="isLocked" class="h-14 rounded-xl bg-red-500/10 hover:bg-red-500/20 text-red-400 font-bold border border-red-500/20 flex items-center justify-center disabled:opacity-50">
                  C
              </button>
              <button @click="handleBackspace" :disabled="isLocked" class="h-14 rounded-xl bg-white/5 hover:bg-white/10 text-white font-bold border border-white/10 flex items-center justify-center disabled:opacity-50">
                  ⌫
              </button>
          </div>

          <!-- Submit Button -->
          <button 
            @click="handleSubmit"
            :disabled="currentLength < 5 || isLoading || isLocked"
            class="w-full p-4 rounded-xl border-none text-base font-semibold cursor-pointer mb-4 text-white transition-all hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed bg-gradient-to-br from-indigo-500 to-purple-600 shadow-lg shadow-indigo-500/20"
          >
            <span v-if="isLoading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              Vérification...
            </span>
            <span v-else>🔓 Vérifier</span>
          </button>

          <!-- Forgot PIN -->
          <p class="text-xs text-gray-500 text-center">
            <a href="#" class="text-indigo-400 hover:text-indigo-300">
              PIN oublié ?
            </a>
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { usePin } from '~/composables/usePin'

const { showVerifyModal, isLoading, verifyPin, closeModals, executePendingAction } = usePin()

const pinDigits = ref(['', '', '', '', ''])
const errorMessage = ref('')
const attemptsLeft = ref<number | null>(null)
const isLocked = ref(false)
const shuffledKeys = ref<number[]>([1, 2, 3, 4, 5, 6, 7, 8, 9, 0])

const currentLength = computed(() => pinDigits.value.filter(d => d !== '').length)

// Shuffle keys (Fisher-Yates)
const shuffleKeys = () => {
    const keys = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
    for (let i = keys.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [keys[i], keys[j]] = [keys[j], keys[i]];
    }
    shuffledKeys.value = keys
}

// Reset when modal opens
watch(() => showVerifyModal.value, (isOpen) => {
  if (isOpen) {
    pinDigits.value = ['', '', '', '', '']
    errorMessage.value = ''
    attemptsLeft.value = null
    isLocked.value = false
    shuffleKeys()
  }
})

const handleKeyPress = (key: number) => {
    if (currentLength.value < 5 && !isLocked.value) {
        pinDigits.value[currentLength.value] = key.toString()
        errorMessage.value = ''
    }
}

const handleBackspace = () => {
    if (currentLength.value > 0 && !isLocked.value) {
        pinDigits.value[currentLength.value - 1] = ''
        errorMessage.value = ''
    }
}

const handleClear = () => {
    if (!isLocked.value) {
        pinDigits.value = ['', '', '', '', '']
        errorMessage.value = ''
    }
}

const handleSubmit = async () => {
  const pin = pinDigits.value.join('')
  if (pin.length !== 5 || isLocked.value) return
  
  const result = await verifyPin(pin)
  
  if (result.success) {
    executePendingAction()
  } else {
    errorMessage.value = result.message
    
    if (result.attemptsLeft !== undefined) {
      attemptsLeft.value = result.attemptsLeft
    }
    
    if (result.locked) {
      isLocked.value = true
    }
    
    // Clear and reshuffle for security
    handleClear()
    shuffleKeys()
  }
}

const handleClose = () => {
  closeModals()
}
</script>

<style scoped>
.pin-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  background: rgba(0,0,0,0.9);
  backdrop-filter: blur(8px);
}

.pin-modal {
  position: relative;
  background: #1a1a2e;
  border-radius: 1.5rem;
  padding: 2rem;
  text-align: center;
  width: 100%;
  max-width: 24rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255,255,255,0.1);
}

.shake {
  animation: shake 0.5s cubic-bezier(.36,.07,.19,.97) both;
}

@keyframes shake {
  10%, 90% { transform: translate3d(-1px, 0, 0); }
  20%, 80% { transform: translate3d(2px, 0, 0); }
  30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
  40%, 60% { transform: translate3d(4px, 0, 0); }
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .pin-modal,
.modal-leave-to .pin-modal {
  transform: scale(0.9) translateY(20px);
}
</style>
