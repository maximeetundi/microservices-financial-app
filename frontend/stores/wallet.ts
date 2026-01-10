import { defineStore } from 'pinia'
import { walletAPI } from '~/composables/useApi'

interface Wallet {
    id: string
    currency: string
    balance: string
    type: string
    is_default: boolean
    usd_rate?: number
}

interface WalletState {
    wallets: Wallet[]
    totalBalance: number
    cryptoBalance: number
    loading: boolean
    error: string | null
    lastUpdated: number | null
}

export const useWalletStore = defineStore('wallet', {
    state: (): WalletState => ({
        wallets: [],
        totalBalance: 0,
        cryptoBalance: 0,
        loading: false,
        error: null,
        lastUpdated: null
    }),

    actions: {
        // Helper to estimate rates if not provided by backend
        // In a real app, these should come from a store or API
        getRate(currency: string): number {
            const rates: Record<string, number> = {
                'USD': 1,
                'EUR': 1.08,
                'GBP': 1.27,
                'XOF': 0.00167,
                'XAF': 0.00167,
                'BTC': 43000,
                'ETH': 2200,
            }
            return rates[currency?.toUpperCase()] || 1
        },

        initialize() {
            if (typeof window !== 'undefined') {
                const cached = localStorage.getItem('wallet_store')
                if (cached) {
                    try {
                        const parsed = JSON.parse(cached)
                        this.wallets = parsed.wallets || []
                        this.totalBalance = parsed.totalBalance || 0
                        this.cryptoBalance = parsed.cryptoBalance || 0
                        this.lastUpdated = parsed.lastUpdated
                    } catch (e) {
                        console.error('Failed to parse cached wallet store', e)
                    }
                }
            }
        },

        async fetchWallets() {
            // If we have recent data (e.g. < 1 min), don't show loading spinner but fetch in background
            // If no data, show loading
            if (this.wallets.length === 0) {
                this.loading = true
            }

            try {
                const response = await walletAPI.getAll()
                if (response.data && response.data.wallets) {
                    this.wallets = response.data.wallets
                    this.calculateBalances()
                    this.lastUpdated = Date.now()
                    this.error = null

                    this.persist()
                }
            } catch (err: any) {
                this.error = err.message || 'Failed to fetch wallets'
            } finally {
                this.loading = false
            }
        },

        calculateBalances() {
            let total = 0
            let crypto = 0

            this.wallets.forEach(w => {
                const balance = Number(w.balance) || 0
                // Use usd_rate from backend if available, else fallback
                const rate = w.usd_rate || this.getRate(w.currency)
                const inUSD = balance * rate

                total += inUSD

                if (['BTC', 'ETH', 'USDT', 'USDC', 'BNB', 'XRP', 'SOL', 'CRYPTO'].includes(w.currency?.toUpperCase()) || w.type === 'crypto') {
                    crypto += inUSD
                }
            })

            this.totalBalance = Math.round(total * 100) / 100
            this.cryptoBalance = Math.round(crypto * 100) / 100
        },

        persist() {
            if (typeof window !== 'undefined') {
                localStorage.setItem('wallet_store', JSON.stringify({
                    wallets: this.wallets,
                    totalBalance: this.totalBalance,
                    cryptoBalance: this.cryptoBalance,
                    lastUpdated: this.lastUpdated
                }))
            }
        }
    }
})
