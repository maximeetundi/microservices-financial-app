<template>
  <div
    class="group bg-white dark:bg-slate-900 rounded-xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:border-indigo-200 dark:hover:border-indigo-800 transition-all duration-300"
  >
    <!-- Image Container -->
    <div class="relative aspect-square overflow-hidden bg-gray-50 dark:bg-slate-800">
      <NuxtLink :to="`/shops/${shopSlug}/product/${product.slug}`">
        <img
          v-if="product.images?.length"
          :src="product.images[0]"
          :alt="product.name"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          loading="lazy"
        >
        <div v-else class="w-full h-full flex items-center justify-center text-5xl text-gray-300 dark:text-gray-600">
          📦
        </div>
      </NuxtLink>

      <!-- Badges -->
      <div class="absolute top-2 left-2 flex flex-col gap-1.5">
        <!-- Promo Badge -->
        <span
          v-if="discountPercent > 0"
          class="px-2 py-1 bg-red-500 text-white text-[10px] font-bold rounded-md shadow-sm"
        >
          -{{ discountPercent }}%
        </span>
        <!-- Featured Badge -->
        <span
          v-if="product.is_featured"
          class="px-2 py-1 bg-amber-500 text-white text-[10px] font-bold rounded-md shadow-sm"
        >
          ⭐ Vedette
        </span>
        <!-- New Badge -->
        <span
          v-if="isNew"
          class="px-2 py-1 bg-emerald-500 text-white text-[10px] font-bold rounded-md shadow-sm"
        >
          Nouveau
        </span>
        <!-- Out of Stock -->
        <span
          v-if="product.stock === 0"
          class="px-2 py-1 bg-gray-800 text-white text-[10px] font-bold rounded-md shadow-sm"
        >
          Épuisé
        </span>
      </div>

      <!-- Favorite Button -->
      <button
        @click.prevent="toggleFavorite"
        class="absolute top-2 right-2 w-8 h-8 rounded-full bg-white/90 dark:bg-slate-800/90 shadow-sm flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-200 hover:scale-110"
        :class="isFavorite ? 'text-red-500' : 'text-gray-400 hover:text-red-500'"
      >
        <svg v-if="isFavorite" class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
        </svg>
      </button>

      <!-- Quick Add to Cart (appears on hover) -->
      <div
        class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-2 group-hover:translate-y-0"
      >
        <button
          v-if="product.stock !== 0"
          @click.prevent="addToCart"
          class="w-full py-2 bg-white text-indigo-600 font-semibold text-sm rounded-lg hover:bg-indigo-50 transition-colors flex items-center justify-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
          </svg>
          Ajouter
        </button>
        <span v-else class="block w-full py-2 text-center text-white/80 text-sm font-medium">
          Rupture de stock
        </span>
      </div>
    </div>

    <!-- Info -->
    <div class="p-3">
      <!-- Category Tag -->
      <span
        v-if="product.category_name"
        class="inline-block text-[10px] font-medium text-indigo-600 dark:text-indigo-400 uppercase tracking-wider mb-1"
      >
        {{ product.category_name }}
      </span>

      <!-- Name -->
      <NuxtLink :to="`/shops/${shopSlug}/product/${product.slug}`">
        <h3 class="font-semibold text-gray-900 dark:text-white text-sm mb-1.5 line-clamp-2 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors leading-tight">
          {{ product.name }}
        </h3>
      </NuxtLink>

      <!-- Rating -->
      <div v-if="product.rating || product.reviews_count" class="flex items-center gap-1.5 mb-2">
        <div class="flex text-amber-400 text-xs">
          <span v-for="i in 5" :key="i">{{ i <= Math.round(product.rating || 0) ? '★' : '☆' }}</span>
        </div>
        <span class="text-xs text-gray-400">({{ product.reviews_count || 0 }})</span>
      </div>

      <!-- Price -->
      <div class="flex items-end gap-2">
        <span class="text-lg font-bold text-indigo-600 dark:text-indigo-400">
          {{ formatPrice(product.price) }}
        </span>
        <span
          v-if="product.compare_at_price && product.compare_at_price > product.price"
          class="text-xs text-gray-400 line-through mb-0.5"
        >
          {{ formatPrice(product.compare_at_price) }}
        </span>
      </div>

      <!-- Mobile Add to Cart Button -->
      <button
        v-if="product.stock !== 0"
        @click="addToCart"
        class="md:hidden w-full mt-3 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-lg transition-colors flex items-center justify-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
        </svg>
        Ajouter au panier
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Product } from '~/composables/useShopApi'

interface Props {
  product: Product
  shopSlug: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'toggle-favorite', productId: string): void
  (e: 'add-to-cart', product: Product): void
}>()

// Check if product is new (created within last 7 days)
const isNew = computed(() => {
  if (!props.product.created_at) return false
  const created = new Date(props.product.created_at)
  const now = new Date()
  const diffDays = Math.floor((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays <= 7
})

// Calculate discount percentage
const discountPercent = computed(() => {
  if (!props.product.compare_at_price || props.product.compare_at_price <= props.product.price) {
    return 0
  }
  return Math.round(((props.product.compare_at_price - props.product.price) / props.product.compare_at_price) * 100)
})

// Check if product is in favorites
const isFavorite = computed(() => {
  if (typeof window === 'undefined') return false
  try {
    const favs = JSON.parse(localStorage.getItem(`shop_favorites_${props.shopSlug}`) || '[]')
    return favs.includes(props.product.id)
  } catch {
    return false
  }
})

// Format price
const formatPrice = (amount: number) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: props.product.currency || 'XOF',
    minimumFractionDigits: 0
  }).format(amount)
}

// Toggle favorite
const toggleFavorite = () => {
  emit('toggle-favorite', props.product.id)
}

// Add to cart
const addToCart = () => {
  emit('add-to-cart', props.product)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
