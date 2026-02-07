<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-2xl mx-auto space-y-6">
      <div class="flex items-center space-x-4">
        <button @click="$router.back()" class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Créer une entreprise</h1>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 border border-gray-100 dark:border-gray-700">
        <form @submit.prevent="createEnterprise" class="space-y-6">
          
          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nom de l'entreprise</label>
            <input v-model="form.name" type="text" required 
                   class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
                   placeholder="Ex: Mon Entreprise" />
          </div>

          <!-- Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Type d'activité</label>
            <select v-model="form.type" required
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all">
              <option value="SERVICE">Service Général</option>
              <option value="SCHOOL">École / Éducation</option>
              <option value="TRANSPORT">Transport</option>
              <option value="UTILITY">Eau / Électricité / Gaz</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nombre d'employés</label>
            <select v-model="form.employee_count_range"
                    class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all">
              <option value="1-10">1 - 10 employés</option>
              <option value="11-50">11 - 50 employés</option>
              <option value="51-200">51 - 200 employés</option>
              <option value="201-500">201 - 500 employés</option>
              <option value="500+">Plus de 500 employés</option>
            </select>
          </div>

          <div class="pt-4">
             <button type="submit" 
                     :disabled="loading"
                     class="w-full py-3 px-6 bg-primary-600 hover:bg-primary-700 text-white font-semibold rounded-lg shadow-md hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center">
                <span v-if="loading" class="animate-spin mr-2 h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
                {{ loading ? 'Création...' : 'Créer' }}
             </button>
          </div>

        </form>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { enterpriseAPI } from '@/composables/useApi'

definePageMeta({ middleware: ['auth'], layout: 'dashboard' })

const router = useRouter()
const toast = useToast()

const loading = ref(false)
const form = reactive({
  name: '',
  type: 'SERVICE',
  employee_count_range: '1-10'
})

const createEnterprise = async () => {
  loading.value = true
  try {
    await enterpriseAPI.create(form)
    toast.success('Entreprise créée avec succès !')
    router.push('/enterprise')
  } catch (err) {
    console.error(err)
    toast.error(err.response?.data?.error || 'Échec de la création')
  } finally {
    loading.value = false
  }
}
</script>
