<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden transition-all hover:shadow-md">
    <!-- Service Header -->
    <div class="p-4 bg-gray-50 dark:bg-gray-900 border-b border-gray-100 dark:border-gray-700">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Service ID -->
        <div class="w-32">
          <input v-model="localService.id" 
            placeholder="ID unique" 
            class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm font-mono focus:ring-2 focus:ring-primary-500 focus:border-transparent" />
        </div>

        <!-- Service Name -->
        <div class="flex-1 min-w-[180px]">
          <input v-model="localService.name" 
            placeholder="Nom du service (ex: Inscription)" 
            class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm font-medium focus:ring-2 focus:ring-primary-500 focus:border-transparent" />
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2">
          <button @click="showAdvanced = !showAdvanced" 
            :class="['p-2 rounded-lg transition-colors', showAdvanced ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/30' : 'text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700']"
            title="Options avancées">
            <AdjustmentsHorizontalIcon class="w-5 h-5" />
          </button>
          <button @click="$emit('remove')" 
            class="p-2 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
            title="Supprimer">
            <TrashIcon class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Basic Options Row -->
      <div class="flex flex-wrap items-center gap-3 mt-3">
        <!-- Billing Type -->
        <div class="flex items-center gap-2 bg-white dark:bg-gray-800 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Type:</span>
          <select v-model="localService.billing_type" 
            class="border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0 pr-6 cursor-pointer">
            <option value="FIXED" class="bg-white dark:bg-gray-800 dark:text-white">Forfait fixe</option>
            <option value="USAGE" class="bg-white dark:bg-gray-800 dark:text-white">À la consommation</option>
          </select>
        </div>

        <!-- Frequency -->
        <div class="flex items-center gap-2 bg-white dark:bg-gray-800 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Fréquence:</span>
          <select v-model="localService.billing_frequency" @change="onFrequencyChange"
            class="border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0 pr-6 cursor-pointer">
            <option value="DAILY" class="bg-white dark:bg-gray-800 dark:text-white">Journalier</option>
            <option value="WEEKLY" class="bg-white dark:bg-gray-800 dark:text-white">Hebdomadaire</option>
            <option value="MONTHLY" class="bg-white dark:bg-gray-800 dark:text-white">Mensuel</option>
            <option value="QUARTERLY" class="bg-white dark:bg-gray-800 dark:text-white">Trimestriel</option>
            <option value="ANNUALLY" class="bg-white dark:bg-gray-800 dark:text-white">Annuel</option>
            <option value="ONETIME" class="bg-white dark:bg-gray-800 dark:text-white">Paiement unique</option>
            <option value="CUSTOM" class="bg-white dark:bg-gray-800 dark:text-white">Personnalisé</option>
          </select>
        </div>

        <!-- Unit (for USAGE) -->
        <div v-if="localService.billing_type === 'USAGE'" 
          class="flex items-center gap-2 bg-white dark:bg-gray-800 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Unité:</span>
          <input v-model="localService.unit" 
            placeholder="kWh, m³..." 
            class="w-16 border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0" />
        </div>

        <!-- Pricing Mode (for USAGE) -->
        <div v-if="localService.billing_type === 'USAGE'" 
          class="flex items-center gap-2 bg-white dark:bg-gray-800 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Tarif:</span>
          <select v-model="localService.pricing_mode" 
            class="border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0 pr-6 cursor-pointer">
            <option value="FIXED" class="bg-white dark:bg-gray-800">Fixe</option>
            <option value="TIERED" class="bg-white dark:bg-gray-800">Par palier</option>
            <option value="THRESHOLD" class="bg-white dark:bg-gray-800">Par seuil</option>
          </select>
        </div>

        <!-- Price (only show for FIXED pricing or non-USAGE) -->
        <div v-if="localService.pricing_mode !== 'TIERED' && localService.pricing_mode !== 'THRESHOLD'" 
          class="flex items-center gap-2 bg-gradient-to-r from-primary-50 to-primary-100 dark:from-primary-900/40 dark:to-primary-800/30 px-3 py-2 rounded-xl border border-primary-200 dark:border-primary-700">
          <span class="text-xs text-primary-600 dark:text-primary-400">{{ localService.billing_type === 'USAGE' ? 'Prix/unité:' : (hasPaymentSchedule ? 'Total:' : 'Montant:') }}</span>
          <!-- Editable when no schedule, readonly when schedule exists -->
          <input v-if="!hasPaymentSchedule" 
            v-model.number="localService.base_price" 
            type="number" 
            step="0.01"
            class="w-24 border-0 bg-transparent text-sm font-bold text-primary-700 dark:text-primary-300 focus:ring-0 p-0 text-right" />
          <span v-else class="text-sm font-bold text-primary-700 dark:text-primary-300">{{ scheduleTotalAmount }}</span>
          <span class="text-xs font-bold text-primary-600 dark:text-primary-400">{{ currency }}</span>
        </div>
      </div>

      <!-- Pricing Tiers Builder (for TIERED or THRESHOLD mode) -->
      <div v-if="localService.billing_type === 'USAGE' && (localService.pricing_mode === 'TIERED' || localService.pricing_mode === 'THRESHOLD')"
        class="mt-3 p-3 bg-indigo-50 dark:bg-indigo-900/20 rounded-xl border border-indigo-200 dark:border-indigo-800">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-semibold text-indigo-800 dark:text-indigo-200">
            {{ localService.pricing_mode === 'TIERED' ? 'Paliers de prix' : 'Seuils de facturation' }}
          </span>
          <button @click="addPricingTier" 
            class="text-xs px-2 py-1 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
            + Ajouter
          </button>
        </div>

        <div v-if="!localService.pricing_tiers?.length" class="text-sm text-gray-500 italic text-center py-2">
          Aucun palier configuré
        </div>

        <div v-else class="space-y-2">
          <div v-for="(tier, idx) in localService.pricing_tiers" :key="idx"
            class="bg-white dark:bg-gray-800 p-3 rounded-lg border border-gray-200 dark:border-gray-700">
            <div class="flex flex-wrap items-center gap-2">
              <!-- Tier Label -->
              <input v-model="tier.label" 
                placeholder="Palier 1..." 
                class="w-24 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs" />
              
              <!-- Min -->
              <div class="flex items-center gap-1">
                <span class="text-xs text-gray-500">De:</span>
                <input v-model.number="tier.min_consumption" type="number" 
                  class="w-16 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs text-right" />
              </div>
              
              <!-- Max -->
              <div class="flex items-center gap-1">
                <span class="text-xs text-gray-500">À:</span>
                <input v-model.number="tier.max_consumption" type="number" placeholder="-1=∞"
                  class="w-16 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs text-right" />
              </div>
              
              <!-- Price/Unit -->
              <div class="flex items-center gap-1">
                <span class="text-xs text-gray-500">Prix:</span>
                <input v-model.number="tier.price_per_unit" type="number" step="0.01"
                  class="w-16 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs text-right" />
              </div>

              <!-- Fixed Bonus -->
              <div class="flex items-center gap-1">
                <span class="text-xs text-gray-500">+Fixe:</span>
                <input v-model.number="tier.fixed_bonus" type="number" step="0.01" placeholder="0"
                  class="w-14 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs text-right" />
              </div>

              <!-- Percent Bonus -->
              <div class="flex items-center gap-1">
                <span class="text-xs text-gray-500">+%:</span>
                <input v-model.number="tier.percent_bonus" type="number" step="0.1" placeholder="0"
                  class="w-12 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-xs text-right" />
              </div>

              <!-- Delete -->
              <button @click="removePricingTier(idx)" class="p-1 text-red-400 hover:text-red-600">
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Advanced Options (Collapsible) -->
    <Transition name="slide">
      <div v-if="showAdvanced" class="p-4 bg-gray-50/50 dark:bg-gray-900 space-y-4">
        
        <!-- Custom Interval -->
        <div v-if="localService.billing_frequency === 'CUSTOM'" 
          class="p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
          <div class="flex items-center gap-2 mb-3">
            <CalendarIcon class="w-5 h-5 text-amber-600" />
            <span class="font-semibold text-amber-800 dark:text-amber-200">Calendrier personnalisé</span>
          </div>
          
          <div class="flex items-center gap-4 mb-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="scheduleMode" value="interval" class="text-primary-600 focus:ring-primary-500" />
              <span class="text-sm text-gray-700 dark:text-gray-300">Intervalle fixe</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="scheduleMode" value="schedule" class="text-primary-600 focus:ring-primary-500" />
              <span class="text-sm text-gray-700 dark:text-gray-300">Calendrier de paiement</span>
            </label>
          </div>

          <!-- Fixed Interval -->
          <div v-if="scheduleMode === 'interval'" class="flex items-center gap-2">
            <span class="text-sm text-gray-600">Tous les</span>
            <input v-model.number="localService.custom_interval" 
              type="number" 
              min="1"
              class="w-20 px-3 py-2 rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 text-center" />
            <span class="text-sm text-gray-600">jours</span>
          </div>

          <!-- Payment Schedule -->
          <div v-else class="space-y-3">
            <div v-for="(item, idx) in (localService.payment_schedule || [])" :key="idx" 
              class="flex items-center gap-2 bg-white dark:bg-gray-800 p-3 rounded-lg">
              <input v-model="item.name" placeholder="Période" class="flex-1 px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
              <input v-model="item.start_date" type="date" class="px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
              <span class="text-gray-400">→</span>
              <input v-model="item.end_date" type="date" class="px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
              <input v-model.number="item.amount" type="number" placeholder="Montant" class="w-24 px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm text-right" />
              <button @click="removeScheduleItem(idx)" class="p-1 text-red-400 hover:text-red-600">
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
            <button @click="addScheduleItem" class="text-sm text-primary-600 hover:text-primary-700 font-medium flex items-center gap-1">
              <PlusIcon class="w-4 h-4" /> Ajouter une période
            </button>
          </div>
        </div>

        <!-- Form Builder -->
        <div class="p-4 bg-indigo-50 dark:bg-indigo-900/20 rounded-xl border border-indigo-200 dark:border-indigo-800">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <DocumentTextIcon class="w-5 h-5 text-indigo-600" />
              <span class="font-semibold text-indigo-800 dark:text-indigo-200">Formulaire d'inscription</span>
            </div>
            <button @click="addFormField" class="text-sm text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1">
              <PlusIcon class="w-4 h-4" /> Champ
            </button>
          </div>

          <div v-if="!localService.form_schema?.length" class="text-sm text-gray-500 italic text-center py-4">
            Aucun champ personnalisé (optionnel)
          </div>

          <div v-else class="space-y-2">
            <div v-for="(field, idx) in (localService.form_schema || [])" :key="idx" 
              class="flex items-center gap-2 bg-white dark:bg-gray-800 p-2 rounded-lg">
              <input v-model="field.label" placeholder="Libellé" class="flex-1 px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm" />
              <select v-model="field.type" class="px-2 py-1.5 rounded border border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-white text-sm">
                <option value="text" class="dark:bg-gray-800">Texte</option>
                <option value="number" class="dark:bg-gray-800">Nombre</option>
                <option value="date" class="dark:bg-gray-800">Date</option>
                <option value="email" class="dark:bg-gray-800">Email</option>
                <option value="tel" class="dark:bg-gray-800">Téléphone</option>
                <option value="checkbox" class="dark:bg-gray-800">Case à cocher</option>
                <option value="select" class="dark:bg-gray-800">Liste déroulante</option>
              </select>
              <label class="flex items-center gap-1 text-xs text-gray-500">
                <input type="checkbox" v-model="field.required" class="rounded border-gray-300 text-primary-600" />
                Requis
              </label>
              <button @click="removeFormField(idx)" class="p-1 text-red-400 hover:text-red-600">
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Penalty Config -->
        <div class="p-4 bg-red-50 dark:bg-red-900/20 rounded-xl border border-red-200 dark:border-red-800">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <ExclamationTriangleIcon class="w-5 h-5 text-red-600" />
              <span class="font-semibold text-red-800 dark:text-red-200">Pénalités de retard</span>
            </div>
            <button @click="togglePenalty" 
              :class="['px-3 py-1 rounded-full text-xs font-medium transition-colors', hasPenalty ? 'bg-red-600 text-white' : 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-400']">
              {{ hasPenalty ? 'Activé' : 'Désactivé' }}
            </button>
          </div>

          <div v-if="hasPenalty" class="space-y-3">
            <!-- Row 1: Type, Amounts -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
              <div>
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Type</label>
                <select v-model="localService.penalty_config.type" class="w-full px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm">
                  <option value="FIXED" class="dark:bg-gray-800">Montant fixe</option>
                  <option value="PERCENTAGE" class="dark:bg-gray-800">Pourcentage</option>
                  <option value="HYBRID" class="dark:bg-gray-800">Hybride</option>
                </select>
              </div>
              
              <!-- Amount field (for FIXED and HYBRID) -->
              <div v-if="localService.penalty_config?.type === 'FIXED' || localService.penalty_config?.type === 'HYBRID'">
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Montant</label>
                <input v-model.number="localService.penalty_config.value" type="number" class="w-full px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm" />
              </div>
              
              <!-- Percentage field (for PERCENTAGE and HYBRID) -->
              <div v-if="localService.penalty_config?.type === 'PERCENTAGE' || localService.penalty_config?.type === 'HYBRID'">
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Pourcentage (%)</label>
                <input v-model.number="localService.penalty_config.percentage" type="number" step="0.1" class="w-full px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm" />
              </div>
              
              <div>
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Fréquence d'application</label>
                <select v-model="localService.penalty_config.frequency" class="w-full px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm">
                  <option value="DAILY" class="dark:bg-gray-800">Par jour</option>
                  <option value="WEEKLY" class="dark:bg-gray-800">Par semaine</option>
                  <option value="MONTHLY" class="dark:bg-gray-800">Par mois</option>
                  <option value="QUARTERLY" class="dark:bg-gray-800">Par trimestre</option>
                  <option value="SEMIANNUAL" class="dark:bg-gray-800">Par semestre</option>
                  <option value="ANNUAL" class="dark:bg-gray-800">Par an</option>
                </select>
              </div>
            </div>

            <!-- Row 2: Grace Period & Max Penalty Date -->
            <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
              <div>
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Délai de grâce</label>
                <div class="flex items-center gap-1">
                  <input v-model.number="localService.penalty_config.grace_period" type="number" class="w-16 px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm" />
                  <select v-model="localService.penalty_config.grace_unit" class="flex-1 px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm">
                    <option value="DAYS" class="dark:bg-gray-800">jours</option>
                    <option value="WEEKS" class="dark:bg-gray-800">semaines</option>
                    <option value="MONTHS" class="dark:bg-gray-800">mois</option>
                  </select>
                </div>
              </div>
              
              <div>
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Limite des pénalités</label>
                <div class="flex items-center gap-1">
                  <input v-model.number="localService.penalty_config.max_penalty_months" type="number" class="w-16 px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm" />
                  <span class="text-xs text-gray-500 dark:text-gray-400">mois max</span>
                </div>
              </div>
              
              <div>
                <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">Plafond montant</label>
                <div class="flex items-center gap-1">
                  <input v-model.number="localService.penalty_config.max_penalty_amount" type="number" placeholder="∞" class="w-full px-2 py-1.5 rounded border border-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:text-white text-sm" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { 
  TrashIcon, AdjustmentsHorizontalIcon, CalendarIcon, DocumentTextIcon,
  PlusIcon, XMarkIcon, ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  service: {
    type: Object,
    required: true
  },
  currency: {
    type: String,
    default: 'XOF'
  }
})

