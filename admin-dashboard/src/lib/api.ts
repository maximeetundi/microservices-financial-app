import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';

const api = axios.create({
    baseURL: `${API_URL}/api/v1/admin`,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor for adding auth token
api.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

// Response interceptor for handling errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            if (typeof window !== 'undefined') {
                localStorage.removeItem('admin_token');
                window.location.href = '/login';
            }
        }
        return Promise.reject(error);
    }
);

// Auth
export const login = async (email: string, password: string) => {
    const response = await api.post('/login', { email, password });
    const { token, admin } = response.data;
    if (typeof window !== 'undefined') {
        localStorage.setItem('admin_token', token);
        localStorage.setItem('admin_user', JSON.stringify(admin));
    }
    return response.data;
};

export const logout = () => {
    if (typeof window !== 'undefined') {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
        window.location.href = '/login';
    }
};

export const getCurrentAdmin = () => api.get('/me');

// Dashboard
export const getDashboard = () => api.get('/dashboard');

// Users
export const getUsers = (limit = 50, offset = 0) =>
    api.get(`/users?limit=${limit}&offset=${offset}`);
export const blockUser = (id: string, reason: string) =>
    api.post(`/users/${id}/block`, { reason });
export const unblockUser = (id: string) =>
    api.post(`/users/${id}/unblock`);

// KYC
export const approveKYC = (userId: string, level: string) =>
    api.post(`/kyc/${userId}/approve`, { level });
export const rejectKYC = (userId: string, reason: string) =>
    api.post(`/kyc/${userId}/reject`, { reason });

// Transactions
export const getTransactions = (limit = 50, offset = 0) =>
    api.get(`/transactions?limit=${limit}&offset=${offset}`);
export const blockTransaction = (id: string, reason: string) =>
    api.post(`/transactions/${id}/block`, { reason });
export const refundTransaction = (id: string, reason: string) =>
    api.post(`/transactions/${id}/refund`, { reason });

// Cards
export const getCards = (limit = 50, offset = 0) =>
    api.get(`/cards?limit=${limit}&offset=${offset}`);
export const freezeCard = (id: string, reason: string) =>
    api.post(`/cards/${id}/freeze`, { reason });
export const blockCard = (id: string, reason: string) =>
    api.post(`/cards/${id}/block`, { reason });

// Wallets
export const getWallets = (limit = 50, offset = 0) =>
    api.get(`/wallets?limit=${limit}&offset=${offset}`);
export const freezeWallet = (id: string, reason: string) =>
    api.post(`/wallets/${id}/freeze`, { reason });

// Admins
export const getAdmins = () => api.get('/admins');
export const createAdmin = (data: any) => api.post('/admins', data);
export const updateAdmin = (id: string, data: any) => api.put(`/admins/${id}`, data);
export const deleteAdmin = (id: string) => api.delete(`/admins/${id}`);

// Roles
export const getRoles = () => api.get('/roles');

// Audit Logs
export const getAuditLogs = (limit = 100, offset = 0) =>
    api.get(`/logs?limit=${limit}&offset=${offset}`);

// Support Tickets (via support-service proxied through Kong)
const supportApi = axios.create({
    baseURL: `${API_URL}/support-service/api/v1`,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Use same interceptors for support API
supportApi.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

export const getSupportTickets = (limit = 50, offset = 0) =>
    supportApi.get(`/tickets?limit=${limit}&offset=${offset}`);

export const getSupportTicket = (id: string) =>
    supportApi.get(`/tickets/${id}`);

export const getTicketMessages = (ticketId: string) =>
    supportApi.get(`/tickets/${ticketId}/messages`);

export const sendTicketMessage = (ticketId: string, message: string) =>
    supportApi.post(`/tickets/${ticketId}/messages`, { message });

export const closeTicket = (ticketId: string) =>
    supportApi.post(`/tickets/${ticketId}/close`);

export const getSupportStats = () =>
    supportApi.get('/stats');

export default api;
