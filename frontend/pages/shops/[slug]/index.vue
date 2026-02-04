<template>
  <div class="space-y-8">
    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-48 rounded-2xl"></div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="i in 8" :key="i" class="animate-pulse bg-gray-100 dark:bg-slate-800 h-64 rounded-xl"></div>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="text-center py-16">
      <div class="text-6xl mb-4">😞</div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white">Boutique non trouvée</h3>
      <NuxtLink to="/shops" class="mt-4 inline-block px-6 py-3 bg-indigo-600 text-white rounded-xl font-bold">
        Retour aux boutiques
      </NuxtLink>
    </div>

    <div v-else-if="shop">
      <!-- Banner -->
      <div class="relative h-48 md:h-64 rounded-2xl overflow-hidden mb-8 shadow-lg">
        <img v-if="shop.banner_url" :src="shop.banner_url" class="w-full h-full object-cover" alt="">
        <div v-else class="w-full h-full bg-gradient-to-br from-indigo-500 to-purple-600"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"></div>
        
        <!-- Shop Info -->
        <div class="absolute bottom-0 left-0 right-0 p-6 flex items-end gap-4">
          <div class="w-20 h-20 rounded-xl bg-white dark:bg-slate-800 border-4 border-white dark:border-slate-700 shadow-xl overflow-hidden flex-shrink-0">
            <img v-if="shop.logo_url" :src="shop.logo_url" class="w-full h-full object-cover" alt="">
            <div v-else class="w-full h-full flex items-center justify-center text-3xl bg-gradient-to-br from-indigo-500 to-purple-500 text-white">
              {{ shop.name.charAt(0) }}
            </div>
          </div>
          <div class="flex-1 text-white">
            <h1 class="text-2xl md:text-3xl font-bold mb-1">{{ shop.name }}</h1>
            <p v-if="shop.description" class="text-white/90 text-sm line-clamp-2 max-w-2xl">{{ shop.description }}</p>
          </div>
          
          <!-- QR Code -->
          <button v-if="shop.qr_code" @click="showQR = true" class="p-3 bg-white/20 backdrop-blur rounded-xl hover:bg-white/30 transition-colors border border-white/30 text-white">
            <span class="text-2xl">📱</span>
          </button>
        </div>
      </div>

      <!-- Categories Filter -->
      <div v-if="categories.length > 0" class="flex gap-2 overflow-x-auto pb-4 mb-6 scrollbar-hide">
        <button 
          @click="selectCategory('')"
          :class="[
            'px-5 py-2.5 rounded-full whitespace-nowrap font-medium transition-all duration-200 border',
            !selectedCategory 
              ? 'bg-indigo-600 text-white border-indigo-600 shadow-md transform scale-105' 
              : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300 border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700'
          ]"
        >
          Tous
        </button>
        <button 
          v-for="cat in categories" 
          :key="cat.id"
          @click="selectCategory(cat.slug)"
          :class="[
            'px-5 py-2.5 rounded-full whitespace-nowrap font-medium transition-all duration-200 border',
            selectedCategory === cat.slug 
              ? 'bg-indigo-600 text-white border-indigo-600 shadow-md transform scale-105' 
              : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300 border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700'
          ]"
        >
          {{ cat.name }} <span class="ml-1 opacity-75 text-xs">({{ cat.product_count }})</span>
        </button>
      </div>

      <!-- Products Grid -->
      <div v-if="products.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div 
          v-for="product in products" 
          :key="product.id"
          @click="openProduct(product)"
          class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-xl hover:-translate-y-1 hover:border-indigo-500/30 transition-all duration-300 cursor-pointer flex flex-col h-full"
        >
          <!-- Image -->
          <div class="aspect-[4/3] bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
            <img v-if="product.images?.length" :src="product.images[0]" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" alt="">
            <div v-else class="w-full h-full flex items-center justify-center text-4xl text-gray-300">📦</div>
            
            <!-- Badges -->
            <div class="absolute top-3 left-3 flex flex-col gap-2">
                 <div v-if="product.is_featured" class="px-2.5 py-1 bg-amber-500/90 backdrop-blur text-white text-xs font-bold rounded-lg shadow-sm">
                    ⭐ Featured
                 </div>
                 <div v-if="product.stock === 0" class="px-2.5 py-1 bg-red-500/90 backdrop-blur text-white text-xs font-bold rounded-lg shadow-sm">
                    Épuisé
                 </div>
            </div>

             <!-- Quick Actions (Hover) -->
             <div class="absolute bottom-3 right-3 translate-y-10 opacity-0 group-hover:translate-y-0 group-hover:opacity-100 transition-all duration-300 flex gap-2">
                 <div class="w-8 h-8 rounded-full bg-white dark:bg-slate-800 flex items-center justify-center shadow-md text-gray-600 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400">
                     👁️
                 </div>
             </div>
          </div>
          
          <!-- Info -->
          <div class="p-4 flex-1 flex flex-col">
             <!-- Tags? -->
             <div v-if="product.tags?.length" class="flex gap-1 mb-2 overflow-hidden">
                 <span v-for="tag in product.tags.slice(0,2)" :key="tag" class="text-[10px] px-2 py-0.5 bg-gray-100 dark:bg-slate-800 text-gray-500 dark:text-gray-400 rounded-md">
                     {{ tag }}
                 </span>
             </div>

            <h3 class="font-bold text-gray-900 dark:text-white text-base mb-1 line-clamp-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
              {{ product.name }}
            </h3>
            
            <div class="mt-auto pt-3 flex items-center justify-between">
              <div class="flex flex-col">
                  <span class="text-lg font-bold text-indigo-600 dark:text-indigo-400">
                    {{ formatPrice(product.price, product.currency) }}
                  </span>
                  <span v-if="product.compare_at_price" class="text-xs text-gray-400 line-through">
                    {{ formatPrice(product.compare_at_price, product.currency) }}
                  </span>
              </div>
              <button class="w-8 h-8 rounded-lg bg-indigo-50 dark:bg-slate-800 text-indigo-600 flex items-center justify-center hover:bg-indigo-600 hover:text-white transition-colors">
                  +
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty Products -->
      <div v-else class="text-center py-24 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-gray-800">
        <div class="w-24 h-24 mx-auto mb-6 bg-gray-50 dark:bg-slate-800 rounded-full flex items-center justify-center text-6xl">📦</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Aucun produit trouvé</h3>
        <p class="text-gray-500 max-w-xs mx-auto">Cette boutique n'a pas encore de produits dans cette catégorie.</p>
        <button v-if="selectedCategory" @click="selectedCategory = ''" class="mt-6 px-6 py-2 text-indigo-600 font-medium hover:bg-indigo-50 dark:hover:bg-slate-800 rounded-lg transition-colors">
            Voir tous les produits
        </button>
      </div>

    </div>

    <!-- QR Modal -->
    <ShareQRCodeModal
      v-model="showQR"
      :qr-data="windowLocation"
      :display-code="shop?.slug"
      :title="shop?.name"
      subtitle="Scanner pour visiter"
      :share-text="'Visitez la boutique ' + shop?.name"
    />
  </div>
