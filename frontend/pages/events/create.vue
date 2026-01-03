<template>
  <NuxtLayout name="dashboard">
    <div class="create-event-page">
      <div class="page-header">
        <NuxtLink to="/events" class="back-link">← Retour</NuxtLink>
        <h1 class="page-title">🎪 Créer un événement</h1>
      </div>

      <form @submit.prevent="createEvent" class="event-form">
        <!-- Basic Info -->
        <section class="form-section">
          <h2 class="section-title">📋 Informations générales</h2>
          
          <div class="form-group">
            <label>Titre de l'événement *</label>
            <input v-model="form.title" type="text" placeholder="Ex: Concert de Jazz" required />
          </div>

          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" rows="4" placeholder="Décrivez votre événement..."></textarea>
          </div>

          <div class="form-group">
            <label>Lieu</label>
            <input v-model="form.location" type="text" placeholder="Ex: Palais des Congrès, Paris" />
          </div>

          <div class="form-group">
            <label>Image de couverture</label>
            <div class="image-upload-wrapper">
              <input 
                type="file" 
                accept="image/*" 
                @change="handleCoverImageUpload" 
                ref="coverImageInput"
                class="file-input"
              />
              <div class="upload-area" @click="$refs.coverImageInput.click()">
                <div v-if="!coverImagePreview" class="upload-placeholder">
                  <span class="upload-icon">📷</span>
                  <span>Cliquez pour choisir une image</span>
                </div>
                <img v-else :src="coverImagePreview" class="cover-preview" />
              </div>
              <button v-if="coverImagePreview" type="button" @click="removeCoverImage" class="remove-image-btn">✕ Supprimer</button>
            </div>
          </div>
        </section>

        <!-- Dates -->
        <section class="form-section">
          <h2 class="section-title">📅 Dates</h2>
          
          <div class="form-row">
            <div class="form-group">
              <label>Début de l'événement *</label>
              <input v-model="form.start_date" type="datetime-local" required />
            </div>
            <div class="form-group">
              <label>Fin de l'événement *</label>
              <input v-model="form.end_date" type="datetime-local" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Début des ventes *</label>
              <input v-model="form.sale_start_date" type="datetime-local" required />
            </div>
            <div class="form-group">
              <label>Fin des ventes *</label>
              <input v-model="form.sale_end_date" type="datetime-local" required />
            </div>
          </div>

          <div class="form-group">
            <label>Devise des prix *</label>
            <select v-model="form.currency" required>
              <option value="XOF">XOF - Franc CFA (BCEAO)</option>
              <option value="XAF">XAF - Franc CFA (BEAC)</option>
              <option value="EUR">EUR - Euro</option>
              <option value="USD">USD - Dollar américain</option>
              <option value="GBP">GBP - Livre sterling</option>
              <option value="MAD">MAD - Dirham marocain</option>
              <option value="TND">TND - Dinar tunisien</option>
              <option value="NGN">NGN - Naira nigérian</option>
              <option value="GHS">GHS - Cedi ghanéen</option>
              <option value="KES">KES - Shilling kenyan</option>
            </select>
          </div>
        </section>

        <!-- Form Fields -->
        <section class="form-section">
          <h2 class="section-title">📝 Formulaire d'inscription</h2>
          <p class="section-desc">Définissez les informations à collecter lors de l'achat</p>
          
          <div class="fields-list">
            <div v-for="(field, index) in form.form_fields" :key="index" class="field-item-wrapper">
              <div class="field-item">
                <input v-model="field.label" type="text" placeholder="Nom du champ" class="field-label" />
                <select v-model="field.type" class="field-type" @change="onFieldTypeChange(field)">
                  <option value="text">Texte</option>
                  <option value="email">Email</option>
                  <option value="phone">Téléphone</option>
                  <option value="number">Nombre</option>
                  <option value="select">Liste déroulante</option>
                </select>
                <label class="field-required">
                  <input type="checkbox" v-model="field.required" /> Obligatoire
                </label>
                <button type="button" @click="removeField(index)" class="remove-btn">✕</button>
              </div>
              <!-- Options for select type -->
              <div v-if="field.type === 'select'" class="field-options">
                <label class="options-label">Options (séparées par des virgules) :</label>
                <input 
                  v-model="field.optionsText" 
                  type="text" 
                  placeholder="Ex: Option 1, Option 2, Option 3"
                  class="options-input"
                  @input="updateFieldOptions(field)"
                />
              </div>
            </div>
          </div>
          <button type="button" @click="addField" class="add-btn">+ Ajouter un champ</button>
        </section>

        <!-- Ticket Tiers -->
        <section class="form-section">
          <h2 class="section-title">🎫 Niveaux de tickets</h2>
          <p class="section-desc">Créez différents niveaux avec des prix et avantages distincts</p>
          
          <div class="tiers-list">
            <div v-for="(tier, index) in form.ticket_tiers" :key="index" class="tier-card">
              <div class="tier-header">
                <button type="button" @click="openIconPicker(index)" class="icon-btn">
                  {{ tier.icon || '🎫' }}
                </button>
                <input v-model="tier.name" type="text" placeholder="Nom du niveau" class="tier-name" required />
                <input v-model="tier.color" type="color" class="tier-color" />
                <button type="button" @click="removeTier(index)" class="remove-btn">✕</button>
              </div>
              <div class="tier-body">
                <div class="tier-row">
                  <div class="form-group">
                    <label>Prix ({{ form.currency }}) *</label>
                    <input v-model.number="tier.price" type="number" min="0" required />
                  </div>
                  <div class="form-group">
                    <label>Quantité (-1 = illimité)</label>
                    <input v-model.number="tier.quantity" type="number" min="-1" />
                  </div>
                </div>
                <div class="form-group">
                  <label>Description</label>
                  <textarea v-model="tier.description" rows="2" placeholder="Avantages inclus..."></textarea>
                </div>
              </div>
            </div>
          </div>
          <button type="button" @click="addTier" class="add-btn">+ Ajouter un niveau</button>
        </section>

        <!-- Submit -->
        <div class="form-actions">
          <button type="submit" :disabled="submitting" class="btn-submit">
            <span v-if="submitting">Création en cours...</span>
            <span v-else>🎉 Créer l'événement</span>
          </button>
        </div>
      </form>

      <!-- Icon Picker Modal -->
      <Teleport to="body">
        <div v-if="showIconPicker" class="modal-overlay" @click.self="showIconPicker = false">
          <div class="icon-picker-modal">
            <h3>Choisir une icône</h3>
            <div class="icons-grid">
              <button 
                v-for="icon in availableIcons" 
                :key="icon" 
                type="button"
                @click="selectIcon(icon)"
                class="icon-option"
              >
                {{ icon }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ticketAPI } from '~/composables/useApi'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const router = useRouter()
