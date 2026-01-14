<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-8 min-h-screen bg-transparent">
      
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
            
            <button v-if="!currentEnterprise" @click="openCreateModal" 
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
              <div v-if="ent.logo" class="w-14 h-14 rounded-xl border border-gray-200 p-1">
                  <img :src="ent.logo" class="w-full h-full object-cover rounded-lg">
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
        </div>
      </div>

      <!-- Enterprise Dashboard View -->
      <div v-else class="space-y-8 animate-fade-in">
         
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

         <!-- OVERVIEW TAB (New Quick Actions) -->
         <div v-if="currentTab === 'Overview'" class="space-y-8">
             <!-- Welcome Banner -->
             <div class="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl p-8 text-white shadow-lg relative overflow-hidden">
                 <div class="absolute right-0 top-0 h-full w-1/2 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-10"></div>
                 <h2 class="text-2xl font-bold relative z-10">Bienvenue sur {{ currentEnterprise.name }}</h2>
                 <p class="text-indigo-100 mt-2 relative z-10 max-w-xl">Accédez rapidement à vos outils de gestion. Configurez vos services, gérez votre personnel et suivez vos encaissements en temps réel.</p>
             </div>

             <!-- Quick Actions Grid -->
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <!-- Add Service Action -->
                    <div @click="goToSettingsAndAddService" class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all cursor-pointer border border-gray-100 dark:border-gray-700 group">
                        <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 text-blue-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                            <BoltIcon class="w-6 h-6" />
                        </div>
                        <h3 class="font-bold text-lg dark:text-white mb-1">Nouveau Service</h3>
                        <p class="text-sm text-gray-500">Ajoutez une prestation, une classe ou un abonnement.</p>
                    </div>

                    <!-- Add Member Action -->
                    <div @click="currentTab = 'Employees'" class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all cursor-pointer border border-gray-100 dark:border-gray-700 group">
                        <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 text-green-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                            <UserPlusIcon class="w-6 h-6" />
                        </div>
                        <h3 class="font-bold text-lg dark:text-white mb-1">Inviter Membre</h3>
                        <p class="text-sm text-gray-500">Ajoutez des employés ou gestionnaires à votre équipe.</p>
                    </div>

                    <!-- Billing Action -->
                    <div @click="currentTab = 'Billing'" class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm hover:shadow-md transition-all cursor-pointer border border-gray-100 dark:border-gray-700 group">
                        <div class="w-12 h-12 bg-purple-100 dark:bg-purple-900/30 text-purple-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                            <BanknotesIcon class="w-6 h-6" />
                        </div>
                        <h3 class="font-bold text-lg dark:text-white mb-1">Facturation</h3>
                        <p class="text-sm text-gray-500">Générez des factures ou saisissez des consommations.</p>
                    </div>
                </div>

             <!-- Stats / Recent Activity (Placeholder) -->
             <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                 <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
                     <h3 class="font-bold text-gray-900 dark:text-white mb-4">Aperçu Rapide</h3>
                     <div class="grid grid-cols-2 gap-4">
                         <div class="p-4 bg-gray-50 dark:bg-gray-750 rounded-xl">
                             <span class="text-sm text-gray-500">Total Abonnés</span>
                             <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">--</p>
                         </div>
                         <div class="p-4 bg-gray-50 dark:bg-gray-750 rounded-xl">
                             <span class="text-sm text-gray-500">Services Actifs</span>
                             <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">{{ currentEnterprise.custom_services?.length || 0 }}</p>
                         </div>
                     </div>
                 </div>
                 
                 <!-- QR Code Preview Mini -->
                 <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700 flex items-center justify-between">
                     <div>
                         <h3 class="font-bold text-gray-900 dark:text-white">QR Code Public</h3>
                         <p class="text-sm text-gray-500 mt-1 max-w-xs">Permettez à vos clients de s'abonner en scannant votre code unique.</p>
                         <button @click="showQRModal = true" class="mt-4 text-primary-600 font-medium hover:underline text-sm">Voir les codes &rarr;</button>
                     </div>
                     <div class="w-24 h-24 bg-white p-2 rounded-lg border border-gray-200">
                         <img :src="`/enterprise-service/api/v1/enterprises/${currentEnterprise.id}/qrcode`" class="w-full h-full object-contain">
                     </div>
                 </div>
             </div>
         </div>

         <!-- Employee Tab -->
         <div v-if="currentTab === 'Employees'" class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
            <div class="flex justify-between items-center mb-6">
               <h3 class="font-bold text-xl dark:text-white">Annuaire des Employés</h3>
               <button class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-xl text-sm font-medium shadow-sm transition-colors flex items-center gap-2">
                   <UserPlusIcon class="w-4 h-4" /> Inviter
               </button>
            </div>
            <!-- Mock Table -->
             <div class="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
                 <table class="w-full text-left text-sm text-gray-500">
                   <thead class="bg-gray-50 dark:bg-gray-900/50 text-gray-700 dark:text-gray-300">
                      <tr>
                         <th class="px-6 py-4 font-semibold">Nom</th>
                         <th class="px-6 py-4 font-semibold">Rôle</th>
                         <th class="px-6 py-4 font-semibold">Statut</th>
                         <th class="px-6 py-4 font-semibold text-right">Actions</th>
                      </tr>
                   </thead>
                   <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
                      <tr v-for="emp in employees" :key="emp.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                         <td class="px-6 py-4 font-medium text-gray-900 dark:text-white">{{ emp.first_name }} {{ emp.last_name }}</td>
                         <td class="px-6 py-4">{{ emp.profession }}</td>
                         <td class="px-6 py-4">
                            <span :class="{'bg-green-100 text-green-700 border-green-200': emp.status === 'ACTIVE', 'bg-yellow-100 text-yellow-700 border-yellow-200': emp.status === 'PENDING_INVITE'}" class="px-2.5 py-0.5 rounded-full text-xs font-medium border">
                               {{ formatEmployeeStatus(emp.status) }}
                            </span>
                         </td>
                         <td class="px-6 py-4 text-right">
                             <button class="text-gray-400 hover:text-gray-600 ml-2">
                                <EllipsisHorizontalIcon class="w-5 h-5" />
                             </button>
                         </td>
                      </tr>
                      <tr v-if="employees.length === 0">
                         <td colspan="4" class="px-6 py-12 text-center text-gray-400 bg-gray-50/50 dark:bg-gray-900/50">
                             <UserGroupIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
                             Aucun employé trouvé.
                         </td>
                      </tr>
                   </tbody>
                </table>
             </div>
         </div>

         <!-- CLIENTS Tab -->
         <div v-if="currentTab === 'Clients'" class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
            <!-- (Content remains similar but with updated styling, will rely on generic wrapper for now to save tokens, assuming similar table updates) -->
            <!-- ... Clients Tab Logic ... -->
            <!-- Re-inserting Clients content with container improvements -->
             <div class="flex justify-between items-center mb-6">
                <h3 class="font-bold text-xl dark:text-white">Gestion des Abonnés</h3>
                <div class="flex gap-2">
                    <button @click="downloadExport" class="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl text-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors flex items-center gap-2">
                        <ArrowDownTrayIcon class="w-4 h-4" /> Exporter
                    </button>
                    <button @click="showAddClientModal = true" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-sm font-medium shadow-sm transition-colors flex items-center gap-2">
                        <PlusIcon class="w-4 h-4" /> Nouveau Client
                    </button>
                </div>
            </div>

            <!-- Clients List (Simplified for brevity in this replace, ensuring styling matches) -->
              <div v-if="isLoading" class="text-center py-12">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
              </div>
              <div v-else>
                <div class="mb-4 flex gap-2">
                    <select v-model="selectedClientFilterService" @change="fetchClients" class="border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white rounded-lg px-3 py-2 text-sm shadow-sm focus:ring-primary-500 focus:border-primary-500">
                        <option value="">Tous les services</option>
                        <option v-for="svc in currentEnterprise.custom_services" :key="svc.id" :value="svc.id">{{ svc.name }}</option>
                    </select>
                </div>
                <div class="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
                    <table class="w-full text-left text-sm text-gray-500">
                        <thead class="bg-gray-50 dark:bg-gray-900/50 text-gray-700 dark:text-gray-300">
                            <tr>
                                <th class="px-6 py-3 font-semibold">Nom du Client</th>
                                <th class="px-6 py-3 font-semibold">Matricule</th>
                                <th class="px-6 py-3 font-semibold">Service</th>
                                <th class="px-6 py-3 font-semibold">Détails</th>
                                <th class="px-6 py-3 font-semibold">Proch. Fact.</th>
                                <th class="px-6 py-3 font-semibold text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
                            <tr v-for="sub in clientSubscriptions" :key="sub.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                <td class="px-6 py-4 font-medium text-gray-900 dark:text-white">{{ sub.client_name }}</td>
                                <td class="px-6 py-4 text-xs font-mono">
                                    <span v-if="sub.external_id" class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-gray-600 dark:text-gray-300">{{ sub.external_id }}</span>
                                    <span v-else class="text-gray-400">-</span>
                                </td>
                                <td class="px-6 py-4">
                                     {{ getServiceName(sub.service_id) }}
                                </td>
                                <td class="px-6 py-4">
                                    <div class="flex flex-col">
                                        <span class="text-gray-900 dark:text-white font-medium">{{ sub.amount }} {{ currentEnterprise.settings?.currency }}</span>
                                        <span class="text-xs text-gray-400 capitalize">{{ formatBillingFrequency(sub.billing_frequency) }}</span>
                                    </div>
                                </td>
                                <td class="px-6 py-4 text-gray-500">
                                    {{ new Date(sub.next_billing_at).toLocaleDateString() }}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <button class="text-red-500 hover:text-red-700 text-xs font-medium hover:underline">Résilier</button>
                                </td>
                            </tr>
                            <tr v-if="clientSubscriptions.length === 0">
                                <td colspan="6" class="px-6 py-12 text-center text-gray-400 bg-gray-50/50 dark:bg-gray-900/50">
                                    <UsersIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
                                    Aucun abonné pour le moment.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
              </div>
         </div>

         <!-- Add Client Modal -->
         <div v-if="showAddClientModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
            <div class="bg-white dark:bg-gray-800 rounded-xl max-w-lg w-full p-6 shadow-2xl">
                <h3 class="font-bold text-lg mb-4 dark:text-white">Ajouter un Client / Élève</h3>
                
                <!-- Step 1: Search User -->
                <div class="mb-4">
                    <label class="block text-sm font-medium mb-1">Rechercher un utilisateur (App Mobile)</label>
                    <div class="flex gap-2">
                        <input v-model="clientSearchQuery" placeholder="Téléphone ou Email" class="flex-1 border rounded px-3 py-2 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white">
                        <button @click="searchUser" :disabled="isSearchingUser" class="px-3 py-2 bg-blue-100 text-blue-700 rounded text-sm font-medium">
                            {{ isSearchingUser ? '...' : '🔍' }}
                        </button>
                    </div>
                    <div v-if="foundUser" class="mt-2 p-2 bg-green-50 text-green-700 rounded text-sm flex justify-between items-center">
                        <span>✅ {{ foundUser.name }}</span>
                        <button @click="foundUser = null; clientSearchQuery = ''" class="text-xs text-gray-500">X</button>
                    </div>
                </div>

                <!-- Step 2: Select Service -->
                <div v-if="foundUser" class="mb-4">
                    <label class="block text-sm font-medium mb-1">Service ou Classe</label>
                    <select v-model="newSubscription.service_id" class="w-full border rounded px-3 py-2 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white">
                        <option value="">Sélectionner...</option>
                        <option v-for="svc in currentEnterprise.custom_services" :key="svc.id" :value="svc.id">
                            {{ svc.name }}
                        </option>
                    </select>
                </div>
                
                <!-- External ID (Matricule) -->
                <div v-if="foundUser" class="mb-4">
                    <label class="block text-sm font-medium mb-1">ID Externe / Matricule (Optionnel)</label>
                    <input v-model="newSubscription.external_id" placeholder="ex: MAT-2024-001 ou Compteur #123" class="w-full border rounded px-3 py-2 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white">
                </div>
                
                <!-- Dynamic Form Fields -->
                <div v-if="foundUser && selectedServiceFormSchema && selectedServiceFormSchema.length" class="mb-4 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                    <h4 class="text-xs font-bold text-gray-500 uppercase mb-2">Informations Complémentaires</h4>
                    <div class="space-y-3">
                        <div v-for="field in selectedServiceFormSchema" :key="field.key">
                            <label class="block text-xs font-medium mb-1">
                                {{ field.label }} <span v-if="field.required" class="text-red-500">*</span>
                            </label>
                            
                            <input v-if="field.type !== 'select'" 
                                   v-model="newSubscription.form_data[field.key]" 
                                   :type="field.type" 
                                   :required="field.required"
                                   class="w-full border rounded px-2 py-1 text-sm dark:bg-gray-600 dark:border-gray-500 dark:text-white">
                                   
                            <select v-else 
                                    v-model="newSubscription.form_data[field.key]" 
                                    class="w-full border rounded px-2 py-1 text-sm dark:bg-gray-600 dark:border-gray-500 dark:text-white">
                                    <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                            </select>
                        </div>
                    </div>
                </div>

                <!-- Step 3: Specific Details (School - Legacy Support) -->
                 <div v-if="foundUser && currentEnterprise.type === 'SCHOOL' && !selectedServiceFormSchema?.length" class="mb-4 grid grid-cols-2 gap-2">
                     <div>
                         <label class="block text-xs font-medium mb-1">Nom de l'élève</label>
                         <input v-model="newSubscription.student_name" placeholder="Si différent du parent" class="w-full border rounded px-2 py-1 text-sm dark:bg-gray-700 dark:border-gray-600">
                     </div>
                     <div>
                         <label class="block text-xs font-medium mb-1">Classe</label>
                         <select v-model="newSubscription.class_id" class="w-full border rounded px-2 py-1 text-sm dark:bg-gray-700 dark:border-gray-600">
                             <option v-for="cls in currentEnterprise.school_config?.classes" :key="cls.name" :value="cls.name">{{ cls.name }}</option>
                         </select>
                     </div>
                 </div>

                <div class="flex justify-end gap-2 mt-6">
                    <button @click="showAddClientModal = false" class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg">Annuler</button>
                    <button @click="createClientSubscription" :disabled="!foundUser || !newSubscription.service_id" class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50">
                        Ajouter
                    </button>
                </div>
            </div>
         </div>

         <!-- Billing Tab -->
         <div v-if="currentTab === 'Billing'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <h3 class="font-semibold text-lg dark:text-white mb-4">Facturation en masse</h3>
            
            <div class="mb-6 space-y-4">
                <div v-if="currentEnterprise.type === 'UTILITY' || currentEnterprise.custom_services?.length">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Service Concerné</label>
                    <select v-model="selectedImportService" @change="loadManualSubscribers" class="w-full md:w-1/3 rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                        <option value="">Sélectionner un service...</option>
                        <option v-for="svc in currentEnterprise.custom_services" :key="svc.id" :value="svc.id">
                            {{ svc.name }} ({{ svc.billing_type === 'USAGE' ? 'Compteur' : 'Fixe' }})
                        </option>
                        <option v-if="currentEnterprise.type === 'SCHOOL'" value="tuition">Scolarité (Général)</option>
                        <option v-if="currentEnterprise.type === 'TRANSPORT'" value="subscription">Abonnement (Général)</option>
                    </select>
                </div>
                
                <!-- Billing Mode Toggle -->
                <div v-if="selectedImportService" class="flex gap-2">
                    <button @click="billingMode = 'IMPORT'" 
                            :class="{'bg-blue-100 text-blue-700 border-blue-200': billingMode === 'IMPORT', 'bg-gray-100 text-gray-600': billingMode !== 'IMPORT'}"
                            class="px-4 py-2 rounded-lg text-sm font-medium border flex items-center gap-2">
                            <span>📄</span> Importer CSV
                    </button>
                    <button @click="billingMode = 'MANUAL'; loadManualSubscribers()" 
                            :class="{'bg-purple-100 text-purple-700 border-purple-200': billingMode === 'MANUAL', 'bg-gray-100 text-gray-600': billingMode !== 'MANUAL'}"
                            class="px-4 py-2 rounded-lg text-sm font-medium border flex items-center gap-2">
                            <span>✍️</span> Saisie Manuelle
                    </button>
                </div>

                <!-- CSV Import Mode -->
                <div v-if="billingMode === 'IMPORT' && selectedImportService">
                    <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center bg-gray-50 dark:bg-gray-750">
                        <input type="file" @change="handleFileUpload" accept=".csv" class="hidden" id="invoice-upload">
                        <label for="invoice-upload" class="cursor-pointer">
                            <div class="text-4xl mb-2">📄</div>
                            <p class="text-gray-500 mb-2">Cliquez pour sélectionner un fichier CSV</p>
                            <p class="text-xs text-gray-400">Colonnes requises: ID Client [0], Montant [1] OU Consommation [2]</p>
                            <div v-if="importFile" class="mt-4 text-sm font-medium text-green-600">
                                Fichier sélectionné: {{ importFile.name }}
                            </div>
                        </label>
                    </div>

                    <div class="flex justify-end mt-4">
                        <button @click="uploadInvoiceFile" :disabled="!importFile || !selectedImportService" class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50">
                            Importer et Générer
                        </button>
                    </div>
                </div>

                <!-- Manual Entry Mode -->
                <div v-if="billingMode === 'MANUAL' && selectedImportService">
                    <div v-if="isLoading" class="text-center py-8 text-gray-500">Chargement des abonnés...</div>
                    <div v-else-if="manualSubscribers.length === 0" class="text-center py-8 text-gray-500">Aucun abonné trouvé éligible.</div>
                    <div v-else>
                         <div class="overflow-x-auto border rounded-lg dark:border-gray-700">
                            <table class="w-full text-sm text-left">
                                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                                    <tr>
                                        <th class="px-4 py-3">Client</th>
                                        <th class="px-4 py-3">ID Abonnement</th>
                                        <th class="px-4 py-3 text-center">Consommation</th>
                                        <th class="px-4 py-3 text-center">Montant (Optionnel)</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="sub in manualSubscribers" :key="sub.id" class="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
                                        <td class="px-4 py-3 font-medium">{{ sub.client_name || 'Inconnu' }}</td>
                                        <td class="px-4 py-3 text-gray-500 text-xs">{{ sub.id.substring(0,8) }}...</td>
                                        <td class="px-4 py-3 text-center">
                                            <input v-model.number="manualEntries[sub.id].consumption" type="number" step="0.1" class="w-24 px-2 py-1 text-right border rounded bg-white dark:bg-gray-900 border-gray-300 dark:border-gray-600">
                                        </td>
                                        <td class="px-4 py-3 text-center">
                                            <input v-model.number="manualEntries[sub.id].amount" type="number" step="50" class="w-28 px-2 py-1 text-right border rounded bg-white dark:bg-gray-900 border-gray-300 dark:border-gray-600">
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                         </div>
                         <div class="flex justify-between items-center mt-4">
                             <div class="text-xs text-gray-500">
                                {{ manualSubscribers.length }} abonnés chargés
                             </div>
                             <div class="flex gap-2">
                                 <button @click="loadManualSubscribers" class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50">
                                    Actualiser
                                 </button>
                                 <button @click="submitManualBatch" :disabled="isGeneratingManual" class="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50">
                                     {{ isGeneratingManual ? 'Génération...' : 'Valider et Générer' }}
                                 </button>
                             </div>
                         </div>
                    </div>
                </div>
         </div>
      </div>
      
         <!-- Settings Tab -->
         <div v-if="currentTab === 'Settings'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <h3 class="font-semibold text-lg dark:text-white mb-6">Configuration de l'Entreprise</h3>

            <!-- General Settings -->
            <div class="mb-8 p-4 border rounded-lg dark:border-gray-700">
                <h4 class="font-medium mb-4 dark:text-gray-200">Paramètres Généraux</h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Nom</label>
                        <input v-model="currentEnterprise.name" type="text" class="mt-1 block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                    </div>
                     <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Devise</label>
                        <select v-model="currentEnterprise.settings.currency" class="mt-1 block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                            <option value="XOF">XOF (FCFA)</option>
                            <option value="EUR">EUR (€)</option>
                            <option value="USD">USD ($)</option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- SCHOOL Specific Settings -->
            <div v-if="currentEnterprise.type === 'SCHOOL'" class="mb-8 p-4 border rounded-lg dark:border-gray-700 bg-blue-50 dark:bg-blue-900/10">
                <h4 class="font-medium mb-4 text-blue-800 dark:text-blue-300 flex items-center gap-2">
                    <span>🎓 Configuration École</span>
                </h4>
                
                <!-- Classes -->
                <div class="mb-6">
                    <div class="flex justify-between items-center mb-2">
                        <label class="text-sm font-medium dark:text-gray-300">Classes</label>
                        <button @click="addClass" class="text-xs bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-700">+ Ajouter Classe</button>
                    </div>
                    <div v-if="!currentEnterprise.school_config?.classes?.length" class="text-sm text-gray-500 italic">Aucune classe configurée.</div>
                    <div v-else class="space-y-3">
                        <div v-for="(cls, idx) in currentEnterprise.school_config.classes" :key="idx" class="flex gap-2 items-start bg-white dark:bg-gray-800 p-3 rounded shadow-sm">
                            <input v-model="cls.name" placeholder="Nom (ex: CP)" class="flex-1 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                            <input v-model.number="cls.total_fees" type="number" placeholder="Frais Total" class="w-32 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                            <button @click="removeClass(idx)" class="text-red-500 hover:text-red-700 text-sm px-2">Suppr.</button>
                        </div>
                        
                        <!-- Tranches / Echéances Management -->
                        <div class="ml-4 pl-4 border-l-2 border-gray-200 dark:border-gray-600 mt-2">
                             <div class="flex items-center gap-2 mb-2">
                                <label class="text-xs font-semibold text-gray-500 uppercase">Échéances / Tranches (Optionnel)</label>
                                <button @click="addTranche(cls)" class="text-xs bg-gray-200 dark:bg-gray-600 px-2 py-0.5 rounded hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-700 dark:text-gray-300">+ Ajouter</button>
                             </div>
                             <div class="space-y-2">
                                <div v-for="(tr, trIdx) in cls.tranches" :key="trIdx" class="flex items-center gap-2">
                                    <input v-model="tr.name" placeholder="Nom (ex: 1ère Tranche)" class="flex-1 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <input v-model.number="tr.amount" type="number" placeholder="Montant" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <input v-model="tr.due_date" type="date" class="w-32 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <button @click="removeTranche(cls, trIdx)" class="text-red-400 hover:text-red-600 text-xs font-bold px-1">×</button>
                                </div>
                             </div>
                             <div v-if="cls.tranches?.length" class="text-xs text-blue-600 mt-1 font-medium">
                                Total Tranches: {{ cls.tranches.reduce((sum, t) => sum + (t.amount || 0), 0) }} / {{ cls.total_fees }}
                             </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- TRANSPORT Specific Settings -->
            <div v-if="currentEnterprise.type === 'TRANSPORT'" class="mb-8 p-4 border rounded-lg dark:border-gray-700 bg-green-50 dark:bg-green-900/10">
                <h4 class="font-medium mb-4 text-green-800 dark:text-green-300 flex items-center gap-2">
                    <span>🚌 Configuration Transport</span>
                </h4>

                <!-- Routes -->
                <div class="mb-6">
                    <div class="flex justify-between items-center mb-2">
                        <label class="text-sm font-medium dark:text-gray-300">Lignes / Trajets</label>
                        <button @click="addRoute" class="text-xs bg-green-600 text-white px-2 py-1 rounded hover:bg-green-700">+ Ajouter Ligne</button>
                    </div>
                    <div v-if="!currentEnterprise.transport_config?.routes?.length" class="text-sm text-gray-500 italic">Aucune ligne configurée.</div>
                     <div v-else class="space-y-3">
                        <div v-for="(route, idx) in currentEnterprise.transport_config.routes" :key="idx" class="flex gap-2 items-start bg-white dark:bg-gray-800 p-3 rounded shadow-sm">
                            <input v-model="route.name" placeholder="Nom (ex: Ligne 14)" class="flex-1 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                            <input v-model.number="route.base_price" type="number" placeholder="Prix Base" class="w-32 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                            <button @click="removeRoute(idx)" class="text-red-500 hover:text-red-700 text-sm px-2">Suppr.</button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Service Groups Management -->
            <div class="mb-8 p-4 border rounded-lg dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                <h4 class="font-medium mb-4 dark:text-gray-200 flex items-center gap-2">
                    <span>⚡ Groupes & Services</span>
                </h4>
                
                <div class="mb-6 flex justify-between items-center">
                    <p class="text-sm text-gray-500">Organisez vos services par catégories (ex: Scolarité, Transport).</p>
                    <button @click="addServiceGroup" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700 flex items-center gap-2">
                        <PlusIcon class="w-4 h-4" /> Nouveau Groupe
                    </button>
                </div>

                <div v-if="!currentEnterprise.service_groups?.length" class="text-center py-8 text-gray-400 border-2 border-dashed border-gray-300 rounded-xl">
                    Aucun groupe de services configuré.
                </div>

                <!-- Groups Loop -->
                <div v-for="(group, gIdx) in currentEnterprise.service_groups" :key="group.id" class="mb-6 p-4 border rounded-xl dark:border-gray-600 bg-white dark:bg-gray-800 shadow-sm">
                     
                     <!-- Group Header Config -->
                     <div class="flex flex-col md:flex-row gap-4 mb-4 pb-4 border-b border-gray-100 dark:border-gray-700">
                        <div class="flex-1">
                            <label class="block text-xs font-semibold text-gray-500 uppercase mb-1">Nom du Groupe</label>
                            <input v-model="group.name" placeholder="ex: Scolarité" class="w-full rounded-md border-gray-300 text-base font-bold dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                        </div>
                        <div class="w-full md:w-32">
                            <label class="block text-xs font-semibold text-gray-500 uppercase mb-1">Devise</label>
                            <select v-model="group.currency" class="w-full rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-2 border font-mono">
                                <option v-for="curr in currencies" :key="curr.code" :value="curr.code">
                                    {{ curr.code }} - {{ curr.name }}
                                </option>
                            </select>
                        </div>
                        <div class="flex items-center pt-5">
                            <label class="flex items-center gap-2 cursor-pointer">
                                <input type="checkbox" v-model="group.is_private" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600">
                                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Privé (Inv. seule)</span>
                            </label>
                        </div>
                        <div class="pt-5">
                             <button @click="removeServiceGroup(gIdx)" class="text-red-500 hover:bg-red-50 p-2 rounded-lg transition-colors" title="Supprimer le groupe">
                                <TrashIcon class="w-5 h-5" />
                             </button>
                        </div>
                     </div>

                     <!-- Services in Group -->
                     <div class="space-y-4 pl-0 md:pl-4">
                        <div v-if="!group.services?.length" class="text-sm text-gray-400 italic mb-2">
                            Aucun service dans ce groupe.
                        </div>

                        <div v-for="(svc, idx) in group.services" :key="svc.uid" class="flex flex-col gap-2 p-3 rounded shadow-sm border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700/50">
                             <div class="flex gap-2">
                                <input v-model="svc.id" placeholder="ID (ex: water_usage)" class="w-1/4 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border font-mono">
                                <input v-model="svc.name" placeholder="Nom (ex: Consommation Eau)" class="flex-1 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                <button @click="removeCustomService(group, svc)" class="text-red-500 hover:text-red-700 text-sm px-2">Suppr.</button>
                             </div>
                             
                             <!-- Config Fields -->
                             <div class="flex gap-2 items-center flex-wrap">
                                <select v-model="svc.billing_type" class="w-32 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <option value="FIXED">Fixe / Forfait</option>
                                    <option value="USAGE">À l'usage</option>
                                </select>
                                <select v-model="svc.billing_frequency" class="w-32 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <option value="DAILY">Journalier</option>
                                    <option value="WEEKLY">Hebdomadaire</option>
                                    <option value="MONTHLY">Mensuel</option>
                                    <option value="ANNUALLY">Annuel</option>
                                    <option value="CUSTOM">Personnalisé</option>
                                    <option value="ONETIME">Ponctuel</option>
                                </select>
                                
                                <input v-if="svc.billing_frequency === 'CUSTOM' && !svc.use_schedule" v-model.number="svc.custom_interval" type="number" placeholder="Jours" class="w-20 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                
                                <input v-if="svc.billing_type === 'USAGE'" v-model="svc.unit" placeholder="Unité" class="w-20 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                
                                <div class="relative">
                                    <input v-model.number="svc.base_price" type="number" step="0.01" :placeholder="svc.billing_type === 'USAGE' ? 'Prix Unit.' : 'Montant'" class="w-28 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white pl-2 pr-8 py-1 border">
                                    <span class="absolute right-2 top-1.5 text-xs text-gray-500 font-bold">{{ group.currency }}</span>
                                </div>
                             </div>

                             <!-- Payment Schedule UI -->
                             <div v-if="svc.billing_frequency === 'CUSTOM'" class="mt-2 ml-2 pl-4 border-l-2 border-orange-200 dark:border-orange-800">
                                   <!-- ... Schedule Content ... -->
                                   <div class="flex items-center gap-2 mb-2">
                                       <label class="text-xs font-semibold text-gray-500 uppercase">Calendrier de Paiement</label>
                                       <div class="flex items-center gap-2 ml-4">
                                           <button @click="toggleSchedule(svc)" :class="{'bg-orange-100 text-orange-700': svc.use_schedule, 'text-gray-400': !svc.use_schedule}" class="text-xs px-2 py-0.5 rounded border border-gray-200 dark:border-gray-600">
                                               {{ svc.use_schedule ? 'Mode: Calendrier' : 'Mode: Intervalle' }}
                                           </button>
                                       </div>
                                   </div>
                                   
                                   <div v-if="svc.use_schedule">
                                         <div v-for="(item, sIdx) in svc.payment_schedule" :key="sIdx" class="flex items-center gap-2 mb-1">
                                            <input v-model="item.name" placeholder="Nom" class="flex-1 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <input v-model="item.start_date" type="date" class="w-28 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <span class="text-xs text-gray-400">au</span>
                                            <input v-model="item.end_date" type="date" class="w-28 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <input v-model.number="item.amount" type="number" placeholder="Montant" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <button @click="removeScheduleItem(svc, sIdx)" class="text-red-400 hover:text-red-600 text-xs font-bold px-1">×</button>
                                         </div>
                                         <button @click="addScheduleItem(svc)" class="text-xs text-blue-600 hover:text-blue-800 mt-1">+ Ajouter une période</button>
                                   </div>
                             </div>

                             <!-- Form Builder & Penalty -->
                             <div class="flex flex-col gap-2">
                                <!-- Form Builder -->
                                 <div class="mt-2 ml-2 pl-4 border-l-2 border-indigo-200 dark:border-indigo-800">
                                       <div class="flex items-center gap-2 mb-2">
                                           <label class="text-xs font-semibold text-gray-500 uppercase">Formulaire</label>
                                           <button @click="addFormField(svc)" class="text-xs bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded border border-indigo-200">+ Champ</button>
                                       </div>
                                       <div v-if="svc.form_schema?.length">
                                             <div v-for="(field, fIdx) in svc.form_schema" :key="fIdx" class="flex items-center gap-2 mb-1 bg-gray-50 dark:bg-gray-700 p-1 rounded">
                                                <input v-model="field.label" placeholder="Libellé" class="flex-1 rounded-md border-gray-300 text-xs dark:bg-gray-600 dark:border-gray-500 dark:text-white px-2 py-1 border">
                                                <select v-model="field.type" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-600 dark:border-gray-500 dark:text-white px-1 py-1 border">
                                                    <option value="text">Texte</option>
                                                    <option value="number">Nombre</option>
                                                    <option value="date">Date</option>
                                                    <option value="select">Liste</option>
                                                </select>
                                                <label class="flex items-center gap-1 text-xs text-gray-500">
                                                    <input type="checkbox" v-model="field.required">
                                                    <span>Req.</span>
                                                </label>
                                                <button @click="removeFormField(svc, fIdx)" class="text-red-400 hover:text-red-600 text-xs font-bold px-1">×</button>
                                             </div>
                                       </div>
                                 </div>

                                 <!-- Penalty Configuration -->
                                <div class="mt-2 ml-2 pl-4 border-l-2 border-red-200 dark:border-red-800">
                                   <!-- ... Same Penalty UI ... -->
                                   <div class="flex items-center gap-2 mb-2">
                                       <label class="text-xs font-semibold text-gray-500 uppercase">Pénalités</label>
                                       <button @click="togglePenalty(svc)" :class="{'bg-red-100 text-red-700 border-red-200': svc.penalty_config, 'text-gray-400 border-gray-200': !svc.penalty_config}" class="text-xs px-2 py-0.5 rounded border">
                                           {{ svc.penalty_config ? 'Configuré' : 'Aucune' }}
                                       </button>
                                   </div>
                                   <div v-if="svc.penalty_config" class="bg-red-50 dark:bg-red-900/10 p-2 rounded-lg space-y-2">
                                       <div class="flex flex-col md:flex-row gap-2">
                                           <select v-model="svc.penalty_config.type" class="rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                               <option value="FIXED">Montant Fixe</option>
                                               <option value="PERCENTAGE">Pourcentage (%)</option>
                                               <option value="HYBRID">Hybride</option>
                                           </select>
                                            <div class="flex gap-2">
                                                <input v-model.number="svc.penalty_config.value" type="number" :placeholder="svc.penalty_config.type === 'FIXED' ? 'Montant' : '%'" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                                <input v-if="svc.penalty_config.type === 'HYBRID'" v-model.number="svc.penalty_config.value_fixed" type="number" placeholder="+ Fixe" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            </div>
                                       </div>
                                       <div class="flex gap-2 items-center flex-wrap">
                                           <span class="text-xs text-gray-500">Fréq:</span>
                                           <select v-model="svc.penalty_config.frequency" class="rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                               <option value="DAILY">Jour</option>
                                               <option value="WEEKLY">Semaine</option>
                                               <option value="MONTHLY">Mois</option>
                                               <option value="ONETIME">Unique</option>
                                           </select>
                                           <span class="text-xs text-gray-500">Grâce:</span>
                                           <input v-model.number="svc.penalty_config.grace_period" type="number" placeholder="Jours" class="w-16 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                       </div>
                                   </div>
                                </div>
                             </div>
                        </div>

                        <div class="pt-2">
                            <button @click="addCustomService(group)" class="w-full py-2 border-2 border-dashed border-gray-200 dark:border-gray-600 rounded-lg text-sm text-gray-500 hover:border-gray-400 hover:text-gray-700 transition-colors">
                                + Ajouter un Service dans {{ group.name }}
                            </button>
                        </div>
                     </div>
                </div>
            </div>

            <!-- Save Actions -->
            <div class="flex justify-end pt-4 border-t dark:border-gray-700">
                <button @click="saveSettings" :disabled="isSaving" class="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 flex items-center gap-2">
                    <span v-if="isSaving" class="animate-spin text-white">⟳</span>
                    {{ isSaving ? 'Enregistrement...' : 'Enregistrer les modifications' }}
                </button>
            </div>
         </div>
      </div>
      
      <!-- QR Code Modal -->
      <div v-if="showQRModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white dark:bg-gray-800 rounded-xl w-full max-w-2xl p-6 space-y-6 max-h-[90vh] overflow-y-auto">
             <div class="flex justify-between items-center">
                <h2 class="text-xl font-bold dark:text-white">Codes QR - {{ currentEnterprise?.name }}</h2>
                <button @click="showQRModal = false" class="text-gray-500 hover:text-gray-700">✕</button>
             </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                 <div class="col-span-full bg-gradient-to-br from-gray-900 to-gray-800 rounded-2xl p-8 text-white relative overflow-hidden shadow-xl">
                    <div class="relative z-10 flex flex-col md:flex-row items-center justify-between gap-6">
                        <div>
                            <h3 class="text-2xl font-bold mb-2">QR Codes & Paiements</h3>
                            <p class="text-gray-300 max-w-lg">
                                Générez et téléchargez les codes QR pour votre entreprise, vos groupes de services ou des services spécifiques. 
                                Permettez à vos clients de s'abonner et de payer instantanément.
                            </p>
                        </div>
                        <button @click="showQRModal = true" class="px-6 py-3 bg-white text-gray-900 font-bold rounded-xl hover:bg-gray-100 transition-colors flex items-center gap-2 shadow-lg">
                            <QrCodeIcon class="w-6 h-6" />
                            Afficher les Codes QR
                        </button>
                    </div>
                    <!-- Background decoration -->
                    <div class="absolute -right-20 -bottom-20 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl"></div>
                 </div>
             </div>
        </div>
      </div>

      <!-- Sophisticated QR Modal -->
      <div v-if="showQRModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div class="bg-gray-900 rounded-3xl w-full max-w-4xl max-h-[90vh] overflow-hidden shadow-2xl flex flex-col md:flex-row border border-gray-700">
              
              <!-- Sidebar / Controls -->
              <div class="w-full md:w-1/3 bg-gray-800 p-6 flex flex-col gap-6 border-r border-gray-700">
                  <div>
                      <h3 class="text-xl font-bold text-white mb-4">Type de Code</h3>
                      <div class="space-y-2">
                          <button @click="qrTab = 'ENTERPRISE'" :class="{'bg-primary-600 text-white': qrTab === 'ENTERPRISE', 'bg-gray-700 text-gray-300 hover:bg-gray-600': qrTab !== 'ENTERPRISE'}" class="w-full text-left px-4 py-3 rounded-xl font-medium transition-all flex items-center gap-3">
                              <BuildingOfficeIcon class="w-5 h-5" /> Entreprise (Global)
                          </button>
                          <button @click="qrTab = 'GROUP'" :class="{'bg-primary-600 text-white': qrTab === 'GROUP', 'bg-gray-700 text-gray-300 hover:bg-gray-600': qrTab !== 'GROUP'}" class="w-full text-left px-4 py-3 rounded-xl font-medium transition-all flex items-center gap-3">
                              <FolderIcon class="w-5 h-5" /> Groupe de Services
                          </button>
                          <button @click="qrTab = 'SERVICE'" :class="{'bg-primary-600 text-white': qrTab === 'SERVICE', 'bg-gray-700 text-gray-300 hover:bg-gray-600': qrTab !== 'SERVICE'}" class="w-full text-left px-4 py-3 rounded-xl font-medium transition-all flex items-center gap-3">
                              <TagIcon class="w-5 h-5" /> Service Spécifique
                          </button>
                      </div>
                  </div>

                  <!-- Filters -->
                  <div v-if="qrTab === 'GROUP' || qrTab === 'SERVICE'" class="space-y-4 animate-fadeIn">
                       <div>
                           <label class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 block">Choisir un Groupe</label>
                           <select v-model="selectedQRGroup" class="w-full bg-gray-700 border-none rounded-lg text-white px-3 py-2 focus:ring-2 focus:ring-primary-500">
                               <option v-for="grp in currentEnterprise.service_groups.filter(g => !g.is_private)" :key="grp.id" :value="grp.id">
                                   {{ grp.name }}
                               </option>
                           </select>
                       </div>
                       
                       <div v-if="qrTab === 'SERVICE' && selectedQRGroup">
                           <label class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 block">Choisir un Service</label>
                           <select v-model="selectedQRService" class="w-full bg-gray-700 border-none rounded-lg text-white px-3 py-2 focus:ring-2 focus:ring-primary-500">
                               <option v-for="svc in currentEnterprise.service_groups.find(g => g.id === selectedQRGroup)?.services" :key="svc.id" :value="svc.id">
                                   {{ svc.name }}
                               </option>
                           </select>
                       </div>
                  </div>

                  <div class="mt-auto">
                      <button @click="showQRModal = false" class="text-gray-400 hover:text-white text-sm flex items-center gap-2">
                          <ArrowLeftIcon class="w-4 h-4" /> Retour au tableau de bord
                      </button>
                  </div>
              </div>

              <!-- Preview Area -->
              <div class="w-full md:w-2/3 p-8 flex flex-col items-center justify-center bg-gray-900 relative">
                  <!-- QR Content -->
                  <div class="bg-white p-4 rounded-2xl shadow-2xl transform transition-all hover:scale-105 duration-300">
                      <img :src="getCurrentQRLink()" alt="QR Code" class="w-64 h-64 object-contain">
                  </div>
                  
                  <div class="mt-8 text-center space-y-2">
                       <h2 class="text-2xl font-bold text-white">{{ getQRTitle() }}</h2>
                       <p class="text-gray-400 max-w-md mx-auto">{{ getQRDescription() }}</p>
                  </div>

                  <!-- Actions -->
                  <div class="flex gap-4 mt-8">
                       <a :href="getCurrentQRLink()" :download="`qr_${qrTab.toLowerCase()}.png`" class="px-6 py-3 bg-primary-600 text-white rounded-xl font-bold hover:bg-primary-500 transition-colors flex items-center gap-2 shadow-lg shadow-primary-900/50">
                           <ArrowDownTrayIcon class="w-5 h-5" />
                           Télécharger PNG
                       </a>
                       <div class="bg-gray-800 rounded-xl px-4 py-3 flex items-center gap-3 border border-gray-700">
                           <span class="font-mono text-primary-400 font-bold tracking-wider">{{ getQRLabel() }}</span>
                           <button @click="copyCode(getQRLabel())" class="text-gray-400 hover:text-white">
                               <DocumentDuplicateIcon class="w-5 h-5" />
                           </button>
                       </div>
                  </div>
              </div>
          </div>
      </div>

      <!-- Add Service Modal -->
      <div v-if="showAddServiceModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
          <div class="bg-white dark:bg-gray-800 rounded-xl w-full max-w-md p-6 space-y-4">
              <h3 class="font-bold text-lg dark:text-white">Nouveau Service</h3>
              
              <div>
                  <label class="block text-sm font-medium mb-1 dark:text-gray-300">Groupe de Service</label>
                  <select v-model="newService.group_id" class="w-full rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2">
                      <option v-for="grp in currentEnterprise.service_groups" :key="grp.id" :value="grp.id">
                          {{ grp.name }}
                      </option>
                  </select>
              </div>

              <div>
                  <label class="block text-sm font-medium mb-1 dark:text-gray-300">Nom du Service</label>
                  <input v-model="newService.name" placeholder="ex: Frais de scolarité, Cantine" class="w-full rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2">
              </div>

              <div class="grid grid-cols-2 gap-4">
                  <div>
                      <label class="block text-sm font-medium mb-1 dark:text-gray-300">Type</label>
                      <select v-model="newService.billing_type" class="w-full rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 text-sm">
                           <option value="FIXED">Fixe</option>
                           <option value="USAGE">Usage</option>
                      </select>
                  </div>
                  <div>
                      <label class="block text-sm font-medium mb-1 dark:text-gray-300">Fréquence</label>
                      <select v-model="newService.billing_frequency" class="w-full rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 text-sm">
                           <option value="DAILY">Journalier</option>
                           <option value="WEEKLY">Hebdo</option>
                           <option value="MONTHLY">Mensuel</option>
                           <option value="ANNUALLY">Annuel</option>
                           <option value="ONETIME">Unique</option>
                      </select>
                  </div>
              </div>

              <div>
                  <label class="block text-sm font-medium mb-1 dark:text-gray-300">Prix de base</label>
                  <div class="relative">
                       <input v-model.number="newService.base_price" type="number" class="w-full rounded-lg border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:text-white pl-3 pr-12 py-2">
                       <span class="absolute right-3 top-2 text-gray-500 text-sm">
                           {{ currentEnterprise.service_groups?.find(g => g.id === newService.group_id)?.currency || 'XOF' }}
                       </span>
                  </div>
              </div>

              <div class="flex justify-end gap-3 pt-4">
                  <button @click="showAddServiceModal = false" class="px-4 py-2 text-gray-500 hover:text-gray-700">Annuler</button>
                  <button @click="confirmAddService" class="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">Créer</button>
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
                    <input v-model="newEnterprise.name" type="text" required placeholder="Ex: Ma Super École" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                </div>
                
                <div>
                     <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Logo de l'entreprise</label>
                     <div class="mt-1 flex items-center gap-4">
                        <div v-if="newEnterprise.logoPreview" class="w-16 h-16 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden border border-gray-200 dark:border-gray-600">
                           <img :src="newEnterprise.logoPreview" alt="Logo" class="w-full h-full object-cover">
                        </div>
                        <div v-else class="w-16 h-16 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-gray-400 border border-dashed border-gray-300 dark:border-gray-600">
                           📷
                        </div>
                        <input type="file" @change="handleLogoUpload" accept="image/*" class="text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100 dark:file:bg-gray-700 dark:file:text-white">
                     </div>
                </div>

                <div>
                   <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Nombre moyen d'employés</label>
                   <select v-model="newEnterprise.employee_count_range" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border">
                       <option value="1-10">1 - 10 employés</option>
                       <option value="11-50">11 - 50 employés</option>
                       <option value="51-200">51 - 200 employés</option>
                       <option value="201-500">201 - 500 employés</option>
                       <option value="500+">Plus de 500 employés</option>
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
import { ref, onMounted, watch, computed } from 'vue' 
import { enterpriseAPI, useApi } from '@/composables/useApi'
import { 
    Squares2X2Icon, 
    UsersIcon, 
    UserGroupIcon, 
    BanknotesIcon, 
    DocumentTextIcon, 
    Cog6ToothIcon, 
    ChevronRightIcon, 
    QrCodeIcon, 
    PlusIcon, 
    ArrowRightCircleIcon, 
    CheckCircleIcon,
    BuildingOffice2Icon,
    BoltIcon,
    UserPlusIcon,
    EllipsisHorizontalIcon,
    ArrowDownTrayIcon,
    TrashIcon,
    BuildingOfficeIcon,
    FolderIcon,
    TagIcon,
    ArrowLeftIcon,
    DocumentDuplicateIcon
} from '@heroicons/vue/24/outline'

