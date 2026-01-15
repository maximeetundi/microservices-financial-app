<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <Squares2X2Icon class="w-6 h-6 text-primary-500" />
          Groupes de Services
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Organisez vos services en catégories pour une meilleure gestion
        </p>
      </div>
      <button @click="addGroup" 
        class="px-4 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl text-sm font-semibold shadow-lg shadow-primary-500/25 transition-all flex items-center gap-2 transform hover:-translate-y-0.5">
        <PlusIcon class="w-5 h-5" />
        Nouveau Groupe
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="!groups?.length" 
      class="text-center py-16 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-2xl bg-gray-50/50 dark:bg-gray-800/50">
      <FolderPlusIcon class="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
      <h4 class="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">Aucun groupe de services</h4>
      <p class="text-gray-500 text-sm mb-6 max-w-sm mx-auto">
        Créez des groupes pour organiser vos services (ex: Scolarité, Transport, Abonnements)
      </p>
      <button @click="addGroup" 
        class="px-5 py-2.5 bg-primary-600 hover:bg-primary-700 text-white rounded-xl text-sm font-medium transition-colors">
        Créer mon premier groupe
      </button>
    </div>

    <!-- Groups Grid -->
    <div v-else class="grid grid-cols-1 gap-6">
      <TransitionGroup name="list">
        <div v-for="(group, gIdx) in groups" :key="group.id" 
          class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden transition-all hover:shadow-lg">
          
          <!-- Group Header -->
          <div class="p-5 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-gray-50 to-white dark:from-gray-900 dark:to-gray-800">
            <div class="flex flex-wrap items-center gap-4">
              <!-- Group Name -->
              <div class="flex-1 min-w-[200px]">
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Nom du groupe</label>
                <div class="relative">
                  <FolderIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input v-model="group.name" 
                    placeholder="ex: Scolarité" 
                    class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white font-medium focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" />
                </div>
              </div>

              <!-- Currency -->
              <div class="w-40">
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Devise</label>
                <select v-model="group.currency" 
                  class="w-full px-3 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white font-mono font-medium focus:ring-2 focus:ring-primary-500 focus:border-transparent cursor-pointer">
                  <option v-for="curr in currencies" :key="curr.code" :value="curr.code" class="bg-white dark:bg-gray-800 dark:text-white">
                    {{ curr.code }} - {{ curr.symbol }}
                  </option>
                </select>
              </div>

              <!-- Wallet Destination -->
              <div class="w-48">
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Wallet destination</label>
                <select v-model="group.wallet_id" 
                  class="w-full px-3 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white font-medium focus:ring-2 focus:ring-primary-500 focus:border-transparent cursor-pointer">
                  <option value="" class="bg-white dark:bg-gray-800 dark:text-white">-- Wallet par défaut --</option>
                  <option v-for="wallet in availableWallets" :key="wallet.id" :value="wallet.id" class="bg-white dark:bg-gray-800 dark:text-white">
                    {{ wallet.currency }} - {{ wallet.label || 'Principal' }}
                  </option>
                </select>
              </div>

              <!-- Private Toggle -->
              <div class="flex items-center">
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="group.is_private" class="sr-only peer" />
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary-600"></div>
                  <span class="ms-3 text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-1.5">
                    <LockClosedIcon v-if="group.is_private" class="w-4 h-4 text-amber-500" />
                    <LockOpenIcon v-else class="w-4 h-4 text-gray-400" />
                    {{ group.is_private ? 'Privé' : 'Public' }}
                  </span>
                </label>
              </div>

              <!-- Delete Button -->
              <button @click="removeGroup(gIdx)" 
                class="p-2.5 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-xl transition-colors"
                title="Supprimer ce groupe">
                <TrashIcon class="w-5 h-5" />
              </button>

              <!-- Collapse Toggle -->
              <button @click="toggleGroupCollapse(group)" 
                class="p-2.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl transition-colors">
                <ChevronDownIcon :class="['w-5 h-5 transition-transform', group._collapsed ? '-rotate-90' : '']" />
              </button>
            </div>

            <!-- Group Stats -->
            <div v-if="!group._collapsed" class="flex gap-4 mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
              <div class="text-sm">
                <span class="text-gray-500">Services:</span>
                <span class="font-bold text-gray-900 dark:text-white ml-1">{{ group.services?.length || 0 }}</span>
              </div>
              <div class="text-sm">
                <span class="text-gray-500">Visibilité:</span>
                <span :class="['font-medium ml-1', group.is_private ? 'text-amber-600' : 'text-green-600']">
                  {{ group.is_private ? 'Invitation seulement' : 'QR Code public' }}
                </span>
              </div>
            </div>
          </div>

          <!-- Services List (Collapsible) -->
          <div v-show="!group._collapsed" class="p-5">
            <div v-if="!group.services?.length" 
              class="text-center py-8 text-gray-400 border border-dashed border-gray-200 dark:border-gray-700 rounded-xl">
              <BoltIcon class="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p class="text-sm">Aucun service dans ce groupe</p>
            </div>

            <!-- Services Grid -->
            <div v-else class="space-y-4">
              <ServiceCard 
                v-for="(svc, sIdx) in group.services" 
                :key="svc.uid || sIdx"
                :service="svc"
                :currency="group.currency"
                @update="updateService(group, sIdx, $event)"
                @remove="removeService(group, sIdx)"
              />
            </div>

            <!-- Add Service Button -->
            <button @click="addService(group)" 
              class="w-full mt-4 py-3 border-2 border-dashed border-gray-200 dark:border-gray-600 rounded-xl text-sm text-gray-500 hover:text-primary-600 hover:border-primary-300 dark:hover:border-primary-700 transition-colors flex items-center justify-center gap-2">
              <PlusCircleIcon class="w-5 h-5" />
              Ajouter un service à "{{ group.name }}"
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup>
import { 
  Squares2X2Icon, PlusIcon, FolderPlusIcon, FolderIcon, 
  TrashIcon, ChevronDownIcon, BoltIcon, PlusCircleIcon,
  LockClosedIcon, LockOpenIcon
} from '@heroicons/vue/24/outline'
import ServiceCard from './ServiceCard.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  availableWallets: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const groups = computed({
  get: () => props.modelValue || [],
  set: (val) => emit('update:modelValue', val)
})

