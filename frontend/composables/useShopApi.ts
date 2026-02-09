import api from './useApi'

export interface Shop {
    id: string
    owner_id: string
    owner_type: string
    name: string
    slug: string
    description: string
    logo_url: string
    banner_url: string
    is_public: boolean
    wallet_id: string
    currency: string
    managers: ShopManager[]
    tags: string[]
    qr_code: string
    status: string
    settings: ShopSettings
    stats: ShopStats
    trust_badges?: ShopTrustBadge[]
    created_at: string
    updated_at: string
}

export interface ShopTrustBadge {
    key: string
    icon: string
    title: string
    subtitle: string
    enabled: boolean
    order: number
}

export interface ShopManager {
    user_id: string
    email: string
    first_name: string
    last_name: string
    role: string
    permissions: string[]
    status: string
}

export interface ShopSettings {
    allow_pickup: boolean
    allow_delivery: boolean
    delivery_fee: number
    min_order_amount: number
    max_order_amount: number
    auto_accept_orders: boolean
}

export interface ShopStats {
    total_products: number
    total_orders: number
    total_revenue: number
    average_rating: number
    total_reviews: number
}

export interface Product {
    id: string
    shop_id: string
    category_id: string
    name: string
    slug: string
    description: string
    short_desc: string
    price: number
    compare_at_price: number
    currency: string
    images: string[]
    is_digital?: boolean
    digital_file_url?: string
    license_text?: string
    stock: number
    is_customizable: boolean
    custom_fields: CustomField[]
    tags: string[]
    qr_code: string
    status: string
    is_featured: boolean
    sold_count: number
    view_count: number
    created_at: string
    updated_at: string
}

export interface CustomField {
    name: string
    type: string
    required: boolean
    options?: string[]
    placeholder?: string
}

export interface Category {
    id: string
    shop_id: string
    parent_id?: string
    name: string
    slug: string
    description: string
    image_url: string
    icon?: string
    qr_code: string
    order: number
    is_active: boolean
    product_count: number
    children?: Category[]
}

export interface Order {
    id: string
    order_number: string
    shop_id: string
    shop_name: string
    buyer_id: string
    buyer_name: string
    items: OrderItem[]
    sub_total: number
    delivery_fee: number
    total_amount: number
    currency: string
    converted_amount: number
    buyer_currency: string
    transaction_id: string
    payment_status: string
    order_status: string
    delivery_type: string
    shipping_address?: Address
    notes: string
    created_at: string
}

export interface OrderItem {
    product_id: string
    product_name: string
    product_image: string
    quantity: number
    unit_price: number
    total_price: number
    custom_values?: Record<string, string>
}

export interface Address {
    street: string
    city: string
    state: string
    country: string
    postal_code: string
}

export interface CartItem {
    product_id: string
    variant_id?: string
    quantity: number
    custom_values?: Record<string, string>
    product?: Product
}

export interface Review {
    id: string
    product_id: string
    user_id: string
    user_name: string
    rating: number
    comment: string
    created_at: string
}

export interface ReviewListResponse {
    reviews: Review[]
    total: number
    page: number
}

const baseUrl = '/shop-service/api/v1'