const { authApi } = useApi()

// Translation Helpers
const formatEnterpriseType = (type) => {
    const map = {
        'SME': 'PME / Standard',
        'SCHOOL': 'École',
        'TRANSPORT': 'Transport',
        'UTILITY': 'Service Public',
        'NGO': 'ONG'
    }
    return map[type] || type
}

const formatEmployeeStatus = (status) => {
    const map = {
        'ACTIVE': 'Actif',
        'PENDING_INVITE': 'Invitation en attente',
        'INACTIVE': 'Inactif'
    }
    return map[status] || status
}

const formatBillingFrequency = (freq) => {
    const map = {
        'DAILY': 'Quotidien',
        'WEEKLY': 'Hebdomadaire',
        'MONTHLY': 'Mensuel',
        'ANNUALLY': 'Annuel',
        'CUSTOM': 'Personnalisé',
        'ONETIME': 'Ponctuel'
    }
    return map[freq] || freq
}

// Tab Icons Helper
const getTabIcon = (tab) => {
    switch (tab) {
        case 'Overview': return Squares2X2Icon
        case 'Employees': return UsersIcon
        case 'Clients': return UserGroupIcon
        case 'Payroll': return BanknotesIcon
        case 'Billing': return DocumentTextIcon
        case 'Settings': return Cog6ToothIcon
        default: return Squares2X2Icon
    }
}



