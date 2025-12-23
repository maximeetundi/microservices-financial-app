<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="modalState.isOpen" 
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        @click.self="handleBackdropClick"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        
        <!-- Modal Content -->
        <div class="relative bg-surface rounded-2xl shadow-2xl max-w-md w-full p-6 transform transition-all">
          <!-- Icon -->
          <div class="flex justify-center mb-4">
            <div :class="iconClasses">
              <span class="text-3xl">{{ iconEmoji }}</span>
            </div>
          </div>
          
          <!-- Title -->
          <h3 class="text-xl font-bold text-center text-base mb-2">
            {{ modalState.title }}
          </h3>
          
          <!-- Message -->
          <p class="text-center text-muted mb-6">
            {{ modalState.message }}
          </p>
          
          <!-- Buttons -->
          <div class="flex gap-3" :class="modalState.showCancel ? 'justify-center' : 'justify-center'">
            <button 
              v-if="modalState.showCancel"
              @click="handleCancel"
              class="flex-1 px-4 py-3 rounded-xl font-medium text-muted bg-secondary-100 dark:bg-secondary-800 hover:bg-secondary-200 dark:hover:bg-secondary-700 transition-colors"
            >
              {{ modalState.cancelText }}
            </button>
            <button 
              @click="handleConfirm"
              :class="confirmButtonClasses"
              class="flex-1 px-4 py-3 rounded-xl font-medium transition-colors"
            >
              {{ modalState.confirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useModal } from '~/composables/useModal'

const { state: modalState, close } = useModal()

const iconClasses = computed(() => {
  const baseClasses = 'w-16 h-16 rounded-full flex items-center justify-center'
  switch (modalState.type) {
    case 'success':
      return `${baseClasses} bg-green-100 dark:bg-green-900/30`
    case 'error':
      return `${baseClasses} bg-red-100 dark:bg-red-900/30`
    case 'warning':
      return `${baseClasses} bg-yellow-100 dark:bg-yellow-900/30`
    case 'confirm':
      return `${baseClasses} bg-blue-100 dark:bg-blue-900/30`
    default:
      return `${baseClasses} bg-primary-100 dark:bg-primary-900/30`
  }
})

const iconEmoji = computed(() => {
  switch (modalState.type) {
    case 'success':
      return '✅'
    case 'error':
      return '❌'
    case 'warning':
      return '⚠️'
    case 'confirm':
      return '❓'
    default:
      return 'ℹ️'
  }
})

const confirmButtonClasses = computed(() => {
  switch (modalState.type) {
    case 'success':
      return 'bg-green-600 hover:bg-green-700 text-white'
    case 'error':
      return 'bg-red-600 hover:bg-red-700 text-white'
    case 'warning':
      return 'bg-yellow-600 hover:bg-yellow-700 text-white'
    case 'confirm':
      return 'bg-primary-600 hover:bg-primary-700 text-white'
    default:
      return 'bg-primary-600 hover:bg-primary-700 text-white'
  }
})

const handleBackdropClick = () => {
  if (!modalState.persistent) {
    close(false)
  }
}

const handleConfirm = () => {
  close(true)
}

const handleCancel = () => {
  close(false)
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.9) translateY(20px);
}
</style>