const submitting = ref(false)
const showIconPicker = ref(false)
const selectedTierIndex = ref(0)
const coverImagePreview = ref(null)
const coverImageFile = ref(null)
const availableIcons = ref([
  '⭐', '🌟', '✨', '💎', '👑', '🏆', '🎖️', '🥇', '🥈', '🥉',
  '🎫', '🎟️', '🎪', '🎭', '🎬', '🎵', '🎸', '🎤', '🎧', '🎹',
  '🔥', '💫', '⚡', '🌈', '🎯', '🚀', '💥', '🎉', '🎊', '🎁',
])

const form = reactive({
  title: '',
  description: '',
  location: '',
  cover_image: '',
  start_date: '',
  end_date: '',
  sale_start_date: '',
  sale_end_date: '',
  currency: 'XOF',
  form_fields: [
    { name: 'full_name', label: 'Nom complet', type: 'text', required: true },
    { name: 'email', label: 'Email', type: 'email', required: true },
    { name: 'phone', label: 'Téléphone', type: 'phone', required: false },
  ],
  ticket_tiers: [
    { name: 'Standard', icon: '🎫', price: 5000, quantity: -1, description: 'Accès standard', color: '#6366f1' },
  ]
})

const addField = () => {
  form.form_fields.push({
    name: `field_${Date.now()}`,
    label: '',
    type: 'text',
    required: false,
    options: [],
    optionsText: ''
  })
}

