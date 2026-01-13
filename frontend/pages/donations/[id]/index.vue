<template>
  <NuxtLayout name="dashboard">
    <div class="event-detail-page"> 
      
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center py-20">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="!campaign && !loading" class="text-center py-20">
        <div class="text-6xl mb-4">❌</div>
        <h2 class="text-2xl font-bold mb-4">Campagne non trouvée</h2>
        <NuxtLink to="/donations" class="text-indigo-600 hover:underline">Retour aux dons</NuxtLink>
      </div>
    
      <template v-else-if="campaign">
        <!-- Hero Section -->
        <div class="relative h-64 md:h-80 w-full overflow-hidden rounded-b-3xl mb-8 group">
          <div class="absolute inset-0 bg-slate-900/50 z-10"></div>
          <img v-if="campaign.image_url" :src="campaign.image_url" class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" alt="Campaign">
          <div v-else class="w-full h-full bg-gradient-to-r from-blue-600 to-indigo-700 flex items-center justify-center text-8xl">
             🤲
          </div>
          
          <div class="absolute inset-0 z-20 container mx-auto px-4 flex flex-col justify-end pb-8">
             <NuxtLink to="/donations" class="text-white/80 hover:text-white mb-4 flex items-center gap-2 transition-colors w-fit">
                ← Retour
             </NuxtLink>
             
             <div class="flex flex-wrap items-end justify-between gap-4">
                 <div>
                    <span class="inline-block px-3 py-1 rounded-full text-xs font-bold uppercase mb-3"
                          :class="campaign.status === 'active' ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30' : 'bg-red-500/20 text-red-300 border border-red-500/30'">
                        {{ campaign.status === 'active' ? 'En cours' : (campaign.status === 'cancelled' ? 'Annulée' : 'Terminée') }}
                    </span>
                    <h1 class="text-3xl md:text-5xl font-bold text-white mb-2">{{ campaign.title }}</h1>
                    <div class="flex items-center gap-3 text-white/90">
                        <span class="flex items-center gap-1">📅 {{ formatDate(campaign.created_at) }}</span>
                        <span class="w-1 h-1 bg-white/50 rounded-full"></span>
                        <span class="flex items-center gap-1">👤 Par {{ creatorName }}</span>
                    </div>
                 </div>
             </div>
          </div>
        </div>

        <div class="container mx-auto px-4 pb-20">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <!-- Left: Info & Donors -->
                <div class="lg:col-span-2 space-y-8">
                    <!-- About -->
                    <section class="bg-white dark:bg-slate-800 rounded-3xl p-8 border border-gray-100 dark:border-gray-700 shadow-sm">
                        <h3 class="text-xl font-bold mb-4 flex items-center gap-2 text-gray-900 dark:text-white">
                            📝 À propos
                        </h3>
                        <p class="whitespace-pre-line text-gray-600 dark:text-gray-300 leading-relaxed text-lg">
                            {{ campaign.description }}
                        </p>
                    </section>

                    <!-- Creator Code Section -->
                    <section v-if="isValidCode" class="bg-indigo-50 dark:bg-slate-800/80 rounded-3xl p-8 border border-indigo-100 dark:border-indigo-900/30 relative overflow-hidden">
                        <div class="relative z-10">
                            <h3 class="text-xl font-bold mb-6 flex items-center gap-2 text-indigo-900 dark:text-indigo-100">
                                🔲 Partagez cette campagne
                            </h3>
                            <div class="flex flex-wrap gap-4 items-center">
                                <div class="bg-white dark:bg-slate-900 px-6 py-3 rounded-xl border border-indigo-200 dark:border-indigo-800 flex items-center gap-3 cursor-pointer hover:shadow-md transition-shadow" @click="copyCode(displayCode)">
                                    <span class="font-mono font-bold text-2xl text-indigo-600 dark:text-indigo-400 tracking-wider">{{ displayCode }}</span>
                                    <span class="text-gray-400 text-sm">📋</span>
                                </div>
                                <button @click="openQRModal" class="px-5 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-600/20 flex items-center gap-2">
                                    🔍 QR Code & Options
                                </button>
                            </div>
                            <p class="mt-4 text-sm text-indigo-600/70 dark:text-indigo-400">
                                Utilisez ce code pour recevoir des dons directement.
                            </p>
                        </div>
                        <!-- Decorative background -->
                        <div class="absolute -right-10 -bottom-10 w-48 h-48 bg-indigo-500/10 rounded-full blur-3xl"></div>
                    </section>
                    
                    <!-- View Donors Button -->
                    <section>
                         <button @click="navigateTo(`/donations/${campaign.id}/donors`)" 
                                 class="w-full py-6 bg-white dark:bg-slate-800 rounded-3xl border border-gray-100 dark:border-gray-700 shadow-sm hover:border-indigo-200 dark:hover:border-indigo-800 transition-all group flex items-center justify-between px-8">
                             <div class="flex items-center gap-4">
                                 <span class="w-12 h-12 rounded-full bg-yellow-100 dark:bg-yellow-900/30 flex items-center justify-center text-2xl group-hover:scale-110 transition-transform">🏆</span>
                                 <div class="text-left">
                                     <h3 class="text-xl font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">Mur des Donateurs</h3>
                                     <p class="text-gray-500">Voir la liste complète et les statistiques</p>
                                 </div>
                             </div>
                             <span class="text-2xl text-gray-300 group-hover:translate-x-2 transition-transform">→</span>
                         </button>
                    </section>
                </div>

                <!-- Right: Actions & Progress -->
                <div class="space-y-6">
                    <!-- Progress Card -->
                    <div class="bg-white dark:bg-slate-800 rounded-3xl p-6 border border-gray-100 dark:border-gray-700 shadow-xl shadow-indigo-900/5 sticky top-24">
                        <div class="mb-6">
                            <div class="flex items-end gap-2 mb-1">
                                <span class="text-4xl font-bold text-indigo-600 dark:text-indigo-400">{{ formatAmount(campaign.collected_amount, campaign.currency) }}</span>
                                <span class="text-gray-500 dark:text-gray-400 mb-1.5 font-medium">récoltés</span>
                            </div>
                             <div class="w-full bg-gray-100 dark:bg-gray-700 h-4 rounded-full overflow-hidden">
                                 <div class="bg-gradient-to-r from-indigo-500 to-purple-600 h-full rounded-full transition-all duration-1000 ease-out" 
                                      :style="{ width: getProgress(campaign) + '%' }"></div>
                            </div>
                            <div class="flex justify-between mt-2 text-sm font-medium">
                                <span class="text-gray-500">Objectif: {{ campaign.target_amount > 0 ? formatAmount(campaign.target_amount, campaign.currency) : 'Illimité' }}</span>
                                <span class="text-indigo-600 dark:text-indigo-400">{{ Math.round(getProgress(campaign)) }}%</span>
                            </div>
                        </div>

                        <div class="space-y-3">
                            <button @click="openDonateModal" class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg shadow-indigo-600/30 transition-all hover:-translate-y-1 flex justify-center items-center gap-2 text-lg">
                                <span>❤️</span> Faire un don
                            </button>
    
                            <button @click="shareCampaign" class="w-full py-3 bg-white dark:bg-slate-700 text-indigo-600 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-900/50 rounded-xl font-bold hover:bg-indigo-50 dark:hover:bg-slate-600 transition-colors flex justify-center items-center gap-2">
                                <span>🔗</span> Partager
                            </button>
                        </div>

                         <!-- Organizer Actions -->
                        <div v-if="isCreator" class="mt-8 pt-6 border-t border-gray-100 dark:border-gray-700">
                             <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-4">Espace Organisateur</h4>
                             
                             <div class="grid grid-cols-2 gap-3 mb-3">
                                 <NuxtLink :to="`/donations/${campaign.id}/edit`" class="py-2 px-3 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-200 rounded-lg font-bold text-sm text-center hover:bg-gray-200 transition-colors">
                                     ✏️ Modifier
                                 </NuxtLink>
                                 <button @click="navigateTo('/scan')" class="py-2 px-3 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-200 rounded-lg font-bold text-sm hover:bg-gray-200 transition-colors">
                                     📷 Scanner (QR)
                                 </button>
                             </div>

                             <div v-if="campaign.status === 'active'" class="space-y-2">
                                <button @click="confirmCancelCampaign" class="w-full py-2 px-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg font-bold text-sm hover:bg-red-100 transition-colors flex justify-center items-center gap-2">
                                    🛑 Annuler la campagne
                                </button>
                             </div>
                             
                             <div v-else class="text-center p-2 mb-2 bg-red-50 dark:bg-red-900/20 text-red-600 rounded-lg text-sm font-bold">
                                 Statut: {{ campaign.status === 'cancelled' ? 'Annulée' : 'Terminée' }}
                             </div>

                             <button @click="confirmDelete" class="w-full mt-2 py-2 text-gray-400 hover:text-red-500 text-xs font-bold transition-colors">
                                 🗑️ Supprimer définitivement
                             </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
      </template>

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
                     <div v-if="campaign?.form_schema?.length && !donationForm.isAnonymous" class="mb-6 border-t border-gray-100 dark:border-gray-800 pt-4">
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
                            {{ processing ? 'Confirmer' : 'Confirmer le Don' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
      </Teleport>

      <!-- QR Code Modal -->
      <Teleport to="body">
        <div v-if="showQRModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm" @click.self="showQRModal = false">
            <div class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-sm overflow-hidden animate-in fade-in zoom-in duration-300">
                <div class="p-4 bg-gray-50 dark:bg-slate-800 flex justify-end">
                    <button @click="showQRModal = false" class="text-gray-500 hover:text-gray-700 dark:hover:text-white transition-colors">✕</button>
                </div>
                
                <div class="p-8 text-center">
                    <div class="mb-8 relative inline-block">
                        <div class="bg-white p-4 rounded-2xl shadow-xl">
                            <!-- Generated QR Code -->
                            <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="QR Code" class="w-48 h-48 object-contain">
                            <div v-else class="w-48 h-48 bg-gray-100 flex items-center justify-center text-gray-300 text-6xl rounded-lg">📱</div>
                        </div>
                    </div>

                    <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">{{ campaign?.title }}</h3>
                    <p class="text-gray-500 text-sm mb-8 font-mono bg-gray-100 dark:bg-slate-800 py-2 rounded-lg select-all">
                        {{ displayCode }}
                    </p>

                    <div class="grid grid-cols-2 gap-3">
                        <button @click="copyCode(displayCode)" class="py-3 bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-xl font-bold hover:bg-gray-200 transition-colors flex items-center justify-center gap-2">
                            <span>📋</span> Copier
                        </button>
                         <button @click="downloadQRCode" class="py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors flex items-center justify-center gap-2">
                            <span>⬇️</span> PNG
                        </button>
                         <button @click="shareCampaign" class="col-span-2 py-3 bg-white border border-gray-200 dark:bg-slate-800 dark:border-gray-700 text-indigo-600 dark:text-indigo-400 rounded-xl font-bold hover:bg-gray-50 transition-colors flex items-center justify-center gap-2">
                            <span>📤</span> Partager le lien
                        </button>
                    </div>
                </div>
            </div>
        </div>
      </Teleport>

      <!-- Success Modal (Simplified) -->
        <Teleport to="body">
            <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md" @click.self="showSuccessModal = false">
                <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-sm shadow-2xl p-6 border border-gray-100 dark:border-gray-800 animate-in fade-in zoom-in duration-300 text-center">
                    <div class="w-20 h-20 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 flex items-center justify-center mb-6 mx-auto animate-bounce">
                        <span class="text-4xl">🎉</span>
                    </div>
                    <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Merci pour votre Don !</h3>
                    <p class="text-gray-500 dark:text-gray-400 mb-8">
                        Le montant de <strong>{{ formatAmount(lastDonationAmount, campaign?.currency) }}</strong> a été envoyé avec succès.
                    </p>
                    <button @click="closeSuccess" class="w-full px-4 py-3 bg-emerald-600 text-white rounded-xl hover:bg-emerald-700 transition-colors font-bold shadow-lg shadow-emerald-600/20">
                        Génial !
                    </button>
                </div>
            </div>
        </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useApi } from '~/composables/useApi'
import { useAuthStore } from '~/stores/auth'
import { usePin } from '~/composables/usePin'
import QRCode from 'qrcode'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const { requirePin } = usePin()
const route = useRoute()
const router = useRouter()
const apiContext = useApi()
const { donationApi, walletApi } = apiContext

definePageMeta({
  middleware: 'auth'
})

const campaignId = computed(() => route.params.id as string)
const campaign = ref<any>(null)
const loading = ref(true)
const creatorName = ref('') 
const qrCodeUrl = ref('')

const isCreator = computed(() => campaign.value && user.value && campaign.value.creator_id === user.value.id)

const displayCode = computed(() => {
    return campaign.value?.campaign_code || campaign.value?.id || ''
})

const isValidCode = computed(() => !!displayCode.value)

// Donate Modal State
const showDonateModal = ref(false)
const showQRModal = ref(false)
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
        const campRes = await donationApi.getCampaign(campaignId.value)
        campaign.value = campRes.data

        // Always generate QR code from URL for easy scanning
        generateQRCode(window.location.href)

        // Fetch Creator Name
        if (campaign.value.creator_id) {
            try {
                const userRes = await apiContext.userApi.getById(campaign.value.creator_id)
                const u = userRes.data
                creatorName.value = u.first_name || u.last_name ? `${u.first_name} ${u.last_name}` : u.username || `Utilisateur ${campaign.value.creator_id.slice(0,6)}`
            } catch (e) {
                console.warn("Could not fetch creator info", e)
            }
        }
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const generateQRCode = async (text: string) => {
    try {
        qrCodeUrl.value = await QRCode.toDataURL(text, { width: 300, margin: 2, color: { dark: '#4f46e5', light: '#ffffff' } })
    } catch (e) {
        console.error("QR Env err", e)
    }
}

const openQRModal = () => {
    showQRModal.value = true
    if (!qrCodeUrl.value) {
        generateQRCode(window.location.href)
    }
}

const shareCampaign = async () => {
    const url = window.location.href
    const title = campaign.value?.title || 'Soutenez cette campagne !'
    
    if (navigator.share) {
        try {
            await navigator.share({
                title: title,
                text: campaign.value?.description?.slice(0, 100) + '...',
                url: url
            })
        } catch (e) { console.log('Share cancelled') }
    } else {
        copyCode(url)
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

    // Dynamic fields validation
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
            loadData() 
            
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
    if (!code) return alert("Code non disponible")
    navigator.clipboard.writeText(code)
    alert('Code copié !')
}

const downloadQRCode = () => {
  if (!qrCodeUrl.value) {
    alert('QR Code non disponible')
    return
  }
  
  const link = document.createElement('a')
  link.href = qrCodeUrl.value
  link.download = `campaign-${displayCode.value}.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const confirmCancelCampaign = async () => {
    if (!confirm('Êtes-vous sûr de vouloir annuler cette campagne ? Cela remboursera automatiquement tous les donateurs. C\'est irréversible.')) return

    try {
        await donationApi.cancelCampaign(campaign.value.id)
        alert('Campagne annulée initée. Les remboursement sont en cours.')
        loadData()
    } catch (e: any) {
        alert(e.response?.data?.error || "Erreur lors de l'annulation")
    }
}

const confirmDelete = async () => {
     if (!confirm('ATTENTION: Supprimer définitivement cette campagne ? Cette action est irréversible.')) return
     // TODO: Implement actual delete if supported by backend
     alert("Fonctionnalité de suppression en cours de développement.")
}

onMounted(() => {
    loadData()
})
</script>
