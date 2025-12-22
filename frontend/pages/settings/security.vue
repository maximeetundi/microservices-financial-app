<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="mb-10">
        <h1 class="text-3xl font-bold text-base mb-2">Sécurité 🔒</h1>
        <p class="text-muted">Protégez votre compte et gérez vos appareils</p>
      </div>

      <!-- Password Section -->
      <div class="glass-card mb-6">
        <div class="flex items-center gap-3 mb-4">
          <span class="text-2xl">🔑</span>
          <h2 class="text-lg font-bold text-base">Mot de passe</h2>
        </div>
        
        <div class="space-y-4">
          <button 
            @click="showPasswordModal = true"
            class="w-full p-4 rounded-xl bg-surface-hover hover:bg-primary/10 transition-colors flex items-center justify-between group"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
                <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
                </svg>
              </div>
              <div class="text-left">
                <p class="font-medium text-base">Changer le mot de passe</p>
                <p class="text-sm text-muted">Dernière modification il y a 30 jours</p>
              </div>
            </div>
            <svg class="w-5 h-5 text-muted group-hover:text-primary transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- 2FA Section -->
      <div class="glass-card mb-6">
        <div class="flex items-center gap-3 mb-4">
          <span class="text-2xl">🛡️</span>
          <h2 class="text-lg font-bold text-base">Authentification à deux facteurs (2FA)</h2>
        </div>
        
        <div class="space-y-4">
          <!-- 2FA Toggle -->
          <div class="p-4 rounded-xl bg-surface-hover flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg" :class="twoFactorEnabled ? 'bg-success/10' : 'bg-secondary-200'">
                <div class="w-full h-full flex items-center justify-center">
                  <svg v-if="twoFactorEnabled" class="w-5 h-5 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                  </svg>
                  <svg v-else class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                  </svg>
                </div>
              </div>
              <div>
                <p class="font-medium text-base">Authentification 2FA</p>
                <p class="text-sm" :class="twoFactorEnabled ? 'text-success' : 'text-muted'">
                  {{ twoFactorEnabled ? 'Activé - Votre compte est protégé' : 'Désactivé - Recommandé pour plus de sécurité' }}
                </p>
              </div>
            </div>
            <button 
              @click="twoFactorEnabled ? disable2FA() : show2FASetup = true"
              class="px-4 py-2 rounded-lg font-medium transition-colors"
              :class="twoFactorEnabled 
                ? 'bg-error/10 text-error hover:bg-error/20' 
                : 'bg-primary text-white hover:bg-primary-hover'"
            >
              {{ twoFactorEnabled ? 'Désactiver' : 'Activer' }}
            </button>
          </div>

          <!-- Recovery Codes -->
          <button 
            v-if="twoFactorEnabled"
            @click="showRecoveryCodes = true"
            class="w-full p-4 rounded-xl bg-surface-hover hover:bg-primary/10 transition-colors flex items-center justify-between group"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-warning/10 flex items-center justify-center">
                <svg class="w-5 h-5 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
              </div>
              <div class="text-left">
                <p class="font-medium text-base">Codes de récupération</p>
                <p class="text-sm text-muted">Télécharger vos codes de secours</p>
              </div>
            </div>
            <svg class="w-5 h-5 text-muted group-hover:text-primary transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Sessions Section -->
      <div class="glass-card mb-6">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-3">
            <span class="text-2xl">📱</span>
            <h2 class="text-lg font-bold text-base">Appareils connectés</h2>
          </div>
          <button 
            @click="revokeAllSessions"
            class="text-sm text-error hover:text-error-hover font-medium"
          >
            Tout déconnecter
          </button>
        </div>
        
        <div v-if="loadingSessions" class="flex justify-center py-8">
          <div class="loading-spinner w-8 h-8"></div>
        </div>

        <div v-else class="space-y-3">
          <div 
            v-for="session in sessions" 
            :key="session.id"
            class="p-4 rounded-xl bg-surface-hover flex items-center justify-between"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="session.is_current ? 'bg-primary/10' : 'bg-secondary-200'"
              >
                <svg v-if="session.device_type === 'mobile'" class="w-5 h-5" :class="session.is_current ? 'text-primary' : 'text-muted'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
                </svg>
                <svg v-else class="w-5 h-5" :class="session.is_current ? 'text-primary' : 'text-muted'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                </svg>
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <p class="font-medium text-base">{{ session.device_name || 'Appareil inconnu' }}</p>
                  <span v-if="session.is_current" class="px-2 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-medium">
                    Actuel
                  </span>
                </div>
                <p class="text-sm text-muted">
                  {{ session.location || 'Localisation inconnue' }} • {{ formatSessionDate(session.last_active) }}
                </p>
              </div>
            </div>
            <button 
              v-if="!session.is_current"
              @click="revokeSession(session.id)"
              class="p-2 rounded-lg hover:bg-error/10 transition-colors group"
            >
              <svg class="w-5 h-5 text-muted group-hover:text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
              </svg>
            </button>
          </div>

          <div v-if="sessions.length === 0" class="text-center py-8 text-muted">
            <p>Aucune session active</p>
          </div>
        </div>
      </div>

      <!-- Login History -->
      <div class="glass-card mb-6">
        <div class="flex items-center gap-3 mb-4">
          <span class="text-2xl">📋</span>
          <h2 class="text-lg font-bold text-base">Historique de connexion</h2>
        </div>
        
        <div class="space-y-2">
          <div 
            v-for="(log, index) in loginHistory" 
            :key="index"
            class="p-3 rounded-lg flex items-center justify-between"
            :class="log.success ? 'bg-surface-hover' : 'bg-error/5'"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-full flex items-center justify-center"
                :class="log.success ? 'bg-success/10' : 'bg-error/10'"
              >
                <svg v-if="log.success" class="w-4 h-4 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <svg v-else class="w-4 h-4 text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium text-base">{{ log.success ? 'Connexion réussie' : 'Tentative échouée' }}</p>
                <p class="text-xs text-muted">{{ log.location }} • {{ log.ip }}</p>
              </div>
            </div>
            <span class="text-xs text-muted">{{ formatLogDate(log.timestamp) }}</span>
          </div>
        </div>
      </div>

      <!-- Danger Zone -->
      <div class="glass-card border-2 border-error/20">
        <div class="flex items-center gap-3 mb-4">
          <span class="text-2xl">⚠️</span>
          <h2 class="text-lg font-bold text-error">Zone Danger</h2>
        </div>
        
        <p class="text-sm text-muted mb-4">Ces actions sont irréversibles. Procédez avec prudence.</p>
        
        <div class="space-y-3">
          <button 
            @click="showDeleteAccountModal = true"
            class="w-full p-4 rounded-xl border-2 border-error/30 hover:bg-error/5 transition-colors flex items-center gap-3 group"
          >
            <svg class="w-5 h-5 text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            <span class="font-medium text-error">Supprimer mon compte</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Change Password Modal -->
    <Teleport to="body">
      <div v-if="showPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="showPasswordModal = false"></div>
        <div class="relative glass-card w-full max-w-md p-6 animate-fade-in-up">
          <h3 class="text-xl font-bold text-base mb-6">Changer le mot de passe</h3>
          
          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-muted mb-2">Mot de passe actuel</label>
              <input 
                v-model="passwordForm.current" 
                type="password" 
                class="input-field w-full" 
                placeholder="••••••••"
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-muted mb-2">Nouveau mot de passe</label>
              <input 
                v-model="passwordForm.new" 
                type="password" 
                class="input-field w-full" 
                placeholder="••••••••"
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-muted mb-2">Confirmer</label>
              <input 
                v-model="passwordForm.confirm" 
                type="password" 
                class="input-field w-full" 
                placeholder="••••••••"
                required
              />
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="showPasswordModal = false" class="flex-1 btn-secondary py-3">
                Annuler
              </button>
              <button type="submit" class="flex-1 btn-premium py-3" :disabled="changingPassword">
                <span v-if="changingPassword" class="loading-spinner w-5 h-5"></span>
                <span v-else>Confirmer</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- 2FA Setup Modal -->
    <Teleport to="body">
      <div v-if="show2FASetup" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="show2FASetup = false"></div>
        <div class="relative glass-card w-full max-w-md p-6 animate-fade-in-up">
          <h3 class="text-xl font-bold text-base mb-4">Configurer 2FA</h3>
          
          <div v-if="setupStep === 1" class="space-y-4">
            <p class="text-muted">Scannez ce QR code avec votre application d'authentification (Google Authenticator, Authy...)</p>
            
            <div class="flex justify-center py-4">
              <div class="w-48 h-48 bg-white rounded-xl flex items-center justify-center">
                <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="QR Code" class="w-40 h-40" />
                <div v-else class="loading-spinner w-8 h-8"></div>
              </div>
            </div>

            <div class="bg-surface-hover p-3 rounded-lg">
              <p class="text-xs text-muted mb-1">Clé secrète (si vous ne pouvez pas scanner)</p>
              <code class="text-sm font-mono text-base break-all">{{ totpSecret }}</code>
            </div>
            
            <button @click="setupStep = 2" class="w-full btn-premium py-3">
              Suivant
            </button>
          </div>

          <div v-else class="space-y-4">
            <p class="text-muted">Entrez le code à 6 chiffres affiché dans votre application</p>
            
            <input 
              v-model="verifyCode" 
              type="text" 
              maxlength="6"
              class="input-field w-full text-center text-2xl tracking-widest" 
              placeholder="000000"
            />
            
            <div class="flex gap-3">
              <button @click="setupStep = 1" class="flex-1 btn-secondary py-3">
                Retour
              </button>
              <button @click="verify2FACode" class="flex-1 btn-premium py-3" :disabled="verifying2FA">
                <span v-if="verifying2FA" class="loading-spinner w-5 h-5"></span>
                <span v-else>Activer</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { userAPI } from '@/composables/useApi'

