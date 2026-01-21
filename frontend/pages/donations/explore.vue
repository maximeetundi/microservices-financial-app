<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto py-8 px-4">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">🌍 Explorer les campagnes</h1>
          <p class="text-gray-500 dark:text-gray-400">Découvrez et soutenez des causes inspirantes.</p>
        </div>
        <button @click="navigateTo('/donations')" class="px-4 py-2 bg-white dark:bg-slate-800 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-700 rounded-xl font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors flex items-center gap-2">
            <span>👤</span> Mes campagnes
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="flex gap-4 mb-8 overflow-x-auto pb-2">
         <div class="relative flex-1 min-w-[200px]">
             <input v-model="searchQuery" type="text" placeholder="Rechercher une cause..." class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 transition-shadow">
             <span class="absolute left-3 top-3 text-gray-400">🔍</span>
         </div>
         <select v-model="filterStatus" class="px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800">
             <option value="active">En cours</option>
             <option value="all">Toutes</option>
             <option value="completed">Terminées</option>
         </select>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 h-80 rounded-2xl"></div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredCampaigns.length === 0" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
          <div class="text-6xl mb-4">🔍</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">Aucune campagne trouvée</h3>
          <p class="text-gray-500 mt-2">D'autres campagnes apparaîtront ici bientôt.</p>
      </div>

      <!-- Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div v-for="campaign in filteredCampaigns" :key="campaign.id" 
             @click="navigateTo(`/donations/${campaign.id}`)"
             class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:border-indigo-500/30 transition-all duration-300 cursor-pointer flex flex-col h-full">
            
            <!-- Image -->
            <div class="h-48 bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
                <img v-if="campaign.image_url" :src="campaign.image_url" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" alt="Campaign">
                <div v-else class="w-full h-full flex items-center justify-center text-4xl bg-gradient-to-br from-indigo-500/10 to-purple-500/10">
                    🌍
                </div>
                <!-- Badge -->
                <div class="absolute top-3 right-3 px-3 py-1 rounded-full text-xs font-bold backdrop-blur-md"
                     :class="campaign.status === 'active' ? 'bg-emerald-500/90 text-white' : 'bg-gray-500/90 text-white'">
                    {{ campaign.status === 'active' ? 'En cours' : 'Terminée' }}
                </div>
            </div>

            <!-- Content -->
            <div class="p-5 flex-1 flex flex-col">
                <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2 line-clamp-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
                    {{ campaign.title }}
                </h3>
                <p class="text-gray-500 dark:text-gray-400 text-sm mb-4 line-clamp-2 flex-1">
                    {{ campaign.description }}
                </p>

                <!-- Progress -->
                <div class="bg-gray-100 dark:bg-slate-800 rounded-full h-2.5 mb-2 overflow-hidden">
                    <div class="bg-gradient-to-r from-indigo-500 to-purple-500 h-full rounded-full transition-all duration-1000"
                         :style="{ width: getProgress(campaign) + '%' }"></div>
                </div>
                <div class="flex justify-between items-end text-sm">
                    <div>
                        <p class="font-bold text-gray-900 dark:text-white">{{ formatAmount(campaign.collected_amount, campaign.currency) }}</p>
                        <p class="text-xs text-gray-500">récoltés</p>
                    </div>
                    <div class="text-right" v-if="campaign.target_amount > 0">
                        <p class="font-medium text-gray-600 dark:text-gray-300">{{ Math.round(getProgress(campaign)) }}%</p>
                        <p class="text-xs text-gray-500">de {{ formatAmount(campaign.target_amount, campaign.currency) }}</p>
                    </div>
                    <div class="text-right" v-else>
                         <p class="font-medium text-emerald-600">Sans limite</p>
                    </div>
                </div>
            </div>
            
            <div class="p-4 border-t border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-slate-800/50 flex justify-between items-center group-hover:bg-indigo-50 dark:group-hover:bg-indigo-900/10 transition-colors">
                <span class="text-xs font-medium text-indigo-600 dark:text-indigo-400">Soutenir le projet →</span>
                <span class="text-xs text-gray-400">{{ formatDate(campaign.created_at) }}</span>
            </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useApi } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const apiContext = useApi()
const { donationApi } = apiContext
const authStore = useAuthStore()
const user = computed(() => authStore.user)

definePageMeta({
  middleware: 'auth'
})

const loading = ref(true)
const campaigns = ref<any[]>([])
const searchQuery = ref('')
const filterStatus = ref('active')

const filteredCampaigns = computed(() => {
    let list = Array.isArray(campaigns.value) ? campaigns.value : []
    
    // Filter out user's own campaigns
    if (user.value?.id) {
        list = list.filter(c => c.creator_id !== user.value.id)
    }

    if (filterStatus.value !== 'all') {
        list = list.filter(c => c.status === filterStatus.value)
    }
    if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase()
        list = list.filter(c => c.title.toLowerCase().includes(q) || c.description.toLowerCase().includes(q))
    }
    return list
})

const loadCampaigns = async () => {
    loading.value = true
    try {
        // Fetch public campaigns (limit higher for exploration)
        const res = await donationApi.getCampaigns(100)
        // Ensure data is array (handle {campaigns: []} or [])
        const data = res.data?.campaigns || res.data || []
        campaigns.value = Array.isArray(data) ? data : []
    } catch (e) {
        console.error("Failed to load campaigns", e)
        campaigns.value = [] // Reset on error
    } finally {
        loading.value = false
    }
}

const getProgress = (campaign: any) => {
    if (!campaign.target_amount || campaign.target_amount <= 0) return 100 
    return Math.min(100, (campaign.collected_amount / campaign.target_amount) * 100)
}

const formatAmount = (amount: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount || 0)
}

const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString()
}

onMounted(() => {
    loadCampaigns()
})
</script>
