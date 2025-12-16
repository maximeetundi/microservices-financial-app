<template>
  <div></div>
</template>

<script setup>
// Redirect to dashboard or login based on auth status
const router = useRouter()

onMounted(() => {
  // Check if user is authenticated
  const token = useCookie('auth_token')
  
  if (token.value) {
    router.push('/dashboard')
  } else {
    router.push('/auth/login')
  }
})
</script>
