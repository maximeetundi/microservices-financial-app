<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeModal"></div>

      <!-- Modal -->
      <div class="relative bg-white dark:bg-slate-800 rounded-2xl shadow-2xl max-w-md w-full max-h-[90vh] overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-slate-700">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">
            {{ currentStep === 'success' ? '✅ Dépôt Réussi' : 'Recharger Compte' }}
          </h2>
          <button @click="closeModal" class="p-2 hover:bg-gray-100 dark:hover:bg-slate-700 rounded-full transition-colors">
            <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div ref="contentRef" class="p-6 overflow-y-auto max-h-[calc(90vh-180px)]">

          <!-- Step 1: Select Provider -->
          <div v-if="currentStep === 'provider'" class="space-y-4">
            <!-- Amount Input -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Montant à déposer
              </label>
              <div class="relative">
                <input
                  v-model.number="amount"
                  type="number"
                  min="100"
                  :max="maxAmount"
                  class="w-full px-4 py-3 pr-16 rounded-xl border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  placeholder="5000"
                />
                <span class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400 font-medium">
                  {{ currency }}
                </span>
              </div>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Min: {{ formatAmount(minAmount) }} {{ currency }} - Max: {{ formatAmount(maxAmount) }} {{ currency }}
              </p>
            </div>

            <!-- Provider Categories -->
            <div class="space-y-4">
              <!-- Mobile Money -->
              <div v-if="mobileMoneyProviders.length > 0">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3 flex items-center gap-2">
                  <span>📱</span> MOBILE MONEY
                </h3>
                <div class="grid gap-3">
                  <button
                    v-for="provider in mobileMoneyProviders"
                    :key="provider.name"
                    @click="selectProvider(provider)"
                    :disabled="!provider.deposit_enabled"
                    class="flex items-center gap-4 p-4 rounded-xl border-2 transition-all"
                    :class="selectedProvider?.name === provider.name
                      ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
                      : 'border-gray-200 dark:border-slate-600 hover:border-indigo-300 dark:hover:border-indigo-600'"
                  >
                    <img
                      :src="getProviderLogo(provider)"
                      :alt="provider.display_name"
                      class="w-12 h-12 rounded-lg object-contain bg-white p-1"
                    />
                    <div class="flex-1 text-left">
                      <p class="font-semibold text-gray-900 dark:text-white">{{ provider.display_name }}</p>
                      <p class="text-sm text-gray-500 dark:text-gray-400">
                        {{ provider.is_demo_mode ? 'Mode Démo' : 'Paiement instantané' }}
                      </p>
                    </div>
                    <div v-if="selectedProvider?.name === provider.name" class="text-indigo-500">
                      <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                      </svg>
                    </div>
                  </button>
                </div>
              </div>

              <!-- Cards & International -->
              <div v-if="cardProviders.length > 0">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3 flex items-center gap-2">
                  <span>💳</span> CARTE BANCAIRE
                </h3>
                <div class="grid gap-3">
                  <button
                    v-for="provider in cardProviders"
                    :key="provider.name"
                    @click="selectProvider(provider)"
                    :disabled="!provider.deposit_enabled"
                    class="flex items-center gap-4 p-4 rounded-xl border-2 transition-all"
                    :class="selectedProvider?.name === provider.name
                      ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
                      : 'border-gray-200 dark:border-slate-600 hover:border-indigo-300 dark:hover:border-indigo-600'"
                  >
                    <img
                      :src="getProviderLogo(provider)"
                      :alt="provider.display_name"
                      class="w-12 h-12 rounded-lg object-contain bg-white p-1"
                    />
                    <div class="flex-1 text-left">
                      <p class="font-semibold text-gray-900 dark:text-white">{{ provider.display_name }}</p>
                      <p class="text-sm text-gray-500 dark:text-gray-400">Paiement sécurisé</p>
                    </div>
                  </button>
                </div>
              </div>

              <!-- Crypto Ramp -->
              <div v-if="cryptoProviders.length > 0">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3 flex items-center gap-2">
                  <span>🪙</span> CRYPTO
                </h3>
                <div class="grid gap-3">
                  <button
                    v-for="provider in cryptoProviders"
                    :key="provider.name"
                    @click="selectProvider(provider)"
                    class="flex items-center gap-4 p-4 rounded-xl border-2 transition-all"
                    :class="selectedProvider?.name === provider.name
                      ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
                      : 'border-gray-200 dark:border-slate-600 hover:border-indigo-300'"
                  >
                    <img
                      :src="getProviderLogo(provider)"
                      :alt="provider.display_name"
                      class="w-12 h-12 rounded-lg object-contain bg-white p-1"
                    />
                    <div class="flex-1 text-left">
                      <p class="font-semibold text-gray-900 dark:text-white">{{ provider.display_name }}</p>
                      <p class="text-sm text-gray-500 dark:text-gray-400">Acheter crypto → Fiat</p>
                    </div>
                  </button>
                </div>
              </div>

              <!-- No Providers Message -->
              <div v-if="providers.length === 0 && !loadingProviders" class="text-center py-8">
                <p class="text-gray-500 dark:text-gray-400">
                  Aucun moyen de paiement disponible pour votre pays.
                </p>
              </div>

              <!-- Loading -->
              <div v-if="loadingProviders" class="flex justify-center py-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
              </div>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl">
              <p class="text-red-600 dark:text-red-400 text-sm">{{ error }}</p>
            </div>
          </div>

          <!-- Step: PayPal Checkout (JS SDK Buttons) -->
          <div v-if="currentStep === 'paypal'" class="space-y-4">
            <div class="text-center mb-2">
              <img
                :src="getProviderLogo(selectedProvider)"
                :alt="selectedProvider?.display_name"
                class="w-16 h-16 mx-auto rounded-xl bg-white p-2 mb-3"
              />
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                Payer avec PayPal
              </h3>
              <p class="text-gray-500 dark:text-gray-400">
                Montant: {{ formatAmount(amount) }} {{ currency }}
              </p>
            </div>

            <div
              v-if="paypalError"
              class="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl"
            >
              <p class="text-red-600 dark:text-red-400 text-sm">{{ paypalError }}</p>
            </div>

            <div v-if="paypalLoading" class="flex justify-center py-4">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
            </div>

            <div id="paypal-buttons" class="min-h-[40px]"></div>

            <button
              @click="currentStep = 'provider'"
              class="text-indigo-500 hover:text-indigo-600 text-sm font-medium"
            >
              ← Changer de méthode
            </button>
          </div>

          <!-- Step: Choose Mode (Sandbox / Live) -->
          <div v-if="currentStep === 'mode'" class="space-y-4">
            <div class="text-center mb-6">
              <img
                :src="getProviderLogo(selectedProvider)"
                :alt="selectedProvider?.display_name"
                class="w-16 h-16 mx-auto rounded-xl bg-white p-2 mb-3"
              />
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                Choisir le mode
              </h3>
              <p class="text-gray-500 dark:text-gray-400">
                {{ selectedProvider?.display_name }}
              </p>
            </div>

            <div class="grid gap-3">
              <button
                @click="selectMode(false)"
                class="flex items-center gap-4 p-4 rounded-xl border-2 transition-all"
                :class="selectedIsTestMode === false
                  ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
                  : 'border-gray-200 dark:border-slate-600 hover:border-indigo-300 dark:hover:border-indigo-600'"
              >
                <span class="text-2xl">🟦</span>
                <div class="flex-1 text-left">
                  <p class="font-semibold text-gray-900 dark:text-white">Live</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Production</p>
                </div>
              </button>

              <button
                @click="selectMode(true)"
                class="flex items-center gap-4 p-4 rounded-xl border-2 transition-all"
                :class="selectedIsTestMode === true
                  ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
                  : 'border-gray-200 dark:border-slate-600 hover:border-indigo-300 dark:hover:border-indigo-600'"
              >
                <span class="text-2xl">🧪</span>
                <div class="flex-1 text-left">
                  <p class="font-semibold text-gray-900 dark:text-white">Sandbox</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Test</p>
                </div>
              </button>
            </div>

            <button
              @click="currentStep = 'provider'"
              class="text-indigo-500 hover:text-indigo-600 text-sm font-medium"
            >
              ← Changer de méthode
            </button>
          </div>

          <!-- Step 2: Phone Number for Mobile Money -->
          <div v-if="currentStep === 'phone'" class="space-y-4">
            <div class="text-center mb-6">
              <img
                :src="getProviderLogo(selectedProvider)"
                :alt="selectedProvider?.display_name"
                class="w-16 h-16 mx-auto rounded-xl bg-white p-2 mb-3"
              />
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ selectedProvider?.display_name }}
              </h3>
              <p class="text-gray-500 dark:text-gray-400">
                Montant: {{ formatAmount(amount) }} {{ currency }}
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Numéro de recharge
              </label>

              <div v-if="depositNumbersLoading" class="flex justify-center py-4">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-500"></div>
              </div>

              <div v-else-if="depositNumbers.length > 0" class="space-y-3">
                <select
                  v-model="selectedDepositNumberId"
                  @change="onSelectDepositNumber"
                  class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                >
                  <option v-for="n in depositNumbers" :key="n.id" :value="n.id">
                    {{ n.label ? `${n.label} - ${n.phone}` : n.phone }}
                  </option>
                </select>

                <div class="p-4 rounded-xl border border-gray-200 dark:border-slate-600 bg-gray-50 dark:bg-slate-900/30">
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">Ajouter un nouveau numéro</p>
                  <div class="space-y-2">
                    <input
                      v-model="newDepositNumberPhone"
                      type="tel"
                      class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                      placeholder="+225 07 XX XX XX XX"
                    />
                    <input
                      v-model="newDepositNumberLabel"
                      type="text"
                      class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                      placeholder="Label (optionnel)"
                    />
                    <button
                      @click="addDepositNumber"
                      class="w-full px-4 py-3 bg-indigo-500 text-white rounded-xl hover:bg-indigo-600 transition-colors"
                    >
                      Ajouter ce numéro
                    </button>
                  </div>
                </div>
              </div>

              <div v-else class="space-y-2">
                <input
                  v-model="phoneNumber"
                  type="tel"
                  class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500"
                  placeholder="+225 07 XX XX XX XX"
                />
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  Astuce: ajoute ce numéro dans tes paramètres pour le retrouver plus vite.
                </p>
              </div>
            </div>

            <button
              @click="currentStep = 'provider'"
              class="text-indigo-500 hover:text-indigo-600 text-sm font-medium"
            >
              ← Changer de méthode
            </button>
          </div>

          <!-- Step 3: Processing -->
          <div v-if="currentStep === 'processing'" class="text-center py-8">
            <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-indigo-500 mx-auto mb-4"></div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              Initialisation du paiement...
            </h3>
            <p class="text-gray-500 dark:text-gray-400">
              Veuillez patienter
            </p>
          </div>

          <!-- Step 4: Redirect Info -->
          <div v-if="currentStep === 'redirect'" class="text-center py-8">
            <div class="w-16 h-16 bg-indigo-100 dark:bg-indigo-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              Continuer vers {{ selectedProvider?.display_name }}
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-4">
              Pour finaliser votre dépôt, le paiement doit être validé sur la page sécurisée du prestataire.
            </p>
            <p v-if="redirectReason" class="text-sm text-gray-500 dark:text-gray-400 mb-4">
              {{ redirectReason }}
            </p>

            <div class="max-w-md mx-auto text-left bg-gray-50 dark:bg-slate-800 rounded-xl p-4 border border-gray-100 dark:border-gray-700">
              <div class="flex justify-between text-sm text-gray-600 dark:text-gray-300">
                <span>Montant</span>
                <span class="font-semibold">{{ formatAmount(amount) }} {{ currency }}</span>
              </div>
              <div class="flex justify-between text-sm text-gray-600 dark:text-gray-300 mt-2">
                <span>Prestataire</span>
                <span class="font-semibold">{{ selectedProvider?.display_name || selectedProvider?.name }}</span>
              </div>
              <div class="flex justify-between text-sm text-gray-600 dark:text-gray-300 mt-2">
                <span>Transaction</span>
                <span class="font-mono text-xs">{{ transactionId }}</span>
              </div>
            </div>
            <p class="text-sm text-gray-400 dark:text-gray-500">
              Un nouvel onglet sera ouvert.
            </p>
            <button
              v-if="paymentUrl"
              @click="continueToPayment"
              class="inline-block mt-4 px-6 py-3 bg-indigo-500 text-white rounded-xl hover:bg-indigo-600 transition-colors"
            >
              Valider et continuer
            </button>

            <button
              @click="cancelRedirect"
              class="block mx-auto mt-3 text-indigo-500 hover:text-indigo-600 text-sm font-medium"
            >
              ← Annuler / Retour
            </button>
            <a
              v-if="paymentUrl"
              :href="paymentUrl"
              target="_blank"
              class="inline-block mt-3 px-6 py-3 bg-white dark:bg-slate-900 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-200 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors"
            >
              Ouvrir manuellement
            </a>
          </div>

          <!-- Step 5: Pending -->
          <div v-if="currentStep === 'pending'" class="text-center py-8">
            <div class="w-16 h-16 bg-yellow-100 dark:bg-yellow-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              Paiement en attente
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-4">
              Veuillez compléter le paiement sur la page de {{ selectedProvider?.display_name }}.
              <br/>Votre compte sera crédité automatiquement.
            </p>
            <p class="text-sm text-gray-400 dark:text-gray-500">
              Transaction: {{ transactionId }}
            </p>
            <button
              @click="checkStatus"
              class="mt-4 px-4 py-2 text-indigo-500 hover:text-indigo-600 font-medium"
            >
              Vérifier le statut
            </button>
          </div>

          <!-- Step 6: Success -->
          <div v-if="currentStep === 'success'" class="text-center py-8">
            <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4 animate-bounce-once">
              <svg class="w-10 h-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">
              Dépôt Réussi! 🎉
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-2">
              {{ formatAmount(amount) }} {{ currency }} ont été ajoutés à votre compte.
            </p>
            <p v-if="newBalance" class="text-lg font-semibold text-indigo-500">
              Nouveau solde: {{ formatAmount(newBalance) }} {{ currency }}
            </p>
          </div>

          <!-- Step: Failed -->
          <div v-if="currentStep === 'failed'" class="text-center py-8">
            <div class="w-16 h-16 bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              Échec du paiement
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-4">
              {{ error || 'Le paiement a échoué. Veuillez réessayer.' }}
            </p>
            <button
              @click="resetForm"
              class="px-6 py-2 bg-indigo-500 text-white rounded-xl hover:bg-indigo-600 transition-colors"
            >
              Réessayer
            </button>
          </div>
        </div>

        <!-- Footer Actions -->
        <div v-if="['provider', 'mode', 'phone'].includes(currentStep)" class="p-6 border-t border-gray-200 dark:border-slate-700 bg-gray-50 dark:bg-slate-900/50">
          <button
            @click="initiateDeposit"
            :disabled="!canProceed || loading"
            class="w-full py-4 px-6 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg shadow-indigo-500/30 hover:shadow-xl hover:shadow-indigo-500/40 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="loading" class="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></span>
            <span v-else>
              Déposer {{ formatAmount(amount) }} {{ currency }}
            </span>
          </button>
        </div>

        <div v-if="currentStep === 'success'" class="p-6 border-t border-gray-200 dark:border-slate-700">
          <button
            @click="closeModal"
            class="w-full py-3 px-6 bg-indigo-500 text-white font-semibold rounded-xl hover:bg-indigo-600 transition-colors"
          >
            Fermer
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Flutterwave SDK -->
  <div id="flutterwave-container"></div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useWalletStore } from '~/stores/wallet'