export const shopAPI = {
    // Shops
    listShops: (page = 1, pageSize = 20, search = '') => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (search) params.append('search', search)
        return api.get(`${baseUrl}/shops?${params}`)
    },
    getShop: (slug: string) => api.get(`${baseUrl}/shops/${slug}`),
    getMyShops: () => api.get(`${baseUrl}/my-shops`),
    createShop: (data: Partial<Shop>) => api.post(`${baseUrl}/shops`, data),
    updateShop: (id: string, data: Partial<Shop>) => api.put(`${baseUrl}/shops/${id}`, data),
    updateTrustBadges: (id: string, badges: ShopTrustBadge[]) => api.put(`${baseUrl}/shops/${id}/trust-badges`, { badges }),
    deleteShop: (id: string) => api.delete(`${baseUrl}/shops/${id}`),

    // Products
    listProducts: (shopSlug: string, page = 1, pageSize = 20, options: { category?: string, status?: string, search?: string } = {}) => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (options.category) params.append('category', options.category)
        if (options.status) params.append('status', options.status)
        if (options.search) params.append('search', options.search)
        return api.get(`${baseUrl}/shops/${shopSlug}/products?${params}`)
    },
    getProduct: (shopSlug: string, productSlug: string) => api.get(`${baseUrl}/shops/${shopSlug}/products/${productSlug}`),
    getProductById: (id: string) => api.get(`${baseUrl}/products/${id}`),
    createProduct: (data: Partial<Product> & { shop_id: string }) => api.post(`${baseUrl}/products`, data),
    updateProduct: (id: string, data: Partial<Product>) => api.put(`${baseUrl}/products/${id}`, data),
    deleteProduct: (id: string) => api.delete(`${baseUrl}/products/${id}`),

    // Categories
    listCategories: (shopSlug: string) => api.get(`${baseUrl}/shops/${shopSlug}/categories`),
    listCategoriesTree: (shopSlug: string) => api.get(`${baseUrl}/shops/${shopSlug}/categories/tree`),
    createCategory: (data: Partial<Category> & { shop_id: string }) => api.post(`${baseUrl}/categories`, data),
    updateCategory: (id: string, data: Partial<Category>) => api.put(`${baseUrl}/categories/${id}`, data),
    deleteCategory: (id: string) => api.delete(`${baseUrl}/categories/${id}`),

    // Orders
    createOrder: (data: { shop_id: string, items: CartItem[], wallet_id: string, delivery_type: string, shipping_address?: Address, notes?: string }) =>
        api.post(`${baseUrl}/orders`, data),
    listMyOrders: (page = 1, pageSize = 20) => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        return api.get(`${baseUrl}/orders?${params}`)
    },
    getOrder: (id: string) => api.get(`${baseUrl}/orders/${id}`),
    listShopOrders: (shopId: string, page = 1, pageSize = 20, status = '') => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (status) params.append('status', status)
        return api.get(`${baseUrl}/shop-orders/${shopId}?${params}`)
    },
    updateOrderStatus: (id: string, status: string, trackingNumber = '', sellerNotes = '') =>
        api.put(`${baseUrl}/orders/${id}/status`, { status, tracking_number: trackingNumber, seller_notes: sellerNotes }),
    refundOrder: (id: string, reason: string) => api.post(`${baseUrl}/orders/${id}/refund`, { reason }),

    // Upload
    uploadMedia: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return api.post(`${baseUrl}/upload`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        })
    },

    uploadDigitalFile: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return api.post(`${baseUrl}/upload/digital`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        })
    },

    getDigitalDownload: (id: string) => api.get(`${baseUrl}/products/${id}/digital-download`),

    // Manager management
    inviteManager: (shopId: string, email: string, role: string, permissions: string[]) =>
        api.post(`${baseUrl}/shops/${shopId}/managers`, { email, role, permissions }),
    removeManager: (shopId: string, userId: string) => api.delete(`${baseUrl}/shops/${shopId}/managers/${userId}`),

    // Client Invitations (for private shops)
    inviteClient: (shopId: string, data: { email: string, phone?: string, first_name?: string, last_name?: string, notes?: string, discount?: number }) =>
        api.post(`${baseUrl}/shops/${shopId}/clients`, data),
    listShopClients: (shopId: string, page = 1, pageSize = 20, status = '') => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (status) params.append('status', status)
        return api.get(`${baseUrl}/shops/${shopId}/clients?${params}`)
    },
    revokeClientAccess: (shopId: string, clientId: string) => api.delete(`${baseUrl}/shops/${shopId}/clients/${clientId}`),
    getMyInvitations: () => api.get(`${baseUrl}/my-invitations`),
    acceptInvitation: (invitationId: string, pin: string) => api.post(`${baseUrl}/invitations/accept`, { invitation_id: invitationId, pin }),
    declineInvitation: (invitationId: string) => api.delete(`${baseUrl}/invitations/${invitationId}`),
    getMyPrivateShops: () => api.get(`${baseUrl}/my-private-shops`),

    // Reviews
    listReviews: (productId: string, page = 1, pageSize = 10) => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        return api.get(`${baseUrl}/products/${productId}/reviews?${params}`)
    },
    createReview: (productId: string, data: { rating: number, comment: string }) =>
        api.post(`${baseUrl}/products/${productId}/reviews`, data),
}

