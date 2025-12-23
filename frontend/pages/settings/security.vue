<template>
  <NuxtLayout name="dashboard">
    <div class="security-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>🔒 Sécurité</h1>
        <p>Protégez votre compte et gérez vos appareils</p>
      </div>

      <!-- Security Score -->
      <div class="score-card">
        <div class="score-circle" :class="scoreClass">
          <span class="score-value">{{ securityScore }}%</span>
        </div>
        <div class="score-info">
          <h3>Niveau de sécurité: {{ scoreLabel }}</h3>
          <p>{{ scoreDescription }}</p>
        </div>
      </div>

      <!-- Password Section -->
      <div class="section">
        <h3>🔑 MOT DE PASSE</h3>
        <button @click="showPasswordModal = true" class="action-btn">
          <span class="icon">🔐</span>
          <div class="content">
            <strong>Changer le mot de passe</strong>
            <span>Dernière modification: {{ passwordLastChanged }}</span>
          </div>
          <span class="arrow">→</span>
        </button>
      </div>

      <!-- PIN Section -->
      <div class="section">
        <h3>🔢 PIN DE SÉCURITÉ</h3>
        <div class="action-btn" :class="hasPin ? 'configured' : 'not-configured'">
          <span class="icon">{{ hasPin ? '✓' : '⚠️' }}</span>
          <div class="content">
            <strong>PIN à 5 chiffres</strong>
            <span>{{ hasPin ? 'Configuré - requis pour les opérations sensibles' : 'Non configuré' }}</span>
          </div>
          <button v-if="hasPin" @click="showChangePinModal = true" class="small-btn">Modifier</button>
          <button v-else @click="showSetupPinModal = true" class="small-btn primary">Configurer</button>
        </div>
      </div>

      <!-- 2FA Section -->
      <div class="section">
        <h3>🛡️ AUTHENTIFICATION À DEUX FACTEURS (2FA)</h3>
        
        <div class="twofa-card" :class="twoFactorEnabled ? 'enabled' : 'disabled'">
          <div class="twofa-header">
            <span class="icon">{{ twoFactorEnabled ? '✅' : '🔓' }}</span>
            <div class="info">
              <strong>Application d'authentification (TOTP)</strong>
              <span>{{ twoFactorEnabled ? 'Activé - Votre compte est protégé' : 'Désactivé - Recommandé' }}</span>
            </div>
          </div>
          
          <button 
            @click="twoFactorEnabled ? showDisable2FAModal = true : start2FASetup()"
            class="twofa-btn"
            :class="twoFactorEnabled ? 'danger' : 'primary'"
          >
            {{ twoFactorEnabled ? 'Désactiver' : 'Activer 2FA' }}
          </button>
        </div>

        <div v-if="twoFactorEnabled" class="recovery-link">
          <button @click="showRecoveryCodes = true" class="action-btn small">
            <span class="icon">📄</span>
            <div class="content">
              <strong>Codes de récupération</strong>
              <span>Télécharger vos codes de secours</span>
            </div>
            <span class="arrow">→</span>
          </button>
        </div>
      </div>

      <!-- Active Sessions -->
      <div class="section">
        <div class="section-header">
          <h3>📱 APPAREILS CONNECTÉS</h3>
          <button @click="revokeAllSessions" class="text-danger small">Tout déconnecter</button>
        </div>
        
        <div v-if="loadingSessions" class="loading">
          <div class="spinner"></div>
        </div>
        
        <div v-else class="sessions-list">
          <div v-for="session in sessions" :key="session.id" class="session-item enhanced">
            <div class="device-icon">{{ getDeviceIcon(session.device_type) }}</div>
            <div class="session-info">
              <div class="session-header">
                <span class="session-name">
                  {{ session.browser || 'Navigateur inconnu' }}
                  <span v-if="session.is_current" class="current-badge">Actuel</span>
                </span>
              </div>
              <div class="session-details">
                <span class="detail-item">
                  <span class="icon">💻</span> {{ session.os || 'OS inconnu' }}
                </span>
                <span class="detail-item">
                  <span class="icon">🌐</span> {{ session.ip_address || 'IP inconnue' }}
                </span>
              </div>
              <div class="session-meta">
                {{ session.location || 'Localisation inconnue' }} • {{ formatSessionDate(session.last_active || session.created_at) }}
              </div>
            </div>
            <div class="session-actions">
              <button @click="showSessionDetails(session)" class="info-btn" title="Voir les détails">ℹ️</button>
              <button v-if="!session.is_current" @click="revokeSession(session.id)" class="revoke-btn" title="Déconnecter cet appareil">✕</button>
            </div>
          </div>
          
          <div v-if="sessions.length === 0" class="empty">Aucune session active</div>
        </div>
      </div>

      <!-- Session Details Modal -->
      <div v-if="selectedSession" class="modal-overlay" @click="selectedSession = null">
        <div class="modal-content session-modal" @click.stop>
          <div class="modal-header">
            <h3>📱 Détails de l'appareil</h3>
            <button @click="selectedSession = null" class="close-btn">✕</button>
          </div>
          
          <div class="session-detail-grid">
            <div class="detail-row">
              <span class="label">Appareil</span>
              <span class="value">{{ getDeviceIcon(selectedSession.device_type) }} {{ selectedSession.device_name || 'Inconnu' }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Système</span>
              <span class="value">{{ selectedSession.os || 'Inconnu' }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Navigateur</span>
              <span class="value">{{ selectedSession.browser || 'Inconnu' }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Adresse IP</span>
              <span class="value code">{{ selectedSession.ip_address || 'Inconnue' }}</span>
            </div>
             <div class="detail-row">
              <span class="label">Localisation</span>
              <span class="value">{{ selectedSession.location || 'Inconnue' }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Première connexion</span>
              <span class="value">{{ formatSessionDate(selectedSession.created_at) }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Dernière activité</span>
              <span class="value">{{ formatSessionDate(selectedSession.last_active) }}</span>
            </div>
            <div class="detail-row">
              <span class="label">Statut</span>
              <span class="value status-badge" :class="{ 'current': selectedSession.is_current }">
                {{ selectedSession.is_current ? '🟢 Session Actuelle' : '⚪ Connecté' }}
              </span>
            </div>
          </div>

          <div class="modal-actions" v-if="!selectedSession.is_current">
            <button @click="revokeSession(selectedSession.id); selectedSession = null" class="btn-danger full-width">
              Déconnecter cet appareil
            </button>
          </div>
        </div>
      </div>

      <!-- Change Password Modal -->
      <div v-if="showPasswordModal" class="modal-overlay" @click="showPasswordModal = false">
        <div class="modal-content" @click.stop>
          <h3>🔐 Changer le mot de passe</h3>
          <form @submit.prevent="changePassword">
            <div class="form-group">
              <label>Mot de passe actuel</label>
              <input v-model="passwordForm.current" type="password" placeholder="••••••••" required>
            </div>
            <div class="form-group">
              <label>Nouveau mot de passe</label>
              <input v-model="passwordForm.new" type="password" placeholder="••••••••" required>
            </div>
            <div class="form-group">
              <label>Confirmer</label>
              <input v-model="passwordForm.confirm" type="password" placeholder="••••••••" required>
            </div>
            <p v-if="passwordError" class="error">{{ passwordError }}</p>
            <div class="modal-actions">
              <button type="button" @click="showPasswordModal = false" class="btn-cancel">Annuler</button>
              <button type="submit" :disabled="changingPassword" class="btn-confirm">
                {{ changingPassword ? 'Modification...' : 'Confirmer' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 2FA Setup Modal -->
      <div v-if="show2FASetup" class="modal-overlay" @click="show2FASetup = false">
        <div class="modal-content twofa-modal" @click.stop>
          <h3>🛡️ Configurer l'authentification 2FA</h3>
          
          <!-- Step 1: Show QR Code -->
          <div v-if="setupStep === 1" class="setup-step">
            <p>Scannez ce QR code avec une application d'authentification:</p>
            <div class="app-icons">
              <span title="Google Authenticator">🔐</span>
              <span title="Authy">Authy</span>
              <span title="Microsoft Authenticator">📱</span>
            </div>
            
            <div class="qr-container">
              <div v-if="loadingQR" class="spinner"></div>
              <img v-else-if="qrCodeUrl" :src="qrCodeUrl" alt="QR Code" class="qr-code">
              <div v-else-if="qrError" class="qr-error">
                <span>⚠️</span>
                <p>{{ qrError }}</p>
                <button @click="start2FASetup" class="retry-btn">Réessayer</button>
              </div>
              <div v-else class="qr-loading">
                <div class="spinner"></div>
              </div>
            </div>

            <div v-if="totpSecret" class="secret-box">
              <label>Clé secrète (entrée manuelle):</label>
              <code>{{ totpSecret }}</code>
              <button @click="copySecret" class="copy-btn">📋 Copier</button>
            </div>

            <button @click="setupStep = 2" :disabled="!totpSecret" class="btn-next">Suivant →</button>
          </div>

          <!-- Step 2: Verify Code -->
          <div v-if="setupStep === 2" class="setup-step">
            <p>Entrez le code à 6 chiffres affiché dans l'application:</p>
            
            <div class="code-input-container">
              <input
                v-for="(_, i) in 6"
                :key="i"
                :ref="el => codeInputs[i] = el"
                type="text"
                maxlength="1"
                inputmode="numeric"
                class="code-input"
                :value="verifyCode[i] || ''"
                @input="handleCodeInput($event, i)"
                @keydown="handleCodeKeydown($event, i)"
              >
            </div>

            <p v-if="verifyError" class="error">{{ verifyError }}</p>

            <div class="modal-actions">
              <button @click="setupStep = 1" class="btn-cancel">← Retour</button>
              <button @click="verify2FACode" :disabled="verifying2FA || verifyCode.length < 6" class="btn-confirm">
                {{ verifying2FA ? 'Vérification...' : '✓ Activer 2FA' }}
              </button>
            </div>
          </div>

          <!-- Step 3: Success + Recovery Codes -->
          <div v-if="setupStep === 3" class="setup-step success-step">
            <div class="success-icon">✅</div>
            <h4>2FA Activé avec succès!</h4>
            <p>Sauvegardez ces codes de récupération en lieu sûr. Ils vous permettront d'accéder à votre compte si vous perdez votre téléphone.</p>
            
            <div class="recovery-codes">
              <code v-for="code in recoveryCodes" :key="code">{{ code }}</code>
            </div>

            <div class="modal-actions">
              <button @click="downloadRecoveryCodes" class="btn-cancel">📥 Télécharger</button>
              <button @click="finish2FASetup" class="btn-confirm">Terminer</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Setup PIN Modal -->
      <div v-if="showSetupPinModal" class="modal-overlay" @click="showSetupPinModal = false">
        <div class="modal-content" @click.stop>
          <h3>🔢 Configurer votre PIN</h3>
          <p>Créez un PIN à 5 chiffres pour sécuriser les opérations sensibles.</p>
          
          <div class="pin-form">
            <label>Nouveau PIN</label>
            <div class="pin-inputs" ref="newPinInputs">
              <input
                v-for="(_, i) in 5"
                :key="'new-'+i"
                type="password"
                maxlength="1"
                inputmode="numeric"
                :value="newPin[i] || ''"
                @input="handleNewPinInput($event, i)"
                @keydown="handlePinKeydown($event, i, 'new')"
              >
            </div>
            
            <label>Confirmer le PIN</label>
            <div class="pin-inputs" ref="confirmPinInputs">
              <input
                v-for="(_, i) in 5"
                :key="'confirm-'+i"
                type="password"
                maxlength="1"
                inputmode="numeric"
                :value="confirmPin[i] || ''"
                @input="handleConfirmPinInput($event, i)"
                @keydown="handlePinKeydown($event, i, 'confirm')"
              >
            </div>
          </div>

          <p v-if="pinError" class="error">{{ pinError }}</p>

          <div class="modal-actions">
            <button @click="showSetupPinModal = false" class="btn-cancel">Annuler</button>
            <button @click="setupPin" :disabled="settingPin" class="btn-confirm">
              {{ settingPin ? 'Configuration...' : 'Configurer' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { userAPI } from '~/composables/useApi'

// Security Status
const hasPin = ref(false)
const twoFactorEnabled = ref(false)
const sessions = ref([])
const selectedSession = ref(null)
const loadingSessions = ref(true)
const passwordLastChanged = ref('il y a 30 jours')

function showSessionDetails(session) {
  selectedSession.value = session
}

// Computed Security Score
const securityScore = computed(() => {
  let score = 30 // Base score
  if (hasPin.value) score += 25
  if (twoFactorEnabled.value) score += 35
  if (sessions.value.length <= 2) score += 10
  return Math.min(score, 100)
})

const scoreClass = computed(() => {
  if (securityScore.value >= 80) return 'excellent'
  if (securityScore.value >= 50) return 'good'
  return 'weak'
})

const scoreLabel = computed(() => {
  if (securityScore.value >= 80) return 'Excellent'
  if (securityScore.value >= 50) return 'Bon'
  return 'Faible'
})

const scoreDescription = computed(() => {
  if (securityScore.value >= 80) return 'Votre compte est bien protégé'
  if (securityScore.value >= 50) return 'Activez 2FA pour améliorer la sécurité'
  return 'Configurez PIN et 2FA pour sécuriser votre compte'
})

// Password Modal
const showPasswordModal = ref(false)
const changingPassword = ref(false)
const passwordError = ref('')
const passwordForm = ref({ current: '', new: '', confirm: '' })

// PIN Modal
const showSetupPinModal = ref(false)
const showChangePinModal = ref(false)
const settingPin = ref(false)
const pinError = ref('')
const newPin = ref('')
const confirmPin = ref('')
const newPinInputs = ref(null)
const confirmPinInputs = ref(null)

// PIN Input Handlers - auto-advance to next field
function handleNewPinInput(event, index) {
  const value = event.target.value.replace(/\D/g, '') // Only digits
  newPin.value = newPin.value.substring(0, index) + value + newPin.value.substring(index + 1)
  
  // Auto-advance to next input
  if (value && index < 4 && newPinInputs.value) {
    const inputs = newPinInputs.value.querySelectorAll('input')
    if (inputs[index + 1]) {
      nextTick(() => inputs[index + 1].focus())
    }
  }
  // Move to confirm PIN inputs when newPin is complete
  else if (value && index === 4 && confirmPinInputs.value) {
    const inputs = confirmPinInputs.value.querySelectorAll('input')
    if (inputs[0]) {
      nextTick(() => inputs[0].focus())
    }
  }
}

function handleConfirmPinInput(event, index) {
  const value = event.target.value.replace(/\D/g, '') // Only digits
  confirmPin.value = confirmPin.value.substring(0, index) + value + confirmPin.value.substring(index + 1)
  
  // Auto-advance to next input
  if (value && index < 4 && confirmPinInputs.value) {
    const inputs = confirmPinInputs.value.querySelectorAll('input')
    if (inputs[index + 1]) {
      nextTick(() => inputs[index + 1].focus())
    }
  }
}

function handlePinKeydown(event, index, type) {
  // Handle backspace - go to previous input
  if (event.key === 'Backspace') {
    const container = type === 'new' ? newPinInputs.value : confirmPinInputs.value
    const currentValue = type === 'new' ? newPin.value[index] : confirmPin.value[index]
    
    if (!currentValue && index > 0 && container) {
      const inputs = container.querySelectorAll('input')
      if (inputs[index - 1]) {
        nextTick(() => inputs[index - 1].focus())
      }
    }
  }
}

// 2FA Setup
const show2FASetup = ref(false)
const showDisable2FAModal = ref(false)
const showRecoveryCodes = ref(false)
const setupStep = ref(1)
const loadingQR = ref(false)
const qrCodeUrl = ref('')
const qrError = ref('')
const totpSecret = ref('')
const verifyCode = ref('')
const verifyError = ref('')
const verifying2FA = ref(false)
const codeInputs = ref([])
const recoveryCodes = ref([])

// Load data on mount
onMounted(async () => {
  await loadSecurityStatus()
  await loadSessions()
})

async function loadSecurityStatus() {
  try {
    const [profileRes, pinRes] = await Promise.all([
      userAPI.getProfile(),
      userAPI.checkPinStatus()
    ])
    
    twoFactorEnabled.value = profileRes.data?.two_fa_enabled || false
    hasPin.value = pinRes.data?.has_pin || false
  } catch (e) {
    console.error('Error loading security status:', e)
  }
}

async function loadSessions() {
  loadingSessions.value = true
  try {
    const res = await userAPI.getSessions()
    sessions.value = res.data?.sessions || []
  } catch (e) {
    // Demo fallback
    sessions.value = [
      { id: '1', device_name: 'Chrome - Windows', device_type: 'desktop', location: 'Paris', last_active: new Date(), is_current: true }
    ]
  } finally {
    loadingSessions.value = false
  }
}

async function changePassword() {
  passwordError.value = ''
  
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    passwordError.value = 'Les mots de passe ne correspondent pas'
    return
  }
  
  if (passwordForm.value.new.length < 8) {
    passwordError.value = 'Le mot de passe doit contenir au moins 8 caractères'
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
    passwordError.value = e.response?.data?.error || 'Erreur lors de la modification'
  } finally {
    changingPassword.value = false
  }
}

async function setupPin() {
  pinError.value = ''
  
  if (newPin.value.length !== 5 || confirmPin.value.length !== 5) {
    pinError.value = 'Le PIN doit contenir 5 chiffres'
    return
  }
  
  if (newPin.value !== confirmPin.value) {
    pinError.value = 'Les PINs ne correspondent pas'
    return
  }

  settingPin.value = true
  try {
    await userAPI.setupPin({ pin: newPin.value, confirm_pin: confirmPin.value })
    hasPin.value = true
    showSetupPinModal.value = false
    newPin.value = ''
    confirmPin.value = ''
    alert('PIN configuré avec succès!')
  } catch (e) {
    pinError.value = e.response?.data?.error || 'Erreur lors de la configuration'
  } finally {
    settingPin.value = false
  }
}

async function start2FASetup() {
  show2FASetup.value = true
  setupStep.value = 1
  loadingQR.value = true
  verifyCode.value = ''
  verifyError.value = ''
  qrError.value = ''
  qrCodeUrl.value = ''
  totpSecret.value = ''
  
  try {
    const res = await userAPI.enable2FA()
    // Backend may return qr_code or qr_code_url
    qrCodeUrl.value = res.data?.qr_code || res.data?.qr_code_url || ''
    totpSecret.value = res.data?.secret || ''
    
    if (!qrCodeUrl.value && !totpSecret.value) {
      qrError.value = 'Erreur lors de la génération du QR code'
    }
  } catch (e) {
    console.error('2FA setup error:', e)
    qrError.value = e.response?.data?.error || 'Erreur de connexion au serveur'
  } finally {
    loadingQR.value = false
  }
}

function handleCodeInput(e, index) {
  const value = e.target.value.replace(/\D/g, '')
  const newCode = verifyCode.value.split('')
  newCode[index] = value
  verifyCode.value = newCode.join('')
  
  if (value && index < 5) {
    nextTick(() => codeInputs.value[index + 1]?.focus())
  }
}

function handleCodeKeydown(e, index) {
  if (e.key === 'Backspace' && !verifyCode.value[index] && index > 0) {
    nextTick(() => codeInputs.value[index - 1]?.focus())
  }
}

async function verify2FACode() {
  if (verifyCode.value.length !== 6) return
  if (!totpSecret.value) {
    verifyError.value = 'Erreur: clé secrète non disponible'
    return
  }
  
  verifying2FA.value = true
  verifyError.value = ''
  
  try {
    const res = await userAPI.verify2FA({ 
      code: verifyCode.value,
      secret: totpSecret.value  // Pass secret to backend for storage
    })
    recoveryCodes.value = res.data?.recovery_codes || []
    setupStep.value = 3
    twoFactorEnabled.value = true
  } catch (e) {
    verifyError.value = e.response?.data?.error || 'Code invalide. Veuillez réessayer.'
  } finally {
    verifying2FA.value = false
  }
}

function copySecret() {
  navigator.clipboard.writeText(totpSecret.value)
  alert('Clé copiée!')
}

function downloadRecoveryCodes() {
  const content = 'Codes de récupération 2FA\n\n' + recoveryCodes.value.join('\n')
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'recovery-codes.txt'
  a.click()
}

function finish2FASetup() {
  show2FASetup.value = false
  setupStep.value = 1
}

async function revokeSession(sessionId) {
  try {
    await userAPI.revokeSession(sessionId)
  } catch (e) {}
  sessions.value = sessions.value.filter(s => s.id !== sessionId)
}

async function revokeAllSessions() {
  if (!confirm('Déconnecter tous les appareils sauf celui-ci?')) return
  try {
    await userAPI.revokeAllSessions()
  } catch (e) {}
  sessions.value = sessions.value.filter(s => s.is_current)
}

function formatSessionDate(date) {
  if (!date) return 'Inconnu'
  const d = new Date(date)
  const diff = Date.now() - d.getTime()
  if (diff < 60000) return 'Maintenant'
  if (diff < 3600000) return `Il y a ${Math.floor(diff / 60000)} min`
  if (diff < 86400000) return `Il y a ${Math.floor(diff / 3600000)}h`
  return d.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getDeviceIcon(type) {
  switch (type) {
    case 'mobile': return '📱'
    case 'tablet': return '📲'
    case 'desktop': return '💻'
    default: return '❓'
  }
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.security-page {
  width: 100%;
  max-width: 700px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 1.5rem;
}

.back-link {
  color: #888;
  text-decoration: none;
  font-size: 0.875rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0.5rem 0 0.25rem 0;
}

.page-header > p {
  color: #888;
  font-size: 0.875rem;
  margin: 0;
}

/* Score Card */
.score-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  margin-bottom: 1.5rem;
}

.score-circle {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.score-circle.excellent { background: rgba(34, 197, 94, 0.2); border: 3px solid #22c55e; }
.score-circle.good { background: rgba(249, 115, 22, 0.2); border: 3px solid #f97316; }
.score-circle.weak { background: rgba(239, 68, 68, 0.2); border: 3px solid #ef4444; }

.score-value {
  font-size: 1rem;
  font-weight: 700;
  color: #fff;
}

.score-info h3 {
  font-size: 1rem;
  color: #fff;
  margin: 0;
}

.score-info p {
  font-size: 0.75rem;
  color: #888;
  margin: 0.25rem 0 0 0;
}

/* Sections */
.section {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  padding: 1.25rem;
  margin-bottom: 1rem;
}

.section h3 {
  font-size: 0.75rem;
  font-weight: 600;
  color: #888;
  margin: 0 0 1rem 0;
  letter-spacing: 0.05em;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.section-header h3 {
  margin: 0;
}

.text-danger {
  color: #ef4444;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.75rem;
}

/* Action Button */
.action-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 0.875rem;
  color: #fff;
  cursor: pointer;
  text-align: left;
}

.action-btn .icon {
  font-size: 1.5rem;
}

.action-btn .content {
  flex: 1;
}

.action-btn .content strong {
  display: block;
  font-size: 0.9375rem;
}

.action-btn .content span {
  font-size: 0.75rem;
  color: #888;
}

.action-btn .arrow {
  color: #888;
}

.action-btn.configured {
  border-color: rgba(34, 197, 94, 0.3);
}

.action-btn.not-configured {
  border-color: rgba(249, 115, 22, 0.3);
}

.small-btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  border: none;
  background: rgba(255,255,255,0.1);
  color: #fff;
  font-size: 0.75rem;
  cursor: pointer;
}

.small-btn.primary {
  background: #6366f1;
}

/* 2FA Card */
.twofa-card {
  padding: 1rem;
  border-radius: 0.875rem;
}

.twofa-card.enabled {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.twofa-card.disabled {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
}

.twofa-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.twofa-header .icon {
  font-size: 1.5rem;
}

.twofa-header .info strong {
  display: block;
  color: #fff;
  font-size: 0.9375rem;
}

.twofa-header .info span {
  font-size: 0.75rem;
  color: #888;
}

.twofa-btn {
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 0.625rem;
  font-weight: 600;
  cursor: pointer;
}

.twofa-btn.primary {
  background: #6366f1;
  color: #fff;
}

.twofa-btn.danger {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.recovery-link {
  margin-top: 0.75rem;
}

/* Sessions */
.loading {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.sessions-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: rgba(255,255,255,0.03);
  border-radius: 0.625rem;
}

.device-icon {
  font-size: 1.25rem;
}

.session-info {
  flex: 1;
}

.session-name {
  font-size: 0.875rem;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.current-badge {
  font-size: 0.625rem;
  padding: 0.125rem 0.375rem;
  background: rgba(99, 102, 241, 0.2);
  color: #6366f1;
  border-radius: 0.25rem;
}

.session-meta {
  font-size: 0.75rem;
  color: #888;
}

.revoke-btn {
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border-radius: 0.375rem;
  cursor: pointer;
}

.empty {
  text-align: center;
  color: #888;
  padding: 1rem;
}

/* Modals */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.modal-content {
  background: #1a1a2e;
  border-radius: 1rem;
  padding: 1.5rem;
  max-width: 420px;
  width: 100%;
}

.modal-content h3 {
  font-size: 1.25rem;
  color: #fff;
  margin: 0 0 1rem 0;
}

.modal-content p {
  color: #888;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.75rem;
  color: #888;
  margin-bottom: 0.375rem;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 0.625rem;
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
}

.error {
  color: #ef4444;
  font-size: 0.75rem;
  margin-bottom: 1rem;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

.btn-cancel, .btn-confirm {
  flex: 1;
  padding: 0.75rem;
  border-radius: 0.625rem;
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel {
  background: rgba(255,255,255,0.1);
  color: #fff;
}

.btn-confirm {
  background: #6366f1;
  color: #fff;
}

.btn-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 2FA Setup Modal */
.twofa-modal {
  max-width: 450px;
}

.setup-step {
  text-align: center;
}

.qr-container {
  display: flex;
  justify-content: center;
  margin: 1.5rem 0;
}

.qr-code {
  width: 180px;
  height: 180px;
  border-radius: 0.75rem;
}

.qr-placeholder {
  width: 180px;
  height: 180px;
  background: rgba(255,255,255,0.05);
  border-radius: 0.75rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.qr-placeholder span {
  font-size: 3rem;
}

.qr-placeholder p {
  margin: 0.5rem 0 0 0;
}

.secret-box {
  background: rgba(255,255,255,0.05);
  padding: 1rem;
  border-radius: 0.75rem;
  margin: 1rem 0;
  text-align: left;
}

.secret-box label {
  display: block;
  font-size: 0.75rem;
  color: #888;
  margin-bottom: 0.25rem;
}

.secret-box code {
  display: block;
  color: #6366f1;
  font-size: 0.875rem;
  word-break: break-all;
}

.copy-btn {
  margin-top: 0.5rem;
  padding: 0.375rem 0.75rem;
  background: rgba(99, 102, 241, 0.15);
  color: #6366f1;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  cursor: pointer;
}

.btn-next {
  width: 100%;
  padding: 0.875rem;
  background: #6366f1;
  color: #fff;
  border: none;
  border-radius: 0.625rem;
  font-weight: 600;
  cursor: pointer;
}

.code-input-container {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin: 1.5rem 0;
}

.code-input {
  width: 2.5rem;
  height: 3rem;
  border: 2px solid rgba(99, 102, 241, 0.3);
  border-radius: 0.5rem;
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 1.25rem;
  text-align: center;
  outline: none;
}

.code-input:focus {
  border-color: #6366f1;
}

.success-step .success-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.success-step h4 {
  color: #22c55e;
  margin: 0 0 0.5rem 0;
}

.recovery-codes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
  margin: 1rem 0;
}

.recovery-codes code {
  padding: 0.5rem;
  background: rgba(255,255,255,0.05);
  border-radius: 0.375rem;
  color: #fff;
  font-size: 0.75rem;
}

/* PIN Form */
.pin-form {
  margin: 1rem 0;
}

.pin-form label {
  display: block;
  font-size: 0.75rem;
  color: #888;
  margin-bottom: 0.5rem;
}

.pin-inputs {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 1rem;
}

.pin-inputs input {
  width: 2.5rem;
  height: 3rem;
  border: 2px solid rgba(255,255,255,0.1);
  border-radius: 0.5rem;
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 1.25rem;
  text-align: center;
  outline: none;
}

.pin-inputs input:focus {
  border-color: #6366f1;
}

/* App icons for 2FA setup */
.app-icons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin: 0.75rem 0;
  font-size: 0.875rem;
  color: #888;
}

.app-icons span {
  padding: 0.25rem 0.5rem;
  background: rgba(255,255,255,0.05);
  border-radius: 0.375rem;
}

/* QR error state */
.qr-error {
  width: 180px;
  height: 180px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.75rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 1rem;
}

.qr-error span {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.qr-error p {
  font-size: 0.75rem;
  color: #ef4444;
  margin: 0 0 0.75rem 0;
}

.retry-btn {
  padding: 0.375rem 0.75rem;
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  cursor: pointer;
}

.qr-loading {
  width: 180px;
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.05);
  border-radius: 0.75rem;
}

/* Enhanced Session List */
.session-item.enhanced {
  align-items: flex-start;
  padding: 1.25rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.05);
  border-radius: 1rem;
  margin-bottom: 0.75rem;
}

.session-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.session-header .session-name {
  font-size: 1rem;
  margin: 0;
}

.session-details {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.8rem;
  color: #94a3b8;
  background: rgba(255,255,255,0.05);
  padding: 0.25rem 0.6rem;
  border-radius: 0.5rem;
}

.detail-item .icon {
  font-size: 0.9rem;
  opacity: 0.7;
}

.session-meta {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 0.25rem;
}

.session-actions {
  display: flex;
  gap: 0.5rem;
}

.info-btn {
  background: rgba(255,255,255,0.05);
  border: none;
  color: #fff;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s;
}

.info-btn:hover {
  background: rgba(255,255,255,0.1);
  transform: scale(1.1);
}

.session-modal {
  max-width: 500px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
}

.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 1.25rem;
  cursor: pointer;
}

.session-detail-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: rgba(255,255,255,0.03);
  border-radius: 0.5rem;
}

.detail-row .label {
  color: #94a3b8;
  font-size: 0.9rem;
}

.detail-row .value {
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.detail-row .value.code {
  font-family: monospace;
  background: rgba(0,0,0,0.3);
  padding: 0.2rem 0.5rem;
  border-radius: 0.25rem;
}

.status-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 1rem;
  font-size: 0.8rem;
  background: rgba(255,255,255,0.1);
}

.status-badge.current {
  background: rgba(34, 197, 94, 0.2);
  color: #4ade80;
}

.full-width {
  width: 100%;
  padding: 0.75rem;
  font-weight: 600;
  display: block;
}
</style>
