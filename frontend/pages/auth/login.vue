<template>
  <div class="min-h-screen flex items-center justify-center p-4" style="background: linear-gradient(135deg, #0f0f1a 0%, #1a1a2e 50%, #16213e 100%);">
    <!-- Animated Background -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-indigo-500/20 rounded-full blur-3xl animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-500/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
    </div>

    <div class="relative w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-6">
          <span class="text-4xl">🏦</span>
        </div>
        <h1 class="text-3xl font-bold text-white mb-2">CryptoBank</h1>
        <p class="text-white/60">Connectez-vous à votre espace sécurisé</p>
      </div>

      <!-- Login Card -->
      <div class="glass-card p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-white/80 mb-2">
              Adresse email
            </label>
            <input
              id="email"
              v-model="loginForm.email"
              type="email"
              autocomplete="email"
              required
              class="input-premium"
              placeholder="exemple@email.com"
            />
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-white/80 mb-2">
              Mot de passe
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                class="input-premium pr-12"
                placeholder="••••••••"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/80 transition-colors"
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
          <div v-if="show2FA">
            <label for="totp" class="block text-sm font-medium text-white/80 mb-2">
              Code d'authentification
            </label>
            <input
              id="totp"
              v-model="loginForm.totp_code"
              type="text"
              maxlength="6"
              pattern="[0-9]{6}"
              required
              class="input-premium text-center text-lg tracking-widest"
              placeholder="000000"
            />
            <p class="mt-2 text-sm text-white/50">
              Entrez le code à 6 chiffres de votre application
            </p>
          </div>

          <!-- Remember & Forgot -->
          <div class="flex items-center justify-between">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="loginForm.remember_me"
                type="checkbox"
                class="w-4 h-4 rounded border-white/20 bg-white/10 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-sm text-white/60">Se souvenir de moi</span>
            </label>
            <NuxtLink to="/auth/forgot-password" class="text-sm text-indigo-400 hover:text-indigo-300">
              Mot de passe oublié ?
            </NuxtLink>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="p-4 rounded-xl bg-red-500/20 border border-red-500/30">
            <div class="flex items-center gap-3">
              <svg class="w-5 h-5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <p class="text-sm text-red-300">{{ errorMessage }}</p>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-premium w-full py-4 text-lg"
          >
            <span v-if="isLoading" class="flex items-center justify-center gap-2">
              <div class="loading-spinner w-5 h-5"></div>
              Connexion...
            </span>
            <span v-else>Se connecter</span>
          </button>

          <!-- Register Link -->
          <p class="text-center text-white/60">
            Pas encore de compte ?
            <NuxtLink to="/auth/register" class="text-indigo-400 hover:text-indigo-300 font-medium">
              Créer un compte
            </NuxtLink>
          </p>
        </form>
      </div>

      <!-- Security Notice -->
      <div class="mt-6 p-4 rounded-xl bg-white/5 border border-white/10">
        <div class="flex items-start gap-3">
          <svg class="w-5 h-5 text-indigo-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
          </svg>
          <p class="text-sm text-white/50">
            <strong class="text-white/70">Sécurité bancaire</strong> - Votre connexion est protégée par un chiffrement de niveau bancaire.
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div class="mt-8 flex justify-center gap-6 text-sm text-white/40">
        <NuxtLink to="/legal/privacy" class="hover:text-white/60">Confidentialité</NuxtLink>
        <NuxtLink to="/legal/terms" class="hover:text-white/60">CGU</NuxtLink>
        <NuxtLink to="/support" class="hover:text-white/60">Aide</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: false,
  auth: false
})

const authStore = useAuthStore()

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