// State
const twoFactorEnabled = ref(false)
const sessions = ref([])
const loadingSessions = ref(true)
const loginHistory = ref([])

// Password Modal
const showPasswordModal = ref(false)
const changingPassword = ref(false)
const passwordForm = ref({
  current: '',
  new: '',
  confirm: ''
})

// 2FA Setup
const show2FASetup = ref(false)
const setupStep = ref(1)
const qrCodeUrl = ref('')
const totpSecret = ref('')
const verifyCode = ref('')
const verifying2FA = ref(false)
const showRecoveryCodes = ref(false)

// Delete Account
const showDeleteAccountModal = ref(false)

// Fetch data on mount
onMounted(async () => {
  await loadSessions()
  await loadLoginHistory()
  await check2FAStatus()
})

async function loadSessions() {
  loadingSessions.value = true
  try {
    const res = await userAPI.getSessions()
    sessions.value = res.data.sessions || []
  } catch (e) {
    // Demo data fallback
    sessions.value = [
      { id: '1', device_name: 'Chrome - Windows', device_type: 'desktop', location: 'Paris, France', last_active: new Date(), is_current: true },
      { id: '2', device_name: 'iPhone 15 Pro', device_type: 'mobile', location: 'Lyon, France', last_active: new Date(Date.now() - 3600000), is_current: false },
    ]
  } finally {
    loadingSessions.value = false
  }
}

