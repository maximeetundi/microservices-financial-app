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

                    <!-- Target Amount -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Objectif (Optionnel)</label>
                            <div class="relative">
                                <input v-model.number="form.targetAmount" type="number" class="w-full pl-4 pr-12 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="0">
                                <span class="absolute right-4 top-3.5 text-gray-500 font-bold">XOF</span>
                            </div>
                            <p class="text-xs text-gray-400 mt-1">Laissez vide pour une cagnotte sans limite (0).</p>
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

                    <!-- Image URL (MVP) -->
                    <div>
                        <label class="block text-sm font-bold text-gray-900 dark:text-white mb-2">Image de couverture (URL)</label>
                        <input v-model="form.imageUrl" type="url" class="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500" placeholder="https://example.com/image.jpg">
                        <div v-if="form.imageUrl" class="mt-4 h-40 rounded-xl overflow-hidden bg-gray-100">
                            <img :src="form.imageUrl" class="w-full h-full object-cover">
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
                                <span class="block text-xs text-gray-500">Permettre aux donateurs de donner mensuellement ou annuellement.</span>
                            </div>
                            <input type="checkbox" v-model="form.allowRecurring" class="w-6 h-6 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                        </label>
                        
                        <label class="flex items-center gap-3 p-3 rounded-xl border border-gray-200 dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
                            <div class="bg-indigo-100 dark:bg-indigo-900/30 p-2 rounded-lg text-indigo-600">
                                🕵️
                            </div>
                            <div class="flex-1">
                                <span class="block font-bold text-gray-900 dark:text-white">Autoriser les dons anonymes</span>
                                <span class="block text-xs text-gray-500">Les donateurs peuvent choisir de masquer leur identité.</span>
                            </div>
                            <input type="checkbox" v-model="form.allowAnonymous" class="w-6 h-6 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500">
                        </label>
                    </div>

                    <!-- Submit -->
                    <div class="pt-6">
                        <button type="submit" :disabled="loading" class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold text-lg shadow-xl shadow-indigo-600/20 transition-all hover:scale-[1.02] active:scale-95 disabled:opacity-50 disabled:scale-100">
                            {{ loading ? 'Création...' : 'Lancer la campagne 🚀' }}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useApi } from '~/composables/useApi'
const { donationApi } = useApi()

definePageMeta({
  middleware: 'auth'
})

const loading = ref(false)
const form = reactive({
    title: '',
    description: '',
    targetAmount: null as number | null,
    currency: 'XOF',
    imageUrl: '',
    allowRecurring: true,
    allowAnonymous: true
})

const createCampaign = async () => {
    loading.value = true
    try {
        const payload = {
            title: form.title,
            description: form.description,
            target_amount: form.targetAmount || 0,
            currency: form.currency,
            image_url: form.imageUrl,
            allow_recurring: form.allowRecurring,
            allow_anonymous: form.allowAnonymous
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
