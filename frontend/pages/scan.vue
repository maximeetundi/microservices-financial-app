<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-xl mx-auto py-8 px-4">
      <!-- Header -->
      <div class="text-center mb-10">
        <div class="w-20 h-20 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-primary to-primary-600 flex items-center justify-center shadow-lg shadow-primary/30">
          <span class="text-4xl">📱</span>
        </div>
        <h1 class="text-2xl font-bold text-base mb-2">Payer un marchand</h1>
        <p class="text-muted">Scannez ou entrez le code de paiement</p>
      </div>

      <!-- Method Tabs -->
      <div class="flex gap-2 mb-6">
        <button 
          @click="activeTab = 'code'"
          class="flex-1 py-3 px-4 rounded-xl font-medium transition-all"
          :class="activeTab === 'code' 
            ? 'bg-primary text-white shadow-lg shadow-primary/30' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          ⌨️ Saisir le code
        </button>
        <button 
          @click="activeTab = 'image'"
          class="flex-1 py-3 px-4 rounded-xl font-medium transition-all"
          :class="activeTab === 'image' 
            ? 'bg-primary text-white shadow-lg shadow-primary/30' 
            : 'bg-surface-hover text-muted hover:text-base'"
        >
          🖼️ Uploader QR
        </button>
      </div>

      <!-- Code Entry Tab -->
      <div v-if="activeTab === 'code'" class="glass-card">
        <form @submit.prevent="submitCode" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-muted mb-2">Code de paiement</label>
            <input
              v-model="paymentCode"
              type="text"
              placeholder="pay_abc123..."
              class="input-field text-center text-lg font-mono"
              autofocus
            />
          </div>
          <p class="text-xs text-muted text-center">
            Le code se trouve sous le QR code du marchand
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
              <span class="text-5xl mb-4 block">📷</span>
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
      <div class="mt-8 text-center">
        <p class="text-sm text-muted">
          Vous n'avez pas de code ?
          <NuxtLink to="/merchant" class="text-primary hover:underline">
            Créer un paiement marchand
          </NuxtLink>
        </p>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import QrScanner from 'qr-scanner'

const activeTab = ref('code')
const paymentCode = ref('')
const selectedImage = ref(null)
const selectedFile = ref(null)
const isDragging = ref(false)
const loading = ref(false)
const scanning = ref(false)
const error = ref('')
const fileInput = ref(null)

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
  
  // Extract payment ID from code
  let code = paymentCode.value.trim()
  
  // Handle URL format
  if (code.includes('/pay/')) {
    code = code.split('/pay/').pop()
  }
  
  // Navigate to payment page
  navigateTo(`/pay/${code}`)
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
      let code = result.data
      
      // Parse if JSON
      try {
        const json = JSON.parse(code)
        code = json.payment_id || json.id || code
      } catch {
        // Not JSON, try URL parsing
        if (code.includes('/pay/')) {
          code = code.split('/pay/').pop()
        }
      }
      
      navigateTo(`/pay/${code}`)
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
</style>
