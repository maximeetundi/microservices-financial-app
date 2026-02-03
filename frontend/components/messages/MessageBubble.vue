<template>
  <div :class="['flex', message.isMine ? 'justify-end' : 'justify-start', 'animate-fade-in mb-1']">
    <div :class="[
      'max-w-[85%] md:max-w-[70%] rounded-2xl shadow-sm overflow-hidden break-words relative',
      message.isMine 
        ? 'bg-gradient-to-br from-green-500 to-green-600 text-white rounded-br-md' 
        : 'bg-white dark:bg-gray-800 text-gray-900 dark:text-white rounded-bl-md'
    ]">
      <!-- Reply Reference (if replying to a message) -->
      <div v-if="message.reply_to" :class="['px-3 pt-2 pb-1 border-l-4', message.isMine ? 'border-white/40 bg-white/10' : 'border-green-500 bg-gray-50 dark:bg-gray-700']">
        <p class="text-xs opacity-70 truncate">{{ message.reply_to.sender_name }}</p>
        <p class="text-xs truncate">{{ message.reply_to.content }}</p>
      </div>

      <!-- Sender name for group chats -->
      <div v-if="!message.isMine && showSenderName" :class="['px-4 pt-2 text-xs font-semibold', message.isMine ? 'text-white/80' : 'text-green-600 dark:text-green-400']">
        {{ message.senderName }}
      </div>

      <!-- Image Message -->
      <div v-if="message.message_type === 'image' && message.attachment" class="cursor-pointer group" @click="$emit('openImage', message.attachment.file_url)">
        <div class="relative">
          <img :src="message.attachment.file_url" :alt="message.attachment.file_name" class="max-w-full h-auto" loading="lazy" />
          <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors flex items-center justify-center">
            <svg class="w-10 h-10 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
            </svg>
          </div>
        </div>
        <div v-if="message.content" class="px-3 py-2">
          <p class="text-sm break-words">{{ message.content }}</p>
        </div>
      </div>

      <!-- Audio Message (Voice Note) -->
      <div v-else-if="message.message_type === 'audio' && message.attachment" class="px-4 py-3">
        <div class="flex items-center gap-3">
          <!-- Avatar for voice note -->
          <div :class="['w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0', message.isMine ? 'bg-white/20' : 'bg-green-100 dark:bg-green-900/30']">
            <svg class="w-6 h-6" :class="message.isMine ? 'text-white' : 'text-green-600'" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z" clip-rule="evenodd" />
            </svg>
          </div>
          
          <button @click="toggleAudioPlay" :class="['w-10 h-10 rounded-full flex items-center justify-center transition-all transform hover:scale-105', message.isMine ? 'bg-white/20 hover:bg-white/30' : 'bg-green-500 hover:bg-green-600 text-white']">
            <svg v-if="!isPlaying" class="w-5 h-5 ml-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zM7 8a1 1 0 012 0v4a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v4a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
          </button>
          
          <div class="flex-1 min-w-0">
            <!-- Waveform visualization -->
            <div class="flex items-center gap-0.5 h-8">
              <div v-for="i in 25" :key="i" 
                :class="['w-1 rounded-full transition-all', message.isMine ? 'bg-white/40' : 'bg-gray-300 dark:bg-gray-600']"
                :style="{ 
                  height: `${Math.random() * 100}%`,
                  opacity: audioProgress > (i / 25 * 100) ? 1 : 0.4
                }">
              </div>
            </div>
            <div class="flex justify-between items-center mt-1">
              <span :class="['text-xs', message.isMine ? 'text-white/70' : 'text-gray-500']">
                {{ formatDuration(currentTime || 0) }}
              </span>
              <span :class="['text-xs', message.isMine ? 'text-white/70' : 'text-gray-500']">
                {{ formatDuration(message.attachment.duration || 0) }}
              </span>
            </div>
          </div>
          <audio ref="audioPlayer" :src="message.attachment.file_url" @timeupdate="updateAudioProgress" @ended="onAudioEnded" @loadedmetadata="onAudioLoaded" />
        </div>
      </div>

      <!-- Document Message -->
      <div v-else-if="message.message_type === 'document' && message.attachment" class="px-4 py-3">
        <div class="flex items-center gap-3 p-2 rounded-lg" :class="message.isMine ? 'bg-white/10' : 'bg-gray-50 dark:bg-gray-700'">
          <div :class="['w-12 h-12 rounded-xl flex items-center justify-center', getDocumentColor(message.attachment.mime_type)]">
            <span class="text-2xl">{{ getFileIcon(message.attachment.mime_type) }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{{ message.attachment.file_name }}</p>
            <p :class="['text-xs mt-0.5', message.isMine ? 'text-white/70' : 'text-gray-500']">
              {{ formatFileSize(message.attachment.file_size) }} • {{ getFileExtension(message.attachment.file_name) }}
            </p>
          </div>
          <a :href="message.attachment.file_url" download target="_blank" 
            :class="['w-10 h-10 rounded-full flex items-center justify-center transition-all hover:scale-105', message.isMine ? 'bg-white/20 hover:bg-white/30' : 'bg-green-500 hover:bg-green-600 text-white']">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
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
        <p class="text-sm break-words whitespace-pre-wrap leading-relaxed">{{ message.content }}</p>
      </div>

      <!-- Timestamp & Read Status -->
      <div :class="['px-3 pb-1.5 flex items-center gap-1 justify-end', message.isMine ? 'text-white/60' : 'text-gray-400']">
        <span class="text-[10px]">{{ formatTime(message.created_at || message.createdAt) }}</span>
        
        <!-- Read Status Icons (WhatsApp style) -->
        <template v-if="message.isMine">
          <!-- Sending -->
          <svg v-if="message.status === 'sending'" class="w-4 h-4 animate-pulse" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
          </svg>
          <!-- Sent (single check) -->
          <svg v-else-if="message.status === 'sent' || !message.read_at" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
          <!-- Delivered (double check gray) -->
          <div v-else-if="message.status === 'delivered'" class="flex -space-x-2">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </div>
          <!-- Read (double check blue) -->
          <div v-else-if="message.read_at" class="flex -space-x-2">
            <svg class="w-4 h-4 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
            <svg class="w-4 h-4 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </div>
        </template>
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
const currentTime = ref(0)
const audioDuration = ref(0)

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
  currentTime.value = audioPlayer.value.currentTime
  const progress = (audioPlayer.value.currentTime / audioPlayer.value.duration) * 100
  audioProgress.value = progress || 0
}

