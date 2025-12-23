export default defineNuxtRouteMiddleware((to, from) => {
    // Skip for auth pages
    if (to.path.startsWith('/auth')) return

    // Use useCookie for SSR compatibility - use only 'accessToken' for consistency
    const accessToken = useCookie('accessToken')
    const refreshToken = useCookie('refreshToken')

    // On client side, sync localStorage to cookies first
    if (process.client) {
        const localAccessToken = localStorage.getItem('accessToken')
        const localRefreshToken = localStorage.getItem('refreshToken')

        // If we have tokens in localStorage but not in cookies, sync them
        if (localAccessToken && !accessToken.value) {
            accessToken.value = localAccessToken
        }
        if (localRefreshToken && !refreshToken.value) {
            refreshToken.value = localRefreshToken
        }

        // Check if user has valid tokens (either in cookie or localStorage)
        if (localAccessToken || accessToken.value) {
            return // User is authenticated
        }
    } else {
        // Server-side: only check cookies
        if (accessToken.value) {
            return // User is authenticated
        }
    }

    // No valid authentication found, redirect to login
    return navigateTo('/auth/login')
})
