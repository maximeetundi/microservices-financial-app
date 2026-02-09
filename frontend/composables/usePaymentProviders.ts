/**
 * Composable for loading payment providers/aggregators
 * Filters by user country and current IP location
 */
import { ref, computed } from 'vue'

interface PaymentProvider {
    id: string
    name: string
    display_name: string
    provider_type: string
    logo_url?: string
    is_demo_mode: boolean
    fee_percentage: number
    fee_fixed: number
    min_amount: number
    max_amount: number
    // Derived
    displayLabel: string
    icon: string
    color: string
    category: 'mobile_money' | 'card' | 'bank' | 'international' | 'demo'
}

// Mapping provider codes to simple display names
const PROVIDER_DISPLAY_NAMES: Record<string, { label: string; icon: string; color: string; category: PaymentProvider['category'] }> = {
    // Mobile Money Africa
    'orange_money': { label: 'Orange Money', icon: '🟠', color: 'orange', category: 'mobile_money' },
    'orange_money_ci': { label: 'Orange Money', icon: '🟠', color: 'orange', category: 'mobile_money' },
    'orange_money_cm': { label: 'Orange Money', icon: '🟠', color: 'orange', category: 'mobile_money' },
    'orange_money_sn': { label: 'Orange Money', icon: '🟠', color: 'orange', category: 'mobile_money' },

    'mtn_momo': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_ci': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_cm': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_sn': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_bj': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_gh': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },
    'mtn_za': { label: 'MTN MoMo', icon: '🟡', color: 'yellow', category: 'mobile_money' },

    'wave': { label: 'Wave', icon: '🌊', color: 'blue', category: 'mobile_money' },
    'wave_ci': { label: 'Wave', icon: '🌊', color: 'blue', category: 'mobile_money' },
    'wave_sn': { label: 'Wave', icon: '🌊', color: 'blue', category: 'mobile_money' },

    'moov': { label: 'Moov Money', icon: '🔵', color: 'blue', category: 'mobile_money' },
    'moov_ci': { label: 'Moov Money', icon: '🔵', color: 'blue', category: 'mobile_money' },
    'moov_bj': { label: 'Moov Money', icon: '🔵', color: 'blue', category: 'mobile_money' },
    'moov_tg': { label: 'Moov Money', icon: '🔵', color: 'blue', category: 'mobile_money' },

    // Aggregators (Tiers)
    'cinetpay': { label: 'CinetPay', icon: '💳', color: 'green', category: 'mobile_money' },
    'flutterwave': { label: 'Flutterwave', icon: '🦋', color: 'orange', category: 'international' },
    'lygos': { label: 'Lygos', icon: '🦁', color: 'purple', category: 'mobile_money' },
    'fedapay': { label: 'FedaPay', icon: '🔄', color: 'blue', category: 'mobile_money' },
    'yellowcard': { label: 'Yellow Card', icon: '💛', color: 'yellow', category: 'mobile_money' },

    // International
    'paypal': { label: 'PayPal', icon: '🅿️', color: 'blue', category: 'international' },
    'stripe': { label: 'Carte Bancaire (Stripe)', icon: '💳', color: 'purple', category: 'card' },

    // Bank
    'bank_transfer': { label: 'Virement Bancaire', icon: '🏦', color: 'emerald', category: 'bank' },

    // Demo
    'demo': { label: 'Mode Test', icon: '🧪', color: 'gray', category: 'demo' },
}

