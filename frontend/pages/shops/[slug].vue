<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto py-8 px-4">
      <div class="mb-6">
        <NuxtLink to="/shops" class="inline-flex items-center text-gray-500 hover:text-indigo-600 transition-colors">
          <span class="mr-2">←</span> Retour au Marketplace
        </NuxtLink>
      </div>
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
        <div class="relative h-48 md:h-64 rounded-2xl overflow-hidden mb-8">
          <img v-if="shop.banner_url" :src="shop.banner_url" class="w-full h-full object-cover" alt="">
          <div v-else class="w-full h-full bg-gradient-to-br from-indigo-500 to-purple-600"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"></div>
          
          <!-- Shop Info -->
          <div class="absolute bottom-0 left-0 right-0 p-6 flex items-end gap-4">
            <div class="w-20 h-20 rounded-xl bg-white dark:bg-slate-800 border-4 border-white shadow-xl overflow-hidden flex-shrink-0">
              <img v-if="shop.logo_url" :src="shop.logo_url" class="w-full h-full object-cover" alt="">
              <div v-else class="w-full h-full flex items-center justify-center text-3xl bg-gradient-to-br from-indigo-500 to-purple-500 text-white">
                {{ shop.name.charAt(0) }}
              </div>
            </div>
            <div class="flex-1 text-white">
              <h1 class="text-2xl md:text-3xl font-bold mb-1">{{ shop.name }}</h1>
              <p v-if="shop.description" class="text-white/80 text-sm line-clamp-2">{{ shop.description }}</p>
            </div>
            
            <!-- QR Code -->
            <button v-if="shop.qr_code" @click="showQR = true" class="p-3 bg-white/20 backdrop-blur rounded-xl hover:bg-white/30 transition-colors">
              <span class="text-2xl">📱</span>
            </button>
          </div>
        </div>

        <!-- Categories Filter -->
        <div v-if="categories.length > 0" class="flex gap-2 overflow-x-auto pb-4 mb-6">
          <button 
            @click="selectedCategory = ''"
            :class="[
              'px-4 py-2 rounded-full whitespace-nowrap font-medium transition-colors',
              selectedCategory === '' 
                ? 'bg-indigo-600 text-white' 
                : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700'
            ]"
          >
            Tous
          </button>
          <button 
            v-for="cat in categories" 
            :key="cat.id"
            @click="selectedCategory = cat.slug"
            :class="[
              'px-4 py-2 rounded-full whitespace-nowrap font-medium transition-colors',
              selectedCategory === cat.slug 
                ? 'bg-indigo-600 text-white' 
                : 'bg-white dark:bg-slate-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700'
            ]"
          >
            {{ cat.name }} ({{ cat.product_count }})
          </button>
        </div>

        <!-- Products Grid -->
        <div v-if="products.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 md:gap-6">
          <div 
            v-for="product in products" 
            :key="product.id"
            @click="openProduct(product)"
            class="group bg-white dark:bg-slate-900 rounded-xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-lg hover:border-indigo-500/30 transition-all cursor-pointer"
          >
            <!-- Image -->
            <div class="aspect-square bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
              <img v-if="product.images?.length" :src="product.images[0]" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" alt="">
              <div v-else class="w-full h-full flex items-center justify-center text-4xl">📦</div>
              
              <!-- Featured Badge -->
              <div v-if="product.is_featured" class="absolute top-2 left-2 px-2 py-1 bg-amber-500 text-white text-xs font-bold rounded-lg">
                ⭐ Featured
              </div>
              
              <!-- Stock Badge -->
              <div v-if="product.stock === 0" class="absolute top-2 right-2 px-2 py-1 bg-red-500 text-white text-xs font-bold rounded-lg">
                Épuisé
              </div>
            </div>
            
            <!-- Info -->
            <div class="p-3">
              <h3 class="font-semibold text-gray-900 dark:text-white text-sm mb-1 line-clamp-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
                {{ product.name }}
              </h3>
              <div class="flex items-center justify-between">
                <span class="text-lg font-bold text-indigo-600 dark:text-indigo-400">
                  {{ formatPrice(product.price, product.currency) }}
                </span>
                <span v-if="product.compare_at_price" class="text-sm text-gray-400 line-through">
                  {{ formatPrice(product.compare_at_price, product.currency) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty Products -->
        <div v-else class="text-center py-12 bg-white dark:bg-slate-800 rounded-2xl">
          <div class="text-5xl mb-4">📦</div>
          <p class="text-gray-500">Aucun produit dans cette boutique</p>
        </div>

        <!-- Cart FAB -->
        <div v-if="cartStore.itemCount > 0" class="fixed bottom-6 right-6 z-50">
          <NuxtLink 
            to="/cart" 
            class="flex items-center gap-3 px-6 py-3 bg-indigo-600 text-white rounded-full shadow-xl hover:bg-indigo-700 transition-colors"
          >
            <span class="text-xl">🛒</span>
            <span class="font-bold">{{ cartStore.itemCount }} articles</span>
            <span class="px-2 py-1 bg-white/20 rounded-lg">{{ formatPrice(cartStore.subtotal, shop.currency) }}</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Product Modal -->
      <Teleport to="body">
        <div v-if="selectedProduct" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="selectedProduct = null">
          <div class="bg-white dark:bg-slate-900 rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <!-- Images -->
            <div class="aspect-video bg-gray-100 dark:bg-slate-800 relative">
              <img v-if="selectedProduct.images?.length" :src="selectedProduct.images[selectedImageIndex]" class="w-full h-full object-contain" alt="">
              <button @click="selectedProduct = null" class="absolute top-4 right-4 p-2 bg-black/30 text-white rounded-full">✕</button>
              
              <!-- Image Navigation -->
              <div v-if="selectedProduct.images?.length > 1" class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                <button 
                  v-for="(img, i) in selectedProduct.images" 
                  :key="i"
                  @click="selectedImageIndex = i"
                  :class="['w-2 h-2 rounded-full', i === selectedImageIndex ? 'bg-white' : 'bg-white/50']"
                ></button>
              </div>
            </div>
            
            <div class="p-6">
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ selectedProduct.name }}</h2>
              <p class="text-gray-500 dark:text-gray-400 mb-4">{{ selectedProduct.description }}</p>
              
              <div class="flex items-center gap-4 mb-6">
                <span class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">
                  {{ formatPrice(selectedProduct.price, selectedProduct.currency) }}
                </span>
                <span v-if="selectedProduct.compare_at_price" class="text-xl text-gray-400 line-through">
                  {{ formatPrice(selectedProduct.compare_at_price, selectedProduct.currency) }}
                </span>
              </div>
              
              <!-- Quantity -->
              <div class="flex items-center gap-4 mb-6">
                <span class="text-gray-700 dark:text-gray-300">Quantité:</span>
                <div class="flex items-center border border-gray-200 dark:border-gray-700 rounded-lg">
                  <button @click="quantity > 1 && quantity--" class="px-4 py-2 hover:bg-gray-100 dark:hover:bg-slate-800">-</button>
                  <span class="px-4 py-2 font-bold">{{ quantity }}</span>
                  <button @click="quantity++" class="px-4 py-2 hover:bg-gray-100 dark:hover:bg-slate-800">+</button>
                </div>
              </div>
              
              <!-- Add to Cart -->
              <button 
                @click="addToCart"
                :disabled="selectedProduct.stock === 0"
                class="w-full btn-premium py-4 text-lg disabled:opacity-50"
              >
                {{ selectedProduct.stock === 0 ? 'Épuisé' : 'Ajouter au panier' }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- QR Modal -->
      <Teleport to="body">
        <div v-if="showQR && shop?.qr_code" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="showQR = false">
          <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 text-center">
            <img :src="shop.qr_code" alt="QR Code" class="w-64 h-64 mx-auto mb-4">
            <p class="text-gray-500">Scannez pour partager cette boutique</p>
            <button @click="showQR = false" class="mt-4 px-6 py-2 bg-gray-100 dark:bg-slate-800 rounded-lg">Fermer</button>
          </div>
        </div>
      </Teleport>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useShopApi, type Shop, type Product, type Category } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'

