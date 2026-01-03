<template>
  <NuxtLayout name="dashboard">
    <div class="flex items-center justify-center min-h-[50vh]">
      <div class="loading-spinner"></div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'auth'
})

const route = useRoute()

// Redirect to the main event page
onMounted(() => {
  const eventId = route.params.id as string
  navigateTo(`/events/${eventId}`, { replace: true })
})
</script>

<style scoped>
.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
