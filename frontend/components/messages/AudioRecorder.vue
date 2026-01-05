<template>
  <div v-if="isRecording" class="flex items-center gap-3 px-4 py-3 bg-red-50 dark:bg-red-900/20 border-t border-red-200 dark:border-red-800">
    <div class="flex-1 flex items-center gap-3">
      <!-- Recording Animation -->
      <div class="flex items-center gap-2">
        <div class="w-3 h-3 rounded-full bg-red-500 animate-pulse"></div>
        <span class="text-sm font-medium text-red-600 dark:text-red-400">Enregistrement...</span>
      </div>

      <!-- Timer -->
      <span class="text-sm text-gray-600 dark:text-gray-400">{{ formatRecordingTime }}</span>

      <!-- Waveform visualization (simple) -->
      <div class="flex-1 flex items-center gap-1 h-8">
        <div v-for="i in 20" :key="i" class="flex-1 bg-red-400 rounded-full transition-all" :style="{ height: Math.random() * 100 + '%' }"></div>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex gap-2">
      <button @click="$emit('cancel')" class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 transition-colors">
        <svg class="w-5 h-5 text-gray-700 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
      <button @click="stopRecording" class="p-2 rounded-full bg-green-500 hover:bg-green-600 transition-colors">
        <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  isRecording: boolean
}>()

const emit = defineEmits(['audioRecorded', 'cancel'])

let mediaRecorder: MediaRecorder | null = null
let audioChunks: Blob[] = []
const recordingStartTime = ref(0)
const recordingDuration = ref(0)
let timer: NodeJS.Timeout | null = null

const formatRecordingTime = computed(() => {
  const seconds = Math.floor(recordingDuration.value / 1000)
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
})

const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder = new MediaRecorder(stream)
    audioChunks = []

    mediaRecorder.ondataavailable = (e) => {
      audioChunks.push(e.data)
    }

    mediaRecorder.onstop = () => {
      const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })
      const duration = recordingDuration.value / 1000 // in seconds
      emit('audioRecorded', { blob: audioBlob, duration })
      
      // Stop all tracks
      stream.getTracks().forEach(track => track.stop())
    }

    mediaRecorder.start()
    recordingStartTime.value = Date.now()
    
    // Update timer every 100ms
    timer = setInterval(() => {
      recordingDuration.value = Date.now() - recordingStartTime.value
      
      // Max 60 seconds
      if (recordingDuration.value >= 60000) {
        stopRecording()
      }
    }, 100)
  } catch (err) {
    console.error('Failed to start recording:', err)
    alert('Impossible d\'accéder au microphone. Veuillez autoriser l\'accès.')
    emit('cancel')
  }
}

const stopRecording = () => {
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
  }
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

onMounted(() => {
  if (props.isRecording) {
    startRecording()
  }
})

onUnmounted(() => {
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
  }
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>
