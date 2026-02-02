<template>
  <NuxtLayout name="enterprise">
    <BillingTab :enterprise="enterprise" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const BillingTab = defineAsyncComponent(() => import('@/components/enterprise/BillingTab.vue'))

const enterprise = inject('enterprise', ref(null))
</script>
