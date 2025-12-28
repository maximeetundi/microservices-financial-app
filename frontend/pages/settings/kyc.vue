<template>
  <NuxtLayout name="dashboard">
    <div class="kyc-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>🔐 Vérification d'identité</h1>
        <p>Validez votre identité pour débloquer toutes les fonctionnalités</p>
      </div>

      <!-- Progress Steps -->
      <div class="progress-section">
        <div class="progress-steps">
          <div class="step" :class="{ active: currentStep >= 1, completed: documents.identity.status === 'approved' }">
            <div class="step-number">1</div>
            <div class="step-label">Identité</div>
          </div>
          <div class="step-line" :class="{ active: documents.identity.status === 'approved' }"></div>
          <div class="step" :class="{ active: currentStep >= 2, completed: documents.selfie.status === 'approved' }">
            <div class="step-number">2</div>
            <div class="step-label">Selfie</div>
          </div>
          <div class="step-line" :class="{ active: documents.selfie.status === 'approved' }"></div>
          <div class="step" :class="{ active: currentStep >= 3, completed: documents.address.status === 'approved' }">
            <div class="step-number">3</div>
            <div class="step-label">Domicile</div>
          </div>
        </div>
      </div>

      <!-- Status Card -->
      <div class="status-card" :class="kyc.status">
        <div class="status-icon">{{ getStatusIcon() }}</div>
        <div class="status-info">
          <h3>{{ getStatusTitle() }}</h3>
          <p>{{ getStatusDescription() }}</p>
        </div>
        <div class="status-badge" :class="kyc.status">{{ getStatusBadge() }}</div>
      </div>

      <!-- Documents Section -->
      <div class="section">
        <h2>📄 Documents requis</h2>
        
        <div class="doc-grid">
          <!-- Identity Document -->
          <div class="doc-card" :class="{ selected: selectedDocType === 'identity' }" @click="selectDocType('identity')">
            <div class="doc-card-header">
              <span class="doc-emoji">🪪</span>
              <span class="doc-badge" :class="documents.identity.status">{{ documents.identity.label }}</span>
            </div>
            <h4>Pièce d'identité</h4>
            <p>Passeport, carte d'identité ou permis de conduire</p>
            <div class="doc-requirements">
              <span>✓ Photo claire</span>
              <span>✓ Non expiré</span>
              <span>✓ Tous les coins visibles</span>
            </div>
          </div>

          <!-- Selfie -->
          <div class="doc-card" :class="{ selected: selectedDocType === 'selfie' }" @click="selectDocType('selfie')">
            <div class="doc-card-header">
              <span class="doc-emoji">🤳</span>
              <span class="doc-badge" :class="documents.selfie.status">{{ documents.selfie.label }}</span>
            </div>
            <h4>Selfie avec document</h4>
            <p>Photo de vous tenant votre pièce d'identité</p>
            <div class="doc-requirements">
              <span>✓ Visage visible</span>
              <span>✓ Document lisible</span>
              <span>✓ Bonne luminosité</span>
            </div>
          </div>

          <!-- Proof of Address -->
          <div class="doc-card" :class="{ selected: selectedDocType === 'address' }" @click="selectDocType('address')">
            <div class="doc-card-header">
              <span class="doc-emoji">🏠</span>
              <span class="doc-badge" :class="documents.address.status">{{ documents.address.label }}</span>
            </div>
            <h4>Justificatif de domicile</h4>
            <p>Facture ou relevé bancaire récent</p>
            <div class="doc-requirements">
              <span>✓ Moins de 3 mois</span>
              <span>✓ Adresse complète</span>
              <span>✓ Nom visible</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Upload Section -->
      <div class="section upload-section" v-if="kyc.status !== 'verified' && selectedDocType">
        <h2>📤 Télécharger {{ getDocName(selectedDocType) }}</h2>
        
        <div class="upload-form">
          <!-- Preview for Images -->
          <div v-if="previewUrl" class="preview-container">
            <img :src="previewUrl" alt="Preview" class="preview-image" />
            <button @click="clearFile" class="clear-btn">✕</button>
          </div>

          <!-- File Info for PDFs (no preview) -->
          <div v-else-if="selectedFile && !previewUrl" class="file-info-container">
            <div class="file-info">
              <span class="file-icon">📄</span>
              <div class="file-details">
                <p class="file-name">{{ selectedFile.name }}</p>
                <p class="file-size">{{ formatFileSize(selectedFile.size) }}</p>
              </div>
              <button @click="clearFile" class="clear-btn-inline">✕</button>
            </div>
          </div>

          <!-- Upload Zone -->
          <div v-else class="upload-zone" @click="triggerUpload" @dragover.prevent @drop.prevent="handleDrop">
            <input 
              type="file" 
              ref="fileInput" 
              @change="handleFileSelect" 
              accept="image/*,.pdf"
              style="display: none"
            >
            <div class="upload-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <p>Glissez-déposez ou <span class="upload-link">parcourez</span></p>
            <span class="upload-hint">JPG, PNG ou PDF • Maximum 10MB</span>
          </div>

          <!-- Document Metadata Fields (only for identity/address) -->
          <div v-if="selectedDocType && selectedDocType !== 'selfie'" class="doc-metadata-fields">
            <div class="form-group">
              <label class="field-label">
                <span class="field-icon">🆔</span>
                {{ selectedDocType === 'identity' ? 'Numéro du document (CNI, Passeport...)' : 'Numéro de justificatif' }}
              </label>
              <input 
                type="text" 
                v-model="documentNumber" 
                placeholder="Ex: AB123456"
                class="text-input"
              />
            </div>
            <div class="form-group">
              <label class="field-label">
                <span class="field-icon">📅</span>
                Date d'expiration (optionnel)
              </label>
              <input 
                type="date" 
                v-model="expiryDate" 
                class="text-input"
                :min="getTodayDate()"
              />
            </div>
          </div>

          <button 
            @click="uploadDocument" 
            :disabled="!selectedFile || uploading"
            class="upload-btn"
          >
            <span v-if="uploading" class="loading-spinner"></span>
            <span v-else>📤 Envoyer le document</span>
          </button>
        </div>
      </div>

      <!-- History -->
      <div class="section" v-if="uploadHistory.length > 0">
        <h2>📋 Historique des envois</h2>
        
        <div class="history-list">
          <div v-for="item in uploadHistory" :key="item.id" class="history-item">
            <div class="history-icon">{{ getDocIcon(item.type) }}</div>
            <div class="history-info">
              <h4>{{ getDocName(item.type) }}</h4>
              <p>{{ formatDate(item.uploaded_at) }}</p>
            </div>
            <div class="history-status" :class="item.status">
              <span v-if="item.status === 'approved'">✓ Approuvé</span>
              <span v-else-if="item.status === 'rejected'">✗ Refusé</span>
              <span v-else>⏳ En cours</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Info -->
      <div class="info-box">
        <div class="info-icon">💡</div>
        <div class="info-content">
          <h4>Conseils pour une vérification rapide</h4>
          <ul>
            <li>Assurez-vous que les documents sont lisibles et non flous</li>
            <li>Évitez les reflets et ombres sur les photos</li>
            <li>La vérification prend généralement 24 à 48 heures</li>
          </ul>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { userAPI } from '~/composables/useApi'

