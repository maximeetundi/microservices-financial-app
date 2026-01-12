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
                    
                    <!-- Right Column (Progress & Share) (Desktop) -->
                    <div class="w-full md:w-80 flex flex-col gap-6">
                        <!-- Progress Card -->
                        <div class="bg-gray-50 dark:bg-slate-800/50 rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
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

                        <!-- Share Card (Creator Only) -->
                        <div v-if="isCreator" class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800 shadow-sm">
                            <h4 class="font-bold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                               <span>📱</span> Partager (Visible par vous seul)
                            </h4>
                            
                            <!-- QR Code -->
                            <div v-if="campaign.qr_code_url" class="flex justify-center mb-4 bg-white p-2 rounded-xl border border-gray-100">
                                 <img :src="campaign.qr_code_url" class="w-32 h-32 object-contain" alt="QR Code">
                            </div>

                            <!-- Code -->
                            <div v-if="campaign.code" class="bg-gray-50 dark:bg-slate-800 p-3 rounded-xl flex justify-between items-center mb-4">
                                <span class="font-mono font-bold text-indigo-600 dark:text-indigo-400 text-lg">{{ campaign.code }}</span>
                                <button @click="copyCode(campaign.code)" class="text-xs bg-white dark:bg-slate-700 text-gray-700 dark:text-gray-300 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 hover:bg-gray-50 transition-colors">Copier</button>
                            </div>
                            
                            <p class="text-xs text-center text-gray-500 leading-relaxed">
                                Scannez ce code pour tester ou partagez-le avec vos donateurs.
                            </p>
                        </div>
                    </div>
                </div>

                <!-- Description -->
                <div class="prose dark:prose-invert max-w-none mb-12">
                    <h3 class="text-xl font-bold mb-4">À propos</h3>
                    <p class="whitespace-pre-line text-gray-600 dark:text-gray-300 leading-relaxed">{{ campaign.description }}</p>
                </div>

                <!-- Wall of Donors (Refined) -->
                <div>
                    <h3 class="text-2xl font-bold mb-8 flex items-center gap-3 text-gray-900 dark:text-white">
                        <span class="bg-yellow-100 dark:bg-yellow-900/30 p-2 rounded-lg">🏆</span> Mur des donateurs
                    </h3>
                    
                    <div v-if="donations.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        <div v-for="donation in donations" :key="donation.id" 
                             class="group relative overflow-hidden p-6 rounded-2xl bg-white dark:bg-slate-800 border border-gray-100 dark:border-gray-700 hover:shadow-xl hover:-translate-y-1 transition-all duration-300">
                            
                            <div class="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-indigo-500/10 to-transparent rounded-bl-full -mr-6 -mt-6"></div>

                            <div class="flex items-start gap-4 mb-4">
                                <div class="w-12 h-12 rounded-full flex-shrink-0 flex items-center justify-center text-xl font-bold shadow-inner"
                                     :class="donation.is_anonymous ? 'bg-gray-100 dark:bg-gray-700 text-gray-500' : 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-600'">
                                    {{ donation.is_anonymous ? '?' : (donation.donor_name?.[0]?.toUpperCase() || 'B') }}
                                </div>
                                <div>
                                    <p class="font-bold text-gray-900 dark:text-white text-lg">
                                        {{ donation.is_anonymous ? 'Donateur Anonyme' : (donation.donor_name || 'Bienfaiteur') }}
                                    </p>
                                    <p class="text-sm text-gray-500">{{ formatDate(donation.created_at) }}</p>
                                </div>
                            </div>

                            <div class="flex justify-between items-end">
                                <div class="flex-1 pr-4">
                                    <p v-if="donation.message" class="text-sm text-gray-600 dark:text-gray-300 italic relative pl-3 border-l-2 border-indigo-200 dark:border-indigo-800 line-clamp-2">
                                        {{ donation.message }}
                                    </p>
                                </div>
                                <div class="font-bold text-xl text-emerald-600 tabular-nums">
                                    +{{ formatAmount(donation.amount, donation.currency) }}
                                </div>
                            </div>
                        </div>
                    </div>
                    <div v-else class="text-center py-12 px-6 bg-gray-50 dark:bg-slate-800/30 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
                        <div class="text-4xl mb-3">🌱</div>
                        <h4 class="text-lg font-medium text-gray-900 dark:text-white">Soyez le premier à planter une graine !</h4>
                        <p class="text-gray-500">Votre soutien peut déclencher une vague de générosité.</p>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Donate Modal -->
    <Teleport to="body">
        <div v-if="showDonateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showDonateModal = false">
            <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-md shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200 border border-gray-100 dark:border-gray-800 flex flex-col max-h-[90vh]">
                <div class="p-6 overflow-y-auto custom-scrollbar"> 
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

                    <!-- Anonymous Toggle -->
                    <div class="mb-6">
                        <label class="flex items-center gap-3 p-3 rounded-xl border border-gray-200 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                            <input type="checkbox" v-model="donationForm.isAnonymous" class="w-5 h-5 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                            <div>
                                <span class="block font-medium text-gray-900 dark:text-white">Rester anonyme</span>
                                <span class="block text-xs text-gray-500">Masquer mon nom et ignorer les détails facultatifs.</span>
                            </div>
                        </label>
                    </div>

                    <!-- Dynamic Fields (Hidden if Anonymous) -->
                    <div v-if="campaign?.form_schema?.length && !donationForm.isAnonymous" class="mb-6 border-t border-gray-100 dark:border-gray-800 pt-4 animate-in slide-in-from-top-2">
                        <h4 class="font-bold text-sm mb-4 text-gray-900 dark:text-white">Informations complémentaires</h4>
                        <div v-for="field in campaign.form_schema" :key="field.name" class="mb-4">
                            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                                {{ field.label }} <span v-if="field.required" class="text-red-500">*</span>
                            </label>
                            
                            <select v-if="field.type === 'select'" v-model="donationForm.formData[field.name]" class="w-full px-4 py-2 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500">
                                <option value="" disabled>Sélectionner</option>
                                <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                            </select>
                            
                            <input v-else :type="field.type === 'number' ? 'number' : 'text'" 
                                   v-model="donationForm.formData[field.name]"
                                   class="w-full px-4 py-2 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500"
                                   :placeholder="field.label">
                        </div>
                    </div>

                    <!-- Wallet Selection -->
                    <div class="mb-8">
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
                            {{ processing ? 'Confirmer le Don' : 'Confirmer le Don' }}
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
import { useAuth } from '~/composables/useAuth'
import { usePin } from '~/composables/usePin'