const emit = defineEmits(['update', 'remove'])

// Initialize with safe defaults
const getInitialService = (svc) => {
  return {
    ...svc,
    payment_schedule: svc.payment_schedule || [],
    custom_interval: svc.custom_interval || 30,
    form_schema: svc.form_schema || [],
    penalty_config: svc.penalty_config || null,
    pricing_mode: svc.pricing_mode || 'FIXED',
    pricing_tiers: svc.pricing_tiers || []
  }
}

const localService = ref(getInitialService(props.service))
const showAdvanced = ref(false)
const scheduleMode = ref(props.service.payment_schedule?.length ? 'schedule' : 'interval')

const hasPenalty = computed(() => !!localService.value.penalty_config)

// Check if using payment schedule with periods
const hasPaymentSchedule = computed(() => {
  return localService.value.billing_frequency === 'CUSTOM' && 
         scheduleMode.value === 'schedule' && 
         localService.value.payment_schedule?.length > 0
})

// Calculate total amount from all periods
const scheduleTotalAmount = computed(() => {
  if (!hasPaymentSchedule.value) return 0
  return (localService.value.payment_schedule || []).reduce((sum, item) => sum + (item.amount || 0), 0)
})

// Auto-update base_price when schedule changes
watch(scheduleTotalAmount, (newTotal) => {
  if (hasPaymentSchedule.value && newTotal > 0) {
    localService.value.base_price = newTotal
  }
})

