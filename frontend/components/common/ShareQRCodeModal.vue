<template>
  <Teleport to="body">
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm" @click.self="close">
      <div class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-sm overflow-hidden animate-in fade-in zoom-in duration-300 border border-gray-100 dark:border-gray-800 shadow-2xl">
        <!-- Header -->
        <div class="p-4 bg-gray-50 dark:bg-slate-800/50 flex justify-between items-center border-b border-gray-100 dark:border-gray-800">
          <h3 class="font-bold text-gray-900 dark:text-white">{{ title || 'Code QR' }}</h3>
          <button @click="close" class="text-gray-500 hover:text-gray-700 dark:hover:text-white transition-colors w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-200 dark:hover:bg-slate-700">✕</button>
        </div>
        
        <div class="p-8 text-center">
          <!-- QR Code Container -->
          <div class="mb-6 relative inline-block group">
            <div class="bg-white p-4 rounded-2xl shadow-xl border border-gray-100 transition-transform duration-300 group-hover:scale-105">
              <img v-if="qrUrl" :src="qrUrl" alt="QR Code" class="w-48 h-48 object-contain">
              <div v-else class="w-48 h-48 bg-gray-100 flex items-center justify-center text-gray-300 rounded-lg">
                <div class="animate-pulse flex flex-col items-center">
                  <span class="text-4xl mb-2">📱</span>
                  <span class="text-xs">Génération...</span>
                </div>
              </div>
            </div>
            <!-- Logo Overlay (Optional) -->
            <div v-if="qrUrl" class="absolute inset-0 flex items-center justify-center pointer-events-none">
              <div class="w-10 h-10 bg-white rounded-full p-1 shadow-md">
                 <img src="/logo-icon.png" onerror="this.style.display='none'" class="w-full h-full object-contain rounded-full">
              </div>
            </div>
          </div>

          <!-- Title/Subtitle -->
          <h3 v-if="subtitle" class="text-lg font-bold text-gray-900 dark:text-white mb-2">{{ subtitle }}</h3>
          
          <!-- Code Display -->
          <div v-if="displayCode" class="mb-8 font-mono bg-indigo-50 dark:bg-slate-800 text-indigo-700 dark:text-indigo-300 py-3 px-4 rounded-xl border border-indigo-100 dark:border-indigo-900/50 select-all flex items-center justify-center gap-2 group cursor-pointer" @click="copyCode">
             <span class="text-lg font-bold tracking-wider">{{ displayCode }}</span>
             <span class="text-indigo-400 opacity-0 group-hover:opacity-100 transition-opacity">📋</span>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-3">
             <button @click="copyCode" class="py-3 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-xl font-bold hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors flex items-center justify-center gap-2">
                <span>📋</span> Copier
             </button>
             <button @click="downloadQR" class="py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors flex items-center justify-center gap-2 shadow-lg shadow-indigo-600/20">
                <span>⬇️</span> PNG
             </button>
             <button @click="shareLink" class="col-span-2 py-3 bg-white border border-gray-200 dark:bg-slate-800 dark:border-gray-700 text-indigo-600 dark:text-indigo-400 rounded-xl font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors flex items-center justify-center gap-2">
                <span>📤</span> Partager le lien
             </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import QRCode from 'qrcode'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: 'Partager'
  },
  subtitle: {
    type: String,
    default: ''
  },
  qrData: {
    type: String,
    required: true
  },
  displayCode: {
    type: String,
    default: ''
  },
  downloadName: {
    type: String,
    default: 'qrcode'
  },
  shareTitle: {
      type: String,
      default: 'Regardez ça !'
  },
  shareText: {
      type: String,
      default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const qrUrl = ref('')

// Generate QR Code when modal opens or data changes
watchEffect(async () => {
    if (props.modelValue && props.qrData) {
        try {
            qrUrl.value = await QRCode.toDataURL(props.qrData, { 
                width: 400, 
                margin: 2, 
                color: { dark: '#4f46e5', light: '#ffffff' },
                errorCorrectionLevel: 'H'
            })
        } catch (e) {
            console.error('QR Gen Error:', e)
        }
    }
})

const close = () => {
  emit('update:modelValue', false)
}

const copyCode = () => {
  const textToCopy = props.displayCode || props.qrData
  navigator.clipboard.writeText(textToCopy)
  alert('Copié dans le presse-papier !')
}

const downloadQR = () => {
  if (!qrUrl.value) return
  
  const link = document.createElement('a')
  link.href = qrUrl.value
  link.download = `${props.downloadName}.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const shareLink = async () => {
    const url = props.qrData.startsWith('http') ? props.qrData : window.location.href
    
    if (navigator.share) {
        try {
            await navigator.share({
                title: props.shareTitle,
                text: props.shareText,
                url: url
            })
        } catch (e) { console.log('Share cancelled') }
    } else {
        copyCode()
    }
}
</script>
