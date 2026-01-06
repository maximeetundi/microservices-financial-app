<template>
  <div>
    <!-- Image Preview -->
    <div v-if="selectedImage" class="px-4 py-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="relative inline-block">
        <img :src="selectedImage.preview" class="max-h-32 rounded-lg" />
        <button @click="removeImage" class="absolute -top-2 -right-2 w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center hover:bg-red-600 transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <input v-model="imageCaption" type="text" placeholder="Ajouter une légende..." class="mt-2 w-full px-3 py-2 rounded-lg bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 text-sm" />
    </div>

    <!-- Document Preview -->
    <div v-if="selectedDocument" class="px-4 py-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="flex items-center gap-3 p-3 bg-white dark:bg-gray-700 rounded-lg">
        <div class="w-10 h-10 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
          <span class="text-xl">📄</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium truncate">{{ selectedDocument.file.name }}</p>
          <p class="text-xs text-gray-500">{{ formatFileSize(selectedDocument.file.size) }}</p>
        </div>
        <button @click="removeDocument" class="p-1 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors">
          <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Audio Recorder -->
    <AudioRecorder v-if="isRecording" :is-recording="isRecording" @audioRecorded="handleAudioRecorded" @cancel="isRecording = false" />

    <!-- Input Area -->
    <div v-if="!isRecording" class="bg-gray-50 dark:bg-gray-800 px-4 py-3 border-t border-gray-200 dark:border-gray-700">
      <div class="flex items-center gap-3">
        <!-- Emoji Picker -->
        <EmojiPicker @select="insertEmoji" />

        <!-- Attachment Menu -->
        <AttachmentMenu @fileSelected="handleFileSelected" />

        <!-- Text Input -->
        <input 
          v-model="message" 
          @keyup.enter="sendMessage" 
          type="text" 
          :placeholder="isUploading ? 'Envoi en cours...' : 'Tapez un message...'"
          :disabled="isUploading"
          class="flex-1 px-4 py-2 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 focus:ring-2 focus:ring-green-500 disabled:opacity-50" 
        />

        <!-- Send or Microphone Button -->
        <button 
          v-if="message.trim() || selectedImage || selectedDocument" 
          @click="sendMessage" 
          :disabled="isUploading"
          class="w-10 h-10 rounded-full bg-green-500 hover:bg-green-600 text-white flex items-center justify-center transition-colors disabled:opacity-50"
        >
          <svg v-if="!isUploading" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
          </svg>
          <svg v-else class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </button>
        <button v-else @click="startAudioRecording" class="w-10 h-10 rounded-full bg-green-500 hover:bg-green-600 text-white flex items-center justify-center transition-colors">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AttachmentMenu from './AttachmentMenu.vue'
import AudioRecorder from './AudioRecorder.vue'
import EmojiPicker from './EmojiPicker.vue'
import api from '~/composables/useApi'

const props = defineProps<{
  conversationId?: string
  associationId?: string
}>()

const emit = defineEmits(['messageSent'])

const message = ref('')
const selectedImage = ref<any>(null)
const selectedDocument = ref<any>(null)
const imageCaption = ref('')
const isUploading = ref(false)
const isRecording = ref(false)

const insertEmoji = (emoji: string) => {
  message.value += emoji
}

const handleFileSelected = async ({ file, type }: { file: File, type: string }) => {
  if (type === 'image') {
    selectedImage.value = {
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

    if (props.associationId) {
      await api.post(`/association-service/api/v1/associations/${props.associationId}/chat`, messageData)
    } else if (props.conversationId) {
      await api.post(`/messaging-service/api/v1/conversations/${props.conversationId}/messages`, messageData)
    }

    // Clear inputs
    message.value = ''
    removeImage()
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
  if (!message.value.trim() && !selectedImage.value && !selectedDocument.value) return

  if (selectedImage.value) {
    await uploadAndSend(selectedImage.value.file, 'image', imageCaption.value)
  } else if (selectedDocument.value) {
    await uploadAndSend(selectedDocument.value.file, 'document')
  } else {
    // Text only
    try {
      const messageData = {
        content: message.value,
        message_type: 'text'
      }

      if (props.associationId) {
        await api.post(`/association-service/api/v1/associations/${props.associationId}/chat`, messageData)
      } else if (props.conversationId) {
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