// ... Tabs ...

// Client Management State
const showAddClientModal = ref(false)
const clientSubscriptions = ref([])
const selectedClientFilterService = ref('')
const clientSearchQuery = ref('')
const foundUser = ref(null)
const isSearchingUser = ref(false)
const showQRModal = ref(false)
const newSubscription = ref({
    service_id: '',
    student_name: '',
    class_id: '',
    external_id: '', // Matricule
    form_data: {} // Dynamic Data
})

const selectedServiceFormSchema = computed(() => {
    if (!newSubscription.value.service_id || !currentEnterprise.value) return []
    const svc = currentEnterprise.value.custom_services.find(s => s.id === newSubscription.value.service_id)
    return svc ? svc.form_schema : []
})

const downloadExport = () => {
    // Simple CSV Export
    if (!clientSubscriptions.value.length) return
    let csv = 'Nom Client,Service/Classe,Matricule,Frequence,Montant,Prochaine Facturation\n'
    clientSubscriptions.value.forEach(sub => {
        const serviceName = getServiceName(sub.service_id)
        const matricule = sub.external_id || (sub.school_details?.student_name ? sub.school_details.student_name : 'N/A')
        csv += `"${sub.client_name}","${serviceName}","${matricule}","${sub.billing_frequency}","${sub.amount}","${new Date(sub.next_billing_at).toLocaleDateString()}"\n`
    })
    
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `clients_export_${currentEnterprise.value.name}_${new Date().toISOString().split('T')[0]}.csv`
    a.click()
}

