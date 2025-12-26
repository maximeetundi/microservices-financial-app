<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="modalState.isOpen" 
        class="modal-overlay"
        @click.self="handleBackdropClick"
      >
        <!-- Backdrop -->
        <div class="modal-backdrop"></div>
        
        <!-- Modal Content -->
        <div class="modal-container">
          <!-- Close Button (for info/success modals) -->
          <button 
            v-if="!modalState.showCancel && !modalState.persistent"
            @click="handleConfirm"
            class="modal-close-btn"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
          
          <!-- Icon with animated ring -->
          <div class="modal-icon-wrapper">
            <div :class="iconRingClasses">
              <div :class="iconBgClasses">
                <span class="modal-icon-emoji">{{ iconEmoji }}</span>
              </div>
            </div>
          </div>
          
          <!-- Title -->
          <h3 class="modal-title">{{ modalState.title }}</h3>
          
          <!-- Message -->
          <p class="modal-message">{{ modalState.message }}</p>
          
          <!-- Buttons -->
          <div class="modal-buttons">
            <button 
              v-if="modalState.showCancel"
              @click="handleCancel"
              class="modal-btn modal-btn-secondary"
            >
              {{ modalState.cancelText }}
            </button>
            <button 
              @click="handleConfirm"
              :class="confirmButtonClasses"
              class="modal-btn"
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

const iconRingClasses = computed(() => {
  const base = 'modal-icon-ring'
  switch (modalState.type) {
    case 'success': return `${base} ring-success`
    case 'error': return `${base} ring-error`
    case 'warning': return `${base} ring-warning`
    case 'confirm': return `${base} ring-info`
    default: return `${base} ring-info`
  }
})

const iconBgClasses = computed(() => {
  const base = 'modal-icon-bg'
  switch (modalState.type) {
    case 'success': return `${base} bg-success`
    case 'error': return `${base} bg-error`
    case 'warning': return `${base} bg-warning`
    case 'confirm': return `${base} bg-info`
    default: return `${base} bg-info`
  }
})

const iconEmoji = computed(() => {
  switch (modalState.type) {
    case 'success': return '✓'
    case 'error': return '✕'
    case 'warning': return '!'
    case 'confirm': return '?'
    default: return 'i'
  }
})

const confirmButtonClasses = computed(() => {
  const base = 'modal-btn-primary'
  switch (modalState.type) {
    case 'success': return `${base} btn-success`
    case 'error': return `${base} btn-error`
    case 'warning': return `${base} btn-warning`
    default: return `${base} btn-primary`
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
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
}

.modal-container {
  position: relative;
  background: linear-gradient(180deg, var(--color-surface) 0%, var(--color-surface-hover) 100%);
  border-radius: 24px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  max-width: 380px;
  width: 100%;
  padding: 2rem 1.5rem;
  text-align: center;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-close-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.5rem;
  border-radius: 50%;
  color: var(--color-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.modal-close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-base);
}

.modal-icon-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.modal-icon-ring {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse-ring 2s infinite;
}

.modal-icon-ring.ring-success { background: rgba(16, 185, 129, 0.15); }
.modal-icon-ring.ring-error { background: rgba(239, 68, 68, 0.15); }
.modal-icon-ring.ring-warning { background: rgba(245, 158, 11, 0.15); }
.modal-icon-ring.ring-info { background: rgba(59, 130, 246, 0.15); }

.modal-icon-bg {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: bold;
  color: white;
}

.modal-icon-bg.bg-success { background: linear-gradient(135deg, #10B981, #059669); }
.modal-icon-bg.bg-error { background: linear-gradient(135deg, #EF4444, #DC2626); }
.modal-icon-bg.bg-warning { background: linear-gradient(135deg, #F59E0B, #D97706); }
.modal-icon-bg.bg-info { background: linear-gradient(135deg, #3B82F6, #2563EB); }

.modal-icon-emoji {
  font-family: system-ui, -apple-system, sans-serif;
  text-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.modal-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--color-base);
  margin-bottom: 0.75rem;
  line-height: 1.3;
}

.modal-message {
  font-size: 0.9375rem;
  color: #FFFFFF;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.modal-buttons {
  display: flex;
  gap: 0.75rem;
}

.modal-btn {
  flex: 1;
  padding: 0.875rem 1.25rem;
  border-radius: 14px;
  font-weight: 600;
  font-size: 0.9375rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.modal-btn-secondary {
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-muted);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-btn-secondary:hover {
  background: rgba(255, 255, 255, 0.12);
  color: var(--color-base);
}

.modal-btn-primary {
  color: white;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.25);
}

.modal-btn-primary.btn-success {
  background: linear-gradient(135deg, #10B981, #059669);
}
.modal-btn-primary.btn-success:hover {
  background: linear-gradient(135deg, #059669, #047857);
}

.modal-btn-primary.btn-error {
  background: linear-gradient(135deg, #EF4444, #DC2626);
}
.modal-btn-primary.btn-error:hover {
  background: linear-gradient(135deg, #DC2626, #B91C1C);
}

.modal-btn-primary.btn-warning {
  background: linear-gradient(135deg, #F59E0B, #D97706);
}
.modal-btn-primary.btn-warning:hover {
  background: linear-gradient(135deg, #D97706, #B45309);
}

.modal-btn-primary.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
}
.modal-btn-primary.btn-primary:hover {
  background: linear-gradient(135deg, #5a67d8, #6b46c1);
}

@keyframes pulse-ring {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.05); opacity: 0.8; }
}

/* Transition animations */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9) translateY(20px);
  opacity: 0;
}

/* Mobile optimizations */
@media (max-width: 480px) {
  .modal-overlay {
    padding: 0.75rem;
    align-items: flex-end;
  }
  
  .modal-container {
    border-radius: 24px 24px 16px 16px;
    padding: 1.5rem 1.25rem;
    max-width: 100%;
  }
  
  .modal-icon-ring {
    width: 70px;
    height: 70px;
  }
  
  .modal-icon-bg {
    width: 52px;
    height: 52px;
    font-size: 1.5rem;
  }
  
  .modal-title {
    font-size: 1.25rem;
  }
  
  .modal-message {
    font-size: 0.875rem;
  }
  
  .modal-btn {
    padding: 0.875rem 1rem;
  }
}
</style>
