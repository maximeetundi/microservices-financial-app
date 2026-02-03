<template>
  <div class="flex flex-col lg:flex-row min-h-[calc(100vh-8rem)] lg:h-[calc(100vh-4rem)] lg:overflow-hidden bg-surface rounded-xl shadow-xl border border-secondary-200 dark:border-secondary-800">
    <!-- Desktop Sidebar -->
    <div class="hidden lg:block h-full w-64 border-r border-secondary-200 dark:border-secondary-800">
      <ShopSidebar />
    </div>

    <!-- Mobile Header for Sub-nav -->
    <div class="lg:hidden w-full bg-surface border-b border-secondary-200 dark:border-secondary-800 p-4 flex items-center justify-between sticky top-0 z-20">
      <span class="font-bold text-lg text-primary-600 dark:text-primary-400">Menu Boutique</span>
      <button 
        @click="isMobileMenuOpen = !isMobileMenuOpen" 
        class="px-4 py-2 rounded-lg bg-secondary-100 dark:bg-secondary-800 text-sm font-medium hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-colors"
      >
        {{ isMobileMenuOpen ? 'Fermer' : 'Menu' }}
      </button>
    </div>

    <!-- Mobile Sidebar Drawer -->
    <div v-if="isMobileMenuOpen" class="lg:hidden fixed inset-0 z-50 flex" style="top: 0;">
      <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="isMobileMenuOpen = false"></div>
      <div class="relative flex-1 max-w-[280px] w-full bg-surface shadow-2xl h-full overflow-y-auto">
        <div class="p-4 border-b border-secondary-200 dark:border-secondary-800 flex justify-between items-center">
            <span class="font-bold text-lg">Menu</span>
            <button @click="isMobileMenuOpen = false" class="p-2 rounded-full hover:bg-secondary-100 dark:hover:bg-secondary-800">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
        </div>
        <ShopSidebar />
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 lg:overflow-y-auto w-full h-full scrollbar-thin bg-base">
      <div class="p-4 lg:p-8 w-full max-w-full">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ShopSidebar from './ShopSidebar.vue'

const isMobileMenuOpen = ref(false)
</script>

<style>
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}
.dark .scrollbar-thin::-webkit-scrollbar-thumb {
  background: #475569;
}
</style>
