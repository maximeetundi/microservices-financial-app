<template>
  <NuxtLayout name="enterprise">
    <ServicesTab :enterprise="enterprise" :enterprise-id="enterpriseId?.value || ''" @update="handleUpdate" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const ServicesTab = defineAsyncComponent(() => import('@/components/enterprise/ServicesTab.vue'))

const enterprise = inject('enterprise', ref(null))
const enterpriseId = inject('enterpriseId', ref(''))

const handleUpdate = (updated) => {
  enterprise.value = updated
}
</script>
