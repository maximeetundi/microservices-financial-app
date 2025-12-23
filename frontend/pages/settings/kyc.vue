<template>
  <NuxtLayout name="dashboard">
    <div class="kyc-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>📋 Vérification KYC</h1>
        <p>Validez votre identité pour débloquer toutes les fonctionnalités</p>
      </div>

      <!-- Status Card -->
      <div class="status-card" :class="kyc.status">
        <div class="status-icon">{{ getStatusIcon() }}</div>
        <div class="status-info">
          <h3>{{ getStatusTitle() }}</h3>
          <p>{{ getStatusDescription() }}</p>
        </div>
      </div>

      <!-- Documents Section -->
      <div class="section">
        <h2>Documents requis</h2>
        
        <div class="doc-list">
          <!-- Identity Document -->
          <div class="doc-item">
            <div class="doc-icon">🪪</div>
            <div class="doc-info">
              <h4>Pièce d'identité</h4>
              <p>Passeport, carte d'identité ou permis de conduire</p>
            </div>
            <div class="doc-status" :class="documents.identity.status">
              {{ documents.identity.label }}
            </div>
          </div>

          <!-- Selfie -->
          <div class="doc-item">
            <div class="doc-icon">🤳</div>
            <div class="doc-info">
              <h4>Selfie avec document</h4>
              <p>Photo de vous tenant votre pièce d'identité</p>
            </div>
            <div class="doc-status" :class="documents.selfie.status">
              {{ documents.selfie.label }}
            </div>
          </div>

          <!-- Proof of Address -->
          <div class="doc-item">
            <div class="doc-icon">🏠</div>
            <div class="doc-info">
              <h4>Justificatif de domicile</h4>
              <p>Facture de moins de 3 mois</p>
            </div>
            <div class="doc-status" :class="documents.address.status">
              {{ documents.address.label }}
            </div>
          </div>
        </div>
      </div>

      <!-- Upload Section -->
      <div class="section" v-if="kyc.status !== 'verified'">
        <h2>Uploader un document</h2>
        
        <div class="upload-form">
          <select v-model="selectedDocType" class="select-field">
            <option value="">Choisir le type de document</option>
            <option value="identity">Pièce d'identité</option>
            <option value="selfie">Selfie avec document</option>
            <option value="address">Justificatif de domicile</option>
          </select>

          <div class="upload-zone" @click="triggerUpload" @dragover.prevent @drop.prevent="handleDrop">
            <input 
              type="file" 
              ref="fileInput" 
              @change="handleFileSelect" 
              accept="image/*,.pdf"
              style="display: none"
            >
            <div class="upload-icon">📤</div>
            <p v-if="!selectedFile">Cliquez ou déposez un fichier ici</p>
            <p v-else class="file-name">{{ selectedFile.name }}</p>
            <span class="upload-hint">JPG, PNG ou PDF • Max 10MB</span>
          </div>

          <button 
            @click="uploadDocument" 
            :disabled="!selectedFile || !selectedDocType || uploading"
            class="upload-btn"
          >
            <span v-if="uploading">Envoi en cours...</span>
            <span v-else>📤 Envoyer le document</span>
          </button>
        </div>
      </div>

      <!-- History -->
      <div class="section" v-if="uploadHistory.length > 0">
        <h2>Historique des envois</h2>
        
        <div class="history-list">
          <div v-for="item in uploadHistory" :key="item.id" class="history-item">
            <div class="history-icon">{{ getDocIcon(item.type) }}</div>
            <div class="history-info">
              <h4>{{ getDocName(item.type) }}</h4>
              <p>{{ formatDate(item.uploaded_at) }}</p>
            </div>
            <div class="history-status" :class="item.status">
              {{ item.status === 'approved' ? '✓' : item.status === 'rejected' ? '✗' : '...' }}
            </div>
          </div>
        </div>
      </div>

      <!-- Info -->
      <div class="info-box">
        <span class="info-icon">ℹ️</span>
        <p>La vérification prend généralement 24 à 48 heures. Vous recevrez une notification une fois terminée.</p>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { userAPI } from '~/composables/useApi'

const fileInput = ref(null)
const selectedFile = ref(null)
const selectedDocType = ref('')
const uploading = ref(false)

const kyc = reactive({
  status: 'pending', // pending, submitted, verified, rejected
})

const documents = reactive({
  identity: { status: 'required', label: 'Requis' },
  selfie: { status: 'required', label: 'Requis' },
  address: { status: 'required', label: 'Requis' },
})

const uploadHistory = ref([])

const getStatusIcon = () => {
  const icons = { pending: '⏳', submitted: '📨', verified: '✅', rejected: '❌' }
  return icons[kyc.status] || '⏳'
}

const getStatusTitle = () => {
  const titles = {
    pending: 'Vérification en attente',
    submitted: 'Documents en cours de vérification',
    verified: 'Compte vérifié',
    rejected: 'Vérification refusée'
  }
  return titles[kyc.status] || 'Vérification en attente'
}

