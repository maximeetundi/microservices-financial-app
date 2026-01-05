<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-[9999] bg-black/50 flex items-center justify-center p-4" @click="$emit('close')">
      <div class="bg-white dark:bg-gray-800 rounded-2xl max-w-md w-full p-6 shadow-2xl" @click.stop>
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">Nouvelle conversation</h3>
          <button @click="$emit('close')" class="w-8 h-8 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
          Recherchez un utilisateur par email ou numéro de téléphone
        </p>

        <!-- Search Input -->
        <div class="relative mb-4">
          <input 
            v-model="searchQuery" 
            @input="searchUsers"
            type="text" 
            placeholder="Email ou numéro..." 
            class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          />
          <svg v-if="searching" class="animate-spin w-5 h-5 absolute right-3 top-3.5 text-gray-400" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <!-- Results -->
        <div v-if="searchResults.length > 0" class="space-y-2 max-h-64 overflow-y-auto">
          <div 
            v-for="user in searchResults" 
            :key="user.id"
            @click="selectUser(user)"
            class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer transition-colors"
          >
            <div class="w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center font-bold">
              {{ user.name?.[0]?.toUpperCase() || 'U' }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-medium text-gray-900 dark:text-white">{{ user.name || 'Utilisateur' }}</div>
              <div class="text-sm text-gray-500 truncate">{{ user.email || user.phone }}</div>
            </div>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </div>

        <div v-else-if="searchQuery && !searching" class="text-center py-8 text-gray-500">
          Aucun utilisateur trouvé
        </div>

        <div v-else-if="!searchQuery" class="text-center py-8 text-gray-400 text-sm">
          Commencez à taper pour rechercher...
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import api from '~/composables/useApi'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits(['close', 'userSelected'])

const searchQuery = ref('')
const searchResults = ref<any[]>([])
const searching = ref(false)
let searchTimeout: NodeJS.Timeout | null = null

const searchUsers = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  if (!searchQuery.value || searchQuery.value.length < 3) {
    searchResults.value = []
    return
  }

  searching.value = true
  
  searchTimeout = setTimeout(async () => {
    try {
      const res = await api.get('/auth-service/api/v1/users/search', {
        params: { q: searchQuery.value }
      })
      searchResults.value = res.data?.users || []
    } catch (err) {
      console.error('Search error:', err)
      searchResults.value = []
    } finally {
      searching.value = false
    }
  }, 300)
}

const selectUser = async (user: any) => {
  try {
    // Create conversation
    const res = await api.post('/messaging-service/api/v1/conversations', {
      participant_id: user.id
    })
    
    emit('userSelected', {
      ...res.data,
      name: user.name,
      email: user.email
    })
    emit('close')
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur lors de la création de la conversation')
  }
}

watch(() => props.show, (newVal) => {
  if (!newVal) {
    searchQuery.value = ''
    searchResults.value = []
  }
})
</script>