import { depositNumbersAPI } from '~/composables/useApi'

// Props
const props = defineProps<{
  isOpen: boolean
  walletId?: string
  currency?: string
}>()

const emit = defineEmits(['close', 'success'])

// Stores
const authStore = useAuthStore()
const walletStore = useWalletStore()

// State
const currentStep = ref<'provider' | 'mode' | 'phone' | 'processing' | 'redirect' | 'paypal' | 'pending' | 'success' | 'failed'>('provider')
const amount = ref(5000)
const phoneNumber = ref('')
const depositNumbers = ref<any[]>([])
const selectedDepositNumberId = ref<string>('')
const depositNumbersLoading = ref(false)
const newDepositNumberPhone = ref('')
const newDepositNumberLabel = ref('')
const selectedProvider = ref<any>(null)
const selectedIsTestMode = ref<boolean | null>(null)
const providers = ref<any[]>([])
const loadingProviders = ref(false)
const loading = ref(false)
const error = ref('')
const transactionId = ref('')
const paymentUrl = ref('')
const newBalance = ref<number | null>(null)
const sdkConfig = ref<any>(null)
const paypalLoading = ref(false)
const paypalError = ref('')
const isMobile = ref(false)
const contentRef = ref<HTMLElement | null>(null)

// Config
const minAmount = 100
const maxAmount = 5000000
const currency = computed(() => props.currency || 'XOF')
const userCountry = computed(() => authStore.user?.country || 'CI')

