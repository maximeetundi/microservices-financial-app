<template>
  <ShopLayout>
    <NuxtPage />
  </ShopLayout>
</template>

<script setup lang="ts">
import ShopLayout from '@/components/shops/ShopLayout.vue'

definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>
