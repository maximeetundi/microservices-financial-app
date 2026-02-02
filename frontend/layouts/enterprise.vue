<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileMenuOpen" 
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      @click="isMobileMenuOpen = false"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-white dark:bg-slate-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out lg:transform-none flex flex-col"
      :class="isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo / Enterprise Name -->
      <div class="h-16 flex items-center px-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center gap-3 flex-1 min-w-0">
           <div v-if="enterprise?.logo" class="h-10 w-10 rounded-xl overflow-hidden flex-shrink-0">
             <img :src="enterprise.logo" class="h-full w-full object-cover" />
           </div>
           <div v-else class="h-10 w-10 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center text-white font-bold text-lg flex-shrink-0">
             {{ enterprise?.name?.charAt(0)?.toUpperCase() || 'E' }}
           </div>
           <div class="min-w-0 flex-1">
             <span class="font-bold text-gray-900 dark:text-white truncate block">
               {{ enterprise?.name || 'Chargement...' }}
             </span>
             <span class="text-xs text-gray-500 dark:text-gray-400">
               {{ userRoleLabel }}
             </span>
           </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
        <!-- Overview - visible to all -->
        <NuxtLink 
          :to="`/enterprise/${enterpriseId}`"
          class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
          :class="isExactActive ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
          @click="isMobileMenuOpen = false"
        >
          <Squares2X2Icon class="h-5 w-5 mr-3 transition-colors" :class="isExactActive ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
          Aperçu
        </NuxtLink>

        <!-- ============== ADMIN MENU ============== -->
        <template v-if="isAdmin">
          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Équipe</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/employees`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('employees') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <UsersIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('employees') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Employés
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/clients`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('clients') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <UserGroupIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('clients') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Abonnés
          </NuxtLink>

          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Finances</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/wallets`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('wallets') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <WalletIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('wallets') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Portefeuilles
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/payroll`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('payroll') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <BanknotesIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('payroll') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Paie
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/billing`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('billing') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <DocumentTextIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('billing') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Facturation
          </NuxtLink>

          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Configuration</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/services`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('services') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <FolderIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('services') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Services
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/security`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('security') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <ShieldCheckIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('security') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Sécurité
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/settings`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('settings') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <Cog6ToothIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('settings') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Paramètres
          </NuxtLink>
        </template>

        <!-- ============== EMPLOYEE MENU ============== -->
        <template v-else-if="isEmployee">
          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Mon Espace</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/my-profile`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('my-profile') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <UserIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('my-profile') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Mon Profil
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/my-salary`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('my-salary') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <BanknotesIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('my-salary') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Mes Salaires
          </NuxtLink>

          <!-- Show approvals if employee has manager role -->
          <NuxtLink 
            v-if="userEmployee?.role === 'MANAGER'"
            :to="`/enterprise/${enterpriseId}/approvals`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('approvals') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <ClipboardDocumentCheckIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('approvals') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Approbations
          </NuxtLink>
        </template>

        <!-- ============== CLIENT/SUBSCRIBER MENU ============== -->
        <template v-else-if="isClient">
          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Mes Abonnements</p>
          </div>

          <!-- Dynamic services list -->
          <NuxtLink 
            v-for="sub in mySubscriptions" 
            :key="sub.id"
            :to="`/enterprise/${enterpriseId}/subscription/${sub.id}`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive(`subscription/${sub.id}`) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <TagIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive(`subscription/${sub.id}`) ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            {{ sub.service_name || 'Service' }}
          </NuxtLink>

          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Paiements</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/my-invoices`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('my-invoices') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <DocumentTextIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('my-invoices') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Mes Factures
          </NuxtLink>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/payment-history`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="isActive('payment-history') ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700'"
            @click="isMobileMenuOpen = false"
          >
            <ClockIcon class="h-5 w-5 mr-3 transition-colors" :class="isActive('payment-history') ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 group-hover:text-gray-500'" />
            Historique
          </NuxtLink>
        </template>

        <!-- ============== VISITOR/UNKNOWN MENU ============== -->
        <template v-else>
          <div class="pt-4 pb-2">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4">Services</p>
          </div>

          <NuxtLink 
            :to="`/enterprise/${enterpriseId}/services-public`"
            class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors group text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-slate-700"
            @click="isMobileMenuOpen = false"
          >
            <FolderIcon class="h-5 w-5 mr-3 text-gray-400 group-hover:text-gray-500" />
            Voir les services
          </NuxtLink>
        </template>
      </nav>

      <!-- Bottom Actions -->
      <div class="p-4 border-t border-gray-200 dark:border-gray-700 space-y-2">
         <NuxtLink 
           to="/enterprise"
           class="flex items-center justify-center w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-slate-700 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors"
         >
           <ArrowLeftIcon class="h-4 w-4 mr-2" />
           {{ isAdmin ? 'Mes Entreprises' : 'Retour' }}
         </NuxtLink>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Mobile Header -->
      <header class="bg-white dark:bg-slate-800 shadow-sm lg:hidden h-16 flex items-center justify-between px-4 z-30 border-b border-gray-200 dark:border-gray-700">
        <button @click="isMobileMenuOpen = true" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
          <Bars3Icon class="h-6 w-6" />
        </button>
        <span class="font-bold text-gray-900 dark:text-white">{{ enterprise?.name }}</span>
        <div class="w-6"></div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-8">
        <div class="max-w-7xl mx-auto">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watchEffect, provide, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { enterpriseAPI } from '@/composables/useApi'