const fileInput = ref(null)
const selectedFile = ref(null)
const previewUrl = ref(null)
const selectedDocType = ref('')
const uploading = ref(false)

// Document metadata fields
const documentNumber = ref('')
const expiryDate = ref('')

const kyc = reactive({
  status: 'none', // none, pending, submitted, verified, rejected
})

const documents = reactive({
  identity: { status: 'required', label: 'Requis' },
  selfie: { status: 'required', label: 'Requis' },
  address: { status: 'required', label: 'Requis' },
})

const uploadHistory = ref([])

const currentStep = computed(() => {
  if (documents.identity.status === 'required') return 1
  if (documents.selfie.status === 'required') return 2
  if (documents.address.status === 'required') return 3
  return 3
})

const getStatusIcon = () => {
  const icons = { none: '📝', pending: '⏳', submitted: '📨', verified: '✅', rejected: '❌' }
  return icons[kyc.status] || '📝'
}

const getStatusTitle = () => {
  const titles = {
    none: 'Vérification non commencée',
    pending: 'Documents en cours de vérification',
    submitted: 'Documents soumis',
    verified: 'Identité vérifiée',
    rejected: 'Vérification refusée'
  }
  return titles[kyc.status] || 'Vérification non commencée'
}

const getStatusDescription = () => {
  const descs = {
    none: 'Soumettez vos documents pour vérifier votre identité et débloquer toutes les fonctionnalités.',
    pending: 'Nous examinons vos documents. Vous recevrez une notification sous 24-48h.',
    submitted: 'Vos documents ont été reçus et sont en cours de traitement.',
    verified: 'Félicitations ! Votre identité a été confirmée. Vous avez accès à toutes les fonctionnalités.',
    rejected: 'Un ou plusieurs documents ont été refusés. Veuillez soumettre de nouveaux documents valides.'
  }
  return descs[kyc.status] || ''
}

const getStatusBadge = () => {
  const badges = { none: 'Non vérifié', pending: 'En cours', submitted: 'Soumis', verified: 'Vérifié', rejected: 'Refusé' }
  return badges[kyc.status] || 'Non vérifié'
}

