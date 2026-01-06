<template>
  <div class="video-player-wrapper relative rounded-lg overflow-hidden bg-black">
    <!-- Thumbnail/Preview with Play Button -->
    <div v-if="!showPlayer" class="relative cursor-pointer" @click="loadPlayer">
      <!-- Video Thumbnail or Placeholder -->
      <div class="aspect-video bg-gray-800 flex items-center justify-center">
        <video 
          v-if="src" 
          :src="src" 
          class="w-full h-full object-cover" 
          preload="metadata"
          muted
          @loadeddata="handleVideoLoaded"
        ></video>
        <div v-else class="text-gray-500 flex flex-col items-center">
          <svg class="w-12 h-12" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
          <span class="text-sm mt-2">Vidéo</span>
        </div>
      </div>
      
      <!-- Play Button Overlay -->
      <div class="absolute inset-0 flex items-center justify-center bg-black/30 hover:bg-black/40 transition-colors">
        <div class="w-16 h-16 rounded-full bg-white/90 flex items-center justify-center shadow-lg transform hover:scale-110 transition-transform">
          <svg class="w-8 h-8 text-gray-800 ml-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
      </div>
      
      <!-- Duration Badge -->
      <div v-if="duration" class="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
        {{ formatDuration(duration) }}
      </div>
    </div>

    <!-- Plyr Video Player -->
    <div v-else class="aspect-video">
      <video 
        ref="videoElement" 
        :src="src" 
        class="plyr-video w-full h-full"
        playsinline
        @loadedmetadata="handleMetadata"
      >
        <source :src="src" :type="mimeType" />
      </video>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'

interface Props {
  src: string
  mimeType?: string
  poster?: string
}

const props = withDefaults(defineProps<Props>(), {
  mimeType: 'video/mp4',
  poster: ''
})

const showPlayer = ref(false)
const videoElement = ref<HTMLVideoElement | null>(null)
const duration = ref(0)
let player: any = null

const handleVideoLoaded = (e: Event) => {
  const video = e.target as HTMLVideoElement
  if (video.duration) {
    duration.value = video.duration
  }
}

const handleMetadata = (e: Event) => {
  const video = e.target as HTMLVideoElement
  if (video.duration) {
    duration.value = video.duration
  }
}

const formatDuration = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const loadPlayer = async () => {
  showPlayer.value = true
  
  await nextTick()
  
  if (videoElement.value && process.client) {
    try {
      const Plyr = (await import('plyr')).default
      
      player = new Plyr(videoElement.value, {
        controls: [
          'play-large', 
          'play', 
          'progress', 
          'current-time', 
          'mute', 
          'volume', 
          'fullscreen'
        ],
        settings: [],
        resetOnEnd: true,
        clickToPlay: true,
        hideControls: true,
        keyboard: { focused: true, global: false },
        tooltips: { controls: false, seek: true },
        i18n: {
          play: 'Lecture',
          pause: 'Pause',
          mute: 'Muet',
          unmute: 'Son',
          enterFullscreen: 'Plein écran',
          exitFullscreen: 'Quitter plein écran',
          seek: 'Chercher'
        }
      })
      
      // Auto-play when loaded
      player.on('ready', () => {
        player.play()
      })
    } catch (e) {
      console.error('Failed to load Plyr:', e)
    }
  }
}

onBeforeUnmount(() => {
  if (player) {
    player.destroy()
    player = null
  }
})
</script>

<style>
@import 'plyr/dist/plyr.css';

.video-player-wrapper .plyr {
  --plyr-color-main: #25D366;
  --plyr-video-background: #000;
}

.video-player-wrapper .plyr__control--overlaid {
  background: rgba(37, 211, 102, 0.9);
}

.video-player-wrapper .plyr__control--overlaid:hover {
  background: #25D366;
}
</style>
