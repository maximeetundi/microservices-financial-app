<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Mon Profil</h1>
      <p class="text-gray-500 dark:text-gray-400">Vos informations en tant qu'employé</p>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-emerald-500 border-t-transparent rounded-full mx-auto"></div>
    </div>

    <div v-else-if="employee" class="space-y-6">
      <!-- Profile Card -->
      <div class="bg-white dark:bg-slate-800 rounded-xl shadow-sm p-6">
        <div class="flex items-center gap-4 mb-6">
          <div class="h-16 w-16 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center text-white text-2xl font-bold">
            {{ employee.first_name?.charAt(0) }}{{ employee.last_name?.charAt(0) }}
          </div>
          <div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ employee.first_name }} {{ employee.last_name }}
            </h2>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="roleClass">
              {{ roleLabel }}
            </span>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Email</label>
            <p class="text-gray-900 dark:text-white">{{ employee.email || '-' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Téléphone</label>
            <p class="text-gray-900 dark:text-white">{{ employee.phone || '-' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Poste</label>
            <p class="text-gray-900 dark:text-white">{{ employee.position || 'Non défini' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Département</label>
            <p class="text-gray-900 dark:text-white">{{ employee.department || 'Non défini' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Date d'embauche</label>
            <p class="text-gray-900 dark:text-white">{{ formatDate(employee.hire_date) }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Statut</label>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="employee.status === 'ACTIVE' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-800'">
              {{ employee.status === 'ACTIVE' ? 'Actif' : employee.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Salary Info (if available) -->
      <div v-if="employee.salary_config" class="bg-white dark:bg-slate-800 rounded-xl shadow-sm p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Configuration Salariale</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 dark:bg-slate-700 rounded-lg p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Salaire de base</p>
            <p class="text-xl font-bold text-gray-900 dark:text-white">
              {{ formatAmount(employee.salary_config.base_salary) }} {{ employee.salary_config.currency }}
            </p>
          </div>
          <div class="bg-gray-50 dark:bg-slate-700 rounded-lg p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Type de paiement</p>
            <p class="text-xl font-bold text-gray-900 dark:text-white">
              {{ paymentTypeLabel }}
            </p>
          </div>
          <div class="bg-gray-50 dark:bg-slate-700 rounded-lg p-4">
            <p class="text-sm text-gray-500 dark:text-gray-400">Prochain paiement</p>
            <p class="text-xl font-bold text-gray-900 dark:text-white">
              {{ formatDate(employee.salary_config.next_payment_date) }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-12 text-gray-500 dark:text-gray-400">
      Impossible de charger votre profil
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted } from 'vue'

definePageMeta({
  layout: 'enterprise',
})

const userEmployee = inject<any>('userEmployee')
const employee = computed(() => userEmployee?.value)
const loading = ref(false)

const roleClass = computed(() => {
  const role = employee.value?.role
  if (role === 'OWNER') return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400'
  if (role === 'ADMIN') return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
  if (role === 'MANAGER') return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400'
  return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
})

const roleLabel = computed(() => {
  const map: Record<string, string> = {
    'OWNER': 'Propriétaire',
    'ADMIN': 'Administrateur',
    'MANAGER': 'Manager',
    'EMPLOYEE': 'Employé'
  }
  return map[employee.value?.role] || employee.value?.role
})

const paymentTypeLabel = computed(() => {
  const type = employee.value?.salary_config?.payment_type
  const map: Record<string, string> = {
    'MONTHLY': 'Mensuel',
    'BIWEEKLY': 'Bi-hebdomadaire',
    'WEEKLY': 'Hebdomadaire',
  }
  return map[type] || type || 'Non défini'
})

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('fr-FR', { dateStyle: 'long' })
}

const formatAmount = (amount: number) => {
  if (!amount) return '0'
  return amount.toLocaleString('fr-FR')
}
</script>
