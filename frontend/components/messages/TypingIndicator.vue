<template>
  <div class="flex justify-start animate-fade-in mb-1">
    <div class="flex items-center gap-2 px-4 py-3 bg-white dark:bg-gray-800 rounded-2xl rounded-bl-md shadow-sm">
      <div class="flex gap-1">
        <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0ms"></span>
        <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 150ms"></span>
        <span class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 300ms"></span>
      </div>
      <span class="text-xs text-gray-500 dark:text-gray-400 ml-1">{{ name }} écrit...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  name?: string
}>()
</script>

<style scoped>
@keyframes bounce {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-4px);
  }
}

.animate-bounce {
  animation: bounce 1s infinite;
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.2s ease-out;
}
</style>
