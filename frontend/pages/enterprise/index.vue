<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <!-- Build Header -->
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Portail Entreprise</h1>
          <p class="text-gray-500 dark:text-gray-400">Gérez votre entreprise, vos employés et la facturation</p>
        </div>
        <button @click="openCreateModal" class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors">
          Créer une nouvelle entreprise
        </button>
      </div>

      <!-- Enterprise List (Selection) -->
      <div v-if="!currentEnterprise" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Mock list for MVP or fetch using user ID -->
        <div v-for="ent in enterprises" :key="ent.id" @click="selectEnterprise(ent)" 
             class="cursor-pointer p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 dark:border-gray-700">
          <div class="flex items-center space-x-4 mb-4">
             <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-xl">
               {{ ent.name.charAt(0) }}
             </div>
             <div>
               <h3 class="font-semibold text-lg dark:text-white">{{ ent.name }}</h3>
               <span class="text-xs px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded-full text-gray-600 dark:text-gray-300">{{ ent.type }}</span>
             </div>
          </div>
          <div class="flex justify-between text-sm text-gray-500">
             <span>{{ ent.employees_count || 0 }} Employés</span>
             <span>Statut: Actif</span>
          </div>
        </div>
        
        <!-- Empty State -->
        <div v-if="enterprises.length === 0 && !isLoading" class="col-span-full text-center py-12">
            <p class="text-gray-500">Aucune entreprise trouvée.</p>
        </div>
      </div>

      <!-- Details View (Once selected) -->
      <div v-else class="space-y-6">
         <button @click="currentEnterprise = null" class="text-sm text-gray-500 hover:text-gray-700 underline mb-4">
            &larr; Retour à la liste
         </button>

         <!-- Tabs -->
         <div class="flex space-x-1 bg-gray-100 dark:bg-gray-800 p-1 rounded-lg w-fit">
            <button v-for="tab in tabs" :key="tab" @click="currentTab = tab"
               :class="['px-4 py-2 rounded-md text-sm font-medium transition-all', 
                        currentTab === tab ? 'bg-white dark:bg-gray-700 shadow text-gray-900 dark:text-white' : 'text-gray-500 hover:text-gray-700']">
              {{ tabLabels[tab] }}
            </button>
         </div>

         <!-- Employee Tab -->
         <div v-if="currentTab === 'Employees'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <div class="flex justify-between mb-4">
               <h3 class="font-semibold text-lg dark:text-white">Employés</h3>
               <button class="px-3 py-1 bg-green-600 text-white rounded-md text-sm">Inviter un employé</button>
            </div>
            <!-- Mock Table -->
             <table class="w-full text-left text-sm text-gray-500">
               <thead class="bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                  <tr>
                     <th class="px-4 py-3">Nom</th>
                     <th class="px-4 py-3">Rôle</th>
                     <th class="px-4 py-3">Statut</th>
                     <th class="px-4 py-3">Actions</th>
                  </tr>
               </thead>
               <tbody>
                  <tr v-for="emp in employees" :key="emp.id" class="border-b dark:border-gray-700">
                     <td class="px-4 py-3">{{ emp.first_name }} {{ emp.last_name }}</td>
                     <td class="px-4 py-3">{{ emp.profession }}</td>
                     <td class="px-4 py-3">
                        <span :class="{'text-green-600': emp.status === 'ACTIVE', 'text-yellow-600': emp.status === 'PENDING_INVITE'}">
                           {{ emp.status }}
                        </span>
                     </td>
                     <td class="px-4 py-3">...</td>
                  </tr>
                  <tr v-if="employees.length === 0">
                     <td colspan="4" class="px-4 py-8 text-center text-gray-400">Aucun employé trouvé.</td>
                  </tr>
               </tbody>
            </table>
         </div>

         <!-- Billing Tab -->
         <div v-if="currentTab === 'Billing'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <h3 class="font-semibold text-lg dark:text-white mb-4">Facturation en masse</h3>
            <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors">
               <p class="text-gray-500">Glissez-déposez votre fichier CSV/Excel ici pour générer des factures</p>
               <button class="mt-4 px-4 py-2 border border-gray-300 rounded-md text-sm">Télécharger un fichier</button>
            </div>
         </div>
      </div>
      
      <!-- Create Enterprise Modal -->
      <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-xl w-full max-w-lg p-6 space-y-4">
            <h2 class="text-xl font-bold dark:text-white">Créer une nouvelle entreprise</h2>
            
            <form @submit.prevent="handleCreateEnterprise" class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Nom de l'entreprise</label>
                    <input v-model="newEnterprise.name" type="text" required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Numéro d'enregistrement (NIF/RCCM)</label>
                    <input v-model="newEnterprise.registration_number" type="text" required class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                </div>
                
                <div>
                   <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Type</label>
                   <select v-model="newEnterprise.type" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                       <option value="SME">PME (Petite/Moyenne Entreprise)</option>
                       <option value="LARGE">Grande Entreprise</option>
                       <option value="STARTUP">Startup</option>
                   </select>
                </div>

                <div class="flex justify-end gap-3 mt-6">
                    <button type="button" @click="showCreateModal = false" class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                        Annuler
                    </button>
                    <button type="submit" :disabled="isCreating" class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50">
                        {{ isCreating ? 'Création...' : 'Créer' }}
                    </button>
                </div>
            </form>
        </div>
      </div>

    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue' // Explicit imports
import { enterpriseAPI } from '@/composables/useApi'

const tabs = ['Overview', 'Employees', 'Payroll', 'Billing', 'Settings']
const tabLabels = {
    'Overview': 'Aperçu',
    'Employees': 'Employés',
    'Payroll': 'Paie',
    'Billing': 'Facturation',
    'Settings': 'Paramètres'
}

const currentTab = ref('Overview')
const currentEnterprise = ref(null)

const enterprises = ref([])
const employees = ref([])
const isLoading = ref(true)

// Create Modal State
const showCreateModal = ref(false)
const isCreating = ref(false)
const newEnterprise = ref({
    name: '',
    registration_number: '',
    type: 'SME'
})

const fetchEnterprises = async () => {
   try {
      console.log('Fetching enterprises...')
      const { data } = await enterpriseAPI.list()
      enterprises.value = data
   } catch (error) {
      console.error('Failed to fetch enterprises', error)
      enterprises.value = []
   } finally {
      isLoading.value = false
   }
}

onMounted(() => {
   fetchEnterprises()
})

const selectEnterprise = (ent) => {
  currentEnterprise.value = ent
}

const fetchEmployees = async () => {
   if (!currentEnterprise.value) return
   try {
      const { data } = await enterpriseAPI.listEmployees(currentEnterprise.value.id)
      employees.value = data
   } catch (error) {
      console.error('Failed to fetch employees', error)
      employees.value = []
   }
}

watch(currentTab, (newTab) => {
   if (newTab === 'Employees') {
      fetchEmployees()
   }
})

// Initial fetch if already on tab (e.g. deep link or reload)
watch(currentEnterprise, (newEnt) => {
   if (newEnt && currentTab.value === 'Employees') {
      fetchEmployees()
   }
})

const openCreateModal = () => {
  showCreateModal.value = true
  newEnterprise.value = { name: '', registration_number: '', type: 'SME' }
}

const handleCreateEnterprise = async () => {
    isCreating.value = true
    try {
        await enterpriseAPI.create(newEnterprise.value)
        showCreateModal.value = false
        // Refresh list
        await fetchEnterprises()
    } catch (error) {
        console.error('Failed to create', error)
        alert('Erreur lors de la création')
    } finally {
        isCreating.value = false
    }
}
</script>
