import axios from 'axios'

// Flag to prevent infinite redirect loop
let isLoggingOut = false

// Function to reset the logout flag (called after successful login)
export const resetLogoutFlag = () => {
    isLoggingOut = false
}

// Use runtime config in Nuxt context, fallback to production API
const getApiUrl = () => {
    if (typeof window !== 'undefined') {
        // Client-side: use the API URL
        return 'https://api.app.maximeetundi.store'
    }
    // Server-side or fallback
    return process.env.API_BASE_URL || 'http://api-gateway:8080'
}

const API_URL = getApiUrl()

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Request interceptor
api.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('accessToken')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
    }
    return config
})

// Response interceptor
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true

            const refreshToken = localStorage.getItem('refreshToken')
            if (refreshToken) {
                try {
                    const response = await axios.post(`${API_URL}/auth-service/api/v1/auth/refresh`, {
                        refresh_token: refreshToken
                    })
                    const { access_token, refresh_token } = response.data
                    localStorage.setItem('accessToken', access_token)
                    localStorage.setItem('refreshToken', refresh_token)
                    originalRequest.headers.Authorization = `Bearer ${access_token}`
                    // Reset logout flag on successful refresh
                    isLoggingOut = false
                    return api(originalRequest)
                } catch {
                    // Prevent infinite loop - only redirect once
                    if (!isLoggingOut) {
                        isLoggingOut = true
                        localStorage.removeItem('accessToken')
                        localStorage.removeItem('refreshToken')
                        if (typeof window !== 'undefined') {
                            window.location.href = '/auth/login'
                        }
                    }
                }
            } else if (!isLoggingOut) {
                // No refresh token, redirect to login
                isLoggingOut = true
                localStorage.removeItem('accessToken')
                localStorage.removeItem('refreshToken')
                if (typeof window !== 'undefined') {
                    window.location.href = '/auth/login'
                }
            }
        }
        return Promise.reject(error)
    }
)

// ========== Auth ==========
export const authAPI = {
    login: (email: string, password: string, twoFaCode?: string) =>
        api.post('/auth-service/api/v1/auth/login', { email, password, two_fa_code: twoFaCode }),
    register: (data: any) => api.post('/auth-service/api/v1/auth/register', data),
    refresh: (refreshToken: string) => api.post('/auth-service/api/v1/auth/refresh', { refresh_token: refreshToken }),
    logout: () => api.post('/auth-service/api/v1/auth/logout'),
    forgotPassword: (email: string) => api.post('/auth-service/api/v1/auth/forgot-password', { email }),
    resetPassword: (token: string, password: string) =>
        api.post('/auth-service/api/v1/auth/reset-password', { token, new_password: password }),
}

// ========== User ==========
export const userAPI = {
    getProfile: () => api.get('/auth-service/api/v1/users/profile'),
    updateProfile: (data: any) => api.put('/auth-service/api/v1/users/profile', data),
    changePassword: (data: { current_password: string, new_password: string }) =>
        api.post('/auth-service/api/v1/auth/change-password', data),

    // 2FA
    enable2FA: () => api.post('/auth-service/api/v1/auth/2fa/enable'),
    verify2FA: (data: { code: string }) => api.post('/auth-service/api/v1/auth/2fa/verify', data),
    disable2FA: () => api.post('/auth-service/api/v1/auth/2fa/disable'),

    // Sessions
    getSessions: () => api.get('/auth-service/api/v1/auth/sessions'),
    revokeSession: (sessionId: string) => api.delete(`/auth-service/api/v1/auth/sessions/${sessionId}`),
    revokeAllSessions: () => api.delete('/auth-service/api/v1/auth/sessions'),

    // User lookup
    lookup: (query: { email?: string, phone?: string }) => {
        const params = new URLSearchParams()
        if (query.email) params.append('email', query.email)
        if (query.phone) params.append('phone', query.phone)
        return api.get(`/auth-service/api/v1/users/lookup?${params.toString()}`)
    },
}

