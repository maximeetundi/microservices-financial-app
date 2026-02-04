<template>
  <div class="max-w-7xl mx-auto pb-16">
    <!-- Breadcrumb -->
    <nav class="flex items-center text-sm text-gray-500 mb-8 overflow-x-auto whitespace-nowrap">
        <NuxtLink :to="`/shops/${shopSlug}`" class="hover:text-indigo-600 transition-colors">Boutique</NuxtLink>
        <span class="mx-2">/</span>
        <span class="text-gray-900 dark:text-gray-300 font-medium truncate max-w-[200px]">{{ product?.name || 'Chargement...' }}</span>
    </nav>
    
    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-12">
        <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-[500px] rounded-2xl"></div>
        <div class="space-y-6">
            <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-10 w-3/4 rounded-xl"></div>
            <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-6 w-1/4 rounded-xl"></div>
            <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-32 w-full rounded-xl"></div>
            <div class="animate-pulse bg-gray-100 dark:bg-slate-800 h-16 w-full rounded-xl"></div>
        </div>
    </div>

    <!-- Error -->
    <div v-else-if="error || !product" class="text-center py-24">
         <div class="text-6xl mb-4">😕</div>
         <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">Produit introuvable</h1>
         <NuxtLink :to="`/shops/${shopSlug}`" class="btn-primary">Retour à la boutique</NuxtLink>
    </div>

    <div v-else>
        <!-- Product Main Section -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 mb-16">
            <!-- Image Gallery -->
            <div class="space-y-4">
                <div class="aspect-square bg-gray-50 dark:bg-slate-900 rounded-3xl overflow-hidden border border-gray-100 dark:border-gray-800 relative group">
                    <img v-if="product.images?.length" :src="product.images[selectedImageIndex]" class="w-full h-full object-contain p-4 group-hover:scale-105 transition-transform duration-500" alt="">
                    <div v-else class="w-full h-full flex items-center justify-center text-6xl text-gray-300">📦</div>
                    
                    <!-- Badges -->
                    <div class="absolute top-4 left-4 flex flex-col gap-2">
                         <div v-if="product.is_featured" class="px-3 py-1.5 bg-amber-500 text-white text-xs font-bold rounded-lg shadow-sm">
                            ⭐ Featured
                         </div>
                         <div v-if="product.stock === 0" class="px-3 py-1.5 bg-red-500 text-white text-xs font-bold rounded-lg shadow-sm">
                            Épuisé
                         </div>
                    </div>
                </div>
                
                <!-- Thumbnails -->
                <div v-if="product.images?.length > 1" class="flex gap-4 overflow-x-auto pb-2 scrollbar-thin">
                    <button 
                        v-for="(img, i) in product.images" 
                        :key="i"
                        @click="selectedImageIndex = i"
                        :class="[
                            'w-20 h-20 rounded-xl border-2 flex-shrink-0 overflow-hidden transition-all',
                            selectedImageIndex === i ? 'border-indigo-600 shadow-md transform scale-105' : 'border-transparent hover:border-gray-300 dark:hover:border-gray-700 bg-gray-50 dark:bg-slate-900'
                        ]"
                    >
                        <img :src="img" class="w-full h-full object-cover" alt="">
                    </button>
                </div>
            </div>

            <!-- Info -->
            <div>
                <div class="mb-2 flex items-center gap-2">
                     <span v-if="product.status === 'Active'" class="px-2.5 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 text-xs font-bold rounded-full uppercase tracking-wide">En Stock</span>
                     <!-- Category? -->
                </div>
                
                <h1 class="text-3xl lg:text-4xl font-bold text-gray-900 dark:text-white mb-4 leading-tight">{{ product.name }}</h1>
                
                <div class="flex items-end gap-4 mb-6">
                    <span class="text-4xl font-bold text-indigo-600 dark:text-indigo-400">
                        {{ formatPrice(product.price, product.currency) }}
                    </span>
                    <span v-if="product.compare_at_price" class="text-xl text-gray-400 line-through mb-1">
                        {{ formatPrice(product.compare_at_price, product.currency) }}
                    </span>
                    <span v-if="product.compare_at_price" class="text-sm font-bold text-green-600 dark:text-green-400 mb-2 px-2 py-1 bg-green-50 dark:bg-green-900/20 rounded-lg">
                        -{{ Math.round(((product.compare_at_price - product.price) / product.compare_at_price) * 100) }}%
                    </span>
                </div>

                <div class="prose dark:prose-invert text-gray-600 dark:text-gray-300 mb-8 max-w-none">
                    <p>{{ product.description }}</p>
                </div>

                <!-- Quantity -->
                <div class="flex items-center gap-6 mb-8 p-4 bg-gray-50 dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-gray-800">
                    <span class="text-sm font-bold text-gray-500 uppercase tracking-wider">Quantité</span>
                    <div class="flex items-center gap-4">
                        <button 
                            @click="quantity > 1 && quantity--" 
                            class="w-10 h-10 rounded-xl bg-white dark:bg-slate-800 shadow-sm border border-gray-200 dark:border-gray-700 hover:border-indigo-500 text-lg font-bold transition-all"
                        >-</button>
                        <span class="w-8 text-center text-xl font-bold">{{ quantity }}</span>
                        <button 
                            @click="quantity++" 
                            class="w-10 h-10 rounded-xl bg-white dark:bg-slate-800 shadow-sm border border-gray-200 dark:border-gray-700 hover:border-indigo-500 text-lg font-bold transition-all"
                        >+</button>
                    </div>
                </div>

                <!-- Actions -->
                <div class="flex flex-col sm:flex-row gap-4">
                    <button 
                        @click="addToCart"
                        :disabled="product.stock === 0"
                        class="flex-1 py-4 px-8 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-lg shadow-indigo-500/30 transform active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3"
                    >
                        <span class="text-2xl">🛒</span>
                        <span>{{ product.stock === 0 ? 'Rupture de stock' : 'Ajouter au panier' }}</span>
                    </button>
                    
                    <button 
                        @click="toggleFavorite"
                        :class="[
                            'p-4 rounded-xl border-2 transition-all flex items-center justify-center',
                            isFavorite 
                                ? 'border-red-500 text-red-500 bg-red-50 dark:bg-red-900/20' 
                                : 'border-gray-200 dark:border-gray-700 hover:border-red-300 text-gray-400 hover:text-red-500'
                        ]"
                    >
                        <span class="text-2xl">{{ isFavorite ? '❤️' : '🤍' }}</span>
                    </button>
                </div>

                <!-- Trust Badges -->
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-8 pt-8 border-t border-gray-100 dark:border-gray-800">
                     <div class="flex flex-col items-center text-center gap-1 p-2 rounded-lg bg-gray-50 dark:bg-slate-900/50">
                         <span class="text-2xl">🔒</span>
                         <span class="text-[10px] font-bold text-gray-500 uppercase">Paiement Sécurisé</span>
                     </div>
                     <div class="flex flex-col items-center text-center gap-1 p-2 rounded-lg bg-gray-50 dark:bg-slate-900/50">
                         <span class="text-2xl">⚡</span>
                         <span class="text-[10px] font-bold text-gray-500 uppercase">Livraison Rapide</span>
                     </div>
                     <div class="flex flex-col items-center text-center gap-1 p-2 rounded-lg bg-gray-50 dark:bg-slate-900/50">
                         <span class="text-2xl">🛡️</span>
                         <span class="text-[10px] font-bold text-gray-500 uppercase">Garantie Vendeur</span>
                     </div>
                     <div class="flex flex-col items-center text-center gap-1 p-2 rounded-lg bg-gray-50 dark:bg-slate-900/50">
                         <span class="text-2xl">💬</span>
                         <span class="text-[10px] font-bold text-gray-500 uppercase">Support 24/7</span>
                     </div>
                </div>
            </div>
        </div>

        <!-- Reviews & Details Tabs -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 mb-16 shadow-sm border border-gray-100 dark:border-gray-800">
            <div class="flex border-b border-gray-200 dark:border-gray-700 mb-8">
                <button 
                  @click="activeTab = 'description'"
                  :class="['px-6 py-4 font-bold text-lg border-b-2 transition-colors', activeTab === 'description' ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300']"
                >
                  Description
                </button>
                <button 
                  @click="activeTab = 'reviews'"
                  :class="['px-6 py-4 font-bold text-lg border-b-2 transition-colors', activeTab === 'reviews' ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300']"
                >
                  Avis ({{ reviews.length }})
                </button>
            </div>

            <div v-if="activeTab === 'description'">
                <h3 class="font-bold text-xl mb-4">À propos de ce produit</h3>
                <p class="text-gray-600 dark:text-gray-300 leading-relaxed max-w-3xl">
                    {{ product.description }}
                    <br><br>
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
                </p>
            </div>

            <div v-if="activeTab === 'reviews'" class="space-y-8">
                <div v-for="review in reviews" :key="review.id" class="flex gap-4">
                    <div class="w-12 h-12 rounded-full bg-indigo-100 flex-shrink-0 flex items-center justify-center font-bold text-indigo-600">
                        {{ review.user_name ? review.user_name.charAt(0) : 'U' }}
                    </div>
                    <div>
                        <div class="flex items-center gap-2 mb-1">
                            <h4 class="font-bold dark:text-white">{{ review.user_name }}</h4>
                            <span class="text-xs text-gray-500">{{ new Date(review.created_at).toLocaleDateString() }}</span>
                        </div>
                        <div class="flex text-amber-500 text-sm mb-2">
                             <span v-for="i in 5" :key="i">{{ i <= review.rating ? '★' : '☆' }}</span>
                        </div>
                        <p class="text-gray-600 dark:text-gray-300">{{ review.comment }}</p>
                    </div>
                </div>
                
                <!-- Add Review Form -->
                <div class="pt-8 mt-8 border-t border-gray-100 dark:border-gray-800">
                    <h4 class="font-bold mb-4">Laisser un avis</h4>
                    <div class="flex gap-2 mb-4">
                        <button 
                            v-for="i in 5" 
                            :key="i" 
                            type="button"
                            @click="newReviewRating = i"
                            class="text-2xl transition-colors"
                            :class="i <= newReviewRating ? 'text-amber-500' : 'text-gray-300 hover:text-amber-400'"
                        >★</button>
                    </div>
                    <textarea 
                        v-model="newReviewComment"
                        placeholder="Votre avis sur ce produit..." 
                        class="w-full p-4 rounded-xl bg-gray-50 dark:bg-slate-800 border-none focus:ring-2 focus:ring-indigo-500 mb-4 h-32 text-gray-900 dark:text-gray-200"
                    ></textarea>
                    <button 
                        @click="submitReview" 
                        :disabled="!newReviewComment || submittingReview"
                        class="px-6 py-2 bg-indigo-600 text-white rounded-lg font-bold disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {{ submittingReview ? 'Envoi...' : 'Publier' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Related Products -->
        <div v-if="relatedProducts.length > 0">
            <h2 class="text-2xl font-bold mb-8 flex items-center gap-2">
                <span>📦</span> Produits similaires
            </h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
                 <!-- Reuse product card logic or component -->
                 <div 
                    v-for="prod in relatedProducts" 
                    :key="prod.id"
                    @click="navigateTo(`/shops/${shopSlug}/product/${prod.slug}`)"
                    class="group bg-white dark:bg-slate-900 rounded-2xl overflow-hidden border border-gray-100 dark:border-gray-800 hover:shadow-lg cursor-pointer transition-all"
                  >
                     <div class="aspect-square bg-gray-100 dark:bg-slate-800 relative overflow-hidden">
                        <img v-if="prod.images?.length" :src="prod.images[0]" class="w-full h-full object-cover group-hover:scale-105 transition-transform" alt="">
                     </div>
                     <div class="p-4">
                        <h3 class="font-bold text-sm mb-1 line-clamp-1 break-all">{{ prod.name }}</h3>
                        <span class="text-indigo-600 dark:text-indigo-400 font-bold">
                            {{ formatPrice(prod.price, prod.currency) }}
                        </span>
                     </div>
                 </div>
            </div>
        </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { useShopApi, type Shop, type Product, type Review } from '~/composables/useShopApi'
