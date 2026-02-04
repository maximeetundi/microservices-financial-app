<template>
  <!-- List Mode -->
  <div 
    v-if="listMode"
    @click="navigateTo(`/shops/${shopSlug}/product/${product.slug}`)"
    class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:border-indigo-500/50 transition-all duration-300 cursor-pointer flex"
  >
    <!-- Image -->
    <div class="w-48 h-48 bg-gray-100 dark:bg-slate-800 flex-shrink-0 relative overflow-hidden">
      <img 
        v-if="product.images?.length" 
        :src="product.images[0]" 
        class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" 
        :alt="product.name"
      >
      <div v-else class="w-full h-full flex items-center justify-center text-4xl text-gray-300">📦</div>
      
      <!-- Badges -->
      <div class="absolute top-2 left-2 flex flex-col gap-1">
        <span v-if="product.is_featured" class="px-2 py-0.5 bg-amber-500 text-white text-xs font-bold rounded">⭐ Vedette</span>
        <span v-if="product.compare_at_price" class="px-2 py-0.5 bg-red-500 text-white text-xs font-bold rounded">-{{ discountPercent }}%</span>
      </div>
    </div>
    
    <!-- Info -->
    <div class="flex-1 p-6 flex flex-col">
      <div class="flex-1">
        <!-- Tags -->
        <div v-if="product.tags?.length" class="flex gap-1 mb-2 flex-wrap">
          <span 
            v-for="tag in product.tags.slice(0, 3)" 
            :key="tag" 
            class="text-xs px-2 py-0.5 bg-indigo-50 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 rounded"
          >
            #{{ tag }}
          </span>
        </div>

        <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2 group-hover:text-indigo-600 transition-colors">
          {{ product.name }}
        </h3>
        
        <p class="text-sm text-gray-500 line-clamp-2 mb-4">{{ product.description }}</p>
        
        <!-- Rating -->
        <div v-if="product.average_rating" class="flex items-center gap-2 mb-4">
          <div class="flex text-amber-400">
            <span v-for="i in 5" :key="i">{{ i <= Math.round(product.average_rating) ? '★' : '☆' }}</span>
          </div>
          <span class="text-sm text-gray-500">({{ product.review_count || 0 }} avis)</span>
        </div>
      </div>
      
      <!-- Price & Actions -->
      <div class="flex items-center justify-between">
        <div>
          <span class="text-2xl font-bold text-indigo-600 dark:text-indigo-400">
            {{ formatPrice(product.price, product.currency) }}
          </span>
          <span v-if="product.compare_at_price" class="ml-2 text-gray-400 line-through">
            {{ formatPrice(product.compare_at_price, product.currency) }}
          </span>
        </div>
        <div class="flex gap-2">
          <button 
            @click.stop="$emit('toggle-favorite', product.id)"
            class="w-12 h-12 rounded-xl border border-gray-200 dark:border-gray-700 flex items-center justify-center hover:border-red-500 hover:text-red-500 transition-colors"
            :class="isFavorite ? 'text-red-500 border-red-500 bg-red-50 dark:bg-red-900/20' : 'text-gray-400'"
          >
            {{ isFavorite ? '❤️' : '🤍' }}
          </button>
          <button 
            @click.stop="addToCart"
            :disabled="product.stock === 0"
            class="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            🛒 Ajouter
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Grid Mode (Default) -->
  <div 
    v-else
    @click="navigateTo(`/shops/${shopSlug}/product/${product.slug}`)"
    class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-2xl hover:-translate-y-2 hover:border-indigo-500/50 transition-all duration-300 cursor-pointer flex flex-col h-full"
  >
    <!-- Image -->
    <div class="aspect-square bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
      <img 
        v-if="product.images?.length" 
        :src="product.images[0]" 
        class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" 
        :alt="product.name"
      >
      <div v-else class="w-full h-full flex items-center justify-center text-5xl text-gray-300">📦</div>
      
      <!-- Badges -->
      <div class="absolute top-3 left-3 flex flex-col gap-2">
        <div v-if="product.is_featured" class="px-2.5 py-1 bg-amber-500 text-white text-xs font-bold rounded-lg shadow-lg">
          ⭐ Vedette
        </div>
        <div v-if="product.compare_at_price" class="px-2.5 py-1 bg-red-500 text-white text-xs font-bold rounded-lg shadow-lg">
          -{{ discountPercent }}%
        </div>
        <div v-if="product.stock === 0" class="px-2.5 py-1 bg-gray-800 text-white text-xs font-bold rounded-lg">
          Épuisé
        </div>
      </div>

      <!-- Favorite Button -->
      <button 
        @click.stop="$emit('toggle-favorite', product.id)"
        class="absolute top-3 right-3 w-10 h-10 bg-white dark:bg-slate-800 rounded-full flex items-center justify-center shadow-lg opacity-0 group-hover:opacity-100 transition-all hover:scale-110"
        :class="isFavorite ? 'text-red-500' : 'text-gray-400 hover:text-red-500'"
      >
        {{ isFavorite ? '❤️' : '🤍' }}
      </button>

      <!-- Quick Actions -->
      <div class="absolute bottom-3 left-3 right-3 translate-y-10 opacity-0 group-hover:translate-y-0 group-hover:opacity-100 transition-all duration-300">
        <button 
          @click.stop="addToCart"
          :disabled="product.stock === 0"
          class="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl shadow-lg flex items-center justify-center gap-2 disabled:opacity-50"
        >
          🛒 Ajouter au panier
        </button>
      </div>
    </div>
    
    <!-- Info -->
    <div class="p-4 flex-1 flex flex-col">
      <!-- Tags -->
      <div v-if="product.tags?.length" class="flex gap-1 mb-2 flex-wrap">
        <span 
          v-for="tag in product.tags.slice(0, 2)" 
          :key="tag" 
          class="text-[10px] px-2 py-0.5 bg-indigo-50 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 rounded-md font-medium"
        >
          #{{ tag }}
        </span>
      </div>

      <!-- Name -->
      <h3 class="font-bold text-gray-900 dark:text-white text-base mb-2 line-clamp-2 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
        {{ product.name }}
      </h3>
      
      <!-- Rating -->
      <div v-if="product.average_rating" class="flex items-center gap-1 mb-2">
        <div class="flex text-amber-400 text-sm">
          <span v-for="i in 5" :key="i">{{ i <= Math.round(product.average_rating) ? '★' : '☆' }}</span>
        </div>
        <span class="text-xs text-gray-500">({{ product.review_count || 0 }})</span>
      </div>
      
      <!-- Price -->
      <div class="mt-auto pt-3 flex items-center justify-between">
        <div class="flex flex-col">
          <span class="text-lg font-bold text-indigo-600 dark:text-indigo-400">
            {{ formatPrice(product.price, product.currency) }}
          </span>
          <span v-if="product.compare_at_price" class="text-xs text-gray-400 line-through">
            {{ formatPrice(product.compare_at_price, product.currency) }}
          </span>
        </div>
        <button 
          @click.stop="addToCart"
          :disabled="product.stock === 0"
          class="w-10 h-10 rounded-xl bg-indigo-100 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 flex items-center justify-center hover:bg-indigo-600 hover:text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed font-bold text-lg"
        >
          +
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Product } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'

const props = defineProps<{
  product: Product
  shopSlug: string
  listMode?: boolean
}>()

const emit = defineEmits(['toggle-favorite'])
const cartStore = useCartStore()

const discountPercent = computed(() => {
  if (!props.product.compare_at_price) return 0
  return Math.round(((props.product.compare_at_price - props.product.price) / props.product.compare_at_price) * 100)
})

const isFavorite = computed(() => {
  if (typeof window === 'undefined') return false
  const favs = JSON.parse(localStorage.getItem(`shop_favorites_${props.shopSlug}`) || '[]')
  return favs.includes(props.product.id)
})

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const addToCart = () => {
  cartStore.addItem(props.product, 1)
}
</script>
