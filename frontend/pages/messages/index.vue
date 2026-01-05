<template>
  <NuxtLayout name="dashboard">
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Messages</h1>
          <p class="text-gray-600 dark:text-gray-400">Conversations privées</p>
        </div>
      </div>

      <!-- Conversations List -->
      <div class="bg-surface rounded-2xl border border-secondary-200 dark:border-secondary-700">
        <div class="p-6">
          <div class="text-center py-12 text-gray-500">
            <span class="text-6xl mb-4 block">💬</span>
            <p class="text-lg">Messagerie utilisateur à utilisateur</p>
            <p class="text-sm mt-2">Fonctionnalité à venir - Backend prêt</p>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>
