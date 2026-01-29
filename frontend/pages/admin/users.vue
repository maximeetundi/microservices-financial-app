<template>
  <NuxtLayout name="admin">
    <div class="p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">Utilisateurs</h1>
        <p class="text-slate-400">Gestion des utilisateurs et de leurs permissions.</p>
      </div>

      <!-- Content -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
      </div>

      <div v-else class="bg-slate-800/50 backdrop-blur-xl rounded-2xl border border-slate-700/50 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead class="bg-slate-900/50 text-slate-400 text-xs uppercase font-medium">
              <tr>
                <th class="px-6 py-4">Utilisateur</th>
                <th class="px-6 py-4">Contact</th>
                <th class="px-6 py-4">Rôle</th>
                <th class="px-6 py-4">Statut</th>
                <th class="px-6 py-4">KYC</th>
                <th class="px-6 py-4">Dernière connexion</th>
                <th class="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr v-for="user in users" :key="user.id" class="hover:bg-slate-700/20 transition-colors">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-indigo-500/20 flex items-center justify-center text-indigo-400 font-bold text-sm">
                      {{ user.name.charAt(0) }}
                    </div>
                    <div>
                      <div class="font-medium text-white">{{ user.name }}</div>
                      <div class="text-xs text-slate-500">ID: {{ user.id.substring(0, 8) }}...</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="text-sm text-slate-300">{{ user.email }}</div>
                  <div class="text-xs text-slate-500">{{ user.phone }}</div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="user.role === 'admin' ? 'bg-purple-500/10 text-purple-400' : 'bg-slate-700 text-slate-300'">
                    {{ user.role }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="user.is_active ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400'">
                    {{ user.is_active ? 'Actif' : 'Inactif' }}
                  </span>
                </td>
                <td class="px-6 py-4">
                   <span class="px-2 py-1 rounded-full text-xs font-medium" 
                    :class="{
                      'bg-emerald-500/10 text-emerald-400': user.kyc_level >= 2,
                      'bg-amber-500/10 text-amber-400': user.kyc_level === 1,
                      'bg-slate-700 text-slate-300': user.kyc_level === 0
                    }">
                    Niveau {{ user.kyc_level }}
                  </span>
                </td>
                <td class="px-6 py-4 text-xs text-slate-400">
                  {{ formatDate(user.last_login) }}
                </td>
                <td class="px-6 py-4 text-right">
                  <button class="text-indigo-400 hover:text-indigo-300 font-medium text-sm">Détails</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="px-6 py-4 border-t border-slate-700/50 flex justify-between items-center">
          <button @click="prevPage" :disabled="offset === 0" class="text-slate-400 hover:text-white disabled:opacity-50">
            ← Précédent
          </button>
           <span class="text-slate-500 text-sm">Page {{ currentPage }} / {{ totalPages }}</span>
          <button @click="nextPage" :disabled="offset + limit >= total" class="text-slate-400 hover:text-white disabled:opacity-50">
            Suivant →
          </button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useApi } from '@/composables/useApi'

const { adminUserAPI } = useApi()

const loading = ref(true)
const users = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)

const currentPage = computed(() => Math.floor(offset.value / limit.value) + 1)
const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)

const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await adminUserAPI.getUsers(limit.value, offset.value)
    users.value = response.data?.users || []
    total.value = response.data?.total || 0
  } catch (error) {
    console.error("Failed to fetch users", error)
  } finally {
    loading.value = false
  }
}

const nextPage = () => {
  if (offset.value + limit.value < total.value) {
    offset.value += limit.value
    fetchUsers()
  }
}

const prevPage = () => {
  if (offset.value >= limit.value) {
    offset.value -= limit.value
    fetchUsers()
  }
}

const formatDate = (dateString) => {
  if (!dateString || dateString.startsWith('0001')) return 'Jamais'
  return new Date(dateString).toLocaleString('fr-FR')
}

onMounted(() => {
  fetchUsers()
})
</script>
