<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-bold">Tontine Téléphonée</h3>
      <button @click="showCreate = true" class="btn-primary">+ Nouveau tour</button>
    </div>

    <div class="space-y-3">
      <div v-for="round in rounds" :key="round.id"
        class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div class="flex justify-between items-start mb-3">
          <div>
            <div class="font-medium">Tour #{{ round.round_number }}</div>
            <div class="text-sm text-gray-500">Bénéficiaire: {{ round.beneficiary_name }}</div>
          </div>
          <span :class="['px-2 py-1 rounded text-xs', round.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100']">
            {{ round.status }}
          </span>
        </div>

        <div class="mb-3">
          <div class="text-lg font-bold text-indigo-600">{{ formatAmount(round.total_collected) }}</div>
          <div class="text-sm text-gray-500">Total collecté</div>
        </div>

        <div v-if="round.status === 'active'" class="flex gap-2">
          <button @click="viewPledges(round.id)" class="flex-1 btn-secondary text-sm">Voir les promesses</button>
          <button @click="makePledge(round.id)" class="flex-1 btn-primary text-sm">Promettre</button>
        </div>
      </div>
    </div>

    <!-- Create Round Modal -->
    <div v-if="showCreate" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold mb-4">Nouveau tour</h3>

        <div class="space-y-4">
          <div>
            <label class="block font-medium mb-2">Bénéficiaire (qui "bouffe")</label>
            <select v-model="newRound.beneficiary_id" class="input w-full">
              <option v-for="m in members" :key="m.id" :value="m.id">{{ m.user_id }}</option>
            </select>
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button @click="showCreate = false" class="flex-1 btn-secondary">Annuler</button>
          <button @click="createRound" class="flex-1 btn-primary">Créer</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { associationAPI } from '~/composables/useApi'

const props = defineProps<{
  associationId: string
}>()

const rounds = ref<any[]>([])
const members = ref<any[]>([])
const showCreate = ref(false)

const newRound = ref({
  beneficiary_id: '',
})

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'XOF' }).format(amount || 0)
}

const loadRounds = async () => {
  try {
    const res = await associationAPI.getCalledRounds(props.associationId)
    rounds.value = res.data
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

const createRound = async () => {
  try {
    await associationAPI.createCalledRound(props.associationId, newRound.value.beneficiary_id)
    showCreate.value = false
    loadRounds()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur')
  }
}

const makePledge = (roundId: string) => {
  const amount = prompt('Montant de votre promesse:')
  if (amount) {
    associationAPI.makePledge(roundId, parseFloat(amount))
      .then(() => loadRounds())
      .catch((err: any) => alert(err.response?.data?.error || 'Erreur'))
  }
}

const viewPledges = (roundId: string) => {
  // Navigate to pledges view or open modal
  alert('Vue des promesses (à implémenter)')
}

onMounted(() => {
  loadRounds()
  loadMembers()
})
</script>
