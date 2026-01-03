<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-3xl mx-auto space-y-6">
      <div class="flex items-center space-x-4 mb-6">
        <button @click="navigateTo('/associations')" class="p-2 rounded-full hover:bg-surface-hover transition-colors">
          <ArrowLeftIcon class="w-6 h-6 text-muted" />
        </button>
        <div>
          <h1 class="text-2xl font-bold text-base">Créer une association</h1>
          <p class="text-muted">Configurez votre nouvelle communauté financière</p>
        </div>
      </div>

      <div class="bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700 overflow-hidden">
      <!-- Progress Steps -->
      <div class="flex border-b border-gray-100 dark:border-gray-700">
        <div v-for="(step, index) in steps" :key="index" 
             class="flex-1 py-4 text-center border-b-2 transition-colors cursor-pointer"
             :class="[
               currentStep === index ? 'border-indigo-600 text-indigo-600 font-medium' : 
               currentStep > index ? 'border-green-500 text-green-500' : 'border-transparent text-gray-400'
             ]"
             @click="canNavigateTo(index) ? currentStep = index : null">
          <div class="flex items-center justify-center space-x-2">
            <div v-if="currentStep > index" class="w-6 h-6 rounded-full bg-green-100 text-green-600 flex items-center justify-center">
              <CheckIcon class="w-4 h-4" />
            </div>
            <span v-else class="w-6 h-6 rounded-full flex items-center justify-center text-xs border"
                  :class="currentStep === index ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-300 text-gray-500'">
              {{ index + 1 }}
            </span>
            <span>{{ step.label }}</span>
          </div>
        </div>
      </div>

      <div class="p-8">
        <form @submit.prevent="handleSubmit">
          <!-- Step 1: Basic Info -->
          <div v-if="currentStep === 0" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nom de l'association</label>
              <input v-model="form.name" type="text" placeholder="Ex: Tontine des Amis" 
                     class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
                     required />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Description</label>
              <textarea v-model="form.description" rows="3" placeholder="Quel est le but de cette association ?" 
                        class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Devise principale</label>
              <select v-model="form.currency" 
                      class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all">
                <option value="XOF">XOF - Franc CFA (BCEAO)</option>
                <option value="XAF">XAF - Franc CFA (BEAC)</option>
                <option value="GNF">GNF - Franc Guinéen</option>
                <option value="USD">USD - Dollar Américain</option>
                <option value="EUR">EUR - Euro</option>
              </select>
            </div>
          </div>

          <!-- Step 2: Type Selection -->
          <div v-if="currentStep === 1" class="space-y-6">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">Type d'association</label>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div v-for="type in associationTypes" :key="type.value"
                   @click="form.type = type.value"
                   class="cursor-pointer border-2 rounded-xl p-4 hover:border-indigo-400 transition-all relative"
                   :class="form.type === type.value ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-900/20' : 'border-gray-200 dark:border-gray-700'">
                <div v-if="form.type === type.value" class="absolute top-3 right-3 text-indigo-600">
                  <CheckCircleIcon class="w-6 h-6" />
                </div>
                <div class="w-10 h-10 rounded-lg bg-white dark:bg-gray-800 shadow-sm flex items-center justify-center mb-3">
                  <component :is="type.icon" class="w-6 h-6 text-indigo-600" />
                </div>
                <h3 class="font-bold text-gray-900 dark:text-white">{{ type.label }}</h3>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ type.description }}</p>
              </div>
            </div>
          </div>

          <!-- Step 3: Rules & Config -->
          <div v-if="currentStep === 2" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Configuration Règlements</h3>
            
            <div v-if="form.type === 'tontine'">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Montant de la cotisation (par tour)</label>
              <div class="relative rounded-md shadow-sm">
                <input v-model.number="form.rules.contribution_amount" type="number" 
                       class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
                       placeholder="0.00" />
                <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                  <span class="text-gray-500 sm:text-sm">{{ form.currency }}</span>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-1">Montant fixe que chaque membre doit payer à chaque période</p>
            </div>

            <div v-if="form.type === 'tontine' || form.type === 'savings'">
               <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 mt-4">Fréquence des cotisations</label>
               <select v-model="form.rules.frequency"
                       class="w-full px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all">
                 <option value="weekly">Hebdomadaire (Chaque semaine)</option>
                 <option value="monthly">Mensuelle (Chaque mois)</option>
                 <option value="biweekly">Bi-mensuelle (Toutes les 2 semaines)</option>
                 <option value="daily">Journalière</option>
               </select>
            </div>

            <!-- Generic rules for all types -->
             <div class="mt-4">
               <label class="flex items-center space-x-3">
                 <input v-model="form.rules.loans_enabled" type="checkbox" class="h-5 w-5 text-indigo-600 rounded focus:ring-indigo-500 border-gray-300 dark:border-gray-600" />
                 <span class="text-gray-900 dark:text-white">Autoriser les prêts internes</span>
               </label>
               <p class="text-xs text-gray-500 ml-8 mt-1">Permettre aux membres d'emprunter de l'argent de la caisse</p>
             </div>

             <div v-if="form.rules.loans_enabled" class="mt-4 pl-8 border-l-2 border-indigo-100 dark:border-indigo-900 space-y-4">
               <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Taux d'intérêt (%)</label>
                  <input v-model.number="form.rules.loan_interest_rate" type="number" step="0.1"
                        class="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700" />
               </div>
             </div>
          </div>

          <!-- Actions -->
          <div class="mt-8 flex justify-between pt-6 border-t border-gray-100 dark:border-gray-700">
            <button v-if="currentStep > 0" type="button" @click="currentStep--" 
                    class="px-6 py-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 font-medium transition-colors">
              Précédent
            </button>
            <div v-else></div> <!-- Spacer -->

            <button v-if="currentStep < steps.length - 1" type="button" @click="nextStep"
                    class="px-6 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-medium transition-colors shadow-lg shadow-indigo-500/30">
              Suivant
            </button>
            <button v-else type="submit" :disabled="loading"
                    class="px-6 py-2 rounded-lg bg-green-600 hover:bg-green-700 text-white font-medium transition-colors shadow-lg shadow-green-500/30 flex items-center">
              <span v-if="loading" class="animate-spin mr-2 h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
              Créer l'association
            </button>
          </div>
        </form>
      </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ArrowLeftIcon, 
  CheckIcon, 
  CheckCircleIcon,
  UserGroupIcon, 
  ArrowTrendingUpIcon, 
  BanknotesIcon, 
  ScaleIcon 
} from '@heroicons/vue/24/outline'

