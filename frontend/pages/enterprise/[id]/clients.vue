<template>
  <NuxtLayout name="enterprise">
    <ClientsTab :enterprise="enterprise" :enterprise-id="enterpriseId?.value || ''" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
const ClientsTab = defineAsyncComponent(() => import('@/components/enterprise/ClientsTab.vue'))

const enterprise = inject('enterprise', ref(null))
const enterpriseId = inject('enterpriseId', ref(''))
</script>