// Prevent infinite loop: only emit if actually changed by user
let isUpdatingFromProps = false

watch(localService, (val) => {
  if (!isUpdatingFromProps) {
    emit('update', val)
  }
}, { deep: true })

watch(() => props.service, (val) => {
  // Only update if the external prop actually changed
  if (JSON.stringify(val) !== JSON.stringify(localService.value)) {
    isUpdatingFromProps = true
    localService.value = getInitialService(val)
    // Reset flag after Vue's next tick to allow future user changes
    nextTick(() => {
      isUpdatingFromProps = false
    })
  }
}, { deep: true })

const togglePenalty = () => {
  if (localService.value.penalty_config) {
    localService.value.penalty_config = null
  } else {
    localService.value.penalty_config = {
      type: 'PERCENTAGE',
      value: 0,
      percentage: 10,
      frequency: 'MONTHLY',
      grace_period: 5,
      grace_unit: 'DAYS',
      max_penalty_months: 6, // Default 6 months max
      max_penalty_amount: 0 // 0 = no limit
    }
  }
}

const addFormField = () => {
  if (!localService.value.form_schema) localService.value.form_schema = []
  localService.value.form_schema.push({
    key: `field_${Date.now()}`,
    label: '',
    type: 'text',
    required: false
  })
}

const removeFormField = (idx) => {
  localService.value.form_schema.splice(idx, 1)
}

