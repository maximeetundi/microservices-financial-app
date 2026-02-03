<template>
  <NuxtLayout name="dashboard">
    <div class="w-full py-6 px-6">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">🛒 Mon Panier</h1>

      <!-- Empty Cart -->
      <div v-if="cartStore.isEmpty" class="text-center py-16 bg-white dark:bg-slate-800 rounded-2xl">
        <div class="text-6xl mb-4">🛒</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Votre panier est vide</h3>
        <p class="text-gray-500 mb-6">Découvrez nos boutiques pour trouver des produits</p>
        <NuxtLink to="/shops" class="btn-premium px-6 py-3">
          Explorer les boutiques
        </NuxtLink>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Cart Items -->
        <div class="lg:col-span-2 space-y-4">
          <div class="bg-white dark:bg-slate-800 rounded-xl p-4 mb-4">
            <p class="text-sm text-gray-500">
              Articles de <span class="font-bold text-indigo-600">{{ cartStore.shopName }}</span>
            </p>
          </div>

          <div 
            v-for="item in cartStore.items" 
            :key="item.product_id"
            class="bg-white dark:bg-slate-900 rounded-xl p-4 flex gap-4 border border-gray-100 dark:border-gray-800"
          >
            <!-- Image -->
            <div class="w-20 h-20 rounded-lg bg-gray-100 dark:bg-slate-800 overflow-hidden flex-shrink-0">
              <img v-if="item.product?.images?.[0]" :src="item.product.images[0]" class="w-full h-full object-cover" alt="">
              <div v-else class="w-full h-full flex items-center justify-center text-2xl">📦</div>
            </div>
            
            <!-- Info -->
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-900 dark:text-white truncate">{{ item.product?.name }}</h3>
              <p class="text-indigo-600 dark:text-indigo-400 font-bold">
                {{ formatPrice(item.product?.price || 0, cartStore.shopCurrency || 'XOF') }}
              </p>
              
              <!-- Quantity -->
              <div class="flex items-center gap-2 mt-2">
                <button 
                  @click="updateQuantity(item.product_id, item.quantity - 1)"
                  class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-slate-700 hover:bg-gray-200 dark:hover:bg-slate-600"
                >
                  -
                </button>
                <span class="w-8 text-center font-bold">{{ item.quantity }}</span>
                <button 
                  @click="updateQuantity(item.product_id, item.quantity + 1)"
                  class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-slate-700 hover:bg-gray-200 dark:hover:bg-slate-600"
                >
                  +
                </button>
              </div>
            </div>
            
            <!-- Total & Remove -->
            <div class="text-right">
              <p class="font-bold text-gray-900 dark:text-white">
                {{ formatPrice((item.product?.price || 0) * item.quantity, cartStore.shopCurrency || 'XOF') }}
              </p>
              <button @click="removeItem(item.product_id)" class="text-red-500 text-sm mt-2 hover:underline">
                Supprimer
              </button>
            </div>
          </div>
        </div>

        <!-- Order Summary -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-slate-900 rounded-xl p-6 border border-gray-100 dark:border-gray-800 sticky top-4">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Résumé</h3>
            
            <div class="space-y-3 mb-6">
              <div class="flex justify-between text-gray-600 dark:text-gray-400">
                <span>Sous-total</span>
                <span>{{ formatPrice(cartStore.subtotal, cartStore.shopCurrency || 'XOF') }}</span>
              </div>
              <div class="flex justify-between text-gray-600 dark:text-gray-400">
                <span>Livraison</span>
                <span>À calculer</span>
              </div>
              <div class="border-t border-gray-100 dark:border-gray-800 pt-3 flex justify-between text-lg font-bold text-gray-900 dark:text-white">
                <span>Total</span>
                <span>{{ formatPrice(cartStore.subtotal, cartStore.shopCurrency || 'XOF') }}</span>
              </div>
            </div>

            <!-- Wallet Selection -->
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Payer avec
              </label>
              <select v-model="selectedWalletId" class="w-full input-premium">
                <option value="">Sélectionner un portefeuille</option>
                <option v-for="w in wallets" :key="w.id" :value="w.id">
                  {{ w.currency }} - {{ formatPrice(w.balance, w.currency) }}
                </option>
              </select>
            </div>

            <!-- Delivery Type -->
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Mode de livraison
              </label>
              <div class="flex gap-2">
                <button 
                  @click="deliveryType = 'pickup'"
                  :class="['flex-1 py-2 rounded-lg border', deliveryType === 'pickup' ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-600' : 'border-gray-200 dark:border-gray-700']"
                >
                  🏃 Retrait
                </button>
                <button 
                  @click="deliveryType = 'delivery'"
                  :class="['flex-1 py-2 rounded-lg border', deliveryType === 'delivery' ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-600' : 'border-gray-200 dark:border-gray-700']"
                >
                  🚚 Livraison
                </button>
              </div>
            </div>

            <button 
              @click="checkout"
              :disabled="!selectedWalletId || processing"
              class="w-full btn-premium py-4 text-lg disabled:opacity-50"
            >
              {{ processing ? 'Traitement...' : 'Payer maintenant' }}
            </button>

            <button @click="cartStore.clearCart()" class="w-full mt-3 text-red-500 text-sm hover:underline">
              Vider le panier
            </button>
          </div>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useCartStore } from '~/stores/cart'
import { useShopApi } from '~/composables/useShopApi'
import { useApi } from '~/composables/useApi'

definePageMeta({ middleware: 'auth' })

const cartStore = useCartStore()
const shopApi = useShopApi()
const { walletApi } = useApi()

const wallets = ref<any[]>([])
const selectedWalletId = ref('')
const deliveryType = ref('pickup')
const processing = ref(false)

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount)
}

const updateQuantity = (productId: string, quantity: number) => {
  cartStore.updateQuantity(productId, quantity)
}

const removeItem = (productId: string) => {
  cartStore.removeItem(productId)
}

const loadWallets = async () => {
  try {
    const result = await walletApi.getWallets()
    wallets.value = result.wallets || []
  } catch (e) {
    console.error('Failed to load wallets', e)
  }
}

const checkout = async () => {
  if (!cartStore.shopId || !selectedWalletId.value) return
  
  processing.value = true
  try {
    const order = await shopApi.createOrder({
      shop_id: cartStore.shopId,
      items: cartStore.items.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
        custom_values: item.custom_values
      })),
      wallet_id: selectedWalletId.value,
      delivery_type: deliveryType.value
    })
    
    cartStore.clearCart()
    navigateTo(`/orders/${order.id}`)
  } catch (e: any) {
    alert(e.message || 'Échec du paiement')
  } finally {
    processing.value = false
  }
}

onMounted(() => {
  cartStore.loadFromStorage()
  loadWallets()
})
</script>