export function usePaymentProviders() {
    const providers = ref<PaymentProvider[]>([])
    const loading = ref(false)
    const error = ref('')
    const userCountry = ref('') // Will be set from user profile
    const ipCountry = ref('')

    // Detect country from IP
    const detectIpCountry = async () => {
        try {
            const res = await fetch('https://ip-api.com/json/?fields=countryCode')
            const data = await res.json()
            if (data.countryCode) {
                ipCountry.value = data.countryCode
            }
        } catch (e) {
            console.warn('Could not detect IP country:', e)
        }
    }

    // Get user profile country
    const getUserCountry = async () => {
        try {
            const config = useRuntimeConfig()
            const API_URL = config.public.apiBaseUrl || 'http://localhost:8000'
            const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : null

            if (!token) return

            const res = await fetch(`${API_URL}/auth-service/api/v1/users/profile`, {
                headers: { 'Authorization': `Bearer ${token}` }
            })

            if (res.ok) {
                const data = await res.json()
                if (data.user?.country || data.country) {
                    userCountry.value = data.user?.country || data.country
                }
            }
        } catch (e) {
            console.warn('Could not get user profile country:', e)
        }
    }

    // Load available payment methods for multiple countries
    const loadProviders = async (countries?: string[]) => {
        loading.value = true
        error.value = ''

        try {
            const config = useRuntimeConfig()
            const API_URL = config.public.apiBaseUrl || 'http://localhost:8000'

			const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : null

            // Build countries list from params or auto-detect
            let countryList = countries || []

            if (countryList.length === 0) {
                // Try to get both profile country and IP country
                if (!userCountry.value) await getUserCountry()
                if (!ipCountry.value) await detectIpCountry()

                // Add user profile country (priority)
                if (userCountry.value) countryList.push(userCountry.value)

                // Add IP country if different
                if (ipCountry.value && ipCountry.value !== userCountry.value) {
                    countryList.push(ipCountry.value)
                }

                // Fallback to default
                if (countryList.length === 0) countryList.push('CI')
            }

            // Build URL with multiple country params
            const params = new URLSearchParams()
            countryList.forEach(c => params.append('country', c))

            // Use Transfer Service Aggregator Proxy instead of direct Admin Service call
            const url = `${API_URL}/transfer-service/api/v1/aggregators?${params.toString()}`
            console.log('[Debug] Fetching Payment Methods (Transfer Service):', url)
            const response = await fetch(url, {
				headers: token ? { 'Authorization': `Bearer ${token}` } : undefined
			})

            if (!response.ok) {
                throw new Error('Failed to load payment methods')
            }

            const data = await response.json()
            const methods = data.aggregators || []

            // Transform to our format
            providers.value = methods.map((m: any) => {
                // Transfer Service returns 'code' (e.g. mtn_momo) and 'name' (e.g. MTN Mobile Money)
                const code = m.code || m.name
                const mapping = PROVIDER_DISPLAY_NAMES[code] || PROVIDER_DISPLAY_NAMES[code.toLowerCase()] || {
                    label: m.name || code,
                    icon: '💰',
                    color: 'gray',
                    category: 'mobile_money' as const
                }

                return {
                    id: m.code || m.id,
                    name: code,
                    display_name: m.name || m.display_name,
                    provider_type: m.provider_type || 'aggregator',
                    logo_url: m.logo_url,
                    is_demo_mode: m.is_demo_mode || false,
                    fee_percentage: m.fee_percent || m.fee_percentage || 0,
                    fee_fixed: m.fee_fixed || 0,
                    min_amount: m.min_amount || 0,
                    max_amount: m.max_amount || 1000000,
                    displayLabel: mapping.label,
                    icon: mapping.icon,
                    color: mapping.color,
                    category: m.is_demo_mode ? 'demo' : mapping.category,
                }
            })

            // Always add demo if not present and needed
            const hasDemo = providers.value.some(p => p.is_demo_mode || p.name === 'demo')
            if (!hasDemo) {
                // We'll check from another endpoint or add manually
            }

        } catch (e: any) {
            console.error('Error loading payment providers:', e)
            error.value = e.message || 'Erreur de chargement'

            // Fallback to demo only
            providers.value = [{
                id: 'demo',
                name: 'demo',
                display_name: 'Demo Provider',
                provider_type: 'demo',
                is_demo_mode: true,
                fee_percentage: 0,
                fee_fixed: 0,
                min_amount: 100,
                max_amount: 1000000,
                displayLabel: 'Mode Test',
                icon: '🧪',
                color: 'gray',
                category: 'demo',
            }]
        } finally {
            loading.value = false
        }
    }

    // Group providers by category
    const groupedProviders = computed(() => {
        const groups: Record<string, { title: string; providers: PaymentProvider[] }> = {
            mobile_money: { title: '📱 Mobile Money', providers: [] },
            card: { title: '💳 Carte Bancaire', providers: [] },
            bank: { title: '🏦 Virement', providers: [] },
            international: { title: '🌍 International', providers: [] },
            demo: { title: '🧪 Test', providers: [] },
        }

        for (const p of providers.value) {
            if (groups[p.category]) {
                groups[p.category].providers.push(p)
            }
        }

        // Filter out empty groups
        return Object.entries(groups)
            .filter(([_, g]) => g.providers.length > 0)
            .map(([key, g]) => ({ key, ...g }))
    })

    // Get color class for a provider
    const getColorClass = (color: string, type: 'border' | 'bg' | 'text' = 'border') => {
        const map: Record<string, Record<string, string>> = {
            orange: { border: 'border-orange-500', bg: 'bg-orange-50 dark:bg-orange-500/10', text: 'text-orange-500' },
            yellow: { border: 'border-yellow-500', bg: 'bg-yellow-50 dark:bg-yellow-500/10', text: 'text-yellow-500' },
            blue: { border: 'border-blue-500', bg: 'bg-blue-50 dark:bg-blue-500/10', text: 'text-blue-500' },
            green: { border: 'border-emerald-500', bg: 'bg-emerald-50 dark:bg-emerald-500/10', text: 'text-emerald-500' },
            emerald: { border: 'border-emerald-500', bg: 'bg-emerald-50 dark:bg-emerald-500/10', text: 'text-emerald-500' },
            purple: { border: 'border-purple-500', bg: 'bg-purple-50 dark:bg-purple-500/10', text: 'text-purple-500' },
            gray: { border: 'border-gray-500', bg: 'bg-gray-50 dark:bg-gray-500/10', text: 'text-gray-500' },
        }
        return map[color]?.[type] || map.gray[type]
    }

    return {
        providers,
        groupedProviders,
        loading,
        error,
        userCountry,
        ipCountry,
        detectIpCountry,
        getUserCountry,
        loadProviders,
        getColorClass,
    }
}
