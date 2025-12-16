import axios from 'axios'

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
                    return api(originalRequest)
                } catch {
                    localStorage.removeItem('accessToken')
                    localStorage.removeItem('refreshToken')
                    if (typeof window !== 'undefined') {
                        window.location.href = '/auth/login'
                    }
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
    changePassword: (oldPassword: string, newPassword: string) =>
        api.post('/auth-service/api/v1/users/change-password', { old_password: oldPassword, new_password: newPassword }),
    setup2FA: () => api.post('/auth-service/api/v1/users/2fa/setup'),
    verify2FA: (code: string) => api.post('/auth-service/api/v1/users/2fa/verify', { code }),
    disable2FA: (code: string) => api.post('/auth-service/api/v1/users/2fa/disable', { code }),
}

// ========== Wallets ==========
export const walletAPI = {
    getAll: () => api.get('/wallet-service/api/v1/wallets'),
    getWallets: () => api.get('/wallet-service/api/v1/wallets'),
    get: (id: string) => api.get(`/wallet-service/api/v1/wallets/${id}`),
    create: (data: { currency: string, wallet_type: string, name?: string, description?: string }) =>
        api.post('/wallet-service/api/v1/wallets/', data),
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
export const exchangeAPI = {
    getRates: (baseCurrency?: string) =>
        api.get(`/exchange-service/api/v1/exchange/rates${baseCurrency ? `?base=${baseCurrency}` : ''}`),
    getRate: (from: string, to: string) => api.get(`/exchange-service/api/v1/exchange/rate?from=${from}&to=${to}`),
    convert: (fromCurrency: string, toCurrency: string, amount: number) =>
        api.post('/exchange-service/api/v1/exchange/convert', {
            from_currency: fromCurrency,
            to_currency: toCurrency,
            amount,
        }),
    getCryptoRates: () => api.get('/exchange-service/api/v1/exchange/crypto/rates'),
    buyCrypto: (currency: string, amount: number, paymentMethod: string) =>
        api.post('/exchange-service/api/v1/exchange/crypto/buy', { currency, amount, payment_method: paymentMethod }),
    sellCrypto: (currency: string, amount: number, destinationWalletId: string) =>
        api.post('/exchange-service/api/v1/exchange/crypto/sell', {
            currency,
            amount,
            destination_wallet_id: destinationWalletId,
        }),
    getHistory: (limit = 50) => api.get(`/exchange-service/api/v1/exchange/history?limit=${limit}`),
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
    }
}

export default api

