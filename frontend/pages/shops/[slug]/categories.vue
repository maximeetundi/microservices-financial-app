<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Catégories</h1>
        <p class="text-gray-500">Explorez nos {{ categories.length }} catégories de produits</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div v-for="i in 8" :key="i" class="animate-pulse">
        <div class="bg-gray-200 dark:bg-slate-800 aspect-square rounded-2xl mb-4"></div>
        <div class="h-4 bg-gray-200 dark:bg-slate-800 rounded"></div>
      </div>
    </div>

    <!-- Categories Grid -->
    <div v-else-if="categories.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <NuxtLink 
        v-for="cat in categories" 
        :key="cat.id"
        :to="`/shops/${shopSlug}?category=${cat.slug}`"
        class="group relative bg-white dark:bg-slate-900 rounded-3xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-2xl hover:-translate-y-2 transition-all duration-300"
      >
        <!-- Image -->
        <div class="aspect-square bg-gradient-to-br from-indigo-100 to-purple-100 dark:from-slate-800 dark:to-slate-700 relative overflow-hidden">
          <img 
            v-if="cat.image_url" 
            :src="cat.image_url" 
            class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" 
            :alt="cat.name"
          >
          <div v-else class="w-full h-full flex items-center justify-center text-6xl group-hover:scale-125 transition-transform duration-300">
            {{ cat.icon || '📦' }}
          </div>
          
          <!-- Overlay -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
          
          <!-- Product Count Badge -->
          <div class="absolute top-4 right-4 px-3 py-1.5 bg-white/90 dark:bg-slate-800/90 backdrop-blur rounded-full text-sm font-bold text-gray-900 dark:text-white shadow-lg">
            {{ cat.product_count || 0 }} produits
          </div>
        </div>
        
        <!-- Info -->
        <div class="p-5">
          <h3 class="font-bold text-xl text-gray-900 dark:text-white mb-2 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
            {{ cat.name }}
          </h3>
          <p v-if="cat.description" class="text-sm text-gray-500 line-clamp-2">
            {{ cat.description }}
          </p>
          <div class="mt-4 flex items-center text-indigo-600 dark:text-indigo-400 text-sm font-medium group-hover:gap-2 transition-all">
            <span>Explorer</span>
            <span class="group-hover:translate-x-1 transition-transform">→</span>
          </div>
        </div>
      </NuxtLink>
    </div>

    <!-- Empty -->
    <div v-else class="text-center py-24 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-gray-800">
      <div class="text-6xl mb-4">📂</div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucune catégorie</h3>
      <p class="text-gray-500">Cette boutique n'a pas encore de catégories.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import { useShopApi, type Category } from '~/composables/useShopApi'

definePageMeta({
  layout: 'shop-layout'
})

const route = useRoute()
const shopApi = useShopApi()

const shopSlug = computed(() => route.params.slug as string)
const categories = inject<Ref<Category[]>>('categories', ref([]))
const loading = ref(false)

// If categories not injected, load them
onMounted(async () => {
  if (categories.value.length === 0) {
    loading.value = true
    try {
      const result = await shopApi.listCategories(shopSlug.value)
      categories.value = result.categories || []
    } catch (e) {
      console.error('Failed to load categories', e)
    } finally {
      loading.value = false
    }
  }
})
</script>
