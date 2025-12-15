export default defineNuxtRouteMiddleware((to, from) => {
    // Skip middleware on server
    if (process.server) return

    // Skip for auth pages
    if (to.path.startsWith('/auth')) return

    // Check authentication
    const token = localStorage.getItem('accessToken')

    if (!token) {
        return navigateTo('/auth/login')
    }
})
