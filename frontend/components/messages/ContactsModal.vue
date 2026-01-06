<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60" @click.self="close">
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-lg max-h-[80vh] flex flex-col overflow-hidden">
        <!-- Header -->
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between bg-gradient-to-r from-green-500 to-teal-600">
          <h2 class="text-xl font-bold text-white">Mes Contacts</h2>
          <button @click="close" class="text-white/80 hover:text-white transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Add Contact Form -->
        <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
          <div class="flex gap-2">
            <input 
              v-model="newContact.phone" 
              type="tel" 
              placeholder="Téléphone" 
              class="flex-1 px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm"
            />
            <input 
              v-model="newContact.name" 
              type="text" 
              placeholder="Nom" 
              class="flex-1 px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm"
            />
            <button 
              @click="addContact" 
              :disabled="!newContact.name || (!newContact.phone && !newContact.email)"
              class="px-4 py-2 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 text-white rounded-lg text-sm font-medium transition-colors"
            >
              Ajouter
            </button>
          </div>
          <input 
            v-model="newContact.email" 
            type="email" 
            placeholder="Email (optionnel)" 
            class="mt-2 w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm"
          />
        </div>

        <!-- Contacts List -->
        <div class="flex-1 overflow-y-auto">
          <div v-if="loading" class="flex items-center justify-center py-12">
            <div class="animate-spin w-8 h-8 border-2 border-green-500 border-t-transparent rounded-full"></div>
          </div>

          <div v-else-if="contacts.length === 0" class="flex flex-col items-center justify-center py-12 text-gray-500">
            <svg class="w-16 h-16 mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <p class="font-medium">Aucun contact</p>
            <p class="text-sm">Ajoutez des contacts manuellement ou synchronisez depuis votre mobile</p>
          </div>

          <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
            <div 
              v-for="contact in contacts" 
              :key="contact.id" 
              class="flex items-center gap-4 px-6 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              <div class="w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center font-bold text-lg">
                {{ contact.name[0]?.toUpperCase() }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="font-medium text-gray-900 dark:text-white truncate">{{ contact.name }}</p>
                <p class="text-sm text-gray-500 truncate">{{ contact.phone || contact.email }}</p>
              </div>
              <button 
                @click="startConversation(contact)" 
                class="p-2 text-green-600 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-full transition-colors"
                title="Démarrer une conversation"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </button>
              <button 
                @click="deleteContact(contact.id)" 
                class="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-full transition-colors"
                title="Supprimer"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900 text-center text-sm text-gray-500">
          {{ contacts.length }} contact(s) • Synchro depuis l'app mobile
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { contactsAPI } from '~/composables/useApi'
import api from '~/composables/useApi'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits(['close', 'startConversation'])

interface Contact {
  id: string
  phone: string
  email: string
  name: string
}

const contacts = ref<Contact[]>([])
const loading = ref(false)
const newContact = ref({
  phone: '',
  email: '',
  name: ''
})

const close = () => emit('close')

const loadContacts = async () => {
  loading.value = true
  try {
    const res = await contactsAPI.getAll()
    contacts.value = res.data?.contacts || []
  } catch (e) {
    console.error('Failed to load contacts:', e)
  } finally {
    loading.value = false
  }
}

const addContact = async () => {
  if (!newContact.value.name) return
  if (!newContact.value.phone && !newContact.value.email) return
  
  try {
    await contactsAPI.add(newContact.value)
    newContact.value = { phone: '', email: '', name: '' }
    await loadContacts()
  } catch (e: any) {
    alert(e.response?.data?.error || 'Erreur lors de l\'ajout')
  }
}

const deleteContact = async (id: string) => {
  if (!confirm('Supprimer ce contact ?')) return
  try {
    await contactsAPI.delete(id)
    contacts.value = contacts.value.filter(c => c.id !== id)
  } catch (e) {
    console.error('Failed to delete contact:', e)
  }
}

const startConversation = async (contact: Contact) => {
  try {
    // Search for user by phone/email
    const searchQuery = contact.phone || contact.email
    const res = await api.get('/auth-service/api/v1/users/search', {
      params: { q: searchQuery }
    })
    
    const users = res.data?.users || []
    if (users.length > 0) {
      emit('startConversation', {
        ...users[0],
        contactName: contact.name // Pass the contact name to use
      })
    } else {
      alert('Cet utilisateur n\'est pas encore inscrit sur la plateforme')
    }
  } catch (e) {
    console.error('Failed to start conversation:', e)
  }
}

watch(() => props.show, (val) => {
  if (val) loadContacts()
})

onMounted(() => {
  if (props.show) loadContacts()
})
</script>