// Computed
const mobileMoneyProviders = computed(() =>
  providers.value.filter(p =>
    ['mobile_money', 'all'].includes(p.provider_type) &&
    !['stripe', 'paypal'].includes(p.name)
  )
)

const cardProviders = computed(() =>
  providers.value.filter(p =>
    ['card', 'wallet', 'international'].includes(p.provider_type) ||
    ['stripe', 'paypal'].includes(p.name)
  )
)

const cryptoProviders = computed(() =>
  providers.value.filter(p => p.provider_type === 'crypto_ramp')
)

const canProceed = computed(() => {
  if (!selectedProvider.value) return false
  if (selectedIsTestMode.value === null) return false
  if (amount.value < minAmount || amount.value > maxAmount) return false
  if (currentStep.value === 'phone' && !selectedDepositNumberId.value && !phoneNumber.value) return false
  return true
})

// API Base URL
const API_URL = (() => {
  try {
    const config = useRuntimeConfig()
    return config.public.apiBaseUrl || 'https://api.app.tech-afm.com'
  } catch {
    return 'https://api.app.tech-afm.com'
  }
})()

// Methods
const fetchProviders = async () => {
  loadingProviders.value = true
  error.value = ''

  try {
    await authStore.initializeAuth()
    const token = authStore.accessToken || (typeof window !== 'undefined' ? localStorage.getItem('accessToken') : '')
    if (!token) {
      throw new Error('Veuillez vous reconnecter pour voir les moyens de paiement.')
    }
    const response = await fetch(`${API_URL}/transfer-service/api/v1/aggregators/deposit?country=${userCountry.value}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) throw new Error('Erreur lors du chargement des méthodes de paiement')

    const data = await response.json()
    providers.value = data.aggregators || []

    // Auto-select demo if only option
    if (providers.value.length === 1) {
      selectedProvider.value = providers.value[0]
    }
  } catch (err: any) {
    error.value = err.message
    console.error('Failed to fetch providers:', err)
  } finally {
    loadingProviders.value = false
  }
}

const fetchDepositNumbers = async () => {
  if (!authStore.accessToken) return
  depositNumbersLoading.value = true
  try {
    const res = await depositNumbersAPI.list()
    depositNumbers.value = res.data?.numbers || []

    // Auto-seed with main phone if list is empty
    if (depositNumbers.value.length === 0 && authStore.user?.phone) {
      await depositNumbersAPI.create({
        country: userCountry.value,
        phone: authStore.user.phone,
        label: 'Numéro principal',
        is_default: true,
      })
      const seeded = await depositNumbersAPI.list()
      depositNumbers.value = seeded.data?.numbers || []
    }

    const def = depositNumbers.value.find((n: any) => n.is_default)
    if (def?.id) {
      selectedDepositNumberId.value = def.id
      phoneNumber.value = def.phone
    } else if (depositNumbers.value.length > 0 && depositNumbers.value[0]?.id) {
      selectedDepositNumberId.value = depositNumbers.value[0].id
      phoneNumber.value = depositNumbers.value[0].phone
    } else {
      selectedDepositNumberId.value = ''
    }
  } catch (_e) {
    depositNumbers.value = []
    selectedDepositNumberId.value = ''
  } finally {
    depositNumbersLoading.value = false
  }
}

const onSelectDepositNumber = () => {
  const n = depositNumbers.value.find((x: any) => x.id === selectedDepositNumberId.value)
  phoneNumber.value = n?.phone || ''
}

const addDepositNumber = async () => {
  if (!newDepositNumberPhone.value) return
  try {
    const res = await depositNumbersAPI.create({
      country: userCountry.value,
      phone: newDepositNumberPhone.value,
      label: newDepositNumberLabel.value,
    })
    const created = res.data?.number
    await fetchDepositNumbers()
    if (created?.id) {
      selectedDepositNumberId.value = created.id
      onSelectDepositNumber()
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || "Erreur lors de l'ajout du numéro"
  } finally {
    newDepositNumberPhone.value = ''
    newDepositNumberLabel.value = ''
  }
}

const selectProvider = (provider: any) => {
  selectedProvider.value = provider
  selectedIsTestMode.value = null
  currentStep.value = 'mode'
  error.value = ''
}

const selectMode = (isTest: boolean) => {
  selectedIsTestMode.value = isTest
  error.value = ''
}

const getProviderLogo = (provider: any) => {
  if (!provider) return '/icons/aggregators/default.svg'
  if (provider.logo_url) return provider.logo_url
  return `/icons/aggregators/${provider.name}.svg`
}

const formatAmount = (value: number) => {
  return new Intl.NumberFormat('fr-FR').format(value)
}

const isMobileMoneySelected = computed(() => {
  if (!selectedProvider.value) return false
  return mobileMoneyProviders.value.some(p => p.name === selectedProvider.value.name)
})

const initiateDeposit = async () => {
  if (!selectedProvider.value || !amount.value) return

  if (selectedIsTestMode.value === null) {
    currentStep.value = 'mode'
    return
  }

  if (currentStep.value === 'mode') {
    if (isMobileMoneySelected.value) {
      currentStep.value = 'phone'
      return
    }
  }

  // For mobile money, check phone number
  if (isMobileMoneySelected.value) {
    if (currentStep.value === 'provider') {
      currentStep.value = 'phone'
      return
    }
  }

  loading.value = true
  error.value = ''
  currentStep.value = 'processing'

  try {
    const returnUrl = `${window.location.origin}/wallet?deposit_callback=true`
    const cancelUrl = `${window.location.origin}/wallet?deposit_cancelled=true`

    const response = await fetch(`${API_URL}/transfer-service/api/v1/deposits/initiate`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.accessToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        user_id: authStore.user?.id,
        wallet_id: props.walletId,
        amount: amount.value,
        currency: currency.value,
        provider: selectedProvider.value.name,
        country: userCountry.value,
        is_test_mode: selectedIsTestMode.value,
        email: authStore.user?.email,
        phone: isMobileMoneySelected.value ? phoneNumber.value : undefined,
        deposit_number_id: isMobileMoneySelected.value ? selectedDepositNumberId.value : undefined,
        return_url: returnUrl,
        cancel_url: cancelUrl
      })
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Erreur lors de l\'initiation du paiement')
    }

    transactionId.value = data.transaction_id
    paymentUrl.value = data.payment_url
    sdkConfig.value = data.sdk_config

    // Handle based on status
    if (data.status === 'completed' || data.status === 'instant_success') {
      // Demo or instant payment
      newBalance.value = data.new_balance
      currentStep.value = 'success'
      emit('success', { transactionId: transactionId.value, amount: amount.value })
      walletStore.fetchWallets()
    } else if (data.payment_url) {
      // Redirect to payment page or use SDK
      handlePaymentFlow(data)
    } else {
      currentStep.value = 'pending'
    }
  } catch (err: any) {
    error.value = err.message
    currentStep.value = 'failed'
    console.error('Deposit failed:', err)
  } finally {
    loading.value = false
  }
}

const handlePaymentFlow = (data: any) => {
  const provider = String(selectedProvider.value?.name || '').toLowerCase().trim()

  // Try SDK first, fallback to redirect
  switch (provider) {
    case 'flutterwave':
      if (!isMobile.value && window.FlutterwaveCheckout && sdkConfig.value?.public_key) {
        openFlutterwaveModal(data)
      } else {
        redirectToPayment(data.payment_url, 'Le paiement intégré n\'est pas disponible sur cet appareil. Vous pouvez continuer via la page du prestataire.')
      }
      break

    case 'paystack':
      if (!isMobile.value && window.PaystackPop && sdkConfig.value?.public_key) {
        openPaystackModal(data)
      } else {
        redirectToPayment(data.payment_url, 'Le paiement intégré n\'est pas disponible sur cet appareil. Vous pouvez continuer via la page du prestataire.')
      }
      break

    case 'stripe':
      // Stripe Checkout redirect
      redirectToPayment(data.payment_url, 'Stripe utilise une page de paiement hébergée pour finaliser la transaction.')
      break

    case 'paypal':
      // PayPal JS SDK Buttons (fallback to redirect if not possible)
      if (sdkConfig.value?.public_key) {
        openPayPalButtons(data)
      } else {
        redirectToPayment(data.payment_url, 'Le paiement intégré PayPal n\'est pas disponible sur cet appareil. Vous pouvez continuer via la page PayPal.')
      }
      break

    default:
      // All other providers: redirect
      redirectToPayment(data.payment_url, 'Ce moyen de paiement nécessite une validation externe sur la page du prestataire.')
  }
}

const redirectReason = ref('')

const redirectToPayment = (url: string, reason = '') => {
  currentStep.value = 'redirect'
  paymentUrl.value = url
  redirectReason.value = reason
}

const continueToPayment = () => {
  if (!paymentUrl.value) return
  window.open(paymentUrl.value, '_blank')
  currentStep.value = 'pending'
}

const cancelRedirect = () => {
  paymentUrl.value = ''
  redirectReason.value = ''
  transactionId.value = ''
  currentStep.value = 'provider'
}

const loadPayPalSdk = (clientId: string) => {
  return new Promise<void>((resolve, reject) => {
    if (window.paypal) {
      resolve()
      return
    }

    const existing = document.getElementById('paypal-sdk') as HTMLScriptElement | null
    if (existing) {
      existing.addEventListener('load', () => resolve())
      existing.addEventListener('error', () => reject(new Error('PayPal SDK load error')))
      return
    }

    const script = document.createElement('script')
    script.id = 'paypal-sdk'
    script.async = true
    const sdkCurrency = sdkConfig.value?.currency || currency.value
    script.src = `https://www.paypal.com/sdk/js?client-id=${encodeURIComponent(clientId)}&currency=${encodeURIComponent(sdkCurrency)}`
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('PayPal SDK load error'))
    document.head.appendChild(script)
  })
}

