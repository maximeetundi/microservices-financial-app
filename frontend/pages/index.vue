<template>
  <div class="min-h-screen bg-base flex items-center justify-center">
    <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

// If user reaches this page with auth, redirect to dashboard
const token = useCookie('accessToken')
const authToken = useCookie('auth_token')

if (token.value || authToken.value) {
  navigateTo('/dashboard')
}
</script>

