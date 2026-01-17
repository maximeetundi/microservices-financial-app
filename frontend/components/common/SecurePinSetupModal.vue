<template>
  <Teleport to="body">
    <div v-if="isOpen" class="pin-overlay focus-trap" @click.self="handleClose">
      <div class="pin-modal animate-in fade-in zoom-in duration-200" :class="{ 'shake': error }">
        <div class="pin-icon mb-4 text-5xl">🔐</div>
        <h2 class="text-2xl font-bold mb-2 text-white">{{ title }}</h2>
        <p class="text-sm mb-2 text-gray-400">
             {{ step === 1 ? 'Créez un code PIN à 5 chiffres' : 'Confirmez votre code PIN' }}
             <br><span class="text-xs text-indigo-400 mt-1 block">🛡️ Clavier sécurisé anti-keylogger activé</span>
        </p>
        
        <!-- Step Indicator -->
        <div class="flex justify-center gap-2 mb-4">
          <div class="w-3 h-3 rounded-full" :class="step === 1 ? 'bg-indigo-500' : 'bg-gray-600'"></div>
          <div class="w-3 h-3 rounded-full" :class="step === 2 ? 'bg-indigo-500' : 'bg-gray-600'"></div>
        </div>
        
        <!-- Display Only Inputs -->
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

        <p v-if="error" class="text-red-500 text-sm mb-4 font-medium bg-red-500/10 p-2 rounded-lg border border-red-500/20">{{ error }}</p>
        
        <!-- Randomized Keypad -->
        <div class="grid grid-cols-3 gap-3 mb-6 max-w-[280px] mx-auto select-none">
            <button v-for="key in shuffledKeys" :key="key"
                    @click="handleKeyPress(key)"
                    class="h-14 rounded-xl bg-white/5 hover:bg-white/10 active:bg-indigo-600 active:scale-95 transition-all text-xl font-bold text-white border border-white/10 flex items-center justify-center shadow-lg"
            >
                {{ key }}
            </button>
            <button @click="handleClear" class="h-14 rounded-xl bg-red-500/10 hover:bg-red-500/20 text-red-400 font-bold border border-red-500/20 flex items-center justify-center">
                C
            </button>
            <button @click="handleBackspace" class="h-14 rounded-xl bg-white/5 hover:bg-white/10 text-white font-bold border border-white/10 flex items-center justify-center">
                ⌫
            </button>
        </div>

        <button 
          @click="handleSubmit" 
          :disabled="currentLength < 5 || loading" 
          class="w-full p-4 rounded-xl border-none text-base font-semibold cursor-pointer mb-4 text-white transition-all hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed bg-gradient-to-br from-indigo-500 to-purple-600 shadow-lg shadow-indigo-500/20"
        >
          <span v-if="loading" class="flex items-center justify-center gap-2">
            <svg class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            Configuration...
          </span>
          <span v-else>{{ step === 1 ? '→ Suivant' : '✓ Configurer le PIN' }}</span>
        </button>

        <button @click="handleClose" class="text-sm text-gray-400 hover:text-white transition-colors">
          Annuler
        </button>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { userAPI } from '~/composables/useApi'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Configurer votre PIN'
  }
})

const emit = defineEmits(['update:isOpen', 'success', 'close'])

const step = ref(1) // 1 = Enter PIN, 2 = Confirm PIN
const pinDigits = ref(['', '', '', '', ''])
const firstPin = ref('') // Store the first entry
const error = ref('')
const loading = ref(false)
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

// Reset state when modal opens
watch(() => props.isOpen, async (newVal) => {
  if (newVal) {
    step.value = 1
    pinDigits.value = ['', '', '', '', '']
    firstPin.value = ''
    error.value = ''
    shuffleKeys() // Randomize layout on every open
  }
})

const handleKeyPress = (key: number) => {
    if (currentLength.value < 5) {
        pinDigits.value[currentLength.value] = key.toString()
        error.value = ''
    }
}

const handleBackspace = () => {
    if (currentLength.value > 0) {
        pinDigits.value[currentLength.value - 1] = ''
        error.value = ''
    }
}

const handleClear = () => {
    pinDigits.value = ['', '', '', '', '']
    error.value = ''
}

const handleSubmit = async () => {
  const pin = pinDigits.value.join('')
  if (pin.length !== 5) return
  
  if (step.value === 1) {
    // Move to confirmation step
    firstPin.value = pin
    pinDigits.value = ['', '', '', '', '']
    step.value = 2
    shuffleKeys() // Reshuffle for confirmation
    return
  }
  
  // Step 2: Verify PINs match
  if (pin !== firstPin.value) {
    error.value = 'Les PINs ne correspondent pas'
    handleClear()
    shuffleKeys()
    return
  }
  
  // PINs match - submit to backend
  loading.value = true
  error.value = ''
  
  try {
    await userAPI.setupPin({ pin: pin, confirm_pin: pin })
    emit('success')
    emit('update:isOpen', false)
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Erreur lors de la configuration'
    handleClear()
    step.value = 1
    firstPin.value = ''
    shuffleKeys()
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  emit('close')
  emit('update:isOpen', false)
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
  background: #1a1a2e;
  border-radius: 1.5rem;
  padding: 2rem;
  text-align: center;
  width: 100%;
  max-width: 24rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255,255,255,0.1);
  transform: scale(1);
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
</style>
