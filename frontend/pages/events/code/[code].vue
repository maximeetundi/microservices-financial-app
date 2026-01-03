<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-xl mx-auto py-8 px-4">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-20">
        <div class="loading-spinner w-12 h-12 mx-auto mb-4"></div>
        <p class="text-muted">Recherche de l'événement...</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-20">
        <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-error/10 flex items-center justify-center">
          <span class="text-5xl">❌</span>
        </div>
        <h1 class="text-2xl font-bold text-base mb-2">Événement non trouvé</h1>
        <p class="text-muted mb-6">Le code "{{ eventCode }}" ne correspond à aucun événement.</p>
        <NuxtLink to="/events" class="btn-premium inline-block px-6 py-3">
          Voir tous les événements
        </NuxtLink>
      </div>

      <!-- Success - Redirect message (user will be redirected automatically) -->
      <div v-else-if="event" class="text-center py-20">
        <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-success/10 flex items-center justify-center">
          <span class="text-5xl">✓</span>
        </div>
        <h1 class="text-2xl font-bold text-base mb-2">Événement trouvé!</h1>
        <p class="text-muted mb-4">{{ event.title }}</p>
        <p class="text-sm text-muted">Redirection en cours...</p>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ticketAPI } from '~/composables/useApi'

definePageMeta({
  layout: false
})

const route = useRoute()
const router = useRouter()
const eventCode = route.params.code

const loading = ref(true)
const error = ref(false)
const event = ref(null)

onMounted(async () => {
  try {
    // Search for event by code
    const res = await ticketAPI.getEventByCode(eventCode)
    event.value = res.data?.event || res.data
    
    if (event.value?.id) {
      // Redirect to event detail page
      setTimeout(() => {
        router.push(`/events/${event.value.id}`)
      }, 500)
    } else {
      error.value = true
    }
  } catch (e) {
    console.error('Event lookup failed:', e)
    error.value = true
  } finally {
    loading.value = false
  }
})
</script>
