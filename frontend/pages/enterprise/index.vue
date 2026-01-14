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
         
         <!-- Clients Tab -->
         <div v-if="currentTab === 'Clients'" class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
            <div class="flex justify-between mb-4">
                <h3 class="font-semibold text-lg dark:text-white">Abonnés & Élèves</h3>
                <div class="flex gap-2">
                    <button @click="downloadExport" class="px-3 py-1 border border-gray-300 text-gray-700 rounded-md text-sm hover:bg-gray-50 flex items-center gap-2">
                        <span>📥</span> Exporter CSV
                    </button>
                    <button @click="showAddClientModal = true" class="px-3 py-1 bg-purple-600 text-white rounded-md text-sm">+ Ajouter un Client</button>
                </div>
            </div>

            <!-- Clients List -->
             <!-- We reuse manualSubscribers style list, or fetch distinct -->
             <div v-if="isLoading" class="text-center py-8 text-gray-400">Chargement...</div>
             <div v-else>
                <div class="mb-4 flex gap-2">
                    <select v-model="selectedClientFilterService" @change="fetchClients" class="border rounded px-2 py-1 text-sm">
                        <option value="">Tous les services</option>
                        <option v-for="svc in currentEnterprise.custom_services" :key="svc.id" :value="svc.id">{{ svc.name }}</option>
                    </select>
                </div>
                <table class="w-full text-left text-sm text-gray-500">
                    <thead class="bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                        <tr>
                            <th class="px-4 py-3">Nom du Client</th>
                            <th class="px-4 py-3">Matricule / ID</th>
                            <th class="px-4 py-3">Service / Classe</th>
                            <th class="px-4 py-3">Fréquence/Montant</th>
                            <th class="px-4 py-3">Proch. Facturation</th>
                            <th class="px-4 py-3">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="sub in clientSubscriptions" :key="sub.id" class="border-b dark:border-gray-700">
                            <td class="px-4 py-3 font-medium">{{ sub.client_name }}</td>
                            <td class="px-4 py-3 text-xs font-mono bg-gray-50 dark:bg-gray-900 rounded px-1">
                                {{ sub.external_id || '-' }}
                            </td>
                             <td class="px-4 py-3">
                                 <!-- Resolving Service Name -->
                                 {{ getServiceName(sub.service_id) }}
                                 <span v-if="sub.school_details" class="text-xs text-gray-400 block">
                                     {{ sub.school_details.student_name }} ({{ sub.school_details.class_id }})
                                 </span>
                             </td>
                             <td class="px-4 py-3">
                                 {{ sub.billing_frequency }} - {{ sub.amount }} XOF
                             </td>
                             <td class="px-4 py-3">
                                 {{ new Date(sub.next_billing_at).toLocaleDateString() }}
                             </td>
                             <td class="px-4 py-3">
                                 <button class="text-xs text-red-500 hover:underline">Résilier</button>
                             </td>
                        </tr>
                        <tr v-if="clientSubscriptions.length === 0">
                            <td colspan="6" class="text-center py-8 text-gray-400">Aucun abonné trouvé.</td>
                        </tr>
                    </tbody>
                </table>
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

            <!-- Generic/Custom Services (Utilities, etc.) -->
            <div class="mb-8 p-4 border rounded-lg dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                <h4 class="font-medium mb-4 dark:text-gray-200 flex items-center gap-2">
                    <span>⚡ Services Personnalisés (Eau, Élec, Abonnements)</span>
                </h4>
                <div class="mb-6">
                    <div class="flex justify-between items-center mb-2">
                        <label class="text-sm font-medium dark:text-gray-300">Services Définis</label>
                        <button @click="addCustomService" class="text-xs bg-gray-600 text-white px-2 py-1 rounded hover:bg-gray-700">+ Ajouter Service</button>
                    </div>
                    <div v-if="!currentEnterprise.custom_services?.length" class="text-sm text-gray-500 italic">Aucun service configuré.</div>
                     <div v-else class="space-y-4">
                        <div v-for="(svc, idx) in currentEnterprise.custom_services" :key="idx" class="flex flex-col gap-2 bg-white dark:bg-gray-800 p-3 rounded shadow-sm border border-gray-200 dark:border-gray-600">
                             <div class="flex gap-2">
                                <input v-model="svc.id" placeholder="ID (ex: water_usage)" class="w-1/4 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border font-mono">
                                <input v-model="svc.name" placeholder="Nom (ex: Consommation Eau)" class="flex-1 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                <button @click="removeCustomService(idx)" class="text-red-500 hover:text-red-700 text-sm px-2">Suppr.</button>
                             </div>
                             <div class="flex gap-2 items-center">
                                <select v-model="svc.billing_type" class="w-1/4 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <option value="FIXED">Fixe / Forfait</option>
                                    <option value="USAGE">À l'usage (Compteur)</option>
                                </select>
                                <select v-model="svc.billing_frequency" class="w-1/4 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                    <option value="DAILY">Journalier</option>
                                    <option value="WEEKLY">Hebdomadaire</option>
                                    <option value="MONTHLY">Mensuel</option>
                                    <option value="ANNUALLY">Annuel</option>
                                    <option value="CUSTOM">Pér. Personnalisée</option>
                                    <option value="ONETIME">Ponctuel</option>
                                </select>
                                <input v-if="svc.billing_frequency === 'CUSTOM' && !svc.use_schedule" v-model.number="svc.custom_interval" type="number" placeholder="Jours" class="w-16 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                <input v-if="svc.billing_type === 'USAGE'" v-model="svc.unit" placeholder="Unité (ex: kWh)" class="w-20 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                <input v-model.number="svc.base_price" type="number" step="0.01" :placeholder="svc.billing_type === 'USAGE' ? 'Prix Unitaire' : 'Montant Fixe'" class="w-32 rounded-md border-gray-300 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                <span class="text-xs text-gray-500">{{ currentEnterprise.settings?.currency || 'XOF' }}</span>
                             </div>

                             <!-- Payment Schedule UI (Custom Frequency) -->
                             <div v-if="svc.billing_frequency === 'CUSTOM'" class="mt-2 ml-2 pl-4 border-l-2 border-orange-200 dark:border-orange-800">
                                   <div class="flex items-center gap-2 mb-2">
                                       <label class="text-xs font-semibold text-gray-500 uppercase">Calendrier de Paiement</label>
                                       <div class="flex items-center gap-2 ml-4">
                                           <button @click="toggleSchedule(svc)" :class="{'bg-orange-100 text-orange-700': svc.use_schedule, 'text-gray-400': !svc.use_schedule}" class="text-xs px-2 py-0.5 rounded border border-gray-200 dark:border-gray-600">
                                               {{ svc.use_schedule ? 'Mode: Calendrier (Dates)' : 'Mode: Intervalle (Jours)' }}
                                           </button>
                                       </div>
                                   </div>
                                   
                                   <div v-if="svc.use_schedule">
                                         <div v-for="(item, sIdx) in svc.payment_schedule" :key="sIdx" class="flex items-center gap-2 mb-1">
                                            <input v-model="item.name" placeholder="Nom (ex: Session 1)" class="flex-1 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <input v-model="item.start_date" type="date" class="w-28 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <span class="text-xs text-gray-400">au</span>
                                            <input v-model="item.end_date" type="date" class="w-28 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <input v-model.number="item.amount" type="number" placeholder="Montant" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-700 dark:border-gray-600 dark:text-white px-2 py-1 border">
                                            <button @click="removeScheduleItem(svc, sIdx)" class="text-red-400 hover:text-red-600 text-xs font-bold px-1">×</button>
                                         </div>
                                         <button @click="addScheduleItem(svc)" class="text-xs text-blue-600 hover:text-blue-800 mt-1">+ Ajouter une période</button>
                                   </div>
                             </div>

                             <!-- Custom Form Builder (Client Side) -->
                             <div class="mt-2 ml-2 pl-4 border-l-2 border-indigo-200 dark:border-indigo-800">
                                   <div class="flex items-center gap-2 mb-2">
                                       <label class="text-xs font-semibold text-gray-500 uppercase">Formulaire d'Inscription (Optionnel)</label>
                                       <button @click="addFormField(svc)" class="text-xs bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded border border-indigo-200">+ Ajouter Champ</button>
                                   </div>
                                   <div v-if="svc.form_schema?.length">
                                         <div v-for="(field, fIdx) in svc.form_schema" :key="fIdx" class="flex items-center gap-2 mb-1 bg-gray-50 dark:bg-gray-700 p-1 rounded">
                                            <input v-model="field.label" placeholder="Libellé (ex: Nom de l'enfant)" class="flex-1 rounded-md border-gray-300 text-xs dark:bg-gray-600 dark:border-gray-500 dark:text-white px-2 py-1 border">
                                            <input v-model="field.key" placeholder="Clé (ex: child_name)" class="w-24 rounded-md border-gray-300 text-xs dark:bg-gray-600 dark:border-gray-500 dark:text-white px-2 py-1 border font-mono">
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
                       <option value="SCHOOL">École / Éducation</option>
                       <option value="TRANSPORT">Transport</option>
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
// Explicit imports including computed
import { ref, onMounted, watch, computed } from 'vue' 
import { enterpriseAPI, useApi } from '@/composables/useApi'

