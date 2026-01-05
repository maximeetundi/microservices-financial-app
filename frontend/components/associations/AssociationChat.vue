<template>
  <div  class="bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700 h-[600px] flex flex-col">
    <div class="p-4 border-b border-gray-200 dark:border-gray-700">
      <h3 class="font-bold text-lg">Chat</h3>
    </div>

    <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3">
      <div v-for="msg in messages" :key="msg.id"
        :class="['flex', msg.sender_id === currentUserId ? 'justify-end' : 'justify-start']">
        <div :class="['max-w-[70%] rounded-lg p-3', msg.sender_id === currentUserId ? 'bg-indigo-600 text-white' : 'bg-gray-100 dark:bg-gray-800']">
          <div v-if="msg.sender_id !== currentUserId" class="text-xs font-medium mb-1">{{ msg.sender_name }}</div>
          <div>{{ msg.content }}</div>
          <div class="text-xs opacity-70 mt-1">{{ formatTime(msg.created_at) }}</div>
        </div>
      </div>
    </div>

    <div class="p-4 border-t border-gray-200 dark:border-gray-700">
      <div class="flex gap-2">
        <input v-model="newMessage" @keyup.enter="sendMessage" type="text" placeholder="Votre message..."
          class="flex-1 input" />
        <button @click="sendMessage" class="btn-primary">Envoyer</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { associationAPI } from '~/composables/useApi'

const props = defineProps<{
  associationId: string
}>()

const messages = ref<any[]>([])
const newMessage = ref('')
const currentUserId = ref('current-user') // TODO: get from auth
const messagesContainer = ref<HTMLElement>()

const loadMessages = async () => {
  try {
    const res = await associationAPI.getMessages(props.associationId)
    messages.value = res.data
    await nextTick()
    scrollToBottom()
  } catch (err) {
    console.error(err)
  }
}

const sendMessage = async () => {
  if (!newMessage.value.trim()) return

  try {
    await associationAPI.sendMessage(props.associationId, newMessage.value)
    newMessage.value = ''
    loadMessages()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erreur')
  }
}

const formatTime = (date: string) => {
  return new Date(date).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

onMounted(() => {
  loadMessages()
  setInterval(loadMessages, 5000) // Refresh every 5s
})
</script>