const removeField = (index) => {
  form.form_fields.splice(index, 1)
}

const onFieldTypeChange = (field) => {
  if (field.type === 'select') {
    field.options = field.options || []
    field.optionsText = field.optionsText || ''
  }
}

const updateFieldOptions = (field) => {
  if (field.optionsText) {
    field.options = field.optionsText.split(',').map(opt => opt.trim()).filter(opt => opt)
  } else {
    field.options = []
  }
}

// Image upload handlers
const handleCoverImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    coverImageFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      coverImagePreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const removeCoverImage = () => {
  coverImageFile.value = null
  coverImagePreview.value = null
  form.cover_image = ''
}

const addTier = () => {
  form.ticket_tiers.push({
    name: '',
    icon: '🎟️',
    price: 0,
    quantity: -1,
    description: '',
    color: '#8b5cf6'
  })
}

const removeTier = (index) => {
  if (form.ticket_tiers.length > 1) {
    form.ticket_tiers.splice(index, 1)
  }
}

const openIconPicker = (index) => {
  selectedTierIndex.value = index
  showIconPicker.value = true
}

const selectIcon = (icon) => {
  form.ticket_tiers[selectedTierIndex.value].icon = icon
  showIconPicker.value = false
}

const createEvent = async () => {
  submitting.value = true
  try {
    let coverImageUrl = form.cover_image
    
    // Upload cover image if selected
    if (coverImageFile.value) {
      try {
        const uploadRes = await ticketAPI.uploadImage(coverImageFile.value)
        if (uploadRes.data?.url) {
          coverImageUrl = uploadRes.data.url
        }
      } catch (uploadErr) {
        console.warn('Image upload failed, continuing without image:', uploadErr)
      }
    }
    
    // Format dates to ISO
    const payload = {
      ...form,
      cover_image: coverImageUrl,
      start_date: new Date(form.start_date).toISOString(),
      end_date: new Date(form.end_date).toISOString(),
      sale_start_date: new Date(form.sale_start_date).toISOString(),
      sale_end_date: new Date(form.sale_end_date).toISOString(),
      form_fields: form.form_fields.map(f => ({
        ...f,
        name: f.label.toLowerCase().replace(/\s+/g, '_')
      }))
    }
    
    const res = await ticketAPI.createEvent(payload)
    if (res.data?.event?.id) {
      router.push(`/events/${res.data.event.id}`)
    }
  } catch (e) {
    console.error('Failed to create event:', e)
    alert(e.response?.data?.error || 'Erreur lors de la création')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const res = await ticketAPI.getIcons()
    if (res.data?.icons) {
      availableIcons.value = res.data.icons
    }
  } catch (e) {
    // Use default icons
  }
})
</script>

<style scoped>
.create-event-page {
  padding: 24px;
  max-width: 900px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.back-link {
  color: var(--text-muted);
  text-decoration: none;
  font-size: 14px;
  display: inline-block;
  margin-bottom: 8px;
}

.back-link:hover {
  color: var(--text-primary);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  color: var(--text-primary);
}

.form-section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--text-primary);
}

.section-desc {
  color: var(--text-muted);
  font-size: 14px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--bg-tertiary, var(--surface));
  color: var(--text-primary);
  font-size: 14px;
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: var(--text-muted);
}

.form-group select option {
  background: #1a1a2e;
  color: #ffffff;
  padding: 10px;
}

