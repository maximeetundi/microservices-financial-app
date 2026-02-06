<template>
  <!-- Cart Floating Button -->
  <div v-if="cartItemCount > 0" class="fixed bottom-6 right-6 z-50">
    <button
      @click="showCartModal = true"
      class="relative w-16 h-16 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-full shadow-lg shadow-orange-500/30 flex items-center justify-center hover:scale-105 active:scale-95 transition-all duration-300 group"
    >
      <!-- Cart Icon -->
      <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
      </svg>
      
      <!-- Badge with item count -->
      <span class="absolute -top-1 -right-1 bg-red-500 text-white text-xs font-bold rounded-full w-6 h-6 flex items-center justify-center animate-pulse">
        {{ cartItemCount }}
      </span>
      
      <!-- Hover effect -->
      <div class="absolute inset-0 rounded-full bg-white opacity-0 group-hover:opacity-20 transition-opacity"></div>
    </button>
  </div>

  <!-- Cart Modal -->
  <Teleport to="body">
    <div
      v-if="showCartModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <!-- Backdrop -->
      <div
        class="absolute inset-0 bg-black/50 backdrop-blur-sm"
        @click="showCartModal = false"
      ></div>
      
      <!-- Modal Content -->
      <div class="relative bg-white dark:bg-surface rounded-2xl shadow-2xl w-full max-w-md max-h-[80vh] flex flex-col animate-fade-in-up">
        <!-- Header -->
        <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-secondary-800">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">
            Votre Panier
          </h2>
          <button
            @click="showCartModal = false"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-secondary-800 transition-colors"
          >
            <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        
        <!-- Cart Items -->
        <div class="flex-1 overflow-y-auto p-6">
          <div v-if="cartItems.length === 0" class="text-center py-8">
            <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/>
            </svg>
            <p class="text-gray-500 dark:text-gray-400">Votre panier est vide</p>
          </div>
          
          <div v-else class="space-y-4">
            <div
              v-for="item in cartItems"
              :key="item.id"
              class="flex gap-4 p-4 bg-gray-50 dark:bg-secondary-800/50 rounded-xl"
            >
              <!-- Product Image -->
              <img
                :src="item.image || '/placeholder-product.png'"
                :alt="item.name"
                class="w-16 h-16 rounded-lg object-cover bg-gray-200 dark:bg-secondary-700"
              />
              
              <!-- Product Info -->
              <div class="flex-1 min-w-0">
                <h3 class="font-semibold text-gray-900 dark:text-white truncate">
                  {{ item.name }}
                </h3>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ item.shop_name }}
                </p>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-lg font-bold text-orange-600">
                    {{ formatMoney(item.price * item.quantity) }}
                  </span>
                  <div class="flex items-center gap-2">
                    <!-- Quantity controls -->
                    <button
                      @click="updateQuantity(item.id, item.quantity - 1)"
                      class="w-8 h-8 rounded-lg bg-gray-200 dark:bg-secondary-700 hover:bg-gray-300 dark:hover:bg-secondary-600 transition-colors flex items-center justify-center"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
                      </svg>
                    </button>
                    <span class="w-8 text-center font-medium">{{ item.quantity }}</span>
                    <button
                      @click="updateQuantity(item.id, item.quantity + 1)"
                      class="w-8 h-8 rounded-lg bg-gray-200 dark:bg-secondary-700 hover:bg-gray-300 dark:hover:bg-secondary-600 transition-colors flex items-center justify-center"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                      </svg>
                    </button>
                    <!-- Remove button -->
                    <button
                      @click="removeFromCart(item.id)"
                      class="w-8 h-8 rounded-lg bg-red-100 dark:bg-red-900/30 hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors flex items-center justify-center text-red-600"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Footer -->
        <div v-if="cartItems.length > 0" class="border-t border-gray-200 dark:border-secondary-800 p-6 space-y-4">
          <!-- Total -->
          <div class="flex items-center justify-between">
            <span class="text-lg font-semibold text-gray-900 dark:text-white">Total</span>
            <span class="text-2xl font-bold text-orange-600">{{ formatMoney(cartTotal) }}</span>
          </div>
          
          <!-- Action Buttons -->
          <div class="flex gap-3">
            <button
              @click="showCartModal = false"
              class="flex-1 px-4 py-3 bg-gray-100 dark:bg-secondary-800 text-gray-700 dark:text-gray-300 rounded-xl font-medium hover:bg-gray-200 dark:hover:bg-secondary-700 transition-colors"
            >
              Continuer
            </button>
            <button
              @click="goToCheckout"
              class="flex-1 px-4 py-3 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-xl font-medium hover:from-orange-600 hover:to-orange-700 transition-all duration-200 shadow-lg shadow-orange-500/30"
            >
              Payer
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '~/stores/cart'

const router = useRouter()
const cartStore = useCartStore()
const showCartModal = ref(false)

// Computed properties
const cartItems = computed(() => cartStore.items || [])
const cartItemCount = computed(() => cartStore.totalItems || 0)
const cartTotal = computed(() => cartStore.total || 0)

// Methods
const formatMoney = (amount) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'XOF'
  }).format(amount)
}

const updateQuantity = (itemId, newQuantity) => {
  if (newQuantity <= 0) {
    removeFromCart(itemId)
  } else {
    cartStore.updateQuantity(itemId, newQuantity)
  }
}

const removeFromCart = (itemId) => {
  cartStore.removeItem(itemId)
}

const goToCheckout = () => {
  showCartModal.value = false
  // Rediriger vers la page panier (quand elle existera)
  router.push('/cart')
}

// Initialize cart on mount
onMounted(() => {
  cartStore.loadCart()
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
