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
        <p class="text-white/60">Créez votre compte sécurisé</p>
      </div>

      <!-- Register Card -->
      <div class="glass-card p-8">
        <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- Name Fields -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="firstName" class="block text-sm font-medium text-white/80 mb-2">
                Prénom
              </label>
              <input
                id="firstName"
                v-model="registerForm.first_name"
                type="text"
                required
                class="input-premium"
                placeholder="John"
              />
            </div>
            <div>
              <label for="lastName" class="block text-sm font-medium text-white/80 mb-2">
                Nom
              </label>
              <input
                id="lastName"
                v-model="registerForm.last_name"
                type="text"
                required
                class="input-premium"
                placeholder="Doe"
              />
            </div>
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-white/80 mb-2">
              Adresse email
            </label>
            <input
              id="email"
              v-model="registerForm.email"
              type="email"
              autocomplete="email"
              required
              class="input-premium"
              placeholder="exemple@email.com"
            />
          </div>

          <!-- Phone -->
          <div>
            <label for="phone" class="block text-sm font-medium text-white/80 mb-2">
              Téléphone (optionnel)
            </label>
            <input
              id="phone"
              v-model="registerForm.phone"
              type="tel"
              class="input-premium"
              placeholder="+33 6 12 34 56 78"
            />
          </div>

          <!-- Date of Birth -->
          <div>
            <label for="dateOfBirth" class="block text-sm font-medium text-white/80 mb-2">
              Date de naissance
            </label>
            <input
              id="dateOfBirth"
              v-model="registerForm.date_of_birth"
              type="date"
              required
              class="input-premium"
            />
          </div>

          <!-- Country -->
          <div>
            <label for="country" class="block text-sm font-medium text-white/80 mb-2">
              Pays
            </label>
            <select
              id="country"
              v-model="registerForm.country"
              required
              class="input-premium"
            >
              <option value="" disabled>Sélectionnez votre pays</option>
              <option value="FRA">France</option>
              <option value="BEL">Belgique</option>
              <option value="CHE">Suisse</option>
              <option value="CAN">Canada</option>
              <option value="USA">États-Unis</option>
              <option value="GBR">Royaume-Uni</option>
              <option value="DEU">Allemagne</option>
              <option value="ESP">Espagne</option>
              <option value="ITA">Italie</option>
              <option value="CMR">Cameroun</option>
              <option value="SEN">Sénégal</option>
              <option value="CIV">Côte d'Ivoire</option>
              <option value="MAR">Maroc</option>
              <option value="TUN">Tunisie</option>
              <option value="DZA">Algérie</option>
              <option value="COD">RD Congo</option>
              <option value="COG">Congo</option>
              <option value="GAB">Gabon</option>
              <option value="MLI">Mali</option>
              <option value="BFA">Burkina Faso</option>
              <option value="NER">Niger</option>
              <option value="TGO">Togo</option>
              <option value="BEN">Bénin</option>
              <option value="GIN">Guinée</option>
              <option value="MDG">Madagascar</option>
              <option value="HTI">Haïti</option>
              <option value="MRT">Mauritanie</option>
              <option value="RWA">Rwanda</option>
              <option value="BDI">Burundi</option>
              <option value="COM">Comores</option>
              <option value="MUS">Maurice</option>
              <option value="PRT">Portugal</option>
              <option value="NLD">Pays-Bas</option>
              <option value="LUX">Luxembourg</option>
              <option value="MCO">Monaco</option>
            </select>
          </div>


          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-white/80 mb-2">
              Mot de passe
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="registerForm.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                required
                minlength="8"
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
            <p class="mt-1 text-xs text-white/40">Minimum 8 caractères</p>
          </div>

          <!-- Confirm Password -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-white/80 mb-2">
              Confirmer le mot de passe
            </label>
            <input
              id="confirmPassword"
              v-model="registerForm.confirm_password"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="new-password"
              required
              class="input-premium"
              placeholder="••••••••"
            />
          </div>

          <!-- Terms -->
          <div class="flex items-start gap-3">
            <input
              id="terms"
              v-model="acceptTerms"
              type="checkbox"
              required
              class="mt-1 w-4 h-4 rounded border-white/20 bg-white/10 text-indigo-500 focus:ring-indigo-500"
            />
            <label for="terms" class="text-sm text-white/60">
              J'accepte les <a href="#" class="text-indigo-400 hover:underline">conditions d'utilisation</a> et la <a href="#" class="text-indigo-400 hover:underline">politique de confidentialité</a>
            </label>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="p-3 rounded-lg bg-red-500/20 border border-red-500/30">
            <p class="text-sm text-red-400">{{ error }}</p>
          </div>

          <!-- Success Message -->
          <div v-if="success" class="p-3 rounded-lg bg-green-500/20 border border-green-500/30">
            <p class="text-sm text-green-400">{{ success }}</p>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading || !acceptTerms"
            class="btn-premium w-full flex items-center justify-center"
          >
            <svg v-if="loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Création du compte...' : 'Créer mon compte' }}
          </button>
        </form>

        <!-- Divider -->
        <div class="my-6 flex items-center">
          <div class="flex-1 border-t border-white/10"></div>
          <span class="px-4 text-sm text-white/40">ou</span>
          <div class="flex-1 border-t border-white/10"></div>
        </div>

        <!-- Login Link -->
        <p class="text-center text-white/60">
          Déjà un compte ?
          <NuxtLink to="/auth/login" class="text-indigo-400 hover:text-indigo-300 font-medium">
            Se connecter
          </NuxtLink>
        </p>
      </div>

      <!-- Security Badge -->
      <div class="mt-6 text-center">
        <div class="inline-flex items-center gap-2 text-white/40 text-sm">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
          Connexion sécurisée 256-bit SSL
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

definePageMeta({
  layout: false
})

const router = useRouter()
const config = useRuntimeConfig()

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
    const response = await $fetch(`${apiUrl}/auth-service/api/v1/auth/register`, {
      method: 'POST',
      body: {
        email: registerForm.value.email,
        password: registerForm.value.password,
        first_name: registerForm.value.first_name,
        last_name: registerForm.value.last_name,
        phone: registerForm.value.phone || undefined,
        date_of_birth: registerForm.value.date_of_birth,
        country: registerForm.value.country
      }
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
    error.value = err.data?.message || err.message || 'Erreur lors de l\'inscription. Veuillez réessayer.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.glass-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
}

.input-premium {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: white;
  transition: all 0.3s;
}

.input-premium::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.input-premium:focus {
  outline: none;
  border-color: rgba(99, 102, 241, 0.5);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.btn-premium {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 12px;
  color: white;
  font-weight: 600;
  transition: all 0.3s;
}

.btn-premium:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.4);
}

.btn-premium:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
