import { create } from 'zustand';

interface Admin {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    role: { name: string };
    is_active: boolean;
}

interface AuthState {
    admin: Admin | null;
    permissions: string[];
    isAuthenticated: boolean;
    setAdmin: (admin: Admin, permissions: string[]) => void;
    logout: () => void;
    hasPermission: (permission: string) => boolean;
    hasAnyPermission: (...permissions: string[]) => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
    admin: null,
    permissions: [],
    isAuthenticated: false,

    setAdmin: (admin, permissions) => set({
        admin,
        permissions,
        isAuthenticated: true
    }),

    logout: () => set({
        admin: null,
        permissions: [],
        isAuthenticated: false
    }),

    hasPermission: (permission) => {
        return get().permissions.includes(permission);
    },

    hasAnyPermission: (...permissions) => {
        const userPerms = get().permissions;
        return permissions.some(p => userPerms.includes(p));
    },
}));

// Permission constants
export const Permissions = {
    // Users
    VIEW_USERS: 'users.view',
    BLOCK_USERS: 'users.block',

    // KYC
    VIEW_KYC: 'kyc.view',
    APPROVE_KYC: 'kyc.approve',
    REJECT_KYC: 'kyc.reject',

    // Transactions
    VIEW_TRANSACTIONS: 'transactions.view',
    BLOCK_TRANSACTIONS: 'transactions.block',
    REFUND_TRANSACTIONS: 'transactions.refund',

    // Cards
    VIEW_CARDS: 'cards.view',
    FREEZE_CARDS: 'cards.freeze',
    BLOCK_CARDS: 'cards.block',

    // Wallets
    VIEW_WALLETS: 'wallets.view',
    FREEZE_WALLETS: 'wallets.freeze',

    // Admins
    VIEW_ADMINS: 'admins.view',
    CREATE_ADMINS: 'admins.create',
    UPDATE_ADMINS: 'admins.update',
    DELETE_ADMINS: 'admins.delete',

    // System
    VIEW_LOGS: 'system.logs',
    VIEW_ANALYTICS: 'analytics.view',
};