import { useCartStore } from '~/stores/cart'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'shop-customer'
})

const route = useRoute()
const shopApi = useShopApi()
const cartStore = useCartStore()
const authStore = useAuthStore()

const shopSlug = computed(() => route.params.slug as string)
const productSlug = computed(() => route.params.product_slug as string)

const loading = ref(true)
const error = ref(false)
const product = ref<Product | null>(null)
const selectedImageIndex = ref(0)
const quantity = ref(1)
const isFavorite = ref(false)
const activeTab = ref('description')
const relatedProducts = ref<Product[]>([])

// Reviews State
const reviews = ref<Review[]>([])
const reviewsLoading = ref(false)
const newReviewRating = ref(5)
const newReviewComment = ref('')
const submittingReview = ref(false)

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const toggleFavorite = () => {
    isFavorite.value = !isFavorite.value
    // Logic to save to local storage or API would go here
}

const addToCart = async () => {
    if (product.value) {
        try {
             const shop = await shopApi.getShop(shopSlug.value)
             cartStore.setShopInfo(shop.id, shop.slug, shop.name, shop.currency)
             cartStore.addItem(product.value, quantity.value)
             alert('Produit ajouté au panier !')
        } catch (e) {
            console.error(e)
        }
    }
}

const loadReviews = async () => {
    if (!product.value) return
    reviewsLoading.value = true
    try {
        const result = await shopApi.listReviews(product.value.id)
        reviews.value = result.reviews || []
    } catch (e) {
        console.error('Failed to load reviews', e)
    } finally {
        reviewsLoading.value = false
    }
}