const { user } = useAuth()
const { requirePin } = usePin()
const route = useRoute()
const apiContext = useApi()
const { donationApi, walletApi } = apiContext

definePageMeta({
  middleware: 'auth'
})

const campaignId = computed(() => route.params.id as string)
const campaign = ref<any>(null)
const donations = ref<any[]>([])
const loading = ref(true)
const creatorName = ref('') 

const isCreator = computed(() => campaign.value && user.value && campaign.value.creator_id === user.value.id)

// Donate Modal State
const showDonateModal = ref(false)
const showSuccessModal = ref(false)
const processing = ref(false)
const wallets = ref<any[]>([])
const lastDonationAmount = ref(0)
const walletsLoaded = ref(false)

const donationForm = reactive({
    amount: null as number | null,
    frequency: 'one_time',
    isAnonymous: false,
    message: '',
    walletId: '',
    formData: {} as Record<string, any>
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
    showDonateModal.value = true
    if (!walletsLoaded.value) {
        try {
            const res = await walletApi.getMyWallets() 
            wallets.value = res.data?.wallets || res.data || []
            walletsLoaded.value = true
             if (wallets.value.length > 0) {
                 donationForm.walletId = wallets.value[0].id
            }
        } catch(e) { console.error(e) }
    }
}

const isValidDonation = computed(() => {
    return donationForm.amount && donationForm.amount > 0 && donationForm.walletId
})

const processDonation = async () => {
    if (!isValidDonation.value) return

    // Validate Dynamic Fields (Skip if Anonymous)
    if (campaign.value.form_schema && !donationForm.isAnonymous) {
        for (const field of campaign.value.form_schema) {
             if (field.required) {
                 const val = donationForm.formData[field.name]
                 if (!val || val.toString().trim() === '') {
                     alert(`Le champ "${field.label}" est obligatoire.`)
                     return
                 }
             }
        }
    }

    // Use Global PIN Modal
    requirePin(async (pin) => {
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
                form_data: donationForm.isAnonymous ? {} : donationForm.formData, 
                pin: pin 
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
    })
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
    return new Date(date).toLocaleDateString('fr-FR')
}

const copyCode = (code: string) => {
    navigator.clipboard.writeText(code)
    alert('Code copié !')
}

onMounted(() => {
    loadData()
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: rgba(156, 163, 175, 0.5); border-radius: 20px; }
</style>
