<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Paramètres</h1>
          <p class="text-gray-600 dark:text-gray-400">{{ association?.name }}</p>
        </div>
        <button @click="navigateTo(`/associations/${id}`)" class="btn-secondary">
          Retour
        </button>
      </div>

      <div class="grid grid-cols-4 gap-6">
        <!-- Sidebar -->
        <div class="col-span-1 space-y-1">
          <button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id"
            :class="['w-full text-left px-4 py-3 rounded-lg transition-colors', activeTab === tab.id ? 'bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 font-medium' : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800']">
            {{ tab.name }}
          </button>
        </div>

        <!-- Content -->
        <div class="col-span-3 bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700 p-6">
          <!-- Rôles Personnalisés -->
          <div v-if="activeTab === 'roles'">
            <div class="flex justify-between items-center mb-6">
              <h2 class="text-xl font-bold">Rôles personnalisés</h2>
              <button @click="showCreateRole = true" class="btn-primary">+ Nouveau rôle</button>
            </div>

            <div class="space-y-3">
              <div v-for="role in roles" :key="role.id"
                class="flex justify-between items-center p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">{{ role.name }}</div>
                  <div class="text-sm text-gray-500">{{ role.permissions.length}} permissions</div>
                </div>
                <button v-if="!role.is_default" @click="deleteRole(role.id)" class="text-red-600 hover:text-red-700">
                  Supprimer
                </button>
              </div>
            </div>
          </div>

          <!-- Approuveurs Multi-Signature -->
          <div v-else-if="activeTab === 'approvers'">
            <h2 class="text-xl font-bold mb-6">Approuveurs (4/5 requis)</h2>
            <p class="text-gray-600 mb-4">Sélectionnez exactement 5 membres qui auront le pouvoir d'approuver les prêts et distributions.</p>

            <div class="space-y-3">
              <div v-for="(member, i) in members" :key="member.id"
                class="flex justify-between items-center p-4 border border-gray-200 dark:border-gray-700 rounded-lg cursor-pointer hover:bg-gray-50"
                @click="toggleApprover(member.id)">
                <div class="flex items-center space-x-3">
                  <input type="checkbox" :checked="selectedApprovers.includes(member.id)" class="w-5 h-5"
                    @click.stop="toggleApprover(member.id)" />
                  <div>
                    <div class="font-medium">{{ member.user_id }}</div>
                    <div class="text-sm text-gray-500">{{ member.role }}</div>
                  </div>
                </div>
              </div>
            </div>

            <button @click="saveApprovers" :disabled="selectedApprovers.length !== 5"
              class="mt-6 btn-primary disabled:opacity-50">
              Enregistrer les approuve urs ({{ selectedApprovers.length }}/5)
            </button>
          </div>

          <!-- Paramètres Généraux -->
          <div v-else-if="activeTab === 'general'">
            <h2 class="text-xl font-bold mb-6">Paramètres généraux</h2>

            <div class="space-y-6">
              <div>
                <label class="block font-medium mb-2">Frais de retard (montant fixe)</label>
                <input v-model.number="settings.late_fee_amount" type="number" class="input w-full" />
              </div>

              <div>
                <label class="block font-medium mb-2">Frais de nourriture</label>
                <input v-model.number="settings.food_fee" type="number" class="input w-full" />
              </div>

              <div>
                <label class="block font-medium mb-2">Frais de boisson</label>
                <input v-model.number="settings.drink_fee" type="number" class="input w-full" />
              </div>

              <div>
                <label class="flex items-center space-x-3">
                  <input v-model="chatEnabled" type="checkbox" class="w-5 h-5" />
                  <span>Activer le chat d'association</span>
                </label>
              </div>

              <div v-if="chatEnabled">
                <label class="flex items-center space-x-3">
                  <input v-model="chatAdminOnly" type="checkbox" class="w-5 h-5" />
                  <span>Restreindre le chat aux admins uniquement</span>
                </label>
              </div>

              <button @click="saveSettings" class="btn-primary">Enregistrer</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Role Modal -->
    <div v-if="showCreateRole" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Créer un rôle</h3>

        <div class="space-y-4">
          <div>
            <label class="block font-medium mb-2">Nom du rôle</label>
            <input v-model="newRole.name" type="text" class="input w-full" />
          </div>

          <div>
            <label class="block font-medium mb-2">Permissions</label>
            <div class="space-y-2">
              <label v-for="perm in availablePermissions" :key="perm.value" class="flex items-center space-x-2">
                <input v-model="newRole.permissions" :value="perm.value" type="checkbox" class="w-4 h-4" />
                <span class="text-sm">{{ perm.label }}</span>
              </label>
            </div>
          </div>
        </div>

        <div class="flex space-x-3 mt-6">
          <button @click="showCreateRole = false" class="flex-1 btn-secondary">Annuler</button>
          <button @click="createRole" class="flex-1 btn-primary">Créer</button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { associationAPI } from '~/composables/useApi'