const openPayPalButtons = async (_data: any) => {
  paypalError.value = ''
  paypalLoading.value = true
  currentStep.value = 'paypal'

  try {
    const clientId = sdkConfig.value?.public_key
    if (!clientId) {
      throw new Error('PayPal client_id manquant (sdk_config.public_key)')
    }

    await loadPayPalSdk(clientId)

    const container = document.getElementById('paypal-buttons')
    if (!container) {
      throw new Error('PayPal container introuvable')
    }
    container.innerHTML = ''

    if (!window.paypal?.Buttons) {
      throw new Error('PayPal SDK non initialisé')
    }

    window.paypal.Buttons({
      createOrder: async () => {
        const resp = await fetch(`${API_URL}/transfer-service/api/v1/deposits/paypal/create-order`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${authStore.accessToken}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ transaction_id: transactionId.value })
        })

        const payload = await resp.json().catch(() => ({}))
        if (!resp.ok) {
          throw new Error(payload.error || 'Impossible de créer la commande PayPal')
        }

        return payload.order_id
      },

      onApprove: async (data: any) => {
        const resp = await fetch(`${API_URL}/transfer-service/api/v1/deposits/paypal/capture`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${authStore.accessToken}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            transaction_id: transactionId.value,
            order_id: data?.orderID
          })
        })

        const payload = await resp.json().catch(() => ({}))
        if (!resp.ok) {
          throw new Error(payload.error || 'Capture PayPal échouée')
        }

        if (payload.status === 'completed') {
          currentStep.value = 'success'
          emit('success', { transactionId: transactionId.value, amount: amount.value })
          walletStore.fetchWallets()
        } else {
          currentStep.value = 'pending'
        }
      },

      onError: (err: any) => {
        console.error('PayPal Buttons error:', err)
        paypalError.value = 'Erreur PayPal. Veuillez réessayer.'
      }
    }).render('#paypal-buttons')
  } catch (e: any) {
    console.error('PayPal init error:', e)
    paypalError.value = e?.message || 'Erreur lors de l\'initialisation PayPal'

		if (paymentUrl.value) {
			redirectToPayment(paymentUrl.value, 'Le paiement intégré PayPal n\'est pas disponible sur cet appareil. Vous pouvez continuer via la page PayPal.')
		}
  } finally {
    paypalLoading.value = false
  }
}

