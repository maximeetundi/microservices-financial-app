<template>
  <NuxtLayout name="enterprise">
    <EmployeesTab :enterprise="enterprise" />
  </NuxtLayout>
</template>

<script setup>
import { inject, ref } from 'vue'
import EmployeesTab from '@/components/enterprise/EmployeesTab.vue'

const enterprise = inject('enterprise', ref(null))
</script>