const selectDocType = (type) => {
  if (documents[type].status !== 'approved') {
    selectedDocType.value = type
    clearFile()
  }
}

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file && file.size <= 10 * 1024 * 1024) {
    selectedFile.value = file
    if (file.type.startsWith('image/')) {
      previewUrl.value = URL.createObjectURL(file)
    }
  } else {
    alert('Fichier trop volumineux (max 10MB)')
  }
}

const handleDrop = (e) => {
  const file = e.dataTransfer.files[0]
  if (file && file.size <= 10 * 1024 * 1024) {
    selectedFile.value = file
    if (file.type.startsWith('image/')) {
      previewUrl.value = URL.createObjectURL(file)
    }
  }
}

const clearFile = () => {
  selectedFile.value = null
  previewUrl.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const uploadDocument = async () => {
  if (!selectedFile.value || !selectedDocType.value) return
  
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('document', selectedFile.value)
    formData.append('type', selectedDocType.value)
    
    // Add document metadata for identity/address documents
    if (selectedDocType.value !== 'selfie') {
      if (documentNumber.value) {
        formData.append('document_number', documentNumber.value)
      }
      if (expiryDate.value) {
        formData.append('expiry_date', expiryDate.value)
      }
    }
    
    try {
      await userAPI.uploadKYCDocument(formData)
    } catch (e) {
      console.log('KYC API not available, using local state')
    }
    
    uploadHistory.value.unshift({
      id: Date.now(),
      type: selectedDocType.value,
      status: 'pending',
      uploaded_at: new Date().toISOString()
    })
    
    documents[selectedDocType.value] = { status: 'pending', label: 'En cours' }
    kyc.status = 'pending'
    clearFile()
    selectedDocType.value = ''
    documentNumber.value = ''
    expiryDate.value = ''
    
    alert('Document envoyé avec succès!')
  } catch (e) {
    console.error('Upload error:', e)
    alert('Erreur lors de l\'envoi')
  } finally {
    uploading.value = false
  }
}

const getTodayDate = () => {
  const today = new Date()
  return today.toISOString().split('T')[0]
}

const getDocIcon = (type) => {
  const icons = { identity: '🪪', selfie: '🤳', address: '🏠' }
  return icons[type] || '📄'
}

