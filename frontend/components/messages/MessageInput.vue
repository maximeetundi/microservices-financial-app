<template>
  <div>
    <!-- Image Preview -->
    <div v-if="selectedImage" class="px-4 py-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="relative inline-block">
        <img :src="selectedImage.preview" class="max-h-32 rounded-xl shadow-md" />
        <button @click="removeImage" class="absolute -top-2 -right-2 w-7 h-7 rounded-full bg-red-500 text-white flex items-center justify-center hover:bg-red-600 transition-all hover:scale-110 shadow-lg">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <input v-model="imageCaption" type="text" placeholder="Ajouter une légende..." 
        class="mt-2 w-full px-4 py-2.5 rounded-xl bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 text-sm focus:ring-2 focus:ring-green-500 focus:border-transparent" />
    </div>

    <!-- Document Preview -->
    <div v-if="selectedDocument" class="px-4 py-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="flex items-center gap-3 p-3 bg-white dark:bg-gray-700 rounded-xl shadow-sm">
        <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
          <span class="text-2xl">📄</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 dark:text-white truncate">{{ selectedDocument.file.name }}</p>
          <p class="text-xs text-gray-500">{{ formatFileSize(selectedDocument.file.size) }}</p>
        </div>
        <button @click="removeDocument" class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors">
          <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Video Preview -->
    <div v-if="selectedVideo" class="px-4 py-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="relative inline-block">
        <video :src="selectedVideo.preview" class="max-h-40 rounded-xl shadow-md" controls></video>
        <button @click="removeVideo" class="absolute -top-2 -right-2 w-7 h-7 rounded-full bg-red-500 text-white flex items-center justify-center hover:bg-red-600 transition-all hover:scale-110 shadow-lg">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <p class="text-xs text-gray-500 mt-2">{{ formatFileSize(selectedVideo.file.size) }}</p>
    </div>

    <!-- Audio Recorder -->
    <AudioRecorder v-if="isRecording" :is-recording="isRecording" @audioRecorded="handleAudioRecorded" @cancel="isRecording = false" />

    <!-- Input Area -->
    <div v-if="!isRecording" class="bg-gray-50 dark:bg-gray-800 px-2 md:px-3 py-2 md:py-3 border-t border-gray-200 dark:border-gray-700">
      <div class="flex items-center gap-1.5 md:gap-2">
        <!-- Emoji Picker -->
        <EmojiPicker @select="insertEmoji" />

        <!-- Attachment Menu -->
        <AttachmentMenu @fileSelected="handleFileSelected" />

        <!-- Text Input -->
        <div class="flex-1 relative min-w-0">
          <input 
            ref="inputRef"
            v-model="message" 
            @keyup.enter="sendMessage" 
            @input="handleTyping"
            @focus="handleFocus"
            @blur="handleBlur"
            type="text" 
            :placeholder="isUploading ? 'Envoi en cours...' : 'Tapez un message...'"
            :disabled="isUploading"
            class="w-full px-3 md:px-4 py-2 md:py-2.5 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm transition-all disabled:opacity-50 pr-10" 
          />
          <!-- Character count when typing long message -->
          <span v-if="message.length > 100" class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-gray-400">
            {{ message.length }}/1000
          </span>
        </div>

        <!-- Send or Microphone Button -->
        <button 
          v-if="message.trim() || selectedImage || selectedDocument || selectedVideo" 
          @click="sendMessage" 
          :disabled="isUploading"
          class="w-10 h-10 md:w-11 md:h-11 flex-shrink-0 rounded-full bg-gradient-to-br from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white flex items-center justify-center transition-all disabled:opacity-50 shadow-lg hover:shadow-xl hover:scale-105"
        >
          <svg v-if="!isUploading" class="w-4 h-4 md:w-5 md:h-5 transform rotate-45" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
          </svg>
          <svg v-else class="animate-spin w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </button>
        <button v-else @click="startAudioRecording" 
          class="w-10 h-10 md:w-11 md:h-11 flex-shrink-0 rounded-full bg-gradient-to-br from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white flex items-center justify-center transition-all shadow-lg hover:shadow-xl hover:scale-105">
          <svg class="w-4 h-4 md:w-5 md:h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import AttachmentMenu from './AttachmentMenu.vue'
import AudioRecorder from './AudioRecorder.vue'
import EmojiPicker from './EmojiPicker.vue'
import api from '~/composables/useApi'

