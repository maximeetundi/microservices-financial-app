export default defineNuxtRouteMiddleware((to, from) => {
    // Only run on client side
    if (process.server) return

    const adminToken = localStorage.getItem('adminToken')

    // If no token and trying to access protected admin pages
    if (!adminToken && to.path.startsWith('/admin') && to.path !== '/admin/login') {
        return navigateTo('/admin/login')
    }

    // If has token and trying to access login page, redirect to dashboard
    if (adminToken && to.path === '/admin/login') {
        return navigateTo('/admin/dashboard')
    }
})
