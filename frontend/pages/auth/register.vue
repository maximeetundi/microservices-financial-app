<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden" :class="isDark ? 'bg-[#0f0f1a]' : 'bg-gradient-to-br from-slate-50 to-indigo-50'">
    <!-- Animated Background -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[20%] -right-[10%] w-[70vw] h-[70vw] bg-indigo-600/20 rounded-full blur-[100px] animate-blob"></div>
      <div class="absolute top-[20%] -left-[10%] w-[50vw] h-[50vw] bg-purple-600/20 rounded-full blur-[100px] animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-[20%] right-[20%] w-[60vw] h-[60vw] bg-blue-600/10 rounded-full blur-[100px] animate-blob animation-delay-4000"></div>
    </div>

    <!-- Main Container -->
    <div class="relative w-full max-w-lg z-10 animate-fade-in-up">
      <!-- Logo -->
      <div class="text-center mb-6">
        <NuxtLink to="/" class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 mb-4 shadow-lg shadow-indigo-500/30 transform hover:scale-105 transition-transform duration-300">
          <span class="text-3xl">🏦</span>
        </NuxtLink>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-1 tracking-tight">Créer votre compte</h1>
        <p class="text-gray-600 dark:text-indigo-200/70 text-sm">Rejoignez Zekora en quelques étapes</p>
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

      <!-- Progress Steps -->
      <div class="flex items-center justify-center gap-2 mb-6">
        <template v-for="step in 3" :key="step">
          <div 
            class="flex items-center justify-center w-10 h-10 rounded-full font-bold text-sm transition-all duration-300"
            :class="currentStep >= step 
              ? 'bg-gradient-to-br from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30' 
              : 'bg-black/10 dark:bg-white/10 text-gray-400 dark:text-white/50'"
          >
            <svg v-if="currentStep > step" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span v-else>{{ step }}</span>
          </div>
          <div v-if="step < 3" class="w-12 h-1 rounded-full transition-all duration-300" :class="currentStep > step ? 'bg-indigo-500' : 'bg-black/10 dark:bg-white/10'"></div>
        </template>
      </div>

      <!-- Step Labels -->
      <div class="flex justify-between text-xs text-gray-500 dark:text-slate-400 mb-8 px-4">
        <span :class="currentStep >= 1 ? 'text-indigo-400' : ''">Identité</span>
        <span :class="currentStep >= 2 ? 'text-indigo-400' : ''">Contact</span>
        <span :class="currentStep >= 3 ? 'text-indigo-400' : ''">Sécurité</span>
      </div>

      <!-- Register Card -->
      <div class="glass-card p-8 border border-gray-200 dark:border-white/10 shadow-2xl backdrop-blur-xl bg-white/90 dark:bg-white/5 relative overflow-hidden">
        <form @submit.prevent="handleSubmit">
            <!-- Step 1: Identity -->
          <div v-show="currentStep === 1" class="space-y-5">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Prénom</label>
                <input
                  v-model="form.first_name"
                  type="text"
                  required
                  class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                  placeholder="John"
                />
              </div>
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Nom</label>
                <input
                  v-model="form.last_name"
                  type="text"
                  required
                  class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                  placeholder="Doe"
                />
              </div>
            </div>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Date de naissance</label>
              <input
                v-model="form.date_of_birth"
                type="date"
                required
                class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white [color-scheme:light] dark:[color-scheme:dark]"
              />
            </div>
          </div>

          <!-- Step 2: Contact -->
          <div v-show="currentStep === 2" class="space-y-5">
            <div class="space-y-2">
              <label class="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-indigo-100">
                <svg class="w-4 h-4 text-indigo-500 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                </svg>
                Email
              </label>
              <input
                v-model="form.email"
                type="email"
                required
                class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                placeholder="exemple@email.com"
              />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Pays</label>
                <select
                  v-model="form.country"
                  required
                  class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white appearance-none"
                >
                  <option value="" disabled class="text-gray-400">Sélectionner</option>
                  <option v-for="country in countries" :key="country.code" :value="country.code" class="text-gray-900 dark:text-white bg-white dark:bg-slate-900">
                    {{ country.name }}
                  </option>
                </select>
              </div>
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Téléphone</label>
                <input
                  v-model="form.phone"
                  type="tel"
                  class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                  placeholder="+33 6..."
                />
              </div>
            </div>
          </div>

          <!-- Step 3: Security -->
          <div v-show="currentStep === 3" class="space-y-5">
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Mot de passe</label>
              <div class="relative">
                <input
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  minlength="8"
                  class="input-premium pr-12 w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                  placeholder="••••••••"
                />
                <button type="button" @click="showPassword = !showPassword" class="absolute right-4 top-1/2 -translate-y-1/2 text-indigo-400 hover:text-indigo-600 dark:hover:text-white transition-colors">
                  <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path></svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                </button>
              </div>
              <!-- Password Strength -->
              <div class="flex gap-1 mt-2">
                <div v-for="i in 4" :key="i" class="h-1 flex-1 rounded-full transition-all" :class="passwordStrength >= i ? strengthColors[passwordStrength] : 'bg-gray-200 dark:bg-white/10'"></div>
              </div>
              <p class="text-xs" :class="strengthTextColors[passwordStrength]">{{ strengthTexts[passwordStrength] }}</p>
            </div>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-indigo-100">Confirmer</label>
              <input
                v-model="form.confirm_password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="input-premium w-full bg-gray-50 dark:bg-black/20 focus:bg-white dark:focus:bg-black/30 border-gray-300 dark:border-white/10 focus:border-indigo-500 dark:focus:border-indigo-500/50 text-gray-900 dark:text-white"
                placeholder="••••••••"
              />
            </div>
            <label class="flex items-start gap-3 cursor-pointer group pt-2">
              <div class="relative flex items-center pt-0.5">
                <input v-model="acceptTerms" type="checkbox" class="peer appearance-none w-5 h-5 rounded border border-gray-300 dark:border-white/20 bg-gray-50 dark:bg-white/5 checked:bg-indigo-500 checked:border-indigo-500 transition-all cursor-pointer"/>
                <svg class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-3.5 h-3.5 text-white opacity-0 peer-checked:opacity-100 transition-opacity pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
              </div>
              <span class="text-sm text-gray-600 dark:text-indigo-200/80 group-hover:text-indigo-600 dark:group-hover:text-white transition-colors">
                J'accepte les <NuxtLink to="/legal/terms" class="text-indigo-500 dark:text-indigo-400 hover:underline">conditions</NuxtLink> et la <NuxtLink to="/legal/privacy" class="text-indigo-500 dark:text-indigo-400 hover:underline">confidentialité</NuxtLink>
              </span>
            </label>
          </div>

          <!-- Error/Success Messages -->
          <div v-if="error" class="mt-4 p-3 rounded-xl bg-red-500/10 border border-red-500/20 flex items-center gap-3">
            <div class="p-1.5 rounded-full bg-red-500/20 text-red-400">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </div>
            <p class="text-sm text-red-600 dark:text-red-200">{{ error }}</p>
          </div>
          <div v-if="success" class="mt-4 p-3 rounded-xl bg-green-500/10 border border-green-500/20 flex items-center gap-3">
            <div class="p-1.5 rounded-full bg-green-500/20 text-green-400">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
            </div>
            <p class="text-sm text-green-600 dark:text-green-200">{{ success }}</p>
          </div>

          <!-- Navigation Buttons -->
          <div class="flex gap-3 mt-6">
            <button
              v-if="currentStep > 1"
              type="button"
              @click="currentStep--"
              class="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 dark:text-white bg-gray-100 dark:bg-white/10 hover:bg-gray-200 dark:hover:bg-white/20 transition-all"
            >
              Retour
            </button>
            <button
              v-if="currentStep < 3"
              type="button"
              @click="nextStep"
              class="flex-1 btn-premium py-3 font-semibold"
            >
              Continuer
            </button>
            <button
              v-else
              type="submit"
              :disabled="loading || !acceptTerms"
              class="flex-1 btn-premium py-3 font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>
                Création...
              </span>
              <span v-else>Créer mon compte</span>
            </button>
          </div>

          <!-- Login Link -->
          <div class="text-center pt-4">
            <p class="text-gray-500 dark:text-indigo-200/60 text-sm">
              Déjà un compte ?
              <NuxtLink to="/auth/login" class="text-indigo-500 dark:text-indigo-400 hover:text-indigo-600 dark:hover:text-indigo-300 font-semibold hover:underline ml-1">Se connecter</NuxtLink>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useCountries } from '~/composables/useCountries'

