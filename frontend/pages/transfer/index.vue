<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Envoyer de l'argent 💸</h1>
        <p class="text-gray-900/60">Transferts internationaux, Mobile Money, et virements bancaires</p>
      </div>

      <!-- Transfer Type Selector -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <button 
          v-for="type in transferTypes" :key="type.id"
          @click="selectedType = type.id"
          class="glass-card p-6 text-left transition-all"
          :class="selectedType === type.id ? 'ring-2 ring-indigo-500' : ''">
          <span class="text-3xl mb-3 block">{{ type.icon }}</span>
          <h3 class="font-semibold text-gray-900 mb-1">{{ type.name }}</h3>
          <p class="text-sm text-gray-900/50">{{ type.description }}</p>
        </button>
      </div>

      <!-- Transfer Form -->
      <div class="glass-card p-8">
        <form @submit.prevent="submitTransfer" class="space-y-6">
          <!-- Amount -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Montant à envoyer</label>
            <div class="flex gap-4">
              <input
                v-model.number="form.amount"
                type="number"
                min="1"
                step="0.01"
                required
                class="input-premium flex-1 text-2xl"
                placeholder="0.00"
              />
              <select v-model="form.currency" class="input-premium w-32">
                <option value="USD">USD</option>
                <option value="EUR">EUR</option>
                <option value="GBP">GBP</option>
                <option value="XOF">XOF</option>
                <option value="XAF">XAF</option>
              </select>
            </div>
          </div>

          <!-- Recipient -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">
              {{ getRecipientLabel }}
            </label>
            <input
              v-model="form.recipient"
              type="text"
              required
              class="input-premium"
              :placeholder="getRecipientPlaceholder"
            />
          </div>

          <!-- Country (for Mobile Money) -->
          <div v-if="selectedType === 'mobile'">
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Pays du destinataire</label>
            <select v-model="form.country" class="input-premium">
              <option value="CI">🇨🇮 Côte d'Ivoire</option>
              <option value="SN">🇸🇳 Sénégal</option>
              <option value="CM">🇨🇲 Cameroun</option>
              <option value="GH">🇬🇭 Ghana</option>
              <option value="KE">🇰🇪 Kenya</option>
              <option value="NG">🇳🇬 Nigeria</option>
            </select>
          </div>

          <!-- Bank (for Wire) -->
          <div v-if="selectedType === 'wire'">
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Banque du destinataire</label>
            <input
              v-model="form.bankName"
              type="text"
              class="input-premium"
              placeholder="Nom de la banque"
            />
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-900/80 mb-2">Motif du transfert</label>
            <input
              v-model="form.description"
              type="text"
              class="input-premium"
              placeholder="ex: Soutien familial"
            />
          </div>

          <!-- Summary -->
          <div class="p-6 rounded-xl bg-white/5">
            <h4 class="font-semibold text-gray-900 mb-4">Résumé</h4>
            <div class="space-y-3 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-900/60">Montant envoyé</span>
                <span class="text-gray-900 font-medium">{{ formatMoney(form.amount || 0, form.currency) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-900/60">Frais de transfert</span>
                <span class="text-gray-900 font-medium">{{ formatMoney(fees, form.currency) }}</span>
              </div>
              <div class="border-t border-white/10 pt-3 flex justify-between">
                <span class="text-gray-900 font-semibold">Total débité</span>
                <span class="text-gray-900 font-bold text-lg">{{ formatMoney((form.amount || 0) + fees, form.currency) }}</span>
              </div>
              <div class="flex justify-between text-emerald-400">
                <span>Le destinataire reçoit</span>
                <span class="font-semibold">{{ formatMoney(receivedAmount, localCurrency) }}</span>
              </div>
            </div>
          </div>

          <!-- Submit -->
          <button type="submit" :disabled="loading" class="btn-premium w-full py-4 text-lg">
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <div class="loading-spinner w-5 h-5"></div>
              Envoi en cours...
            </span>
            <span v-else>Envoyer {{ formatMoney(form.amount || 0, form.currency) }}</span>
          </button>
        </form>
      </div>

      <!-- Recent Transfers -->
      <div class="glass-card p-6 mt-8">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">Transferts récents</h3>
        <div class="space-y-4">
          <div v-for="tx in recentTransfers" :key="tx.id" 
              class="flex items-center justify-between p-4 rounded-xl bg-white/5">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-indigo-500/20 flex items-center justify-center">
                <span class="text-lg">💸</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">{{ tx.recipient }}</p>
                <p class="text-sm text-gray-900/50">{{ tx.type }} • {{ formatDate(tx.date) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span class="badge" :class="getStatusClass(tx.status)">{{ tx.status }}</span>
              <p class="font-semibold text-gray-900">{{ formatMoney(tx.amount, tx.currency) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed } from 'vue'

const transferTypes = [
  { id: 'mobile', name: 'Mobile Money', icon: '📱', description: 'Orange, MTN, Wave, M-Pesa' },
  { id: 'wire', name: 'Virement bancaire', icon: '🏦', description: 'Compte bancaire' },
  { id: 'crypto', name: 'Crypto', icon: '₿', description: 'BTC, ETH, USDT' }
]

const selectedType = ref('mobile')

const form = ref({
  amount: null,
  currency: 'USD',
  recipient: '',
  country: 'CI',
  bankName: '',
  description: ''
})

const loading = ref(false)

const fees = computed(() => {
  const amount = form.value.amount || 0
  if (selectedType.value === 'mobile') return Math.max(0.99, amount * 0.01)
  if (selectedType.value === 'wire') return Math.max(4.99, amount * 0.02)
  return amount * 0.005
})

const localCurrency = computed(() => {
  if (selectedType.value === 'mobile') {
    const currencies = { CI: 'XOF', SN: 'XOF', CM: 'XAF', GH: 'GHS', KE: 'KES', NG: 'NGN' }
    return currencies[form.value.country] || 'XOF'
  }
  return form.value.currency
})

const receivedAmount = computed(() => {
  const amount = form.value.amount || 0
  const rates = { USD: 600, EUR: 655, GBP: 760 }
  if (localCurrency.value === 'XOF' || localCurrency.value === 'XAF') {
    return amount * (rates[form.value.currency] || 600)
  }
  return amount
})

const getRecipientLabel = computed(() => {
  if (selectedType.value === 'mobile') return 'Numéro de téléphone'
  if (selectedType.value === 'wire') return 'IBAN / Numéro de compte'
  return 'Adresse du portefeuille'
})

const getRecipientPlaceholder = computed(() => {
  if (selectedType.value === 'mobile') return '+225 07 12 34 56 78'
  if (selectedType.value === 'wire') return 'FR76 1234 5678 9012 3456'
  return '0x742d35Cc6634C0532925a3b844Bc9e7595f2bD'
})

// Recent transfers - will be loaded from API
const recentTransfers = ref([])

const formatMoney = (amount, currency = 'USD') => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

const formatDate = (date) => new Date(date).toLocaleDateString('fr-FR')

const getStatusClass = (status) => {
  return status === 'completed' ? 'badge-success' : status === 'pending' ? 'badge-warning' : 'badge-danger'
}

const submitTransfer = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 2000))
  loading.value = false
  alert('Transfert envoyé avec succès!')
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>