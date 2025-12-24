<template>
  <NuxtLayout name="dashboard">
    <div class="profile-page">
      <!-- PIN Not Configured - Force Setup -->
      <div v-if="needsPinSetup" class="pin-overlay">
        <div class="pin-modal">
          <div class="pin-icon">⚠️</div>
          <h2>Configuration requise</h2>
          <p>Pour accéder à vos informations personnelles, vous devez d'abord configurer votre PIN de sécurité à 5 chiffres.</p>
          
          <NuxtLink to="/settings/security" class="setup-btn">
            🔐 Configurer mon PIN
          </NuxtLink>

          <NuxtLink to="/settings" class="back-link">← Retour aux paramètres</NuxtLink>
        </div>
      </div>

      <!-- PIN Verification Required -->
      <div v-else-if="!isUnlocked" class="pin-overlay">
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
const needsPinSetup = ref(false)
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
      // No PIN configured, force user to set it up
      needsPinSetup.value = true
    } else {
      // PIN required, focus first input
      needsPinSetup.value = false
      nextTick(() => pinInputs.value[0]?.focus())
    }
  } catch (e) {
    // If error checking PIN status, require PIN setup to be safe
    needsPinSetup.value = true
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
/* ========== Profile Page Styles ========== */
.profile-page {
  @apply w-full max-w-3xl mx-auto;
}

/* PIN Overlay */
.pin-overlay {
  @apply fixed inset-0 flex items-center justify-center z-50 p-4;
  background: rgba(0,0,0,0.9);
}

.pin-modal {
  @apply rounded-3xl p-8 text-center w-full max-w-sm shadow-2xl;
  background: #1a1a2e; /* Always dark for security modal */
}

.pin-icon {
  @apply text-5xl mb-4;
}

.pin-modal h2 {
  @apply text-2xl font-bold mb-2 text-white;
}

.pin-modal > p {
  @apply text-sm mb-6 text-gray-400;
}

.pin-input-container {
  @apply flex gap-3 justify-center mb-4;
}

.pin-input {
  @apply w-12 h-14 border-2 rounded-xl text-2xl text-center outline-none transition-all;
  border-color: rgba(99, 102, 241, 0.3);
  background: rgba(255,255,255,0.05);
  color: #fff;
}

.pin-input:focus {
  @apply border-indigo-500;
  background: rgba(99, 102, 241, 0.1);
}

.pin-error {
  @apply text-red-500 text-sm mb-4;
}

.verify-btn {
  @apply w-full p-4 rounded-xl border-none text-base font-semibold cursor-pointer mb-4 text-white hover:opacity-90 transition-opacity;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

.verify-btn:disabled {
  @apply opacity-50 cursor-not-allowed;
}

.setup-btn {
  @apply block w-full p-4 rounded-xl border-none text-base font-semibold text-center no-underline mb-4 text-white hover:opacity-90 transition-opacity;
  background: linear-gradient(135deg, #f97316, #ef4444);
}

.back-link {
  @apply text-sm no-underline hover:text-white transition-colors;
  color: #888;
}

/* Header */
.page-header {
  @apply mb-6;
}

.page-header h1 {
  @apply text-2xl font-bold mb-1;
  color: #1e293b;
}

.dark .page-header h1 {
  color: #fff;
}

.page-header > p {
  @apply text-sm m-0;
  color: #64748b;
}

.dark .page-header > p {
  color: #94a3b8;
}

/* Notice Card */
.notice-card {
  @apply flex gap-3 p-4 rounded-2xl mb-6 border;
  background: #eff6ff;
  border-color: #dbeafe;
}

.dark .notice-card {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.2);
}

.notice-icon {
  @apply text-2xl;
}

.notice-card strong {
  @apply text-sm block;
  color: #2563eb;
}

.dark .notice-card strong {
  color: #3b82f6;
}

.notice-card p {
  @apply text-xs mt-1;
  color: #64748b;
}

.dark .notice-card p {
  color: #94a3b8;
}

/* Loading */
.loading {
  @apply text-center py-12;
}

.spinner {
  @apply w-12 h-12 border-4 rounded-full mx-auto mb-4;
  border-color: rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading p {
  color: #64748b;
}

.dark .loading p {
  color: #94a3b8;
}

/* Profile Header Card */
.profile-header-card {
  @apply flex items-center gap-4 p-6 rounded-2xl mb-6 border;
  background: white;
  border-color: #e2e8f0;
}

.dark .profile-header-card {
  background: rgba(255,255,255,0.03);
  border-color: rgba(255,255,255,0.08);
}

.avatar {
  @apply w-16 h-16 rounded-2xl flex items-center justify-center text-2xl font-bold text-white flex-shrink-0;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

.basic-info h2 {
  @apply text-xl font-bold m-0;
  color: #1e293b;
}

.dark .basic-info h2 {
  color: #fff;
}

.basic-info .email {
  @apply text-sm mt-1 mb-2;
  color: #64748b;
}

.dark .basic-info .email {
  color: #94a3b8;
}

.badges {
  @apply flex flex-wrap gap-2;
}

.badge {
  @apply px-2 py-1 rounded-md text-[10px] font-bold uppercase;
}

.badge.verified {
  background: #f0fdf4;
  color: #16a34a;
}

.dark .badge.verified {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.badge.pending {
  background: #fff7ed;
  color: #ea580c;
}

.dark .badge.pending {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

/* Content Sections */
.section {
  @apply p-5 rounded-2xl mb-4 border;
  background: white;
  border-color: #e2e8f0;
}

.dark .section {
  background: rgba(255,255,255,0.03);
  border-color: rgba(255,255,255,0.08);
}

.section h3 {
  @apply text-xs font-bold mb-4 uppercase tracking-wider;
  color: #64748b;
}

.dark .section h3 {
  color: #94a3b8;
}

.info-grid {
  @apply grid grid-cols-2 gap-4;
}

.info-item {
  @apply flex flex-col gap-1;
}

.info-item.full {
  @apply col-span-2;
}

.info-item label {
  @apply text-xs uppercase font-medium;
  color: #64748b;
}

.dark .info-item label {
  color: #94a3b8;
}

.info-item span {
  @apply text-sm font-medium;
  color: #1e293b;
}

.dark .info-item span {
  color: #fff;
}

.kyc-level {
  @apply text-indigo-500 font-bold !important;
}

.status-active {
  @apply text-emerald-500 !important;
}

.status-inactive {
  @apply text-red-500 !important;
}

/* Verifications Grid */
.verifications {
  @apply grid grid-cols-2 gap-3;
}

.verify-item {
  @apply flex items-center gap-2 p-3 rounded-xl border;
  background: #f8fafc;
  border-color: transparent;
}

.dark .verify-item {
  background: rgba(255,255,255,0.03);
}

.verify-item .icon {
  @apply text-xl;
}

.verify-item .label {
  @apply flex-1 text-sm font-medium;
  color: #334155;
}

.dark .verify-item .label {
  color: #fff;
}

.verify-item .status {
  @apply text-[10px] font-bold uppercase px-2 py-1 rounded-md;
}

.verify-item.verified .status {
  background: #f0fdf4;
  color: #16a34a;
}

.dark .verify-item.verified .status {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.verify-item.pending .status {
  background: #f1f5f9;
  color: #64748b;
}

.dark .verify-item.pending .status {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
}

/* Support Section */
.support-section {
  @apply text-center py-6;
}

.support-btn {
  @apply inline-block px-6 py-4 rounded-xl font-bold text-white no-underline shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/30 transition-all hover:-translate-y-0.5;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

.support-section p {
  @apply text-xs mt-3;
  color: #64748b;
}

.dark .support-section p {
  color: #94a3b8;
}

@media (max-width: 480px) {
  .info-grid, .verifications {
    @apply grid-cols-1;
  }
  
  .info-item.full {
    @apply col-span-1;
  }
}
</style>