async function loadLoginHistory() {
  // Demo data - in production, fetch from API
  loginHistory.value = [
    { success: true, location: 'Paris, France', ip: '192.168.1.1', timestamp: new Date() },
    { success: true, location: 'Paris, France', ip: '192.168.1.1', timestamp: new Date(Date.now() - 86400000) },
    { success: false, location: 'Londres, UK', ip: '10.0.0.1', timestamp: new Date(Date.now() - 172800000) },
  ]
}

async function check2FAStatus() {
  try {
    const res = await userAPI.getProfile()
    twoFactorEnabled.value = res.data.two_factor_enabled || false
  } catch (e) {
    console.error(e)
  }
}

async function changePassword() {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    alert('Les mots de passe ne correspondent pas')
    return
  }

  changingPassword.value = true
  try {
    await userAPI.changePassword({
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.new
    })
    showPasswordModal.value = false
    passwordForm.value = { current: '', new: '', confirm: '' }
    alert('Mot de passe modifié avec succès!')
  } catch (e) {
    alert('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    changingPassword.value = false
  }
}

async function enable2FA() {
  try {
    const res = await userAPI.enable2FA()
    qrCodeUrl.value = res.data.qr_code_url || ''
    totpSecret.value = res.data.secret || ''
    show2FASetup.value = true
    setupStep.value = 1
  } catch (e) {
    // Demo fallback
    totpSecret.value = 'DEMO123SECRET456'
    show2FASetup.value = true
  }
}

async function verify2FACode() {
  if (verifyCode.value.length !== 6) return

  verifying2FA.value = true
  try {
    await userAPI.verify2FA({ code: verifyCode.value })
    twoFactorEnabled.value = true
    show2FASetup.value = false
    verifyCode.value = ''
    alert('2FA activé avec succès!')
  } catch (e) {
    alert('Code invalide')
  } finally {
    verifying2FA.value = false
  }
}

async function disable2FA() {
  if (!confirm('Êtes-vous sûr de vouloir désactiver 2FA?')) return

  try {
    await userAPI.disable2FA()
    twoFactorEnabled.value = false
    alert('2FA désactivé')
  } catch (e) {
    alert('Erreur: ' + e.message)
  }
}

async function revokeSession(sessionId) {
  try {
    await userAPI.revokeSession(sessionId)
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
  } catch (e) {
    // Remove from UI anyway
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
  }
}

async function revokeAllSessions() {
  if (!confirm('Déconnecter tous les appareils sauf celui-ci?')) return

  try {
    await userAPI.revokeAllSessions()
    sessions.value = sessions.value.filter(s => s.is_current)
    alert('Tous les appareils ont été déconnectés')
  } catch (e) {
    alert('Erreur lors de la déconnexion')
  }
}

