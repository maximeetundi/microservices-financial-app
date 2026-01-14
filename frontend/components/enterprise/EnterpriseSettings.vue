<template>
  <div class="space-y-8">
    <!-- Enterprise Profile Section -->
    <section class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-gray-50 to-white dark:from-gray-800 dark:to-gray-750">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <BuildingOffice2Icon class="w-5 h-5 text-primary-500" />
          Profil de l'entreprise
        </h3>
      </div>
      
      <div class="p-6 space-y-6">
        <!-- Logo & Name Row -->
        <div class="flex flex-col md:flex-row gap-6">
          <!-- Logo Upload -->
          <div class="flex flex-col items-center">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Logo</label>
            <div class="relative group">
              <div 
                v-if="enterprise.logo" 
                class="w-28 h-28 rounded-2xl overflow-hidden border-2 border-gray-200 dark:border-gray-600 bg-white shadow-sm">
                <img :src="enterprise.logo" alt="Logo" class="w-full h-full object-cover" />
              </div>
              <div 
                v-else 
                class="w-28 h-28 rounded-2xl bg-gradient-to-br from-primary-100 to-primary-50 dark:from-primary-900/30 dark:to-primary-800/20 flex items-center justify-center text-4xl font-bold text-primary-600 dark:text-primary-400 border-2 border-primary-200 dark:border-primary-700">
                {{ enterprise.name?.charAt(0) || '?' }}
              </div>
              <label class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover:opacity-100 rounded-2xl cursor-pointer transition-opacity">
                <div class="text-white text-center">
                  <CameraIcon class="w-6 h-6 mx-auto mb-1" />
                  <span class="text-xs">Modifier</span>
                </div>
                <input type="file" @change="handleLogoUpload" accept="image/*" class="hidden" />
              </label>
            </div>
            <button v-if="enterprise.logo" @click="removeLogo" class="mt-2 text-xs text-red-500 hover:text-red-700">
              Supprimer
            </button>
          </div>

          <!-- Enterprise Details -->
          <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nom de l'entreprise</label>
              <input 
                v-model="enterprise.name" 
                type="text" 
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent" 
                placeholder="Ma Super Entreprise" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Type d'activité</label>
              <select 
                v-model="enterprise.type" 
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                <option value="SERVICE">Service Général</option>
                <option value="SCHOOL">École / Éducation</option>
                <option value="TRANSPORT">Transport</option>
                <option value="UTILITY">Eau / Électricité / Gaz</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Numéro d'enregistrement</label>
              <input 
                v-model="enterprise.registration_number" 
                type="text" 
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent" 
                placeholder="RCCM / NIF (optionnel)" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Effectif</label>
              <select 
                v-model="enterprise.employee_count_range" 
                class="w-full px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                <option value="1-10">1 - 10 employés</option>
                <option value="11-50">11 - 50 employés</option>
                <option value="51-200">51 - 200 employés</option>
                <option value="201-500">201 - 500 employés</option>
                <option value="500+">Plus de 500 employés</option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Payroll Settings -->
    <section class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-green-50 to-white dark:from-green-900/20 dark:to-gray-800">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <BanknotesIcon class="w-5 h-5 text-green-500" />
          Configuration de la Paie
        </h3>
      </div>
      
      <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Jour de versement</label>
          <div class="flex items-center gap-2">
            <span class="text-gray-500">Le</span>
            <select 
              v-model.number="enterprise.settings.payroll_date" 
              class="px-4 py-2.5 rounded-xl border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent">
              <option v-for="d in 28" :key="d" :value="d">{{ d }}</option>
            </select>
            <span class="text-gray-500">de chaque mois</span>
          </div>
        </div>

        <div class="flex items-center">
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="enterprise.settings.auto_pay_salaries" class="sr-only peer" />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-green-300 dark:peer-focus:ring-green-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-green-600"></div>
            <span class="ms-3 text-sm font-medium text-gray-700 dark:text-gray-300">
              Paiement automatique des salaires
            </span>
          </label>
        </div>
      </div>
    </section>

    <!-- School-Specific Settings -->
    <section v-if="enterprise.type === 'SCHOOL'" class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-blue-50 to-white dark:from-blue-900/20 dark:to-gray-800">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <AcademicCapIcon class="w-5 h-5 text-blue-500" />
            Configuration École
          </h3>
          <button @click="addClass" class="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-1">
            <PlusIcon class="w-4 h-4" /> Ajouter une classe
          </button>
        </div>
      </div>
      
      <div class="p-6">
        <div v-if="!enterprise.school_config?.classes?.length" class="text-center py-8 text-gray-400">
          <AcademicCapIcon class="w-12 h-12 mx-auto mb-2 opacity-50" />
          <p>Aucune classe configurée</p>
        </div>

        <div v-else class="space-y-4">
          <div v-for="(cls, idx) in enterprise.school_config.classes" :key="idx" 
            class="bg-gray-50 dark:bg-gray-750 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
            
            <div class="flex flex-wrap gap-4 items-center mb-4">
              <div class="flex-1 min-w-[150px]">
                <input v-model="cls.name" placeholder="Nom de la classe (ex: CP, 6ème)" 
                  class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white font-medium" />
              </div>
              <div class="w-40">
                <div class="flex items-center gap-2">
                  <input v-model.number="cls.total_fees" type="number" placeholder="Frais total" 
                    class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white" />
                  <span class="text-sm text-gray-500 font-medium">XOF</span>
                </div>
              </div>
              <button @click="removeClass(idx)" class="p-2 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg">
                <TrashIcon class="w-5 h-5" />
              </button>
            </div>

            <!-- Tranches -->
            <div class="pl-4 border-l-2 border-blue-200 dark:border-blue-700 space-y-2">
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-gray-500 uppercase">Tranches de paiement</span>
                <button @click="addTranche(cls)" class="text-xs text-blue-600 hover:text-blue-700 font-medium flex items-center gap-1">
                  <PlusIcon class="w-3 h-3" /> Ajouter
                </button>
              </div>
              <div v-for="(tr, trIdx) in cls.tranches" :key="trIdx" class="flex items-center gap-2">
                <input v-model="tr.name" placeholder="1ère Tranche" class="flex-1 px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                <input v-model.number="tr.amount" type="number" placeholder="Montant" class="w-24 px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm text-right" />
                <input v-model="tr.due_date" type="date" class="px-2 py-1.5 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm" />
                <button @click="removeTranche(cls, trIdx)" class="p-1 text-red-400 hover:text-red-600">
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
              <div v-if="cls.tranches?.length" class="text-xs text-blue-600 font-medium mt-2">
                Total tranches: {{ cls.tranches.reduce((sum, t) => sum + (t.amount || 0), 0).toLocaleString() }} / {{ cls.total_fees?.toLocaleString() || 0 }} XOF
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Transport-Specific Settings -->
    <section v-if="enterprise.type === 'TRANSPORT'" class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="p-6 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-yellow-50 to-white dark:from-yellow-900/20 dark:to-gray-800">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <TruckIcon class="w-5 h-5 text-yellow-500" />
            Configuration Transport
          </h3>
          <button @click="addRoute" class="px-3 py-1.5 bg-yellow-600 hover:bg-yellow-700 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-1">
            <PlusIcon class="w-4 h-4" /> Ajouter une ligne
          </button>
        </div>
      </div>
      
      <div class="p-6">
        <div v-if="!enterprise.transport_config?.routes?.length" class="text-center py-8 text-gray-400">
          <TruckIcon class="w-12 h-12 mx-auto mb-2 opacity-50" />
          <p>Aucune ligne configurée</p>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="(route, idx) in enterprise.transport_config.routes" :key="idx" 
            class="bg-gray-50 dark:bg-gray-750 rounded-xl p-4 border border-gray-200 dark:border-gray-700 flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 flex items-center justify-center font-bold">
              {{ idx + 1 }}
            </div>
            <div class="flex-1">
              <input v-model="route.name" placeholder="Ligne 14 - Centre-Gare" 
                class="w-full px-3 py-2 rounded-lg border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-gray-900 dark:text-white font-medium mb-2" />
              <div class="flex items-center gap-2">
                <span class="text-xs text-gray-500">Tarif:</span>
                <input v-model.number="route.base_price" type="number" 
                  class="w-24 px-2 py-1 rounded border-gray-200 dark:border-gray-600 dark:bg-gray-700 text-sm text-right" />
                <span class="text-xs text-gray-500">XOF</span>
              </div>
            </div>
            <button @click="removeRoute(idx)" class="p-2 text-red-400 hover:text-red-600">
              <TrashIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Service Groups Manager -->
    <ServiceGroupsManager v-model="enterprise.service_groups" />
  </div>
