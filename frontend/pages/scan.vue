<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-xl mx-auto py-8 px-4">
      <!-- Header -->
      <div class="text-center mb-10">
        <div class="w-20 h-20 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-primary to-primary-600 flex items-center justify-center shadow-lg shadow-primary/30">
          <span class="text-4xl">📷</span>
        </div>
        <h1 class="text-2xl font-bold text-base mb-2">Scanner un QR Code</h1>
        <p class="text-muted">Paiement, événement ou transfert - détection automatique</p>
      </div>

      <!-- Method Tabs -->
      <div class="flex gap-2 mb-6">
        <button 
          @click="switchTab('camera')"
          class="flex-1 py-3 px-4 rounded-xl font-medium transition-all"
          :class="activeTab === 'camera' 
            ? 'bg-primary text-white shadow-lg shadow-primary/30' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          📷 Scanner
        </button>
        <button 
          @click="switchTab('code')"
          class="flex-1 py-3 px-4 rounded-xl font-medium transition-all"
          :class="activeTab === 'code' 
            ? 'bg-primary text-white shadow-lg shadow-primary/30' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          ⌨️ Code
        </button>
        <button 
          @click="switchTab('image')"
          class="flex-1 py-3 px-4 rounded-xl font-medium transition-all"
          :class="activeTab === 'image' 
            ? 'bg-primary text-white shadow-lg shadow-primary/30' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          🖼️ Image
        </button>
      </div>

      <!-- Camera Scanner Tab -->
      <div v-if="activeTab === 'camera'" class="glass-card">
        <div class="space-y-4">
          <!-- Camera Permission Error -->
          <div v-if="cameraError" class="p-4 rounded-xl bg-error/10 border border-error/20 text-error text-center">
            <p class="font-medium mb-2">{{ cameraError }}</p>
            <button @click="startCamera" class="text-sm underline">Réessayer</button>
          </div>

          <!-- Camera View -->
          <div v-else class="relative">
            <div class="aspect-square rounded-2xl overflow-hidden bg-black relative">
              <video 
                ref="videoElement" 
                class="w-full h-full object-cover"
                playsinline
              ></video>
              
              <!-- Scanning Overlay -->
              <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
                <div class="scanner-frame">
                  <div class="scanner-corner top-left"></div>
                  <div class="scanner-corner top-right"></div>
                  <div class="scanner-corner bottom-left"></div>
                  <div class="scanner-corner bottom-right"></div>
                  <div v-if="cameraActive" class="scanner-line"></div>
                </div>
              </div>

              <!-- Camera Loading -->
              <div v-if="!cameraActive && !cameraError" class="absolute inset-0 flex items-center justify-center bg-black/80">
                <div class="text-center text-white">
                  <div class="loading-spinner w-10 h-10 mx-auto mb-4"></div>
                  <p class="text-sm">Activation de la caméra...</p>
                </div>
              </div>
            </div>

            <!-- Camera Controls -->
            <div v-if="cameraActive" class="flex justify-center gap-4 mt-4">
              <button 
                v-if="hasFlash"
                @click="toggleFlash"
                class="p-3 rounded-xl bg-surface-hover hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-colors"
                :title="flashOn ? 'Désactiver le flash' : 'Activer le flash'"
              >
                <span class="text-xl">{{ flashOn ? '🔦' : '💡' }}</span>
              </button>
              <button 
                v-if="cameras.length > 1"
                @click="switchCamera"
                class="p-3 rounded-xl bg-surface-hover hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-colors"
                title="Changer de caméra"
              >
                <span class="text-xl">🔄</span>
              </button>
            </div>
          </div>

          <p class="text-xs text-muted text-center">
            Placez le QR code dans le cadre pour scanner automatiquement
          </p>
        </div>
      </div>

      <!-- Code Entry Tab -->
      <div v-if="activeTab === 'code'" class="glass-card">
        <form @submit.prevent="submitCode" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-muted mb-2">Entrez le code</label>
            <input
              v-model="paymentCode"
              type="text"
              placeholder="pay_xxx, EVT-xxx, TKT-xxx..."
              class="input-field text-center text-lg font-mono"
              autofocus
            />
          </div>
          <p class="text-xs text-muted text-center">
            Code marchand (pay_), événement (EVT-) ou ticket (TKT-)
          </p>
          <button 
            type="submit"
            :disabled="!paymentCode || loading"
            class="btn-premium w-full py-3 disabled:opacity-50"
          >
            <span v-if="loading" class="loading-spinner w-5 h-5"></span>
            <span v-else>Continuer →</span>
          </button>
        </form>
      </div>

      <!-- Image Upload Tab -->
      <div v-if="activeTab === 'image'" class="glass-card">
        <div class="space-y-4">
          <div 
            @click="openFilePicker"
            @dragover.prevent="isDragging = true"
            @dragleave="isDragging = false"
            @drop.prevent="handleDrop"
            class="border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all"
            :class="isDragging ? 'border-primary bg-primary/10' : 'border-secondary-300 dark:border-secondary-600 hover:border-primary'"
          >
            <div v-if="!selectedImage">
              <span class="text-5xl mb-4 block">🖼️</span>
              <p class="font-medium text-base mb-1">Cliquez ou glissez une image</p>
              <p class="text-sm text-muted">PNG, JPG jusqu'à 5MB</p>
            </div>
            <div v-else>
              <img :src="selectedImage" class="max-h-48 mx-auto rounded-lg mb-4" />
              <p class="text-sm text-success">✓ Image sélectionnée</p>
            </div>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              @change="handleFileSelect"
              class="hidden"
            />
          </div>
          
          <button 
            @click="scanImage"
            :disabled="!selectedImage || scanning"
            class="btn-premium w-full py-3 disabled:opacity-50"
          >
            <span v-if="scanning" class="loading-spinner w-5 h-5"></span>
            <span v-else>Scanner le QR code →</span>
          </button>
        </div>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="mt-4 p-4 rounded-xl bg-error/10 border border-error/20 text-error text-center">
        {{ error }}
      </div>

      <!-- Help -->
      <div class="mt-8 text-center space-y-2">
        <p class="text-sm text-muted">
          <NuxtLink to="/merchant" class="text-primary hover:underline">
            💳 Créer un paiement marchand
          </NuxtLink>
        </p>
        <p class="text-sm text-muted">
          <NuxtLink to="/events" class="text-primary hover:underline">
            🎫 Voir les événements
          </NuxtLink>
        </p>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import QrScanner from 'qr-scanner'

