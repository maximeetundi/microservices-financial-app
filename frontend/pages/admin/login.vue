<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl mb-4">
          <span class="text-3xl">🛡️</span>
        </div>
        <h1 class="text-3xl font-bold text-white">Admin Panel</h1>
        <p class="text-slate-400 mt-2">Connectez-vous à votre espace administrateur</p>
      </div>

      <!-- Login Form -->
      <div class="glass-card p-8">
        <form @submit.prevent="handleLogin">
          <!-- Error Message -->
          <div v-if="error" class="mb-4 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
            {{ error }}
          </div>

          <!-- Email -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-slate-300 mb-2">Email</label>
            <input 
              v-model="form.email" 
              type="email" 
              required
              placeholder="admin@example.com"
              class="w-full px-4 py-3 bg-slate-800 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
            />
          </div>

          <!-- Password -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-slate-300 mb-2">Mot de passe</label>
            <div class="relative">
              <input 
                v-model="form.password" 
                :type="showPassword ? 'text' : 'password'" 
                required
                placeholder="••••••••"
                class="w-full px-4 py-3 bg-slate-800 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all pr-12"
              />
              <button 
                type="button" 
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white transition-colors"
              >
                <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button 
            type="submit" 
            :disabled="loading"
            class="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-semibold rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="loading" class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            {{ loading ? 'Connexion...' : 'Se connecter' }}
          </button>
        </form>

        <!-- Demo Credentials -->
        <div class="mt-6 p-4 bg-slate-800/50 rounded-xl border border-slate-700/50">
          <p class="text-xs text-slate-400 mb-2">Identifiants de démonstration :</p>
          <div class="flex flex-col gap-1 text-sm">
            <div class="flex justify-between">
              <span class="text-slate-500">Email:</span>
              <span class="text-slate-300 font-mono">admin@crypto-bank.com</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-500">Mot de passe:</span>
              <span class="text-slate-300 font-mono">Admin123!</span>
            </div>
          </div>
          <button 
            @click="fillDemoCredentials" 
            class="mt-3 w-full py-2 text-sm text-indigo-400 hover:text-indigo-300 hover:bg-indigo-500/10 rounded-lg transition-colors"
          >
            Remplir automatiquement
          </button>
        </div>
      </div>

      <!-- Back to app -->
      <div class="text-center mt-6">
        <NuxtLink to="/dashboard" class="text-slate-400 hover:text-white transition-colors text-sm">
          ← Retour à l'application
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})
const loading = ref(false)
const error = ref('')
const showPassword = ref(false)

const API_URL = 'https://api.app.tech-afm.com'

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await axios.post(`${API_URL}/admin-service/api/v1/admin/login`, {
      email: form.email,
      password: form.password
    })

    if (response.data.token) {
      localStorage.setItem('adminToken', response.data.token)
      localStorage.setItem('adminUser', JSON.stringify(response.data.admin))
      router.push('/admin/dashboard')
    }
  } catch (e) {
    console.error('Login error:', e)
    if (e.response?.data?.error) {
      error.value = e.response.data.error
    } else {
      error.value = 'Erreur de connexion. Vérifiez vos identifiants.'
    }
  } finally {
    loading.value = false
  }
}

const fillDemoCredentials = () => {
  form.email = 'admin@crypto-bank.com'
  form.password = 'Admin123!'
}

definePageMeta({
  layout: false
})
</script>

<style scoped>
.glass-card {
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(71, 85, 105, 0.5);
  border-radius: 1.5rem;
}
</style>