const fetchClients = async () => {
    if (!currentEnterprise.value) return
    try {
        isLoading.value = true
        const { data } = await enterpriseAPI.getSubscriptions(currentEnterprise.value.id, selectedClientFilterService.value)
        // Default Sort by Name
        clientSubscriptions.value = (data || []).sort((a, b) => (a.client_name || '').localeCompare(b.client_name || ''))
    } catch (e) {
        console.error('Failed to fetch clients', e)
    } finally {
        isLoading.value = false
    }
}

const searchUser = async () => {
    if (!clientSearchQuery.value) return
    isSearchingUser.value = true
    foundUser.value = null
    try {
        const { data } = await authApi.lookup(clientSearchQuery.value)
        if (data) foundUser.value = data
        else alert('Utilisateur non trouvé')
    } catch (e) {
        console.error(e)
        alert('Utilisateur non trouvé')
    } finally {
        isSearchingUser.value = false
    }
}

const copyCode = (code) => {
    navigator.clipboard.writeText(code)
    alert('Code copié : ' + code)
}

const createClientSubscription = async () => {
    if (!currentEnterprise.value || !foundUser.value) return
    try {
        const payload = {
            client_id: foundUser.value.id,
            client_name: foundUser.value.name, // Snapshot
            ...newSubscription.value
        }
        
        // Resolve Service Amount/Freq
        const svc = currentEnterprise.value.custom_services.find(s => s.id === payload.service_id)
        if (svc) {
            payload.amount = svc.base_price
            payload.billing_frequency = svc.billing_frequency
        }

        await enterpriseAPI.createSubscription(currentEnterprise.value.id, payload)
        
        showAddClientModal.value = false
        alert('Client ajouté avec succès')
        // Reset
        newSubscription.value = { service_id: '', student_name: '', class_id: '' }
        foundUser.value = null
        clientSearchQuery.value = ''
        fetchClients()
    } catch (e) {
        console.error(e)
        alert("Erreur lors de l'ajout")
    }
}

