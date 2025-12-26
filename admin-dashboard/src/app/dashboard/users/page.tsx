'use client';

import { useEffect, useState } from 'react';
import { getUsers, blockUser, unblockUser, unlockUserPin } from '@/lib/api';
import { format } from 'date-fns';
import {
    UsersIcon,
    ShieldCheckIcon,
    NoSymbolIcon,
    KeyIcon,
    MagnifyingGlassIcon,
    XMarkIcon,
    ExclamationTriangleIcon,
    LockClosedIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    LockOpenIcon,
} from '@heroicons/react/24/outline';

interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone: string;
    is_active: boolean;
    kyc_level: number | string;
    created_at: string;
    pin_locked?: boolean;
}

// Toast Component
interface ToastProps {
    message: string;
    type: 'success' | 'error' | 'info';
    onClose: () => void;
}

function Toast({ message, type, onClose }: ToastProps) {
    useEffect(() => {
        const timer = setTimeout(onClose, 4000);
        return () => clearTimeout(timer);
    }, [onClose]);

    const colors = {
        success: 'from-green-500 to-emerald-600',
        error: 'from-red-500 to-pink-600',
        info: 'from-indigo-500 to-purple-600',
    };

    return (
        <div className={`fixed bottom-6 right-6 z-50 animate-slide-up`}>
            <div className={`flex items-center gap-3 px-5 py-4 bg-gradient-to-r ${colors[type]} text-white rounded-xl shadow-2xl`}>
                {type === 'success' && <CheckCircleIcon className="w-5 h-5" />}
                {type === 'error' && <ExclamationTriangleIcon className="w-5 h-5" />}
                <span className="font-medium">{message}</span>
                <button onClick={onClose} className="ml-2 hover:bg-white/20 p-1 rounded-lg transition-colors">
                    <XMarkIcon className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
}

// Modal Component
interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    icon?: React.ReactNode;
    iconColor?: string;
    children: React.ReactNode;
}

function Modal({ isOpen, onClose, title, icon, iconColor = 'bg-indigo-100 text-indigo-600', children }: ModalProps) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-black/50 backdrop-blur-sm animate-fade-in"
                onClick={onClose}
            />
            {/* Modal */}
            <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-md animate-slide-up overflow-hidden">
                {/* Header */}
                <div className="flex items-center gap-4 p-6 border-b border-gray-100">
                    {icon && (
                        <div className={`w-12 h-12 rounded-xl flex items-center justify-center ${iconColor}`}>
                            {icon}
                        </div>
                    )}
                    <h3 className="text-xl font-bold text-gray-900">{title}</h3>
                    <button
                        onClick={onClose}
                        className="ml-auto p-2 rounded-lg hover:bg-gray-100 transition-colors"
                    >
                        <XMarkIcon className="w-5 h-5 text-gray-500" />
                    </button>
                </div>
                {/* Content */}
                <div className="p-6">
                    {children}
                </div>
            </div>
        </div>
    );
}