const props = defineProps<{
  conversationId?: string
}>()

const emit = defineEmits(['messageSent', 'typing'])

const inputRef = ref<HTMLInputElement>()
const message = ref('')
const selectedImage = ref<any>(null)
const selectedVideo = ref<any>(null)
const selectedDocument = ref<any>(null)
const imageCaption = ref('')
const isUploading = ref(false)
const isRecording = ref(false)

// Typing indicator management
let typingTimeout: ReturnType<typeof setTimeout> | null = null

const handleTyping = () => {
  emit('typing', true)
  
  // Clear existing timeout
  if (typingTimeout) {
    clearTimeout(typingTimeout)
  }
  
  // Stop typing after 2 seconds of inactivity
  typingTimeout = setTimeout(() => {
    emit('typing', false)
  }, 2000)
}

const handleFocus = () => {
  // Could emit focus event if needed
}

const handleBlur = () => {
  if (typingTimeout) {
    clearTimeout(typingTimeout)
  }
  emit('typing', false)
}

const insertEmoji = (emoji: string) => {
  message.value += emoji
  inputRef.value?.focus()
}

const handleFileSelected = async ({ file, type }: { file: File, type: string }) => {
  if (type === 'image') {
    selectedImage.value = {
      file,
      preview: URL.createObjectURL(file)
    }
  } else if (type === 'video') {
    selectedVideo.value = {
      file,
      preview: URL.createObjectURL(file)
    }
  } else {
    selectedDocument.value = { file }
  }
}

const removeImage = () => {
  if (selectedImage.value?.preview) {
    URL.revokeObjectURL(selectedImage.value.preview)
  }
  selectedImage.value = null
  imageCaption.value = ''
}

const removeVideo = () => {
  if (selectedVideo.value?.preview) {
    URL.revokeObjectURL(selectedVideo.value.preview)
  }
  selectedVideo.value = null
}

const removeDocument = () => {
  selectedDocument.value = null
}

const startAudioRecording = () => {
  isRecording.value = true
}

const handleAudioRecorded = async ({ blob, duration }: { blob: Blob, duration: number }) => {
  isRecording.value = false
  
  // Upload audio and send message
  const audioFile = new File([blob], `audio-${Date.now()}.webm`, { type: 'audio/webm' })
  await uploadAndSend(audioFile, 'audio', '', duration)
}

const uploadAndSend = async (file: File, messageType: string, caption: string = '', duration?: number) => {
  isUploading.value = true
  
  try {
    // Upload file to MinIO
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', messageType)

    const uploadRes = await api.post('/messaging-service/api/v1/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    const { url, filename, size, mime_type } = uploadRes.data

    // Send message with attachment
    const messageData = {
      content: caption || message.value,
      message_type: messageType,
      attachment: {
        file_url: url,
        file_name: filename,
        file_size: size,
        mime_type,
        ...(duration && { duration })
      }
    }

    if (props.conversationId) {
      await api.post(`/messaging-service/api/v1/conversations/${props.conversationId}/messages`, messageData)
    }

    // Clear inputs
    message.value = ''
    removeImage()
    removeVideo()
    removeDocument()
    
    emit('messageSent')
  } catch (err: any) {
    console.error('Upload/send error:', err)
    alert(err.response?.data?.error || 'Erreur lors de l\'envoi')
  } finally {
    isUploading.value = false
  }
}

const sendMessage = async () => {
  if (!message.value.trim() && !selectedImage.value && !selectedVideo.value && !selectedDocument.value) return

  // Stop typing indicator
  emit('typing', false)
  if (typingTimeout) {
    clearTimeout(typingTimeout)
  }

  if (selectedImage.value) {
    await uploadAndSend(selectedImage.value.file, 'image', imageCaption.value)
  } else if (selectedVideo.value) {
    await uploadAndSend(selectedVideo.value.file, 'video')
  } else if (selectedDocument.value) {
    await uploadAndSend(selectedDocument.value.file, 'document')
  } else {
    // Text only
    try {
      const messageData = {
        content: message.value,
        message_type: 'text'
      }

      if (props.conversationId) {
        await api.post(`/messaging-service/api/v1/conversations/${props.conversationId}/messages`, messageData)
      }

      message.value = ''
      emit('messageSent')
    } catch (err: any) {
      alert(err.response?.data?.error || 'Erreur')
    }
  }
}

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>