// ========== Wallets ==========
export const walletAPI = {
    getAll: () => api.get('/wallet-service/api/v1/wallets'),
    getWallets: () => api.get('/wallet-service/api/v1/wallets'),
    get: (id: string) => api.get(`/wallet-service/api/v1/wallets/${id}`),
    create: (data: { currency: string, wallet_type: string, name?: string, description?: string }) =>
        api.post('/wallet-service/api/v1/wallets', data),
    getBalance: (id: string) => api.get(`/wallet-service/api/v1/wallets/${id}/balance`),
    getTransactions: (id: string, limit = 50, offset = 0) =>
        api.get(`/wallet-service/api/v1/wallets/${id}/transactions?limit=${limit}&offset=${offset}`),
    deposit: (id: string, amount: number, method: string) =>
        api.post(`/wallet-service/api/v1/wallets/${id}/deposit`, { amount, method }),
    withdraw: (id: string, amount: number, destination: string) =>
        api.post(`/wallet-service/api/v1/wallets/${id}/withdraw`, { amount, destination }),
}

// ========== Transfers ==========
export const transferAPI = {
    getAll: (limit = 50, offset = 0) => api.get(`/transfer-service/api/v1/transfers?limit=${limit}&offset=${offset}`),
    get: (id: string) => api.get(`/transfer-service/api/v1/transfers/${id}`),
    create: (data: {
        type: string
        amount: number
        currency: string
        recipient: string
        description?: string
    }) => api.post('/transfer-service/api/v1/transfers', data),
    getBanks: (country: string) => api.get(`/transfer-service/api/v1/transfers/banks?country=${country}`),
    getMobileOperators: (country: string) => api.get(`/transfer-service/api/v1/transfers/mobile-operators?country=${country}`),
    validateRecipient: (type: string, recipient: string) =>
        api.post('/transfer-service/api/v1/transfers/validate-recipient', { type, recipient }),
    getFees: (type: string, amount: number, currency: string) =>
        api.get(`/transfer-service/api/v1/transfers/fees?type=${type}&amount=${amount}&currency=${currency}`),
}

// ========== Cards ==========
export const cardAPI = {
    getAll: () => api.get('/card-service/api/v1/cards'),
    get: (id: string) => api.get(`/card-service/api/v1/cards/${id}`),
    create: (data: { type: string; currency: string; name: string }) => api.post('/card-service/api/v1/cards', data),
    activate: (id: string) => api.post(`/card-service/api/v1/cards/${id}/activate`),
    freeze: (id: string) => api.post(`/card-service/api/v1/cards/${id}/freeze`),
    unfreeze: (id: string) => api.post(`/card-service/api/v1/cards/${id}/unfreeze`),
    setLimit: (id: string, limitType: string, amount: number) =>
        api.post(`/card-service/api/v1/cards/${id}/limits`, { limit_type: limitType, amount }),
    setPin: (id: string, pin: string) => api.post(`/card-service/api/v1/cards/${id}/pin`, { pin }),
    getTransactions: (id: string, limit = 50) =>
        api.get(`/card-service/api/v1/cards/${id}/transactions?limit=${limit}`),
    topUp: (id: string, amount: number, sourceWalletId: string) =>
        api.post(`/card-service/api/v1/cards/${id}/topup`, { amount, source_wallet_id: sourceWalletId }),
}