type FilterTab = 'all' | 'active' | 'blocked' | 'pin_locked';

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState<FilterTab>('all');

    // Modal states
    const [blockModal, setBlockModal] = useState<{ isOpen: boolean; user: User | null }>({ isOpen: false, user: null });
    const [unblockModal, setUnblockModal] = useState<{ isOpen: boolean; user: User | null }>({ isOpen: false, user: null });
    const [unlockPinModal, setUnlockPinModal] = useState<{ isOpen: boolean; user: User | null }>({ isOpen: false, user: null });
    const [blockReason, setBlockReason] = useState('');

    // Toast state
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);

    const showToast = (message: string, type: 'success' | 'error' | 'info') => {
        setToast({ message, type });
    };

    const fetchUsers = async () => {
        try {
            const response = await getUsers();
            setUsers(response.data.users || []);
        } catch (error) {
            console.error('Failed to fetch users:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const handleBlock = async () => {
        if (!blockModal.user || !blockReason.trim()) return;

        setActionLoading(blockModal.user.id);
        try {
            await blockUser(blockModal.user.id, blockReason);
            await fetchUsers();
            setBlockModal({ isOpen: false, user: null });
            setBlockReason('');
            showToast(`Compte de ${blockModal.user.first_name} ${blockModal.user.last_name} bloqué avec succès`, 'success');
        } catch (error) {
            console.error('Block error:', error);
            showToast('Erreur lors du blocage du compte', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnblock = async () => {
        if (!unblockModal.user) return;

        setActionLoading(unblockModal.user.id);
        try {
            await unblockUser(unblockModal.user.id);
            await fetchUsers();
            showToast(`Compte de ${unblockModal.user.first_name} ${unblockModal.user.last_name} débloqué`, 'success');
            setUnblockModal({ isOpen: false, user: null });
        } catch (error) {
            console.error('Unblock error:', error);
            showToast('Erreur lors du déblocage du compte', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnlockPin = async () => {
        if (!unlockPinModal.user) return;

        setActionLoading(`pin-${unlockPinModal.user.id}`);
        try {
            await unlockUserPin(unlockPinModal.user.id);
            await fetchUsers();
            showToast(`PIN de ${unlockPinModal.user.first_name} ${unlockPinModal.user.last_name} débloqué`, 'success');
            setUnlockPinModal({ isOpen: false, user: null });
        } catch (error: any) {
            console.error('Unlock PIN error:', error);
            showToast('Erreur lors du déblocage du PIN', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    // Filter users based on search and tab
    const filteredUsers = users.filter(user => {
        const matchesSearch =
            user.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.first_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.last_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.phone?.includes(searchQuery);

        switch (activeTab) {
            case 'active':
                return matchesSearch && user.is_active;
            case 'blocked':
                return matchesSearch && !user.is_active;
            case 'pin_locked':
                return matchesSearch && user.pin_locked;
            default:
                return matchesSearch;
        }
    });

    // Compute stats
    const activeUsers = users.filter(u => u.is_active).length;
    const blockedUsers = users.filter(u => !u.is_active).length;
    const kycVerified = users.filter(u => String(u.kyc_level) === '3' || String(u.kyc_level) === 'verified').length;
    const pinLocked = users.filter(u => u.pin_locked).length;

    const tabs: { key: FilterTab; label: string; count: number; color: string }[] = [
        { key: 'all', label: 'Tous', count: users.length, color: 'indigo' },
        { key: 'active', label: 'Actifs', count: activeUsers, color: 'green' },
        { key: 'blocked', label: 'Bloqués', count: blockedUsers, color: 'red' },
        { key: 'pin_locked', label: 'PIN Bloqué', count: pinLocked, color: 'orange' },
    ];

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-center">
                    <div className="w-16 h-16 rounded-full border-4 border-indigo-100 border-t-indigo-500 animate-spin mx-auto"></div>
                    <p className="mt-4 text-gray-500 font-medium">Chargement des utilisateurs...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Toast */}
            {toast && (
                <Toast
                    message={toast.message}
                    type={toast.type}
                    onClose={() => setToast(null)}
                />
            )}

            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-slate-900 to-slate-600 bg-clip-text text-transparent">
                        Gestion des Utilisateurs
                    </h1>
                    <p className="text-slate-500 mt-1">
                        Gérez les comptes utilisateurs, vérifiez le KYC et débloquez les accès
                    </p>
                </div>
                <button
                    onClick={() => fetchUsers()}
                    className="btn-secondary flex items-center gap-2"
                >
                    <ArrowPathIcon className="w-4 h-4" />
                    Actualiser
                </button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
                <div className="stat-card-primary flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <UsersIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Total</p>
                        <p className="text-2xl font-bold">{users.length}</p>
                    </div>
                </div>
                <div className="stat-card-success flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <ShieldCheckIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Actifs</p>
                        <p className="text-2xl font-bold">{activeUsers}</p>
                    </div>
                </div>
                <div className="stat-card-danger flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <NoSymbolIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Bloqués</p>
                        <p className="text-2xl font-bold">{blockedUsers}</p>
                    </div>
                </div>
                <div className="stat-card-warning flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <LockClosedIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">PIN Bloqués</p>
                        <p className="text-2xl font-bold">{pinLocked}</p>
                    </div>
                </div>
                <div className="stat-card-info flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <KeyIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">KYC Vérifiés</p>
                        <p className="text-2xl font-bold">{kycVerified}</p>
                    </div>
                </div>
            </div>

            {/* Tabs & Search */}
            <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
                {/* Tabs */}
                <div className="flex gap-2 overflow-x-auto pb-2 lg:pb-0">
                    {tabs.map(tab => (
                        <button
                            key={tab.key}
                            onClick={() => setActiveTab(tab.key)}
                            className={`px-4 py-2.5 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${activeTab === tab.key
                                    ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                                    : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200 shadow-sm'
                                }`}
                        >
                            {tab.label}
                            <span className={`ml-2 px-2 py-0.5 rounded-full text-xs ${activeTab === tab.key ? 'bg-white/20' : 'bg-gray-100'
                                }`}>
                                {tab.count}
                            </span>
                        </button>
                    ))}
                </div>

                {/* Search Bar */}
                <div className="relative w-full lg:w-80">
                    <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                    <input
                        type="text"
                        placeholder="Rechercher..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="input pl-12"
                    />
                </div>
            </div>

            {/* Users Table */}
            <div className="table-container">
                <table className="w-full min-w-[900px]">
                    <thead>
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Utilisateur</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Email</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Téléphone</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">KYC</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Inscrit le</th>
                            <th className="text-right px-6 py-4 text-sm font-semibold text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {filteredUsers.map((user) => (
                            <tr key={user.id} className="group hover:bg-indigo-50/50 transition-colors">
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold shadow-md">
                                            {user.first_name?.[0]}{user.last_name?.[0]}
                                        </div>
                                        <div>
                                            <p className="font-semibold text-slate-900">
                                                {user.first_name} {user.last_name}
                                            </p>
                                            <p className="text-xs text-slate-400 truncate max-w-[150px] font-mono">
                                                {user.id}
                                            </p>
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-slate-600">{user.email}</td>
                                <td className="px-6 py-4 text-slate-600 font-mono text-sm">
                                    {user.phone || <span className="text-slate-300">—</span>}
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getKYCBadge(user.kyc_level)}`}>
                                        {getKYCLabel(user.kyc_level)}
                                    </span>
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex flex-col gap-1">
                                        <div className="flex items-center gap-2">
                                            <span className={`status-dot ${user.is_active ? 'status-dot-success' : 'status-dot-danger'}`}></span>
                                            <span className={user.is_active ? 'text-green-700 font-medium' : 'text-red-700 font-medium'}>
                                                {user.is_active ? 'Actif' : 'Bloqué'}
                                            </span>
                                        </div>
                                        {user.pin_locked && (
                                            <span className="text-xs text-orange-600 flex items-center gap-1">
                                                <LockClosedIcon className="w-3 h-3" /> PIN bloqué
                                            </span>
                                        )}
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {user.created_at ? format(new Date(user.created_at), 'dd/MM/yyyy') : '—'}
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center justify-end gap-2">
                                        <button
                                            onClick={() => setUnlockPinModal({ isOpen: true, user })}
                                            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-amber-700 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors"
                                            title="Débloquer le PIN"
                                        >
                                            <LockOpenIcon className="w-3.5 h-3.5" />
                                            PIN
                                        </button>
                                        {user.is_active ? (
                                            <button
                                                onClick={() => setBlockModal({ isOpen: true, user })}
                                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 hover:bg-red-100 rounded-lg transition-colors"
                                            >
                                                <NoSymbolIcon className="w-3.5 h-3.5" />
                                                Bloquer
                                            </button>
                                        ) : (
                                            <button
                                                onClick={() => setUnblockModal({ isOpen: true, user })}
                                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-green-700 bg-green-50 hover:bg-green-100 rounded-lg transition-colors"
                                            >
                                                <CheckCircleIcon className="w-3.5 h-3.5" />
                                                Débloquer
                                            </button>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {filteredUsers.length === 0 && (
                    <div className="text-center py-16">
                        <UsersIcon className="w-12 h-12 text-slate-300 mx-auto mb-4" />
                        <p className="text-slate-500 font-medium">Aucun utilisateur trouvé</p>
                        <p className="text-slate-400 text-sm mt-1">
                            {searchQuery ? 'Essayez une autre recherche' : 'Les utilisateurs apparaîtront ici'}
                        </p>
                    </div>
                )}
            </div>

            {/* Block User Modal */}
            <Modal
                isOpen={blockModal.isOpen}
                onClose={() => { setBlockModal({ isOpen: false, user: null }); setBlockReason(''); }}
                title="Bloquer l'utilisateur"
                icon={<ExclamationTriangleIcon className="w-6 h-6" />}
                iconColor="bg-red-100 text-red-600"
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-red-50 rounded-xl border border-red-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {blockModal.user?.first_name?.[0]}{blockModal.user?.last_name?.[0]}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{blockModal.user?.first_name} {blockModal.user?.last_name}</p>
                            <p className="text-sm text-gray-500">{blockModal.user?.email}</p>
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Raison du blocage <span className="text-red-500">*</span>
                        </label>
                        <textarea
                            value={blockReason}
                            onChange={(e) => setBlockReason(e.target.value)}
                            placeholder="Décrivez la raison du blocage..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-red-500 focus:ring-2 focus:ring-red-500/20 outline-none transition-all resize-none"
                        />
                    </div>
                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => { setBlockModal({ isOpen: false, user: null }); setBlockReason(''); }}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleBlock}
                            disabled={!blockReason.trim() || actionLoading === blockModal.user?.id}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-red-500/25"
                        >
                            {actionLoading === blockModal.user?.id ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Blocage...
                                </span>
                            ) : 'Confirmer le blocage'}
                        </button>
                    </div>
                </div>
            </Modal>

            {/* Unblock User Modal */}
            <Modal
                isOpen={unblockModal.isOpen}
                onClose={() => setUnblockModal({ isOpen: false, user: null })}
                title="Débloquer l'utilisateur"
                icon={<ShieldCheckIcon className="w-6 h-6" />}
                iconColor="bg-green-100 text-green-600"
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-green-50 rounded-xl border border-green-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {unblockModal.user?.first_name?.[0]}{unblockModal.user?.last_name?.[0]}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{unblockModal.user?.first_name} {unblockModal.user?.last_name}</p>
                            <p className="text-sm text-gray-500">{unblockModal.user?.email}</p>
                        </div>
                    </div>
                    <p className="text-sm text-gray-600 bg-gray-50 p-4 rounded-xl">
                        L'utilisateur pourra de nouveau accéder à son compte et effectuer des transactions.
                    </p>
                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => setUnblockModal({ isOpen: false, user: null })}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleUnblock}
                            disabled={actionLoading === unblockModal.user?.id}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-green-500/25"
                        >
                            {actionLoading === unblockModal.user?.id ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Déblocage...
                                </span>
                            ) : 'Confirmer le déblocage'}
                        </button>
                    </div>
                </div>
            </Modal>

            {/* Unlock PIN Modal */}
            <Modal
                isOpen={unlockPinModal.isOpen}
                onClose={() => setUnlockPinModal({ isOpen: false, user: null })}
                title="Débloquer le PIN"
                icon={<KeyIcon className="w-6 h-6" />}
                iconColor="bg-amber-100 text-amber-600"
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-amber-50 rounded-xl border border-amber-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {unlockPinModal.user?.first_name?.[0]}{unlockPinModal.user?.last_name?.[0]}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{unlockPinModal.user?.first_name} {unlockPinModal.user?.last_name}</p>
                            <p className="text-sm text-gray-500">{unlockPinModal.user?.email}</p>
                        </div>
                    </div>
                    <div className="bg-amber-50 border border-amber-100 p-4 rounded-xl">
                        <div className="flex items-start gap-3">
                            <ExclamationTriangleIcon className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                            <div className="text-sm text-amber-800">
                                <p className="font-medium">Information</p>
                                <p className="mt-1">Le PIN peut être bloqué après trop de tentatives incorrectes. Cette action réinitialise le compteur de tentatives.</p>
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => setUnlockPinModal({ isOpen: false, user: null })}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleUnlockPin}
                            disabled={actionLoading === `pin-${unlockPinModal.user?.id}`}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-amber-500/25"
                        >
                            {actionLoading === `pin-${unlockPinModal.user?.id}` ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Déblocage...
                                </span>
                            ) : 'Débloquer le PIN'}
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

function getKYCBadge(level: number | string): string {
    switch (String(level)) {
        case '3':
        case 'verified':
            return 'badge-success';
        case '2':
            return 'badge-info';
        case '1':
            return 'badge-warning';
        case 'pending':
            return 'badge-warning';
        default:
            return 'badge-gray';
    }
}

function getKYCLabel(level: number | string): string {
    switch (String(level)) {
        case '3':
        case 'verified':
            return '✓ Vérifié';
        case '2':
            return 'Niveau 2';
        case '1':
            return 'Niveau 1';
        case 'pending':
            return 'En attente';
        case '0':
            return 'Non vérifié';
        default:
            return `Niveau ${level}`;
    }
}
