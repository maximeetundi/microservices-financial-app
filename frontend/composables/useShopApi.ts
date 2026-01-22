import { useApi } from './useApi'

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
    created_at: string
    updated_at: string
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
    // Local UI enrichment
    product?: Product
}

export function useShopApi() {
    const { apiCall } = useApi()
    const baseUrl = 'shop-service/api/v1'

    // Shops
    const listShops = async (page = 1, pageSize = 20, search = '') => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (search) params.append('search', search)
        return apiCall<{ shops: Shop[], total: number, page: number, total_pages: number }>(`${baseUrl}/shops?${params}`)
    }

    const getShop = async (slug: string) => {
        return apiCall<Shop>(`${baseUrl}/shops/${slug}`)
    }

    const getMyShops = async () => {
        return apiCall<{ shops: Shop[] }>(`${baseUrl}/my-shops`)
    }

    const createShop = async (data: Partial<Shop>) => {
        return apiCall<Shop>(`${baseUrl}/shops`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    const updateShop = async (id: string, data: Partial<Shop>) => {
        return apiCall<Shop>(`${baseUrl}/shops/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        })
    }

    const deleteShop = async (id: string) => {
        return apiCall(`${baseUrl}/shops/${id}`, { method: 'DELETE' })
    }

    // Products
    const listProducts = async (shopSlug: string, page = 1, pageSize = 20, options: { category?: string, status?: string, search?: string } = {}) => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (options.category) params.append('category', options.category)
        if (options.status) params.append('status', options.status)
        if (options.search) params.append('search', options.search)
        return apiCall<{ products: Product[], total: number, page: number, total_pages: number }>(`${baseUrl}/shops/${shopSlug}/products?${params}`)
    }

    const getProduct = async (shopSlug: string, productSlug: string) => {
        return apiCall<Product>(`${baseUrl}/shops/${shopSlug}/products/${productSlug}`)
    }

    const createProduct = async (data: Partial<Product> & { shop_id: string }) => {
        return apiCall<Product>(`${baseUrl}/products`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    const updateProduct = async (id: string, data: Partial<Product>) => {
        return apiCall<Product>(`${baseUrl}/products/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        })
    }

    const deleteProduct = async (id: string) => {
        return apiCall(`${baseUrl}/products/${id}`, { method: 'DELETE' })
    }

    // Categories
    const listCategories = async (shopSlug: string) => {
        return apiCall<{ categories: Category[] }>(`${baseUrl}/shops/${shopSlug}/categories`)
    }

    const listCategoriesTree = async (shopSlug: string) => {
        return apiCall<{ categories: Category[] }>(`${baseUrl}/shops/${shopSlug}/categories/tree`)
    }

    const createCategory = async (data: Partial<Category> & { shop_id: string }) => {
        return apiCall<Category>(`${baseUrl}/categories`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    const updateCategory = async (id: string, data: Partial<Category>) => {
        return apiCall<Category>(`${baseUrl}/categories/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        })
    }

    const deleteCategory = async (id: string) => {
        return apiCall(`${baseUrl}/categories/${id}`, { method: 'DELETE' })
    }

    // Orders
    const createOrder = async (data: { shop_id: string, items: CartItem[], wallet_id: string, delivery_type: string, shipping_address?: Address, notes?: string }) => {
        return apiCall<Order>(`${baseUrl}/orders`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    const listMyOrders = async (page = 1, pageSize = 20) => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        return apiCall<{ orders: Order[], total: number, page: number, total_pages: number }>(`${baseUrl}/orders?${params}`)
    }

    const getOrder = async (id: string) => {
        return apiCall<Order>(`${baseUrl}/orders/${id}`)
    }

    const listShopOrders = async (shopId: string, page = 1, pageSize = 20, status = '') => {
        const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
        if (status) params.append('status', status)
        return apiCall<{ orders: Order[], total: number, page: number, total_pages: number }>(`${baseUrl}/shop-orders/${shopId}?${params}`)
    }

    const updateOrderStatus = async (id: string, status: string, trackingNumber = '', sellerNotes = '') => {
        return apiCall<Order>(`${baseUrl}/orders/${id}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status, tracking_number: trackingNumber, seller_notes: sellerNotes }),
        })
    }

    const refundOrder = async (id: string, reason: string) => {
        return apiCall(`${baseUrl}/orders/${id}/refund`, {
            method: 'POST',
            body: JSON.stringify({ reason }),
        })
    }

    // Upload
    const uploadMedia = async (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return apiCall<{ url: string }>(`${baseUrl}/upload`, {
            method: 'POST',
            body: formData,
            headers: {}, // Let browser set content-type for FormData
        })
    }

    // Manager management
    const inviteManager = async (shopId: string, email: string, role: string, permissions: string[]) => {
        return apiCall(`${baseUrl}/shops/${shopId}/managers`, {
            method: 'POST',
            body: JSON.stringify({ email, role, permissions }),
        })
    }

    const removeManager = async (shopId: string, userId: string) => {
        return apiCall(`${baseUrl}/shops/${shopId}/managers/${userId}`, { method: 'DELETE' })
    }

    return {
        // Shops
        listShops,
        getShop,
        getMyShops,
        createShop,
        updateShop,
        deleteShop,
        // Products
        listProducts,
        getProduct,
        createProduct,
        updateProduct,
        deleteProduct,
        // Categories
        listCategories,
        listCategoriesTree,
        createCategory,
        updateCategory,
        deleteCategory,
        // Orders
        createOrder,
        listMyOrders,
        getOrder,
        listShopOrders,
        updateOrderStatus,
        refundOrder,
        // Upload
        uploadMedia,
        // Managers
        inviteManager,
        removeManager,
    }
}
