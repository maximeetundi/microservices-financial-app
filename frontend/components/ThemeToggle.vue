<template>
  <button 
    @click="toggleTheme" 
    class="relative p-2 rounded-full transition-colors duration-300 hover:bg-black/5 dark:hover:bg-white/10 group"
    aria-label="Toggle Theme"
  >
    <!-- Sun Icon -->
    <svg 
      v-if="$colorMode.value === 'light'"
      xmlns="http://www.w3.org/2000/svg" 
      class="w-6 h-6 text-orange-500 transform transition-transform duration-500 group-hover:rotate-90" 
      fill="none" 
      viewBox="0 0 24 24" 
      stroke="currentColor"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
    </svg>

    <!-- Moon Icon -->
    <svg 
      v-else
      xmlns="http://www.w3.org/2000/svg" 
      class="w-6 h-6 text-blue-400 transform transition-transform duration-500 group-hover:-rotate-12" 
      fill="none" 
      viewBox="0 0 24 24" 
      stroke="currentColor"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
    </svg>
  </button>
</template>

<script setup>
const colorMode = useColorMode()

const toggleTheme = () => {
  colorMode.preference = colorMode.value === 'light' ? 'dark' : 'light'
}
</script>