// Composable wrapper for Vue usage
export function useShopApi() {
    return {
        // Shops
        listShops: async (page = 1, pageSize = 20, search = '') => {
            const response = await shopAPI.listShops(page, pageSize, search)
            return response.data
        },
        getShop: async (slug: string) => {
            const response = await shopAPI.getShop(slug)
            return response.data
        },
        getMyShops: async () => {
            const response = await shopAPI.getMyShops()
            return response.data
        },
        createShop: async (data: Partial<Shop>) => {
            const response = await shopAPI.createShop(data)
            return response.data
        },
        updateShop: async (id: string, data: Partial<Shop>) => {
            const response = await shopAPI.updateShop(id, data)
            return response.data
        },
        updateTrustBadges: async (id: string, badges: ShopTrustBadge[]) => {
            const response = await shopAPI.updateTrustBadges(id, badges)
            return response.data
        },
        deleteShop: async (id: string) => {
            const response = await shopAPI.deleteShop(id)
            return response.data
        },
        // Products
        listProducts: async (shopSlug: string, page = 1, pageSize = 20, options = {}) => {
            const response = await shopAPI.listProducts(shopSlug, page, pageSize, options)
            return response.data
        },
        getProduct: async (shopSlug: string, productSlug: string) => {
            const response = await shopAPI.getProduct(shopSlug, productSlug)
            return response.data
        },
        getProductById: async (id: string) => {
            const response = await shopAPI.getProductById(id)
            return response.data
        },
        createProduct: async (data: Partial<Product> & { shop_id: string }) => {
            const response = await shopAPI.createProduct(data)
            return response.data
        },
        updateProduct: async (id: string, data: Partial<Product>) => {
            const response = await shopAPI.updateProduct(id, data)
            return response.data
        },
        deleteProduct: async (id: string) => {
            const response = await shopAPI.deleteProduct(id)
            return response.data
        },
        // Categories
        listCategories: async (shopSlug: string) => {
            const response = await shopAPI.listCategories(shopSlug)
            return response.data
        },
        listCategoriesTree: async (shopSlug: string) => {
            const response = await shopAPI.listCategoriesTree(shopSlug)
            return response.data
        },
        createCategory: async (data: Partial<Category> & { shop_id: string }) => {
            const response = await shopAPI.createCategory(data)
            return response.data
        },
        updateCategory: async (id: string, data: Partial<Category>) => {
            const response = await shopAPI.updateCategory(id, data)
            return response.data
        },
        deleteCategory: async (id: string) => {
            const response = await shopAPI.deleteCategory(id)
            return response.data
        },
        // Orders
        createOrder: async (data: { shop_id: string, items: CartItem[], wallet_id: string, delivery_type: string, shipping_address?: Address, notes?: string }) => {
            const response = await shopAPI.createOrder(data)
            return response.data
        },
        listMyOrders: async (page = 1, pageSize = 20) => {
            const response = await shopAPI.listMyOrders(page, pageSize)
            return response.data
        },
        getOrder: async (id: string) => {
            const response = await shopAPI.getOrder(id)
            return response.data
        },
        listShopOrders: async (shopId: string, page = 1, pageSize = 20, status = '') => {
            const response = await shopAPI.listShopOrders(shopId, page, pageSize, status)
            return response.data
        },
        updateOrderStatus: async (id: string, status: string, trackingNumber = '', sellerNotes = '') => {
            const response = await shopAPI.updateOrderStatus(id, status, trackingNumber, sellerNotes)
            return response.data
        },
        refundOrder: async (id: string, reason: string) => {
            const response = await shopAPI.refundOrder(id, reason)
            return response.data
        },
        // Upload
        uploadMedia: async (file: File) => {
            const response = await shopAPI.uploadMedia(file)
            return response.data
        },

        uploadDigitalFile: async (file: File) => {
            const response = await shopAPI.uploadDigitalFile(file)
            return response.data
        },

        getDigitalDownload: async (id: string) => {
            const response = await shopAPI.getDigitalDownload(id)
            return response.data
        },
        // Managers
        inviteManager: async (shopId: string, email: string, role: string, permissions: string[]) => {
            const response = await shopAPI.inviteManager(shopId, email, role, permissions)
            return response.data
        },
        removeManager: async (shopId: string, userId: string) => {
            const response = await shopAPI.removeManager(shopId, userId)
            return response.data
        },
        // Client Invitations
        inviteClient: async (shopId: string, data: { email: string, phone?: string, first_name?: string, last_name?: string, notes?: string, discount?: number }) => {
            const response = await shopAPI.inviteClient(shopId, data)
            return response.data
        },
        listShopClients: async (shopId: string, page = 1, pageSize = 20, status = '') => {
            const response = await shopAPI.listShopClients(shopId, page, pageSize, status)
            return response.data
        },
        revokeClientAccess: async (shopId: string, clientId: string) => {
            const response = await shopAPI.revokeClientAccess(shopId, clientId)
            return response.data
        },
        getMyInvitations: async () => {
            const response = await shopAPI.getMyInvitations()
            return response.data
        },
        acceptInvitation: async (invitationId: string, pin: string) => {
            const response = await shopAPI.acceptInvitation(invitationId, pin)
            return response.data
        },
        declineInvitation: async (invitationId: string) => {
            const response = await shopAPI.declineInvitation(invitationId)
            return response.data
        },
        getMyPrivateShops: async () => {
            const response = await shopAPI.getMyPrivateShops()
            return response.data
        },
        // Reviews
        listReviews: async (productId: string, page = 1, pageSize = 10) => {
            const response = await shopAPI.listReviews(productId, page, pageSize)
            return response.data
        },
        createReview: async (productId: string, data: { rating: number, comment: string }) => {
            const response = await shopAPI.createReview(productId, data)
            return response.data
        },
    }
}
