import { defineStore } from 'pinia'
import { exchangeAPI, fiatAPI } from '~/composables/useApi'

interface CryptoRate {
    pair: string
    symbol?: string // Optional, as some endpoints might return this
    from_currency?: string
    to_currency?: string
    rate: number
    change_24h: number
    volume_24h: number
    last_updated: string
    // Add other fields if necessary
    price?: number
}

interface FiatRate {
    rate: number
    bid: number
    ask: number
    change_24h: number
    last_update: number
}

interface ExchangeState {
    cryptoRates: CryptoRate[]
    fiatRates: Record<string, FiatRate>
    loading: boolean
    error: string | null
    lastUpdated: number | null
}

export const useExchangeStore = defineStore('exchange', {
    state: (): ExchangeState => ({
        cryptoRates: [],
        fiatRates: {},
        loading: false,
        error: null,
        lastUpdated: null
    }),

    getters: {
        getRate: (state) => (from: string, to: string = 'USD'): number => {
            from = from.toUpperCase()
            to = to.toUpperCase()

            if (from === to) return 1

            // Helper to get rate against USD
            const getPriceInUSD = (currency: string): number => {
                if (currency === 'USD') return 1

                // Check fiat rates (assuming base is USD)
                if (state.fiatRates[currency]) {
                    // If base is USD, the rate in fiatRates (e.g. EUR) is USD -> EUR?
                    // Let's check the handler: h.rateService.GetRate(baseCurrency, targetCurrency)
                    // If base=USD, GetFiatRates returns rates for EUR, GBP, etc.
                    // usually this means 1 USD = X EUR. 
                    // So Price of 1 EUR in USD is 1 / X.
                    return 1 / state.fiatRates[currency].rate
                }

                // Check crypto rates
                const cryptoPair = state.cryptoRates.find(r => r.pair === `${currency}/USD` || r.pair === `${currency}USD` || r.pair === `${currency}-USD`)
                if (cryptoPair) return cryptoPair.rate

                // Check reverse crypto pair (USD/BTC - unlikely but possible)
                const reverseCryptoPair = state.cryptoRates.find(r => r.pair === `USD/${currency}` || r.pair === `USD${currency}`)
                if (reverseCryptoPair) return 1 / reverseCryptoPair.rate

                return 0
            }

            const fromRate = getPriceInUSD(from)
            const toRate = getPriceInUSD(to)

            if (fromRate && toRate) {
                return fromRate / toRate
            }

            // Fallback for hardcoded common values if store is empty (prevent NaN on init)
            // This mirrors the original wallet.ts logic but should be temporary
            const fallbackRates: Record<string, number> = {
                'USD': 1,
                'EUR': 1.08,
                'GBP': 1.27,
                'XOF': 0.00167,
                'XAF': 0.00167, // 1/600 approx
                'BTC': 43000,
                'ETH': 2200,
                'SOL': 95,
            }
            const f = fallbackRates[from]
            const t = fallbackRates[to]
            if (f && t) return f / t

            return 0
        }
    },

    actions: {
        async fetchRates() {
            this.loading = true
            this.error = null
            try {
                // Fetch both in parallel
                const [cryptoRes, fiatRes] = await Promise.allSettled([
                    exchangeAPI.getRates(),
                    fiatAPI.getRates('USD')
                ])

                if (cryptoRes.status === 'fulfilled' && cryptoRes.value.data?.rates) {
                    // Normalize rates to ensure 'pair' exists
                    this.cryptoRates = cryptoRes.value.data.rates.map((r: any) => {
                        let pair = r.pair || r.symbol
                        if (!pair && r.from_currency && r.to_currency) {
                            pair = `${r.from_currency}/${r.to_currency}`
                        }

                        return {
                            ...r,
                            pair: pair || '', // Ensure pair is never undefined
                            symbol: pair || '', // Populate symbol too for convenience
                        }
                    }).filter((r: any) => r.pair)
                }

                if (fiatRes.status === 'fulfilled' && fiatRes.value.data?.rates) {
                    // Check structure based on handler: rates[targetCurrency] = { rate: ... }
                    // If fiatAPI.getRates('USD') returns rates: { EUR: { rate: 0.9, ... } }
                    // This means 1 USD = 0.9 EUR
                    this.fiatRates = fiatRes.value.data.rates
                }

                this.lastUpdated = Date.now()
            } catch (err: any) {
                console.error("Failed to fetch rates:", err)
                this.error = err.message || "Failed to fetch rates"
            } finally {
                this.loading = false
            }
        }
    }
})
