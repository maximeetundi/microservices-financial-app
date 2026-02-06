<template>
  <NuxtPage />
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>
