<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <FolderIcon class="w-6 h-6 text-primary-500" />
            Groupes de Services
          </h2>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            Configurez les services que votre entreprise propose à ses clients
          </p>
        </div>
        
        <button 
          @click="addGroup" 
          class="px-5 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl font-medium transition-all shadow-lg shadow-primary-500/25 flex items-center gap-2">
          <PlusIcon class="w-5 h-5" />
          Nouveau Groupe
        </button>
      </div>
    </div>

    <!-- Service Groups Manager -->  
    <ServiceGroupsManager 
      v-model="serviceGroups" 
      :available-wallets="wallets"
    />
    
    <!-- Save Button -->
    <div class="flex justify-end">
      <button 
        @click="saveServices" 
        :disabled="isSaving"
        class="px-6 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl font-medium hover:from-primary-700 hover:to-primary-800 disabled:opacity-50 flex items-center gap-2 shadow-lg shadow-primary-500/25 transition-all">
        <span v-if="isSaving" class="animate-spin">⟳</span>
        {{ isSaving ? 'Enregistrement...' : 'Enregistrer les services' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { FolderIcon, PlusIcon } from '@heroicons/vue/24/outline'
import { enterpriseAPI, walletAPI } from '@/composables/useApi'
import ServiceGroupsManager from './ServiceGroupsManager.vue'

const props = defineProps({
  enterprise: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

const isSaving = ref(false)
const serviceGroups = ref([])
const wallets = ref([])

// Initialize with enterprise data
watch(() => props.enterprise, (ent) => {
  if (ent?.service_groups) {
    serviceGroups.value = JSON.parse(JSON.stringify(ent.service_groups))
  }
}, { immediate: true })

onMounted(async () => {
  try {
    // Fetch user's wallets to use as payment destinations
    const res = await walletAPI.getMyWallets()
    wallets.value = res.data
  } catch (error) {
    console.error('Failed to fetch wallets:', error)
  }
})

const addGroup = () => {
  serviceGroups.value.push({
    id: `group_${Date.now()}`,
    name: '',
    currency: 'XOF',
    is_private: false,
    services: []
  })
}

const saveServices = async () => {
  isSaving.value = true
  try {
    const updated = {
      ...props.enterprise,
      service_groups: serviceGroups.value
    }
    await enterpriseAPI.update(props.enterprise.id, updated)
    emit('update', updated)
  } catch (error) {
    console.error('Failed to save services:', error)
  } finally {
    isSaving.value = false
  }
}
</script>
