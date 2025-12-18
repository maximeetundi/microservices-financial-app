<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden bg-[#0f0f1a]">
    <!-- Animated Background -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[20%] -right-[10%] w-[70vw] h-[70vw] bg-indigo-600/20 rounded-full blur-[100px] animate-blob"></div>
      <div class="absolute top-[20%] -left-[10%] w-[50vw] h-[50vw] bg-purple-600/20 rounded-full blur-[100px] animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-[20%] right-[20%] w-[60vw] h-[60vw] bg-blue-600/10 rounded-full blur-[100px] animate-blob animation-delay-4000"></div>
    </div>

    <!-- Main Container -->
    <div class="relative w-full max-w-lg z-10 animate-fade-in-up">
      <!-- Logo -->
      <div class="text-center mb-8">
        <NuxtLink to="/" class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-6 shadow-lg shadow-indigo-500/30 transform hover:scale-105 transition-transform duration-300">
          <span class="text-4xl">🏦</span>
        </NuxtLink>
        <h1 class="text-3xl font-bold text-white mb-2 tracking-tight">CryptoBank</h1>
        <p class="text-indigo-200/80">Créez votre compte sécurisé</p>
      </div>

      <!-- Register Card -->
      <div class="glass-card p-8 border border-white/10 shadow-2xl backdrop-blur-xl bg-white/5 relative overflow-hidden group">
        <!-- Shine Effect -->
        <div class="absolute inset-0 bg-gradient-to-tr from-white/0 via-white/5 to-white/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000 pointer-events-none"></div>

        <form @submit.prevent="handleRegister" class="space-y-5 relative z-10">
          <!-- Name Fields -->
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="firstName" class="block text-sm font-medium text-indigo-100">Prénom</label>
              <input
                id="firstName"
                v-model="registerForm.first_name"
                type="text"
                required
                class="input-premium w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="John"
              />
            </div>
            <div class="space-y-2">
              <label for="lastName" class="block text-sm font-medium text-indigo-100">Nom</label>
              <input
                id="lastName"
                v-model="registerForm.last_name"
                type="text"
                required
                class="input-premium w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="Doe"
              />
            </div>
          </div>

          <!-- Email -->
          <div class="space-y-2">
            <label for="email" class="block text-sm font-medium text-indigo-100">Adresse email</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-indigo-300">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                </svg>
              </div>
              <input
                id="email"
                v-model="registerForm.email"
                type="email"
                autocomplete="email"
                required
                class="input-premium pl-11 w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="exemple@email.com"
              />
            </div>
          </div>

          <!-- Phone & Country Grid (Modified for better fit) -->
          <div class="grid grid-cols-2 gap-4">
             <!-- Country -->
             <div class="space-y-2">
              <label for="country" class="block text-sm font-medium text-indigo-100">Pays</label>
              <select
                id="country"
                v-model="registerForm.country"
                required
                class="input-premium w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50 appearance-none"
              >
                <option value="" disabled class="bg-slate-900 text-gray-400">Sélectionner</option>
                <option v-for="country in countries" :key="country.code" :value="country.code" class="bg-slate-900 text-white">
                  {{ country.name }}
                </option>
              </select>
            </div>
             <!-- Phone -->
             <div class="space-y-2">
              <label for="phone" class="block text-sm font-medium text-indigo-100">Téléphone</label>
              <input
                id="phone"
                v-model="registerForm.phone"
                type="tel"
                class="input-premium w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="+33 6..."
              />
            </div>
          </div>

          <!-- Date of Birth -->
           <div class="space-y-2">
            <label for="dateOfBirth" class="block text-sm font-medium text-indigo-100">Date de naissance</label>
            <input
              id="dateOfBirth"
              v-model="registerForm.date_of_birth"
              type="date"
              required
              class="input-premium w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50 [color-scheme:dark]"
            />
          </div>

          <!-- Password -->
          <div class="space-y-2">
            <label for="password" class="block text-sm font-medium text-indigo-100">Mot de passe</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-indigo-300">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                id="password"
                v-model="registerForm.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                required
                minlength="8"
                class="input-premium pl-11 pr-12 w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
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

          <!-- Confirm Password -->
          <div class="space-y-2">
            <label for="confirmPassword" class="block text-sm font-medium text-indigo-100">Confirmer</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-indigo-300">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                id="confirmPassword"
                v-model="registerForm.confirm_password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                required
                class="input-premium pl-11 w-full bg-black/20 focus:bg-black/30 border-white/10 focus:border-indigo-500/50"
                placeholder="••••••••"
              />
            </div>
          </div>

          <!-- Terms -->
          <label class="flex items-start gap-3 cursor-pointer group pt-2">
            <div class="relative flex items-center pt-0.5">
              <input
                v-model="acceptTerms"
                type="checkbox"
                class="peer appearance-none w-5 h-5 rounded border border-white/20 bg-white/5 checked:bg-indigo-500 checked:border-indigo-500 transition-all cursor-pointer"
              />
              <svg class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-3.5 h-3.5 text-white opacity-0 peer-checked:opacity-100 transition-opacity pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <span class="text-sm text-indigo-200/80 group-hover:text-white transition-colors">
              J'accepte les <a href="#" class="text-indigo-400 hover:text-indigo-300 hover:underline">conditions d'utilisation</a>
            </span>
          </label>

          <!-- Messages -->
          <div v-if="error" class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 backdrop-blur-sm animate-shake">
            <div class="flex items-center gap-3">
               <div class="p-2 rounded-full bg-red-500/20 text-red-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <p class="text-sm text-red-200">{{ error }}</p>
            </div>
          </div>
          <div v-if="success" class="p-4 rounded-xl bg-green-500/10 border border-green-500/20 backdrop-blur-sm">
            <div class="flex items-center gap-3">
               <div class="p-2 rounded-full bg-green-500/20 text-green-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
              </div>
              <p class="text-sm text-green-200">{{ success }}</p>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading || !acceptTerms"
            class="btn-premium w-full py-4 text-white font-bold text-lg shadow-lg shadow-indigo-500/25 hover:shadow-indigo-500/40 transform hover:-translate-y-0.5 transition-all duration-300 disabled:opacity-70 disabled:cursor-not-allowed group relative overflow-hidden"
          >
            <div class="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
            <span class="flex items-center justify-center gap-2 relative z-10">
              <svg v-if="loading" class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? 'En cours...' : 'Créer mon compte' }}
            </span>
          </button>

          <!-- Login Link -->
          <div class="text-center pt-2">
             <p class="text-indigo-200/60 text-sm">
              Déjà un compte ?
              <NuxtLink to="/auth/login" class="text-indigo-400 hover:text-indigo-300 font-semibold hover:underline transition-all ml-1">
                Se connecter
              </NuxtLink>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useCountries } from '~/composables/useCountries'