.form-group select {
  color-scheme: dark;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

/* Image Upload Styles */
.image-upload-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.file-input {
  display: none;
}

.upload-area {
  border: 2px dashed var(--border);
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--bg-tertiary, var(--surface));
}

.upload-area:hover {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
}

.upload-icon {
  font-size: 48px;
}

.remove-image-btn {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid #ef4444;
  color: #ef4444;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  align-self: flex-start;
}

.remove-image-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.cover-preview {
  max-width: 100%;
  max-height: 200px;
  border-radius: 12px;
  object-fit: cover;
}

/* Field Items */
.fields-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.field-item-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-options {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
  background: var(--bg-tertiary, var(--surface));
  border-radius: 8px;
  margin-left: 20px;
  border-left: 3px solid #6366f1;
}

.options-label {
  font-size: 12px;
  color: var(--text-muted);
}

.options-input {
  background: var(--bg-tertiary, var(--surface));
  color: var(--text-primary);
  border: 1px solid var(--border);
  padding: 10px 12px;
  border-radius: 8px;
  width: 100%;
}

.options-input::placeholder {
  color: var(--text-muted);
}

.field-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--bg-tertiary, var(--surface));
  padding: 12px;
  border-radius: 10px;
  border: 1px solid var(--border);
}

.field-item input,
.field-item select {
  background: var(--surface);
  color: var(--text-primary);
  border: 1px solid var(--border);
  padding: 10px 12px;
  border-radius: 8px;
}

.field-item input::placeholder {
  color: var(--text-muted);
}

.field-item select option {
  background: #1a1a2e;
  color: #ffffff;
}

.field-item select {
  color-scheme: dark;
}

.field-label {
  flex: 1;
}

.field-type {
  width: 150px;
}

.field-required {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  white-space: nowrap;
  color: var(--text-primary);
}

.remove-btn {
  background: rgba(239, 68, 68, 0.1);
  border: none;
  color: #ef4444;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
}

.remove-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

.add-btn {
  background: rgba(99, 102, 241, 0.1);
  border: 1px dashed #6366f1;
  color: #6366f1;
  padding: 12px 20px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  width: 100%;
}

.add-btn:hover {
  background: rgba(99, 102, 241, 0.2);
}

/* Tier Cards */
.tiers-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 16px;
}

.tier-card {
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  background: var(--bg-tertiary, var(--surface));
}

.tier-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--surface);
}

.icon-btn {
  font-size: 28px;
  background: var(--bg-tertiary, var(--surface));
  border: 2px dashed var(--border);
  border-radius: 10px;
  width: 50px;
  height: 50px;
  cursor: pointer;
}

.icon-btn:hover {
  border-color: #6366f1;
}

.tier-name {
  flex: 1;
  font-size: 16px;
  font-weight: 600;
  background: var(--surface);
  color: var(--text-primary);
  border: 1px solid var(--border);
  padding: 10px 12px;
  border-radius: 8px;
}

.tier-name::placeholder {
  color: var(--text-muted);
}

.tier-color {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.tier-body {
  padding: 16px;
}

.tier-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

/* Submit */
.form-actions {
  text-align: center;
  padding: 24px 0;
}

.btn-submit {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  padding: 16px 48px;
  border-radius: 12px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Icon Picker Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.icon-picker-modal {
  background: var(--surface);
  border-radius: 16px;
  padding: 24px;
  max-width: 400px;
  width: 90%;
}

.icon-picker-modal h3 {
  margin: 0 0 20px;
  text-align: center;
}

.icons-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}

.icon-option {
  font-size: 28px;
  background: var(--surface-hover);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-option:hover {
  background: rgba(99, 102, 241, 0.2);
  border-color: #6366f1;
  transform: scale(1.1);
}

@media (max-width: 768px) {
  .form-row, .tier-row {
    grid-template-columns: 1fr;
  }
  
  .field-item {
    flex-wrap: wrap;
  }
  
  .field-type {
    width: 100%;
  }
}
</style>