const route = useRoute()
const id = route.params.id as string

const activeTab = ref('roles')
const association = ref<any>(null)
const roles = ref<any[]>([])
const members = ref<any[]>([])
const selectedApprovers = ref<string[]>([])
const showCreateRole = ref(false)
const chatEnabled = ref(true)
const chatAdminOnly = ref(false)

const settings = ref({
  late_fee_amount: 0,
  food_fee: 2000,
  drink_fee: 2000,
})

const tabs = [
  { id: 'roles', name: 'Rôles' },
  { id: 'approvers', name: 'Approuveurs' },
  { id: 'general', name: 'Général' },
]

const availablePermissions = [
  { value: 'manage_members', label: 'Gérer les membres' },
  { value: 'manage_meetings', label: 'Gérer les réunions' },
  { value: 'approve_loans', label: 'Approuver les prêts' },
  { value: 'approve_distributions', label: 'Approuver les distributions' },
  { value: 'view_treasury', label: 'Voir la trésorerie' },
  { value: 'record_contributions', label: 'Enregistrer les cotisations' },
  { value: 'manage_roles', label: 'Gérer les rôles' },
  { value: 'manage_chat', label: 'Gérer le chat' },
]

const newRole = ref({
  name: '',
  permissions: [] as string[],
})

const loadRoles = async () => {
  try {
    const res = await associationAPI.getRoles(id)
    roles.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const loadMembers = async () => {
  try {
    const res = await associationAPI.getMembers(id)
    members.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const loadApprovers = async () => {
  try {
    const res = await associationAPI.getApprovers(id)
    selectedApprovers.value = res.data.map((a: any) => a.member_id)
  } catch (err) {
    console.error(err)
  }
}

const createRole = async () => {
  try {
    await associationAPI.createRole(id, newRole.value)
    showCreateRole.value = false
    newRole.value = { name: '', permissions: [] }
    loadRoles()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur')
  }
}

const deleteRole = async (roleId: string) => {
  if (confirm('Supprimer ce rôle ?')) {
    try {
      await associationAPI.deleteRole(id, roleId)
      loadRoles()
    } catch (err: any) {
      alert(err.response?.data?.error || 'Erreur')
    }
  }
}

const toggleApprover = (memberId: string) => {
  const index = selectedApprovers.value.indexOf(memberId)
  if (index > -1) {
    selectedApprovers.value.splice(index, 1)
  } else if (selectedApprovers.value.length < 5) {
    selectedApprovers.value.push(memberId)
  }
}

const saveApprovers = async () => {
  try {
    await associationAPI.setApprovers(id, selectedApprovers.value)
    alert('Approuveurs enregistrés !')
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur')
  }
}

const saveSettings = async () => {
  try {
    await associationAPI.update(id, { settings: settings.value, chat_enabled: chatEnabled.value, chat_admin_only: chatAdminOnly.value })
    alert('Paramètres enregistrés !')
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur')
  }
}

onMounted(async () => {
  const res = await associationAPI.get(id)
  association.value = res.data
  loadRoles()
  loadMembers()
  loadApprovers()
})
</script>