const addScheduleItem = () => {
  if (!localService.value.payment_schedule) localService.value.payment_schedule = []
  localService.value.payment_schedule.push({
    name: '',
    start_date: '',
    end_date: '',
    amount: 0
  })
}

const removeScheduleItem = (idx) => {
  localService.value.payment_schedule.splice(idx, 1)
}

// Pricing Tier Helpers
const addPricingTier = () => {
  if (!localService.value.pricing_tiers) localService.value.pricing_tiers = []
  const tierCount = localService.value.pricing_tiers.length
  const lastMax = tierCount > 0 ? localService.value.pricing_tiers[tierCount - 1].max_consumption : 0
  
  localService.value.pricing_tiers.push({
    label: `Palier ${tierCount + 1}`,
    min_consumption: lastMax > 0 ? lastMax : 0,
    max_consumption: -1, // -1 means unlimited
    price_per_unit: localService.value.base_price || 0,
    fixed_bonus: 0,
    percent_bonus: 0
  })
}

const removePricingTier = (idx) => {
  localService.value.pricing_tiers.splice(idx, 1)
}

// Auto-open advanced options when selecting CUSTOM frequency
const onFrequencyChange = () => {
  if (localService.value.billing_frequency === 'CUSTOM') {
    showAdvanced.value = true
    // Initialize payment_schedule if not exists
    if (!localService.value.payment_schedule) {
      localService.value.payment_schedule = []
    }
    // Initialize custom_interval if not exists
    if (!localService.value.custom_interval) {
      localService.value.custom_interval = 30
    }
  }
}
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
}
.slide-enter-to,
.slide-leave-from {
  opacity: 1;
  max-height: 500px;
}
</style>
