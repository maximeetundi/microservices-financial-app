<template>
  <NuxtLayout name="enterprise">
    <WalletsTab :enterprise="enterprise" @update="handleUpdate" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const WalletsTab = defineAsyncComponent(() => import('@/components/enterprise/WalletsTab.vue'))

const enterprise = inject('enterprise', ref(null))

const handleUpdate = (updated) => {
  enterprise.value = updated
}
</script>
