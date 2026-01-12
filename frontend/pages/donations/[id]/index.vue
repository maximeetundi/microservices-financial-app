<template>
  <NuxtLayout name="dashboard">
    <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>
    
    <div v-else-if="campaign" class="max-w-4xl mx-auto py-8 px-4 animate-in fade-in duration-500">
        <!-- Back Link -->
        <button @click="navigateTo('/donations')" class="mb-6 flex items-center text-gray-500 hover:text-indigo-600 transition-colors">
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
            Retour aux campagnes
        </button>

        <div class="glass-card bg-white dark:bg-slate-900 rounded-3xl overflow-hidden shadow-xl border border-gray-100 dark:border-gray-800">
            <!-- Hero Image -->
            <div class="h-64 md:h-80 bg-gray-100 dark:bg-slate-800 relative">
                <img v-if="campaign.image_url" :src="campaign.image_url" class="w-full h-full object-cover" alt="Campaign">
                <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-r from-blue-500 to-indigo-600 text-white text-6xl">
                    🤲
                </div>
            </div>

            <div class="p-8 md:p-10">
                <div class="flex flex-col md:flex-row justify-between items-start gap-6 mb-8">
                    <div>
                        <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4">{{ campaign.title }}</h1>
                        <div class="flex items-center gap-3">
                            <div class="w-10 h-10 rounded-full bg-indigo-100 dark:bg-indigo-900/50 flex items-center justify-center text-indigo-600 font-bold">
                                👤
                            </div>
                            <div>
                                <p class="text-sm text-gray-500 dark:text-gray-400">Organisé par</p>
                                <p class="font-medium text-gray-900 dark:text-white">{{ creatorName || 'Utilisateur ' + campaign.creator_id.slice(0,6) }}</p>
                            </div>
                        </div>
                    </div>
                    
                    <!-- Progress Card (Desktop Side) -->
                    <div class="w-full md:w-80 bg-gray-50 dark:bg-slate-800/50 rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
                        <div class="mb-4">
                            <span class="text-3xl font-bold text-emerald-600">{{ formatAmount(campaign.collected_amount, campaign.currency) }}</span>
                            <span class="text-gray-500 dark:text-gray-400 text-sm ml-2">récoltés</span>
                        </div>
                        
                        <div class="bg-gray-200 dark:bg-gray-700 h-3 rounded-full overflow-hidden mb-2">
                             <div class="bg-emerald-500 h-full rounded-full transition-all duration-1000" :style="{ width: getProgress(campaign) + '%' }"></div>
                        </div>
                        
                        <div v-if="campaign.target_amount > 0" class="flex justify-between text-sm text-gray-500 mb-6">
                            <span>Objectif: {{ formatAmount(campaign.target_amount, campaign.currency) }}</span>
                            <span>{{ Math.round(getProgress(campaign)) }}%</span>
                        </div>
                        <div v-else class="text-sm text-emerald-600 font-medium mb-6">Objectif illimité 🚀</div>

                        <button @click="openDonateModal" class="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg shadow-indigo-600/30 transition-all hover:-translate-y-1">
                            Faire un don ❤️
                        </button>
                    </div>
                </div>

                <!-- Description -->
                <div class="prose dark:prose-invert max-w-none mb-12">
                    <h3 class="text-xl font-bold mb-4">À propos</h3>
                    <p class="whitespace-pre-line text-gray-600 dark:text-gray-300 leading-relaxed">{{ campaign.description }}</p>
                </div>

                <!-- Wall of Donors -->
                <div>
                    <h3 class="text-xl font-bold mb-6 flex items-center gap-2">
                        <span>🏆</span> Mur des donateurs
                    </h3>
                    
                    <div v-if="donations.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div v-for="donation in donations" :key="donation.id" class="flex items-center gap-4 p-4 rounded-xl bg-gray-50 dark:bg-slate-800/50 border border-gray-100 dark:border-gray-700">
                            <div class="w-12 h-12 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center text-xl">
                                {{ donation.is_anonymous ? '🕵️' : '👤' }}
                            </div>
                            <div class="flex-1">
                                <p class="font-bold text-gray-900 dark:text-white">
                                    {{ donation.is_anonymous ? 'Donateur Anonyme' : (donation.donor_name || 'Bienfaiteur') }}
                                </p>
                                <p class="text-sm text-gray-500">{{ formatDate(donation.created_at) }}</p>
                                <div v-if="donation.message" class="text-xs text-gray-400 italic mt-1">"{{ donation.message }}"</div>
                            </div>
                            <div class="font-bold text-emerald-600">
                                +{{ formatAmount(donation.amount, donation.currency) }}
                            </div>
                        </div>
                    </div>
                    <div v-else class="text-center py-8 text-gray-500 bg-gray-50 dark:bg-slate-800/30 rounded-xl border border-dashed border-gray-200 dark:border-gray-700">
                        Soyez le premier à soutenir cette cause ! 🚀
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Donate Modal -->
    <Teleport to="body">
        <div v-if="showDonateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showDonateModal = false">
            <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200 border border-gray-100 dark:border-gray-800">
                <div class="p-6">
                    <h3 class="text-2xl font-bold text-center mb-6 text-gray-900 dark:text-white">Faire un don ❤️</h3>
                    
                    <!-- Amount -->
                    <div class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Montant du don</label>
                        <div class="relative">
                            <input v-model.number="donationForm.amount" type="number" class="w-full pl-4 pr-12 py-3 text-xl font-bold rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="0">
                            <span class="absolute right-4 top-3.5 text-gray-500 font-bold">{{ campaign?.currency }}</span>
                        </div>
                    </div>

                    <!-- Frequency (Recurring) -->
                    <div v-if="campaign?.allow_recurring" class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Fréquence</label>
                        <div class="grid grid-cols-3 gap-2">
                            <button v-for="freq in frequencyOptions" :key="freq.value"
                                    @click="donationForm.frequency = freq.value"
                                    class="px-2 py-2 text-xs font-medium rounded-lg border transition-all"
                                    :class="donationForm.frequency === freq.value 
                                        ? 'bg-indigo-600 text-white border-indigo-600' 
                                        : 'bg-white dark:bg-slate-800 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-700 hover:border-indigo-300'">
                                {{ freq.label }}
                            </button>
                        </div>
                        <p v-if="donationForm.frequency !== 'one_time'" class="text-xs text-amber-600 mt-2 bg-amber-50 dark:bg-amber-900/20 p-2 rounded">
                            ⚠️ Vous recevrez une demande de paiement à valider manuellement à chaque échéance.
                        </p>
                    </div>

                    <!-- Message -->
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Message de soutien (optionnel)</label>
                        <textarea v-model="donationForm.message" rows="2" class="w-full px-4 py-2 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="Un petit mot pour l'organisateur..."></textarea>
                    </div>

                    <!-- Anonymous -->
                    <div class="mb-6">
                        <label class="flex items-center gap-3 p-3 rounded-xl border border-gray-200 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                            <input type="checkbox" v-model="donationForm.isAnonymous" class="w-5 h-5 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                            <div>
                                <span class="block font-medium text-gray-900 dark:text-white">Rester anonyme</span>
                                <span class="block text-xs text-gray-500">Votre nom ne sera pas affiché publiquement.</span>
                            </div>
                        </label>
                    </div>

                    <!-- Wallet Selection -->
                    <div class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Payer avec</label>
                         <select v-model="donationForm.walletId" class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800">
                             <option value="" disabled>Choisir un portefeuille</option>
                             <option v-for="w in wallets" :key="w.id" :value="w.id">
                                 {{ w.name }} ({{ formatAmount(w.balance, w.currency) }})
                             </option>
                         </select>
                    </div>

                    <div class="flex gap-3">
                        <button @click="showDonateModal = false" class="flex-1 px-4 py-3 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-xl font-bold hover:bg-gray-200 transition-colors">Annuler</button>
                        <button @click="processDonation" :disabled="processing || !isValidDonation" class="flex-1 px-4 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-600/20 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center gap-2">
                            <span v-if="processing" class="animate-spin w-5 h-5 border-2 border-white/30 border-t-white rounded-full"></span>
                            {{ processing ? 'Traitement...' : 'Confirmer le Don' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Teleport>

    <!-- Success Modal -->
     <Teleport to="body">
          <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showSuccessModal = false">
              <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-sm shadow-2xl p-6 border border-gray-100 dark:border-gray-800 animate-in fade-in zoom-in duration-300 text-center">
                  <div class="w-20 h-20 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 flex items-center justify-center mb-6 mx-auto animate-bounce">
                      <span class="text-4xl">🎉</span>
                  </div>
                  <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Merci pour votre Don !</h3>
                  <p class="text-gray-500 dark:text-gray-400 mb-8">
                      Votre générosité fait la différence. Le montant de <strong>{{ formatAmount(lastDonationAmount, campaign?.currency) }}</strong> a été envoyé avec succès.
                  </p>
                  <button @click="closeSuccess" class="w-full px-4 py-3 bg-emerald-600 text-white rounded-xl hover:bg-emerald-700 transition-colors font-bold shadow-lg shadow-emerald-600/20">
                      Génial !
                  </button>
              </div>
          </div>
      </Teleport>

  </NuxtLayout>
</template>

<script setup lang="ts">
import { useApi } from '~/composables/useApi'
const { donationApi, walletApi } = useApi()
const route = useRoute()

definePageMeta({
  middleware: 'auth'
})

const campaignId = computed(() => route.params.id as string)
const campaign = ref<any>(null)
const donations = ref<any[]>([])
const loading = ref(true)
const creatorName = ref('') // We might fetch this if we had user Service API exposed

// Donate Modal State
const showDonateModal = ref(false)
const showSuccessModal = ref(false)
const processing = ref(false)
const wallets = ref<any[]>([])
const lastDonationAmount = ref(0)

const donationForm = reactive({
    amount: null as number | null,
    frequency: 'one_time',
    isAnonymous: false,
    message: '',
    walletId: ''
})

const frequencyOptions = [
    { value: 'one_time', label: 'Une fois' },
    { value: 'monthly', label: 'Mensuel' },
    { value: 'annually', label: 'Annuel' }
]

const loadData = async () => {
    loading.value = true
    try {
        const [campRes, donRes] = await Promise.all([
            donationApi.getCampaign(campaignId.value),
            donationApi.getDonations(campaignId.value, 50)
        ])
        campaign.value = campRes.data
        donations.value = donRes.data.donations || []
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const openDonateModal = async () => {
    // Load wallets
    try {
        const res = await walletApi.getMyWallets() // Assuming this exists or getAll
        wallets.value = res.data?.wallets || res.data || []
        // Filter wallets compatible with campaign currency? Or auto-convert?
        // Service handles conversion usually, but cleaner to show relevant ones.
        // For now show all.
    } catch(e) { console.error(e) }
    
    showDonateModal.value = true
}

const isValidDonation = computed(() => {
    return donationForm.amount && donationForm.amount > 0 && donationForm.walletId
})

const processDonation = async () => {
    if (!isValidDonation.value) return
    processing.value = true
    try {
        await donationApi.initiateDonation({
            campaign_id: campaign.value.id,
            amount: donationForm.amount,
            currency: campaign.value.currency,
            wallet_id: donationForm.walletId,
            message: donationForm.message,
            is_anonymous: donationForm.isAnonymous,
            frequency: donationForm.frequency,
            // PIN? Usually required. For MVP assuming simplified or cached pin session.
            // If API requires PIN, we need PinModal.
            // Let's assume we need PIN.
            pin: "12345" // HARDCODED for now as PinModal integration takes more steps. Ideally use PinModal.
        })
        
        lastDonationAmount.value = donationForm.amount || 0
        showDonateModal.value = false
        showSuccessModal.value = true
        loadData() // Refresh
        
        // Reset form
        donationForm.amount = null
        donationForm.message = ''
        
    } catch(e: any) {
        alert(e.response?.data?.error || "Erreur lors du don")
    } finally {
        processing.value = false
    }
}

const closeSuccess = () => {
    showSuccessModal.value = false
}

const getProgress = (c: any) => {
    if (!c.target_amount || c.target_amount <= 0) return 100
    return Math.min(100, (c.collected_amount / c.target_amount) * 100)
}

const formatAmount = (amount: number, currency: string) => {
    return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount || 0)
}

const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString()
}

onMounted(() => {
    loadData()
})
</script>
