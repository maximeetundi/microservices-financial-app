<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-2xl mx-auto space-y-6">
      <div class="flex items-center space-x-4">
        <button @click="$router.back()" class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Create New Enterprise</h1>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-100 dark:border-gray-700">
        <form @submit.prevent="createEnterprise" class="space-y-6">
          
          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Enterprise Name</label>
            <input v-model="form.name" type="text" required 
                   class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
                   placeholder="e.g. Acme Corp" />
          </div>

          <!-- Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Business Type</label>
            <select v-model="form.type" required
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all">
              <option value="GENERIC">General Service</option>
              <option value="SCHOOL">School / Education</option>
              <option value="TRANSPORT">Public Transport</option>
              <option value="UTILITY">Utility (Water, Energy)</option>
            </select>
          </div>

          <!-- Modules -->
          <div class="space-y-3">
             <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Active Modules</label>
             <div class="flex items-center space-x-3">
                <input type="checkbox" v-model="form.has_payroll" id="payroll" class="rounded text-primary-600 focus:ring-primary-500 w-5 h-5 bg-gray-100 border-gray-300">
                <label for="payroll" class="text-gray-700 dark:text-gray-300">Payroll Management</label>
             </div>
             <div class="flex items-center space-x-3">
                <input type="checkbox" v-model="form.has_invoicing" id="invoicing" class="rounded text-primary-600 focus:ring-primary-500 w-5 h-5 bg-gray-100 border-gray-300">
                <label for="invoicing" class="text-gray-700 dark:text-gray-300">Recurring Billing / Invoicing</label>
             </div>
          </div>

          <div class="pt-4">
             <button type="submit" 
                     :disabled="loading"
                     class="w-full py-3 px-6 bg-primary-600 hover:bg-primary-700 text-white font-semibold rounded-lg shadow-md hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center">
                <span v-if="loading" class="animate-spin mr-2 h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
                {{ loading ? 'Creating...' : 'Launch Enterprise' }}
             </button>
          </div>

        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
const { enterpriseApi } = useApi()
const router = useRouter()
const toast = useToast()

const loading = ref(false)
const form = reactive({
  name: '',
  type: 'GENERIC',
  has_payroll: true,
  has_invoicing: true
})

const createEnterprise = async () => {
  loading.value = true
  try {
    const res = await enterpriseApi.create(form)
    toast.success('Enterprise created successfully!')
    router.push('/enterprise')
  } catch (err) {
    console.error(err)
    toast.error(err.response?.data?.error || 'Failed to create enterprise')
  } finally {
    loading.value = false
  }
}
</script>
