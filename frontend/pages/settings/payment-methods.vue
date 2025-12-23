<template>
  <NuxtLayout name="dashboard">
    <div class="payment-methods-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>💳 Moyens de paiement</h1>
        <p>Gérez vos cartes et comptes bancaires</p>
      </div>

      <!-- Cards Section -->
      <div class="section">
        <div class="section-header">
          <h2>Mes cartes</h2>
          <NuxtLink to="/cards/new" class="add-btn">+ Ajouter</NuxtLink>
        </div>
        
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
        </div>

        <div v-else-if="cards.length === 0" class="empty-state">
          <span class="empty-icon">💳</span>
          <p>Aucune carte enregistrée</p>
          <NuxtLink to="/cards/new" class="link-btn">Créer une carte</NuxtLink>
        </div>

        <div v-else class="cards-list">
          <div v-for="card in cards" :key="card.id" class="card-item">
            <div class="card-preview" :class="card.is_virtual ? 'virtual' : 'physical'">
              <span class="card-type">{{ card.is_virtual ? 'Virtuelle' : 'Physique' }}</span>
              <span class="card-number">•••• {{ card.card_number?.slice(-4) || '0000' }}</span>
            </div>
            <div class="card-info">
              <h4>{{ card.is_virtual ? 'Carte Virtuelle' : 'Carte Physique' }}</h4>
              <p>Expire {{ card.expiry_month }}/{{ card.expiry_year }}</p>
            </div>
            <div class="card-status" :class="card.status">
              {{ card.status === 'active' ? 'Active' : 'Inactive' }}
            </div>
          </div>
        </div>
      </div>

      <!-- Bank Accounts Section -->
      <div class="section">
        <div class="section-header">
          <h2>Comptes bancaires</h2>
          <button @click="showAddBank = true" class="add-btn">+ Ajouter</button>
        </div>

        <div v-if="bankAccounts.length === 0" class="empty-state">
          <span class="empty-icon">🏦</span>
          <p>Aucun compte bancaire lié</p>
          <button @click="showAddBank = true" class="link-btn">Ajouter un compte</button>
        </div>

        <div v-else class="bank-list">
          <div v-for="bank in bankAccounts" :key="bank.id" class="bank-item">
            <div class="bank-icon">🏦</div>
            <div class="bank-info">
              <h4>{{ bank.bank_name }}</h4>
              <p>{{ bank.account_number }}</p>
            </div>
            <button @click="removeBank(bank.id)" class="remove-btn">✕</button>
          </div>
        </div>
      </div>

      <!-- Mobile Money Section -->
      <div class="section">
        <div class="section-header">
          <h2>Mobile Money</h2>
          <button @click="showAddMobile = true" class="add-btn">+ Ajouter</button>
        </div>

        <div v-if="mobileAccounts.length === 0" class="empty-state">
          <span class="empty-icon">📱</span>
          <p>Aucun compte Mobile Money</p>
          <button @click="showAddMobile = true" class="link-btn">Ajouter</button>
        </div>

        <div v-else class="mobile-list">
          <div v-for="mobile in mobileAccounts" :key="mobile.id" class="mobile-item">
            <div class="mobile-icon">{{ getOperatorIcon(mobile.operator) }}</div>
            <div class="mobile-info">
              <h4>{{ mobile.operator }}</h4>
              <p>{{ mobile.phone_number }}</p>
            </div>
            <button @click="removeMobile(mobile.id)" class="remove-btn">✕</button>
          </div>
        </div>
      </div>

      <!-- Add Bank Modal -->
      <div v-if="showAddBank" class="modal-overlay" @click="showAddBank = false">
        <div class="modal-content" @click.stop>
          <h3>🏦 Ajouter un compte bancaire</h3>
          <div class="form-group">
            <label>Nom de la banque</label>
            <input v-model="newBank.name" type="text" placeholder="Ex: Société Générale">
          </div>
          <div class="form-group">
            <label>IBAN / Numéro de compte</label>
            <input v-model="newBank.iban" type="text" placeholder="Ex: SN01 1234 5678...">
          </div>
          <div class="modal-actions">
            <button @click="showAddBank = false" class="btn-cancel">Annuler</button>
            <button @click="addBank" class="btn-confirm">Ajouter</button>
          </div>
        </div>
      </div>

      <!-- Add Mobile Modal -->
      <div v-if="showAddMobile" class="modal-overlay" @click="showAddMobile = false">
        <div class="modal-content" @click.stop>
          <h3>📱 Ajouter Mobile Money</h3>
          <div class="form-group">
            <label>Opérateur</label>
            <select v-model="newMobile.operator">
              <option value="">Choisir...</option>
              <option value="Orange Money">Orange Money</option>
              <option value="Wave">Wave</option>
              <option value="MTN MoMo">MTN Mobile Money</option>
              <option value="Free Money">Free Money</option>
            </select>
          </div>
          <div class="form-group">
            <label>Numéro de téléphone</label>
            <input v-model="newMobile.phone" type="tel" placeholder="+221 77 123 45 67">
          </div>
          <div class="modal-actions">
            <button @click="showAddMobile = false" class="btn-cancel">Annuler</button>
            <button @click="addMobile" class="btn-confirm">Ajouter</button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { cardAPI } from '~/composables/useApi'

