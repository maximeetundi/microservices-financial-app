<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6 min-h-screen">
      
      <!-- Top Navigation / Breadcrumbs -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <nav v-if="currentEnterprise" class="flex items-center text-sm text-gray-500 mb-2">
            <button @click="currentEnterprise = null" class="hover:text-primary-600 transition-colors flex items-center gap-1">
              <Squares2X2Icon class="w-4 h-4" /> Entreprises
            </button>
            <ChevronRightIcon class="w-4 h-4 mx-2" />
            <span class="font-medium text-gray-900 dark:text-white">{{ currentEnterprise.name }}</span>
          </nav>
          
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white tracking-tight">
            {{ currentEnterprise ? 'Tableau de Bord' : 'Portail Entreprise' }}
          </h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            {{ currentEnterprise ? 'Gérez les activités de ' + currentEnterprise.name : 'Pilotez l\'ensemble de vos structures professionnelles' }}
          </p>
        </div>

        <div class="flex gap-3">
          <button v-if="currentEnterprise" @click="showQRModal = true" 
            class="px-5 py-2.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-200 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 transition-all shadow-sm hover:shadow flex items-center gap-2 font-medium">
            <QrCodeIcon class="w-5 h-5 text-primary-600" />
            <span>Codes QR</span>
          </button>
          
          <button v-if="!currentEnterprise" @click="showCreateModal = true" 
            class="px-5 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white rounded-xl transition-all shadow-md hover:shadow-lg flex items-center gap-2 font-medium transform hover:-translate-y-0.5">
            <PlusIcon class="w-5 h-5" />
            <span>Nouvelle Entreprise</span>
          </button>
        </div>
      </div>

      <!-- Enterprise Selection List -->
      <div v-if="!currentEnterprise" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="ent in enterprises" :key="ent.id" @click="selectEnterprise(ent)" 
          class="group cursor-pointer relative bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 dark:border-gray-700">
          
          <div class="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
            <ArrowRightCircleIcon class="w-6 h-6 text-primary-500" />
          </div>

          <div class="flex items-center gap-4 mb-6">
            <div v-if="ent.logo" class="w-14 h-14 rounded-xl border border-gray-200 p-1 overflow-hidden">
              <img :src="ent.logo" class="w-full h-full object-cover rounded-lg" />
            </div>
            <div v-else class="w-14 h-14 rounded-xl bg-gradient-to-br from-primary-100 to-blue-50 dark:from-primary-900/30 dark:to-blue-900/10 flex items-center justify-center text-primary-600 dark:text-primary-400 font-bold text-xl ring-1 ring-primary-100 dark:ring-primary-800">
              {{ ent.name.charAt(0) }}
            </div>
            <div>
              <h3 class="font-bold text-lg text-gray-900 dark:text-white group-hover:text-primary-600 transition-colors">{{ ent.name }}</h3>
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 mt-1">
                {{ formatEnterpriseType(ent.type) }}
              </span>
            </div>
          </div>
          
          <div class="flex justify-between items-center text-sm text-gray-500 border-t dark:border-gray-700 pt-4">
            <span class="flex items-center gap-1"><UsersIcon class="w-4 h-4" /> {{ ent.employees_count || 0 }} Membres</span>
            <span class="flex items-center gap-1 text-green-600 dark:text-green-400"><CheckCircleIcon class="w-4 h-4" /> Actif</span>
          </div>
        </div>
        
        <div v-if="enterprises.length === 0 && !isLoading" class="col-span-full py-16 text-center bg-white dark:bg-gray-800 rounded-2xl border-2 border-dashed border-gray-200 dark:border-gray-700">
          <BuildingOffice2Icon class="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">Aucune entreprise</h3>
          <p class="text-gray-500 mt-1">Commencez par créer votre première structure.</p>
          <button @click="showCreateModal = true" class="mt-4 px-5 py-2.5 bg-primary-600 hover:bg-primary-700 text-white rounded-xl font-medium transition-colors">
            Créer une entreprise
          </button>
        </div>
      </div>

      <!-- Enterprise Dashboard View -->
      <div v-else class="space-y-6 animate-fade-in">
        
        <!-- Navigation Tabs -->
        <div class="bg-white dark:bg-gray-800 p-1.5 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 inline-flex overflow-x-auto max-w-full">
          <button v-for="tab in tabs" :key="tab" @click="currentTab = tab"
            :class="['px-5 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap flex items-center gap-2', 
              currentTab === tab 
                ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 shadow-sm ring-1 ring-primary-100 dark:ring-primary-800' 
                : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/50']">
            <component :is="getTabIcon(tab)" class="w-4 h-4" />
            {{ tabLabels[tab] }}
          </button>
        </div>

        <!-- Tab Content -->
        <EnterpriseOverview 
          v-if="currentTab === 'Overview'" 
          :enterprise="currentEnterprise"
          :stats="enterpriseStats"
          @navigate="currentTab = $event"
          @show-qr="showQRModal = true" />

        <EmployeesTab 
          v-if="currentTab === 'Employees'" 
          :enterprise="currentEnterprise" />

        <ClientsTab 
          v-if="currentTab === 'Clients'" 
          :enterprise="currentEnterprise" />

        <ServicesTab 
          v-if="currentTab === 'Services'" 
          :enterprise="currentEnterprise"
          @update="handleEnterpriseUpdate" />

        <PayrollTab 
          v-if="currentTab === 'Payroll'" 
          :enterprise-id="currentEnterprise.id" />

        <BillingTab 
          v-if="currentTab === 'Billing'" 
          :enterprise="currentEnterprise" />

        <WalletsTab 
          v-if="currentTab === 'Wallets'" 
          :enterprise="currentEnterprise"
          @update="handleEnterpriseUpdate" />

        <SecurityTab 
          v-if="currentTab === 'Security'" 
          :enterprise="currentEnterprise"
          @update="handleEnterpriseUpdate" />

        <OrganisationTab 
          v-if="currentTab === 'Organisation'" 
          :enterprise="currentEnterprise"
          :employees="enterpriseEmployees"
          @update="handleEnterpriseUpdate"
          @save="saveSettings" />

        <EnterpriseSettings 
          v-if="currentTab === 'Settings'" 
          v-model="currentEnterprise"
          @upload-logo="handleLogoUpload" />

        <!-- Save Button (for Settings) -->
        <div v-if="currentTab === 'Settings'" class="flex justify-end pt-4 border-t dark:border-gray-700">
          <button @click="saveSettings" :disabled="isSaving" 
            class="px-6 py-2.5 bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-xl font-medium hover:from-primary-700 hover:to-primary-800 disabled:opacity-50 flex items-center gap-2 shadow-lg shadow-primary-500/25 transition-all">
            <span v-if="isSaving" class="animate-spin">⟳</span>
            {{ isSaving ? 'Enregistrement...' : 'Enregistrer les modifications' }}
          </button>
        </div>
      </div>

      <!-- QR Code Modal -->
      <QRCodeModal 
        :is-open="showQRModal" 
        :enterprise="currentEnterprise"
        @close="showQRModal = false" />

      <!-- Create Enterprise Modal -->
      <CreateEnterpriseModal 
        v-if="showCreateModal"
        @close="showCreateModal = false"
        @created="handleEnterpriseCreated" />
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { enterpriseAPI } from '@/composables/useApi'
import { 
  Squares2X2Icon, UsersIcon, UserGroupIcon, BanknotesIcon, 
  DocumentTextIcon, Cog6ToothIcon, ChevronRightIcon, QrCodeIcon, 
  PlusIcon, ArrowRightCircleIcon, CheckCircleIcon, BuildingOffice2Icon,
  FolderIcon, WalletIcon, ShieldCheckIcon
} from '@heroicons/vue/24/outline'

