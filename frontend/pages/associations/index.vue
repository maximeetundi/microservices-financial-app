<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-surface p-6 rounded-2xl border border-secondary-200 dark:border-secondary-700">
        <div>
          <h1 class="text-2xl font-bold bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">Mes Associations</h1>
          <p class="text-muted mt-1">Gérez vos tontines, groupes d'épargne et crédits mutuels</p>
        </div>
        <button @click="navigateTo('/associations/create')" class="btn-premium flex items-center space-x-2">
          <PlusIcon class="w-5 h-5" />
          <span>Créer une association</span>
        </button>
      </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in computedStats" :key="stat.name" class="bg-surface p-4 rounded-2xl border border-secondary-200 dark:border-secondary-700">
        <div class="flex items-center justify-between mb-2">
          <span class="text-muted text-sm">{{ stat.name }}</span>
          <component :is="stat.icon" class="w-5 h-5 text-primary" />
        </div>
        <div class="text-2xl font-bold text-base">{{ stat.value }}</div>
        <div v-if="stat.change" class="text-xs text-success mt-1 flex items-center">
          <ArrowTrendingUpIcon class="w-3 h-3 mr-1" />
          {{ stat.change }}
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="i in 3" :key="i" class="bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700 h-64 animate-pulse">
        <div class="h-4 bg-secondary-200 dark:bg-secondary-700 rounded w-3/4 m-4"></div>
        <div class="h-32 bg-secondary-200 dark:bg-secondary-700 m-4 rounded"></div>
        <div class="h-8 bg-secondary-200 dark:bg-secondary-700 m-4 rounded"></div>
      </div>
    </div>

    <!-- No Associations -->
    <div v-else-if="associations.length === 0" class="text-center py-12 bg-surface rounded-2xl border border-dashed border-secondary-300 dark:border-secondary-600">
      <UserGroupIcon class="w-16 h-16 mx-auto text-muted mb-4" />
      <h3 class="text-lg font-medium text-base">Aucune association</h3>
      <p class="text-muted mb-6">Vous n'avez pas encore rejoint d'association.</p>
      <button @click="navigateTo('/associations/create')" class="text-primary font-medium hover:underline">
        Créer ma première association &rarr;
      </button>
    </div>

    <!-- Associations List -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="association in associations" :key="association.id" 
           @click="navigateTo(`/associations/${association.id}`)"
           class="group bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700 hover:border-primary/50 hover:shadow-lg hover:shadow-primary/10 transition-all cursor-pointer overflow-hidden relative">
        
        <div class="absolute top-0 right-0 p-4">
          <span :class="[
            'px-2 py-1 rounded-full text-xs font-medium',
            getStatusColor(association.status)
          ]">
            {{ getStatusLabel(association.status) }}
          </span>
        </div>

        <div class="p-6">
          <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <component :is="getTypeIcon(association.type)" class="w-6 h-6 text-primary" />
          </div>
          
          <h3 class="text-lg font-bold text-base mb-1 group-hover:text-primary transition-colors">{{ association.name }}</h3>
          <p class="text-sm text-muted mb-4 line-clamp-2">{{ association.description }}</p>
          
          <div class="flex items-center justify-between text-sm bg-surface-hover p-3 rounded-xl">
            <div class="flex items-center text-muted">
              <UsersIcon class="w-4 h-4 mr-1" />
              <span>{{ association.total_members || 0 }} membres</span>
            </div>
            <div class="font-bold text-base">
              {{ formatCurrency(association.treasury_balance || 0, association.currency) }}
            </div>
          </div>
        </div>
        
        <div class="px-6 py-4 border-t border-secondary-200 dark:border-secondary-700 bg-surface-hover flex justify-between items-center group-hover:bg-primary/5 transition-colors">
          <span class="text-xs text-muted">{{ getTypeLabel(association.type) }}</span>
          <ArrowRightIcon class="w-4 h-4 text-muted group-hover:text-primary group-hover:translate-x-1 transition-all" />
        </div>
      </div>
    </div>
  </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  PlusIcon, 
  UserGroupIcon, 
  UsersIcon, 
  BanknotesIcon, 
  ScaleIcon,
  ArrowTrendingUpIcon,
  ArrowRightIcon
} from '@heroicons/vue/24/outline'

import { associationAPI } from '~/composables/useApi'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const loading = ref(true)
const associations = ref<any[]>([])

// Compute stats from real data
const computedStats = computed(() => {
  const totalBalance = associations.value.reduce((sum, a) => sum + (a.treasury_balance || 0), 0)
  const activeCount = associations.value.filter(a => a.status === 'active').length
  const totalMembers = associations.value.reduce((sum, a) => sum + (a.total_members || 0), 0)
  const currency = associations.value[0]?.currency || 'XOF'
  
  return [
    { 
      name: 'Total Épargné', 
      value: formatCurrency(totalBalance, currency), 
      icon: BanknotesIcon 
    },
    { 
      name: 'Associations Actives', 
      value: String(activeCount), 
      icon: UserGroupIcon 
    },
    { 
      name: 'Membres Total', 
      value: String(totalMembers), 
      icon: UsersIcon 
    },
    { 
      name: 'Associations', 
      value: String(associations.value.length), 
      icon: ScaleIcon 
    }
  ]
})

const getTypeIcon = (type: string) => {
  switch(type) {
    case 'tontine': return ArrowTrendingUpIcon
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
    case 'active': return 'bg-success/20 text-success'
    case 'suspended': return 'bg-warning/20 text-warning'
    case 'closed': return 'bg-error/20 text-error'
    default: return 'bg-secondary-200 text-secondary-700'
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
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount || 0)
}

onMounted(async () => {
  try {
    const response = await associationAPI.getAll()
    associations.value = response.data?.associations || response.data || []
  } catch (err) {
    console.error('Failed to load associations', err)
    associations.value = []
  } finally {
    loading.value = false
  }
})
</script>

