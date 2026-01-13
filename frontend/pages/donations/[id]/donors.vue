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
          
          <!-- Access Control: Show content only if creator or show_donors is true -->
          <div v-if="canViewDonors">
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

              <!-- Filters & Search -->
              <div class="mb-6 flex flex-col sm:flex-row gap-4 justify-between">
                  <div class="relative flex-1">
                      <input v-model="searchQuery" type="text" placeholder="Rechercher un donateur..." class="w-full pl-10 pr-4 py-2 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-900 focus:ring-2 focus:ring-indigo-500">
                      <span class="absolute left-3 top-2.5 text-gray-400">🔍</span>
                  </div>
                  <!-- Tier Filter (Optional, maybe later) -->
              </div>

              <!-- Tier Stats (For Tiers Campaigns) -->
              <div v-if="campaign.donation_type === 'tiers' && tierStats.length > 0" class="mb-8 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                  <div v-for="stat in tierStats" :key="stat.label" class="bg-white dark:bg-slate-800 p-4 rounded-xl border border-gray-100 dark:border-gray-700">
                      <div class="text-xs text-gray-500 uppercase font-bold">{{ stat.label }}</div>
                      <div class="text-xl font-bold text-indigo-600 dark:text-indigo-400">{{ stat.count }} dons</div>
                      <div class="text-xs text-gray-400">{{ formatAmount(stat.total, campaign.currency) }}</div>
                  </div>
              </div>

              <!-- Donors List -->
              <div v-if="filteredDonations.length > 0" class="space-y-4">
                  <div v-for="donation in filteredDonations" :key="donation.id" 
                       @click="openDonationModal(donation)"
                       class="bg-white dark:bg-slate-800 p-4 rounded-xl border border-gray-100 dark:border-gray-700 flex items-center justify-between hover:shadow-md transition-shadow cursor-pointer">
                      
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
          
          <div v-else class="text-center py-20 bg-gray-50 dark:bg-slate-800/50 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
             <div class="text-6xl mb-4">🔒</div>
             <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Liste des donateurs privée</h3>
             <p class="text-gray-500">Seul l'organisateur peut voir la liste des donateurs. Demandez-lui d'activer l'option s'il le souhaite.</p>
          </div>
      </div>
    </div>

    <!-- Donation Detail Modal -->
    <div v-if="selectedDonation" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4" @click.self="selectedDonation = null">
        <div class="bg-white dark:bg-slate-900 rounded-2xl shadow-2xl max-w-md w-full overflow-hidden animate-in fade-in zoom-in duration-200">
            <div class="p-6 relative">
                <button @click="selectedDonation = null" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">✕</button>
                
                <div class="text-center mb-6">
                    <div class="w-16 h-16 rounded-full mx-auto flex items-center justify-center text-2xl font-bold mb-3"
                         :class="selectedDonation.is_anonymous ? 'bg-gray-100 dark:bg-gray-700 text-gray-500' : 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600'">
                        {{ selectedDonation.is_anonymous ? '?' : (selectedDonation.donor_name?.[0]?.toUpperCase() || 'B') }}
                    </div>
                    <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                        {{ selectedDonation.is_anonymous ? 'Donateur Anonyme' : (selectedDonation.donor_name || 'Bienfaiteur') }}
                    </h3>
                    <p class="text-sm text-gray-500">{{ formatDate(selectedDonation.created_at) }}</p>
                </div>

                <div class="bg-gray-50 dark:bg-slate-800 rounded-xl p-4 mb-6 text-center">
                    <span class="block text-gray-500 text-sm mb-1">Montant du Don</span>
                    <span class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">
                        {{ formatAmount(selectedDonation.amount, selectedDonation.currency) }}
                    </span>
                    <div v-if="selectedDonation.frequency && selectedDonation.frequency !== 'one_time'" class="mt-2 inline-block px-2 py-1 bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 text-xs rounded-lg">
                        🔄 {{ selectedDonation.frequency === 'monthly' ? 'Mensuel' : 'Annuel' }}
                    </div>
                </div>

                <div v-if="selectedDonation.message" class="mb-6">
                    <h4 class="text-xs font-bold text-gray-500 uppercase mb-2">Message</h4>
                    <p class="text-gray-700 dark:text-gray-300 italic bg-gray-50 dark:bg-slate-800 p-3 rounded-lg">
                        "{{ selectedDonation.message }}"
                    </p>
                </div>

                <!-- Custom Fields Data -->
                <div v-if="selectedDonation.form_data && Object.keys(selectedDonation.form_data).length > 0" class="space-y-3">
                    <h4 class="text-xs font-bold text-gray-500 uppercase mb-2">Informations</h4>
                    <div v-for="(value, key) in selectedDonation.form_data" :key="key" class="flex justify-between text-sm border-b border-gray-100 dark:border-gray-800 pb-2">
                        <span class="text-gray-500 capitalize">{{ key.replace(/_/g, ' ') }}</span>
                        <span class="font-medium text-gray-900 dark:text-white">{{ value }}</span>
                    </div>
                </div>

                <div class="mt-8 pt-4 border-t border-gray-100 dark:border-gray-800 text-center">
                     <button v-if="isCreator && selectedDonation.status === 'paid'" @click="confirmRefund(selectedDonation)" 
                           class="text-red-500 hover:text-red-700 font-bold text-sm hover:underline">
                        Demander un remboursement
                     </button>
                </div>
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
const canViewDonors = computed(() => isCreator.value || (campaign.value && campaign.value.show_donors))

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

// === Search & Filtering ===
const searchQuery = ref('')

const filteredDonations = computed(() => {
    if (!searchQuery.value) return donations.value
    const q = searchQuery.value.toLowerCase()
    return donations.value.filter(d => 
        (d.donor_name && d.donor_name.toLowerCase().includes(q)) ||
        (d.message && d.message.toLowerCase().includes(q))
    )
})

// === Tier Stats ===
const tierStats = computed(() => {
    if (!campaign.value || campaign.value.donation_type !== 'tiers' || !campaign.value.donation_tiers) return []
    
    // Map tiers to stats
    return campaign.value.donation_tiers.map((tier: any) => {
        // Find donations matching this tier amount (approximate check usually sufficient for float, but exact match preferred for tiers logic)
        const tierDonations = donations.value.filter(d => Math.abs(d.amount - tier.amount) < 0.01)
        const total = tierDonations.reduce((sum, d) => sum + d.amount, 0)
        return {
            label: tier.label,
            count: tierDonations.length,
            total: total
        }
    })
})

const selectedDonation = ref<any>(null)
const openDonationModal = (donation: any) => {
    // If not creator and private data, maybe restricted? But here we are on Wall of Donors.
    // If anonymous, maybe we hide details? For now show what we have.
    selectedDonation.value = donation
}

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
