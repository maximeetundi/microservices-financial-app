<template>
  <NuxtLayout name="enterprise">
    <BillingTab :enterprise="enterprise" :enterprise-id="enterpriseId?.value || ''" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const BillingTab = defineAsyncComponent(() => import('@/components/enterprise/BillingTab.vue'))

const enterprise = inject('enterprise', ref(null))
const enterpriseId = inject('enterpriseId', ref(''))
</script>
