<template>
  <NuxtLayout name="enterprise">
    <ClientsTab :enterprise="enterprise" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const ClientsTab = defineAsyncComponent(() => import('@/components/enterprise/ClientsTab.vue'))

const enterprise = inject('enterprise', ref(null))
</script>
