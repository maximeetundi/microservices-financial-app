<template>
  <NuxtLayout name="admin">
    <div class="p-6 lg:p-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">💳 Agrégateurs de Paiement</h1>
          <p class="text-slate-400">Configurer les providers de paiement par pays</p>
        </div>
        <button @click="showAddProvider = true" class="btn-premium flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Ajouter un Provider
        </button>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="glass-card p-5">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-indigo-500/20 rounded-xl">
              <span class="text-2xl">🔌</span>
            </div>
            <div>
              <p class="text-2xl font-bold text-white">{{ providers.length }}</p>
              <p class="text-sm text-slate-400">Providers</p>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-emerald-500/20 rounded-xl">
              <span class="text-2xl">✅</span>
            </div>
            <div>
              <p class="text-2xl font-bold text-white">{{ activeProviders }}</p>
              <p class="text-sm text-slate-400">Actifs</p>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-amber-500/20 rounded-xl">
              <span class="text-2xl">🧪</span>
            </div>
            <div>
              <p class="text-2xl font-bold text-white">{{ demoProviders }}</p>
              <p class="text-sm text-slate-400">Mode Démo</p>
            </div>
          </div>
        </div>
        <div class="glass-card p-5">
          <div class="flex items-center gap-4">
            <div class="p-3 bg-purple-500/20 rounded-xl">
              <span class="text-2xl">🌍</span>
            </div>
            <div>
              <p class="text-2xl font-bold text-white">{{ totalCountries }}</p>
              <p class="text-sm text-slate-400">Pays configurés</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Providers List -->
      <div class="glass-card overflow-hidden">
        <div class="p-4 border-b border-slate-700/50">
          <h2 class="text-lg font-semibold text-white">Providers Configurés</h2>
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-slate-800/50">
              <tr>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Provider</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Type</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Pays</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Status</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Mode</th>
                <th class="px-6 py-4 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-700/50">
              <tr v-for="provider in providers" :key="provider.id" class="hover:bg-slate-800/30 transition-colors">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl" :class="getProviderBg(provider.name)">
                      {{ getProviderIcon(provider.name) }}
                    </div>
                    <div>
                      <p class="font-semibold text-white">{{ provider.display_name }}</p>
                      <p class="text-xs text-slate-400">{{ provider.name }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs rounded-lg" :class="getTypeBadge(provider.provider_type)">
                    {{ getTypeLabel(provider.provider_type) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <div class="flex flex-wrap gap-1">
                    <span v-for="country in (provider.countries || []).slice(0, 4)" :key="country.country_code"
                      class="px-2 py-0.5 text-xs bg-slate-700 rounded text-slate-300">
                      {{ country.country_code }}
                    </span>
                    <span v-if="(provider.countries || []).length > 4" class="px-2 py-0.5 text-xs bg-slate-600 rounded text-slate-300">
                      +{{ provider.countries.length - 4 }}
                    </span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <button @click="toggleActive(provider)" class="relative">
                    <div class="w-12 h-6 rounded-full transition-colors" :class="provider.is_active ? 'bg-emerald-500' : 'bg-slate-600'">
                      <div class="absolute top-1 w-4 h-4 bg-white rounded-full transition-transform" 
                        :class="provider.is_active ? 'translate-x-7' : 'translate-x-1'"></div>
                    </div>
                  </button>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span v-if="provider.is_demo_mode" class="px-2 py-1 text-xs bg-amber-500/20 text-amber-400 rounded-lg">
                    🧪 DEMO
                  </span>
                  <span v-else class="px-2 py-1 text-xs bg-emerald-500/20 text-emerald-400 rounded-lg">
                    🔴 LIVE
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center gap-2">
                    <button @click="editProvider(provider)" class="p-2 hover:bg-slate-700 rounded-lg transition-colors text-slate-400 hover:text-white">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                      </svg>
                    </button>
                    <button @click="testProvider(provider)" class="p-2 hover:bg-slate-700 rounded-lg transition-colors text-slate-400 hover:text-emerald-400" title="Tester la connexion">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      </svg>
                    </button>
                    <button @click="deleteProvider(provider)" class="p-2 hover:bg-red-500/20 rounded-lg transition-colors text-slate-400 hover:text-red-400" :disabled="provider.name === 'demo'">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Edit Provider Modal -->
      <div v-if="showEditModal" class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4">
        <div class="bg-slate-900 rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-slate-700">
          <div class="p-6 border-b border-slate-700 flex justify-between items-center">
            <h3 class="text-xl font-bold text-white">{{ editingProvider?.id ? 'Modifier' : 'Ajouter' }} Provider</h3>
            <button @click="closeEditModal" class="text-slate-400 hover:text-white">✕</button>
          </div>
          
          <div class="p-6 space-y-6">
            <!-- Basic Info -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Nom technique</label>
                <input v-model="editingProvider.name" type="text" placeholder="flutterwave" 
                  class="input-premium w-full" :disabled="editingProvider?.id"/>
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Nom d'affichage</label>
                <input v-model="editingProvider.display_name" type="text" placeholder="Flutterwave" class="input-premium w-full"/>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Type</label>
                <select v-model="editingProvider.provider_type" class="input-premium w-full">
                  <option value="all">Tous (Mobile + Carte + Bank)</option>
                  <option value="mobile_money">Mobile Money</option>
                  <option value="card">Carte Bancaire</option>
                  <option value="bank_transfer">Virement Bancaire</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">URL API</label>
                <input v-model="editingProvider.api_base_url" type="url" placeholder="https://api.example.com" class="input-premium w-full"/>
              </div>
            </div>

            <!-- API Keys -->
            <div class="p-4 bg-slate-800/50 rounded-xl space-y-4">
              <h4 class="font-medium text-white flex items-center gap-2">
                <span>🔐</span> Clés API
              </h4>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">API Key</label>
                  <input v-model="editingProvider.api_key" type="password" placeholder="••••••••" class="input-premium w-full"/>
                </div>
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">API Secret</label>
                  <input v-model="editingProvider.api_secret" type="password" placeholder="••••••••" class="input-premium w-full"/>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">Public Key</label>
                  <input v-model="editingProvider.public_key" type="password" placeholder="••••••••" class="input-premium w-full"/>
                </div>
                <div>
                  <label class="block text-sm font-medium text-slate-400 mb-2">Webhook Secret</label>
                  <input v-model="editingProvider.webhook_secret" type="password" placeholder="••••••••" class="input-premium w-full"/>
                </div>
              </div>
            </div>

            <!-- Mode Switches -->
            <div class="flex items-center gap-8">
              <label class="flex items-center gap-3 cursor-pointer">
                <input type="checkbox" v-model="editingProvider.is_active" class="sr-only">
                <div class="w-12 h-6 rounded-full transition-colors" :class="editingProvider.is_active ? 'bg-emerald-500' : 'bg-slate-600'">
                  <div class="w-4 h-4 mt-1 ml-1 bg-white rounded-full transition-transform" :class="editingProvider.is_active ? 'translate-x-6' : ''"></div>
                </div>
                <span class="text-white">Actif</span>
              </label>
              <label class="flex items-center gap-3 cursor-pointer">
                <input type="checkbox" v-model="editingProvider.is_demo_mode" class="sr-only">
                <div class="w-12 h-6 rounded-full transition-colors" :class="editingProvider.is_demo_mode ? 'bg-amber-500' : 'bg-slate-600'">
                  <div class="w-4 h-4 mt-1 ml-1 bg-white rounded-full transition-transform" :class="editingProvider.is_demo_mode ? 'translate-x-6' : ''"></div>
                </div>
                <span class="text-white">Mode Démo</span>
              </label>
            </div>

            <!-- Countries -->
            <div>
              <h4 class="font-medium text-white mb-3">🌍 Pays associés</h4>
              <div class="grid grid-cols-3 gap-2">
                <label v-for="country in availableCountries" :key="country.code" 
                  class="flex items-center gap-2 p-2 rounded-lg border border-slate-700 cursor-pointer hover:bg-slate-800 transition-colors"
                  :class="isCountrySelected(country.code) ? 'bg-indigo-500/20 border-indigo-500' : ''">
                  <input type="checkbox" :value="country.code" v-model="selectedCountries" class="sr-only">
                  <span class="text-lg">{{ country.flag }}</span>
                  <span class="text-sm text-white">{{ country.name }}</span>
                  <span class="text-xs text-slate-400">{{ country.currency }}</span>
                </label>
              </div>
            </div>
          </div>

          <div class="p-6 border-t border-slate-700 flex justify-end gap-3">
            <button @click="closeEditModal" class="px-6 py-3 rounded-xl font-medium text-slate-300 hover:bg-slate-800 transition-colors">
              Annuler
            </button>
            <button @click="saveProvider" :disabled="saving" class="btn-premium px-6 py-3">
              {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminPaymentAPI } from '~/composables/useApi'

const providers = ref([])
const showEditModal = ref(false)
const showAddProvider = ref(false)
const editingProvider = ref({})
const selectedCountries = ref([])
const saving = ref(false)
const loading = ref(false)

const availableCountries = [
  { code: 'CI', name: "Côte d'Ivoire", currency: 'XOF', flag: '🇨🇮' },
  { code: 'SN', name: 'Sénégal', currency: 'XOF', flag: '🇸🇳' },
  { code: 'CM', name: 'Cameroun', currency: 'XAF', flag: '🇨🇲' },
  { code: 'NG', name: 'Nigeria', currency: 'NGN', flag: '🇳🇬' },
  { code: 'GH', name: 'Ghana', currency: 'GHS', flag: '🇬🇭' },
  { code: 'BF', name: 'Burkina Faso', currency: 'XOF', flag: '🇧🇫' },
  { code: 'ML', name: 'Mali', currency: 'XOF', flag: '🇲🇱' },
  { code: 'KE', name: 'Kenya', currency: 'KES', flag: '🇰🇪' },
  { code: 'ZA', name: 'Afrique du Sud', currency: 'ZAR', flag: '🇿🇦' },
]

const activeProviders = computed(() => providers.value.filter(p => p.is_active).length)
const demoProviders = computed(() => providers.value.filter(p => p.is_demo_mode).length)
const totalCountries = computed(() => {
  const countries = new Set()
  providers.value.forEach(p => (p.countries || []).forEach(c => countries.add(c.country_code)))
  return countries.size
})

const getProviderIcon = (name) => {
  const icons = {
    demo: '🧪', flutterwave: '🦋', cinetpay: '💳', paystack: '⚡',
    orange_money: '🟠', mtn_momo: '🟡', wave: '🌊', stripe: '💜'
  }
  return icons[name] || '💰'
}

const getProviderBg = (name) => {
  const bgs = {
    demo: 'bg-slate-600', flutterwave: 'bg-orange-500/20', cinetpay: 'bg-blue-500/20', 
    paystack: 'bg-cyan-500/20', orange_money: 'bg-orange-500/20', mtn_momo: 'bg-yellow-500/20', 
    wave: 'bg-blue-500/20', stripe: 'bg-purple-500/20'
  }
  return bgs[name] || 'bg-slate-600'
}

const getTypeBadge = (type) => {
  const badges = {
    all: 'bg-indigo-500/20 text-indigo-400',
    mobile_money: 'bg-emerald-500/20 text-emerald-400',
    card: 'bg-purple-500/20 text-purple-400',
    bank_transfer: 'bg-blue-500/20 text-blue-400'
  }
  return badges[type] || 'bg-slate-500/20 text-slate-400'
}

const getTypeLabel = (type) => {
  const labels = { all: 'Tous', mobile_money: 'Mobile Money', card: 'Carte', bank_transfer: 'Virement' }
  return labels[type] || type
}

const isCountrySelected = (code) => selectedCountries.value.includes(code)

const loadProviders = async () => {
  loading.value = true
  try {
    const response = await adminPaymentAPI.getProviders()
    if (response.data?.providers) {
      providers.value = response.data.providers
    }
  } catch (e) {
    console.error('Error loading providers:', e)
    // Fallback to mock data if API fails
    providers.value = [
      { id: '1', name: 'demo', display_name: 'Mode Démo', provider_type: 'all', is_active: true, is_demo_mode: true, countries: [
        { country_code: 'CI' }, { country_code: 'SN' }, { country_code: 'NG' }, { country_code: 'GH' }, { country_code: 'CM' }
      ]},
      { id: '2', name: 'flutterwave', display_name: 'Flutterwave', provider_type: 'all', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'NG' }, { country_code: 'GH' }, { country_code: 'KE' }, { country_code: 'CI' }
      ]},
      { id: '3', name: 'cinetpay', display_name: 'CinetPay', provider_type: 'mobile_money', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'CI' }, { country_code: 'SN' }, { country_code: 'CM' }, { country_code: 'BF' }
      ]},
      { id: '4', name: 'paystack', display_name: 'Paystack', provider_type: 'all', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'NG' }, { country_code: 'GH' }
      ]},
      { id: '5', name: 'orange_money', display_name: 'Orange Money', provider_type: 'mobile_money', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'CI' }, { country_code: 'SN' }
      ]},
      { id: '6', name: 'mtn_momo', display_name: 'MTN MoMo', provider_type: 'mobile_money', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'CI' }, { country_code: 'CM' }, { country_code: 'GH' }
      ]},
      { id: '7', name: 'wave', display_name: 'Wave', provider_type: 'mobile_money', is_active: false, is_demo_mode: true, countries: [
        { country_code: 'SN' }, { country_code: 'CI' }
      ]},
      { id: '8', name: 'stripe', display_name: 'Stripe', provider_type: 'card', is_active: false, is_demo_mode: true, countries: [] },
    ]
  } finally {
    loading.value = false
  }
}

