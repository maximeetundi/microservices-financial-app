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

// Auth-service admin endpoints (via Kong)
const authApi = axios.create({
    baseURL: `${API_URL}/auth-service/api/v1/admin`,
    headers: {
        'Content-Type': 'application/json',
    },
});
authApi.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

// Unlock user PIN (reset failed attempts)
export const unlockUserPin = (userId: string) =>
    authApi.post(`/users/${userId}/unlock-pin`);

// KYC
export const approveKYC = (userId: string, level: string) =>
    api.post(`/kyc/${userId}/approve`, { level });
export const rejectKYC = (userId: string, reason: string) =>
    api.post(`/kyc/${userId}/reject`, { reason });
export const getUserKYCDocuments = (userId: string) =>
    api.get(`/users/${userId}/kyc/documents`);
// Get secure presigned URL for viewing/downloading KYC document
export const getKYCDocumentURL = (filePath: string) =>
    api.post('/kyc/document-url', { file_path: filePath });
export const getKYCDownloadURL = (filePath: string, downloadName?: string) =>
    api.post('/kyc/download-url', { file_path: filePath, download_name: downloadName });

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

// Fees & Settings
// Fees
export const getFeeConfigs = () => api.get('/fees');
export const updateFeeConfig = (key: string, data: any) => api.put(`/fees/${key}`, data);
export const createFeeConfig = (data: any) => api.post('/fees', data);

// System Settings
export const getSystemSettings = () => api.get('/settings');
export const updateSystemSetting = (data: any) => api.put('/settings', data);

// Platform Accounts & Wallets
export const getPlatformAccounts = () => api.get('/platform/accounts');
export const createPlatformAccount = (data: any) => api.post('/platform/accounts', data);
export const creditPlatformAccount = (id: string, data: any) => api.post(`/platform/accounts/${id}/credit`, { ...data, account_id: id });
export const debitPlatformAccount = (id: string, data: any) => api.post(`/platform/accounts/${id}/debit`, { ...data, account_id: id });

export const getPlatformCryptoWallets = () => api.get('/platform/crypto-wallets');
export const createPlatformCryptoWallet = (data: any) => api.post('/platform/crypto-wallets', data);
export const syncPlatformCryptoWallet = (id: string) => api.put(`/platform/crypto-wallets/${id}/sync`);
export const consolidateUserFunds = (data: { target_type: string, amount: number, currency: string }) => api.post('/platform/crypto-wallets/consolidate', data);
export const transferPlatformFunds = (data: { source_id: string, target_id: string, amount: number, description?: string }) => api.post('/platform/crypto-wallets/transfer', data);

export const getPlatformTransactions = (limit = 50, offset = 0) => api.get(`/platform/transactions?limit=${limit}&offset=${offset}`);
export const getPlatformWalletTransactions = (walletId: string, limit = 50, offset = 0) => api.get(`/platform/crypto-wallets/${walletId}/transactions?limit=${limit}&offset=${offset}`);
export const getPlatformReconciliation = () => api.get('/platform/reconciliation');



// Support Tickets - route through the same API gateway as admin
// The support-service is accessible via Kong at /support-service
// In production, use the main API gateway (same as admin but replace admin-service path)
const getBaseApiUrl = () => {
    // For production, use the same base URL as admin but without the /api/v1/admin path
    // API_URL is like https://api.admin.tech-afm.com
    // We need to go through Kong which routes /support-service to the support service
    if (API_URL.includes('localhost')) {
        // Local dev: support-service is on port 8089
        return `http://localhost:8080/support-service/api/v1`;
    }
    // Production: use the app API gateway (not admin gateway)
    // Replace api.admin with api.app to go through the main Kong gateway
    const appApiUrl = API_URL.replace('api.admin', 'api.app').replace('/api/v1/admin', '');
    return `${appApiUrl}/support-service/api/v1`;
};

