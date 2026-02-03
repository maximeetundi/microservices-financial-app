<template>
  <div class="managers-page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">👥 Équipe</h1>
          <p class="page-subtitle">Gérez les accès à votre boutique</p>
        </div>
        <button @click="showInviteModal = true" class="btn-create">
          <span class="icon">+</span>
          Inviter un membre
        </button>
      </div>

      <!-- Roles Info -->
      <div class="roles-info">
        <div class="role-card">
          <div class="role-icon admin">👑</div>
          <div class="role-details">
            <h4>Administrateur</h4>
            <p>Accès complet sauf suppression boutique</p>
          </div>
        </div>
        <div class="role-card">
          <div class="role-icon editor">✏️</div>
          <div class="role-details">
            <h4>Éditeur</h4>
            <p>Gère produits et commandes</p>
          </div>
        </div>
        <div class="role-card">
          <div class="role-icon viewer">👁️</div>
          <div class="role-details">
            <h4>Observateur</h4>
            <p>Lecture seule</p>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement de l'équipe...</p>
      </div>

      <div v-else class="members-list">
        <!-- Owner Card -->
        <div class="member-card owner">
          <div class="member-avatar owner-avatar">
            {{ ownerInitials }}
          </div>
          <div class="member-info">
            <h3>Propriétaire</h3>
            <p class="member-role">Accès complet</p>
          </div>
          <span class="role-badge owner">👑 Owner</span>
        </div>

        <!-- Managers -->
        <div 
          v-for="manager in managers" 
          :key="manager.user_id || manager.email" 
          class="member-card"
        >
          <div class="member-avatar">
            {{ (manager.first_name?.[0] || manager.email?.[0] || '?').toUpperCase() }}
          </div>
          <div class="member-info">
            <h3>{{ manager.first_name ? `${manager.first_name} ${manager.last_name}` : manager.email }}</h3>
            <p class="member-email">{{ manager.email }}</p>
          </div>
          <div class="member-meta">
            <span :class="['status-badge', manager.status]">
              {{ manager.status === 'active' ? '✓ Actif' : '⏳ En attente' }}
            </span>
            <span :class="['role-badge', manager.role]">{{ getRoleLabel(manager.role) }}</span>
          </div>
          <button @click="removeManager(manager)" class="btn-remove" title="Révoquer l'accès">
            🗑️
          </button>
        </div>

        <!-- Empty State -->
        <div v-if="managers.length === 0" class="empty-state">
          <div class="empty-icon">👥</div>
          <h3>Aucun membre</h3>
          <p>Invitez des collaborateurs pour gérer votre boutique ensemble</p>
          <button @click="showInviteModal = true" class="btn-primary">Inviter un membre</button>
        </div>
      </div>

      <!-- Invite Modal -->
      <Teleport to="body">
        <div v-if="showInviteModal" class="modal-overlay" @click.self="showInviteModal = false">
          <div class="modal-content">
            <button class="close-btn" @click="showInviteModal = false">✕</button>
            <h2 class="modal-title">Inviter un membre</h2>
            
            <form @submit.prevent="inviteManager" class="modal-form">
              <div class="form-group">
                <label>Email du membre *</label>
                <input 
                  v-model="inviteForm.email" 
                  type="email" 
                  required
                  placeholder="colleague@exemple.com"
                  class="input"
                >
              </div>

              <div class="form-group">
                <label>Rôle *</label>
                <div class="role-selector">
                  <label 
                    v-for="role in availableRoles" 
                    :key="role.value"
                    :class="['role-option', inviteForm.role === role.value && 'active']"
                  >
                    <input 
                      type="radio" 
                      :value="role.value" 
                      v-model="inviteForm.role"
                      hidden
                    >
                    <span class="role-option-icon">{{ role.icon }}</span>
                    <span class="role-option-label">{{ role.label }}</span>
                    <span class="role-option-desc">{{ role.desc }}</span>
                  </label>
                </div>
              </div>

              <div class="modal-actions">
                <button type="button" @click="showInviteModal = false" class="btn-secondary">
                  Annuler
                </button>
                <button type="submit" class="btn-primary" :disabled="submitting">
                  <span v-if="submitting" class="spinner-small"></span>
                  Envoyer l'invitation
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type ShopManager } from '@/composables/useShopApi'
import ShopLayout from '@/components/shops/ShopLayout.vue'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const loading = ref(true)
const showInviteModal = ref(false)
const submitting = ref(false)
const shopId = ref('')
const managers = ref<ShopManager[]>([])

const inviteForm = ref({
  email: '',
  role: 'editor'
})

const availableRoles = [
  { value: 'admin', label: 'Administrateur', icon: '👑', desc: 'Accès complet' },
  { value: 'editor', label: 'Éditeur', icon: '✏️', desc: 'Produits & commandes' },
  { value: 'viewer', label: 'Observateur', icon: '👁️', desc: 'Lecture seule' }
]