watch(currentStep, async () => {
	await nextTick()
	if (contentRef.value) {
		contentRef.value.scrollTop = 0
	}
})

// Flutterwave SDK Integration
const openFlutterwaveModal = (data: any) => {
  if (!window.FlutterwaveCheckout) {
    console.error('Flutterwave SDK not loaded')
    redirectToPayment(data.payment_url)
    return
  }

  window.FlutterwaveCheckout({
    public_key: sdkConfig.value.public_key,
    tx_ref: transactionId.value,
    amount: amount.value,
    currency: currency.value,
    payment_options: 'card,mobilemoney,ussd',
    customer: {
      email: authStore.user?.email,
      phone_number: phoneNumber.value,
      name: `${authStore.user?.first_name} ${authStore.user?.last_name}`
    },
    customizations: {
      title: 'Zekora - Recharge',
      description: `Recharge de ${formatAmount(amount.value)} ${currency.value}`,
      logo: '/logo.png'
    },
    callback: (response: any) => {
      if (response.status === 'successful' || response.status === 'completed') {
        currentStep.value = 'success'
        emit('success', { transactionId: transactionId.value, amount: amount.value })
        walletStore.fetchWallets()
      } else {
        currentStep.value = 'pending'
      }
    },
    onclose: () => {
      if (currentStep.value === 'processing') {
        currentStep.value = 'pending'
      }
    }
  })
}