function formatSessionDate(date) {
  if (!date) return 'Inconnu'
  const d = new Date(date)
  const diff = Date.now() - d.getTime()
  if (diff < 60000) return 'Maintenant'
  if (diff < 3600000) return `Il y a ${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `Il y a ${Math.floor(diff / 3600000)}h`
  return d.toLocaleDateString('fr-FR')
}

function formatLogDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric', 
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

// State
const twoFactorEnabled = ref(false)
const sessions = ref([])
const loadingSessions = ref(true)
const loginHistory = ref([])

// Password Modal
const showPasswordModal = ref(false)
const changingPassword = ref(false)
const passwordForm = ref({
  current: '',
  new: '',
  confirm: ''
})

// 2FA Setup
const show2FASetup = ref(false)
const setupStep = ref(1)
const qrCodeUrl = ref('')
const totpSecret = ref('')
const verifyCode = ref('')
const verifying2FA = ref(false)
const showRecoveryCodes = ref(false)

// Delete Account
const showDeleteAccountModal = ref(false)

// Fetch data on mount
onMounted(async () => {
  await loadSessions()
  await loadLoginHistory()
  await check2FAStatus()
})

async function loadSessions() {
  loadingSessions.value = true
  try {
    const res = await authApi.getSessions()
    sessions.value = res.data.sessions || []
  } catch (e) {
    // Demo data fallback
    sessions.value = [
      { id: '1', device_name: 'Chrome - Windows', device_type: 'desktop', location: 'Paris, France', last_active: new Date(), is_current: true },
      { id: '2', device_name: 'iPhone 15 Pro', device_type: 'mobile', location: 'Lyon, France', last_active: new Date(Date.now() - 3600000), is_current: false },
    ]
  } finally {
    loadingSessions.value = false
  }
}

async function loadLoginHistory() {
  // Demo data - in production, fetch from API
  loginHistory.value = [
    { success: true, location: 'Paris, France', ip: '192.168.1.1', timestamp: new Date() },
    { success: true, location: 'Paris, France', ip: '192.168.1.1', timestamp: new Date(Date.now() - 86400000) },
    { success: false, location: 'Londres, UK', ip: '10.0.0.1', timestamp: new Date(Date.now() - 172800000) },
  ]
}

async function check2FAStatus() {
  try {
    const res = await authApi.getProfile()
    twoFactorEnabled.value = res.data.two_factor_enabled || false
  } catch (e) {
    console.error(e)
  }
}

async function changePassword() {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    alert('Les mots de passe ne correspondent pas')
    return
  }

  changingPassword.value = true
  try {
    await authApi.changePassword({
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.new
    })
    showPasswordModal.value = false
    passwordForm.value = { current: '', new: '', confirm: '' }
    alert('Mot de passe modifié avec succès!')
  } catch (e) {
    alert('Erreur: ' + (e.response?.data?.error || e.message))
  } finally {
    changingPassword.value = false
  }
}

async function enable2FA() {
  try {
    const res = await authApi.enable2FA()
    qrCodeUrl.value = res.data.qr_code_url || ''
    totpSecret.value = res.data.secret || ''
    show2FASetup.value = true
    setupStep.value = 1
  } catch (e) {
    // Demo fallback
    totpSecret.value = 'DEMO123SECRET456'
    show2FASetup.value = true
  }
}

async function verify2FACode() {
  if (verifyCode.value.length !== 6) return

  verifying2FA.value = true
  try {
    await authApi.verify2FA({ code: verifyCode.value })
    twoFactorEnabled.value = true
    show2FASetup.value = false
    verifyCode.value = ''
    alert('2FA activé avec succès!')
  } catch (e) {
    alert('Code invalide')
  } finally {
    verifying2FA.value = false
  }
}

async function disable2FA() {
  if (!confirm('Êtes-vous sûr de vouloir désactiver 2FA?')) return

  try {
    await authApi.disable2FA()
    twoFactorEnabled.value = false
    alert('2FA désactivé')
  } catch (e) {
    alert('Erreur: ' + e.message)
  }
}

async function revokeSession(sessionId) {
  try {
    await authApi.revokeSession(sessionId)
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
  } catch (e) {
    // Remove from UI anyway
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
  }
}

async function revokeAllSessions() {
  if (!confirm('Déconnecter tous les appareils sauf celui-ci?')) return

  try {
    await authApi.revokeAllSessions()
    sessions.value = sessions.value.filter(s => s.is_current)
    alert('Tous les appareils ont été déconnectés')
  } catch (e) {
    alert('Erreur lors de la déconnexion')
  }
}

function formatSessionDate(date) {
  if (!date) return 'Inconnu'
  const d = new Date(date)
  const diff = Date.now() - d.getTime()
  if (diff < 60000) return 'Maintenant'
  if (diff < 3600000) return `Il y a ${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `Il y a ${Math.floor(diff / 3600000)}h`
  return d.toLocaleDateString('fr-FR')
}

function formatLogDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric', 
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