</template>

<script setup>
import { 
  BuildingOffice2Icon, CameraIcon, BanknotesIcon, AcademicCapIcon, 
  TruckIcon, PlusIcon, TrashIcon, XMarkIcon
} from '@heroicons/vue/24/outline'
import ServiceGroupsManager from './ServiceGroupsManager.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'upload-logo'])

const enterprise = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// Ensure nested objects exist
if (!enterprise.value.settings) {
  enterprise.value.settings = { payroll_date: 25, auto_pay_salaries: false }
}
if (!enterprise.value.service_groups) {
  enterprise.value.service_groups = []
}

// Logo handlers
const handleLogoUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    emit('upload-logo', file)
  }
}

const removeLogo = () => {
  enterprise.value.logo = ''
}

// School Methods
const addClass = () => {
  if (!enterprise.value.school_config) enterprise.value.school_config = { classes: [] }
  enterprise.value.school_config.classes.push({ name: '', total_fees: 0, tranches: [] })
}

const removeClass = (idx) => {
  enterprise.value.school_config.classes.splice(idx, 1)
}

const addTranche = (cls) => {
  if (!cls.tranches) cls.tranches = []
  cls.tranches.push({ name: '', amount: 0, due_date: '' })
}

const removeTranche = (cls, idx) => {
  cls.tranches.splice(idx, 1)
}

// Transport Methods
const addRoute = () => {
  if (!enterprise.value.transport_config) enterprise.value.transport_config = { routes: [], zones: [] }
  enterprise.value.transport_config.routes.push({ name: '', base_price: 0 })
}

const removeRoute = (idx) => {
  enterprise.value.transport_config.routes.splice(idx, 1)
}
</script>
