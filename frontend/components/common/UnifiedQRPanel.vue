<template>
  <div class="text-center">
    <!-- QR Code Container -->
    <div class="mb-6 relative inline-block group">
      <div class="bg-white p-4 rounded-2xl shadow-xl border border-gray-100 transition-transform duration-300 hover:scale-105">
        <div v-if="loading" class="w-48 h-48 flex items-center justify-center text-gray-300 bg-gray-50 rounded-lg">
          <div class="animate-pulse flex flex-col items-center">
            <span class="text-4xl mb-2">📱</span>
            <span class="text-xs">Génération...</span>
          </div>
        </div>
        <img v-else-if="qrSrc" :src="qrSrc" alt="QR Code" class="w-48 h-48 object-contain" />
        <div v-else class="w-48 h-48 flex items-center justify-center text-gray-300 bg-gray-50 rounded-lg">
          <span class="text-xs">Code non disponible</span>
        </div>
      </div>
    </div>
    
    <!-- Title/Subtitle -->
    <h3 v-if="title" class="text-lg font-bold text-gray-900 dark:text-white mb-2">{{ title }}</h3>
    <p v-if="subtitle" class="text-sm text-gray-500 dark:text-gray-400 mb-6">{{ subtitle }}</p>
    
    <!-- Code Display -->
    <div v-if="code" class="mb-8 font-mono bg-indigo-50 dark:bg-slate-800 text-indigo-700 dark:text-indigo-300 py-3 px-4 rounded-xl border border-indigo-100 dark:border-indigo-900/50 flex items-center justify-between group cursor-pointer hover:bg-indigo-100 dark:hover:bg-slate-700 transition-colors" @click="copyCode">
      <span class="font-bold text-lg tracking-wider">{{ code }}</span>
      <span class="text-indigo-400 group-hover:text-indigo-600 dark:group-hover:text-white transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
        </svg>
      </span>
    </div>
    
    <!-- Actions -->
    <div class="grid grid-cols-2 gap-3">
      <button @click="download" class="flex flex-col items-center justify-center p-3 rounded-xl bg-gray-50 hover:bg-gray-100 dark:bg-slate-800 dark:hover:bg-slate-700 text-gray-700 dark:text-gray-300 transition-colors border border-gray-100 dark:border-gray-700 group">
        <span class="mb-1 p-2 bg-white dark:bg-slate-700 rounded-full shadow-sm group-hover:scale-110 transition-transform">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </span>
        <span class="text-xs font-semibold">Télécharger</span>
      </button>
      
      <button @click="share" class="flex flex-col items-center justify-center p-3 rounded-xl bg-gray-50 hover:bg-gray-100 dark:bg-slate-800 dark:hover:bg-slate-700 text-gray-700 dark:text-gray-300 transition-colors border border-gray-100 dark:border-gray-700 group">
        <span class="mb-1 p-2 bg-white dark:bg-slate-700 rounded-full shadow-sm group-hover:scale-110 transition-transform">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
        </span>
        <span class="text-xs font-semibold">Partager</span>
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  qrSrc: {
    type: String,
    default: ''
  },
  code: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['copy-code', 'download', 'share'])

const copyCode = () => {
    emit('copy-code')
}

const download = () => {
  if (!props.qrSrc) return
  emit('download')
}

const share = () => {
  emit('share')
}
</script>