definePageMeta({
  layout: false
})

const router = useRouter()
const config = useRuntimeConfig()
const { countries, getCurrencyByCountry } = useCountries()

const registerForm = ref({
  email: '',
  password: '',
  confirm_password: '',
  first_name: '',
  last_name: '',
  phone: '',
  date_of_birth: '',
  country: ''
})

const showPassword = ref(false)
const acceptTerms = ref(false)
const loading = ref(false)
const error = ref('')
const success = ref('')

const handleRegister = async () => {
  error.value = ''
  success.value = ''

  // Validate passwords match
  if (registerForm.value.password !== registerForm.value.confirm_password) {
    error.value = 'Les mots de passe ne correspondent pas'
    return
  }

  // Validate password strength
  if (registerForm.value.password.length < 8) {
    error.value = 'Le mot de passe doit contenir au moins 8 caractères'
    return
  }

  loading.value = true

  try {
    const apiUrl = config.public.apiBaseUrl || 'https://api.app.maximeetundi.store'
    // Format date as RFC3339
    const dateOfBirth = registerForm.value.date_of_birth ? `${registerForm.value.date_of_birth}T00:00:00Z` : ''
    
    // Debug logging
    console.log('Using country:', registerForm.value.country)
    console.log('Derived currency:', getCurrencyByCountry(registerForm.value.country))
    
    const payload = {
      email: registerForm.value.email,
      password: registerForm.value.password,
      first_name: registerForm.value.first_name,
      last_name: registerForm.value.last_name,
      phone: registerForm.value.phone || undefined,
      date_of_birth: dateOfBirth,
      country: registerForm.value.country,
      currency: getCurrencyByCountry(registerForm.value.country)
    }
    
    console.log('Registration Payload:', payload)

    const response = await $fetch(`${apiUrl}/auth-service/api/v1/auth/register`, {
      method: 'POST',
      body: payload
    })

    success.value = 'Compte créé avec succès ! Redirection...'
    
    // Store token if provided
    if (response.access_token) {
      localStorage.setItem('accessToken', response.access_token)
      useCookie('accessToken').value = response.access_token
    }

    // Redirect to login or dashboard
    setTimeout(() => {
      if (response.access_token) {
        router.push('/dashboard')
      } else {
        router.push('/auth/login')
      }
    }, 1500)

  } catch (err) {
    console.error('Registration error:', err)
    // Extract error message from various response formats
    const errorData = err.data || err.response?.data || {}
    const errorMessage = errorData.error || errorData.message || err.message || 'Erreur lors de l\'inscription.'
    
    // Map backend errors to user-friendly messages
    if (errorMessage.includes('email already registered')) {
      error.value = 'Cet email est déjà utilisé. Veuillez vous connecter ou utiliser un autre email.'
    } else if (errorMessage.includes('phone already registered')) {
      error.value = 'Ce numéro de téléphone est déjà utilisé.'
    } else {
      error.value = errorMessage
    }
  } finally {
    loading.value = false
  }
}
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
