<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
      <div>
        <h1 class="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">Mes Associations</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Gérez vos tontines, groupes d'épargne et crédits mutuels</p>
      </div>
      <button @click="navigateTo('/associations/create')" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg flex items-center space-x-2 transition-colors">
        <PlusIcon class="w-5 h-5" />
        <span>Créer une association</span>
      </button>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in stats" :key="stat.name" class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
        <div class="flex items-center justify-between mb-2">
          <span class="text-gray-500 dark:text-gray-400 text-sm">{{ stat.name }}</span>
          <component :is="stat.icon" class="w-5 h-5 text-indigo-500" />
        </div>
        <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ stat.value }}</div>
        <div class="text-xs text-green-500 mt-1 flex items-center">
          <ArrowTrendingUpIcon class="w-3 h-3 mr-1" />
          {{ stat.change }}
        </div>
      </div>
    </div>

    <!-- Associations List -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-if="loading" v-for="i in 3" :key="i" class="bg-white dark:bg-gray-800 rounded-lg shadow-sm h-64 animate-pulse">
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 m-4"></div>
        <div class="h-32 bg-gray-200 dark:bg-gray-700 m-4 rounded"></div>
        <div class="h-8 bg-gray-200 dark:bg-gray-700 m-4 rounded"></div>
      </div>

      <div v-else-if="associations.length === 0" class="col-span-full text-center py-12 bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-600">
        <UserGroupIcon class="w-16 h-16 mx-auto text-gray-400 mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">Aucune association</h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">Vous n'avez pas encore rejoint d'association.</p>
        <button @click="navigateTo('/associations/create')" class="text-indigo-600 font-medium hover:text-indigo-500">
          Créer ma première association &rarr;
        </button>
      </div>

      <div v-else v-for="association in associations" :key="association.id" 
           @click="navigateTo(`/associations/${association.id}`)"
           class="group bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700 hover:shadow-md transition-all cursor-pointer overflow-hidden relative">
        
        <div class="absolute top-0 right-0 p-4">
          <span :class="[
            'px-2 py-1 rounded-full text-xs font-medium',
            getStatusColor(association.status)
          ]">
            {{ getStatusLabel(association.status) }}
          </span>
        </div>

        <div class="p-6">
          <div class="w-12 h-12 rounded-xl bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <component :is="getTypeIcon(association.type)" class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
          </div>
          
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-1 group-hover:text-indigo-600 transition-colors">{{ association.name }}</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-4 line-clamp-2">{{ association.description }}</p>
          
          <div class="flex items-center justify-between text-sm text-gray-600 dark:text-gray-300 bg-gray-50 dark:bg-gray-750 p-3 rounded-lg">
            <div class="flex items-center">
              <UsersIcon class="w-4 h-4 mr-1 text-gray-400" />
              <span>{{ association.total_members }} membres</span>
            </div>
            <div class="font-bold text-gray-900 dark:text-white">
              {{ formatCurrency(association.treasury_balance, association.currency) }}
            </div>
          </div>
        </div>
        
        <div class="px-6 py-4 border-t border-gray-100 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-800/50 flex justify-between items-center group-hover:bg-indigo-50/30 dark:group-hover:bg-indigo-900/10 transition-colors">
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ getTypeLabel(association.type) }}</span>
          <ArrowRightIcon class="w-4 h-4 text-gray-400 group-hover:text-indigo-600 group-hover:translate-x-1 transition-all" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  PlusIcon, 
  UserGroupIcon, 
  UsersIcon, 
  BanknotesIcon, 
  ScaleIcon,
  ArrowTrendingUpIcon,
  ArrowRightIcon,
  BriefcaseIcon
} from '@heroicons/vue/24/outline'

const { associationApi } = useApi()

const loading = ref(true)
const associations = ref<any[]>([])
const stats = ref([
  { name: 'Total Épargné', value: '2.5M FCFA', icon: BanknotesIcon, change: '+12%' },
  { name: 'Associations Actives', value: '3', icon: UserGroupIcon, change: '+1' },
  { name: 'Prêts en cours', value: '150K FCFA', icon: ScaleIcon, change: '0%' }, // ScaleIcon instead of CreditCardIcon for loans/balance
  { name: 'Contributions ce mois', value: '50K FCFA', icon: ArrowTrendingUpIcon, change: '+5%' }
])

const getTypeIcon = (type: string) => {
  switch(type) {
    case 'tontine': return ArrowTrendingUpIcon // Cycle/Turn
    case 'savings': return BanknotesIcon
    case 'credit': return ScaleIcon
    default: return UserGroupIcon
  }
}

const getTypeLabel = (type: string) => {
  switch(type) {
    case 'tontine': return 'Tontine Rotative'
    case 'savings': return 'Groupe d\'Épargne'
    case 'credit': return 'Crédit Mutuel'
    default: return 'Association'
  }
}

const getStatusColor = (status: string) => {
  switch(status) {
    case 'active': return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case 'suspended': return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'closed': return 'bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400'
    default: return 'bg-gray-100 text-gray-700'
  }
}

const getStatusLabel = (status: string) => {
  switch(status) {
    case 'active': return 'Actif'
    case 'suspended': return 'Suspendu'
    case 'closed': return 'Clôturé'
    default: return status
  }
}

const formatCurrency = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

onMounted(async () => {
  try {
    const response = await associationApi.getAll()
    associations.value = response.data
  } catch (err) {
    console.error('Failed to load associations', err)
    // Mock data for demo if API fails/not ready
    associations.value = [
      {
        id: '1',
        name: 'Tontine Famille Toure',
        type: 'tontine',
        description: 'Tontine mensuelle pour la famille et les amis proches.',
        total_members: 12,
        treasury_balance: 1200000,
        currency: 'XOF',
        status: 'active'
      },
      {
        id: '2',
        name: 'Épargne Commerçants',
        type: 'savings',
        description: 'Caisse de solidarité pour les commerçants du marché.',
        total_members: 45,
        treasury_balance: 5500000,
        currency: 'XOF',
        status: 'active'
      }
    ]
  } finally {
    loading.value = false
  }
})
</script>
