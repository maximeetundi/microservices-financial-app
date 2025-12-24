<template>
  <NuxtLayout name="dashboard">
    <div class="settings-page">
      <!-- Header -->
      <div class="page-header">
        <h1>⚙️ Paramètres</h1>
        <p>Gérez votre compte et vos préférences</p>
      </div>

      <!-- Settings Grid -->
      <div class="settings-grid">
        <!-- Profile -->
        <NuxtLink to="/settings/profile" class="settings-card">
          <div class="card-icon blue">👤</div>
          <div class="card-content">
            <h3>Profil</h3>
            <p>Informations personnelles, coordonnées</p>
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>

        <!-- Security -->
        <NuxtLink to="/settings/security" class="settings-card">
          <div class="card-icon green">🔒</div>
          <div class="card-content">
            <h3>Sécurité</h3>
            <p>Mot de passe, 2FA, sessions</p>
          </div>
          <div class="card-status" :class="{ active: securityScore >= 80 }">
            {{ securityScore }}%
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>

        <!-- KYC -->
        <NuxtLink to="/settings/kyc" class="settings-card">
          <div class="card-icon purple">📋</div>
          <div class="card-content">
            <h3>Vérification KYC</h3>
            <p>Documents d'identité, validation</p>
          </div>
          <div class="card-badge" :class="kycStatus.class">
            {{ kycStatus.label }}
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>

        <!-- Preferences -->
        <NuxtLink to="/settings/preferences" class="settings-card">
          <div class="card-icon orange">🎨</div>
          <div class="card-content">
            <h3>Préférences</h3>
            <p>Thème, langue, notifications</p>
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>

        <!-- Notifications -->
        <NuxtLink to="/settings/notifications" class="settings-card">
          <div class="card-icon pink">🔔</div>
          <div class="card-content">
            <h3>Notifications</h3>
            <p>Alertes email, push, SMS</p>
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>

        <!-- Payment Methods -->
        <NuxtLink to="/settings/payment-methods" class="settings-card">
          <div class="card-icon teal">💳</div>
          <div class="card-content">
            <h3>Moyens de paiement</h3>
            <p>Cartes, comptes bancaires</p>
          </div>
          <div class="card-arrow">→</div>
        </NuxtLink>
      </div>

      <!-- Quick Actions -->
      <div class="quick-section">
        <h2>Actions rapides</h2>
        <div class="quick-actions">
          <button @click="exportData" class="quick-btn">
            📥 Exporter mes données
          </button>
          <button @click="showDeleteModal = true" class="quick-btn danger">
            🗑️ Supprimer mon compte
          </button>
        </div>
      </div>

      <!-- App Info -->
      <div class="app-info">
        <p>Zekora v1.0.0</p>
        <div class="info-links">
          <NuxtLink to="/support">Aide</NuxtLink>
          <span>•</span>
          <a href="#">Conditions</a>
          <span>•</span>
          <a href="#">Confidentialité</a>
        </div>
      </div>

      <!-- Delete Account Modal -->
      <div v-if="showDeleteModal" class="modal-overlay" @click="showDeleteModal = false">
        <div class="modal-content" @click.stop>
          <h3>⚠️ Supprimer votre compte</h3>
          <p>Cette action est irréversible. Toutes vos données seront supprimées.</p>
          <div class="modal-actions">
            <button @click="showDeleteModal = false" class="btn-cancel">Annuler</button>
            <button @click="deleteAccount" class="btn-delete">Supprimer</button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { userAPI } from '~/composables/useApi'

const showDeleteModal = ref(false)
const securityScore = ref(60)
const kycVerified = ref(false)

const kycStatus = computed(() => {
  if (kycVerified.value) {
    return { label: 'Vérifié', class: 'verified' }
  }
  return { label: 'En attente', class: 'pending' }
})

const exportData = async () => {
  try {
    // Trigger data export
    alert('Un email vous sera envoyé avec vos données.')
  } catch (e) {
    console.error('Export error:', e)
  }
}

const deleteAccount = async () => {
  try {
    // Delete account logic
    alert('Veuillez contacter le support pour supprimer votre compte.')
    showDeleteModal.value = false
  } catch (e) {
    console.error('Delete error:', e)
  }
}