// Components
import EnterpriseOverview from '@/components/enterprise/EnterpriseOverview.vue'
import EmployeesTab from '@/components/enterprise/EmployeesTab.vue'
import EnterpriseSettings from '@/components/enterprise/EnterpriseSettings.vue'
import QRCodeModal from '@/components/enterprise/QRCodeModal.vue'

// Lazy loaded components (create stubs for now)
const ClientsTab = defineAsyncComponent(() => import('@/components/enterprise/ClientsTab.vue'))
const ServicesTab = defineAsyncComponent(() => import('@/components/enterprise/ServicesTab.vue'))
const PayrollTab = defineAsyncComponent(() => import('@/components/enterprise/PayrollTab.vue'))
const BillingTab = defineAsyncComponent(() => import('@/components/enterprise/BillingTab.vue'))
const WalletsTab = defineAsyncComponent(() => import('@/components/enterprise/WalletsTab.vue'))
const SecurityTab = defineAsyncComponent(() => import('@/components/enterprise/SecurityTab.vue'))
const OrganisationTab = defineAsyncComponent(() => import('@/components/enterprise/OrganisationTab.vue'))
const CreateEnterpriseModal = defineAsyncComponent(() => import('@/components/enterprise/CreateEnterpriseModal.vue'))

// State
const enterprises = ref([])
const currentEnterprise = ref(null)
const currentTab = ref('Overview')
const isLoading = ref(true)
const isSaving = ref(false)
const showQRModal = ref(false)
const showCreateModal = ref(false)
const employeesCount = ref(0)
const clientsCount = ref(0)
const myEmployee = ref(null) // Current user's employee record for permission checks
const enterpriseEmployees = ref([]) // All employees list for OrganisationTab

// Role-based tab access - Wallets and Services require admin permission
const adminOnlyTabs = ['Security', 'Organisation', 'Settings', 'Payroll', 'Wallets', 'Services']
const allTabs = ['Overview', 'Employees', 'Clients', 'Services', 'Wallets', 'Payroll', 'Billing', 'Security', 'Organisation', 'Settings']

