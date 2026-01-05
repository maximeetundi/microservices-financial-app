import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const API_URL = 'https://api.app.maximeetundi.store';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor
api.interceptors.request.use(async (config) => {
    const token = await AsyncStorage.getItem('accessToken');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Response interceptor
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 401) {
            await AsyncStorage.removeItem('accessToken');
            await AsyncStorage.removeItem('refreshToken');
            // Navigate to login
        }
        return Promise.reject(error);
    }
);

// Association API
export const associationAPI = {
    getAll: () => api.get('/association-service/api/v1/associations'),
    get: (id: string) => api.get(`/association-service/api/v1/associations/${id}`),
    create: (data: any) => api.post('/association-service/api/v1/associations', data),
    getMembers: (id: string) => api.get(`/association-service/api/v1/associations/${id}/members`),
    getTreasury: (id: string) => api.get(`/association-service/api/v1/associations/${id}/treasury`),
    getMeetings: (id: string) => api.get(`/association-service/api/v1/associations/${id}/meetings`),
    getLoans: (id: string) => api.get(`/association-service/api/v1/associations/${id}/loans`),

    // New endpoints
    getRoles: (id: string) => api.get(`/association-service/api/v1/associations/${id}/roles`),
    createRole: (id: string, data: any) => api.post(`/association-service/api/v1/associations/${id}/roles`, data),

    getApprovers: (id: string) => api.get(`/association-service/api/v1/associations/${id}/approvers`),
    setApprovers: (id: string, memberIds: string[]) =>
        api.post(`/association-service/api/v1/associations/${id}/approvers`, { member_ids: memberIds }),

    getPendingApprovals: (id: string) => api.get(`/association-service/api/v1/associations/${id}/approvals`),
    voteOnApproval: (requestId: string, vote: 'approve' | 'reject', comment?: string) =>
        api.post(`/association-service/api/v1/approvals/${requestId}/vote`, { vote, comment }),

    getMessages: (id: string, limit: number = 50) =>
        api.get(`/association-service/api/v1/associations/${id}/chat?limit=${limit}`),
    sendMessage: (id: string, content: string) =>
        api.post(`/association-service/api/v1/associations/${id}/chat`, { content }),

    getSolidarityEvents: (id: string) => api.get(`/association-service/api/v1/associations/${id}/solidarity`),
    createSolidarityEvent: (id: string, data: any) =>
        api.post(`/association-service/api/v1/associations/${id}/solidarity`, data),

    getCalledRounds: (id: string) => api.get(`/association-service/api/v1/associations/${id}/rounds`),
    createCalledRound: (id: string, beneficiaryId: string) =>
        api.post(`/association-service/api/v1/associations/${id}/rounds`, { beneficiary_id: beneficiaryId }),
    makePledge: (roundId: string, amount: number) =>
        api.post(`/association-service/api/v1/rounds/${roundId}/pledge`, { amount }),
    getPledges: (roundId: string) => api.get(`/association-service/api/v1/rounds/${roundId}/pledges`),
};

export default api;