const getServiceName = (id) => {
    if (!currentEnterprise.value) return id
    const svc = currentEnterprise.value.custom_services.find(s => s.id === id)
    if (svc) return svc.name
    // Check school classes? No, they don't have ID per se, just names.
    return id
}

// Watch tab change to fetch clients


// ... Existing Code ...

const currencies = [
    { code: 'XOF', name: 'Franc CFA (BCEAO)' },
    { code: 'XAF', name: 'Franc CFA (BEAC)' },
    { code: 'EUR', name: 'Euro' },
    { code: 'USD', name: 'Dollar Américain' },
    { code: 'GBP', name: 'Livre Sterling' },
    { code: 'CAD', name: 'Dollar Canadien' },
    { code: 'CHF', name: 'Franc Suisse' },
    { code: 'CNY', name: 'Yuan Renminbi' },
    { code: 'JPY', name: 'Yen Japonais' },
    { code: 'AED', name: 'Dirham EAU' },
    { code: 'SAR', name: 'Riyal Saoudien' },
    { code: 'NGN', name: 'Naira Nigérian' },
    { code: 'GHS', name: 'Cedi Ghanéen' },
    { code: 'KES', name: 'Shilling Kenyan' },
    { code: 'ZAR', name: 'Rand Sud-Africain' },
    { code: 'MAD', name: 'Dirham Marocain' },
    { code: 'EGP', name: 'Livre Égyptienne' },
    { code: 'CDF', name: 'Franc Congolais' },
    { code: 'GNF', name: 'Franc Guinéen' },
    { code: 'RWF', name: 'Franc Rwandais' },
    { code: 'TND', name: 'Dinar Tunisien' },
    { code: 'DZD', name: 'Dinar Algérien' },
    { code: 'ETB', name: 'Birr Éthiopien' },
    { code: 'INR', name: 'Roupie Indienne' },
    { code: 'BRL', name: 'Réal Brésilien' },
    { code: 'RUB', name: 'Rouble Russe' },
    { code: 'TRY', name: 'Lire Turque' },
]