const currencies = [
  { code: 'XOF', symbol: 'FCFA', name: 'Franc CFA (BCEAO)' },
  { code: 'XAF', symbol: 'FCFA', name: 'Franc CFA (BEAC)' },
  { code: 'EUR', symbol: '€', name: 'Euro' },
  { code: 'USD', symbol: '$', name: 'Dollar US' },
  { code: 'GBP', symbol: '£', name: 'Livre Sterling' },
  { code: 'CAD', symbol: 'C$', name: 'Dollar Canadien' },
  { code: 'CHF', symbol: 'CHF', name: 'Franc Suisse' },
  { code: 'NGN', symbol: '₦', name: 'Naira' },
  { code: 'GHS', symbol: 'GH₵', name: 'Cedi' },
  { code: 'KES', symbol: 'KSh', name: 'Shilling Kenyan' },
  { code: 'ZAR', symbol: 'R', name: 'Rand' },
  { code: 'MAD', symbol: 'DH', name: 'Dirham Marocain' },
  { code: 'CDF', symbol: 'FC', name: 'Franc Congolais' },
  { code: 'GNF', symbol: 'GNF', name: 'Franc Guinéen' },
]

const addGroup = () => {
  const newGroup = {
    id: crypto.randomUUID(),
    name: '',
    is_private: false,
    currency: 'XOF',
    services: [],
    _collapsed: false
  }
  groups.value = [...groups.value, newGroup]
}

const removeGroup = (idx) => {
  if (confirm('Supprimer ce groupe et tous ses services ?')) {
    const updated = [...groups.value]
    updated.splice(idx, 1)
    groups.value = updated
  }
}

const toggleGroupCollapse = (group) => {
  group._collapsed = !group._collapsed
}

const addService = (group) => {
  if (!group.services) group.services = []
  group.services.push({
    id: '',
    name: '',
    billing_type: 'FIXED',
    billing_frequency: 'MONTHLY',
    base_price: 0,
    unit: '',
    uid: Date.now(),
    form_schema: [],
    penalty_config: null
  })
}

const removeService = (group, idx) => {
  group.services.splice(idx, 1)
}

const updateService = (group, idx, updatedService) => {
  group.services[idx] = updatedService
}
</script>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
