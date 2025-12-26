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
}

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');

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

    const handleBlock = async (userId: string) => {
        const reason = prompt('Raison du blocage:');
        if (!reason) return;

        setActionLoading(userId);
        try {
            await blockUser(userId, reason);
            await fetchUsers();
            alert('Utilisateur bloqué avec succès');
        } catch (error) {
            alert('Erreur lors du blocage');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnblock = async (userId: string) => {
        setActionLoading(userId);
        try {
            await unblockUser(userId);
            await fetchUsers();
            alert('Utilisateur débloqué avec succès');
        } catch (error) {
            alert('Erreur lors du déblocage');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnlockPin = async (userId: string) => {
        if (!confirm('Êtes-vous sûr de vouloir débloquer le PIN de cet utilisateur ?')) return;

        setActionLoading(`pin-${userId}`);
        try {
            await unlockUserPin(userId);
            alert('PIN débloqué avec succès');
        } catch (error: any) {
            if (error.response?.status === 403) {
                alert('Accès refusé - Droits admin requis');
            } else if (error.response?.status === 404) {
                alert('Utilisateur non trouvé');
            } else {
                alert('Erreur lors du déblocage du PIN');
            }
        } finally {
            setActionLoading(null);
        }
    };

    // Filter users based on search
    const filteredUsers = users.filter(user =>
        user.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.first_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.last_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.phone?.includes(searchQuery)
    );

    // Compute stats
    const activeUsers = users.filter(u => u.is_active).length;
    const blockedUsers = users.filter(u => !u.is_active).length;
    const kycVerified = users.filter(u => String(u.kyc_level) === '3' || String(u.kyc_level) === 'verified').length;

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="w-16 h-16 rounded-full border-4 border-indigo-100 border-t-indigo-500 animate-spin"></div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
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
                    ↻ Actualiser
                </button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="stat-card-primary flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <UsersIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Total Utilisateurs</p>
                        <p className="text-2xl font-bold">{users.length}</p>
                    </div>
                </div>
                <div className="stat-card-success flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <ShieldCheckIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Comptes Actifs</p>
                        <p className="text-2xl font-bold">{activeUsers}</p>
                    </div>
                </div>
                <div className="stat-card-danger flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <NoSymbolIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Comptes Bloqués</p>
                        <p className="text-2xl font-bold">{blockedUsers}</p>
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

            {/* Search Bar */}
            <div className="relative">
                <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                <input
                    type="text"
                    placeholder="Rechercher par nom, email ou téléphone..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="input pl-12"
                />
            </div>

            {/* Users Table */}
            <div className="table-container">
                <table className="w-full">
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
                            <tr key={user.id} className="group">
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold shadow-md">
                                            {user.first_name?.[0]}{user.last_name?.[0]}
                                        </div>
                                        <div>
                                            <p className="font-semibold text-slate-900">
                                                {user.first_name} {user.last_name}
                                            </p>
                                            <p className="text-xs text-slate-400 truncate max-w-[180px] font-mono">
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
                                    <div className="flex items-center gap-2">
                                        <span className={`status-dot ${user.is_active ? 'status-dot-success' : 'status-dot-danger'}`}></span>
                                        <span className={user.is_active ? 'text-green-700' : 'text-red-700'}>
                                            {user.is_active ? 'Actif' : 'Bloqué'}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {user.created_at ? format(new Date(user.created_at), 'dd/MM/yyyy') : '—'}
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center justify-end gap-2 opacity-70 group-hover:opacity-100 transition-opacity">
                                        <button
                                            onClick={() => handleUnlockPin(user.id)}
                                            disabled={actionLoading === `pin-${user.id}`}
                                            className="btn-warning text-xs px-3 py-1.5"
                                            title="Débloquer le PIN"
                                        >
                                            {actionLoading === `pin-${user.id}` ? (
                                                <span className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                            ) : (
                                                '🔓 PIN'
                                            )}
                                        </button>
                                        {user.is_active ? (
                                            <button
                                                onClick={() => handleBlock(user.id)}
                                                disabled={actionLoading === user.id}
                                                className="btn-danger text-xs px-3 py-1.5"
                                            >
                                                {actionLoading === user.id ? (
                                                    <span className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                                ) : (
                                                    '⛔ Bloquer'
                                                )}
                                            </button>
                                        ) : (
                                            <button
                                                onClick={() => handleUnblock(user.id)}
                                                disabled={actionLoading === user.id}
                                                className="btn-success text-xs px-3 py-1.5"
                                            >
                                                {actionLoading === user.id ? (
                                                    <span className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                                ) : (
                                                    '✓ Débloquer'
                                                )}
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