</template>

<script setup lang="ts">
import { useShopApi, type Shop, type Product, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import ShareQRCodeModal from '~/components/common/ShareQRCodeModal.vue'

definePageMeta({
  layout: 'shop-store'
})

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const cartStore = useCartStore()

const slug = computed(() => route.params.slug as string)

const loading = ref(true)
const error = ref(false)
const shop = ref<Shop | null>(null)
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const showQR = ref(false)
const windowLocation = ref('')

// Computed from route for reactivity
const selectedCategory = computed(() => route.query.category as string || '')
const searchQuery = computed(() => route.query.search as string || '')

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const openProduct = (product: Product) => {
  navigateTo(`/shops/${slug.value}/product/${product.slug}`)
}

const selectCategory = (catSlug: string) => {
    // If clicking same category, clear it (toggle)
    const newCat = selectedCategory.value === catSlug ? undefined : catSlug
    router.push({ 
        path: route.path, 
        query: { 
            ...route.query, 
            category: newCat,
            // Reset search when changing category? Maybe not.
        } 
    })
}

const loadProducts = async () => {
    loading.value = true
    try {
        const options: any = {}
        if (selectedCategory.value) options.category = selectedCategory.value
        if (searchQuery.value) options.search = searchQuery.value
        
        // Note: listProducts(shopSlug, page, pageSize, options)
        const result = await shopApi.listProducts(slug.value, 1, 100, options)
        products.value = result.products || []
    } catch (e) {
        console.error('Failed to load products', e)
    } finally {
        loading.value = false
    }
}

const loadShop = async () => {
  loading.value = true
  error.value = false
  try {
    shop.value = await shopApi.getShop(slug.value)
    if (typeof window !== 'undefined') {
        windowLocation.value = window.location.href
    }
    
    // Load categories
    const catResult = await shopApi.listCategories(slug.value)
    categories.value = catResult.categories || []

    // Load Initial Products
    await loadProducts()
    
  } catch (e) {
    console.error('Failed to load shop', e)
    error.value = true
    loading.value = false
  }
}

// Watch for URL changes to reload products
watch(() => route.query, () => {
    if (shop.value) {
        loadProducts()
    }
})

onMounted(() => {
  loadShop()
})
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
    display: none;
}
.scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>