// Visible tabs based on role
const tabs = computed(() => {
  if (!myEmployee.value) return allTabs // Default to all for owner
  const role = myEmployee.value.role
  if (role === 'OWNER' || role === 'ADMIN') {
    return allTabs
  }
  // STANDARD employees: limited access
  return allTabs.filter(tab => !adminOnlyTabs.includes(tab))
})

const tabLabels = {
  'Overview': 'Aperçu',
  'Employees': 'Employés',
  'Clients': 'Abonnés',
  'Services': 'Services',
  'Wallets': 'Portefeuilles',
  'Payroll': 'Paie',
  'Billing': 'Facturation',
  'Security': 'Sécurité',
  'Organisation': 'Organisation',
  'Settings': 'Paramètres'
}

const getTabIcon = (tab) => {
  switch (tab) {
    case 'Overview': return Squares2X2Icon
    case 'Employees': return UsersIcon
    case 'Clients': return UserGroupIcon
    case 'Services': return FolderIcon
    case 'Wallets': return WalletIcon
    case 'Payroll': return BanknotesIcon
    case 'Billing': return DocumentTextIcon
    case 'Security': return ShieldCheckIcon
    case 'Settings': return Cog6ToothIcon
    default: return Squares2X2Icon
  }
}

// Computed
const enterpriseStats = computed(() => {
  if (!currentEnterprise.value) return { employees: 0, clients: 0, services: 0, revenue: 0 }
  
  const totalServices = (currentEnterprise.value.service_groups || [])
    .reduce((sum, g) => sum + (g.services?.length || 0), 0)
  
  return {
    employees: employeesCount.value,
    clients: clientsCount.value,
    services: totalServices,
    revenue: 0 // TODO: Fetch from API
  }
})

// Methods
const formatEnterpriseType = (type) => {
  const map = {
    'SERVICE': 'Service',
    'SCHOOL': 'École',
    'TRANSPORT': 'Transport',
    'UTILITY': 'Utilitaire'
  }
  return map[type] || type
}

const fetchEnterprises = async () => {
  try {
    isLoading.value = true
    const { data } = await enterpriseAPI.list()
    enterprises.value = data || []
  } catch (error) {
    console.error('Failed to fetch enterprises', error)
    enterprises.value = []
  } finally {
    isLoading.value = false
  }
}

const selectEnterprise = async (ent) => {
  // Ensure nested objects exist
  if (!ent.settings) ent.settings = { payroll_date: 25, auto_pay_salaries: false }
  if (ent.type === 'SCHOOL' && !ent.school_config) ent.school_config = { classes: [] }
  if (ent.type === 'TRANSPORT' && !ent.transport_config) ent.transport_config = { routes: [], zones: [] }
  if (!ent.service_groups) ent.service_groups = []
  
  currentEnterprise.value = JSON.parse(JSON.stringify(ent))
  
  // Fetch current user's employee record for role-based access
  try {
    const { data: empData } = await enterpriseAPI.getMyEmployee(ent.id || ent._id)
    myEmployee.value = empData
  } catch (e) {
    console.error('Failed to fetch my employee record:', e)
    myEmployee.value = null // Default to full access if can't fetch (owner might not have employee record)
  }
  
  // Fetch employees list (for count and OrganisationTab)
  try {
    const { data: employees } = await enterpriseAPI.listEmployees(ent.id || ent._id)
    enterpriseEmployees.value = employees || []
    employeesCount.value = employees?.length || 0
  } catch (e) {
    console.error('Failed to fetch employees:', e)
    enterpriseEmployees.value = []
    employeesCount.value = 0
  }
  
  // Fetch clients count (clients are subscriptions in the backend)
  try {
    const { data: clients } = await enterpriseAPI.getSubscriptions(ent.id || ent._id)
    clientsCount.value = clients?.length || 0
  } catch (e) {
    console.error('Failed to fetch clients count:', e)
    clientsCount.value = 0
  }
}

const handleLogoUpload = async (file) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const { data } = await enterpriseAPI.uploadLogo(formData)
    currentEnterprise.value.logo = data.url
  } catch (e) {
    console.error('Failed to upload logo', e)
    alert('Erreur lors de l\'upload du logo')
  }
}

const saveSettings = async () => {
  if (!currentEnterprise.value) return
  isSaving.value = true
  try {
    await enterpriseAPI.update(currentEnterprise.value.id, currentEnterprise.value)
    await fetchEnterprises()
    alert('Sauvegardé avec succès')
  } catch (error) {
    console.error('Failed to save settings', error)
    alert('Erreur lors de la sauvegarde')
  } finally {
    isSaving.value = false
  }
}

const handleEnterpriseCreated = () => {
  showCreateModal.value = false
  fetchEnterprises()
}

const handleEnterpriseUpdate = (updated) => {
  currentEnterprise.value = updated
  fetchEnterprises()
}

// Lifecycle
onMounted(fetchEnterprises)
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
