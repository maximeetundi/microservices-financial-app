<template>
  <NuxtLayout name="dashboard">
    <div class="payment-methods-page">
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>📱 Numéros de recharge</h1>
        <p>Gérez votre répertoire de numéros pour les dépôts Mobile Money (max 8)</p>
      </div>

      <div class="section">
        <div class="section-header">
          <h2>Mes numéros</h2>
          <button @click="showAdd = true" class="add-btn">+ Ajouter</button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
        </div>

        <div v-else-if="numbers.length === 0" class="empty-state">
          <span class="empty-icon">📱</span>
          <p>Aucun numéro enregistré</p>
          <button @click="showAdd = true" class="link-btn">Ajouter un numéro</button>
        </div>

        <div v-else class="mobile-list">
          <div v-for="n in numbers" :key="n.id" class="mobile-item">
            <div class="mobile-icon">📞</div>
            <div class="mobile-info">
              <h4>
                {{ n.label || 'Numéro' }}
                <span v-if="n.is_default" class="badge">Par défaut</span>
              </h4>
              <p>{{ n.phone }} • {{ n.country }}</p>
            </div>
            <button @click="removeNumber(n.id)" class="remove-btn">✕</button>
          </div>
        </div>

        <div v-if="error" class="error-box">
          {{ error }}
        </div>
      </div>

      <div v-if="showAdd" class="modal-overlay" @click="showAdd = false">
        <div class="modal-content" @click.stop>
          <h3>📱 Ajouter un numéro</h3>
          <div class="form-group">
            <label>Pays</label>
            <input v-model="newNumber.country" type="text" placeholder="CI / CMR / SEN ...">
          </div>
          <div class="form-group">
            <label>Numéro</label>
            <input v-model="newNumber.phone" type="tel" placeholder="+225 07 XX XX XX XX">
          </div>
          <div class="form-group">
            <label>Label (optionnel)</label>
            <input v-model="newNumber.label" type="text" placeholder="Ex: Mon MTN">
          </div>
          <div class="form-group">
            <label class="checkbox">
              <input v-model="newNumber.is_default" type="checkbox">
              Définir comme numéro par défaut
            </label>
          </div>
          <div class="modal-actions">
            <button @click="showAdd = false" class="btn-cancel">Annuler</button>
            <button @click="addNumber" class="btn-confirm">Ajouter</button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { depositNumbersAPI } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()

const loading = ref(true)
const numbers = ref<any[]>([])
const error = ref('')
const showAdd = ref(false)

const newNumber = reactive({
  country: '',
  phone: '',
  label: '',
  is_default: false,
})

const loadNumbers = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await depositNumbersAPI.list()
    numbers.value = res.data?.numbers || []

    if (numbers.value.length === 0 && authStore.user?.phone) {
      await depositNumbersAPI.create({
        country: authStore.user?.country || 'CI',
        phone: authStore.user.phone,
        label: 'Numéro principal',
        is_default: true,
      })
      const seeded = await depositNumbersAPI.list()
      numbers.value = seeded.data?.numbers || []
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Erreur de chargement'
  } finally {
    loading.value = false
  }
}

const addNumber = async () => {
  if (!newNumber.country || !newNumber.phone) return
  error.value = ''
  try {
    await depositNumbersAPI.create({
      country: newNumber.country,
      phone: newNumber.phone,
      label: newNumber.label,
      is_default: newNumber.is_default,
    })
    showAdd.value = false
    newNumber.country = ''
    newNumber.phone = ''
    newNumber.label = ''
    newNumber.is_default = false
    await loadNumbers()
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Erreur lors de la création'
  }
}

const removeNumber = async (id: string) => {
  error.value = ''
  try {
    await depositNumbersAPI.remove(id)
    await loadNumbers()
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Suppression impossible'
  }
}

onMounted(async () => {
  await loadNumbers()
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
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.03);
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
  to {
    transform: rotate(360deg);
  }
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
}

.mobile-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.mobile-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 0.875rem;
}

.mobile-icon {
  font-size: 1.5rem;
}

.mobile-info {
  flex: 1;
}

.mobile-info h4 {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  margin: 0;
}

.mobile-info p {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

.badge {
  margin-left: 0.5rem;
  padding: 0.125rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
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

.error-box {
  margin-top: 1rem;
  padding: 0.75rem;
  border-radius: 0.75rem;
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
  font-size: 0.875rem;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
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

.checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border-radius: 0.625rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

.btn-cancel,
.btn-confirm {
  flex: 1;
  padding: 0.75rem;
  border-radius: 0.625rem;
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-confirm {
  background: #6366f1;
  color: #fff;
}
</style>
