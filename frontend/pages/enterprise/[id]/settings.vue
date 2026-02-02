<template>
  <NuxtLayout name="enterprise">
    <div class="space-y-6">
      <EnterpriseSettings v-model="enterpriseData" @upload-logo="handleLogoUpload" />
      
      <!-- Save Button -->
      <div class="flex justify-end pt-4 border-t border-gray-200 dark:border-gray-700">
        <button 
          @click="saveSettings" 
          :disabled="isSaving" 
          class="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-medium disabled:opacity-50 flex items-center gap-2 transition-colors"
        >
          <span v-if="isSaving" class="animate-spin">⟳</span>
          {{ isSaving ? 'Enregistrement...' : 'Enregistrer les modifications' }}
        </button>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, inject, watch } from 'vue'
import { useRoute } from 'vue-router'
import { enterpriseAPI } from '@/composables/useApi'
import EnterpriseSettings from '@/components/enterprise/EnterpriseSettings.vue'

const route = useRoute()
const enterprise = inject('enterprise', ref(null))
const enterpriseData = ref(null)
const isSaving = ref(false)

// Clone enterprise data for editing
watch(enterprise, (val) => {
  if (val) {
    enterpriseData.value = JSON.parse(JSON.stringify(val))
  }
}, { immediate: true })

const handleLogoUpload = async (file) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const { data } = await enterpriseAPI.uploadLogo(formData)
    enterpriseData.value.logo = data.url
  } catch (e) {
    console.error('Failed to upload logo', e)
    alert('Erreur lors de l\'upload du logo')
  }
}

const saveSettings = async () => {
  if (!enterpriseData.value) return
  isSaving.value = true
  try {
    await enterpriseAPI.update(route.params.id, enterpriseData.value)
    enterprise.value = enterpriseData.value
    alert('Sauvegardé avec succès')
  } catch (error) {
    console.error('Failed to save settings', error)
    alert('Erreur lors de la sauvegarde')
  } finally {
    isSaving.value = false
  }
}
</script>