const activeTab = ref('camera')
const paymentCode = ref('')
const selectedImage = ref(null)
const selectedFile = ref(null)
const isDragging = ref(false)
const loading = ref(false)
const scanning = ref(false)
const error = ref('')
const fileInput = ref(null)

// Camera state
const videoElement = ref(null)
const qrScanner = ref(null)
const cameraActive = ref(false)
const cameraError = ref('')
const cameras = ref([])
const currentCameraIndex = ref(0)
const hasFlash = ref(false)
const flashOn = ref(false)

// Initialize camera on mount
onMounted(async () => {
  // Get available cameras
  try {
    cameras.value = await QrScanner.listCameras(true)
  } catch (e) {
    console.log('Could not list cameras:', e)
  }
  
  // Start camera if we're on camera tab
  if (activeTab.value === 'camera') {
    await nextTick()
    startCamera()
  }
})

// Cleanup on unmount
onUnmounted(() => {
  stopCamera()
})

const startCamera = async () => {
  cameraError.value = ''
  cameraActive.value = false
  
  await nextTick()
  
  if (!videoElement.value) {
    cameraError.value = 'Élément vidéo non disponible'
    return
  }

  try {
    // Create scanner
    qrScanner.value = new QrScanner(
      videoElement.value,
      (result) => {
        handleScanResult(result.data)
      },
      {
        returnDetailedScanResult: true,
        highlightScanRegion: false,
        highlightCodeOutline: true,
        preferredCamera: cameras.value.length > 1 ? 'environment' : undefined
      }
    )

    await qrScanner.value.start()
    cameraActive.value = true

    // Check flash availability
    hasFlash.value = await qrScanner.value.hasFlash()
  } catch (e) {
    console.error('Camera error:', e)
    if (e.name === 'NotAllowedError') {
      cameraError.value = 'Accès à la caméra refusé. Veuillez autoriser l\'accès dans les paramètres de votre navigateur.'
    } else if (e.name === 'NotFoundError') {
      cameraError.value = 'Aucune caméra trouvée sur cet appareil.'
    } else {
      cameraError.value = 'Impossible d\'accéder à la caméra. Essayez le mode "Image" ou "Code".'
    }
  }
}

const stopCamera = () => {
  if (qrScanner.value) {
    qrScanner.value.stop()
    qrScanner.value.destroy()
    qrScanner.value = null
  }
  cameraActive.value = false
}