const tabs = ['Overview', 'Employees', 'Clients', 'Payroll', 'Billing', 'Settings']
const tabLabels = {
    'Overview': 'Aperçu',
    'Employees': 'Employés',
    'Clients': 'Abonnés / Élèves',
    'Payroll': 'Paie',
    'Billing': 'Facturation',
    'Settings': 'Paramètres'
}

const currentTab = ref('Overview')
const currentEnterprise = ref(null)

const enterprises = ref([])
const employees = ref([])
const isLoading = ref(true)
const isSaving = ref(false)

// Import State
const selectedImportService = ref('')
const importFile = ref(null)

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
      enterprises.value = data || []
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
  // Ensure nested objects exist for v-model binding
  if (!ent.settings) ent.settings = { currency: 'XOF', auto_pay_salaries: false }
  
  if (ent.type === 'SCHOOL' && !ent.school_config) ent.school_config = { classes: [] }
  if (ent.type === 'TRANSPORT' && !ent.transport_config) ent.transport_config = { routes: [], zones: [] }
  if (!ent.service_groups) ent.service_groups = [] // Ensure service_groups exists
  
  currentEnterprise.value = JSON.parse(JSON.stringify(ent)) // Deep copy to avoid mutating list directly
  selectedImportService.value = '' // Reset import selection
}

