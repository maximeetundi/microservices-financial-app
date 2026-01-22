<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <BellIcon class="w-6 h-6 text-primary-500" />
          Notifications
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Suivez les activités importantes de votre entreprise
        </p>
      </div>
      <button @click="refresh" 
        class="text-sm text-primary-600 hover:text-primary-700 font-medium flex items-center gap-1">
        <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
        Actualiser
      </button>
    </div>

    <!-- Notifications List -->
    <div class="space-y-4">
      
      <!-- Loading State -->
      <div v-if="isLoading && notifications.length === 0" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
        <p class="mt-4 text-gray-500">Chargement...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="notifications.length === 0" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700">
        <BellIcon class="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <h4 class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">Aucune notification</h4>
        <p class="text-gray-500 text-sm">Vous êtes à jour !</p>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div v-for="notif in notifications" :key="notif.id" 
          class="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-100 dark:border-gray-700 shadow-sm flex gap-4 transition-all hover:shadow-md">
          
          <div :class="[getIconBgClass(notif.type), 'w-10 h-10 rounded-full flex items-center justify-center shrink-0']">
            <component :is="getIcon(notif.type)" class="w-5 h-5 text-white" />
          </div>
          
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-start gap-2">
              <h4 class="font-medium text-gray-900 dark:text-white truncate pr-2">{{ notif.title }}</h4>
              <span class="text-xs text-gray-400 whitespace-nowrap">{{ formatDate(notif.created_at) }}</span>
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 break-words">{{ notif.message }}</p>
            
            <!-- Optional Action Button (if supported by data) -->
            <div v-if="notif.data?.action_url" class="mt-3">
              <a :href="notif.data.action_url" class="text-xs font-medium text-primary-600 hover:text-primary-700 flex items-center gap-1">
                Voir détails <ArrowRightIcon class="w-3 h-3" />
              </a>
            </div>
          </div>

          <div v-if="!notif.is_read" class="shrink-0 self-center">
            <div class="w-2 h-2 rounded-full bg-primary-600"></div>
          </div>
        </div>

        <!-- Load More -->
        <div v-if="hasMore" class="text-center pt-4">
          <button @click="loadMore" :disabled="isLoading"
            class="px-4 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50">
            {{ isLoading ? 'Chargement...' : 'Charger plus' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { enterpriseAPI } from '@/composables/useApi'
import { 
  BellIcon, 
  ArrowPathIcon, 
  CheckCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon,
  UserGroupIcon,
  BanknotesIcon,
  ArrowRightIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  enterprise: {
    type: Object,
    required: true
  }
})

const notifications = ref([])
const isLoading = ref(false)
const offset = ref(0)
const limit = 20
const hasMore = ref(true)

// Icons mapping based on type
const getIcon = (type) => {
  switch (type) {
    case 'success': return CheckCircleIcon
    case 'warning': return ExclamationTriangleIcon
    case 'employee_invite': return UserGroupIcon
    case 'payroll': return BanknotesIcon
    default: return InformationCircleIcon
  }
}

const getIconBgClass = (type) => {
  switch (type) {
    case 'success': return 'bg-green-500'
    case 'warning': return 'bg-amber-500'
    case 'error': return 'bg-red-500'
    case 'employee_invite': return 'bg-purple-500'
    case 'payroll': return 'bg-blue-500'
    default: return 'bg-primary-500'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('fr-FR', {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit'
  }).format(date)
}

const fetchNotifications = async (reset = false) => {
  if (isLoading.value) return
  isLoading.value = true
  
  if (reset) {
    offset.value = 0
    notifications.value = []
    hasMore.value = true
  }

  try {
    // Assuming backend returns { notifications: [], ... }
    const { data } = await enterpriseAPI.getNotifications(props.enterprise.id, limit, offset.value)
    
    // Adjust logic based on actual response structure. 
    // Usually it returns data directly or data.notifications
    const newNotifs = data.notifications || data || []
    
    if (newNotifs.length < limit) {
      hasMore.value = false
    }
    
    if (reset) {
      notifications.value = newNotifs
    } else {
      notifications.value = [...notifications.value, ...newNotifs]
    }
    
    offset.value += limit

  } catch (error) {
    console.error('Failed to fetch notifications', error)
  } finally {
    isLoading.value = false
  }
}

const loadMore = () => fetchNotifications(false)
const refresh = () => fetchNotifications(true)

onMounted(() => {
  fetchNotifications(true)
})

watch(() => props.enterprise, () => {
  fetchNotifications(true)
}, { deep: true })

</script>