const submitReview = async () => {
    if (!product.value) return
    if (!authStore.isLoggedIn) {
        alert("Veuillez vous connecter pour laisser un avis.")
        return
    }
    
    submittingReview.value = true
    try {
        await shopApi.createReview(product.value.id, {
            rating: newReviewRating.value,
            comment: newReviewComment.value
        })
        newReviewComment.value = ''
        newReviewRating.value = 5
        await loadReviews() // Reload reviews
        alert("Merci pour votre avis !")
    } catch (e) {
        console.error('Failed to submit review', e)
        alert("Erreur lors de l'envoi de l'avis.")
    } finally {
        submittingReview.value = false
    }
}

const loadProduct = async () => {
  loading.value = true
  error.value = false
  try {
    product.value = await shopApi.getProduct(shopSlug.value, productSlug.value)
    
    // Load related
    if (product.value?.category_id) {
         const related = await shopApi.listProducts(shopSlug.value, 1, 4, { category: product.value.category_id })
         relatedProducts.value = related.products.filter(p => p.id !== product.value?.id) || []
    }

    // Load reviews
    if (product.value) {
        loadReviews()
    }

  } catch (e) {
    console.error('Failed to load product', e)
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  cartStore.loadFromStorage()
  loadProduct()
})
</script>

<style scoped>
.scrollbar-thin::-webkit-scrollbar {
    height: 4px;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
    @apply bg-gray-300 dark:bg-gray-700 rounded-full;
}
</style>