const fetchEmployees = async () => {
   if (!currentEnterprise.value) return
   try {
      const { data } = await enterpriseAPI.listEmployees(currentEnterprise.value.id)
      employees.value = data || []
   } catch (error) {
      console.error('Failed to fetch employees', error)
      employees.value = []
   }
}



const openCreateModal = () => {
  showCreateModal.value = true
  newEnterprise.value = { 
      name: '', 
      employee_count_range: '1-10',
      logo: null,
      logoPreview: null
  }
}

const handleLogoUpload = (event) => {
    const file = event.target.files[0]
    if (file) {
        newEnterprise.value.logo = file
        newEnterprise.value.logoPreview = URL.createObjectURL(file)
    }
}

const handleCreateEnterprise = async () => {
    isCreating.value = true
    try {
        let logoUrl = ''
        if (newEnterprise.value.logo && newEnterprise.value.logo instanceof File) {
            const formData = new FormData()
            formData.append('file', newEnterprise.value.logo)
            const { data } = await enterpriseAPI.uploadLogo(formData)
            logoUrl = data.url
        }

        const payload = {
            name: newEnterprise.value.name,
            employee_count_range: newEnterprise.value.employee_count_range,
            logo: logoUrl
        }

        await enterpriseAPI.create(payload)
        showCreateModal.value = false
        await fetchEnterprises()
    } catch (error) {
        console.error('Failed to create', error)
        alert('Erreur lors de la création')
    } finally {
        isCreating.value = false
    }
}

// Settings Logic
const addClass = () => {
    if (!currentEnterprise.value.school_config) currentEnterprise.value.school_config = { classes: [] }
    currentEnterprise.value.school_config.classes.push({ name: '', total_fees: 0, tranches: [] })
}

const removeClass = (index) => {
    currentEnterprise.value.school_config.classes.splice(index, 1)
}

const addRoute = () => {
     if (!currentEnterprise.value.transport_config) currentEnterprise.value.transport_config = { routes: [], zones: [] }
     currentEnterprise.value.transport_config.routes.push({ name: '', base_price: 0 })
}

const removeRoute = (index) => {
    currentEnterprise.value.transport_config.routes.splice(index, 1)
}

const addTranche = (cls) => {
    if (!cls.tranches) cls.tranches = []
    cls.tranches.push({ name: '', amount: 0, due_date: '' })
}

const removeTranche = (cls, index) => {
    cls.tranches.splice(index, 1)
}

const toggleSchedule = (svc) => {
    svc.use_schedule = !svc.use_schedule
    if (svc.use_schedule && !svc.payment_schedule) {
        svc.payment_schedule = []
    }
}

const addScheduleItem = (svc) => {
    if (!svc.payment_schedule) svc.payment_schedule = []
    svc.payment_schedule.push({ name: '', start_date: '', end_date: '', amount: 0 })
}

const removeScheduleItem = (svc, idx) => {
    svc.payment_schedule.splice(idx, 1)
}

// Form Builder Logic
const addFormField = (svc) => {
    if (!svc.form_schema) svc.form_schema = []
    svc.form_schema.push({ key: '', label: '', type: 'text', required: false })
}

const removeFormField = (svc, idx) => {
    svc.form_schema.splice(idx, 1)
}

const togglePenalty = (svc) => {
    if (svc.penalty_config) {
        svc.penalty_config = null
    } else {
        svc.penalty_config = {
            type: 'PERCENTAGE',
            value: 10,
            frequency: 'ONETIME',
            grace_period: 5
        }
    }
}

// Auto-Calculate Total for Custom Schedules
// Auto-Calculate Total for Custom Schedules
watch(() => currentEnterprise.value?.service_groups, (groups) => {
    if (!groups) return
    groups.forEach(group => {
        if (!group.services) return
        group.services.forEach(svc => {
            if (svc.billing_frequency === 'CUSTOM' && svc.use_schedule && svc.payment_schedule?.length) {
                const total = svc.payment_schedule.reduce((sum, item) => sum + (parseFloat(item.amount) || 0), 0)
                if (total > 0 && total !== svc.base_price) {
                    svc.base_price = total
                }
            }
        })
    })
}, { deep: true })

const saveSettings = async () => {
    if (!currentEnterprise.value) return
    isSaving.value = true
    try {
        // ... (existing save logic in component, handled by updateEnterprise call in template usually? No, the saveSettings function was cut off in previous view)
        // Re-implementing save logic or just ensuring it precedes new methods?
        // Ah, saveSettings starts at 468.
        await enterpriseAPI.update(currentEnterprise.value.id, currentEnterprise.value)
        fetchEnterprises() // Refresh list
        // Close ? No, just notify
        alert('Sauvegardé avec succès')
    } catch (error) {
        console.error('Failed to save settings', error)
        alert('Erreur lors de la sauvegarde')
    } finally {
        isSaving.value = false
    }
}

// Manual Entry State
const billingMode = ref('IMPORT') // 'IMPORT' | 'MANUAL'
const manualSubscribers = ref([])
const manualEntries = ref({}) // subId -> { amount: 0, consumption: 0 }
const isGeneratingManual = ref(false)

