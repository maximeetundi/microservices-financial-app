<template>
  <NuxtLayout name="enterprise">
    <ServicesTab :enterprise="enterprise" @update="handleUpdate" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const ServicesTab = defineAsyncComponent(() => import('@/components/enterprise/ServicesTab.vue'))

const enterprise = inject('enterprise', ref(null))

const handleUpdate = (updated) => {
  enterprise.value = updated
}
</script>