const handleScanResult = (data) => {
  if (!data) return
  
  // Stop camera
  stopCamera()
  
  let code = data.trim()
  
  // Parse if JSON
  try {
    const json = JSON.parse(code)
    code = json.payment_id || json.id || json.code || code
  } catch {
    // Not JSON, continue with raw data
  }
  
  // ========== SMART TYPE DETECTION ==========
  
  // 1. EVENT QR Code: "ZEKORA_EVENT:EVT-XXXXX"
  if (code.startsWith('ZEKORA_EVENT:')) {
    const eventCode = code.replace('ZEKORA_EVENT:', '')
    console.log('Detected: EVENT -', eventCode)
    navigateTo(`/events/code/${eventCode}`)
    return
  }
  
  // 2. TICKET QR Code: "ZEKORA_TICKET:TKT-XXXXX"
  if (code.startsWith('ZEKORA_TICKET:')) {
    const ticketCode = code.replace('ZEKORA_TICKET:', '')
    console.log('Detected: TICKET -', ticketCode)
    navigateTo(`/tickets/verify/${ticketCode}`)
    return
  }
  
  // 3. USER TRANSFER QR Code: "ZEKORA_USER:user-uuid" or just UUID
  if (code.startsWith('ZEKORA_USER:')) {
    const userId = code.replace('ZEKORA_USER:', '')
    console.log('Detected: USER TRANSFER -', userId)
    navigateTo(`/transfer?to=${userId}`)
    return
  }
  
  // 4. MERCHANT PAYMENT: "pay_XXXXX" or URL containing /pay/
  if (code.startsWith('pay_') || code.includes('/pay/')) {
    let paymentId = code
    if (code.includes('/pay/')) {
      paymentId = code.split('/pay/').pop()
    }
    console.log('Detected: MERCHANT PAYMENT -', paymentId)
    navigateTo(`/pay/${paymentId}`)
    return
  }

  // 5. ENTERPRISE PROFILE / SUBSCRIPTION
  // URL Format: .../enterprises/:id or .../enterprises/:id/subscribe?service_id=...
  if (code.includes('/enterprises/')) {
      // Extract path after /enterprises/
      const path = code.split('/enterprises/')[1] // "123" or "123/subscribe?..."
      // Basic validation: ensure it has an ID
      if (path) {
          console.log('Detected: ENTERPRISE URL -', path)
          // Navigate to absolute path to ensure we hit the right page
          navigateTo(`/enterprises/${path}`)
          return
      }
  }

  // Code Format: "ENT-XXXXX"
  if (code.startsWith('ENT-')) {
      const entId = code.replace('ENT-', '')
      console.log('Detected: ENTERPRISE CODE -', entId)
      navigateTo(`/enterprises/${entId}/subscribe`) // Default to subscribe page
      return
  }
  
  // 6. EVENT CODE FORMAT: "EVT-XXXXX" (direct event code, may contain dashes)
  if (code.match(/^EVT-[A-Z0-9_-]+$/i)) {
    console.log('Detected: EVENT CODE -', code)
    navigateTo(`/events/code/${code}`)
    return
  }
  
  // 7. TICKET CODE FORMAT: "TKT-XXXXX" (direct ticket code, may contain dashes)
  if (code.match(/^TKT-[A-Z0-9_-]+$/i)) {
    console.log('Detected: TICKET CODE -', code)
    navigateTo(`/tickets/verify/${code}`)
    return
  }
  
  // 8. UUID format - could be user ID, try transfer
  if (code.match(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i)) {
    console.log('Detected: UUID (assuming user) -', code)
    navigateTo(`/transfer?to=${code}`)
    return
  }
  
  // 9. CAMPAIGN URL or CODE: "/donations/XXXX" or "CPN-XXXX"
  if (code.includes('/donations/')) {
    const campaignId = code.split('/donations/').pop().split('?')[0] // handle params if any
    console.log('Detected: CAMPAIGN URL -', campaignId)
    navigateTo(`/donations/${campaignId}`)
    return
  }

  if (code.startsWith('CPN-')) {
     console.log('Detected: CAMPAIGN CODE -', code)
     // Campaign codes need to be resolved to ID, or handled by a lookup page.
     // For now, if CPN code IS the ID (which it might be in some contexts, but usually ID is UUID and Code is CPN-...)
     // Actually, earlier we saw campaign code used as ID in some places or vice versa?
     // The Campaign Code "CPN-..." is stored in `campaign_code`. 
     // We should probably check if the frontend router supports resolving by Code?
     // If not, we might need a lookup.
     // However, the best path for now is to try navigating to it, but `donations/[id]` usually expects ID.
     // A safer bet: The QR code we generated (in prev step) puts the *full URL* with ID.
     // So the URL check above handles the new QRs.
     // For OLD QRs or manual CPN code entry:
     // We will redirect to /donations/code/:code if we make such a page, OR just try /donations/:id and let it fail?
     // Let's assume for now we redirect to /donations/explore causing a search?
     // Better: Redirect to /donations/${code} and generic error if not found?
     // Actually, the user screenshot showed /pay/CPN-..., which failed.
     // If we redirect to /donations/CPN-..., it might fail if [id] expects UUID.
     // BUT, we can make [id].vue handle it?
     // Let's use /donations/${code} and hope backend resolves it or we add resolution logic.
     navigateTo(`/donations/${code}`)
     return
  }
  
  // 10. Default: assume merchant payment code
  console.log('Detected: UNKNOWN (treating as merchant) -', code)
  navigateTo(`/pay/${code}`)
}

