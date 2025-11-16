<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Global Loading -->
    <div v-if="pending" class="fixed inset-0 bg-white bg-opacity-80 flex items-center justify-center z-50">
      <div class="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
    </div>

    <!-- App Content -->
    <NuxtPage />

    <!-- Global Notifications -->
    <Teleport to="body">
      <div id="notifications-container"></div>
    </Teleport>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { onMounted } from 'vue'

const { $router } = useNuxtApp()
const authStore = useAuthStore()
const pending = ref(false)

// Initialize auth state
onMounted(async () => {
  pending.value = true
  try {
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Auth initialization failed:', error)
  } finally {
    pending.value = false
  }
})

// Global navigation guard
watch(() => $router.currentRoute.value, (to, from) => {
  // Check authentication for protected routes
  if (to.meta?.requiresAuth && !authStore.isAuthenticated) {
    navigateTo('/auth/login')
  }
})
</script>

<style>
html {
  scroll-behavior: smooth;
}

body {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}
</style>