// Paystack SDK Integration
const openPaystackModal = (data: any) => {
  if (!window.PaystackPop) {
    console.error('Paystack SDK not loaded')
    redirectToPayment(data.payment_url)
    return
  }

  const handler = window.PaystackPop.setup({
    key: sdkConfig.value.public_key,
    email: authStore.user?.email,
    amount: amount.value * 100, // Paystack uses kobo/pesewas
    currency: currency.value,
    ref: transactionId.value,
    channels: ['card', 'bank', 'ussd', 'mobile_money'],
    metadata: {
      user_id: authStore.user?.id,
      wallet_id: props.walletId
    },
    callback: (response: any) => {
      if (response.status === 'success') {
        currentStep.value = 'success'
        emit('success', { transactionId: transactionId.value, amount: amount.value })
        walletStore.fetchWallets()
      } else {
        currentStep.value = 'pending'
      }
    },
    onClose: () => {
      if (currentStep.value === 'processing') {
        currentStep.value = 'pending'
      }
    }
  })

  handler.openIframe()
}

const checkStatus = async () => {
  if (!transactionId.value) return

  try {
    const response = await fetch(`${API_URL}/transfer-service/api/v1/deposits/${transactionId.value}/status`, {
      headers: {
        'Authorization': `Bearer ${authStore.accessToken}`
      }
    })

    const data = await response.json()

    if (data.status === 'completed') {
      newBalance.value = data.new_balance
      currentStep.value = 'success'
      emit('success', { transactionId: transactionId.value, amount: amount.value })
      walletStore.fetchWallets()
    } else if (data.status === 'failed' || data.status === 'expired') {
      error.value = data.status_message || 'Le paiement a échoué'
      currentStep.value = 'failed'
    } else if (data.status === 'cancelled') {
      error.value = 'Paiement annulé'
      currentStep.value = 'failed'
    }
    // If still pending, do nothing
  } catch (err) {
    console.error('Failed to check status:', err)
  }
}

