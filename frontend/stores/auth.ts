import { defineStore } from 'pinia'
import api, { resetLogoutFlag } from '~/composables/useApi'

interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  phone: string
  country: string
  kycStatus: string
  kycLevel: number
  isActive: boolean
  twoFaEnabled: boolean
  emailVerified: boolean
  phoneVerified: boolean
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
    isLoading: false
  }),

  actions: {
    async initializeAuth() {
      // Skip initialization if already on auth pages
      if (typeof window !== 'undefined' && window.location.pathname.startsWith('/auth')) {
        return
      }

      // Skip if already authenticated with valid user data
      if (this.isAuthenticated && this.user) {
        return
      }

      // Skip if already loading to prevent concurrent initialization
      if (this.isLoading) {
        return
      }

      const token = localStorage.getItem('accessToken')
      const refreshToken = localStorage.getItem('refreshToken')

      if (token && refreshToken) {
        this.accessToken = token
        this.refreshToken = refreshToken
        // Set authenticated immediately - don't wait for profile
        this.isAuthenticated = true

        // Sync to cookies for SSR middleware
        if (typeof document !== 'undefined') {
          document.cookie = `accessToken=${token}; path=/; max-age=86400; SameSite=Lax`
          document.cookie = `refreshToken=${refreshToken}; path=/; max-age=604800; SameSite=Lax`
        }

        // Only fetch profile if we don't already have user data
        if (!this.user) {
          try {
            this.isLoading = true
            await this.fetchUserProfile()
          } catch (error: any) {
            // Profile fetch failed - check if it's a 401 (token expired)
            if (error.response?.status === 401) {
              // Token is invalid, let the API interceptor handle refresh
              console.warn('Profile fetch returned 401, token may need refresh')
            } else {
              // Other error - keep auth state, user can retry
              console.warn('Profile fetch failed, keeping auth state:', error)
            }
          } finally {
            this.isLoading = false
          }
        }
      }
    },

    async login(email: string, password: string, twoFaCode?: string) {
      this.isLoading = true
      try {
        const response = await api.post('/auth-service/api/v1/auth/login', {
          email,
          password,
          two_fa_code: twoFaCode
        })

        const { access_token, refresh_token, user } = response.data

        this.accessToken = access_token
        this.refreshToken = refresh_token
        this.user = user
        this.isAuthenticated = true

        // Store in localStorage
        localStorage.setItem('accessToken', access_token)
        localStorage.setItem('refreshToken', refresh_token)

        // Also set cookies for SSR/middleware compatibility
        if (typeof document !== 'undefined') {
          document.cookie = `accessToken=${access_token}; path=/; max-age=86400; SameSite=Lax`
          document.cookie = `refreshToken=${refresh_token}; path=/; max-age=604800; SameSite=Lax`
        }

        // Reset the logout flag so API calls work properly
        resetLogoutFlag()

        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          error: error.response?.data?.error || 'Login failed',
          requires2FA: error.response?.data?.requires_2fa
        }
      } finally {
        this.isLoading = false
      }
    },

    async register(userData: any) {
      this.isLoading = true
      try {
        const response = await api.post('/auth-service/api/v1/auth/register', userData)
        return { success: true, data: response.data }
      } catch (error: any) {
        return {
          success: false,
          error: error.response?.data?.error || 'Registration failed'
        }
      } finally {
        this.isLoading = false
      }
    },

    async logout() {
      try {
        if (this.accessToken) {
          await api.post('/auth-service/api/v1/auth/logout', {}, {
            headers: { Authorization: `Bearer ${this.accessToken}` }
          })
        }
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.user = null
        this.accessToken = null
        this.refreshToken = null
        this.isAuthenticated = false

        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')

        navigateTo('/auth/login')
      }
    },

    async refreshAccessToken() {
      if (!this.refreshToken) {
        this.logout()
        return false
      }

      try {
        const response = await api.post('/auth-service/api/v1/auth/refresh', {
          refresh_token: this.refreshToken
        })

        const { access_token, refresh_token } = response.data
        this.accessToken = access_token
        this.refreshToken = refresh_token

        localStorage.setItem('accessToken', access_token)
        localStorage.setItem('refreshToken', refresh_token)

        return true
      } catch (error) {
        this.logout()
        return false
      }
    },

    async fetchUserProfile() {
      if (!this.accessToken) return

      try {
        const response = await api.get('/auth-service/api/v1/users/profile', {
          headers: { Authorization: `Bearer ${this.accessToken}` }
        })

        this.user = response.data
        this.isAuthenticated = true
      } catch (error) {
        throw error
      }
    },

    async forgotPassword(email: string) {
      try {
        await api.post('/auth-service/api/v1/auth/forgot-password', { email })
        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          error: error.response?.data?.error || 'Failed to send reset email'
        }
      }
    },

    async resetPassword(token: string, newPassword: string) {
      try {
        await api.post('/auth-service/api/v1/auth/reset-password', {
          token,
          new_password: newPassword
        })
        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          error: error.response?.data?.error || 'Failed to reset password'
        }
      }
    }
  }
})