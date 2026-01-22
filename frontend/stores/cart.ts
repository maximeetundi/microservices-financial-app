import { defineStore } from 'pinia'
import type { Product } from '~/composables/useShopApi'

interface CartItem {
    product_id: string
    quantity: number
    custom_values?: Record<string, string>
    product?: Product
}

interface CartState {
    items: CartItem[]
    shopId: string | null
    shopSlug: string | null
    shopName: string | null
    shopCurrency: string | null
}

export const useCartStore = defineStore('cart', {
    state: (): CartState => ({
        items: [],
        shopId: null,
        shopSlug: null,
        shopName: null,
        shopCurrency: null,
    }),

    getters: {
        itemCount: (state) => state.items.reduce((sum, item) => sum + item.quantity, 0),

        subtotal: (state) => state.items.reduce((sum, item) => {
            const price = item.product?.price || 0
            return sum + (price * item.quantity)
        }, 0),

        isEmpty: (state) => state.items.length === 0,
    },

    actions: {
        addItem(product: Product, quantity: number = 1, customValues?: Record<string, string>) {
            if (this.shopId && this.shopId !== product.shop_id) {
                this.clearCart()
            }

            if (!this.shopId) {
                this.shopId = product.shop_id
            }

            const existingIndex = this.items.findIndex(
                item => item.product_id === product.id && JSON.stringify(item.custom_values) === JSON.stringify(customValues)
            )

            if (existingIndex >= 0) {
                this.items[existingIndex].quantity += quantity
            } else {
                this.items.push({
                    product_id: product.id,
                    quantity,
                    custom_values: customValues,
                    product,
                })
            }

            this.saveToStorage()
        },

        updateQuantity(productId: string, quantity: number) {
            const item = this.items.find(i => i.product_id === productId)
            if (item) {
                if (quantity <= 0) {
                    this.removeItem(productId)
                } else {
                    item.quantity = quantity
                    this.saveToStorage()
                }
            }
        },

        removeItem(productId: string) {
            this.items = this.items.filter(i => i.product_id !== productId)
            if (this.items.length === 0) {
                this.shopId = null
                this.shopSlug = null
                this.shopName = null
                this.shopCurrency = null
            }
            this.saveToStorage()
        },

        clearCart() {
            this.items = []
            this.shopId = null
            this.shopSlug = null
            this.shopName = null
            this.shopCurrency = null
            this.saveToStorage()
        },

        setShopInfo(shopId: string, shopSlug: string, shopName: string, currency: string) {
            this.shopId = shopId
            this.shopSlug = shopSlug
            this.shopName = shopName
            this.shopCurrency = currency
        },

        saveToStorage() {
            if (process.client) {
                localStorage.setItem('cart', JSON.stringify({
                    items: this.items,
                    shopId: this.shopId,
                    shopSlug: this.shopSlug,
                    shopName: this.shopName,
                    shopCurrency: this.shopCurrency,
                }))
            }
        },

        loadFromStorage() {
            if (process.client) {
                const saved = localStorage.getItem('cart')
                if (saved) {
                    try {
                        const data = JSON.parse(saved)
                        this.items = data.items || []
                        this.shopId = data.shopId
                        this.shopSlug = data.shopSlug
                        this.shopName = data.shopName
                        this.shopCurrency = data.shopCurrency
                    } catch (e) {
                        console.error('Failed to load cart', e)
                    }
                }
            }
        },
    },
})