const { authApi } = useApi()

// ... Tabs ...

// Client Management State
const showAddClientModal = ref(false)
const clientSubscriptions = ref([])
const selectedClientFilterService = ref('')
const clientSearchQuery = ref('')
const foundUser = ref(null)
const isSearchingUser = ref(false)
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
        clientSubscriptions.value = data.sort((a, b) => (a.client_name || '').localeCompare(b.client_name || ''))
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
  // Ensure nested objects exist for v-model binding
  if (!ent.settings) ent.settings = { currency: 'XOF', auto_pay_salaries: false }
  
  if (ent.type === 'SCHOOL' && !ent.school_config) ent.school_config = { classes: [] }
  if (ent.type === 'TRANSPORT' && !ent.transport_config) ent.transport_config = { routes: [], zones: [] }
  if (!ent.custom_services) ent.custom_services = [] // Ensure custom_services exists
  
  currentEnterprise.value = JSON.parse(JSON.stringify(ent)) // Deep copy to avoid mutating list directly
  selectedImportService.value = '' // Reset import selection
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



const openCreateModal = () => {
  showCreateModal.value = true
  newEnterprise.value = { name: '', registration_number: '', type: 'SME' }
}

const handleCreateEnterprise = async () => {
    isCreating.value = true
    try {
        await enterpriseAPI.create(newEnterprise.value)
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

const addCustomService = () => {
    if (!currentEnterprise.value.custom_services) currentEnterprise.value.custom_services = []
    currentEnterprise.value.custom_services.push({ id: '', name: '', billing_type: 'FIXED', billing_frequency: 'MONTHLY', base_price: 0, unit: '' })
}

const removeCustomService = (index) => {
    currentEnterprise.value.custom_services.splice(index, 1)
}

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
