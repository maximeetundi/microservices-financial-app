'use client';

import { useEffect, useState } from 'react';
import { getUsers, approveKYC, rejectKYC } from '@/lib/api';

interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone: string;
    kyc_level: number;
    kyc_status: string;
    is_active: boolean;
    created_at: string;
}

export default function KYCPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [allUsers, setAllUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [filter, setFilter] = useState('pending');

    const fetchUsers = async () => {
        try {
            setLoading(true);
            const response = await getUsers();
            const userList = response.data?.users || [];
            setAllUsers(userList);
        } catch (error) {
            console.error('Failed to fetch users:', error);
            setAllUsers([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    useEffect(() => {
        // Filter users based on selected filter
        if (filter === 'all') {
            setUsers(allUsers);
        } else {
            setUsers(allUsers.filter((u: User) => u.kyc_status === filter));
        }
    }, [filter, allUsers]);

    const handleApprove = async (userId: string) => {
        if (!confirm('Approuver cette demande KYC ?')) return;

        setActionLoading(userId);
        try {
            await approveKYC(userId, 'verified');
            await fetchUsers();
            alert('KYC approuvé avec succès');
        } catch (error) {
            console.error('Failed to approve KYC:', error);
            alert('Erreur lors de l\'approbation');
        } finally {
            setActionLoading(null);
        }
    };

    const handleReject = async (userId: string) => {
        const reason = prompt('Raison du rejet:');
        if (!reason) return;

        setActionLoading(userId);
        try {
            await rejectKYC(userId, reason);
            await fetchUsers();
            alert('KYC rejeté');
        } catch (error) {
            console.error('Failed to reject KYC:', error);
            alert('Erreur lors du rejet');
        } finally {
            setActionLoading(null);
        }
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'pending':
                return 'bg-yellow-100 text-yellow-800';
            case 'verified':
                return 'bg-green-100 text-green-800';
            case 'rejected':
                return 'bg-red-100 text-red-800';
            case 'none':
                return 'bg-gray-100 text-gray-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'pending': return 'En attente';
            case 'verified': return 'Vérifié';
            case 'rejected': return 'Rejeté';
            case 'none': return 'Non soumis';
            default: return status || 'Inconnu';
        }
    };

    const formatDate = (dateString: string) => {
        if (!dateString) return '-';
        try {
            return new Date(dateString).toLocaleDateString('fr-FR', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric'
            });
        } catch {
            return '-';
        }
    };

    const pendingCount = allUsers.filter(u => u.kyc_status === 'pending').length;
    const verifiedCount = allUsers.filter(u => u.kyc_status === 'verified').length;
    const rejectedCount = allUsers.filter(u => u.kyc_status === 'rejected').length;

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900">Vérification KYC</h1>
                    <p className="text-gray-500">Gérez les demandes de vérification d'identité</p>
                </div>
                <button
                    onClick={fetchUsers}
                    className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
                >
                    ↻ Actualiser
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                    <p className="text-gray-500 text-sm">Total utilisateurs</p>
                    <p className="text-2xl font-bold text-slate-900">{allUsers.length}</p>
                </div>
                <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm border-l-4 border-l-yellow-500">
                    <p className="text-gray-500 text-sm">En attente</p>
                    <p className="text-2xl font-bold text-yellow-600">{pendingCount}</p>
                </div>
                <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm border-l-4 border-l-green-500">
                    <p className="text-gray-500 text-sm">Vérifiés</p>
                    <p className="text-2xl font-bold text-green-600">{verifiedCount}</p>
                </div>
                <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm border-l-4 border-l-red-500">
                    <p className="text-gray-500 text-sm">Rejetés</p>
                    <p className="text-2xl font-bold text-red-600">{rejectedCount}</p>
                </div>
            </div>

            {/* Filters */}
            <div className="flex gap-2 flex-wrap">
                {[
                    { key: 'pending', label: `🕐 En attente (${pendingCount})` },
                    { key: 'verified', label: `✅ Vérifiés (${verifiedCount})` },
                    { key: 'rejected', label: `❌ Rejetés (${rejectedCount})` },
                    { key: 'all', label: 'Tous' }
                ].map(f => (
                    <button
                        key={f.key}
                        onClick={() => setFilter(f.key)}
                        className={`px-4 py-2 rounded-lg text-sm font-medium transition ${filter === f.key
                                ? 'bg-primary-600 text-white'
                                : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
                            }`}
                    >
                        {f.label}
                    </button>
                ))}
            </div>

            {/* Content */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                </div>
            ) : users.length === 0 ? (
                <div className="bg-white rounded-xl p-12 text-center border border-gray-200">
                    <p className="text-4xl mb-4">📋</p>
                    <p className="text-gray-500">Aucune demande {filter === 'pending' ? 'en attente' : ''}</p>
                </div>
            ) : (
                <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                    <table className="w-full">
                        <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                                <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Utilisateur</th>
                                <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Contact</th>
                                <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Niveau</th>
                                <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Statut</th>
                                <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Inscription</th>
                                <th className="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {users.map((user) => (
                                <tr key={user.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <div className={`w-10 h-10 rounded-full flex items-center justify-center font-medium ${user.kyc_status === 'pending' ? 'bg-yellow-100 text-yellow-600' :
                                                    user.kyc_status === 'verified' ? 'bg-green-100 text-green-600' :
                                                        'bg-gray-100 text-gray-600'
                                                }`}>
                                                {(user.first_name?.[0] || '')}{(user.last_name?.[0] || '')}
                                            </div>
                                            <div>
                                                <p className="font-medium text-slate-900">
                                                    {user.first_name} {user.last_name}
                                                </p>
                                                <p className="text-sm text-gray-500">{user.id?.slice(0, 8)}...</p>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <p className="text-sm text-slate-900">{user.email}</p>
                                        <p className="text-sm text-gray-500">{user.phone || '-'}</p>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                            Niveau {user.kyc_level || 0}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadge(user.kyc_status)}`}>
                                            {getStatusLabel(user.kyc_status)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">
                                        {formatDate(user.created_at)}
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        {user.kyc_status === 'pending' && (
                                            <div className="flex gap-2 justify-end">
                                                <button
                                                    onClick={() => handleApprove(user.id)}
                                                    disabled={actionLoading === user.id}
                                                    className="px-3 py-1.5 bg-green-600 text-white rounded-lg text-sm hover:bg-green-700 disabled:opacity-50"
                                                >
                                                    {actionLoading === user.id ? '...' : '✓ Approuver'}
                                                </button>
                                                <button
                                                    onClick={() => handleReject(user.id)}
                                                    disabled={actionLoading === user.id}
                                                    className="px-3 py-1.5 bg-red-600 text-white rounded-lg text-sm hover:bg-red-700 disabled:opacity-50"
                                                >
                                                    ✗ Rejeter
                                                </button>
                                            </div>
                                        )}
                                        {user.kyc_status === 'verified' && (
                                            <span className="text-green-600 text-sm">✓ Approuvé</span>
                                        )}
                                        {user.kyc_status === 'rejected' && (
                                            <button
                                                onClick={() => handleApprove(user.id)}
                                                disabled={actionLoading === user.id}
                                                className="px-3 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50"
                                            >
                                                Réexaminer
                                            </button>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}