const supportApi = axios.create({
    baseURL: getBaseApiUrl(),
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

// Admin support endpoints - using /admin/support for admin access
export const getSupportTickets = (limit = 50, offset = 0, status?: string) =>
    supportApi.get(`/admin/support/conversations?limit=${limit}&offset=${offset}${status ? `&status=${status}` : ''}`);

export const getSupportTicket = (id: string) =>
    supportApi.get(`/admin/support/conversations/${id}`);

export const getTicketMessages = (conversationId: string) =>
    supportApi.get(`/admin/support/conversations/${conversationId}/messages`);

export const assignAgent = (conversationId: string, agentId: string) =>
    supportApi.put(`/admin/support/conversations/${conversationId}/assign`, { agent_id: agentId });

export const closeTicket = (conversationId: string, rating?: number, feedback?: string) =>
    supportApi.put(`/admin/support/conversations/${conversationId}/close`, { rating, feedback });

export const getSupportStats = () =>
    supportApi.get('/admin/support/stats');

export const getSupportAgents = () =>
    supportApi.get('/admin/support/agents');

export const uploadSupportFile = (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return supportApi.post('/admin/support/upload', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const sendTicketMessage = (conversationId: string, message: string, attachments?: string[]) =>
    supportApi.post(`/admin/support/conversations/${conversationId}/messages`, {
        content: message,
        content_type: 'text',
        attachments: attachments || []
    });

// ========== Ticket/Events API ==========
const ticketApi = axios.create({
    baseURL: (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088').replace('/api/v1/admin', ''),
    headers: { 'Content-Type': 'application/json' },
});

ticketApi.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

// Get all events (admin view)
export const getAllEvents = (limit = 50, offset = 0) =>
    ticketApi.get(`/ticket-service/api/v1/events/active?limit=${limit}&offset=${offset}`);

// Get event details
export const getEventDetails = (eventId: string) =>
    ticketApi.get(`/ticket-service/api/v1/events/${eventId}`);

// Get event statistics
export const getEventStats = (eventId: string) =>
    ticketApi.get(`/ticket-service/api/v1/events/${eventId}/stats`);

// Get event tickets
export const getEventTickets = (eventId: string, limit = 50, offset = 0) =>
    ticketApi.get(`/ticket-service/api/v1/events/${eventId}/tickets?limit=${limit}&offset=${offset}`);

// Verify ticket
export const verifyTicket = (code: string) =>
    ticketApi.post('/ticket-service/api/v1/tickets/verify', { ticket_code: code });

// Mark ticket as used
export const useTicket = (ticketId: string) =>
    ticketApi.post(`/ticket-service/api/v1/tickets/${ticketId}/use`);

// Admin: Update event status (suspend, cancel, activate)
export const updateEventStatus = (eventId: string, status: string) =>
    ticketApi.put(`/ticket-service/api/v1/events/${eventId}`, { status });

// Admin: Suspend event (mark as cancelled for scam/fraud)
export const suspendEvent = (eventId: string) =>
    ticketApi.put(`/ticket-service/api/v1/events/${eventId}`, { status: 'cancelled' });

// Admin: Delete event entirely
export const deleteEvent = (eventId: string) =>
    ticketApi.delete(`/ticket-service/api/v1/events/${eventId}`);

// ========== Associations API ==========
const associationApi = axios.create({
    baseURL: API_URL,
    headers: { 'Content-Type': 'application/json' },
});

associationApi.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('admin_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

// Get all associations
export const getAllAssociations = (limit = 50, offset = 0) =>
    associationApi.get(`/association-service/api/v1/associations?limit=${limit}&offset=${offset}`);

// Get association details
export const getAssociationDetails = (id: string) =>
    associationApi.get(`/association-service/api/v1/associations/${id}`);

// Get association members
export const getAssociationMembers = (id: string) =>
    associationApi.get(`/association-service/api/v1/associations/${id}/members`);

// Get association treasury
export const getAssociationTreasury = (id: string) =>
    associationApi.get(`/association-service/api/v1/associations/${id}/treasury`);

// Suspend association
export const suspendAssociation = (id: string) =>
    associationApi.put(`/association-service/api/v1/associations/${id}`, { status: 'suspended' });

// Reactivate association
export const activateAssociation = (id: string) =>
    associationApi.put(`/association-service/api/v1/associations/${id}`, { status: 'active' });

export default api;


