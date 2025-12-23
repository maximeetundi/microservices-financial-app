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
.settings-page {
  width: 100%;
  max-width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 1.5rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.page-header p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

.settings-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 2rem;
}

.settings-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
}

.settings-card:active {
  background: rgba(255,255,255,0.08);
}

.card-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.875rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  flex-shrink: 0;
}

.card-icon.blue { background: rgba(59, 130, 246, 0.15); }
.card-icon.green { background: rgba(34, 197, 94, 0.15); }
.card-icon.purple { background: rgba(168, 85, 247, 0.15); }
.card-icon.orange { background: rgba(249, 115, 22, 0.15); }
.card-icon.pink { background: rgba(236, 72, 153, 0.15); }
.card-icon.teal { background: rgba(20, 184, 166, 0.15); }

.card-content {
  flex: 1;
  min-width: 0;
}

.card-content h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.125rem 0;
}

.card-content p {
  font-size: 0.75rem;
  color: #888;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-status {
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.card-status.active {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.card-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
}

.card-badge.verified {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.card-badge.pending {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.card-arrow {
  color: #666;
  font-size: 1.25rem;
}

.quick-section {
  margin-bottom: 2rem;
}

.quick-section h2 {
  font-size: 0.875rem;
  font-weight: 600;
  color: #888;
  margin: 0 0 0.75rem 0;
  text-transform: uppercase;
}

.quick-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.quick-btn {
  padding: 0.75rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: transparent;
  color: #fff;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: rgba(255,255,255,0.05);
}

.quick-btn.danger {
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

.quick-btn.danger:hover {
  background: rgba(239, 68, 68, 0.1);
}

.app-info {
  text-align: center;
  padding: 2rem 0;
  color: #666;
  font-size: 0.75rem;
}

.info-links {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.info-links a {
  color: #888;
  text-decoration: none;
}

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
  max-width: 400px;
  width: 100%;
}

.modal-content h3 {
  font-size: 1.25rem;
  margin: 0 0 0.75rem 0;
  color: #fff;
}

.modal-content p {
  color: #888;
  margin: 0 0 1.5rem 0;
  font-size: 0.875rem;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
}

.btn-cancel, .btn-delete {
  flex: 1;
  padding: 0.75rem;
  border-radius: 0.75rem;
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel {
  background: rgba(255,255,255,0.1);
  color: #fff;
}

.btn-delete {
  background: #ef4444;
  color: #fff;
}

@media (min-width: 640px) {
  .settings-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }

  .settings-card:hover {
    transform: translateY(-2px);
    border-color: rgba(99, 102, 241, 0.3);
  }
}
</style>
