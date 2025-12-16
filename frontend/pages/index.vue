<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center">
    <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
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