const loading = ref(true)
const cards = ref([])
const bankAccounts = ref([])
const mobileAccounts = ref([])

const showAddBank = ref(false)
const showAddMobile = ref(false)

const newBank = reactive({ name: '', iban: '' })
const newMobile = reactive({ operator: '', phone: '' })

const getOperatorIcon = (operator) => {
  const icons = {
    'Orange Money': '🟠',
    'Wave': '🔵',
    'MTN MoMo': '🟡',
    'Free Money': '🟢'
  }
  return icons[operator] || '📱'
}

const addBank = () => {
  if (!newBank.name || !newBank.iban) return
  
  bankAccounts.value.push({
    id: Date.now(),
    bank_name: newBank.name,
    account_number: '••••' + newBank.iban.slice(-4)
  })
  
  newBank.name = ''
  newBank.iban = ''
  showAddBank.value = false
}

const removeBank = (id) => {
  bankAccounts.value = bankAccounts.value.filter(b => b.id !== id)
}

const addMobile = () => {
  if (!newMobile.operator || !newMobile.phone) return
  
  mobileAccounts.value.push({
    id: Date.now(),
    operator: newMobile.operator,
    phone_number: newMobile.phone
  })
  
  newMobile.operator = ''
  newMobile.phone = ''
  showAddMobile.value = false
}

const removeMobile = (id) => {
  mobileAccounts.value = mobileAccounts.value.filter(m => m.id !== id)
}

onMounted(async () => {
  try {
    const res = await cardAPI.getAll()
    if (res.data?.cards) {
      cards.value = res.data.cards
    }
  } catch (e) {
    console.error('Error loading cards:', e)
  } finally {
    loading.value = false
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.payment-methods-page {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 1.5rem;
}

.back-link {
  display: inline-block;
  color: #888;
  text-decoration: none;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
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

.section {
  margin-bottom: 2rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.section-header h2 {
  font-size: 0.75rem;
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  margin: 0;
}

.add-btn {
  padding: 0.375rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(99, 102, 241, 0.3);
  background: transparent;
  color: #6366f1;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
}

.loading-state, .empty-state {
  text-align: center;
  padding: 2rem;
  background: rgba(255,255,255,0.03);
  border-radius: 1rem;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 2.5rem;
  display: block;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.empty-state p {
  color: #888;
  margin: 0 0 0.75rem 0;
}

.link-btn {
  color: #6366f1;
  background: none;
  border: none;
  font-size: 0.875rem;
  cursor: pointer;
  text-decoration: none;
}

.cards-list, .bank-list, .mobile-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.card-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
}

.card-preview {
  width: 3.5rem;
  height: 2.25rem;
  border-radius: 0.375rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 0.5rem;
  color: #fff;
}

.card-preview.virtual {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

.card-preview.physical {
  background: linear-gradient(135deg, #1e1e2f, #3d3d5c);
}

.card-type {
  font-size: 0.375rem;
  opacity: 0.7;
}

.card-number {
  font-size: 0.5rem;
  font-family: monospace;
}

.card-info, .bank-info, .mobile-info {
  flex: 1;
  min-width: 0;
}

.card-info h4, .bank-info h4, .mobile-info h4 {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  margin: 0;
}

.card-info p, .bank-info p, .mobile-info p {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

.card-status {
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
}

.card-status.active {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.card-status.inactive {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
}

.bank-item, .mobile-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 0.875rem;
}

.bank-icon, .mobile-icon {
  font-size: 1.5rem;
}

.remove-btn {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 0.5rem;
  border: none;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  cursor: pointer;
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
  margin: 0 0 1rem 0;
  color: #fff;
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

.form-group input, .form-group select {
  width: 100%;
  padding: 0.75rem;
  border-radius: 0.625rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
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
</style>
