<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
      <div class="bg-gradient-to-r from-indigo-600 to-purple-600 px-6 py-8">
        <div class="flex items-start justify-between">
          <div class="flex items-center space-x-4">
            <button @click="navigateTo('/associations')" class="p-2 rounded-full bg-white/10 hover:bg-white/20 transition-colors">
              <ArrowLeftIcon class="w-5 h-5 text-white" />
            </button>
            <div>
              <h1 class="text-2xl font-bold text-white">{{ association?.name || 'Chargement...' }}</h1>
              <p class="text-indigo-100 mt-1">{{ getTypeLabel(association?.type) }}</p>
            </div>
          </div>
          <span class="px-3 py-1 rounded-full text-sm font-medium bg-white/20 text-white">
            {{ association?.status === 'active' ? 'Actif' : association?.status }}
          </span>
        </div>
      </div>
      
      <!-- Quick Stats Bar -->
      <div class="grid grid-cols-4 divide-x divide-gray-100 dark:divide-gray-700">
        <div class="p-4 text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ association?.total_members || 0 }}</div>
          <div class="text-xs text-gray-500">Membres</div>
        </div>
        <div class="p-4 text-center">
          <div class="text-2xl font-bold text-indigo-600">{{ formatCurrency(association?.treasury_balance || 0) }}</div>
          <div class="text-xs text-gray-500">Caisse</div>
        </div>
        <div class="p-4 text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.contributions_this_month || 0 }}</div>
          <div class="text-xs text-gray-500">Cotisations ce mois</div>
        </div>
        <div class="p-4 text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.pending_loans || 0 }}</div>
          <div class="text-xs text-gray-500">Prêts en attente</div>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
      <div class="border-b border-gray-100 dark:border-gray-700">
        <nav class="flex space-x-8 px-6" aria-label="Tabs">
          <button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id"
                  class="py-4 px-1 border-b-2 font-medium text-sm transition-colors flex items-center space-x-2"
                  :class="activeTab === tab.id ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700'">
            <component :is="tab.icon" class="w-5 h-5" />
            <span>{{ tab.name }}</span>
          </button>
        </nav>
      </div>

      <div class="p-6">
        <!-- Members Tab -->
        <div v-if="activeTab === 'members'">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">Membres</h3>
            <button class="text-sm text-indigo-600 hover:text-indigo-500 font-medium">+ Inviter</button>
          </div>
          <div class="space-y-3">
            <div v-for="member in members" :key="member.id" 
                 class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-750 rounded-lg">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center">
                  <span class="text-indigo-600 font-medium">{{ getInitials(member.user_name) }}</span>
                </div>
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">{{ member.user_name }}</div>
                  <div class="text-xs text-gray-500">{{ getRoleLabel(member.role) }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-sm font-medium text-gray-900 dark:text-white">{{ formatCurrency(member.contributions_paid) }}</div>
                <div class="text-xs text-gray-500">cotisé</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Treasury Tab -->
        <div v-if="activeTab === 'treasury'">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">Trésorerie</h3>
            <button @click="showContributionModal = true" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
              + Cotisation
            </button>
          </div>
          
          <div class="grid grid-cols-3 gap-4 mb-6">
            <div class="bg-green-50 dark:bg-green-900/20 p-4 rounded-lg">
              <div class="text-sm text-green-600">Total Cotisations</div>
              <div class="text-2xl font-bold text-green-700">{{ formatCurrency(treasury.total_contributions) }}</div>
            </div>
            <div class="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg">
              <div class="text-sm text-red-600">Total Prêts</div>
              <div class="text-2xl font-bold text-red-700">{{ formatCurrency(treasury.total_loans) }}</div>
            </div>
            <div class="bg-indigo-50 dark:bg-indigo-900/20 p-4 rounded-lg">
              <div class="text-sm text-indigo-600">Solde Caisse</div>
              <div class="text-2xl font-bold text-indigo-700">{{ formatCurrency(treasury.total_balance) }}</div>
            </div>
          </div>

          <h4 class="font-medium text-gray-900 dark:text-white mb-3">Dernières transactions</h4>
          <div class="space-y-2">
            <div v-for="tx in treasury.transactions?.slice(0, 10)" :key="tx.id" 
                 class="flex items-center justify-between py-3 border-b border-gray-100 dark:border-gray-700 last:border-0">
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 rounded-full flex items-center justify-center"
                     :class="tx.type === 'contribution' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'">
                  <ArrowUpIcon v-if="tx.type === 'contribution'" class="w-4 h-4" />
                  <ArrowDownIcon v-else class="w-4 h-4" />
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-900 dark:text-white">{{ tx.description }}</div>
                  <div class="text-xs text-gray-500">{{ formatDate(tx.created_at) }}</div>
                </div>
              </div>
              <div :class="tx.type === 'contribution' ? 'text-green-600' : 'text-red-600'" class="font-medium">
                {{ tx.type === 'contribution' ? '+' : '-' }}{{ formatCurrency(tx.amount) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Meetings Tab -->
        <div v-if="activeTab === 'meetings'">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">Réunions</h3>
            <button class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
              + Planifier
            </button>
          </div>
          <div v-if="meetings.length === 0" class="text-center py-12 text-gray-500">
            <CalendarIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p>Aucune réunion planifiée</p>
          </div>
          <div v-else class="space-y-3">
            <div v-for="meeting in meetings" :key="meeting.id" class="p-4 border border-gray-100 dark:border-gray-700 rounded-lg">
              <div class="flex justify-between items-start">
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">{{ meeting.title }}</div>
                  <div class="text-sm text-gray-500 mt-1">{{ meeting.location }}</div>
                </div>
                <span class="px-2 py-1 rounded-full text-xs font-medium"
                      :class="meeting.status === 'scheduled' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'">
                  {{ meeting.status === 'scheduled' ? 'Planifiée' : 'Terminée' }}
                </span>
              </div>
              <div class="mt-3 text-sm text-indigo-600 font-medium">
                📅 {{ formatDate(meeting.date) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Loans Tab -->
        <div v-if="activeTab === 'loans'">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">Prêts</h3>
            <button @click="showLoanModal = true" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
              + Demander un prêt
            </button>
          </div>
          <div class="text-center py-12 text-gray-500">
            <BanknotesIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p>Aucun prêt en cours</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <ContributionModal 
      :show="showContributionModal" 
      :association-id="id" 
      :currency="association?.currency || 'XOF'"
      @close="showContributionModal = false"
      @success="handleContributionSuccess"
    />
    <LoanRequestModal
      :show="showLoanModal"
      :association-id="id"
      :currency="association?.currency || 'XOF'"
      :interest-rate="5"
      @close="showLoanModal = false"
      @success="handleLoanSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  ArrowLeftIcon,
  UsersIcon,
  BanknotesIcon,
  CalendarIcon,
  ScaleIcon,
  ArrowUpIcon,
  ArrowDownIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const { associationApi } = useApi()

const id = route.params.id as string
const association = ref<any>(null)
const members = ref<any[]>([])
const meetings = ref<any[]>([])
const treasury = ref<any>({ total_balance: 0, total_contributions: 0, total_loans: 0, transactions: [] })
const stats = ref({ contributions_this_month: 0, pending_loans: 0 })
const activeTab = ref('members')
const showContributionModal = ref(false)
const showLoanModal = ref(false)

const tabs = [
  { id: 'members', name: 'Membres', icon: UsersIcon },
  { id: 'treasury', name: 'Trésorerie', icon: BanknotesIcon },
  { id: 'meetings', name: 'Réunions', icon: CalendarIcon },
  { id: 'loans', name: 'Prêts', icon: ScaleIcon }
]

const handleContributionSuccess = async () => {
  // Reload treasury data
  try {
    const treasuryRes = await associationApi.getTreasury(id)
    treasury.value = treasuryRes.data || treasury.value
    const assocRes = await associationApi.get(id)
    association.value = assocRes.data
  } catch (err) {
    console.error('Failed to refresh', err)
  }
}

const handleLoanSuccess = () => {
  // Reload loans - will show in pending
  activeTab.value = 'loans'
}

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    tontine: 'Tontine Rotative',
    savings: "Groupe d'Épargne",
    credit: 'Crédit Mutuel',
    general: 'Association'
  }
  return labels[type] || type
}

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    president: 'Président',
    treasurer: 'Trésorier',
    secretary: 'Secrétaire',
    member: 'Membre'
  }
  return labels[role] || role
}

const getInitials = (name: string) => {
  if (!name) return '?'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: association.value?.currency || 'XOF' }).format(amount || 0)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const [assocRes, membersRes, treasuryRes] = await Promise.all([
      associationApi.get(id),
      associationApi.getMembers(id),
      associationApi.getTreasury(id)
    ])
    association.value = assocRes.data
    members.value = membersRes.data || []
    treasury.value = treasuryRes.data || treasury.value
  } catch (err) {
    console.error('Failed to load association', err)
    // Mock data for demo
    association.value = {
      id,
      name: 'Tontine Famille Toure',
      type: 'tontine',
      total_members: 12,
      treasury_balance: 1200000,
      currency: 'XOF',
      status: 'active'
    }
    members.value = [
      { id: '1', user_name: 'Mamadou Toure', role: 'president', contributions_paid: 150000 },
      { id: '2', user_name: 'Fatou Diallo', role: 'treasurer', contributions_paid: 150000 },
      { id: '3', user_name: 'Ibrahim Kone', role: 'member', contributions_paid: 100000 }
    ]
    treasury.value = {
      total_balance: 1200000,
      total_contributions: 1500000,
      total_loans: 300000,
      transactions: [
        { id: '1', type: 'contribution', amount: 50000, description: 'Cotisation Janvier - Mamadou', created_at: '2024-01-15' },
        { id: '2', type: 'loan', amount: 100000, description: 'Prêt à Ibrahim', created_at: '2024-01-10' }
      ]
    }
  }
})
</script>