const editProvider = (provider) => {
  editingProvider.value = { ...provider }
  selectedCountries.value = (provider.countries || []).map(c => c.country_code)
  showEditModal.value = true
}

const toggleActive = async (provider) => {
  const newStatus = !provider.is_active
  try {
    await adminPaymentAPI.toggleStatus(provider.id, newStatus)
    provider.is_active = newStatus
  } catch (e) {
    console.error('Toggle status error:', e)
    provider.is_active = newStatus // Optimistic update even if API fails
  }
}

const testProvider = async (provider) => {
  try {
    const response = await adminPaymentAPI.testConnection(provider.id)
    if (response.data?.success) {
      alert(`✅ Connexion réussie pour ${provider.display_name}\nMode: ${response.data.mode}`)
    } else {
      alert(`❌ Échec de connexion pour ${provider.display_name}`)
    }
  } catch (e) {
    alert(`❌ Erreur de test: ${e.message}`)
  }
}

const deleteProvider = async (provider) => {
  if (provider.name === 'demo') return
  if (confirm(`Supprimer ${provider.display_name} ?`)) {
    try {
      await adminPaymentAPI.deleteProvider(provider.id)
      providers.value = providers.value.filter(p => p.id !== provider.id)
    } catch (e) {
      console.error('Delete error:', e)
      alert('Erreur lors de la suppression')
    }
  }
}

const closeEditModal = () => {
  showEditModal.value = false
  editingProvider.value = {}
  selectedCountries.value = []
}

const saveProvider = async () => {
  saving.value = true
  try {
    if (editingProvider.value.id) {
      await adminPaymentAPI.updateProvider(editingProvider.value.id, editingProvider.value)
    } else {
      await adminPaymentAPI.createProvider(editingProvider.value)
    }
    closeEditModal()
    await loadProviders()
  } catch (e) {
    console.error('Save error:', e)
    alert('Erreur lors de la sauvegarde')
  } finally {
    saving.value = false
  }
}

onMounted(loadProviders)

definePageMeta({
  layout: false,
  middleware: 'admin-auth'
})
</script>
