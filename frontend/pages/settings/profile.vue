<template>
  <NuxtLayout name="dashboard">
    <div class="profile-page">
      <!-- PIN Verification Modal -->
      <div v-if="!isUnlocked" class="pin-overlay">
        <div class="pin-modal">
          <div class="pin-icon">🔐</div>
          <h2>Vérification requise</h2>
          <p>Entrez votre PIN de sécurité pour accéder à vos informations personnelles</p>
          
          <div class="pin-input-container">
            <input
              v-for="(_, i) in 5"
              :key="i"
              :ref="el => pinInputs[i] = el"
              type="password"
              maxlength="1"
              inputmode="numeric"
              class="pin-input"
              :value="pinDigits[i]"
              @input="handlePinInput($event, i)"
              @keydown="handlePinKeydown($event, i)"
            >
          </div>

          <p v-if="pinError" class="pin-error">{{ pinError }}</p>
          
          <button @click="verifyPin" :disabled="pinDigits.join('').length < 5 || verifyingPin" class="verify-btn">
            {{ verifyingPin ? 'Vérification...' : '🔓 Déverrouiller' }}
          </button>

          <NuxtLink to="/settings" class="back-link">← Retour aux paramètres</NuxtLink>
        </div>
      </div>

      <!-- Profile Content (after PIN verified) -->
      <div v-else>
        <!-- Header -->
        <div class="page-header">
          <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
          <h1>👤 Mon profil</h1>
          <p>Vos informations personnelles</p>
        </div>

        <!-- Security Notice -->
        <div class="notice-card">
          <span class="notice-icon">🔒</span>
          <div>
            <strong>Protection anti-fraude activée</strong>
            <p>La modification des informations nécessite une vérification par le support.</p>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="loading">
          <div class="spinner"></div>
          <p>Chargement de votre profil...</p>
        </div>

        <template v-else>
          <!-- Avatar & Basic Info -->
          <div class="profile-header-card">
            <div class="avatar">{{ userInitials }}</div>
            <div class="basic-info">
              <h2>{{ profile.first_name }} {{ profile.last_name }}</h2>
              <p class="email">{{ profile.email }}</p>
              <div class="badges">
                <span v-if="profile.email_verified" class="badge verified">✓ Email vérifié</span>
                <span v-if="profile.phone_verified" class="badge verified">✓ Téléphone vérifié</span>
                <span v-if="profile.kyc_status === 'verified'" class="badge verified">✓ KYC vérifié</span>
                <span v-else class="badge pending">⏳ KYC en attente</span>
              </div>
            </div>
          </div>

          <!-- Personal Information -->
          <div class="section">
            <h3>👤 INFORMATIONS PERSONNELLES</h3>
            <div class="info-grid">
              <div class="info-item">
                <label>Prénom</label>
                <span>{{ profile.first_name || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Nom</label>
                <span>{{ profile.last_name || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Email</label>
                <span>{{ profile.email || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Téléphone</label>
                <span>{{ profile.phone || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Date de naissance</label>
                <span>{{ formatBirthDate(profile.date_of_birth) }}</span>
              </div>
              <div class="info-item">
                <label>Nationalité</label>
                <span>{{ getCountryName(profile.country) }}</span>
              </div>
            </div>
          </div>

          <!-- Address -->
          <div class="section">
            <h3>🏠 ADRESSE</h3>
            <div class="info-grid">
              <div class="info-item full">
                <label>Adresse</label>
                <span>{{ profile.address || 'Non renseignée' }}</span>
              </div>
              <div class="info-item">
                <label>Ville</label>
                <span>{{ profile.city || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Code postal</label>
                <span>{{ profile.postal_code || '—' }}</span>
              </div>
              <div class="info-item">
                <label>Pays</label>
                <span>{{ getCountryName(profile.country) }}</span>
              </div>
            </div>
          </div>

          <!-- Account Info -->
          <div class="section">
            <h3>📊 INFORMATIONS DU COMPTE</h3>
            <div class="info-grid">
              <div class="info-item">
                <label>Membre depuis</label>
                <span>{{ formatDate(profile.created_at) }}</span>
              </div>
              <div class="info-item">
                <label>Dernière connexion</label>
                <span>{{ formatDate(profile.last_login_at) }}</span>
              </div>
              <div class="info-item">
                <label>Niveau KYC</label>
                <span class="kyc-level">Niveau {{ profile.kyc_level || 0 }}</span>
              </div>
              <div class="info-item">
                <label>Statut du compte</label>
                <span :class="profile.is_active ? 'status-active' : 'status-inactive'">
                  {{ profile.is_active ? '✓ Actif' : '✗ Inactif' }}
                </span>
              </div>
            </div>
          </div>

          <!-- Verification Status -->
          <div class="section">
            <h3>✅ VÉRIFICATIONS</h3>
            <div class="verifications">
              <div class="verify-item" :class="profile.email_verified ? 'verified' : 'pending'">
                <span class="icon">📧</span>
                <span class="label">Email</span>
                <span class="status">{{ profile.email_verified ? 'Vérifié' : 'En attente' }}</span>
              </div>
              <div class="verify-item" :class="profile.phone_verified ? 'verified' : 'pending'">
                <span class="icon">📱</span>
                <span class="label">Téléphone</span>
                <span class="status">{{ profile.phone_verified ? 'Vérifié' : 'En attente' }}</span>
              </div>
              <div class="verify-item" :class="profile.two_fa_enabled ? 'verified' : 'pending'">
                <span class="icon">🛡️</span>
                <span class="label">2FA</span>
                <span class="status">{{ profile.two_fa_enabled ? 'Activé' : 'Désactivé' }}</span>
              </div>
              <div class="verify-item" :class="profile.has_pin ? 'verified' : 'pending'">
                <span class="icon">🔑</span>
                <span class="label">PIN</span>
                <span class="status">{{ profile.has_pin ? 'Configuré' : 'Non configuré' }}</span>
              </div>
            </div>
          </div>

          <!-- Contact Support -->
          <div class="support-section">
            <NuxtLink to="/support" class="support-btn">
              📞 Contacter le support pour modifier mes informations
            </NuxtLink>
            <p>Délai de traitement: 24-48h</p>
          </div>
        </template>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { userAPI } from '~/composables/useApi'

const isUnlocked = ref(false)
const loading = ref(true)
const pinDigits = ref(['', '', '', '', ''])
const pinInputs = ref([])
const pinError = ref('')
const verifyingPin = ref(false)

const profile = ref({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  date_of_birth: null,
  country: '',
  address: '',
  city: '',
  postal_code: '',
  kyc_status: 'pending',
  kyc_level: 0,
  email_verified: false,
  phone_verified: false,
  two_fa_enabled: false,
  has_pin: false,
  is_active: true,
  created_at: null,
  last_login_at: null,
})

const userInitials = computed(() => {
  const first = profile.value.first_name?.[0] || ''
  const last = profile.value.last_name?.[0] || ''
  return (first + last).toUpperCase() || '?'
})

const handlePinInput = (e, index) => {
  const value = e.target.value.replace(/\D/g, '')
  pinDigits.value[index] = value
  pinError.value = ''
  
  if (value && index < 4) {
    nextTick(() => pinInputs.value[index + 1]?.focus())
  }
}

const handlePinKeydown = (e, index) => {
  if (e.key === 'Backspace' && !pinDigits.value[index] && index > 0) {
    nextTick(() => pinInputs.value[index - 1]?.focus())
  }
}

const verifyPin = async () => {
  const pin = pinDigits.value.join('')
  if (pin.length !== 5) return
  
  verifyingPin.value = true
  pinError.value = ''
  
  try {
    const res = await userAPI.verifyPin({ pin })
    if (res.data?.valid) {
      isUnlocked.value = true
      loadProfile()
    } else {
      pinError.value = res.data?.message || 'PIN incorrect'
      pinDigits.value = ['', '', '', '', '']
      nextTick(() => pinInputs.value[0]?.focus())
    }
  } catch (e) {
    pinError.value = 'Erreur de vérification'
    console.error('PIN verification error:', e)
  } finally {
    verifyingPin.value = false
  }
}

const loadProfile = async () => {
  loading.value = true
  try {
    const res = await userAPI.getProfile()
    if (res.data) {
      profile.value = { ...profile.value, ...res.data }
    }
  } catch (e) {
    console.error('Error loading profile:', e)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric',
    month: 'long', 
    year: 'numeric' 
  })
}

const formatBirthDate = (date) => {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric',
    month: 'long', 
    year: 'numeric' 
  })
}

const getCountryName = (code) => {
  if (!code) return '—'
  const countries = {
    'SEN': 'Sénégal', 'CIV': "Côte d'Ivoire", 'FRA': 'France',
    'USA': 'États-Unis', 'GBR': 'Royaume-Uni', 'MLI': 'Mali',
    'BFA': 'Burkina Faso', 'NGA': 'Nigeria', 'GHA': 'Ghana',
    'CMR': 'Cameroun', 'BEN': 'Bénin', 'TGO': 'Togo'
  }
  return countries[code] || code
}

onMounted(async () => {
  // Check if user has PIN configured
  try {
    const pinStatus = await userAPI.checkPinStatus()
    if (!pinStatus.data?.has_pin) {
      // No PIN configured, skip verification
      isUnlocked.value = true
      loadProfile()
    } else {
      // PIN required, focus first input
      nextTick(() => pinInputs.value[0]?.focus())
    }
  } catch (e) {
    // If error, allow access (fallback)
    isUnlocked.value = true
    loadProfile()
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.profile-page {
  width: 100%;
  max-width: 700px;
  margin: 0 auto;
}

.pin-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.pin-modal {
  background: #1a1a2e;
  border-radius: 1.5rem;
  padding: 2rem;
  text-align: center;
  max-width: 400px;
  width: 100%;
}

.pin-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.pin-modal h2 {
  font-size: 1.5rem;
  color: #fff;
  margin: 0 0 0.5rem 0;
}

.pin-modal > p {
  color: #888;
  font-size: 0.875rem;
  margin: 0 0 1.5rem 0;
}

.pin-input-container {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
  margin-bottom: 1rem;
}

.pin-input {
  width: 3rem;
  height: 3.5rem;
  border: 2px solid rgba(99, 102, 241, 0.3);
  border-radius: 0.75rem;
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 1.5rem;
  text-align: center;
  outline: none;
  transition: all 0.2s;
}

.pin-input:focus {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.pin-error {
  color: #ef4444;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.verify-btn {
  width: 100%;
  padding: 1rem;
  border: none;
  border-radius: 0.875rem;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 1rem;
}

.verify-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.back-link {
  color: #888;
  text-decoration: none;
  font-size: 0.875rem;
}

.page-header {
  margin-bottom: 1.5rem;
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

.notice-card {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 1rem;
  margin-bottom: 1.5rem;
}

.notice-icon {
  font-size: 1.5rem;
}

.notice-card strong {
  color: #3b82f6;
  font-size: 0.875rem;
}

.notice-card p {
  color: #888;
  font-size: 0.75rem;
  margin: 0.25rem 0 0 0;
}

.loading {
  text-align: center;
  padding: 3rem;
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading p {
  color: #888;
}

.profile-header-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  margin-bottom: 1.5rem;
}

.avatar {
  width: 4rem;
  height: 4rem;
  border-radius: 1rem;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.basic-info h2 {
  font-size: 1.25rem;
  color: #fff;
  margin: 0;
}

.basic-info .email {
  color: #888;
  font-size: 0.875rem;
  margin: 0.25rem 0 0.5rem 0;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.625rem;
  font-weight: 600;
}

.badge.verified {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.badge.pending {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

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

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-item.full {
  grid-column: span 2;
}

.info-item label {
  font-size: 0.75rem;
  color: #666;
  text-transform: uppercase;
}

.info-item span {
  font-size: 0.9375rem;
  color: #fff;
}

.kyc-level {
  color: #6366f1 !important;
  font-weight: 600;
}

.status-active {
  color: #22c55e !important;
}

.status-inactive {
  color: #ef4444 !important;
}

.verifications {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.verify-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: rgba(255,255,255,0.03);
  border-radius: 0.75rem;
}

.verify-item .icon {
  font-size: 1.25rem;
}

.verify-item .label {
  flex: 1;
  font-size: 0.875rem;
  color: #fff;
}

.verify-item .status {
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
}

.verify-item.verified .status {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.verify-item.pending .status {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
}

.support-section {
  text-align: center;
  padding: 1.5rem 0;
}

.support-btn {
  display: inline-block;
  padding: 1rem 1.5rem;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  text-decoration: none;
  border-radius: 0.875rem;
  font-weight: 600;
}

.support-section p {
  color: #666;
  font-size: 0.75rem;
  margin-top: 0.75rem;
}

@media (max-width: 480px) {
  .info-grid, .verifications {
    grid-template-columns: 1fr;
  }
  
  .info-item.full {
    grid-column: span 1;
  }
}
</style>