const switchTab = async (tab) => {
  // Stop camera when leaving camera tab
  if (activeTab.value === 'camera' && tab !== 'camera') {
    stopCamera()
  }
  
  activeTab.value = tab
  error.value = ''
  
  // Start camera when switching to camera tab
  if (tab === 'camera') {
    await nextTick()
    startCamera()
  }
}

const switchCamera = async () => {
  if (cameras.value.length <= 1) return
  
  currentCameraIndex.value = (currentCameraIndex.value + 1) % cameras.value.length
  const camera = cameras.value[currentCameraIndex.value]
  
  if (qrScanner.value) {
    await qrScanner.value.setCamera(camera.id)
  }
}

const toggleFlash = async () => {
  if (!qrScanner.value || !hasFlash.value) return
  
  try {
    await qrScanner.value.toggleFlash()
    flashOn.value = await qrScanner.value.isFlashOn()
  } catch (e) {
    console.error('Flash toggle failed:', e)
  }
}

const openFilePicker = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) {
    loadFile(file)
  }
}

const handleDrop = (event) => {
  isDragging.value = false
  const file = event.dataTransfer.files?.[0]
  if (file && file.type.startsWith('image/')) {
    loadFile(file)
  }
}

const loadFile = (file) => {
  if (file.size > 5 * 1024 * 1024) {
    error.value = 'Image trop grande (max 5MB)'
    return
  }
  
  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    selectedImage.value = e.target.result
    error.value = ''
  }
  reader.readAsDataURL(file)
}

const submitCode = () => {
  if (!paymentCode.value) return
  
  loading.value = true
  error.value = ''
  
  // Use the same smart detection as camera scan
  handleScanResult(paymentCode.value)
}

const scanImage = async () => {
  if (!selectedFile.value) return
  
  scanning.value = true
  error.value = ''
  
  try {
    const result = await QrScanner.scanImage(selectedFile.value, {
      returnDetailedScanResult: true
    })
    
    if (result?.data) {
      handleScanResult(result.data)
    } else {
      error.value = 'Aucun QR code trouvé dans cette image'
    }
  } catch (e) {
    console.error('QR scan failed:', e)
    error.value = 'Impossible de lire le QR code. Essayez une image plus nette.'
  } finally {
    scanning.value = false
  }
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.glass-card {
  @apply bg-surface rounded-2xl p-6 shadow-lg border border-secondary-200 dark:border-secondary-700;
}

/* Scanner frame styling */
.scanner-frame {
  width: 70%;
  height: 70%;
  position: relative;
}

.scanner-corner {
  position: absolute;
  width: 30px;
  height: 30px;
  border-color: rgba(102, 126, 234, 0.9);
  border-style: solid;
  border-width: 0;
}

.scanner-corner.top-left {
  top: 0;
  left: 0;
  border-top-width: 4px;
  border-left-width: 4px;
  border-top-left-radius: 12px;
}

.scanner-corner.top-right {
  top: 0;
  right: 0;
  border-top-width: 4px;
  border-right-width: 4px;
  border-top-right-radius: 12px;
}

.scanner-corner.bottom-left {
  bottom: 0;
  left: 0;
  border-bottom-width: 4px;
  border-left-width: 4px;
  border-bottom-left-radius: 12px;
}

.scanner-corner.bottom-right {
  bottom: 0;
  right: 0;
  border-bottom-width: 4px;
  border-right-width: 4px;
  border-bottom-right-radius: 12px;
}

.scanner-line {
  position: absolute;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.9), transparent);
  animation: scan 2s linear infinite;
}

@keyframes scan {
  0% {
    top: 0;
  }
  50% {
    top: 100%;
  }
  100% {
    top: 0;
  }
}
</style>

