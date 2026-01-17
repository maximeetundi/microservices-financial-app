<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
      <p class="text-gray-500">Redirection...</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// Redirect to generic approve page
onMounted(() => {
  const employeeId = route.query.id
  if (employeeId) {
    router.replace(`/approve/${employeeId}?type=invitation`)
  } else {
    router.replace('/enterprise')
  }
})
</script>