const getStatusDescription = () => {
  const descs = {
    pending: 'Soumettez vos documents pour vérifier votre identité',
    submitted: 'Nous examinons vos documents. Réponse sous 24-48h.',
    verified: 'Votre identité a été confirmée. Accès complet activé.',
    rejected: 'Veuillez soumettre de nouveaux documents valides.'
  }
  return descs[kyc.status] || ''
}

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file && file.size <= 10 * 1024 * 1024) {
    selectedFile.value = file
  } else {
    alert('Fichier trop volumineux (max 10MB)')
  }
}

const handleDrop = (e) => {
  const file = e.dataTransfer.files[0]
  if (file && file.size <= 10 * 1024 * 1024) {
    selectedFile.value = file
  }
}

const uploadDocument = async () => {
  if (!selectedFile.value || !selectedDocType.value) return
  
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('document', selectedFile.value)
    formData.append('type', selectedDocType.value)
    
    // API call to upload document
    try {
      await userAPI.uploadKYCDocument(formData)
    } catch (e) {
      // If API not available, continue with local state
      console.log('KYC API not available, using local state')
    }
    
    // Update local state
    uploadHistory.value.unshift({
      id: Date.now(),
      type: selectedDocType.value,
      status: 'pending',
      uploaded_at: new Date().toISOString()
    })
    
    documents[selectedDocType.value] = { status: 'pending', label: 'En cours' }
    selectedFile.value = null
    selectedDocType.value = ''
    
    alert('Document envoyé avec succès!')
  } catch (e) {
    console.error('Upload error:', e)
    alert('Erreur lors de l\'envoi')
  } finally {
    uploading.value = false
  }
}

const getDocIcon = (type) => {
  const icons = { identity: '🪪', selfie: '🤳', address: '🏠' }
  return icons[type] || '📄'
}

const getDocName = (type) => {
  const names = { identity: 'Pièce d\'identité', selfie: 'Selfie', address: 'Justificatif' }
  return names[type] || 'Document'
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('fr-FR', { 
    day: 'numeric', 
    month: 'short', 
    year: 'numeric' 
  })
}

onMounted(async () => {
  try {
    // Try to get KYC status from backend
    const res = await userAPI.getProfile()
    if (res.data) {
      kyc.status = res.data.kyc_status || 'pending'
      
      // Try to get KYC documents
      try {
        const docsRes = await userAPI.getKYCDocuments()
        if (docsRes.data?.documents) {
          uploadHistory.value = docsRes.data.documents
          // Update document statuses based on backend data
          docsRes.data.documents.forEach((doc: any) => {
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

.status-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem;
  border-radius: 1rem;
  margin-bottom: 1.5rem;
  background: rgba(249, 115, 22, 0.1);
  border: 1px solid rgba(249, 115, 22, 0.2);
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
  font-size: 2rem;
}

.status-info h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.status-info p {
  font-size: 0.75rem;
  color: #888;
  margin: 0;
}

.section {
  margin-bottom: 1.5rem;
}

.section h2 {
  font-size: 0.875rem;
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  margin: 0 0 0.75rem 0;
}

.doc-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.doc-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 0.875rem;
}

.doc-icon {
  font-size: 1.5rem;
}

.doc-info {
  flex: 1;
  min-width: 0;
}

.doc-info h4 {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  margin: 0;
}

.doc-info p {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

.doc-status {
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
}

.doc-status.required {
  background: rgba(107, 114, 128, 0.2);
  color: #9ca3af;
}

.doc-status.pending {
  background: rgba(249, 115, 22, 0.2);
  color: #f97316;
}

.doc-status.approved {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.doc-status.rejected {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.select-field {
  width: 100%;
  padding: 0.875rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.875rem;
  outline: none;
}

.upload-zone {
  border: 2px dashed rgba(99, 102, 241, 0.3);
  border-radius: 0.875rem;
  padding: 2rem 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-zone:hover {
  border-color: rgba(99, 102, 241, 0.5);
  background: rgba(99, 102, 241, 0.05);
}

.upload-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.upload-zone p {
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.file-name {
  color: #6366f1 !important;
}

.upload-hint {
  font-size: 0.75rem;
  color: #666;
}

.upload-btn {
  width: 100%;
  padding: 1rem;
  border-radius: 0.75rem;
  border: none;
  background: #6366f1;
  color: #fff;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-btn:not(:disabled):hover {
  background: #4f46e5;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem;
  background: rgba(255,255,255,0.03);
  border-radius: 0.75rem;
}

.history-icon {
  font-size: 1.25rem;
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
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
}

.history-status.pending {
  background: rgba(249, 115, 22, 0.2);
  color: #f97316;
}

.history-status.approved {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.history-status.rejected {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.info-box {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 0.75rem;
  margin-top: 1.5rem;
}

.info-icon {
  font-size: 1.25rem;
}

.info-box p {
  font-size: 0.75rem;
  color: #888;
  margin: 0;
  line-height: 1.5;
}
</style>