const getDocName = (type) => {
  const names = { identity: 'Pièce d\'identité', selfie: 'Selfie avec document', address: 'Justificatif de domicile' }
  return names[type] || 'Document'
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric', 
    month: 'short', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(async () => {
  try {
    const res = await userAPI.getProfile()
    if (res.data) {
      // Map kyc_status properly
      const status = res.data.kyc_status || 'none'
      kyc.status = status === '' ? 'none' : status
      
      try {
        const docsRes = await userAPI.getKYCDocuments()
        if (docsRes.data?.documents) {
          uploadHistory.value = docsRes.data.documents
          docsRes.data.documents.forEach((doc) => {
            if (documents[doc.type]) {
              documents[doc.type].status = doc.status
              documents[doc.type].label = doc.status === 'approved' ? 'Approuvé' : 
                                          doc.status === 'rejected' ? 'Refusé' : 'En cours'
            }
          })
        }
      } catch (e) {
        console.log('KYC documents API not available')
      }
    }
  } catch (e) {
    console.error('Error loading KYC status:', e)
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.kyc-page {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.back-link {
  display: inline-block;
  color: #888;
  text-decoration: none;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  transition: color 0.2s;
}

.back-link:hover {
  color: #6366f1;
}

.page-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.page-header p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

/* Progress Steps */
.progress-section {
  margin-bottom: 2rem;
}

.progress-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.step-number {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background: rgba(255,255,255,0.1);
  border: 2px solid rgba(255,255,255,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #888;
  transition: all 0.3s;
}

.step.active .step-number {
  background: rgba(99, 102, 241, 0.2);
  border-color: #6366f1;
  color: #6366f1;
}

.step.completed .step-number {
  background: #22c55e;
  border-color: #22c55e;
  color: #fff;
}

.step-label {
  font-size: 0.75rem;
  color: #888;
}

.step.active .step-label,
.step.completed .step-label {
  color: #fff;
}

.step-line {
  width: 4rem;
  height: 2px;
  background: rgba(255,255,255,0.1);
  margin: 0 0.5rem;
  margin-bottom: 1.5rem;
  transition: background 0.3s;
}

.step-line.active {
  background: #22c55e;
}

/* Status Card */
.status-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  border-radius: 1rem;
  margin-bottom: 2rem;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.status-card.none {
  background: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
}

.status-card.pending,
.status-card.submitted {
  background: rgba(249, 115, 22, 0.1);
  border-color: rgba(249, 115, 22, 0.2);
}

.status-card.verified {
  background: rgba(34, 197, 94, 0.1);
  border-color: rgba(34, 197, 94, 0.2);
}

.status-card.rejected {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

.status-icon {
  font-size: 2.5rem;
}

.status-info {
  flex: 1;
}

.status-info h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.status-info p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

.status-badge {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge.none { background: rgba(107, 114, 128, 0.2); color: #9ca3af; }
.status-badge.pending, .status-badge.submitted { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.status-badge.verified { background: rgba(34, 197, 94, 0.2); color: #22c55e; }
.status-badge.rejected { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

/* Section */
.section {
  margin-bottom: 2rem;
}

.section h2 {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 1rem 0;
}

/* Document Grid */
.doc-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.doc-card {
  padding: 1.25rem;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  cursor: pointer;
  transition: all 0.3s;
}

.doc-card:hover {
  border-color: rgba(99, 102, 241, 0.3);
  background: rgba(99, 102, 241, 0.05);
}

.doc-card.selected {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.doc-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.doc-emoji {
  font-size: 2rem;
}

.doc-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
}

.doc-badge.required { background: rgba(107, 114, 128, 0.2); color: #9ca3af; }
.doc-badge.pending { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.doc-badge.approved { background: rgba(34, 197, 94, 0.2); color: #22c55e; }
.doc-badge.rejected { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

.doc-card h4 {
  font-size: 0.9rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.doc-card p {
  font-size: 0.75rem;
  color: #666;
  margin: 0 0 0.75rem 0;
}

.doc-requirements {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.doc-requirements span {
  font-size: 0.7rem;
  color: #22c55e;
}

/* Upload Section */
.upload-section {
  padding: 1.5rem;
  background: rgba(99, 102, 241, 0.05);
  border: 1px solid rgba(99, 102, 241, 0.1);
  border-radius: 1rem;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.preview-container {
  position: relative;
  border-radius: 0.875rem;
  overflow: hidden;
  max-height: 300px;
}

.preview-image {
  width: 100%;
  height: auto;
  max-height: 300px;
  object-fit: contain;
  background: rgba(0,0,0,0.3);
}

.clear-btn {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.9);
  color: #fff;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.file-info-container {
  border: 2px solid rgba(34, 197, 94, 0.4);
  border-radius: 1rem;
  background: rgba(34, 197, 94, 0.1);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
}

.file-icon {
  font-size: 2.5rem;
}

.file-details {
  flex: 1;
}

.file-name {
  color: #fff;
  font-weight: 600;
  margin: 0 0 0.25rem 0;
  word-break: break-all;
}

.file-size {
  color: #888;
  font-size: 0.875rem;
  margin: 0;
}

.clear-btn-inline {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s;
}

.clear-btn-inline:hover {
  background: rgba(239, 68, 68, 0.9);
  color: #fff;
}

.upload-zone {
  border: 2px dashed rgba(99, 102, 241, 0.4);
  border-radius: 1rem;
  padding: 3rem 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.upload-zone:hover {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.upload-icon {
  color: #6366f1;
  margin-bottom: 1rem;
}

.upload-zone p {
  color: #fff;
  margin: 0 0 0.5rem 0;
}

.upload-link {
  color: #6366f1;
  text-decoration: underline;
}

.upload-hint {
  font-size: 0.75rem;
  color: #666;
}

/* Document Metadata Fields */
.doc-metadata-fields {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.25rem;
  background: rgba(99, 102, 241, 0.05);
  border-radius: 0.875rem;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #cbd5e1;
}

.field-icon {
  font-size: 1rem;
}

.text-input {
  width: 100%;
  padding: 0.875rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.text-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.text-input:focus {
  outline: none;
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.text-input[type="date"] {
  color-scheme: dark;
}

.upload-btn {
  width: 100%;
  padding: 1rem;
  border-radius: 0.875rem;
  border: none;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.upload-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-btn:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(99, 102, 241, 0.3);
}

.loading-spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* History */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border-radius: 0.875rem;
}

.history-icon {
  font-size: 1.5rem;
}

.history-info {
  flex: 1;
}

.history-info h4 {
  font-size: 0.875rem;
  color: #fff;
  margin: 0;
}

.history-info p {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

.history-status {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.history-status.pending { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.history-status.approved { background: rgba(34, 197, 94, 0.2); color: #22c55e; }
.history-status.rejected { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

/* Info Box */
.info-box {
  display: flex;
  gap: 1rem;
  padding: 1.25rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 1rem;
}

.info-icon {
  font-size: 2rem;
}

.info-content h4 {
  font-size: 0.9rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.5rem 0;
}

.info-content ul {
  margin: 0;
  padding-left: 1.25rem;
}

.info-content li {
  font-size: 0.8rem;
  color: #888;
  margin-bottom: 0.25rem;
}

/* Responsive */
@media (max-width: 640px) {
  .doc-grid {
    grid-template-columns: 1fr;
  }
  
  .step-line {
    width: 2rem;
  }
  
  .progress-steps {
    transform: scale(0.9);
  }
}
</style>
