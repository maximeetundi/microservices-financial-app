<template>
  <NuxtLayout name="enterprise">
    <EmployeesTab :enterprise="enterprise" :enterprise-id="enterpriseId?.value || ''" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
import EmployeesTab from '@/components/enterprise/EmployeesTab.vue'

const enterprise = inject('enterprise', ref(null))
const enterpriseId = inject('enterpriseId', ref(''))
</script>