const ownerInitials = computed(() => 'OP')

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    admin: '👑 Admin',
    editor: '✏️ Éditeur',
    viewer: '👁️ Observateur'
  }
  return labels[role] || role
}

const fetchShopData = async () => {
  try {
    loading.value = true
    const shop = await shopApi.getShop(slug)
    shopId.value = shop.id
    managers.value = shop.managers || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const inviteManager = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    let perms: string[] = []
    if (inviteForm.value.role === 'admin') perms = ['*']
    else if (inviteForm.value.role === 'editor') perms = ['products.*', 'orders.*']
    else perms = ['read']

    await shopApi.inviteManager(shopId.value, inviteForm.value.email, inviteForm.value.role, perms)
    
    await fetchShopData()
    showInviteModal.value = false
    inviteForm.value.email = ''
    inviteForm.value.role = 'editor'
  } catch (e: any) {
    alert('Erreur: ' + (e.message || 'Impossible d\'inviter'))
  } finally {
    submitting.value = false
  }
}

const removeManager = async (manager: ShopManager) => {
  if (!confirm(`Révoquer l'accès pour ${manager.email} ?`)) return
  try {
    await shopApi.removeManager(shopId.value, manager.user_id || '')
    await fetchShopData()
  } catch (e: any) {
    alert('Erreur: ' + (e.message || 'Impossible de révoquer'))
  }
}

onMounted(fetchShopData)

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>

<style scoped>
.managers-page {
  padding: 24px;
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.page-subtitle {
  color: var(--text-muted, #6b7280);
  margin-top: 4px;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  padding: 12px 24px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  font-weight: 600;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

/* Roles Info */
.roles-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.role-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
}

.role-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.role-icon.admin { background: #fef3c7; }
.role-icon.editor { background: #dbeafe; }
.role-icon.viewer { background: #f3f4f6; }

.role-details h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.role-details p {
  font-size: 12px;
  color: var(--text-muted, #6b7280);
  margin: 2px 0 0;
}

/* Members List */
.members-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  transition: all 0.2s;
}

.member-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.member-card.owner {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border-color: rgba(99, 102, 241, 0.3);
}

.member-avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #e5e7eb 0%, #d1d5db 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 18px;
  color: #6b7280;
}

.owner-avatar {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
}

.member-info {
  flex: 1;
  min-width: 0;
}

.member-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.member-email, .member-role {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
  margin: 2px 0 0;
}

.member-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.active { background: #d1fae5; color: #059669; }
.status-badge.pending { background: #fef3c7; color: #d97706; }

.role-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: #f3f4f6;
  color: #374151;
}

.role-badge.owner { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); color: white; }
.role-badge.admin { background: #fef3c7; color: #92400e; }
.role-badge.editor { background: #dbeafe; color: #1e40af; }

.btn-remove {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  opacity: 0.5;
  transition: all 0.2s;
}

.btn-remove:hover {
  background: #fee2e2;
  opacity: 1;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  color: var(--text-primary, #1f2937);
  margin: 0 0 8px;
}

.empty-state p {
  color: var(--text-muted, #6b7280);
  margin: 0 0 24px;
}

/* Loading */
.loading-state {
  text-align: center;
  padding: 60px;
  color: var(--text-muted, #6b7280);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border, #e5e7eb);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.modal-content {
  background: var(--surface, white);
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  padding: 32px;
  position: relative;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text-muted, #6b7280);
}

.modal-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: var(--text-primary, #1f2937);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary, #1f2937);
}

.input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 15px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
}

.input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.role-selector {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.role-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 2px solid var(--border, #e5e7eb);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.role-option:hover {
  border-color: #6366f1;
}

.role-option.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.role-option-icon {
  font-size: 20px;
}

.role-option-label {
  font-weight: 600;
  color: var(--text-primary, #1f2937);
}

.role-option-desc {
  margin-left: auto;
  font-size: 13px;
  color: var(--text-muted, #6b7280);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn-primary {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-secondary {
  padding: 12px 24px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  color: var(--text-primary, #1f2937);
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

/* Dark Mode */
:global(.dark) .member-card,
:global(.dark) .role-card,
:global(.dark) .empty-state {
  background: #1e293b;
  border-color: #334155;
}

:global(.dark) .member-card.owner {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(139, 92, 246, 0.15) 100%);
}

:global(.dark) .modal-content {
  background: #1e293b;
}

:global(.dark) .role-option {
  border-color: #475569;
}

:global(.dark) .status-badge.active { background: rgba(5, 150, 105, 0.2); color: #34d399; }
:global(.dark) .status-badge.pending { background: rgba(217, 119, 6, 0.2); color: #fbbf24; }

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .member-card {
    flex-wrap: wrap;
  }
  
  .member-meta {
    width: 100%;
    justify-content: flex-start;
    margin-top: 8px;
  }
}
</style>
