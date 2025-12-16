export default defineNuxtRouteMiddleware((to, from) => {
    // Skip for auth pages
    if (to.path.startsWith('/auth')) return

    // Use useCookie for SSR compatibility
    const token = useCookie('auth_token')
    const accessToken = useCookie('accessToken')

    // Check authentication - check both cookie names for compatibility
    if (!token.value && !accessToken.value) {
        // On client side, also check localStorage as fallback
        if (process.client) {
            const localToken = localStorage.getItem('accessToken')
            if (localToken) {
                // Sync localStorage to cookie
                accessToken.value = localToken
                return
            }
        }
        return navigateTo('/auth/login')
    }
})
