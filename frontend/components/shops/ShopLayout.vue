<template>
  <div class="flex h-[calc(100vh-8rem)] lg:h-[calc(100vh-4rem)] overflow-hidden bg-white dark:bg-slate-900 rounded-none md:rounded-2xl shadow-xl ring-1 ring-black/5 dark:ring-white/10">
    <!-- Desktop Sidebar -->
    <div class="hidden lg:block h-full">
      <ShopSidebar />
    </div>

    <!-- Mobile Header for Sub-nav -->
    <div class="lg:hidden w-full bg-white dark:bg-slate-900 border-b border-gray-200 dark:border-gray-700 p-4 flex items-center justify-between z-10">
      <span class="font-bold text-gray-900 dark:text-white">Menu Boutique</span>
      <button @click="isMobileMenuOpen = !isMobileMenuOpen" class="text-indigo-600 font-medium text-sm">
        {{ isMobileMenuOpen ? 'Fermer' : 'Menu' }}
      </button>
    </div>

    <!-- Mobile Sidebar Drawer -->
    <div v-if="isMobileMenuOpen" class="lg:hidden absolute inset-0 z-40 flex">
      <div class="absolute inset-0 bg-black/50" @click="isMobileMenuOpen = false"></div>
      <div class="relative flex-1 max-w-xs w-full bg-white dark:bg-slate-900 shadow-2xl">
        <ShopSidebar />
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-y-auto h-full scrollbar-thin bg-gray-50/50 dark:bg-slate-900/50">
      <div class="p-6">
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
