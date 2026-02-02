<template>
  <NuxtLayout name="enterprise">
    <SecurityTab :enterprise="enterprise" @update="handleUpdate" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const SecurityTab = defineAsyncComponent(() => import('@/components/enterprise/SecurityTab.vue'))

const enterprise = inject('enterprise', ref(null))

const handleUpdate = (updated) => {
  enterprise.value = updated
}
</script>
