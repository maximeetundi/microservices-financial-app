/**
 * Safe number formatting utility
 * Handles null/undefined values gracefully
 */
export const safeNumber = (value: any, fallback: number = 0): number => {
    if (value == null || value === '' || isNaN(value)) return fallback
    return Number(value)
}

export const formatNumber = (value: any, fallback: string = '0'): string => {
    const num = safeNumber(value)
    return num.toLocaleString()
}

export const formatPrice = (value: any, fallback: string = '0'): string => {
    const num = safeNumber(value)
    if (num >= 1000) return num.toLocaleString()
    if (num >= 1) return num.toFixed(2)
    if (num > 0) return num.toFixed(6)
    return fallback
}

export const formatPercent = (value: any, fallback: string = '0.00'): string => {
    const num = safeNumber(value)
    return num.toFixed(2)
}

export const formatMoney = (value: any, currency: string = 'USD'): string => {
    const num = safeNumber(value)
    try {
        return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(num)
    } catch {
        return `${num.toLocaleString()} ${currency}`
    }
}

export const useFormatters = () => {
    return {
        safeNumber,
        formatNumber,
        formatPrice,
        formatPercent,
        formatMoney
    }
}