definePageMeta({
  layout: false,
  middleware: 'auth'
})

const { associationApi } = useApi()
const router = useRouter()
const loading = ref(false)
const currentStep = ref(0)

const steps = [
  { label: 'Informations' },
  { label: 'Type' },
  { label: 'Règles' }
]

const form = reactive({
  name: '',
  description: '',
  currency: 'XOF',
  type: 'tontine',
  rules: {
    contribution_amount: 0,
    frequency: 'monthly',
    loans_enabled: false,
    loan_interest_rate: 5
  }
})

const associationTypes = [
  { value: 'tontine', label: 'Tontine Rotative', description: 'Chaque membre cotise, et un bénéficiaire reçoit la totalité à tour de rôle.', icon: ArrowTrendingUpIcon },
  { value: 'savings', label: 'Groupe d\'Épargne', description: 'Cotisations libres ou fixes. L\'argent est accumulé dans une caisse commune.', icon: BanknotesIcon },
  { value: 'credit', label: 'Crédit Mutuel', description: 'Épargne orientée vers l\'octroi de crédits aux membres avec intérêts.', icon: ScaleIcon },
  { value: 'general', label: 'Association Générale', description: 'Gestion simple de trésorerie pour clubs, amicales, ou syndicats.', icon: UserGroupIcon }
]

const nextStep = () => {
  if (currentStep.value === 0 && (!form.name || !form.currency)) {
    // Basic validation
    return 
  }
  currentStep.value++
}

const canNavigateTo = (index: number) => {
  return index < currentStep.value
}

const handleSubmit = async () => {
  loading.value = true
  try {
    const response = await associationApi.create(form)
    // Success notification could be added here
    await router.push('/associations')
  } catch (err) {
    console.error('Failed to create association', err)
    // Show error
  } finally {
    loading.value = false
  }
}
</script>
