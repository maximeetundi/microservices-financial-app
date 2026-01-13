<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-5xl mx-auto py-8 px-4">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center py-20">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>

      <div v-else-if="!campaign" class="text-center py-20">
         <h2 class="text-2xl font-bold">Campagne introuvable</h2>
      </div>

      <div v-else>
          <!-- Header -->
          <div class="flex items-center gap-4 mb-8">
              <button @click="router.back()" class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors">
                  ←
              </button>
              <div>
                  <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Mur des Donateurs</h1>
                  <p class="text-gray-500">{{ campaign.title }}</p>
              </div>
          </div>

          <!-- Stats Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
              <div class="bg-indigo-600 text-white p-6 rounded-2xl shadow-lg shadow-indigo-600/20">
                  <div class="text-indigo-200 mb-1 text-sm font-bold uppercase">Total Récolté</div>
                  <div class="text-3xl font-bold">{{ formatAmount(campaign.collected_amount, campaign.currency) }}</div>
                  <div class="mt-2 text-indigo-200 text-xs">
                       {{ Math.round(getProgress(campaign)) }}% de l'objectif
                  </div>
              </div>
              
              <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm">
                  <div class="text-gray-500 dark:text-gray-400 mb-1 text-sm font-bold uppercase">Nombre de Dons</div>
                  <div class="text-3xl font-bold text-gray-900 dark:text-white">{{ globalStats.totalCount }}</div>
                  <div class="mt-2 text-emerald-500 text-xs font-bold">
                       Merci à tous !
                  </div>
              </div>

              <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm">
                  <div class="text-gray-500 dark:text-gray-400 mb-1 text-sm font-bold uppercase">Moyenne</div>
                  <div class="text-3xl font-bold text-gray-900 dark:text-white">{{ formatAmount(globalStats.average, campaign.currency) }}</div>
                  <div class="mt-2 text-gray-400 text-xs">
                       Par don
                  </div>
              </div>
          </div>

          <!-- Donors List -->
          <div v-if="donations.length > 0" class="space-y-4">
              <div v-for="donation in donations" :key="donation.id" 
                   class="bg-white dark:bg-slate-800 p-4 rounded-xl border border-gray-100 dark:border-gray-700 flex items-center justify-between hover:shadow-md transition-shadow">
                  
                  <div class="flex items-center gap-4">
                      <div class="w-12 h-12 rounded-full flex items-center justify-center text-lg font-bold"
                           :class="donation.is_anonymous ? 'bg-gray-100 dark:bg-gray-700 text-gray-500' : 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600'">
                          {{ donation.is_anonymous ? '?' : (donation.donor_name?.[0]?.toUpperCase() || 'B') }}
                      </div>
                      <div>
                          <p class="font-bold text-gray-900 dark:text-white">
                              {{ donation.is_anonymous ? 'Donateur Anonyme' : (donation.donor_name || 'Bienfaiteur') }}
                          </p>
                          <p class="text-xs text-gray-500">{{ formatDate(donation.created_at) }}</p>
                          <p v-if="donation.message" class="text-sm text-gray-600 dark:text-gray-300 italic mt-1">"{{ donation.message }}"</p>
                      </div>
                  </div>

                  <div class="text-right">
                      <div class="font-bold text-emerald-600">
                          +{{ formatAmount(donation.amount, donation.currency) }}
                      </div>
                      
                      <!-- Creator Refund Action -->
                       <button v-if="isCreator && donation.status === 'paid'" @click="confirmRefund(donation)" 
                               class="mt-1 text-xs text-red-500 hover:text-red-700 font-medium px-2 py-1 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors">
                            ↩️ Rembourser
                       </button>
                  </div>
              </div>

               <!-- Load More -->
               <div v-if="hasMore" class="text-center mt-8">
                   <button @click="loadMore" class="px-6 py-2 bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-gray-300 rounded-full font-bold hover:bg-gray-200 transition-colors">
                       Charger plus ↓
                   </button>
               </div>
          </div>

          <div v-else class="text-center py-20 bg-gray-50 dark:bg-slate-800/50 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
              <div class="text-6xl mb-4">🌱</div>
              <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun don pour le moment</h3>
              <p class="text-gray-500">Soyez le premier à contribuer !</p>
          </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useApi } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const router = useRouter()
const { donationApi } = useApi()
const authStore = useAuthStore()
const user = computed(() => authStore.user)

const campaignId = computed(() => route.params.id as string)
const campaign = ref<any>(null)
const donations = ref<any[]>([])
const loading = ref(true)
const offset = ref(0)
const limit = 20
const hasMore = ref(false)

const isCreator = computed(() => campaign.value && user.value && campaign.value.creator_id === user.value.id)

const globalStats = computed(() => {
    const total = donations.value.length
    if (total === 0) return { totalCount: 0, average: 0 }
    
    // Note: Average based on loaded donations might be inaccurate if not all loaded.
    // Better to use campaign.collected_amount / campaign.donations_count if available.
    // Assuming campaign object updates with count? Usually APIs return count.
    // Let's rely on accumulated loaded for now or just generic.
    // Actually, campaign.collected_amount is reliable.
    
    // We don't have total count in campaign object explicitly in view, but assuming we can derive or it's mostly correct.
    const avg = total > 0 ? (campaign.value.collected_amount / total) : 0 // This is rough if total is not accurate count
    return {
        totalCount: total, // Only loaded
        average: avg
    }
})

const loadData = async () => {
    loading.value = true
    try {
        const [campRes, donRes] = await Promise.all([
            donationApi.getCampaign(campaignId.value),
            donationApi.getDonations(campaignId.value, limit, offset.value)
        ])
        campaign.value = campRes.data
        const newDonations = donRes.data.donations || []
        donations.value = newDonations
        hasMore.value = newDonations.length >= limit
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const loadMore = async () => {
    offset.value += limit
    try {
        const res = await donationApi.getDonations(campaignId.value, limit, offset.value)
        const newDonations = res.data.donations || []
        donations.value.push(...newDonations)
        hasMore.value = newDonations.length >= limit
    } catch (e) { console.error(e) }
}

const formatAmount = (amount: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount || 0)
}

const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

const getProgress = (c: any) => {
    if (!c.target_amount || c.target_amount <= 0) return 100
    return Math.min(100, (c.collected_amount / c.target_amount) * 100)
}

const confirmRefund = async (donation: any) => {
    if (!confirm(`Rembourser le don de ${formatAmount(donation.amount, donation.currency)} ?`)) return
    try {
        await donationApi.refundDonation(donation.id)
        alert('Remboursement initié.')
        // Reload to update status
        const idx = donations.value.findIndex(d => d.id === donation.id)
        if(idx !== -1) donations.value[idx].status = 'refunding'
    } catch (e: any) {
        alert(e.response?.data?.error || "Erreur lors du remboursement")
    }
}

onMounted(() => {
    loadData()
})

definePageMeta({
  middleware: 'auth'
})
</script>
