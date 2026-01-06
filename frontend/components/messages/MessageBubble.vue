<template>
  <div :class="['flex', message.isMine ? 'justify-end' : 'justify-start', 'animate-fade-in']">
    <div :class="['max-w-[85%] md:max-w-[70%] rounded-lg shadow-md overflow-hidden break-words', message.isMine ? 'bg-green-500 text-white' : 'bg-white dark:bg-gray-800 text-gray-900 dark:text-white']">
      <!-- Sender name for group chats -->
      <div v-if="!message.isMine && showSenderName" :class="['px-4 pt-2 text-xs font-medium', message.isMine ? 'text-white/70' : 'text-green-600']">
        {{ message.senderName }}
      </div>

      <!-- Image Message -->
      <div v-if="message.message_type === 'image' && message.attachment" class="cursor-pointer" @click="$emit('openImage', message.attachment.file_url)">
        <img :src="message.attachment.file_url" :alt="message.attachment.file_name" class="max-w-full h-auto" loading="lazy" />
        <div v-if="message.content" class="px-4 py-2">
          <p class="text-sm break-words">{{ message.content }}</p>
        </div>
      </div>

      <!-- Audio Message -->
      <div v-else-if="message.message_type === 'audio' && message.attachment" class="px-4 py-3">
        <div class="flex items-center gap-3">
          <button @click="toggleAudioPlay" :class="['w-10 h-10 rounded-full flex items-center justify-center transition-colors', message.isMine ? 'bg-white/20 hover:bg-white/30' : 'bg-green-500 hover:bg-green-600 text-white']">
            <svg v-if="!isPlaying" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M5 4a2 2 0 012-2h6a2 2 0 012 2v12a2 2 0 01-2 2H7a2 2 0 01-2-2V4z" />
            </svg>
          </button>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <div class="flex-1 h-1 bg-white/20 rounded-full overflow-hidden">
                <div class="h-full bg-white transition-all" :style="{ width: audioProgress + '%' }"></div>
              </div>
              <span class="text-xs opacity-80 min-w-[40px]">{{ formatDuration(message.attachment.duration || 0) }}</span>
            </div>
          </div>
          <audio ref="audioPlayer" :src="message.attachment.file_url" @timeupdate="updateAudioProgress" @ended="isPlaying = false" />
        </div>
      </div>

      <!-- Document Message -->
      <div v-else-if="message.message_type === 'document' && message.attachment" class="px-4 py-3">
        <div class="flex items-center gap-3">
          <div :class="['w-12 h-12 rounded-lg flex items-center justify-center', message.isMine ? 'bg-white/20' : 'bg-gray-100']">
            <span class="text-2xl">{{ getFileIcon(message.attachment.mime_type) }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{{ message.attachment.file_name }}</p>
            <p :class="['text-xs', message.isMine ? 'text-white/70' : 'text-gray-500']">{{ formatFileSize(message.attachment.file_size) }}</p>
          </div>
          <a :href="message.attachment.file_url" download :class="['w-8 h-8 rounded-full flex items-center justify-center transition-colors', message.isMine ? 'bg-white/20 hover:bg-white/30' : 'bg-gray-100 hover:bg-gray-200']">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path d="M13 8V2H7v6H2l8 8 8-8h-5zM0 18h20v2H0v-2z" />
            </svg>
          </a>
        </div>
      </div>

      <!-- Video Message -->
      <div v-else-if="message.message_type === 'video' && message.attachment" class="p-1">
        <VideoPlayer 
          :src="message.attachment.file_url" 
          :mime-type="message.attachment.mime_type"
          :poster="message.attachment.thumbnail_url"
        />
        <div v-if="message.content" class="px-3 py-2">
          <p class="text-sm break-words">{{ message.content }}</p>
        </div>
      </div>

      <!-- Text Message -->
      <div v-else class="px-4 py-2">
        <p class="text-sm break-words whitespace-pre-wrap">{{ message.content }}</p>
      </div>

      <!-- Timestamp -->
      <div :class="['px-4 pb-2 text-xs flex items-center gap-1 justify-end', message.isMine ? 'text-white/70' : 'text-gray-500']">
        <span>{{ formatTime(message.created_at || message.createdAt) }}</span>
        <svg v-if="message.isMine" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" />
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import VideoPlayer from './VideoPlayer.vue'

const props = defineProps<{
  message: any
  showSenderName?: boolean
}>()

defineEmits(['openImage'])

const audioPlayer = ref<HTMLAudioElement>()
const isPlaying = ref(false)
const audioProgress = ref(0)

const toggleAudioPlay = () => {
  if (!audioPlayer.value) return
  if (isPlaying.value) {
    audioPlayer.value.pause()
  } else {
    audioPlayer.value.play()
  }
  isPlaying.value = !isPlaying.value
}

const updateAudioProgress = () => {
  if (!audioPlayer.value) return
  const progress = (audioPlayer.value.currentTime / audioPlayer.value.duration) * 100
  audioProgress.value = progress || 0
}

const formatTime = (date: any) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

const formatDuration = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const getFileIcon = (mimeType: string) => {
  if (mimeType?.includes('pdf')) return '📄'
  if (mimeType?.includes('word')) return '📝'
  if (mimeType?.includes('excel') || mimeType?.includes('spreadsheet')) return '📊'
  if (mimeType?.includes('powerpoint') || mimeType?.includes('presentation')) return '📑'
  if (mimeType?.includes('zip') || mimeType?.includes('rar')) return '🗜️'
  return '📎'
}
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}
</style>