const resetForm = () => {
  currentStep.value = 'provider'
  selectedProvider.value = null
  selectedIsTestMode.value = null
  phoneNumber.value = ''
  error.value = ''
  transactionId.value = ''
  paymentUrl.value = ''
  newBalance.value = null
  paypalError.value = ''
  paypalLoading.value = false
}

const closeModal = () => {
  if (currentStep.value === 'success') {
    walletStore.fetchWallets()
  }
  resetForm()
  emit('close')
}

// Load SDK scripts
const loadSDKScripts = () => {
  // Flutterwave SDK
  if (!document.getElementById('flutterwave-sdk')) {
    const flwScript = document.createElement('script')
    flwScript.id = 'flutterwave-sdk'
    flwScript.src = 'https://checkout.flutterwave.com/v3.js'
    flwScript.async = true
    document.head.appendChild(flwScript)
  }

  // Paystack SDK
  if (!document.getElementById('paystack-sdk')) {
    const psScript = document.createElement('script')
    psScript.id = 'paystack-sdk'
    psScript.src = 'https://js.paystack.co/v1/inline.js'
    psScript.async = true
    document.head.appendChild(psScript)
  }
}

// Watch for modal open
watch(() => props.isOpen, async (isOpen) => {
  if (isOpen) {
    await authStore.initializeAuth()
    await fetchProviders()
    await fetchDepositNumbers()
  } else {
    resetForm()
  }
})

onMounted(() => {
  try {
    const isCoarse = window.matchMedia?.('(pointer: coarse)')?.matches || false
    const isSmallScreen = typeof window.innerWidth === 'number' ? window.innerWidth < 768 : false
    isMobile.value = isSmallScreen || isCoarse
  } catch (_e) {
    isMobile.value = false
  }
})

// Check URL params for callback
onMounted(() => {
  loadSDKScripts()

  const urlParams = new URLSearchParams(window.location.search)
  if (urlParams.get('deposit_callback') === 'true') {
    // Payment completed, refresh wallet
    walletStore.fetchWallets()
  }
})

// Type declarations for SDK globals
declare global {
  interface Window {
    FlutterwaveCheckout: (config: any) => void
    PaystackPop: {
      setup: (config: any) => { openIframe: () => void }
    }
    paypal?: any
  }
}
</script>

<style scoped>
@keyframes bounce-once {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.animate-bounce-once {
  animation: bounce-once 0.5s ease-in-out;
}

.animate-fade-in-up {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
