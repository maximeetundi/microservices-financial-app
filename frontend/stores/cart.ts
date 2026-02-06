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

// Helper to get user-specific cart key
const getCartKey = (): string => {
    if (typeof window === 'undefined') return 'cart'
    try {
        const authData = localStorage.getItem('auth')
        if (authData) {
            const auth = JSON.parse(authData)
            if (auth.user?.id) {
                return `cart_${auth.user.id}`
            }
        }
    } catch (e) {
        console.error('Failed to get user ID for cart key', e)
    }
    return 'cart_guest'
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
                const key = getCartKey()
                localStorage.setItem(key, JSON.stringify({
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
                const key = getCartKey()
                const saved = localStorage.getItem(key)
                if (saved) {
                    try {
                        const data = JSON.parse(saved)
                        const rawItems = Array.isArray(data.items) ? data.items : []

                        this.items = rawItems
                            .map((raw: any) => {
                                const product_id = raw?.product_id || raw?.id
                                const quantity = typeof raw?.quantity === 'number' ? raw.quantity : 0
                                if (!product_id || quantity <= 0) return null

                                const product = raw?.product || (raw?.name ? {
                                    id: product_id,
                                    shop_id: data.shopId || raw?.shop_id || raw?.shopId || 'unknown',
                                    name: raw.name,
                                    price: raw.price || 0,
                                    images: raw.image ? [raw.image] : [],
                                } : undefined)

                                return {
                                    product_id,
                                    quantity,
                                    custom_values: raw?.custom_values,
                                    product,
                                } as CartItem
                            })
                            .filter(Boolean) as CartItem[]

                        this.shopId = data.shopId || null
                        this.shopSlug = data.shopSlug || null
                        this.shopName = data.shopName || null
                        this.shopCurrency = data.shopCurrency || null

                        if (this.items.length === 0) {
                            this.shopId = null
                            this.shopSlug = null
                            this.shopName = null
                            this.shopCurrency = null
                        }
                    } catch (e) {
                        console.error('Failed to load cart', e)
                        this.clearCart()
                    }
                }
            }
        },
    },
})

