<template>
  <div class="relative">
    <!-- Attachment Button -->
    <button @click="showMenu = !showMenu" :class="['p-2 rounded-full transition-colors', showMenu ? 'bg-gray-200 dark:bg-gray-700' : 'hover:bg-gray-100 dark:hover:bg-gray-800']">
      <svg class="w-6 h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
      </svg>
    </button>

    <!-- Attachment Menu -->
    <Teleport to="body">
      <div v-if="showMenu" class="fixed inset-0 z-40" @click="showMenu = false"></div>
    </Teleport>
    
    <div v-if="showMenu" class="absolute bottom-full mb-2 left-0 bg-white dark:bg-gray-800 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-700 py-2 z-50 min-w-[200px]">
      <!-- Image -->
      <button @click="selectImage" class="w-full px-4 py-3 flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
        <div class="w-10 h-10 rounded-full bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center">
          <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <div class="flex-1 text-left">
          <div class="font-medium text-sm">Photo</div>
          <div class="text-xs text-gray-500">Depuis la galerie</div>
        </div>
      </button>

      <!-- Document -->
      <button @click="selectDocument" class="w-full px-4 py-3 flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
        <div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
          <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <div class="flex-1 text-left">
          <div class="font-medium text-sm">Document</div>
          <div class="text-xs text-gray-500">PDF, DOC, etc.</div>
        </div>
      </button>

      <!-- Camera (future) -->
      <button disabled class="w-full px-4 py-3 flex items-center gap-3 opacity-50 cursor-not-allowed">
        <div class="w-10 h-10 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center">
          <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div class="flex-1 text-left">
          <div class="font-medium text-sm">Caméra</div>
          <div class="text-xs text-gray-500">Bientôt disponible</div>
        </div>
      </button>
    </div>

    <!-- Hidden File Inputs -->
    <input ref="imageInput" type="file" accept="image/*" @change="handleImageSelect" class="hidden" />
    <input ref="documentInput" type="file" accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.zip,.rar" @change="handleDocumentSelect" class="hidden" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['fileSelected'])

const showMenu = ref(false)
const imageInput = ref<HTMLInputElement>()
const documentInput = ref<HTMLInputElement>()

const selectImage = () => {
  showMenu.value = false
  imageInput.value?.click()
}

const selectDocument = () => {
  showMenu.value = false
  documentInput.value?.click()
}

const handleImageSelect = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) {
    emit('fileSelected', { file, type: 'image' })
  }
}

const handleDocumentSelect = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) {
    emit('fileSelected', { file, type: 'document' })
  }
}
</script>