const loadManualSubscribers = async () => {
    if (!currentEnterprise.value || !selectedImportService.value) return
    try {
        isLoading.value = true // Reuse global loader or local?
        const { data } = await enterpriseAPI.getSubscriptions(currentEnterprise.value.id, selectedImportService.value)
        manualSubscribers.value = data
        
        // Initialize entries
        const initial = {}
        data.forEach(sub => {
            // Check if existing entry? No, reset for new batch
             initial[sub.id] = { amount: 0, consumption: 0 }
        })
        manualEntries.value = initial
    } catch (e) {
        console.error('Failed to load subs', e)
    } finally {
        isLoading.value = false
    }
}

const submitManualBatch = async () => {
    if (!confirm('Voulez-vous générer ces factures ?')) return
    isGeneratingManual.value = true
    try {
        const items = Object.entries(manualEntries.value).map(([subId, val]) => ({
            subscription_id: subId,
            amount: parseFloat(val.amount) || 0,
            consumption: parseFloat(val.consumption) || 0
        })).filter(i => i.amount > 0 || i.consumption > 0)

        if (items.length === 0) {
            alert('Aucune donnée saisie')
            return
        }

        await enterpriseAPI.createBatchInvoices(currentEnterprise.value.id, items)
        alert('Factures générées avec succès (En attente de validation)')
        // Reset or leave? Leave for now.
    } catch (e) {
        console.error(e)
        alert('Erreur lors de la génération')
    } finally {
        isGeneratingManual.value = false
    }
}

const addServiceGroup = () => {
    if (!currentEnterprise.value.service_groups) currentEnterprise.value.service_groups = []
    currentEnterprise.value.service_groups.push({
        id: crypto.randomUUID(),
        name: 'Nouveau Groupe',
        is_private: false,
        currency: currentEnterprise.value.settings?.currency || 'XOF',
        services: []
    })
}

const removeServiceGroup = (idx) => {
    if (confirm('Supprimer ce groupe et tous ses services ?')) {
        currentEnterprise.value.service_groups.splice(idx, 1)
    }
}

const addCustomService = (group) => {
    if (!group.services) group.services = []
    group.services.push({ 
        id: '', 
        name: '', 
        // category removed
        billing_type: 'FIXED', 
        billing_frequency: 'MONTHLY', 
        base_price: 0, 
        unit: '',
        uid: Date.now() 
    })
}

const removeCustomService = (group, svc) => {
    const idx = group.services.indexOf(svc)
    if (idx > -1) {
        group.services.splice(idx, 1)
    }
}

// GroupedServices computed is removed as we iterate state directly

const handleFileUpload = (event) => {
    importFile.value = event.target.files[0]
}

const uploadInvoiceFile = async () => {
    if (!importFile.value || !selectedImportService.value) {
        alert('Veuillez sélectionner un service et un fichier')
        return
    }
    
    const formData = new FormData()
    formData.append('file', importFile.value)
    formData.append('service_id', selectedImportService.value)
    
    // Default Mapping (User should arguably configure this, but hardcoding for MVP based on UI hint)
    // 0: Client ID, 1: Amount, 2: Consumption
    formData.append('col_client_idx', "0") 
    formData.append('col_amount_idx', "1") 
    formData.append('col_consumption_idx', "2")
    
    try {
        await enterpriseAPI.importInvoices(currentEnterprise.value.id, formData)
         alert('Import réussi ! Les factures ont été générées en brouillon.')
    } catch (error) {
         console.error('Upload failed', error)
         alert("Erreur lors de l'import: " + (error.response?.data?.error || error.message))
    }
}

// Add Service Modal Logic
const showAddServiceModal = ref(false)
const newService = ref({
    group_id: '',
    name: '',
    billing_type: 'FIXED', 
    billing_frequency: 'MONTHLY', 
    base_price: 0, 
    unit: ''
})

const goToSettingsAndAddService = () => {
    // If no groups, user must create one first
    if (!currentEnterprise.value.service_groups || currentEnterprise.value.service_groups.length === 0) {
        alert("Veuillez d'abord créer un groupe de services (ex: 'Scolarité', 'Transport') dans l'onglet Paramètres.")
        currentTab.value = 'Settings'
        return
    }
    
    // Open Modal
    newService.value = {
        group_id: currentEnterprise.value.service_groups[0].id, // Default to first
        name: '',
        billing_type: 'FIXED', 
        billing_frequency: 'MONTHLY', 
        base_price: 0, 
        unit: ''
    }
    showAddServiceModal.value = true
}

const confirmAddService = () => {
    if (!newService.value.group_id || !newService.value.name) {
        alert('Veuillez remplir le nom et choisir un groupe')
        return
    }
    
    const group = currentEnterprise.value.service_groups.find(g => g.id === newService.value.group_id)
    if (!group) return
    
    if (!group.services) group.services = []
    
    group.services.push({
        id: newService.value.name.toLowerCase().replace(/\s+/g, '_'), 
        name: newService.value.name,
        billing_type: newService.value.billing_type,
        billing_frequency: newService.value.billing_frequency,
        base_price: newService.value.base_price,
        unit: newService.value.unit,
        uid: Date.now()
    })
    
    showAddServiceModal.value = false
    alert('Service ajouté ! N\'oubliez pas d\'enregistrer les modifications.')
    currentTab.value = 'Settings' // Go to settings so they can save
}

// QR Modal Logic
const qrTab = ref('ENTERPRISE') // 'ENTERPRISE' | 'GROUP' | 'SERVICE'
const selectedQRGroup = ref('')
const selectedQRService = ref('')

const getCurrentQRLink = () => {
    if (!currentEnterprise.value) return ''
    const baseUrl = `/enterprise-service/api/v1/enterprises/${currentEnterprise.value.id}`
    
    if (qrTab.value === 'ENTERPRISE') {
        return `${baseUrl}/qrcode`
    } else if (qrTab.value === 'GROUP') {
        if (!selectedQRGroup.value) return ''
        return `${baseUrl}/groups/${selectedQRGroup.value}/qrcode`
    } else if (qrTab.value === 'SERVICE') {
        if (!selectedQRService.value) return ''
        return `${baseUrl}/services/${selectedQRService.value}/qrcode`
    }
    return ''
}

const getQRTitle = () => {
    if (qrTab.value === 'ENTERPRISE') return currentEnterprise.value?.name
    if (qrTab.value === 'GROUP') {
        return currentEnterprise.value?.service_groups?.find(g => g.id === selectedQRGroup.value)?.name || 'Sélectionner un groupe'
    }
    if (qrTab.value === 'SERVICE') {
         const group = currentEnterprise.value?.service_groups?.find(g => g.id === selectedQRGroup.value)
         return group?.services?.find(s => s.id === selectedQRService.value)?.name || 'Sélectionner un service'
    }
    return ''
}

const getQRDescription = () => {
    if (qrTab.value === 'ENTERPRISE') return "Code QR principal de l'entreprise. Scannez pour voir tous les services publics."
    if (qrTab.value === 'GROUP') return "Code QR pour ce groupe de services. Scannez pour voir uniquement les services de cette catégorie."
    if (qrTab.value === 'SERVICE') return "Code QR direct pour ce service. Scannez pour souscrire directement."
    return ''
}

const getQRLabel = () => {
    if (qrTab.value === 'ENTERPRISE') return `ENT-${currentEnterprise.value?.id}`
    if (qrTab.value === 'GROUP') return `GRP-${selectedQRGroup.value?.substring(0,8)}`
    if (qrTab.value === 'SERVICE') return `SVC-${selectedQRService.value}`
    return 'CODE'
}

watch(qrTab, () => { 
    // Reset selections on tab change if needed, or keep for convenience
    if (qrTab.value === 'GROUP' && !selectedQRGroup.value) {
        // Auto select first public group
        const grp = currentEnterprise.value.service_groups.find(g => !g.is_private)
        if (grp) selectedQRGroup.value = grp.id
    }
})

watch(selectedQRGroup, () => {
     if (qrTab.value === 'SERVICE') {
         selectedQRService.value = '' // Reset service when group changes
     }
})

// Watchers moved to end to avoid initialization issues
watch(currentTab, (newTab) => {
   if (newTab === 'Employees') {
      fetchEmployees()
   }
   if (newTab === 'Clients') {
      fetchClients()
   }
})

// Initial fetch if already on tab
watch(currentEnterprise, (newEnt) => {
   if (newEnt && currentTab.value === 'Employees') {
      fetchEmployees()
   }
   if (newEnt && currentTab.value === 'Clients') {
       fetchClients()
   }
})
</script>
