<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/80 backdrop-blur-sm" @click="$emit('close')"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-gray-900 rounded-3xl w-full max-w-5xl max-h-[90vh] overflow-hidden shadow-2xl border border-gray-700 flex flex-col lg:flex-row">
          
          <!-- Sidebar -->
          <div class="w-full lg:w-80 bg-gray-800 p-6 flex flex-col gap-6 border-b lg:border-b-0 lg:border-r border-gray-700">
            <!-- Header -->
            <div>
              <h3 class="text-xl font-bold text-white mb-2">Codes QR</h3>
              <p class="text-sm text-gray-400">Partagez ces codes pour permettre à vos clients de s'abonner</p>
            </div>

            <!-- Type Selection -->
            <div class="space-y-2">
              <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Type de code</label>
              
              <button @click="qrType = 'ENTERPRISE'" 
                :class="['w-full text-left px-4 py-3.5 rounded-xl font-medium transition-all flex items-center gap-3', 
                  qrType === 'ENTERPRISE' 
                    ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white shadow-lg shadow-primary-900/50' 
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600']">
                <BuildingOfficeIcon class="w-5 h-5" />
                <div>
                  <div>Entreprise (Global)</div>
                  <div class="text-xs opacity-70">Tous les services publics</div>
                </div>
              </button>
              
              <button @click="qrType = 'GROUP'" 
                :class="['w-full text-left px-4 py-3.5 rounded-xl font-medium transition-all flex items-center gap-3', 
                  qrType === 'GROUP' 
                    ? 'bg-gradient-to-r from-indigo-600 to-indigo-700 text-white shadow-lg shadow-indigo-900/50' 
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600']">
                <FolderIcon class="w-5 h-5" />
                <div>
                  <div>Groupe de Services</div>
                  <div class="text-xs opacity-70">Une catégorie spécifique</div>
                </div>
              </button>
              
              <button @click="qrType = 'SERVICE'" 
                :class="['w-full text-left px-4 py-3.5 rounded-xl font-medium transition-all flex items-center gap-3', 
                  qrType === 'SERVICE' 
                    ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white shadow-lg shadow-purple-900/50' 
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600']">
                <TagIcon class="w-5 h-5" />
                <div>
                  <div>Service Spécifique</div>
                  <div class="text-xs opacity-70">Un service direct</div>
                </div>
              </button>
            </div>

            <!-- Filters -->
            <div v-if="qrType === 'GROUP' || qrType === 'SERVICE'" class="space-y-4">
              <div>
                <label class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 block">Groupe</label>
                <select v-model="selectedGroup" 
                  class="w-full bg-gray-700 border-none rounded-xl text-white px-4 py-2.5 focus:ring-2 focus:ring-primary-500">
                  <option value="">-- Sélectionner --</option>
                  <option v-for="grp in publicGroups" :key="grp.id" :value="grp.id">
                    {{ grp.name }}
                  </option>
                </select>
              </div>

              <div v-if="qrType === 'SERVICE' && selectedGroup">
                <label class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 block">Service</label>
                <select v-model="selectedService" 
                  class="w-full bg-gray-700 border-none rounded-xl text-white px-4 py-2.5 focus:ring-2 focus:ring-primary-500">
                  <option value="">-- Sélectionner --</option>
                  <option v-for="svc in selectedGroupServices" :key="svc.id" :value="svc.id">
                    {{ svc.name }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Close Button -->
            <button @click="$emit('close')" 
              class="mt-auto text-gray-400 hover:text-white text-sm flex items-center gap-2 transition-colors">
              <ArrowLeftIcon class="w-4 h-4" />
              Retour au tableau de bord
            </button>
          </div>

          <!-- Preview Area -->
          <div class="flex-1 p-8 flex flex-col items-center justify-center bg-gradient-to-br from-gray-900 to-gray-800 relative overflow-hidden">
            <!-- Background Pattern -->
            <div class="absolute inset-0 opacity-5 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSIjZmZmIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxjaXJjbGUgY3g9IjIwIiBjeT0iMjAiIHI9IjIiLz48L2c+PC9zdmc+')]"></div>
            
            <!-- QR Code Display -->
            <div class="relative z-10">
              <div v-if="!canShowQR" class="text-center text-gray-500">
                <QrCodeIcon class="w-32 h-32 mx-auto mb-4 opacity-30" />
                <p>Sélectionnez un {{ qrType === 'SERVICE' ? 'service' : 'groupe' }} pour afficher le QR code</p>
              </div>
              
              <div v-else class="text-center">
                <!-- QR Container with glow -->
                <div class="relative inline-block">
                  <div class="absolute -inset-4 bg-gradient-to-r from-primary-500 to-purple-500 rounded-3xl blur-xl opacity-30 animate-pulse"></div>
                  <div class="relative bg-white p-4 rounded-2xl shadow-2xl transform transition-all hover:scale-105 duration-300">
                    <div v-if="isGenerating" class="w-56 h-56 md:w-64 md:h-64 flex items-center justify-center">
                      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                    </div>
                    <img v-else-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="QR Code" class="w-56 h-56 md:w-64 md:h-64 object-contain" />
                  </div>
                </div>

                <!-- Title & Description -->
                <div class="mt-8 space-y-2">
                  <h2 class="text-2xl font-bold text-white">{{ displayTitle }}</h2>
                  <p class="text-gray-400 max-w-md mx-auto text-sm">{{ displayDescription }}</p>
                </div>

                <!-- Code Badge -->
                <div class="mt-6 inline-flex items-center gap-3 bg-gray-800 rounded-xl px-5 py-3 border border-gray-700">
                  <span class="font-mono text-primary-400 font-bold tracking-wider text-lg">{{ displayCode }}</span>
                  <button @click="copyCode" 
                    class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
                    title="Copier">
                    <DocumentDuplicateIcon class="w-5 h-5" />
                  </button>
                </div>

                <!-- Actions -->
                <div class="flex justify-center gap-4 mt-8">
                  <button @click="downloadQR"
                    class="px-6 py-3 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl font-semibold hover:from-primary-500 hover:to-primary-600 transition-all flex items-center gap-2 shadow-lg shadow-primary-900/50">
                    <ArrowDownTrayIcon class="w-5 h-5" />
                    Télécharger PNG
                  </button>
                  <button @click="copyLink" 
                    class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-xl font-semibold transition-all flex items-center gap-2">
                    <LinkIcon class="w-5 h-5" />
                    Copier le lien
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { 
  BuildingOfficeIcon, FolderIcon, TagIcon, ArrowLeftIcon,
  QrCodeIcon, DocumentDuplicateIcon, ArrowDownTrayIcon, LinkIcon
} from '@heroicons/vue/24/outline'
import QRCode from 'qrcode'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  enterprise: {
    type: Object,
    required: false,
    default: null
  },
  enterpriseId: {
    type: String,
    default: ''
  },
  baseUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const qrType = ref('ENTERPRISE')
const selectedGroup = ref('')
const selectedService = ref('')
const qrCodeDataUrl = ref('')
const isGenerating = ref(false)

// Computed
const publicGroups = computed(() => {
  return (props.enterprise?.service_groups || []).filter(g => !g.is_private)
})

const selectedGroupServices = computed(() => {
  if (!selectedGroup.value) return []
  const group = props.enterprise?.service_groups?.find(g => g.id === selectedGroup.value)
  return group?.services || []
})

const canShowQR = computed(() => {
  if (qrType.value === 'ENTERPRISE') return true
  if (qrType.value === 'GROUP') return !!selectedGroup.value
  if (qrType.value === 'SERVICE') return !!selectedService.value
  return false
})

const subscriptionLink = computed(() => {
  const entId = props.enterpriseId || props.enterprise?.id || props.enterprise?._id
  if (!entId) return ''
  const origin = typeof window !== 'undefined' ? window.location.origin : 'https://app.tech-afm.com'
  let url = `${origin}/enterprises/${entId}/subscribe`
  
  if (qrType.value === 'SERVICE' && selectedService.value) {
    url += `?service_id=${selectedService.value}`
  } else if (qrType.value === 'GROUP' && selectedGroup.value) {
    url += `?group_id=${selectedGroup.value}`
  }
  return url
})

const displayTitle = computed(() => {
  if (qrType.value === 'ENTERPRISE') return props.enterprise?.name
  if (qrType.value === 'GROUP') {
    const grp = props.enterprise?.service_groups?.find(g => g.id === selectedGroup.value)
    return grp?.name || 'Groupe'
  }
  if (qrType.value === 'SERVICE') {
    for (const grp of (props.enterprise?.service_groups || [])) {
      const svc = grp.services?.find(s => s.id === selectedService.value)
      if (svc) return svc.name
    }
  }
  return 'Code QR'
})

const displayDescription = computed(() => {
  if (qrType.value === 'ENTERPRISE') {
    return "Scannez ce code pour découvrir tous les services de l'entreprise"
  }
  if (qrType.value === 'GROUP') {
    return "Scannez pour voir les services de cette catégorie"
  }
  if (qrType.value === 'SERVICE') {
    return "Scannez pour souscrire directement à ce service"
  }
  return ''
})

const displayCode = computed(() => {
  const entId = props.enterpriseId || props.enterprise?.id || props.enterprise?._id
  if (qrType.value === 'ENTERPRISE') return `ENT-${entId || ''}`
  if (qrType.value === 'GROUP') return `GRP-${selectedGroup.value}`
  if (qrType.value === 'SERVICE') return `SVC-${selectedService.value}`
  return 'CODE'
})

// Generate QR code client-side
const generateQRCode = async () => {
  if (!canShowQR.value) {
    qrCodeDataUrl.value = ''
    return
  }
  
  isGenerating.value = true
  try {
    const url = subscriptionLink.value
    if (!url || String(url).trim().length === 0) {
      qrCodeDataUrl.value = ''
      return
    }
    qrCodeDataUrl.value = await QRCode.toDataURL(url, {
      width: 300,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    })
  } catch (err) {
    console.error('Failed to generate QR code:', err)
    qrCodeDataUrl.value = ''
  } finally {
    isGenerating.value = false
  }
}

// Methods
const copyCode = () => {
  navigator.clipboard.writeText(displayCode.value)
  alert('Code copié : ' + displayCode.value)
}

const copyLink = () => {
  navigator.clipboard.writeText(subscriptionLink.value)
  alert('Lien copié !')
}

const downloadQR = () => {
  if (!qrCodeDataUrl.value) return
  const link = document.createElement('a')
  link.download = `qr_${qrType.value.toLowerCase()}_${displayCode.value}.png`
  link.href = qrCodeDataUrl.value
  link.click()
}

// Watchers
watch([qrType, selectedGroup, selectedService, () => props.isOpen], () => {
  if (props.isOpen) {
    generateQRCode()
  }
}, { immediate: true })

watch(qrType, () => {
  if (qrType.value === 'ENTERPRISE') {
    selectedGroup.value = ''
    selectedService.value = ''
  } else if (qrType.value === 'GROUP') {
    selectedService.value = ''
    // Auto-select first public group
    if (!selectedGroup.value && publicGroups.value.length > 0) {
      selectedGroup.value = publicGroups.value[0].id
    }
  }
})

watch(selectedGroup, () => {
  selectedService.value = ''
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.95);
}
</style>