// ========== Exchange ==========
// ========== Exchange ==========
export const exchangeAPI = {
    getMarkets: () => api.get('/exchange-service/api/v1/markets'),
    getRates: () => api.get('/exchange-service/api/v1/rates'),
    getRate: (from: string, to: string) => api.get(`/exchange-service/api/v1/rates/${from}/${to}`),

    // Quote and Execute flow
    getQuote: (fromCurrency: string, toCurrency: string, amount: number, side: 'from' | 'to' = 'from') => {
        const payload: any = {
            from_currency: fromCurrency,
            to_currency: toCurrency,
        }
        if (side === 'from') payload.from_amount = amount
        else payload.to_amount = amount
        return api.post('/exchange-service/api/v1/quote', payload)
    },

    executeExchange: (quoteId: string, fromWalletId: string, toWalletId: string) =>
        api.post('/exchange-service/api/v1/execute', {
            quote_id: quoteId,
            from_wallet_id: fromWalletId,
            to_wallet_id: toWalletId
        }),

    // Legacy/Composite convert for simple UI (may need to be implemented in backend or stitched here)
    // For now, replacing convert with getQuote+execute would require UI changes. 
    // If UI expects single call, we should probably keep distinct functions or fix UI.
    // Given the UI in index.vue calls convert(), let's map it to getQuote for now to at least return the "to_amount".
    convert: async (fromCurrency: string, toCurrency: string, amount: number) => {
        // This is a helper for the quick converter which just wants the rate/amount
        return api.post('/exchange-service/api/v1/quote', {
            from_currency: fromCurrency,
            to_currency: toCurrency,
            from_amount: amount
        })
    },

    getCryptoRates: () => api.get('/exchange-service/api/v1/rates'), // Fallback to all rates

    buyCrypto: (currency: string, amount: number, paymentMethod: string, orderType = 'market', limitPrice?: number) =>
        api.post('/exchange-service/api/v1/trading/buy', {
            currency, amount, payment_method: paymentMethod, order_type: orderType, limit_price: limitPrice
        }),

    sellCrypto: (currency: string, amount: number, destinationWalletId: string, orderType = 'market', limitPrice?: number) =>
        api.post('/exchange-service/api/v1/trading/sell', {
            currency, amount, destination_wallet_id: destinationWalletId, order_type: orderType, limit_price: limitPrice
        }),

    getHistory: (limit = 50) => api.get(`/exchange-service/api/v1/history?limit=${limit}`),
    getOrders: () => api.get('/exchange-service/api/v1/trading/orders'),
    getTradingPortfolio: () => api.get('/exchange-service/api/v1/trading/portfolio'),
}

// ========== Dashboard ==========
export const dashboardAPI = {
    getSummary: () => api.get('/wallet-service/api/v1/dashboard/summary'),
    getRecentActivity: (limit = 10) => api.get(`/wallet-service/api/v1/dashboard/activity?limit=${limit}`),
    getPortfolio: () => api.get('/wallet-service/api/v1/dashboard/portfolio'),
    getStats: (period: string) => api.get(`/wallet-service/api/v1/dashboard/stats?period=${period}`),
}

// ========== Merchant Payments ==========
export const merchantAPI = {
    // Create a payment request
    createPayment: (data: {
        type: 'fixed' | 'variable' | 'invoice'
        wallet_id: string
        amount?: number
        min_amount?: number
        max_amount?: number
        currency: string
        title: string
        description?: string
        expires_in_minutes?: number
        reusable?: boolean
    }) => api.post('/wallet-service/api/v1/merchant/payments', data),

    // Get all payment requests
    getPayments: (limit = 20, offset = 0) =>
        api.get(`/wallet-service/api/v1/merchant/payments?limit=${limit}&offset=${offset}`),

    // Get payment history
    getHistory: (limit = 20, offset = 0) =>
        api.get(`/wallet-service/api/v1/merchant/payments/history?limit=${limit}&offset=${offset}`),

    // Cancel a payment request
    cancelPayment: (id: string) => api.delete(`/wallet-service/api/v1/merchant/payments/${id}`),

    // Get QR code for a payment
    getQRCode: (id: string) => api.get(`/wallet-service/api/v1/payments/${id}/qr`),

    // Quick payment (simplified)
    quickPay: (data: {
        wallet_id: string
        amount?: number
        currency: string
        description?: string
        never_expires?: boolean
    }) => api.post('/wallet-service/api/v1/merchant/quick-pay', data),
}

// ========== Payment (for customers) ==========
export const paymentAPI = {
    // Get payment details (public - for scanning)
    getPaymentDetails: (id: string) => api.get(`/wallet-service/api/v1/pay/${id}`),

    // Pay a payment request
    payPayment: (id: string, data: {
        from_wallet_id: string
        amount?: number
    }) => api.post(`/wallet-service/api/v1/payments/${id}/pay`, data),

    // Get available payment methods for a country
    getPaymentMethods: (countryCode = 'CI') =>
        api.get(`/admin-service/api/v1/payment-methods?country=${countryCode}`),
}