const onAudioEnded = () => {
  isPlaying.value = false
  audioProgress.value = 0
  currentTime.value = 0
}

const onAudioLoaded = () => {
  if (audioPlayer.value) {
    audioDuration.value = audioPlayer.value.duration
  }
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
  if (!bytes) return ''
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const getFileExtension = (filename: string) => {
  return filename?.split('.').pop()?.toUpperCase() || 'FILE'
}

const getFileIcon = (mimeType: string) => {
  if (mimeType?.includes('pdf')) return '📄'
  if (mimeType?.includes('word')) return '📝'
  if (mimeType?.includes('excel') || mimeType?.includes('spreadsheet')) return '📊'
  if (mimeType?.includes('powerpoint') || mimeType?.includes('presentation')) return '📑'
  if (mimeType?.includes('zip') || mimeType?.includes('rar') || mimeType?.includes('7z')) return '🗜️'
  if (mimeType?.includes('text')) return '📃'
  return '📎'
}

const getDocumentColor = (mimeType: string) => {
  if (mimeType?.includes('pdf')) return 'bg-red-100 dark:bg-red-900/30'
  if (mimeType?.includes('word')) return 'bg-blue-100 dark:bg-blue-900/30'
  if (mimeType?.includes('excel') || mimeType?.includes('spreadsheet')) return 'bg-green-100 dark:bg-green-900/30'
  if (mimeType?.includes('powerpoint')) return 'bg-orange-100 dark:bg-orange-900/30'
  if (mimeType?.includes('zip')) return 'bg-yellow-100 dark:bg-yellow-900/30'
  return 'bg-gray-100 dark:bg-gray-700'
}
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.animate-fade-in {
  animation: fade-in 0.25s ease-out;
}
</style>
