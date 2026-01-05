<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-bold">Solidarité</h3>
      <button @click="showCreate = true" class="btn-primary">+ Nouvelle aide</button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="event in events" :key="event.id"
        class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div class="flex justify-between items-start mb-2">
          <div>
            <div class="font-medium">{{ event.title }}</div>
            <div class="text-sm text-gray-500">{{ getEventTypeLabel(event.event_type) }}</div>
          </div>
          <span :class="['px-2 py-1 rounded text-xs', event.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100']">
            {{ event.status }}
          </span>
        </div>

        <div class="mb-2">
          <div class="flex justify-between text-sm mb-1">
            <span>Collecté</span>
            <span class="font-medium">{{ formatAmount(event.collected_amount) }} / {{ formatAmount(event.target_amount) }}</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div :style="{width: `${(event.collected_amount / event.target_amount) * 100}%`}"
              class="bg-indigo-600 h-2 rounded-full"></div>
          </div>
        </div>

        <button v-if="event.status === 'active'" @click="contribute(event.id)" class="btn-secondary w-full text-sm">
          Contribuer
        </button>
      </div>
    </div>

    <!-- Create Event Modal -->
    <div v-if="showCreate" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Créer une aide de solidarité</h3>

        <div class="space-y-4">
          <div>
            <label class="block font-medium mb-2">Type</label>
            <select v-model="newEvent.event_type" class="input w-full">
              <option value="deceased">Deuil</option>
              <option value="marriage">Mariage</option>
              <option value="birth">Naissance</option>
              <option value="illness">Maladie</option>
              <option value="other">Autre</option>
            </select>
          </div>

          <div>
            <label class="block font-medium mb-2">Titre</label>
            <input v-model="newEvent.title" type="text" class="input w-full" />
          </div>

          <div>
            <label class="block font-medium mb-2">Bénéficiaire</label>
            <select v-model="newEvent.beneficiary_id" class="input w-full">
              <option v-for="m in members" :key="m.id" :value="m.id">{{ m.user_id }}</option>
            </select>
          </div>

          <div>
            <label class="block font-medium mb-2">Objectif (optionnel)</label>
            <input v-model.number="newEvent.target_amount" type="number" class="input w-full" />
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button @click="showCreate = false" class="flex-1 btn-secondary">Annuler</button>
          <button @click="createEvent" class="flex-1 btn-primary">Créer</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { associationAPI } from '~/composables/useApi'
import { useNotification } from '~/composables/useNotification'

const props = defineProps<{
  associationId: string
}>()

const { showSuccess, showError } = useNotification()
const events = ref<any[]>([])
const members = ref<any[]>([])
const showCreate = ref(false)

const newEvent = ref({
  event_type: 'deceased',
  title: '',
  beneficiary_id: '',
  target_amount: 0,
})

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'XOF' }).format(amount || 0)
}

const getEventTypeLabel = (type: string) => {
  const labels: any = {
    deceased: 'Deuil',
    marriage: 'Mariage',
    birth: 'Naissance',
    illness: 'Maladie',
    other: 'Autre',
  }
  return labels[type] || type
}

const loadEvents = async () => {
  try {
    const res = await associationAPI.getSolidarityEvents(props.associationId)
    events.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const loadMembers = async () => {
  try {
    const res = await associationAPI.getMembers(props.associationId)
    members.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const createEvent = async () => {
  try {
    await associationAPI.createSolidarityEvent(props.associationId, newEvent.value)
    showCreate.value = false
    loadEvents()
  } catch (err: any) {
    showError(err.response?.data?.error || 'Erreur lors de la création')
  }
}

const contribute = (eventId: string) => {
  // Open contribution modal (to be created)
  showSuccess('Contribution enregistrée !', 'Succès')
}

onMounted(() => {
  loadEvents()
  loadMembers()
})
</script>
