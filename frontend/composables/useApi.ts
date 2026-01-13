import axios from 'axios'

// Flag to prevent infinite redirect loop
let isLoggingOut = false

// Function to reset the logout and refresh flags (called after successful login)
export const resetLogoutFlag = () => {
    isLoggingOut = false
    // Also reset refresh state in case there was a pending refresh
    isRefreshingGlobal = false
    refreshSubscribersGlobal = []
}

// Global refresh state (declared here, used in interceptor)
let isRefreshingGlobal = false
let refreshSubscribersGlobal: Array<(token: string) => void> = []

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

// List of endpoints that should NOT trigger auto-logout on 401
const AUTH_ENDPOINTS = [
    '/auth/login',
    '/auth/register',
    '/auth/refresh',
    '/auth/forgot-password',
    '/auth/reset-password'
]

// Check if a URL is an auth endpoint
const isAuthEndpoint = (url: string): boolean => {
    return AUTH_ENDPOINTS.some(endpoint => url.includes(endpoint))
}

// Subscribe to refresh completion
const subscribeToRefresh = (callback: (token: string) => void) => {
    refreshSubscribersGlobal.push(callback)
}

// Notify all subscribers with new token
const onRefreshComplete = (token: string) => {
    refreshSubscribersGlobal.forEach(callback => callback(token))
    refreshSubscribersGlobal = []
}

// Notify all subscribers of refresh failure
const onRefreshFailed = () => {
    refreshSubscribersGlobal = []
}

// Response interceptor
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config
        const requestUrl = originalRequest?.url || ''

        // Skip 401 handling for auth endpoints to prevent logout loops
        if (isAuthEndpoint(requestUrl)) {
            return Promise.reject(error)
        }

        if (error.response?.status === 401 && !originalRequest._retry && !isLoggingOut) {
            originalRequest._retry = true

            const refreshToken = localStorage.getItem('refreshToken')
            if (!refreshToken) {
                // No refresh token available, must logout
                if (!isLoggingOut) {
                    console.warn('No refresh token available, logging out')
                    isLoggingOut = true
                    localStorage.removeItem('accessToken')
                    localStorage.removeItem('refreshToken')
                    if (typeof document !== 'undefined') {
                        document.cookie = 'accessToken=; path=/; max-age=0'
                        document.cookie = 'refreshToken=; path=/; max-age=0'
                    }
                    if (typeof window !== 'undefined') {
                        window.location.href = '/auth/login'
                    }
                }
                return Promise.reject(error)
            }

            // If already refreshing, wait for the refresh to complete
            if (isRefreshingGlobal) {
                return new Promise((resolve) => {
                    subscribeToRefresh((token: string) => {
                        originalRequest.headers.Authorization = `Bearer ${token}`
                        resolve(api(originalRequest))
                    })
                })
            }

            // Start refreshing
            isRefreshingGlobal = true

            try {
                const response = await axios.post(`${API_URL}/auth-service/api/v1/auth/refresh`, {
                    refresh_token: refreshToken
                })
                const { access_token, refresh_token } = response.data

                // Update tokens in localStorage
                localStorage.setItem('accessToken', access_token)
                localStorage.setItem('refreshToken', refresh_token)

                // Update cookies for SSR compatibility
                if (typeof document !== 'undefined') {
                    document.cookie = `accessToken=${access_token}; path=/; max-age=86400; SameSite=Lax`
                    document.cookie = `refreshToken=${refresh_token}; path=/; max-age=604800; SameSite=Lax`
                }

                // Reset logout flag on successful refresh
                isLoggingOut = false
                isRefreshingGlobal = false

                // Notify all waiting requests
                onRefreshComplete(access_token)

                originalRequest.headers.Authorization = `Bearer ${access_token}`
                return api(originalRequest)
            } catch (refreshError) {
                console.warn('Token refresh failed:', refreshError)
                isRefreshingGlobal = false
                onRefreshFailed()

                // Only logout if not already logging out
                if (!isLoggingOut) {
                    isLoggingOut = true
                    localStorage.removeItem('accessToken')
                    localStorage.removeItem('refreshToken')
                    if (typeof document !== 'undefined') {
                        document.cookie = 'accessToken=; path=/; max-age=0'
                        document.cookie = 'refreshToken=; path=/; max-age=0'
                    }
                    if (typeof window !== 'undefined') {
                        window.location.href = '/auth/login'
                    }
                }
                return Promise.reject(refreshError)
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
    getPublicKey: () => api.get('/auth-service/api/v1/auth/public-key'),
}

// ========== User ==========
export const userAPI = {
    getProfile: () => api.get('/auth-service/api/v1/users/profile'),
    updateProfile: (data: any) => api.put('/auth-service/api/v1/users/profile', data),
    changePassword: (data: { current_password: string, new_password: string }) =>
        api.post('/auth-service/api/v1/users/change-password', data),

    // 2FA - routes are under /users/2fa/ in backend
    enable2FA: () => api.post('/auth-service/api/v1/users/2fa/setup'),
    verify2FA: (data: { code: string, secret?: string }) => api.post('/auth-service/api/v1/users/2fa/verify', data),
    disable2FA: (data?: { code: string }) => api.post('/auth-service/api/v1/users/2fa/disable', data || {}),

    // Sessions
    getSessions: () => api.get('/auth-service/api/v1/sessions'),
    revokeSession: (sessionId: string) => api.delete(`/auth-service/api/v1/sessions/${sessionId}`),
    revokeAllSessions: () => api.delete('/auth-service/api/v1/sessions'),

    // PIN (5-digit transaction security PIN)
    checkPinStatus: () => api.get(`/auth-service/api/v1/users/pin/status?t=${Date.now()}`),
    setupPin: (data: { pin: string, confirm_pin: string }) =>
        api.post('/auth-service/api/v1/users/pin/setup', data),
    verifyPin: (data: { pin: string }) =>
        api.post('/auth-service/api/v1/users/pin/verify', data),
    changePin: (data: { current_pin: string, new_pin: string, confirm_pin: string }) =>
        api.post('/auth-service/api/v1/users/pin/change', data),

    // User lookup
    getById: (id: string) => api.get(`/auth-service/api/v1/users/${id}`),

    lookup: (query: { email?: string, phone?: string }) => {
        const params = new URLSearchParams()
        if (query.email) params.append('email', query.email)
        if (query.phone) params.append('phone', query.phone)
        return api.get(`/auth-service/api/v1/users/lookup?${params.toString()}`)
    },

    // KYC
    getKYCStatus: () => api.get('/auth-service/api/v1/users/kyc/status'),
    uploadKYCDocument: (formData: FormData) =>
        api.post('/auth-service/api/v1/users/kyc/documents', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        }),
    getKYCDocuments: () => api.get('/auth-service/api/v1/users/kyc/documents'),

    // Preferences
    getPreferences: () => api.get('/auth-service/api/v1/users/preferences'),
    updatePreferences: (data: any) => api.put('/auth-service/api/v1/users/preferences', data),

    // Notification Preferences
    getNotificationPrefs: () => api.get('/auth-service/api/v1/users/notifications/preferences'),
    updateNotificationPrefs: (data: any) => api.put('/auth-service/api/v1/users/notifications/preferences', data),
}

