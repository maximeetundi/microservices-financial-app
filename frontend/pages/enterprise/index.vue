<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Portail Entreprise</h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">Pilotez l'ensemble de vos structures professionnelles</p>
        </div>
        <button 
          @click="showCreateModal = true" 
          class="inline-flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white rounded-xl font-medium transition-all shadow-lg shadow-emerald-500/25"
        >
          <PlusIcon class="w-5 h-5" />
          Nouvelle Entreprise
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center py-16">
        <div class="w-12 h-12 border-4 border-emerald-200 border-t-emerald-600 rounded-full animate-spin"></div>
      </div>

      <!-- Enterprise List -->
      <div v-else-if="enterprises.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <NuxtLink 
          v-for="ent in enterprises" 
          :key="ent.id || ent._id" 
          :to="`/enterprise/${ent.id || ent._id}`"
          class="group relative bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 dark:border-gray-700"
        >
          <div class="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
            <ArrowRightCircleIcon class="w-6 h-6 text-emerald-500" />
          </div>

          <div class="flex items-center gap-4 mb-6">
            <div v-if="ent.logo" class="w-14 h-14 rounded-xl border border-gray-200 dark:border-gray-700 p-1 overflow-hidden flex-shrink-0">
              <img :src="ent.logo" class="w-full h-full object-cover rounded-lg" />
            </div>
            <div v-else class="w-14 h-14 rounded-xl bg-gradient-to-br from-emerald-100 to-teal-50 dark:from-emerald-900/30 dark:to-teal-900/10 flex items-center justify-center text-emerald-600 dark:text-emerald-400 font-bold text-xl ring-1 ring-emerald-100 dark:ring-emerald-800 flex-shrink-0">
              {{ ent.name?.charAt(0) || 'E' }}
            </div>
            <div class="min-w-0 flex-1">
              <h3 class="font-bold text-lg text-gray-900 dark:text-white group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors truncate">
                {{ ent.name }}
              </h3>
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 mt-1">
                {{ formatEnterpriseType(ent.type) }}
              </span>
            </div>
          </div>
          
          <div class="flex justify-between items-center text-sm text-gray-500 border-t border-gray-100 dark:border-gray-700 pt-4">
            <span class="flex items-center gap-1">
              <UsersIcon class="w-4 h-4" /> 
              {{ ent.employees_count || 0 }} Membres
            </span>
            <span class="flex items-center gap-1 text-green-600 dark:text-green-400">
              <CheckCircleIcon class="w-4 h-4" /> Actif
            </span>
          </div>
        </NuxtLink>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl border-2 border-dashed border-gray-200 dark:border-gray-700">
        <BuildingOffice2Icon class="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">Aucune entreprise</h3>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Commencez par créer votre première structure.</p>
        <button 
          @click="showCreateModal = true" 
          class="mt-4 px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-medium transition-colors"
        >
          Créer une entreprise
        </button>
      </div>

      <!-- Create Enterprise Modal -->
      <CreateEnterpriseModal 
        v-if="showCreateModal"
        @close="showCreateModal = false"
        @created="handleEnterpriseCreated" 
      />
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { enterpriseAPI } from '@/composables/useApi'
import { 
  PlusIcon, 
  UsersIcon, 
  CheckCircleIcon, 
  ArrowRightCircleIcon,
  BuildingOffice2Icon
} from '@heroicons/vue/24/outline'

const CreateEnterpriseModal = defineAsyncComponent(() => import('@/components/enterprise/CreateEnterpriseModal.vue'))

const enterprises = ref([])
const isLoading = ref(true)
const showCreateModal = ref(false)

const formatEnterpriseType = (type) => {
  const map = {
    'SERVICE': 'Service',
    'GENERIC': 'Général',
    'SCHOOL': 'École',
    'TRANSPORT': 'Transport',
    'UTILITY': 'Utilitaire'
  }
  return map[type] || type || 'Entreprise'
}

const fetchEnterprises = async () => {
  try {
    isLoading.value = true
    const { data } = await enterpriseAPI.list()
    enterprises.value = data || []
  } catch (error) {
    console.error('Failed to fetch enterprises', error)
    enterprises.value = []
  } finally {
    isLoading.value = false
  }
}

const handleEnterpriseCreated = () => {
  showCreateModal.value = false
  fetchEnterprises()
}

onMounted(fetchEnterprises)
</script>
