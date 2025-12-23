<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden" :class="isDark ? 'bg-[#0f0f1a]' : 'bg-gradient-to-br from-slate-50 to-indigo-50'">
    <!-- Animated Background -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[20%] -right-[10%] w-[70vw] h-[70vw] bg-indigo-600/20 rounded-full blur-[100px] animate-blob"></div>
      <div class="absolute top-[20%] -left-[10%] w-[50vw] h-[50vw] bg-purple-600/20 rounded-full blur-[100px] animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-[20%] right-[20%] w-[60vw] h-[60vw] bg-blue-600/10 rounded-full blur-[100px] animate-blob animation-delay-4000"></div>
    </div>

    <!-- Main Container -->
    <div class="relative w-full max-w-md z-10 animate-fade-in-up">
      <!-- Logo -->
      <div class="text-center mb-8">
        <NuxtLink to="/" class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-6 shadow-lg shadow-indigo-500/30 transform hover:scale-105 transition-transform duration-300">
          <span class="text-4xl">🏦</span>
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2 tracking-tight">Zekora</h1>
        <p class="text-gray-600 dark:text-indigo-200/80">Connectez-vous à votre espace sécurisé</p>
      </div>

      <!-- Theme Toggle Button -->
      <button @click="toggleTheme" class="absolute top-4 right-4 p-3 rounded-xl bg-white/80 dark:bg-white/10 hover:bg-white dark:hover:bg-white/20 border border-gray-200 dark:border-white/10 transition-all shadow-lg z-20">
        <svg v-if="isDark" class="w-5 h-5 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clip-rule="evenodd" />
        </svg>
        <svg v-else class="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
          <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
        </svg>
      </button>

      <!-- Login Card -->
      <div class="glass-card p-8 border border-gray-200 dark:border-white/10 shadow-2xl backdrop-blur-xl bg-white/90 dark:bg-white/5 relative overflow-hidden group">
        <!-- Shine Effect -->
        <div class="absolute inset-0 bg-gradient-to-tr from-white/0 via-white/5 to-white/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000 pointer-events-none"></div>

        <form @submit.prevent="handleLogin" class="space-y-6 relative z-10">
          <!-- Email -->
          <div class="space-y-2">
            <label for="email" class="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-indigo-100">
              <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
              </svg>
              Adresse email
            </label>
            <input
              id="email"
              v-model="loginForm.email"
              type="email"
              autocomplete="email"
              required
              class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="exemple@email.com"
            />
          </div>

          <!-- Password -->
          <div class="space-y-2">
            <label for="password" class="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-indigo-100">
              <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              Mot de passe
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                class="input-premium pr-12 w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
                placeholder="••••••••"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-indigo-300 hover:text-white transition-colors focus:outline-none"
              >
                <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
              </button>
            </div>
          </div>

          <!-- 2FA Field -->
          <div v-if="show2FA" class="space-y-2 animate-fade-in-up">
            <label for="totp" class="block text-sm font-medium text-indigo-100">
              Code d'authentification
            </label>
            <div class="relative">
              <input
                id="totp"
                v-model="loginForm.totp_code"
                type="text"
                maxlength="6"
                pattern="[0-9]{6}"
                required
                class="input-premium text-center text-2xl tracking-[0.5em] font-mono bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="000000"
              />
            </div>
            <p class="text-xs text-indigo-300 text-center">
              Entrez le code à 6 chiffres de votre application
            </p>
          </div>

          <!-- Remember & Forgot -->
          <div class="flex items-center justify-between pt-2">
            <label class="flex items-center gap-2 cursor-pointer group">
              <div class="relative flex items-center">
                <input
                  v-model="loginForm.remember_me"
                  type="checkbox"
                  class="peer appearance-none w-5 h-5 rounded border border-white/20 bg-white/5 checked:bg-indigo-500 checked:border-indigo-500 transition-all cursor-pointer"
                />
                <svg class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-3.5 h-3.5 text-white opacity-0 peer-checked:opacity-100 transition-opacity pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <span class="text-sm text-indigo-200 group-hover:text-white transition-colors">Se souvenir de moi</span>
            </label>
            <NuxtLink to="/auth/forgot-password" class="text-sm text-indigo-400 hover:text-indigo-300 font-medium hover:underline transition-all">
              Mot de passe oublié ?
            </NuxtLink>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 backdrop-blur-sm animate-shake">
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-full bg-red-500/20 text-red-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <p class="text-sm text-red-200 font-medium">{{ errorMessage }}</p>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-premium w-full py-4 text-white font-bold text-lg shadow-lg shadow-indigo-500/25 hover:shadow-indigo-500/40 transform hover:-translate-y-0.5 transition-all duration-300 disabled:opacity-70 disabled:cursor-not-allowed group relative overflow-hidden"
          >
            <div class="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
            <span v-if="isLoading" class="flex items-center justify-center gap-2 relative z-10">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Connexion...
            </span>
            <span v-else class="relative z-10">Se connecter</span>
          </button>

          <!-- Register Link -->
          <div class="text-center pt-2">
            <p class="text-indigo-200/60 text-sm">
              Pas encore de compte ?
              <NuxtLink to="/auth/register" class="text-indigo-400 hover:text-indigo-300 font-semibold hover:underline transition-all ml-1">
                Créer un compte
              </NuxtLink>
            </p>
          </div>
        </form>
      </div>

      <!-- Footer -->
      <div class="mt-8 flex justify-center gap-6 text-xs font-medium text-indigo-200/40">
        <NuxtLink to="/legal/privacy" class="hover:text-indigo-300 transition-colors">Confidentialité</NuxtLink>
        <span class="text-indigo-200/20">•</span>
        <NuxtLink to="/legal/terms" class="hover:text-indigo-300 transition-colors">CGU</NuxtLink>
        <span class="text-indigo-200/20">•</span>
        <NuxtLink to="/support" class="hover:text-indigo-300 transition-colors">Aide</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: false,
  auth: false
})

const authStore = useAuthStore()
const colorMode = useColorMode()

const isDark = computed(() => colorMode.value === 'dark')
const toggleTheme = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}

const loginForm = ref({
  email: '',
  password: '',
  totp_code: '',
  remember_me: false
})

const showPassword = ref(false)
const show2FA = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  if (isLoading.value) return

  isLoading.value = true
  errorMessage.value = ''

  try {
    const result = await authStore.login(
      loginForm.value.email,
      loginForm.value.password,
      loginForm.value.totp_code || undefined
    )

    if (result.requires2FA) {
      show2FA.value = true
      errorMessage.value = 'Veuillez entrer votre code d\'authentification.'
      return
    }

    if (result.success) {
      await navigateTo('/dashboard')
    } else {
      errorMessage.value = result.error || 'Identifiants incorrects. Veuillez réessayer.'
    }
  } catch (error) {
    console.error('Login error:', error)
    errorMessage.value = 'Une erreur est survenue. Veuillez réessayer.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  document.getElementById('email')?.focus()
})
</script>

<style scoped>
/* Blob Animation */
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob {
  animation: blob 7s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
.animation-delay-4000 {
  animation-delay: 4s;
}

/* Shake Animation for Error */
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
.animate-shake {
  animation: shake 0.5s cubic-bezier(.36,.07,.19,.97) both;
}

/* Autofill Hack for Dark Theme */
input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px rgba(0, 0, 0, 0.2) inset !important;
  -webkit-text-fill-color: white !important;
  transition: background-color 5000s ease-in-out 0s;
  background-color: transparent !important;
}

/* Custom Scrollbar if needed */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}
</style>