import { useAuthStore } from '@/stores/auth'
import { 
  Squares2X2Icon,
  UsersIcon,
  UserGroupIcon,
  UserIcon,
  WalletIcon,
  BanknotesIcon,
  DocumentTextIcon,
  FolderIcon,
  ShieldCheckIcon,
  Cog6ToothIcon,
  ArrowLeftIcon,
  Bars3Icon,
  TagIcon,
  ClockIcon,
  ClipboardDocumentCheckIcon
} from '@heroicons/vue/24/outline'

const route = useRoute()
const authStore = useAuthStore()
const enterpriseId = computed(() => route.params.id as string)
const isMobileMenuOpen = ref(false)
const enterprise = ref<any>(null)
const userEmployee = ref<any>(null)
const mySubscriptions = ref<any[]>([])
const isLoading = ref(true)

// Role detection
const isAdmin = computed(() => {
  if (!userEmployee.value) return false
  const role = userEmployee.value.role
  return role === 'OWNER' || role === 'ADMIN'
})

const isEmployee = computed(() => {
  if (!userEmployee.value) return false
  const role = userEmployee.value.role
  return role === 'EMPLOYEE' || role === 'MANAGER'
})

const isClient = computed(() => {
  // Client = has subscriptions but no employee record
  return !userEmployee.value && mySubscriptions.value.length > 0
})

const userRoleLabel = computed(() => {
  if (isAdmin.value) {
    return userEmployee.value?.role === 'OWNER' ? 'Propriétaire' : 'Administrateur'
  } else if (isEmployee.value) {
    return userEmployee.value?.role === 'MANAGER' ? 'Manager' : 'Employé'
  } else if (isClient.value) {
    return 'Abonné'
  }
  return 'Visiteur'
})

const isExactActive = computed(() => {
  return route.path === `/enterprise/${enterpriseId.value}`
})

const isActive = (section: string) => {
  return route.path.includes(`/enterprise/${enterpriseId.value}/${section}`)
}

// Fetch enterprise and user's role
watchEffect(async () => {
  if (!enterpriseId.value) return
  isLoading.value = true
  
  try {
    // Get enterprise data
    const { data } = await enterpriseAPI.get(enterpriseId.value)
    enterprise.value = data
    
    // Get current user's employee record (if any)
    try {
      const empRes = await enterpriseAPI.getMyEmployee(enterpriseId.value)
      userEmployee.value = empRes.data
    } catch {
      userEmployee.value = null
    }
    
    // Get user's subscriptions (if any)
    try {
      const subRes = await enterpriseAPI.getMySubscriptions(enterpriseId.value)
      mySubscriptions.value = subRes.data || []
    } catch {
      mySubscriptions.value = []
    }
  } catch (e) {
    console.error('Layout failed to load enterprise', e)
  } finally {
    isLoading.value = false
  }
})

// Provide to child pages
provide('enterprise', enterprise)
provide('enterpriseId', enterpriseId)
provide('userEmployee', userEmployee)
provide('userRole', computed(() => userEmployee.value?.role || 'VISITOR'))
provide('isAdmin', isAdmin)
provide('isEmployee', isEmployee)
provide('isClient', isClient)
</script>