// ========== Wallets ==========
export const walletAPI = {
    getAll: () => api.get('/wallet-service/api/v1/wallets'),
    getWallets: () => api.get('/wallet-service/api/v1/wallets'),
    getMyWallets: () => api.get('/wallet-service/api/v1/wallets'),
    get: (id: string) => api.get(`/wallet-service/api/v1/wallets/${id}`),
    create: (data: { currency: string, wallet_type: string, name?: string, description?: string }) =>
        api.post('/wallet-service/api/v1/wallets', data),
    createWallet: (data: { currency: string, wallet_type: string, name?: string, description?: string }) =>
        api.post('/wallet-service/api/v1/wallets', data),
    getBalance: (id: string) => api.get(`/wallet-service/api/v1/wallets/${id}/balance`),
    getTransactions: (id: string, limit = 50, offset = 0) =>
        api.get(`/wallet-service/api/v1/wallets/${id}/transactions?limit=${limit}&offset=${offset}`),
    deposit: (id: string, amount: number, method: string) =>
        api.post(`/wallet-service/api/v1/wallets/${id}/deposit`, { amount, method }),
    withdraw: (id: string, amount: number, destination: string) =>
        api.post(`/wallet-service/api/v1/wallets/${id}/withdraw`, { amount, destination }),

    // Crypto wallet methods
    generateCryptoWallet: (currency: string) =>
        api.post('/wallet-service/api/v1/crypto/generate', { currency, wallet_type: 'crypto', name: `${currency} Wallet` }),
    getCryptoAddress: (walletId: string) =>
        api.get(`/wallet-service/api/v1/crypto/${walletId}/address`),
    sendCrypto: (walletId: string, data: { to_address: string, amount: number, note?: string }) =>
        api.post(`/wallet-service/api/v1/crypto/${walletId}/send`, data),
    estimateCryptoFee: (walletId: string, amount: number) =>
        api.post(`/wallet-service/api/v1/crypto/${walletId}/estimate-fee`, { amount }),
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

    // Cancellation & Reversal
    cancelTransfer: (id: string, reason: string) => api.post(`/transfer-service/api/v1/transfers/${id}/cancel`, { reason }),
    reverseTransfer: (id: string, reason: string) => api.post(`/transfer-service/api/v1/transfers/${id}/reverse`, { reason }),
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

// ========== Fiat Exchange (Fiat-to-Fiat) ==========
export const fiatAPI = {
    // Get quote for fiat-to-fiat conversion
    getQuote: (fromCurrency: string, toCurrency: string, amount: number) =>
        api.post('/exchange-service/api/v1/fiat/quote', {
            from_currency: fromCurrency,
            to_currency: toCurrency,
            amount: amount
        }),

    // Execute fiat-to-fiat exchange (using quote is optional, direct exchange supported)
    executeExchange: (data: {
        from_wallet_id: string,
        to_wallet_id: string,
        from_currency: string,
        to_currency: string,
        amount: number
    }) => api.post('/exchange-service/api/v1/fiat/exchange', data),

    // Get fiat rates
    getRates: (base = 'USD') => api.get(`/exchange-service/api/v1/fiat/rates?base=${base}`),

    // Get specific fiat rate
    getRate: (from: string, to: string) => api.get(`/exchange-service/api/v1/fiat/rates/${from}/${to}`),

    // Get supported currencies
    getCurrencies: () => api.get('/exchange-service/api/v1/fiat/currencies'),

    // Fiat converter (quick calculation)
    convert: (from: string, to: string, amount: number) =>
        api.get(`/exchange-service/api/v1/fiat/convert?from=${from}&to=${to}&amount=${amount}`),
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

// ========== Admin Fees ==========
export const adminFeeAPI = {
    getWalletFees: () => api.get('/wallet-service/api/v1/admin/fees'),
    updateWalletFee: (data: any) => api.put('/wallet-service/api/v1/admin/fees', data),
    getExchangeFees: () => api.get('/exchange-service/api/v1/admin/fees'),
    updateExchangeFee: (data: any) => api.put('/exchange-service/api/v1/admin/fees', data),
}

// ========== Support ==========
export const supportAPI = {
    // Get all conversations for the current user
    getTickets: (limit = 20, offset = 0) =>
        api.get(`/support-service/api/v1/support/conversations?limit=${limit}&offset=${offset}`),

    // Get a specific conversation
    getTicket: (id: string) => api.get(`/support-service/api/v1/support/conversations/${id}`),

    // Create a new support conversation
    createTicket: (data: {
        subject: string
        description: string
        category: string
        priority?: string
        agent_type?: 'ai' | 'human'
    }) => api.post('/support-service/api/v1/support/conversations', {
        agent_type: data.agent_type || 'ai',
        subject: data.subject,
        category: data.category,
        message: data.description,
        priority: data.priority || 'medium'
    }),

    // Send a message to a conversation
    sendMessage: (conversationId: string, message: string, attachments: string[] = []) =>
        api.post(`/support-service/api/v1/support/conversations/${conversationId}/messages`, {
            content: message,
            content_type: 'text',
            attachments: attachments
        }),

    // Upload a file
    uploadFile: (formData: FormData) =>
        api.post('/support-service/api/v1/support/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        }),

    // Get messages for a conversation
    getMessages: (conversationId: string) =>
        api.get(`/support-service/api/v1/support/conversations/${conversationId}/messages`),

    // Escalate to human agent
    escalate: (conversationId: string, reason: string) =>
        api.post(`/support-service/api/v1/support/conversations/${conversationId}/escalate`, { reason }),

    // Close a conversation with rating
    closeTicket: (conversationId: string, rating?: number, feedback?: string) =>
        api.put(`/support-service/api/v1/support/conversations/${conversationId}/close`, { rating, feedback }),
}




// ========== Notifications ==========
export const notificationAPI = {
    getAll: (options: { limit?: number; offset?: number } = {}) => {
        const limit = options.limit || 20
        const offset = options.offset || 0
        return api.get(`/notification-service/api/v1/notifications?limit=${limit}&offset=${offset}`)
    },
    getUnreadCount: () => api.get('/notification-service/api/v1/notifications/unread-count'),
    markAsRead: (id: string) => api.post(`/notification-service/api/v1/notifications/${id}/read`),
    markAllAsRead: () => api.post('/notification-service/api/v1/notifications/read-all'),
    delete: (id: string) => api.delete(`/notification-service/api/v1/notifications/${id}`),
}

// ========== Tickets & Events ==========
export const ticketAPI = {
    // Event management (organizer)
    createEvent: (data: any) => api.post('/ticket-service/api/v1/events', data),
    getMyEvents: (limit = 20, offset = 0) => api.get(`/ticket-service/api/v1/events?limit=${limit}&offset=${offset}`),
    getEvent: (id: string) => api.get(`/ticket-service/api/v1/events/${id}`),
    getEventByCode: (code: string) => api.get(`/ticket-service/api/v1/events/code/${code}`),
    updateEvent: (id: string, data: any) => api.put(`/ticket-service/api/v1/events/${id}`, data),
    deleteEvent: (id: string) => api.delete(`/ticket-service/api/v1/events/${id}`),
    publishEvent: (id: string) => api.post(`/ticket-service/api/v1/events/${id}/publish`),
    getEventTickets: (id: string, limit = 50, offset = 0) => api.get(`/ticket-service/api/v1/events/${id}/tickets?limit=${limit}&offset=${offset}`),
    getEventStats: (id: string) => api.get(`/ticket-service/api/v1/events/${id}/stats`),

    // Public events
    getActiveEvents: (limit = 20, offset = 0) => api.get(`/ticket-service/api/v1/events/active?limit=${limit}&offset=${offset}`),
    getPublicEvent: (id: string) => api.get(`/ticket-service/api/v1/events/public/${id}`),

    // Ticket purchase
    purchaseTicket: (data: { event_id: string; tier_id: string; form_data: any; wallet_id: string; pin: string }) =>
        api.post('/ticket-service/api/v1/tickets/purchase', data),
    getMyTickets: (limit = 20, offset = 0) => api.get(`/ticket-service/api/v1/tickets?limit=${limit}&offset=${offset}`),
    getTicket: (id: string) => api.get(`/ticket-service/api/v1/tickets/${id}`),

    // Ticket verification
    verifyTicket: (code: string) => api.post('/ticket-service/api/v1/tickets/verify', { ticket_code: code }),
    useTicket: (id: string) => api.post(`/ticket-service/api/v1/tickets/${id}/use`),

    // Available icons
    getIcons: () => api.get('/ticket-service/api/v1/icons'),

    // Upload image
    uploadImage: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return api.post('/ticket-service/api/v1/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
    },

    // Refund and Cancellation
    refundTicket: (id: string, reason: string) => api.post(`/ticket-service/api/v1/tickets/${id}/refund`, { reason }),
    cancelEvent: (id: string, reason: string) => api.post(`/ticket-service/api/v1/events/${id}/cancel`, { reason }),
}

// ========== Donations ==========
export const donationAPI = {
    // Campaigns
    createCampaign: (data: any) => api.post('/donation-service/api/v1/campaigns', data),
    getCampaigns: (limit = 20, offset = 0) => api.get(`/donation-service/api/v1/campaigns?limit=${limit}&offset=${offset}`),
    getMyCampaigns: (creatorID: string) => api.get(`/donation-service/api/v1/campaigns?creator_id=${creatorID}`),
    getCampaign: (id: string) => api.get(`/donation-service/api/v1/campaigns/${id}`),
    updateCampaign: (id: string, data: any) => api.put(`/donation-service/api/v1/campaigns/${id}`, data),
    cancelCampaign: (id: string, reason: string) => api.post(`/donation-service/api/v1/campaigns/${id}/cancel`, { reason }),

    // Donations
    initiateDonation: (data: any) => api.post('/donation-service/api/v1/donations', data),
    refundDonation: (id: string, reason: string) => api.post(`/donation-service/api/v1/donations/${id}/refund`, { reason }),
    getDonations: (campaignID: string, limit = 20, offset = 0) =>
        api.get(`/donation-service/api/v1/donations?campaign_id=${campaignID}&limit=${limit}&offset=${offset}`),

    // Upload
    uploadMedia: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return api.post('/donation-service/api/v1/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
    },
}

// === CONTACTS API ===
export const contactsAPI = {
    // Get all user's synced contacts
    getAll: () => api.get('/auth-service/api/v1/users/contacts'),

    // Add a single contact manually
    add: (data: { phone?: string; email?: string; name: string }) =>
        api.post('/auth-service/api/v1/users/contacts', data),

    // Bulk sync contacts from mobile
    sync: (contacts: Array<{ phone?: string; email?: string; name: string }>) =>
        api.post('/auth-service/api/v1/users/contacts/sync', { contacts }),

    // Delete a contact
    delete: (id: string) => api.delete(`/auth-service/api/v1/users/contacts/${id}`),

    // Lookup contact name by phone or email
    lookup: (phone?: string, email?: string) =>
        api.get('/auth-service/api/v1/users/contacts/lookup', { params: { phone, email } }),
}

export const useApi = () => {
    return {
        authApi: authAPI,
        userApi: userAPI,
        walletApi: walletAPI,
        transferApi: transferAPI,
        cardApi: cardAPI,
        exchangeApi: exchangeAPI,
        fiatApi: fiatAPI,
        dashboardApi: dashboardAPI,
        merchantApi: merchantAPI,
        paymentApi: paymentAPI,
        adminPaymentApi: adminPaymentAPI,
        adminFeeApi: adminFeeAPI,
        supportApi: supportAPI,
        notificationApi: notificationAPI,
        ticketApi: ticketAPI,
        contactsApi: contactsAPI,
        donationApi: donationAPI,
    }
}

// Default export for direct api access (used by stores)
export default api
