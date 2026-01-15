<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden transition-all hover:shadow-md">
    <!-- Service Header -->
    <div class="p-4 bg-gray-50 dark:bg-gray-750 border-b border-gray-100 dark:border-gray-700">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Service ID -->
        <div class="w-32">
          <input v-model="localService.id" 
            placeholder="ID unique" 
            class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm font-mono focus:ring-2 focus:ring-primary-500 focus:border-transparent" />
        </div>

        <!-- Service Name -->
        <div class="flex-1 min-w-[180px]">
          <input v-model="localService.name" 
            placeholder="Nom du service (ex: Inscription)" 
            class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm font-medium focus:ring-2 focus:ring-primary-500 focus:border-transparent" />
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
        <div class="flex items-center gap-2 bg-white dark:bg-gray-700/80 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Type:</span>
          <select v-model="localService.billing_type" 
            class="border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0 pr-6 cursor-pointer dark:bg-gray-700">
            <option value="FIXED" class="dark:bg-gray-800 dark:text-white">Forfait fixe</option>
            <option value="USAGE" class="dark:bg-gray-800 dark:text-white">À la consommation</option>
          </select>
        </div>

        <!-- Frequency -->
        <div class="flex items-center gap-2 bg-white dark:bg-gray-700/80 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Fréquence:</span>
          <select v-model="localService.billing_frequency" @change="onFrequencyChange"
            class="border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0 pr-6 cursor-pointer dark:bg-gray-700">
            <option value="DAILY" class="dark:bg-gray-800 dark:text-white">Journalier</option>
            <option value="WEEKLY" class="dark:bg-gray-800 dark:text-white">Hebdomadaire</option>
            <option value="MONTHLY" class="dark:bg-gray-800 dark:text-white">Mensuel</option>
            <option value="QUARTERLY" class="dark:bg-gray-800 dark:text-white">Trimestriel</option>
            <option value="ANNUALLY" class="dark:bg-gray-800 dark:text-white">Annuel</option>
            <option value="ONETIME" class="dark:bg-gray-800 dark:text-white">Paiement unique</option>
            <option value="CUSTOM" class="dark:bg-gray-800 dark:text-white">Personnalisé</option>
          </select>
        </div>

        <!-- Unit (for USAGE) -->
        <div v-if="localService.billing_type === 'USAGE'" 
          class="flex items-center gap-2 bg-white dark:bg-gray-700/80 px-3 py-2 rounded-xl border border-gray-200 dark:border-gray-600">
          <span class="text-xs text-gray-500 dark:text-gray-400">Unité:</span>
          <input v-model="localService.unit" 
            placeholder="kWh, m³..." 
            class="w-16 border-0 bg-transparent text-sm font-medium text-gray-900 dark:text-white focus:ring-0 p-0" />
        </div>

        <!-- Price -->
        <div class="flex items-center gap-2 bg-gradient-to-r from-primary-50 to-primary-100 dark:from-primary-900/40 dark:to-primary-800/30 px-3 py-2 rounded-xl border border-primary-200 dark:border-primary-700">
          <span class="text-xs text-primary-600 dark:text-primary-400">{{ localService.billing_type === 'USAGE' ? 'Prix/unité:' : 'Montant:' }}</span>
          <input v-model.number="localService.base_price" 
            type="number" 
            step="0.01"
            class="w-24 border-0 bg-transparent text-sm font-bold text-primary-700 dark:text-primary-300 focus:ring-0 p-0 text-right" />
          <span class="text-xs font-bold text-primary-600 dark:text-primary-400">{{ currency }}</span>
        </div>
      </div>
    </div>

    <!-- Advanced Options (Collapsible) -->
    <Transition name="slide">
      <div v-if="showAdvanced" class="p-4 bg-gray-50/50 dark:bg-gray-850 space-y-4">
        
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

          <div v-if="hasPenalty" class="grid grid-cols-2 md:grid-cols-4 gap-3">
            <div>
              <label class="block text-xs text-gray-500 mb-1">Type</label>
              <select v-model="localService.penalty_config.type" class="w-full px-2 py-1.5 rounded border-gray-200 dark:bg-gray-700 dark:border-gray-600 text-sm">
                <option value="FIXED">Montant fixe</option>
                <option value="PERCENTAGE">Pourcentage</option>
                <option value="HYBRID">Hybride</option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">{{ localService.penalty_config?.type === 'PERCENTAGE' ? 'Pourcentage' : 'Montant' }}</label>
              <input v-model.number="localService.penalty_config.value" type="number" class="w-full px-2 py-1.5 rounded border-gray-200 dark:bg-gray-700 dark:border-gray-600 text-sm" />
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">Fréquence</label>
              <select v-model="localService.penalty_config.frequency" class="w-full px-2 py-1.5 rounded border-gray-200 dark:bg-gray-700 dark:border-gray-600 text-sm">
                <option value="ONETIME">Une fois</option>
                <option value="DAILY">Par jour</option>
                <option value="WEEKLY">Par semaine</option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">Délai de grâce</label>
              <div class="flex items-center gap-1">
                <input v-model.number="localService.penalty_config.grace_period" type="number" class="w-full px-2 py-1.5 rounded border-gray-200 dark:bg-gray-700 dark:border-gray-600 text-sm" />
                <span class="text-xs text-gray-500">jours</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
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
    penalty_config: svc.penalty_config || null
  }
}

const localService = ref(getInitialService(props.service))
const showAdvanced = ref(false)
const scheduleMode = ref(props.service.payment_schedule?.length ? 'schedule' : 'interval')

const hasPenalty = computed(() => !!localService.value.penalty_config)

watch(localService, (val) => {
  emit('update', val)
}, { deep: true })

watch(() => props.service, (val) => {
  localService.value = getInitialService(val)
}, { deep: true })

const togglePenalty = () => {
  if (localService.value.penalty_config) {
    localService.value.penalty_config = null
  } else {
    localService.value.penalty_config = {
      type: 'PERCENTAGE',
      value: 10,
      frequency: 'ONETIME',
      grace_period: 5
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