const route = useRoute()
const shopApi = useShopApi()
const cartStore = useCartStore()

const slug = computed(() => route.params.slug as string)

const loading = ref(true)
const error = ref(false)
const shop = ref<Shop | null>(null)
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const selectedCategory = ref('')
const selectedProduct = ref<Product | null>(null)
const selectedImageIndex = ref(0)
const quantity = ref(1)
const showQR = ref(false)

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const openProduct = (product: Product) => {
  selectedProduct.value = product
  selectedImageIndex.value = 0
  quantity.value = 1
}

const addToCart = () => {
  if (selectedProduct.value && shop.value) {
    cartStore.setShopInfo(shop.value.id, shop.value.slug, shop.value.name, shop.value.currency)
    cartStore.addItem(selectedProduct.value, quantity.value)
    selectedProduct.value = null
  }
}

const loadShop = async () => {
  loading.value = true
  error.value = false
  try {
    shop.value = await shopApi.getShop(slug.value)
    const [prodResult, catResult] = await Promise.all([
      shopApi.listProducts(slug.value, 1, 100),
      shopApi.listCategories(slug.value)
    ])
    products.value = prodResult.products || []
    categories.value = catResult.categories || []
  } catch (e) {
    console.error('Failed to load shop', e)
    error.value = true
  } finally {
    loading.value = false
  }
}

watch(selectedCategory, async () => {
  if (!shop.value) return
  loading.value = true
  try {
    const result = await shopApi.listProducts(slug.value, 1, 100, { category: selectedCategory.value })
    products.value = result.products || []
  } finally {
    loading.value = false
  }
})

onMounted(() => {
  cartStore.loadFromStorage()
  loadShop()
})
</script>