definePageMeta({ layout: false })

const router = useRouter()
const config = useRuntimeConfig()
const { countries, getCurrencyByCountry, getDialCodeByCountry } = useCountries()
const colorMode = useColorMode()

const isDark = computed(() => colorMode.value === 'dark')
const toggleTheme = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}

const currentStep = ref(1)
const showPassword = ref(false)
const acceptTerms = ref(false)
const loading = ref(false)
const error = ref('')
const success = ref('')

const form = ref({
  first_name: '',
  last_name: '',
  date_of_birth: '',
  email: '',
  country: '',
  phone: '',
  password: '',
  confirm_password: ''
})

// Watch country selection to auto-fill phone dial code
watch(() => form.value.country, (newCountry) => {
  if (newCountry) {
    const dialCode = getDialCodeByCountry(newCountry)
    // Only auto-fill if phone is empty or just contains a dial code
    if (!form.value.phone || form.value.phone.match(/^\+[\d-]+$/)) {
      form.value.phone = dialCode + ' '
    }
  }
})

// Password strength
const passwordStrength = computed(() => {
  const p = form.value.password
  if (!p) return 0
  let score = 0
  if (p.length >= 8) score++
  if (/[A-Z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++
  return score
})

const strengthColors = { 1: 'bg-red-500', 2: 'bg-orange-500', 3: 'bg-yellow-500', 4: 'bg-green-500' }
const strengthTextColors = { 0: 'text-slate-500', 1: 'text-red-400', 2: 'text-orange-400', 3: 'text-yellow-400', 4: 'text-green-400' }
const strengthTexts = { 0: 'Entrez un mot de passe', 1: 'Faible', 2: 'Moyen', 3: 'Bon', 4: 'Excellent' }

const nextStep = () => {
  error.value = ''
  if (currentStep.value === 1) {
    if (!form.value.first_name || !form.value.last_name || !form.value.date_of_birth) {
      error.value = 'Veuillez remplir tous les champs'
      return
    }
  }
  if (currentStep.value === 2) {
    if (!form.value.email || !form.value.country) {
      error.value = 'Email et pays sont requis'
      return
    }
  }
  currentStep.value++
}

const handleSubmit = async () => {
  error.value = ''
  success.value = ''
  
  if (form.value.password !== form.value.confirm_password) {
    error.value = 'Les mots de passe ne correspondent pas'
    return
  }
  if (form.value.password.length < 8) {
    error.value = 'Mot de passe: 8 caractères minimum'
    return
  }

  loading.value = true
  try {
    const apiUrl = config.public.apiBaseUrl || 'https://api.app.tech-afm.com'
    const dateOfBirth = form.value.date_of_birth ? `${form.value.date_of_birth}T00:00:00Z` : ''
    
    const response = await $fetch(`${apiUrl}/auth-service/api/v1/auth/register`, {
      method: 'POST',
      body: {
        email: form.value.email,
        password: form.value.password,
        first_name: form.value.first_name,
        last_name: form.value.last_name,
        phone: form.value.phone || undefined,
        date_of_birth: dateOfBirth,
        country: form.value.country,
        currency: getCurrencyByCountry(form.value.country)
      }
    })

    success.value = 'Compte créé ! Redirection...'
    if (response.access_token) {
      localStorage.setItem('accessToken', response.access_token)
      useCookie('accessToken').value = response.access_token
    }
    setTimeout(() => router.push(response.access_token ? '/dashboard' : '/auth/login'), 1500)
  } catch (err) {
    const msg = err.data?.error || err.data?.message || err.message || 'Erreur lors de l\'inscription'
    error.value = msg.includes('email already') ? 'Email déjà utilisé' : msg.includes('phone already') ? 'Téléphone déjà utilisé' : msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob { animation: blob 7s infinite; }
.animation-delay-2000 { animation-delay: 2s; }
.animation-delay-4000 { animation-delay: 4s; }

input:-webkit-autofill {
  -webkit-box-shadow: 0 0 0 30px rgba(0, 0, 0, 0.2) inset !important;
  -webkit-text-fill-color: white !important;
}

.bg-gradient-to-br input:-webkit-autofill {
  -webkit-box-shadow: 0 0 0 30px #f9fafb inset !important;
  -webkit-text-fill-color: #1f2937 !important;
}
</style>
