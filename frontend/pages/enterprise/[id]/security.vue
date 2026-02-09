<template>
  <NuxtLayout name="enterprise">
    <SecurityTab :enterprise="enterprise" :enterprise-id="enterpriseId?.value || ''" @update="handleUpdate" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const SecurityTab = defineAsyncComponent(() => import('@/components/enterprise/SecurityTab.vue'))

const enterprise = inject('enterprise', ref(null))
const enterpriseId = inject('enterpriseId', ref(''))

const handleUpdate = (updated) => {
  enterprise.value = updated
}
</script>
