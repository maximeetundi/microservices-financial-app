<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 flex items-center justify-center z-[9999]">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="close"></div>
      
      <!-- Modal -->
      <div class="relative bg-white dark:bg-gray-800 rounded-2xl shadow-2xl max-w-md w-full mx-4 overflow-hidden animate-scale-in">
        <!-- Icon -->
        <div :class="['p-6 flex items-center justify-center', iconBgColor]">
          <div class="text-6xl">{{ icon }}</div>
        </div>
        
        <!-- Content -->
        <div class="p-6 text-center">
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">{{ title }}</h3>
          <p class="text-gray-600 dark:text-gray-400">{{ message }}</p>
        </div>
        
        <!-- Actions -->
        <div class="p-6 border-t border-gray-200 dark:border-gray-700 flex gap-3">
          <button 
            v-if="showCancel"
            @click="handleCancel" 
            class="flex-1 px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 font-medium hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            {{ cancelText }}
          </button>
          <button 
            @click="handleConfirm" 
            :class="['flex-1 px-4 py-3 rounded-lg font-medium transition-colors', buttonClass]"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  visible: boolean
  type?: 'success' | 'error' | 'warning' | 'info' | 'confirm'
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  showCancel?: boolean
}>()

const emit = defineEmits(['confirm', 'cancel', 'close'])

const icon = computed(() => {
  switch (props.type) {
    case 'success': return '✅'
    case 'error': return '❌'
    case 'warning': return '⚠️'
    case 'confirm': return '❓'
    default: return 'ℹ️'
  }
})

const iconBgColor = computed(() => {
  switch (props.type) {
    case 'success': return 'bg-green-50 dark:bg-green-900/20'
    case 'error': return 'bg-red-50 dark:bg-red-900/20'
    case 'warning': return 'bg-orange-50 dark:bg-orange-900/20'
    default: return 'bg-blue-50 dark:bg-blue-900/20'
  }
})

const buttonClass = computed(() => {
  switch (props.type) {
    case 'success': return 'bg-green-600 hover:bg-green-700 text-white'
    case 'error': return 'bg-red-600 hover:bg-red-700 text-white'
    case 'warning': return 'bg-orange-600 hover:bg-orange-700 text-white'
    default: return 'bg-indigo-600 hover:bg-indigo-700 text-white'
  }
})

const handleConfirm = () => {
  emit('confirm')
  emit('close')
}

const handleCancel = () => {
  emit('cancel')
  emit('close')
}

const close = () => {
  if (!props.showCancel) {
    emit('close')
  }
}
</script>

<style scoped>
@keyframes scale-in {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-scale-in {
  animation: scale-in 0.2s ease-out;
}
</style>
