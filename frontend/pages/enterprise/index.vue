<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <!-- Build Header -->
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Enterprise Portal</h1>
          <p class="text-gray-500 dark:text-gray-400">Manage your business, employees, and billing</p>
        </div>
        <button @click="openCreateModal" class="px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors">
          Create New Enterprise
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
             <span>{{ ent.employees_count || 0 }} Employees</span>
             <span>Status: Active</span>
          </div>
        </div>
      </div>

      <!-- Details View (Once selected) -->
      <div v-else class="space-y-6">
         <button @click="currentEnterprise = null" class="text-sm text-gray-500 hover:text-gray-700 underline mb-4">
            &larr; Back to List
         </button>

         <!-- Tabs -->
         <div class="flex space-x-1 bg-gray-100 dark:bg-gray-800 p-1 rounded-lg w-fit">
            <button v-for="tab in tabs" :key="tab" @click="currentTab = tab"
               :class="['px-4 py-2 rounded-md text-sm font-medium transition-all', 
                        currentTab === tab ? 'bg-white dark:bg-gray-700 shadow text-gray-900 dark:text-white' : 'text-gray-500 hover:text-gray-700']">
              {{ tab }}
            </button>
         </div>

         <!-- Employee Tab -->
         <div v-if="currentTab === 'Employees'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <div class="flex justify-between mb-4">
               <h3 class="font-semibold text-lg dark:text-white">Employees</h3>
               <button class="px-3 py-1 bg-green-600 text-white rounded-md text-sm">Invite Employee</button>
            </div>
            <!-- Mock Table -->
             <table class="w-full text-left text-sm text-gray-500">
               <thead class="bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                  <tr>
                     <th class="px-4 py-3">Name</th>
                     <th class="px-4 py-3">Role</th>
                     <th class="px-4 py-3">Status</th>
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
                     <td colspan="4" class="px-4 py-8 text-center text-gray-400">No employees found.</td>
                  </tr>
               </tbody>
            </table>
         </div>

         <!-- Billing Tab -->
         <div v-if="currentTab === 'Billing'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <h3 class="font-semibold text-lg dark:text-white mb-4">Bulk Invoicing</h3>
            <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors">
               <p class="text-gray-500">Drag and drop your CSV/Excel file here to generate invoices</p>
               <button class="mt-4 px-4 py-2 border border-gray-300 rounded-md text-sm">Upload File</button>
            </div>
         </div>
      </div>

    </div>
  </NuxtLayout>
</template>

<script setup>
import { enterpriseAPI } from '@/composables/useApi'

const tabs = ['Overview', 'Employees', 'Payroll', 'Billing', 'Settings']
const currentTab = ref('Overview')
const currentEnterprise = ref(null)

const enterprises = ref([])
const employees = ref([])
const isLoading = ref(true)

const fetchEnterprises = async () => {
   try {
      console.log('Fetching enterprises...')
      const { data } = await enterpriseAPI.list()
      enterprises.value = data
   } catch (error) {
      console.error('Failed to fetch enterprises', error)
      // Fallback empty or error state
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
  // Logic to open modal
}
</script>
