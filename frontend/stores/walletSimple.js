// Simple wallet store without TypeScript complications
import { defineStore } from 'pinia'
import { walletAPI, systemConfigAPI } from '~/composables/useApi'
import { useExchangeStore } from './exchange'
import { usePin } from '~/composables/usePinSimple'

export const useWalletStore = defineStore('wallet', {
    state: () => ({
        wallets: [],
        totalBalance: 0,
        cryptoBalance: 0,
        loading: false,
        error: null,
        lastUpdated: null,
        testnetEnabled: false,
        pinVerified: false
    }),

    actions: {
        // Helper to get rates from exchange store
        getRate(currency) {
            const exchangeStore = useExchangeStore()
            return exchangeStore.getRate(currency, 'USD')
        },

        // SECURE: Verify balance with PIN locally before sensitive operations
        async verifyBalanceWithPin(pin) {
            const { verifyPin } = usePin()
            
            // First verify PIN locally (never sends PIN to backend)
            const result = await verifyPin(pin)
            if (!result.valid) {
                return { success: false, message: result.message || 'PIN incorrect' }
            }

            // PIN is correct, allow access to balance
            this.pinVerified = true
            
            // Set timeout to reset verification (5 minutes)
            setTimeout(() => {
                this.pinVerified = false
            }, 5 * 60 * 1000)

            return { 
                success: true, 
                message: 'Balance vérifié avec succès',
                balance: this.totalBalance 
            }
        },

        // Check if PIN verification is required
        isPinVerificationRequired() {
            return !this.pinVerified
        },

        // Reset PIN verification state
        resetPinVerification() {
            this.pinVerified = false
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
                        this.testnetEnabled = parsed.testnetEnabled || false
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
                const [walletResponse, configResponse, _] = await Promise.all([
                    walletAPI.getAll(),
                    systemConfigAPI.getPublicConfig(),
                    exchangeStore.fetchRates()
                ])

                if (configResponse.data) {
                    this.testnetEnabled = configResponse.data.testnet_enabled
                }

                if (walletResponse.data && walletResponse.data.wallets) {
                    this.wallets = walletResponse.data.wallets.map((w) => ({
                        ...w,
                        type: w.wallet_type || w.type || 'fiat',
                        wallet_type: w.wallet_type || w.type || 'fiat'
                    }))
                    this.calculateBalances()
                    this.lastUpdated = Date.now()
                    this.error = null

                    this.persist()
                }
            } catch (err) {
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
                    lastUpdated: this.lastUpdated,
                    testnetEnabled: this.testnetEnabled
                }))
            }
        }
    }
})
