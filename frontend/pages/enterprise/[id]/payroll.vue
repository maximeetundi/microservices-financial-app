<template>
  <NuxtLayout name="enterprise">
    <PayrollTab :enterprise-id="enterpriseId" />
  </NuxtLayout>
</template>

<script setup>
import { inject, computed } from 'vue'
import { useRoute } from 'vue-router'
const PayrollTab = defineAsyncComponent(() => import('@/components/enterprise/PayrollTab.vue'))

const route = useRoute()
const enterpriseId = computed(() => route.params.id)
</script>