onMounted(async () => {
  try {
    const res = await userAPI.getProfile()
    if (res.data) {
      kycVerified.value = res.data.kyc_verified || false
      // Calculate security score based on 2FA, phone verified, etc.
      let score = 40
      if (res.data.phone_verified) score += 20
      if (res.data.email_verified) score += 20
      if (res.data.two_factor_enabled) score += 20
      securityScore.value = score
    }
  } catch (e) {
    console.error('Error loading profile:', e)
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
/* ========== Settings Page Styles ========== */
.settings-page {
  @apply w-full max-w-full p-0;
}

.page-header {
  @apply mb-6;
}

.page-header h1 {
  @apply text-2xl font-bold mb-1;
  color: #1e293b; /* Light mode text */
}

.dark .page-header h1 {
  color: #fff;
}

.page-header p {
  @apply text-sm m-0;
  color: #64748b;
}

.dark .page-header p {
  color: #94a3b8;
}

.settings-grid {
  @apply flex flex-col gap-3 mb-8;
}

/* Settings Card */
.settings-card {
  @apply flex items-center gap-4 p-4 rounded-2xl transition-all duration-200 no-underline cursor-pointer;
  /* Light mode */
  background: white;
  border: 1px solid #e2e8f0; /* Slate-200 */
  color: #1e293b;
}

.dark .settings-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  color: white;
}

.settings-card:active {
  @apply scale-[0.98];
}

.settings-card:hover {
  border-color: #cbd5e1; /* Slate-300 */
  background: #f8fafc; /* Slate-50 */
}

.dark .settings-card:hover {
  border-color: rgba(99, 102, 241, 0.3);
  background: rgba(255,255,255,0.05);
}

/* Icons */
.card-icon {
  @apply w-12 h-12 rounded-xl flex items-center justify-center text-2xl flex-shrink-0;
}

/* Icon Colors - Light/Dark compatible opacity */
.card-icon.blue { @apply bg-blue-50 text-blue-600 dark:bg-blue-500/15 dark:text-blue-400; }
.card-icon.green { @apply bg-emerald-50 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400; }
.card-icon.purple { @apply bg-purple-50 text-purple-600 dark:bg-purple-500/15 dark:text-purple-400; }
.card-icon.orange { @apply bg-orange-50 text-orange-600 dark:bg-orange-500/15 dark:text-orange-400; }
.card-icon.pink { @apply bg-pink-50 text-pink-600 dark:bg-pink-500/15 dark:text-pink-400; }
.card-icon.teal { @apply bg-teal-50 text-teal-600 dark:bg-teal-500/15 dark:text-teal-400; }

.card-content {
  @apply flex-1 min-w-0;
}

.card-content h3 {
  @apply text-base font-semibold mb-0.5;
  color: #1e293b;
}

.dark .card-content h3 {
  color: white;
}

.card-content p {
  @apply text-xs m-0 truncate;
  color: #64748b;
}

.dark .card-content p {
  color: #94a3b8;
}

/* Status Badge */
.card-status {
  @apply px-2 py-1 rounded-lg text-xs font-semibold;
  background: #fef2f2; /* Red-50 */
  color: #ef4444;
}

.dark .card-status {
  background: rgba(239, 68, 68, 0.15);
}

.card-status.active {
  background: #f0fdf4; /* Green-50 */
  color: #22c55e;
}

.dark .card-status.active {
  background: rgba(34, 197, 94, 0.15);
}

/* Badge */
.card-badge {
  @apply px-2 py-1 rounded-lg text-[10px] font-bold uppercase;
}

.card-badge.verified {
  background: #f0fdf4;
  color: #22c55e;
}

.dark .card-badge.verified {
  background: rgba(34, 197, 94, 0.15);
}

.card-badge.pending {
  background: #fff7ed;
  color: #f97316;
}

.dark .card-badge.pending {
  background: rgba(249, 115, 22, 0.15);
}

.card-arrow {
  @apply text-xl;
  color: #cbd5e1;
}

.dark .card-arrow {
  color: #666; /* or white/20 */
}

/* Quick App Actions */
.quick-section {
  @apply mb-8;
}

.quick-section h2 {
  @apply text-sm font-semibold mb-3 uppercase tracking-wider;
  color: #64748b;
}

.dark .quick-section h2 {
  color: #94a3b8;
}

.quick-actions {
  @apply flex flex-wrap gap-3;
}

.quick-btn {
  @apply px-4 py-3 rounded-xl text-sm font-medium cursor-pointer transition-all border;
  background: white;
  border-color: #e2e8f0;
  color: #1e293b;
}

.dark .quick-btn {
  background: transparent;
  border-color: rgba(255,255,255,0.1);
  color: #fff;
}

.quick-btn:hover {
  background: #f8fafc;
}

.dark .quick-btn:hover {
  background: rgba(255,255,255,0.05);
}

.quick-btn.danger {
  color: #ef4444;
  border-color: #fecaca;
}

.dark .quick-btn.danger {
  border-color: rgba(239, 68, 68, 0.3);
}

.quick-btn.danger:hover {
  background: #fef2f2;
}

.dark .quick-btn.danger:hover {
  background: rgba(239, 68, 68, 0.1);
}

/* App Info */
.app-info {
  @apply text-center py-8 text-xs;
  color: #94a3b8;
}

.info-links {
  @apply flex justify-center gap-2 mt-2;
}

.info-links a {
  color: #64748b;
  @apply no-underline hover:text-indigo-500 transition-colors;
}

.dark .info-links a {
  color: #94a3b8;
  @apply hover:text-white;
}

/* Modal */
.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center z-50 p-4;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(4px);
}

.modal-content {
  @apply rounded-2xl p-6 w-full max-w-sm shadow-2xl;
  background: white;
}

.dark .modal-content {
  background: #1a1a2e;
}

.modal-content h3 {
  @apply text-xl font-bold mb-3;
  color: #1e293b;
}

.dark .modal-content h3 {
  color: white;
}

.modal-content p {
  @apply text-sm mb-6;
  color: #64748b;
}

.dark .modal-content p {
  color: #94a3b8;
}

.modal-actions {
  @apply flex gap-3;
}

.btn-cancel {
  @apply block flex-1 py-3 rounded-xl border-none text-sm font-semibold cursor-pointer;
  background: #f1f5f9;
  color: #1e293b;
}

.dark .btn-cancel {
  background: rgba(255,255,255,0.1);
  color: white;
}

.btn-delete {
  @apply block flex-1 py-3 rounded-xl border-none text-sm font-semibold cursor-pointer;
  background: #ef4444;
  color: white;
}

@media (min-width: 640px) {
  .settings-grid {
    @apply grid grid-cols-2;
  }
}
</style>
