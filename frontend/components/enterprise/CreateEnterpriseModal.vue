<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="$emit('close')"></div>
      
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl">
        <button @click="$emit('close')" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
          <XMarkIcon class="w-5 h-5" />
        </button>

        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <BuildingOffice2Icon class="w-6 h-6 text-primary-500" />
          Créer une Entreprise
        </h2>
        
        <form @submit.prevent="createEnterprise" class="space-y-5">
          <!-- Logo Upload -->
          <div class="flex items-center gap-4">
            <div class="relative group">
              <div v-if="form.logoPreview" class="w-20 h-20 rounded-2xl overflow-hidden border-2 border-gray-200 dark:border-gray-600">
                <img :src="form.logoPreview" alt="Logo" class="w-full h-full object-cover" />
              </div>
              <div v-else class="w-20 h-20 rounded-2xl bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-gray-400 border-2 border-dashed border-gray-300 dark:border-gray-600">
                <CameraIcon class="w-8 h-8" />
              </div>
              <label class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover:opacity-100 rounded-2xl cursor-pointer transition-opacity">
                <span class="text-white text-xs">Changer</span>
                <input type="file" @change="handleLogoUpload" accept="image/*" class="hidden" />
              </label>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-700 dark:text-gray-300">Logo de l'entreprise</p>
              <p class="text-xs text-gray-500">Optionnel - Format: JPG, PNG</p>
            </div>
          </div>

          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nom de l'entreprise *</label>
            <input v-model="form.name" type="text" required placeholder="Ex: Ma Super École"
              class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent" />
          </div>

          <!-- Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Type d'activité</label>
            <select v-model="form.type"
              class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent">
              <option value="SERVICE">Service Général</option>
              <option value="SCHOOL">École / Éducation</option>
              <option value="TRANSPORT">Transport</option>
              <option value="UTILITY">Eau / Électricité / Gaz</option>
            </select>
          </div>

          <!-- Employee Range -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nombre d'employés</label>
            <select v-model="form.employee_count_range"
              class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent">
              <option value="1-10">1 - 10 employés</option>
              <option value="11-50">11 - 50 employés</option>
              <option value="51-200">51 - 200 employés</option>
              <option value="201-500">201 - 500 employés</option>
              <option value="500+">Plus de 500 employés</option>
            </select>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-4">
            <button type="button" @click="$emit('close')"
              class="flex-1 px-4 py-2.5 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors font-medium">
              Annuler
            </button>
            <button type="submit" :disabled="isCreating"
              class="flex-1 px-4 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl font-medium hover:from-primary-700 hover:to-primary-800 disabled:opacity-50 transition-all">
              {{ isCreating ? 'Création...' : 'Créer l\'entreprise' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { BuildingOffice2Icon, XMarkIcon, CameraIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI } from '@/composables/useApi'

const emit = defineEmits(['close', 'created'])

const isCreating = ref(false)
const form = reactive({
  name: '',
  type: 'SERVICE',
  employee_count_range: '1-10',
  logo: null,
  logoPreview: null
})

const handleLogoUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // Compress image before storing
  try {
    const compressedFile = await compressImage(file, 500, 500, 0.8)
    form.logo = compressedFile
    form.logoPreview = URL.createObjectURL(compressedFile)
    console.log(`Logo compressed: ${(file.size / 1024).toFixed(1)}KB → ${(compressedFile.size / 1024).toFixed(1)}KB`)
  } catch (error) {
    console.error('Failed to compress image:', error)
    // Fallback to original file if compression fails
    form.logo = file
    form.logoPreview = URL.createObjectURL(file)
  }
}

// Compress image using canvas
const compressImage = (file, maxWidth, maxHeight, quality) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        // Calculate new dimensions
        let width = img.width
        let height = img.height
        
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = Math.round(width * ratio)
          height = Math.round(height * ratio)
        }
        
        // Create canvas and draw resized image
        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)
        
        // Convert to blob
        canvas.toBlob((blob) => {
          if (blob) {
            const compressedFile = new File([blob], file.name, {
              type: 'image/jpeg',
              lastModified: Date.now()
            })
            resolve(compressedFile)
          } else {
            reject(new Error('Failed to create blob'))
          }
        }, 'image/jpeg', quality)
      }
      img.onerror = reject
      img.src = e.target.result
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

const createEnterprise = async () => {
  if (!form.name) return
  isCreating.value = true
  try {
    let logoUrl = ''
    if (form.logo instanceof File) {
      const formData = new FormData()
      formData.append('file', form.logo)
      const { data } = await enterpriseAPI.uploadLogo(formData)
      logoUrl = data.url
    }

    const payload = {
      name: form.name,
      type: form.type,
      employee_count_range: form.employee_count_range,
      logo: logoUrl
    }

    await enterpriseAPI.create(payload)
    emit('created')
  } catch (e) {
    console.error(e)
    // Handle specific error types
    let errorMsg = 'Erreur inconnue'
    if (e.response?.status === 413) {
      errorMsg = 'Le fichier logo est trop volumineux. Veuillez utiliser une image de moins de 5 MB.'
    } else if (e.message === 'Network Error') {
      errorMsg = 'Erreur réseau. Vérifiez votre connexion ou réduisez la taille du logo.'
    } else {
      errorMsg = e.response?.data?.error || e.message
    }
    alert('Erreur: ' + errorMsg)
  } finally {
    isCreating.value = false
  }
}
</script>
