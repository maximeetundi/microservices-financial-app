import { defineStore } from 'pinia'
import { walletAPI } from '~/composables/useApi'
import { useExchangeStore } from './exchange'

interface Wallet {
    id: string
    currency: string
    balance: string
    type: string
    is_default: boolean
    usd_rate?: number
    balanceUSD?: number
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
        // Helper to get rates from exchange store
        getRate(currency: string): number {
            const exchangeStore = useExchangeStore()
            return exchangeStore.getRate(currency, 'USD')
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
                // Fetch wallets and rates in parallel
                const exchangeStore = useExchangeStore()
                const [walletResponse, _] = await Promise.all([
                    walletAPI.getAll(),
                    exchangeStore.fetchRates()
                ])

                if (walletResponse.data && walletResponse.data.wallets) {
                    this.wallets = walletResponse.data.wallets.map((w: any) => ({
                        ...w,
                        type: w.wallet_type || w.type || 'fiat',
                        wallet_type: w.wallet_type || w.type || 'fiat'
                    }))
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
            const exchangeStore = useExchangeStore()

            this.wallets.forEach(w => {
                const balance = Number(w.balance) || 0
                // Use usd_rate from backend if available, else get from exchange store
                const rate = w.usd_rate || exchangeStore.getRate(w.currency, 'USD')
                const inUSD = balance * rate

                w.balanceUSD = inUSD

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
