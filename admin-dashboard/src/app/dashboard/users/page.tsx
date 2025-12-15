'use client';

import { useEffect, useState } from 'react';
import { getUsers, blockUser, unblockUser } from '@/lib/api';
import { format } from 'date-fns';

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
            alert('Utilisateur bloqué');
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
            alert('Utilisateur débloqué');
        } catch (error) {
            alert('Erreur lors du déblocage');
        } finally {
            setActionLoading(null);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-8 flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900">Utilisateurs</h1>
                    <p className="text-slate-500 mt-1">{users.length} utilisateurs enregistrés</p>
                </div>
            </div>

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Utilisateur</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Email</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Téléphone</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">KYC</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Inscrit le</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {users.map((user) => (
                            <tr key={user.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-medium">
                                            {user.first_name?.[0]}{user.last_name?.[0]}
                                        </div>
                                        <div>
                                            <p className="font-medium text-slate-900">{user.first_name} {user.last_name}</p>
                                            <p className="text-sm text-slate-500 truncate max-w-[200px]">{user.id}</p>
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-slate-700">{user.email}</td>
                                <td className="px-6 py-4 text-slate-700">{user.phone || '-'}</td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getKYCBadge(user.kyc_level)}`}>
                                        Niveau {user.kyc_level}
                                    </span>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${user.is_active ? 'badge-success' : 'badge-danger'}`}>
                                        {user.is_active ? 'Actif' : 'Bloqué'}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {user.created_at ? format(new Date(user.created_at), 'dd/MM/yyyy') : '-'}
                                </td>
                                <td className="px-6 py-4 text-right">
                                    {user.is_active ? (
                                        <button
                                            onClick={() => handleBlock(user.id)}
                                            disabled={actionLoading === user.id}
                                            className="btn-danger text-sm px-3 py-1"
                                        >
                                            {actionLoading === user.id ? '...' : 'Bloquer'}
                                        </button>
                                    ) : (
                                        <button
                                            onClick={() => handleUnblock(user.id)}
                                            disabled={actionLoading === user.id}
                                            className="btn-primary text-sm px-3 py-1"
                                        >
                                            {actionLoading === user.id ? '...' : 'Débloquer'}
                                        </button>
                                    )}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {users.length === 0 && (
                    <div className="text-center py-12 text-slate-500">
                        Aucun utilisateur trouvé
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
        case 'pending':
            return 'badge-warning';
        default:
            return 'badge-warning';
    }
}
