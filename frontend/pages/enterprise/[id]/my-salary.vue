<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Mes Salaires</h1>
      <p class="text-gray-500 dark:text-gray-400">Historique de vos paiements de salaire</p>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-emerald-500 border-t-transparent rounded-full mx-auto"></div>
    </div>

    <div v-else class="space-y-6">
      <!-- Current Salary Card -->
      <div v-if="employee?.salary_config" class="bg-gradient-to-r from-emerald-500 to-teal-600 rounded-xl p-6 text-white">
        <p class="text-sm opacity-80 mb-1">Salaire mensuel net</p>
        <p class="text-3xl font-bold">
          {{ formatAmount(employee.salary_config.base_salary) }} {{ employee.salary_config.currency }}
        </p>
        <p class="text-sm opacity-80 mt-2">Prochain paiement prévu le {{ formatDate(employee.salary_config.next_payment_date) }}</p>
      </div>

      <!-- Salary History -->
      <div class="bg-white dark:bg-slate-800 rounded-xl shadow-sm overflow-hidden">
        <div class="p-4 border-b border-gray-200 dark:border-gray-700">
          <h3 class="font-semibold text-gray-900 dark:text-white">Historique des paiements</h3>
        </div>
        
        <div v-if="salaryHistory.length === 0" class="p-8 text-center text-gray-500 dark:text-gray-400">
          <BanknotesIcon class="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>Aucun paiement de salaire enregistré</p>
        </div>

        <div v-else class="divide-y divide-gray-200 dark:divide-gray-700">
          <div v-for="payment in salaryHistory" :key="payment.id" class="p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-slate-700/50">
            <div class="flex items-center gap-4">
              <div class="h-10 w-10 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
                <CheckCircleIcon class="h-5 w-5 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <p class="font-medium text-gray-900 dark:text-white">
                  Salaire {{ formatMonth(payment.period_month, payment.period_year) }}
                </p>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  Payé le {{ formatDate(payment.executed_at) }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold text-gray-900 dark:text-white">
                {{ formatAmount(payment.net_pay) }} {{ payment.currency || 'XOF' }}
              </p>
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="payment.status === 'SUCCESS' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-800'">
                {{ payment.status === 'SUCCESS' ? 'Payé' : payment.status }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, inject, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { enterpriseAPI } from '@/composables/useApi'
import { BanknotesIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'

definePageMeta({
  layout: 'enterprise',
})

const route = useRoute()
const enterpriseId = computed(() => route.params.id as string)
const userEmployee = inject<any>('userEmployee')
const employee = computed(() => userEmployee?.value)

const loading = ref(true)
const salaryHistory = ref<any[]>([])

const fetchSalaryHistory = async () => {
  loading.value = true
  try {
    const { data } = await enterpriseAPI.getMySalary(enterpriseId.value)
    salaryHistory.value = data?.payments || []
  } catch (e) {
    console.error('Failed to load salary history', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchSalaryHistory)

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('fr-FR', { dateStyle: 'long' })
}

const formatMonth = (month: number, year: number) => {
  const months = ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Septembre', 'Octobre', 'Novembre', 'Décembre']
  return `${months[month - 1]} ${year}`
}

const formatAmount = (amount: number) => {
  if (!amount) return '0'
  return amount.toLocaleString('fr-FR')
}
</script>
