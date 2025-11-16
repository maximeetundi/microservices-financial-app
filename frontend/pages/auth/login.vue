<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div class="text-center">
        <div class="flex justify-center">
          <div class="w-16 h-16 bg-blue-600 rounded-xl flex items-center justify-center">
            <span class="text-white font-bold text-xl">CB</span>
          </div>
        </div>
        <h2 class="mt-6 text-3xl font-bold text-gray-900">Welcome back</h2>
        <p class="mt-2 text-sm text-gray-600">
          Sign in to your Crypto Bank account
        </p>
      </div>

      <!-- Login Form -->
      <div class="bg-white rounded-xl shadow-lg p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          
          <!-- Email Field -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email address
            </label>
            <input
              id="email"
              v-model="loginForm.email"
              type="email"
              autocomplete="email"
              required
              class="w-full px-3 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter your email"
            />
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                class="w-full px-3 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10"
                placeholder="Enter your password"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg v-if="showPassword" class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                </svg>
                <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
              </button>
            </div>
          </div>

          <!-- 2FA Field (if enabled) -->
          <div v-if="show2FA">
            <label for="totp" class="block text-sm font-medium text-gray-700 mb-2">
              Two-Factor Authentication Code
            </label>
            <input
              id="totp"
              v-model="loginForm.totp_code"
              type="text"
              maxlength="6"
              pattern="[0-9]{6}"
              required
              class="w-full px-3 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter 6-digit code"
            />
            <p class="mt-2 text-sm text-gray-500">
              Enter the 6-digit code from your authenticator app
            </p>
          </div>

          <!-- Remember Me & Forgot Password -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember-me"
                v-model="loginForm.remember_me"
                type="checkbox"
                class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label for="remember-me" class="ml-2 block text-sm text-gray-700">
                Remember me
              </label>
            </div>

            <div class="text-sm">
              <NuxtLink to="/auth/forgot-password" class="font-medium text-blue-600 hover:text-blue-500">
                Forgot your password?
              </NuxtLink>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="bg-red-50 border border-red-200 rounded-lg p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-red-800">{{ errorMessage }}</p>
              </div>
            </div>
          </div>

          <!-- Submit Button -->
          <div>
            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg v-if="isLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              {{ isLoading ? 'Signing in...' : 'Sign in' }}
            </button>
          </div>

          <!-- Sign Up Link -->
          <div class="text-center">
            <p class="text-sm text-gray-600">
              Don't have an account?
              <NuxtLink to="/auth/register" class="font-medium text-blue-600 hover:text-blue-500">
                Create account
              </NuxtLink>
            </p>
          </div>
        </form>

        <!-- Social Login (Optional) -->
        <div class="mt-6">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-300"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-white text-gray-500">Or continue with</span>
            </div>
          </div>

          <div class="mt-6 grid grid-cols-2 gap-3">
            <button
              @click="loginWithGoogle"
              class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
            >
              <svg class="w-5 h-5" viewBox="0 0 24 24">
                <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              <span class="ml-2">Google</span>
            </button>

            <button
              @click="loginWithApple"
              class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12.017 0C6.624 0 2.246 4.377 2.246 9.771s4.378 9.771 9.771 9.771 9.771-4.377 9.771-9.771S17.41 0 12.017 0zm2.392 6.9c.378 0 .756.189.945.567.189.378.189.756 0 1.134-.189.378-.567.567-.945.567s-.756-.189-.945-.567c-.189-.378-.189-.756 0-1.134.189-.378.567-.567.945-.567zm-4.784 0c.378 0 .756.189.945.567.189.378.189.756 0 1.134-.189.378-.567.567-.945.567s-.756-.189-.945-.567c-.189-.378-.189-.756 0-1.134.189-.378.567-.567.945-.567z"/>
              </svg>
              <span class="ml-2">Apple</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Security Notice -->
      <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-blue-800">
              <strong>Security Notice:</strong> Your connection is protected with bank-level encryption. 
              Never share your login credentials with anyone.
            </p>
          </div>
        </div>
      </div>

      <!-- Footer Links -->
      <div class="text-center text-sm text-gray-500 space-x-4">
        <NuxtLink to="/legal/privacy" class="hover:text-gray-700">Privacy Policy</NuxtLink>
        <span>•</span>
        <NuxtLink to="/legal/terms" class="hover:text-gray-700">Terms of Service</NuxtLink>
        <span>•</span>
        <NuxtLink to="/support" class="hover:text-gray-700">Support</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// Page meta
definePageMeta({
  layout: false,
  auth: false
})

// Reactive data
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

// Methods
const handleLogin = async () => {
  if (isLoading.value) return

  isLoading.value = true
  errorMessage.value = ''

  try {
    // First attempt login with email and password
    const response = await $fetch('/api/auth/login', {
      method: 'POST',
      body: {
        email: loginForm.value.email,
        password: loginForm.value.password,
        totp_code: loginForm.value.totp_code,
        remember_me: loginForm.value.remember_me
      }
    })

    if (response.requires_2fa && !loginForm.value.totp_code) {
      // Show 2FA field
      show2FA.value = true
      errorMessage.value = 'Please enter your two-factor authentication code.'
      return
    }

    // Login successful
    if (response.token) {
      // Store auth token (you might want to use a more secure method)
      const authStore = useAuthStore() // Assuming you have an auth store
      authStore.setToken(response.token)
      authStore.setUser(response.user)

      // Redirect to dashboard
      await navigateTo('/dashboard')
    }

  } catch (error) {
    console.error('Login error:', error)
    
    if (error.data?.message) {
      errorMessage.value = error.data.message
    } else if (error.data?.error) {
      errorMessage.value = error.data.error
    } else {
      errorMessage.value = 'Login failed. Please check your credentials and try again.'
    }

    // Reset 2FA if error occurs
    if (show2FA.value) {
      loginForm.value.totp_code = ''
    }
  } finally {
    isLoading.value = false
  }
}

const loginWithGoogle = async () => {
  try {
    // Implement Google OAuth login
    window.location.href = '/api/auth/google'
  } catch (error) {
    console.error('Google login error:', error)
    errorMessage.value = 'Google login failed. Please try again.'
  }
}

const loginWithApple = async () => {
  try {
    // Implement Apple OAuth login
    window.location.href = '/api/auth/apple'
  } catch (error) {
    console.error('Apple login error:', error)
    errorMessage.value = 'Apple login failed. Please try again.'
  }
}

// Auto-focus email field
onMounted(() => {
  document.getElementById('email')?.focus()
})
</script>

<style scoped>
/* Custom animations */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.form-container {
  animation: fadeIn 0.6s ease-out;
}
</style>