// ========== Admin Payment Providers ==========
export const adminPaymentAPI = {
    // Get all payment providers
    getProviders: () => api.get('/admin-service/api/v1/admin/payment-providers'),

    // Get single provider
    getProvider: (id: string) => api.get(`/admin-service/api/v1/admin/payment-providers/${id}`),

    // Create provider
    createProvider: (data: any) => api.post('/admin-service/api/v1/admin/payment-providers', data),

    // Update provider
    updateProvider: (id: string, data: any) =>
        api.put(`/admin-service/api/v1/admin/payment-providers/${id}`, data),

    // Delete provider
    deleteProvider: (id: string) =>
        api.delete(`/admin-service/api/v1/admin/payment-providers/${id}`),

    // Toggle active status
    toggleStatus: (id: string, isActive: boolean) =>
        api.post(`/admin-service/api/v1/admin/payment-providers/${id}/toggle-status`, { is_active: isActive }),

    // Toggle demo mode
    toggleDemo: (id: string, isDemoMode: boolean) =>
        api.post(`/admin-service/api/v1/admin/payment-providers/${id}/toggle-demo`, { is_demo_mode: isDemoMode }),

    // Test connection
    testConnection: (id: string) =>
        api.post(`/admin-service/api/v1/admin/payment-providers/${id}/test`),

    // Add country
    addCountry: (id: string, country: any) =>
        api.post(`/admin-service/api/v1/admin/payment-providers/${id}/countries`, country),

    // Remove country
    removeCountry: (id: string, countryCode: string) =>
        api.delete(`/admin-service/api/v1/admin/payment-providers/${id}/countries/${countryCode}`),
}

// ========== Support ==========
export const supportAPI = {
    // Get all tickets for the current user
    getTickets: (limit = 20, offset = 0) =>
        api.get(`/support-service/api/v1/tickets?limit=${limit}&offset=${offset}`),

    // Get a specific ticket
    getTicket: (id: string) => api.get(`/support-service/api/v1/tickets/${id}`),

    // Create a new support ticket
    createTicket: (data: {
        subject: string
        description: string
        category: string
        priority?: string
    }) => api.post('/support-service/api/v1/tickets', data),

    // Add a message to a ticket
    sendMessage: (ticketId: string, message: string) =>
        api.post(`/support-service/api/v1/tickets/${ticketId}/messages`, { message }),

    // Get messages for a ticket
    getMessages: (ticketId: string) =>
        api.get(`/support-service/api/v1/tickets/${ticketId}/messages`),

    // Close a ticket
    closeTicket: (ticketId: string) =>
        api.post(`/support-service/api/v1/tickets/${ticketId}/close`),
}

// ========== Composable ==========
export function useApi() {
    return {
        authApi: authAPI,
        userApi: userAPI,
        walletApi: walletAPI,
        transferApi: transferAPI,
        cardApi: cardAPI,
        exchangeApi: exchangeAPI,
        dashboardApi: dashboardAPI,
        merchantApi: merchantAPI,
        paymentApi: paymentAPI,
        adminPaymentApi: adminPaymentAPI,
        supportApi: supportAPI,
        notificationApi: notificationAPI,
    }
}

// ========== Notifications ==========
export const notificationAPI = {
    getAll: (limit = 20, offset = 0) => api.get(`/notification-service/api/v1/notifications?limit=${limit}&offset=${offset}`),
    getUnreadCount: () => api.get('/notification-service/api/v1/notifications/unread-count'),
    markAsRead: (id: string) => api.post(`/notification-service/api/v1/notifications/${id}/read`),
    markAllAsRead: () => api.post('/notification-service/api/v1/notifications/read-all'),
    delete: (id: string) => api.delete(`/notification-service/api/v1/notifications/${id}`),
}

export default api

