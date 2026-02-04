<template>
  <div class="max-w-7xl mx-auto pb-16">
    <div class="mb-8">
      <NuxtLink :to="`/shops/${slug}`" class="text-sm text-gray-500 hover:text-indigo-600 mb-2 inline-block">← Retour à la boutique</NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Toutes les catégories</h1>
      <p class="text-gray-500 mt-2">Explorez nos collections et trouvez ce dont vous avez besoin</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
       <div v-for="i in 8" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 rounded-2xl h-48"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="categories.length === 0" class="text-center py-24 bg-gray-50 dark:bg-slate-900 rounded-3xl">
        <div class="text-6xl mb-4">📂</div>
        <p class="text-gray-500 font-medium">Aucune catégorie disponible</p>
    </div>

    <!-- Grid -->
    <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
       <NuxtLink 
         v-for="cat in categories" 
         :key="cat.id"
         :to="`/shops/${slug}?category=${cat.slug}`"
         class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
       >
          <!-- Image -->
          <div class="aspect-[4/3] bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
             <img v-if="cat.image_url" :src="cat.image_url" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" alt="">
             <div v-else class="w-full h-full flex items-center justify-center text-4xl bg-gradient-to-br from-indigo-50 to-purple-50 dark:from-indigo-900/20 dark:to-purple-900/20 text-indigo-300">
               🗂️
             </div>
             
             <!-- Count Badge -->
             <div class="absolute top-3 right-3 px-3 py-1 bg-black/60 backdrop-blur text-white text-xs font-bold rounded-full">
                {{ cat.product_count }} articles
             </div>
          </div>
          
          <!-- Info -->
          <div class="p-6">
             <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2 group-hover:text-indigo-600 transition-colors">{{ cat.name }}</h3>
             <p class="text-sm text-gray-500 line-clamp-2 mb-4">{{ cat.description }}</p>
             <span class="text-sm font-bold text-indigo-600 dark:text-indigo-400 group-hover:underline">Explorer →</span>
          </div>
       </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useShopApi, type Category } from '~/composables/useShopApi'

definePageMeta({
  layout: 'shop-store'
})

const route = useRoute()
const shopApi = useShopApi()

const slug = computed(() => route.params.slug as string)
const loading = ref(true)
const categories = ref<Category[]>([])

const loadCategories = async () => {
  loading.value = true
  try {
    const result = await shopApi.listCategories(slug.value)
    categories.value = result.categories || []
  } catch (e) {
    console.error('Failed to load categories', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>
