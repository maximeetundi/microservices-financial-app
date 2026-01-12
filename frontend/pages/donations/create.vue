<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-3xl mx-auto py-8 px-4">
        <!-- Header -->
        <div class="mb-8">
            <button @click="navigateTo('/donations')" class="text-gray-500 hover:text-indigo-600 mb-4 flex items-center">
                <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
                Annuler
            </button>
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Créer une Cagnotte 🚀</h1>
            <p class="text-gray-500 mt-2">Mobilisez votre communauté autour d'un projet qui vous tient à cœur.</p>
        </div>

        <div class="bg-white dark:bg-slate-900 rounded-3xl shadow-xl border border-gray-100 dark:border-gray-800 overflow-hidden">
            <div class="p-8">
                <form @submit.prevent="createCampaign" class="space-y-6">
                    
                    <!-- Title -->
                    <div>
                        <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Titre de la campagne</label>
                        <input v-model="form.title" type="text" class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="Ex: Soutien pour l'école locale..." required>
                    </div>

                    <!-- Description -->
                    <div>
                        <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Histoire & Description</label>
                        <textarea v-model="form.description" rows="5" class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="Racontez pourquoi cette cause est importante..." required></textarea>
                    </div>

                    <!-- Media Uploads -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <!-- Image Upload -->
                        <div>
                            <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Image de couverture 📷</label>
                            <input type="file" ref="coverInput" accept="image/*" class="hidden" @change="handleImageUpload">
                            
                            <div v-if="!coverImagePreview && !form.imageUrl" 
                                 @click="$refs.coverInput.click()"
                                 class="border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-xl p-8 text-center cursor-pointer hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/10 transition-all h-48 flex flex-col items-center justify-center">
                                <div class="text-4xl mb-2">🖼️</div>
                                <span class="text-sm text-gray-500 font-medium">Cliquez pour importer une image</span>
                            </div>

                            <div v-else class="relative h-48 rounded-xl overflow-hidden group">
                                <img :src="coverImagePreview || form.imageUrl" class="w-full h-full object-cover">
                                <button type="button" @click="removeImage" class="absolute top-2 right-2 bg-red-500 text-white p-2 rounded-full shadow-lg opacity-0 group-hover:opacity-100 transition-opacity">
                                    ✕
                                </button>
                            </div>
                        </div>

                        <!-- Video Upload -->
                        <div>
                            <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Vidéo de présentation 🎥 (Optionnel)</label>
                            <input type="file" ref="videoInput" accept="video/*" class="hidden" @change="handleVideoUpload">
                            
                            <div v-if="!videoPreview && !form.videoUrl" 
                                 @click="$refs.videoInput.click()"
                                 class="border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-xl p-8 text-center cursor-pointer hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/10 transition-all h-48 flex flex-col items-center justify-center">
                                <div class="text-4xl mb-2">🎬</div>
                                <span class="text-sm text-gray-500 font-medium">Cliquez pour importer une vidéo</span>
                            </div>

                            <div v-else class="relative h-48 rounded-xl overflow-hidden bg-black group">
                                <video :src="videoPreview || form.videoUrl" class="w-full h-full object-contain" controls></video>
                                <button type="button" @click="removeVideo" class="absolute top-2 right-2 bg-red-500 text-white p-2 rounded-full shadow-lg opacity-0 group-hover:opacity-100 transition-opacity z-10">
                                    ✕
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Target Amount -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Objectif (Optionnel)</label>
                            <div class="relative">
                                <input v-model.number="form.targetAmount" type="number" class="w-full pl-4 pr-12 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="0">
                                <span class="absolute right-4 top-3.5 text-gray-500 font-bold">XOF</span>
                            </div>
                        </div>
                        <div>
                            <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Devise</label>
                            <select v-model="form.currency" class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500">
                                <option value="XOF">XOF (CFA)</option>
                                <option value="USD">USD ($)</option>
                                <option value="EUR">EUR (€)</option>
                            </select>
                        </div>
                    </div>

                    <!-- Dynamic Form Builder -->
                    <div class="border-t border-gray-100 dark:border-gray-800 pt-6">
                        <div class="flex justify-between items-center mb-4">
                            <div>
                                <h3 class="text-lg font-bold text-gray-900 dark:text-white">📝 Formulaire Donateur</h3>
                                <p class="text-xs text-gray-500">Quelles infos voulez-vous demander ?</p>
                            </div>
                            <button type="button" @click="addField" class="text-indigo-600 hover:text-indigo-700 font-semibold text-sm">
                                + Ajouter un champ
                            </button>
                        </div>
                        
                        <div class="space-y-3">
                            <div v-for="(field, index) in form.form_fields" :key="index" class="p-4 bg-gray-50 dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-gray-700">
                                <div class="flex gap-3 items-start">
                                    <div class="flex-1 grid grid-cols-1 sm:grid-cols-2 gap-3">
                                        <input v-model="field.label" type="text" placeholder="Nom du champ (Ex: Taille T-shirt)" class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-slate-900 text-sm">
                                        <select v-model="field.type" class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-slate-900 text-sm">
                                            <option value="text">Texte</option>
                                            <option value="email">Email</option>
                                            <option value="phone">Téléphone</option>
                                            <option value="number">Nombre</option>
                                            <option value="select">Liste déroulante</option>
                                            <option value="checkbox">Case à cocher</option>
                                        </select>
                                    </div>
                                    <label class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mt-2">
                                        <input type="checkbox" v-model="field.required" class="rounded text-indigo-600">
                                        Obligatoire
                                    </label>
                                    <button type="button" @click="removeField(index)" class="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors">✕</button>
                                </div>
                                <!-- Options for select -->
                                <div v-if="field.type === 'select'" class="mt-3 pl-4 border-l-2 border-indigo-500">
                                    <label class="block text-xs text-gray-500 mb-1">Options (séparées par une virgule)</label>
                                    <input v-model="field.optionsText" @input="updateFieldOptions(field)" type="text" placeholder="Rouge, Vert, Bleu" class="w-full px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-slate-900 text-sm">
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Options -->
                    <div class="space-y-4 pt-4 border-t border-gray-100 dark:border-gray-800">
                        <label class="flex items-center gap-3 p-3 rounded-xl border border-gray-200 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                            <div class="bg-indigo-100 dark:bg-indigo-900/30 p-2 rounded-lg text-indigo-600">
                                🔄
                            </div>
                            <div class="flex-1">
                                <span class="block font-bold text-gray-900 dark:text-white">Autoriser les dons récurrents</span>
                                <span class="block text-xs text-gray-500">Permettre aux donateurs de donner mensuellement.</span>
                            </div>
                            <input type="checkbox" v-model="form.allowRecurring" class="w-6 h-6 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                        </label>
                        
                        <label class="flex items-center gap-3 p-3 rounded-xl border border-gray-200 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                            <div class="bg-indigo-100 dark:bg-indigo-900/30 p-2 rounded-lg text-indigo-600">
                                🕵️
                            </div>
                            <div class="flex-1">
                                <span class="block font-bold text-gray-900 dark:text-white">Autoriser les dons anonymes</span>
                                <span class="block text-xs text-gray-500">Les donateurs peuvent masquer leur identité.</span>
                            </div>
                            <input type="checkbox" v-model="form.allowAnonymous" class="w-6 h-6 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                        </label>
                    </div>

                    <!-- Submit -->
                    <div class="pt-6">
                        <button type="submit" :disabled="loading" class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold text-lg shadow-xl shadow-indigo-600/20 transition-all hover:scale-[1.02] active:scale-95 disabled:opacity-50 disabled:scale-100">
                            {{ loading ? 'Création de la campagne...' : 'Lancer la campagne 🚀' }}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useApi } from '~/composables/useApi'
const { donationApi } = useApi()

definePageMeta({
  middleware: 'auth'
})

const loading = ref(false)
const coverImageFile = ref<File | null>(null)
const coverImagePreview = ref<string | null>(null)
const videoFile = ref<File | null>(null)
const videoPreview = ref<string | null>(null)

const form = reactive({
    title: '',
    description: '',
    targetAmount: null as number | null,
    currency: 'XOF',
    imageUrl: '',
    videoUrl: '', // New field
    allowRecurring: true,
    allowAnonymous: true,
    form_fields: [
        { name: 'full_name', label: 'Nom complet', type: 'text', required: true, options: [] as string[], optionsText: '' },
        { name: 'email', label: 'Email', type: 'email', required: true, options: [] as string[], optionsText: '' },
    ]
})

const addField = () => {
  form.form_fields.push({
    name: `field_${Date.now()}`,
    label: '',
    type: 'text',
    required: false,
    options: [],
    optionsText: ''
  })
}

const removeField = (index: number) => {
  form.form_fields.splice(index, 1)
}

const updateFieldOptions = (field: any) => {
  if (field.optionsText) {
    field.options = field.optionsText.split(',').map((opt: string) => opt.trim()).filter((opt: string) => opt)
  } else {
    field.options = []
  }
}

const handleImageUpload = (event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) {
        const file = input.files[0]
        coverImageFile.value = file
        const reader = new FileReader()
        reader.onload = (e) => {
            coverImagePreview.value = e.target?.result as string
        }
        reader.readAsDataURL(file)
    }
}

const removeImage = () => {
    coverImageFile.value = null
    coverImagePreview.value = null
    form.imageUrl = ''
}

const handleVideoUpload = (event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) {
        const file = input.files[0]
        // Check size (50MB limit logic can be here or backend)
        if (file.size > 50 * 1024 * 1024) {
             alert("La vidéo est trop volumineuse (Max 50MB)")
             return
        }
        videoFile.value = file
        videoPreview.value = URL.createObjectURL(file)
    }
}

const removeVideo = () => {
    videoFile.value = null
    videoPreview.value = null
    form.videoUrl = ''
}

const createCampaign = async () => {
    loading.value = true
    try {
        let finalImageUrl = form.imageUrl
        let finalVideoUrl = form.videoUrl

        // Upload Image
        if (coverImageFile.value) {
            try {
                const res = await donationApi.uploadMedia(coverImageFile.value)
                if (res.data?.url) finalImageUrl = res.data.url
            } catch (e) {
                console.error("Image upload failed", e)
            }
        }

        // Upload Video
        if (videoFile.value) {
            try {
                const res = await donationApi.uploadMedia(videoFile.value)
                if (res.data?.url) finalVideoUrl = res.data.url
            } catch (e) {
                console.error("Video upload failed", e)
            }
        }

        const payload = {
            title: form.title,
            description: form.description,
            target_amount: form.targetAmount || 0,
            currency: form.currency,
            image_url: finalImageUrl,
            video_url: finalVideoUrl,
            allow_recurring: form.allowRecurring,
            allow_anonymous: form.allowAnonymous,
            form_schema: form.form_fields.map(f => ({
                name: f.label ? f.label.toLowerCase().replace(/\s+/g, '_') : f.name,
                label: f.label,
                type: f.type,
                required: f.required,
                options: f.options
            }))
        }
        
        const res = await donationApi.createCampaign(payload)
        const campaign = res.data?.campaign || res.data
        
        navigateTo(`/donations/${campaign.id}`)
        
    } catch (e: any) {
        alert(e.response?.data?.error || "Erreur lors de la création")
    } finally {
        loading.value = false
    }
}
</script>
